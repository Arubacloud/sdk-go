# Aruba Cloud SDK for Go

Welcome to the official Go SDK for the Aruba Cloud API. This SDK provides a convenient and powerful way for Go developers to interact with the Aruba Cloud API. The primary goal is to simplify the management of cloud resources, allowing you to programmatically create, read, update, and delete resources such as compute instances, virtual private clouds (VPCs), block storage, and more.

## 1. Quick Start

Getting started with the SDK is straightforward. You need to import the `aruba` package, create a client with your credentials, and then you can start making API calls.

The first and most fundamental resource in the Aruba Cloud is the **Project**. All other resources belong to a project.

Here’s how to initialize the SDK client and create your first project:

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

## 2. Usage Details

### 2.1. Config Options

The SDK client is configured using an `Options` object. The easiest way to get started is with `aruba.DefaultOptions`, which sets up a production-ready client.

```go
// Creates a client with default settings for production
options := aruba.DefaultOptions("your-client-id", "your-client-secret")

// You can further customize the options using a fluent API
options.WithNativeLogger() // Enable built-in logging

// Create the client from the options
arubaClient, err := aruba.NewClient(options)
```

Key configuration areas include:

-   **Authentication**: By default, the SDK uses the provided `clientID` and `clientSecret` to automatically manage OAuth2 tokens. For advanced scenarios, you can use `WithToken()` for static tokens or configure a Vault repository.
-   **Logging**: Logging is disabled by default. You can enable the built-in logger with `WithNativeLogger()` or provide your own `logger.Logger` implementation with `WithCustomLogger()`.
-   **HTTP Client**: The SDK uses `http.DefaultClient`. You can supply your own `*http.Client` with custom settings (like timeouts) using `WithCustomHTTPClient()`.

For a comprehensive guide on all available SDK configuration options, including detailed explanations of mutual exclusions and side effects, please refer to the [Options Documentation](./doc/OPTIONS.md).

### 2.2. Performing Calls, Setting Filters, and Handling Responses

#### Performing API Calls

API calls are organized into logical groups (e.g., `Compute`, `Network`, `Storage`). You can access these groups from the main `arubaClient` object.

The general pattern is: `arubaClient.From<Group>().<Resource>().<Action>()`

```go
// Example: List all Cloud Servers in a project
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, nil)

// Example: Get a specific VPC
vpc, err := arubaClient.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)

// Example: Delete a Block Storage volume
deleteResp, err := arubaClient.FromStorage().Volumes().Delete(ctx, projectID, volumeID, nil)
```

#### Setting Filters

When listing resources, you can use filters to narrow down the results. The SDK provides a `FilterBuilder` to construct complex filter expressions programmatically.

Filters follow the format `field:operator:value`. Multiple filters are combined with `,` (for AND) or `;` (for OR).

**Example: Find all 'running' cloud servers with more than 2 CPU cores.**

```go
import "github.com/Arubacloud/sdk-go/pkg/types"

// Construct the filter string
filterStr := types.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 2).
    Build() // Builds the string: "status:eq:running,cpu:gt:2"

// Create request parameters and assign the filter
params := &types.RequestParameters{
    Filter: &filterStr,
}

// Perform the call
resp, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, params)
```

Supported operators include `Equal`, `NotEqual`, `GreaterThan`, `LessThan`, `In`, `Contains`, `StartsWith`, and `EndsWith`.

For a comprehensive guide on advanced filtering capabilities, including programmatic filter building and practical examples, please refer to the [Filtering Guide](./doc/FILTERS.md).

#### Handling Responses

All API calls return a `*types.Response[T]` and an `error`.

1.  **Always check the `error` first.** This indicates a network problem, a context timeout, or a client-side issue.
2.  **If `error` is `nil`, check the response status.** The `Response` object contains helper methods like `IsSuccess()` (for 2xx statuses) and `IsError()` (for 4xx/5xx statuses).

The generic `Response` struct looks like this:

```go
type Response[T any] struct {
    Data         *T              // Populated on success (2xx)
    Error        *ErrorResponse  // Populated on failure (4xx/5xx)
    HTTPResponse *http.Response  // The raw HTTP response
    StatusCode   int             // The HTTP status code
    RawBody      []byte          // The raw response body
}
```

**Standard response handling pattern:**

```go
resp, err := arubaClient.FromNetwork().VPCs().Get(ctx, projectID, "non-existent-vpc", nil)

// 1. Check for network/client errors
if err != nil {
    log.Fatalf("Request failed: %v", err)
}

// 2. Check for success (2xx)
if resp.IsSuccess() {
    // Safely access the response data
    fmt.Printf("VPC Name: %s\n", *resp.Data.Metadata.Name)
    return
}

// 3. Handle API errors (4xx/5xx)
if resp.IsError() && resp.Error != nil {
    // Log structured error details
    log.Printf("API Error Occurred (Status %d):", resp.StatusCode)
    log.Printf("  Title: %s", stringValue(resp.Error.Title))
    log.Printf("  Detail: %s", stringValue(resp.Error.Detail))

    // You can also inspect specific status codes
    if resp.StatusCode == 404 {
        log.Println("  Diagnosis: The resource was not found.")
    }
}
```

For a comprehensive guide on response handling, including detailed error scenarios and best practices, please refer to the [Response Handling Guide](./doc/RESPONSE_HANDLING.md).

### 2.3. SDK Client API Reference

#### Groups

The SDK is divided into the following service groups, accessible from the `arubaClient`:

-   `FromAudit()`: Access audit events.
-   `FromCompute()`: Manage virtual machines and SSH keys.
-   `FromContainer()`: Manage Kubernetes (KaaS) and Container Registry services.
-   `FromDatabase()`: Manage Database-as-a-Service (DBaaS) instances.
-   `FromMetric()`: Access metrics and alerts.
-   `FromNetwork()`: Manage VPCs, subnets, security groups, elastic IPs, etc.
-   `FromProject()`: Manage projects.
-   `FromSchedule()`: Manage scheduled jobs.
-   `FromSecurity()`: Manage KMS keys.
-   `FromStorage()`: Manage block storage volumes, snapshots, backups, and restores.

#### Resources

Each group provides clients for specific resources. For example:

-   `arubaClient.FromCompute().CloudServers()`
-   `arubaClient.FromCompute().KeyPairs()`
-   `arubaClient.FromNetwork().VPCs()`
-   `arubaClient.FromNetwork().Subnets()`
-   `arubaClient.FromStorage().Volumes()`

Each resource client provides methods for CRUD operations (`Create`, `List`, `Get`, `Update`, `Delete`), where applicable.

For a comprehensive guide on all API groups, resource clients, and their available operations, please refer to the [API Groups and Resources Documentation](./doc/RESOURCES.md).

#### Data Types

All request bodies, response objects, and data models are defined in the `pkg/types` package. They are named to correspond with the resources they represent.

-   **Requests**: Typically end in `Request` (e.g., `types.VPCRequest`, `types.CloudServerRequest`).
-   **Responses**: Typically end in `Response` or `List` (e.g., `types.VPCResponse`, `types.VPCList`).
-   **Shared Structures**: Common structures like `ResourceMetadataRequest`, `ResourceMetadataResponse`, and `ResourceStatus` are used across many types.

For a comprehensive guide on all SDK data types, including their structure and usage in requests and responses, please refer to the [SDK Data Types Documentation](./doc/TYPES.md).