package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// runDeleteExample demonstrates how to delete all resources.
// Run with: go run ./examples/all-resources/ -mode=delete -clientID=… -clientSecret=… -projectID=…
func runDeleteExample(clientID, clientSecret, projectID string, debug bool) {
	arubaClient, err := aruba.NewClient(buildClientOptions(clientID, clientSecret, debug))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	fmt.Println("\n=== Delete Example ===")

	proj, err := arubaClient.FromProject().Get(ctx, aruba.URI("/projects/"+projectID))
	if err != nil {
		log.Fatalf("Error fetching project: %v", err)
	}

	resources := fetchAllResources(ctx, arubaClient, proj)

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

	deleteAllResources(ctx, arubaClient, resources)

	fmt.Println("\n=== Delete Example Complete ===")
}

func fetchAllResources(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) *ResourceCollection {
	resources := &ResourceCollection{
		Project: proj,
	}

	fmt.Println("Fetching all resources...")

	serverList, err := arubaClient.FromCompute().CloudServers().List(ctx, proj)
	if err == nil && serverList.Total() > 0 {
		firstServer := serverList.Items()[0]
		serverResp, err := arubaClient.FromCompute().CloudServers().Get(ctx, firstServer)
		if err == nil {
			resources.CloudServer = serverResp
			fmt.Printf("✓ Found Cloud Server: %s\n", serverResp.Name())
		}
	}

	kaasList, err := arubaClient.FromContainer().KaaS().List(ctx, proj)
	if err == nil && kaasList.Total() > 0 {
		kaasResp, err := arubaClient.FromContainer().KaaS().Get(ctx, kaasList.Items()[0])
		if err == nil {
			resources.KaaS = kaasResp
			fmt.Printf("✓ Found KaaS: %s\n", kaasResp.Name())
		}
	}

	dbaasListResp, err := arubaClient.FromDatabase().DBaaS().List(ctx, proj)
	if err == nil && dbaasListResp.Total() > 0 {
		dbaasResp, err := arubaClient.FromDatabase().DBaaS().Get(ctx, dbaasListResp.Items()[0])
		if err == nil {
			resources.DBaaS = dbaasResp
			fmt.Printf("✓ Found DBaaS: %s\n", dbaasResp.Name())

			dbList, err := arubaClient.FromDatabase().Databases().List(ctx, dbaasResp)
			if err == nil && dbList.Total() > 0 {
				resources.Database = dbList.Items()[0]
				fmt.Printf("✓ Found Database: %s\n", resources.Database.Name())

				grantList, err := arubaClient.FromDatabase().Grants().List(ctx, resources.Database)
				if err == nil && grantList.Total() > 0 {
					resources.Grant = grantList.Items()[0]
					fmt.Printf("✓ Found Grant: %s\n", resources.Grant.ID())
				}
			}

			userList, err := arubaClient.FromDatabase().Users().List(ctx, dbaasResp)
			if err == nil && userList.Total() > 0 {
				resources.DBaaSUser = userList.Items()[0]
				fmt.Printf("✓ Found DBaaS User: %s\n", resources.DBaaSUser.Username())
			}
		}
	}

	jobList, err := arubaClient.FromSchedule().Jobs().List(ctx, proj)
	if err == nil {
		for _, j := range jobList.Items() {
			switch j.JobType() {
			case aruba.JobTypeRecurring:
				resources.JobRecurring = j
				fmt.Printf("✓ Found Recurring Job: %s\n", j.Name())
			case aruba.JobTypeOneShot:
				resources.JobOneShot = j
				fmt.Printf("✓ Found OneShot Job: %s\n", j.Name())
			}
		}
	}

	keyPairList, err := arubaClient.FromCompute().KeyPairs().List(ctx, proj)
	if err == nil && keyPairList.Total() > 0 {
		firstKP := keyPairList.Items()[0]
		keyPairResp, err := arubaClient.FromCompute().KeyPairs().Get(ctx, firstKP)
		if err == nil {
			resources.KeyPair = keyPairResp
			fmt.Printf("✓ Found Key Pair: %s\n", keyPairResp.Name())
		}
	}

	vpcList, err := arubaClient.FromNetwork().VPCs().List(ctx, proj)
	if err == nil && vpcList.Total() > 0 {
		firstVPC := vpcList.Items()[0]
		vpcResp, err := arubaClient.FromNetwork().VPCs().Get(ctx, firstVPC)
		if err == nil {
			resources.VPC = vpcResp
			fmt.Printf("✓ Found VPC: %s\n", vpcResp.Name())

			sgList, err := arubaClient.FromNetwork().SecurityGroups().List(ctx, firstVPC)
			if err == nil && sgList.Total() > 0 {
				sg := sgList.Items()[0]
				resources.SecurityGroup = sg
				fmt.Printf("✓ Found Security Group: %s\n", sg.Name())

				ruleList, err := arubaClient.FromNetwork().SecurityGroupRules().List(ctx, sg)
				if err == nil {
					for _, r := range ruleList.Items() {
						switch r.Direction() {
						case aruba.RuleDirectionIngress:
							resources.SecurityRulesIngress = append(resources.SecurityRulesIngress, r)
							fmt.Printf("✓ Found Ingress Rule: %s\n", r.Name())
						case aruba.RuleDirectionEgress:
							resources.SecurityRuleEgress = r
							fmt.Printf("✓ Found Egress Rule: %s\n", r.Name())
						}
					}
				}
			}

			subnetList, err := arubaClient.FromNetwork().Subnets().List(ctx, resources.VPC)
			if err == nil {
				for _, sn := range subnetList.Items() {
					switch sn.Type() {
					case aruba.SubnetTypeBasic:
						resources.SubnetBasic = sn
						fmt.Printf("✓ Found Basic Subnet: %s\n", sn.Name())
					case aruba.SubnetTypeAdvanced:
						resources.SubnetAdvanced = sn
						fmt.Printf("✓ Found Advanced Subnet: %s\n", sn.Name())
					}
				}
			}
		}
	}

	backupList, err := arubaClient.FromStorage().Backups().List(ctx, proj)
	if err == nil && backupList.Total() > 0 {
		backupItem := backupList.Items()[0]
		backupResp, err := arubaClient.FromStorage().Backups().Get(ctx, backupItem)
		if err == nil {
			resources.Backup = backupResp
			fmt.Printf("✓ Found Backup: %s\n", backupResp.Name())
		}
		restoreList, err := arubaClient.FromStorage().Restores().List(ctx, backupItem)
		if err == nil && restoreList.Total() > 0 {
			restoreItem := restoreList.Items()[0]
			restoreResp, err := arubaClient.FromStorage().Restores().Get(ctx, restoreItem)
			if err == nil {
				resources.Restore = restoreResp
				fmt.Printf("✓ Found Restore: %s\n", restoreResp.Name())
			}
		}
	}

	containerRegistryList, err := arubaClient.FromContainer().ContainerRegistry().List(ctx, proj)
	if err == nil && containerRegistryList.Total() > 0 {
		first := containerRegistryList.Items()[0]
		containerRegistryResp, err := arubaClient.FromContainer().ContainerRegistry().Get(ctx, first)
		if err == nil {
			resources.ContainerRegistry = containerRegistryResp
			fmt.Printf("✓ Found Container Registry: %s\n", containerRegistryResp.Name())
		}
	}

	kmsList, err := arubaClient.FromSecurity().KMS().List(ctx, proj)
	if err == nil && kmsList.Total() > 0 {
		kmsResp, err := arubaClient.FromSecurity().KMS().Get(ctx, kmsList.Items()[0])
		if err == nil {
			resources.KMS = kmsResp
			fmt.Printf("✓ Found KMS: %s\n", kmsResp.Name())

			keysList, err := arubaClient.FromSecurity().Keys().List(ctx, kmsResp)
			if err == nil && keysList.Total() > 0 {
				keyResp, err := arubaClient.FromSecurity().Keys().Get(ctx, keysList.Items()[0])
				if err == nil {
					resources.KMSKey = keyResp
					fmt.Printf("✓ Found KMS Key: %s\n", keyResp.Name())
				}
			}

			kmipList, err := arubaClient.FromSecurity().Kmips().List(ctx, kmsResp)
			if err == nil && kmipList.Total() > 0 {
				kmipResp, err := arubaClient.FromSecurity().Kmips().Get(ctx, kmipList.Items()[0])
				if err == nil {
					resources.Kmip = kmipResp
					fmt.Printf("✓ Found KMIP: %s\n", kmipResp.Name())
				}
			}
		}
	}

	return resources
}

