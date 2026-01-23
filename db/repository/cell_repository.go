package repository

import (
	"context"
	"database/sql"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

// CellRepository defines the data access methods for a cell.
type CellRepository interface {
	CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error)
	GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error)
	GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error)
	UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error)
	DeleteCell(ctx context.Context, id uuid.UUID) error
	UpdateCells(ctx context.Context, notebookID uuid.UUID, req *models.UpdateCellsRequest) error

	CreateCellOutput(ctx context.Context, output *models.CellOutput) (*models.CellOutput, error)
	GetCellOutputsByCellID(ctx context.Context, cellID uuid.UUID) ([]*models.CellOutput, error)
	DeleteCellOutput(ctx context.Context, id uuid.UUID) error
}

type cellRepository struct {
	db *pgxpool.Pool
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

func (r *cellRepository) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	query := `
		SELECT id, notebook_id, cell_index, cell_name, cell_type, source, execution_count
		FROM cells
		WHERE id = $1;
	`
	row := r.db.QueryRow(ctx, query, id)

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
		return nil, err
	}
	cell.ID = models.StringUUID(scannedID)

	return &cell, nil
}

func (r *cellRepository) GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error) {
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

func (r *cellRepository) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	query := `
		UPDATE cells
		SET cell_index = $2, cell_name = $3, cell_type = $4, source = $5, execution_count = $6
		WHERE id = $1
		RETURNING id, notebook_id, cell_index, cell_name, cell_type, source, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		cell.ID.ToUUID(),
		cell.CellIndex,
		cell.CellName,
		cell.CellType,
		cell.Source,
		cell.ExecutionCount,
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
		return nil, err
	}
	updatedCell.ID = models.StringUUID(scannedID)

	return &updatedCell, nil
}

func (r *cellRepository) DeleteCell(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cells WHERE id = $1", id)
	return err
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

	// 1. Handle Deletions
	if len(req.CellsToDelete) > 0 {
		cellsToDeleteUUIDs := make([]uuid.UUID, len(req.CellsToDelete))
		for i, id := range req.CellsToDelete {
			cellsToDeleteUUIDs[i] = id.ToUUID()
		}
		r.Logger.Info().Interface("cells_to_delete", cellsToDeleteUUIDs).Msg("Deleting cells")
		if _, err := tx.Exec(ctx, "DELETE FROM cells WHERE id = ANY($1)", cellsToDeleteUUIDs); err != nil {
			r.Logger.Error().Err(err).Msg("Failed to delete cells")
			return err // Rollback will be called by defer
		}
	}

	// Create a map for quick index lookup, now using string keys for the map
	orderMap := make(map[string]int) // Changed key type
	for i, id := range req.UpdatedOrder {
		orderMap[id.ToUUID().String()] = i // Convert StringUUID to string for map key
	}

	// 2. Handle Upserts and Reordering in one go
	if len(req.CellsToUpsert) > 0 {
		r.Logger.Info().Interface("cells_to_upsert", req.CellsToUpsert).Msg("Upserting cells")
		for idStr, cellData := range req.CellsToUpsert { // idStr is string
			cellUUID, err := uuid.Parse(idStr) // Parse string ID to UUID
			if err != nil {
				r.Logger.Error().Err(err).Str("cell_id_string", idStr).Msg("Failed to parse cell ID string to UUID")
				return err
			}

			cellIndex, ok := orderMap[idStr] // Lookup using string ID
			if !ok {
				cellIndex = 9999
				r.Logger.Warn().Str("cell_id", idStr).Msg("Cell ID found in upsert map but not in updated order. Assigning high index.")
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
			r.Logger.Info().
				Str("cell_id", idStr).
				Int("cell_index", cellIndex).
				Msg("Executing upsert for cell")
			if _, err := tx.Exec(ctx, query, cellUUID, notebookID, cellData.CellType, cellData.Source, nullCellName, cellData.ExecutionCount, cellIndex); err != nil {
				r.Logger.Error().Err(err).Str("cell_id", idStr).Msg("Failed to upsert cell")
				return err // Rollback
			}
		}
	} else if len(req.UpdatedOrder) > 0 {
		// This block handles reordering-only updates
		r.Logger.Info().Msg("Performing reorder-only update")
		for i, cellID := range req.UpdatedOrder {
			if _, err := tx.Exec(ctx, "UPDATE cells SET cell_index = $1 WHERE id = $2", i, cellID.ToUUID()); err != nil {
				r.Logger.Error().Err(err).Str("cell_id", cellID.ToUUID().String()).Msg("Failed to reorder cell")
				return err
			}
		}
	}

	r.Logger.Info().Msg("Committing transaction")
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

func (r *cellRepository) DeleteCellOutput(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cell_outputs WHERE id = $1", id)
	return err
}
