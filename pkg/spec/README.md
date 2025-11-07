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

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/compute"
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
    "github.com/Arubacloud/sdk-go/pkg/spec/storage"
    "github.com/Arubacloud/sdk-go/pkg/spec/database"
    "github.com/Arubacloud/sdk-go/pkg/spec/security"
    "github.com/Arubacloud/sdk-go/pkg/spec/metric"
    "github.com/Arubacloud/sdk-go/pkg/spec/audit"
    "github.com/Arubacloud/sdk-go/pkg/spec/schedule"
)

func main() {
    // Initialize the client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Compute services
    cloudServerApi := compute.NewCloudServerService(c)
    keyPairApi := compute.NewKeyPairService(c)
    
    // Network services
    vpcApi := network.NewVPCService(c)
    subnetApi := network.NewSubnetService(c)
    elasticIPApi := network.NewElasticIPService(c)
    
    // Storage services
    blockStorageApi := storage.NewBlockStorageService(c)
    snapshotApi := storage.NewSnapshotService(c)
    
    // Database services
    dbaasApi := database.NewDBaaSService(c)
    databaseApi := database.NewDatabaseService(c)
    
    // Security services
    kmsApi := security.NewKMSService(c)
    
    // Metric services
    metricApi := metric.NewMetricService(c)
    alertApi := metric.NewAlertService(c)
    
    // Audit services
    eventApi := audit.NewEventService(c)

    // Schedule services
    jobApi := schedule.NewJobService(c)
    
    // Example: List virtual machines
    resp, err := cloudServerApi.ListCloudServers(ctx, projectID, nil)
    if err != nil {
        log.Fatalf("Failed to list CloudServers: %v", err)
    }
    defer resp.Body.Close()
}
```

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