package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createCloudServer provisions a cloud server with all dependencies and waits until Ready.
func createCloudServer(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) *aruba.CloudServer {
	fmt.Println("--- Cloud Server ---")

	if err := waitForDependencies(ctx, "Cloud Server", map[string]waitFunc{
		"VPC":            resources.VPC.WaitUntilActive,
		"Subnet":         resources.SubnetBasic.WaitUntilActive,
		"Security Group": resources.SecurityGroup.WaitUntilActive,
		"Elastic IP":     resources.CloudServerEIP.WaitUntilNotUsed,
		"Block Storage":  resources.CloudServerBlockStorage.WaitUntilNotUsed,
		"Key Pair":       resources.KeyPair.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
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
		IntoProject(resources.Project).
		WithName(resourceName(NameCloudServer)).
		AddTag("virtualmachine").
		AddTag("container").
		InRegion(defaultRegion).
		InZone(defaultZone).
		OfFlavor("CSO4A8").
		WithVPC(resources.VPC).
		WithElasticIP(resources.CloudServerEIP).
		WithBootVolume(resources.CloudServerBlockStorage).
		WithKeyPair(resources.KeyPair).
		AddSubnet(resources.SubnetBasic).
		AddSecurityGroup(resources.SecurityGroup).
		WithUserData(userData)

	cs, err := arubaClient.FromCompute().CloudServers().Create(ctx, cs)
	if err != nil {
		log.Fatalf("Error creating Cloud Server: %s", formatErr(err))
		return nil
	}

	fmt.Printf("✓ Created Cloud Server: %s (Zone: %s, Flavor: %s)\n",
		cs.Name(),
		cs.Zone(),
		cs.Flavor())

	if err := cs.WaitUntilReady(ctx); err != nil {
		log.Printf("Cloud Server %s did not become Ready: %v", cs.Name(), err)
	}

	waitPostDependencies(ctx, "Cloud Server", map[string]waitFunc{
		"Elastic IP":    resources.CloudServerEIP.WaitUntilUsed,
		"Block Storage": resources.CloudServerBlockStorage.WaitUntilUsed,
	})

	return cs
}

// deleteCloudServer deletes a cloud server and blocks until the platform confirms
// it is fully gone. Cloud Server deletion is async: the HTTP call returns quickly
// but the server keeps running (and holding references to SG, Subnet, Block Storage,
// Elastic IP) until the platform completes teardown. Calling waitUntilGone here
// prevents the subsequent deletes from racing against that async termination.
func deleteCloudServer(ctx context.Context, arubaClient aruba.Client, cs *aruba.CloudServer) {
	fmt.Println("--- Deleting Cloud Server ---")

	if err := arubaClient.FromCompute().CloudServers().Delete(ctx, cs); err != nil {
		log.Printf("Error deleting cloud server: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted cloud server: %s\n", cs.Name())
	waitUntilGone(ctx, "cloud server "+cs.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromCompute().CloudServers().Get(ctx, cs)
		return err
	})
}
