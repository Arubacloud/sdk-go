package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createSnapshot creates a volume snapshot from the given block storage and waits until Ready.
func createSnapshot(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, bs *aruba.BlockStorage) *aruba.Snapshot {
	fmt.Println("--- Snapshot ---")

	if err := waitForDependencies(ctx, "Snapshot", map[string]waitFunc{
		"Block Storage": bs.WaitUntilReady,
	}); err != nil {
		printDepWaitError("Snapshot", err)
		return nil
	}

	snap := aruba.NewSnapshot().
		IntoProject(proj).
		Named(resourceName(NameSnapshot)).
		AddTag("backup").
		AddTag("snapshot").
		InRegion(aruba.RegionITBGBergamo).
		WithBillingPeriod(aruba.BillingPeriodHour).
		FromVolume(bs)

	snap, err := arubaClient.FromStorage().Snapshots().Create(ctx, snap)
	if err != nil {
		printCreateError("Snapshot", err)
		return nil
	}
	printCreated("Snapshot", snap.Name(), snap.ID())

	waitUntilSelfReady(ctx, "Snapshot", snap.Name(), snap.WaitUntilReady)

	return snap
}

// deleteSnapshot tears down the snapshot and waits until it is fully gone.
func deleteSnapshot(ctx context.Context, arubaClient aruba.Client, snap *aruba.Snapshot) {
	printDeleteBanner("Snapshot")
	if err := arubaClient.FromStorage().Snapshots().Delete(ctx, snap); err != nil {
		printDeleteError("Snapshot", err)
		return
	}
	printDeleteSubmitted("Snapshot", snap.Name())
	waitUntilGone(ctx, "Snapshot "+snap.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromStorage().Snapshots().Get(ctx, snap)
		return err
	})
}
