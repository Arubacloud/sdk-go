# SDK Example

This directory contains example code demonstrating how to use the Aruba Cloud SDK.

## Files

- **main.go** - Main entry point with mode selection (create/update/delete) and create implementation
- **update.go** - Update mode implementation with functions for updating resources
- **delete.go** - Delete mode implementation with functions for deleting resources

## Usage

### Creating Resources

Creates a complete infrastructure from scratch:

```bash
go run .
# or explicitly:
go run . -mode=create
```

This will:
1. Create a Project
2. Create an Elastic IP
3. Create Block Storage
4. Create a Snapshot from Block Storage
5. Create a VPC
6. Create a Subnet in the VPC
7. Create a Security Group
8. Create Security Group Rules
9. Create an SSH Key Pair
10. Create a DBaaS instance
11. Create a KaaS cluster
12. Create a Cloud Server (optional)

### Updating Resources

Updates existing resources in a project:

```bash
PROJECT_ID=your-project-id go run . -mode=update
```

Updates:
- Project metadata (name, tags, description)
- DBaaS storage size and autoscaling settings
- KaaS cluster node pool size (increases from 3 to 5 nodes)

### Deleting Resources

Deletes all resources in a project (with confirmation):

```bash
PROJECT_ID=your-project-id go run . -mode=delete
```

⚠️ **Warning**: This will delete ALL resources found in the specified project.
You will be prompted to confirm by typing 'yes'.

Deletes resources in this order (respecting dependencies):
1. Cloud Server
2. KaaS Cluster
3. DBaaS Instance
4. SSH Key Pair
5. Security Group Rule
6. Security Group
7. Subnet
8. VPC
9. Snapshot
10. Block Storage
11. Elastic IP
12. Project

## Configuration

Update the credentials in the respective files:

```go
config := &client.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   &http.Client{Timeout: 30 * time.Second},
    Debug:        true,
}
```

## Environment Variables

- `PROJECT_ID` - Required for update and delete modes. Specifies the project to operate on.

## Examples

```bash
# Create all resources
go run .

# Update resources in a specific project
PROJECT_ID=my-project-123 go run . -mode=update

# Delete all resources in a project (with confirmation)
PROJECT_ID=my-project-123 go run . -mode=delete
```

## Notes

- The SDK automatically handles token management (OAuth2 client credentials flow)
- Some operations (like VPC creation) include automatic polling for resource readiness
- Resources are created with proper dependency ordering
- Delete operations respect resource dependencies by deleting in reverse order
- Update and delete modes fetch existing resources from the API before operating on them
