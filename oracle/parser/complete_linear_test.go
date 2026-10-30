package parser

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestCompletionScopeStaysLinear: completion read the references of each
// SELECT by scanning the statement from its start, found each matching
// parenthesis by a scan, and expanded a star over a derived table again for
// every query around it. Completion through 1,000 nested derived tables
// took 20s to 35s, and the CTE chain, set operations, scalar subqueries,
// and derived-table list below took 9s to 21s; each takes milliseconds
// now.
func TestCompletionScopeStaysLinear(t *testing.T) {
	derived := strings.Repeat("SELECT * FROM (", 1000) + "SELECT a FROM t" + strings.Repeat(") x", 1000)
	var cte strings.Builder
	cte.WriteString("WITH c0 AS (SELECT a FROM t)")
	for i := 1; i <= 10000; i++ {
		cte.WriteString(", c" + strconv.Itoa(i) + " AS (SELECT * FROM c" + strconv.Itoa(i-1) + ")")
	}
	cte.WriteString(" SELECT * FROM c10000 WHERE ")
	unionGroups := "SELECT * FROM " + strings.Repeat("(SELECT a FROM t UNION SELECT b FROM u) x, ", 10000) + "(SELECT  FROM v) y"
	for _, tc := range []struct {
		name, sql string
		caret     int
	}{
		{"derived tables, caret after them", derived, len(derived)},
		{"derived tables, caret in the outer select list", derived, len("SELECT ")},
		{"derived tables, caret in the inner select list", derived, strings.Index(derived, "SELECT a") + len("SELECT ")},
		{"CTE chain", cte.String(), cte.Len()},
		{"set operations", "SELECT a FROM t" + strings.Repeat(" UNION SELECT a FROM t", 20000) + " WHERE ", -1},
		{"scalar subqueries", strings.Repeat("SELECT (", 10000) + "SELECT a FROM t" + strings.Repeat(") FROM t", 10000), -1},
		{"derived-table list", "SELECT * FROM " + strings.Repeat("(SELECT a FROM t) x, ", 10000) + "t WHERE ", -1},
		{"derived tables with set operations", unionGroups, strings.LastIndex(unionGroups, "SELECT ") + len("SELECT ")},
		{"IN subqueries with set operations", "SELECT * FROM t WHERE " + strings.Repeat("a IN (SELECT b FROM u UNION SELECT c FROM w) AND ", 10000), -1},
	} {
		caret := tc.caret
		if caret < 0 {
			caret = len(tc.sql)
		}
		done := make(chan *CompletionContext, 1)
		go func() {
			done <- CollectCompletion(tc.sql, caret)
		}()
		select {
		case got := <-done:
			if len(got.Scope.LocalReferences) == 0 {
				t.Errorf("%s: no local references", tc.name)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: completion in %d bytes took over 2s", tc.name, len(tc.sql))
		}
	}
	// The outer select list sees the outer derived table, whose star
	// expands through every level to the innermost select list.
	got := CollectCompletion(derived, len("SELECT "))
	if refs := got.Scope.LocalReferences; len(refs) != 1 || strings.Join(refs[0].Columns, ",") != "A" {
		t.Errorf("columns through 1,000 derived tables = %+v, want one reference with [A]", refs)
	}
}
