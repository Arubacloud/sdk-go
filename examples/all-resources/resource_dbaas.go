package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createDBaaS(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, vpc *aruba.VPC, subnet *aruba.Subnet, sg *aruba.SecurityGroup, eip *aruba.ElasticIP) *aruba.DBaaS {
	fmt.Println("--- DBaaS ---")

	if err := waitForDependencies(ctx, "DBaaS", map[string]waitFunc{
		"VPC":            vpc.WaitUntilActive,
		"Subnet":         subnet.WaitUntilActive,
		"Security Group": sg.WaitUntilActive,
		"Elastic IP":     eip.WaitUntilNotUsed,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	d := aruba.NewDBaaS().
		IntoProject(proj).
		WithName(resourceName(NameDBaaS)).
		AddTag("database").
		AddTag("mysql").
		InRegion("ITBG-Bergamo").
		InZone("ITBG-1").
		OfEngine("mysql-8.0").
		OfFlavor("DBO2A8").
		WithSizeGB(10).
		WithAutoscaling(2, 5).
		WithBillingPeriod("Hour").
		WithVPC(vpc).
		WithSubnet(subnet).
		WithSecurityGroup(sg).
		WithElasticIP(eip)
	if err := d.Err(); err != nil {
		log.Printf("Error building DBaaS request: %v", err)
		return nil
	}

	result, err := arubaClient.FromDatabase().DBaaS().Create(ctx, d)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to create DBaaS - Status: %d, Error: %s", httpErr.StatusCode, httpErr.Error())
		} else {
			log.Printf("Error creating DBaaS: %v", err)
		}
		return nil
	}

	fmt.Printf("✓ Created DBaaS: %s (Engine: %s, Flavor: %s, Storage: %d GB)\n",
		result.Name(), result.Engine(), result.Flavor(), result.SizeGB())

	if err := result.WaitUntilReady(ctx, aruba.WithTimeout(20*time.Minute), aruba.WithRetries(120)); err != nil {
		log.Printf("DBaaS %s did not become Ready: %v", result.Name(), err)
	}

	waitPostDependencies(ctx, "DBaaS", map[string]waitFunc{
		"Elastic IP": eip.WaitUntilUsed,
	})

	return result
}

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

func deleteDBaaS(ctx context.Context, arubaClient aruba.Client, d *aruba.DBaaS) {
	fmt.Println("--- Deleting DBaaS ---")

	err := arubaClient.FromDatabase().DBaaS().Delete(ctx, d)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete DBaaS - Status: %d, Error: %s", httpErr.StatusCode, httpErr.Error())
		} else {
			log.Printf("Error deleting DBaaS: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted DBaaS: %s\n", d.Name())
}
