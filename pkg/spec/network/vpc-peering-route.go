package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// VpcPeeringRouteService implements the VpcPeeringRouteAPI interface
type VpcPeeringRouteService struct {
	client *client.Client
}

// NewVpcPeeringRouteService creates a new VpcPeeringRouteService
func NewVpcPeeringRouteService(client *client.Client) *VpcPeeringRouteService {
	return &VpcPeeringRouteService{
		client: client,
	}
}

// ListVpcPeeringRoutes retrieves all VPC peering routes for a VPC peering connection
func (s *VpcPeeringRouteService) ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetVpcPeeringRoute retrieves a specific VPC peering route by ID
func (s *VpcPeeringRouteService) GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteId == "" {
		return nil, fmt.Errorf("VPC peering route ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateVpcPeeringRoute creates or updates a VPC peering route
func (s *VpcPeeringRouteService) CreateOrUpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteVpcPeeringRoute deletes a VPC peering route by ID
func (s *VpcPeeringRouteService) DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteId == "" {
		return nil, fmt.Errorf("VPC peering route ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, projectId, vpcId, vpcPeeringId, vpcPeeringRouteId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
