package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createDBaaS provisions a DBaaS instance with all dependencies and waits until Ready.
func createDBaaS(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, vpc *aruba.VPC, subnet *aruba.Subnet, sg *aruba.SecurityGroup, eip *aruba.ElasticIP) *aruba.DBaaS {
	fmt.Println("--- DBaaS ---")

	if err := waitForDependencies(ctx, "DBaaS", map[string]depEntry{
		"VPC":            dep(vpc, vpc.WaitUntilActive),
		"Subnet":         dep(subnet, subnet.WaitUntilActive),
		"Security Group": dep(sg, sg.WaitUntilActive),
		"Elastic IP":     dep(eip, eip.WaitUntilNotUsed),
	}); err != nil {
		printDepWaitError("DBaaS", err)
		return nil
	}

	d := aruba.NewDBaaS().
		InProject(proj).
		Named(resourceName(NameDBaaS)).
		Tagged("database").
		Tagged("mysql").
		InRegion(aruba.RegionITBGBergamo).
		InZone(aruba.ZoneITBG1).
		OfEngine(aruba.DatabaseEngineMySQL80).
		OfFlavor(aruba.DBaaSFlavorDBO4A8).
		SizedGB(10).
		WithAutoscaling(2, 5).
		BilledHourly().
		WithVPC(vpc).
		WithSubnet(subnet).
		WithSecurityGroup(sg).
		WithElasticIP(eip)
	if err := d.Err(); err != nil {
		log.Printf("✗ Invalid DBaaS request: %v", err)
		return nil
	}

	result, err := arubaClient.FromDatabase().DBaaS().Create(ctx, d)
	if err != nil {
		printCreateError("DBaaS", err)
		return nil
	}
	printCreated("DBaaS", result.Name(), result.DBaaSID())

	waitUntilSelfReady(ctx, "DBaaS", result.Name(), result, result.WaitUntilReady, longWaitOpts...)

	waitPostDependencies(ctx, "DBaaS", map[string]depEntry{
		"Elastic IP": dep(eip, eip.WaitUntilUsed),
	})

	return result
}

// updateDBaaS applies name, tag, and storage-size changes to the DBaaS instance.
func updateDBaaS(ctx context.Context, arubaClient aruba.Client, d *aruba.DBaaS) {
	fmt.Println("--- Updating DBaaS ---")

	// Mutate what needs updating; networking URIs are already hydrated from the
	// prior Get call so they round-trip automatically into the Update request.
	d.Named(updatedName(d.Name())).
		RetaggedAs("database", "mysql", "updated").
		SizedGB(25) // Increased from 20 to 25 GB

	result, err := arubaClient.FromDatabase().DBaaS().Update(ctx, d)
	if err != nil {
		log.Printf("Error updating DBaaS: %v", err)
		return
	}

	fmt.Printf("✓ Updated DBaaS: %s (Storage: %d GB)\n", result.Name(), result.SizeGB())
}

// deleteDBaaS tears down the DBaaS instance and waits until it is fully gone.
func deleteDBaaS(ctx context.Context, arubaClient aruba.Client, d *aruba.DBaaS) {
	printDeleteBanner("DBaaS")
	if err := arubaClient.FromDatabase().DBaaS().Delete(ctx, d); err != nil {
		printDeleteError("DBaaS", err)
		return
	}
	printDeleteSubmitted("DBaaS", d.Name())
	waitUntilGone(ctx, "DBaaS "+d.Name(), d.WaitUntilGone)
}
