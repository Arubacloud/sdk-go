# OAuth2 Client Credentials Flow - Implementation Details

## Overview

This SDK implements OAuth2 **Client Credentials Grant** flow for machine-to-machine authentication. This flow is designed for server-to-server communication where no user interaction is required.

## Key Characteristics

### No Refresh Token
Unlike the Authorization Code flow, client credentials flow **does not** provide a refresh token. Instead:
- A new access token is requested using the same `client_id` and `client_secret`
- The `refresh_expires_in` field in the response is always `0`
- There is no `refresh_token` field in the response

### Token Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│  1. Request Token                                            │
│     POST /oauth2/token                                       │
│     grant_type=client_credentials                           │
│     client_id=xxx                                            │
│     client_secret=yyy                                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Receive Token Response                                   │
│     {                                                        │
│       "access_token": "eyJ...",                             │
│       "token_type": "Bearer",                               │
│       "expires_in": 3600,                                   │
│       "refresh_expires_in": 0,  ← Always 0                 │
│       "scope": "email"                                      │
│     }                                                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Use Token in API Calls                                   │
│     Authorization: Bearer eyJ...                             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Token About to Expire (< 5 min remaining)                │
│     Automatically request new token                          │
│     (using same client_id and client_secret)                 │
└─────────────────────────────────────────────────────────────┘
```

## Aruba Cloud Token Response

### Actual Response Format

```json
{
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIs...",
    "expires_in": 3600,
    "refresh_expires_in": 0,
    "token_type": "Bearer",
    "not-before-policy": 0,
    "scope": "email"
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | JWT Bearer token for API authentication |
| `expires_in` | int | Token lifetime in seconds (typically 3600 = 1 hour) |
| `refresh_expires_in` | int | Always 0 for client credentials flow |
| `token_type` | string | Always "Bearer" |
| `not-before-policy` | int | Timestamp before which token is not valid (usually 0) |
| `scope` | string | Granted scopes (e.g., "email") |

### JWT Token Claims

The access token is a JWT that contains claims like:

```json
{
  "exp": 1761927735,           // Expiration timestamp
  "iat": 1761924135,           // Issued at timestamp
  "jti": "b2dad6f5-...",       // JWT ID
  "iss": "https://login.aruba.it/auth/realms/...",
  "sub": "861f9a93-...",       // Subject (client ID)
  "typ": "Bearer",
  "azp": "cmp-7ffe0eee-...",   // Authorized party
  "scope": "email",
  "company": "ARU",
  "tenant": "ARU",
  "client_id": "cmp-7ffe0eee-..."
}
```

## SDK Implementation

### TokenManager

The `TokenManager` handles all token operations:

```go
type TokenManager struct {
    tokenIssuerURL     string        // OAuth2 token endpoint
    clientID           string        // Your client ID
    clientSecret       string        // Your client secret
    httpClient         *http.Client  // HTTP client for requests
    tokenRefreshBuffer time.Duration // When to refresh (default: 5 min before expiry)
    
    // Protected by mutex
    mu          sync.RWMutex
    accessToken string        // Current JWT token
    expiresAt   time.Time     // When token expires
}
```

### Automatic Token Management

The SDK automatically:

1. **Obtains initial token** when `NewClient()` is called
2. **Caches the token** for reuse
3. **Checks expiration** before each API call
4. **Refreshes proactively** 5 minutes before expiry
5. **Thread-safe** access for concurrent requests

### Usage Example

```go
// Configure SDK
config := &client.Config{
    BaseURL:        "https://api.arubacloud.com",
    TokenIssuerURL: "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token",
    ClientID:       "cmp-7ffe0eee-e45b-41c5-864b-3b178eeacb2d",
    ClientSecret:   "your-secret-here",
}

// Create SDK client (automatically gets token)
sdk, err := client.NewClient(config)

// Use SDK - token automatically managed
result, err := sdk.Network().ListNetworks(ctx)
```

### Manual Token Access

If you need direct access to the token:

```go
// Get current valid token (auto-refreshes if needed)
token, err := sdk.GetToken(ctx)

// Get token info
token, expiresAt, isValid := sdk.GetTokenInfo()

// Check remaining time
remaining := sdk.GetRemainingTime()
```

## Token Refresh Strategy

### Refresh Buffer

The SDK uses a configurable **refresh buffer** (default: 5 minutes):

```
Token Lifetime: 3600 seconds (1 hour)
Refresh Buffer: 300 seconds (5 minutes)
Effective Use : 3300 seconds (55 minutes)

Timeline:
0s ─────────── 3300s ──────── 3600s
│              │               │
Obtained       Refresh         Expires
               Triggered
```

### Benefits

1. **Prevents expired token errors** in API calls
2. **Allows time for refresh** even under heavy load
3. **Configurable** based on your needs

```go
config := &client.Config{
    // ... other fields ...
    TokenRefreshBuffer: 10 * time.Minute, // Refresh 10 min before expiry
}
```

## Error Handling

### Token Request Failures

```go
sdk, err := client.NewClient(config)
if err != nil {
    // Could be:
    // - Invalid credentials
    // - Network error
    // - Token issuer unavailable
    log.Fatalf("Failed to authenticate: %v", err)
}
```

### Token Refresh Failures

If token refresh fails during an API call, the SDK will:
1. Return the error to the caller
2. Keep the old (expired) token
3. Retry refresh on next API call

## Security Considerations

### Credentials Storage

⚠️ **Never hardcode credentials in source code**

Use environment variables or secure secret management:

```go
config := &client.Config{
    BaseURL:        os.Getenv("ARUBA_BASE_URL"),
    TokenIssuerURL: os.Getenv("ARUBA_TOKEN_URL"),
    ClientID:       os.Getenv("ARUBA_CLIENT_ID"),
    ClientSecret:   os.Getenv("ARUBA_CLIENT_SECRET"),
}
```

### Token Storage

- Tokens are stored **in memory only**
- Never persisted to disk
- Automatically cleared when client is destroyed
- Use `ClearToken()` to manually clear if needed

## Differences from Other OAuth2 Flows

### vs Authorization Code Flow

| Feature | Client Credentials | Authorization Code |
|---------|-------------------|-------------------|
| User interaction | ❌ None | ✅ Required |
| Refresh token | ❌ No | ✅ Yes |
| Use case | Server-to-server | User applications |
| Token refresh | Re-request with credentials | Use refresh_token |

### vs Password Grant

| Feature | Client Credentials | Password Grant |
|---------|-------------------|----------------|
| User credentials | ❌ Not needed | ✅ Required |
| Security | ✅ Better | ⚠️ Less secure |
| Recommended | ✅ Yes | ❌ Deprecated |

## Testing

### Mock Token Server

```go
tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    resp := client.TokenResponse{
        AccessToken: "test-token",
        TokenType:   "Bearer",
        ExpiresIn:   3600,
        RefreshExpiresIn: 0,
    }
    json.NewEncoder(w).Encode(resp)
}))
defer tokenServer.Close()

config := &client.Config{
    TokenIssuerURL: tokenServer.URL,
    // ... other config ...
}
```

## Troubleshooting

### Common Issues

1. **"invalid_client" error**
   - Check client_id and client_secret
   - Verify credentials are URL-encoded properly

2. **Token expires immediately**
   - Check system clock synchronization
   - Verify token_issuer_url is correct

3. **Frequent token refreshes**
   - Normal if TokenRefreshBuffer is large
   - Reduce buffer time if needed

### Debug Logging

```go
// Get token info for debugging
token, expiresAt, isValid := sdk.GetTokenInfo()
log.Printf("Token valid: %v, Expires: %v, Remaining: %v",
    isValid, expiresAt, sdk.GetRemainingTime())
```

## References

- [OAuth 2.0 RFC 6749 - Client Credentials](https://datatracker.ietf.org/doc/html/rfc6749#section-4.4)
- [JWT RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)
- [Aruba Cloud API Documentation](https://www.arubacloud.com/)
