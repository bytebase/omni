// Package parser implements a hand-written GoogleSQL lexer and (in future DAG
// nodes) a recursive-descent parser. This node ships only the lexer.
//
// GoogleSQL is the SQL dialect shared by Google BigQuery and Google Cloud
// Spanner; one omni parser serves both. The lexer is a faithful port of the
// legacy ANTLR GoogleSQLLexer.g4 (itself a hand-port of ZetaSQL's
// flex_tokenizer.l), covering all 308 word-keywords, the operator/punctuation
// set (including the pipe |> token), string/bytes/numeric/identifier literals
// (raw + triple-quoted + backtick), and the three comment styles (-- # /* */).
package parser

import (
	"strconv"

	"github.com/bytebase/omni/googlesql/ast"
)

// LexError describes a single lexing failure with its source location.
//
// Lex errors are non-fatal: the lexer collects them in a slice (Lexer.Errors)
// and emits a synthetic tokInvalid token at each failure site so consumers can
// choose to halt on the first error or proceed with best-effort parsing.
type LexError struct {
	Loc ast.Loc
	Msg string
}

// Standard lex error messages. Tests assert against these constants so a
// reword in one place propagates everywhere.
const (
	errUnterminatedString       = "unterminated string literal"
	errUnterminatedTripleString = "unterminated triple-quoted string literal"
	errUnterminatedBytes        = "unterminated bytes literal"
	errUnterminatedComment      = "unterminated block comment"
	errUnterminatedIdentifier   = "unterminated quoted identifier"
	errInvalidByte              = "invalid byte"
)

// Lexer is a GoogleSQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring of a larger document) and
// call NextToken until it returns Token{Type: tokEOF}. Lex errors are
// collected in a slice (Errors); each error is accompanied by a synthetic
// tokInvalid token at the failure site so consumers can halt on the first
// error or proceed with best-effort parsing.
type Lexer struct {
	input      string
	pos        int // current byte offset (one past the last consumed byte)
	start      int // start byte of the token currently being scanned
	errors     []LexError
	baseOffset int // added to token/error Loc when returned via the public API
}

// NewLexer constructs a Lexer for the given input. Token positions are
// zero-based byte offsets into input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset constructs a Lexer whose emitted token Loc values are
// shifted by baseOffset. Use this when tokenizing a substring of a larger
// document and you want Loc values to refer to positions in the larger
// document rather than the substring (e.g. a parser-entry that splits a
// multi-statement script then lexes each piece while preserving absolute
// positions).
//
// The shift is applied lazily when tokens and errors leave the Lexer via
// NextToken() and Errors(); internal scan helpers use local (unshifted)
// positions against input.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// Errors returns all lex errors collected so far, in source order. Each error
// is accompanied by a tokInvalid token in the stream. If the Lexer was
// constructed with NewLexerWithOffset, the returned errors have Loc values
// shifted by baseOffset (the internal slice retains unshifted positions).
func (l *Lexer) Errors() []LexError {
	if l.baseOffset == 0 {
		return l.errors
	}
	shifted := make([]LexError, len(l.errors))
	for i, e := range l.errors {
		shifted[i] = LexError{
			Loc: ast.Loc{Start: e.Loc.Start + l.baseOffset, End: e.Loc.End + l.baseOffset},
			Msg: e.Msg,
		}
	}
	return shifted
}

// Tokenize is a one-shot convenience that runs a lexer to EOF and returns the
// full token stream and error list. Useful for tests and for callers that
// don't need streaming.
func Tokenize(input string) (tokens []Token, errors []LexError) {
	l := NewLexer(input)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}

// NextToken advances the lexer and returns the next token. At end of input it
// returns Token{Type: tokEOF}; subsequent calls continue to return EOF.
//
// If the Lexer was constructed with NewLexerWithOffset, the returned token's
// Loc values are shifted by baseOffset.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	if l.baseOffset != 0 {
		tok.Loc.Start += l.baseOffset
		tok.Loc.End += l.baseOffset
	}
	return tok
}

// nextTokenInner returns the next token with LOCAL (unshifted) positions. The
// public NextToken wraps it to apply baseOffset.
func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'' || ch == '"':
		// Plain string literal (no prefix).
		return l.scanString(false)
	case ch == '`':
		return l.scanBacktickIdent()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.peekDigit(1):
		return l.scanNumber()
	case isIdentStart(ch):
		// May be an r/R/b/B-prefixed string or bytes literal, otherwise an
		// identifier or keyword.
		if tok, ok := l.tryScanPrefixedLiteral(); ok {
			return tok
		}
		return l.scanIdentOrKeyword()
	}
	return l.scanOperatorOrPunct(ch)
}

