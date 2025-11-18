package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListVpcPeeringRoutes retrieves all VPC peering routes for a VPC peering connection
func (s *Service) ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteList], error) {
	s.client.Logger().Debugf("Listing VPC peering routes for VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpcPeeringRouteListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpcPeeringRouteListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VPCPeeringRouteList](httpResp)
}

// GetVpcPeeringRoute retrieves a specific VPC peering route by ID
func (s *Service) GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Getting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCPeeringRoute(project, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpcPeeringRouteGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpcPeeringRouteGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VPCPeeringRouteResponse](httpResp)
}

// CreateVpcPeeringRoute creates a new VPC peering route
func (s *Service) CreateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VPCPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Creating VPC peering route in VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpcPeeringRouteCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpcPeeringRouteCreateAPIVersion
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

	response := &schema.Response[schema.VPCPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VPCPeeringRouteResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateVpcPeeringRoute updates an existing VPC peering route
func (s *Service) UpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body schema.VPCPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Updating VPC peering route: %s in VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCPeeringRoute(project, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpcPeeringRouteUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpcPeeringRouteUpdateAPIVersion
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

	response := &schema.Response[schema.VPCPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VPCPeeringRouteResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteVpcPeeringRoute deletes a VPC peering route by ID
func (s *Service) DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, projectId)

	if err := schema.ValidateVPCPeeringRoute(projectId, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, projectId, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpcPeeringRouteDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpcPeeringRouteDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[any](httpResp)
}
