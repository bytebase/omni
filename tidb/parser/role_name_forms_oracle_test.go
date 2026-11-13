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
	for _, c := range roleNameCases {
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

// TestRoleNameSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5 in each
// position of a role name, which is what keeps tidbPlainIdentifiers true: a keyword added to
// the lexer that TiDB does not have fails here until it is listed.
func TestRoleNameSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, position := range roleNamePositions {
		for _, word := range registeredKeywords() {
			sql := strings.ReplaceAll(position, "%s", word)
			accepts, err := tc.prepareSyntax(sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			omniAccepts, perr := parsesAsOneStatement(sql)
			if accepts != omniAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, accepts, sql, perr)
			}
		}
	}
}
