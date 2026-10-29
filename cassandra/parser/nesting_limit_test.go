package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/cassandra/ast"
)

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

func insert(value string) string { return "INSERT INTO t (a) VALUES (" + value + ")" }

// nestingShapes are the constructs the parser descends into recursively,
// n levels deep.
var nestingShapes = []struct {
	name string
	cql  func(n int) string
}{
	{"list literal", func(n int) string { return insert(nest(n, "[", "1", "]")) }},
	{"set literal", func(n int) string { return insert(nest(n, "{", "1", "}")) }},
	{"map literal", func(n int) string { return insert(nest(n, "{1: ", "1", "}")) }},
	{"tuple literal", func(n int) string { return insert(nest(n, "(", "1", ")")) }},
	{"function call", func(n int) string { return insert(nest(n, "f(", "1", ")")) }},
	{"tuple function argument", func(n int) string { return insert(nest(n, "f((", "1", "))")) }},
	{"CAST", func(n int) string { return insert(nest(n, "CAST(", "1", " AS int)")) }},
	{"index access", func(n int) string { return insert(nest(n, "a[", "0", "]")) }},
	{"nested type", func(n int) string {
		return "CREATE TABLE t (a " + nest(n, "frozen<list<", "int", ">>") + " PRIMARY KEY)"
	}},
	{"INITCOND", func(n int) string {
		return "CREATE AGGREGATE agg(int) SFUNC f STYPE int INITCOND " + nest(n, "(", "0", ")")
	}},
	{"unclosed lists", func(n int) string { return "INSERT INTO t (a) VALUES (" + strings.Repeat("[", n) }},
	{"unclosed tuples", func(n int) string { return "INSERT INTO t (a) VALUES " + strings.Repeat("(", n) }},
	{"unclosed function calls", func(n int) string { return "INSERT INTO t (a) VALUES (" + strings.Repeat("f(", n) }},
	{"unclosed types", func(n int) string { return "CREATE TABLE t (a " + strings.Repeat("frozen<", n) }},
	{"unclosed INITCOND", func(n int) string {
		return "CREATE AGGREGATE agg(int) SFUNC f STYPE int INITCOND " + strings.Repeat("(", n)
	}},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 880KB of '[' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a parse that retried exponentially at the limit fails the
// test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		if _, err := Parse(tc.cql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			cql := tc.cql(n)
			done := make(chan error, 1)
			go func() {
				_, err := Parse(cql)
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
// a loop, so only the tree, not the parse, gets deep. Index access is the
// grammar's one such loop: WHERE relations, BATCH statements, and the other
// repeated clauses are kept in flat lists.
var chainShapes = []struct {
	name string
	cql  func(n int) string
}{
	{"index access", func(n int) string { return insert("a" + strings.Repeat("[0]", n)) }},
	{"index access in SET", func(n int) string {
		return "UPDATE t SET a = b" + strings.Repeat("[0]", n) + " WHERE k = 1"
	}},
}

// TestTreeHeightLimit: a chain longer than maxTreeHeight is rejected,
// because every consumer that walks the tree recursively would overflow the
// stack on it (ast.Walk did at three million terms). A chain within the
// limit still parses.
func TestTreeHeightLimit(t *testing.T) {
	const first = "USE ks; "
	for _, tc := range chainShapes {
		if _, err := Parse(tc.cql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		_, err := Parse(first + tc.cql(maxTreeHeight+1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingErrorMessage || pe.Position != len(first) {
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
	check := func(name, cql string, minHeight int) {
		list, err := Parse(cql)
		if err != nil || len(list.Items) == 0 {
			return
		}
		h := treeHeight(list.Items[0])
		if h > 16*len(cql) {
			t.Errorf("%s: tree %d deep from %d bytes", name, h, len(cql))
		}
		if h < minHeight {
			t.Errorf("%s: tree %d deep, want at least %d", name, h, minHeight)
		}
	}
	for _, tc := range nestingShapes {
		check(tc.name, tc.cql(1000), 0)
	}
	for _, tc := range chainShapes {
		check(tc.name, tc.cql(1000), 1000)
	}
}
