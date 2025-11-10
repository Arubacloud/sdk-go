# Audit Package

The `audit` package provides Go client interfaces for managing Aruba Cloud audit services, including audit event tracking.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### AuditAPI

The unified `AuditAPI` interface provides all audit-related operations:

**Event Operations** - View audit events  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/audit"
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
    
    // Create unified audit service
    auditService := audit.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use auditService for all audit operations
}
```

### List Audit Events

```go
// Use the unified service
auditService := audit.NewService(sdk)

resp, err := auditService.ListEvents(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```
