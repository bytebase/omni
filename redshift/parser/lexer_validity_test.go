package parser

import (
	"errors"
	"testing"
)

// TestLexerRejectsWhatScanLRejects is pg/parser's test, for the lexer this
// fork shares: it holds each token form PostgreSQL's scanner (scan.l, and
// base_yylex's Unicode handling in parser.c) rejects to the message,
// SQLSTATE and byte offset PostgreSQL 17.12 reports. A position of -1 is an
// error PostgreSQL reports without one.
func TestLexerRejectsWhatScanLRejects(t *testing.T) {
	tests := []struct {
		sql  string
		msg  string
		code string
		pos  int
	}{
		// {xeunicodefail}: \u needs four hex digits, \U eight.
		{`SELECT E'\u00'`, `invalid Unicode escape`, "22025", 9},
		{`SELECT E'\U0000004'`, `invalid Unicode escape`, "22025", 9},
		{`SELECT E'\u'`, `invalid Unicode escape`, "22025", 9},
		// addunicode: zero and values past U+10FFFF.
		{`SELECT E'\u0000'`, `invalid Unicode escape value at or near "\u0000"`, "", 9},
		{`SELECT E'\U00110000'`, `invalid Unicode escape value at or near "\U00110000"`, "", 9},
		// Surrogates: a second alone, a first without its second.
		{`SELECT E'\uDE00'`, `invalid Unicode surrogate pair at or near "\uDE00"`, "", 9},
		{`SELECT E'\uD83D'`, `invalid Unicode surrogate pair at or near "'"`, "", 15},
		{`SELECT E'\uD83Dx'`, `invalid Unicode surrogate pair at or near "x"`, "", 15},
		{`SELECT E'\uD83DA'`, `invalid Unicode surrogate pair at or near "A"`, "", 15},
		{`SELECT E'\uD83D\u00'`, `invalid Unicode escape`, "22025", 15},
		{`SELECT E'\uD83D`, `invalid Unicode surrogate pair`, "", 15},
		// pg_verifymbstr: a NUL or invalid UTF-8 made by an escape.
		{`SELECT E'\0'`, `invalid byte sequence for encoding "UTF8": 0x00`, "22021", -1},
		{`SELECT E'\400'`, `invalid byte sequence for encoding "UTF8": 0x00`, "22021", -1},
		{`SELECT E'\xff'`, `invalid byte sequence for encoding "UTF8": 0xff`, "22021", -1},
		{`SELECT E'\xc3' || 'x'`, `invalid byte sequence for encoding "UTF8": 0xc3`, "22021", -1},
		{"SELECT E'a'\n'\\xff'", `invalid byte sequence for encoding "UTF8": 0xff`, "22021", -1},
		// str_udeescape points at the escape's offset in the literal.
		{`SELECT U&'\00'`, `invalid Unicode escape`, "", 10},
		{`SELECT U&'\+00006'`, `invalid Unicode escape`, "", 10},
		{`SELECT U&'\'`, `invalid Unicode escape`, "", 10},
		{`SELECT U&'\0000'`, `invalid Unicode escape value`, "", 10},
		{`SELECT U&'\+110000'`, `invalid Unicode escape value`, "", 10},
		{`SELECT U&'\D83D'`, `invalid Unicode surrogate pair`, "", 15},
		{`SELECT U&'\DE00'`, `invalid Unicode surrogate pair`, "", 10},
		{`SELECT U&'\D83Dx'`, `invalid Unicode surrogate pair`, "", 15},
		{`SELECT U&"\00"`, `invalid Unicode escape`, "", 10},
		{`SELECT U&'a''b\0000'`, `invalid Unicode escape value`, "", 13},
		{"SELECT U&'a'\n'\\0000'", `invalid Unicode escape value`, "", 11},
		// UESCAPE takes a one-byte string that is no hex digit, '+',
		// quote or space.
		{`SELECT U&'a' UESCAPE '+'`, `invalid Unicode escape character at or near "'+'"`, "", 21},
		{`SELECT U&'a' UESCAPE 'ab'`, `invalid Unicode escape character at or near "'ab'"`, "", 21},
		{`SELECT U&'a' UESCAPE 'é'`, `invalid Unicode escape character at or near "'é'"`, "", 21},
		{`SELECT U&'a' UESCAPE ''''`, `invalid Unicode escape character at or near "''''"`, "", 21},
		{`SELECT U&'a' UESCAPE`, `UESCAPE must be followed by a simple string literal`, "", 20},
		{`SELECT U&'a' UESCAPE 1`, `UESCAPE must be followed by a simple string literal at or near "1"`, "", 21},
		{`SELECT U&'a' UESCAPE U&'!'`, `UESCAPE must be followed by a simple string literal at or near "U&'!'"`, "", 21},
		{`SELECT U&'a' UESCAPE B'1'`, `UESCAPE must be followed by a simple string literal at or near "B'1'"`, "", 21},
		// The token after a Unicode literal is lexed before it is decoded.
		{`SELECT U&'\0000' 1_`, `trailing junk after numeric literal at or near "1_"`, "", 17},
		{`SELECT U&'\0000' 'abc`, `unterminated quoted string at or near "'abc"`, "", 17},
		{`SELECT U&""`, `zero-length delimited identifier at or near "U&"""`, "", 7},
		{`SELECT U&"" UESCAPE`, `zero-length delimited identifier at or near "U&"""`, "", 7},
		// {decinteger}: an underscore only between two digits.
		{`SELECT 1__000`, `trailing junk after numeric literal at or near "1__000"`, "", 7},
		{`SELECT 1_`, `trailing junk after numeric literal at or near "1_"`, "", 7},
		{`SELECT 1._5`, `trailing junk after numeric literal at or near "1._5"`, "", 7},
		{`SELECT 1_.5`, `trailing junk after numeric literal at or near "1_"`, "", 7},
		{`SELECT 1.5_`, `trailing junk after numeric literal at or near "1.5_"`, "", 7},
		{`SELECT .5_`, `trailing junk after numeric literal at or near ".5_"`, "", 7},
		{`SELECT 1e+1_`, `trailing junk after numeric literal at or near "1e+1_"`, "", 7},
		// A construct opened by the last bytes of input is unterminated.
		{`SELECT '`, `unterminated quoted string at or near "'"`, "", 7},
		{`SELECT E'`, `unterminated quoted string at or near "E'"`, "", 7},
		{`SELECT "`, `unterminated quoted identifier at or near """`, "", 7},
		{`SELECT $$`, `unterminated dollar-quoted string at or near "$$"`, "", 7},
		{`SELECT B'`, `unterminated bit string literal at or near "B'"`, "", 7},
		{`SELECT U&'`, `unterminated quoted string at or near "U&'"`, "", 7},
		{`SELECT /*`, `unterminated /* comment at or near "/*"`, "", 7},
		{"SELECT 'a'\n'", "unterminated quoted string at or near \"'a'\n'\"", "", 7},
	}
	for _, tt := range tests {
		_, err := Parse(tt.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a *ParseError", tt.sql, err)
			continue
		}
		if pe.Message != tt.msg || pe.Code != tt.code || pe.Position != tt.pos {
			t.Errorf("Parse(%q) = %q code %q at %d, want %q code %q at %d",
				tt.sql, pe.Message, pe.Code, pe.Position, tt.msg, tt.code, tt.pos)
		}
	}
}

