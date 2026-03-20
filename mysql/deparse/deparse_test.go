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

func TestDeparse_Section_2_4_PrecedenceParenthesization(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Higher precedence preserved: * binds tighter than +
		{"higher_prec_preserved", "a + b * c", "(`a` + (`b` * `c`))"},
		// Lower precedence grouping: explicit parens force + before *
		{"lower_prec_grouping", "(a + b) * c", "((`a` + `b`) * `c`)"},
		// Mixed precedence: * first, then +
		{"mixed_precedence", "a * b + c", "((`a` * `b`) + `c`)"},
		// Deeply nested left-associative chaining
		{"deeply_nested", "a + b + c + a + b + c", "(((((`a` + `b`) + `c`) + `a`) + `b`) + `c`)"},
		// Parenthesized expression passthrough — ParenExpr unwrapped, BinaryExpr provides outer parens
		{"paren_passthrough", "(a + b)", "(`a` + `b`)"},
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

func TestDeparse_Section_2_5_ComparisonPredicates(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// IN list
		{"in_list", "a IN (1,2,3)", "(`a` in (1,2,3))"},
		// NOT IN
		{"not_in_list", "a NOT IN (1,2,3)", "(`a` not in (1,2,3))"},
		// BETWEEN
		{"between", "a BETWEEN 1 AND 10", "(`a` between 1 and 10)"},
		// NOT BETWEEN
		{"not_between", "a NOT BETWEEN 1 AND 10", "(`a` not between 1 and 10)"},
		// LIKE
		{"like", "a LIKE 'foo%'", "(`a` like 'foo%')"},
		// LIKE with ESCAPE
		{"like_escape", "a LIKE 'x' ESCAPE '\\\\'", "(`a` like 'x' escape '\\\\')"},
		// IS NULL
		{"is_null", "a IS NULL", "(`a` is null)"},
		// IS NOT NULL
		{"is_not_null", "a IS NOT NULL", "(`a` is not null)"},
		// IS TRUE
		{"is_true", "a IS TRUE", "(`a` is true)"},
		// IS FALSE
		{"is_false", "a IS FALSE", "(`a` is false)"},
		// IS UNKNOWN
		{"is_unknown", "a IS UNKNOWN", "(`a` is unknown)"},
		// ROW comparison
		{"row_comparison", "ROW(a,b) = ROW(1,2)", "(row(`a`,`b`) = row(1,2))"},
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

func TestDeparse_Section_2_6_CaseCastConvert(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Searched CASE
		{"searched_case", "CASE WHEN a > 0 THEN 'pos' ELSE 'zero' END",
			"(case when (`a` > 0) then 'pos' else 'zero' end)"},
		// Searched CASE multiple WHEN
		{"searched_case_multi", "CASE WHEN a>0 THEN 'a' WHEN b>0 THEN 'b' ELSE 'c' END",
			"(case when (`a` > 0) then 'a' when (`b` > 0) then 'b' else 'c' end)"},
		// Simple CASE
		{"simple_case", "CASE a WHEN 1 THEN 'one' ELSE 'other' END",
			"(case `a` when 1 then 'one' else 'other' end)"},
		// CASE without ELSE
		{"case_no_else", "CASE WHEN a > 0 THEN 'pos' END",
			"(case when (`a` > 0) then 'pos' end)"},
		// CAST to CHAR
		{"cast_char", "CAST(a AS CHAR)", "cast(`a` as char charset utf8mb4)"},
		// CAST to CHAR(N)
		{"cast_char_n", "CAST(a AS CHAR(10))", "cast(`a` as char(10) charset utf8mb4)"},
		// CAST to BINARY
		{"cast_binary", "CAST(a AS BINARY)", "cast(`a` as char charset binary)"},
		// CAST to SIGNED
		{"cast_signed", "CAST(a AS SIGNED)", "cast(`a` as signed)"},
		// CAST to UNSIGNED
		{"cast_unsigned", "CAST(a AS UNSIGNED)", "cast(`a` as unsigned)"},
		// CAST to DECIMAL
		{"cast_decimal", "CAST(a AS DECIMAL(10,2))", "cast(`a` as decimal(10,2))"},
		// CAST to DATE
		{"cast_date", "CAST(a AS DATE)", "cast(`a` as date)"},
		// CAST to DATETIME
		{"cast_datetime", "CAST(a AS DATETIME)", "cast(`a` as datetime)"},
		// CAST to JSON
		{"cast_json", "CAST(a AS JSON)", "cast(`a` as json)"},
		// CONVERT USING
		{"convert_using", "CONVERT(a USING utf8mb4)", "convert(`a` using utf8mb4)"},
		// CONVERT type — rewritten to cast
		{"convert_type_char", "CONVERT(a, CHAR)", "cast(`a` as char charset utf8mb4)"},
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

func TestDeparse_Section_2_7_OtherExpressions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// INTERVAL: INTERVAL 1 DAY + a → (`a` + interval 1 day) — operand order swapped
		{"interval_add", "INTERVAL 1 DAY + a", "(`a` + interval 1 day)"},
		// COLLATE: a COLLATE utf8mb4_bin → (`a` collate utf8mb4_bin)
		{"collate", "a COLLATE utf8mb4_bin", "(`a` collate utf8mb4_bin)"},
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

func TestDeparse_Section_3_1_FunctionsAndRewrites(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Regular function calls — lowercase, no space after comma
		{"simple_concat", "CONCAT(a, b)", "concat(`a`,`b`)"},
		{"nested_functions", "CONCAT(UPPER(TRIM(a)), LOWER(b))", "concat(upper(trim(`a`)),lower(`b`))"},
		{"ifnull", "IFNULL(a, 0)", "ifnull(`a`,0)"},
		{"coalesce", "COALESCE(a, b, 0)", "coalesce(`a`,`b`,0)"},
		{"nullif", "NULLIF(a, 0)", "nullif(`a`,0)"},
		{"if_func", "IF(a > 0, 'yes', 'no')", "if((`a` > 0),'yes','no')"},
		{"abs", "ABS(a)", "abs(`a`)"},
		{"greatest", "GREATEST(a, b)", "greatest(`a`,`b`)"},
		{"least", "LEAST(a, b)", "least(`a`,`b`)"},

		// Function name rewrites
		{"substring_to_substr", "SUBSTRING(a, 1, 3)", "substr(`a`,1,3)"},
		{"current_timestamp_no_parens", "CURRENT_TIMESTAMP", "now()"},
		{"current_timestamp_parens", "CURRENT_TIMESTAMP()", "now()"},
		{"current_date", "CURRENT_DATE", "curdate()"},
		{"current_time", "CURRENT_TIME", "curtime()"},
		{"current_user", "CURRENT_USER", "current_user()"},
		{"now_stays", "NOW()", "now()"},
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

func TestDeparse_Section_3_2_SpecialFunctionForms(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// TRIM simple — no direction, single arg
		{"trim_simple", "TRIM(a)", "trim(`a`)"},
		// TRIM LEADING
		{"trim_leading", "TRIM(LEADING 'x' FROM a)", "trim(leading 'x' from `a`)"},
		// TRIM TRAILING
		{"trim_trailing", "TRIM(TRAILING 'x' FROM a)", "trim(trailing 'x' from `a`)"},
		// TRIM BOTH
		{"trim_both", "TRIM(BOTH 'x' FROM a)", "trim(both 'x' from `a`)"},
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

func TestDeparse_Section_3_3_AggregateFunctions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// COUNT(*) — MySQL 8.0 rewrites * to 0
		{"count_star", "COUNT(*)", "count(0)"},
		// COUNT(col)
		{"count_col", "COUNT(a)", "count(`a`)"},
		// COUNT(DISTINCT col)
		{"count_distinct", "COUNT(DISTINCT a)", "count(distinct `a`)"},
		// SUM
		{"sum", "SUM(a)", "sum(`a`)"},
		// AVG
		{"avg", "AVG(a)", "avg(`a`)"},
		// MAX
		{"max", "MAX(a)", "max(`a`)"},
		// MIN
		{"min", "MIN(a)", "min(`a`)"},
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

func TestDeparse_Section_3_4_GroupConcat(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic GROUP_CONCAT — default separator always shown
		{"basic", "GROUP_CONCAT(a)", "group_concat(`a` separator ',')"},
		// With ORDER BY — ASC shown explicitly
		{"with_order_by", "GROUP_CONCAT(a ORDER BY a)", "group_concat(`a` order by `a` ASC separator ',')"},
		// With explicit SEPARATOR
		{"with_separator", "GROUP_CONCAT(a SEPARATOR ';')", "group_concat(`a` separator ';')"},
		// DISTINCT + ORDER BY DESC + SEPARATOR — full combination
		{"distinct_order_desc_separator", "GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')",
			"group_concat(distinct `a` order by `a` DESC separator ';')"},
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

func TestDeparse_Section_3_5_WindowFunctions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// ROW_NUMBER with ORDER BY
		{"row_number_over_order_by",
			"ROW_NUMBER() OVER (ORDER BY a)",
			"row_number() OVER (ORDER BY `a` )"},
		// SUM with PARTITION BY and ORDER BY
		{"sum_over_partition_order",
			"SUM(b) OVER (PARTITION BY a ORDER BY b)",
			"sum(`b`) OVER (PARTITION BY `a` ORDER BY `b` )"},
		// Frame clause: ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		{"sum_over_frame",
			"SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
			"sum(`b`) OVER (PARTITION BY `a` ORDER BY `b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW )"},
		// Named window reference
		{"named_window_ref",
			"ROW_NUMBER() OVER w",
			"row_number() OVER w"},
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

func TestDeparse_Section_3_6_OperatorToFunctionRewrites(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// REGEXP → regexp_like()
		{"regexp", "a REGEXP 'pattern'", "regexp_like(`a`,'pattern')"},
		// NOT REGEXP → (not(regexp_like()))
		{"not_regexp", "a NOT REGEXP 'p'", "(not(regexp_like(`a`,'p')))"},
		// -> → json_extract()
		{"json_extract", "a->'$.key'", "json_extract(`a`,'$.key')"},
		// ->> → json_unquote(json_extract())
		{"json_unquote", "a->>'$.key'", "json_unquote(json_extract(`a`,'$.key'))"},
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

func TestDeparse_NilNode(t *testing.T) {
	got := Deparse(nil)
	if got != "" {
		t.Errorf("Deparse(nil) = %q, want empty string", got)
	}
}
