---
sidebar_position: 3
---

# Resources

This page is the exhaustive reference for every resource wrapper in the `pkg/aruba` package. For each wrapper you will find:

1. The accessor chain to reach it from `arubaClient`
2. A copy-paste-ready `Create` snippet
3. The response accessor methods available on the returned wrapper

For the end-to-end lifecycle walkthrough (how `Create`, `Get`, `Update`, `List`, `Delete`, and polling fit together) see the [API Walkthrough](./walkthrough).

---

## Conventions

Every resource follows the same shape:

```go
// 1. Reach the sub-client
client := arubaClient.FromX().Y()

// 2. Build the request inline and create
result, err := client.Create(ctx,
    aruba.NewX().
        IntoParent(parentRef).   // scope to project / VPC / etc.
        Named("my-resource").
        Tagged("env-prod").
        WithFoo(...))

// 3. Wait for async resources to become ready
if err := result.WaitUntilReady(ctx); err != nil { … }

// 4. Read response accessors
fmt.Println(result.ID(), result.Name(), result.State())
```

- `aruba.NewX()` — factory constructor for every resource builder
- `IntoFoo(ref)` — binds the parent scope; accepts any `aruba.Ref` (hydrated wrapper or `aruba.URI("…")`)
- `WithFoo(...)` — fluent setters; errors are deferred until `Create`/`Update`
- `WaitUntilReady(ctx, opts...)` — available on resources marked **async** below; see [Async / Await](./async) for full options
- `aruba.URI(s)` — wraps a raw string path into a `Ref` (see [API Walkthrough](./walkthrough#5-get-a-specific-resource))

:::info Tag format
The Aruba API validates tag values against `^[A-Za-z0-9-]{4,30}$`: **alphanumerics and hyphens only, length 4 to 30**. Colons, dots, underscores, spaces, and other punctuation are rejected with `400 — One or more validation error occurred`. The SDK does not validate tag values client-side, so an invalid tag only fails when the request reaches the server.
:::

### Reading wrapper state

Every wrapper promotes the most-used response fields to flat accessors. Prefer these over `wrapper.Raw().Properties.X`:

```go
fmt.Println(result.ID())        // UUID
fmt.Println(result.Name())      // resource name
fmt.Println(result.State())     // lifecycle state
fmt.Println(result.Region())    // region slug
fmt.Println(result.RawJSON())   // full JSON wire payload (for --output json)
fmt.Println(result.RawYAML())   // full YAML wire payload (for --output yaml)
```

Resource-specific scalars (e.g. `cs.Subnets()`, `vpnRoute.CloudSubnetCIDR()`, `kaas.PodCIDR()`) are documented in each resource's **Response accessors** section below. See [Response Handling](./response-handling#reading-wrapper-state) for the complete accessor taxonomy.

Each resource section also lists its **Setters** (chainable builder methods grouped by the canonical chain order from `ai/CONVENTIONS.md`) and a link to the runnable example in `examples/all-resources/`.

---

## Project

```go
arubaClient.FromProject()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`

> Project is **not** async — it is synchronously ready after `Create` returns. No `WaitUntilReady` call is needed.

```go
proj, err := arubaClient.FromProject().Create(
    ctx,
    aruba.NewProject().
        Named("my-project").
        Tagged("env-prod").
        DescribedAs("Production project").
        NotDefault())
if err != nil {
    log.Fatalf("Create project: %v", err)
}
fmt.Printf("✓ Project: %s (ID: %s)\n", proj.Name(), proj.ID())
```

**Response accessors**:
- `ID()` — resource UUID
- `URI()` — full resource path (e.g. `/projects/abc-123`)
- `Name()` — project name
- `Description()` — project description
- `IsDefault()` — whether this is the default project
- `Tags()` — `[]string` tag list
- `CreatedAt()`, `UpdatedAt()` — timestamps
- `CreatedBy()`, `UpdatedBy()` — actor identifier (e.g. `aru-297647`) of the creator/last updater
- `CreatedUser()`, `UpdatedUser()` — display name of the creator/last updater
- `Raw()` — underlying `*types.ProjectResponse` wire struct
- `RawJSON()` / `RawYAML()` — serialized payload for `--output json/yaml` flags
- `RawRequest()` — `types.ProjectRequest` for round-trip `Get → Update` flows

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Descriptive scalars*: `DescribedAs(string)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_project.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_project.go)
:::

---

## Audit

```go
arubaClient.FromAudit().Events()
```

**Supported operations**: `List`

Audit Events are read-only. There is no `Create` constructor — use `List` with a project `Ref` and optional `aruba.WithFilter(…)` to query the audit trail.

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

**Response accessors**:
- `ID()` — event UUID
- `URI()` — resource path
- `ResourceURI()` — URI of the resource the event relates to
- `Action()` — action string (e.g. `"Create"`, `"Delete"`)
- `Timestamp()` — event time
- `User()` — user identifier who triggered the event
- `Raw()` — underlying wire struct

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Compute

### Cloud Server

```go
arubaClient.FromCompute().CloudServers()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`, `PowerOn`, `PowerOff`, `SetPassword`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

A Cloud Server depends on network resources (VPC, Subnet, Security Group), an Elastic IP, a Boot Volume (Block Storage), and a Key Pair. Create those first and pass the hydrated wrappers as `Ref` parameters.

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
fmt.Printf("✓ Cloud Server: %s (zone: %s, flavor: %s)\n", cs.Name(), cs.Zone(), cs.Flavor())
```

**Power and password actions** (require a hydrated wrapper from `Create`/`Get`):

```go
if err := cs.PowerOff(ctx); err != nil { log.Fatalf("PowerOff: %v", err) }
if err := cs.PowerOn(ctx);  err != nil { log.Fatalf("PowerOn: %v", err) }
if err := cs.SetPassword(ctx, "NewStr0ngP@ss!"); err != nil { log.Fatalf("SetPassword: %v", err) }
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CloudServerID()` — provider-assigned server ID
- `Zone()` — availability zone
- `Flavor()` — compute flavor slug
- `FlavorRaw()` — full flavor struct
- `VPC()` — `aruba.Ref` of the attached VPC
- `BootVolume()` — `aruba.Ref` of the boot volume
- `KeyPair()` — `aruba.Ref` of the key pair
- `NetworkInterfaces()` — slice of network interface descriptors
- `Template()` — image/template used at boot
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()` — from `statusMixin`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
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

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_cloud_server.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_cloud_server.go)
:::

---

### Key Pair

```go
arubaClient.FromCompute().KeyPairs()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: no.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KeyPairID()` — provider-assigned key ID
- `PublicKey()` — public key string
- `Region()` — region slug
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithPublicKey(string)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_key_pair.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_key_pair.go)
:::

---

## Container

### KaaS (Kubernetes as a Service)

```go
arubaClient.FromContainer().KaaS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`, `DownloadKubeconfig`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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
fmt.Printf("✓ KaaS cluster: %s (k8s: %s)\n", k.Name(), k.KubernetesVersion())
```

**Download kubeconfig** (requires a hydrated wrapper):

```go
kubeconfig, err := k.DownloadKubeconfig(ctx)
if err != nil {
    log.Fatalf("DownloadKubeconfig: %v", err)
}
// kubeconfig is a []byte YAML kubeconfig
```

**Node pool builder** — `aruba.NewNodePool()`:
- `Named(name)` — pool name
- `WithCount(n)` — number of nodes
- `OfInstance(flavor)` — node instance flavor
- `InZone(zone)` — availability zone
- `WithAutoscaling(min, max)` — enable autoscaling

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KaaSID()` — provider-assigned cluster ID
- `VPC()`, `Subnet()` — `aruba.Ref` to attached network resources
- `SecurityGroupName()` — name of the applied security group
- `KubernetesVersion()` — Kubernetes version string
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithKubernetesVersion(KubernetesVersion)`, `WithPodCIDR(string)`, `WithMaxStorageQuotaGB(int)`, `WithIdentity(string, string)`
- *Attached config*: `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithNodeCIDR(string, string)`, `WithNodePools(...*NodePool)`, `WithoutNodePools()`, `ReplaceNodePools(...*NodePool)`
- *Boolean state*: `HighlyAvailable()`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_kaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kaas.go)
:::

---

### Container Registry

```go
arubaClient.FromContainer().ContainerRegistry()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`. This resource can take 20–40 minutes to converge — use a generous wait budget.

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
fmt.Printf("✓ Registry: %s (public IP: %s)\n", reg.Name(), reg.PublicIP())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `ContainerRegistryID()` — provider-assigned registry ID
- `ElasticIP()` — public endpoint URI
- `VPC()`, `Subnet()`, `SecurityGroup()`, `BlockStorage()` — `aruba.Ref` to attached resources
- `AdminUsername()` — registry admin user
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfSize(ContainerRegistrySizeFlavor)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithAdminUsername(string)`
- *Attached config*: `WithElasticIP(Ref)`, `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithBlockStorage(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_container_registry.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_container_registry.go)
:::

---

## Database

### DBaaS (Database as a Service)

```go
arubaClient.FromDatabase().DBaaS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DBaaSID()` — provider-assigned instance ID
- `Engine()` — engine identifier (`DatabaseEngine` constant)
- `EngineRaw()` — full engine struct
- `Flavor()` — flavor identifier (`DBaaSFlavor` constant)
- `FlavorRaw()` — full flavor struct
- `SizeGB()` — storage size in GB
- `AutoscalingEnabled()` — bool
- `VPC()`, `Subnet()`, `SecurityGroup()`, `ElasticIP()` — `aruba.Ref` to networking resources
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfEngine(DatabaseEngine)`, `OfFlavor(DBaaSFlavor)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Descriptive scalars*: `SizedGB(int)`, `WithAutoscaling(min, max int)`, `WithoutAutoscaling()`
- *Attached config*: `WithVPC(Ref)`, `WithSubnet(Ref)`, `WithSecurityGroup(Ref)`, `WithElasticIP(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_dbaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas.go)
:::

---

### Database

```go
arubaClient.FromDatabase().Databases()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DatabaseID()` — provider-assigned database ID
- `DBaaSID()` — parent DBaaS ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Containment*: `InDBaaS(Ref)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_database.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_database.go)
:::

