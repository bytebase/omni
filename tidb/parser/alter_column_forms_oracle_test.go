package parser

import "testing"

// What TiDB v8.5.5 reads after ALTER [COLUMN] col_name. ALTER COLUMN sets or
// drops a default and nothing else: no SET NOT NULL or DROP NOT NULL, which are
// PostgreSQL's, and no column visibility (view_visibility_forms_oracle_test.go).
// A column's nullability changes through MODIFY or CHANGE. Every verdict was read
// from pingcap/tidb:v8.5.5 on 2026-11-06.
//
// Ref: the "ALTER" arms of AlterTableSpec in TiDB v8.5.5 pkg/parser/parser.y:
// `"ALTER" ColumnKeywordOpt ColumnName "SET" "DEFAULT" SignedLiteral`, `... "SET"
// "DEFAULT" '(' Expression ')'` and `... "DROP" "DEFAULT"`.

type alterColumnFormCase struct {
	sql    string
	accept bool
}

var alterColumnFormCases = []alterColumnFormCase{
	// nullability, which ALTER COLUMN does not change.
	{"ALTER TABLE t ALTER COLUMN b SET NOT NULL", false},
	{"ALTER TABLE t ALTER b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT NULL", false},
	{"ALTER TABLE t ALTER b DROP NOT NULL", false},
	{"alter table t alter column b set not null", false},
	{"ALTER TABLE t ALTER COLUMN b SET NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NULL", false},
	{"ALTER TABLE t ALTER COLUMN b SET NOT", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT", false},
	{"ALTER TABLE t ALTER COLUMN t.b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT NULL, ALTER COLUMN b DROP DEFAULT", false},

	// a column name with no action, or an action with nothing after it.
	{"ALTER TABLE t ALTER COLUMN b SET", false},
	{"ALTER TABLE t ALTER COLUMN b DROP", false},
	{"ALTER TABLE t ALTER COLUMN b", false},
	{"ALTER TABLE t ALTER b", false},

	// the actions ALTER COLUMN has.
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT 1", true},
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL", true},
	{"ALTER TABLE t ALTER b DROP DEFAULT", true},
	{"ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2", true},

	// how a column's nullability does change.
	{"ALTER TABLE t MODIFY b INT NOT NULL", true},
	{"ALTER TABLE t MODIFY COLUMN b INT NULL", true},
	{"ALTER TABLE t CHANGE b b INT NOT NULL", true},
}

func TestAlterColumnFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range alterColumnFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbRelease, accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, c.accept, c.sql, perr)
			}
		})
	}
}

// TestAlterColumnForms holds the parser to the recorded verdicts without the
// server, so a regression fails however the suite is run.
func TestAlterColumnForms(t *testing.T) {
	for _, c := range alterColumnFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
