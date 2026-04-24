# Aruba Cloud SDK for Go

[![Build](https://github.com/Arubacloud/sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Arubacloud/sdk-go/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Arubacloud/sdk-go)](https://github.com/Arubacloud/sdk-go/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/Arubacloud/sdk-go)](https://github.com/Arubacloud/sdk-go/releases)
[![Codecov](https://codecov.io/gh/Arubacloud/sdk-go/branch/main/graph/badge.svg)](https://codecov.io/gh/Arubacloud/sdk-go)
[![License](https://img.shields.io/github/license/Arubacloud/sdk-go)](LICENSE)

> **Note**: This SDK is currently in its **Alpha** stage. The API is not yet stable, and breaking changes may be introduced in future releases without prior notice.

Official Go SDK for the Aruba Cloud API. Manage cloud resources — compute instances, VPCs, storage, databases, and more — from Go code.

## Highlights

- Strongly-typed request/response structs for all resources
- Automatic OAuth2 token management with configurable refresh strategies
- Generic `Response[T]` wrapper with consistent error handling across all calls
- Built-in async polling helpers for long-running operations
- Multi-tenant client management out of the box

## Table of Contents

- [1. Quick Start](#1-quick-start)
- [2. Usage Details](#2-usage-details)
  - [2.1. Config Options](#21-config-options)
  - [2.2. Performing Calls, Setting Filters, and Handling Responses](#22-performing-calls-setting-filters-and-handling-responses)
  - [2.3. SDK Client API Reference](#23-sdk-client-api-reference)
  - [2.4. Handling Asynchronous Operations](#24-handling-asynchronous-operations)
  - [2.5. Data Types](#25-data-types)

## 1. Quick Start

Import the `aruba` package, create a client with your credentials, and start making API calls.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	aruba_types "github.com/Arubacloud/sdk-go/pkg/types"
)

func main() {
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions("your-client-id", "your-client-secret"))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	createResp, err := arubaClient.FromProject().Create(ctx, aruba_types.ProjectRequest{
		Metadata: aruba_types.ResourceMetadataRequest{
			Name: "my-first-project",
			Tags: []string{"go-sdk", "quick-start"},
		},
	}, nil)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}
	if !createResp.IsSuccess() {
		log.Fatalf("API error %d: %s", createResp.StatusCode, *createResp.Error.Title)
	}

	fmt.Printf("Created project: %s\n", *createResp.Data.Metadata.ID)
}
```

## 2. Usage Details

### 2.1. Config Options

Configure the client via the `Options` fluent builder. `aruba.DefaultOptions` covers the most common case.

```go
options := aruba.DefaultOptions("your-client-id", "your-client-secret")
options.WithNativeLogger()          // enable built-in logging
options.WithCustomHTTPClient(hc)    // supply a custom *http.Client
```

Key areas:

- **Authentication** — OAuth2 client credentials by default. Use `WithToken()` for a static token, or configure Vault-backed credentials for secrets management.
- **Logging** — disabled by default. Enable with `WithNativeLogger()` or inject your own `logger.Logger`.
- **HTTP client** — defaults to `http.DefaultClient`. Override with `WithCustomHTTPClient()` to set timeouts or a custom transport.

For the full options reference see [`docs/website/docs/options.md`](docs/website/docs/options.md).

### 2.2. Performing Calls, Setting Filters, and Handling Responses

#### API calls

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, nil)
vpc, err    := arubaClient.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)
_, err       = arubaClient.FromStorage().Volumes().Delete(ctx, projectID, volumeID, nil)
```

#### Filters

```go
// Filters follow the format "field:operator:value"; combine with "," (AND) or ";" (OR).
filter := "status:eq:running,cpu:gt:2"

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID,
    &types.RequestParameters{Filter: &filter})
```

Supported operators: `eq`, `ne`, `gt`, `lt`, `in`, `contains`, `startswith`, `endswith`.

For the full filtering guide see [`docs/website/docs/filters.md`](docs/website/docs/filters.md).

#### Responses

Every call returns `*types.Response[T]` and `error`. Check `error` first (network/client issues), then `IsSuccess()` / `IsError()` for the API result.

```go
resp, err := arubaClient.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)
if err != nil {
    log.Fatalf("request failed: %v", err)
}
if resp.IsSuccess() {
    fmt.Println(*resp.Data.Metadata.Name)
} else {
    log.Printf("API error %d: %s", resp.StatusCode, *resp.Error.Title)
}
```

For detailed error handling patterns see [`docs/website/docs/response-handling.md`](docs/website/docs/response-handling.md).

### 2.3. SDK Client API Reference

#### Service groups

| Accessor | Resources |
|---|---|
| `FromAudit()` | Audit events |
| `FromCompute()` | Cloud servers, key pairs |
| `FromContainer()` | Kubernetes (KaaS), container registries |
| `FromDatabase()` | Database-as-a-Service instances |
| `FromMetric()` | Metrics and alerts |
| `FromNetwork()` | VPCs, subnets, security groups, elastic IPs, load balancers, VPN tunnels, VPC peerings |
| `FromProject()` | Projects |
| `FromSchedule()` | Scheduled jobs |
| `FromSecurity()` | KMS keys |
| `FromStorage()` | Volumes, snapshots, backups, restores |

Each resource client provides `Create`, `List`, `Get`, `Update`, `Delete` where applicable.

For the full resource listing see [`docs/website/docs/resources.md`](docs/website/docs/resources.md).

### 2.4. Handling Asynchronous Operations

Many cloud operations are asynchronous. Use `async.DefaultWaitFor` to poll until the resource reaches the desired state.

```go
future := async.DefaultWaitFor(
    ctx,
    func(ctx context.Context) (*types.Response[types.CloudServerResponse], error) {
        return arubaClient.FromCompute().CloudServers().Get(ctx, projectID, serverID, nil)
    },
    func(resp *types.Response[types.CloudServerResponse]) (bool, error) {
        return resp.IsSuccess() && resp.Data.Properties.Status == "running", nil
    },
)

result, err := future.Await(ctx)
if err != nil {
    log.Fatalf("wait failed: %v", err)
}
fmt.Printf("server status: %s\n", result.Data.Properties.Status)
```

`async.DefaultWaitFor` retries 10 times with a 10 s delay and a 60 s total timeout. Use `async.WaitFor` for custom retry counts, delay, and timeout.

### 2.5. Data Types

All models live in `pkg/types`:

- **Requests** end in `Request` — e.g. `types.VPCRequest`, `types.CloudServerRequest`
- **Single responses** end in `Response` — e.g. `types.VPCResponse`
- **Collection responses** end in `List` — e.g. `types.VPCList`
- **Shared structures** — `ResourceMetadataRequest`, `ResourceMetadataResponse`, `ResourceStatus`

For the full type reference see [`docs/website/docs/types.md`](docs/website/docs/types.md).
