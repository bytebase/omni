package deparse

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ast "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// repeatTerms writes first, then sep and term n-1 times.
func repeatTerms(n int, first, sep, term string) string {
	var b strings.Builder
	b.WriteString(first)
	for i := 1; i < n; i++ {
		b.WriteString(sep)
		b.WriteString(term)
	}
	return b.String()
}

// numberedTerms writes first, then format once for each i in [1, n).
func numberedTerms(n int, first, format string) string {
	var b strings.Builder
	b.WriteString(first)
	for i := 1; i < n; i++ {
		fmt.Fprintf(&b, format, i)
	}
	return b.String()
}

// parseSelectOrFatal parses sql and returns its one SELECT.
func parseSelectOrFatal(t *testing.T, name, sql string) *ast.SelectStmt {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("%s: Parse: %v", name, err)
	}
	return list.Items[0].(*ast.SelectStmt)
}

// TestLongChainsDeparseStayLinear guards the single builder every node is
// written into. Building each level of a chain as a string of its own copied
// the text below it once per level, so 1 + 1 + ... took a second at 10,000
// terms and 24 seconds at 100,000. Each shape takes milliseconds when linear
// and tens of seconds otherwise.
func TestLongChainsDeparseStayLinear(t *testing.T) {
	const n = 80000
	for _, tc := range []struct {
		name, sql, marker string
	}{
		{"sum", "SELECT " + repeatTerms(n, "1", " + ", "1"), " + 1)"},
		{"UNION ALL", repeatTerms(n, "SELECT 1", " UNION ALL ", "SELECT 1"), " union all "},
		{"AND", "SELECT * FROM t WHERE " + repeatTerms(n, "a = 1", " AND ", "a = 1"), " and "},
		{"OR", "SELECT * FROM t WHERE " + repeatTerms(n, "a = 1", " OR ", "a = 1"), " or "},
		{"join", "SELECT 1 FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%[1]d ON t%[1]d.a = t0.a"), " join "},
		{"comma join", "SELECT 1 FROM " + numberedTerms(n, "t AS t0", ", t AS t%d"), " join "},
		{"USING list", "SELECT 1 FROM t AS x JOIN t AS y USING (" + numberedTerms(n, "c0", ", c%d") + ")", " and "},
	} {
		stmt := parseSelectOrFatal(t, tc.name, tc.sql)
		done := make(chan string, 1)
		go func() { done <- DeparseSelect(stmt) }()
		select {
		case out := <-done:
			if got := strings.Count(out, tc.marker); got != n-1 {
				t.Errorf("%s: deparsed text holds %q %d times, want %d", tc.name, tc.marker, got, n-1)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: DeparseSelect of a %d-term chain took over 10s", tc.name, n)
		}
	}
}

// TestResolveLongJoinChainStaysLinear guards the scope lookups behind column
// qualification: finding the entry for each column reference or t.* by
// scanning every table of the FROM clause made a view over a join chain take
// time quadratic in its length. Each shape takes milliseconds when linear and
// tens of seconds otherwise.
func TestResolveLongJoinChainStaysLinear(t *testing.T) {
	const n = 50000
	r := &Resolver{Lookup: func(name string) *ResolverTable {
		return &ResolverTable{Name: name, Columns: []ResolverColumn{{Name: "a", Position: 1}, {Name: "b", Position: 2}}}
	}}
	for _, tc := range []struct {
		name, sql, marker string
		want              int
	}{
		{"qualified columns", "SELECT t0.a FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%[1]d ON t%[1]d.a = T%[1]d.b"), "`.`b`", n - 1},
		{"unqualified columns", "SELECT a FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%d ON a = b"), "(`t0`.`a` = `t0`.`b`)", n - 1},
		{"qualified stars", "SELECT " + numberedTerms(n, "t0.*", ", T%d.*") + " FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%d"), "` AS `b`", n},
		// A NATURAL or USING join finds the columns it joins on in its
		// inputs; SELECT * lists each such column once.
		{"USING chain", "SELECT * FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%d USING (a)"), "` AS `b`", n},
		{"USING chain on two columns", "SELECT * FROM " + numberedTerms(n, "t AS t0", " JOIN t AS t%d USING (b, a)"), "` AS `b`", 1},
		{"NATURAL chain", "SELECT * FROM " + numberedTerms(n, "t AS t0", " NATURAL JOIN t AS t%d"), "(`t0`.`b` = `t", n - 1},
		{"right-nested USING", "SELECT * FROM " + numberedTerms(9000, "", "t AS t%d JOIN (") + "t AS t0" + strings.Repeat(") USING (a)", 8999), "` AS `b`", 9000},
	} {
		stmt := parseSelectOrFatal(t, tc.name, tc.sql)
		done := make(chan string, 1)
		go func() { done <- DeparseSelect(r.Resolve(stmt)) }()
		select {
		case out := <-done:
			if got := strings.Count(out, tc.marker); got != tc.want {
				t.Errorf("%s: resolved text holds %q %d times, want %d", tc.name, tc.marker, got, tc.want)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: Resolve over a %d-table join chain took over 10s", tc.name, n)
		}
	}
}
