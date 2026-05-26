---
sidebar_position: 2
---

# API Walkthrough

The Aruba Cloud Go SDK gives you a single import — `github.com/Arubacloud/sdk-go/pkg/aruba` — that exposes a fluent builder API for every cloud resource. You construct a resource description with a `aruba.NewX()` builder chain, pass it to the appropriate client method (`Create`, `Get`, `Update`, `Delete`, or `List`), and work with the typed wrapper that comes back.

Resources are scoped to a **Project**, and child resources reference their parents via the `aruba.Ref` interface. You never have to extract or thread raw ID strings by hand: pass the hydrated wrapper (returned by `Create` or `Get`) directly as a `Ref` parameter to builder methods like `InProject(proj)`, `InVPC(vpc)`, or `InSecurityGroup(sg)`.

This page walks through the core CRUD lifecycle on a minimal example — Project + VPC + Subnet. Every other resource follows the exact same shape. See [Resources](./resources) for copy-paste-ready snippets for all supported resources.

---

## 1. Initialise the Client

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/Arubacloud/sdk-go/pkg/aruba"
)

func main() {
    arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()
}
```

`aruba.NewClient` accepts an `*aruba.Options` value. `aruba.DefaultOptions(clientID, clientSecret)` is the fastest way to get started; see [Configuration Options](./options) for Vault credentials, Redis token caching, custom loggers, and more.

The returned `aruba.Client` is a façade that exposes domain-specific sub-clients:
`FromProject()`, `FromAudit()`, `FromCompute()`, `FromContainer()`, `FromDatabase()`, `FromMetric()`, `FromNetwork()`, `FromSchedule()`, `FromSecurity()`, `FromStorage()`.

---

## 2. Provision Resources

Resources are created inline: build the request with `aruba.NewX()`, pass it directly to `Create`. The returned wrapper carries the resource ID and URI — pass it as a `Ref` to child resource builders.

### Project

The Project is the top-level container. Every other resource belongs to a project. It is synchronously ready after `Create` returns — no polling required.

```go
proj, err := arubaClient.FromProject().Create(
    ctx,
    aruba.NewProject().
        Named("my-project").
        Tagged("go-sdk").
        DescribedAs("Created via the Aruba Cloud Go SDK").
        NotDefault())
if err != nil {
    log.Fatalf("Create project: %v", err)
}
fmt.Printf("✓ Project created: %s (ID: %s)\n", proj.Name(), proj.ID())
```

### VPC

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(
    ctx,
    aruba.NewVPC().
        Named("my-vpc").
        Tagged("network").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        NotDefault().
        WithoutPreset())
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}
fmt.Printf("✓ VPC created: %s\n", vpc.Name())

// Most resources are asynchronous — wait until they reach Active state.
// See "7. Wait for Readiness" below for options and details.
if err := vpc.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

`InProject(proj)` accepts any `aruba.Ref` — it binds the project scope without requiring you to extract a raw ID string.

### Subnet

```go
subnet, err := arubaClient.FromNetwork().Subnets().Create(
    ctx,
    aruba.NewSubnet().
        OfType(aruba.SubnetTypeAdvanced).
        Named("my-subnet").
        Tagged("network").
        InVPC(vpc).
        InRegion(aruba.RegionITBGBergamo).
        WithCIDR("192.168.1.0/25").
        WithDHCP(aruba.NewSubnetDHCP().
            Enabled().
            WithRange("192.168.1.10", 50).
            WithRoutes(aruba.SubnetDHCPRoute{Address: "10.0.0.0/8", Gateway: "192.168.1.1"}).
            WithDNSServers("8.8.8.8", "8.8.4.4")).
        NotDefault())
if err != nil {
    log.Fatalf("Create subnet: %v", err)
}
fmt.Printf("✓ Subnet created: %s (CIDR: %s)\n", subnet.Name(), subnet.CIDR())

