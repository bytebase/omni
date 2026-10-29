package catalog

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestViewChainsStayLinear: CREATE VIEW over a long left-deep chain, and the
// definition deparsed from the view, take time and memory linear in the
// chain. The analyzer nests each arm of a UNION one query level inside the
// next, and each shape here was quadratic or worse before:
//
//   - UNION ALL: recording the view's dependencies copied the stack of query
//     levels at every arm (30,000 arms took 21s and allocated 9 GB), and
//     deparse copied the whole definition to trim spaces at every line.
//   - UNION ALL of tables: looking each table name up as a CTE searched the
//     level of every arm enclosing it.
//   - UNION ALL of outer columns: resolving the column did the same, with a
//     Var allocated per level, and deparse stepped up the levels one by one.
//   - concatenation: the collation check rescanned the chain below every
//     level for a COLLATE.
//   - joins: every column an ON condition names rebuilt the set of tables
//     the condition sees, listing them in time quadratic in the joins. Each
//     join lists every column of its inputs, as in PostgreSQL, so memory
//     grows quadratically there and only the time is checked.
//
// Linear growth allocates about 4x for 4x the terms; quadratic, 16x.
func TestViewChainsStayLinear(t *testing.T) {
	// repeat returns first followed by n-1 copies of each, in which %[1]d,
	// if any, is the copy's number.
	repeat := func(n int, first, each string) string {
		if !strings.Contains(each, "%") {
			return first + strings.Repeat(each, n-1)
		}
		var b strings.Builder
		b.WriteString(first)
		for i := 1; i < n; i++ {
			fmt.Fprintf(&b, each, i)
		}
		return b.String()
	}
	for _, tc := range []struct {
		name      string
		body      func(n int) string
		n         int
		budget    time.Duration
		quadratic bool // memory grows quadratically by nature
	}{
		{"UNION ALL", func(n int) string { return repeat(n, "SELECT 1", " UNION ALL SELECT 1") }, 20000, 10 * time.Second, false},
		{"UNION ALL of tables", func(n int) string { return repeat(n, "SELECT a FROM t", " UNION ALL SELECT a FROM t") }, 20000, 10 * time.Second, false},
		{"UNION ALL of outer columns", func(n int) string {
			return repeat(n, "SELECT (SELECT t.a", " UNION ALL SELECT t.a") + ") FROM t"
		}, 20000, 10 * time.Second, false},
		{"concatenation", func(n int) string { return repeat(n, "SELECT 'a'", " || 'a'") }, 20000, 10 * time.Second, false},
		{"joins", func(n int) string { return repeat(n, "SELECT 1 FROM t t0", " JOIN t t%[1]d ON t%[1]d.a = t0.a") }, 3000, 15 * time.Second, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			small := defineViewWithin(t, tc.body(tc.n), tc.budget, true)
			if tc.quadratic {
				return
			}
			large := defineViewWithin(t, tc.body(4*tc.n), tc.budget, true)
			if large.total() > 8*small.total() {
				t.Errorf("%d terms allocated %d bytes, %d terms %d bytes", tc.n, small.total(), 4*tc.n, large.total())
			}
		})
	}
}

// TestNestedSubqueryViewStaysLinear: the collation check of a view column
// follows the column down through the subqueries it comes from. It copied
// the stack of query levels at each one, so a column of 2,000 nested
// subqueries copied two million pointers. The definition is left out: its
// indentation grows with the nesting, as in PostgreSQL.
func TestNestedSubqueryViewStaysLinear(t *testing.T) {
	body := func(n int) string {
		return strings.Repeat("SELECT * FROM (", n) + "SELECT a, a AS b, a AS c, a AS d, a AS e, a AS f, a AS g, a AS h FROM t" + strings.Repeat(") s", n)
	}
	const n = 500
	small := defineViewWithin(t, body(n), 10*time.Second, false)
	large := defineViewWithin(t, body(4*n), 10*time.Second, false)
	if large.exec > 8*small.exec {
		t.Errorf("%d nested subqueries allocated %d bytes, %d allocated %d bytes", n, small.exec, 4*n, large.exec)
	}
}

// viewCost is what defineViewWithin allocated: to create the view, and to
// deparse its definition.
type viewCost struct{ exec, deparse uint64 }

func (c viewCost) total() uint64 { return c.exec + c.deparse }

// defineViewWithin creates view v AS body over table t(a int) in a fresh
// catalog and, when deparse is set, deparses its definition, failing the
// test when that takes longer than budget.
func defineViewWithin(t *testing.T, body string, budget time.Duration, deparse bool) viewCost {
	t.Helper()
	c := New()
	if _, err := c.Exec("CREATE TABLE t (a int)", nil); err != nil {
		t.Fatal(err)
	}
	type outcome struct {
		cost viewCost
		err  error
	}
	done := make(chan outcome, 1)
	go func() {
		var m0, m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m0)
		results, err := c.Exec("CREATE VIEW v AS "+body, nil)
		for _, r := range results {
			if err == nil {
				err = r.Error
			}
		}
		runtime.ReadMemStats(&m1)
		if err == nil && deparse {
			_, err = c.GetViewDefinition("public", "v")
		}
		runtime.ReadMemStats(&m2)
		done <- outcome{viewCost{exec: m1.TotalAlloc - m0.TotalAlloc, deparse: m2.TotalAlloc - m1.TotalAlloc}, err}
	}()
	select {
	case o := <-done:
		if o.err != nil {
			t.Fatal(o.err)
		}
		return o.cost
	case <-time.After(budget):
		if deparse {
			t.Fatalf("CREATE VIEW and its definition took over %v", budget)
		}
		t.Fatalf("CREATE VIEW took over %v", budget)
	}
	return viewCost{}
}
