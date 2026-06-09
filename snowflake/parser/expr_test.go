package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseExpr constructs a Parser from input and calls parseExpr.
func testParseExpr(input string) (ast.Node, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseExpr()
}

// ---------------------------------------------------------------------------
// 1. Literals
// ---------------------------------------------------------------------------

func TestExpr_Literals(t *testing.T) {
	tests := []struct {
		input    string
		wantKind ast.LiteralKind
		wantIval int64
		wantBval bool
	}{
		{"42", ast.LitInt, 42, false},
		{"0", ast.LitInt, 0, false},
		{"3.14", ast.LitFloat, 0, false},
		{"TRUE", ast.LitBool, 0, true},
		{"FALSE", ast.LitBool, 0, false},
		{"NULL", ast.LitNull, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node, err := testParseExpr(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			lit, ok := node.(*ast.Literal)
			if !ok {
				t.Fatalf("expected *ast.Literal, got %T", node)
			}
			if lit.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", lit.Kind, tt.wantKind)
			}
			if tt.wantKind == ast.LitInt && lit.Ival != tt.wantIval {
				t.Errorf("Ival = %d, want %d", lit.Ival, tt.wantIval)
			}
			if tt.wantKind == ast.LitBool && lit.Bval != tt.wantBval {
				t.Errorf("Bval = %v, want %v", lit.Bval, tt.wantBval)
			}
		})
	}
}

func TestExpr_StringLiteral(t *testing.T) {
	node, err := testParseExpr("'hello'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lit, ok := node.(*ast.Literal)
	if !ok {
		t.Fatalf("expected *ast.Literal, got %T", node)
	}
	if lit.Kind != ast.LitString {
		t.Errorf("Kind = %v, want LitString", lit.Kind)
	}
	if lit.Value != "hello" {
		t.Errorf("Value = %q, want %q", lit.Value, "hello")
	}
}

// ---------------------------------------------------------------------------
// 2. Column refs
// ---------------------------------------------------------------------------

func TestExpr_ColumnRef(t *testing.T) {
	tests := []struct {
		input     string
		wantParts int
	}{
		{"col", 1},
		{"t.col", 2},
		{"s.t.col", 3},
		{"db.s.t.col", 4},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node, err := testParseExpr(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			cr, ok := node.(*ast.ColumnRef)
			if !ok {
				t.Fatalf("expected *ast.ColumnRef, got %T", node)
			}
			if len(cr.Parts) != tt.wantParts {
				t.Errorf("len(Parts) = %d, want %d", len(cr.Parts), tt.wantParts)
			}
		})
	}
}

