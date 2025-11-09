# Response Handling Guide

## Overview

The SDK uses a generic `Response[T]` type that properly handles both success and error responses from the API. Response parsing is centralized through the `ParseResponseBody[T]` function for consistency.

## Response Structure

```go
type Response[T any] struct {
    Data         *T              // Populated for 2xx responses
    Error        *ErrorResponse  // Populated for 4xx/5xx responses
    HTTPResponse *http.Response  // The underlying HTTP response
    StatusCode   int             // HTTP status code
    Headers      http.Header     // Response headers
    RawBody      []byte          // Raw response body (always available)
}
```

## Centralized Response Parsing

All service methods use the `ParseResponseBody[T]` function to handle response parsing:

```go
func ParseResponseBody[T any](httpResp *http.Response) (*Response[T], error) {
    // Reads the response body
    // Creates the Response[T] wrapper
    // Parses into Data for 2xx responses
    // Parses into Error for 4xx/5xx responses
    return response, nil
}
```

### Implementation in Services

Service methods simply call `ParseResponseBody` after the HTTP request:

```go
func (s *VPCService) GetVPC(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VpcResponse], error) {
    // ... prepare request ...
    
    httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
    if err != nil {
        return nil, err
    }
    defer httpResp.Body.Close()

    return schema.ParseResponseBody[schema.VpcResponse](httpResp)
}
```

**Benefits:**
- ✅ Eliminates code duplication across all services
- ✅ Ensures consistent error handling
- ✅ Simplifies service implementations
- ✅ Makes updates easier to maintain

## Error Response Structure

```go
type ErrorResponse struct {
    Type       *string                 // URI reference for the problem type
    Title      *string                 // Short, human-readable summary
    Status     *int32                  // HTTP status code
    Detail     *string                 // Human-readable explanation
    Instance   *string                 // URI for this specific occurrence
    Extensions map[string]interface{}  // Additional dynamic properties
}
```

## Response Handling Pattern

### 1. Basic Pattern

```go
resp, err := api.CreateResource(ctx, projectID, request, nil)
if err != nil {
    // Network error, context timeout, or SDK error
    log.Fatalf("Request failed: %v", err)
}

if resp.IsSuccess() {
    // 2xx - Success response
    fmt.Printf("Created: %s\n", *resp.Data.Metadata.Name)
} else if resp.IsError() && resp.Error != nil {
    // 4xx/5xx - API error response
    log.Printf("API Error: %s - %s", 
        stringValue(resp.Error.Title), 
        stringValue(resp.Error.Detail))
}
```

### 2. Complete Error Handling

```go
resp, err := api.GetResource(ctx, projectID, resourceID, nil)
if err != nil {
    return fmt.Errorf("request failed: %w", err)
}

switch {
case resp.IsSuccess():
    // Handle success - resp.Data is populated
    resource := resp.Data
    fmt.Printf("Resource: %s (Status: %s)\n", 
        *resource.Metadata.Name, 
        *resource.Status.State)
    return nil

case resp.StatusCode == 404:
    // Handle not found
    return fmt.Errorf("resource not found")

case resp.StatusCode == 400:
    // Handle validation errors
    if resp.Error != nil {
        return fmt.Errorf("validation error: %s", stringValue(resp.Error.Detail))
    }
    return fmt.Errorf("bad request: %s", string(resp.RawBody))

case resp.IsError():
    // Handle other errors
    if resp.Error != nil {
        return fmt.Errorf("API error (%d): %s - %s", 
            resp.StatusCode,
            stringValue(resp.Error.Title),
            stringValue(resp.Error.Detail))
    }
    return fmt.Errorf("unexpected error (%d): %s", resp.StatusCode, string(resp.RawBody))

default:
    // Unexpected status code
    return fmt.Errorf("unexpected status %d", resp.StatusCode)
}
```

### 3. Accessing Error Details

```go
if resp.IsError() && resp.Error != nil {
    // Standard fields
    log.Printf("Error Type: %s", stringValue(resp.Error.Type))
    log.Printf("Error Title: %s", stringValue(resp.Error.Title))
    log.Printf("Error Detail: %s", stringValue(resp.Error.Detail))
    log.Printf("Status Code: %d", int32Value(resp.Error.Status))
    
    // Access custom extensions (e.g., validation errors)
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

### 4. Raw Body Access

```go
// Always available for debugging
log.Printf("Raw response: %s", string(resp.RawBody))

// Useful for logging full responses during development
if !resp.IsSuccess() {
    log.Printf("Request failed with status %d: %s", 
        resp.StatusCode, 
        string(resp.RawBody))
}
```

## Helper Functions

```go
// Safe pointer dereference helpers
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

## Common Error Scenarios

### 400 Bad Request - Validation Errors

```go
resp, err := api.CreateResource(ctx, projectID, invalidRequest, nil)
if err != nil {
    log.Fatalf("Request failed: %v", err)
}

if resp.StatusCode == 400 && resp.Error != nil {
    fmt.Printf("Validation failed: %s\n", stringValue(resp.Error.Title))
    
    // Check for field-level errors in Extensions
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

// Use resp.Data
resource := resp.Data
```

### 500 Internal Server Error

```go
if resp.StatusCode >= 500 {
    // Server error - may want to retry
    if resp.Error != nil {
        log.Printf("Server error: %s", stringValue(resp.Error.Detail))
    }
    
    // Log trace ID for support
    if traceID, ok := resp.Error.Extensions["traceId"].(string); ok {
        log.Printf("Trace ID: %s", traceID)
    }
}
```

## Best Practices

1. **Always check for network errors first** (`err != nil`)
2. **Use `IsSuccess()` to check for 2xx responses** before accessing `Data`
3. **Use `IsError()` to check for 4xx/5xx responses** before accessing `Error`
4. **Check if `Error` field is non-nil** before dereferencing
5. **Use helper functions** to safely dereference pointer fields
6. **Keep `RawBody` available** for debugging and logging
7. **Log trace IDs** from error responses for support requests

## Testing Response Handling

```go
func TestResourceCreation(t *testing.T) {
    resp, err := api.CreateResource(ctx, projectID, request, nil)
    
    // Check no network error
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    
    // Check successful response
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
    
    // Validate response data
    if resp.Data == nil {
        t.Fatal("Expected data to be populated")
    }
    
    if resp.Data.Metadata.Name == nil {
        t.Fatal("Expected resource name")
    }
}
```

## Complete Example

See `cmd/example/main.go` for a complete working example demonstrating:
- Project creation
- Resource creation (Elastic IP, Block Storage, VPC)
- Success and error handling
- Accessing typed response data
- Accessing error details
