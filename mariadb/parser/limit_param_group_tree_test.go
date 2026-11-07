package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestLimitParamMarkerTree checks that a `?` LIMIT operand parses into a
// ParamMarkerExpr whose Loc spans the single `?`, in each LIMIT position.
// Acceptance is container-verified by TestLimitParamMarkerOracle.
func TestLimitParamMarkerTree(t *testing.T) {
	spansOneQuestionMark := func(sql string, e ast.ExprNode) {
		t.Helper()
		pm, ok := e.(*ast.ParamMarkerExpr)
		if !ok {
			t.Errorf("%q: operand = %T, want *ParamMarkerExpr", sql, e)
			return
		}
		if pm.Loc.End != pm.Loc.Start+1 || pm.Loc.Start < 0 || pm.Loc.End > len(sql) || sql[pm.Loc.Start] != '?' {
			t.Errorf("%q: param Loc = {%d,%d} does not span a single `?`", sql, pm.Loc.Start, pm.Loc.End)
		}
	}

	cases := []struct {
		sql        string
		wantOffset bool
	}{
		{"SELECT * FROM t LIMIT ?", false},
		{"SELECT * FROM t LIMIT ?, ?", true},
		{"SELECT * FROM t LIMIT ? OFFSET ?", true},
	}
	for _, c := range cases {
		list, err := Parse(c.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.sql, err)
			continue
		}
		sel := list.Items[0].(*ast.SelectStmt)
		if sel.Limit == nil || sel.Limit.Count == nil {
			t.Errorf("%q: no LIMIT count", c.sql)
			continue
		}
		spansOneQuestionMark(c.sql, sel.Limit.Count)
		if c.wantOffset {
			spansOneQuestionMark(c.sql, sel.Limit.Offset)
		}
	}

	// A `?` is a LIMIT operand only, not a general expression.
	if _, err := Parse("SELECT ?"); err == nil {
		t.Errorf("SELECT ? should still be rejected")
	}
}

// TestGroupByDirectionTree checks that GROUP BY items carry an optional
// ASC/DESC direction (MariaDB shares ORDER BY's ByItem rule), with ASC the
// default. Acceptance is container-verified by TestGroupByDirectionOracle.
func TestGroupByDirectionTree(t *testing.T) {
	cases := []struct {
		sql      string
		wantDesc []bool
	}{
		{"SELECT a FROM t GROUP BY a", []bool{false}},
		{"SELECT a FROM t GROUP BY a ASC", []bool{false}},
		{"SELECT a FROM t GROUP BY a DESC", []bool{true}},
		{"SELECT a, b FROM t GROUP BY a ASC, b DESC", []bool{false, true}},
	}
	for _, c := range cases {
		list, err := Parse(c.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.sql, err)
			continue
		}
		sel := list.Items[0].(*ast.SelectStmt)
		if len(sel.GroupBy) != len(c.wantDesc) {
			t.Errorf("%q: %d GROUP BY items, want %d", c.sql, len(sel.GroupBy), len(c.wantDesc))
			continue
		}
		for i, want := range c.wantDesc {
			if sel.GroupBy[i].Desc != want {
				t.Errorf("%q: item %d Desc = %v, want %v", c.sql, i, sel.GroupBy[i].Desc, want)
			}
		}
	}

	// ASC/DESC is at most one per item.
	if _, err := Parse("SELECT a FROM t GROUP BY a ASC DESC"); err == nil {
		t.Errorf("GROUP BY a ASC DESC should be rejected")
	}
}
