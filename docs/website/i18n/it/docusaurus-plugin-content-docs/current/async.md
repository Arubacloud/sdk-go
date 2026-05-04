---
sidebar_position: 3
---

# Async / Await

La maggior parte delle operazioni API di Aruba Cloud sono **asincrone**: la chiamata HTTP ritorna rapidamente (con un `201 Created` o `200 OK`), ma la risorsa continua a transitare tra stati in background — `Creating` → `Active`, oppure `Updating` → `Active`, oppure `Deleting` → rimossa — per secondi o svariati minuti.

L'SDK espone tre livelli per gestire questa situazione:

| Livello | Quando usarlo |
|---------|---------------|
| `WaitUntilActive(ctx)` | 95% dei casi — blocca finché la risorsa è pronta |
| `WaitUntilState(ctx, target)` | Attendi uno stato specifico (es. `"Stopped"`) |
| `pkg/async.WaitFor` + `AsyncClient.Await` | Avanzato — avvia il polling in una goroutine in background, fai altro lavoro, raccogli il risultato in seguito |

---

## `WaitUntilActive`

Dopo qualsiasi `Create`, `Update` o `Get`, chiama `WaitUntilActive` sul wrapper restituito per bloccare finché la risorsa raggiunge lo stato `"Active"`:

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}

if err := vpc.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

`WaitUntilActive` effettua il polling dell'API ripetutamente con un ritardo fisso. Quando la risorsa entra in uno **stato terminale di errore** noto (es. `"Error"`, `"Failed"`), ritorna immediatamente con un errore descrittivo invece di esaurire tutti i tentativi.

Vedi la [Guida al Walkthrough API](./walkthrough) per esempi completi di Create + polling + Update + Delete.

### Personalizzare il comportamento di polling

Tre opzioni di chiamata permettono di sovrascrivere i valori predefiniti:

```go
if err := vpc.WaitUntilActive(ctx,
    aruba.WithRetries(30),              // max iterazioni di polling (default: 60)
    aruba.WithBaseDelay(5*time.Second), // ritardo fisso tra i poll (default: 10s)
    aruba.WithTimeout(3*time.Minute),   // scadenza rigida (default: 600s)
); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

Il limite effettivo è `min(retries × baseDelay, timeout)`. Con i valori predefiniti: `min(60×10s, 600s) = 600s`.

---

## `WaitUntilState`

Usa `WaitUntilState` quando devi attendere uno stato diverso da `"Active"`:

```go
// Attendi che un Cloud Server si fermi completamente dopo PowerOff
if err := cs.WaitUntilState(ctx, "Stopped"); err != nil {
    log.Fatalf("Cloud Server did not stop: %v", err)
}
```

```go
// Attendi che un'istanza DBaaS finisca un aggiornamento in corso
if err := db.WaitUntilState(ctx, "Active",
    aruba.WithRetries(120),
    aruba.WithBaseDelay(15*time.Second),
); err != nil {
    log.Fatalf("DBaaS did not return to Active after update: %v", err)
}
```

Si applica lo stesso comportamento di uscita anticipata per stati di errore terminali: se la risorsa raggiunge `"Error"` o `"Failed"` mentre si aspetta `"Stopped"`, la chiamata ritorna immediatamente con un errore che indica sia lo stato attuale che quello atteso.

---

## Accessor di Stato

Ogni wrapper che supporta il polling espone anche accessor di stato dettagliati. Puoi leggerli in qualsiasi momento dopo una chiamata `Create`, `Get`, `Update` o `List`:

| Metodo | Restituisce | Utilizzo tipico |
|--------|-------------|-----------------|
| `State()` | `string` — stato corrente | Logging, diramazione condizionale |
| `PreviousState()` | `string` — stato prima dell'ultima transizione | Post-mortem dopo un'attesa fallita |
| `FailureReason()` | `string` — testo di errore fornito dal server | Mostrare all'utente / log di alert |
| `IsDisabled()` | `bool` | Bloccare operazioni quando il server disabilita una risorsa |
| `DisableReasons()` | `[]string` | Spiegare perché una risorsa è disabilitata |

Un pattern comune — chiamare `WaitUntilActive` e, in caso di errore, allegare la motivazione di fallimento del server all'errore:

```go
if err := vpc.WaitUntilActive(ctx); err != nil {
    reason := vpc.FailureReason()
    if reason != "" {
        log.Fatalf("VPC failed: %v (reason: %s)", err, reason)
    }
    log.Fatalf("VPC polling failed: %v", err)
}
```

---

## Risorse che Supportano il Polling

I seguenti wrapper di risorse incorporano il mixin di polling e supportano `WaitUntilActive`, `WaitUntilState` e gli accessor di stato:

- **Compute**: `CloudServer`
- **Container**: `KaaS`, `ContainerRegistry`
- **Database**: `DBaaS`
- **Network**: `VPC`, `Subnet`, `SecurityGroup`, `SecurityRule`, `ElasticIP`
- **Security**: `KMS`, `Kmip`
- **Storage**: `BlockStorage`, `Snapshot`

> **Il Progetto non supporta il polling.** È pronto in modo sincrono immediatamente dopo che `Create` ritorna — non è necessaria né disponibile alcuna chiamata `WaitUntilActive`.

---

## Avvertenze

### Wrapper idratato obbligatorio

`WaitUntilActive` e `WaitUntilState` funzionano solo su wrapper che sono stati **restituiti da una chiamata adapter** (`Create`, `Get`, `Update` o `List`). Chiamare uno dei due metodi su un builder di richiesta appena costruito restituisce:

```
WaitUntilState: refresh callback not set; resource must be produced by an adapter (Create/Get/Update/List) to support polling
```

Usa sempre il wrapper restituito dalla chiamata API:

```go
// Corretto — vpc è stato restituito da Create
vpc, err := arubaClient.FromNetwork().VPCs().Create(ctx, myVPC)
vpc.WaitUntilActive(ctx)

