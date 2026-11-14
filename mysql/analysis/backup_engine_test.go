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

// roundTrip is a table, the statements that change it, and whether its
// restore is refused.
type roundTrip struct {
	table      string
	setup      string
	script     string
	restoreErr string
}

var roundTrips = []roundTrip{
	{
		// Each statement keeps its own ORDER BY and LIMIT in the backup.
		table: "limits",
		setup: "CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b));" +
			"INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40)",
		script: "DELETE FROM limits WHERE id = 1; DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1",
	},
	{
		table: "update_limits",
		setup: "CREATE TABLE update_limits (id INT PRIMARY KEY, a INT, g INT AS (a * 2) STORED);" +
			"INSERT INTO update_limits (id, a) VALUES (1, 1), (2, 2), (3, 3)",
		script: "UPDATE update_limits SET a = a + 100 WHERE id = 1; UPDATE update_limits SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1",
	},
	{
		// Both aliases of a self-join are backed up and restored.
		table: "self_join",
		setup: "CREATE TABLE self_join (id INT PRIMARY KEY, a INT, b INT);" +
			"INSERT INTO self_join VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)",
		script: "UPDATE self_join x JOIN self_join y ON x.id = y.id + 1 SET x.a = 10, y.b = 20",
	},
	{
		// The restore quotes the statement in a comment that "*/" must not end.
		table: "comments",
		setup: "CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10));" +
			"INSERT INTO comments VALUES (1, '*/'), (2, 'x')",
		script: "DELETE /* note */ FROM comments WHERE note = '*/'",
	},
	{
		// An arm that opens with a WITH clause keeps it inside its parentheses.
		table: "cte_arms",
		setup: "CREATE TABLE cte_arms (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO cte_arms VALUES (1, 1), (2, 2), (3, 3)",
		script: "WITH c AS (SELECT 1 AS id) DELETE t FROM cte_arms t JOIN c ON t.id = c.id; WITH c AS (SELECT 3 AS id) DELETE t FROM cte_arms t JOIN c ON t.id = c.id",
	},
	{
		// A table aliased with a CTE's name is still backed up.
		table: "cte_alias",
		setup: "CREATE TABLE cte_alias (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO cte_alias VALUES (1, 1), (2, 2)",
		script: "WITH c AS (SELECT 1 AS id) DELETE c FROM cte_alias AS c WHERE c.id = 1",
	},
	{
		// A unique key on a generated column changes with the UPDATE, so it
		// cannot match the backed-up row to the changed one.
		table: "generated_key",
		setup: "CREATE TABLE generated_key (id INT, a INT, g INT AS (a + 1), UNIQUE KEY uk_g (g));" +
			"INSERT INTO generated_key (id, a) VALUES (1, 1), (2, 2)",
		script:     "UPDATE generated_key SET a = a + 10 WHERE id = 1",
		restoreErr: "no disjoint unique key found",
	},
	{
		table: "expression_key",
		setup: "CREATE TABLE expression_key (id INT, a INT, UNIQUE KEY uk_a ((a + 1)));" +
			"INSERT INTO expression_key VALUES (1, 1), (2, 2)",
		script:     "UPDATE expression_key SET a = a + 10 WHERE id = 1",
		restoreErr: "no disjoint unique key found",
	},
}

// TestBackupRestoreRoundTrip runs each backup, its statements and their
// restore on MySQL, and requires the table's rows back.
func TestBackupRestoreRoundTrip(t *testing.T) {
	db := testcontainer.OpenDatabase(t, "mysql_analysis_backup")
	database := testcontainer.DatabaseName("mysql_analysis_backup")
	for _, rt := range roundTrips {
		runScript(t, db, rt.setup)
		before := tableRows(t, db, rt.table)
		catalog := engineCatalog(t, db, database, rt.table)

		backups, err := analysis.TransformDMLToSelect(rt.script, analysis.BackupOptions{
			DefaultDatabase: database,
			BackupDatabase:  database,
			TablePrefix:     "bak",
			Catalog:         catalog,
		})
		require.NoError(t, err, rt.script)
		require.Len(t, backups, 1, rt.script)
		runScript(t, db, sortUnionArms(backups[0].Statement))
		runScript(t, db, rt.script)
		require.NotEqual(t, before, tableRows(t, db, rt.table), rt.script)

		restored, err := restoreSQL(t, rt.script, database, backups[0], catalog)
		if rt.restoreErr != "" {
			require.ErrorContains(t, err, rt.restoreErr, rt.script)
			continue
		}
		require.NoError(t, err, rt.script)
		runScript(t, db, restored)
		require.Equal(t, before, tableRows(t, db, rt.table), rt.script)
	}
}

// restoreSQL builds the restore of backup as Bytebase does: from the
// statements in its range that name the table, the first DML decides.
func restoreSQL(t *testing.T, script, database string, backup analysis.BackupStatement, catalog analysis.Catalog) (string, error) {
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
	node, err := analysis.FirstDML(original)
	require.NoError(t, err)
	require.NotNil(t, node, original)
	_, columns, err := analysis.ClassifyColumns(catalog, false, database, backup.SourceTable)
	require.NoError(t, err)
	return analysis.RestoreSQL(node, analysis.RestoreOptions{
		Statement:      original,
		Database:       database,
		Table:          backup.SourceTable,
		BackupDatabase: database,
		BackupTable:    backup.BackupTable,
		Columns:        columns,
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

// tableRows returns every row of table, ordered, as text.
func tableRows(t *testing.T, db *sql.DB, table string) []string {
	t.Helper()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM `%s` ORDER BY 1, 2", table))
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
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

// engineCatalog reads table's columns and keys from information_schema.
func engineCatalog(t *testing.T, db *sql.DB, database, table string) analysis.Catalog {
	t.Helper()
	tableMetadata := &metadata.TableMetadata{Name: table}
	rows, err := db.Query("SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", database, table)
	require.NoError(t, err)
	for rows.Next() {
		var name, generation string
		require.NoError(t, rows.Scan(&name, &generation))
		column := &metadata.ColumnMetadata{Name: name}
		if generation != "" {
			column.Generation = &metadata.GenerationMetadata{Expression: generation}
		}
		tableMetadata.Columns = append(tableMetadata.Columns, column)
	}
	require.NoError(t, rows.Close())
	rows, err = db.Query("SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX", database, table)
	require.NoError(t, err)
	indexes := map[string]*metadata.IndexMetadata{}
	for rows.Next() {
		var name, part string
		var nonUnique bool
		require.NoError(t, rows.Scan(&name, &nonUnique, &part))
		index := indexes[name]
		if index == nil {
			index = &metadata.IndexMetadata{Name: name, Primary: name == "PRIMARY", Unique: !nonUnique}
			indexes[name] = index
			tableMetadata.Indexes = append(tableMetadata.Indexes, index)
		}
		index.Expressions = append(index.Expressions, part)
	}
	require.NoError(t, rows.Close())
	db2 := &backupDatabase{name: database, tables: []*metadata.TableMetadata{tableMetadata}}
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) { return []string{database}, nil },
		GetDatabase: func(name string) (analysis.Database, error) {
			if name == database {
				return db2, nil
			}
			return nil, nil
		},
	}
}
