package parser

import (
	"strings"
	"testing"
)

// parseExpectError parses sql and asserts that an error is returned whose
// message contains expectedSubstring.
func parseExpectError(t *testing.T, sql, expectedSubstring string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("expected error for %q, got nil", sql)
	}
	if !strings.Contains(err.Error(), expectedSubstring) {
		t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), expectedSubstring, sql)
	}
}

func TestParseError_Section_1_1_ErrorInfrastructure(t *testing.T) {
	// Verify that parsing invalid SQL returns an error (basic harness test).
	// The exact error message content will be improved in Phase 2+.
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		// Truncated CREATE returns error.
		{"truncated_create", "CREATE", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.sql)
			}
			if tc.contains != "" && !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), tc.contains, tc.sql)
			}
		})
	}
}

func TestSyntaxErrorAtCur_EOF(t *testing.T) {
	p := &Parser{
		lexer: NewLexer(""),
	}
	p.advance()

	err := p.syntaxErrorAtCur()
	if !strings.Contains(err.Error(), "at end of input") {
		t.Errorf("expected 'at end of input', got %q", err.Error())
	}
	if err.Line == 0 {
		t.Error("expected Line to be set")
	}
}

func TestSyntaxErrorAtTok_NonEOF(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT"),
	}
	p.advance()

	err := p.syntaxErrorAtTok(p.cur)
	if !strings.Contains(err.Error(), "at or near") {
		t.Errorf("expected 'at or near', got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "SELECT") {
		t.Errorf("expected token text 'SELECT' in error, got %q", err.Error())
	}
	if err.Line != 1 || err.Column != 1 {
		t.Errorf("expected line 1, column 1, got line %d, column %d", err.Line, err.Column)
	}
}

func TestParseError_PositionFormat(t *testing.T) {
	pe := &ParseError{
		Message:  "syntax error at end of input",
		Position: 10,
		Line:     1,
		Column:   11,
	}
	got := pe.Error()
	if got != "syntax error at end of input (line 1, column 11)" {
		t.Errorf("unexpected Error() output: %q", got)
	}
}

func TestParseError_NoPosition(t *testing.T) {
	pe := &ParseError{
		Message:  "unexpected token",
		Position: 5,
	}
	got := pe.Error()
	if got != "unexpected token" {
		t.Errorf("expected plain message, got %q", got)
	}
}

func TestLineCol(t *testing.T) {
	p := &Parser{lexer: NewLexer("SELECT\n  1 + 2")}
	// offset 0 -> line 1, col 1
	l, c := p.lineCol(0)
	if l != 1 || c != 1 {
		t.Errorf("offset 0: got line %d col %d, want 1 1", l, c)
	}
	// offset 7 (the space after newline) -> line 2, col 1
	l, c = p.lineCol(7)
	if l != 2 || c != 1 {
		t.Errorf("offset 7: got line %d col %d, want 2 1", l, c)
	}
	// offset 9 (the '1' on line 2) -> line 2, col 3
	l, c = p.lineCol(9)
	if l != 2 || c != 3 {
		t.Errorf("offset 9: got line %d col %d, want 2 3", l, c)
	}
}
