package parser

import "testing"

// TestWildcardFieldFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// wildcardFieldCases (container test).
func TestWildcardFieldFormsTiDBOracle(t *testing.T) {
	holdFormsToTiDB(t, wildcardFieldCases)
}

// holdFormsToTiDB asserts that a live TiDB v8.5.5 gives each case the verdict it
// records, and that omni gives it too. The verdict comes from PREPARE, which parses
// the statement without running it.
func holdFormsToTiDB(t *testing.T, cases []tableSampleCase) {
	t.Helper()
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = cases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, c.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, c.accept, c.sql, perr)
			}
		})
	}
}
