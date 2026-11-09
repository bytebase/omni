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

// TestSpacedFuncViewUpdatable pins that a spaced generic call does not make a
// view non-updatable: `count (a)` is a generic scalar stored-function call, not
// the COUNT aggregate, so a view selecting a real column beside it stays
// updatable, where the adjacent aggregate makes the view non-updatable. (The
// engine would reject such a view at creation because `count` has no function,
// so this holds omni's inference consistent with its catalog analysis rather
// than to an oracle.)
func TestSpacedFuncViewUpdatable(t *testing.T) {
	c := New()
	mustExec(t, c, `CREATE DATABASE testdb; USE testdb;
		CREATE TABLE t (a INT);
		CREATE VIEW v_spaced AS SELECT a, count (a) AS x FROM t;
		CREATE VIEW v_agg AS SELECT a, count(*) AS x FROM t;`)
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("testdb not found")
	}
	if v := db.Views["v_spaced"]; v == nil || !v.IsUpdatable {
		t.Errorf("v_spaced IsUpdatable: got %v, want true (spaced generic call is not an aggregate)", v != nil && v.IsUpdatable)
	}
	if v := db.Views["v_agg"]; v == nil || v.IsUpdatable {
		t.Errorf("v_agg IsUpdatable: got %v, want false (count(*) is the aggregate)", v != nil && v.IsUpdatable)
	}
}

// TestSpacedFuncGeneratedExprRejected pins that a spaced generic call — a generic
// stored-function call — is refused in a generated-column expression and a check
// constraint, matching mysql:8.0 (3763 and 3814). Admitting it would also store
// the space-dropped, invalid builtin form (`position(a)`). A deterministic scalar
// builtin is still allowed.
func TestSpacedFuncGeneratedExprRejected(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		reject bool
	}{
		{"CREATE TABLE t (a INT, b INT AS (abs(a)))", false},
		{"CREATE TABLE t (a INT, b INT AS (position (a)))", true},
		{"CREATE TABLE t (a INT, b INT AS (sum (a)))", true},
		{"CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (position (a)))", true},
		{"CREATE TABLE t (a INT, b INT, CHECK (sum (a) > 0))", true},
		{"CREATE TABLE t (a INT, b INT CHECK (position (a) > 0))", true},
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d;")
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.sql, err)
		}
		var execErr error
		for _, r := range results {
			if r.Error != nil {
				execErr = r.Error
			}
		}
		if (execErr != nil) != tc.reject {
			t.Errorf("Exec(%q) error = %v, want reject = %v", tc.sql, execErr, tc.reject)
		}
	}
}
