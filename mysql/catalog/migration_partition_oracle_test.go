package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the partition generator (correctness-protocol.md gate 2 + the gate-1
// spine), against the LIVE MySQL engines (5.7 :13307, 8.0 :13306). Two properties:
//
//  1. Apply-correctness: for each (from, to) partition transition the node covers, the generated
//     plan applied to a real `from` database yields a table whose canonical form equals `to`
//     (in both directions). Covers defining partitioning on a fresh/existing table, removing it,
//     and re-partitioning (count / type / boundary / subpartitions).
//  2. Migration idempotence: a partitioned table in its real stored form generates an EMPTY plan
//     (no-op), and the full apply→dump→reload→diff→generate round-trip emits nothing.
//
// The apply assertion (assertPartitionApplyCorrect) is SELF-CONTAINED rather than reusing the
// shared assertApplyCorrect, so it can give each probe its OWN apply database (loadOnePartTable)
// instead of the harness-wide fixed `diffdb`. That keeps these tests hermetic when the whole
// package runs together: the shared apply harness contends on one `diffdb` across the generate-core
// and index apply suites (a pre-existing pooled-`USE` race), and adding more `diffdb` users widens
// that window — a per-probe database sidesteps it entirely. The harness reuses connectOracle /
// NormalizerFor / serverCharsetFor / both / only; it skips cleanly when the engines are unreachable.

// loadOnePartTable wraps a CREATE TABLE in a NAMED database (so each probe gets an isolated apply
// DB, not the shared `diffdb`) whose default charset matches the oracle box, and returns the named
// table. Mirrors loadOneTable but with a caller-chosen database name.
func loadOnePartTable(t *testing.T, dbName, serverCharset, createSQL, table string) (*Catalog, *Table) {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s", dbName, serverCharset, dbName, createSQL)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("LoadSQL failed for %q: %v", createSQL, err)
	}
	for _, db := range cat.Databases() {
		if tt := db.GetTable(table); tt != nil {
			return cat, tt
		}
	}
	t.Fatalf("table %q not found after load of %q", table, createSQL)
	return nil, nil
}

// assertPartitionApplyCorrect builds `from` and `to` from the engine's own readbacks (authentic
// stored forms), generates the plan, applies it to a from-state database, reads the result back,
// and asserts the result canonicalizes equal to `to` (both directions). Each probe uses its own
// database (pa_<slug>), on ONE dedicated connection, so the test is hermetic.
func assertPartitionApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p migrationProbe) {
	t.Helper()

	slug := strings.ReplaceAll(p.id, "-", "_")
	applyDB := "pa_" + slug
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	// One pinned connection for ALL of this probe's statements (readbacks + apply), so the
	// pool never lands a `USE` on a different connection than the next statement.
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	// readbackVia applies a CREATE in a throwaway db on this connection and returns SHOW CREATE.
	readbackVia := func(throwaway, createSQL, table string) (string, bool) {
		for _, s := range []string{
			"DROP DATABASE IF EXISTS " + throwaway,
			fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", throwaway, sc),
			"USE " + throwaway,
			createSQL,
		} {
			if _, err := conn.ExecContext(ctx, s); err != nil {
				t.Logf("[%s] %s readback setup failed (may be expected): %q: %v", o.name, p.id, s, err)
				return "", false
			}
		}
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", throwaway, table))
		if err := row.Scan(&name, &ddl); err != nil {
			t.Logf("[%s] %s SHOW CREATE failed: %v", o.name, p.id, err)
			return "", false
		}
		return ddl, true
	}

	// Build `to` and `from` catalogs from authentic readbacks, wrapped in applyDB so the plan
	// qualifies tables as `applyDB`.`t`.
	var toCat, fromCat *Catalog
	tableInTo := strings.TrimSpace(p.toDDL) != ""
	tableInFrom := strings.TrimSpace(p.fromDDL) != ""

	if tableInTo {
		rb, ok := readbackVia("pto_"+slug, p.toDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
		}
		toCat, _ = loadOnePartTable(t, applyDB, sc, rb, p.table)
	} else {
		toCat = New()
	}
	if tableInFrom {
		rb, ok := readbackVia("pfrom_"+slug, p.fromDDL, p.table)
		if !ok {
			t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
		}
		fromCat, _ = loadOnePartTable(t, applyDB, sc, rb, p.table)
	} else {
		fromCat = New()
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Build the real database in state `from` and apply the plan.
	setup := []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
	}
	if tableInFrom {
		setup = append(setup, p.fromDDL)
	}
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	readback := func(table string) (string, bool) {
		var name, ddl string
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", applyDB, table))
		if err := row.Scan(&name, &ddl); err != nil {
			return "", false
		}
		return ddl, true
	}

	if !tableInTo {
		if _, ok := readback(p.table); ok {
			t.Errorf("[%s] %s: table %s still exists after DROP plan:\n%s", o.name, p.id, p.table, plan.SQL())
		}
		return
	}

	resultRB, ok := readback(p.table)
	if !ok {
		t.Fatalf("[%s] %s: result table %s missing after apply:\n%s", o.name, p.id, p.table, plan.SQL())
	}
	// Load the result into the SAME applyDB so it shares the table identity (database.table) with
	// toCat — loading into a different db name would phantom a whole-table DROP+ADD in the diff.
	resultCat, _ := loadOnePartTable(t, applyDB, sc, resultRB, p.table)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  diff: %s",
			o.name, p.id, plan.SQL(), strings.TrimSpace(resultRB), describeDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeDiff(d))
	}
}

