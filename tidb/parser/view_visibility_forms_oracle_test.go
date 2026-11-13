package parser

import "testing"

// What TiDB v8.5.5 reads of MySQL's column visibility, ALTER VIEW and view check
// options. TiDB has no column visibility: VISIBLE and INVISIBLE are index options
// and the tail of `ALTER INDEX`, nothing else. It has no ALTER VIEW in any
// spelling. Its view check option names CASCADED or LOCAL and spells out CHECK
// OPTION. Every verdict was read from pingcap/tidb:v8.5.5 on 2026-11-06.
//
// Ref: TiDB v8.5.5 pkg/parser/parser.y: ColumnOption has no "VISIBLE" or
// "INVISIBLE" arm; AlterTableSpec has `"ALTER" "INDEX" Identifier IndexInvisible`
// and no column visibility arm; ViewCheckOption is `"WITH" "CASCADED" "CHECK"
// "OPTION" | "WITH" "LOCAL" "CHECK" "OPTION"`; no rule starts with "ALTER"
// followed by "VIEW", ViewAlgorithm, ViewDefiner or ViewSQLSecurity.

type viewVisibilityFormCase struct {
	sql    string
	accept bool
}

var viewVisibilityFormCases = []viewVisibilityFormCase{
	// column visibility in a column definition.
	{"CREATE TABLE t (a INT, b INT INVISIBLE)", false},
	{"CREATE TABLE t (a INT, b INT VISIBLE)", false},
	{"CREATE TABLE t (a INT, b INT NOT NULL INVISIBLE DEFAULT 1)", false},
	{"CREATE TABLE t (a INT PRIMARY KEY INVISIBLE)", false},
	{"CREATE TABLE t (a INT, b INT AS (a + 1) VIRTUAL INVISIBLE)", false},
	{"CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL VISIBLE)", false},
	{"ALTER TABLE t ADD COLUMN b INT INVISIBLE", false},
	{"ALTER TABLE t MODIFY b INT INVISIBLE", false},
	{"ALTER TABLE t CHANGE b c INT VISIBLE", false},

	// column visibility in ALTER COLUMN.
	{"ALTER TABLE t ALTER COLUMN b SET VISIBLE", false},
	{"ALTER TABLE t ALTER COLUMN b SET INVISIBLE", false},
	{"ALTER TABLE t ALTER b SET INVISIBLE", false},
	{"ALTER TABLE t ALTER t.b SET INVISIBLE", false},
	{"ALTER TABLE t ALTER b SET DEFAULT 1", true},
	{"ALTER TABLE t ALTER b DROP DEFAULT", true},
	{"ALTER TABLE t ALTER b SET", false},
	{"ALTER TABLE t ALTER b DROP", false},
	{"ALTER TABLE t ALTER b", false},

	// index visibility stays.
	{"ALTER TABLE t ALTER INDEX i VISIBLE", true},
	{"ALTER TABLE t ALTER INDEX i INVISIBLE", true},
	{"ALTER TABLE t ALTER INDEX i", false},
	{"CREATE TABLE t (a INT, INDEX i (a) INVISIBLE)", true},
	{"CREATE INDEX i ON t (a) INVISIBLE", true},

	// the words as names.
	{"CREATE TABLE visible (invisible INT)", true},
	{"SELECT 1 FROM t AS visible", true},

	// view check options.
	{"CREATE VIEW v AS SELECT 1 WITH CHECK OPTION", false},
	{"CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION", true},
	{"CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION", true},
	{"CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK", false},
	{"CREATE VIEW v AS SELECT 1 WITH LOCAL", false},
	{"CREATE VIEW v AS SELECT 1 WITH CASCADED", false},
	{"CREATE OR REPLACE ALGORITHM = MERGE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v (a) AS SELECT 1 WITH LOCAL CHECK OPTION", true},

	// ALTER VIEW in every spelling.
	{"ALTER VIEW v AS SELECT 1", false},
	{"ALTER VIEW v (a) AS SELECT 1 WITH CASCADED CHECK OPTION", false},
	{"ALTER ALGORITHM = MERGE VIEW v AS SELECT 1", false},
	{"ALTER DEFINER = CURRENT_USER VIEW v AS SELECT 1", false},
	{"ALTER SQL SECURITY INVOKER VIEW v AS SELECT 1", false},
	{"ALTER DEFINER = CURRENT_USER EVENT e ENABLE", false},
	{"ALTER ALGORITHM = INPLACE TABLE t ADD c INT", false},
}

func TestViewVisibilityFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range viewVisibilityFormCases {
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

// TestViewVisibilityForms holds the parser to the recorded verdicts without the
// server, so a regression fails however the suite is run.
func TestViewVisibilityForms(t *testing.T) {
	for _, c := range viewVisibilityFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
