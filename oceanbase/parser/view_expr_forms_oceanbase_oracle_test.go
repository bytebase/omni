package parser

import (
	"context"
	"testing"
)

// TestViewExprFormsOceanBaseOracle asserts omni agrees with a live OceanBase on
// viewExprFormCases. It runs on the package's own OceanBase (see
// sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestViewExprFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range viewExprFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts := prepareAcceptsOceanBase(t, ctx, db, c.sql)
			if accepts != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
