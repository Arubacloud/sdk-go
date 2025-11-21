package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	aruba "github.com/Arubacloud/sdk-go"
	"github.com/Arubacloud/sdk-go/internal/restclient"
)

// runDeleteExample demonstrates how to delete all resources
// To run: PROJECT_ID=your-project go run . -mode=delete
func runDeleteExample() {
	config := &restclient.Config{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		Debug:        true,
	}

	// Initialize the SDK
	arubaClient, err := aruba.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	fmt.Println("\n=== Delete Example ===")

	// Get project ID from environment
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("Please set PROJECT_ID environment variable")
	}

	// Fetch all existing resources
	resources := fetchAllResources(ctx, arubaClient, projectID)

	// Confirm deletion
	fmt.Printf("\n⚠️  WARNING: This will delete ALL resources in project: %s\n", projectID)
	fmt.Print("Type 'yes' to confirm: ")
	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		log.Fatalf("Internal Error: %v", err)
	}
	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	// Delete all resources
	deleteAllResources(ctx, arubaClient, resources)

	fmt.Println("\n=== Delete Example Complete ===")
}

// fetchAllResources retrieves all existing resources from the API
func fetchAllResources(ctx context.Context, arubaClient aruba.Client, projectID string) *ResourceCollection {
	resources := &ResourceCollection{
		ProjectID: projectID,
	}

	fmt.Println("Fetching all resources...")

	// Fetch Cloud Servers
	serverList, err := arubaClient.FromCompute().CloudServers().List(ctx, projectID, nil)
	if err == nil && serverList.IsSuccess() && len(serverList.Data.Values) > 0 {
		// Get the first server details
		serverID := serverList.Data.Values[0].Metadata.Name
		serverResp, err := arubaClient.FromCompute().CloudServers().Get(ctx, projectID, serverID, nil)
		if err == nil && serverResp.IsSuccess() {
			resources.CloudServerResp = serverResp
			fmt.Printf("✓ Found Cloud Server: %s\n", serverResp.Data.Metadata.Name)
		}
	}

	// Fetch KaaS clusters
	kaasList, err := arubaClient.FromContainer().KaaS().List(ctx, projectID, nil)
	if err == nil && kaasList.IsSuccess() && len(kaasList.Data.Values) > 0 {
		kaasID := *kaasList.Data.Values[0].Metadata.ID
		kaasResp, err := arubaClient.FromContainer().KaaS().Get(ctx, projectID, kaasID, nil)
		if err == nil && kaasResp.IsSuccess() {
			resources.KaaSResp = kaasResp
			fmt.Printf("✓ Found KaaS: %s\n", *kaasResp.Data.Metadata.Name)
		}
	}

	// Fetch DBaaS instances
	dbaasListResp, err := arubaClient.FromDatabase().DBaaS().List(ctx, projectID, nil)
	if err == nil && dbaasListResp.IsSuccess() && len(dbaasListResp.Data.Values) > 0 {
		dbaasID := *dbaasListResp.Data.Values[0].Metadata.ID
		dbaasResp, err := arubaClient.FromDatabase().DBaaS().Get(ctx, projectID, dbaasID, nil)
		if err == nil && dbaasResp.IsSuccess() {
			resources.DBaaSResp = dbaasResp
			fmt.Printf("✓ Found DBaaS: %s\n", *dbaasResp.Data.Metadata.Name)
		}
	}

	// Fetch Key Pairs
	keyPairList, err := arubaClient.FromCompute().KeyPairs().List(ctx, projectID, nil)
	if err == nil && keyPairList.IsSuccess() && len(keyPairList.Data.Values) > 0 {
		keyPairID := *keyPairList.Data.Values[0].Metadata.ID
		keyPairResp, err := arubaClient.FromCompute().KeyPairs().Get(ctx, projectID, keyPairID, nil)
		if err == nil && keyPairResp.IsSuccess() {
			resources.KeyPairResp = keyPairResp
			fmt.Printf("✓ Found Key Pair: %s\n", *keyPairResp.Data.Metadata.Name)
		}
	}

	// Fetch VPCs and their resources
	vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, projectID, nil)
	if err == nil && vpcList.IsSuccess() && len(vpcList.Data.Values) > 0 {
		vpcID := *vpcList.Data.Values[0].Metadata.ID
		vpcResp, err := arubaClient.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)
		if err == nil && vpcResp.IsSuccess() {
			resources.VPCResp = vpcResp
			fmt.Printf("✓ Found VPC: %s\n", *vpcResp.Data.Metadata.Name)

			// Fetch Security Groups in VPC
			sgList, err := arubaClient.FromNetwork().SecurityGroups().List(ctx, projectID, vpcID, nil)
			if err == nil && sgList.IsSuccess() && len(sgList.Data.Values) > 0 {
				sgID := *sgList.Data.Values[0].Metadata.ID
				sgResp, err := arubaClient.FromNetwork().SecurityGroups().Get(ctx, projectID, vpcID, sgID, nil)
				if err == nil && sgResp.IsSuccess() {
					resources.SecurityGroupResp = sgResp
					fmt.Printf("✓ Found Security Group: %s\n", *sgResp.Data.Metadata.Name)

					// Fetch Security Group Rules
					ruleList, err := arubaClient.FromNetwork().SecurityGroupRules().List(ctx, projectID, vpcID, sgID, nil)
					if err == nil && ruleList.IsSuccess() && len(ruleList.Data.Values) > 0 {
						ruleID := *ruleList.Data.Values[0].Metadata.ID
						ruleResp, err := arubaClient.FromNetwork().SecurityGroupRules().Get(ctx, projectID, vpcID, sgID, ruleID, nil)
						if err == nil && ruleResp.IsSuccess() {
							resources.SecurityRuleResp = ruleResp
							fmt.Printf("✓ Found Security Rule: %s\n", *ruleResp.Data.Metadata.Name)
						}
					}
				}
			}

			// Fetch Subnets in VPC
			subnetList, err := arubaClient.FromNetwork().Subnets().List(ctx, projectID, vpcID, nil)
			if err == nil && subnetList.IsSuccess() && len(subnetList.Data.Values) > 0 {
				subnetID := *subnetList.Data.Values[0].Metadata.ID
				subnetResp, err := arubaClient.FromNetwork().Subnets().Get(ctx, projectID, vpcID, subnetID, nil)
				if err == nil && subnetResp.IsSuccess() {
					resources.SubnetResp = subnetResp
					fmt.Printf("✓ Found Subnet: %s\n", *subnetResp.Data.Metadata.Name)
				}
			}
		}
	}

	// Fetch Snapshots
	snapshotList, err := arubaClient.FromStorage().Snapshots().List(ctx, projectID, nil)
	if err == nil && snapshotList.IsSuccess() && len(snapshotList.Data.Values) > 0 {
		snapshotID := *snapshotList.Data.Values[0].Metadata.ID
		snapshotResp, err := arubaClient.FromStorage().Snapshots().Get(ctx, projectID, snapshotID, nil)
		if err == nil && snapshotResp.IsSuccess() {
			resources.SnapshotResp = snapshotResp
			fmt.Printf("✓ Found Snapshot: %s\n", *snapshotResp.Data.Metadata.Name)
		}
	}

	// Fetch Block Storage
	blockStorageList, err := arubaClient.FromStorage().Volumes().List(ctx, projectID, nil)
	if err == nil && blockStorageList.IsSuccess() && len(blockStorageList.Data.Values) > 0 {
		blockStorageID := *blockStorageList.Data.Values[0].Metadata.ID
		blockStorageResp, err := arubaClient.FromStorage().Volumes().Get(ctx, projectID, blockStorageID, nil)
		if err == nil && blockStorageResp.IsSuccess() {
			resources.BlockStorageResp = blockStorageResp
			fmt.Printf("✓ Found Block Storage: %s\n", *blockStorageResp.Data.Metadata.Name)
		}
	}

	// Fetch Elastic IPs
	elasticIPList, err := arubaClient.FromNetwork().ElasticIPs().List(ctx, projectID, nil)
	if err == nil && elasticIPList.IsSuccess() && len(elasticIPList.Data.Values) > 0 {
		elasticIPID := *elasticIPList.Data.Values[0].Metadata.ID
		elasticIPResp, err := arubaClient.FromNetwork().ElasticIPs().Get(ctx, projectID, elasticIPID, nil)
		if err == nil && elasticIPResp.IsSuccess() {
			resources.ElasticIPResp = elasticIPResp
			fmt.Printf("✓ Found Elastic IP: %s\n", *elasticIPResp.Data.Metadata.Name)
		}
	}

	return resources
}

