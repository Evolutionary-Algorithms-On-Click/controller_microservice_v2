package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

// CellRepository defines the data access methods for a cell, ensuring ownership.
type CellRepository interface {
	CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error)
	GetCellByID(ctx context.Context, id uuid.UUID, userID string) (*models.Cell, error)
	GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error)
	UpdateCell(ctx context.Context, cell *models.Cell, userID string) (*models.Cell, error)
	DeleteCell(ctx context.Context, id uuid.UUID, userID string) error
	UpdateCells(ctx context.Context, notebookID uuid.UUID, req *models.UpdateCellsRequest) error

	CreateCellOutput(ctx context.Context, output *models.CellOutput) (*models.CellOutput, error)
	GetCellOutputsByCellID(ctx context.Context, cellID uuid.UUID) ([]*models.CellOutput, error)
	GetCellOutputByID(ctx context.Context, outputID uuid.UUID, userID string) (*models.CellOutput, error)
	DeleteCellOutput(ctx context.Context, id uuid.UUID, userID string) error
}

type cellRepository struct {
	db     *pgxpool.Pool
	Logger zerolog.Logger
}

func NewCellRepository(db *pgxpool.Pool, logger zerolog.Logger) CellRepository {
	return &cellRepository{db: db, Logger: logger}
}

func (r *cellRepository) CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	query := `
		INSERT INTO cells (id, notebook_id, cell_index, cell_name, cell_type, source, execution_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, notebook_id, cell_index, cell_name, cell_type, source, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		cell.ID.ToUUID(),
		cell.NotebookID,
		cell.CellIndex,
		cell.CellName,
		cell.CellType,
		cell.Source,
		cell.ExecutionCount,
	)

	var createdCell models.Cell
	var id uuid.UUID
	err := row.Scan(
		&id,
		&createdCell.NotebookID,
		&createdCell.CellIndex,
		&createdCell.CellName,
		&createdCell.CellType,
		&createdCell.Source,
		&createdCell.ExecutionCount,
	)
	if err != nil {
		return nil, err
	}
	createdCell.ID = models.StringUUID(id)

	return &createdCell, nil
}

func (r *cellRepository) GetCellByID(
	ctx context.Context,
	id uuid.UUID,
	userID string,
) (*models.Cell, error) {
	query := `
		SELECT c.id, c.notebook_id, c.cell_index, c.cell_name, c.cell_type, c.source, c.execution_count
		FROM cells c
		JOIN notebooks n ON c.notebook_id = n.id
		JOIN problem_statements ps ON n.problem_statement_id = ps.id
		WHERE c.id = $1 AND ps.created_by = $2;
	`
	row := r.db.QueryRow(ctx, query, id, userID)

	var cell models.Cell
	var scannedID uuid.UUID
	err := row.Scan(
		&scannedID,
		&cell.NotebookID,
		&cell.CellIndex,
		&cell.CellName,
		&cell.CellType,
		&cell.Source,
		&cell.ExecutionCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("cell not found or not owned by user")
		}
		return nil, err
	}
	cell.ID = models.StringUUID(scannedID)

	return &cell, nil
}

func (r *cellRepository) GetCellsByNotebookID(
	ctx context.Context,
	notebookID uuid.UUID,
) ([]*models.Cell, error) {
	// Ownership check is expected to happen in the controller/module before this call
	query := `
		SELECT id, notebook_id, cell_index, cell_name, cell_type, source, execution_count
		FROM cells
		WHERE notebook_id = $1
		ORDER BY cell_index;
	`
	rows, err := r.db.Query(ctx, query, notebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cells []*models.Cell
	for rows.Next() {
		var cell models.Cell
		var scannedID uuid.UUID
		err := rows.Scan(
			&scannedID,
			&cell.NotebookID,
			&cell.CellIndex,
			&cell.CellName,
			&cell.CellType,
			&cell.Source,
			&cell.ExecutionCount,
		)
		if err != nil {
			return nil, err
		}
		cell.ID = models.StringUUID(scannedID)
		cells = append(cells, &cell)
	}

	return cells, nil
}

func (r *cellRepository) UpdateCell(ctx context.Context, cell *models.Cell, userID string) (*models.Cell, error) {
	query := `
		UPDATE cells
		SET cell_index = $2, cell_name = $3, cell_type = $4, source = $5, execution_count = $6
		WHERE id = $1 AND notebook_id IN (
			SELECT n.id FROM notebooks n
			JOIN problem_statements ps ON n.problem_statement_id = ps.id
			WHERE ps.created_by = $7
		)
		RETURNING id, notebook_id, cell_index, cell_name, cell_type, source, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		cell.ID.ToUUID(),
		cell.CellIndex,
		cell.CellName,
		cell.CellType,
		cell.Source,
		cell.ExecutionCount,
		userID,
	)

	var updatedCell models.Cell
	var scannedID uuid.UUID
	err := row.Scan(
		&scannedID,
		&updatedCell.NotebookID,
		&updatedCell.CellIndex,
		&updatedCell.CellName,
		&updatedCell.CellType,
		&updatedCell.Source,
		&updatedCell.ExecutionCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("cell not found or not owned by user")
		}
		return nil, err
	}
	updatedCell.ID = models.StringUUID(scannedID)

	return &updatedCell, nil
}

