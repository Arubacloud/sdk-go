# Repository Organization

## Module

`github.com/Arubacloud/sdk-go` — Official Go SDK for the Aruba Cloud API (Go 1.24+, currently Alpha / unstable API).

## Package layout

```
pkg/aruba/           — Public client entry point and Options builder
pkg/types/           — All request/response data models
pkg/async/           — Polling utilities for long-running operations
pkg/multitenant/     — Multi-tenant client management
internal/clients/    — Service-specific HTTP client implementations (one dir per service)
internal/impl/       — Pluggable subsystems: auth, interceptor, logger
internal/restclient/ — Low-level HTTP execution layer
cmd/example/         — Usage examples (excluded from linting)
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
