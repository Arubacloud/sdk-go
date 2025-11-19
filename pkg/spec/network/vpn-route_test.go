package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

func TestListVpnRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpntunnels/tunnel-123/routes" {
			w.WriteHeader(http.StatusOK)
			resp := types.VPNRouteList{
				ListResponse: types.ListResponse{Total: 1},
				Values: []types.VPNRouteResponse{
					{
						Metadata: types.ResourceMetadataResponse{
							Name: types.StringPtr("route-1"),
						},
						Properties: types.VPNRoutePropertiesResponse{
							CloudSubnet:  "10.0.0.0/24",
							OnPremSubnet: "192.168.1.0/24",
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.ListVpnRoutes(context.Background(), "test-project", "tunnel-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
		t.Errorf("expected 1 vpn route")
	}
}

func TestGetVpnRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpntunnels/tunnel-123/routes/route-1" {
			w.WriteHeader(http.StatusOK)
			resp := types.VPNRouteResponse{
				Metadata: types.ResourceMetadataResponse{
					Name: types.StringPtr("route-1"),
				},
				Properties: types.VPNRoutePropertiesResponse{
					CloudSubnet:  "10.0.0.0/24",
					OnPremSubnet: "192.168.1.0/24",
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.GetVpnRoute(context.Background(), "test-project", "tunnel-123", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "route-1" {
		t.Errorf("expected route name 'route-1'")
	}
}

func TestDeleteVpnRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpntunnels/tunnel-123/routes/route-1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	_, err = svc.DeleteVpnRoute(context.Background(), "test-project", "tunnel-123", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
