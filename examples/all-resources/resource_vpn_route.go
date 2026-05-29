package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createVPNRoute creates a VPN route under the given tunnel, routing the Basic Subnet's
// server-assigned CIDR to the on-premises subnet 192.168.1.0/24. The Basic Subnet must
// reach Active so its CIDR is available before the route is created.
func createVPNRoute(ctx context.Context, arubaClient aruba.Client, tunnel *aruba.VPNTunnel, subnet *aruba.Subnet) *aruba.VPNRoute {
	printBanner("VPN Route", "")

	if err := waitForDependencies(ctx, "VPN Route", map[string]depEntry{
		"VPN Tunnel":     dep(tunnel, tunnel.WaitUntilActive),
		"Subnet (Basic)": dep(subnet, subnet.WaitUntilActive),
	}); err != nil {
		printDepWaitError("VPN Route", err)
		return nil
	}

	name := resourceName(NameVPNRoute)
	route := aruba.NewVPNRoute().
		Named(name).
		Tagged("vpn-route", "route").
		InVPNTunnel(tunnel).
		InRegion(aruba.RegionITBGBergamo).
		WithCloudSubnet(subnet.CIDR()).
		WithOnPremSubnet("192.168.1.0/24")

	created, err := arubaClient.FromNetwork().VPNRoutes().Create(ctx, route)
	if err != nil {
		printCreateError("VPN Route", err)
		return nil
	}
	printCreated("VPN Route", created.Name(), created.ID(),
		"cloudSubnet="+created.CloudSubnetCIDR(),
		"onPremSubnet="+created.OnPremSubnet())

	waitUntilSelfReady(ctx, "VPN Route", created.Name(), created, created.WaitUntilReady)

	// GET the route back and print the raw JSON — used to inspect the CloudSubnet
	// wire shape that the API returns on GET (for issue #308 diagnosis).
	got, err := arubaClient.FromNetwork().VPNRoutes().Get(ctx, created)
	if err != nil {
		fmt.Printf("⚠ VPN Route GET failed (non-fatal): %v\n", err)
	} else {
		fmt.Printf("VPN Route GET CloudSubnetCIDR: %q\n", got.CloudSubnetCIDR())
		fmt.Printf("VPN Route GET RawJSON:\n%s\n", string(got.RawJSON()))
	}

	return created
}

// deleteVPNRoute deletes a VPN route and waits for it to be fully removed.
func deleteVPNRoute(ctx context.Context, arubaClient aruba.Client, route *aruba.VPNRoute) {
	printDeleteBanner("VPN Route")
	if err := arubaClient.FromNetwork().VPNRoutes().Delete(ctx, route); err != nil {
		printDeleteError("VPN Route", err)
		return
	}
	printDeleteSubmitted("VPN Route", route.Name())
	waitUntilGone(ctx, "VPN Route "+route.Name(), route.WaitUntilGone)
}
