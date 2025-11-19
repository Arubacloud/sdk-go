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

// ListVPCs retrieves all VPCs for a project
func (s *Service) ListVPCs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.VPCList], error) {
	s.client.Logger().Debugf("Listing VPCs for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCNetworksPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCList](httpResp)
}

// GetVPC retrieves a specific VPC by ID
func (s *Service) GetVPC(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCResponse], error) {
	s.client.Logger().Debugf("Getting VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCNetworkPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCResponse](httpResp)
}

// CreateVPC creates a new VPC
func (s *Service) CreateVPC(ctx context.Context, project string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error) {
	s.client.Logger().Debugf("Creating VPC in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCNetworksPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCCreateAPIVersion
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

	response := &types.Response[types.VPCResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCResponse
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

// UpdateVPC updates an existing VPC
func (s *Service) UpdateVPC(ctx context.Context, project string, vpcId string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error) {
	s.client.Logger().Debugf("Updating VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCNetworkPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCUpdateAPIVersion
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

	response := &types.Response[types.VPCResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCResponse
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

// DeleteVPC deletes a VPC by ID
func (s *Service) DeleteVPC(ctx context.Context, projectId string, vpcId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPC: %s in project: %s", vpcId, projectId)

	if err := types.ValidateProjectAndResource(projectId, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCNetworkPath, projectId, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCDeleteAPIVersion
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
