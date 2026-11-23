package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/doris/internal/testcontainer"
	"github.com/bytebase/omni/driver"
)

const dumpCorpusPath = "testdata/dump/corpus.sql"

// TestDump holds the dump of the corpus to testdata/dump/default.sql. The
// corpus states a key model, a distribution, properties, partitions, inline
// indexes, views, materialized views and a procedure, so the golden is what the
// server states for every shape a schema dump covers.
func TestDump(t *testing.T) {
	d := openTestDriver(t, "omni_dump_corpus")
	execDumpScript(t, d.db, readDumpCorpus(t))
	awaitSyncMaterializedView(t, d.databaseName, "sync_mv_base", "smv_total")
	awaitSyncMaterializedView(t, d.databaseName, "select_mv_base", "select")
	got := dumpScript(t, d)

	golden := filepath.Join("testdata", "dump", "default.sql")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Errorf("Dump differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestDumpRestores dumps the corpus, recreates the database, applies the dump,
// and requires the second dump to be the first: the dump is a restore script,
// so replaying it states the same schema. It also proves the order: the corpus
// names a_child_mv and a_child_view after the objects they read, so name order
// would not restore.
func TestDumpRestores(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_dump_roundtrip")
	execDumpScript(t, d.db, readDumpCorpus(t))
	awaitSyncMaterializedView(t, d.databaseName, "sync_mv_base", "smv_total")
	awaitSyncMaterializedView(t, d.databaseName, "select_mv_base", "select")
	first := dumpScript(t, d)

	for _, stmt := range []string{
		"DROP DATABASE " + quoteBacktick(d.databaseName),
		"CREATE DATABASE " + quoteBacktick(d.databaseName),
		"USE " + quoteBacktick(d.databaseName),
	} {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	execDumpScript(t, d.db, first)
	awaitSyncMaterializedView(t, d.databaseName, "sync_mv_base", "smv_total")
	awaitSyncMaterializedView(t, d.databaseName, "select_mv_base", "select")

	if diff := cmp.Diff(first, dumpScript(t, d)); diff != "" {
		t.Errorf("restoring the dump changes it (-dumped +restored):\n%s", diff)
	}
}

// awaitSyncMaterializedView waits, on a connection a recording does not keep,
// for the table's synchronous materialized views to appear: the statement
// returns before the asynchronous schema change adds them.
func awaitSyncMaterializedView(t *testing.T, dbName, table string, want ...string) {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		deadline := time.Now().Add(3 * time.Minute)
		for {
			pending, err := rollupChangePending(ctx, db, dbName, table)
			if err != nil {
				return err
			}
			if !pending {
				names, err := dumpIndexNames(ctx, db, "DESC "+quoteBacktick(dbName)+"."+quoteBacktick(table)+" ALL")
				if err != nil {
					return err
				}
				have := make(map[string]bool, len(names))
				for _, name := range names {
					have[name] = true
				}
				missing := false
				for _, name := range want {
					if !have[name] {
						missing = true
					}
				}
				if !missing {
					return nil
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("%s does not carry its synchronous materialized views %v", table, want)
			}
			time.Sleep(time.Second)
		}
	})
}

// TestDumpMissing holds a database the connection names but the server does not
// have to the not-found error, rather than to an empty script.
func TestDumpMissing(t *testing.T) {
	d := openTestDriver(t, "omni_dump_missing")
	if _, err := d.db.ExecContext(context.Background(), "DROP DATABASE "+quoteBacktick(d.databaseName)); err != nil {
		t.Fatal(err)
	}
	var dumped strings.Builder
	if err := d.Dump(context.Background(), &dumped, nil); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("Dump of a missing database = %v, want ErrNotFound", err)
	}
}

// TestDumpGaps holds each shape the pinned server's own DDL does not carry, so
// a dump cannot restore it. A check asserts the limitation as it stands, so
// closing the gap is what fails the entry and removes it.
func TestDumpGaps(t *testing.T) {
	for _, gap := range dumpGaps {
		t.Run(gap.name, func(t *testing.T) {
			t.Logf("gap: %s", gap.reason)
			gap.check(t)
		})
	}
}

// dumpGaps lists the shapes a dump cannot round-trip, and why.
var dumpGaps = []struct {
	name   string
	reason string
	check  func(t *testing.T)
}{
	{
		name: "variadic-function",
		reason: "A variadic function's definition does not replay. SHOW FULL FUNCTIONS states its " +
			"signature with the ellipsis against the last type (fv(int...)), and SHOW CREATE FUNCTION " +
			"fv(int...) is a syntax error: the grammar takes the ellipsis as its own argument " +
			"(fv(int, ...)). The accepted form states CREATE ALIAS FUNCTION fv(int...) ..., which the " +
			"parser refuses when applied. The dump therefore omits a variadic function.",
		check: func(t *testing.T) {
			const (
				db       = "omni_dump_gap_variadic"
				function = "dump_gap_arity"
			)
			ctx := context.Background()
			d := openTestDriver(t, db)
			execDumpScript(t, d.db, "CREATE ALIAS FUNCTION "+function+"(INT, ...) WITH PARAMETER(a) AS abs(a)")

			// The catalogue states the compact form, and the accepted argument
			// list still states a definition; only applying it fails.
			rows, err := d.db.QueryContext(ctx, "SHOW FULL FUNCTIONS FROM "+quoteBacktick(db))
			if err != nil {
				t.Fatal(err)
			}
			var signature string
			for rows.Next() {
				columns, err := rows.Columns()
				if err != nil {
					t.Fatal(err)
				}
				values := make([]sql.NullString, len(columns))
				dest := make([]any, len(columns))
				for i := range values {
					dest[i] = &values[i]
				}
				if err := rows.Scan(dest...); err != nil {
					t.Fatal(err)
				}
				if strings.HasPrefix(values[0].String, function+"(") {
					signature = values[0].String
				}
			}
			rows.Close()
			if err := rows.Err(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(signature, "...") {
				t.Errorf("SHOW FULL FUNCTIONS states %q, not the compact variadic form; the entry's reason no longer holds", signature)
			}

			if dumped := dumpScript(t, d); strings.Contains(dumped, function) {
				t.Errorf("the dump carries the variadic function now; remove this entry and test it")
			}

			var statedSignature, statement string
			if err := d.db.QueryRowContext(ctx, "SHOW CREATE FUNCTION "+function+"(int, ...)").Scan(&statedSignature, &statement); err != nil {
				t.Fatalf("SHOW CREATE FUNCTION %s(int, ...): %v", function, err)
			}
			if !strings.Contains(statement, "CREATE ALIAS FUNCTION") {
				t.Errorf("SHOW CREATE FUNCTION stated %q; the entry's reason no longer holds", statement)
			}

			// The form the catalogue states is not a statement the server takes.
			if rows, err := d.db.QueryContext(ctx, "SHOW CREATE FUNCTION "+signature); err == nil {
				rows.Close()
				t.Errorf("SHOW CREATE FUNCTION %s is accepted now; remove this entry and dump the function", signature)
			}

			// Nor does the definition the accepted form states replay: apply it
			// to the same database, once the function it names is gone.
			if _, err := d.db.ExecContext(ctx, "DROP FUNCTION "+function+"(int, ...)"); err != nil {
				t.Fatalf("DROP FUNCTION %s(int, ...): %v", function, err)
			}
			if _, err := d.db.ExecContext(ctx, statement); err == nil {
				t.Errorf("the stated definition replays now (%q); remove this entry and dump the function", statement)
			}
		},
	},
	{
		name: "view-definition-database-qualifier",
		reason: "SHOW CREATE VIEW and SHOW CREATE MATERIALIZED VIEW state every reference qualified " +
			"with the source database (`internal`.`db`.`table`), and the dump writes the statement as " +
			"the server states it, so the dump restores only into a database of the same name: under " +
			"another name its views read the source database.",
		check: func(t *testing.T) {
			const db = "omni_dump_gap_relocate"
			d := openTestDriver(t, db)
			execDumpScript(t, d.db, `CREATE TABLE t (
	k INT NOT NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1");
CREATE MATERIALIZED VIEW mv
BUILD DEFERRED REFRESH COMPLETE ON MANUAL
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k FROM t;
CREATE VIEW v (k) AS SELECT k FROM t`)

			dumped := dumpScript(t, d)
			if !strings.Contains(dumped, "CREATE VIEW `v`") || !strings.Contains(dumped, "CREATE MATERIALIZED VIEW mv") {
				t.Fatalf("the dump leaves out the view or the materialized view:\n%s", dumped)
			}
			if source := "`internal`.`" + db + "`"; !strings.Contains(dumped, source) {
				t.Errorf("the dump no longer states the source database %s, so it may be relocatable; remove this entry and test the restore under another name", source)
			}
		},
	},
	{
		name: "rollup",
		reason: "CREATE TABLE accepts ROLLUP and the table carries the index (DESC t ALL lists r1), " +
			"but nothing states the rollup's definition: SHOW CREATE TABLE has no ROLLUP clause, and " +
			"SHOW CREATE MATERIALIZED VIEW r1 ON t answers no row. The dump writes what the server " +
			"states, so a table declared with a rollup dumps without it.",
		check: func(t *testing.T) {
			const (
				db    = "omni_dump_gap_rollup"
				table = "r"
			)
			d := openTestDriver(t, db)
			execDumpScript(t, d.db, `CREATE TABLE r (
	k INT NOT NULL,
	v INT NULL,
	d DATETIME NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
DISTRIBUTED BY HASH(k) BUCKETS 1
ROLLUP (r1 (v, k))
PROPERTIES("replication_num" = "1")`)
			awaitSyncMaterializedView(t, db, table, "r1")

			ctx := context.Background()
			names, err := dumpIndexNames(ctx, d.db, "DESC "+quoteBacktick(db)+"."+quoteBacktick(table)+" ALL")
			if err != nil {
				t.Fatal(err)
			}
			rollupListed := false
			for _, name := range names {
				if name == "r1" {
					rollupListed = true
				}
			}
			if !rollupListed {
				t.Errorf("the table no longer carries the rollup index %v; the entry's premise no longer holds", names)
			}

			var tableName, viewName, statement string
			err = d.db.QueryRowContext(ctx,
				"SHOW CREATE MATERIALIZED VIEW r1 ON "+quoteBacktick(db)+"."+quoteBacktick(table)).
				Scan(&tableName, &viewName, &statement)
			if err == nil {
				t.Errorf("SHOW CREATE MATERIALIZED VIEW states the rollup now (%q); remove this entry and dump it", statement)
			} else if !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("SHOW CREATE MATERIALIZED VIEW r1 ON %s: %v", table, err)
			}

			if dumped := dumpScript(t, d); strings.Contains(dumped, "ROLLUP") || strings.Contains(dumped, "r1") {
				t.Errorf("the dump carries the rollup now; remove this entry")
			}
		},
	},
}

func readDumpCorpus(t *testing.T) string {
	t.Helper()
	corpus, err := os.ReadFile(dumpCorpusPath)
	if err != nil {
		t.Fatal(err)
	}
	return string(corpus)
}

// dumpScript returns the driver's dump of its database.
func dumpScript(t *testing.T, d *Driver) string {
	t.Helper()
	var dumped strings.Builder
	if err := d.Dump(context.Background(), &dumped, nil); err != nil {
		t.Fatal(err)
	}
	return dumped.String()
}

// execDumpScript runs a script as a client does, cutting a procedure's body by
// its DELIMITER directive.
func execDumpScript(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, command := range splitExecuteCommands(script) {
		if _, err := db.ExecContext(context.Background(), command.Text); err != nil {
			t.Fatalf("%s:\n%v", command.Text, err)
		}
	}
}

// rollupChangePending reports whether the table has a rollup schema change that
// has not finished: a synchronous materialized view is one, and Doris refuses
// another ALTER on the table until it does.
func rollupChangePending(ctx context.Context, db *sql.DB, dbName, table string) (bool, error) {
	query := "SHOW ALTER TABLE ROLLUP FROM " + quoteBacktick(dbName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return false, err
	}
	tableColumn, stateColumn := -1, -1
	for i, column := range columns {
		switch column {
		case "TableName":
			tableColumn = i
		case "State":
			stateColumn = i
		}
	}
	if tableColumn < 0 || stateColumn < 0 {
		return false, fmt.Errorf("%s states %v", query, columns)
	}
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return false, err
		}
		if values[tableColumn].String != table {
			continue
		}
		switch values[stateColumn].String {
		case "FINISHED", "CANCELLED":
		default:
			return true, nil
		}
	}
	return false, rows.Err()
}

