package aruba

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Arubacloud/sdk-go/internal/ports/interceptor"
	"github.com/Arubacloud/sdk-go/internal/ports/logger"
)

// LoggerType defines the supported logging strategies.
type LoggerType int

const (
	// LoggerNoLog disables all SDK logging.
	LoggerNoLog LoggerType = iota
	// LoggerNative uses the SDK's built-in standard library logger.
	LoggerNative
	// loggerCustom indicates a user-provided logger implementation is in use.
	// This is set automatically when WithCustomLogger is called.
	loggerCustom
)

// tokenManagerOptions holds internal configuration for the authentication subsystem.
type tokenManagerOptions struct {
	// tokenIssuerURL is the Aruba OAuth2 token endpoint URL.
	tokenIssuerURL string

	// clientID is the OAuth2 client ID.
	// It is used for Client Credentials authentication or to retrieve
	// credentials from a Vault repository.
	clientID string

	// clientSecret is the OAuth2 client secret.
	// It is mandatory if Vault is not used.
	// Mutually exclusive with vaultCredentialsRepositoryOptions.
	clientSecret string

	// tokenExpirationDriftSeconds defines the safety buffer subtracted from
	// a token's expiration time to prevent race conditions.
	// Ignored if no persistent repository proxy is configured.
	tokenExpirationDriftSeconds uint32

	// vaultCredentialsRepositoryOptions contains configuration for HashiCorp Vault.
	// Mutually exclusive with clientSecret.
	vaultCredentialsRepositoryOptions *vaultCredentialsRepositoryOptions

	// redisTokenRepositoryOptions contains configuration for a Redis token cache.
	// Mutually exclusive with fileTokenRepositoryOptions.
	redisTokenRepositoryOptions *redisTokenRepositoryOptions

	// fileTokenRepositoryOptions contains configuration for a file-system token cache.
	// Mutually exclusive with redisTokenRepositoryOptions.
	fileTokenRepositoryOptions *fileTokenRepositoryOptions
}

// vaultCredentialsRepositoryOptions configures the Vault connection.
type vaultCredentialsRepositoryOptions struct {
	// vaultURI is the address of the Vault server (e.g., "https://vault.example.com:8200").
	vaultURI  string
	kvMount   string
	kvPath    string
	namespace string
	rolePath  string
	roleID    string
	secretID  string
}

// redisTokenRepositoryOptions configures the Redis connection.
type redisTokenRepositoryOptions struct {
	// redisURI is the connection string for the Redis cluster.
	// Format: "redis://<user>:<pass>@localhost:6379/<db>"
	redisURI string
}

// fileTokenRepositoryOptions configures local file storage for tokens.
type fileTokenRepositoryOptions struct {
	// baseDir is the directory path where JSON token files will be stored.
	baseDir string
}

// userDefinedDependenciesOptions holds dependencies injected by the user.
type userDefinedDependenciesOptions struct {
	httpClient *http.Client
	logger     logger.Logger
	middleware interceptor.Interceptor
}

// Options is the configuration builder for the Aruba Cloud Client.
// It uses a fluent API pattern to chain configuration settings.
type Options struct {
	// baseURL is the root URL for Aruba REST API calls.
	baseURL string

	// loggerType indicates the logging strategy.
	loggerType LoggerType

	// tokenManager contains authentication-specific settings.
	tokenManager tokenManagerOptions

	// userDefinedDependencies contains injected components.
	userDefinedDependencies userDefinedDependenciesOptions
}

// NewOptions creates a new, empty configuration builder.
func NewOptions() *Options {
	return &Options{}
}

//
// Basic Options Helpers

// WithBaseURL overrides the default Aruba Cloud API URL.
func (o *Options) WithBaseURL(baseURL string) *Options {
	o.baseURL = baseURL
	return o
}

// WithTokenIssuerURL overrides the default OAuth2 token endpoint.
func (o *Options) WithTokenIssuerURL(tokenIssuerURL string) *Options {
	o.tokenManager.tokenIssuerURL = tokenIssuerURL
	return o
}

