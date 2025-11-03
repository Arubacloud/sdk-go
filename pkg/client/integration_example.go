package client

// This file contains examples of how to integrate generated clients
// After running 'make generate', you can integrate the generated clients here

// Example integration (uncomment after code generation):
//
// import (
//     "context"
//     "github.com/Arubacloud/sdk-go/pkg/generated/network"
// )
//
// type NetworkClient struct {
//     client *network.ClientWithResponses
//     sdk    *Client
// }
//
// // Network returns the Network resource provider client
// func (c *Client) Network() *NetworkClient {
//     if c.networkClient == nil {
//         client, _ := network.NewClientWithResponses(
//             c.config.BaseURL,
//             network.WithHTTPClient(c.config.HTTPClient),
//             network.WithRequestEditorFn(c.RequestEditorFn()),
//         )
//         c.networkClient = &NetworkClient{
//             client: client,
//             sdk:    c,
//         }
//     }
//     return c.networkClient
// }
//
// // Example method wrapper
// func (n *NetworkClient) ListNetworks(ctx context.Context) (*network.NetworkListResponse, error) {
//     resp, err := n.client.ListNetworksWithResponse(ctx)
//     if err != nil {
//         return nil, err
//     }
//     if resp.StatusCode() != 200 {
//         return nil, NewError(resp.StatusCode(), "failed to list networks", resp.Body, nil)
//     }
//     return resp.JSON200, nil
// }

// Instructions:
// 1. Run 'make generate' to generate client code from Swagger files
// 2. Import the generated packages (e.g., "github.com/Arubacloud/sdk-go/pkg/generated/network")
// 3. Create wrapper types for each resource provider
// 4. Add accessor methods to the Client struct
// 5. Implement high-level methods that wrap the generated client calls
// 6. Add proper error handling and type conversions
