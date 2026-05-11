package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// runUpdateExample demonstrates how to update existing resources.
// Run with: go run ./examples/all-resources/ -mode=update -clientID=… -clientSecret=… -projectID=…
func runUpdateExample(clientID, clientSecret, projectID string, debug bool) {
	arubaClient, err := aruba.NewClient(buildClientOptions(clientID, clientSecret, debug))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("\n=== Update Example ===")

	proj, err := arubaClient.FromProject().Get(ctx, aruba.URI("/projects/"+projectID))
	if err != nil {
		log.Fatalf("Error fetching project: %v", err)
	}

	resources := fetchExistingResources(ctx, arubaClient, proj)

	updateAllResources(ctx, arubaClient, resources)

	fmt.Println("\n=== Update Example Complete ===")
}

func fetchExistingResources(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) *ResourceCollection {
	resources := &ResourceCollection{
		Project: proj,
	}

	fmt.Println("Fetching existing resources...")

	dbaasListResp, err := arubaClient.FromDatabase().DBaaS().List(ctx, proj)
	if err == nil && dbaasListResp.Total() > 0 {
		dbaasResp, err := arubaClient.FromDatabase().DBaaS().Get(ctx, dbaasListResp.Items()[0])
		if err == nil {
			resources.DBaaS = dbaasResp
			fmt.Printf("✓ Found DBaaS: %s\n", dbaasResp.Name())
		}
	}

	kaasList, err := arubaClient.FromContainer().KaaS().List(ctx, proj)
	if err == nil && kaasList.Total() > 0 {
		kaasResp, err := arubaClient.FromContainer().KaaS().Get(ctx, kaasList.Items()[0])
		if err == nil {
			resources.KaaS = kaasResp
			fmt.Printf("✓ Found KaaS: %s\n", kaasResp.Name())
		}
	}

	return resources
}

func updateAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Updating Resources ===")

	updateProject(ctx, arubaClient, resources.Project)

	if resources.DBaaS != nil && resources.DBaaS.DBaaSID() != "" {
		updateDBaaS(ctx, arubaClient, resources.DBaaS)
	}

	if resources.KaaS != nil && resources.KaaS.KaaSID() != "" {
		updateKaaS(ctx, arubaClient, resources.KaaS)
	}

	fmt.Println("\n=== Update Complete ===")
}
