package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// Lexer is a Snowflake SQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring of a larger document)
// and call NextToken until it returns Token{Type: tokEOF}. Lex errors
// are collected in a slice (Errors); each error is accompanied by a
// synthetic tokInvalid token at the failure site so consumers can
// choose to halt on the first error or proceed with best-effort parsing.
type Lexer struct {
	input      string
	pos        int // current byte offset (one past the last consumed byte)
	start      int // start byte of the token currently being scanned
	errors     []LexError
	baseOffset int // added to token Loc.Start/End and error Loc.Start/End when returned via public API
}

// NewLexer constructs a Lexer for the given input. Token positions are
// zero-based byte offsets into input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset constructs a Lexer whose emitted token Loc values
// are shifted by baseOffset. Use this when tokenizing a substring of a
// larger document and you want Loc values to refer to positions in the
// larger document rather than the substring. The F4 parser-entry uses
// this to plumb absolute positions through F3's split-then-parse
// pipeline.
//
// The shift is applied lazily when tokens and errors leave the Lexer
// via NextToken() and Errors(). Internal scan helpers continue to use
// local (unshifted) positions against input.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// Errors returns all lex errors collected so far. Errors are appended in
// source order; each error is accompanied by a tokInvalid token in the
// stream.
//
// If the Lexer was constructed with NewLexerWithOffset, the returned
// errors have Loc values shifted by baseOffset. The internal l.errors
// slice retains unshifted local positions.
func (l *Lexer) Errors() []LexError {
	if l.baseOffset == 0 {
		return l.errors
	}
	shifted := make([]LexError, len(l.errors))
	for i, e := range l.errors {
		shifted[i] = LexError{
			Loc: ast.Loc{
				Start: e.Loc.Start + l.baseOffset,
				End:   e.Loc.End + l.baseOffset,
			},
			Msg: e.Msg,
		}
	}
	return shifted
}

// Tokenize is a one-shot convenience that runs a lexer to EOF and returns
// the full token stream and error list. Useful for tests and for callers
// that don't need streaming.
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

// NextToken advances the lexer and returns the next token. At end of input
// it returns Token{Type: tokEOF}; subsequent calls continue to return EOF.
//
// If the Lexer was constructed with NewLexerWithOffset, the returned
// token's Loc values are shifted by baseOffset.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	if l.baseOffset != 0 {
		tok.Loc.Start += l.baseOffset
		tok.Loc.End += l.baseOffset
	}
	return tok
}

// nextTokenInner returns the next token with LOCAL (unshifted) positions.
// It is the original body of the lexer dispatch — the public NextToken
// above is a thin wrapper that applies baseOffset before returning.
func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespaceAndComments()
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
	case ch == '$':
		return l.scanDollar()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.peekDigit():
		return l.scanNumber()
	case (ch == 'X' || ch == 'x') && l.peek(1) == '\'':
		l.pos++ // consume X
		tok := l.scanString()
		tok.XPrefix = true
		tok.Loc.Start = l.start
		return tok
	case isIdentStart(ch):
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

// peekDigit reports whether the byte at l.pos+1 is an ASCII digit.
func (l *Lexer) peekDigit() bool {
	c := l.peek(1)
	return c >= '0' && c <= '9'
}

// isIdentStart reports whether ch may begin an unquoted identifier.
// Snowflake's ID rule is [A-Za-z_]; the case-insensitive flag in the
// legacy grammar already covers both cases.
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isIdentCont reports whether ch may continue an unquoted identifier.
// Snowflake's ID rule allows letters, digits, underscore, @, and $.
func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '@' || ch == '$'
}

// skipWhitespaceAndComments advances l.pos past any whitespace, line
// comments (-- or //), and block comments (/* ... */ with nesting).
// Comments are dropped silently; whitespace produces no token.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == ' ', ch == '\t', ch == '\n', ch == '\r':
			l.pos++
		case ch == '-' && l.peek(1) == '-':
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '/':
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '*':
			l.scanBlockComment()
		default:
			return
		}
	}
}

// scanBlockComment advances past a /* ... */ block comment. Block comments
// nest in Snowflake (the legacy grammar uses recursive
// '/*' (SQL_COMMENT | .)*? '*/'). On unterminated comment, appends a
// LexError but does NOT emit a token (comments are channel-HIDDEN).
func (l *Lexer) scanBlockComment() {
	start := l.pos
	l.pos += 2 // consume /*
	depth := 1
	for l.pos < len(l.input) && depth > 0 {
		if l.input[l.pos] == '/' && l.peek(1) == '*' {
			depth++
			l.pos += 2
		} else if l.input[l.pos] == '*' && l.peek(1) == '/' {
			depth--
			l.pos += 2
		} else {
			l.pos++
		}
	}
	if depth > 0 {
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: l.pos},
			Msg: errUnterminatedComment,
		})
	}
}

