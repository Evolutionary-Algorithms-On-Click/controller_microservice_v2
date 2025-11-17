package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/rs/zerolog"
)

// ProblemController holds the dependencies for the problem statement handlers.
type ProblemController struct {
	ProblemModule *modules.ProblemModule
	Logger        zerolog.Logger
}

// NewProblemController creates and returns a new ProblemController.
func NewProblemController(problemModule *modules.ProblemModule, logger zerolog.Logger) *ProblemController {
	return &ProblemController{
		ProblemModule: problemModule,
		Logger:        logger,
	}
}

// CreateProblemHandler handles POST /api/v1/problems
func (c *ProblemController) CreateProblemHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// For now, we'll use a hardcoded user ID for testing.
	// TODO: Replace with actual user_id from auth middleware.
	const hardcodedUserID = "123e4567-e89b-12d3-a456-426614174000"

	var req models.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	problemStatement, err := c.ProblemModule.CreateProblem(ctx, &req, hardcodedUserID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to create problem statement")
		http.Error(w, fmt.Sprintf("error creating problem statemnet: %v", err), http.StatusInternalServerError)
		return
	}
	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, problemStatement, &c.Logger)
}

// ListProblemsHandler handles GET /api/v1/problems
func (c *ProblemController) ListProblemsHandler(w http.ResponseWriter, r *http.Request) {
	// For now, we'll use a hardcoded user ID for testing.
	// TODO: Replace with actual user_id from auth middleware.
	const hardcodedUserID = "123e4567-e89b-12d3-a456-426614174000"
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	problems, err := c.ProblemModule.GetProblemsByUserID(ctx, hardcodedUserID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to list problems by user id")
		http.Error(w, "failed to retrieve problems", http.StatusInternalServerError)
		return
	}

	// If no problems are found, it should return an empty list, not an error.
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, problems, &c.Logger)
}

// GetProblemByIDHandler handles GET /api/v1/problems/{id}
func (c *ProblemController) GetProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	problem, err := c.ProblemModule.GetProblemByID(ctx, problemID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get problem by id")
		if err == sql.ErrNoRows {
			http.Error(w, "problem not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to retrive problem", http.StatusInternalServerError)
		}
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, problem, &c.Logger)
}

// UpdateProblemByIDHandler handles PUT /api/v1/problems/{id}
func (c *ProblemController) UpdateProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// For now, we'll use a hardcoded user ID for testing.
	// TODO: Replace with actual user_id from auth middleware.
	const hardcodedUserID = "123e4567-e89b-12d3-a456-426614174000"

	var req models.UpdateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedProblem, err := c.ProblemModule.UpdateProblem(ctx, problemID, &req, hardcodedUserID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to update problem")
		// This could be a not found error or an authorization error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, updatedProblem, &c.Logger)
}

// DeleteProblemByIDHandler handles DELETE /api/v1/problems/{id}
func (c *ProblemController) DeleteProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// For now, we'll use a hardcoded user ID for testing.
	// TODO: Replace with actual user_id from auth middleware.
	const hardcodedUserID = "123e4567-e89b-12d3-a456-426614174000"

	err := c.ProblemModule.DeleteProblem(ctx, problemID, hardcodedUserID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to delete problem")
		// This could be a not found error or an authorization error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
