package main

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createKeyPair uploads an SSH public key and waits until Ready.
func createKeyPair(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref) *aruba.KeyPair {
	fmt.Println("--- SSH Key Pair ---")

	sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA2No7At0tgHrcZTL0kGWyLLUqPKfOhD9hGdNV9PbJxhjOGNFxcwdQ9wCXsJ3RQaRHBuGIgVodDurrlqzxFK86yCHMgXT2YLHF0j9P4m9GDiCfOK6msbFb89p5xZExjwD2zK+w68r7iOKZeRB2yrznW5TD3KDemSPIQQIVcyLF+yxft49HWBTI3PVQ4rBVOBJ2PdC9SAOf7CYnptW24CRrC0h85szIDwMA+Kmasfl3YGzk4MxheHrTO8C40aXXpieJ9S2VQA4VJAMRyAboptIK0cKjBYrbt5YkEL0AlyBGPIu6MPYr5K/MHyDunDi9yc7VYRYRR0f46MBOSqMUiGPnMw=="

	kp := aruba.NewKeyPair().
		IntoProject(proj).
		WithName(resourceName(NameKeyPair)).
		AddTag("ssh-access").
		AddTag("ingress").
		InRegion(aruba.RegionITBGBergamo).
		WithPublicKey(sshPublicKey)

	kp, err := arubaClient.FromCompute().KeyPairs().Create(ctx, kp)
	if err != nil {
		printCreateError("SSH Key Pair", err)
		return nil
	}
	printCreated("SSH Key Pair", kp.Name(), kp.KeyPairID())
	if err := kp.WaitUntilReady(ctx); err != nil {
		printSelfWaitError("SSH Key Pair", kp.Name(), err)
	}

	return kp
}

// deleteKeyPair removes the SSH key pair and waits until it is fully gone.
func deleteKeyPair(ctx context.Context, arubaClient aruba.Client, kp *aruba.KeyPair) {
	printDeleteBanner("SSH Key Pair")
	if err := arubaClient.FromCompute().KeyPairs().Delete(ctx, kp); err != nil {
		printDeleteError("SSH Key Pair", err)
		return
	}
	printDeleteSubmitted("SSH Key Pair", kp.Name())
	waitUntilGone(ctx, "SSH Key Pair "+kp.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromCompute().KeyPairs().Get(ctx, kp)
		return err
	})
}
