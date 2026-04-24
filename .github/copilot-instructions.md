# GitHub Copilot Instructions

This file provides guidance to GitHub Copilot when working with this repository.
Detailed guidance is split across files in the `ai/` folder:

| Task type | File |
|---|---|
| Building, running tests, linting, CI, reference docs | [`ai/DEVEX.md`](../ai/DEVEX.md) |
| Finding files, understanding package structure, service groups | [`ai/REPO.md`](../ai/REPO.md) |
| Implementing features, understanding design patterns, adding resources | [`ai/ARCHITECTURE.md`](../ai/ARCHITECTURE.md) |
| Naming, code style, type conventions, error handling rules | [`ai/CONVENTIONS.md`](../ai/CONVENTIONS.md) |
| Known bugs, risks, and refactoring backlog | [`ai/TECH_DEBT.md`](../ai/TECH_DEBT.md) |

---

## Quick reference

### Module & Go version
`github.com/Arubacloud/sdk-go` — Go 1.24+

### Package layout
- `pkg/aruba/` — public client entry point and Options builder
- `pkg/types/` — all request/response data models
- `pkg/async/` — polling utilities for long-running operations
- `pkg/multitenant/` — multi-tenant client management
- `internal/clients/` — service-specific HTTP client implementations (one dir per service)
- `internal/impl/` — pluggable subsystems: auth, interceptor, logger
- `internal/restclient/` — low-level HTTP execution layer

### Key naming conventions
| Kind | Suffix | Example |
|---|---|---|
| Input/request struct | `Request` | `CloudServerRequest` |
| Single output struct | `Response` | `CloudServerResponse` |
| Collection output struct | `List` (embeds `ListResponse`) | `CloudServerList` |
| Service client interface | none | `CloudServersClient` |
| Service client implementation | `Impl` (unexported) | `cloudServersClientImpl` |
| Constructor | `New<TypeName>` | `NewCloudServersClientImpl` |

### Service method structure
Every resource method must follow this sequence exactly:
1. `c.client.Logger().Debugf(...)` — log operation and key IDs
2. `types.Validate*(...)` — fail fast before any HTTP call
3. If `params == nil`, initialize `params = &types.RequestParameters{}`
4. If `params.APIVersion == nil`, set the domain-specific version constant
5. `params.ToQueryParams()` / `params.ToHeaders()`
6. Marshal body with `json.Marshal` if applicable
7. `c.client.DoRequest(ctx, method, path, body, queryParams, headers)`
8. `defer httpResp.Body.Close()`
9. `types.ParseResponseBody[T](httpResp)`

### Error handling
- Use `resp.IsSuccess()` / `resp.IsError()` — never inspect raw HTTP status codes in business logic
- Validation errors return as the `error` return value: `fmt.Errorf("project cannot be empty")`
- Wrap errors with `fmt.Errorf("context: %w", err)`

### Testing
- Each resource file has a paired `_test.go` in the same package
- Tests use `httptest.NewServer` to mock the `/token` endpoint and the resource endpoint
- Use `noop.NoOpLogger{}` in all tests
- Use subtests: `t.Run("scenario", func(t *testing.T) { ... })`
- All methods on impl types use pointer receivers

### Build commands
```bash
make build        # go build ./...
make test         # run all tests with race detection and coverage
make test-short   # quick tests without coverage
make lint         # go fmt + go vet
make verify       # lint + test (recommended before committing)
make all          # tidy + lint + build + test
```
