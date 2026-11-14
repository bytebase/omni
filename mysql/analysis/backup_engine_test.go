package analysis_test

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	"github.com/bytebase/omni/mysql/parser"
)

// TestMain runs the tests through the tape harness, and removes the servers a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// roundTrip is a script whose backup and restore must give the tables' rows
// back, or that prior backup or the restore refuses.
type roundTrip struct {
	name   string
	setup  string
	script string
	// tables are compared before the script and after the restore; the
	// first is the one backed up. "other." names the second database.
	tables       []string
	transformErr string
	restoreErr   string
}

// otherDatabase stands for the second database in a case's text.
const otherDatabase = "omni_t_mysql_analysis_backup_other"

var roundTrips = []roundTrip{
	{
		// Each statement keeps its own ORDER BY and LIMIT in the backup.
		name: "limits",
		setup: "CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b));" +
			"INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40)",
		script: "DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1; DELETE FROM limits WHERE id = 1",
		tables: []string{"limits"},
	},
	{
		// The backup reads both statements' rows before either runs, so the
		// second LIMIT would delete a row it never read.
		name:         "overlapping_limits",
		setup:        "CREATE TABLE overlapping_limits (id INT PRIMARY KEY);INSERT INTO overlapping_limits VALUES (1), (2)",
		script:       "DELETE FROM overlapping_limits ORDER BY id LIMIT 1; DELETE FROM overlapping_limits ORDER BY id LIMIT 1",
		tables:       []string{"overlapping_limits"},
		transformErr: "prior backup cannot handle a statement with LIMIT after another statement",
	},
	{
		// The restore undoes every column any of the UPDATEs changes.
		name: "updates",
		setup: "CREATE TABLE updates (id INT PRIMARY KEY, a INT, b INT, g INT AS (a * 2) STORED);" +
			"INSERT INTO updates (id, a, b) VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)",
		script: "UPDATE updates SET a = 10 WHERE id = 1; UPDATE updates SET b = 20 WHERE id = 2",
		tables: []string{"updates"},
	},
	{
		// Both aliases of a self-join are backed up and restored.
		name: "self_join",
		setup: "CREATE TABLE self_join (id INT PRIMARY KEY, a INT, b INT);" +
			"INSERT INTO self_join VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)",
		script: "UPDATE self_join x JOIN self_join y ON x.id = y.id + 1 SET x.a = 10, y.b = 20",
		tables: []string{"self_join"},
	},
	{
		// A row the join repeats is backed up once, so the keyed backup table
		// takes it.
		name: "one_to_many",
		setup: "CREATE TABLE one_to_many (id INT PRIMARY KEY, a INT);CREATE TABLE children (id INT PRIMARY KEY, parent INT);" +
			"INSERT INTO one_to_many VALUES (1, 1), (2, 2);INSERT INTO children VALUES (1, 1), (2, 1), (3, 2)",
		script: "DELETE p FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1",
		tables: []string{"one_to_many", "children"},
	},
	{
		// Without a row key a repeated row cannot be told from an equal one.
		name: "keyless_join",
		setup: "CREATE TABLE keyless_join (id INT, a INT);CREATE TABLE keyless_children (id INT PRIMARY KEY, parent INT);" +
			"INSERT INTO keyless_join VALUES (1, 1), (1, 1);INSERT INTO keyless_children VALUES (1, 1), (2, 1)",
		script:       "DELETE p FROM keyless_join p JOIN keyless_children c ON c.parent = p.id",
		tables:       []string{"keyless_join"},
		transformErr: "prior backup cannot back up a table without a primary key or a unique key of non-null columns from a statement that joins other tables",
	},
	{
		// A keyless table's equal rows are all backed up and restored.
		name:   "keyless_rows",
		setup:  "CREATE TABLE keyless_rows (id INT, a INT);INSERT INTO keyless_rows VALUES (1, 1), (1, 1), (2, 2)",
		script: "DELETE FROM keyless_rows WHERE id = 1",
		tables: []string{"keyless_rows"},
	},
	{
		// An invisible column is copied and restored too.
		name:   "invisible",
		setup:  "CREATE TABLE invisible (id INT PRIMARY KEY, secret INT INVISIBLE DEFAULT 0);INSERT INTO invisible (id, secret) VALUES (1, 42), (2, 43)",
		script: "DELETE FROM invisible WHERE id = 1",
		tables: []string{"invisible"},
	},
	{
		// A parenthesized join keeps both of its parentheses.
		name: "parens",
		setup: "CREATE TABLE parens (id INT PRIMARY KEY, a INT);CREATE TABLE paren_peers (id INT PRIMARY KEY);" +
			"INSERT INTO parens VALUES (1, 1), (2, 2);INSERT INTO paren_peers VALUES (1), (2)",
		script: "UPDATE (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id SET parens.a = 9 WHERE parens.id = 1",
		tables: []string{"parens"},
	},
	{
		// Names that hold a backtick stay quoted.
		name:   "backtick",
		setup:  "CREATE TABLE `back``tick` (id INT PRIMARY KEY, `c``ol` INT);INSERT INTO `back``tick` VALUES (1, 1), (2, 2)",
		script: "DELETE FROM `back``tick` WHERE `c``ol` = 1",
		tables: []string{"back`tick"},
	},
	{
		// An unqualified SET column belongs to the table that has it in its
		// own database, here other.u, though the default database's t has a
		// column of that name.
		name: "other_database",
		setup: "CREATE TABLE t (id INT PRIMARY KEY, a INT);CREATE TABLE u (id INT PRIMARY KEY, b INT);" +
			"DROP DATABASE IF EXISTS " + otherDatabase + ";CREATE DATABASE " + otherDatabase + ";" +
			"CREATE TABLE " + otherDatabase + ".t (id INT PRIMARY KEY, b INT);CREATE TABLE " + otherDatabase + ".u (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO " + otherDatabase + ".t VALUES (1, 1);INSERT INTO " + otherDatabase + ".u VALUES (1, 1)",
		script: "UPDATE " + otherDatabase + ".t JOIN " + otherDatabase + ".u ON " + otherDatabase + ".t.id = " + otherDatabase + ".u.id SET a = 7",
		tables: []string{"other.u", "other.t"},
	},
	{
		// The restore quotes the statement in a comment that "*/" must not end.
		name: "comments",
		setup: "CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10));" +
			"INSERT INTO comments VALUES (1, '*/'), (2, 'x')",
		script: "DELETE /* note */ FROM comments WHERE note = '*/'",
		tables: []string{"comments"},
	},
	{
		// An arm that opens with a WITH clause keeps it inside its parentheses.
		name: "cte_arms",
		setup: "CREATE TABLE cte_arms (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO cte_arms VALUES (1, 1), (2, 2), (3, 3)",
		script: "WITH c AS (SELECT 1 AS id) DELETE t FROM cte_arms t JOIN c ON t.id = c.id; WITH c AS (SELECT 3 AS id) DELETE t FROM cte_arms t JOIN c ON t.id = c.id",
		tables: []string{"cte_arms"},
	},
	{
		// A table aliased with a CTE's name is still backed up.
		name: "cte_alias",
		setup: "CREATE TABLE cte_alias (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO cte_alias VALUES (1, 1), (2, 2)",
		script: "WITH c AS (SELECT 1 AS id) DELETE c FROM cte_alias AS c WHERE c.id = 1",
		tables: []string{"cte_alias"},
	},
	{
		// A unique key on a generated column changes with the UPDATE, so it
		// cannot match the backed-up row to the changed one.
		name: "generated_key",
		setup: "CREATE TABLE generated_key (id INT, a INT, g INT AS (a + 1), UNIQUE KEY uk_g (g));" +
			"INSERT INTO generated_key (id, a) VALUES (1, 1), (2, 2)",
		script:     "UPDATE generated_key SET a = a + 10 WHERE id = 1",
		tables:     []string{"generated_key"},
		restoreErr: "no disjoint unique key found",
	},
	{
		name: "expression_key",
		setup: "CREATE TABLE expression_key (id INT, a INT, UNIQUE KEY uk_a ((a + 1)));" +
			"INSERT INTO expression_key VALUES (1, 1), (2, 2)",
		script:     "UPDATE expression_key SET a = a + 10 WHERE id = 1",
		tables:     []string{"expression_key"},
		restoreErr: "no disjoint unique key found",
	},
	{
		// NULL keys never collide, so a NULL row's backup would be inserted
		// beside the changed row.
		name:       "nullable_key",
		setup:      "CREATE TABLE nullable_key (u INT NULL, v INT, UNIQUE KEY uk_u (u));INSERT INTO nullable_key VALUES (NULL, 1), (2, 2)",
		script:     "UPDATE nullable_key SET v = 3 WHERE u IS NULL",
		tables:     []string{"nullable_key"},
		restoreErr: "no disjoint unique key found",
	},
}

