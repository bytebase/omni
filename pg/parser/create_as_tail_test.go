package parser

import "testing"

// TestIsCreateStmtTailStartMatchesSlice checks isCreateStmtTailStart against
// every keyword: it accepts exactly createStmtTailLeadTokens.
func TestIsCreateStmtTailStartMatchesSlice(t *testing.T) {
	for _, kw := range Keywords {
		p := &Parser{cur: Token{Type: kw.Token, Str: kw.Name}}
		if got, want := p.isCreateStmtTailStart(), createStmtTailLeadSet[kw.Token]; got != want {
			t.Errorf("isCreateStmtTailStart(%q) = %v, want %v", kw.Name, got, want)
		}
	}
	for _, tok := range createStmtTailLeadTokens {
		p := &Parser{cur: Token{Type: tok}}
		if !p.isCreateStmtTailStart() {
			t.Errorf("isCreateStmtTailStart must accept token %d", tok)
		}
	}
}
