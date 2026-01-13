package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

/*************  ✨ Windsurf Command ⭐  *************/
// main demonstrates a modular approach for creating cloud infrastructure.
// The code is organized into focused, reusable functions rather than one monolithic main function.
// The example shows how to create all resources, poll for their Active state, and then print a summary.
//
// Usage:
//   go run . [flags]
//
// Flags:
//   -mode string
//        Operation mode: create, update, or delete (default "create")
//
// Examples:
//   go run .                                  # Create resources
//   PROJECT_ID=my-project go run . -mode=update    # Update resources
//   PROJECT_ID=my-project go run . -mode=delete    # Delete resources
/*******  916ef78b-f0e3-4a8d-8711-6783ddf0996d  *******/

const (
	defaultClientID     = "client_id"
	defaultClientSecret = "client_secret"
)

func main() {
	mode := flag.String("mode", "create", "Operation mode: create, update, or delete")
	flag.Parse()

	switch *mode {
	case "create":
		runCreateExample(defaultClientID, defaultClientSecret)
	case "update":
		runUpdateExample(defaultClientID, defaultClientSecret)
	case "delete":
		runDeleteExample(defaultClientID, defaultClientSecret)
	default:
		log.Fatalf("Unknown mode: %s. Use 'create', 'update', or 'delete'", *mode)
	}
}

func runCreateExample(clientID, clientSecret string) {
	// Initialize the SDK (automatically obtains JWT token)
	arubaClient, err := aruba.NewClient(aruba.DefaultOptions(clientID, clientSecret))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with timeout - increased to handle multiple resource creations and polling
	// With multiple resources polling (30 attempts × 5s each = 150s per resource), we need sufficient time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	fmt.Println("\n=== SDK Create Example ===")

	// Create all resources
	resources := createAllResources(ctx, arubaClient)

	// Print summary
	printResourceSummary(resources)
}

// ResourceCollection holds all created resources
type ResourceCollection struct {
	ProjectID         string
	ElasticIPResp     *types.Response[types.ElasticIPResponse]
	BlockStorageResp  *types.Response[types.BlockStorageResponse]
	SnapshotResp      *types.Response[types.SnapshotResponse]
	BackupResp        *types.Response[types.StorageBackupResponse]
	RestoreResp       *types.Response[types.RestoreResponse]
	ContainerRegistry *types.Response[types.ContainerRegistryResponse]
	VPCResp           *types.Response[types.VPCResponse]
	SubnetResp        *types.Response[types.SubnetResponse]
	SecurityGroupResp *types.Response[types.SecurityGroupResponse]
	SecurityRuleResp  *types.Response[types.SecurityRuleResponse]
	KeyPairResp       *types.Response[types.KeyPairResponse]
	DBaaSResp         *types.Response[types.DBaaSResponse]
	KaaSResp          *types.Response[types.KaaSResponse]
	CloudServerResp   *types.Response[types.CloudServerResponse]
}

// createAllResources creates all resources in the correct order
func createAllResources(ctx context.Context, arubaClient aruba.Client) *ResourceCollection {
	resources := &ResourceCollection{}

	// 1. Create Project
	resources.ProjectID = createProject(ctx, arubaClient)

	// 2. Create Elastic IP
	resources.ElasticIPResp = createElasticIP(ctx, arubaClient, resources.ProjectID)

	// 3. Create Block Storage (SDK handles waiting for it to be ready)
	resources.BlockStorageResp = createBlockStorage(ctx, arubaClient, resources.ProjectID)

	// 4. Create Snapshot from Block Storage (SDK waits for BlockStorage to be ready)
	resources.SnapshotResp = createSnapshot(ctx, arubaClient, resources.ProjectID, resources.BlockStorageResp)

	// 5. Create VPC (SDK handles waiting for it to be active)
	resources.VPCResp = createVPC(ctx, arubaClient, resources.ProjectID)

	// 6. Create Subnet in VPC (SDK waits for VPC to be active)
	resources.SubnetResp = createSubnet(ctx, arubaClient, resources.ProjectID, resources.VPCResp)

	// 7. Create Security Group (SDK waits for VPC to be active)
	resources.SecurityGroupResp = createSecurityGroup(ctx, arubaClient, resources.ProjectID, resources.VPCResp)

	// 8. Create Security Group Rule (SDK waits for SecurityGroup to be active)
	resources.SecurityRuleResp = createSecurityGroupRule(ctx, arubaClient, resources.ProjectID, resources.VPCResp, resources.SecurityGroupResp)

	// 9. Create SSH Key Pair
	resources.KeyPairResp = createKeyPair(ctx, arubaClient, resources.ProjectID)

	// 10. Create DBaaS
	resources.DBaaSResp = createDBaaS(ctx, arubaClient, resources.ProjectID, resources.VPCResp, resources.SubnetResp, resources.SecurityGroupResp)

	// 11. Create KaaS
	resources.KaaSResp = createKaaS(ctx, arubaClient, resources.ProjectID, resources.VPCResp, resources.SubnetResp, resources.SecurityGroupResp)

	// 12. Create Cloud Server (commented out)
	resources.CloudServerResp = createCloudServer(ctx, arubaClient, resources)

	// 13. Create Container Registry
	resources.ContainerRegistry = createContainerRegistry(ctx, arubaClient, resources.ProjectID)

	// 14. Create Storage Backup
	resources.BackupResp = createStorageBackup(ctx, arubaClient, resources.ProjectID, stringValue(resources.BlockStorageResp.Data.Metadata.ID))

	// 15. Create Restore from Backup
	resources.RestoreResp = createRestore(ctx, arubaClient, resources.ProjectID, resources.BackupResp)

	return resources
}

