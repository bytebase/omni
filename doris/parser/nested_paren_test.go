package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/doris/ast"
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

// TestParenthesizedSetOpsStayLinear: each '(' around a query asks whether
// the query's leading SELECT has a WITH, which keeps the group as a
// GroupedQuery scope boundary. Every enclosing '(' walked the set operations
// inside it again, so 9,990 parens around a 90,000-arm chain took 20s.
func TestParenthesizedSetOpsStayLinear(t *testing.T) {
	const parens, arms = 9000, 60000
	chain := strings.Repeat(" UNION ALL SELECT 1", arms-1)
	for _, tc := range []struct {
		name    string
		sql     string
		grouped bool // the statement keeps a GroupedQuery for the WITH
	}{
		{"parens around a chain", strings.Repeat("(", parens) + "SELECT 1" + chain + strings.Repeat(")", parens), false},
		{"WITH inside parens around a chain", strings.Repeat("(", parens) + "WITH c AS (SELECT 1) SELECT * FROM c" + chain + strings.Repeat(")", parens), true},
		{"parenthesized left operands", strings.Repeat("(", parens) + "SELECT 1" + strings.Repeat(") UNION ALL SELECT 1", parens), false},
	} {
		if err := parseWithin(t, tc.name, tc.sql); err != nil {
			t.Fatalf("%s: Parse error = %v", tc.name, err)
		}
		file, _ := Parse(tc.sql)
		if _, ok := file.Stmts[0].(*ast.GroupedQuery); ok != tc.grouped {
			t.Errorf("%s: statement is %T, want GroupedQuery %v", tc.name, file.Stmts[0], tc.grouped)
		}
	}
}
