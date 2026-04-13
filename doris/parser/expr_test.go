package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// parseExprFrom creates a Parser from input and calls parseExpr.
func parseExprFrom(input string) (ast.Node, error) {
	p := makeParser(input)
	return p.parseExpr()
}

// mustParseExpr is a test helper that calls parseExprFrom and fatals on error.
func mustParseExpr(t *testing.T, input string) ast.Node {
	t.Helper()
	node, err := parseExprFrom(input)
	if err != nil {
		t.Fatalf("parseExpr(%q) error: %v", input, err)
	}
	if node == nil {
		t.Fatalf("parseExpr(%q) returned nil", input)
	}
	return node
}

// ---------------------------------------------------------------------------
// Literals
// ---------------------------------------------------------------------------

func TestExprLiteralInt(t *testing.T) {
	node := mustParseExpr(t, "42")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitInt {
		t.Errorf("kind = %d, want LitInt", lit.Kind)
	}
	if lit.Value != "42" {
		t.Errorf("value = %q, want %q", lit.Value, "42")
	}
}

func TestExprLiteralFloat(t *testing.T) {
	node := mustParseExpr(t, "3.14")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitFloat {
		t.Errorf("kind = %d, want LitFloat", lit.Kind)
	}
	if lit.Value != "3.14" {
		t.Errorf("value = %q, want %q", lit.Value, "3.14")
	}
}

func TestExprLiteralString(t *testing.T) {
	node := mustParseExpr(t, "'hello'")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitString {
		t.Errorf("kind = %d, want LitString", lit.Kind)
	}
	if lit.Value != "hello" {
		t.Errorf("value = %q, want %q", lit.Value, "hello")
	}
}

func TestExprLiteralTrue(t *testing.T) {
	node := mustParseExpr(t, "TRUE")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitBool {
		t.Errorf("kind = %d, want LitBool", lit.Kind)
	}
}

func TestExprLiteralFalse(t *testing.T) {
	node := mustParseExpr(t, "FALSE")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitBool {
		t.Errorf("kind = %d, want LitBool", lit.Kind)
	}
}

func TestExprLiteralNull(t *testing.T) {
	node := mustParseExpr(t, "NULL")
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitNull {
		t.Errorf("kind = %d, want LitNull", lit.Kind)
	}
}

// ---------------------------------------------------------------------------
// Column references
// ---------------------------------------------------------------------------

func TestExprColumnRefSimple(t *testing.T) {
	node := mustParseExpr(t, "col")
	ref, ok := node.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", node)
	}
	if len(ref.Name.Parts) != 1 || ref.Name.Parts[0] != "col" {
		t.Errorf("parts = %v, want [col]", ref.Name.Parts)
	}
}

func TestExprColumnRefQualified(t *testing.T) {
	node := mustParseExpr(t, "t.col")
	ref, ok := node.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", node)
	}
	if len(ref.Name.Parts) != 2 || ref.Name.Parts[0] != "t" || ref.Name.Parts[1] != "col" {
		t.Errorf("parts = %v, want [t, col]", ref.Name.Parts)
	}
}

func TestExprColumnRefThreePart(t *testing.T) {
	node := mustParseExpr(t, "db.t.col")
	ref, ok := node.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", node)
	}
	if len(ref.Name.Parts) != 3 {
		t.Errorf("parts = %v, want [db, t, col]", ref.Name.Parts)
	}
}

// ---------------------------------------------------------------------------
// Arithmetic
// ---------------------------------------------------------------------------

func TestExprArithmeticSimple(t *testing.T) {
	node := mustParseExpr(t, "1 + 2")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAdd {
		t.Errorf("op = %v, want BinAdd", bin.Op)
	}
}

func TestExprArithmeticPrecedence(t *testing.T) {
	// a * b + c should parse as (a * b) + c
	node := mustParseExpr(t, "a * b + c")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAdd {
		t.Errorf("top op = %v, want BinAdd", bin.Op)
	}
	leftBin, ok := bin.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.BinaryExpr, got %T", bin.Left)
	}
	if leftBin.Op != ast.BinMul {
		t.Errorf("left op = %v, want BinMul", leftBin.Op)
	}
}

