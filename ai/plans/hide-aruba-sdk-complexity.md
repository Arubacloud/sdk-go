# Hide Aruba SDK Models Complexity — Wrapper Layer Plan

GitHub issues addressed:
- **#14** *(parent)* — Isolate SDK developer-facing models from the underlying Aruba REST API models.
- **#15** — Move current models to inside the REST Adapter. *Becomes trivial after #16; treated as a follow-up.*
- **#16** — Abstract and fancy models to developer-facing layer. *Bulk of the work. This plan.*

---

## Context

End users of the SDK are exposed to the full REST schema. To set tags on a VPC today they must construct a four-level nested literal (`VPCRequest{Metadata: RegionalResourceMetadataRequest{ResourceMetadataRequest{Name:..., Tags:[...]}, Location: LocationRequest{Value:"..."}}, Properties:...}`), thread `projectID` strings positionally through every call, and after every response fish a `*string` URI out of `*resp.Data.Metadata.URI` to feed the next call. There is no fail-early validation, the `*types.RequestParameters` last-arg is `nil` in every example, and pagination is not abstracted. The REST schema is verbose and confusing without adding value to the SDK consumer.

This plan introduces a developer-facing wrapper layer that:

- presents one fluent type per resource (`*aruba.VPC`, `*aruba.Subnet`, ...) hiding the nested `Metadata`/`Properties` boxing,
- carries hierarchy (`IntoProject(p)`, `IntoVPC(v)`) and body cross-references (`WithVPC(v)`, `AddSubnet(s)`) typed instead of as URI strings,
- unifies request and response into a single round-trippable struct (so `Create(...)` returns a wrapper that can be re-fed into `Update(...)` after `.AddTag(...)`),
- validates as early as possible and surfaces server errors via `wrapper.RawError()` even when `Create` returns an error.

After the wrapper layer is shipped, the public API will no longer reference `pkg/types`, making issue #15 (move types to `internal/`) a mechanical rename PR.

---

## Decisions locked (from /plan Q&A)

