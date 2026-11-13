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
	for _, c := range setTransactionCases {
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

// TestSetVariableNameSweepTiDBOracle holds every registered keyword to a live TiDB v8.5.5
// in each position of a variable name, because a word that SET refuses and every other
// statement reads fails nowhere else.
func TestSetVariableNameSweepTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	words := make([]string, 0, len(keywords))
	for word := range keywords {
		words = append(words, word)
	}
	sort.Strings(words)
	for _, position := range setVariableNamePositions {
		for _, word := range words {
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
