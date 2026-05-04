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
        WithName("my-resource").
        AddTag("env:prod").
        WithFoo(...))

// 3. Wait for async resources to become ready
if err := result.WaitUntilActive(ctx); err != nil { … }

// 4. Read response accessors
fmt.Println(result.ID(), result.Name(), result.State())
```

- `aruba.NewX()` — factory constructor for every resource builder
- `IntoFoo(ref)` — binds the parent scope; accepts any `aruba.Ref` (hydrated wrapper or `aruba.URI("…")`)
- `WithFoo(...)` — fluent setters; errors are deferred until `Create`/`Update`
- `WaitUntilActive(ctx, opts...)` — available on resources marked **async** below; see [Async / Await](./async) for full options
- `aruba.URI(s)` — wraps a raw string path into a `Ref` (see [API Walkthrough](./walkthrough#5-get-a-specific-resource))

---

## Project

```go
arubaClient.FromProject()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`

> Project is **not** async — it is synchronously ready after `Create` returns. No `WaitUntilActive` call is needed.

```go
proj, err := arubaClient.FromProject().Create(
    ctx,
    aruba.NewProject().
        WithName("my-project").
        WithDescription("Production project").
        AddTag("env:prod").
        WithDefault(false))
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

---

## Compute

### Cloud Server

```go
arubaClient.FromCompute().CloudServers()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`, `PowerOn`, `PowerOff`, `SetPassword`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

A Cloud Server depends on network resources (VPC, Subnet, Security Group), an Elastic IP, a Boot Volume (Block Storage), and a Key Pair. Create those first and pass the hydrated wrappers as `Ref` parameters.

```go
cs, err := arubaClient.FromCompute().CloudServers().Create(
    ctx,
    aruba.NewCloudServer().
        IntoProject(proj).
        WithName("my-server").
        AddTag("env:prod").
        InRegion("ITBG-Bergamo").
        InZone("ITBG-1").
        WithFlavor("CSO2A4").
        WithVPC(vpc).
        AddSubnet(subnet).
        AddSecurityGroup(sg).
        WithElasticIP(eip).
        WithBootVolume(blockStorage).
        WithKeyPair(keyPair))
if err != nil {
    log.Fatalf("Create Cloud Server: %v", err)
}

if err := cs.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Cloud Server did not become Active: %v", err)
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

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
        IntoProject(proj).
        WithName("my-keypair").
        InRegion("ITBG-Bergamo").
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

---

## Container

### KaaS (Kubernetes as a Service)

```go
arubaClient.FromContainer().KaaS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`, `DownloadKubeconfig`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
k, err := arubaClient.FromContainer().KaaS().Create(
    ctx,
    aruba.NewKaaS().
        IntoProject(proj).
        WithName("my-cluster").
        AddTag("env:prod").
        WithLocation("ITBG-Bergamo").
        WithVPC(vpc).
        WithSubnet(subnet).
        WithSecurityGroupName("my-security-group").
        WithNodeCIDR("10.100.0.0/16", "node-cidr").
        WithPodCIDR("10.200.0.0/16").
        WithKubernetesVersion("1.32.3").
        WithHA(true).
        WithBillingPeriod("Hour").
        AddNodePool(aruba.NewNodePool().
            Named("default-pool").
            WithCount(3).
            OfInstance("K4A8").
            InZone("ITBG-1")))
if err != nil {
    log.Fatalf("Create KaaS: %v", err)
}

if err := k.WaitUntilActive(ctx); err != nil {
    log.Fatalf("KaaS did not become Active: %v", err)
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Container Registry

```go
arubaClient.FromContainer().ContainerRegistry()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
reg, err := arubaClient.FromContainer().ContainerRegistry().Create(
    ctx,
    aruba.NewContainerRegistry().
        IntoProject(proj).
        WithName("my-registry").
        AddTag("env:prod").
        WithVPC(vpc).
        WithSubnet(subnet).
        WithSecurityGroup(sg).
        WithPublicIP(eip).
        WithBlockStorage(blockStorage).
        WithAdminUsername("admin").
        WithSize(100).
        WithBillingPeriod("Hour"))
