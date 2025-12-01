package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NotebookRepository interface {
	CreateNotebook(ctx context.Context, req *models.CreateNotebookRequest) (*models.Notebook, error)
	ListNotebooks(ctx context.Context, filters map[string]string) ([]models.Notebook, error)
	GetNotebookByID(ctx context.Context, id string) (*models.Notebook, error)
	UpdateNotebook(ctx context.Context, id string, req *models.UpdateNotebookRequest) (*models.Notebook, error)
	DeleteNotebook(ctx context.Context, id string) error
	SaveNotebookCells(ctx context.Context, notebookID string, req *models.SaveCellsRequest) error
}

type notebookRepository struct {
	pool *pgxpool.Pool
}

func NewNotebookRepository(pool *pgxpool.Pool) NotebookRepository {
	return &notebookRepository{
		pool: pool,
	}
}

func (r *notebookRepository) CreateNotebook(ctx context.Context, req *models.CreateNotebookRequest) (*models.Notebook, error) {
	id := uuid.New().String()
	now := time.Now().UTC()

	query := `
		INSERT INTO notebooks (id, title, context_minio_url, problem_statement_id, created_at, last_modified_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, context_minio_url, problem_statement_id, created_at, last_modified_at;
	`

	row := r.pool.QueryRow(ctx, query,
		id,
		req.Title,
		nil, // TODO: Should include logic for context minIO url
		req.ProblemStatementID,
		now,
		now,
	)

	var nb models.Notebook
	if err := row.Scan(
		&nb.ID,
		&nb.Title,
		&nb.ContextMinioURL,
		&nb.ProblemStatementID,
		&nb.CreatedAt,
		&nb.LastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &nb, nil
}

func (r *notebookRepository) ListNotebooks(ctx context.Context, filters map[string]string) ([]models.Notebook, error) {
	query := `SELECT id, title, context_minio_url, problem_statement_id, created_at, last_modified_at FROM notebooks`
	args := []any{}
	where := ""

	if filters != nil {
		i := 1
		for key, val := range filters {
			if i == 1 {
				where += " WHERE "
			} else {
				where += " AND "
			}
			where += key + " = $" + strconv.Itoa(i)
			args = append(args, val)
			i++
		}
	}

	rows, err := r.pool.Query(ctx, query+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebooks []models.Notebook
	for rows.Next() {
		var nb models.Notebook
		if err := rows.Scan(
			&nb.ID,
			&nb.Title,
			&nb.ContextMinioURL,
			&nb.ProblemStatementID,
			&nb.CreatedAt,
			&nb.LastModifiedAt,
		); err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}

	return notebooks, nil
}

func (r *notebookRepository) GetNotebookByID(ctx context.Context, id string) (*models.Notebook, error) {
	notebookUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid notebook ID format")
	}

	query := `
		SELECT
			n.id, n.title, n.context_minio_url, n.problem_statement_id, n.created_at, n.last_modified_at,
			c.id, c.notebook_id, c.cell_index, c.cell_type, c.source, c.execution_count,
			co.id, co.cell_id, co.output_index, co.type, co.data_json, co.minio_url, co.execution_count,
			er.id, er.source_cell_id, er.start_time, er.end_time, er.status,
			cv.id, cv.evolution_run_id, cv.code, cv.metric, cv.is_best, cv.generation, cv.parent_variant_id
		FROM
			notebooks n
		LEFT JOIN
			cells c ON n.id = c.notebook_id
		LEFT JOIN
			cell_outputs co ON c.id = co.cell_id
		LEFT JOIN
			evolution_runs er ON c.id = er.source_cell_id
		LEFT JOIN
			cell_variations cv ON er.id = cv.evolution_run_id
		WHERE
			n.id = $1
		ORDER BY
			c.cell_index ASC, 
			co.output_index ASC, 
			er.start_time ASC, 
			cv.generation ASC;
	`

	rows, err := r.pool.Query(ctx, query, notebookUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebook models.Notebook
	cellMap := make(map[uuid.UUID]*models.Cell)
	runMap := make(map[uuid.UUID]*models.EvolutionRun)
	var orderedCellIDs []uuid.UUID

	for rows.Next() {
		var (
			cellID            uuid.NullUUID
			cellNotebookID    uuid.NullUUID
			cellIndex         sql.NullInt32
			cellType          sql.NullString
			cellSource        sql.NullString
			cellExecCount     sql.NullInt32
			outputID          uuid.NullUUID
			outputCellID      uuid.NullUUID
			outputIndex       sql.NullInt32
			outputType        sql.NullString
			outputDataJSON    []byte
			outputMinioURL    sql.NullString
			outputExecCount   sql.NullInt32
			erID              uuid.NullUUID
			erSourceCellID    uuid.NullUUID
			erStartTime       sql.NullTime
			erEndTime         sql.NullTime
			erStatus          sql.NullString
			cvID              uuid.NullUUID
			cvEvolutionRunID  uuid.NullUUID
			cvCode            sql.NullString
			cvMetric          sql.NullFloat64
			cvIsBest          sql.NullBool
			cvGeneration      sql.NullInt32
			cvParentVariantID uuid.NullUUID
		)

		if err := rows.Scan(
			&notebook.ID, &notebook.Title, &notebook.ContextMinioURL, &notebook.ProblemStatementID, &notebook.CreatedAt, &notebook.LastModifiedAt,
			&cellID, &cellNotebookID, &cellIndex, &cellType, &cellSource, &cellExecCount,
			&outputID, &outputCellID, &outputIndex, &outputType, &outputDataJSON, &outputMinioURL, &outputExecCount,
			&erID, &erSourceCellID, &erStartTime, &erEndTime, &erStatus,
			&cvID, &cvEvolutionRunID, &cvCode, &cvMetric, &cvIsBest, &cvGeneration, &cvParentVariantID,
		); err != nil {
			return nil, err
		}

		if cellID.Valid {
			if _, exists := cellMap[cellID.UUID]; !exists {
				cellMap[cellID.UUID] = &models.Cell{
					ID:             cellID.UUID,
					NotebookID:     cellNotebookID.UUID,
					CellIndex:      int(cellIndex.Int32),
					CellType:       cellType.String,
					Source:         cellSource.String,
					ExecutionCount: int(cellExecCount.Int32),
					Outputs:        []models.CellOutput{},
					EvolutionRuns:  []models.EvolutionRun{},
				}
				orderedCellIDs = append(orderedCellIDs, cellID.UUID)
			}
		}

		if outputID.Valid {
			if cell, exists := cellMap[outputCellID.UUID]; exists {
				cell.Outputs = append(cell.Outputs, models.CellOutput{
					ID:             outputID.UUID,
					CellID:         outputCellID.UUID,
					OutputIndex:    int(outputIndex.Int32),
					Type:           outputType.String,
					DataJSON:       outputDataJSON,
					MinioURL:       outputMinioURL.String,
					ExecutionCount: int(outputExecCount.Int32),
				})
			}
		}

		if erID.Valid {
			if _, exists := runMap[erID.UUID]; !exists {
				run := models.EvolutionRun{
					ID:           erID.UUID,
					SourceCellID: erSourceCellID.UUID,
					StartTime:    erStartTime.Time,
					Status:       erStatus.String,
					Variations:   []models.CellVariation{},
				}
				if erEndTime.Valid {
					run.EndTime = &erEndTime.Time
				}
				runMap[erID.UUID] = &run
			}
		}

		if cvID.Valid {
			if run, exists := runMap[cvEvolutionRunID.UUID]; exists {
				variation := models.CellVariation{
					ID:             cvID.UUID,
					EvolutionRunID: cvEvolutionRunID.UUID,
					Code:           cvCode.String,
					Metric:         cvMetric.Float64,
					IsBest:         cvIsBest.Bool,
					Generation:     int(cvGeneration.Int32),
				}
				if cvParentVariantID.Valid {
					variation.ParentVariantID = &cvParentVariantID.UUID
				}
				run.Variations = append(run.Variations, variation)
			}
		}
	}

	if notebook.ID == "" {
		return nil, errors.New("notebook not found")
	}

	for _, run := range runMap {
		if cell, exists := cellMap[run.SourceCellID]; exists {
			cell.EvolutionRuns = append(cell.EvolutionRuns, *run)
		}
	}

	notebook.Cells = make([]models.Cell, len(orderedCellIDs))
	for i, cellID := range orderedCellIDs {
		notebook.Cells[i] = *cellMap[cellID]
	}

	return &notebook, nil
}

func (r *notebookRepository) UpdateNotebook(ctx context.Context, id string, req *models.UpdateNotebookRequest) (*models.Notebook, error) {
	setClause := ""
	args := []any{}
	argIndex := 1

	if req.Title != nil {
		setClause += "title = $" + strconv.Itoa(argIndex)
		args = append(args, *req.Title)
		argIndex++
	}

	if setClause == "" {
		return r.GetNotebookByID(ctx, id)
	}

	setClause += ", last_modified_at = $" + strconv.Itoa(argIndex)
	args = append(args, time.Now().UTC())
	argIndex++

	args = append(args, id)

	query := `
		UPDATE notebooks
		SET ` + setClause + `
		WHERE id = $` + strconv.Itoa(argIndex) + `
		RETURNING id, title, context_minio_url, problem_statement_id, created_at, last_modified_at;
	`

	row := r.pool.QueryRow(ctx, query, args...)

	var nb models.Notebook
	if err := row.Scan(
		&nb.ID,
		&nb.Title,
		&nb.ContextMinioURL,
		&nb.ProblemStatementID,
		&nb.CreatedAt,
		&nb.LastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &nb, nil
}

func (r *notebookRepository) DeleteNotebook(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM notebooks WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("notebook not found")
	}
	return nil
}

func (r *notebookRepository) SaveNotebookCells(ctx context.Context, notebookID string, req *models.SaveCellsRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() {
		// If commit already ran successfully, rollback will return ErrTxClosed — that’s fine.
		_ = tx.Rollback(ctx)
	}()

	var changed bool = false

	if len(req.CellsToDelete) > 0 {
		changed = true
		deleteQuery := "DELETE FROM cells WHERE id = ANY($1)"
		if _, err := tx.Exec(ctx, deleteQuery, req.CellsToDelete); err != nil {
			return fmt.Errorf("could not delete cells: %w", err)
		}
	}

	batch := &pgx.Batch{}

	if len(req.CellsToUpsert) > 0 {
		changed = true
		upsertQuery := `
			INSERT INTO cells (id, notebook_id, cell_type, source)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				cell_type = EXCLUDED.cell_type,
				source = EXCLUDED.source;
		`
		for idStr, cellData := range req.CellsToUpsert {
			cellUUID, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("invalid UUID in CellsToUpsert: %s", idStr)
			}
			batch.Queue(upsertQuery, cellUUID, notebookID, cellData.CellType, cellData.Source)
		}
	}

	if len(req.UpdatedOrder) > 0 {
		changed = true
		updateIndexQuery := `UPDATE cells SET cell_index = $1 WHERE id = $2`
		for i, cellIDStr := range req.UpdatedOrder {
			cellUUID, err := uuid.Parse(cellIDStr)
			if err != nil {
				return fmt.Errorf("invalid UUID in UpdatedOrder: %s", cellIDStr)
			}
			batch.Queue(updateIndexQuery, i, cellUUID)
		}
	}

	if batch.Len() > 0 {
		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return fmt.Errorf("failed to execute batch for saving cells: %w", err)
		}
	}

	if changed {
		updateTimeQuery := "UPDATE notebooks SET last_modified_at = $1 WHERE id = $2"
		if _, err := tx.Exec(ctx, updateTimeQuery, time.Now().UTC(), notebookID); err != nil {
			return fmt.Errorf("could not update notebook timestamp: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

