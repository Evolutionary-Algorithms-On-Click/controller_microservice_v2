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

func NewClient(baseUrl string, authToken string) (*Client, error) {
	// 1. Validate input for baseUrl and authToken.
	if strings.TrimSpace(baseUrl) == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	if strings.TrimSpace(authToken) == "" {
		return nil, fmt.Errorf("auth token cannot be empty")
	}

	if _, err := url.ParseRequestURI(baseUrl); err != nil {
		return nil, fmt.Errorf("invalid base URL format: %s", err.Error())
	}
	parsedUrl, _ := url.Parse(baseUrl)
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return nil, fmt.Errorf("base URL must use http or https protocol")
	}

	c := &Client{
		baseURL: baseUrl,
		token:   authToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	if err := c.CheckApiVersion(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Kernel Gateway: %s", err.Error())
	}

	return c, nil
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
