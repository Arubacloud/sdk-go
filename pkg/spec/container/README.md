# Container Package

This package provides services for managing Aruba Cloud Container resources.

## Services

### KaaS (Kubernetes as a Service)

Manage Kubernetes clusters with the KaaS service:

```go
import "github.com/Arubacloud/sdk-go/pkg/spec/container"

kaasAPI := container.NewKaaSService(sdk)

// Create a KaaS cluster
kaasResp, err := kaasAPI.CreateKaaS(ctx, projectID, kaasRequest, nil)

// Get cluster details
kaasResp, err := kaasAPI.GetKaaS(ctx, projectID, kaasID, nil)

// List clusters
kaasList, err := kaasAPI.ListKaaS(ctx, projectID, nil)

// Delete cluster
_, err := kaasAPI.DeleteKaaS(ctx, projectID, kaasID, nil)
```

## Features

- Create and manage Kubernetes clusters
- Configure node pools with different instance types
- High availability (HA) support
- Kubernetes version management
- VPC, Subnet, and Security Group integration
- Persistent storage configuration
