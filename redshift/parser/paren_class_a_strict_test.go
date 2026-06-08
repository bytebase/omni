package parser

import (
	"strings"
	"testing"
)

// TestClassAStrictDispatchRejections pins the behavior of the 7 Class A
// `default: return nil, nil` sites flagged in PARSER_DISPATCH_AUDIT.md §2.
// Prior to the fix each site silently returned (nil, nil), which meant the
// caller either swallowed the nil, appended it into a list, or crashed on
// a downstream type assertion. All 7 now raise a syntax error via
// p.syntaxErrorAtCur() in the default branch.
//
// Each sub-test here corresponds to one site and feeds the minimum input
// that used to exercise the silent path.
func TestClassAStrictDispatchRejections(t *testing.T) {
	cases := []struct {
		name       string
		sql        string
		wantErrSub string // substring of err.Error() we expect
	}{
		// set.go:282 — parseGenericSetOrFromCurrent default.
		// `SET foo <bad>` with no TO / = / FROM used to return (nil, nil)
		// and the SET-caller silently produced an empty stmt list.
		{
			name:       "set.go:282 SET var_name missing TO/=/FROM",
			sql:        "SET foo baz",
			wantErrSub: "syntax error",
		},
		// Same site, via ALTER SYSTEM — used to crash on the type assert
		// `setstmt.(*nodes.VariableSetStmt)` inside parseAlterSystemStmt.
		{
			name:       "set.go:282 ALTER SYSTEM SET var_name missing TO/=",
			sql:        "ALTER SYSTEM SET foo baz",
			wantErrSub: "syntax error",
		},

		// extension.go:106 — parseAlterExtensionStmt default.
		// `ALTER EXTENSION foo <bad>` used to silently return nil.
		{
			name:       "extension.go:106 ALTER EXTENSION unknown action",
			sql:        "ALTER EXTENSION foo RENAME TO bar",
			wantErrSub: "syntax error",
		},

		// create_table.go:550 — parseColConstraintElem default.
		// `CREATE TABLE t (c int CONSTRAINT name <bad>)` used to silently
		// drop the constraint clause.
		{
			name:       "create_table.go:550 CONSTRAINT name bad-element",
			sql:        "CREATE TABLE t (c int CONSTRAINT mycon FOREIGN KEY)",
			wantErrSub: "syntax error",
		},

		// alter_misc.go:480 — parseAlterTypeCmd default.
		// Upstream dispatch routes ADD/DROP/ALTER into parseAlterCompositeType
		// → parseAlterTypeCmds, which then loops on ','. Feeding a trailing
		// ', <bad>' forces the second call to parseAlterTypeCmd to hit the
		// default branch. Previously a nil command was appended into
		// AlterTableStmt.Cmds; now this rejects.
		{
			name:       "alter_misc.go:480 ALTER TYPE cmds comma-bad",
			sql:        "ALTER TYPE foo ADD ATTRIBUTE a int, GIBBERISH",
			wantErrSub: "syntax error",
		},

		// utility.go:96 — parseExplainableStmt default.
		// `EXPLAIN <gibberish>` used to produce an ExplainStmt with
		// Query=nil. PG rejects; we do too now.
		{
			name:       "utility.go:96 EXPLAIN non-explainable",
			sql:        "EXPLAIN GRANT ALL ON t TO u",
			wantErrSub: "syntax error",
		},

		// define.go:613 — parseOpclassItem default.
		// Malformed item in CREATE OPERATOR CLASS body used to append a
		// nil item into CreateOpClassStmt.Items.
		{
			name:       "define.go:613 CREATE OPERATOR CLASS bad item",
			sql:        "CREATE OPERATOR CLASS oc DEFAULT FOR TYPE int USING btree AS UNKNOWN 1 foo",
			wantErrSub: "syntax error",
		},

		// define.go:725 — parseOpclassDrop default.
		// Malformed drop item in ALTER OPERATOR FAMILY ... DROP used to
		// append a nil item.
		{
			name:       "define.go:725 ALTER OPERATOR FAMILY DROP bad item",
			sql:        "ALTER OPERATOR FAMILY of USING btree DROP UNKNOWN 1 (int, int)",
			wantErrSub: "syntax error",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Errorf("Parse(%q): expected error containing %q, got %v",
					tc.sql, tc.wantErrSub, err)
			}
		})
	}
}
