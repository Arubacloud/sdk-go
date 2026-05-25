# Working at Low Level

## Why this page exists

The SDK is designed around a **single-import principle**: importing only `github.com/Arubacloud/sdk-go/pkg/aruba` covers ~99.9% of real-world use cases.

The wrapper surface — typed accessors, `WaitUntil*` helpers, `RawJSON()` / `RawYAML()` — is designed so you rarely need to reach into the underlying wire structs. For the 0.1% of cases where you do, this page collects every escape hatch that requires a second import.

> These patterns are **intentional and supported** — they are escape hatches, not workarounds. When you hit a case not covered by the wrapper API, reach for these rather than opening a feature request.

---

## Accessing non-promoted wire fields with `Raw()`

Every wrapper exposes a `Raw()` method that returns the underlying `*types.XxxResponse` struct. Use it when you need a field that hasn't been promoted to the wrapper surface:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

raw := vpc.Raw()                          // *types.VPCResponse
fmt.Println(raw.Properties.IsDefault)     // field not promoted to the wrapper
```

For lists, `Raw()` returns `any` and you type-assert to the concrete list type:

```go
vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil { /* … */ }

raw := vpcList.Raw().(*types.VPCList)     // requires pkg/types import
fmt.Println("server total:", raw.Total)   // same as vpcList.Total() — shown for illustration
fmt.Println("self link:", raw.Self)
```

> **Prefer wrapper accessors for serialisation.** If your goal is JSON or YAML output, use `vpc.RawJSON()` / `vpc.RawYAML()` (or the `List[T]` equivalents) — no `pkg/types` import required.
>
> ```go
> fmt.Println(string(vpcList.RawJSON()))  // JSON without pkg/types
> fmt.Println(string(vpcList.RawYAML()))  // YAML without pkg/types
> ```

---

## Inspecting structured validation errors

`*aruba.HTTPError` is the error type for all 4xx/5xx responses. Its `ErrResp` field is a `*types.ErrorResponse` and holds structured RFC 7807 details — including a `[]types.ValidationError` slice for field-level 400 errors:

```go
import (
    "errors"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

_, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) && httpErr.ErrResp != nil {
        fmt.Printf("title:  %s\n", derefStr(httpErr.ErrResp.Title))
        fmt.Printf("detail: %s\n", derefStr(httpErr.ErrResp.Detail))

        // Field-level validation errors — require types.ValidationError
        for _, ve := range httpErr.ErrResp.Errors {
            fmt.Printf("  field %s: %s\n", ve.Field, ve.Message)
        }

        // TraceID for support requests
        fmt.Printf("trace-id: %s\n", derefStr(httpErr.ErrResp.TraceID))
    }
}
```

### `MetadataValidationError`

A `*types.MetadataValidationError` is returned (alongside a non-nil wrapper) when an API response is missing required metadata fields (`id` or `uri`). Use `errors.As` to detect it:

```go
import (
    "errors"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

result, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var mvErr *types.MetadataValidationError
    if errors.As(err, &mvErr) {
        // result is non-nil and partially hydrated; mvErr lists the missing fields
        fmt.Printf("metadata incomplete: %v\n", mvErr)
        fmt.Printf("ID so far: %s\n", result.ID())
    }
}
```

---

## Iterating `LinkedResources()`

Every resource wrapper exposes `LinkedResources()` which returns `[]types.LinkedResource`. Each entry has a `URI string` and a `StrictCorrelation bool`:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

for _, lr := range vpc.LinkedResources() {
    fmt.Println("linked URI:", lr.URI)
    if lr.StrictCorrelation {
        fmt.Println("  → strict correlation (lifecycle-linked)")
    }
}
```

> If you only need the linked URI strings — for example, to pass them back to another SDK call as `aruba.URI(lr.URI)` — you don't need `types.LinkedResource` at all:
>
> ```go
> for _, lr := range vpc.LinkedResources() {
>     ref := aruba.URI(lr.URI)   // no pkg/types import
>     _ = ref
> }
> ```

---

## Inspecting request bodies before they are sent

Every wrapper exposes `RawRequest()` which returns the wire-level request struct (`*types.XxxRequest`). This is useful for debugging or for feeding the request into another tool:

```go
import (
    "encoding/json"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc := aruba.NewVPC().
    Named("my-vpc").
    InProject(proj).
    InRegion(aruba.RegionITBGBergamo).
    AsDefault()

req := vpc.RawRequest()               // types.VPCRequest — requires pkg/types import
b, _ := json.MarshalIndent(req, "", "  ")
fmt.Println(string(b))
```

---

## Background polling with `pkg/async` {#background-polling-with-pkgasync}

`WaitUntilReady`, `WaitUntilActive`, and `WaitUntilStates` block the calling goroutine. If you need to **start multiple waits concurrently**, or **poll an arbitrary condition** (not just a resource state), use the lower-level `pkg/async` package directly.

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
        var state types.State
        if resp.Data.Properties != nil && resp.Data.Properties.Status != nil &&
            resp.Data.Properties.Status.State != nil {
            state = *resp.Data.Properties.Status.State
        }
        return state == types.StateActive, nil
    },
)

futureVPC2 := async.DefaultWaitFor(ctx, /* same pattern for vpc2 */)

// Block for both results
resp1, err1 := futureVPC1.Await(ctx)
resp2, err2 := futureVPC2.Await(ctx)
```

`DefaultWaitFor` uses the package defaults: `DefaultRetries=60`, `DefaultBaseDelay=10s`, `DefaultTimeout=600s`. Use `async.WaitFor(ctx, retries, baseDelay, timeout, call, check)` to override.

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

## What does NOT require `pkg/types`

The following are all available via a single `pkg/aruba` import — no second import needed:

| What you need | `pkg/aruba` surface |
|---|---|
| State constants (`StateActive`, `StateStopped`, …) | `aruba.StateActive`, `aruba.StateStopped`, … |
| Region / zone constants | `aruba.RegionITBGBergamo`, `aruba.ZoneITBG1`, … |
| Billing period | `aruba.BillingPeriodHour`, `aruba.BillingPeriodMonth`, … |
| All compute, storage, network, security enums | See `pkg/aruba/aliases.go` |
| Wait for state transitions | `wrapper.WaitUntilReady(ctx)`, `WaitUntilActive`, `WaitUntilStates` |
| Serialise a response to JSON / YAML | `wrapper.RawJSON()`, `wrapper.RawYAML()` |
| HTTP envelope introspection | `wrapper.StatusCode()`, `.Headers()`, `.RawHTTP()`, `.RawError()` |
| Pagination | `list.Total()`, `.HasNext()`, `.Next(ctx)`, `.All(ctx, yield)` |
| HTTP error details | `*aruba.HTTPError` — `StatusCode`, `ErrResp.Title`, `ErrResp.Detail`, `ErrResp.TraceID` |

---

## See Also

- [Response Handling](./response-handling) — typed HTTP errors, envelope accessors, `RawJSON`/`RawYAML`
- [Async / Await](./async) — `WaitUntilReady`, `WaitUntilStates`, polling options
- [API Walkthrough](./walkthrough) — full Create + poll + Update + Delete lifecycle example

---

```go
// Helper used in examples above
func derefStr(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
```
