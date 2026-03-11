package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// ParseAndCheck is a test helper that parses SQL and returns the AST.
// It fails the test if parsing produces an error.
func ParseAndCheck(t *testing.T, sql string) *ast.List {
	t.Helper()
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): unexpected error: %v", sql, err)
	}
	return result
}

// ParseShouldFail is a test helper that expects parsing to fail.
func ParseShouldFail(t *testing.T, sql string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("Parse(%q): expected error, got success", sql)
	}
}

// TestParseEmpty tests parsing of empty input.
func TestParseEmpty(t *testing.T) {
	result := ParseAndCheck(t, "")
	if result.Len() != 0 {
		t.Fatalf("expected 0 statements, got %d", result.Len())
	}
}

// TestParseSemicolon tests parsing of standalone semicolons.
func TestParseSemicolon(t *testing.T) {
	result := ParseAndCheck(t, ";;;")
	if result.Len() != 0 {
		t.Fatalf("expected 0 statements, got %d", result.Len())
	}
}

// TestLexerBasicTokens tests basic lexer functionality.
func TestLexerBasicTokens(t *testing.T) {
	tests := []struct {
		input string
		typ   int
		str   string
	}{
		{"123", tokICONST, "123"},
		{"'hello'", tokSCONST, "hello"},
		{"\"MyCol\"", tokQIDENT, "MyCol"},
		{":param", tokBIND, "param"},
		{":1", tokBIND, "1"},
		{":=", tokASSIGN, ":="},
		{"=>", tokASSOC, "=>"},
		{"..", tokDOTDOT, ".."},
		{"**", tokEXPON, "**"},
		{"||", tokCONCAT, "||"},
		{"<=", tokLESSEQ, "<="},
		{">=", tokGREATEQ, ">="},
		{"!=", tokNOTEQ, "!="},
		{"<>", tokNOTEQ, "<>"},
		{"~=", tokNOTEQ, "~="},
		{"^=", tokNOTEQ, "^="},
		{"<<", tokLABELOPEN, "<<"},
		{">>", tokLABELCLOSE, ">>"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("type: got %d, want %d", tok.Type, tc.typ)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerStrings tests string literal lexing.
func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input string
		str   string
	}{
		{"'hello'", "hello"},
		{"'it''s'", "it's"},
		{"''", ""},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tokSCONST {
				t.Errorf("type: got %d, want %d", tok.Type, tokSCONST)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerQQuote tests q-quote string lexing.
func TestLexerQQuote(t *testing.T) {
	tests := []struct {
		input string
		str   string
	}{
		{"q'[hello world]'", "hello world"},
		{"q'{it's}'", "it's"},
		{"q'(test)'", "test"},
		{"q'<angle>'", "angle"},
		{"q'!bang!'", "bang"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tokSCONST {
				t.Errorf("type: got %d, want %d", tok.Type, tokSCONST)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
		})
	}
}

// TestLexerNCharLiteral tests N'...' national character literal.
func TestLexerNCharLiteral(t *testing.T) {
	lex := NewLexer("N'hello'")
	tok := lex.NextToken()
	if tok.Type != tokNCHARLIT {
		t.Errorf("type: got %d, want %d", tok.Type, tokNCHARLIT)
	}
	if tok.Str != "hello" {
		t.Errorf("str: got %q, want %q", tok.Str, "hello")
	}
}

// TestLexerHints tests hint lexing.
func TestLexerHints(t *testing.T) {
	lex := NewLexer("/*+ FULL(t) */")
	tok := lex.NextToken()
	if tok.Type != tokHINT {
		t.Errorf("type: got %d, want %d", tok.Type, tokHINT)
	}
	if tok.Str != "FULL(t)" {
		t.Errorf("str: got %q, want %q", tok.Str, "FULL(t)")
	}
}

// TestLexerComments tests that comments are skipped.
func TestLexerComments(t *testing.T) {
	lex := NewLexer("-- line comment\n42")
	tok := lex.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Errorf("expected 42, got type=%d ival=%d", tok.Type, tok.Ival)
	}

	lex = NewLexer("/* block comment */ 42")
	tok = lex.NextToken()
	if tok.Type != tokICONST || tok.Ival != 42 {
		t.Errorf("expected 42, got type=%d ival=%d", tok.Type, tok.Ival)
	}
}

// TestLexerKeywords tests keyword recognition.
func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input string
		typ   int
	}{
		{"SELECT", kwSELECT},
		{"select", kwSELECT},
		{"Select", kwSELECT},
		{"FROM", kwFROM},
		{"WHERE", kwWHERE},
		{"INSERT", kwINSERT},
		{"UPDATE", kwUPDATE},
		{"DELETE", kwDELETE},
		{"MERGE", kwMERGE},
		{"CREATE", kwCREATE},
		{"ALTER", kwALTER},
		{"DROP", kwDROP},
		{"CONNECT", kwCONNECT},
		{"PRIOR", kwPRIOR},
		{"DECODE", kwDECODE},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("keyword %q: got type %d, want %d", tc.input, tok.Type, tc.typ)
			}
		})
	}
}

// TestLexerNumbers tests number literal lexing.
func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input string
		typ   int
		str   string
		ival  int64
	}{
		{"0", tokICONST, "0", 0},
		{"42", tokICONST, "42", 42},
		{"3.14", tokFCONST, "3.14", 0},
		{".5", tokFCONST, ".5", 0},
		{"1e10", tokFCONST, "1e10", 0},
		{"1.5E-3", tokFCONST, "1.5E-3", 0},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			lex := NewLexer(tc.input)
			tok := lex.NextToken()
			if tok.Type != tc.typ {
				t.Errorf("type: got %d, want %d", tok.Type, tc.typ)
			}
			if tok.Str != tc.str {
				t.Errorf("str: got %q, want %q", tok.Str, tc.str)
			}
			if tc.typ == tokICONST && tok.Ival != tc.ival {
				t.Errorf("ival: got %d, want %d", tok.Ival, tc.ival)
			}
		})
	}
}

// TestNodeToString tests AST serialization.
func TestNodeToString(t *testing.T) {
	n := &ast.NumberLiteral{Val: "42", Ival: 42, Loc: ast.Loc{Start: 0, End: -1}}
	s := ast.NodeToString(n)
	if s == "" {
		t.Error("NodeToString returned empty string")
	}
}
