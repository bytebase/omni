package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// Lexer is a hand-written tokenizer for PartiQL source code.
//
// Single-pass scanner. The caller drives it via Next(); each call
// returns one token. At end of input or after a lex error, Next()
// returns Token{Type: tokEOF, ...}. The first error encountered is
// stored in Err and all subsequent Next() calls return tokEOF.
type Lexer struct {
	input string // source text
	pos   int    // current read position (next byte to consume)
	start int    // byte offset of token currently being scanned
	Err   error  // first error encountered, nil if none
}

// NewLexer creates a Lexer for the given source string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token from the input.
// At end of input or after a lex error, returns Token{Type: tokEOF, ...}.
// After Err is set, all subsequent calls return tokEOF.
func (l *Lexer) Next() Token {
	// EOF tokens use l.pos for both Start and End: no token is under
	// construction yet, so l.start may still reflect the previous call.
	// Scan helpers use l.loc() for the {l.start, l.pos} range instead.
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.skipWhitespaceAndComments()
	if l.Err != nil {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString()
	case ch == '"':
		return l.scanQuotedIdent()
	case ch == '`':
		return l.scanIonLiteral()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]):
		return l.scanNumber() // leading-dot decimal: .5
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	default:
		return l.scanOperator()
	}
}

// skipWhitespaceAndComments advances l.pos past whitespace, line comments,
// and block comments. All three are on the HIDDEN channel per the grammar
// and never appear in the token stream.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace.
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		// Line comment: -- to end of line.
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		// Block comment: /* ... */ (non-nested, greedy-shortest).
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			start := l.pos
			l.pos += 2
			closed := false
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					closed = true
					break
				}
				l.pos++
			}
			if !closed {
				l.Err = fmt.Errorf("unterminated block comment at position %d", start)
				return
			}
			continue
		}

		break
	}
}

// ============================================================================
// Character class helpers.
// ============================================================================

// isIdentStart reports whether ch can begin a PartiQL identifier.
// PartiQL identifiers start with [a-zA-Z_$].
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' || ch == '$'
}

// isIdentContinue reports whether ch can appear in a PartiQL identifier
// after the first character. Adds digits to isIdentStart.
func isIdentContinue(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}

// isDigit reports whether ch is an ASCII decimal digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isHexDigit reports whether ch is an Ion HEX_DIGIT (g4 533-534: [0-9A-F]).
// caseInsensitive=true means lowercase a-f match as well.
func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

// hasHexRun reports whether the n bytes of l-input starting at index `at`
// are all HEX_DIGITs and in range. Used to validate the fixed-width hex
// runs in Ion HEX_ESCAPE / UNICODE_ESCAPE.
func hasHexRun(s string, at, n int) bool {
	if at+n > len(s) {
		return false
	}
	for i := 0; i < n; i++ {
		if !isHexDigit(s[at+i]) {
			return false
		}
	}
	return true
}

// lowerASCII folds an ASCII letter to lowercase, leaving other bytes
// unchanged. Mirrors the grammar's caseInsensitive=true handling for the
// single-letter escape codes.
func lowerASCII(ch byte) byte {
	if ch >= 'A' && ch <= 'Z' {
		return ch + ('a' - 'A')
	}
	return ch
}

// ============================================================================
// Token location helper.
// ============================================================================

// loc returns an ast.Loc covering the byte range [l.start, l.pos).
// This is the standard "current token" range used by every scan helper
// when constructing the returned Token.
func (l *Lexer) loc() ast.Loc {
	return ast.Loc{Start: l.start, End: l.pos}
}

// ============================================================================
// Scan helpers — one per leading-character class dispatched by Next().
// ============================================================================

