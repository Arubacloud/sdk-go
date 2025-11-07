# Database Package

The `database` package provides Go client interfaces for managing Aruba Cloud database services (DBaaS - Database as a Service).

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

### DBaaSAPI

Manage database instances with full CRUD operations:
- List all database instances in a project
- Get details of a specific database instance
- Create or update a database instance
- Delete a database instance

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/database"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interface
    var dbaasAPI database.DBaaSAPI = database.NewDBaaSService(c)
}
```

### Database Management

#### List Database Instances

```go
resp, err := dbaasAPI.ListDBaaS(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list databases: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode == 200 {
    fmt.Println("Databases retrieved successfully")
}
```