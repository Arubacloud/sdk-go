package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

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
		InRegion("ITBG-Bergamo").
		WithName(resourceName(NameContainerRegistry)).
		WithVPC(resources.VPC).
		WithSubnet(resources.SubnetBasic).
		WithSecurityGroup(resources.SecurityGroup).
		WithElasticIP(resources.ContainerRegistryEIP).
		WithBlockStorage(resources.ContainerRegistryStorage).
		WithBillingPeriod("Hour").
		WithAdminUsername("adminuser").
		OfSize(aruba.ContainerRegistrySizeFlavorSmall)

	resp, err := arubaClient.FromContainer().ContainerRegistry().Create(ctx, r)
	if err != nil {
		log.Printf("Error creating container registry: %v", err)
		return nil
	}
	fmt.Printf("✓ Created container registry: %s\n", resp.Name())

	if err := resp.WaitUntilReady(ctx, aruba.WithTimeout(20*time.Minute), aruba.WithRetries(120)); err != nil {
		log.Printf("Container Registry %s did not become Ready: %v", resp.Name(), err)
	}

	waitPostDependencies(ctx, "Container Registry", map[string]waitFunc{
		"Elastic IP":    resources.ContainerRegistryEIP.WaitUntilUsed,
		"Block Storage": resources.ContainerRegistryStorage.WaitUntilUsed,
	})

	return resp
}

func deleteContainerRegistry(ctx context.Context, arubaClient aruba.Client, r *aruba.ContainerRegistry) {
	fmt.Println("--- Deleting Container Registry ---")
	err := arubaClient.FromContainer().ContainerRegistry().Delete(ctx, r)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete container registry - Status: %d, Error: %s",
				httpErr.StatusCode,
				stringValue(httpErr.ErrResp.Title))
		} else {
			log.Printf("Error deleting container registry: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted Container Registry: %s\n", r.ContainerRegistryID())
}
