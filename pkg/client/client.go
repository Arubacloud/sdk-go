package client

import (
	"context"
	"fmt"
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

	// Compute APIs
	CloudServer schema.CloudServerAPI
	KaaS        schema.KaaSAPI

	// Network APIs
	Vpc        schema.VpcAPI
	Subnet     schema.SubnetAPI
	VpcPeering schema.VpcPeeringAPI
	VpcRoute   schema.VpcRouteAPI
	VpnTunnel  schema.VpnTunnelAPI
	ElasticIp  schema.ElasticIpAPI

	// Security APIs
	SecurityGroup schema.SecurityGroupAPI

	// Storage APIs
	BlockStorage schema.BlockStorageAPI
	Snapshot     schema.SnapshotAPI

	// Database APIs
	DBaaS schema.DBaaSAPI

	// Schedule APIs
	ScheduleJob schema.ScheduleJobAPI
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
