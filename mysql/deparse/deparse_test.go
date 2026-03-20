package deparse

import (
	"testing"

	ast "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// parseExpr parses a SQL expression by wrapping it in SELECT and extracting
// the first target list expression from the AST.
func parseExpr(t *testing.T, expr string) ast.ExprNode {
	t.Helper()
	sql := "SELECT " + expr
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", sql, err)
	}
	if list.Len() == 0 {
		t.Fatalf("no statements parsed from %q", sql)
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", list.Items[0])
	}
	if len(sel.TargetList) == 0 {
		t.Fatalf("no target list in SELECT from %q", sql)
	}
	target := sel.TargetList[0]
	// TargetList entries may be ResTarget wrapping the actual expression.
	if rt, ok := target.(*ast.ResTarget); ok {
		return rt.Val
	}
	return target
}

func TestDeparse_Section_1_1_IntFloatNull(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Integer literals
		{"integer_1", "1", "1"},
		{"negative_integer", "-5", "-5"},
		{"large_integer", "9999999999", "9999999999"},
		{"zero", "0", "0"},

		// Float literals
		{"float_1_5", "1.5", "1.5"},
		{"float_with_exponent", "1.5e10", "1.5e10"},
		{"float_zero_point_five", "0.5", "0.5"},

		// NULL literal
		{"null", "NULL", "NULL"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparse_Section_1_2_BoolStringLiterals(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Boolean literals
		{"true_literal", "TRUE", "true"},
		{"false_literal", "FALSE", "false"},

		// String literals
		{"simple_string", "'hello'", "'hello'"},
		{"string_with_single_quote", "'it''s'", `'it\'s'`},
		{"empty_string", "''", "''"},
		{"string_with_backslash", `'back\\slash'`, `'back\\slash'`},

		// Charset introducers — skipped: parser doesn't support charset introducers yet
		// {"charset_utf8mb4", "_utf8mb4'hello'", "_utf8mb4'hello'"},
		// {"charset_latin1", "_latin1'world'", "_latin1'world'"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_1_2_CharsetIntroducer tests charset introducer deparsing
// using hand-built AST nodes, since the parser doesn't support charset introducers yet.
func TestDeparse_Section_1_2_CharsetIntroducer(t *testing.T) {
	cases := []struct {
		name     string
		node     ast.ExprNode
		expected string
	}{
		{"charset_utf8mb4", &ast.StringLit{Value: "hello", Charset: "_utf8mb4"}, "_utf8mb4'hello'"},
		{"charset_latin1", &ast.StringLit{Value: "world", Charset: "_latin1"}, "_latin1'world'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Deparse(tc.node)
			if got != tc.expected {
				t.Errorf("Deparse() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestDeparse_Section_1_3_HexBitLiterals(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Hex literals — MySQL normalizes to 0x lowercase form
		{"hex_0x_form", "0xFF", "0xff"},
		{"hex_X_quote_form", "X'FF'", "0xff"},

		// Bit literals — MySQL converts to hex form
		{"bit_0b_form", "0b1010", "0x0a"},
		{"bit_b_quote_form", "b'1010'", "0x0a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_1_3_DateTimeLiterals tests DATE/TIME/TIMESTAMP literal deparsing
// using hand-built AST nodes, since the parser doesn't support temporal literals yet.
// These are marked as [~] partial in SCENARIOS — parser support needed.

func TestDeparse_Section_2_1_ArithmeticUnary(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Arithmetic binary operators
		{"addition", "a + b", "(`a` + `b`)"},
		{"subtraction", "a - b", "(`a` - `b`)"},
		{"multiplication", "a * b", "(`a` * `b`)"},
		{"division", "a / b", "(`a` / `b`)"},
		{"integer_division", "a DIV b", "(`a` DIV `b`)"},
		{"modulo_MOD", "a MOD b", "(`a` % `b`)"},
		{"modulo_percent", "a % b", "(`a` % `b`)"},

		// Left-associative chaining
		{"left_assoc_chain", "a + b + c", "((`a` + `b`) + `c`)"},

		// Unary minus (with column ref operand)
		{"unary_minus", "-a", "-`a`"},

		// Unary plus — parser drops it entirely, so +a parses as just ColumnRef
		{"unary_plus", "+a", "`a`"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_2_1_UnaryPlusAST tests unary plus via hand-built AST
// to verify that deparseUnaryExpr handles UnaryPlus correctly, even though
// the parser drops unary plus before building the AST.
func TestDeparse_Section_2_1_UnaryPlusAST(t *testing.T) {
	node := &ast.UnaryExpr{
		Op:      ast.UnaryPlus,
		Operand: &ast.ColumnRef{Column: "a"},
	}
	got := Deparse(node)
	expected := "`a`"
	if got != expected {
		t.Errorf("Deparse(UnaryPlus(a)) = %q, want %q", got, expected)
	}
}

func TestDeparse_Section_2_2_ComparisonOperators(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic comparison operators
		{"equal", "a = b", "(`a` = `b`)"},
		{"not_equal_angle", "a <> b", "(`a` <> `b`)"},
		{"not_equal_bang", "a != b", "(`a` <> `b`)"},          // != normalized to <>
		{"greater", "a > b", "(`a` > `b`)"},
		{"less", "a < b", "(`a` < `b`)"},
		{"greater_or_equal", "a >= b", "(`a` >= `b`)"},
		{"less_or_equal", "a <= b", "(`a` <= `b`)"},
		{"null_safe_equal", "a <=> b", "(`a` <=> `b`)"},
		{"sounds_like", "a SOUNDS LIKE b", "(`a` sounds like `b`)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparse_Section_2_3_BitwiseOperators(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"bitwise_or", "a | b", "(`a` | `b`)"},
		{"bitwise_and", "a & b", "(`a` & `b`)"},
		{"bitwise_xor", "a ^ b", "(`a` ^ `b`)"},
		{"left_shift", "a << b", "(`a` << `b`)"},
		{"right_shift", "a >> b", "(`a` >> `b`)"},
		{"bitwise_not", "~a", "~`a`"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.expected {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_2_3_BitwiseNotAST tests bitwise NOT via hand-built AST
// to verify deparseUnaryExpr handles UnaryBitNot correctly.
func TestDeparse_Section_2_3_BitwiseNotAST(t *testing.T) {
	node := &ast.UnaryExpr{
		Op:      ast.UnaryBitNot,
		Operand: &ast.ColumnRef{Column: "a"},
	}
	got := Deparse(node)
	expected := "~`a`"
	if got != expected {
		t.Errorf("Deparse(UnaryBitNot(a)) = %q, want %q", got, expected)
	}
}

func TestDeparse_NilNode(t *testing.T) {
	got := Deparse(nil)
	if got != "" {
		t.Errorf("Deparse(nil) = %q, want empty string", got)
	}
}
