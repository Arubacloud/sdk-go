---
sidebar_position: 3
---

# Async / Await

Most Aruba Cloud API operations are **asynchronous**: the HTTP call returns quickly (with a `201 Created` or `200 OK`), but the resource keeps transitioning through states in the background — `Creating` → `Active`, or `Updating` → `Active`, or `Deleting` → gone — for seconds to several minutes.

The SDK exposes three layers for dealing with this:

| Layer | When to use |
|-------|-------------|
| `WaitUntilActive(ctx)` | 95% of cases — block until the resource is ready |
| `WaitUntilState(ctx, target)` | Wait for any named state (e.g. `"Stopped"`) |
| `pkg/async.WaitFor` + `AsyncClient.Await` | Advanced — start polling in a background goroutine, do other work, collect the result later |

---

## `WaitUntilActive`

After any `Create`, `Update`, or `Get`, call `WaitUntilActive` on the returned wrapper to block until the resource reaches the `"Active"` state:

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}

if err := vpc.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

`WaitUntilActive` polls the API repeatedly with a fixed delay. When the resource enters a known **error terminal state** (e.g. `"Error"`, `"Failed"`), it returns immediately with a descriptive error rather than exhausting all retries.

See the [API Walkthrough](./walkthrough) for full Create + poll + Update + Delete examples.

### Tuning poll behaviour

Three call options let you override the defaults:

```go
if err := vpc.WaitUntilActive(ctx,
    aruba.WithRetries(30),              // max polling iterations (default: 60)
    aruba.WithBaseDelay(5*time.Second), // fixed delay between polls (default: 10s)
    aruba.WithTimeout(3*time.Minute),   // hard deadline (default: 600s)
); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

The effective ceiling is `min(retries × baseDelay, timeout)`. At the defaults that is `min(60×10s, 600s) = 600s`.

---

## `WaitUntilState`

Use `WaitUntilState` when you need to wait for a state other than `"Active"`:

```go
// Wait for a Cloud Server to fully stop after PowerOff
if err := cs.WaitUntilState(ctx, "Stopped"); err != nil {
    log.Fatalf("Cloud Server did not stop: %v", err)
}
```

```go
// Wait until a DBaaS instance finishes an in-progress update
if err := db.WaitUntilState(ctx, "Active",
    aruba.WithRetries(120),
    aruba.WithBaseDelay(15*time.Second),
); err != nil {
    log.Fatalf("DBaaS did not return to Active after update: %v", err)
}
```

The same error-terminal-state early exit applies: if the resource reaches `"Error"` or `"Failed"` while you are waiting for `"Stopped"`, the call returns immediately with an error that names both the actual state and the target state.

---

## Status Accessors

Every wrapper that supports polling also exposes fine-grained status accessors. You can read these at any time after a `Create`, `Get`, `Update`, or `List` call:

| Method | Returns | Typical use |
|--------|---------|-------------|
| `State()` | `string` — current state | Logging, conditional branching |
| `PreviousState()` | `string` — state before the last transition | Post-mortem after a failed wait |
| `FailureReason()` | `string` — server-supplied error text | Surface to end user / log alert |
| `IsDisabled()` | `bool` | Gate operations when server disables a resource |
| `DisableReasons()` | `[]string` | Explain why a resource is disabled |

A common pattern — call `WaitUntilActive`, and if it fails, attach the server's failure reason to the error:

```go
if err := vpc.WaitUntilActive(ctx); err != nil {
    reason := vpc.FailureReason()
    if reason != "" {
        log.Fatalf("VPC failed: %v (reason: %s)", err, reason)
    }
    log.Fatalf("VPC polling failed: %v", err)
}
```

---

## Resources That Support Polling

The following resource wrappers embed the polling mixin and support `WaitUntilActive`, `WaitUntilState`, and the status accessors:

- **Compute**: `CloudServer`
- **Container**: `KaaS`, `ContainerRegistry`
- **Database**: `DBaaS`
- **Network**: `VPC`, `Subnet`, `SecurityGroup`, `SecurityRule`, `ElasticIP`
- **Security**: `KMS`, `Kmip`
- **Storage**: `BlockStorage`, `Snapshot`

> **Project does not support polling.** It is synchronously ready immediately after `Create` returns — no `WaitUntilActive` call is needed or available.

---

## Caveats

### Hydrated wrapper required

`WaitUntilActive` and `WaitUntilState` only work on wrappers that were **returned by an adapter call** (`Create`, `Get`, `Update`, or `List`). Calling either method on a freshly-built request builder returns:

```
WaitUntilState: refresh callback not set; resource must be produced by an adapter (Create/Get/Update/List) to support polling
```

Always use the wrapper returned by the API call:

```go
// Correct — vpc was returned by Create
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, myVPC)
vpc.WaitUntilActive(ctx)

