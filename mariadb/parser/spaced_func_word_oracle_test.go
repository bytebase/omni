package parser

import "testing"

// TestMariaDBSpacedFuncWordsMatchEngine holds every row of spacedFuncWordCases to
// a live mariadb:11.8.8. prepareParses parses the statement and stops, telling a
// parse 1064 from a later semantic error (a missing column or function) by where
// it was raised, so a generic call to a name the server has no function for still
// counts as parsed.
func TestMariaDBSpacedFuncWordsMatchEngine(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range spacedFuncWordCases {
		if accepts := o.prepareAccepts(t, c.sql); accepts != c.accept {
			t.Errorf("mariadb:11.8.8 accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
