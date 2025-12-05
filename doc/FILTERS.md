# Filtering Guide

This guide explains how to use filters with the Aruba Cloud Go SDK to query and filter resources based on various criteria.

## Overview

The SDK provides a powerful and flexible filtering system that follows the Aruba Cloud API filter specification. Filters allow you to:

- Query resources based on field values
- Use comparison operators (equal, greater than, less than, etc.)
- Combine multiple conditions with logical operators (AND/OR)
- Perform pattern matching (contains, starts with, ends with)
- Filter by multiple values (IN, NOT IN)

## Filter Format

Filters follow this format: `field:operator:value`

- **Field**: The resource field to filter on (e.g., `status`, `name`, `cpu`)
- **Operator**: The comparison operator (e.g., `eq`, `gt`, `like`)
- **Value**: The value to compare against

Multiple filters are combined using:
- `,` (comma) for **AND** operations
- `;` (semicolon) for **OR** operations

## Supported Operators

| Operator | Code   | Description              | Example                    |
|----------|--------|--------------------------|----------------------------|
| Equal    | `eq`   | Exact match              | `status:eq:running`        |
| Not Equal| `ne`   | Not equal to             | `status:ne:stopped`        |
| Greater  | `gt`   | Greater than             | `cpu:gt:2`                 |
| Greater/Equal | `gte` | Greater than or equal | `memory:gte:4096`         |
| Less     | `lt`   | Less than                | `disk:lt:100`              |
| Less/Equal | `lte` | Less than or equal      | `cpu:lte:8`                |
| In       | `in`   | Value in list            | `region:in:us-east,us-west`|
| Not In   | `nin`  | Value not in list        | `status:nin:failed,error`  |
| Contains | `like` | Substring match          | `name:like:prod`           |
| Starts With | `sw` | Prefix match            | `name:sw:web-`             |
| Ends With | `ew`  | Suffix match             | `name:ew:-prod`            |

## Using FilterBuilder

The SDK provides a `FilterBuilder` for constructing complex filter expressions programmatically.

### Simple Filters

#### Single Condition

```go
import "github.com/Arubacloud/sdk-go/pkg/spec/schema"

// Filter by status
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    Build()

params := &schema.ListParams{
    Filter: &filter,
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

Result: `status:eq:running`

#### Multiple AND Conditions

```go
// Filter by status AND cpu AND memory
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 2).
    GreaterThanOrEqual("memory", 4096).
    Build()
```

Result: `status:eq:running,cpu:gt:2,memory:gte:4096`

### OR Conditions

```go
// Filter by status = running OR status = starting
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    Or().
    Equal("status", "starting").
    Build()
```

Result: `status:eq:running;status:eq:starting`

### Complex Filters (AND + OR)

```go
// (environment = production AND memory >= 4096) OR (tier = premium AND region IN [us-east-1, us-west-2])
filter := schema.NewFilterBuilder().
    Equal("environment", "production").
    GreaterThanOrEqual("memory", 4096).
    Or().
    Equal("tier", "premium").
    In("region", "us-east-1", "us-west-2").
    Build()
```

Result: `environment:eq:production,memory:gte:4096;tier:eq:premium,region:in:us-east-1,us-west-2`

## Filter Methods

### Comparison Methods

```go
fb := schema.NewFilterBuilder()

// Equality
fb.Equal("field", "value")           // field = value
fb.NotEqual("field", "value")        // field != value

// Numeric comparisons
fb.GreaterThan("field", 100)         // field > 100
fb.GreaterThanOrEqual("field", 100)  // field >= 100
fb.LessThan("field", 100)            // field < 100
fb.LessThanOrEqual("field", 100)     // field <= 100

// List operations
fb.In("field", "val1", "val2", "val3")     // field IN (val1, val2, val3)
fb.NotIn("field", "val1", "val2")          // field NOT IN (val1, val2)

