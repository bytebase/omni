package parser

import (
	"context"
	"testing"
)

// TestTableSampleFormsOceanBaseOracle asserts omni agrees with a live OceanBase on
// tableSampleCases. startParserOracle boots MySQL in this package, so a real verdict
// needs the server itself, which CI has no room for. This reads OCEANBASE_WALK_DSN and
// skips without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run TableSample -count=1
func TestTableSampleFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range tableSampleCases {
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
