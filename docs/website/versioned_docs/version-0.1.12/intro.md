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

Getting started with the SDK is straightforward. You need to import the `aruba` package, create a client with your credentials, and then you can start making API calls.

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
	aruba_types "github.com/Arubacloud/sdk-go/pkg/types"
)

func main() {
	// Your API credentials
	clientID := "your-client-id"
	clientSecret := "your-client-secret"

	// 1. Initialize the SDK Client using default options
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 2. Define the request to create a new Project
	projectReq := aruba_types.ProjectRequest{
		Metadata: aruba_types.ResourceMetadataRequest{
			Name: "my-first-project",
			Tags: []string{"go-sdk", "quick-start"},
		},
		Properties: aruba_types.ProjectPropertiesRequest{
			Description: stringPtr("A project created with the Go SDK"),
		},
	}

	// 3. Create the Project
	fmt.Println("Creating a new project...")
	createResp, err := arubaClient.FromProject().Create(ctx, projectReq, nil)
	if err != nil {
		// Handle network errors or client-side validation errors
		log.Fatalf("Error creating project: %v", err)
	}

	// 4. Check the API response
	if !createResp.IsSuccess() {
		// Handle API errors (e.g., 4xx or 5xx statuses)
		log.Fatalf("API Error: Failed to create project - Status: %d, Title: %s, Detail: %s",
			createResp.StatusCode,
			stringValue(createResp.Error.Title),
			stringValue(createResp.Error.Detail))
	}

	projectID := *createResp.Data.Metadata.ID
	fmt.Printf("✓ Successfully created project with ID: %s\n", projectID)
}

// Helper functions for pointers
func stringPtr(s string) *string { return &s }
func stringValue(s *string) string {
	if s == nil { return "" }
	return *s
}
```

## Next Steps

- Learn about [Configuration Options](./options) to customize your SDK client
- Explore [API Resources](./resources) to see what you can manage
- Understand [Response Handling](./response-handling) for robust error handling
- Check out [Data Types](./types) for detailed type information
- Learn about [Filtering](./filters) to query resources efficiently

