package analysis_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/mariadb/analysis"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
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

var backupDB = &backupDatabase{
	name: "db",
	tables: []*metadata.TableMetadata{
		{Name: "t", Columns: backupColumns("a", "b", "c"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}}},
		{Name: "t1", Columns: backupColumns("a", "b")},
		{Name: "t2", Columns: backupColumns("c", "d")},
		{Name: "tg", Columns: backupColumns("a", "b", "g_c"), Indexes: []*metadata.IndexMetadata{{Name: "uk", Unique: true, Expressions: []string{"B"}}}},
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
	tests := []struct {
		statement string
		want      []analysis.BackupStatement
	}{
		{
			statement: "UPDATE t SET a = 1 WHERE b = 2 ORDER BY c LIMIT 3;",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db`  SELECT `t`.* FROM t WHERE b = 2 ORDER BY c LIMIT 3;",
				SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 0, ByteEnd: 49,
			}},
		},
		{
			// Two statements on one table share a backup; the range spans both.
			statement: "SELECT 1; DELETE FROM t WHERE a = 1; DELETE FROM db.t WHERE a = 2;",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db`  SELECT `t`.* FROM t WHERE a = 1\n  UNION DISTINCT\n  SELECT `t`.* FROM db.t WHERE a = 2;",
				SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 9, ByteEnd: 65,
			}},
		},
		{
			// Generated columns are left out of the copy.
			statement: "DELETE x FROM tg AS x JOIN t1 ON x.a = t1.a",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_tg_db` LIKE `db`.`tg`;\nINSERT INTO `bak`.`pre_tg_db` (`a`,`b`)  SELECT `x`.`a`,`x`.`b` FROM tg AS x JOIN t1 ON x.a = t1.a;",
				SourceTable: "tg", BackupTable: "pre_tg_db", ByteStart: 0, ByteEnd: 43,
			}},
		},
		{
			// The CTE is kept, a CTE target is not, and an unqualified SET column
			// goes to the table that has it.
			statement: "WITH c AS (SELECT a FROM t2) UPDATE t1 JOIN t2 ON t1.a = t2.c JOIN c ON c.a = t1.a SET d = 1, c.a = 2",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_t2_db` LIKE `db`.`t2`;\nINSERT INTO `bak`.`pre_t2_db`  WITH c AS (SELECT a FROM t2) SELECT `t2`.* FROM t1 JOIN t2 ON t1.a = t2.c JOIN c ON c.a = t1.a;",
				SourceTable: "t2", BackupTable: "pre_t2_db", ByteStart: 0, ByteEnd: 101,
			}},
		},
		{
			statement: "UPDATE t SET a = 1; UPDATE t1 SET a = 1",
			want: []analysis.BackupStatement{
				{
					Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db`  SELECT `t`.* FROM t;",
					SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 0, ByteEnd: 18,
				},
				{
					Statement:   "CREATE TABLE `bak`.`pre_t1_db` LIKE `db`.`t1`;\nINSERT INTO `bak`.`pre_t1_db`  SELECT `t1`.* FROM t1;",
					SourceTable: "t1", BackupTable: "pre_t1_db", ByteStart: 19, ByteEnd: 39,
				},
			},
		},
		{statement: "INSERT INTO t VALUES (1); SELECT 1"},
	}
	for _, tc := range tests {
		got, err := transform(t, tc.statement)
		require.NoError(t, err, tc.statement)
		require.Equal(t, tc.want, got, tc.statement)
	}
}

func TestTransformDMLToSelectErrors(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"UPDATE t SET", "failed to prepare transformation: failed to parse sql: "},
		{"UPDATE t1 JOIN t2 SET t9.a = 1", `failed to prepare transformation: failed to extract tables: cannot extract reference table: no matched updated table "t9" in referenced table list`},
		{"UPDATE t SET a = 1; DELETE FROM t", "prior backup cannot handle mixed DML statements on the same table db.t"},
		{"DELETE FROM nodb.t", "failed to generate SQL for table nodb.t: failed to classify columns: resource not found: database: nodb"},
		{"DELETE FROM missing", "failed to generate SQL for table db.missing: failed to classify columns: table missing not found in schema"},
	}
	for _, tc := range tests {
		_, err := transform(t, tc.statement)
		require.Error(t, err, tc.statement)
		require.True(t, strings.HasPrefix(err.Error(), tc.want), "%s: %v", tc.statement, err)
	}
}

func TestTransformDMLToSelectOptions(t *testing.T) {
	var asked []string
	got, err := analysis.TransformDMLToSelect("UPDATE `T` SET a = 1", analysis.BackupOptions{
		DefaultDatabase: "db",
		BackupDatabase:  "bak",
		TablePrefix:     strings.Repeat("回", 70),
		Columns: func(database, table string) ([]string, []string, error) {
			asked = append(asked, database+"."+table)
			return []string{"g"}, []string{"x"}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"db.T"}, asked)
	// The backup table's name is cut to 64 characters, not bytes.
	require.Equal(t, strings.Repeat("回", 64), got[0].BackupTable)
	require.Contains(t, got[0].Statement, "(`x`)  SELECT `T`.`x` FROM `T`;")

	_, err = analysis.TransformDMLToSelect("DELETE FROM t", analysis.BackupOptions{
		DefaultDatabase: "db",
		Columns: func(string, string) ([]string, []string, error) {
			return nil, nil, errors.New("no columns")
		},
	})
	require.EqualError(t, err, "failed to generate SQL for table db.t: failed to classify columns: no columns")
}

func TestClassifyColumns(t *testing.T) {
	generated, normal, err := analysis.ClassifyColumns(backupCatalog(), false, "Db", "TG")
	require.NoError(t, err)
	require.Equal(t, []string{"g_c"}, generated)
	require.Equal(t, []string{"a", "b"}, normal)

	// Case-sensitive lookups find neither the listed "DB" nor the table "TG".
	_, _, err = analysis.ClassifyColumns(backupCatalog(), true, "db", "tg")
	var notFound *analysis.NotFoundError
	require.ErrorAs(t, err, &notFound)
	_, _, err = analysis.ClassifyColumns(backupCatalog(), true, "DB", "TG")
	require.EqualError(t, err, "table TG not found in schema")

	catalog := backupCatalog()
	catalog.ListDatabaseNames = func() ([]string, error) { return nil, errors.New("down") }
	_, _, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, "failed to list databases names: down")

	catalog = backupCatalog()
	catalog.GetDatabase = func(string) (analysis.Database, error) { return nil, errors.New("down") }
	_, _, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, `failed to get database metadata for database "DB": down`)

	_, _, err = analysis.ClassifyColumns(analysis.Catalog{}, false, "db", "t")
	require.Error(t, err)
}

func TestExtractTables(t *testing.T) {
	statement := "UPDATE t1 a JOIN t2 b ON a.a = b.c SET d = 1"
	list, err := parser.Parse(statement)
	require.NoError(t, err)
	// Without metadata an unqualified SET column could be any table's.
	got, err := analysis.ExtractTables("db", list.Items[0], statement, 3, nil)
	require.NoError(t, err)
	require.Len(t, got, 2)
	got, err = analysis.ExtractTables("db", list.Items[0], statement, 3, backupDB)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, analysis.TableReference{Database: "db", Table: "t2", Alias: "b", StatementType: analysis.DMLTypeUpdate}, *got[0].Table)
	require.Equal(t, 3, got[0].Offset)
	require.Equal(t, statement, got[0].Statement)
}

func TestTransformDMLToSelectReturning(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"DELETE FROM t WHERE a = 1 RETURNING a", "SELECT `t`.* FROM t WHERE a = 1;"},
		{"DELETE FROM t WHERE a = 'RETURNING' ORDER BY a LIMIT 1 RETURNING /* RETURNING */ *", "SELECT `t`.* FROM t WHERE a = 'RETURNING' ORDER BY a LIMIT 1;"},
		{"DELETE FROM t RETURNING `returning`, b + 1 AS x", "SELECT `t`.* FROM t;"},
	}
	for _, tc := range tests {
		got, err := transform(t, tc.statement)
		require.NoError(t, err, tc.statement)
		require.Len(t, got, 1)
		require.True(t, strings.HasSuffix(got[0].Statement, tc.want), "%s: %s", tc.statement, got[0].Statement)
	}
}