func createContainerRegistry(ctx context.Context, arubaClient aruba.Client, projectID string) *types.Response[types.ContainerRegistryResponse] {
	fmt.Println("--- Container Registry ---")

	req := types.ContainerRegistryRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "test-registry",
			},
		},
		Properties: types.ContainerRegistryPropertiesRequest{
			VPC:             types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"},
			Subnet:          types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"},
			SecurityGroup:   types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"},
			PublicIp:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"},
			BlockStorage:    types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"},
			BillingPlan:     &types.BillingPeriodResource{BillingPeriod: "Hour"},
			AdminUser:       &types.UserCredential{Username: "admin"},
			ConcurrentUsers: types.StringPtr("100"),
		},
	}

	resp, err := arubaClient.FromContainer().ContainerRegistry().Create(ctx, projectID, req, nil)
	if err != nil {
		log.Printf("Error creating container registry: %v", err)
		os.Exit(1)
	} else if !resp.IsSuccess() {
		log.Printf("Failed to create container registry - Status: %d, Error: %s, Detail: %s",
			resp.StatusCode,
			stringValue(resp.Error.Title),
			stringValue(resp.Error.Detail))
		os.Exit(1)
	}
	if resp.Data != nil {
		fmt.Printf("✓ Created container registry: %s\n", *resp.Data.Metadata.Name)
	} else {
		fmt.Println("Warning: resp.Data is nil")
	}

	return resp
}

func createRestore(ctx context.Context, arubaClient aruba.Client, s string, storageBackupResponse *types.Response[types.StorageBackupResponse]) *types.Response[types.RestoreResponse] {
	fmt.Println("--- Restore ---")

	if storageBackupResponse == nil || storageBackupResponse.Data == nil || storageBackupResponse.Data.Metadata.ID == nil {
		fmt.Println("⚠ Skipping Restore creation - Backup not available")
		return nil
	}

	backupID := *storageBackupResponse.Data.Metadata.ID

	req := types.RestoreRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "test-restore",
			},
		},
		Properties: types.RestorePropertiesRequest{
			Target: types.ReferenceResource{URI: storageBackupResponse.Data.Properties.Origin.URI},
		},
	}

	resp, err := arubaClient.FromStorage().Restores().Create(ctx, s, backupID, req, nil)
	if err != nil {
		log.Printf("Error creating restore: %v", err)
		os.Exit(1)
	} else if !resp.IsSuccess() {
		log.Printf("Failed to create restore - Status: %d, Error: %s, Detail: %s",
			resp.StatusCode,
			stringValue(resp.Error.Title),
			stringValue(resp.Error.Detail))
		os.Exit(1)
	}
	if resp.Data != nil {
		fmt.Printf("✓ Created restore: %s\n", *resp.Data.Metadata.Name)
	} else {
		fmt.Println("Warning: resp.Data is nil")
	}

	return resp
}

func createStorageBackup(ctx context.Context, arubaClient aruba.Client, projectID string, blockStorageID string) *types.Response[types.StorageBackupResponse] {
	fmt.Println("--- Storage Backup ---")

	req := types.StorageBackupRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "test-backup",
			},
		},
		Properties: types.StorageBackupPropertiesRequest{
			StorageBackupType: types.StorageBackupTypeFull,
			Origin:            types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/" + blockStorageID},
			RetentionDays:     types.IntPtr(10),
			BillingPeriod:     types.StringPtr("Monthly"),
		},
	}

	resp, err := arubaClient.FromStorage().Backups().Create(ctx, projectID, req, nil)
	if err != nil {
		log.Printf("Error creating storage backup: %v", err)
		os.Exit(1)
	} else if !resp.IsSuccess() {
		log.Printf("Failed to create storage backup - Status: %d, Error: %s, Detail: %s",
			resp.StatusCode,
			stringValue(resp.Error.Title),
			stringValue(resp.Error.Detail))
		os.Exit(1)
	}
	if resp.Data != nil {
		fmt.Printf("✓ Created storage backup: %s\n", *resp.Data.Metadata.Name)
	} else {
		fmt.Println("Warning: resp.Data is nil")
	}

	return resp
}

