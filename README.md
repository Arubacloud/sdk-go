# Aruba Cloud Go SDK

A Go SDK for interacting with Aruba Cloud REST APIs. This SDK provides a clean, type-safe interface for all Aruba Cloud services following the oapi-codegen style.

## Overview

This SDK follows a microservices architecture where each resource type has its own API client. The SDK provides:

- **Type-safe API clients** with HTTP-level operations (returns `*http.Response`)
- **Domain-specific interfaces** for each resource type
- **Unified client** that aggregates all resource providers
- **Automatic JWT authentication** via OAuth2 client credentials flow
- **Request editor support** for middleware and customization
- **Context support** for cancellation and timeouts

## Project Structure

```
sdk-go/
├── pkg/
│   ├── client/              # Main SDK client
│   │   ├── client.go        # Client with all API providers
│   │   └── token.go         # OAuth2 token manager
│   ├── spec/
│   │   ├── schema/          # Shared types and interfaces
│   │   │   ├── api.go       # API interfaces (VpcAPI, SubnetAPI, etc.)
│   │   │   ├── types.go     # Common types and request/response structs
│   │   │   └── common.go    # Common parameters and enums
│   │   ├── compute/         # Compute service implementations
│   │   │   ├── cloudserver.go # CloudServer API client
│   │   │   └── kaas.go        # Kubernetes API client
│   │   ├── network/         # Network service implementations
│   │   │   ├── vpc.go
│   │   │   ├── subnet.go
│   │   │   ├── elasticip.go
│   │   │   ├── vpcpeering.go
│   │   │   ├── vpcroute.go
│   │   │   └── vpntunnel.go
│   │   ├── security/        # Security service implementations
│   │   │   └── securitygroup.go
│   │   ├── storage/         # Storage service implementations
│   │   │   ├── blockstorage.go
│   │   │   └── snapshot.go
│   │   ├── database/        # Database service implementations
│   │   │   └── dbaas.go
│   │   └── schedule/        # Schedule service implementations
│   │       └── schedulejob.go
├── examples/                # Usage examples
├── go.mod
└── README.md
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
    defer resp.Body.Close()
    
    // Parse response
    var result schema.CloudServerListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d cloud servers\n", len(result.Items))
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
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 || resp.StatusCode == 201 {
        fmt.Println("Cloud server created successfully")
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
    defer resp.Body.Close()
    
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
    defer resp.Body.Close()
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
    defer resp.Body.Close()
    
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
    defer resp.Body.Close()
}
```

## Available Resources

### Compute
- **CloudServer** - Virtual machine management
- **KaaS** - Kubernetes cluster management

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
```

## Error Handling

All API methods return `*http.Response` and `error`. You should check both:

```go
resp, err := sdk.CloudServer.GetCloudServer(ctx, "my-project", "my-server")
if err != nil {
    log.Printf("Request failed: %v", err)
    return err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    log.Printf("API error: %d - %s", resp.StatusCode, string(body))
    return fmt.Errorf("unexpected status: %d", resp.StatusCode)
}

// Parse success response
var result schema.CloudServerResponse
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    return err
}
```

## API Interface Pattern

All resource clients follow the same pattern:

```go
type ResourceAPI interface {
    // List resources with pagination and filtering
    List{Resource}(ctx, project, params, ...editors) (*http.Response, error)
    
    // Get a single resource by name
    Get{Resource}(ctx, project, name, ...editors) (*http.Response, error)
    
    // Create or update a resource (upsert)
    CreateOrUpdate{Resource}(ctx, project, name, params, body, ...editors) (*http.Response, error)
    CreateOrUpdate{Resource}WithBody(ctx, project, name, params, contentType, body, ...editors) (*http.Response, error)
    
    // Delete a resource
    Delete{Resource}(ctx, project, name, params, ...editors) (*http.Response, error)
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

1. **Always close response bodies**:
   ```go
   resp, err := sdk.CloudServer.GetCloudServer(ctx, project, name)
   if err != nil {
       return err
   }
   defer resp.Body.Close()
   ```

2. **Check HTTP status codes**:
   ```go
   if resp.StatusCode != http.StatusOK {
       // Handle error
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