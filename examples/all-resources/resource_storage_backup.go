package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createStorageBackup provisions a full backup of the given block storage volume and waits until Ready.
func createStorageBackup(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, bs *aruba.BlockStorage) *aruba.StorageBackup {
	fmt.Println("--- Storage Backup ---")

	if err := waitForDependencies(ctx, "Storage Backup", map[string]waitFunc{
		"Block Storage": bs.WaitUntilReady,
	}); err != nil {
		printDepWaitError("Storage Backup", err)
		return nil
	}

	b := aruba.NewStorageBackup().
		IntoProject(proj).
		WithName(resourceName(NameStorageBackup)).
		InRegion(aruba.RegionITBGBergamo).
		OfType(aruba.StorageBackupTypeFull).
		WithRetentionDays(10).
		WithBillingPeriod(aruba.BillingPeriodMonth).
		FromVolume(bs)

	result, err := arubaClient.FromStorage().Backups().Create(ctx, b)
	if err != nil {
		printCreateError("Storage Backup", err)
		return nil
	}
	printCreated("Storage Backup", result.Name(), result.BackupID())

	if err := result.WaitUntilReady(ctx); err != nil {
		printSelfWaitError("Storage Backup", result.Name(), err)
	}

	return result
}

// deleteStorageBackup tears down the storage backup.
func deleteStorageBackup(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup) {
	fmt.Println("--- Deleting Storage Backup ---")
	err := arubaClient.FromStorage().Backups().Delete(ctx, b)
	if err != nil {
		log.Printf("Error deleting storage backup: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted storage backup: %s\n", b.Name())
}
