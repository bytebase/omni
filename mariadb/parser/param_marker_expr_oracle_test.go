package parser

import "testing"

// TestParamMarkerExprPositionsOracle asserts omni agrees with a live MariaDB
// 11.8.8 on a `?` parameter marker in the SimpleExpr positions beyond LIMIT.
// MariaDB accepts `?` as a bare operand (param_marker in simple_expr) and as
// the LEAD/LAG offset, but — unlike MySQL — rejects it as a window-frame bound,
// whose grammar is a numeric literal or INTERVAL. MariaDB accepts `?` only
// under PREPARE, so the oracle is read through prepareParses. The analyzer
// boundary's rejection of a standalone marker is pinned by the catalog tests.
func TestParamMarkerExprPositionsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT ?", true},
		{"SELECT 1 + ?", true},
		{"SELECT LEAD(a, ?) OVER () FROM t", true},
		// MariaDB rejects a `?` frame bound (MySQL/TiDB/OceanBase accept it).
		{"SELECT ROW_NUMBER() OVER (ROWS ? PRECEDING) FROM t", false},
		{"SELECT ROW_NUMBER() OVER (ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM t", false},
	} {
		if got := o.prepareParses(t, c.sql); got != c.accept {
			t.Fatalf("MariaDB prepareParses(%q)=%v, want %v", c.sql, got, c.accept)
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != c.accept {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
				omniAccepts, c.accept, c.sql, perr)
		}
	}
}
