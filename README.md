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
    "net/http"
    "time"
    
    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/project"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create SDK configuration with OAuth2 client credentials
    config := &client.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        HTTPClient:   &http.Client{Timeout: 30 * time.Second},
        Debug:        false, // Set to true for debug logging
    }
    
    // Initialize the SDK client (automatically obtains JWT token)
    sdk, err := client.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    // Use the SDK with context
    sdk = sdk.WithContext(ctx)
    
    // Initialize service clients
    projectAPI := project.NewProjectService(sdk)
    
    // Use the service
    createProject(ctx, projectAPI)
}
```

### Working with Projects

```go
func createProject(ctx context.Context, projectAPI *project.ProjectService) {
    // Create a new project
    projectReq := schema.ProjectRequest{
        Metadata: schema.ResourceMetadataRequest{
            Name: "my-production-project",
            Tags: []string{"production", "arubacloud-sdk"},
        },
        Properties: schema.ProjectPropertiesRequest{
            Description: stringPtr("My production project"),
            Default:     false,
        },
    }
    
    resp, err := projectAPI.CreateProject(ctx, projectReq, nil)
    if err != nil {
        log.Fatalf("Error creating project: %v", err)
    }
    
    if !resp.IsSuccess() {
        log.Fatalf("Failed to create project: status %d, error: %s", 
            resp.StatusCode, stringValue(resp.Error.Title))
    }
    
    projectID := *resp.Data.Metadata.Id
    fmt.Printf("✓ Created project: %s (ID: %s)\n", 
        *resp.Data.Metadata.Name, projectID)
    
    // Update the project
    projectReq.Properties.Description = stringPtr("Updated description")
    updateResp, err := projectAPI.UpdateProject(ctx, projectID, projectReq, nil)
    if err != nil {
        log.Fatalf("Error updating project: %v", err)
    }
    
    if updateResp.IsSuccess() {
        fmt.Printf("✓ Updated project: %s\n", *updateResp.Data.Metadata.Name)
    }
}

func stringPtr(s string) *string {
    return &s
}

