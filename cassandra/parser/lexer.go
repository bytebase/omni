package parser

import (
	"strings"
	"unicode/utf8"
)

// Lexer tokenizes CQL input.
type Lexer struct {
	input   string
	lineIdx lineIndex
	pos     int
	Err     error
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, lineIdx: buildLineIndex(input)}
}

func (l *Lexer) makeError(msg string, start, end int) *ParseError {
	line, col := offsetToLineCol(l.lineIdx, start)
	nearEnd := start + 30
	if nearEnd > end {
		nearEnd = end
	}
	return &ParseError{
		Message: msg,
		Loc:     locFromOffsets(start, end),
		Line:    line,
		Column:  col,
		Near:    l.input[start:nearEnd],
	}
}

// Next returns the next token from the input.
func (l *Lexer) Next() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
	}

	start := l.pos
	ch := l.input[l.pos]

	switch {
	case ch == '\'':
		return l.scanString(start)
	case ch == '"':
		return l.scanQuotedIdentifier(start)
	case ch == '$' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '$':
		return l.scanCodeBlock(start)
	case ch == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X'):
		return l.scanHex(start)
	case isDigit(ch):
		return l.scanNumber(start)
	case isIdentStart(ch):
		return l.scanIdentOrKeyword(start)
	default:
		return l.scanOperator(start)
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\f' {
			l.pos++
			continue
		}
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}
		if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			l.pos += 2
			for l.pos+1 < len(l.input) {
				if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
					l.pos += 2
					break
				}
				l.pos++
			}
			if l.pos >= len(l.input) {
				// Unterminated block comment — just stop skipping
			}
			continue
		}
		break
	}
}

func (l *Lexer) scanString(start int) Token {
	l.pos++ // skip opening '
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				b.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokSTRING, Str: b.String(), Loc: start, End: l.pos}
		}
		b.WriteByte(ch)
		l.pos++
	}
	l.Err = l.makeError("unterminated string literal", start, l.pos)
	return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
}

func (l *Lexer) scanQuotedIdentifier(start int) Token {
	l.pos++ // skip opening "
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				b.WriteByte('"')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokQUOTED, Str: b.String(), Loc: start, End: l.pos}
		}
		b.WriteByte(ch)
		l.pos++
	}
	l.Err = l.makeError("unterminated quoted identifier", start, l.pos)
	return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
}

func (l *Lexer) scanCodeBlock(start int) Token {
	l.pos += 2 // skip opening $$
	idx := strings.Index(l.input[l.pos:], "$$")
	if idx < 0 {
		l.Err = l.makeError("unterminated code block", start, len(l.input))
		l.pos = len(l.input)
		return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
	}
	val := l.input[l.pos : l.pos+idx]
	l.pos += idx + 2
	return Token{Type: tokCODEBLOCK, Str: val, Loc: start, End: l.pos}
}

func (l *Lexer) scanHex(start int) Token {
	l.pos += 2 // skip 0x
	for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: tokHEX, Str: l.input[start:l.pos], Loc: start, End: l.pos}
}

func (l *Lexer) scanNumber(start int) Token {
	isFloat := false
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	// UUID check: digits followed by hex letters (a-f) may form the first 8-char group.
	numLen := l.pos - start
	if numLen < 8 && l.pos < len(l.input) && isHexLetter(l.input[l.pos]) {
		savedPos := l.pos
		for l.pos < len(l.input) && l.pos-start < 8 && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		if l.pos-start == 8 {
			str := l.input[start:l.pos]
			if isUUIDCandidate(str, l) {
				return l.scanUUID(start, str)
			}
		}
		l.pos = savedPos
	}

	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		next := l.pos + 1
		if next < len(l.input) && isDigit(l.input[next]) {
			isFloat = true
			l.pos++ // skip .
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		}
	}
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		ePos := l.pos
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		if l.pos >= len(l.input) || !isDigit(l.input[l.pos]) {
			l.pos = ePos
		} else {
			isFloat = true
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.pos++
			}
		}
	}
	str := l.input[start:l.pos]

	// A number like 550e8400 may actually be the first group of a UUID.
	if isUUIDCandidate(str, l) {
		return l.scanUUID(start, str)
	}

	tok := tokINTEGER
	if isFloat {
		tok = tokFLOAT
	}
	return Token{Type: tok, Str: str, Loc: start, End: l.pos}
}