// deleteAllResources deletes all resources in reverse order of creation
// This ensures dependencies are respected (e.g., delete cloud server before VPC)
func deleteAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Deleting Resources ===")

	// Delete in reverse order of creation to respect dependencies

	// 12. Delete Cloud Server (if created)
	if resources.CloudServerResp != nil && resources.CloudServerResp.Data != nil && resources.CloudServerResp.Data.Metadata.Name != "" {
		deleteCloudServer(ctx, arubaClient, resources.ProjectID, &resources.CloudServerResp.Data.Metadata.Name)
	}

	// 11. Delete KaaS (if created)
	if resources.KaaSResp != nil && resources.KaaSResp.Data != nil {
		deleteKaaS(ctx, arubaClient, resources.ProjectID, *resources.KaaSResp.Data.Metadata.ID)
	}

	// 10. Delete DBaaS (if created)
	if resources.DBaaSResp != nil && resources.DBaaSResp.Data != nil {
		deleteDBaaS(ctx, arubaClient, resources.ProjectID, *resources.DBaaSResp.Data.Metadata.ID)
	}

	// 9. Delete SSH Key Pair (if created)
	if resources.KeyPairResp != nil && resources.KeyPairResp.Data != nil {
		deleteKeyPair(ctx, arubaClient, resources.ProjectID, *resources.KeyPairResp.Data.Metadata.ID)
	}

	// 8. Delete Security Group Rule (if created)
	if resources.SecurityRuleResp != nil && resources.SecurityRuleResp.Data != nil && resources.VPCResp != nil {
		deleteSecurityGroupRule(ctx, arubaClient, resources.ProjectID, *resources.VPCResp.Data.Metadata.ID,
			*resources.SecurityGroupResp.Data.Metadata.ID, *resources.SecurityRuleResp.Data.Metadata.ID)
	}

	// 7. Delete Security Group (if created)
	if resources.SecurityGroupResp != nil && resources.SecurityGroupResp.Data != nil && resources.VPCResp != nil {
		deleteSecurityGroup(ctx, arubaClient, resources.ProjectID, *resources.VPCResp.Data.Metadata.ID,
			*resources.SecurityGroupResp.Data.Metadata.ID)
	}

	// 6. Delete Subnet (if created)
	if resources.SubnetResp != nil && resources.SubnetResp.Data != nil && resources.VPCResp != nil {
		deleteSubnet(ctx, arubaClient, resources.ProjectID, *resources.VPCResp.Data.Metadata.ID,
			*resources.SubnetResp.Data.Metadata.ID)
	}

	// 5. Delete VPC (if created)
	if resources.VPCResp != nil && resources.VPCResp.Data != nil {
		deleteVPC(ctx, arubaClient, resources.ProjectID, *resources.VPCResp.Data.Metadata.ID)
	}

	// 4. Delete Snapshot (if created)
	if resources.SnapshotResp != nil && resources.SnapshotResp.Data != nil {
		deleteSnapshot(ctx, arubaClient, resources.ProjectID, *resources.SnapshotResp.Data.Metadata.ID)
	}

	// 3. Delete Block Storage (if created)
	if resources.BlockStorageResp != nil && resources.BlockStorageResp.Data != nil {
		deleteBlockStorage(ctx, arubaClient, resources.ProjectID, *resources.BlockStorageResp.Data.Metadata.ID)
	}

	// 2. Delete Elastic IP (if created)
	if resources.ElasticIPResp != nil && resources.ElasticIPResp.Data != nil {
		deleteElasticIP(ctx, arubaClient, resources.ProjectID, *resources.ElasticIPResp.Data.Metadata.ID)
	}

	// 1. Delete Project (last - after all resources are deleted)
	if resources.ProjectID != "" {
		deleteProject(ctx, arubaClient, resources.ProjectID)
	}

	fmt.Println("\n=== Delete Complete ===")
}

