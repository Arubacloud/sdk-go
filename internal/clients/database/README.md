# Database Package

The `database` package provides Go client interfaces for managing Aruba Cloud database services (DBaaS), including MySQL and PostgreSQL clusters, databases, users, grants, and backups.

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### DatabaseAPI

The unified `DatabaseAPI` interface provides all database-related operations:

**DBaaS Operations** - Manage database clusters (MySQL/PostgreSQL)  
**Database Operations** - Manage databases within clusters  
**User Operations** - Manage database users  
**Grant Operations** - Manage user permissions  
**Backup Operations** - Manage database backups  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/database"
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
    
    // Create unified database service
    databaseService := database.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use databaseService for all database operations
}
```

### List DBaaS

```go
// Use the unified service
databaseService := database.NewService(sdk)

resp, err := databaseService.ListDBaaS(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d items\n", len(resp.Data.Values))
}
```

### Create DBaaS

```go
dbaasReq := schema.DBaaSRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "my-database",
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.DBaaSPropertiesRequest{
        Engine:  "MySQL",
        Version: "8.0",
        // ... other properties
    },
}

createResp, err := databaseService.CreateDBaaS(ctx, projectID, dbaasReq, nil)
if err != nil {
    log.Fatalf("Failed to create: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created: %s\n", *createResp.Data.Metadata.Id)
}
```
