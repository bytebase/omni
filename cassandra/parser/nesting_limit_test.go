package parser

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: about 880KB of '[' were enough. The same
// constructs nested 1,000 deep still parse. Each deep parse runs against a
// deadline, so a parse that retried exponentially at the limit fails the
// test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	insert := func(value string) string { return "INSERT INTO t (a) VALUES (" + value + ")" }
	for _, tc := range []struct {
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
	} {
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
