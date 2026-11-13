package deparse

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

// TestDeparseKeepsTableClauses: the clauses after a table's name survive a
// deparse, each alone and all together, and what comes out parses and deparses
// to itself. A deparse that dropped a SAMPLE or an AS OF SNAPSHOT
// would read the whole current table.
func TestDeparseKeepsTableClauses(t *testing.T) {
	for _, c := range []struct{ input, expected string }{
		{"SELECT a FROM t SAMPLE (1)", "select `a` AS `a` from `t` sample(1)"},
		{"SELECT a FROM t SAMPLE BLOCK ALL (1.5) SEED (7)", "select `a` AS `a` from `t` sample block all(1.5) seed(7)"},
		{"SELECT a FROM t SAMPLE HYBRID INCR (10) x", "select `a` AS `a` from `t` sample hybrid incr(10) `x`"},
		{"SELECT a FROM t AS OF SNAPSHOT 1", "select `a` AS `a` from `t` as of snapshot 1"},
		{"SELECT a FROM t AS OF SNAPSHOT 1 + 2 x", "select `a` AS `a` from `t` as of snapshot (1 + 2) `x`"},
		{"SELECT a FROM t AS OF PROCTIME()", "select `a` AS `a` from `t` as of proctime()"},
		{"SELECT a FROM t PARTITION (p0, p1)", "select `a` AS `a` from `t` partition(`p0`,`p1`)"},
		{"SELECT a FROM t PARTITION (`d` = '2026', h = 1)", "select `a` AS `a` from `t` partition(`d` = '2026',`h` = 1)"},
		{"SELECT a FROM t USE INDEX (i1, PRIMARY) IGNORE KEY FOR ORDER BY (i2)", "select `a` AS `a` from `t` use index(`i1`,primary) ignore index for order by(`i2`)"},
		{"SELECT a FROM t FORCE INDEX FOR JOIN (i1) JOIN u USE KEY () ON t.a = u.a", "select `a` AS `a` from (`t` force index for join(`i1`) join `u` use index() on((`t`.`a` = `u`.`a`)))"},
		{
			"SELECT a FROM d.t PARTITION (p0) SAMPLE (1) SEED (2) AS OF SNAPSHOT 3 AS x USE INDEX (i1)",
			"select `a` AS `a` from `d`.`t` partition(`p0`) sample(1) seed(2) as of snapshot 3 `x` use index(`i1`)",
		},
	} {
		list, err := parser.Parse(c.input)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.input, err)
			continue
		}
		got := DeparseSelect(list.Items[0].(*ast.SelectStmt))
		if got != c.expected {
			t.Errorf("DeparseSelect(%q)\n got %s\nwant %s", c.input, got, c.expected)
		}
		// What comes out parses, and deparses to itself: no clause is lost on the
		// way back in.
		again, err := parser.Parse(got)
		if err != nil {
			t.Errorf("Parse(%q), the deparse of %q: %v", got, c.input, err)
			continue
		}
		if twice := DeparseSelect(again.Items[0].(*ast.SelectStmt)); twice != got {
			t.Errorf("deparse of %q is not stable:\n once  %s\n twice %s", c.input, got, twice)
		}
	}
}
