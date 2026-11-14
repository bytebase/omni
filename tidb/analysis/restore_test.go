package analysis_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/tidb/analysis"
)

func restore(t *testing.T, statement, table string) (string, error) {
	t.Helper()
	columnSet, err := analysis.NormalColumnSet(backupCatalog(), "db", table)
	require.NoError(t, err)
	nodes, err := analysis.MatchingDMLs(statement, "db", table, columnSet)
	require.NoError(t, err)
	require.NotEmpty(t, nodes)
	columns, err := analysis.ClassifyColumns(backupCatalog(), false, "db", table)
	require.NoError(t, err)
	return analysis.RestoreSQL(nodes, analysis.RestoreOptions{
		Statement:      statement,
		Database:       "db",
		Table:          table,
		BackupDatabase: "bak",
		BackupTable:    "pre_" + table,
		Columns:        columns.Normal,
		Catalog:        backupCatalog(),
	})
}

func TestRestoreSQL(t *testing.T) {
	got, err := restore(t, "DELETE FROM tg WHERE a = 1", "tg")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM tg WHERE a = 1\n*/\nINSERT INTO `db`.`tg` (`a`, `b`) SELECT `a`, `b` FROM `bak`.`pre_tg`;", got)

	// Every UPDATE of the bundle contributes its columns, once each.
	got, err = restore(t, "UPDATE t SET b = 1; DELETE FROM t; UPDATE t X JOIN t1 ON x.a = t1.a SET x.C = 1, B = 2, t1.b = 3", "t")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nUPDATE t SET b = 1; DELETE FROM t; UPDATE t X JOIN t1 ON x.a = t1.a SET x.C = 1, B = 2, t1.b = 3\n*/\nINSERT INTO `db`.`t` (`a`, `b`, `c`) SELECT `a`, `b`, `c` FROM `bak`.`pre_t` ON DUPLICATE KEY UPDATE `b` = VALUES(`b`), `C` = VALUES(`C`);", got)

	// A unique key on a generated column, or with no columns, never counts.
	_, err = restore(t, "UPDATE tg SET b = 1", "tg")
	require.EqualError(t, err, "no disjoint unique key found for db.tg")
	_, err = restore(t, "UPDATE tg SET a = 1", "tg")
	require.NoError(t, err)
	// g_c reads a, so an expression on a can move uk(g_c)'s values.
	_, err = restore(t, "UPDATE tg SET a = a + 1", "tg")
	require.EqualError(t, err, "prior backup cannot restore db.tg: the UPDATE changes unique key uk, so the restore could match another row")
}

func TestNormalColumnSet(t *testing.T) {
	got, err := analysis.NormalColumnSet(backupCatalog(), "db", "tg")
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"a": true, "b": true}, got)

	_, err = analysis.NormalColumnSet(backupCatalog(), "nodb", "t")
	require.EqualError(t, err, "database metadata is nil for nodb")
	_, err = analysis.NormalColumnSet(backupCatalog(), "db", "missing")
	require.EqualError(t, err, "table metadata is nil for db.missing")
	catalog := backupCatalog()
	catalog.GetDatabase = func(string) (analysis.Database, error) { return nil, errors.New("down") }
	_, err = analysis.NormalColumnSet(catalog, "db", "t")
	require.EqualError(t, err, "failed to get database metadata for db: down")
}

func TestMutatesTable(t *testing.T) {
	columns := map[string]bool{"a": true, "b": true, "c": true}
	tests := []struct {
		statement string
		want      bool
	}{
		{"UPDATE t SET a = 1", true},
		{"UPDATE DB.T SET a = 1", true},
		{"UPDATE t JOIN t2 ON t.a = t2.c SET d = 1", false},
		{"UPDATE t JOIN t2 ON t.a = t2.c SET t2.d = 1", false},
		{"UPDATE t x JOIN t2 ON x.a = t2.c SET X.b = 1", true},
		{"UPDATE db.t x JOIN other.t y ON x.a = y.a SET y.b = 1", false},
		{"DELETE t2 FROM t2 JOIN t ON t2.c = t.a", false},
		{"DELETE x FROM t2 JOIN t AS x ON t2.c = x.a", true},
		{"DELETE FROM t", true},
		{"SELECT * FROM t", false},
	}
	for _, tc := range tests {
		got, err := analysis.MutatesTable(tc.statement, "db", "t", columns)
		require.NoError(t, err)
		require.Equal(t, tc.want, got, tc.statement)
	}
	_, err := analysis.MutatesTable("UPDATE", "db", "t", columns)
	require.Error(t, err)

	nodes, err := analysis.MatchingDMLs("UPDATE t SET a = 1; UPDATE t1 SET a = 1; DELETE FROM t", "db", "t", columns)
	require.NoError(t, err)
	require.Len(t, nodes, 2)
}

