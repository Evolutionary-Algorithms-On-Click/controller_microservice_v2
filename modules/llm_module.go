package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
)

// LlmModule encapsulates the business logic for proxying requests to the LLM service.
type LlmModule struct {
	Repo repository.LlmRepository
}

// NewLlmModule creates and returns a new LlmModule.
func NewLlmModule(repo repository.LlmRepository) *LlmModule {
	return &LlmModule{
		Repo: repo,
	}
}

// GenerateNotebook validates and proxies the generate request.
func (m *LlmModule) GenerateNotebook(ctx context.Context, body io.Reader) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if _, hasSessionID := requestData["session_id"]; !hasSessionID {
		notebookID, hasNotebookID := requestData["notebook_id"]
		if !hasNotebookID {
			return nil, fmt.Errorf("request body must contain 'session_id' or 'notebook_id'")
		}

		notebookIDStr, isString := notebookID.(string)
		if !isString || notebookIDStr == "" {
			return nil, fmt.Errorf("'notebook_id' must be a non-empty string")
		}
		requestData["session_id"] = notebookIDStr
	}

	finalBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to re-encode request body: %w", err)
	}

	return m.Repo.GenerateNotebook(ctx, bytes.NewBuffer(finalBodyBytes))
}

// ModifyNotebook validates and proxies the modify request.
func (m *LlmModule) ModifyNotebook(ctx context.Context, sessionID string, body io.Reader) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if instruction, ok := requestData["instruction"].(string); !ok || instruction == "" {
		return nil, fmt.Errorf("request body must contain a non-empty 'instruction' string")
	}

	currentNotebook, ok := requestData["current_notebook"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("request body must contain a 'current_notebook' object")
	}
	if _, ok := currentNotebook["cells"].([]interface{}); !ok {
		return nil, fmt.Errorf("'current_notebook' object must contain a 'cells' array")
	}

	return m.Repo.ModifyNotebook(ctx, sessionID, bytes.NewBuffer(bodyBytes))
}

// FixNotebook validates and proxies the fix request.
func (m *LlmModule) FixNotebook(ctx context.Context, sessionID string, body io.Reader) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if traceback, ok := requestData["traceback"].(string); !ok || traceback == "" {
		return nil, fmt.Errorf("request body must contain a non-empty 'traceback' string")
	}

	currentNotebook, ok := requestData["current_notebook"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("request body must contain a 'current_notebook' object")
	}
	if _, ok := currentNotebook["cells"].([]interface{}); !ok {
		return nil, fmt.Errorf("'current_notebook' object must contain a 'cells' array")
	}

	return m.Repo.FixNotebook(ctx, sessionID, bytes.NewBuffer(bodyBytes))
}
