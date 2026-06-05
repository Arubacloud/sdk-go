# Conventions

## Package boundary

- Public API surface lives in `pkg/` — types, interfaces, and the `NewClient` entry point
- Concrete implementations live in `internal/` — not importable by external modules
- Only interfaces and data types are exported; concrete impl structs are always unexported
- **Prefer `aruba.X` over `types.X` in public surface.** When `pkg/aruba` needs to refer to a `pkg/types` enum or struct, alias it in `pkg/aruba/aliases.go` (`type Foo = types.Foo` plus matching constants). This preserves the single-import contract — see `ai/ARCHITECTURE.md` › "Single-import design principle".

## Naming

### Types

| Kind | Suffix | Example |
|---|---|---|
| Input/request struct | `Request` | `CloudServerRequest` |
| Single output struct | `Response` | `CloudServerResponse` |
| Collection output struct | `ListResponse` (embeds `ListResponse`) | `CloudServerListResponse` |
| Struct used on both request AND response sides | `Common` | `LinkedResourceCommon`, `BillingPlanCommon` |
| Service client interface | none (domain name) | `CloudServersClient` |
| Service client implementation | `Impl` (unexported) | `cloudServersClientImpl` |
| Constructor | `New<TypeName>` | `NewCloudServersClientImpl` |

The `*Result` suffix is **not** used — it was a legacy straggler and has been removed. Enum/scalar types (`State`, `Region`, `BillingPeriod`, `RuleProtocol`, …) carry no Request/Response/Common suffix because their wire role is determined by the parent field. See `pkg/types/doc.go` for the canonical in-source statement of this rule.

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

| Connector | Role | Example |
|---|---|---|
| `New<Resource>()` + `.Named(name)` | construct, then name | `NewCloudServer().Named("web-01")` |
| `In<Parent>` / `In<Geo>` | containment & placement | `InProject`, `InVPC`, `InRegion`, `InZone` |
| `Of<Classifier>` | type / sizing class / role | `OfFlavor`, `OfType`, `OfEngine`, `OfSize`, `OfAlgorithm`, `OfInstance`, `OfRole` |
| `From<Source>` | source / origin | `FromImage`, `FromSnapshot`, `FromVolume`, `FromBackup`, `FromDBaaS`, `FromDatabase` |
| `With<Noun>` | accompaniment / attached typed config | `WithVPC`, `WithElasticIP`, `WithCIDR`, `WithKubernetesVersion` |
| participial (`<Verb>ing…`) | active relationship | `BootingFrom`, `UsingKeyPair`, `Targeting`, `PeeredWith` |
| `On<Noun>` | network placement | `OnSubnets` |
| `Tagged` / `Untagged` / `RetaggedAs` | tag-set mutation — variadic, append / remove / replace | `Tagged("prod","sdk")` |
| `SizedGB(int)` / `RetainedForDays(int)` / `DescribedAs(string)` / `BilledBy(BillingPeriod)` | descriptive value phrases | `SizedGB(20)`, `BilledBy(BillingPeriodHour)` |
| no-arg adjective/participle | boolean state | `Enabled()`, `Disabled()`, `HighlyAvailable()` |
| `As<Adj>()` / `Not<Adj>()` + `Is<Adj>()` getter | tri-state `*bool` idiom | `AsDefault`/`NotDefault`/`IsDefault`, `AsBootable`/`NotBootable`/`IsBootable` |
| `With<X>()` / `Without<X>()` | no-arg boolean pair | `WithPreset`/`WithoutPreset`, `WithoutAutoscaling`, `WithoutNodePools` |

**Collection setters** (`OnSubnets`, `WithSecurityGroups`, `WithNodePools`, `WithSteps`, `WithRoutes`, `WithDNSServers`) are variadic and plural-named — a single call can carry one or many items, and repeated calls append. `Replace<Plural>` does wholesale replacement; `Without<Plural>()` clears.

**Important:** bare `Set*` is reserved for non-chainable, side-effecting action methods (`*CloudServer.SetPassword(ctx, pw, opts...)` in `resource_cloud_server.go`) — never use it for a builder setter.

### Canonical chain order

Every builder chain should read like a natural English sentence — imagine an *Eloquent Old Grandma* at a cloud booking office placing an order: "I'd like a CloudServer of flavor Ubuntu 24.04, named 'my-server-01', tagged 'prod', in project Foo, in region Italy, booting from volume V, with VPC X, billed hourly."

Setters cluster into 13 ordered buckets. Within a chain, appear top-to-bottom in this order; skip buckets that don't apply.

