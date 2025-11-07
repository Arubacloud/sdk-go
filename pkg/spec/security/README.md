# Security Package

The `security` package provides Go client interfaces for managing Aruba Cloud security services, including KMS (Key Management Service) for encryption key management.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)

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
    var kmsAPI security.KMSAPI = security.NewKMSService(c)
}
```

### KMS Key Management

#### List KMS Keys

```go
resp, err := kmsAPI.ListKMSKeys(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list KMS keys: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode == 200 {
    fmt.Println("KMS keys retrieved successfully")
}
```

