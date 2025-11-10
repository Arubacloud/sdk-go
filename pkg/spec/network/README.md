# Network Package

The `network` package provides Go client interfaces for managing Aruba Cloud network services, including VPCs, subnets, security groups, elastic IPs, load balancers, VPN tunnels, and VPC peering.

## Table of Contents

- [Installation](#installation)
- [Available Services](#available-services)
- [Usage Examples](#usage-examples)

## Installation

```bash
go get github.com/Arubacloud/sdk-go
```

## Available Services

### NetworkAPI

The unified `NetworkAPI` interface provides all network-related operations:

**VPC Operations** - Manage Virtual Private Clouds  
**Subnet Operations** - Manage subnets within VPCs  
**Security Group Operations** - Manage security groups for network access control  
**Security Group Rule Operations** - Manage security group rules  
**Elastic IP Operations** - Manage elastic IP addresses  
**Load Balancer Operations** - View load balancers (read-only)  
**VPC Peering Operations** - Manage VPC peering connections  
**VPC Peering Route Operations** - Manage routes for VPC peering connections  
**VPN Tunnel Operations** - Manage VPN tunnels  
**VPN Route Operations** - Manage VPN routes  

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
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
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
    
    // Create unified network service
    networkService := network.NewService(sdk)
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Now use networkService for all network operations
}
```

### VPC Management

#### List VPCs

```go
// Use the unified service for all network operations
networkService := network.NewService(sdk)

resp, err := networkService.ListVPCs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list VPCs: %v", err)
}

// Check response status
if resp.IsSuccess() {
    fmt.Printf("Found %d VPCs\n", len(resp.Data.Values))
    for _, vpc := range resp.Data.Values {
        fmt.Printf("VPC: %s - CIDR: %s\n", 
            *vpc.Metadata.Name, 
            vpc.Properties.Cidr)
    }
}
    }
}
```

#### Create VPC

```go
vpcReq := schema.VpcRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "my-vpc",
            Tags: []string{"production"},
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.VpcPropertiesRequest{
        Properties: &schema.VpcProperties{
            Default: boolPtr(false),
            Preset:  boolPtr(true),
        },
    },
}

createResp, err := networkService.CreateVPC(ctx, projectID, vpcReq, nil)
if err != nil {
    log.Fatalf("Failed to create VPC: %v", err)
}

if createResp.IsSuccess() {
    fmt.Printf("Created VPC: %s\n", *createResp.Data.Metadata.Id)
}
```

### Subnet Management

#### List Subnets

```go
vpcID := "vpc-123"
resp, err := networkService.ListSubnets(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("Failed to list subnets: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d subnets\n", len(resp.Data.Values))
}
```

#### Create Subnet

```go
subnetReq := schema.SubnetRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "my-subnet",
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.SubnetPropertiesRequest{
        Type:    schema.SubnetTypeAdvanced,
        Default: true,
        Network: &schema.SubnetNetwork{
            Address: "192.168.1.0/25",
        },
    },
}

createResp, err := networkService.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)
```

### Security Group Management

#### List Security Groups

```go
vpcID := "vpc-123"
resp, err := networkService.ListSecurityGroups(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("Failed to list security groups: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d security groups\n", len(resp.Data.Values))
}
```

### Security Group Rule Management

#### Create Security Group Rule

```go
vpcID := "vpc-123"
securityGroupID := "sg-456"

ruleReq := schema.SecurityRuleRequest{
    Metadata: schema.RegionalResourceMetadataRequest{
        ResourceMetadataRequest: schema.ResourceMetadataRequest{
            Name: "allow-ssh",
        },
        Location: schema.LocationRequest{
            Value: "ITBG-Bergamo",
        },
    },
    Properties: schema.SecurityRulePropertiesRequest{
        Direction: schema.RuleDirectionIngress,
        Protocol:  "TCP",
        Port:      "22",
        Target: &schema.RuleTarget{
            Kind:  schema.EndpointTypeIP,
            Value: "0.0.0.0/0",
        },
    },
}

resp, err := networkService.CreateSecurityGroupRule(ctx, projectID, vpcID, securityGroupID, ruleReq, nil)
```

### Elastic IP Management

#### List Elastic IPs

```go
resp, err := networkService.ListElasticIPs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list elastic IPs: %v", err)
}

if resp.IsSuccess() {
    fmt.Printf("Found %d elastic IPs\n", len(resp.Data.Values))
}
```

### All Operations Available

The unified `NetworkAPI` service provides access to all network operations:

- **VPC**: Create, Get, Update, Delete, List
- **Subnet**: Create, Get, Update, Delete, List  
- **SecurityGroup**: Create, Get, Update, Delete, List
- **SecurityGroupRule**: Create, Get, Update, Delete, List
- **ElasticIP**: Create, Get, Update, Delete, List
- **LoadBalancer**: Get, List
- **VpcPeering**: Create, Get, Update, Delete, List
- **VpcPeeringRoute**: Create, Get, Update, Delete, List
- **VpnTunnel**: Create, Get, Update, Delete, List
- **VpnRoute**: Create, Get, Update, Delete, List

### VPC Peering Route Management

#### List VPC Peering Routes

```go
vpcID := "vpc-123"
vpcPeeringID := "peering-456"
resp, err := vpcPeeringRouteAPI.ListVpcPeeringRoutes(ctx, projectID, vpcID, vpcPeeringID, nil)
if err != nil {
    log.Fatalf("Failed to list VPC peering routes: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d VPC peering routes\n", len(resp.Data.Values))
}
```

### VPN Tunnel Management

#### List VPN Tunnels

```go
resp, err := vpnTunnelAPI.ListVpnTunnels(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list VPN tunnels: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d VPN tunnels\n", len(resp.Data.Values))
}
```

### VPN Route Management

#### List VPN Routes

```go
vpnTunnelID := "vpn-123"
resp, err := vpnRouteAPI.ListVpnRoutes(ctx, projectID, vpnTunnelID, nil)
if err != nil {
    log.Fatalf("Failed to list VPN routes: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d VPN routes\n", len(resp.Data.Values))
}
```

## Resource Hierarchy

```
Project
├── VPC
│   ├── Subnet
│   ├── Security Group
│   │   └── Security Group Rule
│   └── VPC Peering
│       └── VPC Peering Route
├── Elastic IP
├── Load Balancer (read-only)
└── VPN Tunnel
    └── VPN Route
```
