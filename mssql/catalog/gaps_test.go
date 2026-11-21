package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"
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
}

// TestMigrationGaps pins the changes a migration refuses or cannot compare
// exactly. Each subtest fails when its gap closes.
func TestMigrationGaps(t *testing.T) {
	t.Run("omitted column collation", func(t *testing.T) {
		// SQL Server stores the database's default collation on character
		// columns, so columnsEqual ignores collation differences whenever
		// either side omits COLLATE. Without the database's default collation
		// carried from the driver and catalog context, adding an explicit
		// COLLATE onto a column whose source omitted COLLATE produces no
		// ALTER COLUMN.
		from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(20) NULL);`)
		require.NoError(t, err)
		to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(20) COLLATE Latin1_General_CI_AS NULL);`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.NotContains(t, sql, "COLLATE Latin1_General_CI_AS")
	})
}
