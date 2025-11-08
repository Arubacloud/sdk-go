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

### VpcAPI
Manage Virtual Private Clouds (VPCs)

### SubnetAPI
Manage subnets within VPCs

### SecurityGroupAPI
Manage security groups for network access control

### SecurityGroupRuleAPI
Manage security group rules

### ElasticIPAPI
Manage elastic IP addresses

### LoadBalancerAPI
Manage load balancers (read-only)

### VpcPeeringAPI
Manage VPC peering connections

### VpcPeeringRouteAPI
Manage routes for VPC peering connections

### VpnTunnelAPI
Manage VPN tunnels

### VpnRouteAPI
Manage VPN routes

## Usage Examples

### Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Arubacloud/sdk-go/pkg/client"
    "github.com/Arubacloud/sdk-go/pkg/spec/network"
    "github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func main() {
    // Create a new client
    c := client.NewClient("https://api.arubacloud.com", "your-api-key")
    
    ctx := context.Background()
    projectID := "my-project-id"
    
    // Initialize API interfaces
    var vpcAPI network.VpcAPI = network.NewVPCService(c)
    var subnetAPI network.SubnetAPI = network.NewSubnetService(c)
    var securityGroupAPI network.SecurityGroupAPI = network.NewSecurityGroupService(c)
    var securityGroupRuleAPI network.SecurityGroupRuleAPI = network.NewSecurityGroupRuleService(c)
    var elasticIPAPI network.ElasticIPAPI = network.NewElasticIPService(c)
    var loadBalancerAPI network.LoadBalancerAPI = network.NewLoadBalancerService(c)
    var vpcPeeringAPI network.VpcPeeringAPI = network.NewVpcPeeringService(c)
    var vpcPeeringRouteAPI network.VpcPeeringRouteAPI = network.NewVpcPeeringRouteService(c)
    var vpnTunnelAPI network.VpnTunnelAPI = network.NewVpnTunnelService(c)
    var vpnRouteAPI network.VpnRouteAPI = network.NewVpnRouteService(c)
}
```

### VPC Management

#### List VPCs

```go
resp, err := vpcAPI.ListVPCs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list VPCs: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d VPCs\n", len(resp.Data.Values))
    for _, vpc := range resp.Data.Values {
        fmt.Printf("VPC: %s - CIDR: %s\n", vpc.Metadata.Name, vpc.Properties.Cidr)
    }
}
```

### Subnet Management

#### List Subnets

```go
vpcID := "vpc-123"
resp, err := subnetAPI.ListSubnets(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("Failed to list subnets: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d subnets\n", len(resp.Data.Values))
}
```

### Security Group Management

#### List Security Groups

```go
vpcID := "vpc-123"
resp, err := securityGroupAPI.ListSecurityGroups(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("Failed to list security groups: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d security groups\n", len(resp.Data.Values))
}
```

### Security Group Rule Management

#### List Security Group Rules

```go
vpcID := "vpc-123"
securityGroupID := "sg-456"
resp, err := securityGroupRuleAPI.ListSecurityGroupRules(ctx, projectID, vpcID, securityGroupID, nil)
if err != nil {
    log.Fatalf("Failed to list security group rules: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d security group rules\n", len(resp.Data.Values))
}
```

### Elastic IP Management

#### List Elastic IPs

```go
resp, err := elasticIPAPI.ListElasticIPs(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list elastic IPs: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d elastic IPs\n", len(resp.Data.Values))
}
```

### Load Balancer Management

#### List Load Balancers

```go
resp, err := loadBalancerAPI.ListLoadBalancers(ctx, projectID, nil)
if err != nil {
    log.Fatalf("Failed to list load balancers: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d load balancers\n", len(resp.Data.Values))
}
```

### VPC Peering Management

#### List VPC Peerings

```go
vpcID := "vpc-123"
resp, err := vpcPeeringAPI.ListVpcPeerings(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("Failed to list VPC peerings: %v", err)
}

// Access response data
if resp.IsSuccess() {
    fmt.Printf("Found %d VPC peerings\n", len(resp.Data.Values))
}
```

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
