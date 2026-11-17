package driver

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

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

	if diff := cmp.Diff(first, dumpScript(t, d)); diff != "" {
		t.Errorf("restoring the dump changes it (-dumped +restored):\n%s", diff)
	}
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
		name: "function-definition",
		reason: "Doris states a user-defined function's definition with SHOW CREATE FUNCTION name(types), " +
			"but nothing states the argument types: SHOW FUNCTIONS lists a name only, and " +
			"information_schema.routines and parameters have no row for a UDF. The dump cannot build " +
			"the signature a definition is asked with, so it omits every function.",
		check: func(t *testing.T) {
			const (
				db       = "omni_dump_gap_function"
				function = "dump_gap_abs"
			)
			d := openTestDriver(t, db)
			execDumpScript(t, d.db, "CREATE ALIAS FUNCTION "+function+"(INT) WITH PARAMETER(a) AS abs(a)")

			if dumped := dumpScript(t, d); strings.Contains(dumped, function) {
				t.Errorf("the dump carries a function definition now; remove this entry and test it")
			}

			ctx := context.Background()
			// The definition is stateable once the signature is known: the gap is
			// that nothing states the signature, not that the server keeps it.
			var signature, statement string
			if err := d.db.QueryRowContext(ctx, "SHOW CREATE FUNCTION "+function+"(INT)").Scan(&signature, &statement); err != nil {
				t.Fatalf("SHOW CREATE FUNCTION %s(INT): %v", function, err)
			}
			if !strings.Contains(statement, "CREATE ALIAS FUNCTION") {
				t.Errorf("SHOW CREATE FUNCTION stated %q; the entry's reason no longer holds", statement)
			}

			rows, err := d.db.QueryContext(ctx, "SHOW FUNCTIONS LIKE '"+function+"'")
			if err != nil {
				t.Fatal(err)
			}
			columns, err := rows.Columns()
			rows.Close()
			if err != nil {
				t.Fatal(err)
			}
			if len(columns) != 1 {
				t.Errorf("SHOW FUNCTIONS states %v now, so a signature may be discoverable; remove this entry and dump the function", columns)
			}

			var count int
			if err := d.db.QueryRowContext(ctx,
				"SELECT COUNT(*) FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND SPECIFIC_NAME = ?",
				db, function).Scan(&count); err != nil {
				t.Fatal(err)
			}
			if count != 0 {
				t.Errorf("information_schema.routines lists the function now; remove this entry and dump it")
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
		reason: "CREATE TABLE accepts ROLLUP and the server states no rollup afterwards: SHOW CREATE TABLE " +
			"has no ROLLUP clause, and SHOW PROC lists no index for the table. The dump writes what the " +
			"server states, so a table declared with a rollup dumps without one.",
		check: func(t *testing.T) {
			const db = "omni_dump_gap_rollup"
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

			if stmt := dumpShowCreate(t, d.db, "SHOW CREATE TABLE "+quoteBacktick(db)+".r"); strings.Contains(stmt, "ROLLUP") {
				t.Errorf("SHOW CREATE TABLE states the rollup now; remove this entry and dump it:\n%s", stmt)
			}
			if dumped := dumpScript(t, d); strings.Contains(dumped, "ROLLUP") {
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
// parser's view of the server's statement rather than from its text.
func TestDumpDefinitionReads(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		catalog  string
		database string
		want     []string
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := dumpDefinitionReads(test.stmt, test.catalog, test.database)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("reads differ (-want +got):\n%s", diff)
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
