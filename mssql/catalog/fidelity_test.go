package catalog_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/driver"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestMain accounts for this package's tapes and removes the SQL Server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// fidelityCase is one schema definition that must survive a trip through a real
// SQL Server: the server stores it, the driver syncs it, the catalog renders it,
// the server runs that rendering, and the catalog renders the same definition
// again.
type fidelityCase struct {
	name string
	sql  string
}

func fidelityCases() []fidelityCase {
	return []fidelityCase{
		{
			name: "columns-keys-constraints",
			sql: `CREATE SCHEMA [fid_cols];
GO
CREATE TABLE [fid_cols].[customers] (
    [id] INT NOT NULL,
    [name] NVARCHAR(100) NOT NULL,
    [email] NVARCHAR(255) NULL,
    [created_at] DATETIME2(3) NOT NULL CONSTRAINT [DF_customers_created_at] DEFAULT (SYSUTCDATETIME()),
    [balance] DECIMAL(12,2) NOT NULL CONSTRAINT [DF_customers_balance] DEFAULT (0),
    CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC)
);
GO
CREATE TABLE [fid_cols].[orders] (
    [id] INT IDENTITY(1,1) NOT NULL,
    [customer_id] INT NOT NULL,
    [status] VARCHAR(20) COLLATE Latin1_General_CS_AS NOT NULL,
    [note] NVARCHAR(50) NULL,
    [total] MONEY NOT NULL CONSTRAINT [DF_orders_total] DEFAULT (0),
    CONSTRAINT [PK_orders] PRIMARY KEY NONCLUSTERED ([id] ASC),
    CONSTRAINT [UQ_orders_customer_status] UNIQUE CLUSTERED ([customer_id] ASC, [status] ASC),
    CONSTRAINT [CK_orders_total] CHECK ([total] >= 0),
    CONSTRAINT [FK_orders_customers] FOREIGN KEY ([customer_id]) REFERENCES [fid_cols].[customers] ([id]) ON DELETE CASCADE
);
GO
CREATE NONCLUSTERED INDEX [IX_orders_note] ON [fid_cols].[orders] ([note] DESC, [status] ASC);
GO
CREATE UNIQUE NONCLUSTERED INDEX [UX_orders_status_id] ON [fid_cols].[orders] ([status] ASC, [id] DESC);
GO
`,
		},
		{
			name: "views-functions-procedures-sequences",
			sql: `CREATE SCHEMA [fid_routines];
GO
CREATE TABLE [fid_routines].[items] (
    [id] INT NOT NULL,
    [name] NVARCHAR(100) NOT NULL,
    CONSTRAINT [PK_items] PRIMARY KEY ([id])
);
GO
CREATE VIEW [fid_routines].[v_items] AS
SELECT [id], [name] FROM [fid_routines].[items];
GO
CREATE VIEW [fid_routines].[v_item_names] AS
SELECT [name] FROM [fid_routines].[v_items];
GO
CREATE FUNCTION [fid_routines].[item_name] (@id INT)
RETURNS NVARCHAR(100)
AS
BEGIN
    RETURN (SELECT [name] FROM [fid_routines].[items] WHERE [id] = @id)
END
GO
CREATE PROCEDURE [fid_routines].[get_items]
AS
BEGIN
    SELECT [id], [name] FROM [fid_routines].[items];
END
GO
CREATE SEQUENCE [fid_routines].[item_seq] AS BIGINT;
GO
`,
		},
		{
			name: "cross-schema",
			sql: `CREATE SCHEMA [fid_a];
GO
CREATE SCHEMA [fid_b];
GO
CREATE TABLE [fid_a].[parents] (
    [id] INT NOT NULL,
    CONSTRAINT [PK_parents] PRIMARY KEY ([id])
);
GO
CREATE TABLE [fid_b].[children] (
    [id] INT NOT NULL,
    [parent_id] INT NOT NULL,
    CONSTRAINT [PK_children] PRIMARY KEY ([id]),
    CONSTRAINT [FK_children_parents] FOREIGN KEY ([parent_id]) REFERENCES [fid_a].[parents] ([id]) ON UPDATE NO ACTION
);
GO
`,
		},
	}
}

func TestDefinitionFidelity(t *testing.T) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	for _, tc := range fidelityCases() {
		t.Run(tc.name, func(t *testing.T) {
			execScript(t, db, tc.sql)
			stored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			def, err := catalog.DatabaseDefinition(stored, catalog.DefinitionOptions{})
			require.NoError(t, err)

			execScript(t, db, dropScript(stored))

			execScript(t, db, def)
			restored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			restoredDef, err := catalog.DatabaseDefinition(filterTo(restored, stored), catalog.DefinitionOptions{})
			require.NoError(t, err)
			require.Equal(t, def, restoredDef)
		})
	}
}

