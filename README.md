# Aruba Cloud Go SDK

A Go SDK for interacting with Aruba Cloud REST APIs. This SDK provides a clean, type-safe interface for all Aruba Cloud services following the oapi-codegen style.

## Overview

This SDK follows a microservices architecture where each resource type has its own API client. The SDK provides:

- **Type-safe API clients** with generic `Response[T]` wrapper (parsed data + HTTP metadata)
- **Domain-specific interfaces** for each resource type
- **Unified client** that aggregates all resource providers
- **Automatic JWT authentication** via OAuth2 client credentials flow
- **Request editor support** for middleware and customization
- **Context support** for cancellation and timeouts

## Project Structure

```
sdk-go/
├── cmd/
│   └── example/             # Example usage
│       ├── main.go
│       └── README.md
├── pkg/
│   ├── client/              # Main SDK client
│   │   ├── client.go        # Client with all API providers
│   │   ├── client_test.go
│   │   ├── error.go         # Error handling
│   │   ├── middleware.go    # Request middleware
│   │   ├── params.go        # Parameter helpers
│   │   ├── providers.go     # Service providers
│   │   ├── token.go         # OAuth2 token manager
│   │   └── token_test.go
│   └── spec/
│       ├── schema/          # Shared types and schemas
│       │   ├── resource.go  # Generic Response[T] wrapper
│       │   ├── parameters.go # Common parameters
│       │   ├── error.go     # Error types
│       │   └── *.go         # Resource schemas
│       ├── audit/           # Audit service
│       │   ├── interface.go
│       │   ├── event.go
│       │   ├── path.go
│       │   └── README.md
│       ├── compute/         # Compute service
│       │   ├── interface.go
│       │   ├── cloudserver.go
│       │   ├── keypair.go
│       │   ├── path.go
│       │   └── README.md
│       ├── database/        # Database service
│       │   ├── interface.go
│       │   ├── dbaas.go
│       │   ├── database.go
│       │   ├── user.go
│       │   ├── grant.go
│       │   ├── backup.go
│       │   ├── path.go
│       │   └── README.md
│       ├── metric/          # Metrics service
│       │   ├── interface.go
│       │   ├── metric.go
│       │   ├── alert.go
│       │   ├── path.go
│       │   └── README.md
│       ├── network/         # Network service
│       │   ├── interface.go
│       │   ├── vpc.go
│       │   ├── subnet.go
│       │   ├── elastic-ip.go
│       │   ├── load-balancer.go
│       │   ├── security-group.go
│       │   ├── security-group-rule.go
│       │   ├── vpc-peering.go
│       │   ├── vpc-peering-route.go
│       │   ├── vpn-tunnel.go
│       │   ├── vpn-route.go
│       │   ├── path.go
│       │   └── README.md
│       ├── project/         # Project service
│       │   ├── interface.go
│       │   ├── path.go
│       │   └── README.md
│       ├── schedule/        # Schedule service
│       │   ├── interface.go
│       │   ├── job.go
│       │   ├── path.go
│       │   └── README.md
│       ├── security/        # Security service
│       │   ├── interface.go
│       │   ├── kms.go
│       │   ├── path.go
│       │   └── README.md
│       └── storage/         # Storage service
│           ├── interface.go
│           ├── block-storage.go
│           ├── snapshot.go
│           ├── path.go
│           └── README.md
├── tools/                   # Build tools
│   ├── go.mod
│   └── tools.go
├── go.mod
├── Makefile
├── README.md
├── DEVELOPMENT.md
├── FILTERS.md               # Filtering documentation
├── OAUTH2.md                # OAuth2 documentation
├── QUICKREF.md              # Quick reference
└── SDK_READY.md
```

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Quick Start

