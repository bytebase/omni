package parser

import (
	"testing"
)

// ---- helpers ----------------------------------------------------------------

// tokenize is a thin wrapper that returns only token kinds, dropping EOF.
func tokenizeKinds(t *testing.T, input string) []TokenKind {
	t.Helper()
	tokens, _ := Tokenize(input)
	out := make([]TokenKind, 0, len(tokens))
	for _, tok := range tokens {
		if tok.Kind != tokEOF {
			out = append(out, tok.Kind)
		}
	}
	return out
}

// singleToken lexes input and asserts exactly one non-EOF token is returned,
// returning it for further inspection.
func singleToken(t *testing.T, input string) Token {
	t.Helper()
	tokens, errs := Tokenize(input)
	// expect at most one non-EOF token
	nonEOF := make([]Token, 0, 2)
	for _, tok := range tokens {
		if tok.Kind != tokEOF {
			nonEOF = append(nonEOF, tok)
		}
	}
	if len(nonEOF) != 1 {
		t.Fatalf("singleToken(%q): got %d non-EOF tokens, want 1", input, len(nonEOF))
	}
	if len(errs) > 0 {
		t.Fatalf("singleToken(%q): unexpected errors: %v", input, errs)
	}
	return nonEOF[0]
}

// assertKinds verifies token kinds ignoring EOF.
func assertKinds(t *testing.T, input string, want []TokenKind) {
	t.Helper()
	got := tokenizeKinds(t, input)
	if len(got) != len(want) {
		t.Fatalf("assertKinds(%q): got %d tokens %v, want %d %v", input, len(got), got, len(want), want)
	}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("assertKinds(%q) token[%d]: got kind %d, want %d", input, i, got[i], k)
		}
	}
}

// ---- 1. Keywords ------------------------------------------------------------

func TestKeywords_BasicReserved(t *testing.T) {
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"SELECT", kwSELECT},
		{"select", kwSELECT},
		{"FROM", kwFROM},
		{"from", kwFROM},
		{"WHERE", kwWHERE},
		{"AND", kwAND},
		{"OR", kwOR},
		{"NOT", kwNOT},
		{"CREATE", kwCREATE},
		{"TABLE", kwTABLE},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("keyword %q: got kind %d, want %d", tc.input, tok.Kind, tc.want)
		}
	}
}

func TestKeywords_DorisCaseInsensitive(t *testing.T) {
	// Mix of cases for Doris-specific keywords.
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"DISTRIBUTED", kwDISTRIBUTED},
		{"Distributed", kwDISTRIBUTED},
		{"distributed", kwDISTRIBUTED},
		{"AGGREGATE", kwAGGREGATE},
		{"aggregate", kwAGGREGATE},
		{"BUCKETS", kwBUCKETS},
		{"buckets", kwBUCKETS},
		{"MATERIALIZED", kwMATERIALIZED},
		{"materialized", kwMATERIALIZED},
		{"COMMENT", kwCOMMENT},
		{"comment", kwCOMMENT},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("keyword %q: got kind %d (%s), want %d (%s)",
				tc.input, tok.Kind, TokenName(tok.Kind), tc.want, TokenName(tc.want))
		}
	}
}

func TestKeywords_StrFieldPreservedAsInput(t *testing.T) {
	// Keyword tokens carry the original (mixed-case) text in Str.
	tok := singleToken(t, "Select")
	if tok.Str != "Select" {
		t.Errorf("Str field: got %q, want %q", tok.Str, "Select")
	}
}

// ---- 2. Identifiers ---------------------------------------------------------

func TestIdentifiers_Unquoted(t *testing.T) {
	cases := []string{"my_table", "col1", "_priv", "$alias", "notAKeyword123"}
	for _, input := range cases {
		tok := singleToken(t, input)
		if tok.Kind != tokIdent {
			t.Errorf("ident %q: got kind %d, want tokIdent", input, tok.Kind)
		}
		if tok.Str != input {
			t.Errorf("ident %q: Str = %q, want %q", input, tok.Str, input)
		}
	}
}

