# Conventions

## Package boundary

- Public API surface lives in `pkg/` — types, interfaces, and the `NewClient` entry point
- Concrete implementations live in `internal/` — not importable by external modules
- Only interfaces and data types are exported; concrete impl structs are always unexported

## Naming

### Types

| Kind | Suffix | Example |
|---|---|---|
| Input/request struct | `Request` | `CloudServerRequest` |
| Single output struct | `Response` | `CloudServerResponse` |
| Collection output struct | `List` (embeds `ListResponse`) | `CloudServerList` |
| Service client interface | none (domain name) | `CloudServersClient` |
| Service client implementation | `Impl` (unexported) | `cloudServersClientImpl` |
| Constructor | `New<TypeName>` | `NewCloudServersClientImpl` |

### Files

- Types in `pkg/types/`: `<domain>.<resource>.go` — e.g., `compute.cloudserver.go`, `network.vpc.go`
- Client files in `internal/clients/<domain>/`: one file per resource — e.g., `cloudserver.go`
- Each client file has a paired test file: `cloudserver_test.go`
- Shared per-domain constants split into two files:
  - `path.go` — API path constants (e.g., `CloudServersPath`)
  - `version.go` — per-operation API version constants (e.g., `ComputeCloudServerList`)
- Optional `common.go` for helpers shared across resources in the same domain

### Constants

- Path constants: `PascalCase` (e.g., `CloudServersPath`, `CloudServerPath`)
- API version constants: `<Domain><Resource><Action>` (e.g., `ComputeCloudServerList`, `NetworkVPCCreate`)

## Method receivers

All methods on impl types use **pointer receivers**: `func (c *cloudServersClientImpl) List(...)`. No value receivers.

## Service method structure

Every resource method must follow this sequence exactly:

1. `c.client.Logger().Debugf("...")` — log operation and key IDs at the start
2. `types.Validate*(...)` — fail fast before any HTTP call
3. If `params == nil`, initialize `params = &types.RequestParameters{}`
4. If `params.APIVersion == nil`, set the domain-specific version constant
5. `params.ToQueryParams()` / `params.ToHeaders()` — convert parameters
6. Marshal body with `json.Marshal` if the method takes a body
7. `c.client.DoRequest(ctx, method, path, body, queryParams, headers)`
8. `defer httpResp.Body.Close()`
9. `types.ParseResponseBody[T](httpResp, c.client.Logger())` for standard responses; manual unmarshal only for complex cases (the logger parameter is required — it logs non-JSON error bodies at `Debug` level)

## Error handling

- Check `resp.IsSuccess()` / `resp.IsError()` before accessing `resp.Data`; never inspect raw HTTP status codes in business logic
- Validation errors (pre-request) return as the `error` return value using `fmt.Errorf("project cannot be empty")`
- API errors (HTTP 4xx/5xx) are unmarshaled into `resp.Error` (`*types.ErrorResponse`); field-level details are in `resp.Error.Errors []ValidationError`
- Use `fmt.Errorf("context: %w", err)` for error wrapping (Go 1.13+ idiom)

## RequestParameters nil-safety

`RequestParameters` is always a pointer parameter. Service methods must handle a `nil` input — create a new struct rather than panicking. Never assume the caller provided an API version.

## Logging

Logger interface (`internal/ports/logger/logger.go`) has four methods: `Debugf`, `Infof`, `Warnf`, `Errorf`. Access it via `c.client.Logger()`.

- Use `Debugf` at the start of every service method
- Never log a raw `Authorization` header; the restclient already redacts it as `Bearer [REDACTED]`
- Use `noop.NoOpLogger{}` in all tests

## Testing

- Each resource file has a paired `_test.go` in the same package
- Tests use `httptest.NewServer` to mock both the `/token` endpoint and the resource API endpoint
- Test setup: `restclient.NewClient(baseURL, httpClient, interceptor, logger)` with `noop.NoOpLogger{}`
- Use subtests: `t.Run("scenario", func(t *testing.T) { ... })`
- Mocks generated with MockGen (`go.uber.org/mock/gomock`); generated files are named `zz_mock_<type>_test.go`

