package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createCloudServer provisions a cloud server with all dependencies and waits until Ready.
func createCloudServer(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) *aruba.CloudServer {
	fmt.Println("--- Cloud Server ---")

	if err := waitForDependencies(ctx, "Cloud Server", map[string]depEntry{
		"VPC":            dep(resources.VPC, resources.VPC.WaitUntilActive),
		"Subnet":         dep(resources.SubnetBasic, resources.SubnetBasic.WaitUntilActive),
		"Security Group": dep(resources.SecurityGroup, resources.SecurityGroup.WaitUntilActive),
		"Elastic IP":     dep(resources.CloudServerEIP, resources.CloudServerEIP.WaitUntilNotUsed),
		"Block Storage":  dep(resources.CloudServerBlockStorage, resources.CloudServerBlockStorage.WaitUntilNotUsed),
		"Key Pair":       dep(resources.KeyPair, resources.KeyPair.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Cloud Server", err)
		return nil
	}

	// Example cloud-init content: update packages and create a welcome file
	cloudInitContent := `#cloud-config
package_update: true
package_upgrade: true
write_files:
  - path: /etc/motd
    content: |
      Welcome to Aruba Cloud Server!
      This server was initialized with cloud-init.
    owner: root:root
    permissions: '0644'
`
	// Base64 encode the cloud-init content
	userData := base64.StdEncoding.EncodeToString([]byte(cloudInitContent))

	cs := aruba.NewCloudServer().
		OfFlavor(aruba.CloudServerFlavorCSO4A8).
		Named(resourceName(NameCloudServer)).
		Tagged("virtualmachine", "container").
		InProject(resources.Project).
		InRegion(aruba.RegionITBGBergamo).
		InZone(aruba.ZoneITBG1).
		WithUserData(userData).
		BootingFrom(resources.CloudServerBlockStorage).
		WithVPC(resources.VPC).
		OnSubnets(resources.SubnetBasic).
		WithSecurityGroups(resources.SecurityGroup).
		WithElasticIP(resources.CloudServerEIP).
		UsingKeyPair(resources.KeyPair).
		BilledBy(aruba.BillingPeriodHour)

	cs, err := arubaClient.FromCompute().CloudServers().Create(ctx, cs)
	if err != nil {
		printCreateError("Cloud Server", err)
		return nil
	}
	printCreated("Cloud Server", cs.Name(), cs.CloudServerID())

	waitUntilSelfReady(ctx, "Cloud Server", cs.Name(), cs, cs.WaitUntilReady)

	waitPostDependencies(ctx, "Cloud Server", map[string]depEntry{
		"Elastic IP":    dep(resources.CloudServerEIP, resources.CloudServerEIP.WaitUntilUsed),
		"Block Storage": dep(resources.CloudServerBlockStorage, resources.CloudServerBlockStorage.WaitUntilUsed),
	})

	return cs
}

// deleteCloudServer deletes a cloud server and blocks until the platform confirms
// it is fully gone. Cloud Server deletion is async: the HTTP call returns quickly
// but the server keeps running (and holding references to SG, Subnet, Block Storage,
// Elastic IP) until the platform completes teardown. Calling waitUntilGone here
// prevents the subsequent deletes from racing against that async termination.
func deleteCloudServer(ctx context.Context, arubaClient aruba.Client, cs *aruba.CloudServer) {
	printDeleteBanner("Cloud Server")
	if err := arubaClient.FromCompute().CloudServers().Delete(ctx, cs); err != nil {
		printDeleteError("Cloud Server", err)
		return
	}
	printDeleteSubmitted("Cloud Server", cs.Name())
	waitUntilGone(ctx, "Cloud Server "+cs.Name(), cs.WaitUntilGone)
}
