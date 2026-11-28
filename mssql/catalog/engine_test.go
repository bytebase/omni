package catalog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	omnianalysis "github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

func TestMSSQLEngineExtractMetadata(t *testing.T) {
	script := "CREATE TABLE [dbo].[t] ([id] INT NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([id]));"
	got, err := Engine.ExtractMetadata(context.Background(), script)
	require.NoError(t, err)
	want, err := GetDatabaseMetadata(script)
	require.NoError(t, err)
	require.Equal(t, want.String(), got.String())

	_, err = Engine.ExtractMetadata(context.Background(), "CREATE TABLE [dbo].[t] (")
	require.Error(t, err)
}

func mssqlSnapshot(caseSensitive bool, schemas ...*metadata.SchemaMetadata) *omnianalysis.Database {
	return &omnianalysis.Database{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: schemas}, ObjectCaseSensitive: caseSensitive}
}

func mssqlTable(name string, skip bool, columnType string) *metadata.TableMetadata {
	return &metadata.TableMetadata{Name: name, SkipDump: skip, Columns: []*metadata.ColumnMetadata{{Name: "id", Type: columnType}}}
}

func TestMSSQLEngineMetadataMigration(t *testing.T) {
	ctx := context.Background()
	source := mssqlSnapshot(false, &metadata.SchemaMetadata{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable("t", false, "int")}})

	migration, err := Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{To: source})
	require.NoError(t, err)
	require.Contains(t, migration, "CREATE TABLE [dbo].[t]")

	migration, err = Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{From: source, To: source})
	require.NoError(t, err)
	require.Empty(t, migration)
}

// TestMSSQLEngineMetadataMigrationCaseSensitivity covers a skipped "T" and a
// target "t": one skipped object, so no change, unless either side is
// case-sensitive, where they are a case-only rename.
func TestMSSQLEngineMetadataMigrationCaseSensitivity(t *testing.T) {
	snapshot := func(caseSensitive bool, name string, skip bool, columnType string) *omnianalysis.Database {
		return mssqlSnapshot(caseSensitive, &metadata.SchemaMetadata{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable(name, skip, columnType)}})
	}
	for _, tc := range []struct {
		from, to bool
		refused  bool
	}{
		{from: false, to: false, refused: false},
		{from: true, to: false, refused: true},
		{from: false, to: true, refused: true},
	} {
		migration, err := Engine.MetadataMigration(context.Background(), omnianalysis.MetadataMigrationRequest{
			From: snapshot(tc.from, "T", true, "int"),
			To:   snapshot(tc.to, "t", false, "bigint"),
		})
		if tc.refused {
			require.ErrorContains(t, err, "case-only rename")
			continue
		}
		require.NoError(t, err)
		require.Empty(t, migration)
	}
}

func TestMSSQLEngineMetadataMigrationRefusesCaseOnlyIdentities(t *testing.T) {
	snapshot := func(caseSensitive bool, tables ...*metadata.TableMetadata) *omnianalysis.Database {
		return mssqlSnapshot(caseSensitive, &metadata.SchemaMetadata{Name: "dbo", Tables: tables})
	}
	table := func(name string, columns ...string) *metadata.TableMetadata {
		t := &metadata.TableMetadata{Name: name}
		for _, column := range columns {
			t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: column, Type: "int"})
		}
		return t
	}
	for _, tc := range []struct {
		name     string
		from, to *omnianalysis.Database
		err      string
	}{
		{"collision", snapshot(true, table("T", "id")), snapshot(true, table("T", "id"), table("t", "id")),
			`case-sensitive database has "dbo.T" and "dbo.t", which the catalog cannot distinguish`},
		{"rename", snapshot(true, table("T", "id")), snapshot(false, table("t", "id")),
			`case-only rename "dbo.T" to "dbo.t" is not supported by the catalog diff`},
		{"column rename", snapshot(false, table("t", "id")), snapshot(true, table("t", "ID")),
			`case-only rename "dbo.t.id" to "dbo.t.ID" is not supported by the catalog diff`},
		{"insensitive", snapshot(false, table("T", "id")), snapshot(false, table("t", "id")), ""},
	} {
		migration, err := Engine.MetadataMigration(context.Background(), omnianalysis.MetadataMigrationRequest{From: tc.from, To: tc.to})
		if tc.err == "" {
			require.NoError(t, err, tc.name)
			require.Empty(t, migration, tc.name)
			continue
		}
		require.EqualError(t, err, tc.err, tc.name)
	}
}

// TestMSSQLEngineMetadataMigrationDropsSkippedObjectsWithTheirSchema covers a
// removed schema: SQL Server refuses to drop a schema that still holds a table.
func TestMSSQLEngineMetadataMigrationDropsSkippedObjectsWithTheirSchema(t *testing.T) {
	source := mssqlSnapshot(false, &metadata.SchemaMetadata{Name: "app", Tables: []*metadata.TableMetadata{mssqlTable("t", true, "int")}})
	migration, err := Engine.MetadataMigration(context.Background(), omnianalysis.MetadataMigrationRequest{From: source, To: mssqlSnapshot(false)})
	require.NoError(t, err)
	require.Contains(t, migration, "DROP TABLE")
	require.True(t, source.Metadata.Schemas[0].Tables[0].SkipDump, "the request's metadata stays unchanged")
}

func TestMSSQLReconcileSkipDump(t *testing.T) {
	nilOld, nilNew := reconcileSkipDump(nil, nil, false)
	require.Empty(t, nilOld.GetSchemas())
	require.Empty(t, nilNew.GetSchemas())

	meta := func(name string, skip bool) *metadata.DatabaseSchemaMetadata {
		return &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable(name, skip, "int")}}}}
	}
	_, reconciled := reconcileSkipDump(meta("T", true), meta("t", false), true)
	require.Len(t, reconciled.Schemas[0].Tables, 2)
	require.Equal(t, "t", reconciled.Schemas[0].Tables[0].Name)
	require.Equal(t, "T", reconciled.Schemas[0].Tables[1].Name)
}

func TestMSSQLEngineMetadataMigrationSkipSchemas(t *testing.T) {
	ctx := context.Background()
	source := mssqlSnapshot(false,
		&metadata.SchemaMetadata{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable("t", false, "int")}},
		&metadata.SchemaMetadata{Name: "archive", Tables: []*metadata.TableMetadata{mssqlTable("old", false, "int")}})
	target := mssqlSnapshot(false, &metadata.SchemaMetadata{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable("t", false, "int")}})

	migration, err := Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{From: source, To: target})
	require.NoError(t, err)
	require.Contains(t, migration, "[archive]")

	migration, err = Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{From: source, To: target, SkipSchemas: []string{"archive"}})
	require.NoError(t, err)
	require.Empty(t, migration)
}
