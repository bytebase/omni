package parser

import (
	"strings"
	"testing"
	"time"
)

// TestFromParensStayLinear guards the records behind the two readings of a
// FROM-clause '(' (parseTablePrimary): the base reading, a parenthesized
// expression or subquery, and the wrapped reading, `( tableReference )`,
// which runs when the base reading fails. Both readings parsed every '('
// nested inside, so where each level reads both ways the work doubled per
// level: a syntax error inside 18 nested derived tables took seconds, and so
// did the valid `((SELECT ...) CROSS JOIN t)` nested 18 deep. Where only the
// outer levels fail, as in a left-deep join group, each level reparsed the
// groups inside it, and 3,300 levels took seconds. These take milliseconds
// now; the doubling shapes would never return, so each parse runs against a
// deadline instead of being timed afterwards.
func TestFromParensStayLinear(t *testing.T) {
	// Each level of these costs a few levels of maxNestingDepth, so 3,300
	// is about as deep as they go.
	const depth = 3300
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"failing derived tables", nest(depth/4, "SELECT * FROM (", "SELECT * FROM t WHERE", ")"), true},
		{"derived tables failing after the alias", nest(depth/4, "SELECT * FROM (", "SELECT * FROM t", ") x y"), true},
		{"derived tables around a failing join group", nest(depth/4, "SELECT * FROM (", "t CROSS JOIN", ")"), true},
		{"joined subqueries", "SELECT * FROM " + nest(depth/4, "((SELECT * FROM ", "t", ") CROSS JOIN t)"), false},
		{"aliased subqueries", "SELECT * FROM " + nest(depth/4, "((SELECT * FROM ", "t", ") x)"), false},
		{"failing joined subqueries", "SELECT * FROM " + nest(depth/4, "((SELECT * FROM ", "t WHERE", ") CROSS JOIN t)"), true},
		{"failing expression", "SELECT * FROM " + nest(depth, "(", "t +", ")"), true},
		{"failing join group", "SELECT * FROM " + nest(depth, "(", "t CROSS JOIN", ")"), true},
		{"left-deep join group", "SELECT * FROM " + nest(depth, "(", "t", " CROSS JOIN t)"), false},
		{"left-deep aliased group", "SELECT * FROM " + nest(depth, "(", "t", ") x"), true},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			if (err != nil) != tc.wantErr || isNestingError(err) {
				t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse took over 2s", tc.name)
		}
	}
}
