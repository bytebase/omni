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