// dumpShowCreate returns the statement of a two-column SHOW CREATE form.
func dumpShowCreate(t *testing.T, db *sql.DB, query string) string {
	t.Helper()
	var name, stmt string
	if err := db.QueryRowContext(context.Background(), query).Scan(&name, &stmt); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return stmt
}

// TestDumpDefinitionReads holds the references a definition names to the
// objects in the dumped database, so the dump's order is built from the
// parser's view of the server's statement rather than from its text, and the
// functions it calls, which the dump has to create before the definition
// replays.
func TestDumpDefinitionReads(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		catalog  string
		database string
		want     []string
		wantCall []string
	}{
		{
			name:     "view reads a table",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `internal`.`db`.`t`.`k` AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
		},
		{
			name:     "view reads a view",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `internal`.`db`.`w`.`k` AS `k` FROM `internal`.`db`.`w`",
			catalog:  "internal",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "another database is not the dump's object",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `k` FROM `other`.`t`",
			catalog:  "internal",
			database: "db",
			want:     nil,
		},
		{
			name:     "another catalog with the same database name is not the dump's object",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `k` FROM `external`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     nil,
		},
		{
			name:     "a CTE is not an object",
			stmt:     "CREATE VIEW `v` (k) AS WITH `w` AS (SELECT 1 AS k) SELECT k FROM w",
			catalog:  "internal",
			database: "db",
			want:     nil,
		},
		{
			name:     "a CTE body reads a view",
			stmt:     "CREATE VIEW `v` (k) AS WITH `c` AS (SELECT `internal`.`db`.`w`.`k` AS `k` FROM `internal`.`db`.`w`) SELECT `c`.`k` AS `k` FROM c",
			catalog:  "internal",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "a set operation's operand reads a view",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `internal`.`db`.`w`.`k` AS `k` FROM `internal`.`db`.`w` UNION ALL SELECT `internal`.`db`.`w`.`k` AS `k` FROM `internal`.`db`.`w`",
			catalog:  "internal",
			database: "db",
			want:     []string{"w", "w"},
		},
		{
			name:     "a subquery reads a view",
			stmt:     "CREATE VIEW `v` (k) AS SELECT (SELECT max(`internal`.`db`.`w`.`k`) FROM `internal`.`db`.`w`) AS `k`",
			catalog:  "internal",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "materialized view reads a table",
			stmt:     "CREATE MATERIALIZED VIEW mv (k) BUILD DEFERRED REFRESH COMPLETE ON MANUAL DISTRIBUTED BY HASH(k) BUCKETS 1 AS SELECT `internal`.`db`.`t`.`k` AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
		},
		{
			name:     "a database function's call is qualified with its database",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `db`.`alias_abs`(`internal`.`db`.`t`.`k`) AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
			wantCall: []string{"alias_abs"},
		},
		{
			name:     "a database function's call is qualified with its catalog and database",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `internal`.`db`.`alias_abs`(`internal`.`db`.`t`.`k`) AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
			wantCall: []string{"alias_abs"},
		},
		{
			name:     "a built-in call is bare and is not a routine",
			stmt:     "CREATE VIEW `v` (k) AS SELECT abs(`internal`.`db`.`t`.`k`) AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
		},
		{
			name:     "another database's function is not the dump's",
			stmt:     "CREATE VIEW `v` (k) AS SELECT `other`.`alias_abs`(`internal`.`db`.`t`.`k`) AS `k` FROM `internal`.`db`.`t`",
			catalog:  "internal",
			database: "db",
			want:     []string{"t"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, calls, err := dumpDefinitionReads(test.stmt, test.catalog, test.database)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("reads differ (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantCall, calls); diff != "" {
				t.Errorf("calls differ (-want +got):\n%s", diff)
			}
		})
	}
}