## Documentation

- Package comments: `// Package <name> provides ...`
- Type comments: start with the type name — `// CloudServersClient is the interface for ...`
- Method comments: start with the method name — `// List retrieves all cloud servers ...`
- Document all exported interfaces and their methods individually
- Path and version constants are self-documenting by name; minimal comments needed

## Wrapper-layer conventions (`pkg/aruba/`)

### Per-resource file layout

Every `resource_<name>.go` is divided into three banner-separated sections:

```
// ---- Wrapper ----                       chainable builder + mixin embeds
// ---- Low-level client interface ----    contract the adapter depends on
// ---- Adapter ----                       translates wrapper ↔ types.<X>Request/Response
```

A per-resource `<resource>IDsFromRef(ref Ref)` helper lives at the bottom and is used by adapter `Get` / `Delete` / `List`. Pure sub-builders (JobStep, NodePool, VPNIKE, …) have only the Wrapper section.

### Setter-verb taxonomy

All chainable setters follow `func (rcv *T) Verb(...) *T` and return the receiver. Verbs cluster into a fixed vocabulary:

| Verb prefix | Use | Example |
|---|---|---|
| `Named(...)` | set the resource name | `*Database.Named("mydb")` |
| `With<Prop>(...)` | generic scalar property setter | `*KaaS.WithKubernetesVersion(v)` |
| `Into<Parent>(...)` | bind to parent via `Ref` | `*Database.IntoDBaaS(d)` |
| `In<Region\|Zone>(...)` | geographical placement | `*VPC.InRegion(aruba.RegionITBGBergamo)` |
| `Of<Type\|Flavor\|...>(...)` | typing / sizing / affinity | `*CloudServer.OfFlavor(aruba.CloudServerFlavorCSO2A4)` |
| `From<Source>(...)` | source / origin `Ref` | `*Snapshot.FromVolume(bs)` |
| `Add<X>(...)` / `Replace<X>s(...)` / `Remove<X>(...)` | collection mutation | `*VPNTunnel.AddTag("env=prod")` |
| `Set<X>()` / `Unset<X>()` | explicit boolean toggle (builder) | `*BlockStorage.SetBootable()` |
| `Is<X>()` (read) / `As<X>()`, `Not<X>()` (write) | paired-bool idiom for `*bool` fields | `*SecurityGroup.AsDefault()`, `*VPC.NotDefault()` |

**Important:** `Set*` on its own (not `Unset*`) is reserved for non-chainable, side-effecting action methods (`*CloudServer.SetPassword(ctx, pw, opts...)` in `resource_cloud_server.go`) — never use it for a builder setter.

### Wire-level translation rules

The wrapper API never echoes a wire field name verbatim if a clearer Go-idiomatic name exists:

- **Units explicit in the method name:** `WithSizeGB(int)` (wire field `size` / JSON `sizeGb`) on `*BlockStorage` and `*DBaaS`.
- **Paired-bool idiom for tri-state `*bool` flags:** `IsDefault()` read + `AsDefault()` / `NotDefault()` write on SecurityGroup, VPC, Subnet, Project — preserves the `nil` / `false` / `true` distinction across the wire.
- **`Ref`-accepting setters hide URI assembly:** `FromVolume(Ref)`, `FromSnapshot(Ref)`, `FromDBaaS(Ref)`, `FromDatabase(Ref)` — an empty URI from the `Ref` goes on the error sink, not on the wire.
- **Sub-builder errors propagate:** `WithIKESettings(*VPNIKE)` drains the sub-builder's `errMixin` into the tunnel's accumulator (`resource_vpn_tunnel.go`).
- **Shape-collapsing:** `*Job.OneShotAt(t)` sets both `JobType` and `ScheduleAt`; `WithCron + RecurringUntil` sets `JobType` + `Cron` + `ExecuteUntil`. Mode conflicts record a setter-time error via `requireMode` in `resource_job.go`.

