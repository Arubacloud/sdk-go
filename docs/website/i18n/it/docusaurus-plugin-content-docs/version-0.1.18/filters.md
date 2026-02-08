# Guida al Filtraggio

Questa guida spiega come utilizzare i filtri con l'SDK Aruba Cloud Go per interrogare e filtrare le risorse in base a vari criteri.

## Panoramica

L'SDK fornisce un sistema di filtraggio potente e flessibile che segue la specifica dei filtri dell'API Aruba Cloud. I filtri ti consentono di:

- Interrogare le risorse in base ai valori dei campi
- Utilizzare operatori di confronto (uguale, maggiore di, minore di, ecc.)
- Combinare più condizioni con operatori logici (AND/OR)
- Eseguire corrispondenze di pattern (contiene, inizia con, termina con)
- Filtrare per più valori (IN, NOT IN)

## Formato Filtro

I filtri seguono questo formato: `field:operator:value`

- **Field**: Il campo della risorsa su cui filtrare (ad esempio, `status`, `name`, `cpu`)
- **Operator**: L'operatore di confronto (ad esempio, `eq`, `gt`, `like`)
- **Value**: Il valore con cui confrontare

I filtri multipli sono combinati utilizzando:
- `,` (virgola) per operazioni **AND**
- `;` (punto e virgola) per operazioni **OR**

## Operatori Supportati

| Operatore | Codice   | Descrizione              | Esempio                    |
|----------|--------|--------------------------|----------------------------|
| Uguale    | `eq`   | Corrispondenza esatta              | `status:eq:running`        |
| Diverso| `ne`   | Diverso da             | `status:ne:stopped`        |
| Maggiore  | `gt`   | Maggiore di             | `cpu:gt:2`                 |
| Maggiore/Uguale | `gte` | Maggiore o uguale | `memory:gte:4096`         |
| Minore     | `lt`   | Minore di                | `disk:lt:100`              |
| Minore/Uguale | `lte` | Minore o uguale      | `cpu:lte:8`                |
| In       | `in`   | Valore nell'elenco            | `region:in:us-east,us-west`|
| Non In   | `nin`  | Valore non nell'elenco        | `status:nin:failed,error`  |
| Contiene | `like` | Corrispondenza sottostringa          | `name:like:prod`           |
| Inizia Con | `sw` | Corrispondenza prefisso            | `name:sw:web-`             |
| Termina Con | `ew`  | Corrispondenza suffisso             | `name:ew:-prod`            |

## Utilizzo di FilterBuilder

L'SDK fornisce un `FilterBuilder` per costruire espressioni di filtro complesse in modo programmatico.

### Filtri Semplici

#### Condizione Singola

```go
import "github.com/Arubacloud/sdk-go/pkg/types"

// Filtra per stato
filter := types.NewFilterBuilder().
    Equal("status", "running").
    Build()

params := &types.RequestParameters{
    Filter: &filter,
}

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

Risultato: `status:eq:running`

#### Condizioni AND Multiple

```go
// Filtra per stato AND cpu AND memoria
filter := types.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 2).
    GreaterThanOrEqual("memory", 4096).
    Build()
```

Risultato: `status:eq:running,cpu:gt:2,memory:gte:4096`

### Condizioni OR

```go
// Filtra per stato = running OR stato = starting
filter := types.NewFilterBuilder().
    Equal("status", "running").
    Or().
    Equal("status", "starting").
    Build()
```

Risultato: `status:eq:running;status:eq:starting`

### Filtri Complessi (AND + OR)

```go
// (environment = production AND memory >= 4096) OR (tier = premium AND region IN [us-east-1, us-west-2])
filter := types.NewFilterBuilder().
    Equal("environment", "production").
    GreaterThanOrEqual("memory", 4096).
    Or().
    Equal("tier", "premium").
    In("region", "us-east-1", "us-west-2").
    Build()
```

Risultato: `environment:eq:production,memory:gte:4096;tier:eq:premium,region:in:us-east-1,us-west-2`

## Metodi Filtro

### Metodi di Confronto

```go
fb := types.NewFilterBuilder()

// Uguaglianza
fb.Equal("field", "value")           // field = value
fb.NotEqual("field", "value")        // field != value

// Confronti numerici
fb.GreaterThan("field", 100)         // field > 100
fb.GreaterThanOrEqual("field", 100)  // field >= 100
fb.LessThan("field", 100)            // field < 100
fb.LessThanOrEqual("field", 100)     // field <= 100

// Operazioni su elenchi
fb.In("field", "val1", "val2", "val3")     // field IN (val1, val2, val3)
fb.NotIn("field", "val1", "val2")          // field NOT IN (val1, val2)

// Corrispondenza pattern stringa
fb.Contains("field", "substring")    // field LIKE %substring%
fb.StartsWith("field", "prefix")     // field LIKE prefix%
fb.EndsWith("field", "suffix")       // field LIKE %suffix
```

### Operatori Logici

```go
fb := types.NewFilterBuilder()

