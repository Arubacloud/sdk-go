# SDK Configuration Options

The Aruba Cloud SDK for Go is configured using a fluent API that allows you to chain option setters to build a
client configuration. This guide details all available options, grouped by topic.

## Getting Started

<p>The easiest way to configure the client is by using <code>DefaultOptions</code>, which provides a
production-ready setup for the most common use case (Client Credentials authentication).</p>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>DefaultOptions(clientID, clientSecret)</code></td>
      <td>Creates a standard configuration for the production environment using your API Client ID and Secret.</td>
      <td>This is the recommended starting point for most applications. It sets the production API and token URLs
      and configures authentication.</td>
    </tr>
    <tr>
      <td><code>NewOptions()</code></td>
      <td>Creates a new, empty configuration builder. You must manually configure all required settings,
      including authentication.</td>
      <td>Use this if you need complete control over the configuration from scratch.</td>
    </tr>
  </tbody>
</table>

<h2>General Configuration</h2>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithBaseURL(baseURL)</code></td>
      <td>Overrides the default Aruba Cloud API URL.</td>
      <td>The default is <code>https://api.arubacloud.com</code>. You should only change this if you need to
      target a different environment.</td>
    </tr>
    <tr>
      <td><code>WithDefaultBaseURL()</code></td>
      <td>A helper that sets the API URL to the production default.</td>
      <td>Called by <code>DefaultOptions()</code>.</td>
    </tr>
  </tbody>
</table>

<h2>Authentication</h2>

<p>Authentication is a critical part of the configuration. You must choose <b>one</b> of the following methods.</p>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithClientCredentials(clientID, clientSecret)</code></td>
      <td><b>(Recommended)</b> Configures the SDK to use the OAuth2 Client Credentials flow. The SDK will
      automatically manage fetching and renewing the access token.</td>
      <td><b>Mutual Exclusion</b>: Cannot be used with <code>WithToken()</code> or
      <code>WithVaultCredentialsRepository()</code>. This is the standard and most secure method for
      service-to-service authentication.</td>
    </tr>
    <tr>
      <td><code>WithToken(token)</code></td>
      <td>Configures the SDK to use a static, pre-existing OAuth2 access token.</td>
      <td><b>Mutual Exclusion</b>: Cannot be used with any other authentication or token issuer option.<br/>
      <b>Warning</b>: The SDK will <b>not</b> renew this token. Once it expires, all API calls will fail. This method
      is only recommended for short-lived clients or simple, atomic scripts.</td>
    </tr>
    <tr>
      <td><code>WithVaultCredentialsRepository(...)</code></td>
      <td>Configures the SDK to fetch the <code>clientID</code> and <code>clientSecret</code> from a HashiCorp
      Vault instance. The SDK will then use these credentials to manage the access token.</td>
      <td><b>Mutual Exclusion</b>: Cannot be used with <code>WithClientCredentials()</code> or <code>WithToken()</code>.<br/>
      <b>Parameters</b>: <code>vaultURI</code>, <code>kvMount</code>, <code>kvPath</code>, <code>namespace</code>,
      <code>rolePath</code>, <code>roleID</code>, <code>secretID</code>.</td>
    </tr>
    <tr>
      <td><code>WithTokenIssuerURL(url)</code></td>
      <td>Overrides the default URL for the OAuth2 token endpoint.</td>
      <td>The default is <code>https://login.aruba.it/...</code>. You should only change this if you need to
      target a different authentication endpoint.</td>
    </tr>
    <tr>
      <td><code>WithSecurityScopes(scopes ...)</code></td>
      <td>Sets the security scopes to be claimed during authentication.</td>
      <td>Replaces any previously set scopes.</td>
    </tr>
    <tr>
      <td><code>WithAdditionalSecurityScopes(scopes ...)</code></td>
      <td>Appends additional security scopes to the existing list.</td>
      <td>Does not override existing scopes.</td>
    </tr>
  </tbody>
</table>

<h2>Token Caching (Optional)</h2>

