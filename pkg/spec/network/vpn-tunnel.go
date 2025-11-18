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

// ListVpnTunnels retrieves all VPN tunnels for a project
func (s *Service) ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelList], error) {
	s.client.Logger().Debugf("Listing VPN tunnels for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNTunnelsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VPNTunnelListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNTunnelListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VPNTunnelList](httpResp)
}

// GetVpnTunnel retrieves a specific VPN tunnel by ID
func (s *Service) GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error) {
	s.client.Logger().Debugf("Getting VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := schema.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNTunnelPath, project, vpnTunnelId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VPNTunnelGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNTunnelGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.VPNTunnelResponse](httpResp)
}

// CreateVpnTunnel creates a new VPN tunnel
func (s *Service) CreateVpnTunnel(ctx context.Context, project string, body schema.VPNTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error) {
	s.client.Logger().Debugf("Creating VPN tunnel in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNTunnelsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VPNTunnelCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNTunnelCreateAPIVersion
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

	response := &schema.Response[schema.VPNTunnelResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VPNTunnelResponse
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

// UpdateVpnTunnel updates an existing VPN tunnel
func (s *Service) UpdateVpnTunnel(ctx context.Context, project string, vpnTunnelId string, body schema.VPNTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error) {
	s.client.Logger().Debugf("Updating VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if err := schema.ValidateProjectAndResource(project, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNTunnelPath, project, vpnTunnelId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VPNTunnelUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNTunnelUpdateAPIVersion
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

	response := &schema.Response[schema.VPNTunnelResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VPNTunnelResponse
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

// DeleteVpnTunnel deletes a VPN tunnel by ID
func (s *Service) DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPN tunnel: %s in project: %s", vpnTunnelId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, vpnTunnelId, "VPN tunnel ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VPNTunnelPath, projectId, vpnTunnelId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &VPNTunnelDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &VPNTunnelDeleteAPIVersion
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
