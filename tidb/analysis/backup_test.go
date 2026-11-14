package analysis_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/analysis"
)

// backupDatabase is an in-memory Database whose table names match exactly.
type backupDatabase struct {
	name   string
	tables []*metadata.TableMetadata
}

func (d *backupDatabase) Name() string          { return d.name }
func (*backupDatabase) HasSchema() bool         { return true }
func (*backupDatabase) ListViewNames() []string { return nil }

func (*backupDatabase) GetView(string) *metadata.ViewMetadata { return nil }

func (d *backupDatabase) ListTableNames() []string {
	var names []string
	for _, table := range d.tables {
		names = append(names, table.Name)
	}
	return names
}

func (d *backupDatabase) GetTable(name string) *metadata.TableMetadata {
	for _, table := range d.tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func backupColumns(names ...string) []*metadata.ColumnMetadata {
	var columns []*metadata.ColumnMetadata
	for _, name := range names {
		column := &metadata.ColumnMetadata{Name: name}
		if strings.HasPrefix(name, "g_") {
			column.Generation = &metadata.GenerationMetadata{Expression: "a + b"}
		}
		columns = append(columns, column)
	}
	return columns
}

func uniqueKey(columns ...string) *metadata.IndexMetadata {
	return &metadata.IndexMetadata{Name: "uk", Unique: true, Expressions: columns}
}

var backupDB = &backupDatabase{
	name: "db",
	tables: []*metadata.TableMetadata{
		{Name: "t", Columns: backupColumns("a", "b", "c"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"a"}}}},
		{Name: "t1", Columns: backupColumns("a", "b")},
		{Name: "t2", Columns: backupColumns("c", "d")},
		{Name: "tg", Columns: backupColumns("a", "b", "g_c"), Indexes: []*metadata.IndexMetadata{uniqueKey("g_c"), uniqueKey(), uniqueKey("B")}},
	},
}

func backupCatalog() analysis.Catalog {
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) { return []string{"other", "DB"}, nil },
		GetDatabase: func(name string) (analysis.Database, error) {
			if strings.EqualFold(name, "db") {
				return backupDB, nil
			}
			return nil, nil
		},
	}
}

func transform(t *testing.T, statement string) ([]analysis.BackupStatement, error) {
	t.Helper()
	return analysis.TransformDMLToSelect(statement, analysis.BackupOptions{
		DefaultDatabase: "db",
		BackupDatabase:  "bak",
		TablePrefix:     "pre",
		Catalog:         backupCatalog(),
	})
}