// TestBackupRestoreRoundTrip runs each backup, its statements and their
// restore on MySQL, and requires the tables' rows back.
func TestBackupRestoreRoundTrip(t *testing.T) {
	db := testcontainer.OpenDatabase(t, "mysql_analysis_backup")
	database := testcontainer.DatabaseName("mysql_analysis_backup")
	name := func(table string) (string, string) {
		if rest, ok := strings.CutPrefix(table, "other."); ok {
			return otherDatabase, rest
		}
		return database, table
	}
	rowsOf := func(table string) []string {
		database, table := name(table)
		return tableRows(t, db, database, table)
	}
	for _, rt := range roundTrips {
		runScript(t, db, rt.setup)
		var before [][]string
		for _, table := range rt.tables {
			before = append(before, rowsOf(table))
		}
		catalog := engineCatalog(t, db, database, otherDatabase)

		backups, err := analysis.TransformDMLToSelect(rt.script, analysis.BackupOptions{
			DefaultDatabase: database,
			BackupDatabase:  database,
			TablePrefix:     "bak",
			Catalog:         catalog,
		})
		if rt.transformErr != "" {
			require.ErrorContains(t, err, rt.transformErr, rt.name)
			continue
		}
		require.NoError(t, err, rt.name)
		require.Len(t, backups, 1, rt.name)
		runScript(t, db, sortUnionArms(backups[0].Statement))
		runScript(t, db, rt.script)
		sourceDatabase, sourceTable := name(rt.tables[0])
		require.Equal(t, sourceTable, backups[0].SourceTable, rt.name)
		require.NotEqual(t, before[0], tableRows(t, db, sourceDatabase, sourceTable), rt.name)

		restored, err := restoreSQL(t, rt.script, sourceDatabase, database, backups[0], catalog)
		if rt.restoreErr != "" {
			require.ErrorContains(t, err, rt.restoreErr, rt.name)
			continue
		}
		require.NoError(t, err, rt.name)
		runScript(t, db, restored)
		for i, table := range rt.tables {
			require.Equal(t, before[i], rowsOf(table), rt.name)
		}
	}
}