if err != nil {
    log.Fatalf("Create ContainerRegistry: %v", err)
}

if err := reg.WaitUntilActive(ctx); err != nil {
    log.Fatalf("ContainerRegistry did not become Active: %v", err)
}
fmt.Printf("✓ Registry: %s (public IP: %s)\n", reg.Name(), reg.PublicIP())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `ContainerRegistryID()` — provider-assigned registry ID
- `PublicIP()` — public endpoint IP
- `VPC()`, `Subnet()`, `SecurityGroup()`, `BlockStorage()` — `aruba.Ref` to attached resources
- `AdminUsername()` — registry admin user
- `ConcurrentUsers()` — configured concurrent user limit
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

## Database

### DBaaS (Database as a Service)

```go
arubaClient.FromDatabase().DBaaS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
db, err := arubaClient.FromDatabase().DBaaS().Create(
    ctx,
    aruba.NewDBaaS().
        IntoProject(proj).
        WithName("my-database").
        AddTag("env:prod").
        InRegion("ITBG-Bergamo").
        InZone("ITBG-1").
        WithEngine("mysql-8.0").
        WithFlavor("DBO2A4").
        WithStorage(20).
        WithBillingPeriod("Hour").
        WithAutoscaling(true).
        WithNetworking(vpc, subnet, sg, eip))
if err != nil {
    log.Fatalf("Create DBaaS: %v", err)
}

if err := db.WaitUntilActive(ctx); err != nil {
    log.Fatalf("DBaaS did not become Active: %v", err)
}
fmt.Printf("✓ DBaaS: %s (engine: %s)\n", db.Name(), db.Engine())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DBaaSID()` — provider-assigned instance ID
- `Engine()` — engine slug (e.g. `"mysql-8.0"`)
- `EngineRaw()` — full engine struct
- `Flavor()` — flavor slug
- `FlavorRaw()` — full flavor struct
- `Storage()` — storage size in GB
- `Autoscaling()` — bool
- `VPC()`, `Subnet()`, `SecurityGroup()`, `ElasticIP()` — `aruba.Ref` to networking resources
- `BillingPeriod()` — billing cadence
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Database

```go
arubaClient.FromDatabase().Databases()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
database, err := arubaClient.FromDatabase().Databases().Create(
    ctx,
    aruba.NewDatabase().
        IntoDBaaS(db).
        WithName("my-app-db").
        AddTag("app:backend"))
if err != nil {
    log.Fatalf("Create Database: %v", err)
}

if err := database.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Database did not become Active: %v", err)
}
fmt.Printf("✓ Database: %s\n", database.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `DatabaseID()` — provider-assigned database ID
- `DBaaSID()` — parent DBaaS ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### User

```go
arubaClient.FromDatabase().Users()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
user, err := arubaClient.FromDatabase().Users().Create(
    ctx,
    aruba.NewUser().
        IntoDBaaS(db).
        WithName("app_user").
        WithPassword("Str0ngP@ssword!").
        AddTag("app:backend"))
if err != nil {
    log.Fatalf("Create User: %v", err)
}

