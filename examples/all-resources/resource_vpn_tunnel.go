package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createVPNTunnel creates a Site-to-Site VPN tunnel with IKEv2/AES-256/SHA-256 settings.
// The peer IP and PSK use placeholder values suitable for a dry-run; replace them with
// real values when running against a live peer.
func createVPNTunnel(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, vpc *aruba.VPC, subnet *aruba.Subnet, eip *aruba.ElasticIP) *aruba.VPNTunnel {
	name := resourceName(NameVPNTunnel)
	fmt.Printf("--- VPN Tunnel (%s) ---\n", name)

	ipcfg := aruba.NewVPNIPConfig().
		WithVPC(vpc).
		WithSubnet(subnet.Name(), subnet.CIDR()).
		WithElasticIP(eip)

	tunnel := aruba.NewVPNTunnel().
		OfType(aruba.VPNTypeSiteToSite).
		Named(name).
		Tagged("vpn-tunnel", "network").
		InProject(proj).
		InRegion(aruba.RegionITBGBergamo).
		WithIPConfig(ipcfg).
		WithVPNClientProtocol(aruba.VPNClientProtocolIKEv2).
		WithPeerClientPublicIP("203.0.113.1"). // TEST-NET-3, replace with real peer IP
		WithIKESettings(aruba.NewVPNIKE().
			WithEncryption(aruba.IKEEncryptionAES256).
			WithHash(aruba.IKEHashSHA256).
			WithDHGroup(aruba.IKEDHGroup14).
			WithLifetimeSeconds(86400).
			WithDPDAction(aruba.IKEDPDActionRestart).
			WithDPDIntervalSeconds(30).
			WithDPDTimeoutSeconds(120)).
		WithESPSettings(aruba.NewVPNESP().
			WithEncryption(aruba.ESPEncryptionAES256).
			WithHash(aruba.ESPHashSHA256).
			WithPFS(aruba.ESPPFSGroupDHGroup14).
			WithLifetimeSeconds(3600)).
		WithPSKSettings(aruba.NewVPNPSK().
			WithCloudSite("cloud-side").
			WithOnPremSite("onprem-side").
			WithKey("example-psk-replace-in-production")).
		BilledBy(aruba.BillingPeriodHour)

	created, err := arubaClient.FromNetwork().VPNTunnels().Create(ctx, tunnel)
	if err != nil {
		printCreateError("VPN Tunnel", err)
		return nil
	}
	printCreated("VPN Tunnel", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "VPN Tunnel", created.Name(), created, created.WaitUntilReady)

	return created
}

// deleteVPNTunnel deletes a VPN Tunnel and waits for it to be fully removed.
func deleteVPNTunnel(ctx context.Context, arubaClient aruba.Client, tunnel *aruba.VPNTunnel) {
	printDeleteBanner("VPN Tunnel")
	if err := arubaClient.FromNetwork().VPNTunnels().Delete(ctx, tunnel); err != nil {
		printDeleteError("VPN Tunnel", err)
		return
	}
	printDeleteSubmitted("VPN Tunnel", tunnel.Name())
	waitUntilGone(ctx, "VPN Tunnel "+tunnel.Name(), tunnel.WaitUntilGone)
}
