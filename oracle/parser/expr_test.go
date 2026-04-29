package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseExprLiterals tests parsing of literal expressions.
func TestParseExprLiterals(t *testing.T) {
	tests := []struct {
		input string
		check func(t *testing.T, e ast.ExprNode)
	}{
		{"42", func(t *testing.T, e ast.ExprNode) {
			n, ok := e.(*ast.NumberLiteral)
			if !ok {
				t.Fatalf("expected NumberLiteral, got %T", e)
			}
			if n.Ival != 42 {
				t.Errorf("expected 42, got %d", n.Ival)
			}
		}},
		{"3.14", func(t *testing.T, e ast.ExprNode) {
			n, ok := e.(*ast.NumberLiteral)
			if !ok {
				t.Fatalf("expected NumberLiteral, got %T", e)
			}
			if !n.IsFloat {
				t.Error("expected IsFloat=true")
			}
		}},
		{"'hello'", func(t *testing.T, e ast.ExprNode) {
			n, ok := e.(*ast.StringLiteral)
			if !ok {
				t.Fatalf("expected StringLiteral, got %T", e)
			}
			if n.Val != "hello" {
				t.Errorf("expected hello, got %q", n.Val)
			}
		}},
		{"NULL", func(t *testing.T, e ast.ExprNode) {
			if _, ok := e.(*ast.NullLiteral); !ok {
				t.Fatalf("expected NullLiteral, got %T", e)
			}
		}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			e, parseErr1 := p.parseExpr()
			if parseErr1 != nil {
				t.Fatalf("parse: %v", parseErr1)
			}
			if e == nil {
				t.Fatal("expected non-nil expression")
			}
			tc.check(t, e)
		})
	}
}

// TestParseExprBinary tests binary operator expressions.
func TestParseExprBinary(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{"1 + 2", "+"},
		{"1 - 2", "-"},
		{"3 * 4", "*"},
		{"6 / 3", "/"},
		{"2 ** 3", "**"},
		{"'a' || 'b'", "||"},
		{"1 = 2", "="},
		{"1 != 2", "!="},
		{"1 <> 2", "<>"},
		{"1 < 2", "<"},
		{"1 > 2", ">"},
		{"1 <= 2", "<="},
		{"1 >= 2", ">="},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			e, parseErr2 := p.parseExpr()
			if parseErr2 != nil {
				t.Fatalf("parse: %v", parseErr2)
			}
			be, ok := e.(*ast.BinaryExpr)
			if !ok {
				t.Fatalf("expected BinaryExpr, got %T", e)
			}
			if be.Op != tc.op {
				t.Errorf("expected op %q, got %q", tc.op, be.Op)
			}
		})
	}
}

// TestParseExprPrecedence tests operator precedence.
func TestParseExprPrecedence(t *testing.T) {
	// 1 + 2 * 3 should parse as 1 + (2 * 3)
	p := newTestParser("1 + 2 * 3")
	e, parseErr1 := p.parseExpr()
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", e)
	}
	if be.Op != "+" {
		t.Errorf("expected top-level op +, got %q", be.Op)
	}
	right, ok := be.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be BinaryExpr, got %T", be.Right)
	}
	if right.Op != "*" {
		t.Errorf("expected right op *, got %q", right.Op)
	}
}

// TestParseExprUnary tests unary expressions.
func TestParseExprUnary(t *testing.T) {
	p := newTestParser("-42")
	e, parseErr2 := p.parseExpr()
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	ue, ok := e.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", e)
	}
	if ue.Op != "-" {
		t.Errorf("expected op -, got %q", ue.Op)
	}
}

// TestParseExprBoolOps tests AND, OR, NOT.
func TestParseExprBoolOps(t *testing.T) {
	p := newTestParser("1 = 1 AND 2 = 2")
	e, parseErr3 := p.parseExpr()
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	be, ok := e.(*ast.BoolExpr)
	if !ok {
		t.Fatalf("expected BoolExpr, got %T", e)
	}
	if be.Boolop != ast.BOOL_AND {
		t.Errorf("expected BOOL_AND, got %d", be.Boolop)
	}
}

