# Tipi di Dati SDK

L'SDK Aruba Cloud per Go utilizza un insieme coerente di tipi di dati per richieste e risposte in tutti i gruppi API. Questo documento descrive i tipi generici e condivisi che formano la base dei modelli dell'SDK, seguiti da una suddivisione dei tipi specifici per ogni gruppo API.

## Tipi Generici e Condivisi

La maggior parte dei tipi specifici delle risorse sono composti da questi elementi fondamentali. Comprenderli è fondamentale per utilizzare l'SDK in modo efficace.

<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>File</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>Response[T]</code></td>
      <td><code>resource.go</code></td>
      <td>Un wrapper generico per tutte le risposte delle chiamate API. Contiene i dati analizzati con successo (<code>Data *T</code>) per risposte 2xx o un errore strutturato (<code>Error *ErrorResponse</code>) per risposte 4xx/5xx.</td>
    </tr>
    <tr>
      <td><code>ErrorResponse</code></td>
      <td><code>error.go</code></td>
      <td>Una struct che segue RFC 7807 Problem Details for HTTP APIs. Fornisce informazioni di errore strutturate, inclusi titolo, dettaglio e codice di stato.</td>
    </tr>
    <tr>
      <td><code>ListResponse</code></td>
      <td><code>resource.go</code></td>
      <td>Una struct incorporata utilizzata in tutte le risposte di lista (ad esempio, <code>VPCList</code>, <code>CloudServerList</code>). Fornisce dettagli di paginazione, incluso il numero totale di elementi e i link alle pagine successive/precedenti.</td>
    </tr>
    <tr>
      <td><code>ResourceMetadataRequest</code></td>
      <td><code>resource.go</code></td>
      <td>Una struct standard per denominare una risorsa e assegnare tag durante la creazione. È tipicamente incorporata nel tipo <code>...Request</code> specifico di una risorsa.</td>
    </tr>
    <tr>
      <td><code>RegionalResourceMetadataRequest</code></td>
      <td><code>resource.go</code></td>
      <td>Estende <code>ResourceMetadataRequest</code> aggiungendo un campo obbligatorio <code>Location</code>. Utilizzato per creare risorse che devono essere collegate a una regione specifica del data center.</td>
    </tr>
    <tr>
      <td><code>ResourceMetadataResponse</code></td>
      <td><code>resource.go</code></td>
      <td>La struttura di metadati standard restituita per ogni risorsa. Include l'<code>ID</code> univoco della risorsa, <code>URI</code>, <code>Name</code>, <code>Location</code>, timestamp di creazione/aggiornamento e altro ancora.</td>
    </tr>
    <tr>
      <td><code>ResourceStatus</code></td>
      <td><code>resource.go</code></td>
      <td>Una struct comune restituita con le risposte delle risorse, che indica lo stato corrente della risorsa (ad esempio, "Active", "Creating", "Error").</td>
    </tr>
    <tr>
      <td><code>ReferenceResource</code></td>
      <td><code>resource.go</code></td>
      <td>Una struct semplice utilizzata per collegarsi a un'altra risorsa tramite il suo <code>URI</code> univoco. Questo è comunemente utilizzato nei corpi delle richieste per specificare dipendenze.</td>
    </tr>
    <tr>
        <td><code>RequestParameters</code></td>
        <td><code>parameters.go</code></td>
        <td>Una struct utilizzata per fornire parametri di query opzionali per le chiamate API, come filtri, ordinamento, paginazione (limit/offset) e versioning API.</td>
    </tr>
  </tbody>
</table>

## Tipi per Gruppo API

Le seguenti sezioni dettagliano i tipi principali di richiesta e risposta per ogni gruppo API.

### Tipi Progetto
*File: `project.project.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>ProjectRequest</code></td>
      <td>Create/Update</td>
      <td>Il payload per creare o aggiornare un progetto.</td>
    </tr>
    <tr>
      <td><code>ProjectResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa progetto.</td>
    </tr>
    <tr>
      <td><code>ProjectList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di progetti.</td>
    </tr>
  </tbody>
</table>

### Tipi Audit
*File: `audit.event.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>AuditEvent</code></td>
      <td>List</td>
      <td>Rappresenta un singolo record di evento di audit.</td>
    </tr>
    <tr>
      <td><code>AuditEventListResponse</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di eventi di audit.</td>
    </tr>
  </tbody>
</table>

