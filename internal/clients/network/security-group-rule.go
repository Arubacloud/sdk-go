package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type securityGroupRulesClientImpl struct {
	client               *restclient.Client
	securityGroupsClient *securityGroupsClientImpl
}

// NewService creates a new unified Network service
func NewSecurityGroupRulesClientImpl(client *restclient.Client, securityGroupsClient *securityGroupsClientImpl) *securityGroupRulesClientImpl {
	return &securityGroupRulesClientImpl{
		client:               client,
		securityGroupsClient: securityGroupsClient,
	}
}

// List retrieves all security group rules for a security group
func (c *securityGroupRulesClientImpl) List(ctx context.Context, projectID string, vpcID string, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error) {
	c.client.Logger().Debugf("Listing security group rules for security group: %s in VPC: %s in project: %s", securityGroupID, vpcID, projectID)

	if err := types.ValidateVPCResource(projectID, vpcID, securityGroupID, "security group ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, projectID, vpcID, securityGroupID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityRuleList](httpResp)
}

// Get retrieves a specific security group rule by ID
func (c *securityGroupRulesClientImpl) Get(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	c.client.Logger().Debugf("Getting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleID, securityGroupID, vpcID, projectID)

	if err := types.ValidateSecurityGroupRule(projectID, vpcID, securityGroupID, securityGroupRuleID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, projectID, vpcID, securityGroupID, securityGroupRuleID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SecurityRuleResponse](httpResp)
}

// Create creates a new security group rule
// The SDK automatically waits for the SecurityGroup to become Active before creating the rule
func (c *securityGroupRulesClientImpl) Create(ctx context.Context, projectID string, vpcID string, securityGroupID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	c.client.Logger().Debugf("Creating security group rule in security group: %s in VPC: %s in project: %s", securityGroupID, vpcID, projectID)

	if err := types.ValidateVPCResource(projectID, vpcID, securityGroupID, "security group ID"); err != nil {
		return nil, err
	}

	// Wait for SecurityGroup to become Active before creating rule
	err := waitForSecurityGroupActive(ctx, *c.securityGroupsClient, projectID, vpcID, securityGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for SecurityGroup to become active: %w", err)
	}

	path := fmt.Sprintf(SecurityGroupRulesPath, projectID, vpcID, securityGroupID)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
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

// Update updates an existing security group rule
func (c *securityGroupRulesClientImpl) Update(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error) {
	c.client.Logger().Debugf("Updating security group rule: %s in security group: %s in VPC: %s in project: %s", securityGroupRuleID, securityGroupID, vpcID, projectID)

	if err := types.ValidateSecurityGroupRule(projectID, vpcID, securityGroupID, securityGroupRuleID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, projectID, vpcID, securityGroupID, securityGroupRuleID)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
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

// Delete deletes a security group rule by ID
func (c *securityGroupRulesClientImpl) Delete(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting security group rule: %s from security group: %s in VPC: %s in project: %s", securityGroupRuleID, securityGroupID, vpcID, projectID)

	if err := types.ValidateSecurityGroupRule(projectID, vpcID, securityGroupID, securityGroupRuleID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SecurityGroupRulePath, projectID, vpcID, securityGroupID, securityGroupRuleID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SecurityRuleDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SecurityRuleDeleteAPIVersion
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