func TestRestoreSQLSelfJoin(t *testing.T) {
	got, err := restore(t, "UPDATE t x JOIN t y ON x.a = y.a + 1 SET x.b = 10, y.c = 20", "t")
	require.NoError(t, err)
	require.Contains(t, got, "ON DUPLICATE KEY UPDATE `b` = VALUES(`b`), `c` = VALUES(`c`);")
}

func TestRestoreSQLEscapesCommentEnd(t *testing.T) {
	got, err := restore(t, "DELETE /* note */ FROM t WHERE a = '*/'", "t")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE /* note * / FROM t WHERE a = '* /'\n*/\nINSERT INTO `db`.`t` (`a`, `b`, `c`) SELECT `a`, `b`, `c` FROM `bak`.`pre_t`;", got)
}

func TestRestoreSQLSkipsNullableKeys(t *testing.T) {
	_, err := restore(t, "UPDATE nokey SET b = 1", "nokey")
	require.EqualError(t, err, "no disjoint unique key found for db.nokey")
}

func TestRestoreSQLQuotesIdentifiers(t *testing.T) {
	nodes, err := analysis.MatchingDMLs("DELETE FROM `a``b`", "d`b", "a`b", nil)
	require.NoError(t, err)
	got, err := analysis.RestoreSQL(nodes, analysis.RestoreOptions{
		Statement:      "DELETE FROM `a``b`",
		Database:       "d`b",
		Table:          "a`b",
		BackupDatabase: "b`ak",
		BackupTable:    "p`t",
		Columns:        []string{"c`d", "e"},
	})
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM `a``b`\n*/\nINSERT INTO `d``b`.`a``b` (`c``d`, `e`) SELECT `c``d`, `e` FROM `b``ak`.`p``t`;", got)
}

func TestRestoreSQLRestoresOnUpdateColumns(t *testing.T) {
	got, err := restore(t, "UPDATE stamped SET v = 2 WHERE id = 1", "stamped")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got, "ON DUPLICATE KEY UPDATE `v` = VALUES(`v`), `ts` = VALUES(`ts`);"), got)
}

func TestRestoreSQLRefusesChangedUniqueKeys(t *testing.T) {
	// A key's old value can sit in another row once the UPDATEs move values
	// between rows: by an expression, or across statements.
	for _, statement := range []string{
		"UPDATE swapped SET u = u + 1",
		"UPDATE swapped SET u = 30 WHERE id = 1; UPDATE swapped SET u = 10 WHERE id = 2",
	} {
		_, err := restore(t, statement, "swapped")
		require.EqualError(t, err, "prior backup cannot restore db.swapped: the UPDATE changes unique key uk_u, so the restore could match another row", statement)
	}
	_, err := restore(t, "UPDATE swapped SET w = CONCAT(w, 'x')", "swapped")
	require.EqualError(t, err, "prior backup cannot restore db.swapped: the UPDATE changes unique key uk_w, so the restore could match another row")
	// The server stamps ts on every changed row, which no literal pins.
	_, err = restore(t, "UPDATE stamped_key SET v = 1 WHERE id = 1", "stamped_key")
	require.EqualError(t, err, "prior backup cannot restore db.stamped_key: the UPDATE changes unique key uk_ts, so the restore could match another row")
	// One statement setting a key to a literal leaves the old value free.
	for _, statement := range []string{
		"UPDATE swapped SET u = 30 WHERE id = 1",
		"UPDATE swapped SET w = 'x' WHERE id = 1",
		"UPDATE swapped SET v = 1",
	} {
		_, err := restore(t, statement, "swapped")
		require.NoError(t, err, statement)
	}
}
