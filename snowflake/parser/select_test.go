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
	from0, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0]: expected *ast.TableRef, got %T", sel.From[0])
	}
	if from0.Name.Name.Name != "t1" {
		t.Errorf("from[0] = %q, want %q", from0.Name.Name.Name, "t1")
	}
	if !from0.Alias.IsEmpty() {
		t.Errorf("from[0] alias should be empty, got %q", from0.Alias.Name)
	}
	from1, ok := sel.From[1].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1]: expected *ast.TableRef, got %T", sel.From[1])
	}
	if from1.Name.Name.Name != "t2" {
		t.Errorf("from[1] = %q, want %q", from1.Name.Name.Name, "t2")
	}
	if from1.Alias.Name != "x" {
		t.Errorf("from[1] alias = %q, want %q", from1.Alias.Name, "x")
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
	from0, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0]: expected *ast.TableRef, got %T", sel.From[0])
	}
	if !from0.Alias.IsEmpty() {
		t.Errorf("table alias should be empty, got %q", from0.Alias.Name)
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
	from0, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0]: expected *ast.TableRef, got %T", sel.From[0])
	}
	if from0.Alias.Name != "x" {
		t.Errorf("from[0] alias = %q, want %q", from0.Alias.Name, "x")
	}
	from1, ok := sel.From[1].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1]: expected *ast.TableRef, got %T", sel.From[1])
	}
	if from1.Alias.Name != "y" {
		t.Errorf("from[1] alias = %q, want %q", from1.Alias.Name, "y")
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
	from0, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0]: expected *ast.TableRef, got %T", sel.From[0])
	}
	name := from0.Name
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

// ---------------------------------------------------------------------------
// T1.7 Set operator tests
// ---------------------------------------------------------------------------

// testParseSetOp is a helper that parses input and returns the first statement
// as *ast.SetOperationStmt plus any errors.
func testParseSetOp(input string) (*ast.SetOperationStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	node, ok := result.File.Stmts[0].(*ast.SetOperationStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a SetOperationStmt"})
	}
	return node, result.Errors
}

// TestSetOp_Union verifies UNION parsing.
func TestSetOp_Union(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 UNION SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("Op = %v, want SetOpUnion", node.Op)
	}
	if node.All {
		t.Error("All should be false")
	}
	if node.ByName {
		t.Error("ByName should be false")
	}
	if _, ok := node.Left.(*ast.SelectStmt); !ok {
		t.Errorf("Left: expected *SelectStmt, got %T", node.Left)
	}
	if _, ok := node.Right.(*ast.SelectStmt); !ok {
		t.Errorf("Right: expected *SelectStmt, got %T", node.Right)
	}
}

// TestSetOp_UnionAll verifies UNION ALL parsing.
func TestSetOp_UnionAll(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 UNION ALL SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("Op = %v, want SetOpUnion", node.Op)
	}
	if !node.All {
		t.Error("All should be true")
	}
	if node.ByName {
		t.Error("ByName should be false")
	}
}

// TestSetOp_UnionAllByName verifies UNION ALL BY NAME parsing.
func TestSetOp_UnionAllByName(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 UNION ALL BY NAME SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("Op = %v, want SetOpUnion", node.Op)
	}
	if !node.All {
		t.Error("All should be true")
	}
	if !node.ByName {
		t.Error("ByName should be true")
	}
}

// TestSetOp_Except verifies EXCEPT parsing.
func TestSetOp_Except(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 EXCEPT SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpExcept {
		t.Errorf("Op = %v, want SetOpExcept", node.Op)
	}
}

// TestSetOp_Minus verifies MINUS is treated as EXCEPT.
func TestSetOp_Minus(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 MINUS SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpExcept {
		t.Errorf("Op = %v, want SetOpExcept (MINUS maps to EXCEPT)", node.Op)
	}
}

// TestSetOp_Intersect verifies INTERSECT parsing.
func TestSetOp_Intersect(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 INTERSECT SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpIntersect {
		t.Errorf("Op = %v, want SetOpIntersect", node.Op)
	}
}

