package parser

import (
	"context"
	"testing"
)

// TestAsOfFormsOceanBaseOracle asserts omni agrees with a live OceanBase on
// asOfCases. startParserOracle boots MySQL in this package, so a real verdict
// needs the server itself: oceanBaseOracle runs on the package's own OceanBase (see
// sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run AsOf -count=1
func TestAsOfFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range asOfCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
