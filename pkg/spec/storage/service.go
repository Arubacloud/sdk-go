package storage

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the StorageAPI interface for all Storage operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Storage service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

// waitForBlockStorageActive waits for a Block Storage volume to become Active or NotUsed before proceeding
func (s *Service) waitForBlockStorageActive(ctx context.Context, projectID, volumeID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := s.GetBlockStorageVolume(ctx, projectID, volumeID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("BlockStorage state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	config := client.DefaultPollingConfig()
	// BlockStorage can be "Active" (attached) or "NotUsed" (unattached but ready)
	config.SuccessStates = []string{"Active", "NotUsed"}

	return s.client.WaitForResourceState(ctx, "BlockStorage", volumeID, getter, config)
}