// peek returns the byte at l.pos+offset, or 0 if past end of input.
func (l *Lexer) peek(offset int) byte {
	if l.pos+offset >= len(l.input) {
		return 0
	}
	return l.input[l.pos+offset]
}

// peekDigit reports whether the byte at l.pos+offset is an ASCII digit.
func (l *Lexer) peekDigit(offset int) bool {
	c := l.peek(offset)
	return c >= '0' && c <= '9'
}

// isIdentStart reports whether ch may begin an unquoted identifier.
// GoogleSQL UNQUOTED_IDENTIFIER: [A-Z_] (case-insensitive ⇒ [A-Za-z_]).
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isIdentCont reports whether ch may continue an unquoted identifier.
// GoogleSQL UNQUOTED_IDENTIFIER continuation: [A-Z0-9_]. Note: unlike
// Snowflake, '@' and '$' are NOT identifier characters in GoogleSQL.
func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

// skipWhitespaceAndComments advances l.pos past whitespace, dash comments
// (--), pound comments (#), and block comments (/* */). Comments are dropped
// silently (channel HIDDEN in the grammar); whitespace produces no token.
//
// GoogleSQL specifics:
//   - WHITESPACE: [ \t\f\r\n] — includes form-feed.
//   - DASH_COMMENT: '--' to end-of-line.
//   - POUND_COMMENT: '#' to end-of-line (BigQuery style).
//   - BLOCK_COMMENT: '/**/' | '/*' ~[!] .*? '*/' — does NOT nest, and a
//     comment opener immediately followed by '!' (i.e. "/*!") is NOT a comment.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == ' ', ch == '\t', ch == '\f', ch == '\r', ch == '\n':
			l.pos++
		case ch == '-' && l.peek(1) == '-':
			l.pos += 2
			l.skipToEndOfLine()
		case ch == '#':
			l.pos++
			l.skipToEndOfLine()
		case ch == '/' && l.peek(1) == '*' && l.peek(2) != '!':
			l.scanBlockComment()
		default:
			return
		}
	}
}

// skipToEndOfLine advances l.pos to just past the next newline (or EOF). Used
// by -- and # line comments.
func (l *Lexer) skipToEndOfLine() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' && l.input[l.pos] != '\r' {
		l.pos++
	}
	// Consume the line terminator (\r, \n, or \r\n) so the comment span
	// matches DASH_COMMENT/POUND_COMMENT, which optionally include it.
	if l.pos < len(l.input) && l.input[l.pos] == '\r' {
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '\n' {
		l.pos++
	}
}

// scanBlockComment advances past a /* ... */ block comment. GoogleSQL block
// comments do NOT nest (BLOCK_COMMENT: '/*' ~[!] .*? '*/', non-greedy): the
// comment ends at the FIRST '*/'. On an unterminated comment, appends a
// LexError but emits no token (comments are HIDDEN). The caller has already
// verified the opener is "/*" not followed by '!'.
func (l *Lexer) scanBlockComment() {
	start := l.pos
	l.pos += 2 // consume /*
	for l.pos < len(l.input) {
		if l.input[l.pos] == '*' && l.peek(1) == '/' {
			l.pos += 2
			return
		}
		l.pos++
	}
	// EOF before closing */.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedComment,
	})
}

// tryScanPrefixedLiteral checks whether the identifier-shaped run starting at
// l.pos is actually a string or bytes literal prefix (r/R for raw strings,
// b/B for bytes, and the raw-bytes combinations rb/br/RB/BR/...). If so it
// scans the literal and returns (tok, true). Otherwise it returns
// (Token{}, false) and leaves l.pos unchanged, so the caller falls through to
// identifier/keyword scanning.
//
// GoogleSQLLexer.g4:
//
//	STRING_LITERAL: R? ( SQTEXT | DQTEXT | SQ3TEXT | DQ3TEXT )
//	BYTES_LITERAL:  (B | R B | B R) ( ... )
//
// The prefix is case-insensitive and is immediately followed by a quote.
func (l *Lexer) tryScanPrefixedLiteral() (Token, bool) {
	c0 := lower(l.input[l.pos])
	c1 := lower(l.peek(1))
	c2 := l.peek(2)

	switch {
	// Single-letter raw string prefix: r' r" r''' r"""
	case c0 == 'r' && isQuote(l.peek(1)):
		l.pos++ // consume r
		tok := l.scanString(true)
		tok.Loc.Start = l.start
		return tok, true
	// Single-letter bytes prefix: b' b" ...
	case c0 == 'b' && isQuote(l.peek(1)):
		l.pos++ // consume b
		tok := l.scanBytes(false)
		tok.Loc.Start = l.start
		return tok, true
	// Two-letter raw-bytes prefixes: rb' br' (case-insensitive) followed by a quote.
	case (c0 == 'r' && c1 == 'b' || c0 == 'b' && c1 == 'r') && isQuote(c2):
		l.pos += 2 // consume the two-letter prefix
		tok := l.scanBytes(true)
		tok.Loc.Start = l.start
		return tok, true
	}
	return Token{}, false
}