func TestExprArithmeticParens(t *testing.T) {
	// (a + b) * c
	node := mustParseExpr(t, "(a + b) * c")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinMul {
		t.Errorf("top op = %v, want BinMul", bin.Op)
	}
	paren, ok := bin.Left.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.ParenExpr, got %T", bin.Left)
	}
	innerBin, ok := paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected paren inner to be *ast.BinaryExpr, got %T", paren.Expr)
	}
	if innerBin.Op != ast.BinAdd {
		t.Errorf("inner op = %v, want BinAdd", innerBin.Op)
	}
}

func TestExprDivMod(t *testing.T) {
	node := mustParseExpr(t, "10 / 3")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinDiv {
		t.Errorf("op = %v, want BinDiv", bin.Op)
	}

	node2 := mustParseExpr(t, "10 % 3")
	bin2, ok := node2.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node2)
	}
	if bin2.Op != ast.BinMod {
		t.Errorf("op = %v, want BinMod", bin2.Op)
	}
}

func TestExprDIVKeyword(t *testing.T) {
	node := mustParseExpr(t, "10 DIV 3")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinIntDiv {
		t.Errorf("op = %v, want BinIntDiv", bin.Op)
	}
}

// ---------------------------------------------------------------------------
// Comparison
// ---------------------------------------------------------------------------

func TestExprComparisonEq(t *testing.T) {
	node := mustParseExpr(t, "a = 1")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinEq {
		t.Errorf("op = %v, want BinEq", bin.Op)
	}
}

func TestExprComparisonNe(t *testing.T) {
	node := mustParseExpr(t, "a <> b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinNe {
		t.Errorf("op = %v, want BinNe", bin.Op)
	}
}

func TestExprComparisonNeBang(t *testing.T) {
	node := mustParseExpr(t, "a != b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinNe {
		t.Errorf("op = %v, want BinNe", bin.Op)
	}
}

func TestExprComparisonNullSafe(t *testing.T) {
	node := mustParseExpr(t, "a <=> b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinNullSafeEq {
		t.Errorf("op = %v, want BinNullSafeEq", bin.Op)
	}
}

func TestExprComparisonGe(t *testing.T) {
	node := mustParseExpr(t, "a >= 1")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinGe {
		t.Errorf("op = %v, want BinGe", bin.Op)
	}
}

func TestExprComparisonLt(t *testing.T) {
	node := mustParseExpr(t, "a < 1")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinLt {
		t.Errorf("op = %v, want BinLt", bin.Op)
	}
}

func TestExprComparisonGt(t *testing.T) {
	node := mustParseExpr(t, "a > 1")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinGt {
		t.Errorf("op = %v, want BinGt", bin.Op)
	}
}

func TestExprComparisonLe(t *testing.T) {
	node := mustParseExpr(t, "a <= 1")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinLe {
		t.Errorf("op = %v, want BinLe", bin.Op)
	}
}

// ---------------------------------------------------------------------------
// Logical operators
// ---------------------------------------------------------------------------

func TestExprLogicalAnd(t *testing.T) {
	node := mustParseExpr(t, "a AND b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAnd {
		t.Errorf("op = %v, want BinAnd", bin.Op)
	}
}

func TestExprLogicalOr(t *testing.T) {
	node := mustParseExpr(t, "a OR b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinOr {
		t.Errorf("op = %v, want BinOr", bin.Op)
	}
}

func TestExprLogicalXor(t *testing.T) {
	node := mustParseExpr(t, "a XOR b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinXor {
		t.Errorf("op = %v, want BinXor", bin.Op)
	}
}

func TestExprLogicalNot(t *testing.T) {
	node := mustParseExpr(t, "NOT a")
	un, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if un.Op != ast.UnaryNot {
		t.Errorf("op = %v, want UnaryNot", un.Op)
	}
}

func TestExprLogicalPrecedence(t *testing.T) {
	// a AND b OR c should parse as (a AND b) OR c
	node := mustParseExpr(t, "a AND b OR c")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinOr {
		t.Errorf("top op = %v, want BinOr", bin.Op)
	}
	leftBin, ok := bin.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.BinaryExpr, got %T", bin.Left)
	}
	if leftBin.Op != ast.BinAnd {
		t.Errorf("left op = %v, want BinAnd", leftBin.Op)
	}
}

