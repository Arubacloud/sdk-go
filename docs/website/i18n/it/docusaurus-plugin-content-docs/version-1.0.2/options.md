# Opzioni di Configurazione SDK

L'SDK Aruba Cloud per Go è configurato mediante un'API fluente che consente di concatenare i setter di opzioni per costruire una configurazione client. Questa guida dettaglia tutte le opzioni disponibili, raggruppate per argomento.

## Iniziare

<p>Il modo più semplice per configurare il client è utilizzare <code>DefaultOptions</code>, che fornisce una configurazione pronta per la produzione per il caso d'uso più comune (autenticazione Client Credentials).</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>DefaultOptions(clientID, clientSecret)</code></td>
      <td>Crea una configurazione standard per l'ambiente di produzione utilizzando il Client ID e il Secret API.</td>
      <td>Questo è il punto di partenza consigliato per la maggior parte delle applicazioni. Imposta l'API di produzione e gli URL dei token e configura l'autenticazione.</td>
    </tr>
    <tr>
      <td><code>NewOptions()</code></td>
      <td>Crea un nuovo builder di configurazione vuoto. È necessario configurare manualmente tutte le impostazioni richieste, inclusa l'autenticazione.</td>
      <td>Usa questa opzione se hai bisogno di controllo completo sulla configurazione da zero.</td>
    </tr>
  </tbody>
</table>

## Configurazione Generale

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithBaseURL(baseURL)</code></td>
      <td>Sovrascrive l'URL API Aruba Cloud predefinito.</td>
      <td>Il valore predefinito è <code>https://api.arubacloud.com</code>. Modificalo solo se devi puntare a un ambiente diverso.</td>
    </tr>
    <tr>
      <td><code>WithDefaultBaseURL()</code></td>
      <td>Helper che imposta l'URL API al valore predefinito di produzione.</td>
      <td>Chiamato da <code>DefaultOptions()</code>.</td>
    </tr>
  </tbody>
</table>

## Autenticazione

<p>L'autenticazione è una parte critica della configurazione. È necessario scegliere <b>uno</b> dei seguenti metodi.</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithDefaultTokenManagerSchema(clientID, clientSecret)</code></td>
      <td>Configura l'autenticazione Client Credentials standard con l'URL del token issuer predefinito e senza caching persistente. Chiamato internamente da <code>DefaultOptions()</code>.</td>
      <td>Da usare quando si vuole ripristinare il token manager alle impostazioni di fabbrica prima di applicare sovrascritture selettive.</td>
    </tr>
    <tr>
      <td><code>WithClientCredentials(clientID, clientSecret)</code></td>
      <td><b>(Consigliato)</b> Configura l'SDK per utilizzare il flusso OAuth2 Client Credentials. L'SDK gestirà automaticamente il recupero e il rinnovo del token di accesso.</td>
      <td><b>Esclusione reciproca</b>: non può essere usato con <code>WithToken()</code> o <code>WithVaultCredentialsRepository()</code>. È il metodo standard e più sicuro per l'autenticazione service-to-service.</td>
    </tr>
    <tr>
      <td><code>WithToken(token)</code></td>
      <td>Configura l'SDK per utilizzare un token di accesso OAuth2 statico preesistente.</td>
      <td><b>Esclusione reciproca</b>: non può essere usato con nessun'altra opzione di autenticazione o token issuer.<br/>
      <b>Avviso</b>: L'SDK <b>non</b> rinnoverà questo token. Una volta scaduto, tutte le chiamate API falliranno. Questo metodo è consigliato solo per client di breve durata o script semplici e atomici.</td>
    </tr>
    <tr>
      <td><code>WithVaultCredentialsRepository(...)</code></td>
      <td>Configura l'SDK per recuperare <code>clientID</code> e <code>clientSecret</code> da un'istanza HashiCorp Vault. L'SDK utilizzerà poi queste credenziali per gestire il token di accesso.</td>
      <td><b>Esclusione reciproca</b>: non può essere usato con <code>WithClientCredentials()</code> o <code>WithToken()</code>.<br/>
      <b>Parametri</b>: <code>vaultURI</code>, <code>kvMount</code>, <code>kvPath</code>, <code>namespace</code>,
      <code>rolePath</code>, <code>roleID</code>, <code>secretID</code>.</td>
    </tr>
    <tr>
      <td><code>WithTokenIssuerURL(url)</code></td>
      <td>Sovrascrive l'URL predefinito per l'endpoint del token OAuth2.</td>
      <td>Il valore predefinito è <code>https://mylogin.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token</code>. Modificalo solo se devi puntare a un endpoint di autenticazione diverso.</td>
    </tr>
    <tr>
      <td><code>WithSecurityScopes(scopes ...)</code></td>
      <td>Imposta gli scope di sicurezza da richiedere durante l'autenticazione.</td>
      <td>Sostituisce qualsiasi scope precedentemente configurato.</td>
    </tr>
    <tr>
      <td><code>WithAdditionalSecurityScopes(scopes ...)</code></td>
      <td>Aggiunge scope di sicurezza supplementari all'elenco esistente.</td>
      <td>Non sovrascrive gli scope già presenti.</td>
    </tr>
  </tbody>