// createProject creates and updates a project
func createProject(ctx context.Context, arubaClient aruba.Client) string {
	fmt.Println("--- Project Management ---")

	projectReq := types.ProjectRequest{
		Metadata: types.ResourceMetadataRequest{
			Name: "seca-sdk-example",
			Tags: []string{"production", "arubacloud-sdk"},
		},
		Properties: types.ProjectPropertiesRequest{
			Description: stringPtr("My production project"),
			Default:     false,
		},
	}

	createResp, err := arubaClient.FromProject().Create(ctx, projectReq, nil)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	} else if !createResp.IsSuccess() {
		log.Fatalf("Failed to create project, status code: %d and error title: %s", createResp.StatusCode, stringValue(createResp.Error.Title))
	}
	projectID := *createResp.Data.Metadata.ID
	fmt.Printf("✓ Created project with ID: %s\n", projectID)

	// Update the project
	updateResp, err := arubaClient.FromProject().Update(ctx, projectID, projectReq, nil)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		os.Exit(1)
	} else if !updateResp.IsSuccess() {
		log.Printf("Failed to update project, status code: %d and error title: %s", updateResp.StatusCode, stringValue(updateResp.Error.Title))
		os.Exit(1)
	}
	fmt.Printf("✓ Updated project: %s\n", *updateResp.Data.Metadata.Name)

	return projectID
}

// createElasticIP creates an Elastic IP
func createElasticIP(ctx context.Context, arubaClient aruba.Client, projectID string) *types.Response[types.ElasticIPResponse] {
	fmt.Println("--- Elastic IP ---")

	elasticIPReq := types.ElasticIPRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-elastic-ip",
				Tags: []string{"network", "public"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.ElasticIPPropertiesRequest{
			BillingPlan: types.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	elasticIPResp, err := arubaClient.FromNetwork().ElasticIPs().Create(ctx, projectID, elasticIPReq, nil)
	if err != nil {
		log.Printf("Error creating Elastic IP: %v", err)
		os.Exit(1)
	} else if !elasticIPResp.IsSuccess() {
		log.Printf("Failed to create Elastic IP - Status: %d, Error: %s, Detail: %s",
			elasticIPResp.StatusCode,
			stringValue(elasticIPResp.Error.Title),
			stringValue(elasticIPResp.Error.Detail))
		os.Exit(1)
	}
	fmt.Printf("✓ Created Elastic IP: %s (ObjectID: %s)\n", *elasticIPResp.Data.Metadata.Name, *elasticIPResp.Data.Metadata.ID)

	return elasticIPResp
}

// createBlockStorage creates a block storage volume
// The SDK automatically waits for it to become Active or NotUsed
func createBlockStorage(ctx context.Context, arubaClient aruba.Client, projectID string) *types.Response[types.BlockStorageResponse] {
	fmt.Println("--- Block Storage ---")

	blockStorageReq := types.BlockStorageRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-block-storage",
				Tags: []string{"storage", "data"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.BlockStoragePropertiesRequest{
			SizeGB:        20,
			Type:          types.BlockStorageTypeStandard,
			Zone:          stringPtr("ITBG-1"),
			BillingPeriod: "Hour",
			Bootable:      boolPtr(true),
			Image:         stringPtr("LU22-001"),
		},
	}

	blockStorageResp, err := arubaClient.FromStorage().Volumes().Create(ctx, projectID, blockStorageReq, nil)
	if err != nil {
		log.Printf("Error creating block storage: %v", err)
		os.Exit(1)
	} else if !blockStorageResp.IsSuccess() {
		log.Printf("Failed to create block storage - Status: %d, Error: %s, Detail: %s",
			blockStorageResp.StatusCode,
			stringValue(blockStorageResp.Error.Title),
			stringValue(blockStorageResp.Error.Detail))
		os.Exit(1)
	}
	if blockStorageResp.Data != nil {
		fmt.Printf("✓ Created block storage: %s (%d GB, %s)\n", *blockStorageResp.Data.Metadata.Name,
			blockStorageResp.Data.Properties.SizeGB,
			blockStorageResp.Data.Properties.Type)
	} else {
		fmt.Println("Warning: blockStorageResp.Data is nil")
	}

	return blockStorageResp
}

// createSnapshot creates a snapshot from block storage
func createSnapshot(ctx context.Context, arubaClient aruba.Client, projectID string, blockStorageResp *types.Response[types.BlockStorageResponse]) *types.Response[types.SnapshotResponse] {
	fmt.Println("--- Snapshot ---")

	snapshotReq := types.SnapshotRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-snapshot",
				Tags: []string{"backup", "snapshot"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.SnapshotPropertiesRequest{
			BillingPeriod: stringPtr("Hour"),
			Volume: types.ReferenceResource{
				URI: *blockStorageResp.Data.Metadata.URI,
			},
		},
	}

	snapshotResp, err := arubaClient.FromStorage().Snapshots().Create(ctx, projectID, snapshotReq, nil)
	if err != nil {
		log.Printf("Error creating snapshot: %v", err)
		os.Exit(1)
	} else if !snapshotResp.IsSuccess() {
		log.Printf("Failed to create snapshot - Status: %d, Error: %s, Detail: %s",
			snapshotResp.StatusCode,
			stringValue(snapshotResp.Error.Title),
			stringValue(snapshotResp.Error.Detail))
		os.Exit(1)
	}
	if snapshotResp.Data != nil {
		fmt.Printf("✓ Created snapshot: %s from volume %s\n",
			*snapshotResp.Data.Metadata.Name,
			*blockStorageResp.Data.Metadata.Name)
	} else {
		fmt.Println("Warning: snapshotResp.Data is nil")
	}

	return snapshotResp
}

// createVPC creates a VPC
// The SDK automatically waits for it to become Active for dependent operations
func createVPC(ctx context.Context, arubaClient aruba.Client, projectID string) *types.Response[types.VPCResponse] {
	fmt.Println("--- VPC ---")

	vpcReq := types.VPCRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-vpc",
				Tags: []string{"network", "infrastructure"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.VPCPropertiesRequest{
			Properties: &types.VPCProperties{
				Default: boolPtr(false),
				Preset:  boolPtr(false),
			},
		},
	}

	vpcResp, err := arubaClient.FromNetwork().VPCs().Create(ctx, projectID, vpcReq, nil)
	if err != nil {
		log.Printf("Error creating VPC: %v", err)
		os.Exit(1)
	} else if vpcResp.IsError() && vpcResp.Error != nil {
		log.Printf("Failed to create VPC - Status: %d, Error: %s, Detail: %s",
			vpcResp.StatusCode,
			stringValue(vpcResp.Error.Title),
			stringValue(vpcResp.Error.Detail))
		os.Exit(1)
	}

	if vpcResp.Data != nil && vpcResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created VPC: %s (Default: %t)\n",
			*vpcResp.Data.Metadata.Name,
			vpcResp.Data.Properties.Default)
	} else {
		fmt.Println("Warning: vpcResp.Data or vpcResp.Data.Metadata.Name is nil")
	}

	return vpcResp
}

