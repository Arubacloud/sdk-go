# Aruba Cloud SDK for Go

[![Build](https://github.com/Arubacloud/sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Arubacloud/sdk-go/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Arubacloud/sdk-go)](https://github.com/Arubacloud/sdk-go/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/Arubacloud/sdk-go)](https://github.com/Arubacloud/sdk-go/releases)
[![Codecov](https://codecov.io/gh/Arubacloud/sdk-go/branch/main/graph/badge.svg)](https://codecov.io/gh/Arubacloud/sdk-go)
[![License](https://img.shields.io/github/license/Arubacloud/sdk-go)](LICENSE)

See [`CHANGELOG.md`](CHANGELOG.md) for the full release history and the v0.1.x → v0.2.x branch policy.

> **Note**: This SDK is currently in its **Alpha** stage. The API is not yet stable, and breaking changes may be introduced in future releases without prior notice.

Official Go SDK for the Aruba Cloud API. Manage cloud resources — compute instances, VPCs, storage, databases, and more — from Go code.

## Highlights

- **Fluent builder API** — construct any resource with chained, type-safe setters: `aruba.NewVPC().IntoProject(p).Named("prod")...`
- **Single import** — `pkg/aruba` re-exports every typed enum and factory; `pkg/types` is reserved for advanced escape hatches.
- **Built-in async polling** — `wrapper.WaitUntilReady(ctx)` covers 95% of long-running operations; specialized waits cover the rest.
- **Configurable authentication** — OAuth2 client credentials by default, plus static tokens and Vault-backed credential storage.
- **Multi-tenant client management** out of the box via `pkg/multitenant`.

## Table of Contents

- [1. Quick Start](#1-quick-start)
- [2. Usage Details](#2-usage-details)
  - [2.1. Configuring the Client](#21-configuring-the-client)
  - [2.2. Making Calls, Filtering, and Handling Responses](#22-making-calls-filtering-and-handling-responses)
  - [2.3. SDK Client API Reference](#23-sdk-client-api-reference)
  - [2.4. Async Operations](#24-async-operations)
  - [2.5. Multi-Tenancy](#25-multi-tenancy)
  - [2.6. Typed Enums & Escape Hatches](#26-typed-enums--escape-hatches)

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
)

func main() {
	arubaClient, err := aruba.NewClient(
		aruba.DefaultOptions("your-client-id", "your-client-secret"),
	)
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	proj := aruba.NewProject().
		Named("my-first-project").
		WithDescription("Created from the Go SDK Quick Start").
		AddTag("go-sdk")

	if _, err := arubaClient.FromProject().Create(ctx, proj); err != nil {
		log.Fatalf("create project: %v", err)
	}

	fmt.Printf("Created project %s (id=%s)\n", proj.Name(), proj.ID())
}
```

Project create is synchronous — no `WaitUntilReady` needed. For an end-to-end Project → VPC → Subnet → CloudServer flow with async waits and cleanup, see [`docs/website/docs/walkthrough.md`](docs/website/docs/walkthrough.md) or the runnable [`examples/all-resources/`](examples/all-resources/).

## 2. Usage Details

### 2.1. Configuring the Client

`aruba.DefaultOptions(clientID, clientSecret)` is the recommended starting point. For full manual control, build from `aruba.NewOptions()`.

```go
opts := aruba.DefaultOptions("your-client-id", "your-client-secret").
    WithNativeLogger().                 // enable built-in logging
    WithCustomHTTPClient(myHTTPClient)  // custom *http.Client / transport
```

Authentication (mutually exclusive):

- `WithClientCredentials(id, secret)` — OAuth2 client credentials (default via `DefaultOptions`).
- `WithToken(staticToken)` — supply a pre-issued bearer token.
- `WithVaultCredentialsRepository(...)` — fetch credentials from HashiCorp Vault per request.

Logging: `WithNoLogs()` (default), `WithNativeLogger()`, or `WithCustomLogger(myLogger)`.

Token caching (optional): `WithRedisTokenRepositoryFromURI(...)` or `WithFileTokenRepositoryFromBaseDir(...)`.

For the full options reference see [`docs/website/docs/options.md`](docs/website/docs/options.md).

### 2.2. Making Calls, Filtering, and Handling Responses

#### Calls

Resource clients expose a uniform `Create / List / Get / Update / Delete` shape:

```go
proj := aruba.NewProject().Named("staging")
_, err := arubaClient.FromProject().Create(ctx, proj)

list, err := arubaClient.FromCompute().CloudServers().List(ctx, proj)
vpc, err  := arubaClient.FromNetwork().VPCs().Get(ctx, vpcRef)
_, err     = arubaClient.FromStorage().Volumes().Delete(ctx, volRef)
```

Any `Ref` works as input — a wrapper, a sub-client list item, or `aruba.URI("/projects/<id>")` to bootstrap from a string.

#### Filters & paging

Filters, sort, and paging are passed as variadic `CallOption`s:

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,cpu:gt:2"),
    aruba.WithSort("name:asc"),
    aruba.WithLimit(50),
)
```

Filter syntax: `field:operator:value` joined by `,` (AND) or `;` (OR).
Operators: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`, `in`, `nin`, `like`, `sw`, `ew`.
See [`docs/website/docs/filters.md`](docs/website/docs/filters.md) for the full reference.

#### Responses

Wrapper methods return `(wrapper, error)`. API errors come back as `*aruba.HTTPError`:

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpcRef)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        log.Fatalf("API error %d: %s", httpErr.StatusCode, httpErr.Message)
    }
    log.Fatalf("transport error: %v", err)
}
fmt.Println(vpc.Name())
```

Setter-time errors are deferred — surface them with `wrapper.Err()` or by calling `Create`/`Update`. See [`docs/website/docs/response-handling.md`](docs/website/docs/response-handling.md).

### 2.3. SDK Client API Reference

| Domain | Sub-clients |
|---|---|
| `FromAudit()` | `Events()` |
| `FromCompute()` | `CloudServers()`, `KeyPairs()` |
| `FromContainer()` | `KaaS()`, `ContainerRegistry()` |
| `FromDatabase()` | `DBaaS()`, `Databases()`, `Backups()`, `Users()`, `Grants()` |
| `FromMetric()` | `Alerts()`, `Metrics()` |
| `FromNetwork()` | `ElasticIPs()`, `LoadBalancers()`, `SecurityGroups()`, `SecurityGroupRules()`, `Subnets()`, `VPCs()`, `VPCPeerings()`, `VPCPeeringRoutes()`, `VPNTunnels()`, `VPNRoutes()` |
| `FromProject()` | (single client; `Create`, `List`, `Get`, `Update`, `Delete` directly) |
| `FromSchedule()` | `Jobs()` |
| `FromSecurity()` | `KMS()`, `Keys()`, `Kmips()` |
| `FromStorage()` | `Volumes()`, `Snapshots()`, `Backups()`, `Restores()` |

Every leaf client provides `Create`, `List`, `Get`, `Update`, `Delete` (where applicable).

For per-resource builders (`aruba.NewVPC()`, `aruba.NewKaaS()`, …) and the full type catalog see [`docs/website/docs/resources.md`](docs/website/docs/resources.md).

### 2.4. Async Operations

Most cloud resources reach their terminal state asynchronously. Every wrapper that has a status mixin exposes:

- `WaitUntilReady(ctx, opts...)` — accepts `Active`, `NotUsed`, `InUse`, `Used`. The 95% answer.
- `WaitUntilActive(ctx, opts...)` — strictly `Active`.
- `WaitUntilStates(ctx, []string{"Foo", "Bar"}, opts...)` — arbitrary target list.
- Specialized: `BlockStorage`/`ElasticIP` add `WaitUntilNotUsed`/`WaitUntilUsed`; `Kmip` adds `WaitUntilCertificateAvailable`.

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx,
    aruba.NewVPC().IntoProject(proj).Named("prod-vpc"))
if err != nil { log.Fatal(err) }

if err := vpc.WaitUntilReady(ctx); err != nil {
    log.Fatalf("vpc never reached Active: %v", err)
}
```

Tunable via `aruba.WithRetries(n)`, `aruba.WithBaseDelay(d)`, `aruba.WithTimeout(d)`.
Defaults: 60 retries × 10 s × 600 s ceiling.

For concurrent waits or polling on arbitrary conditions (e.g. polling for an HTTP 404 after `Delete`), use `pkg/async.DefaultWaitFor` directly — see [`docs/website/docs/async.md`](docs/website/docs/async.md) and the `waitUntilGone` helper in [`examples/all-resources/common.go`](examples/all-resources/common.go).

### 2.5. Multi-Tenancy

`pkg/multitenant` keeps an in-memory `tenant → aruba.Client` registry. Use it when one process serves many Aruba accounts (e.g. a reconciler).

```go
mt := multitenant.New()

// Static credentials per tenant
mt.NewFromOptions("tenant-1",
    aruba.DefaultOptions(tenant1ID, tenant1Secret))

// Vault-backed credentials sharing a base options template
base := aruba.NewOptions().WithDefaultBaseURL().WithDefaultTokenIssuerURL()
mt.NewFromOptions("tenant-2",
    base.DeepCopy().WithVaultCredentialsRepository(vaultURI, kvMount, "tenant-2", ns, rolePath, roleID, secretID))

client, ok := mt.Get("tenant-1")
```

Optional periodic eviction of idle tenants:

```go
multitenant.StartCleanupRoutine(ctx, mt, 5*time.Minute, 1*time.Hour)
```

For the full reconciler-style cache-or-create pattern see [`docs/website/docs/multitenancy.md`](docs/website/docs/multitenancy.md) and [`examples/all-resources/orchestrator_multitenancy.go`](examples/all-resources/orchestrator_multitenancy.go).

### 2.6. Typed Enums & Escape Hatches

Typed constants for regions, zones, flavors, billing periods, Kubernetes versions, VPN crypto, and more are re-exported from `pkg/aruba` — you almost never need a second import:

```go
cs := aruba.NewCloudServer().
    IntoProject(proj).
    InRegion(aruba.RegionITBGBergamo).
    InZone(aruba.ZoneITBG1).
    OfFlavor(aruba.CloudServerFlavorCSO4A8).
    WithBillingPeriod(aruba.BillingPeriodHour)
```

When the wrapper doesn't yet expose a wire field, fall back via `wrapper.Raw()`:

```go
raw := vpc.Raw()                    // *types.VPCResponse
fmt.Println(raw.Properties.Whatever)
```

The full enum catalog and resource→type cross-reference lives in [`docs/website/docs/resources.md`](docs/website/docs/resources.md).
