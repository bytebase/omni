package parser

import (
	"strings"
	"testing"
	"time"
)

// TestNestedParensStayLinear guards the subquery validation: strict mode
// validates every subquery body, and validating a body by reparsing its text
// reparsed every nested body once per level, so 2,000 nested subqueries took
// seconds. Each shape below takes milliseconds when linear and minutes
// otherwise.
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
		{"right-deep scalar", "SELECT " + nest("((1 + ", "1", "))"), false},
		{"join group", "SELECT * FROM " + nest("(", "t", ")"), false},
		{"left-deep join group", "SELECT * FROM " + strings.Repeat("(", depth) + "t" + strings.Repeat(" JOIN t ON 1 = 1)", depth), false},
		{"scalar subquery", "SELECT " + nest("(SELECT ", "1", ")"), false},
		{"subquery comparison", "SELECT " + strings.Repeat("(", depth) + "SELECT 1" + strings.Repeat(") = 0", depth), false},
		{"EXISTS", "SELECT 1 FROM t WHERE " + nest("EXISTS (SELECT 1 FROM t WHERE ", "a = 1", ")"), false},
		{"unclosed EXISTS", "SELECT 1 FROM t WHERE " + strings.Repeat("EXISTS (SELECT 1 FROM t WHERE ", depth) + "a = 1", true},
		// Every subquery body is invalid at its end, after its nested body.
		{"invalid EXISTS", "SELECT 1 FROM t WHERE " + nest("EXISTS (SELECT 1 FROM t WHERE ", "a = 1", " x y)"), true},
		{"IN subquery", "SELECT 1 FROM t WHERE a IN " + nest("(SELECT a FROM t WHERE a IN ", "(1)", ")"), false},
		{"FROM subquery", "SELECT * FROM " + nest("(SELECT * FROM ", "t", ") x"), false},
	} {
		start := time.Now()
		if _, err := Parse(tc.sql); (err != nil) != tc.wantErr {
			t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
		}
		if elapsed := time.Since(start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse through %d nested parens took %v", tc.name, depth, elapsed)
		}
	}
}
