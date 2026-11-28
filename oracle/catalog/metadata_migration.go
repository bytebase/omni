package catalog

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

// Engine is Oracle's implementation of the omni/analysis metadata migration
// and extraction interfaces. It lives here rather than in oracle/analysis,
// which this package imports.
var Engine engine

type engine struct{}

var (
	_ analysis.MetadataMigrator  = Engine
	_ analysis.MetadataExtractor = Engine
)

// MetadataMigration loads both snapshots and renders the migration between
// them, ending a non-empty migration with a newline.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	from, err := loadMigrationSnapshot(ctx, request.From)
	if err != nil {
		return "", fmt.Errorf("failed to load source schema: %w", err)
	}
	to, err := loadMigrationSnapshot(ctx, request.To)
	if err != nil {
		return "", fmt.Errorf("failed to load target schema: %w", err)
	}
	diff := Diff(from, to)
	if diff.IsEmpty() {
		return "", nil
	}
	migration := GenerateMigration(from, to, diff).SQL()
	if migration != "" && !strings.HasSuffix(migration, "\n") {
		migration += "\n"
	}
	return migration, nil
}

// ExtractMetadata returns the schema a DDL script defines, as
// GetDatabaseMetadata does.
func (engine) ExtractMetadata(_ context.Context, script string) (*metadatapb.DatabaseSchemaMetadata, error) {
	return GetDatabaseMetadata(script)
}

func loadMigrationSnapshot(ctx context.Context, database *analysis.Database) (*Catalog, error) {
	meta := &metadatapb.DatabaseSchemaMetadata{}
	if database != nil && database.Metadata != nil {
		meta = database.Metadata
	}
	c := New()
	if _, err := c.LoadMetadata(ctx, stripUnrenderedFacts(meta)); err != nil {
		return nil, err
	}
	return c, nil
}

// stripUnrenderedFacts drops identity clauses and index visibility, which the
// migration leaves out until bytebase/bytebase#22119 renders them.
func stripUnrenderedFacts(meta *metadatapb.DatabaseSchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
	out := proto.CloneOf(meta)
	for _, schema := range out.GetSchemas() {
		for _, table := range schema.GetTables() {
			for _, column := range table.GetColumns() {
				column.IsIdentity = false
				column.IdentityGeneration = metadatapb.ColumnMetadata_BY_DEFAULT
				column.IdentitySeed = 0
				column.IdentityIncrement = 0
			}
			resetIndexVisibility(table.GetIndexes())
		}
		for _, view := range schema.GetMaterializedViews() {
			resetIndexVisibility(view.GetIndexes())
		}
	}
	return out
}

func resetIndexVisibility(indexes []*metadatapb.IndexMetadata) {
	for _, index := range indexes {
		index.Visible = true
	}
}
