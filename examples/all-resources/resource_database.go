package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createDatabase(ctx context.Context, arubaClient aruba.Client, dbaas *aruba.DBaaS) *aruba.Database {
	fmt.Println("--- DBaaS: Database ---")

	if err := waitForDependencies(ctx, "Database", map[string]waitFunc{
		"DBaaS": dbaas.WaitUntilReady,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	// MySQL identifier rules forbid hyphens, so the hyphenated resourceName()
	// helper is bypassed here. A database name only needs to be unique within
	// its DBaaS instance, and each example run creates a fresh DBaaS.
	db := aruba.NewDatabase().
		IntoDBaaS(dbaas).
		WithName(NameDatabase)

	res, err := arubaClient.FromDatabase().Databases().Create(ctx, db)
	if err != nil {
		log.Printf("Error creating Database: %v", err)
		return nil
	}
	fmt.Printf("✓ Created Database: %s\n", res.Name())
	return res
}

func deleteDatabase(ctx context.Context, arubaClient aruba.Client, db *aruba.Database) {
	fmt.Println("--- Deleting Database ---")
	if err := arubaClient.FromDatabase().Databases().Delete(ctx, db); err != nil {
		log.Printf("Error deleting Database: %v", err)
		return
	}
	fmt.Printf("✓ Deleted Database: %s\n", db.Name())
}
