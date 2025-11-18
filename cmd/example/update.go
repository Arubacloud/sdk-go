package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sdkgo "github.com/Arubacloud/sdk-go"
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// runUpdateExample demonstrates how to update existing resources
// To run: PROJECT_ID=your-project go run . -mode=update
func runUpdateExample() {
	config := &restclient.Config{
		ClientID:     "clientId",
		ClientSecret: "clientSecret",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		Debug:        true,
	}

	// Initialize the SDK
	sdk, err := sdkgo.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	sdk.Client = sdk.Client.WithContext(ctx)

	fmt.Println("\n=== Update Example ===")

	// Get project ID from environment
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("Please set PROJECT_ID environment variable")
	}

	// Fetch existing resources
	resources := fetchExistingResources(ctx, sdk, projectID)

	// Update all resources
	updateAllResources(ctx, sdk, resources)

	fmt.Println("\n=== Update Example Complete ===")
}

// fetchExistingResources retrieves existing resources from the API for updating
func fetchExistingResources(ctx context.Context, sdk *sdkgo.Client, projectID string) *ResourceCollection {
	resources := &ResourceCollection{
		ProjectID: projectID,
	}

	fmt.Println("Fetching existing resources...")

	// Fetch DBaaS instances
	dbaasListResp, err := sdk.Database.ListDBaaS(ctx, projectID, nil)
	if err == nil && dbaasListResp.IsSuccess() && len(dbaasListResp.Data.Values) > 0 {
		// Get the first DBaaS instance details
		dbaasID := *dbaasListResp.Data.Values[0].Metadata.ID
		dbaasResp, err := sdk.Database.GetDBaaS(ctx, projectID, dbaasID, nil)
		if err == nil && dbaasResp.IsSuccess() {
			resources.DBaaSResp = dbaasResp
			fmt.Printf("✓ Found DBaaS: %s\n", *dbaasResp.Data.Metadata.Name)
		}
	}

	// Fetch KaaS clusters
	kaasList, err := sdk.Container.ListKaaS(ctx, projectID, nil)
	if err == nil && kaasList.IsSuccess() && len(kaasList.Data.Values) > 0 {
		// Get the first KaaS cluster details
		kaasID := *kaasList.Data.Values[0].Metadata.ID
		kaasResp, err := sdk.Container.GetKaaS(ctx, projectID, kaasID, nil)
		if err == nil && kaasResp.IsSuccess() {
			resources.KaaSResp = kaasResp
			fmt.Printf("✓ Found KaaS: %s\n", *kaasResp.Data.Metadata.Name)
		}
	}

	return resources
}

// updateAllResources updates all resources that support update operations
func updateAllResources(ctx context.Context, sdk *sdkgo.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Updating Resources ===")

	// Update Project
	updateProject(ctx, sdk, resources.ProjectID)

	// Update DBaaS (if created)
	if resources.DBaaSResp != nil && resources.DBaaSResp.Data != nil {
		updateDBaaS(ctx, sdk, resources.ProjectID, resources.DBaaSResp)
	}

	// Update KaaS (if created)
	if resources.KaaSResp != nil && resources.KaaSResp.Data != nil {
		updateKaaS(ctx, sdk, resources.ProjectID, resources.KaaSResp)
	}

	// Update KMS Key (if you create one in the future)
	// updateKMSKey(ctx, sdk, resources.ProjectID, resources.KMSResp)

	fmt.Println("\n=== Update Complete ===")
}

// updateProject updates a project
func updateProject(ctx context.Context, sdk *sdkgo.Client, projectID string) {
	fmt.Println("--- Updating Project ---")

	projectReq := schema.ProjectRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "seca-sdk-example-updated",
			Tags: []string{"production", "arubacloud-sdk", "updated"},
		},
		Properties: schema.ProjectPropertiesRequest{
			Description: stringPtr("My production project - UPDATED"),
			Default:     false,
		},
	}

	updateResp, err := sdk.Project.UpdateProject(ctx, projectID, projectReq, nil)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		return
	} else if !updateResp.IsSuccess() {
		log.Printf("Failed to update project, status code: %d and error title: %s", updateResp.StatusCode, stringValue(updateResp.Error.Title))
		return
	}
	fmt.Printf("✓ Updated project: %s\n", *updateResp.Data.Metadata.Name)
}

