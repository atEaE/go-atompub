package atompub

import (
	"bytes"
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
	verbose       bool
}

// ClientOption defines a function type for configuring the Client
type ClientOption func(*Client)

// WithVerbose enables or disables verbose logging
func WithVerbose(verbose bool) ClientOption {
	return func(c *Client) {
		c.verbose = verbose
	}
}

// WithUserAgent sets a custom User-Agent header for the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// NewClient creates a new AtomPub client.
func NewClient(auth Authenticator, opts ...ClientOption) *Client {
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

	for _, opt := range opts {
		opt(client)
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("do request: %w", err)
	}
	if c.verbose {
		if err := c.dumpResponseBody(resp); err != nil {
			return resp, fmt.Errorf("dump response body: %w", err)
		}
	}
	return resp, nil
}

// dumpResponseBody prints the response body for debugging purposes
func (c *Client) dumpResponseBody(resp *http.Response) error {
	if resp.Body == nil {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	fmt.Printf("Response Body:\n%s\n", string(body))

	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

// GetServiceDocument retrieves and parses a Service Document (RFC 5023 8.  Service Documents)
func (c *Client) GetServiceDocument(ctx context.Context, url string) (*GetServiceDocumentResponse, error) {
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
	return &GetServiceDocumentResponse{
		Headers: resp.Header,
		Body:    sdoc,
	}, nil
}

type GetServiceDocumentResponse struct {
	Headers http.Header
	Body    ServiceDocument
}

// CreateEntry creates a new entry in a collection
func (c *Client) CreateEntry(ctx context.Context, collectionURL string, entry *Entry) (*CreateEntryResponse, error) {
	body, err := xml.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("marshal entry: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, collectionURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdEntry Entry
	if err := xml.NewDecoder(resp.Body).Decode(&createdEntry); err != nil {
		return nil, fmt.Errorf("decode created entry: %w", err)
	}
	return &CreateEntryResponse{
		Headers: resp.Header,
		Body:    createdEntry,
	}, nil
}

type CreateEntryResponse struct {
	Headers http.Header
	Body    Entry
}
