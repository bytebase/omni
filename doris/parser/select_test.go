package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseSelect parses input and returns the first statement as *ast.SelectStmt.
func mustParseSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.SelectStmt", input, file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Basic SELECT
// ---------------------------------------------------------------------------

func TestSelectLiteral(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT 1")
	if len(stmt.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(stmt.Items))
	}
	item := stmt.Items[0]
	if item.Star {
		t.Error("item.Star = true, want false")
	}
	lit, ok := item.Expr.(*ast.Literal)
	if !ok {
		t.Fatalf("item.Expr = %T, want *ast.Literal", item.Expr)
	}
	if lit.Kind != ast.LitInt || lit.Value != "1" {
		t.Errorf("literal = %v %q, want LitInt 1", lit.Kind, lit.Value)
	}
}

func TestSelectColumnList(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a, b, c FROM t")
	if len(stmt.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(stmt.Items))
	}
	for i, name := range []string{"a", "b", "c"} {
		item := stmt.Items[i]
		ref, ok := item.Expr.(*ast.ColumnRef)
		if !ok {
			t.Fatalf("item[%d].Expr = %T, want *ast.ColumnRef", i, item.Expr)
		}
		if ref.Name.Parts[0] != name {
			t.Errorf("item[%d] name = %q, want %q", i, ref.Name.Parts[0], name)
		}
	}
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	tbl, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if tbl.Name.Parts[0] != "t" {
		t.Errorf("table name = %q, want %q", tbl.Name.Parts[0], "t")
	}
}

func TestSelectStar(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM t")
	if len(stmt.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(stmt.Items))
	}
	if !stmt.Items[0].Star {
		t.Error("item.Star = false, want true")
	}
	if stmt.Items[0].TableName != nil {
		t.Error("item.TableName should be nil for bare *")
	}
}

func TestSelectQualifiedStar(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT t.* FROM t")
	if len(stmt.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(stmt.Items))
	}
	item := stmt.Items[0]
	if !item.Star {
		t.Error("item.Star = false, want true")
	}
	if item.TableName == nil {
		t.Fatal("item.TableName = nil, want non-nil")
	}
	if item.TableName.Parts[0] != "t" {
		t.Errorf("table name = %q, want %q", item.TableName.Parts[0], "t")
	}
}

func TestSelectAliases(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a AS col1, b col2 FROM t")
	if len(stmt.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(stmt.Items))
	}
	if stmt.Items[0].Alias != "col1" {
		t.Errorf("item[0].Alias = %q, want %q", stmt.Items[0].Alias, "col1")
	}
	if stmt.Items[1].Alias != "col2" {
		t.Errorf("item[1].Alias = %q, want %q", stmt.Items[1].Alias, "col2")
	}
}

func TestSelectDistinct(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT DISTINCT a FROM t")
	if !stmt.Distinct {
		t.Error("Distinct = false, want true")
	}
	if len(stmt.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(stmt.Items))
	}
}

// ---------------------------------------------------------------------------
// WHERE clause
// ---------------------------------------------------------------------------

func TestSelectWhere(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t WHERE a > 1")
	if stmt.Where == nil {
		t.Fatal("Where = nil, want non-nil")
	}
	binExpr, ok := stmt.Where.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Where = %T, want *ast.BinaryExpr", stmt.Where)
	}
	if binExpr.Op != ast.BinGt {
		t.Errorf("Where.Op = %v, want BinGt", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// GROUP BY clause
// ---------------------------------------------------------------------------

func TestSelectGroupBy(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a, COUNT(*) FROM t GROUP BY a")
	if len(stmt.GroupBy) != 1 {
		t.Fatalf("GroupBy = %d, want 1", len(stmt.GroupBy))
	}
	ref, ok := stmt.GroupBy[0].(*ast.ColumnRef)
	if !ok {
		t.Fatalf("GroupBy[0] = %T, want *ast.ColumnRef", stmt.GroupBy[0])
	}
	if ref.Name.Parts[0] != "a" {
		t.Errorf("GroupBy[0] name = %q, want %q", ref.Name.Parts[0], "a")
	}
	// Check that COUNT(*) is in the select list
	if len(stmt.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(stmt.Items))
	}
	fc, ok := stmt.Items[1].Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("item[1].Expr = %T, want *ast.FuncCallExpr", stmt.Items[1].Expr)
	}
	if !fc.Star {
		t.Error("COUNT(*) Star = false, want true")
	}
}