// updateDBaaS updates a DBaaS instance
func updateDBaaS(ctx context.Context, sdk *sdkgo.Client, projectID string, dbaasResp *schema.Response[schema.DBaaSResponse]) {
	fmt.Println("--- Updating DBaaS ---")

	dbaasID := *dbaasResp.Data.Metadata.ID

	// Update with new storage size
	dbaasReq := schema.DBaaSRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-dbaas-updated",
				Tags: []string{"database", "mysql", "updated"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.DBaaSPropertiesRequest{
			Engine: &schema.DBaaSEngine{
				ID:         stringPtr("mysql-8.0"),
				DataCenter: stringPtr("ITBG-1"),
			},
			Flavor: &schema.DBaaSFlavor{
				Name: stringPtr("DBO2A4"),
			},
			Storage: &schema.DBaaSStorage{
				SizeGB: int32Ptr(25), // Increased from 20 to 25 GB
			},
			BillingPlan: &schema.DBaaSBillingPlan{
				BillingPeriod: stringPtr("Hour"),
			},
			Networking: &schema.DBaaSNetworking{
				VPCURI:           &dbaasResp.Data.Properties.Networking.VPC.URI,
				SubnetURI:        &dbaasResp.Data.Properties.Networking.Subnet.URI,
				SecurityGroupURI: &dbaasResp.Data.Properties.Networking.SecurityGroup.URI,
			},
			Autoscaling: &schema.DBaaSAutoscaling{
				Enabled:        boolPtr(true),
				AvailableSpace: int32Ptr(25), // Updated
				StepSize:       int32Ptr(15), // Increased step size
			},
		},
	}

	updateResp, err := sdk.Database.UpdateDBaaS(ctx, projectID, dbaasID, dbaasReq, nil)
	if err != nil {
		log.Printf("Error updating DBaaS: %v", err)
		return
	} else if !updateResp.IsSuccess() {
		log.Printf("Failed to update DBaaS - Status: %d, Error: %s, Detail: %s",
			updateResp.StatusCode,
			stringValue(updateResp.Error.Title),
			stringValue(updateResp.Error.Detail))
		return
	}

	if updateResp.Data != nil && updateResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Updated DBaaS: %s (Storage: %d GB)\n",
			*updateResp.Data.Metadata.Name,
			int32Value(updateResp.Data.Properties.Storage.SizeGB))
	}
}

// updateKaaS updates a KaaS cluster
func updateKaaS(ctx context.Context, sdk *sdkgo.Client, projectID string, kaasResp *schema.Response[schema.KaaSResponse]) {
	fmt.Println("--- Updating KaaS Cluster ---")

	kaasID := *kaasResp.Data.Metadata.ID

	// Update with modified node pool
	kaasReq := schema.KaaSRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-kaas-cluster-updated",
				Tags: []string{"kubernetes", "container", "updated"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.KaaSPropertiesRequest{
			Preset: false,
			VPC: schema.ReferenceResource{
				URI: kaasResp.Data.Properties.VPC.URI,
			},
			Subnet: schema.ReferenceResource{
				URI: kaasResp.Data.Properties.Subnet.URI,
			},
			SecurityGroup: schema.SecurityGroupProperties{
				Name: kaasResp.Data.Properties.SecurityGroup.Name,
			},
			NodeCIDR: schema.NodeCIDRProperties{
				Name:    kaasResp.Data.Properties.NodeCIDR.Name,
				Address: kaasResp.Data.Properties.NodeCIDR.Address,
			},
			KubernetesVersion: schema.KubernetesVersionInfo{
				Value: kaasResp.Data.Properties.KubernetesVersion.Value,
			},
			NodePools: []schema.NodePoolProperties{
				{
					Name:     "default-pool",
					Nodes:    5, // Increased from 3 to 5 nodes
					Instance: "K4A8",
					Zone:     "ITBG-1",
				},
			},
			HA: true,
			BillingPlan: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	updateResp, err := sdk.Container.UpdateKaaS(ctx, projectID, kaasID, kaasReq, nil)
	if err != nil {
		log.Printf("Error updating KaaS cluster: %v", err)
		return
	} else if !updateResp.IsSuccess() {
		log.Printf("Failed to update KaaS cluster - Status: %d, Error: %s, Detail: %s",
			updateResp.StatusCode,
			stringValue(updateResp.Error.Title),
			stringValue(updateResp.Error.Detail))
		return
	}

	if updateResp.Data != nil && updateResp.Data.Metadata.Name != nil {
		totalNodes := int32(0)
		for _, pool := range updateResp.Data.Properties.NodePools {
			totalNodes += pool.Nodes
		}
		fmt.Printf("✓ Updated KaaS cluster: %s (Total Nodes: %d)\n",
			*updateResp.Data.Metadata.Name,
			totalNodes)
	}
}