// TestParseExprNot tests NOT expression.
func TestParseExprNot(t *testing.T) {
	p := newTestParser("NOT 1 = 1")
	e, parseErr4 := p.parseExpr()
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	be, ok := e.(*ast.BoolExpr)
	if !ok {
		t.Fatalf("expected BoolExpr, got %T", e)
	}
	if be.Boolop != ast.BOOL_NOT {
		t.Errorf("expected BOOL_NOT, got %d", be.Boolop)
	}
}

// TestParseExprCase tests CASE expression.
func TestParseExprCase(t *testing.T) {
	// Simple CASE
	p := newTestParser("CASE x WHEN 1 THEN 'a' WHEN 2 THEN 'b' ELSE 'c' END")
	e, parseErr5 := p.parseExpr()
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	ce, ok := e.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected CaseExpr, got %T", e)
	}
	if ce.Arg == nil {
		t.Error("expected non-nil Arg for simple CASE")
	}
	if ce.Whens.Len() != 2 {
		t.Errorf("expected 2 WHEN clauses, got %d", ce.Whens.Len())
	}
	if ce.Default == nil {
		t.Error("expected non-nil Default")
	}
}

// TestParseExprSearchedCase tests searched CASE expression.
func TestParseExprSearchedCase(t *testing.T) {
	p := newTestParser("CASE WHEN x > 0 THEN 'pos' WHEN x < 0 THEN 'neg' END")
	e, parseErr6 := p.parseExpr()
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	ce, ok := e.(*ast.CaseExpr)
	if !ok {
		t.Fatalf("expected CaseExpr, got %T", e)
	}
	if ce.Arg != nil {
		t.Error("expected nil Arg for searched CASE")
	}
	if ce.Whens.Len() != 2 {
		t.Errorf("expected 2 WHEN clauses, got %d", ce.Whens.Len())
	}
}

// TestParseExprDecode tests DECODE expression.
func TestParseExprDecode(t *testing.T) {
	p := newTestParser("DECODE(status, 'A', 'Active', 'I', 'Inactive', 'Unknown')")
	e, parseErr7 := p.parseExpr()
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	de, ok := e.(*ast.DecodeExpr)
	if !ok {
		t.Fatalf("expected DecodeExpr, got %T", e)
	}
	if de.Arg == nil {
		t.Error("expected non-nil Arg")
	}
	if de.Pairs.Len() != 2 {
		t.Errorf("expected 2 pairs, got %d", de.Pairs.Len())
	}
	if de.Default == nil {
		t.Error("expected non-nil Default")
	}
}

// TestParseExprBetween tests BETWEEN expression.
func TestParseExprBetween(t *testing.T) {
	p := newTestParser("x BETWEEN 1 AND 10")
	e, parseErr8 := p.parseExpr()
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	be, ok := e.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected BetweenExpr, got %T", e)
	}
	if be.Not {
		t.Error("expected Not=false")
	}
}

// TestParseExprNotBetween tests NOT BETWEEN expression.
func TestParseExprNotBetween(t *testing.T) {
	p := newTestParser("x NOT BETWEEN 1 AND 10")
	e, parseErr9 := p.parseExpr()
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	be, ok := e.(*ast.BetweenExpr)
	if !ok {
		t.Fatalf("expected BetweenExpr, got %T", e)
	}
	if !be.Not {
		t.Error("expected Not=true")
	}
}

// TestParseExprIn tests IN expression.
func TestParseExprIn(t *testing.T) {
	p := newTestParser("x IN (1, 2, 3)")
	e, parseErr10 := p.parseExpr()
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	ie, ok := e.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected InExpr, got %T", e)
	}
	if ie.Not {
		t.Error("expected Not=false")
	}
	if ie.List.Len() != 3 {
		t.Errorf("expected 3 items, got %d", ie.List.Len())
	}
}