func TestIdentifiers_BacktickQuoted(t *testing.T) {
	tok := singleToken(t, "`my column`")
	if tok.Kind != tokQuotedIdent {
		t.Fatalf("backtick ident: got kind %d, want tokQuotedIdent", tok.Kind)
	}
	if tok.Str != "my column" {
		t.Errorf("backtick ident Str: got %q, want %q", tok.Str, "my column")
	}
}

func TestIdentifiers_BacktickDoubledEscape(t *testing.T) {
	// `` inside backtick literal -> single backtick in value
	tok := singleToken(t, "`a``b`")
	if tok.Kind != tokQuotedIdent {
		t.Fatalf("backtick ident: got kind %d", tok.Kind)
	}
	if tok.Str != "a`b" {
		t.Errorf("backtick doubled escape: got %q, want %q", tok.Str, "a`b")
	}
}

func TestIdentifiers_BacktickKeywordName(t *testing.T) {
	// Backtick-quoting a reserved keyword produces tokQuotedIdent, not a keyword.
	tok := singleToken(t, "`select`")
	if tok.Kind != tokQuotedIdent {
		t.Errorf("`select`: got kind %d, want tokQuotedIdent", tok.Kind)
	}
}

func TestIdentifiers_UnterminatedBacktick(t *testing.T) {
	tokens, errs := Tokenize("`unterminated")
	if len(errs) == 0 {
		t.Fatal("unterminated backtick: expected error")
	}
	if errs[0].Msg != errUnterminatedQuoted {
		t.Errorf("error msg: got %q, want %q", errs[0].Msg, errUnterminatedQuoted)
	}
	// Should emit tokInvalid
	found := false
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			found = true
			break
		}
	}
	if !found {
		t.Error("unterminated backtick: expected tokInvalid token")
	}
}

// ---- 3. String literals -----------------------------------------------------

func TestStrings_SingleQuoted(t *testing.T) {
	tok := singleToken(t, `'hello world'`)
	if tok.Kind != tokString {
		t.Fatalf("single-quoted string: got kind %d", tok.Kind)
	}
	if tok.Str != "hello world" {
		t.Errorf("Str: got %q, want %q", tok.Str, "hello world")
	}
}

func TestStrings_DoubleQuoted(t *testing.T) {
	tok := singleToken(t, `"hello world"`)
	if tok.Kind != tokString {
		t.Fatalf("double-quoted string: got kind %d", tok.Kind)
	}
	if tok.Str != "hello world" {
		t.Errorf("Str: got %q, want %q", tok.Str, "hello world")
	}
}

func TestStrings_EscapeSequences(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{`'\n'`, "\n"},
		{`'\t'`, "\t"},
		{`'\r'`, "\r"},
		{`'\0'`, string([]byte{0})},
		{`'\\'`, "\\"},
		{`'\''`, "'"},
		{`'\"'`, `"`},
		{`'\b'`, string([]byte{0x08})},
		{`'\Z'`, string([]byte{0x1A})},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tokString {
			t.Errorf("escape %s: got kind %d, want tokString", tc.input, tok.Kind)
			continue
		}
		if tok.Str != tc.want {
			t.Errorf("escape %s: Str = %q, want %q", tc.input, tok.Str, tc.want)
		}
	}
}

func TestStrings_DoubledQuoteEscape(t *testing.T) {
	// '' inside single-quoted string -> single '
	tok := singleToken(t, `'it''s'`)
	if tok.Kind != tokString {
		t.Fatalf("doubled-quote: kind %d", tok.Kind)
	}
	if tok.Str != "it's" {
		t.Errorf("doubled-quote: Str = %q, want %q", tok.Str, "it's")
	}
}

func TestStrings_NewlineAllowed(t *testing.T) {
	// Doris (unlike Snowflake) allows real newlines inside string literals.
	tok := singleToken(t, "'line1\nline2'")
	if tok.Kind != tokString {
		t.Fatalf("newline in string: kind %d", tok.Kind)
	}
	if tok.Str != "line1\nline2" {
		t.Errorf("newline in string: Str = %q", tok.Str)
	}
}

