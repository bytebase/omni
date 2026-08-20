package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
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

// ---------------------------------------------------------------------------
// EXTRACT expressions
// ---------------------------------------------------------------------------

func TestExprExtractUnits(t *testing.T) {
	// Every unit documented for EXTRACT: plain identifiers, non-reserved unit
	// keywords, and the compound units that lex as reserved keywords.
	units := []string{
		"YEAR", "QUARTER", "MONTH", "WEEK", "DAY", "HOUR", "MINUTE", "SECOND",
		"YEAR_MONTH", "DAY_HOUR", "DAY_MINUTE", "DAY_SECOND", "DAY_MICROSECOND",
		"HOUR_MINUTE", "HOUR_SECOND", "HOUR_MICROSECOND",
		"MINUTE_SECOND", "MINUTE_MICROSECOND", "SECOND_MICROSECOND",
		"DAYOFWEEK", "DOW", "DAYOFYEAR", "DOY",
	}
	for _, unit := range units {
		t.Run(unit, func(t *testing.T) {
			node := mustParseExpr(t, "EXTRACT("+unit+" FROM d)")
			ex, ok := node.(*ast.ExtractExpr)
			if !ok {
				t.Fatalf("expected *ast.ExtractExpr, got %T", node)
			}
			if ex.Unit != unit {
				t.Errorf("unit = %q, want %q", ex.Unit, unit)
			}
			if _, ok := ex.Expr.(*ast.ColumnRef); !ok {
				t.Errorf("expr = %T, want *ast.ColumnRef", ex.Expr)
			}
		})
	}
}

func TestExprExtractUnitNormalized(t *testing.T) {
	node := mustParseExpr(t, "extract(year from d)")
	ex, ok := node.(*ast.ExtractExpr)
	if !ok {
		t.Fatalf("expected *ast.ExtractExpr, got %T", node)
	}
	if ex.Unit != "YEAR" {
		t.Errorf("unit = %q, want YEAR", ex.Unit)
	}
}

func TestExprExtractNestedFuncCall(t *testing.T) {
	node := mustParseExpr(t, "EXTRACT(YEAR FROM MONTHS_ADD(NOW(), -1))")
	ex, ok := node.(*ast.ExtractExpr)
	if !ok {
		t.Fatalf("expected *ast.ExtractExpr, got %T", node)
	}
	fc, ok := ex.Expr.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expr = %T, want *ast.FuncCallExpr", ex.Expr)
	}
	if got := fc.Name.Parts[len(fc.Name.Parts)-1]; got != "MONTHS_ADD" {
		t.Errorf("func name = %q, want MONTHS_ADD", got)
	}
}

func TestExprExtractInsideCTEStatement(t *testing.T) {
	// The reported shape: EXTRACT inside a CTE body, spread over several lines.
	sql := `WITH c AS (
  SELECT year_num
  FROM currency_record
  WHERE year_num = EXTRACT(
    YEAR
    FROM
      MONTHS_ADD(NOW(), -1)
  )
)
SELECT * FROM c`
	file, errs := Parse(sql)
	if len(errs) > 0 {
		t.Fatalf("Parse returned errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Stmts len = %d, want 1", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.SelectStmt); !ok {
		t.Fatalf("stmt = %T, want *ast.SelectStmt", file.Stmts[0])
	}
}

func TestExprExtractSyntaxErrors(t *testing.T) {
	for _, input := range []string{
		"EXTRACT(YEAR d)",     // missing FROM
		"EXTRACT(FROM d)",     // missing unit
		"EXTRACT(YEAR FROM)",  // missing source expression
		"EXTRACT YEAR FROM d", // missing parentheses
	} {
		if _, err := parseExprFrom(input); err == nil {
			t.Errorf("parseExpr(%q) = nil error, want syntax error", input)
		}
	}
}

// ---------------------------------------------------------------------------
// Reserved keywords usable as function names (functionNameIdentifier)
// ---------------------------------------------------------------------------

func TestExprFunctionNameKeywords(t *testing.T) {
	tests := []struct {
		input string
		name  string
		args  int
	}{
		{"TRIM(s)", "TRIM", 1},
		{"TRIM(s, 'x')", "TRIM", 2},
		{"LEFT(s, 2)", "LEFT", 2},
		{"RIGHT(s, 2)", "RIGHT", 2},
		{"IF(a, 1, 2)", "IF", 3},
		{"DATABASE()", "DATABASE", 0},
		{"SCHEMA()", "SCHEMA", 0},
		{"ADD(1, 2)", "ADD", 2},
		{"REGEXP(1)", "REGEXP", 1},
		{"LIKE(1)", "LIKE", 1},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := mustParseExpr(t, tt.input)
			fc, ok := node.(*ast.FuncCallExpr)
			if !ok {
				t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
			}
			if got := fc.Name.Parts[len(fc.Name.Parts)-1]; !strings.EqualFold(got, tt.name) {
				t.Errorf("name = %q, want %q", got, tt.name)
			}
			if len(fc.Args) != tt.args {
				t.Errorf("args = %d, want %d", len(fc.Args), tt.args)
			}
		})
	}
}

