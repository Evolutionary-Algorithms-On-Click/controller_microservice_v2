package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Thanus-Kumaar/controller_microservice_v2/middleware"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	CellTypeMarkdown = "markdown"
	CellTypeCode     = "code"
	CellTypeRaw      = "raw"
)

const (
	OutputTypeStream        = "stream"
	OutputTypeDisplayData   = "display_data"
	OutputTypeExecuteResult = "execute_result"
	OutputTypeError         = "error"
)

var validCellTypes = map[string]struct{}{
	CellTypeMarkdown: {},
	CellTypeCode:     {},
	CellTypeRaw:      {},
}

var validOutputTypes = map[string]struct{}{
	OutputTypeStream:        {},
	OutputTypeDisplayData:   {},
	OutputTypeExecuteResult: {},
	OutputTypeError:         {},
}

// CellController holds the dependencies for the cell handlers.
type CellController struct {
	Module         *modules.CellModule
	Logger         zerolog.Logger
	NotebookModule *modules.NotebookModule // Added NotebookModule
}

// NewCellController creates and returns a new CellController.
func NewCellController(module *modules.CellModule, logger zerolog.Logger, notebookModule *modules.NotebookModule) *CellController {
	return &CellController{
		Module:         module,
		Logger:         logger,
		NotebookModule: notebookModule,
	}
}

func (c *CellController) UpdateCellsHandler(w http.ResponseWriter, r *http.Request) {
	notebookIDStr := r.PathValue("notebook_id")
	notebookID, err := uuid.Parse(notebookIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "Invalid notebook ID"},
			&c.Logger,
		)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for cell update")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Verify ownership of the notebook
	_, err = c.NotebookModule.GetNotebookByID(r.Context(), notebookIDStr, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookIDStr).Msg("Notebook not found or not owned by user for cell update")
		http.Error(w, "Notebook not found or not owned by user", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to read request body")
		pkg.WriteJSONResponseWithLogger(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to read request body"},
			&c.Logger,
		)
		return
	}
	c.Logger.Info().RawJSON("raw_request_body", body).Msg("Received request to update cells")

	var req models.UpdateCellsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to decode request body")
		pkg.WriteJSONResponseWithLogger(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request body"},
			&c.Logger,
		)
		return
	}

	c.Logger.Info().Interface("decoded_request", req).Msg("Decoded request body")

	if err := c.Module.UpdateCells(r.Context(), notebookID, &req, user.ID); err != nil { // Pass userID to module
		c.Logger.Error().Err(err).Msg("Failed to update cells")
		pkg.WriteJSONResponseWithLogger(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to update cells"},
			&c.Logger,
		)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, map[string]string{"status": "success"}, &c.Logger)
}
func (c *CellController) CreateCellHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for cell creation")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateCellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, &c.Logger)
		return
	}

	// Verify ownership of the parent notebook
	notebookIDStr := req.NotebookID.String()
	_, err := c.NotebookModule.GetNotebookByID(r.Context(), notebookIDStr, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookIDStr).Msg("Notebook not found or not owned by user for cell creation")
		http.Error(w, "Notebook not found or not owned by user", http.StatusNotFound)
		return
	}

	if _, ok := validCellTypes[req.CellType]; !ok {
		allowedTypes := make([]string, 0, len(validCellTypes))
		for k := range validCellTypes {
			allowedTypes = append(allowedTypes, k)
		}
		err_msg := fmt.Sprintf("Invalid cell type: '%s'. Allowed types are: %s", req.CellType, strings.Join(allowedTypes, ", "))
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": err_msg}, &c.Logger)
		return
	}

	cell, err := c.Module.CreateCell(r.Context(), &req, user.ID) // Pass userID to module
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to create cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create cell"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, cell, &c.Logger)
}

