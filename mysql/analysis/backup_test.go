package analysis_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/parser"
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
		{Name: "t1", Columns: backupColumns("a", "b"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}}},
		{Name: "t2", Columns: backupColumns("c", "d"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"c"}}}},
		{Name: "tg", Columns: backupColumns("a", "b", "g_c"), Indexes: []*metadata.IndexMetadata{{Name: "uk", Unique: true, Expressions: []string{"B"}}}},
		{Name: "stamped", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "ts", OnUpdate: "CURRENT_TIMESTAMP"}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}},
		{Name: "stamped_key", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "ts", OnUpdate: "CURRENT_TIMESTAMP"}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}, {Name: "uk_ts", Unique: true, Expressions: []string{"ts"}}}},
		{Name: "swapped", Columns: backupColumns("id", "u", "v", "w"), Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}},
			{Name: "uk_u", Unique: true, Expressions: []string{"u"}},
			{Name: "uk_w", Unique: true, Expressions: []string{"(lower(`w`))"}},
		}},
		{Name: "parent", Columns: backupColumns("id", "v"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}},
		{Name: "child", Columns: backupColumns("id", "parent_id"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk_cascade", Columns: []string{"parent_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "CASCADE", OnUpdate: "RESTRICT"}}},
		{Name: "bumped", Columns: backupColumns("id", "v"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, Triggers: []*metadata.TriggerMetadata{{Name: "tr_bump", Event: "UPDATE", Timing: "BEFORE"}}},
		{Name: "logged", Columns: backupColumns("id", "v"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, Triggers: []*metadata.TriggerMetadata{{Name: "tr_log", Event: "delete", Timing: "AFTER"}}},
		{Name: "stamped_insert", Columns: backupColumns("id", "v"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, Triggers: []*metadata.TriggerMetadata{{Name: "tr_insert", Event: "INSERT", Timing: "BEFORE"}}},
		{Name: "cased", Columns: backupColumns("a"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"a"}}}},
		{Name: "CASED", Columns: backupColumns("a"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"a"}}}},
		{Name: "stamped_parent", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "v"}, {Name: "ts", OnUpdate: "CURRENT_TIMESTAMP"}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}, {Name: "uk_ts", Unique: true, Expressions: []string{"ts"}}}},
		{Name: "stamped_child", Columns: backupColumns("id", "parent_ts"), Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}}, ForeignKeys: []*metadata.ForeignKeyMetadata{{Name: "fk_ts", Columns: []string{"parent_ts"}, ReferencedTable: "stamped_parent", ReferencedColumns: []string{"ts"}, OnDelete: "RESTRICT", OnUpdate: "CASCADE"}}},
		{Name: "summed", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "a"}, {Name: "b"}, {Name: "g", Generation: &metadata.GenerationMetadata{Expression: "`a` + `b`"}}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}, {Name: "uk_g", Unique: true, Expressions: []string{"g"}}}},
		{Name: "chained", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "a"}, {Name: "g1", Generation: &metadata.GenerationMetadata{Expression: "`a` + 1"}}, {Name: "g2", Generation: &metadata.GenerationMetadata{Expression: "`g1` + 1"}}}, Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Expressions: []string{"id"}}, {Name: "uk_g2", Unique: true, Expressions: []string{"g2"}}}},
		{Name: "nokey", Columns: []*metadata.ColumnMetadata{{Name: "a", Nullable: true}, {Name: "b"}}, Indexes: []*metadata.IndexMetadata{{Name: "uk", Unique: true, Expressions: []string{"a"}}}},
		{Name: "tk", Columns: backupColumns("a", "b", "g_c"), Indexes: []*metadata.IndexMetadata{
			{Name: "uk_generated", Unique: true, Expressions: []string{"g_c"}},
			{Name: "uk_expression", Unique: true, Expressions: []string{"(`a` + 1)"}},
			{Name: "uk_empty", Unique: true},
		}},
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
				Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db` (`a`,`b`,`c`)  SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE b = 2 ORDER BY c LIMIT 3;",
				SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 0, ByteEnd: 49,
			}},
		},
		{
			// Two statements on one table share a backup; the range spans both.
			statement: "SELECT 1; DELETE FROM t WHERE a = 1; DELETE FROM db.t WHERE a = 2;",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db` (`a`,`b`,`c`)  (SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE a = 1)\n  UNION DISTINCT\n  (SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM db.t WHERE a = 2);",
				SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 9, ByteEnd: 65,
			}},
		},
		{
			// Generated columns are left out of the copy, and a row a join
			// repeats is copied once.
			statement: "DELETE x FROM tg AS x JOIN t1 ON x.a = t1.a",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_tg_db` LIKE `db`.`tg`;\nINSERT INTO `bak`.`pre_tg_db` (`a`,`b`)  SELECT DISTINCT `x`.`a`,`x`.`b` FROM tg AS x JOIN t1 ON x.a = t1.a;",
				SourceTable: "tg", BackupTable: "pre_tg_db", ByteStart: 0, ByteEnd: 43,
			}},
		},
		{
			// The CTE is kept, a CTE target is not, and an unqualified SET column
			// goes to the table that has it.
			statement: "WITH c AS (SELECT a FROM t2) UPDATE t1 JOIN t2 ON t1.a = t2.c JOIN c ON c.a = t1.a SET d = 1, c.a = 2",
			want: []analysis.BackupStatement{{
				Statement:   "CREATE TABLE `bak`.`pre_t2_db` LIKE `db`.`t2`;\nINSERT INTO `bak`.`pre_t2_db` (`c`,`d`)  WITH c AS (SELECT a FROM t2) SELECT DISTINCT `t2`.`c`,`t2`.`d` FROM t1 JOIN t2 ON t1.a = t2.c JOIN c ON c.a = t1.a;",
				SourceTable: "t2", BackupTable: "pre_t2_db", ByteStart: 0, ByteEnd: 101,
			}},
		},
		{
			statement: "UPDATE t SET a = 1; UPDATE t1 SET a = 1",
			want: []analysis.BackupStatement{
				{
					Statement:   "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db` (`a`,`b`,`c`)  SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t;",
					SourceTable: "t", BackupTable: "pre_t_db", ByteStart: 0, ByteEnd: 18,
				},
				{
					Statement:   "CREATE TABLE `bak`.`pre_t1_db` LIKE `db`.`t1`;\nINSERT INTO `bak`.`pre_t1_db` (`a`,`b`)  SELECT `t1`.`a`,`t1`.`b` FROM t1;",
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
		{"DELETE t9 FROM t1 JOIN t2", `failed to prepare transformation: failed to extract tables: cannot extract reference table: no matched table "t9" in referenced table list`},
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
		Columns: func(database, table string) (analysis.TableColumns, error) {
			asked = append(asked, database+"."+table)
			return analysis.TableColumns{Generated: []string{"g"}, Normal: []string{"x"}}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"db.T"}, asked)
	// The backup table's name is cut to 64 characters, not bytes.
	require.Equal(t, strings.Repeat("回", 64), got[0].BackupTable)
	require.Contains(t, got[0].Statement, "(`x`)  SELECT `T`.`x` FROM `T`;")

	_, err = analysis.TransformDMLToSelect("DELETE FROM t", analysis.BackupOptions{
		DefaultDatabase: "db",
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{}, errors.New("no columns")
		},
	})
	require.EqualError(t, err, "failed to generate SQL for table db.t: failed to classify columns: no columns")
}

