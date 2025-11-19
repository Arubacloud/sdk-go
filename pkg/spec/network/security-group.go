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

// ListSecurityGroups retrieves all security groups for a VPC
func (s *Service) ListSecurityGroups(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error) {
	s.client.Logger().Debugf("Listing security groups for VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityGroupListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityGroupListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityGroupList](httpResp)
}

// GetSecurityGroup retrieves a specific security group by ID
func (s *Service) GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	s.client.Logger().Debugf("Getting security group: %s from VPC: %s in project: %s", securityGroupId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, securityGroupId, "security group ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupPath, project, vpcId, securityGroupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityGroupGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityGroupGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityGroupResponse](httpResp)
}

// CreateSecurityGroup creates a new security group in a VPC
// The SDK automatically waits for the VPC to become Active before creating the security group
func (s *Service) CreateSecurityGroup(ctx context.Context, project string, vpcId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	s.client.Logger().Debugf("Creating security group in VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	// Wait for VPC to become Active before creating security group
	err := s.waitForVPCActive(ctx, project, vpcId)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for VPC to become active: %w", err)
	}

	path := fmt.Sprintf(SecurityGroupsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityGroupCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityGroupCreateAPIVersion
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

	response := &types.Response[types.SecurityGroupResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SecurityGroupResponse
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

// UpdateSecurityGroup updates an existing security group
func (s *Service) UpdateSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	s.client.Logger().Debugf("Updating security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, securityGroupId, "security group ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupPath, project, vpcId, securityGroupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityGroupUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityGroupUpdateAPIVersion
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

	response := &types.Response[types.SecurityGroupResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SecurityGroupResponse
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

// DeleteSecurityGroup deletes a security group by ID
func (s *Service) DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting security group: %s from VPC: %s in project: %s", securityGroupId, vpcId, projectId)

	if err := types.ValidateVPCResource(projectId, vpcId, securityGroupId, "security group ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupPath, projectId, vpcId, securityGroupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityGroupDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityGroupDeleteAPIVersion
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
