package parser

import "testing"

func collectToEOF(p *Parser) []Token {
	var toks []Token
	for p.cur.Type != tokEOF {
		toks = append(toks, p.cur)
		p.advance()
	}
	return toks
}

func assertSameTokens(t *testing.T, a, b []Token) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("token count: first=%d second=%d", len(a), len(b))
	}
	for i := range a {
		if a[i].Type != b[i].Type || a[i].Str != b[i].Str {
			t.Errorf("token %d: first={%d,%q} second={%d,%q}", i, a[i].Type, a[i].Str, b[i].Type, b[i].Str)
		}
	}
}

// Snapshot at '(', scan to EOF, restore, rescan — identical token stream.
func TestSnapshotRestore_RoundTrip(t *testing.T) {
	p := &Parser{lexer: NewLexerWithOffset("(SELECT 1) UNION (SELECT 2)", 0)}
	p.advance() // cur = '('
	snap := p.snapshotTokenStream()
	first := collectToEOF(p)
	p.restoreTokenStream(snap)
	assertSameTokens(t, first, collectToEOF(p))
}

// A conditional comment /*!NNNNN ... */ is spliced into l.input in place during
// lexing. Restore (which reverts l.input) must yield an identical rescan.
func TestSnapshotRestore_AcrossConditionalComment(t *testing.T) {
	p := &Parser{lexer: NewLexerWithOffset("(SELECT /*!40000 a */ FROM t)", 0)}
	p.advance() // cur = '('
	snap := p.snapshotTokenStream()
	first := collectToEOF(p)
	p.restoreTokenStream(snap)
	assertSameTokens(t, first, collectToEOF(p))
}