// scanString scans a single-quoted PartiQL string literal and returns a
// tokSCONST token. The only escape mechanism in PartiQL strings is the
// doubled-apostrophe form: two consecutive apostrophes in a row stand for
// a single apostrophe in the decoded value. There are no backslash escapes.
// Token.Str holds the fully decoded string content (without the surrounding
// quotes). Next() has already set l.start to the offset of the opening quote
// before calling this function, so scanString must not modify l.start.
func (l *Lexer) scanString() Token {
	l.pos++ // skip the opening apostrophe
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			// Two apostrophes in a row encode a single apostrophe.
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			// Single closing apostrophe — end of the literal.
			l.pos++
			return Token{Type: tokSCONST, Str: buf.String(), Loc: l.loc()}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	// Reached end of input without finding the closing apostrophe.
	l.Err = fmt.Errorf("unterminated string literal at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanQuotedIdent scans a double-quoted PartiQL identifier and returns a
// tokIDENT_QUOTED token. Quoted identifiers preserve case and are NOT looked
// up in the keyword map — they are always identifiers, never keywords. The
// only escape mechanism is the doubled double-quote form: "" within the
// identifier stands for a single " in the decoded name. Token.Str holds the
// decoded identifier text (without the surrounding quotes). Next() has
// already set l.start to the offset of the opening quote before calling this
// function, so scanQuotedIdent must not modify l.start.
func (l *Lexer) scanQuotedIdent() Token {
	l.pos++ // skip the opening double-quote
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			// Two double-quotes in a row encode a single double-quote.
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			// Single closing double-quote — end of the identifier.
			l.pos++
			return Token{Type: tokIDENT_QUOTED, Str: buf.String(), Loc: l.loc()}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	// Reached end of input without finding the closing double-quote.
	l.Err = fmt.Errorf("unterminated quoted identifier at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanIdentOrKeyword consumes an unquoted identifier and looks it up in
// the keywords map. If the lowercased text matches a keyword, returns
// that keyword token; otherwise returns tokIDENT.
//
// Token.Str preserves the original case (so the parser/AST can render
// identifiers as written). Keyword matching is case-insensitive per
// PartiQLLexer.g4 caseInsensitive=true.
//
// Grammar: IDENTIFIER : [A-Z$_][A-Z0-9$_]*;
//
//	(with caseInsensitive=true expanding to [a-zA-Z$_][a-zA-Z0-9$_]*)
func (l *Lexer) scanIdentOrKeyword() Token {
	for l.pos < len(l.input) && isIdentContinue(l.input[l.pos]) {
		l.pos++
	}
	raw := l.input[l.start:l.pos]
	lower := strings.ToLower(raw)
	if tt, ok := keywords[lower]; ok {
		return Token{
			Type: tt,
			Str:  raw,
			Loc:  l.loc(),
		}
	}
	return Token{
		Type: tokIDENT,
		Str:  raw,
		Loc:  l.loc(),
	}
}

// scanNumber consumes an integer or decimal literal. Returns tokICONST
// for plain integers and tokFCONST for any number with a decimal point
// or scientific exponent. Token.Str is the raw source text.
//
// Grammar:
//
//	LITERAL_INTEGER : DIGIT+;
//	LITERAL_DECIMAL :
//	    DIGIT+ '.' DIGIT* ([e] [+-]? DIGIT+)?
//	  | '.' DIGIT+ ([e] [+-]? DIGIT+)?
//	  | DIGIT+ ([e] [+-]? DIGIT+)?
//	  ;
//
// (caseInsensitive=true means [e] matches both 'e' and 'E'.)
func (l *Lexer) scanNumber() Token {
	isFloat := false

	// Leading-dot form (.5).
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	} else {
		// Integer part.
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
		// Optional fraction.
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		}
	}

	// Optional scientific exponent.
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	tt := tokICONST
	if isFloat {
		tt = tokFCONST
	}
	return Token{
		Type: tt,
		Str:  l.input[l.start:l.pos],
		Loc:  l.loc(),
	}
}

// scanOperator consumes a one or two-character operator or punctuation
// token. Two-character operators (<=, >=, <>, <<, >>, ||, !=) are
// matched first via lookahead; otherwise the single-character cases
// fall through. Unrecognized characters set l.Err.
func (l *Lexer) scanOperator() Token {
	ch := l.input[l.pos]
	l.pos++

	// Two-character lookahead.
	if l.pos < len(l.input) {
		next := l.input[l.pos]
		switch {
		case ch == '<' && next == '=':
			l.pos++
			return Token{Type: tokLT_EQ, Str: "<=", Loc: l.loc()}
		case ch == '<' && next == '>':
			l.pos++
			return Token{Type: tokNEQ, Str: "<>", Loc: l.loc()}
		case ch == '<' && next == '<':
			l.pos++
			return Token{Type: tokANGLE_DOUBLE_LEFT, Str: "<<", Loc: l.loc()}
		case ch == '>' && next == '=':
			l.pos++
			return Token{Type: tokGT_EQ, Str: ">=", Loc: l.loc()}
		case ch == '>' && next == '>':
			l.pos++
			return Token{Type: tokANGLE_DOUBLE_RIGHT, Str: ">>", Loc: l.loc()}
		case ch == '|' && next == '|':
			l.pos++
			return Token{Type: tokCONCAT, Str: "||", Loc: l.loc()}
		case ch == '!' && next == '=':
			l.pos++
			return Token{Type: tokNEQ, Str: "!=", Loc: l.loc()}
		}
	}

	// Single-character operators / punctuation.
	switch ch {
	case '+':
		return Token{Type: tokPLUS, Str: "+", Loc: l.loc()}
	case '-':
		return Token{Type: tokMINUS, Str: "-", Loc: l.loc()}
	case '*':
		return Token{Type: tokASTERISK, Str: "*", Loc: l.loc()}
	case '/':
		return Token{Type: tokSLASH_FORWARD, Str: "/", Loc: l.loc()}
	case '%':
		return Token{Type: tokPERCENT, Str: "%", Loc: l.loc()}
	case '^':
		return Token{Type: tokCARET, Str: "^", Loc: l.loc()}
	case '~':
		return Token{Type: tokTILDE, Str: "~", Loc: l.loc()}
	case '@':
		return Token{Type: tokAT_SIGN, Str: "@", Loc: l.loc()}
	case '=':
		return Token{Type: tokEQ, Str: "=", Loc: l.loc()}
	case '<':
		return Token{Type: tokLT, Str: "<", Loc: l.loc()}
	case '>':
		return Token{Type: tokGT, Str: ">", Loc: l.loc()}
	case '(':
		return Token{Type: tokPAREN_LEFT, Str: "(", Loc: l.loc()}
	case ')':
		return Token{Type: tokPAREN_RIGHT, Str: ")", Loc: l.loc()}
	case '[':
		return Token{Type: tokBRACKET_LEFT, Str: "[", Loc: l.loc()}
	case ']':
		return Token{Type: tokBRACKET_RIGHT, Str: "]", Loc: l.loc()}
	case '{':
		return Token{Type: tokBRACE_LEFT, Str: "{", Loc: l.loc()}
	case '}':
		return Token{Type: tokBRACE_RIGHT, Str: "}", Loc: l.loc()}
	case ':':
		return Token{Type: tokCOLON, Str: ":", Loc: l.loc()}
	case ';':
		return Token{Type: tokCOLON_SEMI, Str: ";", Loc: l.loc()}
	case ',':
		return Token{Type: tokCOMMA, Str: ",", Loc: l.loc()}
	case '.':
		return Token{Type: tokPERIOD, Str: ".", Loc: l.loc()}
	case '?':
		return Token{Type: tokQUESTION_MARK, Str: "?", Loc: l.loc()}
	}

	l.Err = fmt.Errorf("unexpected character %q at position %d", ch, l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanIonLiteral consumes a backtick-delimited inline Ion value: `…`.
//
// On the opening backtick the lexer enters a dedicated Ion sub-mode
// (the hand-written analogue of PartiQLLexer.g4's `pushMode(ION)`):
// it consumes an Ion value verbatim until a *standalone* backtick —
// the ION_CLOSURE that pops the mode (g4 line 428) — and emits one
// tokION_LITERAL. Token.Str is the verbatim inner content between the
// backticks (no Ion parsing, no decoding); Token.Loc covers the entire
// `…` range including both backticks.
//
// The whole point of the sub-mode is that a backtick appearing INSIDE
// an Ion sub-token does not terminate the literal. Per the ION mode
// grammar (g4 lines 406-430), the sub-tokens that may legally contain a
// backtick — and therefore must be consumed whole — are:
//
//   - // line comment           (ION_INLINE_COMMENT, g4 408-409; EOF-terminable)
//   - /* ... */ block comment   (ION_BLOCK_COMMENT,  g4 411-412)
//   - "..." short string        (SHORT_QUOTED_STRING, g4 417-419)
//   - triple-quoted long string (LONG_QUOTED_STRING,  g4 421-423)
//   - '...' quoted symbol       (QUOTED_SYMBOL,       g4 425-426)
//
// Any other byte is consumed verbatim (ION_ANY, g4 430).
//
// ANTLR lexer rules are maximal-munch WITH FALLBACK: a construct only
// consumes its bytes if it can reach an accept state. A double-quote,
// single-quote, triple-quote, or block-comment opener that never finds its
// closer does NOT match its multi-byte rule; instead each of those bytes
// degrades to the single-byte ION_ANY rule and scanning continues.
// Concretely, the input backtick + "abc + backtick lexes as ION_ANY for the
// '"' then 'a','b','c' then ION_CLOSURE — a *complete* literal with inner
// content "abc, not an error. We mirror that by rewinding to one byte past
// the opener and continuing on any sub-scanner failure. The literal is
// unterminated only when EOF is reached with no standalone closing backtick
// (e.g. backtick + abc with no closing backtick). This was verified
// differentially against the generated ANTLR lexer
// (github.com/bytebase/parser/partiql) used as the oracle; see the lob note.
//
// Note on Ion blobs/clobs (ION_BLOB, g4 414-415): the grammar's ION_BLOB
// rule matches base64 + whitespace only, and base64 contains no backtick.
// A clob such as {{ "}}" }} is therefore NOT a single ION_BLOB token in
// this grammar — ANTLR tokenizes it as ION_ANY '{', ION_ANY '{', a
// SHORT_QUOTED_STRING for "}}", then ION_ANY for the trailing }}, and a
// standalone backtick closes. The }} inside the quoted string is protected
// by the string rule, and the braces are plain ION_ANY. Hence there is no
// lob special case here: '{' and '}' are ION_ANY, and the quoted-string
// scanners (invoked anywhere a quote appears, including between braces)
// consume any }} inside a string as content for free.
//
// The AWS DynamoDB PartiQL corpus has zero real Ion literals (the only
// two backtick uses are placeholder skeletons filtered out of the corpus
// smoke tests), so these forms come from the PartiQL/Ion specs.
func (l *Lexer) scanIonLiteral() Token {
	l.pos++ // skip opening backtick
	contentStart := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		// save marks the opener of any multi-byte construct so we can
		// rewind to one byte past it on a failed match (ION_ANY fallback).
		save := l.pos
		switch {
		case ch == '`':
			// Standalone backtick — ION_CLOSURE. Pop the sub-mode.
			content := l.input[contentStart:l.pos]
			l.pos++ // skip closing backtick
			return Token{Type: tokION_LITERAL, Str: content, Loc: l.loc()}

		case ch == '{' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '{':
			// Possible ION_BLOB: {{ base64 }} (g4 414-415). ANTLR matches
			// this whole region as ONE maximal-munch token, so any byte
			// inside it — including '//' or '/*', which base64 may legally
			// contain ('/' is a BASE_64_CHAR) — is blob CONTENT and never a
			// comment. Only a *valid* base64 lob matches; if it does, consume
			// it whole. Otherwise the '{' degrades to ION_ANY (so '{{' alone,
			// a clob whose content is a quoted string, '{{//}}', etc. fall
			// through to the per-byte handling that mirrors ANTLR).
			if !l.scanIonBlob() {
				l.pos = save + 1 // the leading '{' is ION_ANY
			}

		case ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/':
			// Line comment to newline or EOF. If it reaches EOF it
			// swallows the closing backtick; the loop then exits and the
			// literal is reported unterminated — matching ANTLR, where
			// ION_INLINE_COMMENT consumes through EOF and the mode is
			// never popped.
			l.scanIonLineComment()

		case ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*':
			if !l.scanIonBlockComment() {
				// No closing */: /* does not match ION_BLOCK_COMMENT;
				// degrade the '/' to ION_ANY and continue.
				l.pos = save + 1
			}

		case ch == '"':
			if !l.scanIonQuoted('"') {
				// No closing ": the '"' is ION_ANY, not a string opener.
				l.pos = save + 1
			}

		case ch == '\'' && l.matchesTripleQuote():
			// A '''-prefixed run is a long string only if it actually
			// closes with another '''. On failure, ANTLR maximal-munch
			// falls back to QUOTED_SYMBOL (the quotes scan as (possibly
			// empty) single-quoted symbols), and a symbol that itself
			// cannot close degrades to ION_ANY. Rewind down that ladder
			// rather than erroring.
			if !l.scanIonLongString() {
				l.pos = save
				if !l.scanIonQuoted('\'') {
					l.pos = save + 1 // the leading "'" is ION_ANY
				}
			}

		case ch == '\'':
			if !l.scanIonQuoted('\'') {
				// No closing ': the "'" is ION_ANY, not a symbol opener.
				l.pos = save + 1
			}

		default:
			l.pos++ // ION_ANY
		}
	}
	return l.unterminatedIon()
}

// unterminatedIon records the unterminated-Ion-literal error (anchored
// to the opening backtick at l.start) and returns the EOF sentinel.
func (l *Lexer) unterminatedIon() Token {
	l.Err = fmt.Errorf("unterminated Ion literal at position %d", l.start)
	return Token{Type: tokEOF, Loc: ast.Loc{Start: l.start, End: l.start}}
}

// scanIonBlob attempts to match a complete Ion BLOB starting at l.pos and,
// on success, advances l.pos past the closing '}}' and returns true. On any
// failure l.pos is left unchanged and it returns false. l.pos must point at
// the first '{' of a '{{'.
//
// Grammar (g4 414-415, 468-489), with caseInsensitive=true:
//
//	ION_BLOB        : '{{' (BASE_64_QUARTET | WS)* BASE_64_PAD? WS* '}}'
//	BASE_64_QUARTET : C WS* C WS* C WS* C
//	BASE_64_PAD     : C WS* C WS* C WS* '='   |   C WS* C WS* '=' WS* '='
//	C (BASE_64_CHAR): [0-9A-Za-z+/]
//	WS              : [ \r\n\t]
//
// Because base64 chars, '=', and '}' are pairwise disjoint and WS is allowed
// freely between the significant characters, a valid lob body is exactly:
// some base64 chars whose count is a multiple of four (the quartets),
// followed by an optional pad contributing three chars + one '=' or two
// chars + two '=', then '}}'. This matcher consumes base64 chars and WS,
// counts the base64 chars, validates the trailing pad shape, and requires
// the closing '}}'. ANTLR matches ION_BLOB maximally only when it can reach
// '}}'; if it cannot, the caller treats the leading '{' as ION_ANY.
func (l *Lexer) scanIonBlob() bool {
	p := l.pos + 2 // past '{{'
	chars := 0     // count of base64 chars seen so far (excludes '=')
	eq := 0        // count of '=' padding chars seen (only valid at the end)
	for p < len(l.input) {
		c := l.input[p]
		switch {
		case c == ' ' || c == '\t' || c == '\r' || c == '\n':
			p++ // WS, allowed anywhere between significant chars
		case c == '}':
			// Candidate LOB_END. Require the second '}' and a structurally
			// valid character count: the base64 chars must form whole
			// quartets, adjusted for the pad ('=' count).
			if p+1 < len(l.input) && l.input[p+1] == '}' && validBase64Counts(chars, eq) {
				l.pos = p + 2 // consume through '}}'
				return true
			}
			return false
		case c == '=':
			eq++
			if eq > 2 {
				return false // a pad has at most two '='
			}
			p++
		case isBase64Char(c):
			if eq > 0 {
				return false // base64 char after '=' is not a valid pad
			}
			chars++
			p++
		default:
			return false // anything else (incl. a backtick or quote) breaks the lob
		}
	}
	return false // reached EOF without a closing '}}'
}

// validBase64Counts reports whether `chars` base64 characters followed by
// `eq` '=' padding characters form a valid Ion base64 lob body:
//   - no padding: chars is a non-negative multiple of 4 (quartets only);
//   - one '=' (BASE_64_PAD1: C C C '='): the three pad chars complete a
//     quartet, so chars ≡ 3 (mod 4) and chars >= 3;
//   - two '=' (BASE_64_PAD2: C C '=' '='): chars ≡ 2 (mod 4) and chars >= 2.
func validBase64Counts(chars, eq int) bool {
	switch eq {
	case 0:
		return chars%4 == 0
	case 1:
		return chars >= 3 && chars%4 == 3
	case 2:
		return chars >= 2 && chars%4 == 2
	default:
		return false
	}
}

// isBase64Char reports whether ch is an Ion BASE_64_CHAR (g4 488-489:
// [0-9A-Z+/]); caseInsensitive=true also admits a-z.
func isBase64Char(ch byte) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z') ||
		ch == '+' || ch == '/'
}

// scanIonLineComment consumes an Ion // line comment (g4 408-409).
// The comment runs to the next newline or to EOF — EOF is a legal
// terminator for a line comment, so l.pos may land at len(l.input).
// Positioned at the first '/'.
func (l *Lexer) scanIonLineComment() {
	l.pos += 2 // skip //
	for l.pos < len(l.input) && l.input[l.pos] != '\n' && l.input[l.pos] != '\r' {
		l.pos++
	}
}

// scanIonBlockComment consumes an Ion /* … */ block comment (g4 411-412)
// and reports whether the closing */ was found. Positioned at the first
// '/'. Returns false if EOF is reached before */.
func (l *Lexer) scanIonBlockComment() bool {
	l.pos += 2 // skip /*
	for l.pos+1 < len(l.input) {
		if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
			l.pos += 2
			return true
		}
		l.pos++
	}
	l.pos = len(l.input)
	return false
}

// ionEscapeLen reports the byte length consumed by a valid Ion TEXT_ESCAPE
// beginning at l.input[at] (including the leading backslash), or 0 if the
// bytes at `at` are not a valid escape. l.input[at] must be a backslash.
//
// TEXT_ESCAPE = COMMON_ESCAPE | HEX_ESCAPE | UNICODE_ESCAPE (g4 465-466):
//
//   - COMMON_ESCAPE  '\' [abtnfrv?0'"/\] | '\' ION_NEWLINE  (g4 501-519)
//   - HEX_ESCAPE     '\x' HEX_DIGIT HEX_DIGIT                (g4 521-522)
//   - UNICODE_ESCAPE '\u' HHHH | '\U000' HHHH H | '\U0010' HHHH (g4 524-528)
//
// caseInsensitive=true (g4 line 4) applies to every literal char in these
// rules, so the escape letters (a,b,t,n,f,r,v,x,u) and HEX_DIGIT [0-9A-F]
// all match either case. In particular '\u' and '\U' are the SAME prefix:
// the generated lexer accepts '\U' HHHH as the 6-byte short UNICODE_ESCAPE
// (verified differentially against the antlr_fallback oracle).
//
// The two longer '\U000…' / '\U0010…' (10-byte) UNICODE_ESCAPE alternatives
// are NOT distinguished here: they only ever differ from the 6-byte short
// form by absorbing four extra HEX_DIGIT bytes, and a backtick (the sole
// byte that can close an Ion literal) is never a HEX_DIGIT — so whether
// those four bytes are part of the escape or are plain string content, the
// surrounding string stays open across them identically and the literal's
// boundary is byte-for-byte the same. Token.Str is verbatim (we never decode
// the escape), so reporting 6 here is observationally equivalent to 10.
//
// A backslash that does NOT begin one of these is not an escape at all: the
// surrounding string/symbol rule cannot match (the backslash is excluded
// from every *_TEXT_ALLOWED set), so the token fails and the caller degrades
// the opener to ION_ANY.
func (l *Lexer) ionEscapeLen(at int) int {
	if at+1 >= len(l.input) {
		return 0 // dangling backslash at EOF: not a complete escape
	}
	c := l.input[at+1]
	switch lowerASCII(c) {
	case 'a', 'b', 't', 'n', 'f', 'r', 'v', '?', '0', '\'', '"', '/', '\\':
		return 2 // COMMON_ESCAPE single-character code
	case 'x':
		// HEX_ESCAPE: '\x' HEX_DIGIT HEX_DIGIT.
		if hasHexRun(l.input, at+2, 2) {
			return 4
		}
		return 0
	case 'u':
		// UNICODE_ESCAPE: '\u' / '\U' followed by a HEX_DIGIT_QUARTET.
		if hasHexRun(l.input, at+2, 4) {
			return 6
		}
		return 0
	case '\r':
		// '\' ION_NEWLINE line continuation: \r\n counts as one newline.
		if at+2 < len(l.input) && l.input[at+2] == '\n' {
			return 3
		}
		return 2
	case '\n':
		return 2
	}
	return 0
}

// ionShortTextDisallowsRaw reports whether b is a raw byte that the Ion
// SHORT string / quoted symbol grammars exclude from their content sets
// (STRING_SHORT_TEXT_ALLOWED g4 451-456, SYMBOL_TEXT_ALLOWED g4 494-499).
// Both rules allow U+0020..U+FFFF (minus the delimiter and backslash,
// handled by the caller) plus exactly the WS_NOT_NL bytes below U+0020 —
// U+0009 tab, U+000B vertical tab, U+000C form feed (g4 536-541).
// Every OTHER C0 control byte (U+0000..U+0008, U+000A LF, U+000D CR,
// U+000E..U+001F) is excluded; a raw occurrence makes the rule fail to
// match (it must be written as an escape instead), so the caller degrades
// the opener to ION_ANY. Bytes >= 0x20 — including UTF-8 lead/continuation
// bytes (>= 0x80) and DEL/C1 (which the g4 ranges admit) — are content.
// This is byte-for-byte the same boundary as the generated ANTLR lexer
// (verified differentially against the antlr_fallback oracle). LONG strings
// use a DIFFERENT (wider) allowed set — they additionally admit raw newlines
// LF/CR — see ionLongTextDisallowsRaw / scanIonLongString.
func ionShortTextDisallowsRaw(b byte) bool {
	if b >= 0x20 {
		return false
	}
	switch b {
	case '\t', '\v', '\f': // U+0009, U+000B, U+000C — WS_NOT_NL, allowed
		return false
	default:
		return true
	}
}

// ionLongTextDisallowsRaw reports whether b is a raw byte that the Ion LONG
// string grammar (STRING_LONG_TEXT_ALLOWED, g4 459-463) excludes from its
// content set. A long string admits U+0020..U+FFFF (minus the delimiter and
// backslash, handled by the caller) plus WS = [space CR LF tab] (the
// WHITESPACE fragment, g4 390-391). Unlike the SHORT set (ionShortTextDisallowsRaw),
// the long set DOES admit raw newlines (U+000A LF, U+000D CR) — long strings
// span lines by design. Below U+0020 the allowed bytes are therefore
// U+0009 tab, U+000A LF, U+000B vtab, U+000C form feed, U+000D CR; every OTHER
// C0 control (U+0000 NUL, U+0001..U+0008, U+000E..U+001F) is excluded.
//
// IMPORTANT: U+000B (vtab) and U+000C (form feed) ARE admitted here, even
// though the g4 *text* (WS = [ \r\n\t]) would exclude them. The GENERATED
// ANTLR lexer — the executable oracle (correctness-protocol.md: the oracle
// decides, antlr's own grammar-text is only a hint) — accepts both as long
// string content, mirroring the SHORT set (where they are WS_NOT_NL). This
// was confirmed byte-for-byte differentially against the generated lexer at
// github.com/bytebase/parser/partiql: a long string carrying a raw VT/FF
// before an inner backtick keeps the long string open across that backtick
// (the literal spans it), whereas a raw NUL/US/etc. fails the rule and the
// backtick closes the literal. See the divergence note in scanIonLongString.
// Bytes >= 0x20 — including UTF-8 lead/continuation bytes (>= 0x80) and DEL/C1
// — are content.
func ionLongTextDisallowsRaw(b byte) bool {
	if b >= 0x20 {
		return false
	}
	switch b {
	case '\t', '\n', '\v', '\f', '\r': // U+0009,000A,000B,000C,000D — allowed
		return false
	default:
		return true
	}
}

// scanIonQuoted consumes a single-character-delimited Ion text token —
// a double-quoted short string (g4 417-419) or a single-quoted symbol
// (g4 425-426) — and reports whether the closing quote was found. A
// backslash escapes the following byte ONLY when it forms a valid Ion
// TEXT_ESCAPE (g4 465-466); an invalid escape (e.g. '\q') means the token
// cannot match and the caller degrades the opener to ION_ANY. A raw
// newline or other disallowed control byte (ionShortTextDisallowsRaw) is
// likewise NOT permitted as content — it must be escaped — so encountering
// one is a match failure, not content (this prevents a stray raw newline
// from holding the string open across a backtick and swallowing the SQL
// that follows). Positioned at the opening quote. Returns false if EOF is
// reached before the closer, on a dangling/invalid backslash, or on a
// disallowed raw control byte — in every such case the opener is ION_ANY.
func (l *Lexer) scanIonQuoted(quote byte) bool {
	l.pos++ // skip opening quote
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == '\\':
			// Escape: only a valid TEXT_ESCAPE keeps the token going. An
			// invalid or dangling escape means this is not a string/symbol
			// at all (the '\' is excluded from the allowed set), so report
			// failure and let the caller treat the opener as ION_ANY.
			n := l.ionEscapeLen(l.pos)
			if n == 0 {
				return false
			}
			l.pos += n
		case ch == quote:
			l.pos++ // skip closing quote
			return true
		case ionShortTextDisallowsRaw(ch):
			// Raw newline / disallowed control byte: not valid SHORT-string
			// or symbol content, so the rule cannot match. Fail and let the
			// caller degrade the opener to ION_ANY (so a later quote does
			// NOT close this token across an intervening backtick).
			return false
		default:
			l.pos++
		}
	}
	return false
}

