package parser

import (
	"strings"
	"testing"
	"time"
)

// TestNestedParensStayLinear guards the memos behind the FROM-clause paren
// speculation (parenQueryFail and parenRelation): the parser tries a query at
// each '(' of a relation before it parses a relation, and answering by
// reparsing made 10,000 nested parens take seconds and a nested aliased
// subquery double in cost per level. Each shape below takes milliseconds when
// linear. A regressed shape may never return, so each parse runs against the
// deadline instead of being timed afterwards.
func TestNestedParensStayLinear(t *testing.T) {
	const depth = 20000
	nest := func(open, mid, closing string) string {
		return strings.Repeat(open, depth) + mid + strings.Repeat(closing, depth)
	}
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"join group", "SELECT * FROM " + nest("(", "t", ")"), false},
		{"unclosed join group", "SELECT * FROM " + strings.Repeat("(", depth) + "t", true},
		{"left-deep join group", "SELECT * FROM " + strings.Repeat("(", depth) + "t" + strings.Repeat(" JOIN t ON 1)", depth), false},
		{"right-deep join group", "SELECT * FROM " + nest("((t JOIN ", "t", " ON 1))"), false},
		{"derived table", "SELECT * FROM " + nest("(", "SELECT 1", ")"), false},
		{"aliased derived table", "SELECT * FROM " + nest("(", "(SELECT 1) x", ")"), false},
		{"aliased subquery chain", "SELECT * FROM " + nest("((SELECT * FROM ", "t", ") x)"), false},
		{"failing subquery chain", "SELECT * FROM " + nest("((SELECT * FROM ", "t WHERE", ") x)"), true},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			if (err != nil) != tc.wantErr {
				t.Errorf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse through %d nested parens took over 2s", tc.name, depth)
		}
	}
}
