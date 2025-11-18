package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sdkgo "github.com/Arubacloud/sdk-go"
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
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
func main() {
	mode := flag.String("mode", "create", "Operation mode: create, update, or delete")
	flag.Parse()

	switch *mode {
	case "create":
		runCreateExample()
	case "update":
		runUpdateExample()
	case "delete":
		runDeleteExample()
	default:
		log.Fatalf("Unknown mode: %s. Use 'create', 'update', or 'delete'", *mode)
	}
}

func runCreateExample() {
	config := &restclient.Config{
		ClientID:     "clientId",
		ClientSecret: "clientSecret",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		Debug:        true,
	}

	// Initialize the SDK (automatically obtains JWT token)
	sdk, err := sdkgo.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with timeout - increased to handle multiple resource creations and polling
	// With multiple resources polling (30 attempts × 5s each = 150s per resource), we need sufficient time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Use the SDK with context
	sdk.Client = sdk.Client.WithContext(ctx)

	fmt.Println("\n=== SDK Create Example ===")

	// Create all resources
	resources := createAllResources(ctx, sdk)

	// Print summary
	printResourceSummary(resources)
}

// ResourceCollection holds all created resources
type ResourceCollection struct {
	ProjectID         string
	ElasticIPResp     *schema.Response[schema.ElasticIPResponse]
	BlockStorageResp  *schema.Response[schema.BlockStorageResponse]
	SnapshotResp      *schema.Response[schema.SnapshotResponse]
	VPCResp           *schema.Response[schema.VPCResponse]
	SubnetResp        *schema.Response[schema.SubnetResponse]
	SecurityGroupResp *schema.Response[schema.SecurityGroupResponse]
	SecurityRuleResp  *schema.Response[schema.SecurityRuleResponse]
	KeyPairResp       *schema.Response[schema.KeyPairResponse]
	DBaaSResp         *schema.Response[schema.DBaaSResponse]
	KaaSResp          *schema.Response[schema.KaaSResponse]
	CloudServerResp   *schema.Response[schema.CloudServerResponse]
}

// createAllResources creates all resources in the correct order
func createAllResources(ctx context.Context, sdk *sdkgo.Client) *ResourceCollection {
	resources := &ResourceCollection{}

	// 1. Create Project
	resources.ProjectID = createProject(ctx, sdk)

	// 2. Create Elastic IP
	resources.ElasticIPResp = createElasticIP(ctx, sdk, resources.ProjectID)

	// 3. Create Block Storage (SDK handles waiting for it to be ready)
	resources.BlockStorageResp = createBlockStorage(ctx, sdk, resources.ProjectID)

	// 4. Create Snapshot from Block Storage (SDK waits for BlockStorage to be ready)
	resources.SnapshotResp = createSnapshot(ctx, sdk, resources.ProjectID, resources.BlockStorageResp)

	// 5. Create VPC (SDK handles waiting for it to be active)
	resources.VPCResp = createVPC(ctx, sdk, resources.ProjectID)

	// 6. Create Subnet in VPC (SDK waits for VPC to be active)
	resources.SubnetResp = createSubnet(ctx, sdk, resources.ProjectID, resources.VPCResp)

	// 7. Create Security Group (SDK waits for VPC to be active)
	resources.SecurityGroupResp = createSecurityGroup(ctx, sdk, resources.ProjectID, resources.VPCResp)

	// 8. Create Security Group Rule (SDK waits for SecurityGroup to be active)
	resources.SecurityRuleResp = createSecurityGroupRule(ctx, sdk, resources.ProjectID, resources.VPCResp, resources.SecurityGroupResp)

	// 9. Create SSH Key Pair
	resources.KeyPairResp = createKeyPair(ctx, sdk, resources.ProjectID)

	// 10. Create DBaaS
	resources.DBaaSResp = createDBaaS(ctx, sdk, resources.ProjectID, resources.VPCResp, resources.SubnetResp, resources.SecurityGroupResp)

	// 11. Create KaaS
	resources.KaaSResp = createKaaS(ctx, sdk, resources.ProjectID, resources.VPCResp, resources.SubnetResp, resources.SecurityGroupResp)

	// 12. Create Cloud Server (commented out)
	resources.CloudServerResp = createCloudServer(ctx, sdk, resources)

	return resources
}

