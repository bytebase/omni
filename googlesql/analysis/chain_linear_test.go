package analysis

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestAnalysisChainsStayLinear guards the left-deep chains the parser builds
// with loops: JOIN chains, a parenthesized comma join, set-operation chains, a
// field-access chain, a sum. Classify and GetQuerySpan redid the work of every
// level below each level: they copied the accumulated join projection at
// every join, flattened the rest of a field chain at every field, and
// re-unioned the sources gathered so far at every set-operation arm and every
// column reference. Classify also resolved the lineage it never reads, which
// for a USING chain whose keys change is cubic. At these sizes each shape took
// from several seconds to half a minute; it takes milliseconds now.
func TestAnalysisChainsStayLinear(t *testing.T) {
	const joins, fields, arms = 12000, 20000, 6000
	// numbered repeats format with i = 1..n, every verb given i.
	numbered := func(n int, format string) string {
		var b strings.Builder
		verbs := strings.Count(format, "%d")
		for i := 1; i <= n; i++ {
			args := make([]any, verbs)
			for j := range args {
				args[j] = i
			}
			fmt.Fprintf(&b, format, args...)
		}
		return b.String()
	}
	for _, tc := range []struct {
		name string
		sql  string
		// classifyOnly skips GetQuerySpan: its result for this shape is
		// quadratic in size (every star's except list grows with the keys after
		// it), so only Classify, which reads no lineage, is linear.
		classifyOnly bool
	}{
		{name: "JOIN ON", sql: "SELECT * FROM t" + strings.Repeat(" JOIN t ON true", joins)},
		{name: "CROSS JOIN", sql: "SELECT * FROM t" + strings.Repeat(" CROSS JOIN t", joins)},
		{name: "JOIN USING one key", sql: "SELECT * FROM t" + strings.Repeat(" JOIN t USING (a)", joins)},
		{name: "parenthesized comma join", sql: "SELECT * FROM (t" + strings.Repeat(", t", joins) + ")"},
		// Nesting, which the parser bounds, limits this one to fewer joins.
		{name: "right-nested JOIN", sql: "SELECT * FROM " + strings.Repeat("t JOIN (", 8000) + "t" + strings.Repeat(") ON true", 8000)},
		{name: "JOIN USING changing keys", sql: "SELECT * FROM t0" + numbered(1000, " JOIN t%d USING (a%d)"), classifyOnly: true},
		{name: "field access on a call", sql: "SELECT f(a)" + strings.Repeat(".b", fields)},
		{name: "field access on a subscript", sql: "SELECT 1 FROM t WHERE a[0]" + strings.Repeat(".b", fields) + " = 1"},
		{name: "UNION ALL", sql: "SELECT a FROM t0" + numbered(arms, " UNION ALL SELECT a FROM t%d")},
		{name: "UNION ALL BY NAME", sql: "SELECT a FROM t0" + numbered(arms, " UNION ALL BY NAME SELECT a FROM t%d")},
		{name: "sum", sql: "SELECT a0" + numbered(arms, " + a%d") + " FROM t"},
	} {
		done := make(chan error, 1)
		go func() {
			if got := ClassifySQL(tc.sql, DialectBigQuery); got != Select {
				done <- fmt.Errorf("ClassifySQL = %v, want %v", got, Select)
				return
			}
			if tc.classifyOnly {
				done <- nil
				return
			}
			_, err := GetQuerySpan(tc.sql, DialectBigQuery)
			done <- err
		}()
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: ClassifySQL and GetQuerySpan took over 2s", tc.name)
		}
	}
}