### Tipi Compute
*File: `compute.cloudserver.go`, `compute.keypair.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>CloudServerRequest</code></td>
      <td>Create/Update</td>
      <td>Il payload per creare o aggiornare un Cloud Server. Il campo <code>Properties</code> contiene <code>CloudServerPropertiesRequest</code>, che include un campo opzionale <code>UserData</code> (stringa nullable) per contenuto cloud-init codificato in base64 per inizializzare il server.</td>
    </tr>
    <tr>
      <td><code>CloudServerResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa Cloud Server.</td>
    </tr>
    <tr>
      <td><code>CloudServerList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di Cloud Server.</td>
    </tr>
    <tr>
      <td><code>CloudServerPasswordRequest</code></td>
      <td>SetPassword</td>
      <td>Il payload per impostare o modificare una password del Cloud Server.</td>
    </tr>
    <tr>
      <td><code>KeyPairRequest</code></td>
      <td>Create</td>
      <td>Il payload per creare una coppia di chiavi SSH.</td>
    </tr>
    <tr>
      <td><code>KeyPairResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa coppia di chiavi SSH.</td>
    </tr>
    <tr>
      <td><code>KeyPairListResponse</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di coppie di chiavi SSH.</td>
    </tr>
  </tbody>
</table>

### Tipi Container
*File: `container.kaas.go`, `container.containerregistry.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>KaaSRequest</code></td>
      <td>Create</td>
      <td>Il payload per creare un cluster KaaS.</td>
    </tr>
    <tr>
      <td><code>KaaSUpdateRequest</code></td>
      <td>Update</td>
      <td>Il payload per aggiornare un cluster KaaS (campi limitati).</td>
    </tr>
    <tr>
      <td><code>KaaSResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa cluster KaaS.</td>
    </tr>
    <tr>
      <td><code>KaaSList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di cluster KaaS.</td>
    </tr>
    <tr>
      <td><code>KaaSKubeconfigResponse</code></td>
      <td>DownloadKubeconfig</td>
      <td>Rappresenta la risposta di download del file kubeconfig con nome file e contenuto base64.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryRequest</code></td>
      <td>Create/Update</td>
      <td>Il payload per creare o aggiornare un Container Registry.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa Container Registry.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di Container Registry.</td>
    </tr>
  </tbody>
</table>

### Tipi Database
*File: `database.dbaas.go`, `database.database.go`, ecc.*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>DBaaSRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare/aggiornare un'istanza DBaaS.</td>
    </tr>
    <tr>
      <td><code>DBaaSResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola istanza DBaaS.</td>
    </tr>
    <tr>
      <td><code>DatabaseRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare/aggiornare un database all'interno di un'istanza DBaaS.</td>
    </tr>
    <tr>
      <td><code>UserRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare/aggiornare un utente per un'istanza DBaaS.</td>
    </tr>
    <tr>
      <td><code>GrantRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per gestire i permessi utente su un database.</td>
    </tr>
    <tr>
      <td><code>BackupRequest</code></td>
      <td>Create</td>
      <td>Payload per creare un backup di un'istanza DBaaS.</td>
    </tr>
  </tbody>
</table>

### Tipi Metric
*File: `metrics.alert.go`, `metrics.metric.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>AlertResponse</code></td>
      <td>List</td>
      <td>Rappresenta un singolo alert di monitoraggio.</td>
    </tr>
    <tr>
      <td><code>AlertsListResponse</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di alert.</td>
    </tr>
    <tr>
      <td><code>MetricResponse</code></td>
      <td>List</td>
      <td>Rappresenta un insieme di punti dati time-series per una metrica specifica.</td>
    </tr>
    <tr>
      <td><code>MetricListResponse</code></td>
      <td>List</td>
      <td>Rappresenta un elenco di metriche.</td>
    </tr>
  </tbody>
</table>

### Tipi Network
*File: `network.vpc.go`, `network.subnet.go`, ecc.*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>VPCRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un Virtual Private Cloud.</td>
    </tr>
    <tr>
      <td><code>SubnetRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare una Subnet all'interno di una VPC.</td>
    </tr>
    <tr>
      <td><code>ElasticIPRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un Elastic IP.</td>
    </tr>
    <tr>
      <td><code>SecurityGroupRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un Security Group.</td>
    </tr>
    <tr>
      <td><code>SecurityRuleRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare una regola all'interno di un Security Group.</td>
    </tr>
    <tr>
      <td><code>LoadBalancerResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola risorsa Load Balancer.</td>
    </tr>
     <tr>
      <td><code>VPCPeeringRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare una connessione VPC Peering.</td>
    </tr>
     <tr>
      <td><code>VPNTunnelRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un Tunnel VPN.</td>
    </tr>
  </tbody>
</table>

### Tipi Schedule
*File: `schedule.job.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>JobRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un job programmato.</td>
    </tr>
    <tr>
      <td><code>JobResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta un singolo job programmato.</td>
    </tr>
    <tr>
      <td><code>JobList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di job programmati.</td>
    </tr>
  </tbody>
</table>

### Tipi Security
*File: `security.kms.go`*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>KmsRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare una chiave KMS.</td>
    </tr>
    <tr>
      <td><code>KmsResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta una singola chiave KMS.</td>
    </tr>
    <tr>
      <td><code>KmsList</code></td>
      <td>List</td>
      <td>Rappresenta un elenco paginato di chiavi KMS.</td>
    </tr>
  </tbody>
</table>

### Tipi Storage
*File: `storage.block-storage.go`, `storage.snapshot.go`, ecc.*
<table>
  <thead>
    <tr>
      <th>Nome Tipo</th>
      <th>Utilizzo</th>
      <th>Descrizione</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>BlockStorageRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un volume di storage a blocchi.</td>
    </tr>
    <tr>
      <td><code>BlockStorageResponse</code></td>
      <td>Get/List</td>
      <td>Rappresenta un singolo volume di storage a blocchi.</td>
    </tr>
    <tr>
      <td><code>SnapshotRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare uno snapshot di volume.</td>
    </tr>
    <tr>
      <td><code>StorageBackupRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per creare o aggiornare un backup di volume.</td>
    </tr>
    <tr>
      <td><code>RestoreRequest</code></td>
      <td>Create/Update</td>
      <td>Payload per ripristinare un backup su un volume.</td>
    </tr>
  </tbody>
</table>
