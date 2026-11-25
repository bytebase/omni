package analysis_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/cosmosdb/analysis"
	"github.com/bytebase/omni/enginetest/cputime"
)

// TestAnalyzeChainsStayLinear: field-path extraction returned a slice per
// subexpression and copied it into its parent, so every level re-copied the
// paths below it and a chain of n terms took time quadratic in n: 16,000
// ORs took seconds and 100,000 ran past two minutes. Each expression now
// appends to one slice. A value bound by FROM is shared, not flattened, so
// one that doubles the binding below it at every level stays linear too.
// The chains run to near maxTreeHeight and the nested
// shapes to near maxNestingDepth; each Analyze runs against a deadline.
func TestAnalyzeChainsStayLinear(t *testing.T) {
	const chain, nested = 90000, 9000
	where := func(expr string) string { return "SELECT * FROM c WHERE " + expr }
	for _, tc := range []struct {
		name      string
		sql       string
		wantPaths int    // distinct field paths in WHERE, or in the one projection
		wantLast  string // the last of them
	}{
		{"OR", where(strings.Repeat("c.a OR ", chain) + "c.z"), 2, "c.z"},
		{"AND", where(strings.Repeat("c.a AND ", chain) + "c.z"), 2, "c.z"},
		{"addition", where(strings.Repeat("c.a + ", chain) + "c.z"), 2, "c.z"},
		{"BETWEEN", where("c.a" + strings.Repeat(" BETWEEN c.b AND c.z", chain/2)), 3, "c.z"},
		{"LIKE", where("c.a" + strings.Repeat(" LIKE c.z", chain)), 2, "c.z"},
		{"IN", where("c.a" + strings.Repeat(" IN (c.z)", chain)), 2, "c.z"},
		{"projected OR", "SELECT " + strings.Repeat("c.a OR ", chain) + "c.z FROM c", 2, "c.z"},
		{"property of an OR", where("(" + strings.Repeat("c.a OR ", chain) + "c.z).x"), 2, "c.z.x"},
		{"properties of an object", where(`{"k": c.z}` + strings.Repeat(".k", chain)), 1, "c.z" + strings.Repeat(".k", 127)},
		{"properties of nested objects", where(strings.Repeat(`{"k": `, nested) + "c.z" + strings.Repeat("}.k", nested)), 1, "c.z"},
		{"nested scalar subqueries", "SELECT " + strings.Repeat("(SELECT VALUE ", nested/2) + "c" + strings.Repeat(".k", nested/2) + strings.Repeat(")", nested/2) + " AS v FROM c", 1, "c" + strings.Repeat(".k", 128)},
		{"a binding doubled at every level", "SELECT VALUE s FROM " + strings.Repeat(`(SELECT VALUE {"a": s, "b": s} FROM `, nested/10) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", nested/10), 1, "c.z"},
		{"a sum of a binding doubled at every level", "SELECT VALUE s FROM " + strings.Repeat("(SELECT VALUE s + s FROM ", nested/10) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", nested/10), 1, "c.z"},
		{"a computed binding doubled at every level", "SELECT VALUE s FROM " + strings.Repeat("(SELECT VALUE (s + s).k FROM ", nested/10) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", nested/10), 1, "c.z" + strings.Repeat(".k", 127)},
		{"selectors through a union shared at every level", "SELECT VALUE s" + strings.Repeat(".k", 300) + " FROM " + strings.Repeat(`(SELECT VALUE s ?? {"k": s} FROM `, 300) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", 300), 128, "c.z"},
		{"equal accesses at every level", "SELECT VALUE s FROM " + strings.Repeat("(SELECT VALUE s.k ?? s.k FROM ", 300) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", 300), 1, "c.z" + strings.Repeat(".k", 127)},
		{"equal selections made in different steps", "SELECT VALUE s FROM " + strings.Repeat("(SELECT VALUE s.k.k ?? (SELECT VALUE s.k).k FROM ", 300) + "(SELECT VALUE c.z FROM c) s" + strings.Repeat(") s", 300), 1, "c.z" + strings.Repeat(".k", 127)},
		{"SELECT * over a union shared at every level", "SELECT * FROM " + strings.Repeat(`(SELECT VALUE s ?? (s ?? {"k": 0}) FROM `, 300) + `(SELECT VALUE {"k": c.z} FROM c) s` + strings.Repeat(") s", 300), 1, "c.z"},
		{"elements of nested arrays", where(strings.Repeat("[", nested) + "c.z" + strings.Repeat("][0]", nested)), 1, "c.z"},
		{"ternary", where(strings.Repeat("c.a ? c.b : ", nested) + "c.z"), 3, "c.z"},
		{"function calls", where(strings.Repeat("CONCAT(", nested) + "c.a" + strings.Repeat(", c.z)", nested)), 2, "c.z"},
		{"arrays", where(strings.Repeat("[", nested) + "c.a" + strings.Repeat(", c.z]", nested)), 2, "c.z"},
		{"objects", where(strings.Repeat(`{"k": `, nested) + "c.a" + strings.Repeat(`, "j": c.z}`, nested)), 2, "c.z"},
	} {
		sel := parseSelect(t, tc.sql)
		var qa *analysis.QueryAnalysis
		cputime.Within(t, 2*time.Second, fmt.Sprintf("%s: Analyze", tc.name), func() {
			qa = analysis.Analyze(sel)
		})
		paths := qa.Predicates
		if len(qa.Projections) == 1 {
			paths = qa.Projections[0].SourcePaths
		}
		if len(paths) != tc.wantPaths || paths[len(paths)-1].String() != tc.wantLast {
			t.Errorf("%s: %d paths, the last %q; want %d, the last %q", tc.name, len(paths), paths[len(paths)-1], tc.wantPaths, tc.wantLast)
		}
	}
}
