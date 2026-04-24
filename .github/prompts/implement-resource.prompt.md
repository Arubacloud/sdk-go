---
mode: agent
description: Implement a new resource end-to-end in the SDK following all project conventions.
tools: [codebase, editFiles, runCommands]
---

You are implementing a new resource in the Aruba Cloud Go SDK. Follow every step below exactly and in order. Do not skip steps.

## Inputs required
Ask the user for the following if not already provided:
- **Domain** (e.g. `compute`, `network`, `storage`) — which service group this belongs to
- **Resource name** (e.g. `snapshot`, `firewall-rule`) — singular, kebab-case
- **Operations** — which of: `List`, `Get`, `Create`, `Update`, `Delete` to implement

---

## Step 1 — Define types in `pkg/types/`

Create (or append to) `pkg/types/<domain>.<resource>.go`.

Rules:
- Request struct name: `<PascalResource>Request`
- Single-item response: `<PascalResource>Response`
- Collection response: `<PascalResource>List` (must embed `types.ListResponse`)
- All fields use pointer types where the value is optional
- No JSON tags unless you are certain of the wire format

Example skeleton:
```go
package types

// <PascalResource>Request is the input for creating or updating a <resource>.
type <PascalResource>Request struct {
    Name        *string `json:"name,omitempty"`
    Description *string `json:"description,omitempty"`
}

// <PascalResource>Response represents a single <resource> returned by the API.
type <PascalResource>Response struct {
    ID   *string `json:"id,omitempty"`
    Name *string `json:"name,omitempty"`
}

// <PascalResource>List is the paginated collection of <resource> items.
type <PascalResource>List struct {
    ListResponse
    Items []<PascalResource>Response `json:"items,omitempty"`
}
```

---

## Step 2 — Add path constants to `internal/clients/<domain>/path.go`

Append:
```go
const (
    <PascalResource>sPath = "/projects/{projectID}/<resources>"
    <PascalResource>Path  = "/projects/{projectID}/<resources>/{resourceID}"
)
```

---

## Step 3 — Add API version constants to `internal/clients/<domain>/version.go`

Add one constant per operation:
```go
const (
    <Domain><PascalResource>List   = "YYYY-MM-DD"
    <Domain><PascalResource>Get    = "YYYY-MM-DD"
    <Domain><PascalResource>Create = "YYYY-MM-DD"
    <Domain><PascalResource>Update = "YYYY-MM-DD"
    <Domain><PascalResource>Delete = "YYYY-MM-DD"
)
```
Use the same date as other constants in this file unless you know the real version.

---

## Step 4 — Create the resource client file `internal/clients/<domain>/<resource>.go`

Structure:
```go
package <domain>

import (...)

// <PascalResource>Client is the interface for managing <resource> resources.
type <PascalResource>Client interface {
    List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.<PascalResource>List], error)
    Get(ctx context.Context, projectID, resourceID string, params *types.RequestParameters) (*types.Response[types.<PascalResource>Response], error)
    Create(ctx context.Context, projectID string, body *types.<PascalResource>Request, params *types.RequestParameters) (*types.Response[types.<PascalResource>Response], error)
    Update(ctx context.Context, projectID, resourceID string, body *types.<PascalResource>Request, params *types.RequestParameters) (*types.Response[types.<PascalResource>Response], error)
    Delete(ctx context.Context, projectID, resourceID string, params *types.RequestParameters) (*types.Response[types.<PascalResource>Response], error)
}

type <camelResource>ClientImpl struct {
    client restclient.Client
}

func New<PascalResource>ClientImpl(client restclient.Client) *<camelResource>ClientImpl {
    return &<camelResource>ClientImpl{client: client}
}
```

Each method must follow this exact sequence:
1. `c.client.Logger().Debugf("...")` — log the operation and key IDs
2. `types.Validate*(projectID, ...)` — fail fast
3. If `params == nil`, set `params = &types.RequestParameters{}`
4. If `params.APIVersion == nil`, set `params.APIVersion = ptr(<Domain><PascalResource><Action>)`
5. `queryParams := params.ToQueryParams()`; `headers := params.ToHeaders()`
6. Marshal body with `json.Marshal` (Create/Update only)
7. `httpResp, err := c.client.DoRequest(ctx, http.Method*, path, body, queryParams, headers)`
8. `defer httpResp.Body.Close()`
9. `return types.ParseResponseBody[types.<PascalResource>Response](httpResp)`

---

## Step 5 — Expose from the service group

In `internal/clients/<domain>/<group>.go`, add a field and accessor:
```go
type <Domain>Client struct {
    ...
    <camelResource> *<camelResource>ClientImpl
}

func (c *<Domain>Client) <PascalResource>() <PascalResource>Client {
    return c.<camelResource>
}
```

Wire it in the constructor.

---

## Step 6 — Wire to `pkg/aruba/` if needed

If the accessor is already in `pkg/aruba/<domain>.go`, add the method delegation. Otherwise add it.

---

## Step 7 — Write tests

Create `internal/clients/<domain>/<resource>_test.go` with a subtest for every operation. See the `write-tests` agent for the test template.

---

## Step 8 — Verify

Run:
```bash
make verify
```
Fix any lint or test failures before finishing.
