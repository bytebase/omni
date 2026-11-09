package parser

import "testing"

// TestSpacedFuncWordsMatchEngine holds every row of spacedFuncWordCases to a live
// mysql:8.0. prepareParses parses the statement and stops, telling a parse 1064
// from a later semantic error (a missing column or function) by where it was
// raised, so a generic call to a name the server has no function for still counts
// as parsed.
func TestSpacedFuncWordsMatchEngine(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range spacedFuncWordCases {
		if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
			t.Errorf("mysql:8.0 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
