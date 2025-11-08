# Database Package

The `database` package provides Go client interfaces for managing Aruba Cloud database services (DBaaS - Database as a Service), including databases, users, grants, and backups.

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

### DatabaseAPI

Manage databases within a DBaaS instance:
- List all databases in a DBaaS instance
- Get details of a specific database
- Create or update a database
- Delete a database

### UserAPI

Manage database users:
- List all users in a project
- Get details of a specific user
- Create or update a user
- Delete a user

### GrantAPI

Manage database permissions/grants:
- List all grants for a database
- Get details of a specific grant
- Create or update a grant
- Delete a grant

### BackupAPI

Manage database backups:
- List all backups in a project
- Get details of a specific backup
- Create a backup
- Delete a backup

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
    
    // Initialize API interfaces
    var dbaasAPI database.DBaaSAPI = database.NewDBaaSService(c)
    var databaseAPI database.DatabaseAPI = database.NewDatabaseService(c)
    var userAPI database.UserAPI = database.NewUserService(c)
    var grantAPI database.GrantAPI = database.NewGrantService(c)
    var backupAPI database.BackupAPI = database.NewBackupService(c)
}
```

### DBaaS Management

#### List DBaaS Instances

```go
resp, err := dbaasAPI.ListDBaaS(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list databases: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d DBaaS instances\n", len(resp.Data.Values))
    for _, dbaas := range resp.Data.Values {
        fmt.Printf("DBaaS: %s - Engine: %s\n", dbaas.Metadata.Name, dbaas.Properties.Engine)
    }
}
```

### Database Management

#### List Databases

```go
dbaasID := "dbaas-123"
resp, err := databaseAPI.ListDatabases(ctx, projectID, dbaasID, nil)
if err != nil {
    log.Fatalf("Failed to list databases: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d databases\n", len(resp.Data.Values))
}
```

### User Management

#### List Users

```go
dbaasID := "dbaas-123"
resp, err := userAPI.ListUsers(ctx, projectID, dbaasID, nil)
if err != nil {
    log.Fatalf("Failed to list users: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d users\n", len(resp.Data.Values))
}
```

### Grant Management

#### List Grants

```go
dbaasID := "dbaas-123"
databaseID := "db-456"
resp, err := grantAPI.ListGrants(ctx, projectID, dbaasID, databaseID, nil)
if err != nil {
    log.Fatalf("Failed to list grants: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d grants\n", len(resp.Data.Values))
}
```

### Backup Management

#### List Backups

```go
resp, err := backupAPI.ListBackups(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list backups: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d backups\n", len(resp.Data.Values))
}
```

## Resource Hierarchy

```
Project
└── DBaaS Instance
    ├── Database
    │   └── Grant
    └── User
└── Backup
```
