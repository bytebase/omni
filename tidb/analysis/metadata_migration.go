package analysis

import (
	"context"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/catalog"
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

// scratchDatabase names the database a schema script runs in when it selects
// none of its own. ExtractMetadata adds a numeric suffix when the script
// mentions the name, so it never collides with a user's database or alias.
const scratchDatabase = "_bytebase_schema_scratch"

// MetadataMigration renders the migration between two snapshots, after the
// checks omni/analysis.MySQLMetadataMigration documents.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	return metadataMigration.MetadataMigration(ctx, request)
}

// ExtractMetadata runs a schema script on an empty catalog and returns the
// database it leaves selected.
func (engine) ExtractMetadata(ctx context.Context, script string) (*metadata.DatabaseSchemaMetadata, error) {
	scratch := analysis.NameAbsentFrom(scratchDatabase, script)
	return analysis.MySQLMetadataExtraction{
		Database:   scratch,
		NewCatalog: func() analysis.MySQLScriptCatalog { return walkThroughCatalog{catalog.New()} },
		Export: func(c analysis.MySQLScriptCatalog) *metadata.DatabaseSchemaMetadata {
			return exportCurrentDatabase(c, scratch)
		},
	}.ExtractMetadata(ctx, script)
}

// exportCurrentDatabase returns the database the script left selected. The
// scratch database is returned schema-relative: unnamed, and with its
// qualifier off view bodies. A database the script names keeps its
// qualifiers, as information_schema.VIEWS syncs them.
func exportCurrentDatabase(c analysis.MySQLScriptCatalog, scratch string) *metadata.DatabaseSchemaMetadata {
	database := c.CurrentDatabase()
	meta := c.DatabaseMetadata(database)
	if database == scratch {
		meta.Name = ""
		for _, s := range meta.GetSchemas() {
			for _, v := range s.GetViews() {
				v.Definition = analysis.StripMySQLDatabaseQualifier(v.Definition, scratch)
			}
		}
	}
	return meta
}
