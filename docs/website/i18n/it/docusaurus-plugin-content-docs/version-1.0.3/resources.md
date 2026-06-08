---
sidebar_position: 3
---

# Risorse

Questa pagina è il riferimento esaustivo per ogni wrapper di risorsa nel pacchetto `pkg/aruba`. Per ogni wrapper troverai:

1. La catena di accessor per raggiungerlo da `arubaClient`
2. Uno snippet `Create` pronto all'uso
3. I metodi accessor di risposta disponibili sul wrapper restituito

Per il walkthrough end-to-end del ciclo di vita (come `Create`, `Get`, `Update`, `List`, `Delete` e il polling si integrano) vedi la [Guida al Walkthrough API](./walkthrough).

---

## Convenzioni

Ogni risorsa segue la stessa struttura:

```go
// 1. Raggiungi il sub-client
client := arubaClient.FromX().Y()

// 2. Costruisci la richiesta inline e crea
result, err := client.Create(ctx,
    aruba.NewX().
        IntoParent(parentRef).   // scope al progetto / VPC / ecc.
        Named("my-resource").
        Tagged("env-prod").
        WithFoo(...))

// 3. Attendi che le risorse asincrone diventino pronte
if err := result.WaitUntilReady(ctx); err != nil { … }

// 4. Leggi gli accessor di risposta
fmt.Println(result.ID(), result.Name(), result.State())
```