// TestOrderDumpDefinitions holds definitions to a dependency order: a
// definition follows every definition in the set that it reads, and the order
// does not depend on the order the objects were collected in. A reference cycle
// ends in name order rather than dropping a definition.
func TestOrderDumpDefinitions(t *testing.T) {
	objects := []*dumpObject{
		{kind: dumpView, name: "a_child", reads: []string{"z_parent"}},
		{kind: dumpView, name: "z_parent"},
		{kind: dumpMaterializedView, name: "m_child", reads: []string{"m_parent"}},
		{kind: dumpMaterializedView, name: "m_parent"},
		{kind: dumpView, name: "b_cycle", reads: []string{"c_cycle"}},
		{kind: dumpView, name: "c_cycle", reads: []string{"b_cycle"}},
	}
	names := make(map[string]bool, len(objects))
	for _, o := range objects {
		names[o.name] = true
	}
	// The cycle cannot be satisfied; it holds only that no definition is
	// dropped and that the order is stable.
	cycle := map[string]bool{"b_cycle": true, "c_cycle": true}

	order := dumpObjectNames(orderDumpDefinitions(objects))
	if len(order) != len(objects) {
		t.Fatalf("ordered %d of %d definitions: %v", len(order), len(objects), order)
	}
	position := make(map[string]int, len(order))
	for i, name := range order {
		position[name] = i
	}
	for _, o := range objects {
		for _, read := range o.reads {
			if names[read] && !cycle[o.name] && !cycle[read] && position[read] > position[o.name] {
				t.Errorf("%s reads %s, but is ordered before it: %v", o.name, read, order)
			}
		}
	}

	reversed := make([]*dumpObject, 0, len(objects))
	for i := len(objects) - 1; i >= 0; i-- {
		reversed = append(reversed, objects[i])
	}
	if again := dumpObjectNames(orderDumpDefinitions(reversed)); !slices.Equal(again, order) {
		t.Errorf("order depends on the collection order: %v, then %v", order, again)
	}
}

