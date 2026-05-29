# Repository Organization

## Module

`github.com/Arubacloud/sdk-go` — Official Go SDK for the Aruba Cloud API (Go 1.24+, currently Alpha / unstable API).

## Package layout

```
pkg/aruba/               — Public client entry point, Options builder, and wrapper layer (pkg/aruba/doc.go)
pkg/aruba/resource_*.go  — Fluent builder wrappers + adapters; New* factory constructors are co-located in each file
pkg/aruba/mixin_common.go   — Cross-cutting mixins (errMixin, metadataMixin, regionalMixin, zonalMixin, responseMetadataMixin, linkedMixin, httpEnvelopeMixin)
pkg/aruba/mixin_scoped.go   — Parent-scope mixins (projectScopedMixin, vpcScopedMixin, securityGroupScopedMixin, dbaasScopedMixin, databaseScopedMixin, backupScopedMixin, kmsScopedMixin, vpnTunnelScopedMixin, vpcPeeringScopedMixin)
pkg/aruba/mixin_status.go   — Lifecycle-state mixin (statusMixin, WaitUntilActive, WaitUntilReady, WaitUntilStates)
pkg/aruba/ref.go         — Ref interface + extractID + parseURIIDs
pkg/aruba/list.go        — Generic List[T Wrapper] paginated container
pkg/aruba/list_pagination.go — Pagination helpers
pkg/aruba/aliases.go     — Typed enum constants re-exported from pkg/types
pkg/aruba/billing_period.go — BillingPeriod enum and helpers
pkg/aruba/call_options.go   — CallOption and per-call option types
pkg/aruba/errors.go      — *HTTPError wrapper
pkg/aruba/client_root.go    — Root Client interface (exposes FromCompute(), FromNetwork(), etc.)
pkg/aruba/client_<domain>.go — Per-domain service-group interfaces (client_compute.go, client_network.go, …)
pkg/types/           — All request/response/common data models (Request/Response/Common naming rule); includes typed State with predicate methods. Package convention codified in pkg/types/doc.go.
pkg/async/           — Polling utilities for long-running operations (pkg/async/doc.go)
pkg/multitenant/     — Multi-tenant client management (pkg/multitenant/doc.go)
pkg/util/middleware/ — HTTP interceptor helpers (WithCustomHeaders, WithUserAgent, etc.); package name is restclient
internal/clients/    — Service-specific HTTP client implementations (one dir per service)
internal/impl/       — Pluggable subsystems: auth, interceptor, logger
internal/restclient/ — Low-level HTTP execution layer
examples/all-resources/ — Usage examples (excluded from linting)
docs/                — Testing docs, versioning scripts, and Docusaurus site
docs/website/        — Docusaurus documentation site (multi-locale, includes Italian)
```

## Service groups

The top-level `Client` interface exposes ten service groups:

| Accessor | Service |
|---|---|
| `FromCompute()` | Cloud servers, key pairs |
| `FromNetwork()` | Networking resources |
| `FromStorage()` | Storage resources |
| `FromProject()` | Project management |
| `FromDatabase()` | Managed databases |
| `FromContainer()` | Container services |
| `FromAudit()` | Audit logs |
| `FromMetric()` | Metrics |
| `FromSchedule()` | Scheduled jobs |
| `FromSecurity()` | Security resources |

Each group is backed by an implementation under `internal/clients/<service>/`.

## Wrapper layer (`pkg/aruba/`)

Each resource type has a self-contained triplet inside a single `resource_<name>.go` file: **Wrapper** (chainable builder), **low-level client interface** (adapter contract), **Adapter** (bridges wrapper ↔ `internal/clients/<x>`). Pure sub-builders with no CRUD of their own — `JobStep`, `NodePool`, `VPNIKE`, `VPNESP`, `VPNPSK`, `VPNIPConfig`, `SubnetDHCP` — follow only the Wrapper section.

Resources split into two wire-shape families:

- **Family A** — `Metadata{Properties{...}}` envelope, regional, `statusMixin`. The large majority of resources.
- **Family B** — flat request body, no metadata envelope, no tags/region/status. Set: Database, Key, Kmip, User, Grant.

See `ai/ARCHITECTURE.md` § "Wrapper layer" for the full pattern, mixin catalogue, and non-standard cases.
