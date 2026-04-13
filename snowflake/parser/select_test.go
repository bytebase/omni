package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseSelectStmt is a helper that parses input via ParseBestEffort and
// returns the first statement as *ast.SelectStmt plus any errors.
func testParseSelectStmt(input string) (*ast.SelectStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	sel, ok := result.File.Stmts[0].(*ast.SelectStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a SelectStmt"})
	}
	return sel, result.Errors
}

// ---------------------------------------------------------------------------
// 1. SELECT 1 — simplest
// ---------------------------------------------------------------------------

func TestSelect_Literal(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel == nil {
		t.Fatal("expected SelectStmt, got nil")
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if sel.Targets[0].Star {
		t.Error("target should not be star")
	}
	lit, ok := sel.Targets[0].Expr.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", sel.Targets[0].Expr)
	}
	if lit.Kind != ast.LitInt || lit.Ival != 1 {
		t.Errorf("literal = %v/%d, want LitInt/1", lit.Kind, lit.Ival)
	}
	if len(sel.From) != 0 {
		t.Errorf("from = %d, want 0", len(sel.From))
	}
	if sel.Where != nil {
		t.Error("where should be nil")
	}
}

// ---------------------------------------------------------------------------
// 2. SELECT 1, 2, 3 — multiple targets
// ---------------------------------------------------------------------------

func TestSelect_MultipleTargets(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT 1, 2, 3")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 3 {
		t.Fatalf("targets = %d, want 3", len(sel.Targets))
	}
	for i, want := range []int64{1, 2, 3} {
		lit, ok := sel.Targets[i].Expr.(*ast.Literal)
		if !ok {
			t.Fatalf("target[%d]: expected *ast.Literal, got %T", i, sel.Targets[i].Expr)
		}
		if lit.Ival != want {
			t.Errorf("target[%d]: Ival = %d, want %d", i, lit.Ival, want)
		}
	}
}

// ---------------------------------------------------------------------------
// 3. SELECT a AS x, b AS y FROM t — aliases with AS
// ---------------------------------------------------------------------------

