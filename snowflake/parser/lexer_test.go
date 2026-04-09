package parser

import "testing"

// expect describes one expected token in the output stream. Loc is omitted
// here — position assertions live in position_test.go.
type expect struct {
	Type int
	Str  string
}

// runLexerCases drives a table-driven test. For each row, it tokenizes
// row.input and asserts the resulting non-EOF tokens match row.want by
// (Type, Str). Errors must be empty.
func runLexerCases(t *testing.T, cases []struct {
	name  string
	input string
	want  []expect
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, errs := Tokenize(c.input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			// Strip the trailing EOF token before comparison.
			if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
				t.Fatalf("expected stream to end with tokEOF, got %+v", tokens)
			}
			tokens = tokens[:len(tokens)-1]
			if len(tokens) != len(c.want) {
				t.Fatalf("token count mismatch: got %d, want %d\n got=%+v\nwant=%+v",
					len(tokens), len(c.want), tokens, c.want)
			}
			for i, tok := range tokens {
				if tok.Type != c.want[i].Type {
					t.Errorf("token[%d] Type = %d (%s), want %d (%s)",
						i, tok.Type, TokenName(tok.Type), c.want[i].Type, TokenName(c.want[i].Type))
				}
				if tok.Str != c.want[i].Str {
					t.Errorf("token[%d] Str = %q, want %q", i, tok.Str, c.want[i].Str)
				}
			}
		})
	}
}