if err := user.WaitUntilActive(ctx); err != nil {
    log.Fatalf("User did not become Active: %v", err)
}
fmt.Printf("✓ User: %s\n", user.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `UserID()` — provider-assigned user ID
- `DBaaSID()` — parent DBaaS ID
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Grant

```go
arubaClient.FromDatabase().Grants()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
grant, err := arubaClient.FromDatabase().Grants().Create(
    ctx,
    aruba.NewGrant().
        IntoDatabase(database).
        WithName("app_user-grant").
        WithPrivileges("ALL"))
if err != nil {
    log.Fatalf("Create Grant: %v", err)
}

if err := grant.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Grant did not become Active: %v", err)
}
fmt.Printf("✓ Grant: %s (privileges: %s)\n", grant.Name(), grant.Privileges())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `GrantID()` — provider-assigned grant ID
- `DatabaseID()` — parent Database ID
- `Privileges()` — privilege string
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### DBaaS Backup

```go
arubaClient.FromDatabase().DBaaSBackups()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
backup, err := arubaClient.FromDatabase().DBaaSBackups().Create(
    ctx,
    aruba.NewDBaaSBackup().
        IntoProject(proj).
        WithName("my-db-backup").
        WithDBaaS(db).
        WithType("Full").
        WithRetentionDays(30).
        AddTag("backup"))
if err != nil {
    log.Fatalf("Create DBaaSBackup: %v", err)
}

if err := backup.WaitUntilActive(ctx); err != nil {
    log.Fatalf("DBaaS Backup did not become Active: %v", err)
}
fmt.Printf("✓ DBaaS Backup: %s\n", backup.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BackupID()` — provider-assigned backup ID
- `DBaaSID()` — source DBaaS ID
- `Type()` — backup type string
- `RetentionDays()` — retention period
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

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

---

## Network

### VPC

```go
arubaClient.FromNetwork().VPCs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
vpc, err := arubaClient.FromNetwork().VPCs().Create(
    ctx,
    aruba.NewVPC().
        IntoProject(proj).
        WithName("my-vpc").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithDefault(false).
        WithPreset(false))
if err != nil {
    log.Fatalf("Create VPC: %v", err)
}

if err := vpc.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPC did not become Active: %v", err)
}
fmt.Printf("✓ VPC: %s\n", vpc.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPCID()` — provider-assigned VPC ID
- `Region()` — region slug
- `IsDefault()`, `IsPreset()` — flags
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Subnet

```go
arubaClient.FromNetwork().Subnets()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

`WithType` accepts `string(aruba.SubnetTypeBasic)` or `string(aruba.SubnetTypeAdvanced)`.

`aruba.NewSubnetDHCP()` is a sub-builder for DHCP configuration. Attach it with `WithDHCP(...)`.

```go
subnet, err := arubaClient.FromNetwork().Subnets().Create(
    ctx,
    aruba.NewSubnet().
        IntoVPC(vpc).
        WithName("my-subnet").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithType(string(aruba.SubnetTypeAdvanced)).
        WithDefault(false).
        WithCIDR("192.168.1.0/25").
        WithDHCP(aruba.NewSubnetDHCP().
            Enabled().
            WithRange("192.168.1.10", 50).
            AddRoute("0.0.0.0/0", "192.168.1.1").
            AddDNS("8.8.8.8").
            AddDNS("8.8.4.4")))
if err != nil {
    log.Fatalf("Create Subnet: %v", err)
}

if err := subnet.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Subnet did not become Active: %v", err)
}
fmt.Printf("✓ Subnet: %s (CIDR: %s)\n", subnet.Name(), subnet.CIDR())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SubnetID()` — provider-assigned subnet ID
- `Type()` — subnet type string
- `CIDR()` — CIDR block
- `DHCP()` — DHCP configuration
- `IsDefault()` — bool
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Elastic IP

```go
arubaClient.FromNetwork().ElasticIPs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
eip, err := arubaClient.FromNetwork().ElasticIPs().Create(
    ctx,
    aruba.NewElasticIP().
        IntoProject(proj).
        WithName("my-eip").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithBillingPeriod("Hour"))
if err != nil {
    log.Fatalf("Create ElasticIP: %v", err)
}

if err := eip.WaitUntilActive(ctx); err != nil {
    log.Fatalf("ElasticIP did not become Active: %v", err)
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Security Group

```go
arubaClient.FromNetwork().SecurityGroups()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
sg, err := arubaClient.FromNetwork().SecurityGroups().Create(
    ctx,
    aruba.NewSecurityGroup().
        IntoVPC(vpc).
        WithName("my-security-group").
        AddTag("security").
        WithDefault(false))
if err != nil {
    log.Fatalf("Create SecurityGroup: %v", err)
}

if err := sg.WaitUntilActive(ctx); err != nil {
    log.Fatalf("SecurityGroup did not become Active: %v", err)
}
fmt.Printf("✓ Security Group: %s (ID: %s)\n", sg.Name(), sg.ID())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SecurityGroupID()` — provider-assigned group ID
- `Default()` — bool
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Security Rule

```go
arubaClient.FromNetwork().SecurityGroupRules()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

`WithDirection` accepts `string(aruba.RuleDirectionIngress)` or `string(aruba.RuleDirectionEgress)`.

> **Caveat**: `WithTargetCIDR` and `WithTargetSecurityGroup` are mutually exclusive. Setting both records a setter-time error that surfaces on `Create`.

```go
rule, err := arubaClient.FromNetwork().SecurityGroupRules().Create(
    ctx,
    aruba.NewSecurityRule().
        IntoSecurityGroup(sg).
        WithName("allow-ssh").
        AddTag("ssh").
        InRegion("ITBG-Bergamo").
        WithDirection(string(aruba.RuleDirectionIngress)).
        WithProtocol("TCP").
        WithPort("22").
        WithTargetCIDR("0.0.0.0/0"))
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
- `Raw()` — underlying wire struct

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

---

### VPC Peering

```go
arubaClient.FromNetwork().VPCPeerings()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
peering, err := arubaClient.FromNetwork().VPCPeerings().Create(
    ctx,
    aruba.NewVPCPeering().
        IntoVPC(vpc).
        WithName("my-peering").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithPeerVPC(aruba.URI("/projects/"+peerProjectID+"/vpcs/"+peerVPCID)))
if err != nil {
    log.Fatalf("Create VPCPeering: %v", err)
}

if err := peering.WaitUntilActive(ctx); err != nil {
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### VPC Peering Route

```go
arubaClient.FromNetwork().VPCPeeringRoutes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
route, err := arubaClient.FromNetwork().VPCPeeringRoutes().Create(
    ctx,
    aruba.NewVPCPeeringRoute().
        IntoVPCPeering(peering).
        WithName("my-peering-route").
        AddTag("network").
        InRegion("ITBG-Bergamo").
        WithCIDR("10.0.0.0/8").
        WithTarget(aruba.URI("/projects/"+projectID+"/vpcs/"+vpcID)))
if err != nil {
    log.Fatalf("Create VPCPeeringRoute: %v", err)
}

if err := route.WaitUntilActive(ctx); err != nil {
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### VPN Tunnel

```go
arubaClient.FromNetwork().VPNTunnels()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

VPN Tunnel sub-builders:
- `aruba.NewVPNIKE()` — IKE phase 1 parameters
- `aruba.NewVPNESP()` — ESP phase 2 parameters
- `aruba.NewVPNPSK()` — pre-shared key configuration

```go
tunnel, err := arubaClient.FromNetwork().VPNTunnels().Create(
    ctx,
    aruba.NewVPNTunnel().
        IntoProject(proj).
        WithName("my-vpn-tunnel").
        AddTag("vpn").
        InRegion("ITBG-Bergamo").
        WithRemoteGateway("203.0.113.1").
        WithIKE(aruba.NewVPNIKE().
            WithEncryption(string(aruba.VPNEncryptionAES256)).
            WithHash(string(aruba.VPNHashSHA256)).
            WithDHGroup(string(aruba.VPNDHGroup14))).
        WithESP(aruba.NewVPNESP().
            WithEncryption(string(aruba.VPNEncryptionAES256)).
            WithHash(string(aruba.VPNHashSHA256))).
        WithPSK(aruba.NewVPNPSK().
            WithKey("my-pre-shared-key").
            WithID("tunnel-id")))
if err != nil {
    log.Fatalf("Create VPNTunnel: %v", err)
}

if err := tunnel.WaitUntilActive(ctx); err != nil {
    log.Fatalf("VPNTunnel did not become Active: %v", err)
}
fmt.Printf("✓ VPN Tunnel: %s (gateway: %s)\n", tunnel.Name(), tunnel.RemoteGateway())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `VPNTunnelID()` — provider-assigned tunnel ID
- `RemoteGateway()` — remote gateway IP
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### VPN Route

```go
arubaClient.FromNetwork().VPNRoutes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
vpnRoute, err := arubaClient.FromNetwork().VPNRoutes().Create(
    ctx,
    aruba.NewVPNRoute().
        IntoVPNTunnel(tunnel).
        WithName("my-vpn-route").
        AddTag("vpn").
        InRegion("ITBG-Bergamo").
        WithCIDR("10.0.0.0/8").
        WithTarget(aruba.URI("/projects/"+projectID+"/vpcs/"+vpcID)))
if err != nil {
    log.Fatalf("Create VPNRoute: %v", err)
}

if err := vpnRoute.WaitUntilActive(ctx); err != nil {
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

## Schedule

### Job

```go
arubaClient.FromSchedule().Jobs()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

`WithType` accepts `string(aruba.JobTypeOneShot)` or `string(aruba.JobTypeRecurring)`.

For recurring jobs, `WithRecurrence` accepts `aruba.RecurrenceTypeHourly`, `aruba.RecurrenceTypeDaily`, `aruba.RecurrenceTypeWeekly`, `aruba.RecurrenceTypeMonthly`, or `aruba.RecurrenceTypeCustom`.

```go
job, err := arubaClient.FromSchedule().Jobs().Create(
    ctx,
    aruba.NewJob().
        IntoProject(proj).
        WithName("my-job").
        AddTag("automation").
        WithType(string(aruba.JobTypeRecurring)).
        WithRecurrence(string(aruba.RecurrenceTypeDaily)))
if err != nil {
    log.Fatalf("Create Job: %v", err)
}
fmt.Printf("✓ Job: %s (type: %s)\n", job.Name(), job.Type())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `JobID()` — provider-assigned job ID
- `Type()` — job type string
- `Recurrence()` — recurrence type string
- `Steps()` — configured job steps
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `Raw()` — underlying wire struct

---

## Security

### KMS (Key Management Service)

```go
arubaClient.FromSecurity().KMS()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
kms, err := arubaClient.FromSecurity().KMS().Create(
    ctx,
    aruba.NewKMS().
        IntoProject(proj).
        WithName("my-kms").
        AddTag("security").
        InRegion("ITBG-Bergamo").
        WithBillingPeriod("Hour"))
if err != nil {
    log.Fatalf("Create KMS: %v", err)
}

if err := kms.WaitUntilActive(ctx); err != nil {
    log.Fatalf("KMS did not become Active: %v", err)
}
fmt.Printf("✓ KMS: %s (ID: %s)\n", kms.Name(), kms.ID())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `KMSID()` — provider-assigned KMS instance ID
- `BillingPeriod()` — billing cadence
- `Region()` — region slug
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Key

```go
arubaClient.FromSecurity().Keys()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — `State()` and `FailureReason()` are available.

`WithAlgorithm` accepts `string(aruba.KeyAlgorithmAes)` or `string(aruba.KeyAlgorithmRsa)`.

```go
key, err := arubaClient.FromSecurity().Keys().Create(
    ctx,
    aruba.NewKey().
        IntoKMS(kms).
        WithName("my-encryption-key").
        AddTag("security").
        WithAlgorithm(string(aruba.KeyAlgorithmAes)))
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
- `Raw()` — underlying wire struct

---

### Kmip

```go
arubaClient.FromSecurity().Kmips()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` (or wait for `"CertificateAvailable"`) after `Create`.

```go
km, err := arubaClient.FromSecurity().Kmips().Create(
    ctx,
    aruba.NewKmip().
        IntoKMS(kms).
        WithName("my-kmip").
        AddTag("security"))
if err != nil {
    log.Fatalf("Create Kmip: %v", err)
}

// Wait for certificate to be available
if err := km.WaitUntilState(ctx, string(aruba.ServiceStatusCertificateAvailable)); err != nil {
    log.Fatalf("Kmip certificate not available: %v", err)
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
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

## Storage

### Block Storage (Volume)

```go
arubaClient.FromStorage().Volumes()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

`WithType` accepts `aruba.BlockStorageTypeStandard` or `aruba.BlockStorageTypePerformance`.

```go
bs, err := arubaClient.FromStorage().Volumes().Create(
    ctx,
    aruba.NewBlockStorage().
        IntoProject(proj).
        WithName("my-volume").
        AddTag("storage").
        InRegion("ITBG-Bergamo").
        InZone("ITBG-1").
        WithSize(20).
        WithType(aruba.BlockStorageTypeStandard).
        WithBillingPeriod("Hour").
        WithBootable(true).
        WithImage("LU22-001"))
if err != nil {
    log.Fatalf("Create BlockStorage: %v", err)
}

if err := bs.WaitUntilActive(ctx); err != nil {
    log.Fatalf("BlockStorage did not become Active: %v", err)
}
fmt.Printf("✓ Volume: %s (%d GB)\n", bs.Name(), bs.Size())
```

To create a volume **from a snapshot**, use `FromSnapshot(snapshot)` instead of `WithImage`:

```go
bs, err := arubaClient.FromStorage().Volumes().Create(
    ctx,
    aruba.NewBlockStorage().
        IntoProject(proj).
        WithName("restored-volume").
        InRegion("ITBG-Bergamo").
        InZone("ITBG-1").
        WithSize(20).
        WithType(aruba.BlockStorageTypeStandard).
        WithBillingPeriod("Hour").
        WithBootable(true).
        FromSnapshot(snapshot))
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BlockStorageID()` — provider-assigned volume ID
- `Size()` — size in GB
- `Type()` — storage type string
- `Zone()` — availability zone
- `BillingPeriod()` — billing cadence
- `Bootable()` — bool
- `Image()` — image reference
- `SnapshotURI()` — source snapshot URI (if created from snapshot)
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Snapshot

```go
arubaClient.FromStorage().Snapshots()
```

**Supported operations**: `Create`, `List`, `Get`, `Update`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
snap, err := arubaClient.FromStorage().Snapshots().Create(
    ctx,
    aruba.NewSnapshot().
        IntoProject(proj).
        WithName("my-snapshot").
        AddTag("backup").
        InRegion("ITBG-Bergamo").
        WithBillingPeriod("Hour").
        OfVolume(bs))
if err != nil {
    log.Fatalf("Create Snapshot: %v", err)
}

if err := snap.WaitUntilActive(ctx); err != nil {
    log.Fatalf("Snapshot did not become Active: %v", err)
}
fmt.Printf("✓ Snapshot: %s\n", snap.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `SnapshotID()` — provider-assigned snapshot ID
- `Size()` — snapshot size in GB
- `Type()` — storage type
- `Zone()` — availability zone
- `BillingPeriod()` — billing cadence
- `Bootable()` — bool
- `VolumeURI()` — source volume URI
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Storage Backup

```go
arubaClient.FromStorage().Backups()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

`WithType` accepts `aruba.StorageBackupTypeFull` or `aruba.StorageBackupTypeIncremental`.

```go
backup, err := arubaClient.FromStorage().Backups().Create(
    ctx,
    aruba.NewStorageBackup().
        IntoProject(proj).
        WithName("my-backup").
        AddTag("backup").
        WithOrigin(bs).
        WithType(aruba.StorageBackupTypeFull).
        WithRetentionDays(30).
        WithBillingPeriod("Hour"))
if err != nil {
    log.Fatalf("Create StorageBackup: %v", err)
}

if err := backup.WaitUntilActive(ctx); err != nil {
    log.Fatalf("StorageBackup did not become Active: %v", err)
}
fmt.Printf("✓ Storage Backup: %s\n", backup.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `BackupID()` — provider-assigned backup ID
- `Type()` — backup type string
- `RetentionDays()` — retention period
- `OriginURI()` — source volume URI
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

---

### Storage Restore

```go
arubaClient.FromStorage().Restores()
```

**Supported operations**: `Create`, `List`, `Get`, `Delete`
**Async**: yes — call `WaitUntilActive(ctx)` after `Create`.

```go
restore, err := arubaClient.FromStorage().Restores().Create(
    ctx,
    aruba.NewStorageRestore().
        IntoBackup(backup).
        WithName("my-restore").
        AddTag("restore").
        WithTarget(aruba.URI(backup.OriginURI())))
if err != nil {
    log.Fatalf("Create StorageRestore: %v", err)
}

if err := restore.WaitUntilActive(ctx); err != nil {
    log.Fatalf("StorageRestore did not become Active: %v", err)
}
fmt.Printf("✓ Storage Restore: %s\n", restore.Name())
```

**Response accessors**:
- `ID()`, `URI()`, `Name()`, `Tags()`
- `RestoreID()` — provider-assigned restore ID
- `TargetURI()` — target volume URI
- `State()`, `FailureReason()`, `PreviousState()`, `IsDisabled()`, `DisableReasons()`
- `WaitUntilActive(ctx, opts...)`, `WaitUntilState(ctx, target, opts...)`
- `Raw()` — underlying wire struct

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

All enum types are re-exported from `pkg/aruba` — no extra import needed.

### Network

| Constant | Value |
|----------|-------|
| `aruba.RuleDirectionIngress` | `"Ingress"` |
| `aruba.RuleDirectionEgress` | `"Egress"` |
| `aruba.EndpointTypeIP` | `"Ip"` |
| `aruba.EndpointTypeSecurityGroup` | `"SecurityGroup"` |
| `aruba.SubnetTypeBasic` | `"Basic"` |
| `aruba.SubnetTypeAdvanced` | `"Advanced"` |

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
| `aruba.KeyStatusActive` | `"Active"` |
| `aruba.KeyStatusInCreation` | `"InCreation"` |
| `aruba.ServiceStatusActive` | `"Active"` |
| `aruba.ServiceStatusCertificateAvailable` | `"CertificateAvailable"` |

### Schedule

| Constant | Value |
|----------|-------|
| `aruba.JobTypeOneShot` | `"OneShot"` |
| `aruba.JobTypeRecurring` | `"Recurring"` |
| `aruba.RecurrenceTypeHourly` | `"Hourly"` |
| `aruba.RecurrenceTypeDaily` | `"Daily"` |
| `aruba.RecurrenceTypeWeekly` | `"Weekly"` |
| `aruba.RecurrenceTypeMonthly` | `"Monthly"` |

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
| `ResourceStatus` | `resource.go` | State field |
| `ReferenceResource` | `resource.go` | `{uri: "…"}` link to another resource |
| `RequestParameters` | `parameters.go` | Low-level filter/sort/limit/offset struct (prefer `CallOption` helpers) |
| `ProjectRequest` / `ProjectResponse` / `ProjectList` | `project.project.go` | |
| `VPCRequest` / `VPCResponse` / `VPCList` | `network.vpc.go` | |
| `SubnetRequest` / `SubnetResponse` / `SubnetList` | `network.subnet.go` | |
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
| `JobRequest` / `JobResponse` / `JobList` | `schedule.job.go` | |
| `AlertResponse` / `AlertsListResponse` | `metrics.alert.go` | |
| `MetricResponse` / `MetricListResponse` | `metrics.metric.go` | |
| `AuditEvent` / `AuditEventListResponse` | `audit.event.go` | |
| `VPCPeeringRequest` / `VPCPeeringResponse` | `network.vpc-peering.go` | |
| `VPNTunnelRequest` / `VPNTunnelResponse` | `network.vpn-tunnel.go` | |
| `VPNRouteRequest` / `VPNRouteResponse` | `network.vpn-route.go` | |
| `LoadBalancerResponse` | `network.load-balancer.go` | |
