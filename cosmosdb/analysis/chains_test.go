package analysis_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/cosmosdb/analysis"
)

// TestAnalyzeChainsStayLinear: field-path extraction returned a slice per
// subexpression and copied it into its parent, so every level re-copied the
// paths below it and a chain of n terms took time quadratic in n: 16,000
// ORs took seconds and 100,000 ran past two minutes. Each expression now
// appends to one slice. The chains run to near maxTreeHeight and the nested
// shapes to near maxNestingDepth; each Analyze runs against a deadline.
func TestAnalyzeChainsStayLinear(t *testing.T) {
	const chain, nested = 90000, 9000
	where := func(expr string) string { return "SELECT * FROM c WHERE " + expr }
	for _, tc := range []struct {
		name      string
		sql       string
		wantPaths int    // field paths in WHERE, or in the one projection
		wantLast  string // the last of them
	}{
		{"OR", where(strings.Repeat("c.a OR ", chain) + "c.z"), chain + 1, "c.z"},
		{"AND", where(strings.Repeat("c.a AND ", chain) + "c.z"), chain + 1, "c.z"},
		{"addition", where(strings.Repeat("c.a + ", chain) + "c.z"), chain + 1, "c.z"},
		{"BETWEEN", where("c.a" + strings.Repeat(" BETWEEN c.b AND c.z", chain/2)), chain + 1, "c.z"},
		{"LIKE", where("c.a" + strings.Repeat(" LIKE c.z", chain)), chain + 1, "c.z"},
		{"IN", where("c.a" + strings.Repeat(" IN (c.z)", chain)), chain + 1, "c.z"},
		{"projected OR", "SELECT " + strings.Repeat("c.a OR ", chain) + "c.z FROM c", chain + 1, "c.z"},
		{"property of an OR", where("(" + strings.Repeat("c.a OR ", chain) + "c.z).x"), chain + 1, "c.z.x"},
		{"ternary", where(strings.Repeat("c.a ? c.b : ", nested) + "c.z"), 2*nested + 1, "c.z"},
		{"function calls", where(strings.Repeat("CONCAT(", nested) + "c.a" + strings.Repeat(", c.z)", nested)), nested + 1, "c.z"},
		{"arrays", where(strings.Repeat("[", nested) + "c.a" + strings.Repeat(", c.z]", nested)), nested + 1, "c.z"},
		{"objects", where(strings.Repeat(`{"k": `, nested) + "c.a" + strings.Repeat(`, "j": c.z}`, nested)), nested + 1, "c.z"},
	} {
		sel := parseSelect(t, tc.sql)
		done := make(chan *analysis.QueryAnalysis, 1)
		go func() { done <- analysis.Analyze(sel) }()
		select {
		case qa := <-done:
			paths := qa.Predicates
			if len(qa.Projections) == 1 {
				paths = qa.Projections[0].SourcePaths
			}
			if len(paths) != tc.wantPaths || paths[len(paths)-1].String() != tc.wantLast {
				t.Errorf("%s: %d paths, the last %q; want %d, the last %q", tc.name, len(paths), paths[len(paths)-1], tc.wantPaths, tc.wantLast)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Analyze took over 2s", tc.name)
		}
	}
}
