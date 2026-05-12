package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createContainerRegistry provisions a container registry with all dependencies and waits until Ready.
func createContainerRegistry(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) *aruba.ContainerRegistry {
	fmt.Println("--- Container Registry ---")

	if err := waitForDependencies(ctx, "Container Registry", map[string]waitFunc{
		"VPC":            resources.VPC.WaitUntilActive,
		"Subnet":         resources.SubnetBasic.WaitUntilActive,
		"Security Group": resources.SecurityGroup.WaitUntilActive,
		"Block Storage":  resources.ContainerRegistryStorage.WaitUntilNotUsed,
		"Elastic IP":     resources.ContainerRegistryEIP.WaitUntilNotUsed,
	}); err != nil {
		printDepWaitError("Container Registry", err)
		return nil
	}

	r := aruba.NewContainerRegistry().
		IntoProject(resources.Project).
		WithName(resourceName(NameContainerRegistry)).
		InRegion(aruba.RegionITBGBergamo).
		OfSize(aruba.ContainerRegistrySizeFlavorSmall).
		WithAdminUsername("adminuser").
		WithBillingPeriod(aruba.BillingPeriodHour).
		WithVPC(resources.VPC).
		WithSubnet(resources.SubnetBasic).
		WithSecurityGroup(resources.SecurityGroup).
		WithElasticIP(resources.ContainerRegistryEIP).
		WithBlockStorage(resources.ContainerRegistryStorage)

	resp, err := arubaClient.FromContainer().ContainerRegistry().Create(ctx, r)
	if err != nil {
		printCreateError("Container Registry", err)
		return nil
	}
	printCreated("Container Registry", resp.Name(), resp.ContainerRegistryID())

	if err := resp.WaitUntilReady(ctx, longWaitOpts...); err != nil {
		printSelfWaitError("Container Registry", resp.Name(), err)
	}

	waitPostDependencies(ctx, "Container Registry", map[string]waitFunc{
		"Elastic IP":    resources.ContainerRegistryEIP.WaitUntilUsed,
		"Block Storage": resources.ContainerRegistryStorage.WaitUntilUsed,
	})

	return resp
}

// deleteContainerRegistry tears down the container registry.
func deleteContainerRegistry(ctx context.Context, arubaClient aruba.Client, r *aruba.ContainerRegistry) {
	fmt.Println("--- Deleting Container Registry ---")
	err := arubaClient.FromContainer().ContainerRegistry().Delete(ctx, r)
	if err != nil {
		log.Printf("Error deleting container registry: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted Container Registry: %s\n", r.ContainerRegistryID())
}