func TestMigrationFidelity(t *testing.T) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	for _, tc := range fidelityCases() {
		t.Run(tc.name, func(t *testing.T) {
			execScript(t, db, tc.sql)
			stored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			def, err := catalog.DatabaseDefinition(stored, catalog.DefinitionOptions{})
			require.NoError(t, err)

			from := catalog.New()
			require.NoError(t, from.LoadMetadata(ctx, stored))
			target := mutate(stored, "fid_added_"+strings.ReplaceAll(tc.name, "-", "_"))
			to := catalog.New()
			require.NoError(t, to.LoadMetadata(ctx, target))
			diff := catalog.Diff(from, to)
			plan := catalog.GenerateMigration(from, to, diff)
			require.False(t, plan.IsEmpty())

			execScript(t, db, plan.SQL())
			restored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			want, err := catalog.DatabaseDefinition(target, catalog.DefinitionOptions{})
			require.NoError(t, err)
			got, err := catalog.DatabaseDefinition(filterTo(restored, target), catalog.DefinitionOptions{})
			require.NoError(t, err)
			require.Equal(t, want, got, "migration %q from\n%s", plan.SQL(), def)
		})
	}
}

// mutate adds a schema and table to the stored snapshot, so every case proves
// the migration creates objects on a server that already holds the source.
func mutate(stored *metadata.DatabaseSchemaMetadata, addedSchema string) *metadata.DatabaseSchemaMetadata {
	target := cloneMetadata(stored)
	target.Schemas = append(target.Schemas, &metadata.SchemaMetadata{
		Name: addedSchema,
		Tables: []*metadata.TableMetadata{{
			Name: "added",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", Nullable: false},
				// SQL Server fills the database collation on a character column
				// written without one, so the target holds the stored form.
				{Name: "note", Type: "nvarchar(20)", Nullable: true, Collation: defaultCollation(stored)},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PK_added", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}, Descending: []bool{false}},
				{Name: "IX_added_note", Type: "NONCLUSTERED", Expressions: []string{"note"}, Descending: []bool{false}},
			},
			CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "CK_added_note", Expression: "([note] IS NULL OR [note]<>'')"}},
		}},
	})
	slices.SortFunc(target.Schemas, func(a, b *metadata.SchemaMetadata) int { return strings.Compare(a.Name, b.Name) })
	return target
}

// defaultCollation returns the collation the server stored on a character
// column, which is the database's default.
func defaultCollation(snapshot *metadata.DatabaseSchemaMetadata) string {
	for _, schema := range snapshot.Schemas {
		for _, table := range schema.Tables {
			for _, column := range table.Columns {
				if column.Collation != "" {
					return column.Collation
				}
			}
		}
	}
	return "SQL_Latin1_General_CP1_CI_AS"
}

func filterTo(synced, wanted *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
	names := make(map[string]bool)
	for _, schema := range wanted.Schemas {
		names[strings.ToLower(schema.Name)] = true
	}
	filtered := &metadata.DatabaseSchemaMetadata{Name: synced.Name, DefaultSchema: synced.DefaultSchema}
	for _, schema := range synced.Schemas {
		if names[strings.ToLower(schema.Name)] {
			filtered.Schemas = append(filtered.Schemas, schema)
		}
	}
	return filtered
}