| # | Bucket | Setters |
|---|---|---|
| 1 | **noun** | `New<X>()` |
| 2 | **classifier** | `OfFlavor`, `OfType`, `OfEngine`, `OfSize`, `OfAlgorithm`, `OfInstance`, `OfRole` |
| 3 | **name** | `Named` |
| 4 | **labels** | `Tagged(…)` — variadic, **single call** |
| 5 | **containment** | `InProject`, `InVPC`, `InSecurityGroup`, `InDBaaS`, `InDatabase`, `InKMS`, `InVPNTunnel`, `InVPCPeering` |
| 6 | **geography** | `InRegion`, `InZone` |
| 7 | **descriptive scalars** | `SizedGB`, `WithCIDR`, `WithKubernetesVersion`, `WithPodCIDR`, `WithNodeCIDR`, `WithPublicKey`, `WithAdminUsername`, `WithUsername`, `WithPassword`, `WithUserData`, `DescribedAs`, `RetainedForDays`, `WithMaxStorageQuotaGB`, `WithAutoscaling`, `WithCron`, `RecurringUntil`, `OneShotAt`, `WithAction`, `WithVerb`, `WithPort`, `WithProtocol`, `WithDirection`, `WithPeerClientPublicIP` |
| 8 | **origin** | `FromImage`, `FromVolume`, `FromBackup`, `FromSnapshot`, `FromDBaaS`, `FromDatabase`, `BootingFrom` |
| 9 | **attached config** | `WithVPC`, `WithSubnet`, `WithSecurityGroup`, `WithSecurityGroups`, `WithElasticIP`, `WithBlockStorage`, `WithDHCP`, `WithNodePools`, `WithSteps`, `WithIKESettings`, `WithESPSettings`, `WithPSKSettings`, `WithPeerVPC`, `WithTarget` |
| 10 | **network placement** | `OnSubnets` |
| 11 | **active relationship** | `UsingKeyPair`, `Targeting`, `TargetingCIDR`, `ForUser`, `ToVolume` |
| 12 | **boolean / policy state** | `Enabled`/`Disabled`, `HighlyAvailable`, `AsBootable`/`NotBootable`, `AsDefault`/`NotDefault`, `WithPreset`/`WithoutPreset`, `WithoutAutoscaling` |
| 13 | **billing** | `BilledBy(period)` |

**Variadic-collapse rule:** when the same variadic setter would be called twice on the same builder, collapse to a single call — `Tagged("a").Tagged("b")` → `Tagged("a", "b")`. Applies to `Tagged`, `WithDNSServers`, and any other plural-named collection setter.

Per-resource canonical chains are exercised in `examples/all-resources/` — treat those files as the executable specification of the correct ordering.

### Wire-level translation rules

The wrapper API never echoes a wire field name verbatim if a clearer Go-idiomatic name exists:

- **Units explicit in the method name:** `SizedGB(int)` (wire field `size` / JSON `sizeGb`) on `*BlockStorage` and `*DBaaS`.
- **Paired-bool idiom for tri-state `*bool` flags:** `IsDefault()` read + `AsDefault()` / `NotDefault()` write on SecurityGroup, VPC, Subnet, Project — preserves the `nil` / `false` / `true` distinction across the wire.
- **`Ref`-accepting setters hide URI assembly:** `FromVolume(Ref)`, `FromSnapshot(Ref)`, `FromDBaaS(Ref)`, `FromDatabase(Ref)` — an empty URI from the `Ref` goes on the error sink, not on the wire.
- **Sub-builder errors propagate:** `WithIKESettings(*VPNIKE)` drains the sub-builder's `errMixin` into the tunnel's accumulator (`resource_vpn_tunnel.go`).
- **Shape-collapsing:** `*Job.OneShotAt(t)` sets both `JobType` and `ScheduleAt`; `WithCron + RecurringUntil` sets `JobType` + `Cron` + `ExecuteUntil`. Mode conflicts record a setter-time error via `requireMode` in `resource_job.go`.

### Typed-enum aliases (`pkg/aruba/aliases.go`)

Public callers never pass raw strings. Every domain that accepts an enum has typed aliases in `aliases.go`:

| Domain | Alias prefix | Typical setter |
|---|---|---|
| Geography | `aruba.Region*`, `aruba.Zone*` | `InRegion(...)`, `InZone(...)` |
| Billing | `aruba.BillingPeriod*` | `BilledBy(BillingPeriod)` |
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

### Getter taxonomy

Every Family-A wrapper exposes a standard set of read accessors. Callers should always use these instead of reaching into `wrapper.Raw().Properties.X`.