func TestLexer_SingleCharOperators(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"plus", "+", []expect{{int('+'), ""}}},
		{"minus", "-", []expect{{int('-'), ""}}},
		{"star", "*", []expect{{int('*'), ""}}},
		{"slash", "/", []expect{{int('/'), ""}}},
		{"percent", "%", []expect{{int('%'), ""}}},
		{"tilde", "~", []expect{{int('~'), ""}}},
		{"lparen", "(", []expect{{int('('), ""}}},
		{"rparen", ")", []expect{{int(')'), ""}}},
		{"lbracket", "[", []expect{{int('['), ""}}},
		{"rbracket", "]", []expect{{int(']'), ""}}},
		{"lbrace", "{", []expect{{int('{'), ""}}},
		{"rbrace", "}", []expect{{int('}'), ""}}},
		{"comma", ",", []expect{{int(','), ""}}},
		{"semi", ";", []expect{{int(';'), ""}}},
		{"dot", ".", []expect{{int('.'), ""}}},
		{"at", "@", []expect{{int('@'), ""}}},
		{"colon", ":", []expect{{int(':'), ""}}},
		{"lt", "<", []expect{{int('<'), ""}}},
		{"gt", ">", []expect{{int('>'), ""}}},
		{"eq", "=", []expect{{int('='), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_MultiCharOperators(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"double_colon", "::", []expect{{tokDoubleColon, ""}}},
		{"concat", "||", []expect{{tokConcat, ""}}},
		{"arrow", "->", []expect{{tokArrow, ""}}},
		{"flow", "->>", []expect{{tokFlow, ""}}},
		{"assoc", "=>", []expect{{tokAssoc, ""}}},
		{"not_eq_bang", "!=", []expect{{tokNotEq, ""}}},
		{"not_eq_angle", "<>", []expect{{tokNotEq, ""}}},
		{"less_eq", "<=", []expect{{tokLessEq, ""}}},
		{"greater_eq", ">=", []expect{{tokGreaterEq, ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_Identifiers(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple_lower", "abc", []expect{{tokIdent, "abc"}}},
		{"simple_upper", "ABC", []expect{{tokIdent, "ABC"}}},
		{"underscore_start", "_x", []expect{{tokIdent, "_x"}}},
		{"with_digits", "x123", []expect{{tokIdent, "x123"}}},
		{"with_at", "x@y", []expect{{tokIdent, "x@y"}}},
		{"with_dollar", "x$y", []expect{{tokIdent, "x$y"}}},
		{"quoted_simple", `"my_col"`, []expect{{tokQuotedIdent, "my_col"}}},
		{"quoted_with_space", `"my col"`, []expect{{tokQuotedIdent, "my col"}}},
		{"quoted_doubled_quote", `"a""b"`, []expect{{tokQuotedIdent, `a"b`}}},
		{"quoted_empty", `""`, []expect{{tokQuotedIdent, ""}}},
		{"variable_simple", "$x", []expect{{tokVariable, "x"}}},
		{"variable_numeric", "$1", []expect{{tokVariable, "1"}}},
		{"variable_long", "$ABC_123", []expect{{tokVariable, "ABC_123"}}},
		{"bare_dollar", "$", []expect{{int('$'), ""}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_StringLiterals(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"simple", `'hello'`, []expect{{tokString, "hello"}}},
		{"empty", `''`, []expect{{tokString, ""}}},
		{"escape_n", `'a\nb'`, []expect{{tokString, "a\nb"}}},
		{"escape_t", `'a\tb'`, []expect{{tokString, "a\tb"}}},
		{"escape_backslash", `'a\\b'`, []expect{{tokString, `a\b`}}},
		{"escape_quote", `'a\'b'`, []expect{{tokString, "a'b"}}},
		{"doubled_quote", `'a''b'`, []expect{{tokString, "a'b"}}},
		{"dollar_simple", `$$hello$$`, []expect{{tokString, "hello"}}},
		{"dollar_empty", `$$$$`, []expect{{tokString, ""}}},
		{"dollar_with_quotes", `$$can't$$`, []expect{{tokString, "can't"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_HexBinaryLiterals(t *testing.T) {
	// X'...' is tokString with XPrefix=true. We can't use runLexerCases
	// for this since it doesn't compare XPrefix.
	tokens, errs := Tokenize("X'48656C6C6F'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 2 || tokens[0].Type != tokString || !tokens[0].XPrefix {
		t.Fatalf("expected single tokString with XPrefix=true, got %+v", tokens)
	}
	if tokens[0].Str != "48656C6C6F" {
		t.Errorf("Str = %q, want %q", tokens[0].Str, "48656C6C6F")
	}

	// Lowercase x prefix should also work.
	tokens, errs = Tokenize("x'DEADBEEF'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if !tokens[0].XPrefix {
		t.Errorf("lowercase x prefix not detected")
	}
}

func TestLexer_NumericLiterals(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"int_zero", "0", []expect{{tokInt, "0"}}},
		{"int_small", "42", []expect{{tokInt, "42"}}},
		{"int_large", "1234567890", []expect{{tokInt, "1234567890"}}},
		{"float_simple", "1.5", []expect{{tokFloat, "1.5"}}},
		{"float_trailing_dot", "1.", []expect{{tokFloat, "1."}}},
		{"float_leading_dot", ".5", []expect{{tokFloat, ".5"}}},
		{"real_simple", "1e10", []expect{{tokReal, "1e10"}}},
		{"real_pos_exp", "1e+10", []expect{{tokReal, "1e+10"}}},
		{"real_neg_exp", "1e-10", []expect{{tokReal, "1e-10"}}},
		{"real_float_base", "1.5e10", []expect{{tokReal, "1.5e10"}}},
		{"real_dot_base", ".5e10", []expect{{tokReal, ".5e10"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_NumericInt64Overflow(t *testing.T) {
	// Snowflake's NUMBER(38, 0) accepts arbitrary-precision integers.
	// When a literal exceeds int64, scanNumber downgrades to tokFloat.
	tokens, errs := Tokenize("99999999999999999999999999999999999999")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens (literal + EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != tokFloat {
		t.Errorf("overflow literal Type = %d (%s), want tokFloat",
			tokens[0].Type, TokenName(tokens[0].Type))
	}
	if tokens[0].Str != "99999999999999999999999999999999999999" {
		t.Errorf("overflow literal Str = %q, want full source text", tokens[0].Str)
	}
	if tokens[0].Ival != 0 {
		t.Errorf("overflow literal Ival = %d, want 0 (unset)", tokens[0].Ival)
	}
}

func TestLexer_Keywords(t *testing.T) {
	// Spot-check a few keywords from each category. The full coverage is
	// in keyword_completeness_test.go (which checks every keyword from the
	// legacy .g4).
	cases := []struct {
		name  string
		input string
		want  []expect
	}{
		{"select_lower", "select", []expect{{kwSELECT, "select"}}},
		{"select_upper", "SELECT", []expect{{kwSELECT, "SELECT"}}},
		{"select_mixed", "SeLeCt", []expect{{kwSELECT, "SeLeCt"}}},
		{"from", "FROM", []expect{{kwFROM, "FROM"}}},
		{"where", "WHERE", []expect{{kwWHERE, "WHERE"}}},
		{"create", "CREATE", []expect{{kwCREATE, "CREATE"}}},
		{"table", "TABLE", []expect{{kwTABLE, "TABLE"}}},
		{"varchar", "VARCHAR", []expect{{kwVARCHAR, "VARCHAR"}}},
		{"variant", "VARIANT", []expect{{kwVARIANT, "VARIANT"}}},
		{"declare", "DECLARE", []expect{{kwDECLARE, "DECLARE"}}},
	}
	runLexerCases(t, cases)
}

func TestLexer_KeywordReservedClassification(t *testing.T) {
	// Spot-check IsReservedKeyword against known reserved/non-reserved entries.
	reserved := []string{"SELECT", "FROM", "WHERE", "CREATE", "DROP", "TABLE", "JOIN", "WITH"}
	for _, kw := range reserved {
		if !IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = false, want true", kw)
		}
	}
	nonReserved := []string{"ALERT", "VARCHAR", "VARIANT", "STAGE", "DECLARE"}
	for _, kw := range nonReserved {
		if IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = true, want false", kw)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	// All three comment forms should produce empty token streams (just EOF).
	cases := []string{
		"-- line comment\n",
		"// line comment\n",
		"/* block comment */",
		"/* nested /* inner */ outer */",
		"/* a */ /* b */",
		"-- comment 1\n-- comment 2\n",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			tokens, errs := Tokenize(input)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %+v", errs)
			}
			if len(tokens) != 1 || tokens[0].Type != tokEOF {
				t.Errorf("expected only EOF token, got %+v", tokens)
			}
		})
	}
}

func TestLexer_CommentsBetweenTokens(t *testing.T) {
	// Comments between two real tokens should be invisible.
	tokens, errs := Tokenize("SELECT /* hidden */ FROM")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens (SELECT FROM EOF), got %d: %+v", len(tokens), tokens)
	}
	if tokens[0].Type != kwSELECT || tokens[1].Type != kwFROM {
		t.Errorf("expected [SELECT, FROM, EOF], got %+v", tokens)
	}
}

func TestLexer_NewLexerWithOffset(t *testing.T) {
	// Tokenize "   SELECT   " with baseOffset=100. The kwSELECT token
	// should report Loc.Start = 103, Loc.End = 109 (positions shifted
	// by 100 relative to the segment's local positions 3..9).
	tokens, errs := TokenizeWithOffset("   SELECT   ", 100)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	// Find the kwSELECT token.
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
	// An unterminated string should produce a LexError whose Loc is
	// shifted by baseOffset.
	_, errs := TokenizeWithOffset("'unterminated", 50)
	if len(errs) != 1 {
		t.Fatalf("expected 1 lex error, got %d: %+v", len(errs), errs)
	}
	// The error's Loc.Start should be 50 (the position of the opening
	// quote in the full document).
	if errs[0].Loc.Start != 50 {
		t.Errorf("lex error Loc.Start = %d, want 50", errs[0].Loc.Start)
	}
}

func TestLexer_NewLexer_UnchangedBehavior(t *testing.T) {
	// Regression check: the existing NewLexer should produce zero-based
	// positions, same as before the baseOffset refactor.
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

// TokenizeWithOffset is a test helper: one-shot tokenize with a base
// offset. Mirrors the existing Tokenize helper but constructs via
// NewLexerWithOffset.
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
