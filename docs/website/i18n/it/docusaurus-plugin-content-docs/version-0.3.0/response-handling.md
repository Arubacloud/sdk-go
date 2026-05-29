# Guida alla Gestione delle Risposte

## Panoramica

Il layer wrapper dell'SDK gestisce automaticamente il parsing delle risposte e la segnalazione degli errori. Ogni metodo CRUD restituisce o un wrapper popolato (in caso di successo) o un errore. Raramente è necessario ispezionare direttamente l'envelope HTTP grezzo — ma gli strumenti per farlo sono sempre disponibili.

## Pattern Base

Ogni metodo wrapper restituisce `(wrapper, error)`. L'errore è non-nil sia per errori di rete che per errori a livello API (4xx / 5xx).

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx,
    aruba.URI("/projects/<projectID>/providers/Aruba.Network/vpcs/<vpcID>"),
)
if err != nil {
    log.Fatalf("Get VPC failed: %v", err)
}
fmt.Printf("VPC: %s (state: %s)\n", vpc.Name(), vpc.State())
```

## Errori HTTP Tipizzati

Quando l'API restituisce una risposta 4xx o 5xx, l'SDK la racchiude in `*aruba.HTTPError`. Usa `errors.As` per ispezionare il codice di stato e il corpo dell'errore strutturato:

```go
import "errors"

vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("API error %d: %s\n", httpErr.StatusCode, httpErr.Error())
        if httpErr.ErrResp != nil {
            fmt.Printf("  title:  %s\n", derefStr(httpErr.ErrResp.Title))
            fmt.Printf("  detail: %s\n", derefStr(httpErr.ErrResp.Detail))
            for _, ve := range httpErr.ErrResp.Errors {
                fmt.Printf("  field %s: %s\n", ve.Field, ve.Message)
            }
        }
    } else {
        // Errore di rete, timeout del contesto, ecc.
        log.Fatalf("Request failed: %v", err)
    }
}
```

## Pattern Completo di Gestione degli Errori

```go
proj, err := arubaClient.FromProject().Get(ctx, ref)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        switch httpErr.StatusCode {
        case 404:
            return fmt.Errorf("project not found")
        case 400:
            return fmt.Errorf("bad request: %s", derefStr(httpErr.ErrResp.Detail))
        default:
            return fmt.Errorf("API error (%d): %s", httpErr.StatusCode, httpErr.Error())
        }
    }
    return fmt.Errorf("request failed: %w", err)
}
// proj è popolato — usalo direttamente
fmt.Printf("Project: %s (tags: %v)\n", proj.Name(), proj.Tags())
```

## Accessori dell'Envelope HTTP

Ogni wrapper prodotto da una chiamata Create / Get / Update / List espone il suo envelope HTTP grezzo:

```go
// Dopo qualsiasi chiamata CRUD:
proj, err := arubaClient.FromProject().Create(ctx, p)
// …

fmt.Println("Status:", proj.StatusCode())
fmt.Println("Headers:", proj.Headers())
rawResp, rawBody := proj.RawHTTP()
fmt.Println("Raw body:", string(rawBody))
fmt.Println("HTTP status:", rawResp.StatusCode)
fmt.Println("Error body (if any):", proj.RawError())
```

## Accesso alla Risposta Wire Grezza

Ogni wrapper ha un metodo `Raw()` che restituisce lo struct di risposta tipizzato sottostante da `pkg/types`. È utile per i campi non ancora promossi alla superficie del wrapper:

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, ref)
if err != nil { /* … */ }

raw := vpc.Raw()                         // struct wire sottostante
fmt.Println(raw.Properties.IsDefault)    // campo non sul wrapper
```

### Convenienza JSON / YAML

Per flag CLI come `--output json` / `--output yaml`, ogni wrapper espone slice di byte già serializzate:

```go
fmt.Println(string(vpc.RawJSON()))   // payload codificato in JSON
fmt.Println(string(vpc.RawYAML()))   // payload codificato in YAML
```

Restituisce `nil` se il wrapper non è ancora stato popolato (receiver a valore zero).

## Risposte List

`List[T]` espone la stessa superficie di introspezione dei wrapper per risorsa singola, senza mai dover importare `pkg/types`:

```go
vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil { /* … */ }

// Paginazione e conteggi — accessor tipizzati sul wrapper.
fmt.Println("server total:", vpcList.Total())
if vpcList.HasNext() {
    nextPage, _ := vpcList.Next(ctx)
    _ = nextPage
}

// Envelope HTTP — stessi accessor dei wrapper per risorsa singola.
fmt.Println("status:", vpcList.StatusCode())
fmt.Println("trace-id:", vpcList.Headers().Get("X-Trace-Id"))
_, body := vpcList.RawHTTP()
fmt.Println("raw body bytes:", len(body))
```

### Convenienza JSON / YAML

Anche `List[T]` espone `RawJSON()` e `RawYAML()` per il payload della lista tipizzata:

```go
fmt.Println(string(vpcList.RawJSON()))   // payload codificato in JSON
fmt.Println(string(vpcList.RawYAML()))   // payload codificato in YAML
```

Restituisce `nil` quando la lista non ha payload (`Raw() == nil`).

> **Accedere a campi non promossi.** Se hai bisogno di un campo non esposto dalla superficie del wrapper, vedi [Lavorare a Basso Livello](./working-at-low-level) — copre il cast alla struct wire tipizzata e gli altri escape hatch che richiedono l'import di `pkg/types`.

## Errori in Fase di Setter

I setter del builder fluente non restituiscono mai errori — li registrano internamente. L'errore emerge alla prima chiamata `Create` o `Update`. Puoi anche controllarlo in anticipo:

```go
rule := aruba.NewSecurityRule().
    InSecurityGroup(sg).
    TargetingCIDR("0.0.0.0/0").
    TargetingSecurityGroup(otherSG) // setter in conflitto — registra un errore

if err := rule.Err(); err != nil {
    log.Fatalf("Bad rule configuration: %v", err)
}
```

## Best Practice

1. **Controlla sempre `err` prima** — copre sia gli errori di rete che quelli API.
2. **Usa `errors.As(err, &httpErr)`** per ottenere dettagli strutturati sugli errori 4xx/5xx.
3. **Controlla `httpErr.ErrResp.Errors`** per messaggi di validazione a livello di campo sugli errori 400.
4. **Usa `httpErr.ErrResp.TraceID`** quando apri una richiesta di supporto.
5. **Usa `.Raw()` con parsimonia** — preferisci gli accessor tipizzati del wrapper.
6. **Controlla `wrapper.Err()` prima di Create/Update** quando la catena di builder è lunga.

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
