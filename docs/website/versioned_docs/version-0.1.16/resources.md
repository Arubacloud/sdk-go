# API Groups and Resources

The Aruba Cloud SDK is organized into several API groups, each corresponding to a specific service area (e.g., Compute, Network, Storage). You can access these groups from the main `arubaClient` object. This document provides a comprehensive list of all available groups and the resources they manage.

## Project

The Project is the top-level resource under which all other resources are organized. The client for managing projects is accessed directly.

<table>
  <thead>
    <tr>
      <th>Client Accessor</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>arubaClient.FromProject()</code></td>
      <td>Manages projects, which are containers for all other cloud resources.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Audit

Manages access to audit trail information.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Events()</code></td>
      <td>Provides access to audit events for a project.</td>
      <td><code>List</code></td>
    </tr>
  </tbody>
</table>

## Compute

Manages virtual machines and related resources.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.CloudServers()</code></td>
      <td>Manages virtual machine instances (Cloud Servers).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code>, <code>PowerOn</code>, <code>PowerOff</code>, <code>SetPassword</code></td>
    </tr>
    <tr>
      <td><code>.KeyPairs()</code></td>
      <td>Manages SSH key pairs for server access.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Container

Manages container-based services.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.KaaS()</code></td>
      <td>Manages Kubernetes as a Service (KaaS) clusters.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code>, <code>DownloadKubeconfig</code></td>
    </tr>
    <tr>
      <td><code>.ContainerRegistry()</code></td>
      <td>Manages private container registries.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Database

Manages Database as a Service (DBaaS) and its sub-resources.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.DBaaS()</code></td>
      <td>Manages DBaaS instances (e.g., MySQL, PostgreSQL).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Databases()</code></td>
      <td>Manages individual databases within a DBaaS instance.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Backups()</code></td>
      <td>Manages backups of DBaaS instances.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Users()</code></td>
      <td>Manages database users for a DBaaS instance.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Grants()</code></td>
      <td>Manages user permissions (grants) on databases.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Metric

Provides access to monitoring data and alerts.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Alerts()</code></td>
      <td>Provides access to monitoring alerts.</td>
      <td><code>List</code></td>
    </tr>
    <tr>
      <td><code>.Metrics()</code></td>
      <td>Provides access to time-series monitoring data for resources.</td>
      <td><code>List</code></td>
    </tr>
  </tbody>
</table>

## Network

Manages all networking resources.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.VPCs()</code></td>
      <td>Manages Virtual Private Clouds (VPCs).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Subnets()</code></td>
      <td>Manages subnets within a VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.ElasticIPs()</code></td>
      <td>Manages public, static IP addresses (Elastic IPs).</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.SecurityGroups()</code></td>
      <td>Manages security groups (firewalls) within a VPC.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.SecurityGroupRules()</code></td>
      <td>Manages individual rules within a security group.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.LoadBalancers()</code></td>
      <td>Manages load balancers.</td>
      <td><code>List</code>, <code>Get</code></td>
    </tr>
    <tr>
      <td><code>.VPCPeerings()</code></td>
      <td>Manages peering connections between two VPCs.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPCPeeringRoutes()</code></td>
      <td>Manages routes for a VPC peering connection.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPNTunnels()</code></td>
      <td>Manages Site-to-Site VPN tunnels.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.VPNRoutes()</code></td>
      <td>Manages routes for a VPN tunnel.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Schedule

Manages scheduled, automated jobs.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Jobs()</code></td>
      <td>Manages scheduled jobs (one-shot or recurring) that can perform actions on resources.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Security

Manages security-related services.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.KMSKeys()</code></td>
      <td>Manages Key Management Service (KMS) keys.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

## Storage

Manages block storage and related data protection resources.

<table>
  <thead>
    <tr>
      <th>Resource Client</th>
      <th>Description</th>
      <th>Available Operations</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>.Volumes()</code></td>
      <td>Manages block storage volumes.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Snapshots()</code></td>
      <td>Manages point-in-time snapshots of block storage volumes.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Backups()</code></td>
      <td>Manages backups of block storage volumes.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
    <tr>
      <td><code>.Restores()</code></td>
      <td>Manages the restoration of a backup to a volume.</td>
      <td><code>Create</code>, <code>List</code>, <code>Get</code>, <code>Update</code>, <code>Delete</code></td>
    </tr>
  </tbody>
</table>

