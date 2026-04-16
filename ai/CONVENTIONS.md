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
9. `types.ParseResponseBody[T](httpResp)` for standard responses; manual unmarshal only for complex cases

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