// ---------------------------------------------------------------------------
// HAVING clause
// ---------------------------------------------------------------------------

func TestSelectHaving(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a, COUNT(*) FROM t GROUP BY a HAVING COUNT(*) > 1")
	if stmt.Having == nil {
		t.Fatal("Having = nil, want non-nil")
	}
	binExpr, ok := stmt.Having.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Having = %T, want *ast.BinaryExpr", stmt.Having)
	}
	if binExpr.Op != ast.BinGt {
		t.Errorf("Having.Op = %v, want BinGt", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// ORDER BY clause
// ---------------------------------------------------------------------------

func TestSelectOrderByDesc(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t ORDER BY a DESC")
	if len(stmt.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(stmt.OrderBy))
	}
	if !stmt.OrderBy[0].Desc {
		t.Error("OrderBy[0].Desc = false, want true")
	}
}

func TestSelectOrderByNullsFirst(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t ORDER BY a ASC NULLS FIRST")
	if len(stmt.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(stmt.OrderBy))
	}
	item := stmt.OrderBy[0]
	if item.Desc {
		t.Error("OrderBy[0].Desc = true, want false")
	}
	if item.NullsFirst == nil {
		t.Fatal("OrderBy[0].NullsFirst = nil, want non-nil")
	}
	if !*item.NullsFirst {
		t.Error("OrderBy[0].NullsFirst = false, want true")
	}
}

func TestSelectOrderByNullsLast(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t ORDER BY a DESC NULLS LAST")
	if len(stmt.OrderBy) != 1 {
		t.Fatalf("OrderBy = %d, want 1", len(stmt.OrderBy))
	}
	item := stmt.OrderBy[0]
	if !item.Desc {
		t.Error("OrderBy[0].Desc = false, want true")
	}
	if item.NullsFirst == nil {
		t.Fatal("OrderBy[0].NullsFirst = nil, want non-nil")
	}
	if *item.NullsFirst {
		t.Error("OrderBy[0].NullsFirst = true, want false")
	}
}

// ---------------------------------------------------------------------------
// LIMIT / OFFSET
// ---------------------------------------------------------------------------

func TestSelectLimit(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t LIMIT 10")
	if stmt.Limit == nil {
		t.Fatal("Limit = nil, want non-nil")
	}
	lit, ok := stmt.Limit.(*ast.Literal)
	if !ok {
		t.Fatalf("Limit = %T, want *ast.Literal", stmt.Limit)
	}
	if lit.Value != "10" {
		t.Errorf("Limit.Value = %q, want %q", lit.Value, "10")
	}
	if stmt.Offset != nil {
		t.Error("Offset should be nil")
	}
}

func TestSelectLimitOffset(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t LIMIT 10 OFFSET 5")
	if stmt.Limit == nil {
		t.Fatal("Limit = nil, want non-nil")
	}
	lit, ok := stmt.Limit.(*ast.Literal)
	if !ok {
		t.Fatalf("Limit = %T, want *ast.Literal", stmt.Limit)
	}
	if lit.Value != "10" {
		t.Errorf("Limit.Value = %q, want %q", lit.Value, "10")
	}
	if stmt.Offset == nil {
		t.Fatal("Offset = nil, want non-nil")
	}
	offLit, ok := stmt.Offset.(*ast.Literal)
	if !ok {
		t.Fatalf("Offset = %T, want *ast.Literal", stmt.Offset)
	}
	if offLit.Value != "5" {
		t.Errorf("Offset.Value = %q, want %q", offLit.Value, "5")
	}
}

// ---------------------------------------------------------------------------
// FROM clause variations
// ---------------------------------------------------------------------------

func TestSelectCommaJoin(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t1, t2 WHERE t1.id = t2.id")
	if len(stmt.From) != 2 {
		t.Fatalf("from = %d, want 2", len(stmt.From))
	}
	tbl1, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if tbl1.Name.Parts[0] != "t1" {
		t.Errorf("from[0] name = %q, want %q", tbl1.Name.Parts[0], "t1")
	}
	tbl2, ok := stmt.From[1].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[1] = %T, want *ast.TableRef", stmt.From[1])
	}
	if tbl2.Name.Parts[0] != "t2" {
		t.Errorf("from[1] name = %q, want %q", tbl2.Name.Parts[0], "t2")
	}
}