// lower returns the ASCII-lowercased byte.
func lower(ch byte) byte {
	if ch >= 'A' && ch <= 'Z' {
		return ch + ('a' - 'A')
	}
	return ch
}

// isQuote reports whether ch opens a string/bytes body (single or double quote).
func isQuote(ch byte) bool {
	return ch == '\'' || ch == '"'
}

// scanString reads a string literal body beginning at the opening quote
// (l.pos points at the quote). The optional R prefix has already been consumed
// by the caller; raw indicates whether it was present. l.start must already
// point at the first byte of the whole literal (prefix or quote).
//
// Supports all four quote forms: '...', "...", ”'...”', """...""". Triple-
// quoted strings may span newlines and may contain unescaped single quotes;
// single/double-quoted strings may not span an unescaped newline.
//
// Token.Str is the verbatim body with the surrounding quotes stripped but
// escape sequences PRESERVED (the lexer does not unescape — escape resolution
// is a later concern, matching the ZetaSQL tokenizer which produces a raw
// token image). On an unterminated literal, appends a LexError and returns a
// tokInvalid token covering the bad span.
func (l *Lexer) scanString(raw bool) Token {
	quote := l.input[l.pos]
	triple := l.peek(1) == quote && l.peek(2) == quote
	if triple {
		return l.scanQuotedBody(quote, true, tokString, raw, errUnterminatedTripleString)
	}
	return l.scanQuotedBody(quote, false, tokString, raw, errUnterminatedString)
}

// scanBytes reads a bytes literal body. The b/rb/br prefix has been consumed
// by the caller; raw indicates whether an r was part of the prefix. Mirrors
// scanString but tags the token tokBytes.
func (l *Lexer) scanBytes(raw bool) Token {
	quote := l.input[l.pos]
	triple := l.peek(1) == quote && l.peek(2) == quote
	return l.scanQuotedBody(quote, triple, tokBytes, raw, errUnterminatedBytes)
}

// scanQuotedBody is the shared scanner for single/double, single/triple-quoted
// string and bytes literals. l.pos points at the first opening quote byte.
//
//   - quote is '\” or '"'.
//   - triple selects the ”'...”' / """...""" form.
//   - tokType is tokString or tokBytes.
//   - raw sets Token.IsRaw.
//   - unterminatedMsg is the LexError message on EOF/newline before close.
//
// Escapes (\x, including line continuations \<newline>) are recognized so a
// backslash-escaped closing quote does not terminate the literal — but the
// escape bytes are kept verbatim in Str (no unescaping). In a raw literal a
// backslash is an ordinary byte; however the legacy grammar still treats
// "\\" + quote as a non-terminator (ANY_ESCAPE applies even with the R
// prefix), so escape handling for finding the end is identical for raw and
// non-raw — only the downstream interpretation differs.
func (l *Lexer) scanQuotedBody(quote byte, triple bool, tokType int, raw bool, unterminatedMsg string) Token {
	start := l.start
	closeLen := 1
	if triple {
		closeLen = 3
	}
	l.pos += closeLen // consume the opening quote(s)
	bodyStart := l.pos

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == '\\':
			// Escape: consume the backslash and the next byte (ANY_ESCAPE).
			// This applies to raw literals too for the purpose of locating the
			// terminator (the grammar's R? prefix sits in front of the same
			// SQTEXT/DQTEXT bodies, which use ANY_ESCAPE). At EOF after the
			// backslash, fall through to the unterminated handler below.
			if l.pos+1 >= len(l.input) {
				l.pos++ // consume the trailing backslash
				return l.unterminatedQuoted(start, unterminatedMsg)
			}
			l.pos += 2
		case ch == quote:
			if triple {
				if l.peek(1) == quote && l.peek(2) == quote {
					body := l.input[bodyStart:l.pos]
					l.pos += 3 // consume closing '''
					return Token{Type: tokType, Str: body, IsRaw: raw, Loc: ast.Loc{Start: start, End: l.pos}}
				}
				// A lone quote inside a triple-quoted literal is content.
				l.pos++
			} else {
				body := l.input[bodyStart:l.pos]
				l.pos++ // consume closing quote
				return Token{Type: tokType, Str: body, IsRaw: raw, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		case (ch == '\n' || ch == '\r') && !triple:
			// Single/double-quoted literals cannot span a raw newline.
			return l.unterminatedQuoted(start, unterminatedMsg)
		default:
			l.pos++
		}
	}
	// EOF before the closing quote.
	return l.unterminatedQuoted(start, unterminatedMsg)
}