func TestStrings_UnterminatedSingleQuote(t *testing.T) {
	_, errs := Tokenize("'unterminated")
	if len(errs) == 0 {
		t.Fatal("unterminated single-quote: expected error")
	}
	if errs[0].Msg != errUnterminatedString {
		t.Errorf("error msg: got %q, want %q", errs[0].Msg, errUnterminatedString)
	}
}

func TestStrings_UnterminatedDoubleQuote(t *testing.T) {
	_, errs := Tokenize(`"unterminated`)
	if len(errs) == 0 {
		t.Fatal("unterminated double-quote: expected error")
	}
	if errs[0].Msg != errUnterminatedString {
		t.Errorf("error msg: got %q, want %q", errs[0].Msg, errUnterminatedString)
	}
}

// ---- 4. Numeric literals ----------------------------------------------------

func TestNumbers_Integer(t *testing.T) {
	tok := singleToken(t, "42")
	if tok.Kind != tokInt {
		t.Fatalf("integer: got kind %d", tok.Kind)
	}
	if tok.Ival != 42 {
		t.Errorf("Ival: got %d, want 42", tok.Ival)
	}
}

func TestNumbers_Zero(t *testing.T) {
	tok := singleToken(t, "0")
	if tok.Kind != tokInt || tok.Ival != 0 {
		t.Errorf("zero: kind=%d Ival=%d", tok.Kind, tok.Ival)
	}
}

func TestNumbers_Decimal(t *testing.T) {
	for _, input := range []string{"3.14", ".5", "0.0"} {
		tok := singleToken(t, input)
		if tok.Kind != tokFloat {
			t.Errorf("decimal %q: got kind %d, want tokFloat", input, tok.Kind)
		}
	}
}

func TestNumbers_Exponent(t *testing.T) {
	for _, input := range []string{"1e10", "1E10", "1.5e-10", "2.0E+3"} {
		tok := singleToken(t, input)
		if tok.Kind != tokFloat {
			t.Errorf("exponent %q: got kind %d, want tokFloat", input, tok.Kind)
		}
		if tok.Str != input {
			t.Errorf("exponent %q: Str = %q", input, tok.Str)
		}
	}
}

func TestNumbers_HexPrefix(t *testing.T) {
	tok := singleToken(t, "0xFF")
	if tok.Kind != tokInt {
		t.Fatalf("0xFF: got kind %d, want tokInt", tok.Kind)
	}
	if tok.Ival != 255 {
		t.Errorf("0xFF Ival: got %d, want 255", tok.Ival)
	}
}

func TestNumbers_HexPrefixLower(t *testing.T) {
	tok := singleToken(t, "0xdeadbeef")
	if tok.Kind != tokInt {
		t.Fatalf("0xdeadbeef: got kind %d, want tokInt", tok.Kind)
	}
	if tok.Ival != 0xdeadbeef {
		t.Errorf("0xdeadbeef Ival: got %d, want %d", tok.Ival, int64(0xdeadbeef))
	}
}

func TestNumbers_BinaryPrefix(t *testing.T) {
	tok := singleToken(t, "0b101")
	if tok.Kind != tokInt {
		t.Fatalf("0b101: got kind %d, want tokInt", tok.Kind)
	}
	if tok.Ival != 5 {
		t.Errorf("0b101 Ival: got %d, want 5", tok.Ival)
	}
}

func TestNumbers_OverflowBecomesFloat(t *testing.T) {
	// A number too large for int64 should become tokFloat.
	tok := singleToken(t, "99999999999999999999999999")
	if tok.Kind != tokFloat {
		t.Errorf("overflow: got kind %d, want tokFloat", tok.Kind)
	}
}

// ---- 5. Operators -----------------------------------------------------------

