# Response Handling Guide

## Overview

The SDK wrapper layer handles response parsing and error surfacing automatically. Every CRUD method returns either a populated wrapper (on success) or an error. You rarely need to inspect the raw HTTP envelope — but the tools to do so are all there when you need them.

## Basic Pattern

Every wrapper method returns `(wrapper, error)`. The error is non-nil for both network failures and API-level errors (4xx / 5xx).

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx,
    aruba.URI("/projects/<projectID>/providers/Aruba.Network/vpcs/<vpcID>"),
)
if err != nil {
    log.Fatalf("Get VPC failed: %v", err)
}
fmt.Printf("VPC: %s (state: %s)\n", vpc.Name(), vpc.State())
```

## Typed HTTP Errors

When the API returns a 4xx or 5xx response, the SDK wraps it in `*aruba.HTTPError`. Use `errors.As` to inspect the status code and structured error body:

```go
import "errors"

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("API error %d: %s\n", httpErr.StatusCode, httpErr.Error())
        if httpErr.ErrResp != nil {
            fmt.Printf("  title:  %s\n", derefStr(httpErr.ErrResp.Title))
            fmt.Printf("  detail: %s\n", derefStr(httpErr.ErrResp.Detail))
            for _, ve := range httpErr.ErrResp.Errors {
                fmt.Printf("  field %s: %s\n", ve.Field, ve.Message)
            }
        }
    } else {
        // Network error, context timeout, etc.
        log.Fatalf("Request failed: %v", err)
    }
}
```

## Complete Error Handling Pattern

```go
proj, err := arubaClient.FromProject().Get(ctx, ref)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        switch httpErr.StatusCode {
        case 404:
            return fmt.Errorf("project not found")
        case 400:
            return fmt.Errorf("bad request: %s", derefStr(httpErr.ErrResp.Detail))
        default:
            return fmt.Errorf("API error (%d): %s", httpErr.StatusCode, httpErr.Error())
        }
    }
    return fmt.Errorf("request failed: %w", err)
}
// proj is populated — use it directly
fmt.Printf("Project: %s (tags: %v)\n", proj.Name(), proj.Tags())
```

## HTTP Envelope Accessors

Every wrapper produced by a Create / Get / Update / List call exposes its raw HTTP envelope:

```go
// After any CRUD call:
proj, err := arubaClient.FromProject().Create(ctx, p)
// …

fmt.Println("Status:", proj.StatusCode())
fmt.Println("Headers:", proj.Headers())
rawResp, rawBody := proj.RawHTTP()
fmt.Println("Raw body:", string(rawBody))
fmt.Println("HTTP status:", rawResp.StatusCode)
fmt.Println("Error body (if any):", proj.RawError())
```

## Accessing the Raw Wire Response

Every wrapper has a `Raw()` method that returns the underlying typed response struct from `pkg/types`. This is useful for fields not yet promoted to the wrapper surface:

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

raw := vpc.Raw()                         // underlying wire struct
fmt.Println(raw.Properties.IsDefault)    // field not on the wrapper
```

## List Responses

`List[T]` exposes the same introspection surface as single-resource wrappers:

```go
vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil { /* … */ }

// Wire payload — JSON-marshalable; contains pagination links and values.
raw := vpcList.Raw().(*types.VPCList)
fmt.Println("server total:", raw.Total)

// HTTP envelope — same accessors as single-resource wrappers.
fmt.Println("status:", vpcList.StatusCode())
fmt.Println("trace-id:", vpcList.Headers().Get("X-Trace-Id"))
_, body := vpcList.RawHTTP()
fmt.Println("raw body bytes:", len(body))
```

`Raw()` returns only the JSON-safe payload (`*types.XxxList`). The HTTP envelope is exposed via separate accessors because `*http.Response` is not JSON-serialisable.

## Setter-Time Errors

Fluent builder setters never return errors — instead they record them internally. The error surfaces the first time you call `Create` or `Update`. You can also check eagerly:

```go
rule := aruba.NewSecurityRule().
    IntoSecurityGroup(sg).
    WithTargetCIDR("0.0.0.0/0").
    WithTargetSecurityGroup(otherSG) // conflicting setter — records an error

if err := rule.Err(); err != nil {
    log.Fatalf("Bad rule configuration: %v", err)
}
```

## Best Practices

1. **Always check `err` first** — it covers both network failures and API errors.
2. **Use `errors.As(err, &httpErr)`** to get structured error details on 4xx/5xx.
3. **Check `httpErr.ErrResp.Errors`** for field-level validation messages on 400.
4. **Use `httpErr.ErrResp.TraceID`** when filing a support request.
5. **Use `.Raw()`** sparingly — prefer the typed wrapper accessors.
6. **Check `wrapper.Err()` before Create/Update** when the builder chain is long.

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
