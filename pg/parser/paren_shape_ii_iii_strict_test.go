package parser

import (
	"strings"
	"testing"
)

// TestShapeIIAndIIIStrictRejections pins the behavior of the 13 silent-accept
// sites flagged under Shape II (stale-default struct return) and Shape III
// (partial consume / error-swallow) in PARSER_DISPATCH_AUDIT.md §3.
//
// Prior to the fix each site either:
//   - returned a zero-value &nodes.AlterTableCmd{} (Subtype defaults to 0 ==
//     AT_AddColumn, producing a phantom AddColumn ghost), or
//   - swallowed a real error returned from a helper parser
//     (`result, _ := ...`), or
//   - returned a silent enum default for an unknown keyword.
//
// All 13 now surface a syntax error. Each sub-test feeds the minimal input
// that used to exercise the silent path.
func TestShapeIIAndIIIStrictRejections(t *testing.T) {
	cases := []struct {
		name       string
		sql        string
		wantErrSub string
	}{
		// -------------------------------------------------------------------
		// Shape II — stale-default AlterTableCmd{} returns (alter_table.go)
		// -------------------------------------------------------------------

		// alter_table.go:718 — parseAlterTableCmd default.
		// Unknown alter_table_cmd lead keyword used to produce a phantom
		// AT_AddColumn command appended into Cmds.
		{
			name:       "alter_table.go parseAlterTableCmd unknown lead keyword",
			sql:        "ALTER TABLE t FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:993 — parseAlterColumnAction default.
		// `ALTER TABLE t ALTER COLUMN c <junk>` used to yield the AT_AddColumn
		// ghost. Note: `RESTART` alone (without SET GENERATED) is a valid
		// identity option, so we feed a truly unknown keyword.
		{
			name:       "alter_table.go parseAlterColumnAction unknown action",
			sql:        "ALTER TABLE t ALTER COLUMN c FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:1102 — parseAlterColumnSet default.
		// `ALTER TABLE t ALTER COLUMN c SET <junk>` used to yield the ghost.
		{
			name:       "alter_table.go parseAlterColumnSet unknown SET action",
			sql:        "ALTER TABLE t ALTER COLUMN c SET FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:1491 — parseAlterTableSet SET WITHOUT <junk>.
		// Only CLUSTER / OIDS are legal after WITHOUT.
		{
			name:       "alter_table.go SET WITHOUT unknown",
			sql:        "ALTER TABLE t SET WITHOUT FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:1527 — parseAlterTableSet SET WITH <junk>.
		// Only OIDS is legal after WITH.
		{
			name:       "alter_table.go SET WITH unknown",
			sql:        "ALTER TABLE t SET WITH FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:1536 — parseAlterTableSet outer default.
		// `ALTER TABLE t SET <junk>` fallthrough.
		{
			name:       "alter_table.go SET unknown",
			sql:        "ALTER TABLE t SET FOOBAR",
			wantErrSub: "syntax error",
		},

		// alter_table.go:1827 — parseOneSeqOptElem fallback.
		// Reached via ALTER TABLE t ALTER COLUMN c SET GENERATED … then
		// RESTART would match but `SET <bad-seq-opt>` would not. We trigger
		// it via ALTER TABLE t ALTER COLUMN c RESTART SET <junk> path —
		// actually parseAlterIdentityColumnOption's SET branch calls
		// parseOneSeqOptElem. Feed a SET with an unrecognized seq opt.
		{
			name:       "alter_table.go parseOneSeqOptElem unknown seq opt",
			sql:        "ALTER TABLE t ALTER COLUMN c SET GENERATED ALWAYS SET FOOBAR",
			wantErrSub: "syntax error",
		},

		// -------------------------------------------------------------------
		// Shape III — partial-consume + success
		// -------------------------------------------------------------------

		// trigger.go:313 — parseEnableTrigger (ALTER EVENT TRIGGER tail).
		// Unknown keyword after the trigger name used to silently return
		// TRIGGER_FIRES_ON_ORIGIN (the default).
		{
			name:       "trigger.go parseEnableTrigger unknown keyword",
			sql:        "ALTER EVENT TRIGGER trig FOOBAR",
			wantErrSub: "syntax error",
		},

		// utility.go:223 — parseCallStmt without a following `(`.
		// Used to blindly call parseFuncApplication which advances whatever
		// token follows funcname.
		{
			name:       "utility.go parseCallStmt missing open paren",
			sql:        "CALL foo",
			wantErrSub: "syntax error",
		},

		// database.go:265 — parseSetResetClause SET branch.
		// ALTER DATABASE name SET <junk> previously swallowed parseSetRest
		// error and yielded AlterDatabaseSetStmt with nil Setstmt.
		{
			name:       "database.go parseSetResetClause SET unknown",
			sql:        "ALTER DATABASE db SET 123",
			wantErrSub: "syntax error",
		},

		// database.go:273 — parseSetResetClause RESET branch.
		// ALTER DATABASE name RESET <junk> (non-ident) previously accepted.
		{
			name:       "database.go parseSetResetClause RESET unknown",
			sql:        "ALTER DATABASE db RESET 123",
			wantErrSub: "syntax error",
		},

		// maintenance.go:184 — parseClusterStmt qualified-name error in
		// the CLUSTER '(' opts ')' form.
		{
			name:       "maintenance.go parseClusterStmt paren-opts bad relation",
			sql:        "CLUSTER (VERBOSE) 123",
			wantErrSub: "syntax error",
		},

		// maintenance.go:228 — parseClusterStmt qualified-name error in
		// the CLUSTER opt_verbose name form. parseQualifiedName returning
		// an error used to be swallowed. Note: the name path starts only
		// when isColId() says "we have a name lead"; if the first token
		// is a colid we try parseQualifiedName. A trailing '.' with no
		// continuation name triggers the error return inside
		// parseQualifiedName.
		{
			name:       "maintenance.go parseClusterStmt bareword relation with trailing dot",
			sql:        "CLUSTER VERBOSE t.",
			wantErrSub: "syntax error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected syntax error, got nil\n  SQL: %s", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Fatalf("error %q does not contain %q\n  SQL: %s",
					err.Error(), tc.wantErrSub, tc.sql)
			}
		})
	}
}