// String pattern matching
fb.Contains("field", "substring")    // field LIKE %substring%
fb.StartsWith("field", "prefix")     // field LIKE prefix%
fb.EndsWith("field", "suffix")       // field LIKE %suffix
```

### Logical Operators

```go
fb := schema.NewFilterBuilder()

// Default is AND
fb.Equal("field1", "value1").
   Equal("field2", "value2")  // field1 = value1 AND field2 = value2

// Explicit OR
fb.Equal("field1", "value1").
   Or().
   Equal("field2", "value2")  // field1 = value1 OR field2 = value2

// Mix AND and OR
fb.Equal("field1", "value1").
   Equal("field2", "value2").  // Group 1: AND
   Or().
   Equal("field3", "value3")   // Group 2: OR
```

## Helper Functions

### Quick Filters

```go
// Simple equality filter
filter := schema.FilterEqual("status", "running")
// Result: status:eq:running

// Single filter with custom operator
filter := schema.BuildFilter("cpu", schema.FilterOpGreaterThan, 4)
// Result: cpu:gt:4
```

### Combining Multiple Filters

```go
// Multiple AND conditions
filter := schema.FilterAnd(
    schema.Filter{Field: "status", Operator: schema.FilterOpEqual, Value: "running"},
    schema.Filter{Field: "cpu", Operator: schema.FilterOpGreaterThan, Value: 2},
    schema.Filter{Field: "memory", Operator: schema.FilterOpGreaterThanOrEqual, Value: 4096},
)
// Result: status:eq:running,cpu:gt:2,memory:gte:4096

// Multiple OR conditions
filter := schema.FilterOr(
    schema.Filter{Field: "status", Operator: schema.FilterOpEqual, Value: "running"},
    schema.Filter{Field: "status", Operator: schema.FilterOpEqual, Value: "starting"},
)
// Result: status:eq:running;status:eq:starting
```

## Practical Examples

### Filter Active Cloud Servers

```go
// List all running cloud servers with at least 4GB RAM
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    GreaterThanOrEqual("memory", 4096).
    Build()

