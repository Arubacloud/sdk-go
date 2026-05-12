# Architecture

## Client construction (`pkg/aruba/`)

The public entry point is `NewClient(options *Options) (Client, error)` in `pkg/aruba/aruba.go`, which delegates to `buildClient()` in `builder.go`.

`buildClient()` performs cascading construction in this order:
1. Validate options via `options.validate()`
2. Build REST client via `buildRESTClient()` — internally constructs:
   - HTTP client (`buildHTTPClient()`) — defaults to `http.DefaultClient`
   - Logger (`buildLogger()`)
   - Middleware (`buildMiddleware()`) — builds the token manager and binds it as the last interceptor
3. Build each of the 10 service group clients sequentially

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

## Wrapper layer (`pkg/aruba/`)

The `pkg/aruba/` package is a chainable, error-accumulating, fluent builder façade layered over `internal/clients/*`. Users never call the low-level clients directly; they construct typed wrappers, set properties via chained setters, then pass the wrapper to an adapter method (`Create`, `Get`, `Update`, `Delete`, `List`).

### Triplet pattern

Every `resource_<name>.go` is divided into three banner-separated sections:

```
// ---- Wrapper ----             chainable builder struct + mixin embeds
// ---- Low-level client interface ----   contract the adapter depends on
// ---- Adapter ----             bridges wrapper ↔ internal/clients/<x>
```

- **Wrapper** embeds the relevant mixins (see below), holds a typed `*types.<X>Response`, and exposes chainable setters + read accessors.
- **Low-level interface** is declared inline so only the adapter depends on the concrete `internal/clients/*` impl — allows the adapter to be unit-tested with a mock.
- **Adapter** (`<x>ClientAdapter`) translates wrapper ↔ wire types. Constructor `new<X>ClientAdapter(rest *restclient.Client)` wires to the concrete impl (e.g. `database.NewDatabasesClientImpl(rest)` at `resource_database.go:147`). Adapters are instantiated in `pkg/aruba/builder.go`.

Canonical examples: `resource_database.go` (Family B, minimal), `resource_cloud_server.go` (Family A with action dispatch and `setRefresh` wiring).

### Mixins (`pkg/aruba/mixins.go`)

| Mixin | Responsibility |
|---|---|
| `errMixin` | Setter-time error accumulator; `Err()` returns joined errors; setters always return the receiver so the chain never short-circuits. |
| `metadataMixin` | Request-side name + tag set; `toMetadata()` emits `types.ResourceMetadataRequest`. |
| `regionalMixin` | Holds `Region`; `toLocation()` emits `types.LocationRequest`. Embedded by `zonalMixin`. |
| `zonalMixin` | Extends `regionalMixin` with an availability-zone pointer (wire field `dataCenter`). |
| `projectScopedMixin` | Direct child of a Project; `intoProject(Ref)` extracts `projectID` via `extractID`. |
| `vpcScopedMixin` | Direct child of a VPC; inherits `projectID` from parent. |
| `securityGroupScopedMixin` | Direct child of a SecurityGroup; inherits `vpcID` + `projectID`. |
| `dbaasScopedMixin` | Direct child of a DBaaS; inherits `projectID`. |
| `databaseScopedMixin` | Direct child of a Database (grandchild of DBaaS). |
| `backupScopedMixin` | Direct child of a StorageBackup. |
| `kmsScopedMixin` | Direct child of a KMS instance. |
| `vpnTunnelScopedMixin` | Direct child of a VPN tunnel; tolerates both `vpn-tunnels` and `vpnTunnels` URI forms. |
| `vpcPeeringScopedMixin` | Direct child of a VPC peering; inherits `vpcID` + `projectID`. |
| `responseMetadataMixin` | Holds the post-server `*types.ResourceMetadataResponse`; exposes `ID()`, `RespURI()`, `CreatedAt()`, `Version()`, etc. |
| `statusMixin` | Holds `*types.ResourceStatus` + a `refresh` callback + `terminalStates` map; powers `WaitUntilActive`, `WaitUntilReady`, `WaitUntilStates`. |
| `linkedMixin` | Stores `[]types.LinkedResource` returned by the API. |
| `httpEnvelopeMixin` | Captures StatusCode / Headers / RawBody / `*http.Response` / parsed ErrorResponse after every adapter call. |

`populateHTTPEnvelope[T]` is a package-level generic function (Go does not allow generic methods on structs; `mixins.go:816`).

### Family A vs. Family B