</table>

## Caching del Token (Opzionale)

<p>Per migliorare le prestazioni e la resilienza, l'SDK può memorizzare nella cache il token di accesso su uno store esterno. È fortemente consigliato nelle applicazioni di produzione per evitare di richiedere un nuovo token a ogni avvio.</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithRedisTokenRepositoryFromURI(redisURI)</code></td>
      <td>Configura un'istanza Redis come cache persistente per il token di accesso, specificando l'URI completo.</td>
      <td><b>Esclusione reciproca</b>: non può essere usato con le opzioni <code>WithFileTokenRepository...</code>.<br/>
      Il formato URI è <code>redis://&lt;user&gt;:&lt;pass&gt;@&lt;host&gt;:&lt;port&gt;/&lt;db&gt;</code>.</td>
    </tr>
    <tr>
      <td><code>WithRedisTokenRepositoryFromStandardURI()</code></td>
      <td>Configura il caching Redis usando l'URI locale standard (<code>redis://admin:admin@localhost:6379/0</code>).</td>
      <td>Scorciatoia per lo sviluppo locale. Non imposta il drift di scadenza; abbinare con
      <code>WithStandardTokenExpirationDriftSeconds()</code> o usare <code>WithStandardRedisTokenRepository()</code>.</td>
    </tr>
    <tr>
      <td><code>WithFileTokenRepositoryFromBaseDir(baseDir)</code></td>
      <td>Configura una directory locale per memorizzare il token di accesso in un file JSON.</td>
      <td><b>Esclusione reciproca</b>: non può essere usato con le opzioni <code>WithRedisTokenRepository...</code>. Il processo SDK deve avere i permessi di lettura/scrittura sulla directory specificata.</td>
    </tr>
    <tr>
      <td><code>WithFileTokenRepositoryFromStandardBaseDir()</code></td>
      <td>Configura il caching su file in <code>/tmp/sdk-go</code>.</td>
      <td>Scorciatoia per lo sviluppo locale. Non imposta il drift di scadenza; abbinare con
      <code>WithStandardTokenExpirationDriftSeconds()</code> o usare <code>WithStandardFileTokenRepository()</code>.</td>
    </tr>
    <tr>
      <td><code>WithTokenExpirationDriftSeconds(seconds)</code></td>
      <td>Imposta un buffer di sicurezza (in secondi) per considerare un token scaduto prima che lo sia effettivamente.</td>
      <td>Previene race condition in cui l'SDK usa un token che scade appena prima del completamento della richiesta API. Il valore predefinito è 300 secondi (5 minuti). Questa opzione non ha effetto se non è configurato un meccanismo di caching (Redis o File).</td>
    </tr>
    <tr>
      <td><code>WithStandardTokenExpirationDriftSeconds()</code></td>
      <td>Imposta il drift di scadenza a 300 secondi (5 minuti).</td>
      <td>Equivalente a <code>WithTokenExpirationDriftSeconds(300)</code>.</td>
    </tr>
    <tr>
      <td><code>WithStandardRedisTokenRepository()</code></td>
      <td>Helper che configura il caching Redis su <code>localhost:6379</code> e imposta il drift di scadenza standard (300s).</td>
      <td>Scorciatoia comoda per lo sviluppo locale.</td>
    </tr>
    <tr>
      <td><code>WithStandardFileTokenRepository()</code></td>
      <td>Helper che configura il caching su file in <code>/tmp/sdk-go</code> e imposta il drift di scadenza standard (300s).</td>
      <td>Scorciatoia comoda per lo sviluppo locale.</td>
    </tr>
  </tbody>