<p>For improved performance and resilience, the SDK can cache the access token to an external store. This is highly
recommended for production applications to avoid fetching a new token on every startup.</p>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithRedisTokenRepositoryFromURI(redisURI)</code></td>
      <td>Configures a Redis instance as a persistent cache for the access token.</td>
      <td><b>Mutual Exclusion</b>: Cannot be used with <code>WithFileTokenRepository...</code> options.<br/>
      The URI format is <code>redis://&lt;user&gt;:&lt;pass&gt;@&lt;host&gt;:&lt;port&gt;/&lt;db&gt;</code>.</td>
    </tr>
    <tr>
      <td><code>WithFileTokenRepositoryFromBaseDir(baseDir)</code></td>
      <td>Configures a local directory to store the access token in a JSON file.</td>
      <td><b>Mutual Exclusion</b>: Cannot be used with <code>WithRedisTokenRepository...</code> options. The SDK process
      must have read/write permissions to the specified directory.</td>
    </tr>
    <tr>
      <td><code>WithTokenExpirationDriftSeconds(seconds)</code></td>
      <td>Sets a safety buffer (in seconds) to treat a token as expired before it actually does.</td>
      <td>This prevents race conditions where the SDK uses a token that expires just before the API request
      completes. The default is 300 seconds (5 minutes). This option has no effect if a caching
      mechanism (Redis or File) is not configured.</td>
    </tr>
    <tr>
      <td><code>WithStandardRedisTokenRepository()</code></td>
      <td>Helper to configure Redis caching on <code>localhost:6379</code> and sets the standard expiration
      drift (300s).</td>
      <td>A convenient shortcut for local development.</td>
    </tr>
    <tr>
      <td><code>WithStandardFileTokenRepository()</code></td>
      <td>Helper to configure file caching in <code>/tmp/sdk-go</code> and sets the standard expiration
      drift (300s).</td>
      <td>A convenient shortcut for local development.</td>
    </tr>
  </tbody>
</table>

<h2>Logging</h2>

<p>Logging is disabled by default.</p>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithNoLogs()</code></td>
      <td>Disables all SDK logging.</td>
      <td>This is the default behavior.</td>
    </tr>
    <tr>
      <td><code>WithNativeLogger()</code></td>
      <td>Enables the SDK's built-in logger, which uses the standard <code>log</code> package.</td>
      <td>Useful for basic debugging.</td>
    </tr>
    <tr>
      <td><code>WithLoggerType(loggerType)</code></td>
      <td>A more general function to set the logger type.</td>
      <td><code>loggerType</code> can be <code>LoggerNoLog</code> or <code>LoggerNative</code>.</td>
    </tr>
    <tr>
      <td><code>WithCustomLogger(logger)</code></td>
      <td>Injects a custom logger that conforms to the <code>ports/logger.Logger</code> interface.</td>
      <td>Allows integration with your application's logging framework (e.g., Logrus, Zap). Setting this
      automatically sets the logger type to <code>loggerCustom</code>.</td>
    </tr>
  </tbody>
</table>

<h2>Advanced / Custom Dependencies</h2>

<p>These options are for advanced use cases where you need to inject your own custom components into the SDK's
workflow.</p>

<table>
  <thead>
    <tr>
      <th>Option Setter</th>
      <th>Description</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithCustomHTTPClient(client)</code></td>
      <td>Injects a pre-configured <code>*http.Client</code>.</td>
      <td>Useful for setting custom timeouts, transport, or other HTTP-level configurations.</td>
    </tr>
    <tr>
      <td><code>WithCustomMiddleware(middleware)</code></td>
      <td>Injects a custom middleware that conforms to the <code>ports/interceptor.Interceptor</code> interface.</td>
      <td>Allows you to add custom logic (like request/response logging or tracing) into the SDK's HTTP
      call chain. The SDK's authentication middleware will be automatically bound to the end of your
      custom middleware chain.</td>
    </tr>
  </tbody>
</table>