params := &schema.ListParams{
    Filter: &filter,
    Limit:  ptrInt(50),
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### Filter by Multiple Regions

```go
// List resources in US East or US West regions
filter := schema.NewFilterBuilder().
    In("region", "us-east-1", "us-east-2", "us-west-1", "us-west-2").
    Build()

params := &schema.ListParams{
    Filter: &filter,
}

resp, err := sdk.Network.ListVpcs(ctx, "my-project", params)
```

### Filter by Name Pattern

```go
// List all production web servers
filter := schema.NewFilterBuilder().
    StartsWith("name", "web-").
    Contains("environment", "prod").
    Build()

params := &schema.ListParams{
    Filter: &filter,
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### Complex Business Logic

```go
// Find servers that are either:
// - Production servers with high resources (cpu >= 8 AND memory >= 16GB)
// - OR Development servers in specific regions
filter := schema.NewFilterBuilder().
    Equal("environment", "production").
    GreaterThanOrEqual("cpu", 8).
    GreaterThanOrEqual("memory", 16384).
    Or().
    Equal("environment", "development").
    In("region", "us-east-1", "eu-west-1").
    Build()

params := &schema.ListParams{
    Filter: &filter,
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### Filter Deleted Resources

```go
// Include deleted resources in results
accept := schema.AcceptHeaderJsonDeletedTrue

// List all cloud servers including deleted ones
params := &schema.ListParams{
    Accept: &accept,
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### Filter Only Deleted Resources

```go
// Show only deleted resources
accept := schema.AcceptHeaderJsonDeletedOnly

params := &schema.ListParams{
    Accept: &accept,
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### Combining Filters with Labels

```go
// Filter by both labels and field filters
labels := schema.LabelSelector("app=nginx,tier=frontend")
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 2).
    Build()

params := &schema.ListParams{
    Labels: &labels,
    Filter: &filter,
    Limit:  ptrInt(100),
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

## Filter Operators Reference

### Numeric Filters

```go
// CPU cores
schema.NewFilterBuilder().
    Equal("cpu", 4).              // Exactly 4 cores
    GreaterThan("cpu", 2).        // More than 2 cores
    GreaterThanOrEqual("cpu", 4). // 4 or more cores
    LessThan("cpu", 8).           // Fewer than 8 cores
    LessThanOrEqual("cpu", 4).    // 4 or fewer cores
    In("cpu", 2, 4, 8).          // 2, 4, or 8 cores
    Build()
```

### String Filters

```go
// Name matching
schema.NewFilterBuilder().
    Equal("name", "web-server-01").        // Exact match
    NotEqual("name", "web-server-01").     // Not equal
    Contains("name", "prod").              // Contains "prod"
    StartsWith("name", "web-").            // Starts with "web-"
    EndsWith("name", "-prod").             // Ends with "-prod"
    In("name", "server1", "server2").     // In list
    Build()
```

### Status Filters

```go
// Resource status
schema.NewFilterBuilder().
    Equal("status", "running").                           // Running
    NotEqual("status", "failed").                        // Not failed
    In("status", "running", "starting", "stopping").     // Multiple states
    NotIn("status", "failed", "error", "deleted").       // Exclude states
    Build()
```

## Best Practices

### 1. Use Specific Filters

Be as specific as possible to reduce the amount of data transferred:

```go
// Good: Specific filter
filter := schema.NewFilterBuilder().
    Equal("status", "running").
    Equal("region", "us-east-1").
    Build()

// Less efficient: No filter, process all results
resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", nil)
```

### 2. Combine with Pagination

Use filters with pagination for large result sets:

```go
filter := schema.NewFilterBuilder().
    Equal("environment", "production").
    Build()

params := &schema.ListParams{
    Filter: &filter,
    Limit:  ptrInt(50),
}

resp, err := sdk.Compute.ListCloudServers(ctx, "my-project", params)
```

### 3. Validate Filter Logic

Test your filter expressions to ensure they produce the expected results:

```go
fb := schema.NewFilterBuilder().
    Equal("status", "running").
    GreaterThan("cpu", 4)

filterStr := fb.Build()
fmt.Println("Filter:", filterStr)
// Output: Filter: status:eq:running,cpu:gt:4
```

### 4. Use Type-Safe Values

Use the correct types for filter values:

```go
// Good: Correct types
fb.Equal("cpu", 4)           // int
fb.Equal("status", "running") // string
fb.GreaterThan("memory", 4096) // int

// Avoid: String representation of numbers when numbers are expected
fb.Equal("cpu", "4") // May not work as expected
```

### 5. Handle Special Characters

The SDK automatically URL-encodes filter values:

```go
// String with spaces and special characters
filter := schema.NewFilterBuilder().
    Contains("description", "production environment").
    Build()
// Automatically URL-encoded
```

## Troubleshooting

### Filter Not Working

**Problem**: Filter doesn't return expected results

**Solutions**:
1. Check field names match the API schema exactly (case-sensitive)
2. Verify the operator is appropriate for the field type
3. Ensure values are the correct type (string, int, bool)
4. Print the filter string to debug: `fmt.Println(fb.Build())`

### Empty Results

**Problem**: Filter returns no results

**Solutions**:
1. Verify the filter logic is correct
2. Try simpler filters to isolate the issue
3. Check if the field supports the operator being used
4. Test without filters to confirm resources exist

### Complex Filter Issues

**Problem**: Complex AND/OR logic not working as expected

**Solutions**:
1. Break down complex filters into simpler parts
2. Test each condition separately
3. Remember: commas (`,`) = AND, semicolons (`;`) = OR
4. Use parentheses in your mind to understand grouping

```go
// This: status:eq:running,cpu:gt:2;memory:gte:4096
// Means: (status = running AND cpu > 2) OR (memory >= 4096)

// Not: status = running AND (cpu > 2 OR memory >= 4096)
```

## API Reference

For complete API documentation and additional filter options, see:
- [Aruba Cloud API - Filters](https://api.arubacloud.com/docs/filters)
