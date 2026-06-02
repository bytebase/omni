package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// Lexer is a Trino SQL tokenizer. Construct via NewLexer (or NewLexerWithOffset
// when tokenizing a substring of a larger document) and call NextToken until it
// returns Token{Kind: tokEOF}. Lexing is byte-oriented; multibyte UTF-8 runes
// are permitted inside quoted identifiers and string literals, but an *unquoted*
// identifier is ASCII-only (matching TrinoLexer.g4's LETTER_ = [A-Z]).
//
// Lex errors are accumulated in Errors(); each error is reported alongside a
// tokInvalid token so the statement splitter and parser can either halt or
// proceed with best-effort recovery (Trino's own grammar keeps an UNRECOGNIZED_
// catch-all for exactly this reason). The lexer never panics on malformed
// input.
type Lexer struct {
	input      string
	pos        int // current byte offset
	start      int // start byte of the token currently being scanned
	errors     []LexError
	baseOffset int // added to every emitted Loc (and to error Locs)
}

// NewLexer creates a lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset creates a lexer whose emitted Loc values are shifted by
// baseOffset. Used when tokenizing one statement carved out of a multi-statement
// document so positions stay anchored to the original text.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// Errors returns all lex errors collected so far, with positions shifted by
// baseOffset when applicable.
func (l *Lexer) Errors() []LexError {
	if l.baseOffset == 0 {
		return l.errors
	}
	shifted := make([]LexError, len(l.errors))
	for i, e := range l.errors {
		shifted[i] = LexError{
			Msg: e.Msg,
			Loc: ast.Loc{Start: e.Loc.Start + l.baseOffset, End: e.Loc.End + l.baseOffset},
		}
	}
	return shifted
}

// Tokenize is a one-shot convenience that lexes the entire input, returning all
// tokens (terminated by a tokEOF) and any lex errors.
func Tokenize(input string) ([]Token, []LexError) {
	l := NewLexer(input)
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind == tokEOF {
			break
		}
	}
	return tokens, l.Errors()
}

// NextToken returns the next token, with positions shifted by baseOffset.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	if l.baseOffset != 0 {
		tok.Loc.Start += l.baseOffset
		tok.Loc.End += l.baseOffset
	}
	return tok
}

func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Kind: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
	}
	l.start = l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString()
	case ch == '"':
		return l.scanQuotedIdent()
	case ch == '`':
		return l.scanBackquotedIdent()
	case (ch == 'U' || ch == 'u') && l.peek(1) == '&' && l.peek(2) == '\'':
		return l.scanUnicodeString()
	case (ch == 'X' || ch == 'x') && l.peek(1) == '\'':
		return l.scanBinaryLiteral()
	case ch >= '0' && ch <= '9':
		return l.scanNumberOrDigitIdent()
	case ch == '.' && isDigit(l.peek(1)):
		return l.scanNumberOrDigitIdent()
	case isIdentStart(ch):
		return l.scanIdentOrKeyword()
	}
	return l.scanOperator(ch)
}

// --- Helpers ---------------------------------------------------------------

func (l *Lexer) peek(offset int) byte {
	if l.pos+offset >= len(l.input) {
		return 0
	}
	return l.input[l.pos+offset]
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

// isIdentStart reports whether ch may begin an unquoted identifier.
//
// TrinoLexer.g4: IDENTIFIER_ = (LETTER_ | '_') ... with LETTER_ = [A-Z]
// under caseInsensitive = true, i.e. ASCII letters and '_' only. Non-ASCII
// bytes are deliberately NOT admitted: Trino (and the legacy grammar) route
// them to the UNRECOGNIZED_ catch-all, so an unquoted non-ASCII identifier such
// as `é` is a lex error in both truth1 and truth2 — the identifier must be
// double-quoted. Admitting high bytes here would be an unadjudicated divergence
// that wrongly accepts SQL Trino rejects.
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isIdentCont reports whether ch may continue an unquoted identifier.
func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}

func (l *Lexer) addError(msg string, start, end int) {
	l.errors = append(l.errors, LexError{Msg: msg, Loc: ast.Loc{Start: start, End: end}})
}

// --- Whitespace & Comments -------------------------------------------------

