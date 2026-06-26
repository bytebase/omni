package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the partition differ (correctness-protocol.md gates 1 & 2), against the
// LIVE MySQL engines (5.7 :13307, 8.0 :13306). Partitioning is the diff-core idempotence gap the
// breadth node closes: a partitioned table's SHOW CREATE emits a canonical partition clause whose
// surface form differs from the user DDL (per-partition `ENGINE = ...` echo, 8.0's backtick-quoted
// lower-cased expression, `PARTITIONS N` defaulting), and diff-core deliberately deferred it, so
// before this node a partitioned table is NOT idempotent.
//
// Two properties are proven mechanically for every partition variant:
//
//  1. Idempotence (the spine): a schema in its real stored form (the engine's SHOW CREATE
//     readback, reloaded) self-diffs empty, AND the user DDL vs the engine's stored form diffs
//     empty (PartitionChanged=false). A non-empty diff is a partition-canonicalization bug.
//  2. Real changes ARE detected (gate-2 complement): two genuinely different partition specs
//     (and partitioned-vs-not) diff with PartitionChanged=true. Without this, a differ that
//     always returned false would pass the idempotence gate vacuously.
//
// The harness reuses assertDiffEmptyAgainstReadback / connectOracle / serverCharsetFor / both /
// only from the diff oracle tests; it skips cleanly when the engines are unreachable.

// partitionDiffProbes enumerates the single-table partition FORMS that must round-trip empty: a
// user-form DDL whose engine-stored partition clause differs (engine echo, expr quoting, default
// expansion) but must canonicalize equal. Every PARTITION BY kind is covered, on every version
// that supports it.
func partitionDiffProbes() []diffProbe {
	return []diffProbe{
		// RANGE over an expression — the headline expr-canonicalization case (8.0 backtick-quotes
		// + lower-cases `year(`dt`)`, 5.7 keeps `YEAR(dt)`, user wrote `YEAR(dt)`).
		{"range-expr", "CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)", "t", both()},
		// RANGE COLUMNS (single + multi) — /*!50500*/ comment, double-space `RANGE  COLUMNS`.
		{"range-columns", "CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id,a)) PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (MAXVALUE))", "t", both()},
		{"range-columns-multi", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b)) PARTITION BY RANGE COLUMNS(a,b) (PARTITION p0 VALUES LESS THAN (10,10), PARTITION p1 VALUES LESS THAN (MAXVALUE,MAXVALUE))", "t", both()},
		// LIST over an expression.
		{"list-expr", "CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3), PARTITION pb VALUES IN (4,5,6))", "t", both()},
		// LIST COLUMNS with multi-column tuple values IN ((1,1),(2,2)).
		{"list-columns-multi", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b)) PARTITION BY LIST COLUMNS(a,b) (PARTITION p0 VALUES IN ((1,1),(2,2)), PARTITION p1 VALUES IN ((3,3),(4,4)))", "t", both()},
		// HASH with PARTITIONS N (auto-gen defs → `PARTITIONS 4`, no per-partition engine echo).
		{"hash-n", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4", "t", both()},
		// HASH with explicit named partitions (each echoed with ENGINE = InnoDB).
		{"hash-explicit", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) (PARTITION pa, PARTITION pb, PARTITION pc)", "t", both()},
		// LINEAR HASH.
		{"linear-hash", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR HASH (id) PARTITIONS 4", "t", both()},
		// KEY over the PK (no explicit column list → loader fills from PK).
		{"key-implicit", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY () PARTITIONS 3", "t", both()},
		// KEY over an explicit column.
		{"key-explicit", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY (id) PARTITIONS 3", "t", both()},
		// KEY ALGORITHM=2 — the DEFAULT algorithm, STRIPPED from SHOW CREATE (`KEY (id)`), so the
		// explicit user form must fold to the stripped stored form (keyAlgorithmDefault fold).
		{"key-algorithm-2", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 3", "t", both()},
		// KEY ALGORITHM=1 — the LEGACY algorithm, ECHOED by SHOW CREATE via a /*!50611 ALGORITHM = 1
		// */ split comment, so it must round-trip WITHOUT being folded.
		{"key-algorithm-1", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY KEY ALGORITHM=1 (id) PARTITIONS 3", "t", both()},
		// LINEAR KEY.
		{"linear-key", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 2", "t", both()},
		// Per-partition COMMENT (echoed as `COMMENT = '...' ENGINE = ...`).
		{"part-comment", "CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3) COMMENT 'first', PARTITION pb VALUES IN (4,5,6) COMMENT 'second')", "t", both()},
		// Subpartitions via SUBPARTITIONS N (default subpartition names psp).
		{"subpart-n", "CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)", "t", both()},
		// Subpartitions with explicit subpartition definitions (named, parenthesized per partition).
		{"subpart-explicit", "CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) (PARTITION p0 VALUES LESS THAN (2010) (SUBPARTITION s0a, SUBPARTITION s0b), PARTITION p1 VALUES LESS THAN MAXVALUE (SUBPARTITION s1a, SUBPARTITION s1b))", "t", both()},
		// Subpartition BY KEY.
		{"subpart-key", "CREATE TABLE t (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id,d)) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY KEY (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN MAXVALUE)", "t", both()},
		// MyISAM-partitioned table (5.7 only — 8.0 rejects native MyISAM partitioning): the
		// per-partition engine echo is ENGINE = MyISAM, so the empty→table-engine fold must use
		// the table's engine, not a hardcoded InnoDB.
		{"myisam-explicit", "CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, PRIMARY KEY (id,c)) ENGINE=MyISAM PARTITION BY LIST (c) (PARTITION pa VALUES IN (1,2,3), PARTITION pb VALUES IN (4,5,6))", "t", only(MySQL57)},
	}
}

// TestOracle_PartitionDiffIdempotence proves gates 1 & 2 for every partition form: the user DDL
// vs its engine readback diffs EMPTY (PartitionChanged=false), and the stored form self-diffs
// empty, on every supported version.
func TestOracle_PartitionDiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, probe := range partitionDiffProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "part_" + strings.ReplaceAll(probe.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)
			})
		}
	}
}

