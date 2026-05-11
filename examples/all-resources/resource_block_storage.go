package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createBlockStorage creates a block storage volume with the given name
func createBlockStorage(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, name string) *aruba.BlockStorage {
	fmt.Printf("--- Block Storage (%s) ---\n", name)

	bs := aruba.NewBlockStorage().
		IntoProject(proj).
		WithName(name).
		AddTag("storage").
		AddTag("data").
		InRegion("ITBG-Bergamo").
		InZone("ITBG-1").
		WithSizeGB(20).
		OfType(aruba.BlockStorageTypeStandard).
		WithBillingPeriod("Hour").
		SetBootable().
		FromImage("LU22-001")

	bs, err := arubaClient.FromStorage().Volumes().Create(ctx, bs)
	if err != nil {
		log.Printf("Error creating block storage: %v", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created block storage: %s (%d GB, %s)\n", bs.Name(), bs.SizeGB(), bs.Type())

	if err := bs.WaitUntilReady(ctx); err != nil {
		log.Printf("Block Storage %s did not become Ready: %v", bs.Name(), err)
	}

	return bs
}

// deleteBlockStorage deletes a block storage volume and waits for the platform
// to confirm removal. Project deletion fails with 400 if a volume is still in
// Deleting state.
func deleteBlockStorage(ctx context.Context, arubaClient aruba.Client, bs *aruba.BlockStorage) {
	fmt.Println("--- Deleting Block Storage ---")

	if err := arubaClient.FromStorage().Volumes().Delete(ctx, bs); err != nil {
		log.Printf("Error deleting block storage: %v", err)
		return
	}
	fmt.Printf("✓ Deleted block storage: %s\n", bs.ID())
	waitUntilGone(ctx, "block storage "+bs.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromStorage().Volumes().Get(ctx, bs)
		return err
	})
}
