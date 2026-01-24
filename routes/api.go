package routes

import (
	"net/http"
	"os"

	"github.com/Thanus-Kumaar/controller_microservice_v2/controllers"
	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	// "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/middleware" // New import
)

func RegisterAPIRoutes(mux *http.ServeMux, c *jupyterclient.Client) {

	// Initialize Repositories
	notebookRepo := repository.NewNotebookRepository(db.Pool)
	llmRepo := repository.NewLlmProxy(os.Getenv("LLM_MICROSERVICE_URL"))
	sessionRepo := repository.NewSessionRepository(db.Pool)
	problemRepo := repository.NewProblemRepository(db.Pool)
	cellRepo := repository.NewCellRepository(db.Pool, *pkg.Logger)

	// Initialize Modules
	notebookModule := modules.NewNotebookModule(notebookRepo)
	llmModule := modules.NewLlmModule(llmRepo)
	sessionModule := modules.NewSessionModule(sessionRepo, c, *pkg.Logger)
	problemModule := modules.NewProblemModule(problemRepo)
	cellModule := modules.NewCellModule(cellRepo, *pkg.Logger)

	// Initialize Controllers
	notebookController := controllers.NewNotebookController(notebookModule, pkg.Logger)
	sessionController := controllers.NewSessionController(sessionModule, *pkg.Logger)
	llmController := controllers.NewLlmController(llmModule, *pkg.Logger)
	problemController := controllers.NewProblemController(problemModule, *pkg.Logger)
	cellController := controllers.NewCellController(cellModule, *pkg.Logger)
	kernelController := controllers.NewKernelController(c, *pkg.Logger, cellRepo)

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
	mux.HandleFunc("PATCH /api/v1/notebooks/{notebook_id}/cells", cellController.UpdateCellsHandler)

	// Session Routes
	// mux.Handle("POST /api/v1/sessions", authMiddleware.Authenticate(http.HandlerFunc(sessionController.CreateSessionHandler))) // Applied middleware
	mux.HandleFunc("POST /api/v1/sessions", sessionController.CreateSessionHandler)
	mux.HandleFunc("GET /api/v1/sessions", sessionController.ListSessionsHandler)
	mux.HandleFunc("GET /api/v1/sessions/{id}", sessionController.GetSessionByIDHandler)
	mux.HandleFunc("PUT /api/v1/sessions/{id}", sessionController.UpdateSessionByIDHandler)
	mux.HandleFunc("DELETE /api/v1/sessions/{id}", sessionController.DeleteSessionByIDHandler)

	// Cell Routes
	mux.HandleFunc("POST /api/v1/notebooks/{notebook_id}/cells", cellController.CreateCellHandler)
	mux.HandleFunc("GET /api/v1/notebooks/{notebook_id}/cells", cellController.GetCellsByNotebookIDHandler)
	mux.HandleFunc("GET /api/v1/cells/{cell_id}", cellController.GetCellByIDHandler)
	mux.HandleFunc("PUT /api/v1/cells/{cell_id}", cellController.UpdateCellHandler)
	mux.HandleFunc("DELETE /api/v1/cells/{cell_id}", cellController.DeleteCellHandler)

	// Cell Output Routes
	mux.HandleFunc("POST /api/v1/cells/{cell_id}/outputs", cellController.CreateCellOutputHandler)
	mux.HandleFunc("GET /api/v1/cells/{cell_id}/outputs", cellController.GetCellOutputsByCellIDHandler)
	mux.HandleFunc("DELETE /api/v1/outputs/{output_id}", cellController.DeleteCellOutputHandler)

	// Llm Routes
	mux.HandleFunc("POST /api/v1/llm/generate", llmController.GenerateNotebookHandler)
	mux.HandleFunc("POST /api/v1/llm/modify", llmController.ModifyNotebookHandler)
	mux.HandleFunc("POST /api/v1/llm/fix", llmController.FixNotebookHandler)

	// Kernel Routes
	mux.HandleFunc("POST /api/v1/kernels", kernelController.StartKernelHandler)
	mux.HandleFunc("GET /api/v1/kernels", kernelController.ListKernelsHandler)
	mux.HandleFunc("GET /api/v1/kernels/{id}", kernelController.GetKernelInfoHandler)
	mux.HandleFunc("DELETE /api/v1/kernels/{id}", kernelController.DeleteKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/interrupt", kernelController.InterruptKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/restart", kernelController.RestartKernelHandler)
	mux.HandleFunc("GET /api/v1/kernels/{id}/channels", kernelController.KernelChannelsHandler)
}