func TestSelectTableAlias(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t AS alias")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	tbl, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if tbl.Name.Parts[0] != "t" {
		t.Errorf("table name = %q, want %q", tbl.Name.Parts[0], "t")
	}
	if tbl.Alias != "alias" {
		t.Errorf("table alias = %q, want %q", tbl.Alias, "alias")
	}
}

func TestSelectTableAliasImplicit(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t myalias")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	tbl, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if tbl.Alias != "myalias" {
		t.Errorf("table alias = %q, want %q", tbl.Alias, "myalias")
	}
}

func TestSelectQualifiedTableName(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM db.t")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	tbl, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if len(tbl.Name.Parts) != 2 {
		t.Fatalf("table name parts = %d, want 2", len(tbl.Name.Parts))
	}
	if tbl.Name.Parts[0] != "db" {
		t.Errorf("table name[0] = %q, want %q", tbl.Name.Parts[0], "db")
	}
	if tbl.Name.Parts[1] != "t" {
		t.Errorf("table name[1] = %q, want %q", tbl.Name.Parts[1], "t")
	}
}

// ---------------------------------------------------------------------------
// Full combination
// ---------------------------------------------------------------------------

func TestSelectFullCombo(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t WHERE a > 1 ORDER BY a LIMIT 10")
	if stmt.Where == nil {
		t.Error("Where = nil")
	}
	if len(stmt.OrderBy) != 1 {
		t.Errorf("OrderBy = %d, want 1", len(stmt.OrderBy))
	}
	if stmt.Limit == nil {
		t.Error("Limit = nil")
	}
}

// ---------------------------------------------------------------------------
// JOIN (basic)
// ---------------------------------------------------------------------------

func TestSelectInnerJoin(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t1 JOIN t2 ON t1.id = t2.id")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	join, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if join.Type != ast.JoinInner {
		t.Errorf("join type = %d, want JoinInner", join.Type)
	}
	if join.On == nil {
		t.Error("join.On = nil, want non-nil")
	}
}

func TestSelectLeftJoin(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t1 LEFT JOIN t2 ON t1.id = t2.id")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	join, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if join.Type != ast.JoinLeft {
		t.Errorf("join type = %d, want JoinLeft", join.Type)
	}
}

func TestSelectCrossJoin(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t1 CROSS JOIN t2")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	join, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if join.Type != ast.JoinCross {
		t.Errorf("join type = %d, want JoinCross", join.Type)
	}
}

// ---------------------------------------------------------------------------
// QUALIFY clause
// ---------------------------------------------------------------------------

