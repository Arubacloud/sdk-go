# Aruba Cloud SDK for Go

> **Note**: This SDK is currently in its **Alpha** stage. The API is not yet stable, and breaking changes may be introduced in future releases without prior notice. Please use with caution and be prepared for updates.

### Table of Contents

- [1. Quick Start](#1-quick-start)
- [2. Usage Details](#2-usage-details)
  - [2.1. Config Options](#21-config-options)
  - [2.2. Performing Calls, Setting Filters, and Handling Responses](#22-performing-calls-setting-filters-and-handling-responses)
  - [2.3. SDK Client API Reference](#23-sdk-client-api-reference)
  - [2.4. Handling Asynchronous Operations](#24-handling-asynchronous-operations)
  - [2.5. Data Types](#25-data-types)

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

### 2.4. Handling Asynchronous Operations

Many operations in a cloud environment are asynchronous. For instance, when you create a new Cloud Server, the API call might return successfully, but the server itself takes a few minutes to be provisioned and ready. The SDK provides helpers in the `pkg/async` package to manage these scenarios by polling for a desired state.

The primary tool for this is the `async.WaitFor` function. It repeatedly executes an API call until a specific condition is met or a timeout is reached. A simpler `async.DefaultWaitFor` is also available for common use cases.

**Example: Waiting for a Cloud Server to be 'running'**

Let's say you've just created a Cloud Server and want to wait until it is fully provisioned and has a status of `"running"`.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/Arubacloud/sdk-go/pkg/async"
	aruba_types "github.com/Arubacloud/sdk-go/pkg/types"
)

// Assume arubaClient is initialized and projectID and serverID are available.
var arubaClient *aruba.Client
var projectID = "your-project-id"
var serverID = "your-server-id"

func main() {
	// Initialize client, create server, etc.
	// ...

	fmt.Printf("Waiting for server %s to become ready...\n", serverID)

	// Use DefaultWaitFor for simple polling with default settings.
	// It will retry every 10s for up to 10 times with a total timeout of 60s.
	waitFuture := async.DefaultWaitFor(
		context.Background(),
		// This is the function that will be called repeatedly.
		// It should fetch the resource's current state.
		func(ctx context.Context) (*aruba_types.Response[aruba_types.CloudServerResponse], error) {
			return arubaClient.FromCompute().CloudServers().Get(ctx, projectID, serverID, nil)
		},
		// This function checks if the desired state has been reached.
		// It returns 'true' to stop waiting, or 'false' to continue.
		func(resp *aruba_types.Response[aruba_types.CloudServerResponse]) (bool, error) {
			if resp.IsSuccess() {
				// Replace "running" with the actual desired status of the resource.
				isReady := resp.Data.Properties.Status == "running"
				if isReady {
					fmt.Println("✓ Server is now running.")
				}
				return isReady, nil
			}
			// Continue waiting if the resource isn't ready or if there's a temporary API issue.
			return false, nil
		},
	)

	// Await blocks until the wait operation is complete (or fails).
	finalResp, err := waitFuture.Await(context.Background())
	if err != nil {
		log.Fatalf("Error waiting for server: %v", err)
	}

	// Now you can safely work with the resource, knowing it is in the desired state.
	fmt.Printf("Server %s is ready. Status: %s\n", *finalResp.Data.Metadata.ID, finalResp.Data.Properties.Status)
}
```

For more control over the polling behavior, you can use the `async.WaitFor` function directly. This allows you to specify custom retry counts, delays, and a total timeout.

```go
// Example with custom parameters
waitFutureCustom := async.WaitFor(
    context.Background(),
    15,                 // retries
    5*time.Second,      // delay between retries
    2*time.Minute,      // total timeout
    func(ctx context.Context) (*aruba_types.Response[aruba_types.CloudServerResponse], error) {
		return arubaClient.FromCompute().CloudServers().Get(ctx, projectID, serverID, nil)
    },
    func(resp *aruba_types.Response[aruba_types.CloudServerResponse]) (bool, error) {
        // The check function returns true when the desired state is reached.
		return resp.IsSuccess() && resp.Data.Properties.Status == "running", nil
    },
)

// Await the result
finalRespCustom, err := waitFutureCustom.Await(context.Background())
if err != nil {
    log.Fatalf("Error waiting for server with custom settings: %v", err)
}
fmt.Printf("Server %s is ready after custom wait. Status: %s\n", *finalRespCustom.Data.Metadata.ID, finalRespCustom.Data.Properties.Status)
```

### 2.5. Data Types

All request bodies, response objects, and data models are defined in the `pkg/types` package. They are named to correspond with the resources they represent.

-   **Requests**: Typically end in `Request` (e.g., `types.VPCRequest`, `types.CloudServerRequest`).
-   **Responses**: Typically end in `Response` or `List` (e.g., `types.VPCResponse`, `types.VPCList`).
-   **Shared Structures**: Common structures like `ResourceMetadataRequest`, `ResourceMetadataResponse`, and `ResourceStatus` are used across many types.

For a comprehensive guide on all SDK data types, including their structure and usage in requests and responses, please refer to the [SDK Data Types Documentation](./doc/TYPES.md).