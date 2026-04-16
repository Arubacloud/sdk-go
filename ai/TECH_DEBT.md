# TECH_DEBT.md — Technical Debt & Refactoring Backlog

Issues are grouped by severity. Address Critical items before new features ship; High items before any public release.

## Resolved

None yet.

---

## Critical

### TD-001 · Parameter order swap in file token repository constructor — [#111](https://github.com/Arubacloud/sdk-go/issues/111)
`pkg/aruba/builder.go` calls `NewFileTokenRepository(options.baseDir, clientID)` but the function signature is `NewFileTokenRepository(clientID, baseDir string)`. The arguments are swapped: the base directory is used as the client ID and vice versa, producing a wrong token file path and breaking file-based token persistence entirely.

**Fix:** Change the call to `NewFileTokenRepository(clientID, options.baseDir)`.

---

### TD-002 · Static token is never stored — access token parameter silently ignored — [#112](https://github.com/Arubacloud/sdk-go/issues/112)
`internal/impl/auth/tokenrepository/memory/memory.go` — `NewTokenRepositoryWithAccessToken(accessToken string)` ignores its parameter and returns `&TokenRepository{}`. Any client created with `WithToken()` will always get an empty token, causing all requests to fail authentication silently.

**Fix:** Initialize the struct with `token: &auth.Token{AccessToken: accessToken}`.

---

### TD-003 · Race condition: `lastUsage` written under read lock in Multitenant — [#113](https://github.com/Arubacloud/sdk-go/issues/113)
`pkg/multitenant/multitenant.go` — `Get()`, `MustGet()`, and `GetOrNil()` all hold `RLock` while writing `e.lastUsage = time.Now()`. Mutating a struct field through a map value under a read lock is a data race detected by the Go race detector.

**Fix:** Upgrade to `Lock()` for the update, or store `lastUsage` as an `atomic.Int64` (Unix nanoseconds) so it can be updated without a write lock.

---

### TD-004 · Interceptor `Bind()` and `Intercept()` are not thread-safe — [#114](https://github.com/Arubacloud/sdk-go/issues/114)
`internal/impl/interceptor/standard/standard.go` — `Bind()` appends to `interceptFuncs` without any synchronization. Concurrent calls to `Bind()` and `Intercept()` produce an unsynchronized read/write on a slice — a data race. While `Bind` is today called only at construction time, the interface contract does not prevent concurrent use.

**Fix:** Add a `sync.RWMutex` to the struct. `Bind()` takes a write lock; `Intercept()` copies the slice under a read lock then iterates the copy.

---

## High

### TD-005 · Typo in builder function name: `buildDetebaseClient` — [#115](https://github.com/Arubacloud/sdk-go/issues/115)
`pkg/aruba/builder.go` — the function is named `buildDetebaseClient` instead of `buildDatabaseClient`. Harmless at runtime but breaks searchability and violates naming consistency.

**Fix:** Rename to `buildDatabaseClient`.

---

### TD-006 · Goroutine leak in `WaitFor` on context cancellation — [#116](https://github.com/Arubacloud/sdk-go/issues/116)
`pkg/async/async_client.go` — the goroutine launched by `WaitFor()` writes to `resultCh` unconditionally. If the caller's context is cancelled and `Await()` is never called (or returns early), the goroutine blocks forever on the channel send. The channel buffer is 1, so if the slot is already filled, the goroutine leaks.

**Fix:** Wrap the channel send in a `select` with `ctx.Done()`.

---

### TD-007 · Variable shadowing creates unreachable `AsyncClient` in `WaitFor` nil-call path — [#117](https://github.com/Arubacloud/sdk-go/issues/117)
`pkg/async/async_client.go` — when `callFunc == nil`, a new inner `asyncClient` variable (`:=`) shadows the outer one. The outer variable is discarded; the inner one is returned. The code works by accident but is fragile and confusing.

**Fix:** Remove the inner `:=` and reuse the outer `asyncClient`.

---

### TD-008 · Polling loop always sleeps before the first attempt — [#118](https://github.com/Arubacloud/sdk-go/issues/118)
`internal/restclient/polling.go` — `time.Sleep(config.Interval)` is called at the top of each iteration, including the very first one. This adds an unnecessary 5-second delay before the initial state check. Additionally, the last iteration checks `attempt == config.MaxAttempts` and returns a timeout error after having just called `getter`, discarding the final state.

**Fix:** Move the sleep to the bottom of the loop (or skip when `attempt == 1`) and always check the state before emitting a timeout error.

---

### TD-009 · Caller headers can override SDK-controlled `Content-Type` — [#119](https://github.com/Arubacloud/sdk-go/issues/119)
`internal/restclient/client.go` — the code sets `Content-Type: application/json` then overwrites headers with caller-supplied values using `req.Header.Set(k, v)`. A caller can silently override `Content-Type` (and in principle `Authorization`), breaking server-side request parsing.

**Fix:** Apply caller headers before SDK-controlled headers, or explicitly protect `Content-Type` and `Authorization` from being overridden.

---

### TD-010 · Create/Update methods duplicate `ParseResponseBody` logic across ~15 client files — [#120](https://github.com/Arubacloud/sdk-go/issues/120)
`internal/clients/compute/cloudserver.go`, `internal/clients/compute/keypair.go`, `internal/clients/network/vpc.go`, `internal/clients/network/security-group.go`, and ~11 more — List/Get/Delete methods call `types.ParseResponseBody[T]()`. Create/Update methods in the same files manually re-implement the same logic: marshal body → `DoRequest` → `io.ReadAll` → construct `Response[T]` wrapper → unmarshal success/error branches. This is ~40 lines duplicated per mutating method, totalling 2 000+ lines across the codebase. Any bug fix to the parsing logic must be applied in all locations.

