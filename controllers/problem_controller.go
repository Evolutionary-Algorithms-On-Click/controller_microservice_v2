package controllers

import (
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
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
	// TODO: Implement logic to create a new problem.
	// 1. Extract user_id from request context (from auth middleware).
	// 2. Deserialize the JSON request body into models.CreateProblemRequest.
	// 3. Call c.ProblemModule.CreateProblem.
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// ListProblemsHandler handles GET /api/v1/problems
func (c *ProblemController) ListProblemsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to list all problems for a user.
	// 1. Extract the user_id from the request context.
	// 2. Call c.ProblemModule.GetProblemsByUserID.
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// GetProblemByIDHandler handles GET /api/v1/problems/{id}
func (c *ProblemController) GetProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	// TODO: Implement logic to get a problem by its ID.
	// 1. Call c.ProblemModule.GetProblemByID with problemID.
	http.Error(w, "Not Implemented: "+problemID, http.StatusNotImplemented)
}

// UpdateProblemByIDHandler handles PUT /api/v1/problems/{id}
func (c *ProblemController) UpdateProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	// TODO: Implement logic to update a problem.
	// 1. Deserialize the JSON request body.
	// 2. Call c.ProblemModule.UpdateProblem with problemID and request body.
	http.Error(w, "Not Implemented: "+problemID, http.StatusNotImplemented)
}

// DeleteProblemByIDHandler handles DELETE /api/v1/problems/{id}
func (c *ProblemController) DeleteProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	problemID := r.PathValue("id")
	// TODO: Implement logic to delete a problem.
	// 1. Call c.ProblemModule.DeleteProblem with problemID.
	http.Error(w, "Not Implemented: "+problemID, http.StatusNotImplemented)
}