func TestOperators_MultiChar(t *testing.T) {
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"<=", tokLessEq},
		{">=", tokGreaterEq},
		{"<>", tokNotEq},
		{"!=", tokNotEq},
		{"<=>", tokNullSafeEq},
		{"&&", tokLogicalAnd},
		{"||", tokDoublePipes},
		{"<<", tokShiftLeft},
		{">>", tokShiftRight},
		{"->", tokArrow},
		{"@@", tokDoubleAt},
		{":=", tokAssign},
		{"...", tokDotDotDot},
		{"==", int('=')}, // == is normalized to =
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("operator %q: got kind %d (%s), want %d (%s)",
				tc.input, tok.Kind, TokenName(tok.Kind), tc.want, TokenName(tc.want))
		}
	}
}

func TestOperators_SingleChar(t *testing.T) {
	singles := []byte{'+', '-', '*', '/', '%', '~', '^', '(', ')', ',', ';', '&', '|', '<', '>', '=', '.', '@', ':', '!'}
	for _, ch := range singles {
		input := string(ch)
		tok := singleToken(t, input)
		if tok.Kind != int(ch) {
			t.Errorf("single-char %q: got kind %d, want %d", input, tok.Kind, int(ch))
		}
	}
}

func TestOperators_DorisSpecialAlternatives(t *testing.T) {
	// !< maps to >= and !> maps to <=
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"!<", tokGreaterEq},
		{"!>", tokLessEq},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("operator %q: got kind %d, want %d", tc.input, tok.Kind, tc.want)
		}
	}
}

// ---- 6. Comments ------------------------------------------------------------

func TestComments_DashDashSpace(t *testing.T) {
	// "-- comment" is stripped; remaining tokens are just the integer.
	assertKinds(t, "42 -- this is a comment", []TokenKind{tokInt})
}

func TestComments_DashDashTab(t *testing.T) {
	assertKinds(t, "1 --\tcomment", []TokenKind{tokInt})
}

func TestComments_DashDashNewline(t *testing.T) {
	assertKinds(t, "1 --\nstill", []TokenKind{tokInt, tokIdent})
}

func TestComments_DashDashEOF(t *testing.T) {
	// -- at end of input with nothing after: treated as line comment.
	assertKinds(t, "1 --", []TokenKind{tokInt})
}

func TestComments_DoubleSlash(t *testing.T) {
	assertKinds(t, "1 // comment\n2", []TokenKind{tokInt, tokInt})
}

func TestComments_Hash(t *testing.T) {
	assertKinds(t, "1 # comment\n2", []TokenKind{tokInt, tokInt})
}

func TestComments_BlockComment(t *testing.T) {
	assertKinds(t, "1 /* block */ 2", []TokenKind{tokInt, tokInt})
}

func TestComments_BlockCommentMultiline(t *testing.T) {
	assertKinds(t, "1 /* line1\nline2 */ 2", []TokenKind{tokInt, tokInt})
}

func TestComments_UnterminatedBlock(t *testing.T) {
	_, errs := Tokenize("/* unterminated")
	if len(errs) == 0 {
		t.Fatal("unterminated block comment: expected error")
	}
	if errs[0].Msg != errUnterminatedComment {
		t.Errorf("error msg: got %q, want %q", errs[0].Msg, errUnterminatedComment)
	}
}

func TestComments_HintNotConsumedAsComment(t *testing.T) {
	// /*+ must NOT be skipped as a block comment.
	kinds := tokenizeKinds(t, "/*+ hint */")
	if len(kinds) == 0 {
		t.Fatal("hint: expected tokens, got none")
	}
	if kinds[0] != tokHintStart {
		t.Errorf("hint: first token kind %d, want tokHintStart (%d)", kinds[0], tokHintStart)
	}
}

// ---- 7. Hints ---------------------------------------------------------------

func TestHints_StartAndEnd(t *testing.T) {
	assertKinds(t, "/*+ USE_HASH(t) */", []TokenKind{
		tokHintStart,
		tokIdent,  // USE_HASH (or keyword if recognized)
		int('('),
		tokIdent, // t
		int(')'),
		tokHintEnd,
	})
}

func TestHints_StartToken(t *testing.T) {
	tokens, _ := Tokenize("/*+")
	if tokens[0].Kind != tokHintStart {
		t.Errorf("/*+: first token kind %d, want tokHintStart", tokens[0].Kind)
	}
}

