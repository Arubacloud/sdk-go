# Schedule Package

The `schedule` package provides Go client interfaces for managing Aruba Cloud scheduled jobs and automation tasks.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ScheduleJobAPI

Manage scheduled jobs with full CRUD operations:
- List all scheduled jobs in a project
- Get details of a specific job
- Create or update a job
- Delete a job

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/schedule"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
        // Initialize API interface
    var jobAPI schedule.ScheduleJobAPI = schedule.NewJobService(c)
}
```

### Job Management

#### List Scheduled Jobs

```go
resp, err := jobAPI.ListScheduleJobs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list jobs: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d scheduled jobs\n", len(resp.Data.Values))
    for _, job := range resp.Data.Values {
        fmt.Printf("Job: %s - Type: %s\n", job.Metadata.Name, job.Properties.JobType)
    }
}
```

}
```

### Job Management

#### List Scheduled Jobs

```go
resp, err := jobAPI.ListJobs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list jobs: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode == 200 {
    fmt.Println("Jobs retrieved successfully")
}
```