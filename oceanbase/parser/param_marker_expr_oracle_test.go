package parser

import (
	"context"
	"testing"
)

// TestParamMarkerExprPositionsOracle asserts omni agrees with a live OceanBase
// (MySQL mode) on a `?` parameter marker in the SimpleExpr positions beyond
// LIMIT: a bare operand (param_marker in simple_expr), the offset of LEAD/LAG,
// and a window-frame bound. OceanBase accepts `?` only under PREPARE, so the
// oracle is read through prepareAcceptsOceanBase. The analyzer boundary's
// rejection of a standalone marker is pinned by the catalog tests.
func TestParamMarkerExprPositionsOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT ?", true},
		{"SELECT 1 + ?", true},
		{"SELECT LEAD(a, ?) OVER () FROM t", true},
		{"SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t", true},
		{"SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t", true},
	} {
		if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
			t.Fatalf("OceanBase prepareAccepts(%q)=%v, want %v", c.sql, got, c.accept)
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != c.accept {
			t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
				omniAccepts, c.accept, c.sql, perr)
		}
	}
}
