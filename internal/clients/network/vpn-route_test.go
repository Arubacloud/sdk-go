package network

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListVpnRoutes(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"total":1,"values":[{"metadata":{"name":"route-1"},"properties":{"cloudSubnet":"10.0.0.0/24","onPremSubnet":"192.168.1.0/24"}}]}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.List(context.Background(), "test-project", "tunnel-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 vpn route")
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.List(context.Background(), "", "tunnel-123", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnTunnelID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.List(context.Background(), "test-project", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.List(context.Background(), "test-project", "tunnel-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.List(context.Background(), "test-project", "tunnel-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.List(context.Background(), "test-project", "tunnel-123", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"total":0,"values":[]}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.List(context.Background(), "test-project", "tunnel-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestGetVpnRoute(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"metadata":{"name":"route-1"},"properties":{"cloudSubnet":"10.0.0.0/24","onPremSubnet":"192.168.1.0/24"}}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Get(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "route-1" {
			t.Errorf("expected route name 'route-1', got %v", resp.Data.Metadata.Name)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Get(context.Background(), "", "tunnel-123", "route-1", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnTunnelID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Get(context.Background(), "test-project", "", "route-1", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnRouteID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Get(context.Background(), "test-project", "tunnel-123", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Get(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Get(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Get(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Get(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestCreateVpnRoute(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"metadata":{"name":"new-route"}}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Create(context.Background(), "test-project", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Create(context.Background(), "", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnTunnelID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Create(context.Background(), "test-project", "", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Create(context.Background(), "test-project", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		// TODO(TD-010): Create/Update's manual response build silently swallows non-JSON
		// unmarshal errors (diverges from ParseResponseBody which logs at DEBUG).
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Create(context.Background(), "test-project", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Create(context.Background(), "test-project", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Create(context.Background(), "test-project", "tunnel-123", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateVpnRoute(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"metadata":{"name":"updated-route"}}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Update(context.Background(), "test-project", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Update(context.Background(), "", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnTunnelID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Update(context.Background(), "test-project", "", "route-1", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnRouteID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Update(context.Background(), "test-project", "tunnel-123", "", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Update(context.Background(), "test-project", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		// TODO(TD-010): Create/Update's manual response build silently swallows non-JSON
		// unmarshal errors (diverges from ParseResponseBody which logs at DEBUG).
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Update(context.Background(), "test-project", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Update(context.Background(), "test-project", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Update(context.Background(), "test-project", "tunnel-123", "route-1", types.VPNRouteRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteVpnRoute(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "", "tunnel-123", "route-1", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnTunnelID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "test-project", "", "route-1", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpnRouteID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		resp, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewVPNRoutesClientImpl(c)
		_, err := svc.Delete(context.Background(), "test-project", "tunnel-123", "route-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
