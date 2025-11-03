# Aruba Cloud Go SDK - Development Guide

This guide explains how to develop and extend the SDK after generating client code from Swagger files.

## Initial Setup

1. **Add Swagger Files**: Place your OpenAPI/Swagger JSON files in the `swagger/` directory
2. **Generate Code**: Run `make generate` to create client code
3. **Review Generated Code**: Check `pkg/generated/` for the generated types and clients

## Integrating Generated Clients

### Step 1: Import Generated Package

After running `make generate`, import the generated package in `pkg/client/providers.go`:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/generated/network"
)
```

### Step 2: Add Provider Field to Client

Add a field for your resource provider in the `Client` struct:

```go
type Client struct {
    config         *Config
    ctx            context.Context
    tokenManager   *TokenManager
    networkClient  *NetworkClient  // Add this
}
```

### Step 3: Create Provider Wrapper

Create a wrapper type for your resource provider:

```go
type NetworkClient struct {
    client *network.ClientWithResponses
    sdk    *Client
}
```

### Step 4: Add Accessor Method

Add a method to access your provider:

```go
func (c *Client) Network() *NetworkClient {
    if c.networkClient == nil {
        client, _ := network.NewClientWithResponses(
            c.config.BaseURL,
            network.WithHTTPClient(c.config.HTTPClient),
            network.WithRequestEditorFn(c.RequestEditorFn()),
        )
        c.networkClient = &NetworkClient{
            client: client,
            sdk:    c,
        }
    }
    return c.networkClient
}
```

### Step 5: Create High-Level Methods

Wrap the generated methods with high-level, user-friendly methods:

```go
func (n *NetworkClient) ListNetworks(ctx context.Context) ([]*network.Network, error) {
    resp, err := n.client.ListNetworksWithResponse(ctx)
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode() != 200 {
        return nil, NewError(
            resp.StatusCode(),
            "failed to list networks",
            resp.Body,
            nil,
        )
    }
    
    return resp.JSON200.Networks, nil
}

func (n *NetworkClient) GetNetwork(ctx context.Context, id string) (*network.Network, error) {
    resp, err := n.client.GetNetworkWithResponse(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode() != 200 {
        return nil, NewError(
            resp.StatusCode(),
            fmt.Sprintf("failed to get network %s", id),
            resp.Body,
            nil,
        )
    }
    
    return resp.JSON200, nil
}
```

## Example: Complete Integration

Here's a complete example for the Network service:

```go
// File: pkg/client/network.go
package client

import (
    "context"
    "fmt"
    
    "github.com/Arubacloud/sdk-go/pkg/generated/network"
)

type NetworkClient struct {
    client *network.ClientWithResponses
    sdk    *Client
}

func (c *Client) Network() *NetworkClient {
    if c.networkClient == nil {
        client, _ := network.NewClientWithResponses(
            c.config.BaseURL,
            network.WithHTTPClient(c.config.HTTPClient),
            network.WithRequestEditorFn(c.RequestEditorFn()),
        )
        c.networkClient = &NetworkClient{
            client: client,
            sdk:    c,
        }
    }
    return c.networkClient
}

func (n *NetworkClient) ListNetworks(ctx context.Context, params *network.ListNetworksParams) ([]*network.Network, error) {
    resp, err := n.client.ListNetworksWithResponse(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to list networks: %w", err)
    }
    
    if resp.StatusCode() != 200 {
        return nil, NewError(resp.StatusCode(), "failed to list networks", resp.Body, nil)
    }
    
    return resp.JSON200.Networks, nil
}

func (n *NetworkClient) CreateNetwork(ctx context.Context, req network.CreateNetworkRequest) (*network.Network, error) {
    resp, err := n.client.CreateNetworkWithResponse(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create network: %w", err)
    }
    
    if resp.StatusCode() != 201 {
        return nil, NewError(resp.StatusCode(), "failed to create network", resp.Body, nil)
    }
    
    return resp.JSON201, nil
}
```

## Testing Your Integration

Create tests for your wrapper:

```go
// File: pkg/client/network_test.go
package client

import (
    "context"
    "testing"
)

func TestNetworkClient_ListNetworks(t *testing.T) {
    // Setup mock server
    tokenServer := setupMockTokenServer(t)
    defer tokenServer.Close()
    
    // Create SDK client
    config := &Config{
        BaseURL:        "https://api.example.com",
        TokenIssuerURL: tokenServer.URL,
        ClientID:       "test-id",
        ClientSecret:   "test-secret",
    }
    
    sdk, err := NewClient(config)
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }
    
    // Test your methods
    ctx := context.Background()
    networks, err := sdk.Network().ListNetworks(ctx, nil)
    // Add assertions...
}
```

## Best Practices

1. **Error Handling**: Always wrap errors with context
2. **Context Propagation**: Accept `context.Context` as the first parameter
3. **Nil Checks**: Check for nil responses before accessing fields
4. **Type Safety**: Use the generated types from the Swagger schemas
5. **Documentation**: Add godoc comments to all public methods
6. **Validation**: Validate input parameters before making API calls
7. **Logging**: Consider adding logging for debugging (optional)

## Code Generation Workflow

```bash
# 1. Add or update Swagger files
cp new-service.json swagger/

# 2. Generate client code
make generate

# 3. Update pkg/client/providers.go with new integration

# 4. Run tests
make test

# 5. Format and lint
make lint

# 6. Build
make build
```

## Continuous Integration

The generated code should not be committed to git. Add to `.gitignore`:

```gitignore
pkg/generated/
```

In CI/CD, always run `make generate` before building:

```yaml
# Example GitHub Actions
steps:
  - name: Generate code
    run: make generate
  - name: Test
    run: make test
  - name: Build
    run: make build
```