func TestClassifyColumns(t *testing.T) {
	columns, err := analysis.ClassifyColumns(backupCatalog(), false, "Db", "TG")
	require.NoError(t, err)
	require.Equal(t, analysis.TableColumns{Generated: []string{"g_c"}, Normal: []string{"a", "b"}, RowKey: true, RowKeyColumns: []string{"B"}}, columns)

	// Case-sensitive lookups find neither the listed "DB" nor the table "TG".
	_, err = analysis.ClassifyColumns(backupCatalog(), true, "db", "tg")
	var notFound *analysis.NotFoundError
	require.ErrorAs(t, err, &notFound)
	_, err = analysis.ClassifyColumns(backupCatalog(), true, "DB", "TG")
	require.EqualError(t, err, "table TG not found in schema")

	catalog := backupCatalog()
	catalog.ListDatabaseNames = func() ([]string, error) { return nil, errors.New("down") }
	_, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, "failed to list databases names: down")

	catalog = backupCatalog()
	catalog.GetDatabase = func(string) (analysis.Database, error) { return nil, errors.New("down") }
	_, err = analysis.ClassifyColumns(catalog, false, "db", "t")
	require.EqualError(t, err, `failed to get database metadata for database "DB": down`)

	_, err = analysis.ClassifyColumns(analysis.Catalog{}, false, "db", "t")
	require.Error(t, err)
}

