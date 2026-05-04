# Filtering Guide

This guide explains how to use filters with the Aruba Cloud Go SDK to query and filter resources based on various criteria.

## Overview

The SDK provides a flexible filtering system via `CallOption` helpers. Pass them directly to `List` (and other read operations) without constructing intermediate parameter structs.

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,cpu:gt:2"),
    aruba.WithSort("name:asc"),
    aruba.WithLimit(50),
)
```

Available call options:

| Option | Description |
|--------|-------------|
| `aruba.WithFilter(expr string)` | Server-side filter expression |
| `aruba.WithSort(expr string)` | Sort expression |
| `aruba.WithLimit(n int)` | Page size |
| `aruba.WithOffset(n int)` | Pagination offset |
| `aruba.WithProjection(expr string)` | Field projection |

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
| Equal    | `eq`   | Exact match              | `status:eq:Active`         |
| Not Equal| `ne`   | Not equal to             | `status:ne:Error`          |
| Greater  | `gt`   | Greater than             | `cpu:gt:2`                 |
| Greater/Equal | `gte` | Greater than or equal | `memory:gte:4096`         |
| Less     | `lt`   | Less than                | `disk:lt:100`              |
| Less/Equal | `lte` | Less than or equal      | `cpu:lte:8`                |
| In       | `in`   | Value in list            | `region:in:us-east,us-west`|
| Not In   | `nin`  | Value not in list        | `status:nin:Error,Failed`  |
| Contains | `like` | Substring match          | `name:like:prod`           |
| Starts With | `sw` | Prefix match            | `name:sw:web-`             |
| Ends With | `ew`  | Suffix match             | `name:ew:-prod`            |

## Simple Filters

### Single Condition

```go
// List Active cloud servers
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active"),
)
```

### Multiple AND Conditions

```go
// Active servers with at least 2 vCPUs and 4 GB RAM
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,cpu:gt:2,memory:gte:4096"),
)
```

Result expression: `status:eq:Active,cpu:gt:2,memory:gte:4096`

### OR Conditions

```go
// Servers that are Active OR Starting
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active;status:eq:Starting"),
)
```

Result expression: `status:eq:Active;status:eq:Starting`

### Complex Filters (AND + OR)

```go
// (environment=production AND memory>=4096) OR (tier=premium AND region IN [us-east-1,us-west-2])
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("environment:eq:production,memory:gte:4096;tier:eq:premium,region:in:us-east-1,us-west-2"),
)
```

Result expression: `environment:eq:production,memory:gte:4096;tier:eq:premium,region:in:us-east-1,us-west-2`

## Practical Examples

### Filter Active Cloud Servers

```go
// List running servers with at least 4 GB RAM, page size 50
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,memory:gte:4096"),
    aruba.WithLimit(50),
)
if err != nil {
    log.Fatalf("List failed: %v", err)
}
fmt.Printf("Found %d servers\n", servers.Total())
for _, s := range servers.Items() {
    fmt.Println("-", s.Name())
}
```

### Filter VPCs by Region

```go
// List VPCs in specific data centers
vpcs, err := arubaClient.FromNetwork().VPCs().List(ctx, proj,
    aruba.WithFilter("location:in:ITBG-Bergamo,ITMI-Milan"),
)
```

### Filter by Name Pattern

```go
// All cloud servers whose name starts with "web-"
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("name:sw:web-"),
)
```

### Complex Business Logic

```go
// Production servers with high resources OR development servers in specific regions
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter(
        "environment:eq:production,cpu:gte:8,memory:gte:16384" +
        ";environment:eq:development,region:in:ITBG-Bergamo,ITMI-Milan",
    ),
)
```

## Sorting

```go
// Sort by name ascending
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithSort("name:asc"),
)

// Sort by creation date descending
servers, err = arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithSort("createdAt:desc"),
)
```

## Pagination

```go
const pageSize = 25

// First page
page1, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithLimit(pageSize),
    aruba.WithOffset(0),
)

// Second page
page2, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithLimit(pageSize),
    aruba.WithOffset(pageSize),
)

fmt.Printf("Total resources: %d\n", page1.Total())
```

## Best Practices

### Use Specific Filters

Be as specific as possible to reduce the amount of data transferred:

```go
// Good: specific filter — only the resources you need
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("status:eq:Active,region:eq:ITBG-Bergamo"),
)

// Less efficient: no filter — fetches everything
servers, err = arubaClient.FromCompute().CloudServers().List(ctx, proj)
```

### Combine Filter and Pagination

```go
servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter("environment:eq:production"),
    aruba.WithLimit(50),
    aruba.WithOffset(0),
)
```

### Validate Filter Strings

Print the filter string to debug unexpected results:

```go
filter := "status:eq:Active,cpu:gt:4"
fmt.Println("Filter:", filter)
// Output: Filter: status:eq:Active,cpu:gt:4

servers, err := arubaClient.FromCompute().CloudServers().List(ctx, proj,
    aruba.WithFilter(filter),
)
```

## Troubleshooting

### Filter Not Working

**Problem**: Filter doesn't return expected results

**Solutions**:
1. Check field names match the API schema exactly (case-sensitive)
2. Verify the operator is appropriate for the field type
3. Ensure values are the correct format (e.g., `Active` not `active`)
4. Print the filter string to debug

### Empty Results

**Problem**: Filter returns no results

**Solutions**:
1. Verify the filter logic is correct
2. Try simpler filters to isolate the issue
3. Check if the field supports the operator being used
4. List without filters to confirm resources exist

### Complex Filter Issues

**Problem**: Complex AND/OR logic not working as expected

**Solutions**:
1. Break down complex filters into simpler parts
2. Test each condition separately
3. Remember: commas (`,`) = AND, semicolons (`;`) = OR
4. Use parentheses in your mind to understand grouping

```
// This: status:eq:Active,cpu:gt:2;memory:gte:4096
// Means: (status = Active AND cpu > 2) OR (memory >= 4096)

// Not: status = Active AND (cpu > 2 OR memory >= 4096)
```
