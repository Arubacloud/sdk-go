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

func TestListBackups(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups" {
				w.WriteHeader(http.StatusOK)
				resp := types.StorageBackupList{
					ListResponse: types.ListResponse{Total: 1},
					Values: []types.StorageBackupResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("test-backup"),
							},
							Properties: types.StorageBackupPropertiesResult{
								Type: types.StorageBackupTypeFull,
								Origin: types.ReferenceResource{
									URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123",
								},
								RetentionDays: types.IntPtr(10),
								BillingPeriod: types.StringPtr("Monthly"),
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("active"),
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewBackupClientImpl(c)
		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 Backup")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "test-backup" {
			t.Errorf("expected name 'test-backup'")
		}
		if resp.Data.Values[0].Properties.Type != types.StorageBackupTypeFull {
			t.Errorf("expected Type 'Full'")
		}
		if resp.Data.Values[0].Properties.Origin.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123" {
			t.Errorf("expected Origin URI match")
		}
		if resp.Data.Values[0].Properties.RetentionDays == nil || *resp.Data.Values[0].Properties.RetentionDays != 10 {
			t.Errorf("expected RetentionDays 10")
		}
		if resp.Data.Values[0].Properties.BillingPeriod == nil || *resp.Data.Values[0].Properties.BillingPeriod != "Monthly" {
			t.Errorf("expected BillingPeriod 'Monthly'")
		}
	})
}

func TestGetBackup(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.StorageBackupResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("test-backup"),
						ID:   types.StringPtr("backup-123"),
					},
					Properties: types.StorageBackupPropertiesResult{
						Type:          types.StorageBackupTypeFull,
						Origin:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123"},
						RetentionDays: types.IntPtr(10),
						BillingPeriod: types.StringPtr("Monthly"),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("active"),
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
		svc := NewBackupClientImpl(c)
		resp, err := svc.Get(context.Background(), "test-project", "backup-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-backup" {
			t.Errorf("expected name 'test-backup'")
		}
		if resp.Data.Properties.Type != types.StorageBackupTypeFull {
			t.Errorf("expected Type 'Full'")
		}
		if resp.Data.Properties.Origin.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123" {
			t.Errorf("expected Origin URI match")
		}
		if resp.Data.Properties.RetentionDays == nil || *resp.Data.Properties.RetentionDays != 10 {
			t.Errorf("expected RetentionDays 10")
		}
		if resp.Data.Properties.BillingPeriod == nil || *resp.Data.Properties.BillingPeriod != "Monthly" {
			t.Errorf("expected BillingPeriod 'Monthly'")
		}
	})
}

func TestCreateBackup(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups" {
				w.WriteHeader(http.StatusCreated)
				resp := types.StorageBackupResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-backup"),
						ID:   types.StringPtr("backup-456"),
					},
					Properties: types.StorageBackupPropertiesResult{
						Type:          types.StorageBackupTypeFull,
						Origin:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-456"},
						RetentionDays: types.IntPtr(20),
						BillingPeriod: types.StringPtr("Yearly"),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("creating"),
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
		svc := NewBackupClientImpl(c)
		body := types.StorageBackupRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-backup",
				},
			},
			Properties: types.StorageBackupPropertiesRequest{
				StorageBackupType: types.StorageBackupTypeFull,
				Origin:            types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-456"},
				RetentionDays:     types.IntPtr(20),
				BillingPeriod:     types.StringPtr("Yearly"),
			},
		}
		resp, err := svc.Create(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-backup" {
			t.Errorf("expected name 'new-backup'")
		}
		if resp.Data.Properties.Type != types.StorageBackupTypeFull {
			t.Errorf("expected Type 'Full'")
		}
		if resp.Data.Properties.Origin.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-456" {
			t.Errorf("expected Origin URI match")
		}
		if resp.Data.Properties.RetentionDays == nil || *resp.Data.Properties.RetentionDays != 20 {
			t.Errorf("expected RetentionDays 20")
		}
		if resp.Data.Properties.BillingPeriod == nil || *resp.Data.Properties.BillingPeriod != "Yearly" {
			t.Errorf("expected BillingPeriod 'Yearly'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateBackup(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.StorageBackupResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("updated-backup"),
						ID:   types.StringPtr("backup-123"),
					},
					Properties: types.StorageBackupPropertiesResult{
						Type:          types.StorageBackupTypeIncremental,
						Origin:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123"},
						RetentionDays: types.IntPtr(30),
						BillingPeriod: types.StringPtr("Monthly"),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("updating"),
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
		svc := NewBackupClientImpl(c)
		body := types.StorageBackupRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "updated-backup",
				},
			},
			Properties: types.StorageBackupPropertiesRequest{
				StorageBackupType: types.StorageBackupTypeIncremental,
				Origin:            types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123"},
				RetentionDays:     types.IntPtr(30),
				BillingPeriod:     types.StringPtr("Monthly"),
			},
		}
		resp, err := svc.Update(context.Background(), "test-project", "backup-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-backup" {
			t.Errorf("expected name 'updated-backup'")
		}
		if resp.Data.Properties.Type != types.StorageBackupTypeIncremental {
			t.Errorf("expected Type 'Incremental'")
		}
		if resp.Data.Properties.Origin.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/volume-123" {
			t.Errorf("expected Origin URI match")
		}
		if resp.Data.Properties.RetentionDays == nil || *resp.Data.Properties.RetentionDays != 30 {
			t.Errorf("expected RetentionDays 30")
		}
		if resp.Data.Properties.BillingPeriod == nil || *resp.Data.Properties.BillingPeriod != "Monthly" {
			t.Errorf("expected BillingPeriod 'Monthly'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "updating" {
			t.Errorf("expected state 'updating'")
		}
	})
}

func TestDeleteBackup(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/backups/backup-123" {
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
		svc := NewBackupClientImpl(c)
		_, err = svc.Delete(context.Background(), "test-project", "backup-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