**Family A** — the standard shape: `Metadata{Properties{...}}` envelope on the wire, embeds `metadataMixin` + a regional mixin + `statusMixin` + `responseMetadataMixin`. Most resources (CloudServer, KaaS, DBaaS, VPC, BlockStorage, Job, KMS, …). Canonical reference: `resource_kaas.go:18`:

> `// Family A: regional, Metadata/Properties envelope, location-aware.`

**Family B** — flat request body, no Metadata/Properties boxing, no tags, no location, no `metadataMixin`, no `statusMixin`. Resources: Database, Key, Kmip, User, Grant. Canonical reference: `resource_database.go:18`:

> `// Family B: flat request (no Metadata/Properties boxing, no metadataMixin, no tags, no location).`

Family B sub-variant — **no-Update**: Key and Kmip additionally omit the `Update` operation; the service-group interfaces (`security.go:33-49`) deliberately exclude it. Reflective guards in `resource_key_test.go:664` and `resource_kmip_test.go:896` enforce this at test time.

**Identity quirk in Family B:** `DatabaseResponse` carries no server-side `id` — the name IS the path identifier (`resource_database.go:23`). `Key` returns `KeyResponse.KeyID`; `Kmip` returns `KmipResponse.ID`. All three construct `URI()` client-side from ancestor IDs + the resource ID (e.g. `resource_database.go:55-61`).

### `Ref` interface (`pkg/aruba/ref.go`)

```go
type Ref interface { URI() string; ID() string }
```

- Every typed wrapper satisfies `Ref`. `aruba.URI(s)` returns an opaque `uriRef` for raw-URI callers (`factories.go:8`).
- `extractID(ref, typedExtractor, segment)` — tries the typed `withXxxID` interface first, then falls back to `parseURIIDs` which splits a URI by path segment (`ref.go:63`).
- 25 unexported `withXxxID` interfaces (`withProjectID`, `withDBaaSID`, `withKMSID`, etc.) allow adapters to extract a parent's typed ID without coupling to a concrete wrapper type (`ref.go:75-99`).
- Per-resource `<resource>IDsFromRef(ref)` helpers unwrap deep parent chains — e.g. `databaseIDsFromRef` returns `(projectID, dbaasID, databaseID)` at `resource_database.go:298`.

### `List[T Wrapper]` (`pkg/aruba/list.go`)

Generic paginated container, constrained to `Wrapper { URI(); ID() }`. Carries `items`, `total`, pagination link URLs (`self/prev/next/first/last`), `callerOpts`, `raw` HTTP envelope, and a `refetch` callback. Navigation methods: `Items()`, `Total()`, `HasNext()`, `Next(ctx)`, `All(ctx, yield)`. The `refetch` callback is structurally wired in every adapter but currently returns a "not yet wired" error — pagination requires re-calling `List` with updated `CallOption` paging parameters.

### Wait helpers and async

`statusMixin` provides three wait methods (`mixins.go:739-787`), all backed by `pkg/async.WaitFor[any]`:

- `WaitUntilActive` — targets state `"Active"`.
- `WaitUntilReady` — accepts `"Active"`, `"NotUsed"`, `"InUse"`, or `"Used"` (covers attach/detach resources in either steady state).
- `WaitUntilStates(ctx, targets, opts...)` — the work-horse; requires a `refresh` callback installed by the adapter.

Adapters install the `refresh` closure post-`Create`/`Get`/`List` (e.g. `resource_cloud_server.go:476-485`). The closure re-`Get`s the resource and hydrates the same wrapper in place so each polling tick sees the updated state.

Defaults: 60 retries × 10s base delay × 600s timeout (`mixins.go:658-664`).

**Per-resource specialised waiters:**
- `*Kmip.WaitUntilCertificateAvailable` (`resource_kmip.go:69`) — Family B has no `statusMixin`; drives `async.WaitFor` directly against `KmipResponse.Status` with explicit terminal map `kmipTerminalStates` (`resource_kmip.go:41`).
- `*BlockStorage.WaitUntilUsed` / `WaitUntilNotUsed` (`resource_block_storage.go:266-275`) and `*ElasticIP` equivalents — attach/detach lifecycle, three positive terminals (`InUse`, `Used`, `NotUsed`).

### HTTP envelope and typed `*HTTPError`

After every adapter call, `populateHTTPEnvelope[T]` snapshots StatusCode / Headers / RawBody / `*http.Response` / parsed `*types.ErrorResponse` onto the wrapper's `httpEnvelopeMixin`. On non-2xx the adapter returns:

