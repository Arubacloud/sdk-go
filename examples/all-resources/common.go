package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// ---------------------------------------------------------------------------
// Name constants, vars, and helpers
// ---------------------------------------------------------------------------

// Fixed (resource-specific) parts of every name produced by createAllResources.
// Every constant is unique within this set so that two resources never share
// a generated name when the same prefix and suffix are used.
const (
	NameProject                   = "project"
	NameElasticIPCS               = "cs-eip"
	NameElasticIPDBaaS            = "dbaas-eip"
	NameElasticIPCR               = "cr-eip"
	NameElasticIPVPN              = "vpn-eip"
	NameBlockStorageCS            = "cs-bs"
	NameBlockStorageCR            = "cr-bs"
	NameBlockStorageRestoreTarget = "restore-target-bs"
	NameSnapshot                  = "snapshot"
	NameVPC                       = "vpc"
	NameSubnetAdvanced            = "subnet-advanced"
	NameSubnetBasic               = "subnet-basic"
	NameSecurityGroup             = "sg"
	NameSGRuleSSH                 = "rule-ssh"
	NameSGRuleHTTP                = "rule-http"
	NameSGRuleHTTPS               = "rule-https"
	NameSGRuleMySQL               = "rule-mysql"
	NameSGRuleEgress              = "rule-egress"
	NameKeyPair                   = "keypair"
	NameDBaaS                     = "dbaas"
	NameKaaS                      = "kaas"
	NameKaaSSecurityGroup         = "kaas-sg"
	NameKaaSNodeCIDR              = "kaas-node-cidr"
	NameNodePool                  = "pool"
	NameCloudServer               = "cs"
	NameContainerRegistry         = "cr"
	NameStorageBackup             = "backup"
	NameStorageRestore            = "restore"
	NameKMS                       = "kms"
	NameKMSKey                    = "kms-key"
	NameKmip                      = "kmip"
	NameDatabase                  = "testdb"
	NameDBaaSUser                 = "restapi"
	NameGrant                     = "grant"
	NameJobRecurring              = "job-recurring"
	NameJobOneShot                = "job-oneshot"
	NameVPNTunnel                 = "vpn-tunnel"
	NameVPNRoute                  = "vpn-route"
)

// namePrefix and nameSuffix are set once in main() before any create/update
// flow runs.
var (
	namePrefix string
	nameSuffix string
)

// resourceName composes the canonical name for a resource as:
//
//	<prefix>-<fixedPart>-<suffix>
//
// This is the only place in examples/all-resources that should build resource names.
func resourceName(fixedPart string) string {
	return fmt.Sprintf("%s-%s-%s", namePrefix, fixedPart, nameSuffix)
}

// updatedName appends "-updated" to current idempotently so that re-running
// update mode on an already-updated resource leaves the name unchanged
// instead of producing "name-updated-updated".
func updatedName(current string) string {
	const sfx = "-updated"
	if strings.HasSuffix(current, sfx) {
		return current
	}
	return current + sfx
}

// generateRandomSuffix returns 8 lower-case hex characters drawn from
// crypto/rand, giving ~4 billion possible suffixes per prefix.
func generateRandomSuffix() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Failed to generate random suffix: %v", err)
	}
	return hex.EncodeToString(b)
}

// ---------------------------------------------------------------------------
// Client options
// ---------------------------------------------------------------------------

// buildClientOptions returns DefaultOptions with optional verbose logging.
func buildClientOptions(clientID, clientSecret string, debug bool) *aruba.Options {
	opts := aruba.DefaultOptions(clientID, clientSecret)
	if debug {
		opts = opts.WithNativeLogger()
	}
	return opts
}

// ---------------------------------------------------------------------------
// ResourceCollection
// ---------------------------------------------------------------------------

