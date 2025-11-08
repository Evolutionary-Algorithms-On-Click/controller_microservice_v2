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
	notebookModule := modules.NewNotebookModule()

	problemController := controllers.NewProblemController(problemModule, *pkg.Logger)
	kernelController := controllers.NewKernelController(c, *pkg.Logger)
	notebookController := controllers.NewNotebookController(notebookModule, pkg.Logger)

	// Register the handler functions with API versioning (v1)

	// Problem Routes
	mux.HandleFunc("POST /api/v1/problems", problemController.CreateProblemHandler)
	mux.HandleFunc("GET /api/v1/problems", problemController.ListProblemsHandler)
	mux.HandleFunc("GET /api/v1/problems/{id}", problemController.GetProblemByIDHandler)
	mux.HandleFunc("PUT /api/v1/problems/{id}", problemController.UpdateProblemByIDHandler)
	mux.HandleFunc("DELETE /api/v1/problems/{id}", problemController.DeleteProblemByIDHandler)

	// Notebook Routes
	mux.HandleFunc("POST /api/v1/notebooks", notebookController.CreateNotebookHandler)
	mux.HandleFunc("GET /api/v1/notebooks", notebookController.ListNotebooksHandler)
	mux.HandleFunc("GET /api/v1/notebooks/{id}", notebookController.GetNotebookByIDHandler)
	mux.HandleFunc("PUT /api/v1/notebooks/{id}", notebookController.UpdateNotebookByIDHandler)
	mux.HandleFunc("DELETE /api/v1/notebooks/{id}", notebookController.DeleteNotebookByIDHandler)

	// Kernel Routes
	mux.HandleFunc("POST /api/v1/kernels", kernelController.StartKernelHandler)
	mux.HandleFunc("GET /api/v1/kernels", kernelController.ListKernelsHandler)
	mux.HandleFunc("GET /api/v1/kernels/{id}", kernelController.GetKernelInfoHandler)
	mux.HandleFunc("DELETE /api/v1/kernels/{id}", kernelController.DeleteKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/interrupt", kernelController.InterruptKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/restart", kernelController.RestartKernelHandler)
}