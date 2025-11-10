# SDK Client with Direct Service Access

The Aruba Cloud Go SDK now provides direct access to all services from the main client!

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"
    
    sdkgo "github.com/Arubacloud/sdk-go"
    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Initialize the SDK client
    config := &client.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        HTTPClient:   &http.Client{Timeout: 30 * time.Second},
        Debug:        true,
    }
    
    sdk, err := sdkgo.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    
    ctx := context.Background()
    projectID := "my-project"
    
    // Use services directly from the SDK client!
    
    // Compute operations
    servers, err := sdk.Compute.ListCloudServers(ctx, projectID, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Network operations
    vpcs, err := sdk.Network.ListVPCs(ctx, projectID, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Storage operations
    volumes, err := sdk.Storage.ListBlockStorageVolumes(ctx, projectID, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Database operations
    databases, err := sdk.Database.ListDBaaS(ctx, projectID, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Project operations
    projects, err := sdk.Project.ListProjects(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Available Services

Access all services directly from the SDK client:

- **`sdk.Compute`** - CloudServers, KeyPairs
- **`sdk.Network`** - VPCs, Subnets, SecurityGroups, ElasticIPs, LoadBalancers, VPN, Peering
- **`sdk.Storage`** - BlockStorage, Snapshots
- **`sdk.Database`** - DBaaS, Databases, Users, Grants, Backups
- **`sdk.Container`** - KaaS (Kubernetes)
- **`sdk.Security`** - KMS Keys
- **`sdk.Metric`** - Metrics, Alerts
- **`sdk.Audit`** - Events
- **`sdk.Schedule`** - Jobs
- **`sdk.Project`** - Projects

## Example: Create Complete Infrastructure

```go
func createInfrastructure(sdk *sdkgo.Client, ctx context.Context, projectID string) {
    // Create VPC
    vpcReq := schema.VpcRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-vpc",
            },
            Location: schema.LocationRequest{Value: "ITBG-Bergamo"},
        },
        Properties: schema.VpcPropertiesRequest{
            Properties: &schema.VpcProperties{
                Default: boolPtr(false),
                Preset:  boolPtr(true),
            },
        },
    }
    
    vpcResp, err := sdk.Network.CreateVPC(ctx, projectID, vpcReq, nil)
    if err != nil || !vpcResp.IsSuccess() {
        log.Fatalf("Failed to create VPC: %v", err)
    }
    
    // Create Subnet
    vpcID := *vpcResp.Data.Metadata.Id
    subnetReq := schema.SubnetRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-subnet",
            },
            Location: schema.LocationRequest{Value: "ITBG-Bergamo"},
        },
        Properties: schema.SubnetPropertiesRequest{
            Type:    schema.SubnetTypeAdvanced,
            Default: true,
            Network: &schema.SubnetNetwork{Address: "192.168.1.0/25"},
            DHCP:    &schema.SubnetDHCP{Enabled: true},
        },
    }
    
    subnetResp, err := sdk.Network.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)
    if err != nil || !subnetResp.IsSuccess() {
        log.Fatalf("Failed to create subnet: %v", err)
    }
    
    // Create Block Storage
    storageReq := schema.BlockStorageRequest{
        Metadata: schema.RegionalResourceMetadataRequest{
            ResourceMetadataRequest: schema.ResourceMetadataRequest{
                Name: "my-volume",
            },
            Location: schema.LocationRequest{Value: "ITBG-Bergamo"},
        },
        Properties: schema.BlockStoragePropertiesRequest{
            SizeGB:        20,
            Type:          schema.BlockStorageTypeStandard,
            Zone:          "ITBG-1",
            BillingPeriod: "Hour",
            Bootable:      boolPtr(true),
            Image:         stringPtr("LU24-001"), // Ubuntu 24.04
        },
    }
    
    storageResp, err := sdk.Storage.CreateBlockStorageVolume(ctx, projectID, storageReq, nil)
    if err != nil || !storageResp.IsSuccess() {
        log.Fatalf("Failed to create storage: %v", err)
    }
}

func boolPtr(b bool) *bool { return &b }
func stringPtr(s string) *string { return &s }
```

## Benefits

✅ **Simpler API** - No need to create service instances manually  
✅ **Direct Access** - All services available immediately: `sdk.Network`, `sdk.Compute`, etc.  
✅ **Type-Safe** - Full IntelliSense/autocomplete support  
✅ **Consistent** - Same pattern across all services  
✅ **Clean Code** - Less boilerplate, more readable  

## Migration from Old Pattern

### Before (manual service creation):
```go
import (
    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
)

sdk, _ := client.NewClient(config)

// Had to create services manually
computeService := compute.NewService(sdk)
networkService := network.NewService(sdk)

servers, _ := computeService.ListCloudServers(ctx, projectID, nil)
vpcs, _ := networkService.ListVPCs(ctx, projectID, nil)
```

### After (direct service access):
```go
import (
    sdkgo "github.com/Arubacloud/sdk-go"
    "github.com/Arubacloud/sdk-go/pkg/client"
)

sdk, _ := sdkgo.NewClient(config)

// Services available directly!
servers, _ := sdk.Compute.ListCloudServers(ctx, projectID, nil)
vpcs, _ := sdk.Network.ListVPCs(ctx, projectID, nil)
```

Much cleaner! 🎉