// ResourceCollection holds all created resources.
type ResourceCollection struct {
	Project                  *aruba.Project
	CloudServerEIP           *aruba.ElasticIP    // used exclusively by CloudServer
	DBaaSEIP                 *aruba.ElasticIP    // used exclusively by DBaaS
	ContainerRegistryEIP     *aruba.ElasticIP    // used exclusively by ContainerRegistry
	VPNTunnelEIP             *aruba.ElasticIP    // used exclusively by VPN Tunnel ipConfigurations
	CloudServerBlockStorage  *aruba.BlockStorage // boot volume for CloudServer
	ContainerRegistryStorage *aruba.BlockStorage // storage for ContainerRegistry
	RestoreTargetStorage     *aruba.BlockStorage // dedicated unattached volume used as Restore destination
	Snapshot                 *aruba.Snapshot
	Backup                   *aruba.StorageBackup
	Restore                  *aruba.StorageRestore
	ContainerRegistry        *aruba.ContainerRegistry
	VPC                      *aruba.VPC
	SubnetAdvanced           *aruba.Subnet
	SubnetBasic              *aruba.Subnet
	SecurityGroup            *aruba.SecurityGroup
	SecurityRulesIngress     []*aruba.SecurityRule
	SecurityRuleEgress       *aruba.SecurityRule
	KeyPair                  *aruba.KeyPair
	DBaaS                    *aruba.DBaaS
	KaaS                     *aruba.KaaS
	CloudServer              *aruba.CloudServer
	KMS                      *aruba.KMS
	KMSKey                   *aruba.Key
	Kmip                     *aruba.Kmip
	KmipCert                 *aruba.KmipCertificate
	Database                 *aruba.Database
	DBaaSUser                *aruba.User
	Grant                    *aruba.Grant
	JobRecurring             *aruba.Job
	JobOneShot               *aruba.Job
	VPNTunnel                *aruba.VPNTunnel
	VPNRoute                 *aruba.VPNRoute

	// Populated by fetchAllResources for the delete flow. Contains every
	// BlockStorage and ElasticIP in the project so the pre-delete inventory
	// is complete before the user confirms deletion.
	BlockStorages []*aruba.BlockStorage
	ElasticIPs    []*aruba.ElasticIP
}

// ---------------------------------------------------------------------------
// Wait helpers
// ---------------------------------------------------------------------------

// longWaitOpts is the wait-option set for resources whose Ready transition
// routinely exceeds the SDK default (DBaaS, ContainerRegistry).
var longWaitOpts = []aruba.WaitOption{
	aruba.WithTimeout(40 * time.Minute),
	aruba.WithRetries(240),
}

// deleteOpTimeout caps each individual Delete + waitUntilGone in the delete
// orchestrator. The SDK's pkg/async.DefaultWaitFor already enforces a 10-min
// internal ceiling per wait; this cap is the orchestrator-level guard so a
// single slow resource (or one that never returns 404) cannot consume budget
// that subsequent deletes need.
const deleteOpTimeout = 12 * time.Minute

// withDeleteDeadline derives a child ctx with deleteOpTimeout from parent and
// runs fn against it. Parent cancellation (Ctrl-C, parent timeout) still
// propagates through the child, so this only bounds the upper limit — it
// does not extend the parent budget.
func withDeleteDeadline(parent context.Context, fn func(context.Context)) {
	ctx, cancel := context.WithTimeout(parent, deleteOpTimeout)
	defer cancel()
	fn(ctx)
}

type waitFunc func(context.Context, ...aruba.WaitOption) error

// stateReporter is satisfied by every statusMixin-backed wrapper.
// It lets wait helpers print the final state and failure reason on error.
type stateReporter interface {
	State() aruba.State
	FailureReason() string
}

// depEntry pairs a wait function with its optional state reporter so
// dependency helpers can emit state/reason detail on failure.
type depEntry struct {
	wait     waitFunc
	reporter stateReporter // may be nil
}

// dep constructs a depEntry. The reporter is typically the resource itself.
func dep(r stateReporter, w waitFunc) depEntry { return depEntry{wait: w, reporter: r} }