func (l *Lexer) scanIdentOrKeyword(start int) Token {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if isIdentPart(ch) {
			l.pos++
		} else {
			break
		}
	}
	str := l.input[start:l.pos]
	lower := strings.ToLower(str)

	// Check if this looks like a UUID: 8-4-4-4-12 hex pattern.
	if isUUIDCandidate(str, l) {
		return l.scanUUID(start, str)
	}

	if tok, ok := keywords[lower]; ok {
		return Token{Type: tok, Str: str, Loc: start, End: l.pos}
	}

	return Token{Type: tokIDENT, Str: str, Loc: start, End: l.pos}
}

// isUUIDCandidate checks if the current identifier followed by upcoming chars
// forms a UUID pattern (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx).
func isUUIDCandidate(str string, l *Lexer) bool {
	if len(str) != 8 {
		return false
	}
	for _, c := range str {
		if !isHexDigitRune(c) {
			return false
		}
	}
	// Check if followed by - and more hex groups
	remaining := l.input[l.pos:]
	if len(remaining) < 28 { // -xxxx-xxxx-xxxx-xxxxxxxxxxxx
		return false
	}
	// Pattern: -XXXX-XXXX-XXXX-XXXXXXXXXXXX
	if remaining[0] != '-' {
		return false
	}
	expected := []int{4, 4, 4, 12}
	off := 1
	for i, groupLen := range expected {
		for j := 0; j < groupLen; j++ {
			if off >= len(remaining) || !isHexDigitByte(remaining[off]) {
				return false
			}
			off++
		}
		if i < len(expected)-1 {
			if off >= len(remaining) || remaining[off] != '-' {
				return false
			}
			off++
		}
	}
	// Make sure the UUID isn't followed by more identifier chars
	if off < len(remaining) && isIdentPart(remaining[off]) {
		return false
	}
	return true
}

func (l *Lexer) scanUUID(start int, firstGroup string) Token {
	// We already consumed the first 8-char group. Consume the rest.
	// Pattern: -xxxx-xxxx-xxxx-xxxxxxxxxxxx
	groups := []int{4, 4, 4, 12}
	for _, groupLen := range groups {
		l.pos++ // skip -
		l.pos += groupLen
	}
	return Token{Type: tokUUID, Str: l.input[start:l.pos], Loc: start, End: l.pos}
}

func (l *Lexer) scanOperator(start int) Token {
	ch := l.input[l.pos]
	l.pos++
	switch ch {
	case '.':
		return Token{Type: tokDOT, Str: ".", Loc: start, End: l.pos}
	case ',':
		return Token{Type: tokCOMMA, Str: ",", Loc: start, End: l.pos}
	case ';':
		return Token{Type: tokSEMI, Str: ";", Loc: start, End: l.pos}
	case ':':
		return Token{Type: tokCOLON, Str: ":", Loc: start, End: l.pos}
	case '(':
		return Token{Type: tokLPAREN, Str: "(", Loc: start, End: l.pos}
	case ')':
		return Token{Type: tokRPAREN, Str: ")", Loc: start, End: l.pos}
	case '{':
		return Token{Type: tokLBRACE, Str: "{", Loc: start, End: l.pos}
	case '}':
		return Token{Type: tokRBRACE, Str: "}", Loc: start, End: l.pos}
	case '[':
		return Token{Type: tokLBRACK, Str: "[", Loc: start, End: l.pos}
	case ']':
		return Token{Type: tokRBRACK, Str: "]", Loc: start, End: l.pos}
	case '*':
		return Token{Type: tokSTAR, Str: "*", Loc: start, End: l.pos}
	case '+':
		return Token{Type: tokPLUS, Str: "+", Loc: start, End: l.pos}
	case '-':
		return Token{Type: tokMINUS, Str: "-", Loc: start, End: l.pos}
	case '!':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
			return Token{Type: tokNE, Str: "!=", Loc: start, End: l.pos}
		}
		return Token{Type: tokILLEGAL, Str: "!", Loc: start, End: l.pos}
	case '?':
		return Token{Type: tokQMARK, Str: "?", Loc: start, End: l.pos}
	case '=':
		return Token{Type: tokEQ, Str: "=", Loc: start, End: l.pos}
	case '<':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
			return Token{Type: tokLTE, Str: "<=", Loc: start, End: l.pos}
		}
		return Token{Type: tokLT, Str: "<", Loc: start, End: l.pos}
	case '>':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.pos++
			return Token{Type: tokGTE, Str: ">=", Loc: start, End: l.pos}
		}
		return Token{Type: tokGT, Str: ">", Loc: start, End: l.pos}
	default:
		_, size := utf8.DecodeRuneInString(l.input[start:])
		l.pos = start + size
		return Token{Type: tokILLEGAL, Str: l.input[start:l.pos], Loc: start, End: l.pos}
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isHexLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isHexDigitRune(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isHexDigitByte(ch byte) bool {
	return isHexDigit(ch)
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch)
}
