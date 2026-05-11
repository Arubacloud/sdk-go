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
		log.Printf("%v", err)
		return nil
	}

	r := aruba.NewContainerRegistry().
		IntoProject(resources.Project).
		InRegion(aruba.RegionITBGBergamo).
		WithName(resourceName(NameContainerRegistry)).
		WithVPC(resources.VPC).
		WithSubnet(resources.SubnetBasic).
		WithSecurityGroup(resources.SecurityGroup).
		WithElasticIP(resources.ContainerRegistryEIP).
		WithBlockStorage(resources.ContainerRegistryStorage).
		WithBillingPeriod(aruba.BillingPeriodHour).
		WithAdminUsername("adminuser").
		OfSize(aruba.ContainerRegistrySizeFlavorSmall)

	resp, err := arubaClient.FromContainer().ContainerRegistry().Create(ctx, r)
	if err != nil {
		log.Fatalf("Error creating container registry: %s", formatErr(err))
		return nil
	}
	fmt.Printf("✓ Created container registry: %s\n", resp.Name())

	if err := resp.WaitUntilReady(ctx, longWaitOpts...); err != nil {
		log.Printf("Container Registry %s did not become Ready: %v", resp.Name(), err)
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
