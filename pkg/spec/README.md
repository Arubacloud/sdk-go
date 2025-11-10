# Spec Package

The `spec` package provides Go client interfaces for all Aruba Cloud services, organized by category. Each category contains specific APIs for managing different aspects of the cloud platform.

## Table of Contents

- [Overview](#overview)
- [Available Categories](#available-categories)
- [Installation](#installation)
- [Quick Start](#quick-start)

## Overview

The spec package is organized into the following categories, each with dedicated interfaces and implementations:

- **Compute** - CloudServers, KeyPairs
- **Network** - VPCs, subnets, security groups, and connectivity
- **Storage** - Block storage volumes and snapshots
- **Database** - Database-as-a-Service (DBaaS) management
- **Security** - KMS key management and encryption
- **Metric** - Monitoring metrics and alerts
- **Audit** - Audit logs and event tracking
- **Schedule** - Schedule Jobs

## Available Categories

### [Compute](./compute/README.md)

Manage compute resources including:
- **Cloud Servers** - Create and manage VMs
- **KeyPairs** - KeyPair

[View Compute Documentation →](./compute/README.md)

### [Network](./network/README.md)

Manage networking resources including:
- **VPCs** - Virtual Private Cloud management
- **Subnets** - Subnet configuration
- **Security Groups** - Network access control
- **Elastic IPs** - Public IP address management
- **Load Balancers** - Load balancing services
- **VPN Tunnels** - VPN connectivity
- **VPC Peering** - VPC interconnection

[View Network Documentation →](./network/README.md)

### [Storage](./storage/README.md)

Manage storage resources including:
- **Block Storage** - Persistent block volumes
- **Snapshots** - Volume snapshots and backups

[View Storage Documentation →](./storage/README.md)

### [Database](./database/README.md)

Manage database services including:
- **DBaaS** - Database instances
- **Databases** - Database management within instances
- **Users** - Database user management
- **Grants** - Permission and access control
- **Backups** - Database backup management

[View Database Documentation →](./database/README.md)

### [Security](./security/README.md)

Manage security services including:
- **KMS Keys** - Key Management Service for encryption

[View Security Documentation →](./security/README.md)

### [Metric](./metric/README.md)

Access monitoring and alerting services including:
- **Metrics** - Performance and resource metrics
- **Alerts** - Alert notifications and history

[View Metric Documentation →](./metric/README.md)

### [Audit](./audit/README.md)

Access audit and compliance services including:
- **Events** - Audit log events and activity tracking

[View Audit Documentation →](./audit/README.md)

### [Schedule](./schedule/README.md)

Schedule or on-demand run jobs:
- **Jobs** - Configure a job to run

[View Schedule Documentation →](./schedule/README.md)

### [Project](./project/README.md)

Manage project resources:
- **Projects** - Project management and configuration

[View Project Documentation →](./project/README.md)


## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    sdkgo "github.com/Arubacloud/sdk-go"
    "github.com/Arubacloud/sdk-go/pkg/client"
)

func main() {
    // Initialize the SDK client
    config := &client.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        Debug:        false,
    }
    
    sdk, err := sdkgo.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // All services are immediately available through the SDK client
    
    // Example: List cloud servers
    resp, err := sdk.Compute.ListCloudServers(ctx, projectID, nil)
    if err != nil {
        log.Fatalf("Failed to list CloudServers: %v", err)
    }
    defer resp.Body.Close()
    
    // Example: List VPCs
    vpcResp, err := sdk.Network.ListVPCs(ctx, projectID, nil)
    if err != nil {
        log.Fatalf("Failed to list VPCs: %v", err)
    }
    defer vpcResp.Body.Close()
}
```

### Direct Service Access

All services are immediately available through the SDK client without any additional initialization:

```go
// Services are pre-initialized and ready to use
sdk.Compute    // Access compute resources (CloudServers, KeyPairs)
sdk.Network    // Access network resources (VPCs, Subnets, etc.)
sdk.Storage    // Access storage resources (Block Storage, Snapshots)
sdk.Database   // Access database services
sdk.Security   // Access security services (KMS)
sdk.Metric     // Access metrics and alerts
sdk.Audit      // Access audit logs
sdk.Schedule   // Access scheduled jobs
sdk.Project    // Access projects
```

This simplifies initialization and provides a cleaner API. 
## Package Structure

```
spec/
├── compute/          # Compute resources
├── network/          # Network resources 
├── storage/          # Storage resources
├── database/         # Database services
├── security/         # Security services
├── metric/           # Monitoring and alerts
├── audit/            # Audit logs and events
├── schedule/         # Schedule Jobs
└── schema/           # Shared data structures and types
```