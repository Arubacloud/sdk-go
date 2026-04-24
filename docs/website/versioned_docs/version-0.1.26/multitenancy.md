---
id: multitenancy
title: Multitenancy
---

The `pkg/multitenant` package provides an in-memory tenant-to-client registry for the Aruba Cloud SDK. It is useful when your application serves multiple tenants and each tenant needs its own `aruba.Client`.

## Overview

The package exposes:

- A `Multitenant` interface to create, store, retrieve, and clean up tenant clients
- A default implementation backed by a map and mutex
- A cleanup routine helper to periodically remove stale tenants

Core files:

- `pkg/multitenant/multitenant.go`
- `pkg/multitenant/cleanup_routine.go`

## Main Interface

The `Multitenant` interface supports these operations:

- `New(tenant string) error`: create client using template options
- `NewFromOptions(tenant string, options *aruba.Options) error`: create client from explicit options
- `Add(tenant string, client aruba.Client)`: inject an existing client
- `Get(tenant string) (aruba.Client, bool)`: retrieve client with existence flag
- `MustGet(tenant string) aruba.Client`: retrieve or terminate process if missing
- `GetOrNil(tenant string) aruba.Client`: retrieve or return `nil`
- `CleanUp(from time.Duration)`: delete inactive tenants

## Creating a Manager

### Empty manager

Use this when you want to add tenant clients manually or via `NewFromOptions`:

```go
mt := multitenant.New()
```

### Manager with template

Use this when tenants share a common base configuration:

```go
opts := aruba.DefaultOptions(clientID, clientSecret)
mt := multitenant.NewWithTemplate(opts)

// Later:
if err := mt.New("tenant-a"); err != nil {
    // handle error
}
```

## Usage Patterns

### Add an existing client

```go
client, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
if err != nil {
    // handle error
}

mt.Add("tenant-a", client)
```

### Create from tenant-specific options

```go
tenantOpts := aruba.DefaultOptions(tenantClientID, tenantClientSecret)
if err := mt.NewFromOptions("tenant-a", tenantOpts); err != nil {
    // handle error
}
```

### Retrieve a client

```go
client, ok := mt.Get("tenant-a")
if !ok {
    // tenant not found
}
```

If you require strict existence:

```go
client := mt.MustGet("tenant-a")
```

## Automatic Cleanup Routine

The package also includes `StartCleanupRoutine` in `cleanup_routine.go`. It runs `CleanUp` periodically in a background goroutine.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

stopCleanup := multitenant.StartCleanupRoutine(
    ctx,
    mt,
    5*time.Minute,   // tick interval
    24*time.Hour,    // remove tenants inactive for 24h
)
defer stopCleanup()
```

Defaults:

- `tickInterval`: 1 hour if zero/negative
- `fromDuration`: 24 hours if zero/negative

## Notes

- This implementation is in-memory and process-local.
- Tenant lifecycle is based on `lastUsage`.
- `CleanUp` removes stale and invalid entries (`nil` entry/client).

## Example Usage (`cmd/example/multitenancy.go`)

For a complete example see:

- `cmd/example/multitenancy.go`

Key snippet (cache + per-tenant Vault credentials):

```go
c, ok := r.multiTenantClient.Get(tenant)
if ok {
	return c, nil
}

options := aruba.NewOptions().
	WithBaseURL(r.config.APIGateway).
	WithDefaultTokenIssuerURL().
	WithVaultCredentialsRepository(
		r.config.VaultAddress,
		r.config.KVMount,
		tenant, // tenant -> kvPath (e.g. ARU-297647)
		r.config.Namespace,
		r.config.RolePath,
		r.config.RoleID,
		r.config.RoleSecret,
	)

client, err := aruba.NewClient(options)
if err != nil {
	return nil, err
}
r.multiTenantClient.Add(tenant, client)
return client, nil
```