// Wrong — myVPC is a request builder, not an adapter response
myVPC := aruba.NewVPC().WithName("x")
myVPC.WaitUntilActive(ctx) // returns "refresh callback not set"
```

### Constant poll cadence

Polling uses a **fixed delay** (no exponential backoff). If you are hitting API rate limits, increase `WithBaseDelay` rather than expecting the SDK to back off automatically.

### Context cancellation

All polling respects the `ctx` deadline and cancellation. If the context expires mid-poll the call returns `ctx.Err()` (typically `context.DeadlineExceeded` or `context.Canceled`).

---

## Advanced: Background Polling with `pkg/async`

`WaitUntilActive` and `WaitUntilState` block the calling goroutine. If you need to **start multiple waits concurrently**, or **poll an arbitrary condition** (not just a resource state), use the lower-level `pkg/async` package directly.

`pkg/async` is a public package — import it alongside `pkg/aruba`:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/async"
    "github.com/Arubacloud/sdk-go/pkg/types"
)
```

### `WaitFor` — start a background future

`async.WaitFor` launches a polling goroutine immediately and returns an `*async.AsyncClient[T]` (a future). You call `.Await(ctx)` later to block for the result:

```go
// Start polling VPC1 and VPC2 concurrently
futureVPC1 := async.DefaultWaitFor(ctx,
    func(ctx context.Context) (*types.Response[types.VPCResponse], error) {
        return arubaClient.FromNetwork().VPCs().Get(ctx, vpc1)
    },
    func(resp *types.Response[types.VPCResponse]) (bool, error) {
        if resp == nil || resp.Data == nil {
            return false, nil
        }
        state := ""
        if resp.Data.Properties != nil && resp.Data.Properties.Status != nil &&
            resp.Data.Properties.Status.State != nil {
            state = *resp.Data.Properties.Status.State
        }
        return state == "Active", nil
    },
)

futureVPC2 := async.DefaultWaitFor(ctx, /* same pattern for vpc2 */)

// Block for both results
resp1, err1 := futureVPC1.Await(ctx)
resp2, err2 := futureVPC2.Await(ctx)
```

`DefaultWaitFor` uses the same defaults as `WaitUntilActive`: 60 retries, 10s delay, 600s timeout. Use `async.WaitFor(ctx, retries, baseDelay, timeout, call, check)` to override.

### `WaitFor` signature

```go
func WaitFor[T any](
    ctx         context.Context,
    retries     int,
    baseDelay   time.Duration,
    timeout     time.Duration,
    call        func(ctx context.Context) (*types.Response[T], error),
    check       func(*types.Response[T]) (bool, error),
) *AsyncClient[T]
```

- `call` — the polling function, called once per iteration.
- `check` — returns `(true, nil)` to signal success, `(true, error)` to signal terminal failure, `(false, nil)` to keep polling.
- If `check` is `nil`, any non-nil `response.Data` is treated as success.

### `AsyncClient.Await`

```go
func (f *AsyncClient[T]) Await(ctx context.Context) (*types.Response[T], error)
```

Blocks until the background goroutine sends its result or `ctx` is cancelled. Subsequent calls return the **cached** result immediately — safe to call multiple times.

> `pkg/async` works directly with the `pkg/types` wire structs. This is the only layer of the SDK where you'll interact with `types.Response[T]` and `types.*Response` types directly.

---

## See Also

- [API Walkthrough](./walkthrough) — full Create + `WaitUntilActive` + Update + Delete lifecycle example
- [Response Handling](./response-handling) — how `*aruba.HTTPError` propagates through `WaitUntilActive` when the API returns 4xx/5xx
