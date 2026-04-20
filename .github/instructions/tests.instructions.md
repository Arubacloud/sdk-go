---
applyTo: "**/*_test.go"
---

# Test conventions

You are writing or editing a test file in this SDK.

## Setup pattern
- Always use `httptest.NewServer` — no real network calls in tests
- The test server must handle both `/token` (returns a dummy JWT) and the resource endpoint
- Use `http.NewServeMux()` with explicit `HandleFunc` registrations, not a catch-all handler
- Call `t.Cleanup(srv.Close)` immediately after creating the server

## Logger
- Always use `noop.NoOpLogger{}` from `internal/impl/logger/noop` — never a real logger
- Pass it as `logger` to `restclient.NewClient`

## Test structure
- Use subtests for every scenario: `t.Run("descriptive scenario name", func(t *testing.T) { ... })`
- Minimum scenarios per method: **happy path**, **validation error (empty required ID)**, **API 4xx error**
- Group subtests under a `Test<PascalResource><Method>` top-level function

## Assertions
- Check `resp.IsSuccess()` for 2xx — never `resp.StatusCode == 200`
- Check `resp.IsError()` for 4xx/5xx
- Use `t.Fatalf` for unexpected errors to stop the subtest immediately
- Do not use third-party assertion libraries — stdlib `testing` only

## Constructing the impl under test
```go
logger := &noop.NoOpLogger{}
interceptor := standard.NewInterceptor()
rc := restclient.NewClient(srv.URL, http.DefaultClient, interceptor, logger)
client := New<PascalResource>ClientImpl(rc)
```

## Mock token endpoint
```go
mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"access_token":"test","token_type":"Bearer","expires_in":3600}`))
})
```

## Generated mocks
- Mocks are generated with MockGen (`go.uber.org/mock/gomock`)
- Generated files are named `zz_mock_<Type>_test.go` in the same package
- Do not hand-write mock structs — regenerate with `go generate`
