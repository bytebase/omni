package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// TestDefinitionGaps pins the objects a definition cannot restore yet. Each
// subtest fails when its gap closes, so the case moves into the fidelity tests
// instead of staying here forever.
func TestDefinitionGaps(t *testing.T) {
	t.Run("indexed view", func(t *testing.T) {
		// The catalog rejects an index on a view, and the driver's sync has no
		// field in ViewMetadata to carry one, so an indexed view cannot be
		// rendered or round-tripped.
		_, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL);
GO
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT [id] FROM [dbo].[t];
GO
CREATE UNIQUE CLUSTERED INDEX [UX_v_id] ON [dbo].[v] ([id]);
`)
		require.Error(t, err)
		require.Contains(t, err.Error(), "indexes on views are not supported")
	})

	t.Run("transitive early function dependency", func(t *testing.T) {
		// Only the function a table default names is hoisted; the function it
		// calls is left for later, so the earlier CREATE FUNCTION refers to a
		// function that does not exist yet.
		meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Name: "dbo",
			Tables: []*metadata.TableMetadata{{
				Name: "t",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int", Nullable: false},
					{Name: "v", Type: "int", Nullable: false, Default: "([dbo].[outer_value]())"},
				},
				Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
			}},
			Functions: []*metadata.FunctionMetadata{
				{Name: "outer_value", Definition: "CREATE FUNCTION [dbo].[outer_value]() RETURNS int AS BEGIN RETURN [dbo].[inner_value]() END"},
				{Name: "inner_value", Definition: "CREATE FUNCTION [dbo].[inner_value]() RETURNS int AS BEGIN RETURN 1 END"},
			},
		}}}
		def, err := DatabaseDefinition(meta, DefinitionOptions{})
		require.NoError(t, err)
		outerAt := strings.Index(def, "CREATE FUNCTION [dbo].[outer_value]")
		innerAt := strings.Index(def, "CREATE FUNCTION [dbo].[inner_value]")
		require.GreaterOrEqual(t, outerAt, 0)
		require.GreaterOrEqual(t, innerAt, 0)
		require.Less(t, outerAt, innerAt, "gap closed: hoist transitive function dependencies")
	})

	t.Run("schema-bound function before its prerequisite table", func(t *testing.T) {
		// A table default needs the function, and a schema-bound function needs
		// the table it reads. The definition hoists every function a table
		// default names ahead of all tables, which puts this function before
		// its prerequisite table. SQL Server rejects that schema binding.
		meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Name: "dbo",
			Tables: []*metadata.TableMetadata{
				{
					Name:    "lookup",
					Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
					Indexes: []*metadata.IndexMetadata{{Name: "PK_lookup", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
				},
				{
					Name: "t",
					Columns: []*metadata.ColumnMetadata{
						{Name: "id", Type: "int", Nullable: false},
						{Name: "v", Type: "int", Nullable: false, Default: "([dbo].[lookup_value]())"},
					},
					Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
				},
			},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "lookup_value",
				Definition: "CREATE FUNCTION [dbo].[lookup_value]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[lookup]) END",
			}},
		}}}
		def, err := DatabaseDefinition(meta, DefinitionOptions{})
		require.NoError(t, err)
		functionAt := strings.Index(def, "CREATE FUNCTION [dbo].[lookup_value]")
		tableAt := strings.Index(def, "CREATE TABLE [dbo].[lookup]")
		require.GreaterOrEqual(t, functionAt, 0)
		require.GreaterOrEqual(t, tableAt, 0)
		require.Less(t, functionAt, tableAt, "gap closed: order the function after the table it reads")
	})
}

// TestMigrationGaps pins the changes a migration refuses or cannot compare
// exactly. Each subtest fails when its gap closes.
func TestMigrationGaps(t *testing.T) {
	t.Run("function used by a CHECK constraint", func(t *testing.T) {
		from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() = 1));
`)
		require.NoError(t, err)
		to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() = 2));
