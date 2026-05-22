package main

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createAdvancedSubnet provisions an advanced subnet with DHCP in the VPC and waits until Ready.
func createAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	printBanner("Subnet", "Advanced")

	if err := waitForDependencies(ctx, "Subnet (Advanced)", map[string]depEntry{
		"VPC": dep(vpc, vpc.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Subnet (Advanced)", err)
		return nil
	}

	subnet := aruba.NewSubnet().
		IntoVPC(vpc).
		Named(resourceName(NameSubnetAdvanced)).
		AddTag("network").
		AddTag("subnet").
		InRegion(aruba.RegionITBGBergamo).
		OfType(aruba.SubnetTypeAdvanced).
		WithCIDR("10.0.1.0/24").
		WithDHCP(aruba.NewSubnetDHCP().
			Enabled().
			WithRange("10.0.1.10", 100).
			AddDNS("8.8.8.8").
			AddDNS("8.8.4.4"))

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		printCreateError("Subnet (Advanced)", err)
		return nil
	}
	printCreated("Subnet (Advanced)", result.Name(), result.ID())

	waitUntilSelfReady(ctx, "Subnet (Advanced)", result.Name(), result, result.WaitUntilReady)

	return result
}

// createBasicSubnet provisions a basic subnet in the VPC and waits until Ready.
func createBasicSubnet(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.Subnet {
	printBanner("Subnet", "Basic")

	if err := waitForDependencies(ctx, "Subnet (Basic)", map[string]depEntry{
		"VPC": dep(vpc, vpc.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Subnet (Basic)", err)
		return nil
	}

	subnet := aruba.NewSubnet().
		IntoVPC(vpc).
		Named(resourceName(NameSubnetBasic)).
		AddTag("network").
		AddTag("subnet").
		InRegion(aruba.RegionITBGBergamo).
		OfType(aruba.SubnetTypeBasic)

	result, err := arubaClient.FromNetwork().Subnets().Create(ctx, subnet)
	if err != nil {
		printCreateError("Subnet (Basic)", err)
		return nil
	}
	printCreated("Subnet (Basic)", result.Name(), result.ID())

	waitUntilSelfReady(ctx, "Subnet (Basic)", result.Name(), result, result.WaitUntilReady)

	return result
}

// deleteAdvancedSubnet tears down the advanced subnet and waits until gone.
func deleteAdvancedSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	printDeleteBanner("Subnet (Advanced)")
	if err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet); err != nil {
		printDeleteError("Subnet (Advanced)", err)
		return
	}
	printDeleteSubmitted("Subnet (Advanced)", subnet.Name())
	waitUntilGone(ctx, "Subnet (Advanced) "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}

// deleteBasicSubnet tears down the basic subnet and waits until gone.
func deleteBasicSubnet(ctx context.Context, arubaClient aruba.Client, subnet *aruba.Subnet) {
	printDeleteBanner("Subnet (Basic)")
	if err := arubaClient.FromNetwork().Subnets().Delete(ctx, subnet); err != nil {
		printDeleteError("Subnet (Basic)", err)
		return
	}
	printDeleteSubmitted("Subnet (Basic)", subnet.Name())
	waitUntilGone(ctx, "Subnet (Basic) "+subnet.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().Subnets().Get(ctx, subnet)
		return err
	})
}