// TestParseExprNotIn tests NOT IN expression.
func TestParseExprNotIn(t *testing.T) {
	p := newTestParser("x NOT IN (1, 2)")
	e, parseErr11 := p.parseExpr()
	if parseErr11 != nil {
		t.Fatalf("parse: %v", parseErr11)
	}
	ie, ok := e.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected InExpr, got %T", e)
	}
	if !ie.Not {
		t.Error("expected Not=true")
	}
}

// TestParseExprLike tests LIKE expression.
func TestParseExprLike(t *testing.T) {
	p := newTestParser("name LIKE '%test%'")
	e, parseErr12 := p.parseExpr()
	if parseErr12 != nil {
		t.Fatalf("parse: %v", parseErr12)
	}
	le, ok := e.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected LikeExpr, got %T", e)
	}
	if le.Not {
		t.Error("expected Not=false")
	}
	if le.Type != ast.LIKE_LIKE {
		t.Errorf("expected LIKE_LIKE, got %d", le.Type)
	}
}

// TestParseExprLikeEscape tests LIKE with ESCAPE.
func TestParseExprLikeEscape(t *testing.T) {
	p := newTestParser("name LIKE '%\\%%' ESCAPE '\\'")
	e, parseErr13 := p.parseExpr()
	if parseErr13 != nil {
		t.Fatalf("parse: %v", parseErr13)
	}
	le, ok := e.(*ast.LikeExpr)
	if !ok {
		t.Fatalf("expected LikeExpr, got %T", e)
	}
	if le.Escape == nil {
		t.Error("expected non-nil Escape")
	}
}

// TestParseExprIsNull tests IS NULL and IS NOT NULL.
func TestParseExprIsNull(t *testing.T) {
	tests := []struct {
		input string
		not   bool
		test  string
	}{
		{"x IS NULL", false, "NULL"},
		{"x IS NOT NULL", true, "NULL"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			e, parseErr3 := p.parseExpr()
			if parseErr3 != nil {
				t.Fatalf("parse: %v", parseErr3)
			}
			ie, ok := e.(*ast.IsExpr)
			if !ok {
				t.Fatalf("expected IsExpr, got %T", e)
			}
			if ie.Not != tc.not {
				t.Errorf("expected Not=%v, got %v", tc.not, ie.Not)
			}
			if ie.Test != tc.test {
				t.Errorf("expected Test=%q, got %q", tc.test, ie.Test)
			}
		})
	}
}

// TestParseExprCast tests CAST expression.
func TestParseExprCast(t *testing.T) {
	p := newTestParser("CAST(salary AS NUMBER(10,2))")
	e, parseErr14 := p.parseExpr()
	if parseErr14 != nil {
		t.Fatalf("parse: %v", parseErr14)
	}
	ce, ok := e.(*ast.CastExpr)
	if !ok {
		t.Fatalf("expected CastExpr, got %T", e)
	}
	if ce.Arg == nil {
		t.Error("expected non-nil Arg")
	}
	if ce.TypeName == nil {
		t.Error("expected non-nil TypeName")
	}
}

// TestParseExprFuncCall tests function call expression.
func TestParseExprFuncCall(t *testing.T) {
	p := newTestParser("NVL(salary, 0)")
	e, parseErr15 := p.parseExpr()
	if parseErr15 != nil {
		t.Fatalf("parse: %v", parseErr15)
	}
	fc, ok := e.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected FuncCallExpr, got %T", e)
	}
	if fc.FuncName.Name != "NVL" {
		t.Errorf("expected NVL, got %q", fc.FuncName.Name)
	}
	if fc.Args.Len() != 2 {
		t.Errorf("expected 2 args, got %d", fc.Args.Len())
	}
}

// TestParseExprCountStar tests COUNT(*).
func TestParseExprCountStar(t *testing.T) {
	p := newTestParser("COUNT(*)")
	e, parseErr16 := p.parseExpr()
	if parseErr16 != nil {
		t.Fatalf("parse: %v", parseErr16)
	}
	fc, ok := e.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected FuncCallExpr, got %T", e)
	}
	if !fc.Star {
		t.Error("expected Star=true")
	}
}

