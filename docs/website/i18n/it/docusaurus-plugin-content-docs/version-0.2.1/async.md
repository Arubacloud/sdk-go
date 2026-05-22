---
sidebar_position: 3
---

# Async / Await

La maggior parte delle operazioni API di Aruba Cloud sono **asincrone**: la chiamata HTTP ritorna rapidamente (con un `201 Created` o `200 OK`), ma la risorsa continua a transitare tra stati in background — `Creating` → `Active`, oppure `Updating` → `Active`, oppure `Deleting` → rimossa — per secondi o svariati minuti.

L'SDK espone tre livelli per gestire questa situazione:

| Livello | Quando usarlo |
|---------|---------------|
| `WaitUntilReady(ctx)` | 95% dei casi — blocca finché la risorsa è pronta (accetta `Active`, `Running`, `Stopped`, `NotUsed`, `Reserved`, `InUse`, `Used`) |
| `WaitUntilActive(ctx)` | Quando hai specificamente bisogno solo dello stato `Active` |
| `WaitUntilStates(ctx, []types.State{...}, opts...)` | Attendi uno o più stati specifici (es. `[]types.State{types.StateStopped}`) |
| `pkg/async.WaitFor` + `AsyncClient.Await` | Avanzato — avvia il polling in una goroutine in background, fai altro lavoro, raccogli il risultato in seguito |

---

## `WaitUntilReady`

Dopo qualsiasi `Create`, `Update` o `Get`, chiama `WaitUntilReady` sul wrapper restituito per bloccare finché la risorsa raggiunge uno qualsiasi dei 7 stati stabili: `Active`, `Running`, `Stopped`, `NotUsed`, `Reserved`, `InUse`, o `Used`.

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}

