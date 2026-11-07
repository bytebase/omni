package parser

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// TestLimitParamMarkerTree checks that a `?` LIMIT operand parses into a
// ParamMarkerExpr whose Loc spans the single `?`, in each LIMIT position.
// MySQL accepts `?` as a LIMIT operand only in a prepared statement; the
// acceptance is container-verified by TestLimitParamMarkerOracle.
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
