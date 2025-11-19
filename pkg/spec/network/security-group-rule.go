package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListSecurityGroupRules retrieves all security group rules for a security group
func (s *Service) ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error) {
	s.client.Logger().Debugf("Listing security group rules for security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, securityGroupId, "security group ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, project, vpcId, securityGroupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityRuleList](httpResp)
}

// GetSecurityGroupRule retrieves a specific security group rule by ID
func (s *Service) GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Getting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, project)

	if err := types.ValidateSecurityGroupRule(project, vpcId, securityGroupId, securityGroupRuleId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, project, vpcId, securityGroupId, securityGroupRuleId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityRuleResponse](httpResp)
}

// CreateSecurityGroupRule creates a new security group rule
// The SDK automatically waits for the SecurityGroup to become Active before creating the rule
func (s *Service) CreateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Creating security group rule in security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, securityGroupId, "security group ID"); err != nil {
		return nil, err
	}

	// Wait for SecurityGroup to become Active before creating rule
	err := s.waitForSecurityGroupActive(ctx, project, vpcId, securityGroupId)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for SecurityGroup to become active: %w", err)
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, project, vpcId, securityGroupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleCreateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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

	response := &types.Response[types.SecurityRuleResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SecurityRuleResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateSecurityGroupRule updates an existing security group rule
func (s *Service) UpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	s.client.Logger().Debugf("Updating security group rule: %s in security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, project)

	if err := types.ValidateSecurityGroupRule(project, vpcId, securityGroupId, securityGroupRuleId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, project, vpcId, securityGroupId, securityGroupRuleId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleUpdateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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

	response := &types.Response[types.SecurityRuleResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SecurityRuleResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteSecurityGroupRule deletes a security group rule by ID
func (s *Service) DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleId, securityGroupId, vpcId, projectId)

	if err := types.ValidateSecurityGroupRule(projectId, vpcId, securityGroupId, securityGroupRuleId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, projectId, vpcId, securityGroupId, securityGroupRuleId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}
