package analysis_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestMain runs the tests through the tape harness, and removes the servers a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

const backupSchema = "bbdataarchive"

// roundTrip is a script whose backup and restore must give its table's rows
// back. The table is the one in schema the backup covers.
type roundTrip struct {
	name   string
	setup  string
	script string
	schema string
	table  string
}

var roundTrips = []roundTrip{
	{
		name:   "delete",
		setup:  "CREATE TABLE plain (id INT PRIMARY KEY, a INT); INSERT INTO plain VALUES (1, 1), (2, 2)",
		script: "DELETE FROM plain WHERE id = 1",
		schema: "public",
		table:  "plain",
	},
	{
		name:   "update",
		setup:  "CREATE TABLE plain_update (id INT PRIMARY KEY, a INT); INSERT INTO plain_update VALUES (1, 1), (2, 2)",
		script: "UPDATE plain_update SET a = 10 WHERE id = 2",
		schema: "public",
		table:  "plain_update",
	},
	{
		// Each statement's WITH clause decides its own rows.
		name:  "cte_per_statement",
		setup: "CREATE TABLE cte_t (id INT PRIMARY KEY); INSERT INTO cte_t VALUES (1), (2), (3)",
		script: "WITH x AS (SELECT 1 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id;" +
			"WITH x AS (SELECT 2 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id",
		schema: "public",
		table:  "cte_t",
	},
	{
		// Every name the backup and restore write holds a '"'.
		name: "quoted_delete",
		setup: `CREATE SCHEMA "s""q"; CREATE TABLE "s""q"."t""q" ("i""d" INT CONSTRAINT "p""k" PRIMARY KEY, "a""b" INT);` +
			`INSERT INTO "s""q"."t""q" VALUES (1, 1), (2, 2)`,
		script: `DELETE FROM "s""q"."t""q" AS "x""y" WHERE "x""y"."i""d" = 1`,
		schema: `s"q`,
		table:  `t"q`,
	},
	{
		name: "quoted_update",
		setup: `CREATE SCHEMA "u""q"; CREATE TABLE "u""q"."t""q" ("i""d" INT CONSTRAINT "p""k2" PRIMARY KEY, "a""b" INT);` +
			`INSERT INTO "u""q"."t""q" VALUES (1, 1), (2, 2)`,
		script: `UPDATE "u""q"."t""q" SET "a""b" = 10 WHERE "i""d" = 2`,
		schema: `u"q`,
		table:  `t"q`,
	},
}

// TestBackupRestoreRoundTrip runs each backup, its script and the restore
// on PostgreSQL, and requires the table's rows back.
func TestBackupRestoreRoundTrip(t *testing.T) {
	db := testcontainer.Open(t)
	exec(t, db, "CREATE SCHEMA "+backupSchema)
	for _, rt := range roundTrips {
		exec(t, db, rt.setup)
		before := tableRows(t, db, rt.schema, rt.table)

		stmts, err := omnipg.Parse(rt.script)
		require.NoError(t, err, rt.name)
		backups, err := analysis.TransformDMLToSelect(stmts, rt.script, analysis.BackupOptions{
			SearchPath:     []string{"public"},
			SearchRelation: func(searchPath []string, name string) string { return searchRelation(t, db, searchPath, name) },
			TargetSchema:   backupSchema,
			TablePrefix:    "_bb_",
		})
		require.NoError(t, err, rt.name)
		require.Len(t, backups, 1, rt.name)
		require.Equal(t, rt.schema, backups[0].SourceSchema, rt.name)
		require.Equal(t, rt.table, backups[0].SourceTableName, rt.name)
		exec(t, db, backups[0].Statement)
		exec(t, db, rt.script)
		require.NotEqual(t, before, tableRows(t, db, rt.schema, rt.table), rt.name)

		// Bytebase restores a backup by its first statement's type.
		dmls := analysis.BackedUpDMLs(stmts)
		require.NotEmpty(t, dmls, rt.name)
		table, primaryKey := tableMetadata(t, db, rt.schema, rt.table)
		restore, err := analysis.RestoreSQL(dmls[0].Node, analysis.RestoreOptions{
			OriginalSchema: rt.schema,
			OriginalTable:  rt.table,
			BackupSchema:   backupSchema,
			BackupTable:    backups[0].TargetTableName,
			Table:          table,
			PrimaryKey:     primaryKey,
			OriginalSQL:    rt.script,
		})
		require.NoError(t, err, rt.name)
		exec(t, db, restore)
		require.Equal(t, before, tableRows(t, db, rt.schema, rt.table), rt.name)
	}
}

func exec(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	_, err := db.Exec(script)
	require.NoError(t, err, script)
}

func quote(name string) string { return `"` + strings.ReplaceAll(name, `"`, `""`) + `"` }

// searchRelation returns the first schema of searchPath holding a table
// named name, as Bytebase's metadata lookup does.
func searchRelation(t *testing.T, db *sql.DB, searchPath []string, name string) string {
	t.Helper()
	for _, schema := range searchPath {
		var found bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)", schema, name).Scan(&found)
		require.NoError(t, err)
		if found {
			return schema
		}
	}
	return ""
}

// tableRows returns every row of schema.table, ordered, as text.
func tableRows(t *testing.T, db *sql.DB, schema, table string) []string {
	t.Helper()
	rows, err := db.Query(fmt.Sprintf("SELECT x::text FROM %s.%s AS x ORDER BY 1", quote(schema), quote(table)))
	require.NoError(t, err)
	defer rows.Close()
	var result []string
	for rows.Next() {
		var row string
		require.NoError(t, rows.Scan(&row))
		result = append(result, row)
	}
	require.NoError(t, rows.Err())
	return result
}

// tableMetadata reads schema.table's columns and unique indexes, and its
// primary key, as Bytebase's sync records them.
func tableMetadata(t *testing.T, db *sql.DB, schema, table string) (*metadata.TableMetadata, *metadata.IndexMetadata) {
	t.Helper()
	result := &metadata.TableMetadata{Name: table}
	rows, err := db.Query(`SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum`, schema, table)
	require.NoError(t, err)
	for rows.Next() {
		var name string
		var generated bool
		require.NoError(t, rows.Scan(&name, &generated))
		column := &metadata.ColumnMetadata{Name: name}
		if generated {
			column.Generation = &metadata.GenerationMetadata{}
		}
		result.Columns = append(result.Columns, column)
	}
	require.NoError(t, rows.Close())

	var primaryKey *metadata.IndexMetadata
	rows, err = db.Query(`SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname`, schema, table)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var name, columns string
		var primary, unique, constraint bool
		require.NoError(t, rows.Scan(&name, &primary, &unique, &constraint, &columns))
		index := &metadata.IndexMetadata{Name: name, Primary: primary, Unique: unique, IsConstraint: constraint, Expressions: strings.Split(columns, "\x1f")}
		result.Indexes = append(result.Indexes, index)
		if primary {
			primaryKey = index
		}
	}
	require.NoError(t, rows.Err())
	return result, primaryKey
}
