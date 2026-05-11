package main

import (
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/Arubacloud/sdk-go/pkg/multitenant"
)

// runMultitenancyExample runs two simple multitenancy examples:
// 1) tenant clients from key-pair credentials
// 2) tenant clients from Vault credentials repository
func runMultitenancyExample() {
	runMultitenancyKeyPairExample()
	runMultitenancyVaultExample()
	runMultitenancyReconcilerStyleExample()
}

// runMultitenancyKeyPairExample creates one SDK client per tenant using
// aruba.DefaultOptions(clientID, clientSecret).
func runMultitenancyKeyPairExample() {

	// Tenant 1 credentials
	clientIDTenant1 := "tenant1_client_id"
	clientSecretTenant1 := "tenant1_client_secret"
	tenant1Options := aruba.DefaultOptions(clientIDTenant1, clientSecretTenant1)

	log.Printf("tenant-1 credentials: %s, %s", clientIDTenant1, clientSecretTenant1)
	// Tenant 2 credentials
	clientIDTenant2 := "tenant2_client_id"
	clientSecretTenant2 := "tenant2_client_secret"
	tenant2Options := aruba.DefaultOptions(clientIDTenant2, clientSecretTenant2)

	mt := multitenant.New()

	if err := mt.NewFromOptions("tenant-1", tenant1Options); err != nil {
		log.Fatalf("failed to initialize tenant-1 client from key pair: %v", err)
	}
	if err := mt.NewFromOptions("tenant-2", tenant2Options); err != nil {
		log.Fatalf("failed to initialize tenant-2 client from key pair: %v", err)
	}

	tenant1Client, ok := mt.Get("tenant-1")
	if !ok || tenant1Client == nil {
		log.Fatal("tenant-1 client not found")
	}
	tenant2Client, ok := mt.Get("tenant-2")
	if !ok || tenant2Client == nil {
		log.Fatal("tenant-2 client not found")
	}

	fmt.Println("multitenancy key-pair example initialized successfully")
}

// runMultitenancyVaultExample creates one SDK client per tenant using
// the same base SDK options, then customizes Vault settings per tenant.
//
// Example Vault UI path:
// http://localhost:56163/ui/vault/secrets/kv/kv/ARU-297647
// In this case:
// - kvMount = "kv"
// - kvPath  = "ARU-000000"
// - under the kvPath we have two secrets "client-id" "client-secret"
func runMultitenancyVaultExample() {
	// Shared base options for all tenants (logging, scopes, repositories, etc.).
	baseOptions := aruba.NewOptions().
		WithNoLogs()

	// Shared Vault auth/mount settings.
	kvMount := "kv"
	namespace := ""
	rolePath := "approle"
	roleID := "shared-role-id"
	secretID := "shared-secret-id"

	// Tenant-specific Vault endpoints/paths.
	// If you use a single Vault instance, keep vaultURI the same and change only kvPath.
	tenant1Vault := baseOptions.DeepCopy().
		WithVaultCredentialsRepository("http://vault0.default.svc.cluster.local:8200", kvMount, "ARU-000000", namespace, rolePath, roleID, secretID)
	tenant2Vault := baseOptions.DeepCopy().
		WithVaultCredentialsRepository("http://vault0.default.svc.cluster.local:8200", kvMount, "ARU-123456", namespace, rolePath, roleID, secretID)

	mt := multitenant.New()

	if err := mt.NewFromOptions("ARU-000000", tenant1Vault); err != nil {
		log.Fatalf("failed to initialize ARU-000000 client from vault options: %v", err)
	}
	if err := mt.NewFromOptions("ARU-000001", tenant2Vault); err != nil {
		log.Fatalf("failed to initialize ARU-000001 client from vault options: %v", err)
	}

	fmt.Println("multitenancy vault example initialized successfully")
}

// multitenancyExampleConfig mimics operator/reconciler configuration.
type multitenancyExampleConfig struct {
	APIGateway string

	// Vault settings
	VaultAddress string
	KVMount      string
	Namespace    string
	RolePath     string
	RoleID       string
	RoleSecret   string
}

// multitenancyExampleReconciler demonstrates the real-world usage pattern:
// cache hit -> return client, otherwise build options, create client, cache it.
type multitenancyExampleReconciler struct {
	config            multitenancyExampleConfig
	multiTenantClient multitenant.Multitenant
}

func runMultitenancyReconcilerStyleExample() {
	r := &multitenancyExampleReconciler{
		config: multitenancyExampleConfig{
			APIGateway:   "https://api.arubacloud.com",
			VaultAddress: "http://vault0.default.svc.cluster.local:8200",
			KVMount:      "kw",
			Namespace:    "",
			RolePath:     "approle",
			RoleID:       "shared-role-id",
			RoleSecret:   "shared-secret-id",
		},
		multiTenantClient: multitenant.New(),
	}

	if _, err := r.ArubaClient("ARU-000000"); err != nil {
		log.Fatalf("reconciler-style example failed: %v", err)
	}
	fmt.Println("multitenancy reconciler-style example initialized successfully")
}

// ArubaClient returns an authenticated Aruba cloud API client scoped to the tenant associated
// with the given resource. It first checks if a client for the tenant already exists in
// the multi-tenant client cache, and if not, it creates a new client using the Reconciler's
// configuration (either Vault-based or direct credentials) and adds it to the cache for
// future use. Errors during client creation are returned to the caller for handling.
func (r *multitenancyExampleReconciler) ArubaClient(tenant string) (aruba.Client, error) {
	c, ok := r.multiTenantClient.Get(tenant)
	if ok {
		return c, nil
	}

	options := aruba.NewOptions().WithBaseURL(r.config.APIGateway).WithDefaultTokenIssuerURL()
	options = options.WithVaultCredentialsRepository(
		r.config.VaultAddress,
		r.config.KVMount,
		tenant, // tenant maps to kvPath (e.g. ARU-000000)
		r.config.Namespace,
		r.config.RolePath,
		r.config.RoleID,
		r.config.RoleSecret,
	)

	arubaClient, err := aruba.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aruba client: %w", err)
	}

	r.multiTenantClient.Add(tenant, arubaClient)

	return arubaClient, nil
}
