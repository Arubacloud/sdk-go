# Compute Package

The `compute` package provides Go client interfaces for managing Aruba Cloud compute resources, including cloud servers and SSH key pairs.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)
  - [Cloud Server Management](#cloud-server-management)
  - [Key Pair Management](#key-pair-management)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### CloudServerAPI

Manage cloud server instances with full CRUD operations:
- List all cloud servers in a project
- Get details of a specific cloud server
- Create or update a cloud server
- Delete a cloud server

### KeyPairAPI

Manage SSH key pairs for secure server access:
- List all key pairs in a project
- Get details of a specific key pair
- Create or update a key pair
- Delete a key pair

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Use the API interfaces
    var cloudServerAPI compute.CloudServerAPI
    var keyPairAPI compute.KeyPairAPI
    
    // Initialize services
    cloudServerAPI = compute.NewCloudServerService(c)
    keyPairAPI = compute.NewKeyPairService(c)
}
```

### Cloud Server Management

#### List Cloud Servers

```go
// Using the CloudServerAPI interface
var cloudServerAPI compute.CloudServerAPI = compute.NewCloudServerService(c)

resp, err := cloudServerAPI.ListCloudServers(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list cloud servers: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode == 200 {
    fmt.Println("Cloud servers retrieved successfully")
}
```

### Key Pair Management

#### Using the KeyPairAPI Interface

```go
// Initialize the KeyPairAPI interface
var keyPairAPI compute.KeyPairAPI = compute.NewKeyPairService(c)
```

#### List Key Pairs

```go
resp, err := keyPairAPI.ListKeyPairs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list key pairs: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode == 200 {
    fmt.Println("Key pairs retrieved successfully")
}
```

