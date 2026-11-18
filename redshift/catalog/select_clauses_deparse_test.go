package catalog

import (
	"strings"
	"testing"
)

// TestViewGroupByAllDeparse pins Redshift's GROUP BY ALL through the catalog:
// the clause carries no list of its own, and a stored view keeps it rather than
// deparsing as a query with no grouping at all.
func TestViewGroupByAllDeparse(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_group_all AS SELECT name, sum(val) FROM t1 GROUP BY ALL")
	if !strings.Contains(def, "GROUP BY ALL") {
		t.Fatalf("view definition = %q, want GROUP BY ALL in it", def)
	}

	// A grouped query that names its columns still deparses with the list.
	plain := viewDef(t, c, "CREATE VIEW v_group_list AS SELECT name, sum(val) FROM t1 GROUP BY name")
	if !strings.Contains(plain, "GROUP BY name") || strings.Contains(plain, "ALL") {
		t.Fatalf("view definition = %q, want GROUP BY name", plain)
	}
}

// TestViewConnectByDeparse pins Redshift's START WITH / CONNECT BY through the
// catalog, PRIOR included: the conditions are analyzed, deparsed back in their
// place, and their columns count as the view's dependencies.
func TestViewConnectByDeparse(t *testing.T) {
	c := setupPhase4(t)
	def := viewDef(t, c, "CREATE VIEW v_hier AS SELECT id, name FROM t1 START WITH name IS NULL CONNECT BY PRIOR id = id")
	if !strings.Contains(def, "START WITH") {
		t.Fatalf("view definition = %q, want START WITH in it", def)
	}
	// The deparser parenthesizes the prefix operator, so the condition reads
	// (PRIOR id) = id: the word and its operand are what matter.
	if !strings.Contains(def, "CONNECT BY") || !strings.Contains(def, "PRIOR id") {
		t.Fatalf("view definition = %q, want CONNECT BY PRIOR id in it", def)
	}

	// A CONNECT BY with no START WITH deparses with the one clause it has.
	only := viewDef(t, c, "CREATE VIEW v_hier_connect AS SELECT id FROM t1 CONNECT BY id > 0")
	if strings.Contains(only, "START WITH") || !strings.Contains(only, "CONNECT BY") || !strings.Contains(only, "id > 0") {
		t.Fatalf("view definition = %q, want CONNECT BY alone", only)
	}
}
