package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// Lexer is a Doris SQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring) and call NextToken
// until it returns Token{Kind: tokEOF}. Lex errors are collected in
// Errors; each error is accompanied by a tokInvalid token so consumers
// can halt or proceed with best-effort parsing.
type Lexer struct {
	input      string
	pos        int // current byte offset
	start      int // start byte of current token
	errors     []LexError
	baseOffset int // added to Loc positions on return
}

// NewLexer creates a lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NewLexerWithOffset creates a lexer whose emitted Loc values are shifted
// by baseOffset. Used for multi-statement parsing.
func NewLexerWithOffset(input string, baseOffset int) *Lexer {
	return &Lexer{input: input, baseOffset: baseOffset}
}

// Errors returns all lex errors collected so far, with positions shifted
// by baseOffset if applicable.
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

// Tokenize is a one-shot convenience that lexes the entire input.
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

// NextToken returns the next token with baseOffset-shifted positions.
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
		return l.scanString('\'')
	case ch == '"':
		return l.scanString('"')
	case ch == '`':
		return l.scanBacktickIdent()
	case ch >= '0' && ch <= '9':
		return l.scanNumber()
	case ch == '.' && l.peekDigit():
		return l.scanNumber()
	case (ch == 'X' || ch == 'x') && l.peek(1) == '\'':
		return l.scanHexLiteral()
	case (ch == 'B' || ch == 'b') && l.peek(1) == '\'':
		return l.scanBitLiteral()
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

func (l *Lexer) peekDigit() bool {
	c := l.peek(1)
	return c >= '0' && c <= '9'
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' || ch > 127
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

// --- Whitespace & Comments -------------------------------------------------

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			l.pos++
		case ch == '-' && l.peek(1) == '-' && l.isLineCommentAfterDash():
			// MySQL-compatible: -- must be followed by space/tab/newline/EOF
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '/':
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '#':
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		case ch == '/' && l.peek(1) == '*' && l.peek(2) != '+':
			// Block comment (but NOT hint /*+)
			l.scanBlockComment()
		default:
			return
		}
	}
}

// isLineCommentAfterDash checks that the byte after -- is space/tab/newline/EOF.
func (l *Lexer) isLineCommentAfterDash() bool {
	c := l.peek(2)
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == 0
}

func (l *Lexer) scanBlockComment() {
	start := l.pos
	l.pos += 2 // consume /*
	for l.pos < len(l.input)-1 {
		if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
			l.pos += 2
			return
		}
		l.pos++
	}
	l.pos = len(l.input)
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnterminatedComment,
	})
}

// --- String Scanning -------------------------------------------------------

// scanString reads a single-quoted or double-quoted string literal.
// Supports backslash escapes and doubled-quote escapes.
// Unlike Snowflake, Doris allows newlines inside string literals.
func (l *Lexer) scanString(quote byte) Token {
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
			return Token{Kind: tokString, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
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
			case 'Z':
				sb.WriteByte(0x1A)
			default:
				sb.WriteByte(esc)
			}
			l.pos++
			continue
		}
		// Newlines are allowed inside Doris string literals.
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// scanBacktickIdent reads a backtick-quoted identifier. Only doubled-backtick
// escape is supported (`` → `). No backslash escapes.
func (l *Lexer) scanBacktickIdent() Token {
	start := l.start
	l.pos++ // consume opening backtick
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '`' {
			if l.peek(1) == '`' {
				sb.WriteByte('`')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Kind: tokQuotedIdent, Str: sb.String(), Loc: ast.Loc{Start: start, End: l.pos}}
		}
		sb.WriteByte(ch)
		l.pos++
	}
	l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedQuoted})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}

// --- Hex/Bit Literal Scanning -----------------------------------------------

func (l *Lexer) scanHexLiteral() Token {
	start := l.start
	l.pos += 2 // consume X and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}
	content := l.input[contentStart:l.pos]
	l.pos++ // consume closing quote
	return Token{Kind: tokHexLiteral, Str: content, Loc: ast.Loc{Start: start, End: l.pos}}
}

func (l *Lexer) scanBitLiteral() Token {
	start := l.start
	l.pos += 2 // consume B and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\'' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		l.errors = append(l.errors, LexError{Loc: ast.Loc{Start: start, End: l.pos}, Msg: errUnterminatedString})
		return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
	}
	content := l.input[contentStart:l.pos]
	l.pos++ // consume closing quote
	return Token{Kind: tokBitLiteral, Str: content, Loc: ast.Loc{Start: start, End: l.pos}}
}

// --- Number Scanning -------------------------------------------------------

