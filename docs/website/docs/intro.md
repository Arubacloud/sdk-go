---
sidebar_position: 1
---

# Quick Start

> **Note**: This SDK is currently in its **Alpha** stage. The API is not yet stable, and breaking changes may be introduced in future releases without prior notice. Please use with caution and be prepared for updates.

Welcome to the official Go SDK for the Aruba Cloud API. This SDK provides a convenient and powerful way for Go developers to interact with the Aruba Cloud API. The primary goal is to simplify the management of cloud resources, allowing you to programmatically create, read, update, and delete resources such as compute instances, virtual private clouds (VPCs), block storage, and more.

## Installation

Add the SDK to your Go project:

```bash
go get github.com/Arubacloud/sdk-go
```

## Getting Started

Getting started with the SDK is straightforward. Import the `aruba` package, create a client with your credentials, and start making API calls using the wrapper builder pattern.

The first and most fundamental resource in the Aruba Cloud is the **Project**. All other resources belong to a project.

Here's how to initialize the SDK client and create your first project:

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
	// Your API credentials
	clientID := "your-client-id"
	clientSecret := "your-client-secret"

	// 1. Initialize the SDK client using default options
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 2. Create the project — build inline and pass to Create
	fmt.Println("Creating a new project...")
	proj, err := arubaClient.FromProject().Create(
		ctx,
		aruba.NewProject().
			Named("my-first-project").
			DescribedAs("A project with the Go SDK").
			Tagged("go-sdk").
			Tagged("quick-start"))
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}

	fmt.Printf("✓ Successfully created project: %s (ID: %s)\n", proj.Name(), proj.ID())
}
```

## Next Steps

- Read the [API Walkthrough](./walkthrough) for an end-to-end guide through the full resource lifecycle
- Learn about [Configuration Options](./options) to customize your SDK client
- Explore [API Resources](./resources) to see what you can manage
- Understand [Response Handling](./response-handling) for robust error handling
- Learn about [Filtering](./filters) to query resources efficiently
- Read about [Multitenancy](./multitenancy) to manage tenant-specific clients
