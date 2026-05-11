package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

func createKMS(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref) *aruba.KMS {
	fmt.Println("--- KMS Instance ---")

	k := aruba.NewKMS().
		IntoProject(proj).
		WithName(resourceName(NameKMS)).
		AddTag("security").
		AddTag("encryption").
		InRegion("ITBG-Bergamo").
		WithBillingPeriod("Hour")

	result, err := arubaClient.FromSecurity().KMS().Create(ctx, k)
	if err != nil {
		log.Printf("Error creating KMS: %v", err)
		return nil
	}

	fmt.Printf("✓ Created KMS instance: %s (ID: %s)\n", result.Name(), result.KMSID())

	if err := result.WaitUntilReady(ctx); err != nil {
		log.Printf("KMS %s did not become Ready: %v", result.Name(), err)
	}

	return result
}

func createKMSKey(ctx context.Context, arubaClient aruba.Client, kmsParent *aruba.KMS) *aruba.Key {
	fmt.Println("--- KMS Cryptographic Key ---")

	if err := waitForDependencies(ctx, "KMS Key", map[string]waitFunc{
		"KMS": kmsParent.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	key := aruba.NewKey().
		IntoKMS(kmsParent).
		WithName(resourceName(NameKMSKey)).
		OfAlgorithm(aruba.KeyAlgorithmAes)

	result, err := arubaClient.FromSecurity().Keys().Create(ctx, key)
	if err != nil {
		log.Printf("Error creating Key: %v", err)
		return nil
	}

	fmt.Printf("✓ Created Key: %s (Algorithm: %s, Type: %s)\n",
		result.Name(),
		result.Algorithm(),
		result.Type())

	// Creating a Key transitions the parent KMS out of Active. Wait for it to settle
	// before returning so that createKmip can safely fire next.
	waitPostDependencies(ctx, "KMS Key", map[string]waitFunc{
		"KMS": kmsParent.WaitUntilActive,
	})

	return result
}

func createKmip(ctx context.Context, arubaClient aruba.Client, kmsParent *aruba.KMS) *aruba.Kmip {
	fmt.Println("--- KMIP Service ---")

	if err := waitForDependencies(ctx, "KMIP", map[string]waitFunc{
		"KMS": kmsParent.WaitUntilActive,
	}); err != nil {
		log.Printf("%v", err)
		return nil
	}

	km := aruba.NewKmip().IntoKMS(kmsParent).WithName(resourceName(NameKmip))

	created, err := arubaClient.FromSecurity().Kmips().Create(ctx, km)
	if err != nil {
		log.Printf("Error creating KMIP service: %v", err)
		return nil
	}

	fmt.Printf("✓ Created KMIP service: %s (ID: %s, Status: %s)\n",
		created.Name(),
		created.KmipID(),
		created.KmipStatus())

	if err := created.WaitUntilReady(ctx); err != nil {
		log.Printf("KMIP %s did not become Ready: %v", created.Name(), err)
	}

	return created
}

func downloadKmipCertificate(ctx context.Context, arubaClient aruba.Client, kmip *aruba.Kmip) *aruba.KmipCertificate {
	fmt.Println("--- KMIP Certificate Download ---")

	fmt.Println("⏳ Waiting for KMIP certificate to become available...")
	if err := kmip.WaitUntilCertificateAvailable(ctx); err != nil {
		log.Printf("KMIP certificate did not become available: %v", err)
		return nil
	}
	fmt.Println("✓ KMIP certificate is now available")

	cert, err := arubaClient.FromSecurity().Kmips().Download(ctx, kmip)
	if err != nil {
		log.Printf("Error downloading KMIP certificate: %v", err)
		return nil
	}

	if cert != nil {
		fmt.Printf("✓ Downloaded KMIP certificate\n")
		fmt.Printf("  - Key length: %d bytes\n", len(cert.Key()))
		fmt.Printf("  - Cert length: %d bytes\n", len(cert.Cert()))
	}

	return cert
}

func deleteKMS(ctx context.Context, arubaClient aruba.Client, k *aruba.KMS) {
	fmt.Println("--- Deleting KMS Instance ---")

	if err := arubaClient.FromSecurity().KMS().Delete(ctx, k); err != nil {
		log.Printf("Error deleting KMS: %v", err)
		return
	}
	fmt.Printf("✓ Deleted KMS instance: %s\n", k.KMSID())
}

func deleteKMSKey(ctx context.Context, arubaClient aruba.Client, key *aruba.Key) {
	fmt.Println("--- Deleting KMS Key ---")

	err := arubaClient.FromSecurity().Keys().Delete(ctx, key)
	if err != nil {
		log.Printf("Error deleting KMS key: %v", err)
		return
	}
	fmt.Printf("✓ Deleted KMS key: %s\n", key.KeyID())
}

func deleteKmip(ctx context.Context, arubaClient aruba.Client, km *aruba.Kmip) {
	fmt.Println("--- Deleting KMIP Service ---")

	err := arubaClient.FromSecurity().Kmips().Delete(ctx, km)
	if err != nil {
		log.Printf("Error deleting KMIP service: %v", err)
		return
	}
	fmt.Printf("✓ Deleted KMIP service: %s\n", km.KmipID())
}
