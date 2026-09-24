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

func checkLexerError(t *testing.T, sql, wantMsg string, wantPos int) {
	t.Helper()
	list, err := Parse(sql)
	if err == nil {
		n := 0
		if list != nil {
			n = len(list.Items)
		}
		t.Fatalf("Parse(%q) returned %d statements and no error", sql, n)
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("Parse(%q) error is %T, want *ParseError", sql, err)
	}
	if pe.Message != wantMsg {
		t.Errorf("Parse(%q) message = %q, want %q", sql, pe.Message, wantMsg)
	}
	if pe.Position != wantPos {
		t.Errorf("Parse(%q) position = %d, want %d", sql, pe.Position, wantPos)
	}
}

// TestParseLexerErrorBehindLookahead covers a lexer failure that is hit while
// peeking past a NOT/WITH token for reclassification: the error token sits in
// the lookahead buffer, whether the statement then ends cleanly or the grammar
// rejects the EOF token it sees. The lexer's report wins in both cases.
func TestParseLexerErrorBehindLookahead(t *testing.T) {
	checkLexerError(t, "SELECT 1; WITH 'abc", `unterminated quoted string at or near "'abc"`, 15)
	checkLexerError(t, "SELECT NOT 'abc", `unterminated quoted string at or near "'abc"`, 11)
	checkLexerError(t, "SELECT 1 WITH 'abc", `unterminated quoted string at or near "'abc"`, 14)
}

// TestParseLexerErrorAfterConsumedKeyword: a dispatcher that consumes its
// leading keyword and then returns no statement at the erroneous EOF must not
// look like a clean end of input.
func TestParseLexerErrorAfterConsumedKeyword(t *testing.T) {
	checkLexerError(t, "ALTER 'abc", `unterminated quoted string at or near "'abc"`, 6)
	checkLexerError(t, "DROP /* abc", `unterminated /* comment at or near "/* abc"`, 5)
	checkLexerError(t, "SELECT 1; SET /* abc", `unterminated /* comment at or near "/* abc"`, 14)
	checkLexerError(t, "SELECT 1; ALTER $$abc", `unterminated dollar-quoted string at or near "$$abc"`, 16)
}

// TestParseTrailingJunkQuotesJunk: PostgreSQL's junk rules match the literal
// plus the whole trailing identifier, so the diagnostic shows what made the
// token invalid rather than the valid prefix. Expectations checked against
// PostgreSQL 17.
func TestParseTrailingJunkQuotesJunk(t *testing.T) {
	checkLexerError(t, "SELECT $1foo", `trailing junk after parameter at or near "$1foo"`, 7)
	checkLexerError(t, "SELECT 123abc", `trailing junk after numeric literal at or near "123abc"`, 7)
	checkLexerError(t, "SELECT 123abc + 1", `trailing junk after numeric literal at or near "123abc"`, 7)
	checkLexerError(t, "SELECT 1.5abc", `trailing junk after numeric literal at or near "1.5abc"`, 7)
	checkLexerError(t, "SELECT 1e5abc", `trailing junk after numeric literal at or near "1e5abc"`, 7)
	checkLexerError(t, "SELECT 0x1fzz", `trailing junk after numeric literal at or near "0x1fzz"`, 7)
	checkLexerError(t, "SELECT 0o17zz", `trailing junk after numeric literal at or near "0o17zz"`, 7)
	checkLexerError(t, "SELECT 0b101zz", `trailing junk after numeric literal at or near "0b101zz"`, 7)
	// Flex's longest match: {integer_junk} ("0" + identifier) beats
	// {hexfail} unless the run after the prefix is empty or a lone "_".
	checkLexerError(t, "SELECT 0xg", `trailing junk after numeric literal at or near "0xg"`, 7)
	checkLexerError(t, "SELECT 0xg + 1", `trailing junk after numeric literal at or near "0xg"`, 7)
	checkLexerError(t, "SELECT 0o8foo", `trailing junk after numeric literal at or near "0o8foo"`, 7)
	checkLexerError(t, "SELECT 0b2foo", `trailing junk after numeric literal at or near "0b2foo"`, 7)
	checkLexerError(t, "SELECT 0o9", `trailing junk after numeric literal at or near "0o9"`, 7)
	checkLexerError(t, "SELECT 0b12", `trailing junk after numeric literal at or near "0b12"`, 7)
	checkLexerError(t, "SELECT 0x_g", `trailing junk after numeric literal at or near "0x_g"`, 7)
	checkLexerError(t, "SELECT 0x1_", `trailing junk after numeric literal at or near "0x1_"`, 7)
	checkLexerError(t, "SELECT 0x1__2", `trailing junk after numeric literal at or near "0x1__2"`, 7)
	checkLexerError(t, "SELECT 0x", `invalid hexadecimal integer at or near "0x"`, 7)
	checkLexerError(t, "SELECT 0x ", `invalid hexadecimal integer at or near "0x"`, 7)
	checkLexerError(t, "SELECT 0x_", `invalid hexadecimal integer at or near "0x_"`, 7)
	checkLexerError(t, "SELECT 0o", `invalid octal integer at or near "0o"`, 7)
	checkLexerError(t, "SELECT 0b", `invalid binary integer at or near "0b"`, 7)
	for _, sql := range []string{"SELECT 0x_1", "SELECT 0x1_f", "SELECT 0o1_7", "SELECT 0b1_0", "SELECT 0XFF"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) = %v, want a valid radix literal", sql, err)
		}
	}
	checkLexerError(t, "SELECT 1e", `trailing junk after numeric literal at or near "1e"`, 7)
	checkLexerError(t, "SELECT 1e foo", `trailing junk after numeric literal at or near "1e"`, 7)
	checkLexerError(t, "SELECT 1efoo", `trailing junk after numeric literal at or near "1efoo"`, 7)
	checkLexerError(t, "SELECT 1ee5", `trailing junk after numeric literal at or near "1ee5"`, 7)
	checkLexerError(t, "SELECT 1e_x", `trailing junk after numeric literal at or near "1e_x"`, 7)
	checkLexerError(t, "SELECT 1E$x", `trailing junk after numeric literal at or near "1E$x"`, 7)
	checkLexerError(t, "SELECT 1.5efoo", `trailing junk after numeric literal at or near "1.5efoo"`, 7)
	checkLexerError(t, "SELECT 1e+foo", `trailing junk after numeric literal at or near "1e+"`, 7)
	checkLexerError(t, "SELECT 1e-", `trailing junk after numeric literal at or near "1e-"`, 7)
	checkLexerError(t, "SELECT 1.5e+foo", `trailing junk after numeric literal at or near "1.5e+"`, 7)
	checkLexerError(t, "SELECT 123日本", `trailing junk after numeric literal at or near "123日本"`, 7)
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
