package parser

import (
	"strings"
	"testing"
)

// TestRoleNameFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on roleNameCases
// (container test). The verdict comes from PREPARE, which parses the statement without
// running it.
func TestRoleNameFormsTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range roleNameCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = roleNameCases[i-1].sql
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

// TestRoleNameSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5 in each
// position of a role name, which is what keeps tidbPlainIdentifiers true: a keyword added to
// the lexer that TiDB does not have fails here until it is listed.
func TestRoleNameSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	previous := "<start of sweep>"
	for _, position := range roleNamePositions {
		for _, word := range registeredKeywords() {
			sql := strings.ReplaceAll(position, "%s", word)
			accepts, _, err := tc.syntaxVerdict(t, sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			previous = sql
			omniAccepts, perr := parsesAsOneStatement(sql)
			if accepts != omniAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, accepts, sql, perr)
			}
		}
	}
}