// ---- 8. Hex/Bit literals ----------------------------------------------------

func TestHexLiteral_Uppercase(t *testing.T) {
	tok := singleToken(t, "X'FF'")
	if tok.Kind != tokHexLiteral {
		t.Fatalf("X'FF': kind %d, want tokHexLiteral", tok.Kind)
	}
	if tok.Str != "FF" {
		t.Errorf("X'FF' Str: got %q, want %q", tok.Str, "FF")
	}
}

func TestHexLiteral_Lowercase(t *testing.T) {
	tok := singleToken(t, "x'ff'")
	if tok.Kind != tokHexLiteral {
		t.Fatalf("x'ff': kind %d, want tokHexLiteral", tok.Kind)
	}
	if tok.Str != "ff" {
		t.Errorf("x'ff' Str: got %q, want %q", tok.Str, "ff")
	}
}

func TestBitLiteral_Uppercase(t *testing.T) {
	tok := singleToken(t, "B'101'")
	if tok.Kind != tokBitLiteral {
		t.Fatalf("B'101': kind %d, want tokBitLiteral", tok.Kind)
	}
	if tok.Str != "101" {
		t.Errorf("B'101' Str: got %q, want %q", tok.Str, "101")
	}
}

func TestBitLiteral_Lowercase(t *testing.T) {
	tok := singleToken(t, "b'0'")
	if tok.Kind != tokBitLiteral {
		t.Fatalf("b'0': kind %d", tok.Kind)
	}
	if tok.Str != "0" {
		t.Errorf("b'0' Str: got %q, want %q", tok.Str, "0")
	}
}

func TestHexLiteral_DistinctFromHexPrefix(t *testing.T) {
	// X'FF' -> tokHexLiteral, 0xFF -> tokInt
	tokX := singleToken(t, "X'FF'")
	tok0x := singleToken(t, "0xFF")
	if tokX.Kind != tokHexLiteral {
		t.Errorf("X'FF': want tokHexLiteral, got %d", tokX.Kind)
	}
	if tok0x.Kind != tokInt {
		t.Errorf("0xFF: want tokInt, got %d", tok0x.Kind)
	}
}

// ---- 9. Error recovery ------------------------------------------------------

func TestErrors_UnterminatedString(t *testing.T) {
	tokens, errs := Tokenize("'oops")
	if len(errs) == 0 {
		t.Fatal("unterminated string: expected error")
	}
	if errs[0].Msg != errUnterminatedString {
		t.Errorf("error msg: got %q", errs[0].Msg)
	}
	// tokInvalid must be emitted
	found := false
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			found = true
		}
	}
	if !found {
		t.Error("expected tokInvalid token for unterminated string")
	}
}

func TestErrors_UnterminatedBlockComment(t *testing.T) {
	_, errs := Tokenize("SELECT /* no end")
	if len(errs) == 0 {
		t.Fatal("unterminated block comment: expected error")
	}
	if errs[0].Msg != errUnterminatedComment {
		t.Errorf("error msg: got %q", errs[0].Msg)
	}
}

func TestErrors_UnknownCharacter(t *testing.T) {
	tokens, errs := Tokenize("\x01")
	if len(errs) == 0 {
		t.Fatal("unknown char: expected error")
	}
	if errs[0].Msg != errUnknownChar {
		t.Errorf("error msg: got %q", errs[0].Msg)
	}
	found := false
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			found = true
		}
	}
	if !found {
		t.Error("expected tokInvalid token for unknown char")
	}
}

func TestErrors_RecoveryAfterError(t *testing.T) {
	// After an unterminated string, the lexer should still emit subsequent tokens.
	tokens, errs := Tokenize("'oops 42")
	if len(errs) == 0 {
		t.Fatal("expected lex error")
	}
	// The tokInvalid should be present.
	foundInvalid := false
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			foundInvalid = true
		}
	}
	if !foundInvalid {
		t.Error("expected tokInvalid in error recovery")
	}
}

// ---- 10. Position tracking --------------------------------------------------

