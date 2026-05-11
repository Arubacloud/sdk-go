package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createSnapshot creates a snapshot from block storage
func createSnapshot(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, bs *aruba.BlockStorage) *aruba.Snapshot {
	fmt.Println("--- Snapshot ---")

	snap := aruba.NewSnapshot().
		IntoProject(proj).
		WithName(resourceName("snapshot")).
		AddTag("backup").
		AddTag("snapshot").
		InRegion("ITBG-Bergamo").
		WithBillingPeriod("Hour").
		FromVolume(bs)

	snap, err := arubaClient.FromStorage().Snapshots().Create(ctx, snap)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to create snapshot - Status: %d, Error: %s",
				httpErr.StatusCode,
				stringValue(httpErr.ErrResp.Title))
		} else {
			log.Printf("Error creating snapshot: %v", err)
		}
		os.Exit(1)
	}
	fmt.Printf("✓ Created snapshot: %s from volume %s\n", snap.Name(), bs.Name())

	return snap
}

// deleteSnapshot deletes a snapshot
func deleteSnapshot(ctx context.Context, arubaClient aruba.Client, snap *aruba.Snapshot) {
	fmt.Println("--- Deleting Snapshot ---")

	err := arubaClient.FromStorage().Snapshots().Delete(ctx, snap)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete snapshot - Status: %d, Error: %s",
				httpErr.StatusCode,
				stringValue(httpErr.ErrResp.Title))
		} else {
			log.Printf("Error deleting snapshot: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted snapshot: %s\n", snap.ID())
}
