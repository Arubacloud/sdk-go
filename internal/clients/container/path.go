package container

const (
	// KaaSPath is the base path for KaaS operations
	KaaSPath = "/projects/%s/providers/Aruba.Container/kaas"

	// KaaSItemPath is the path for a specific KaaS cluster
	KaaSItemPath = "/projects/%s/providers/Aruba.Container/kaas/%s"

	// ContainerRegistryPath is the base path for container registry operations
	ContainerRegistryPath = "/projects/%s/providers/Aruba.Container/registries"

	// ContainerRegistryItemPath is the path for a specific container registry
	ContainerRegistryItemPath = "/projects/%s/providers/Aruba.Container/registries/%s"
)
