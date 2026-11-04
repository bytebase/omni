package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// What PLAN REPLAYER does to the rules.
//
// The statement carries upstream's ExplainableStmt, so it reaches ALTER TABLE
// and the rules divide the way they do for EXPLAIN: the ANALYZE form runs the
// change and every rule reports it, while the plain dump writes a dump and no
// rule reports it. Both verdicts are the engine's, measured on TiDB v8.5.5 --
// PLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN c adds the
// column and the plain dump does not (tidb/catalog/plan_replayer_test.go pins
// the catalog's half of the same split).

// TestPlanReplayerOfADDL walks every DDL rule through both forms, as
// explain_test.go walks them through EXPLAIN.
func TestPlanReplayerOfADDL(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, tc := range []struct {
		rule    review.Rule
		alter   string
		message string
	}{
		{review.DisallowDropObject, "ALTER TABLE t DROP COLUMN a", "Drops column `a` of `t`."},
		{review.DisallowRename, "ALTER TABLE t RENAME TO u", "Renames table `t` to `u`."},
		{review.DisallowTruncate, "ALTER TABLE t TRUNCATE PARTITION p0", "Truncates partition `p0` of `t`."},
		{review.DisallowDropConstraint, "ALTER TABLE t DROP PRIMARY KEY", "Drops the primary key of `t`."},
		{review.RequirePrimaryKey, "ALTER TABLE t DROP PRIMARY KEY", "Removes the primary key of `t`, and the change doesn't add one back."},
	} {
		t.Run(string(tc.rule), func(t *testing.T) {
			// The change itself is reported, which is what the rest is
			// measured against.
			if got := run(t, tc.rule, tc.alter+";", review.Change{}, shop); len(got) != 1 || got[0].Message != tc.message {
				t.Fatalf("%s: findings = %+v, want one: %q", tc.alter, got, tc.message)
			}
			// A plain dump changes nothing, so no rule reports it.
			for _, prefix := range []string{
				"PLAN REPLAYER DUMP EXPLAIN ",
				"PLAN REPLAYER DUMP WITH STATS AS OF TIMESTAMP '2024-01-01' EXPLAIN ",
			} {
				if got := run(t, tc.rule, prefix+tc.alter+";", review.Change{}, shop); len(got) != 0 {
					t.Errorf("%s%s: findings = %+v, want none: it writes a dump and makes no change",
						prefix, tc.alter, got)
				}
			}
			// The ANALYZE form runs it, so every rule reports it as it
			// reports the change itself.
			for _, prefix := range []string{
				"PLAN REPLAYER DUMP EXPLAIN ANALYZE ",
				"PLAN REPLAYER DUMP WITH STATS AS OF TIMESTAMP '2024-01-01' EXPLAIN ANALYZE ",
			} {
				got := run(t, tc.rule, prefix+tc.alter+";", review.Change{}, shop)
				if len(got) != 1 || got[0].Message != tc.message {
					t.Errorf("%s%s: findings = %+v, want one: %q", prefix, tc.alter, got, tc.message)
				}
			}
		})
	}
}

// TestPlanReplayerAnalyzeSettlesALaterStatement pins the other direction of
// the statements a rule follows across a change: the dump's ALTER TABLE adds
// the primary key a keyless CREATE TABLE needs, and the plain dump does not.
func TestPlanReplayerAnalyzeSettlesALaterStatement(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, sql := range []string{
		"CREATE TABLE x (a int);\nPLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE x ADD PRIMARY KEY (a);",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t ADD PRIMARY KEY (a);",
		"ALTER TABLE t DROP PRIMARY KEY;\nPLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD PRIMARY KEY (a);",
	} {
		if got := run(t, review.RequirePrimaryKey, sql, review.Change{}, shop); len(got) != 0 {
			t.Errorf("%s\nfindings = %+v, want none: the change ends with a primary key", sql, got)
		}
	}
	sql := "CREATE TABLE x (a int);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE x ADD PRIMARY KEY (a);"
	if got := run(t, review.RequirePrimaryKey, sql, review.Change{}, shop); len(got) != 1 {
		t.Errorf("%s\nfindings = %+v, want one: a plain dump adds no key", sql, got)
	}
}

