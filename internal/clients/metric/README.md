# Metric Package

The `metric` package provides Go client interfaces for managing Aruba Cloud monitoring services, including metrics and alerts.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### MetricAPI

The unified `MetricAPI` interface provides all metric-related operations:

**Metric Operations** - View metrics  
**Alert Operations** - View alerts  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/metric"
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
    
    // Create unified metric service
    metricService := metric.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use metricService for all metric operations
}
```

### List Metrics

```go
// Use the unified service
metricService := metric.NewService(sdk)

resp, err := metricService.ListMetrics(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```
