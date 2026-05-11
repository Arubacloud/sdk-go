---
sidebar_position: 2
---

# Guida al Walkthrough API

L'SDK Go di Aruba Cloud fornisce un singolo import — `github.com/Arubacloud/sdk-go/pkg/aruba` — che espone un'API fluente con pattern builder per ogni risorsa cloud. Si costruisce la descrizione della risorsa con una catena `aruba.NewX()`, la si passa al metodo del client appropriato (`Create`, `Get`, `Update`, `Delete` o `List`), e si lavora con il wrapper tipizzato restituito.

Le risorse sono organizzate in un **Progetto**, e le risorse figlio referenziano i propri genitori tramite l'interfaccia `aruba.Ref`. Non è mai necessario estrarre o passare manualmente stringhe di ID grezzi: si passa direttamente il wrapper idratato (restituito da `Create` o `Get`) come parametro `Ref` ai metodi builder come `IntoProject(proj)`, `IntoVPC(vpc)` o `IntoSecurityGroup(sg)`.

Questa pagina illustra il ciclo CRUD completo su un esempio minimale — Project + VPC + Subnet. Ogni altra risorsa segue esattamente la stessa struttura. Vedi [Risorse](./resources) per snippet pronti all'uso per tutte le risorse supportate.

---

## 1. Inizializzare il Client

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/Arubacloud/sdk-go/pkg/aruba"
)

func main() {
    arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()
}
```

`aruba.NewClient` accetta un valore `*aruba.Options`. `aruba.DefaultOptions(clientID, clientSecret)` è il modo più rapido per iniziare; vedi [Opzioni di Configurazione](./options) per credenziali Vault, caching token Redis, logger personalizzati e altro.

Il `aruba.Client` restituito è una facciata che espone sub-client specifici per dominio:
`FromProject()`, `FromAudit()`, `FromCompute()`, `FromContainer()`, `FromDatabase()`, `FromMetric()`, `FromNetwork()`, `FromSchedule()`, `FromSecurity()`, `FromStorage()`.

---

## 2. Provisioning delle Risorse

Le risorse vengono create inline: si costruisce la richiesta con `aruba.NewX()` e la si passa direttamente a `Create`. Il wrapper restituito porta l'ID e l'URI della risorsa — lo si passa come `Ref` ai builder delle risorse figlio.

### Progetto

Il Progetto è il contenitore di livello superiore. Ogni altra risorsa appartiene a un progetto. È pronto in modo sincrono dopo che `Create` ritorna — non è necessario il polling.

```go
proj, err := arubaClient.FromProject().Create(
    ctx,
    aruba.NewProject().
        WithName("my-project").
        WithDescription("Creato tramite l'SDK Go di Aruba Cloud").
        AddTag("go-sdk").
        WithDefault(false))
if err != nil {
    log.Fatalf("Create project: %v", err)
}
fmt.Printf("✓ Progetto creato: %s (ID: %s)\n", proj.Name(), proj.ID())
```

### VPC

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(
    ctx,
    aruba.NewVPC().
        IntoProject(proj).
        WithName("my-vpc").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithDefault(false).
        WithPreset(false))
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}
fmt.Printf("✓ VPC creata: %s\n", vpc.Name())

// La maggior parte delle risorse è asincrona — attendi che raggiungano lo stato Active.
// Vedi "7. Attendere la Disponibilità" per le opzioni e i dettagli.
if err := vpc.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

`IntoProject(proj)` accetta qualsiasi `aruba.Ref` — lega lo scope del progetto senza richiedere l'estrazione di un ID stringa grezzo.

### Subnet

```go
subnet, err := arubaClient.FromNetwork().Subnets().Create(
    ctx,
    aruba.NewSubnet().
        IntoVPC(vpc).
        WithName("my-subnet").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithType(string(aruba.SubnetTypeAdvanced)).
        WithDefault(false).
        WithCIDR("192.168.1.0/25").
        WithDHCP(aruba.NewSubnetDHCP().
            Enabled().
            WithRange("192.168.1.10", 50).
            AddRoute("10.0.0.0/8", "192.168.1.1").
            AddDNS("8.8.8.8").
            AddDNS("8.8.4.4")))
if err != nil {
    log.Fatalf("Create subnet: %v", err)
}
fmt.Printf("✓ Subnet creata: %s (CIDR: %s)\n", subnet.Name(), subnet.CIDR())