// TestParseExprCountDistinct tests COUNT(DISTINCT col).
func TestParseExprCountDistinct(t *testing.T) {
	p := newTestParser("COUNT(DISTINCT status)")
	e, parseErr17 := p.parseExpr()
	if parseErr17 != nil {
		t.Fatalf("parse: %v", parseErr17)
	}
	fc, ok := e.(*ast.FuncCallExpr)
	if !ok {
		t.Fatalf("expected FuncCallExpr, got %T", e)
	}
	if !fc.Distinct {
		t.Error("expected Distinct=true")
	}
}

// TestParseExprParens tests parenthesized expressions.
func TestParseExprParens(t *testing.T) {
	p := newTestParser("(1 + 2) * 3")
	e, parseErr18 := p.parseExpr()
	if parseErr18 != nil {
		t.Fatalf("parse: %v", parseErr18)
	}
	be, ok := e.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", e)
	}
	if be.Op != "*" {
		t.Errorf("expected top-level op *, got %q", be.Op)
	}
}

// TestParseExprBindVar tests bind variable in expressions.
func TestParseExprBindVar(t *testing.T) {
	p := newTestParser(":param1")
	e, parseErr19 := p.parseExpr()
	if parseErr19 != nil {
		t.Fatalf("parse: %v", parseErr19)
	}
	bv, ok := e.(*ast.BindVariable)
	if !ok {
		t.Fatalf("expected BindVariable, got %T", e)
	}
	if bv.Name != "param1" {
		t.Errorf("expected param1, got %q", bv.Name)
	}
}

// TestParseExprColumnRef tests column reference in expressions.
func TestParseExprColumnRef(t *testing.T) {
	p := newTestParser("e.salary")
	e, parseErr20 := p.parseExpr()
	if parseErr20 != nil {
		t.Fatalf("parse: %v", parseErr20)
	}
	cr, ok := e.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected ColumnRef, got %T", e)
	}
	if cr.Table != "E" {
		t.Errorf("expected table E, got %q", cr.Table)
	}
	if cr.Column != "SALARY" {
		t.Errorf("expected column SALARY, got %q", cr.Column)
	}
}

// TestParseExprOrPrecedence tests OR has lower precedence than AND.
func TestParseExprOrPrecedence(t *testing.T) {
	// a = 1 OR b = 2 AND c = 3 should be a = 1 OR (b = 2 AND c = 3)
	p := newTestParser("a = 1 OR b = 2 AND c = 3")
	e, parseErr21 := p.parseExpr()
	if parseErr21 != nil {
		t.Fatalf("parse: %v", parseErr21)
	}
	be, ok := e.(*ast.BoolExpr)
	if !ok {
		t.Fatalf("expected BoolExpr, got %T", e)
	}
	if be.Boolop != ast.BOOL_OR {
		t.Errorf("expected BOOL_OR, got %d", be.Boolop)
	}
}

// TestParseExprPrior tests PRIOR unary operator.
func TestParseExprPrior(t *testing.T) {
	p := newTestParser("PRIOR employee_id")
	e, parseErr22 := p.parseExpr()
	if parseErr22 != nil {
		t.Fatalf("parse: %v", parseErr22)
	}
	ue, ok := e.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", e)
	}
	if ue.Op != "PRIOR" {
		t.Errorf("expected op PRIOR, got %q", ue.Op)
	}
}

// TestParseExprStar tests standalone star expression.
func TestParseExprStar(t *testing.T) {
	p := newTestParser("*")
	e, parseErr23 := p.parseExpr()
	if parseErr23 != nil {
		t.Fatalf("parse: %v", parseErr23)
	}
	_, ok := e.(*ast.Star)
	if !ok {
		t.Fatalf("expected Star, got %T", e)
	}
}