func TestExtractTables(t *testing.T) {
	statement := "UPDATE t1 a JOIN t2 b ON a.a = b.c SET d = 1"
	list, err := parser.Parse(statement)
	require.NoError(t, err)
	// Without metadata an unqualified SET column could be any table's.
	got, err := analysis.ExtractTables("db", list.Items[0], statement, 3, analysis.Catalog{})
	require.NoError(t, err)
	require.Len(t, got, 2)
	got, err = analysis.ExtractTables("db", list.Items[0], statement, 3, backupCatalog())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, analysis.TableReference{Database: "db", Table: "t2", Alias: "b", StatementType: analysis.DMLTypeUpdate}, *got[0].Table)
	require.Equal(t, 3, got[0].Offset)
	require.Equal(t, statement, got[0].Statement)
}

func TestTransformDMLToSelectKeepsEachLimit(t *testing.T) {
	got, err := transform(t, "DELETE FROM t WHERE a > 2 ORDER BY b LIMIT 2; DELETE FROM t WHERE a = 1")
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "CREATE TABLE `bak`.`pre_t_db` LIKE `db`.`t`;\nINSERT INTO `bak`.`pre_t_db` (`a`,`b`,`c`)  (SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE a > 2 ORDER BY b LIMIT 2)\n  UNION DISTINCT\n  (SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE a = 1);", got[0].Statement)

	// A single statement's SELECT needs no parentheses.
	got, err = transform(t, "DELETE FROM t WHERE a = 1 LIMIT 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "  SELECT `t`.`a`,`t`.`b`,`t`.`c` FROM t WHERE a = 1 LIMIT 1;"), got[0].Statement)

	// A later statement's LIMIT may pick rows the backup, read before any
	// statement runs, never saw.
	for _, script := range []string{
		"DELETE FROM t ORDER BY a LIMIT 1; DELETE FROM t ORDER BY a LIMIT 1",
		"UPDATE t SET b = 1 WHERE a = 1; UPDATE t SET b = b + 1 ORDER BY b LIMIT 1",
	} {
		_, err = transform(t, script)
		require.EqualError(t, err, "prior backup cannot handle a statement with LIMIT after another statement on the same table db.t", script)
	}
}

