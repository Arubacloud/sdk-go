# SDK Architecture Diagrams — Quick Reference

All diagrams use [Mermaid](https://mermaid.js.org/) syntax. Render in GitHub, VS Code (Mermaid extension), or [mermaid.live](https://mermaid.live).

---

## D1 — Layered Architecture Overview

```mermaid
graph TB
    subgraph UserCode["User Code"]
        UC["import pkg/aruba"]
    end

    subgraph PublicLayer["Public Layer (pkg/)"]
        PKG["pkg/aruba\nClient · Options · Wrappers · List · Ref · Aliases"]
        TYPES["pkg/types\nRequest / Response / State / ErrorResponse"]
        ASYNC["pkg/async\nWaitFor · AsyncClient"]
        MULTI["pkg/multitenant\nMultitenant fleet manager"]
    end

    subgraph InternalLayer["Internal Layer (internal/)"]
        CLIENTS["internal/clients/<service>\n10 service implementations"]
        REST["internal/restclient\nDoRequest (HTTP execution)"]
        AUTH["internal/impl/auth\nTokenManager · Repositories · OAuth2"]
        MW["internal/impl/interceptor\nMiddleware chain"]
    end

    subgraph External["External Systems"]
        API["Aruba Cloud REST API"]
        VAULT["HashiCorp Vault"]
        REDIS["Redis token store"]
    end

    UC --> PKG
    PKG --> TYPES & ASYNC & CLIENTS
    CLIENTS --> REST
    REST --> MW --> AUTH
    AUTH --> VAULT & REDIS
    REST --> API

    style PublicLayer fill:#dbeafe,stroke:#3b82f6
    style InternalLayer fill:#fef3c7,stroke:#f59e0b
    style External fill:#dcfce7,stroke:#22c55e
```

---

## D2 — Client Construction Flow

```mermaid
flowchart TD
    A["NewClient(options)"] --> B["options.validate()"]
    B --> D["buildRESTClient()"]
    D --> D1["buildHTTPClient()"]
    D --> D2["buildLogger()"]
    D --> D3["buildMiddleware()\n+ TokenManager.BindTo (last)"]
    D1 & D2 & D3 --> E["Build 10 service group clients"]
    E --> G["aruba.Client  ← returned to caller"]

    style A fill:#3b82f6,color:#fff
    style G fill:#22c55e,color:#fff
```

---

## D3 — OAuth2 Token Injection (Double-Checked Locking)

```mermaid
sequenceDiagram
    participant REQ as HTTP Request
    participant TM as TokenManager
    participant REPO as TokenRepository
    participant CONN as OAuth2 Connector

    REQ->>TM: Intercept(ctx, req)
    TM->>REPO: FetchToken() [read lock]
    alt token valid
        REPO-->>TM: token
    else token missing/expired
        TM->>TM: acquire write lock
        TM->>REPO: FetchToken() [double-check]
        alt goroutine already refreshed
            REPO-->>TM: new token (reuse)
        else still stale
            TM->>CONN: RequestToken()
            CONN-->>TM: new token
            TM->>REPO: SaveToken() + increment ticket
        end
    end
    TM->>REQ: inject Authorization: Bearer <token>
```

---

## D4 — HTTP Request Lifecycle

```mermaid
sequenceDiagram
    participant C as Caller
    participant A as Adapter
    participant RC as restclient
    participant TM as TokenManager
    participant API as Aruba Cloud API

    C->>A: Create(ctx, wrapper)
    A->>A: Err() check · validate IDs · toRequest()
    A->>RC: DoRequest(ctx, POST, path, body)
    RC->>RC: Build URL · log · create request
    RC->>RC: Attach query params · set Content-Type
    RC->>TM: middleware.Intercept → inject Bearer
    TM-->>RC: authorized request
    RC->>API: httpClient.Do(req)
    API-->>RC: HTTP response
    RC->>RC: log status · re-wrap body
    RC-->>A: *http.Response
    A->>A: ParseResponseBody · populateHTTPEnvelope · fromResponse
    A-->>C: (*Wrapper, nil)
```

---

## D5 — Wrapper Triplet Pattern

```mermaid
graph LR
    subgraph File["resource_cloud_server.go"]
        W["WRAPPER\nfluent builder + mixins\nNewCloudServer().Named().InProject()…"]
        I["LOW-LEVEL INTERFACE\ncloudServersLowLevelClient\n(mockable contract)"]
        AD["ADAPTER\nbridges wrapper ↔ internal/clients"]
    end

    W -->|toRequest()| AD
    AD -->|fromResponse()| W
    AD --> I --> IMPL["internal/clients/compute\ncloudServersClientImpl"]

    style W fill:#dbeafe,stroke:#3b82f6
    style I fill:#fef3c7,stroke:#f59e0b
    style AD fill:#dcfce7,stroke:#22c55e
```

---

## D6 — Mixin Composition

```mermaid
graph TD
    EM["errMixin\nError accumulator"] 
    MM["metadataMixin\nName + tags"]
    RM["regionalMixin\nRegion"]
    ZM["zonalMixin\n→ regionalMixin\nRegion + Zone"]
    REM["responseMetadataMixin\nID · URI · CreatedAt"]
    HEM["httpEnvelopeMixin\nStatusCode · RawBody"]
    REFM["refreshMixin\nrefresh callback\nWaitUntilGone()"]
    STAM["statusMixin\n→ refreshMixin\nWaitUntilActive/Ready/States"]
    PSM["projectScopedMixin\nInProject()"]
    VSM["vpcScopedMixin\n→ projectScopedMixin"]

    ZM --> RM
    STAM --> REFM
    VSM --> PSM

    FA["Family A Wrapper\n(CloudServer, VPC, DBaaS, …)"]
    FA --- EM & MM & ZM & STAM & REM & HEM & PSM

    FB["Family B Wrapper\n(Database, Key, User, …)"]
    FB --- EM & HEM & REFM

    style FA fill:#ede9fe,stroke:#7c3aed
    style FB fill:#fce7f3,stroke:#ec4899
```

---

## D7 — Async Polling State Machine

```mermaid
stateDiagram-v2
    [*] --> Polling: WaitUntilStates(ctx, targets)

    Polling --> Success: state ∈ targets
    Polling --> TerminalFailure: state.IsFailure()
    Polling --> Retry: state empty OR IsTransitory()
    Polling --> TerminalError: settled non-target state

    Retry --> Polling: sleep baseDelay (10s default)
    Retry --> Timeout: context deadline exceeded

    Success --> [*]
    TerminalFailure --> [*]
    TerminalError --> [*]
    Timeout --> [*]
```

---

## D8 — Service Group Map

```mermaid
mindmap
  root((aruba.Client))
    FromCompute
      CloudServers
      KeyPairs
    FromNetwork
      VPCs
      Subnets
      SecurityGroups
      ElasticIPs
      LoadBalancers
      VPNTunnels
      VPCPeerings
    FromStorage
      BlockStorages
      Snapshots
      StorageBackups
    FromDatabase
      DBaaS
      Databases
      Users
      Grants
    FromContainer
      KaaS
      ContainerRegistries
    FromSecurity
      KMS
      Keys
      Kmip
    FromProject
      Projects
    FromAudit
      AuditEvents
    FromMetric
      Alerts
    FromSchedule
      Jobs
```

---

## D9 — Multi-Tenant Fleet

```mermaid
graph TD
    TMPL["Template Options\n(shallow-copied singletons:\n*http.Client · logger · middleware)"]
    MT["Multitenant\nmap[tenantID → {client, lastUsage}]"]

    TMPL -->|NewWithTemplate| MT

    MT -->|New('tenant-a')| CA["aruba.Client\ntenant-a"]
    MT -->|New('tenant-b')| CB["aruba.Client\ntenant-b"]
    MT -->|New('tenant-c')| CC["aruba.Client\ntenant-c"]

    CLEAN["StartCleanupRoutine\ntick: 1h · idle threshold: 24h"]
    CLEAN -->|CleanUp()| MT

    style MT fill:#dbeafe,stroke:#3b82f6
    style CLEAN fill:#fef3c7,stroke:#f59e0b
```

---

## D10 — Error Handling Flow

```mermaid
flowchart TD
    SETTER["Setter call\ne.g. .InProject(nil)"]
    SETTER --> ERRMIXIN["addErr() → errMixin\nchain continues"]

    ADAPTER["Adapter.Create()"]
    ADAPTER --> CHECK{"Err() != nil?"}
    CHECK -->|yes| RET1["return wrapper, accumulated errors"]
    CHECK -->|no| VALID["Validate IDs\nfmt.Errorf(...)"]
    VALID --> HTTP["DoRequest()"]
    HTTP --> STATUS{"IsSuccess()?"}
    STATUS -->|yes| PARSE["ParseResponseBody → Data"]
    STATUS -->|no| HTTPERR["&HTTPError{StatusCode, Body, ErrResp}\npopulateHTTPEnvelope (envelope retained)"]

    PARSE --> OK["return wrapper, nil"]
    HTTPERR --> ERRRET["return wrapper, *HTTPError"]

    style OK fill:#22c55e,color:#fff
    style RET1 fill:#ef4444,color:#fff
    style ERRRET fill:#ef4444,color:#fff
```
