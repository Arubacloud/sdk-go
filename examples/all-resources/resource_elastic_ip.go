package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createElasticIP creates an Elastic IP with the given name
func createElasticIP(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, name string) *aruba.ElasticIP {
	fmt.Printf("--- Elastic IP (%s) ---\n", name)

	eip := aruba.NewElasticIP().
		IntoProject(proj).
		Named(name).
		AddTag("network").
		AddTag("public").
		InRegion(aruba.RegionITBGBergamo).
		WithBillingPeriod(aruba.BillingPeriodHour)

	created, err := arubaClient.FromNetwork().ElasticIPs().Create(ctx, eip)
	if err != nil {
		printCreateError("Elastic IP", err)
		return nil
	}
	printCreated("Elastic IP", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "Elastic IP", created.Name(), created, created.WaitUntilReady)

	return created
}

// deleteElasticIP deletes an Elastic IP and waits for the platform to confirm
// removal. Project deletion fails with 400 if an Elastic IP is still in
// Deleting state.
func deleteElasticIP(ctx context.Context, arubaClient aruba.Client, eip *aruba.ElasticIP) {
	printDeleteBanner("Elastic IP")
	if err := arubaClient.FromNetwork().ElasticIPs().Delete(ctx, eip); err != nil {
		printDeleteError("Elastic IP", err)
		return
	}
	printDeleteSubmitted("Elastic IP", eip.Name())
	waitUntilGone(ctx, "Elastic IP "+eip.Name(), eip.WaitUntilGone)
}