// scanIonLongString consumes an Ion triple-quoted long string (g4 421-423)
// and reports whether the closing triple-quote was found. Long strings may
// span newlines and honor backslash escapes, but only VALID Ion TEXT_ESCAPEs
// (g4 447-448, 465-466); an invalid escape means the long-string rule cannot
// match (STRING_LONG_TEXT_ALLOWED excludes the backslash), so the caller
// degrades the opener. A raw C0 control byte outside the long-string allowed
// set (ionLongTextDisallowsRaw — NUL etc.; raw LF/CR/VT/FF ARE allowed) is
// likewise NOT valid content: it must be escaped, so encountering one is a
// match failure, not content. Without this, a long string carrying a raw NUL
// before an inner backtick would stay open ACROSS that backtick and close
// only at a LATER ''' — hiding the standalone backtick (the real ION_CLOSURE)
// and swallowing the SQL that follows (Codex round-4 P2). On failure the
// caller degrades the opening quotes through QUOTED_SYMBOL to ION_ANY (the
// same maximal-munch fallback ANTLR applies), so the first standalone backtick
// closes the literal. Positioned at the first of the three opening quotes.
// Returns false if EOF is reached before the closing triple-quote, on an
// invalid/dangling escape, or on a disallowed raw control byte.
func (l *Lexer) scanIonLongString() bool {
	l.pos += 3 // skip opening '''
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == '\\':
			n := l.ionEscapeLen(l.pos)
			if n == 0 {
				return false
			}
			l.pos += n
		case ch == '\'' && l.matchesTripleQuote():
			l.pos += 3 // skip closing '''
			return true
		case ionLongTextDisallowsRaw(ch):
			// Disallowed raw C0 control byte (e.g. NUL): not valid long-string
			// content, so the rule cannot match. Fail and let the caller
			// degrade the opener (so a later ''' does NOT close this token
			// across an intervening standalone backtick).
			return false
		default:
			l.pos++
		}
	}
	return false
}

// matchesTripleQuote reports whether the three bytes at l.pos are three
// single-quotes (the Ion long-string delimiter).
// Caller must ensure l.pos+3 <= len(l.input) for the read to be in range.
func (l *Lexer) matchesTripleQuote() bool {
	return l.pos+3 <= len(l.input) &&
		l.input[l.pos] == '\'' &&
		l.input[l.pos+1] == '\'' &&
		l.input[l.pos+2] == '\''
}
