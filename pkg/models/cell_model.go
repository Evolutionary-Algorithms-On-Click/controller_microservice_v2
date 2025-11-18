package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Cell represents a single cell within a notebook.
type Cell struct {
	ID             uuid.UUID `json:"id"`
	NotebookID     uuid.UUID `json:"notebook_id"`
	CellIndex      int       `json:"cell_index"`
	CellType       string    `json:"cell_type"`
	Source         string    `json:"source"`
	ExecutionCount int       `json:"execution_count"`
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
	NotebookID uuid.UUID `json:"notebook_id"`
	CellIndex  int       `json:"cell_index"`
	CellType   string    `json:"cell_type"`
	Source     string    `json:"source"`
}

// UpdateCellRequest defines the structure for a request to update a cell.
type UpdateCellRequest struct {
	CellIndex      *int    `json:"cell_index,omitempty"`
	CellType       *string `json:"cell_type,omitempty"`
	Source         *string `json:"source,omitempty"`
	ExecutionCount *int    `json:"execution_count,omitempty"`
}

// CreateCellOutputRequest defines the structure for a request to create a new cell output.
type CreateCellOutputRequest struct {
	CellID      uuid.UUID       `json:"cell_id"`
	OutputIndex int             `json:"output_index"`
	Type        string          `json:"type"`
	DataJSON    json.RawMessage `json:"data_json"`
	MinioURL    string          `json:"minio_url"`
}
