package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createSecurityGroup provisions a security group inside the VPC and waits until Ready.
func createSecurityGroup(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.SecurityGroup {
	printBanner("Security Group", "")

	if err := waitForDependencies(ctx, "Security Group", map[string]depEntry{
		"VPC": dep(vpc, vpc.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Security Group", err)
		return nil
	}

	sg := aruba.NewSecurityGroup().
		Named(resourceName(NameSecurityGroup)).
		Tagged("security", "network").
		InVPC(vpc)

	created, err := arubaClient.FromNetwork().SecurityGroups().Create(ctx, sg)
	if err != nil {
		printCreateError("Security Group", err)
		return nil
	}
	printCreated("Security Group", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "Security Group", created.Name(), created, created.WaitUntilReady)

	return created
}

// createSecurityGroupIngressRule provisions an ingress rule on the security group.
func createSecurityGroupIngressRule(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup, name, tag string, protocol aruba.RuleProtocol, port string) *aruba.SecurityRule {
	fmt.Printf("--- Security Rule (Ingress/%s) ---\n", name)

	if err := waitForDependencies(ctx, "Security Rule (Ingress)", map[string]depEntry{
		"Security Group": dep(sg, sg.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Security Rule (Ingress)", err)
		return nil
	}

	rule := aruba.NewSecurityRule().
		Named(name).
		Tagged(tag, "ingress").
		InSecurityGroup(sg).
		InRegion(aruba.RegionITBGBergamo).
		WithDirection(aruba.RuleDirectionIngress).
		WithProtocol(protocol).
		WithPort(port).
		TargetingCIDR("0.0.0.0/0")

	created, err := arubaClient.FromNetwork().SecurityGroupRules().Create(ctx, rule)
	if err != nil {
		printCreateError("Security Rule (Ingress)", err)
		return nil
	}
	printCreated("Security Rule (Ingress)", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "Security Rule (Ingress)", created.Name(), created, created.WaitUntilReady)

	return created
}

// createSecurityGroupEgressRule allows all outbound traffic from the security group.
// Without this, DBaaS and other resources cannot initiate outbound connections.
func createSecurityGroupEgressRule(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup) *aruba.SecurityRule {
	printBanner("Security Rule", "Egress")

	if err := waitForDependencies(ctx, "Security Rule (Egress)", map[string]depEntry{
		"Security Group": dep(sg, sg.WaitUntilActive),
	}); err != nil {
		printDepWaitError("Security Rule (Egress)", err)
		return nil
	}

	rule := aruba.NewSecurityRule().
		Named(resourceName(NameSGRuleEgress)).
		Tagged("egress").
		InSecurityGroup(sg).
		InRegion(aruba.RegionITBGBergamo).
		WithDirection(aruba.RuleDirectionEgress).
		WithProtocol(aruba.RuleProtocolANY).
		TargetingCIDR("0.0.0.0/0")

	created, err := arubaClient.FromNetwork().SecurityGroupRules().Create(ctx, rule)
	if err != nil {
		printCreateError("Security Rule (Egress)", err)
		return nil
	}
	printCreated("Security Rule (Egress)", created.Name(), created.ID())

	waitUntilSelfReady(ctx, "Security Rule (Egress)", created.Name(), created, created.WaitUntilReady)

	return created
}

// deleteSecurityGroup deletes a security group and waits for the platform to
// confirm removal. VPC deletion fails with 400 if a security group is still
// in Deleting state, so we block here before the caller proceeds to deleteVPC.
func deleteSecurityGroup(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup) {
	printDeleteBanner("Security Group")
	if err := arubaClient.FromNetwork().SecurityGroups().Delete(ctx, sg); err != nil {
		printDeleteError("Security Group", err)
		return
	}
	printDeleteSubmitted("Security Group", sg.Name())
	waitUntilGone(ctx, "Security Group "+sg.Name(), sg.WaitUntilGone)
}

// deleteSecurityGroupRule removes the security group rule and waits until it is fully gone.
func deleteSecurityGroupRule(ctx context.Context, arubaClient aruba.Client, rule *aruba.SecurityRule) {
	printDeleteBanner("Security Group Rule")
	if err := arubaClient.FromNetwork().SecurityGroupRules().Delete(ctx, rule); err != nil {
		printDeleteError("Security Group Rule", err)
		return
	}
	printDeleteSubmitted("Security Group Rule", rule.Name())
	waitUntilGone(ctx, "Security Group Rule "+rule.Name(), rule.WaitUntilGone)
}