// partitionMigrationProbes enumerates the partition transition FORMS this node covers. Each is
// proven on the listed versions: the generated DDL, applied to a real `from` database, must yield
// a table whose canonical form equals `to`.
func partitionMigrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- Define partitioning on a freshly created table (DiffAdd → CREATE + ALTER PARTITION BY) ----
		{"create-range", "t", "",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)", both()},
		{"create-range-columns", "t", "",
			"CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id,a)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))", both()},
		{"create-list", "t", "",
			"CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3), PARTITION pb VALUES IN (4,5,6))", both()},
		{"create-hash-n", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4", both()},
		{"create-key", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 3", both()},
		{"create-linear-key", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 2", both()},
		{"create-subpart-n", "t", "",
			"CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)", both()},
		{"create-part-comment", "t", "",
			"CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3) COMMENT 'first', PARTITION pb VALUES IN (4,5,6) COMMENT 'second')", both()},

		// ---- Add partitioning to an existing (unpartitioned) table ----
		{"add-partitioning-hash", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4", both()},
		{"add-partitioning-range", "t",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt))",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)", both()},

		// ---- Remove partitioning ----
		{"remove-partitioning", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)", both()},
		{"remove-partitioning-range", "t",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt))", both()},

		// ---- Re-partition (change the scheme on an existing partitioned table) ----
		{"repartition-count", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 8", both()},
		{"repartition-type", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4", both()},
		{"repartition-range-bound", "t",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2015), PARTITION pmax VALUES LESS THAN MAXVALUE)", both()},
		{"repartition-add-range", "t",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)", both()},
		{"repartition-toggle-linear", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 4", both()},

		// ---- A non-partition table change leaves partitioning untouched (partitioned both sides,
		//      only a column is added): the partition generator must emit NOTHING for it. ----
		{"add-column-keep-partitioning", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4", both()},

		// ---- REMOVE PARTITIONING while DROPPING a column the OLD partition function references:
		//      the REMOVE must run (PhasePre) before the DROP COLUMN, else MySQL rejects dropping a
		//      column bound by the live partition function. Regression guard for the ordering fix. ----
		{"remove-partitioning-drop-key-column", "t",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)", both()},

		// ---- KEY ALGORITHM=2 (the stripped default) defined on a fresh table: the generated
		//      PARTITION BY must apply and read back equal to `to` (which folds ALGORITHM=2 away). ----
		{"create-key-algorithm-2", "t", "",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 3", both()},
	}
}

// TestOracle_PartitionApplyCorrectness proves gate 2 for every partition transition: the generated
// DDL, applied to a real `from` database, yields a table whose canonical form equals `to`.
func TestOracle_PartitionApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, p := range partitionMigrationProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				assertPartitionApplyCorrect(t, o, n, p)
			})
		}
	}
}

// TestOracle_PartitionMigrationIdempotence proves the gate-1 spine for partitioning: a partitioned
// table in its real stored form generates an EMPTY plan, and the user-form-vs-stored round-trip
// plan is empty too. A non-empty no-op plan is a partition normalization/ordering bug.
func TestOracle_PartitionMigrationIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range partitionDiffProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				slug := strings.ReplaceAll(probe.id, "-", "_")
				cat, _, rb, ok := o.catalogFromReadback(t, "pgen_idem_"+slug, probe.create, probe.table)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, probe.id)
				}
				diff := DiffWithNormalizer(cat, cat, n)
				plan := GenerateMigrationWithNormalizer(cat, cat, diff, n)
				if plan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY NO-OP PLAN for %s (partition normalization/ordering bug):\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(rb), plan.SQL())
				}
				userCat, _ := loadOneTable(t, serverCharsetFor(o.version), probe.create, probe.table)
				rtDiff := DiffWithNormalizer(userCat, cat, n)
				rtPlan := GenerateMigrationWithNormalizer(userCat, cat, rtDiff, n)
				if rtPlan.SQL() != "" {
					t.Errorf("[%s] NON-EMPTY ROUND-TRIP PLAN (user vs stored) for %s:\n  user:   %s\n  stored: %s\n  plan:\n%s",
						o.name, probe.id, strings.TrimSpace(probe.create), strings.TrimSpace(rb), rtPlan.SQL())
				}
			})
		}
	}
}
