package compute

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// Client provides access to all compute resources
type Client struct {
	client *client.Client

	// Compute resource clients
	CloudServer schema.CloudServerAPI
	KaaS        schema.KaaSAPI
}

// NewClient creates a new compute client
func NewClient(c *client.Client) *Client {
	return &Client{
		client:      c,
		CloudServer: NewCloudServerClient(c),
		KaaS:        NewKaaSClient(c),
	}
}

// Config returns the underlying SDK configuration
func (c *Client) Config() *client.Config {
	return c.client.Config()
}

// HTTPClient returns the HTTP client
func (c *Client) HTTPClient() *client.Client {
	return c.client
}