func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
```

### Working with Network Resources

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
)

func createNetworkInfrastructure(ctx context.Context, sdk *client.Client, projectID string) {
    // Create Elastic IP
    elasticIPAPI := network.NewElasticIPService(sdk)
    
    elasticIPReq := schema.ElasticIpRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-elastic-ip",
                Tags: []string{"network", "public"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.ElasticIpPropertiesRequest{
            BillingPlan: schema.BillingPeriodResource{
                BillingPeriod: "Hour",
            },
        },
    }
    
    elasticIPResp, err := elasticIPAPI.CreateElasticIP(ctx, projectID, elasticIPReq, nil)
    if err != nil || !elasticIPResp.IsSuccess() {
        log.Fatalf("Failed to create Elastic IP: %v", err)
    }
    fmt.Printf("✓ Created Elastic IP: %s\n", *elasticIPResp.Data.Metadata.Name)
    
    // Create VPC
    vpcAPI := network.NewVPCService(sdk)
    
    vpcReq := schema.VpcRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-vpc",
                Tags: []string{"network", "infrastructure"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.VpcPropertiesRequest{
            Properties: &schema.VpcProperties{
                Default: boolPtr(false),
                Preset:  boolPtr(true),
            },
        },
    }
    
    vpcResp, err := vpcAPI.CreateVPC(ctx, projectID, vpcReq, nil)
    if err != nil || vpcResp.IsError() {
        log.Fatalf("Failed to create VPC: %v", err)
    }
    fmt.Printf("✓ Created VPC: %s\n", *vpcResp.Data.Metadata.Name)
    
    // Wait for VPC to become active
    vpcID := *vpcResp.Data.Metadata.Id
    waitForResourceActive(ctx, vpcAPI, projectID, vpcID, "VPC")
    
    // Create Subnet
    subnetAPI := network.NewSubnetService(sdk)
    
    subnetReq := schema.SubnetRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-subnet",
                Tags: []string{"network", "subnet"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.SubnetPropertiesRequest{
            Type:    schema.SubnetTypeAdvanced,
            Default: true,
            Network: &schema.SubnetNetwork{
                Address: "192.168.1.0/25",
            },
            DHCP: &schema.SubnetDHCP{
                Enabled: true,
            },
        },
    }
    
    subnetResp, err := subnetAPI.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)
    if err != nil || subnetResp.IsError() {
        log.Fatalf("Failed to create subnet: %v", err)
    }
    fmt.Printf("✓ Created Subnet: %s (Network: %s)\n",
        *subnetResp.Data.Metadata.Name,
        subnetResp.Data.Properties.Network.Address)
    
    // Create Security Group
    securityGroupAPI := network.NewSecurityGroupService(sdk)
    
    sgReq := schema.SecurityGroupRequest{
        Metadata: schema.ResourceMetadataRequest{
            Name: "my-security-group",
            Tags: []string{"security", "network"},
        },
        Properties: schema.SecurityGroupPropertiesRequest{
            Default: boolPtr(false),
        },
    }
    
    sgResp, err := securityGroupAPI.CreateSecurityGroup(ctx, projectID, vpcID, sgReq, nil)
    if err != nil || sgResp.IsError() {
        log.Fatalf("Failed to create security group: %v", err)
    }
    fmt.Printf("✓ Created Security Group: %s\n", *sgResp.Data.Metadata.Name)
    
    sgID := *sgResp.Data.Metadata.Id
    waitForResourceActive(ctx, securityGroupAPI, projectID, vpcID, sgID, "SecurityGroup")
    
    // Create Security Group Rule (allow SSH)
    securityRuleAPI := network.NewSecurityGroupRuleService(sdk)
    
    ruleReq := schema.SecurityRuleRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "allow-ssh",
                Tags: []string{"ssh-access", "ingress"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.SecurityRulePropertiesRequest{
            Direction: schema.RuleDirectionIngress,
            Protocol:  "TCP",
            Port:      "22",
            Target: &schema.RuleTarget{
                Kind:  schema.EndpointTypeIP,
                Value: "0.0.0.0/0",
            },
        },
    }
    
    ruleResp, err := securityRuleAPI.CreateSecurityGroupRule(ctx, projectID, vpcID, sgID, ruleReq, nil)
    if err != nil || ruleResp.IsError() {
        log.Fatalf("Failed to create security rule: %v", err)
    }
    fmt.Printf("✓ Created Security Rule: %s (Port: %s)\n",
        *ruleResp.Data.Metadata.Name,
        ruleResp.Data.Properties.Port)
}

func boolPtr(b bool) *bool {
    return &b
}

func waitForResourceActive(ctx context.Context, api interface{}, projectID, resourceID, resourceType string) {
    maxAttempts := 30
    pollInterval := 5 * time.Second
    
    fmt.Printf("⏳ Waiting for %s to become active...\n", resourceType)
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        time.Sleep(pollInterval)
        
        // Check resource status based on type
        // Implementation depends on resource type
        fmt.Printf("  %s state check (attempt %d/%d)\n", resourceType, attempt, maxAttempts)
        
        // Break when active
        if attempt == maxAttempts {
            log.Fatalf("Timeout waiting for %s to become active", resourceType)
        }
    }
    fmt.Printf("✓ %s is now active\n", resourceType)
}
```

### Working with Block Storage

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