`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "UNSUPPORTED")
		require.Contains(t, sql, "is used by a CHECK constraint")
	})

	t.Run("generated default parenthesization", func(t *testing.T) {
		// SQL Server stores its own parenthesization for a default, so a
		// desired-state script and a synced snapshot can spell the same default
		// differently. The catalog compares the text and cannot normalize it.
		require.False(t, expressionsEqual("(0)", "((0))"), "gap closed: normalize generated defaults")
		require.False(t, expressionsEqual("GETDATE()", "(getdate())"), "gap closed: normalize generated defaults")
	})

	t.Run("function name inside a default literal", func(t *testing.T) {
		// The cleanup matches function calls outside literals, so a literal that
		// spells the name keeps its constraint.
		from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] varchar(20) NOT NULL CONSTRAINT [DF_v] DEFAULT ('foo'), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
		require.NoError(t, err)
		to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] varchar(20) NOT NULL CONSTRAINT [DF_v] DEFAULT ('foo'), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.NotContains(t, sql, "DROP CONSTRAINT [DF_v]")
	})

	t.Run("computed column comment after its rebuild", func(t *testing.T) {
		// Rebuilding a computed column drops its extended properties, but the
		// comment phase compares the old and new text and skips an unchanged
		// comment, so the rebuilt column has none.
		from := New()
		require.NoError(t, from.LoadMetadata(context.Background(), commentRebuildMetadata("AS ([a] + 1) PERSISTED")))
		to := New()
		require.NoError(t, to.LoadMetadata(context.Background(), commentRebuildMetadata("AS ([a] + 2) PERSISTED")))
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "DROP COLUMN [c]")
		require.Contains(t, sql, "sp_addextendedproperty")
	})

	t.Run("schema-bound view calling a replaced function", func(t *testing.T) {
		// The function replacement does not consider an unchanged schema-bound
		// view, so the drop runs while the view still references the function.
		from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS n;
`)
		require.NoError(t, err)
		to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS n;
`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "UNSUPPORTED")
		require.NotContains(t, sql, "DROP FUNCTION [dbo].[f]")
	})

	t.Run("function default from another schema", func(t *testing.T) {
		// Default cleanup matches the function name without its schema, so
		// dropping a.f also drops a default bound to b.f.
		from, err := LoadSQL(`
CREATE SCHEMA [a];
CREATE SCHEMA [b];
GO
CREATE FUNCTION [a].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE FUNCTION [b].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [b].[t] ([id] int NOT NULL, [v] int NOT NULL CONSTRAINT [DF_v] DEFAULT ([b].[f]()), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
		require.NoError(t, err)
		to, err := LoadSQL(`
CREATE SCHEMA [b];
GO
CREATE FUNCTION [b].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [b].[t] ([id] int NOT NULL, [v] int NOT NULL CONSTRAINT [DF_v] DEFAULT ([b].[f]()), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "DROP CONSTRAINT [DF_v]", "gap closed: match a default's function by schema")
	})

	t.Run("schema-bound function reads a dropped table", func(t *testing.T) {
		// The table drop phase runs before the function drops, so a
		// schema-bound function still references its table.
		from, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
GO
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[t]) END;
`)
		require.NoError(t, err)
		to := New()
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Less(t, strings.Index(sql, "DROP TABLE [dbo].[t]"), strings.Index(sql, "DROP FUNCTION [dbo].[f]"), "gap closed: drop the function before its table or refuse")
	})

	t.Run("omitted column collation", func(t *testing.T) {
		// SQL Server stores the database's default collation on the column, so
		// a desired column without COLLATE mismatches the synced column.
		from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(20) NULL);`)
		require.NoError(t, err)
		to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(20) COLLATE Latin1_General_CI_AS NULL);`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "COLLATE Latin1_General_CI_AS", "gap closed: resolve an omitted collation to the database default")
	})
}

func commentRebuildMetadata(computedType string) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "t",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", Nullable: false},
				{Name: "a", Type: "int", Nullable: true},
				{Name: "c", Type: computedType, Nullable: true, Comment: "note"},
			},
			Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
		}},
	}}}
}
