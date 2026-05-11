package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createAdvancedSubnet provisions an advanced subnet with DHCP in the VPC and waits until Ready.
func createAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	fmt.Println("--- Network: Subnet (Advanced) ---")

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
		InRegion(defaultRegion).
		OfType(aruba.SubnetTypeAdvanced).
		WithCIDR("10.0.1.0/24").
		WithDHCP(aruba.NewSubnetDHCP().
			Enabled().
			WithRange("10.0.1.10", 100).
			AddDNS("8.8.8.8").
			AddDNS("8.8.4.4"))

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		log.Fatalf("Error creating advanced subnet: %s", formatErr(err))
		return nil
	}
	fmt.Printf("✓ Created Advanced Subnet: %s (Type: %s, Network: %s)\n",
		result.Name(), result.Type(), result.CIDR())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Advanced Subnet %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

// createBasicSubnet provisions a basic subnet in the VPC and waits until Ready.
func createBasicSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	fmt.Println("--- Network: Subnet (Basic) ---")

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
		InRegion(defaultRegion).
		OfType(aruba.SubnetTypeBasic)

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		log.Fatalf("Error creating basic subnet: %s", formatErr(err))
		return nil
	}
	fmt.Printf("✓ Created Basic Subnet: %s (Type: %s, Network: %s)\n",
		result.Name(), result.Type(), result.CIDR())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Basic Subnet %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

// deleteAdvancedSubnet tears down the advanced subnet and waits until gone.
func deleteAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	fmt.Println("--- Deleting Advanced Subnet ---")

	err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet)
	if err != nil {
		log.Printf("Error deleting advanced subnet: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted advanced subnet: %s\n", subnet.ID())
	waitUntilGone(ctx, "advanced subnet "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}

// deleteBasicSubnet tears down the basic subnet and waits until gone.
func deleteBasicSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	fmt.Println("--- Deleting Basic Subnet ---")

	err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet)
	if err != nil {
		log.Printf("Error deleting basic subnet: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted basic subnet: %s\n", subnet.ID())
	waitUntilGone(ctx, "basic subnet "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}