func TestExpr_QuotedColumnRef(t *testing.T) {
	node, err := testParseExpr(`"My Column"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cr, ok := node.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", node)
	}
	if len(cr.Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(cr.Parts))
	}
	if cr.Parts[0].Name != "My Column" {
		t.Errorf("Name = %q, want %q", cr.Parts[0].Name, "My Column")
	}
	if !cr.Parts[0].Quoted {
		t.Error("expected Quoted = true")
	}
}

// ---------------------------------------------------------------------------
// 3. Star
// ---------------------------------------------------------------------------

func TestExpr_Star(t *testing.T) {
	node, err := testParseExpr("*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	star, ok := node.(*ast.StarExpr)
	if !ok {
		t.Fatalf("expected *ast.StarExpr, got %T", node)
	}
	if star.Qualifier != nil {
		t.Error("expected nil Qualifier for bare *")
	}
}

func TestExpr_QualifiedStar(t *testing.T) {
	node, err := testParseExpr("t.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	star, ok := node.(*ast.StarExpr)
	if !ok {
		t.Fatalf("expected *ast.StarExpr, got %T", node)
	}
	if star.Qualifier == nil {
		t.Fatal("expected non-nil Qualifier")
	}
	if star.Qualifier.Name.Name != "t" {
		t.Errorf("Qualifier.Name = %q, want %q", star.Qualifier.Name.Name, "t")
	}
}

// ---------------------------------------------------------------------------
// 4. Arithmetic precedence
// ---------------------------------------------------------------------------

func TestExpr_ArithmeticPrecedence(t *testing.T) {
	// 1 + 2 * 3 should produce BinaryExpr(+, Literal(1), BinaryExpr(*, Literal(2), Literal(3)))
	node, err := testParseExpr("1 + 2 * 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	add, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if add.Op != ast.BinAdd {
		t.Errorf("Op = %v, want BinAdd", add.Op)
	}

	// Left should be Literal(1)
	left, ok := add.Left.(*ast.Literal)
	if !ok {
		t.Fatalf("left: expected *ast.Literal, got %T", add.Left)
	}
	if left.Ival != 1 {
		t.Errorf("left.Ival = %d, want 1", left.Ival)
	}

	// Right should be BinaryExpr(*, 2, 3)
	mul, ok := add.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("right: expected *ast.BinaryExpr, got %T", add.Right)
	}
	if mul.Op != ast.BinMul {
		t.Errorf("right.Op = %v, want BinMul", mul.Op)
	}
	right2, _ := mul.Left.(*ast.Literal)
	right3, _ := mul.Right.(*ast.Literal)
	if right2 == nil || right2.Ival != 2 {
		t.Errorf("right.Left = %v, want Literal(2)", mul.Left)
	}
	if right3 == nil || right3.Ival != 3 {
		t.Errorf("right.Right = %v, want Literal(3)", mul.Right)
	}
}

func TestExpr_ArithmeticLeftAssociativity(t *testing.T) {
	// 1 - 2 - 3 should produce BinaryExpr(-, BinaryExpr(-, 1, 2), 3)
	node, err := testParseExpr("1 - 2 - 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outer, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if outer.Op != ast.BinSub {
		t.Errorf("Op = %v, want BinSub", outer.Op)
	}
	inner, ok := outer.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("left: expected *ast.BinaryExpr, got %T", outer.Left)
	}
	if inner.Op != ast.BinSub {
		t.Errorf("inner.Op = %v, want BinSub", inner.Op)
	}
}

// ---------------------------------------------------------------------------
// 5. Parenthesized expressions
// ---------------------------------------------------------------------------

func TestExpr_ParenExpr(t *testing.T) {
	// (1 + 2) * 3 should produce BinaryExpr(*, ParenExpr(BinaryExpr(+, 1, 2)), 3)
	node, err := testParseExpr("(1 + 2) * 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mul, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if mul.Op != ast.BinMul {
		t.Errorf("Op = %v, want BinMul", mul.Op)
	}
	paren, ok := mul.Left.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("left: expected *ast.ParenExpr, got %T", mul.Left)
	}
	add, ok := paren.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("paren.Expr: expected *ast.BinaryExpr, got %T", paren.Expr)
	}
	if add.Op != ast.BinAdd {
		t.Errorf("paren.Expr.Op = %v, want BinAdd", add.Op)
	}
}

// ---------------------------------------------------------------------------
// 6. Comparison operators
// ---------------------------------------------------------------------------

func TestExpr_Comparison(t *testing.T) {
	tests := []struct {
		input  string
		wantOp ast.BinaryOp
	}{
		{"a = b", ast.BinEq},
		{"a <> b", ast.BinNe},
		{"a != b", ast.BinNe},
		{"a < b", ast.BinLt},
		{"a > b", ast.BinGt},
		{"a <= b", ast.BinLe},
		{"a >= b", ast.BinGe},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node, err := testParseExpr(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			binExpr, ok := node.(*ast.BinaryExpr)
			if !ok {
				t.Fatalf("expected *ast.BinaryExpr, got %T", node)
			}
			if binExpr.Op != tt.wantOp {
				t.Errorf("Op = %v, want %v", binExpr.Op, tt.wantOp)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 7. Logical operators
// ---------------------------------------------------------------------------

func TestExpr_LogicalPrecedence(t *testing.T) {
	// a AND b OR c => BinaryExpr(OR, BinaryExpr(AND, a, b), c)
	node, err := testParseExpr("a AND b OR c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	orExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if orExpr.Op != ast.BinOr {
		t.Errorf("top Op = %v, want BinOr", orExpr.Op)
	}
	andExpr, ok := orExpr.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("left: expected *ast.BinaryExpr, got %T", orExpr.Left)
	}
	if andExpr.Op != ast.BinAnd {
		t.Errorf("left.Op = %v, want BinAnd", andExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// 8. NOT prefix
// ---------------------------------------------------------------------------

func TestExpr_NotPrefix(t *testing.T) {
	node, err := testParseExpr("NOT a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if unary.Op != ast.UnaryNot {
		t.Errorf("Op = %v, want UnaryNot", unary.Op)
	}
}

func TestExpr_NotNotPrefix(t *testing.T) {
	node, err := testParseExpr("NOT NOT a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outer, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if outer.Op != ast.UnaryNot {
		t.Errorf("Op = %v, want UnaryNot", outer.Op)
	}
	inner, ok := outer.Expr.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("inner: expected *ast.UnaryExpr, got %T", outer.Expr)
	}
	if inner.Op != ast.UnaryNot {
		t.Errorf("inner.Op = %v, want UnaryNot", inner.Op)
	}
}

// ---------------------------------------------------------------------------
// 9. Unary minus and plus
// ---------------------------------------------------------------------------

func TestExpr_UnaryMinus(t *testing.T) {
	node, err := testParseExpr("-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if unary.Op != ast.UnaryMinus {
		t.Errorf("Op = %v, want UnaryMinus", unary.Op)
	}
}

func TestExpr_UnaryPlus(t *testing.T) {
	node, err := testParseExpr("+a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected *ast.UnaryExpr, got %T", node)
	}
	if unary.Op != ast.UnaryPlus {
		t.Errorf("Op = %v, want UnaryPlus", unary.Op)
	}
}

// ---------------------------------------------------------------------------
// 10. Concat
// ---------------------------------------------------------------------------

func TestExpr_Concat(t *testing.T) {
	// a || b || c => left-assoc: BinaryExpr(||, BinaryExpr(||, a, b), c)
	node, err := testParseExpr("a || b || c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outer, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if outer.Op != ast.BinConcat {
		t.Errorf("Op = %v, want BinConcat", outer.Op)
	}
	inner, ok := outer.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("left: expected *ast.BinaryExpr, got %T", outer.Left)
	}
	if inner.Op != ast.BinConcat {
		t.Errorf("inner.Op = %v, want BinConcat", inner.Op)
	}
}

// ---------------------------------------------------------------------------
// 11. CAST / TRY_CAST / ::
// ---------------------------------------------------------------------------

func TestExpr_Cast(t *testing.T) {
	node, err := testParseExpr("CAST(x AS INT)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if cast.TryCast {
		t.Error("TryCast should be false for CAST")
	}
	if cast.ColonColon {
		t.Error("ColonColon should be false for CAST")
	}
	if cast.TypeName.Kind != ast.TypeInt {
		t.Errorf("TypeName.Kind = %v, want TypeInt", cast.TypeName.Kind)
	}
}

func TestExpr_TryCast(t *testing.T) {
	node, err := testParseExpr("TRY_CAST(x AS VARCHAR(100))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if !cast.TryCast {
		t.Error("TryCast should be true for TRY_CAST")
	}
	if cast.TypeName.Kind != ast.TypeVarchar {
		t.Errorf("TypeName.Kind = %v, want TypeVarchar", cast.TypeName.Kind)
	}
}

func TestExpr_ColonColonCast(t *testing.T) {
	node, err := testParseExpr("x::INT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if !cast.ColonColon {
		t.Error("ColonColon should be true for ::")
	}
	if cast.TypeName.Kind != ast.TypeInt {
		t.Errorf("TypeName.Kind = %v, want TypeInt", cast.TypeName.Kind)
	}
}

// ---------------------------------------------------------------------------
// 12. CASE simple
// ---------------------------------------------------------------------------

func TestExpr_CaseSimple(t *testing.T) {
	node, err := testParseExpr("CASE x WHEN 1 THEN 'a' WHEN 2 THEN 'b' ELSE 'c' END")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ce, ok := node.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected *ast.CaseExpr, got %T", node)
	}
	if ce.Kind != ast.CaseSimple {
		t.Errorf("Kind = %v, want CaseSimple", ce.Kind)
	}
	if ce.Operand == nil {
		t.Error("Operand should be non-nil for simple CASE")
	}
	if len(ce.Whens) != 2 {
		t.Errorf("len(Whens) = %d, want 2", len(ce.Whens))
	}
	if ce.Else == nil {
		t.Error("Else should be non-nil")
	}
}

// ---------------------------------------------------------------------------
// 13. CASE searched
// ---------------------------------------------------------------------------

func TestExpr_CaseSearched(t *testing.T) {
	node, err := testParseExpr("CASE WHEN x > 0 THEN 'pos' ELSE 'neg' END")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ce, ok := node.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected *ast.CaseExpr, got %T", node)
	}
	if ce.Kind != ast.CaseSearched {
		t.Errorf("Kind = %v, want CaseSearched", ce.Kind)
	}
	if ce.Operand != nil {
		t.Error("Operand should be nil for searched CASE")
	}
	if len(ce.Whens) != 1 {
		t.Errorf("len(Whens) = %d, want 1", len(ce.Whens))
	}
}

func TestExpr_CaseNoElse(t *testing.T) {
	node, err := testParseExpr("CASE WHEN x > 0 THEN 'pos' END")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ce, ok := node.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected *ast.CaseExpr, got %T", node)
	}
	if ce.Else != nil {
		t.Error("Else should be nil when absent")
	}
}

// ---------------------------------------------------------------------------
// 14. IFF
// ---------------------------------------------------------------------------

func TestExpr_Iff(t *testing.T) {
	node, err := testParseExpr("IFF(x > 0, 'pos', 'neg')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	iff, ok := node.(*ast.IffExpr)
	if !ok {
		t.Fatalf("expected *ast.IffExpr, got %T", node)
	}
	if iff.Cond == nil || iff.Then == nil || iff.Else == nil {
		t.Error("Cond, Then, Else must all be non-nil")
	}
}

// ---------------------------------------------------------------------------
// 15. IS NULL / IS NOT NULL / IS DISTINCT FROM / IS NOT DISTINCT FROM
// ---------------------------------------------------------------------------

func TestExpr_IsNull(t *testing.T) {
	node, err := testParseExpr("x IS NULL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if is.Not {
		t.Error("Not should be false")
	}
	if !is.Null {
		t.Error("Null should be true")
	}
}

func TestExpr_IsNotNull(t *testing.T) {
	node, err := testParseExpr("x IS NOT NULL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if !is.Not {
		t.Error("Not should be true")
	}
	if !is.Null {
		t.Error("Null should be true")
	}
}

func TestExpr_IsDistinctFrom(t *testing.T) {
	node, err := testParseExpr("x IS DISTINCT FROM y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if is.Not {
		t.Error("Not should be false")
	}
	if is.Null {
		t.Error("Null should be false for DISTINCT FROM")
	}
	if is.DistinctFrom == nil {
		t.Error("DistinctFrom should be non-nil")
	}
}

func TestExpr_IsNotDistinctFrom(t *testing.T) {
	node, err := testParseExpr("x IS NOT DISTINCT FROM y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	is, ok := node.(*ast.IsExpr)
	if !ok {
		t.Fatalf("expected *ast.IsExpr, got %T", node)
	}
	if !is.Not {
		t.Error("Not should be true")
	}
	if is.Null {
		t.Error("Null should be false for DISTINCT FROM")
	}
	if is.DistinctFrom == nil {
		t.Error("DistinctFrom should be non-nil")
	}
}

// ---------------------------------------------------------------------------
// 16. BETWEEN / NOT BETWEEN
// ---------------------------------------------------------------------------

func TestExpr_Between(t *testing.T) {
	node, err := testParseExpr("x BETWEEN 1 AND 10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	be, ok := node.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected *ast.BetweenExpr, got %T", node)
	}
	if be.Not {
		t.Error("Not should be false")
	}
	if be.Low == nil || be.High == nil {
		t.Error("Low and High must be non-nil")
	}
}

func TestExpr_NotBetween(t *testing.T) {
	node, err := testParseExpr("x NOT BETWEEN 1 AND 10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	be, ok := node.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected *ast.BetweenExpr, got %T", node)
	}
	if !be.Not {
		t.Error("Not should be true")
	}
}

// ---------------------------------------------------------------------------
// 17. IN / NOT IN
// ---------------------------------------------------------------------------

func TestExpr_In(t *testing.T) {
	node, err := testParseExpr("x IN (1, 2, 3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in, ok := node.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", node)
	}
	if in.Not {
		t.Error("Not should be false")
	}
	if len(in.Values) != 3 {
		t.Errorf("len(Values) = %d, want 3", len(in.Values))
	}
}

func TestExpr_NotIn(t *testing.T) {
	node, err := testParseExpr("x NOT IN (1, 2, 3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in, ok := node.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", node)
	}
	if !in.Not {
		t.Error("Not should be true")
	}
}

// ---------------------------------------------------------------------------
// 18. LIKE / ILIKE / RLIKE / REGEXP / LIKE ANY / ESCAPE
// ---------------------------------------------------------------------------

func TestExpr_Like(t *testing.T) {
	node, err := testParseExpr("x LIKE 'pat%'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpLike {
		t.Errorf("Op = %v, want LikeOpLike", le.Op)
	}
	if le.Not {
		t.Error("Not should be false")
	}
	if le.Any {
		t.Error("Any should be false")
	}
}

func TestExpr_NotLike(t *testing.T) {
	node, err := testParseExpr("x NOT LIKE 'pat%'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if !le.Not {
		t.Error("Not should be true")
	}
}

func TestExpr_ILike(t *testing.T) {
	node, err := testParseExpr("x ILIKE '%pat%'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpILike {
		t.Errorf("Op = %v, want LikeOpILike", le.Op)
	}
}

func TestExpr_RLike(t *testing.T) {
	node, err := testParseExpr("x RLIKE '.*'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpRLike {
		t.Errorf("Op = %v, want LikeOpRLike", le.Op)
	}
}

func TestExpr_LikeEscape(t *testing.T) {
	node, err := testParseExpr("x LIKE '%' ESCAPE '#'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Escape == nil {
		t.Error("Escape should be non-nil")
	}
}

func TestExpr_LikeAny(t *testing.T) {
	node, err := testParseExpr("x LIKE ANY ('a%', 'b%')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if !le.Any {
		t.Error("Any should be true")
	}
	if len(le.AnyValues) != 2 {
		t.Errorf("len(AnyValues) = %d, want 2", len(le.AnyValues))
	}
}

func TestExpr_ILikeAny(t *testing.T) {
	node, err := testParseExpr("x ILIKE ANY ('a%', 'b%')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpILike {
		t.Errorf("Op = %v, want LikeOpILike", le.Op)
	}
	if !le.Any {
		t.Error("Any should be true")
	}
}

func TestExpr_Regexp(t *testing.T) {
	node, err := testParseExpr("x REGEXP '.*'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpRegexp {
		t.Errorf("Op = %v, want LikeOpRegexp", le.Op)
	}
}

func TestExpr_NotRegexp(t *testing.T) {
	node, err := testParseExpr("x NOT REGEXP '.*'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpRegexp {
		t.Errorf("Op = %v, want LikeOpRegexp", le.Op)
	}
	if !le.Not {
		t.Error("Not should be true")
	}
}

// ---------------------------------------------------------------------------
// 19. Function calls
// ---------------------------------------------------------------------------

func TestExpr_FuncCallNoArgs(t *testing.T) {
	node, err := testParseExpr("f()")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Name.Name.Name != "f" {
		t.Errorf("Name = %q, want %q", fc.Name.Name.Name, "f")
	}
	if len(fc.Args) != 0 {
		t.Errorf("len(Args) = %d, want 0", len(fc.Args))
	}
}

func TestExpr_FuncCallWithArgs(t *testing.T) {
	node, err := testParseExpr("f(1, 2)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 2 {
		t.Errorf("len(Args) = %d, want 2", len(fc.Args))
	}
}

func TestExpr_CountStar(t *testing.T) {
	node, err := testParseExpr("COUNT(*)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if !fc.Star {
		t.Error("Star should be true for COUNT(*)")
	}
}

func TestExpr_CountDistinct(t *testing.T) {
	node, err := testParseExpr("COUNT(DISTINCT x)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if !fc.Distinct {
		t.Error("Distinct should be true")
	}
	if len(fc.Args) != 1 {
		t.Errorf("len(Args) = %d, want 1", len(fc.Args))
	}
}

func TestExpr_TrimFunc(t *testing.T) {
	node, err := testParseExpr("TRIM(x)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 1 {
		t.Errorf("len(Args) = %d, want 1", len(fc.Args))
	}
}

func TestExpr_QualifiedFunc(t *testing.T) {
	node, err := testParseExpr("my_schema.my_func(1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Name.Schema.Name != "my_schema" {
		t.Errorf("Schema = %q, want %q", fc.Name.Schema.Name, "my_schema")
	}
	if fc.Name.Name.Name != "my_func" {
		t.Errorf("Name = %q, want %q", fc.Name.Name.Name, "my_func")
	}
}

// ---------------------------------------------------------------------------
// 20. Window functions
// ---------------------------------------------------------------------------

func TestExpr_WindowRowNumber(t *testing.T) {
	node, err := testParseExpr("ROW_NUMBER() OVER (PARTITION BY a ORDER BY b)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over == nil {
		t.Fatal("Over should be non-nil for window function")
	}
	if len(fc.Over.PartitionBy) != 1 {
		t.Errorf("len(PartitionBy) = %d, want 1", len(fc.Over.PartitionBy))
	}
	if len(fc.Over.OrderBy) != 1 {
		t.Errorf("len(OrderBy) = %d, want 1", len(fc.Over.OrderBy))
	}
}

func TestExpr_WindowWithFrame(t *testing.T) {
	node, err := testParseExpr("SUM(x) OVER (ORDER BY y ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over == nil {
		t.Fatal("Over should be non-nil")
	}
	if fc.Over.Frame == nil {
		t.Fatal("Frame should be non-nil")
	}
	if fc.Over.Frame.Kind != ast.FrameRows {
		t.Errorf("Frame.Kind = %v, want FrameRows", fc.Over.Frame.Kind)
	}
	if fc.Over.Frame.Start.Kind != ast.BoundUnboundedPreceding {
		t.Errorf("Frame.Start.Kind = %v, want BoundUnboundedPreceding", fc.Over.Frame.Start.Kind)
	}
	if fc.Over.Frame.End.Kind != ast.BoundCurrentRow {
		t.Errorf("Frame.End.Kind = %v, want BoundCurrentRow", fc.Over.Frame.End.Kind)
	}
}

func TestExpr_WindowRangeFrame(t *testing.T) {
	node, err := testParseExpr("SUM(x) OVER (ORDER BY y RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over.Frame.Kind != ast.FrameRange {
		t.Errorf("Frame.Kind = %v, want FrameRange", fc.Over.Frame.Kind)
	}
	if fc.Over.Frame.Start.Kind != ast.BoundPreceding {
		t.Errorf("Frame.Start.Kind = %v, want BoundPreceding", fc.Over.Frame.Start.Kind)
	}
	if fc.Over.Frame.End.Kind != ast.BoundFollowing {
		t.Errorf("Frame.End.Kind = %v, want BoundFollowing", fc.Over.Frame.End.Kind)
	}
}

func TestExpr_WindowGroupsFrame(t *testing.T) {
	node, err := testParseExpr("SUM(x) OVER (ORDER BY y GROUPS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over.Frame.Kind != ast.FrameGroups {
		t.Errorf("Frame.Kind = %v, want FrameGroups", fc.Over.Frame.Kind)
	}
	if fc.Over.Frame.Start.Kind != ast.BoundCurrentRow {
		t.Errorf("Frame.Start.Kind = %v, want BoundCurrentRow", fc.Over.Frame.Start.Kind)
	}
	if fc.Over.Frame.End.Kind != ast.BoundUnboundedFollowing {
		t.Errorf("Frame.End.Kind = %v, want BoundUnboundedFollowing", fc.Over.Frame.End.Kind)
	}
}

func TestExpr_WindowOrderByNulls(t *testing.T) {
	node, err := testParseExpr("RANK() OVER (ORDER BY a DESC NULLS LAST)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over == nil {
		t.Fatal("Over should be non-nil")
	}
	if len(fc.Over.OrderBy) != 1 {
		t.Fatalf("len(OrderBy) = %d, want 1", len(fc.Over.OrderBy))
	}
	item := fc.Over.OrderBy[0]
	if !item.Desc {
		t.Error("Desc should be true")
	}
	if item.NullsFirst == nil {
		t.Fatal("NullsFirst should be non-nil")
	}
	if *item.NullsFirst {
		t.Error("NullsFirst should be false (NULLS LAST)")
	}
}

// ---------------------------------------------------------------------------
// 21. JSON access / Array access / Dot access
// ---------------------------------------------------------------------------

func TestExpr_ColonAccess(t *testing.T) {
	node, err := testParseExpr("v:field")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	acc, ok := node.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("expected *ast.AccessExpr, got %T", node)
	}
	if acc.Kind != ast.AccessColon {
		t.Errorf("Kind = %v, want AccessColon", acc.Kind)
	}
	if acc.Field.Name != "field" {
		t.Errorf("Field.Name = %q, want %q", acc.Field.Name, "field")
	}
}

func TestExpr_BracketAccess(t *testing.T) {
	node, err := testParseExpr("v[0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	acc, ok := node.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("expected *ast.AccessExpr, got %T", node)
	}
	if acc.Kind != ast.AccessBracket {
		t.Errorf("Kind = %v, want AccessBracket", acc.Kind)
	}
	if acc.Index == nil {
		t.Error("Index should be non-nil")
	}
}

func TestExpr_ChainedAccess(t *testing.T) {
	// v:field.subfield should parse as AccessExpr(Dot, AccessExpr(Colon, v, field), subfield)
	node, err := testParseExpr("v:field.subfield")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dot, ok := node.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("expected *ast.AccessExpr, got %T", node)
	}
	if dot.Kind != ast.AccessDot {
		t.Errorf("Kind = %v, want AccessDot", dot.Kind)
	}
	if dot.Field.Name != "subfield" {
		t.Errorf("Field.Name = %q, want %q", dot.Field.Name, "subfield")
	}
	inner, ok := dot.Expr.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("inner: expected *ast.AccessExpr, got %T", dot.Expr)
	}
	if inner.Kind != ast.AccessColon {
		t.Errorf("inner.Kind = %v, want AccessColon", inner.Kind)
	}
}

// ---------------------------------------------------------------------------
// 22. Array literal / JSON literal
// ---------------------------------------------------------------------------

func TestExpr_ArrayLiteral(t *testing.T) {
	node, err := testParseExpr("[1, 2, 3]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := node.(*ast.ArrayLiteralExpr)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteralExpr, got %T", node)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("len(Elements) = %d, want 3", len(arr.Elements))
	}
}

func TestExpr_EmptyArrayLiteral(t *testing.T) {
	node, err := testParseExpr("[]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := node.(*ast.ArrayLiteralExpr)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteralExpr, got %T", node)
	}
	if len(arr.Elements) != 0 {
		t.Errorf("len(Elements) = %d, want 0", len(arr.Elements))
	}
}

func TestExpr_JsonLiteral(t *testing.T) {
	node, err := testParseExpr("{'key': 'val'}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	jl, ok := node.(*ast.JsonLiteralExpr)
	if !ok {
		t.Fatalf("expected *ast.JsonLiteralExpr, got %T", node)
	}
	if len(jl.Pairs) != 1 {
		t.Fatalf("len(Pairs) = %d, want 1", len(jl.Pairs))
	}
	if jl.Pairs[0].Key != "key" {
		t.Errorf("Key = %q, want %q", jl.Pairs[0].Key, "key")
	}
}

func TestExpr_JsonLiteralMultiPair(t *testing.T) {
	node, err := testParseExpr("{'a': 1, 'b': 2}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	jl, ok := node.(*ast.JsonLiteralExpr)
	if !ok {
		t.Fatalf("expected *ast.JsonLiteralExpr, got %T", node)
	}
	if len(jl.Pairs) != 2 {
		t.Errorf("len(Pairs) = %d, want 2", len(jl.Pairs))
	}
}

// ---------------------------------------------------------------------------
// 23. Lambda expressions
// ---------------------------------------------------------------------------

func TestExpr_Lambda(t *testing.T) {
	node, err := testParseExpr("x -> x + 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lam, ok := node.(*ast.LambdaExpr)
	if !ok {
		t.Fatalf("expected *ast.LambdaExpr, got %T", node)
	}
	if len(lam.Params) != 1 {
		t.Errorf("len(Params) = %d, want 1", len(lam.Params))
	}
	if lam.Params[0].Name != "x" {
		t.Errorf("Params[0].Name = %q, want %q", lam.Params[0].Name, "x")
	}
	if lam.Body == nil {
		t.Error("Body should be non-nil")
	}
}

// ---------------------------------------------------------------------------
// 24. Nested expressions
// ---------------------------------------------------------------------------

func TestExpr_NestedCastIff(t *testing.T) {
	node, err := testParseExpr("CAST(IFF(a > b, a, b) AS INT)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	_, ok = cast.Expr.(*ast.IffExpr)
	if !ok {
		t.Fatalf("cast.Expr: expected *ast.IffExpr, got %T", cast.Expr)
	}
}

func TestExpr_NestedFuncCalls(t *testing.T) {
	node, err := testParseExpr("f(g(x))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outer, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(outer.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(outer.Args))
	}
	inner, ok := outer.Args[0].(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("inner: expected *ast.FuncCallExpr, got %T", outer.Args[0])
	}
	if inner.Name.Name.Name != "g" {
		t.Errorf("inner.Name = %q, want %q", inner.Name.Name.Name, "g")
	}
}

func TestExpr_CastPostfixPrecedence(t *testing.T) {
	// a + b * c :: INT
	// :: binds tighter than + and *, so this is: a + (b * (c :: INT))
	// Actually :: is bpPostfix (80) > bpMul (60) > bpAdd (50)
	// So: a + (b * (c::INT))
	node, err := testParseExpr("a + b * c :: INT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	add, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if add.Op != ast.BinAdd {
		t.Errorf("Op = %v, want BinAdd", add.Op)
	}
	mul, ok := add.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("right: expected *ast.BinaryExpr, got %T", add.Right)
	}
	if mul.Op != ast.BinMul {
		t.Errorf("right.Op = %v, want BinMul", mul.Op)
	}
	_, ok = mul.Right.(*ast.CastExpr)
	if !ok {
		t.Fatalf("mul.Right: expected *ast.CastExpr, got %T", mul.Right)
	}
}

// ---------------------------------------------------------------------------
// 25. COLLATE
// ---------------------------------------------------------------------------

func TestExpr_Collate(t *testing.T) {
	node, err := testParseExpr("x COLLATE utf8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	col, ok := node.(*ast.CollateExpr)
	if !ok {
		t.Fatalf("expected *ast.CollateExpr, got %T", node)
	}
	if col.Collation != "utf8" {
		t.Errorf("Collation = %q, want %q", col.Collation, "utf8")
	}
}

func TestExpr_CollateString(t *testing.T) {
	node, err := testParseExpr("x COLLATE 'en_US'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	col, ok := node.(*ast.CollateExpr)
	if !ok {
		t.Fatalf("expected *ast.CollateExpr, got %T", node)
	}
	if col.Collation != "en_US" {
		t.Errorf("Collation = %q, want %q", col.Collation, "en_US")
	}
}

// ---------------------------------------------------------------------------
// 25b. Oracle-style outer-join marker (+)
// ---------------------------------------------------------------------------

// mustOuterJoin parses input and asserts the result is an *ast.OuterJoinExpr.
func mustOuterJoin(t *testing.T, input string) *ast.OuterJoinExpr {
	t.Helper()
	node, err := testParseExpr(input)
	if err != nil {
		t.Fatalf("parse %q: unexpected error: %v", input, err)
	}
	oj, ok := node.(*ast.OuterJoinExpr)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.OuterJoinExpr", input, node)
	}
	return oj
}

func TestExpr_OuterJoinMarker_NoSpace(t *testing.T) {
	oj := mustOuterJoin(t, "t2.c2(+)")
	cr, ok := oj.Operand.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("Operand = %T, want *ast.ColumnRef", oj.Operand)
	}
	if len(cr.Parts) != 2 || cr.Parts[0].Name != "t2" || cr.Parts[1].Name != "c2" {
		t.Errorf("Operand parts = %+v, want [t2 c2]", cr.Parts)
	}
	// Loc spans the operand through the closing ')'.
	if oj.Loc.Start != cr.Loc.Start {
		t.Errorf("Loc.Start = %d, want operand start %d", oj.Loc.Start, cr.Loc.Start)
	}
	if oj.Loc.End <= cr.Loc.End {
		t.Errorf("Loc.End = %d, want > operand end %d (must cover `(+)`)", oj.Loc.End, cr.Loc.End)
	}
}

func TestExpr_OuterJoinMarker_WithSpace(t *testing.T) {
	// A space before the marker is allowed: `t2.c2 (+)`.
	oj := mustOuterJoin(t, "t2.c2 (+)")
	if _, ok := oj.Operand.(*ast.ColumnRef); !ok {
		t.Fatalf("Operand = %T, want *ast.ColumnRef", oj.Operand)
	}
}

func TestExpr_OuterJoinMarker_SingleColumn(t *testing.T) {
	// Unqualified column also accepts the marker.
	oj := mustOuterJoin(t, "c(+)")
	cr, ok := oj.Operand.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("Operand = %T, want *ast.ColumnRef", oj.Operand)
	}
	if len(cr.Parts) != 1 || cr.Parts[0].Name != "c" {
		t.Errorf("Operand parts = %+v, want [c]", cr.Parts)
	}
}

func TestExpr_OuterJoinMarker_InComparison(t *testing.T) {
	// `t1.c1 = t2.c2(+)` — the marker binds tightly to the right operand.
	node, err := testParseExpr("t1.c1 = t2.c2(+)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("got %T, want *ast.BinaryExpr", node)
	}
	if bin.Op != ast.BinEq {
		t.Errorf("Op = %v, want BinEq", bin.Op)
	}
	if _, ok := bin.Right.(*ast.OuterJoinExpr); !ok {
		t.Errorf("Right = %T, want *ast.OuterJoinExpr", bin.Right)
	}
	if _, ok := bin.Left.(*ast.OuterJoinExpr); ok {
		t.Errorf("Left wrongly wrapped in *ast.OuterJoinExpr")
	}
}

func TestExpr_OuterJoinMarker_MultipleInWhere(t *testing.T) {
	// Two markers in one WHERE (both spaced), via the full statement parser.
	sql := "SELECT t1.c1, t2.c2 FROM t1, t2 WHERE t1.c1 = t2.c2 (+) AND t1.c3 = t2.c4 (+);"
	f, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	if len(f.Stmts) != 1 {
		t.Fatalf("len(Stmts) = %d, want 1", len(f.Stmts))
	}
	markers := 0
	ast.Inspect(f, func(n ast.Node) bool {
		if _, ok := n.(*ast.OuterJoinExpr); ok {
			markers++
		}
		return true
	})
	if markers != 2 {
		t.Errorf("OuterJoinExpr count = %d, want 2", markers)
	}
}

// Regression: a genuine function call `count(*)` must stay a FuncCallExpr and
// never be mistaken for the `(+)` marker.
func TestExpr_OuterJoinMarker_CountStarUnchanged(t *testing.T) {
	fc := mustFuncCall(t, "count(*)")
	if !fc.Star {
		t.Errorf("count(*) Star = false, want true")
	}
}

// Regression: a plain call `f(a)` is unaffected by the marker lookahead.
func TestExpr_OuterJoinMarker_PlainCallUnchanged(t *testing.T) {
	fc := mustFuncCall(t, "f(a)")
	if len(fc.Args) != 1 {
		t.Errorf("f(a) Args = %d, want 1", len(fc.Args))
	}
}

// Regression: `f(+1)` is a call whose single argument is a unary-plus literal,
// NOT an outer-join marker (the token after '+' is an int, not ')').
func TestExpr_OuterJoinMarker_UnaryPlusArgUnchanged(t *testing.T) {
	fc := mustFuncCall(t, "f(+1)")
	if len(fc.Args) != 1 {
		t.Fatalf("f(+1) Args = %d, want 1", len(fc.Args))
	}
	if _, ok := fc.Args[0].(*ast.UnaryExpr); !ok {
		t.Errorf("f(+1) arg0 = %T, want *ast.UnaryExpr", fc.Args[0])
	}
}

// Regression: a parenthesized `(a+1)` expression is unchanged (it is not a
// column ref followed by `(+)`).
func TestExpr_OuterJoinMarker_ParenExprUnchanged(t *testing.T) {
	node, err := testParseExpr("(a+1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := node.(*ast.ParenExpr); !ok {
		t.Errorf("(a+1) = %T, want *ast.ParenExpr", node)
	}
}

// ---------------------------------------------------------------------------
// 26. Error cases
// ---------------------------------------------------------------------------

func TestExpr_Error_UnclosedParen(t *testing.T) {
	_, err := testParseExpr("(1 + 2")
	if err == nil {
		t.Fatal("expected error for unclosed paren")
	}
}

func TestExpr_Error_MissingThenInCase(t *testing.T) {
	_, err := testParseExpr("CASE WHEN x > 0 ELSE 'neg' END")
	if err == nil {
		t.Fatal("expected error for missing THEN in CASE")
	}
}

func TestExpr_Error_UnterminatedBetween(t *testing.T) {
	_, err := testParseExpr("x BETWEEN 1")
	if err == nil {
		t.Fatal("expected error for unterminated BETWEEN")
	}
}

func TestExpr_SubqueryExpr(t *testing.T) {
	node, err := testParseExpr("(SELECT 1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	subq, ok := node.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("expected *ast.SubqueryExpr, got %T", node)
	}
	sel, ok := subq.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt inside SubqueryExpr, got %T", subq.Query)
	}
	if len(sel.Targets) != 1 {
		t.Errorf("targets = %d, want 1", len(sel.Targets))
	}
}

func TestExpr_ExistsExpr(t *testing.T) {
	node, err := testParseExpr("EXISTS (SELECT 1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exists, ok := node.(*ast.ExistsExpr)
	if !ok {
		t.Fatalf("expected *ast.ExistsExpr, got %T", node)
	}
	sel, ok := exists.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt inside ExistsExpr, got %T", exists.Query)
	}
	if len(sel.Targets) != 1 {
		t.Errorf("targets = %d, want 1", len(sel.Targets))
	}
}

// ---------------------------------------------------------------------------
// ParseExpr public helper
// ---------------------------------------------------------------------------

func TestParseExpr_Public(t *testing.T) {
	node, errs := ParseExpr("1 + 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	_, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
}

func TestParseExpr_TrailingToken(t *testing.T) {
	_, errs := ParseExpr("1 + 2 extra")
	if len(errs) == 0 {
		t.Fatal("expected error for trailing token")
	}
}

// ---------------------------------------------------------------------------
// WITHIN GROUP support
// ---------------------------------------------------------------------------

func TestExpr_WithinGroup(t *testing.T) {
	node, err := testParseExpr("LISTAGG(x, ',') WITHIN GROUP (ORDER BY y)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.OrderBy) != 1 {
		t.Errorf("len(OrderBy) = %d, want 1 for WITHIN GROUP", len(fc.OrderBy))
	}
}

// ---------------------------------------------------------------------------
// Additional edge cases
// ---------------------------------------------------------------------------

func TestExpr_ModOperator(t *testing.T) {
	node, err := testParseExpr("10 % 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinMod {
		t.Errorf("Op = %v, want BinMod", bin.Op)
	}
}

func TestExpr_DivOperator(t *testing.T) {
	node, err := testParseExpr("10 / 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinDiv {
		t.Errorf("Op = %v, want BinDiv", bin.Op)
	}
}

func TestExpr_NotILike(t *testing.T) {
	node, err := testParseExpr("x NOT ILIKE '%foo%'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpILike {
		t.Errorf("Op = %v, want LikeOpILike", le.Op)
	}
	if !le.Not {
		t.Error("Not should be true")
	}
}

func TestExpr_NotRLike(t *testing.T) {
	node, err := testParseExpr("x NOT RLIKE '.*'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	le, ok := node.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected *ast.LikeExpr, got %T", node)
	}
	if le.Op != ast.LikeOpRLike {
		t.Errorf("Op = %v, want LikeOpRLike", le.Op)
	}
	if !le.Not {
		t.Error("Not should be true")
	}
}

func TestExpr_MultipleColonCast(t *testing.T) {
	// x::INT::VARCHAR should parse as CastExpr(CastExpr(x, INT), VARCHAR)
	node, err := testParseExpr("x::INT::VARCHAR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outer, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if outer.TypeName.Kind != ast.TypeVarchar {
		t.Errorf("outer TypeName.Kind = %v, want TypeVarchar", outer.TypeName.Kind)
	}
	inner, ok := outer.Expr.(*ast.CastExpr)
	if !ok {
		t.Fatalf("inner: expected *ast.CastExpr, got %T", outer.Expr)
	}
	if inner.TypeName.Kind != ast.TypeInt {
		t.Errorf("inner TypeName.Kind = %v, want TypeInt", inner.TypeName.Kind)
	}
}

func TestExpr_WindowSingleBound(t *testing.T) {
	// Window frame with single bound (no BETWEEN)
	node, err := testParseExpr("SUM(x) OVER (ORDER BY y ROWS UNBOUNDED PRECEDING)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if fc.Over == nil || fc.Over.Frame == nil {
		t.Fatal("expected frame")
	}
	if fc.Over.Frame.Start.Kind != ast.BoundUnboundedPreceding {
		t.Errorf("Start.Kind = %v, want BoundUnboundedPreceding", fc.Over.Frame.Start.Kind)
	}
}

// ---------------------------------------------------------------------------
// $-references ($N positional column refs, $<name> variable refs)
//
// Oracle: official Snowflake docs (truth1) + the previously-filtered corpus
// forms (copyDollarLimited / execCallDollarLimited / SET $var). The legacy
// SnowflakeParser.g4 is not vendored in this tree, so docs + corpus are the
// authoritative sources. The lexer collapses both forms into a single
// tokVariable (leading '$' stripped); Positional is derived from the text.
// ---------------------------------------------------------------------------

func TestExpr_DollarPositional(t *testing.T) {
	node, err := testParseExpr("$1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dr, ok := node.(*ast.DollarRef)
	if !ok {
		t.Fatalf("expected *ast.DollarRef, got %T", node)
	}
	if dr.Name != "1" {
		t.Errorf("Name = %q, want %q", dr.Name, "1")
	}
	if !dr.Positional {
		t.Error("Positional should be true for $1")
	}
	if dr.Qualifier != nil {
		t.Error("expected nil Qualifier for bare $1")
	}
	// Loc must span the whole token including the leading '$'.
	if dr.Loc.Start != 0 || dr.Loc.End != 2 {
		t.Errorf("Loc = %+v, want {0 2}", dr.Loc)
	}
}

func TestExpr_DollarPositionalMultiDigit(t *testing.T) {
	node, err := testParseExpr("$42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dr, ok := node.(*ast.DollarRef)
	if !ok {
		t.Fatalf("expected *ast.DollarRef, got %T", node)
	}
	if dr.Name != "42" || !dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"42\" true", dr.Name, dr.Positional)
	}
}

func TestExpr_DollarNamedVar(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
	}{
		{"$min", "min"},
		{"$Variable1", "Variable1"},
		{"$ABC_123", "ABC_123"},
		{"$_x", "_x"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node, err := testParseExpr(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			dr, ok := node.(*ast.DollarRef)
			if !ok {
				t.Fatalf("expected *ast.DollarRef, got %T", node)
			}
			if dr.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", dr.Name, tt.wantName)
			}
			if dr.Positional {
				t.Errorf("Positional should be false for %s", tt.input)
			}
			if dr.Qualifier != nil {
				t.Error("expected nil Qualifier")
			}
		})
	}
}

func TestExpr_DollarQualifiedPositional(t *testing.T) {
	// d.$1 — positional column ref qualified by table alias d (COPY transform).
	node, err := testParseExpr("d.$1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dr, ok := node.(*ast.DollarRef)
	if !ok {
		t.Fatalf("expected *ast.DollarRef, got %T", node)
	}
	if dr.Name != "1" || !dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"1\" true", dr.Name, dr.Positional)
	}
	if dr.Qualifier == nil {
		t.Fatal("expected non-nil Qualifier for d.$1")
	}
	if dr.Qualifier.Name.Name != "d" {
		t.Errorf("Qualifier.Name = %q, want %q", dr.Qualifier.Name.Name, "d")
	}
	// Loc spans from the qualifier start to the $-token end.
	if dr.Loc.Start != 0 || dr.Loc.End != 4 {
		t.Errorf("Loc = %+v, want {0 4}", dr.Loc)
	}
}

func TestExpr_DollarQualifiedSchemaTable(t *testing.T) {
	// s.t.$1 — two-part qualifier on a positional ref.
	node, err := testParseExpr("s.t.$1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dr, ok := node.(*ast.DollarRef)
	if !ok {
		t.Fatalf("expected *ast.DollarRef, got %T", node)
	}
	if dr.Qualifier == nil {
		t.Fatal("expected non-nil Qualifier")
	}
	if dr.Qualifier.Schema.Name != "s" || dr.Qualifier.Name.Name != "t" {
		t.Errorf("Qualifier = %s.%s, want s.t", dr.Qualifier.Schema.Name, dr.Qualifier.Name.Name)
	}
}

func TestExpr_DollarColonCast(t *testing.T) {
	// $1:num::number — positional ref, semi-structured colon access, then cast.
	// Composes through the infix loop: CastExpr(AccessExpr(Colon, DollarRef, num)).
	node, err := testParseExpr("$1:num::number")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cast, ok := node.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected *ast.CastExpr, got %T", node)
	}
	if !cast.ColonColon {
		t.Error("ColonColon should be true")
	}
	acc, ok := cast.Expr.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("expected *ast.AccessExpr, got %T", cast.Expr)
	}
	if acc.Kind != ast.AccessColon || acc.Field.Name != "num" {
		t.Errorf("access = %v %q, want Colon num", acc.Kind, acc.Field.Name)
	}
	dr, ok := acc.Expr.(*ast.DollarRef)
	if !ok {
		t.Fatalf("expected *ast.DollarRef, got %T", acc.Expr)
	}
	if dr.Name != "1" || !dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"1\" true", dr.Name, dr.Positional)
	}
}

func TestExpr_DollarArithmetic(t *testing.T) {
	// $1 + $2 — DollarRefs compose as arithmetic operands.
	node, err := testParseExpr("$1 + $2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	if bin.Op != ast.BinAdd {
		t.Errorf("Op = %v, want BinAdd", bin.Op)
	}
	if _, ok := bin.Left.(*ast.DollarRef); !ok {
		t.Errorf("Left = %T, want *ast.DollarRef", bin.Left)
	}
	if _, ok := bin.Right.(*ast.DollarRef); !ok {
		t.Errorf("Right = %T, want *ast.DollarRef", bin.Right)
	}
}

func TestExpr_DollarInArithmeticWithVar(t *testing.T) {
	// 2 * $min — the SET corpus form (set/example_06.sql value).
	node, err := testParseExpr("2 * $min")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", node)
	}
	dr, ok := bin.Right.(*ast.DollarRef)
	if !ok {
		t.Fatalf("Right = %T, want *ast.DollarRef", bin.Right)
	}
	if dr.Name != "min" || dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"min\" false", dr.Name, dr.Positional)
	}
}

func TestExpr_DollarAsFuncArg(t *testing.T) {
	// TO_DATE($1) — positional ref as a function argument (copy-into-location).
	node, err := testParseExpr("TO_DATE($1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(fc.Args))
	}
	dr, ok := fc.Args[0].(*ast.DollarRef)
	if !ok {
		t.Fatalf("Args[0] = %T, want *ast.DollarRef", fc.Args[0])
	}
	if dr.Name != "1" || !dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"1\" true", dr.Name, dr.Positional)
	}
}

func TestExpr_DollarVarAsFuncArg(t *testing.T) {
	// p($v) argument form — the execCallDollarLimited corpus case.
	node, err := testParseExpr("p($v)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected *ast.FuncCallExpr, got %T", node)
	}
	if len(fc.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(fc.Args))
	}
	dr, ok := fc.Args[0].(*ast.DollarRef)
	if !ok {
		t.Fatalf("Args[0] = %T, want *ast.DollarRef", fc.Args[0])
	}
	if dr.Name != "v" || dr.Positional {
		t.Errorf("got Name=%q Positional=%v, want \"v\" false", dr.Name, dr.Positional)
	}
}

func TestExpr_DollarBracketAccess(t *testing.T) {
	// $1[0] — positional ref with array subscript (semi-structured).
	node, err := testParseExpr("$1[0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	acc, ok := node.(*ast.AccessExpr)
	if !ok {
		t.Fatalf("expected *ast.AccessExpr, got %T", node)
	}
	if acc.Kind != ast.AccessBracket {
		t.Errorf("Kind = %v, want AccessBracket", acc.Kind)
	}
	if _, ok := acc.Expr.(*ast.DollarRef); !ok {
		t.Errorf("Expr = %T, want *ast.DollarRef", acc.Expr)
	}
}

// Negative: a lone '$' with no following name/number is not a DollarRef — the
// lexer emits a bare '$' ASCII token which the expression parser must reject.
func TestExpr_DollarBareIsError(t *testing.T) {
	_, err := testParseExpr("$")
	if err == nil {
		t.Fatal("expected an error for a lone '$'")
	}
}

// Negative: '$' followed by whitespace then a number is still a lone '$' (the
// digits are a separate token) and must error.
func TestExpr_DollarSpaceNumberIsError(t *testing.T) {
	_, err := testParseExpr("$ 1")
	if err == nil {
		t.Fatal("expected an error for '$ 1' (bare $ followed by separate number)")
	}
}

// ---------------------------------------------------------------------------
// Named (keyword) function arguments:  f(name => value)
// ---------------------------------------------------------------------------

// namedArg fetches the i-th argument of a function call, asserting it is a
// named argument (*ast.CallArg with a non-zero Name) and returning it.
func namedArg(t *testing.T, fc *ast.FuncCallExpr, i int) *ast.CallArg {
	t.Helper()
	if i >= len(fc.Args) {
		t.Fatalf("arg index %d out of range (len=%d)", i, len(fc.Args))
	}
	ca, ok := fc.Args[i].(*ast.CallArg)
	if !ok {
		t.Fatalf("arg %d type = %T, want *ast.CallArg", i, fc.Args[i])
	}
	if ca.Name.Name == "" {
		t.Fatalf("arg %d has empty Name, want a named argument", i)
	}
	return ca
}

func mustFuncCall(t *testing.T, input string) *ast.FuncCallExpr {
	t.Helper()
	node, err := testParseExpr(input)
	if err != nil {
		t.Fatalf("parse %q: unexpected error: %v", input, err)
	}
	fc, ok := node.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.FuncCallExpr", input, node)
	}
	return fc
}

func TestExpr_NamedArgSingle(t *testing.T) {
	fc := mustFuncCall(t, "INFER_SCHEMA(LOCATION => '@stage')")
	if len(fc.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(fc.Args))
	}
	ca := namedArg(t, fc, 0)
	if ca.Name.Name != "LOCATION" {
		t.Errorf("Name = %q, want LOCATION", ca.Name.Name)
	}
	lit, ok := ca.Value.(*ast.Literal)
	if !ok || lit.Kind != ast.LitString || lit.Value != "@stage" {
		t.Errorf("Value = %#v, want string literal '@stage'", ca.Value)
	}
}

func TestExpr_NamedArgNoSpaces(t *testing.T) {
	// LOCATION=>'@s' with no surrounding whitespace (=> is its own token).
	fc := mustFuncCall(t, "INFER_SCHEMA(LOCATION=>'@s', FILE_FORMAT=>'fmt')")
	if len(fc.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(fc.Args))
	}
	if a := namedArg(t, fc, 0); a.Name.Name != "LOCATION" {
		t.Errorf("arg0 Name = %q, want LOCATION", a.Name.Name)
	}
	if a := namedArg(t, fc, 1); a.Name.Name != "FILE_FORMAT" {
		t.Errorf("arg1 Name = %q, want FILE_FORMAT", a.Name.Name)
	}
}

func TestExpr_NamedArgMultiple(t *testing.T) {
	fc := mustFuncCall(t, "FLATTEN(INPUT => a.b, PATH => 'contact', OUTER => TRUE)")
	if len(fc.Args) != 3 {
		t.Fatalf("len(Args) = %d, want 3", len(fc.Args))
	}
	for i, want := range []string{"INPUT", "PATH", "OUTER"} {
		if a := namedArg(t, fc, i); a.Name.Name != want {
			t.Errorf("arg%d Name = %q, want %q", i, a.Name.Name, want)
		}
	}
}

func TestExpr_NamedArgMixedPositionalThenNamed(t *testing.T) {
	// Positional args first, then named — the common Snowflake form.
	fc := mustFuncCall(t, "f(1, x, TYPE => 'STREAMING')")
	if len(fc.Args) != 3 {
		t.Fatalf("len(Args) = %d, want 3", len(fc.Args))
	}
	if _, ok := fc.Args[0].(*ast.CallArg); ok {
		t.Errorf("arg0 should be positional, got *ast.CallArg")
	}
	if _, ok := fc.Args[1].(*ast.CallArg); ok {
		t.Errorf("arg1 should be positional, got *ast.CallArg")
	}
	if a := namedArg(t, fc, 2); a.Name.Name != "TYPE" {
		t.Errorf("arg2 Name = %q, want TYPE", a.Name.Name)
	}
}

func TestExpr_NamedArgInterleaved(t *testing.T) {
	// Lenient: named, then positional, then named (Snowflake itself is stricter,
	// but the parser accepts interleaving).
	fc := mustFuncCall(t, "f(A => 1, 2, B => 3)")
	if len(fc.Args) != 3 {
		t.Fatalf("len(Args) = %d, want 3", len(fc.Args))
	}
	if a := namedArg(t, fc, 0); a.Name.Name != "A" {
		t.Errorf("arg0 Name = %q, want A", a.Name.Name)
	}
	if _, ok := fc.Args[1].(*ast.CallArg); ok {
		t.Errorf("arg1 should be positional")
	}
	if a := namedArg(t, fc, 2); a.Name.Name != "B" {
		t.Errorf("arg2 Name = %q, want B", a.Name.Name)
	}
}

func TestExpr_NamedArgNested(t *testing.T) {
	// A function call whose single argument is itself a named-arg call:
	// WRAP(DATA_SOURCE(TYPE => 'STREAMING')). (TABLE is a reserved word and only
	// legal as a FROM table-function, so the bare-expression nesting uses a
	// non-reserved outer name; the TABLE(DATA_SOURCE(...)) statement form is
	// covered by TestStmt_NamedArgInTableFunction.)
	fc := mustFuncCall(t, "WRAP(DATA_SOURCE(TYPE => 'STREAMING'))")
	if len(fc.Args) != 1 {
		t.Fatalf("outer len(Args) = %d, want 1", len(fc.Args))
	}
	inner, ok := fc.Args[0].(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("outer arg0 = %T, want *ast.FuncCallExpr (DATA_SOURCE)", fc.Args[0])
	}
	if inner.Name.Name.Name != "DATA_SOURCE" {
		t.Errorf("inner func = %q, want DATA_SOURCE", inner.Name.Name.Name)
	}
	ca := namedArg(t, inner, 0)
	if ca.Name.Name != "TYPE" {
		t.Errorf("inner arg Name = %q, want TYPE", ca.Name.Name)
	}
	if lit, ok := ca.Value.(*ast.Literal); !ok || lit.Value != "STREAMING" {
		t.Errorf("inner arg value = %#v, want literal 'STREAMING'", ca.Value)
	}
}

// TestStmt_NamedArgInTableFunction exercises the literal TABLE(DATA_SOURCE(TYPE
// => 'STREAMING')) table-function form (pipe / COPY streaming source) at the
// statement level, where TABLE() is a legal FROM table function.
func TestStmt_NamedArgInTableFunction(t *testing.T) {
	sql := "SELECT $1 FROM TABLE(DATA_SOURCE(TYPE => 'STREAMING'))"
	result, errs := parseSingle(sql, 0)
	if len(errs) > 0 {
		t.Fatalf("parse %q: %v", sql, errs)
	}
	// Walk the AST and confirm a DATA_SOURCE call with a named TYPE argument
	// appears.
	found := false
	ast.Inspect(result, func(n ast.Node) bool {
		fc, ok := n.(*ast.FuncCallExpr)
		if !ok || fc.Name.Name.Name != "DATA_SOURCE" {
			return true
		}
		if len(fc.Args) == 1 {
			if ca, ok := fc.Args[0].(*ast.CallArg); ok && ca.Name.Name == "TYPE" {
				found = true
			}
		}
		return true
	})
	if !found {
		t.Errorf("did not find DATA_SOURCE(TYPE => ...) named arg in %q", sql)
	}
}

func TestExpr_NamedArgValueIsComplexExpr(t *testing.T) {
	// The value of a named arg may be any expression: a function call, a cast,
	// a $N ref, and a colon path all compose through.
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, v ast.Node)
	}{
		{
			name:  "func-call value",
			input: "CLONE_AT(TIMESTAMP => TO_TIMESTAMP_TZ('x', 'fmt'))",
			check: func(t *testing.T, v ast.Node) {
				if _, ok := v.(*ast.FuncCallExpr); !ok {
					t.Errorf("value = %T, want *ast.FuncCallExpr", v)
				}
			},
		},
		{
			name:  "cast value",
			input: "f(X => $1::number)",
			check: func(t *testing.T, v ast.Node) {
				if _, ok := v.(*ast.CastExpr); !ok {
					t.Errorf("value = %T, want *ast.CastExpr", v)
				}
			},
		},
		{
			name:  "dollar-ref value",
			input: "f(X => $1)",
			check: func(t *testing.T, v ast.Node) {
				if _, ok := v.(*ast.DollarRef); !ok {
					t.Errorf("value = %T, want *ast.DollarRef", v)
				}
			},
		},
		{
			name:  "colon-path value",
			input: "f(X => col:field)",
			check: func(t *testing.T, v ast.Node) {
				if v == nil {
					t.Errorf("value is nil")
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fc := mustFuncCall(t, tc.input)
			ca := namedArg(t, fc, 0)
			tc.check(t, ca.Value)
		})
	}
}

func TestExpr_NamedArgLoc(t *testing.T) {
	// Loc must span from the name's start to the value's end.
	input := "f(NAME => 123)"
	fc := mustFuncCall(t, input)
	ca := namedArg(t, fc, 0)
	// "NAME" starts at offset 2; the value "123" ends at offset 13.
	if ca.Loc.Start != 2 {
		t.Errorf("Loc.Start = %d, want 2 (start of NAME)", ca.Loc.Start)
	}
	if ca.Loc.End != 13 {
		t.Errorf("Loc.End = %d, want 13 (end of 123)", ca.Loc.End)
	}
}

// ---------------------------------------------------------------------------
// Named-arg negatives + regression
// ---------------------------------------------------------------------------

func TestExpr_NamedArgBareAssocIsError(t *testing.T) {
	// `f(=> 1)` — `=>` with no name must error cleanly, not loop.
	if _, err := testParseExpr("f(=> 1)"); err == nil {
		t.Fatal("expected an error for f(=> 1)")
	}
}

func TestExpr_NamedArgMissingValueIsError(t *testing.T) {
	// `f(x =>)` — name with no value must error cleanly.
	if _, err := testParseExpr("f(x =>)"); err == nil {
		t.Fatal("expected an error for f(x =>)")
	}
}

func TestExpr_NamedArgMissingValueBeforeCommaIsError(t *testing.T) {
	// `f(x =>, y)` — the value is missing before the comma.
	if _, err := testParseExpr("f(x =>, y)"); err == nil {
		t.Fatal("expected an error for f(x =>, y)")
	}
}

// Regression: an ordinary positional call whose argument is a bare identifier
// must NOT be mistaken for a named argument (no `=>` follows).
func TestExpr_PositionalIdentArgsUnchanged(t *testing.T) {
	fc := mustFuncCall(t, "f(a, b)")
	if len(fc.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(fc.Args))
	}
	for i := range fc.Args {
		if _, ok := fc.Args[i].(*ast.CallArg); ok {
			t.Errorf("arg %d wrongly parsed as a named *ast.CallArg", i)
		}
	}
}

// Regression: the `->` (tokArrow) lambda/operator is distinct from `=>` and must
// keep parsing unchanged inside a higher-order function argument.
func TestExpr_ArrowOperatorUnchanged(t *testing.T) {
	// TRANSFORM(arr, x -> x + 1): the second argument is a lambda using ->, not
	// a named argument. It must parse as a LambdaExpr, never a *ast.CallArg.
	fc := mustFuncCall(t, "TRANSFORM(arr, x -> x + 1)")
	if len(fc.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(fc.Args))
	}
	if _, ok := fc.Args[1].(*ast.CallArg); ok {
		t.Errorf("arg1 (x -> x + 1) wrongly parsed as a named *ast.CallArg")
	}
	if _, ok := fc.Args[1].(*ast.LambdaExpr); !ok {
		t.Errorf("arg1 = %T, want *ast.LambdaExpr", fc.Args[1])
	}
}

// Regression: a comparison `a >= b` must not be confused with a named arg — its
// operator is its own token, distinct from tokAssoc. (Sanity at the arg-list
// level; a plain function is used because IFF is a dedicated AST node.)
func TestExpr_GreaterEqualArgUnchanged(t *testing.T) {
	fc := mustFuncCall(t, "f(a >= b, 1, 0)")
	if len(fc.Args) != 3 {
		t.Fatalf("len(Args) = %d, want 3", len(fc.Args))
	}
	if _, ok := fc.Args[0].(*ast.CallArg); ok {
		t.Errorf("arg0 (a >= b) wrongly parsed as a named *ast.CallArg")
	}
	if _, ok := fc.Args[0].(*ast.BinaryExpr); !ok {
		t.Errorf("arg0 = %T, want *ast.BinaryExpr (a >= b)", fc.Args[0])
	}
}
