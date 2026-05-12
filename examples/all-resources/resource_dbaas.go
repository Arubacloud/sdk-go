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

	if err := waitForDependencies(ctx, "DBaaS", map[string]waitFunc{
		"VPC":            vpc.WaitUntilActive,
		"Subnet":         subnet.WaitUntilActive,
		"Security Group": sg.WaitUntilActive,
		"Elastic IP":     eip.WaitUntilNotUsed,
	}); err != nil {
		printDepWaitError("DBaaS", err)
		return nil
	}

	d := aruba.NewDBaaS().
		IntoProject(proj).
		WithName(resourceName(NameDBaaS)).
		AddTag("database").
		AddTag("mysql").
		InRegion(aruba.RegionITBGBergamo).
		InZone(aruba.ZoneITBG1).
		OfEngine(aruba.DatabaseEngineMySQL80).
		OfFlavor(aruba.DBaaSFlavorDBO2A8).
		WithSizeGB(10).
		WithAutoscaling(2, 5).
		WithBillingPeriod(aruba.BillingPeriodHour).
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

	if err := result.WaitUntilReady(ctx, longWaitOpts...); err != nil {
		printSelfWaitError("DBaaS", result.Name(), err)
	}

	waitPostDependencies(ctx, "DBaaS", map[string]waitFunc{
		"Elastic IP": eip.WaitUntilUsed,
	})

	return result
}

// updateDBaaS applies name, tag, and storage-size changes to the DBaaS instance.
func updateDBaaS(ctx context.Context, arubaClient aruba.Client, d *aruba.DBaaS) {
	fmt.Println("--- Updating DBaaS ---")

	// Mutate what needs updating; networking URIs are already hydrated from the
	// prior Get call so they round-trip automatically into the Update request.
	d.WithName(updatedName(d.Name())).
		ReplaceTags("database", "mysql", "updated").
		WithSizeGB(25) // Increased from 20 to 25 GB

	result, err := arubaClient.FromDatabase().DBaaS().Update(ctx, d)
	if err != nil {
		log.Printf("Error updating DBaaS: %v", err)
		return
	}

	fmt.Printf("✓ Updated DBaaS: %s (Storage: %d GB)\n", result.Name(), result.SizeGB())
}

// deleteDBaaS tears down the DBaaS instance.
func deleteDBaaS(ctx context.Context, arubaClient aruba.Client, d *aruba.DBaaS) {
	fmt.Println("--- Deleting DBaaS ---")

	err := arubaClient.FromDatabase().DBaaS().Delete(ctx, d)
	if err != nil {
		log.Printf("Error deleting DBaaS: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted DBaaS: %s\n", d.Name())
}
