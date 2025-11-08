package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
func (s *SecurityGroupRuleService) ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleList], error) {
	s.client.Logger().Debugf("Listing security group rules for security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

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

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SecurityRuleList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.SecurityRuleList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetSecurityGroupRule retrieves a specific security group rule by ID
func (s *SecurityGroupRuleService) GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Getting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, project)

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

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SecurityRuleResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.SecurityRuleResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateSecurityGroupRule creates a new security group rule
func (s *SecurityGroupRuleService) CreateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Creating security group rule in security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

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

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SecurityRuleResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.SecurityRuleResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateSecurityGroupRule updates an existing security group rule
func (s *SecurityGroupRuleService) UpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Updating security group rule: %s in security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, project)

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

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SecurityRuleResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.SecurityRuleResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteSecurityGroupRule deletes a security group rule by ID
func (s *SecurityGroupRuleService) DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, projectId)

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

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}
