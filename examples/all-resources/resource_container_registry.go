package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createContainerRegistry provisions a container registry with all dependencies and waits until Ready.
func createContainerRegistry(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) *aruba.ContainerRegistry {
	fmt.Println("--- Container Registry ---")

	if err := waitForDependencies(ctx, "Container Registry", map[string]depEntry{
		"VPC":            dep(resources.VPC, resources.VPC.WaitUntilActive),
		"Subnet":         dep(resources.SubnetBasic, resources.SubnetBasic.WaitUntilActive),
		"Security Group": dep(resources.SecurityGroup, resources.SecurityGroup.WaitUntilActive),
		"Block Storage":  dep(resources.ContainerRegistryStorage, resources.ContainerRegistryStorage.WaitUntilNotUsed),
		"Elastic IP":     dep(resources.ContainerRegistryEIP, resources.ContainerRegistryEIP.WaitUntilNotUsed),
	}); err != nil {
		printDepWaitError("Container Registry", err)
		return nil
	}

	r := aruba.NewContainerRegistry().
		InProject(resources.Project).
		Named(resourceName(NameContainerRegistry)).
		InRegion(aruba.RegionITBGBergamo).
		OfSize(aruba.ContainerRegistrySizeFlavorSmall).
		WithAdminUsername("adminuser").
		BilledHourly().
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

	waitUntilSelfReady(ctx, "Container Registry", resp.Name(), resp, resp.WaitUntilReady, longWaitOpts...)

	waitPostDependencies(ctx, "Container Registry", map[string]depEntry{
		"Elastic IP":    dep(resources.ContainerRegistryEIP, resources.ContainerRegistryEIP.WaitUntilUsed),
		"Block Storage": dep(resources.ContainerRegistryStorage, resources.ContainerRegistryStorage.WaitUntilUsed),
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
	waitUntilGone(ctx, "Container Registry "+r.Name(), r.WaitUntilGone)
}
