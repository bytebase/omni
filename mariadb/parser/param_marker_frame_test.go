package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestParamMarkerExprPositions pins omni's parse verdicts for a `?` parameter
// marker in the SimpleExpr positions beyond LIMIT, which diverge by engine.
// MariaDB accepts `?` as a bare operand and as the LEAD/LAG offset, but — unlike
// MySQL, TiDB and OceanBase — rejects it as a window-frame bound (its bound is a
// numeric literal or INTERVAL). The live agreement is in
// TestParamMarkerExprPositionsOracle; this test needs no container.
func TestParamMarkerExprPositions(t *testing.T) {
	accept := []string{
		"SELECT ?",
		"SELECT 1 + ?",
		"SELECT a FROM t WHERE a = ?",
		"SELECT LEAD(a, ?) OVER () FROM t",
		"DO LEAD(a, ?) OVER w",
	}
	for _, sql := range accept {
		list, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q): %v, want accept", sql, err)
			continue
		}
		var found bool
		ast.Inspect(list.Items[0], func(n ast.Node) bool {
			if _, ok := n.(*ast.ParamMarkerExpr); ok {
				found = true
			}
			return true
		})
		if !found {
			t.Errorf("%q: no ParamMarkerExpr in the tree", sql)
		}
	}

	// A bare `?` is not a valid MariaDB window-frame bound (its grammar has no
	// paramMarker alternative there), in either bound of a BETWEEN.
	reject := []string{
		"SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t",
		"SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t",
		"SELECT SUM(a) OVER (ORDER BY a ROWS BETWEEN UNBOUNDED PRECEDING AND ? FOLLOWING) FROM t",
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want a rejection of the ? frame bound", sql)
		}
	}
}