// createSubnet creates a subnet in a VPC
func createSubnet(ctx context.Context, arubaClient aruba.Client, projectID string, vpcResp *types.Response[types.VPCResponse]) *types.Response[types.SubnetResponse] {
	fmt.Println("\n--- Network: Subnet ---")

	vpcID := *vpcResp.Data.Metadata.ID

	subnetReq := types.SubnetRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-subnet",
				Tags: []string{"network", "subnet"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.SubnetPropertiesRequest{
			Type:    types.SubnetTypeAdvanced,
			Default: false,
			Network: &types.SubnetNetwork{
				Address: "192.168.1.0/25",
			},
			DHCP: &types.SubnetDHCP{
				Enabled: true,
				Range: &types.SubnetDHCPRange{
					Start: "192.168.1.10",
					Count: 50,
				},
				Routes: []types.SubnetDHCPRoute{
					{
						Address: "0.0.0.0/0",
						Gateway: "192.168.1.1",
					},
				},
				DNS: []string{"8.8.8.8", "8.8.4.4"},
			},
		},
	}

	subnetResp, err := arubaClient.FromNetwork().Subnets().Create(ctx, projectID, vpcID, subnetReq, nil)
	if err != nil {
		log.Printf("Error creating subnet: %v", err)
	} else if subnetResp.IsError() && subnetResp.Error != nil {
		log.Printf("Failed to create subnet - Status: %d, Error: %s, Detail: %s",
			subnetResp.StatusCode,
			stringValue(subnetResp.Error.Title),
			stringValue(subnetResp.Error.Detail))
	} else if subnetResp.Data != nil && subnetResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created Subnet: %s (Type: %s, Network: %s)\n",
			*subnetResp.Data.Metadata.Name,
			subnetResp.Data.Properties.Type,
			subnetResp.Data.Properties.Network.Address)
	}

	return subnetResp
}