// TestPlanReplayerOfARowChange pins the rule that judges a row change rather
// than a schema one, which the dump reaches through its DML arms.
func TestPlanReplayerOfARowChange(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, tc := range []struct {
		stmt    string
		message string
	}{
		{"UPDATE t SET a = 1", "`UPDATE` of `t` has no `WHERE` clause."},
		{"DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause."},
	} {
		t.Run(tc.stmt, func(t *testing.T) {
			got := run(t, review.RequireWhere, "PLAN REPLAYER DUMP EXPLAIN ANALYZE "+tc.stmt+";", review.Change{}, shop)
			if len(got) != 1 || got[0].Message != tc.message {
				t.Errorf("the analyzed dump of %s: findings = %+v, want one: %q", tc.stmt, got, tc.message)
			}
			if got := run(t, review.RequireWhere, "PLAN REPLAYER DUMP EXPLAIN "+tc.stmt+";", review.Change{}, shop); len(got) != 0 {
				t.Errorf("the plain dump of %s: findings = %+v, want none: it changes no row", tc.stmt, got)
			}
		})
	}
}

// TestPlanReplayerIsWhatThePriorBackupReads pins the hole this family would
// open in prior backup, for the reason trace_test.go gives: Bytebase's backup
// reads the top-level node, so a wrapped UPDATE is skipped and the change
// runs with no backup.
func TestPlanReplayerIsWhatThePriorBackupReads(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema(), BackupDatabaseExists: true}}
	backup := review.Change{PriorBackup: true}
	const msg = "Prior backup can't back up `PLAN REPLAYER DUMP EXPLAIN ANALYZE` of a statement that changes rows.\n\nRemove `PLAN REPLAYER DUMP EXPLAIN ANALYZE`, or turn off **Prior backup**."
	for _, stmt := range []string{
		"UPDATE t SET a = 1 WHERE id = 1",
		"DELETE FROM t WHERE id = 1",
		"REPLACE INTO t VALUES (1, 2, 3)",
	} {
		got := run(t, review.PriorBackup, "PLAN REPLAYER DUMP EXPLAIN ANALYZE "+stmt+";", backup, shop)
		if len(got) != 1 || got[0].Message != msg {
			t.Errorf("the analyzed dump of %s: findings = %+v, want one: %q", stmt, got, msg)
		}
		// A plain dump runs nothing, so there is nothing to back up.
		if got := run(t, review.PriorBackup, "PLAN REPLAYER DUMP EXPLAIN "+stmt+";", backup, shop); len(got) != 0 {
			t.Errorf("the plain dump of %s: findings = %+v, want none: it changes no row", stmt, got)
		}
	}
	// A plain INSERT only adds rows, so it needs no backup even dumped.
	if got := run(t, review.PriorBackup, "PLAN REPLAYER DUMP EXPLAIN ANALYZE INSERT INTO t VALUES (1, 2, 3);", backup, shop); len(got) != 0 {
		t.Errorf("the analyzed dump of an INSERT: findings = %+v, want none", got)
	}
}

// TestPlanReplayerKeepsTheWalkGoing pins the walk rather than a rule: the
// lowering switch has a fallback that leaves the walk claiming nothing from
// there on, so an unmodelled PlanReplayerStmt would quietly stop the name
// tracking for every statement after it. That holds for the forms that carry
// no statement at all, which is where an unmodelled node is easiest to miss.
//
// LOAD is not among them, and TestPlanReplayerLoadStopsTheWalk says why.
func TestPlanReplayerKeepsTheWalkGoing(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, first := range []string{
		"PLAN REPLAYER DUMP EXPLAIN SELECT 1",
		"PLAN REPLAYER DUMP EXPLAIN SLOW QUERY LIMIT 10",
		"PLAN REPLAYER DUMP EXPLAIN 'sql.txt'",
		"PLAN REPLAYER CAPTURE '123' '456'",
		"PLAN REPLAYER CAPTURE REMOVE '123' '456'",
	} {
		sql := first + ";\nCREATE TABLE x (a int);\nALTER TABLE x DROP COLUMN a;"
		got := run(t, review.DisallowDropObject, sql, review.Change{}, shop)
		if len(got) != 1 || got[0].Message != "Drops column `a` of `x`." {
			t.Errorf("after %s: findings = %+v, want one: %q", first, got, "Drops column `a` of `x`.")
		}
	}
	// And the table a dumped statement names is still checked.
	if got := run(t, review.Syntax, "PLAN REPLAYER DUMP EXPLAIN INSERT INTO nosuch VALUES (1);", review.Change{}, shop); len(got) != 0 {
		t.Errorf("syntax findings = %+v, want none", got)
	}
}