// waitForDependencies blocks until every entry in deps reaches its ready state.
func waitForDependencies(ctx context.Context, resourceLabel string, deps map[string]depEntry) error {
	fmt.Printf("⏳ Waiting for %s dependencies...\n", resourceLabel)
	for label, d := range deps {
		if err := d.wait(ctx); err != nil {
			stateDetail := ""
			if d.reporter != nil {
				stateDetail = fmt.Sprintf(" (state=%s reason=%q)", d.reporter.State(), d.reporter.FailureReason())
			}
			return fmt.Errorf("%s dep %s: %w%s", resourceLabel, label, err, stateDetail)
		}
		fmt.Printf("   ✓ %s ready\n", label)
	}
	fmt.Printf("✓ %s dependencies are READY\n", resourceLabel)
	return nil
}

// waitPostDependencies waits for downstream effects after a resource is created
// (e.g., an Elastic IP or Block Storage transitioning to Used state).
func waitPostDependencies(ctx context.Context, resourceLabel string, deps map[string]depEntry) {
	fmt.Printf("⏳ Waiting for %s post-dependencies...\n", resourceLabel)
	for label, d := range deps {
		if err := d.wait(ctx); err != nil {
			stateDetail := ""
			if d.reporter != nil {
				stateDetail = fmt.Sprintf(" (state=%s reason=%q)", d.reporter.State(), d.reporter.FailureReason())
			}
			log.Printf("%s post-dep %s: %v%s", resourceLabel, label, err, stateDetail)
			continue
		}
		fmt.Printf("   ✓ %s ready\n", label)
	}
	fmt.Printf("✓ %s post-dependencies confirmed\n", resourceLabel)
}

// waitUntilSelfReady wraps a resource's WaitUntilReady call with entry and
// success lifecycle messages. Logs entry, success, and error in one place so
// every resource_*.go create func collapses its boilerplate to a single call.
func waitUntilSelfReady(ctx context.Context, pretty, name string, reporter stateReporter, wait waitFunc, opts ...aruba.WaitOption) {
	fmt.Printf("⏳ Waiting for %s %s to become Ready...\n", pretty, name)
	if err := wait(ctx, opts...); err != nil {
		printSelfWaitError(pretty, name, reporter, err)
		return
	}
	fmt.Printf("✓ %s %s is Ready\n", pretty, name)
}

