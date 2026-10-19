package parser

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 880KB of '(' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a parse that retried exponentially at the limit fails the
// test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") + " FROM c" }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1 FROM c" }},
		{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true FROM c" }},
		{"bitwise NOT", func(n int) string { return "SELECT " + strings.Repeat("~", n) + "1 FROM c" }},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") + " FROM c" }},
		{"ternary", func(n int) string { return "SELECT " + strings.Repeat("true ? 1 : ", n) + "1 FROM c" }},
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
