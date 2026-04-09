package parser

import (
	"testing"
)

func TestRecovery_TwoUnterminatedStrings(t *testing.T) {
	// 'unterm1' is actually terminated; 'unterm2 is not. So this only
	// produces one real error. Use a multiline form for two errors.
	input := "'unterm1\n'unterm2"
	tokens, errs := Tokenize(input)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %+v", len(errs), errs)
	}
	if len(tokens) < 1 {
		t.Errorf("expected at least 1 token, got 0")
	}
}

func TestRecovery_KeywordBetweenErrors(t *testing.T) {
	// SELECT preserved between two invalid bytes.
	input := "| SELECT |"
	tokens, errs := Tokenize(input)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %+v", len(errs), errs)
	}
	// Stream should contain at least: tokInvalid, kwSELECT, tokInvalid, EOF
	hasSelect := false
	for _, tok := range tokens {
		if tok.Type == kwSELECT {
			hasSelect = true
			break
		}
	}
	if !hasSelect {
		t.Errorf("expected kwSELECT in stream, got %+v", tokens)
	}
}

func TestRecovery_ContinuesPastInvalidByte(t *testing.T) {
	// A NUL byte in the middle of input should be reported as one error
	// and lexing should continue past it.
	input := "SELECT \x00 FROM"
	tokens, errs := Tokenize(input)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	hasSelect := false
	hasFrom := false
	for _, tok := range tokens {
		if tok.Type == kwSELECT {
			hasSelect = true
		}
		if tok.Type == kwFROM {
			hasFrom = true
		}
	}
	if !hasSelect || !hasFrom {
		t.Errorf("expected SELECT and FROM preserved, got %+v", tokens)
	}
}

func TestRecovery_ReachesEOF(t *testing.T) {
	// Lexer must always reach EOF, even with errors.
	input := "| | | |"
	tokens, _ := Tokenize(input)
	if tokens[len(tokens)-1].Type != tokEOF {
		t.Errorf("did not reach EOF; last token = %+v", tokens[len(tokens)-1])
	}
}
