package analysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestPGMigrationSearchPath(t *testing.T) {
	for setting, want := range map[string][]string{
		"":                               nil,
		" , ":                            nil,
		`"$user", public`:                {"$user", "public"},
		`$USER, PUBLIC`:                  {"$user", "public"},
		`'$USER'`:                        {"$USER"},
		`public, pg_catalog`:             {"public", "pg_catalog"},
		`"My""Schema", 'o''k', "", info`: {`My"Schema`, "o'k", "info"},
		`"a,b", c`:                       {"a,b", "c"},
	} {
		got := migrationSearchPath(&analysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{SearchPath: setting}})
		require.Equal(t, want, got, setting)
	}
	require.Nil(t, migrationSearchPath(nil))
}

func TestPGEngineMetadataMigration(t *testing.T) {
	table := func(schema, name string) *metadata.SchemaMetadata {
		return &metadata.SchemaMetadata{Name: schema, Tables: []*metadata.TableMetadata{{
			Name:    name,
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
		}}}
	}
	to := &analysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{table("public", "t"), table("bbdataarchive", "b")}}}

	migration, err := Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{To: to, SkipSchemas: []string{"bbdataarchive"}})
	require.NoError(t, err)
	require.Contains(t, migration, `CREATE TABLE "public"."t"`)
	require.NotContains(t, migration, "bbdataarchive")

	migration, err = Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{From: to, To: to})
	require.NoError(t, err)
	require.Empty(t, migration)
}

// TestPGEngineMetadataMigrationKeepsPgCatalogPosition covers a user type named
// like a built-in: listed after public, pg_catalog lets public.text shadow
// text, which $user makes ambiguous.
func TestPGEngineMetadataMigrationKeepsPgCatalogPosition(t *testing.T) {
	snapshot := func(searchPath string, holder bool) *analysis.Database {
		schemas := []*metadata.SchemaMetadata{{
			Name: "public",
			CompositeTypes: []*metadata.CompositeTypeMetadata{{
				Name:       "text",
				Attributes: []*metadata.CompositeTypeAttribute{{Name: "value", Type: "integer"}},
			}},
		}}
		if holder {
			schemas = append(schemas, &metadata.SchemaMetadata{Name: "app", Tables: []*metadata.TableMetadata{{
				Name:    "holder",
				Columns: []*metadata.ColumnMetadata{{Name: "value", Type: "text", Position: 1}},
			}}})
		}
		return &analysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{SearchPath: searchPath, Schemas: schemas}}
	}
	migrate := func(searchPath string) error {
		_, err := Engine.MetadataMigration(context.Background(), analysis.MetadataMigrationRequest{
			From: snapshot(searchPath, false),
			To:   snapshot(searchPath, true),
		})
		return err
	}
	require.NoError(t, migrate(`"$user", public`))
	require.NoError(t, migrate(`pg_catalog, "$user", public`))
	require.Error(t, migrate(`"$user", public, pg_catalog`))
}
