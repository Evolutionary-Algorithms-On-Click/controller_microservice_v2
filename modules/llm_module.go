package modules

import (
	"context"
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

// GenerateNotebook calls the repository to proxy the generate request.
func (m *LlmModule) GenerateNotebook(ctx context.Context, body io.Reader) (*http.Response, error) {
	return m.Repo.GenerateNotebook(ctx, body)
}

// ModifyNotebook calls the repository to proxy the modify request.
func (m *LlmModule) ModifyNotebook(ctx context.Context, sessionID string, body io.Reader) (*http.Response, error) {
	return m.Repo.ModifyNotebook(ctx, sessionID, body)
}

// FixNotebook calls the repository to proxy the fix request.
func (m *LlmModule) FixNotebook(ctx context.Context, sessionID string, body io.Reader) (*http.Response, error) {
	return m.Repo.FixNotebook(ctx, sessionID, body)
}
