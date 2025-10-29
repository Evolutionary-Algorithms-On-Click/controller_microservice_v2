package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/rs/zerolog"
)

// KernelController holds the dependencies required by the handlers.
type KernelController struct {
	JupyterClient *jupyterclient.Client
	Logger        zerolog.Logger
}

// NewKernelController creates and returns a new KernelController instance.
func NewKernelController(client *jupyterclient.Client, logger zerolog.Logger) *KernelController {
	return &KernelController{
		JupyterClient: client,
		Logger:        logger,
	}
}

// writeJSONResponse is a helper function to format and write a JSON response.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

// KernelFunctionsHandler acts as a single entry point for all /api/kernels routes.
// Handles:
//   - GET /api/kernels                     → List all kernels
//   - POST /api/kernels                    → Start new kernel
//   - GET /api/kernels/{id}                → Get kernel info
//   - DELETE /api/kernels/{id}             → Delete kernel
//   - POST /api/kernels/{id}/interrupt     → Interrupt kernel
//   - POST /api/kernels/{id}/restart       → Restart kernel
func (c *KernelController) KernelFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	// Expected structures:
	// [api, kernels]
	// [api, kernels, {id}]
	// [api, kernels, {id}, interrupt|restart]
	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "kernels" {
		http.Error(w, "Invalid route", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	switch len(pathParts) {
	case 2:
		// /api/kernels
		c.handleKernelList(w, r, ctx)
		return

	case 3:
		// /api/kernels/{id}
		kernelID := pathParts[2]
		c.handleSingleKernel(w, r, ctx, kernelID)
		return

	case 4:
		// /api/kernels/{id}/{action}
		kernelID := pathParts[2]
		action := pathParts[3]
		c.handleKernelAction(w, r, ctx, kernelID, action)
		return

	default:
		http.Error(w, "Invalid kernel path", http.StatusNotFound)
		return
	}
}

// handleKernelList handles both GET (list) and POST (start) requests to /api/kernels.
func (c *KernelController) handleKernelList(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	switch r.Method {

	case http.MethodGet:
		// --- List all running kernels ---
		kernels, err := c.JupyterClient.GetKernels(ctx)
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to retrieve running kernels list")
			http.Error(w, "Error retrieving kernel list", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, kernels)
		return

	case http.MethodPost:
		// --- Start a new kernel ---
		var reqBody struct {
			Language string `json:"language"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if reqBody.Language == "" {
			http.Error(w, "Missing 'language' field", http.StatusBadRequest)
			return
		}

		kernel, err := c.JupyterClient.StartKernel(ctx, reqBody.Language)
		if err != nil {
			c.Logger.Error().Err(err).Str("language", reqBody.Language).Msg("Failed to start kernel")
			http.Error(w, fmt.Sprintf("Error starting kernel: %v", err), http.StatusInternalServerError)
			return
		}

		writeJSONResponse(w, http.StatusCreated, kernel)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// handleSingleKernel handles GET and DELETE for /api/kernels/{id}.
func (c *KernelController) handleSingleKernel(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	switch r.Method {
	case http.MethodGet:
		info, err := c.JupyterClient.GetKernelInfo(ctx, kernelID)
		if err != nil {
			c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to retrieve kernel info")
			http.Error(w, "Kernel not found or service error", http.StatusNotFound)
			return
		}
		writeJSONResponse(w, http.StatusOK, info)

	case http.MethodDelete:
		err := c.JupyterClient.DeleteKernel(ctx, kernelID)
		if err != nil {
			c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to delete kernel")
			http.Error(w, "Error deleting kernel", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusNoContent, nil)

	default:
		http.Error(w, "Method not allowed for single kernel resource", http.StatusMethodNotAllowed)
	}
}

// handleKernelAction routes POST /api/kernels/{id}/interrupt|restart requests.
func (c *KernelController) handleKernelAction(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID, action string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed for kernel action", http.StatusMethodNotAllowed)
		return
	}

	switch action {
	case "interrupt":
		if err := c.JupyterClient.InterruptKernel(ctx, kernelID); err != nil {
			c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to interrupt kernel")
			http.Error(w, "Error interrupting kernel", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusNoContent, nil)

	case "restart":
		info, err := c.JupyterClient.RestartKernel(ctx, kernelID)
		if err != nil {
			c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to restart kernel")
			http.Error(w, "Error restarting kernel", http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, info)

	default:
		http.Error(w, "Unknown kernel action", http.StatusNotFound)
	}
}