// createProject creates and updates a project
func createProject(ctx context.Context, sdk *sdkgo.Client) string {
	fmt.Println("--- Project Management ---")

	projectReq := schema.ProjectRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "seca-sdk-example",
			Tags: []string{"production", "arubacloud-sdk"},
		},
		Properties: schema.ProjectPropertiesRequest{
			Description: stringPtr("My production project"),
			Default:     false,
		},
	}

	createResp, err := sdk.Project.CreateProject(ctx, projectReq, nil)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	} else if !createResp.IsSuccess() {
		log.Fatalf("Failed to create project, status code: %d and error title: %s", createResp.StatusCode, stringValue(createResp.Error.Title))
	}
	projectID := *createResp.Data.Metadata.ID
	fmt.Printf("✓ Created project with ID: %s\n", projectID)

	// Update the project
	updateResp, err := sdk.Project.UpdateProject(ctx, projectID, projectReq, nil)
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
func createElasticIP(ctx context.Context, sdk *sdkgo.Client, projectID string) *schema.Response[schema.ElasticIPResponse] {
	fmt.Println("--- Elastic IP ---")

	elasticIPReq := schema.ElasticIPRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-elastic-ip",
				Tags: []string{"network", "public"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.ElasticIPPropertiesRequest{
			BillingPlan: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	elasticIPResp, err := sdk.Network.CreateElasticIP(ctx, projectID, elasticIPReq, nil)
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
	fmt.Printf("✓ Created Elastic IP: %s (ObjectId: %s)\n", *elasticIPResp.Data.Metadata.Name, *elasticIPResp.Data.Metadata.ID)

	return elasticIPResp
}

// createBlockStorage creates a block storage volume
// The SDK automatically waits for it to become Active or NotUsed
func createBlockStorage(ctx context.Context, sdk *sdkgo.Client, projectID string) *schema.Response[schema.BlockStorageResponse] {
	fmt.Println("--- Block Storage ---")

	blockStorageReq := schema.BlockStorageRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-block-storage",
				Tags: []string{"storage", "data"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.BlockStoragePropertiesRequest{
			SizeGB:        20,
			Type:          schema.BlockStorageTypeStandard,
			Zone:          "ITBG-1",
			BillingPeriod: "Hour",
			Bootable:      boolPtr(true),
			Image:         stringPtr("LU22-001"),
		},
	}

	blockStorageResp, err := sdk.Storage.CreateBlockStorageVolume(ctx, projectID, blockStorageReq, nil)
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
func createSnapshot(ctx context.Context, sdk *sdkgo.Client, projectID string, blockStorageResp *schema.Response[schema.BlockStorageResponse]) *schema.Response[schema.SnapshotResponse] {
	fmt.Println("--- Snapshot ---")

	snapshotReq := schema.SnapshotRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-snapshot",
				Tags: []string{"backup", "snapshot"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.SnapshotPropertiesRequest{
			BillingPeriod: stringPtr("Hour"),
			Volume: schema.ReferenceResource{
				URI: *blockStorageResp.Data.Metadata.URI,
			},
		},
	}

	snapshotResp, err := sdk.Storage.CreateSnapshot(ctx, projectID, snapshotReq, nil)
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
func createVPC(ctx context.Context, sdk *sdkgo.Client, projectID string) *schema.Response[schema.VPCResponse] {
	fmt.Println("--- VPC ---")

	vpcReq := schema.VPCRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-vpc",
				Tags: []string{"network", "infrastructure"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.VPCPropertiesRequest{
			Properties: &schema.VPCProperties{
				Default: boolPtr(false),
				Preset:  boolPtr(false),
			},
		},
	}

	vpcResp, err := sdk.Network.CreateVPC(ctx, projectID, vpcReq, nil)
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
func createSubnet(ctx context.Context, sdk *sdkgo.Client, projectID string, vpcResp *schema.Response[schema.VPCResponse]) *schema.Response[schema.SubnetResponse] {
	fmt.Println("\n--- Network: Subnet ---")

	vpcID := *vpcResp.Data.Metadata.ID

	subnetReq := schema.SubnetRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-subnet",
				Tags: []string{"network", "subnet"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.SubnetPropertiesRequest{
			Type:    schema.SubnetTypeAdvanced,
			Default: false,
			Network: &schema.SubnetNetwork{
				Address: "192.168.1.0/25",
			},
			DHCP: &schema.SubnetDHCP{
				Enabled: true,
			},
		},
	}

	subnetResp, err := sdk.Network.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)
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
func createSecurityGroup(ctx context.Context, sdk *sdkgo.Client, projectID string, vpcResp *schema.Response[schema.VPCResponse]) *schema.Response[schema.SecurityGroupResponse] {
	fmt.Println("\n--- Network: Security Group ---")

	vpcID := *vpcResp.Data.Metadata.ID

	sgReq := schema.SecurityGroupRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "my-security-group",
			Tags: []string{"security", "network"},
		},
		Properties: schema.SecurityGroupPropertiesRequest{
			Default: boolPtr(false),
		},
	}

	sgResp, err := sdk.Network.CreateSecurityGroup(ctx, projectID, vpcID, sgReq, nil)
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
func createSecurityGroupRule(ctx context.Context, sdk *sdkgo.Client, projectID string, vpcResp *schema.Response[schema.VPCResponse], sgResp *schema.Response[schema.SecurityGroupResponse]) *schema.Response[schema.SecurityRuleResponse] {
	if sgResp == nil || sgResp.Data == nil {
		fmt.Println("⚠ Skipping security rule creation - Security Group not available")
		return nil
	}

	fmt.Println("\n--- Network: Security Group Rule ---")

	vpcID := *vpcResp.Data.Metadata.ID
	sgID := *sgResp.Data.Metadata.ID

	ruleReq := schema.SecurityRuleRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "allow-ssh",
				Tags: []string{"ssh-access", "ingress"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.SecurityRulePropertiesRequest{
			Direction: schema.RuleDirectionIngress,
			Protocol:  "TCP",
			Port:      "22",
			Target: &schema.RuleTarget{
				Kind:  schema.EndpointTypeIP,
				Value: "0.0.0.0/0",
			},
		},
	}

	ruleResp, err := sdk.Network.CreateSecurityGroupRule(ctx, projectID, vpcID, sgID, ruleReq, nil)
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
func createKeyPair(ctx context.Context, sdk *sdkgo.Client, projectID string) *schema.Response[schema.KeyPairResponse] {
	fmt.Println("--- SSH Key Pair ---")

	sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA2No7At0tgHrcZTL0kGWyLLUqPKfOhD9hGdNV9PbJxhjOGNFxcwdQ9wCXsJ3RQaRHBuGIgVodDurrlqzxFK86yCHMgXT2YLHF0j9P4m9GDiCfOK6msbFb89p5xZExjwD2zK+w68r7iOKZeRB2yrznW5TD3KDemSPIQQIVcyLF+yxft49HWBTI3PVQ4rBVOBJ2PdC9SAOf7CYnptW24CRrC0h85szIdwMA+Kmasfl3YGzk4MxheHrTO8C40aXXpieJ9S2VQA4VJAMRyAboptIK0cKjBYrbt5YkEL0AlyBGPIu6MPYr5K/MHyDunDi9yc7VYRYRR0f46MBOSqMUiGPnMw=="

	keyPairReq := schema.KeyPairRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "allow-ssh",
				Tags: []string{"ssh-access", "ingress"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.KeyPairPropertiesRequest{
			Value: sshPublicKey,
		},
	}

	keyPairResp, err := sdk.Compute.CreateKeyPair(ctx, projectID, keyPairReq, nil)
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
func createDBaaS(ctx context.Context, sdk *sdkgo.Client, projectID string, vpcResp *schema.Response[schema.VPCResponse], subnetResp *schema.Response[schema.SubnetResponse], sgResp *schema.Response[schema.SecurityGroupResponse]) *schema.Response[schema.DBaaSResponse] {
	fmt.Println("--- DBaaS ---")

	// Only create DBaaS if VPC, Subnet, and Security Group are available
	if vpcResp == nil || vpcResp.Data == nil || vpcResp.Data.Metadata.URI == nil ||
		subnetResp == nil || subnetResp.Data == nil || subnetResp.Data.Metadata.URI == nil ||
		sgResp == nil || sgResp.Data == nil || sgResp.Data.Metadata.URI == nil {
		fmt.Println("⚠ Skipping DBaaS creation - VPC, Subnet, or Security Group not available")
		return nil
	}

	dbaasReq := schema.DBaaSRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-dbaas",
				Tags: []string{"database", "mysql"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.DBaaSPropertiesRequest{
			Engine: &schema.DBaaSEngine{
				ID:         stringPtr("mysql-8.0"),
				DataCenter: stringPtr("ITBG-1"),
			},
			Flavor: &schema.DBaaSFlavor{
				Name: stringPtr("DBO2A4"),
			},
			Storage: &schema.DBaaSStorage{
				SizeGB: int32Ptr(20),
			},
			BillingPlan: &schema.DBaaSBillingPlan{
				BillingPeriod: stringPtr("Hour"),
			},
			Networking: &schema.DBaaSNetworking{
				VPCURI:           vpcResp.Data.Metadata.URI,
				SubnetURI:        subnetResp.Data.Metadata.URI,
				SecurityGroupURI: sgResp.Data.Metadata.URI,
			},
			Autoscaling: &schema.DBaaSAutoscaling{
				Enabled:        boolPtr(true),
				AvailableSpace: int32Ptr(20),
				StepSize:       int32Ptr(10),
			},
		},
	}

	dbaasResp, err := sdk.Database.CreateDBaaS(ctx, projectID, dbaasReq, nil)
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
func createKaaS(ctx context.Context, sdk *sdkgo.Client, projectID string, vpcResp *schema.Response[schema.VPCResponse], subnetResp *schema.Response[schema.SubnetResponse], sgResp *schema.Response[schema.SecurityGroupResponse]) *schema.Response[schema.KaaSResponse] {
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
		subnetCheckResp, err := sdk.Network.GetSubnet(ctx, projectID, vpcID, subnetID, nil)
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

	kaasReq := schema.KaaSRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-kaas-cluster",
				Tags: []string{"kubernetes", "container"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.KaaSPropertiesRequest{
			Preset: false,
			VPC: schema.ReferenceResource{
				URI: *vpcResp.Data.Metadata.URI,
			},
			Subnet: schema.ReferenceResource{
				URI: *subnetResp.Data.Metadata.URI,
			},
			SecurityGroup: schema.SecurityGroupProperties{
				Name: "sg-name-for-kaas",
			},
			NodeCIDR: schema.NodeCIDRProperties{
				Name:    "my-node-cidr",
				Address: "10.100.0.0/16",
			},
			KubernetesVersion: schema.KubernetesVersionInfo{
				Value: "1.32.3",
			},
			NodePools: []schema.NodePoolProperties{
				{
					Name:     "default-pool",
					Nodes:    3,
					Instance: "K4A8",
					Zone:     "ITBG-1",
				},
			},
			HA: true,
			BillingPlan: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	kaasResp, err := sdk.Container.CreateKaaS(ctx, projectID, kaasReq, nil)
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
		fmt.Printf("✓ Created KaaS cluster: %s (K8s: %s, Nodes: %d, HA: %t)\n",
			*kaasResp.Data.Metadata.Name,
			kaasResp.Data.Properties.KubernetesVersion.Value,
			len(kaasResp.Data.Properties.NodePools),
			kaasResp.Data.Properties.HA)
	}

	return kaasResp
}

// createCloudServer creates a cloud server instance
func createCloudServer(ctx context.Context, sdk *sdkgo.Client, resources *ResourceCollection) *schema.Response[schema.CloudServerResponse] {
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

	cloudServerReq := schema.CloudServerRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-cloudserver",
				Tags: []string{"virtualmachine", "container"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.CloudServerPropertiesRequest{
			Zone:       "ITBG-1",
			VPCPreset:  false,
			FlavorName: stringPtr("CSO2A4"),
			VPC: schema.ReferenceResource{
				URI: *resources.VPCResp.Data.Metadata.URI,
			},
			ElastcIP: schema.ReferenceResource{
				URI: *resources.ElasticIPResp.Data.Metadata.URI,
			},
			BootVolume: schema.ReferenceResource{
				URI: *resources.BlockStorageResp.Data.Metadata.URI,
			},
			KeyPair: schema.ReferenceResource{
				URI: *resources.KeyPairResp.Data.Metadata.URI,
			},
			Subnets: []schema.ReferenceResource{
				{
					URI: *resources.SubnetResp.Data.Metadata.URI,
				},
			},
			SecurityGroups: []schema.ReferenceResource{
				{
					URI: *resources.SecurityGroupResp.Data.Metadata.URI,
				},
			},
		},
	}

	cloudServerResp, err := sdk.Compute.CreateCloudServer(ctx, resources.ProjectID, cloudServerReq, nil)
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

	if cloudServerResp.Data != nil && cloudServerResp.Data.Metadata.Name != "" {
		fmt.Printf("✓ Created Cloud Server: %s (Zone: %s, Flavor: %s)\n",
			cloudServerResp.Data.Metadata.Name,
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

	if resources.CloudServerResp != nil && resources.CloudServerResp.Data != nil && resources.CloudServerResp.Data.Metadata.Name != "" {
		fmt.Println("- Cloud Server:", resources.CloudServerResp.Data.Metadata.Name)
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
