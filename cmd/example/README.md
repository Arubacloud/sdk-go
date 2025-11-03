# Example Usage

This directory contains example code demonstrating how to use the Aruba Cloud Go SDK.

## Basic Usage

See `main.go` for a complete example of:

- Configuring the SDK with OAuth2 client credentials
- Initializing the client (automatic token acquisition)
- Using contexts for request management
- Accessing resource providers

## Running the Example

```bash
# Set your credentials as environment variables
export ARUBA_CLIENT_ID="your-client-id"
export ARUBA_CLIENT_SECRET="your-client-secret"
export ARUBA_BASE_URL="https://api.arubacloud.com"
export ARUBA_TOKEN_URL="https://auth.arubacloud.com/oauth2/token"

# Run the example
go run main.go
```

## With Generated Clients

After generating client code from your Swagger files:

```bash
make generate
```

You can use the resource providers like this:

```go
// Initialize SDK
sdk, err := client.NewClient(config)

// Use Network service
network := sdk.Network()
networks, err := network.ListNetworks(ctx)

// Use Compute service
compute := sdk.Compute()
vms, err := compute.ListVMs(ctx)
```

## Environment Configuration

Create a `.env` file (not committed to git):

```env
ARUBA_BASE_URL=https://api.arubacloud.com
ARUBA_TOKEN_URL=https://auth.arubacloud.com/oauth2/token
ARUBA_CLIENT_ID=your-client-id
ARUBA_CLIENT_SECRET=your-client-secret
```
