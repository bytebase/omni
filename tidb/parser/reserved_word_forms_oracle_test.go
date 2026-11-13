package parser

import (
	"strings"
	"testing"
)

// TestReservedWordFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// reservedWordCases (container test). The verdict comes from PREPARE, which parses the
// statement without running it.
func TestReservedWordFormsTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range reservedWordCases {
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

// TestReservedWordSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5 in each
// name position, which is what keeps keywordCategories TiDB's and not MySQL's. A word added to
// the lexer with the wrong category fails here.
func TestReservedWordSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, position := range reservedWordPositions {
		for _, word := range registeredKeywords() {
			sql := strings.ReplaceAll(position, "%s", word)
			accepts, err := tc.prepareSyntax(sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			omniAccepts, perr := parsesAsOneStatement(sql)
			if reservedWordSweepGap[sql] {
				if accepts == omniAccepts {
					t.Errorf("%q now agrees: drop it from reservedWordSweepGap", sql)
				}
				continue
			}
			if accepts != omniAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, accepts, sql, perr)
			}
		}
	}
}
