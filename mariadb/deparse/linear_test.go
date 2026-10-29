package deparse

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ast "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// TestLongChainsDeparseStayLinear guards the single builder every node is
// written into. Building each level of a chain as a string of its own copied
// the text below it once per level, so 1 + 1 + ... took a second at 10,000
// terms and 24 seconds at 100,000. Each shape takes milliseconds when linear
// and tens of seconds otherwise.
func TestLongChainsDeparseStayLinear(t *testing.T) {
	const n = 80000
	repeat := func(first, sep, term string) string {
		var b strings.Builder
		b.WriteString(first)
		for i := 1; i < n; i++ {
			b.WriteString(sep)
			b.WriteString(term)
		}
		return b.String()
	}
	// numbered writes first, then format once for each i in [1, n).
	numbered := func(first, format string) string {
		var b strings.Builder
		b.WriteString(first)
		for i := 1; i < n; i++ {
			fmt.Fprintf(&b, format, i)
		}
		return b.String()
	}
	for _, tc := range []struct {
		name, sql, marker string
	}{
		{"sum", "SELECT " + repeat("1", " + ", "1"), " + 1)"},
		{"UNION ALL", repeat("SELECT 1", " UNION ALL ", "SELECT 1"), " union all "},
		{"AND", "SELECT * FROM t WHERE " + repeat("a = 1", " AND ", "a = 1"), " and "},
		{"OR", "SELECT * FROM t WHERE " + repeat("a = 1", " OR ", "a = 1"), " or "},
		{"join", "SELECT 1 FROM " + numbered("t AS t0", " JOIN t AS t%[1]d ON t%[1]d.a = t0.a"), " join "},
		{"comma join", "SELECT 1 FROM " + numbered("t AS t0", ", t AS t%d"), " join "},
		{"USING list", "SELECT 1 FROM t AS x JOIN t AS y USING (" + numbered("c0", ", c%d") + ")", " and "},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("%s: Parse: %v", tc.name, err)
		}
		stmt := list.Items[0].(*ast.SelectStmt)
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
