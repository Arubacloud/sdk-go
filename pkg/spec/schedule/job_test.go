package schedule

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListScheduleJobs(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusOK)
				resp := schema.JobList{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.JobResponse{
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("daily-backup"),
								ID:   schema.StringPtr("job-123"),
							},
							Properties: schema.JobPropertiesResponse{
								Enabled:       true,
								JobType:       schema.TypeJobRecurring,
								Cron:          schema.StringPtr("0 2 * * *"),
								ExecuteUntil:  schema.StringPtr("2025-12-31T23:59:59Z"),
								NextExecution: schema.StringPtr("2025-11-12T02:00:00Z"),
								Steps: []schema.JobStepResponse{
									{
										Name:        schema.StringPtr("backup-step"),
										ResourceURI: schema.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-123"),
										ActionURI:   schema.StringPtr("/snapshot"),
										HttpVerb:    schema.StringPtr("POST"),
									},
								},
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("active"),
							},
						},
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("oneshot-task"),
								ID:   schema.StringPtr("job-456"),
							},
							Properties: schema.JobPropertiesResponse{
								Enabled:    true,
								JobType:    schema.TypeJobOneShot,
								ScheduleAt: schema.StringPtr("2025-11-15T10:00:00Z"),
								Steps: []schema.JobStepResponse{
									{
										Name:        schema.StringPtr("shutdown-step"),
										ResourceURI: schema.StringPtr("/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123"),
										ActionURI:   schema.StringPtr("/stop"),
										HttpVerb:    schema.StringPtr("POST"),
									},
								},
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("scheduled"),
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

		resp, err := svc.ListScheduleJobs(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 jobs")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "daily-backup" {
			t.Errorf("expected name 'daily-backup'")
		}
		if resp.Data.Values[0].Properties.JobType != schema.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
		if resp.Data.Values[1].Properties.JobType != schema.TypeJobOneShot {
			t.Errorf("expected one-shot job type")
		}
	})

	t.Run("empty_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusOK)
				resp := schema.JobList{
					ListResponse: schema.ListResponse{Total: 0},
					Values:       []schema.JobResponse{},
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

		resp, err := svc.ListScheduleJobs(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected empty list")
		}
	})
}

func TestGetScheduleJob(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.JobResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("daily-backup"),
						ID:   schema.StringPtr("job-123"),
					},
					Properties: schema.JobPropertiesResponse{
						Enabled:       true,
						JobType:       schema.TypeJobRecurring,
						Cron:          schema.StringPtr("0 2 * * *"),
						ExecuteUntil:  schema.StringPtr("2025-12-31T23:59:59Z"),
						NextExecution: schema.StringPtr("2025-11-12T02:00:00Z"),
						Recurrency:    (*schema.RecurrenceType)(schema.StringPtr(string(schema.RecurrenceTypeDaily))),
						Steps: []schema.JobStepResponse{
							{
								Name:         schema.StringPtr("backup-step"),
								ResourceURI:  schema.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-123"),
								ActionURI:    schema.StringPtr("/snapshot"),
								ActionName:   schema.StringPtr("CreateSnapshot"),
								Typology:     schema.StringPtr("Aruba.Storage/block-storages"),
								TypologyName: schema.StringPtr("Block Storage"),
								HttpVerb:     schema.StringPtr("POST"),
								Body:         schema.StringPtr(`{"name":"daily-snapshot"}`),
							},
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

		resp, err := svc.GetScheduleJob(context.Background(), "test-project", "job-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "daily-backup" {
			t.Errorf("expected name 'daily-backup'")
		}
		if resp.Data.Properties.JobType != schema.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
		if resp.Data.Properties.Cron == nil || *resp.Data.Properties.Cron != "0 2 * * *" {
			t.Errorf("expected cron '0 2 * * *'")
		}
		if len(resp.Data.Properties.Steps) != 1 {
			t.Errorf("expected 1 step")
		}
	})
}

