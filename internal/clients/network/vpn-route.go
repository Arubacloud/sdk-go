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

type vpnRoutesClientImpl struct {
	client *restclient.Client
}

// NewService creates a new unified Network service
func NewVPNRoutesClientImpl(client *restclient.Client) *vpnRoutesClientImpl {
	return &vpnRoutesClientImpl{
		client: client,
	}
}

// List retrieves all VPN routes for a VPN tunnel
func (c *vpnRoutesClientImpl) List(ctx context.Context, project string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[types.VPNRouteList], error) {
	c.client.Logger().Debugf("Listing VPN routes for VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := types.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNRoutesPath, project, vpnTunnelId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPNRouteListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNRouteListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPNRouteList](httpResp)
}

// Get retrieves a specific VPN route by ID
func (c *vpnRoutesClientImpl) Get(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error) {
	c.client.Logger().Debugf("Getting VPN route: %s from VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, project)

	if err := types.ValidateVPNRoute(project, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNRoutePath, project, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPNRouteGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNRouteGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.VPNRouteResponse](httpResp)
}

// Create creates a new VPN route in a VPN tunnel
func (c *vpnRoutesClientImpl) Create(ctx context.Context, project string, vpnTunnelId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error) {
	c.client.Logger().Debugf("Creating VPN route in VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := types.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNRoutesPath, project, vpnTunnelId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPNRouteCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNRouteCreateAPIVersion
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

	response := &types.Response[types.VPNRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPNRouteResponse
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

// Update updates an existing VPN route
func (c *vpnRoutesClientImpl) Update(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error) {
	c.client.Logger().Debugf("Updating VPN route: %s in VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, project)

	if err := types.ValidateVPNRoute(project, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNRoutePath, project, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPNRouteUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNRouteUpdateAPIVersion
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

	response := &types.Response[types.VPNRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.VPNRouteResponse
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

// Delete deletes a VPN route by ID
func (c *vpnRoutesClientImpl) Delete(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting VPN route: %s from VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, projectId)

	if err := types.ValidateVPNRoute(projectId, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNRoutePath, projectId, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &VPNRouteDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNRouteDeleteAPIVersion
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