// Il predefinito è AND
fb.Equal("field1", "value1").
   Equal("field2", "value2")  // field1 = value1 AND field2 = value2

// OR esplicito
fb.Equal("field1", "value1").
   Or().
   Equal("field2", "value2")  // field1 = value1 OR field2 = value2

// Mescola AND e OR
fb.Equal("field1", "value1").
   Equal("field2", "value2").  // Gruppo 1: AND
   Or().
   Equal("field3", "value3")   // Gruppo 2: OR
```

## Esempi Pratici

### Filtra Cloud Server Attivi

```go
// Elenca tutti i cloud server in esecuzione con almeno 4GB di RAM
filter := types.NewFilterBuilder().
    Equal("status", "running").
    GreaterThanOrEqual("memory", 4096).
    Build()

params := &types.RequestParameters{
    Filter: &filter,
    Limit:  types.Int32Ptr(50),
}

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

### Filtra per Più Regioni

```go
// Elenca risorse nelle regioni US East o US West
filter := types.NewFilterBuilder().
    In("region", "us-east-1", "us-east-2", "us-west-1", "us-west-2").
    Build()

params := &types.RequestParameters{
    Filter: &filter,
}

resp, err := arubaClient.FromNetwork().VPCs().List(ctx, projectID, params)
```

### Filtra per Pattern Nome

```go
// Elenca tutti i server web di produzione
filter := types.NewFilterBuilder().
    StartsWith("name", "web-").
    Contains("environment", "prod").
    Build()

params := &types.RequestParameters{
    Filter: &filter,
}

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

### Logica di Business Complessa

```go
// Trova server che sono:
// - Server di produzione con risorse elevate (cpu >= 8 AND memory >= 16GB)
// - OR Server di sviluppo in regioni specifiche
filter := types.NewFilterBuilder().
    Equal("environment", "production").
    GreaterThanOrEqual("cpu", 8).
    GreaterThanOrEqual("memory", 16384).
    Or().
    Equal("environment", "development").
    In("region", "us-east-1", "eu-west-1").
    Build()

params := &types.RequestParameters{
    Filter: &filter,
}

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

## Best Practices

### 1. Usa Filtri Specifici

Sii il più specifico possibile per ridurre la quantità di dati trasferiti:

```go
// Buono: Filtro specifico
filter := types.NewFilterBuilder().
    Equal("status", "running").
    Equal("region", "us-east-1").
    Build()

// Meno efficiente: Nessun filtro, elabora tutti i risultati
resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, nil)
```

### 2. Combina con Paginazione

Usa i filtri con la paginazione per set di risultati grandi:

```go
filter := types.NewFilterBuilder().
    Equal("environment", "production").
    Build()

params := &types.RequestParameters{
    Filter: &filter,
    Limit:  types.Int32Ptr(50),
}

resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

### 3. Valida la Logica del Filtro

Testa le tue espressioni di filtro per assicurarti che producano i risultati attesi:

```go
fb := types.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 4)

filterStr := fb.Build()
fmt.Println("Filter:", filterStr)
// Output: Filter: status:eq:running,cpu:gt:4
```

### 4. Usa Valori Type-Safe

Usa i tipi corretti per i valori del filtro:

```go
// Buono: Tipi corretti
fb.Equal("cpu", 4)           // int
fb.Equal("status", "running") // string
fb.GreaterThan("memory", 4096) // int

// Evita: Rappresentazione stringa di numeri quando sono attesi numeri
fb.Equal("cpu", "4") // Potrebbe non funzionare come previsto
```

## Risoluzione Problemi

### Filtro Non Funziona

**Problema**: Il filtro non restituisce i risultati attesi

**Soluzioni**:
1. Controlla che i nomi dei campi corrispondano esattamente allo schema API (case-sensitive)
2. Verifica che l'operatore sia appropriato per il tipo di campo
3. Assicurati che i valori siano del tipo corretto (string, int, bool)
4. Stampa la stringa del filtro per il debug: `fmt.Println(fb.Build())`

### Risultati Vuoti

**Problema**: Il filtro non restituisce risultati

**Soluzioni**:
1. Verifica che la logica del filtro sia corretta
2. Prova filtri più semplici per isolare il problema
3. Controlla se il campo supporta l'operatore utilizzato
4. Testa senza filtri per confermare che le risorse esistano

### Problemi con Filtri Complessi

**Problema**: La logica AND/OR complessa non funziona come previsto

**Soluzioni**:
1. Suddividi i filtri complessi in parti più semplici
2. Testa ogni condizione separatamente
3. Ricorda: virgole (`,`) = AND, punto e virgola (`;`) = OR
4. Usa parentesi nella tua mente per capire il raggruppamento

```go
// Questo: status:eq:running,cpu:gt:2;memory:gte:4096
// Significa: (status = running AND cpu > 2) OR (memory >= 4096)

// Non: status = running AND (cpu > 2 OR memory >= 4096)
```
