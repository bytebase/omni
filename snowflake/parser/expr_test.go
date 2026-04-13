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
