package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds the configuration for the SDK client
type Config struct {
	// BaseURL is the base URL for all API calls
	BaseURL string
	// HTTPClient is the HTTP client to use for requests
	HTTPClient *http.Client
	// TokenIssuerURL is the URL to obtain JWT tokens (OAuth2 token endpoint)
	TokenIssuerURL string
	// ClientID for OAuth2 client credentials flow
	ClientID string
	// ClientSecret for OAuth2 client credentials flow
	ClientSecret string
	// Headers are additional headers to include in all requests
	Headers map[string]string
	// TokenRefreshBuffer is the time before expiry to refresh the token (default: 5 minutes)
	TokenRefreshBuffer time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		HTTPClient:         http.DefaultClient,
		Headers:            make(map[string]string),
		TokenRefreshBuffer: 5 * time.Minute,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("BaseURL is required")
	}
	if c.HTTPClient == nil {
		return fmt.Errorf("HTTPClient is required")
	}
	if c.TokenIssuerURL == "" {
		return fmt.Errorf("TokenIssuerURL is required for authentication")
	}
	if c.ClientID == "" {
		return fmt.Errorf("ClientID is required for authentication")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("ClientSecret is required for authentication")
	}
	if c.TokenRefreshBuffer == 0 {
		c.TokenRefreshBuffer = 5 * time.Minute
	}
	return nil
}

// Client is the main SDK client that aggregates all resource providers
type Client struct {
	config       *Config
	ctx          context.Context
	tokenManager *TokenManager
}

// NewClient creates a new SDK client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize token manager
	tokenManager := NewTokenManager(
		config.TokenIssuerURL,
		config.ClientID,
		config.ClientSecret,
		config.HTTPClient,
		config.TokenRefreshBuffer,
	)

	client := &Client{
		config:       config,
		ctx:          context.Background(),
		tokenManager: tokenManager,
	}

	// Fetch initial token
	if err := tokenManager.RefreshToken(client.ctx); err != nil {
		return nil, fmt.Errorf("failed to obtain initial token: %w", err)
	}

	// Initialize all API clients - these will be created in respective packages
	// For example: client.CloudServer = compute.NewCloudServerClient(client)
	// TODO: Initialize API clients once implementations are ready

	return client, nil
}

// WithContext returns a new client with the given context
func (c *Client) WithContext(ctx context.Context) *Client {
	newClient := *c
	newClient.ctx = ctx
	return &newClient
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// Context returns the client context
func (c *Client) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

// HTTPClient returns the HTTP client
func (c *Client) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

// GetToken returns the current valid JWT token, refreshing if necessary
func (c *Client) GetToken(ctx context.Context) (string, error) {
	return c.tokenManager.GetToken(ctx)
}

// DoRequest performs an HTTP request with automatic authentication token injection
func (c *Client) DoRequest(ctx context.Context, method, path string, body io.Reader, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	// Build full URL
	url := c.config.BaseURL + path

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add additional headers before authentication headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Use RequestEditorFn to add authentication and custom headers from config
	editorFn := c.RequestEditorFn()
	if err := editorFn(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	// Execute request
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// RequestEditorFn returns a function that adds the Bearer token to requests
func (c *Client) RequestEditorFn() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		token, err := c.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Add custom headers
		for k, v := range c.config.Headers {
			req.Header.Set(k, v)
		}

		return nil
	}
}
