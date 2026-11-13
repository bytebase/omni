package parser

import "testing"

// TestAlterTableClauseFormsOracle asserts omni agrees with the pinned TiDB on
// alterTableClauseCases (container test). The verdict comes from PREPARE,
// which parses the statement without running it: 1064 is the parser's
// rejection, and any other server error means the text parsed and failed
// later.
//
// The expectations in alterTableClauseCases are what pingcap/tidb:v8.5.5
// answered on 2026-11-04, so an image upgrade that changes a verdict fails
// here instead of passing silently.
func TestAlterTableClauseFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range alterTableClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbRelease, accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, c.accept, c.sql, perr)
			}
		})
	}
}
