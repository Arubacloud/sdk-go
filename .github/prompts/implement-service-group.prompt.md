---
mode: agent
description: Scaffold a brand-new service group (domain) in the SDK end-to-end.
tools: [codebase, editFiles, runCommands]
---

You are adding a completely new service group to the Aruba Cloud Go SDK. Follow every step below in order.

## Inputs required
Ask the user for:
- **Domain name** — singular, lowercase (e.g. `billing`, `dns`)
- **Service group display name** — PascalCase (e.g. `Billing`, `DNS`)
- **First resource(s)** to implement inside this group

---

## Step 1 — Create `internal/clients/<domain>/`

Create these files:

### `path.go`
```go
package <domain>

const (
    // Example — add real paths per resource
    // <PascalResource>sPath = "/projects/{projectID}/<resources>"
)
```

### `version.go`
```go
package <domain>

const (
    // Example — add per-operation API version constants per resource
    // <Domain><PascalResource>List = "YYYY-MM-DD"
)
```

### `<domain>.go` — the service group aggregator
```go
package <domain>

import "github.com/Arubacloud/sdk-go/internal/restclient"

// <Domain>Client exposes all <domain> resource clients.
type <Domain>Client struct {
    client restclient.Client
    // Add one field per resource client impl
}

// New<Domain>Client constructs the service group.
func New<Domain>Client(client restclient.Client) *<Domain>Client {
    return &<Domain>Client{
        client: client,
        // Initialise each resource impl here
    }
}
```

---

## Step 2 — Add the public interface to `pkg/aruba/`

Create `pkg/aruba/<domain>.go`:
```go
package aruba

import "<module>/internal/clients/<domain>"

// <Domain>Client is the public interface for <domain> resources.
type <Domain>Client interface {
    // Add one accessor per resource, e.g.:
    // <PascalResource>() <domain>.<PascalResource>Client
}
```

---

## Step 3 — Expose from the top-level `Client`

In `pkg/aruba/client.go`, add the accessor to the `Client` interface:
```go
From<Domain>() <Domain>Client
```

In `pkg/aruba/aruba.go` (the `clientImpl` struct), add a field and implement the accessor:
```go
<camelDomain> <Domain>Client

func (c *clientImpl) From<Domain>() <Domain>Client {
    return c.<camelDomain>
}
```

---

## Step 4 — Wire in `builder.go`

In `buildClient()`, after the last existing service group construction, add:
```go
<camelDomain>Client := <domain>.New<Domain>Client(restClient)
```

And include it in the `clientImpl` construction.

---

## Step 5 — Implement the first resource(s)

For each resource, run the `implement-resource` agent with the new domain name.

---

## Step 6 — Verify

```bash
make verify
```

Fix all errors before finishing.
