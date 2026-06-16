package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// Binary/hex literals (X'..', B'..') and the BINARY unary operator.

func TestHexLiteral(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT X'4142' FROM t")
	lit, ok := stmt.Items[0].Expr.(*ast.Literal)
	if !ok {
		t.Fatalf("item expr = %T, want *ast.Literal", stmt.Items[0].Expr)
	}
	if lit.Kind != ast.LitHex {
		t.Errorf("Kind = %v, want LitHex", lit.Kind)
	}
}

func TestBitLiteral(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT B'101' FROM t")
	lit, ok := stmt.Items[0].Expr.(*ast.Literal)
	if !ok {
		t.Fatalf("item expr = %T, want *ast.Literal", stmt.Items[0].Expr)
	}
	if lit.Kind != ast.LitBit {
		t.Errorf("Kind = %v, want LitBit", lit.Kind)
	}
}

func TestBinaryOperator(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT BINARY col1 FROM t")
	un, ok := stmt.Items[0].Expr.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("item expr = %T, want *ast.UnaryExpr", stmt.Items[0].Expr)
	}
	if un.Op != ast.UnaryBinary {
		t.Errorf("Op = %v, want UnaryBinary", un.Op)
	}
	if _, ok := un.Expr.(*ast.ColumnRef); !ok {
		t.Errorf("operand = %T, want *ast.ColumnRef", un.Expr)
	}
}

// TestBinaryOperatorInWhere confirms BINARY binds at the lowest precedence:
// BINARY name = 'abc' parses as BINARY(name = 'abc').
func TestBinaryOperatorInWhere(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM t WHERE BINARY name = 'abc'")
	un, ok := stmt.Where.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("where = %T, want *ast.UnaryExpr", stmt.Where)
	}
	if un.Op != ast.UnaryBinary {
		t.Errorf("Op = %v, want UnaryBinary", un.Op)
	}
	if _, ok := un.Expr.(*ast.BinaryExpr); !ok {
		t.Errorf("operand = %T, want *ast.BinaryExpr (name = 'abc')", un.Expr)
	}
}

// TestHexInteger confirms the already-supported 0x form still parses.
func TestHexInteger(t *testing.T) {
	_ = mustParseSelect(t, "SELECT 0x4142 FROM t")
}