// skipWhitespaceAndComments consumes the hidden-channel tokens of
// TrinoLexer.g4: WS_ ([ \r\n\t]+), SIMPLE_COMMENT_ ('--' to end of line), and
// BRACKETED_COMMENT_ ('/*' ... '*/', non-nesting). Unlike MySQL/Doris, Trino's
// '--' starts a comment regardless of the following character, and bracketed
// comments do not nest ('.*?' is non-greedy, stopping at the first '*/').
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			l.pos++
		case ch == '-' && l.peek(1) == '-':
			// SIMPLE_COMMENT_: -- to end of line (the newline is consumed too).
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' && l.input[l.pos] != '\r' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '*':
			l.scanBracketedComment()
		default:
			return
		}
	}
}

// scanBracketedComment consumes a non-nesting '/* ... */' comment.
func (l *Lexer) scanBracketedComment() {
	start := l.pos
	l.pos += 2 // consume /*
	for l.pos+1 < len(l.input) {
		if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
			l.pos += 2
			return
		}
		l.pos++
	}
	// Reached EOF without a closing */.
	l.pos = len(l.input)
	l.addError(errUnterminatedComment, start, l.pos)
}

// --- String / Identifier scanning ------------------------------------------

// scanString reads a single-quoted string literal (STRING_). Trino is
// standard-conforming: the only in-string escape is a doubled quote (”) which
// denotes one literal quote; backslashes are ordinary characters. Newlines are
// permitted inside the literal.
func (l *Lexer) scanString() Token {
	start := l.start
	l.pos++ // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.peek(1) == '\'' {
				sb.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: tokString, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.addError(errUnterminatedString, start, l.pos)
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanUnicodeString reads a U&'...' unicode string literal (UNICODE_STRING_).
// The leading U& has already been confirmed by the caller. The optional
// "UESCAPE 'c'" suffix and the \XXXX / \+XXXXXX escapes are part of the
// parser's string-rule semantics (string_ -> unicodeStringLiteral), not the
// lexer; the lexer captures the raw inter-quote bytes verbatim, with ” meaning
// a literal quote (consistent with STRING_).
func (l *Lexer) scanUnicodeString() Token {
	start := l.start
	l.pos += 2 // consume U&
	l.pos++    // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.peek(1) == '\'' {
				sb.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: tokUnicodeString, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.addError(errUnterminatedString, start, l.pos)
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanBinaryLiteral reads an X'...' binary literal (BINARY_LITERAL_). Per the
// grammar comment, any byte except a quote is admitted between the quotes; the
// content is validated as hex when the AST is constructed, so a more
// descriptive error can be produced there. There is no ” escape inside X'...'.
func (l *Lexer) scanBinaryLiteral() Token {
	start := l.start
	l.pos += 2 // consume X and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		l.addError(errUnterminatedBinary, start, l.pos)
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}
	content := l.input[contentStart:l.pos]
	l.pos++ // consume closing quote
	return Token{Kind: tokBinaryLiteral, Str: content, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanQuotedIdent reads a "..." identifier (QUOTED_IDENTIFIER_). A doubled
// double-quote ("") denotes one literal double-quote in the identifier.
func (l *Lexer) scanQuotedIdent() Token {
	return l.scanDelimitedIdent('"', tokQuotedIdent)
}

// scanBackquotedIdent reads a `...` identifier (BACKQUOTED_IDENTIFIER_). A
// doubled backtick (“) denotes one literal backtick.
func (l *Lexer) scanBackquotedIdent() Token {
	return l.scanDelimitedIdent('`', tokBackquotedIdent)
}

func (l *Lexer) scanDelimitedIdent(quote byte, kind TokenKind) Token {
	start := l.start
	l.pos++ // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == quote {
			if l.peek(1) == quote {
				sb.WriteByte(quote)
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: kind, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.addError(errUnterminatedQuoted, start, l.pos)
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanIdentOrKeyword reads an unquoted identifier (IDENTIFIER_) and resolves it
// to a keyword token when it matches one (case-insensitively). The original
// source text is preserved in Str for both kinds.
func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.start
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	word := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}
	if kw, ok := KeywordToken(word); ok {
		return Token{Kind: kw, Str: word, Loc: loc}
	}
	return Token{Kind: tokIdent, Str: word, Loc: loc}
}

// --- Numbers / digit-leading identifiers -----------------------------------

// scanNumberOrDigitIdent lexes a run beginning with a digit (or a dot followed
// by a digit) into one of INTEGER_VALUE_, DECIMAL_VALUE_, DOUBLE_VALUE_, or
// DIGIT_IDENTIFIER_, faithfully reproducing ANTLR's maximal-munch + rule-order
// disambiguation:
//
//	INTEGER_VALUE_:   DIGIT+
//	DECIMAL_VALUE_:   DIGIT+ '.' DIGIT* | '.' DIGIT+
//	DOUBLE_VALUE_:    DIGIT+ ('.' DIGIT*)? EXPONENT | '.' DIGIT+ EXPONENT
//	DIGIT_IDENTIFIER_: DIGIT (LETTER | DIGIT | '_')+
//
// ANTLR picks the rule that consumes the longest input; on a tie the
// earliest-declared rule wins (numbers precede DIGIT_IDENTIFIER_ in the
// grammar). Concretely: a pure-digit run is INTEGER_VALUE_; a digit run with a
// fractional/exponent part is DECIMAL/DOUBLE; but if a letter or underscore
// follows in a way the numeric rules cannot absorb, the whole run is a
// DIGIT_IDENTIFIER_ because that rule then matches more characters.
func (l *Lexer) scanNumberOrDigitIdent() Token {
	start := l.start

	// First, compute how far the numeric rules can reach.
	numEnd, numKind := l.scanNumericExtent(start)

	// Then, compute how far DIGIT_IDENTIFIER_ can reach. It requires a leading
	// digit (not a leading dot) and at least one trailing LETTER/DIGIT/'_'.
	idEnd := -1
	if isDigit(l.input[start]) {
		j := start
		for j < len(l.input) && (isIdentCont(l.input[j])) {
			j++
		}
		// Must contain at least one char beyond the first to satisfy the '+'.
		if j > start+1 {
			idEnd = j
		}
	}

	// Maximal munch: longer match wins; tie goes to the numeric rule (declared
	// first in the grammar).
	if idEnd > numEnd {
		l.pos = idEnd
		return Token{Kind: tokDigitIdent, Str: l.input[start:idEnd], Loc: ast.Loc{Start: start, End: idEnd}}
	}

	l.pos = numEnd
	text := l.input[start:numEnd]
	loc := ast.Loc{Start: start, End: numEnd}
	switch numKind {
	case tokInteger:
		val, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			// Overflow: still a valid INTEGER_VALUE_ token to the grammar; the
			// numeric value is left at 0 and validated downstream if needed.
			return Token{Kind: tokInteger, Str: text, Loc: loc}
		}
		return Token{Kind: tokInteger, Str: text, Ival: val, Loc: loc}
	case tokDecimal:
		return Token{Kind: tokDecimal, Str: text, Loc: loc}
	default: // tokDouble
		return Token{Kind: tokDouble, Str: text, Loc: loc}
	}
}

// scanNumericExtent computes, without moving l.pos, the end offset and kind of
// the longest INTEGER/DECIMAL/DOUBLE match starting at start. The caller has
// guaranteed input[start] is a digit, or a '.' immediately followed by a digit.
func (l *Lexer) scanNumericExtent(start int) (end int, kind TokenKind) {
	i := start
	n := len(l.input)
	kind = tokInteger

	if l.input[i] == '.' {
		// '.' DIGIT+ form — fractional, no integer part.
		i++ // consume '.'
		for i < n && isDigit(l.input[i]) {
			i++
		}
		kind = tokDecimal
	} else {
		// DIGIT+ integer part.
		for i < n && isDigit(l.input[i]) {
			i++
		}
		// Optional '.' DIGIT* fractional part. A '.' followed by another '.'
		// (e.g. an unlikely '1..') still consumes a single '.' here, matching
		// DECIMAL_VALUE_'s 'DIGIT+ "." DIGIT*'.
		if i < n && l.input[i] == '.' {
			i++ // consume '.'
			for i < n && isDigit(l.input[i]) {
				i++
			}
			kind = tokDecimal
		}
	}

	// Optional exponent: EXPONENT_ = 'E' [+-]? DIGIT+. The exponent only counts
	// if it is well-formed (at least one digit after E and the optional sign);
	// otherwise the numeric match stops before the 'E', which lets a trailing
	// 'E...' be absorbed by DIGIT_IDENTIFIER_ instead.
	if i < n && (l.input[i] == 'e' || l.input[i] == 'E') {
		j := i + 1
		if j < n && (l.input[j] == '+' || l.input[j] == '-') {
			j++
		}
		if j < n && isDigit(l.input[j]) {
			for j < n && isDigit(l.input[j]) {
				j++
			}
			i = j
			kind = tokDouble
		}
	}
	return i, kind
}

// --- Operators / punctuation ----------------------------------------------

// scanOperator lexes operator and punctuation tokens. Single-byte punctuation
// is returned with Kind == int(byte); multi-byte operators use the named tok*
// constants. An unrecognized byte yields tokInvalid plus a LexError, mirroring
// TrinoLexer.g4's UNRECOGNIZED_ catch-all (which exists so the splitter can
// recover).
func (l *Lexer) scanOperator(ch byte) Token {
	start := l.start
	l.pos++

	switch ch {
	case '<':
		switch l.peek(0) {
		case '>':
			l.pos++
			return l.op(tokNotEq, start) // <>
		case '=':
			l.pos++
			return l.op(tokLessEq, start) // <=
		case '-':
			l.pos++
			return l.op(tokLeftArrow, start) // <-
		}
		return l.op(int('<'), start)
	case '>':
		if l.peek(0) == '=' {
			l.pos++
			return l.op(tokGreaterEq, start) // >=
		}
		return l.op(int('>'), start)
	case '!':
		if l.peek(0) == '=' {
			l.pos++
			return l.op(tokNotEq, start) // !=
		}
		// Bare '!' is not a Trino token; fall through to the catch-all.
		l.addError(errUnknownChar, start, l.pos)
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	case '|':
		if l.peek(0) == '|' {
			l.pos++
			return l.op(tokConcat, start) // ||
		}
		return l.op(int('|'), start) // VBAR_
	case '=':
		if l.peek(0) == '>' {
			l.pos++
			return l.op(tokDoubleArrow, start) // =>
		}
		return l.op(int('='), start)
	case '-':
		if l.peek(0) == '>' {
			l.pos++
			return l.op(tokArrow, start) // ->
		}
		if l.peek(0) == '}' {
			l.pos++
			return l.op(tokRCurlyHyphen, start) // -}
		}
		return l.op(int('-'), start)
	case '{':
		if l.peek(0) == '-' {
			l.pos++
			return l.op(tokLCurlyHyphen, start) // {-
		}
		return l.op(int('{'), start)
	case '?':
		return Token{Kind: tokQuestion, Loc: ast.Loc{Start: start, End: l.pos}} // QUESTION_MARK_
	}

	// Single-byte operators and punctuation from TrinoLexer.g4:
	// EQ_ '=', LT_ '<', GT_ '>', PLUS_ '+', MINUS_ '-', ASTERISK_ '*',
	// SLASH_ '/', PERCENT_ '%', SEMICOLON_ ';', DOT_ '.', COLON_ ':',
	// COMMA_ ',', LPAREN_ '(', RPAREN_ ')', LSQUARE_ '[', RSQUARE_ ']',
	// LCURLY_ '{', RCURLY_ '}', VBAR_ '|', DOLLAR_ '$', CARET_ '^'.
	//
	// NOTE on ':' — the legacy grammar declares COLON_ : '_:', a known
	// copy/paste bug: the intended COLON_ lexeme is a single ':' (the SqlBase
	// colon used for JSON_OBJECT key:value pairs and MATCH_RECOGNIZE labels).
	// The two-character '_:' literal is the wrong token shape (it would only ever
	// match the byte pair "_:" , never a lone colon). The migration divergence
	// ledger adjudicated the fix to ':'; this lexer emits ':' as COLON_
	// (Kind == int(':')). See divergence record trino/lexer "COLON_ token literal".
	switch ch {
	case '+', '*', '/', '%', ';', '.', ':', ',', '(', ')', '[', ']', '}', '$', '^':
		return l.op(int(ch), start)
	}

	l.addError(errUnknownChar, start, l.pos)
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// op builds an operator/punctuation token spanning [start, l.pos).
func (l *Lexer) op(kind TokenKind, start int) Token {
	return Token{Kind: kind, Loc: ast.Loc{Start: start, End: l.pos}}
}
