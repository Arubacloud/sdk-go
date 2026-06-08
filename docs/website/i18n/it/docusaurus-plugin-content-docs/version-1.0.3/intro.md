---
sidebar_position: 1
---

# Guida Rapida

Benvenuto nell'SDK ufficiale Go per l'API Aruba Cloud. Questo SDK fornisce un modo comodo e potente per gli sviluppatori Go di interagire con l'API Aruba Cloud. L'obiettivo principale è semplificare la gestione delle risorse cloud, permettendoti di creare, leggere, aggiornare ed eliminare programmaticamente risorse come istanze di calcolo, virtual private cloud (VPC), storage a blocchi e altro ancora.

## Installazione

Aggiungi l'SDK al tuo progetto Go:

```bash
go get github.com/Arubacloud/sdk-go@latest
```

## Iniziare

Iniziare con l'SDK è semplice. Importa il pacchetto `aruba`, crea un client con le tue credenziali e inizia a effettuare chiamate API usando il pattern builder fluente.

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
)

func main() {
	// Le tue credenziali API
	clientID := "your-client-id"
	clientSecret := "your-client-secret"

	// 1. Inizializza il client SDK utilizzando le opzioni predefinite
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Crea un contesto con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 2. Crea il Progetto — costruisci inline e passa a Create
	fmt.Println("Creating a new project...")
	proj, err := arubaClient.FromProject().Create(
		ctx,
		aruba.NewProject().
			Named("my-first-project").
			Tagged("go-sdk", "quick-start").
			DescribedAs("Un progetto creato con l'SDK Go"))
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}

	fmt.Printf("✓ Progetto creato con successo: %s (ID: %s)\n", proj.Name(), proj.ID())
}
```

## Prossimi Passi

- Leggi la [Guida al Walkthrough API](./walkthrough) per un percorso completo attraverso il ciclo di vita delle risorse
- Scopri le [Opzioni di Configurazione](./options) per personalizzare il tuo client SDK
- Esplora le [Risorse API](./resources) per vedere cosa puoi gestire
- Comprendi la [Gestione delle Risposte](./response-handling) per una gestione robusta degli errori
- Scopri il [Filtraggio](./filters) per interrogare le risorse in modo efficiente
- Leggi [Multitenancy](./multitenancy) per gestire client specifici per tenant
- Vedi [Utilizzo a Basso Livello](./working-at-low-level) per funzionalità avanzate (`pkg/types`, `pkg/async`)
