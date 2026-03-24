# Multitenancy

Il package `pkg/multitenant` fornisce un registro in-memory tenant-to-client per l'SDK Aruba Cloud. E utile quando la tua applicazione gestisce piu tenant e ogni tenant richiede un proprio `aruba.Client`.

## Panoramica

Il package espone:

- Un'interfaccia `Multitenant` per creare, salvare, recuperare e pulire i client tenant
- Un'implementazione predefinita basata su mappa e mutex
- Un helper per routine di cleanup periodico dei tenant inattivi

File principali:

- `pkg/multitenant/multitenant.go`
- `pkg/multitenant/cleanup_routine.go`

## Interfaccia Principale

L'interfaccia `Multitenant` supporta queste operazioni:

- `New(tenant string) error`: crea client usando opzioni template
- `NewFromOptions(tenant string, options *aruba.Options) error`: crea client da opzioni esplicite
- `Add(tenant string, client aruba.Client)`: aggiunge un client gia inizializzato
- `Get(tenant string) (aruba.Client, bool)`: recupera client con flag di esistenza
- `MustGet(tenant string) aruba.Client`: recupera o termina il processo se assente
- `GetOrNil(tenant string) aruba.Client`: recupera o restituisce `nil`
- `CleanUp(from time.Duration)`: elimina tenant inattivi

## Creazione del Manager

### Manager vuoto

Usalo quando vuoi aggiungere client manualmente o tramite `NewFromOptions`:

```go
mt := multitenant.New()
```

### Manager con template

Usalo quando i tenant condividono una configurazione base comune:

```go
opts := aruba.DefaultOptions(clientID, clientSecret)
mt := multitenant.NewWithTemplate(opts)

// In seguito:
if err := mt.New("tenant-a"); err != nil {
    // gestisci errore
}
```

## Pattern di Utilizzo

### Aggiungere un client esistente

```go
client, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
if err != nil {
    // gestisci errore
}

mt.Add("tenant-a", client)
```

### Creare da opzioni specifiche tenant

```go
tenantOpts := aruba.DefaultOptions(tenantClientID, tenantClientSecret)
if err := mt.NewFromOptions("tenant-a", tenantOpts); err != nil {
    // gestisci errore
}
```

### Recuperare un client

```go
client, ok := mt.Get("tenant-a")
if !ok {
    // tenant non trovato
}
```

Se richiedi esistenza obbligatoria:

```go
client := mt.MustGet("tenant-a")
```

## Cleanup Automatico

Il package include anche `StartCleanupRoutine` in `cleanup_routine.go`. Esegue `CleanUp` periodicamente in una goroutine in background.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

stopCleanup := multitenant.StartCleanupRoutine(
    ctx,
    mt,
    5*time.Minute,   // intervallo di esecuzione
    24*time.Hour,    // rimuove tenant inattivi da 24h
)
defer stopCleanup()
```

Valori di default:

- `tickInterval`: 1 ora se zero/negativo
- `fromDuration`: 24 ore se zero/negativo

## Note

- Implementazione in-memory, locale al processo.
- Ciclo di vita tenant basato su `lastUsage`.
- `CleanUp` rimuove anche entry non valide (`nil` entry/client).

## Esempio di Utilizzo (`cmd/example/multitenancy.go`)

Per un esempio completo vedi:

- `cmd/example/multitenancy.go`

Snippet principale (cache + credenziali Vault per tenant):

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
		tenant, // tenant -> kvPath (es. ARU-297647)
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
