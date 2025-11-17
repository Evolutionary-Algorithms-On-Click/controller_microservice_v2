package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg" // Added pkg import
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// SessionController holds the dependencies for the session handlers.
type SessionController struct {
	Module *modules.SessionModule
	Logger zerolog.Logger
}

// NewSessionController creates and returns a new SessionController.
func NewSessionController(module *modules.SessionModule, logger zerolog.Logger) *SessionController {
	return &SessionController{
		Module: module,
		Logger: logger,
	}
}

// CreateSessionHandler handles POST /api/v1/sessions
func (c *SessionController) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second) // Increased timeout for kernel start
	defer cancel()

	// userID, ok := r.Context().Value("userID").(string)
	// if !ok || userID == "" {
	// 	c.Logger.Error().Msg("userID not found in context after authentication")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	userID := "123e4567-e89b-12d3-a456-426614174000" // Hardcoded for testing

	var req models.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	session, err := c.Module.CreateSession(ctx, userID, req.NotebookID, req.Language)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to create session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, session, &c.Logger)
}
// ListSessionsHandler handles GET /api/v1/sessions
func (c *SessionController) ListSessionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// userID, ok := r.Context().Value("userID").(string)
	// if !ok || userID == "" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	userIDStr := "123e4567-e89b-12d3-a456-426614174000" // Hardcoded for testing

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	sessions, err := c.Module.ListSessions(ctx, userID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to list sessions")
		http.Error(w, "failed to list sessions", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, sessions, &c.Logger)
}

// GetSessionByIDHandler handles GET /api/v1/sessions/{id}
func (c *SessionController) GetSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid session ID format", http.StatusBadRequest)
		return
	}

	// userID, ok := r.Context().Value("userID").(string)
	// if !ok || userID == "" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	userIDStr := "123e4567-e89b-12d3-a456-426614174000" // Hardcoded for testing

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	session, err := c.Module.GetSessionByID(ctx, id, userID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get session by id")
		// TODO: Check for pgx.ErrNoRows and return 404
		http.Error(w, "failed to get session", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, session, &c.Logger)
}

// UpdateSessionByIDHandler handles PUT /api/v1/sessions/{id}
func (c *SessionController) UpdateSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid session ID format", http.StatusBadRequest)
		return
	}

	// userID, ok := r.Context().Value("userID").(string)
	// if !ok || userID == "" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	userIDStr := "123e4567-e89b-12d3-a456-426614174000" // Hardcoded for testing

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// TODO: Validate the status value (e.g., "active", "closed")

	session, err := c.Module.UpdateSessionStatus(ctx, id, userID, req.Status)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to update session status")
		// TODO: Check for pgx.ErrNoRows and return 404
		http.Error(w, "failed to update session", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, session, &c.Logger)
}

// DeleteSessionByIDHandler handles DELETE /api/v1/sessions/{id}
func (c *SessionController) DeleteSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid session ID format", http.StatusBadRequest)
		return
	}

	// userID, ok := r.Context().Value("userID").(string)
	// if !ok || userID == "" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	userIDStr := "123e4567-e89b-12d3-a456-426614174000" // Hardcoded for testing

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	err = c.Module.DeleteSession(ctx, id, userID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to delete session")
		// TODO: Check for pgx.ErrNoRows and return 404
		http.Error(w, "failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