func TestCreateScheduleJob(t *testing.T) {
	t.Run("successful_create_recurring", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.JobResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("weekly-cleanup"),
						ID:   schema.StringPtr("job-789"),
					},
					Properties: schema.JobPropertiesResponse{
						Enabled:       true,
						JobType:       schema.TypeJobRecurring,
						Cron:          schema.StringPtr("0 3 * * 0"),
						ExecuteUntil:  schema.StringPtr("2026-01-01T00:00:00Z"),
						NextExecution: schema.StringPtr("2025-11-17T03:00:00Z"),
						Recurrency:    (*schema.RecurrenceType)(schema.StringPtr(string(schema.RecurrenceTypeWeekly))),
						Steps: []schema.JobStepResponse{
							{
								Name:        schema.StringPtr("cleanup-old-snapshots"),
								ResourceURI: schema.StringPtr("/projects/test-project/providers/Aruba.Storage/snapshots"),
								ActionURI:   schema.StringPtr("/cleanup"),
								HttpVerb:    schema.StringPtr("POST"),
							},
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

		body := schema.JobRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "weekly-cleanup",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.JobPropertiesRequest{
				Enabled:      true,
				JobType:      schema.TypeJobRecurring,
				Cron:         schema.StringPtr("0 3 * * 0"),
				ExecuteUntil: schema.StringPtr("2026-01-01T00:00:00Z"),
				Steps: []schema.JobStep{
					{
						Name:        schema.StringPtr("cleanup-old-snapshots"),
						ResourceURI: "/projects/test-project/providers/Aruba.Storage/snapshots",
						ActionURI:   "/cleanup",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.CreateScheduleJob(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "weekly-cleanup" {
			t.Errorf("expected name 'weekly-cleanup'")
		}
		if resp.Data.Properties.JobType != schema.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
	})

	t.Run("successful_create_oneshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.JobResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("maintenance-window"),
						ID:   schema.StringPtr("job-999"),
					},
					Properties: schema.JobPropertiesResponse{
						Enabled:    true,
						JobType:    schema.TypeJobOneShot,
						ScheduleAt: schema.StringPtr("2025-11-20T22:00:00Z"),
						Steps: []schema.JobStepResponse{
							{
								Name:        schema.StringPtr("stop-servers"),
								ResourceURI: schema.StringPtr("/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123"),
								ActionURI:   schema.StringPtr("/stop"),
								HttpVerb:    schema.StringPtr("POST"),
							},
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("scheduled"),
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

		body := schema.JobRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "maintenance-window",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.JobPropertiesRequest{
				Enabled:    true,
				JobType:    schema.TypeJobOneShot,
				ScheduleAt: schema.StringPtr("2025-11-20T22:00:00Z"),
				Steps: []schema.JobStep{
					{
						Name:        schema.StringPtr("stop-servers"),
						ResourceURI: "/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123",
						ActionURI:   "/stop",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.CreateScheduleJob(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "maintenance-window" {
			t.Errorf("expected name 'maintenance-window'")
		}
		if resp.Data.Properties.JobType != schema.TypeJobOneShot {
			t.Errorf("expected one-shot job type")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "scheduled" {
			t.Errorf("expected state 'scheduled'")
		}
	})
}

func TestUpdateScheduleJob(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.JobResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("updated-backup"),
						ID:   schema.StringPtr("job-123"),
					},
					Properties: schema.JobPropertiesResponse{
						Enabled:        false,
						JobType:        schema.TypeJobRecurring,
						Cron:           schema.StringPtr("0 4 * * *"),
						ExecuteUntil:   schema.StringPtr("2025-12-31T23:59:59Z"),
						NextExecution:  schema.StringPtr("2025-11-12T04:00:00Z"),
						DeactiveReason: (*schema.DeactiveReasonDto)(schema.StringPtr(string(schema.DeactiveReasonManual))),
						Steps: []schema.JobStepResponse{
							{
								Name:        schema.StringPtr("updated-backup-step"),
								ResourceURI: schema.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-456"),
								ActionURI:   schema.StringPtr("/snapshot"),
								HttpVerb:    schema.StringPtr("POST"),
							},
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("inactive"),
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

		body := schema.JobRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "updated-backup",
				},
				Location: schema.LocationRequest{Value: "it-eur"},
			},
			Properties: schema.JobPropertiesRequest{
				Enabled:      false,
				JobType:      schema.TypeJobRecurring,
				Cron:         schema.StringPtr("0 4 * * *"),
				ExecuteUntil: schema.StringPtr("2025-12-31T23:59:59Z"),
				Steps: []schema.JobStep{
					{
						Name:        schema.StringPtr("updated-backup-step"),
						ResourceURI: "/projects/test-project/providers/Aruba.Storage/block-storages/vol-456",
						ActionURI:   "/snapshot",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.UpdateScheduleJob(context.Background(), "test-project", "job-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-backup" {
			t.Errorf("expected name 'updated-backup'")
		}
		if resp.Data.Properties.Enabled {
			t.Errorf("expected job to be disabled")
		}
		if resp.Data.Properties.Cron == nil || *resp.Data.Properties.Cron != "0 4 * * *" {
			t.Errorf("expected updated cron '0 4 * * *'")
		}
	})
}

func TestDeleteScheduleJob(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
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

		_, err = svc.DeleteScheduleJob(context.Background(), "test-project", "job-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