// WithClientID sets the OAuth2 Client ID.
func (o *Options) WithClientID(clientID string) *Options {
	o.tokenManager.clientID = clientID
	return o
}

// WithClientSecret sets the OAuth2 Client Secret.
// Side Effect: Removes any previously configured Vault options.
func (o *Options) WithClientSecret(clientSecret string) *Options {
	o.tokenManager.clientSecret = clientSecret
	// Ensure mutual exclusivity
	o.tokenManager.vaultCredentialsRepositoryOptions = nil
	return o
}

// WithClientCredentials is a helper to set both Client ID and Secret.
func (o *Options) WithClientCredentials(clientID string, clientSecret string) *Options {
	return o.WithClientID(clientID).WithClientSecret(clientSecret)
}

// WithLoggerType sets the logging strategy.
// Side Effect: Removes any custom logger previously set.
func (o *Options) WithLoggerType(loggerType LoggerType) *Options {
	o.loggerType = loggerType
	o.userDefinedDependencies.logger = nil
	return o
}

//
// Default Options Values and Helpers

const (
	defaultBaseURL        = "https://api.arubacloud.com"
	defaultLoggerType     = LoggerNoLog
	defaultTokenIssuerURL = "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token"
)

// WithDefaultBaseURL sets the URL to the production Aruba Cloud API.
func (o *Options) WithDefaultBaseURL() *Options {
	o.baseURL = defaultBaseURL
	return o
}

// WithDefaultTokenIssuerURL sets the URL to the production IDP.
func (o *Options) WithDefaultTokenIssuerURL() *Options {
	o.tokenManager.tokenIssuerURL = defaultTokenIssuerURL
	return o
}

// WithDefaultTokenManagerSchema configures standard Client Credentials auth
// without any persistent caching (Redis/File).
func (o *Options) WithDefaultTokenManagerSchema(clientID string, clientSecret string) *Options {
	o.tokenManager.fileTokenRepositoryOptions = nil
	o.tokenManager.redisTokenRepositoryOptions = nil
	return o.WithDefaultTokenIssuerURL().WithClientCredentials(clientID, clientSecret)
}

// WithDefaultLogger sets the logger type to "NoLog".
func (o *Options) WithDefaultLogger() *Options {
	o.loggerType = defaultLoggerType
	o.userDefinedDependencies.logger = nil
	return o
}

// DefaultOptions creates a ready-to-use configuration for the production environment
// using Client Credentials.
func DefaultOptions(clientID string, clientSecret string) *Options {
	return NewOptions().
		WithDefaultBaseURL().
		WithDefaultLogger().
		WithDefaultTokenManagerSchema(clientID, clientSecret)
}

//
// Logger Options Helpers

// WithNativeLogger enables the standard library logger.
func (o *Options) WithNativeLogger() *Options {
	o.loggerType = LoggerNative
	return o
}

// WithNoLogs disables logging.
func (o *Options) WithNoLogs() *Options {
	o.loggerType = LoggerNoLog
	return o
}

//
// Token Manager Options Helpers

const (
	stdRedisURI                           = "redis://admin:admin@localhost:6379/0"
	stdFileBaseDir                        = "/tmp/sdk-go"
	stdTokenExpirationDriftSeconds uint32 = 300
)

// WithVaultCredentialsRepository configures the SDK to fetch secrets from HashiCorp Vault.
// Side Effect: Clears any manually set Client Secret.
func (o *Options) WithVaultCredentialsRepository(
	vaultURI string,
	kvMount string,
	kvPath string,
	namespace string,
	rolePath string,
	roleID string,
	secretID string,
) *Options {
	o.tokenManager.vaultCredentialsRepositoryOptions = &vaultCredentialsRepositoryOptions{
		vaultURI:  vaultURI,
		kvMount:   kvMount,
		kvPath:    kvPath,
		namespace: namespace,
		rolePath:  rolePath,
		roleID:    roleID,
		secretID:  secretID,
	}

	o.tokenManager.clientSecret = ""
	return o
}

