package jupyterclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// struct for interacting with the Jupyter Kernel Gateway.
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new Jupyter Gateway client
func NewClient(baseURL string, token string) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("auth token cannot be empty")
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid base URL format: %s", err.Error())
	}
	parsedUrl, _ := url.Parse(baseURL)
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return nil, fmt.Errorf("base URL must use http or https protocol")
	}
	return &Client{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// GetGatewayURL returns the base URL of the Jupyter Gateway.
func (c *Client) GetGatewayURL() string {
	return c.baseURL
}

// GetAuthToken returns the authentication token for the Jupyter Gateway.
func (c *Client) GetAuthToken() string {
	return c.token
}


// CheckApiVersion sends a GET request to the /api endpoint to verify server connectivity.
// It returns an error if the request fails or the response status is not 200 OK.
func (c *Client) CheckApiVersion(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %s", resp.Status)
	}

	return nil
}
