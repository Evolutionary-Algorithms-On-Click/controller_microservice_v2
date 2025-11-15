package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user's session in the database.
type Session struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	NotebookID      uuid.UUID `json:"notebook_id"`
	CurrentKernelID uuid.UUID `json:"current_kernel_id"`
	Status          string    `json:"status"`
	LastActiveAt    time.Time `json:"last_active_at"`
}

// CreateSessionRequest is the struct for the request body to create a new session.
type CreateSessionRequest struct {
	NotebookID string `json:"notebook_id" binding:"required"`
	Language string `json:"language" binding:"required"`
}