func TestSelect_AliasWithAS(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a AS x, b AS y FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
	if sel.Targets[0].Alias.Name != "x" {
		t.Errorf("target[0] alias = %q, want %q", sel.Targets[0].Alias.Name, "x")
	}
	if sel.Targets[1].Alias.Name != "y" {
		t.Errorf("target[1] alias = %q, want %q", sel.Targets[1].Alias.Name, "y")
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

// ---------------------------------------------------------------------------
// 4. SELECT a x FROM t — alias without AS
// ---------------------------------------------------------------------------

func TestSelect_AliasWithoutAS(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a x FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if sel.Targets[0].Alias.Name != "x" {
		t.Errorf("alias = %q, want %q", sel.Targets[0].Alias.Name, "x")
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

// ---------------------------------------------------------------------------
// 5. SELECT * FROM t — star target
// ---------------------------------------------------------------------------

func TestSelect_Star(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if !sel.Targets[0].Star {
		t.Error("target should be star")
	}
	star, ok := sel.Targets[0].Expr.(*ast.StarExpr)
	if !ok {
		t.Fatalf("expected *ast.StarExpr, got %T", sel.Targets[0].Expr)
	}
	if star.Qualifier != nil {
		t.Error("star qualifier should be nil for bare *")
	}
}

// ---------------------------------------------------------------------------
// 6. SELECT * EXCLUDE (a, b) FROM t — EXCLUDE
// ---------------------------------------------------------------------------

func TestSelect_StarExclude(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * EXCLUDE (a, b) FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	target := sel.Targets[0]
	if !target.Star {
		t.Error("target should be star")
	}
	if len(target.Exclude) != 2 {
		t.Fatalf("exclude = %d, want 2", len(target.Exclude))
	}
	if target.Exclude[0].Name != "a" {
		t.Errorf("exclude[0] = %q, want %q", target.Exclude[0].Name, "a")
	}
	if target.Exclude[1].Name != "b" {
		t.Errorf("exclude[1] = %q, want %q", target.Exclude[1].Name, "b")
	}
}

// ---------------------------------------------------------------------------
// 7. SELECT t.* FROM t — qualified star
// ---------------------------------------------------------------------------

func TestSelect_QualifiedStar(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT t.* FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	target := sel.Targets[0]
	if !target.Star {
		t.Error("target should be star")
	}
	star, ok := target.Expr.(*ast.StarExpr)
	if !ok {
		t.Fatalf("expected *ast.StarExpr, got %T", target.Expr)
	}
	if star.Qualifier == nil {
		t.Fatal("star qualifier should not be nil")
	}
	if star.Qualifier.Name.Name != "t" {
		t.Errorf("qualifier = %q, want %q", star.Qualifier.Name.Name, "t")
	}
}

// ---------------------------------------------------------------------------
// 8. SELECT DISTINCT a FROM t
// ---------------------------------------------------------------------------

func TestSelect_Distinct(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT DISTINCT a FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Distinct {
		t.Error("expected Distinct = true")
	}
	if sel.All {
		t.Error("expected All = false")
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
}

// ---------------------------------------------------------------------------
// 9. SELECT TOP 10 a FROM t
// ---------------------------------------------------------------------------

func TestSelect_Top(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT TOP 10 a FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Top == nil {
		t.Fatal("expected Top to be set")
	}
	lit, ok := sel.Top.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Top, got %T", sel.Top)
	}
	if lit.Ival != 10 {
		t.Errorf("Top = %d, want 10", lit.Ival)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
}

// ---------------------------------------------------------------------------
// 10. SELECT a FROM t1, t2 AS x — comma FROM with alias
// ---------------------------------------------------------------------------

func TestSelect_CommaFrom(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t1, t2 AS x")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 2 {
		t.Fatalf("from = %d, want 2", len(sel.From))
	}
	if sel.From[0].Name.Name.Name != "t1" {
		t.Errorf("from[0] = %q, want %q", sel.From[0].Name.Name.Name, "t1")
	}
	if !sel.From[0].Alias.IsEmpty() {
		t.Errorf("from[0] alias should be empty, got %q", sel.From[0].Alias.Name)
	}
	if sel.From[1].Name.Name.Name != "t2" {
		t.Errorf("from[1] = %q, want %q", sel.From[1].Name.Name.Name, "t2")
	}
	if sel.From[1].Alias.Name != "x" {
		t.Errorf("from[1] alias = %q, want %q", sel.From[1].Alias.Name, "x")
	}
}

// ---------------------------------------------------------------------------
// 11. SELECT a FROM t WHERE a > 0
// ---------------------------------------------------------------------------

func TestSelect_Where(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t WHERE a > 0")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Where == nil {
		t.Fatal("expected Where to be set")
	}
	binExpr, ok := sel.Where.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr for Where, got %T", sel.Where)
	}
	if binExpr.Op != ast.BinGt {
		t.Errorf("where op = %v, want BinGt", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// 12. SELECT a, COUNT(*) FROM t GROUP BY a — normal GROUP BY
// ---------------------------------------------------------------------------

func TestSelect_GroupByNormal(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a, COUNT(*) FROM t GROUP BY a")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.GroupBy == nil {
		t.Fatal("expected GroupBy to be set")
	}
	if sel.GroupBy.Kind != ast.GroupByNormal {
		t.Errorf("GroupBy.Kind = %v, want GroupByNormal", sel.GroupBy.Kind)
	}
	if len(sel.GroupBy.Items) != 1 {
		t.Fatalf("GroupBy.Items = %d, want 1", len(sel.GroupBy.Items))
	}
}

// ---------------------------------------------------------------------------
// 13. GROUP BY CUBE / ROLLUP / GROUPING SETS / ALL variants
// ---------------------------------------------------------------------------

func TestSelect_GroupByCube(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t GROUP BY CUBE (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.GroupBy == nil {
		t.Fatal("expected GroupBy to be set")
	}
	if sel.GroupBy.Kind != ast.GroupByCube {
		t.Errorf("GroupBy.Kind = %v, want GroupByCube", sel.GroupBy.Kind)
	}
	if len(sel.GroupBy.Items) != 2 {
		t.Fatalf("GroupBy.Items = %d, want 2", len(sel.GroupBy.Items))
	}
}

func TestSelect_GroupByRollup(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t GROUP BY ROLLUP (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.GroupBy == nil {
		t.Fatal("expected GroupBy to be set")
	}
	if sel.GroupBy.Kind != ast.GroupByRollup {
		t.Errorf("GroupBy.Kind = %v, want GroupByRollup", sel.GroupBy.Kind)
	}
	if len(sel.GroupBy.Items) != 2 {
		t.Fatalf("GroupBy.Items = %d, want 2", len(sel.GroupBy.Items))
	}
}

func TestSelect_GroupByGroupingSets(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t GROUP BY GROUPING SETS (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.GroupBy == nil {
		t.Fatal("expected GroupBy to be set")
	}
	if sel.GroupBy.Kind != ast.GroupByGroupingSets {
		t.Errorf("GroupBy.Kind = %v, want GroupByGroupingSets", sel.GroupBy.Kind)
	}
	if len(sel.GroupBy.Items) != 2 {
		t.Fatalf("GroupBy.Items = %d, want 2", len(sel.GroupBy.Items))
	}
}

func TestSelect_GroupByAll(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t GROUP BY ALL")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.GroupBy == nil {
		t.Fatal("expected GroupBy to be set")
	}
	if sel.GroupBy.Kind != ast.GroupByAll {
		t.Errorf("GroupBy.Kind = %v, want GroupByAll", sel.GroupBy.Kind)
	}
	if len(sel.GroupBy.Items) != 0 {
		t.Errorf("GroupBy.Items = %d, want 0 for ALL", len(sel.GroupBy.Items))
	}
}

// ---------------------------------------------------------------------------
// 14. SELECT ... HAVING COUNT(*) > 1
// ---------------------------------------------------------------------------

func TestSelect_Having(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a, COUNT(*) FROM t GROUP BY a HAVING COUNT(*) > 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Having == nil {
		t.Fatal("expected Having to be set")
	}
	binExpr, ok := sel.Having.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr for Having, got %T", sel.Having)
	}
	if binExpr.Op != ast.BinGt {
		t.Errorf("Having op = %v, want BinGt", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// 15. SELECT ... QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1
// ---------------------------------------------------------------------------

func TestSelect_Qualify(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Qualify == nil {
		t.Fatal("expected Qualify to be set")
	}
	// QUALIFY ROW_NUMBER() OVER (...) = 1 is a BinaryExpr (= comparison)
	binExpr, ok := sel.Qualify.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr for Qualify, got %T", sel.Qualify)
	}
	if binExpr.Op != ast.BinEq {
		t.Errorf("Qualify op = %v, want BinEq", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// 16. SELECT ... ORDER BY a DESC NULLS LAST
// ---------------------------------------------------------------------------

func TestSelect_OrderBy(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t ORDER BY a DESC NULLS LAST")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(sel.OrderBy))
	}
	item := sel.OrderBy[0]
	if !item.Desc {
		t.Error("expected Desc = true")
	}
	if item.NullsFirst == nil {
		t.Fatal("expected NullsFirst to be set")
	}
	if *item.NullsFirst {
		t.Error("expected NullsFirst = false (NULLS LAST)")
	}
}

// ---------------------------------------------------------------------------
// 17. SELECT ... LIMIT 10 OFFSET 5
// ---------------------------------------------------------------------------

func TestSelect_LimitOffset(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t LIMIT 10 OFFSET 5")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Limit == nil {
		t.Fatal("expected Limit to be set")
	}
	litLimit, ok := sel.Limit.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Limit, got %T", sel.Limit)
	}
	if litLimit.Ival != 10 {
		t.Errorf("Limit = %d, want 10", litLimit.Ival)
	}
	if sel.Offset == nil {
		t.Fatal("expected Offset to be set")
	}
	litOffset, ok := sel.Offset.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Offset, got %T", sel.Offset)
	}
	if litOffset.Ival != 5 {
		t.Errorf("Offset = %d, want 5", litOffset.Ival)
	}
}

// ---------------------------------------------------------------------------
// 18. SELECT ... FETCH FIRST 10 ROWS ONLY
// ---------------------------------------------------------------------------

func TestSelect_Fetch(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t FETCH FIRST 10 ROWS ONLY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Fetch == nil {
		t.Fatal("expected Fetch to be set")
	}
	litCount, ok := sel.Fetch.Count.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Fetch.Count, got %T", sel.Fetch.Count)
	}
	if litCount.Ival != 10 {
		t.Errorf("Fetch.Count = %d, want 10", litCount.Ival)
	}
}

func TestSelect_FetchNext(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t FETCH NEXT 5 ROW ONLY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Fetch == nil {
		t.Fatal("expected Fetch to be set")
	}
	litCount, ok := sel.Fetch.Count.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Fetch.Count, got %T", sel.Fetch.Count)
	}
	if litCount.Ival != 5 {
		t.Errorf("Fetch.Count = %d, want 5", litCount.Ival)
	}
}

func TestSelect_OffsetFetch(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t OFFSET 3 FETCH FIRST 10 ROWS ONLY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Offset == nil {
		t.Fatal("expected Offset to be set")
	}
	litOff, ok := sel.Offset.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Offset, got %T", sel.Offset)
	}
	if litOff.Ival != 3 {
		t.Errorf("Offset = %d, want 3", litOff.Ival)
	}
	if sel.Fetch == nil {
		t.Fatal("expected Fetch to be set")
	}
}

// ---------------------------------------------------------------------------
// 19. WITH cte AS (SELECT 1 AS x) SELECT * FROM cte — basic CTE
// ---------------------------------------------------------------------------

func TestSelect_CTE(t *testing.T) {
	sel, errs := testParseSelectStmt("WITH cte AS (SELECT 1 AS x) SELECT * FROM cte")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel == nil {
		t.Fatal("expected SelectStmt, got nil")
	}
	if len(sel.With) != 1 {
		t.Fatalf("With = %d, want 1", len(sel.With))
	}
	cte := sel.With[0]
	if cte.Name.Name != "cte" {
		t.Errorf("CTE name = %q, want %q", cte.Name.Name, "cte")
	}
	if cte.Recursive {
		t.Error("expected Recursive = false")
	}
	if cte.Query == nil {
		t.Fatal("CTE query is nil")
	}
	innerSel, ok := cte.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("CTE query: expected *ast.SelectStmt, got %T", cte.Query)
	}
	if len(innerSel.Targets) != 1 {
		t.Errorf("CTE query targets = %d, want 1", len(innerSel.Targets))
	}
	// Check outer SELECT
	if len(sel.Targets) != 1 {
		t.Fatalf("outer targets = %d, want 1", len(sel.Targets))
	}
	if !sel.Targets[0].Star {
		t.Error("outer target should be star")
	}
}

func TestSelect_CTERecursive(t *testing.T) {
	input := "WITH RECURSIVE r (n) AS (SELECT 1) SELECT * FROM r"
	sel, errs := testParseSelectStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.With) != 1 {
		t.Fatalf("With = %d, want 1", len(sel.With))
	}
	cte := sel.With[0]
	if !cte.Recursive {
		t.Error("expected Recursive = true")
	}
	if cte.Name.Name != "r" {
		t.Errorf("CTE name = %q, want %q", cte.Name.Name, "r")
	}
	if len(cte.Columns) != 1 {
		t.Fatalf("CTE columns = %d, want 1", len(cte.Columns))
	}
	if cte.Columns[0].Name != "n" {
		t.Errorf("CTE column[0] = %q, want %q", cte.Columns[0].Name, "n")
	}
}

func TestSelect_MultipleCTEs(t *testing.T) {
	input := "WITH a AS (SELECT 1), b AS (SELECT 2) SELECT * FROM a, b"
	sel, errs := testParseSelectStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.With) != 2 {
		t.Fatalf("With = %d, want 2", len(sel.With))
	}
	if sel.With[0].Name.Name != "a" {
		t.Errorf("CTE[0] name = %q, want %q", sel.With[0].Name.Name, "a")
	}
	if sel.With[1].Name.Name != "b" {
		t.Errorf("CTE[1] name = %q, want %q", sel.With[1].Name.Name, "b")
	}
}

