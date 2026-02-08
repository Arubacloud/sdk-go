# Opzioni di Configurazione SDK

L'SDK Aruba Cloud per Go è configurato utilizzando un'API fluente che ti consente di concatenare i setter di opzioni per costruire una configurazione client. Questa guida dettaglia tutte le opzioni disponibili, raggruppate per argomento.

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
      <td>Crea una configurazione standard per l'ambiente di produzione utilizzando il tuo Client ID e Secret API.</td>
      <td>Questo è il punto di partenza consigliato per la maggior parte delle applicazioni. Imposta l'API di produzione e gli URL dei token e configura l'autenticazione.</td>
    </tr>
    <tr>
      <td><code>NewOptions()</code></td>
      <td>Crea un nuovo builder di configurazione vuoto. Devi configurare manualmente tutte le impostazioni richieste, inclusa l'autenticazione.</td>
      <td>Usa questo se hai bisogno di controllo completo sulla configurazione da zero.</td>
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
      <td>Il predefinito è <code>https://api.arubacloud.com</code>. Dovresti cambiarlo solo se devi puntare a un ambiente diverso.</td>
    </tr>
    <tr>
      <td><code>WithDefaultBaseURL()</code></td>
      <td>Un helper che imposta l'URL API al predefinito di produzione.</td>
      <td>Chiamato da <code>DefaultOptions()</code>.</td>
    </tr>
  </tbody>
</table>

## Autenticazione

<p>L'autenticazione è una parte critica della configurazione. Devi scegliere <b>uno</b> dei seguenti metodi.</p>

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
      <td><code>WithClientCredentials(clientID, clientSecret)</code></td>
      <td><b>(Consigliato)</b> Configura l'SDK per utilizzare il flusso OAuth2 Client Credentials. L'SDK gestirà automaticamente il recupero e il rinnovo del token di accesso.</td>
      <td><b>Esclusione Reciproca</b>: Non può essere utilizzato con <code>WithToken()</code> o <code>WithVaultCredentialsRepository()</code>. Questo è il metodo standard e più sicuro per l'autenticazione service-to-service.</td>
    </tr>
    <tr>
      <td><code>WithToken(token)</code></td>
      <td>Configura l'SDK per utilizzare un token di accesso OAuth2 statico preesistente.</td>
      <td><b>Esclusione Reciproca</b>: Non può essere utilizzato con qualsiasi altra opzione di autenticazione o token issuer.<br/> <b>Avviso</b>: L'SDK <b>non</b> rinnoverà questo token. Una volta scaduto, tutte le chiamate API falliranno. Questo metodo è consigliato solo per client di breve durata o script semplici e atomici.</td>
    </tr>
    <tr>
      <td><code>WithVaultCredentialsRepository(...)</code></td>
      <td>Configura l'SDK per recuperare <code>clientID</code> e <code>clientSecret</code> da un'istanza HashiCorp Vault. L'SDK utilizzerà quindi queste credenziali per gestire il token di accesso.</td>
      <td><b>Esclusione Reciproca</b>: Non può essere utilizzato con <code>WithClientCredentials()</code> o <code>WithToken()</code>.<br/> <b>Parametri</b>: <code>vaultURI</code>, <code>kvMount</code>, <code>kvPath</code>, <code>namespace</code>, <code>rolePath</code>, <code>roleID</code>, <code>secretID</code>.</td>
    </tr>
    <tr>
      <td><code>WithTokenIssuerURL(url)</code></td>
      <td>Sovrascrive l'URL predefinito per l'endpoint del token OAuth2.</td>
      <td>Il predefinito è <code>https://login.aruba.it/...</code>. Dovresti cambiarlo solo se devi puntare a un endpoint di autenticazione diverso.</td>
    </tr>
    <tr>
      <td><code>WithSecurityScopes(scopes ...)</code></td>
      <td>Imposta gli scope di sicurezza da richiedere durante l'autenticazione.</td>
      <td>Sostituisce qualsiasi scope precedentemente impostato.</td>
    </tr>
    <tr>
      <td><code>WithAdditionalSecurityScopes(scopes ...)</code></td>
      <td>Aggiunge scope di sicurezza aggiuntivi all'elenco esistente.</td>
      <td>Non sovrascrive gli scope esistenti.</td>
    </tr>
  </tbody>