// waitAllReady blocks until every polling-aware resource in the collection reaches
// its terminal-success state. Logs but does not abort on per-resource failure so the
// summary still prints what succeeded.
func waitAllReady(ctx context.Context, r *ResourceCollection) {
	fmt.Println("\n=== Final readiness check ===")
	type waiter struct {
		label    string
		wait     waitFunc
		reporter stateReporter // may be nil
	}
	var ws []waiter
	if r.VPC != nil {
		ws = append(ws, waiter{"VPC", r.VPC.WaitUntilReady, r.VPC})
	}
	if r.SubnetAdvanced != nil {
		ws = append(ws, waiter{"Subnet (Advanced)", r.SubnetAdvanced.WaitUntilReady, r.SubnetAdvanced})
	}
	if r.SubnetBasic != nil {
		ws = append(ws, waiter{"Subnet (Basic)", r.SubnetBasic.WaitUntilReady, r.SubnetBasic})
	}
	if r.SecurityGroup != nil {
		ws = append(ws, waiter{"SecurityGroup", r.SecurityGroup.WaitUntilReady, r.SecurityGroup})
	}
	for _, rule := range r.SecurityRulesIngress {
		if rule != nil {
			label := "SecurityRule:" + rule.Name()
			ws = append(ws, waiter{label, rule.WaitUntilReady, rule})
		}
	}
	if r.SecurityRuleEgress != nil {
		ws = append(ws, waiter{"SecurityRuleEgress", r.SecurityRuleEgress.WaitUntilReady, r.SecurityRuleEgress})
	}
	if r.KeyPair != nil {
		ws = append(ws, waiter{"KeyPair", r.KeyPair.WaitUntilReady, r.KeyPair})
	}
	if r.CloudServerEIP != nil {
		ws = append(ws, waiter{"CloudServer EIP", r.CloudServerEIP.WaitUntilReady, r.CloudServerEIP})
	}
	if r.DBaaSEIP != nil {
		ws = append(ws, waiter{"DBaaS EIP", r.DBaaSEIP.WaitUntilReady, r.DBaaSEIP})
	}
	if r.ContainerRegistryEIP != nil {
		ws = append(ws, waiter{"CR EIP", r.ContainerRegistryEIP.WaitUntilReady, r.ContainerRegistryEIP})
	}
	if r.CloudServerBlockStorage != nil {
		ws = append(ws, waiter{"CloudServer BS", r.CloudServerBlockStorage.WaitUntilReady, r.CloudServerBlockStorage})
	}
	if r.ContainerRegistryStorage != nil {
		ws = append(ws, waiter{"CR BS", r.ContainerRegistryStorage.WaitUntilReady, r.ContainerRegistryStorage})
	}
	if r.RestoreTargetStorage != nil {
		ws = append(ws, waiter{"Restore-target BS", r.RestoreTargetStorage.WaitUntilReady, r.RestoreTargetStorage})
	}
	if r.DBaaS != nil {
		ws = append(ws, waiter{"DBaaS", func(ctx context.Context, _ ...aruba.WaitOption) error {
			return r.DBaaS.WaitUntilReady(ctx, longWaitOpts...)
		}, r.DBaaS})
	}
	if r.KaaS != nil {
		ws = append(ws, waiter{"KaaS", r.KaaS.WaitUntilReady, r.KaaS})
	}
	if r.CloudServer != nil {
		ws = append(ws, waiter{"CloudServer", r.CloudServer.WaitUntilReady, r.CloudServer})
	}
	if r.ContainerRegistry != nil {
		ws = append(ws, waiter{"ContainerRegistry", func(ctx context.Context, _ ...aruba.WaitOption) error {
			return r.ContainerRegistry.WaitUntilReady(ctx, longWaitOpts...)
		}, r.ContainerRegistry})
	}
	if r.Backup != nil {
		ws = append(ws, waiter{"StorageBackup", r.Backup.WaitUntilReady, r.Backup})
	}
	if r.Restore != nil {
		ws = append(ws, waiter{"StorageRestore", r.Restore.WaitUntilReady, r.Restore})
	}
	if r.KMS != nil {
		ws = append(ws, waiter{"KMS", r.KMS.WaitUntilReady, r.KMS})
	}
	if r.Kmip != nil {
		ws = append(ws, waiter{"KMIP", r.Kmip.WaitUntilReady, nil})
	}

	for _, w := range ws {
		if err := w.wait(ctx); err != nil {
			detail := err.Error()
			if w.reporter != nil && w.reporter.State() != "" {
				detail = fmt.Sprintf("%s (state=%s reason=%q)", detail, w.reporter.State(), w.reporter.FailureReason())
			}
			log.Printf("✗ %s not Ready: %s", w.label, detail)
		} else {
			fmt.Printf("✓ %s Ready\n", w.label)
		}
	}
}

// waitUntilGone blocks until the resource's WaitUntilGone reports it is fully
// deleted (Get returns HTTP 404), wrapping the call with lifecycle messages.
func waitUntilGone(ctx context.Context, label string, wait waitFunc) {
	fmt.Printf("⏳ Waiting for %s to fully terminate...\n", label)
	if err := wait(ctx); err != nil {
		log.Printf("Wait for %s to terminate: %v", label, err)
		return
	}
	fmt.Printf("✓ %s is fully gone\n", label)
}

// ---------------------------------------------------------------------------
// Error formatting
// ---------------------------------------------------------------------------

// formatErr returns a human-friendly error string. If err wraps an
// *aruba.HTTPError, the status code and API title are surfaced when available;
// otherwise err.Error() is returned verbatim.
func formatErr(err error) string {
	var httpErr *aruba.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.ErrResp != nil && httpErr.ErrResp.Title != nil {
			return fmt.Sprintf("HTTP %d: %s", httpErr.StatusCode, *httpErr.ErrResp.Title)
		}
		return fmt.Sprintf("HTTP %d: %s", httpErr.StatusCode, string(httpErr.Body))
	}
	return err.Error()
}

