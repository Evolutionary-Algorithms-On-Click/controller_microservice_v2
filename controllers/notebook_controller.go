package controllers

import (
	"time"
	"encoding/json"
	"context"
	"fmt"
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg" // Added pkg import
	"github.com/rs/zerolog"
)

// NotebookController handles notebook related endpoints.
type NotebookController struct {
	NotebookModule *modules.NotebookModule
	Logger         *zerolog.Logger
}

func NewNotebookController(module *modules.NotebookModule, logger *zerolog.Logger) *NotebookController {
	return &NotebookController{
		NotebookModule: module,
		Logger:         logger,
	}
}



// CreateNotebookHandler handles POST /api/v1/notebooks
func (c *NotebookController) CreateNotebookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req models.CreateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	nb, err := c.NotebookModule.CreateNotebook(ctx, &req)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to create notebook")
		http.Error(w, fmt.Sprintf("error creating notebook: %v", err), http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, nb, c.Logger)
}

// ListNotebooksHandler handles GET /api/v1/notebooks
func (c *NotebookController) ListNotebooksHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// TODO: parse limit/offset/created_by/problem_statement_id from r.URL.Query()
	nbs, err := c.NotebookModule.ListNotebooks(ctx, nil)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to list notebooks")
		http.Error(w, "error listing notebooks", http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, nbs, c.Logger)
}

// GetNotebookByIDHandler handles GET /api/v1/notebooks/{id}
func (c *NotebookController) GetNotebookByIDHandler(w http.ResponseWriter, r *http.Request) {
	notebookID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	nb, err := c.NotebookModule.GetNotebookByID(ctx, notebookID)
	if err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("get notebook failed")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, nb, c.Logger)
}

// UpdateNotebookByIDHandler handles PUT /api/v1/notebooks/{id}
func (c *NotebookController) UpdateNotebookByIDHandler(w http.ResponseWriter, r *http.Request) {
	notebookID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req models.UpdateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updated, err := c.NotebookModule.UpdateNotebook(ctx, notebookID, &req)
	if err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("update notebook failed")
		http.Error(w, "error updating notebook", http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, updated, c.Logger)
}

// DeleteNotebookByIDHandler handles DELETE /api/v1/notebooks/{id}
func (c *NotebookController) DeleteNotebookByIDHandler(w http.ResponseWriter, r *http.Request) {
	notebookID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := c.NotebookModule.DeleteNotebook(ctx, notebookID); err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("delete notebook failed")
		http.Error(w, "error deleting notebook", http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusNoContent, nil, c.Logger)
}