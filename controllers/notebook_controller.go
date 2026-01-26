package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/middleware"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg" // Added pkg import
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
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

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("userID not found in context for notebook creation")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	nb, err := c.NotebookModule.CreateNotebook(ctx, &req, user.ID)
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

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("userID not found in context for listing notebooks")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	filters := make(map[string]string)
	query := r.URL.Query()

	if createdBy := query.Get("created_by"); createdBy != "" {
		filters["created_by"] = createdBy
	}
	if problemStatementID := query.Get("problem_statement_id"); problemStatementID != "" {
		filters["problem_statement_id"] = problemStatementID
	}
	if limit := query.Get("limit"); limit != "" {
		filters["limit"] = limit
	}
	if offset := query.Get("offset"); offset != "" {
		filters["offset"] = offset
	}

	nbs, err := c.NotebookModule.ListNotebooks(ctx, filters, user.ID)
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

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("userID not found in context for getting notebook by ID")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	nb, err := c.NotebookModule.GetNotebookByID(ctx, notebookID, user.ID)
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

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("userID not found in context for updating notebook by ID")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.UpdateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updated, err := c.NotebookModule.UpdateNotebook(ctx, notebookID, &req, user.ID)
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

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("userID not found in context for deleting notebook by ID")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	if err := c.NotebookModule.DeleteNotebook(ctx, notebookID, user.ID); err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("delete notebook failed")
		http.Error(w, "error deleting notebook", http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusNoContent, nil, c.Logger)
}