// WithTokenExpirationDriftSeconds sets the safety buffer for token expiration.
func (o *Options) WithTokenExpirationDriftSeconds(tokenExpirationDriftSeconds uint32) *Options {
	o.tokenManager.tokenExpirationDriftSeconds = tokenExpirationDriftSeconds
	return o
}

// WithStandardTokenExpirationDriftSeconds sets the drift to 300 seconds (5 minutes).
func (o *Options) WithStandardTokenExpirationDriftSeconds() *Options {
	return o.WithTokenExpirationDriftSeconds(stdTokenExpirationDriftSeconds)
}

// WithRedisTokenRepositoryFromURI configures a Redis cluster for token caching.
// Side Effect: Disables File Token Repository.
func (o *Options) WithRedisTokenRepositoryFromURI(redisURI string) *Options {
	o.tokenManager.redisTokenRepositoryOptions = &redisTokenRepositoryOptions{
		redisURI: redisURI,
	}
	o.tokenManager.fileTokenRepositoryOptions = nil
	return o
}

// WithRedisTokenRepositoryFromStandardURI configures Redis using localhost defaults.
func (o *Options) WithRedisTokenRepositoryFromStandardURI() *Options {
	return o.WithRedisTokenRepositoryFromURI(stdRedisURI)
}

// WithStandardRedisTokenRepository configures localhost Redis with standard drift settings.
func (o *Options) WithStandardRedisTokenRepository() *Options {
	return o.WithRedisTokenRepositoryFromStandardURI().WithStandardTokenExpirationDriftSeconds()
}

// WithFileTokenRepositoryFromBaseDir configures a directory for storing token files.
// Side Effect: Disables Redis Token Repository.
func (o *Options) WithFileTokenRepositoryFromBaseDir(baseDir string) *Options {
	o.tokenManager.fileTokenRepositoryOptions = &fileTokenRepositoryOptions{
		baseDir: baseDir,
	}
	o.tokenManager.redisTokenRepositoryOptions = nil
	return o
}

// WithFileTokenRepositoryFromStandardBaseDir configures file storage in /tmp/sdk-go.
func (o *Options) WithFileTokenRepositoryFromStandardBaseDir() *Options {
	return o.WithFileTokenRepositoryFromBaseDir(stdFileBaseDir)
}

// WithStandardFileTokenRepository configures /tmp storage with standard drift settings.
func (o *Options) WithStandardFileTokenRepository() *Options {
	return o.WithFileTokenRepositoryFromStandardBaseDir().WithStandardTokenExpirationDriftSeconds()
}

//
// User-Defined Dependency Options Helpers

// WithCustomHTTPClient allows injecting a pre-configured *http.Client.
func (o *Options) WithCustomHTTPClient(client *http.Client) *Options {
	o.userDefinedDependencies.httpClient = client
	return o
}

// WithCustomLogger allows injecting a custom logger.Logger implementation.
func (o *Options) WithCustomLogger(logger logger.Logger) *Options {
	o.loggerType = loggerCustom
	o.userDefinedDependencies.logger = logger
	return o
}

// WithCustomMiddleware allows injecting a custom interceptor.Interceptor.
func (o *Options) WithCustomMiddleware(middleware interceptor.Interceptor) *Options {
	o.userDefinedDependencies.middleware = middleware
	return o
}

//
// Validation

// validate checks the root Options structure.
func (o *Options) validate() error {
	var errs []error

	if err := validateURL(o.baseURL, "base API URL"); err != nil {
		errs = append(errs, err)
	}

	if err := o.tokenManager.validate(); err != nil {
		errs = append(errs, err)
	}

	if o.loggerType == loggerCustom && o.userDefinedDependencies.logger == nil {
		errs = append(
			errs,
			errors.New(
				"logger type is set to 'Custom' but no custom logger implementation was provided via WithCustomLogger()",
			),
		)
	}

	return errors.Join(errs...)
}