func dumpObjectNames(objects []*dumpObject) []string {
	names := make([]string, 0, len(objects))
	for _, o := range objects {
		names = append(names, o.name)
	}
	return names
}

// TestDumpGlobalFunction holds the dump of a view that calls a global function:
// the dump states the function, with IF NOT EXISTS, before the view, and leaves
// out a global function no definition calls. The script restores twice to the
// same bytes: on the server it came from, where the function is still there,
// and on one without the function, which the test makes by dropping it.
func TestDumpGlobalFunction(t *testing.T) {
	ctx := context.Background()
	dropGlobalFunctions := func() {
		testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
			return testcontainer.DropInstanceObjects(ctx, db)
		})
	}
	dropGlobalFunctions()
	t.Cleanup(dropGlobalFunctions)
	d := openTestDriver(t, "omni_dump_global_function")
	execDumpScript(t, d.db, `CREATE GLOBAL ALIAS FUNCTION dump_global_called(INT) WITH PARAMETER(x) AS abs(x);
CREATE GLOBAL ALIAS FUNCTION dump_global_called(INT, INT) WITH PARAMETER(x, y) AS greatest(x, y);
CREATE GLOBAL ALIAS FUNCTION dump_global_unused(INT) WITH PARAMETER(x) AS abs(x);
CREATE ALIAS FUNCTION dump_local(INT) WITH PARAMETER(x) AS abs(x);
CREATE VIEW calls_functions AS SELECT dump_global_called(-1) AS g, dump_local(-2) AS l`)
	first := dumpScript(t, d)
	for _, want := range []string{
		"-- Global function structure for dump_global_called(int)\n--\nCREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int)",
		"-- Global function structure for dump_global_called(int, int)\n--\nCREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int, int)",
		"-- Function structure for dump_local(int)",
	} {
		if !strings.Contains(first, want) {
			t.Errorf("the dump does not state %q:\n%s", want, first)
		}
	}
	if strings.Contains(first, "dump_global_unused") {
		t.Errorf("the dump states a global function no definition calls:\n%s", first)
	}
	if global, view := strings.Index(first, "dump_global_called"), strings.Index(first, "CREATE VIEW"); global < 0 || view < global {
		t.Errorf("the global function does not precede the view that calls it:\n%s", first)
	}

	restore := func(dropFunctions bool) {
		t.Helper()
		if dropFunctions {
			dropGlobalFunctions()
		}
		for _, stmt := range []string{
			"DROP DATABASE " + quoteBacktick(d.databaseName),
			"CREATE DATABASE " + quoteBacktick(d.databaseName),
			"USE " + quoteBacktick(d.databaseName),
		} {
			if _, err := d.db.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("%s: %v", stmt, err)
			}
		}
		execDumpScript(t, d.db, first)
		if diff := cmp.Diff(first, dumpScript(t, d)); diff != "" {
			t.Errorf("restoring the dump (global functions dropped: %v) changes it (-dumped +restored):\n%s", dropFunctions, diff)
		}
	}
	restore(false)
	restore(true)
}

