package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createRestore(ctx context.Context, arubaClient aruba.Client, b *aruba.StorageBackup, target *aruba.BlockStorage) *aruba.StorageRestore {
	fmt.Println("--- Restore ---")

	if err := waitForDependencies(ctx, "Restore", map[string]waitFunc{
		"Backup":               b.WaitUntilReady,
		"Target Block Storage": target.WaitUntilNotUsed,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	r := aruba.NewStorageRestore().
		IntoBackup(b).
		InRegion("ITBG-Bergamo").
		WithName(resourceName(NameStorageRestore)).
		ToVolume(target)

	r, err := arubaClient.FromStorage().Restores().Create(ctx, r)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to create restore - Status: %d, Body: %s", httpErr.StatusCode, string(httpErr.Body))
		} else {
			log.Printf("Error creating restore: %v", err)
		}
		return nil
	}
	fmt.Printf("✓ Created restore: %s\n", r.Name())

	if err := r.WaitUntilReady(ctx); err != nil {
		log.Printf("Storage Restore %s did not become Ready: %v", r.Name(), err)
	}

	return r
}

// deleteRestore deletes a restore resource
func deleteRestore(ctx context.Context, arubaClient aruba.Client, r *aruba.StorageRestore) {
	fmt.Println("--- Deleting Restore ---")
	err := arubaClient.FromStorage().Restores().Delete(ctx, r)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete restore - Status: %d, Body: %s", httpErr.StatusCode, string(httpErr.Body))
		} else {
			log.Printf("Error deleting restore: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted Restore: %s\n", r.Name())
}