// TestSetOp_Chained verifies left-associative chaining:
// SELECT 1 UNION SELECT 2 UNION SELECT 3
// → SetOperationStmt{Left: SetOperationStmt{...}, Right: SelectStmt}
func TestSetOp_Chained(t *testing.T) {
	node, errs := testParseSetOp("SELECT 1 UNION SELECT 2 UNION SELECT 3")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("outer Op = %v, want SetOpUnion", node.Op)
	}
	inner, ok := node.Left.(*ast.SetOperationStmt)
	if !ok {
		t.Fatalf("Left: expected *SetOperationStmt (left-assoc), got %T", node.Left)
	}
	if inner.Op != ast.SetOpUnion {
		t.Errorf("inner Op = %v, want SetOpUnion", inner.Op)
	}
	if _, ok := inner.Left.(*ast.SelectStmt); !ok {
		t.Errorf("inner.Left: expected *SelectStmt, got %T", inner.Left)
	}
	if _, ok := inner.Right.(*ast.SelectStmt); !ok {
		t.Errorf("inner.Right: expected *SelectStmt, got %T", inner.Right)
	}
	if _, ok := node.Right.(*ast.SelectStmt); !ok {
		t.Errorf("outer Right: expected *SelectStmt, got %T", node.Right)
	}
}

// TestSetOp_Parenthesized verifies (SELECT 1) UNION (SELECT 2).
func TestSetOp_Parenthesized(t *testing.T) {
	node, errs := testParseSetOp("(SELECT 1) UNION (SELECT 2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("Op = %v, want SetOpUnion", node.Op)
	}
}

// TestSetOp_WithCTE verifies WITH cte AS (SELECT 1) SELECT * FROM cte UNION SELECT 2.
func TestSetOp_WithCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) SELECT * FROM cte UNION SELECT 2"
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected at least one statement")
	}
	node, ok := result.File.Stmts[0].(*ast.SetOperationStmt)
	if !ok {
		t.Fatalf("expected *SetOperationStmt, got %T", result.File.Stmts[0])
	}
	if node.Op != ast.SetOpUnion {
		t.Errorf("Op = %v, want SetOpUnion", node.Op)
	}
	// Left side should be a SelectStmt (with With CTE attached).
	leftSel, ok := node.Left.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Left: expected *SelectStmt, got %T", node.Left)
	}
	if len(leftSel.With) != 1 {
		t.Errorf("With = %d, want 1 CTE", len(leftSel.With))
	}
	if _, ok := node.Right.(*ast.SelectStmt); !ok {
		t.Errorf("Right: expected *SelectStmt, got %T", node.Right)
	}
}

// ===========================================================================
// JOIN tests (T1.5)
// ===========================================================================

// ---------------------------------------------------------------------------
// J1. Basic INNER JOIN
// ---------------------------------------------------------------------------

func TestJoin_InnerJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if join.Natural {
		t.Error("expected Natural = false")
	}
	if join.Directed {
		t.Error("expected Directed = false")
	}

	// Left and Right should be TableRef
	left, ok := join.Left.(*ast.TableRef)
	if !ok {
		t.Fatalf("left: expected *ast.TableRef, got %T", join.Left)
	}
	if left.Name.Name.Name != "t1" {
		t.Errorf("left = %q, want %q", left.Name.Name.Name, "t1")
	}
	right, ok := join.Right.(*ast.TableRef)
	if !ok {
		t.Fatalf("right: expected *ast.TableRef, got %T", join.Right)
	}
	if right.Name.Name.Name != "t2" {
		t.Errorf("right = %q, want %q", right.Name.Name.Name, "t2")
	}

	// ON condition
	if join.On == nil {
		t.Fatal("expected ON condition")
	}
	_, ok = join.On.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("ON: expected *ast.BinaryExpr, got %T", join.On)
	}
}

func TestJoin_ExplicitInnerJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if join.On == nil {
		t.Fatal("expected ON condition")
	}
}

// ---------------------------------------------------------------------------
// J2. LEFT JOIN
// ---------------------------------------------------------------------------

func TestJoin_LeftJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %v, want JoinLeft", join.Type)
	}
	if join.On == nil {
		t.Fatal("expected ON condition")
	}
}

func TestJoin_LeftOuterJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %v, want JoinLeft", join.Type)
	}
}

// ---------------------------------------------------------------------------
// J3. RIGHT OUTER JOIN
// ---------------------------------------------------------------------------