func manageStorage(ctx context.Context, sdk *client.Client, projectID string) {
    storageAPI := storage.NewBlockStorageService(sdk)
    
    // Create block storage volume (bootable with Ubuntu image)
    blockStorageReq := schema.BlockStorageRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-block-storage",
                Tags: []string{"storage", "data"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.BlockStoragePropertiesRequest{
            SizeGB:        20,
            Type:          schema.BlockStorageTypeStandard,
            Zone:          "ITBG-1",
            BillingPeriod: "Hour",
            Bootable:      boolPtr(true),
            Image:         stringPtr("LU24-001"), // Ubuntu 24.04 LTS
        },
    }
    
    blockStorageResp, err := storageAPI.CreateBlockStorageVolume(ctx, projectID, blockStorageReq, nil)
    if err != nil || !blockStorageResp.IsSuccess() {
        log.Fatalf("Failed to create block storage: %v", err)
    }
    
    fmt.Printf("✓ Created block storage: %s (%d GB, %s)\n",
        *blockStorageResp.Data.Metadata.Name,
        blockStorageResp.Data.Properties.SizeGB,
        blockStorageResp.Data.Properties.Type)
    
    // Wait for block storage to become ready (Active or NotUsed state)
    blockStorageID := *blockStorageResp.Data.Metadata.Id
    maxAttempts := 30
    pollInterval := 5 * time.Second
    
    fmt.Println("⏳ Waiting for Block Storage to become active...")
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        time.Sleep(pollInterval)
        
        getBlockStorageResp, err := storageAPI.GetBlockStorageVolume(ctx, projectID, blockStorageID, nil)
        if err != nil {
            log.Printf("Error checking Block Storage status: %v", err)
            continue
        }
        
        if getBlockStorageResp.Data != nil && getBlockStorageResp.Data.Status.State != nil {
            state := *getBlockStorageResp.Data.Status.State
            fmt.Printf("  Block Storage state: %s (attempt %d/%d)\n", state, attempt, maxAttempts)
            
            // Block storage can be "Active" (attached) or "NotUsed" (unattached but ready)
            if state == "Active" || state == "NotUsed" {
                fmt.Printf("✓ Block Storage is now ready (state: %s)\n", state)
                break
            } else if state == "Failed" || state == "Error" {
                log.Fatalf("Block Storage creation failed with state: %s", state)
            }
        }
        
        if attempt == maxAttempts {
            log.Fatalf("Timeout waiting for Block Storage to become ready")
        }
    }
    
    // Create snapshot from block storage
    snapshotAPI := storage.NewSnapshotService(sdk)
    
    snapshotReq := schema.SnapshotRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-snapshot",
                Tags: []string{"backup", "snapshot"},
            },
            Location: schema.LocationRequest{
                Value: "ITBG-Bergamo",
            },
        },
        Properties: schema.SnapshotPropertiesRequest{
            BillingPeriod: stringPtr("Hour"),
            Volume: schema.ReferenceResource{
                Uri: *blockStorageResp.Data.Metadata.Uri,
            },
        },
    }
    
    snapshotResp, err := snapshotAPI.CreateSnapshot(ctx, projectID, snapshotReq, nil)
    if err != nil || !snapshotResp.IsSuccess() {
        log.Fatalf("Failed to create snapshot: %v", err)
    }
    
    fmt.Printf("✓ Created snapshot: %s from volume %s\n",
        *snapshotResp.Data.Metadata.Name,
        *blockStorageResp.Data.Metadata.Name)
}
```

### Working with Compute Resources

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
)

func createCloudServer(ctx context.Context, sdk *client.Client, projectID string, 
    vpcResp *schema.Response[schema.VpcResponse],
    elasticIPResp *schema.Response[schema.ElasticIpResponse],
    blockStorageResp *schema.Response[schema.BlockStorageResponse]) {
    
    // Create SSH Key Pair
    keyPairAPI := compute.NewKeyPairService(sdk)
    
    sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAAB... your-public-key"
    
    keyPairReq := schema.KeyPairRequest{
        Metadata: schema.ResourceMetadataRequest{
            Name: "my-ssh-keypair",
            Tags: []string{"compute", "access"},
        },
        Properties: schema.KeyPairPropertiesRequest{
            Value: sshPublicKey,
        },
    }
    
    keyPairResp, err := keyPairAPI.CreateKeyPair(ctx, projectID, keyPairReq, nil)
    if err != nil || !keyPairResp.IsSuccess() {
        log.Fatalf("Failed to create SSH key pair: %v", err)
    }
    fmt.Printf("✓ Created SSH Key Pair: %s\n", keyPairResp.Data.Metadata.Name)
    
    // Create Cloud Server using all resources
    cloudServerAPI := compute.NewCloudServerService(sdk)
    
    // Construct KeyPair URI
    keyPairUri := "/projects/" + projectID + "/providers/Aruba.Compute/keypairs/" + keyPairResp.Data.Metadata.Name
    
    cloudServerReq := schema.CloudServerRequest{
        Metadata: schema.ResourceMetadataRequest{
            Name: "my-cloud-server",
            Tags: []string{"compute", "production"},
        },
        Properties: schema.CloudServerPropertiesRequest{
            Zone: "ITBG-1",
            Vpc: schema.ReferenceResource{
                Uri: *vpcResp.Data.Metadata.Uri,
            },
            VpcPreset:  true,
            FlavorName: stringPtr("GP.2x2"), // 2 vCPU, 2GB RAM
            ElastcIp: schema.ReferenceResource{
                Uri: *elasticIPResp.Data.Metadata.Uri,
            },
            BootVolume: schema.ReferenceResource{
                Uri: *blockStorageResp.Data.Metadata.Uri,
            },
            KeyPair: schema.ReferenceResource{
                Uri: keyPairUri,
            },
        },
    }
    
    cloudServerResp, err := cloudServerAPI.CreateCloudServer(ctx, projectID, cloudServerReq, nil)
    if err != nil || !cloudServerResp.IsSuccess() {
        log.Fatalf("Failed to create cloud server: %v", err)
    }
    
    fmt.Printf("✓ Created Cloud Server: %s (Flavor: %s, Zone: %s)\n",
        cloudServerResp.Data.Metadata.Name,
        cloudServerResp.Data.Properties.Flavor.Name,
        cloudServerResp.Data.Properties.Zone)
}
```