if err := subnet.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Subnet did not become Active: %v", err)
}
```

`aruba.NewSubnetDHCP()` is a sub-builder for DHCP configuration. Attach it to the subnet with `WithDHCP(...)`.

`OfType` accepts `aruba.SubnetTypeBasic` or `aruba.SubnetTypeAdvanced` (typed constants — no string cast needed).

> Every other resource — Security Groups, Elastic IPs, Block Storage, Cloud Servers, KaaS clusters, DBaaS instances, and more — follows the exact same `NewX()` → `IntoParent(ref)` → `Create(ctx, ...)` → `WaitUntilReady(ctx)` shape. See [Resources](./resources) for the full list with copy-paste-ready snippets.

---

## 3. Update an Existing Resource

Fetch the resource first, mutate via setters, then call `Update`. The response wrapper from `Get` carries all internal state (parent URIs, network refs, etc.) that round-trips automatically into the `Update` request.

```go
// Fetch
vpc, err = arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
if err != nil {
    log.Fatalf("Get VPC: %v", err)
}

// Mutate
vpc.Named("my-vpc-updated").
    RetaggedAs("network", "updated")

// Update
updated, err := arubaClient.FromNetwork().VPCs().Update(ctx, vpc)
if err != nil {
    log.Fatalf("Update VPC: %v", err)
}
fmt.Printf("✓ VPC updated: %s\n", updated.Name())
```

> **Important**: Always call `Get` before `Update`. Calling `Update` on a freshly-built wrapper (with no prior `Create` or `Get`) returns an error: `"Update: resource has no ID"`.

---

## 4. List Existing Resources

`List` takes a parent `Ref` and returns a `*aruba.List[T]`. Iterate the items with `Items()`:

```go
list, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil {
    log.Fatalf("List VPCs: %v", err)
}
fmt.Println("Total VPCs:", list.Total())
for _, v := range list.Items() {
    fmt.Println("-", v.Name(), v.ID())
}
```

Items in the list are lightweight wrappers — they carry the resource ID and URI, so you can pass them directly to `Get`, `Update`, or `Delete` as a `Ref`:

```go
for _, v := range list.Items() {
    full, err := arubaClient.FromNetwork().VPCs().Get(ctx, v)
    // full has all fields populated
}
```

For server-side filtering, sorting, and pagination see [Filters](./filters).

---

## 5. Get a Specific Resource

Use `Get` when you have a `Ref` (a hydrated wrapper, or a `*aruba.List[T]` item):

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
if err != nil {
    log.Fatalf("Get VPC: %v", err)
}
```

### The `aruba.URI(…)` escape hatch

When you only have a resource identifier as a string — for example, read from an environment variable or external config — wrap it in `aruba.URI(…)` to satisfy the `aruba.Ref` interface:

```go
projectID := os.Getenv("PROJECT_ID")

// Bootstrap a typed wrapper from a string ID
proj, err := arubaClient.FromProject().Get(ctx, aruba.URI("/projects/"+projectID))
if err != nil {
    log.Fatalf("Get project: %v", err)
}

// Now proj is fully hydrated — use it as a Ref for child resources
vpcs, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
```

`aruba.URI(s)` returns a lightweight `Ref` that the SDK uses to extract ancestor IDs from the URI path segments. Any valid resource URI works — the SDK parses it internally.

---

## 6. Tear Down (Reverse Order)

Delete children before parents. The Aruba Cloud API returns **HTTP 400** when you try to delete a parent that still has live or still-deleting children — not 409/422. The safe pattern is to issue each child delete, then wait until the resource is fully gone before moving up the dependency chain.

`WaitUntilGone` blocks until the resource no longer exists — that is, until its `Get` returns HTTP 404. Call it on the hydrated wrapper you already hold (`Delete` itself returns only an `error`, no wrapper):