if err := vpc.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPC did not become ready: %v", err)
}
```

`WaitUntilReady` effettua il polling dell'API ripetutamente con un ritardo fisso. Quando la risorsa entra in uno **stato terminale di errore** noto (es. `"Error"`, `"Failed"`), ritorna immediatamente con un errore descrittivo invece di esaurire tutti i tentativi.

`WaitUntilActive` è disponibile quando hai specificamente bisogno dello stato `"Active"` — ad esempio dopo un'operazione di power-on.

Vedi la [Guida al Walkthrough API](./walkthrough) per esempi completi di Create + polling + Update + Delete.

### Personalizzare il comportamento di polling

Tre opzioni di chiamata permettono di sovrascrivere i valori predefiniti:

```go
if err := vpc.WaitUntilReady(ctx,
    aruba.WithRetries(30),              // max iterazioni di polling (default: 60)
    aruba.WithBaseDelay(5*time.Second), // ritardo fisso tra i poll (default: 10s)
    aruba.WithTimeout(3*time.Minute),   // scadenza rigida (default: 600s)
); err != nil {
    log.Fatalf("VPC did not become ready: %v", err)
}
```

Il limite effettivo è `min(retries × baseDelay, timeout)`. Con i valori predefiniti: `min(60×10s, 600s) = 600s`.

Per risorse a lunga esecuzione (Container Registry, DBaaS, KaaS) che possono richiedere 20–40 minuti, usa un budget maggiore:

```go
longWait := []aruba.WaitOption{
    aruba.WithTimeout(40 * time.Minute),
    aruba.WithRetries(240),
}
if err := reg.WaitUntilReady(ctx, longWait...); err != nil {
    log.Fatalf("ContainerRegistry did not become ready: %v", err)
}
```

---

## `WaitUntilStates`

Usa `WaitUntilStates` quando devi attendere uno o più stati specifici — ad esempio, attendere lo stato di stop dopo un'operazione di power-off:

```go
// Attendi che un Cloud Server si fermi completamente dopo PowerOff
if err := cs.WaitUntilStates(ctx, []types.State{types.StateStopped}); err != nil {
    log.Fatalf("Cloud Server did not stop: %v", err)
}
```

```go
// Attendi che un'istanza DBaaS finisca un aggiornamento in corso
if err := db.WaitUntilStates(ctx, []types.State{types.StateActive},
    aruba.WithRetries(120),
    aruba.WithBaseDelay(15*time.Second),
); err != nil {
    log.Fatalf("DBaaS did not return to Active after update: %v", err)
}
```

Si applica lo stesso comportamento di uscita anticipata per gli stati di errore terminali: se la risorsa raggiunge `"Error"` o `"Failed"` mentre si aspetta `"Stopped"`, la chiamata ritorna immediatamente con un errore che indica sia lo stato attuale che gli stati attesi.

`WaitUntilActive` e `WaitUntilReady` sono wrapper di convenienza attorno a `WaitUntilStates`:
- `WaitUntilActive(ctx, opts...)` — equivalente a `WaitUntilStates(ctx, []types.State{types.StateActive}, opts...)`
- `WaitUntilReady(ctx, opts...)` — equivalente a `WaitUntilStates(ctx, []types.State{types.StateActive, types.StateRunning, types.StateStopped, types.StateNotUsed, types.StateReserved, types.StateInUse, types.StateUsed}, opts...)`

---

## Accessor di Stato

Ogni wrapper che supporta il polling espone anche accessor di stato dettagliati. Puoi leggerli in qualsiasi momento dopo una chiamata `Create`, `Get`, `Update` o `List`:

| Metodo | Restituisce | Utilizzo tipico |
|--------|-------------|-----------------|
| `State()` | `types.State` — stato corrente | Logging, diramazione condizionale |
| `PreviousState()` | `types.State` — stato prima dell'ultima transizione | Post-mortem dopo un'attesa fallita |
| `FailureReason()` | `string` — testo di errore fornito dal server | Mostrare all'utente / log di alert |
| `IsDisabled()` | `bool` | Bloccare operazioni quando il server disabilita una risorsa |
| `DisableReasons()` | `[]string` | Spiegare perché una risorsa è disabilitata |

Un pattern comune — chiamare `WaitUntilReady` e, in caso di errore, allegare la motivazione di fallimento del server all'errore:

```go
if err := vpc.WaitUntilReady(ctx); err != nil {
    reason := vpc.FailureReason()
    if reason != "" {
        log.Fatalf("VPC failed: %v (reason: %s)", err, reason)
    }
    log.Fatalf("VPC polling failed: %v", err)
}
```

---

## Risorse che Supportano il Polling

I seguenti wrapper di risorse supportano `WaitUntilReady`, `WaitUntilActive`, `WaitUntilStates` e gli accessor di stato. Le risorse contrassegnate con un metodo wait speciale espongono un'ulteriore forma nominata.

| Risorsa | Wait speciale | Note |
|---------|---------------|------|
| `CloudServer` | — | `WaitUntilReady` → `Active` |
| `KaaS` | — | `WaitUntilReady` → `Active`; può richiedere 10–20 min |
| `ContainerRegistry` | — | `WaitUntilReady` → `Active`; può richiedere 20–40 min |
| `DBaaS` | — | `WaitUntilReady` → `Active`; può richiedere 5–15 min |
| `Database` | — | |
| `User` | — | |
| `Grant` | — | |
| `VPC` | — | |
| `Subnet` | — | |
| `SecurityGroup` | — | |
| `SecurityRule` | — | |
| `ElasticIP` | `WaitUntilNotUsed`, `WaitUntilUsed` | Delegano a `WaitUntilStates` |
| `BlockStorage` | `WaitUntilNotUsed`, `WaitUntilUsed` | Delegano a `WaitUntilStates` |
| `Snapshot` | — | |
| `StorageBackup` | — | |
| `StorageRestore` | — | |
| `VPCPeering` | — | |
| `VPCPeeringRoute` | — | |
| `VPNTunnel` | — | |
| `VPNRoute` | — | |
| `KMS` | — | |
| `Kmip` | `WaitUntilCertificateAvailable` | Waiter personalizzato (Family B — nessun `statusMixin`); fa il polling di `KmipResponse.Status` direttamente |

> **Il Progetto non supporta il polling.** È pronto in modo sincrono immediatamente dopo che `Create` ritorna — non è necessaria né disponibile alcuna chiamata `WaitUntilActive`.

---

## Avvertenze

### Wrapper idratato obbligatorio

`WaitUntilReady`, `WaitUntilActive` e `WaitUntilStates` funzionano solo su wrapper che sono stati **restituiti da una chiamata adapter** (`Create`, `Get`, `Update` o `List`). Chiamare uno qualsiasi di questi metodi su un builder di richiesta appena costruito restituisce:

```
WaitUntilStates: refresh callback not set; resource must be produced by an adapter (Create/Get/Update/List) to support polling
```

Usa sempre il wrapper restituito dalla chiamata API:

```go
// Corretto — vpc è stato restituito da Create
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, myVPC)
vpc.WaitUntilReady(ctx)

