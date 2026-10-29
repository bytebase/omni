package parser

import (
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	saved := maxNestingDepth
	maxNestingDepth = 1 << 20
	t.Cleanup(func() { maxNestingDepth = saved })
}

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

// nestingShapes are the constructs the parser descends into recursively,
// n levels deep.
var nestingShapes = []struct {
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
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: a few hundred kilobytes of '(' were enough. The
// same constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		if _, err := Parse(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
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

// treeHeight returns how many nodes deep the tree under root is.
func treeHeight(root nodes.Node) int {
	lo, hi := 1, 1
	for treeHeightExceeds(root, hi) {
		lo, hi = hi+1, hi*2
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if treeHeightExceeds(root, mid) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// chainShapes build left-deep chains of n terms: the parser builds each with
// a loop, so only the tree, not the parse, gets deep.
var chainShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1+", n) + "1" }},
	{"concatenation", func(n int) string { return "SELECT " + strings.Repeat("'a'||", n) + "'a'" }},
	{"OR", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a OR ", n) + "a" }},
	{"AND", func(n int) string { return "SELECT 1 WHERE " + strings.Repeat("a AND ", n) + "a" }},
	{"casts", func(n int) string { return "SELECT 1" + strings.Repeat("::int", n) }},
	{"subscripts", func(n int) string { return "SELECT a" + strings.Repeat("[1]", n) }},
	{"IS NULL", func(n int) string { return "SELECT a" + strings.Repeat(" IS NULL", n) }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT 1 UNION ALL ", n) + "SELECT 1" }},
	{"joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON true", n) }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it (the catalog did at three million terms). A
// chain within the limit still parses.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		_, err := Parse(tc.sql(maxTreeHeight + 1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Code != "54001" || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, err)
		}
	}
}

// TestTreeHeightPerByte: treeTooHigh skips statements shorter than
// maxTreeHeight/16 bytes, on the ground that no construct adds more than 16
// levels of tree per byte of text. Check that across the recursive
// constructs and the chains.
func TestTreeHeightPerByte(t *testing.T) {
	check := func(name, sql string) {
		list, err := Parse(sql)
		if err != nil || len(list.Items) == 0 {
			return
		}
		if h := treeHeight(list.Items[0]); h > 16*len(sql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(sql))
		}
	}
	for _, tc := range nestingShapes {
		check(tc.name, tc.sql(1000))
	}
	for _, tc := range chainShapes {
		check(tc.name, tc.sql(1000))
	}
}