## Available Resources

The SDK provides service clients for managing various Aruba Cloud resources:

### Compute
- **CloudServer** - Virtual machine management (create, update, delete, list)
- **KeyPair** - SSH key pair management

### Network
- **VPC** - Virtual Private Cloud management
- **Subnet** - Subnet management within VPCs
- **ElasticIP** - Elastic IP address allocation and management
- **SecurityGroup** - Security group management
- **SecurityGroupRule** - Security rules for controlling traffic
- **VPCPeering** - VPC peering connections
- **VPCPeeringRoute** - Routes for VPC peering
- **VPNTunnel** - VPN tunnel management
- **VPNRoute** - Routes for VPN tunnels
- **LoadBalancer** - Load balancer management

### Storage
- **BlockStorage** - Block storage volume management
- **Snapshot** - Volume snapshot management

### Database
- **DBaaS** - Database as a Service management
- **Database** - Individual database management
- **DatabaseUser** - Database user management
- **DatabaseGrant** - Database grant/permission management
- **Backup** - Database backup management

### Project
- **Project** - Project management and configuration

### Security
- **KMS** - Key Management Service

### Monitoring
- **Metric** - Metrics and monitoring data
- **Alert** - Alert configuration and management

### Automation
- **ScheduleJob** - Scheduled task management

### Audit
- **Event** - Audit event tracking and retrieval

## Authentication

The SDK uses **OAuth2 Client Credentials Flow** to obtain JWT Bearer tokens automatically.

### Automatic Token Management

```go
config := &client.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   &http.Client{Timeout: 30 * time.Second},
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
- ✅ **Debug logging** support for troubleshooting

### Advanced Configuration

```go
config := &client.Config{
    ClientID:           "your-client-id",
    ClientSecret:       "your-client-secret",
    HTTPClient:         &http.Client{Timeout: 30 * time.Second},
    TokenRefreshBuffer: 10 * time.Minute, // Refresh 10 min before expiry
    Debug:              true,              // Enable debug logging
}

sdk, err := client.NewClient(config)
```

## Request Customization

### Using Request Editors

Request editors allow you to modify HTTP requests before they are sent. This is useful for adding custom headers, implementing middleware, or debugging.

```go
// Add custom headers or modify request
customEditor := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Request-ID", "123456")
    req.Header.Set("X-Application", "my-app")
    return nil
}

// Pass request editor to any API method (last parameter)
resp, err := projectAPI.GetProject(
    ctx,
    "my-project",
    customEditor, // request editor
)

