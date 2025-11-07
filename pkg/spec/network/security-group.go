package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// SecurityGroupService implements the SecurityGroupAPI interface
type SecurityGroupService struct {
	client *client.Client
}

// NewSecurityGroupService creates a new SecurityGroupService
func NewSecurityGroupService(client *client.Client) *SecurityGroupService {
	return &SecurityGroupService{
		client: client,
	}
}

// ListSecurityGroups retrieves all security groups for a VPC
func (s *SecurityGroupService) ListSecurityGroups(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetSecurityGroup retrieves a specific security group by ID
func (s *SecurityGroupService) GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupPath, project, vpcId, securityGroupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateSecurityGroup creates or updates a security group
func (s *SecurityGroupService) CreateOrUpdateSecurityGroup(ctx context.Context, project string, vpcId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteSecurityGroup deletes a security group by ID
func (s *SecurityGroupService) DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupPath, projectId, vpcId, securityGroupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