if err := subnet.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Subnet did not become Active: %v", err)
}
```

`aruba.NewSubnetDHCP()` è un sub-builder per la configurazione DHCP. Si allega alla subnet con `WithDHCP(...)`.

`WithType` accetta `string(aruba.SubnetTypeBasic)` o `string(aruba.SubnetTypeAdvanced)`.

> Ogni altra risorsa — Security Group, Elastic IP, Block Storage, Cloud Server, cluster KaaS, istanze DBaaS e altro — segue esattamente la stessa struttura `NewX()` → `IntoParent(ref)` → `Create(ctx, ...)` → `WaitUntilActive(ctx)`. Vedi [Risorse](./resources) per l'elenco completo con snippet pronti all'uso.

---

## 3. Aggiornare una Risorsa Esistente

Prima si recupera la risorsa, poi la si muta tramite setter, poi si chiama `Update`. Il wrapper di risposta da `Get` porta tutto lo stato interno (URI padre, riferimenti di rete, ecc.) che viene trasferito automaticamente nella richiesta `Update`.

```go
// Recupera
vpc, err = arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
if err != nil {
    log.Fatalf("Get VPC: %v", err)
}

// Muta
vpc.WithName("my-vpc-updated").
    ReplaceTags("network", "updated")

// Aggiorna
updated, err := arubaClient.FromNetwork().VPCs().Update(ctx, vpc)
if err != nil {
    log.Fatalf("Update VPC: %v", err)
}
fmt.Printf("✓ VPC aggiornata: %s\n", updated.Name())
```

> **Importante**: Chiama sempre `Get` prima di `Update`. Chiamare `Update` su un wrapper appena costruito (senza un precedente `Create` o `Get`) restituisce un errore: `"Update: resource has no ID"`.

---

## 4. Elencare Risorse Esistenti

`List` accetta un `Ref` genitore e restituisce un `*aruba.List[T]`. Si iterano gli elementi con `Items()`:

```go
list, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
if err != nil {
    log.Fatalf("List VPCs: %v", err)
}
fmt.Println("Totale VPC:", list.Total())
for _, v := range list.Items() {
    fmt.Println("-", v.Name(), v.ID())
}
```

Gli elementi nella lista sono wrapper leggeri — portano l'ID e l'URI della risorsa, così puoi passarli direttamente a `Get`, `Update` o `Delete` come `Ref`:

```go
for _, v := range list.Items() {
    full, err := arubaClient.FromNetwork().VPCs().Get(ctx, v)
    // full ha tutti i campi popolati
}
```

Per il filtraggio lato server, ordinamento e paginazione vedi [Filtri](./filters).

---

## 5. Ottenere una Risorsa Specifica

Usa `Get` quando hai un `Ref` (un wrapper idratato, o un elemento di `*aruba.List[T]`):

```go
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
if err != nil {
    log.Fatalf("Get VPC: %v", err)
}
```

### L'escape hatch `aruba.URI(…)`

Quando hai solo un identificatore di risorsa come stringa — ad esempio, letto da una variabile d'ambiente o da una configurazione esterna — avvolgilo in `aruba.URI(…)` per soddisfare l'interfaccia `aruba.Ref`:

```go
projectID := os.Getenv("PROJECT_ID")

// Bootstrap di un wrapper tipizzato da un ID stringa
proj, err := arubaClient.FromProject().Get(ctx, aruba.URI("/projects/"+projectID))
if err != nil {
    log.Fatalf("Get project: %v", err)
}

// Ora proj è completamente idratato — usalo come Ref per le risorse figlio
vpcs, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
```

`aruba.URI(s)` restituisce un `Ref` leggero che l'SDK usa per estrarre gli ID degli antenati dai segmenti di percorso URI. Qualsiasi URI di risorsa valido funziona — l'SDK lo analizza internamente.

---

## 6. Eliminazione (Ordine Inverso)

Elimina i figli prima dei genitori. L'API Aruba Cloud restituisce **HTTP 400** quando si tenta di eliminare un genitore che ha ancora risorse figlio attive o in fase di eliminazione — non 409/422. Il pattern sicuro è emettere ogni delete sul figlio, poi attendere che la risorsa sia completamente sparita (HTTP 404) prima di salire nella catena delle dipendenze.

Usa `pkg/async.WaitFor` per attendere il 404 — centralizza la logica di retry/timeout/cadenza:

```go
import (
    "errors"
    "net/http"

    "github.com/Arubacloud/sdk-go/pkg/aruba"
    "github.com/Arubacloud/sdk-go/pkg/async"
    "github.com/Arubacloud/sdk-go/pkg/types"
)

// waitUntilGone blocca finché il Get della risorsa restituisce HTTP 404.
func waitUntilGone(ctx context.Context, poll func(context.Context) error) error {
    const gone = "gone"
    fut := async.DefaultWaitFor(ctx,
        func(ctx context.Context) (*types.Response[string], error) {
            err := poll(ctx)
            if err == nil {
                return &types.Response[string]{}, nil // esiste ancora
            }
            var httpErr *aruba.HTTPError
            if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
                return &types.Response[string]{Data: &[]string{gone}[0]}, nil // eliminata
            }
            return nil, err // transitorio — riprova
        },
        func(resp *types.Response[string]) (bool, error) {
            return resp != nil && resp.Data != nil, nil
        },
    )
    _, err := fut.Await(ctx)
    return err
}
```

Poi elimina in ordine inverso delle dipendenze, attendendo che ogni figlio sparisca completamente prima di eliminare il suo genitore:

```go
// subnet → VPC → progetto
if err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet); err != nil {
    log.Printf("Delete subnet: %v", err)
} else {
    waitUntilGone(ctx, func(ctx context.Context) error {
        _, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
        return err
    })
}