func TestExprFunctionNameKeywordRequiresCall(t *testing.T) {
	// The keyword is only a function name when directly applied — `LEFT JOIN`
	// must still parse as a join, and a bare `LEFT` is still not an identifier.
	if _, errs := Parse("SELECT * FROM a LEFT JOIN b ON a.x = b.x"); len(errs) > 0 {
		t.Errorf("LEFT JOIN broke: %v", errs)
	}
	if _, errs := Parse("SELECT * FROM a RIGHT JOIN b ON a.x = b.x"); len(errs) > 0 {
		t.Errorf("RIGHT JOIN broke: %v", errs)
	}
	if _, errs := Parse("SELECT LEFT FROM t"); len(errs) == 0 {
		t.Error("bare LEFT as a column parsed, want syntax error")
	}
}

func TestExprNilaryFunctionPrecision(t *testing.T) {
	for _, input := range []string{
		"CURRENT_TIMESTAMP(3)", "CURRENT_TIME(3)", "LOCALTIME(3)", "LOCALTIMESTAMP(3)",
	} {
		t.Run(input, func(t *testing.T) {
			node := mustParseExpr(t, input)
			fc, ok := node.(*ast.FuncCallExpr)
			if !ok {
				t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
			}
			if len(fc.Args) != 1 {
				t.Fatalf("args = %d, want 1", len(fc.Args))
			}
		})
	}
	// The bare and empty-paren forms must keep working.
	for _, input := range []string{"CURRENT_TIMESTAMP", "CURRENT_DATE()"} {
		node := mustParseExpr(t, input)
		if fc, ok := node.(*ast.FuncCallExpr); !ok || len(fc.Args) != 0 {
			t.Errorf("%s = %T with args, want a zero-arg FuncCallExpr", input, node)
		}
	}
}

// ---------------------------------------------------------------------------
// Literals: hex, bit, array, map
// ---------------------------------------------------------------------------

func TestExprHexAndBitLiterals(t *testing.T) {
	hex := mustParseExpr(t, "x'1f'")
	if lit, ok := hex.(*ast.Literal); !ok || lit.Kind != ast.LitHex {
		t.Errorf("x'1f' = %T (%v), want LitHex literal", hex, hex)
	}
	bit := mustParseExpr(t, "b'101'")
	if lit, ok := bit.(*ast.Literal); !ok || lit.Kind != ast.LitBit {
		t.Errorf("b'101' = %T, want LitBit literal", bit)
	}
}

func TestExprArrayLiteral(t *testing.T) {
	node := mustParseExpr(t, "[1, 2, 3]")
	lit, ok := node.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteral, got %T", node)
	}
	if len(lit.Elements) != 3 {
		t.Errorf("elements = %d, want 3", len(lit.Elements))
	}
	if empty, ok := mustParseExpr(t, "[]").(*ast.ArrayLiteral); !ok || len(empty.Elements) != 0 {
		t.Error("[] did not parse as an empty array literal")
	}
}

func TestExprMapLiteral(t *testing.T) {
	node := mustParseExpr(t, "{'a': 1, 'b': 2}")
	lit, ok := node.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("expected *ast.MapLiteral, got %T", node)
	}
	if len(lit.Entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(lit.Entries))
	}
	if lit.Entries[0].Key == nil || lit.Entries[0].Value == nil {
		t.Error("first entry has a nil key or value")
	}
}

// ---------------------------------------------------------------------------
// Session and user variables
// ---------------------------------------------------------------------------

