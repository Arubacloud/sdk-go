package client

import (
	"encoding/json"
	"testing"
)

func TestTokenResponse_Unmarshal(t *testing.T) {
	// Real response from Aruba Cloud API
	realResponse := `{
		"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICItYVRnU0pVSnkxMUVhNjY2c3NZTFIxYzlDSlhqWG9FMi1xdFVXRHozYVdJIn0.eyJleHAiOjE3NjE5Mjc3MzUsImlhdCI6MTc2MTkyNDEzNSwianRpIjoiYjJkYWQ2ZjUtYThmMi00YzE4LWEyNzctNWJkZTE0NTJhODljIiwiaXNzIjoiaHR0cHM6Ly9sb2dpbi5hcnViYS5pdC9hdXRoL3JlYWxtcy9jbXAtbmV3LWFwaWtleSIsInN1YiI6Ijg2MWY5YTkzLWFkMTctNGU5Ny1hZmU5LWNmZWNhOTMyZThkOSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNtcC03ZmZlMGVlZS1lNDViLTQxYzUtODY0Yi0zYjE3OGVlYWNiMmQiLCJhY3IiOiIxIiwic2NvcGUiOiJlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiY2xpZW50SG9zdCI6IjgyLjE5Mi4xMzEuMSIsImNvbXBhbnkiOiJBUlUiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhcnUtMjk3NjQ3IiwiZ3JvdXBNZW1iZXJzaGlwIjpbIlN0YW5kYXJkIl0sImNsaWVudEFkZHJlc3MiOiI4Mi4xOTIuMTMxLjEiLCJjdXJyZW5jeUNvZGUiOiJFVVIiLCJjbGllbnRfaWQiOiJjbXAtN2ZmZTBlZWUtZTQ1Yi00MWM1LTg2NGItM2IxNzhlZWFjYjJkIiwidGVuYW50IjoiQVJVIn0.HN8MF7aHmCiuEcnDDc_Muur_35o49rMK5_LBkFzSJJPxqCnsM1jqq4PY2Ul70KsWvJ4yz86Urot2WK8zcc2i6LhUyvh8OZNOrPniCeu96_RHwhQEuw3URMxeA8mqyeEkmesco6SX5fWd1zktkhuEOiy11jMVNV3qs5OGMkvKZnPthkCead8o_U2jY_0ok5fe0tweQPvGDhhFahqjM5COzf5c9eM2_lmGt8R-dzaxMWYLaVYZk6Y1Ff2uT77LLIAOzJCtb_8cnywFDD5DujCXBuYPQzm9559d5PDHxRMikYCafp4DEynyfw0xtgV1NtXvrLNXroDZYj7nTZmXW3mXXQ",
		"expires_in": 3600,
		"refresh_expires_in": 0,
		"token_type": "Bearer",
		"not-before-policy": 0,
		"scope": "email"
	}`

	var response TokenResponse
	err := json.Unmarshal([]byte(realResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Validate fields
	if response.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}

	if response.TokenType != "Bearer" {
		t.Errorf("Expected TokenType 'Bearer', got '%s'", response.TokenType)
	}

	if response.ExpiresIn != 3600 {
		t.Errorf("Expected ExpiresIn 3600, got %d", response.ExpiresIn)
	}

	if response.RefreshExpiresIn != 0 {
		t.Errorf("Expected RefreshExpiresIn 0 (client credentials), got %d", response.RefreshExpiresIn)
	}

	if response.Scope != "email" {
		t.Errorf("Expected Scope 'email', got '%s'", response.Scope)
	}

	if response.NotBeforePolicy != 0 {
		t.Errorf("Expected NotBeforePolicy 0, got %d", response.NotBeforePolicy)
	}

	t.Logf("Successfully parsed token response")
	t.Logf("Token starts with: %s...", response.AccessToken[:50])
	t.Logf("Expires in: %d seconds", response.ExpiresIn)
}

func TestTokenManager_GetRemainingTime(t *testing.T) {
	tokenServer := setupMockTokenServer(t)
	defer tokenServer.Close()

	tm := NewTokenManager(
		tokenServer.URL,
		"test-client-id",
		"test-client-secret",
		nil,
	)

	// Before getting token
	remaining := tm.GetRemainingTime()
	if remaining != 0 {
		t.Errorf("Expected remaining time 0 before token, got %v", remaining)
	}

	// Get token
	ctx := testContext()
	_, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	// After getting token
	remaining = tm.GetRemainingTime()
	if remaining <= 0 {
		t.Errorf("Expected positive remaining time after getting token, got %v", remaining)
	}

	t.Logf("Token remaining time: %v", remaining)
}

func TestTokenManager_GetTokenInfo(t *testing.T) {
	tokenServer := setupMockTokenServer(t)
	defer tokenServer.Close()

	tm := NewTokenManager(
		tokenServer.URL,
		"test-client-id",
		"test-client-secret",
		nil,
	)

	ctx := testContext()

	// Get token
	_, err := tm.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	// Get token info
	token, expiresAt, isValid := tm.GetTokenInfo()

	if token == "" {
		t.Error("Expected non-empty token")
	}

	if expiresAt.IsZero() {
		t.Error("Expected non-zero expiration time")
	}

	if !isValid {
		t.Error("Expected token to be valid")
	}

	if len(token) > 20 {
		t.Logf("Token: %s...", token[:20])
	} else {
		t.Logf("Token: %s", token)
	}
	t.Logf("Expires at: %v", expiresAt)
	t.Logf("Is valid: %v", isValid)
}
