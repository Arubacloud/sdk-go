package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createRestore provisions a restore of the given backup into the target block storage and waits until Ready.
func createRestore(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup, target *aruba.BlockStorage) *aruba.StorageRestore {
	fmt.Println("--- Storage Restore ---")

	if err := waitForDependencies(ctx, "Restore", map[string]waitFunc{
		"Backup":               b.WaitUntilReady,
		"Target Block Storage": target.WaitUntilNotUsed,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	r := aruba.NewStorageRestore().
		IntoBackup(b).
		InRegion(aruba.RegionITBGBergamo).
		WithName(resourceName(NameStorageRestore)).
		ToVolume(target)

	r, err := arubaClient.FromStorage().Restores().Create(ctx, r)
	if err != nil {
		log.Fatalf("Error creating restore: %s", formatErr(err))
	}
	fmt.Printf("✓ Created restore: %s\n", r.Name())

	if err := r.WaitUntilReady(ctx); err != nil {
		log.Printf("Storage Restore %s did not become Ready: %v", r.Name(), err)
	}

	return r
}

// deleteRestore tears down the restore resource.
func deleteRestore(ctx context.Context, arubaClient aruba.Client, r *aruba.StorageRestore) {
	fmt.Println("--- Deleting Storage Restore ---")
	err := arubaClient.FromStorage().Restores().Delete(ctx, r)
	if err != nil {
		log.Printf("Error deleting restore: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted storage restore: %s\n", r.Name())
}