// unterminatedQuoted records an unterminated-literal error spanning
// [start, l.pos) and returns a tokInvalid token over that span.
func (l *Lexer) unterminatedQuoted(start int, msg string) Token {
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: msg})
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanBacktickIdent reads a `backtick-quoted` identifier (BQTEXT). The body
// may contain any byte except an unescaped backtick, backslash-newline pair,
// or bare newline; escapes (\x) are recognized so a \` does not terminate the
// identifier. Token.Str is the verbatim body with the backticks stripped (no
// unescaping). On an unterminated identifier, appends a LexError and returns a
// tokInvalid token (mapping the grammar's UNCLOSED_ESCAPED_IDENTIFIER).
func (l *Lexer) scanBacktickIdent() Token {
	start := l.start
	l.pos++ // consume opening backtick
	bodyStart := l.pos
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '\\':
			if l.pos+1 >= len(l.input) {
				l.pos++
				return l.unterminatedQuoted(start, errUnterminatedIdentifier)
			}
			l.pos += 2
		case '`':
			body := l.input[bodyStart:l.pos]
			l.pos++ // consume closing backtick
			return Token{Type: tokIdentifier, Str: body, Loc: ast.Loc{Start: start, End: l.pos}}
		case '\n', '\r':
			// Backtick identifiers cannot span a bare newline (BQTEXT_0 body
			// excludes \r and \n).
			return l.unterminatedQuoted(start, errUnterminatedIdentifier)
		default:
			l.pos++
		}
	}
	return l.unterminatedQuoted(start, errUnterminatedIdentifier)
}

// scanNumber reads an integer or floating-point literal.
//
// GoogleSQLLexer.g4:
//
//	INTEGER_LITERAL:        DECIMAL_DIGITS | HEX_DIGITS   (HEX_DIGITS: '0x' [0-9a-f]+)
//	FLOATING_POINT_LITERAL: digits '.' digits? (E [+-]? digits)?
//	                        | digits? '.' digits (E [+-]? digits)?
//	                        | digits E [+-]? digits
//
// Token.Str is the verbatim source text. tokInteger also populates Ival (both
// decimal and 0x-hex). The lexer leaves a leading sign to the parser (the
// grammar's signed forms are assembled there), so this scanner never consumes
// a leading +/-.
func (l *Lexer) scanNumber() Token {
	start := l.start

	// Hex integer: 0x[0-9a-f]+ . Only when at least one hex digit follows;
	// otherwise treat "0" as a decimal integer and let the rest re-lex.
	if l.input[l.pos] == '0' && (l.peek(1) == 'x' || l.peek(1) == 'X') && isHexDigit(l.peek(2)) {
		l.pos += 2 // consume 0x
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		text := l.input[start:l.pos]
		ival, _ := strconv.ParseInt(text[2:], 16, 64)
		return Token{Type: tokInteger, Str: text, Ival: ival, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	isFloat := false

	// Leading-dot form (.5). The caller only routes here when a digit follows
	// the dot, so this consumes at least one fractional digit.
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		l.consumeDigits()
	} else {
		l.consumeDigits() // integer part
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			l.consumeDigits() // optional fractional part
		}
	}

	// Optional exponent: E [+-]? digits. Only an exponent with at least one
	// digit is part of the literal; a trailing E with no digits is not (it
	// would be re-lexed as an identifier).
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		save := l.pos
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		if l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			isFloat = true
			l.consumeDigits()
		} else {
			// Not a valid exponent — back out.
			l.pos = save
		}
	}

	text := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}
	if isFloat {
		return Token{Type: tokFloat, Str: text, Loc: loc}
	}
	ival, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		// Overflow of int64 — keep the verbatim text, leave Ival zero. (The
		// GoogleSQL grammar's INTEGER_LITERAL has no width bound; downstream
		// type analysis handles out-of-range values.)
		return Token{Type: tokInteger, Str: text, Loc: loc}
	}
	return Token{Type: tokInteger, Str: text, Ival: ival, Loc: loc}
}