if resp.IsSuccess() {
    fmt.Printf("Project: %s\n", *resp.Data.Metadata.Name)
}
```

### Service-Level Initialization

Each resource type has its own service client that must be initialized separately:

```go
// Initialize service clients as needed
projectAPI := project.NewProjectService(sdk)
vpcAPI := network.NewVPCService(sdk)
storageAPI := storage.NewBlockStorageService(sdk)
computeAPI := compute.NewCloudServerService(sdk)

// Use the service clients
projectResp, err := projectAPI.GetProject(ctx, "my-project", nil)
vpcResp, err := vpcAPI.GetVPC(ctx, "my-project", "my-vpc", nil)
```

## Debug Logging

The SDK supports comprehensive debug logging for troubleshooting API interactions. When enabled, it logs:

- **Service Operations**: High-level operations with resource identifiers
- **HTTP Requests**: Method, URL, headers (with token redaction), query parameters, and request body
- **HTTP Responses**: Status code, headers, and response body

### Enable Debug Logging

```go
config := &client.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Debug:        true,  // Enable debug logging to stdout/stderr
}

sdk, err := client.NewClient(config)
```

### Custom Logger

You can provide a custom logger implementing the `Logger` interface:

```go
type Logger interface {
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}

// Use custom logger
config := &client.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Logger:       myCustomLogger,  // Your logger implementation
}
```

### Example Debug Output

```
[DEBUG] Initializing SDK client
[DEBUG] Successfully obtained initial token
[DEBUG] Listing cloud servers for project: my-project
[DEBUG] Making GET request to https://api.arubacloud.com/v1/projects/my-project/cloudservers
[DEBUG] Added query parameters: map[limit:10]
[DEBUG] Request headers (pre-auth): map[]
[DEBUG] Request headers (final): map[Authorization:Bearer [REDACTED] Content-Type:application/json X-Application:aruba-sdk-example]
[DEBUG] Received response with status: 200 OK
[DEBUG] Response headers: map[Content-Type:[application/json] X-Request-Id:[abc-123]]
[DEBUG] Response body: {"values":[...],"metadata":{...}}
```

## Error Handling

All API methods return a generic `Response[T]` wrapper and `error`. The response includes both parsed data and HTTP metadata.

### Response[T] Structure

The `Response[T]` type provides:
- `Data *T` - Parsed response data (for 2xx responses)
- `Error *ErrorResponse` - RFC 7807 error details (for 4xx/5xx responses)
- `HTTPResponse *http.Response` - Full HTTP response
- `StatusCode int` - HTTP status code
- `Headers http.Header` - Response headers
- `RawBody []byte` - Raw response body
- `IsSuccess() bool` - Returns true for 2xx status codes
- `IsError() bool` - Returns true for 4xx/5xx status codes

### Basic Error Handling

```go
resp, err := projectAPI.GetProject(ctx, "my-project", nil)
if err != nil {
    log.Printf("Request failed: %v", err)
    return err
}

// Check status using helper methods
if resp.IsError() {
    log.Printf("API error: %d - %s: %s", 
        resp.StatusCode,
        stringValue(resp.Error.Title),
        stringValue(resp.Error.Detail))
    return fmt.Errorf("API error: %d", resp.StatusCode)
}

if resp.IsSuccess() {
    // Access parsed data directly
    fmt.Printf("Project: %s (ID: %s)\n", 
        *resp.Data.Metadata.Name,
        *resp.Data.Metadata.Id)
}
```

### RFC 7807 Error Response

For 4xx and 5xx responses, the `Error` field contains RFC 7807-compliant error details:

```go
resp, err := vpcAPI.CreateVPC(ctx, projectID, vpcReq, nil)
if err != nil {
    log.Fatalf("Request failed: %v", err)
}

if resp.IsError() && resp.Error != nil {
    fmt.Printf("Error creating VPC:\n")
    fmt.Printf("  Status: %d\n", resp.StatusCode)
    fmt.Printf("  Title: %s\n", stringValue(resp.Error.Title))
    fmt.Printf("  Detail: %s\n", stringValue(resp.Error.Detail))
    fmt.Printf("  Type: %s\n", stringValue(resp.Error.Type))
    if resp.Error.Instance != nil {
        fmt.Printf("  Instance: %s\n", *resp.Error.Instance)
    }
}
```

### Accessing HTTP Metadata

```go
resp, err := storageAPI.GetBlockStorageVolume(ctx, projectID, volumeID, nil)
if err != nil {
    return err
}

