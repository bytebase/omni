package parser

import "testing"

// TestPriorityFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// priorityCases (container test). The verdict comes from PREPARE, which parses
// the statement without running it.
func TestPriorityFormsTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range priorityCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbRelease, accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, c.accept, c.sql, perr)
			}
		})
	}
}
