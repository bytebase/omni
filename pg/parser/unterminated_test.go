package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestParseUnterminatedConstructs checks that an unterminated lexical
// construct is reported wherever it appears. At the start of a statement the
// lexer hands the parser an EOF token with Err set, which used to look like
// end of input and was dropped silently.
func TestParseUnterminatedConstructs(t *testing.T) {
	constructs := []struct {
		name string
		text string // the unterminated construct, running to end of input
		msg  string
	}{
		{"quoted string", "'abc", "unterminated quoted string"},
		{"escape string", "E'abc", "unterminated quoted string"},
		{"quoted identifier", `"abc`, "unterminated quoted identifier"},
		{"block comment", "/* abc", "unterminated /* comment"},
		{"nested block comment", "/* a /* b */", "unterminated /* comment"},
		{"dollar-quoted string", "$$abc", "unterminated dollar-quoted string"},
		{"tagged dollar-quoted string", "$fn$abc$f$", "unterminated dollar-quoted string"},
	}
	positions := []struct {
		name   string
		prefix string
	}{
		{"at statement start", ""},
		{"after a previous statement", "SELECT 1; "},
		{"after a previous statement on a new line", "SELECT 1;\n\n  "},
		{"mid-statement", "SELECT "},
		{"mid-statement after a previous statement", "SELECT 1; SELECT 1, "},
	}
	for _, c := range constructs {
		for _, pos := range positions {
			t.Run(c.name+" "+pos.name, func(t *testing.T) {
				sql := pos.prefix + c.text
				list, err := Parse(sql)
				if err == nil {
					n := 0
					if list != nil {
						n = len(list.Items)
					}
					t.Fatalf("Parse(%q) returned %d statements and no error", sql, n)
				}
				if list != nil {
					t.Errorf("Parse(%q) returned a list alongside the error", sql)
				}
				var pe *ParseError
				if !errors.As(err, &pe) {
					t.Fatalf("Parse(%q) error is %T, want *ParseError", sql, err)
				}
				// PostgreSQL quotes the input from the construct's start to
				// the end (scanner_yyerror prints from yylloc).
				want := fmt.Sprintf("%s at or near \"%s\"", c.msg, c.text)
				if pe.Message != want {
					t.Errorf("Parse(%q) message = %q, want %q", sql, pe.Message, want)
				}
				if pe.Position != len(pos.prefix) {
					t.Errorf("Parse(%q) position = %d, want %d", sql, pe.Position, len(pos.prefix))
				}
				if !strings.Contains(err.Error(), "SQLSTATE 42601") {
					t.Errorf("Parse(%q) error = %q, want SQLSTATE 42601", sql, err.Error())
				}
			})
		}
	}
}

// TestParseLexerErrorBehindLookahead covers a lexer failure that is hit while
// peeking past a NOT/WITH token for reclassification: the error token sits in
// the lookahead buffer, and the report must still point at the construct.
func TestParseLexerErrorBehindLookahead(t *testing.T) {
	sql := "SELECT 1; WITH 'abc"
	_, err := Parse(sql)
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("Parse(%q) error = %v, want *ParseError", sql, err)
	}
	if want := `unterminated quoted string at or near "'abc"`; pe.Message != want {
		t.Errorf("message = %q, want %q", pe.Message, want)
	}
	if pe.Position != 15 {
		t.Errorf("position = %d, want 15", pe.Position)
	}
}

// TestParseZeroLengthIdentifierAtStatementStart: the other lexer failure that
// yields an EOF token. PostgreSQL quotes just the empty identifier.
func TestParseZeroLengthIdentifierAtStatementStart(t *testing.T) {
	sql := `SELECT 1; "" x`
	_, err := Parse(sql)
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("Parse(%q) error = %v, want *ParseError", sql, err)
	}
	if want := `zero-length delimited identifier at or near """"`; pe.Message != want {
		t.Errorf("message = %q, want %q", pe.Message, want)
	}
	if pe.Position != 10 {
		t.Errorf("position = %d, want 10", pe.Position)
	}
}