// ---------------------------------------------------------------------------
// Output helpers
// ---------------------------------------------------------------------------

// printBanner emits a resource creation banner: `--- {pretty} ---` or
// `--- {pretty} ({qualifier}) ---` when qualifier is non-empty.
func printBanner(pretty, qualifier string) {
	if qualifier == "" {
		fmt.Printf("--- %s ---\n", pretty)
	} else {
		fmt.Printf("--- %s (%s) ---\n", pretty, qualifier)
	}
}

// printPhase emits a numbered phase header for the create orchestrator.
func printPhase(n, total int, title string) {
	fmt.Printf("\n--- Phase %d/%d: %s ---\n", n, total, title)
}

// printCreated emits a creation success line. Extras are appended as
// comma-separated key=value pairs after the ID.
func printCreated(pretty, name, id string, extras ...string) {
	if len(extras) == 0 {
		fmt.Printf("✓ Created %s: %s (ID: %s)\n", pretty, name, id)
	} else {
		fmt.Printf("✓ Created %s: %s (ID: %s, %s)\n", pretty, name, id, strings.Join(extras, ", "))
	}
}

// printCreateError logs a create failure via log.Printf.
func printCreateError(pretty string, err error) {
	log.Printf("✗ Failed to create %s: %s", pretty, formatErr(err))
}

// printDepWaitError logs a pre-create dependency wait failure.
func printDepWaitError(pretty string, err error) {
	log.Printf("✗ %s dependency wait failed: %v", pretty, err)
}

// describeWaitFailure rewrites the wait-error string so the printed reason is
// actionable. statusMixin.WaitUntilStates already returns rich messages of the
// form `resource entered terminal error state "Failed" (targets [...])` —
// passed through verbatim. A bare `context deadline exceeded` (returned when
// the SDK's WaitFor hits its retry/timeout budget while the resource is still
// transitioning) gets an extra note so the reader knows the resource was still
// mid-flight, not failed.
func describeWaitFailure(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return err.Error() + " (timed out — wait budget exhausted before a terminal state was observed)"
	}
	return err.Error()
}

// printSelfWaitError logs a post-create self-readiness wait failure with state/reason.
func printSelfWaitError(pretty, name string, reporter stateReporter, err error) {
	detail := describeWaitFailure(err)
	if reporter != nil && reporter.State() != "" {
		detail = fmt.Sprintf("%s (state=%s reason=%q)", detail, reporter.State(), reporter.FailureReason())
	}
	log.Printf("✗ %s %s did not become Ready: %s", pretty, name, detail)
}

// printDeleteBanner emits a delete section header: `--- Deleting {pretty} ---`.
func printDeleteBanner(pretty string) {
	fmt.Printf("--- Deleting %s ---\n", pretty)
}

// printDeleteSubmitted emits the "delete accepted" line.
func printDeleteSubmitted(pretty, idOrName string) {
	fmt.Printf("→ Submitted delete for %s: %s\n", pretty, idOrName)
}

// printDeleteError logs a delete-call failure via log.Printf.
func printDeleteError(pretty string, err error) {
	log.Printf("✗ Failed to delete %s: %s", pretty, formatErr(err))
}

// summaryRow formats a single resource summary line.
// When name and id are both empty the resource was not created.
func summaryRow(label, name, id string, extras ...string) string {
	if name == "" && id == "" {
		return fmt.Sprintf("- %s: <not created>", label)
	}
	tail := id
	if len(extras) > 0 {
		tail = id + ", " + strings.Join(extras, ", ")
	}
	return fmt.Sprintf("- %s: %s (ID: %s)", label, name, tail)
}

// eipRow returns a summary line for an ElasticIP (which may be nil).
func eipRow(label string, eip *aruba.ElasticIP) string {
	if eip != nil {
		return summaryRow(label, eip.Name(), eip.ID())
	}
	return summaryRow(label, "", "")
}