func TestTransformDMLToSelectNeedsRowKey(t *testing.T) {
	// nokey has no primary key and only a unique key on a nullable column, so
	// a join's repeated row and two equal rows look alike.
	got, err := transform(t, "DELETE FROM nokey WHERE a = 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "  SELECT `nokey`.`a`,`nokey`.`b` FROM nokey WHERE a = 1;"), got[0].Statement)
	_, err = transform(t, "DELETE nokey FROM nokey JOIN t ON nokey.a = t.a")
	require.EqualError(t, err, "failed to generate SQL for table db.nokey: prior backup cannot back up a table without a primary key or a unique key of non-null columns from a statement that joins other tables")
	_, err = transform(t, "DELETE FROM nokey WHERE a = 1; DELETE FROM nokey WHERE b = 2")
	require.EqualError(t, err, "failed to generate SQL for table db.nokey: prior backup cannot combine several statements on a table without a primary key or a unique key of non-null columns")
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
	require.Equal(t, "CREATE TABLE `b``ak`.`p_a``b_d``b` LIKE `d``b`.`a``b`;\nINSERT INTO `b``ak`.`p_a``b_d``b` (`c``d`)  SELECT `a``b`.`c``d` FROM `a``b` WHERE `c``d` = 1;", got[0].Statement)
}

func TestTransformDMLToSelectBalancesTableRefs(t *testing.T) {
	got, err := transform(t, "UPDATE (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a SET t1.b = 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a;"), got[0].Statement)
	got, err = transform(t, "DELETE t1 FROM (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a WHERE t2.d = ')'")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM (t1 JOIN t2 ON t1.a = t2.c) JOIN t ON t.a = t1.a WHERE t2.d = ')';"), got[0].Statement)
}

func TestTransformDMLToSelectResolvesEachDatabase(t *testing.T) {
	// other's u has a, which db's t also has; the SET column is other.u's.
	catalog := backupCatalog()
	catalog.GetDatabase = func(name string) (analysis.Database, error) {
		switch name {
		case "db", "DB":
			return backupDB, nil
		case "other":
			return otherDB, nil
		}
		return nil, nil
	}
	got, err := analysis.TransformDMLToSelect("UPDATE other.t JOIN other.u ON other.t.id = other.u.id SET a = 1", analysis.BackupOptions{
		DefaultDatabase: "db",
		Catalog:         catalog,
		Columns: func(string, string) (analysis.TableColumns, error) {
			return analysis.TableColumns{Normal: []string{"id"}, RowKey: true}, nil
		},
	})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "u", got[0].SourceTable)

	// Without other's metadata either table may own the column.
	got, err = transform(t, "UPDATE other.t JOIN other.u ON other.t.id = other.u.id SET a = 1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "resource not found: database: other")
}

var otherDB = &backupDatabase{
	name: "other",
	tables: []*metadata.TableMetadata{
		{Name: "t", Columns: backupColumns("id", "b")},
		{Name: "u", Columns: backupColumns("id", "a")},
	},
}

func TestTransformDMLToSelectCTEAlias(t *testing.T) {
	tests := []struct {
		statement string
		want      []string
	}{
		// A table aliased with a CTE's name is still a target.
		{"WITH c AS (SELECT 1 AS a) DELETE c FROM t AS c WHERE c.a = 1", []string{"t"}},
		{"WITH c AS (SELECT 1 AS a) UPDATE t AS c JOIN t1 ON c.a = t1.a SET c.b = 1", []string{"t"}},
		// A table named with its database is never the CTE of that name.
		{"WITH t1 AS (SELECT 1 AS a) DELETE db.t1 FROM db.t1 JOIN t2 ON db.t1.a = t2.c", []string{"t1"}},
		{"WITH t1 AS (SELECT 1 AS a) UPDATE db.t1 JOIN t2 ON db.t1.a = t2.c SET db.t1.b = 1", []string{"t1"}},
		// The CTE itself is not.
		{"WITH c AS (SELECT 1 AS a) DELETE t1, c FROM t1 JOIN c ON t1.a = c.a", []string{"t1"}},
		{"WITH c AS (SELECT 1 AS a) UPDATE t1 JOIN c ON t1.a = c.a SET t1.b = 1, c.a = 2", []string{"t1"}},
	}
	for _, tc := range tests {
		got, err := transform(t, tc.statement)
		require.NoError(t, err, tc.statement)
		var tables []string
		for _, backup := range got {
			tables = append(tables, backup.SourceTable)
		}
		require.Equal(t, tc.want, tables, tc.statement)
	}
}

