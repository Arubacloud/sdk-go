package main

import (
	"context"
	"fmt"
	"log"
	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createVPC provisions a VPC and waits until Ready.
func createVPC(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref) *aruba.VPC {
	fmt.Println("--- VPC ---")

	vpc := aruba.NewVPC().
		IntoProject(proj).
		WithName(resourceName(NameVPC)).
		AddTag("network").AddTag("infrastructure").
		InRegion(aruba.RegionITBGBergamo).
		WithPreset(false)

	created, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
	if err != nil {
		log.Fatalf("Error creating VPC: %s", formatErr(err))
	}

	fmt.Printf("✓ Created VPC: %s (Default: %t)\n", created.Name(), created.IsDefault())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("VPC %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

// deleteVPC tears down the VPC and waits until gone.
func deleteVPC(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) {
	fmt.Println("--- Deleting VPC ---")

	err := arubaClient.FromNetwork().VPCs().Delete(ctx, vpc)
	if err != nil {
		log.Printf("Error deleting VPC: %v", err)
		return
	}
	fmt.Printf("✓ Deleted VPC: %s\n", vpc.ID())
	waitUntilGone(ctx, "VPC "+vpc.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
		return err
	})
}