// Errato — myVPC è un builder di richiesta, non una risposta adapter
myVPC := aruba.NewVPC().Named("x")
myVPC.WaitUntilReady(ctx) // restituisce "refresh callback not set"
```

### Cadenza di polling costante

Il polling usa un **ritardo fisso** (senza backoff esponenziale). Se stai raggiungendo i limiti di rate dell'API, aumenta `WithBaseDelay` invece di aspettarti che l'SDK rallenti automaticamente.

### Cancellazione del context

Tutto il polling rispetta la scadenza e la cancellazione del `ctx`. Se il context scade durante il polling, la chiamata restituisce `ctx.Err()` (tipicamente `context.DeadlineExceeded` o `context.Canceled`).

---

## Avanzato: Polling in Background con `pkg/async`

`WaitUntilReady`, `WaitUntilActive` e `WaitUntilStates` bloccano la goroutine chiamante. Se devi **avviare più attese in modo concorrente**, o **fare il polling di una condizione arbitraria** (non solo uno stato di risorsa), usa direttamente il pacchetto `pkg/async` di livello inferiore.

`pkg/async` è un pacchetto pubblico — importalo insieme a `pkg/aruba`:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/async"
    "github.com/Arubacloud/sdk-go/pkg/types"
)
```

### `WaitFor` — avvia un future in background

`async.WaitFor` avvia immediatamente una goroutine di polling e restituisce un `*async.AsyncClient[T]` (un future). Chiami `.Await(ctx)` in seguito per bloccare e ottenere il risultato:

```go
// Avvia il polling di VPC1 e VPC2 in modo concorrente
futureVPC1 := async.DefaultWaitFor(ctx,
    func(ctx context.Context) (*types.Response[types.VPCResponse], error) {
        return arubaClient.FromNetwork().VPCs().Get(ctx, vpc1)
    },
    func(resp *types.Response[types.VPCResponse]) (bool, error) {
        if resp == nil || resp.Data == nil {
            return false, nil
        }
        var state types.State
        if resp.Data.Properties != nil && resp.Data.Properties.Status != nil &&
            resp.Data.Properties.Status.State != nil {
            state = *resp.Data.Properties.Status.State
        }
        return state == types.StateActive, nil
    },
)

futureVPC2 := async.DefaultWaitFor(ctx, /* stesso pattern per vpc2 */)

// Blocca e attendi entrambi i risultati
resp1, err1 := futureVPC1.Await(ctx)
resp2, err2 := futureVPC2.Await(ctx)
```

`DefaultWaitFor` usa i valori predefiniti del pacchetto: `DefaultRetries=60`, `DefaultBaseDelay=10s`, `DefaultTimeout=600s`. Usa `async.WaitFor(ctx, retries, baseDelay, timeout, call, check)` per sovrascriverli.

### Firma di `WaitFor`

```go
func WaitFor[T any](
    ctx         context.Context,
    retries     int,
    baseDelay   time.Duration,
    timeout     time.Duration,
    call        func(ctx context.Context) (*types.Response[T], error),
    check       func(*types.Response[T]) (bool, error),
) *AsyncClient[T]
```

- `call` — la funzione di polling, chiamata una volta per ogni iterazione.
- `check` — restituisce `(true, nil)` per segnalare il successo, `(true, error)` per segnalare un fallimento terminale, `(false, nil)` per continuare il polling.
- Se `check` è `nil`, qualsiasi `response.Data` non-nil viene trattato come successo.

### `AsyncClient.Await`

```go
func (f *AsyncClient[T]) Await(ctx context.Context) (*types.Response[T], error)
```

Blocca finché la goroutine in background invia il suo risultato oppure il `ctx` viene cancellato. Chiamate successive restituiscono il risultato **in cache** immediatamente — sicuro da chiamare più volte.

> `pkg/async` lavora direttamente con le struct wire di `pkg/types`. Questo è l'unico livello dell'SDK dove interagirai direttamente con `types.Response[T]` e i tipi `types.*Response`.

---

## Vedi Anche

- [Guida al Walkthrough API](./walkthrough) — esempio completo del ciclo di vita Create + `WaitUntilReady` + Update + Delete
- [Gestione delle Risposte](./response-handling) — come `*aruba.HTTPError` si propaga attraverso `WaitUntilReady` quando l'API restituisce 4xx/5xx
