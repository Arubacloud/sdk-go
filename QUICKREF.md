# Aruba Cloud Go SDK - Quick Reference

## Overview
A Go SDK for Aruba Cloud REST APIs with automatic OAuth2 JWT token management.

## Key Features
✅ OAuth2 Client Credentials Flow with automatic token refresh
✅ Auto-generated clients from Swagger/OpenAPI JSON files
✅ Thread-safe token management
✅ Context support for request cancellation
✅ Comprehensive error handling
✅ Extensible architecture for multiple microservices

## Quick Start

### 1. Add Swagger Files
```bash
cp your-service.json swagger/
```

### 2. Generate Clients
```bash
make generate
```

### 3. Use the SDK
```go
config := &client.Config{
    BaseURL:        "https://api.arubacloud.com",
    TokenIssuerURL: "https://auth.arubacloud.com/oauth2/token",
    ClientID:       "your-client-id",
    ClientSecret:   "your-client-secret",
}

sdk, err := client.NewClient(config)
// SDK automatically manages JWT tokens
```

## File Structure
```
sdk-go/
├── swagger/                    # Swagger JSON files (input)
├── pkg/
│   ├── generated/             # Auto-generated code (git-ignored)
│   │   └── {service}/         # One per swagger file
│   │       ├── types.go       # Data types
│   │       └── client.go      # API client
│   └── client/                # SDK implementation
│       ├── client.go          # Main client
│       ├── token.go           # OAuth2 token manager
│       ├── providers.go       # Resource provider wrappers
│       ├── middleware.go      # Request helpers
│       └── error.go           # Error handling
├── config/                     # Code generation configs
├── tools/                      # Development tools
├── cmd/example/               # Usage examples
└── Makefile                   # Build automation
```

## Authentication Flow
1. Client initialization → Request token from TokenIssuerURL
2. Token stored with expiry time
3. Before each API call → Check if token valid
4. If expired/expiring → Automatically refresh
5. Add Bearer token to request headers

## Common Commands
```bash
make generate    # Generate code from Swagger files
make build       # Build the project
make test        # Run tests
make lint        # Run linters
make clean       # Clean generated files
make all         # Generate, build, and test
```

## Configuration Options
```go
type Config struct {
    BaseURL            string              // API base URL (required)
    TokenIssuerURL     string              // OAuth2 token endpoint (required)
    ClientID           string              // OAuth2 client ID (required)
    ClientSecret       string              // OAuth2 client secret (required)
    HTTPClient         *http.Client        // HTTP client (optional)
    Headers            map[string]string   // Custom headers (optional)
    TokenRefreshBuffer time.Duration       // Token refresh timing (default: 5min)
}
```

## Token Management
- **Automatic**: Token obtained on client initialization
- **Thread-safe**: Safe for concurrent use
- **Auto-refresh**: Refreshes before expiration
- **Configurable**: Adjust refresh buffer time
- **Manual access**: `token, err := sdk.GetToken(ctx)`

## Error Handling
```go
result, err := sdk.SomeService().SomeMethod(ctx)
if err != nil {
    if sdkErr, ok := err.(*client.Error); ok {
        log.Printf("Status: %d, Message: %s", sdkErr.StatusCode, sdkErr.Message)
    }
    return err
}
```

## Integration Steps
1. Place Swagger JSON in `swagger/`
2. Run `make generate`
3. Create wrapper in `pkg/client/`
4. Import generated package
5. Add provider to Client struct
6. Implement high-level methods

## Example Integration
```go
// pkg/client/network.go
import "github.com/Arubacloud/sdk-go/pkg/generated/network"

type NetworkClient struct {
    client *network.ClientWithResponses
}

func (c *Client) Network() *NetworkClient {
    client, _ := network.NewClientWithResponses(
        c.config.BaseURL,
        network.WithHTTPClient(c.config.HTTPClient),
        network.WithRequestEditorFn(c.RequestEditorFn()),
    )
    return &NetworkClient{client: client}
}
```

## Support
- README.md - Complete documentation
- DEVELOPMENT.md - Integration guide
- cmd/example/ - Usage examples

## Version
Go 1.24+