func (r *cellRepository) DeleteCell(ctx context.Context, id uuid.UUID, userID string) error {
	query := `
		DELETE FROM cells
		WHERE id = $1 AND notebook_id IN (
			SELECT n.id FROM notebooks n
			JOIN problem_statements ps ON n.problem_statement_id = ps.id
			WHERE ps.created_by = $2
		);
	`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("cell not found or not owned by user")
	}
	return nil
}

func (r *cellRepository) UpdateCells(ctx context.Context, notebookID uuid.UUID, req *models.UpdateCellsRequest) error {
	r.Logger.Info().
		Str("notebook_id", notebookID.String()).
		Int("delete_count", len(req.CellsToDelete)).
		Int("upsert_count", len(req.CellsToUpsert)).
		Int("order_count", len(req.UpdatedOrder)).
		Msg("Updating cells in repository")

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// update notebook requirements if provided
	if req.Requirements != nil {
		r.Logger.Info().
			Str("notebook_id", notebookID.String()).
			Msg("Updating notebook requirements")

		_, err := tx.Exec(ctx, "UPDATE notebooks SET requirements = $1 WHERE id = $2", *req.Requirements, notebookID)
		if err != nil {
			r.Logger.Error().Err(err).Msg("Failed to update notebook requirements")
			return err
		}
	}

	// fetch existing cell ids to find orphans
	rows, err := tx.Query(ctx, "SELECT id FROM cells WHERE notebook_id = $1", notebookID)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Failed to fetch existing cell IDs")
		return err
	}
	defer rows.Close()

	existingIDs := make(map[string]bool)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return err
		}
		existingIDs[id.String()] = true
	}

	// map order for quick lookup and identify cells to keep
	orderMap := make(map[string]int)
	keepIDs := make(map[string]bool)
	for i, id := range req.UpdatedOrder {
		idStr := id.ToUUID().String()
		orderMap[idStr] = i
		keepIDs[idStr] = true
	}

	// determine cells to delete: explicitly requested + orphans
	idsToDeleteMap := make(map[uuid.UUID]bool)
	for _, id := range req.CellsToDelete {
		idsToDeleteMap[id.ToUUID()] = true
	}
	for idStr := range existingIDs {
		if !keepIDs[idStr] {
			u, _ := uuid.Parse(idStr)
			idsToDeleteMap[u] = true
		}
	}

	// delete all orphaned and explicitly removed cells
	if len(idsToDeleteMap) > 0 {
		deleteSlice := make([]uuid.UUID, 0, len(idsToDeleteMap))
		for id := range idsToDeleteMap {
			deleteSlice = append(deleteSlice, id)
		}
		if _, err := tx.Exec(ctx, "DELETE FROM cells WHERE id = ANY($1)", deleteSlice); err != nil {
			r.Logger.Error().Err(err).Msg("Failed to delete cells")
			return err
		}
	}

	upsertedIDs := make(map[string]bool)

	// handle upserts for modified or new cells
	if len(req.CellsToUpsert) > 0 {
		for idStr, cellData := range req.CellsToUpsert {
			upsertedIDs[idStr] = true
			cellUUID, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}

			cellIndex, ok := orderMap[idStr]
			if !ok {
				// fallback to end of list if order is missing
				cellIndex = len(req.UpdatedOrder)
				r.Logger.Warn().Str("cell_id", idStr).Msg("Cell ID missing from order list, appending to end")
			}

			var nullCellName sql.NullString
			if cellData.CellName != nil {
				nullCellName.String = *cellData.CellName
				nullCellName.Valid = true
			}

			query := `
                INSERT INTO cells (id, notebook_id, cell_type, source, cell_name, execution_count, cell_index)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (id) DO UPDATE
                SET source = $4, cell_name = $5, execution_count = $6, cell_index = $7;
            `
			if _, err := tx.Exec(ctx, query, cellUUID, notebookID, cellData.CellType, cellData.Source, nullCellName, cellData.ExecutionCount, cellIndex); err != nil {
				r.Logger.Error().Err(err).Str("cell_id", idStr).Msg("Failed to upsert cell")
				return err
			}
		}
	}

	// update indices for cells that were only reordered
	for idStr, index := range orderMap {
		if !upsertedIDs[idStr] {
			cellUUID, _ := uuid.Parse(idStr)
			if _, err := tx.Exec(ctx, "UPDATE cells SET cell_index = $1 WHERE id = $2", index, cellUUID); err != nil {
				r.Logger.Error().Err(err).Str("cell_id", idStr).Msg("Failed to reorder cell")
				return err
			}
		}
	}

	return tx.Commit(ctx)
}