func TestSelectQualify(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a FROM t QUALIFY ROW_NUMBER() OVER() = 1")
	if stmt.Qualify == nil {
		t.Fatal("Qualify = nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// Multiple statements
// ---------------------------------------------------------------------------

func TestSelectMultipleStatements(t *testing.T) {
	file, errs := Parse("SELECT 1; SELECT 2")
	if len(errs) > 0 {
		t.Fatalf("errors: %v", errs)
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("stmts = %d, want 2", len(file.Stmts))
	}
	for i, s := range file.Stmts {
		if _, ok := s.(*ast.SelectStmt); !ok {
			t.Errorf("stmt[%d] = %T, want *ast.SelectStmt", i, s)
		}
	}
}

// ---------------------------------------------------------------------------
// Node tag verification
// ---------------------------------------------------------------------------

func TestSelectStmtTag(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT 1")
	if stmt.Tag() != ast.T_SelectStmt {
		t.Errorf("Tag() = %v, want T_SelectStmt", stmt.Tag())
	}
}

// ---------------------------------------------------------------------------
// Loc verification
// ---------------------------------------------------------------------------

func TestSelectLoc(t *testing.T) {
	input := "SELECT a FROM t"
	stmt := mustParseSelect(t, input)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// Constructs previously hidden by the trailing-token swallow (BYT-10084)
// ---------------------------------------------------------------------------

func TestSelectLimitOffsetComma(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM tbl LIMIT 5, 10")
	lit, ok := stmt.Limit.(*ast.Literal)
	if !ok || lit.Value != "10" {
		t.Fatalf("Limit = %+v, want literal 10", stmt.Limit)
	}
	off, ok := stmt.Offset.(*ast.Literal)
	if !ok || off.Value != "5" {
		t.Fatalf("Offset = %+v, want literal 5", stmt.Offset)
	}

	// The uint64-max count from the legacy corpus must round-trip textually.
	stmt = mustParseSelect(t, "SELECT * FROM tbl LIMIT 95, 18446744073709551615")
	lit, ok = stmt.Limit.(*ast.Literal)
	if !ok || lit.Value != "18446744073709551615" {
		t.Fatalf("Limit = %+v, want literal 18446744073709551615", stmt.Limit)
	}
}

func TestSelectStringAlias(t *testing.T) {
	// AS "string": the alias used to be dropped together with everything
	// after it, so the span lost the FROM table.
	stmt := mustParseSelect(t, `SELECT *, (price * 0.8) AS "20%" FROM tb_book`)
	if len(stmt.Items) != 2 || stmt.Items[1].Alias != "20%" {
		t.Fatalf("Items[1].Alias = %+v, want 20%%", stmt.Items)
	}
	if len(stmt.From) != 1 {
		t.Fatalf("From = %+v, want tb_book to survive the alias", stmt.From)
	}

	stmt = mustParseSelect(t, "SELECT 1 AS 'x'")
	if stmt.Items[0].Alias != "x" {
		t.Errorf("Alias = %q, want x", stmt.Items[0].Alias)
	}

	// The empty alias AS '' is engine-valid and distinct from no alias:
	// presence is carried by Aliased, not by the string value.
	stmt = mustParseSelect(t, "SELECT c AS '' FROM t")
	if !stmt.Items[0].Aliased || stmt.Items[0].Alias != "" {
		t.Errorf("AS '': Aliased=%v Alias=%q, want true and empty", stmt.Items[0].Aliased, stmt.Items[0].Alias)
	}
	stmt = mustParseSelect(t, "SELECT c FROM t")
	if stmt.Items[0].Aliased {
		t.Error("unaliased item reports Aliased=true")
	}
}

func TestSelectGroupByWithRollup(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP")
	if !stmt.GroupByWithRollup {
		t.Error("GroupByWithRollup = false, want true")
	}
	if len(stmt.GroupBy) != 1 {
		t.Errorf("GroupBy len = %d, want 1", len(stmt.GroupBy))
	}

	stmt = mustParseSelect(t, "SELECT a FROM t GROUP BY a")
	if stmt.GroupByWithRollup {
		t.Error("GroupByWithRollup = true, want false")
	}
}

func TestSelectGroupByGroupingSets(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT a, SUM(b) FROM t GROUP BY GROUPING SETS ((a, b), (a), ())")
	if len(stmt.GroupBy) != 1 {
		t.Fatalf("GroupBy len = %d, want 1", len(stmt.GroupBy))
	}
	gs, ok := stmt.GroupBy[0].(*ast.GroupingSetsExpr)
	if !ok {
		t.Fatalf("GroupBy[0] = %T, want *ast.GroupingSetsExpr", stmt.GroupBy[0])
	}
	if len(gs.Sets) != 3 || len(gs.Sets[0]) != 2 || len(gs.Sets[1]) != 1 || len(gs.Sets[2]) != 0 {
		t.Fatalf("Sets shape = %v, want [2 1 0]", gs.Sets)
	}

	// Like CUBE, GROUPING SETS is the whole grouping specification.
	_, errs := Parse("SELECT a FROM t GROUP BY GROUPING SETS ((a)), b")
	if len(errs) == 0 {
		t.Error("GROUPING SETS with a trailing item parsed, want error")
	}

	// A column that happens to be named grouping still groups normally.
	stmt = mustParseSelect(t, "SELECT grouping FROM t GROUP BY grouping")
	if _, ok := stmt.GroupBy[0].(*ast.ColumnRef); !ok {
		t.Errorf("GroupBy[0] = %T, want *ast.ColumnRef", stmt.GroupBy[0])
	}
}

func TestSelectFromTableValuedFunction(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM BACKENDS()")
	ref, ok := stmt.From[0].(*ast.TableRef)
	if !ok || ref.Func == nil {
		t.Fatalf("From[0] = %+v, want TableRef with Func", stmt.From[0])
	}
	if len(ref.Func.Name.Parts) != 1 || ref.Func.Name.Parts[0] != "BACKENDS" {
		t.Errorf("Func.Name = %+v, want BACKENDS", ref.Func.Name)
	}

	stmt = mustParseSelect(t, `SELECT * FROM numbers("number" = "10") n`)
	ref = stmt.From[0].(*ast.TableRef)
	if ref.Func == nil || len(ref.Func.Args) != 1 {
		t.Fatalf("Func = %+v, want one argument", ref.Func)
	}
	if ref.Alias != "n" {
		t.Errorf("Alias = %q, want n", ref.Alias)
	}
}

func TestSelectFromTabletAndTableSample(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM t1 TABLET(10001) TABLESAMPLE(1000 ROWS) REPEATABLE 2 LIMIT 1000")
	ref := stmt.From[0].(*ast.TableRef)
	if len(ref.TabletIDs) != 1 || ref.TabletIDs[0] != 10001 {
		t.Errorf("TabletIDs = %v, want [10001]", ref.TabletIDs)
	}
	if ref.Sample == nil || ref.Sample.Unit != "ROWS" {
		t.Fatalf("Sample = %+v, want 1000 ROWS", ref.Sample)
	}
	if lit, ok := ref.Sample.Value.(*ast.Literal); !ok || lit.Value != "1000" {
		t.Errorf("Sample.Value = %+v, want 1000", ref.Sample.Value)
	}
	if seed, ok := ref.Sample.Seed.(*ast.Literal); !ok || seed.Value != "2" {
		t.Errorf("Sample.Seed = %+v, want 2", ref.Sample.Seed)
	}
	if stmt.Limit == nil {
		t.Error("LIMIT after the sample clause was lost")
	}

	stmt = mustParseSelect(t, "SELECT * FROM t TABLESAMPLE(20 PERCENT)")
	if ref := stmt.From[0].(*ast.TableRef); ref.Sample == nil || ref.Sample.Unit != "PERCENT" {
		t.Errorf("Sample = %+v, want 20 PERCENT", ref.Sample)
	}

	stmt = mustParseSelect(t, "SELECT * FROM t TABLESAMPLE()")
	if ref := stmt.From[0].(*ast.TableRef); ref.Sample == nil || ref.Sample.Value != nil {
		t.Errorf("Sample = %+v, want empty TABLESAMPLE()", ref.Sample)
	}

	// TABLET sits before the alias, TABLESAMPLE after it.
	stmt = mustParseSelect(t, "SELECT * FROM t1 TABLET(1) x TABLESAMPLE(10 ROWS)")
	if ref := stmt.From[0].(*ast.TableRef); ref.Alias != "x" || len(ref.TabletIDs) != 1 || ref.Sample == nil {
		t.Errorf("ref = %+v, want tablet+alias+sample", stmt.From[0])
	}
}

func TestSelectTableFunctionNameWalkedOnce(t *testing.T) {
	// TableRef.Name mirrors Func.Name for a table-valued function (same
	// *ObjectName); the walker must visit it once, not twice.
	stmt := mustParseSelect(t, "SELECT * FROM BACKENDS()")
	count := 0
	ast.Inspect(stmt, func(n ast.Node) bool {
		if _, ok := n.(*ast.ObjectName); ok {
			count++
		}
		return true
	})
	if count != 1 {
		t.Errorf("ObjectName visited %d times, want 1", count)
	}
}

func TestSelectQualifiedColumnContinuations(t *testing.T) {
	// A select item starting with a qualified column must flow through the
	// general expression parser: the old fast path returned a bare ColumnRef
	// and handed any continuation to the trailing-token swallow.
	stmt := mustParseSelect(t, "SELECT t.secret_col[1] FROM sensitive_table t")
	if _, ok := stmt.Items[0].Expr.(*ast.ElementAtExpr); !ok {
		t.Fatalf("Items[0].Expr = %T, want *ast.ElementAtExpr", stmt.Items[0].Expr)
	}
	if len(stmt.From) != 1 {
		t.Fatal("FROM clause lost after qualified subscript")
	}

	stmt = mustParseSelect(t, "SELECT t.a + 1 FROM t")
	if _, ok := stmt.Items[0].Expr.(*ast.BinaryExpr); !ok {
		t.Fatalf("Items[0].Expr = %T, want *ast.BinaryExpr", stmt.Items[0].Expr)
	}
	if len(stmt.From) != 1 {
		t.Fatal("FROM clause lost after qualified arithmetic")
	}

	// The plain and star forms keep their shapes.
	stmt = mustParseSelect(t, "SELECT t.a AS x FROM t")
	if _, ok := stmt.Items[0].Expr.(*ast.ColumnRef); !ok || stmt.Items[0].Alias != "x" {
		t.Fatalf("Items[0] = %+v, want aliased ColumnRef", stmt.Items[0])
	}
	stmt = mustParseSelect(t, "SELECT t.* FROM t")
	if !stmt.Items[0].Star || stmt.Items[0].TableName == nil {
		t.Fatalf("Items[0] = %+v, want qualified star", stmt.Items[0])
	}
	stmt = mustParseSelect(t, "SELECT db.f(1) FROM t")
	if _, ok := stmt.Items[0].Expr.(*ast.FuncCallExpr); !ok {
		t.Fatalf("Items[0].Expr = %T, want *ast.FuncCallExpr", stmt.Items[0].Expr)
	}
}

func TestParenQueryStatements(t *testing.T) {
	// Top-level parenthesized queries are engine-valid; the parens are
	// grouping only, so the inner node comes back directly.
	file, errs := Parse("(SELECT 1)")
	if len(errs) != 0 {
		t.Fatalf("(SELECT 1) errors: %v", errs)
	}
	if _, ok := file.Stmts[0].(*ast.SelectStmt); !ok {
		t.Fatalf("stmt = %T, want *ast.SelectStmt", file.Stmts[0])
	}
	file, errs = Parse("((SELECT 1))")
	if len(errs) != 0 {
		t.Fatalf("((SELECT 1)) errors: %v", errs)
	}
	file, errs = Parse("(SELECT 1) UNION (SELECT 2)")
	if len(errs) != 0 {
		t.Fatalf("union errors: %v", errs)
	}
	if _, ok := file.Stmts[0].(*ast.SetOpStmt); !ok {
		t.Fatalf("stmt = %T, want *ast.SetOpStmt", file.Stmts[0])
	}
	if _, errs := Parse("(INSERT INTO t VALUES (1))"); len(errs) == 0 {
		t.Error("(INSERT ...) parsed, want error")
	}
}

func TestParenQuerySetOpTails(t *testing.T) {
	// A set-op tail follows any operand kind inside the parens
	// (engine-verified accepts).
	for _, sql := range []string{
		"((SELECT 1) UNION SELECT 2)",
		"((SELECT 1) UNION SELECT 2 UNION SELECT 3)",
		"((SELECT 1) INTERSECT SELECT 2)",
	} {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Fatalf("%s errors: %v", sql, errs)
		}
		if _, ok := file.Stmts[0].(*ast.SetOpStmt); !ok {
			t.Fatalf("%s: stmt = %T, want *ast.SetOpStmt", sql, file.Stmts[0])
		}
	}

	// A WITH-bearing group keeps its scope-boundary wrapper.
	file, errs := Parse("(WITH c AS (SELECT 1) SELECT 1 UNION SELECT 2)")
	if len(errs) != 0 {
		t.Fatalf("paren WITH union errors: %v", errs)
	}
	g, ok := file.Stmts[0].(*ast.GroupedQuery)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.GroupedQuery (CTE scope boundary)", file.Stmts[0])
	}
	if _, ok := g.Query.(*ast.SetOpStmt); !ok {
		t.Fatalf("Query = %T, want *ast.SetOpStmt", g.Query)
	}
}

func TestParenBoundsCTEScope(t *testing.T) {
	// (WITH c AS (...) SELECT 1) UNION SELECT * FROM c: the engine scopes c
	// to the parens ("Table [c] does not exist", container-verified), so the
	// left arm stays wrapped in the boundary node.
	file, errs := Parse("(WITH c AS (SELECT 1) SELECT 1) UNION SELECT * FROM c")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	setOp, ok := file.Stmts[0].(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.SetOpStmt", file.Stmts[0])
	}
	if _, ok := setOp.Left.(*ast.GroupedQuery); !ok {
		t.Fatalf("Left = %T, want *ast.GroupedQuery (scope boundary)", setOp.Left)
	}

	// A CTE body may itself be parenthesized (engine-verified).
	if _, errs := Parse("WITH c AS ((SELECT 1)) SELECT * FROM c"); len(errs) != 0 {
		t.Fatalf("paren CTE body errors: %v", errs)
	}
}

func TestParenQueryTrailingClauses(t *testing.T) {
	// Outer ORDER BY / LIMIT attach to the grouped query rather than being
	// consumed and dropped.
	file, errs := Parse("(SELECT 1) ORDER BY 1 LIMIT 5")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	sel, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.SelectStmt", file.Stmts[0])
	}
	if len(sel.OrderBy) != 1 || sel.Limit == nil {
		t.Fatalf("clauses dropped: OrderBy=%d Limit=%v", len(sel.OrderBy), sel.Limit)
	}

	file, errs = Parse("((SELECT 1)) LIMIT 5")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if sel := file.Stmts[0].(*ast.SelectStmt); sel.Limit == nil {
		t.Fatal("LIMIT dropped on nested paren query")
	}

	file, errs = Parse("(SELECT 1) UNION (SELECT 2) ORDER BY 1 LIMIT 5")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	setOp, ok := file.Stmts[0].(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.SetOpStmt", file.Stmts[0])
	}
	if len(setOp.OrderBy) != 1 || setOp.Limit == nil {
		t.Fatalf("clauses dropped: OrderBy=%d Limit=%v", len(setOp.OrderBy), setOp.Limit)
	}

	// The strict trailing-token check still applies after the clauses.
	if _, errs := Parse("(SELECT 1) ORDER BY 1 LIMIT 5 x"); len(errs) == 0 {
		t.Error("junk after trailing clauses parsed, want error")
	}
}

func TestQueryTailStatementLevel(t *testing.T) {
	// Trailing clauses after a set operation whose right operand is
	// parenthesized (engine-verified accepts).
	file, errs := Parse("SELECT 1 UNION (SELECT 2) LIMIT 5")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	setOp := file.Stmts[0].(*ast.SetOpStmt)
	if setOp.Limit == nil {
		t.Fatal("LIMIT not attached to statement-level set operation")
	}

	file, errs = Parse("SELECT 1 UNION (SELECT 2) ORDER BY 1 LIMIT 5")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	setOp = file.Stmts[0].(*ast.SetOpStmt)
	if len(setOp.OrderBy) != 1 || setOp.Limit == nil {
		t.Fatalf("clauses dropped: OrderBy=%d Limit=%v", len(setOp.OrderBy), setOp.Limit)
	}

	// WITH ... SELECT continues with set operators at statement level.
	file, errs = Parse("WITH c AS (SELECT 1) SELECT 1 UNION SELECT 2")
	if len(errs) != 0 {
		t.Fatalf("WITH union errors: %v", errs)
	}
	setOp, ok := file.Stmts[0].(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.SetOpStmt", file.Stmts[0])
	}
	left, ok := setOp.Left.(*ast.SelectStmt)
	if !ok || left.With == nil {
		t.Fatalf("WITH clause lost from left arm (%T, With=%v)", setOp.Left, ok && left.With != nil)
	}

	// A CTE body takes a set-op tail too.
	file, errs = Parse("WITH c AS (SELECT 1 UNION SELECT 2) SELECT * FROM c")
	if len(errs) != 0 {
		t.Fatalf("CTE body union errors: %v", errs)
	}
	sel := file.Stmts[0].(*ast.SelectStmt)
	if _, ok := sel.With.CTEs[0].Query.(*ast.SetOpStmt); !ok {
		t.Fatalf("CTE query = %T, want *ast.SetOpStmt", sel.With.CTEs[0].Query)
	}

	// The engine is lenient about clause order and repetition
	// (container-verified). A second group lands on a wrapper so the two
	// stay ordered: LIMIT-first-then-ORDER is not ORDER-then-LIMIT.
	file, errs = Parse("SELECT 1 LIMIT 5 ORDER BY 1")
	if len(errs) != 0 {
		t.Fatalf("LIMIT-then-ORDER errors: %v", errs)
	}
	g, ok := file.Stmts[0].(*ast.GroupedQuery)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.GroupedQuery", file.Stmts[0])
	}
	if len(g.OrderBy) != 1 || g.Query.(*ast.SelectStmt).Limit == nil {
		t.Fatalf("clause layers wrong: outer OrderBy=%d inner Limit=%v", len(g.OrderBy), g.Query.(*ast.SelectStmt).Limit)
	}
	if _, errs := Parse("SELECT 1 ORDER BY 1 ORDER BY 2"); len(errs) != 0 {
		t.Fatalf("repeated ORDER BY errors: %v", errs)
	}

	// Junk after the clauses still errors.
	if _, errs := Parse("SELECT 1 UNION (SELECT 2) GARBAGE"); len(errs) == 0 {
		t.Error("junk after set operation parsed, want error")
	}
}

