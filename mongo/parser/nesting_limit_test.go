package parser

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNestingLimit: every value the parser descends into recursively stops
// at maxNestingDepth with an error instead of overflowing the Go stack, which
// killed the process: about 880KB of '[' were enough. The same values nested
// 1,000 deep still parse. Each deep parse runs against a deadline, so a parse
// that retried exponentially at the limit fails the test instead of hanging
// it; ParseBestEffort, which skips past the failed statement, runs too.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	find := func(arg string) string { return "db.c.find(" + arg + ")" }
	for _, tc := range []struct {
		name  string
		input func(n int) string
	}{
		{"document", func(n int) string { return find(nest(n, "{a: ", "1", "}")) }},
		{"array", func(n int) string { return find(nest(n, "[", "1", "]")) }},
		{"documents in arrays", func(n int) string { return find(nest(n, "[{a: ", "1", "}]")) }},
		{"helper call", func(n int) string { return find(nest(n, "NumberInt(", "1", ")")) }},
		{"identifier helper call", func(n int) string { return find(nest(n, "Long(", "1", ")")) }},
		{"method call", func(n int) string { return find(nest(n, "x.f(", "1", ")")) }},
		{"native function argument", func(n int) string { return "sleep(" + nest(n, "[", "1", "]") + ")" }},
		{"unclosed documents", func(n int) string { return "db.c.find(" + strings.Repeat("{a: ", n) }},
		{"unclosed arrays", func(n int) string { return "db.c.find(" + strings.Repeat("[", n) }},
		{"unclosed helper calls", func(n int) string { return "db.c.find(" + strings.Repeat("NumberInt(", n) }},
	} {
		if _, err := Parse(tc.input(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			input := tc.input(n)
			done := make(chan error, 1)
			go func() {
				_, err := Parse(input)
				ParseBestEffort(input)
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
