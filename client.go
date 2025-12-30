package atompub

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient    *http.Client
	authenticator Authenticator
	userAgent     string
}

// NewClient creates a new AtomPub client.
func NewClient(auth Authenticator) *Client {
	if auth == nil {
		auth = &NoAuth{}
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authenticator: auth,
		userAgent:     "go-atompub",
	}
	return client
}

// doRequest executes an HTTP request
func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	if err := c.authenticator.Authenticate(req); err != nil {
		return nil, fmt.Errorf("authenticate request: %w", err)
	}
	return c.httpClient.Do(req)
}

// GetServiceDocument retrieves and parses a Service Document
func (c *Client) GetServiceDocument(ctx context.Context, url string) (*ServiceDocument, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sdoc ServiceDocument
	if err := xml.NewDecoder(resp.Body).Decode(&sdoc); err != nil {
		return nil, fmt.Errorf("decode service document: %w", err)
	}
	return &sdoc, nil
}