func TestGroupedQueryWrapOnConflict(t *testing.T) {
	// A trailing group that repeats a clause the query already carries wraps
	// in GroupedQuery instead of overwriting — both layers stay in the tree.
	file, errs := Parse("(SELECT 1 ORDER BY 1) ORDER BY 2")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	g, ok := file.Stmts[0].(*ast.GroupedQuery)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.GroupedQuery", file.Stmts[0])
	}
	if len(g.OrderBy) != 1 {
		t.Fatalf("outer OrderBy = %d, want 1", len(g.OrderBy))
	}
	inner, ok := g.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", g.Query)
	}
	if len(inner.OrderBy) != 1 {
		t.Fatal("inner ORDER BY overwritten by the outer group")
	}

	file, errs = Parse("SELECT 1 LIMIT 1 LIMIT 2")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	g, ok = file.Stmts[0].(*ast.GroupedQuery)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.GroupedQuery", file.Stmts[0])
	}
	if g.Limit == nil || g.Query.(*ast.SelectStmt).Limit == nil {
		t.Fatal("one of the LIMIT layers was lost")
	}

	// Different clause kinds also wrap: (SELECT a FROM t LIMIT 1) ORDER BY a
	// limits first and then orders, which is not the same operation as
	// SELECT a FROM t ORDER BY a LIMIT 1.
	file, errs = Parse("(SELECT a FROM t LIMIT 1) ORDER BY a")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	g, ok = file.Stmts[0].(*ast.GroupedQuery)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.GroupedQuery", file.Stmts[0])
	}
	if len(g.OrderBy) != 1 || g.Query.(*ast.SelectStmt).Limit == nil {
		t.Fatal("cross-kind clause layers lost")
	}

	// No clause on the node — no wrapper: the group attaches in place.
	file, errs = Parse("(SELECT 1) ORDER BY 1")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if _, ok := file.Stmts[0].(*ast.SelectStmt); !ok {
		t.Fatalf("stmt = %T, want *ast.SelectStmt (no wrapper without conflict)", file.Stmts[0])
	}

	// Both clause layers are reachable from Walk.
	file, errs = Parse("SELECT 1 ORDER BY a ORDER BY b")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	count := 0
	ast.Inspect(file.Stmts[0], func(n ast.Node) bool {
		if n != nil && n.Tag() == ast.T_ColumnRef {
			count++
		}
		return true
	})
	if count != 2 {
		t.Errorf("ColumnRef visited %d times, want 2 (both ORDER BY layers)", count)
	}
}

func TestParenQueryLocCoversParens(t *testing.T) {
	// The parens are grouping only, but NodeLoc must cover the delimiters
	// and any trailing clauses, or source extraction drops them.
	for _, sql := range []string{
		"(SELECT 1)",
		"((SELECT 1))",
		"(SELECT 1) ORDER BY 1 LIMIT 5",
		"((SELECT 1) UNION SELECT 2)",
		"SELECT 1 UNION (SELECT 2) LIMIT 5",
		"(SELECT 1 ORDER BY 1) ORDER BY 2",
		"(SELECT 1 LIMIT 1) ORDER BY 2",
		"(WITH c AS (SELECT 1) SELECT 1)",
	} {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Fatalf("%s errors: %v", sql, errs)
		}
		loc := ast.NodeLoc(file.Stmts[0])
		if loc.Start != 0 || loc.End != len(sql) {
			t.Errorf("%s: NodeLoc = [%d,%d), want [0,%d)", sql, loc.Start, loc.End, len(sql))
		}
	}
}
