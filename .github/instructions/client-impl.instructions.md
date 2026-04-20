---
applyTo: "internal/clients/**"
---

# Service client implementation conventions

You are editing a service client inside `internal/clients/<domain>/`.

## Struct & constructor
- Impl struct is unexported: `type <camelResource>ClientImpl struct { client restclient.Client }`
- Constructor returns a concrete pointer: `func New<PascalResource>ClientImpl(client restclient.Client) *<camelResource>ClientImpl`
- All methods use pointer receivers: `func (c *<camelResource>ClientImpl) List(...)`

## Mandatory method sequence — never deviate
1. `c.client.Logger().Debugf("...")` — log the operation and key IDs first
2. `types.Validate*(...)` — fail fast before any HTTP call
3. If `params == nil` → `params = &types.RequestParameters{}`
4. If `params.APIVersion == nil` → set the domain-specific version constant from `version.go`
5. `queryParams := params.ToQueryParams()` / `headers := params.ToHeaders()`
6. Marshal body: `bodyBytes, err := json.Marshal(body)` (Create/Update only; handle error immediately)
7. `httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, bodyBytes, queryParams, headers)`
8. `if err != nil { return nil, fmt.Errorf("...resource...: %w", err) }`
9. `defer httpResp.Body.Close()`
10. `return types.ParseResponseBody[types.<PascalResource>Response](httpResp)`

## Path construction
- Use `strings.NewReplacer("{projectID}", projectID, "{resourceID}", resourceID).Replace(<Path>)` or `fmt.Sprintf`
- Never hardcode path strings inside method bodies — always reference the constant from `path.go`

## Cross-resource dependencies
- If this resource requires another resource to be in an active state first, accept the dependency as a constructor parameter (concrete `*<otherResource>ClientImpl`, not the interface)
- Call the internal helper before the HTTP call, e.g. `c.securityGroups.waitForActive(ctx, projectID, sgID)`

## Error handling
- Validation errors: `return nil, fmt.Errorf("projectID cannot be empty")`
- Wrapped transport/marshal errors: `return nil, fmt.Errorf("<resource> <operation>: %w", err)`
- Never inspect raw HTTP status codes; use `resp.IsSuccess()` / `resp.IsError()`

## Path and version constants
- Path constants live in `path.go` — PascalCase, e.g. `CloudServersPath`
- Version constants live in `version.go` — `<Domain><PascalResource><Action>`, e.g. `ComputeCloudServerList`
