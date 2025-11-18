package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListBlockStorageVolumes(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages" {
				w.WriteHeader(http.StatusOK)
				resp := schema.BlockStorageList{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.BlockStorageResponse{
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("data-volume"),
								ID:   schema.StringPtr("vol-123"),
							},
							Properties: schema.BlockStoragePropertiesResponse{
								SizeGB:        100,
								BillingPeriod: "Hour",
								Zone:          "it-eur-1",
								Type:          schema.BlockStorageTypePerformance,
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("active"),
							},
						},
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("backup-volume"),
								ID:   schema.StringPtr("vol-456"),
							},
							Properties: schema.BlockStoragePropertiesResponse{
								SizeGB:        200,
								BillingPeriod: "Hour",
								Zone:          "it-eur-1",
								Type:          schema.BlockStorageTypeStandard,
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("active"),
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

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.ListBlockStorageVolumes(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 volumes")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "data-volume" {
			t.Errorf("expected name 'data-volume'")
		}
		if resp.Data.Values[0].Properties.Type != schema.BlockStorageTypePerformance {
			t.Errorf("expected performance type")
		}
	})
}

func TestGetBlockStorageVolume(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.BlockStorageResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("my-volume"),
						ID:   schema.StringPtr("vol-123"),
					},
					Properties: schema.BlockStoragePropertiesResponse{
						SizeGB:        150,
						BillingPeriod: "Hour",
						Zone:          "it-eur-1",
						Type:          schema.BlockStorageTypePerformance,
						Bootable:      schema.BoolPtr(false),
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("active"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.GetBlockStorageVolume(context.Background(), "test-project", "vol-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-volume" {
			t.Errorf("expected name 'my-volume'")
		}
		if resp.Data.Properties.SizeGB != 150 {
			t.Errorf("expected size 150GB")
		}
	})
}

func TestCreateBlockStorageVolume(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.BlockStorageResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("new-volume"),
						ID:   schema.StringPtr("vol-789"),
					},
					Properties: schema.BlockStoragePropertiesResponse{
						SizeGB:        50,
						BillingPeriod: "Hour",
						Zone:          "it-eur-1",
						Type:          schema.BlockStorageTypeStandard,
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("creating"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		body := schema.BlockStorageRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "new-volume",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.BlockStoragePropertiesRequest{
				SizeGB:        50,
				BillingPeriod: "Hour",
				Zone:          "it-eur-1",
				Type:          schema.BlockStorageTypeStandard,
			},
		}

		resp, err := svc.CreateBlockStorageVolume(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-volume" {
			t.Errorf("expected name 'new-volume'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestDeleteBlockStorageVolume(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		_, err = svc.DeleteBlockStorageVolume(context.Background(), "test-project", "vol-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestListSnapshots(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots" {
				w.WriteHeader(http.StatusOK)
				resp := schema.SnapshotList{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.SnapshotResponse{
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("backup-snapshot-1"),
								ID:   schema.StringPtr("snap-123"),
							},
							Properties: schema.SnapshotPropertiesResponse{
								SizeGb:        schema.Int32Ptr(100),
								BillingPeriod: schema.StringPtr("Hour"),
								Zone:          "it-eur-1",
								Type:          schema.BlockStorageTypePerformance,
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("available"),
							},
						},
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("backup-snapshot-2"),
								ID:   schema.StringPtr("snap-456"),
							},
							Properties: schema.SnapshotPropertiesResponse{
								SizeGb:        schema.Int32Ptr(200),
								BillingPeriod: schema.StringPtr("Hour"),
								Zone:          "it-eur-1",
								Type:          schema.BlockStorageTypeStandard,
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("available"),
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

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.ListSnapshots(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 snapshots")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "backup-snapshot-1" {
			t.Errorf("expected name 'backup-snapshot-1'")
		}
	})
}

func TestGetSnapshot(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots/snap-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.SnapshotResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("my-snapshot"),
						ID:   schema.StringPtr("snap-123"),
					},
					Properties: schema.SnapshotPropertiesResponse{
						SizeGb:        schema.Int32Ptr(150),
						BillingPeriod: schema.StringPtr("Hour"),
						Zone:          "it-eur-1",
						Type:          schema.BlockStorageTypePerformance,
						Volume: &schema.VolumeInfo{
							Uri:  schema.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123"),
							Name: schema.StringPtr("source-volume"),
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("available"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.GetSnapshot(context.Background(), "test-project", "snap-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-snapshot" {
			t.Errorf("expected name 'my-snapshot'")
		}
		if resp.Data.Properties.Volume == nil || resp.Data.Properties.Volume.Name == nil || *resp.Data.Properties.Volume.Name != "source-volume" {
			t.Errorf("expected volume name 'source-volume'")
		}
	})
}

func TestCreateSnapshot(t *testing.T) {
	t.Skip("Skipping test because CreateSnapshot waits for block storage to be active")
	t.Run("successful_create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			// Mock the GET request to check volume status (waitForBlockStorageActive)
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.BlockStorageResponse{
					Metadata: schema.ResourceMetadataResponse{
						ID: schema.StringPtr("vol-123"),
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("active"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.SnapshotResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("new-snapshot"),
						ID:   schema.StringPtr("snap-789"),
					},
					Properties: schema.SnapshotPropertiesResponse{
						SizeGb:        schema.Int32Ptr(50),
						BillingPeriod: schema.StringPtr("Hour"),
						Zone:          "it-eur-1",
						Type:          schema.BlockStorageTypeStandard,
						Volume: &schema.VolumeInfo{
							Uri: schema.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123"),
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("creating"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		body := schema.SnapshotRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "new-snapshot",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.SnapshotPropertiesRequest{
				BillingPeriod: schema.StringPtr("Hour"),
				Volume: schema.ReferenceResource{
					Uri: "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123",
				},
			},
		}

		resp, err := svc.CreateSnapshot(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-snapshot" {
			t.Errorf("expected name 'new-snapshot'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestDeleteSnapshot(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots/snap-123" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		_, err = svc.DeleteSnapshot(context.Background(), "test-project", "snap-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
