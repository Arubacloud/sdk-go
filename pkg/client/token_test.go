package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testContext returns a context for testing
func testContext() context.Context {
	return context.Background()
}

func TestTokenManager_GetToken(t *testing.T) {
	// Create a mock token server
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}

		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}

		if r.FormValue("client_id") != "test-client-id" {
			t.Errorf("Expected client_id=test-client-id, got %s", r.FormValue("client_id"))
		}

		if r.FormValue("client_secret") != "test-client-secret" {
			t.Errorf("Expected client_secret=test-client-secret, got %s", r.FormValue("client_secret"))
		}

		// Return a valid token response
		resp := TokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer tokenServer.Close()

	tm := NewTokenManager(
		tokenServer.URL,
		"test-client-id",
		"test-client-secret",
		http.DefaultClient,
		5*time.Minute,
	)

	ctx := testContext()

	// Test getting a token
	token, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	if token != "test-access-token" {
		t.Errorf("Expected token 'test-access-token', got '%s'", token)
	}

	// Test token caching (should not make a new request)
	token2, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() (cached) error = %v", err)
	}

	if token2 != token {
		t.Errorf("Cached token should be the same")
	}
}

func TestTokenManager_RefreshToken(t *testing.T) {
	callCount := 0

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		resp := TokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer tokenServer.Close()

	tm := NewTokenManager(
		tokenServer.URL,
		"test-client-id",
		"test-client-secret",
		http.DefaultClient,
		5*time.Minute,
	)

	ctx := testContext()

	if err := tm.RefreshToken(ctx); err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call to token server, got %d", callCount)
	}

	if !tm.IsTokenValid() {
		t.Error("Token should be valid after refresh")
	}
}

func TestTokenManager_TokenExpiration(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a token that expires in 2 seconds
		resp := TokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   2,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer tokenServer.Close()

	tm := NewTokenManager(
		tokenServer.URL,
		"test-client-id",
		"test-client-secret",
		http.DefaultClient,
		1*time.Millisecond, // Minimal refresh buffer for testing
	)

	ctx := testContext()

	if err := tm.RefreshToken(ctx); err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	// Small delay to ensure token is set
	time.Sleep(10 * time.Millisecond)

	// Check token details
	token, expiresAt, isValid := tm.GetTokenInfo()
	now := time.Now()
	timeWithBuffer := now.Add(tm.tokenRefreshBuffer)
	beforeCheck := timeWithBuffer.Before(expiresAt)
	t.Logf("Token: %s", token)
	t.Logf("ExpiresAt: %v", expiresAt)
	t.Logf("Now: %v", now)
	t.Logf("Buffer: %v", tm.tokenRefreshBuffer)
	t.Logf("Now+Buffer: %v", timeWithBuffer)
	t.Logf("(Now+Buffer).Before(ExpiresAt): %v", beforeCheck)
	t.Logf("IsValid from method: %v", isValid)

	if !isValid {
		t.Errorf("Token should be valid immediately after refresh")
	}

	// Wait for token to expire
	time.Sleep(3 * time.Second)

	if tm.IsTokenValid() {
		t.Error("Token should be invalid after expiration")
	}
}

func TestTokenManager_ErrorHandling(t *testing.T) {
	// Test with a server that returns an error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_client"}`))
	}))
	defer errorServer.Close()

	tm := NewTokenManager(
		errorServer.URL,
		"invalid-client-id",
		"invalid-client-secret",
		http.DefaultClient,
		5*time.Minute,
	)

	ctx := testContext()

	_, err := tm.GetToken(ctx)
	if err == nil {
		t.Error("Expected error with invalid credentials, got nil")
	}
}

func TestNewClient_WithTokenManager(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := TokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer tokenServer.Close()

	config := &Config{
		BaseURL:        "https://api.example.com",
		TokenIssuerURL: tokenServer.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		HTTPClient:     http.DefaultClient,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	// Test getting token through client
	token, err := client.GetToken(context.Background())
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	if token != "test-access-token" {
		t.Errorf("Expected token 'test-access-token', got '%s'", token)
	}
}