```go
&HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
```

(`pkg/aruba/errors.go:11-15`). The wrapper retains the envelope for diagnostics (`cs.StatusCode()`, `cs.Headers()`, `cs.RawError()`) even after an error.

### Error accumulation

`errMixin.addErr` collects setter-time errors without breaking the chain. Every adapter `Create` and `Update` opens with:

```go
if err := X.Err(); err != nil { return X, err }
```

Followed by per-adapter friendly validation. Family B resources with deep parent chains do the most validation (Database, Key, Kmip, User, Grant check ProjectID, parent-scope ID, and name individually). Family A adapters mostly only check `ProjectID() != ""`.

Sub-builder errors are drained into the parent's `errMixin` at attachment time (e.g. `WithIKESettings(*VPNIKE)` at `resource_vpn_tunnel.go:84-89`), so `job.AddStep(step)` propagates any accumulated step errors into the job.

### Client integration

The central `Client` interface in `pkg/aruba/client.go` exposes ten `From<Domain>()` accessors. Each returns a per-domain interface (defined in `audit.go`, `compute.go`, `container.go`, `database.go`, `metric.go`, `network.go`, `project.go`, `schedule.go`, `security.go`, `storage.go`) whose methods return per-resource service-group interfaces.

Call chain: `arubaClient.FromCompute().CloudServers().Create(ctx, cs)` → `cloudServersClientAdapter.Create` → `compute.NewCloudServersClientImpl(rest).Create`. Adapters are constructed in `pkg/aruba/builder.go` via `build<Domain>Client` → `new<Resource>ClientAdapter(rest)`.

### Non-standard cases and translation mechanisms

**Resources without a public `Named` setter:**
- `LoadBalancer` — name comes only from the response; `named()` is called inside `fromResponse` (`resource_load_balancer.go:65`).
- `Alert` and `AuditEvent` — read-only / list-only; `URI()` returns `""`.
- `User` — uses `WithUsername` instead; the username IS the path identifier (`resource_user.go:46`).
- `Grant` — no name setter at all; the opaque server-supplied grant ID is recoverable from a URI Ref only (`resource_grant.go:22-30`).

**Deep parent chains (4-level):**
- Grant → Database → DBaaS → Project.
- SecurityRule → SecurityGroup → VPC → Project.
- VPCPeeringRoute → VPCPeering → VPC → Project.

**Body-side parent refs (vs. path-side):** Snapshot, StorageBackup, StorageRestore, DBaaSBackup are `IntoProject(...)` but reference their source resource in the wire body via `FromVolume(Ref)` / `FromDBaaS(Ref)` / `FromDatabase(Ref)`. An empty URI from the `Ref` goes onto the error sink, not the wire.

**Job lifecycle quirk:** Jobs persist as historical records after Delete. `jobTerminalStates` (`resource_job.go:318-322`) enumerates only `Active`, `Error`, `Failed` — no `Deleted` or `Cancelled`. Polling for 404 after Delete always exhausts the wait budget. The canonical reference explaining this is `examples/all-resources/resource_job.go:65-77`.

**Shape-collapsing setters** — one wrapper method sets multiple wire fields:
- `*Job.OneShotAt(t)` sets `JobType=OneShot` + `ScheduleAt`; `*Job.WithCron(expr)` + `RecurringUntil(t)` set `JobType=Recurring` + `Cron` + `ExecuteUntil`. All three are mode-locked via `requireMode` (`resource_job.go:108`).
- `*VPNTunnel.WithIKESettings/ESPSettings/PSKSettings/IPConfig` each attach a sub-builder and drain its `errMixin` into the tunnel (`resource_vpn_tunnel.go:82-123`).
- `*SecurityRule.WithTargetCIDR` / `WithTargetSecurityGroup` set the same wire `target` field but stamp different `Kind` values; calling both records an error (`resource_security_rule.go:77-100`).
- `*NodePool.WithAutoscaling(min, max)` sets `autoscaling=true` + `minCount` + `maxCount` in one call (`resource_kaas_nodepool.go:41`).

**Sub-builders without an adapter** — used only inside a parent, no CRUD: `JobStep` (inside `Job.AddStep`), `NodePool` (inside `KaaS.AddNodePool`), `VPNIKE` / `VPNESP` / `VPNPSK` / `VPNIPConfig` (inside `VPNTunnel`), `SubnetDHCP` (inside `Subnet`). Each has its own `errMixin` drained into the parent at attachment time.
