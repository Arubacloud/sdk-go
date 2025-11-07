package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// VPCNetworkService implements the VPCNetworkAPI interface
type VPCNetworkService struct {
	client *client.Client
}

// NewVPCNetworkService creates a new VPCNetworkService
func NewVPCNetworkService(client *client.Client) *VPCNetworkService {
	return &VPCNetworkService{
		client: client,
	}
}

// ListVPCNetworks retrieves all VPC networks for a project
func (s *VPCNetworkService) ListVPCNetworks(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(VPCNetworksPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetVPCNetwork retrieves a specific VPC network by ID
func (s *VPCNetworkService) GetVPCNetwork(ctx context.Context, project string, vpcNetworkId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcNetworkId == "" {
		return nil, fmt.Errorf("VPC network ID cannot be empty")
	}

	path := fmt.Sprintf(VPCNetworkPath, project, vpcNetworkId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateVPCNetwork creates or updates a VPC network
func (s *VPCNetworkService) CreateOrUpdateVPCNetwork(ctx context.Context, project string, body schema.VPCNetworkRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(VPCNetworksPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteVPCNetwork deletes a VPC network by ID
func (s *VPCNetworkService) DeleteVPCNetwork(ctx context.Context, projectId string, vpcNetworkId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcNetworkId == "" {
		return nil, fmt.Errorf("VPC network ID cannot be empty")
	}

	path := fmt.Sprintf(VPCNetworkPath, projectId, vpcNetworkId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