// validate checks the Token Manager configuration, including mutual exclusions.
func (tm *tokenManagerOptions) validate() error {
	var errs []error

	//
	// Basic Fields

	if err := validateURL(tm.tokenIssuerURL, "token issuer URL"); err != nil {
		errs = append(errs, err)
	}

	if strings.TrimSpace(tm.clientID) == "" {
		errs = append(errs, errors.New("client ID is required"))
	}

	//
	// Credentials Mutual Exclusion & Validity

	hasClientSecret := strings.TrimSpace(tm.clientSecret) != ""
	hasVault := tm.vaultCredentialsRepositoryOptions != nil

	if hasClientSecret && hasVault {
		errs = append(
			errs,
			errors.New(
				"configuration conflict: cannot use both Client Secret and Vault Repository for credentials; please choose one",
			),
		)
	} else if !hasClientSecret && !hasVault {
		errs = append(
			errs,
			errors.New(
				"missing credentials: must provide either a Client Secret or Vault Repository configuration",
			),
		)
	} else if hasVault {
		if err := tm.vaultCredentialsRepositoryOptions.validate(); err != nil {
			errs = append(errs, fmt.Errorf("vault configuration error: %w", err))
		}
	}

	//
	// Token Cache Mutual Exclusion & Validity

	// Note: It is Valid for both Redis and File to be nil: implies no
	// persistence/caching.
	hasRedis := tm.redisTokenRepositoryOptions != nil
	hasFile := tm.fileTokenRepositoryOptions != nil

	if hasRedis && hasFile {
		errs = append(
			errs,
			errors.New(
				"configuration conflict: cannot use both Redis and File System for token caching; please choose one",
			),
		)
	}

	if hasRedis {
		if err := tm.redisTokenRepositoryOptions.validate(); err != nil {
			errs = append(errs, fmt.Errorf("redis configuration error: %w", err))
		}
	}

	if hasFile {
		if err := tm.fileTokenRepositoryOptions.validate(); err != nil {
			errs = append(errs, fmt.Errorf("file repository configuration error: %w", err))
		}
	}

	return errors.Join(errs...)
}

// validate checks Vault specific options.
func (v *vaultCredentialsRepositoryOptions) validate() error {
	var errs []error

	if err := validateURL(v.vaultURI, "vault URI"); err != nil {
		errs = append(errs, err)
	}

	if strings.TrimSpace(v.roleID) == "" {
		errs = append(errs, errors.New("vault Role ID is required"))
	}

	if strings.TrimSpace(v.secretID) == "" {
		errs = append(errs, errors.New("vault Secret ID is required"))
	}

	if strings.TrimSpace(v.kvMount) == "" {
		errs = append(errs, errors.New("vault KV Mount path is required"))
	}

	if strings.TrimSpace(v.kvPath) == "" {
		errs = append(errs, errors.New("vault KV Secret path is required"))
	}

	return errors.Join(errs...)
}

// validate checks Redis specific options.
func (r *redisTokenRepositoryOptions) validate() error {
	u, err := url.ParseRequestURI(r.redisURI)
	if err != nil {
		return fmt.Errorf("invalid redis URI format: %w", err)
	}

	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return fmt.Errorf("invalid redis URI scheme '%s': must be 'redis://' or 'rediss://'", u.Scheme)
	}

	if u.Host == "" {
		return errors.New("invalid redis URI: missing host address")
	}

	return nil
}

// validate checks File Repository options.
func (f *fileTokenRepositoryOptions) validate() error {
	// We rely on string length.
	// Note: We do not check if the directory exists here (os.Stat) because
	// the application might have permissions to create it later.
	// We only validate that the configuration string is sensible.
	path := strings.TrimSpace(f.baseDir)
	if path == "" {
		return errors.New("base directory path cannot be empty")
	}

	// Simple check for potentially dangerous or invalid paths (optional)
	// Example: prevents root directory usage if desired, though usually ignored in SDKs.
	// if path == "/" { return errors.New("cannot use root directory") }

	return nil
}

//
// Validation Helper Functions

// validateURL parses a string to ensure it is a valid absolute URL (HTTP/HTTPS).
func validateURL(rawURL, fieldName string) error {
	if strings.TrimSpace(rawURL) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("%s is malformed: %w", fieldName, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s has invalid scheme '%s': must be http or https", fieldName, u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("%s is missing a host", fieldName)
	}

	return nil
}