// deleteAllResources deletes all resources in reverse order of creation
// to ensure dependencies are respected.
func deleteAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Deleting Resources ===")

	if resources.Restore != nil && resources.Restore.RestoreID() != "" {
		deleteRestore(ctx, arubaClient, resources.Restore)
	}

	if resources.Backup != nil && resources.Backup.ID() != "" {
		deleteStorageBackup(ctx, arubaClient, resources.Backup)
	}

	if resources.Snapshot != nil && resources.Snapshot.ID() != "" {
		deleteSnapshot(ctx, arubaClient, resources.Snapshot)
	}

	if resources.ContainerRegistry != nil && resources.ContainerRegistry.ContainerRegistryID() != "" {
		deleteContainerRegistry(ctx, arubaClient, resources.ContainerRegistry)
	}

	if resources.KMS != nil && resources.KMS.KMSID() != "" {
		if resources.Kmip != nil && resources.Kmip.KmipID() != "" {
			deleteKmip(ctx, arubaClient, resources.Kmip)
		}
		if resources.KMSKey != nil && resources.KMSKey.KeyID() != "" {
			deleteKMSKey(ctx, arubaClient, resources.KMSKey)
		}
		deleteKMS(ctx, arubaClient, resources.KMS)
	}

	if resources.JobRecurring != nil && resources.JobRecurring.JobID() != "" {
		deleteJob(ctx, arubaClient, resources.JobRecurring, "Recurring")
	}
	if resources.JobOneShot != nil && resources.JobOneShot.JobID() != "" {
		deleteJob(ctx, arubaClient, resources.JobOneShot, "OneShot")
	}

	if resources.CloudServer != nil && resources.CloudServer.CloudServerID() != "" {
		deleteCloudServer(ctx, arubaClient, resources.CloudServer)
	}

	if resources.KaaS != nil && resources.KaaS.KaaSID() != "" {
		deleteKaaS(ctx, arubaClient, resources.KaaS)
	}

	if resources.Grant != nil && resources.Grant.ID() != "" {
		deleteGrant(ctx, arubaClient, resources.Grant)
	}
	if resources.DBaaSUser != nil && resources.DBaaSUser.ID() != "" {
		deleteDBaaSUser(ctx, arubaClient, resources.DBaaSUser)
	}
	if resources.Database != nil && resources.Database.ID() != "" {
		deleteDatabase(ctx, arubaClient, resources.Database)
	}

	if resources.DBaaS != nil && resources.DBaaS.DBaaSID() != "" {
		deleteDBaaS(ctx, arubaClient, resources.DBaaS)
	}

	if resources.KeyPair != nil && resources.KeyPair.KeyPairID() != "" {
		deleteKeyPair(ctx, arubaClient, resources.KeyPair)
	}

	for _, r := range resources.SecurityRulesIngress {
		if r != nil && r.ID() != "" {
			deleteSecurityGroupRule(ctx, arubaClient, r)
		}
	}
	if resources.SecurityRuleEgress != nil && resources.SecurityRuleEgress.ID() != "" {
		deleteSecurityGroupRule(ctx, arubaClient, resources.SecurityRuleEgress)
	}

	if resources.SecurityGroup != nil && resources.SecurityGroup.ID() != "" {
		deleteSecurityGroup(ctx, arubaClient, resources.SecurityGroup)
	}

	if resources.SubnetBasic != nil {
		deleteBasicSubnet(ctx, arubaClient, resources.SubnetBasic)
	}
	if resources.SubnetAdvanced != nil {
		deleteAdvancedSubnet(ctx, arubaClient, resources.SubnetAdvanced)
	}

	if resources.VPC != nil {
		deleteVPC(ctx, arubaClient, resources.VPC)
	}

	bsList, err := arubaClient.FromStorage().Volumes().List(ctx, resources.Project)
	if err == nil {
		for _, bs := range bsList.Items() {
			deleteBlockStorage(ctx, arubaClient, bs)
		}
	}

	eipList, err := arubaClient.FromNetwork().ElasticIPs().List(ctx, resources.Project)
	if err == nil {
		for _, eip := range eipList.Items() {
			deleteElasticIP(ctx, arubaClient, eip)
		}
	}

	if resources.Project != nil {
		deleteProject(ctx, arubaClient, resources.Project)
	}

	fmt.Println("\n=== Delete Complete ===")
}