// ---------------------------------------------------------------------------
// 20. SELECT (SELECT 1) — SubqueryExpr
// ---------------------------------------------------------------------------

func TestSelect_SubqueryExpr(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT (SELECT 1)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	subq, ok := sel.Targets[0].Expr.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr, got %T", sel.Targets[0].Expr)
	}
	innerSel, ok := subq.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt inside SubqueryExpr, got %T", subq.Query)
	}
	if len(innerSel.Targets) != 1 {
		t.Errorf("inner targets = %d, want 1", len(innerSel.Targets))
	}
}

// ---------------------------------------------------------------------------
// 21. SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t2) — ExistsExpr
// ---------------------------------------------------------------------------

func TestSelect_ExistsExpr(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Where == nil {
		t.Fatal("expected Where to be set")
	}
	exists, ok := sel.Where.(*ast.ExistsExpr)
	if !ok {
		t.Fatalf("expected *ast.ExistsExpr for Where, got %T", sel.Where)
	}
	innerSel, ok := exists.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt inside ExistsExpr, got %T", exists.Query)
	}
	if len(innerSel.From) != 1 {
		t.Errorf("inner from = %d, want 1", len(innerSel.From))
	}
}

// ---------------------------------------------------------------------------
// 22. SELECT * FROM t WHERE a IN (SELECT b FROM t2) — IN-subquery
// ---------------------------------------------------------------------------