- `aruba.NewX()` — factory constructor per ogni builder di risorsa
- `IntoFoo(ref)` — lega lo scope del genitore; accetta qualsiasi `aruba.Ref` (wrapper idratato o `aruba.URI("…")`)
- `WithFoo(...)` — setter fluenti; gli errori sono differiti fino a `Create`/`Update`
- `WaitUntilReady(ctx, opts...)` — disponibile sulle risorse marcate **async** qui sotto; vedi [Async / Await](./async) per le opzioni complete
- `aruba.URI(s)` — avvolge un percorso stringa grezzo in un `Ref` (vedi [Guida al Walkthrough API](./walkthrough#5-ottenere-una-risorsa-specifica))

:::info Formato dei tag
L'API di Aruba valida i valori dei tag contro `^[A-Za-z0-9-]{4,30}$`: **solo caratteri alfanumerici e trattini, lunghezza da 4 a 30**. Due punti, punti, underscore, spazi e altra punteggiatura vengono rifiutati con `400 — One or more validation error occurred`. L'SDK non valida i tag lato client, quindi un tag non valido fallisce solo quando la richiesta raggiunge il server.
:::

### Lettura dello stato del wrapper

Ogni wrapper promuove i campi di risposta più usati a accessor piatti. Preferiscili rispetto a `wrapper.Raw().Properties.X`:

```go
fmt.Println(result.ID())        // UUID
fmt.Println(result.Name())      // nome della risorsa
fmt.Println(result.State())     // stato del ciclo di vita
fmt.Println(result.Region())    // slug della regione
fmt.Println(result.RawJSON())   // payload JSON wire completo (per --output json)
fmt.Println(result.RawYAML())   // payload YAML wire completo (per --output yaml)
```

Gli scalari specifici della risorsa (es. `cs.Subnets()`, `vpnRoute.CloudSubnetCIDR()`, `kaas.PodCIDR()`) sono documentati nella sezione **Accessor di risposta** di ogni risorsa qui sotto. Vedi [Gestione delle risposte](./response-handling#lettura-dello-stato-del-wrapper) per la tassonomia completa degli accessor.

Ogni sezione di risorsa elenca anche i suoi **Setter** (metodi builder concatenabili raggruppati per l'ordine canonico della catena da `ai/CONVENTIONS.md`) e un collegamento all'esempio eseguibile in `examples/all-resources/`.

---

## Progetto

```go
arubaClient.FromProject()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`

> Il Progetto **non** è asincrono — è pronto in modo sincrono dopo che `Create` ritorna. Non è necessaria alcuna chiamata `WaitUntilReady`.

```go
proj, err := arubaClient.FromProject().Create(
    ctx,
    aruba.NewProject().
        Named("my-project").
        Tagged("env-prod").
        DescribedAs("Progetto di produzione").
        NotDefault())
if err != nil {
    log.Fatalf("Create project: %v", err)
}
fmt.Printf("✓ Progetto: %s (ID: %s)\n", proj.Name(), proj.ID())
```

**Accessor di risposta**:
- `ID()` — UUID della risorsa
- `URI()` — percorso completo della risorsa (es. `/projects/abc-123`)
- `Name()` — nome del progetto
- `Description()` — descrizione del progetto
- `IsDefault()` — se questo è il progetto predefinito
- `Tags()` — lista di tag `[]string`
- `CreatedAt()`, `UpdatedAt()` — timestamp
- `CreatedBy()`, `UpdatedBy()` — identificatore dell'attore che ha creato/aggiornato la risorsa (es. `aru-297647`)
- `CreatedUser()`, `UpdatedUser()` — nome visualizzato dell'utente che ha creato/aggiornato la risorsa
- `Raw()` — struct wire sottostante `*types.ProjectResponse`
- `RawJSON()` / `RawYAML()` — payload serializzato per i flag `--output json/yaml`
- `RawRequest()` — `types.ProjectRequest` per i flussi round-trip `Get → Update`

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Descriptive scalars*: `DescribedAs(string)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_project.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_project.go)
:::

---

## Audit

```go
arubaClient.FromAudit().Events()
```

**Operazioni supportate**: `List`

Gli Audit Event sono in sola lettura. Non esiste un costruttore `Create` — usa `List` con un `Ref` di progetto e opzionalmente `aruba.WithFilter(…)` per interrogare l'audit trail.

```go
list, err := arubaClient.FromAudit().Events().List(ctx, proj,
    aruba.WithLimit(50),
    aruba.WithFilter("action eq 'Create'"))
if err != nil {
    log.Fatalf("List events: %v", err)
}
for _, e := range list.Items() {
    fmt.Println(e.ID(), e.Action(), e.Timestamp())
}
```

**Accessor di risposta**:
- `ID()` — UUID dell'evento
- `URI()` — percorso della risorsa
- `ResourceURI()` — URI della risorsa a cui l'evento si riferisce
- `Action()` — stringa dell'azione (es. `"Create"`, `"Delete"`)
- `Timestamp()` — ora dell'evento
- `User()` — identificatore utente che ha scatenato l'evento
- `Raw()` — struct wire sottostante

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Compute

### Cloud Server

```go
arubaClient.FromCompute().CloudServers()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`, `PowerOn`, `PowerOff`, `SetPassword`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

Un Cloud Server dipende da risorse di rete (VPC, Subnet, Security Group), un Elastic IP, un Boot Volume (Block Storage) e un Key Pair. Crea prima queste risorse e passa i wrapper idratati come parametri `Ref`.

```go
cs, err := arubaClient.FromCompute().CloudServers().Create(
    ctx,
    aruba.NewCloudServer().
        OfFlavor(aruba.CloudServerFlavorCSO2A4).
        Named("my-server").
        Tagged("env-prod").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        InZone(aruba.ZoneITBG1).
        BootingFrom(blockStorage).
        WithVPC(vpc).
        OnSubnets(subnet).
        WithSecurityGroups(sg).
        WithElasticIP(eip).
        UsingKeyPair(keyPair))
if err != nil {
    log.Fatalf("Create Cloud Server: %v", err)
}

if err := cs.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Cloud Server did not become Ready: %v", err)
}
fmt.Printf("✓ Cloud Server: %s (zona: %s, flavor: %s)\n", cs.Name(), cs.Zone(), cs.Flavor())
```

**Azioni di alimentazione e password** (richiedono un wrapper idratato da `Create`/`Get`):

```go
if err := cs.PowerOff(ctx); err != nil { log.Fatalf("PowerOff: %v", err) }
if err := cs.PowerOn(ctx);  err != nil { log.Fatalf("PowerOn: %v", err) }
if err := cs.SetPassword(ctx, "NewStr0ngP@ss!"); err != nil { log.Fatalf("SetPassword: %v", err) }
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CloudServerID()` — ID server assegnato dal provider
- `Zone()` — zona di disponibilità
- `Flavor()` — slug del flavor di calcolo
- `FlavorRaw()` — struct flavor completa
- `VPC()` — `aruba.Ref` della VPC collegata
- `BootVolume()` — `aruba.Ref` del volume di boot
- `KeyPair()` — `aruba.Ref` della key pair
- `NetworkInterfaces()` — slice di descrittori di interfacce di rete
- `Template()` — immagine/template usata al boot
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()` — da `statusMixin`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfFlavor(CloudServerFlavor)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Descriptive scalars*: `WithUserData(string)`
- *Origin*: `BootingFrom(Ref)`
- *Attached config*: `WithVPC(Ref)`, `WithSecurityGroups(...Ref)`, `WithElasticIP(Ref)`
- *Network placement*: `OnSubnets(...Ref)`
- *Active relationship*: `UsingKeyPair(Ref)`
- *Boolean state*: `WithVPCPreset()`, `WithoutVPCPreset()`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_cloud_server.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_cloud_server.go)
:::

---

### Key Pair

```go
arubaClient.FromCompute().KeyPairs()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: no.

```go
kp, err := arubaClient.FromCompute().KeyPairs().Create(
    ctx,
    aruba.NewKeyPair().
        Named("my-keypair").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        WithPublicKey("ssh-rsa AAAAB3NzaC1yc2E..."))
if err != nil {
    log.Fatalf("Create KeyPair: %v", err)
}
fmt.Printf("✓ KeyPair: %s (ID: %s)\n", kp.Name(), kp.ID())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KeyPairID()` — ID chiave assegnato dal provider
- `PublicKey()` — stringa della chiave pubblica
- `Region()` — slug della regione
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithPublicKey(string)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_key_pair.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_key_pair.go)
:::

---

## Container

### KaaS (Kubernetes as a Service)

```go
arubaClient.FromContainer().KaaS()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`, `DownloadKubeconfig`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
k, err := arubaClient.FromContainer().KaaS().Create(
    ctx,
    aruba.NewKaaS().
        Named("my-cluster").
        Tagged("env-prod").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        WithKubernetesVersion(aruba.KubernetesVersion1323).
        WithPodCIDR("10.200.0.0/16").
        WithNodeCIDR("10.100.0.0/16", "node-cidr").
        WithVPC(vpc).
        WithSubnet(subnet).
        WithSecurityGroup(sg).
        WithNodePools(aruba.NewNodePool().
            OfInstance(aruba.NodePoolInstanceK4A8).
            Named("default-pool").
            InZone(aruba.ZoneITBG1).
            WithCount(3)).
        HighlyAvailable().
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create KaaS: %v", err)
}

if err := k.WaitUntilReady(ctx); err != nil {
    log.Fatalf("KaaS did not become Ready: %v", err)
}
fmt.Printf("✓ Cluster KaaS: %s (k8s: %s)\n", k.Name(), k.KubernetesVersion())
```

**Download del kubeconfig** (richiede un wrapper idratato):

```go
kubeconfig, err := k.DownloadKubeconfig(ctx)
if err != nil {
    log.Fatalf("DownloadKubeconfig: %v", err)
}
// kubeconfig è un []byte YAML kubeconfig
```

**Builder del node pool** — `aruba.NewNodePool()`:
- `Named(name)` — nome del pool
- `WithCount(n)` — numero di nodi
- `OfInstance(flavor)` — flavor dell'istanza nodo
- `InZone(zone)` — zona di disponibilità
- `WithAutoscaling(min, max)` — abilita l'autoscaling

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KaaSID()` — ID cluster assegnato dal provider
- `VPC()`, `Subnet()` — `aruba.Ref` alle risorse di rete collegate
- `SecurityGroupName()` — nome del security group applicato
- `KubernetesVersion()` — stringa della versione Kubernetes
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithKubernetesVersion(KubernetesVersion)`, `WithPodCIDR(string)`, `WithMaxStorageQuotaGB(int)`, `WithIdentity(string, string)`
- *Attached config*: `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithNodeCIDR(string, string)`, `WithNodePools(...*NodePool)`, `WithoutNodePools()`, `ReplaceNodePools(...*NodePool)`
- *Boolean state*: `HighlyAvailable()`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_kaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kaas.go)
:::

---

### Container Registry

```go
arubaClient.FromContainer().ContainerRegistry()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`. Questa risorsa può impiegare 20–40 minuti per convergere — utilizza un budget di attesa generoso.

```go
reg, err := arubaClient.FromContainer().ContainerRegistry().Create(
    ctx,
    aruba.NewContainerRegistry().
        OfSize(aruba.ContainerRegistrySizeFlavorSmall).
        Named("my-registry").
        Tagged("env-prod").
        InProject(proj).
        WithAdminUsername("admin").
        WithVPC(vpc).
        WithSubnet(subnet).
        WithSecurityGroup(sg).
        WithElasticIP(eip).
        WithBlockStorage(blockStorage).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create ContainerRegistry: %v", err)
}

if err := reg.WaitUntilReady(ctx); err != nil {
    log.Fatalf("ContainerRegistry did not become Ready: %v", err)
}
fmt.Printf("✓ Registry: %s (IP pubblico: %s)\n", reg.Name(), reg.PublicIP())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `ContainerRegistryID()` — ID registry assegnato dal provider
- `ElasticIP()` — URI dell'endpoint pubblico
- `VPC()`, `Subnet()`, `SecurityGroup()`, `BlockStorage()` — `aruba.Ref` alle risorse collegate
- `AdminUsername()` — utente amministratore del registry
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfSize(ContainerRegistrySizeFlavor)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithAdminUsername(string)`
- *Attached config*: `WithElasticIP(Ref)`, `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithBlockStorage(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_container_registry.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_container_registry.go)
:::

---

## Database

### DBaaS (Database as a Service)

```go
arubaClient.FromDatabase().DBaaS()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
db, err := arubaClient.FromDatabase().DBaaS().Create(
    ctx,
    aruba.NewDBaaS().
        OfEngine(aruba.DatabaseEngineMySQL80).
        OfFlavor(aruba.DBaaSFlavorDBO2A4).
        Named("my-database").
        Tagged("env-prod").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        InZone(aruba.ZoneITBG1).
        SizedGB(20).
        WithVPC(vpc).
        WithSubnet(subnet).
        WithSecurityGroup(sg).
        WithElasticIP(eip).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create DBaaS: %v", err)
}

if err := db.WaitUntilReady(ctx); err != nil {
    log.Fatalf("DBaaS did not become Ready: %v", err)
}
fmt.Printf("✓ DBaaS: %s (engine: %s)\n", db.Name(), db.Engine())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DBaaSID()` — ID istanza assegnato dal provider
- `Engine()` — identificatore dell'engine (costante `DatabaseEngine`)
- `EngineRaw()` — struct engine completa
- `Flavor()` — identificatore del flavor (costante `DBaaSFlavor`)
- `FlavorRaw()` — struct flavor completa
- `SizeGB()` — dimensione dello storage in GB
- `AutoscalingEnabled()` — bool
- `VPC()`, `Subnet()`, `SecurityGroup()`, `ElasticIP()` — `aruba.Ref` alle risorse di rete
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfEngine(DatabaseEngine)`, `OfFlavor(DBaaSFlavor)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Descriptive scalars*: `SizedGB(int)`, `WithAutoscaling(min, max int)`, `WithoutAutoscaling()`
- *Attached config*: `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithElasticIP(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_dbaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas.go)
:::

---

### Database

```go
arubaClient.FromDatabase().Databases()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
database, err := arubaClient.FromDatabase().Databases().Create(
    ctx,
    aruba.NewDatabase().
        Named("my-app-db").
        Tagged("app-backend").
        InDBaaS(db))
if err != nil {
    log.Fatalf("Create Database: %v", err)
}

if err := database.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Database did not become Ready: %v", err)
}
fmt.Printf("✓ Database: %s\n", database.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DatabaseID()` — ID database assegnato dal provider
- `DBaaSID()` — ID DBaaS padre
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Containment*: `InDBaaS(Ref)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_database.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_database.go)
:::

---

### User

```go
arubaClient.FromDatabase().Users()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
user, err := arubaClient.FromDatabase().Users().Create(
    ctx,
    aruba.NewUser().
        Tagged("app-backend").
        InDBaaS(db).
        WithUsername("app_user").
        WithPassword("Str0ngP@ssword!"))
if err != nil {
    log.Fatalf("Create User: %v", err)
}

if err := user.WaitUntilReady(ctx); err != nil {
    log.Fatalf("User did not become Ready: %v", err)
}
fmt.Printf("✓ Utente: %s\n", user.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `UserID()` — ID utente assegnato dal provider
- `Username()` — nome utente del database
- `DBaaSID()` — ID DBaaS padre
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `WithUsername(string)`
- *Containment*: `InDBaaS(Ref)`
- *Descriptive scalars*: `WithPassword(string)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_dbaas_user.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas_user.go)
:::

---

### Grant

```go
arubaClient.FromDatabase().Grants()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
grant, err := arubaClient.FromDatabase().Grants().Create(
    ctx,
    aruba.NewGrant().
        OfRole("liteadmin").
        InDatabase(database).
        ForUser("app_user"))
if err != nil {
    log.Fatalf("Create Grant: %v", err)
}

if err := grant.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Grant did not become Ready: %v", err)
}
fmt.Printf("✓ Grant: %s (privilegi: %s)\n", grant.Name(), grant.Privileges())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `GrantID()` — ID grant assegnato dal provider
- `DatabaseID()` — ID Database padre
- `Privileges()` — stringa dei privilegi
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Containment*: `InDatabase(Ref)`
- *Active relationship*: `ForUser(string)`, `OfRole(string)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_grant.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_grant.go)
:::

---

### DBaaS Backup

```go
arubaClient.FromDatabase().DBaaSBackups()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
backup, err := arubaClient.FromDatabase().DBaaSBackups().Create(
    ctx,
    aruba.NewDBaaSBackup().
        Named("my-db-backup").
        Tagged("backup").
        InProject(proj).
        FromDBaaS(db).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create DBaaSBackup: %v", err)
}

if err := backup.WaitUntilReady(ctx); err != nil {
    log.Fatalf("DBaaS Backup did not become Ready: %v", err)
}
fmt.Printf("✓ DBaaS Backup: %s\n", backup.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DBaaSBackupID()` — ID backup assegnato dal provider
- `DBaaSURI()` — URI del DBaaS sorgente
- `DatabaseURI()` — URI del Database sorgente (se applicabile)
- `SizeGB()` — dimensione del backup in GB
- `Zone()` — zona di disponibilità
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`, `FromDBaaS(Ref)`, `FromDatabase(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Le operazioni di DBaaS Backup sono trattate nell'esempio DBaaS: [`examples/all-resources/resource_dbaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas.go)
:::

---

## Metric

### Alert

```go
arubaClient.FromMetric().Alerts()
```

**Operazioni supportate**: `List`

Gli Alert sono in sola lettura. Usa `List` con un `Ref` di progetto per interrogare gli alert attivi.

```go
list, err := arubaClient.FromMetric().Alerts().List(ctx, proj)
if err != nil {
    log.Fatalf("List Alerts: %v", err)
}
for _, a := range list.Items() {
    fmt.Println(a.ID(), a.Name(), a.IsActive())
}
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`
- `Threshold()` — valore soglia dell'alert
- `Action()` — azione scatenata all'alert
- `IsActive()` — bool
- `Raw()` — struct wire sottostante

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### Metric

```go
arubaClient.FromMetric().Metrics()
```

**Operazioni supportate**: `List`

Le Metric sono risultati di query di serie temporali in sola lettura.

```go
list, err := arubaClient.FromMetric().Metrics().List(ctx, proj,
    aruba.WithFilter("resource eq '"+cs.URI()+"'"))
if err != nil {
    log.Fatalf("List Metrics: %v", err)
}
for _, m := range list.Items() {
    fmt.Println(m.ID(), m.Name())
}
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`
- `Raw()` — struct wire sottostante

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Network

### VPC

```go
arubaClient.FromNetwork().VPCs()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(
    ctx,
    aruba.NewVPC().
        Named("my-vpc").
        Tagged("network").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        NotDefault().
        WithoutPreset())
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}

if err := vpc.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPC did not become Ready: %v", err)
}
fmt.Printf("✓ VPC: %s\n", vpc.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPCID()` — ID VPC assegnato dal provider
- `Region()` — slug della regione
- `IsDefault()`, `IsPreset()` — flag
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Boolean state*: `AsDefault()`, `NotDefault()`, `WithPreset()`, `WithoutPreset()`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_vpc.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_vpc.go)
:::

