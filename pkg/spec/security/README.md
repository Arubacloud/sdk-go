# Security Package

The `security` package provides Go client interfaces for managing Aruba Cloud security services, including Key Management Service (KMS).

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### SecurityAPI

The unified `SecurityAPI` interface provides all security-related operations:

**KMS Operations** - Manage encryption keys  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/security"
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
    
    // Create unified security service
    securityService := security.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use securityService for all security operations
}
```

### List KMS Key

```go
// Use the unified service
securityService := security.NewService(sdk)

resp, err := securityService.ListKMSKeys(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create KMS Key

```go
kmsReq := schema.KMSRequest{
    Metadata: schema.ResourceMetadataRequest{
        Name: "my-encryption-key",
    },
    Properties: schema.KMSPropertiesRequest{
        // ... key properties
    },
}

createResp, err := securityService.CreateKMSKey(ctx, projectID, kmsReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
