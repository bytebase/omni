package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// What TRACE does to the rules.
//
// TRACE runs the statement it traces, so every rule has to see through it:
// a rule that switched on the top-level node would stop checking anything a
// change script wraps in TRACE, which is a rule silently turned off rather
// than a rule that disagrees. TRACE PLAN is the other half -- it traces the
// optimizer and runs nothing -- so a rule that judged what a TRACE PLAN
// carries would report a change the server never makes.
//
// Both verdicts are the engine's, measured on TiDB v8.5.5: TRACE SET
// foreign_key_checks = 0 leaves the session variable at 0 and TRACE PLAN SET
// foreign_key_checks = 0 leaves it at 1 (tidb/catalog/trace_test.go pins the
// catalog's half of the same split).
//
// TRACE carries no DDL -- upstream's TraceableStmt names none -- so the rules
// reachable through it are the ones that judge row changes and predicates.
// The DDL rules are reachable through PLAN REPLAYER instead
// (plan_replayer_test.go).

// TestTraceOfARowChange is the rule the family would silently switch off:
// an UPDATE or DELETE with no WHERE clause, wrapped in each form of TRACE.
func TestTraceOfARowChange(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, tc := range []struct {
		stmt    string
		message string
	}{
		{"UPDATE t SET a = 1", "`UPDATE` of `t` has no `WHERE` clause."},
		{"DELETE FROM t", "`DELETE` from `t` has no `WHERE` clause."},
	} {
		t.Run(tc.stmt, func(t *testing.T) {
			// The change itself is reported, which is what the rest is
			// measured against.
			if got := run(t, review.RequireWhere, tc.stmt+";", review.Change{}, shop); len(got) != 1 || got[0].Message != tc.message {
				t.Fatalf("%s: findings = %+v, want one: %q", tc.stmt, got, tc.message)
			}
			// A TRACE runs it, so the rule reports it exactly as it reports
			// the unwrapped statement.
			for _, prefix := range []string{"TRACE ", "TRACE FORMAT = 'row' ", "TRACE FORMAT = 'json' "} {
				got := run(t, review.RequireWhere, prefix+tc.stmt+";", review.Change{}, shop)
				if len(got) != 1 || got[0].Message != tc.message {
					t.Errorf("%s%s: findings = %+v, want one: %q", prefix, tc.stmt, got, tc.message)
				}
			}
			// A TRACE PLAN runs nothing, so the rule leaves it alone.
			for _, prefix := range []string{"TRACE PLAN ", "TRACE PLAN TARGET = 'estimation' "} {
				if got := run(t, review.RequireWhere, prefix+tc.stmt+";", review.Change{}, shop); len(got) != 0 {
					t.Errorf("%s%s: findings = %+v, want none: TRACE PLAN traces the optimizer and changes no row",
						prefix, tc.stmt, got)
				}
			}
		})
	}
}

// TestTraceOfANullComparison pins a rule that judges a predicate rather than
// a change. `= NULL` is never true wherever it is written, so this rule reads
// through every wrapper, TRACE PLAN included: the statement does not have to
// run for the comparison to be wrong.
func TestTraceOfANullComparison(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	const stmt = "SELECT * FROM t WHERE a = NULL"
	for _, prefix := range []string{"", "TRACE ", "TRACE FORMAT = 'json' ", "TRACE PLAN ", "TRACE PLAN TARGET = 'estimation' "} {
		got := run(t, review.RequireIsNull, prefix+stmt+";", review.Change{}, shop)
		if len(got) != 1 {
			t.Errorf("%s%s: findings = %+v, want one", prefix, stmt, got)
		}
	}
}

// TestTraceIsWhatThePriorBackupReads pins the hole this family would open in
// prior backup. Bytebase's own backup reads the statement's top-level node
// (backend/plugin/parser/tidb/backup.go, extractTables), so a wrapped UPDATE
// reaches its default arm and is skipped -- and the change then runs with no
// backup at all. Before TRACE parsed, the review refused the statement as a
// syntax error and the change never ran; now that it parses, this rule is
// what stands in its place.
func TestTraceIsWhatThePriorBackupReads(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema(), BackupDatabaseExists: true}}
	backup := review.Change{PriorBackup: true}
	const msg = "Prior backup can't back up `TRACE` of a statement that changes rows.\n\nRemove `TRACE`, or turn off **Prior backup**."
	for _, tc := range []struct {
		stmt string
		// backupable is whether the unwrapped statement is one the backup
		// can read. For those, the wrapped statement's finding is the
		// wrapper's and nothing else; REPLACE is refused either way.
		backupable bool
	}{
		{stmt: "UPDATE t SET a = 1 WHERE id = 1", backupable: true},
		{stmt: "DELETE FROM t WHERE id = 1", backupable: true},
		{stmt: "REPLACE INTO t VALUES (1, 2, 3)"},
	} {
		stmt := tc.stmt
		if tc.backupable {
			if got := run(t, review.PriorBackup, stmt+";", backup, shop); len(got) != 0 {
				t.Errorf("%s: findings = %+v, want none", stmt, got)
			}
		}
		for _, prefix := range []string{"TRACE ", "TRACE FORMAT = 'json' "} {
			got := run(t, review.PriorBackup, prefix+stmt+";", backup, shop)
			if len(got) != 1 || got[0].Message != msg {
				t.Errorf("%s%s: findings = %+v, want one: %q", prefix, stmt, got, msg)
			}
		}
		// A TRACE PLAN runs nothing, so there is nothing to back up.
		if got := run(t, review.PriorBackup, "TRACE PLAN "+stmt+";", backup, shop); len(got) != 0 {
			t.Errorf("TRACE PLAN %s: findings = %+v, want none: it changes no row", stmt, got)
		}
	}
	// A plain INSERT only adds rows, so it needs no backup even traced.
	if got := run(t, review.PriorBackup, "TRACE INSERT INTO t VALUES (1, 2, 3);", backup, shop); len(got) != 0 {
		t.Errorf("TRACE INSERT: findings = %+v, want none", got)
	}
}

// TestTraceKeepsTheWalkGoing pins the walk rather than a rule: the lowering
// switch has a fallback that leaves the walk claiming nothing from there on,
// so an unmodelled TraceStmt would quietly stop the name tracking for every
// statement after it. The table a traced INSERT needs is still checked, and
// the statement after the TRACE is still judged.
func TestTraceKeepsTheWalkGoing(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	// The traced INSERT's table has to exist, as an untraced one's does.
	for _, sql := range []string{
		"INSERT INTO nosuch VALUES (1);",
		"TRACE INSERT INTO nosuch VALUES (1);",
	} {
		if got := run(t, review.Syntax, sql, review.Change{}, shop); len(got) != 0 {
			t.Fatalf("%s: syntax findings = %+v, want none", sql, got)
		}
	}
	// A statement after a TRACE is judged as it is after anything else: the
	// walk did not stop at the TRACE.
	const sql = "TRACE SELECT 1;\nCREATE TABLE x (a int);\nALTER TABLE x DROP COLUMN a;"
	got := run(t, review.DisallowDropObject, sql, review.Change{}, shop)
	if len(got) != 1 || got[0].Message != "Drops column `a` of `x`." {
		t.Errorf("%s\nfindings = %+v, want one: %q", sql, got, "Drops column `a` of `x`.")
	}
}
