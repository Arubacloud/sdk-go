# SDK Data Types

The Aruba Cloud SDK for Go uses a consistent set of data types for requests and responses across all API groups. This document outlines the generic, shared types that form the basis of the SDK's models, followed by a breakdown of the specific types for each API group.

## Generic & Shared Types

Most resource-specific types are composed of these fundamental building blocks. Understanding them is key to using the SDK effectively.

<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>File</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>Response[T]</code></td>
      <td><code>resource.go</code></td>
      <td>A generic wrapper for all API call responses. It contains either the successfully parsed data (<code>Data *T</code>) for 2xx responses or a structured error (<code>Error *ErrorResponse</code>) for 4xx/5xx responses.</td>
    </tr>
    <tr>
      <td><code>ErrorResponse</code></td>
      <td><code>error.go</code></td>
      <td>A struct that follows RFC 7807 Problem Details for HTTP APIs. It provides structured error information, including a title, detail, and status code.</td>
    </tr>
    <tr>
      <td><code>ListResponse</code></td>
      <td><code>resource.go</code></td>
      <td>An embedded struct used in all list responses (e.g., <code>VPCList</code>, <code>CloudServerList</code>). It provides pagination details, including the total number of items and links to the next/previous pages.</td>
    </tr>
    <tr>
      <td><code>ResourceMetadataRequest</code></td>
      <td><code>resource.go</code></td>
      <td>A standard struct for naming a resource and assigning tags during creation. It is typically embedded in a resource's specific <code>...Request</code> type.</td>
    </tr>
    <tr>
      <td><code>RegionalResourceMetadataRequest</code></td>
      <td><code>resource.go</code></td>
      <td>Extends <code>ResourceMetadataRequest</code> by adding a mandatory <code>Location</code> field. Used for creating resources that must be tied to a specific data center region.</td>
    </tr>
    <tr>
      <td><code>ResourceMetadataResponse</code></td>
      <td><code>resource.go</code></td>
      <td>The standard metadata structure returned for every resource. It includes the resource's unique <code>ID</code>, <code>URI</code>, <code>Name</code>, <code>Location</code>, creation/update timestamps, and more.</td>
    </tr>
    <tr>
      <td><code>ResourceStatus</code></td>
      <td><code>resource.go</code></td>
      <td>A common struct returned with resource responses, indicating the current state of the resource (e.g., "Active", "Creating", "Error").</td>
    </tr>
    <tr>
      <td><code>ReferenceResource</code></td>
      <td><code>resource.go</code></td>
      <td>A simple struct used to link to another resource by its unique <code>URI</code>. This is commonly used in request bodies to specify dependencies.</td>
    </tr>
    <tr>
        <td><code>RequestParameters</code></td>
        <td><code>parameters.go</code></td>
        <td>A struct used to provide optional query parameters for API calls, such as filtering, sorting, pagination (limit/offset), and API versioning.</td>
    </tr>
  </tbody>
</table>

## Types by API Group

The following sections detail the primary request and response types for each API group.

### Project Types
*File: `project.project.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>ProjectRequest</code></td>
      <td>Create/Update</td>
      <td>The payload for creating or updating a project.</td>
    </tr>
    <tr>
      <td><code>ProjectResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single project resource.</td>
    </tr>
    <tr>
      <td><code>ProjectList</code></td>
      <td>List</td>
      <td>Represents a paginated list of projects.</td>
    </tr>
  </tbody>
</table>

### Audit Types
*File: `audit.event.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>AuditEvent</code></td>
      <td>List</td>
      <td>Represents a single audit event record.</td>
    </tr>
    <tr>
      <td><code>AuditEventListResponse</code></td>
      <td>List</td>
      <td>Represents a paginated list of audit events.</td>
    </tr>
  </tbody>
</table>

### Compute Types
*Files: `compute.cloudserver.go`, `compute.keypair.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>CloudServerRequest</code></td>
      <td>Create/Update</td>
      <td>The payload for creating or updating a Cloud Server. The <code>Properties</code> field contains <code>CloudServerPropertiesRequest</code>, which includes an optional <code>UserData</code> field (nullable string) for base64-encoded cloud-init content to initialize the server.</td>
    </tr>
    <tr>
      <td><code>CloudServerResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single Cloud Server resource.</td>
    </tr>
    <tr>
      <td><code>CloudServerList</code></td>
      <td>List</td>
      <td>Represents a paginated list of Cloud Servers.</td>
    </tr>
    <tr>
      <td><code>CloudServerPasswordRequest</code></td>
      <td>SetPassword</td>
      <td>The payload for setting or changing a Cloud Server password.</td>
    </tr>
    <tr>
      <td><code>KeyPairRequest</code></td>
      <td>Create</td>
      <td>The payload for creating an SSH Key Pair.</td>
    </tr>
    <tr>
      <td><code>KeyPairResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single SSH Key Pair resource.</td>
    </tr>
    <tr>
      <td><code>KeyPairListResponse</code></td>
      <td>List</td>
      <td>Represents a paginated list of SSH Key Pairs.</td>
    </tr>
  </tbody>
</table>

### Container Types
*Files: `container.kaas.go`, `container.containerregistry.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>KaaSRequest</code></td>
      <td>Create</td>
      <td>The payload for creating a KaaS cluster.</td>
    </tr>
    <tr>
      <td><code>KaaSUpdateRequest</code></td>
      <td>Update</td>
      <td>The payload for updating a KaaS cluster (limited fields).</td>
    </tr>
    <tr>
      <td><code>KaaSResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single KaaS cluster resource.</td>
    </tr>
    <tr>
      <td><code>KaaSList</code></td>
      <td>List</td>
      <td>Represents a paginated list of KaaS clusters.</td>
    </tr>
    <tr>
      <td><code>KaaSKubeconfigResponse</code></td>
      <td>DownloadKubeconfig</td>
      <td>Represents the kubeconfig file download response with filename and base64 content.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryRequest</code></td>
      <td>Create/Update</td>
      <td>The payload for creating or updating a Container Registry.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single Container Registry resource.</td>
    </tr>
    <tr>
      <td><code>ContainerRegistryList</code></td>
      <td>List</td>
      <td>Represents a paginated list of Container Registries.</td>
    </tr>
  </tbody>
