package network

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/restclient"
)

// waitForVPCActive waits for a VPC to become Active before proceeding
func waitForVPCActive(ctx context.Context, vpcClient *vpcsClientImpl, projectID, vpcID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := vpcClient.Get(ctx, projectID, vpcID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("VPC state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	return vpcClient.client.WaitForResourceState(ctx, "VPC", vpcID, getter, restclient.DefaultPollingConfig())
}

// waitForSecurityGroupActive waits for a Security Group to become Active before proceeding
func waitForSecurityGroupActive(ctx context.Context, securityGroupsClient securityGroupsClientImpl, projectID, vpcID, sgID string) error {
	getter := func(ctx context.Context) (string, error) {
		resp, err := securityGroupsClient.Get(ctx, projectID, vpcID, sgID, nil)
		if err != nil {
			return "", err
		}
		if resp.Data == nil || resp.Data.Status.State == nil {
			return "", fmt.Errorf("SecurityGroup state is nil")
		}
		return *resp.Data.Status.State, nil
	}

	return securityGroupsClient.client.WaitForResourceState(ctx, "SecurityGroup", sgID, getter, restclient.DefaultPollingConfig())
}