func TestExprLogicalAndSymbol(t *testing.T) {
	node := mustParseExpr(t, "a && b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAnd {
		t.Errorf("op = %v, want BinAnd", bin.Op)
	}
}

func TestExprLogicalOrSymbol(t *testing.T) {
	node := mustParseExpr(t, "a || b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinOr {
		t.Errorf("op = %v, want BinOr", bin.Op)
	}
}

// ---------------------------------------------------------------------------
// IS expressions
// ---------------------------------------------------------------------------

func TestExprIsNull(t *testing.T) {
	node := mustParseExpr(t, "a IS NULL")
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if is.Not {
		t.Error("expected Not=false")
	}
	if is.IsWhat != "NULL" {
		t.Errorf("IsWhat = %q, want %q", is.IsWhat, "NULL")
	}
}

func TestExprIsNotNull(t *testing.T) {
	node := mustParseExpr(t, "a IS NOT NULL")
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if !is.Not {
		t.Error("expected Not=true")
	}
	if is.IsWhat != "NULL" {
		t.Errorf("IsWhat = %q, want %q", is.IsWhat, "NULL")
	}
}

func TestExprIsTrue(t *testing.T) {
	node := mustParseExpr(t, "a IS TRUE")
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if is.IsWhat != "TRUE" {
		t.Errorf("IsWhat = %q, want %q", is.IsWhat, "TRUE")
	}
}

func TestExprIsNotFalse(t *testing.T) {
	node := mustParseExpr(t, "a IS NOT FALSE")
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if !is.Not {
		t.Error("expected Not=true")
	}
	if is.IsWhat != "FALSE" {
		t.Errorf("IsWhat = %q, want %q", is.IsWhat, "FALSE")
	}
}

// ---------------------------------------------------------------------------
// BETWEEN
// ---------------------------------------------------------------------------

func TestExprBetween(t *testing.T) {
	node := mustParseExpr(t, "a BETWEEN 1 AND 10")
	btw, ok := node.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected *ast.BetweenExpr, got %T", node)
	}
	if btw.Not {
		t.Error("expected Not=false")
	}

	// Verify low and high are literals
	lowLit, ok := btw.Low.(*ast.Literal)
	if !ok {
		t.Fatalf("expected Low to be *ast.Literal, got %T", btw.Low)
	}
	if lowLit.Value != "1" {
		t.Errorf("low = %q, want %q", lowLit.Value, "1")
	}

	highLit, ok := btw.High.(*ast.Literal)
	if !ok {
		t.Fatalf("expected High to be *ast.Literal, got %T", btw.High)
	}
	if highLit.Value != "10" {
		t.Errorf("high = %q, want %q", highLit.Value, "10")
	}
}

func TestExprNotBetween(t *testing.T) {
	node := mustParseExpr(t, "a NOT BETWEEN 1 AND 10")
	btw, ok := node.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected *ast.BetweenExpr, got %T", node)
	}
	if !btw.Not {
		t.Error("expected Not=true")
	}
}

// ---------------------------------------------------------------------------
// IN
// ---------------------------------------------------------------------------

func TestExprIn(t *testing.T) {
	node := mustParseExpr(t, "a IN (1, 2, 3)")
	in, ok := node.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", node)
	}
	if in.Not {
		t.Error("expected Not=false")
	}
	if len(in.Values) != 3 {
		t.Errorf("len(Values) = %d, want 3", len(in.Values))
	}
}

func TestExprNotIn(t *testing.T) {
	node := mustParseExpr(t, "a NOT IN (1, 2)")
	in, ok := node.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", node)
	}
	if !in.Not {
		t.Error("expected Not=true")
	}
	if len(in.Values) != 2 {
		t.Errorf("len(Values) = %d, want 2", len(in.Values))
	}
}

