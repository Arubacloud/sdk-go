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
| `WaitUntilGone(ctx)` | Dopo `Delete` — blocca finché il `Get` della risorsa restituisce HTTP 404 (completamente rimossa) |
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
if err := db.WaitUntilActive(ctx,
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

## `WaitUntilGone`

Usa `WaitUntilGone` dopo una chiamata `Delete` per bloccare finché la risorsa è completamente rimossa — ovvero finché il suo `Get` restituisce HTTP 404:

```go
if err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet); err != nil {
    log.Printf("Delete subnet: %v", err)
} else if err := subnet.WaitUntilGone(ctx); err != nil {
    log.Printf("Subnet not gone: %v", err)
}
```

`WaitUntilGone` è disponibile su ogni wrapper di risorsa che supporta il polling (vedi [Risorse che Supportano il Polling](#risorse-che-supportano-il-polling) in basso). Accetta le stesse `WaitOption` di `WaitUntilReady`. Qualsiasi errore da `Get` diverso da HTTP 404 viene trattato come transitorio e riprovato; un 404 segnala il successo.

`Project` non ha supporto al polling e quindi nessun `WaitUntilGone`. Viene eliminato per ultimo, senza figli da attendere.

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

I seguenti wrapper di risorse supportano `WaitUntilReady`, `WaitUntilActive`, `WaitUntilStates`, `WaitUntilGone` e gli accessor di stato. Le risorse contrassegnate con un metodo wait speciale espongono un'ulteriore forma nominata.

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

`WaitUntilReady`, `WaitUntilActive`, `WaitUntilStates` e `WaitUntilGone` funzionano solo su wrapper che sono stati **restituiti da una chiamata adapter** (`Create`, `Get`, `Update` o `List`). Chiamare uno qualsiasi di questi metodi su un builder di richiesta appena costruito restituisce:

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

## Avanzato: polling concorrente e personalizzato

`WaitUntilReady`, `WaitUntilActive` e `WaitUntilStates` bloccano la goroutine chiamante. Quando devi **avviare più attese in modo concorrente**, o **fare il polling di una condizione arbitraria** (non solo uno stato di risorsa), scendi al livello di `pkg/async`. Quel layer lavora direttamente con `*types.Response[T]` ed è documentato separatamente — consulta [Lavorare a Basso Livello](./working-at-low-level#background-polling-with-pkgasync).

---

## Vedi Anche

- [Guida al Walkthrough API](./walkthrough) — esempio completo del ciclo di vita Create + `WaitUntilReady` + Update + Delete
- [Gestione delle Risposte](./response-handling) — come `*aruba.HTTPError` si propaga attraverso `WaitUntilReady` quando l'API restituisce 4xx/5xx
- [Lavorare a Basso Livello](./working-at-low-level) — polling in background con `pkg/async`, accesso ai campi wire non promossi
