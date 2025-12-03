package aruba

import (
	"errors"
	"net/http"

	"github.com/Arubacloud/sdk-go/internal/ports/interceptor"
	"github.com/Arubacloud/sdk-go/internal/ports/logger"
)

type options struct {
	// baseURL is the base URL for all calls to Aruba REST API
	baseURL string

	// loggerType indicates the type of log to use
	loggerType LoggerType

	//
	// Authentication and Token Management

	// authOptions
	authOptions *authOptions

	//
	// User-built dependencies

	// httpClient
	httpClient *http.Client

	// logger
	logger logger.Logger

	// middleware
	middleware interceptor.Interceptor
}

func NewOprions() *options {
	return &options{}
}

// Default values
const (
	defaultBaseURL        = "https://api.arubacloud.com"
	defaultTokenIssuerURL = "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token"
	defaultRedisURI       = "redis://admin:admin@localhost:6379/0"
	defaultFileBaseDir    = "/tmp/sdk-go"
)

func DefaultOptions() *options {
	return &options{}
}

func (o *options) parse() error {
	return errors.New("not implemented")
}

//
// Base config

func (o *options) WithBaseURL(baseURL string) *options {
	o.baseURL = baseURL
	return o
}

func (o *options) WithDefaultBaseURL() *options {
	o.baseURL = defaultBaseURL

	return o
}

//
// Logginig

type LoggerType int

const (
	LoggerNoLog LoggerType = iota
	LoggerNative
	loggerCustom // this should be kept private
)

func (o *options) WithLoggerType(loggerType LoggerType) *options {
	o.loggerType = loggerType

	return o
}

func (o *options) WithNativeLogger() *options {
	o.loggerType = LoggerNative

	return o
}

func (o *options) WithNoLogs() *options {
	o.loggerType = LoggerNoLog

	return o
}

func (o *options) WithCustomLogger(logger logger.Logger) *options {
	o.loggerType = loggerCustom

	o.logger = logger

	return o
}

//
// Authentication

type authOptions struct {
	// tokenIssuerURL is the URL Aruba OAuth2 token endpoint
	tokenIssuerURL string

	// clientID is the OAuth2 client id used in the client credentials
	// authentication schema.
	//
	// However, it is also used to retrieve the proper credentials from a Vault
	// credential repository.
	clientID string

	// clientSecret is the OAuth2 client secret used in the client credentials
	// authentication schema.
	//
	// It mandatory to be set if no Vault credential repository is used.
	//
	// It is mutualy excludent with vaultCredentialsRepositoryOptions.
	clientSecret string

	// vaultCredentialsRepositoryOptions contains all the configuration parameters
	// necessary to use Vault as credentials repository.
	//
	// It is mutualy excludent with clientSecret.
	vaultCredentialsRepositoryOptions *vaultCredentialsRepositoryOptions

	// redisTokenRepositoryOptions contains all the configuration parameters
	// necessary to use a Redis cluster as token repository.
	//
	// If it is set, so a RedisTokenRepository will be added to the
	// TokenManager chain.
	//
	// It is mutualy excludent with fileTokenRepositoryOptions.
	redisTokenRepositoryOptions *redisTokenRepositoryOptions

	// fileTokenRepositoryOptions contains all the configuration parameters
	// necessary to use the local file system as token repository.
	//
	// If it is set, so a FileTokenRepository will be added to the
	// TokenManager chain.
	//
	// It is mutualy excludent with redisTokenRepositoryOptions.
	fileTokenRepositoryOptions *fileTokenRepositoryOptions
}

type vaultCredentialsRepositoryOptions struct {
	//address:port
	vaultURI  string
	kvMount   string
	kvPath    string
	namespace string
	rolePath  string
	roleID    string
	secretID  string
}

type redisTokenRepositoryOptions struct {
	// redisURI hold the URI to connect to a Redis cluster.
	//
	// Format: "redis://<user>:<pass>@localhost:6379/<db>"
	redisURI string
}

type fileTokenRepositoryOptions struct {
	// baseDir is the path to the directory where stored token json files are
	// located.
	baseDir string
}
