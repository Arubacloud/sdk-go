package network

import (
	"context"
	"fmt"
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
func (s *VpnTunnelService) ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetVpnTunnel retrieves a specific VPN tunnel by ID
func (s *VpnTunnelService) GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateVpnTunnel creates or updates a VPN tunnel
func (s *VpnTunnelService) CreateOrUpdateVpnTunnel(ctx context.Context, project string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteVpnTunnel deletes a VPN tunnel by ID
func (s *VpnTunnelService) DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
