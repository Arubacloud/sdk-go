package main

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createGrant grants the given user a role on the specified database.
func createGrant(ctx context.Context, arubaClient aruba.Client, db *aruba.Database, user *aruba.User) *aruba.Grant {
	printBanner("DBaaS Grant", "")

	g := aruba.NewGrant().
		OfRole("liteadmin").
		InDatabase(db).
		ForUser(user.Username())

	res, err := arubaClient.FromDatabase().Grants().Create(ctx, g)
	if err != nil {
		printCreateError("DBaaS Grant", err)
		return nil
	}
	printCreated("DBaaS Grant", res.Username(), res.ID())
	return res
}

// deleteGrant revokes the grant and waits until it is fully gone.
func deleteGrant(ctx context.Context, arubaClient aruba.Client, g *aruba.Grant) {
	printDeleteBanner("Grant")
	if err := arubaClient.FromDatabase().Grants().Delete(ctx, g); err != nil {
		printDeleteError("Grant", err)
		return
	}
	printDeleteSubmitted("Grant", g.ID())
	waitUntilGone(ctx, "Grant "+g.ID(), g.WaitUntilGone)
}
