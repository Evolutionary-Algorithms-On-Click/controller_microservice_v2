package routes

import (
	"net/http"
	"os"

	"github.com/Thanus-Kumaar/controller_microservice_v2/controllers"
	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/middleware"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
)

func RegisterAPIRoutes(mux *http.ServeMux, c *jupyterclient.Client) {

	// Initialize Repositories
	notebookRepo := repository.NewNotebookRepository(db.Pool)
	llmRepo := repository.NewLlmProxy(os.Getenv("LLM_MICROSERVICE_URL"))
	sessionRepo := repository.NewSessionRepository(db.Pool)
	problemRepo := repository.NewProblemRepository(db.Pool).WithLogger(*pkg.Logger)
	cellRepo := repository.NewCellRepository(db.Pool, *pkg.Logger)

	userDataDir := os.Getenv("USER_DATA_DIR")
	if userDataDir == "" {
		userDataDir = "/mnt/user_data"
	}
	fileModule := modules.NewFileModule(userDataDir)

	// Initialize Modules
	notebookModule := modules.NewNotebookModule(notebookRepo, problemRepo)
	llmModule := modules.NewLlmModule(llmRepo)
	sessionModule := modules.NewSessionModule(sessionRepo, c, *pkg.Logger, notebookRepo)
	problemModule := modules.NewProblemModule(problemRepo, notebookRepo, fileModule, *pkg.Logger) // Pass the logger here
	cellModule := modules.NewCellModule(cellRepo, *pkg.Logger)

	// Initialize Controllers
	notebookController := controllers.NewNotebookController(notebookModule, pkg.Logger)
	sessionController := controllers.NewSessionController(sessionModule, *pkg.Logger)
	llmController := controllers.NewLlmController(llmModule, *pkg.Logger)
	problemController := controllers.NewProblemController(problemModule, *pkg.Logger)
	cellController := controllers.NewCellController(cellModule, *pkg.Logger, notebookModule)
	kernelController := controllers.NewKernelController(c, *pkg.Logger, cellRepo)
	fileController := controllers.NewFileController(fileModule, *pkg.Logger)

	// Register the handler functions with API versioning (v1)

	// Problem Routes
	mux.Handle("POST /api/v1/problems",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.CreateProblemHandler)))
	mux.Handle("GET /api/v1/problems",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.ListProblemsHandler)))
	mux.Handle("GET /api/v1/problems/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.GetProblemByIDHandler)))
	mux.Handle("PUT /api/v1/problems/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.UpdateProblemByIDHandler)))
	mux.Handle("DELETE /api/v1/problems/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.DeleteProblemByIDHandler)))

	// VolPE Routes
	mux.Handle("POST /api/v1/submission/submit",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.SubmitNotebookHandler)))
	mux.Handle("GET /api/v1/submission/results/{problemId}",
		middleware.AuthMiddleware(http.HandlerFunc(problemController.GetSubmissionResultsHandler)))

	// Notebook Routes
	mux.Handle("POST /api/v1/notebooks",
		middleware.AuthMiddleware(http.HandlerFunc(notebookController.CreateNotebookHandler)))
	mux.Handle("GET /api/v1/notebooks",
		middleware.AuthMiddleware(http.HandlerFunc(notebookController.ListNotebooksHandler)))
	mux.Handle("GET /api/v1/notebooks/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(notebookController.GetNotebookByIDHandler)))
	mux.Handle("PUT /api/v1/notebooks/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(notebookController.UpdateNotebookByIDHandler)))
	mux.Handle("DELETE /api/v1/notebooks/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(notebookController.DeleteNotebookByIDHandler)))
	mux.Handle("PATCH /api/v1/notebooks/{notebook_id}/cells",
		middleware.AuthMiddleware(http.HandlerFunc(cellController.UpdateCellsHandler)))

	// Session Routes
	mux.Handle("POST /api/v1/sessions",
		middleware.AuthMiddleware(http.HandlerFunc(sessionController.CreateSessionHandler)))
	mux.Handle("GET /api/v1/sessions",
		middleware.AuthMiddleware(http.HandlerFunc(sessionController.ListSessionsHandler)))
	mux.Handle("GET /api/v1/sessions/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(sessionController.GetSessionByIDHandler)))
	mux.Handle("PUT /api/v1/sessions/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(sessionController.UpdateSessionByIDHandler)))
	mux.Handle("DELETE /api/v1/sessions/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(sessionController.DeleteSessionByIDHandler)))

	// User file Routes
	mux.Handle("POST /api/v1/sessions/{session_id}/files",
		middleware.AuthMiddleware(http.HandlerFunc(fileController.UploadFileHandler)))
	mux.Handle("GET /api/v1/sessions/{session_id}/files",
		middleware.AuthMiddleware(http.HandlerFunc(fileController.ListFilesHandler)))
	mux.Handle("DELETE /api/v1/sessions/{session_id}/files/{filename}",
		middleware.AuthMiddleware(http.HandlerFunc(fileController.DeleteFileHandler)))

	// Cell Routes
	mux.Handle("POST /api/v1/notebooks/{notebook_id}/cells", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.CreateCellHandler)))
	mux.Handle("GET /api/v1/notebooks/{notebook_id}/cells", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.GetCellsByNotebookIDHandler)))
	mux.Handle("GET /api/v1/cells/{cell_id}", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.GetCellByIDHandler)))
	mux.Handle("PUT /api/v1/cells/{cell_id}", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.UpdateCellHandler)))
	mux.Handle("DELETE /api/v1/cells/{cell_id}", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.DeleteCellHandler)))

	// Cell Output Routes
	mux.Handle("POST /api/v1/cells/{cell_id}/outputs", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.CreateCellOutputHandler)))
	mux.Handle("GET /api/v1/cells/{cell_id}/outputs", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.GetCellOutputsByCellIDHandler)))
	mux.Handle("DELETE /api/v1/outputs/{output_id}", 
		middleware.AuthMiddleware(http.HandlerFunc(cellController.DeleteCellOutputHandler)))

	// Llm Routes
	mux.Handle("POST /api/v1/llm/generate",
		middleware.AuthMiddleware(http.HandlerFunc(llmController.GenerateNotebookHandler)))
	mux.Handle("POST /api/v1/llm/modify",
		middleware.AuthMiddleware(http.HandlerFunc(llmController.ModifyNotebookHandler)))
	mux.Handle("POST /api/v1/llm/fix",
		middleware.AuthMiddleware(http.HandlerFunc(llmController.FixNotebookHandler)))

	// Kernel Routes
	mux.HandleFunc("POST /api/v1/kernels", kernelController.StartKernelHandler)
	mux.HandleFunc("GET /api/v1/kernels", kernelController.ListKernelsHandler)
	mux.HandleFunc("GET /api/v1/kernels/{id}", kernelController.GetKernelInfoHandler)
	mux.HandleFunc("DELETE /api/v1/kernels/{id}", kernelController.DeleteKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/interrupt", kernelController.InterruptKernelHandler)
	mux.HandleFunc("POST /api/v1/kernels/{id}/restart", kernelController.RestartKernelHandler)
	mux.HandleFunc("GET /api/v1/kernels/{id}/channels", kernelController.KernelChannelsHandler)
}
