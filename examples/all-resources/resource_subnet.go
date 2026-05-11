package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	fmt.Println("\n--- Network: Subnet (Advanced) ---")

	if err := waitForDependencies(ctx, "Advanced Subnet", map[string]waitFunc{
		"VPC": vpc.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	subnet := aruba.NewSubnet().
		IntoVPC(vpc).
		WithName(resourceName(NameSubnetAdvanced)).
		AddTag("network").
		AddTag("subnet").
		InRegion("ITBG-Bergamo").
		OfType(aruba.SubnetTypeAdvanced).
		WithCIDR("10.0.1.0/24").
		WithDHCP(aruba.NewSubnetDHCP().
			Enabled().
			WithRange("10.0.1.10", 100).
			AddDNS("8.8.8.8").
			AddDNS("8.8.4.4"))

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		log.Printf("Error creating advanced subnet: %v", err)
		return result
	}
	fmt.Printf("✓ Created Advanced Subnet: %s (Type: %s, Network: %s)\n",
		result.Name(), result.Type(), result.CIDR())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Advanced Subnet %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

func createBasicSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	fmt.Println("\n--- Network: Subnet (Basic) ---")

	if err := waitForDependencies(ctx, "Basic Subnet", map[string]waitFunc{
		"VPC": vpc.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	subnet := aruba.NewSubnet().
		IntoVPC(vpc).
		WithName(resourceName(NameSubnetBasic)).
		AddTag("network").
		AddTag("subnet").
		InRegion("ITBG-Bergamo").
		OfType(aruba.SubnetTypeBasic)

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		log.Printf("Error creating basic subnet: %v", err)
		return result
	}
	fmt.Printf("✓ Created Basic Subnet: %s (Type: %s, Network: %s)\n",
		result.Name(), result.Type(), result.CIDR())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Basic Subnet %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

func deleteAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	fmt.Println("--- Deleting Advanced Subnet ---")

	err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet)
	if err != nil {
		log.Printf("Error deleting advanced subnet: %v", err)
		return
	}
	fmt.Printf("✓ Deleted advanced subnet: %s\n", subnet.ID())
	waitUntilGone(ctx, "advanced subnet "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}

func deleteBasicSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	fmt.Println("--- Deleting Basic Subnet ---")

	err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet)
	if err != nil {
		log.Printf("Error deleting basic subnet: %v", err)
		return
	}
	fmt.Printf("✓ Deleted basic subnet: %s\n", subnet.ID())
	waitUntilGone(ctx, "basic subnet "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}
