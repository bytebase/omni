package parser

import "testing"

// TestIsNumericOnlyStartMatchesSlice checks isNumericOnlyStart: it accepts
// numericOnlyLeadTokens and no keyword.
func TestIsNumericOnlyStartMatchesSlice(t *testing.T) {
	for _, kw := range Keywords {
		p := &Parser{cur: Token{Type: kw.Token, Str: kw.Name}}
		if p.isNumericOnlyStart() {
			t.Errorf("isNumericOnlyStart(%q) = true", kw.Name)
		}
	}
	for _, tok := range numericOnlyLeadTokens {
		p := &Parser{cur: Token{Type: tok}}
		if !p.isNumericOnlyStart() {
			t.Errorf("isNumericOnlyStart must accept token %d", tok)
		}
	}
	for _, tok := range []int{IDENT, SCONST, ',', '('} {
		p := &Parser{cur: Token{Type: tok}}
		if p.isNumericOnlyStart() {
			t.Errorf("isNumericOnlyStart accepts token %d", tok)
		}
	}
}