// restoreSQL builds the restore of backup as Bytebase does: from the
// statements in its range that name the table, the first DML decides.
func restoreSQL(t *testing.T, script, database, backupDatabase string, backup analysis.BackupStatement, catalog analysis.Catalog) (string, error) {
	t.Helper()
	var statements []string
	for _, segment := range parser.Split(script) {
		if segment.ByteStart < backup.ByteStart || segment.ByteEnd > backup.ByteEnd {
			continue
		}
		ok, err := analysis.ReferencesTable(segment.Text, database, backup.SourceTable)
		require.NoError(t, err)
		if ok {
			statements = append(statements, segment.Text+";")
		}
	}
	original := strings.Join(statements, "\n")
	nodes, err := analysis.DMLs(original)
	require.NoError(t, err)
	require.NotEmpty(t, nodes, original)
	columns, err := analysis.ClassifyColumns(catalog, false, database, backup.SourceTable)
	require.NoError(t, err)
	return analysis.RestoreSQL(nodes, analysis.RestoreOptions{
		Statement:      original,
		Database:       database,
		Table:          backup.SourceTable,
		BackupDatabase: backupDatabase,
		BackupTable:    backup.BackupTable,
		Columns:        columns.Normal,
		Catalog:        catalog,
	})
}

// sortUnionArms orders a backup's UNION arms, whose order follows map
// iteration, so the tape sees the same request each run.
func sortUnionArms(statement string) string {
	const union = "\n  UNION DISTINCT\n"
	head, body, ok := strings.Cut(statement, "  (")
	if !ok {
		return statement
	}
	arms := strings.Split(strings.TrimSuffix("  ("+body, ";"), union)
	sort.Strings(arms)
	return head + strings.Join(arms, union) + ";"
}

