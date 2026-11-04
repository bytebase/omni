package parser

import (
	"context"
	"testing"
)

// TestGrantFormsOnOceanBase holds grant_forms_test.go's verdicts to OceanBase
// itself. oracle_test.go's shared oracle is MySQL 8.0, which answers for the
// grammar this engine inherits but not for this family: the two servers
// disagree on the auth clauses, on the REQUIRE and resource clauses of a
// GRANT, and on what a host may hold.
//
// OceanBase needs about 8 GB, which CI does not have, so this reads
// OCEANBASE_WALK_DSN and skips without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestGrantFormsOnOceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range append(append([]grantFormCase{}, grantFormCases...), grantKnownDivergences...) {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := engineAccepts(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			want := c.accept
			for _, d := range grantKnownDivergences {
				if d.sql == c.sql {
					// A recorded divergence is a statement the engine
					// and this parser answer differently, so the
					// engine's answer is the other one.
					want = !c.accept
					break
				}
			}
			if accepts != want {
				t.Errorf("OceanBase %ss %q, this package records %s",
					grantVerdict(accepts), c.sql, grantVerdict(want))
			}
		})
	}
}
