package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Restore tests
func TestListRestores(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123/restores" {
				w.WriteHeader(http.StatusOK)
				resp := map[string]interface{}{
					"values": []map[string]interface{}{
						{
							"metadata": map[string]interface{}{
								"name": "test-restore",
							},
							"properties": map[string]interface{}{
								"destinationVolume": map[string]interface{}{
									"uri": "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789",
								},
							},
						},
					},
					"total": 1,
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		backupClient := NewBackupClientImpl(c)
		svc := NewRestoreClientImpl(c, backupClient)
		resp, err := svc.List(context.Background(), "test-project", "backup-123", &types.RequestParameters{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatalf("resp is nil")
		}
		if resp.Data == nil {
			t.Logf("Raw response body: %s", string(resp.RawBody))
			t.Fatalf("resp.Data is nil")
		}
		if resp.Data.Values == nil {
			t.Fatalf("resp.Data.Values is nil")
		}
		if len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 Restore")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "test-restore" {
			t.Errorf("expected name 'test-restore'")
		}
		if resp.Data.Values[0].Properties.Destination.URI == "" {
			t.Errorf("expected Destination URI to be set")
		}
	})
}

func TestGetRestore(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123/restores/restore-123" {
				w.WriteHeader(http.StatusOK)
				resp := map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "test-restore",
						"id":   "restore-123",
					},
					"properties": map[string]interface{}{
						"destinationVolume": map[string]interface{}{
							"uri": "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789",
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		backupClient := NewBackupClientImpl(c)
		svc := NewRestoreClientImpl(c, backupClient)
		resp, err := svc.Get(context.Background(), "test-project", "backup-123", "restore-123", &types.RequestParameters{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Logf("Raw response body: %s", string(resp.RawBody))
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-restore" {
			t.Errorf("expected name 'test-restore'")
		}
		if resp.Data.Properties.Destination.URI == "" {
			t.Errorf("expected Destination URI to be set")
		}
	})
}

func TestCreateRestore(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123/restores" {
				w.WriteHeader(http.StatusCreated)
				resp := map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "new-restore",
						"id":   "restore-456",
					},
					"properties": map[string]interface{}{
						"destinationVolume": map[string]interface{}{
							"uri": "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789",
						},
					},
					"status": map[string]interface{}{
						"state": "creating",
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		backupClient := NewBackupClientImpl(c)
		svc := NewRestoreClientImpl(c, backupClient)
		body := types.RestoreRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-restore",
				},
			},
			Properties: types.RestorePropertiesRequest{
				Target: types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789"},
			},
		}
		resp, err := svc.Create(context.Background(), "test-project", "backup-123", body, &types.RequestParameters{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Logf("Raw response body: %s", string(resp.RawBody))
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-restore" {
			t.Errorf("expected name 'new-restore'")
		}
		if resp.Data.Properties.Destination.URI == "" {
			t.Errorf("expected Destination URI to be set")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateRestore(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123/restores/restore-123" {
				w.WriteHeader(http.StatusOK)
				resp := map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "updated-restore",
						"id":   "restore-123",
					},
					"properties": map[string]interface{}{
						"destinationVolume": map[string]interface{}{
							"uri": "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789",
						},
					},
					"status": map[string]interface{}{
						"state": "updating",
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		backupClient := NewBackupClientImpl(c)
		svc := NewRestoreClientImpl(c, backupClient)
		body := types.RestoreRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "updated-restore",
				},
			},
			Properties: types.RestorePropertiesRequest{
				Target: types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-789"},
			},
		}
		resp, err := svc.Update(context.Background(), "test-project", "backup-123", "restore-123", body, &types.RequestParameters{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Logf("Raw response body: %s", string(resp.RawBody))
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-restore" {
			t.Errorf("expected name 'updated-restore'")
		}
		if resp.Data.Properties.Destination.URI == "" {
			t.Errorf("expected Destination URI to be set")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "updating" {
			t.Errorf("expected state 'updating'")
		}
	})
}

func TestDeleteRestore(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123/restores/restore-123" {
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		backupClient := NewBackupClientImpl(c)
		svc := NewRestoreClientImpl(c, backupClient)
		_, err = svc.Delete(context.Background(), "test-project", "backup-123", "restore-123", &types.RequestParameters{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
