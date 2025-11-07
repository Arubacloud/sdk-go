package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// SubnetService implements the SubnetAPI interface
type SubnetService struct {
	client *client.Client
}

// NewSubnetService creates a new SubnetService
func NewSubnetService(client *client.Client) *SubnetService {
	return &SubnetService{
		client: client,
	}
}

// ListSubnets retrieves all subnets for a VPC
func (s *SubnetService) ListSubnets(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetSubnet retrieves a specific subnet by ID
func (s *SubnetService) GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if subnetId == "" {
		return nil, fmt.Errorf("subnet ID cannot be empty")
	}

	path := fmt.Sprintf(SubnetPath, project, vpcId, subnetId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateSubnet creates or updates a subnet
func (s *SubnetService) CreateOrUpdateSubnet(ctx context.Context, project string, vpcId string, body schema.SubnetRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteSubnet deletes a subnet by ID
func (s *SubnetService) DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if subnetId == "" {
		return nil, fmt.Errorf("subnet ID cannot be empty")
	}

	path := fmt.Sprintf(SubnetPath, projectId, vpcId, subnetId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