if resp.IsSuccess() {
    // Access parsed data
    fmt.Printf("Volume: %s (Size: %d GB)\n", 
        *resp.Data.Metadata.Name,
        resp.Data.Properties.SizeGB)
    
    // Access HTTP metadata if needed
    fmt.Printf("Status Code: %d\n", resp.StatusCode)
    fmt.Printf("Content-Type: %s\n", resp.Headers.Get("Content-Type"))
    fmt.Printf("X-Request-Id: %s\n", resp.Headers.Get("X-Request-Id"))
}
```

## API Interface Pattern

All resource services follow a consistent pattern with these common operations:

### Service Initialization

```go
// Each resource type has its own service client
import (
    "github.com/Arubacloud/sdk-go/pkg/spec/project"
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
    "github.com/Arubacloud/sdk-go/pkg/spec/storage"
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
)

// Initialize services
projectAPI := project.NewProjectService(sdk)
vpcAPI := network.NewVPCService(sdk)
storageAPI := storage.NewBlockStorageService(sdk)
computeAPI := compute.NewCloudServerService(sdk)
```

### Common Operations

Most resource services provide these operations:

```go
// Create a resource
Create{Resource}(ctx, projectID, request, ...editors) (*schema.Response[schema.{Resource}Response], error)

// Get a single resource by ID
Get{Resource}(ctx, projectID, resourceID, ...editors) (*schema.Response[schema.{Resource}Response], error)

// Update a resource
Update{Resource}(ctx, projectID, resourceID, request, ...editors) (*schema.Response[schema.{Resource}Response], error)

// Delete a resource
Delete{Resource}(ctx, projectID, resourceID, ...editors) (*schema.Response[any], error)

// List resources (if applicable)
List{Resources}(ctx, projectID, params, ...editors) (*schema.Response[schema.{Resource}ListResponse], error)
```

### Hierarchical Resources

Some resources are nested within others (e.g., Subnets within VPCs):

```go
// Create a subnet within a VPC
subnetAPI := network.NewSubnetService(sdk)
resp, err := subnetAPI.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)

// Create a security rule within a security group
ruleAPI := network.NewSecurityGroupRuleService(sdk)
resp, err := ruleAPI.CreateSecurityGroupRule(ctx, projectID, vpcID, securityGroupID, ruleReq, nil)
```

### Request Editors (Optional Last Parameter)

All methods accept optional request editors as variadic parameters:

```go
// Without request editor
resp, err := projectAPI.GetProject(ctx, projectID, nil)

// With request editor
customEditor := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Request-ID", "12345")
    return nil
}
resp, err := projectAPI.GetProject(ctx, projectID, customEditor)
```

## Context Support

All operations support context for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := projectAPI.GetProject(ctx, "my-project", nil)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel after 5 seconds
}()

resp, err := vpcAPI.GetVPC(ctx, projectID, vpcID, nil)

// Context propagation with SDK
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

// Apply context to SDK client - all operations will use this context
sdk = sdk.WithContext(ctx)
```

## Resource State Polling

Many resources require time to become active after creation. Use polling to wait for resources to reach the desired state:

```go
func waitForResourceActive(ctx context.Context, getFunc func() (*schema.Response[ResourceResponse], error), 
    resourceType string) error {
    
    maxAttempts := 30
    pollInterval := 5 * time.Second
    
    fmt.Printf("⏳ Waiting for %s to become active...\n", resourceType)
    
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        time.Sleep(pollInterval)
        
        resp, err := getFunc()
        if err != nil {
            log.Printf("Error checking %s status: %v", resourceType, err)
            continue
        }
        
        if resp.Data != nil && resp.Data.Status.State != nil {
            state := *resp.Data.Status.State
            fmt.Printf("  %s state: %s (attempt %d/%d)\n", resourceType, state, attempt, maxAttempts)
            
            switch state {
            case "Active":
                fmt.Printf("✓ %s is now active\n", resourceType)
                return nil
            case "Failed", "Error":
                return fmt.Errorf("%s creation failed with state: %s", resourceType, state)
            }
        }
        
        if attempt == maxAttempts {
            return fmt.Errorf("timeout waiting for %s to become active", resourceType)
        }
    }
    
    return nil
}

// Example usage: Wait for VPC to become active
err := waitForResourceActive(ctx, func() (*schema.Response[schema.VpcResponse], error) {
    return vpcAPI.GetVPC(ctx, projectID, vpcID, nil)
}, "VPC")
if err != nil {
    log.Fatalf("Failed to wait for VPC: %v", err)
}
```

