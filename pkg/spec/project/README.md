# Project Package

The `project` package provides Go client interfaces for managing Aruba Cloud projects.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ProjectAPI

Manage projects with full CRUD operations:
- List all projects
- Get details of a specific project
- Create or update a project
- Delete a project

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/project"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    
    // Initialize API interface
    var projectAPI project.ProjectAPI = project.NewProjectService(c)
}
```

### Project Management

#### List Projects

```go
resp, err := projectAPI.ListProjects(ctx, nil)
if err != nil {
    log.Fatalf("Failed to list projects: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d projects\n", len(resp.Data.Values))
    for _, project := range resp.Data.Values {
        fmt.Printf("Project: %s\n", project.Metadata.Name)
    }
}

// Access HTTP metadata
fmt.Printf("Status Code: %d\n", resp.StatusCode)
```

