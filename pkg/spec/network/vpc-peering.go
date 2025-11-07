package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// VpcPeeringService implements the VpcPeeringAPI interface
type VpcPeeringService struct {
	client *client.Client
}

// NewVpcPeeringService creates a new VpcPeeringService
func NewVpcPeeringService(client *client.Client) *VpcPeeringService {
	return &VpcPeeringService{
		client: client,
	}
}

// ListVpcPeerings retrieves all VPC peerings for a VPC
func (s *VpcPeeringService) ListVpcPeerings(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetVpcPeering retrieves a specific VPC peering by ID
func (s *VpcPeeringService) GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringPath, project, vpcId, vpcPeeringId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateVpcPeering creates or updates a VPC peering
func (s *VpcPeeringService) CreateOrUpdateVpcPeering(ctx context.Context, project string, vpcId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteVpcPeering deletes a VPC peering by ID
func (s *VpcPeeringService) DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringPath, projectId, vpcId, vpcPeeringId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
