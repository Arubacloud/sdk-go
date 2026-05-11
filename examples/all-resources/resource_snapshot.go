package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createSnapshot creates a volume snapshot from the given block storage and waits until Ready.
func createSnapshot(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, bs *aruba.BlockStorage) *aruba.Snapshot {
	fmt.Println("--- Snapshot ---")

	if err := waitForDependencies(ctx, "Snapshot", map[string]waitFunc{
		"Block Storage": bs.WaitUntilReady,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	snap := aruba.NewSnapshot().
		IntoProject(proj).
		WithName(resourceName(NameSnapshot)).
		AddTag("backup").
		AddTag("snapshot").
		InRegion(aruba.RegionITBGBergamo).
		WithBillingPeriod(aruba.BillingPeriodHour).
		FromVolume(bs)

	snap, err := arubaClient.FromStorage().Snapshots().Create(ctx, snap)
	if err != nil {
		log.Fatalf("Error creating snapshot: %s", formatErr(err))
	}
	fmt.Printf("✓ Created snapshot: %s from volume %s\n", snap.Name(), bs.Name())

	if err := snap.WaitUntilReady(ctx); err != nil {
		log.Printf("Snapshot %s did not become Ready: %v", snap.Name(), err)
	}

	return snap
}

// deleteSnapshot tears down the snapshot.
func deleteSnapshot(ctx context.Context, arubaClient aruba.Client, snap *aruba.Snapshot) {
	fmt.Println("--- Deleting Snapshot ---")

	err := arubaClient.FromStorage().Snapshots().Delete(ctx, snap)
	if err != nil {
		log.Printf("Error deleting snapshot: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted snapshot: %s\n", snap.ID())
}