// TestPlanReplayerLoadStopsTheWalk pins the one form that is not a harmless
// wrapper. PLAN REPLAYER LOAD puts back the schemas its dump carries, and the
// dump is a file no rule can read, so the walk cannot keep its names accurate
// past it and drops out instead -- where RESTORE and RECOVER already are, for
// the same reason.
//
// Measured on pingcap/tidb:v8.5.5: a dump taken over `t` and loaded after `t`
// was dropped recreates the table with both its columns, and loading the same
// dump again answers 1050. So a table the walk believes is missing may be one
// the dump created.
func TestPlanReplayerLoadStopsTheWalk(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	// Without the LOAD the walk reports the missing table, which is what
	// makes the silence after it meaningful rather than incidental.
	const alter = "ALTER TABLE nosuch ADD COLUMN b int;"
	if got := run(t, review.WalkThrough, alter, review.Change{}, shop); len(got) != 1 {
		t.Fatalf("%s: findings = %+v, want one: the table is not in the schema", alter, got)
	}
	if got := run(t, review.WalkThrough, "PLAN REPLAYER LOAD '/tmp/a.zip';\n"+alter, review.Change{}, shop); len(got) != 0 {
		t.Errorf("after a LOAD: findings = %+v, want none: the dump may have created the table", got)
	}
	// The other forms read no dump, so they do not excuse it.
	for _, first := range []string{
		"PLAN REPLAYER DUMP EXPLAIN SELECT 1",
		"PLAN REPLAYER DUMP EXPLAIN SLOW QUERY LIMIT 10",
		"PLAN REPLAYER CAPTURE '123' '456'",
	} {
		if got := run(t, review.WalkThrough, first+";\n"+alter, review.Change{}, shop); len(got) != 1 {
			t.Errorf("after %s: findings = %+v, want one: it creates nothing", first, got)
		}
	}
}

// TestPlanReplayerSlowQueryWhereIsAPredicate pins the dump's own WHERE as a
// predicate: `= NULL` there is never true and selects no slow query, so it is
// judged like any other WHERE. The rule collects predicate roots by statement
// type, so a new WHERE field is invisible to it until its statement is named.
func TestPlanReplayerSlowQueryWhereIsAPredicate(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	// The bare predicate is reported, which is the baseline.
	if got := run(t, review.RequireIsNull, "SELECT * FROM t WHERE a = NULL;", review.Change{}, shop); len(got) != 1 {
		t.Fatalf("findings = %+v, want one", got)
	}
	for _, sql := range []string{
		"PLAN REPLAYER DUMP EXPLAIN SLOW QUERY WHERE a = NULL;",
		"PLAN REPLAYER DUMP EXPLAIN ANALYZE SLOW QUERY WHERE a = NULL LIMIT 10;",
		"PLAN REPLAYER DUMP EXPLAIN SLOW QUERY WHERE a = NULL ORDER BY t LIMIT 10;",
	} {
		if got := run(t, review.RequireIsNull, sql, review.Change{}, shop); len(got) != 1 {
			t.Errorf("%s: findings = %+v, want one", sql, got)
		}
	}
	// A WHERE that compares properly is no finding, so the arm is reading
	// the expression rather than reporting on sight.
	if got := run(t, review.RequireIsNull, "PLAN REPLAYER DUMP EXPLAIN SLOW QUERY WHERE a IS NULL;", review.Change{}, shop); len(got) != 0 {
		t.Errorf("IS NULL: findings = %+v, want none", got)
	}
}
