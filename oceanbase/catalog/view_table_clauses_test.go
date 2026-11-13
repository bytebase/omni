package catalog

import (
	"strings"
	"testing"
)

// TestViewDefinitionKeepsTableClauses: a view's stored definition keeps the
// clauses after its table's name. A definition without the SAMPLE or the AS OF
// SNAPSHOT would read the whole current table, which is another view.
func TestViewDefinitionKeepsTableClauses(t *testing.T) {
	c := newCatalogWithDB(t, "testdb")
	mustExec(t, c, "CREATE TABLE t (a INT, KEY i1 (a)) PARTITION BY HASH (a) PARTITIONS 2")
	for name, tc := range map[string]struct{ body, want string }{
		"v_sample":    {"SELECT a FROM t SAMPLE (1) SEED (2)", "sample(1) seed(2)"},
		"v_snapshot":  {"SELECT a FROM t AS OF SNAPSHOT 1", "as of snapshot 1"},
		"v_partition": {"SELECT a FROM t PARTITION (p0)", "partition(`p0`)"},
		"v_hint":      {"SELECT a FROM t USE INDEX (i1)", "use index(`i1`)"},
	} {
		mustExec(t, c, "CREATE VIEW "+name+" AS "+tc.body)
		if v := requireView(t, c, name); !strings.Contains(v.Definition, tc.want) {
			t.Errorf("%s: definition %q does not hold %q", name, v.Definition, tc.want)
		}
	}
}
