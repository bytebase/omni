package parser

import (
	"strings"
	"testing"
	"time"
)

// raiseNestingLimit lifts maxNestingDepth for a test that nests deeper on
// purpose, and restores it when the test ends.
func raiseNestingLimit(t *testing.T) {
	setNestingLimit(t, 1<<20)
}

// setNestingLimit sets maxNestingDepth to n until the test ends.
func setNestingLimit(t *testing.T, n int) {
	saved := maxNestingDepth
	maxNestingDepth = n
	t.Cleanup(func() { maxNestingDepth = saved })
}

// hasNestingError reports whether err carries the error enterNesting returns.
func hasNestingError(err error) bool {
	for _, e := range AllErrors(err) {
		if e.Message == nestingErrorMessage {
			return true
		}
	}
	return false
}

// parseWithin runs Parse on sql and reports whether it returned within limit,
// so a parse that retries exponentially at the nesting limit fails the test
// instead of hanging the run.
func parseWithin(sql string, limit time.Duration) (finished bool, err error) {
	done := make(chan error, 1)
	go func() {
		_, err := Parse(sql)
		done <- err
	}()
	select {
	case err := <-done:
		return true, err
	case <-time.After(limit):
		return false, nil
	}
}

// TestNestingLimit: every construct the parser descends into recursively
// stops at maxNestingDepth with an error instead of overflowing the Go stack,
// which killed the process: half a megabyte of '(' was enough. The same
// constructs nested 1,000 deep still parse.
func TestNestingLimit(t *testing.T) {
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"parenthesized expression", func(n int) string { return "SELECT " + nest(n, "(", "1", ")") }},
		{"unary minus", func(n int) string { return "SELECT " + strings.Repeat("- ", n) + "1" }},
		{"NOT", func(n int) string { return "SELECT " + strings.Repeat("NOT ", n) + "true" }},
		{"right-deep operators", func(n int) string { return "SELECT " + nest(n, "1 + (", "1", ")") }},
		{"function call", func(n int) string { return "SELECT " + nest(n, "f(", "1", ")") }},
		{"CAST", func(n int) string { return "SELECT " + nest(n, "CAST(", "1", " AS INT64)") }},
		{"CASE", func(n int) string { return "SELECT " + nest(n, "CASE WHEN ", "true", " THEN 1 END") }},
		{"array constructor", func(n int) string { return "SELECT " + nest(n, "ARRAY[", "1", "]") }},
		{"STRUCT constructor", func(n int) string { return "SELECT " + nest(n, "STRUCT(", "1", ")") }},
		{"tuple", func(n int) string { return "SELECT " + nest(n, "(1, ", "1", ")") }},
		{"subscript", func(n int) string { return "SELECT " + nest(n, "a[", "0", "]") }},
		{"braced constructor", func(n int) string { return "SELECT NEW t {" + nest(n, "a {", "b: 1", "}") + "}" }},
		{"scalar subquery", func(n int) string { return nest(n, "SELECT (", "SELECT 1", ")") }},
		{"EXISTS", func(n int) string { return nest(n, "SELECT EXISTS (", "SELECT 1", ")") }},
		{"IN subquery", func(n int) string { return nest(n, "SELECT 1 IN (", "SELECT 1", ")") }},
		{"ARRAY subquery", func(n int) string { return nest(n, "SELECT ARRAY(", "SELECT 1", ")") }},
		{"invalid subquery bodies", func(n int) string { return nest(n, "SELECT (", "SELECT 1", " x y)") }},
		{"derived table", func(n int) string { return nest(n, "SELECT * FROM (", "SELECT 1", ") x") }},
		{"join group", func(n int) string { return "SELECT * FROM " + nest(n, "(t JOIN ", "t", " ON true)") }},
		{"parenthesized query", func(n int) string { return nest(n, "(", "SELECT 1", ")") }},
		{"set operation operand", func(n int) string { return nest(n, "SELECT 1 UNION ALL (", "SELECT 1", ")") }},
		{"CTE", func(n int) string { return nest(n, "WITH c AS (", "SELECT 1", ") SELECT * FROM c") }},
		{"ARRAY type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "ARRAY<", "INT64", ">") + ")" }},
		{"STRUCT type", func(n int) string { return "SELECT CAST(NULL AS " + nest(n, "STRUCT<a ", "INT64", ">") + ")" }},
		{"BEGIN block", func(n int) string { return nest(n, "BEGIN ", "SELECT 1;", " END;") }},
		{"IF", func(n int) string { return nest(n, "IF true THEN ", "SELECT 1;", " END IF;") }},
		{"LOOP", func(n int) string { return nest(n, "LOOP ", "SELECT 1;", " END LOOP;") }},
		{"nested DML", func(n int) string {
			return "UPDATE t SET " + nest(n, "(UPDATE a SET ", "x = 1", " WHERE true)") + " WHERE true"
		}},
		{"graph label negation", func(n int) string { return "GRAPH g MATCH (n IS " + strings.Repeat("!", n) + "a) RETURN n" }},
		{"graph label group", func(n int) string { return "GRAPH g MATCH (n IS " + nest(n, "(", "a", ")") + ") RETURN n" }},
		{"graph path", func(n int) string { return "GRAPH g MATCH " + nest(n, "(", "(n)", ")") + " RETURN n" }},
		{"unclosed parentheses", func(n int) string { return "SELECT " + strings.Repeat("(", n) }},
		{"unclosed subqueries", func(n int) string { return strings.Repeat("SELECT (", n) }},
		{"unclosed EXISTS", func(n int) string {
			return "SELECT 1 FROM t WHERE " + strings.Repeat("EXISTS (SELECT 1 FROM t WHERE ", n)
		}},
		{"unclosed join groups", func(n int) string { return "SELECT * FROM " + strings.Repeat("(", n) }},
		{"unclosed BEGIN blocks", func(n int) string { return strings.Repeat("BEGIN ", n) + "SELECT 1" }},
	} {
		if !strings.HasPrefix(tc.name, "unclosed") && !strings.HasPrefix(tc.name, "invalid") {
			if _, err := Parse(tc.sql(1000)); err != nil {
				t.Errorf("%s: 1,000 levels: %v", tc.name, err)
			}
		}
		for _, n := range []int{maxNestingDepth + 1, 200000} {
			finished, err := parseWithin(tc.sql(n), 10*time.Second)
			if !finished {
				t.Errorf("%s: %d levels: Parse did not return within 10s", tc.name, n)
			} else if !hasNestingError(err) {
				t.Errorf("%s: %d levels: err = %v, want the nesting error", tc.name, n, err)
			}
		}
	}
}

