package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/client"
)

func main() {
	// Configure the SDK with OAuth2 client credentials
	config := &client.Config{
		BaseURL:            "https://api.arubacloud.com",
		TokenIssuerURL:     "https://auth.arubacloud.com/oauth2/token",
		ClientID:           "your-client-id",
		ClientSecret:       "your-client-secret",
		TokenRefreshBuffer: 5 * time.Minute,
		Headers: map[string]string{
			"X-Application": "aruba-sdk-example",
		},
	}

	// Initialize the SDK (automatically obtains JWT token)
	sdk, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use the SDK with context
	sdk = sdk.WithContext(ctx)

	// Example: Get current token (optional - SDK manages this automatically)
	token, err := sdk.GetToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("Successfully authenticated. Token: %s...\n", token[:20])

	// Example: Use resource providers (after code generation)
	// network := sdk.Network()
	// result, err := network.ListNetworks(ctx)
	// if err != nil {
	//     log.Fatalf("Failed to list networks: %v", err)
	// }
	// fmt.Printf("Networks: %+v\n", result)

	fmt.Println("SDK initialized successfully!")
}
