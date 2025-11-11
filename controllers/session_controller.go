package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
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

	var req models.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	session, err := c.Module.CreateSession(ctx, req.NotebookID, req.Language)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to create session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponseWithLogger(w, http.StatusCreated, session, &c.Logger)
}

// ListSessionsHandler handles GET /api/v1/sessions
func (c *SessionController) ListSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to list all sessions.
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// GetSessionByIDHandler handles GET /api/v1/sessions/{id}
func (c *SessionController) GetSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to get a session by its ID.
	http.Error(w, "Not Implemented: get by id "+id, http.StatusNotImplemented)
}

// UpdateSessionByIDHandler handles PUT /api/v1/sessions/{id}
func (c *SessionController) UpdateSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to update a session.
	http.Error(w, "Not Implemented: update by id "+id, http.StatusNotImplemented)
}

// DeleteSessionByIDHandler handles DELETE /api/v1/sessions/{id}
func (c *SessionController) DeleteSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to delete a session.
	http.Error(w, "Not Implemented: delete by id "+id, http.StatusNotImplemented)
}
