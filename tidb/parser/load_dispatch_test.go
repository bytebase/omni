package parser

import "testing"

// TestLoadDispatch: LOAD takes DATA, XML, INDEX or STATS after it, and any
// other word is a syntax error rather than a LOAD DATA.
func TestLoadDispatch(t *testing.T) {
	if _, err := Parse("LOAD DATABASE INFILE 's' INTO TABLE t"); err == nil {
		t.Error("Parse(LOAD DATABASE ...) succeeded, want a syntax error")
	}
	if _, err := Parse("LOAD DATA INFILE 's' INTO TABLE t"); err != nil {
		t.Errorf("Parse(LOAD DATA ...): %v", err)
	}
}