// Errato — myVPC è un builder di richiesta, non una risposta adapter
myVPC := aruba.NewVPC().WithName("x")
myVPC.WaitUntilActive(ctx) // restituisce "refresh callback not set"
```

### Cadenza di polling costante

Il polling usa un **ritardo fisso** (senza backoff esponenziale). Se stai raggiungendo i limiti di rate dell'API, aumenta `WithBaseDelay` invece di aspettarti che l'SDK rallenti automaticamente.

### Cancellazione del context

Tutto il polling rispetta la scadenza e la cancellazione del `ctx`. Se il context scade durante il polling, la chiamata restituisce `ctx.Err()` (tipicamente `context.DeadlineExceeded` o `context.Canceled`).

---

## Avanzato: Polling in Background con `pkg/async`

`WaitUntilActive` e `WaitUntilState` bloccano la goroutine chiamante. Se devi **avviare più attese in modo concorrente**, o **fare il polling di una condizione arbitraria** (non solo uno stato di risorsa), usa direttamente il pacchetto `pkg/async` di livello inferiore.

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
        state := ""
        if resp.Data.Properties != nil && resp.Data.Properties.Status != nil &&
            resp.Data.Properties.Status.State != nil {
            state = *resp.Data.Properties.Status.State
        }
        return state == "Active", nil
    },
)

futureVPC2 := async.DefaultWaitFor(ctx, /* stesso pattern per vpc2 */)

// Blocca e attendi entrambi i risultati
resp1, err1 := futureVPC1.Await(ctx)
resp2, err2 := futureVPC2.Await(ctx)
```

`DefaultWaitFor` usa gli stessi valori predefiniti di `WaitUntilActive`: 60 tentativi, 10s di ritardo, 600s di timeout. Usa `async.WaitFor(ctx, retries, baseDelay, timeout, call, check)` per sovrascriverli.

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

- [Guida al Walkthrough API](./walkthrough) — esempio completo del ciclo di vita Create + `WaitUntilActive` + Update + Delete
- [Gestione delle Risposte](./response-handling) — come `*aruba.HTTPError` si propaga attraverso `WaitUntilActive` quando l'API restituisce 4xx/5xx
