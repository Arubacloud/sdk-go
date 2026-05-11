package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createElasticIP creates an Elastic IP with the given name
func createElasticIP(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, name string) *aruba.ElasticIP {
	fmt.Printf("--- Elastic IP (%s) ---\n", name)

	eip := aruba.NewElasticIP().
		IntoProject(proj).
		WithName(name).
		AddTag("network").
		AddTag("public").
		InRegion("ITBG-Bergamo").
		WithBillingPeriod("Hour")

	created, err := arubaClient.FromNetwork().ElasticIPs().Create(ctx, eip)
	if err != nil {
		log.Printf("Error creating Elastic IP: %v", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created Elastic IP: %s (ObjectID: %s)\n", created.Name(), created.ID())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("Elastic IP %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

// deleteElasticIP deletes an Elastic IP and waits for the platform to confirm
// removal. Project deletion fails with 400 if an Elastic IP is still in
// Deleting state.
func deleteElasticIP(ctx context.Context, arubaClient aruba.Client, eip *aruba.ElasticIP) {
	fmt.Println("--- Deleting Elastic IP ---")

	if err := arubaClient.FromNetwork().ElasticIPs().Delete(ctx, eip); err != nil {
		log.Printf("Error deleting Elastic IP: %v", err)
		return
	}
	fmt.Printf("✓ Deleted Elastic IP: %s\n", eip.ID())
	waitUntilGone(ctx, "Elastic IP "+eip.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().ElasticIPs().Get(ctx, eip)
		return err
	})
}
