# Architecture

## Client construction (`pkg/aruba/`)

The public entry point is `NewClient(options *Options) (Client, error)` in `pkg/aruba/aruba.go`, which delegates to `buildClient()` in `builder.go`.

`buildClient()` performs cascading construction in this order:
1. Validate options via `options.validate()`
2. Build logger (`buildLogger()`)
3. Build REST client (`buildRESTClient()`) — injects logger and HTTP client
4. Build token manager (`buildTokenManager()`) — binds itself as the last interceptor
5. Build each of the 10 service group clients sequentially

`pkg/aruba.Options` is a fluent builder (~40 methods). Key injection points:
- `WithCustomHTTPClient(*http.Client)` — defaults to `http.DefaultClient`
- `WithCustomLogger(logger.Logger)` / `WithNativeLogger()` / `WithNoLogs()`
- `WithCustomMiddleware(interceptor.Interceptor)` — defaults to `standard.NewInterceptor()`
- `WithToken(token)` or `WithClientCredentials(clientID, secret)` — selects auth strategy

`pkg/aruba.Client` exposes 10 service group accessors (`FromCompute()`, `FromNetwork()`, etc.). Each returns an interface backed by an unexported impl in `internal/clients/<service>/`.

**Cross-client injection:** Some service clients receive other concrete impl clients at build time to enforce resource pre-conditions. For example, `SecurityGroupRulesClientImpl` holds a `*securityGroupsClientImpl` and calls `waitForSecurityGroupActive()` before creating a rule. These dependencies are always concrete types, not interfaces, because they call internal methods not on any interface.

## HTTP request lifecycle (`internal/restclient/`)

`restclient.Client.DoRequest(ctx, method, path, body, queryParams, headers)` follows these steps:

1. Construct full URL from `baseURL + path`
2. Log request details (auth header is redacted as `Bearer [REDACTED]`)
3. Create `*http.Request` with context
4. Attach query parameters to URL
5. Set `Content-Type: application/json` if body is present
6. Merge caller-supplied headers
7. **Run middleware chain** via `middleware.Intercept(ctx, req)` — this is where the auth token is injected
8. Execute via `httpClient.Do(req)`
9. Log response status and headers; re-wrap body for caller (logging consumed the stream)
10. Return `*http.Response`

## Interceptor/middleware chain (`internal/impl/interceptor/`)

The `Interceptor` interface has two methods: `Bind(...InterceptFunc)` and `Intercept(ctx, req)`. `InterceptFunc` is `func(ctx context.Context, r *http.Request) error`.

The standard implementation collects a slice of `InterceptFunc` values and executes them in order on each request. Execution stops on the first error.

The token manager always binds itself **last** via `BindTo(interceptable)`, so auth injection is always the final middleware step. Custom middleware added by the caller via `WithCustomMiddleware` runs before the token manager.

## Auth subsystem (`internal/impl/auth/`)

Core interfaces (defined in `internal/ports/auth/auth.go`):

```
TokenManager       — binds as interceptor, injects Bearer token on each request
TokenRepository    — FetchToken / SaveToken (multiple backends)
ProviderConnector  — RequestToken (OAuth2 client credentials flow)
CredentialsRepository — FetchCredentials (static memory or Vault)
```

**Token injection with double-checked locking:**
1. Read lock: fetch token from repository; capture ticket counter
2. If missing or expired and a connector is configured:
   - Acquire write lock
   - If ticket changed (another goroutine refreshed), reuse the new token
   - Otherwise: call `connector.RequestToken()`, `repository.SaveToken()`, increment ticket
3. Inject `Authorization: Bearer <token>` header

**Token repository implementations:**
- **Memory** — standalone in-memory store; supports configurable `expirationDriftSeconds` safety buffer
- **Memory proxy** — wraps a persistent store (write-through on save, read-through on miss)
- **File** — persists tokens to `<baseDir>/<clientID>.token.json` with `0o600` permissions
- **Redis** — stores tokens by a key derived from `clientID`

`NewTokenProxyWithRandomExpirationDriftSeconds(persistent, maxDrift)` randomizes the expiry drift to avoid synchronized refresh storms across a fleet.

**Credentials repository implementations:**
- **Memory** — holds static `ClientID` + `ClientSecret`
- **Vault** — fetches credentials from HashiCorp Vault using AppRole authentication (KV v2)

The OAuth2 connector (`internal/impl/auth/providerconnector/oauth2/`) uses `golang.org/x/oauth2/clientcredentials` (Client Credentials flow, RFC 6749). HTTP 401 maps to `auth.ErrAuthenticationFailed`, HTTP 403 to `auth.ErrInsufficientPrivileges`.

