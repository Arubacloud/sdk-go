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

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	fmt.Println("\n=== Delete Example ===")

	proj, err := arubaClient.FromProject().Get(ctx, aruba.URI("/projects/"+projectID))
	if err != nil {
		log.Fatalf("Error fetching project: %v", err)
	}

	resources := fetchAllResources(ctx, arubaClient, proj)

	printDeletionInventory(resources)

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

	snapshotList, err := arubaClient.FromStorage().Snapshots().List(ctx, proj)
	if err == nil && snapshotList.Total() > 0 {
		snapResp, err := arubaClient.FromStorage().Snapshots().Get(ctx, snapshotList.Items()[0])
		if err == nil {
			resources.Snapshot = snapResp
			fmt.Printf("✓ Found Snapshot: %s\n", snapResp.Name())
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

	// All BlockStorages in the project (boot volumes, attached, unattached — catches every stray).
	bsList, err := arubaClient.FromStorage().Volumes().List(ctx, proj)
	if err == nil {
		for _, bs := range bsList.Items() {
			resources.BlockStorages = append(resources.BlockStorages, bs)
			fmt.Printf("✓ Found Block Storage: %s\n", bs.Name())
		}
	}

	// All ElasticIPs in the project.
	eipList, err := arubaClient.FromNetwork().ElasticIPs().List(ctx, proj)
	if err == nil {
		for _, eip := range eipList.Items() {
			resources.ElasticIPs = append(resources.ElasticIPs, eip)
			fmt.Printf("✓ Found Elastic IP: %s\n", eip.Name())
		}
	}

	return resources
}

// printDeletionInventory prints a concise summary of everything fetchAllResources found.
func printDeletionInventory(r *ResourceCollection) {
	fmt.Println("\n=== Inventory ===")

	count := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	fmt.Printf("  Project: %d\n", count(r.Project != nil))

	fmt.Printf("  KMS: %d  KMS Key: %d  KMIP: %d\n",
		count(r.KMS != nil), count(r.KMSKey != nil), count(r.Kmip != nil))

	fmt.Printf("  Backup: %d  Restore: %d  Snapshot: %d\n",
		count(r.Backup != nil), count(r.Restore != nil), count(r.Snapshot != nil))

	fmt.Printf("  Container Registry: %d\n", count(r.ContainerRegistry != nil))

	jobs := 0
	if r.JobRecurring != nil {
		jobs++
	}
	if r.JobOneShot != nil {
		jobs++
	}
	fmt.Printf("  Cloud Server: %d  KaaS: %d  Jobs: %d\n",
		count(r.CloudServer != nil), count(r.KaaS != nil), jobs)

	grants := count(r.Grant != nil)
	users := count(r.DBaaSUser != nil)
	dbs := count(r.Database != nil)
	fmt.Printf("  DBaaS: %d (Database: %d  User: %d  Grant: %d)\n",
		count(r.DBaaS != nil), dbs, users, grants)

	fmt.Printf("  VPC: %d  Subnets: %d  Security Group: %d  Rules: %d  Key Pair: %d\n",
		count(r.VPC != nil),
		count(r.SubnetAdvanced != nil)+count(r.SubnetBasic != nil),
		count(r.SecurityGroup != nil),
		len(r.SecurityRulesIngress)+count(r.SecurityRuleEgress != nil),
		count(r.KeyPair != nil))

	fmt.Printf("  Block Storages: %d  Elastic IPs: %d\n",
		len(r.BlockStorages), len(r.ElasticIPs))
}

// deleteAllResources deletes all resources in strict inverse of the creation order
// defined in orchestrator_create.go. Each resource waits until fully gone (HTTP 404)
// before the next dependent delete proceeds — this prevents the API from rejecting
// deletes of resources still referenced by a peer in "Deleting" state.
func deleteAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Deleting Resources ===")

	// Phase 1/7: KMS stack (inverse of create Phase 7: KMS → KMSKey → KMIP).
	printPhase(1, 7, "KMS stack")

	if resources.Kmip != nil && resources.Kmip.KmipID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKmip(c, arubaClient, resources.Kmip) })
	}
	if resources.KMSKey != nil && resources.KMSKey.KeyID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKMSKey(c, arubaClient, resources.KMSKey) })
	}
	if resources.KMS != nil && resources.KMS.KMSID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKMS(c, arubaClient, resources.KMS) })
	}

	// Phase 2/7: Backup & restore (inverse of create Phase 6: Restore → Backup).
	// RestoreTargetStorage is a BlockStorage; it falls into Phase 6 with all other volumes.
	printPhase(2, 7, "Backup & restore")

	if resources.Restore != nil && resources.Restore.RestoreID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteRestore(c, arubaClient, resources.Restore) })
	}
	if resources.Backup != nil && resources.Backup.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteStorageBackup(c, arubaClient, resources.Backup) })
	}

	// Phase 3/7: Compute & container platforms (inverse of create Phase 5:
	// ContainerRegistry → Jobs → CloudServer → KaaS).
	printPhase(3, 7, "Compute & container platforms")

	if resources.ContainerRegistry != nil && resources.ContainerRegistry.ContainerRegistryID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) {
			deleteContainerRegistry(c, arubaClient, resources.ContainerRegistry)
		})
	}
	if resources.JobOneShot != nil && resources.JobOneShot.JobID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteJob(c, arubaClient, resources.JobOneShot, "One-Shot") })
	}
	if resources.JobRecurring != nil && resources.JobRecurring.JobID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteJob(c, arubaClient, resources.JobRecurring, "Recurring") })
	}
	if resources.CloudServer != nil && resources.CloudServer.CloudServerID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteCloudServer(c, arubaClient, resources.CloudServer) })
	}
	if resources.KaaS != nil && resources.KaaS.KaaSID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKaaS(c, arubaClient, resources.KaaS) })
	}

	// Phase 4/7: Database stack (inverse of create Phase 4:
	// Grant → DBaaSUser → Database → DBaaS).
	printPhase(4, 7, "Database stack")

	if resources.Grant != nil && resources.Grant.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteGrant(c, arubaClient, resources.Grant) })
	}
	if resources.DBaaSUser != nil && resources.DBaaSUser.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteDBaaSUser(c, arubaClient, resources.DBaaSUser) })
	}
	if resources.Database != nil && resources.Database.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteDatabase(c, arubaClient, resources.Database) })
	}
	if resources.DBaaS != nil && resources.DBaaS.DBaaSID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteDBaaS(c, arubaClient, resources.DBaaS) })
	}

	// Phase 5/7: VPC-scoped network (inverse of create Phase 3:
	// KeyPair → Rules (egress, ingress) → SecurityGroup → Subnets (Basic, Advanced)).
	printPhase(5, 7, "VPC-scoped network")

	if resources.KeyPair != nil && resources.KeyPair.KeyPairID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKeyPair(c, arubaClient, resources.KeyPair) })
	}
	if resources.SecurityRuleEgress != nil && resources.SecurityRuleEgress.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) {
			deleteSecurityGroupRule(c, arubaClient, resources.SecurityRuleEgress)
		})
	}
	for _, r := range resources.SecurityRulesIngress {
		if r != nil && r.ID() != "" {
			r := r
			withDeleteDeadline(ctx, func(c context.Context) { deleteSecurityGroupRule(c, arubaClient, r) })
		}
	}
	if resources.SecurityGroup != nil && resources.SecurityGroup.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteSecurityGroup(c, arubaClient, resources.SecurityGroup) })
	}
	if resources.SubnetBasic != nil {
		withDeleteDeadline(ctx, func(c context.Context) { deleteBasicSubnet(c, arubaClient, resources.SubnetBasic) })
	}
	if resources.SubnetAdvanced != nil {
		withDeleteDeadline(ctx, func(c context.Context) { deleteAdvancedSubnet(c, arubaClient, resources.SubnetAdvanced) })
	}

	// Phase 6/7: Independent network & storage primitives (inverse of create Phase 2:
	// VPC → Snapshot → BlockStorages → ElasticIPs).
	printPhase(6, 7, "Independent network & storage primitives")

	if resources.VPC != nil {
		withDeleteDeadline(ctx, func(c context.Context) { deleteVPC(c, arubaClient, resources.VPC) })
	}
	if resources.Snapshot != nil && resources.Snapshot.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteSnapshot(c, arubaClient, resources.Snapshot) })
	}
	for _, bs := range resources.BlockStorages {
		bs := bs
		withDeleteDeadline(ctx, func(c context.Context) { deleteBlockStorage(c, arubaClient, bs) })
	}
	for _, eip := range resources.ElasticIPs {
		eip := eip
		withDeleteDeadline(ctx, func(c context.Context) { deleteElasticIP(c, arubaClient, eip) })
	}

	// Phase 7/7: Account & isolation (inverse of create Phase 1).
	printPhase(7, 7, "Account & isolation")

	if resources.Project != nil {
		withDeleteDeadline(ctx, func(c context.Context) { deleteProject(c, arubaClient, resources.Project) })
	}

	fmt.Println("\n=== Delete Complete ===")
}
