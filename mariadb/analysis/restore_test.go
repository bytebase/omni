package analysis_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/mariadb/analysis"
	"github.com/bytebase/omni/mariadb/ast"
)

func restore(t *testing.T, statement, table string) (string, error) {
	t.Helper()
	node, err := analysis.FirstDML(statement)
	require.NoError(t, err)
	require.NotNil(t, node)
	_, normal, err := analysis.ClassifyColumns(backupCatalog(), false, "db", table)
	require.NoError(t, err)
	return analysis.RestoreSQL(node, analysis.RestoreOptions{
		Statement:      statement,
		Database:       "db",
		Table:          table,
		BackupDatabase: "bak",
		BackupTable:    "pre_" + table,
		Columns:        normal,
		Catalog:        backupCatalog(),
	})
}

func TestRestoreSQL(t *testing.T) {
	got, err := restore(t, "DELETE FROM tg WHERE a = 1", "tg")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM tg WHERE a = 1\n*/\nINSERT INTO `db`.`tg` (`a`, `b`) SELECT `a`, `b` FROM `bak`.`pre_tg`;", got)

	got, err = restore(t, "UPDATE t x JOIN t1 ON x.a = t1.a SET x.b = 1, c = 2, t1.b = 3", "t")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nUPDATE t x JOIN t1 ON x.a = t1.a SET x.b = 1, c = 2, t1.b = 3\n*/\nINSERT INTO `db`.`t` (`a`, `b`, `c`) SELECT `a`, `b`, `c` FROM `bak`.`pre_t` ON DUPLICATE KEY UPDATE `b` = VALUES(`b`), `c` = VALUES(`c`);", got)

	// The only unique key is the column the UPDATE changes.
	_, err = restore(t, "UPDATE t SET a = 1", "t")
	require.EqualError(t, err, "no disjoint unique key found for db.t")
	_, err = restore(t, "UPDATE tg SET b = 1", "tg")
	require.EqualError(t, err, "no disjoint unique key found for db.tg")

	long := "DELETE FROM t WHERE a = '" + strings.Repeat("é", 2000) + "'"
	got, err = restore(t, long, "t")
	require.NoError(t, err)
	require.Contains(t, got, strings.Repeat("é", 975)+"...\n*/")
}

func TestFirstDML(t *testing.T) {
	node, err := analysis.FirstDML("SELECT 1; INSERT INTO t VALUES (1); DELETE FROM t; UPDATE t SET a = 1")
	require.NoError(t, err)
	require.IsType(t, &ast.DeleteStmt{}, node)
	node, err = analysis.FirstDML("SELECT 1")
	require.NoError(t, err)
	require.Nil(t, node)
	_, err = analysis.FirstDML("DELETE FROM")
	require.Error(t, err)
}

func TestReferencesTable(t *testing.T) {
	tests := []struct {
		statement string
		want      bool
	}{
		{"UPDATE t SET a = 1", true},
		{"UPDATE db.T SET a = 1", true},
		{"UPDATE other.t SET a = 1", false},
		{"DELETE t1 FROM t1 JOIN t ON t1.a = t.a", true},
		{"UPDATE t1 JOIN (t2 JOIN t ON t2.c = t.a) ON t1.a = t2.c SET t1.b = 1", true},
		{"SELECT * FROM t", false},
	}
	for _, tc := range tests {
		got, err := analysis.ReferencesTable(tc.statement, "db", "t")
		require.NoError(t, err)
		require.Equal(t, tc.want, got, tc.statement)
	}
	_, err := analysis.ReferencesTable("UPDATE", "db", "t")
	require.Error(t, err)
}

func TestRestoreSQLSelfJoin(t *testing.T) {
	// Each alias of the table contributes the columns it sets; the other
	// table's do not.
	got, err := restore(t, "UPDATE t x JOIN t y ON x.a = y.a + 1 JOIN t1 ON t1.a = x.a SET x.b = 10, y.c = 20, t1.b = 30", "t")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got, "ON DUPLICATE KEY UPDATE `b` = VALUES(`b`), `c` = VALUES(`c`);"), got)

	// An alias of the same-named table in another database is not the table.
	got, err = restore(t, "UPDATE db.t x JOIN other.t y ON x.a = y.a SET x.b = 1, y.c = 2", "t")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got, "ON DUPLICATE KEY UPDATE `b` = VALUES(`b`);"), got)
}

func TestRestoreSQLEscapesCommentEnd(t *testing.T) {
	got, err := restore(t, "DELETE /* note */ FROM t WHERE a = '*/'", "t")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE /* note * / FROM t WHERE a = '* /'\n*/\nINSERT INTO `db`.`t` (`a`, `b`, `c`) SELECT `a`, `b`, `c` FROM `bak`.`pre_t`;", got)
}

func TestRestoreSQLSkipsKeysThatCanChange(t *testing.T) {
	// Every unique key of tk is on a generated column, an expression, or
	// nothing, so none identifies the row the UPDATE changed.
	_, err := restore(t, "UPDATE tk SET b = 1", "tk")
	require.EqualError(t, err, "no disjoint unique key found for db.tk")
}
