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

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/starrocks/internal/testcontainer"
)

const dumpCorpusPath = "testdata/dump/corpus.sql"

// dumpCorpusMaterializedViews is the corpus's asynchronous materialized views.
var dumpCorpusMaterializedViews = []string{"a_child_mv", "mv_on_view", "z_parent_mv"}

// TestDump holds the dump of the corpus to testdata/dump/default.sql. The
// corpus states every key model, a distribution, properties, partitions an
// inline index, views and asynchronous materialized views, so the golden is what
// the server states for every shape a schema dump covers.
func TestDump(t *testing.T) {
	d := openTestDriver(t, "omni_dump_corpus")
	execDumpScript(t, d.db, readDumpCorpus(t))
	awaitMaterializedViews(t, d.databaseName, dumpCorpusMaterializedViews)
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
// and requires the second dump to be the first: the dump is a restore script, so
// replaying it states the same schema. It also proves the order: the corpus
// names a_child_mv, a_child_view and aa_union_view after the objects they read,
// so name order would not restore.
func TestDumpRestores(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_dump_roundtrip")
	execDumpScript(t, d.db, readDumpCorpus(t))
	awaitMaterializedViews(t, d.databaseName, dumpCorpusMaterializedViews)
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
	awaitMaterializedViews(t, d.databaseName, dumpCorpusMaterializedViews)

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

// TestDumpGaps holds each shape the pinned server's own DDL does not carry, so a
// dump cannot restore it. A check asserts the limitation as it stands, so
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
		name: "view-definition-database-qualifier",
		reason: "SHOW CREATE VIEW and SHOW CREATE MATERIALIZED VIEW state every reference qualified " +
			"with the source database (`db`.`table`), and the dump writes the statement as the server " +
			"states it, so the dump restores only into a database of the same name: under another " +
			"name its views read the source database. The server states no catalog, because a " +
			"definition is read in the session's own catalog.",
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
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k FROM t;
CREATE VIEW v (k) AS SELECT k FROM t`)
			awaitMaterializedViews(t, db, []string{"mv"})

			dumped := dumpScript(t, d)
			if !strings.Contains(dumped, "CREATE VIEW `v`") || !strings.Contains(dumped, "CREATE MATERIALIZED VIEW `mv`") {
				t.Fatalf("the dump leaves out the view or the materialized view:\n%s", dumped)
			}
			if source := "`" + db + "`.`t`"; !strings.Contains(dumped, source) {
				t.Errorf("the dump no longer states the source database %s, so it may be relocatable; remove this entry and test the restore under another name", source)
			}
		},
	},
	{
		name: "rollup",
		reason: "A rollup is a table's synchronous materialized view: CREATE TABLE accepts ROLLUP and the " +
			"table carries the index (DESC t ALL lists r1), but no DDL states its definition. SHOW CREATE " +
			"TABLE has no ROLLUP clause, and SHOW CREATE MATERIALIZED VIEW takes only a name and answers " +
			"that r1 is not a table: the server has no ON form for a rollup, which is a Doris statement. " +
			"The one definition the server states is SHOW MATERIALIZED VIEWS' text column, a statement the " +
			"server takes, but not the rollup that was declared: it adds a GROUP BY, so it restores an " +
			"aggregate-key rollup where the declared index of a duplicate-key table is duplicate-key, under " +
			"the same name and columns. The dump writes what the server states, so a table declared with a " +
			"rollup dumps without it.",
		check: func(t *testing.T) {
			const (
				db     = "omni_dump_gap_rollup"
				apply  = "omni_dump_gap_rollup_apply"
				table  = "r"
				rollup = "r1"
			)
			ctx := context.Background()
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
			awaitRollup(t, db, table, rollup)
			// The definition the server states for the rollup, which the check
			// below applies. It is read in a recording only: the listing it comes
			// from carries a server id, which a tape cannot hold.
			const statedDefinition = "create materialized view r1 as select v,k from r group by v,k"
			awaitRollupDefinition(t, db, rollup, statedDefinition)

			names, err := dumpIndexNames(ctx, d.db, "DESC "+quoteBacktick(db)+"."+quoteBacktick(table)+" ALL")
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Contains(names, rollup) {
				t.Errorf("the table no longer carries the rollup index %v; the entry's premise no longer holds", names)
			}
			declared, err := dumpIndexKeyType(ctx, d.db, db, table, rollup)
			if err != nil {
				t.Fatal(err)
			}
			if declared != "DUP_KEYS" {
				t.Errorf("the declared rollup is %s, not DUP_KEYS; the entry's reason no longer holds", declared)
			}

			// The name-only form SHOW CREATE MATERIALIZED VIEW takes answers that
			// a rollup is not a table, and the ON form Doris states is not a
			// statement this server takes.
			if rows, err := d.db.QueryContext(ctx, "SHOW CREATE MATERIALIZED VIEW "+quoteBacktick(db)+"."+quoteBacktick(rollup)); err == nil {
				rows.Close()
				t.Errorf("SHOW CREATE MATERIALIZED VIEW states the rollup now; remove this entry and dump it")
			}
			if rows, err := d.db.QueryContext(ctx, "SHOW CREATE MATERIALIZED VIEW "+quoteBacktick(rollup)+" ON "+quoteBacktick(db)+"."+quoteBacktick(table)); err == nil {
				rows.Close()
				t.Errorf("SHOW CREATE MATERIALIZED VIEW ... ON is a statement the server takes now; remove this entry and dump it")
			} else if !strings.Contains(err.Error(), "syntax error") {
				t.Fatalf("SHOW CREATE MATERIALIZED VIEW ... ON: %v, want a syntax error", err)
			}

			dumped := dumpScript(t, d)
			if strings.Contains(dumped, rollup) || strings.Contains(dumped, "ROLLUP") {
				t.Errorf("the dump carries the rollup now; remove this entry")
			}

			// Apply the definition the server states for the rollup, which is
			// the materialized-view form with a GROUP BY, to a table of the
			// declared table's shape in a database of its own. The declared
			// index is duplicate-key; the stated form restores an aggregate-key
			// rollup.
			if _, err := d.db.ExecContext(ctx, "CREATE DATABASE "+quoteBacktick(apply)); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if _, err := d.db.ExecContext(context.Background(), "DROP DATABASE "+quoteBacktick(apply)); err != nil {
					t.Errorf("drop %s: %v", apply, err)
				}
			}()
			// The definition names its base table unqualified, so it is applied
			// in the database that holds the table it creates the rollup on.
			for _, stmt := range []string{
				"USE " + quoteBacktick(apply),
				"CREATE TABLE " + quoteBacktick(apply) + ".`r` (k INT NOT NULL, v INT NULL, d DATETIME NULL) ENGINE=OLAP DUPLICATE KEY(k) DISTRIBUTED BY HASH(k) BUCKETS 1 PROPERTIES(\"replication_num\" = \"1\")",
				statedDefinition,
				"USE " + quoteBacktick(db),
			} {
				if _, err := d.db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			awaitRollup(t, apply, table, rollup)
			restored, err := dumpIndexKeyType(ctx, d.db, apply, table, rollup)
			if err != nil {
				t.Fatal(err)
			}
			if restored != "AGG_KEYS" {
				t.Errorf("the stated definition now restores a %s rollup; if it restores the declared key model, remove this entry and dump the rollup", restored)
			}
		},
	},
	{
		name: "function-definition",
		reason: "A user-defined function is stated by no DDL. SHOW CREATE FUNCTION is not a statement the " +
			"server takes, and SHOW FULL FUNCTIONS states a signature, a return type and a properties " +
			"object whose fid, md5 and language are the catalogue's, not a CREATE FUNCTION. The dump " +
			"writes what the server states, so it does not carry a function; fails if the server states " +
			"one, or if the dump starts carrying it.",
		check: func(t *testing.T) {
			const (
				db       = "omni_dump_gap_function"
				function = "dump_gap_fn"
			)
			ctx := context.Background()
			d := openTestDriver(t, db)
			// The object file is read, not run: the front end only checksums it,
			// so a file the pinned image carries stands in for a Python archive.
			execDumpScript(t, d.db, `CREATE FUNCTION `+function+`(INT) RETURNS INT PROPERTIES("type" = "Python", "file" = "file:///data/deploy/starrocks/fe/conf/fe.conf", "symbol" = "f", "input_type" = "INT")`)

			if dumped := dumpScript(t, d); strings.Contains(dumped, function) {
				t.Errorf("the dump carries the function now; remove this entry and test it")
			}

			// No statement states the definition: the name-only form is not a
			// statement the server takes.
			rows, err := d.db.QueryContext(ctx, "SHOW CREATE FUNCTION "+function+"(INT)")
			if err == nil {
				rows.Close()
				t.Errorf("SHOW CREATE FUNCTION is a statement the server takes now; remove this entry and dump the function")
			} else if want := "No viable statement for input 'SHOW CREATE FUNCTION'"; !strings.Contains(err.Error(), want) {
				t.Fatalf("SHOW CREATE FUNCTION %s(INT): %v, want %q", function, err, want)
			}

			// The function exists: the same signature is the one the catalogue
			// would state, and it is what a drop takes.
			if _, err := d.db.ExecContext(ctx, "DROP FUNCTION "+function+"(INT)"); err != nil {
				t.Fatalf("DROP FUNCTION %s(INT): %v", function, err)
			}
		},
	},
	{
		name: "global-function-definition",
		reason: "A view that calls a global function dumps without it. SHOW CREATE VIEW states the call " +
			"bare (dump_gap_global(1)), and no statement states a global function's definition: SHOW " +
			"CREATE GLOBAL FUNCTION is not a statement the server takes, and SHOW FULL GLOBAL FUNCTIONS " +
			"states the catalogue's properties, as for a database's function. Doris states the definition, " +
			"and its dump carries the global functions a view calls; this dump restores such a view only " +
			"on a server that has the function. Fails if the server states the definition.",
		check: func(t *testing.T) {
			const (
				db       = "omni_dump_gap_global_function"
				function = "dump_gap_global"
			)
			ctx := context.Background()
			dropGlobalFunctions := func() {
				testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
					return testcontainer.DropInstanceObjects(ctx, db)
				})
			}
			dropGlobalFunctions()
			t.Cleanup(dropGlobalFunctions)
			d := openTestDriver(t, db)
			execDumpScript(t, d.db, `CREATE GLOBAL FUNCTION `+function+`(INT) RETURNS INT PROPERTIES("type" = "Python", "file" = "file:///data/deploy/starrocks/fe/conf/fe.conf", "symbol" = "f", "input_type" = "INT");
CREATE VIEW calls_global (g) AS SELECT `+function+`(1)`)

			dumped := dumpScript(t, d)
			if !strings.Contains(dumped, function+"(1)") {
				t.Fatalf("the dump does not state the view's call to %s:\n%s", function, dumped)
			}
			if strings.Contains(dumped, "CREATE GLOBAL FUNCTION") {
				t.Errorf("the dump carries the global function now; remove this entry and test it")
			}

			rows, err := d.db.QueryContext(ctx, "SHOW CREATE GLOBAL FUNCTION "+function+"(INT)")
			if err == nil {
				rows.Close()
				t.Errorf("SHOW CREATE GLOBAL FUNCTION is a statement the server takes now; remove this entry and dump the function")
			} else if want := "No viable statement for input 'SHOW CREATE GLOBAL'"; !strings.Contains(err.Error(), want) {
				t.Fatalf("SHOW CREATE GLOBAL FUNCTION %s(INT): %v, want %q", function, err, want)
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

// execDumpScript runs a script as a client does, cutting a statement's body by
// its DELIMITER directive.
func execDumpScript(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, command := range splitExecuteCommands(script) {
		if _, err := db.ExecContext(context.Background(), command.Text); err != nil {
			t.Fatalf("%s:\n%v", command.Text, err)
		}
	}
}

// awaitRollup waits, on a connection a recording does not keep, for the table's
// rollup to appear: the statement returns before the asynchronous schema change
// adds it.
func awaitRollup(t *testing.T, dbName, table, rollup string) {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		deadline := time.Now().Add(3 * time.Minute)
		for {
			names, err := dumpIndexNames(ctx, db, "DESC "+quoteBacktick(dbName)+"."+quoteBacktick(table)+" ALL")
			if err != nil {
				return err
			}
			if slices.Contains(names, rollup) {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("%s does not carry the rollup %s", table, rollup)
			}
			time.Sleep(time.Second)
		}
	})
}

// awaitMaterializedViews waits, on a connection a recording does not keep, for
// the database's materialized-view listing to name every view in want. The
// front end's listing fails with a front-end RPC error until a back end reports
// the views, which a fresh server does some seconds after an asynchronous
// materialized view (or a rollup) is created.
func awaitMaterializedViews(t *testing.T, dbName string, want []string) {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		deadline := time.Now().Add(3 * time.Minute)
		for {
			have, err := queryMaterializedViewNames(ctx, db, dbName)
			if err == nil {
				missing := false
				for _, name := range want {
					if !slices.Contains(have, name) {
						missing = true
					}
				}
				if !missing {
					return nil
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("%s does not list its materialized views %v", dbName, want)
			}
			time.Sleep(time.Second)
		}
	})
}

// queryMaterializedViewNames returns the names the database's materialized-view
// listing states, or the server's error for a listing it cannot answer yet.
func queryMaterializedViewNames(ctx context.Context, db *sql.DB, dbName string) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME", dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// awaitRollupDefinition holds, in a recording, the definition SHOW MATERIALIZED
// VIEWS states for a rollup to want. Its answer carries a server id that a new
// server does not reuse, so a tape cannot hold it.
func awaitRollupDefinition(t *testing.T, dbName, rollup, want string) {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		deadline := time.Now().Add(3 * time.Minute)
		for {
			got, err := statedRollupDefinition(ctx, db, dbName, rollup)
			if err == nil {
				if got != want {
					return fmt.Errorf("%s states the rollup %s's definition as %q, not %q", dbName, rollup, got, want)
				}
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("%s states no definition for the rollup %s: %w", dbName, rollup, err)
			}
			time.Sleep(time.Second)
		}
	})
}

// statedRollupDefinition returns the definition SHOW MATERIALIZED VIEWS states
// for a rollup: its text column, which the listing answers as soon as the index
// exists.
func statedRollupDefinition(ctx context.Context, db *sql.DB, dbName, rollup string) (string, error) {
	query := "SHOW MATERIALIZED VIEWS FROM " + quoteBacktick(dbName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	name, text := -1, -1
	for i, column := range columns {
		switch column {
		case "name":
			name = i
		case "text":
			text = i
		}
	}
	if name < 0 || text < 0 {
		return "", fmt.Errorf("%s states %v", query, columns)
	}
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return "", err
		}
		if values[name].String == rollup {
			return strings.TrimSpace(values[text].String), nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("%s states no %s", query, rollup)
}

// dumpIndexNames returns the first column of a description of a table's
// indexes: the index name, repeated for each column it holds.
func dumpIndexNames(ctx context.Context, q interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, query string) ([]string, error) {
	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		name := values[0].String
		if name == "" || (len(names) > 0 && names[len(names)-1] == name) {
			continue
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// dumpIndexKeyType returns the key model DESC ... ALL states for one index of a
// table: DUP_KEYS, AGG_KEYS, UNIQUE_KEYS or PRIMARY_KEYS.
func dumpIndexKeyType(ctx context.Context, db *sql.DB, dbName, table, index string) (string, error) {
	query := "DESC " + quoteBacktick(dbName) + "." + quoteBacktick(table) + " ALL"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	name, keyType := -1, -1
	for i, column := range columns {
		switch column {
		case "IndexName":
			name = i
		case "IndexKeysType":
			keyType = i
		}
	}
	if name < 0 || keyType < 0 {
		return "", fmt.Errorf("%s states %v", query, columns)
	}
	current := ""
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return "", err
		}
		if values[name].String != "" {
			current = values[name].String
		}
		if current == index {
			return values[keyType].String, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("%s does not list the index %s", query, index)
}

// TestDumpDefinitionReads holds the references a definition names to the objects
// in the dumped database, so the dump's order is built from the parser's view of
// the server's statement rather than from its text.
func TestDumpDefinitionReads(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		database string
		want     []string
	}{
		{
			name:     "view reads a table",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT `db`.`t`.`k`\nFROM `db`.`t`",
			database: "db",
			want:     []string{"t"},
		},
		{
			name:     "view reads a view",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT `db`.`w`.`k`\nFROM `db`.`w`",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "view without SECURITY reads a table",
			stmt:     "CREATE VIEW `v` (`k`) AS SELECT `db`.`t`.`k` FROM `db`.`t`",
			database: "db",
			want:     []string{"t"},
		},
		{
			name:     "another database is not the dump's object",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT `k`\nFROM `other`.`t`",
			database: "db",
			want:     nil,
		},
		{
			name:     "a CTE is not an object",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS WITH `c` (`k`) AS (SELECT 1 AS `k`) SELECT `c`.`k`\nFROM `c`",
			database: "db",
			want:     nil,
		},
		{
			name:     "a CTE body reads a view",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS WITH `c` (`k`) AS (SELECT `db`.`w`.`k`\nFROM `db`.`w`) SELECT `c`.`k`\nFROM `c`",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "a set operation's operand reads a view",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS (SELECT `db`.`w`.`k`\nFROM `db`.`w`) UNION ALL (SELECT `db`.`w`.`k`\nFROM `db`.`w`)",
			database: "db",
			want:     []string{"w", "w"},
		},
		{
			name:     "a scalar subquery reads a view",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT (SELECT max(`db`.`w`.`k`)\nFROM `db`.`w`) AS `k`",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "a FROM subquery reads a view",
			stmt:     "CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT `x`.`k`\nFROM (SELECT `db`.`w`.`k`\nFROM `db`.`w`) `x`",
			database: "db",
			want:     []string{"w"},
		},
		{
			name:     "materialized view reads a view",
			stmt:     "CREATE MATERIALIZED VIEW `mv` (`k`, `s`)\nDISTRIBUTED BY HASH(`k`) BUCKETS 1 \nREFRESH MANUAL\nAS SELECT `w`.`k`, sum(`w`.`s`) AS `s`\nFROM `db`.`w`\nGROUP BY `w`.`k`",
			database: "db",
			want:     []string{"w"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := dumpDefinitionReads(test.stmt, test.database)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("reads differ (-want +got):\n%s", diff)
			}
		})
	}
}

// TestOrderDumpDefinitions holds definitions to a dependency order: a definition
// follows every definition in the set that it reads, and the order does not
// depend on the order the objects were collected in. A reference cycle ends in
// name order rather than dropping a definition.
func TestOrderDumpDefinitions(t *testing.T) {
	objects := []*dumpObject{
		{kind: dumpView, name: "a_child", reads: []string{"z_parent"}},
		{kind: dumpMaterializedView, name: "z_parent"},
		{kind: dumpMaterializedView, name: "m_child", reads: []string{"m_parent"}},
		{kind: dumpView, name: "m_parent"},
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
