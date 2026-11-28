package catalog

import (
	"context"

	omnianalysis "github.com/bytebase/omni/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

// Engine is CockroachDB's metadata migrator. It lives here rather than in
// cockroachdb/analysis because this package imports that one.
var Engine engine

type engine struct{}

var _ omnianalysis.MetadataMigrator = Engine

func (engine) MetadataMigration(_ context.Context, request omnianalysis.MetadataMigrationRequest) (string, error) {
	return MetadataMigration(databaseMetadata(request.From), databaseMetadata(request.To), MetadataMigrationOptions{SkipSchemas: request.SkipSchemas})
}

func databaseMetadata(database *omnianalysis.Database) *metadatapb.DatabaseSchemaMetadata {
	if database == nil {
		return nil
	}
	return database.Metadata
}
