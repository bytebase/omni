package analysis

import (
	"context"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/catalog"
)

var (
	_ analysis.MetadataMigrator  = Engine
	_ analysis.MetadataExtractor = Engine
)

var metadataMigration = analysis.MySQLMetadataMigration[*catalog.Catalog]{
	Load: func(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, database string) (*catalog.Catalog, []map[string]error, error) {
		c := catalog.New()
		report, err := c.LoadMetadataInto(ctx, meta, database)
		if err != nil {
			return nil, nil, err
		}
		return c, []map[string]error{report.Degraded, report.Missing, report.Incomplete}, nil
	},
	Migration: func(from, to *catalog.Catalog) string {
		diff := catalog.Diff(from, to)
		if diff.IsEmpty() {
			return ""
		}
		return catalog.GenerateMigration(from, to, diff).SQL()
	},
}

var metadataExtraction = analysis.MySQLMetadataExtraction{
	Database:   "tmp",
	NewCatalog: func() analysis.MySQLScriptCatalog { return walkThroughCatalog{catalog.New()} },
}

// MetadataMigration renders the migration between two snapshots, after the
// checks omni/analysis.MySQLMetadataMigration documents.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	return metadataMigration.MetadataMigration(ctx, request)
}

// ExtractMetadata runs a schema script on an empty catalog and returns the
// database it defines.
func (engine) ExtractMetadata(ctx context.Context, script string) (*metadata.DatabaseSchemaMetadata, error) {
	return metadataExtraction.ExtractMetadata(ctx, script)
}