func TestPosition_BasicOffsets(t *testing.T) {
	// "SELECT" starts at 0, ends at 6.
	tokens, _ := Tokenize("SELECT")
	if len(tokens) == 0 || tokens[0].Kind == tokEOF {
		t.Fatal("no tokens")
	}
	tok := tokens[0]
	if tok.Loc.Start != 0 {
		t.Errorf("Loc.Start: got %d, want 0", tok.Loc.Start)
	}
	if tok.Loc.End != 6 {
		t.Errorf("Loc.End: got %d, want 6", tok.Loc.End)
	}
}

func TestPosition_WithLeadingWhitespace(t *testing.T) {
	// "   42" - integer starts at byte 3.
	tokens, _ := Tokenize("   42")
	nonEOF := filterNonEOF(tokens)
	if len(nonEOF) == 0 {
		t.Fatal("no tokens")
	}
	tok := nonEOF[0]
	if tok.Loc.Start != 3 {
		t.Errorf("Loc.Start: got %d, want 3", tok.Loc.Start)
	}
	if tok.Loc.End != 5 {
		t.Errorf("Loc.End: got %d, want 5", tok.Loc.End)
	}
}

func TestPosition_MultipleTokens(t *testing.T) {
	// "a b" — ident 'a' at [0,1), space, ident 'b' at [2,3).
	tokens, _ := Tokenize("a b")
	nonEOF := filterNonEOF(tokens)
	if len(nonEOF) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(nonEOF))
	}
	if nonEOF[0].Loc.Start != 0 || nonEOF[0].Loc.End != 1 {
		t.Errorf("first token loc: [%d,%d), want [0,1)", nonEOF[0].Loc.Start, nonEOF[0].Loc.End)
	}
	if nonEOF[1].Loc.Start != 2 || nonEOF[1].Loc.End != 3 {
		t.Errorf("second token loc: [%d,%d), want [2,3)", nonEOF[1].Loc.Start, nonEOF[1].Loc.End)
	}
}

func TestPosition_BaseOffset(t *testing.T) {
	// NewLexerWithOffset shifts all emitted Locs by baseOffset.
	const base = 100
	l := NewLexerWithOffset("SELECT", base)
	tok := l.NextToken()
	if tok.Kind == tokEOF {
		t.Fatal("expected non-EOF token")
	}
	if tok.Loc.Start != base {
		t.Errorf("shifted Loc.Start: got %d, want %d", tok.Loc.Start, base)
	}
	if tok.Loc.End != base+6 {
		t.Errorf("shifted Loc.End: got %d, want %d", tok.Loc.End, base+6)
	}
}

func TestPosition_BaseOffsetError(t *testing.T) {
	// Errors() should also shift positions by baseOffset.
	const base = 50
	l := NewLexerWithOffset("'unterminated", base)
	for {
		tok := l.NextToken()
		if tok.Kind == tokEOF {
			break
		}
	}
	errs := l.Errors()
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
	if errs[0].Loc.Start != base {
		t.Errorf("shifted error Loc.Start: got %d, want %d", errs[0].Loc.Start, base)
	}
}

func filterNonEOF(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens))
	for _, tok := range tokens {
		if tok.Kind != tokEOF {
			out = append(out, tok)
		}
	}
	return out
}

// ---- 11. Edge cases ---------------------------------------------------------

func TestEdge_EmptyInput(t *testing.T) {
	tokens, errs := Tokenize("")
	if len(errs) != 0 {
		t.Errorf("empty input: unexpected errors %v", errs)
	}
	if len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Errorf("empty input: expected exactly EOF, got %v", tokens)
	}
}

func TestEdge_WhitespaceOnly(t *testing.T) {
	tokens, errs := Tokenize("   \t\n\r  ")
	if len(errs) != 0 {
		t.Errorf("whitespace-only: unexpected errors")
	}
	if len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Errorf("whitespace-only: expected exactly EOF, got %v", tokens)
	}
}

