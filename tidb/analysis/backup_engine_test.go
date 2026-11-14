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

var roundTrips = []roundTrip{
	{
		// More than five statements share a backup, and the first keeps its
		// own ORDER BY and LIMIT in it.
		name: "limits",
		setup: "CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b));" +
			"INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40), (5, 5, 50), (6, 6, 60), (7, 7, 70), (8, 8, 80)",
		script: "DELETE FROM limits WHERE id > 5 ORDER BY id DESC LIMIT 1;\nDELETE FROM limits WHERE id = 1;\nDELETE FROM limits WHERE id = 2;\nDELETE FROM limits WHERE id = 3;\nDELETE FROM limits WHERE id = 4;\nDELETE FROM limits WHERE id = 5;\n",
		tables: []string{"limits"},
	},
	{
		// The backups read both statements' rows before either runs, so the
		// second LIMIT would delete a row neither read.
		name:         "overlapping_limits",
		setup:        "CREATE TABLE overlapping_limits (id INT PRIMARY KEY);INSERT INTO overlapping_limits VALUES (1), (2)",
		script:       "DELETE FROM overlapping_limits ORDER BY id LIMIT 1; DELETE FROM overlapping_limits ORDER BY id LIMIT 1",
		tables:       []string{"overlapping_limits"},
		transformErr: "prior backup cannot handle a statement with LIMIT after another statement",
	},
	{
		// Statements that read the same row back it up once, so the keyed
		// backup table takes it.
		name:   "overlapping_rows",
		setup:  "CREATE TABLE overlapping_rows (id INT PRIMARY KEY, a INT);INSERT INTO overlapping_rows VALUES (1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 6), (7, 7), (8, 8)",
		script: "DELETE FROM overlapping_rows WHERE id IN (1, 2);\nDELETE FROM overlapping_rows WHERE id IN (2, 3);\nDELETE FROM overlapping_rows WHERE id IN (3, 4);\nDELETE FROM overlapping_rows WHERE id IN (4, 5);\nDELETE FROM overlapping_rows WHERE id IN (5, 6);\nDELETE FROM overlapping_rows WHERE id IN (6, 7);\n",
		tables: []string{"overlapping_rows"},
	},
	{
		// The restore undoes every column any of the UPDATEs changes.
		name: "updates",
		setup: "CREATE TABLE updates (id INT PRIMARY KEY, a INT, b INT, g INT AS (a * 2) STORED);" +
			"INSERT INTO updates (id, a, b) VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3), (4, 4, 4), (5, 5, 5), (6, 6, 6)",
		script: "UPDATE updates SET a = 10 WHERE id = 1;\nUPDATE updates SET b = 20 WHERE id = 2;\nUPDATE updates SET a = a + 100 WHERE id = 3;\nUPDATE updates SET a = a + 100 WHERE id = 4;\nUPDATE updates SET a = a + 100 WHERE id = 5;\nUPDATE updates SET a = a + 100 WHERE id = 6;\n",
		tables: []string{"updates"},
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
		// One statement setting a unique key to a literal leaves its old value
		// free, so the restore matches the row by its primary key.
		name:   "single_key",
		setup:  "CREATE TABLE single_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u));INSERT INTO single_key VALUES (1, 10), (2, 20)",
		script: "UPDATE single_key SET u = 30 WHERE id = 1",
		tables: []string{"single_key"},
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
	{
		// The server stamps ts on UPDATE; the restore writes the old stamp back.
		name: "stamped",
		setup: "CREATE TABLE stamped (id INT PRIMARY KEY, v INT, ts TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00' ON UPDATE CURRENT_TIMESTAMP);" +
			"INSERT INTO stamped VALUES (1, 1, '2000-01-01 00:00:00'), (2, 2, '2000-01-01 00:00:00')",
		script: "UPDATE stamped SET v = 5 WHERE id = 1",
		tables: []string{"stamped"},
	},
	{
		// The DELETE reads flags after the UPDATE changes it; its backup
		// would read flags before.
		name: "read_after_write",
		setup: "CREATE TABLE flags (id INT PRIMARY KEY, enabled INT);CREATE TABLE flagged (id INT PRIMARY KEY);" +
			"INSERT INTO flags VALUES (1, 0);INSERT INTO flagged VALUES (1)",
		script:       "UPDATE flags SET enabled = 1 WHERE id = 1; DELETE FROM flagged WHERE id IN (SELECT id FROM flags WHERE enabled = 1)",
		tables:       []string{"flagged"},
		transformErr: "prior backup cannot handle a statement that reads",
	},
	{
		// Deleting a parent deletes its children, which no backup holds.
		name: "cascade",
		setup: "CREATE TABLE fk_parent (id INT PRIMARY KEY);CREATE TABLE fk_child (id INT PRIMARY KEY, parent_id INT, FOREIGN KEY fk_child_parent (parent_id) REFERENCES fk_parent (id) ON DELETE CASCADE);" +
			"INSERT INTO fk_parent VALUES (1);INSERT INTO fk_child VALUES (1, 1)",
		script:       "DELETE FROM fk_parent WHERE id = 1",
		tables:       []string{"fk_parent"},
		transformErr: "is ON DELETE CASCADE",
	},
	{
		// The rows the outer join makes up for t are not backed up.
		name: "outer_join",
		setup: "CREATE TABLE outer_keep (id INT PRIMARY KEY);CREATE TABLE outer_gone (id INT PRIMARY KEY, a INT);" +
			"INSERT INTO outer_keep VALUES (1), (2);INSERT INTO outer_gone VALUES (1, 1)",
		script: "DELETE t FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id",
		tables: []string{"outer_gone"},
	},
	{
		// One statement's two targets share its offset, and truncated to 64
		// characters their backups would get one name.
		name:         "long_names",
		setup:        "CREATE TABLE nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one (id INT PRIMARY KEY);CREATE TABLE nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two (id INT PRIMARY KEY)",
		script:       "DELETE a, b FROM nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one a JOIN nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two b ON a.id = b.id",
		tables:       []string{"nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one"},
		transformErr: "prior backup cannot name the backups of",
	},
	{
		// Shifting u moves one row's old value into another row.
		name:       "swapped_key",
		setup:      "CREATE TABLE swapped_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u));INSERT INTO swapped_key VALUES (1, 10), (2, 20)",
		script:     "UPDATE swapped_key SET u = u + 10 ORDER BY u DESC",
		tables:     []string{"swapped_key"},
		restoreErr: "the UPDATE changes unique key uk_u",
	},
	{
		// Both statements' backups would hold row 1, and its restores would
		// insert it twice.
		name:         "separate_overlap",
		setup:        "CREATE TABLE separate_overlap (id INT PRIMARY KEY, a INT);INSERT INTO separate_overlap VALUES (1, 1)",
		script:       "UPDATE separate_overlap SET a = 2 WHERE id = 1; DELETE FROM separate_overlap WHERE id = 1",
		tables:       []string{"separate_overlap"},
		transformErr: "prior backup cannot back up a DELETE and another statement on the same table",
	},
}

// TestBackupRestoreRoundTrip runs each backup, its statements and their
// restore on TiDB, and requires the tables' rows back.
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
	name := func(table string) (string, string) { return database, table }
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
		catalog := engineCatalog(t, db, database)

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

// restoreSQL builds the restore of backup as Bytebase does: every DML in
// its range that changes the table decides.
func restoreSQL(t *testing.T, script, database, backupDatabase string, backup analysis.BackupStatement, catalog analysis.Catalog) (string, error) {
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
		rows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION", database)
		require.NoError(t, err)
		for rows.Next() {
			var tableName, name, generation, nullable, extra string
			require.NoError(t, rows.Scan(&tableName, &name, &generation, &nullable, &extra))
			column := &metadata.ColumnMetadata{Name: name, Nullable: nullable == "YES"}
			if strings.Contains(strings.ToLower(extra), "on update") {
				column.OnUpdate = "CURRENT_TIMESTAMP"
			}
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
		rows, err = db.Query("SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION", database)
		require.NoError(t, err)
		keys := map[string]*metadata.ForeignKeyMetadata{}
		for rows.Next() {
			var tableName, name, column, referencedTable, referencedColumn, onUpdate, onDelete string
			require.NoError(t, rows.Scan(&tableName, &name, &column, &referencedTable, &referencedColumn, &onUpdate, &onDelete))
			fk := keys[tableName+"."+name]
			if fk == nil {
				fk = &metadata.ForeignKeyMetadata{Name: name, ReferencedTable: referencedTable, OnUpdate: onUpdate, OnDelete: onDelete}
				keys[tableName+"."+name] = fk
				table(tableName).ForeignKeys = append(table(tableName).ForeignKeys, fk)
			}
			fk.Columns = append(fk.Columns, column)
			fk.ReferencedColumns = append(fk.ReferencedColumns, referencedColumn)
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
