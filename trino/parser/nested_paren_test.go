package parser

import (
	"errors"
	"reflect"
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

// TestRawBodyChainsStayLinear guards the in-place parse of the subquery
// bodies whose parse in place used to fail while the parse of their text,
// which validates them, did not: a body over a lex error, one ending in white
// space the lexer rejects (the text is trimmed of it), and a WITH FUNCTION
// statement. Each was captured raw and validated by parsing its text, which
// validated every body inside it the same way, so 1,000 nested levels took
// seconds. These take milliseconds, and the statement fails with the error
// the innermost body reports, as before.
func TestRawBodyChainsStayLinear(t *testing.T) {
	raiseNestingLimit(t)
	nest := func(depth int, open, mid, closing string) string {
		return strings.Repeat(open, depth) + mid + strings.Repeat(closing, depth)
	}
	const routine = "WITH FUNCTION f() RETURNS int RETURN 1 SELECT "
	for _, tc := range []struct {
		name string
		sql  string
		at   string // the first error is at the first occurrence of at
		msg  string // with this message
	}{
		{"lex error under subqueries", "SELECT " + nest(10000, "(SELECT ", "! 1", ")"), "!", "expected expression, found INVALID"},
		{"lex error under EXISTS", "SELECT " + nest(10000, "EXISTS (SELECT ", "! 1", ")"), "!", "expected expression, found INVALID"},
		{"lex error in every subquery", "SELECT " + nest(10000, "(SELECT (SELECT !), ", "1", ")"), "!", "expected expression, found INVALID"},
		{"no-break space before every ')'", "SELECT " + nest(10000, "(SELECT ", "1", "  )"), " ", "unrecognized character"},
		{"vertical tab before every ')'", "SELECT " + nest(10000, "(SELECT ", "1", "\v)"), "\v", "unrecognized character"},
		{"WITH FUNCTION subqueries", "SELECT " + nest(5000, "("+routine, "1", ")"), "WITH FUNCTION f() RETURNS int RETURN 1 SELECT 1", "subquery must contain exactly one query"},
		{"WITH FUNCTION subqueries over a lex error", "SELECT " + nest(5000, "("+routine, "! 1", ")"), "!", "expected expression, found INVALID"},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			ParseBestEffort(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("%s: Parse error = %v, want a ParseError", tc.name, err)
			}
			if want := strings.Index(tc.sql, tc.at); pe.Position != want || pe.Message != tc.msg {
				t.Errorf("%s: first error %q at %d, want %q at %d", tc.name, pe.Message, pe.Position, tc.msg, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse and ParseBestEffort of %d bytes took over 2s", tc.name, len(tc.sql))
		}
	}
}

// TestInPlaceBodyErrors pins the errors of statements whose subquery and
// PREPARE bodies are validated by their in-place parse (see
// TestRawBodyChainsStayLinear and TestNestedPreparesStayLinear): every
// error, in the order parsing each body's text again reported them. A body's
// lex errors follow the errors of the bodies inside it; white space trimmed
// from a body is the enclosing text's; a PREPARE body's errors are reported
// as at their offset in the trimmed body.
func TestInPlaceBodyErrors(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []ParseError
	}{
		{"SELECT EXISTS (SELECT (SELECT 1 FROM t WHERE a IN (SELECT \u0085))), (SELECT ! 1)", []ParseError{
			{Position: 57, End: 57, Message: "expected expression, found end of input"},
			{Position: 58, End: 59, Message: "unrecognized character"},
			{Position: 59, End: 60, Message: "unrecognized character"},
			{Position: 73, End: 74, Message: "expected expression, found INVALID"},
			{Position: 73, End: 74, Message: "unrecognized character"},
		}},
		{"SELECT (SELECT (SELECT 1  )), (SELECT # 1)", []ParseError{
			{Position: 25, End: 26, Message: "unrecognized character"},
			{Position: 26, End: 27, Message: "unrecognized character"},
			{Position: 39, End: 40, Message: "expected expression, found INVALID"},
			{Position: 39, End: 40, Message: "unrecognized character"},
		}},
		{"SELECT (WITH FUNCTION f() RETURNS int RETURN 1 SELECT (WITH FUNCTION f() RETURNS int RETURN 1 SELECT 1))", []ParseError{
			{Position: 55, End: 102, Message: "subquery must contain exactly one query"},
		}},
		{"PREPARE q FROM PREPARE q FROM PREPARE q FROM TABLE s.t \"q\u0085", []ParseError{
			{Position: 15, End: 22, Message: strings.Repeat(prepareInvalid, 3) + "syntax error at or near INVALID"},
			{Position: 55, End: 59, Message: "unterminated quoted identifier"},
		}},
		{"PREPARE p FROM  PREPARE q FROM SELECT 1 1", []ParseError{
			{Position: 15, End: 16, Message: strings.Repeat(prepareInvalid, 2) + "syntax error at or near 1"},
			{Position: 15, End: 16, Message: "unrecognized character"},
			{Position: 16, End: 17, Message: "unrecognized character"},
		}},
		{"PREPARE p FROM SELECT 'ab  ", []ParseError{
			{Position: 15, End: 21, Message: prepareInvalid + "expected expression, found INVALID"},
			{Position: 22, End: 28, Message: "unterminated string literal"},
		}},
		{"PREPARE p FROM  /* c */ ", []ParseError{
			{Position: 15, End: 16, Message: "unrecognized character"},
			{Position: 16, End: 17, Message: "unrecognized character"},
			{Position: 24, End: 25, Message: "unrecognized character"},
			{Position: 25, End: 26, Message: "unrecognized character"},
		}},
	} {
		_, err := Parse(tc.sql)
		if got := AllErrors(err); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Parse(%q) errors:\n got %+v\nwant %+v", tc.sql, got, tc.want)
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
