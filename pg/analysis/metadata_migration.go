package analysis

import (
	"context"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/catalog"
)

var _ analysis.MetadataMigrator = Engine

// MetadataMigration generates the migration with catalog.MetadataMigration,
// each side resolving names on its own synced search_path.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	return catalog.MetadataMigration(ctx, databaseMetadata(request.From), databaseMetadata(request.To), catalog.MetadataMigrationOptions{
		SourceSearchPath: migrationSearchPath(request.From),
		TargetSearchPath: migrationSearchPath(request.To),
		SkipSchemas:      request.SkipSchemas,
	})
}

func databaseMetadata(database *analysis.Database) *metadata.DatabaseSchemaMetadata {
	if database == nil {
		return nil
	}
	return database.Metadata
}

// migrationSearchPath returns the database's search_path entries with $user
// and system schemas kept, since pg_catalog's position decides whether a user
// schema shadows a built-in.
func migrationSearchPath(database *analysis.Database) []string {
	return analysis.SearchPathEntries(databaseMetadata(database).GetSearchPath())
}
