package repository

import (
	"context"
	"fmt"

	"io"
	"net/http"
)

// llmProxy implements the LlmRepository interface by acting as a proxy to the Python LLM service.
type llmProxy struct {
	BaseURL string
	Client  *http.Client
}

type LlmRepository interface {
	GenerateNotebook(
		ctx context.Context,
		body io.Reader,
	) (*http.Response, error)
	ModifyNotebook(
		ctx context.Context,
		body io.Reader,
	) (*http.Response, error)
	FixNotebook(
		ctx context.Context,
		body io.Reader,
	) (*http.Response, error)
}

// NewLlmProxy creates a new llmProxy.
func NewLlmProxy(baseURL string) LlmRepository {
	return &llmProxy{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// GenerateNotebook proxies the request to the /generate endpoint of the LLM service.
func (p *llmProxy) GenerateNotebook(
	ctx context.Context,
	body io.Reader,
) (*http.Response, error) {
	targetURL := fmt.Sprintf("%s/v1/generate", p.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create generate request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call generate endpoint: %w", err)
	}

	return resp, nil
}

// ModifyNotebook proxies the request to the /modify endpoint.
func (p *llmProxy) ModifyNotebook(
	ctx context.Context,
	body io.Reader,
) (*http.Response, error) {
	targetURL := fmt.Sprintf("%s/v1/modify", p.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create modify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call modify endpoint: %w", err)
	}

	return resp, nil
}

// FixNotebook proxies the request to the /sessions/{session_id}/fix endpoint.
func (p *llmProxy) FixNotebook(
	ctx context.Context,
	body io.Reader,
) (*http.Response, error) {
	targetURL := fmt.Sprintf("%s/v1/fix", p.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create fix request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call fix endpoint: %w", err)
	}

	return resp, nil
}
