package catalog

import (
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
}
