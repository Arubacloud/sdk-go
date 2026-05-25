package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createKaaS provisions a Kubernetes-as-a-Service cluster with all dependencies and waits until Ready.
func createKaaS(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, vpc *aruba.VPC, subnet *aruba.Subnet) *aruba.KaaS {
	printBanner("KaaS Cluster", "")

	if err := waitForDependencies(ctx, "KaaS Cluster", map[string]depEntry{
		"VPC":    dep(vpc, vpc.WaitUntilActive),
		"Subnet": dep(subnet, subnet.WaitUntilActive),
	}); err != nil {
		printDepWaitError("KaaS Cluster", err)
		return nil
	}

	kaasSG := aruba.NewSecurityGroup().
		Named(resourceName(NameKaaSSecurityGroup))

	k := aruba.NewKaaS().
		Named(resourceName(NameKaaS)).
		Tagged("kubernetes", "container").
		InProject(proj).
		InRegion(aruba.RegionITBGBergamo).
		WithKubernetesVersion(aruba.KubernetesVersion1341).
		WithPodCIDR("10.0.3.0/24").
		WithNodeCIDR("172.16.0.0/16", resourceName(NameKaaSNodeCIDR)).
		WithVPC(vpc).
		WithSubnet(subnet).
		WithSecurityGroup(kaasSG).
		WithNodePools(aruba.NewNodePool().
			OfInstance(aruba.NodePoolInstanceK2A4).
			Named(resourceName(NameNodePool)).
			InZone(aruba.ZoneITBG1).
			WithCount(2).
			WithAutoscaling(1, 5)).
		HighlyAvailable().
		BilledBy(aruba.BillingPeriodHour)

	result, err := arubaClient.FromContainer().KaaS().Create(ctx, k)
	if err != nil {
		printCreateError("KaaS Cluster", err)
		return nil
	}
	printCreated("KaaS Cluster", result.Name(), result.KaaSID())

	waitUntilSelfReady(ctx, "KaaS Cluster", result.Name(), result, result.WaitUntilReady)

	return result
}

// updateKaaS applies name, tag, storage-quota, and node-pool changes to the cluster.
func updateKaaS(ctx context.Context, arubaClient aruba.Client, k *aruba.KaaS) {
	fmt.Println("--- Updating KaaS Cluster ---")

	// Mutate only the fields exposed by KaaSUpdateRequest.
	// Networking URIs and CIDRs are immutable after creation.
	k.Named(updatedName(k.Name())).
		RetaggedAs("kubernetes", "container", "updated").
		WithMaxStorageQuotaGB(100).
		WithNodePools(aruba.NewNodePool().
			OfInstance(aruba.NodePoolInstanceK2A4).
			Named(resourceName(NameNodePool)).
			InZone(aruba.ZoneITBG1).
			WithCount(5).
			WithAutoscaling(1, 5)).
		HighlyAvailable().
		BilledBy(aruba.BillingPeriodHour)

	result, err := arubaClient.FromContainer().KaaS().Update(ctx, k)
	if err != nil {
		log.Printf("Error updating KaaS cluster: %s", formatErr(err))
		return
	}

	fmt.Printf("✓ Updated KaaS cluster: %s (K8s: %s)\n", result.Name(), result.KubernetesVersion())
}

// deleteKaaS tears down the KaaS cluster and waits until it is fully gone.
func deleteKaaS(ctx context.Context, arubaClient aruba.Client, k *aruba.KaaS) {
	printDeleteBanner("KaaS Cluster")
	if err := arubaClient.FromContainer().KaaS().Delete(ctx, k); err != nil {
		printDeleteError("KaaS Cluster", err)
		return
	}
	printDeleteSubmitted("KaaS Cluster", k.Name())
	waitUntilGone(ctx, "KaaS Cluster "+k.Name(), k.WaitUntilGone)
}