// bsRow returns a summary line for a BlockStorage (which may be nil).
func bsRow(label string, bs *aruba.BlockStorage) string {
	if bs != nil {
		return summaryRow(label, bs.Name(), bs.ID())
	}
	return summaryRow(label, "", "")
}

// printResourceSummary prints a summary of all created resources.
func printResourceSummary(resources *ResourceCollection) {
	fmt.Println("\n=== SDK Example Complete ===")
	fmt.Println("Resources created:")

	// Project
	if resources.Project != nil {
		fmt.Println(summaryRow("Project", resources.Project.Name(), resources.Project.ID()))
	} else {
		fmt.Println(summaryRow("Project", "", ""))
	}

	// Elastic IPs
	fmt.Println(eipRow("ElasticIP (Cloud Server)", resources.CloudServerEIP))
	fmt.Println(eipRow("ElasticIP (DBaaS)", resources.DBaaSEIP))
	fmt.Println(eipRow("ElasticIP (Container Registry)", resources.ContainerRegistryEIP))
	fmt.Println(eipRow("ElasticIP (VPN Tunnel)", resources.VPNTunnelEIP))

	// Block Storages
	fmt.Println(bsRow("BlockStorage (Cloud Server)", resources.CloudServerBlockStorage))
	fmt.Println(bsRow("BlockStorage (Container Registry)", resources.ContainerRegistryStorage))
	fmt.Println(bsRow("BlockStorage (Restore Target)", resources.RestoreTargetStorage))

	// Snapshot / Backup / Restore
	if resources.Snapshot != nil {
		fmt.Println(summaryRow("Snapshot", resources.Snapshot.Name(), resources.Snapshot.ID()))
	} else {
		fmt.Println(summaryRow("Snapshot", "", ""))
	}
	if resources.Backup != nil {
		fmt.Println(summaryRow("Storage Backup", resources.Backup.Name(), resources.Backup.BackupID()))
	} else {
		fmt.Println(summaryRow("Storage Backup", "", ""))
	}
	if resources.Restore != nil {
		fmt.Println(summaryRow("Storage Restore", resources.Restore.Name(), resources.Restore.RestoreID()))
	} else {
		fmt.Println(summaryRow("Storage Restore", "", ""))
	}

	// Network
	if resources.VPC != nil {
		fmt.Println(summaryRow("VPC", resources.VPC.Name(), resources.VPC.ID()))
	} else {
		fmt.Println(summaryRow("VPC", "", ""))
	}
	if resources.SubnetAdvanced != nil {
		fmt.Println(summaryRow("Subnet (Advanced)", resources.SubnetAdvanced.Name(), resources.SubnetAdvanced.ID()))
	} else {
		fmt.Println(summaryRow("Subnet (Advanced)", "", ""))
	}
	if resources.SubnetBasic != nil {
		fmt.Println(summaryRow("Subnet (Basic)", resources.SubnetBasic.Name(), resources.SubnetBasic.ID()))
	} else {
		fmt.Println(summaryRow("Subnet (Basic)", "", ""))
	}

	// Security
	if resources.SecurityGroup != nil {
		fmt.Println(summaryRow("Security Group", resources.SecurityGroup.Name(), resources.SecurityGroup.ID()))
	} else {
		fmt.Println(summaryRow("Security Group", "", ""))
	}
	for _, r := range resources.SecurityRulesIngress {
		if r != nil {
			fmt.Println(summaryRow("Security Rule (Ingress/"+r.Name()+")", r.Name(), r.ID()))
		}
	}
	if resources.SecurityRuleEgress != nil {
		fmt.Println(summaryRow("Security Rule (Egress)", resources.SecurityRuleEgress.Name(), resources.SecurityRuleEgress.ID()))
	} else {
		fmt.Println(summaryRow("Security Rule (Egress)", "", ""))
	}

	// SSH Key Pair
	if resources.KeyPair != nil {
		fmt.Println(summaryRow("SSH Key Pair", resources.KeyPair.Name(), resources.KeyPair.KeyPairID()))
	} else {
		fmt.Println(summaryRow("SSH Key Pair", "", ""))
	}

	// Database stack
	if resources.DBaaS != nil {
		fmt.Println(summaryRow("DBaaS", resources.DBaaS.Name(), resources.DBaaS.DBaaSID()))
	} else {
		fmt.Println(summaryRow("DBaaS", "", ""))
	}
	if resources.Database != nil && resources.Database.ID() != "" {
		fmt.Printf("- DBaaS Database: %s\n", resources.Database.Name())
	} else {
		fmt.Println("- DBaaS Database: <not created>")
	}
	if resources.DBaaSUser != nil && resources.DBaaSUser.ID() != "" {
		fmt.Printf("- DBaaS User: %s\n", resources.DBaaSUser.Username())
	} else {
		fmt.Println("- DBaaS User: <not created>")
	}
	if resources.Grant != nil && resources.Grant.Username() != "" {
		fmt.Printf("- DBaaS Grant: %s on %s (%s)\n",
			resources.Grant.Username(), resources.Grant.DatabaseName(), resources.Grant.RoleName())
	} else {
		fmt.Println("- DBaaS Grant: <not created>")
	}

	// Compute & container platforms
	if resources.KaaS != nil {
		fmt.Println(summaryRow("KaaS Cluster", resources.KaaS.Name(), resources.KaaS.KaaSID()))
	} else {
		fmt.Println(summaryRow("KaaS Cluster", "", ""))
	}
	if resources.CloudServer != nil {
		fmt.Println(summaryRow("Cloud Server", resources.CloudServer.Name(), resources.CloudServer.CloudServerID()))
	} else {
		fmt.Println(summaryRow("Cloud Server", "", ""))
	}
	if resources.ContainerRegistry != nil {
		fmt.Println(summaryRow("Container Registry", resources.ContainerRegistry.Name(), resources.ContainerRegistry.ContainerRegistryID()))
	} else {
		fmt.Println(summaryRow("Container Registry", "", ""))
	}

	// Schedule jobs
	if resources.JobRecurring != nil {
		fmt.Println(summaryRow("Job (Recurring)", resources.JobRecurring.Name(), resources.JobRecurring.JobID()))
	} else {
		fmt.Println(summaryRow("Job (Recurring)", "", ""))
	}
	if resources.JobOneShot != nil {
		fmt.Println(summaryRow("Job (One-Shot)", resources.JobOneShot.Name(), resources.JobOneShot.JobID()))
	} else {
		fmt.Println(summaryRow("Job (One-Shot)", "", ""))
	}

	// VPN stack
	if resources.VPNTunnel != nil {
		fmt.Println(summaryRow("VPN Tunnel", resources.VPNTunnel.Name(), resources.VPNTunnel.VPNTunnelID()))
	} else {
		fmt.Println(summaryRow("VPN Tunnel", "", ""))
	}
	if resources.VPNRoute != nil {
		fmt.Println(summaryRow("VPN Route", resources.VPNRoute.Name(), resources.VPNRoute.VPNRouteID(),
			"cloudSubnet="+resources.VPNRoute.CloudSubnetCIDR()))
	} else {
		fmt.Println(summaryRow("VPN Route", "", ""))
	}

	// KMS stack
	if resources.KMS != nil {
		fmt.Println(summaryRow("KMS Instance", resources.KMS.Name(), resources.KMS.KMSID()))
	} else {
		fmt.Println(summaryRow("KMS Instance", "", ""))
	}
	if resources.KMSKey != nil {
		fmt.Println(summaryRow("KMS Key", resources.KMSKey.Name(), resources.KMSKey.KeyID()))
	} else {
		fmt.Println(summaryRow("KMS Key", "", ""))
	}
	if resources.Kmip != nil {
		fmt.Println(summaryRow("KMIP Service", resources.Kmip.Name(), resources.Kmip.KmipID()))
	} else {
		fmt.Println(summaryRow("KMIP Service", "", ""))
	}
	if resources.KmipCert != nil {
		fmt.Println("- KMIP Certificate: downloaded")
	} else {
		fmt.Println("- KMIP Certificate: <not created>")
	}
}