// dropScript empties the schemas a snapshot holds so its definition can be
// replayed, leaving the database itself in place.
func dropScript(snapshot *metadata.DatabaseSchemaMetadata) string {
	var lines []string
	schemas := slices.Clone(snapshot.Schemas)
	slices.SortFunc(schemas, func(a, b *metadata.SchemaMetadata) int { return strings.Compare(a.Name, b.Name) })
	for _, schema := range schemas {
		views := slices.Clone(schema.Views)
		slices.SortFunc(views, func(a, b *metadata.ViewMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, view := range views {
			lines = append(lines, fmt.Sprintf("DROP VIEW IF EXISTS [%s].[%s];", schema.Name, view.Name))
		}
		procedures := slices.Clone(schema.Procedures)
		slices.SortFunc(procedures, func(a, b *metadata.ProcedureMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, procedure := range procedures {
			lines = append(lines, fmt.Sprintf("DROP PROCEDURE IF EXISTS [%s].[%s];", schema.Name, procedure.Name))
		}
		functions := slices.Clone(schema.Functions)
		slices.SortFunc(functions, func(a, b *metadata.FunctionMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, function := range functions {
			lines = append(lines, fmt.Sprintf("DROP FUNCTION IF EXISTS [%s].[%s];", schema.Name, function.Name))
		}
		sequences := slices.Clone(schema.Sequences)
		slices.SortFunc(sequences, func(a, b *metadata.SequenceMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, sequence := range sequences {
			lines = append(lines, fmt.Sprintf("DROP SEQUENCE IF EXISTS [%s].[%s];", schema.Name, sequence.Name))
		}
	}
	// Foreign keys go before any table, because they cross schemas.
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			fks := slices.Clone(table.ForeignKeys)
			slices.SortFunc(fks, func(a, b *metadata.ForeignKeyMetadata) int { return strings.Compare(a.Name, b.Name) })
			for _, fk := range fks {
				lines = append(lines, fmt.Sprintf("ALTER TABLE [%s].[%s] DROP CONSTRAINT IF EXISTS [%s];", schema.Name, table.Name, fk.Name))
			}
		}
	}
	for _, schema := range schemas {
		tables := slices.Clone(schema.Tables)
		slices.SortFunc(tables, func(a, b *metadata.TableMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, table := range tables {
			lines = append(lines, fmt.Sprintf("DROP TABLE IF EXISTS [%s].[%s];", schema.Name, table.Name))
		}
	}
	for _, schema := range schemas {
		if !strings.EqualFold(schema.Name, "dbo") {
			lines = append(lines, fmt.Sprintf("DROP SCHEMA IF EXISTS [%s];", schema.Name))
		}
	}
	return strings.Join(lines, "\n")
}

func cloneMetadata(in *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
	return proto.Clone(in).(*metadata.DatabaseSchemaMetadata)
}

// execScript sends a GO-separated script as batches.
func execScript(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	var batch strings.Builder
	flush := func() {
		text := strings.TrimSpace(batch.String())
		batch.Reset()
		if text == "" {
			return
		}
		if _, err := db.ExecContext(context.Background(), text); err != nil {
			t.Fatalf("%s:\n%v", text, err)
		}
	}
	for _, line := range strings.Split(script, "\n") {
		if strings.EqualFold(strings.TrimSpace(line), "GO") {
			flush()
			continue
		}
		batch.WriteString(line)
		batch.WriteString("\n")
	}
	flush()
}

// TestColumnDefaultChangeFidelity round-trips a column type change that
// keeps its bound default: SQL Server refuses ALTER COLUMN until the default
// is dropped, and the migration restores it afterward.
func TestColumnDefaultChangeFidelity(t *testing.T) {
	columnChangeFidelity(t, `CREATE SCHEMA [fid_default];
GO
CREATE TABLE [fid_default].[t] (
    [id] INT NOT NULL,
    [v] INT NOT NULL CONSTRAINT [DF_t_v] DEFAULT (0),
    CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
GO
`, `CREATE SCHEMA [fid_default];
GO
CREATE TABLE [fid_default].[t] (
    [id] INT NOT NULL,
    [v] BIGINT NOT NULL CONSTRAINT [DF_t_v] DEFAULT (0),
    CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
GO
`)
}

// TestComputedNullabilityFidelity round-trips a computed column whose explicit
// NOT NULL changes. ALTER COLUMN cannot express that, so the migration drops
// and recreates the column.
func TestComputedNullabilityFidelity(t *testing.T) {
	columnChangeFidelity(t, `CREATE SCHEMA [fid_computed];
GO
CREATE TABLE [fid_computed].[t] (
    [id] INT NOT NULL,
    [a] INT NULL,
    [c] AS ([a] + 1) PERSISTED,
    CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
GO
`, `CREATE SCHEMA [fid_computed];
GO
CREATE TABLE [fid_computed].[t] (
    [id] INT NOT NULL,
    [a] INT NULL,
    [c] AS ([a] + 1) PERSISTED NOT NULL,
    CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
GO
`)
}

func columnChangeFidelity(t *testing.T, sourceSQL, targetSQL string) {
	t.Helper()
	db := testcontainer.Open(t)
	ctx := context.Background()
	execScript(t, db, sourceSQL)
	stored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
	require.NoError(t, err)
	from := catalog.New()
	require.NoError(t, from.LoadMetadata(ctx, stored))

	execScript(t, db, dropScript(stored))
	execScript(t, db, targetSQL)
	target, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
	require.NoError(t, err)
	to := catalog.New()
	require.NoError(t, to.LoadMetadata(ctx, target))
	plan := catalog.GenerateMigration(from, to, catalog.Diff(from, to))
	require.False(t, plan.IsEmpty())

	execScript(t, db, dropScript(target))
	execScript(t, db, sourceSQL)
	execScript(t, db, plan.SQL())
	restored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
	require.NoError(t, err)
	want, err := catalog.DatabaseDefinition(target, catalog.DefinitionOptions{})
	require.NoError(t, err)
	got, err := catalog.DatabaseDefinition(filterTo(restored, target), catalog.DefinitionOptions{})
	require.NoError(t, err)
	require.Equal(t, want, got, "migration %q", plan.SQL())
}
