# Aruba Cloud Go SDK

A Go SDK for interacting with Aruba Cloud REST APIs. This SDK is automatically generated from OpenAPI/Swagger specifications and provides a clean, type-safe interface for all Aruba Cloud services.

## Overview

This SDK follows a microservices architecture where each Swagger file represents a resource provider (microservice). The SDK provides:

- **Type-safe API clients** generated from Swagger/OpenAPI specifications
- **Domain-specific interfaces** for each resource provider
- **Unified client** that aggregates all resource providers
- **Middleware support** for authentication, headers, and request customization
- **Context support** for cancellation and timeouts

## Project Structure

```
sdk-go/
├── swagger/              # OpenAPI/Swagger JSON specifications
│   └── network.json     # Example: Network resource provider spec
├── pkg/
│   ├── generated/       # Auto-generated code from Swagger files
│   │   └── network/     # Generated network client and types
│   └── client/          # SDK client implementation
│       ├── client.go    # Main SDK client
│       ├── providers.go # Resource provider wrappers
│       ├── middleware.go# Request middleware/interceptors
│       └── error.go     # Error handling
├── config/              # Code generation configurations
│   ├── types.yaml      # Config for generating types
│   └── client.yaml     # Config for generating clients
├── tools/              # Development tools
│   ├── go.mod          # Tools dependencies
│   └── tools.go        # Tools import
├── Makefile            # Build automation
└── README.md           # This file
```

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Quick Start

### 1. Generate Client Code

First, place your Swagger JSON files in the `swagger/` directory, then generate the client code:

```bash
# Generate all client code from Swagger files
make generate

# Or run all steps including mocks and tests
make all
```

### 2. Use the SDK

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Arubacloud/sdk-go/pkg/client"
)

func main() {
    // Create SDK configuration with OAuth2 client credentials
    config := &client.Config{
        BaseURL:        "https://api.arubacloud.com",
        TokenIssuerURL: "https://auth.arubacloud.com/oauth2/token",
        ClientID:       "your-client-id",
        ClientSecret:   "your-client-secret",
    }
    
    // Initialize the SDK client (automatically obtains JWT token)
    sdk, err := client.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use with context
    ctx := context.Background()
    sdk = sdk.WithContext(ctx)
    
    // The SDK automatically manages JWT token refresh
    // Use resource providers
    // Example: network := sdk.Network()
    // result, err := network.ListNetworks(ctx)
}
```

## Development

### Prerequisites

- Go 1.24 or higher
- Make

### Setup

1. Clone the repository:
```bash
git clone https://github.com/Arubacloud/sdk-go.git
cd sdk-go
```

2. Install dependencies:
```bash
go mod download
cd tools && go mod download && cd ..
```

### Makefile Commands

```bash
# Generate code from Swagger files
make generate

# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Run security checks
make sec

# Build the project
make build

# Generate mocks for testing
make mock

# Clean generated files
make clean

# Run all (generate, mock, build, test)
make all
```

### Adding a New Resource Provider

1. Add your Swagger JSON file to the `swagger/` directory:
```bash
cp my-service.json swagger/
```

2. Generate the client code:
```bash
make generate
```

3. This will create:
   - `pkg/generated/my-service/types.go` - Data types
   - `pkg/generated/my-service/client.go` - API client

4. Create a wrapper in `pkg/client/providers.go`:
```go
type MyServiceClient struct {
    client *myservice.ClientWithResponses
}

func (c *Client) MyService() *MyServiceClient {
    // Implementation
}
```

### Code Generation Configuration

The SDK uses `oapi-codegen` for generating code from OpenAPI/Swagger specifications:

- **config/types.yaml**: Configuration for generating data types
- **config/client.yaml**: Configuration for generating API clients

Both configurations are applied to each Swagger file in the `swagger/` directory.

## Authentication

The SDK uses **OAuth2 Client Credentials Flow** to obtain JWT Bearer tokens automatically.

### OAuth2 Client Credentials (Recommended)

The SDK automatically handles token acquisition and refresh:

```go
config := &client.Config{
    BaseURL:        "https://api.arubacloud.com",
    TokenIssuerURL: "https://auth.arubacloud.com/oauth2/token",
    ClientID:       "your-client-id",
    ClientSecret:   "your-client-secret",
}

sdk, err := client.NewClient(config)
// SDK automatically obtains and refreshes JWT tokens
```

### Token Management Features

- **Automatic token acquisition** on client initialization
- **Automatic token refresh** when tokens are about to expire
- **Thread-safe** token caching and refresh
- **Configurable refresh buffer** (default: 5 minutes before expiry)

### Advanced Configuration

```go
config := &client.Config{
    BaseURL:            "https://api.arubacloud.com",
    TokenIssuerURL:     "https://auth.arubacloud.com/oauth2/token",
    ClientID:           "your-client-id",
    ClientSecret:       "your-client-secret",
    TokenRefreshBuffer: 10 * time.Minute, // Refresh 10 min before expiry
    Headers: map[string]string{
        "X-Custom-Header": "value",
    },
}
```

### Manual Token Access

You can also manually access the current token if needed:

```go
token, err := sdk.GetToken(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Current JWT: %s\n", token)
```

## Error Handling

The SDK provides structured error handling:

```go
result, err := sdk.Network().GetNetwork(ctx, networkID)
if err != nil {
    if sdkErr, ok := err.(*client.Error); ok {
        fmt.Printf("Status: %d\n", sdkErr.StatusCode)
        fmt.Printf("Message: %s\n", sdkErr.Message)
        fmt.Printf("Body: %s\n", sdkErr.Body)
    }
    return err
}
```

## Testing

### Run Tests
```bash
make test
```

### Generate Mocks
```bash
make mock
```

Mocks are generated using mockery and placed in the `mock/` directory.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linters: `make lint test`
5. Commit your changes
6. Push to the branch
7. Create a Pull Request

## License

See [LICENSE](LICENSE) file for details.

## Architecture

### Resource Providers (Microservices)

Each Swagger file represents a separate microservice/resource provider:
- **Network Service** (`swagger/network.json`) - Network management
- **Compute Service** - VM and compute resources (add your swagger file)
- **Storage Service** - Storage management (add your swagger file)
- etc.

### Domain Interfaces

The SDK provides high-level domain interfaces that abstract the underlying API calls:

```go
// Main SDK client
type Client struct {
    config *Config
    // Resource provider clients
    networkClient *NetworkClient
    computeClient *ComputeClient
    // ... more providers
}

// Each provider has its own interface
type NetworkClient struct {
    client *network.ClientWithResponses
}
```

### Code Generation Flow

1. **Source**: Swagger JSON files in `swagger/` directory
2. **Generator**: `oapi-codegen` with custom configurations
3. **Output**: 
   - Types in `pkg/generated/{service}/types.go`
   - Client in `pkg/generated/{service}/client.go`
4. **Wrapper**: Hand-written wrappers in `pkg/client/providers.go`

## Support

For issues, questions, or contributions, please open an issue on GitHub.