</table>

### Database Types
*Files: `database.dbaas.go`, `database.database.go`, etc.*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>DBaaSRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating/updating a DBaaS instance.</td>
    </tr>
    <tr>
      <td><code>DBaaSResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single DBaaS instance.</td>
    </tr>
    <tr>
      <td><code>DatabaseRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating/updating a database within a DBaaS instance.</td>
    </tr>
    <tr>
      <td><code>UserRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating/updating a user for a DBaaS instance.</td>
    </tr>
    <tr>
      <td><code>GrantRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for managing user permissions on a database.</td>
    </tr>
    <tr>
      <td><code>BackupRequest</code></td>
      <td>Create</td>
      <td>Payload for creating a backup of a DBaaS instance.</td>
    </tr>
  </tbody>
</table>

### Metric Types
*Files: `metrics.alert.go`, `metrics.metric.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>AlertResponse</code></td>
      <td>List</td>
      <td>Represents a single monitoring alert.</td>
    </tr>
    <tr>
      <td><code>AlertsListResponse</code></td>
      <td>List</td>
      <td>Represents a paginated list of alerts.</td>
    </tr>
    <tr>
      <td><code>MetricResponse</code></td>
      <td>List</td>
      <td>Represents a set of time-series data points for a specific metric.</td>
    </tr>
    <tr>
      <td><code>MetricListResponse</code></td>
      <td>List</td>
      <td>Represents a list of metrics.</td>
    </tr>
  </tbody>
</table>

### Network Types
*Files: `network.vpc.go`, `network.subnet.go`, etc.*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>VPCRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a Virtual Private Cloud.</td>
    </tr>
    <tr>
      <td><code>SubnetRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a Subnet within a VPC.</td>
    </tr>
    <tr>
      <td><code>ElasticIPRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating an Elastic IP.</td>
    </tr>
    <tr>
      <td><code>SecurityGroupRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a Security Group.</td>
    </tr>
    <tr>
      <td><code>SecurityRuleRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a rule within a Security Group.</td>
    </tr>
    <tr>
      <td><code>LoadBalancerResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single Load Balancer resource.</td>
    </tr>
     <tr>
      <td><code>VPCPeeringRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a VPC Peering connection.</td>
    </tr>
     <tr>
      <td><code>VPNTunnelRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a VPN Tunnel.</td>
    </tr>
  </tbody>
</table>

### Schedule Types
*File: `schedule.job.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>JobRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a scheduled job.</td>
    </tr>
    <tr>
      <td><code>JobResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single scheduled job.</td>
    </tr>
    <tr>
      <td><code>JobList</code></td>
      <td>List</td>
      <td>Represents a paginated list of scheduled jobs.</td>
    </tr>
  </tbody>
</table>

### Security Types
*File: `security.kms.go`*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>KmsRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a KMS key.</td>
    </tr>
    <tr>
      <td><code>KmsResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single KMS key.</td>
    </tr>
    <tr>
      <td><code>KmsList</code></td>
      <td>List</td>
      <td>Represents a paginated list of KMS keys.</td>
    </tr>
  </tbody>
</table>

### Storage Types
*Files: `storage.block-storage.go`, `storage.snapshot.go`, etc.*
<table>
  <thead>
    <tr>
      <th>Type Name</th>
      <th>Usage</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>BlockStorageRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a block storage volume.</td>
    </tr>
    <tr>
      <td><code>BlockStorageResponse</code></td>
      <td>Get/List</td>
      <td>Represents a single block storage volume.</td>
    </tr>
    <tr>
      <td><code>SnapshotRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a volume snapshot.</td>
    </tr>
    <tr>
      <td><code>StorageBackupRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for creating or updating a volume backup.</td>
    </tr>
    <tr>
      <td><code>RestoreRequest</code></td>
      <td>Create/Update</td>
      <td>Payload for restoring a backup to a volume.</td>
    </tr>
  </tbody>
</table>