func TestExprInSubquery(t *testing.T) {
	node := mustParseExpr(t, "a IN (SELECT id FROM t)")
	in, ok := node.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", node)
	}
	if len(in.Values) != 1 {
		t.Fatalf("len(Values) = %d, want 1", len(in.Values))
	}
	subq, ok := in.Values[0].(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr, got %T", in.Values[0])
	}
	if subq.RawText != "SELECT id FROM t" {
		t.Errorf("RawText = %q, want %q", subq.RawText, "SELECT id FROM t")
	}
}

// ---------------------------------------------------------------------------
// LIKE
// ---------------------------------------------------------------------------

func TestExprLike(t *testing.T) {
	node := mustParseExpr(t, "a LIKE '%foo%'")
	like, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if like.Not {
		t.Error("expected Not=false")
	}
	pat, ok := like.Pattern.(*ast.Literal)
	if !ok {
		t.Fatalf("expected Pattern to be *ast.Literal, got %T", like.Pattern)
	}
	if pat.Value != "%foo%" {
		t.Errorf("pattern = %q, want %q", pat.Value, "%foo%")
	}
}

func TestExprNotLike(t *testing.T) {
	node := mustParseExpr(t, "a NOT LIKE 'bar'")
	like, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if !like.Not {
		t.Error("expected Not=true")
	}
}

func TestExprLikeEscape(t *testing.T) {
	node := mustParseExpr(t, "a NOT LIKE 'bar' ESCAPE '\\\\'")
	like, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if like.Escape == nil {
		t.Fatal("expected non-nil Escape")
	}
}

// ---------------------------------------------------------------------------
// REGEXP / RLIKE
// ---------------------------------------------------------------------------

func TestExprRegexp(t *testing.T) {
	node := mustParseExpr(t, "a REGEXP '^[0-9]+'")
	re, ok := node.(*ast.RegexpExpr)
	if !ok {
		t.Fatalf("expected *ast.RegexpExpr, got %T", node)
	}
	if re.Not {
		t.Error("expected Not=false")
	}
}

func TestExprRLike(t *testing.T) {
	node := mustParseExpr(t, "a RLIKE '^abc'")
	re, ok := node.(*ast.RegexpExpr)
	if !ok {
		t.Fatalf("expected *ast.RegexpExpr, got %T", node)
	}
	if re.Not {
		t.Error("expected Not=false")
	}
}

func TestExprNotRegexp(t *testing.T) {
	node := mustParseExpr(t, "a NOT REGEXP '^[0-9]+'")
	re, ok := node.(*ast.RegexpExpr)
	if !ok {
		t.Fatalf("expected *ast.RegexpExpr, got %T", node)
	}
	if !re.Not {
		t.Error("expected Not=true")
	}
}

// ---------------------------------------------------------------------------
// CAST / TRY_CAST
// ---------------------------------------------------------------------------

func TestExprCast(t *testing.T) {
	node := mustParseExpr(t, "CAST(a AS INT)")
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if cast.TryCast {
		t.Error("expected TryCast=false")
	}
	if cast.TypeName.Name != "INT" {
		t.Errorf("type = %q, want %q", cast.TypeName.Name, "INT")
	}
}

func TestExprTryCast(t *testing.T) {
	node := mustParseExpr(t, "TRY_CAST(a AS VARCHAR(255))")
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if !cast.TryCast {
		t.Error("expected TryCast=true")
	}
	if cast.TypeName.Name != "VARCHAR" {
		t.Errorf("type = %q, want %q", cast.TypeName.Name, "VARCHAR")
	}
	if len(cast.TypeName.Params) != 1 || cast.TypeName.Params[0] != 255 {
		t.Errorf("params = %v, want [255]", cast.TypeName.Params)
	}
}

func TestExprCastComplex(t *testing.T) {
	node := mustParseExpr(t, "CAST(a + b AS DECIMAL(10,2))")
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	// Inner expression should be a + b
	bin, ok := cast.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected inner to be *ast.BinaryExpr, got %T", cast.Expr)
	}
	if bin.Op != ast.BinAdd {
		t.Errorf("inner op = %v, want BinAdd", bin.Op)
	}
	if cast.TypeName.Name != "DECIMAL" {
		t.Errorf("type = %q, want %q", cast.TypeName.Name, "DECIMAL")
	}
	if len(cast.TypeName.Params) != 2 || cast.TypeName.Params[0] != 10 || cast.TypeName.Params[1] != 2 {
		t.Errorf("params = %v, want [10,2]", cast.TypeName.Params)
	}
}