// createSecurityGroup creates a security group
// The SDK automatically waits for the VPC to become Active before creating the group
func createSecurityGroup(ctx context.Context, arubaClient aruba.Client, projectID string, vpcResp *types.Response[types.VPCResponse]) *types.Response[types.SecurityGroupResponse] {
	fmt.Println("\n--- Network: Security Group ---")

	vpcID := *vpcResp.Data.Metadata.ID

	sgReq := types.SecurityGroupRequest{
		Metadata: types.ResourceMetadataRequest{
			Name: "my-security-group",
			Tags: []string{"security", "network"},
		},
		Properties: types.SecurityGroupPropertiesRequest{
			Default: boolPtr(false),
		},
	}

	sgResp, err := arubaClient.FromNetwork().SecurityGroups().Create(ctx, projectID, vpcID, sgReq, nil)
	if err != nil {
		log.Printf("Error creating security group: %v", err)
		return nil
	} else if sgResp.IsError() && sgResp.Error != nil {
		log.Printf("Failed to create security group - Status: %d, Error: %s, Detail: %s",
			sgResp.StatusCode,
			stringValue(sgResp.Error.Title),
			stringValue(sgResp.Error.Detail))
		return nil
	}

	if sgResp.Data != nil && sgResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created Security Group: %s\n", *sgResp.Data.Metadata.Name)
	}

	return sgResp
}

// createSecurityGroupRule creates a security group rule
func createSecurityGroupRule(ctx context.Context, arubaClient aruba.Client, projectID string, vpcResp *types.Response[types.VPCResponse], sgResp *types.Response[types.SecurityGroupResponse]) *types.Response[types.SecurityRuleResponse] {
	if sgResp == nil || sgResp.Data == nil {
		fmt.Println("⚠ Skipping security rule creation - Security Group not available")
		return nil
	}

	fmt.Println("\n--- Network: Security Group Rule ---")

	vpcID := *vpcResp.Data.Metadata.ID
	sgID := *sgResp.Data.Metadata.ID

	ruleReq := types.SecurityRuleRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "allow-ssh",
				Tags: []string{"ssh-access", "ingress"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.SecurityRulePropertiesRequest{
			Direction: types.RuleDirectionIngress,
			Protocol:  "TCP",
			Port:      "22",
			Target: &types.RuleTarget{
				Kind:  types.EndpointTypeIP,
				Value: "0.0.0.0/0",
			},
		},
	}

	ruleResp, err := arubaClient.FromNetwork().SecurityGroupRules().Create(ctx, projectID, vpcID, sgID, ruleReq, nil)
	if err != nil {
		log.Printf("Error creating security rule: %v", err)
	} else if ruleResp.IsError() && ruleResp.Error != nil {
		log.Printf("Failed to create security rule - Status: %d, Error: %s, Detail: %s",
			ruleResp.StatusCode,
			stringValue(ruleResp.Error.Title),
			stringValue(ruleResp.Error.Detail))
	} else if ruleResp.Data != nil && ruleResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created Security Rule: %s (Direction: %s, Protocol: %s, Port: %s, Target: %s)\n",
			*ruleResp.Data.Metadata.Name,
			ruleResp.Data.Properties.Direction,
			ruleResp.Data.Properties.Protocol,
			ruleResp.Data.Properties.Port,
			ruleResp.Data.Properties.Target.Value)
	}

	return ruleResp
}

// createKeyPair creates an SSH key pair
func createKeyPair(ctx context.Context, arubaClient aruba.Client, projectID string) *types.Response[types.KeyPairResponse] {
	fmt.Println("--- SSH Key Pair ---")

	sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA2No7At0tgHrcZTL0kGWyLLUqPKfOhD9hGdNV9PbJxhjOGNFxcwdQ9wCXsJ3RQaRHBuGIgVodDurrlqzxFK86yCHMgXT2YLHF0j9P4m9GDiCfOK6msbFb89p5xZExjwD2zK+w68r7iOKZeRB2yrznW5TD3KDemSPIQQIVcyLF+yxft49HWBTI3PVQ4rBVOBJ2PdC9SAOf7CYnptW24CRrC0h85szIDwMA+Kmasfl3YGzk4MxheHrTO8C40aXXpieJ9S2VQA4VJAMRyAboptIK0cKjBYrbt5YkEL0AlyBGPIu6MPYr5K/MHyDunDi9yc7VYRYRR0f46MBOSqMUiGPnMw=="

	keyPairReq := types.KeyPairRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "allow-ssh",
				Tags: []string{"ssh-access", "ingress"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.KeyPairPropertiesRequest{
			Value: sshPublicKey,
		},
	}

	keyPairResp, err := arubaClient.FromCompute().KeyPairs().Create(ctx, projectID, keyPairReq, nil)
	if err != nil {
		log.Printf("Error creating SSH key pair: %v", err)
	} else if !keyPairResp.IsSuccess() {
		log.Printf("Failed to create SSH key pair - Status: %d, Error: %s, Detail: %s",
			keyPairResp.StatusCode,
			stringValue(keyPairResp.Error.Title),
			stringValue(keyPairResp.Error.Detail))
	} else if keyPairResp.Data != nil && *keyPairResp.Data.Metadata.Name != "" {
		fmt.Printf("✓ Created SSH Key Pair: %s\n", *keyPairResp.Data.Metadata.Name)
	}

	return keyPairResp
}

