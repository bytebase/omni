package parser

import (
	"sort"
	"strings"
	"testing"
)

// TestSetTransactionFormsTiDBOracle asserts omni agrees with a live TiDB v8.5.5 on
// setTransactionCases (container test). The verdict comes from PREPARE, which parses
// the statement without running it.
func TestSetTransactionFormsTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range setTransactionCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = setTransactionCases[i-1].sql
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

// TestSetVariableNameSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5
// in each position of a variable name, because a word that SET refuses and every other
// statement reads fails nowhere else.
func TestSetVariableNameSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	words := make([]string, 0, len(keywords))
	for word := range keywords {
		words = append(words, word)
	}
	sort.Strings(words)
	previous := "<start of sweep>"
	for _, position := range setVariableNamePositions {
		for _, word := range words {
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
