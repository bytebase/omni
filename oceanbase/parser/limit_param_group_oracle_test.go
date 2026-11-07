package parser

import (
	"context"
	"testing"
)

// TestLimitParamMarkerOracle asserts omni agrees with a live OceanBase on a `?`
// parameter marker as a LIMIT/OFFSET operand. OceanBase accepts `?` only in a
// prepared statement, so the oracle is read through PREPARE
// (prepareAcceptsOceanBase). PREPARE accepts `?` in any position, so it cannot
// measure omni's deliberate restriction of `?` to the LIMIT operand; that is
// pinned by TestLimitParamMarkerTree.
func TestLimitParamMarkerOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT * FROM t LIMIT ?", true},
		{"SELECT * FROM t LIMIT ?, ?", true},
		{"SELECT * FROM t LIMIT ? OFFSET ?", true},
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

// TestGroupByDirectionOracle asserts omni agrees with a live OceanBase on an
// ASC/DESC direction on a GROUP BY item. OceanBase (MySQL mode) keeps the
// direction that MySQL 8.0 removed; at most one direction per item.
func TestGroupByDirectionOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT a FROM t GROUP BY a ASC", true},
		{"SELECT a FROM t GROUP BY a DESC", true},
		{"SELECT a, b FROM t GROUP BY a ASC, b DESC", true},
		{"SELECT a FROM t GROUP BY a ASC DESC", false},
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
