package models

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

// Cell represents a single cell within a notebook.
type Cell struct {
	ID             uuid.UUID      `json:"id"`
	NotebookID     uuid.UUID      `json:"notebook_id"`
	CellIndex      int            `json:"cell_index"`
	CellName       sql.NullString `json:"cell_name"`
	CellType       string         `json:"cell_type"`
	Source         string         `json:"source"`
	ExecutionCount int            `json:"execution_count"`
	Outputs        []CellOutput   `json:"outputs,omitempty"`
	EvolutionRuns  []EvolutionRun `json:"evolution_runs,omitempty"`
}

// CellOutput represents the output of a cell execution.
type CellOutput struct {
	ID             uuid.UUID       `json:"id"`
	CellID         uuid.UUID       `json:"cell_id"`
	OutputIndex    int             `json:"output_index"`
	Type           string          `json:"type"`
	DataJSON       json.RawMessage `json:"data_json"`
	MinioURL       string          `json:"minio_url"`
	ExecutionCount int             `json:"execution_count"`
}

// CreateCellRequest defines the structure for a request to create a new cell.
type CreateCellRequest struct {
	NotebookID uuid.UUID `json:"notebook_id" binding:"required"`
	CellIndex  int       `json:"cell_index" binding:"required"`
	CellName   string    `json:"cell_name"`
	CellType   string    `json:"cell_type" binding:"required"`
	Source     string    `json:"source"`
}

// UpdateCellRequest defines the structure for a request to update a cell.
type UpdateCellRequest struct {
	CellIndex      *int    `json:"cell_index,omitempty"`
	CellName       *string `json:"cell_name,omitempty"`
	CellType       *string `json:"cell_type,omitempty"`
	Source         *string `json:"source,omitempty"`
	ExecutionCount *int    `json:"execution_count,omitempty"`
}

// UpdateCellsRequest defines the structure for a bulk cell update request.
type UpdateCellsRequest struct {
	UpdatedOrder  []uuid.UUID                  `json:"updated_order"`
	CellsToDelete []uuid.UUID                  `json:"cells_to_delete"`
	CellsToUpsert map[uuid.UUID]CellDataForUpsert `json:"cells_to_upsert"`
}

// CellDataForUpsert represents the data for a cell to be upserted.
type CellDataForUpsert struct {
	CellType       string         `json:"cell_type"`
	Source         string         `json:"source"`
	CellName       sql.NullString `json:"cell_name"`
	ExecutionCount int            `json:"execution_count"`
}

// CreateCellOutputRequest defines the structure for a request to create a new cell output.
type CreateCellOutputRequest struct {
	CellID      uuid.UUID       `json:"cell_id"`
	OutputIndex int             `json:"output_index"`
	Type        string          `json:"type"`
	DataJSON    json.RawMessage `json:"data_json"`
	MinioURL    string          `json:"minio_url"`
}
