package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

func TestPosition_TokenAtStart(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	if got := tokens[0].Loc; got != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", got)
	}
}

func TestPosition_TokenAfterWhitespace(t *testing.T) {
	tokens, _ := Tokenize("   SELECT")
	if got := tokens[0].Loc; got != (ast.Loc{Start: 3, End: 9}) {
		t.Errorf("Loc = %+v, want {3, 9}", got)
	}
}

func TestPosition_MultiCharOperator(t *testing.T) {
	tokens, _ := Tokenize("a >= b")
	// expect: ident, >=, ident
	if tokens[1].Type != tokGreaterEqual {
		t.Fatalf("token[1] = %+v, expected >=", tokens[1])
	}
	if tokens[1].Loc != (ast.Loc{Start: 2, End: 4}) {
		t.Errorf("Loc = %+v, want {2, 4}", tokens[1].Loc)
	}
}

func TestPosition_Pipe(t *testing.T) {
	tokens, _ := Tokenize("a |> b")
	if tokens[1].Type != tokPipe {
		t.Fatalf("token[1] = %+v, expected |>", tokens[1])
	}
	if tokens[1].Loc != (ast.Loc{Start: 2, End: 4}) {
		t.Errorf("Loc = %+v, want {2, 4}", tokens[1].Loc)
	}
}

func TestPosition_StringLiteral(t *testing.T) {
	tokens, _ := Tokenize("'hello'")
	// Loc.Start is the opening quote; Loc.End is one past the closing quote.
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 7}) {
		t.Errorf("Loc = %+v, want {0, 7}", tokens[0].Loc)
	}
}

func TestPosition_RawString(t *testing.T) {
	// r'abc' — Loc covers the r prefix through the closing quote (6 bytes).
	tokens, _ := Tokenize(`r'abc'`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

func TestPosition_TripleString(t *testing.T) {
	// '''abc''' is 9 bytes.
	tokens, _ := Tokenize(`'''abc'''`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 9}) {
		t.Errorf("Loc = %+v, want {0, 9}", tokens[0].Loc)
	}
}

func TestPosition_BytesLiteral(t *testing.T) {
	// b'abc' is 6 bytes (prefix + quotes + body).
	tokens, _ := Tokenize(`b'abc'`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

func TestPosition_BacktickIdent(t *testing.T) {
	// `my col` is 8 bytes (backticks included in the span).
	tokens, _ := Tokenize("`my col`")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 8}) {
		t.Errorf("Loc = %+v, want {0, 8}", tokens[0].Loc)
	}
}

func TestPosition_Number(t *testing.T) {
	cases := []struct {
		input string
		want  ast.Loc
	}{
		{"42", ast.Loc{Start: 0, End: 2}},
		{"0x1f", ast.Loc{Start: 0, End: 4}},
		{"1.5", ast.Loc{Start: 0, End: 3}},
		{".5", ast.Loc{Start: 0, End: 2}},
		{"1e10", ast.Loc{Start: 0, End: 4}},
		{"1.5e-10", ast.Loc{Start: 0, End: 7}},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, _ := Tokenize(c.input)
			if got := tokens[0].Loc; got != c.want {
				t.Errorf("Loc = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestPosition_EOFToken(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	last := tokens[len(tokens)-1]
	if last.Type != tokEOF {
		t.Fatalf("last token Type = %d, want tokEOF", last.Type)
	}
	if last.Loc != (ast.Loc{Start: 6, End: 6}) {
		t.Errorf("EOF Loc = %+v, want {6, 6}", last.Loc)
	}
}

func TestPosition_EOFTokenEmptyInput(t *testing.T) {
	tokens, _ := Tokenize("")
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token (just EOF), got %d", len(tokens))
	}
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 0}) {
		t.Errorf("EOF Loc = %+v, want {0, 0}", tokens[0].Loc)
	}
}

func TestLexer_NewLexerWithOffset(t *testing.T) {
	// Tokenize "   SELECT   " with baseOffset=100. kwSELECT should report
	// Loc.Start = 103, Loc.End = 109.
	tokens, errs := TokenizeWithOffset("   SELECT   ", 100)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	var selectTok *Token
	for i, tok := range tokens {
		if tok.Type == kwSELECT {
			selectTok = &tokens[i]
			break
		}
	}
	if selectTok == nil {
		t.Fatalf("kwSELECT not found in tokens: %+v", tokens)
	}
	if selectTok.Loc.Start != 103 || selectTok.Loc.End != 109 {
		t.Errorf("kwSELECT Loc = %+v, want {103, 109}", selectTok.Loc)
	}
}

func TestLexer_NewLexerWithOffset_LexErrors(t *testing.T) {
	// An unterminated string should produce a LexError whose Loc is shifted by
	// baseOffset.
	_, errs := TokenizeWithOffset("'unterminated", 50)
	if len(errs) != 1 {
		t.Fatalf("expected 1 lex error, got %d: %+v", len(errs), errs)
	}
	if errs[0].Loc.Start != 50 {
		t.Errorf("lex error Loc.Start = %d, want 50", errs[0].Loc.Start)
	}
}

func TestLexer_NewLexer_UnchangedBehavior(t *testing.T) {
	// Regression: NewLexer produces zero-based positions.
	tokens, errs := Tokenize("SELECT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) < 1 || tokens[0].Type != kwSELECT {
		t.Fatalf("expected first token kwSELECT, got %+v", tokens)
	}
	if tokens[0].Loc.Start != 0 || tokens[0].Loc.End != 6 {
		t.Errorf("kwSELECT Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

// TokenizeWithOffset is a test helper: one-shot tokenize with a base offset.
func TokenizeWithOffset(input string, baseOffset int) (tokens []Token, errors []LexError) {
	l := NewLexerWithOffset(input, baseOffset)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}
