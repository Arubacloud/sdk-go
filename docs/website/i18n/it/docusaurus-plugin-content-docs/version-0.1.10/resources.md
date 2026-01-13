# Gruppi API e Risorse

L'SDK Aruba Cloud è organizzato in diversi gruppi API, ciascuno corrispondente a un'area di servizio specifica (ad esempio, Compute, Network, Storage). Puoi accedere a questi gruppi dall'oggetto principale `arubaClient`. Questo documento fornisce un elenco completo di tutti i gruppi disponibili e delle risorse che gestiscono.

## Progetto

Il Progetto è la risorsa di livello superiore sotto la quale sono organizzate tutte le altre risorse. Il client per la gestione dei progetti è accessibile direttamente.

<table>
  <thead>
    <tr>
      <th>Accessore Client</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>arubaClient.FromProject()</code></td>
      <td>Gestisce i progetti, che sono contenitori per tutte le altre risorse cloud.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Audit

Gestisce l'accesso alle informazioni di audit trail.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Events()</code></td>
      <td>Fornisce l'accesso agli eventi di audit per un progetto.</td>
      <td><code>List</code></td>
    </tr>
  </tbody>
</table>

## Compute

Gestisce macchine virtuali e risorse correlate.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.CloudServers()</code></td>
      <td>Gestisce istanze di macchine virtuali (Cloud Server).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code>, <code>PowerOn</code>, <code>PowerOff</code>, <code>SetPassword</code></td>
    </tr>
    <tr>
      <td><code>.KeyPairs()</code></td>
      <td>Gestisce coppie di chiavi SSH per l'accesso al server.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Container

Gestisce servizi basati su container.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.KaaS()</code></td>
      <td>Gestisce cluster Kubernetes as a Service (KaaS).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code>, <code>DownloadKubeconfig</code></td>
    </tr>
    <tr>
      <td><code>.ContainerRegistry()</code></td>
      <td>Gestisce registri container privati.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Database

Gestisce Database as a Service (DBaaS) e le sue sotto-risorse.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.DBaaS()</code></td>
      <td>Gestisce istanze DBaaS (ad esempio, MySQL, PostgreSQL).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Databases()</code></td>
      <td>Gestisce database individuali all'interno di un'istanza DBaaS.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Backups()</code></td>
      <td>Gestisce i backup delle istanze DBaaS.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Users()</code></td>
      <td>Gestisce gli utenti del database per un'istanza DBaaS.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Grants()</code></td>
      <td>Gestisce i permessi utente (grant) sui database.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Metric

Fornisce l'accesso ai dati di monitoraggio e agli alert.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Alerts()</code></td>
      <td>Fornisce l'accesso agli alert di monitoraggio.</td>
      <td><code>List</code></td>
    </tr>
    <tr>
      <td><code>.Metrics()</code></td>
      <td>Fornisce l'accesso ai dati di monitoraggio time-series per le risorse.</td>
      <td><code>List</code></td>
    </tr>
  </tbody>
</table>

## Network

Gestisce tutte le risorse di rete.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.VPCs()</code></td>
      <td>Gestisce Virtual Private Cloud (VPC).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Subnets()</code></td>
      <td>Gestisce le subnet all'interno di una VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.ElasticIPs()</code></td>
      <td>Gestisce indirizzi IP pubblici statici (Elastic IP).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.SecurityGroups()</code></td>
      <td>Gestisce i gruppi di sicurezza (firewall) all'interno di una VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.SecurityGroupRules()</code></td>
      <td>Gestisce le regole individuali all'interno di un gruppo di sicurezza.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.LoadBalancers()</code></td>
      <td>Gestisce i bilanciatori di carico.</td>
      <td><code>List</code>, <code>Get</code></td>
    </tr>
    <tr>
      <td><code>.VPCPeerings()</code></td>
      <td>Gestisce le connessioni di peering tra due VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPCPeeringRoutes()</code></td>
      <td>Gestisce le rotte per una connessione di peering VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPNTunnels()</code></td>
      <td>Gestisce tunnel VPN Site-to-Site.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPNRoutes()</code></td>
      <td>Gestisce le rotte per un tunnel VPN.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Schedule

Gestisce job automatizzati programmati.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Jobs()</code></td>
      <td>Gestisce job programmati (one-shot o ricorrenti) che possono eseguire azioni sulle risorse.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Security

Gestisce servizi relativi alla sicurezza.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.KMSKeys()</code></td>
      <td>Gestisce le chiavi del Key Management Service (KMS).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Storage

Gestisce lo storage a blocchi e le risorse di protezione dati correlate.

<table>
  <thead>
    <tr>
      <th>Client Risorsa</th>
      <th>Descrizione</th>
      <th>Operazioni Disponibili</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Volumes()</code></td>
      <td>Gestisce i volumi di storage a blocchi.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Snapshots()</code></td>
      <td>Gestisce snapshot point-in-time dei volumi di storage a blocchi.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Backups()</code></td>
      <td>Gestisce i backup dei volumi di storage a blocchi.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Restores()</code></td>
      <td>Gestisce il ripristino di un backup su un volume.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>
