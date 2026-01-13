# Guida alla Gestione delle Risposte

## Panoramica

L'SDK utilizza un tipo generico `Response[T]` che gestisce correttamente sia le risposte di successo che quelle di errore dall'API. L'analisi delle risposte è centralizzata tramite la funzione `ParseResponseBody[T]` per coerenza.

## Struttura della Risposta

```go
type Response[T any] struct {
    Data         *T              // Popolato per risposte 2xx
    Error        *ErrorResponse  // Popolato per risposte 4xx/5xx
    HTTPResponse *http.Response  // La risposta HTTP sottostante
    StatusCode   int             // Il codice di stato HTTP
    Headers      http.Header     // Intestazioni della risposta
    RawBody      []byte          // Corpo della risposta raw (sempre disponibile)
}
```

## Analisi Centralizzata delle Risposte

Tutti i metodi del servizio utilizzano la funzione `ParseResponseBody[T]` per gestire l'analisi delle risposte:

```go
func ParseResponseBody[T any](httpResp *http.Response) (*Response[T], error) {
    // Legge il corpo della risposta
    // Crea il wrapper Response[T]
    // Analizza in Data per risposte 2xx
    // Analizza in Error per risposte 4xx/5xx
    return response, nil
}
```

### Implementazione nei Servizi

I metodi del servizio chiamano semplicemente `ParseResponseBody` dopo la richiesta HTTP:

```go
func (s *VPCService) GetVPC(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VpcResponse], error) {
    // ... prepara la richiesta ...
    
    httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
    if err != nil {
        return nil, err
    }
    defer httpResp.Body.Close()

    return schema.ParseResponseBody[schema.VpcResponse](httpResp)
}
```

**Vantaggi:**
- ✅ Elimina la duplicazione del codice in tutti i servizi
- ✅ Garantisce una gestione degli errori coerente
- ✅ Semplifica le implementazioni dei servizi
- ✅ Rende gli aggiornamenti più facili da mantenere

## Struttura della Risposta di Errore

```go
type ErrorResponse struct {
    Type       *string                 // Riferimento URI per il tipo di problema
    Title      *string                 // Riepilogo breve e leggibile
    Status     *int32                  // Codice di stato HTTP
    Detail     *string                 // Spiegazione leggibile
    Instance   *string                 // URI per questa occorrenza specifica
    Extensions map[string]interface{}  // Proprietà dinamiche aggiuntive
}
```

## Pattern di Gestione delle Risposte

### 1. Pattern Base

```go
resp, err := api.CreateResource(ctx, projectID, request, nil)
if err != nil {
    // Errore di rete, timeout del contesto o errore SDK
    log.Fatalf("Request failed: %v", err)
}

if resp.IsSuccess() {
    // 2xx - Risposta di successo
    fmt.Printf("Created: %s\n", *resp.Data.Metadata.Name)
} else if resp.IsError() && resp.Error != nil {
    // 4xx/5xx - Risposta di errore API
    log.Printf("API Error: %s - %s", 
        stringValue(resp.Error.Title), 
        stringValue(resp.Error.Detail))
}
```

### 2. Gestione Completa degli Errori

```go
resp, err := api.GetResource(ctx, projectID, resourceID, nil)
if err != nil {
    return fmt.Errorf("request failed: %w", err)
}

switch {
case resp.IsSuccess():
    // Gestisci il successo - resp.Data è popolato
    resource := resp.Data
    fmt.Printf("Resource: %s (Status: %s)\n", 
        *resource.Metadata.Name, 
        *resource.Status.State)
    return nil

case resp.StatusCode == 404:
    // Gestisci non trovato
    return fmt.Errorf("resource not found")

case resp.StatusCode == 400:
    // Gestisci errori di validazione
    if resp.Error != nil {
        return fmt.Errorf("validation error: %s", stringValue(resp.Error.Detail))
    }
    return fmt.Errorf("bad request: %s", string(resp.RawBody))

case resp.IsError():
    // Gestisci altri errori
    if resp.Error != nil {
        return fmt.Errorf("API error (%d): %s - %s", 
            resp.StatusCode,
            stringValue(resp.Error.Title),
            stringValue(resp.Error.Detail))
    }
    return fmt.Errorf("unexpected error (%d): %s", resp.StatusCode, string(resp.RawBody))

default:
    // Codice di stato inaspettato
    return fmt.Errorf("unexpected status %d", resp.StatusCode)
}
```

### 3. Accesso ai Dettagli dell'Errore

