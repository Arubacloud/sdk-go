package restclient

// This file will contain wrappers for generated clients
// Each resource provider (microservice) will have its own wrapper

// Example structure (will be populated after code generation):
//
// type NetworkClient struct {
//     client *network.ClientWithResponses
// }
//
// func (c *Client) Network() *NetworkClient {
//     if c.networkClient == nil {
//         c.networkClient = &NetworkClient{
//             client: network.NewClientWithResponses(c.config.BaseURL, c.config.HTTPClient),
//         }
//     }
//     return c.networkClient
// }
