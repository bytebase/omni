package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/oceanbase/ast"
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

func proc(body string) string { return "CREATE PROCEDURE p() BEGIN " + body + " END" }

// nestingShapes are the constructs the parser descends into recursively,
// n levels deep.
var nestingShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
	{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
	{"unary plus", func(n int) string { return "SELECT " + strings.Repeat("+ ", n) + "1" }},
	{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "1" }},
	{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
	{"assignment chain", func(n int) string { return "SELECT " + strings.Repeat("@a := ", n) + "1" }},
	{"function call", func(n int) string { return "SELECT " + nest(n, "abs(", "1", ")") }},
	{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS SIGNED)") }},
	{"CONVERT", func(n int) string { return "SELECT " + nest(n, "CONVERT(", "1", ", SIGNED)") }},
	{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "1", " THEN 1 END") }},
	{"ROW", func(n int) string { return "SELECT " + nest(n, "ROW(", "1", ", 1)") }},
	{"row constructor", func(n int) string { return "SELECT " + nest(n, "(", "1", ", 1)") }},
	{"IN list", func(n int) string { return "SELECT " + nest(n, "1 IN (", "1", ")") }},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
	{"EXISTS", func(n int) string { return nest(n, "SELECT EXISTS (", "SELECT 1", ")") }},
	{"IN subquery", func(n int) string { return nest(n, "SELECT 1 IN (", "SELECT 1", ")") }},
	{"quantified subquery", func(n int) string { return nest(n, "SELECT 1 = ANY (", "SELECT 1", ")") }},
	{"subquery comparison", func(n int) string {
		return "SELECT " + strings.Repeat("(", n) + "SELECT 1" + strings.Repeat(") = 0", n)
	}},
	{"subquery operand", func(n int) string {
		return "SELECT " + nest(n, "((SELECT 1) LIKE 1 ESCAPE ", "1", ")")
	}},
	{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
	{"LATERAL derived table", func(n int) string {
		return nest(n, "SELECT * FROM t, LATERAL (", "SELECT 1", ") x")
	}},
	{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON 1)") }},
	{"parenthesized table", func(n int) string { return "SELECT * FROM " + nest(n, "(", "t", ")") }},
	{"parenthesized select", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
	{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
	// A TABLE or VALUES source in the ORDER BY of another: while the source
	// and the SelectStmt over it both held the ORDER BY, the tree-height walk
	// doubled at each level.
	{"TABLE in ORDER BY", func(n int) string { return nest(n, "TABLE t ORDER BY (", "SELECT 1", ")") }},
	{"VALUES in ORDER BY", func(n int) string { return nest(n, "VALUES ROW(1) ORDER BY (", "SELECT 1", ")") }},
	{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
	{"JSON_TABLE nested columns", func(n int) string {
		return "SELECT * FROM JSON_TABLE('{}', '$' COLUMNS (" +
			nest(n, "NESTED PATH '$' COLUMNS (", "a INT PATH '$'", ")") + ")) AS jt"
	}},
	{"BEGIN block", func(n int) string { return proc(nest(n, "BEGIN ", "SELECT 1;", " END;")) }},
	{"IF", func(n int) string { return proc(nest(n, "IF 1 THEN ", "SELECT 1;", " END IF;")) }},
	{"CASE statement", func(n int) string {
		return proc(nest(n, "CASE WHEN 1 THEN ", "SELECT 1;", " END CASE;"))
	}},
	{"WHILE", func(n int) string { return proc(nest(n, "WHILE 1 DO ", "SELECT 1;", " END WHILE;")) }},
	{"REPEAT", func(n int) string { return proc(nest(n, "REPEAT ", "SELECT 1;", " UNTIL 1 END REPEAT;")) }},
	{"LOOP", func(n int) string { return proc(nest(n, "LOOP ", "SELECT 1;", " END LOOP;")) }},
	{"handler body", func(n int) string {
		return proc(nest(n, "BEGIN DECLARE CONTINUE HANDLER FOR SQLEXCEPTION ", "SELECT 1;", " END;"))
	}},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
	{"unclosed subquery comparisons", func(n int) string {
		return "SELECT " + strings.Repeat("(", n) + "SELECT 1"
	}},
	{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
	{"unclosed parenthesized select", func(n int) string { return strings.Repeat("(", n) + "SELECT 1" }},
	{"unclosed blocks", func(n int) string { return proc(strings.Repeat("BEGIN ", n)) }},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: a few hundred kilobytes of '(' were enough. The
// same constructs nested 1,000 deep still parse. Each deep parse runs against
// a deadline, so a speculative parse that retried at every level past the
// limit fails the test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		if err := parseWithin(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, briefErr(err))
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			err := parseWithin(tc.sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, briefErr(err))
			}
		}
	}
}

// parseWithin parses sql, then runs completion at its end, and fails when
// either takes over ten seconds.
func parseWithin(sql string) error {
	parsed := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		parsed <- err
	}()
	var err error
	select {
	case err = <-parsed:
	case <-time.After(10 * time.Second):
		return errors.New("Parse took over 10s")
	}
	collected := make(chan struct{})
	go func() {
		Collect(sql, len(sql))
		close(collected)
	}()
	select {
	case <-collected:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("Collect took over 10s")
	}
}

