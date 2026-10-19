package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
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
	{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") + " FROM c" }},
	{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1 FROM c" }},
	{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true FROM c" }},
	{"bitwise NOT", func(n int) string { return "SELECT " + strings.Repeat("~", n) + "1 FROM c" }},
	{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") + " FROM c" }},
	{"ternary", func(n int) string { return "SELECT " + strings.Repeat("true ? 1 : ", n) + "1 FROM c" }},
	{"coalesce", func(n int) string { return "SELECT " + strings.Repeat("c.a ?? ", n) + "1 FROM c" }},
	{"function call", func(n int) string { return "SELECT " + nest(n, "ABS(", "1", ")") + " FROM c" }},
	{"UDF call", func(n int) string { return "SELECT " + nest(n, "udf.f(", "1", ")") + " FROM c" }},
	{"array literal", func(n int) string { return "SELECT " + nest(n, "[", "1", "]") + " FROM c" }},
	{"object literal", func(n int) string { return "SELECT " + nest(n, "{\"a\": ", "1", "}") + " FROM c" }},
	{"bracket index", func(n int) string { return "SELECT c" + nest(n, "[", "0", "]") + " FROM c" }},
	{"IN list", func(n int) string { return "SELECT * FROM c WHERE " + nest(n, "1 IN (", "1", ")") }},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT VALUE 1", ") FROM c") }},
	{"EXISTS", func(n int) string { return nest(n, "SELECT EXISTS(", "SELECT VALUE 1", ") FROM c") }},
	{"ARRAY subquery", func(n int) string { return nest(n, "SELECT ARRAY(", "SELECT VALUE 1", ") FROM c") }},
	{"FROM subquery", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT * FROM c", ")") }},
	{"JOIN subquery", func(n int) string { return nest(n, "SELECT * FROM c JOIN (", "SELECT * FROM c", ")") }},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT * FROM (", n) }},
	{"unclosed arrays", func(n int) string { return "SELECT " + strings.Repeat("[", n) }},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 880KB of '(' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a parse that retried exponentially at the limit fails the
// test instead of hanging it.
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
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1+", n) + "1 FROM c" }},
	{"product", func(n int) string { return "SELECT " + strings.Repeat("1*", n) + "1 FROM c" }},
	{"concatenation", func(n int) string { return "SELECT " + strings.Repeat("'a'||", n) + "'a' FROM c" }},
	{"bitwise OR", func(n int) string { return "SELECT " + strings.Repeat("1|", n) + "1 FROM c" }},
	{"shift", func(n int) string { return "SELECT " + strings.Repeat("1<<", n) + "1 FROM c" }},
	{"comparisons", func(n int) string { return "SELECT * FROM c WHERE " + strings.Repeat("c.a = ", n) + "c.a" }},
	{"OR", func(n int) string { return "SELECT * FROM c WHERE " + strings.Repeat("c.a OR ", n) + "c.a" }},
	{"AND", func(n int) string { return "SELECT * FROM c WHERE " + strings.Repeat("c.a AND ", n) + "c.a" }},
	{"IN", func(n int) string { return "SELECT * FROM c WHERE c.a" + strings.Repeat(" IN (1)", n) }},
	{"LIKE", func(n int) string { return "SELECT * FROM c WHERE c.a" + strings.Repeat(" LIKE 'a'", n) }},
	{"BETWEEN", func(n int) string { return "SELECT * FROM c WHERE c.a" + strings.Repeat(" BETWEEN 1 AND 2", n) }},
	{"member access", func(n int) string { return "SELECT c" + strings.Repeat(".a", n) + " FROM c" }},
	{"brackets", func(n int) string { return "SELECT c" + strings.Repeat("[0]", n) + " FROM c" }},
	{"FROM member access", func(n int) string { return "SELECT * FROM c" + strings.Repeat(".a", n) }},
	{"FROM brackets", func(n int) string { return "SELECT * FROM c" + strings.Repeat("[0]", n) }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it (analysis did at a million terms). A chain within
// the limit still parses.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		_, err := Parse("  " + tc.sql(maxTreeHeight+1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingErrorMessage || pe.Position != 2 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, err)
		}
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
