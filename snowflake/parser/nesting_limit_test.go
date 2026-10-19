package parser

import (
	"errors"
	"strings"
	"testing"
	"time"
)

const nestingErr = "statement is nested too deeply"

// parseWithin parses sql with Parse and fails the test if that takes over 10s,
// so a retry that turns exponential at the limit fails instead of hanging.
func parseWithin(t *testing.T, name string, n int, sql string) error {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("%s: %d levels: Parse took over 10s", name, n)
		return nil
	}
}

func nest(n int, open, mid, closing string) string {
	return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: 470KB of unclosed '(' were enough. The same
// constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
		{"unary plus", func(n int) string { return "SELECT " + strings.Repeat("+ ", n) + "1" }},
		{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true" }},
		{"PRIOR", func(n int) string {
			return "SELECT a FROM t CONNECT BY " + strings.Repeat("PRIOR ", n) + "a = b"
		}},
		{"CONNECT_BY_ROOT", func(n int) string {
			return "SELECT " + strings.Repeat("CONNECT_BY_ROOT ", n) + "a FROM t CONNECT BY PRIOR a = b"
		}},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
		{"function call", func(n int) string { return "SELECT " + nest(n, "f(", "1", ")") }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS INT)") }},
		{"TRY_CAST", func(n int) string { return "SELECT " + nest(n, "TRY_CAST(", "'1'", " AS INT)") }},
		{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN true THEN ", "1", " END") }},
		{"IFF", func(n int) string { return "SELECT " + nest(n, "IFF(true, ", "1", ", 0)") }},
		{"array literal", func(n int) string { return "SELECT " + nest(n, "[", "1", "]") }},
		{"object literal", func(n int) string { return "SELECT " + nest(n, "{'a': ", "1", "}") }},
		{"lambda", func(n int) string { return "SELECT f(" + strings.Repeat("x -> ", n) + "1)" }},
		{"subscript", func(n int) string { return "SELECT " + nest(n, "a[", "1", "]") }},
		{"IN list", func(n int) string { return "SELECT " + nest(n, "1 IN (", "1", ")") }},
		{"window ORDER BY", func(n int) string { return "SELECT " + nest(n, "f() OVER (ORDER BY ", "a", ")") }},
		{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
		{"EXISTS", func(n int) string { return nest(n, "SELECT 1 WHERE EXISTS (", "SELECT 1", ")") }},
		{"IN subquery", func(n int) string { return nest(n, "SELECT 1 WHERE 1 IN (", "SELECT 1", ")") }},
		{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ")") }},
		{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON true)") }},
		{"LATERAL", func(n int) string { return "SELECT * FROM t, " + strings.Repeat("LATERAL ", n) + "u" }},
		{"parenthesized query", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
		{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION (", "SELECT 1", ")") }},
		{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT 1") }},
		{"nested type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "ARRAY(", "INT", ")") + ")" }},
		{"column type", func(n int) string { return "CREATE TABLE t (c " + nest(n, "ARRAY(", "INT", ")") + ")" }},
		{"option group", func(n int) string {
			return "COPY INTO t FROM @s FILE_FORMAT = " + nest(n, "(TYPE = ", "CSV", ")")
		}},
		{"EXPLAIN", func(n int) string { return strings.Repeat("EXPLAIN ", n) + "SELECT 1" }},
		{"expression in a scripting block", func(n int) string {
			return "BEGIN RETURN " + nest(n, "(", "1", ")") + "; END"
		}},
		{"query in a scripting block", func(n int) string {
			return "BEGIN " + nest(n, "SELECT (", "SELECT 1", ")") + "; END"
		}},
		{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
		{"unclosed function calls", func(n int) string { return "SELECT " + strings.Repeat("f(", n) }},
		{"unclosed array literals", func(n int) string { return "SELECT " + strings.Repeat("[", n) }},
		{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
		{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
		{"unclosed parenthesized queries", func(n int) string { return strings.Repeat("(", n) }},
	} {
		if !strings.HasPrefix(tc.name, "unclosed") {
			if err := parseWithin(t, tc.name, 1000, tc.sql(1000)); err != nil {
				t.Errorf("%s: 1,000 levels: %v", tc.name, err)
			}
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			err := parseWithin(t, tc.name, n, tc.sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != nestingErr {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// TestNestingLimitEntryPoints: ParseExpr and ParseDataType start their own
// parser and are bounded the same way.
func TestNestingLimitEntryPoints(t *testing.T) {
	for _, n := range []int{maxNestingDepth + 1, 200000} {
		if _, errs := ParseExpr(nest(n, "(", "1", ")")); len(errs) == 0 || errs[0].Message != nestingErr {
			t.Errorf("ParseExpr: %d levels: errs = %v, want the nesting error", n, errs)
		}
		if _, errs := ParseDataType(nest(n, "ARRAY(", "INT", ")")); len(errs) == 0 || errs[0].Message != nestingErr {
			t.Errorf("ParseDataType: %d levels: errs = %v, want the nesting error", n, errs)
		}
	}
}

// TestNestingLimitLeavesScriptDepth: Snowflake Scripting blocks keep their
// own, lower ceiling (maxScriptDepth) and its error, however deep they nest.
func TestNestingLimitLeavesScriptDepth(t *testing.T) {
	for name, sql := range map[string]func(n int) string{
		"blocks":       func(n int) string { return nest(n, "BEGIN ", "RETURN 1;", " END;") },
		"IF":           func(n int) string { return "BEGIN " + nest(n, "IF (true) THEN ", "RETURN 1;", " END IF;") + " END" },
		"unclosed IFs": func(n int) string { return "BEGIN " + strings.Repeat("IF (true) THEN ", n) },
	} {
		for _, n := range []int{maxScriptDepth + 1, maxNestingDepth + 1, 200000} {
			err := parseWithin(t, name, n, sql(n))
			var pe *ParseError
			if !errors.As(err, &pe) || pe.Message != "snowflake scripting block nested too deeply" {
				t.Errorf("%s: %d levels: err = %v, want the scripting depth error", name, n, err)
			}
		}
	}
}

// TestNestingLimitInEmbeddedBody: a task body that fails to parse is kept
// verbatim, and one nested past the limit is no different. Nesting tasks in
// tasks must not make that fallback retry exponentially.
func TestNestingLimitInEmbeddedBody(t *testing.T) {
	for _, n := range []int{maxNestingDepth + 1, 100000} {
		sql := strings.Repeat("CREATE TASK t AS ", n) + "SELECT 1"
		if err := parseWithin(t, "nested tasks", n, sql); err != nil {
			t.Errorf("nested tasks: %d levels: %v", n, err)
		}
	}
}
