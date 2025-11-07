package database

// API path constants for database resources
const (
	// DBaaS paths
	DBaaSPath     = "/projects/%s/providers/Aruba.Database/dbaas"
	DBaaSItemPath = "/projects/%s/providers/Aruba.Database/dbaas/%s"

	// Database paths
	DatabaseInstancesPath = "/projects/%s/providers/Aruba.Database/dbaas/%s/databases"
	DatabaseInstancePath  = "/projects/%s/providers/Aruba.Database/dbaas/%s/databases/%s"

	// Backup paths
	BackupsPath = "/projects/%s/providers/Aruba.Database/backups"
	BackupPath  = "/projects/%s/providers/Aruba.Database/backups/%s"

	// GrantDatabase Paths
	GrantsPath    = "/projects/%s/providers/Aruba.Database/dbaas/%s/databases/%s/grants"
	GrantItemPath = "/projects/%s/providers/Aruba.Database/dbaas/%s/databases/%s/grants/%s"

	// User paths
	UsersPath    = "/projects/%s/providers/Aruba.Database/dbaas/%s/users"
	UserItemPath = "/projects/%s/providers/Aruba.Database/dbaas/%s/users/%s"
)
