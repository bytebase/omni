package parser

import "testing"

// ---- helpers ----------------------------------------------------------------

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

func filterNonEOF(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens))
	for _, tok := range tokens {
		if tok.Kind != tokEOF {
			out = append(out, tok)
		}
	}
	return out
}

// singleToken lexes input and asserts exactly one non-EOF token with no errors.
func singleToken(t *testing.T, input string) Token {
	t.Helper()
	tokens, errs := Tokenize(input)
	nonEOF := filterNonEOF(tokens)
	if len(nonEOF) != 1 {
		t.Fatalf("singleToken(%q): got %d non-EOF tokens %v, want 1", input, len(nonEOF), nonEOF)
	}
	if len(errs) > 0 {
		t.Fatalf("singleToken(%q): unexpected errors: %v", input, errs)
	}
	return nonEOF[0]
}

func assertKinds(t *testing.T, input string, want []TokenKind) {
	t.Helper()
	got := tokenizeKinds(t, input)
	if len(got) != len(want) {
		t.Fatalf("assertKinds(%q): got %d tokens %v, want %d %v", input, len(got), got, len(want), want)
	}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("assertKinds(%q) token[%d]: got %s, want %s", input, i, TokenName(got[i]), TokenName(k))
		}
	}
}

func hasInvalid(tokens []Token) bool {
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			return true
		}
	}
	return false
}

// ---- 1. Keywords ------------------------------------------------------------

func TestKeywords_BasicReserved(t *testing.T) {
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"SELECT", kwSELECT}, {"select", kwSELECT}, {"SeLeCt", kwSELECT},
		{"FROM", kwFROM}, {"WHERE", kwWHERE}, {"AND", kwAND}, {"OR", kwOR},
		{"NOT", kwNOT}, {"CREATE", kwCREATE}, {"TABLE", kwTABLE},
		{"JOIN", kwJOIN}, {"UNION", kwUNION}, {"INTERSECT", kwINTERSECT},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("keyword %q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
	}
}

func TestKeywords_NonReserved(t *testing.T) {
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"ARRAY", kwARRAY}, {"MAP", kwMAP}, {"SHOW", kwSHOW}, {"GRANT", kwGRANT},
		{"MERGE", kwMERGE}, {"MATERIALIZED", kwMATERIALIZED}, {"SESSION", kwSESSION},
		{"json", kwJSON}, {"comment", kwCOMMENT},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("keyword %q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
	}
}

func TestKeywords_Compound(t *testing.T) {
	// Compound keywords are single lexer tokens (one IDENTIFIER_-shaped run with
	// embedded underscores), not multiple tokens.
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"CURRENT_DATE", kwCURRENT_DATE},
		{"CURRENT_TIMESTAMP", kwCURRENT_TIMESTAMP},
		{"CURRENT_USER", kwCURRENT_USER},
		{"MATCH_RECOGNIZE", kwMATCH_RECOGNIZE},
		{"JSON_OBJECT", kwJSON_OBJECT},
		{"TRY_CAST", kwTRY_CAST},
		{"LOCALTIMESTAMP", kwLOCALTIMESTAMP},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("compound keyword %q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
	}
}

func TestKeywords_TextStringAlias(t *testing.T) {
	// The literal text STRING is Trino's STRING type alias (legacy TEXT_STRING_).
	tok := singleToken(t, "STRING")
	if tok.Kind != kwSTRING {
		t.Errorf("STRING: got %s, want kwSTRING", TokenName(tok.Kind))
	}
	if IsReserved(kwSTRING) {
		t.Error("STRING should be non-reserved")
	}
}

func TestKeywords_StrPreservesSourceCase(t *testing.T) {
	tok := singleToken(t, "Select")
	if tok.Str != "Select" {
		t.Errorf("keyword Str: got %q, want %q", tok.Str, "Select")
	}
}

// ---- 2. Identifiers ---------------------------------------------------------

func TestIdent_Unquoted(t *testing.T) {
	for _, input := range []string{"my_table", "col1", "_priv", "notAKeyword123", "_"} {
		tok := singleToken(t, input)
		if tok.Kind != tokIdent {
			t.Errorf("ident %q: got %s, want IDENTIFIER", input, TokenName(tok.Kind))
		}
		if tok.Str != input {
			t.Errorf("ident %q: Str = %q, want %q", input, tok.Str, input)
		}
	}
}

