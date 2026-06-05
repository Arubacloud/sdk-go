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

	// VPN tunnel + first route (the delete orchestrator deletes the route then the tunnel).
	vpnTunnelList, err := arubaClient.FromNetwork().VPNTunnels().List(ctx, proj)
	if err == nil && vpnTunnelList.Total() > 0 {
		tunnelRef := vpnTunnelList.Items()[0]
		tunnelResp, err := arubaClient.FromNetwork().VPNTunnels().Get(ctx, tunnelRef)
		if err == nil {
			resources.VPNTunnel = tunnelResp
			fmt.Printf("✓ Found VPN Tunnel: %s\n", tunnelResp.Name())

			routeList, err := arubaClient.FromNetwork().VPNRoutes().List(ctx, tunnelResp)
			if err == nil && routeList.Total() > 0 {
				routeResp, err := arubaClient.FromNetwork().VPNRoutes().Get(ctx, routeList.Items()[0])
				if err == nil {
					resources.VPNRoute = routeResp
					fmt.Printf("✓ Found VPN Route: %s (cloudSubnet=%s)\n",
						routeResp.Name(), routeResp.CloudSubnetCIDR())
				}
			}
		}
	}

	return resources
}

// printDeletionInventory lists every resource found by fetchAllResources, with
// CreatedBy and CreatedAt where the resource carries audit metadata.
func printDeletionInventory(r *ResourceCollection) {
	fmt.Println("\n=== Inventory ===")

	// Project
	if r.Project != nil {
		fmt.Println(summaryRow("Project", r.Project.Name(), r.Project.ID(),
			auditExtras(r.Project.CreatedBy(), r.Project.CreatedAt())...))
	} else {
		fmt.Println(summaryRow("Project", "", ""))
	}

	// VPC & network
	if r.VPC != nil {
		fmt.Println(summaryRow("VPC", r.VPC.Name(), r.VPC.ID(),
			auditExtras(r.VPC.CreatedBy(), r.VPC.CreatedAt())...))
	}
	if r.SubnetAdvanced != nil {
		fmt.Println(summaryRow("Subnet (Advanced)", r.SubnetAdvanced.Name(), r.SubnetAdvanced.ID(),
			auditExtras(r.SubnetAdvanced.CreatedBy(), r.SubnetAdvanced.CreatedAt())...))
	}
	if r.SubnetBasic != nil {
		fmt.Println(summaryRow("Subnet (Basic)", r.SubnetBasic.Name(), r.SubnetBasic.ID(),
			auditExtras(r.SubnetBasic.CreatedBy(), r.SubnetBasic.CreatedAt())...))
	}
	if r.SecurityGroup != nil {
		fmt.Println(summaryRow("Security Group", r.SecurityGroup.Name(), r.SecurityGroup.ID(),
			auditExtras(r.SecurityGroup.CreatedBy(), r.SecurityGroup.CreatedAt())...))
	}
	for _, rule := range r.SecurityRulesIngress {
		if rule != nil {
			fmt.Println(summaryRow("Security Rule (Ingress/"+rule.Name()+")", rule.Name(), rule.ID(),
				auditExtras(rule.CreatedBy(), rule.CreatedAt())...))
		}
	}
	if r.SecurityRuleEgress != nil {
		fmt.Println(summaryRow("Security Rule (Egress)", r.SecurityRuleEgress.Name(), r.SecurityRuleEgress.ID(),
			auditExtras(r.SecurityRuleEgress.CreatedBy(), r.SecurityRuleEgress.CreatedAt())...))
	}
	if r.KeyPair != nil {
		fmt.Println(summaryRow("SSH Key Pair", r.KeyPair.Name(), r.KeyPair.KeyPairID(),
			auditExtras(r.KeyPair.CreatedBy(), r.KeyPair.CreatedAt())...))
	}
	for _, eip := range r.ElasticIPs {
		if eip != nil {
			fmt.Println(summaryRow("Elastic IP", eip.Name(), eip.ID(),
				auditExtras(eip.CreatedBy(), eip.CreatedAt())...))
		}
	}
	for _, bs := range r.BlockStorages {
		if bs != nil {
			fmt.Println(summaryRow("Block Storage", bs.Name(), bs.ID(),
				auditExtras(bs.CreatedBy(), bs.CreatedAt())...))
		}
	}

	// VPN stack
	if r.VPNTunnel != nil {
		fmt.Println(summaryRow("VPN Tunnel", r.VPNTunnel.Name(), r.VPNTunnel.VPNTunnelID(),
			auditExtras(r.VPNTunnel.CreatedBy(), r.VPNTunnel.CreatedAt())...))
	}
	if r.VPNRoute != nil {
		extras := append([]string{"cloudSubnet=" + r.VPNRoute.CloudSubnetCIDR()},
			auditExtras(r.VPNRoute.CreatedBy(), r.VPNRoute.CreatedAt())...)
		fmt.Println(summaryRow("VPN Route", r.VPNRoute.Name(), r.VPNRoute.VPNRouteID(), extras...))
	}

	// KMS stack
	if r.KMS != nil {
		fmt.Println(summaryRow("KMS Instance", r.KMS.Name(), r.KMS.KMSID(),
			auditExtras(r.KMS.CreatedBy(), r.KMS.CreatedAt())...))
	}
	if r.KMSKey != nil {
		fmt.Println(summaryRow("KMS Key", r.KMSKey.Name(), r.KMSKey.KeyID()))
	}
	if r.Kmip != nil {
		fmt.Println(summaryRow("KMIP Service", r.Kmip.Name(), r.Kmip.KmipID()))
	}

	// Backup & restore
	if r.Backup != nil {
		fmt.Println(summaryRow("Storage Backup", r.Backup.Name(), r.Backup.BackupID(),
			auditExtras(r.Backup.CreatedBy(), r.Backup.CreatedAt())...))
	}
	if r.Restore != nil {
		fmt.Println(summaryRow("Storage Restore", r.Restore.Name(), r.Restore.RestoreID(),
			auditExtras(r.Restore.CreatedBy(), r.Restore.CreatedAt())...))
	}
	if r.Snapshot != nil {
		fmt.Println(summaryRow("Snapshot", r.Snapshot.Name(), r.Snapshot.ID(),
			auditExtras(r.Snapshot.CreatedBy(), r.Snapshot.CreatedAt())...))
	}

	// Container Registry
	if r.ContainerRegistry != nil {
		fmt.Println(summaryRow("Container Registry", r.ContainerRegistry.Name(), r.ContainerRegistry.ContainerRegistryID(),
			auditExtras(r.ContainerRegistry.CreatedBy(), r.ContainerRegistry.CreatedAt())...))
	}

	// Compute & schedule
	if r.CloudServer != nil {
		fmt.Println(summaryRow("Cloud Server", r.CloudServer.Name(), r.CloudServer.CloudServerID(),
			auditExtras(r.CloudServer.CreatedBy(), r.CloudServer.CreatedAt())...))
	}
	if r.KaaS != nil {
		fmt.Println(summaryRow("KaaS Cluster", r.KaaS.Name(), r.KaaS.KaaSID(),
			auditExtras(r.KaaS.CreatedBy(), r.KaaS.CreatedAt())...))
	}
	if r.JobRecurring != nil {
		fmt.Println(summaryRow("Job (Recurring)", r.JobRecurring.Name(), r.JobRecurring.JobID(),
			auditExtras(r.JobRecurring.CreatedBy(), r.JobRecurring.CreatedAt())...))
	}
	if r.JobOneShot != nil {
		fmt.Println(summaryRow("Job (One-Shot)", r.JobOneShot.Name(), r.JobOneShot.JobID(),
			auditExtras(r.JobOneShot.CreatedBy(), r.JobOneShot.CreatedAt())...))
	}

	// Database stack
	if r.DBaaS != nil {
		fmt.Println(summaryRow("DBaaS", r.DBaaS.Name(), r.DBaaS.DBaaSID(),
			auditExtras(r.DBaaS.CreatedBy(), r.DBaaS.CreatedAt())...))
	}
	if r.Database != nil && r.Database.ID() != "" {
		fmt.Println(summaryRow("DBaaS Database", r.Database.Name(), r.Database.ID(),
			auditExtras(r.Database.CreatedBy(), r.Database.CreatedAt())...))
	}
	if r.DBaaSUser != nil && r.DBaaSUser.ID() != "" {
		fmt.Println(summaryRow("DBaaS User", r.DBaaSUser.Username(), r.DBaaSUser.ID(),
			auditExtras(r.DBaaSUser.CreatedBy(), r.DBaaSUser.CreatedAt())...))
	}
	if r.Grant != nil && r.Grant.Username() != "" {
		grantDesc := r.Grant.Username() + " on " + r.Grant.DatabaseName() + " (" + r.Grant.RoleName() + ")"
		fmt.Println(summaryRow("DBaaS Grant", grantDesc, r.Grant.ID(),
			auditExtras(r.Grant.CreatedBy(), r.Grant.CreatedAt())...))
	}
}

