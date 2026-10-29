package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/starrocks/ast"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	saved := maxNestingDepth
	maxNestingDepth = 1 << 20
	t.Cleanup(func() { maxNestingDepth = saved })
}

// parseWithin runs Parse, and ParseBestEffort for the crash alone, against a
// deadline, so a parse whose retries blow up at the limit fails the test
// instead of hanging it. It returns Parse's error.
func parseWithin(t *testing.T, name, sql string) error {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		ParseBestEffort(sql)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("%s: Parse and ParseBestEffort took over 10s", name)
		return nil
	}
}

const nestingError = "statement is nested too deeply"

// errMessage returns the message of the first ParseError err carries.
func errMessage(err error) string {
	var pe *ParseError
	if errors.As(err, &pe) {
		return pe.Message
	}
	return ""
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
	{"lambda body", func(n int) string { return "SELECT " + nest(n, "array_map(x -> ", "x", ", a)") }},
	{"subscript", func(n int) string { return "SELECT " + nest(n, "a[", "1", "]") }},
	{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS INT)") }},
	{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "true", " THEN 1 END") }},
	{"array literal", func(n int) string { return "SELECT " + nest(n, "[", "1", "]") }},
	{"map literal", func(n int) string { return "SELECT " + nest(n, "MAP{'a': ", "1", "}") }},
	// A typed literal parses its type speculatively, backtracking to a
	// column comparison (array < ...) when no literal follows.
	{"typed array literal", func(n int) string { return "SELECT " + nest(n, "array<", "int", ">") + "[1]" }},
	{"typed map literal", func(n int) string { return "SELECT " + nest(n, "map<int, ", "int", ">") + "{}" }},
	{"ARRAY type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "ARRAY<", "INT", ">") + ")" }},
	{"MAP type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "MAP<INT, ", "INT", ">") + ")" }},
	{"STRUCT type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "STRUCT<a ", "INT", ">") + ")" }},
	{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
	{"EXISTS", func(n int) string { return nest(n, "SELECT 1 FROM t WHERE EXISTS (", "SELECT 1", ")") }},
	{"IN subquery", func(n int) string { return nest(n, "SELECT 1 FROM t WHERE a IN (", "SELECT a FROM t", ")") }},
	{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
	{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON true)") }},
	{"parenthesized FROM item", func(n int) string { return "SELECT * FROM " + nest(n, "(", "t", ")") }},
	{"parenthesized select", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
	{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
	{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
	// The CTAS query is captured raw and validated by a nested parse,
	// which starts at the depth it was captured at.
	{"CTAS query", func(n int) string { return "CREATE TABLE t AS " + nest(n, "SELECT (", "SELECT 1", ")") }},
	// A body that fails its in-place parse is scanned raw, through every
	// nested paren, and then validated by a nested parse.
	{"invalid raw-scanned subquery body", func(n int) string { return "SELECT (SELECT " + nest(n, "(", "1", ")") + " x y)" }},
	{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
	{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
	{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
	{"unclosed parenthesized selects", func(n int) string { return strings.Repeat("(", n) + "SELECT 1" }},
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: a few hundred kilobytes of '(' were enough. The
// same constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	for _, tc := range nestingShapes {
		// An unclosed or invalid shape is a syntax error at any depth: at
		// 1,000 levels it fails with that error, not the nesting one.
		invalid := strings.HasPrefix(tc.name, "unclosed") || strings.HasPrefix(tc.name, "invalid")
		err := parseWithin(t, tc.name+", 1,000 levels", tc.sql(1000))
		if invalid && (err == nil || errMessage(err) == nestingError) {
			t.Errorf("%s: 1,000 levels: err = %v, want a syntax error", tc.name, err)
		} else if !invalid && err != nil {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		want := nestingError
		if tc.name == "unclosed subqueries" {
			// Each level's body fails its in-place parse at the limit and
			// is then scanned raw, which meets the end of input first.
			want = "unterminated subquery"
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			if err := parseWithin(t, fmt.Sprintf("%s, %d levels", tc.name, n), tc.sql(n)); errMessage(err) != want {
				t.Errorf("%s: %d levels: err = %v, want %q", tc.name, n, err, want)
			}
		}
	}
}

// TestRawScanStaysIterative: a subquery body ParseBestEffort captures as raw
// text, or one whose in-place parse fails in strict mode, is scanned to its
// matching ')' without being parsed. The scan recursed once per nested '('
// and overflowed the stack on its own, with no production on the way to
// count the depth.
func TestRawScanStaysIterative(t *testing.T) {
	const n = 200000
	deep := strings.Repeat("(", n) + "1" + strings.Repeat(")", n)
	for _, tc := range []struct {
		name          string
		sql           string
		strictErr     string // the first error Parse reports
		bestEffortErr bool
	}{
		{"closed", "SELECT (SELECT " + deep + ")", nestingError, false},
		{"unclosed", "SELECT (SELECT " + strings.Repeat("(", n), "unterminated subquery", true},
		{"nested subqueries", "SELECT (" + strings.Repeat("SELECT (", n) + "1" + strings.Repeat(")", n+1), nestingError, false},
		// The in-place parse fails before the parens, so strict mode scans
		// them raw too, and the body's validation reports its syntax error.
		{"invalid before the parens", "SELECT (SELECT 1 x y " + deep + ")", "syntax error at or near y", false},
	} {
		type result struct {
			err        error
			bestEffort *ParseResult
		}
		done := make(chan result, 1)
		go func() {
			_, err := Parse(tc.sql)
			done <- result{err, ParseBestEffort(tc.sql)}
		}()
		select {
		case r := <-done:
			if got := errMessage(r.err); got != tc.strictErr {
				t.Errorf("%s: Parse err = %v, want %q", tc.name, r.err, tc.strictErr)
			}
			if (len(r.bestEffort.Errors) > 0) != tc.bestEffortErr {
				t.Errorf("%s: ParseBestEffort errors = %v, want error %v", tc.name, r.bestEffort.Errors, tc.bestEffortErr)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: Parse and ParseBestEffort took over 10s", tc.name)
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
	{"sum", func(n int) string { return "SELECT " + strings.Repeat("1+", n) + "1" }},
	{"bitwise OR", func(n int) string { return "SELECT " + strings.Repeat("1|", n) + "1" }},
	{"comparisons", func(n int) string { return "SELECT " + strings.Repeat("1 = ", n) + "1" }},
	{"OR", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a OR ", n) + "a" }},
	{"||", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a || ", n) + "a" }},
	{"AND", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a AND ", n) + "a" }},
	{"XOR", func(n int) string { return "SELECT 1 FROM t WHERE " + strings.Repeat("a XOR ", n) + "a" }},
	{"IS NULL", func(n int) string { return "SELECT a" + strings.Repeat(" IS NULL", n) }},
	{"LIKE", func(n int) string { return "SELECT a" + strings.Repeat(" LIKE 'x'", n) }},
	{"REGEXP", func(n int) string { return "SELECT a" + strings.Repeat(" REGEXP 'x'", n) }},
	{"IN", func(n int) string { return "SELECT a" + strings.Repeat(" IN (1)", n) }},
	{"BETWEEN", func(n int) string { return "SELECT a" + strings.Repeat(" BETWEEN 1 AND 2", n) }},
	{"subscripts", func(n int) string { return "SELECT a" + strings.Repeat("[1]", n) }},
	{"UNION ALL", func(n int) string { return strings.Repeat("SELECT 1 UNION ALL ", n) + "SELECT 1" }},
	{"UNION of parenthesized queries", func(n int) string { return strings.Repeat("(SELECT 1) UNION ", n) + "(SELECT 1)" }},
	{"INTERSECT", func(n int) string { return strings.Repeat("SELECT 1 INTERSECT ", n) + "SELECT 1" }},
	{"EXCEPT", func(n int) string { return strings.Repeat("SELECT 1 EXCEPT ", n) + "SELECT 1" }},
	{"joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" JOIN t ON true", n) }},
	{"cross joins", func(n int) string { return "SELECT * FROM t" + strings.Repeat(" CROSS JOIN t", n) }},
	{"comma joins in parens", func(n int) string { return "SELECT * FROM (t" + strings.Repeat(", t", n) + ")" }},
}

// TestTreeHeightLimit: a left-deep chain longer than maxTreeHeight is
// rejected, because every consumer that walks the tree recursively would
// overflow the stack on it. A chain within the limit still parses, and its
// measured height grows with its length: a link the measure skipped would
// let a longer chain through.
func TestTreeHeightLimit(t *testing.T) {
	for _, tc := range chainShapes {
		file, err := Parse(tc.sql(1000))
		if err != nil {
			t.Errorf("%s: 1,000 terms: %v", tc.name, err)
			continue
		}
		if h := treeHeight(file.Stmts[0]); h < 1000 {
			t.Errorf("%s: 1,000 terms measured %d deep", tc.name, h)
		}
		if _, err := Parse(tc.sql(maxTreeHeight / 4)); err != nil {
			t.Errorf("%s: %d terms: %v", tc.name, maxTreeHeight/4, err)
		}
		sql := tc.sql(maxTreeHeight + 1)
		_, err = Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != nestingError || pe.Position != 0 {
			t.Errorf("%s: %d terms: err = %v, want the nesting error at the statement", tc.name, maxTreeHeight+1, err)
		}
		// Best-effort parsing drops the statement and reports the same error.
		if res := ParseBestEffort(sql); len(res.File.Stmts) != 0 || len(res.Errors) == 0 || res.Errors[0].Message != nestingError {
			t.Errorf("%s: %d terms: ParseBestEffort kept %d statements, errors %v", tc.name, maxTreeHeight+1, len(res.File.Stmts), res.Errors)
		}
	}
}

// TestTreeHeightPerByte: treeTooHigh skips statements shorter than
// maxTreeHeight/16 bytes, on the ground that no construct adds more than 16
// levels of tree per byte of text. Check that across the recursive
// constructs and the chains.
func TestTreeHeightPerByte(t *testing.T) {
	check := func(name, sql string) {
		file, err := Parse(sql)
		if err != nil || len(file.Stmts) == 0 {
			return
		}
		if h := treeHeight(file.Stmts[0]); h > 16*len(sql) {
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

// TestTreeHeightCountsSubqueryBodies: Walk treats a subquery placeholder as
// a leaf, but analysis walks the body the strict parse built in place
// (SubqueryExpr.Query) from inside the statement's tree. A sum whose
// leftmost term is a subquery holding another sum is as deep as both
// together, so the two halves count together. ParseBestEffort keeps bodies
// as raw text, so only Parse sees it.
func TestTreeHeightCountsSubqueryBodies(t *testing.T) {
	sql := func(n int) string {
		return "SELECT (SELECT " + strings.Repeat("1+", n/2) + "1)" + strings.Repeat("+1", n/2)
	}
	if _, err := Parse(sql(maxTreeHeight / 4)); err != nil {
		t.Fatalf("%d terms: %v", maxTreeHeight/4, err)
	}
	_, err := Parse(sql(maxTreeHeight + 2))
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Message != nestingError || pe.Position != 0 {
		t.Fatalf("%d terms: err = %v, want the nesting error at the statement", maxTreeHeight+2, err)
	}
}
