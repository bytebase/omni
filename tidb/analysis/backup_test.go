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
		{Name: "t1", Columns: backupColumns("a", "b"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"a"}}}},
		{Name: "t2", Columns: backupColumns("c", "d"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"c"}}}},
		{Name: "stamped", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "ts", OnUpdate: "CURRENT_TIMESTAMP"}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}},
		{Name: "stamped_key", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "ts", OnUpdate: "CURRENT_TIMESTAMP"}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}, {Name: "uk_ts", Unique: true, Expressions: []string{"ts"}}}},
		{Name: "swapped", Columns: backupColumns("id", "u", "v", "w"), Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}},
			{Name: "uk_u", Unique: true, Expressions: []string{"u"}},
			{Name: "uk_w", Unique: true, Expressions: []string{"(lower(`w`))"}},
		}},
		{Name: "parent", Columns: backupColumns("id", "v"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}},
		{Name: "child", Columns: backupColumns("id", "parent_id"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk_cascade", Columns: []string{"parent_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "SET NULL", OnUpdate: "CASCADE"}}},
		{Name: "nokey", Columns: []*metadata.ColumnMetadata{{Name: "a", Nullable: true}, {Name: "b"}}, Indexes: []*metadata.IndexMetadata{uniqueKey("a")}},
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
			Statement:   "CREATE TABLE `bak`.`pre_1_t` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_1_t` (`a`,`b`,`c`) SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE b = 2;",
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
	require.Equal(t, "CREATE TABLE `bak`.`pre_0_t1` LIKE `db`.`t1`;\nINSERT INTO `bak`.`pre_0_t1` (`a`,`b`) WITH c AS (SELECT c FROM t2) SELECT DISTINCT `X`.`a`,`X`.`b` FROM t1 AS X JOIN c ON x.a = c.c;", got[0].Statement)

	// More than five statements on one table share a backup.
	script := strings.Repeat("DELETE FROM t WHERE a = 1;\n", 6)
	got, err = transform(t, script)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "pre_t_db", got[0].BackupTable)
	require.Equal(t, 5, strings.Count(got[0].Statement, "UNION DISTINCT"))
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
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{Normal: []string{"x"}, RowKey: true}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, strings.Repeat("回", 64), got[0].BackupTable)

	_, err = analysis.TransformDMLToSelect("DELETE FROM t", analysis.BackupOptions{
		DefaultDatabase: "db",
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{}, errors.New("no columns")
		},
	})
	require.EqualError(t, err, "failed to classify columns: no columns")
}

func TestClassifyColumns(t *testing.T) {
	columns, err := analysis.ClassifyColumns(backupCatalog(), false, "Db", "tg")
	require.NoError(t, err)
	require.Equal(t, analysis.TableColumns{Generated: []string{"g_c"}, Normal: []string{"a", "b"}, RowKey: true, RowKeyColumns: []string{"g_c"}}, columns)

	_, err = analysis.ClassifyColumns(backupCatalog(), true, "db", "tg")
	var notFound *analysis.NotFoundError
	require.ErrorAs(t, err, &notFound)

	catalog := backupCatalog()
	catalog.ListDatabaseNames = func() ([]string, error) { return nil, errors.New("down") }
	_, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, "failed to list databases names: down")

	_, err = analysis.ClassifyColumns(analysis.Catalog{}, false, "db", "t")
	require.Error(t, err)
}

func TestTransformDMLToSelectKeepsEachLimit(t *testing.T) {
	script := "DELETE FROM t WHERE a > 1 ORDER BY a LIMIT 1;\n" + strings.Repeat("DELETE FROM t WHERE a = 1;\n", 5)
	got, err := transform(t, script)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Contains(t, got[0].Statement, "  (SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE a > 1 ORDER BY a LIMIT 1)\n  UNION DISTINCT\n", got[0].Statement)

	// A later statement's LIMIT may pick rows the backups, read before any
	// statement runs, never saw; separate backups as much as one.
	for _, script := range []string{
		"DELETE FROM t ORDER BY a LIMIT 1; DELETE FROM t ORDER BY a LIMIT 1",
		strings.Repeat("DELETE FROM t WHERE a = 1;\n", 5) + "UPDATE t SET b = 1 ORDER BY a LIMIT 1",
	} {
		_, err = transform(t, script)
		require.EqualError(t, err, "prior backup cannot handle a statement with LIMIT after another statement on the same table db.t", script)
	}
}

