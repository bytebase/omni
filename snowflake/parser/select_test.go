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

// The TOP count is a constant, not a full expression: the `*` after the
// count is the star select target, not a multiplication operator.
func TestSelect_TopStar(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT TOP 125 * FROM t")
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
	if lit.Ival != 125 {
		t.Errorf("Top = %d, want 125", lit.Ival)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
	if !sel.Targets[0].Star {
		t.Error("expected Targets[0].Star = true")
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

func TestSelect_TopColumn(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT TOP 40 c1 FROM t")
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
	if lit.Ival != 40 {
		t.Errorf("Top = %d, want 40", lit.Ival)
	}
	if len(sel.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(sel.Targets))
	}
}

func TestSelect_TopMissingCount(t *testing.T) {
	// TOP with no count expression must error.
	_, errs := testParseSelectStmt("SELECT TOP FROM t")
	if len(errs) == 0 {
		t.Fatal("expected error for TOP without count")
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

// FIRST/NEXT, ROW/ROWS, and ONLY are all optional noise words: FETCH 123
// is equivalent to FETCH FIRST 123 ROWS ONLY.
func TestSelect_FetchBare(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t FETCH 123")
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
	if litCount.Ival != 123 {
		t.Errorf("Fetch.Count = %d, want 123", litCount.Ival)
	}
	// Fetch.Loc must span FETCH..123.
	input := "SELECT a FROM t FETCH 123"
	if got := input[sel.Fetch.Loc.Start:sel.Fetch.Loc.End]; got != "FETCH 123" {
		t.Errorf("Fetch.Loc spans %q, want %q", got, "FETCH 123")
	}
}

func TestSelect_OffsetFetchBare(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t OFFSET 12 FETCH 123")
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
	if litOff.Ival != 12 {
		t.Errorf("Offset = %d, want 12", litOff.Ival)
	}
	if sel.Fetch == nil {
		t.Fatal("expected Fetch to be set")
	}
	litCount, ok := sel.Fetch.Count.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal for Fetch.Count, got %T", sel.Fetch.Count)
	}
	if litCount.Ival != 123 {
		t.Errorf("Fetch.Count = %d, want 123", litCount.Ival)
	}
}

func TestSelect_FetchFirstNoRowsOnly(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t FETCH FIRST 5")
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

func TestSelect_FetchNextRowNoOnly(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT a FROM t FETCH NEXT 7 ROW")
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
	if litCount.Ival != 7 {
		t.Errorf("Fetch.Count = %d, want 7", litCount.Ival)
	}
}

func TestSelect_FetchMissingCount(t *testing.T) {
	// FETCH with no count expression must error.
	_, errs := testParseSelectStmt("SELECT a FROM t FETCH")
	if len(errs) == 0 {
		t.Fatal("expected error for FETCH without count")
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

// ===========================================================================
// gap-select-ext: SELECT * EXCLUDE/RENAME, ORDER BY ALL, trailing comma
// ===========================================================================

// ---------------------------------------------------------------------------
// A. Star column-transforms — EXCLUDE (bare + list), RENAME (bare + list)
// ---------------------------------------------------------------------------

func TestSelect_StarExcludeBare(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * EXCLUDE department_id FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if !target.Star {
		t.Fatal("target should be star")
	}
	if len(target.Exclude) != 1 || target.Exclude[0].Name != "department_id" {
		t.Fatalf("exclude = %v, want [department_id]", target.Exclude)
	}
	if len(target.Rename) != 0 {
		t.Errorf("rename = %v, want empty", target.Rename)
	}
}

func TestSelect_StarRenameBare(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * RENAME department_id AS department FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if !target.Star {
		t.Fatal("target should be star")
	}
	if len(target.Rename) != 1 {
		t.Fatalf("rename = %d, want 1", len(target.Rename))
	}
	if target.Rename[0].Col.Name != "department_id" || target.Rename[0].Alias.Name != "department" {
		t.Errorf("rename[0] = %q AS %q, want department_id AS department",
			target.Rename[0].Col.Name, target.Rename[0].Alias.Name)
	}
}

func TestSelect_StarRenameList(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT * RENAME (department_id AS department, employee_id AS id) FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if len(target.Rename) != 2 {
		t.Fatalf("rename = %d, want 2", len(target.Rename))
	}
	if target.Rename[0].Col.Name != "department_id" || target.Rename[0].Alias.Name != "department" {
		t.Errorf("rename[0] = %q AS %q", target.Rename[0].Col.Name, target.Rename[0].Alias.Name)
	}
	if target.Rename[1].Col.Name != "employee_id" || target.Rename[1].Alias.Name != "id" {
		t.Errorf("rename[1] = %q AS %q", target.Rename[1].Col.Name, target.Rename[1].Alias.Name)
	}
}

func TestSelect_StarExcludeThenRename(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT * EXCLUDE first_name RENAME (department_id AS department, employee_id AS id) FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if len(target.Exclude) != 1 || target.Exclude[0].Name != "first_name" {
		t.Fatalf("exclude = %v, want [first_name]", target.Exclude)
	}
	if len(target.Rename) != 2 {
		t.Fatalf("rename = %d, want 2", len(target.Rename))
	}
}

func TestSelect_QualifiedStarExcludeAndRename(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT employee_table.* EXCLUDE department_id, department_table.* RENAME department_name AS department " +
			"FROM employee_table INNER JOIN department_table ON employee_table.id = department_table.id")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
	t0 := sel.Targets[0]
	if !t0.Star {
		t.Error("target[0] should be star")
	}
	star0, ok := t0.Expr.(*ast.StarExpr)
	if !ok || star0.Qualifier == nil || star0.Qualifier.Name.Name != "employee_table" {
		t.Errorf("target[0] qualifier mismatch: %#v", t0.Expr)
	}
	if len(t0.Exclude) != 1 || t0.Exclude[0].Name != "department_id" {
		t.Errorf("target[0] exclude = %v", t0.Exclude)
	}
	t1 := sel.Targets[1]
	star1, ok := t1.Expr.(*ast.StarExpr)
	if !ok || star1.Qualifier == nil || star1.Qualifier.Name.Name != "department_table" {
		t.Errorf("target[1] qualifier mismatch: %#v", t1.Expr)
	}
	if len(t1.Rename) != 1 || t1.Rename[0].Col.Name != "department_name" || t1.Rename[0].Alias.Name != "department" {
		t.Errorf("target[1] rename = %v", t1.Rename)
	}
}

func TestSelect_StarExcludeLoc(t *testing.T) {
	const sql = "SELECT * EXCLUDE department_id FROM employee_table"
	sel, errs := testParseSelectStmt(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	// The target Loc should span from the '*' through the end of the EXCLUDE col.
	got := sql[target.Loc.Start:target.Loc.End]
	if got != "* EXCLUDE department_id" {
		t.Errorf("target Loc slice = %q, want %q", got, "* EXCLUDE department_id")
	}
}

// Regression for the star-REPLACE silent drop: this statement used to
// parse-prefix to just `SELECT *` with ZERO errors — REPLACE, its pairs, AND
// the FROM clause were dropped from the AST (SelectStmt.From empty), which is
// masking-unsound for query-span consumers. It must now parse fully: Replace
// populated and From non-empty.
func TestSelect_StarReplaceSingle(t *testing.T) {
	const sql = "SELECT * REPLACE (UPPER(SSN) AS SSN) FROM T"
	sel, errs := testParseSelectStmt(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if !target.Star {
		t.Fatal("target should be star")
	}
	if len(target.Replace) != 1 {
		t.Fatalf("replace = %d, want 1", len(target.Replace))
	}
	if target.Replace[0].Col.Name != "SSN" {
		t.Errorf("replace[0].Col = %q, want %q", target.Replace[0].Col.Name, "SSN")
	}
	fn, ok := target.Replace[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("replace[0].Expr = %T, want *ast.FuncCallExpr", target.Replace[0].Expr)
	}
	if fn.Name.Normalize() != "UPPER" {
		t.Errorf("replace[0].Expr func = %q, want UPPER", fn.Name.Normalize())
	}
	// The FROM clause must survive — this is the silent-drop regression check.
	if len(sel.From) != 1 {
		t.Fatalf("From = %d entries, want 1 (FROM clause was dropped)", len(sel.From))
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok || ref.Name == nil || ref.Name.Normalize() != "T" {
		t.Errorf("From[0] = %#v, want table T", sel.From[0])
	}
	// And the statement must span the whole input, not a prefix.
	if got := sql[sel.Loc.Start:sel.Loc.End]; got != sql {
		t.Errorf("stmt Loc slice = %q, want full input", got)
	}
}

func TestSelect_StarReplaceList(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT * REPLACE ('DEPT-' || department_id AS department_id, UPPER(name) AS name) FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if len(target.Replace) != 2 {
		t.Fatalf("replace = %d, want 2", len(target.Replace))
	}
	if target.Replace[0].Col.Name != "department_id" {
		t.Errorf("replace[0].Col = %q, want department_id", target.Replace[0].Col.Name)
	}
	if _, ok := target.Replace[0].Expr.(*ast.BinaryExpr); !ok {
		t.Errorf("replace[0].Expr = %T, want *ast.BinaryExpr (concat)", target.Replace[0].Expr)
	}
	if target.Replace[1].Col.Name != "name" {
		t.Errorf("replace[1].Col = %q, want name", target.Replace[1].Col.Name)
	}
	if len(sel.From) != 1 {
		t.Fatalf("From = %d entries, want 1", len(sel.From))
	}
}

// Combined transforms in Snowflake's documented order: EXCLUDE, REPLACE, RENAME.
func TestSelect_StarExcludeReplaceRename(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT * EXCLUDE first_name REPLACE (UPPER(ssn) AS ssn) RENAME (department_id AS department) FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if len(target.Exclude) != 1 || target.Exclude[0].Name != "first_name" {
		t.Fatalf("exclude = %v, want [first_name]", target.Exclude)
	}
	if len(target.Replace) != 1 || target.Replace[0].Col.Name != "ssn" {
		t.Fatalf("replace = %v, want one ssn pair", target.Replace)
	}
	if len(target.Rename) != 1 || target.Rename[0].Col.Name != "department_id" {
		t.Fatalf("rename = %v, want one department_id pair", target.Rename)
	}
	if len(sel.From) != 1 {
		t.Fatalf("From = %d entries, want 1", len(sel.From))
	}
}

// REPLACE then RENAME (corpus official/select/example_14).
func TestSelect_StarReplaceThenRename(t *testing.T) {
	sel, errs := testParseSelectStmt(
		"SELECT * REPLACE ('DEPT-' || department_id AS department_id) RENAME department_id AS department FROM employee_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	if len(target.Replace) != 1 || target.Replace[0].Col.Name != "department_id" {
		t.Fatalf("replace = %v, want one department_id pair", target.Replace)
	}
	if len(target.Rename) != 1 || target.Rename[0].Alias.Name != "department" {
		t.Fatalf("rename = %v, want department_id AS department", target.Rename)
	}
	if len(sel.From) != 1 {
		t.Fatalf("From = %d entries, want 1", len(sel.From))
	}
}

func TestSelect_QualifiedStarReplace(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT t.* REPLACE (UPPER(x) AS x) FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	star, ok := target.Expr.(*ast.StarExpr)
	if !ok || star.Qualifier == nil || star.Qualifier.Name.Name != "t" {
		t.Fatalf("qualifier mismatch: %#v", target.Expr)
	}
	if len(target.Replace) != 1 || target.Replace[0].Col.Name != "x" {
		t.Fatalf("replace = %v, want one x pair", target.Replace)
	}
	if len(sel.From) != 1 {
		t.Fatalf("From = %d entries, want 1", len(sel.From))
	}
}

func TestSelect_StarReplaceLoc(t *testing.T) {
	const sql = "SELECT * REPLACE (UPPER(ssn) AS ssn) FROM t"
	sel, errs := testParseSelectStmt(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	target := sel.Targets[0]
	got := sql[target.Loc.Start:target.Loc.End]
	if got != "* REPLACE (UPPER(ssn) AS ssn)" {
		t.Errorf("target Loc slice = %q, want %q", got, "* REPLACE (UPPER(ssn) AS ssn)")
	}
}

// Negative: REPLACE without parentheses. Snowflake documents only the
// parenthesized form (the replacement is an arbitrary expression). This must
// be a loud error, not a silent truncation.
func TestSelect_StarReplaceNoParens(t *testing.T) {
	_, err := Parse("SELECT * REPLACE UPPER(ssn) AS ssn FROM t")
	if err == nil {
		t.Fatal("expected error for `* REPLACE` without parentheses")
	}
}

// Negative: REPLACE pair without AS.
func TestSelect_StarReplaceNoAs(t *testing.T) {
	_, err := Parse("SELECT * REPLACE (UPPER(ssn) ssn) FROM t")
	if err == nil {
		t.Fatal("expected error for `* REPLACE (expr col)` without AS")
	}
}

// Negative: empty REPLACE list.
func TestSelect_StarReplaceEmpty(t *testing.T) {
	_, err := Parse("SELECT * REPLACE () FROM t")
	if err == nil {
		t.Fatal("expected error for `* REPLACE ()` with no pairs")
	}
}

// Negative: out-of-documented-order transforms must error loudly rather than
// have the trailing transform (and everything after it, including FROM)
// silently dropped by the statement-entry trailing-token gap.
func TestSelect_StarTransformsOutOfOrder(t *testing.T) {
	for _, sql := range []string{
		"SELECT * RENAME (a AS b) REPLACE (UPPER(c) AS c) FROM t",
		"SELECT * REPLACE (UPPER(c) AS c) EXCLUDE a FROM t",
		"SELECT * RENAME (a AS b) EXCLUDE c FROM t",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("expected error for out-of-order transforms: %s", sql)
		}
	}
}

// Regression: a column literally named "replace" still parses as a plain
// column reference (REPLACE is a non-reserved keyword), and the REPLACE()
// string function is untouched by the star-transform path.
func TestSelect_ReplaceAsColumnAndFunction(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT replace, REPLACE(a, 'x', 'y') FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
	col, ok := sel.Targets[0].Expr.(*ast.ColumnRef)
	if !ok || len(col.Parts) != 1 || col.Parts[0].Name != "replace" {
		t.Errorf("target[0] = %#v, want column replace", sel.Targets[0].Expr)
	}
	fn, ok := sel.Targets[1].Expr.(*ast.FuncCallExpr)
	if !ok || fn.Name.Normalize() != "REPLACE" {
		t.Errorf("target[1] = %#v, want REPLACE() call", sel.Targets[1].Expr)
	}
}

// Negative: `* EXCLUDE` with no column.
func TestSelect_StarExcludeNoColumn(t *testing.T) {
	_, err := Parse("SELECT * EXCLUDE FROM t")
	if err == nil {
		t.Fatal("expected error for `* EXCLUDE` with no column")
	}
}

// Negative: `* RENAME x` with no AS.
func TestSelect_StarRenameNoAs(t *testing.T) {
	_, err := Parse("SELECT * RENAME x FROM t")
	if err == nil {
		t.Fatal("expected error for `* RENAME x` with no AS alias")
	}
}

// Regression: a column literally named "exclude" / "rename" still parses as a
// plain column reference (these are non-reserved keywords).
func TestSelect_ColumnNamedExcludeRename(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT exclude, rename FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
	for i, want := range []string{"exclude", "rename"} {
		col, ok := sel.Targets[i].Expr.(*ast.ColumnRef)
		if !ok {
			t.Fatalf("target[%d] = %T, want *ast.ColumnRef", i, sel.Targets[i].Expr)
		}
		if len(col.Parts) != 1 || col.Parts[0].Name != want {
			t.Errorf("target[%d] = %v, want [%q]", i, col.Parts, want)
		}
		if sel.Targets[i].Star {
			t.Errorf("target[%d] should not be star", i)
		}
	}
}

// ---------------------------------------------------------------------------
// B. ORDER BY ALL
// ---------------------------------------------------------------------------

func TestSelect_OrderByAll(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT col_1, col_2 FROM my_table ORDER BY ALL")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(sel.OrderBy))
	}
	item := sel.OrderBy[0]
	if !item.All {
		t.Error("OrderBy[0].All should be true")
	}
	if item.Expr != nil {
		t.Errorf("OrderBy[0].Expr should be nil for ALL, got %T", item.Expr)
	}
	if item.Desc {
		t.Error("OrderBy[0].Desc should be false (default ASC)")
	}
	if item.NullsFirst != nil {
		t.Error("OrderBy[0].NullsFirst should be nil")
	}
}

func TestSelect_OrderByAllAsc(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM my_sort_example ORDER BY ALL ASC")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	item := sel.OrderBy[0]
	if !item.All {
		t.Error("OrderBy[0].All should be true")
	}
	if item.Desc {
		t.Error("ALL ASC should not be Desc")
	}
}

func TestSelect_OrderByAllDesc(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM my_sort_example ORDER BY ALL DESC")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	item := sel.OrderBy[0]
	if !item.All {
		t.Error("OrderBy[0].All should be true")
	}
	if !item.Desc {
		t.Error("ALL DESC should be Desc")
	}
}

func TestSelect_OrderByAllNullsFirst(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT * FROM my_sort_example ORDER BY ALL NULLS FIRST")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	item := sel.OrderBy[0]
	if !item.All {
		t.Error("OrderBy[0].All should be true")
	}
	if item.NullsFirst == nil || !*item.NullsFirst {
		t.Errorf("ALL NULLS FIRST: NullsFirst = %v, want true", item.NullsFirst)
	}
}

func TestSelect_OrderByAllNullsLast(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT b, s, a FROM my_sort_example ORDER BY ALL NULLS LAST")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	item := sel.OrderBy[0]
	if !item.All {
		t.Error("OrderBy[0].All should be true")
	}
	if item.NullsFirst == nil || *item.NullsFirst {
		t.Errorf("ALL NULLS LAST: NullsFirst = %v, want false", item.NullsFirst)
	}
}

func TestSelect_OrderByAllLoc(t *testing.T) {
	const sql = "SELECT * FROM t ORDER BY ALL"
	sel, errs := testParseSelectStmt(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	item := sel.OrderBy[0]
	got := sql[item.Loc.Start:item.Loc.End]
	if got != "ALL" {
		t.Errorf("OrderBy[0] Loc slice = %q, want %q", got, "ALL")
	}
}

// Negative: `ORDER BY ALL NULLS <garbage>` must error just like the
// expression path does — the ALL branch must not swallow the modifier error.
func TestSelect_OrderByAllNullsGarbageErrors(t *testing.T) {
	_, err := Parse("SELECT a FROM t ORDER BY ALL NULLS SOMETHING")
	if err == nil {
		t.Fatal("expected error for `ORDER BY ALL NULLS SOMETHING`")
	}
}

// Regression: ORDER BY with real expressions still works, and a column named
// "all" inside an expression context is unaffected by the ALL marker (note:
// ALL is reserved, so a *bare* `ORDER BY all` is always the marker; this guards
// the multi-term and modifier paths).
func TestSelect_OrderByExprNotStolenByAll(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT x, y FROM t ORDER BY x, y DESC")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.OrderBy) != 2 {
		t.Fatalf("OrderBy = %d, want 2", len(sel.OrderBy))
	}
	if sel.OrderBy[0].All || sel.OrderBy[1].All {
		t.Error("neither ORDER BY term should be flagged All")
	}
	if sel.OrderBy[0].Expr == nil || sel.OrderBy[1].Expr == nil {
		t.Error("both ORDER BY terms should carry an Expr")
	}
	if !sel.OrderBy[1].Desc {
		t.Error("second ORDER BY term should be DESC")
	}
}

// Regression: `ALL` as the *argument of a function* in ORDER BY (e.g. a hand-
// written ORDER BY over a function) is not the ALL marker because it is not a
// lone term — guard that ORDER BY ALL only fires when ALL is the whole term.
// `ORDER BY x + 1` exercises the expression path adjacent to a possible ALL.
func TestSelect_OrderByExpression(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT x FROM t ORDER BY x + 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(sel.OrderBy))
	}
	if sel.OrderBy[0].All {
		t.Error("ORDER BY x + 1 should not be flagged All")
	}
}

// ---------------------------------------------------------------------------
// C. Trailing comma in SELECT list
// ---------------------------------------------------------------------------

func TestSelect_TrailingCommaBeforeFrom(t *testing.T) {
	sel, errs := testParseSelectStmt("SELECT emp_id, name, dept, FROM employees")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 3 {
		t.Fatalf("targets = %d, want 3 (trailing comma must not add an empty item)", len(sel.Targets))
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

func TestSelect_TrailingCommaNoFrom(t *testing.T) {
	// Trailing comma before end-of-statement (no FROM).
	sel, errs := testParseSelectStmt("SELECT 1, 2,")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.Targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(sel.Targets))
	}
}

func TestSelect_TrailingCommaInSubquery(t *testing.T) {
	// Trailing comma before `)` inside a derived-table subquery.
	sel, errs := testParseSelectStmt("SELECT * FROM (SELECT a, b, FROM t) x")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from = %d, want 1", len(sel.From))
	}
}

// Regression: a genuine empty item mid-list (double comma) must still error.
func TestSelect_DoubleCommaErrors(t *testing.T) {
	_, err := Parse("SELECT a, , b FROM t")
	if err == nil {
		t.Fatal("expected error for empty select item (`a, , b`)")
	}
}

// Regression: a leading comma (`SELECT , a`) must still error.
func TestSelect_LeadingCommaErrors(t *testing.T) {
	_, err := Parse("SELECT , a FROM t")
	if err == nil {
		t.Fatal("expected error for leading comma in select list")
	}
}

// ---------------------------------------------------------------------------
// D. Plain-form regressions (must be unaffected by the new code)
// ---------------------------------------------------------------------------

func TestSelect_PlainStarAndQualifiedStarRegression(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t",
		"SELECT tbl.* FROM tbl",
		"SELECT a, b FROM t",
		"SELECT a, b FROM t ORDER BY x, y DESC",
	} {
		sel, errs := testParseSelectStmt(sql)
		if len(errs) > 0 {
			t.Errorf("%q: unexpected errors: %v", sql, errs)
		}
		if sel == nil {
			t.Errorf("%q: nil select", sql)
		}
	}
}

// Regression: a column/alias literally named "all" is parseable where it is a
// legal identifier — e.g. as an alias via AS (ALL is reserved so it cannot be a
// bare column, but AS "all" quoted, or other non-bare positions are fine). Use
// the AS-quoted alias form which is unambiguously an identifier.
func TestSelect_AllAsQuotedAlias(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT x AS "all" FROM t`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if sel.Targets[0].Alias.Name != "all" {
		t.Errorf("alias = %q, want all", sel.Targets[0].Alias.Name)
	}
}

// ---------------------------------------------------------------------------
// VALUES as a table source  (gap-from-values, official/values/example_01..04)
// ---------------------------------------------------------------------------

// fromTableRef returns sel.From[i] as a *ast.TableRef or fails the test.
func fromTableRef(t *testing.T, sel *ast.SelectStmt, i int) *ast.TableRef {
	t.Helper()
	if len(sel.From) <= i {
		t.Fatalf("from len = %d, want > %d", len(sel.From), i)
	}
	ref, ok := sel.From[i].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[%d] = %T, want *ast.TableRef", i, sel.From[i])
	}
	return ref
}

func TestSelect_ValuesTableSource(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT * FROM (VALUES (1, 'one'), (2, 'two'), (3, 'three'))`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	vc, ok := ref.Subquery.(*ast.ValuesClause)
	if !ok {
		t.Fatalf("Subquery = %T, want *ast.ValuesClause", ref.Subquery)
	}
	if len(vc.Rows) != 3 {
		t.Fatalf("rows = %d, want 3", len(vc.Rows))
	}
	for i, row := range vc.Rows {
		if len(row) != 2 {
			t.Errorf("row %d arity = %d, want 2", i, len(row))
		}
	}
	if !vc.Loc.IsValid() {
		t.Errorf("ValuesClause.Loc invalid: %+v", vc.Loc)
	}
}

func TestSelect_ValuesTableSourcePositionalRef(t *testing.T) {
	// official/values/example_02: positional $2 column ref over a VALUES source.
	sel, errs := testParseSelectStmt(`SELECT column1, $2 FROM (VALUES (1, 'one'), (2, 'two'), (3, 'three'))`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	if _, ok := ref.Subquery.(*ast.ValuesClause); !ok {
		t.Fatalf("Subquery = %T, want *ast.ValuesClause", ref.Subquery)
	}
	// $2 is a positional DollarRef in the select list.
	d, ok := sel.Targets[1].Expr.(*ast.DollarRef)
	if !ok {
		t.Fatalf("target[1] = %T, want *ast.DollarRef", sel.Targets[1].Expr)
	}
	if d.Name != "2" || !d.Positional {
		t.Errorf("DollarRef = %+v, want positional $2", d)
	}
}

func TestSelect_ValuesTableSourceDerivedColumnList(t *testing.T) {
	// official/values/example_04: AS v1 (c1, c2) derived column list.
	sel, errs := testParseSelectStmt(`SELECT c1, c2 FROM (VALUES (1, 'one'), (2, 'two')) AS v1 (c1, c2)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	if ref.Alias.Name != "v1" {
		t.Errorf("alias = %q, want v1", ref.Alias.Name)
	}
	if len(ref.Columns) != 2 || ref.Columns[0].Name != "c1" || ref.Columns[1].Name != "c2" {
		t.Errorf("columns = %+v, want [c1 c2]", ref.Columns)
	}
}

func TestSelect_ValuesTableSourceJoinQualifiedDollar(t *testing.T) {
	// official/values/example_03: two aliased VALUES sources joined, with
	// qualified positional refs (v1.$2, v2.$1).
	sel, errs := testParseSelectStmt(
		`SELECT v1.$2, v2.$2 FROM (VALUES (1, 'one'), (2, 'two')) AS v1 ` +
			`INNER JOIN (VALUES (1, 'One'), (3, 'three')) AS v2 WHERE v2.$1 = v1.$1`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sel.From) != 1 {
		t.Fatalf("from len = %d, want 1 (join)", len(sel.From))
	}
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinExpr", sel.From[0])
	}
	left, ok := join.Left.(*ast.TableRef)
	if !ok || left.Alias.Name != "v1" {
		t.Errorf("join.Left = %T alias=%q, want *ast.TableRef v1", join.Left, aliasOf(join.Left))
	}
	if _, ok := left.Subquery.(*ast.ValuesClause); !ok {
		t.Errorf("join.Left.Subquery = %T, want *ast.ValuesClause", left.Subquery)
	}
}

func aliasOf(n ast.Node) string {
	if ref, ok := n.(*ast.TableRef); ok {
		return ref.Alias.Name
	}
	return ""
}

// Snowflake accepts ragged VALUES rows (differing arity); the parser must not
// reject them. Resolution of column widths is a downstream concern.
func TestSelect_ValuesTableSourceRaggedRows(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT * FROM (VALUES (1, 'a'), (2), (3, 'c', 4))`)
	if len(errs) > 0 {
		t.Fatalf("ragged VALUES rows should parse, got: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	vc := ref.Subquery.(*ast.ValuesClause)
	if got := []int{len(vc.Rows[0]), len(vc.Rows[1]), len(vc.Rows[2])}; got[0] != 2 || got[1] != 1 || got[2] != 3 {
		t.Errorf("row arities = %v, want [2 1 3]", got)
	}
}

// Regression: a plain ( SELECT … ) AS d derived table is unchanged by the
// VALUES additions.
func TestSelect_DerivedSubqueryStillParses(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT d.x FROM (SELECT 1 AS x) AS d`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	if ref.Alias.Name != "d" {
		t.Errorf("alias = %q, want d", ref.Alias.Name)
	}
	if _, ok := ref.Subquery.(*ast.SelectStmt); !ok {
		t.Errorf("Subquery = %T, want *ast.SelectStmt", ref.Subquery)
	}
}

// ---------------------------------------------------------------------------
// $N as a result-set table reference  (gap-from-values, FROM $1)
// ---------------------------------------------------------------------------

func TestSelect_DollarTableRef(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT * FROM $1 ORDER BY 1`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	if ref.DollarN == nil {
		t.Fatalf("DollarN is nil, want $1")
	}
	if ref.DollarN.Name != "1" || !ref.DollarN.Positional {
		t.Errorf("DollarN = %+v, want positional $1", ref.DollarN)
	}
	if ref.Name != nil || ref.Subquery != nil || ref.FuncCall != nil {
		t.Errorf("expected $N-only TableRef, got Name=%v Subquery=%v FuncCall=%v", ref.Name, ref.Subquery, ref.FuncCall)
	}
	if !ref.Loc.IsValid() {
		t.Errorf("TableRef.Loc invalid: %+v", ref.Loc)
	}
}

func TestSelect_DollarTableRefAlias(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT d.x FROM $1 AS d`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ref := fromTableRef(t, sel, 0)
	if ref.DollarN == nil || ref.DollarN.Name != "1" {
		t.Fatalf("DollarN = %+v, want $1", ref.DollarN)
	}
	if ref.Alias.Name != "d" {
		t.Errorf("alias = %q, want d", ref.Alias.Name)
	}
}

// Regression: $1 as a column expression (not a FROM source) still parses as a
// DollarRef in the select list.
func TestSelect_DollarColumnExprUnchanged(t *testing.T) {
	sel, errs := testParseSelectStmt(`SELECT $1, $2 FROM t`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if _, ok := sel.Targets[0].Expr.(*ast.DollarRef); !ok {
		t.Errorf("target[0] = %T, want *ast.DollarRef", sel.Targets[0].Expr)
	}
	ref := fromTableRef(t, sel, 0)
	if ref.Name == nil || ref.Name.String() != "t" {
		t.Errorf("from[0] = %+v, want table t", ref)
	}
}

// Negative: FROM $ with no name is not a valid table reference. The lexer only
// emits tokVariable for $ followed by a name, so a bare $ is a lex error.
func TestSelect_DollarBareIsError(t *testing.T) {
	result := ParseBestEffort(`SELECT * FROM $`)
	if len(result.Errors) == 0 {
		t.Error("expected an error for bare `FROM $`, got none")
	}
}