</table>

## Cache Token (Opzionale)

<p>Per migliorare le prestazioni e la resilienza, l'SDK può memorizzare nella cache il token di accesso in un archivio esterno. Questo è altamente consigliato per le applicazioni di produzione per evitare di recuperare un nuovo token ad ogni avvio.</p>

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
      <td>Configura un'istanza Redis come cache persistente per il token di accesso.</td>
      <td><b>Esclusione Reciproca</b>: Non può essere utilizzato con le opzioni <code>WithFileTokenRepository...</code>.<br/> Il formato URI è <code>redis://&lt;user&gt;:&lt;pass&gt;@&lt;host&gt;:&lt;port&gt;/&lt;db&gt;</code>.</td>
    </tr>
    <tr>
      <td><code>WithFileTokenRepositoryFromBaseDir(baseDir)</code></td>
      <td>Configura una directory locale per memorizzare il token di accesso in un file JSON.</td>
      <td><b>Esclusione Reciproca</b>: Non può essere utilizzato con le opzioni <code>WithRedisTokenRepository...</code>. Il processo SDK deve avere i permessi di lettura/scrittura sulla directory specificata.</td>
    </tr>
    <tr>
      <td><code>WithTokenExpirationDriftSeconds(seconds)</code></td>
      <td>Imposta un buffer di sicurezza (in secondi) per trattare un token come scaduto prima che lo faccia effettivamente.</td>
      <td>Questo previene condizioni di race dove l'SDK utilizza un token che scade appena prima che la richiesta API venga completata. Il predefinito è 300 secondi (5 minuti). Questa opzione non ha effetto se un meccanismo di cache (Redis o File) non è configurato.</td>
    </tr>
    <tr>
      <td><code>WithStandardRedisTokenRepository()</code></td>
      <td>Helper per configurare la cache Redis su <code>localhost:6379</code> e imposta lo drift di scadenza standard (300s).</td>
      <td>Una scorciatoia conveniente per lo sviluppo locale.</td>
    </tr>
    <tr>
      <td><code>WithStandardFileTokenRepository()</code></td>
      <td>Helper per configurare la cache file in <code>/tmp/sdk-go</code> e imposta lo drift di scadenza standard (300s).</td>
      <td>Una scorciatoia conveniente per lo sviluppo locale.</td>
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
      <td>Una funzione più generale per impostare il tipo di logger.</td>
      <td><code>loggerType</code> può essere <code>LoggerNoLog</code> o <code>LoggerNative</code>.</td>
    </tr>
    <tr>
      <td><code>WithCustomLogger(logger)</code></td>
      <td>Inietta un logger personalizzato che rispetta l'interfaccia <code>ports/logger.Logger</code>.</td>
      <td>Consente l'integrazione con il framework di logging della tua applicazione (ad esempio, Logrus, Zap). Impostare questo imposta automaticamente il tipo di logger su <code>loggerCustom</code>.</td>
    </tr>
  </tbody>
</table>

## Avanzato / Dipendenze Personalizzate

<p>Queste opzioni sono per casi d'uso avanzati dove devi iniettare i tuoi componenti personalizzati nel flusso di lavoro dell'SDK.</p>

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
      <td>Ti consente di aggiungere logica personalizzata (come logging o tracing di richiesta/risposta) nella catena di chiamate HTTP dell'SDK. Il middleware di autenticazione dell'SDK sarà automaticamente collegato alla fine della tua catena di middleware personalizzata.</td>
    </tr>
  </tbody>
</table>
