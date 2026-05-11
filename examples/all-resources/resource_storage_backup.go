package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

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
		InRegion("ITBG-Bergamo").
		WithName(resourceName(NameStorageBackup)).
		OfType(aruba.StorageBackupTypeFull).
		FromVolume(bs).
		WithRetentionDays(10).
		WithBillingPeriod("Monthly")

	result, err := arubaClient.FromStorage().Backups().Create(ctx, b)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to create storage backup - Status: %d, Error: %s",
				httpErr.StatusCode,
				httpErr.Error())
		} else {
			log.Printf("Error creating storage backup: %v", err)
		}
		os.Exit(1)
	}
	fmt.Printf("✓ Created storage backup: %s\n", result.Name())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("Storage Backup %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

func deleteBackup(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup) {
	fmt.Println("--- Deleting Backup ---")
	err := arubaClient.FromStorage().Backups().Delete(ctx, b)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete backup - Status: %d, Error: %s",
				httpErr.StatusCode,
				httpErr.Error())
		} else {
			log.Printf("Error deleting backup: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted Backup: %s\n", b.Name())
}
