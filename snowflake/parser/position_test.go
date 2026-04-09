package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
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

func TestPosition_TokenAtEnd(t *testing.T) {
	tokens, _ := Tokenize("SELECT")
	if got := tokens[0].Loc.End; got != 6 {
		t.Errorf("Loc.End = %d, want 6", got)
	}
}

func TestPosition_MultiCharOperator(t *testing.T) {
	tokens, _ := Tokenize("a::b")
	// expect: ident, ::, ident
	if tokens[1].Type != tokDoubleColon {
		t.Fatalf("token[1] = %+v, expected ::", tokens[1])
	}
	if tokens[1].Loc != (ast.Loc{Start: 1, End: 3}) {
		t.Errorf("Loc = %+v, want {1, 3}", tokens[1].Loc)
	}
}

func TestPosition_StringLiteral(t *testing.T) {
	tokens, _ := Tokenize("'hello'")
	// Loc.Start is the opening quote, Loc.End is one past the closing quote.
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 7}) {
		t.Errorf("Loc = %+v, want {0, 7}", tokens[0].Loc)
	}
}

func TestPosition_StringWithEscape(t *testing.T) {
	tokens, _ := Tokenize(`'a\nb'`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 6}) {
		t.Errorf("Loc = %+v, want {0, 6}", tokens[0].Loc)
	}
}

func TestPosition_QuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"my col"`)
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 8}) {
		t.Errorf("Loc = %+v, want {0, 8}", tokens[0].Loc)
	}
}

func TestPosition_DollarString(t *testing.T) {
	tokens, _ := Tokenize("$$hello$$")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 9}) {
		t.Errorf("Loc = %+v, want {0, 9}", tokens[0].Loc)
	}
}

func TestPosition_Variable(t *testing.T) {
	tokens, _ := Tokenize("$abc")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 4}) {
		t.Errorf("Loc = %+v, want {0, 4}", tokens[0].Loc)
	}
}

func TestPosition_Number(t *testing.T) {
	cases := []struct {
		input string
		want  ast.Loc
	}{
		{"42", ast.Loc{Start: 0, End: 2}},
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

func TestPosition_HexLiteral(t *testing.T) {
	// X'48656C6C6F' is 14 bytes. Loc.Start should be the X (offset 0).
	tokens, _ := Tokenize("X'48656C6C6F'")
	if tokens[0].Loc != (ast.Loc{Start: 0, End: 13}) {
		t.Errorf("Loc = %+v, want {0, 13}", tokens[0].Loc)
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
