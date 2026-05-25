# Lavorare a Basso Livello

## Perché questa pagina esiste

L'SDK è progettato attorno a un **principio di import singolo**: importare solo `github.com/Arubacloud/sdk-go/pkg/aruba` copre circa il 99,9% dei casi d'uso reali.

La superficie wrapper — accessor tipizzati, helper `WaitUntil*`, `RawJSON()` / `RawYAML()` — è progettata in modo che raramente sia necessario accedere direttamente alle struct wire sottostanti. Per lo 0,1% dei casi in cui è necessario farlo, questa pagina raccoglie tutte le vie di uscita che richiedono un secondo import.

> Questi pattern sono **intenzionali e supportati** — sono vie di uscita, non workaround. Quando incontri un caso non coperto dalla superficie wrapper, usa questi invece di aprire una feature request.

---

## Accedere ai campi wire non promossi con `Raw()`

Ogni wrapper espone un metodo `Raw()` che restituisce la struct `*types.XxxResponse` sottostante. Usalo quando hai bisogno di un campo non ancora promosso alla superficie del wrapper:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

raw := vpc.Raw()                          // *types.VPCResponse
fmt.Println(raw.Properties.IsDefault)     // campo non promosso al wrapper
```

Per le liste, `Raw()` restituisce `any` e devi fare il type assert al tipo lista concreto:

```go
vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil { /* … */ }

raw := vpcList.Raw().(*types.VPCList)     // richiede l'import di pkg/types
fmt.Println("server total:", raw.Total)   // equivalente a vpcList.Total() — mostrato per illustrazione
fmt.Println("self link:", raw.Self)
```

> **Preferisci gli accessor del wrapper per la serializzazione.** Se il tuo obiettivo è l'output JSON o YAML, usa `vpc.RawJSON()` / `vpc.RawYAML()` (o i metodi equivalenti su `List[T]`) — senza bisogno di importare `pkg/types`.
>
> ```go
> fmt.Println(string(vpcList.RawJSON()))  // JSON senza pkg/types
> fmt.Println(string(vpcList.RawYAML()))  // YAML senza pkg/types
> ```

---

## Ispezionare gli errori di validazione strutturati

`*aruba.HTTPError` è il tipo di errore per tutte le risposte 4xx/5xx. Il suo campo `ErrResp` è un `*types.ErrorResponse` e contiene dettagli strutturati RFC 7807 — inclusa una slice `[]types.ValidationError` per gli errori 400 a livello di campo:

```go
import (
    "errors"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

_, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) && httpErr.ErrResp != nil {
        fmt.Printf("title:  %s\n", derefStr(httpErr.ErrResp.Title))
        fmt.Printf("detail: %s\n", derefStr(httpErr.ErrResp.Detail))

        // Errori di validazione a livello di campo — richiedono types.ValidationError
        for _, ve := range httpErr.ErrResp.Errors {
            fmt.Printf("  campo %s: %s\n", ve.Field, ve.Message)
        }

        // TraceID per le richieste di supporto
        fmt.Printf("trace-id: %s\n", derefStr(httpErr.ErrResp.TraceID))
    }
}
```

### `MetadataValidationError`

Un `*types.MetadataValidationError` viene restituito (insieme a un wrapper non-nil) quando una risposta API manca dei campi di metadati obbligatori (`id` o `uri`). Usa `errors.As` per rilevarlo:

```go
import (
    "errors"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

result, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var mvErr *types.MetadataValidationError
    if errors.As(err, &mvErr) {
        // result è non-nil e parzialmente idratato; mvErr elenca i campi mancanti
        fmt.Printf("metadati incompleti: %v\n", mvErr)
        fmt.Printf("ID finora: %s\n", result.ID())
    }
}
```

---

## Iterare `LinkedResources()`

Ogni wrapper di risorsa espone `LinkedResources()` che restituisce `[]types.LinkedResource`. Ogni elemento ha un `URI string` e un `StrictCorrelation bool`:

```go
import (
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

for _, lr := range vpc.LinkedResources() {
    fmt.Println("linked URI:", lr.URI)
    if lr.StrictCorrelation {
        fmt.Println("  → correlazione stretta (lifecycle-linked)")
    }
}
```

> Se hai bisogno solo delle stringhe URI — ad esempio per passarle a un'altra chiamata SDK come `aruba.URI(lr.URI)` — non hai bisogno di `types.LinkedResource`:
>
> ```go
> for _, lr := range vpc.LinkedResources() {
>     ref := aruba.URI(lr.URI)   // senza import di pkg/types
>     _ = ref
> }
> ```

---

## Ispezionare i body delle richieste prima dell'invio

Ogni wrapper espone `RawRequest()` che restituisce la struct di richiesta a livello wire (`*types.XxxRequest`). È utile per il debug o per passare la richiesta a un altro strumento:

```go
import (
    "encoding/json"
    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

vpc := aruba.NewVPC().
    IntoProject(proj).
    Named("my-vpc").
    InRegion(aruba.RegionITBGBergamo).
    AsDefault()

req := vpc.RawRequest()               // types.VPCRequest — richiede import di pkg/types
b, _ := json.MarshalIndent(req, "", "  ")
fmt.Println(string(b))
```

---

## Polling in background con `pkg/async` {#background-polling-with-pkgasync}

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

## Cosa NON richiede `pkg/types`

I seguenti elementi sono tutti disponibili con un singolo import di `pkg/aruba` — senza bisogno di un secondo import:

| Cosa ti serve | Superficie `pkg/aruba` |
|---|---|
| Costanti di stato (`StateActive`, `StateStopped`, …) | `aruba.StateActive`, `aruba.StateStopped`, … |
| Costanti regione / zona | `aruba.RegionITBGBergamo`, `aruba.ZoneITBG1`, … |
| Periodo di fatturazione | `aruba.BillingPeriodHour`, `aruba.BillingPeriodMonth`, … |
| Tutti gli enum compute, storage, network, security | Vedi `pkg/aruba/aliases.go` |
| Attendere transizioni di stato | `wrapper.WaitUntilReady(ctx)`, `WaitUntilActive`, `WaitUntilStates` |
| Serializzare una risposta in JSON / YAML | `wrapper.RawJSON()`, `wrapper.RawYAML()` |
| Introspezione dell'envelope HTTP | `wrapper.StatusCode()`, `.Headers()`, `.RawHTTP()`, `.RawError()` |
| Paginazione | `list.Total()`, `.HasNext()`, `.Next(ctx)`, `.All(ctx, yield)` |
| Dettagli errore HTTP | `*aruba.HTTPError` — `StatusCode`, `ErrResp.Title`, `ErrResp.Detail`, `ErrResp.TraceID` |

---

## Vedi Anche

- [Gestione delle Risposte](./response-handling) — errori HTTP tipizzati, accessor envelope, `RawJSON`/`RawYAML`
- [Async / Await](./async) — `WaitUntilReady`, `WaitUntilStates`, opzioni di polling
- [Guida al Walkthrough API](./walkthrough) — esempio completo del ciclo di vita Create + polling + Update + Delete

---

```go
// Helper usato negli esempi sopra
func derefStr(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
```
