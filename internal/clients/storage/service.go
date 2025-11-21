package storage

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// waitForBlockStorageActive waits for a Block Storage volume to become Active or NotUsed before proceeding
func waitForBlockStorageActive(ctx context.Context, volumeClient *volumesClientImpl, projectID, volumeID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := volumeClient.Get(ctx, projectID, volumeID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("BlockStorage state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	config := restclient.DefaultPollingConfig()
	// BlockStorage can be "Active" (attached) or "NotUsed" (unattached but ready)
	config.SuccessStates = []string{"Active", "NotUsed"}

	return volumeClient.client.WaitForResourceState(ctx, "BlockStorage", volumeID, getter, config)
}