## Complete Example

See [cmd/example/main.go](cmd/example/main.go) for a comprehensive example that demonstrates:

- Project creation and management
- Network infrastructure setup (VPC, Subnet, Security Groups, Elastic IP)
- Storage management (Block Storage, Snapshots)
- Compute resources (SSH Key Pairs, Cloud Servers)
- Resource state polling
- Proper error handling with `Response[T]`
- Context management with timeouts

The example creates a complete infrastructure stack with proper resource dependencies and state management.

## Best Practices

1. **Always check response status**:
   ```go
   resp, err := projectAPI.GetProject(ctx, projectID, nil)
   if err != nil {
       return err
   }
   
   if !resp.IsSuccess() {
       return fmt.Errorf("API error: %d - %s", 
           resp.StatusCode, stringValue(resp.Error.Title))
   }
   
   // Now safe to use resp.Data
   fmt.Println(*resp.Data.Metadata.Name)
   ```

2. **Use contexts with timeouts for long-running operations**:
   ```go
   // For operations that may take time (create, update, delete)
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
   defer cancel()
   
   sdk = sdk.WithContext(ctx)
   ```

3. **Poll for resource state after creation**:
   ```go
   // Many resources need time to become active
   // Use polling with reasonable intervals (e.g., 5 seconds)
   // Set appropriate max attempts (e.g., 30 attempts = 150 seconds)
   for attempt := 1; attempt <= 30; attempt++ {
       time.Sleep(5 * time.Second)
       resp, err := api.GetResource(ctx, projectID, resourceID, nil)
       if err == nil && resp.Data.Status.State != nil {
           if *resp.Data.Status.State == "Active" {
               break
           }
       }
   }
   ```

4. **Initialize service clients once and reuse**:
   ```go
   // Create SDK client once
   sdk, err := client.NewClient(config)
   
   // Initialize service clients once
   projectAPI := project.NewProjectService(sdk)
   vpcAPI := network.NewVPCService(sdk)
   
   // Reuse for all operations
   resp1, err := projectAPI.GetProject(ctx, "project1", nil)
   resp2, err := projectAPI.GetProject(ctx, "project2", nil)
   ```

5. **Handle errors from RFC 7807 error responses**:
   ```go
   if resp.IsError() && resp.Error != nil {
       log.Printf("API Error - Title: %s, Detail: %s, Type: %s",
           stringValue(resp.Error.Title),
           stringValue(resp.Error.Detail),
           stringValue(resp.Error.Type))
   }
   ```

6. **Use helper functions for pointer types**:
   ```go
   func stringPtr(s string) *string { return &s }
   func boolPtr(b bool) *bool { return &b }
   func intPtr(i int) *int { return &i }
   
   // Usage
   req := schema.ProjectRequest{
       Properties: schema.ProjectPropertiesRequest{
           Description: stringPtr("My project"),
           Default:     false,
       },
   }
   ```

7. **Enable debug logging during development**:
   ```go
   config := &client.Config{
       ClientID:     "your-client-id",
       ClientSecret: "your-client-secret",
       Debug:        true, // Enable to see HTTP requests/responses
   }
   ```

8. **Respect resource dependencies**:
   ```go
   // Create resources in the correct order:
   // 1. Project
   // 2. VPC (wait for Active)
   // 3. Subnet (requires VPC to be Active)
   // 4. Security Group (wait for Active)
   // 5. Cloud Server (requires VPC, Subnet, Security Group)
   ```

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