func TestExprVariableRef(t *testing.T) {
	tests := []struct {
		input  string
		name   string
		scope  string
		system bool
	}{
		{"@@version_comment", "version_comment", "", true},
		{"@@session.query_timeout", "query_timeout", "SESSION", true},
		{"@@global.query_timeout", "query_timeout", "GLOBAL", true},
		{"@uservar", "uservar", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := mustParseExpr(t, tt.input)
			v, ok := node.(*ast.VariableRef)
			if !ok {
				t.Fatalf("expected *ast.VariableRef, got %T", node)
			}
			if !strings.EqualFold(v.Name, tt.name) || v.Scope != tt.scope || v.System != tt.system {
				t.Errorf("got {Name:%q Scope:%q System:%v}, want {%q %q %v}",
					v.Name, v.Scope, v.System, tt.name, tt.scope, tt.system)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Lambdas
// ---------------------------------------------------------------------------

func TestExprLambdaInFunctionArg(t *testing.T) {
	node := mustParseExpr(t, "array_map(x -> x + 1, arr)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 2 {
		t.Fatalf("args = %d, want 2", len(fc.Args))
	}
	lam, ok := fc.Args[0].(*ast.LambdaExpr)
	if !ok {
		t.Fatalf("args[0] = %T, want *ast.LambdaExpr", fc.Args[0])
	}
	if len(lam.Params) != 1 || lam.Params[0] != "x" {
		t.Errorf("params = %v, want [x]", lam.Params)
	}
	if _, ok := lam.Body.(*ast.BinaryExpr); !ok {
		t.Errorf("body = %T, want *ast.BinaryExpr", lam.Body)
	}
}

// ---------------------------------------------------------------------------
// CONVERT / USING charset / BINARY
// ---------------------------------------------------------------------------

func TestExprConvertUsing(t *testing.T) {
	node := mustParseExpr(t, "CONVERT(s USING utf8)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Using != "UTF8" {
		t.Errorf("Using = %q, want UTF8", fc.Using)
	}
	if len(fc.Args) != 1 {
		t.Errorf("args = %d, want 1", len(fc.Args))
	}
}

func TestExprConvertToType(t *testing.T) {
	// CONVERT(expr, type) is CAST spelled differently.
	node := mustParseExpr(t, "CONVERT(s, SIGNED)")
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if cast.TypeName == nil || cast.TypeName.Name != "SIGNED" {
		t.Errorf("type = %+v, want SIGNED", cast.TypeName)
	}
}

func TestExprCharUsing(t *testing.T) {
	node := mustParseExpr(t, "CHAR(65 USING utf8)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Using != "UTF8" {
		t.Errorf("Using = %q, want UTF8", fc.Using)
	}
}

func TestExprBinaryOperator(t *testing.T) {
	node := mustParseExpr(t, "BINARY s")
	un, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if un.Op != ast.UnaryBinary {
		t.Errorf("op = %v, want UnaryBinary", un.Op)
	}
}

// ---------------------------------------------------------------------------
// Statement-level: optimizer hints, GROUP BY CUBE
// ---------------------------------------------------------------------------

func TestStmtOptimizerHintIsSkipped(t *testing.T) {
	for _, sql := range []string{
		"SELECT /*+ SET_VAR(query_timeout = 100) */ 1",
		"SELECT /*+ SET_VAR(a = 1) */ /* ordinary comment */ 1",
		"SELECT /*+ SET_VAR(a = 1) */ a FROM t WHERE a > 1",
	} {
		file, errs := Parse(sql)
		if len(errs) > 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("Parse(%q) produced %d statements, want 1", sql, len(file.Stmts))
		}
	}
}

func TestStmtGroupByCube(t *testing.T) {
	for _, sql := range []string{
		"SELECT a, SUM(b) FROM t GROUP BY CUBE(a)",
		"SELECT a, SUM(b) FROM t GROUP BY CUBE(a, c)",
		"SELECT a, SUM(b) FROM t GROUP BY ROLLUP(a)",
		"SELECT a, SUM(b) FROM t GROUP BY GROUPING SETS ((a), ())",
	} {
		if _, errs := Parse(sql); len(errs) > 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
}

func TestExprLambdaMultiParameter(t *testing.T) {
	node := mustParseExpr(t, "array_map((x, y) -> x + y, a, b)")
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	lam, ok := fc.Args[0].(*ast.LambdaExpr)
	if !ok {
		t.Fatalf("args[0] = %T, want *ast.LambdaExpr", fc.Args[0])
	}
	if len(lam.Params) != 2 || lam.Params[0] != "x" || lam.Params[1] != "y" {
		t.Errorf("params = %v, want [x y]", lam.Params)
	}
}

func TestExprLambdaSingleParameterMustBeBare(t *testing.T) {
	// The engine accepts `x -> ...` and `(x, y) -> ...` but rejects `(x) -> ...`.
	if _, err := parseExprFrom("array_map((x) -> x + 1, arr)"); err == nil {
		t.Error("(x) -> ... parsed, want syntax error")
	}
	// An ordinary parenthesized argument must be unaffected by the speculative
	// lambda parse.
	mustParseExpr(t, "foo((a + b), c)")
	mustParseExpr(t, "foo((a), b)")
}

func TestExprTrailingUsingIsCharOnly(t *testing.T) {
	// CHAR takes a trailing USING on the generic call path; CONVERT has its own
	// production. Every other function must reject it, as the engine does.
	mustParseExpr(t, "CHAR(65 USING utf8)")
	mustParseExpr(t, "CHAR(65, 66 USING utf8)")
	mustParseExpr(t, "CONVERT(s USING utf8)")
	for _, input := range []string{"SUM(a USING utf8)", "foo(a USING utf8)"} {
		if _, err := parseExprFrom(input); err == nil {
			t.Errorf("parseExpr(%q) = nil error, want syntax error", input)
		}
	}
}

func TestStmtGroupByCubeIsWholeSpecification(t *testing.T) {
	// The engine rejects a grouping list after CUBE. Accepting it here would
	// silently discard everything past the comma.
	for _, sql := range []string{
		"SELECT a, b, SUM(c) FROM t GROUP BY CUBE(a), b",
		"SELECT a, SUM(b) FROM t GROUP BY CUBE(a), CUBE(b)",
	} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want syntax error", sql)
		}
	}
}

func TestUnaryBinaryRenders(t *testing.T) {
	if got := ast.UnaryBinary.String(); got != "BINARY" {
		t.Errorf("UnaryBinary.String() = %q, want BINARY", got)
	}
}