func TestSelect_InSubquery(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t WHERE a IN (SELECT b FROM t2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Where == nil {
		t.Fatal("expected Where to be set")
	}
	inExpr, ok := sel.Where.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr for Where, got %T", sel.Where)
	}
	if inExpr.Not {
		t.Error("expected Not = false")
	}
	if len(inExpr.Values) != 1 {
		t.Fatalf("InExpr.Values = %d, want 1", len(inExpr.Values))
	}
	subq, ok := inExpr.Values[0].(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr in InExpr.Values, got %T", inExpr.Values[0])
	}
	innerSel, ok := subq.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt inside SubqueryExpr, got %T", subq.Query)
	}
	if len(innerSel.From) != 1 {
		t.Errorf("inner from = %d, want 1", len(innerSel.From))
	}
}

func TestSelect_NotInSubquery(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t WHERE a NOT IN (SELECT b FROM t2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	inExpr, ok := sel.Where.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr for Where, got %T", sel.Where)
	}
	if !inExpr.Not {
		t.Error("expected Not = true")
	}
}

// ---------------------------------------------------------------------------
// Additional edge cases
// ---------------------------------------------------------------------------

func TestSelect_ImplicitAliasNotConsumeClauseKeyword(t *testing.T) {
	// Verify that FROM is not consumed as an alias.
	sel, errs := testParseSelectStmt("SELECT a FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if !sel.Targets[0].Alias.IsEmpty() {
		t.Errorf("target alias should be empty, got %q", sel.Targets[0].Alias.Name)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

func TestSelect_ImplicitAliasNotConsumeWHERE(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t WHERE a = 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Targets[0].Alias.IsEmpty() {
		t.Errorf("target alias should be empty, got %q", sel.Targets[0].Alias.Name)
	}
	if sel.Where == nil {
		t.Fatal("WHERE should be set")
	}
}

func TestSelect_ImplicitAliasNotConsumeGROUP(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t GROUP BY a")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Targets[0].Alias.IsEmpty() {
		t.Errorf("target alias should be empty, got %q", sel.Targets[0].Alias.Name)
	}
	if sel.GroupBy == nil {
		t.Fatal("GROUP BY should be set")
	}
}

func TestSelect_ImplicitAliasNotConsumeORDER(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t ORDER BY a")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Targets[0].Alias.IsEmpty() {
		t.Errorf("target alias should be empty, got %q", sel.Targets[0].Alias.Name)
	}
	if len(sel.OrderBy) != 1 {
		t.Fatalf("ORDER BY = %d, want 1", len(sel.OrderBy))
	}
}

