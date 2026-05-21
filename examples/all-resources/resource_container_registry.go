package main

import (
	"context"
	"fmt"

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
		Named(resourceName(NameContainerRegistry)).
		InRegion(aruba.RegionITBGBergamo).
		OfSize(aruba.ContainerRegistrySizeFlavorSmall).
		WithAdminUsername("adminuser").
		WithBillingPeriod(aruba.BillingPeriodMonth).
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

	waitUntilSelfReady(ctx, "Container Registry", resp.Name(), resp.WaitUntilReady, longWaitOpts...)

	waitPostDependencies(ctx, "Container Registry", map[string]waitFunc{
		"Elastic IP":    resources.ContainerRegistryEIP.WaitUntilUsed,
		"Block Storage": resources.ContainerRegistryStorage.WaitUntilUsed,
	})

	return resp
}

// deleteContainerRegistry tears down the container registry and waits until it is
// fully gone. BS and EIP detachment races against an in-flight CR teardown without
// this wait, causing the API to reject those deletes.
func deleteContainerRegistry(ctx context.Context, arubaClient aruba.Client, r *aruba.ContainerRegistry) {
	printDeleteBanner("Container Registry")
	if err := arubaClient.FromContainer().ContainerRegistry().Delete(ctx, r); err != nil {
		printDeleteError("Container Registry", err)
		return
	}
	printDeleteSubmitted("Container Registry", r.Name())
	waitUntilGone(ctx, "Container Registry "+r.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromContainer().ContainerRegistry().Get(ctx, r)
		return err
	})
}
