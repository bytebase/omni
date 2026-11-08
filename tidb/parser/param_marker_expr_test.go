package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TestParamMarkerExprPositions checks that a `?` parameter marker parses into a
// ParamMarkerExpr, whose Loc spans the single `?` byte, in the SimpleExpr
// positions beyond LIMIT: a bare operand (SimpleExpr#8, parser.y:7886), the
// offset of LEAD/LAG (OptLeadLagInfo, parser.y:9827), and a window-frame bound
// (WindowFrameStart/Bound, parser.y:9700 and 9729). TiDB v8.5.5 accepts each.
func TestParamMarkerExprPositions(t *testing.T) {
	// firstParamMarker walks n and returns the first ParamMarkerExpr found, so a
	// case need not know the exact tree shape around the marker.
	var firstParamMarker func(n ast.Node) *ast.ParamMarkerExpr
	firstParamMarker = func(n ast.Node) *ast.ParamMarkerExpr {
		var found *ast.ParamMarkerExpr
		ast.Inspect(n, func(x ast.Node) bool {
			if found != nil {
				return false
			}
			if pm, ok := x.(*ast.ParamMarkerExpr); ok {
				found = pm
				return false
			}
			return true
		})
		return found
	}

	cases := []string{
		"DO ?",
		"DO LEAD(a, ?) OVER w",
		"DO ROW_NUMBER() OVER (ROWS ? PRECEDING)",
		"DO ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING)",
		"SELECT 1 + ?",
		"SELECT * FROM t WHERE a = ?",
	}
	for _, sql := range cases {
		list, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
			continue
		}
		pm := firstParamMarker(list.Items[0])
		if pm == nil {
			t.Errorf("%q: no ParamMarkerExpr in the tree", sql)
			continue
		}
		if pm.Loc.End != pm.Loc.Start+1 || pm.Loc.Start < 0 || pm.Loc.End > len(sql) || sql[pm.Loc.Start] != '?' {
			t.Errorf("%q: param Loc = {%d,%d} does not span a single `?`", sql, pm.Loc.Start, pm.Loc.End)
		}
	}
}
