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

// ListVpnRoutes retrieves all VPN routes for a VPN tunnel
func (s *Service) ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteList], error) {
	s.client.Logger().Debugf("Listing VPN routes for VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := schema.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpnRoutesPath, project, vpnTunnelId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpnRouteListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpnRouteListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VpnRouteList](httpResp)
}

// GetVpnRoute retrieves a specific VPN route by ID
func (s *Service) GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error) {
	s.client.Logger().Debugf("Getting VPN route: %s from VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, project)

	if err := schema.ValidateVPNRoute(project, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpnRoutePath, project, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpnRouteGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpnRouteGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VpnRouteResponse](httpResp)
}

// CreateVpnRoute creates a new VPN route in a VPN tunnel
func (s *Service) CreateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error) {
	s.client.Logger().Debugf("Creating VPN route in VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := schema.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpnRoutesPath, project, vpnTunnelId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpnRouteCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpnRouteCreateAPIVersion
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

	response := &schema.Response[schema.VpnRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnRouteResponse
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

// UpdateVpnRoute updates an existing VPN route
func (s *Service) UpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error) {
	s.client.Logger().Debugf("Updating VPN route: %s in VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, project)

	if err := schema.ValidateVPNRoute(project, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpnRoutePath, project, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpnRouteUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpnRouteUpdateAPIVersion
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

	response := &schema.Response[schema.VpnRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnRouteResponse
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

// DeleteVpnRoute deletes a VPN route by ID
func (s *Service) DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPN route: %s from VPN tunnel: %s in project: %s", vpnRouteId, vpnTunnelId, projectId)

	if err := schema.ValidateVPNRoute(projectId, vpnTunnelId, vpnRouteId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpnRoutePath, projectId, vpnTunnelId, vpnRouteId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VpnRouteDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VpnRouteDeleteAPIVersion
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