---

### Subnet

```go
arubaClient.FromNetwork().Subnets()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

`OfType` accetta `aruba.SubnetTypeBasic` o `aruba.SubnetTypeAdvanced` (costanti tipizzate — nessun cast a stringa necessario).

`aruba.NewSubnetDHCP()` è un sub-builder per la configurazione DHCP. Collegalo con `WithDHCP(...)`.

```go
subnet, err := arubaClient.FromNetwork().Subnets().Create(
    ctx,
    aruba.NewSubnet().
        OfType(aruba.SubnetTypeAdvanced).
        Named("my-subnet").
        Tagged("network").
        InVPC(vpc).
        InRegion(aruba.RegionITBGBergamo).
        WithCIDR("192.168.1.0/25").
        WithDHCP(aruba.NewSubnetDHCP().
            Enabled().
            WithRange("192.168.1.10", 50).
            WithRoutes(aruba.SubnetDHCPRouteCommon{Address: "0.0.0.0/0", Gateway: "192.168.1.1"}).
            WithDNSServers("8.8.8.8", "8.8.4.4")).
        NotDefault())
if err != nil {
    log.Fatalf("Create Subnet: %v", err)
}

if err := subnet.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Subnet did not become Ready: %v", err)
}
fmt.Printf("✓ Subnet: %s (CIDR: %s)\n", subnet.Name(), subnet.CIDR())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SubnetID()` — ID subnet assegnato dal provider
- `Type()` — tipo di subnet (costante `SubnetType`)
- `CIDR()` — blocco CIDR
- `DHCP()` — configurazione DHCP
- `IsDefault()` — bool
- `Region()` — slug della regione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfType(SubnetType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithCIDR(string)`
- *Attached config*: `WithDHCP(*SubnetDHCPCommon)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_subnet.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_subnet.go)
:::

---

### Elastic IP

```go
arubaClient.FromNetwork().ElasticIPs()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
eip, err := arubaClient.FromNetwork().ElasticIPs().Create(
    ctx,
    aruba.NewElasticIP().
        Named("my-eip").
        Tagged("network").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create ElasticIP: %v", err)
}