// createDBaaS creates a DBaaS instance
func createDBaaS(ctx context.Context, arubaClient aruba.Client, projectID string, vpcResp *types.Response[types.VPCResponse], subnetResp *types.Response[types.SubnetResponse], sgResp *types.Response[types.SecurityGroupResponse]) *types.Response[types.DBaaSResponse] {
	fmt.Println("--- DBaaS ---")

	// Only create DBaaS if VPC, Subnet, and Security Group are available
	if vpcResp == nil || vpcResp.Data == nil || vpcResp.Data.Metadata.URI == nil ||
		subnetResp == nil || subnetResp.Data == nil || subnetResp.Data.Metadata.URI == nil ||
		sgResp == nil || sgResp.Data == nil || sgResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping DBaaS creation - VPC, Subnet, or Security Group not available")
		return nil
	}

	dbaasReq := types.DBaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-dbaas",
				Tags: []string{"database", "mysql"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.DBaaSPropertiesRequest{
			Engine: &types.DBaaSEngine{
				ID:         stringPtr("mysql-8.0"),
				DataCenter: stringPtr("ITBG-1"),
			},
			Flavor: &types.DBaaSFlavor{
				Name: stringPtr("DBO2A4"),
			},
			Storage: &types.DBaaSStorage{
				SizeGB: int32Ptr(20),
			},
			BillingPlan: &types.DBaaSBillingPlan{
				BillingPeriod: stringPtr("Hour"),
			},
			Networking: &types.DBaaSNetworking{
				VPCURI:           vpcResp.Data.Metadata.URI,
				SubnetURI:        subnetResp.Data.Metadata.URI,
				SecurityGroupURI: sgResp.Data.Metadata.URI,
			},
			Autoscaling: &types.DBaaSAutoscaling{
				Enabled:        boolPtr(true),
				AvailableSpace: int32Ptr(20),
				StepSize:       int32Ptr(10),
			},
		},
	}

	dbaasResp, err := arubaClient.FromDatabase().DBaaS().Create(ctx, projectID, dbaasReq, nil)
	if err != nil {
		log.Printf("Error creating DBaaS: %v", err)
		return nil
	} else if !dbaasResp.IsSuccess() {
		log.Printf("Failed to create DBaaS - Status: %d, Error: %s, Detail: %s",
			dbaasResp.StatusCode,
			stringValue(dbaasResp.Error.Title),
			stringValue(dbaasResp.Error.Detail))
		return nil
	}

	if dbaasResp.Data != nil && dbaasResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created DBaaS: %s (Engine: %s, Flavor: %s, Storage: %d GB)\n",
			*dbaasResp.Data.Metadata.Name,
			stringValue(dbaasResp.Data.Properties.Engine.Type),
			stringValue(dbaasResp.Data.Properties.Flavor.Name),
			int32Value(dbaasResp.Data.Properties.Storage.SizeGB))
	}

	return dbaasResp
}

