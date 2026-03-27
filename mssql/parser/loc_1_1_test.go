package parser

import (
	"testing"
)

// TestTokenEnd verifies that every token produced by the lexer has a correct
// End field (exclusive byte offset past the last byte of the token).
func TestTokenEnd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLoc int
		wantEnd int
		wantTyp int // 0 means "don't check"
	}{
		// EOF
		{name: "EOF empty", input: "", wantLoc: 0, wantEnd: 0, wantTyp: tokEOF},
		{name: "EOF after space", input: "  ", wantLoc: 2, wantEnd: 2, wantTyp: tokEOF},

		// Single-char tokens
		{name: "semicolon", input: ";", wantLoc: 0, wantEnd: 1},
		{name: "lparen", input: "(", wantLoc: 0, wantEnd: 1},
		{name: "rparen", input: ")", wantLoc: 0, wantEnd: 1},
		{name: "comma", input: ",", wantLoc: 0, wantEnd: 1},
		{name: "dot", input: ". ", wantLoc: 0, wantEnd: 1}, // dot not followed by digit
		{name: "equals", input: "= ", wantLoc: 0, wantEnd: 1},

		// Multi-char operators
		{name: "<=", input: "<=", wantLoc: 0, wantEnd: 2},
		{name: ">=", input: ">=", wantLoc: 0, wantEnd: 2},
		{name: "<>", input: "<>", wantLoc: 0, wantEnd: 2},
		{name: "!=", input: "!=", wantLoc: 0, wantEnd: 2},
		{name: "+=", input: "+=", wantLoc: 0, wantEnd: 2},
		{name: "-=", input: "-=", wantLoc: 0, wantEnd: 2},
		{name: "*=", input: "*=", wantLoc: 0, wantEnd: 2},
		{name: "/=", input: "/=", wantLoc: 0, wantEnd: 2},
		{name: "%=", input: "%=", wantLoc: 0, wantEnd: 2},
		{name: "&=", input: "&=", wantLoc: 0, wantEnd: 2},
		{name: "|=", input: "|=", wantLoc: 0, wantEnd: 2},
		{name: "^=", input: "^=", wantLoc: 0, wantEnd: 2},
		{name: "!<", input: "!<", wantLoc: 0, wantEnd: 2},
		{name: "!>", input: "!>", wantLoc: 0, wantEnd: 2},
		{name: "::", input: "::", wantLoc: 0, wantEnd: 2},

		// Integer literals
		{name: "int 123", input: "123", wantLoc: 0, wantEnd: 3, wantTyp: tokICONST},
		{name: "int 0", input: "0 ", wantLoc: 0, wantEnd: 1, wantTyp: tokICONST},
		{name: "int 42 offset", input: "  42", wantLoc: 2, wantEnd: 4, wantTyp: tokICONST},

		// Hex literals
		{name: "hex 0xFF", input: "0xFF", wantLoc: 0, wantEnd: 4, wantTyp: tokICONST},
		{name: "hex 0x1A2B", input: "0x1A2B", wantLoc: 0, wantEnd: 6, wantTyp: tokICONST},

		// Scientific notation
		{name: "sci 1.5e10", input: "1.5e10", wantLoc: 0, wantEnd: 6, wantTyp: tokFCONST},
		{name: "sci 2E3", input: "2E3", wantLoc: 0, wantEnd: 3, wantTyp: tokFCONST},
		{name: "sci 1e+2", input: "1e+2", wantLoc: 0, wantEnd: 4, wantTyp: tokFCONST},

		// Float literals
		{name: "float 3.14", input: "3.14", wantLoc: 0, wantEnd: 4, wantTyp: tokFCONST},
		{name: "float .5", input: ".5", wantLoc: 0, wantEnd: 2, wantTyp: tokFCONST},

		// String literals
		{name: "string 'hello'", input: "'hello'", wantLoc: 0, wantEnd: 7, wantTyp: tokSCONST},
		{name: "string empty", input: "''", wantLoc: 0, wantEnd: 2, wantTyp: tokSCONST},
		{name: "string escaped", input: "'it''s'", wantLoc: 0, wantEnd: 7, wantTyp: tokSCONST},

		// N-string literals
		{name: "nstring N'hello'", input: "N'hello'", wantLoc: 0, wantEnd: 8, wantTyp: tokNSCONST},
		{name: "nstring n'hi'", input: "n'hi'", wantLoc: 0, wantEnd: 5, wantTyp: tokNSCONST},

		// Identifiers
		{name: "ident myTable", input: "myTable", wantLoc: 0, wantEnd: 7, wantTyp: tokIDENT},
		{name: "ident a", input: "a ", wantLoc: 0, wantEnd: 1, wantTyp: tokIDENT},

		// Bracketed identifiers
		{name: "bracket [my col]", input: "[my col]", wantLoc: 0, wantEnd: 8, wantTyp: tokIDENT},
		{name: "bracket escaped [col]]name]", input: "[col]]name]", wantLoc: 0, wantEnd: 11, wantTyp: tokIDENT},

		// Double-quoted identifiers
		{name: "dquote \"my col\"", input: `"my col"`, wantLoc: 0, wantEnd: 8, wantTyp: tokIDENT},

		// Keywords
		{name: "keyword SELECT", input: "SELECT", wantLoc: 0, wantEnd: 6, wantTyp: kwSELECT},
		{name: "keyword FROM", input: "FROM", wantLoc: 0, wantEnd: 4, wantTyp: kwFROM},
		{name: "keyword where", input: "where", wantLoc: 0, wantEnd: 5, wantTyp: kwWHERE},

		// Variables
		{name: "var @var", input: "@var", wantLoc: 0, wantEnd: 4, wantTyp: tokVARIABLE},
		{name: "var @x", input: "@x ", wantLoc: 0, wantEnd: 2, wantTyp: tokVARIABLE},

		// System variables
		{name: "sysvar @@ROWCOUNT", input: "@@ROWCOUNT", wantLoc: 0, wantEnd: 10, wantTyp: tokSYSVARIABLE},
		{name: "sysvar @@VERSION", input: "@@VERSION", wantLoc: 0, wantEnd: 9, wantTyp: tokSYSVARIABLE},

		// Whitespace/comments skipped — next token Loc starts after them
		{name: "space before", input: "   abc", wantLoc: 3, wantEnd: 6},
		{name: "line comment", input: "-- comment\nabc", wantLoc: 11, wantEnd: 14},
		{name: "block comment", input: "/* x */abc", wantLoc: 7, wantEnd: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewLexer(tt.input)
			tok := lex.NextToken()
			if lex.Err != nil {
				t.Fatalf("lexer error: %v", lex.Err)
			}
			if tok.Loc != tt.wantLoc {
				t.Errorf("Loc = %d, want %d", tok.Loc, tt.wantLoc)
			}
			if tok.End != tt.wantEnd {
				t.Errorf("End = %d, want %d", tok.End, tt.wantEnd)
			}
			if tt.wantTyp != 0 && tok.Type != tt.wantTyp {
				t.Errorf("Type = %d, want %d", tok.Type, tt.wantTyp)
			}
		})
	}
}

// TestTokenEndMultipleTokens verifies End values for consecutive tokens in a
// multi-token input, ensuring whitespace is properly skipped.
func TestTokenEndMultipleTokens(t *testing.T) {
	input := "SELECT 123 FROM"
	lex := NewLexer(input)

	// SELECT
	tok := lex.NextToken()
	if tok.Loc != 0 || tok.End != 6 {
		t.Errorf("SELECT: Loc=%d End=%d, want Loc=0 End=6", tok.Loc, tok.End)
	}

	// 123
	tok = lex.NextToken()
	if tok.Loc != 7 || tok.End != 10 {
		t.Errorf("123: Loc=%d End=%d, want Loc=7 End=10", tok.Loc, tok.End)
	}

	// FROM
	tok = lex.NextToken()
	if tok.Loc != 11 || tok.End != 15 {
		t.Errorf("FROM: Loc=%d End=%d, want Loc=11 End=15", tok.Loc, tok.End)
	}

	// EOF
	tok = lex.NextToken()
	if tok.Type != tokEOF || tok.Loc != 15 || tok.End != 15 {
		t.Errorf("EOF: Type=%d Loc=%d End=%d, want Type=0 Loc=15 End=15", tok.Type, tok.Loc, tok.End)
	}
}
