package repository

import (
	"context"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CellRepository defines the data access methods for a cell.
type CellRepository interface {
	CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error)
	GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error)
	GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error)
	UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error)
	DeleteCell(ctx context.Context, id uuid.UUID) error

	CreateCellOutput(ctx context.Context, output *models.CellOutput) (*models.CellOutput, error)
	GetCellOutputsByCellID(ctx context.Context, cellID uuid.UUID) ([]*models.CellOutput, error)
	DeleteCellOutput(ctx context.Context, id uuid.UUID) error
}

type cellRepository struct {
	db *pgxpool.Pool
}

func NewCellRepository(db *pgxpool.Pool) CellRepository {
	return &cellRepository{db: db}
}

func (r *cellRepository) CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	query := `
		INSERT INTO cells (id, notebook_id, cell_index, cell_type, source, execution_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, notebook_id, cell_index, cell_type, source, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		cell.ID,
		cell.NotebookID,
		cell.CellIndex,
		cell.CellType,
		cell.Source,
		cell.ExecutionCount,
	)

	var createdCell models.Cell
	err := row.Scan(
		&createdCell.ID,
		&createdCell.NotebookID,
		&createdCell.CellIndex,
		&createdCell.CellType,
		&createdCell.Source,
		&createdCell.ExecutionCount,
	)
	if err != nil {
		return nil, err
	}

	return &createdCell, nil
}

func (r *cellRepository) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	query := `
		SELECT id, notebook_id, cell_index, cell_type, source, execution_count
		FROM cells
		WHERE id = $1;
	`
	row := r.db.QueryRow(ctx, query, id)

	var cell models.Cell
	err := row.Scan(
		&cell.ID,
		&cell.NotebookID,
		&cell.CellIndex,
		&cell.CellType,
		&cell.Source,
		&cell.ExecutionCount,
	)
	if err != nil {
		return nil, err
	}

	return &cell, nil
}

func (r *cellRepository) GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error) {
	query := `
		SELECT id, notebook_id, cell_index, cell_type, source, execution_count
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
		err := rows.Scan(
			&cell.ID,
			&cell.NotebookID,
			&cell.CellIndex,
			&cell.CellType,
			&cell.Source,
			&cell.ExecutionCount,
		)
		if err != nil {
			return nil, err
		}
		cells = append(cells, &cell)
	}

	return cells, nil
}

func (r *cellRepository) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	query := `
		UPDATE cells
		SET cell_index = $2, cell_type = $3, source = $4, execution_count = $5
		WHERE id = $1
		RETURNING id, notebook_id, cell_index, cell_type, source, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		cell.ID,
		cell.CellIndex,
		cell.CellType,
		cell.Source,
		cell.ExecutionCount,
	)

	var updatedCell models.Cell
	err := row.Scan(
		&updatedCell.ID,
		&updatedCell.NotebookID,
		&updatedCell.CellIndex,
		&updatedCell.CellType,
		&updatedCell.Source,
		&updatedCell.ExecutionCount,
	)
	if err != nil {
		return nil, err
	}

	return &updatedCell, nil
}

func (r *cellRepository) DeleteCell(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cells WHERE id = $1", id)
	return err
}

func (r *cellRepository) CreateCellOutput(ctx context.Context, output *models.CellOutput) (*models.CellOutput, error) {
	query := `
		INSERT INTO cell_outputs (id, cell_id, output_index, type, data_json, minio_url, execution_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, cell_id, output_index, type, data_json, minio_url, execution_count;
	`
	row := r.db.QueryRow(ctx, query,
		output.ID,
		output.CellID,
		output.OutputIndex,
		output.Type,
		output.DataJSON,
		output.MinioURL,
		output.ExecutionCount,
	)

	var createdOutput models.CellOutput
	err := row.Scan(
		&createdOutput.ID,
		&createdOutput.CellID,
		&createdOutput.OutputIndex,
		&createdOutput.Type,
		&createdOutput.DataJSON,
		&createdOutput.MinioURL,
		&createdOutput.ExecutionCount,
	)
	if err != nil {
		return nil, err
	}

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
		err := rows.Scan(
			&output.ID,
			&output.CellID,
			&output.OutputIndex,
			&output.Type,
			&output.DataJSON,
			&output.MinioURL,
			&output.ExecutionCount,
		)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, &output)
	}

	return outputs, nil
}

func (r *cellRepository) DeleteCellOutput(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cell_outputs WHERE id = $1", id)
	return err
}
