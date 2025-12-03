package vault

import (
	"context"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	vaultapi "github.com/hashicorp/vault/api"
)

// CredentialsRepository implements the auth.CredentialsRepository interface.
// It is responsible for fetching credentials from a HashiCorp Vault backend.
type CredentialsRepository struct {
	// Implementation details would go here
	client    VaultClient
	kvMount   string
	kvPath    string
	namespace string
	rolePath  string
	roleID    string
	secretID  string
	renewable bool
	ttl       time.Duration
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
		client:    v,
		kvMount:   cfg.KVMount,
		kvPath:    cfg.KVPath,
		namespace: cfg.Namespace,
		rolePath:  cfg.RolePath,
		roleID:    cfg.RoleID,
		secretID:  cfg.SecretID,
		renewable: false,
		ttl:       0}
}

// FetchCredentials retrieves the Client ID and Secret from Vault.
func (r *CredentialsRepository) FetchCredentials(ctx context.Context) (*auth.Credentials, error) {

	// Implementation to fetch credentials from Vault would go here
	return nil, nil
}

// renewTokenIfNeeded attempts to renew the current Vault token if the repository
// is configured to use a renewable token.
func (r *CredentialsRepository) renewTokenIfNeeded(ctx context.Context) error {
	if r.renewable {
		_, err := r.client.Auth().Token().RenewSelfWithContext(ctx, int(r.ttl.Seconds()))
		if err != nil {
			return err
		}
	}
	return nil
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