// createKaaS creates a KaaS cluster
func createKaaS(ctx context.Context, arubaClient aruba.Client, projectID string, vpcResp *types.Response[types.VPCResponse], subnetResp *types.Response[types.SubnetResponse], sgResp *types.Response[types.SecurityGroupResponse]) *types.Response[types.KaaSResponse] {
	fmt.Println("--- KaaS (Kubernetes) ---")

	// Only create KaaS if VPC, Subnet, and Security Group are available
	if vpcResp == nil || vpcResp.Data == nil || vpcResp.Data.Metadata.URI == nil ||
		subnetResp == nil || subnetResp.Data == nil || subnetResp.Data.Metadata.URI == nil ||
		sgResp == nil || sgResp.Data == nil || sgResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping KaaS creation - VPC, Subnet, or Security Group not available")
		return nil
	}

	// Wait for Subnet to become Active before creating KaaS
	vpcID := *vpcResp.Data.Metadata.ID
	subnetID := *subnetResp.Data.Metadata.ID
	fmt.Println("⏳ Waiting for Subnet to become Active before creating KaaS...")

	// Create a simple polling loop to check Subnet state
	maxAttempts := 30
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		subnetCheckResp, err := arubaClient.FromNetwork().Subnets().Get(ctx, projectID, vpcID, subnetID, nil)
		if err != nil {
			log.Printf("Error checking Subnet state: %v", err)
			break
		}

		if subnetCheckResp.Data != nil && subnetCheckResp.Data.Status.State != nil {
			state := *subnetCheckResp.Data.Status.State
			if state == "Active" {
				fmt.Println("✓ Subnet is now Active")
				break
			}
			if attempt < maxAttempts {
				time.Sleep(5 * time.Second)
			} else {
				fmt.Printf("⚠ Subnet did not become Active after %d attempts (state: %s)\n", maxAttempts, state)
			}
		}
	}

	kaasReq := types.KaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-kaas-cluster",
				Tags: []string{"kubernetes", "container"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.KaaSPropertiesRequest{
			Preset: boolPtr(false),
			VPC: types.ReferenceResource{
				URI: *vpcResp.Data.Metadata.URI,
			},
			Subnet: types.ReferenceResource{
				URI: *subnetResp.Data.Metadata.URI,
			},
			SecurityGroup: types.SecurityGroupProperties{
				Name: "sg-name-for-kaas",
			},
			NodeCIDR: types.NodeCIDRProperties{
				Name:    "my-node-cidr",
				Address: "10.100.0.0/16",
			},
			KubernetesVersion: types.KubernetesVersionInfo{
				Value: "1.32.3",
			},
			NodePools: []types.NodePoolProperties{
				{
					Name:     "default-pool",
					Nodes:    3,
					Instance: "K4A8",
					Zone:     "ITBG-1",
				},
			},
			HA: boolPtr(true),
			BillingPlan: types.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	kaasResp, err := arubaClient.FromContainer().KaaS().Create(ctx, projectID, kaasReq, nil)
	if err != nil {
		log.Printf("Error creating KaaS cluster: %v", err)
		return nil
	} else if !kaasResp.IsSuccess() {
		log.Printf("Failed to create KaaS cluster - Status: %d, Error: %s, Detail: %s",
			kaasResp.StatusCode,
			stringValue(kaasResp.Error.Title),
			stringValue(kaasResp.Error.Detail))
		return nil
	}

	if kaasResp.Data != nil && kaasResp.Data.Metadata.Name != nil {
		nodeCount := 0
		if kaasResp.Data.Properties.NodePools != nil {
			nodeCount = len(*kaasResp.Data.Properties.NodePools)
		}
		haValue := false
		if kaasResp.Data.Properties.HA != nil {
			haValue = *kaasResp.Data.Properties.HA
		}
		k8sVersion := ""
		if kaasResp.Data.Properties.KubernetesVersion.Value != nil {
			k8sVersion = *kaasResp.Data.Properties.KubernetesVersion.Value
		}
		fmt.Printf("✓ Created KaaS cluster: %s (K8s: %s, Nodes: %d, HA: %t)\n",
			*kaasResp.Data.Metadata.Name,
			k8sVersion,
			nodeCount,
			haValue)
	}

	return kaasResp
}

// createCloudServer creates a cloud server instance
func createCloudServer(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) *types.Response[types.CloudServerResponse] {
	fmt.Println("--- Cloud Server ---")

	// Verify all required resources are available
	if resources.VPCResp == nil || resources.VPCResp.Data == nil || resources.VPCResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping Cloud Server creation - VPC not available")
		return nil
	}
	if resources.ElasticIPResp == nil || resources.ElasticIPResp.Data == nil || resources.ElasticIPResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping Cloud Server creation - Elastic IP not available")
		return nil
	}
	if resources.BlockStorageResp == nil || resources.BlockStorageResp.Data == nil || resources.BlockStorageResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping Cloud Server creation - Block Storage not available")
		return nil
	}
	if resources.KeyPairResp == nil || resources.KeyPairResp.Data == nil || *resources.KeyPairResp.Data.Metadata.Name == "" {
		fmt.Println("⚠ Skipping Cloud Server creation - Key Pair not available")
		return nil
	}
	if resources.SubnetResp == nil || resources.SubnetResp.Data == nil || resources.SubnetResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping Cloud Server creation - Subnet not available")
		return nil
	}
	if resources.SecurityGroupResp == nil || resources.SecurityGroupResp.Data == nil || resources.SecurityGroupResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping Cloud Server creation - Security Group not available")
		return nil
	}

	// Example cloud-init content: update packages and create a welcome file
	cloudInitContent := `#cloud-config
package_update: true
package_upgrade: true
write_files:
  - path: /etc/motd
    content: |
      Welcome to Aruba Cloud Server!
      This server was initialized with cloud-init.
    owner: root:root
    permissions: '0644'
`
	// Base64 encode the cloud-init content
	userData := base64.StdEncoding.EncodeToString([]byte(cloudInitContent))

	cloudServerReq := types.CloudServerRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: types.ResourceMetadataRequest{
				Name: "my-cloudserver",
				Tags: []string{"virtualmachine", "container"},
			},
			Location: types.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.CloudServerPropertiesRequest{
			Zone:       "ITBG-1",
			VPCPreset:  false,
			FlavorName: stringPtr("CSO2A4"),
			VPC: types.ReferenceResource{
				URI: *resources.VPCResp.Data.Metadata.URI,
			},
			ElastcIP: types.ReferenceResource{
				URI: *resources.ElasticIPResp.Data.Metadata.URI,
			},
			BootVolume: types.ReferenceResource{
				URI: *resources.BlockStorageResp.Data.Metadata.URI,
			},
			KeyPair: types.ReferenceResource{
				URI: *resources.KeyPairResp.Data.Metadata.URI,
			},
			Subnets: []types.ReferenceResource{
				{
					URI: *resources.SubnetResp.Data.Metadata.URI,
				},
			},
			SecurityGroups: []types.ReferenceResource{
				{
					URI: *resources.SecurityGroupResp.Data.Metadata.URI,
				},
			},
			UserData: stringPtr(userData),
		},
	}

	cloudServerResp, err := arubaClient.FromCompute().CloudServers().Create(ctx, resources.ProjectID, cloudServerReq, nil)
	if err != nil {
		log.Printf("Error creating Cloud Server: %v", err)
		return nil
	} else if !cloudServerResp.IsSuccess() {
		log.Printf("Failed to create Cloud Server - Status: %d, Error: %s, Detail: %s",
			cloudServerResp.StatusCode,
			stringValue(cloudServerResp.Error.Title),
			stringValue(cloudServerResp.Error.Detail))
		return nil
	}

	if cloudServerResp.Data != nil && cloudServerResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created Cloud Server: %s (Zone: %s, Flavor: %s)\n",
			*cloudServerResp.Data.Metadata.Name,
			cloudServerResp.Data.Properties.Zone,
			cloudServerResp.Data.Properties.Flavor.Name)
	}

	return cloudServerResp
}