// TestRawSubqueryScanLoops: a subquery body that fails to parse in place is
// scanned raw to its closing ')', and the scan moved past each nested
// bracket group by calling itself, so a body that failed before its deep
// brackets overflowed the stack. It keeps a stack of groups now.
func TestRawSubqueryScanLoops(t *testing.T) {
	const n = 1000000
	for name, sql := range map[string]string{
		"closed":   "SELECT (SELECT 1 x y " + strings.Repeat("(", n) + strings.Repeat(")", n) + ")",
		"unclosed": "SELECT (SELECT 1 x y " + strings.Repeat("([{", n),
	} {
		finished, err := parseWithin(sql, 10*time.Second)
		if !finished {
			t.Errorf("%s: Parse did not return within 10s", name)
		} else if len(AllErrors(err)) == 0 || hasNestingError(err) {
			t.Errorf("%s: err = %v, want a syntax error", name, err)
		}
	}
}

// TestNestingLimitCountsReparsedBodies: a subquery body kept as text is
// parsed again after the statement, by a nested parser. That parser starts at
// the depth the body sits at, so the limit counts the levels around it. A
// bare SELECT passed as a function argument is kept as text without being
// parsed in place; were each body's count to start over, a chain of them
// would descend without bound, one reparse per level.
func TestNestingLimitCountsReparsedBodies(t *testing.T) {
	setNestingLimit(t, 100)
	sql := "SELECT " + strings.Repeat("f(SELECT ", 1000) + "1" + strings.Repeat(")", 1000)
	if _, err := Parse(sql); !hasNestingError(err) {
		t.Errorf("1,000 nested bare SELECT arguments: err = %v, want the nesting error", err)
	}
}