// ---------------------------------------------------------------------------
// CASE expressions
// ---------------------------------------------------------------------------

func TestExprCaseSearched(t *testing.T) {
	node := mustParseExpr(t, "CASE WHEN a > 1 THEN 'yes' ELSE 'no' END")
	ce, ok := node.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected *ast.CaseExpr, got %T", node)
	}
	if ce.Kind != ast.CaseSearched {
		t.Errorf("kind = %d, want CaseSearched", ce.Kind)
	}
	if ce.Operand != nil {
		t.Error("expected nil Operand for searched CASE")
	}
	if len(ce.Whens) != 1 {
		t.Errorf("len(Whens) = %d, want 1", len(ce.Whens))
	}
	if ce.Else == nil {
		t.Error("expected non-nil Else")
	}
}

func TestExprCaseSimple(t *testing.T) {
	node := mustParseExpr(t, "CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' END")
	ce, ok := node.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected *ast.CaseExpr, got %T", node)
	}
	if ce.Kind != ast.CaseSimple {
		t.Errorf("kind = %d, want CaseSimple", ce.Kind)
	}
	if ce.Operand == nil {
		t.Fatal("expected non-nil Operand for simple CASE")
	}
	if len(ce.Whens) != 2 {
		t.Errorf("len(Whens) = %d, want 2", len(ce.Whens))
	}
	if ce.Else != nil {
		t.Error("expected nil Else")
	}
}

// ---------------------------------------------------------------------------
// Function calls
// ---------------------------------------------------------------------------

func TestExprFuncCountStar(t *testing.T) {
	node := mustParseExpr(t, "COUNT(*)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if !fc.Star {
		t.Error("expected Star=true")
	}
}

func TestExprFuncSum(t *testing.T) {
	node := mustParseExpr(t, "SUM(a)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 1 {
		t.Errorf("len(Args) = %d, want 1", len(fc.Args))
	}
}

func TestExprFuncCoalesce(t *testing.T) {
	node := mustParseExpr(t, "COALESCE(a, b, c)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 3 {
		t.Errorf("len(Args) = %d, want 3", len(fc.Args))
	}
}

func TestExprFuncCountDistinct(t *testing.T) {
	node := mustParseExpr(t, "COUNT(DISTINCT a)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if !fc.Distinct {
		t.Error("expected Distinct=true")
	}
	if len(fc.Args) != 1 {
		t.Errorf("len(Args) = %d, want 1", len(fc.Args))
	}
}

func TestExprFuncNoArgs(t *testing.T) {
	node := mustParseExpr(t, "NOW()")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 0 {
		t.Errorf("len(Args) = %d, want 0", len(fc.Args))
	}
}

// ---------------------------------------------------------------------------
// Nilary functions
// ---------------------------------------------------------------------------

func TestExprCurrentDate(t *testing.T) {
	node := mustParseExpr(t, "CURRENT_DATE")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Name.Parts[0] != "CURRENT_DATE" {
		t.Errorf("name = %q, want %q", fc.Name.Parts[0], "CURRENT_DATE")
	}
}

func TestExprCurrentTimestamp(t *testing.T) {
	node := mustParseExpr(t, "CURRENT_TIMESTAMP")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Name.Parts[0] != "CURRENT_TIMESTAMP" {
		t.Errorf("name = %q, want %q", fc.Name.Parts[0], "CURRENT_TIMESTAMP")
	}
}

// ---------------------------------------------------------------------------
// Unary operators
// ---------------------------------------------------------------------------

func TestExprUnaryMinus(t *testing.T) {
	node := mustParseExpr(t, "-a")
	un, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if un.Op != ast.UnaryMinus {
		t.Errorf("op = %v, want UnaryMinus", un.Op)
	}
}

func TestExprUnaryPlus(t *testing.T) {
	node := mustParseExpr(t, "+a")
	un, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if un.Op != ast.UnaryPlus {
		t.Errorf("op = %v, want UnaryPlus", un.Op)
	}
}