// deleteProject deletes a project
func deleteProject(ctx context.Context, arubaClient aruba.Client, projectID string) {
	fmt.Println("--- Deleting Project ---")

	deleteResp, err := arubaClient.FromProject().Delete(ctx, projectID, nil)
	if err != nil {
		log.Printf("Error deleting project: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete project - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted project: %s\n", projectID)
}

// deleteElasticIP deletes an Elastic IP
func deleteElasticIP(ctx context.Context, arubaClient aruba.Client, projectID, elasticIPID string) {
	fmt.Println("--- Deleting Elastic IP ---")

	deleteResp, err := arubaClient.FromNetwork().ElasticIPs().Delete(ctx, projectID, elasticIPID, nil)
	if err != nil {
		log.Printf("Error deleting Elastic IP: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete Elastic IP - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted Elastic IP: %s\n", elasticIPID)
}

// deleteBlockStorage deletes a block storage volume
func deleteBlockStorage(ctx context.Context, arubaClient aruba.Client, projectID, blockStorageID string) {
	fmt.Println("--- Deleting Block Storage ---")

	deleteResp, err := arubaClient.FromStorage().Volumes().Delete(ctx, projectID, blockStorageID, nil)
	if err != nil {
		log.Printf("Error deleting block storage: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete block storage - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted block storage: %s\n", blockStorageID)
}

// deleteSnapshot deletes a snapshot
func deleteSnapshot(ctx context.Context, arubaClient aruba.Client, projectID, snapshotID string) {
	fmt.Println("--- Deleting Snapshot ---")

	deleteResp, err := arubaClient.FromStorage().Snapshots().Delete(ctx, projectID, snapshotID, nil)
	if err != nil {
		log.Printf("Error deleting snapshot: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete snapshot - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted snapshot: %s\n", snapshotID)
}

// deleteVPC deletes a VPC
func deleteVPC(ctx context.Context, arubaClient aruba.Client, projectID, vpcID string) {
	fmt.Println("--- Deleting VPC ---")

	deleteResp, err := arubaClient.FromNetwork().VPCs().Delete(ctx, projectID, vpcID, nil)
	if err != nil {
		log.Printf("Error deleting VPC: %v", err)
		return
	} else if deleteResp.IsError() {
		log.Printf("Failed to delete VPC - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted VPC: %s\n", vpcID)
}

// deleteSubnet deletes a subnet
func deleteSubnet(ctx context.Context, arubaClient aruba.Client, projectID, vpcID, subnetID string) {
	fmt.Println("--- Deleting Subnet ---")

	deleteResp, err := arubaClient.FromNetwork().Subnets().Delete(ctx, projectID, vpcID, subnetID, nil)
	if err != nil {
		log.Printf("Error deleting subnet: %v", err)
		return
	} else if deleteResp.IsError() {
		log.Printf("Failed to delete subnet - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted subnet: %s\n", subnetID)
}

// deleteSecurityGroup deletes a security group
func deleteSecurityGroup(ctx context.Context, arubaClient aruba.Client, projectID, vpcID, securityGroupID string) {
	fmt.Println("--- Deleting Security Group ---")

	deleteResp, err := arubaClient.FromNetwork().SecurityGroups().Delete(ctx, projectID, vpcID, securityGroupID, nil)
	if err != nil {
		log.Printf("Error deleting security group: %v", err)
		return
	} else if deleteResp.IsError() {
		log.Printf("Failed to delete security group - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted security group: %s\n", securityGroupID)
}

// deleteSecurityGroupRule deletes a security group rule
func deleteSecurityGroupRule(ctx context.Context, arubaClient aruba.Client, projectID, vpcID, securityGroupID, ruleID string) {
	fmt.Println("--- Deleting Security Group Rule ---")

	deleteResp, err := arubaClient.FromNetwork().SecurityGroupRules().Delete(ctx, projectID, vpcID, securityGroupID, ruleID, nil)
	if err != nil {
		log.Printf("Error deleting security rule: %v", err)
		return
	} else if deleteResp.IsError() {
		log.Printf("Failed to delete security rule - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted security rule: %s\n", ruleID)
}

// deleteKeyPair deletes an SSH key pair
func deleteKeyPair(ctx context.Context, arubaClient aruba.Client, projectID, keyPairID string) {
	fmt.Println("--- Deleting SSH Key Pair ---")

	deleteResp, err := arubaClient.FromCompute().KeyPairs().Delete(ctx, projectID, keyPairID, nil)
	if err != nil {
		log.Printf("Error deleting SSH key pair: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete SSH key pair - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted SSH key pair: %s\n", keyPairID)
}

// deleteDBaaS deletes a DBaaS instance
func deleteDBaaS(ctx context.Context, arubaClient aruba.Client, projectID, dbaasID string) {
	fmt.Println("--- Deleting DBaaS ---")

	deleteResp, err := arubaClient.FromDatabase().DBaaS().Delete(ctx, projectID, dbaasID, nil)
	if err != nil {
		log.Printf("Error deleting DBaaS: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete DBaaS - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted DBaaS: %s\n", dbaasID)
}

// deleteKaaS deletes a KaaS cluster
func deleteKaaS(ctx context.Context, arubaClient aruba.Client, projectID, kaasID string) {
	fmt.Println("--- Deleting KaaS Cluster ---")

	deleteResp, err := arubaClient.FromContainer().KaaS().Delete(ctx, projectID, kaasID, nil)
	if err != nil {
		log.Printf("Error deleting KaaS cluster: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete KaaS cluster - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted KaaS cluster: %s\n", kaasID)
}

// deleteCloudServer deletes a cloud server
func deleteCloudServer(ctx context.Context, arubaClient aruba.Client, projectID string, cloudServerID *string) {
	if cloudServerID == nil {
		log.Println("Cloud Server ID is nil, skipping deletion")
		return
	}

	fmt.Println("--- Deleting Cloud Server ---")

	deleteResp, err := arubaClient.FromCompute().CloudServers().Delete(ctx, projectID, *cloudServerID, nil)
	if err != nil {
		log.Printf("Error deleting cloud server: %v", err)
		return
	} else if !deleteResp.IsSuccess() {
		log.Printf("Failed to delete cloud server - Status: %d, Error: %s",
			deleteResp.StatusCode,
			stringValue(deleteResp.Error.Title))
		return
	}
	fmt.Printf("✓ Deleted cloud server: %s\n", *cloudServerID)
}