func (l *Lexer) scanNumber() Token {
	start := l.start

	// Hex: 0x or 0X
	if l.input[l.pos] == '0' && (l.peek(1) == 'x' || l.peek(1) == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		text := l.input[start:l.pos]
		val, err := strconv.ParseInt(text[2:], 16, 64)
		if err != nil {
			return Token{Kind: tokFloat, Str: text, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		return Token{Kind: tokInt, Str: text, Ival: val, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Binary: 0b or 0B (only when not followed by a quote, which would be a bit literal)
	if l.input[l.pos] == '0' && (l.peek(1) == 'b' || l.peek(1) == 'B') && l.peek(2) != '\'' {
		l.pos += 2
		for l.pos < len(l.input) && (l.input[l.pos] == '0' || l.input[l.pos] == '1') {
			l.pos++
		}
		text := l.input[start:l.pos]
		val, err := strconv.ParseInt(text[2:], 2, 64)
		if err != nil {
			return Token{Kind: tokFloat, Str: text, Loc: ast.Loc{Start: start, End: l.pos}}
		}
		return Token{Kind: tokInt, Str: text, Ival: val, Loc: ast.Loc{Start: start, End: l.pos}}
	}

	isFloat := false

	// Leading dot (.5)
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	} else {
		// Integer part
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
		// Optional decimal
		if l.pos < len(l.input) && l.input[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
				l.pos++
			}
		}
	}

	// Optional exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
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

	if isFloat {
		return Token{Kind: tokFloat, Str: text, Loc: loc}
	}
	val, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return Token{Kind: tokFloat, Str: text, Loc: loc}
	}
	return Token{Kind: tokInt, Str: text, Ival: val, Loc: loc}
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// --- Identifier / Keyword Scanning -----------------------------------------

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

// --- Operator Scanning -----------------------------------------------------

func (l *Lexer) scanOperator(ch byte) Token {
	start := l.start
	l.pos++

	switch ch {
	case '<':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				if l.peek(1) == '>' {
					// <=> null-safe equal
					l.pos += 2
					return Token{Kind: tokNullSafeEq, Loc: ast.Loc{Start: start, End: l.pos}}
				}
				l.pos++
				return Token{Kind: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				l.pos++
				return Token{Kind: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '<':
				l.pos++
				return Token{Kind: tokShiftLeft, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '>':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				l.pos++
				return Token{Kind: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				l.pos++
				return Token{Kind: tokShiftRight, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '!':
		if l.pos < len(l.input) {
			switch l.input[l.pos] {
			case '=':
				l.pos++
				return Token{Kind: tokNotEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '<':
				// !< maps to >= per DorisLexer.g4: GTE: '>=' | '!<'
				l.pos++
				return Token{Kind: tokGreaterEq, Loc: ast.Loc{Start: start, End: l.pos}}
			case '>':
				// !> maps to <= per DorisLexer.g4: LTE: '<=' | '!>'
				l.pos++
				return Token{Kind: tokLessEq, Loc: ast.Loc{Start: start, End: l.pos}}
			}
		}
	case '&':
		if l.pos < len(l.input) && l.input[l.pos] == '&' {
			l.pos++
			return Token{Kind: tokLogicalAnd, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '|':
		if l.pos < len(l.input) && l.input[l.pos] == '|' {
			l.pos++
			return Token{Kind: tokDoublePipes, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '@':
		if l.pos < len(l.input) && l.input[l.pos] == '@' {
			l.pos++
			return Token{Kind: tokDoubleAt, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case ':':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
			return Token{Kind: tokAssign, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '.':
		if l.pos+1 < len(l.input) && l.input[l.pos] == '.' && l.input[l.pos+1] == '.' {
			l.pos += 2
			return Token{Kind: tokDotDotDot, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '-':
		if l.pos < len(l.input) && l.input[l.pos] == '>' {
			l.pos++
			return Token{Kind: tokArrow, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '/':
		// /*+ starts a hint (check for * then + at positions pos and pos+1)
		if l.pos+1 < len(l.input) && l.input[l.pos] == '*' && l.input[l.pos+1] == '+' {
			l.pos += 2
			return Token{Kind: tokHintStart, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '*':
		if l.pos < len(l.input) && l.input[l.pos] == '/' {
			l.pos++
			return Token{Kind: tokHintEnd, Loc: ast.Loc{Start: start, End: l.pos}}
		}
	case '?':
		return Token{Kind: tokPlaceholder, Loc: ast.Loc{Start: start, End: l.pos}}
	case '=':
		// == is equivalent to =
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
		}
		return Token{Kind: int('='), Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Valid single-char operators/punctuation
	if ch == '(' || ch == ')' || ch == ',' || ch == ';' || ch == '+' || ch == '-' ||
		ch == '*' || ch == '/' || ch == '%' || ch == '~' || ch == '^' || ch == '&' ||
		ch == '|' || ch == '<' || ch == '>' || ch == '=' || ch == '.' || ch == '@' ||
		ch == ':' || ch == '!' || ch == '[' || ch == ']' || ch == '{' || ch == '}' {
		return Token{Kind: int(ch), Loc: ast.Loc{Start: start, End: l.pos}}
	}

	// Unknown character
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnknownChar,
	})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
