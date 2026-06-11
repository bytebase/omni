package parser

import (
	"testing"
)

func TestRecovery_MultilineStringAndUnterminated(t *testing.T) {
	// Snowflake string constants may span newlines, so "'unterm1\n'unterm2"
	// lexes as ONE multi-line string ("unterm1\n") followed by the identifier
	// unterm2 — zero errors. (Under the historical one-line rule this input
	// produced two unterminated-string errors; that rule falsely rejected
	// valid multi-line strings, e.g. the official-docs JSON literals.)
	input := "'unterm1\n'unterm2"
	tokens, errs := Tokenize(input)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
	if len(tokens) != 3 { // string, ident, EOF
		t.Fatalf("expected 3 tokens (string, ident, EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != tokString || tokens[0].Str != "unterm1\n" {
		t.Errorf("token[0] = %+v, want tokString with content \"unterm1\\n\"", tokens[0])
	}
	if tokens[1].Type != tokIdent || tokens[1].Str != "unterm2" {
		t.Errorf("token[1] = %+v, want tokIdent unterm2", tokens[1])
	}

	// A string with NO closing quote anywhere runs to EOF and reports exactly
	// one unterminated-string error.
	tokens, errs = Tokenize("SELECT 'never closed\nFROM t")
	if len(errs) != 1 {
		t.Errorf("unterminated: expected 1 error, got %d: %+v", len(errs), errs)
	}
	hasInvalid := false
	for _, tok := range tokens {
		if tok.Type == tokInvalid {
			hasInvalid = true
		}
	}
	if !hasInvalid {
		t.Errorf("unterminated: expected a tokInvalid token, got %+v", tokens)
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
