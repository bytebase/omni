package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// What EXPLAIN does to the rules, now that TiDB's EXPLAIN carries ALTER TABLE.
//
// The rules divide in two. One kind walks the statement tree with ast.Inspect,
// and ast.Walk follows what an EXPLAIN carries, so a plain EXPLAIN would hand
// them a change that is not made: EXPLAIN ALTER TABLE t DROP COLUMN a reports
// a plan and drops nothing. Those stop at a non-analyzed ExplainStmt. The
// other kind reads the statement's top-level node, so an EXPLAIN ANALYZE --
// which does run the change -- would hide it; those read through executed().
//
// Both verdicts are the engine's: tidb/parser/explain_forms_test.go measures
// on TiDB v8.5.5 that ANALYZE makes the change and that a plain EXPLAIN does
// not.
func TestExplainOfADDL(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	// Each rule, the statement that triggers it, and the message it gives.
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
			// The change itself is reported, which is what these assertions
			// are measured against.
			if got := run(t, tc.rule, tc.alter+";", review.Change{}, shop); len(got) != 1 || got[0].Message != tc.message {
				t.Fatalf("%s: findings = %+v, want one: %q", tc.alter, got, tc.message)
			}
			// A plain EXPLAIN of it changes nothing, so no rule reports it.
			for _, prefix := range []string{"EXPLAIN ", "EXPLAIN FORMAT = 'brief' ", "EXPLAIN FORMAT = DOT "} {
				if got := run(t, tc.rule, prefix+tc.alter+";", review.Change{}, shop); len(got) != 0 {
					t.Errorf("%s%s: findings = %+v, want none: it reports a plan and makes no change",
						prefix, tc.alter, got)
				}
			}
			// EXPLAIN ANALYZE runs it, so every rule reports it as it
			// reports the change itself.
			for _, prefix := range []string{"EXPLAIN ANALYZE ", "EXPLAIN ANALYZE FORMAT = 'brief' "} {
				got := run(t, tc.rule, prefix+tc.alter+";", review.Change{}, shop)
				if len(got) != 1 || got[0].Message != tc.message {
					t.Errorf("%s%s: findings = %+v, want one: %q", prefix, tc.alter, got, tc.message)
				}
			}
		})
	}
}

// TestExplainAnalyzeSettlesALaterStatement pins the other direction of the
// statements a rule follows across a change: an analyzed ALTER TABLE adds the
// primary key a keyless CREATE TABLE needs, and an analyzed DROP of a key is
// settled by a later ADD.
func TestExplainAnalyzeSettlesALaterStatement(t *testing.T) {
	shop := []review.Target{{Schema: walkSchema()}}
	for _, sql := range []string{
		"CREATE TABLE x (a int);\nEXPLAIN ANALYZE ALTER TABLE x ADD PRIMARY KEY (a);",
		"EXPLAIN ANALYZE ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t ADD PRIMARY KEY (a);",
		"ALTER TABLE t DROP PRIMARY KEY;\nEXPLAIN ANALYZE ALTER TABLE t ADD PRIMARY KEY (a);",
	} {
		if got := run(t, review.RequirePrimaryKey, sql, review.Change{}, shop); len(got) != 0 {
			t.Errorf("%s\nfindings = %+v, want none: the change ends with a primary key", sql, got)
		}
	}
	// A plain EXPLAIN settles nothing, because it adds no key.
	sql := "CREATE TABLE x (a int);\nEXPLAIN ALTER TABLE x ADD PRIMARY KEY (a);"
	if got := run(t, review.RequirePrimaryKey, sql, review.Change{}, shop); len(got) != 1 {
		t.Errorf("%s\nfindings = %+v, want one: a plain EXPLAIN adds no key", sql, got)
	}
}