func TestTransformDMLToSelectNeedsRowKey(t *testing.T) {
	// nokey has no primary key and only a unique key on a nullable column.
	got, err := transform(t, "DELETE FROM nokey WHERE a = 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "SELECT `nokey`.`a`,`nokey`.`b` FROM nokey WHERE a = 1;"), got[0].Statement)
	_, err = transform(t, "DELETE nokey FROM nokey JOIN t ON nokey.a = t.a")
	require.EqualError(t, err, "prior backup cannot back up a table without a primary key or a unique key of non-null columns from a statement that joins other tables")
	_, err = transform(t, strings.Repeat("DELETE FROM nokey WHERE b = 1;\n", 6))
	require.EqualError(t, err, "prior backup cannot combine several statements on a table without a primary key or a unique key of non-null columns")
}

func TestTransformDMLToSelectQuotesIdentifiers(t *testing.T) {
	got, err := analysis.TransformDMLToSelect("DELETE FROM `a``b` WHERE `c``d` = 1", analysis.BackupOptions{
		DefaultDatabase: "d`b",
		BackupDatabase:  "b`ak",
		TablePrefix:     "p",
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{Normal: []string{"c`d"}, RowKey: true}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, "CREATE TABLE `b``ak`.`p_0_a``b` LIKE `d``b`.`a``b`;\nINSERT INTO `b``ak`.`p_0_a``b` (`c``d`) SELECT `a``b`.`c``d` FROM `a``b` WHERE `c``d` = 1;", got[0].Statement)
}

func TestTransformDMLToSelectResolvesEachDatabase(t *testing.T) {
	// The SET column is u's in db2's metadata, though db's t has one of that
	// name too.
	catalog := backupCatalog()
	catalog.GetDatabase = func(name string) (analysis.Database, error) {
		switch name {
		case "db":
			return backupDB, nil
		case "DB":
			return otherDB, nil
		}
		return nil, nil
	}
	got, err := analysis.TransformDMLToSelect("UPDATE DB.t JOIN DB.u ON DB.t.id = DB.u.id SET a = 1", analysis.BackupOptions{
		DefaultDatabase: "db",
		Catalog:         catalog,
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{Normal: []string{"id"}, RowKey: true}, nil
		},
	})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "u", got[0].SourceTable)
}

var otherDB = &backupDatabase{
	name: "DB",
	tables: []*metadata.TableMetadata{
		{Name: "t", Columns: backupColumns("id", "b")},
		{Name: "u", Columns: backupColumns("id", "a")},
	},
}

func TestTransformDMLToSelectCTEAlias(t *testing.T) {
	got, err := transform(t, "WITH c AS (SELECT 1 AS a) DELETE c FROM t AS c WHERE c.a = 1")
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "t", got[0].SourceTable)
}

func TestTransformDMLToSelectRefusesUnsafeScripts(t *testing.T) {
	long := strings.Repeat("x", 62)
	tests := []struct {
		statement string
		want      string
	}{
		{"UPDATE t1 SET b = 1 WHERE a = 1; DELETE FROM t WHERE a IN (SELECT a FROM t1 WHERE b = 1)", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		{"INSERT INTO t1 VALUES (1, 1); UPDATE t JOIN t1 ON t.a = t1.a SET t.b = 1", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		{"DELETE FROM parent WHERE id = 1", "prior backup cannot back up db.parent: foreign key fk_cascade of child is ON DELETE SET NULL, and the rows it changes are not backed up"},
		{"UPDATE parent SET id = 2 WHERE id = 1", "prior backup cannot back up db.parent: foreign key fk_cascade of child is ON UPDATE CASCADE, and the rows it changes are not backed up"},
		// Separate backups both hold a row a DELETE removes.
		{"UPDATE t SET b = 1 WHERE a = 1; DELETE FROM t WHERE a = 1", "prior backup cannot back up a DELETE and another statement on the same table db.t separately"},
		// Both targets of one statement share its offset in their names.
		{"DELETE a, b FROM `" + long + "1` a JOIN `" + long + "2` b ON a.a = b.a", "prior backup cannot name the backups of db." + long + "1 and db." + long + "2 apart: both are pre_0_" + strings.Repeat("x", 58)},
	}
	for _, tc := range tests {
		_, err := analysis.TransformDMLToSelect(tc.statement, analysis.BackupOptions{
			DefaultDatabase: "db",
			BackupDatabase:  "bak",
			TablePrefix:     "pre",
			Catalog:         backupCatalog(),
			Columns: func(database, table string) (analysis.TableColumns, error) {
				if strings.HasPrefix(table, "xxx") {
					return analysis.TableColumns{Normal: []string{"a"}, RowKey: true, RowKeyColumns: []string{"a"}}, nil
				}
				return analysis.ClassifyColumns(backupCatalog(), false, database, table)
			},
		})
		require.EqualError(t, err, tc.want, tc.statement)
	}
}

func TestTransformDMLToSelectKeepsSafeScripts(t *testing.T) {
	for _, script := range []string{
		"UPDATE t SET b = 1 WHERE a = 1; UPDATE t SET c = 2 WHERE b = 1",
		"UPDATE t1 SET b = 1; DELETE FROM t WHERE a = 1",
		"UPDATE parent SET v = 1",
		"DELETE FROM t WHERE a = 1; DELETE FROM t1 WHERE a = 1",
	} {
		_, err := transform(t, script)
		require.NoError(t, err, script)
	}
}

func TestTransformDMLToSelectSkipsOuterJoinRows(t *testing.T) {
	got, err := transform(t, "DELETE t FROM t1 LEFT JOIN t ON t.a = t1.a WHERE t1.b = 1 OR t1.b = 2")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM t1 LEFT JOIN t ON t.a = t1.a WHERE (t1.b = 1 OR t1.b = 2) AND `t`.`a` IS NOT NULL;"), got[0].Statement)
	got, err = transform(t, "DELETE t1 FROM t1 LEFT JOIN t ON t.a = t1.a")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM t1 LEFT JOIN t ON t.a = t1.a;"), got[0].Statement)
}
