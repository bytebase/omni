package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/partiql/ast"
)

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

// nestingShapes are the constructs the parser descends into recursively,
// n levels deep.
var nestingShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") + " FROM t" }},
	{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1 FROM t" }},
	{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true FROM t" }},
	{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") + " FROM t" }},
	{"function call", func(n int) string { return "SELECT " + nest(n, "f(", "1", ")") + " FROM t" }},
	{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS INT)") + " FROM t" }},
	{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "true", " THEN 1 END") + " FROM t" }},
	{"array literal", func(n int) string { return "SELECT " + nest(n, "[", "1", "]") + " FROM t" }},
	{"bag literal", func(n int) string { return "SELECT " + nest(n, "<<", "1", ">>") + " FROM t" }},
	{"tuple literal", func(n int) string { return "SELECT " + nest(n, "{'a': ", "1", "}") + " FROM t" }},
	{"path index", func(n int) string { return "SELECT " + nest(n, "a[", "0", "]") + " FROM t" }},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT a FROM t", ") FROM t") }},
	{"WHERE subquery", func(n int) string { return nest(n, "SELECT a FROM t WHERE (", "SELECT a FROM t", ")") }},
	{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT * FROM t", ")") }},
	{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t CROSS JOIN ", "t", ")") }},
	{"set operation operand", func(n int) string { return nest(n, "SELECT * FROM t UNION (", "SELECT * FROM t", ")") }},
	{"graph pattern group", func(n int) string { return "SELECT * FROM (g MATCH " + nest(n, "(", "(a)", ")") + ")" }},
	{"bracketed graph pattern group", func(n int) string { return "SELECT * FROM (g MATCH " + nest(n, "[", "(a)", "]") + ")" }},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
	{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
	{"unclosed arrays", func(n int) string { return "SELECT " + strings.Repeat("[", n) }},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 300KB of '(' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a speculative parse that retried exponentially at the limit
// fails the test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		if _, err := Parse(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			sql := tc.sql(n)
			done := make(chan error, 1)
			go func() {
				_, err := Parse(sql)
				done <- err
			}()
			select {
			case err := <-done:
				var pe *ParseError
				if !errors.As(err, &pe) || pe.Message != nestingErrorMessage {
					t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
				}
			case <-time.After(10 * time.Second):
				t.Fatalf("%s: %d levels took over 10s", tc.name, n)
			}
		}
	}
}

// treeHeight returns how many nodes deep the tree under root is.
func treeHeight(root ast.Node) int {
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
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1+", n) + "1 FROM t" }},
	{"product", func(n int) string { return "SELECT " + strings.Repeat("1*", n) + "1 FROM t" }},
	{"concatenation", func(n int) string { return "SELECT " + strings.Repeat("'a'||", n) + "'a' FROM t" }},
	{"OR", func(n int) string { return "SELECT * FROM t WHERE " + strings.Repeat("a OR ", n) + "a" }},
	{"AND", func(n int) string { return "SELECT * FROM t WHERE " + strings.Repeat("a AND ", n) + "a" }},
	{"comparisons", func(n int) string { return "SELECT * FROM t WHERE " + strings.Repeat("a = ", n) + "a" }},
	{"IS NULL", func(n int) string { return "SELECT * FROM t WHERE a" + strings.Repeat(" IS NULL", n) }},
	{"IN", func(n int) string { return "SELECT * FROM t WHERE a" + strings.Repeat(" IN [1]", n) }},
	{"LIKE", func(n int) string { return "SELECT * FROM t WHERE a" + strings.Repeat(" LIKE 'a'", n) }},
	{"BETWEEN", func(n int) string { return "SELECT * FROM t WHERE a" + strings.Repeat(" BETWEEN 1 AND 2", n) }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT * FROM t UNION ALL ", n) + "SELECT * FROM t" }},
	{"OUTER EXCEPT", func(n int) string { return strings.Repeat("SELECT * FROM t OUTER EXCEPT ", n) + "SELECT * FROM t" }},
	{"cross joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" CROSS JOIN t", n) }},
	{"comma joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(", t", n) }},
	{"qualified joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON true", n) }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it (ast.NodeToString did at a million terms). A
// chain within the limit still parses.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		_, err := Parse(tc.sql(maxTreeHeight + 1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingErrorMessage || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, err)
		}
	}
	// The other entry points check the tree too.
	sum := strings.Repeat("1+", maxTreeHeight+1) + "1"
	if _, err := NewParser(sum).ParseExpr(); !isNestingError(err) {
		t.Errorf("ParseExpr: err = %v, want the nesting error", err)
	}
	if _, err := NewParser("SELECT " + sum + " FROM t").ParseStatement(); !isNestingError(err) {
		t.Errorf("ParseStatement: err = %v, want the nesting error", err)
	}
	// The error points at the statement it rejects.
	first := "SELECT 1 FROM t; "
	_, err := Parse(first + chainShapes[0].sql(maxTreeHeight+1))
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Message != nestingErrorMessage || pe.Position != len(first) {
		t.Errorf("second statement: err = %v, want the nesting error at offset %d", err, len(first))
	}
}

// TestTreeHeightPerByte: treeTooHigh skips statements shorter than
// maxTreeHeight/16 bytes, on the ground that no construct adds more than 16
// levels of tree per byte of text. Check that across the recursive
// constructs and the chains, and that a chain's height grows with the chain,
// so the height counts every level a consumer recurses through.
func TestTreeHeightPerByte(t *testing.T) {
	check := func(name, sql string, minHeight int) {
		list, err := Parse(sql)
		if err != nil || len(list.Items) == 0 {
			return
		}
		h := treeHeight(list.Items[0])
		if h > 16*len(sql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(sql))
		}
		if h < minHeight {
			t.Errorf("%s: tree %d deep, want at least %d", name, h, minHeight)
		}
	}
	for _, tc := range nestingShapes {
		check(tc.name, tc.sql(1000), 0)
	}
	for _, tc := range chainShapes {
		check(tc.name, tc.sql(1000), 1000)
	}
}
