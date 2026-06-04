# Aruba Cloud Go SDK — Architecture & Capabilities
**Technical Presentation** · June 2026

---

## Table of Contents

1. [What is the SDK?](#1-what-is-the-sdk)
2. [High-Level Architecture](#2-high-level-architecture)
3. [Package Structure](#3-package-structure)
4. [Single-Import Design Principle](#4-single-import-design-principle)
5. [Client Construction](#5-client-construction)
6. [Authentication Subsystem](#6-authentication-subsystem)
7. [HTTP Request Lifecycle](#7-http-request-lifecycle)
8. [Service Groups](#8-service-groups)
9. [Wrapper Layer — Triplet Pattern](#9-wrapper-layer--triplet-pattern)
10. [Mixin System](#10-mixin-system)
11. [Resource Families: A vs. B](#11-resource-families-a-vs-b)
12. [Async Polling & Wait Helpers](#12-async-polling--wait-helpers)
13. [Multi-Tenant Client Management](#13-multi-tenant-client-management)
14. [Key Design Highlights](#14-key-design-highlights)

---

## 1. What is the SDK?

The **Aruba Cloud Go SDK** (`github.com/Arubacloud/sdk-go`) is the official Go client library for the Aruba Cloud API.

| Attribute | Value |
|---|---|
| Language | Go 1.24+ |
| Module | `github.com/Arubacloud/sdk-go` |
| Status | Alpha (API surface may change) |
| Auth | OAuth2 Client Credentials (RFC 6749) |
| Service Groups | 10 (Compute, Network, Storage, Database, Container, Security, Project, Audit, Metrics, Schedule) |
| Wire Format | JSON / REST |

### Core Value Proposition

- **Single-import ergonomics** — `import "github.com/Arubacloud/sdk-go/pkg/aruba"` covers 99.9% of use cases
- **Fluent, chainable builders** — resource construction reads like plain English
- **Built-in async polling** — `WaitUntilActive()`, `WaitUntilReady()`, `WaitUntilStates()` included on every resource
- **Multi-tenant ready** — `pkg/multitenant` manages fleets of per-tenant clients out of the box
- **Extensible auth** — supports static tokens, OAuth2 client credentials, memory/file/Redis/Vault backends

---

## 2. High-Level Architecture

```mermaid
graph TB
    subgraph UserCode["User Code"]
        UC["import pkg/aruba"]
    end

    subgraph PublicLayer["Public Layer (pkg/)"]
        PKG["pkg/aruba\nClient · Options · Wrappers\nList · Ref · Aliases"]
        TYPES["pkg/types\nRequest / Response types\nState · ErrorResponse"]
        ASYNC["pkg/async\nWaitFor · AsyncClient"]
        MULTI["pkg/multitenant\nMultitenant manager"]
    end

    subgraph InternalLayer["Internal Layer (internal/)"]
        CLIENTS["internal/clients/<service>\nCloudServers · VPCs · DBaaS\n10 service implementations"]
        REST["internal/restclient\nDoRequest · HTTP execution"]
        AUTH["internal/impl/auth\nTokenManager · TokenRepo\nCredentialsRepo · Connector"]
        INTERCEPTOR["internal/impl/interceptor\nMiddleware chain"]
        LOGGER["internal/impl/logger\nLogger interface + backends"]
    end

    subgraph ExternalLayer["External Systems"]
        API["Aruba Cloud REST API"]
        VAULT["HashiCorp Vault\n(optional)"]
        REDIS["Redis\n(optional token store)"]
        FILE["File System\n(token persistence)"]
    end

    UC --> PKG
    PKG --> TYPES
    PKG --> ASYNC
    PKG --> CLIENTS
    CLIENTS --> REST
    REST --> INTERCEPTOR
    INTERCEPTOR --> AUTH
    AUTH --> VAULT
    AUTH --> REDIS
    AUTH --> FILE
    REST --> API

    style PublicLayer fill:#dbeafe,stroke:#3b82f6
    style InternalLayer fill:#fef3c7,stroke:#f59e0b
    style ExternalLayer fill:#dcfce7,stroke:#22c55e
```

---

## 3. Package Structure

```
sdk-go/
├── pkg/
│   ├── aruba/              ← Public entry point & wrapper layer
│   │   ├── aruba.go        ← NewClient() + Options
│   │   ├── builder.go      ← buildClient() orchestration
│   │   ├── resource_*.go   ← Fluent wrappers + adapters (one per resource)
│   │   ├── mixin_common.go ← errMixin, metadataMixin, httpEnvelopeMixin, …
│   │   ├── mixin_scoped.go ← projectScopedMixin, vpcScopedMixin, …
│   │   ├── mixin_status.go ← statusMixin: WaitUntilActive/Ready/States
│   │   ├── list.go         ← Generic List[T] paginated container
│   │   ├── aliases.go      ← Typed enum constants re-exported from pkg/types
│   │   ├── ref.go          ← Ref interface, extractID, parseURIIDs
│   │   └── errors.go       ← *HTTPError
│   ├── types/              ← All request/response/common data models
│   ├── async/              ← WaitFor, AsyncClient, polling primitives
│   └── multitenant/        ← Multi-tenant client manager
│
├── internal/
│   ├── clients/            ← Service-specific HTTP client impls
│   │   ├── compute/        ← Cloud servers, key pairs
│   │   ├── network/        ← VPCs, subnets, security groups, …
│   │   ├── storage/        ← Block storage, snapshots, backups
│   │   ├── database/       ← DBaaS, databases, users, grants
│   │   ├── container/      ← KaaS, node pools, registries
│   │   ├── security/       ← KMS, keys, KMIP
│   │   ├── project/        ← Projects
│   │   ├── audit/          ← Audit events
│   │   ├── metric/         ← Metrics & alerts
│   │   └── schedule/       ← Scheduled jobs
│   ├── restclient/         ← Low-level HTTP execution (DoRequest)
│   └── impl/
│       ├── auth/           ← Token manager, repositories, OAuth2 connector
│       ├── interceptor/    ← Middleware chain
│       └── logger/         ← Logger interface + backends
│
├── examples/all-resources/ ← Usage examples (canonical builder chains)
└── docs/website/           ← Docusaurus documentation site
```

---

## 4. Single-Import Design Principle

The fundamental contract: **one import covers 99.9% of real-world usage**.

```go
import "github.com/Arubacloud/sdk-go/pkg/aruba"

client, _ := aruba.NewClient(aruba.NewOptions().
    WithClientCredentials("my-id", "my-secret").
    WithBaseURL("https://api.aruba.cloud"))

cs := aruba.NewCloudServer().
    OfFlavor(aruba.CloudServerFlavorSmall).
    Named("web-01").
    Tagged("prod", "web").
    InProject(projectRef).
    InRegion(aruba.RegionItaly).
    BilledBy(aruba.BillingPeriodHour)

result, err := client.FromCompute().CloudServers().Create(ctx, cs)
```

Three mechanisms enforce this principle:

```mermaid
graph LR
    subgraph SingleImport["pkg/aruba — Single Import Guarantee"]
        A["Re-exported Enums\naliases.go\naruba.StateActive\naruba.RegionItaly\naruba.RuleProtocolTCP"]
        B["Wrapper Serialisation\nraw_marshal.go\nRaw() · RawJSON() · RawYAML()"]
        C["Wait Helpers on Wrapper\nmixin_status.go\nWaitUntilActive()\nWaitUntilReady()\nWaitUntilStates()"]
    end

    subgraph Escape["Residual pkg/types use cases\n(documented in working-at-low-level.md)"]
        D["Structured validation errors\ntypes.ValidationError"]
        E["Non-promoted wire fields\n(deep nested properties)"]
        F["pkg/async background polling\n(concurrent waits)"]
    end

    A -.->|avoids| D
    B -.->|avoids| E
    C -.->|avoids| F
```

---

## 5. Client Construction

```mermaid
flowchart TD
    A["NewClient(options)"] --> B["options.validate()"]
    B --> C{"valid?"}
    C -->|no| ERR["return error"]
    C -->|yes| D["buildRESTClient()"]

    D --> D1["buildHTTPClient()\ndefaults to http.DefaultClient"]
    D --> D2["buildLogger()\nconfigurable backends"]
    D --> D3["buildMiddleware()\nstandard interceptor chain"]
    D3 --> D3a["TokenManager.BindTo(interceptor)\n← always last in chain"]

    D1 & D2 & D3 --> E["Build 10 service group clients\nsequentially"]

    E --> F1["FromCompute()\ncloudServersClientAdapter\nkeyPairsClientAdapter"]
    E --> F2["FromNetwork()\nvpcClientAdapter\nsubnetClientAdapter\n..."]
    E --> F3["FromDatabase()\ndbaasClientAdapter\ndatabaseClientAdapter\n..."]
    E --> F4["... 7 more groups"]

    F1 & F2 & F3 & F4 --> G["aruba.Client\n(returned to caller)"]

    style A fill:#3b82f6,color:#fff
    style G fill:#22c55e,color:#fff
    style ERR fill:#ef4444,color:#fff
```

**Key injection points via `Options`:**

| Option | Default | Purpose |
|---|---|---|
| `WithCustomHTTPClient(*http.Client)` | `http.DefaultClient` | Custom transport, timeouts, TLS |
| `WithClientCredentials(id, secret)` | — | OAuth2 auto-refresh |
| `WithToken(token)` | — | Static bearer token |
| `WithCustomMiddleware(interceptor)` | `standard.NewInterceptor()` | Custom request hooks |
| `WithCustomLogger(logger)` | no-op | Structured logging |

---

## 6. Authentication Subsystem

```mermaid
graph TD
    subgraph Ports["Interfaces (internal/ports/auth/)"]
        TM["TokenManager\n+ Bind as interceptor\n+ Inject Bearer header"]
        TR["TokenRepository\n+ FetchToken\n+ SaveToken"]
        PC["ProviderConnector\n+ RequestToken (OAuth2)"]
        CR["CredentialsRepository\n+ FetchCredentials"]
    end

    subgraph TokenRepos["Token Repository Backends"]
        MEM["Memory\n(in-process, default)"]
        PROXY["Memory Proxy\n(write-through to persistent store)"]
        FILE2["File\n(.token.json, chmod 0600)"]
        REDIS2["Redis\n(distributed fleet)"]
    end

    subgraph CredRepos["Credentials Repository Backends"]
        MEMCRED["Memory\n(static ClientID + Secret)"]
        VAULT2["HashiCorp Vault\n(AppRole + KV v2)"]
    end

    TR --> MEM & PROXY & FILE2 & REDIS2
    CR --> MEMCRED & VAULT2
    PC --> OAUTH["OAuth2 Client Credentials\ngolang.org/x/oauth2/clientcredentials\nRFC 6749"]
    TM --> TR
    TM --> PC
    PC --> CR

    style TM fill:#3b82f6,color:#fff
    style OAUTH fill:#7c3aed,color:#fff
```

**Token Injection — Double-Checked Locking:**

```mermaid
sequenceDiagram
    participant REQ as HTTP Request
    participant TM as TokenManager
    participant REPO as TokenRepository
    participant CONN as ProviderConnector

    REQ->>TM: Intercept(ctx, req)
    TM->>REPO: FetchToken() [read lock]
    alt token valid
        REPO-->>TM: token
    else token missing/expired
        TM->>TM: acquire write lock
        TM->>REPO: FetchToken() [double-check]
        alt another goroutine refreshed
            REPO-->>TM: new token (reuse)
        else still stale
            TM->>CONN: RequestToken()
            CONN-->>TM: new token
            TM->>REPO: SaveToken() + increment ticket
        end
    end
    TM->>REQ: inject Authorization: Bearer <token>
```

> **Fleet safety:** `NewTokenProxyWithRandomExpirationDriftSeconds` randomises the expiry drift to prevent synchronised token-refresh storms across multiple SDK instances.

---

## 7. HTTP Request Lifecycle

```mermaid
sequenceDiagram
    participant C as Client code
    participant A as Adapter (pkg/aruba)
    participant RC as restclient.DoRequest
    participant MW as Middleware chain
    participant TM as TokenManager
    participant API as Aruba Cloud API

    C->>A: Create(ctx, cloudServer)
    A->>A: Err() check + validate IDs
    A->>A: toRequest() → wire body
    A->>RC: DoRequest(ctx, POST, /cloudServers, body, params, headers)

    RC->>RC: 1. Build full URL
    RC->>RC: 2. Log request (Bearer REDACTED)
    RC->>RC: 3. Create *http.Request with context
    RC->>RC: 4. Attach query parameters
    RC->>RC: 5. Set Content-Type: application/json
    RC->>RC: 6. Merge caller headers
    RC->>MW: 7. middleware.Intercept(ctx, req)
    MW->>TM: run token injection (last interceptor)
    TM-->>MW: Authorization: Bearer <token>
    MW-->>RC: request ready
    RC->>API: 8. httpClient.Do(req)
    API-->>RC: HTTP response

    RC->>RC: 9. Log status + headers
    RC->>RC: 10. Re-wrap body (logging consumed stream)
    RC-->>A: *http.Response

    A->>A: ParseResponseBody[T]()
    A->>A: populateHTTPEnvelope()
    A->>A: fromResponse() → hydrate wrapper
    A->>A: install refresh callback
    A-->>C: (*CloudServer, nil)
```

---

## 8. Service Groups

The root `Client` exposes **10 service group accessors**:

```mermaid
graph LR
    ROOT["aruba.Client"] --> CMP["FromCompute()\nCloudServers · KeyPairs"]
    ROOT --> NET["FromNetwork()\nVPCs · Subnets · SecurityGroups\nElasticIPs · LoadBalancers\nVPNTunnels · VPCPeerings\nVPNRoutes · VPCPeeringRoutes"]
    ROOT --> STO["FromStorage()\nBlockStorages · Snapshots\nStorageBackups · StorageRestores"]
    ROOT --> DB["FromDatabase()\nDBaaS · Databases\nUsers · Grants · DBaaSBackups"]
    ROOT --> CON["FromContainer()\nKaaS · ContainerRegistries"]
    ROOT --> SEC["FromSecurity()\nKMS · Keys · Kmip"]
    ROOT --> PRJ["FromProject()\nProjects"]
    ROOT --> AUD["FromAudit()\nAuditEvents"]
    ROOT --> MET["FromMetric()\nAlerts"]
    ROOT --> SCH["FromSchedule()\nJobs"]

    style ROOT fill:#3b82f6,color:#fff
    style CMP fill:#dbeafe
    style NET fill:#dbeafe
    style STO fill:#dbeafe
    style DB fill:#dbeafe
    style CON fill:#dbeafe
    style SEC fill:#dbeafe
    style PRJ fill:#dbeafe
    style AUD fill:#dbeafe
    style MET fill:#dbeafe
    style SCH fill:#dbeafe
```

**Call chain anatomy:**

```
arubaClient.FromCompute().CloudServers().Create(ctx, cs)
     │              │           │           │
     │              │           │           └─ ctx + fluent wrapper
     │              │           └─ CloudServersClient interface
     │              └─ ComputeClient interface
     └─ root aruba.Client
                                    ↓
                        cloudServersClientAdapter.Create()
                                    ↓
                        compute.NewCloudServersClientImpl(rest).Create()
```

---

## 9. Wrapper Layer — Triplet Pattern

Every `resource_<name>.go` follows a strict **three-section layout**:

```mermaid
graph LR
    subgraph TripletFile["resource_cloud_server.go"]
        W["WRAPPER\nChainable builder struct\n+ mixin embeds\n+ typed setters\n+ read accessors\n\nNewCloudServer()\n.Named() .InProject()\n.OfFlavor() .InRegion()\n.BilledBy() …"]
        I["LOW-LEVEL INTERFACE\nAdapter contract\n(mockable in tests)\n\ncloudServersLowLevelClient {\n  List(ctx, params)\n  Create(ctx, req)\n  Get(ctx, id, params)\n  Update(ctx, id, req)\n  Delete(ctx, id)\n}"]
        AD["ADAPTER\nBridges wrapper ↔\ninternal/clients/*\n\ncloudServersClientAdapter\n.Create() .Get()\n.Update() .Delete()\n.List()"]
    end

    W -->|toRequest()| AD
    AD -->|fromResponse()| W
    AD -->|calls| I
    I -->|implemented by| IMPL["internal/clients/compute\ncloudServersClientImpl"]

    style W fill:#dbeafe,stroke:#3b82f6
    style I fill:#fef3c7,stroke:#f59e0b
    style AD fill:#dcfce7,stroke:#22c55e
```

**Fluent builder — setter verb vocabulary:**

| Verb | Role | Example |
|---|---|---|
| `New<X>()` | Construct | `NewCloudServer()` |
| `Named(name)` | Identity | `.Named("web-01")` |
| `Tagged(…)` | Labels | `.Tagged("prod", "web")` |
| `In<Parent/Geo>` | Containment & placement | `.InProject(ref)`, `.InRegion(aruba.RegionItaly)` |
| `Of<Classifier>` | Type / sizing | `.OfFlavor(...)`, `.OfEngine(...)` |
| `From<Source>` | Origin | `.FromImage(ref)`, `.FromVolume(ref)` |
| `With<Noun>` | Attached config | `.WithVPC(ref)`, `.WithElasticIP(ref)` |
| `BilledBy(period)` | Billing | `.BilledBy(aruba.BillingPeriodHour)` |

---

## 10. Mixin System

Mixins are embedded structs providing reusable behaviour across all wrappers.

```mermaid
graph TD
    subgraph Common["mixin_common.go"]
        EM["errMixin\nError accumulator\n.Err() .addErr()"]
        MM["metadataMixin\nName + tags\n.Named() .Tagged()"]
        RM["regionalMixin\nRegion\n.InRegion()"]
        ZM["zonalMixin → regionalMixin\nRegion + Zone\n.InZone()"]
        REM["responseMetadataMixin\nID · URI · CreatedAt · Version"]
        LIM["linkedMixin\nLinkedResources()"]
        HEM["httpEnvelopeMixin\nStatusCode · Headers\nRawBody · RawError()"]
    end

    subgraph Scoped["mixin_scoped.go"]
        PSM["projectScopedMixin\n.InProject(Ref)"]
        VSM["vpcScopedMixin → PSM\n.InVPC(Ref)"]
        SGSM["securityGroupScopedMixin → VSM\n.InSecurityGroup(Ref)"]
        DBSM["dbaasScopedMixin → PSM\n.InDBaaS(Ref)"]
        DASM["databaseScopedMixin → DBSM\n.InDatabase(Ref)"]
        BSM["backupScopedMixin\n.InBackup(Ref)"]
        KSM["kmsScopedMixin → PSM\n.InKMS(Ref)"]
        VPNSM["vpnTunnelScopedMixin → VSM\n.InVPNTunnel(Ref)"]
        VPPSM["vpcPeeringScopedMixin → VSM\n.InVPCPeering(Ref)"]
    end

    subgraph StatusMixins["mixin_status.go / mixin_refresh.go"]
        REFM["refreshMixin\nrefresh callback\n.WaitUntilGone()"]
        STAM["statusMixin → refreshMixin\nState · IsDisabled()\n.WaitUntilActive()\n.WaitUntilReady()\n.WaitUntilStates()"]
    end

    STAM --> REFM
    ZM --> RM
    VSM --> PSM
    SGSM --> VSM
    DBSM --> PSM
    DASM --> DBSM

    subgraph FamilyA["Family A wrapper (e.g. CloudServer)"]
        FA["errMixin\nmetadataMixin\nzonalMixin\nstatusIxin\nresponseMetadataMixin\nlinkedMixin\nhttpEnvelopeMixin\nprojectScopedMixin"]
    end

    subgraph FamilyB["Family B wrapper (e.g. Database)"]
        FB["errMixin\ndbaasScopedMixin\nrefreshMixin\nhttpEnvelopeMixin"]
    end

    style Common fill:#dbeafe,stroke:#3b82f6
    style Scoped fill:#fef3c7,stroke:#f59e0b
    style StatusMixins fill:#dcfce7,stroke:#22c55e
    style FamilyA fill:#ede9fe,stroke:#7c3aed
    style FamilyB fill:#fce7f3,stroke:#ec4899
```

---

## 11. Resource Families: A vs. B

```mermaid
graph LR
    subgraph FamilyA["Family A — Standard Shape"]
        direction TB
        A_WIRE["Wire: Metadata { Properties { … } }"]
        A_FEAT["Features:\n✓ metadataMixin (name + tags)\n✓ regionalMixin / zonalMixin\n✓ statusMixin (WaitUntil*)\n✓ responseMetadataMixin\n✓ linkedMixin"]
        A_EX["Resources:\nCloudServer · VPC · Subnet\nBlockStorage · DBaaS · KaaS\nLoadBalancer · VPNTunnel\nJob · KMS · ElasticIP\n(large majority)"]
        A_WIRE --> A_FEAT --> A_EX
    end

    subgraph FamilyB["Family B — Flat Shape"]
        direction TB
        B_WIRE["Wire: flat JSON body (no envelope)"]
        B_FEAT["Features:\n✗ No Metadata/Properties boxing\n✗ No tags · No location\n✗ No metadataMixin\n✗ No statusMixin\n✓ refreshMixin (WaitUntilGone)"]
        B_EX["Resources:\nDatabase · Key · Kmip\nUser · Grant"]
        B_WIRE --> B_FEAT --> B_EX
    end

    subgraph FamilyBSub["Family B — No-Update Variant"]
        B2["Key · Kmip\nno Update() operation\nenforced by interface\n+ reflective test guards"]
    end

    FamilyB --> FamilyBSub

    style FamilyA fill:#dbeafe,stroke:#3b82f6
    style FamilyB fill:#fef3c7,stroke:#f59e0b
    style FamilyBSub fill:#fce7f3,stroke:#ec4899
```

**Identity quirks in Family B:**

| Resource | ID field | URI construction |
|---|---|---|
| `Database` | name IS the path identifier | client-side: `ancestor-ids + name` |
| `Key` | `KeyResponse.KeyID` | client-side from ancestor IDs |
| `Kmip` | `KmipResponse.ID` | client-side from ancestor IDs |
| `User` | `WithUsername(...)` | name IS path identifier |
| `Grant` | opaque server grant ID | recoverable from URI Ref only |

---

## 12. Async Polling & Wait Helpers

```mermaid
flowchart TD
    subgraph StatusMixin["statusMixin.WaitUntilStates(ctx, targets, opts...)"]
        direction TB
        TICK["polling tick: refresh()"]
        TICK --> R1{"state ∈ targets?"}
        R1 -->|yes| SUCCESS["✓ return nil"]
        R1 -->|no| R2{"state.IsFailure()?"}
        R2 -->|yes| FAIL["✗ terminal error"]
        R2 -->|no| R3{"empty or IsTransitory?"}
        R3 -->|yes| WAIT["sleep baseDelay → retry"]
        WAIT --> TICK
        R3 -->|no| R4["settled non-target state"]
        R4 --> TERR["✗ terminal error\n(fail fast)"]
    end

    subgraph AsyncPkg["pkg/async.WaitFor[T](ctx, retries, baseDelay, timeout, callFunc, checkFunc)"]
        GOROUTINE["goroutine retries callFunc()\nup to retries times\nfixed baseDelay between attempts\ntimeout via context deadline"]
    end

    subgraph Defaults["Default constants"]
        DEF["DefaultRetries  = 60\nDefaultBaseDelay = 10s\nDefaultTimeout  = 600s\n(10 minutes total)"]
    end

    subgraph WaitVariants["Convenience wrappers"]
        WA["WaitUntilActive()\ntarget: StateActive"]
        WR["WaitUntilReady()\ntargets: Active·Running\nStopped·NotUsed\nReserved·InUse·Used"]
        WS["WaitUntilStates(ctx, []types.State{...}, opts...)"]
        WG["WaitUntilGone()\n404 → success\nany other error → transient"]
    end

    WA & WR & WS --> StatusMixin
    WG --> AsyncPkg
    StatusMixin --> AsyncPkg
    AsyncPkg -.-> Defaults

    subgraph Overrides["WaitOption overrides"]
        OV["WithRetries(n)\nWithBaseDelay(d)\nWithTimeout(d)"]
    end

    OV -.->|passed to| WS

    style SUCCESS fill:#22c55e,color:#fff
    style FAIL fill:#ef4444,color:#fff
    style TERR fill:#ef4444,color:#fff
```

**Specialised waiters for edge cases:**

| Waiter | Resource | Trigger |
|---|---|---|
| `WaitUntilCertificateAvailable` | `*Kmip` | Polls `KmipResponse.Status` against explicit terminal map |
| `WaitUntilUsed` / `WaitUntilNotUsed` | `*BlockStorage`, `*ElasticIP` | Attach/detach lifecycle |

---

## 13. Multi-Tenant Client Management

```mermaid
graph TD
    subgraph MT["pkg/multitenant.Multitenant"]
        MAP["map[tenantID → entry]\n{ client, lastUsage }\nsync.RWMutex"]
        NEW["New(tenant)\ndeep-copy template Options\nbuild aruba.Client"]
        GET["Get / MustGet / GetOrNil\nupdates lastUsage on access"]
        ADD["Add(tenant, aruba.Client)\nregister pre-built client"]
        CLEAN["CleanUp(from duration)\nevict idle clients"]
        ROUTINE["StartCleanupRoutine(ctx, tick, from)\nbackground goroutine\ndefault: tick=1h, idle=24h"]
    end

    subgraph Template["NewWithTemplate(*aruba.Options)"]
        TC["Deep-copied per tenant:\n• slices (DNS, tags, …)"]
        TS["Shallow-copied (shared singletons):\n• *http.Client\n• logger\n• middleware"]
    end

    subgraph Usage["Typical Fleet Usage"]
        U1["mt := multitenant.NewWithTemplate(baseOpts)"]
        U2["client := mt.New('tenant-abc')"]
        U3["client.FromCompute().CloudServers().List(ctx, ...)"]
        U1 --> U2 --> U3
    end

    MAP --> GET
    MAP --> CLEAN
    NEW --> MAP
    ADD --> MAP
    ROUTINE --> CLEAN
    Template --> NEW

    style MT fill:#dbeafe,stroke:#3b82f6
    style Template fill:#fef3c7,stroke:#f59e0b
    style Usage fill:#dcfce7,stroke:#22c55e
```

---

## 14. Key Design Highlights

### Design Decisions Summary

| Decision | Rationale |
|---|---|
| **Single public entry point** (`pkg/aruba`) | Eliminates import sprawl; callers rarely need internal packages |
| **Triplet pattern** (Wrapper / Interface / Adapter) | Testability — adapter tested via mock; wrapper tested independently |
| **Mixin composition over inheritance** | Go doesn't have inheritance; mixins compose cleanly and stay independent |
| **Fluent chainable builders** | Readable intent; error accumulation means broken chains don't panic |
| **Double-checked locking in token manager** | Goroutine-safe refresh without serialising every request |
| **Fixed-delay polling** (no exponential backoff) | Predictable behaviour; Aruba API operations have known durations |
| **Family A / Family B split** | Enforces wire-shape discipline; prevents accidental cross-family code reuse |
| **`Ref` interface** (`URI() + ID()`) | Decouples adapters from typed wrappers; enables `aruba.URI("/...")` escape hatch |
| **Response-preferring getters** | `Get → display` works without re-setting fields; server data always wins |
| **`fromResponse` round-trip invariant** | `Get → toRequest() → PUT` roundtrip preserves all server-side fields |

### Error Handling Model

```
Setter-time errors   → accumulated in errMixin (never panic, chain continues)
                         checked at adapter entry: if err := X.Err(); err != nil { return X, err }

Validation errors    → returned as Go error values (pre-HTTP)
                         fmt.Errorf("project cannot be empty")

API errors (4xx/5xx) → unmarshaled into resp.Error (*types.ErrorResponse, RFC 7807)
                         returned as *HTTPError; wrapper retains envelope for diagnostics

Network errors       → surfaced from httpClient.Do(); no special wrapping
```

### Compile-Time Safety

- All enums are **typed strings** (`type State string`, `type Region string`) — no raw string magic
- All enum constants are re-exported in `aliases.go` under the `aruba.*` namespace
- Family B "no-Update" contract enforced by both **service-group interface** definition and **reflective test guards**
- Deep parent chains validated at adapter entry, not silently dropped on the wire

---

*Generated June 2026 · `github.com/Arubacloud/sdk-go`*
