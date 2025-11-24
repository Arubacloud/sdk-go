# Project Package

The `project` package provides Go client interfaces for managing Aruba Cloud projects.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ProjectAPI

The unified `ProjectAPI` interface provides all project-related operations:

**Project Operations** - Manage projects  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/project"
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
    
    // Create unified project service
    projectService := project.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use projectService for all project operations
}
```

### List Project

```go
// Use the unified service
projectService := project.NewService(sdk)

resp, err := projectService.ListProjects(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create Project

```go
projectReq := schema.ProjectRequest{
    Metadata: schema.ResourceMetadataRequest{
        Name: "my-project",
        Tags: []string{"production"},
    },
    Properties: schema.ProjectPropertiesRequest{
        Description: stringPtr("My project"),
    },
}

createResp, err := projectService.CreateProject(ctx, projectReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