func TestIdent_QuotedDoubleQuote(t *testing.T) {
	tok := singleToken(t, `"my column"`)
	if tok.Kind != tokQuotedIdent {
		t.Fatalf(`"my column": got %s, want QUOTED_IDENTIFIER`, TokenName(tok.Kind))
	}
	if tok.Str != "my column" {
		t.Errorf("quoted ident Str: got %q, want %q", tok.Str, "my column")
	}
}

func TestIdent_QuotedDoubledEscape(t *testing.T) {
	// "" inside a quoted identifier -> single ".
	tok := singleToken(t, `"a""b"`)
	if tok.Kind != tokQuotedIdent || tok.Str != `a"b` {
		t.Errorf(`"a""b": kind=%s Str=%q, want QUOTED_IDENTIFIER a"b`, TokenName(tok.Kind), tok.Str)
	}
}

func TestIdent_Backquoted(t *testing.T) {
	tok := singleToken(t, "`my col`")
	if tok.Kind != tokBackquotedIdent || tok.Str != "my col" {
		t.Errorf("backquoted ident: kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestIdent_BackquotedDoubledEscape(t *testing.T) {
	tok := singleToken(t, "`a``b`")
	if tok.Kind != tokBackquotedIdent || tok.Str != "a`b" {
		t.Errorf("backquoted doubled: kind=%s Str=%q, want a`b", TokenName(tok.Kind), tok.Str)
	}
}

func TestIdent_QuotedKeywordIsIdentifier(t *testing.T) {
	// Quoting a reserved keyword yields an identifier token, not the keyword.
	tok := singleToken(t, `"select"`)
	if tok.Kind != tokQuotedIdent {
		t.Errorf(`"select": got %s, want QUOTED_IDENTIFIER`, TokenName(tok.Kind))
	}
}

func TestIdent_DigitLeading(t *testing.T) {
	// DIGIT_IDENTIFIER_: a digit-led run containing a letter or underscore.
	cases := []string{"1a", "12ab", "1_000", "9x9", "0a"}
	for _, input := range cases {
		tok := singleToken(t, input)
		if tok.Kind != tokDigitIdent {
			t.Errorf("digit-ident %q: got %s, want DIGIT_IDENTIFIER", input, TokenName(tok.Kind))
		}
		if tok.Str != input {
			t.Errorf("digit-ident %q: Str = %q", input, tok.Str)
		}
	}
}

func TestIdent_NonASCIIUnquotedRejected(t *testing.T) {
	// TrinoLexer.g4's LETTER_ = [A-Z] (ASCII only). An unquoted non-ASCII
	// identifier byte is routed to UNRECOGNIZED_ in both Trino and the legacy
	// grammar, so omni must flag it (tokInvalid) rather than accept it as an
	// identifier. The same character double-quoted IS a valid identifier.
	tokens, errs := Tokenize("é") // U+00E9, two UTF-8 bytes 0xC3 0xA9
	if len(errs) == 0 || errs[0].Msg != errUnknownChar {
		t.Errorf("unquoted é: expected errUnknownChar, got errs=%v", errs)
	}
	if !hasInvalid(tokens) {
		t.Error("unquoted é: expected tokInvalid")
	}
	// Quoted, it is a valid identifier carrying the UTF-8 content verbatim.
	q := singleToken(t, `"é"`)
	if q.Kind != tokQuotedIdent || q.Str != "é" {
		t.Errorf(`"é": kind=%s Str=%q, want QUOTED_IDENTIFIER é`, TokenName(q.Kind), q.Str)
	}
}

func TestIdent_NonASCIIAfterDigitRejected(t *testing.T) {
	// '1é' — the digit is fine, but the non-ASCII byte is not an identifier
	// continuation, so the run is INTEGER 1 followed by an unrecognized byte.
	tokens, errs := Tokenize("1é")
	if len(errs) == 0 || errs[0].Msg != errUnknownChar {
		t.Errorf("1é: expected errUnknownChar, got %v", errs)
	}
	if tokens[0].Kind != tokInteger || tokens[0].Ival != 1 {
		t.Errorf("1é: first token should be INTEGER 1, got %s", TokenName(tokens[0].Kind))
	}
}

func TestIdent_UnterminatedQuoted(t *testing.T) {
	for _, input := range []string{`"unterminated`, "`unterminated"} {
		tokens, errs := Tokenize(input)
		if len(errs) == 0 || errs[0].Msg != errUnterminatedQuoted {
			t.Errorf("%q: expected errUnterminatedQuoted, got %v", input, errs)
		}
		if !hasInvalid(tokens) {
			t.Errorf("%q: expected tokInvalid", input)
		}
	}
}

// ---- 3. String / Unicode / Binary literals ----------------------------------

func TestString_Basic(t *testing.T) {
	tok := singleToken(t, `'hello world'`)
	if tok.Kind != tokString || tok.Str != "hello world" {
		t.Errorf("string: kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestString_DoubledQuoteEscape(t *testing.T) {
	tok := singleToken(t, `'it''s'`)
	if tok.Kind != tokString || tok.Str != "it's" {
		t.Errorf("doubled-quote string: kind=%s Str=%q, want it's", TokenName(tok.Kind), tok.Str)
	}
}

func TestString_BackslashIsLiteral(t *testing.T) {
	// Trino is standard-conforming: backslash is an ordinary character, NOT an
	// escape (unlike MySQL/Doris). '\n' is backslash + n, two bytes.
	tok := singleToken(t, `'\n'`)
	if tok.Kind != tokString {
		t.Fatalf(`'\n': got %s`, TokenName(tok.Kind))
	}
	if tok.Str != `\n` {
		t.Errorf(`'\n': Str = %q, want literal backslash-n`, tok.Str)
	}
}

func TestString_NewlineAllowed(t *testing.T) {
	tok := singleToken(t, "'line1\nline2'")
	if tok.Kind != tokString || tok.Str != "line1\nline2" {
		t.Errorf("newline string: kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestString_Unterminated(t *testing.T) {
	tokens, errs := Tokenize("'oops")
	if len(errs) == 0 || errs[0].Msg != errUnterminatedString {
		t.Errorf("unterminated string: got errs %v", errs)
	}
	if !hasInvalid(tokens) {
		t.Error("unterminated string: expected tokInvalid")
	}
}

func TestUnicodeString_Basic(t *testing.T) {
	tok := singleToken(t, `U&'\0041'`)
	if tok.Kind != tokUnicodeString {
		t.Fatalf("U&'...': got %s, want UNICODE_STRING", TokenName(tok.Kind))
	}
	// The lexer captures raw inter-quote bytes; escape decoding is the parser's job.
	if tok.Str != `\0041` {
		t.Errorf("unicode string Str: got %q, want raw \\0041", tok.Str)
	}
}

func TestUnicodeString_LowercaseU(t *testing.T) {
	tok := singleToken(t, `u&'abc'`)
	if tok.Kind != tokUnicodeString || tok.Str != "abc" {
		t.Errorf("u&'abc': kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestUnicodeString_NotConfusedWithIdent(t *testing.T) {
	// "U" alone, or "U&" without a quote, is NOT a unicode-string start.
	// "u" is an identifier.
	tok := singleToken(t, "u")
	if tok.Kind != tokIdent {
		t.Errorf("bare u: got %s, want IDENTIFIER", TokenName(tok.Kind))
	}
}

func TestBinaryLiteral_Basic(t *testing.T) {
	tok := singleToken(t, "X'00ff'")
	if tok.Kind != tokBinaryLiteral || tok.Str != "00ff" {
		t.Errorf("X'00ff': kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestBinaryLiteral_LowercaseX(t *testing.T) {
	tok := singleToken(t, "x'ABCD'")
	if tok.Kind != tokBinaryLiteral || tok.Str != "ABCD" {
		t.Errorf("x'ABCD': kind=%s Str=%q", TokenName(tok.Kind), tok.Str)
	}
}

func TestBinaryLiteral_Unterminated(t *testing.T) {
	tokens, errs := Tokenize("X'00")
	if len(errs) == 0 || errs[0].Msg != errUnterminatedBinary {
		t.Errorf("X'00: got errs %v", errs)
	}
	if !hasInvalid(tokens) {
		t.Error("X'00: expected tokInvalid")
	}
}

func TestBinaryLiteral_DistinctFromHexLikeIdent(t *testing.T) {
	// X'FF' is a binary literal; XFF (no quote) is an identifier.
	if tok := singleToken(t, "X'FF'"); tok.Kind != tokBinaryLiteral {
		t.Errorf("X'FF': want BINARY_LITERAL, got %s", TokenName(tok.Kind))
	}
	if tok := singleToken(t, "XFF"); tok.Kind != tokIdent {
		t.Errorf("XFF: want IDENTIFIER, got %s", TokenName(tok.Kind))
	}
}

// ---- 4. Numbers -------------------------------------------------------------

func TestNumber_Integer(t *testing.T) {
	tok := singleToken(t, "42")
	if tok.Kind != tokInteger || tok.Ival != 42 || tok.Str != "42" {
		t.Errorf("42: kind=%s Ival=%d Str=%q", TokenName(tok.Kind), tok.Ival, tok.Str)
	}
}

func TestNumber_Zero(t *testing.T) {
	tok := singleToken(t, "0")
	if tok.Kind != tokInteger || tok.Ival != 0 {
		t.Errorf("0: kind=%s Ival=%d", TokenName(tok.Kind), tok.Ival)
	}
}

func TestNumber_Decimal(t *testing.T) {
	for _, input := range []string{"3.14", ".5", "0.0", "12.", "100.00"} {
		tok := singleToken(t, input)
		if tok.Kind != tokDecimal {
			t.Errorf("decimal %q: got %s, want DECIMAL_VALUE", input, TokenName(tok.Kind))
		}
		if tok.Str != input {
			t.Errorf("decimal %q: Str=%q", input, tok.Str)
		}
	}
}

func TestNumber_Double(t *testing.T) {
	for _, input := range []string{"1e10", "1E10", "1.5e-10", "2.0E+3", ".5e3", "1e+5"} {
		tok := singleToken(t, input)
		if tok.Kind != tokDouble {
			t.Errorf("double %q: got %s, want DOUBLE_VALUE", input, TokenName(tok.Kind))
		}
		if tok.Str != input {
			t.Errorf("double %q: Str=%q", input, tok.Str)
		}
	}
}

func TestNumber_NoSignConsumed(t *testing.T) {
	// Trino's number rule does not lex a leading sign; -5 is MINUS_ then 5.
	assertKinds(t, "-5", []TokenKind{int('-'), tokInteger})
	assertKinds(t, "+5", []TokenKind{int('+'), tokInteger})
}

func TestNumber_DigitIdentMaximalMunch(t *testing.T) {
	// ANTLR maximal-munch + rule-order: pure digits = INTEGER; a trailing letter
	// extends the match into DIGIT_IDENTIFIER; an incomplete exponent ('1e')
	// likewise becomes DIGIT_IDENTIFIER because INTEGER only matches '1'.
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"123", tokInteger},
		{"123a", tokDigitIdent},
		{"1e", tokDigitIdent},   // not a valid DOUBLE; '1e' is a longer DIGIT_IDENTIFIER match than INTEGER '1'
		{"1e10", tokDouble},     // valid exponent
		{"1.5e3", tokDouble},    // DIGIT_IDENTIFIER cannot cross the '.', so DOUBLE is longer
		{"0x1f", tokDigitIdent}, // Trino has no 0x hex numeric form; '0x1f' is a digit-leading identifier
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("%q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
	}
}

func TestNumber_Overflow(t *testing.T) {
	// A pure-digit run too large for int64 is still INTEGER_VALUE to the grammar.
	tok := singleToken(t, "99999999999999999999999999")
	if tok.Kind != tokInteger {
		t.Errorf("overflow: got %s, want INTEGER_VALUE", TokenName(tok.Kind))
	}
}

// ---- 5. Operators / punctuation ---------------------------------------------

func TestOperators_MultiChar(t *testing.T) {
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"<>", tokNotEq}, {"!=", tokNotEq},
		{"<=", tokLessEq}, {">=", tokGreaterEq},
		{"||", tokConcat},
		{"->", tokArrow}, {"<-", tokLeftArrow}, {"=>", tokDoubleArrow},
		{"{-", tokLCurlyHyphen}, {"-}", tokRCurlyHyphen},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("operator %q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
	}
}

func TestOperators_SingleChar(t *testing.T) {
	// Every single-byte operator/punctuation token in TrinoLexer.g4, including
	// COLON_ which the divergence ledger fixes from the bug literal '_:' to ':'.
	singles := []byte{'=', '<', '>', '+', '-', '*', '/', '%', ';', '.', ':', ',',
		'(', ')', '[', ']', '{', '}', '|', '$', '^'}
	for _, ch := range singles {
		tok := singleToken(t, string(ch))
		if tok.Kind != int(ch) {
			t.Errorf("single-char %q: got %s, want %q", string(ch), TokenName(tok.Kind), string(ch))
		}
	}
}

func TestOperators_Colon(t *testing.T) {
	// Divergence ledger item: COLON_ is ':' (legacy grammar's '_:' is a bug).
	tok := singleToken(t, ":")
	if tok.Kind != int(':') {
		t.Errorf("colon: got %s, want ':'", TokenName(tok.Kind))
	}
}

func TestOperators_Question(t *testing.T) {
	tok := singleToken(t, "?")
	if tok.Kind != tokQuestion {
		t.Errorf("?: got %s, want QUESTION_MARK", TokenName(tok.Kind))
	}
}

func TestOperators_DisambiguateMinus(t *testing.T) {
	// '-' alone vs '->' vs '-}'.
	assertKinds(t, "-", []TokenKind{int('-')})
	assertKinds(t, "->", []TokenKind{tokArrow})
	assertKinds(t, "-}", []TokenKind{tokRCurlyHyphen})
	assertKinds(t, "a - b", []TokenKind{tokIdent, int('-'), tokIdent})
}

func TestOperators_DisambiguateCurly(t *testing.T) {
	assertKinds(t, "{", []TokenKind{int('{')})
	assertKinds(t, "{-", []TokenKind{tokLCurlyHyphen})
	assertKinds(t, "}", []TokenKind{int('}')})
	// {- A -} excluded-pattern shape (MATCH_RECOGNIZE).
	assertKinds(t, "{- A -}", []TokenKind{tokLCurlyHyphen, tokIdent, tokRCurlyHyphen})
}

func TestOperators_DisambiguateVbarConcat(t *testing.T) {
	assertKinds(t, "|", []TokenKind{int('|')})
	assertKinds(t, "||", []TokenKind{tokConcat})
	assertKinds(t, "|||", []TokenKind{tokConcat, int('|')})
}

func TestOperators_BareBangIsInvalid(t *testing.T) {
	// '!' alone is not a Trino token (only '!=' is). It must be flagged.
	tokens, errs := Tokenize("!")
	if len(errs) == 0 || errs[0].Msg != errUnknownChar {
		t.Errorf("bare !: expected errUnknownChar, got %v", errs)
	}
	if !hasInvalid(tokens) {
		t.Error("bare !: expected tokInvalid")
	}
}

func TestOperators_TrinoAbsentAreRejected(t *testing.T) {
	// Negative gate: characters that some SQL dialects (MySQL/Doris/pg) tokenize
	// but Trino's lexer does NOT have any token for. TrinoLexer.g4 routes them to
	// the UNRECOGNIZED_ catch-all, so omni must flag each (tokInvalid + error)
	// rather than silently accept an over-permissive token. Notably absent:
	// '&', '~', '@', '#', backslash, and the bare operators ':=', '::', '<<',
	// '>>', '&&' (only their constituent single chars or nothing exist).
	for _, ch := range []string{"&", "~", "@", "#", "\\"} {
		tokens, errs := Tokenize(ch)
		if len(errs) == 0 || errs[0].Msg != errUnknownChar {
			t.Errorf("Trino-absent %q: expected errUnknownChar, got errs=%v", ch, errs)
		}
		if !hasInvalid(tokens) {
			t.Errorf("Trino-absent %q: expected tokInvalid", ch)
		}
	}
}

// ---- 6. Comments ------------------------------------------------------------

func TestComment_SimpleToEOL(t *testing.T) {
	// '-- comment' to end of line. Unlike MySQL/Doris, no space is required after --.
	assertKinds(t, "42 -- this is a comment", []TokenKind{tokInteger})
	assertKinds(t, "1 --no-space-needed\n2", []TokenKind{tokInteger, tokInteger})
}

func TestComment_SimpleNoSpaceStillComment(t *testing.T) {
	// '--5' is a comment in Trino (the whole rest of line), so nothing remains.
	assertKinds(t, "--5", nil)
	assertKinds(t, "1 --5\n2", []TokenKind{tokInteger, tokInteger})
}

func TestComment_SimpleAtEOF(t *testing.T) {
	assertKinds(t, "1 --", []TokenKind{tokInteger})
}

func TestComment_Bracketed(t *testing.T) {
	assertKinds(t, "1 /* block */ 2", []TokenKind{tokInteger, tokInteger})
	assertKinds(t, "1 /* line1\nline2 */ 2", []TokenKind{tokInteger, tokInteger})
}

func TestComment_BracketedNonNesting(t *testing.T) {
	// Trino bracketed comments do NOT nest (the grammar uses '.*?' non-greedy):
	// the first '*/' closes the comment, leaving 'b */ 2' to be lexed.
	// '/* a /* b */' => comment is '/* a /* b */', then ' 2' remains: just '2'.
	assertKinds(t, "/* a /* b */ 2", []TokenKind{tokInteger})
	// With trailing content after the FIRST close: '/* x */ y */' => comment
	// '/* x */', then 'y', then '*' '/'.
	assertKinds(t, "/* x */ y */", []TokenKind{tokIdent, int('*'), int('/')})
}

func TestComment_Unterminated(t *testing.T) {
	_, errs := Tokenize("/* no end")
	if len(errs) == 0 || errs[0].Msg != errUnterminatedComment {
		t.Errorf("unterminated comment: got %v", errs)
	}
}

func TestComment_NoHashOrDoubleSlash(t *testing.T) {
	// Trino does NOT support '#' or '//' comments (only -- and /* */). They lex
	// as ordinary tokens. '#' is unrecognized; '//' is two SLASH_ tokens.
	tokens, errs := Tokenize("#")
	if len(errs) == 0 || !hasInvalid(tokens) {
		t.Errorf("'#': expected unknown-char error + tokInvalid, got errs=%v", errs)
	}
	assertKinds(t, "//", []TokenKind{int('/'), int('/')})
}

// ---- 7. Error recovery ------------------------------------------------------

func TestRecovery_AfterUnterminatedString(t *testing.T) {
	// After an unterminated string the lexer emits tokInvalid then EOF; it does
	// not loop or panic.
	tokens, errs := Tokenize("SELECT 'oops")
	if len(errs) == 0 {
		t.Fatal("expected lex error")
	}
	// SELECT must still have been emitted before the bad string.
	if tokens[0].Kind != kwSELECT {
		t.Errorf("first token: got %s, want SELECT", TokenName(tokens[0].Kind))
	}
	if !hasInvalid(tokens) {
		t.Error("expected tokInvalid")
	}
}

func TestRecovery_UnknownChar(t *testing.T) {
	tokens, errs := Tokenize("\x01")
	if len(errs) == 0 || errs[0].Msg != errUnknownChar {
		t.Errorf("unknown char: got %v", errs)
	}
	if !hasInvalid(tokens) {
		t.Error("unknown char: expected tokInvalid")
	}
}

// ---- 8. Position tracking ---------------------------------------------------

func TestPosition_BasicOffsets(t *testing.T) {
	tok := singleToken(t, "SELECT")
	if tok.Loc.Start != 0 || tok.Loc.End != 6 {
		t.Errorf("SELECT loc: [%d,%d), want [0,6)", tok.Loc.Start, tok.Loc.End)
	}
}

func TestPosition_LeadingWhitespace(t *testing.T) {
	tokens, _ := Tokenize("   42")
	tok := filterNonEOF(tokens)[0]
	if tok.Loc.Start != 3 || tok.Loc.End != 5 {
		t.Errorf("loc: [%d,%d), want [3,5)", tok.Loc.Start, tok.Loc.End)
	}
}

func TestPosition_MultipleTokens(t *testing.T) {
	tokens := filterNonEOF(mustTokens(t, "a.b"))
	want := [][2]int{{0, 1}, {1, 2}, {2, 3}} // a, ., b
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	for i, w := range want {
		if tokens[i].Loc.Start != w[0] || tokens[i].Loc.End != w[1] {
			t.Errorf("token[%d] loc [%d,%d), want [%d,%d)", i, tokens[i].Loc.Start, tokens[i].Loc.End, w[0], w[1])
		}
	}
}

func TestPosition_EOFLoc(t *testing.T) {
	tokens, _ := Tokenize("ab")
	last := tokens[len(tokens)-1]
	if last.Kind != tokEOF || last.Loc.Start != 2 || last.Loc.End != 2 {
		t.Errorf("EOF loc: kind=%s [%d,%d), want EOF [2,2)", TokenName(last.Kind), last.Loc.Start, last.Loc.End)
	}
}

func TestPosition_BaseOffset(t *testing.T) {
	const base = 100
	l := NewLexerWithOffset("SELECT", base)
	tok := l.NextToken()
	if tok.Loc.Start != base || tok.Loc.End != base+6 {
		t.Errorf("shifted loc: [%d,%d), want [%d,%d)", tok.Loc.Start, tok.Loc.End, base, base+6)
	}
}

func TestPosition_BaseOffsetError(t *testing.T) {
	const base = 50
	l := NewLexerWithOffset("'unterminated", base)
	for l.NextToken().Kind != tokEOF { //nolint:revive
	}
	errs := l.Errors()
	if len(errs) == 0 || errs[0].Loc.Start != base {
		t.Errorf("shifted error loc: %v, want Start=%d", errs, base)
	}
}

func mustTokens(t *testing.T, input string) []Token {
	t.Helper()
	tokens, errs := Tokenize(input)
	if len(errs) > 0 {
		t.Fatalf("Tokenize(%q): unexpected errors %v", input, errs)
	}
	return tokens
}

// ---- 9. Edge cases ----------------------------------------------------------

func TestEdge_EmptyInput(t *testing.T) {
	tokens, errs := Tokenize("")
	if len(errs) != 0 {
		t.Errorf("empty: unexpected errors %v", errs)
	}
	if len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Errorf("empty: expected exactly EOF, got %v", tokens)
	}
}

func TestEdge_WhitespaceOnly(t *testing.T) {
	tokens, errs := Tokenize("  \t\r\n  ")
	if len(errs) != 0 || len(tokens) != 1 || tokens[0].Kind != tokEOF {
		t.Errorf("whitespace-only: errs=%v tokens=%v", errs, tokens)
	}
}

func TestEdge_AlwaysEndsWithEOF(t *testing.T) {
	for _, input := range []string{"", "SELECT 1", "   ", "'oops", "/*x", "X'00"} {
		tokens, _ := Tokenize(input)
		if len(tokens) == 0 || tokens[len(tokens)-1].Kind != tokEOF {
			t.Errorf("input %q: last token must be EOF, got %v", input, tokens)
		}
	}
}

func TestEdge_MultiStatementSplitMarkers(t *testing.T) {
	// SplitSQL relies on SEMICOLON_ tokens with positions.
	kinds := tokenizeKinds(t, "SELECT 1; SELECT 2")
	want := []TokenKind{kwSELECT, tokInteger, int(';'), kwSELECT, tokInteger}
	if len(kinds) != len(want) {
		t.Fatalf("multi-statement: got %v, want %v", kinds, want)
	}
	for i, k := range want {
		if kinds[i] != k {
			t.Errorf("multi-statement[%d]: got %s, want %s", i, TokenName(kinds[i]), TokenName(k))
		}
	}
}

// ---- 10. Keyword classification ---------------------------------------------

func TestClass_ReservedSamples(t *testing.T) {
	reserved := []TokenKind{kwSELECT, kwFROM, kwWHERE, kwJOIN, kwTABLE, kwUNION,
		kwAND, kwOR, kwNOT, kwCASE, kwSKIP, kwTRIM, kwLISTAGG, kwCURRENT_DATE}
	for _, k := range reserved {
		if !IsReserved(k) {
			t.Errorf("%s should be reserved", TokenName(k))
		}
	}
}

func TestClass_NonReservedSamples(t *testing.T) {
	nonReserved := []TokenKind{kwARRAY, kwMAP, kwSHOW, kwGRANT, kwMERGE, kwJSON,
		kwMATCH_RECOGNIZE, kwMATERIALIZED, kwSESSION, kwSTRING, kwTRY_CAST, kwCOMMENT}
	for _, k := range nonReserved {
		if IsReserved(k) {
			t.Errorf("%s should be non-reserved", TokenName(k))
		}
	}
}

func TestClass_LiteralsNotReserved(t *testing.T) {
	for _, k := range []TokenKind{tokInteger, tokDecimal, tokDouble, tokString,
		tokUnicodeString, tokBinaryLiteral, tokIdent, tokQuotedIdent,
		tokBackquotedIdent, tokDigitIdent, tokQuestion} {
		if IsReserved(k) {
			t.Errorf("literal/ident kind %s should not be reserved", TokenName(k))
		}
	}
}

func TestClass_Counts(t *testing.T) {
	// Guardrails on the grammar-derived classification, counted directly from
	// TrinoLexer.g4 (292 alphabetic keyword literals + UTF8/UTF16/UTF32 whose
	// literals contain digits = 295 keyword tokens) and the nonReserved rule in
	// TrinoParser.g4 (213 tokens), therefore 295-213 = 82 reserved.
	if got := len(keywordMap); got != 295 {
		t.Errorf("keywordMap size: got %d, want 295", got)
	}
	if got := len(nonReservedKeywords); got != 213 {
		t.Errorf("nonReservedKeywords size: got %d, want 213", got)
	}
	reserved := 0
	seen := map[TokenKind]bool{}
	for _, k := range keywordMap {
		if seen[k] {
			continue
		}
		seen[k] = true
		if IsReserved(k) {
			reserved++
		}
	}
	if reserved != 82 {
		t.Errorf("reserved keyword count: got %d, want 82", reserved)
	}
}

func TestKeywords_UTFEncodingTokens(t *testing.T) {
	// Regression for the digit-bearing keyword literals UTF8/UTF16/UTF32, used by
	// the JSON ... ENCODING clause. They must be (non-reserved) keyword tokens,
	// not plain identifiers.
	cases := []struct {
		input string
		want  TokenKind
	}{
		{"UTF8", kwUTF8}, {"utf8", kwUTF8},
		{"UTF16", kwUTF16}, {"UTF32", kwUTF32},
	}
	for _, tc := range cases {
		tok := singleToken(t, tc.input)
		if tok.Kind != tc.want {
			t.Errorf("%q: got %s, want %s", tc.input, TokenName(tok.Kind), TokenName(tc.want))
		}
		if IsReserved(tc.want) {
			t.Errorf("%q should be non-reserved", tc.input)
		}
	}
}

func TestClass_OutOfRangeKindNotReserved(t *testing.T) {
	// IsReserved must return false for a kind above the keyword range, not treat
	// every value >= 700 as a reserved keyword.
	if IsReserved(maxKeywordKind + 1) {
		t.Error("kind above maxKeywordKind must not be reported reserved")
	}
	if IsReserved(99999) {
		t.Error("arbitrary large kind must not be reported reserved")
	}
}

// ---- 11. KeywordToken / TokenName -------------------------------------------

func TestKeywordToken_CaseInsensitive(t *testing.T) {
	for _, s := range []string{"select", "SELECT", "Select", "SeLeCt"} {
		kind, ok := KeywordToken(s)
		if !ok || kind != kwSELECT {
			t.Errorf("KeywordToken(%q): kind=%d ok=%v", s, kind, ok)
		}
	}
}

func TestKeywordToken_NonKeyword(t *testing.T) {
	if _, ok := KeywordToken("notakeyword"); ok {
		t.Error(`KeywordToken("notakeyword"): expected not found`)
	}
}

func TestTokenName_Samples(t *testing.T) {
	cases := []struct {
		kind TokenKind
		want string
	}{
		{kwSELECT, "SELECT"}, {kwFROM, "FROM"}, {kwSTRING, "STRING"},
		{kwMATCH_RECOGNIZE, "MATCH_RECOGNIZE"},
		{tokEOF, "EOF"}, {tokInvalid, "INVALID"},
		{tokInteger, "INTEGER_VALUE"}, {tokDecimal, "DECIMAL_VALUE"}, {tokDouble, "DOUBLE_VALUE"},
		{tokString, "STRING"}, {tokUnicodeString, "UNICODE_STRING"}, {tokBinaryLiteral, "BINARY_LITERAL"},
		{tokIdent, "IDENTIFIER"}, {tokQuotedIdent, "QUOTED_IDENTIFIER"},
		{tokBackquotedIdent, "BACKQUOTED_IDENTIFIER"}, {tokDigitIdent, "DIGIT_IDENTIFIER"},
		{tokQuestion, "?"},
		{tokNotEq, "<>"}, {tokLessEq, "<="}, {tokGreaterEq, ">="}, {tokConcat, "||"},
		{tokArrow, "->"}, {tokLeftArrow, "<-"}, {tokDoubleArrow, "=>"},
		{tokLCurlyHyphen, "{-"}, {tokRCurlyHyphen, "-}"},
		{int('+'), "+"}, {int(';'), ";"}, {int(':'), ":"}, {int('('), "("},
	}
	for _, tc := range cases {
		if got := TokenName(tc.kind); got != tc.want {
			t.Errorf("TokenName(%d): got %q, want %q", tc.kind, got, tc.want)
		}
	}
}
