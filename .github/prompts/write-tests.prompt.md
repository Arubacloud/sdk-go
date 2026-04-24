---
mode: agent
description: Write complete test coverage for a resource client in internal/clients/.
tools: [codebase, editFiles, runCommands]
---

You are writing tests for an existing resource client in `internal/clients/<domain>/<resource>.go`. Follow every rule below.

## Inputs required
Ask the user for:
- **Domain** and **resource** to test (e.g. `compute`, `cloudserver`)

---

## Rules

- Test file: `internal/clients/<domain>/<resource>_test.go`, same package as the implementation
- Use `httptest.NewServer` — never real HTTP calls
- Use `noop.NoOpLogger{}` from `internal/impl/logger/noop`
- Use `t.Run("scenario", func(t *testing.T) { ... })` for every case
- All method calls use pointer receivers (already enforced by the impl)
- Assert on `resp.IsSuccess()` / `resp.IsError()` — never on raw status codes
- Cover at minimum: **happy path**, **validation error (empty IDs)**, **API error (4xx)**

---

## Test file skeleton

```go
package <domain>

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
    "github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
    "github.com/Arubacloud/sdk-go/internal/restclient"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

func setup<PascalResource>Server(t *testing.T, statusCode int, body interface{}) (*httptest.Server, *<camelResource>ClientImpl) {
    t.Helper()
    mux := http.NewServeMux()
    mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
    })
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        if body != nil {
            _ = json.NewEncoder(w).Encode(body)
        }
    })
    srv := httptest.NewServer(mux)
    t.Cleanup(srv.Close)

    logger := &noop.NoOpLogger{}
    interceptor := standard.NewInterceptor()
    rc := restclient.NewClient(srv.URL, http.DefaultClient, interceptor, logger)
    client := New<PascalResource>ClientImpl(rc)
    return srv, client
}

func Test<PascalResource>List(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        _, client := setup<PascalResource>Server(t, http.StatusOK, types.<PascalResource>List{})
        resp, err := client.List(context.Background(), "proj-1", nil)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if !resp.IsSuccess() {
            t.Fatalf("expected success, got %d", resp.StatusCode)
        }
    })

    t.Run("empty project ID returns validation error", func(t *testing.T) {
        _, client := setup<PascalResource>Server(t, http.StatusOK, nil)
        _, err := client.List(context.Background(), "", nil)
        if err == nil {
            t.Fatal("expected validation error, got nil")
        }
    })

    t.Run("API error 404", func(t *testing.T) {
        _, client := setup<PascalResource>Server(t, http.StatusNotFound, types.ErrorResponse{Title: "not found"})
        resp, err := client.List(context.Background(), "proj-1", nil)
        if err != nil {
            t.Fatalf("unexpected transport error: %v", err)
        }
        if !resp.IsError() {
            t.Fatalf("expected error response, got %d", resp.StatusCode)
        }
    })
}

// Add equivalent Test<PascalResource>Get, Test<PascalResource>Create, etc. following the same pattern.
```

---

## After writing tests

Run:
```bash
make test
```

Fix any failures. Ensure race detector passes (`make test` uses `-race` by default).