### Typed-enum aliases (`pkg/aruba/aliases.go`)

Public callers never pass raw strings. Every domain that accepts an enum has typed aliases in `aliases.go`:

| Domain | Alias prefix | Typical setter |
|---|---|---|
| Geography | `aruba.Region*`, `aruba.Zone*` | `InRegion(...)`, `InZone(...)` |
| Billing | `aruba.BillingPeriod*` | `WithBillingPeriod(...)` |
| Compute | `aruba.CloudServerFlavor*`, `aruba.VolumeImage*` | `OfFlavor(...)`, `FromImage(...)` |
| Storage | `aruba.BlockStorageType*`, `aruba.StorageBackupType*` | `OfType(...)` |
| Database | `aruba.DatabaseEngine*`, `aruba.DBaaSFlavor*` | `OfEngine(...)`, `OfFlavor(...)` |
| Container | `aruba.KubernetesVersion*`, `aruba.NodePoolInstance*`, `aruba.ContainerRegistrySizeFlavor*` | `WithKubernetesVersion(...)`, `OfInstance(...)`, `OfSize(...)` |
| Network | `aruba.SubnetType*`, `aruba.RuleProtocol*`, `aruba.RuleDirection*`, `aruba.EndpointType*`, `aruba.VPNType*`, `aruba.VPNClientProtocol*`, `aruba.IKE*`, `aruba.ESP*` | `OfType(...)`, `WithProtocol(...)`, `WithDirection(...)` |
| Schedule | `aruba.JobType*`, `aruba.HTTPVerb*` | implied by `OneShotAt`/`WithCron`; `WithVerb(...)` |
| Security | `aruba.KeyAlgorithm*`, `aruba.KeyType*`, `aruba.ServiceStatus*` | `OfAlgorithm(...)` |
| Metric | `aruba.ActionType*` | alert setter |

Setter signatures take the alias type, not `string`.

### Adapter conventions

- Constructor: `new<X>ClientAdapter(rest *restclient.Client) *<x>ClientAdapter`. Wired in `pkg/aruba/builder.go`.
- `Create` / `Update` pattern: check `X.Err()` → friendly multi-field nil-ID validation → low-level call → `populateHTTPEnvelope` → `fromResponse` on success → `&HTTPError{...}` on non-2xx.
- `Get` / `Delete` / `List` accept `Ref` (not a typed wrapper) so callers can pass `aruba.URI("/...")` or any `Ref`-satisfying value.
- Family B adapters (Database, Key, Kmip, User, Grant) do the most pre-flight validation because their deep parent chains aren't caught by the API until late; Family A adapters mostly only check `ProjectID() != ""`.

### Wait conventions

- Resources that embed `statusMixin` (`pkg/aruba/mixin_status.go`) get `WaitUntilActive`, `WaitUntilReady`, and `WaitUntilStates(ctx, targets, opts...)` for free.
- Adapters install a `refresh` closure after `Create` / `Get` / `List` — without it, `WaitUntilStates` returns an immediate error.
- Family B resources without `statusMixin` define their own `WaitUntil*` that drive `pkg/async.WaitFor` directly (e.g. `*Kmip.WaitUntilCertificateAvailable` in `resource_kmip.go`).

### URI segment casing

- URI segments use the **exact casing published at <https://api.arubacloud.com/docs/>**. Examples: `securityGroups`, `securityRules`, `loadBalancers`, `keyPairs`, `blockStorages`, `vpnTunnels`, `vpcPeerings`.
- Each resource has **one canonical segment**. `*IDsFromRef` helpers and scoped-mixin lookups use that single form — no fallbacks, no alternative spellings.
- `parseURIIDs` in `pkg/aruba/ref.go` preserves case; never normalise segment keys to lowercase.
- Ref helper templates (`<Resource>Ref(...)`) and `internal/clients/<domain>/path.go` constants must use the same canonical segment so outgoing requests and incoming `Metadata.URI` values remain self-consistent.
