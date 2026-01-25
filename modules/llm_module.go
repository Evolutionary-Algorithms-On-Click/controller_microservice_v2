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
func (m *LlmModule) GenerateNotebook(ctx context.Context, body io.Reader, userID string) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]any
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if err := IsNotebookIDPresent(requestData); err != nil {
		return nil, err
	}

	requestData["user_id"] = userID

	finalBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to re-encode request body: %w", err)
	}

	return m.Repo.GenerateNotebook(ctx, bytes.NewBuffer(finalBodyBytes))
}

// ModifyNotebook validates and proxies the modify request.
func (m *LlmModule) ModifyNotebook(ctx context.Context, sessionID string, body io.Reader, userID string) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]any
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if err := IsNotebookIDPresent(requestData); err != nil {
		return nil, err
	}

	requestData["user_id"] = userID

	if instruction, ok := requestData["instruction"].(string); !ok || instruction == "" {
		return nil, fmt.Errorf("request body must contain a non-empty 'instruction' string")
	}

	currentNotebook, ok := requestData["notebook"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("request body must contain a 'notebook' object")
	}
	if _, ok := currentNotebook["cells"].([]any); !ok {
		return nil, fmt.Errorf("'notebook' object must contain a 'cells' array")
	}

	finalBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to re-encode request body: %w", err)
	}

	return m.Repo.ModifyNotebook(ctx, bytes.NewBuffer(finalBodyBytes))
}

// FixNotebook validates and proxies the fix request.
func (m *LlmModule) FixNotebook(ctx context.Context, sessionID string, body io.Reader, userID string) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestData map[string]any
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body as JSON: %w", err)
	}

	if err := IsNotebookIDPresent(requestData); err != nil {
		return nil, err
	}

	requestData["user_id"] = userID

	if traceback, ok := requestData["traceback"].(string); !ok || traceback == "" {
		return nil, fmt.Errorf("request body must contain a non-empty 'traceback' string")
	}

	currentNotebook, ok := requestData["notebook"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("request body must contain a 'notebook' object")
	}
	if _, ok := currentNotebook["cells"].([]any); !ok {
		return nil, fmt.Errorf("'notebook' object must contain a 'cells' array")
	}

	finalBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to re-encode request body: %w", err)
	}

	return m.Repo.FixNotebook(ctx, bytes.NewBuffer(finalBodyBytes))
}

func IsNotebookIDPresent(requestData map[string]any) error {
	// making sure notebook_id is present
	if _, hasNotebookID := requestData["notebook_id"]; !hasNotebookID {
		return fmt.Errorf("request body must contain 'notebook_id'")
	}
	notebookIDStr, isString := requestData["notebook_id"].(string)
	if !isString || notebookIDStr == "" {
		return fmt.Errorf("'notebook_id' must be a non-empty string")
	}
	return nil
}
