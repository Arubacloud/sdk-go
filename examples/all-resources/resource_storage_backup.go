package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createStorageBackup provisions a full backup of the given block storage volume and waits until Ready.
func createStorageBackup(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, bs *aruba.BlockStorage) *aruba.StorageBackup {
	fmt.Println("--- Storage Backup ---")

	if err := waitForDependencies(ctx, "Storage Backup", map[string]depEntry{
		"Block Storage": dep(bs, bs.WaitUntilReady),
	}); err != nil {
		printDepWaitError("Storage Backup", err)
		return nil
	}

	b := aruba.NewStorageBackup().
		IntoProject(proj).
		Named(resourceName(NameStorageBackup)).
		InRegion(aruba.RegionITBGBergamo).
		OfType(aruba.StorageBackupTypeFull).
		WithRetentionDays(10).
		WithBillingPeriod(aruba.BillingPeriodHour).
		FromVolume(bs)

	result, err := arubaClient.FromStorage().Backups().Create(ctx, b)
	if err != nil {
		printCreateError("Storage Backup", err)
		return nil
	}
	printCreated("Storage Backup", result.Name(), result.BackupID())

	waitUntilSelfReady(ctx, "Storage Backup", result.Name(), result, result.WaitUntilReady)

	return result
}

// deleteStorageBackup tears down the storage backup and waits until it is fully gone.
func deleteStorageBackup(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup) {
	printDeleteBanner("Storage Backup")
	if err := arubaClient.FromStorage().Backups().Delete(ctx, b); err != nil {
		printDeleteError("Storage Backup", err)
		return
	}
	printDeleteSubmitted("Storage Backup", b.Name())
	waitUntilGone(ctx, "Storage Backup "+b.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromStorage().Backups().Get(ctx, b)
		return err
	})
}