if err := eip.WaitUntilReady(ctx); err != nil {
    log.Fatalf("ElasticIP did not become Ready: %v", err)
}
fmt.Printf("✓ Elastic IP: %s (%s)\n", eip.Name(), eip.Address())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `ElasticIPID()` — ID IP assegnato dal provider
- `Address()` — indirizzo IP pubblico allocato
- `BillingPeriod()` — cadenza di fatturazione
- `Region()` — slug della regione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilNotUsed(ctx, opts...)`, `WaitUntilUsed(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_elastic_ip.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_elastic_ip.go)
:::

---

### Security Group

```go
arubaClient.FromNetwork().SecurityGroups()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
sg, err := arubaClient.FromNetwork().SecurityGroups().Create(
    ctx,
    aruba.NewSecurityGroup().
        Named("my-security-group").
        Tagged("security").
        InVPC(vpc).
        NotDefault())
if err != nil {
    log.Fatalf("Create SecurityGroup: %v", err)
}

if err := sg.WaitUntilReady(ctx); err != nil {
    log.Fatalf("SecurityGroup did not become Active: %v", err)
}
fmt.Printf("✓ Security Group: %s (ID: %s)\n", sg.Name(), sg.ID())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SecurityGroupID()` — ID gruppo assegnato dal provider
- `IsDefault()` — bool
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_security_group.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_security_group.go)
:::

---

### Security Rule

```go
arubaClient.FromNetwork().SecurityGroupRules()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — `State()` e `FailureReason()` sono disponibili.

`WithDirection` accetta `aruba.RuleDirectionIngress` o `aruba.RuleDirectionEgress`. `WithProtocol` accetta `aruba.RuleProtocolTCP`, `aruba.RuleProtocolUDP`, `aruba.RuleProtocolICMP` o `aruba.RuleProtocolANY`.

> **Avvertenza**: `TargetingCIDR` e `TargetingSecurityGroup` sono mutuamente esclusivi. Impostare entrambi registra un errore al momento del setter, che emerge al momento di `Create`.

```go
rule, err := arubaClient.FromNetwork().SecurityGroupRules().Create(
    ctx,
    aruba.NewSecurityRule().
        Named("allow-ssh").
        Tagged("ssh-key").
        InSecurityGroup(sg).
        InRegion(aruba.RegionITBGBergamo).
        WithDirection(aruba.RuleDirectionIngress).
        WithProtocol(aruba.RuleProtocolTCP).
        WithPort("22").
        TargetingCIDR("0.0.0.0/0"))