## Types package (`pkg/types/`)

**`Response[T any]`** — generic wrapper returned by every API call:

```go
type Response[T any] struct {
    Data         *T
    Error        *ErrorResponse
    HTTPResponse *http.Response
    StatusCode   int
    Headers      http.Header
    RawBody      []byte
}
func (r *Response[T]) IsSuccess() bool  // StatusCode 2xx
func (r *Response[T]) IsError() bool    // StatusCode 4xx/5xx
```

**`RequestParameters`** — all optional pointer fields:

```go
type RequestParameters struct {
    Filter     *string
    Sort       *string
    Projection *string
    Accept     *AcceptHeader
    Offset     *int32
    Limit      *int32
    APIVersion *string
}
func (r *RequestParameters) ToQueryParams() map[string]string
func (r *RequestParameters) ToHeaders() map[string]string
```

**`ErrorResponse`** (RFC 7807-based): `Type`, `Title`, `Status`, `Detail`, `Instance`, `TraceID`, `Errors []ValidationError`, `Extensions map[string]interface{}`. Custom `UnmarshalJSON` captures unknown keys into `Extensions` for forward compatibility.

**Validation helpers** in `pkg/types/utils.go`:
`ValidateProject`, `ValidateProjectAndResource`, `ValidateDBaaSResource`, `ValidateDatabaseGrant`, `ValidateVPCResource`, `ValidateSecurityGroupRule`, and more.

**`ParseResponseBody[T any](httpResp)`** — utility function that reads the body, unmarshals into `Data` (2xx) or `Error` (4xx/5xx), and stores raw bytes.

## Async polling (`pkg/async/`)

**`AsyncClient[T]`** — holds a channel and a cached `Result[T]` (protected by `sync.Mutex`). `Await(ctx)` blocks until the result arrives and caches it on first call.

**`WaitFor[T](ctx, retries, baseDelay, timeout, callFunc, checkFunc)`** — core polling loop:
- Launches a goroutine retrying `callFunc()` up to `retries` times
- Fixed `baseDelay` between attempts (no exponential backoff — intentional for predictability)
- Enforces `timeout` as a context deadline
- `checkFunc` receives the full `*Response[T]` to decide success

**Defaults** (`DefaultWaitFor`): `retries=10`, `baseDelay=10s`, `timeout=60s`.

## Multitenant client management (`pkg/multitenant/`)

`Multitenant` manages a `map[string]*entry` (tenant ID → client + `lastUsage` timestamp) behind a `sync.RWMutex`.

```go
New(tenant)                              // create from template Options (deep-copied)
NewFromOptions(tenant, *aruba.Options)   // create from custom options
Add(tenant, aruba.Client)                // register a pre-built client
Get / MustGet / GetOrNil                 // all update lastUsage on access
CleanUp(from time.Duration)              // remove entries idle longer than `from`
```

`StartCleanupRoutine(ctx, tickInterval, fromDuration)` runs a background goroutine that calls `CleanUp` on the given interval (default tick: 1 hour, default idle threshold: 24 hours).

`NewWithTemplate(template *aruba.Options)` deep-copies the template for each `New()` call (slices are deep-copied; `*http.Client`, logger, and middleware are shallow-copied as shared singletons).

## Service client standard method flow

Every resource method in `internal/clients/<service>/` follows this sequence:

1. `c.client.Logger().Debugf("...")` — log the operation and key IDs
2. Call `types.Validate*(...)` — fail fast on nil/empty IDs
3. Inject default `APIVersion` if `params.APIVersion == nil`
4. `params.ToQueryParams()` / `params.ToHeaders()`
5. Marshal body with `json.Marshal(body)` (if applicable)
6. `c.client.DoRequest(ctx, method, path, body, queryParams, headers)`
7. `defer httpResp.Body.Close()`
8. `types.ParseResponseBody[T](httpResp)` or manual unmarshal for complex responses

## Adding a new resource

1. Define request/response types in `pkg/types/<domain>.<resource>.go`.
2. Add API path constants to `internal/clients/<service>/path.go`.
3. Add per-operation API version constants to `internal/clients/<service>/version.go`.
4. Create the resource file `internal/clients/<service>/<resource>.go` — define the interface and `*Impl` struct, implement all methods following the standard flow above.
5. Expose the resource from the service group file `internal/clients/<service>/<group>.go`.
6. If the resource depends on another resource's state, accept the dependency as a constructor parameter (concrete impl type).
7. Wire to `pkg/aruba/Client` if this is a new service group.
