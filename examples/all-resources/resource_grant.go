package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createGrant(ctx context.Context, arubaClient aruba.Client, db *aruba.Database, user *aruba.User) *aruba.Grant {
	fmt.Println("--- DBaaS: Grant ---")

	g := aruba.NewGrant().
		IntoDatabase(db).
		WithUsername(user.Username()).
		WithRoleName("liteadmin")

	res, err := arubaClient.FromDatabase().Grants().Create(ctx, g)
	if err != nil {
		log.Printf("Error creating Grant: %v", err)
		return nil
	}
	fmt.Printf("✓ Created Grant: %s on %s (%s)\n", res.Username(), res.DatabaseName(), res.RoleName())
	return res
}

func deleteGrant(ctx context.Context, arubaClient aruba.Client, g *aruba.Grant) {
	fmt.Println("--- Deleting Grant ---")
	if err := arubaClient.FromDatabase().Grants().Delete(ctx, g); err != nil {
		log.Printf("Error deleting Grant: %v", err)
		return
	}
	fmt.Printf("✓ Deleted Grant: %s\n", g.ID())
}
