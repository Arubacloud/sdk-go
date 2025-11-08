# Storage Package

The `storage` package provides Go client interfaces for managing Aruba Cloud storage services, including block storage volumes and snapshots.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### BlockStorageAPI

Manage block storage volumes with full CRUD operations:
- List all block storage volumes in a project
- Get details of a specific volume
- Create a block storage volume
- Delete a block storage volume

### SnapshotAPI

Manage storage snapshots with full CRUD operations:
- List all snapshots in a project
- Get details of a specific snapshot
- Create a snapshot
- Delete a snapshot

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/storage"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interfaces
    var blockStorageAPI storage.BlockStorageAPI = storage.NewBlockStorageService(c)
    var snapshotAPI storage.SnapshotAPI = storage.NewSnapshotService(c)
}
```

### Block Storage Management

#### List Block Storage Volumes

```go
resp, err := blockStorageAPI.ListBlockStorageVolumes(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list block storage volumes: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d block storage volumes\n", len(resp.Data.Values))
    for _, volume := range resp.Data.Values {
        fmt.Printf("Volume: %s - Size: %dGB\n", volume.Metadata.Name, volume.Properties.SizeGb)
    }
}
```

### Snapshot Management

#### List Snapshots

```go
resp, err := snapshotAPI.ListSnapshots(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list snapshots: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d snapshots\n", len(resp.Data.Values))
    for _, snapshot := range resp.Data.Values {
        fmt.Printf("Snapshot: %s\n", snapshot.Metadata.Name)
    }
}
```
