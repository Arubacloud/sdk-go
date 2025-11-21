package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

type securityGroupsClientImpl struct {
	client    *restclient.Client
	vpcClient *vpcsClientImpl
}

// NewService creates a new unified Network service
func NewSecurityGroupsClientImpl(client *restclient.Client, vpcClient *vpcsClientImpl) *securityGroupsClientImpl {
	return &securityGroupsClientImpl{
		client:    client,
		vpcClient: vpcClient,
	}
}

// List retrieves all security groups for a VPC
func (c *securityGroupsClientImpl) List(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error) {
	c.client.Logger().Debugf("Listing security groups for VPC: %s in project: %s", vpcId, project)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityGroupList](httpResp)
}

// Get retrieves a specific security group by ID
func (c *securityGroupsClientImpl) Get(ctx context.Context, project string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	c.client.Logger().Debugf("Getting security group: %s from VPC: %s in project: %s", securityGroupId, vpcId, project)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityGroupResponse](httpResp)
}

// Create creates a new security group in a VPC
// The SDK automatically waits for the VPC to become Active before creating the security group
func (c *securityGroupsClientImpl) Create(ctx context.Context, project string, vpcId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	c.client.Logger().Debugf("Creating security group in VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	// Wait for VPC to become Active before creating security group
	err := waitForVPCActive(ctx, c.vpcClient, project, vpcId)
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

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
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

// Update updates an existing security group
func (c *securityGroupsClientImpl) Update(ctx context.Context, project string, vpcId string, securityGroupId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error) {
	c.client.Logger().Debugf("Updating security group: %s in VPC: %s in project: %s", securityGroupId, vpcId, project)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
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

// Delete deletes a security group by ID
func (c *securityGroupsClientImpl) Delete(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting security group: %s from VPC: %s in project: %s", securityGroupId, vpcId, projectId)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}
