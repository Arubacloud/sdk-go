package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListVPCs(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.VPCList{
				ListResponse: types.ListResponse{Total: 1},
				Values: []types.VPCResponse{
					{Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("vpc-1")}},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		cfg := &restclient.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVPCsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Total != 1 {
			t.Errorf("expected total 1, got %d", resp.Data.Total)
		}
	})
}

func TestGetVPC(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.VPCResponse{
				Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("my-vpc")},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		cfg := &restclient.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVPCsClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-vpc" {
			t.Errorf("expected name 'my-vpc', got '%v'", resp.Data.Metadata.Name)
		}
	})
}

func TestCreateVPC(t *testing.T) {
	// VPC Create doesn't require waiting, so this test should work fine
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.WriteHeader(http.StatusCreated)
			resp := types.VPCResponse{
				Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("new-vpc")},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		cfg := &restclient.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVPCsClientImpl(c)

		req := types.VPCRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "new-vpc"},
				Location:                types.LocationRequest{Value: "ITBG-Bergamo"},
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteVPC(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		cfg := &restclient.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVPCsClientImpl(c)

		_, err = svc.Delete(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
