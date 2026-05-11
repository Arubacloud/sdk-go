package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createDBaaSUser(ctx context.Context, arubaClient aruba.Client, dbaas *aruba.DBaaS) *aruba.User {
	fmt.Println("--- DBaaS: User ---")

	if err := waitForDependencies(ctx, "DBaaS User", map[string]waitFunc{
		"DBaaS": dbaas.WaitUntilReady,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	u := aruba.NewUser().
		IntoDBaaS(dbaas).
		WithUsername(NameDBaaSUser).
		WithPassword("Prova123456789AC@")

	res, err := arubaClient.FromDatabase().Users().Create(ctx, u)
	if err != nil {
		log.Printf("Error creating DBaaS User: %v", err)
		return nil
	}
	fmt.Printf("✓ Created DBaaS User: %s\n", res.Username())
	return res
}

func deleteDBaaSUser(ctx context.Context, arubaClient aruba.Client, u *aruba.User) {
	fmt.Println("--- Deleting DBaaS User ---")
	if err := arubaClient.FromDatabase().Users().Delete(ctx, u); err != nil {
		log.Printf("Error deleting DBaaS User: %v", err)
		return
	}
	fmt.Printf("✓ Deleted DBaaS User: %s\n", u.Username())
}
