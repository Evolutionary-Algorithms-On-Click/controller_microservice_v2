package controllers

import (
	"io"
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/middleware"
	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/rs/zerolog"
	"maps"
)

// LlmController holds the dependencies for the llm proxy handlers.
type LlmController struct {
	Module *modules.LlmModule
	Logger zerolog.Logger
}

// NewLlmController creates and returns a new LlmController.
func NewLlmController(module *modules.LlmModule, logger zerolog.Logger) *LlmController {
	return &LlmController{
		Module: module,
		Logger: logger,
	}
}

// GenerateNotebookHandler handles POST /api/v1/llm/generate
func (c *LlmController) GenerateNotebookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	resp, err := c.Module.GenerateNotebook(ctx, r.Body, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to proxy generate request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the proxied response to our response
	maps.Copy(w.Header(), resp.Header)

	// Write the status code from the proxied response
	w.WriteHeader(resp.StatusCode)

	// Stream the body from the proxied response
	if _, err := io.Copy(w, resp.Body); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to stream response body")
	}
}

// ModifyNotebookHandler handles POST /api/v1/llm/sessions/modify
func (c *LlmController) ModifyNotebookHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	resp, err := c.Module.ModifyNotebook(ctx, sessionID, r.Body, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("Failed to proxy modify request for session %s", sessionID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	maps.Copy(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to stream response body")
	}
}

// FixNotebookHandler handles POST /api/v1/llm/sessions/fix
func (c *LlmController) FixNotebookHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserContextKey).(*middleware.User)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	resp, err := c.Module.FixNotebook(ctx, sessionID, r.Body, user.ID)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("Failed to proxy fix request for session %s", sessionID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	maps.Copy(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to stream response body")
	}
}
