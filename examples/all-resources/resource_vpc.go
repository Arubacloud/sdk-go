package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createVPC provisions a VPC and waits until Ready.
func createVPC(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref) *aruba.VPC {
	fmt.Println("--- VPC ---")

	vpc := aruba.NewVPC().
		IntoProject(proj).
		Named(resourceName(NameVPC)).
		AddTag("network").AddTag("infrastructure").
		InRegion(aruba.RegionITBGBergamo).
		WithPreset(false)

	created, err := arubaClient.FromNetwork().VPCs().Create(ctx, vpc)
	if err != nil {
		printCreateError("VPC", err)
		return nil
	}
	printCreated("VPC", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "VPC", created.Name(), created, created.WaitUntilReady)

	return created
}

// deleteVPC tears down the VPC and waits until gone.
func deleteVPC(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) {
	printDeleteBanner("VPC")
	if err := arubaClient.FromNetwork().VPCs().Delete(ctx, vpc); err != nil {
		printDeleteError("VPC", err)
		return
	}
	printDeleteSubmitted("VPC", vpc.Name())
	waitUntilGone(ctx, "VPC "+vpc.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().VPCs().Get(ctx, vpc)
		return err
	})
}
