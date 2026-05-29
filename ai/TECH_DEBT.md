# TECH_DEBT.md — Technical Debt & Refactoring Backlog

Issues are grouped by severity. Address Critical items before new features ship; High items before any public release.

**Effort scale:** XS < 30 min · S < 2 h · M half-day · L 1–2 days · XL 3+ days

**Impact scale:** Critical — broken in production · High — leak/race/major debt · Medium — edge-case failures or observability · Low — cosmetic or defensive

## Prioritization matrix

| ID | Summary | Severity | Effort | Impact |
|---|---|---|---|---|
| [TD-003](#td-003) | `lastUsage` race under RLock | Critical | S | High |
| [TD-009](#td-009) | Caller headers override `Content-Type` | High | XS | Medium |
| [TD-010](#td-010) | 2 000+ lines duplicated response parsing | High | L | High |
| [TD-012](#td-012) | Expired token injected after failed refresh | Medium | S | High |
| [TD-014](#td-014) | `ParseResponseBody` panics on nil response | Medium | XS | Medium |
| [TD-015](#td-015) | `DefaultWaitFor` timeout too short | Medium | XS | Medium |
| [TD-016](#td-016) | No structured logging | Medium | L | Medium |
| [TD-017](#td-017) | `WARN` writes to stdout | Medium | XS | Low |
| [TD-023](#td-023) | `KaasSecurityGroupPropertiesResponse` mixed casing | Low | XS | Low |
| [TD-024](#td-024) | Stale TODO comments in `pkg/aruba/aruba.go` | Low | XS | Low |
| [TD-025](#td-025) | `pkg/util/middleware` package name is `restclient` | Low | S | Low |
| [TD-026](#td-026) | No lint rule enforces Request/Response/Common naming | Low | S | Medium |
| [TD-027](#td-027) | `Dto` suffix on enum-string constants | Low | XS | Low |
| [TD-020](#td-020) | Test coverage gaps | Low | XL | High |

### Recommended execution order

**Wave 1 — Quick Wins** (XS effort, ship same PR): TD-009, TD-014, TD-015, TD-017, TD-023, TD-024, TD-027

**Wave 2 — High-value focused fixes** (S effort, High+ impact): TD-003, TD-012, TD-025, TD-026

**Wave 3 — Large refactors** (L/XL, plan separately): TD-010, TD-016, TD-020

---

## Critical

### TD-003 · Race condition: `lastUsage` written under read lock in Multitenant — [#113](https://github.com/Arubacloud/sdk-go/issues/113)
`pkg/multitenant/multitenant.go` — `Get()`, `MustGet()`, and `GetOrNil()` all hold `RLock` while writing `e.lastUsage = time.Now()`. Mutating a struct field through a map value under a read lock is a data race detected by the Go race detector.

**Fix:** Upgrade to `Lock()` for the update, or store `lastUsage` as an `atomic.Int64` (Unix nanoseconds) so it can be updated without a write lock.

**Effort:** S — change 3 methods in 1 file; atomic variant needs a small struct change.

**Impact:** High — data race under concurrent multitenant use; triggers the race detector and can cause unpredictable behaviour in production.

---

## High

### TD-009 · Caller headers can override SDK-controlled `Content-Type` — [#119](https://github.com/Arubacloud/sdk-go/issues/119)
`internal/restclient/client.go` — the code sets `Content-Type: application/json` then overwrites headers with caller-supplied values using `req.Header.Set(k, v)`. A caller can silently override `Content-Type` (and in principle `Authorization`), breaking server-side request parsing.

**Fix:** Apply caller headers before SDK-controlled headers, or explicitly protect `Content-Type` and `Authorization` from being overridden.

**Effort:** XS — reorder 3 lines in 1 function.

**Impact:** Medium — silent request corruption if caller passes a `Content-Type` header; unlikely but a correctness guarantee violation.

---

### TD-010 · Create/Update methods duplicate `ParseResponseBody` logic across ~15 client files — [#120](https://github.com/Arubacloud/sdk-go/issues/120)
`internal/clients/compute/cloudserver.go`, `internal/clients/compute/keypair.go`, `internal/clients/network/vpc.go`, `internal/clients/network/security-group.go`, and ~11 more — List/Get/Delete methods call `types.ParseResponseBody[T]()`. Create/Update methods in the same files manually re-implement the same logic: marshal body → `DoRequest` → `io.ReadAll` → construct `Response[T]` wrapper → unmarshal success/error branches. This is ~40 lines duplicated per mutating method, totalling 2 000+ lines across the codebase. Any bug fix to the parsing logic must be applied in all locations.

**Fix:** Add a generic helper to `restclient.Client`:
```go
func DoAndParse[T any](c *Client, ctx context.Context, method, path string,
    body io.Reader, queryParams, headers map[string]string) (*types.Response[T], error)
```
Replace all manual implementations with calls to this helper.

**Effort:** L — design the helper, then update ~30 methods across ~15 files mechanically.

**Impact:** High — eliminates 2 000+ lines of duplication; any bug fix to parsing logic will only need to be applied once.

---

## Medium

### TD-012 · Token manager mishandles post-refresh token fetch in the else-branch — [#122](https://github.com/Arubacloud/sdk-go/issues/122)
`internal/impl/auth/tokenmanager/standard/standard.go` — in the write-lock branch, when the ticket has already changed (another goroutine refreshed), the code re-fetches from the repository. Two failure modes exist:

1. **Nil pointer panic:** If the second `FetchToken()` returns `auth.ErrTokenNotFound`, the error is filtered out by `!errors.Is(err, auth.ErrTokenNotFound)` and not returned. Execution falls through to `token.AccessToken` at the header-injection step with `token == nil`, causing a panic.
2. **Silent expired token injection:** If `FetchToken()` returns a token without error but the token is expired, it is injected into the `Authorization` header without expiration validation.

**Fix:** (1) Return an error (or retry) when `ErrTokenNotFound` is received in the else-branch rather than falling through. (2) Add an expiration check before injecting the token.

**Effort:** S — add nil/expiry checks in 1 function.

**Impact:** High — the nil pointer panic causes an unrecoverable crash under concurrent token refresh with a temporarily unavailable token store; the expired-token path produces silent 401s.

---

### TD-014 · `ParseResponseBody` panics on nil `httpResp` — [#124](https://github.com/Arubacloud/sdk-go/issues/124)
`pkg/types/utils.go` — the function calls `io.ReadAll(httpResp.Body)` without a nil guard. If `DoRequest` returns a nil response alongside a non-nil error and the caller forgets to check the error, the function panics.

**Fix:** Add `if httpResp == nil { return nil, fmt.Errorf("http response is nil") }` at the top of the function.

**Effort:** XS — add 1 nil guard (3 lines).

**Impact:** Medium — defensive fix; prevents a panic from a defensive coding mistake by a future caller.

---

### TD-015 · `DefaultWaitFor` timeout of 60 s is too short for cloud operations — [#125](https://github.com/Arubacloud/sdk-go/issues/125)
`pkg/async/async_client.go` — constants are `DefaultRetries=60`, `DefaultBaseDelay=10s`, `DefaultTimeout=600s`. The original issue filed against a 60 s default has been resolved in code (timeout is now 600 s). However, the constant names (`DefaultRetries=60`) could be confused with a 60 s timeout. The documentation should explicitly state that timeout is 600 s (10 minutes) and retries is the poll-attempt count.

**Status:** Resolved in code; remaining work is documentation clarity.

**Fix:** Add inline godoc to the three constants in `pkg/async/async_client.go` clarifying units.

**Effort:** XS — add 3 short doc comments.

**Impact:** Medium — prevents callers from misreading `DefaultRetries=60` as "60 second timeout".

---

### TD-016 · Logger interface supports only printf-style formatting — no structured logging — [#126](https://github.com/Arubacloud/sdk-go/issues/126)
`internal/ports/logger/logger.go` — all log calls use `%s`/`%v` format strings. In production environments using log aggregators (ELK, Loki, Cloud Logging), querying by structured fields (project ID, resource ID, trace ID) requires string parsing. The SDK's `ErrorResponse` already carries a `TraceID` field that is never emitted as a structured log field.

**Fix:** Add optional structured variants to the logger interface, or adopt `slog` (stdlib since Go 1.21) as the native logger implementation.

**Effort:** L — redesign the logger interface, update the native implementation, and update all ~50 call sites.

**Impact:** Medium — major observability improvement for production use; no functional bug risk.

---

### TD-017 · `WARN` log level writes to `os.Stdout` instead of `os.Stderr` — [#127](https://github.com/Arubacloud/sdk-go/issues/127)
`internal/impl/logger/native/logger.go` — `DEBUG` and `INFO` correctly go to `os.Stdout`; `ERROR` correctly goes to `os.Stderr`. `WARN` also goes to `os.Stdout`, which breaks shell pipelines and container log routers that separate informational output from diagnostic output.

**Fix:** Change `WARN` to write to `os.Stderr`.

**Effort:** XS — change 1 line.

**Impact:** Low — minor log routing correction; no functional impact.

---

## Low

### TD-023 · `KaasSecurityGroupPropertiesResponse` uses mixed PascalCase (`Kaas` vs `KaaS`) — (new)
`pkg/types/container.kaas.go:240` — the response-side struct is named `KaasSecurityGroupPropertiesResponse` while the symmetrical request-side struct at line 91 is named `KaaSSecurityGroupPropertiesRequest` (correct). The Go type name `Kaas` (lowercase second 'a') differs from the acronym's canonical casing `KaaS`. The wire JSON tag is unaffected, but the API surface is inconsistent.

**Fix:** Rename `KaasSecurityGroupPropertiesResponse` → `KaaSSecurityGroupPropertiesResponse` in `pkg/types/container.kaas.go` and its one usage in `resource_kaas_test.go`.

**Effort:** XS — rename 1 type and update 1 test reference.

**Impact:** Low — cosmetic inconsistency; does not affect serialisation.

---

### TD-024 · Stale TODO comments in `pkg/aruba/aruba.go` — (new)
`pkg/aruba/aruba.go` contains four TODO comments describing planned variations of `NewClient` and options loading from file/URL that have not been implemented and are not on any active roadmap:
- `// TODO: Two variations of NewClient`
- `// TODO: DefaultOptions() function`
- `// TODO: LoadOptionsFromPath(path path.Path)`
- `// TODO: LoadOptionsFromURL(url net.URL)`

`DefaultOptions` already exists (in `options.go`). The other two (file/URL loading) may or may not be planned. These comments add noise to the public package entrypoint file.

**Fix:** Remove stale TODOs; if file/URL loading is genuinely planned, open a GitHub issue and remove the in-code notes.

**Effort:** XS — delete 4 comment blocks.

**Impact:** Low — cosmetic; improves readability of the package entrypoint.

---

### TD-025 · `pkg/util/middleware` package is named `restclient`, not `middleware` — (new)
`pkg/util/middleware/middleware.go` declares `package restclient`. The directory path is `pkg/util/middleware/` but the Go package name is `restclient`. Callers import it as `github.com/Arubacloud/sdk-go/pkg/util/middleware` but refer to it as `restclient.WithCustomHeaders(...)`. This diverges from Go's convention that the directory name and package name match and is flagged by `golangci-lint` (gopkg convention linters). The package itself also has an internal TODO: "review the placement of this file".

**Fix:** Either (a) rename the package declaration to `middleware`, or (b) move the file to `internal/` since it wraps an internal interceptor type and has no external callers yet.

**Effort:** S — rename + update any callers.

**Impact:** Low — cosmetic confusion; the package has no known external callers in the repo.

---

### TD-026 · No lint rule enforces the `Request`/`Response`/`Common` naming convention — (new)
The `forbidigo` linter slot exists in `.golangci.yml` but is not configured. Without a lint rule, a new struct added to `pkg/types/` with an old-style suffix (`*Result`, `*List`) will silently violate the convention until caught in code review.

**Fix:** Add a `forbidigo` rule (or a custom `revive` rule) that fails on exported struct names in `pkg/types/` matching `[A-Z][A-Za-z]+List$` or `[A-Z][A-Za-z]+Result$`.

**Effort:** S — configure linter, adjust any existing violations first.

**Impact:** Medium — prevents naming-convention regression; ensures the refactor done in the Unreleased block is durable.

---

### TD-027 · `Dto` suffix on endpoint-type enum strings — (new)
`pkg/types/network.load-balancer.go` (and possibly related files) contain constants like `EndpointTypeDto` and `DeactiveReasonDto` — suffixes inherited from generated code. These are wire-string enum values; the `Dto` suffix is inconsistent with the rest of the enum naming in `pkg/types/` and leaks an internal detail into the public API surface.

**Fix:** Rename affected constants (confirm full list via `git grep -n 'Dto'`). These are enum-string values, not struct types, so they are out of scope for the Request/Response/Common rename but should be cleaned up separately.

**Effort:** XS per constant — alias old names during a deprecation window or do a hard rename.

**Impact:** Low — cosmetic; no functional impact.

---

### TD-020 · Test coverage limited to happy path and empty-ID validation — [#130](https://github.com/Arubacloud/sdk-go/issues/130)
All `internal/clients/*/_test.go` files — existing tests cover successful responses and empty project/resource IDs. Missing: HTTP 4xx/5xx error responses, malformed JSON bodies, network-level errors, `nil` params handling, and request body marshaling failures.

**Fix:** Add table-driven subtests for each error scenario in every client test file using `httptest.NewServer` to return controlled error responses.

**Effort:** XL — ~15 client files × ~5 error scenarios each; requires design of shared test helpers.

**Impact:** High — major regression safety net; issues like TD-011 and TD-014 would have been caught by these tests.

---

## Resolved (historical summary)

| ID | Summary | Version |
|---|---|---|
| TD-001 | Parameter order swap in file token repository constructor | v0.2.x |
| TD-002 | Static token silently ignored in memory repository | v0.2.x |
| TD-004 | `Bind()`/`Intercept()` concurrency — closed as working-as-intended | — |
| TD-005 | Typo `buildDetebaseClient` | v0.2.x |
| TD-006 | Goroutine leak in `WaitFor` — closed as invalid | — |
| TD-007 | Variable shadowing in `WaitFor` nil-call path | v0.2.x |
| TD-008 | Polling loop: sleep before first attempt, discard final state | v0.2.x |
| TD-011 | Silent failure parsing error response body | v0.2.x |
| TD-013 | `SaveToken` increments ticket before confirming persistent write | v0.2.x |
| TD-018 | Injected concrete dependencies not nil-checked in constructors | v0.3.0 |
| TD-019 | Missing compile-time interface satisfaction checks | v0.3.0 |
| TD-021 | Create responses do not validate metadata ID/URI/Name | v1.0.0 |
| TD-022 | Schedule Job `typology` field root cause — `typology` field removed; Terraform confirms shape is correct | v1.0.0 |
