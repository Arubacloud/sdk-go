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

type vpcPeeringRoutesClientImpl struct {
	client *restclient.Client
}

// NewService creates a new unified Network service
func NewVPCPeeringRoutesClientImpl(client *restclient.Client) *vpcPeeringRoutesClientImpl {
	return &vpcPeeringRoutesClientImpl{
		client: client,
	}
}

// List retrieves all VPC peering routes for a VPC peering connection
func (c *vpcPeeringRoutesClientImpl) List(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error) {
	c.client.Logger().Debugf("Listing VPC peering routes for VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringRoutesPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringRouteListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringRouteListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCPeeringRouteList](httpResp)
}

// Get retrieves a specific VPC peering route by ID
func (c *vpcPeeringRoutesClientImpl) Get(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error) {
	c.client.Logger().Debugf("Getting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCPeeringRoute(project, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringRouteGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringRouteGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPCPeeringRouteResponse](httpResp)
}

// Create creates a new VPC peering route
func (c *vpcPeeringRoutesClientImpl) Create(ctx context.Context, project string, vpcId string, vpcPeeringId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error) {
	c.client.Logger().Debugf("Creating VPC peering route in VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringRoutesPath, project, vpcId, vpcPeeringId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringRouteCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringRouteCreateAPIVersion
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

	response := &types.Response[types.VPCPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCPeeringRouteResponse
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

// Update updates an existing VPC peering route
func (c *vpcPeeringRoutesClientImpl) Update(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error) {
	c.client.Logger().Debugf("Updating VPC peering route: %s in VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if err := types.ValidateVPCPeeringRoute(project, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringRouteUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringRouteUpdateAPIVersion
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

	response := &types.Response[types.VPCPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPCPeeringRouteResponse
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

// Delete deletes a VPC peering route by ID
func (c *vpcPeeringRoutesClientImpl) Delete(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, projectId)

	if err := types.ValidateVPCPeeringRoute(projectId, vpcId, vpcPeeringId, vpcPeeringRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPCPeeringRoutePath, projectId, vpcId, vpcPeeringId, vpcPeeringRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPCPeeringRouteDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPCPeeringRouteDeleteAPIVersion
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