```go
// subnet → VPC → project
if err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet); err != nil {
    log.Printf("Delete subnet: %v", err)
} else if err := subnet.WaitUntilGone(ctx); err != nil {
    log.Printf("Subnet not gone: %v", err)
}

if err := arubaClient.FromNetwork().VPCs().Delete(ctx, vpc); err != nil {
    log.Printf("Delete VPC: %v", err)
} else if err := vpc.WaitUntilGone(ctx); err != nil {
    log.Printf("VPC not gone: %v", err)
}

if err := arubaClient.FromProject().Delete(ctx, proj); err != nil {
    log.Printf("Delete project: %v", err)
}
```

`WaitUntilGone` accepts the same `WaitOption`s as `WaitUntilReady` (`WithRetries`, `WithBaseDelay`, `WithTimeout`) and is available on every resource wrapper that supports polling. `Project` has no polling — it is deleted last, with no child left to wait on.

`Delete` accepts any `aruba.Ref` — pass the hydrated wrapper directly or `aruba.URI(…)` if you only have the path.

For a full stack teardown sequence (Security Rules → Security Groups → Subnets → VPC → Cloud Server → Block Storage → Project) see the [Full Example](#full-example) below.

---

## 7. Wait for Readiness

Most cloud operations — Create, Update, scale operations — are **asynchronous**: the HTTP call returns quickly, but the resource keeps transitioning through states (`Creating` → `Active`, `Updating` → `Active`) for seconds to minutes in the background.

The `WaitUntilReady` method on any resource wrapper that embeds `statusMixin` blocks until the resource reaches the `"Active"` state (or returns an error on terminal failure):

```go
if err := vpc.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

Three `WaitOption`s let you override the defaults (60 retries × 10 s base delay × 600 s hard ceiling):

```go
if err := vpc.WaitUntilReady(ctx,
    aruba.WithRetries(30),              // max polling iterations (default: 60)
    aruba.WithBaseDelay(5*time.Second), // fixed delay between polls (default: 10s)
    aruba.WithTimeout(3*time.Minute),   // hard deadline (default: 600s)
); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

For `WaitUntilStates(ctx, []types.State{...}, opts...)` (any target states, not just `"Active"`), status accessors (`State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`), and the low-level `pkg/async.WaitFor` future for concurrent polling, see the [Async / Await](./async) guide.

---

## Caveats

### Setter errors are deferred

Builder setters never return an error — they record it in the wrapper. The error is returned by the first `Create` or `Update` call, or you can check eagerly:

```go
rule := aruba.NewSecurityRule().
    InSecurityGroup(sg).
    TargetingCIDR("0.0.0.0/0").
    TargetingSecurityGroup(otherSG) // conflicting — recorded as error

if err := rule.Err(); err != nil {
    log.Fatalf("Bad rule config: %v", err)
}
```

> **Caveat**: `TargetingCIDR` and `TargetingSecurityGroup` are mutually exclusive. Setting both records a setter-time error that surfaces on `Create`.

### `WaitUntilReady` requires a hydrated wrapper

Calling `WaitUntilReady` on a wrapper you constructed manually (without `Create`/`Get`/`Update`/`List`) returns:

```
WaitUntilStates: refresh callback not set; resource must be produced by an adapter (Create/Get/Update/List) to support polling
```

Always use the wrapper returned by the API call, not the request builder.

### Typed HTTP errors

Non-2xx API responses are returned as `*aruba.HTTPError`. Use `errors.As` to inspect them:

```go
vpc, err = arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        log.Printf("API error %d: %s", httpErr.StatusCode, httpErr.Error())
    } else {
        log.Fatalf("Network error: %v", err)
    }
}
```

See [Response Handling](./response-handling) for the full error handling guide.

---

## Full Example

The `examples/all-resources/` directory in the repository contains a runnable end-to-end example demonstrating all resources:

```bash
go run ./examples/all-resources/ -mode=create -clientID=… -clientSecret=…
go run ./examples/all-resources/ -mode=update -clientID=… -clientSecret=… -projectID=…
go run ./examples/all-resources/ -mode=delete -clientID=… -clientSecret=… -projectID=…

# Add -debug for verbose SDK logging:
go run ./examples/all-resources/ -mode=create -clientID=… -clientSecret=… -debug
```
