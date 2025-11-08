# Metric Package

The `metric` package provides Go client interfaces for managing Aruba Cloud monitoring and alerting services, including metrics and alerts.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### MetricAPI

Retrieve monitoring metrics with read operations:
- List all metrics for a project
- Get details of a specific metric

### AlertAPI

Retrieve monitoring alerts with read operations:
- List all alerts for a project
- Get details of a specific alert

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/metric"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interfaces
    var metricAPI metric.MetricAPI = metric.NewMetricService(c)
    var alertAPI metric.AlertAPI = metric.NewAlertService(c)
}
```

### Metric Management

#### List Metrics

```go
resp, err := metricAPI.ListMetrics(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list metrics: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d metrics\n", len(resp.Data.Values))
    for _, metric := range resp.Data.Values {
        fmt.Printf("Metric: %s\n", metric.Metadata.Name)
    }
}
```

### Alert Management

#### List Alerts

```go
resp, err := alertAPI.ListAlerts(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list alerts: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d alerts\n", len(resp.Data.Values))
    for _, alert := range resp.Data.Values {
        fmt.Printf("Alert: %s - State: %s\n", alert.Metadata.Name, alert.Properties.State)
    }
}
```
