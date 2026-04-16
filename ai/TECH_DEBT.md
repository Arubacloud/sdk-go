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
| [TD-004](#td-004) | Interceptor not thread-safe | Critical | S | Medium |
| [TD-005](#td-005) | Typo `buildDetebaseClient` | High | XS | Low |
| [TD-006](#td-006) | Goroutine leak in `WaitFor` | High | S | High |
| [TD-007](#td-007) | Variable shadowing in `WaitFor` | High | XS | Low |
| [TD-008](#td-008) | Polling sleeps before first attempt | High | S | Medium |
| [TD-009](#td-009) | Caller headers override `Content-Type` | High | XS | Medium |
| [TD-010](#td-010) | 2 000+ lines duplicated response parsing | High | L | High |
| [TD-011](#td-011) | Silent error body parse failure | Medium | S | Medium |
| [TD-012](#td-012) | Expired token injected after failed refresh | Medium | S | High |
| [TD-013](#td-013) | Memory proxy ticket incremented before write | Medium | XS | Medium |
| [TD-014](#td-014) | `ParseResponseBody` panics on nil response | Medium | XS | Medium |
| [TD-015](#td-015) | `DefaultWaitFor` timeout too short | Medium | XS | Medium |
| [TD-016](#td-016) | No structured logging | Medium | L | Medium |
| [TD-017](#td-017) | `WARN` writes to stdout | Medium | XS | Low |
| [TD-018](#td-018) | Nil deps not checked in constructors | Medium | S | Low |
| [TD-019](#td-019) | Missing interface satisfaction checks | Low | S | Low |
| [TD-020](#td-020) | Test coverage gaps | Low | XL | High |

### Recommended execution order

**Wave 1 — Quick Wins** (XS effort, ship same PR): TD-001, TD-002, TD-005, TD-007, TD-009, TD-013, TD-014, TD-015, TD-017

**Wave 2 — High-value focused fixes** (S effort, High+ impact): TD-003, TD-006, TD-012

**Wave 3 — Medium fixes** (S effort, Medium/Low impact): TD-004, TD-008, TD-011, TD-018, TD-019

**Wave 4 — Large refactors** (L/XL, plan separately): TD-010, TD-016, TD-020

---

## Resolved

None yet.

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

### TD-004 · Interceptor `Bind()` and `Intercept()` are not thread-safe — [#114](https://github.com/Arubacloud/sdk-go/issues/114)
`internal/impl/interceptor/standard/standard.go` — `Bind()` appends to `interceptFuncs` without any synchronization. Concurrent calls to `Bind()` and `Intercept()` produce an unsynchronized read/write on a slice — a data race. While `Bind` is today called only at construction time, the interface contract does not prevent concurrent use.

**Fix:** Add a `sync.RWMutex` to the struct. `Bind()` takes a write lock; `Intercept()` copies the slice under a read lock then iterates the copy.

**Effort:** S — add a mutex field and guards in 1 file.

**Impact:** Medium — latent race that does not trigger today (Bind is only called at construction), but the public interface offers no guarantee.

---

## High

### TD-005 · Typo in builder function name: `buildDetebaseClient` — [#115](https://github.com/Arubacloud/sdk-go/issues/115)
`pkg/aruba/builder.go` — the function is named `buildDetebaseClient` instead of `buildDatabaseClient`. Harmless at runtime but breaks searchability and violates naming consistency.

**Fix:** Rename to `buildDatabaseClient`.

**Effort:** XS — rename 1 unexported function.

**Impact:** Low — cosmetic only; no runtime effect.

---

### TD-006 · Goroutine leak in `WaitFor` on context cancellation — [#116](https://github.com/Arubacloud/sdk-go/issues/116)
`pkg/async/async_client.go` — the goroutine launched by `WaitFor()` writes to `resultCh` unconditionally. If the caller's context is cancelled and `Await()` is never called (or returns early), the goroutine blocks forever on the channel send. The channel buffer is 1, so if the slot is already filled, the goroutine leaks.

**Fix:** Wrap the channel send in a `select` with `ctx.Done()`.

**Effort:** S — add a `select` block in 1 goroutine closure.

**Impact:** High — goroutine leak under cancellation in production; accumulated leaks cause memory exhaustion in long-lived services.

---

### TD-007 · Variable shadowing creates unreachable `AsyncClient` in `WaitFor` nil-call path — [#117](https://github.com/Arubacloud/sdk-go/issues/117)
`pkg/async/async_client.go` — when `callFunc == nil`, a new inner `asyncClient` variable (`:=`) shadows the outer one. The outer variable is discarded; the inner one is returned. The code works by accident but is fragile and confusing.

**Fix:** Remove the inner `:=` and reuse the outer `asyncClient`.

**Effort:** XS — delete 1 `:=` keyword.

**Impact:** Low — works today by accident; purely a correctness and readability cleanup.

---

### TD-008 · Polling loop always sleeps before the first attempt — [#118](https://github.com/Arubacloud/sdk-go/issues/118)
`internal/restclient/polling.go` — `time.Sleep(config.Interval)` is called at the top of each iteration, including the very first one. This adds an unnecessary 5-second delay before the initial state check. Additionally, the last iteration checks `attempt == config.MaxAttempts` and returns a timeout error after having just called `getter`, discarding the final state.

**Fix:** Move the sleep to the bottom of the loop (or skip when `attempt == 1`) and always check the state before emitting a timeout error.

**Effort:** S — restructure the loop in 1 function.

**Impact:** Medium — 5 s wasted per poll operation; final state is discarded causing potentially misleading timeout errors.

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

### TD-011 · Silent failure when parsing error response body — [#121](https://github.com/Arubacloud/sdk-go/issues/121)
`pkg/types/utils.go` — in `ParseResponseBody`, when a 4xx/5xx response body cannot be unmarshalled as JSON, the error is swallowed silently (`if err == nil { ... }`). `response.Error` remains `nil`, making it impossible for the caller to understand what the server returned.

**Fix:** Log the unmarshal error at WARN level, or document that `RawBody` should be used as the fallback for non-JSON error bodies.

**Effort:** S — add a logger call or code comment in 1 function.

**Impact:** Medium — improves debuggability of non-JSON error responses (e.g., plain-text 502 from a proxy); no production failure risk.

---

### TD-012 · Token manager injects expired/nil token after failed refresh — [#122](https://github.com/Arubacloud/sdk-go/issues/122)
`internal/impl/auth/tokenmanager/standard/standard.go` — in the write-lock branch, if the ticket has already changed (another goroutine refreshed), the code re-fetches from the repository. If that second fetch fails, execution falls through and injects an empty or expired token into the `Authorization` header with no error returned to the caller.

**Fix:** Check the error from the second `FetchToken()` and return it before reaching the header-injection step.

**Effort:** S — add 1 error check and early return in 1 function.

**Impact:** High — prevents silent auth failure after a transient token storage error; without the fix, the caller receives a 401 with no indication of the root cause.

---

### TD-013 · Memory proxy `SaveToken` increments ticket before confirming persistent write — [#123](https://github.com/Arubacloud/sdk-go/issues/123)
`internal/impl/auth/tokenrepository/memory/memory.go` — `saveTicket++` runs before `r.persistentRepository.SaveToken(...)`. If the persistent write fails, the ticket has already been incremented, invalidating the in-memory cache. Subsequent reads see a miss and re-read from persistent storage, which still has the old token.

**Fix:** Increment the ticket only after a successful persistent write.

**Effort:** XS — move 1 line after the error-check block.

**Impact:** Medium — prevents transient cache corruption when the persistent store is briefly unavailable.

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

### TD-018 · Injected concrete dependencies not nil-checked in constructors — [#128](https://github.com/Arubacloud/sdk-go/issues/128)
`internal/clients/network/security-group.go`, `internal/clients/network/security-group-rule.go` — `SecurityGroupsClientImpl` receives a `*vpcsClientImpl` dependency and `SecurityGroupRulesClientImpl` receives a `*securityGroupsClientImpl`. Neither constructor validates that the injected pointer is non-nil. A nil pointer passed at build time causes a panic deep inside a `Create` call, far from the construction site.

**Fix:** Add explicit nil checks in all constructors that accept dependency pointers.

**Effort:** S — add nil checks to ~4 constructors.

**Impact:** Low — defensive; prevents future wiring bugs from panicking far from their source.

---

## Low

### TD-019 · Missing compile-time interface satisfaction checks — [#129](https://github.com/Arubacloud/sdk-go/issues/129)
Logger implementations (`internal/impl/logger/native/`, `internal/impl/logger/noop/`) and most client `*Impl` types lack `var _ Interface = (*Impl)(nil)` guards. Only the interceptor has this check. A missed method or signature change will only surface at runtime rather than at compile time.

**Fix:** Add `var _ logger.Logger = (*DefaultLogger)(nil)` (and equivalent for each impl type) at package scope in every implementation file.

**Effort:** S — mechanical addition of `var _` lines in ~20 files.

**Impact:** Low — compile-time safety net; no runtime effect until a signature diverges.

---

### TD-020 · Test coverage limited to happy path and empty-ID validation — [#130](https://github.com/Arubacloud/sdk-go/issues/130)
All `internal/clients/*/_test.go` files — existing tests cover successful responses and empty project/resource IDs. Missing: HTTP 4xx/5xx error responses, malformed JSON bodies, network-level errors, `nil` params handling, and request body marshaling failures.

**Fix:** Add table-driven subtests for each error scenario in every client test file using `httptest.NewServer` to return controlled error responses.

**Effort:** XL — ~15 client files × ~5 error scenarios each; requires design of shared test helpers.

**Impact:** High — major regression safety net; issues like TD-011 and TD-014 would have been caught by these tests.

---