</table>

## Logging

<p>Il logging è disabilitato per impostazione predefinita.</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithDefaultLogger()</code></td>
      <td>Ripristina il logger al valore predefinito dell'SDK (nessun log). Chiamato da <code>DefaultOptions()</code>.</td>
      <td>Da usare quando si vuole ripristinare esplicitamente il logging dopo aver configurato un logger personalizzato.</td>
    </tr>
    <tr>
      <td><code>WithNoLogs()</code></td>
      <td>Disabilita tutto il logging dell'SDK.</td>
      <td>Questo è il comportamento predefinito.</td>
    </tr>
    <tr>
      <td><code>WithNativeLogger()</code></td>
      <td>Abilita il logger integrato dell'SDK, che utilizza il pacchetto standard <code>log</code>.</td>
      <td>Utile per il debug di base.</td>
    </tr>
    <tr>
      <td><code>WithLoggerType(loggerType)</code></td>
      <td>Funzione più generale per impostare il tipo di logger.</td>
      <td><code>loggerType</code> può essere <code>LoggerNoLog</code> o <code>LoggerNative</code>.</td>
    </tr>
    <tr>
      <td><code>WithCustomLogger(logger)</code></td>
      <td>Inietta un logger personalizzato che rispetta l'interfaccia <code>ports/logger.Logger</code>.</td>
      <td>Consente l'integrazione con il framework di logging dell'applicazione (ad esempio, Logrus, Zap). Impostando questo valore, il tipo di logger viene automaticamente impostato su <code>loggerCustom</code>.</td>
    </tr>
  </tbody>
</table>

## Identità del Client HTTP

<p>L'SDK imposta automaticamente un header <code>User-Agent</code> su ogni richiesta in uscita, in modo che i log di accesso API possano essere attribuiti alla versione dell'SDK. Per impostazione predefinita il valore dell'header è <code>sdk-go@&lt;version&gt;</code>
(ad esempio <code>sdk-go@1.0.0</code>), derivato dalla costante <code>aruba.Version</code> definita in
<code>pkg/aruba/version.go</code>.</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithUserAgent(ua)</code></td>
      <td>Sovrascrive il valore predefinito dell'header <code>User-Agent</code> inviato con ogni richiesta.</td>
      <td>Da usare quando si costruisce uno strumento sopra l'SDK e si vuole che i log API mostrino l'identità del proprio strumento invece di (o in aggiunta a) la versione grezza dell'SDK. Esempio:
      <code>WithUserAgent("acloud-cli@1.0.0")</code>.</td>
    </tr>
  </tbody>
</table>

## Avanzato / Dipendenze Personalizzate

<p>Queste opzioni sono destinate a casi d'uso avanzati in cui è necessario iniettare componenti personalizzati nel flusso di lavoro dell'SDK.</p>

<table>
  <thead>
    <tr>
      <th>Setter Opzione</th>
      <th>Descrizione</th>
      <th>Note</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>WithCustomHTTPClient(client)</code></td>
      <td>Inietta un <code>*http.Client</code> pre-configurato.</td>
      <td>Utile per impostare timeout personalizzati, transport o altre configurazioni a livello HTTP.</td>
    </tr>
    <tr>
      <td><code>WithCustomMiddleware(middleware)</code></td>
      <td>Inietta un middleware personalizzato che rispetta l'interfaccia <code>ports/interceptor.Interceptor</code>.</td>
      <td>Consente di aggiungere logica personalizzata (come logging o tracing di richiesta/risposta) nella catena di chiamate HTTP dell'SDK. Il middleware di autenticazione dell'SDK verrà automaticamente collegato alla fine della catena di middleware personalizzata.</td>
    </tr>
  </tbody>
</table>
