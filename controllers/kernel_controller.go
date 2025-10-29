package controllers

import (
	"context"
	"encoding/json"
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

// KernelFunctionsHandler routes all /api/kernels and /api/kernels/{id} requests.
func (c *KernelController) KernelFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Determine if the request is for a specific kernel ID
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	
	// Example path: /api/kernels/12345678-abcd-1234-abcd-1234567890ab
	// If pathParts length is 4, the 4th element is the ID
	kernelID := ""
	if len(pathParts) >= 4 && pathParts[len(pathParts)-2] == "kernels" {
		kernelID = pathParts[len(pathParts)-1]
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// --- Route by Kernel ID presence and HTTP Method ---

	if kernelID != "" {
		// Case: /api/kernels/{kernel_id} or /api/kernels/{kernel_id}/action
		c.handleSingleKernel(w, r, ctx, kernelID)
	} else {
		// Case: /api/kernels (List all running kernels)
		c.handleKernelList(w, r, ctx)
	}
}

// handleSingleKernel handles GET, DELETE, and POST for specific kernel actions.
func (c *KernelController) handleSingleKernel(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	// Check for specific actions like /interrupt or /restart
	path := strings.Trim(r.URL.Path, "/")
	
	if strings.HasSuffix(path, "/interrupt") && r.Method == http.MethodPost {
		c.handleInterrupt(w, r, ctx, kernelID)
		return
	}
	if strings.HasSuffix(path, "/restart") && r.Method == http.MethodPost {
		c.handleRestart(w, r, ctx, kernelID)
		return
	}

	// Handle standard ID actions (GET / DELETE)
	switch r.Method {
	case http.MethodGet:
		c.handleGetKernelInfo(w, r, ctx, kernelID)
	case http.MethodDelete:
		c.handleDeleteKernel(w, r, ctx, kernelID)
	default:
		http.Error(w, "Method not allowed for single kernel resource", http.StatusMethodNotAllowed)
		c.Logger.Error().Msg("Method not allowed for single kernel resource")
	}
}

// handleKernelList handles the GET /api/kernels endpoint.
func (c *KernelController) handleKernelList(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	kernels, err := c.JupyterClient.GetKernels(ctx)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to retrieve running kernels list")
		http.Error(w, "Error retrieving kernel list", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, kernels)
}

// --- Specific Action Handlers ---

func (c *KernelController) handleGetKernelInfo(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	info, err := c.JupyterClient.GetKernelInfo(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to retrieve kernel info")
		http.Error(w, "Kernel not found or service error", http.StatusNotFound) // Use 404/500 based on expected error type
		return
	}

	writeJSONResponse(w, http.StatusOK, info)
}

func (c *KernelController) handleDeleteKernel(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	err := c.JupyterClient.DeleteKernel(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to delete kernel")
		http.Error(w, "Error deleting kernel", http.StatusInternalServerError)
		return
	}
	
	// 204 No Content response for successful deletion
	writeJSONResponse(w, http.StatusNoContent, nil)
}

func (c *KernelController) handleInterrupt(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	err := c.JupyterClient.InterruptKernel(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to interrupt kernel")
		http.Error(w, "Error interrupting kernel", http.StatusInternalServerError)
		return
	}
	// 204 No Content response for successful interrupt
	writeJSONResponse(w, http.StatusNoContent, nil)
}

func (c *KernelController) handleRestart(w http.ResponseWriter, r *http.Request, ctx context.Context, kernelID string) {
	info, err := c.JupyterClient.RestartKernel(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to restart kernel")
		http.Error(w, "Error restarting kernel", http.StatusInternalServerError)
		return
	}
	
	writeJSONResponse(w, http.StatusOK, info)
}