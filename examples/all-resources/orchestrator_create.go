package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// runCreateExample demonstrates how to create all resources.
func runCreateExample(clientID, clientSecret string, debug bool) {
	arubaClient, err := aruba.NewClient(buildClientOptions(clientID, clientSecret, debug))
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	fmt.Println("=== SDK Create Example ===")

	resources := createAllResources(ctx, arubaClient)

	printResourceSummary(resources)
}

// createAllResources orchestrates the creation of all resources in dependency order.
func createAllResources(ctx context.Context, arubaClient aruba.Client) *ResourceCollection {
	resources := &ResourceCollection{}

	printPhase(1, 8, "Account & isolation")

	// 1. Create Project.
	resources.Project = createProject(ctx, arubaClient)

	printPhase(2, 8, "Independent network & storage primitives")

	// 2. Create Elastic IPs — one per consumer to avoid attach-state conflicts.
	resources.CloudServerEIP = createElasticIP(ctx, arubaClient, resources.Project, resourceName(NameElasticIPCS))
	resources.DBaaSEIP = createElasticIP(ctx, arubaClient, resources.Project, resourceName(NameElasticIPDBaaS))
	resources.ContainerRegistryEIP = createElasticIP(ctx, arubaClient, resources.Project, resourceName(NameElasticIPCR))

	// 3. Create Block Storage volumes — one per consumer.
	resources.CloudServerBlockStorage = createBlockStorage(ctx, arubaClient, resources.Project, resourceName(NameBlockStorageCS))
	resources.ContainerRegistryStorage = createBlockStorage(ctx, arubaClient, resources.Project, resourceName(NameBlockStorageCR))

	// 4. Snapshot the CloudServer boot volume before it is attached (pre-dep BS Ready included).
	resources.Snapshot = createSnapshot(ctx, arubaClient, resources.Project, resources.CloudServerBlockStorage)

	// 5. Create VPC (self-wait included).
	resources.VPC = createVPC(ctx, arubaClient, resources.Project)

	printPhase(3, 8, "VPC-scoped network")

	// 6. Create Subnets in VPC (pre-dep VPC Active + self-wait included).
	resources.SubnetAdvanced = createAdvancedSubnet(ctx, arubaClient, resources.VPC)
	resources.SubnetBasic = createBasicSubnet(ctx, arubaClient, resources.VPC)

	// 7. Create Security Group (pre-dep VPC Active + self-wait included).
	resources.SecurityGroup = createSecurityGroup(ctx, arubaClient, resources.VPC)

	// 8. Create Security Group Rules (pre-dep SG Active + self-wait included).
	for _, r := range []struct {
		name, tag string
		proto     aruba.RuleProtocol
		port      string
	}{
		{resourceName(NameSGRuleSSH), "ssh", aruba.RuleProtocolTCP, "22"},
		{resourceName(NameSGRuleHTTP), "http", aruba.RuleProtocolTCP, "80"},
		{resourceName(NameSGRuleHTTPS), "https", aruba.RuleProtocolTCP, "443"},
		{resourceName(NameSGRuleMySQL), "mysql", aruba.RuleProtocolTCP, "3306"},
	} {
		if rule := createSecurityGroupIngressRule(ctx, arubaClient, resources.SecurityGroup, r.name, r.tag, r.proto, r.port); rule != nil {
			resources.SecurityRulesIngress = append(resources.SecurityRulesIngress, rule)
		}
	}
	resources.SecurityRuleEgress = createSecurityGroupEgressRule(ctx, arubaClient, resources.SecurityGroup)

	// 9. Create SSH Key Pair (self-wait included).
	resources.KeyPair = createKeyPair(ctx, arubaClient, resources.Project)

	printPhase(4, 8, "Database stack")

	// 10. Create DBaaS (pre-dep waits + self-wait + EIP post-dep all included).
	resources.DBaaS = createDBaaS(ctx, arubaClient, resources.Project, resources.VPC, resources.SubnetBasic, resources.SecurityGroup, resources.DBaaSEIP)

	// 11. Create DBaaS Database, User, and Grant (pre-dep DBaaS Ready).
	if resources.DBaaS != nil {
		resources.Database = createDatabase(ctx, arubaClient, resources.DBaaS)
		resources.DBaaSUser = createDBaaSUser(ctx, arubaClient, resources.DBaaS)
		if resources.Database != nil && resources.DBaaSUser != nil {
			resources.Grant = createGrant(ctx, arubaClient, resources.Database, resources.DBaaSUser)
		}
	}

	printPhase(5, 8, "Compute & container platforms")

	// 12. Create KaaS (pre-dep waits + self-wait included).
	resources.KaaS = createKaaS(ctx, arubaClient, resources.Project, resources.VPC, resources.SubnetBasic)

	// 13. Create Cloud Server (pre-dep waits + self-wait + EIP/BS post-deps all included).
	resources.CloudServer = createCloudServer(ctx, arubaClient, resources)

	// 14. Create Schedule Jobs targeting the CloudServer.
	// Placed before ContainerRegistry so the jobs always run within the outer
	// create context, even if CR's wait consumes its full 20-minute budget.
	if resources.CloudServer != nil {
		resources.JobRecurring = createRecurringJob(ctx, arubaClient, resources.Project, resources.CloudServer)
		resources.JobOneShot = createOneShotJob(ctx, arubaClient, resources.Project, resources.CloudServer)
	}

	// 15. Create Container Registry (pre-dep waits + self-wait + EIP/BS post-deps all included).
	resources.ContainerRegistry = createContainerRegistry(ctx, arubaClient, resources)

	printPhase(6, 8, "Backup & restore")

	// 16. Create Storage Backup (pre-dep + self-wait included).
	resources.Backup = createStorageBackup(ctx, arubaClient, resources.Project, resources.CloudServerBlockStorage)

	// 17. Create a dedicated unattached BlockStorage as the Restore destination.
	//     Restoring into the CloudServer's boot volume is rejected by the API
	//     because that volume is in "InUse"/"Used" state and there is no Detach API.
	resources.RestoreTargetStorage = createBlockStorage(ctx, arubaClient,
		resources.Project, resourceName(NameBlockStorageRestoreTarget))

	// 18. Create Restore from Backup (pre-dep waits + self-wait included).
	resources.Restore = createRestore(ctx, arubaClient, resources.Backup, resources.RestoreTargetStorage)

	printPhase(7, 8, "VPN stack")

	// 19. Create VPN Tunnel (self-wait included).
	resources.VPNTunnel = createVPNTunnel(ctx, arubaClient, resources.Project)

	if resources.VPNTunnel != nil {
		// 20. Create VPN Route under the tunnel (self-wait + GET inspection for #308 included).
		resources.VPNRoute = createVPNRoute(ctx, arubaClient, resources.VPNTunnel)
	}

	printPhase(8, 8, "KMS stack")

	// 19. Create KMS Instance (self-wait included).
	resources.KMS = createKMS(ctx, arubaClient, resources.Project)

	if resources.KMS != nil && resources.KMS.KMSID() != "" {
		// 20. Create KMS Key (pre-dep KMS Active + post-dep KMS re-Active included).
		resources.KMSKey = createKMSKey(ctx, arubaClient, resources.KMS)

		// 21. Create KMIP Service (pre-dep KMS Active + self-wait included).
		resources.Kmip = createKmip(ctx, arubaClient, resources.KMS)

		// 22. Download KMIP Certificate (if KMIP was created).
		if resources.Kmip != nil && resources.Kmip.KmipID() != "" {
			resources.KmipCert = downloadKmipCertificate(ctx, arubaClient, resources.Kmip)
		}
	}

	// Final gate: block until every resource reaches its terminal-success state.
	waitAllReady(ctx, resources)

	return resources
}