### Initialize the SDK Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
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
    
    ctx := context.Background()
    
    // Use the SDK
    listCloudServers(ctx, sdk)
}
```

### Working with Cloud Servers

```go
func listCloudServers(ctx context.Context, sdk *client.Client) {
    // List cloud servers
    resp, err := sdk.CloudServer.ListCloudServers(
        ctx,
        "my-project",
        &schema.ListParams{
            Limit: ptrInt(10),
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Check response status
    if !resp.IsSuccess() {
        log.Fatalf("Failed to list cloud servers: %d", resp.StatusCode)
    }
    
    // Access parsed data
    fmt.Printf("Found %d cloud servers\n", len(resp.Data.Values))
    for _, server := range resp.Data.Values {
        fmt.Printf("- %s (%s)\n", server.Metadata.Name, server.Status.Phase)
    }
}

func createCloudServer(ctx context.Context, sdk *client.Client) {
    // Create a cloud server
    req := schema.CloudServerRequest{
        Spec: schema.CloudServerSpec{
            FlavorName: "small",
            ImageName:  "ubuntu-22.04",
            Zone:       ptrString("it-mil1"),
        },
    }
    
    resp, err := sdk.CloudServer.CreateOrUpdateCloudServer(
        ctx,
        "my-project",
        "my-server",
        nil, // no conditional params
        req,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Cloud server created: %s\n", resp.Data.Metadata.Name)
    }
}

func ptrInt(i int) *schema.LimitParam {
    v := schema.LimitParam(i)
    return &v
}

func ptrString(s string) *string {
    return &s
}
```

### Working with VPCs

```go
func manageVPCs(ctx context.Context, sdk *client.Client) {
    // List VPCs
    resp, err := sdk.Vpc.ListVpcs(
        ctx,
        "my-project",
        &schema.ListParams{
            Labels: ptrLabel("environment=prod"),
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        for _, vpc := range resp.Data.Values {
            fmt.Printf("VPC: %s - %s\n", vpc.Metadata.Name, vpc.Spec.CidrBlock)
        }
    }
    
    // Create a VPC
    vpcReq := schema.VpcRequest{
        Spec: schema.VpcSpec{
            CidrBlock: "10.0.0.0/16",
        },
    }
    
    resp, err = sdk.Vpc.CreateOrUpdateVpc(
        ctx,
        "my-project",
        "my-vpc",
        nil,
        vpcReq,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("VPC created: %s\n", resp.Data.Metadata.Name)
    }
}

func ptrLabel(s string) *schema.LabelSelector {
    v := schema.LabelSelector(s)
    return &v
}
```

### Working with Block Storage

```go
func manageStorage(ctx context.Context, sdk *client.Client) {
    // Create block storage
    storageReq := schema.BlockStorageRequest{
        Spec: schema.BlockStorageSpec{
            SizeGB: 100,
            Type:   ptrString("ssd"),
        },
    }
    
    resp, err := sdk.BlockStorage.CreateOrUpdateBlockStorage(
        ctx,
        "my-project",
        "my-disk",
        nil,
        storageReq,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Block storage created: %s (Size: %dGB)\n", 
            resp.Data.Metadata.Name, resp.Data.Spec.SizeGB)
    }
    
    // Create snapshot
    snapshotReq := schema.SnapshotRequest{
        Spec: schema.SnapshotSpec{
            SourceVolume: "my-disk",
        },
    }
    
    resp, err = sdk.Snapshot.CreateOrUpdateSnapshot(
        ctx,
        "my-project",
        "my-snapshot",
        nil,
        snapshotReq,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Snapshot created: %s\n", resp.Data.Metadata.Name)
    }
}
```

## Available Resources

### Compute
- **CloudServer** - Virtual machine management

### Network
- **Vpc** - Virtual Private Cloud
- **Subnet** - Subnet management
- **ElasticIp** - Elastic IP addresses
- **VpcPeering** - VPC peering connections
- **VpcRoute** - Custom routes
- **VpnTunnel** - VPN tunnels

### Security
- **SecurityGroup** - Security group rules

### Storage
- **BlockStorage** - Block storage volumes
- **Snapshot** - Volume snapshots

### Database
- **DBaaS** - Database as a Service

### Schedule
- **ScheduleJob** - Scheduled tasks

## Authentication

The SDK uses **OAuth2 Client Credentials Flow** to obtain JWT Bearer tokens automatically.

### Automatic Token Management

```go
config := &client.Config{
    BaseURL:        "https://api.arubacloud.com",
    TokenIssuerURL: "https://auth.arubacloud.com/oauth2/token",
    ClientID:       "your-client-id",
    ClientSecret:   "your-client-secret",
}

sdk, err := client.NewClient(config)
// SDK automatically:
// - Obtains initial JWT token on initialization
// - Refreshes token before expiry
// - Adds "Authorization: Bearer <token>" to all requests
```

### Token Management Features

- ✅ **Automatic token acquisition** on client initialization
- ✅ **Automatic token refresh** when tokens are about to expire
- ✅ **Thread-safe** token caching and refresh
- ✅ **Configurable refresh buffer** (default: 5 minutes before expiry)
- ✅ **Custom headers** support

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
        "User-Agent":      "my-app/1.0",
    },
}
```

## Request Customization

### Using Request Editors

```go
// Add custom headers or modify request
customEditor := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Request-ID", "123456")
    return nil
}

resp, err := sdk.CloudServer.GetCloudServer(
    ctx,
    "my-project",
    "my-server",
    customEditor, // pass request editor
)

if resp.IsSuccess() {
    fmt.Printf("Server: %s\n", resp.Data.Metadata.Name)
}
```

### Conditional Requests

```go
// Use if-unmodified-since for optimistic concurrency
params := &schema.CreateOrUpdateParams{
    IfUnmodifiedSince: ptrString("2024-11-03T10:00:00Z"),
}

resp, err := sdk.CloudServer.CreateOrUpdateCloudServer(
    ctx,
    "my-project",
    "my-server",
    params,
    request,
)
```

### Filtering Deleted Resources

```go
// Include deleted resources in list
accept := schema.AcceptHeaderJSONDeletedTrue
params := &schema.ListParams{
    Accept: &accept,
}

resp, err := sdk.Vpc.ListVpcs(ctx, "my-project", params)
if resp.IsSuccess() {
    for _, vpc := range resp.Data.Values {
        fmt.Printf("VPC: %s (Deleted: %v)\n", vpc.Metadata.Name, vpc.Metadata.DeletionTimestamp != nil)
    }
}
```

### Filtering by Labels

```go
// Filter resources by labels
labelSelector := schema.LabelSelector("environment=prod,tier=frontend")
params := &schema.ListParams{
    Labels: &labelSelector,
}

resp, err := sdk.CloudServer.ListCloudServers(ctx, "my-project", params)
if resp.IsSuccess() {
    for _, server := range resp.Data.Values {
        fmt.Printf("Server: %s - Labels: %v\n", 
            server.Metadata.Name, 
            server.Metadata.Labels)
    }
}
```

### Filtering by Field Selectors

```go
// Filter by field values
fieldSelector := schema.FieldSelector("status.phase=Running")
params := &schema.ListParams{
    Fields: &fieldSelector,
}

resp, err := sdk.CloudServer.ListCloudServers(ctx, "my-project", params)
if resp.IsSuccess() {
    fmt.Printf("Found %d running servers\n", len(resp.Data.Values))
}
```

### Pagination

```go
// Paginate through results
limit := schema.LimitParam(50)
params := &schema.ListParams{
    Limit: &limit,
}

resp, err := sdk.Vpc.ListVpcs(ctx, "my-project", params)
if resp.IsSuccess() {
    fmt.Printf("Page 1: %d VPCs\n", len(resp.Data.Values))
    
    // Get next page using continuation token
    if resp.Data.Metadata.Continue != nil {
        params.Continue = resp.Data.Metadata.Continue
        resp, err = sdk.Vpc.ListVpcs(ctx, "my-project", params)
        if resp.IsSuccess() {
            fmt.Printf("Page 2: %d VPCs\n", len(resp.Data.Values))
        }
    }
}
```

## Error Handling

All API methods return a generic `Response[T]` wrapper and `error`. The response includes both parsed data and HTTP metadata:

```go
resp, err := sdk.CloudServer.GetCloudServer(ctx, "my-project", "my-server")
if err != nil {
    log.Printf("Request failed: %v", err)
    return err
}

// Check status using helper methods
if resp.IsError() {
    log.Printf("API error: %d - %s", resp.StatusCode, string(resp.RawBody))
    return fmt.Errorf("unexpected status: %d", resp.StatusCode)
}

if resp.IsSuccess() {
    // Access parsed data directly
    fmt.Printf("Server: %s (Status: %s)\n", 
        resp.Data.Metadata.Name, 
        resp.Data.Status.Phase)
    
    // Access HTTP metadata if needed
    fmt.Printf("Response headers: %v\n", resp.Headers)
}
```

### Response[T] Structure

The `Response[T]` type provides:
- `Data *T` - Parsed JSON response data
- `HTTPResponse *http.Response` - Full HTTP response
- `StatusCode int` - HTTP status code
- `Headers http.Header` - Response headers
- `RawBody []byte` - Raw response body
- `IsSuccess() bool` - Returns true for 2xx status codes
- `IsError() bool` - Returns true for 4xx/5xx status codes

## API Interface Pattern

All resource clients follow the same pattern and return `Response[T]`:

```go
type ResourceAPI interface {
    // List resources with pagination and filtering
    List{Resource}(ctx, project, params, ...editors) (*schema.Response[schema.{Resource}ListResponse], error)
    
    // Get a single resource by name
    Get{Resource}(ctx, project, name, ...editors) (*schema.Response[schema.{Resource}Response], error)
    
    // Create or update a resource (upsert)
    CreateOrUpdate{Resource}(ctx, project, name, params, body, ...editors) (*schema.Response[schema.{Resource}Response], error)
    CreateOrUpdate{Resource}WithBody(ctx, project, name, params, contentType, body, ...editors) (*schema.Response[schema.{Resource}Response], error)
    
    // Delete a resource
    Delete{Resource}(ctx, project, name, params, ...editors) (*schema.Response[any], error)
}
```

## Context Support

All operations support context for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := sdk.CloudServer.ListCloudServers(ctx, "my-project", nil)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel after 5 seconds
}()

resp, err := sdk.CloudServer.GetCloudServer(ctx, "my-project", "my-server")
```

## Best Practices

1. **Check response status using helper methods**:
   ```go
   resp, err := sdk.CloudServer.GetCloudServer(ctx, project, name)
   if err != nil {
       return err
   }
   
   if resp.IsSuccess() {
       // Access parsed data
       fmt.Println(resp.Data.Metadata.Name)
   }
   ```

2. **Access both parsed data and HTTP metadata**:
   ```go
   if resp.IsSuccess() {
       // Use parsed data
       for _, item := range resp.Data.Values {
           fmt.Println(item.Metadata.Name)
       }
       
       // Access HTTP details if needed
       fmt.Printf("Status: %d, Headers: %v\n", resp.StatusCode, resp.Headers)
   }
   ```

3. **Use contexts with timeouts**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

4. **Reuse the SDK client**:
   ```go
   // Create once
   sdk, err := client.NewClient(config)
   // Reuse for all operations
   ```

5. **Handle token refresh errors gracefully**:
   The SDK automatically refreshes tokens, but network issues may occur.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Format code: `go fmt ./...`
6. Commit your changes
7. Push to the branch
8. Create a Pull Request

## License

See [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please open an issue on GitHub.