**Fix:** Add a generic helper to `restclient.Client`:
```go
func DoAndParse[T any](c *Client, ctx context.Context, method, path string,
    body io.Reader, queryParams, headers map[string]string) (*types.Response[T], error)
```
Replace all manual implementations with calls to this helper.

---

## Medium

### TD-011 · Silent failure when parsing error response body — [#121](https://github.com/Arubacloud/sdk-go/issues/121)
`pkg/types/utils.go` — in `ParseResponseBody`, when a 4xx/5xx response body cannot be unmarshalled as JSON, the error is swallowed silently (`if err == nil { ... }`). `response.Error` remains `nil`, making it impossible for the caller to understand what the server returned.

**Fix:** Log the unmarshal error at WARN level, or document that `RawBody` should be used as the fallback for non-JSON error bodies.

---

### TD-012 · Token manager injects expired/nil token after failed refresh — [#122](https://github.com/Arubacloud/sdk-go/issues/122)
`internal/impl/auth/tokenmanager/standard/standard.go` — in the write-lock branch, if the ticket has already changed (another goroutine refreshed), the code re-fetches from the repository. If that second fetch fails, execution falls through and injects an empty or expired token into the `Authorization` header with no error returned to the caller.

**Fix:** Check the error from the second `FetchToken()` and return it before reaching the header-injection step.

---

### TD-013 · Memory proxy `SaveToken` increments ticket before confirming persistent write — [#123](https://github.com/Arubacloud/sdk-go/issues/123)
`internal/impl/auth/tokenrepository/memory/memory.go` — `saveTicket++` runs before `r.persistentRepository.SaveToken(...)`. If the persistent write fails, the ticket has already been incremented, invalidating the in-memory cache. Subsequent reads see a miss and re-read from persistent storage, which still has the old token.

**Fix:** Increment the ticket only after a successful persistent write.

---

### TD-014 · `ParseResponseBody` panics on nil `httpResp` — [#124](https://github.com/Arubacloud/sdk-go/issues/124)
`pkg/types/utils.go` — the function calls `io.ReadAll(httpResp.Body)` without a nil guard. If `DoRequest` returns a nil response alongside a non-nil error and the caller forgets to check the error, the function panics.

**Fix:** Add `if httpResp == nil { return nil, fmt.Errorf("http response is nil") }` at the top of the function.

---

### TD-015 · `DefaultWaitFor` timeout of 60 s is too short for cloud operations — [#125](https://github.com/Arubacloud/sdk-go/issues/125)
`pkg/async/async_client.go` — cloud resource provisioning (VMs, databases, VPCs) routinely takes several minutes. The defaults (`retries=10`, `baseDelay=10s`, `timeout=60s`) mean callers must always pass custom values or operations will time out spuriously.

**Fix:** Raise defaults to at least `retries=36`, `baseDelay=10s`, `timeout=600s`, or expose a named constant set (e.g., `LongOperationDefaults`) for callers to use.

---

### TD-016 · Logger interface supports only printf-style formatting — no structured logging — [#126](https://github.com/Arubacloud/sdk-go/issues/126)
`internal/ports/logger/logger.go` — all log calls use `%s`/`%v` format strings. In production environments using log aggregators (ELK, Loki, Cloud Logging), querying by structured fields (project ID, resource ID, trace ID) requires string parsing. The SDK's `ErrorResponse` already carries a `TraceID` field that is never emitted as a structured log field.

**Fix:** Add optional structured variants to the logger interface, or adopt `slog` (stdlib since Go 1.21) as the native logger implementation.

---

### TD-017 · `WARN` log level writes to `os.Stdout` instead of `os.Stderr` — [#127](https://github.com/Arubacloud/sdk-go/issues/127)
`internal/impl/logger/native/logger.go` — `DEBUG` and `INFO` correctly go to `os.Stdout`; `ERROR` correctly goes to `os.Stderr`. `WARN` also goes to `os.Stdout`, which breaks shell pipelines and container log routers that separate informational output from diagnostic output.

**Fix:** Change `WARN` to write to `os.Stderr`.

---

### TD-018 · Injected concrete dependencies not nil-checked in constructors — [#128](https://github.com/Arubacloud/sdk-go/issues/128)
`internal/clients/network/security-group.go`, `internal/clients/network/security-group-rule.go` — `SecurityGroupsClientImpl` receives a `*vpcsClientImpl` dependency and `SecurityGroupRulesClientImpl` receives a `*securityGroupsClientImpl`. Neither constructor validates that the injected pointer is non-nil. A nil pointer passed at build time causes a panic deep inside a `Create` call, far from the construction site.

**Fix:** Add explicit nil checks in all constructors that accept dependency pointers.

---

## Low

### TD-019 · Missing compile-time interface satisfaction checks — [#129](https://github.com/Arubacloud/sdk-go/issues/129)
Logger implementations (`internal/impl/logger/native/`, `internal/impl/logger/noop/`) and most client `*Impl` types lack `var _ Interface = (*Impl)(nil)` guards. Only the interceptor has this check. A missed method or signature change will only surface at runtime rather than at compile time.

**Fix:** Add `var _ logger.Logger = (*DefaultLogger)(nil)` (and equivalent for each impl type) at package scope in every implementation file.

---

### TD-020 · Test coverage limited to happy path and empty-ID validation — [#130](https://github.com/Arubacloud/sdk-go/issues/130)
All `internal/clients/*/_test.go` files — existing tests cover successful responses and empty project/resource IDs. Missing: HTTP 4xx/5xx error responses, malformed JSON bodies, network-level errors, `nil` params handling, and request body marshaling failures.

**Fix:** Add table-driven subtests for each error scenario in every client test file using `httptest.NewServer` to return controlled error responses.

---
