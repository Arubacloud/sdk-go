---
sidebar_position: 1
---

# Guida Rapida

> **Nota**: Questo SDK è attualmente nella fase **Alpha**. L'API non è ancora stabile e potrebbero essere introdotte modifiche incompatibili nelle versioni future senza preavviso. Si prega di utilizzarlo con cautela e di essere pronti agli aggiornamenti.

Benvenuto nell'SDK ufficiale Go per l'API Aruba Cloud. Questo SDK fornisce un modo comodo e potente per gli sviluppatori Go di interagire con l'API Aruba Cloud. L'obiettivo principale è semplificare la gestione delle risorse cloud, permettendoti di creare, leggere, aggiornare ed eliminare programmaticamente risorse come istanze di calcolo, virtual private cloud (VPC), storage a blocchi e altro ancora.

## Installazione

Aggiungi l'SDK al tuo progetto Go:

```bash
go get github.com/Arubacloud/sdk-go
```

## Iniziare

Iniziare con l'SDK è semplice. Devi importare il pacchetto `aruba`, creare un client con le tue credenziali e poi puoi iniziare a effettuare chiamate API.

La prima e più fondamentale risorsa in Aruba Cloud è il **Progetto**. Tutte le altre risorse appartengono a un progetto.

Ecco come inizializzare il client SDK e creare il tuo primo progetto:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	aruba_types "github.com/Arubacloud/sdk-go/pkg/types"
)

func main() {
	// Le tue credenziali API
	clientID := "your-client-id"
	clientSecret := "your-client-secret"

	// 1. Inizializza il Client SDK utilizzando le opzioni predefinite
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Crea un contesto con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 2. Definisci la richiesta per creare un nuovo Progetto
	projectReq := aruba_types.ProjectRequest{
		Metadata: aruba_types.ResourceMetadataRequest{
			Name: "my-first-project",
			Tags: []string{"go-sdk", "quick-start"},
		},
		Properties: aruba_types.ProjectPropertiesRequest{
			Description: stringPtr("Un progetto creato con l'SDK Go"),
		},
	}

	// 3. Crea il Progetto
	fmt.Println("Creating a new project...")
	createResp, err := arubaClient.FromProject().Create(ctx, projectReq, nil)
	if err != nil {
		// Gestisci errori di rete o errori di validazione lato client
		log.Fatalf("Error creating project: %v", err)
	}

	// 4. Controlla la risposta API
	if !createResp.IsSuccess() {
		// Gestisci errori API (ad esempio, status 4xx o 5xx)
		log.Fatalf("API Error: Failed to create project - Status: %d, Title: %s, Detail: %s",
			createResp.StatusCode,
			stringValue(createResp.Error.Title),
			stringValue(createResp.Error.Detail))
	}

	projectID := *createResp.Data.Metadata.ID
	fmt.Printf("✓ Successfully created project with ID: %s\n", projectID)
}

// Funzioni helper per i puntatori
func stringPtr(s string) *string { return &s }
func stringValue(s *string) string {
	if s == nil { return "" }
	return *s
}
```

## Prossimi Passi

- Scopri le [Opzioni di Configurazione](./options) per personalizzare il tuo client SDK
- Esplora le [Risorse API](./resources) per vedere cosa puoi gestire
- Comprendi la [Gestione delle Risposte](./response-handling) per una gestione robusta degli errori
- Controlla i [Tipi di Dati](./types) per informazioni dettagliate sui tipi
- Scopri il [Filtraggio](./filters) per interrogare le risorse in modo efficiente