func TestExprUnaryBitNot(t *testing.T) {
	node := mustParseExpr(t, "~a")
	un, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if un.Op != ast.UnaryBitNot {
		t.Errorf("op = %v, want UnaryBitNot", un.Op)
	}
}

// ---------------------------------------------------------------------------
// Bitwise operators
// ---------------------------------------------------------------------------

func TestExprBitwiseOr(t *testing.T) {
	node := mustParseExpr(t, "a | b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinBitOr {
		t.Errorf("op = %v, want BinBitOr", bin.Op)
	}
}

func TestExprBitwiseAnd(t *testing.T) {
	node := mustParseExpr(t, "a & b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinBitAnd {
		t.Errorf("op = %v, want BinBitAnd", bin.Op)
	}
}

func TestExprBitwiseXor(t *testing.T) {
	node := mustParseExpr(t, "a ^ b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinBitXor {
		t.Errorf("op = %v, want BinBitXor", bin.Op)
	}
}

func TestExprShiftLeft(t *testing.T) {
	node := mustParseExpr(t, "a << 2")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinShiftLeft {
		t.Errorf("op = %v, want BinShiftLeft", bin.Op)
	}
}

func TestExprShiftRight(t *testing.T) {
	node := mustParseExpr(t, "a >> 2")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinShiftRight {
		t.Errorf("op = %v, want BinShiftRight", bin.Op)
	}
}

// ---------------------------------------------------------------------------
// Nested / complex expressions
// ---------------------------------------------------------------------------

func TestExprNestedComplex(t *testing.T) {
	// a > 0 AND (b IN (1,2) OR c IS NOT NULL)
	node := mustParseExpr(t, "a > 0 AND (b IN (1,2) OR c IS NOT NULL)")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAnd {
		t.Errorf("top op = %v, want BinAnd", bin.Op)
	}

	// Left should be a > 0
	leftBin, ok := bin.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.BinaryExpr, got %T", bin.Left)
	}
	if leftBin.Op != ast.BinGt {
		t.Errorf("left op = %v, want BinGt", leftBin.Op)
	}

	// Right should be a ParenExpr containing an OR
	paren, ok := bin.Right.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("expected right to be *ast.ParenExpr, got %T", bin.Right)
	}
	orBin, ok := paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected paren inner to be *ast.BinaryExpr, got %T", paren.Expr)
	}
	if orBin.Op != ast.BinOr {
		t.Errorf("paren inner op = %v, want BinOr", orBin.Op)
	}
}

func TestExprNotWithPrecedence(t *testing.T) {
	// NOT a OR b should parse as (NOT a) OR b
	node := mustParseExpr(t, "NOT a OR b")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinOr {
		t.Errorf("top op = %v, want BinOr", bin.Op)
	}
	un, ok := bin.Left.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.UnaryExpr, got %T", bin.Left)
	}
	if un.Op != ast.UnaryNot {
		t.Errorf("left op = %v, want UnaryNot", un.Op)
	}
}

// ---------------------------------------------------------------------------
// INTERVAL
// ---------------------------------------------------------------------------

func TestExprInterval(t *testing.T) {
	node := mustParseExpr(t, "INTERVAL 1 DAY")
	iv, ok := node.(*ast.IntervalExpr)
	if !ok {
		t.Fatalf("expected *ast.IntervalExpr, got %T", node)
	}
	if iv.Unit != "DAY" {
		t.Errorf("unit = %q, want %q", iv.Unit, "DAY")
	}
}

// ---------------------------------------------------------------------------
// EXISTS
// ---------------------------------------------------------------------------

func TestExprExists(t *testing.T) {
	node := mustParseExpr(t, "EXISTS (SELECT 1 FROM t)")
	ex, ok := node.(*ast.ExistsExpr)
	if !ok {
		t.Fatalf("expected *ast.ExistsExpr, got %T", node)
	}
	if ex.Subquery == nil {
		t.Fatal("expected non-nil Subquery")
	}
}

// ---------------------------------------------------------------------------
// Parenthesized expressions
// ---------------------------------------------------------------------------

