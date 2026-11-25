package catalog_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
	"github.com/bytebase/omni/metadata"
)

// snapshotOf is a synced snapshot of database name, as CockroachDB states it:
// every relation a view, index or routine body names carries the database.
func snapshotOf(name string) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: name,
		Schemas: []*metadata.SchemaMetadata{{
			Name: "app",
			Tables: []*metadata.TableMetadata{
				{
					Name: "items",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "bigint"},
						{Name: "qty", Type: "integer", Nullable: true},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "items_pkey", Primary: true, Unique: true, Expressions: []string{"id"}, Definition: "CREATE UNIQUE INDEX items_pkey ON " + name + ".app.items USING btree (id ASC)"},
						{Name: "items_qty_idx", Expressions: []string{"qty"}, Definition: "CREATE INDEX items_qty_idx ON " + name + ".app.items USING btree (qty ASC) WHERE (qty > 0)"},
					},
				},
				{
					Name: "events",
					Columns: []*metadata.ColumnMetadata{
						{Name: "note", Type: "text", Nullable: true},
						{Name: "rowid", Type: "bigint", Default: "unique_rowid()", IsInvisible: true},
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "events_pkey", Primary: true, Unique: true, Expressions: []string{"rowid"}, Definition: "CREATE UNIQUE INDEX events_pkey ON " + name + ".app.events USING btree (rowid ASC)"},
					},
				},
			},
			Views: []*metadata.ViewMetadata{{
				Name:              "stocked",
				Definition:        "SELECT id FROM " + name + ".app.items WHERE qty > 0",
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "app", Table: "items", Column: "id"}},
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "total",
				Signature:  "total()",
				Definition: "CREATE FUNCTION app.total()\n\tRETURNS INT8\n\tLANGUAGE SQL\n\tAS $$\n\tSELECT sum(qty) FROM " + name + ".app.items;\n$$",
			}},
		}},
	}
}

func TestDefinitionDropsDatabaseName(t *testing.T) {
	def, err := catalog.DatabaseDefinition(snapshotOf("db_one"), catalog.DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "db_one")
	require.Contains(t, def, "CREATE INDEX items_qty_idx ON app.items (qty ASC) WHERE (qty > 0);")
	require.Contains(t, def, "SELECT id FROM app.items WHERE qty > 0;")
	require.Contains(t, def, "SELECT sum(qty) FROM app.items;")
	// INT4 is stated as "integer", which CockroachDB reads as INT8.
	require.Contains(t, def, `"qty" INT4`)
	// The hidden rowid key is CockroachDB's own.
	require.Contains(t, def, "CREATE TABLE \"app\".\"events\" (\n    \"note\" text\n);")

	sdl, err := catalog.DatabaseDefinition(snapshotOf("db_one"), catalog.DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	require.NotContains(t, sdl, "db_one")

	// One object rendered alone has no database to drop.
	view, err := catalog.ViewDefinition("app", snapshotOf("db_one").Schemas[0].Views[0])
	require.NoError(t, err)
	require.Equal(t, "CREATE VIEW \"app\".\"stocked\" AS \nSELECT id FROM db_one.app.items WHERE qty > 0;\n\n", view)
}

func TestMetadataMigrationAcrossDatabases(t *testing.T) {
	migration, err := catalog.MetadataMigration(snapshotOf("db_one"), snapshotOf("db_two"), catalog.MetadataMigrationOptions{})
	require.NoError(t, err)
	require.Empty(t, migration)
}

func TestMetadataMigrationSkipsSchemas(t *testing.T) {
	target := snapshotOf("db")
	target.Schemas = append(target.Schemas, &metadata.SchemaMetadata{Name: "bbdataarchive", Tables: []*metadata.TableMetadata{{Name: "t"}}})
	migration, err := catalog.MetadataMigration(snapshotOf("db"), target, catalog.MetadataMigrationOptions{SkipSchemas: []string{"bbdataarchive"}})
	require.NoError(t, err)
	require.Empty(t, migration)

	migration, err = catalog.MetadataMigration(snapshotOf("db"), target, catalog.MetadataMigrationOptions{})
	require.NoError(t, err)
	require.Equal(t, "CREATE SCHEMA \"bbdataarchive\";\nCREATE TABLE \"bbdataarchive\".\"t\" (\n)", migration)
}

func TestMetadataMigrationRefuses(t *testing.T) {
	for name, mutate := range map[string]func(*metadata.SchemaMetadata){
		"enum-reorder": func(s *metadata.SchemaMetadata) { s.EnumTypes[0].Values = []string{"b", "a"} },
		"drop-primary-key": func(s *metadata.SchemaMetadata) {
			s.Tables[0].Indexes = s.Tables[0].Indexes[1:]
		},
	} {
		t.Run(name, func(t *testing.T) {
			source := snapshotOf("db")
			source.Schemas[0].EnumTypes = []*metadata.EnumTypeMetadata{{Name: "e", Values: []string{"a", "b"}}}
			target := snapshotOf("db")
			target.Schemas[0].EnumTypes = []*metadata.EnumTypeMetadata{{Name: "e", Values: []string{"a", "b"}}}
			mutate(target.Schemas[0])
			_, err := catalog.MetadataMigration(source, target, catalog.MetadataMigrationOptions{})
			require.Error(t, err)
		})
	}
}

func TestMultiFileDefinition(t *testing.T) {
	files, err := catalog.MultiFileDefinition(snapshotOf("db"))
	require.NoError(t, err)
	var names []string
	for _, file := range files {
		names = append(names, file.Name)
		require.NotContains(t, file.Content, "db.app")
	}
	require.Equal(t, []string{
		"schemas/app/schema.sql",
		"schemas/app/tables/items.sql",
		"schemas/app/tables/events.sql",
		"schemas/app/views/stocked.sql",
		"schemas/app/functions/total.sql",
	}, names)
	require.True(t, strings.HasPrefix(files[1].Content, "CREATE TABLE \"app\".\"items\" (\n    \"id\" bigint NOT NULL,\n    \"qty\" INT4,\n    CONSTRAINT \"items_pkey\" PRIMARY KEY (id ASC)\n);"))

	files, err = catalog.MultiFileDefinition(&metadata.DatabaseSchemaMetadata{})
	require.NoError(t, err)
	require.NotNil(t, files)
	require.Empty(t, files)
}