func TestJoin_RightOuterJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 RIGHT OUTER JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinRight {
		t.Errorf("join type = %v, want JoinRight", join.Type)
	}
}

func TestJoin_RightJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinRight {
		t.Errorf("join type = %v, want JoinRight", join.Type)
	}
}

// ---------------------------------------------------------------------------
// J4. FULL JOIN
// ---------------------------------------------------------------------------

func TestJoin_FullJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 FULL JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinFull {
		t.Errorf("join type = %v, want JoinFull", join.Type)
	}
}

func TestJoin_FullOuterJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinFull {
		t.Errorf("join type = %v, want JoinFull", join.Type)
	}
}

// ---------------------------------------------------------------------------
// J5. CROSS JOIN (no ON)
// ---------------------------------------------------------------------------

func TestJoin_CrossJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 CROSS JOIN t2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinCross {
		t.Errorf("join type = %v, want JoinCross", join.Type)
	}
	if join.On != nil {
		t.Error("CROSS JOIN should not have ON condition")
	}
	if join.Using != nil {
		t.Error("CROSS JOIN should not have USING")
	}
}

// ---------------------------------------------------------------------------
// J6. NATURAL JOIN (no ON/USING)
// ---------------------------------------------------------------------------

func TestJoin_NaturalJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 NATURAL JOIN t2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if !join.Natural {
		t.Error("expected Natural = true")
	}
	if join.On != nil {
		t.Error("NATURAL JOIN should not have ON")
	}
	if join.Using != nil {
		t.Error("NATURAL JOIN should not have USING")
	}
}

func TestJoin_NaturalLeftJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 NATURAL LEFT JOIN t2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %v, want JoinLeft", join.Type)
	}
	if !join.Natural {
		t.Error("expected Natural = true")
	}
}

func TestJoin_NaturalRightJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 NATURAL RIGHT JOIN t2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinRight {
		t.Errorf("join type = %v, want JoinRight", join.Type)
	}
	if !join.Natural {
		t.Error("expected Natural = true")
	}
}

func TestJoin_NaturalFullJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 NATURAL FULL JOIN t2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinFull {
		t.Errorf("join type = %v, want JoinFull", join.Type)
	}
	if !join.Natural {
		t.Error("expected Natural = true")
	}
}

// ---------------------------------------------------------------------------
// J7. USING clause
// ---------------------------------------------------------------------------

func TestJoin_UsingClause(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 JOIN t2 USING (id)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if join.On != nil {
		t.Error("USING join should not have ON")
	}
	if len(join.Using) != 1 {
		t.Fatalf("Using = %d, want 1", len(join.Using))
	}
	if join.Using[0].Name != "id" {
		t.Errorf("Using[0] = %q, want %q", join.Using[0].Name, "id")
	}
}

func TestJoin_UsingMultipleColumns(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 JOIN t2 USING (id, name)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if len(join.Using) != 2 {
		t.Fatalf("Using = %d, want 2", len(join.Using))
	}
	if join.Using[0].Name != "id" {
		t.Errorf("Using[0] = %q, want %q", join.Using[0].Name, "id")
	}
	if join.Using[1].Name != "name" {
		t.Errorf("Using[1] = %q, want %q", join.Using[1].Name, "name")
	}
}

// ---------------------------------------------------------------------------
// J8. Chained JOINs — left-associative
// ---------------------------------------------------------------------------

func TestJoin_ChainedJoins(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 JOIN t2 ON t1.id = t2.id JOIN t3 ON t2.id = t3.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	// Outer join: (t1 JOIN t2) JOIN t3
	outerJoin, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected outer *ast.JoinExpr, got %T", sel.From[0])
	}
	// Left of outer should be inner join
	innerJoin, ok := outerJoin.Left.(*ast.JoinExpr)
	if !ok {
		t.Fatalf("outer.Left: expected *ast.JoinExpr, got %T", outerJoin.Left)
	}
	// innerJoin.Left = t1, innerJoin.Right = t2
	t1, ok := innerJoin.Left.(*ast.TableRef)
	if !ok {
		t.Fatalf("inner.Left: expected *ast.TableRef, got %T", innerJoin.Left)
	}
	if t1.Name.Name.Name != "t1" {
		t.Errorf("inner.Left = %q, want %q", t1.Name.Name.Name, "t1")
	}
	t2, ok := innerJoin.Right.(*ast.TableRef)
	if !ok {
		t.Fatalf("inner.Right: expected *ast.TableRef, got %T", innerJoin.Right)
	}
	if t2.Name.Name.Name != "t2" {
		t.Errorf("inner.Right = %q, want %q", t2.Name.Name.Name, "t2")
	}
	// outerJoin.Right = t3
	t3, ok := outerJoin.Right.(*ast.TableRef)
	if !ok {
		t.Fatalf("outer.Right: expected *ast.TableRef, got %T", outerJoin.Right)
	}
	if t3.Name.Name.Name != "t3" {
		t.Errorf("outer.Right = %q, want %q", t3.Name.Name.Name, "t3")
	}
}