func TestExprParen(t *testing.T) {
	node := mustParseExpr(t, "(42)")
	paren, ok := node.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("expected *ast.ParenExpr, got %T", node)
	}
	lit, ok := paren.Expr.(*ast.Literal)
	if !ok {
		t.Fatalf("expected inner to be *ast.Literal, got %T", paren.Expr)
	}
	if lit.Value != "42" {
		t.Errorf("value = %q, want %q", lit.Value, "42")
	}
}

// ---------------------------------------------------------------------------
// Subquery in expression position
// ---------------------------------------------------------------------------

func TestExprSubquery(t *testing.T) {
	node := mustParseExpr(t, "(SELECT 1)")
	subq, ok := node.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr, got %T", node)
	}
	if subq.RawText != "SELECT 1" {
		t.Errorf("RawText = %q, want %q", subq.RawText, "SELECT 1")
	}
}

// ---------------------------------------------------------------------------
// Location tracking
// ---------------------------------------------------------------------------

func TestExprLoc(t *testing.T) {
	// Verify that the Loc spans the full expression
	node := mustParseExpr(t, "1 + 2")
	loc := ast.NodeLoc(node)
	if loc.Start != 0 || loc.End != 5 {
		t.Errorf("loc = %v, want {0, 5}", loc)
	}
}

// ---------------------------------------------------------------------------
// Walk integration
// ---------------------------------------------------------------------------

func TestExprWalk(t *testing.T) {
	node := mustParseExpr(t, "a + b")
	var tags []ast.NodeTag
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			tags = append(tags, n.Tag())
		}
		return true
	})
	// Expected: BinaryExpr, ColumnRef, ObjectName, ColumnRef, ObjectName
	if len(tags) < 3 {
		t.Errorf("expected at least 3 visited nodes, got %d: %v", len(tags), tags)
	}
	if tags[0] != ast.T_BinaryExpr {
		t.Errorf("first tag = %v, want BinaryExpr", tags[0])
	}
}

// ---------------------------------------------------------------------------
// Edge cases and error handling
// ---------------------------------------------------------------------------

func TestExprError_Empty(t *testing.T) {
	_, err := parseExprFrom("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestExprError_UnterminatedParen(t *testing.T) {
	_, err := parseExprFrom("(1 + 2")
	if err == nil {
		t.Error("expected error for unterminated paren")
	}
}

func TestExprError_MissingOperand(t *testing.T) {
	_, err := parseExprFrom("1 +")
	if err == nil {
		t.Error("expected error for missing operand")
	}
}

// ---------------------------------------------------------------------------
// Comprehensive precedence tests
// ---------------------------------------------------------------------------

func TestExprPrecedence_BitwiseVsComparison(t *testing.T) {
	// a | b = c should parse as (a | b) = c since bitwise | > comparison
	node := mustParseExpr(t, "a | b = c")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinEq {
		t.Errorf("top op = %v, want BinEq", bin.Op)
	}
	leftBin, ok := bin.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.BinaryExpr, got %T", bin.Left)
	}
	if leftBin.Op != ast.BinBitOr {
		t.Errorf("left op = %v, want BinBitOr", leftBin.Op)
	}
}

func TestExprPrecedence_AddVsMul(t *testing.T) {
	// 1 + 2 * 3 should parse as 1 + (2 * 3)
	node := mustParseExpr(t, "1 + 2 * 3")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAdd {
		t.Errorf("top op = %v, want BinAdd", bin.Op)
	}
	rightBin, ok := bin.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be *ast.BinaryExpr, got %T", bin.Right)
	}
	if rightBin.Op != ast.BinMul {
		t.Errorf("right op = %v, want BinMul", rightBin.Op)
	}
}

func TestExprPrecedence_XorVsAnd(t *testing.T) {
	// a XOR b AND c should parse as a XOR (b AND c) since AND > XOR
	node := mustParseExpr(t, "a XOR b AND c")
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinXor {
		t.Errorf("top op = %v, want BinXor", bin.Op)
	}
	rightBin, ok := bin.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be *ast.BinaryExpr, got %T", bin.Right)
	}
	if rightBin.Op != ast.BinAnd {
		t.Errorf("right op = %v, want BinAnd", rightBin.Op)
	}
}

