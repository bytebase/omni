package parser

import (
	"slices"
	"testing"
)

// TestIsSeqOptElemStartMatchesSlice checks isSeqOptElemStart against every
// keyword: it accepts exactly seqOptElemLeadTokens, NO only before CYCLE,
// MAXVALUE or MINVALUE.
func TestIsSeqOptElemStartMatchesSlice(t *testing.T) {
	for _, kw := range Keywords {
		if kw.Token == NO {
			continue
		}
		p := &Parser{cur: Token{Type: kw.Token, Str: kw.Name}}
		if got, want := p.isSeqOptElemStart(), slices.Contains(seqOptElemLeadTokens, kw.Token); got != want {
			t.Errorf("isSeqOptElemStart(%q) = %v, want %v", kw.Name, got, want)
		}
	}
	for sql, want := range map[string]bool{"NO CYCLE": true, "NO MAXVALUE": true, "NO MINVALUE": true, "NO INHERIT": false, "NO FORCE": false} {
		p := &Parser{lexer: NewLexer(sql)}
		p.advance()
		if got := p.isSeqOptElemStart(); got != want {
			t.Errorf("isSeqOptElemStart(%q) = %v, want %v", sql, got, want)
		}
	}
}