// TestLexerAcceptsWhatScanLAccepts covers the valid neighbours of the forms
// above, each accepted by PostgreSQL 17.12, and the values they decode to.
func TestLexerAcceptsWhatScanLAccepts(t *testing.T) {
	tests := []struct {
		sql  string
		want string // the first token's string value
	}{
		{`E'\U0010FFFF'`, "\U0010FFFF"},
		{`E'😀'`, "\U0001F600"},
		{`E'\U0000D83D\U0000DE00'`, "\U0001F600"},
		{`E'\xc3\xa9'`, "é"},
		{`E'\x'`, "x"},
		{`U&'\D83D\DE00'`, "\U0001F600"},
		{`U&'\+01F600'`, "\U0001F600"},
		{`U&'a\\b'`, `a\b`},
		{`U&'!0041' UESCAPE '!'`, "A"},
		{`U&'!0041' uescape $$!$$`, "A"},
		{`U&'!0041' UESCAPE E'!'`, "A"},
		{`U&'!0041' /*c*/ UESCAPE /*d*/ '!'`, "A"},
		{`U&'\0061' UESCAPE '\'`, "a"},
		{`U&"!0041" UESCAPE '!'`, "A"},
		{`1_000`, "1000"},
		{`0.0_1`, "0.01"},
		{`1_0e+1_0`, "10e+10"},
		{`1_000.000_1e1_0`, "1000.0001e10"},
	}
	for _, tt := range tests {
		if _, err := Parse("SELECT " + tt.sql); err != nil {
			t.Errorf("Parse(%q) = %v, want it accepted", "SELECT "+tt.sql, err)
			continue
		}
		l := NewLexer(tt.sql)
		tok := l.NextToken()
		if l.Err != nil || tok.Str != tt.want {
			t.Errorf("lexing %q = %q (err %v), want %q", tt.sql, tok.Str, l.Err, tt.want)
		}
		if next := l.NextToken(); next.Type != lex_EOF || l.Err != nil {
			t.Errorf("lexing %q left token %+v (err %v), want end of input", tt.sql, next, l.Err)
		}
	}
	// "uescapex" is an identifier, not UESCAPE: here the column's alias.
	if _, err := Parse(`SELECT U&'a' uescapex`); err != nil {
		t.Errorf("Parse(U&'a' uescapex) = %v, want it accepted", err)
	}
}