if err != nil {
    log.Fatalf("Create SecurityRule: %v", err)
}
fmt.Printf("✓ Security Rule: %s\n", rule.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SecurityRuleID()` — ID regola assegnato dal provider
- `Direction()` — `"Ingress"` o `"Egress"`
- `Protocol()` — es. `"TCP"`, `"UDP"`, `"ICMP"`
- `Port()` — numero di porta o intervallo
- `TargetKind()` — `"Ip"` o `"SecurityGroup"`
- `TargetValue()` — stringa CIDR o URI del Security Group
- `Region()` — slug della regione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InSecurityGroup(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithDirection(RuleDirection)`, `WithProtocol(RuleProtocol)`, `WithPort(string)`
- *Active relationship*: `TargetingCIDR(string)`, `TargetingSecurityGroup(Ref)`

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### Load Balancer

```go
arubaClient.FromNetwork().LoadBalancers()
```

**Operazioni supportate**: `List`, `Get`

I Load Balancer sono in sola lettura tramite questo SDK — vengono creati e gestiti automaticamente dalla piattaforma Aruba Cloud.

```go
list, err := arubaClient.FromNetwork().LoadBalancers().List(ctx, proj)
if err != nil {
    log.Fatalf("List LoadBalancers: %v", err)
}
for _, lb := range list.Items() {
    fmt.Println(lb.ID(), lb.Name(), lb.Address())
}
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `LoadBalancerID()` — ID LB assegnato dal provider
- `Address()` — indirizzo pubblico
- `VPC()` — `aruba.Ref` alla VPC collegata
- `Region()` — slug della regione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `Raw()` — struct wire sottostante

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPC Peering

```go
arubaClient.FromNetwork().VPCPeerings()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
peering, err := arubaClient.FromNetwork().VPCPeerings().Create(
    ctx,
    aruba.NewVPCPeering().
        Named("my-peering").
        Tagged("network").
        InVPC(vpc).
        InRegion(aruba.RegionITBGBergamo).
        WithPeerVPC(aruba.URI("/projects/"+peerProjectID+"/vpcs/"+peerVPCID)))
if err != nil {
    log.Fatalf("Create VPCPeering: %v", err)
}

if err := peering.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPCPeering did not become Active: %v", err)
}
fmt.Printf("✓ VPC Peering: %s\n", peering.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPCPeeringID()` — ID peering assegnato dal provider
- `VPCID()` — ID VPC sorgente
- `PeerVPC()` — `aruba.Ref` alla VPC peer
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Geography*: `InRegion(Region)`
- *Active relationship*: `PeeredWith(Ref)`

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPC Peering Route

```go
arubaClient.FromNetwork().VPCPeeringRoutes()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
route, err := arubaClient.FromNetwork().VPCPeeringRoutes().Create(
    ctx,
    aruba.NewVPCPeeringRoute().
        Named("my-peering-route").
        Tagged("network").
        InVPCPeering(peering).
        InRegion(aruba.RegionITBGBergamo).
        WithCIDR("10.0.0.0/8").
        WithTarget(aruba.URI("/projects/"+projectID+"/vpcs/"+vpcID)))
if err != nil {
    log.Fatalf("Create VPCPeeringRoute: %v", err)
}

if err := route.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPCPeeringRoute did not become Active: %v", err)
}
fmt.Printf("✓ Peering Route: %s (CIDR: %s)\n", route.Name(), route.CIDR())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CIDR()` — blocco CIDR della route
- `Target()` — `aruba.Ref` al target della route
- `VPCPeeringID()` — ID peering padre
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPCPeering(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithLocalCIDR(string)`, `WithRemoteCIDR(string)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPN Tunnel

```go
arubaClient.FromNetwork().VPNTunnels()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

Sub-builder del VPN Tunnel:
- `aruba.NewVPNIKE()` — parametri fase 1 IKE (`WithEncryption(IKEEncryption)`, `WithHash(IKEHash)`, `WithDHGroup(IKEDHGroup)`, `WithDPDAction(IKEDPDAction)`)
- `aruba.NewVPNESP()` — parametri fase 2 ESP (`WithEncryption(ESPEncryption)`, `WithHash(ESPHash)`, `WithPFS(ESPPFSGroup)`)
- `aruba.NewVPNPSK()` — configurazione pre-shared key (`WithKey(string)`, `WithCloudSite(string)`, `WithOnPremSite(string)`)