func TestTransformDMLToSelect(t *testing.T) {
	// Up to five statements get a backup each, named by statement index.
	got, err := transform(t, "SELECT 1; UPDATE t SET a = 1 WHERE b = 2; DELETE FROM tg WHERE a = 1;")
	require.NoError(t, err)
	require.Equal(t, []analysis.BackupStatement{
		{
			Statement:   "CREATE TABLE `bak`.`pre_1_t` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_1_t` SELECT `t`.* FROM t WHERE b = 2;",
			SourceTable: "t", BackupTable: "pre_1_t", ByteStart: 9, ByteEnd: 40,
		},
		{
			Statement:   "CREATE TABLE `bak`.`pre_2_tg` LIKE `db`.`tg`;\nINSERT INTO `bak`.`pre_2_tg` (`a`,`b`) SELECT `tg`.`a`,`tg`.`b` FROM tg WHERE a = 1;",
			SourceTable: "tg", BackupTable: "pre_2_tg", ByteStart: 41, ByteEnd: 68,
		},
	}, got)

	// A CTE stays in front of its SELECT; qualifiers resolve case-insensitively.
	got, err = transform(t, "WITH c AS (SELECT c FROM t2) UPDATE t1 AS X JOIN c ON x.a = c.c SET x.b = 1")
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "CREATE TABLE `bak`.`pre_0_t1` LIKE `db`.`t1`;\nINSERT INTO `bak`.`pre_0_t1` WITH c AS (SELECT c FROM t2) SELECT `X`.* FROM t1 AS X JOIN c ON x.a = c.c;", got[0].Statement)

	// More than five statements on one table share a backup.
	script := strings.Repeat("DELETE FROM t WHERE a = 1;\n", 6)
	got, err = transform(t, script)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "pre_t_db", got[0].BackupTable)
	require.Equal(t, 5, strings.Count(got[0].Statement, "UNION ALL"))
	require.Equal(t, 0, got[0].ByteStart)
	require.Equal(t, len(script)-2, got[0].ByteEnd)

	// A parenthesized join operand gets its parentheses back.
	got, err = transform(t, "UPDATE (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a SET t1.b = 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a;"), got[0].Statement)

	got, err = transform(t, "INSERT INTO t VALUES (1); EXPLAIN UPDATE t SET a = 1")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestTransformDMLToSelectErrors(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"UPDATE t SET", "failed to prepare transformation: failed to parse sql: "},
		{"UPDATE other.t SET a = 1", `failed to prepare transformation: prior backup does not support cross-database mutations: other.t (task database "db")`},
		{"REPLACE INTO t VALUES (1)", "failed to prepare transformation: failed to extract tables: prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements"},
		{"BATCH ON a LIMIT 10 DELETE FROM t", "failed to prepare transformation: failed to extract tables: prior backup does not support TiDB BATCH (non-transactional DML) statements"},
		{"EXPLAIN ANALYZE DELETE FROM t", "failed to prepare transformation: failed to extract tables: prior backup does not support EXPLAIN ANALYZE of an UPDATE/DELETE/BATCH statement"},
		{"WITH t1 AS (SELECT 1 AS a) DELETE t1 FROM t1", "failed to prepare transformation: failed to extract tables: prior backup could not determine a backup target for the DELETE statement"},
		{strings.Repeat("DELETE FROM t;", 5) + "DELETE FROM t1;", "prior backup cannot handle statements on different tables more than 5"},
		{strings.Repeat("WITH c AS (SELECT 1) DELETE FROM t;", 6), "prior backup does not support WITH (CTE) statements when more than 5 DML statements target the same table"},
		{"DELETE FROM missing", `failed to classify columns: failed to get table metadata for table "missing"`},
	}
	for _, tc := range tests {
		_, err := transform(t, tc.statement)
		require.Error(t, err, tc.statement)
		require.True(t, strings.HasPrefix(err.Error(), tc.want), "%s: %v", tc.statement, err)
	}
}

func TestTransformDMLToSelectOptions(t *testing.T) {
	_, err := analysis.TransformDMLToSelect("UPDATE DB.t SET a = 1", analysis.BackupOptions{DefaultDatabase: "db", CaseSensitive: true, Catalog: backupCatalog()})
	require.EqualError(t, err, `failed to prepare transformation: prior backup does not support cross-database mutations: DB.t (task database "db")`)

	got, err := analysis.TransformDMLToSelect("DELETE FROM t", analysis.BackupOptions{
		DefaultDatabase: "db",
		TablePrefix:     strings.Repeat("回", 70),
		Columns: func(string, string) ([]string, []string, error) {
			return nil, []string{"x"}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, strings.Repeat("回", 64), got[0].BackupTable)

	_, err = analysis.TransformDMLToSelect("DELETE FROM t", analysis.BackupOptions{
		DefaultDatabase: "db",
		Columns: func(string, string) ([]string, []string, error) {
			return nil, nil, errors.New("no columns")
		},
	})
	require.EqualError(t, err, "failed to classify columns: no columns")
}

func TestClassifyColumns(t *testing.T) {
	generated, normal, err := analysis.ClassifyColumns(backupCatalog(), false, "Db", "tg")
	require.NoError(t, err)
	require.Equal(t, []string{"g_c"}, generated)
	require.Equal(t, []string{"a", "b"}, normal)

	_, _, err = analysis.ClassifyColumns(backupCatalog(), true, "db", "tg")
	var notFound *analysis.NotFoundError
	require.ErrorAs(t, err, &notFound)

	catalog := backupCatalog()
	catalog.ListDatabaseNames = func() ([]string, error) { return nil, errors.New("down") }
	_, _, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, "failed to list databases names: down")

	_, _, err = analysis.ClassifyColumns(analysis.Catalog{}, false, "db", "t")
	require.Error(t, err)
}
