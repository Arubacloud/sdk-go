package main

// This example demonstrates how to use the Aruba Cloud SDK with proper error handling.
//
// Response handling pattern:
//   - response.IsSuccess() (2xx): Access response.Data (typed success response)
//   - response.IsError() (4xx/5xx): Access response.Error (ErrorResponse with Title, Detail, etc.)
//   - response.RawBody: Always available for debugging/logging
//
// Example usage:
//   resp, err := api.CreateResource(ctx, ...)
//   if err != nil {
//       // Network or SDK error
//       log.Printf("Error: %v", err)
//   } else if resp.IsSuccess() {
//       // Success - use resp.Data
//       fmt.Printf("Created: %s\n", *resp.Data.Metadata.Name)
//   } else if resp.IsError() && resp.Error != nil {
//       // API error - use resp.Error
//       log.Printf("API Error: %s - %s", *resp.Error.Title, *resp.Error.Detail)
//   }

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/network"
	"github.com/Arubacloud/sdk-go/pkg/spec/project"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
	"github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

func main() {
	config := &client.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		Debug:        true,
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

	fmt.Println("\n=== SDK Examples ===")

	// Initialize service clients
	projectAPI := project.NewProjectService(sdk)
	elasticIPAPI := network.NewElasticIPService(sdk)
	vpcAPI := network.NewVPCService(sdk)
	storageAPI := storage.NewBlockStorageService(sdk)

	// Example: Project Management
	fmt.Println("--- Project Management ---")

	// Create a new project
	projectReq := schema.ProjectRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "my-project",
			Tags: []string{"production", "arubacloud-sdk"},
		},
		Properties: schema.ProjectPropertiesRequest{
			Description: stringPtr("My production project"),
			Default:     false,
		},
	}

	createResp, err := projectAPI.CreateProject(ctx, projectReq, nil)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}

	if !createResp.IsSuccess() {
		log.Fatalf("Failed to create project, status code: %d", createResp.StatusCode)
	}

	projectID := *createResp.Data.Metadata.Id
	fmt.Printf("✓ Created project with ID: %s\n", projectID)

	// Update the project
	updateResp, err := projectAPI.UpdateProject(ctx, projectID, projectReq, nil)
	if err != nil {
		log.Printf("Error updating project: %v", err)
	} else if updateResp.IsSuccess() {
		fmt.Printf("✓ Updated project: %s\n", *updateResp.Data.Metadata.Name)
	}

	// Example: Create Elastic IP
	fmt.Println("\n--- Network: Elastic IP ---")

	elasticIPReq := schema.ElasticIpRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-elastic-ip",
				Tags: []string{"network", "public"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.ElasticIpPropertiesRequest{
			BillingPlan: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	elasticIPResp, err := elasticIPAPI.CreateElasticIP(ctx, projectID, elasticIPReq, nil)
	if err != nil {
		log.Printf("Error creating Elastic IP: %v", err)
	} else if elasticIPResp.IsSuccess() {
		fmt.Printf("✓ Created Elastic IP: %s (ObjectId: %s)\n",
			*elasticIPResp.Data.Metadata.Name, *elasticIPResp.Data.Metadata.Id)
	} else if elasticIPResp.IsError() && elasticIPResp.Error != nil {
		log.Printf("Failed to create Elastic IP - Status: %d, Error: %s, Detail: %s",
			elasticIPResp.StatusCode,
			stringValue(elasticIPResp.Error.Title),
			stringValue(elasticIPResp.Error.Detail))
	}

	// Example: Create Block Storage
	fmt.Println("\n--- Storage: Block Storage ---")

	blockStorageReq := schema.BlockStorageRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-block-storage",
				Tags: []string{"storage", "data"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.BlockStoragePropertiesRequest{
			SizeGB: 10,
			Type:   schema.BlockStorageTypeStandard,
			Zone:   "ITBG-1",
			BillingPeriod: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	blockStorageResp, err := storageAPI.CreateBlockStorageVolume(ctx, projectID, blockStorageReq, nil)
	if err != nil {
		log.Printf("Error creating block storage: %v", err)
	} else if blockStorageResp.IsSuccess() {
		fmt.Printf("✓ Created block storage: %s (%d GB, %s)\n",
			*blockStorageResp.Data.Metadata.Name,
			blockStorageResp.Data.Properties.SizeGB,
			blockStorageResp.Data.Properties.Type)
	} else if blockStorageResp.IsError() && blockStorageResp.Error != nil {
		log.Printf("Failed to create block storage - Status: %d, Error: %s, Detail: %s",
			blockStorageResp.StatusCode,
			stringValue(blockStorageResp.Error.Title),
			stringValue(blockStorageResp.Error.Detail))
	}

	// Example: Create VPC
	fmt.Println("\n--- Network: VPC ---")

	vpcReq := schema.VpcRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-vpc",
				Tags: []string{"network", "infrastructure"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.VpcPropertiesRequest{
			Properties: &schema.VpcProperties{
				Default: boolPtr(false),
				Preset:  boolPtr(true),
			},
		},
	}

	vpcResp, err := vpcAPI.CreateVPC(ctx, projectID, vpcReq, nil)
	if err != nil {
		log.Printf("Error creating VPC: %v", err)
	} else if vpcResp.IsSuccess() {
		fmt.Printf("✓ Created VPC: %s (Default: %t)\n",
			*vpcResp.Data.Metadata.Name,
			vpcResp.Data.Properties.Default)
	} else if vpcResp.IsError() && vpcResp.Error != nil {
		log.Printf("Failed to create VPC - Status: %d, Error: %s, Detail: %s",
			vpcResp.StatusCode,
			stringValue(vpcResp.Error.Title),
			stringValue(vpcResp.Error.Detail))
	}

	fmt.Println("\n=== SDK Example Complete ===")
	fmt.Println("Successfully created project:")
	fmt.Println("- Project ID:", projectID)
	if elasticIPResp != nil && elasticIPResp.IsSuccess() {
		fmt.Println("- Elastic IP:", *elasticIPResp.Data.Metadata.Id)
	}
	if blockStorageResp != nil && blockStorageResp.IsSuccess() {
		fmt.Println("- Block Storage (100 GB):", *blockStorageResp.Data.Metadata.Id)
	}
	if vpcResp != nil && vpcResp.IsSuccess() {
		fmt.Println("- VPC:", *vpcResp.Data.Metadata.Id)
	}
}

// Helper for pointer types
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
