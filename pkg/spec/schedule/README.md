# Schedule Package

The `schedule` package provides Go client interfaces for managing Aruba Cloud scheduled jobs.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ScheduleAPI

The unified `ScheduleAPI` interface provides all schedule-related operations:

**Job Operations** - Manage scheduled jobs  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/schedule"
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
    
    // Create unified schedule service
    scheduleService := schedule.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use scheduleService for all schedule operations
}
```

### List Job

```go
// Use the unified service
scheduleService := schedule.NewService(sdk)

resp, err := scheduleService.ListJobs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create Job

```go
jobReq := schema.JobRequest{
    Metadata: schema.ResourceMetadataRequest{
        Name: "my-scheduled-job",
    },
    Properties: schema.JobPropertiesRequest{
        // ... job properties
    },
}

createResp, err := scheduleService.CreateJob(ctx, projectID, jobReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