```go
if resp.IsError() && resp.Error != nil {
    // Campi standard
    log.Printf("Error Type: %s", stringValue(resp.Error.Type))
    log.Printf("Error Title: %s", stringValue(resp.Error.Title))
    log.Printf("Error Detail: %s", stringValue(resp.Error.Detail))
    log.Printf("Status Code: %d", int32Value(resp.Error.Status))
    
    // Accedi alle estensioni personalizzate (ad esempio, errori di validazione)
    if errors, ok := resp.Error.Extensions["errors"].([]interface{}); ok {
        for _, e := range errors {
            if errMap, ok := e.(map[string]interface{}); ok {
                log.Printf("  Field: %s, Message: %s", 
                    errMap["field"], 
                    errMap["message"])
            }
        }
    }
}
```

### 4. Accesso al Corpo Raw

```go
// Sempre disponibile per il debug
log.Printf("Raw response: %s", string(resp.RawBody))

// Utile per registrare risposte complete durante lo sviluppo
if !resp.IsSuccess() {
    log.Printf("Request failed with status %d: %s", 
        resp.StatusCode, 
        string(resp.RawBody))
}
```

## Funzioni Helper

```go
// Helper per dereferenziare i puntatori in modo sicuro
func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func int32Value(i *int32) int32 {
    if i == nil {
        return 0
    }
    return *i
}

func boolValue(b *bool) bool {
    if b == nil {
        return false
    }
    return *b
}
```

## Scenari di Errore Comuni

### 400 Bad Request - Errori di Validazione

```go
resp, err := api.CreateResource(ctx, projectID, invalidRequest, nil)
if err != nil {
    log.Fatalf("Request failed: %v", err)
}

if resp.StatusCode == 400 && resp.Error != nil {
    fmt.Printf("Validation failed: %s\n", stringValue(resp.Error.Title))
    
    // Controlla errori a livello di campo in Extensions
    if errors, ok := resp.Error.Extensions["errors"].([]interface{}); ok {
        for _, e := range errors {
            if errMap, ok := e.(map[string]interface{}); ok {
                fmt.Printf("  - %s: %s\n", errMap["field"], errMap["message"])
            }
        }
    }
}
```

### 404 Not Found

```go
resp, err := api.GetResource(ctx, projectID, resourceID, nil)
if err != nil {
    return err
}

if resp.StatusCode == 404 {
    return fmt.Errorf("resource %s not found", resourceID)
}

if !resp.IsSuccess() {
    return fmt.Errorf("unexpected error: %d", resp.StatusCode)
}

// Usa resp.Data
resource := resp.Data
```

### 500 Internal Server Error

```go
if resp.StatusCode >= 500 {
    // Errore del server - potrebbe voler riprovare
    if resp.Error != nil {
        log.Printf("Server error: %s", stringValue(resp.Error.Detail))
    }
    
    // Registra trace ID per il supporto
    if traceID, ok := resp.Error.Extensions["traceId"].(string); ok {
        log.Printf("Trace ID: %s", traceID)
    }
}
```

## Best Practices

1. **Controlla sempre gli errori di rete per primi** (`err != nil`)
2. **Usa `IsSuccess()` per controllare le risposte 2xx** prima di accedere a `Data`
3. **Usa `IsError()` per controllare le risposte 4xx/5xx** prima di accedere a `Error`
4. **Controlla se il campo `Error` è non-nil** prima di dereferenziare
5. **Usa funzioni helper** per dereferenziare in modo sicuro i campi puntatore
6. **Mantieni `RawBody` disponibile** per debug e logging
7. **Registra trace ID** dalle risposte di errore per le richieste di supporto

## Test della Gestione delle Risposte

```go
func TestResourceCreation(t *testing.T) {
    resp, err := api.CreateResource(ctx, projectID, request, nil)
    
    // Controlla nessun errore di rete
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    
    // Controlla risposta di successo
    if !resp.IsSuccess() {
        if resp.Error != nil {
            t.Fatalf("API error: %s - %s", 
                stringValue(resp.Error.Title),
                stringValue(resp.Error.Detail))
        }
        t.Fatalf("Unexpected status: %d, body: %s", 
            resp.StatusCode, 
            string(resp.RawBody))
    }
    
    // Valida i dati della risposta
    if resp.Data == nil {
        t.Fatal("Expected data to be populated")
    }
    
    if resp.Data.Metadata.Name == nil {
        t.Fatal("Expected resource name")
    }
}
```