// consumeDigits advances l.pos past a run of ASCII decimal digits.
func (l *Lexer) consumeDigits() {
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }
func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// scanIdentOrKeyword reads an unquoted identifier run and looks it up in the
// keyword map. Returns the kw* token if found (Str = source text, case
// preserved) or tokIdentifier otherwise. The lexer always returns the keyword
// token for a keyword; the reserved/non-reserved distinction is enforced by
// the parser via IsReservedKeyword.
func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.start
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	text := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}
	if t, ok := KeywordToken(text); ok {
		return Token{Type: t, Str: text, Loc: loc}
	}
	return Token{Type: tokIdentifier, Str: text, Loc: loc}
}

// scanOperatorOrPunct handles single- and multi-char operators and
// punctuation per GoogleSQLLexer.g4. Multi-char tokens are matched with the
// longest-prefix rule:
//
//	!=  <>  <=  >=  <<  >>  ->  =>  +=  -=  |>  ||  @@
//
// Single-char tokens (= < > + - * / % ~ ! ^ & ( ) [ ] { } , ; . : | @ ?) use
// the ASCII byte value as their Type. Bytes that are not valid GoogleSQL
// operators/punctuation produce a tokInvalid token and an invalid-byte
// LexError.
func (l *Lexer) scanOperatorOrPunct(ch byte) Token {
	start := l.start
	switch ch {
	case '!':
		if l.peek(1) == '=' {
			return l.emitMulti(start, 2, tokNotEqual)
		}
		return l.emitSingle(start, ch) // EXCLAMATION_OPERATOR
	case '<':
		switch l.peek(1) {
		case '>':
			return l.emitMulti(start, 2, tokNotEqual2)
		case '=':
			return l.emitMulti(start, 2, tokLessEqual)
		case '<':
			return l.emitMulti(start, 2, tokShiftLeft)
		}
		return l.emitSingle(start, ch)
	case '>':
		switch l.peek(1) {
		case '=':
			return l.emitMulti(start, 2, tokGreaterEqual)
		case '>':
			return l.emitMulti(start, 2, tokShiftRight)
		}
		return l.emitSingle(start, ch)
	case '-':
		switch l.peek(1) {
		case '>':
			return l.emitMulti(start, 2, tokArrow)
		case '=':
			return l.emitMulti(start, 2, tokMinusEqual)
		}
		return l.emitSingle(start, ch)
	case '=':
		if l.peek(1) == '>' {
			return l.emitMulti(start, 2, tokFatArrow)
		}
		return l.emitSingle(start, ch)
	case '+':
		if l.peek(1) == '=' {
			return l.emitMulti(start, 2, tokPlusEqual)
		}
		return l.emitSingle(start, ch)
	case '|':
		switch l.peek(1) {
		case '|':
			return l.emitMulti(start, 2, tokBoolOr)
		case '>':
			return l.emitMulti(start, 2, tokPipe)
		}
		return l.emitSingle(start, ch) // STROKE_SYMBOL '|'
	case '@':
		if l.peek(1) == '@' {
			return l.emitMulti(start, 2, tokAtAt)
		}
		return l.emitSingle(start, ch) // AT_SYMBOL '@'
	case '*', '/', '%', '~', '^', '&', '(', ')', '[', ']', '{', '}', ',', ';', '.', ':', '?':
		return l.emitSingle(start, ch)
	}
	// Unknown byte (e.g. '$', which is not a GoogleSQL token).
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: start + 1}, Msg: errInvalidByte})
	l.pos++
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// emitSingle consumes one byte and returns an ASCII single-char token.
func (l *Lexer) emitSingle(start int, ch byte) Token {
	l.pos++
	return Token{Type: int(ch), Loc: ast.Loc{Start: start, End: l.pos}}
}

// emitMulti consumes n bytes and returns a multi-char operator token.
func (l *Lexer) emitMulti(start, n, tokType int) Token {
	l.pos += n
	return Token{Type: tokType, Loc: ast.Loc{Start: start, End: l.pos}}
}
