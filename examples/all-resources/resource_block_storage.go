package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createBlockStorage creates a block storage volume with the given name
func createBlockStorage(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, name string) *aruba.BlockStorage {
	fmt.Printf("--- Block Storage (%s) ---\n", name)

	bs := aruba.NewBlockStorage().
		InProject(proj).
		Named(name).
		Tagged("storage").
		Tagged("data").
		InRegion(aruba.RegionITBGBergamo).
		InZone(aruba.ZoneITBG1).
		OfType(aruba.BlockStorageTypeStandard).
		SizedGB(20).
		AsBootable().
		BilledHourly().
		FromImage(aruba.VolumeImageLU22001)

	bs, err := arubaClient.FromStorage().Volumes().Create(ctx, bs)
	if err != nil {
		printCreateError("Block Storage", err)
		return nil
	}
	printCreated("Block Storage", bs.Name(), bs.ID())

	waitUntilSelfReady(ctx, "Block Storage", bs.Name(), bs, bs.WaitUntilReady)

	return bs
}

// deleteBlockStorage deletes a block storage volume and waits for the platform
// to confirm removal. Project deletion fails with 400 if a volume is still in
// Deleting state.
func deleteBlockStorage(ctx context.Context, arubaClient aruba.Client, bs *aruba.BlockStorage) {
	printDeleteBanner("Block Storage")
	if err := arubaClient.FromStorage().Volumes().Delete(ctx, bs); err != nil {
		printDeleteError("Block Storage", err)
		return
	}
	printDeleteSubmitted("Block Storage", bs.Name())
	waitUntilGone(ctx, "Block Storage "+bs.Name(), bs.WaitUntilGone)
}
