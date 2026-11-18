package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestOptimizerTrace covers optimizerTrace (StarRocks.g4:2584): every mode,
// the optional trace identifier, and the query it wraps.
func TestOptimizerTrace(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		mode       string
		identifier string
	}{
		{"TRACE ALL SELECT 1", "ALL", ""},
		{"TRACE ALL SELECT 1 UNION SELECT 2 INTO OUTFILE 'x'", "ALL", ""},
		{"TRACE LOGS SELECT 1", "LOGS", ""},
		{"TRACE TIMES SELECT 1", "TIMES", ""},
		{"TRACE VALUES SELECT 1", "VALUES", ""},
		{"TRACE REASON SELECT 1", "REASON", ""},
		{"TRACE ALL trace_name SELECT 1 FROM t", "ALL", "trace_name"},
		{"TRACE ALL WITH c AS (SELECT 1) SELECT * FROM c", "ALL", ""},
		{"TRACE ALL WITH c AS (SELECT 1) (SELECT * FROM c)", "ALL", ""},
		{"TRACE ALL WITH c AS (SELECT 1) (WITH d AS (SELECT 2) SELECT * FROM c, d)", "ALL", ""},
		{"TRACE ALL (SELECT 1) INTO OUTFILE 'x'", "ALL", ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := parseForTest(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt, ok := file.Stmts[0].(*ast.OptimizerTraceStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.OptimizerTraceStmt", file.Stmts[0])
			}
			if stmt.Mode != tc.mode || stmt.Identifier != tc.identifier {
				t.Errorf("(Mode, Identifier) = (%q, %q), want (%q, %q)", stmt.Mode, stmt.Identifier, tc.mode, tc.identifier)
			}
			if stmt.Query == nil {
				t.Error("Query is nil")
			}
		})
	}
}

func TestOptimizerTraceInQueryStatementPositions(t *testing.T) {
	for _, sql := range []string{
		"INSERT INTO dst TRACE ALL SELECT * FROM src",
		"INSERT INTO dst EXPLAIN SELECT * FROM src",
		"INSERT INTO dst DESC ANALYZE SELECT * FROM src",
		"CREATE TABLE dst AS TRACE ALL SELECT 1",
		"CREATE TABLE dst AS EXPLAIN SELECT 1",
		"CREATE TABLE dst AS DESC SELECT 1",
	} {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("%q: got %d statements, want 1", sql, len(file.Stmts))
		}
	}
}

// An expression subquery holds queryRelation only; the trace prefix belongs to
// queryStatement, so it cannot stand in for EXISTS or a scalar subquery.
func TestOptimizerTraceRejectsExpressionSubquery(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t WHERE EXISTS (TRACE LOGS SELECT * FROM secret)",
		"SELECT (TRACE ALL SELECT 1)",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

// queryStatement's explain prefix precedes a queryRelation, so explained DML
// is not a valid INSERT source or CTAS body.
func TestExplainedDMLRejectedInQueryStatement(t *testing.T) {
	for _, sql := range []string{
		"INSERT INTO dst EXPLAIN INSERT INTO inner_dst SELECT 1",
		"CREATE TABLE dst AS EXPLAIN INSERT INTO inner_dst SELECT 1",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

// The OUTFILE clause belongs to the wrapper, so no child Loc runs past its
// parent.
func TestTraceOutfileLoc(t *testing.T) {
	file, errs := parseForTest("TRACE ALL (SELECT 1) INTO OUTFILE 'x'")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	trace := file.Stmts[0].(*ast.OptimizerTraceStmt)
	paren, ok := trace.Query.(*ast.ParenSelect)
	if !ok {
		t.Fatalf("Query = %T, want *ast.ParenSelect", trace.Query)
	}
	if paren.Into == nil {
		t.Fatal("ParenSelect.Into is nil")
	}
	if end := ast.NodeLoc(paren.Sel).End; end > paren.Loc.End {
		t.Errorf("child End = %d, parent End = %d", end, paren.Loc.End)
	}
}

// Hoisting a right-arm OUTFILE onto the set operation shrinks the arm.
func TestTraceSetOpOutfileLoc(t *testing.T) {
	file, errs := parseForTest("TRACE ALL SELECT 1 UNION SELECT 2 INTO OUTFILE 'x'")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	trace := file.Stmts[0].(*ast.OptimizerTraceStmt)
	setOp, ok := trace.Query.(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SetOpStmt", trace.Query)
	}
	if setOp.Into == nil {
		t.Fatal("SetOp.Into is nil")
	}
	right := setOp.Right.(*ast.SelectStmt)
	if right.Loc.End > setOp.Into.Loc.Start {
		t.Errorf("right arm End = %d, OUTFILE starts = %d", right.Loc.End, setOp.Into.Loc.Start)
	}
	if right.Into != nil {
		t.Error("right arm still carries Into")
	}
}

func TestOptimizerTraceRejects(t *testing.T) {
	for _, sql := range []string{
		"TRACE SELECT 1",
		"TRACE ALL zzz",
		"TRACE ALL a zzz SELECT 1",
		"TRACE ALL INSERT INTO t SELECT 1",
		"EXPLAIN TRACE ALL SELECT 1",
		"TRACE ALL SELECT 1 INTO OUTFILE 'a' INTO OUTFILE 'b'",
		"TRACE ALL (SELECT 1 INTO OUTFILE 'x')",
		"TRACE ALL SELECT 1 INTO OUTFILE 'a' UNION SELECT 2",
		"TRACE ALL SELECT 1 UNION (SELECT 2 INTO OUTFILE 'x')",
		"DESC (SELECT 1 INTO OUTFILE 'x')",
		"DESC SELECT 1 INTO OUTFILE 'x' UNION SELECT 2",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

func TestExplainedParenthesizedWithQuery(t *testing.T) {
	for _, sql := range []string{
		"INSERT INTO dst EXPLAIN WITH c AS (SELECT 1) (SELECT * FROM c)",
		"CREATE TABLE dst AS EXPLAIN WITH c AS (SELECT 1) (SELECT * FROM c)",
		"DESC (SELECT 1) INTO OUTFILE 'x'",
		"DESCRIBE (SELECT 1) INTO OUTFILE 'x'",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", sql, errs)
		}
	}
}