// deleteAllResources deletes all resources in strict inverse of the creation order
// defined in orchestrator_create.go. Each resource waits until fully gone (HTTP 404)
// before the next dependent delete proceeds — this prevents the API from rejecting
// deletes of resources still referenced by a peer in "Deleting" state.
func deleteAllResources(ctx context.Context, arubaClient aruba.Client, resources *ResourceCollection) {
	fmt.Println("\n=== Deleting Resources ===")

	// Phase 1/8: VPN stack (inverse of create Phase 7: VPNRoute → VPNTunnel).
	printPhase(1, 8, "VPN stack")

	if resources.VPNRoute != nil && resources.VPNRoute.VPNRouteID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteVPNRoute(c, arubaClient, resources.VPNRoute) })
	}
	if resources.VPNTunnel != nil && resources.VPNTunnel.VPNTunnelID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteVPNTunnel(c, arubaClient, resources.VPNTunnel) })
	}

	// Phase 2/8: KMS stack (inverse of create Phase 8: KMS → KMSKey → KMIP).
	printPhase(2, 8, "KMS stack")

	if resources.Kmip != nil && resources.Kmip.KmipID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKmip(c, arubaClient, resources.Kmip) })
	}
	if resources.KMSKey != nil && resources.KMSKey.KeyID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKMSKey(c, arubaClient, resources.KMSKey) })
	}
	if resources.KMS != nil && resources.KMS.KMSID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteKMS(c, arubaClient, resources.KMS) })
	}

	// Phase 3/8: Backup & restore (inverse of create Phase 6: Restore → Backup).
	// RestoreTargetStorage is a BlockStorage; it falls into Phase 7 with all other volumes.
	printPhase(3, 8, "Backup & restore")

	if resources.Restore != nil && resources.Restore.RestoreID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteRestore(c, arubaClient, resources.Restore) })
	}
	if resources.Backup != nil && resources.Backup.ID() != "" {
		withDeleteDeadline(ctx, func(c context.Context) { deleteStorageBackup(c, arubaClient, resources.Backup) })
	}

	// Phase 4/8: Compute & container platforms (inverse of create Phase 5:
	// ContainerRegistry → Jobs → CloudServer → KaaS).
	printPhase(4, 8, "Compute & container platforms")

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

	// Phase 5/8: Database stack (inverse of create Phase 4:
	// Grant → DBaaSUser → Database → DBaaS).
	printPhase(5, 8, "Database stack")

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

	// Phase 6/8: VPC-scoped network (inverse of create Phase 3:
	// KeyPair → Rules (egress, ingress) → SecurityGroup → Subnets (Basic, Advanced)).
	printPhase(6, 8, "VPC-scoped network")

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

	// Phase 7/8: Independent network & storage primitives (inverse of create Phase 2:
	// VPC → Snapshot → BlockStorages → ElasticIPs).
	printPhase(7, 8, "Independent network & storage primitives")

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

	// Phase 8/8: Account & isolation (inverse of create Phase 1).
	printPhase(8, 8, "Account & isolation")

	if resources.Project != nil {
		withDeleteDeadline(ctx, func(c context.Context) { deleteProject(c, arubaClient, resources.Project) })
	}

	fmt.Println("\n=== Delete Complete ===")
}
