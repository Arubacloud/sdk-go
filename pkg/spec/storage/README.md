# Storage Package

The `storage` package provides Go client interfaces for managing Aruba Cloud storage services, including block storage volumes and snapshots.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### StorageAPI

The unified `StorageAPI` interface provides all storage-related operations:

**Block Storage Operations** - Manage block storage volumes  
**Snapshot Operations** - Manage volume snapshots  

## Usage Examples

### Initialize the Service

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/storage"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create SDK client
    config := &client.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        HTTPClient:   &http.Client{Timeout: 30 * time.Second},
    }
    
    sdk, err := client.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Create unified storage service
    storageService := storage.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use storageService for all storage operations
}
```

### List BlockStorage

```go
// Use the unified service
storageService := storage.NewService(sdk)

resp, err := storageService.ListBlockStorageVolumes(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create BlockStorage

```go
volumeReq := schema.BlockStorageRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "my-volume",
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.BlockStoragePropertiesRequest{
        SizeGB: 20,
        Type:   schema.BlockStorageTypeStandard,
    },
}

createResp, err := storageService.CreateBlockStorageVolume(ctx, projectID, volumeReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
