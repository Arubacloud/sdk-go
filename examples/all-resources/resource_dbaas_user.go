package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createDBaaSUser creates a database user inside the given DBaaS instance.
func createDBaaSUser(ctx context.Context, arubaClient aruba.Client, dbaas *aruba.DBaaS) *aruba.User {
	printBanner("DBaaS User", "")

	if err := waitForDependencies(ctx, "DBaaS User", map[string]waitFunc{
		"DBaaS": dbaas.WaitUntilReady,
	}); err != nil {
		printDepWaitError("DBaaS User", err)
		return nil
	}

	u := aruba.NewUser().
		IntoDBaaS(dbaas).
		WithUsername(NameDBaaSUser).
		WithPassword("Prova123456789AC@")

	res, err := arubaClient.FromDatabase().Users().Create(ctx, u)
	if err != nil {
		printCreateError("DBaaS User", err)
		return nil
	}
	printCreated("DBaaS User", res.Username(), res.ID())
	return res
}

// deleteDBaaSUser removes the database user from its DBaaS instance.
func deleteDBaaSUser(ctx context.Context, arubaClient aruba.Client, u *aruba.User) {
	fmt.Println("--- Deleting DBaaS User ---")
	if err := arubaClient.FromDatabase().Users().Delete(ctx, u); err != nil {
		log.Printf("Error deleting DBaaS User: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted DBaaS User: %s\n", u.Username())
}