| Question | Choice |
|---|---|
| Wrapper data model | **Compose** request + response inside the wrapper. Setters mutate request slots (held by mixins); `fromResponse` writes to the response slot. No flat data duplication. |
| Service-client interface evolution | **Replace** existing `Create/Update/Get/Delete/List` signatures wholesale. Breaking change accepted (alpha + `breaking-change` label). |
| Package layout | **Flat in `pkg/aruba/`**. All wrappers, mixins, and factories live directly in the existing public package. New file naming: `resource_<resource>.go`. |
| Setter naming | Idiomatic Go: `WithName`, `AddTag`, `RemoveTag`, `ReplaceTags`, `IntoProject(p)`, `WithVPC(v)`. (Matches user's example; departs from `Options`'s `WithAdditional*` only for slice-append semantics on user-data fields.) |
| Async polling (`WaitUntilActive`) | **In scope** as a dedicated issue (#async-poll). Layered on `statusMixin` + `pkg/async.WaitFor`; additive surface, lands after scaffolding + Project reference. |
| Path-parent setters | **Single direct-parent setter per wrapper.** A wrapper inherits all ancestor IDs from its direct parent (e.g. `SecurityRule.IntoSecurityGroup(sg)` reads `sg.VPCID()` and `sg.ProjectID()` automatically — no `IntoProject` / `IntoVPC` on `SecurityRule`). |

---

## Architecture

### Wrapper data model (Compose)

Every wrapper is a struct that embeds a small set of unexported mixins (each owning one slot of state) plus resource-specific fields. The wrapper produces a `types.XxxRequest` lazily via `toRequest()` at Create/Update time, and absorbs a `*types.XxxResponse` via `fromResponse()` after a server reply. Read accessors prefer the response slot, falling back to the request slot.

```go
// pkg/aruba/resource_vpc.go (illustrative)
type VPC struct {
    metadataMixin            // Name, Tags  -> owns slots written by WithName/AddTag/...
    regionalMixin            // Location.Value
    projectScopedMixin       // parent project ID set by IntoProject
    responseMetadataMixin    // *types.ResourceMetadataResponse (post-reply)
    statusMixin              // *types.ResourceStatus
    linkedMixin              // []types.LinkedResource
    httpEnvelopeMixin        // status, headers, raw body, *types.ErrorResponse
    errMixin                 // setter-time validation accumulator

    // VPC-specific properties
    defaultVPC *bool
    preset     *bool
}

func (v *VPC) toRequest() types.VPCRequest          // builds nested types from mixin slots
func (v *VPC) fromResponse(r *types.VPCResponse)    // writes mixin response slots
```

### Setters (chainable, return `*Self`)

```go
func VPC() *VPC                                        // factory in factories.go
func (v *VPC) IntoProject(p Ref) *VPC                  // path-parent
func (v *VPC) WithName(name string) *VPC               // metadataMixin
func (v *VPC) AddTag(tag string) *VPC                  // metadataMixin (dedup)
func (v *VPC) RemoveTag(tag string) *VPC               // metadataMixin
func (v *VPC) ReplaceTags(tags ...string) *VPC         // metadataMixin
func (v *VPC) WithLocation(loc string) *VPC            // regionalMixin
func (v *VPC) InRegion(region string) *VPC             // regionalMixin alias
func (v *VPC) WithDefault(b bool) *VPC                 // VPC-specific
func (v *VPC) WithPreset(b bool) *VPC                  // VPC-specific
```

### Read accessors (zero-value safe)

```go
func (v *VPC) ID() string             // post-create; "" until first response
func (v *VPC) URI() string            // satisfies the Ref interface
func (v *VPC) Name() string           // resp Name, falls back to request Name
func (v *VPC) Tags() []string
func (v *VPC) Region() string
func (v *VPC) CreatedAt() time.Time
func (v *VPC) State() string          // statusMixin
func (v *VPC) IsDisabled() bool
func (v *VPC) FailureReason() string
func (v *VPC) LinkedResources() []types.LinkedResource

// Power-user escape hatches:
func (v *VPC) Raw() *types.VPCResponse
func (v *VPC) RawRequest() types.VPCRequest
func (v *VPC) RawError() *types.ErrorResponse
func (v *VPC) RawHTTP() (*http.Response, []byte)
func (v *VPC) StatusCode() int
func (v *VPC) Err() error               // joined setter-time errors
```

### Round-trip example (matches user's snippet)

```go
import "github.com/Arubacloud/sdk-go/pkg/aruba"

project, err := arubaClient.FromProject().Create(
    ctx,
    aruba.Project().WithName("my-first-project").AddTag("tag-1").AddTag("tag-2"),
)
// ...
project, err = arubaClient.FromProject().Update(ctx, project.RemoveTag("tag-1").AddTag("tag-3"))

vpc1, err := arubaClient.FromNetwork().VPCs().Create(
    ctx,
    aruba.VPC().IntoProject(project).WithName("my-vpc-1").AddTag("tag-3"),
)

subnet1, err := arubaClient.FromNetwork().Subnets().Create(
    ctx,
    aruba.Subnet().IntoVPC(vpc1).WithName("my-subnet-1").AddTag("tag-3"),
)
```

---

## Service-client interface evolution

Every existing CRUD method is **replaced**, not added alongside. Existing service group accessors (`FromNetwork()`, `VPCs()`, ...) and naming are preserved.

```go
// Existing today (to be removed):
Create(ctx, projectID string, body types.VPCRequest, params *types.RequestParameters)
       (*types.Response[types.VPCResponse], error)

// New shape:
Create(ctx context.Context, vpc *VPC, opts ...CallOption) (*VPC, error)
Update(ctx context.Context, vpc *VPC, opts ...CallOption) (*VPC, error)
Get   (ctx context.Context, ref Ref, opts ...CallOption) (*VPC, error)
Delete(ctx context.Context, ref Ref, opts ...CallOption) error
List  (ctx context.Context, parent Ref, opts ...CallOption) (*List[*VPC], error)
```

`*types.RequestParameters` is replaced by functional `CallOption`s, with a `WithRawParameters(*types.RequestParameters)` escape hatch:

```go
type CallOption func(*callOptions)
func WithFilter(f string) CallOption
func WithSort(s string) CallOption
func WithLimit(n int) CallOption
func WithOffset(n int) CallOption
func WithProjection(p string) CallOption
func WithAPIVersion(v string) CallOption
func WithRawParameters(p *types.RequestParameters) CallOption
```

`Create`/`Update`/`Get` always return the wrapper non-nil even on HTTP error so `wrapper.RawError()` / `RawHTTP()` / `StatusCode()` remain inspectable. `Delete` returns a typed `*HTTPError` carrying status + raw body.

---

## Cross-resource refs

A single `Ref` interface unifies wrappers and ad-hoc URI strings:

```go
type Ref interface {
    URI() string
    ID() string  // "" if URI-only
}

// Escape hatch in factories.go:
func URI(s string) Ref
```

**Path-parent setter (single, direct parent only — ancestors inherited).** Each child wrapper exposes exactly one `IntoX(parent Ref)` setter for its direct parent in the URL hierarchy. The mixin reads the parent's ancestor IDs by attempting an interface assertion (e.g. `parent.(interface{ VPCID() string })`) and falls back to parsing the URI path when the parent is an opaque `aruba.URI(...)` Ref. Validates each ID non-empty at Create.

| Wrapper | Direct path-parent | Ancestor IDs inherited |
|---|---|---|
| VPC, ElasticIP, LoadBalancer, VPNTunnel, KeyPair, BlockStorage, Snapshot, StorageBackup, DBaaS, KaaS, ContainerRegistry, KMS, Job, CloudServer, DBaaSBackup, AuditEvent, Alert, Metric | `IntoProject(p Ref)` | — |
| Subnet, SecurityGroup, VPCPeering | `IntoVPC(v Ref)` | `projectID` |
| SecurityRule | `IntoSecurityGroup(sg Ref)` | `projectID`, `vpcID` |
| VPCPeeringRoute | `IntoVPCPeering(p Ref)` | `projectID`, `vpcID` |
| VPNRoute | `IntoVPNTunnel(t Ref)` | `projectID` *(VPN routes do not carry a vpcID)* |
| Database, User | `IntoDBaaS(d Ref)` | `projectID` |
| Grant | `IntoDatabase(d Ref)` | `projectID`, `dbaasID` |
| Restore | `IntoBackup(b Ref)` | `projectID` |
| Key, Kmip | `IntoKMS(k Ref)` | `projectID` |
| Project | *(root — no parent setter)* | — |

Ancestor lookup contract — every wrapper exposes accessors for the IDs it knows (its own + all inherited):

```go
// Read by children via interface assertions, no concrete imports needed.
type withProjectID       interface{ ProjectID() string }
type withVPCID           interface{ VPCID() string }
type withSecurityGroupID interface{ SecurityGroupID() string }
type withDBaaSID         interface{ DBaaSID() string }
type withDatabaseID      interface{ DatabaseID() string }
type withVPCPeeringID    interface{ VPCPeeringID() string }
type withVPNTunnelID     interface{ VPNTunnelID() string }
type withBackupID        interface{ BackupID() string }
type withKMSID           interface{ KMSID() string }
```

A `*SecurityGroup` therefore exposes `ProjectID()`, `VPCID()`, `SecurityGroupID()`; a child `*SecurityRule` reads all three from its `IntoSecurityGroup(sg)` parent in one call. For an opaque `aruba.URI(string)` Ref the mixin parses the URI path segments to extract whichever ancestor IDs it needs; if a needed ID is missing from the URI shape, a setter-time error is recorded so it surfaces at Create. This eliminates the possibility of constructing a `SecurityRule` whose declared VPC disagrees with its SecurityGroup's VPC — there is no way to declare that contradiction.

**User-facing example (compare to today's positional-args API):**

```go
// Today:
arubaClient.FromNetwork().SecurityGroupRules().Create(
    ctx, projectID, vpcID, securityGroupID, ruleReq, nil,
)

// New:
rule, err := arubaClient.FromNetwork().SecurityGroupRules().Create(
    ctx,
    aruba.SecurityRule().
        IntoSecurityGroup(sg).             // sg already knows vpcID + projectID
        WithDirection("Ingress").
        WithProtocol("tcp").
        WithPort("443").
        WithTargetCIDR("0.0.0.0/0"),
)
```

**Body-ref setters** read `parent.URI()`, validate non-empty at Create:

| Wrapper | Body-ref setters |
|---|---|
| VPCPeering | `WithRemoteVPC(Ref)` |
| VPNTunnel | nested via `aruba.VPNIPConfig().WithVPC(Ref).WithPublicIP(Ref).WithSubnet(name, cidr)` |
| CloudServer | `WithVPC`, `WithBootVolume`, `WithKeyPair`, `WithElasticIP`, `AddSubnet`, `AddSecurityGroup` |
| Snapshot | `OfVolume(Ref)` |
| BlockStorage | `FromSnapshot(Ref)` (optional) |
| StorageBackup | `WithOrigin(Ref)` |
| Restore | `WithTarget(Ref)` |
| DBaaSBackup | `WithDBaaS(Ref)`, `WithDatabase(Ref)` |
| KaaS | `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroupName(string)` *(name, not URI — schema asymmetry)* |
| ContainerRegistry | `WithPublicIP(Ref)`, `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithBlockStorage(Ref)` |
| Job | nested via `aruba.JobStep().OfResource(Ref).WithAction(...)` |
| SecurityRule | `WithTargetSecurityGroup(Ref)` (sets Kind+Value atomically) / `WithTargetCIDR(string)` |
| DBaaS | `WithNetworking(VPC, Subnet, SecurityGroup, ElasticIP Ref)` *(stored as raw URI strings in body — wrapper unwraps)* |

---

## Reusable mixins

Unexported, embedded by value, in `pkg/aruba/mixins.go`:

| Mixin | State | Promotes |
|---|---|---|
| `metadataMixin` | name, tags | `WithName`, `AddTag`, `RemoveTag`, `ReplaceTags`, `Name()`, `Tags()` |
| `regionalMixin` | location | `WithLocation`, `InRegion`, `Region()` |
| `projectScopedMixin` | projectID | `IntoProject`, `ProjectID()` |
| `vpcScopedMixin` | vpcID + projectID *(inherited from VPC parent)* | `IntoVPC`, `VPCID()`, `ProjectID()` |
| `securityGroupScopedMixin` | securityGroupID + vpcID + projectID *(inherited from SG parent)* | `IntoSecurityGroup`, `SecurityGroupID()`, `VPCID()`, `ProjectID()` |
| `dbaasScopedMixin` | dbaasID + projectID | `IntoDBaaS`, `DBaaSID()`, `ProjectID()` |
| `databaseScopedMixin` | databaseID + dbaasID + projectID | `IntoDatabase`, `DatabaseID()`, `DBaaSID()`, `ProjectID()` |
| `backupScopedMixin` | backupID + projectID | `IntoBackup`, `BackupID()`, `ProjectID()` |
| `kmsScopedMixin` | kmsID + projectID | `IntoKMS`, `KMSID()`, `ProjectID()` |
| `vpnTunnelScopedMixin` | vpnTunnelID + projectID | `IntoVPNTunnel`, `VPNTunnelID()`, `ProjectID()` |
| `vpcPeeringScopedMixin` | vpcPeeringID + vpcID + projectID | `IntoVPCPeering`, `VPCPeeringID()`, `VPCID()`, `ProjectID()` |
| `responseMetadataMixin` | `*types.ResourceMetadataResponse` | `ID`, `URI`, `Project`, `CreatedAt`, `UpdatedAt`, `Version` |
| `statusMixin` | `*types.ResourceStatus` + refresh callback | `State`, `IsDisabled`, `DisableReasons`, `FailureReason`, `PreviousState`, `WaitUntilActive(ctx, opts ...WaitOption) error`, `WaitUntilState(ctx, state, opts...) error` |
| `linkedMixin` | `[]types.LinkedResource` | `LinkedResources()` |
| `httpEnvelopeMixin` | status, headers, raw body, `*types.ErrorResponse` | `RawHTTP`, `RawError`, `StatusCode`, `Headers` |
| `errMixin` | `errs []error` | `Err()`, internal `addErr`, `validateAll` |

Mixin name-collision risk in the shared `pkg/aruba` package: mitigated by keeping all mixin types **unexported** (lowercase). The exported `Options.With*` methods do not collide with promoted methods because the wrapper types are distinct (`*VPC.WithName` vs `*Options.WithToken`).

---

## Per-resource wrapper inventory (31 wrappers)

Family A = `Metadata`/`Properties` request shape; Family B = flat request; Family C = read-only (no Create/Update).

| # | Resource | File | Family | Mixins (in addition to errMixin + httpEnvelopeMixin) | Notes |
|---|---|---|---|---|---|
| 1 | Project | `resource_project.go` | A non-regional | metadata, projectScoped (self), responseMetadata | + `WithDescription`, `WithDefault` |
| 2 | VPC | `resource_vpc.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithDefault`, `WithPreset` |
| 3 | Subnet | `resource_subnet.go` | A regional | metadata, regional, vpcScoped *(inherits projectID)*, responseMetadata, status, linked | + DHCP sub-builder, `WithType`, `WithCIDR` |
| 4 | ElasticIP | `resource_elastic_ip.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithBillingPeriod` |
| 5 | LoadBalancer | `resource_load_balancer.go` | C read-only | responseMetadata, status, linked | List + Get only |
| 6 | SecurityGroup | `resource_security_group.go` | A non-regional | metadata, vpcScoped *(inherits projectID)*, responseMetadata, status, linked | + `WithDefault` |
| 7 | SecurityRule | `resource_security_rule.go` | A regional | metadata, regional, securityGroupScoped *(inherits projectID, vpcID)*, responseMetadata, status | + `WithDirection`, `WithProtocol`, `WithPort`, `WithTargetCIDR` / `WithTargetSecurityGroup` |
| 8 | VPCPeering | `resource_vpc_peering.go` | A non-regional | metadata, vpcScoped *(inherits projectID)*, responseMetadata, status | + `WithRemoteVPC` |
| 9 | VPCPeeringRoute | `resource_vpc_peering_route.go` | A non-regional | metadata, vpcPeeringScoped *(inherits projectID, vpcID)*, responseMetadata, status | + `WithLocalCIDR`, `WithRemoteCIDR`, `WithBillingPeriod` |
| 10 | VPNTunnel | `resource_vpn_tunnel.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `aruba.VPNIPConfig()`, `aruba.IKESettings()`, `aruba.ESPSettings()`, `aruba.PSKSettings()` sub-builders |
| 11 | VPNRoute | `resource_vpn_route.go` | A regional | metadata, regional, vpnTunnelScoped *(inherits projectID)*, responseMetadata, status | + `WithCloudSubnet`, `WithOnPremSubnet` |
| 12 | CloudServer | `resource_cloud_server.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + body-refs + actions: `PowerOn(ctx)`, `PowerOff(ctx)`, `SetPassword(ctx, pwd)` |
| 13 | KeyPair | `resource_key_pair.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithPublicKey` (no Update) |
| 14 | BlockStorage (Volume) | `resource_block_storage.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithSize`, `WithType`, `FromSnapshot`, `WithImage`, `WithBootable`, `WithBillingPeriod` |
| 15 | Snapshot | `resource_snapshot.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `OfVolume`, `WithBillingPeriod` |
| 16 | StorageBackup | `resource_storage_backup.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithType`, `WithOrigin`, `WithRetentionDays`, `WithBillingPeriod` |
| 17 | Restore | `resource_restore.go` | A regional | metadata, regional, backupScoped *(inherits projectID)*, responseMetadata, status, linked | + `WithTarget` |
| 18 | DBaaS | `resource_dbaas.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithEngine`, `WithFlavor`, `WithStorage`, `WithNetworking(VPC, Subnet, SG, EIP)`, `WithAutoscaling`, `WithBillingPeriod` |
| 19 | Database | `resource_database.go` | B flat | dbaasScoped *(inherits projectID)*, responseMetadata | bespoke `WithName` (no metadataMixin) |
| 20 | User | `resource_user.go` | B flat | dbaasScoped *(inherits projectID)*, responseMetadata | + `WithUsername`, write-only `WithPassword` |
| 21 | Grant | `resource_grant.go` | B flat | databaseScoped *(inherits projectID, dbaasID)*, responseMetadata | + `WithUserName`, `WithRoleName` |
| 22 | DBaaSBackup | `resource_dbaas_backup.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + body-refs `WithDBaaS`, `WithDatabase`, `WithBillingPeriod` (no Update). *Path-parent is project; DBaaS+Database are body-refs per the REST shape.* |
| 23 | KaaS | `resource_kaas.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `aruba.NodePool()` sub-builder, `WithKubernetesVersion`, `WithVPC`, `WithSubnet`, `WithSecurityGroupName`, `WithNodeCIDR`, `WithPodCIDR`, `WithHA`, `WithStorage`, `WithBillingPeriod`, `WithIdentity`, `WithAPIServerAccessProfile`. Update emits `KaaSUpdateRequest`. + `DownloadKubeconfig(ctx)` action |
| 24 | ContainerRegistry | `resource_container_registry.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + body-refs, `WithAdminUsername`, `WithSize` (`concurrentUsers`) |
| 25 | KMS | `resource_kms.go` | A regional | metadata, regional, projectScoped, responseMetadata, status, linked | + `WithBillingPeriod`. KMSClient retains `Keys()`/`Kmips()` accessors. |
| 26 | Key | `resource_key.go` | B flat | kmsScoped *(inherits projectID)*, responseMetadata | + `WithName`, `WithAlgorithm` (no Update) |
| 27 | Kmip | `resource_kmip.go` | B flat | kmsScoped *(inherits projectID)*, responseMetadata | + `WithName` + `Download(ctx)` action (no Update) |
| 28 | Job | `resource_job.go` | A regional | metadata, regional, projectScoped, responseMetadata, status | + `WithEnabled`, `OneShotAt`, `RecurringUntil`, `WithCron`, `aruba.JobStep()` sub-builder |
| 29 | AuditEvent | `resource_audit_event.go` | C read-only | responseMetadata + event-specific fields | List only |
| 30 | Alert | `resource_alert.go` | C read-only | responseMetadata, status + alert-specific | List only. Wrapper accessors normalize the upstream `Theshold`/`ThesholdExceedence` typos to `Threshold()`/`ThresholdExceedance()` (json tags untouched). |
| 31 | Metric | `resource_metric.go` | C read-only | responseMetadata + datapoints | List only |

Sub-builders (also live in `pkg/aruba/`, separate files):
- `resource_subnet_dhcp.go` — `aruba.SubnetDHCP().Enabled().WithRange(start, count).AddRoute(addr, gw).AddDNS(ip)`
- `resource_vpn_ipconfig.go`, `resource_vpn_ike.go`, `resource_vpn_esp.go`, `resource_vpn_psk.go`
- `resource_kaas_nodepool.go` — `aruba.NodePool().Named(...).OfInstance(...).InRegion(...).WithCount(n).WithAutoscaling(min, max)`
- `resource_job_step.go` — `aruba.JobStep().OfResource(Ref).WithAction(uri).WithVerb(...).WithBody(string)`

---

## `List[T]` wrapper (single generic implementation)

```go
type Wrapper interface {
    URI() string
    ID() string
}

type List[T Wrapper] struct {
    items []T
    total int64
    self, prev, next, first, last string
    callerOpts []CallOption
    refetch    func(ctx context.Context, url string) (*List[T], error)
    raw        any  // *types.Response[types.VPCList] etc.
}

func (l *List[T]) Items() []T
func (l *List[T]) Total() int64
func (l *List[T]) HasNext() bool
func (l *List[T]) HasPrev() bool
func (l *List[T]) Next(ctx context.Context) (*List[T], error)
func (l *List[T]) Prev(ctx context.Context) (*List[T], error)
func (l *List[T]) First(ctx context.Context) (*List[T], error)
func (l *List[T]) Last(ctx context.Context) (*List[T], error)
func (l *List[T]) Cursor() (next, prev string)
func (l *List[T]) Raw() any
func (l *List[T]) All(ctx context.Context, yield func(T) bool) error  // iterate all pages
```

Pagination prefers server-supplied URLs (`Next`/`Prev`) when populated, falls back to `Offset`+`Limit` from the original `CallOption`s. Single generic implementation reused by all 31 resources.

---

## Validation strategy

Three error tiers:

| Tier | Detected | Stored | Surfaced |
|---|---|---|---|
| **Setter-time** (`WithName("")`, `IntoProject(nil)`, conflicting region) | inside the setter | `errMixin.errs` | At `Create`/`Update`: returned via `errors.Join`. Also queryable via `wrapper.Err()` before calling. |
| **Pre-call** (parent IDs non-empty, required body fields, body-ref URI shape sanity) | inside service-client method, before HTTP | wrapped error | Returned as `error` from `Create`/`Update`. |
| **HTTP** (4xx/5xx, network error, JSON parse, `MetadataValidationError` from `Validate()`) | after `DoRequest` | `httpEnvelopeMixin` (status, body, `*ErrorResponse`) | Returned as `error`; **wrapper still returned non-nil** so `RawError()`/`RawHTTP()`/`StatusCode()` are inspectable. |

Re-uses existing validators from `pkg/types/utils.go`: `ValidateProject`, `ValidateProjectAndResource`, `ValidateDBaaSResource`, `ValidateDatabaseGrant`, `ValidateVPCResource`, `ValidateSecurityGroupRule`, `ValidateVPCPeeringRoute`, `ValidateVPNRoute`, `ValidateStorageRestore`.

---

## File layout (final, flat in `pkg/aruba/`)

```
pkg/aruba/
  aruba.go              (today)
  builder.go            (today)
  client.go             (today)
  options.go            (today, untouched)
  assertions_test.go    (today, signatures updated to wrappers)

  audit.go              MODIFIED — flip EventsClient signatures
  compute.go            MODIFIED — flip CloudServers/KeyPairs
  container.go          MODIFIED — flip KaaS/ContainerRegistry
  database.go           MODIFIED — flip DBaaS/Databases/Backups/Users/Grants
  metric.go             MODIFIED — flip Alerts/Metrics
  network.go            MODIFIED — flip 10 sub-clients
  project.go            MODIFIED — flip ProjectClient
  schedule.go           MODIFIED — flip Jobs
  security.go           MODIFIED — flip KMS/Keys/Kmips (preserve KMSClientWrapper Keys()/Kmips())
  storage.go            MODIFIED — flip Snapshots/Volumes/Backups/Restores

  factories.go          NEW — aruba.Project(), aruba.VPC(), aruba.Subnet(), ..., aruba.URI(s),
                              aruba.NodePool(), aruba.JobStep(), aruba.VPNIPConfig(), ...
  call_options.go       NEW — CallOption, WithFilter/WithSort/WithLimit/.../WithRawParameters
  ref.go                NEW — Ref interface + uriRef
  list.go               NEW — generic List[T]
  mixins.go             NEW — all mixin structs (~16 mixins)
  errors.go             NEW — typed wrapper errors (HTTPError, ValidationError)

  resource_project.go            NEW (×31 resource files, naming per the table above)
  resource_vpc.go                NEW
  resource_subnet.go             NEW
  resource_subnet_dhcp.go        NEW (sub-builder)
  resource_elastic_ip.go         NEW
  resource_load_balancer.go      NEW
  resource_security_group.go     NEW
  resource_security_rule.go      NEW
  resource_vpc_peering.go        NEW
  resource_vpc_peering_route.go  NEW
  resource_vpn_tunnel.go         NEW
  resource_vpn_ipconfig.go       NEW (sub-builder)
  resource_vpn_ike.go            NEW (sub-builder)
  resource_vpn_esp.go            NEW (sub-builder)
  resource_vpn_psk.go            NEW (sub-builder)
  resource_vpn_route.go          NEW
  resource_cloud_server.go       NEW
  resource_key_pair.go           NEW
  resource_block_storage.go      NEW
  resource_snapshot.go           NEW
  resource_storage_backup.go     NEW
  resource_restore.go            NEW
  resource_dbaas.go              NEW
  resource_database.go           NEW
  resource_user.go               NEW
  resource_grant.go              NEW
  resource_dbaas_backup.go       NEW
  resource_kaas.go               NEW
  resource_kaas_nodepool.go      NEW (sub-builder)
  resource_container_registry.go NEW
  resource_kms.go                NEW
  resource_key.go                NEW
  resource_kmip.go               NEW
  resource_job.go                NEW
  resource_job_step.go           NEW (sub-builder)
  resource_audit_event.go        NEW
  resource_alert.go              NEW
  resource_metric.go             NEW
```

Each `resource_*.go` is paired with `resource_*_test.go` covering: setter chaining, error accumulation, `toRequest()` round-trip, `fromResponse()` population, `URI()`/`ID()` accessors satisfying `Ref`, ancestor-ID inheritance from the direct-parent setter (including the `aruba.URI(string)` Ref path-parsing fallback).

---

## Migration sequence (~17 PRs)

1. **Scaffolding PR** — `mixins.go`, `ref.go`, `list.go`, `call_options.go`, `errors.go`, `factories.go` (URI factory only). Unit tests via a test-only fake wrapper. Includes URI-path parsing helpers used by the ancestor-inheritance fallback.
2. **Project reference PR** — `resource_project.go` + `resource_project_test.go`, factory `aruba.Project()`, flip `pkg/aruba/project.go` interface, adapt the existing `internal/clients/project/project.go` consumer (preserving the `MetadataValidationError` contract). Add a `docs/wrappers.md` walking through Project end-to-end.
3. **Network domain PR** — 10 wrappers + interface flip in `network.go`. Largest blast-radius first to expose pattern friction (especially the ancestor-inheritance design through SecurityRule and VPCPeeringRoute).
4. **Compute PR** — CloudServer (with action methods) + KeyPair.
5. **Storage PR** — BlockStorage, Snapshot, StorageBackup, Restore.
6. **Database PR** — DBaaS (with networking sub-builder), Database, User, Grant, DBaaSBackup.
7. **Container PR** — KaaS (with NodePool sub-builder + `KaaSUpdateRequest` divergence), ContainerRegistry.
8. **Security PR** — KMS + Key + Kmip; preserve `KMSClientWrapper.Keys()/Kmips()`.
9. **Schedule PR** — Job + JobStep sub-builder.
10. **Audit PR** — read-only AuditEvent.
11. **Metric PR** — read-only Alert (normalize `Theshold` typo at accessor) + Metric.
12. **Async polling PR (`WaitUntilActive`)** — extend `statusMixin` with `WaitUntilActive(ctx, opts...)` and `WaitUntilState(ctx, state, opts...)`, layered on `pkg/async.WaitFor`. Each statused wrapper gains a per-resource `Refresh(ctx)` injected at construction time so the mixin can re-poll without a circular import. Per-resource terminal-state map (e.g. CloudServer terminals: `Active`, `Error`; KaaS terminals: `Active`, `Error`; Snapshot terminals: `Available`, `Error`). Lands after all per-domain PRs so every wrapper benefits at once.
13. **Examples PR** — rewrite `cmd/example/{main,update,delete,multitenancy}.go` to use wrappers end-to-end. Replaces the hand-rolled `time.Sleep(5 * time.Second)` polling loops in `createKaaS` and `downloadKmipCertificate` with `wrapper.WaitUntilActive(ctx)`.
14. **Docs PR** — update `docs/website/docs/{intro,resources,response-handling,filters,types,multitenancy}.md` to lead with wrappers.
15. **Cleanup PR** — remove dead helpers (`stringPtr`/`boolPtr`/`int32Ptr` redefined in examples; any unused `*Result` types).
16. **Release PR** — `CHANGELOG.md`, alpha version bump.
17. **#15 follow-up PR** — move `pkg/types/` into `internal/types/`. Mechanical rename + `goimports` once all public references are gone.

Per-domain PRs may be split further if review burden is too high.

---

## Risks and trade-offs

1. **Maintenance scale (~30 wrappers).** Mitigated by mixins — Family A regional resources should land in ~50 lines of resource-specific code each. Code generation deferred unless mixins prove insufficient (the schema variation across A/B/C families argues against generation today).
2. **Schema asymmetries** in `pkg/types` need bespoke `fromResponse` handling: DBaaS networking (raw URI strings vs `ReferenceResource`); CloudServer `FlavorName` vs `Flavor` object; KaaS NodePool `Instance`/`DataCenter` (string vs object) and JSON-tag mismatches (`nodesPool`/`nodePools`, `nodecidr`/`nodeCidr`); SecurityRule.Target Kind+Value atomicity. Each is documented in its per-resource issue.
3. **Flatness vs nested config tension** — KaaS NodePool, VPN IKE/ESP/PSK/IPConfig, DBaaS Networking, Job Step, Subnet DHCP all need sub-builders. Each sub-builder is itself flat-fluent so the user-facing pattern stays uniform.
4. **`aruba.URI(string)` foot-gun** — passing a wrong-shaped URI to `WithVPC(Ref)` produces a server 422. Mitigation: setter-time prefix-check against expected URI shape per resource type, fail fast with a clear message.
5. **Read-only resources** (LoadBalancer, AuditEvent, Alert, Metric) skip `metadataMixin`/`errMixin` and expose `Get`/`List` only. Resist the temptation to add fake setters for "consistency."
6. **Async polling — in scope as its own issue.** `WaitUntilActive(ctx, opts ...WaitOption) error` and `WaitUntilState(ctx, state, opts ...WaitOption) error` land on `statusMixin`, layered on `pkg/async.WaitFor`. Each statused wrapper accepts a `Refresh(ctx)` callback at construction time so the mixin can re-poll without depending on the service-client packages. Risks to track: per-resource terminal-state list must be authored case by case; default backoff (10s × 10 attempts, matching `pkg/async.DefaultWaitFor`) may need tuning once we measure real provisioning times.

---

## Verification (end-to-end)

After each domain PR:

- `go build ./...` and `go vet ./...` must pass.
- `go test ./pkg/aruba/... ./internal/clients/<domain>/...` must pass with the new wrapper-based interface.
- `pkg/aruba/assertions_test.go` (compile-time interface guard) must compile cleanly against the flipped interfaces.
- Hand-run the rewritten `cmd/example/main.go -mode=create -clientID=... -clientSecret=...` against a sandbox tenant; verify each step's wrapper carries `ID()`, `URI()`, `State() == "Active"` and that `Update` accepts the returned wrapper after `.AddTag(...)`.
- Ancestor inheritance: `aruba.SecurityRule().IntoSecurityGroup(sg)...Create(ctx)` succeeds without any explicit `IntoVPC` / `IntoProject`; assert `rule.SecurityGroupID() == sg.SecurityGroupID()`, `rule.VPCID() == sg.VPCID()`, `rule.ProjectID() == sg.ProjectID()`. Repeat with `aruba.URI("/projects/p/network/vpcs/v/security-groups/s")` Ref to exercise the URI-path-parsing fallback.
- Async polling: `vpc, err := vpcs.Create(ctx, ...); err = vpc.WaitUntilActive(ctx)` resolves once the server reports `State() == "Active"`; deliberately race against a known-slow KaaS provision to confirm the timeout/retry behavior matches `pkg/async.DefaultWaitFor`.
- For `List`: pagination via `.Next(ctx)` walks at least two pages on a tenant with >`Limit` resources.
- Error-path: deliberately send `aruba.VPC().IntoProject(p).WithName("")` to `Create`; assert `(*VPC, error)` returns the wrapper non-nil with `.Err()` populated and **no** HTTP call made.
- Error-path: hit a 4xx (e.g. duplicate name); assert `Create` returns `(*VPC, error)` where `vpc.RawError().Title != ""` and `vpc.StatusCode() == 4xx`.

Final regression: re-run all existing internal-client tests under `internal/clients/...` — they keep using `pkg/types` directly and must remain green throughout.