**Audit actor getters** — `CreatedBy()`, `UpdatedBy()`, `CreatedUser()`, `UpdatedUser()` are promoted from `responseMetadataMixin` onto all Family-A wrappers. They return `""` when the server does not populate the field (e.g. after a `Create` whose response omits `updatedBy`). Family-B wrappers whose wire types carry `createdBy` directly (Database, User, Grant) define their own `CreatedBy()` which shadows the mixin.

**Naming rules:**
- Simple scalar field from `Properties`: use the field name directly — `Name()`, `State()`, `Region()`.
- URI of a linked resource: `<LinkedResource>URI()` (e.g. `VPNTunnelURI()`, `AssociatedResourceURI()`).
- Convenience alias for a commonly-read nested field: keep the same name as the logical concept — `CloudSubnetCIDR()` is an alias for `CloudSubnet()` on `VPNRoute`.
- Collection of linked URIs: plural noun — `Subnets() []string`, `SecurityGroups() []string`.

**Response-preferring fallback rule:** when a field can be set by the caller (via a chainable setter) *and* returned by the server, the getter must read the response field first and fall back to the local pointer. This ensures `Get → display` works without the caller having to re-set anything. Boilerplate:
```go
func (r *Resource) SomeField() string {
    if r.response != nil && r.response.Properties.SomeField != "" {
        return r.response.Properties.SomeField
    }
    return derefString(r.someField)
}
```

**When to expose vs. defer to `Raw()`:** expose a getter whenever the field is needed for display, filtering, or round-trip `Get → Update` flows. Fields that are rarely read, deeply nested, or opaque blobs are fine behind `Raw()`. A getter that just wraps `Raw().Properties.X` with no fallback logic adds no value — skip it and document the field location in the type's doc comment instead.

**`fromResponse` round-trip invariant:** after `fromResponse(resp)`, calling `toRequest()` must reproduce every caller-settable field that the server returned. Rehydrate local pointers / slices inside `fromResponse` when the API echoes them in GET responses. If the API does not echo a field (e.g. password, user-data), the getter must document the limitation.

### Wait conventions

- `refreshMixin` (`pkg/aruba/mixin_refresh.go`) owns the `refresh` callback and `WaitUntilGone`. Adapters install the closure after `Create` / `Get` / `List` — without it, any `WaitUntil*` call returns an immediate error.
- `statusMixin` (`pkg/aruba/mixin_status.go`) embeds `refreshMixin` and adds `WaitUntilActive`, `WaitUntilReady`, and `WaitUntilStates(ctx, targets, opts...)`. All Family-A resources embed `statusMixin`.
- Family B resources embed `refreshMixin` directly (no `statusMixin`) and gain `WaitUntilGone` only. Those that need custom polling define their own `WaitUntil*` driving `pkg/async.WaitFor` (e.g. `*Kmip.WaitUntilCertificateAvailable` in `resource_kmip.go`).

### URI segment casing

- URI segments follow the server-canonical **lowerCamelCase** rule:
  - Single-word / acronym collections stay **lowercase**: `vpcs`, `subnets`, `kaas`, `kms`,
    `dbaas`, `backups`, `restores`, `snapshots`, `registries`, `jobs`.
  - Compound collections use **lowerCamelCase**: `cloudServers`, `keyPairs`, `elasticIps`,
    `securityGroups`, `securityRules`, `blockStorages`, `loadBalancers`, `vpnTunnels`,
    `vpcPeerings`, `vpcPeeringRoutes`, `vpnRoutes`.
  - This was verified against the server's `metadata.uri` echoes in a live end-to-end run
    (2026-05-28). Commit `f548a4f` aligned all `internal/clients/*/path.go` constants to this
    convention.
- **Do not revert to all-lowercase.** Downstream provisioners store and re-emit request URIs
  verbatim — a casing change causes silent provisioning failures. The header comment at the top
  of every `internal/clients/*/path.go` repeats this warning.
- Each resource has **one canonical segment**. `*IDsFromRef` helpers and scoped-mixin lookups use
  that single form — no fallbacks, no alternative spellings.
- `parseURIIDs` in `pkg/aruba/ref.go` preserves case; never normalise segment keys to lowercase.
- Ref helper templates (`<Resource>Ref(...)`) and `internal/clients/<domain>/path.go` constants
  must use the same canonical segment so outgoing requests and incoming `Metadata.URI` values
  remain self-consistent.

### Tag format constraint

The server enforces a minimum tag length of **4 characters** per tag (see `pkg/types/error.go`).
Tags shorter than 4 characters return HTTP 400. In `examples/all-resources/` and any test that
supplies inline tag literals, ensure every string in `Tagged(...)` is ≥4 chars.