func TestSelect_ImplicitAliasNotConsumeLIMIT(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t LIMIT 10")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Targets[0].Alias.IsEmpty() {
		t.Errorf("target alias should be empty, got %q", sel.Targets[0].Alias.Name)
	}
	if sel.Limit == nil {
		t.Fatal("LIMIT should be set")
	}
}

func TestSelect_TableAliasNotConsumeWHERE(t *testing.T) {
	// Verify that WHERE is not consumed as a table alias.
	sel, errs := testParseSelectStmt("SELECT a FROM t WHERE a > 0")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	if !sel.From[0].Alias.IsEmpty() {
		t.Errorf("table alias should be empty, got %q", sel.From[0].Alias.Name)
	}
	if sel.Where == nil {
		t.Fatal("WHERE should be set")
	}
}

func TestSelect_TableImplicitAlias(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t1 x, t2 y WHERE x.a = y.a")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 2 {
		t.Fatalf("from = %d, want 2", len(sel.From))
	}
	if sel.From[0].Alias.Name != "x" {
		t.Errorf("from[0] alias = %q, want %q", sel.From[0].Alias.Name, "x")
	}
	if sel.From[1].Alias.Name != "y" {
		t.Errorf("from[1] alias = %q, want %q", sel.From[1].Alias.Name, "y")
	}
}