func (r *cellRepository) CreateCellOutput(ctx context.Context, output *models.CellOutput) (*models.CellOutput, error) {
	r.Logger.Debug().Str("output_id", output.ID.String()).Str("cell_id", output.CellID.ToUUID().String()).Msg("CellRepository: Creating cell output")
	query := `
		INSERT INTO cell_outputs (id, cell_id, output_index, type, data_json, minio_url, execution_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, cell_id, output_index, type, data_json, minio_url, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		output.ID,
		output.CellID.ToUUID(),
		output.OutputIndex,
		output.Type,
		output.DataJSON,
		output.MinioURL,
		output.ExecutionCount,
	)

	var createdOutput models.CellOutput
	var cellID uuid.UUID
	err := row.Scan(
		&createdOutput.ID,
		&cellID,
		&createdOutput.OutputIndex,
		&createdOutput.Type,
		&createdOutput.DataJSON,
		&createdOutput.MinioURL,
		&createdOutput.ExecutionCount,
	)
	if err != nil {
		r.Logger.Error().Err(err).Msg("CellRepository: Failed to scan created cell output")
		return nil, err
	}
	createdOutput.CellID = models.StringUUID(cellID)
	r.Logger.Info().Str("created_output_id", createdOutput.ID.String()).Msg("CellRepository: Successfully created cell output")
	return &createdOutput, nil
}

func (r *cellRepository) GetCellOutputsByCellID(ctx context.Context, cellID uuid.UUID) ([]*models.CellOutput, error) {
	// Ownership check is expected to happen in the controller/module before this call
	query := `
		SELECT id, cell_id, output_index, type, data_json, minio_url, execution_count
		FROM cell_outputs
		WHERE cell_id = $1
		ORDER BY output_index;
	`
	rows, err := r.db.Query(ctx, query, cellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outputs []*models.CellOutput
	for rows.Next() {
		var output models.CellOutput
		var cellID uuid.UUID
		err := rows.Scan(
			&output.ID,
			&cellID,
			&output.OutputIndex,
			&output.Type,
			&output.DataJSON,
			&output.MinioURL,
			&output.ExecutionCount,
		)
		if err != nil {
			return nil, err
		}
		output.CellID = models.StringUUID(cellID)
		outputs = append(outputs, &output)
	}

	return outputs, nil
}

func (r *cellRepository) GetCellOutputByID(ctx context.Context, outputID uuid.UUID, userID string) (*models.CellOutput, error) {
	query := `
		SELECT co.id, co.cell_id, co.output_index, co.type, co.data_json, co.minio_url, co.execution_count
		FROM cell_outputs co
		JOIN cells c ON co.cell_id = c.id
		JOIN notebooks n ON c.notebook_id = n.id
		JOIN problem_statements ps ON n.problem_statement_id = ps.id
		WHERE co.id = $1 AND ps.created_by = $2;
	`
	row := r.db.QueryRow(ctx, query, outputID, userID)

	var output models.CellOutput
	var cellID uuid.UUID
	err := row.Scan(
		&output.ID,
		&cellID,
		&output.OutputIndex,
		&output.Type,
		&output.DataJSON,
		&output.MinioURL,
		&output.ExecutionCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("output not found or not owned by user")
		}
		return nil, err
	}
	output.CellID = models.StringUUID(cellID)
	return &output, nil
}

func (r *cellRepository) DeleteCellOutput(ctx context.Context, id uuid.UUID, userID string) error {
	query := `
		DELETE FROM cell_outputs
		WHERE id = $1 AND cell_id IN (
			SELECT c.id FROM cells c
			JOIN notebooks n ON c.notebook_id = n.id
			JOIN problem_statements ps ON n.problem_statement_id = ps.id
			WHERE ps.created_by = $2
		);
	`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("output not found or not owned by user")
	}
	return nil
}