// scanString reads a single-quoted string literal. Snowflake supports two
// escape mechanisms inside '...':
//   - backslash escapes: \n \t \r \0 \\ \' \"  (and \<other> → <other>)
//   - doubled-quote escape: ” inside a string is a literal '
//
// Token.Str contains the unescaped content (no surrounding quotes).
// On unterminated string, appends an unterminated-string LexError, emits
// a tokInvalid token covering the bad span, and advances to the next
// newline or EOF.
func (l *Lexer) scanString() Token {
	start := l.start // l.start is the byte offset of the opening quote
	l.pos++          // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.peek(1) == '\'' {
				// Doubled-quote escape: '' → literal '
				sb.WriteByte('\'')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokString,
				Str:  sb.String(),
				Loc:  ast.Loc{Start: start, End: l.pos},
			}
		}
		if ch == '\\' {
			l.pos++
			if l.pos >= len(l.input) {
				break
			}
			esc := l.input[l.pos]
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '0':
				sb.WriteByte(0)
			case '\\':
				sb.WriteByte('\\')
			case '\'':
				sb.WriteByte('\'')
			case '"':
				sb.WriteByte('"')
			case 'b':
				sb.WriteByte(0x08)
			default:
				sb.WriteByte(esc)
			}
			l.pos++
			continue
		}
		if ch == '\n' {
			// Unterminated single-quoted string — strings cannot span newlines.
			l.errors = append(l.errors, LexError{
				Loc: ast.Loc{Start: start, End: l.pos},
				Msg: errUnterminatedString,
			})
			return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	// EOF before closing quote.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedString,
	})
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanQuotedIdent reads a "quoted identifier". Snowflake's quoted identifiers
// preserve case and can contain any character except an unescaped ". The
// only escape is "" → literal ". The empty form "" is also a valid quoted
// identifier (with Str = "").
//
// Token.Str contains the unescaped content (no surrounding quotes).
// On unterminated quoted identifier, appends an unterminated-quoted-identifier
// LexError and emits a tokInvalid token.
func (l *Lexer) scanQuotedIdent() Token {
	start := l.start
	l.pos++ // consume opening "
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			if l.peek(1) == '"' {
				// Doubled-quote escape: "" → literal "
				sb.WriteByte('"')
				l.pos += 2
				continue
			}
			// Closing quote.
			l.pos++
			return Token{
				Type: tokQuotedIdent,
				Str:  sb.String(),
				Loc:  ast.Loc{Start: start, End: l.pos},
			}
		}
		if ch == '\n' {
			l.errors = append(l.errors, LexError{
				Loc: ast.Loc{Start: start, End: l.pos},
				Msg: errUnterminatedQuoted,
			})
			return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedQuoted,
	})
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanDollar handles tokens that begin with '$':
//   - $$...$$ → tokString with raw content (no escapes)
//   - $NAME   → tokVariable with Str = "NAME" (no leading $)
//   - $       → ASCII single-char token (Type = int('$'))
//
// The $$ form must be checked first so $ followed by $... is treated as a
// dollar string opener rather than a variable.
func (l *Lexer) scanDollar() Token {
	start := l.start

	// $$...$$ form
	if l.peek(1) == '$' {
		l.pos += 2 // consume opening $$
		contentStart := l.pos
		for l.pos < len(l.input)-1 {
			if l.input[l.pos] == '$' && l.input[l.pos+1] == '$' {
				content := l.input[contentStart:l.pos]
				l.pos += 2 // consume closing $$
				return Token{
					Type: tokString,
					Str:  content,
					Loc:  ast.Loc{Start: start, End: l.pos},
				}
			}
			l.pos++
		}
		// Unterminated $$ string. Advance to EOF.
		l.pos = len(l.input)
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: l.pos},
			Msg: errUnterminatedDollar,
		})
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// $NAME form: $ followed by [A-Z0-9_]+
	// Snowflake's ID2 rule: '$' [A-Z0-9_]*. Note that an empty $ (with no
	// follow-on identifier characters) falls through to the bare-$ case below.
	l.pos++ // consume $
	nameStart := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
			l.pos++
			continue
		}
		break
	}
	if l.pos > nameStart {
		return Token{
			Type: tokVariable,
			Str:  l.input[nameStart:l.pos],
			Loc:  ast.Loc{Start: start, End: l.pos},
		}
	}

	// Bare $ — emit as ASCII single-char token.
	return Token{
		Type: int('$'),
		Loc:  ast.Loc{Start: start, End: l.pos},
	}
}

