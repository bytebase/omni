package catalog

import (
	"strings"
	"testing"
)

// TestSpacedFuncQueryRoundTrip pins that the catalog carries the generic-call
// identity of a spaced gated keyword through analysis and deparse. TiDB's gated
// set is COUNT and SYSDATE: a space before '(' turns each into a generic
// stored-function call, not the builtin. Without this identity the analyzer
// classifies `count ()` as the COUNT aggregate and deparse renders it `count(*)`,
// and `count (salary, department_id)` as the builtin `count(salary,
// department_id)` — both change meaning and neither re-parses as the builtin.
// roundTrip re-parses and re-analyzes the deparsed SQL, so an invalid rendering
// fails here.
func TestSpacedFuncQueryRoundTrip(t *testing.T) {
	c := setupDeparseTestCatalog(t)
	q1, sql2, _ := roundTrip(t, c,
		"SELECT count (salary), count (id, department_id), count (), sysdate (id) FROM employees")

	if strings.Contains(sql2, "count(*)") {
		t.Errorf("deparse rendered a spaced generic count () as the aggregate: %q", sql2)
	}
	if !strings.Contains(sql2, "count (") || !strings.Contains(sql2, "sysdate (") {
		t.Errorf("deparse dropped the space before '(': %q", sql2)
	}

	for i, want := range []string{"count", "count", "count", "sysdate"} {
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