func TestGlobalFunctionIfNotExists(t *testing.T) {
	tests := []struct{ stmt, want string }{
		{
			stmt: "CREATE GLOBAL ALIAS FUNCTION gf(int) WITH PARAMETER(x) AS abs(`x`)",
			want: "CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS gf(int) WITH PARAMETER(x) AS abs(`x`)",
		},
		{
			stmt: `CREATE GLOBAL AGGREGATE FUNCTION ga(int) RETURNS bigint PROPERTIES ("file"="file:///f.jar", "symbol"="S", "type"="JAVA_UDF")`,
			want: `CREATE GLOBAL AGGREGATE FUNCTION IF NOT EXISTS ga(int) RETURNS bigint PROPERTIES ("file"="file:///f.jar", "symbol"="S", "type"="JAVA_UDF")`,
		},
		{
			stmt: "CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS gf(int) WITH PARAMETER(x) AS abs(`x`)",
			want: "CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS gf(int) WITH PARAMETER(x) AS abs(`x`)",
		},
	}
	for _, test := range tests {
		got, err := globalFunctionIfNotExists(test.stmt)
		if err != nil || got != test.want {
			t.Errorf("globalFunctionIfNotExists(%q) = %q, %v; want %q", test.stmt, got, err, test.want)
		}
	}
	if got, err := globalFunctionIfNotExists("CREATE ALIAS FUNCTION f(int) WITH PARAMETER(x) AS abs(`x`)"); err == nil {
		t.Errorf("globalFunctionIfNotExists of a database function = %q, want an error", got)
	}
}
