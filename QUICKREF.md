# Aruba Cloud Go SDK - Quick Reference

## Overview
A Go SDK for Aruba Cloud REST APIs with automatic OAuth2 JWT token management and comprehensive resource management.

## Key Features
✅ OAuth2 Client Credentials Flow with automatic token refresh
✅ Type-safe API clients for all Aruba Cloud services
✅ Thread-safe token management
✅ Context support for request cancellation
✅ Comprehensive error handling with RFC 7807 Problem Details
✅ Resource state polling and lifecycle management
✅ Modular architecture with 11 service packages

## Quick Start

### 1. Install SDK
```bash
go get github.com/Arubacloud/sdk-go
```

### 2. Initialize Client
```go
import (
    "context"
    "net/http"
    "time"
    
    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/project"
)

config := &client.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   &http.Client{Timeout: 30 * time.Second},
    Debug:        false,
}

sdk, err := client.NewClient(config)
if err != nil {
    log.Fatal(err)
}

ctx := context.WithTimeout(context.Background(), 5*time.Minute)
sdk = sdk.WithContext(ctx)
```

### 3. Use Service APIs
```go
// Project Management
projectAPI := project.NewProjectService(sdk)
resp, err := projectAPI.ListProjects(ctx, nil)

// Network Management
import "github.com/Arubacloud/sdk-go/pkg/spec/network"
vpcAPI := network.NewVPCService(sdk)
vpcResp, err := vpcAPI.CreateVPC(ctx, projectID, vpcRequest, nil)

// Database Management
import "github.com/Arubacloud/sdk-go/pkg/spec/database"
dbaasAPI := database.NewDBaaSService(sdk)
dbResp, err := dbaasAPI.CreateDBaaS(ctx, projectID, dbRequest, nil)

// Container Management (Kubernetes)
import "github.com/Arubacloud/sdk-go/pkg/spec/container"
kaasAPI := container.NewKaaSService(sdk)
clusterResp, err := kaasAPI.CreateKaaS(ctx, projectID, kaasRequest, nil)
```

## File Structure
```
sdk-go/
├── cmd/
│   └── example/               # Complete working examples
│       ├── main.go            # 882 lines - Modular example with 12 resources
│       └── README.md
├── pkg/
│   ├── client/                # SDK Core
│   │   ├── client.go          # Main client with OAuth2
│   │   ├── token.go           # Thread-safe token manager
│   │   ├── error.go           # RFC 7807 error handling
│   │   ├── middleware.go      # Request helpers
│   │   ├── params.go          # Query parameter helpers
│   │   └── providers.go       # (deprecated - use spec packages)
│   └── spec/                  # Service API Implementations
│       ├── audit/             # Audit event tracking
│       ├── compute/           # Cloud servers, keypairs
│       ├── container/         # Kubernetes (KaaS)
│       ├── database/          # MySQL/PostgreSQL (DBaaS)
│       ├── metric/            # Monitoring and alerts
│       ├── network/           # VPC, subnets, security groups, load balancers
│       ├── project/           # Project management
│       ├── schedule/          # Scheduled jobs
│       ├── schema/            # Shared data types
│       ├── security/          # KMS key management
│       └── storage/           # Block storage, snapshots
├── tools/                     # Development tools
├── go.mod
├── Makefile
├── README.md
├── QUICKREF.md                # This file
├── SDK_READY.md               # SDK completion checklist
├── DEVELOPMENT.md
├── FILTERS.md
└── OAUTH2.md
```

## Available Services

| Service | Package | Description |
|---------|---------|-------------|
| **Project** | `pkg/spec/project` | Project creation, update, delete, listing |
| **Network** | `pkg/spec/network` | VPC, Subnet, Security Groups, Elastic IPs, Load Balancers, VPN |
| **Compute** | `pkg/spec/compute` | Cloud Servers, SSH Key Pairs |
| **Storage** | `pkg/spec/storage` | Block Storage volumes, Snapshots |
| **Database** | `pkg/spec/database` | DBaaS (MySQL/PostgreSQL), Users, Grants, Backups |
| **Container** | `pkg/spec/container` | KaaS (Kubernetes clusters) |
| **Security** | `pkg/spec/security` | KMS (Key Management Service) |
| **Metric** | `pkg/spec/metric` | Metrics, Alerts, Monitoring |
| **Audit** | `pkg/spec/audit` | Audit events, Compliance tracking |
| **Schedule** | `pkg/spec/schedule` | Scheduled jobs, Cron management |

## Response Handling

All API methods return `schema.Response[T]` with built-in helpers:

```go
resp, err := api.GetResource(ctx, projectID, resourceID, nil)
if err != nil {
    return err
}

// Check HTTP status
if resp.IsSuccess() {       // 200-299
    fmt.Println(resp.Data)
}

if resp.IsError() {         // 400+
    fmt.Printf("Error: %s - %s\n", 
        stringValue(resp.Error.Title),
        stringValue(resp.Error.Detail))
}

// Access response data
if resp.Data != nil {
    resourceID := *resp.Data.Metadata.Id
    resourceName := *resp.Data.Metadata.Name
}
```
## Authentication Flow
1. Client initialization → Request token from OAuth2 endpoint
2. Token stored with expiry time
3. Before each API call → Check if token valid
4. If expired/expiring → Automatically refresh
5. Add Bearer token to request headers

## Resource State Polling

Many resources require polling until they reach "Active" state:

```go
// Example: Wait for VPC to become active
maxAttempts := 30
pollInterval := 5 * time.Second

for attempt := 1; attempt <= maxAttempts; attempt++ {
    time.Sleep(pollInterval)
    
    getResp, err := vpcAPI.GetVPC(ctx, projectID, vpcID, nil)
    if err != nil {
        continue
    }
    
    if getResp.Data.Status.State != nil {
        state := *getResp.Data.Status.State
        fmt.Printf("VPC state: %s (attempt %d/%d)\n", state, attempt, maxAttempts)
        
        if state == "Active" {
            fmt.Println("✓ VPC is ready")
            break
        }
    }
}
```

Resources that typically require polling:
- VPC (5-30 seconds)
- Security Groups (5-30 seconds)
- Block Storage (30-60 seconds)
- DBaaS clusters (2-5 minutes)
- KaaS clusters (3-10 minutes)

## Common Commands
```bash
make build       # Build the project
make test        # Run tests
make lint        # Run linters
make fmt         # Format code
make clean       # Clean build artifacts
```

## Configuration Options
```go
type Config struct {
    ClientID           string              // OAuth2 client ID (required)
    ClientSecret       string              // OAuth2 client secret (required)
    HTTPClient         *http.Client        // HTTP client (optional, default: 30s timeout)
    Debug              bool                // Enable debug logging (optional)
    TokenRefreshBuffer time.Duration       // Token refresh timing (default: 5min)
}
```

## Helper Functions

The SDK uses pointer types for optional fields. Use these helpers:

```go
func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool { return &b }
func int32Ptr(i int32) *int32 { return &i }

func stringValue(s *string) string {
    if s == nil { return "" }
    return *s
}

func int32Value(i *int32) int32 {
    if i == nil { return 0 }
    return *i
}

// Usage
req := schema.VpcRequest{
    Properties: schema.VpcPropertiesRequest{
        Properties: &schema.VpcProperties{
            Default: boolPtr(false),
            Preset:  boolPtr(true),
        },
    },
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
resp, err := api.GetResource(ctx, projectID, resourceID, nil)
if err != nil {
    return fmt.Errorf("request failed: %w", err)
}

// Check response status
if !resp.IsSuccess() {
    return fmt.Errorf("API error: %d - %s: %s",
        resp.StatusCode,
        stringValue(resp.Error.Title),
        stringValue(resp.Error.Detail))
}

// Safe to use response data
resourceName := *resp.Data.Metadata.Name
```

## Complete Example

See [cmd/example/main.go](cmd/example/main.go) for a modular example (882 lines) that demonstrates:

**Architecture:**
- `main()` - Simple 28-line orchestration function
- `ResourceCollection` - Type-safe struct holding 12 resource types
- `createAllResources()` - Orchestrates creation in dependency order
- 11 individual functions: `createProject()`, `createElasticIP()`, `createBlockStorage()`, etc.
- `printResourceSummary()` - Clean output formatting

**Resources Created:**
1. Project
2. Elastic IP
3. Block Storage (20GB, Ubuntu 24.04)
4. Snapshot
5. VPC (with polling)
6. Subnet (192.168.1.0/25)
7. Security Group (with polling)
8. Security Rule (SSH port 22)
9. SSH Key Pair
10. DBaaS (MySQL 8.0, autoscaling)
11. KaaS (Kubernetes 1.28, 3 nodes, HA)
12. Cloud Server (commented)

**Run the example:**
```bash
cd cmd/example
go run main.go
```

## Support
- **README.md** - Complete documentation with examples
- **DEVELOPMENT.md** - Development and contribution guide
- **FILTERS.md** - Query filtering documentation
- **OAUTH2.md** - OAuth2 implementation details
- **SDK_READY.md** - SDK feature checklist
- **cmd/example/** - Complete working examples

## Version
Go 1.22+
