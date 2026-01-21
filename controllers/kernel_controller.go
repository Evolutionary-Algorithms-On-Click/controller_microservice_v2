package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for now. In production, you'd want to restrict this.
		return true
	},
}

// KernelController holds the dependencies required by the handlers.
type KernelController struct {
	JupyterClient  *jupyterclient.Client
	Logger         zerolog.Logger
	CellRepo       repository.CellRepository
	msgIDCellIDMap map[string]uuid.UUID
	mapMutex       sync.RWMutex
}

// NewKernelController creates and returns a new KernelController instance.
func NewKernelController(client *jupyterclient.Client, logger zerolog.Logger, cellRepo repository.CellRepository) *KernelController {
	return &KernelController{
		JupyterClient:  client,
		Logger:         logger,
		CellRepo:       cellRepo,
		msgIDCellIDMap: make(map[string]uuid.UUID),
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

// ListKernelsHandler handles GET /api/v1/kernels to list all running kernels.
func (c *KernelController) ListKernelsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	kernels, err := c.JupyterClient.GetKernels(ctx)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to retrieve running kernels list")
		http.Error(w, "Error retrieving kernel list", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, kernels)
}

// StartKernelHandler handles POST /api/v1/kernels to start a new kernel.
func (c *KernelController) StartKernelHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

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
}

// GetKernelInfoHandler handles GET /api/v1/kernels/{id} to get info on a single kernel.
func (c *KernelController) GetKernelInfoHandler(w http.ResponseWriter, r *http.Request) {
	kernelID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	info, err := c.JupyterClient.GetKernelInfo(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to retrieve kernel info")
		http.Error(w, "Kernel not found or service error", http.StatusNotFound)
		return
	}
	writeJSONResponse(w, http.StatusOK, info)
}

// DeleteKernelHandler handles DELETE /api/v1/kernels/{id} to delete a kernel.
func (c *KernelController) DeleteKernelHandler(w http.ResponseWriter, r *http.Request) {
	kernelID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := c.JupyterClient.DeleteKernel(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to delete kernel")
		http.Error(w, "Error deleting kernel", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusNoContent, nil)
}

// InterruptKernelHandler handles POST /api/v1/kernels/{id}/interrupt.
func (c *KernelController) InterruptKernelHandler(w http.ResponseWriter, r *http.Request) {
	kernelID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := c.JupyterClient.InterruptKernel(ctx, kernelID); err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to interrupt kernel")
		http.Error(w, "Error interrupting kernel", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusNoContent, nil)
}

// RestartKernelHandler handles POST /api/v1/kernels/{id}/restart.
func (c *KernelController) RestartKernelHandler(w http.ResponseWriter, r *http.Request) {
	kernelID := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	info, err := c.JupyterClient.RestartKernel(ctx, kernelID)
	if err != nil {
		c.Logger.Error().Err(err).Str("kernel_id", kernelID).Msg("Failed to restart kernel")
		http.Error(w, "Error restarting kernel", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, info)
}

// KernelChannelsHandler handles GET /api/v1/kernels/{id}/channels
func (c *KernelController) KernelChannelsHandler(w http.ResponseWriter, r *http.Request) {
	kernelID := r.PathValue("id")

	// Upgrade the client's connection
	feConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to upgrade connection")
		return
	}
	defer feConn.Close()

	// Get the Jupyter Gateway URL and token from the client
	gatewayURL := c.JupyterClient.GetGatewayURL()
	gatewayToken := c.JupyterClient.GetAuthToken()

	// Construct the target websocket URL
	wsURL := "ws" + strings.TrimPrefix(gatewayURL, "http")
	targetURL := fmt.Sprintf("%s/api/kernels/%s/channels", wsURL, kernelID)

	// Create the request headers, including the auth token
	headers := http.Header{}
	headers.Set("Authorization", "token "+gatewayToken)

	c.Logger.Debug().Str("targetURL", targetURL).Interface("headers", headers).Msg("Attempting to dial kernel gateway")
	// Connect to the Jupyter Kernel Gateway
	kgConn, _, err := websocket.DefaultDialer.Dial(targetURL, headers)
	if err != nil {
		c.Logger.Error().Err(err).Str("target", targetURL).Msg("failed to dial kernel gateway")
		if writeErr := feConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "could not connect to kernel")); writeErr != nil {
			c.Logger.Error().Err(writeErr).Msg("failed to write close message to frontend")
		}
		return
	}
	defer kgConn.Close()

	c.Logger.Info().Str("kernel_id", kernelID).Msg("websocket proxy established")

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to proxy messages from Frontend to Kernel Gateway
	go func() {
		defer wg.Done()
		for {
			messageType, p, err := feConn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					c.Logger.Warn().Err(err).Msg("error reading from frontend, closing proxy")
				}
				kgConn.Close()
				return
			}

			// If it's an execute_request, store the msg_id -> cell_id mapping
			var msg jupyterclient.Message
			if err := json.Unmarshal(p, &msg); err == nil && msg.Header.MsgType == "execute_request" {
				var metadata map[string]any
				if err := json.Unmarshal(msg.Metadata, &metadata); err == nil {
					if cellIDStr, ok := metadata["cell_id"].(string); ok {
						if cellID, err := uuid.Parse(cellIDStr); err == nil {
							c.mapMutex.Lock()
							c.msgIDCellIDMap[msg.Header.MsgID] = cellID
							c.mapMutex.Unlock()
						}
					}
				}
			}

			if err := kgConn.WriteMessage(messageType, p); err != nil {
				c.Logger.Warn().Err(err).Msg("error writing to kernel gateway, closing proxy")
				return
			}
			c.Logger.Trace().Str("direction", "FE->KG").Int("size", len(p)).Msg("proxied message")
		}
	}()

	// Goroutine to proxy messages from Kernel Gateway to Frontend
	go func() {
		defer wg.Done()
		for {
			messageType, p, err := kgConn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					c.Logger.Warn().Err(err).Msg("error reading from kernel gateway, closing proxy")
				}
				// Closing the frontend connection will cause the other goroutine's ReadMessage to error out.
				feConn.Close()
				return
			}

			// Asynchronously save the output
			go c.saveCellOutput(p)

			if err := feConn.WriteMessage(messageType, p); err != nil {
				c.Logger.Warn().Err(err).Msg("error writing to frontend, closing proxy")
				return
			}
			c.Logger.Trace().Str("direction", "KG->FE").Int("size", len(p)).Msg("proxied message")
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	c.Logger.Info().Str("kernel_id", kernelID).Msg("websocket proxy closed")
}

