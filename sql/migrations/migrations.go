package migrations

import "embed"

//go:embed up/*.sql
var UpMigrationFiles embed.FS

//go:embed down/*.sql
var DownMigrationFiles embed.FS
