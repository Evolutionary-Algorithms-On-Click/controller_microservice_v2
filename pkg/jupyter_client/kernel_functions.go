package jupyterclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

func (c *Client) StartKernel(ctx context.Context, language string) (*Kernel, error) {
	// hardcoded for temporary reasons
	knownKernels := map[string]bool{
		"python3": true,
		"tfenv":   true,
	}
	if !knownKernels[language] {
		return nil, fmt.Errorf("language '%s' is not a known kernel", language)
	}

	requestBody := StartKernelRequest{Name: language}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal start kernel request: %w", err)
	}

	url := c.baseURL + "/api/kernels"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create start kernel request: %w", err)
	}

	// Adding the required headers.
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute start kernel request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		var errRes ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return nil, fmt.Errorf("failed to decode error response with status %d: %w", res.StatusCode, err)
		}
		// Log the error for debugging purposes.
		pkg.Logger.Error().Err(errors.New(errRes.Message)).Msg("failed to start kernel")
		return nil, fmt.Errorf("failed to start kernel with status %d: %s", res.StatusCode, errRes.Message)
	}

	var kernelInfo Kernel
	if err := json.NewDecoder(res.Body).Decode(&kernelInfo); err != nil {
		return nil, fmt.Errorf("failed to decode successful kernel response: %w", err)
	}

	pkg.Logger.Info().Str("kernel_id", kernelInfo.ID).Msg("kernel started successfully")
	return &kernelInfo, nil
}

func (c *Client) GetKernels(ctx context.Context) (*[]Kernel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/kernels", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %d, unable to parse error body: %w", res.StatusCode, err)
		}
		return nil, fmt.Errorf("API returned non-200 status: %d, reason: %s", res.StatusCode, errResp.Reason)
	}

	var runningKernels []Kernel
	if err := json.NewDecoder(res.Body).Decode(&runningKernels); err != nil {
		return nil, fmt.Errorf("failed to decode kernel list response: %w", err)
	}

	return &runningKernels, nil
}
