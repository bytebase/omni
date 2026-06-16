package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

func mustParseCreateMTMV(t *testing.T, input string) *ast.CreateMTMVStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateMTMVStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.CreateMTMVStmt", input, file.Stmts[0])
	}
	return stmt
}

// StarRocks async-MV scheduling (PR3): REFRESH ASYNC [START('ts')] EVERY
// (INTERVAL n unit) — distinct from the doris ON SCHEDULE EVERY <bare> form.

func TestMVRefreshAsyncEvery(t *testing.T) {
	stmt := mustParseCreateMTMV(t, "CREATE MATERIALIZED VIEW mv REFRESH ASYNC EVERY (INTERVAL 1 DAY) AS SELECT a FROM t")
	if stmt.RefreshTrigger == nil || !stmt.RefreshTrigger.OnSchedule {
		t.Fatalf("RefreshTrigger = %+v, want OnSchedule", stmt.RefreshTrigger)
	}
	if stmt.RefreshTrigger.Interval != "1 DAY" {
		t.Errorf("Interval = %q, want \"1 DAY\"", stmt.RefreshTrigger.Interval)
	}
}

func TestMVRefreshAsyncStartEvery(t *testing.T) {
	stmt := mustParseCreateMTMV(t, "CREATE MATERIALIZED VIEW mv REFRESH ASYNC START('2024-12-01 20:30:00') EVERY (INTERVAL 1 DAY) AS SELECT a FROM t")
	if stmt.RefreshTrigger == nil || stmt.RefreshTrigger.StartsAt == "" {
		t.Fatalf("RefreshTrigger = %+v, want StartsAt set", stmt.RefreshTrigger)
	}
	if stmt.RefreshTrigger.Interval != "1 DAY" {
		t.Errorf("Interval = %q, want \"1 DAY\"", stmt.RefreshTrigger.Interval)
	}
}

func TestMVRefreshImmediateAsyncEvery(t *testing.T) {
	_ = mustParseCreateMTMV(t, "CREATE MATERIALIZED VIEW mv REFRESH IMMEDIATE ASYNC EVERY (INTERVAL 5 MINUTE) AS SELECT a FROM t")
}

// Regression: plain REFRESH ASYNC (no scheduling tail) still parses.
func TestMVRefreshAsyncPlain(t *testing.T) {
	_ = mustParseCreateMTMV(t, "CREATE MATERIALIZED VIEW mv REFRESH ASYNC AS SELECT a FROM t")
}

// The schedule tail is ASYNC-only: MANUAL/INCREMENTAL + EVERY must reject.
func TestMVRefreshManualEveryRejected(t *testing.T) {
	_, errs := Parse("CREATE MATERIALIZED VIEW mv REFRESH MANUAL EVERY (INTERVAL 1 DAY) AS SELECT a FROM t")
	if len(errs) == 0 {
		t.Fatal("expected a parse error for MANUAL + EVERY, got none")
	}
}

// EVERY is mandatory once the async tail starts: START without EVERY must reject.
func TestMVRefreshAsyncStartNoEveryRejected(t *testing.T) {
	_, errs := Parse("CREATE MATERIALIZED VIEW mv REFRESH ASYNC START('2024-12-01 20:30:00') AS SELECT a FROM t")
	if len(errs) == 0 {
		t.Fatal("expected a parse error for START without EVERY, got none")
	}
}

// The no-parens form must REJECT (StarRocks requires EVERY '(' INTERVAL ... ')').
// This was a silent false-accept before the dedicated handler.
func TestMVRefreshAsyncEveryNoParensRejected(t *testing.T) {
	_, errs := Parse("CREATE MATERIALIZED VIEW mv REFRESH ASYNC EVERY INTERVAL 1 DAY AS SELECT a FROM t")
	if len(errs) == 0 {
		t.Fatal("expected a parse error for EVERY without parentheses, got none")
	}
}
