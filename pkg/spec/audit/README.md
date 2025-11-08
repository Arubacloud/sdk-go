# Audit Package

The `audit` package provides Go client interfaces for accessing Aruba Cloud audit logs and event tracking.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### EventAPI

Retrieve audit events with read operations:
- List all audit events for a project
- Get details of a specific audit event

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/audit"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interface
    var eventAPI audit.EventAPI = audit.NewEventService(c)
}
```

### Event Management

#### List Audit Events

```go
resp, err := eventAPI.ListEvents(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list audit events: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d audit events\n", len(resp.Data.Values))
    for _, event := range resp.Data.Values {
        fmt.Printf("Event: %s - %s\n", event.Metadata.Name, event.Properties.EventType)
    }
}

// Access HTTP metadata
fmt.Printf("Status Code: %d\n", resp.StatusCode)
fmt.Printf("Content-Type: %s\n", resp.Headers.Get("Content-Type"))
```


