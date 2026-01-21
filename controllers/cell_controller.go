package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	Module *modules.CellModule
	Logger zerolog.Logger
}



// NewCellController creates and returns a new CellController.
func NewCellController(module *modules.CellModule, logger zerolog.Logger) *CellController {
	return &CellController{
		Module: module,
		Logger: logger,
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

	var req models.UpdateCellsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponseWithLogger(
			w, 
			http.StatusBadRequest, 
			map[string]string{"error": "Invalid request body"}, 
			&c.Logger,
		)
		return
	}

	if err := c.Module.UpdateCells(r.Context(), notebookID, &req); err != nil {
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
	var req models.CreateCellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, &c.Logger)
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

	cell, err := c.Module.CreateCell(r.Context(), &req)
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

	cells, err := c.Module.GetCellsByNotebookID(r.Context(), notebookID)
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

	cell, err := c.Module.GetCellByID(r.Context(), cellID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to get cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusNotFound, map[string]string{"error": "Cell not found"}, &c.Logger)
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

	cell, err := c.Module.UpdateCell(r.Context(), cellID, &req)
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

	if err := c.Module.DeleteCell(r.Context(), cellID); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to delete cell")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete cell"}, &c.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CellController) CreateCellOutputHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCellOutputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, &c.Logger)
		return
	}

	if _, ok := validOutputTypes[req.Type]; !ok {
		allowedTypes := make([]string, 0, len(validOutputTypes))
		for k := range validOutputTypes {
			allowedTypes = append(allowedTypes, k)
		}
		err_msg := fmt.Sprintf("Invalid output type: '%s'. Allowed types are: %s", req.Type, strings.Join(allowedTypes, ", "))
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": err_msg}, &c.Logger)
		return
	}

	output, err := c.Module.CreateCellOutput(r.Context(), &req)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to create cell output")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create cell output"}, &c.Logger)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, output, &c.Logger)
}

func (c *CellController) GetCellOutputsByCellIDHandler(w http.ResponseWriter, r *http.Request) {
	cellIDStr := r.PathValue("cell_id")
	cellID, err := uuid.Parse(cellIDStr)
	if err != nil {
		pkg.WriteJSONResponseWithLogger(w, http.StatusBadRequest, map[string]string{"error": "Invalid cell ID"}, &c.Logger)
		return
	}

	outputs, err := c.Module.GetCellOutputsByCellID(r.Context(), cellID)
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

	if err := c.Module.DeleteCellOutput(r.Context(), outputID); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to delete cell output")
		pkg.WriteJSONResponseWithLogger(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete cell output"}, &c.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