---

### User

```go
arubaClient.FromDatabase().Users()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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
fmt.Printf("✓ User: %s\n", user.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `UserID()` — provider-assigned user ID
- `Username()` — database username
- `DBaaSID()` — parent DBaaS ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `WithUsername(string)`
- *Containment*: `InDBaaS(Ref)`
- *Descriptive scalars*: `WithPassword(string)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_dbaas_user.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas_user.go)
:::

---

### Grant

```go
arubaClient.FromDatabase().Grants()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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
fmt.Printf("✓ Grant: %s (privileges: %s)\n", grant.Name(), grant.Privileges())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `GrantID()` — provider-assigned grant ID
- `DatabaseID()` — parent Database ID
- `Privileges()` — privilege string
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Containment*: `InDatabase(Ref)`
- *Active relationship*: `ForUser(string)`, `OfRole(string)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_grant.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_grant.go)
:::

---

### DBaaS Backup

```go
arubaClient.FromDatabase().DBaaSBackups()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DBaaSBackupID()` — provider-assigned backup ID
- `DBaaSURI()` — source DBaaS URI
- `DatabaseURI()` — source Database URI (if applicable)
- `SizeGB()` — backup size in GB
- `Zone()` — availability zone
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`, `FromDBaaS(Ref)`, `FromDatabase(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
DBaaS Backup operations are covered in the DBaaS example: [`examples/all-resources/resource_dbaas.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_dbaas.go)
:::

---

## Metric

### Alert

```go
arubaClient.FromMetric().Alerts()
```

**Supported operations**: `List`

Alerts are read-only. Use `List` with a project `Ref` to query active alerts.

```go
list, err := arubaClient.FromMetric().Alerts().List(ctx, proj)
if err != nil {
    log.Fatalf("List Alerts: %v", err)
}
for _, a := range list.Items() {
    fmt.Println(a.ID(), a.Name(), a.IsActive())
}
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`
- `Threshold()` — alert threshold value
- `Action()` — action triggered on alert
- `IsActive()` — bool
- `Raw()` — underlying wire struct

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### Metric

```go
arubaClient.FromMetric().Metrics()
```

**Supported operations**: `List`

Metrics are read-only time-series query results.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`
- `Raw()` — underlying wire struct

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Network

### VPC

```go
arubaClient.FromNetwork().VPCs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPCID()` — provider-assigned VPC ID
- `Region()` — region slug
- `IsDefault()`, `IsPreset()` — flags
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Boolean state*: `AsDefault()`, `NotDefault()`, `WithPreset()`, `WithoutPreset()`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_vpc.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_vpc.go)
:::

---

### Subnet

```go
arubaClient.FromNetwork().Subnets()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

`OfType` accepts `aruba.SubnetTypeBasic` or `aruba.SubnetTypeAdvanced` (typed constants — no string cast needed).

`aruba.NewSubnetDHCP()` is a sub-builder for DHCP configuration. Attach it with `WithDHCP(...)`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SubnetID()` — provider-assigned subnet ID
- `Type()` — subnet type (`SubnetType` constant)
- `CIDR()` — CIDR block
- `DHCP()` — DHCP configuration
- `IsDefault()` — bool
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfType(SubnetType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithCIDR(string)`
- *Attached config*: `WithDHCP(*SubnetDHCPCommon)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_subnet.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_subnet.go)
:::

---

### Elastic IP

```go
arubaClient.FromNetwork().ElasticIPs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `ElasticIPID()` — provider-assigned IP ID
- `Address()` — the allocated public IP address
- `BillingPeriod()` — billing cadence
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilNotUsed(ctx, opts...)`, `WaitUntilUsed(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_elastic_ip.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_elastic_ip.go)
:::

---

### Security Group

```go
arubaClient.FromNetwork().SecurityGroups()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SecurityGroupID()` — provider-assigned group ID
- `IsDefault()` — bool
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Boolean state*: `AsDefault()`, `NotDefault()`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_security_group.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_security_group.go)
:::

---

### Security Rule

```go
arubaClient.FromNetwork().SecurityGroupRules()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

`WithDirection` accepts `aruba.RuleDirectionIngress` or `aruba.RuleDirectionEgress`. `WithProtocol` accepts `aruba.RuleProtocolTCP`, `aruba.RuleProtocolUDP`, `aruba.RuleProtocolICMP`, or `aruba.RuleProtocolANY`.

> **Caveat**: `TargetingCIDR` and `TargetingSecurityGroup` are mutually exclusive. Setting both records a setter-time error that surfaces on `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SecurityRuleID()` — provider-assigned rule ID
- `Direction()` — `"Ingress"` or `"Egress"`
- `Protocol()` — e.g. `"TCP"`, `"UDP"`, `"ICMP"`
- `Port()` — port number or range
- `TargetKind()` — `"Ip"` or `"SecurityGroup"`
- `TargetValue()` — CIDR string or Security Group URI
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InSecurityGroup(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithDirection(RuleDirection)`, `WithProtocol(RuleProtocol)`, `WithPort(string)`
- *Active relationship*: `TargetingCIDR(string)`, `TargetingSecurityGroup(Ref)`

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### Load Balancer

```go
arubaClient.FromNetwork().LoadBalancers()
```

**Supported operations**: `List`, `Get`

Load Balancers are read-only through this SDK — they are created and managed by the Aruba Cloud platform automatically.

```go
list, err := arubaClient.FromNetwork().LoadBalancers().List(ctx, proj)
if err != nil {
    log.Fatalf("List LoadBalancers: %v", err)
}
for _, lb := range list.Items() {
    fmt.Println(lb.ID(), lb.Name(), lb.Address())
}
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `LoadBalancerID()` — provider-assigned LB ID
- `Address()` — public address
- `VPC()` — `aruba.Ref` to the attached VPC
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `Raw()` — underlying wire struct

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPC Peering

```go
arubaClient.FromNetwork().VPCPeerings()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPCPeeringID()` — provider-assigned peering ID
- `VPCID()` — source VPC ID
- `PeerVPC()` — `aruba.Ref` to the peer VPC
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPC(Ref)`
- *Geography*: `InRegion(Region)`
- *Active relationship*: `PeeredWith(Ref)`

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPC Peering Route

```go
arubaClient.FromNetwork().VPCPeeringRoutes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CIDR()` — route CIDR block
- `Target()` — `aruba.Ref` to the route target
- `VPCPeeringID()` — parent peering ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPCPeering(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithLocalCIDR(string)`, `WithRemoteCIDR(string)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPN Tunnel

```go
arubaClient.FromNetwork().VPNTunnels()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

VPN Tunnel sub-builders:
- `aruba.NewVPNIKE()` — IKE phase 1 parameters (`WithEncryption(IKEEncryption)`, `WithHash(IKEHash)`, `WithDHGroup(IKEDHGroup)`, `WithDPDAction(IKEDPDAction)`)
- `aruba.NewVPNESP()` — ESP phase 2 parameters (`WithEncryption(ESPEncryption)`, `WithHash(ESPHash)`, `WithPFS(ESPPFSGroup)`)
- `aruba.NewVPNPSK()` — pre-shared key configuration (`WithKey(string)`, `WithCloudSite(string)`, `WithOnPremSite(string)`)

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPNTunnelID()` — provider-assigned tunnel ID
- `PeerClientPublicIP()` — remote peer gateway IP
- `IKE()` — `*aruba.VPNIKE` IKE settings
- `ESP()` — `*aruba.VPNESP` ESP settings
- `PSK()` — `*aruba.VPNPSK` PSK settings
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfType(VPNType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithVPNClientProtocol(VPNClientProtocol)`, `WithPeerClientPublicIP(string)`
- *Attached config*: `WithIPConfig(*VPNIPConfig)`, `WithIKESettings(*VPNIKE)`, `WithESPSettings(*VPNESP)`, `WithPSKSettings(*VPNPSK)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

### VPN Route

```go
arubaClient.FromNetwork().VPNRoutes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `CIDR()` — route CIDR block
- `Target()` — `aruba.Ref` to the route target
- `VPNTunnelID()` — parent VPN Tunnel ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InVPNTunnel(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `WithCloudSubnet(string)`, `WithOnPremSubnet(string)`

:::tip Runnable example
Exercised as part of the orchestrator: [`examples/all-resources/orchestrator_create.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/orchestrator_create.go)
:::

---

## Schedule

### Job

```go
arubaClient.FromSchedule().Jobs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

Use `OneShotAt(t time.Time)` to schedule a one-shot job, or `WithCron(expr string)` for a recurring job on a cron schedule. Use `RecurringUntil(t time.Time)` to set an end date for a recurring job.

```go
// One-shot job — fires once at a specific time
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
fmt.Printf("✓ Job: %s (type: %s)\n", job.Name(), job.JobType())

// Recurring job — fires on a cron schedule
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
fmt.Printf("✓ Recurring Job: %s (cron: %s)\n", cronJob.Name(), cronJob.Cron())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `JobID()` — provider-assigned job ID
- `JobType()` — job type (`types.JobTypeOneShot` or `types.JobTypeRecurring`)
- `Cron()` — cron expression (recurring jobs)
- `IsEnabled()` — bool
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfType(JobType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `OneShotAt(time.Time)`, `StartingAt(time.Time)`, `WithCron(string)`, `RecurringUntil(time.Time)`, `WithSteps(...*JobStep)`
- *Boolean state*: `Enabled()`, `Disabled()`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_job.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_job.go)
:::

---

## Security

### KMS (Key Management Service)

```go
arubaClient.FromSecurity().KMS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KMSID()` — provider-assigned KMS instance ID
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

### Key

```go
arubaClient.FromSecurity().Keys()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

`OfAlgorithm` accepts `aruba.KeyAlgorithmAes` or `aruba.KeyAlgorithmRsa` (typed constants — no string cast needed).

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
fmt.Printf("✓ Key: %s (algorithm: %s)\n", key.Name(), key.Algorithm())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KeyID()` — provider-assigned key ID
- `Algorithm()` — algorithm string
- `Type()` — `"Symmetric"` or `"Asymmetric"`
- `Status()` — key lifecycle status
- `CreationSource()` — how the key was created
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfAlgorithm(KeyAlgorithm)`
- *Name*: `Named(string)`
- *Containment*: `InKMS(Ref)`

:::tip Runnable example
Key operations are covered in the KMS example: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

### Kmip

```go
arubaClient.FromSecurity().Kmips()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`. KMIP's `WaitUntilReady` succeeds on either `"CertificateAvailable"` or `"Active"`. `WaitUntilCertificateAvailable` is an alias for `WaitUntilReady`.

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

**Download the KMIP certificate** (requires a hydrated wrapper):

```go
cert, err := km.Download(ctx)
if err != nil {
    log.Fatalf("Download Kmip certificate: %v", err)
}
fmt.Println("Cert:", cert.Cert())
fmt.Println("Key:",  cert.Key())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KmipID()` — provider-assigned KMIP ID
- `KmipStatus()` — KMIP-specific status
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilCertificateAvailable(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Containment*: `InKMS(Ref)`

:::tip Runnable example
KMIP operations are covered in the KMS example: [`examples/all-resources/resource_kms.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_kms.go)
:::

---

## Storage

### Block Storage (Volume)

```go
arubaClient.FromStorage().Volumes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

`OfType` accepts `aruba.BlockStorageTypeStandard` or `aruba.BlockStorageTypePerformance`. Use `AsBootable()` to mark a volume as bootable; `NotBootable()` to unset. Use `FromImage(imageID)` to specify a base image.

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

To create a volume **from a snapshot**, use `FromSnapshot(snapshot)` instead of `FromImage`:

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BlockStorageID()` — provider-assigned volume ID
- `SizeGB()` — size in GB
- `Type()` — storage type
- `Zone()` — availability zone
- `BillingPeriod()` — billing cadence
- `IsBootable()` — bool
- `Image()` — image reference
- `SnapshotURI()` — source snapshot URI (if created from snapshot)
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilNotUsed(ctx, opts...)`, `WaitUntilUsed(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfType(BlockStorageType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`, `InZone(Zone)`
- *Descriptive scalars*: `SizedGB(int)`
- *Origin*: `FromImage(string)`, `FromSnapshot(Ref)`
- *Boolean state*: `AsBootable()`, `NotBootable()`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_block_storage.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_block_storage.go)
:::

---

### Snapshot

```go
arubaClient.FromStorage().Snapshots()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SnapshotID()` — provider-assigned snapshot ID
- `SizeGB()` — snapshot size in GB
- `Type()` — storage type
- `Zone()` — availability zone
- `BillingPeriod()` — billing cadence
- `Bootable()` — bool
- `VolumeURI()` — source volume URI
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Origin*: `FromVolume(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_snapshot.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_snapshot.go)
:::

---

### Storage Backup

```go
arubaClient.FromStorage().Backups()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

`OfType` accepts `aruba.StorageBackupTypeFull` or `aruba.StorageBackupTypeIncremental`. Use `FromVolume(vol)` to specify the source volume.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BackupID()` — provider-assigned backup ID
- `Type()` — backup type
- `RetentionDays()` — retention period in days
- `OriginURI()` — source volume URI
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Classifier*: `OfType(StorageBackupType)`
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `InProject(Ref)`
- *Geography*: `InRegion(Region)`
- *Descriptive scalars*: `RetainedForDays(int)`
- *Origin*: `FromVolume(Ref)`
- *Billing*: `BilledBy(BillingPeriod)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_storage_backup.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_storage_backup.go)
:::

---

### Storage Restore

```go
arubaClient.FromStorage().Restores()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilReady(ctx)` after `Create`.

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

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `RestoreID()` — provider-assigned restore ID
- `TargetURI()` — target volume URI
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilReady(ctx, opts...)`, `WaitUntilActive(ctx, opts...)`, `WaitUntilStates(ctx, []types.State{...}, opts...)`, `WaitUntilGone(ctx, opts...)`
- `Raw()` — underlying wire struct

**Setters**:
- *Name*: `Named(string)`
- *Labels*: `Tagged(...string)`, `Untagged(...string)`, `RetaggedAs(...string)`
- *Containment*: `FromBackup(Ref)`
- *Geography*: `InRegion(Region)`
- *Origin*: `ToVolume(Ref)`

:::tip Runnable example
Full end-to-end example: [`examples/all-resources/resource_storage_restore.go`](https://github.com/Arubacloud/sdk-go/blob/main/examples/all-resources/resource_storage_restore.go)
:::

---

## Call Options

Pass call options as variadic arguments to any `List`, `Get`, `Create`, `Update`, or `Delete` call:

| Option | Purpose |
|--------|---------|
| `aruba.WithFilter(expr)` | Server-side filter expression |
| `aruba.WithSort(expr)` | Sort expression |
| `aruba.WithLimit(n)` | Page size |
| `aruba.WithOffset(n)` | Pagination offset |
| `aruba.WithProjection(expr)` | Field projection |
| `aruba.WithAPIVersion(v)` | Override API version for this call |

See [Filters](./filters) for filter and sort syntax.

---

## Enum Constants

All enum types are re-exported from `pkg/aruba` — no extra import needed. The canonical list is in `pkg/aruba/aliases.go`.

### Region and Zone

| Constant | Description |
|----------|-------------|
| `aruba.RegionITBGBergamo` | Bergamo (Italy) datacenter |
| `aruba.ZoneITBG1` | Bergamo availability zone 1 |
| `aruba.ZoneITBG2` | Bergamo availability zone 2 |
| `aruba.ZoneITBG3` | Bergamo availability zone 3 |

### Billing

| Constant | Description |
|----------|-------------|
| `aruba.BillingPeriodHour` | Hourly billing |
| `aruba.BillingPeriodMonth` | Monthly billing |

### Network

| Constant | Value |
|----------|-------|
| `aruba.RuleDirectionIngress` | `"Ingress"` |
| `aruba.RuleDirectionEgress` | `"Egress"` |
| `aruba.RuleProtocolTCP` | `"TCP"` |
| `aruba.RuleProtocolUDP` | `"UDP"` |
| `aruba.RuleProtocolICMP` | `"ICMP"` |
| `aruba.RuleProtocolANY` | (wildcard — any protocol) |
| `aruba.SubnetTypeBasic` | `"Basic"` |
| `aruba.SubnetTypeAdvanced` | `"Advanced"` |

### Compute

| Constant | Description |
|----------|-------------|
| `aruba.CloudServerFlavorCSO1A2` | 1 vCPU, 2 GB RAM |
| `aruba.CloudServerFlavorCSO2A4` | 2 vCPU, 4 GB RAM |
| `aruba.CloudServerFlavorCSO4A8` | 4 vCPU, 8 GB RAM |
| `aruba.CloudServerFlavorCSO8A16` | 8 vCPU, 16 GB RAM |
| … (see `aliases.go` for full list) | |

### Container

| Constant | Description |
|----------|-------------|
| `aruba.KubernetesVersion1323` | Kubernetes 1.32.3 |
| `aruba.KubernetesVersion1332` | Kubernetes 1.33.2 |
| `aruba.KubernetesVersion1341` | Kubernetes 1.34.1 |
| `aruba.NodePoolInstanceK2A4` | 2 vCPU, 4 GB RAM |
| `aruba.NodePoolInstanceK4A8` | 4 vCPU, 8 GB RAM |
| `aruba.NodePoolInstanceK8A16` | 8 vCPU, 16 GB RAM |
| … (see `aliases.go` for full list) | |
| `aruba.ContainerRegistrySizeFlavorSmall` | Small concurrent-users tier |
| `aruba.ContainerRegistrySizeFlavorMedium` | Medium concurrent-users tier |
| `aruba.ContainerRegistrySizeFlavorHighPerf` | High-performance tier |

### Database

| Constant | Description |
|----------|-------------|
| `aruba.DatabaseEngineMySQL80` | MySQL 8.0 |
| `aruba.DatabaseEngineMSSQL2022Web` | SQL Server 2022 Web |
| `aruba.DatabaseEngineMSSQL2022Standard` | SQL Server 2022 Standard |
| `aruba.DatabaseEngineMSSQL2022Enterprise` | SQL Server 2022 Enterprise |
| `aruba.DBaaSFlavorDBO1A2` | 1 vCPU, 2 GB RAM |
| `aruba.DBaaSFlavorDBO2A4` | 2 vCPU, 4 GB RAM |
| `aruba.DBaaSFlavorDBO4A8` | 4 vCPU, 8 GB RAM |
| … (see `aliases.go` for full list) | |

### Storage

| Constant | Value |
|----------|-------|
| `aruba.BlockStorageTypeStandard` | `"Standard"` |
| `aruba.BlockStorageTypePerformance` | `"Performance"` |
| `aruba.StorageBackupTypeFull` | `"Full"` |
| `aruba.StorageBackupTypeIncremental` | `"Incremental"` |

### Security

| Constant | Value |
|----------|-------|
| `aruba.KeyAlgorithmAes` | `"Aes"` |
| `aruba.KeyAlgorithmRsa` | `"Rsa"` |
| `aruba.KeyTypeSymmetric` | `"Symmetric"` |
| `aruba.KeyTypeAsymmetric` | `"Asymmetric"` |
| `aruba.ServiceStatusCertificateAvailable` | `"CertificateAvailable"` |

### VPN Crypto

| Constant | Description |
|----------|-------------|
| `aruba.IKEEncryptionAES256` | AES-256 CBC (IKE phase 1) |
| `aruba.IKEHashSHA256` | HMAC-SHA-256 (IKE phase 1) |
| `aruba.IKEDHGroup14` | MODP-2048 Diffie-Hellman group |
| `aruba.ESPEncryptionAES256` | AES-256 CBC (ESP phase 2) |
| `aruba.ESPHashSHA256` | HMAC-SHA-256 (ESP phase 2) |
| `aruba.ESPPFSGroupEnable` | PFS enabled (DH group negotiated) |
| `aruba.ESPPFSGroupDisable` | PFS disabled |
| … (see `aliases.go` for full lists) | |

### Schedule

| Constant | Value |
|----------|-------|
| `aruba.JobTypeOneShot` | `"OneShot"` |
| `aruba.JobTypeRecurring` | `"Recurring"` |

---

## Appendix: Raw Wire Types (`pkg/types`)

The following types are the underlying wire-level structs. You normally access them only via `.Raw()` or `.RawRequest()` on a wrapper, or when building advanced integrations with `pkg/async`. They are also re-exported as `aruba.XxxRequest` / `aruba.XxxResponse` type aliases so you can reference them without an extra import.

| Type | File | Notes |
|------|------|-------|
| `Response[T]` | `resource.go` | Generic HTTP envelope returned by low-level adapters |
| `ErrorResponse` | `error.go` | RFC 7807 structured error |
| `ListResponse` | `resource.go` | Pagination links and total count |
| `ResourceMetadataRequest` | `resource.go` | Name + tags for Create |
| `RegionalResourceMetadataRequest` | `resource.go` | Extends metadata with Location |
| `ResourceMetadataResponse` | `resource.go` | ID, URI, Name, timestamps |
| `ResourceStatusResponse` | `resource.go` | State field |
| `ReferenceResourceCommon` | `resource.go` | `{uri: "…"}` link to another resource |
| `RequestParameters` | `parameters.go` | Low-level filter/sort/limit/offset struct (prefer `CallOption` helpers) |
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
