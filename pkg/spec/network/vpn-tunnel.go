package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// VpnTunnelService implements the VpnTunnelAPI interface
type VpnTunnelService struct {
	client *client.Client
}

// NewVpnTunnelService creates a new VpnTunnelService
func NewVpnTunnelService(client *client.Client) *VpnTunnelService {
	return &VpnTunnelService{
		client: client,
	}
}

// ListVpnTunnels retrieves all VPN tunnels for a project
func (s *VpnTunnelService) ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelList], error) {
	s.client.Logger().Debugf("Listing VPN tunnels for project: %s", project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(VpnTunnelsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.VpnTunnelList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnTunnelList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetVpnTunnel retrieves a specific VPN tunnel by ID
func (s *VpnTunnelService) GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error) {
	s.client.Logger().Debugf("Getting VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}

	path := fmt.Sprintf(VpnTunnelPath, project, vpnTunnelId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.VpnTunnelResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnTunnelResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateVpnTunnel creates a new VPN tunnel
func (s *VpnTunnelService) CreateVpnTunnel(ctx context.Context, project string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error) {
	s.client.Logger().Debugf("Creating VPN tunnel in project: %s", project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(VpnTunnelsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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

	response := &schema.Response[schema.VpnTunnelResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnTunnelResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateVpnTunnel updates an existing VPN tunnel
func (s *VpnTunnelService) UpdateVpnTunnel(ctx context.Context, project string, vpnTunnelId string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error) {
	s.client.Logger().Debugf("Updating VPN tunnel: %s in project: %s", vpnTunnelId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}

	path := fmt.Sprintf(VpnTunnelPath, project, vpnTunnelId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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

	response := &schema.Response[schema.VpnTunnelResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpnTunnelResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteVpnTunnel deletes a VPN tunnel by ID
func (s *VpnTunnelService) DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPN tunnel: %s in project: %s", vpnTunnelId, projectId)

	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}

	path := fmt.Sprintf(VpnTunnelPath, projectId, vpnTunnelId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}
