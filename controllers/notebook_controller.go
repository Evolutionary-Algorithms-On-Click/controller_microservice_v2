package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
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

// helper
func writeJSONResponseWithLogger(w http.ResponseWriter, status int, v interface{}, logger *zerolog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			logger.Error().Err(err).Msg("failed to encode response")
		}
	}
}

// NotebookListHandler handles POST / GET on /api/notebooks
func (c *NotebookController) NotebookListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodPost:
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
		writeJSONResponseWithLogger(w, http.StatusCreated, nb, c.Logger)

	case http.MethodGet:
		// handle list, optional query params
		// TODO: parse limit/offset/created_by/problem_statement_id
		nbs, err := c.NotebookModule.ListNotebooks(ctx, nil)
		if err != nil {
			c.Logger.Error().Err(err).Msg("failed to list notebooks")
			http.Error(w, "error listing notebooks", http.StatusInternalServerError)
			return
		}
		writeJSONResponseWithLogger(w, http.StatusOK, nbs, c.Logger)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// NotebookByIDHandler handles GET/PUT/DELETE on /api/notebooks/{id}
func (c *NotebookController) NotebookByIDHandler(w http.ResponseWriter, r *http.Request) {
	// expecting URL: /api/notebooks/{id} (exact)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	notebookID := parts[2]

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		nb, err := c.NotebookModule.GetNotebookByID(ctx, notebookID)
		if err != nil {
			c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("get notebook failed")
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSONResponseWithLogger(w, http.StatusOK, nb, c.Logger)

	case http.MethodPut:
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
		writeJSONResponseWithLogger(w, http.StatusOK, updated, c.Logger)

	case http.MethodDelete:
		if err := c.NotebookModule.DeleteNotebook(ctx, notebookID); err != nil {
			c.Logger.Error().Err(err).Str("notebook_id", notebookID).Msg("delete notebook failed")
			http.Error(w, "error deleting notebook", http.StatusInternalServerError)
			return
		}
		writeJSONResponseWithLogger(w, http.StatusNoContent, nil, c.Logger)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