func TestSelect_AllClauses(t *testing.T) {
	input := `SELECT DISTINCT a, COUNT(*) AS cnt
FROM t
WHERE a > 0
GROUP BY a
HAVING COUNT(*) > 1
QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1
ORDER BY a DESC
LIMIT 10 OFFSET 5`
	sel, errs := testParseSelectStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.Distinct {
		t.Error("expected Distinct = true")
	}
	if len(sel.Targets) != 2 {
		t.Errorf("targets = %d, want 2", len(sel.Targets))
	}
	if len(sel.From) != 1 {
		t.Errorf("from = %d, want 1", len(sel.From))
	}
	if sel.Where == nil {
		t.Error("where should be set")
	}
	if sel.GroupBy == nil {
		t.Error("groupby should be set")
	}
	if sel.Having == nil {
		t.Error("having should be set")
	}
	if sel.Qualify == nil {
		t.Error("qualify should be set")
	}
	if len(sel.OrderBy) != 1 {
		t.Errorf("orderby = %d, want 1", len(sel.OrderBy))
	}
	if sel.Limit == nil {
		t.Error("limit should be set")
	}
	if sel.Offset == nil {
		t.Error("offset should be set")
	}
}

func TestSelect_QualifiedTableName(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM mydb.myschema.t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	name := sel.From[0].Name
	if name.Database.Name != "mydb" {
		t.Errorf("database = %q, want %q", name.Database.Name, "mydb")
	}
	if name.Schema.Name != "myschema" {
		t.Errorf("schema = %q, want %q", name.Schema.Name, "myschema")
	}
	if name.Name.Name != "t" {
		t.Errorf("name = %q, want %q", name.Name.Name, "t")
	}
}

func TestSelect_SelectAll(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT ALL a FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !sel.All {
		t.Error("expected All = true")
	}
	if sel.Distinct {
		t.Error("expected Distinct = false")
	}
}

func TestSelect_MultiOrderBy(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a, b FROM t ORDER BY a ASC, b DESC NULLS FIRST")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.OrderBy) != 2 {
		t.Fatalf("orderby = %d, want 2", len(sel.OrderBy))
	}
	if sel.OrderBy[0].Desc {
		t.Error("orderby[0] should be ASC")
	}
	if !sel.OrderBy[1].Desc {
		t.Error("orderby[1] should be DESC")
	}
	if sel.OrderBy[1].NullsFirst == nil || !*sel.OrderBy[1].NullsFirst {
		t.Error("orderby[1] should be NULLS FIRST")
	}
}

func TestSelect_LimitOnly(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t LIMIT 10")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Limit == nil {
		t.Fatal("LIMIT should be set")
	}
	if sel.Offset != nil {
		t.Error("OFFSET should be nil")
	}
}

func TestSelect_OffsetOnly(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t OFFSET 5")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Offset == nil {
		t.Fatal("OFFSET should be set")
	}
	if sel.Limit != nil {
		t.Error("LIMIT should be nil")
	}
	if sel.Fetch != nil {
		t.Error("FETCH should be nil")
	}
}

func TestSelect_CTEWithColumns(t *testing.T) {
	input := "WITH cte (x, y) AS (SELECT 1, 2) SELECT * FROM cte"
	sel, errs := testParseSelectStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.With) != 1 {
		t.Fatalf("With = %d, want 1", len(sel.With))
	}
	cte := sel.With[0]
	if len(cte.Columns) != 2 {
		t.Fatalf("CTE columns = %d, want 2", len(cte.Columns))
	}
	if cte.Columns[0].Name != "x" {
		t.Errorf("CTE column[0] = %q, want %q", cte.Columns[0].Name, "x")
	}
	if cte.Columns[1].Name != "y" {
		t.Errorf("CTE column[1] = %q, want %q", cte.Columns[1].Name, "y")
	}
}

func TestSelect_SubqueryWithAlias(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT (SELECT 1) AS val")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if sel.Targets[0].Alias.Name != "val" {
		t.Errorf("alias = %q, want %q", sel.Targets[0].Alias.Name, "val")
	}
	_, ok := sel.Targets[0].Expr.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr, got %T", sel.Targets[0].Expr)
	}
}

func TestSelect_ComplexExpression(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a + b * c, CASE WHEN x > 0 THEN 'y' ELSE 'n' END FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
	// First target is a + b * c
	_, ok := sel.Targets[0].Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("target[0]: expected *ast.BinaryExpr, got %T", sel.Targets[0].Expr)
	}
	// Second target is CASE expression
	_, ok = sel.Targets[1].Expr.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("target[1]: expected *ast.CaseExpr, got %T", sel.Targets[1].Expr)
	}
}

func TestSelect_WithSubqueryInCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) SELECT (SELECT * FROM cte)"
	sel, errs := testParseSelectStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.With) != 1 {
		t.Fatalf("With = %d, want 1", len(sel.With))
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	_, ok := sel.Targets[0].Expr.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected SubqueryExpr, got %T", sel.Targets[0].Expr)
	}
}
