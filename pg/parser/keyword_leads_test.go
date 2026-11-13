package parser

import "testing"

// TestKeywordLeadPredicates checks the lookahead sets keywordLeads decides
// on against the whole Keywords table, and IDENT.
func TestKeywordLeadPredicates(t *testing.T) {
	for _, kw := range Keywords {
		tok := Token{Type: kw.Token, Str: kw.Name}
		colID := kw.Category == UnreservedKeyword || kw.Category == ColNameKeyword
		if got, want := isColIdToken(tok), colID; got != want {
			t.Errorf("isColIdToken(%q) = %v, want %v", kw.Name, got, want)
		}
		if got, want := isOptTableQualifiedNameStart(tok), kw.Token == TABLE || colID; got != want {
			t.Errorf("isOptTableQualifiedNameStart(%q) = %v, want %v", kw.Name, got, want)
		}
		if got, want := isDeallocateTargetStart(tok), kw.Token == ALL || colID; got != want {
			t.Errorf("isDeallocateTargetStart(%q) = %v, want %v", kw.Name, got, want)
		}
		if got, want := isTempStart(tok), kw.Token == TEMP || kw.Token == TEMPORARY; got != want {
			t.Errorf("isTempStart(%q) = %v, want %v", kw.Name, got, want)
		}
	}
	ident := Token{Type: IDENT, Str: "a"}
	if !isColIdToken(ident) || !isOptTableQualifiedNameStart(ident) || !isDeallocateTargetStart(ident) || isTempStart(ident) {
		t.Error("IDENT: want a ColId that leads opt_table qualified_name and DEALLOCATE's target, and no TEMP")
	}
}

// TestVarNameDotNeedsColId: a var_name's dot is followed by a ColId
// (gram.y var_name: var_name '.' ColId), also where TIME, TRANSACTION or
// SESSION is the name.
func TestVarNameDotNeedsColId(t *testing.T) {
	for _, sql := range []string{"SHOW time.", "RESET session.", "SHOW transaction.", "ALTER DATABASE d RESET time.", "SET a. = 1"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
	for _, sql := range []string{"SHOW time", "SHOW time.zone", "RESET session", "FETCH next", "SELECT INTO temp"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}
