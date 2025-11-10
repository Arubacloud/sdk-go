# Container Package

The `container` package provides Go client interfaces for managing Aruba Cloud Kubernetes as a Service (KaaS) clusters.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ContainerAPI

The unified `ContainerAPI` interface provides all container-related operations:

**KaaS Operations** - Manage Kubernetes clusters  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/container"
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
    
    // Create unified container service
    containerService := container.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use containerService for all container operations
}
```

### List KaaS

```go
// Use the unified service
containerService := container.NewService(sdk)

resp, err := containerService.ListKaaS(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create KaaS

```go
kaasReq := schema.KaaSRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "my-k8s-cluster",
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.KaaSPropertiesRequest{
        K8sVersion: "1.28",
        // ... other properties
    },
}

createResp, err := containerService.CreateKaaS(ctx, projectID, kaasReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
