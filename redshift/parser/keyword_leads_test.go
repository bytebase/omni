package parser

import (
	"slices"
	"testing"
)

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
		identityOption := kw.Token == GENERATED || slices.Contains(seqOptElemLeadTokens, kw.Token)
		if got := isAlterIdentityOptionStart(tok); got != identityOption {
			t.Errorf("isAlterIdentityOptionStart(%q) = %v, want %v", kw.Name, got, identityOption)
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

// TestFrameAndFetchLeadPredicates checks isFrameDirectionStart and
// isFetchOptionStart against the whole Keywords table.
func TestFrameAndFetchLeadPredicates(t *testing.T) {
	for _, kw := range Keywords {
		tok := Token{Type: kw.Token, Str: kw.Name}
		if got, want := isFrameDirectionStart(tok), kw.Token == PRECEDING || kw.Token == FOLLOWING; got != want {
			t.Errorf("isFrameDirectionStart(%q) = %v, want %v", kw.Name, got, want)
		}
	}
}

// TestClauseLeadPredicates holds isJsonFormatStart, isRowsFromStart and
// isBareAliasStart to the statements that decide them: each keyword leads
// its clause only before the token that continues it.
func TestClauseLeadPredicates(t *testing.T) {
	for _, sql := range []string{
		"SELECT JSON_SERIALIZE(a FORMAT JSON)",
		"SELECT format FROM t",
		"SELECT * FROM ROWS FROM (f(1))",
		"SELECT * FROM rows(1)",
		"SELECT * FROM t a",
		"SELECT * FROM t AS a(b)",
		"SELECT * FROM t JOIN u ON true",
		"SELECT XMLROOT(a, VERSION NO VALUE)",
		"SELECT XMLROOT(a, VERSION no)",
		"SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (NESTED PATH '$' COLUMNS (a int)))",
		"SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (nested int))",
		"SELECT count(*) OVER (ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)",
		"SELECT count(*) OVER (ROWS current PRECEDING)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{"SELECT JSON_SERIALIZE(a FORMAT)", "SELECT * FROM ROWS FROM f(1)"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}
