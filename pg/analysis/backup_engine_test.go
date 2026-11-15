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
	// transformErr, when set, is the error prior backup refuses with.
	transformErr string
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
		// The restore quotes the statement in a block comment, which
		// PostgreSQL nests: neither "*/" nor "/*" may end or extend it.
		name:   "comment_delimiters",
		setup:  "CREATE TABLE notes (id INT PRIMARY KEY, note TEXT); INSERT INTO notes VALUES (1, '*/'), (2, '/*'), (3, 'x')",
		script: "DELETE FROM notes WHERE note IN ('*/ DROP TABLE canary; /*', '/*/') OR id < 3",
		schema: "public",
		table:  "notes",
	},
	{
		// The UPDATE changes the primary key, so the restore matches rows by
		// a unique index that is no constraint.
		name: "unique_index",
		setup: "CREATE TABLE ui (id INT PRIMARY KEY, code INT, v INT); CREATE UNIQUE INDEX ui_code ON ui (code);" +
			"CREATE UNIQUE INDEX ui_lower ON ui ((code + 1)); INSERT INTO ui VALUES (1, 10, 0), (2, 20, 0)",
		script: "UPDATE ui SET id = id + 100, v = 1 WHERE id = 1",
		schema: "public",
		table:  "ui",
	},
	{
		// The primary key's column needs quoting; the UPDATE changes it.
		name:   "quoted_key_column",
		setup:  `CREATE TABLE mk ("Id" INT PRIMARY KEY, code INT UNIQUE, v INT); INSERT INTO mk VALUES (1, 10, 0), (2, 20, 0)`,
		script: `UPDATE mk SET "Id" = "Id" + 100, v = 1 WHERE "Id" = 1`,
		schema: "public",
		table:  "mk",
	},
	{
		// Each statement names the table by its own alias.
		name:   "alias_per_statement",
		setup:  "CREATE TABLE alias_t (id INT PRIMARY KEY, v INT); INSERT INTO alias_t VALUES (1, 0), (2, 0), (3, 0)",
		script: "UPDATE alias_t AS a SET v = 1 WHERE a.id = 1; UPDATE alias_t AS b SET v = 2 WHERE b.id = 2; UPDATE alias_t SET v = 3 WHERE id = 3",
		schema: "public",
		table:  "alias_t",
	},
	{
		// SET fn.first assigns a field of the composite column fn.
		name: "composite_field",
		setup: "CREATE TYPE fullname AS (first text, last text);" +
			"CREATE TABLE people (id INT PRIMARY KEY, fn fullname, tags INT[]);" +
			"INSERT INTO people VALUES (1, ('a', 'b'), '{1,2}'), (2, ('c', 'd'), '{3}')",
		script: "UPDATE people SET fn.first = 'x', fn.last = 'y', tags[1] = 9 WHERE id = 1",
		schema: "public",
		table:  "people",
	},
	{
		// The backup's name is longer than PostgreSQL's 63 bytes, and its
		// characters take three bytes each.
		name:   "long_multibyte_name",
		setup:  `CREATE TABLE "` + strings.Repeat("表", 20) + `" (id INT PRIMARY KEY); INSERT INTO "` + strings.Repeat("表", 20) + `" VALUES (1), (2)`,
		script: `DELETE FROM "` + strings.Repeat("表", 20) + `" WHERE id = 1`,
		schema: "public",
		table:  strings.Repeat("表", 20),
	},
	{
		// Both backups would be named _bb__a_b_c.
		name:         "colliding_names",
		setup:        "CREATE SCHEMA c; CREATE TABLE c.a_b (id INT PRIMARY KEY); CREATE SCHEMA b_c; CREATE TABLE b_c.a (id INT PRIMARY KEY)",
		script:       "DELETE FROM c.a_b; DELETE FROM b_c.a",
		transformErr: `prior backup cannot back up both "b_c"."a" and "c"."a_b": their backup tables would both be named "_bb__a_b_c"`,
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
	// No restore may drop it.
	exec(t, db, "CREATE TABLE canary (id INT)")
	for _, rt := range roundTrips {
		exec(t, db, rt.setup)

		stmts, err := omnipg.Parse(rt.script)
		require.NoError(t, err, rt.name)
		backups, err := analysis.TransformDMLToSelect(stmts, rt.script, analysis.BackupOptions{
			SearchPath:     []string{"public"},
			SearchRelation: func(searchPath []string, name string) string { return searchRelation(t, db, searchPath, name) },
			TargetSchema:   backupSchema,
			TablePrefix:    "_bb_",
		})
		if rt.transformErr != "" {
			require.EqualError(t, err, rt.transformErr, rt.name)
			continue
		}
		require.NoError(t, err, rt.name)
		before := tableRows(t, db, rt.schema, rt.table)
		require.Len(t, backups, 1, rt.name)
		require.Equal(t, rt.schema, backups[0].SourceSchema, rt.name)
		require.Equal(t, rt.table, backups[0].SourceTableName, rt.name)
		exec(t, db, backups[0].Statement)
		require.True(t, tableExists(t, db, backupSchema, backups[0].TargetTableName), rt.name)
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
		require.Equal(t, "canary", relation(t, db, "canary"), rt.name)
	}
}

func exec(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	_, err := db.Exec(script)
	require.NoError(t, err, script)
}

// tableExists reports whether schema holds a table named exactly name.
func tableExists(t *testing.T, db *sql.DB, schema, name string) bool {
	t.Helper()
	var found bool
	require.NoError(t, db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname::text = $1 AND tablename::text = $2)", schema, name).Scan(&found))
	return found
}

// relation returns name if a relation of that name is on the search path.
func relation(t *testing.T, db *sql.DB, name string) string {
	t.Helper()
	var found sql.NullString
	require.NoError(t, db.QueryRow("SELECT to_regclass($1)::text", name).Scan(&found))
	return found.String
}

func quote(name string) string { return `"` + strings.ReplaceAll(name, `"`, `""`) + `"` }

// searchRelation returns the first schema of searchPath holding a table
// named name, as Bytebase's metadata lookup does.
func searchRelation(t *testing.T, db *sql.DB, searchPath []string, name string) string {
	t.Helper()
	for _, schema := range searchPath {
		var found bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname::text = $1 AND tablename::text = $2)", schema, name).Scan(&found)
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

// tableMetadata reads schema.table's columns and indexes, and its primary
// key, as Bytebase's sync records them: each index part as pg_get_indexdef
// writes it.
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
	rows, err = db.Query(`SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL, pg_get_indexdef(i.oid),
		(SELECT string_agg(pg_get_indexdef(i.oid, k + 1, true), E'\x1f' ORDER BY k) FROM generate_subscripts(x.indkey, 1) AS k WHERE k < x.indnkeyatts)
		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname`, schema, table)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var name, definition, columns string
		var primary, unique, constraint bool
		require.NoError(t, rows.Scan(&name, &primary, &unique, &constraint, &definition, &columns))
		index := &metadata.IndexMetadata{Name: name, Primary: primary, Unique: unique, IsConstraint: constraint, Definition: definition, Expressions: strings.Split(columns, "\x1f")}
		result.Indexes = append(result.Indexes, index)
		if primary {
			primaryKey = index
		}
	}
	require.NoError(t, rows.Err())
	return result, primaryKey
}
