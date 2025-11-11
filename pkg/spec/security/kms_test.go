package security

import (
"context"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/Arubacloud/sdk-go/pkg/client"
"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListKMSKeys(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.URL.Path == "/token" {
w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Security/kms/keys" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KmsList{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.KmsResponse{
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("encryption-key-1"),
								Id:   schema.StringPtr("kms-123"),
							},
							Properties: schema.KmsPropertiesResponse{
								BillingPeriod: schema.BillingPeriodResource{
									BillingPeriod: "Month",
								},
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("active"),
							},
						},
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("encryption-key-2"),
								Id:   schema.StringPtr("kms-456"),
							},
							Properties: schema.KmsPropertiesResponse{
								BillingPeriod: schema.BillingPeriodResource{
									BillingPeriod: "Month",
								},
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

		resp, err := svc.ListKMSKeys(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 KMS keys")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "encryption-key-1" {
			t.Errorf("expected name 'encryption-key-1'")
		}
	})
}

func TestGetKMSKey(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.URL.Path == "/token" {
w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Security/kms/keys/kms-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KmsResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("my-encryption-key"),
						Id:   schema.StringPtr("kms-123"),
					},
					Properties: schema.KmsPropertiesResponse{
						BillingPeriod: schema.BillingPeriodResource{
							BillingPeriod: "Month",
						},
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

		resp, err := svc.GetKMSKey(context.Background(), "test-project", "kms-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-encryption-key" {
			t.Errorf("expected name 'my-encryption-key'")
		}
		if resp.Data.Properties.BillingPeriod.BillingPeriod != "Month" {
			t.Errorf("expected billing period 'Month'")
		}
	})
}

func TestCreateKMSKey(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.URL.Path == "/token" {
w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Security/kms/keys" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.KmsResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("new-encryption-key"),
						Id:   schema.StringPtr("kms-789"),
					},
					Properties: schema.KmsPropertiesResponse{
						BillingPeriod: schema.BillingPeriodResource{
							BillingPeriod: "Month",
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

		body := schema.KmsRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "new-encryption-key",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.KmsPropertiesRequest{
				BillingPeriod: schema.BillingPeriodResource{
					BillingPeriod: "Month",
				},
			},
		}

		resp, err := svc.CreateKMSKey(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-encryption-key" {
			t.Errorf("expected name 'new-encryption-key'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateKMSKey(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.URL.Path == "/token" {
w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Security/kms/keys/kms-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KmsResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("updated-encryption-key"),
						Id:   schema.StringPtr("kms-123"),
					},
					Properties: schema.KmsPropertiesResponse{
						BillingPeriod: schema.BillingPeriodResource{
							BillingPeriod: "Year",
						},
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

		body := schema.KmsRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "updated-encryption-key",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.KmsPropertiesRequest{
				BillingPeriod: schema.BillingPeriodResource{
					BillingPeriod: "Year",
				},
			},
		}

		resp, err := svc.UpdateKMSKey(context.Background(), "test-project", "kms-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-encryption-key" {
			t.Errorf("expected name 'updated-encryption-key'")
		}
		if resp.Data.Properties.BillingPeriod.BillingPeriod != "Year" {
			t.Errorf("expected billing period 'Year'")
		}
	})
}

func TestDeleteKMSKey(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.URL.Path == "/token" {
w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Security/kms/keys/kms-123" {
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

		_, err = svc.DeleteKMSKey(context.Background(), "test-project", "kms-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