if err := arubaClient.FromNetwork().VPCs().Delete(ctx, vpc); err != nil {
    log.Printf("Delete VPC: %v", err)
} else {
    waitUntilGone(ctx, func(ctx context.Context) error {
        _, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
        return err
    })
}

if err := arubaClient.FromProject().Delete(ctx, proj); err != nil {
    log.Printf("Delete project: %v", err)
}
```

`Delete` accetta qualsiasi `aruba.Ref` — puoi passare il wrapper idratato direttamente o `aruba.URI(…)` se hai solo il percorso.

Per una sequenza completa di teardown (Security Rule → Security Group → Subnet → VPC → Cloud Server → Block Storage → Progetto) vedi l'[Esempio Completo](#esempio-completo) in basso.

---

## 7. Attendere la Disponibilità

La maggior parte delle operazioni cloud — Create, Update, operazioni di scaling — sono **asincrone**: la chiamata HTTP ritorna rapidamente, ma la risorsa continua a transitare tra stati (`Creating` → `Active`, `Updating` → `Active`) per secondi o minuti in background.

Il metodo `WaitUntilActive` su qualsiasi wrapper che incorpora `statusMixin` blocca finché la risorsa raggiunge lo stato `"Active"` (o restituisce un errore in caso di fallimento terminale):

```go
if err := vpc.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

Tre `WaitOption` permettono di sovrascrivere i valori predefiniti (60 tentativi × 10 s di ritardo base × 600 s di scadenza rigida):

```go
if err := vpc.WaitUntilActive(ctx,
    aruba.WithRetries(30),              // max iterazioni di polling (default: 60)
    aruba.WithBaseDelay(5*time.Second), // ritardo fisso tra i poll (default: 10s)
    aruba.WithTimeout(3*time.Minute),   // scadenza rigida (default: 600s)
); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
```

Per `WaitUntilStates` (qualsiasi stato target, non solo `"Active"`), gli accessor di stato (`State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`), e il primitivo di basso livello `pkg/async.WaitFor` per il polling concorrente, vedi la guida [Async / Await](./async).

---

## Avvertenze

### Gli errori dei setter sono differiti

I setter del builder non restituiscono mai un errore — lo registrano nel wrapper. L'errore viene restituito dalla prima chiamata `Create` o `Update`, oppure puoi controllarlo in anticipo:

```go
rule := aruba.NewSecurityRule().
    IntoSecurityGroup(sg).
    WithTargetCIDR("0.0.0.0/0").
    WithTargetSecurityGroup(otherSG) // in conflitto — registrato come errore

if err := rule.Err(); err != nil {
    log.Fatalf("Bad rule config: %v", err)
}
```

> **Avvertenza**: `WithTargetCIDR` e `WithTargetSecurityGroup` si escludono a vicenda. Impostarli entrambi registra un errore al momento del setter che emerge su `Create`.

### `WaitUntilActive` richiede un wrapper idratato

Chiamare `WaitUntilActive` su un wrapper costruito manualmente (senza `Create`/`Get`/`Update`/`List`) restituisce:

```
WaitUntilStates: refresh callback not set; call Create/Get/Update/List first
```

Usa sempre il wrapper restituito dalla chiamata API, non il builder della richiesta.

### Errori HTTP tipizzati

Le risposte API non-2xx vengono restituite come `*aruba.HTTPError`. Usa `errors.As` per ispezionarle:

```go
vpc, err = arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
if err != nil {
    var httpErr *aruba.HTTPError
    if errors.As(err, &httpErr) {
        log.Printf("API error %d: %s", httpErr.StatusCode, httpErr.Error())
    } else {
        log.Fatalf("Network error: %v", err)
    }
}
```

Vedi [Gestione delle Risposte](./response-handling) per la guida completa alla gestione degli errori.

---

## Esempio Completo

La directory `examples/all-resources/` nel repository contiene un esempio end-to-end eseguibile che dimostra tutte le risorse:

```bash
go run ./examples/all-resources/ -mode=create -clientID=… -clientSecret=…
go run ./examples/all-resources/ -mode=update -clientID=… -clientSecret=… -projectID=…
go run ./examples/all-resources/ -mode=delete -clientID=… -clientSecret=… -projectID=…

# Aggiungi -debug per il logging verboso dell'SDK:
go run ./examples/all-resources/ -mode=create -clientID=… -clientSecret=… -debug
```
