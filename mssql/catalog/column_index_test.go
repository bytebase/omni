package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// A column's INDEX keys the column, each one a column declares is kept, and a
// column's FOREIGN KEY may name its own columns (BYT-10538).
func TestColumnIndexesAndForeignKeyColumns(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE dbo.p (a INT NOT NULL, b INT NOT NULL, CONSTRAINT pk_p PRIMARY KEY (a, b));
CREATE TABLE dbo.c (x INT INDEX ix_x INDEX ix_x2 (x DESC), y INT NOT NULL CONSTRAINT fk_c FOREIGN KEY (x, y) REFERENCES dbo.p (a, b));
`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "c")
	require.NotNil(t, table)
	indexes := map[string]*metadata.IndexMetadata{}
	for _, index := range table.Indexes {
		indexes[index.Name] = index
	}
	require.Contains(t, indexes, "ix_x")
	require.Equal(t, []string{"x"}, indexes["ix_x"].Expressions)
	require.Equal(t, []bool{false}, indexes["ix_x"].Descending)
	require.Contains(t, indexes, "ix_x2")
	require.Equal(t, []bool{true}, indexes["ix_x2"].Descending)
	require.Len(t, table.ForeignKeys, 1)
	require.Equal(t, []string{"x", "y"}, table.ForeignKeys[0].Columns)
	require.Equal(t, []string{"a", "b"}, table.ForeignKeys[0].ReferencedColumns)

	// ALTER TABLE ... ADD validates the columns the FOREIGN KEY names.
	results, err := c.Exec(`ALTER TABLE dbo.c ADD z INT FOREIGN KEY (x, z) REFERENCES dbo.p (a, b);`, nil)
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error)
	}
	table = c.GetTable("dbo", "c")
	require.Len(t, table.ForeignKeys, 2)
	require.Equal(t, []string{"x", "z"}, table.ForeignKeys[1].Columns)
}