// scanNumber reads a numeric literal and returns one of three token kinds:
//   - tokInt   for digits-only (e.g. 42, 0)
//   - tokFloat for digits with a dot (e.g. 1.5, 1., .5)
//   - tokReal  for digits with an E exponent (e.g. 1e10, 1.5e-10)
//
// All three kinds populate Token.Str with the verbatim source text. tokInt
// additionally populates Token.Ival via strconv.ParseInt; if the integer
// overflows int64, the token is downgraded to tokFloat (still preserving
// Str) — this matches Snowflake's NUMBER(38, 0) arbitrary-precision behavior.
func (l *Lexer) scanNumber() Token {
	start := l.start
	isFloat := false

	// Leading dot form (.5).
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	} else {
		// Integer part.
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
		// Optional decimal point.
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
				l.pos++
			}
		}
	}

	// Optional exponent.
	isReal := false
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isReal = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	text := l.input[start:l.pos]
	loc := ast.Loc{Start: start, End: l.pos}

	if isReal {
		return Token{Type: tokReal, Str: text, Loc: loc}
	}
	if isFloat {
		return Token{Type: tokFloat, Str: text, Loc: loc}
	}

	// Integer. Try strconv.ParseInt; on overflow, downgrade to tokFloat.
	ival, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		// Overflow or other parse failure. Treat as tokFloat (Snowflake
		// NUMBER(38, 0) accepts arbitrary precision).
		return Token{Type: tokFloat, Str: text, Loc: loc}
	}
	return Token{Type: tokInt, Str: text, Ival: ival, Loc: loc}
}

// scanIdentOrKeyword reads an identifier-shaped run and looks it up in the
// keyword map. Returns kw* if found (with Str = source text, preserving
// case) or tokIdent otherwise.
//
// Snowflake's ID rule is [A-Za-z_][A-Za-z0-9_@$]*. The legacy grammar's
// case-insensitive flag is honored by KeywordToken (which lowercases for
// lookup).
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
	return Token{Type: tokIdent, Str: text, Loc: loc}
}

// scanOperatorOrPunct handles single- and multi-char operators and
// punctuation. Multi-char operators (with their first byte and lookahead):
//
//	:: → tokDoubleColon
//	|| → tokConcat
//	-> → tokArrow
//	->>  → tokFlow (must be checked before -> via 2-byte lookahead)
//	=> → tokAssoc
//	!= → tokNotEq
//	<> → tokNotEq
//	<= → tokLessEq
//	>= → tokGreaterEq
//
// Single-char tokens use the ASCII byte value as their Type. Bytes that
// are not valid Snowflake operators or punctuation produce a tokInvalid
// token and an invalid-byte LexError.
func (l *Lexer) scanOperatorOrPunct(ch byte) Token {
	start := l.start
	switch ch {
	case ':':
		if l.peek(1) == ':' {
			l.pos += 2
			return Token{Type: tokDoubleColon, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int(':'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '|':
		if l.peek(1) == '|' {
			l.pos += 2
			return Token{Type: tokConcat, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		// Bare | is not a valid Snowflake operator.
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: start + 1},
			Msg: errInvalidByte,
		})
		l.pos++
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	case '-':
		if l.peek(1) == '>' {
			if l.peek(2) == '>' {
				l.pos += 3
				return Token{Type: tokFlow, Loc: ast.Loc{Start: start, End: l.pos}}
			}
			l.pos += 2
			return Token{Type: tokArrow, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('-'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '=':
		if l.peek(1) == '>' {
			l.pos += 2
			return Token{Type: tokAssoc, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('='), Loc: ast.Loc{Start: start, End: l.pos}}
	case '!':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.errors = append(l.errors, LexError{
			Loc: ast.Loc{Start: start, End: start + 1},
			Msg: errInvalidByte,
		})
		l.pos++
		return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	case '<':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		if l.peek(1) == '>' {
			l.pos += 2
			return Token{Type: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('<'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '>':
		if l.peek(1) == '=' {
			l.pos += 2
			return Token{Type: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		l.pos++
		return Token{Type: int('>'), Loc: ast.Loc{Start: start, End: l.pos}}
	case '+', '*', '/', '%', '~', '(', ')', '[', ']', '{', '}', ',', ';', '.', '@':
		l.pos++
		return Token{Type: int(ch), Loc: ast.Loc{Start: start, End: l.pos}}
	}
	// Unknown byte.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: start + 1},
		Msg: errInvalidByte,
	})
	l.pos++
	return Token{Type: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
