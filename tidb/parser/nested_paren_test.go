package parser

import (
	"strings"
	"testing"
	"time"
)

// TestNestedParensStayLinear guards the memos behind the paren lookaheads: the
// parser asks what a run of '(' opens at every nesting level as it descends,
// and answering by rescanning made 10,000 nested parens take seconds. Each
// shape below takes milliseconds when linear and tens of seconds otherwise.
// The budget below is CPU time, not clock time.
func TestNestedParensStayLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	nest := func(open, mid, closing string) string {
		return strings.Repeat(open, depth) + mid + strings.Repeat(closing, depth)
	}
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"scalar", "SELECT " + nest("(", "1", ")"), false},
		{"unclosed scalar", "SELECT " + strings.Repeat("(", depth) + "1", true},
		{"condition", "SELECT * FROM t WHERE " + nest("(", "a = 1", ")"), false},
		{"IN list", "SELECT * FROM t WHERE a IN " + nest("(", "1", ")"), false},
		{"right-deep scalar", "SELECT " + nest("((1 + ", "1", "))"), false},
		{"unclosed right-deep scalar", "SELECT " + strings.Repeat("((1 + ", depth) + "1", true},
		{"subquery", "SELECT " + nest("(", "SELECT 1", ")"), false},
		{"subquery comparison", "SELECT " + strings.Repeat("(", depth) + "SELECT 1" + strings.Repeat(") = 0", depth), false},
		{"join group", "SELECT * FROM " + nest("(", "t", ")"), false},
		{"right-deep join group", "SELECT * FROM " + nest("((t JOIN ", "t", " ON 1))"), false},
	} {
		start := cpuTime(t)
		if _, err := Parse(tc.sql); (err != nil) != tc.wantErr {
			t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
		}
		if elapsed := cpuSince(t, start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse through %d nested parens took %v of CPU time", tc.name, depth, elapsed)
		}
	}
}
