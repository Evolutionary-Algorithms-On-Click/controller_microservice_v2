package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/middleware"
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
	c.Logger.Info().Msg("received request to create a new problem")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		c.Logger.Error().Msg("user not found in context for problem creation")
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	c.Logger.Info().Str("userID", user.ID).Msg("user authorized for problem creation")

	var req models.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Logger.Error().Err(err).Msg("invalid request body for problem creation")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	c.Logger.Debug().
		Str("title", req.Title).
		RawJSON("description_json", req.DescriptionJSON).
		Msg("decoded create problem request body")

	problemStatement, err := c.ProblemModule.CreateProblem(ctx, &req, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).
			Str("userID", user.ID).
			Str("title", req.Title).
			Msg("failed to create problem statement via module")
		http.Error(w, fmt.Sprintf("error creating problem statemnet: %v", err), http.StatusInternalServerError)
		return
	}

	c.Logger.Info().Str("problemID", problemStatement.ID.String()).Msg("problem statement created successfully")
	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, problemStatement, &c.Logger)
}

// ListProblemsHandler handles GET /api/v1/problems
func (c *ProblemController) ListProblemsHandler(w http.ResponseWriter, r *http.Request) {
	c.Logger.Info().Msg("received request to list problems")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		c.Logger.Error().Msg("user not found in context for listing problems")
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	c.Logger.Info().Str("userID", user.ID).Msg("user authorized for listing problems")

	problems, err := c.ProblemModule.GetProblemsByUserID(ctx, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Str("userID", user.ID).Msg("failed to list problems by user id via module")
		http.Error(w, "failed to retrieve problems", http.StatusInternalServerError)
		return
	}

	c.Logger.Info().
		Str("userID", user.ID).
		Int("problem_count", len(problems)).
		Msg("successfully listed problems for user")
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, problems, &c.Logger)
}

// GetProblemByIDHandler handles GET /api/v1/problems/{id}
func (c *ProblemController) GetProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	c.Logger.Info().Str("problemID", problemID).Msg("received request to get problem by ID")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	problem, err := c.ProblemModule.GetProblemByID(ctx, problemID)
	if err != nil {
		c.Logger.Error().Err(err).Str("problemID", problemID).Msg("failed to get problem by id via module")
		if err.Error() == "problem not found" { // Matching the error string from problem_module
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to retrieve problem", http.StatusInternalServerError)
		return
	}

	c.Logger.Info().Str("problemID", problemID).Msg("successfully retrieved problem by ID")
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, problem, &c.Logger)
}

// UpdateProblemByIDHandler handles PUT /api/v1/problems/{id}
func (c *ProblemController) UpdateProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	c.Logger.Info().Str("problemID", problemID).Msg("received request to update problem by ID")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		c.Logger.Error().Msg("user not found in context for updating problem")
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	c.Logger.Info().Str("userID", user.ID).Msg("user authorized for updating problem")

	var req models.UpdateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Logger.Error().Err(err).Msg("invalid request body for updating problem")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	c.Logger.Debug().
		Str("problemID", problemID).
		Str("userID", user.ID).
		Str("title_update", req.Title).
		RawJSON("description_json_update", req.DescriptionJSON).
		Msg("decoded update problem request body")

	updatedProblem, err := c.ProblemModule.UpdateProblem(ctx, problemID, &req, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).
			Str("problemID", problemID).
			Str("userID", user.ID).
			Msg("failed to update problem via module")
		// This could be a not found error or an authorization error.
		if err.Error() == "problem not found" {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}
		if err.Error() == "user not authorized to update this problem" {
			http.Error(w, "user not authorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Logger.Info().Str("problemID", updatedProblem.ID.String()).Msg("successfully updated problem by ID")
	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, updatedProblem, &c.Logger)
}

// DeleteProblemByIDHandler handles DELETE /api/v1/problems/{id}
func (c *ProblemController) DeleteProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	c.Logger.Info().Str("problemID", problemID).Msg("received request to delete problem by ID")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		c.Logger.Error().Msg("user not found in context for deleting problem")
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	c.Logger.Info().Str("userID", user.ID).Msg("user authorized for deleting problem")

	err := c.ProblemModule.DeleteProblem(ctx, problemID, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).
			Str("problemID", problemID).
			Str("userID", user.ID).
			Msg("failed to delete problem via module")
		// This could be a not found error or an authorization error.
		if err.Error() == "problem not found" {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}
		if err.Error() == "user not authorized to delete this problem" {
			http.Error(w, "user not authorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Logger.Info().Str("problemID", problemID).Msg("successfully deleted problem by ID")
	w.WriteHeader(http.StatusNoContent)
}
