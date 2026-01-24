package models

import (
	"database/sql"
	"time"
)

// Notebook represents the notebooks table.
type Notebook struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	ContextMinioURL    *string        `json:"context_minio_url,omitempty"`
	Requirements       sql.NullString `json:"requirements,omitempty"`
	ProblemStatementID *string        `json:"problem_statement_id,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	LastModifiedAt     time.Time      `json:"last_modified_at"`
	Cells              []Cell         `json:"cells,omitempty"`
}

// CreateNotebookRequest is the payload to create a notebook.
type CreateNotebookRequest struct {
	Title              string  `json:"title" binding:"required"`
	Requirements       *string `json:"requirements,omitempty"`
	ProblemStatementID *string `json:"problem_statement_id,omitempty"`
}

// UpdateNotebookRequest defines updatable fields.
type UpdateNotebookRequest struct {
	Title              *string `json:"title,omitempty"`
	Requirements       *string `json:"requirements,omitempty"`
	ProblemStatementID *string `json:"problem_statement_id,omitempty"`
}
