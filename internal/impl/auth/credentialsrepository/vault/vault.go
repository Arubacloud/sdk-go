package vault

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	vaultapi "github.com/hashicorp/vault/api"
)

// CredentialsRepository implements the auth.CredentialsRepository interface.
// It is responsible for fetching credentials from a HashiCorp Vault backend.
type CredentialsRepository struct {
	// Implementation details would go here
	client     VaultClient
	kvMount    string
	kvPath     string
	namespace  string
	rolePath   string
	roleID     string
	secretID   string
	tokenExist bool
	renewable  bool
	ttl        time.Duration
	mu         sync.RWMutex
}

var _ auth.CredentialsRepository = (*CredentialsRepository)(nil)

// VaultClient is the main interface abstracting the *vaultapi.Client.
// This allows the business logic (CredentialsRepository) to be decoupled from
// the specific Vault SDK implementation, enabling dependency injection and testing.
type VaultClient interface {
	Logical() LogicalAPI
	SetToken(token string)
	KVv2(mount string) KvAPI
	SetNamespace(namespace string)
	Auth() AuthAPI
}

// VaultClientAdapter wraps *vaultapi.Client to conform to the VaultClient interface.
// This is the Adapter pattern in action.
type VaultClientAdapter struct {
	c *vaultapi.Client
}

// Adapter types for the various sub-components of the Vault API.
type logicalAPIAdapter struct {
	l *vaultapi.Logical
}

type kvAPIAdapter struct {
	kv *vaultapi.KVv2
}

type authAPIAdapter struct {
	auth *vaultapi.Auth
}

type authTokenAPIAdapter struct {
	token *vaultapi.TokenAuth
}

// LogicalAPI defines the required methods for interacting with Vault's logical backend.
// This adheres to the Interface Segregation Principle (ISP) by exposing only
// the necessary methods (e.g., Write).
type LogicalAPI interface {
	Write(path string, data map[string]any) (*vaultapi.Secret, error)
}

// KvAPI defines the required methods for interacting with Vault's KVv2 secrets engine.
type KvAPI interface {
	Get(ctx context.Context, path string) (*vaultapi.KVSecret, error)
}

// AuthAPI defines the required methods for interacting with Vault's authentication methods.
type AuthAPI interface {
	Token() AuthTokenAPI
}

// AuthTokenAPI defines the methods for managing tokens (e.g., renewal).
type AuthTokenAPI interface {
	RenewSelfWithContext(ctx context.Context, increment int) (*vaultapi.Secret, error)
}

// NewVaultClientAdapter is the constructor for the VaultClientAdapter.
func NewVaultClientAdapter(c *vaultapi.Client) *VaultClientAdapter {
	return &VaultClientAdapter{c: c}
}

// NewCredentialsRepository creates a new CredentialsRepository that fetches
// credentials from a Vault backend.
func NewCredentialsRepository(v VaultClient, cfg restclient.VaultConfig) *CredentialsRepository {
	return &CredentialsRepository{
		client:     v,
		kvMount:    cfg.KVMount,
		kvPath:     cfg.KVPath,
		namespace:  cfg.Namespace,
		rolePath:   cfg.RolePath,
		roleID:     cfg.RoleID,
		secretID:   cfg.SecretID,
		tokenExist: false,
		renewable:  false,
		ttl:        0}
}

// FetchCredentials retrieves the Client ID and Secret from Vault.
func (r *CredentialsRepository) FetchCredentials(ctx context.Context) (*auth.Credentials, error) {
	r.mu.RLock()

	if !r.tokenExist {
		err := r.loginWithAppRole(ctx)
		if err != nil {
			return nil, auth.ErrAuthenticationFailed
		}
	}
	r.mu.RUnlock()

	// Fetch the secret from Vault KV
	secret, err := r.client.KVv2(r.kvMount).Get(ctx, r.kvPath)
	if err != nil {
		return nil, auth.ErrCredentialsNotFound
	}

	// Extract credentials from the Vault secret
	return getCredentialsFromVaultSecret(secret)
}