// printResourceSummary prints a summary of all created resources
func printResourceSummary(resources *ResourceCollection) {
	fmt.Println("\n=== SDK Example Complete ===")
	fmt.Println("Successfully created resources:")
	fmt.Println("- Project ID:", resources.ProjectID)

	if resources.ElasticIPResp != nil && resources.ElasticIPResp.Data != nil && resources.ElasticIPResp.Data.Metadata.ID != nil {
		fmt.Println("- ElasticIP ID:", *resources.ElasticIPResp.Data.Metadata.ID)
	}

	if resources.BlockStorageResp != nil && resources.BlockStorageResp.Data != nil && resources.BlockStorageResp.Data.Metadata.ID != nil {
		fmt.Println("- Block Storage ID:", *resources.BlockStorageResp.Data.Metadata.ID)
	}

	if resources.SnapshotResp != nil && resources.SnapshotResp.Data != nil && resources.SnapshotResp.Data.Metadata.ID != nil {
		fmt.Println("- Snapshot ID:", *resources.SnapshotResp.Data.Metadata.ID)
	}

	if resources.VPCResp != nil && resources.VPCResp.Data != nil && resources.VPCResp.Data.Metadata.ID != nil {
		fmt.Println("- VPC ID:", *resources.VPCResp.Data.Metadata.ID)
	}

	if resources.SubnetResp != nil && resources.SubnetResp.Data != nil && resources.SubnetResp.Data.Metadata.ID != nil {
		fmt.Println("- Subnet ID:", *resources.SubnetResp.Data.Metadata.ID)
	}

	if resources.SecurityGroupResp != nil && resources.SecurityGroupResp.Data != nil && resources.SecurityGroupResp.Data.Metadata.ID != nil {
		fmt.Println("- Security Group ID:", *resources.SecurityGroupResp.Data.Metadata.ID)
	}

	if resources.SecurityRuleResp != nil && resources.SecurityRuleResp.Data != nil && resources.SecurityRuleResp.Data.Metadata.ID != nil {
		fmt.Println("- Security Rule ID:", *resources.SecurityRuleResp.Data.Metadata.ID)
	}

	if resources.KeyPairResp != nil && resources.KeyPairResp.Data != nil && *resources.KeyPairResp.Data.Metadata.Name != "" {
		fmt.Println("- SSH Key Pair:", resources.KeyPairResp.Data.Metadata.Name)
	}

	if resources.DBaaSResp != nil && resources.DBaaSResp.Data != nil && resources.DBaaSResp.Data.Metadata.ID != nil {
		fmt.Println("- DBaaS ID:", *resources.DBaaSResp.Data.Metadata.ID)
	}

	if resources.KaaSResp != nil && resources.KaaSResp.Data != nil && resources.KaaSResp.Data.Metadata.ID != nil {
		fmt.Println("- KaaS Cluster ID:", *resources.KaaSResp.Data.Metadata.ID)
	}

	if resources.CloudServerResp != nil && resources.CloudServerResp.Data != nil && resources.CloudServerResp.Data.Metadata.Name != nil {
		fmt.Println("- Cloud Server:", *resources.CloudServerResp.Data.Metadata.Name)
	}
}

// Helper for pointer types
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func int32Value(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}
