package routes

import (
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/controllers"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

func RegisterAPIRoutes(mux *http.ServeMux) {
	problemModule := modules.NewProblemModule()
problemController := controllers.NewProblemController(problemModule, *pkg.Logger)

	// Register the handler functions
	mux.HandleFunc("/api/problems", problemController.CreateAndListProblemsHandler)
	mux.HandleFunc("/api/problems/", problemController.ProblemByIDHandler)
	// You might also need a handler for listing all problems
	// mux.HandleFunc("/api/problems/user", problemController.GetProblemsByUserIDHandler)
}
