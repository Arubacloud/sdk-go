package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createRestore provisions a restore of the given backup into the target block storage and waits until Ready.
func createRestore(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup, target *aruba.BlockStorage) *aruba.StorageRestore {
	fmt.Println("--- Storage Restore ---")

	if err := waitForDependencies(ctx, "Storage Restore", map[string]depEntry{
		"Backup":               dep(b, b.WaitUntilReady),
		"Target Block Storage": dep(target, target.WaitUntilNotUsed),
	}); err != nil {
		printDepWaitError("Storage Restore", err)
		return nil
	}

	r := aruba.NewStorageRestore().
		FromBackup(b).
		Named(resourceName(NameStorageRestore)).
		InRegion(aruba.RegionITBGBergamo).
		ToVolume(target)

	r, err := arubaClient.FromStorage().Restores().Create(ctx, r)
	if err != nil {
		printCreateError("Storage Restore", err)
		return nil
	}
	printCreated("Storage Restore", r.Name(), r.RestoreID())

	waitUntilSelfReady(ctx, "Storage Restore", r.Name(), r, r.WaitUntilReady)

	return r
}

// deleteRestore tears down the restore resource and waits until it is fully gone.
func deleteRestore(ctx context.Context, arubaClient aruba.Client, r *aruba.StorageRestore) {
	printDeleteBanner("Storage Restore")
	if err := arubaClient.FromStorage().Restores().Delete(ctx, r); err != nil {
		printDeleteError("Storage Restore", err)
		return
	}
	printDeleteSubmitted("Storage Restore", r.Name())
	waitUntilGone(ctx, "Storage Restore "+r.Name(), r.WaitUntilGone)
}