// ---------------------------------------------------------------------------
// J9. Comma + JOIN mixed
// ---------------------------------------------------------------------------

func TestJoin_CommaAndJoinMixed(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1, t2 JOIN t3 ON t2.id = t3.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 2 {
		t.Fatalf("from = %d, want 2", len(sel.From))
	}
	// From[0] = TableRef(t1)
	ref0, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0]: expected *ast.TableRef, got %T", sel.From[0])
	}
	if ref0.Name.Name.Name != "t1" {
		t.Errorf("from[0] = %q, want %q", ref0.Name.Name.Name, "t1")
	}
	// From[1] = JoinExpr(t2 JOIN t3)
	join1, ok := sel.From[1].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("from[1]: expected *ast.JoinExpr, got %T", sel.From[1])
	}
	left, ok := join1.Left.(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1].Left: expected *ast.TableRef, got %T", join1.Left)
	}
	if left.Name.Name.Name != "t2" {
		t.Errorf("from[1].Left = %q, want %q", left.Name.Name.Name, "t2")
	}
	right, ok := join1.Right.(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1].Right: expected *ast.TableRef, got %T", join1.Right)
	}
	if right.Name.Name.Name != "t3" {
		t.Errorf("from[1].Right = %q, want %q", right.Name.Name.Name, "t3")
	}
}

// ---------------------------------------------------------------------------
// J10. Subquery in FROM
// ---------------------------------------------------------------------------

func TestJoin_SubqueryInFrom(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM (SELECT 1 AS x) AS sub")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("expected *ast.TableRef, got %T", sel.From[0])
	}
	if ref.Subquery == nil {
		t.Fatal("expected Subquery to be set")
	}
	if ref.Name != nil {
		t.Error("Name should be nil for subquery source")
	}
	innerSel, ok := ref.Subquery.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Subquery: expected *ast.SelectStmt, got %T", ref.Subquery)
	}
	if len(innerSel.Targets) != 1 {
		t.Errorf("inner targets = %d, want 1", len(innerSel.Targets))
	}
	if ref.Alias.Name != "sub" {
		t.Errorf("alias = %q, want %q", ref.Alias.Name, "sub")
	}
}

func TestJoin_SubqueryJoinTable(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM (SELECT 1 AS id) AS sub JOIN t2 ON sub.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	left, ok := join.Left.(*ast.TableRef)
	if !ok {
		t.Fatalf("Left: expected *ast.TableRef, got %T", join.Left)
	}
	if left.Subquery == nil {
		t.Fatal("Left.Subquery should be set")
	}
	if left.Alias.Name != "sub" {
		t.Errorf("Left.Alias = %q, want %q", left.Alias.Name, "sub")
	}
}

// ---------------------------------------------------------------------------
// J11. LATERAL subquery
// ---------------------------------------------------------------------------

func TestJoin_LateralSubquery(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1, LATERAL (SELECT t1.id) AS lat")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 2 {
		t.Fatalf("from = %d, want 2", len(sel.From))
	}
	ref, ok := sel.From[1].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1]: expected *ast.TableRef, got %T", sel.From[1])
	}
	if !ref.Lateral {
		t.Error("expected Lateral = true")
	}
	if ref.Subquery == nil {
		t.Fatal("expected Subquery to be set")
	}
	if ref.Alias.Name != "lat" {
		t.Errorf("alias = %q, want %q", ref.Alias.Name, "lat")
	}
}

