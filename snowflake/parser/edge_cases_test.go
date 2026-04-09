package parser

import (
	"strings"
	"testing"
)

func TestEdge_NestedBlockComments(t *testing.T) {
	// 3-deep nesting should produce no tokens and no errors.
	input := "/* a /* b /* c */ */ */"
	tokens, errs := Tokenize(input)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF, got %+v", tokens)
	}
}

func TestEdge_UnterminatedBlockComment(t *testing.T) {
	tokens, errs := Tokenize("/* unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated block comment") {
		t.Errorf("error message = %q, want substring 'unterminated block comment'", errs[0].Msg)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF token (comment is hidden), got %+v", tokens)
	}
}

func TestEdge_MismatchedBlockCommentDepth(t *testing.T) {
	// /* /* */  has depth 2 → 1 → 1 (not back to 0). Unterminated.
	tokens, errs := Tokenize("/* /* */")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if len(tokens) != 1 || tokens[0].Type != tokEOF {
		t.Errorf("expected only EOF token, got %+v", tokens)
	}
}

func TestEdge_UnterminatedString(t *testing.T) {
	tokens, errs := Tokenize("'unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated string literal") {
		t.Errorf("error = %q, want unterminated string", errs[0].Msg)
	}
	// Token stream: tokInvalid, EOF
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_UnterminatedDollarString(t *testing.T) {
	tokens, errs := Tokenize("$$ unterminated")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated $$ string literal") {
		t.Errorf("error = %q, want unterminated $$", errs[0].Msg)
	}
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_UnterminatedQuotedIdent(t *testing.T) {
	tokens, errs := Tokenize(`"unterminated`)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Msg, "unterminated quoted identifier") {
		t.Errorf("error = %q, want unterminated quoted identifier", errs[0].Msg)
	}
	if len(tokens) != 2 || tokens[0].Type != tokInvalid {
		t.Errorf("expected [tokInvalid, EOF], got %+v", tokens)
	}
}

func TestEdge_NumericLeadingDot(t *testing.T) {
	tokens, _ := Tokenize(".5")
	if tokens[0].Type != tokFloat || tokens[0].Str != ".5" {
		t.Errorf("got %+v, want tokFloat .5", tokens[0])
	}
}

func TestEdge_NumericTrailingDot(t *testing.T) {
	tokens, _ := Tokenize("1.")
	if tokens[0].Type != tokFloat || tokens[0].Str != "1." {
		t.Errorf("got %+v, want tokFloat 1.", tokens[0])
	}
}

func TestEdge_NumericExponentVariants(t *testing.T) {
	for _, input := range []string{"1e10", "1E10", "1e+10", "1e-10", ".5e10", "1.5e10"} {
		t.Run(input, func(t *testing.T) {
			tokens, _ := Tokenize(input)
			if tokens[0].Type != tokReal {
				t.Errorf("got %+v, want tokReal", tokens[0])
			}
		})
	}
}

func TestEdge_AmbiguousOperators(t *testing.T) {
	cases := []struct {
		input  string
		opType int
	}{
		{"a::b", tokDoubleColon},
		{"a||b", tokConcat},
		{"a->b", tokArrow},
		{"a->>b", tokFlow},
		{"a=>b", tokAssoc},
		{"a!=b", tokNotEq},
		{"a<>b", tokNotEq},
		{"a<=b", tokLessEq},
		{"a>=b", tokGreaterEq},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			// Expected: [ident, op, ident, EOF]
			if len(tokens) != 4 {
				t.Fatalf("expected 4 tokens, got %d: %+v", len(tokens), tokens)
			}
			if tokens[1].Type != c.opType {
				t.Errorf("op token = %d (%s), want %d (%s)",
					tokens[1].Type, TokenName(tokens[1].Type),
					c.opType, TokenName(c.opType))
			}
		})
	}
}

func TestEdge_UnicodeInQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"héllo"`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != "héllo" {
		t.Errorf("got %+v, want tokQuotedIdent héllo", tokens[0])
	}
}

func TestEdge_UnicodeInString(t *testing.T) {
	tokens, _ := Tokenize(`'café'`)
	if tokens[0].Type != tokString || tokens[0].Str != "café" {
		t.Errorf("got %+v, want tokString café", tokens[0])
	}
}

func TestEdge_DoubledQuoteInQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`"a""b"`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != `a"b` {
		t.Errorf("got %+v, want tokQuotedIdent a\"b", tokens[0])
	}
}

func TestEdge_DoubledQuoteInString(t *testing.T) {
	tokens, _ := Tokenize(`'a''b'`)
	if tokens[0].Type != tokString || tokens[0].Str != "a'b" {
		t.Errorf("got %+v, want tokString a'b", tokens[0])
	}
}

func TestEdge_EmptyQuotedIdent(t *testing.T) {
	tokens, _ := Tokenize(`""`)
	if tokens[0].Type != tokQuotedIdent || tokens[0].Str != "" {
		t.Errorf("got %+v, want tokQuotedIdent (empty)", tokens[0])
	}
}

func TestEdge_BareDollar(t *testing.T) {
	tokens, _ := Tokenize("$")
	if tokens[0].Type != int('$') {
		t.Errorf("got %+v, want int('$')", tokens[0])
	}
}

func TestEdge_BarePipe(t *testing.T) {
	// Bare | is invalid in Snowflake (only || is a token).
	tokens, errs := Tokenize("|")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if tokens[0].Type != tokInvalid {
		t.Errorf("got %+v, want tokInvalid", tokens[0])
	}
}
