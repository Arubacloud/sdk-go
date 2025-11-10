# Compute Package

The `compute` package provides Go client interfaces for managing Aruba Cloud compute resources, including cloud servers and SSH key pairs.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)
  - [Cloud Server Management](#cloud-server-management)
  - [Key Pair Management](#key-pair-management)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### ComputeAPI

The unified `ComputeAPI` interface provides all compute-related operations:

**Cloud Server Operations:**
- List all cloud servers in a project
- Get details of a specific cloud server
- Create a new cloud server
- Update an existing cloud server
- Delete a cloud server

**Key Pair Operations:**
- List all key pairs in a project
- Get details of a specific key pair
- Create a new key pair
- Delete a key pair

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
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
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
    
    // Create unified compute service
    computeService := compute.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use computeService for all compute operations
}
```

### Cloud Server Management

#### List Cloud Servers

```go
// Use the unified service for all operations
computeService := compute.NewService(sdk)

resp, err := computeService.ListCloudServers(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list cloud servers: %v", err)
}

// Check response status
if resp.IsSuccess() {
    fmt.Printf("Found %d cloud servers\n", len(resp.Data.Values))
    for _, server := range resp.Data.Values {
        fmt.Printf("Server: %s - Status: %s\n", 
            *server.Metadata.Name, 
            *server.Status.State)
    }
}

// Access HTTP metadata
fmt.Printf("Status Code: %d\n", resp.StatusCode)
```

#### Create Cloud Server

```go
serverReq := schema.CloudServerRequest{
    Metadata: schema.ResourceMetadataRequest{
        Name: "my-server",
        Tags: []string{"production", "web"},
    },
    Properties: schema.CloudServerPropertiesRequest{
        // Server configuration
        BillingPeriod: "Hour",
        // ... other properties
    },
}

createResp, err := computeService.CreateCloudServer(ctx, projectID, serverReq, nil)
if err != nil {
    log.Fatalf("Failed to create server: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created server: %s\n", *createResp.Data.Metadata.Id)
}
```

### Key Pair Management

#### List Key Pairs

```go
// Same service instance for all compute operations
resp, err := computeService.ListKeyPairs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list key pairs: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d key pairs\n", len(resp.Data.Values))
    for _, keypair := range resp.Data.Values {
        fmt.Printf("KeyPair: %s\n", keypair.Metadata.Name)
    }
}
```

#### Create Key Pair

```go
keyPairReq := schema.KeyPairRequest{
    Metadata: schema.ResourceMetadataRequest{
        Name: "my-ssh-key",
    },
    Properties: schema.KeyPairPropertiesRequest{
        Value: "ssh-rsa AAAAB3NzaC1yc2E...",
    },
}

createResp, err := computeService.CreateKeyPair(ctx, projectID, keyPairReq, nil)
if err != nil {
    log.Fatalf("Failed to create key pair: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created key pair: %s\n", createResp.Data.Metadata.Name)
}
```

## Interface Definition

```go
type ComputeAPI interface {
    // CloudServer operations
    ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerList], error)
    GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
    CreateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
    UpdateCloudServer(ctx context.Context, project string, cloudServerId string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
    DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[any], error)

    // KeyPair operations
    ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairListResponse], error)
    GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
    CreateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
    DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*schema.Response[any], error)
}
```
