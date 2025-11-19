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

// ListSubnets retrieves all subnets for a VPC
func (s *Service) ListSubnets(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.SubnetList], error) {
	s.client.Logger().Debugf("Listing subnets for VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SubnetListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SubnetListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SubnetList](httpResp)
}

// GetSubnet retrieves a specific subnet by ID
func (s *Service) GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error) {
	s.client.Logger().Debugf("Getting subnet: %s from VPC: %s in project: %s", subnetId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, project, vpcId, subnetId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SubnetGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SubnetGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SubnetResponse](httpResp)
}

// CreateSubnet creates a new subnet in a VPC
// The SDK automatically waits for the VPC to become Active before creating the subnet
func (s *Service) CreateSubnet(ctx context.Context, project string, vpcId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error) {
	s.client.Logger().Debugf("Creating subnet in VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	// Wait for VPC to become Active before creating subnet
	err := s.waitForVPCActive(ctx, project, vpcId)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for VPC to become active: %w", err)
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SubnetCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SubnetCreateAPIVersion
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

	response := &types.Response[types.SubnetResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SubnetResponse
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

// UpdateSubnet updates an existing subnet
func (s *Service) UpdateSubnet(ctx context.Context, project string, vpcId string, subnetId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error) {
	s.client.Logger().Debugf("Updating subnet: %s in VPC: %s in project: %s", subnetId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, project, vpcId, subnetId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SubnetUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SubnetUpdateAPIVersion
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

	response := &types.Response[types.SubnetResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.SubnetResponse
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

// DeleteSubnet deletes a subnet by ID
func (s *Service) DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting subnet: %s from VPC: %s in project: %s", subnetId, vpcId, projectId)

	if err := types.ValidateVPCResource(projectId, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, projectId, vpcId, subnetId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SubnetDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SubnetDeleteAPIVersion
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