func (c *CellController) GetCellsByNotebookIDHandler(w http.ResponseWriter, r *http.Request) {
	notebookIDStr := r.PathValue("notebook_id")
	notebookID, err := uuid.Parse(notebookIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid notebook ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for getting cells by notebook ID")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Verify ownership of the parent notebook
	_, err = c.NotebookModule.GetNotebookByID(r.Context(), notebookIDStr, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("notebook_id", notebookIDStr).Msg("Notebook not found or not owned by user for getting cells")
		http.Error(w, "Notebook not found or not owned by user", http.StatusNotFound)
		return
	}

	cells, err := c.Module.GetCellsByNotebookID(r.Context(), notebookID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to get cells")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get cells"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, cells, &c.Logger)
}

func (c *CellController) GetCellByIDHandler(w http.ResponseWriter, r *http.Request) {
	cellIDStr := r.PathValue("cell_id")
	cellID, err := uuid.Parse(cellIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid cell ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for getting cell by ID")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	cell, err := c.Module.GetCellByID(r.Context(), cellID, user.ID) // Pass userID to module
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to get cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusNotFound, map[string]string{"error": "Cell not found or not owned by user"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, cell, &c.Logger)
}

func (c *CellController) UpdateCellHandler(w http.ResponseWriter, r *http.Request) {
	cellIDStr := r.PathValue("cell_id")
	cellID, err := uuid.Parse(cellIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid cell ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for updating cell")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Retrieve the cell to get its notebook ID and implicitly check ownership via module call
	_, err = c.Module.GetCellByID(r.Context(), cellID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("cell_id", cellIDStr).Msg("Failed to retrieve existing cell for update or not owned by user")
		pkg.WriteJSONResponseWithLogger(w, http.StatusNotFound, map[string]string{"error": "Cell not found or not owned by user"}, &c.Logger)
		return
	}

	var req models.UpdateCellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, &c.Logger)
		return
	}

	if req.CellType != nil {
		if _, ok := validCellTypes[*req.CellType]; !ok {
			allowedTypes := make([]string, 0, len(validCellTypes))
			for k := range validCellTypes {
				allowedTypes = append(allowedTypes, k)
			}
			err_msg := fmt.Sprintf("Invalid cell type: '%s'. Allowed types are: %s", *req.CellType, strings.Join(allowedTypes, ", "))
			pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": err_msg}, &c.Logger)
			return
		}
	}

	cell, err := c.Module.UpdateCell(r.Context(), cellID, &req, user.ID) // Pass userID to module
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to update cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update cell"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, cell, &c.Logger)
}

func (c *CellController) DeleteCellHandler(w http.ResponseWriter, r *http.Request) {
	cellIDStr := r.PathValue("cell_id")
	cellID, err := uuid.Parse(cellIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid cell ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for deleting cell")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Verify ownership before deleting
	_, err = c.Module.GetCellByID(r.Context(), cellID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("cell_id", cellIDStr).Msg("Failed to retrieve existing cell for delete or not owned by user")
		pkg.WriteJSONResponseWithLogger(w, http.StatusNotFound, map[string]string{"error": "Cell not found or not owned by user"}, &c.Logger)
		return
	}

	if err := c.Module.DeleteCell(r.Context(), cellID, user.ID); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to delete cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete cell"}, &c.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CellController) CreateCellOutputHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for creating cell output")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateCellOutputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Logger.Error().Err(err).Msg("CreateCellOutputHandler: Invalid request body")
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, &c.Logger)
		return
	}
	c.Logger.Debug().Interface("request_body", req).Msg("CreateCellOutputHandler: Decoded request body")

	// Verify ownership of the cell before creating an output for it
	_, err := c.Module.GetCellByID(r.Context(), req.CellID.ToUUID(), user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("cell_id", req.CellID.ToUUID().String()).Msg("Cell not found or not owned by user for creating output")
		http.Error(w, "Cell not found or not owned by user", http.StatusNotFound)
		return
	}

	if _, ok := validOutputTypes[req.Type]; !ok {
		allowedTypes := make([]string, 0, len(validOutputTypes))
		for k := range validOutputTypes {
			allowedTypes = append(allowedTypes, k)
		}
		err_msg := fmt.Sprintf("Invalid output type: '%s'. Allowed types are: %s", req.Type, strings.Join(allowedTypes, ", "))
		c.Logger.Error().Str("invalid_output_type", req.Type).Msg("CreateCellOutputHandler: Invalid output type")
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": err_msg}, &c.Logger)
		return
	}
	c.Logger.Debug().Str("cell_id", req.CellID.ToUUID().String()).Msg("CreateCellOutputHandler: Calling module to create cell output")
	output, err := c.Module.CreateCellOutput(r.Context(), &req) // Ownership already checked
	if err != nil {
		c.Logger.Error().Err(err).Msg("CreateCellOutputHandler: Failed to create cell output")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create cell output"}, &c.Logger)
		return
	}
	c.Logger.Info().Str("created_output_id", output.ID.String()).Msg("CreateCellOutputHandler: Successfully created cell output")
	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, output, &c.Logger)
}

func (c *CellController) GetCellOutputsByCellIDHandler(w http.ResponseWriter, r *http.Request) {
	cellIDStr := r.PathValue("cell_id")
	cellID, err := uuid.Parse(cellIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid cell ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for getting cell outputs")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Verify ownership of the cell before getting outputs
	_, err = c.Module.GetCellByID(r.Context(), cellID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("cell_id", cellID.String()).Msg("Cell not found or not owned by user for getting outputs")
		http.Error(w, "Cell not found or not owned by user", http.StatusNotFound)
		return
	}

	outputs, err := c.Module.GetCellOutputsByCellID(r.Context(), cellID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to get cell outputs")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get cell outputs"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, outputs, &c.Logger)
}

func (c *CellController) DeleteCellOutputHandler(w http.ResponseWriter, r *http.Request) {
	outputIDStr := r.PathValue("output_id")
	outputID, err := uuid.Parse(outputIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid output ID"}, &c.Logger)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(*middleware.User)
	if !ok || user.ID == "" {
		c.Logger.Error().Msg("user not found in context for deleting cell output")
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// To verify ownership, we must check if the user owns the notebook associated with this output.
	// This requires getting the output, then its cell, then its notebook.
	_, err = c.Module.GetCellOutputByID(r.Context(), outputID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("output_id", outputIDStr).Msg("Output not found or not owned by user")
		pkg.WriteJSONResponseWithLogger(w, http.StatusNotFound, map[string]string{"error": "Output not found or not owned by user"}, &c.Logger)
		return
	}

	if err := c.Module.DeleteCellOutput(r.Context(), outputID, user.ID); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to delete cell output")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete cell output"}, &c.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}