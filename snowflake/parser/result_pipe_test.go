package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// The Snowflake result-pipe operator `->>` feeds a statement's result set into
// a following query, where the prior result is named $1. A SHOW source carries
// the pipe on ShowStmt.Pipe (see show.go / show_test.go); any other source is
// wrapped in a top-level *ast.ResultScanStmt by parseTopStmt. Both rely on the
// $N table reference (FROM $1) added in gap-from-values.

// TestResultPipe_ShowCorpusShapes covers the exact official-corpus SHOW pipes
// closed by gap-from-values (show-tables/example_07..09,
// create-external-volume/example_05, create-warehouse/example_07). They flow
// through ShowStmt.Pipe; the blocker was FROM $1, now parseable.
func TestResultPipe_ShowCorpusShapes(t *testing.T) {
	cases := []string{
		`SHOW TERSE TABLES ->> SELECT * FROM $1 ORDER BY "created_on" DESC`,
		`SHOW TABLES ->> SELECT "name", "database_name", "schema_name", "bytes" FROM $1 ORDER BY "bytes" DESC`,
		`SHOW TABLES IN ACCOUNT ->> SELECT "database_name" || '.' || "schema_name" || '.' || "name" AS fully_qualified_name FROM $1 ORDER BY fully_qualified_name`,
		`SHOW ICEBERG TABLES ->> SELECT * FROM $1 WHERE "external_volume_name" = 'my_external_volume_1'`,
		`SHOW WAREHOUSES LIKE 'test_gen_warehouse' ->> SELECT "name", "resource_constraint" FROM $1`,
	}
	for _, c := range cases {
		node, errs := parseSingle(c, 0)
		if len(errs) > 0 {
			t.Errorf("%q: %v", c, errs)
			continue
		}
		show, ok := node.(*ast.ShowStmt)
		if !ok {
			t.Errorf("%q: got %T, want *ast.ShowStmt", c, node)
			continue
		}
		if show.Pipe == nil {
			t.Errorf("%q: ShowStmt.Pipe is nil", c)
		}
	}
}

// TestResultPipe_GeneralStmt verifies a non-SHOW source is wrapped in a
// *ast.ResultScanStmt with Source and Query set, and a valid span Loc.
func TestResultPipe_GeneralStmt(t *testing.T) {
	node, errs := parseSingle(`SELECT 1 AS x ->> SELECT * FROM $1`, 0)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	rs, ok := node.(*ast.ResultScanStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.ResultScanStmt", node)
	}
	if _, ok := rs.Source.(*ast.SelectStmt); !ok {
		t.Errorf("Source = %T, want *ast.SelectStmt", rs.Source)
	}
	sel, ok := rs.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", rs.Query)
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok || ref.DollarN == nil {
		t.Errorf("Query.From[0] = %T, want *ast.TableRef with DollarN", sel.From[0])
	}
	if !rs.Loc.IsValid() {
		t.Errorf("ResultScanStmt.Loc invalid: %+v", rs.Loc)
	}
}

// TestResultPipe_Chained verifies a ->> b ->> c nests left-associatively into
// two ResultScanStmt levels (outer.Source is the inner ResultScanStmt).
func TestResultPipe_Chained(t *testing.T) {
	node, errs := parseSingle(`SELECT 1 ->> SELECT * FROM $1 ->> SELECT * FROM $1`, 0)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	outer, ok := node.(*ast.ResultScanStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.ResultScanStmt", node)
	}
	if _, ok := outer.Source.(*ast.ResultScanStmt); !ok {
		t.Errorf("outer.Source = %T, want nested *ast.ResultScanStmt", outer.Source)
	}
}

// Negative: `->>` with no following query is an error (parseQueryExpr fails).
func TestResultPipe_NoFollowingQuery(t *testing.T) {
	result := ParseBestEffort(`SELECT 1 ->>`)
	if len(result.Errors) == 0 {
		t.Error("expected an error for trailing `->>` with no query, got none")
	}
}

// Regression: a plain SHOW with no `->>` is unchanged (Pipe stays nil).
func TestResultPipe_ShowWithoutPipeUnchanged(t *testing.T) {
	node, errs := parseSingle(`SHOW TABLES`, 0)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	show, ok := node.(*ast.ShowStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.ShowStmt", node)
	}
	if show.Pipe != nil {
		t.Errorf("ShowStmt.Pipe = %v, want nil for SHOW without ->>", show.Pipe)
	}
}

// Regression: the single-arrow operator `->` (tokArrow, e.g. lambda) is not
// confused with `->>` (tokFlow) by max-munch. A lambda x -> x + 1 still parses
// as an expression, not a result pipe.
func TestResultPipe_SingleArrowUnaffected(t *testing.T) {
	// TRANSFORM-style higher-order function with a lambda using ->.
	sel, errs := testParseSelectStmt(`SELECT FILTER(a, x -> x > 0) FROM t`)
	if len(errs) > 0 {
		t.Fatalf("single-arrow lambda should parse, got: %v", errs)
	}
	if _, ok := sel.Targets[0].Expr.(*ast.FuncCallExpr); !ok {
		t.Errorf("target[0] = %T, want *ast.FuncCallExpr", sel.Targets[0].Expr)
	}
}