// ---------------------------------------------------------------------------
// Multiple expression types in one test (table-driven)
// ---------------------------------------------------------------------------

func TestExprParseSuccess(t *testing.T) {
	// These should all parse successfully without errors.
	tests := []struct {
		name  string
		input string
	}{
		{"int literal", "42"},
		{"float literal", "3.14"},
		{"string literal", "'hello world'"},
		{"true", "TRUE"},
		{"false", "FALSE"},
		{"null", "NULL"},
		{"column ref", "col"},
		{"qualified column", "t.col"},
		{"three-part column", "db.t.col"},
		{"addition", "1 + 2"},
		{"subtraction", "3 - 1"},
		{"multiplication", "2 * 3"},
		{"division", "10 / 2"},
		{"modulo", "10 % 3"},
		{"integer division", "10 DIV 3"},
		{"comparison eq", "a = 1"},
		{"comparison ne", "a <> b"},
		{"comparison ne bang", "a != b"},
		{"comparison lt", "a < 1"},
		{"comparison gt", "a > 1"},
		{"comparison le", "a <= 1"},
		{"comparison ge", "a >= 1"},
		{"null-safe eq", "a <=> b"},
		{"logical and", "a AND b"},
		{"logical or", "a OR b"},
		{"logical xor", "a XOR b"},
		{"logical not", "NOT a"},
		{"is null", "a IS NULL"},
		{"is not null", "a IS NOT NULL"},
		{"is true", "a IS TRUE"},
		{"is not false", "a IS NOT FALSE"},
		{"between", "a BETWEEN 1 AND 10"},
		{"not between", "a NOT BETWEEN 1 AND 10"},
		{"in list", "a IN (1, 2, 3)"},
		{"not in list", "a NOT IN (1, 2)"},
		{"in subquery", "a IN (SELECT 1)"},
		{"like", "a LIKE '%foo%'"},
		{"not like", "a NOT LIKE 'bar'"},
		{"like escape", "a LIKE '%x' ESCAPE '\\\\'"},
		{"regexp", "a REGEXP '^[0-9]+'"},
		{"rlike", "a RLIKE '^abc'"},
		{"not regexp", "a NOT REGEXP '^x'"},
		{"cast int", "CAST(a AS INT)"},
		{"cast varchar", "CAST(a AS VARCHAR(255))"},
		{"try_cast", "TRY_CAST(a AS VARCHAR(255))"},
		{"cast decimal", "CAST(a + b AS DECIMAL(10,2))"},
		{"case searched", "CASE WHEN a > 1 THEN 'yes' ELSE 'no' END"},
		{"case simple", "CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' END"},
		{"case with else", "CASE a WHEN 1 THEN 'one' ELSE 'other' END"},
		{"count star", "COUNT(*)"},
		{"sum", "SUM(a)"},
		{"coalesce", "COALESCE(a, b, c)"},
		{"count distinct", "COUNT(DISTINCT a)"},
		{"func no args", "NOW()"},
		{"current_date", "CURRENT_DATE"},
		{"current_timestamp", "CURRENT_TIMESTAMP"},
		{"unary minus", "-a"},
		{"unary plus", "+a"},
		{"unary bitnot", "~a"},
		{"bitwise or", "a | b"},
		{"bitwise and", "a & b"},
		{"bitwise xor", "a ^ b"},
		{"shift left", "a << 2"},
		{"shift right", "a >> 2"},
		{"paren", "(1 + 2)"},
		{"exists", "EXISTS (SELECT 1 FROM t)"},
		{"interval", "INTERVAL 1 DAY"},
		{"nested", "a > 0 AND (b IN (1,2) OR c IS NOT NULL)"},
		{"double paren", "((a + b))"},
		{"complex arith", "a * b + c / d - e"},
		{"and symbol", "a && b"},
		{"or symbol", "a || b"},
		{"multi comparison", "a = 1 AND b > 2 AND c < 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseExprFrom(tt.input)
			if err != nil {
				t.Errorf("parseExpr(%q) error: %v", tt.input, err)
			}
			if node == nil {
				t.Errorf("parseExpr(%q) returned nil", tt.input)
			}
		})
	}
}