func TestEdge_MultiStatement(t *testing.T) {
	// SELECT 1; SELECT 2 — should produce two SELECT keywords, two ints, one semicolon.
	kinds := tokenizeKinds(t, "SELECT 1; SELECT 2")
	want := []TokenKind{kwSELECT, tokInt, int(';'), kwSELECT, tokInt}
	if len(kinds) != len(want) {
		t.Fatalf("multi-statement: got %v, want %v", kinds, want)
	}
	for i, k := range want {
		if kinds[i] != k {
			t.Errorf("multi-statement token[%d]: got %d, want %d", i, kinds[i], k)
		}
	}
}

func TestEdge_Placeholder(t *testing.T) {
	tok := singleToken(t, "?")
	if tok.Kind != tokPlaceholder {
		t.Errorf("placeholder: got kind %d, want tokPlaceholder", tok.Kind)
	}
}

func TestEdge_PlaceholderInExpression(t *testing.T) {
	assertKinds(t, "WHERE x = ?", []TokenKind{kwWHERE, tokIdent, int('='), tokPlaceholder})
}

// ---- 12. Keyword classification ---------------------------------------------

func TestKeywordClass_SelectIsReserved(t *testing.T) {
	if !IsReserved(kwSELECT) {
		t.Error("SELECT should be reserved")
	}
}

func TestKeywordClass_FromIsReserved(t *testing.T) {
	if !IsReserved(kwFROM) {
		t.Error("FROM should be reserved")
	}
}

func TestKeywordClass_CommentIsNonReserved(t *testing.T) {
	if IsReserved(kwCOMMENT) {
		t.Error("COMMENT should be non-reserved")
	}
}

func TestKeywordClass_MaterializedIsNonReserved(t *testing.T) {
	if IsReserved(kwMATERIALIZED) {
		t.Error("MATERIALIZED should be non-reserved")
	}
}

func TestKeywordClass_BucketsIsNonReserved(t *testing.T) {
	if IsReserved(kwBUCKETS) {
		t.Error("BUCKETS should be non-reserved")
	}
}

func TestKeywordClass_AggregateIsNonReserved(t *testing.T) {
	if IsReserved(kwAGGREGATE) {
		t.Error("AGGREGATE should be non-reserved")
	}
}

func TestKeywordClass_LiteralsAreNotReserved(t *testing.T) {
	for _, kind := range []TokenKind{tokInt, tokFloat, tokString, tokIdent, tokPlaceholder} {
		if IsReserved(kind) {
			t.Errorf("literal token %d (%s) should not be reserved", kind, TokenName(kind))
		}
	}
}

