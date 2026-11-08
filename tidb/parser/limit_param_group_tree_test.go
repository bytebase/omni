package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TestLimitParamMarkerTree checks that a `?` LIMIT operand parses into a
// ParamMarkerExpr whose Loc spans the single `?`, in each LIMIT position.
func TestLimitParamMarkerTree(t *testing.T) {
	// spansOneQuestionMark asserts the node's Loc covers exactly the single `?`
	// byte in sql.
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
		wantOffset bool // Offset is also a ParamMarkerExpr
	}{
		{"SELECT * FROM t LIMIT ?", false},
		{"SELECT * FROM t LIMIT ?, ?", true},
		{"SELECT * FROM t LIMIT ? OFFSET ?", true},
		{"SELECT * FROM t FETCH FIRST ? ROW ONLY", false},
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

	// A `?` is a SimpleExpr alternative (parser.y:7886), so it also parses as a
	// bare select-list operand into a ParamMarkerExpr. The analyzer boundary,
	// not the parser, refuses an unbound marker.
	list, err := Parse("SELECT ?")
	if err != nil {
		t.Fatalf(`Parse("SELECT ?"): %v`, err)
	}
	targets := list.Items[0].(*ast.SelectStmt).TargetList
	if len(targets) != 1 {
		t.Fatalf(`"SELECT ?": %d targets, want 1`, len(targets))
	}
	spansOneQuestionMark("SELECT ?", targets[0])
}

// TestGroupByDirectionTree checks that GROUP BY items carry an optional
// ASC/DESC direction (TiDB's shared ByItem rule), with ASC the default.
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