// loginWithAppRole performs AppRole authentication to obtain a Vault token
func (r *CredentialsRepository) loginWithAppRole(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// another goroutine already logged in
	if r.tokenExist {
		return nil
	}

	// Set the namespace for the Vault client
	if r.namespace != "" {
		r.client.SetNamespace(r.namespace)
	}

	payload := map[string]any{
		"role_id":   r.roleID,
		"secret_id": r.secretID,
	}
	token, err := r.client.Logical().Write(r.rolePath, payload)
	if err != nil {
		return err
	}

	if token == nil || token.Auth == nil || token.Auth.ClientToken == "" {
		return fmt.Errorf("vault approle login response missing token")
	}

	clientToken := token.Auth.ClientToken

	r.client.SetToken(clientToken)
	r.tokenExist = true

	r.ttl = time.Duration(token.Auth.LeaseDuration) * time.Second
	r.renewable = token.Auth.Renewable

	if r.renewable {
		go r.renewTokenBeforeExpiration(ctx)
	}
	return nil
}

// getCredentialsFromVaultSecret extracts Client ID and Secret from a Vault KV secret.
// It assumes the secret data contains "client_id" and "client_secret" keys.
func getCredentialsFromVaultSecret(secret *vaultapi.KVSecret) (*auth.Credentials, error) {
	get := func(key string) (string, error) {
		v, ok := secret.Data[key].(string)
		if !ok {
			return "", auth.ErrCredentialsNotFound
		}
		return v, nil
	}

	clientID, err := get("client_id")
	if err != nil {
		return nil, err
	}

	clientSecret, err := get("client_secret")
	if err != nil {
		return nil, err
	}

	creds := &auth.Credentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	return creds, nil
}

// renewTokenBeforeExpiration attempts to renew the current Vault token if the repository
// is configured to use a renewable token.
func (r *CredentialsRepository) renewTokenBeforeExpiration(ctx context.Context) {
	renewClient := r.client.Auth().Token()

	renewInterval := r.ttl / 2
	if renewInterval <= 0 {
		renewInterval = 30 * time.Second
	}

	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return // stop goroutine immediately
	case <-ticker.C:
		sec, err := renewClient.RenewSelfWithContext(ctx, int(r.ttl/time.Second))
		if err != nil {
			// If renewal fails, the token may be expired → stop renew
			// You can log the error and break or continue depending on policy
			// log.Errorf("Vault token renewal failed: %v", err)
			r.mu.Lock()
			r.tokenExist = false
			r.mu.Unlock()
			return
		}
		if sec != nil && sec.Auth != nil {
			r.mu.Lock()
			r.ttl = time.Duration(sec.Auth.LeaseDuration) * time.Second
			r.renewable = sec.Auth.Renewable
			r.tokenExist = true
			r.client.SetToken(sec.Auth.ClientToken)
			r.mu.Unlock()
		}
		// If token is no longer renewable, end loop
		if !r.renewable {
			return
		}
	}
}

// implement VaultClientAdapter methods
func (v *VaultClientAdapter) SetNamespace(namespace string) {
	v.c.SetNamespace(namespace)
}

func (v *VaultClientAdapter) Auth() AuthAPI {
	return &authAPIAdapter{auth: v.c.Auth()}
}

func (v *VaultClientAdapter) SetToken(token string) {
	v.c.SetToken(token)
}

func (v *VaultClientAdapter) Logical() LogicalAPI {
	return &logicalAPIAdapter{l: v.c.Logical()}
}
func (v *VaultClientAdapter) KVv2(mount string) KvAPI {
	return &kvAPIAdapter{kv: v.c.KVv2(mount)}
}

// implement AuthAPIAdapter methods
func (a *authAPIAdapter) Token() AuthTokenAPI {
	return &authTokenAPIAdapter{token: a.auth.Token()}
}

// implement AuthTokenAPIAdapter methods
func (a *authTokenAPIAdapter) RenewSelfWithContext(ctx context.Context, increment int) (*vaultapi.Secret, error) {
	return a.token.RenewSelfWithContext(ctx, increment)
}

// implement KVAPIAdapter methods
func (k *kvAPIAdapter) Get(ctx context.Context, path string) (*vaultapi.KVSecret, error) {
	return k.kv.Get(ctx, path)
}

// implement LogicalAPIAdapter methods
func (la *logicalAPIAdapter) Write(path string, data map[string]any) (*vaultapi.Secret, error) {
	return la.l.Write(path, data)
}
