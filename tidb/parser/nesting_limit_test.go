package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	saved := maxNestingDepth
	maxNestingDepth = 1 << 20
	t.Cleanup(func() { maxNestingDepth = saved })
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: a few hundred kilobytes of '(' were enough. The
// same constructs nested 1,000 deep still parse. Each deep parse runs against
// a deadline, so a speculative parse that retried at every level past the
// limit fails the test instead of hanging it.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	proc := func(body string) string { return "CREATE PROCEDURE p() BEGIN " + body + " END" }
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
		{"unary plus", func(n int) string { return "SELECT " + strings.Repeat("+ ", n) + "1" }},
		{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "1" }},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
		{"assignment chain", func(n int) string { return "SELECT " + strings.Repeat("@a := ", n) + "1" }},
		{"function call", func(n int) string { return "SELECT " + nest(n, "abs(", "1", ")") }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS SIGNED)") }},
		{"CONVERT", func(n int) string { return "SELECT " + nest(n, "CONVERT(", "1", ", SIGNED)") }},
		{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "1", " THEN 1 END") }},
		{"ROW", func(n int) string { return "SELECT " + nest(n, "ROW(", "1", ", 1)") }},
		{"IN list", func(n int) string { return "SELECT " + nest(n, "1 IN (", "1", ")") }},
		{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
		{"EXISTS", func(n int) string { return nest(n, "SELECT EXISTS (", "SELECT 1", ")") }},
		{"IN subquery", func(n int) string { return nest(n, "SELECT 1 IN (", "SELECT 1", ")") }},
		{"subquery comparison", func(n int) string {
			return "SELECT " + strings.Repeat("(", n) + "SELECT 1" + strings.Repeat(") = 0", n)
		}},
		{"subquery operand", func(n int) string {
			return "SELECT " + nest(n, "((SELECT 1) LIKE 1 ESCAPE ", "1", ")")
		}},
		{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
		{"LATERAL derived table", func(n int) string {
			return nest(n, "SELECT * FROM t, LATERAL (", "SELECT 1", ") x")
		}},
		{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON 1)") }},
		{"parenthesized table", func(n int) string { return "SELECT * FROM " + nest(n, "(", "t", ")") }},
		{"parenthesized select", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
		{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
		{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
		{"JSON_TABLE nested columns", func(n int) string {
			return "SELECT * FROM JSON_TABLE('{}', '$' COLUMNS (" +
				nest(n, "NESTED PATH '$' COLUMNS (", "a INT PATH '$'", ")") + ")) AS jt"
		}},
		{"BEGIN block", func(n int) string { return proc(nest(n, "BEGIN ", "SELECT 1;", " END;")) }},
		{"IF", func(n int) string { return proc(nest(n, "IF 1 THEN ", "SELECT 1;", " END IF;")) }},
		{"CASE statement", func(n int) string {
			return proc(nest(n, "CASE WHEN 1 THEN ", "SELECT 1;", " END CASE;"))
		}},
		{"WHILE", func(n int) string { return proc(nest(n, "WHILE 1 DO ", "SELECT 1;", " END WHILE;")) }},
		{"REPEAT", func(n int) string { return proc(nest(n, "REPEAT ", "SELECT 1;", " UNTIL 1 END REPEAT;")) }},
		{"handler body", func(n int) string {
			return proc(nest(n, "BEGIN DECLARE CONTINUE HANDLER FOR SQLEXCEPTION ", "SELECT 1;", " END;"))
		}},
		{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
		{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
		{"unclosed subquery comparisons", func(n int) string {
			return "SELECT " + strings.Repeat("(", n) + "SELECT 1"
		}},
		{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
		{"unclosed parenthesized select", func(n int) string { return strings.Repeat("(", n) + "SELECT 1" }},
		{"unclosed blocks", func(n int) string { return proc(strings.Repeat("BEGIN ", n)) }},
	} {
		if err := parseWithin(tc.sql(1000)); err != nil && !strings.HasPrefix(tc.name, "unclosed") {
			t.Errorf("%s: 1,000 levels: %v", tc.name, briefErr(err))
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			err := parseWithin(tc.sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, briefErr(err))
			}
		}
	}
}

// parseWithin parses sql, then runs completion at its end, and fails when
// either takes over ten seconds.
func parseWithin(sql string) error {
	parsed := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		parsed <- err
	}()
	var err error
	select {
	case err = <-parsed:
	case <-time.After(10 * time.Second):
		return errors.New("Parse took over 10s")
	}
	collected := make(chan struct{})
	go func() {
		Collect(sql, len(sql))
		close(collected)
	}()
	select {
	case <-collected:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("Collect took over 10s")
	}
}

// briefErr renders err without the source line a ParseError carries, which
// for these inputs is the whole statement.
func briefErr(err error) string {
	var pe *ParseError
	if errors.As(err, &pe) {
		return fmt.Sprintf("%q at offset %d", pe.Message, pe.Position)
	}
	return fmt.Sprint(err)
}
