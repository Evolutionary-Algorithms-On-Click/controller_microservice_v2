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

// CreateAndListProblemsHandler handles POST to create and GET to list all problems.
func (c *ProblemController) CreateAndListProblemsHandler(w http.ResponseWriter, r *http.Request) {
	// A handler MUST always check the method.
	switch r.Method {
	case http.MethodPost:
		// TODO: Implement logic to create a new problem.
		// You would use the same logic as the old CreateProblemHandler.
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	case http.MethodGet:
		// TODO: Implement logic to list all problems for a user.
		// 1. Extract the user_id from the request context (e.g., from a JWT token).
		// 2. Call the ProblemModule's GetProblemsByUserID function.
		// 3. Handle errors and return a JSON response.
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ProblemByIDHandler handles GET, PUT, and DELETE for a specific problem ID.
func (c *ProblemController) ProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract problem_id from the URL path.

	switch r.Method {
	case http.MethodGet:
		// TODO: Implement logic to get a problem by its ID.
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	case http.MethodPut:
		// TODO: Implement logic to update a problem.
		// 1. Deserialize the JSON request body.
		// 2. Call the ProblemModule's UpdateProblem function.
		// 3. Handle errors and return a JSON response.
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	case http.MethodDelete:
		// TODO: Implement logic to delete a problem.
		// 1. Call the ProblemModule's DeleteProblem function.
		// 2. Handle errors and return a 204 No Content response.
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}