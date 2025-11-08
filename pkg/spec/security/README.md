# Security Package

The `security` package provides Go client interfaces for managing Aruba Cloud security services, including KMS (Key Management Service) for encryption key management.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### KMSAPI

Manage KMS encryption keys with full CRUD operations:
- List all KMS keys in a project
- Get details of a specific KMS key
- Create or update a KMS key
- Delete a KMS key

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/security"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interface
    var kmsAPI security.KMSAPI = security.NewKmsKeyService(c)
}
```

### KMS Key Management

#### List KMS Keys

```go
resp, err := kmsAPI.ListKMSKeys(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list KMS keys: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d KMS keys\n", len(resp.Data.Values))
    for _, key := range resp.Data.Values {
        fmt.Printf("Key: %s - Status: %s\n", key.Metadata.Name, key.Status.State)
    }
}
```

