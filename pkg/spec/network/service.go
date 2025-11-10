package network

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the NetworkAPI interface for all Network operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Network service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

// waitForVPCActive waits for a VPC to become Active before proceeding
func (s *Service) waitForVPCActive(ctx context.Context, projectID, vpcID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := s.GetVPC(ctx, projectID, vpcID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("VPC state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	return s.client.WaitForResourceState(ctx, "VPC", vpcID, getter, client.DefaultPollingConfig())
}

// waitForSecurityGroupActive waits for a Security Group to become Active before proceeding
func (s *Service) waitForSecurityGroupActive(ctx context.Context, projectID, vpcID, sgID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := s.GetSecurityGroup(ctx, projectID, vpcID, sgID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("SecurityGroup state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	return s.client.WaitForResourceState(ctx, "SecurityGroup", sgID, getter, client.DefaultPollingConfig())
}
