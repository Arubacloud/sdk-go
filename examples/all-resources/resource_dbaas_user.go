package main

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createDBaaSUser creates a database user inside the given DBaaS instance.
func createDBaaSUser(ctx context.Context, arubaClient aruba.Client, dbaas *aruba.DBaaS) *aruba.User {
	printBanner("DBaaS User", "")

	if err := waitForDependencies(ctx, "DBaaS User", map[string]depEntry{
		"DBaaS": dep(dbaas, dbaas.WaitUntilReady),
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

// deleteDBaaSUser removes the database user and waits until it is fully gone.
func deleteDBaaSUser(ctx context.Context, arubaClient aruba.Client, u *aruba.User) {
	printDeleteBanner("DBaaS User")
	if err := arubaClient.FromDatabase().Users().Delete(ctx, u); err != nil {
		printDeleteError("DBaaS User", err)
		return
	}
	printDeleteSubmitted("DBaaS User", u.Username())
	waitUntilGone(ctx, "DBaaS User "+u.Username(), func(ctx context.Context) error {
		_, err := arubaClient.FromDatabase().Users().Get(ctx, u)
		return err
	})
}
