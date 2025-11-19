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

// ListVpcPeerings retrieves all VPC peerings for a VPC
func (s *Service) ListVpcPeerings(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error) {
	s.client.Logger().Debugf("Listing VPC peerings for VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCPeeringList](httpResp)
}

// GetVpcPeering retrieves a specific VPC peering by ID
func (s *Service) GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error) {
	s.client.Logger().Debugf("Getting VPC peering: %s from VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCPeeringResponse](httpResp)
}

// CreateVpcPeering creates a new VPC peering
func (s *Service) CreateVpcPeering(ctx context.Context, project string, vpcId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error) {
	s.client.Logger().Debugf("Creating VPC peering in VPC: %s in project: %s", vpcId, project)

	if err := types.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringsPath, project, vpcId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringCreateAPIVersion
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

	response := &types.Response[types.VPCPeeringResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCPeeringResponse
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

// UpdateVpcPeering updates an existing VPC peering
func (s *Service) UpdateVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error) {
	s.client.Logger().Debugf("Updating VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringUpdateAPIVersion
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

	response := &types.Response[types.VPCPeeringResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCPeeringResponse
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

// DeleteVpcPeering deletes a VPC peering by ID
func (s *Service) DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPC peering: %s from VPC: %s in project: %s", vpcPeeringId, vpcId, projectId)

	if err := types.ValidateVPCResource(projectId, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringPath, projectId, vpcId, vpcPeeringId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringDeleteAPIVersion
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
