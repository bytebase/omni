package parser

import (
	"strings"
	"testing"
	"time"
)

// TestNestedParensStayLinear guards the memos behind the paren speculation and
// the subquery validation: the parser asks what a '(' opens at every nesting
// level as it descends, and answering by reparsing made 2,000 nested parens
// take seconds. Each shape below takes milliseconds when linear and minutes
// otherwise; the subquery chains doubled per level and never return, so each
// parse runs against the deadline instead of being timed afterwards.
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
		{"unclosed join group", "SELECT * FROM " + strings.Repeat("(", depth) + "t", true},
		{"right-deep join group", "SELECT * FROM " + nest("((t JOIN ", "t", " ON 1))"), false},
		{"left-deep join group", "SELECT * FROM " + strings.Repeat("(", depth) + "t" + strings.Repeat(" JOIN t ON 1)", depth), false},
		{"derived table", "SELECT * FROM " + nest("(", "SELECT 1", ")"), false},
		{"aliased derived table", "SELECT * FROM " + nest("(", "(SELECT 1) x", ")"), false},
		{"aliased subquery chain", "SELECT * FROM " + nest("((SELECT * FROM ", "t", ") x)"), false},
		{"joined subquery chain", "SELECT * FROM " + nest("((SELECT * FROM ", "t", ") JOIN t ON 1)"), false},
		{"failing subquery chain", "SELECT * FROM " + nest("((SELECT * FROM ", "t WHERE", ") x)"), true},
		{"subquery chain with raw bodies", "SELECT * FROM " + nest("((SELECT (SELECT 1 1) FROM ", "t", ") x)"), true},
		{"lex errors in raw groups", "SELECT EXISTS (SELECT " + nest("(! ", "1", ")") + ")", true},
		{"EXISTS", "SELECT " + nest("EXISTS (SELECT (", "1", "))"), false},
		{"unclosed EXISTS", "SELECT " + strings.Repeat("EXISTS (SELECT (", depth) + "1", true},
		// Every subquery body is invalid at its end, after its nested body.
		{"invalid EXISTS", "SELECT " + nest("EXISTS (SELECT (", "1", ") x y)"), true},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			ParseBestEffort(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			if (err != nil) != tc.wantErr {
				t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse and ParseBestEffort through %d nested parens took over 2s", tc.name, depth)
		}
	}
}

// TestSpeculativeReadingsStayLinear guards the records behind the readings the
// parser abandons for another reading of the same tokens: SUBSTRING and
// POSITION, TRIM, a JSON_OBJECT member's KEY, and a ROW type's field list.
// Each level of nesting read the level inside it both ways, so 24 levels took
// 20 seconds, doubling with each level added; these finish in milliseconds.
func TestSpeculativeReadingsStayLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	nest := func(open, mid, closing string) string {
		return strings.Repeat(open, depth) + mid + strings.Repeat(closing, depth)
	}
	castTo := func(typ string) string { return "SELECT CAST(NULL AS " + typ + ")" }
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"SUBSTRING call", "SELECT " + nest("substring(", "x", ", 1)"), false},
		{"SUBSTRING call with a predicate", "SELECT " + nest("substring(", "x", " IS DISTINCT FROM 1, 1)"), false},
		{"failing SUBSTRING call", "SELECT " + nest("substring(", "x +", ", 1)"), true},
		{"failing SUBSTRING FROM", "SELECT " + nest("substring(", "x", " FROM 1 FOR)"), true},
		{"SUBSTRING with raw bodies", "SELECT " + nest("substring(", "x", ", (SELECT 1 1))"), true},
		{"POSITION call", "SELECT " + nest("position(", "'a'", ", 'b')"), false},
		{"failing POSITION IN", "SELECT " + nest("position(", "'a'", " IN)"), true},
		{"TRIM call", "SELECT " + nest("trim(", "x", ", 'a')"), false},
		{"failing TRIM call", "SELECT " + nest("trim(", "x", ", )"), true},
		{"JSON_OBJECT KEY expression", "SELECT " + nest("JSON_OBJECT(KEY - ", "1", " : 1)"), false},
		{"JSON_OBJECT KEY call", "SELECT " + nest("JSON_OBJECT(KEY (", "1", ") : 1)"), false},
		{"failing JSON_OBJECT KEY", "SELECT " + nest("JSON_OBJECT(KEY - ", "1", " VALUE)"), true},
		{"ROW type with integers", castTo(nest("ROW(", "bigint", ", 1)")), false},
		{"ROW type with a trailing name", castTo(nest("ROW(", "bigint", " x)")), true},
		{"failing ROW type", castTo(nest("ROW(", "bigint", ", )")), true},
		{"routine parameter ROW type", "CREATE FUNCTION f(a " + nest("ROW(", "bigint", ", 1)") + ") RETURNS bigint RETURN 1", false},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			ParseBestEffort(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			if (err != nil) != tc.wantErr {
				t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse and ParseBestEffort through %d nested levels took over 2s", tc.name, depth)
		}
	}
}
