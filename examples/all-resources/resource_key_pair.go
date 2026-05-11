package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createKeyPair(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref) *aruba.KeyPair {
	fmt.Println("--- SSH Key Pair ---")

	sshPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA2No7At0tgHrcZTL0kGWyLLUqPKfOhD9hGdNV9PbJxhjOGNFxcwdQ9wCXsJ3RQaRHBuGIgVodDurrlqzxFK86yCHMgXT2YLHF0j9P4m9GDiCfOK6msbFb89p5xZExjwD2zK+w68r7iOKZeRB2yrznW5TD3KDemSPIQQIVcyLF+yxft49HWBTI3PVQ4rBVOBJ2PdC9SAOf7CYnptW24CRrC0h85szIDwMA+Kmasfl3YGzk4MxheHrTO8C40aXXpieJ9S2VQA4VJAMRyAboptIK0cKjBYrbt5YkEL0AlyBGPIu6MPYr5K/MHyDunDi9yc7VYRYRR0f46MBOSqMUiGPnMw=="

	kp := aruba.NewKeyPair().
		IntoProject(proj).
		WithName(resourceName(NameKeyPair)).
		AddTag("ssh-access").
		AddTag("ingress").
		InRegion("ITBG-Bergamo").
		WithPublicKey(sshPublicKey)

	kp, err := arubaClient.FromCompute().KeyPairs().Create(ctx, kp)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to create SSH key pair - Status: %d, Error: %s", httpErr.StatusCode, httpErr.Error())
		} else {
			log.Printf("Error creating SSH key pair: %v", err)
		}
	} else {
		fmt.Printf("✓ Created SSH Key Pair: %s\n", kp.Name())
		if err := kp.WaitUntilReady(ctx); err != nil {
			log.Printf("SSH Key Pair %s did not become Ready: %v", kp.Name(), err)
		}
	}

	return kp
}

func deleteKeyPair(ctx context.Context, arubaClient aruba.Client, kp *aruba.KeyPair) {
	fmt.Println("--- Deleting SSH Key Pair ---")

	err := arubaClient.FromCompute().KeyPairs().Delete(ctx, kp)
	if err != nil {
		var httpErr *aruba.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Failed to delete SSH key pair - Status: %d, Error: %s", httpErr.StatusCode, httpErr.Error())
		} else {
			log.Printf("Error deleting SSH key pair: %v", err)
		}
		return
	}
	fmt.Printf("✓ Deleted SSH key pair: %s\n", kp.Name())
}