// partitionChangeCase is a (from, to) pair whose partition specs genuinely differ; the differ
// must report PartitionChanged=true (and the table as DiffModify). This is the gate-2 complement:
// proof the differ is not a vacuous always-false.
type partitionChangeCase struct {
	id       string
	fromDDL  string
	toDDL    string
	table    string
	versions []Version
}

func partitionChangeCases() []partitionChangeCase {
	return []partitionChangeCase{
		// Add partitioning to an unpartitioned table.
		{"add-partitioning",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4",
			"t", both()},
		// Remove partitioning.
		{"remove-partitioning",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"t", both()},
		// Change the partition count.
		{"change-count",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 8",
			"t", both()},
		// Change the partition TYPE (HASH → KEY).
		{"change-type",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4",
			"t", both()},
		// Change a RANGE boundary.
		{"change-range-bound",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2015), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"t", both()},
		// Add a partition (RANGE).
		{"add-range-partition",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION p1 VALUES LESS THAN (2020), PARTITION pmax VALUES LESS THAN MAXVALUE)",
			"t", both()},
		// LINEAR flag flip.
		{"toggle-linear",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY KEY (id) PARTITIONS 4",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY LINEAR KEY (id) PARTITIONS 4",
			"t", both()},
	}
}

// TestOracle_PartitionChangeDetected proves the gate-2 complement: a real partition-spec change
// (loaded from the engines' own readbacks) is reported as PartitionChanged=true with the table
// marked DiffModify. Run from real readbacks so it exercises the same canonical path as the
// idempotence proof.
func TestOracle_PartitionChangeDetected(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, tc := range partitionChangeCases() {
			if !containsVersion(tc.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, tc.id), func(t *testing.T) {
				slug := strings.ReplaceAll(tc.id, "-", "_")
				fromCat, _, _, ok1 := o.catalogFromReadback(t, "pc_from_"+slug, tc.fromDDL, tc.table)
				toCat, _, _, ok2 := o.catalogFromReadback(t, "pc_to_"+slug, tc.toDDL, tc.table)
				if !ok1 || !ok2 {
					t.Skipf("[%s] could not obtain readbacks for %s", o.name, tc.id)
				}
				d := DiffWithNormalizer(fromCat, toCat, n)
				te := findTable(d, tc.table)
				if te == nil {
					t.Fatalf("[%s] %s: expected table %s to be modified, got empty diff", o.name, tc.id, tc.table)
				}
				if !te.PartitionChanged {
					t.Errorf("[%s] %s: expected PartitionChanged=true, got false (table action %s)", o.name, tc.id, te.Action)
				}
			})
		}
	}
}
