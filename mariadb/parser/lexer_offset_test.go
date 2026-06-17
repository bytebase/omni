package parser

import "testing"

func TestLexerBaseOffset(t *testing.T) {
	input := "SELECT 1"
	baseOffset := 100

	lex := NewLexerWithOffset(input, baseOffset)
	tok := lex.NextToken()

	if tok.Type == tokEOF {
		t.Fatal("expected non-EOF token")
	}
	// First token "SELECT" starts at local offset 0 → absolute 100.
	if tok.Loc != baseOffset {
		t.Errorf("first token Loc = %d, want %d", tok.Loc, baseOffset)
	}

	// Second token "1" starts at local offset 7 → absolute 107.
	tok = lex.NextToken()
	if tok.Loc != baseOffset+7 {
		t.Errorf("second token Loc = %d, want %d", tok.Loc, baseOffset+7)
	}

	// EOF token.
	tok = lex.NextToken()
	if tok.Type != tokEOF {
		t.Errorf("expected EOF, got type %d", tok.Type)
	}
	if tok.Loc != baseOffset+len(input) {
		t.Errorf("EOF Loc = %d, want %d", tok.Loc, baseOffset+len(input))
	}
}

func TestLexerDefaultOffset(t *testing.T) {
	// NewLexer (no offset) should produce Loc == local offset.
	lex := NewLexer("SELECT 1")
	tok := lex.NextToken()
	if tok.Loc != 0 {
		t.Errorf("first token Loc = %d, want 0", tok.Loc)
	}
}
