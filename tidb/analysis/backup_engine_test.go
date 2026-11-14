package analysis_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/analysis"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
	"github.com/bytebase/omni/tidb/parser"
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
		// More than five statements on one table share a backup, and each
		// keeps its own ORDER BY and LIMIT in it.
		table: "limits",
		setup: "CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b));" +
			"INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40), (5, 5, 50), (6, 6, 60), (7, 7, 70), (8, 8, 80), (9, 9, 90), (10, 10, 100), (11, 11, 110), (12, 12, 120)",
		script: "DELETE FROM limits WHERE id IN (1, 2) ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id IN (3, 4) ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id IN (5, 6) ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id IN (7, 8) ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id IN (9, 10) ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id IN (11, 12) ORDER BY id DESC LIMIT 1;\n",
	},
	{
		table: "update_rows",
		setup: "CREATE TABLE update_rows (id INT PRIMARY KEY, a INT, g INT AS (a * 2) STORED);" +
			"INSERT INTO update_rows (id, a) VALUES (1, 1), (2, 2), (3, 3)",
		script: "UPDATE update_rows SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1",
	},
	{
		// The restore quotes the statement in a comment that "*/" must not end.
		table: "comments",
		setup: "CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10));" +
			"INSERT INTO comments VALUES (1, '*/'), (2, 'x')",
		script: "DELETE /* note */ FROM comments WHERE note = '*/'",
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
// restore on TiDB, and requires the table's rows back.
func TestBackupRestoreRoundTrip(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{"DROP DATABASE IF EXISTS test", "CREATE DATABASE test"} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t)
	const database = "test"
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

// restoreSQL builds the restore of backup as Bytebase does: every DML in
// its range that changes the table decides.
func restoreSQL(t *testing.T, script, database string, backup analysis.BackupStatement, catalog analysis.Catalog) (string, error) {
	t.Helper()
	columnSet, err := analysis.NormalColumnSet(catalog, database, backup.SourceTable)
	require.NoError(t, err)
	var statements []string
	for _, segment := range parser.Split(script) {
		if segment.ByteStart < backup.ByteStart || segment.ByteEnd > backup.ByteEnd {
			continue
		}
		ok, err := analysis.MutatesTable(segment.Text, database, backup.SourceTable, columnSet)
		require.NoError(t, err)
		if ok {
			statements = append(statements, segment.Text+";")
		}
	}
	original := strings.Join(statements, "\n")
	nodes, err := analysis.MatchingDMLs(original, database, backup.SourceTable, columnSet)
	require.NoError(t, err)
	require.NotEmpty(t, nodes, original)
	_, columns, err := analysis.ClassifyColumns(catalog, false, database, backup.SourceTable)
	require.NoError(t, err)
	return analysis.RestoreSQL(nodes, analysis.RestoreOptions{
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
	const union = "\n  UNION ALL\n"
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
