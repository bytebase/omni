package parser

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/starrocks/ast"
)

// Lexer is a StarRocks SQL tokenizer. Construct via NewLexer (or
// NewLexerWithOffset when tokenizing a substring) and call NextToken
// until it returns Token{Kind: tokEOF}. Lex errors are collected in
// Errors; each error is accompanied by a tokInvalid token so consumers
// can halt or proceed with best-effort parsing.
//
// A character no rule of StarRocksLex.g4 matches is dropped rather than
// reported: the front end's ANTLR lexer reports one and skips it, so '#', a
// backslash and a NUL separate tokens without ending a line and without a lex
// error. The pinned 3.4.10 server reads `SELECT 1 # x` as `SELECT 1 x` (one
// column named x), `SELECT 1+#2` as `SELECT 1+2`, `SELECT 1 # 2` as a syntax
// error at the 2, and `# x` as a syntax error at the x. Only characters
// outside the Basic Multilingual Plane (an emoji) and the ASCII bytes no
// grammar literal names are unmatched; every other non-ASCII character is a
// LETTER, and so is a byte that is not valid UTF-8, which decodes as U+FFFD.
//
// Dropping is faithful, not lenient: the engine drops the same characters, so
// omni accepts and rejects the same statements it does. The `suite:characters`
// rows of starrocks/parser/testdata/compat.tsv are the probes.
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
	for {
		l.skipWhitespaceAndComments()
		if l.pos >= len(l.input) {
			return Token{Kind: tokEOF, Loc: ast.Loc{Start: l.pos, End: l.pos}}
		}
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if IsUnmatched(r) {
			l.pos += size
			continue
		}
		break
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
		// DOT_IDENTIFIER is '.' DIGIT_IDENTIFIER (StarRocksLex.g4:574) and
		// wins the longest match over DECIMAL_VALUE, so '.1a' is one token
		// and '.1' is a number.
		if _, ok := l.dotIdentEnd(); ok {
			return l.scanDotIdent()
		}
		return l.scanNumber()
	case (ch == 'X' || ch == 'x') && (l.peek(1) == '\'' || l.peek(1) == '"'):
		// BINARY_SINGLE_QUOTED_TEXT and BINARY_DOUBLE_QUOTED_TEXT (StarRocksLex.g4).
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

// isIdentStart reports whether byte ch can start an identifier, which is what
// the dispatch in nextTokenInner asks of the byte at the current position.
// Every byte of a multi-byte character is above 127, so a character that
// reaches the dispatch starts an identifier whatever its later bytes are.
func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' || ch > 127
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

// isIdentRune reports whether r can appear in an identifier: LETTER in
// StarRocksLex.g4 is [a-zA-Z_$\u0080-\uffff], and its DIGIT is [0-9]. A
// character outside the Basic Multilingual Plane is not a LETTER, so it ends
// an identifier where the engine does not have one either.
func isIdentRune(r rune) bool {
	return r < 0x80 && isIdentCont(byte(r)) || r >= 0x80 && r <= 0xFFFF
}

// IsUnmatched reports whether r matches no rule of StarRocksLex.g4 and no
// literal of StarRocks.g4, which the lexer drops instead of rejecting.
//
// LETTER covers \u0080-\uffff: every character in the Basic Multilingual
// Plane is a letter, including U+FFFD, which is what a byte that is not valid
// UTF-8 decodes to in the front end and here. A character above the plane, and
// an ASCII byte no literal names -- '#', a backslash, a control character --
// has no token. Scans that decide whether text holds only blanks ask this
// too, so a dropped character there reads as blank.
func IsUnmatched(r rune) bool {
	if r >= 0x80 {
		return r > 0xFFFF
	}
	switch r {
	case ' ', '\t', '\n', '\r':
		// WS.
		return false
	}
	return !asciiStartsToken[byte(r)]
}

// asciiStartsToken marks the ASCII bytes a rule begins with: the letters, the
// digits and '_'/'$' of LETTER and DIGIT, the quotes and the backtick of the
// literals and identifiers, '.' (DECIMAL_VALUE, DOT, DOTDOTDOT, DOT_IDENTIFIER)
// and every operator and punctuation literal StarRocksLex.g4 and StarRocks.g4
// name. '#' and a backslash are absent, as neither grammar has a rule for
// them.
var asciiStartsToken = func() [128]bool {
	var t [128]bool
	for c := 'a'; c <= 'z'; c++ {
		t[c] = true
	}
	for c := 'A'; c <= 'Z'; c++ {
		t[c] = true
	}
	for c := '0'; c <= '9'; c++ {
		t[c] = true
	}
	for _, c := range "_$'\"`.()[]{},;:+-*/%^~&|<>=!@?" {
		t[c] = true
	}
	return t
}()

// --- Whitespace & Comments -------------------------------------------------

// skipWhitespaceAndComments skips what StarRocksLex.g4 puts on a hidden
// channel: WS, SIMPLE_COMMENT (a `--` to the end of the line) and
// BRACKETED_COMMENT. The grammar has no `//` and no `#` token, so `//` is two
// SLASH_SYMBOLs and `SELECT 1 // x` is a syntax error at the second '/'. `#`
// is dropped one character at a time by NextToken, not ended at.
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
// Unlike Snowflake, StarRocks allows newlines inside string literals.
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
// escape is supported (“ → `). No backslash escapes.
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
	quote := l.input[l.pos+1]
	l.pos += 2 // consume X and opening quote
	contentStart := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != quote {
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
	// DIGIT_IDENTIFIER (StarRocksLex.g4:564) is DIGIT (LETTER | DIGIT | '_')+
	// and outruns INTEGER_VALUE, so '1a' is one identifier. LETTER covers
	// U+0080-U+FFFF, so '1　AS' is one too.
	if end, ok := l.digitIdentEnd(); ok {
		l.pos = end
		return Token{Kind: tokIdent, Str: l.input[start:end], Loc: ast.Loc{Start: start, End: end}}
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

// digitIdentEnd reports where the DIGIT_IDENTIFIER that starts at the digits
// just scanned ends, if one starts there: DIGIT (LETTER | DIGIT | '_')+ needs
// a LETTER, '_' or another digit after the digits already scanned.
func (l *Lexer) digitIdentEnd() (int, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	if !isIdentRune(r) {
		return 0, false
	}
	return l.identRunEnd(l.pos + size), true
}

// dotIdentEnd reports where the DOT_IDENTIFIER (StarRocksLex.g4:574) that
// would start at the current '.' ends, if the text there is one: '.' and a
// DIGIT_IDENTIFIER. DECIMAL_VALUE also matches '.1', so the identifier wins
// only while it is the longer match; DOUBLE_VALUE's exponent ('.1E5') comes
// first in the lexer and keeps the text a number.
func (l *Lexer) dotIdentEnd() (int, bool) {
	i := l.pos + 1
	for i < len(l.input) && l.input[i] >= '0' && l.input[i] <= '9' {
		i++
	}
	if i >= len(l.input) {
		return 0, false
	}
	r, size := utf8.DecodeRuneInString(l.input[i:])
	if r == 'E' {
		j := i + 1
		if j < len(l.input) && (l.input[j] == '+' || l.input[j] == '-') {
			j++
		}
		if j < len(l.input) && l.input[j] >= '0' && l.input[j] <= '9' {
			return 0, false
		}
	}
	if !isIdentRune(r) {
		return 0, false
	}
	return l.identRunEnd(i + size), true
}

// identRunEnd scans the identifier characters from i and returns where they
// end.
func (l *Lexer) identRunEnd(i int) int {
	for i < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[i:])
		if !isIdentRune(r) {
			break
		}
		i += size
	}
	return i
}

// scanDotIdent scans a DOT_IDENTIFIER. Str holds the text after the '.', so
// the token names the identifier part it stands for.
func (l *Lexer) scanDotIdent() Token {
	end, _ := l.dotIdentEnd()
	start := l.pos
	l.pos = end
	return Token{Kind: tokDotIdent, Str: l.input[start+1 : end], Loc: ast.Loc{Start: start, End: end}}
}

// --- Identifier / Keyword Scanning -----------------------------------------

func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.start
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !isIdentRune(r) {
			break
		}
		l.pos += size
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

	// asciiStartsToken and the switch above hold the same set of ASCII bytes,
	// so nothing reaches here; TestTokenizeEveryByte pins that they agree.
	l.errors = append(l.errors, LexError{
		Loc: ast.Loc{Start: start, End: l.pos},
		Msg: errUnknownChar,
	})
	return Token{Kind: tokInvalid, Loc: ast.Loc{Start: start, End: l.pos}}
}