// briefErr renders err without the source line a ParseError carries, which
// for these inputs is the whole statement.
func briefErr(err error) string {
	var pe *ParseError
	if errors.As(err, &pe) {
		return fmt.Sprintf("%q at offset %d", pe.Message, pe.Position)
	}
	return fmt.Sprint(err)
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
	{"OR", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a OR ", n) + "a" }},
	{"AND", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a AND ", n) + "a" }},
	{"XOR", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a XOR ", n) + "a" }},
	{"comparisons", func(n int) string { return "SELECT 1" + strings.Repeat(" = 1", n) }},
	{"JSON access", func(n int) string { return "SELECT a" + strings.Repeat("->'$'", n) + " FROM t" }},
	{"COLLATE", func(n int) string { return "SELECT 'a'" + strings.Repeat(" COLLATE utf8mb4_bin", n) }},
	{"IS NULL", func(n int) string { return "SELECT a" + strings.Repeat(" IS NULL", n) + " FROM t" }},
	{"BETWEEN", func(n int) string { return "SELECT a" + strings.Repeat(" BETWEEN 1 AND 2", n) + " FROM t" }},
	{"IN", func(n int) string { return "SELECT a" + strings.Repeat(" IN (1)", n) + " FROM t" }},
	{"LIKE", func(n int) string { return "SELECT a" + strings.Repeat(" LIKE 'a'", n) + " FROM t" }},
	{"REGEXP", func(n int) string { return "SELECT a" + strings.Repeat(" REGEXP 'a'", n) + " FROM t" }},
	{"SOUNDS LIKE", func(n int) string { return "SELECT a" + strings.Repeat(" SOUNDS LIKE 'a'", n) + " FROM t" }},
	{"MEMBER OF", func(n int) string { return "SELECT 1" + strings.Repeat(" MEMBER OF ('[1]')", n) }},
	{"quantified comparisons", func(n int) string { return "SELECT 1" + strings.Repeat(" = ANY (SELECT 1)", n) }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT 1 UNION ALL ", n) + "SELECT 1" }},
	{"INTERSECT", func(n int) string { return strings.Repeat("SELECT 1 INTERSECT ", n) + "SELECT 1" }},
	{"joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON 1", n) }},
	{"parenthesized table list", func(n int) string { return "SELECT * FROM (" + strings.Repeat("t, ", n) + "t)" }},
	{"INSERT VALUES", func(n int) string { return "INSERT INTO t VALUES (" + strings.Repeat("1+", n) + "1)" }},
	{"VALUES", func(n int) string { return "VALUES ROW(" + strings.Repeat("1+", n) + "1)" }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it. A chain within the limit still parses, and its
// measured height grows with the chain, so the measure sees the levels the
// consumers recurse through.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		list, err := Parse(tc.sql(1000))
		if err != nil {
			t.Errorf("%s: 1,000 terms: %v", tc.name, briefErr(err))
			continue
		}
		if h := treeHeight(list.Items[0]); h <= 1000 {
			t.Errorf("%s: 1,000 terms: tree %d deep, want over 1,000", tc.name, h)
		}
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, briefErr(err))
		}
		_, err = Parse(tc.sql(maxTreeHeight + 1))
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, briefErr(err))
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
