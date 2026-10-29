package catalog

import (
	"fmt"
	"strings"
	"testing"
	"time"
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
