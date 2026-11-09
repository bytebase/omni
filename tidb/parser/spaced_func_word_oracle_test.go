package parser

import "testing"

// TestTiDBSpacedFuncWordsMatchEngine holds every row of spacedFuncWordCases to a
// live TiDB v8.5.5. prepareSyntax prepares the statement and tells a parser 1064
// from a later semantic error with the probe tail, the same verdict the reject
// fuzz and corpus runs use, so a column a builtin would reject when type-checked
// still counts as parsed.
func TestTiDBSpacedFuncWordsMatchEngine(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, c := range spacedFuncWordCases {
		accepts, err := tc.prepareSyntax(c.sql)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", c.sql, err)
		}
		if accepts != c.accept {
			t.Errorf("TiDB accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
