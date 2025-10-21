package handlers

import (
	// "net/http"

	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/rs/zerolog"
)

// Handlers struct holds the dependencies required by API handlers.
type Handlers struct {
	JupyterClient *jupyterclient.Client
	Logger        zerolog.Logger
}

// NewHandlers creates and returns a new Handlers instance.
func NewHandlers(client *jupyterclient.Client, logger zerolog.Logger) *Handlers {
	return &Handlers{
		JupyterClient: client,
		Logger:        logger,
	}
}

// // StartSessionHandler is a handler function for the POST /api/sessions/start endpoint.
// func (h *Handlers) StartSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Implement the logic to start a new session.
// 	// 1. Extract user data from the request (e.g., from a JWT token).
// 	// 2. Call h.JupyterClient.StartKernel() to create a new kernel.
// 	// 3. Store the session and kernel info in the database.
// 	// 4. Respond with the new session ID and WebSocket URL.

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNotImplemented)
// 	w.Write([]byte(`{"error": "not implemented yet"}`))
// }

// // RunCodeHandler is a handler function for the WebSocket connection.
// func (h *Handlers) RunCodeHandler(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Implement the WebSocket proxying logic.
// 	// This will upgrade the connection and proxy messages to the kernel.
// 	// 1. Get the session ID from the request.
// 	// 2. Look up the kernel ID from the database/cache.
// 	// 3. Launch the two goroutines to handle the bidirectional message proxying.

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNotImplemented)
// 	w.Write([]byte(`{"error": "not implemented yet"}`))
// }

// // SaveSessionHandler is a handler function for the POST /api/sessions/save endpoint.
// func (h *Handlers) SaveSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	// TODO: Implement the logic to save a session.
// 	// 1. Deserialize the request body (the delta payload).
// 	// 2. Start a database transaction.
// 	// 3. Update the database with the changes.
// 	// 4. Commit the transaction and return a success response.

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNotImplemented)
// 	w.Write([]byte(`{"error": "not implemented yet"}`))
// }