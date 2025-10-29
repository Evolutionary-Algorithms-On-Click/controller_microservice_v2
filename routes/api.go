package routes

import (
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/controllers"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
)

func RegisterAPIRoutes(mux *http.ServeMux, c *jupyterclient.Client) {
	problemModule := modules.NewProblemModule()
	problemController := controllers.NewProblemController(problemModule, *pkg.Logger)
	kernelController := controllers.NewKernelController(c, *pkg.Logger)

	// Register the handler functions
	mux.HandleFunc("/api/problems", problemController.CreateAndListProblemsHandler)
	mux.HandleFunc("/api/problems/", problemController.ProblemByIDHandler)
	// You might also need a handler for listing all problems
	// mux.HandleFunc("/api/problems/user", problemController.GetProblemsByUserIDHandler)

	// handlers for other routes related to jupyter kernel gateway
	mux.HandleFunc("/api/kernels/", kernelController.KernelFunctionsHandler)
}
