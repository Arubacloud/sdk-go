package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createSecurityGroup(ctx context.Context, arubaClient aruba.Client, vpc *aruba.VPC) *aruba.SecurityGroup {
	fmt.Println("\n--- Network: Security Group ---")

	if err := waitForDependencies(ctx, "Security Group", map[string]waitFunc{
		"VPC": vpc.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	sg := aruba.NewSecurityGroup().
		IntoVPC(vpc).
		WithName(resourceName(NameSecurityGroup)).
		AddTag("security").
		AddTag("network")

	created, err := arubaClient.FromNetwork().SecurityGroups().Create(ctx, sg)
	if err != nil {
		log.Printf("Error creating security group: %v", err)
		return nil
	}
	fmt.Printf("✓ Created Security Group: %s (ObjectID: %s)\n", created.Name(), created.ID())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("Security Group %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

func createSecurityGroupIngressRule(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup, name, tag string, protocol aruba.RuleProtocol, port string) *aruba.SecurityRule {
	fmt.Printf("\n--- Network: Security Group Rule (Ingress/%s) ---\n", name)

	if err := waitForDependencies(ctx, "Ingress Rule", map[string]waitFunc{
		"Security Group": sg.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	rule := aruba.NewSecurityRule().
		IntoSecurityGroup(sg).
		WithName(name).
		AddTag(tag).
		AddTag("ingress").
		InRegion("ITBG-Bergamo").
		WithDirection(aruba.RuleDirectionIngress).
		WithProtocol(protocol).
		WithPort(port).
		WithTargetCIDR("0.0.0.0/0")

	created, err := arubaClient.FromNetwork().SecurityGroupRules().Create(ctx, rule)
	if err != nil {
		log.Printf("Error creating ingress rule %s: %v", name, err)
		return nil
	}
	fmt.Printf("✓ Created Ingress Rule: %s (ID: %s)\n", created.Name(), created.ID())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("Ingress Rule %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

// createSecurityGroupEgressRule allows all outbound traffic from the security group.
// Without this, DBaaS and other resources cannot initiate outbound connections.
func createSecurityGroupEgressRule(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup) *aruba.SecurityRule {
	fmt.Println("\n--- Network: Security Group Rule (Egress) ---")

	if err := waitForDependencies(ctx, "Egress Rule", map[string]waitFunc{
		"Security Group": sg.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	rule := aruba.NewSecurityRule().
		IntoSecurityGroup(sg).
		WithName(resourceName(NameSGRuleEgress)).
		AddTag("egress").
		InRegion("ITBG-Bergamo").
		WithDirection(aruba.RuleDirectionEgress).
		WithProtocol(aruba.RuleProtocolANY).
		WithPort("*").
		WithTargetCIDR("0.0.0.0/0")

	created, err := arubaClient.FromNetwork().SecurityGroupRules().Create(ctx, rule)
	if err != nil {
		log.Printf("Error creating egress rule: %v", err)
		return nil
	}
	fmt.Printf("✓ Created Egress Rule: %s (ID: %s)\n", created.Name(), created.ID())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("Egress Rule %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

// deleteSecurityGroup deletes a security group and waits for the platform to
// confirm removal. VPC deletion fails with 400 if a security group is still
// in Deleting state, so we block here before the caller proceeds to deleteVPC.
func deleteSecurityGroup(ctx context.Context, arubaClient aruba.Client, sg *aruba.SecurityGroup) {
	fmt.Println("--- Deleting Security Group ---")

	err := arubaClient.FromNetwork().SecurityGroups().Delete(ctx, sg)
	if err != nil {
		log.Printf("Error deleting security group: %v", err)
		return
	}
	fmt.Printf("✓ Deleted security group: %s\n", sg.ID())
	waitUntilGone(ctx, "security group "+sg.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromNetwork().SecurityGroups().Get(ctx, sg)
		return err
	})
}

func deleteSecurityGroupRule(ctx context.Context, arubaClient aruba.Client, rule *aruba.SecurityRule) {
	fmt.Println("--- Deleting Security Group Rule ---")

	if err := arubaClient.FromNetwork().SecurityGroupRules().Delete(ctx, rule); err != nil {
		log.Printf("Error deleting security rule: %v", err)
		return
	}
	fmt.Printf("✓ Deleted security rule: %s\n", rule.ID())
}