func TestTransformDMLToSelectRefusesUnsafeScripts(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		// A later statement reads a table an earlier one changes, but every
		// backup reads before the script runs.
		{"UPDATE t1 SET b = 1 WHERE a = 1; DELETE FROM t WHERE a IN (SELECT a FROM t1 WHERE b = 1)", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		{"INSERT INTO t1 VALUES (1, 1); UPDATE t JOIN t1 ON t.a = t1.a SET t.b = 1", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		{"DELETE FROM t1 WHERE a = 1; WITH c AS (SELECT a FROM t1) DELETE t FROM t JOIN c ON t.a = c.a", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		// A foreign key cascades the DELETE to rows the backup leaves out.
		{"DELETE FROM parent WHERE id = 1", "failed to generate SQL for table db.parent: prior backup cannot back up db.parent: foreign key fk_cascade of child is ON DELETE CASCADE, and the rows it changes are not backed up"},
		// A user variable an earlier statement sets picks the rows.
		{"SET @victim = 1; DELETE FROM t WHERE a = @victim", "failed to prepare transformation: prior backup cannot handle a statement that reads @victim after an earlier statement sets it"},
		{"SELECT 1 INTO @v; UPDATE t SET b = 1 WHERE a = @V", "failed to prepare transformation: prior backup cannot handle a statement that reads @V after an earlier statement sets it"},
		{"SELECT @v := 2; UPDATE t JOIN t1 ON t1.a = @v SET t.b = 1", "failed to prepare transformation: prior backup cannot handle a statement that reads @v after an earlier statement sets it"},
		{"SET @v = 1; WITH c AS (SELECT @v AS a) DELETE t FROM t JOIN c ON t.a = c.a", "failed to prepare transformation: prior backup cannot handle a statement that reads @v after an earlier statement sets it"},
		// A clause in an executable comment.
		{"DELETE FROM t /*!80000 WHERE a = 1 */", "prior backup cannot back up t: its statement has clauses in an executable comment"},
		// A system variable an earlier SET changes picks the rows.
		{"SET SESSION auto_increment_increment = 2; DELETE FROM t WHERE @@session.auto_increment_increment = 2", "failed to prepare transformation: prior backup cannot handle a statement that reads @@auto_increment_increment after an earlier statement sets it"},
		{"SET @@SESSION.x = 1; UPDATE t SET b = 1 WHERE a = @@X", "failed to prepare transformation: prior backup cannot handle a statement that reads @@X after an earlier statement sets it"},
		// A target's rows depend on another target an earlier statement changes.
		{"UPDATE t1 SET b = 1; UPDATE t JOIN t1 ON t.a = t1.a SET t.b = 2, t1.b = 2 WHERE t1.b = 1", "failed to prepare transformation: prior backup cannot handle a statement that reads db.t1 after an earlier statement changes it"},
		// A sequence moves when the backup reads it.
		{"DELETE FROM t WHERE a = LASTVAL(s)", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on LASTVAL(), whose value can differ when the backup runs"},
		{"UPDATE t SET b = 1 WHERE a < SETVAL(s, 5)", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on SETVAL(), whose value can differ when the backup runs"},
		// The server's ON UPDATE stamp cascades through a foreign key.
		{"UPDATE stamped_parent SET v = 1 WHERE id = 1", "failed to generate SQL for table db.stamped_parent: prior backup cannot back up db.stamped_parent: foreign key fk_ts of stamped_child is ON UPDATE CASCADE, and the rows it changes are not backed up"},
		{"UPDATE STAMPED_PARENT SET v = 1 WHERE id = 1", "failed to generate SQL for table db.STAMPED_PARENT: prior backup cannot back up db.STAMPED_PARENT: foreign key fk_ts of stamped_child is ON UPDATE CASCADE, and the rows it changes are not backed up"},
		// The backup runs the row selection's := too.
		{"DELETE FROM t WHERE (@n := COALESCE(@n, 0) + 1) = 2", "failed to prepare transformation: prior backup cannot back up a statement whose row selection sets @n, as the backup sets it too"},
		// A USE before the statement moves its unqualified tables.
		{"USE other; DELETE FROM t WHERE a = 1", "failed to prepare transformation: prior backup cannot handle a statement after USE other changes the database"},
		// Two references known by one name.
		{"UPDATE db.t JOIN db2.t ON db.t.a = db2.t.a SET db.t.b = 1", "failed to prepare transformation: prior backup cannot tell apart the tables named t without aliases"},
		{"DELETE db.t FROM db.t JOIN db2.t ON db.t.a = db2.t.a", "failed to prepare transformation: prior backup cannot tell apart the tables named t without aliases"},
		{"DELETE t, t FROM t JOIN t1 ON t.a = t1.a", "failed to prepare transformation: prior backup cannot tell apart the tables named t without aliases"},
		// The backup would draw other random values.
		{"DELETE FROM t WHERE RAND() < 0.5", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on RAND(), whose value can differ when the backup runs"},
		{"DELETE FROM t ORDER BY rand() LIMIT 1", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on RAND(), whose value can differ when the backup runs"},
		{"UPDATE t JOIN t1 ON t1.a = UUID_SHORT() SET t.b = 1", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on UUID_SHORT(), whose value can differ when the backup runs"},
		// The clock moves, and the session's last statement is another one.
		{"DELETE FROM t WHERE a < UNIX_TIMESTAMP() - 60", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on UNIX_TIMESTAMP(), whose value can differ when the backup runs"},
		{"UPDATE t SET b = 1 WHERE a < current_timestamp", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on CURRENT_TIMESTAMP(), whose value can differ when the backup runs"},
		{"INSERT INTO t1 VALUES (1, 1); DELETE FROM t WHERE a = LAST_INSERT_ID()", "failed to prepare transformation: prior backup cannot back up a statement whose rows depend on LAST_INSERT_ID(), whose value can differ when the backup runs"},
		// A trigger on the statement's event, or on the INSERT a restore runs.
		{"UPDATE bumped SET v = 2 WHERE id = 1", "failed to generate SQL for table db.bumped: prior backup cannot back up db.bumped: trigger tr_bump runs on UPDATE, and the backup and its restore do not account for it"},
		{"UPDATE BUMPED SET v = 2 WHERE id = 1", "failed to generate SQL for table db.BUMPED: prior backup cannot back up db.BUMPED: trigger tr_bump runs on UPDATE, and the backup and its restore do not account for it"},
		{"DELETE FROM logged WHERE id = 1", "failed to generate SQL for table db.logged: prior backup cannot back up db.logged: trigger tr_log runs on DELETE, and the backup and its restore do not account for it"},
		{"UPDATE stamped_insert SET v = 2 WHERE id = 1", "failed to generate SQL for table db.stamped_insert: prior backup cannot back up db.stamped_insert: trigger tr_insert runs on INSERT, and the backup and its restore do not account for it"},
		{"DELETE FROM stamped_insert WHERE id = 1", "failed to generate SQL for table db.stamped_insert: prior backup cannot back up db.stamped_insert: trigger tr_insert runs on INSERT, and the backup and its restore do not account for it"},
		// REPLACE and INSERT ... ON DUPLICATE KEY UPDATE overwrite rows.
		{"REPLACE INTO t VALUES (1, 2, 3); DELETE FROM t WHERE b = 2", "failed to prepare transformation: prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements"},
		{"INSERT INTO t1 VALUES (1, 1) ON DUPLICATE KEY UPDATE b = 2", "failed to prepare transformation: prior backup does not support REPLACE or INSERT ... ON DUPLICATE KEY UPDATE statements"},
		// Truncation gives both tables' backups one name.
		{"DELETE FROM `" + strings.Repeat("x", 62) + "1`; DELETE FROM `" + strings.Repeat("x", 62) + "2`", "prior backup cannot name the backups of db." + strings.Repeat("x", 62) + "1 and db." + strings.Repeat("x", 62) + "2 apart: both are pre_" + strings.Repeat("x", 60)},
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
		// Statements on one table read its rows only through the backup.
		"UPDATE t SET b = 1 WHERE a = 1; UPDATE t SET c = 2 WHERE b = 1",
		// A user variable only the SET list reads, or set after its use.
		"SET @v = 1; UPDATE t SET b = @v WHERE a = 1",
		"DELETE FROM t WHERE a = @v; SET @v = 1",
		"SET @@session.v = 1; DELETE FROM t WHERE a = @v",
		"SET @v = 1; DELETE FROM t WHERE a = @@v",
		// A USE of the task database, or after the statement.
		"USE db; DELETE FROM t WHERE a = 1",
		"USE other; USE DB; DELETE FROM t WHERE a = 1",
		"DELETE FROM t WHERE a = 1; USE other",
		// Aliases tell the tables apart, and a CTE is not a table.
		"UPDATE db.t x JOIN db2.t y ON x.a = y.a SET x.b = 1",
		"WITH t AS (SELECT 1 AS a) UPDATE db.t JOIN t ON db.t.a = t.a SET db.t.b = 1",
		// A random value only the SET list draws.
		"UPDATE t SET b = RAND() WHERE a = 1",
		"UPDATE t SET b = (@n := @n + 1) WHERE a = 1",
		"DELETE FROM t WHERE @n = a",
		"SET @@x = 1; DELETE FROM t WHERE a = @x",
		"UPDATE t SET b = 1 WHERE a = 1; UPDATE t x JOIN t y ON x.a = y.b SET x.c = 2, y.c = 3",
		"/*!40000 UPDATE t SET b = 1 */",
		"DELETE FROM t /* WHERE a = 1 */ WHERE a = 2",
		"DELETE FROM t WHERE b = '/*'",
		// MySQL 8.0 does not keep what EXPLAIN ANALYZE changes.
		"EXPLAIN ANALYZE DELETE t FROM t JOIN t1 ON t.a = t1.a",
		"UPDATE t SET b = NOW() WHERE a = 1",
		"DELETE FROM t WHERE a = UNIX_TIMESTAMP(b)",
		// A trigger on events neither the statement nor its restore runs.
		"DELETE FROM bumped WHERE id = 1",
		"UPDATE logged SET v = 1 WHERE id = 1",
		"INSERT IGNORE INTO t1 VALUES (1, 1)",
		// An earlier write to a table the later statement does not read.
		"UPDATE t1 SET b = 1; DELETE FROM t WHERE a = 1",
		// The child's RESTRICT key does not cascade an UPDATE.
		"UPDATE parent SET id = 2 WHERE id = 1",
		// A key on a column the UPDATE leaves alone.
		"UPDATE parent SET v = 1",
	} {
		_, err := transform(t, script)
		require.NoError(t, err, script)
	}
}

func TestTransformDMLToSelectSkipsOuterJoinRows(t *testing.T) {
	got, err := transform(t, "DELETE t FROM t1 LEFT JOIN t ON t.a = t1.a WHERE t1.b = 1 OR t1.b = 2")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM t1 LEFT JOIN t ON t.a = t1.a WHERE (t1.b = 1 OR t1.b = 2) AND `t`.`a` IS NOT NULL;"), got[0].Statement)
	got, err = transform(t, "UPDATE t RIGHT JOIN t1 ON t.a = t1.a SET t.b = 1")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM t RIGHT JOIN t1 ON t.a = t1.a WHERE `t`.`a` IS NOT NULL;"), got[0].Statement)
	// The preserved side of an outer join is backed up as is.
	got, err = transform(t, "DELETE t1 FROM t1 LEFT JOIN t ON t.a = t1.a")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got[0].Statement, "FROM t1 LEFT JOIN t ON t.a = t1.a;"), got[0].Statement)
}
