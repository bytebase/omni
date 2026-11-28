package analysis

import (
	"context"

	"github.com/bytebase/omni/metadata"
)

// MetadataMigrationRequest is two synced snapshots of one database.
type MetadataMigrationRequest struct {
	// From is the schema the migration starts from and To the schema it
	// produces; a nil database or nil Metadata is an empty schema.
	From *Database
	To   *Database
	// SkipSchemas names schemas the migration leaves alone.
	SkipSchemas []string
}

// MetadataMigrator generates the DDL that takes a database from one synced
// snapshot to another.
type MetadataMigrator interface {
	// MetadataMigration returns the migration, or "" when the snapshots
	// define the same schema.
	MetadataMigration(ctx context.Context, request MetadataMigrationRequest) (string, error)
}

// MetadataExtractor builds the metadata a schema's DDL text defines.
type MetadataExtractor interface {
	ExtractMetadata(ctx context.Context, script string) (*metadata.DatabaseSchemaMetadata, error)
}
