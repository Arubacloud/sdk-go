package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// runUpdateExample demonstrates how to update existing resources
// To run: PROJECT_ID=your-project go run . -mode=update
func runUpdateExample(clientID, clientSecret string) {
	// Initialize the SDK
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("\n=== Update Example ===")

	// Get project ID from environment
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("Please set PROJECT_ID environment variable")
	}

	// Fetch existing resources
	resources := fetchExistingResources(ctx, arubaClient, projectID)

	// Update all resources
	updateAllResources(ctx, arubaClient, resources)

	fmt.Println("\n=== Update Example Complete ===")
}

// fetchExistingResources retrieves existing resources from the API for updating
func fetchExistingResources(ctx context.Context, arubaClient aruba.Client, projectID string) *ResourceCollection {
	resources := &ResourceCollection{
		ProjectID: projectID,
	}

	fmt.Println("Fetching existing resources...")

	// Fetch DBaaS instances
	dbaasListResp, err := arubaClient.FromDatabase().DBaaS().List(ctx, projectID, nil)
	if err == nil && dbaasListResp.IsSuccess() && len(dbaasListResp.Data.Values) > 0 {
		// Get the first DBaaS instance details
		dbaasID := *dbaasListResp.Data.Values[0].Metadata.ID
		dbaasResp, err := arubaClient.FromDatabase().DBaaS().Get(ctx, projectID, dbaasID, nil)
		if err == nil && dbaasResp.IsSuccess() {
			resources.DBaaSResp = dbaasResp
			fmt.Printf("✓ Found DBaaS: %s\n", *dbaasResp.Data.Metadata.Name)
		}
	}

	// Fetch KaaS clusters
	kaasList, err := arubaClient.FromContainer().KaaS().List(ctx, projectID, nil)
	if err == nil && kaasList.IsSuccess() && len(kaasList.Data.Values) > 0 {
		// Get the first KaaS cluster details
		kaasID := *kaasList.Data.Values[0].Metadata.ID
		kaasResp, err := arubaClient.FromContainer().KaaS().Get(ctx, projectID, kaasID, nil)
		if err == nil && kaasResp.IsSuccess() {
			resources.KaaSResp = kaasResp
			fmt.Printf("✓ Found KaaS: %s\n", *kaasResp.Data.Metadata.Name)
		}
	}

	return resources
}

// updateAllResources updates all resources that support update operations
func updateAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Updating Resources ===")

	// Update Project
	updateProject(ctx, arubaClient, resources.ProjectID)

	// Update DBaaS (if created)
	if resources.DBaaSResp != nil && resources.DBaaSResp.Data != nil {
		updateDBaaS(ctx, arubaClient, resources.ProjectID, resources.DBaaSResp)
	}

	// Update KaaS (if created)
	if resources.KaaSResp != nil && resources.KaaSResp.Data != nil {
		updateKaaS(ctx, arubaClient, resources.ProjectID, resources.KaaSResp)
	}

	// Update KMS Key (if you create one in the future)
	// updateKMSKey(ctx, sdk, resources.ProjectID, resources.KMSResp)

	fmt.Println("\n=== Update Complete ===")
}

// updateProject updates a project
func updateProject(ctx context.Context, arubaClient aruba.Client, projectID string) {
	fmt.Println("--- Updating Project ---")

	projectReq := types.ProjectRequest{
		Metadata: types.ResourceMetadataRequest{
			Name: "seca-sdk-example-updated",
			Tags: []string{"production", "arubacloud-sdk", "updated"},
		},
		Properties: types.ProjectPropertiesRequest{
			Description: stringPtr("My production project - UPDATED"),
			Default:     false,
		},
	}

	updateResp, err := arubaClient.FromProject().Update(ctx, projectID, projectReq, nil)
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
func updateDBaaS(ctx context.Context, arubaClient aruba.Client, projectID string, dbaasResp *types.Response[types.DBaaSResponse]) {
	fmt.Println("--- Updating DBaaS ---")

	dbaasID := *dbaasResp.Data.Metadata.ID

	// Update with new storage size
	dbaasReq := types.DBaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-dbaas-updated",
				Tags: []string{"database", "mysql", "updated"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.DBaaSPropertiesRequest{
			Engine: &types.DBaaSEngine{
				ID:         stringPtr("mysql-8.0"),
				DataCenter: stringPtr("ITBG-1"),
			},
			Flavor: &types.DBaaSFlavor{
				Name: stringPtr("DBO2A4"),
			},
			Storage: &types.DBaaSStorage{
				SizeGB: int32Ptr(25), // Increased from 20 to 25 GB
			},
			BillingPlan: &types.DBaaSBillingPlan{
				BillingPeriod: stringPtr("Hour"),
			},
			Networking: &types.DBaaSNetworking{
				VPCURI:           &dbaasResp.Data.Properties.Networking.VPC.URI,
				SubnetURI:        &dbaasResp.Data.Properties.Networking.Subnet.URI,
				SecurityGroupURI: &dbaasResp.Data.Properties.Networking.SecurityGroup.URI,
			},
			Autoscaling: &types.DBaaSAutoscaling{
				Enabled:        boolPtr(true),
				AvailableSpace: int32Ptr(25), // Updated
				StepSize:       int32Ptr(15), // Increased step size
			},
		},
	}

	updateResp, err := arubaClient.FromDatabase().DBaaS().Update(ctx, projectID, dbaasID, dbaasReq, nil)
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
func updateKaaS(ctx context.Context, arubaClient aruba.Client, projectID string, kaasResp *types.Response[types.KaaSResponse]) {
	fmt.Println("--- Updating KaaS Cluster ---")

	kaasID := *kaasResp.Data.Metadata.ID

	// Update with modified node pool and Kubernetes version
	kaasUpdateReq := types.KaaSUpdateRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-kaas-cluster-updated",
				Tags: []string{"kubernetes", "container", "updated"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.KaaSPropertiesUpdateRequest{
			KubernetesVersion: types.KubernetesVersionInfoUpdate{
				Value: stringValue(kaasResp.Data.Properties.KubernetesVersion.Value),
			},
			NodePools: []types.NodePoolProperties{
				{
					Name:     "default-pool",
					Nodes:    5, // Increased from 3 to 5 nodes
					Instance: "K4A8",
					Zone:     "ITBG-1",
				},
			},
			HA: boolPtr(true),
			Storage: &types.StorageKubernetes{
				MaxCumulativeVolumeSize: int32Ptr(100),
			},
			BillingPlan: &types.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	updateResp, err := arubaClient.FromContainer().KaaS().Update(ctx, projectID, kaasID, kaasUpdateReq, nil)
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
		if updateResp.Data.Properties.NodePools != nil {
			for _, pool := range *updateResp.Data.Properties.NodePools {
				if pool.Nodes != nil {
					totalNodes += *pool.Nodes
				}
			}
		}
		fmt.Printf("✓ Updated KaaS cluster: %s (Total Nodes: %d)\n",
			*updateResp.Data.Metadata.Name,
			totalNodes)
	}
}
