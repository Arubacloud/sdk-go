# TECH_DEBT.md — Technical Debt & Refactoring Backlog

Issues are grouped by severity. Address Critical items before new features ship; High items before any public release.

**Effort scale:** XS < 30 min · S < 2 h · M half-day · L 1–2 days · XL 3+ days

**Impact scale:** Critical — broken in production · High — leak/race/major debt · Medium — edge-case failures or observability · Low — cosmetic or defensive

## Prioritization matrix

| ID | Summary | Severity | Effort | Impact |
|---|---|---|---|---|
| [TD-001](#td-001) | File token repo param order swap | Critical | XS | Critical |
| [TD-002](#td-002) | Static token silently ignored | Critical | XS | Critical |
| [TD-003](#td-003) | `lastUsage` race under RLock | Critical | S | High |
| [TD-005](#td-005) | Typo `buildDetebaseClient` | High | XS | Low |
| [TD-007](#td-007) | Variable shadowing in `WaitFor` | High | XS | Low |
| [TD-009](#td-009) | Caller headers override `Content-Type` | High | XS | Medium |
| [TD-010](#td-010) | 2 000+ lines duplicated response parsing | High | L | High |
| [TD-012](#td-012) | Expired token injected after failed refresh | Medium | S | High |
| [TD-014](#td-014) | `ParseResponseBody` panics on nil response | Medium | XS | Medium |
| [TD-015](#td-015) | `DefaultWaitFor` timeout too short | Medium | XS | Medium |
| [TD-016](#td-016) | No structured logging | Medium | L | Medium |
| [TD-017](#td-017) | `WARN` writes to stdout | Medium | XS | Low |
| [TD-020](#td-020) | Test coverage gaps | Low | XL | High |
| [TD-021](#td-021) | Create responses validate metadata ID/URI/Name | Low | M | Medium |

### Recommended execution order

**Wave 1 — Quick Wins** (XS effort, ship same PR): TD-001, TD-002, TD-005, TD-007, TD-009, TD-014, TD-015, TD-017

**Wave 2 — High-value focused fixes** (S effort, High+ impact): TD-003, TD-012

**Wave 3 — Medium fixes** (S effort, Medium/Low impact): *(all resolved)*

**Wave 4 — Large refactors** (L/XL, plan separately): TD-010, TD-016, TD-020, TD-021

---

## Resolved

### TD-006 · Goroutine leak in `WaitFor` on context cancellation — [#116](https://github.com/Arubacloud/sdk-go/issues/116) · **Closed as invalid**
The original claim was that the goroutine writes to `resultCh` unconditionally without guarding on `ctx.Done()`. Investigation of `pkg/async/async_client.go` shows:
1. The goroutine checks `ctxTimeout.Done()` via `select` before each retry attempt (lines 115-120) and after each failed attempt before sleeping (lines 146-152).
2. The result channel is buffered (`make(chan Result[T], 1)`), so the single channel send the goroutine performs can never block.

No goroutine leak exists. Issue #116 closed as invalid.

---

### TD-004 · Interceptor `Bind()` and `Intercept()` are not thread-safe — [#114](https://github.com/Arubacloud/sdk-go/issues/114) · **Closed as working as intended**
The original claim was that concurrent calls to `Bind()` and `Intercept()` produce an unsynchronized read/write on the `interceptFuncs` slice. The proposed fix was a `sync.RWMutex` guarding both methods.

Investigation shows that `Bind` is a construction/setup method — it is only ever called from `tokenmanager.NewStandard` (`internal/impl/auth/tokenmanager/standard/standard.go:54`) during client bootstrap, before the interceptor enters the hot request path. A sibling constructor `NewInterceptorWithFuncs` exists for fully-initialized construction. The race only materialises if a caller mutates the interceptor after bootstrap, which is not the intended usage.

The root cause was the absence of a documented contract. Resolved by adding godoc to `Interceptable.Bind` (`internal/ports/interceptor/interceptor.go`) and the standard impl (`internal/impl/interceptor/standard/standard.go`) stating that `Bind` is construction-only and not safe for concurrent use with `Intercept`. Adding a mutex to the hot path to defend against an unsupported usage pattern was rejected as an unnecessary tradeoff. Issue #114 closed as working as intended.

---

### TD-008 · Polling loop always sleeps before the first attempt, discards final state — [#118](https://github.com/Arubacloud/sdk-go/issues/118) · **Resolved**
`internal/restclient/polling.go` — two bugs in `WaitForResourceState`: (1) `time.Sleep(config.Interval)` was called at the top of the loop, wasting a full interval before the first status check; (2) when the last attempt's getter returned an error, the `continue` skipped the timeout-with-state branch, causing the final timeout error to be generic rather than including the last known state.

Resolved by restructuring the loop: the sleep moved to the bottom, guarded by `if attempt < config.MaxAttempts`; a `lastState` / `lastErr` pair is tracked across iterations so the post-loop timeout error always carries the last observed state (or wraps the last getter error when no state was ever retrieved). `slices.Contains` replaces the manual success/failure loops. First polling tests added in `internal/restclient/polling_test.go`. Issue #118 closed.

---

### TD-013 · Memory proxy `SaveToken` increments ticket before confirming persistent write — [#123](https://github.com/Arubacloud/sdk-go/issues/123) · **Resolved**
`internal/impl/auth/tokenrepository/memory/memory.go` — `saveTicket++` ran before `r.persistentRepository.SaveToken(...)`, violating the double-checked-locking invariant that "a changed ticket means the cache was successfully refreshed". On a failed persistent write the ticket was already bumped, causing concurrent `FetchToken` calls to skip the persistent re-fetch and serve a stale (or nil) in-memory token.

Resolved by reordering: persistent write first; on success, increment ticket and update cache; on error, return immediately leaving both unchanged. Two regression tests added to `memory_test.go`: a unit assertion on the `saveTicket` counter after a failed save, and a behavioural subtest that verifies a subsequent `FetchToken` still reaches the persistent store. Issue #123 closed.

---

### TD-019 · Missing compile-time interface satisfaction checks — [#129](https://github.com/Arubacloud/sdk-go/issues/129) · **Resolved**
Guards were missing for all ~24 resource-level client impls under `internal/clients/`. Adding them directly inside those packages would create an import cycle (`pkg/aruba/builder.go` already imports every internal client package). The guards are consolidated in a new file `pkg/aruba/assertions.go`, using each impl's exported constructor with `nil` args to obtain a typed value of the unexported impl type — the exact same assignment the existing `buildXxxClient` return types already checked implicitly.

The security domain is intentionally excluded: `KMSClient`, `KeyClient`, and `KmipClient` in `pkg/aruba/security.go` are type aliases to concrete pointer types, not interfaces, so satisfaction guards would be degenerate.

The issue description incorrectly identified both logger implementations as missing guards; both already had them (`internal/impl/logger/native/logger.go:11`, `internal/impl/logger/noop/logger.go:6`). Three `for i := 0; i < N; i++` loops in test files were also modernized to `for range N` while in the area. Issue #129 closed.

When TD-018 was subsequently resolved (#146), the multi-arg constructor calls in `assertions.go` were updated to use real (non-nil-dep) sentinel instances to avoid init-time panics.

---

### TD-018 · Injected concrete dependencies not nil-checked in constructors — [#128](https://github.com/Arubacloud/sdk-go/issues/128) · **Resolved**
The issue named two constructors (`NewSecurityGroupsClientImpl`, `NewSecurityGroupRulesClientImpl`). A repo-wide grep (`^func New\w+ClientImpl\(.*,.*\)` against `internal/clients/`) found five constructors in total that accept injected `*xxxClientImpl` dependencies beyond the standard `*restclient.Client`: the two network ones named in the issue, `NewSubnetsClientImpl` (also network), `NewSnapshotsClientImpl`, and `NewRestoreClientImpl` (storage).

All five now `panic("... is required and cannot be nil")` when the injected dep is nil. Corresponding panic-on-nil tests added to the paired test files. `pkg/aruba/assertions_test.go` updated to wire real sentinel instances for the affected constructors (avoiding init-time panics). Issue #128 closed.

---

### TD-011 · Silent failure when parsing error response body — [#121](https://github.com/Arubacloud/sdk-go/issues/121) · **Resolved**
`pkg/types/utils.go` — when a 4xx/5xx response body could not be unmarshalled as JSON, the error was silently discarded and `response.Error` remained `nil`.

Resolved by adding a `DebugLogger` interface to `pkg/types` (one method: `Debugf`) and making `ParseResponseBody` accept it as a required parameter. On JSON-unmarshal failure the function now calls `logger.Debugf(...)` naming the HTTP status code, rather than WARN as the issue proposed — non-JSON error bodies are expected behaviour from this API (proxy HTML pages, plain-text 502s). Using a local interface keeps `pkg/types` independent of `internal/ports/logger`. All 31 call-site files across `internal/clients/` were updated to pass `c.client.Logger()`. Four unit tests added in `pkg/types/utils_test.go`. Issue #121 closed.

---

## Critical

### TD-001 · Parameter order swap in file token repository constructor — [#111](https://github.com/Arubacloud/sdk-go/issues/111)
`pkg/aruba/builder.go` calls `NewFileTokenRepository(options.baseDir, clientID)` but the function signature is `NewFileTokenRepository(clientID, baseDir string)`. The arguments are swapped: the base directory is used as the client ID and vice versa, producing a wrong token file path and breaking file-based token persistence entirely.

**Fix:** Change the call to `NewFileTokenRepository(clientID, options.baseDir)`.

**Effort:** XS — swap 2 arguments at 1 call site.

**Impact:** Critical — file-based token persistence is completely broken.

---

### TD-002 · Static token is never stored — access token parameter silently ignored — [#112](https://github.com/Arubacloud/sdk-go/issues/112)
`internal/impl/auth/tokenrepository/memory/memory.go` — `NewTokenRepositoryWithAccessToken(accessToken string)` ignores its parameter and returns `&TokenRepository{}`. Any client created with `WithToken()` will always get an empty token, causing all requests to fail authentication silently.

**Fix:** Initialize the struct with `token: &auth.Token{AccessToken: accessToken}`.

**Effort:** XS — add 1 struct field initialization.

**Impact:** Critical — `WithToken()` static auth is completely broken; every API call returns 401.

---

### TD-003 · Race condition: `lastUsage` written under read lock in Multitenant — [#113](https://github.com/Arubacloud/sdk-go/issues/113)
`pkg/multitenant/multitenant.go` — `Get()`, `MustGet()`, and `GetOrNil()` all hold `RLock` while writing `e.lastUsage = time.Now()`. Mutating a struct field through a map value under a read lock is a data race detected by the Go race detector.

**Fix:** Upgrade to `Lock()` for the update, or store `lastUsage` as an `atomic.Int64` (Unix nanoseconds) so it can be updated without a write lock.

**Effort:** S — change 3 methods in 1 file; atomic variant needs a small struct change.

**Impact:** High — data race under concurrent multitenant use; triggers the race detector and can cause unpredictable behaviour in production.

---

## High

### TD-005 · Typo in builder function name: `buildDetebaseClient` — [#115](https://github.com/Arubacloud/sdk-go/issues/115)
`pkg/aruba/builder.go` — the function is named `buildDetebaseClient` instead of `buildDatabaseClient`. Harmless at runtime but breaks searchability and violates naming consistency.

**Fix:** Rename to `buildDatabaseClient`.

**Effort:** XS — rename 1 unexported function.

**Impact:** Low — cosmetic only; no runtime effect.

---

### TD-007 · Variable shadowing creates unreachable `AsyncClient` in `WaitFor` nil-call path — [#117](https://github.com/Arubacloud/sdk-go/issues/117)
`pkg/async/async_client.go` — when `callFunc == nil`, a new inner `asyncClient` variable (`:=`) shadows the outer one. The outer variable is discarded; the inner one is returned. The code works by accident but is fragile and confusing.

**Fix:** Remove the inner `:=` and reuse the outer `asyncClient`.

**Effort:** XS — delete 1 `:=` keyword.

**Impact:** Low — works today by accident; purely a correctness and readability cleanup.

---

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
`pkg/async/async_client.go` — cloud resource provisioning (VMs, databases, VPCs) routinely takes several minutes. The defaults (`retries=10`, `baseDelay=10s`, `timeout=60s`) mean callers must always pass custom values or operations will time out spuriously.

**Fix:** Raise defaults to at least `retries=36`, `baseDelay=10s`, `timeout=600s`, or expose a named constant set (e.g., `LongOperationDefaults`) for callers to use.

**Effort:** XS — change 3 constants; add a release note as defaults are a breaking change for callers relying on them.

**Impact:** Medium — prevents spurious timeouts for any caller using `DefaultWaitFor` on real cloud resources.

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

### TD-021 · Create responses do not contract-test resource ID exposure (Metadata.ID / URI) — [#175](https://github.com/Arubacloud/sdk-go/issues/175)

All `*.Create` methods returned success even when the API response omitted required identity fields (`metadata.id`, `metadata.uri`, `metadata.name`). Downstream consumers (e.g., acloud-cli) silently received `nil` pointers and fell back to broken workarounds.

**Fix:** Added `ResourceMetadataResponse.Validate()` and called it in every Create method. URI was initially included but later removed: the Aruba Cloud API does not consistently populate `metadata.uri` on Create responses, causing false-positive failures in production (acloud-cli `project create` error). The validator now only requires `id` and `name`. Added `pkg/types/resource_test.go` unit tests and "missing id/name" subtests in each Create test file.

**Effort:** M — 21 impl files + 21 test files; mechanical but thorough.

**Impact:** Medium — eliminates a class of silent nil-pointer bugs at the API boundary.

---

### TD-020 · Test coverage limited to happy path and empty-ID validation — [#130](https://github.com/Arubacloud/sdk-go/issues/130)
All `internal/clients/*/_test.go` files — existing tests cover successful responses and empty project/resource IDs. Missing: HTTP 4xx/5xx error responses, malformed JSON bodies, network-level errors, `nil` params handling, and request body marshaling failures.

**Fix:** Add table-driven subtests for each error scenario in every client test file using `httptest.NewServer` to return controlled error responses.

**Effort:** XL — ~15 client files × ~5 error scenarios each; requires design of shared test helpers.

**Impact:** High — major regression safety net; issues like TD-011 and TD-014 would have been caught by these tests.

---
