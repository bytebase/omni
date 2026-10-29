package parser

import (
	"errors"
	"strings"
	"testing"
)

// TestNestingLimit: nested blocks and control statements stop at
// maxNestingDepth with an error instead of overflowing the Go stack, which
// killed the process. The same constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return "BEGIN " + strings.Repeat(open, n) + mid + strings.Repeat(closing, n) + " END"
	}
	for _, tc := range []struct {
		name string
		body func(n int) string
	}{
		{"block", func(n int) string { return nest(n, "BEGIN ", "NULL;", " END;") }},
		{"IF", func(n int) string { return nest(n, "IF true THEN ", "NULL;", " END IF;") }},
		{"LOOP", func(n int) string { return nest(n, "LOOP ", "EXIT;", " END LOOP;") }},
		{"WHILE", func(n int) string { return nest(n, "WHILE true LOOP ", "EXIT;", " END LOOP;") }},
		{"CASE", func(n int) string { return nest(n, "CASE WHEN true THEN ", "NULL;", " END CASE;") }},
		{"exception handler", func(n int) string {
			return nest(n, "BEGIN NULL; EXCEPTION WHEN others THEN ", "NULL;", " END;")
		}},
		{"unclosed blocks", func(n int) string { return strings.Repeat("BEGIN ", n) }},
	} {
		if _, err := Parse(tc.body(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, err)
		}
		for _, n := range []int{maxNestingDepth + 1, 300000} {
			_, err := Parse(tc.body(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}
