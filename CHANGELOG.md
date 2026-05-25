# Changelog

All notable changes to this project are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## Versioning & Branch Policy

| Branch | Version series | Status |
|--------|----------------|--------|
| `main` | **v0.2.x** | Active development — new features, ongoing releases |
| `legacy` | **v0.1.x** | Maintenance — bug fixes and security patches only |

- **`main` — v0.2.x (current)**: Introduces the `pkg/aruba/` wrapper
  layer — a major, breaking redesign of the public API surface. Adopters
  upgrading from v0.1.x should expect compile-time breakage; see the
  [v0.2.0](#020--2026-05-13) section below for a full
  migration summary.
- **`legacy` — v0.1.x (maintenance)**: Supported for **6 months** after
  the v0.2.0 release date with bug-fix and security-patch releases only
  (tagged `v0.1.29`, `v0.1.30`, …). No new features will be backported.
  Once the support window closes the `legacy` branch will be archived.

---

## [Unreleased]

### Added

- `List[T]` now embeds `httpEnvelopeMixin`, exposing the same HTTP-envelope
  accessors as single-resource wrappers — `StatusCode()`, `Headers()`,
  `RawHTTP()`, `RawError()` — so list responses are fully introspectable on a
  par with `Create`/`Get`/`Update` calls. (#298)
- `types.ListResponse.BaseList()` returns the embedded pagination/total
  metadata; promoted automatically onto every `*types.XxxList` via Go's
  method-promotion rules. (#298)
- `RawJSON() []byte` / `RawYAML() []byte` on every resource wrapper (31
  resources including `KmipCertificate`) and on `List[T]`. Convenience
  marshalers for the typed response payload returned by `Raw()`. Return `nil`
  when the wrapper has no payload. YAML output uses `gopkg.in/yaml.v3`
  (promoted from indirect to direct dependency).
- New documentation page **Working at Low Level** (EN + IT) collecting the
  residual cases that require importing `pkg/types` or `pkg/async`:
  non-promoted wire fields, structured validation errors, `LinkedResources()`
  traversal, and concurrent/custom polling. The `pkg/async` deep-dive moved
  here from `async.md`. Existing `response-handling.md` examples updated to
  use single-import equivalents (`vpcList.Total()`, `RawJSON()`, etc.).

### Fixed

- `List[T].Raw()` previously returned the full `*types.Response[XxxList]`
  envelope, which contains `*http.Response` whose `GetBody` is a function
  field — `json.Marshal(list.Raw())` failed with
  `json: unsupported type: func() (io.ReadCloser, error)` for every resource
  list. `Raw()` now returns only the JSON-safe wire payload (`*types.XxxList`,
  i.e. `resp.Data`). (#298)
- `SecurityGroupsClient.Get`/`Update`/`Delete` rejected every ref produced by
  `aruba.SecurityGroupRef(...)` with `cannot determine security group ID from Ref "..."`,
  because the lookup expected a hyphenated `security-groups` segment that the Ref never
  emits and the live API never returns. The lookup now uses the documented `securityGroups`
  segment as the single canonical form. (#297)

### Changed

- URI path segments across all Ref helpers, internal client paths, and ref-parsing helpers
  now match the canonical casing documented at <https://api.arubacloud.com/docs/>. Affected
  segments: `securityGroups` (was `securitygroups`), `securityRules` (was `securityrules`),
  `loadBalancers` (was `loadbalancers`), `keyPairs` (was `keypairs`), `blockStorages` (was
  `blockstorages`). The Aruba API is case-insensitive on path routing, so outgoing requests
  behave identically; this is purely an alignment to the published API spec. Hyphenated
  test-only fallbacks (`security-groups`, `security-rules`, `vpn-tunnels`, `peerings`) and
  duplicate-form lookups have been removed — each resource now has a single canonical URI
  segment.

---

## [0.2.3] — 2026-05-22

> **Minor release.** Adds `WaitUntilGone` teardown polling and a typed `State` API;
> internal mixin and file-layout restructuring.

### Added

- `WaitUntilGone(ctx, opts...)` on every deletable, pollable resource wrapper — blocks
  until `Get` returns HTTP 404. Backed by a new shared `refreshMixin` that owns the
  adapter `refresh` callback (embedded by `statusMixin` and the Family-B pollable
  resources `Kmip`, `Grant`, `Database`, `User`, `Key`). Accepts the same `WaitOption`s
  as `WaitUntilReady`.
- Typed `types.State` with overlapping predicates (`IsActive`, `IsTerminal`, …) and
  exported `State*` constants (`StateActive`, `StateRunning`, `StateStopped`, …),
  re-exported via `pkg/aruba`. Replaces per-resource string `terminalStates` maps.

### Changed ⚠️ Breaking

- `WaitUntilStates` signature changed from `[]string` to `[]types.State`. Callers
  passing string slices must switch to the `types.State*` constants.
- `pkg/types` convenience pointer helpers removed (`BoolPtr`, `StringPtr`, `IntPtr`,
  `Int32Ptr`, `Int64Ptr`, `Float64Ptr`, `StatePtr`). Use `ptr.To(...)` from
  `k8s.io/utils/ptr` instead. The `[0.2.1]` note referencing `types.BoolPtr(true)` for
  `Job.WithEnabled(false)` should be updated to `ptr.To(true)`.

### Docs

- `WaitUntilGone` documented in the walkthrough (§6 teardown), async, and resources
  guides — EN and IT. The walkthrough no longer hand-rolls a `pkg/async` helper.
- EN/IT doc coherence fixes, versioning structure, and tooling-path updates.
- `ai/` architecture artifacts updated for the mixin split, typed State, `client_` file
  rename, and `refreshMixin`.

### Internal

- `mixins.go` decomposed into `mixin_common.go`, `mixin_scoped.go`, `mixin_status.go`,
  and `mixin_refresh.go`.
- Resource client files renamed under a `client_` prefix; `New*` factories co-located
  with their resource types.
- `examples/all-resources` drops its hand-rolled `waitUntilGone` helper in favour of the
  wrapper method value `x.WaitUntilGone`.
- `refreshMixin` unit tests and adapter `WaitUntilGone` tests added.

---

## [0.2.2] — 2026-05-22

> **Patch release.** Migrates billing configuration from the legacy `billingPeriod`
> string field to a structured `billingPlan` object across all billed resources.

### Changed

- `billingPeriod` wire field replaced by a structured `billingPlan` object across
  `CloudServer`, `ElasticIP`, `KaaS`, `ContainerRegistry`, `VPNTunnel`, `VPCPeeringRoute`,
  and `DBaaSBackup`. New shared `types.BillingPlan` wire type; `DBaaS` uses it directly.
- `DBaaSBillingPlan` alias removed — use `types.BillingPlan` directly.
- `ElasticIP` lowercase billing-period translator removed (superseded by the `billingPlan`
  object).

### Fixed

- DBaaS `User` password is now base64-encoded at the wire boundary, matching the API
  contract.

### Internal

- `examples/all-resources` updated for the new billing fields (billing period Hour, DBaaS
  flavor refreshed); KMS example uses Month; security-rule egress example drops the port.
  gofmt pass over the billing changes.

---

## [0.2.1] — 2026-05-16

> **Patch release.** Stabilises the v0.2.0 wrapper layer based on
> migration feedback from [acloud-cli PR #111](https://github.com/Arubacloud/acloud-cli/pull/111).
> 20 Pre-Live milestone issues resolved.
>
> **Narrow source-incompatible change:** `types.JobPropertiesRequest.Enabled`
> changed from `bool` to `*bool` (required to fix #282 — `omitempty` silently
> dropped `false`). Code that constructs `types.JobPropertiesRequest` literals
> directly (`Enabled: true`) must be updated to `Enabled: types.BoolPtr(true)`.
> Callers using the `pkg/aruba` wrapper layer (`job.WithEnabled(bool)`) are
> unaffected.
> 20 Pre-Live milestone issues resolved. No breaking changes.

### Added

- Typed network `Ref` constructors for all 10 network resources:
  `VPCRef`, `SubnetRef`, `SecurityGroupRef`, `SecurityRuleRef`,
  `ElasticIPRef`, `LoadBalancerRef`, `VPCPeeringRef`,
  `VPCPeeringRouteRef`, `VPNTunnelRef`, `VPNRouteRef` — eliminates
  the need for downstream consumers to hand-build URI strings. (#268)
- `*CloudServer.WithBillingPeriod(period)` fluent setter and matching
  `BillingPeriod` field on `types.CloudServerPropertiesRequest`. (#267)
- `*KaaS.WithSecurityGroupName(name)` convenience for callers that have
  only the SG name (e.g. a CLI flag value). (#278)
- `*KaaS.ClearNodePools()`, `ReplaceNodePools(...)`, `SetNodePools(...)`
  for explicit node-pool control during Update round-trips. (#279)
- `*DBaaSBackup.InZone(zone)` setter for zone selection within a region,
  independent of the region-derived default. (#275)
- `KubernetesVersion1313` re-added as a **deprecated** alias of
  `KubernetesVersion1323`. Will be removed in v0.3.0. (#280)

### Fixed

- `List[T].Next/Prev/First/Last` now follows server-supplied pagination
  links across all 31 resource adapters. Previously returned a "not yet
  wired" stub error. (#269, #271, #277, #281, #283)
- `vpcsClientAdapter.Get` and `List` now backfill `projectID` from the
  caller-supplied `Ref`, so subsequent `Update` calls succeed without a
  defensive `IntoProject` workaround. (#284)
- `projectsClientImpl.Delete` now parses the error response body, so
  `*HTTPError.ErrResp` carries structured error details on failed
  deletes. (#285)
- `Job.WithEnabled(false)` now actually disables a job via Update. The
  wire field `JobPropertiesRequest.Enabled` changed from `bool` to
  `*bool` so `false` is no longer dropped by `omitempty`. (#282)
- `Snapshot` and `DBaaSBackup` terminal state maps now use `"Active"` to
  match the wire form, consistent with other storage resources.
  `WaitUntilActive` no longer times out spuriously on these resources.
  (#270, #274)
- `DBaaS.fromResponse` now back-populates autoscaling fields from the
  server response, so `Get → Update` round-trips preserve previously
  configured autoscaling settings. (#276)
- `BlockStorage.Update` no longer silently sends `bootable=false` when
  the caller never called `SetBootable`/`UnsetBootable`; the field is
  omitted from the PUT body when not explicitly set. (#272)

### Docs

- `StorageRestoreClient.Update` now carries a godoc note explaining that
  platform support for PUT on restore resources is unverified; callers
  should prefer Create+Delete workflows. (#273)
- `KaaS.WithSecurityGroup` documents the `*SecurityGroup` requirement
  and points to `WithSecurityGroupName` for name-only callers. (#278)
- CHANGELOG retro-note: `KubernetesVersion1313` was removed in v0.2.0
  (Kubernetes 1.31.3 left the live catalog); a deprecated alias is
  restored in v0.2.1 and will be removed in v0.3.0. (#280)

### Closed without code change

- #45 — README async import alias. The current README no longer imports
  `pkg/async` in its example; original concern not reproducible. Closed
  as stale.

### Deferred to v0.3.0

- #97 — REST-compliant per-operation request/response model split.
- #110 — Query `FilterBuilder` ergonomics feature.

---

## [0.2.0] — 2026-05-13

> **Breaking change release.** All service-client CRUD interfaces have
> been replaced with wrapper-based signatures. See the Changed and
> Removed sections below before upgrading.

### Added

- `pkg/aruba/` wrapper layer with **31 fluent resource types**:
  `CloudServer`, `BlockStorage`, `Snapshot`, `StorageBackup`,
  `StorageRestore`, `KaaS`, `ContainerRegistry`, `KMS`, `Key`, `Kmip`,
  `DBaaS`, `Database`, `User`, `Grant`, `DBaaSBackup`, `KeyPair`,
  `LoadBalancer`, `SecurityGroup`, `SecurityRule`, `Subnet`, `VPC`,
  `VPCPeering`, `VPCPeeringRoute`, `VPNTunnel`, `VPNRoute`, `ElasticIP`,
  `Job`, `AuditEvent`, `Alert`, `Metric`, `Project`.
- Fluent factory functions for every resource (`aruba.NewKaaS()`,
  `aruba.NewVPC()`, `aruba.NewDBaaS()`, …).
- `Ref` interface and `aruba.URI(s)` escape hatch for raw resource
  references when a wrapper is not yet available.
- `List[T]` generic paginated-collection type.
- `CallOption` functional options replacing `*types.RequestParameters`.
- Async polling on all statused resources:
  `WaitUntilActive`, `WaitUntilStates`, `WaitUntilReady`,
  `WaitUntilNotUsed`, `WaitUntilUsed`.
- Typed enum constants in `pkg/types` (re-exported via `pkg/aruba`):
  `BillingPeriod`, `Region`, `Zone`, `CloudServerFlavor`, `DBaaSFlavor`,
  `DatabaseEngine`, `KubernetesVersion`, `NodePoolInstance`, `VolumeImage`,
  `SubnetType`, `KeyAlgorithm`, `RuleProtocol`, IKE/ESP crypto aliases
  (`IKEEncryption`, `IKEHash`, `IKEDHGroup`, `ESPEncryption`, `ESPHash`,
  `PFSGroup`), `VPNType`, `VPNClientProtocol`, `HTTPVerb`,
  `ContainerRegistrySizeFlavor`.
- `KubernetesVersion1341` constant — Kubernetes 1.34.1.
- `VolumeImage` constants for Aruba's stock BlockStorage templates.
- `HTTPVerb` typed alias for `Job` scheduler.

### Changed ⚠️ Breaking

- **All service-client CRUD interfaces** (`CloudServersClient`,
  `VolumesClient`, …) replaced with wrapper-based signatures. Direct
  use of `internal/clients/` is no longer the primary public API.
- `WithName(string)` renamed to `Named(string)` across the wrapper
  layer for naming consistency.
- Storage/size setters renamed with a `GB` suffix to match units:
  `WithSize` → `WithSizeGB`, `Size()` → `SizeGB()` on `BlockStorage`,
  `Snapshot`, `DBaaSBackup`, `KaaS`, and `DBaaS`.
- VPN lifetime and DPD setters suffixed with `Seconds` and
  standardized to `int`.
- Source/destination/classifier setters renamed to preposition form:
  `From*` / `To*` / `Of*`; classifiers become `OfType`, `OfFlavor`,
  `OfEngine`.
- `WithDefault(bool)` → `AsDefault()` / `NotDefault()`;
  `WithBootable(bool)` → `SetBootable()` / `UnsetBootable()`.
- `BillingPeriod` values translated to lowercase wire form on send.
- DBaaS autoscaling split into three fields:
  `Enabled`, `AvailableSpace`, `StepSize`.
- Storage terminal state `Available` renamed to `Active` to match the
  wire representation.
- Examples entry point moved from `cmd/example/` to
  `examples/all-resources/` and fully rewritten against the wrapper API.
- Interface guards moved from package-level vars to test functions.
- `WaitUntilState` renamed to `WaitUntilStates` (plural) to signal it
  accepts a variadic state list.
- `regionalMixin.location` renamed to `region`; consolidated to
  `inRegion`.
- `Region` struct renamed to `RegionInfo` in `pkg/types`.

### Removed ⚠️ Breaking

- `KubernetesVersion1313` — Kubernetes 1.31.3 is no longer in the live
  catalog. *(A deprecated alias pointing at `KubernetesVersion1323` was
  restored in v0.2.1 and will be removed in v0.3.0.)*
- `restclient.WaitForResourceState` polling primitive; polling now lives
  in wrapper `WaitUntil*` methods.
- Client-side polling in `BlockStorage`, `VPC`, and `SecurityGroup`
  `Create` paths.

### Fixed

- `Failed` state treated as a terminal error across all `WaitUntil*`
  paths.
- `vpnTunnelScopedMixin.intoVPNTunnel` now accepts the production
  camelCase `vpnTunnels` URI segment (#239).
- DBaaS billing plan, flavor, and node-pool name wire encoding; KMIP
  `Active` accepted as `Ready` state; longer wait budget with
  state-aware error propagation.
- VPN tunnel / VPN route / VPC peering route URL paths (#235, #236).
- `VPCPeeringRouteRequest` uses `RegionalResourceMetadataRequest`.
- `CloudServerRequest` omits `ElasticIP` and `KeyPair` when unset.

### Docs

- Documentation site restructured into a guided learning path:
  walkthrough guide + exhaustive resource reference (#225).
- `pkg/aruba` wrapper-layer architecture and naming conventions
  documented in `ai/ARCHITECTURE.md` and `ai/CONVENTIONS.md`.
- Resource reference: tag format rule documented; links to runnable
  examples added.
- `ai/` guidance files (`DEVEX.md`, `REPO.md`, `TECH_DEBT.md`) and
  `CLAUDE.md` added for contributor tooling.

### Tests

- Builder and group-accessor unit tests for `pkg/aruba` (#226).
- Examples coverage excluded from Codecov metrics.

---

## [0.1.28] — 2026-04-29

### Fixed

- `VPCPeeringRouteRequest` now uses `RegionalResourceMetadataRequest`
  instead of the base type; VPN enum constants added to `pkg/types`.
- `CloudServerRequest` omits `ElasticIP` and `KeyPair` fields when not
  set (omitempty alignment).

### CI

- Codecov configured to ignore `cmd/` directory from coverage reports.

### Internal

- gofmt alignment on VPN constants.

---

## [0.1.27] — 2026-04-24

### Fixed

- VPN tunnel, VPN route, and VPC peering route URL paths corrected.
- `VPCPeeringRouteResponse` uses `ResourceMetadataResponse`.

### Docs

- README badges added (Build, Go Version, Release, Codecov, License)
  and structure streamlined (#232).
- README inaccuracies fixed: removed reference to non-existent
  `NewFilterBuilder`; removed unused result variable in examples.

### CI

- `CODECOV_TOKEN` wired up in the CI workflow; coverage upload
  parameters fixed.

### Tests

- Collection-path and metadata-validation tests for VPN tunnel, VPN
  route, and VPC peering route resources.

---

## [0.1.26] — 2026-04-21

### Fixed

- Relaxed Create-response metadata validation: URI field is no longer
  required to be non-empty (matches actual server behavior).

---

## [0.1.25] — 2026-04-20

### Added

- Validate `ID`, `URI`, and `Name` on Create responses.
- Shared `internal/testutil` HTTP test helpers for resource client
  error-path coverage (TD-020).

### Fixed

- Panic on `nil` injected dependencies in all multi-argument client
  constructors.
- `auth.saveTicket`: increment counter only after a successful
  persistent write.
- `restclient`: remove initial poll sleep; preserve final state on
  timeout.
- `logger`: WARN-level messages written to `os.Stderr` instead of
  `os.Stdout`.
- `async.DefaultWaitFor` timeout raised from 60 s → 600 s;
  `async.DefaultRetries` raised to 60.
- Nil guard added to `ParseResponseBody`; non-JSON error bodies logged
  at DEBUG level instead of silently dropped.
- `buildDetebaseClient` typo fixed to `buildDatabaseClient`.
- `buildFileTokenRepository` argument order corrected.
- Preloaded access token actually stored in
  `NewTokenRepositoryWithAccessToken`.
- Per-tenant mutex acquired before setting the last-used timestamp.

### Docs

- Added `CLAUDE.md`, `ai/` guidance files, and a tech-debt backlog
  with effort/impact prioritization matrix.
- `internal/impl/interceptor/standard`: `Bind` documented as
  construction-only; closes TD-004.

### Tests

- Full error-path coverage matrix for all resource clients
  (audit, compute, container, database, metric, network, project,
  schedule, security, storage) migrated to `testutil` (TD-020,
  #149, #150).

### Internal

- Compile-time interface assertions for all resource-level client
  implementations and the logger.

---

## [0.1.24] — 2026-03-26

### Fixed

- `errorResponse` parsing correctly handles `validationErrors` and
  `BadRequest` response shapes.

---

## [0.1.23] — 2026-03-25

### Added

- Multi-tenant client layer (`pkg/multitenant`) with a graceful
  cleanup routine.

### Fixed

- Config struct is deep-copied before constructing the client,
  preventing shared-state mutation.

### Docs

- Multi-tenancy guide and usage examples.

### CI

- golangci-lint upgraded to v2 via action v7.
- gocyclo linter enabled; example packages excluded from lint.

---

## [0.1.22-alpha4] — 2026-03-23 *(pre-release)*

### Internal

- Fix `get` func in multi-tenant layer.

---

## [0.1.22-alpha2] — 2026-03-23 *(pre-release)*

### Fixed

- Deep-copy config before creating the client (regression from
  alpha1).

---

## [0.1.22-alpha1] — 2026-03-12 *(pre-release)*

- Pre-release iteration closing issues #94, #95, #96, #98.

> Note: v0.1.22-alpha3 was skipped during release-CI iteration;
> v0.1.22 stable was superseded by v0.1.23.

---

## [0.1.21] — 2026-02-08

### Added

- Token-issuer URL configurable at runtime (`WithTokenIssuerURL`
  option update).

### Docs

- Italian translation fixes.

---

## [0.1.20] — 2026-01-22

### Fixed

- Typo in `kmsPropertiesResponse` type name.

---

## [0.1.19] — 2026-01-22

### Added

- KMIP certificate download endpoint.

### Fixed

- KMS billing-period type mismatch.

### Docs

- Docusaurus version list sorted in descending order.

---

## [0.1.18] — 2026-01-21

### Docs

- Removed hardcoded version string from Docusaurus config; latest
  version is now selected automatically.

---

## [0.1.17] — 2026-01-21

### Added

- Nested KMS resource group: `Key` and `KMIP` sub-resources under KMS.

---

## [0.1.16] — 2026-01-15

### Fixed

- Missing `Zone` field on `DBaaSPropertiesRequest`.

---

## [0.1.15] — 2026-01-13

### Docs

- Doc-search rendering bug fixed.

---

## [0.1.14] — 2026-01-13

### Added

- Documentation site search support.
- Italian language (`it`) added to the documentation site.
- `CloudServer` `userData` properties.

---

## [0.1.13] — 2026-01-08

### Fixed

- CloudServer response serialization error.
- Latest-version resolution bug in the compute client.

---

## [0.1.12] — 2026-01-07

### Added

- Advanced subnet configuration support.

---

## [0.1.11] — 2025-12-31

### Fixed

- Container registry response serialization error.

---

## [0.1.10] — 2025-12-31

### CI

- Docs-version pipeline fixed.

---

## [0.1.9] — 2025-12-31

### Added

- Container registry support in the client builder.

### Docs

- Multi-release documentation support.

---

## [0.1.8] — 2025-12-30

### Docs

- First iteration of the docs-version release pipeline.

---

## [0.1.7] — 2025-12-30

### Added

- KaaS: `kubeconfig` download and update-request fixes.
- CloudServer: `power-on`, `power-off`, and `set-password` operations.

> Note: v0.1.5 and v0.1.6 were aborted release-CI attempts and were
> never published.

---

## [0.1.4] — 2025-12-29

### Added

- Full KaaS (Kubernetes-as-a-Service) feature set.

---

## [0.1.3] — 2025-12-22

### Fixed

- VPN tunnel type mismatch typo.

---

## [0.1.2] — 2025-12-19

### Fixed

- Snapshot success-state name typo.

---

## [0.1.1] — 2025-12-19

### Fixed

- BlockStorage type alignment.
- Project response type.

---

## [0.1.0] — 2025-12-09

Initial Alpha release of the Aruba Cloud SDK for Go.

---

<!-- compare links -->
[Unreleased]: https://github.com/Arubacloud/sdk-go/compare/v0.2.3...HEAD
[0.2.3]: https://github.com/Arubacloud/sdk-go/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/Arubacloud/sdk-go/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/Arubacloud/sdk-go/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/Arubacloud/sdk-go/compare/v0.1.28...v0.2.0
[0.1.28]: https://github.com/Arubacloud/sdk-go/compare/v0.1.27...v0.1.28
[0.1.27]: https://github.com/Arubacloud/sdk-go/compare/v0.1.26...v0.1.27
[0.1.26]: https://github.com/Arubacloud/sdk-go/compare/v0.1.25...v0.1.26
[0.1.25]: https://github.com/Arubacloud/sdk-go/compare/v0.1.24...v0.1.25
[0.1.24]: https://github.com/Arubacloud/sdk-go/compare/v0.1.23...v0.1.24
[0.1.23]: https://github.com/Arubacloud/sdk-go/compare/v0.1.21...v0.1.23
[0.1.22-alpha4]: https://github.com/Arubacloud/sdk-go/compare/v0.1.22-alpha2...v0.1.22-alpha4
[0.1.22-alpha2]: https://github.com/Arubacloud/sdk-go/compare/v0.1.22-alpha1...v0.1.22-alpha2
[0.1.22-alpha1]: https://github.com/Arubacloud/sdk-go/compare/v0.1.21...v0.1.22-alpha1
[0.1.21]: https://github.com/Arubacloud/sdk-go/compare/v0.1.20...v0.1.21
[0.1.20]: https://github.com/Arubacloud/sdk-go/compare/v0.1.19...v0.1.20
[0.1.19]: https://github.com/Arubacloud/sdk-go/compare/v0.1.18...v0.1.19
[0.1.18]: https://github.com/Arubacloud/sdk-go/compare/v0.1.17...v0.1.18
[0.1.17]: https://github.com/Arubacloud/sdk-go/compare/v0.1.16...v0.1.17
[0.1.16]: https://github.com/Arubacloud/sdk-go/compare/v0.1.15...v0.1.16
[0.1.15]: https://github.com/Arubacloud/sdk-go/compare/v0.1.14...v0.1.15
[0.1.14]: https://github.com/Arubacloud/sdk-go/compare/v0.1.13...v0.1.14
[0.1.13]: https://github.com/Arubacloud/sdk-go/compare/v0.1.12...v0.1.13
[0.1.12]: https://github.com/Arubacloud/sdk-go/compare/v0.1.11...v0.1.12
[0.1.11]: https://github.com/Arubacloud/sdk-go/compare/v0.1.10...v0.1.11
[0.1.10]: https://github.com/Arubacloud/sdk-go/compare/v0.1.9...v0.1.10
[0.1.9]: https://github.com/Arubacloud/sdk-go/compare/v0.1.8...v0.1.9
[0.1.8]: https://github.com/Arubacloud/sdk-go/compare/v0.1.7...v0.1.8
[0.1.7]: https://github.com/Arubacloud/sdk-go/compare/v0.1.4...v0.1.7
[0.1.4]: https://github.com/Arubacloud/sdk-go/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/Arubacloud/sdk-go/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/Arubacloud/sdk-go/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/Arubacloud/sdk-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/Arubacloud/sdk-go/releases/tag/v0.1.0
