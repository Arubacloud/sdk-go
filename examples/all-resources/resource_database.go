package main

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createDatabase provisions a database inside the given DBaaS instance.
func createDatabase(ctx context.Context, arubaClient aruba.Client, dbaas *aruba.DBaaS) *aruba.Database {
	printBanner("DBaaS Database", "")

	if err := waitForDependencies(ctx, "DBaaS Database", map[string]depEntry{
		"DBaaS": dep(dbaas, dbaas.WaitUntilReady),
	}); err != nil {
		printDepWaitError("DBaaS Database", err)
		return nil
	}

	// MySQL identifier rules forbid hyphens, so the hyphenated resourceName()
	// helper is bypassed here. A database name only needs to be unique within
	// its DBaaS instance, and each example run creates a fresh DBaaS.
	db := aruba.NewDatabase().
		IntoDBaaS(dbaas).
		Named(NameDatabase)

	res, err := arubaClient.FromDatabase().Databases().Create(ctx, db)
	if err != nil {
		printCreateError("DBaaS Database", err)
		return nil
	}
	printCreated("DBaaS Database", res.Name(), res.ID())
	return res
}

// deleteDatabase removes the database and waits until it is fully gone.
func deleteDatabase(ctx context.Context, arubaClient aruba.Client, db *aruba.Database) {
	printDeleteBanner("DBaaS Database")
	if err := arubaClient.FromDatabase().Databases().Delete(ctx, db); err != nil {
		printDeleteError("DBaaS Database", err)
		return
	}
	printDeleteSubmitted("DBaaS Database", db.Name())
	waitUntilGone(ctx, "DBaaS Database "+db.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromDatabase().Databases().Get(ctx, db)
		return err
	})
}
