package storage

// API path constants for storage resources
const (

	// Storage Bucket paths
	BlockStoragesPath = "/projects/%s/providers/Aruba.Storage/blockStorages"
	BlockStoragePath  = "/projects/%s/providers/Aruba.Storage/blockStorages/%s"

	//Snapshot paths
	SnapshotsPath = "/projects/%s/providers/Aruba.Storage/snapshots"
	SnapshotPath  = "/projects/%s/providers/Aruba.Storage/snapshots/%s"

	//Backup paths
	BackupsPath = "/projects/%s/providers/Aruba.Storage/backups"
	BackupPath  = "/projects/%s/providers/Aruba.Storage/backups/%s"

	//Restore paths
	RestoresPath = "/projects/%s/providers/Aruba.Storage/backups/%s/restores"
	RestorePath  = "/projects/%s/providers/Aruba.Storage/backups/%s/restores/%s"
)
