package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/Arubacloud/sdk-go/pkg/async"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ---------------------------------------------------------------------------
// Region / zone defaults
// ---------------------------------------------------------------------------

const (
	// defaultRegion is the region every example resource is created in.
	defaultRegion = "ITBG-Bergamo"
	// defaultZone is the zone within defaultRegion for zonal resources.
	defaultZone = "ITBG-1"
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
	NameNodePool                  = "default-pool"
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
}

// ---------------------------------------------------------------------------
// Wait helpers
// ---------------------------------------------------------------------------

// longWaitOpts is the wait-option set for resources whose Ready transition
// routinely exceeds the SDK default (DBaaS, ContainerRegistry).
var longWaitOpts = []aruba.WaitOption{
	aruba.WithTimeout(20 * time.Minute),
	aruba.WithRetries(120),
}

type waitFunc func(context.Context, ...aruba.WaitOption) error

// waitForDependencies blocks until every entry in deps reaches its ready state.
func waitForDependencies(ctx context.Context, resourceLabel string, deps map[string]waitFunc) error {
	fmt.Printf("⏳ Waiting for %s dependencies...\n", resourceLabel)
	for label, wait := range deps {
		if err := wait(ctx); err != nil {
			return fmt.Errorf("%s dep %s: %w", resourceLabel, label, err)
		}
	}
	return nil
}

// waitPostDependencies waits for downstream effects after a resource is created
// (e.g., an Elastic IP or Block Storage transitioning to Used state).
func waitPostDependencies(ctx context.Context, resourceLabel string, deps map[string]waitFunc) {
	for label, wait := range deps {
		if err := wait(ctx); err != nil {
			log.Printf("%s post-dep %s: %v", resourceLabel, label, err)
		}
	}
}

// waitAllReady blocks until every polling-aware resource in the collection reaches
// its terminal-success state. Logs but does not abort on per-resource failure so the
// summary still prints what succeeded.
func waitAllReady(ctx context.Context, r *ResourceCollection) {
	fmt.Println("\n=== Final readiness check ===")
	type waiter struct {
		label string
		wait  waitFunc
	}
	var ws []waiter
	if r.VPC != nil {
		ws = append(ws, waiter{"VPC", r.VPC.WaitUntilReady})
	}
	if r.SubnetAdvanced != nil {
		ws = append(ws, waiter{"Subnet (Advanced)", r.SubnetAdvanced.WaitUntilReady})
	}
	if r.SubnetBasic != nil {
		ws = append(ws, waiter{"Subnet (Basic)", r.SubnetBasic.WaitUntilReady})
	}
	if r.SecurityGroup != nil {
		ws = append(ws, waiter{"SecurityGroup", r.SecurityGroup.WaitUntilReady})
	}
	for _, rule := range r.SecurityRulesIngress {
		if rule != nil {
			label := "SecurityRule:" + rule.Name()
			ws = append(ws, waiter{label, rule.WaitUntilReady})
		}
	}
	if r.SecurityRuleEgress != nil {
		ws = append(ws, waiter{"SecurityRuleEgress", r.SecurityRuleEgress.WaitUntilReady})
	}
	if r.KeyPair != nil {
		ws = append(ws, waiter{"KeyPair", r.KeyPair.WaitUntilReady})
	}
	if r.CloudServerEIP != nil {
		ws = append(ws, waiter{"CloudServer EIP", r.CloudServerEIP.WaitUntilReady})
	}
	if r.DBaaSEIP != nil {
		ws = append(ws, waiter{"DBaaS EIP", r.DBaaSEIP.WaitUntilReady})
	}
	if r.ContainerRegistryEIP != nil {
		ws = append(ws, waiter{"CR EIP", r.ContainerRegistryEIP.WaitUntilReady})
	}
	if r.CloudServerBlockStorage != nil {
		ws = append(ws, waiter{"CloudServer BS", r.CloudServerBlockStorage.WaitUntilReady})
	}
	if r.ContainerRegistryStorage != nil {
		ws = append(ws, waiter{"CR BS", r.ContainerRegistryStorage.WaitUntilReady})
	}
	if r.RestoreTargetStorage != nil {
		ws = append(ws, waiter{"Restore-target BS", r.RestoreTargetStorage.WaitUntilReady})
	}
	if r.DBaaS != nil {
		ws = append(ws, waiter{"DBaaS", func(ctx context.Context, _ ...aruba.WaitOption) error {
			return r.DBaaS.WaitUntilReady(ctx, longWaitOpts...)
		}})
	}
	if r.KaaS != nil {
		ws = append(ws, waiter{"KaaS", r.KaaS.WaitUntilReady})
	}
	if r.CloudServer != nil {
		ws = append(ws, waiter{"CloudServer", r.CloudServer.WaitUntilReady})
	}
	if r.ContainerRegistry != nil {
		ws = append(ws, waiter{"ContainerRegistry", func(ctx context.Context, _ ...aruba.WaitOption) error {
			return r.ContainerRegistry.WaitUntilReady(ctx, longWaitOpts...)
		}})
	}
	if r.Backup != nil {
		ws = append(ws, waiter{"StorageBackup", r.Backup.WaitUntilReady})
	}
	if r.Restore != nil {
		ws = append(ws, waiter{"StorageRestore", r.Restore.WaitUntilReady})
	}
	if r.KMS != nil {
		ws = append(ws, waiter{"KMS", r.KMS.WaitUntilReady})
	}
	if r.Kmip != nil {
		ws = append(ws, waiter{"KMIP", r.Kmip.WaitUntilReady})
	}

	for _, w := range ws {
		if err := w.wait(ctx); err != nil {
			log.Printf("✗ %s not Ready: %v", w.label, err)
		} else {
			fmt.Printf("✓ %s Ready\n", w.label)
		}
	}
}