```go
tunnel, err := arubaClient.FromNetwork().VPNTunnels().Create(
    ctx,
    aruba.NewVPNTunnel().
        Named("my-vpn-tunnel").
        Tagged("vpn-net").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        WithPeerClientPublicIP("203.0.113.1").
        WithIKESettings(aruba.NewVPNIKE().
            WithEncryption(aruba.IKEEncryptionAES256).
            WithHash(aruba.IKEHashSHA256).
            WithDHGroup(aruba.IKEDHGroup14)).
        WithESPSettings(aruba.NewVPNESP().
            WithEncryption(aruba.ESPEncryptionAES256).
            WithHash(aruba.ESPHashSHA256)).
        WithPSKSettings(aruba.NewVPNPSK().
            WithKey("my-pre-shared-key")))
if err != nil {
    log.Fatalf("Create VPNTunnel: %v", err)
}

if err := tunnel.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPNTunnel did not become Active: %v", err)
}
fmt.Printf("✓ VPN Tunnel: %s (gateway: %s)\n", tunnel.Name(), tunnel.PeerClientPublicIP())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPNTunnelID()` — ID tunnel assegnato dal provider
- `PeerClientPublicIP()` — IP gateway del peer remoto
- `IKE()` — impostazioni IKE `*aruba.VPNIKE`
- `ESP()` — impostazioni ESP `*aruba.VPNESP`
- `PSK()` — impostazioni PSK `*aruba.VPNPSK`
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfType(VPNType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithVPNClientProtocol(VPNClientProtocol)`, `WithPeerClientPublicIP(string)`
- *Attached config*: `WithIPConfig(*VPNIPConfig)`, `WithIKESettings(*VPNIKE)`, `WithESPSettings(*VPNESP)`, `WithPSKSettings(*VPNPSK)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPN Route

```go
arubaClient.FromNetwork().VPNRoutes()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
vpnRoute, err := arubaClient.FromNetwork().VPNRoutes().Create(
    ctx,
    aruba.NewVPNRoute().
        Named("my-vpn-route").
        Tagged("vpn-net").
        InVPNTunnel(tunnel).
        InRegion(aruba.RegionITBGBergamo).
        WithCIDR("10.0.0.0/8").
        WithTarget(aruba.URI("/projects/"+projectID+"/vpcs/"+vpcID)))
if err != nil {
    log.Fatalf("Create VPNRoute: %v", err)
}

if err := vpnRoute.WaitUntilReady(ctx); err != nil {
    log.Fatalf("VPNRoute did not become Active: %v", err)
}
fmt.Printf("✓ VPN Route: %s (CIDR: %s)\n", vpnRoute.Name(), vpnRoute.CIDR())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CIDR()` — blocco CIDR della route
- `Target()` — `aruba.Ref` al target della route
- `VPNTunnelID()` — ID VPN Tunnel padre
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPNTunnel(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithCloudSubnet(string)`, `WithOnPremSubnet(string)`

:::tip Esempio eseguibile
Eseguito come parte dell'orchestratore: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Schedule

### Job

```go
arubaClient.FromSchedule().Jobs()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — `State()` e `FailureReason()` sono disponibili.

Usa `OneShotAt(t time.Time)` per pianificare un job one-shot, oppure `WithCron(expr string)` per un job ricorrente con una pianificazione cron. Usa `RecurringUntil(t time.Time)` per impostare una data di fine per un job ricorrente.

```go
// Job one-shot — eseguito una sola volta a un orario specifico
job, err := arubaClient.FromSchedule().Jobs().Create(
    ctx,
    aruba.NewJob().
        Named("my-one-shot-job").
        Tagged("automation").
        InProject(proj).
        OneShotAt(time.Now().Add(10*time.Minute)))
if err != nil {
    log.Fatalf("Create Job: %v", err)
}
fmt.Printf("✓ Job: %s (tipo: %s)\n", job.Name(), job.JobType())

// Job ricorrente — eseguito secondo una pianificazione cron
cronJob, err := arubaClient.FromSchedule().Jobs().Create(
    ctx,
    aruba.NewJob().
        Named("my-recurring-job").
        Tagged("automation").
        InProject(proj).
        WithCron("0 * * * *"))
if err != nil {
    log.Fatalf("Create recurring Job: %v", err)
}
fmt.Printf("✓ Job ricorrente: %s (cron: %s)\n", cronJob.Name(), cronJob.Cron())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `JobID()` — ID job assegnato dal provider
- `JobType()` — tipo di job (`types.JobTypeOneShot` o `types.JobTypeRecurring`)
- `Cron()` — espressione cron (job ricorrenti)
- `IsEnabled()` — bool
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfType(JobType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `OneShotAt(time.Time)`, `StartingAt(time.Time)`, `WithCron(string)`, `RecurringUntil(time.Time)`, `WithSteps(...*JobStep)`
- *Boolean state*: `Enabled()`, `Disabled()`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_job.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_job.go)
:::

---

## Security

### KMS (Key Management Service)

```go
arubaClient.FromSecurity().KMS()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
kms, err := arubaClient.FromSecurity().KMS().Create(
    ctx,
    aruba.NewKMS().
        Named("my-kms").
        Tagged("security").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create KMS: %v", err)
}

if err := kms.WaitUntilReady(ctx); err != nil {
    log.Fatalf("KMS did not become Active: %v", err)
}
fmt.Printf("✓ KMS: %s (ID: %s)\n", kms.Name(), kms.ID())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KMSID()` — ID istanza KMS assegnato dal provider
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

### Key

```go
arubaClient.FromSecurity().Keys()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — `State()` e `FailureReason()` sono disponibili.

`OfAlgorithm` accetta `aruba.KeyAlgorithmAes` o `aruba.KeyAlgorithmRsa` (costanti tipizzate — nessun cast a stringa necessario).

```go
key, err := arubaClient.FromSecurity().Keys().Create(
    ctx,
    aruba.NewKey().
        OfAlgorithm(aruba.KeyAlgorithmAes).
        Named("my-encryption-key").
        Tagged("security").
        InKMS(kms))
if err != nil {
    log.Fatalf("Create Key: %v", err)
}
fmt.Printf("✓ Chiave: %s (algoritmo: %s)\n", key.Name(), key.Algorithm())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KeyID()` — ID chiave assegnato dal provider
- `Algorithm()` — stringa dell'algoritmo
- `Type()` — `"Symmetric"` o `"Asymmetric"`
- `Status()` — stato del ciclo di vita della chiave
- `CreationSource()` — come è stata creata la chiave
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfAlgorithm(KeyAlgorithm)`
- *Name*: `Named(string)`
- *Containment*: `InKMS(Ref)`

:::tip Esempio eseguibile
Le operazioni sulle chiavi sono trattate nell'esempio KMS: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

### Kmip

```go
arubaClient.FromSecurity().Kmips()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`. `WaitUntilReady` del KMIP ha successo su `"CertificateAvailable"` o `"Active"`. `WaitUntilCertificateAvailable` è un alias per `WaitUntilReady`.

```go
km, err := arubaClient.FromSecurity().Kmips().Create(
    ctx,
    aruba.NewKmip().
        Named("my-kmip").
        Tagged("security").
        InKMS(kms))
if err != nil {
    log.Fatalf("Create Kmip: %v", err)
}

if err := km.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Kmip did not become ready: %v", err)
}
fmt.Printf("✓ Kmip: %s\n", km.Name())
```

**Download del certificato KMIP** (richiede un wrapper idratato):

```go
cert, err := km.Download(ctx)
if err != nil {
    log.Fatalf("Download Kmip certificate: %v", err)
}
fmt.Println("Cert:", cert.Cert())
fmt.Println("Key:",  cert.Key())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KmipID()` — ID KMIP assegnato dal provider
- `KmipStatus()` — stato specifico del KMIP
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilCertificateAvailable(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Containment*: `InKMS(Ref)`

:::tip Esempio eseguibile
Le operazioni KMIP sono trattate nell'esempio KMS: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

## Storage

### Block Storage (Volume)

```go
arubaClient.FromStorage().Volumes()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

`OfType` accetta `aruba.BlockStorageTypeStandard` o `aruba.BlockStorageTypePerformance`. Usa `AsBootable()` per contrassegnare un volume come avviabile; `NotBootable()` per annullare l'impostazione. Usa `FromImage(imageID)` per specificare un'immagine di base.

```go
bs, err := arubaClient.FromStorage().Volumes().Create(
    ctx,
    aruba.NewBlockStorage().
        OfType(aruba.BlockStorageTypeStandard).
        Named("my-volume").
        Tagged("storage").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        InZone(aruba.ZoneITBG1).
        SizedGB(20).
        FromImage("LU22-001").
        AsBootable().
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create BlockStorage: %v", err)
}

if err := bs.WaitUntilReady(ctx); err != nil {
    log.Fatalf("BlockStorage did not become Active: %v", err)
}
fmt.Printf("✓ Volume: %s (%d GB)\n", bs.Name(), bs.SizeGB())
```

Per creare un volume **da uno snapshot**, usa `FromSnapshot(snapshot)` invece di `FromImage`:

```go
bs, err := arubaClient.FromStorage().Volumes().Create(
    ctx,
    aruba.NewBlockStorage().
        OfType(aruba.BlockStorageTypeStandard).
        Named("restored-volume").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        InZone(aruba.ZoneITBG1).
        SizedGB(20).
        FromSnapshot(snapshot).
        AsBootable().
        BilledBy(aruba.BillingPeriodHour))
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BlockStorageID()` — ID volume assegnato dal provider
- `SizeGB()` — dimensione in GB
- `Type()` — tipo di storage
- `Zone()` — zona di disponibilità
- `BillingPeriod()` — cadenza di fatturazione
- `IsBootable()` — bool
- `Image()` — riferimento all'immagine
- `SnapshotURI()` — URI dello snapshot sorgente (se creato da snapshot)
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilNotUsed(ctx, opts...)`, `WaitUntilUsed(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfType(BlockStorageType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Descriptive scalars*: `SizedGB(int)`
- *Origin*: `FromImage(string)`, `FromSnapshot(Ref)`
- *Boolean state*: `AsBootable()`, `NotBootable()`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_block_storage.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_block_storage.go)
:::

---

### Snapshot

```go
arubaClient.FromStorage().Snapshots()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Update`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
snap, err := arubaClient.FromStorage().Snapshots().Create(
    ctx,
    aruba.NewSnapshot().
        Named("my-snapshot").
        Tagged("backup").
        InProject(proj).
        InRegion(aruba.RegionITBGBergamo).
        FromVolume(bs).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create Snapshot: %v", err)
}

if err := snap.WaitUntilReady(ctx); err != nil {
    log.Fatalf("Snapshot did not become Active: %v", err)
}
fmt.Printf("✓ Snapshot: %s\n", snap.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SnapshotID()` — ID snapshot assegnato dal provider
- `SizeGB()` — dimensione dello snapshot in GB
- `Type()` — tipo di storage
- `Zone()` — zona di disponibilità
- `BillingPeriod()` — cadenza di fatturazione
- `Bootable()` — bool
- `VolumeURI()` — URI del volume sorgente
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Origin*: `FromVolume(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_snapshot.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_snapshot.go)
:::

---

### Storage Backup

```go
arubaClient.FromStorage().Backups()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

`OfType` accetta `aruba.StorageBackupTypeFull` o `aruba.StorageBackupTypeIncremental`. Usa `FromVolume(vol)` per specificare il volume sorgente.

```go
backup, err := arubaClient.FromStorage().Backups().Create(
    ctx,
    aruba.NewStorageBackup().
        OfType(aruba.StorageBackupTypeFull).
        Named("my-backup").
        Tagged("backup").
        InProject(proj).
        RetainedForDays(30).
        FromVolume(bs).
        BilledBy(aruba.BillingPeriodHour))
if err != nil {
    log.Fatalf("Create StorageBackup: %v", err)
}

if err := backup.WaitUntilReady(ctx); err != nil {
    log.Fatalf("StorageBackup did not become Active: %v", err)
}
fmt.Printf("✓ Storage Backup: %s\n", backup.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BackupID()` — ID backup assegnato dal provider
- `Type()` — tipo di backup
- `RetentionDays()` — periodo di conservazione in giorni
- `OriginURI()` — URI del volume sorgente
- `BillingPeriod()` — cadenza di fatturazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Classifier*: `OfType(StorageBackupType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `RetainedForDays(int)`
- *Origin*: `FromVolume(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_storage_backup.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_storage_backup.go)
:::

---

### Storage Restore

```go
arubaClient.FromStorage().Restores()
```

**Operazioni supportate**: `Create`, `List`, `Get`, `Delete`
**Asincrono**: sì — chiama `WaitUntilReady(ctx)` dopo `Create`.

```go
restore, err := arubaClient.FromStorage().Restores().Create(
    ctx,
    aruba.NewStorageRestore().
        Named("my-restore").
        Tagged("restore").
        FromBackup(backup).
        WithTarget(aruba.URI(backup.OriginURI())))
if err != nil {
    log.Fatalf("Create StorageRestore: %v", err)
}

if err := restore.WaitUntilReady(ctx); err != nil {
    log.Fatalf("StorageRestore did not become Active: %v", err)
}
fmt.Printf("✓ Storage Restore: %s\n", restore.Name())
```

**Accessor di risposta**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `RestoreID()` — ID restore assegnato dal provider
- `TargetURI()` — URI del volume di destinazione
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — struct wire sottostante

**Setter**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `FromBackup(Ref)`
- *Geography*: `InRegion(Region)`
- *Origin*: `ToVolume(Ref)`

:::tip Esempio eseguibile
Esempio end-to-end completo: [`examples/all-resources/resource_storage_restore.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_storage_restore.go)
:::

---

## Opzioni di chiamata

Passa le opzioni di chiamata come argomenti variadic a qualsiasi chiamata `List`, `Get`, `Create`, `Update` o `Delete`:

| Opzione | Scopo |
|--------|---------|
| `aruba.WithFilter(expr)` | Espressione filtro lato server |
| `aruba.WithSort(expr)` | Espressione di ordinamento |
| `aruba.WithLimit(n)` | Dimensione della pagina |
| `aruba.WithOffset(n)` | Offset di paginazione |
| `aruba.WithProjection(expr)` | Proiezione dei campi |
| `aruba.WithAPIVersion(v)` | Sovrascrive la versione API per questa chiamata |

Vedi [Filtri](./filters) per la sintassi dei filtri e degli ordinamenti.

---

## Costanti enum

Tutti i tipi enum sono ri-esportati da `pkg/aruba` — nessun import aggiuntivo necessario. L'elenco canonico si trova in `pkg/aruba/aliases.go`.

### Regione e Zona

| Costante | Descrizione |
|----------|-------------|
| `aruba.RegionITBGBergamo` | Datacenter di Bergamo (Italia) |
| `aruba.ZoneITBG1` | Zona di disponibilità Bergamo 1 |
| `aruba.ZoneITBG2` | Zona di disponibilità Bergamo 2 |
| `aruba.ZoneITBG3` | Zona di disponibilità Bergamo 3 |

### Fatturazione

| Costante | Descrizione |
|----------|-------------|
| `aruba.BillingPeriodHour` | Fatturazione oraria |
| `aruba.BillingPeriodMonth` | Fatturazione mensile |

### Network

| Costante | Valore |
|----------|-------|
| `aruba.RuleDirectionIngress` | `"Ingress"` |
| `aruba.RuleDirectionEgress` | `"Egress"` |
| `aruba.RuleProtocolTCP` | `"TCP"` |
| `aruba.RuleProtocolUDP` | `"UDP"` |
| `aruba.RuleProtocolICMP` | `"ICMP"` |
| `aruba.RuleProtocolANY` | (wildcard — qualsiasi protocollo) |
| `aruba.SubnetTypeBasic` | `"Basic"` |
| `aruba.SubnetTypeAdvanced` | `"Advanced"` |

### Compute

| Costante | Descrizione |
|----------|-------------|
| `aruba.CloudServerFlavorCSO1A2` | 1 vCPU, 2 GB RAM |
| `aruba.CloudServerFlavorCSO2A4` | 2 vCPU, 4 GB RAM |
| `aruba.CloudServerFlavorCSO4A8` | 4 vCPU, 8 GB RAM |
| `aruba.CloudServerFlavorCSO8A16` | 8 vCPU, 16 GB RAM |
| … (vedi `aliases.go` per l'elenco completo) | |

### Container

| Costante | Descrizione |
|----------|-------------|
| `aruba.KubernetesVersion1323` | Kubernetes 1.32.3 |
| `aruba.KubernetesVersion1332` | Kubernetes 1.33.2 |
| `aruba.KubernetesVersion1341` | Kubernetes 1.34.1 |
| `aruba.NodePoolInstanceK2A4` | 2 vCPU, 4 GB RAM |
| `aruba.NodePoolInstanceK4A8` | 4 vCPU, 8 GB RAM |
| `aruba.NodePoolInstanceK8A16` | 8 vCPU, 16 GB RAM |
| … (vedi `aliases.go` per l'elenco completo) | |
| `aruba.ContainerRegistrySizeFlavorSmall` | Livello small di utenti concorrenti |
| `aruba.ContainerRegistrySizeFlavorMedium` | Livello medium di utenti concorrenti |
| `aruba.ContainerRegistrySizeFlavorHighPerf` | Livello ad alte prestazioni |

### Database

| Costante | Descrizione |
|----------|-------------|
| `aruba.DatabaseEngineMySQL80` | MySQL 8.0 |
| `aruba.DatabaseEngineMSSQL2022Web` | SQL Server 2022 Web |
| `aruba.DatabaseEngineMSSQL2022Standard` | SQL Server 2022 Standard |
| `aruba.DatabaseEngineMSSQL2022Enterprise` | SQL Server 2022 Enterprise |
| `aruba.DBaaSFlavorDBO1A2` | 1 vCPU, 2 GB RAM |
| `aruba.DBaaSFlavorDBO2A4` | 2 vCPU, 4 GB RAM |
| `aruba.DBaaSFlavorDBO4A8` | 4 vCPU, 8 GB RAM |
| … (vedi `aliases.go` per l'elenco completo) | |

### Storage

| Costante | Valore |
|----------|-------|
| `aruba.BlockStorageTypeStandard` | `"Standard"` |
| `aruba.BlockStorageTypePerformance` | `"Performance"` |
| `aruba.StorageBackupTypeFull` | `"Full"` |
| `aruba.StorageBackupTypeIncremental` | `"Incremental"` |

### Security

| Costante | Valore |
|----------|-------|
| `aruba.KeyAlgorithmAes` | `"Aes"` |
| `aruba.KeyAlgorithmRsa` | `"Rsa"` |
| `aruba.KeyTypeSymmetric` | `"Symmetric"` |
| `aruba.KeyTypeAsymmetric` | `"Asymmetric"` |
| `aruba.ServiceStatusCertificateAvailable` | `"CertificateAvailable"` |

### Crittografia VPN

| Costante | Descrizione |
|----------|-------------|
| `aruba.IKEEncryptionAES256` | AES-256 CBC (IKE fase 1) |
| `aruba.IKEHashSHA256` | HMAC-SHA-256 (IKE fase 1) |
| `aruba.IKEDHGroup14` | Gruppo Diffie-Hellman MODP-2048 |
| `aruba.ESPEncryptionAES256` | AES-256 CBC (ESP fase 2) |
| `aruba.ESPHashSHA256` | HMAC-SHA-256 (ESP fase 2) |
| `aruba.ESPPFSGroupEnable` | PFS abilitato (gruppo DH negoziato) |
| `aruba.ESPPFSGroupDisable` | PFS disabilitato |
| … (vedi `aliases.go` per gli elenchi completi) | |

### Schedule

| Costante | Valore |
|----------|-------|
| `aruba.JobTypeOneShot` | `"OneShot"` |
| `aruba.JobTypeRecurring` | `"Recurring"` |

---

## Appendice: Tipi wire grezzi (`pkg/types`)

I seguenti tipi sono le struct wire di basso livello sottostanti. Normalmente vi si accede solo tramite `.Raw()` o `.RawRequest()` su un wrapper, oppure quando si costruiscono integrazioni avanzate con `pkg/async`. Sono anche ri-esportati come alias di tipo `aruba.XxxRequest` / `aruba.XxxResponse` in modo da poterli referenziare senza un import aggiuntivo.

| Tipo | File | Note |
|------|------|-------|
| `Response[T]` | `resource.go` | Envelope HTTP generico restituito dagli adapter di basso livello |
| `ErrorResponse` | `error.go` | Errore strutturato RFC 7807 |
| `ListResponse` | `resource.go` | Link di paginazione e conteggio totale |
| `ResourceMetadataRequest` | `resource.go` | Nome + tag per Create |
| `RegionalResourceMetadataRequest` | `resource.go` | Estende i metadati con Location |
| `ResourceMetadataResponse` | `resource.go` | ID, URI, Name, timestamp |
| `ResourceStatusResponse` | `resource.go` | Campo State |
| `ReferenceResourceCommon` | `resource.go` | `{uri: "…"}` link a un'altra risorsa |
| `RequestParameters` | `parameters.go` | Struct filtro/ordinamento/limit/offset di basso livello (preferire gli helper `CallOption`) |
| `ProjectRequest` / `ProjectResponse` / `ProjectListResponse` | `project.project.go` | |
| `VPCRequest` / `VPCResponse` / `VPCListResponse` | `network.vpc.go` | |
| `SubnetRequest` / `SubnetResponse` / `SubnetListResponse` | `network.subnet.go` | |
| `SecurityGroupRequest` / `SecurityGroupResponse` | `network.security-group.go` | |
| `SecurityRuleRequest` / `SecurityRuleResponse` | `network.security-rule.go` | |
| `ElasticIPRequest` / `ElasticIPResponse` | `network.elastic-ip.go` | |
| `CloudServerRequest` / `CloudServerResponse` | `compute.cloudserver.go` | |
| `KeyPairRequest` / `KeyPairResponse` | `compute.keypair.go` | |
| `KaaSRequest` / `KaaSResponse` / `KaaSUpdateRequest` | `container.kaas.go` | |
| `ContainerRegistryRequest` / `ContainerRegistryResponse` | `container.containerregistry.go` | |
| `DBaaSRequest` / `DBaaSResponse` | `database.dbaas.go` | |
| `KmsRequest` / `KmsResponse` | `security.kms.go` | |
| `KeyRequest` / `KeyResponse` | `security.kms.go` | |
| `KmipRequest` / `KmipResponse` / `KmipCertificateResponse` | `security.kms.go` | |
| `BlockStorageRequest` / `BlockStorageResponse` | `storage.block-storage.go` | |
| `SnapshotRequest` / `SnapshotResponse` | `storage.snapshot.go` | |
| `StorageBackupRequest` / `StorageBackupResponse` | `storage.backup.go` | |
| `JobRequest` / `JobResponse` / `JobListResponse` | `schedule.job.go` | |
| `AlertResponse` / `AlertsListResponse` | `metrics.alert.go` | |
| `MetricResponse` / `MetricListResponse` | `metrics.metric.go` | |
| `AuditEvent` / `AuditEventListResponse` | `audit.event.go` | |
| `VPCPeeringRequest` / `VPCPeeringResponse` | `network.vpc-peering.go` | |
| `VPNTunnelRequest` / `VPNTunnelResponse` | `network.vpn-tunnel.go` | |
| `VPNRouteRequest` / `VPNRouteResponse` | `network.vpn-route.go` | |
| `LoadBalancerResponse` | `network.load-balancer.go` | |
