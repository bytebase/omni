package parser

import (
	"strings"
	"testing"
	"time"
)

// collectWithin runs Collect at cursor and fails the test when it takes over
// ten seconds.
func collectWithin(t *testing.T, name, sql string, cursor int) *CandidateSet {
	t.Helper()
	done := make(chan *CandidateSet, 1)
	go func() { done <- Collect(sql, cursor) }()
	select {
	case cs := <-done:
		return cs
	case <-time.After(10 * time.Second):
		t.Fatalf("%s: Collect at offset %d of %d bytes took over 10s", name, cursor, len(sql))
		return nil
	}
}

// TestNestedParensCollectStayLinear is TestNestedParensStayLinear for
// completion. A completion-mode parse of a parenthesized query left
// candidates behind, so its failures were not recorded and the speculation at
// every level of `((((SELECT 1) = 0) = 0) = 0)` reparsed every level below:
// 1,000 levels took 75 ms and 10,000 took 8 s. Each shape takes milliseconds
// when linear and minutes otherwise.
func TestNestedParensCollectStayLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	for _, tc := range []struct{ name, sql string }{
		{"scalar", "SELECT " + nest(depth, "(", "1", ")")},
		{"subquery", "SELECT " + nest(depth, "(", "SELECT 1", ")")},
		{"subquery comparison", "SELECT " + nest(depth, "(", "SELECT 1", ") = 0")},
		{"aliased subquery comparison", "SELECT " + nest(depth, "(", "SELECT a AS x FROM t", ") = 0") + " AS y FROM t ORDER BY "},
		{"join group", "SELECT * FROM " + nest(depth, "(", "t", ")")},
	} {
		inner := strings.Index(tc.sql, "1")
		if inner < 0 {
			inner = strings.LastIndex(tc.sql, "(") + 1
		}
		for _, cursor := range []int{len(tc.sql), inner} {
			collectWithin(t, tc.name, tc.sql, cursor)
		}
	}
	// The candidates at the cursor survive: inside the innermost query and
	// after the whole statement.
	sql := "SELECT " + nest(depth, "(", "SELECT a FROM t WHERE ", ") = 0")
	if cs := collectWithin(t, "inner WHERE", sql, strings.Index(sql, "WHERE ")+len("WHERE ")); !cs.HasRule("columnref") {
		t.Errorf("inner WHERE: no columnref candidate")
	}
	sql = "SELECT " + nest(depth, "(", "SELECT 1", ") = 0") + " "
	if cs := collectWithin(t, "after", sql, len(sql)); !cs.HasToken(kwFROM) {
		t.Errorf("after the statement: no FROM candidate")
	}
}

// TestCollectOverFailingTokensStaysLinear guards lineCol: completion resumes
// after every statement that fails, and each syntax error counted the
// newlines from the start of the input, so a run of n failing tokens took
// time quadratic in n.
func TestCollectOverFailingTokensStaysLinear(t *testing.T) {
	const n = 200000
	for _, sql := range []string{
		"SELECT 1 " + strings.Repeat(") ", n),
		strings.Repeat("=\n", n),
	} {
		collectWithin(t, "failing tokens", sql, len(sql))
	}
}