// ---------------------------------------------------------------------------
// J12. TABLE(FLATTEN(...)) — table function
// ---------------------------------------------------------------------------

func TestJoin_TableFlatten(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM TABLE(FLATTEN(v))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("expected *ast.TableRef, got %T", sel.From[0])
	}
	if ref.FuncCall == nil {
		t.Fatal("expected FuncCall to be set")
	}
	if ref.FuncCall.Name.Name.Name != "FLATTEN" {
		t.Errorf("FuncCall name = %q, want %q", ref.FuncCall.Name.Name.Name, "FLATTEN")
	}
	if ref.Name != nil {
		t.Error("Name should be nil for table function source")
	}
}

func TestJoin_TableFlattenWithAlias(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT f.value FROM TABLE(FLATTEN(v)) AS f")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("expected *ast.TableRef, got %T", sel.From[0])
	}
	if ref.FuncCall == nil {
		t.Fatal("expected FuncCall to be set")
	}
	if ref.Alias.Name != "f" {
		t.Errorf("alias = %q, want %q", ref.Alias.Name, "f")
	}
}

// ---------------------------------------------------------------------------
// J13. ASOF JOIN with MATCH_CONDITION
// ---------------------------------------------------------------------------

func TestJoin_AsofJoinMatchCondition(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 ASOF JOIN t2 MATCH_CONDITION (t1.ts >= t2.ts)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinAsof {
		t.Errorf("join type = %v, want JoinAsof", join.Type)
	}
	if join.MatchCondition == nil {
		t.Fatal("expected MatchCondition to be set")
	}
	_, ok = join.MatchCondition.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("MatchCondition: expected *ast.BinaryExpr, got %T", join.MatchCondition)
	}
}

func TestJoin_AsofJoinMatchConditionAndOn(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 ASOF JOIN t2 MATCH_CONDITION (t1.ts >= t2.ts) ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinAsof {
		t.Errorf("join type = %v, want JoinAsof", join.Type)
	}
	if join.MatchCondition == nil {
		t.Fatal("expected MatchCondition to be set")
	}
	if join.On == nil {
		t.Fatal("expected ON condition to be set")
	}
}

// ---------------------------------------------------------------------------
// J14. DIRECTED JOIN
// ---------------------------------------------------------------------------

func TestJoin_DirectedInnerJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 DIRECTED INNER JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if !join.Directed {
		t.Error("expected Directed = true")
	}
}

func TestJoin_DirectedLeftJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 DIRECTED LEFT JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %v, want JoinLeft", join.Type)
	}
	if !join.Directed {
		t.Error("expected Directed = true")
	}
}

func TestJoin_DirectedBareJoin(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 DIRECTED JOIN t2 ON t1.id = t2.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %v, want JoinInner", join.Type)
	}
	if !join.Directed {
		t.Error("expected Directed = true")
	}
}

// ---------------------------------------------------------------------------
// J15. Error cases
// ---------------------------------------------------------------------------

func TestJoin_ErrorJoinWithoutTable(t *testing.T) {
	// JOIN with no right-hand table — should produce an error
	_, errs := testParseSelectStmt("SELECT * FROM t1 JOIN")
	if len(errs) == 0 {
		t.Fatal("expected error for JOIN without right table")
	}
}

func TestJoin_ErrorSubqueryNoClosingParen(t *testing.T) {
	// Missing closing paren for subquery
	_, errs := testParseSelectStmt("SELECT * FROM (SELECT 1")
	if len(errs) == 0 {
		t.Fatal("expected error for unclosed subquery in FROM")
	}
}

func TestJoin_JoinWithWhereClause(t *testing.T) {
	// Full query with JOIN + WHERE to ensure they compose correctly
	sel, errs := testParseSelectStmt("SELECT * FROM t1 JOIN t2 ON t1.id = t2.id WHERE t1.active = 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
	_, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if sel.Where == nil {
		t.Fatal("expected WHERE to be set")
	}
}

func TestJoin_JoinWithOrderByLimit(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id ORDER BY t1.id LIMIT 10")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("expected *ast.JoinExpr, got %T", sel.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %v, want JoinLeft", join.Type)
	}
	if len(sel.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(sel.OrderBy))
	}
	if sel.Limit == nil {
		t.Fatal("expected LIMIT to be set")
	}
}
