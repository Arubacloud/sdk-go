package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenResponse represents the OAuth2 token response
// Based on OAuth2 client credentials flow (no refresh_token)
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"` // Always 0 for client credentials
	Scope            string `json:"scope,omitempty"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
}

// TokenManager handles OAuth2 token acquisition
//
// Thread Safety:
// - Uses sync.RWMutex to protect accessToken and expiresAt
// - Multiple goroutines can read token simultaneously (RLock)
// - Only one goroutine can obtain a new token at a time (Lock)
// - Prevents "thundering herd" - multiple goroutines won't duplicate token requests

type TokenManager struct {
	tokenIssuerURL string
	clientID       string
	clientSecret   string
	httpClient     *http.Client

	mu          sync.RWMutex // Protects accessToken and expiresAt for thread safety
	accessToken string       // Current JWT token (protected by mu)
	expiresAt   time.Time    // When token expires (protected by mu)
}

// NewTokenManager creates a new token manager
func NewTokenManager(
	tokenIssuerURL string,
	clientID string,
	clientSecret string,
	httpClient *http.Client,
) *TokenManager {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &TokenManager{
		tokenIssuerURL: tokenIssuerURL,
		clientID:       clientID,
		clientSecret:   clientSecret,
		httpClient:     httpClient,
	}
}

// GetToken returns a valid access token, obtaining a new one if necessary
//
// Thread Safety:
// - Uses RLock for fast, concurrent reads when token is valid
// - Multiple goroutines can call this simultaneously with no contention
// - If new token needed, only one goroutine performs HTTP request
// - Safe to call from multiple goroutines concurrently
func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	// Fast path: Read token with read lock (allows concurrent access)
	tm.mu.RLock()
	token := tm.accessToken
	expiresAt := tm.expiresAt
	tm.mu.RUnlock()

	// Check if token is still valid (outside lock - uses local copies)
	if token != "" && time.Now().Before(expiresAt) {
		return token, nil
	}

	// Slow path: Token expired, obtain a new one
	// Note: ObtainToken uses write lock, preventing duplicate requests
	if err := tm.ObtainToken(ctx); err != nil {
		return "", err
	}

	// Read the newly obtained token
	tm.mu.RLock()
	token = tm.accessToken
	tm.mu.RUnlock()

	return token, nil
}

// ObtainToken gets a new access token using client credentials
//
// Thread Safety:
// - Acquires exclusive write lock (Lock) before making HTTP request
// - Blocks all other readers and writers during token acquisition
// - This is intentional: prevents multiple simultaneous HTTP requests
// - If multiple goroutines call this, only first one makes HTTP request
// - Others wait, then see the fresh token and return
// - Safe to call from multiple goroutines concurrently
func (tm *TokenManager) ObtainToken(ctx context.Context) error {
	// Acquire exclusive write lock
	// This blocks all other GetToken() calls until token acquisition completes
	// Prevents "thundering herd" problem
	tm.mu.Lock()
	defer tm.mu.Unlock() // Ensure unlock even if panic occurs

	// Prepare the token request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", tm.clientID)
	data.Set("client_secret", tm.clientSecret)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tm.tokenIssuerURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the token response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("received empty access token")
	}

	// Update the stored token
	tm.accessToken = tokenResp.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// IsTokenValid checks if the current token is still valid
func (tm *TokenManager) IsTokenValid() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tm.accessToken == "" {
		return false
	}

	return time.Now().Before(tm.expiresAt)
}

// GetExpiresAt returns when the current token expires
func (tm *TokenManager) GetExpiresAt() time.Time {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.expiresAt
}

// ClearToken clears the stored token
func (tm *TokenManager) ClearToken() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.accessToken = ""
	tm.expiresAt = time.Time{}
}

// GetTokenInfo returns information about the current token
func (tm *TokenManager) GetTokenInfo() (token string, expiresAt time.Time, isValid bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return tm.accessToken, tm.expiresAt, tm.isTokenValidLocked()
}

// isTokenValidLocked checks token validity without acquiring lock
// Must be called with read or write lock held
func (tm *TokenManager) isTokenValidLocked() bool {
	if tm.accessToken == "" {
		return false
	}
	return time.Now().Before(tm.expiresAt)
}

// GetRemainingTime returns the time until token expiry
func (tm *TokenManager) GetRemainingTime() time.Duration {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tm.accessToken == "" {
		return 0
	}

	remaining := time.Until(tm.expiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
