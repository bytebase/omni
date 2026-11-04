package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TiDB's standalone statistics statements in the WALK_THROUGH walk.
//
// None creates, drops, or renames a name the walk keeps, so each is lowered
// the way ANALYZE TABLE is: followed, recording nothing. What that has to
// earn is the statement after it. A node the walk does not lower returns nil
// from lower(), which ends the walk there, so a DROP STATS in the middle of a
// change script would hide every error below it -- a change that drops a
// column twice would go unreported because the walk never reached the second
// statement.
//
// DROP STATS and CREATE STATISTICS do claim one thing: every table they name
// has to be there. v8.5.5 answers 1146 for `DROP STATS nosuch` and for
// `CREATE STATISTICS s (CARDINALITY) ON nosuch(a)` (read on 2026-11-04), so
// the walk reports a missing table rather than passing a script that would
// fail. LOAD STATS names a file and DROP STATISTICS a statistic, and neither
// is a name the walk holds.
//
// These tests need no container and run unconditionally.

// statisticsWalkStatements are the forms the walk must follow, written
// against walkSchema's table t and its index i.
var statisticsWalkStatements = []string{
	`ANALYZE TABLE t`,
	`ANALYZE TABLE t INDEX`,
	`ANALYZE TABLE t INDEX i WITH 4 BUCKETS, 0.5 SAMPLERATE`,
	`ANALYZE TABLE t ALL COLUMNS`,
	`ANALYZE TABLE t PREDICATE COLUMNS WITH 1024 BUCKETS`,
	`ANALYZE TABLE t COLUMNS a, b`,
	`ANALYZE TABLE t PARTITION p0 INDEX i WITH 4 CMSKETCH DEPTH`,
	`ANALYZE INCREMENTAL TABLE t INDEX i`,
	`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS`,
	`ANALYZE TABLE t DROP HISTOGRAM ON a`,
	`DROP STATS t`,
	`DROP STATS t GLOBAL`,
	`DROP STATS t PARTITION p0`,
	`LOAD STATS '/tmp/stats.json'`,
	`CREATE STATISTICS s1 (CARDINALITY) ON t(a, b)`,
	`CREATE STATISTICS IF NOT EXISTS s2 (DEPENDENCY) ON t(a, b)`,
	`DROP STATISTICS s1`,
}

func TestWalkThroughFollowsStatisticsStatements(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range statisticsWalkStatements {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing: it changes no name.
			noFinding(t, sql+";", shop)

			// The walk carries on past it, so the error below is still
			// found, and it is reported against statement 2 -- the
			// statement's own index, not the statistics statement's.
			oneFinding(t, sql+";\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
				2, "Column `c` already exists in `t`.", shop)

			// And the names a statement above it set are still there
			// below it, so the walk did not lose its state either.
			noFinding(t, "ALTER TABLE t ADD COLUMN c int;\n"+sql+";\nALTER TABLE t DROP COLUMN c;", shop)
		})
	}
}

// TestWalkThroughStatisticsKeepsTheTableName pins that a statistics statement
// in the middle does not make the walk forget a rename either: the table it
// names is read as a name, not as a new object.
func TestWalkThroughStatisticsKeepsTheTableName(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nDROP STATS u;\nALTER TABLE t ADD COLUMN c int;",
		2, "Table `t` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO u;\nANALYZE TABLE u INDEX i;\nCREATE TABLE u (a int);",
		2, "Table `u` already exists.", shop)
}

// TestWalkThroughStatisticsNeedsItsTable pins the one thing DROP STATS and
// CREATE STATISTICS claim: the table has to be there. A table the change
// dropped above them counts as missing, and a table it made above them
// counts as there.
func TestWalkThroughStatisticsNeedsItsTable(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "DROP STATS nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP STATS t, nosuch;", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE STATISTICS s1 (CARDINALITY) ON nosuch(a);", 0, "`nosuch` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE t;\nDROP STATS t;", 1, "`t` doesn't exist.", shop)
	noFinding(t, "CREATE TABLE u (a int);\nDROP STATS u;\nCREATE STATISTICS s1 (CARDINALITY) ON u(a);", shop)
	// A view is a relation here and not a wrong kind: DROP STATS on a view
	// succeeds on v8.5.5, so the walk must not report one.
	noFinding(t, "DROP STATS v;", shop)
	// The claim is a need, not an end to the walk: the ledger carries the
	// failed name forward and reports it once, as it does for ADMIN CHECK
	// TABLE, rather than ending the walk at the statement that named it.
	oneFinding(t, "DROP STATS nosuch;\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
		0, "`nosuch` doesn't exist.", shop)
	// LOAD STATS and DROP STATISTICS claim nothing, so neither reports a
	// name: the file and the statistic are outside what the walk holds.
	noFinding(t, "LOAD STATS '/tmp/nosuch.json';", shop)
	noFinding(t, "DROP STATISTICS nosuch;", shop)
}
