# Guida al Filtraggio

Questa guida spiega come utilizzare i filtri con l'SDK Aruba Cloud Go per interrogare e filtrare le risorse in base a vari criteri.

## Panoramica

L'SDK fornisce un sistema di filtraggio flessibile tramite helper `CallOption`. Passali direttamente a `List` (e ad altre operazioni di lettura) senza costruire struct di parametri intermedi.

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,cpu:gt:2"),
    aruba.WithSort("name:asc"),
    aruba.WithLimit(50),
)
```

Opzioni di chiamata disponibili:

| Opzione | Descrizione |
|---------|-------------|
| `aruba.WithFilter(expr string)` | Espressione di filtro lato server |
| `aruba.WithSort(expr string)` | Espressione di ordinamento |
| `aruba.WithLimit(n int)` | Dimensione della pagina |
| `aruba.WithOffset(n int)` | Offset di paginazione |
| `aruba.WithProjection(expr string)` | Proiezione dei campi |

## Formato Filtro

I filtri seguono questo formato: `campo:operatore:valore`

- **Campo**: Il campo della risorsa su cui filtrare (es. `status`, `name`, `cpu`)
- **Operatore**: L'operatore di confronto (es. `eq`, `gt`, `like`)
- **Valore**: Il valore con cui confrontare

Più filtri vengono combinati usando:
- `,` (virgola) per operazioni **AND**
- `;` (punto e virgola) per operazioni **OR**

## Operatori Supportati

| Operatore | Codice | Descrizione | Esempio |
|-----------|--------|-------------|---------|
| Uguale | `eq` | Corrispondenza esatta | `status:eq:Active` |
| Diverso | `ne` | Non uguale a | `status:ne:Error` |
| Maggiore | `gt` | Maggiore di | `cpu:gt:2` |
| Maggiore/Uguale | `gte` | Maggiore o uguale a | `memory:gte:4096` |
| Minore | `lt` | Minore di | `disk:lt:100` |
| Minore/Uguale | `lte` | Minore o uguale a | `cpu:lte:8` |
| In | `in` | Valore in lista | `region:in:us-east,us-west` |
| Non In | `nin` | Valore non in lista | `status:nin:Error,Failed` |
| Contiene | `like` | Corrispondenza sottostringa | `name:like:prod` |
| Inizia Con | `sw` | Corrispondenza prefisso | `name:sw:web-` |
| Termina Con | `ew` | Corrispondenza suffisso | `name:ew:-prod` |

## Filtri Semplici

### Condizione Singola

```go
// Elenca i cloud server attivi
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active"),
)
```

### Condizioni AND Multiple

```go
// Server attivi con almeno 2 vCPU e 4 GB di RAM
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,cpu:gt:2,memory:gte:4096"),
)
```

Espressione risultante: `status:eq:Active,cpu:gt:2,memory:gte:4096`

### Condizioni OR

```go
// Server che sono Attivi O in Avvio
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active;status:eq:Starting"),
)
```

Espressione risultante: `status:eq:Active;status:eq:Starting`

### Filtri Complessi (AND + OR)

```go
// (environment=production AND memory>=4096) OR (tier=premium AND region IN [us-east-1,us-west-2])
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("environment:eq:production,memory:gte:4096;tier:eq:premium,region:in:us-east-1,us-west-2"),
)
```

## Esempi Pratici

### Filtra Cloud Server Attivi

```go
// Elenca server attivi con almeno 4 GB di RAM, dimensione pagina 50
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,memory:gte:4096"),
    aruba.WithLimit(50),
)
if err != nil {
    log.Fatalf("List failed: %v", err)
}
fmt.Printf("Found %d servers\n", servers.Total())
for _, s := range servers.Items() {
    fmt.Println("-", s.Name())
}
```

### Filtra VPC per Regione

```go
// Elenca VPC in data center specifici
vpcs, err := arubaClient.FromNetwork().VPCs().List(ctx, proj,
    aruba.WithFilter("location:in:ITBG-Bergamo,ITMI-Milan"),
)
```

### Filtra per Pattern Nome

```go
// Tutti i cloud server il cui nome inizia con "web-"
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("name:sw:web-"),
)
```

### Logica di Business Complessa

```go
// Server di produzione con risorse elevate O server di sviluppo in regioni specifiche
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter(
        "environment:eq:production,cpu:gte:8,memory:gte:16384" +
        ";environment:eq:development,region:in:ITBG-Bergamo,ITMI-Milan",
    ),
)
```

## Ordinamento

```go
// Ordina per nome crescente
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithSort("name:asc"),
)

// Ordina per data di creazione decrescente
servers, err = arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithSort("createdAt:desc"),
)
```

## Paginazione

```go
const pageSize = 25

// Prima pagina
page1, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithLimit(pageSize),
    aruba.WithOffset(0),
)

// Seconda pagina
page2, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithLimit(pageSize),
    aruba.WithOffset(pageSize),
)

fmt.Printf("Total resources: %d\n", page1.Total())
```

## Best Practice

### Usa Filtri Specifici

Sii il più specifico possibile per ridurre la quantità di dati trasferiti:

```go
// Bene: filtro specifico — solo le risorse necessarie
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,region:eq:ITBG-Bergamo"),
)

// Meno efficiente: nessun filtro — recupera tutto
servers, err = arubaClient.FromCompute().CloudServers().List(ctx, proj)
```

### Combina Filtro e Paginazione

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("environment:eq:production"),
    aruba.WithLimit(50),
    aruba.WithOffset(0),
)
```

### Valida le Stringhe di Filtro

Stampa la stringa di filtro per debuggare risultati inaspettati:

```go
filter := "status:eq:Active,cpu:gt:4"
fmt.Println("Filter:", filter)
// Output: Filter: status:eq:Active,cpu:gt:4

servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter(filter),
)
```

## Risoluzione dei Problemi

### Filtro Non Funzionante

**Problema**: Il filtro non restituisce i risultati attesi

**Soluzioni**:
1. Controlla che i nomi dei campi corrispondano esattamente allo schema API (sensibile alle maiuscole)
2. Verifica che l'operatore sia appropriato per il tipo di campo
3. Assicurati che i valori siano nel formato corretto (es. `Active` non `active`)
4. Stampa la stringa di filtro per il debug

### Risultati Vuoti

**Problema**: Il filtro non restituisce risultati

**Soluzioni**:
1. Verifica che la logica del filtro sia corretta
2. Prova filtri più semplici per isolare il problema
3. Controlla se il campo supporta l'operatore utilizzato
4. Elenca senza filtri per confermare che le risorse esistano

### Problemi con Filtri Complessi

**Problema**: La logica AND/OR complessa non funziona come previsto

**Soluzioni**:
1. Scomponi i filtri complessi in parti più semplici
2. Testa ogni condizione separatamente
3. Ricorda: virgole (`,`) = AND, punto e virgola (`;`) = OR
4. Usa le parentesi mentalmente per capire il raggruppamento

```
// Questo: status:eq:Active,cpu:gt:2;memory:gte:4096
// Significa: (status = Active AND cpu > 2) OR (memory >= 4096)

// Non: status = Active AND (cpu > 2 OR memory >= 4096)
```
