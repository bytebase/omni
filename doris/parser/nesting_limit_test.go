package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
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
		{"lambda body", func(n int) string { return "SELECT " + nest(n, "array_map(x -> ", "x", ", a)") }},
		{"subscript", func(n int) string { return "SELECT " + nest(n, "a[", "1", "]") }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS INT)") }},
		{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "true", " THEN 1 END") }},
		{"array literal", func(n int) string { return "SELECT " + nest(n, "[", "1", "]") }},
		{"map literal", func(n int) string { return "SELECT " + nest(n, "{'a': ", "1", "}") }},
		{"struct literal", func(n int) string { return "SELECT " + nest(n, "{1, ", "1", "}") }},
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
	} {
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
