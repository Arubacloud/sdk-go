package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createGrant grants the given user a role on the specified database.
func createGrant(ctx context.Context, arubaClient aruba.Client, db *aruba.Database, user *aruba.User) *aruba.Grant {
	fmt.Println("--- DBaaS (Grant) ---")

	g := aruba.NewGrant().
		IntoDatabase(db).
		WithUsername(user.Username()).
		WithRoleName("liteadmin")

	res, err := arubaClient.FromDatabase().Grants().Create(ctx, g)
	if err != nil {
		log.Fatalf("Error creating Grant: %s", formatErr(err))
		return nil
	}
	fmt.Printf("✓ Created Grant: %s on %s (%s)\n", res.Username(), res.DatabaseName(), res.RoleName())
	return res
}

// deleteGrant revokes the grant.
func deleteGrant(ctx context.Context, arubaClient aruba.Client, g *aruba.Grant) {
	fmt.Println("--- Deleting Grant ---")
	if err := arubaClient.FromDatabase().Grants().Delete(ctx, g); err != nil {
		log.Printf("Error deleting Grant: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted Grant: %s\n", g.ID())
}
