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
// target "t": one skipped object, so no change, unless either side compares
// names as written.
func TestMSSQLEngineMetadataMigrationCaseSensitivity(t *testing.T) {
	snapshot := func(caseSensitive bool, name string, skip bool, columnType string) *omnianalysis.Database {
		return mssqlSnapshot(caseSensitive, &metadata.SchemaMetadata{Name: "dbo", Tables: []*metadata.TableMetadata{mssqlTable(name, skip, columnType)}})
	}
	for _, tc := range []struct {
		from, to bool
		changed  bool
	}{
		{from: false, to: false, changed: false},
		{from: true, to: false, changed: true},
		{from: false, to: true, changed: true},
	} {
		migration, err := Engine.MetadataMigration(context.Background(), omnianalysis.MetadataMigrationRequest{
			From: snapshot(tc.from, "T", true, "int"),
			To:   snapshot(tc.to, "t", false, "bigint"),
		})
		require.NoError(t, err)
		require.Equal(t, tc.changed, migration != "", migration)
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