func runScript(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, segment := range parser.Split(script) {
		if segment.Empty() {
			continue
		}
		_, err := db.Exec(segment.Text)
		require.NoError(t, err, segment.Text)
	}
}

func quote(name string) string { return "`" + strings.ReplaceAll(name, "`", "``") + "`" }

// tableRows returns every row of database.table, ordered, as text. It names
// every column, as SELECT * leaves invisible ones out.
func tableRows(t *testing.T, db *sql.DB, database, table string) []string {
	t.Helper()
	var columns, positions []string
	rows, err := db.Query("SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", database, table)
	require.NoError(t, err)
	for rows.Next() {
		var column string
		require.NoError(t, rows.Scan(&column))
		columns = append(columns, quote(column))
		positions = append(positions, fmt.Sprint(len(columns)))
	}
	require.NoError(t, rows.Close())
	rows, err = db.Query(fmt.Sprintf("SELECT %s FROM %s.%s ORDER BY %s", strings.Join(columns, ", "), quote(database), quote(table), strings.Join(positions, ", ")))
	require.NoError(t, err)
	defer rows.Close()
	var result []string
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		var row []string
		for _, value := range values {
			row = append(row, fmt.Sprintf("%v", value))
		}
		result = append(result, strings.Join(row, "|"))
	}
	require.NoError(t, rows.Err())
	return result
}

// engineCatalog reads the databases' columns and keys from
// information_schema.
func engineCatalog(t *testing.T, db *sql.DB, databases ...string) analysis.Catalog {
	t.Helper()
	loaded := map[string]*backupDatabase{}
	for _, database := range databases {
		tables := map[string]*metadata.TableMetadata{}
		var order []*metadata.TableMetadata
		table := func(name string) *metadata.TableMetadata {
			if tables[name] == nil {
				tables[name] = &metadata.TableMetadata{Name: name}
				order = append(order, tables[name])
			}
			return tables[name]
		}
		rows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION", database)
		require.NoError(t, err)
		for rows.Next() {
			var tableName, name, generation, nullable string
			require.NoError(t, rows.Scan(&tableName, &name, &generation, &nullable))
			column := &metadata.ColumnMetadata{Name: name, Nullable: nullable == "YES"}
			if generation != "" {
				column.Generation = &metadata.GenerationMetadata{Expression: generation}
			}
			table(tableName).Columns = append(table(tableName).Columns, column)
		}
		require.NoError(t, rows.Close())
		rows, err = db.Query("SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX", database)
		require.NoError(t, err)
		indexes := map[string]*metadata.IndexMetadata{}
		for rows.Next() {
			var tableName, name, part string
			var nonUnique bool
			require.NoError(t, rows.Scan(&tableName, &name, &nonUnique, &part))
			index := indexes[tableName+"."+name]
			if index == nil {
				index = &metadata.IndexMetadata{Name: name, Primary: name == "PRIMARY", Unique: !nonUnique}
				indexes[tableName+"."+name] = index
				table(tableName).Indexes = append(table(tableName).Indexes, index)
			}
			index.Expressions = append(index.Expressions, part)
		}
		require.NoError(t, rows.Close())
		loaded[database] = &backupDatabase{name: database, tables: order}
	}
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) { return databases, nil },
		GetDatabase: func(name string) (analysis.Database, error) {
			if database, ok := loaded[name]; ok {
				return database, nil
			}
			return nil, nil
		},
	}
}
