package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createKaaS provisions a Kubernetes-as-a-Service cluster with all dependencies and waits until Ready.
func createKaaS(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, vpc *aruba.VPC, subnet *aruba.Subnet) *aruba.KaaS {
	fmt.Println("--- KaaS (Kubernetes) ---")

	if err := waitForDependencies(ctx, "KaaS", map[string]waitFunc{
		"VPC":    vpc.WaitUntilActive,
		"Subnet": subnet.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	kaasSG := aruba.NewSecurityGroup().WithName(resourceName(NameKaaSSecurityGroup))
	k := aruba.NewKaaS().
		IntoProject(proj).
		WithName(resourceName(NameKaaS)).
		AddTag("kubernetes").
		AddTag("container").
		InRegion(defaultRegion).
		WithVPC(vpc).
		WithSubnet(subnet).
		WithSecurityGroup(kaasSG).
		WithNodeCIDR("172.16.0.0/16", resourceName(NameKaaSNodeCIDR)).
		WithKubernetesVersion("1.33.2").
		WithPodCIDR("10.0.3.0/24").
		WithHA(true).
		WithBillingPeriod("Hour").
		AddNodePool(aruba.NewNodePool().
			Named(resourceName(NameNodePool)).
			WithCount(2).
			WithAutoscaling(1, 5).
			OfInstance("K2A4").
			InZone(defaultZone))

	result, err := arubaClient.FromContainer().KaaS().Create(ctx, k)
	if err != nil {
		log.Fatalf("Error creating KaaS cluster: %s", formatErr(err))
		return nil
	}

	fmt.Printf("✓ Created KaaS cluster: %s (K8s: %s)\n",
		result.Name(),
		result.KubernetesVersion())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("KaaS %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

// updateKaaS applies name, tag, storage-quota, and node-pool changes to the cluster.
func updateKaaS(ctx context.Context, arubaClient aruba.Client, k *aruba.KaaS) {
	fmt.Println("--- Updating KaaS Cluster ---")

	// Mutate only the fields exposed by KaaSUpdateRequest.
	// Networking URIs and CIDRs are immutable after creation.
	k.WithName(updatedName(k.Name())).
		ReplaceTags("kubernetes", "container", "updated").
		WithMaxStorageQuotaGB(100).
		WithBillingPeriod("Hour").
		WithHA(true).
		AddNodePool(aruba.NewNodePool().
			Named(resourceName(NameNodePool)).
			WithCount(5).
			WithAutoscaling(1, 5).
			OfInstance("K2A4").
			InZone(defaultZone))

	result, err := arubaClient.FromContainer().KaaS().Update(ctx, k)
	if err != nil {
		log.Printf("Error updating KaaS cluster: %s", formatErr(err))
		return
	}

	fmt.Printf("✓ Updated KaaS cluster: %s (K8s: %s)\n", result.Name(), result.KubernetesVersion())
}

// deleteKaaS tears down the KaaS cluster.
func deleteKaaS(ctx context.Context, arubaClient aruba.Client, k *aruba.KaaS) {
	fmt.Println("--- Deleting KaaS Cluster ---")

	if err := arubaClient.FromContainer().KaaS().Delete(ctx, k); err != nil {
		log.Printf("Error deleting KaaS cluster: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted KaaS cluster: %s\n", k.KaaSID())
}
