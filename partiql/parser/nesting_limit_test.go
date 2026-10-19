package parser

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 300KB of '(' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a speculative parse that retried exponentially at the limit
// fails the test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	for _, tc := range []struct {
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
	} {
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
