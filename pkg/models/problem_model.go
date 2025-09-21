package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProblemStatement represents a problem statement in the database.
type ProblemStatement struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title"`
	DescriptionJSON    []byte    `json:"description_json"`
	CreatedBy          uuid.UUID `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
}

// CreateProblemRequest is the struct for the request body to create a new problem.
type CreateProblemRequest struct {
	Title           string          `json:"title" binding:"required"`
	DescriptionJSON json.RawMessage `json:"description_json" binding:"required"`
}