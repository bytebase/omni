package deparse

import (
	"strings"
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
		{"unary_minus", "-a", "-(`a`)"},

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
		{"bitwise_not", "~a", "~(`a`)"},
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
	expected := "~(`a`)"
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
			"sum(`b`) OVER (PARTITION BY `a` ORDER BY `b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"},
		// Named window reference — MySQL 8.0 backtick-quotes window names
		{"named_window_ref",
			"ROW_NUMBER() OVER w",
			"row_number() OVER `w`"},
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

// parseRewriteDeparse is a helper that parses an expression, applies RewriteExpr, and deparses.
func parseRewriteDeparse(t *testing.T, expr string) string {
	t.Helper()
	node := parseExpr(t, expr)
	rewritten := RewriteExpr(node)
	return Deparse(rewritten)
}

func TestDeparse_Section_4_1_NOTFolding(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// NOT (comparison) → inverted comparison operator
		{"not_gt", "NOT (a > 0)", "(`a` <= 0)"},
		{"not_lt", "NOT (a < 0)", "(`a` >= 0)"},
		{"not_ge", "NOT (a >= 0)", "(`a` < 0)"},
		{"not_le", "NOT (a <= 0)", "(`a` > 0)"},
		{"not_eq", "NOT (a = 0)", "(`a` <> 0)"},
		{"not_ne", "NOT (a <> 0)", "(`a` = 0)"},

		// NOT (non-boolean) → (0 = expr)
		{"not_col", "NOT a", "(0 = `a`)"},
		{"not_add", "NOT (a + 1)", "(0 = (`a` + 1))"},

		// NOT LIKE → not((expr like pattern))  — stays as not() wrapping
		{"not_like", "NOT (a LIKE 'foo%')", "(not((`a` like 'foo%')))"},

		// ! operator on column — same as NOT
		{"bang_col", "!a", "(0 = `a`)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRewriteDeparse(t, tc.input)
			if got != tc.expected {
				t.Errorf("RewriteExpr+Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_4_1_NOTFolding_AST tests NOT folding via hand-built AST nodes
// for cases where the parser may handle NOT differently (e.g., NOT LIKE is parsed
// directly as LikeExpr with Not=true, not as UnaryNot wrapping LikeExpr).
func TestDeparse_Section_4_1_NOTFolding_AST(t *testing.T) {
	// NOT wrapping a LikeExpr (as if we manually constructed this AST)
	t.Run("not_wrapping_like_ast", func(t *testing.T) {
		node := &ast.UnaryExpr{
			Op: ast.UnaryNot,
			Operand: &ast.LikeExpr{
				Expr:    &ast.ColumnRef{Column: "a"},
				Pattern: &ast.StringLit{Value: "foo%"},
			},
		}
		rewritten := RewriteExpr(node)
		got := Deparse(rewritten)
		expected := "(not((`a` like 'foo%')))"
		if got != expected {
			t.Errorf("RewriteExpr+Deparse(NOT(a LIKE 'foo%%')) = %q, want %q", got, expected)
		}
	})

	// NOT wrapping comparison via AST (no ParenExpr wrapper)
	t.Run("not_gt_no_paren_ast", func(t *testing.T) {
		node := &ast.UnaryExpr{
			Op: ast.UnaryNot,
			Operand: &ast.BinaryExpr{
				Op:    ast.BinOpGt,
				Left:  &ast.ColumnRef{Column: "a"},
				Right: &ast.IntLit{Value: 0},
			},
		}
		rewritten := RewriteExpr(node)
		got := Deparse(rewritten)
		expected := "(`a` <= 0)"
		if got != expected {
			t.Errorf("RewriteExpr+Deparse(NOT(a > 0)) = %q, want %q", got, expected)
		}
	})

	// ! on column ref — (0 = `a`)
	t.Run("bang_column_ast", func(t *testing.T) {
		node := &ast.UnaryExpr{
			Op:      ast.UnaryNot,
			Operand: &ast.ColumnRef{Column: "a"},
		}
		rewritten := RewriteExpr(node)
		got := Deparse(rewritten)
		expected := "(0 = `a`)"
		if got != expected {
			t.Errorf("RewriteExpr+Deparse(!a) = %q, want %q", got, expected)
		}
	})

	// NOT on arithmetic — (0 = (a + 1))
	t.Run("not_arithmetic_ast", func(t *testing.T) {
		node := &ast.UnaryExpr{
			Op: ast.UnaryNot,
			Operand: &ast.BinaryExpr{
				Op:    ast.BinOpAdd,
				Left:  &ast.ColumnRef{Column: "a"},
				Right: &ast.IntLit{Value: 1},
			},
		}
		rewritten := RewriteExpr(node)
		got := Deparse(rewritten)
		expected := "(0 = (`a` + 1))"
		if got != expected {
			t.Errorf("RewriteExpr+Deparse(NOT(a+1)) = %q, want %q", got, expected)
		}
	})
}

// TestDeparse_Section_4_2_BooleanContextWrapping tests isBooleanExpr identification
// and boolean context wrapping for AND/OR/XOR operands.
func TestDeparse_Section_4_2_BooleanContextWrapping(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Column ref in AND: non-boolean operands get wrapped
		{"col_ref_in_and", "a AND b", "((0 <> `a`) and (0 <> `b`))"},
		// Arithmetic in AND
		{"arithmetic_in_and", "(a+1) AND b", "((0 <> (`a` + 1)) and (0 <> `b`))"},
		// Function in AND
		{"function_in_and", "ABS(a) AND b", "((0 <> abs(`a`)) and (0 <> `b`))"},
		// CASE in AND
		{"case_in_and", "CASE WHEN a > 0 THEN 1 ELSE 0 END AND b",
			"((0 <> (case when (`a` > 0) then 1 else 0 end)) and (0 <> `b`))"},
		// IF in AND
		{"if_in_and", "IF(a > 0, 1, 0) AND b",
			"((0 <> if((`a` > 0),1,0)) and (0 <> `b`))"},
		// Literal in AND
		{"literal_in_and", "'hello' AND 1", "((0 <> 'hello') and (0 <> 1))"},
		// Comparison NOT wrapped: both sides are boolean
		{"comparison_not_wrapped", "(a > 0) AND (b > 0)", "((`a` > 0) and (`b` > 0))"},
		// Mixed: one boolean, one non-boolean
		{"mixed_bool_nonbool", "(a > 0) AND (b + 1)", "((`a` > 0) and (0 <> (`b` + 1)))"},
		// XOR: non-boolean operands get wrapped
		{"xor_wrapping", "a XOR b", "((0 <> `a`) xor (0 <> `b`))"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRewriteDeparse(t, tc.input)
			if got != tc.expected {
				t.Errorf("RewriteExpr+Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_4_2_IsBooleanExpr tests isBooleanExpr identification via AST.
func TestDeparse_Section_4_2_IsBooleanExpr(t *testing.T) {
	// Test that comparison, IN, BETWEEN, LIKE, IS NULL, AND/OR/NOT/XOR, EXISTS, TRUE/FALSE
	// are all recognized as boolean expressions (not wrapped in AND context).

	booleanCases := []struct {
		name  string
		input string
	}{
		// Comparisons are boolean
		{"comparison_eq", "(a = 1) AND (b = 2)"},
		{"comparison_ne", "(a <> 1) AND (b <> 2)"},
		{"comparison_lt", "(a < 1) AND (b < 2)"},
		{"comparison_gt", "(a > 1) AND (b > 2)"},
		{"comparison_le", "(a <= 1) AND (b <= 2)"},
		{"comparison_ge", "(a >= 1) AND (b >= 2)"},
		{"comparison_nullsafe", "(a <=> 1) AND (b <=> 2)"},
		// IN is boolean
		{"in_is_boolean", "(a IN (1,2)) AND (b IN (3,4))"},
		// BETWEEN is boolean
		{"between_is_boolean", "(a BETWEEN 1 AND 10) AND (b BETWEEN 1 AND 10)"},
		// LIKE is boolean
		{"like_is_boolean", "(a LIKE 'x') AND (b LIKE 'y')"},
		// IS NULL is boolean
		{"is_null_is_boolean", "(a IS NULL) AND (b IS NULL)"},
		// TRUE/FALSE literals are boolean
		{"bool_lit_true", "TRUE AND FALSE"},
	}

	for _, tc := range booleanCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRewriteDeparse(t, tc.input)
			// None of these should contain "(0 <>" wrapping
			if contains0Ne(got) {
				t.Errorf("Boolean expression was incorrectly wrapped: RewriteExpr+Deparse(%q) = %q", tc.input, got)
			}
		})
	}
}

// TestDeparse_Section_4_2_ISTrueFalse tests IS TRUE/IS FALSE wrapping on non-boolean.
func TestDeparse_Section_4_2_ISTrueFalse(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// IS TRUE on non-boolean column: wrap with (0 <> col)
		{"is_true_nonbool", "a IS TRUE", "((0 <> `a`) is true)"},
		// IS FALSE on non-boolean column: wrap with (0 <> col)
		{"is_false_nonbool", "a IS FALSE", "((0 <> `a`) is false)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRewriteDeparse(t, tc.input)
			if got != tc.expected {
				t.Errorf("RewriteExpr+Deparse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparse_Section_4_2_SubqueryInAnd tests subquery wrapping in AND context.
func TestDeparse_Section_4_2_SubqueryInAnd(t *testing.T) {
	// Subquery in AND: non-boolean subquery gets wrapped
	// Note: we build AST directly since parser may not support subqueries in this context easily
	t.Run("subquery_in_and", func(t *testing.T) {
		// Build: (SELECT 1) AND (SELECT 2) — both are SubqueryExpr (non-boolean)
		node := &ast.BinaryExpr{
			Op: ast.BinOpAnd,
			Left: &ast.SubqueryExpr{
				Select: &ast.SelectStmt{
					TargetList: []ast.ExprNode{
						&ast.ResTarget{Val: &ast.IntLit{Value: 1}},
					},
				},
			},
			Right: &ast.SubqueryExpr{
				Select: &ast.SelectStmt{
					TargetList: []ast.ExprNode{
						&ast.ResTarget{Val: &ast.IntLit{Value: 2}},
					},
				},
			},
		}
		rewritten := RewriteExpr(node)
		got := Deparse(rewritten)
		// SubqueryExpr is not boolean, so should be wrapped
		// The exact output depends on SubqueryExpr deparsing — we just verify wrapping happened
		if !contains0Ne(got) {
			t.Errorf("Subquery in AND was NOT wrapped: got %q", got)
		}
	})
}

// contains0Ne checks if the output contains "(0 <>" which indicates boolean wrapping.
func contains0Ne(s string) bool {
	return strings.Contains(s, "(0 <>") || strings.Contains(s, "(0 =")
}

// TestRewriteExpr_NilNode tests that RewriteExpr handles nil gracefully.
func TestRewriteExpr_NilNode(t *testing.T) {
	got := RewriteExpr(nil)
	if got != nil {
		t.Errorf("RewriteExpr(nil) = %v, want nil", got)
	}
}

// parseSelect parses a full SQL SELECT statement and returns the SelectStmt.
func parseSelect(t *testing.T, sql string) *ast.SelectStmt {
	t.Helper()
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
	return sel
}

func TestDeparseSelect_Section_5_1_TargetListAliases(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Single column with FROM
		{"single_column", "SELECT a FROM t", "select `a` AS `a` from `t`"},
		// Multiple columns: comma-separated, no space after comma
		{"multiple_columns", "SELECT a, b, c FROM t", "select `a` AS `a`,`b` AS `b`,`c` AS `c` from `t`"},
		// Column alias with AS
		{"column_alias_as", "SELECT a AS col1 FROM t", "select `a` AS `col1` from `t`"},
		// Column alias without AS (should still produce AS in output)
		{"column_alias_no_as", "SELECT a col1 FROM t", "select `a` AS `col1` from `t`"},
		// Expression alias
		{"expression_alias", "SELECT a + b AS sum_col FROM t", "select (`a` + `b`) AS `sum_col` from `t`"},
		// Auto-alias literal: SELECT 1 → 1 AS `1`
		{"auto_alias_literal", "SELECT 1", "select 1 AS `1`"},
		// Auto-alias expression: SELECT a + b → (`a` + `b`) AS `a + b`
		{"auto_alias_expression", "SELECT a + b FROM t", "select (`a` + `b`) AS `a + b` from `t`"},
		// Auto-alias string literal
		{"auto_alias_string", "SELECT 'hello'", "select 'hello' AS `hello`"},
		// Auto-alias NULL
		{"auto_alias_null", "SELECT NULL", "select NULL AS `NULL`"},
		// Auto-alias function call
		{"auto_alias_func", "SELECT CONCAT(a, b) FROM t", "select concat(`a`,`b`) AS `CONCAT(a, b)` from `t`"},
		// Auto-alias boolean literal
		{"auto_alias_true", "SELECT TRUE", "select true AS `TRUE`"},
		{"auto_alias_false", "SELECT FALSE", "select false AS `FALSE`"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_1_NilStmt(t *testing.T) {
	got := DeparseSelect(nil)
	if got != "" {
		t.Errorf("DeparseSelect(nil) = %q, want empty string", got)
	}
}

func TestDeparseSelect_Section_5_2_FromClause(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Single table
		{"single_table", "SELECT a FROM t", "select `a` AS `a` from `t`"},
		// Table alias with AS — no AS keyword in output for table alias
		{"table_alias_with_as", "SELECT a FROM t AS t1", "select `a` AS `a` from `t` `t1`"},
		// Table alias without AS — same output
		{"table_alias_without_as", "SELECT a FROM t t1", "select `a` AS `a` from `t` `t1`"},
		// Multiple tables (implicit cross join) → explicit join with parens
		{"implicit_cross_join", "SELECT a FROM t1, t2", "select `a` AS `a` from (`t1` join `t2`)"},
		// Three tables implicit cross join
		{"implicit_cross_join_3", "SELECT a FROM t1, t2, t3", "select `a` AS `a` from ((`t1` join `t2`) join `t3`)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_2_DerivedTable(t *testing.T) {
	// Note: Parser currently does not produce SubqueryExpr for derived tables in FROM clause.
	// These tests verify the deparse logic works when given the correct AST manually.
	t.Run("derived_table_subquery_expr", func(t *testing.T) {
		// Manually construct: FROM (SELECT 1 AS a) `d`
		inner := &ast.SelectStmt{
			TargetList: []ast.ExprNode{
				&ast.ResTarget{Name: "a", Val: &ast.IntLit{Value: 1}},
			},
		}
		subq := &ast.SubqueryExpr{
			Select: inner,
			Alias:  "d",
		}
		outer := &ast.SelectStmt{
			TargetList: []ast.ExprNode{
				&ast.ResTarget{Val: &ast.ColumnRef{Column: "a"}},
			},
			From: []ast.TableExpr{subq},
		}
		got := DeparseSelect(outer)
		expected := "select `a` AS `a` from (select 1 AS `a`) `d`"
		if got != expected {
			t.Errorf("DeparseSelect(derived table) =\n  %q\nwant:\n  %q", got, expected)
		}
	})
}

func TestDeparseSelect_Section_5_3_JoinClause(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// INNER JOIN
		{"inner_join", "SELECT a FROM t1 JOIN t2 ON t1.a = t2.a",
			"select `a` AS `a` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))"},
		// LEFT JOIN
		{"left_join", "SELECT a FROM t1 LEFT JOIN t2 ON t1.a = t2.a",
			"select `a` AS `a` from (`t1` left join `t2` on((`t1`.`a` = `t2`.`a`)))"},
		// RIGHT JOIN → LEFT JOIN with table swap
		{"right_join", "SELECT a FROM t1 RIGHT JOIN t2 ON t1.a = t2.a",
			"select `a` AS `a` from (`t2` left join `t1` on((`t1`.`a` = `t2`.`a`)))"},
		// CROSS JOIN → plain join (no ON)
		{"cross_join", "SELECT a FROM t1 CROSS JOIN t2",
			"select `a` AS `a` from (`t1` join `t2`)"},
		// STRAIGHT_JOIN — lowercase
		{"straight_join", "SELECT a FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a",
			"select `a` AS `a` from (`t1` straight_join `t2` on((`t1`.`a` = `t2`.`a`)))"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_3_NaturalJoin(t *testing.T) {
	// NATURAL JOIN — expanded to join without ON (needs Phase 6 resolver for column expansion)
	// For now, verify basic format; full ON expansion requires schema info.
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// NATURAL JOIN → join (no ON without resolver)
		{"natural_join", "SELECT a FROM t1 NATURAL JOIN t2",
			"select `a` AS `a` from (`t1` join `t2`)"},
		// NATURAL LEFT JOIN → left join (no ON without resolver)
		{"natural_left_join", "SELECT a FROM t1 NATURAL LEFT JOIN t2",
			"select `a` AS `a` from (`t1` left join `t2`)"},
		// NATURAL RIGHT JOIN → left join with table swap (no ON without resolver)
		{"natural_right_join", "SELECT a FROM t1 NATURAL RIGHT JOIN t2",
			"select `a` AS `a` from (`t2` left join `t1`)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_3_UsingClause(t *testing.T) {
	// USING — expanded to ON with qualified column references
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// USING single column
		{"using_single", "SELECT a FROM t1 JOIN t2 USING (a)",
			"select `a` AS `a` from (`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestDeparseSelect_Section_5_3_UsingMultiColumn tests USING with multiple columns via AST.
func TestDeparseSelect_Section_5_3_UsingMultiColumn(t *testing.T) {
	// Build AST manually for USING (a, b) since parser may only support single-column USING parsing
	t.Run("using_multi_column", func(t *testing.T) {
		join := &ast.JoinClause{
			Type:  ast.JoinInner,
			Left:  &ast.TableRef{Name: "t1"},
			Right: &ast.TableRef{Name: "t2"},
			Condition: &ast.UsingCondition{
				Columns: []string{"a", "b"},
			},
		}
		sel := &ast.SelectStmt{
			TargetList: []ast.ExprNode{
				&ast.ResTarget{Val: &ast.ColumnRef{Column: "x"}},
			},
			From: []ast.TableExpr{join},
		}
		got := DeparseSelect(sel)
		expected := "select `x` AS `x` from (`t1` join `t2` on(((`t1`.`a` = `t2`.`a`) and (`t1`.`b` = `t2`.`b`))))"
		if got != expected {
			t.Errorf("DeparseSelect(USING multi) =\n  %q\nwant:\n  %q", got, expected)
		}
	})
}

func TestDeparseSelect_Section_5_4_WhereGroupByHavingOrderByLimit(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// WHERE simple
		{"where_simple", "SELECT a FROM t WHERE a > 1",
			"select `a` AS `a` from `t` where (`a` > 1)"},
		// WHERE compound
		{"where_compound", "SELECT a FROM t WHERE a > 1 AND b < 10",
			"select `a` AS `a` from `t` where ((`a` > 1) and (`b` < 10))"},
		// GROUP BY single
		{"group_by_single", "SELECT a FROM t GROUP BY a",
			"select `a` AS `a` from `t` group by `a`"},
		// GROUP BY multiple — no space after comma
		{"group_by_multiple", "SELECT a, b FROM t GROUP BY a, b",
			"select `a` AS `a`,`b` AS `b` from `t` group by `a`,`b`"},
		// GROUP BY WITH ROLLUP
		{"group_by_with_rollup", "SELECT a FROM t GROUP BY a WITH ROLLUP",
			"select `a` AS `a` from `t` group by `a` with rollup"},
		// HAVING
		{"having", "SELECT a FROM t GROUP BY a HAVING COUNT(*) > 1",
			"select `a` AS `a` from `t` group by `a` having (count(0) > 1)"},
		// ORDER BY ASC (default — no explicit ASC in output)
		{"order_by_asc", "SELECT a FROM t ORDER BY a",
			"select `a` AS `a` from `t` order by `a`"},
		// ORDER BY DESC
		{"order_by_desc", "SELECT a FROM t ORDER BY a DESC",
			"select `a` AS `a` from `t` order by `a` desc"},
		// ORDER BY multiple — comma no space
		{"order_by_multiple", "SELECT a, b FROM t ORDER BY a, b DESC",
			"select `a` AS `a`,`b` AS `b` from `t` order by `a`,`b` desc"},
		// LIMIT
		{"limit", "SELECT a FROM t LIMIT 10",
			"select `a` AS `a` from `t` limit 10"},
		// LIMIT with OFFSET — MySQL comma syntax: LIMIT offset,count
		{"limit_with_offset", "SELECT a FROM t LIMIT 10 OFFSET 5",
			"select `a` AS `a` from `t` limit 5,10"},
		// DISTINCT
		{"distinct", "SELECT DISTINCT a FROM t",
			"select distinct `a` AS `a` from `t`"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_6_Subqueries(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Scalar subquery in SELECT list
		// MySQL 8.0 omits column aliases inside scalar subquery body.
		// Auto-alias uses uppercase keywords and unqualified column names.
		{"scalar_subquery", "SELECT (SELECT MAX(a) FROM t) FROM t",
			"select (select max(`a`) from `t`) AS `(SELECT MAX(a) FROM t)` from `t`"},
		// IN subquery in WHERE — MySQL 8.0 omits outer parens and aliases in subquery target list
		{"in_subquery", "SELECT a FROM t WHERE a IN (SELECT a FROM t)",
			"select `a` AS `a` from `t` where `a` in (select `a` from `t`)"},
		// EXISTS in WHERE
		{"exists_subquery", "SELECT a FROM t WHERE EXISTS (SELECT 1 FROM t)",
			"select `a` AS `a` from `t` where exists(select 1 from `t`)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_5_SetOperations(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// UNION
		{"union", "SELECT a FROM t UNION SELECT b FROM t",
			"select `a` AS `a` from `t` union select `b` AS `b` from `t`"},
		// UNION ALL
		{"union_all", "SELECT a FROM t UNION ALL SELECT b FROM t",
			"select `a` AS `a` from `t` union all select `b` AS `b` from `t`"},
		// Multiple UNION: three SELECTs chained flat (left-associative)
		{"multiple_union", "SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t",
			"select `a` AS `a` from `t` union select `b` AS `b` from `t` union select `c` AS `c` from `t`"},
		// INTERSECT
		{"intersect", "SELECT a FROM t INTERSECT SELECT b FROM t",
			"select `a` AS `a` from `t` intersect select `b` AS `b` from `t`"},
		// EXCEPT
		{"except", "SELECT a FROM t EXCEPT SELECT b FROM t",
			"select `a` AS `a` from `t` except select `b` AS `b` from `t`"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparseSelect_Section_5_7_CTE(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Simple CTE
		{
			"simple_cte",
			"WITH cte AS (SELECT 1) SELECT a FROM cte",
			"with `cte` as (select 1 AS `1`) select `a` AS `a` from `cte`",
		},
		// CTE with column list
		{
			"cte_with_columns",
			"WITH cte(x) AS (SELECT 1) SELECT x FROM cte",
			"with `cte` (`x`) as (select 1 AS `1`) select `x` AS `x` from `cte`",
		},
		// RECURSIVE CTE
		{
			"recursive_cte",
			"WITH RECURSIVE cte AS (SELECT 1 UNION ALL SELECT 1) SELECT a FROM cte",
			"with recursive `cte` as (select 1 AS `1` union all select 1 AS `1`) select `a` AS `a` from `cte`",
		},
		// Multiple CTEs
		{
			"multiple_ctes",
			"WITH c1 AS (SELECT 1), c2 AS (SELECT 2) SELECT a FROM c1, c2",
			"with `c1` as (select 1 AS `1`), `c2` as (select 2 AS `2`) select `a` AS `a` from (`c1` join `c2`)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDeparse_Section_5_8_ForUpdate(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"for_update",
			"SELECT a FROM t FOR UPDATE",
			"select `a` AS `a` from `t` for update",
		},
		{
			"for_share",
			"SELECT a FROM t FOR SHARE",
			"select `a` AS `a` from `t` for share",
		},
		{
			"lock_in_share_mode",
			"SELECT a FROM t LOCK IN SHARE MODE",
			"select `a` AS `a` from `t` lock in share mode",
		},
		{
			"for_update_of_table",
			"SELECT a FROM t FOR UPDATE OF t",
			"select `a` AS `a` from `t` for update of `t`",
		},
		{
			"for_update_nowait",
			"SELECT a FROM t FOR UPDATE NOWAIT",
			"select `a` AS `a` from `t` for update nowait",
		},
		{
			"for_update_skip_locked",
			"SELECT a FROM t FOR UPDATE SKIP LOCKED",
			"select `a` AS `a` from `t` for update skip locked",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := parseSelect(t, tc.input)
			got := DeparseSelect(sel)
			if got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}