func TestTokenName_Keywords(t *testing.T) {
	cases := []struct {
		kind TokenKind
		want string
	}{
		{kwSELECT, "SELECT"},
		{kwFROM, "FROM"},
		{kwCOMMENT, "COMMENT"},
		{kwDISTRIBUTED, "DISTRIBUTED"},
		{kwBUCKETS, "BUCKETS"},
	}
	for _, tc := range cases {
		got := TokenName(tc.kind)
		if got != tc.want {
			t.Errorf("TokenName(%d): got %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestTokenName_SpecialTokens(t *testing.T) {
	cases := []struct {
		kind TokenKind
		want string
	}{
		{tokEOF, "EOF"},
		{tokInvalid, "INVALID"},
		{tokInt, "INT"},
		{tokFloat, "FLOAT"},
		{tokString, "STRING"},
		{tokIdent, "IDENT"},
		{tokQuotedIdent, "QUOTED_IDENT"},
		{tokHexLiteral, "HEX_LITERAL"},
		{tokBitLiteral, "BIT_LITERAL"},
		{tokPlaceholder, "PLACEHOLDER"},
	}
	for _, tc := range cases {
		got := TokenName(tc.kind)
		if got != tc.want {
			t.Errorf("TokenName(%d): got %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestTokenName_Operators(t *testing.T) {
	cases := []struct {
		kind TokenKind
		want string
	}{
		{tokLessEq, "<="},
		{tokGreaterEq, ">="},
		{tokNotEq, "<>"},
		{tokNullSafeEq, "<=>"},
		{tokLogicalAnd, "&&"},
		{tokDoublePipes, "||"},
		{tokShiftLeft, "<<"},
		{tokShiftRight, ">>"},
		{tokArrow, "->"},
		{tokDoubleAt, "@@"},
		{tokAssign, ":="},
		{tokDotDotDot, "..."},
		{tokHintStart, "HINT_START"},
		{tokHintEnd, "HINT_END"},
	}
	for _, tc := range cases {
		got := TokenName(tc.kind)
		if got != tc.want {
			t.Errorf("TokenName(%d): got %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestTokenName_ASCIISingleChar(t *testing.T) {
	// Single ASCII printable chars map to their string form.
	for _, ch := range []byte{'+', '-', '*', '/', '(', ')', ';', ','} {
		got := TokenName(int(ch))
		want := string(ch)
		if got != want {
			t.Errorf("TokenName(%d): got %q, want %q", ch, got, want)
		}
	}
}

func TestKeywordToken_CaseInsensitive(t *testing.T) {
	cases := []string{"select", "SELECT", "Select", "SeLeCt"}
	for _, s := range cases {
		kind, ok := KeywordToken(s)
		if !ok {
			t.Errorf("KeywordToken(%q): not found", s)
			continue
		}
		if kind != kwSELECT {
			t.Errorf("KeywordToken(%q): got %d, want kwSELECT", s, kind)
		}
	}
}

func TestKeywordToken_NonKeyword(t *testing.T) {
	_, ok := KeywordToken("notakeyword")
	if ok {
		t.Error("KeywordToken(\"notakeyword\"): expected not found")
	}
}

// ---- 13. MySQL-strict dash-dash --------------------------------------------

func TestDashDash_NoSpaceShouldNotBeComment(t *testing.T) {
	// "--5" is NOT a comment (no space/tab/newline after --), so it produces
	// two minus tokens followed by an integer, not a skipped comment.
	kinds := tokenizeKinds(t, "--5")
	if len(kinds) == 0 {
		t.Fatal("--5: expected tokens, got none (should not be treated as comment)")
	}
	// The integer 5 must appear somewhere in the token stream.
	foundInt := false
	for _, k := range kinds {
		if k == tokInt {
			foundInt = true
		}
	}
	if !foundInt {
		t.Errorf("--5: expected tokInt in output, got %v", kinds)
	}
}

func TestDashDash_WithSpaceIsComment(t *testing.T) {
	// "-- comment\n5" skips the comment; only the integer 5 should remain.
	kinds := tokenizeKinds(t, "-- comment\n5")
	if len(kinds) != 1 || kinds[0] != tokInt {
		t.Errorf("-- comment: expected [tokInt], got %v", kinds)
	}
}

func TestDashDash_AtEOFIsComment(t *testing.T) {
	// "5 --" at EOF: the -- is a valid line comment (EOF counts).
	kinds := tokenizeKinds(t, "5 --")
	if len(kinds) != 1 || kinds[0] != tokInt {
		t.Errorf("5 --: expected [tokInt], got %v", kinds)
	}
}

// ---- 14. Tokenize convenience function --------------------------------------

func TestTokenize_ReturnsErrorsAndTokens(t *testing.T) {
	tokens, errs := Tokenize("SELECT 'unterminated")
	if len(errs) == 0 {
		t.Fatal("expected lex errors")
	}
	// SELECT keyword should still appear.
	found := false
	for _, tok := range tokens {
		if tok.Kind == kwSELECT {
			found = true
		}
	}
	if !found {
		t.Error("Tokenize: SELECT keyword missing despite error after it")
	}
}

func TestTokenize_AlwaysEndsWithEOF(t *testing.T) {
	for _, input := range []string{"", "SELECT 1", "   ", "'oops"} {
		tokens, _ := Tokenize(input)
		if len(tokens) == 0 {
			t.Errorf("input %q: no tokens at all", input)
			continue
		}
		last := tokens[len(tokens)-1]
		if last.Kind != tokEOF {
			t.Errorf("input %q: last token is %d, want tokEOF", input, last.Kind)
		}
	}
}
