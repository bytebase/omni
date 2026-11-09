package catalog

import (
	"strings"
	"testing"
)

// TestSpacedFuncQueryRoundTrip pins that the catalog carries the generic-call
// identity of a spaced gated keyword through analysis and deparse. Without it the
// analyzer classifies `count ()` as the COUNT aggregate and deparse renders it
// `count(*)`, and `max (id, dept)` as the builtin `max(id, dept)` — both change
// meaning and the second does not re-parse. roundTrip re-parses and re-analyzes
// the deparsed SQL, so an invalid rendering fails here.
func TestSpacedFuncQueryRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, sql2, _ := roundTrip(t, c,
		"SELECT count (salary), max (id, department_id), count () FROM employees")

	if strings.Contains(sql2, "count(*)") {
		t.Errorf("deparse rendered a spaced generic count () as the aggregate: %q", sql2)
	}
	if !strings.Contains(sql2, "count (") || !strings.Contains(sql2, "max (") {
		t.Errorf("deparse dropped the space before '(': %q", sql2)
	}

	for i, want := range []string{"count", "max", "count"} {
		fc, ok := q1.TargetList[i].Expr.(*FuncCallExprQ)
		if !ok {
			t.Fatalf("TargetList[%d].Expr: want *FuncCallExprQ, got %T", i, q1.TargetList[i].Expr)
		}
		if fc.Name != want {
			t.Errorf("TargetList[%d] name: want %q, got %q", i, want, fc.Name)
		}
		if !fc.SpacedParen {
			t.Errorf("TargetList[%d] (%s): want SpacedParen set", i, want)
		}
		if fc.IsAggregate {
			t.Errorf("TargetList[%d] (%s): a spaced generic call must not be an aggregate", i, want)
		}
	}
}
