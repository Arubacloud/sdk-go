package restclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/internal/ports/interceptor"
	"github.com/Arubacloud/sdk-go/internal/ports/logger"
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
	// Redis configuration for token storage
	Redis *RedisConfig
	// File Basedir where are stored token json file
	File *FileConfig
	// Vault configuration for credentials retrieval
	Vault *VaultConfig
	// Logger is the logger to use for debug/info messages. If nil, no logging is performed.
	Logger logger.Logger
	// Debug enables debug logging when set to true
	Debug bool
}

// VaultConfig holds the configuration for Vault credentials retrieval
type VaultConfig struct {
	//address:port
	VaultURI  string
	KVMount   string
	KVPath    string
	Namespace string
	RolePath  string
	RoleID    string
	SecretID  string
}

// RedisConfig holds the configuration for Redis token storage
type RedisConfig struct {
	//"redis://<user>:<pass>@localhost:6379/<db>"
	RedisURI string
}

// FileConfig holds the configuration for file-based token storage
type FileConfig struct {
	//directory where stored token json files are located
	BaseDir string
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		HTTPClient:     http.DefaultClient,
		Headers:        make(map[string]string),
		BaseURL:        DefaultBaseURL,
		TokenIssuerURL: DefaultTokenIssuerURL,
		Redis: &RedisConfig{
			RedisURI: DefaultRedisURI,
		},
		File: &FileConfig{
			BaseDir: DefaultFileBaseDir,
		},
	}

}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.HTTPClient == nil {
		return fmt.Errorf("HTTPClient is required")
	}
	if c.TokenIssuerURL == "" {
		c.TokenIssuerURL = DefaultTokenIssuerURL
	}
	if c.ClientID == "" {
		return fmt.Errorf("ClientID is required for authentication")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("ClientSecret is required for authentication")
	}
	return nil
}

// Client is the main SDK client that aggregates all resource providers
type Client struct {
	config       *Config
	tokenManager *TokenManager
	logger       logger.Logger
	middleware   interceptor.Interceptor

	// Service interfaces for all API categories - these will be initialized by the services themselves
	// using a providers pattern to avoid import cycles
	services map[string]interface{}
}

// NewClient creates a new SDK client with the given configuration
func NewClient(config *Config, logger logger.Logger, middleware interceptor.Interceptor) (*Client, error) {
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
	)

	client := &Client{
		config:       config,
		tokenManager: tokenManager,
		logger:       logger,
		middleware:   middleware,
		services:     make(map[string]interface{}),
	}

	logger.Debugf("Initializing SDK client with base URL: %s", config.BaseURL)

	// Obtain initial token
	if err := tokenManager.ObtainToken(context.TODO()); err != nil {
		logger.Errorf("Failed to obtain initial token: %v", err)
		return nil, fmt.Errorf("failed to obtain initial token: %w", err)
	}

	logger.Debugf("Successfully obtained initial token")

	return client, nil
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// HTTPClient returns the HTTP client
func (c *Client) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

// Logger returns the client logger
func (c *Client) Logger() logger.Logger {
	return c.logger
}

// GetToken returns the current valid JWT token, refreshing if necessary
func (c *Client) GetToken(ctx context.Context) (string, error) {
	return c.tokenManager.GetToken(ctx)
}

// DoRequest performs an HTTP request with automatic authentication token injection
func (c *Client) DoRequest(ctx context.Context, method, path string, body io.Reader, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	// Build full URL
	url := c.config.BaseURL + path

	c.logger.Debugf("Making %s request to %s", method, url)

	// Read body for logging if present
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			c.logger.Errorf("Failed to read request body: %v", err)
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		c.logger.Debugf("Request body: %s", string(bodyBytes))
		// Recreate reader for actual request
		body = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		c.logger.Errorf("Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
		c.logger.Debugf("Added query parameters: %v", queryParams)
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add additional headers before authentication headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Log request headers (before auth)
	c.logger.Debugf("Request headers (pre-auth): %v", headers)

	// Use middleware
	if err := c.middleware.Intercept(ctx, req); err != nil {
		c.logger.Errorf("Failed to prepare request: %v", err)
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	// Log all headers after auth (excluding Authorization token for security)
	sanitizedHeaders := make(map[string]string)
	for key, values := range req.Header {
		if key == "Authorization" {
			sanitizedHeaders[key] = "Bearer [REDACTED]"
		} else {
			sanitizedHeaders[key] = values[0]
		}
	}
	c.logger.Debugf("Request headers (final): %v", sanitizedHeaders)

	// Execute request
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		c.logger.Errorf("Request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	c.logger.Debugf("Received response with status: %d %s", resp.StatusCode, resp.Status)

	// Log response headers
	c.logger.Debugf("Response headers: %v", resp.Header)

	// Log response body (for debugging)
	if resp.Body != nil {
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Warnf("Failed to read response body for logging: %v", err)
		} else {
			c.logger.Debugf("Response body: %s", string(respBodyBytes))
			// Recreate the response body so it can be read by the caller
			resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
		}
	}

	return resp, nil
}
