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
	assertTiDBVersion(t, tc)
	for i, c := range reservedWordCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = reservedWordCases[i-1].sql
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

// TestReservedWordSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5 in each
// name position, which is what keeps keywordCategories TiDB's and not MySQL's. A word added to
// the lexer with the wrong category fails here.
func TestReservedWordSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	previous := "<start of sweep>"
	for _, position := range reservedWordPositions {
		for _, word := range registeredKeywords() {
			sql := strings.ReplaceAll(position, "%s", word)
			accepts, _, err := tc.syntaxVerdict(t, sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			previous = sql
			omniAccepts, perr := parsesAsOneStatement(sql)
			if reservedWordSweepGap[sql] {
				if accepts == omniAccepts {
					t.Errorf("%q now agrees: drop it from reservedWordSweepGap", sql)
				}
				continue
			}
			if accepts != omniAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, accepts, sql, perr)
			}
		}
	}
}
