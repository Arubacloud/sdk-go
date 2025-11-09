package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/compute"
	"github.com/Arubacloud/sdk-go/pkg/spec/network"
	"github.com/Arubacloud/sdk-go/pkg/spec/project"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
	"github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

func main() {
	config := &client.Config{
		ClientID:     "clientId",
		ClientSecret: "clientSecret",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		Debug:        true,
	}

	// Initialize the SDK (automatically obtains JWT token)
	sdk, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with timeout - increased to handle multiple resource creations and polling
	// With 3 resources polling (30 attempts × 5s each = 150s per resource), we need sufficient time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Use the SDK with context
	sdk = sdk.WithContext(ctx)

	fmt.Println("\n=== SDK Examples ===")

	// Initialize service clients
	projectAPI := project.NewProjectService(sdk)
	elasticIPAPI := network.NewElasticIPService(sdk)
	vpcAPI := network.NewVPCService(sdk)
	storageAPI := storage.NewBlockStorageService(sdk)

	// Configuration for resource polling
	maxAttempts := 30
	pollInterval := 5 * time.Second

	// Example: Project Management
	fmt.Println("--- Project Management ---")

	// Create a new project
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

	createResp, err := projectAPI.CreateProject(ctx, projectReq, nil)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
		os.Exit(1)
	} else if !createResp.IsSuccess() {
		log.Fatalf("Failed to create project, status code: %d and error title: %s", createResp.StatusCode, stringValue(createResp.Error.Title))
		os.Exit(1)
	}
	projectID := *createResp.Data.Metadata.Id
	fmt.Printf("✓ Created project with ID: %s\n", projectID)

	// Update the project
	updateResp, err := projectAPI.UpdateProject(ctx, projectID, projectReq, nil)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		os.Exit(1)
	} else if !updateResp.IsSuccess() {
		log.Printf("Failed to update project, status code: %d and error title: %s", updateResp.StatusCode, stringValue(updateResp.Error.Title))
		os.Exit(1)
	}
	fmt.Printf("✓ Updated project: %s\n", *updateResp.Data.Metadata.Name)

	// Example: Create Elastic IP
	fmt.Println("\n--- Network: Elastic IP ---")

	elasticIPReq := schema.ElasticIpRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-elastic-ip",
				Tags: []string{"network", "public"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.ElasticIpPropertiesRequest{
			BillingPlan: schema.BillingPeriodResource{
				BillingPeriod: "Hour",
			},
		},
	}

	elasticIPResp, err := elasticIPAPI.CreateElasticIP(ctx, projectID, elasticIPReq, nil)
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
	fmt.Printf("✓ Created Elastic IP: %s (ObjectId: %s)\n", *elasticIPResp.Data.Metadata.Name, *elasticIPResp.Data.Metadata.Id)

	// Example: Create Block Storage
	fmt.Println("\n--- Storage: Block Storage ---")

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
			Image:         stringPtr("LU24-001"),
		},
	}

	blockStorageResp, err := storageAPI.CreateBlockStorageVolume(ctx, projectID, blockStorageReq, nil)
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

	// Wait for Block Storage to become active
	fmt.Println("\n⏳ Waiting for Block Storage to become active...")
	blockStorageID := *blockStorageResp.Data.Metadata.Id
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(pollInterval)

		getBlockStorageResp, err := storageAPI.GetBlockStorageVolume(ctx, projectID, blockStorageID, nil)
		if err != nil {
			log.Printf("Error checking Block Storage status: %v", err)
			continue
		}

		if getBlockStorageResp.Data != nil && getBlockStorageResp.Data.Status.State != nil {
			state := *getBlockStorageResp.Data.Status.State
			fmt.Printf("  Block Storage state: %s (attempt %d/%d)\n", state, attempt, maxAttempts)

			// Block storage can be "Active" (attached) or "NotUsed" (unattached but ready)
			if state == "Active" || state == "NotUsed" {
				fmt.Printf("✓ Block Storage is now ready (state: %s)\n", state)
				break
			} else if state == "Failed" || state == "Error" {
				log.Fatalf("Block Storage creation failed with state: %s", state)
			}
		}

		if attempt == maxAttempts {
			log.Fatalf("Timeout waiting for Block Storage to become ready")
		}
	}

	// Example: Create Snapshot from Block Storage
	fmt.Println("\n--- Storage: Snapshot ---")

	snapshotAPI := storage.NewSnapshotService(sdk)

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
				Uri: *blockStorageResp.Data.Metadata.Uri,
			},
		},
	}

	snapshotResp, err := snapshotAPI.CreateSnapshot(ctx, projectID, snapshotReq, nil)
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

	// Example: Create VPC
	fmt.Println("\n--- Network: VPC ---")

	vpcReq := schema.VpcRequest{
		Metadata: schema.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: schema.ResourceMetadataRequest{
				Name: "my-vpc",
				Tags: []string{"network", "infrastructure"},
			},
			Location: schema.LocationRequest{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: schema.VpcPropertiesRequest{
			Properties: &schema.VpcProperties{
				Default: boolPtr(false),
				Preset:  boolPtr(true),
			},
		},
	}

	vpcResp, err := vpcAPI.CreateVPC(ctx, projectID, vpcReq, nil)
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

	// Wait for VPC to become active
	fmt.Println("\n⏳ Waiting for VPC to become active...")
	vpcID := *vpcResp.Data.Metadata.Id
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(pollInterval)

		getVPCResp, err := vpcAPI.GetVPC(ctx, projectID, vpcID, nil)
		if err != nil {
			log.Printf("Error checking VPC status: %v", err)
			continue
		}

		if getVPCResp.Data != nil && getVPCResp.Data.Status.State != nil {
			state := *getVPCResp.Data.Status.State
			fmt.Printf("  VPC state: %s (attempt %d/%d)\n", state, attempt, maxAttempts)

			if state == "Active" {
				fmt.Println("✓ VPC is now active")
				break
			} else if state == "Failed" || state == "Error" {
				log.Fatalf("VPC creation failed with state: %s", state)
			}
		}

		if attempt == maxAttempts {
			log.Fatalf("Timeout waiting for VPC to become active")
		}
	}

	// Example: Create Subnet in VPC
	fmt.Println("\n--- Network: Subnet ---")

	subnetAPI := network.NewSubnetService(sdk)

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
			Default: true,
			Network: &schema.SubnetNetwork{
				Address: "192.168.1.0/25",
			},
			DHCP: &schema.SubnetDHCP{
				Enabled: true,
			},
		},
	}

	subnetResp, err := subnetAPI.CreateSubnet(ctx, projectID, vpcID, subnetReq, nil)
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

	// Example: Create Security Group
	fmt.Println("\n--- Network: Security Group ---")

	securityGroupAPI := network.NewSecurityGroupService(sdk)

	sgReq := schema.SecurityGroupRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "my-security-group",
			Tags: []string{"security", "network"},
		},
		Properties: schema.SecurityGroupPropertiesRequest{
			Default: boolPtr(false),
		},
	}

	var ruleResp *schema.Response[schema.SecurityRuleResponse]
	sgResp, err := securityGroupAPI.CreateSecurityGroup(ctx, projectID, vpcID, sgReq, nil)
	if err != nil {
		log.Printf("Error creating security group: %v", err)
	} else if sgResp.IsError() && sgResp.Error != nil {
		log.Printf("Failed to create security group - Status: %d, Error: %s, Detail: %s",
			sgResp.StatusCode,
			stringValue(sgResp.Error.Title),
			stringValue(sgResp.Error.Detail))
	} else if sgResp.Data != nil && sgResp.Data.Metadata.Name != nil {
		fmt.Printf("✓ Created Security Group: %s\n", *sgResp.Data.Metadata.Name)

		// Wait for Security Group to become active
		fmt.Println("\n⏳ Waiting for Security Group to become active...")
		sgID := *sgResp.Data.Metadata.Id
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			time.Sleep(pollInterval)

			getSGResp, err := securityGroupAPI.GetSecurityGroup(ctx, projectID, vpcID, sgID, nil)
			if err != nil {
				log.Printf("Error checking Security Group status: %v", err)
				continue
			}

			if getSGResp.Data != nil && getSGResp.Data.Status.State != nil {
				state := *getSGResp.Data.Status.State
				fmt.Printf("  Security Group state: %s (attempt %d/%d)\n", state, attempt, maxAttempts)

				if state == "Active" {
					fmt.Println("✓ Security Group is now active")
					break
				} else if state == "Failed" || state == "Error" {
					log.Fatalf("Security Group creation failed with state: %s", state)
				}
			}

			if attempt == maxAttempts {
				log.Fatalf("Timeout waiting for Security Group to become active")
			}
		}

		// Create security group rule to allow SSH from anywhere
		securityRuleAPI := network.NewSecurityGroupRuleService(sdk)

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

		ruleResp, err = securityRuleAPI.CreateSecurityGroupRule(ctx, projectID, vpcID, sgID, ruleReq, nil)
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
	}

	// Example: Create SSH Key Pair
	fmt.Println("\n--- Compute: SSH Key Pair ---")

	keyPairAPI := compute.NewKeyPairService(sdk)

	sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA2No7At0tgHrcZTL0kGWyLLUqPKfOhD9hGdNV9PbJxhjOGNFxcwdQ9wCXsJ3RQaRHBuGIgVodDurrlqzxFK86yCHMgXT2YLHF0j9P4m9GDiCfOK6msbFb89p5xZExjwD2zK+w68r7iOKZeRB2yrznW5TD3KDemSPIQQIVcyLF+yxft49HWBTI3PVQ4rBVOBJ2PdC9SAOf7CYnptW24CRrC0h85szIdwMA+Kmasfl3YGzk4MxheHrTO8C40aXXpieJ9S2VQA4VJAMRyAboptIK0cKjBYrbt5YkEL0AlyBGPIu6MPYr5K/MHyDunDi9yc7VYRYRR0f46MBOSqMUiGPnMw=="

	keyPairReq := schema.KeyPairRequest{
		Metadata: schema.ResourceMetadataRequest{
			Name: "my-ssh-keypair",
			Tags: []string{"compute", "access"},
		},
		Properties: schema.KeyPairPropertiesRequest{
			Value: sshPublicKey,
		},
	}

	keyPairResp, err := keyPairAPI.CreateKeyPair(ctx, projectID, keyPairReq, nil)
	if err != nil {
		log.Printf("Error creating SSH key pair: %v", err)
	} else if !keyPairResp.IsSuccess() {
		log.Printf("Failed to create SSH key pair - Status: %d, Error: %s, Detail: %s",
			keyPairResp.StatusCode,
			stringValue(keyPairResp.Error.Title),
			stringValue(keyPairResp.Error.Detail))
	} else if keyPairResp.Data != nil && keyPairResp.Data.Metadata.Name != "" {
		fmt.Printf("✓ Created SSH Key Pair: %s\n", keyPairResp.Data.Metadata.Name)
	}

	// Example: Create Cloud Server
	fmt.Println("\n--- Compute: Cloud Server ---")

	cloudServerAPI := compute.NewCloudServerService(sdk)

	var cloudServerResp *schema.Response[schema.CloudServerResponse]

	// Only create cloud server if all required resources are available
	if keyPairResp != nil && keyPairResp.Data != nil && keyPairResp.Data.Metadata.Name != "" {
		// Construct KeyPair URI from the project and keypair name
		keyPairUri := "/projects/" + projectID + "/providers/Aruba.Compute/keypairs/" + keyPairResp.Data.Metadata.Name

		cloudServerReq := schema.CloudServerRequest{
			Metadata: schema.ResourceMetadataRequest{
				Name: "my-cloud-server",
				Tags: []string{"compute", "production"},
			},
			Properties: schema.CloudServerPropertiesRequest{
				Zone: "ITBG-1",
				Vpc: schema.ReferenceResource{
					Uri: *vpcResp.Data.Metadata.Uri,
				},
				VpcPreset:  true,
				FlavorName: stringPtr("CSO4A8"),
				ElastcIp: schema.ReferenceResource{
					Uri: *elasticIPResp.Data.Metadata.Uri,
				},
				BootVolume: schema.ReferenceResource{
					Uri: *blockStorageResp.Data.Metadata.Uri,
				},
				KeyPair: schema.ReferenceResource{
					Uri: keyPairUri,
				},
			},
		}

		// Add subnet if it was created successfully
		if subnetResp != nil && subnetResp.Data != nil && subnetResp.Data.Metadata.Uri != nil {
			cloudServerReq.Properties.Subnets = []schema.ReferenceResource{
				{Uri: *subnetResp.Data.Metadata.Uri},
			}
		}

		// Add security group if it was created successfully
		if sgResp != nil && sgResp.Data != nil && sgResp.Data.Metadata.Uri != nil {
			cloudServerReq.Properties.SecurityGroups = []schema.ReferenceResource{
				{Uri: *sgResp.Data.Metadata.Uri},
			}
		}

		cloudServerResp, err = cloudServerAPI.CreateCloudServer(ctx, projectID, cloudServerReq, nil)
		if err != nil {
			log.Printf("Error creating cloud server: %v", err)
		} else if !cloudServerResp.IsSuccess() {
			log.Printf("Failed to create cloud server - Status: %d, Error: %s, Detail: %s",
				cloudServerResp.StatusCode,
				stringValue(cloudServerResp.Error.Title),
				stringValue(cloudServerResp.Error.Detail))
		} else if cloudServerResp.Data != nil && cloudServerResp.Data.Metadata.Name != "" {
			fmt.Printf("✓ Created Cloud Server: %s (Flavor: %s, Zone: %s)\n",
				cloudServerResp.Data.Metadata.Name,
				cloudServerResp.Data.Properties.Flavor.Name,
				cloudServerResp.Data.Properties.Zone)
		}
	} else {
		fmt.Println("⚠ Skipping cloud server creation - SSH key pair not available")
	}

	fmt.Println("\n=== SDK Example Complete ===")
	fmt.Println("Successfully created project:")
	fmt.Println("- Project ID:", projectID)
	fmt.Println("- ElasticIP ID:", *elasticIPResp.Data.Metadata.Id)
	fmt.Println("- Block Storage ID:", *blockStorageResp.Data.Metadata.Id)
	fmt.Println("- Snapshot ID:", *snapshotResp.Data.Metadata.Id)
	fmt.Println("- VPC ID:", *vpcResp.Data.Metadata.Id)
	if subnetResp != nil && subnetResp.Data != nil && subnetResp.Data.Metadata.Id != nil {
		fmt.Println("- Subnet ID:", *subnetResp.Data.Metadata.Id)
	}
	if sgResp != nil && sgResp.Data != nil && sgResp.Data.Metadata.Id != nil {
		fmt.Println("- Security Group ID:", *sgResp.Data.Metadata.Id)
	}
	if ruleResp != nil && ruleResp.Data != nil && ruleResp.Data.Metadata.Id != nil {
		fmt.Println("- Security Rule ID:", *ruleResp.Data.Metadata.Id)
	}
	if keyPairResp != nil && keyPairResp.Data != nil && keyPairResp.Data.Metadata.Name != "" {
		fmt.Println("- SSH Key Pair:", keyPairResp.Data.Metadata.Name)
	}
	if cloudServerResp != nil && cloudServerResp.Data != nil && cloudServerResp.Data.Metadata.Name != "" {
		fmt.Println("- Cloud Server:", cloudServerResp.Data.Metadata.Name)
	}

}

// Helper for pointer types
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
