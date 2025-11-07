package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// SecurityGroupRuleService implements the SecurityGroupRuleAPI interface
type SecurityGroupRuleService struct {
	client *client.Client
}

// NewSecurityGroupRuleService creates a new SecurityGroupRuleService
func NewSecurityGroupRuleService(client *client.Client) *SecurityGroupRuleService {
	return &SecurityGroupRuleService{
		client: client,
	}
}

// ListSecurityGroupRules retrieves all security group rules for a security group
func (s *SecurityGroupRuleService) ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, project, vpcId, securityGroupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetSecurityGroupRule retrieves a specific security group rule by ID
func (s *SecurityGroupRuleService) GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}
	if securityGroupRuleId == "" {
		return nil, fmt.Errorf("security group rule ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupRulePath, project, vpcId, securityGroupId, securityGroupRuleId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateSecurityGroupRule creates or updates a security group rule
func (s *SecurityGroupRuleService) CreateOrUpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, project, vpcId, securityGroupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteSecurityGroupRule deletes a security group rule by ID
func (s *SecurityGroupRuleService) DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupId == "" {
		return nil, fmt.Errorf("security group ID cannot be empty")
	}
	if securityGroupRuleId == "" {
		return nil, fmt.Errorf("security group rule ID cannot be empty")
	}

	path := fmt.Sprintf(SecurityGroupRulePath, projectId, vpcId, securityGroupId, securityGroupRuleId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
