---
id: changelog
title: Changelog
sidebar_label: Changelog
---

All notable changes to this project are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## Versioning & Branch Policy

| Branch | Version series | Status |
|--------|----------------|--------|
| `main` | **v1.x** | Active development — GA release line |
| `legacy` | **v0.1.x** | Maintenance — bug fixes and security patches only |

- **`main` — v1.x (current)**: GA release line. v1.0.0 stabilises the
  `pkg/aruba/` wrapper layer introduced in v0.2.0: flattened getter taxonomy,
  async polling helpers, automatic User-Agent injection, and the final
  `*Request` / `*Response` / `*Common` naming convention across `pkg/types`.
  See [v1.0.0](#100--2026-05-29) below for migration details.
- **`legacy` — v0.1.x (maintenance)**: Supported for **6 months** after
  the v0.2.0 release date with bug-fix and security-patch releases only.
  No new features will be backported. Once the support window closes the
  `legacy` branch will be archived.

---

## [Unreleased]

---

## [1.0.8] — 2026-07-27

### Added

- **`kvPrefix` support in Vault credentials repository** (`internal/impl/auth/credentialsrepository/vault`, `pkg/aruba`) —
  `WithVaultCredentialsRepository` now accepts a `kvPrefix string` parameter (inserted between
  `kvMount` and `kvPath`). The prefix is prepended to the secret path when calling the Vault KV v2
  `Get` API (`<kvMount>/<kvPrefix>/<kvPath>`), enabling secrets stored under a shared path prefix
  within a single KV mount. Passing an empty string preserves the previous behaviour — the path
  is constructed as `<kvMount>/<kvPath>` with no prefix segment.

---

## [1.0.7] — 2026-07-20

### Fixed

- **`RegionalResourceMetadataRequest.Location` omitted from Update requests when region is unset** (`pkg/types`, `pkg/aruba`) —
  `Location` was typed as `LocationRequest` (value, no `omitempty`), so the SDK always serialised
  `"location":{"value":"..."}` in every Create and Update request body. The SecurityRule Update
  endpoint (and potentially others) reject this field with `400 Validation: Location Invalid` when
  a location is included in an Update. Sending `{"value":""}` also fails — the field must be
  absent entirely.
  Fix: `Location` is now `*LocationRequest` with `json:"location,omitempty"`, and `toLocation()`
  returns `nil` when the region is empty (e.g. after `InRegion("")`) so the field is omitted from
  the serialised JSON. Create behaviour is unchanged — a region is always set before Create so
  `toLocation()` still returns a non-nil pointer and the location is included as before.

- **`SecurityRulePropertiesResponse.Port` deserialisation from object form** (`pkg/types`, `pkg/aruba`) —
  The SecurityRule Get endpoint returns `port` as a plain JSON string (`"80"`), but the Update
  endpoint returns it as a JSON object (`{"value":"80"}`). The field was typed `string`, causing
  unmarshalling the Update response to fail with:
  `json: cannot unmarshal object into Go struct field SecurityRulePropertiesResponse.properties.port of type string`.
  Fix: `Port` in `SecurityRulePropertiesResponse` is now `FlexPort`, a custom type whose
  `UnmarshalJSON` accepts both the string and object forms and normalises to a plain string.
  `MarshalJSON` still emits a plain string so request bodies are unaffected.

---

## [1.0.6] — 2026-07-14

### Added

- **Cloud Server update delta setters** (`pkg/aruba`) — eight new fluent setters record
  association/attachment intent on a `*CloudServer` wrapper for use in `Update` calls:
  `AssociateSubnets`, `DisassociateSubnets`, `AssociateSecurityGroups`,
  `DisassociateSecurityGroups`, `AssociateElasticIPs`, `DisassociateElasticIPs`,
  `AttachDataVolumes`, `DetachDataVolumes`. Repeated calls on the same setter append;
  queues are cleared after a successful `Update`.

- **Four new low-level API endpoints** (`internal/clients/compute`) — `AssociateSubnets`,
  `AssociateSecurityGroups`, `AssociateElasticIPs`, `AttachDetachDataVolumes` map to the
  dedicated Aruba Cloud REST actions (`associateDisassociateSubnets`,
  `associateDisassociateSecurityGroups`, `associateDisassociateElasticIPs`,
  `attachDetachDataVolumes`), replacing the previous practice of bundling all changes into a
  single PUT that only honoured name/tag fields.

- **Four new `pkg/types` request structs** — `CloudServerAssociateSubnetsRequest`,
  `CloudServerAssociateSecurityGroupsRequest`, `CloudServerAssociateElasticIPsRequest`,
  `CloudServerAttachDetachDataVolumesRequest`.

### Changed

- **`CloudServersClient.Update` is now a smart dispatcher** (`pkg/aruba`) — a single
  `Update(ctx, server)` call inspects the wrapper's pending state and fans out to the
  correct API endpoint(s) in sequence:
  - name/tags differ from the last hydrated response → `PUT /cloudServers/:id`
  - subnet delta queues non-empty → `POST …/associateDisassociateSubnets`
  - security-group delta queues non-empty → `POST …/associateDisassociateSecurityGroups`
  - elastic-IP delta queues non-empty → `POST …/associateDisassociateElasticIPs`
  - data-volume delta queues non-empty → `POST …/attachDetachDataVolumes`

  If nothing changed, `Update` is a no-op (zero API calls). Sub-calls are sequential; if
  one fails, the wrapper reflects the state after the last successful call.

  The `CloudServersClient.Update` **signature is unchanged** — this is a behaviour
  improvement only. The previous `Update` silently sent all fields via PUT, which the
  platform ignored for everything except name and tags.

### Removed

- **`ManageSubnets`, `ManageSecurityGroups`, `ManageElasticIPs`, `ManageDataVolumes`** — the
  four action methods added as part of the in-progress v1.0.6 work are not included in the
  final design. The smart `Update` dispatcher supersedes them; client code (Terraform,
  acloud-cli) should use the delta setters + `Update` instead.

---

## [1.0.5] — 2026-07-13

### Fixed

- **DBaaS Backup wire format** (`pkg/types`, `pkg/aruba`) — `BackupPropertiesRequest.Database` and
  `BackupPropertiesResponse.Database` were typed as `ReferenceResourceCommon{URI string}`, causing
  the SDK to send `"database":{"uri":"..."}` on the wire. The backup API requires
  `"database":{"name":"<db-name>"}`.
  Sending the wrong format produced a `400 semantic` error: *Specified Database name is not found.*
  The new `DatabaseNameRef{Name string}` type is now used for both request and response. The
  `toRequest()` helper extracts the plain name from a URI path (last segment after `/databases/`)
  so existing callers using `FromDatabase(aruba.URI("…"))` continue to work without change.

### Added

- **`DatabaseName()` accessor on `DBaaSBackup`** (`pkg/aruba`) — dedicated getter that always
  returns the bare database name. `DatabaseURI()` is preserved for backward compatibility but its
  return value is now the raw stored reference (a full URI when set via `FromDatabase`, or the bare
  name after response hydration). Prefer `DatabaseName()` when you only need the name.

---

## [1.0.4] — 2026-06-08

### Added

- **`EnablePrivateCluster()` on KaaS** (`pkg/aruba`) — preferred boolean-flag setter for private
  cluster mode, following the `HighlyAvailable()` naming convention. Thin alias for the existing
  `WithPrivateCluster()`. `WithAPIServerAccessProfile(*types.KaaSAPIServerAccessProfilePropertiesRequest)`
  is now marked `Deprecated`.

### Fixed

- **`APIServerAuthorizedIPRanges()` returns a defensive copy** — the getter previously returned
  the underlying slice directly, allowing callers to mutate internal wrapper state.

---

## [1.0.3] — 2026-06-08

### Added

- **KaaS fluent API server access profile setters** (`pkg/aruba`) — `WithPrivateCluster()` and
  `WithAuthorizedIPRanges(ranges ...string)` let callers configure the API server access profile
  without importing `pkg/types`. (#323)
- **KaaS API server access profile getters** — `APIServerPrivateCluster() bool` and
  `APIServerAuthorizedIPRanges() []string`.

---

## [1.0.2] — 2026-06-06

### Added

- **Audit/lineage getters** (`pkg/aruba`) — every Family-A wrapper now exposes `CreatedBy()`,
  `UpdatedBy()`, `CreatedUser()`, `UpdatedUser()` via `responseMetadataMixin`. (#316)

### Changed

- Examples (`examples/all-resources`) now print `CreatedBy` / `CreatedAt` in summaries.

### Documentation

- Documented the new lineage getters in getter-taxonomy tables (EN + IT).

---

## [1.0.0] — 2026-05-29

### Added

- **Flattened getter taxonomy** (`pkg/aruba`) — consistent response-side scalar getters across
  all Family-A wrappers.
- **`pkg/aruba.Version` constant** — single source of truth for the SDK version string.
- **Automatic User-Agent injection** (`Options.WithUserAgent`) — every outbound request carries
  `User-Agent: sdk-go@<version>` by default. (#310)

### Changed (Breaking)

All exported structs in `pkg/types/` now carry a strict three-way suffix: `*Request`, `*Response`,
or `*Common`. This is a compile-time breaking change for code referencing `pkg/types` directly;
code using only `pkg/aruba` wrappers is unaffected.

See the root [CHANGELOG.md](https://github.com/Arubacloud/sdk-go/blob/main/CHANGELOG.md#100--2026-05-29)
for the full list of renamed types.

---

## [0.3.0] — 2026-05-26

### Changed (Breaking)

- **Fluent setter vocabulary redesigned** (`pkg/aruba`) — all wrapper setters now follow a
  natural-language taxonomy. See the root [CHANGELOG.md](https://github.com/Arubacloud/sdk-go/blob/main/CHANGELOG.md#030--2026-05-26)
  for the full rename table.

### Added

- `List[T]` embeds `httpEnvelopeMixin` for HTTP-envelope accessors. (#298)
- `RawJSON() []byte` / `RawYAML() []byte` on every resource wrapper and `List[T]`.
- New doc page **Working at Low Level** (EN + IT).

---

## [0.2.3] — 2026-05-22

### Added

- `WaitUntilGone(ctx, opts...)` on every deletable, pollable resource wrapper.
- Typed `types.State` with exported `State*` constants.

### Changed (Breaking)

- `WaitUntilStates` signature changed from `[]string` to `[]types.State`.
- `pkg/types` convenience pointer helpers removed; use `ptr.To(...)` from `k8s.io/utils/ptr`.

---

## [0.2.2] — 2026-05-22

### Changed

- `billingPeriod` wire field replaced by a structured `billingPlan` object across billed resources.

### Fixed

- DBaaS `User` password is now base64-encoded at the wire boundary.

---

## [0.2.1] — 2026-05-16

### Added

- Typed network `Ref` constructors for all 10 network resources. (#268)
- `*CloudServer.WithBillingPeriod(period)` fluent setter. (#267)
- `*KaaS` node-pool control setters (`ClearNodePools`, `ReplaceNodePools`, `SetNodePools`). (#279)

### Fixed

- `List[T].Next/Prev/First/Last` pagination links wired for all 31 adapters. (#269 et al.)
- `Job.WithEnabled(false)` now correctly disables a job via Update. (#282)

---

## [0.2.0] — 2026-05-13

### Added

- `pkg/aruba/` wrapper layer with **31 fluent resource types**.
- `List[T]` generic paginated-collection type.
- Async polling: `WaitUntilActive`, `WaitUntilStates`, `WaitUntilReady`, `WaitUntilNotUsed`, `WaitUntilUsed`.

### Changed (Breaking)

- All service-client CRUD interfaces replaced with wrapper-based signatures.

---

## [0.1.28] — 2026-04-29

### Fixed

- `VPCPeeringRouteRequest` now uses `RegionalResourceMetadataRequest`.
- `CloudServerRequest` omits `ElasticIP` and `KeyPair` fields when not set.

---

## [0.1.27] — 2026-04-24

### Fixed

- VPN tunnel, VPN route, and VPC peering route URL paths corrected.

---

## Older releases

For releases prior to v0.1.27, see the full
[CHANGELOG.md](https://github.com/Arubacloud/sdk-go/blob/main/CHANGELOG.md) on GitHub.