func (c *KernelController) saveCellOutput(p []byte) {
	var msg jupyterclient.Message
	if err := json.Unmarshal(p, &msg); err != nil {
		c.Logger.Debug().Err(err).Msg("failed to unmarshal jupyter message")
		return
	}

	outputTypes := map[string]struct{}{
		"stream":         {},
		"display_data":   {},
		"execute_result": {},
		"error":          {},
		"execute_reply":  {},
	}
	if _, ok := outputTypes[msg.Header.MsgType]; !ok {
		return
	}

	c.mapMutex.RLock()
	cellID, ok := c.msgIDCellIDMap[msg.ParentHeader.MsgID]
	c.mapMutex.RUnlock()

	if !ok {
		c.Logger.Warn().Str("parent_msg_id", msg.ParentHeader.MsgID).Msg("cell_id not found for parent message")
		return
	}

	// Clean up the map when the execution is done
	if msg.Header.MsgType == "execute_reply" {
		c.mapMutex.Lock()
		delete(c.msgIDCellIDMap, msg.ParentHeader.MsgID)
		c.mapMutex.Unlock()
		return // No need to save execute_reply as an output
	}

	var outputData models.CellOutput
	outputData.ID = uuid.New()
	outputData.CellID = cellID
	outputData.Type = msg.Header.MsgType

	switch msg.Header.MsgType {
	case "stream":
		var content jupyterclient.StreamContent
		if err := json.Unmarshal(msg.Content, &content); err == nil {
			outputData.DataJSON, _ = json.Marshal(content)
		}
	case "display_data":
		var content jupyterclient.DisplayDataContent
		if err := json.Unmarshal(msg.Content, &content); err == nil {
			outputData.DataJSON, _ = json.Marshal(content.Data)
		}
	case "execute_result":
		var content jupyterclient.ExecuteResultContent
		if err := json.Unmarshal(msg.Content, &content); err == nil {
			outputData.ExecutionCount = content.ExecutionCount
			outputData.DataJSON, _ = json.Marshal(content.Data)
		}
	case "error":
		var content jupyterclient.ErrorContent
		if err := json.Unmarshal(msg.Content, &content); err == nil {
			outputData.DataJSON, _ = json.Marshal(content)
		}
	}

	outputs, err := c.CellRepo.GetCellOutputsByCellID(context.Background(), cellID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get cell outputs to determine next index")
	}
	outputData.OutputIndex = len(outputs)

	if _, err := c.CellRepo.CreateCellOutput(context.Background(), &outputData); err != nil {
		c.Logger.Error().Err(err).Msg("failed to save cell output")
	} else {
		c.Logger.Info().Str("cell_id", cellID.String()).Str("type", outputData.Type).Msg("successfully saved cell output")
	}
}