// waitUntilGone polls the given Get function until it returns HTTP 404, which is
// the platform's signal that an async delete has fully propagated. The retry,
// fixed-delay, and timeout machinery comes from the SDK's pkg/async.WaitFor —
// the call/check pair encodes the "404 means success" contract.
func waitUntilGone(ctx context.Context, label string, poll func(context.Context) error) {
	fmt.Printf("⏳ Waiting for %s to fully terminate...\n", label)
	const goneSentinel = "gone"
	fut := async.DefaultWaitFor(ctx,
		func(ctx context.Context) (*types.Response[string], error) {
			err := poll(ctx)
			if err == nil {
				// resource still exists — keep polling
				return &types.Response[string]{}, nil
			}
			var httpErr *aruba.HTTPError
			if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
				marker := goneSentinel
				return &types.Response[string]{Data: &marker}, nil
			}
			// transient — record for diagnostics but stay in the loop
			return nil, err
		},
		func(resp *types.Response[string]) (bool, error) {
			return resp != nil && resp.Data != nil && *resp.Data == goneSentinel, nil
		},
	)
	if _, err := fut.Await(ctx); err != nil {
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

// printResourceSummary prints a summary of all created resources.
func printResourceSummary(resources *ResourceCollection) {
	fmt.Println("\n=== SDK Example Complete ===")
	fmt.Println("Successfully created resources:")
	if resources.Project != nil {
		fmt.Println("- Project ID:", resources.Project.ID())
	}

	if resources.CloudServerEIP != nil && resources.CloudServerEIP.ID() != "" {
		fmt.Println("- Cloud Server ElasticIP ID:", resources.CloudServerEIP.ID())
	}
	if resources.DBaaSEIP != nil && resources.DBaaSEIP.ID() != "" {
		fmt.Println("- DBaaS ElasticIP ID:", resources.DBaaSEIP.ID())
	}
	if resources.ContainerRegistryEIP != nil && resources.ContainerRegistryEIP.ID() != "" {
		fmt.Println("- CR ElasticIP ID:", resources.ContainerRegistryEIP.ID())
	}

	if resources.CloudServerBlockStorage != nil && resources.CloudServerBlockStorage.ID() != "" {
		fmt.Println("- Cloud Server Block Storage ID:", resources.CloudServerBlockStorage.ID())
	}
	if resources.ContainerRegistryStorage != nil && resources.ContainerRegistryStorage.ID() != "" {
		fmt.Println("- CR Block Storage ID:", resources.ContainerRegistryStorage.ID())
	}

	if resources.Snapshot != nil && resources.Snapshot.ID() != "" {
		fmt.Printf("- Snapshot: %s (ID: %s)\n", resources.Snapshot.Name(), resources.Snapshot.ID())
	}

	if resources.VPC != nil && resources.VPC.ID() != "" {
		fmt.Println("- VPC ID:", resources.VPC.ID())
	}

	if resources.SubnetAdvanced != nil && resources.SubnetAdvanced.ID() != "" {
		fmt.Println("- Advanced Subnet ID:", resources.SubnetAdvanced.ID())
	}
	if resources.SubnetBasic != nil && resources.SubnetBasic.ID() != "" {
		fmt.Println("- Basic Subnet ID:", resources.SubnetBasic.ID())
	}

	if resources.SecurityGroup != nil && resources.SecurityGroup.ID() != "" {
		fmt.Println("- Security Group ID:", resources.SecurityGroup.ID())
	}

	for _, r := range resources.SecurityRulesIngress {
		if r != nil && r.ID() != "" {
			fmt.Printf("- Security Rule (Ingress/%s) ID: %s\n", r.Name(), r.ID())
		}
	}
	if resources.SecurityRuleEgress != nil && resources.SecurityRuleEgress.ID() != "" {
		fmt.Println("- Security Rule (Egress) ID:", resources.SecurityRuleEgress.ID())
	}

	if resources.KeyPair != nil && resources.KeyPair.KeyPairID() != "" {
		fmt.Printf("- SSH Key Pair: %s (ID: %s)\n", resources.KeyPair.Name(), resources.KeyPair.KeyPairID())
	}

	if resources.DBaaS != nil && resources.DBaaS.DBaaSID() != "" {
		fmt.Printf("- DBaaS: %s (ID: %s)\n", resources.DBaaS.Name(), resources.DBaaS.DBaaSID())
	}

	if resources.KaaS != nil && resources.KaaS.KaaSID() != "" {
		fmt.Println("- KaaS Cluster ID:", resources.KaaS.KaaSID())
	}

	if resources.CloudServer != nil && resources.CloudServer.CloudServerID() != "" {
		fmt.Printf("- Cloud Server: %s (ID: %s)\n", resources.CloudServer.Name(), resources.CloudServer.CloudServerID())
	}

	if resources.KMS != nil && resources.KMS.KMSID() != "" {
		fmt.Println("- KMS Instance ID:", resources.KMS.KMSID())
	}

	if resources.KMSKey != nil && resources.KMSKey.KeyID() != "" {
		fmt.Println("- KMS Key ID:", resources.KMSKey.KeyID())
	}

	if resources.Kmip != nil && resources.Kmip.KmipID() != "" {
		fmt.Println("- KMIP Service ID:", resources.Kmip.KmipID())
	}

	if resources.KmipCert != nil {
		fmt.Println("- KMIP Certificate: Downloaded successfully")
	}

	if resources.Database != nil && resources.Database.ID() != "" {
		fmt.Printf("- DBaaS Database: %s\n", resources.Database.Name())
	}
	if resources.DBaaSUser != nil && resources.DBaaSUser.ID() != "" {
		fmt.Printf("- DBaaS User: %s\n", resources.DBaaSUser.Username())
	}
	if resources.Grant != nil && resources.Grant.ID() != "" {
		fmt.Printf("- DBaaS Grant: %s on %s (%s)\n",
			resources.Grant.Username(), resources.Grant.DatabaseName(), resources.Grant.RoleName())
	}
	if resources.JobRecurring != nil && resources.JobRecurring.JobID() != "" {
		fmt.Println("- Recurring Job ID:", resources.JobRecurring.JobID())
	}
	if resources.JobOneShot != nil && resources.JobOneShot.JobID() != "" {
		fmt.Println("- OneShot Job ID:", resources.JobOneShot.JobID())
	}
}

// stringValue dereferences a *string and returns "" when nil.
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
