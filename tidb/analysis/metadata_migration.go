package analysis

import (
	"context"
	"strings"

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

// scratchDatabase is the database a schema script runs in when it selects
// none of its own, named not to collide with a user's.
const scratchDatabase = "_bytebase_schema_scratch"

var metadataExtraction = analysis.MySQLMetadataExtraction{
	Database:   scratchDatabase,
	NewCatalog: func() analysis.MySQLScriptCatalog { return walkThroughCatalog{catalog.New()} },
	Export:     exportCurrentDatabase,
}

// MetadataMigration renders the migration between two snapshots, after the
// checks omni/analysis.MySQLMetadataMigration documents.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	return metadataMigration.MetadataMigration(ctx, request)
}

// ExtractMetadata runs a schema script on an empty catalog and returns the
// database it leaves selected.
func (engine) ExtractMetadata(ctx context.Context, script string) (*metadata.DatabaseSchemaMetadata, error) {
	return metadataExtraction.ExtractMetadata(ctx, script)
}

// exportCurrentDatabase returns the database the script left selected. The
// scratch database is returned schema-relative: unnamed, and with its
// qualifier off view bodies. A database the script names keeps its
// qualifiers, as information_schema.VIEWS syncs them. The qualifier is removed
// as text, so the same text inside a literal or an alias goes too.
func exportCurrentDatabase(c analysis.MySQLScriptCatalog) *metadata.DatabaseSchemaMetadata {
	database := c.CurrentDatabase()
	meta := c.DatabaseMetadata(database)
	if database == scratchDatabase {
		meta.Name = ""
		prefix := "`" + scratchDatabase + "`."
		for _, s := range meta.GetSchemas() {
			for _, v := range s.GetViews() {
				v.Definition = strings.ReplaceAll(v.Definition, prefix, "")
			}
		}
	}
	return meta
}
