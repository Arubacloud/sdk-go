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
		log.Printf("%v", err)
		return nil
	}

	b := aruba.NewStorageBackup().
		IntoProject(proj).
		InRegion(aruba.RegionITBGBergamo).
		WithName(resourceName(NameStorageBackup)).
		OfType(aruba.StorageBackupTypeFull).
		FromVolume(bs).
		WithRetentionDays(10).
		WithBillingPeriod("Monthly") // no BillingPeriodMonthly constant — see pkg/aruba/aliases.go

	result, err := arubaClient.FromStorage().Backups().Create(ctx, b)
	if err != nil {
		log.Fatalf("Error creating storage backup: %s", formatErr(err))
	}
	fmt.Printf("✓ Created storage backup: %s\n", result.Name())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Storage Backup %s did not become Ready: %v", result.Name(), err)
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
