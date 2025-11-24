package metric

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListMetrics(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/metrics" {
				w.WriteHeader(http.StatusOK)
				resp := types.MetricListResponse{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.MetricResponse{
						{
							ReferenceID:   "resource-123",
							Name:          "cpu_usage",
							ReferenceName: "test-server",
							Metadata: []types.MetricMetadata{
								{
									Field: "unit",
									Value: "percent",
								},
							},
							Data: []types.MetricData{
								{
									Time:    "2024-01-01T00:00:00Z",
									Measure: "45.5",
								},
								{
									Time:    "2024-01-01T00:01:00Z",
									Measure: "48.2",
								},
							},
						},
						{
							ReferenceID:   "resource-123",
							Name:          "memory_usage",
							ReferenceName: "test-server",
							Metadata: []types.MetricMetadata{
								{
									Field: "unit",
									Value: "MB",
								},
							},
							Data: []types.MetricData{
								{
									Time:    "2024-01-01T00:00:00Z",
									Measure: "2048",
								},
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
		svc := NewMetricsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 metrics")
		}
		if resp.Data.Values[0].Name != "cpu_usage" {
			t.Errorf("expected metric name 'cpu_usage', got %s", resp.Data.Values[0].Name)
		}
		if resp.Data.Values[0].ReferenceID != "resource-123" {
			t.Errorf("expected reference ID 'resource-123', got %s", resp.Data.Values[0].ReferenceID)
		}
		if len(resp.Data.Values[0].Data) != 2 {
			t.Errorf("expected 2 data points for cpu_usage")
		}
		if len(resp.Data.Values[1].Metadata) != 1 {
			t.Errorf("expected 1 metadata field for memory_usage")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/metrics" {
				w.WriteHeader(http.StatusOK)
				resp := types.MetricListResponse{
					ListResponse: types.ListResponse{Total: 0},
					Values:       []types.MetricResponse{},
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
		svc := NewMetricsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected 0 metrics")
		}
	})
}

func TestListAlerts(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				w.WriteHeader(http.StatusOK)
				resp := types.AlertsListResponse{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.AlertResponse{
						{
							ID:                 "alert-123",
							EventID:            "event-456",
							EventName:          "High CPU Usage",
							Username:           "user@example.com",
							ServiceCategory:    "Compute",
							ServiceTypology:    "VirtualMachine",
							ResourceID:         "vm-789",
							ServiceName:        "test-vm",
							ResourceTypology:   "CloudServer",
							Metric:             "cpu_usage",
							LastReception:      time.Now(),
							Rule:               "greater_than",
							Theshold:           80,
							UM:                 "%",
							Duration:           "5m",
							ThesholdExceedence: "85%",
							Component:          "CPU",
							Email:              true,
							Panel:              true,
							SMS:                false,
							Hidden:             false,
							ExecutedAlertActions: []types.ExecutedAlertAction{
								{
									ActionType:   types.ActionTypeSendEmail,
									Success:      true,
									ErrorMessage: "",
								},
							},
							Actions: []types.AlertAction{
								{
									Key:        "acknowledge",
									Disabled:   false,
									Executable: true,
								},
							},
						},
						{
							ID:                 "alert-456",
							EventID:            "event-789",
							EventName:          "Low Disk Space",
							Username:           "admin@example.com",
							ServiceCategory:    "Storage",
							ServiceTypology:    "BlockStorage",
							ResourceID:         "disk-321",
							ServiceName:        "test-disk",
							ResourceTypology:   "Volume",
							Metric:             "disk_usage",
							LastReception:      time.Now(),
							Rule:               "greater_than",
							Theshold:           90,
							UM:                 "%",
							Duration:           "10m",
							ThesholdExceedence: "95%",
							Email:              true,
							Panel:              true,
							SMS:                true,
							Hidden:             false,
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
		svc := NewAlertsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 alerts")
		}
		if resp.Data.Values[0].ID != "alert-123" {
			t.Errorf("expected alert ID 'alert-123', got %s", resp.Data.Values[0].ID)
		}
		if resp.Data.Values[0].EventName != "High CPU Usage" {
			t.Errorf("expected event name 'High CPU Usage', got %s", resp.Data.Values[0].EventName)
		}
		if resp.Data.Values[0].Metric != "cpu_usage" {
			t.Errorf("expected metric 'cpu_usage', got %s", resp.Data.Values[0].Metric)
		}
		if resp.Data.Values[0].Theshold != 80 {
			t.Errorf("expected threshold 80, got %d", resp.Data.Values[0].Theshold)
		}
		if !resp.Data.Values[0].Email {
			t.Errorf("expected email notification to be enabled")
		}
		if len(resp.Data.Values[0].ExecutedAlertActions) != 1 {
			t.Errorf("expected 1 executed action")
		}
		if resp.Data.Values[1].SMS != true {
			t.Errorf("expected SMS notification to be enabled for second alert")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				w.WriteHeader(http.StatusOK)
				resp := types.AlertsListResponse{
					ListResponse: types.ListResponse{Total: 0},
					Values:       []types.AlertResponse{},
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
		svc := NewAlertsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected 0 alerts")
		}
	})

	t.Run("with filtering", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				// Verify filter parameters were passed
				filter := r.URL.Query().Get("filter")
				if filter != "resourceID eq 'vm-789'" {
					t.Errorf("expected filter 'resourceID eq 'vm-789'', got %s", filter)
				}

				w.WriteHeader(http.StatusOK)
				resp := types.AlertsListResponse{
					ListResponse: types.ListResponse{Total: 1},
					Values: []types.AlertResponse{
						{
							ID:            "alert-123",
							EventName:     "Filtered Alert",
							ResourceID:    "vm-789",
							Metric:        "cpu_usage",
							LastReception: time.Now(),
							Email:         true,
							Panel:         true,
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
		svc := NewAlertsClientImpl(c)

		params := &types.RequestParameters{
			Filter: types.StringPtr("resourceID eq 'vm-789'"),
		}

		resp, err := svc.List(context.Background(), "test-project", params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 filtered alert")
		}
		if resp.Data.Values[0].ResourceID != "vm-789" {
			t.Errorf("expected resource ID 'vm-789' in filtered result")
		}
	})
}
