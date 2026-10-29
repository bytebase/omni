package catalog

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mysqlparser "github.com/bytebase/omni/mariadb/parser"
)

// chainSQL writes first, then sep and term n-1 times.
func chainSQL(n int, first, sep, term string) string {
	var b strings.Builder
	b.WriteString(first)
	for i := 1; i < n; i++ {
		b.WriteString(sep)
		b.WriteString(term)
	}
	return b.String()
}

// TestLongChainViewsStayLinear creates views over long chains and deparses
// their analyzed queries against a deadline. Both deparsers built each level
// of a chain as a string of their own, copying the text below it once per
// level: CREATE VIEW over 1 + 1 + ... took a second at 10,000 terms and 24
// seconds at 100,000, and DeparseQuery grew the same way.
func TestLongChainViewsStayLinear(t *testing.T) {
	const n = 50000
	for _, tc := range []struct {
		name, sql, defMarker, queryMarker string
	}{
		{"sum", "SELECT " + chainSQL(n, "1", " + ", "1"), " + 1)", " + 1)"},
		{"UNION ALL", chainSQL(n, "SELECT a FROM t", " UNION ALL ", "SELECT a FROM t"), " union all ", " union all "},
		{"AND", "SELECT a FROM t WHERE " + chainSQL(n, "a = 1", " AND ", "a = 1"), " and ", " and "},
		{"OR", "SELECT a FROM t WHERE " + chainSQL(n, "a = 1", " OR ", "a = 1"), " or ", " or "},
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE db; USE db; CREATE TABLE t (a INT)", nil); err != nil {
			t.Fatal(err)
		}
		type result struct {
			view     *View
			deparsed string
			err      error
		}
		done := make(chan result, 1)
		go func() {
			results, err := c.Exec("CREATE VIEW v AS "+tc.sql, nil)
			if err == nil && len(results) == 1 {
				err = results[0].Error
			}
			if err != nil {
				done <- result{err: err}
				return
			}
			v := c.GetDatabase("db").Views["v"]
			if v.AnalyzedQuery == nil {
				done <- result{err: fmt.Errorf("view has no analyzed query")}
				return
			}
			done <- result{view: v, deparsed: DeparseQuery(v.AnalyzedQuery)}
		}()
		select {
		case r := <-done:
			if r.err != nil {
				t.Fatalf("%s: %v", tc.name, r.err)
			}
			if got := strings.Count(r.view.Definition, tc.defMarker); got != n-1 {
				t.Errorf("%s: view definition holds %q %d times, want %d", tc.name, tc.defMarker, got, n-1)
			}
			if got := strings.Count(r.deparsed, tc.queryMarker); got != n-1 {
				t.Errorf("%s: DeparseQuery holds %q %d times, want %d", tc.name, tc.queryMarker, got, n-1)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: CREATE VIEW and DeparseQuery over a %d-term chain took over 10s", tc.name, n)
		}
	}
}

// numberedSQL writes first, then format once for each i in [1, n).
func numberedSQL(n int, first, format string) string {
	var b strings.Builder
	b.WriteString(first)
	for i := 1; i < n; i++ {
		fmt.Fprintf(&b, format, i)
	}
	return b.String()
}

// allocatedBytes returns the number of bytes f allocates.
func allocatedBytes(f func()) uint64 {
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// TestJoinChainAnalysisStaysLinear compares what analyzing a join chain
// allocates at two lengths. Every join RTE lists the columns of both of its
// inputs, and building that list as a fresh copy at every level made the
// analysis of a chain hold a quadratic number of names: 824 MB at 10,000
// joins, and 100,000 joins were killed for memory. A select list naming every
// table of the chain (t0.*, t1.*, ...) copied the scope once per item the
// same way. Four times the joins should cost about four times the bytes, not
// sixteen.
func TestJoinChainAnalysisStaysLinear(t *testing.T) {
	const n = 1000
	for _, tc := range []struct {
		name string
		sql  func(n int) string
	}{
		{"JOIN ON", func(n int) string {
			return "SELECT t0.a FROM " + numberedSQL(n, "t AS t0", " JOIN t AS t%[1]d ON t%[1]d.a = t0.a")
		}},
		{"JOIN USING", func(n int) string {
			return "SELECT * FROM " + numberedSQL(n, "t AS t0", " JOIN t AS t%d USING (a)")
		}},
		{"NATURAL JOIN", func(n int) string {
			return "SELECT * FROM " + numberedSQL(n, "t AS t0", " NATURAL JOIN t AS t%d")
		}},
		{"right-nested JOIN", func(n int) string {
			return "SELECT * FROM " + numberedSQL(n, "", "t AS t%d JOIN (") + "t AS t0" + strings.Repeat(")", n-1)
		}},
		{"t.* per table", func(n int) string {
			return "SELECT " + numberedSQL(n, "t0.*", ", t%d.*") + " FROM " + numberedSQL(n, "t AS t0", " JOIN t AS t%d")
		}},
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE db; USE db; CREATE TABLE t (a INT, b INT)", nil); err != nil {
			t.Fatal(err)
		}
		var allocated [2]uint64
		for i, size := range []int{n, 4 * n} {
			list, err := mysqlparser.Parse(tc.sql(size))
			if err != nil {
				t.Fatalf("%s: Parse: %v", tc.name, err)
			}
			stmt := list.Items[0].(*nodes.SelectStmt)
			var q *Query
			allocated[i] = allocatedBytes(func() { q, err = c.AnalyzeSelectStmt(stmt) })
			if err != nil {
				t.Fatalf("%s: %d joins: %v", tc.name, size, err)
			}
			if got := len(q.RangeTable); got < size {
				t.Fatalf("%s: %d joins analyzed into %d RTEs", tc.name, size, got)
			}
		}
		ratio := float64(allocated[1]) / float64(allocated[0])
		t.Logf("%s: %d joins allocated %d bytes, %d joins %d bytes (%.1fx)", tc.name, n, allocated[0], 4*n, allocated[1], ratio)
		if ratio > 8 {
			t.Errorf("%s: four times the joins allocated %.1f times the bytes", tc.name, ratio)
		}
	}
}
