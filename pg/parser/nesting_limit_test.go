package parser

import (
	"errors"
	"strings"
	"testing"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	saved := maxNestingDepth
	maxNestingDepth = 1 << 20
	t.Cleanup(func() { maxNestingDepth = saved })
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: a few hundred kilobytes of '(' were enough. The
// same constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
		{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true" }},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
		{"function call", func(n int) string { return "SELECT " + nest(n, "f(", "1", ")") }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS int)") }},
		{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "true", " THEN 1 END") }},
		{"ARRAY", func(n int) string { return "SELECT " + nest(n, "ARRAY[", "1", "]") }},
		{"ROW", func(n int) string { return "SELECT " + nest(n, "ROW(", "1", ")") }},
		{"b_expr", func(n int) string { return "CREATE TABLE t (c int DEFAULT " + nest(n, "(", "1", ")") + ")" }},
		{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
		{"EXISTS", func(n int) string { return nest(n, "SELECT 1 WHERE EXISTS (", "SELECT 1", ")") }},
		{"IN subquery", func(n int) string { return nest(n, "SELECT 1 WHERE 1 IN (", "SELECT 1", ")") }},
		{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
		{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON true)") }},
		{"parenthesized select", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
		{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
		{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
		{"GROUP BY grouping sets", func(n int) string { return "SELECT 1 GROUP BY " + nest(n, "GROUPING SETS (", "a", ")") }},
		{"JSON_TABLE nested columns", func(n int) string {
			return "SELECT * FROM JSON_TABLE('{}', '$' COLUMNS (" + nest(n, "NESTED PATH '$' COLUMNS (", "a int", ")") + "))"
		}},
		{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
		{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
		{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
	} {
		if _, err := Parse(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 500000} {
			_, err := Parse(tc.sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Code != "54001" {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// TestLexerLoopsOverCommentsAndStrings: the lexer moved from a comment or a
// string piece to the next token by calling itself, one call per comment, so
// a long script of comments overflowed the stack.
func TestLexerLoopsOverCommentsAndStrings(t *testing.T) {
	for name, sql := range map[string]string{
		"block comments":      "SELECT 1 " + strings.Repeat("/* c */ ", 2000000),
		"line comments":       "SELECT 1 " + strings.Repeat("-- c\n", 2000000),
		"string continuation": "SELECT 'a'" + strings.Repeat("\n'a'", 2000000),
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
