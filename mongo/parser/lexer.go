package parser

import (
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes mongosh input character by character.
type Lexer struct {
	input string
	pos   int
	start int
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// String literals
	if ch == '"' || ch == '\'' {
		return l.scanString(ch)
	}

	// Numbers: digits, or leading dot followed by digit
	if ch >= '0' && ch <= '9' {
		return l.scanNumber()
	}
	if ch == '.' && l.pos+1 < len(l.input) && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9' {
		return l.scanNumber()
	}

	// Negative numbers: minus followed by digit or dot-digit
	if ch == '-' {
		if l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			if next >= '0' && next <= '9' {
				return l.scanNumber()
			}
			if next == '.' && l.pos+2 < len(l.input) && l.input[l.pos+2] >= '0' && l.input[l.pos+2] <= '9' {
				return l.scanNumber()
			}
		}
	}

	// Regex literal /pattern/flags
	// Disambiguate from division: regex can appear at start of input or after certain tokens.
	if ch == '/' {
		if l.canBeRegex() {
			return l.scanRegex()
		}
	}

	// Identifiers (including $-prefixed)
	if isIdentStart(ch) || ch == '$' {
		return l.scanIdentOrKeyword()
	}
	// Unicode identifiers
	if ch >= 0x80 {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if r != utf8.RuneError && (unicode.IsLetter(r) || r == '_') {
			return l.scanIdentOrKeyword()
		}
	}

	// Multi-character operators
	if l.pos+1 < len(l.input) {
		next := l.input[l.pos+1]
		switch {
		case ch == '=' && next == '=':
			if l.pos+2 < len(l.input) && l.input[l.pos+2] == '=' {
				l.pos += 3
				return Token{Type: tokEqEqEq, Str: "===", Loc: l.start, End: l.pos}
			}
			l.pos += 2
			return Token{Type: tokEqEq, Str: "==", Loc: l.start, End: l.pos}
		case ch == '!' && next == '=':
			if l.pos+2 < len(l.input) && l.input[l.pos+2] == '=' {
				l.pos += 3
				return Token{Type: tokNotEqEq, Str: "!==", Loc: l.start, End: l.pos}
			}
			l.pos += 2
			return Token{Type: tokNotEq, Str: "!=", Loc: l.start, End: l.pos}
		case ch == '<' && next == '=':
			l.pos += 2
			return Token{Type: tokLessEq, Str: "<=", Loc: l.start, End: l.pos}
		case ch == '>' && next == '=':
			l.pos += 2
			return Token{Type: tokGreaterEq, Str: ">=", Loc: l.start, End: l.pos}
		case ch == '&' && next == '&':
			l.pos += 2
			return Token{Type: tokAnd, Str: "&&", Loc: l.start, End: l.pos}
		case ch == '|' && next == '|':
			l.pos += 2
			return Token{Type: tokOr, Str: "||", Loc: l.start, End: l.pos}
		case ch == '=' && next == '>':
			l.pos += 2
			return Token{Type: tokArrow, Str: "=>", Loc: l.start, End: l.pos}
		case ch == '.' && next == '.' && l.pos+2 < len(l.input) && l.input[l.pos+2] == '.':
			l.pos += 3
			return Token{Type: tokSpread, Str: "...", Loc: l.start, End: l.pos}
		case ch == '?' && next == '.':
			l.pos += 2
			return Token{Type: tokOptChain, Str: "?.", Loc: l.start, End: l.pos}
		case ch == '?' && next == '?':
			l.pos += 2
			return Token{Type: tokNullCoal, Str: "??", Loc: l.start, End: l.pos}
		case ch == '+' && next == '+':
			l.pos += 2
			return Token{Type: tokPlusPlus, Str: "++", Loc: l.start, End: l.pos}
		case ch == '-' && next == '-':
			l.pos += 2
			return Token{Type: tokMinusMinus, Str: "--", Loc: l.start, End: l.pos}
		case ch == '+' && next == '=':
			l.pos += 2
			return Token{Type: tokPlusEq, Str: "+=", Loc: l.start, End: l.pos}
		case ch == '-' && next == '=':
			l.pos += 2
			return Token{Type: tokMinusEq, Str: "-=", Loc: l.start, End: l.pos}
		case ch == '*' && next == '*':
			l.pos += 2
			return Token{Type: tokPower, Str: "**", Loc: l.start, End: l.pos}
		case ch == '*' && next == '=':
			l.pos += 2
			return Token{Type: tokStarEq, Str: "*=", Loc: l.start, End: l.pos}
		case ch == '/' && next == '=':
			l.pos += 2
			return Token{Type: tokSlashEq, Str: "/=", Loc: l.start, End: l.pos}
		case ch == '%' && next == '=':
			l.pos += 2
			return Token{Type: tokPercentEq, Str: "%=", Loc: l.start, End: l.pos}
		case ch == '<' && next == '<':
			l.pos += 2
			return Token{Type: tokShiftLeft, Str: "<<", Loc: l.start, End: l.pos}
		case ch == '>' && next == '>':
			l.pos += 2
			return Token{Type: tokShiftRight, Str: ">>", Loc: l.start, End: l.pos}
		}
	}

	// Single character tokens — returned as their ASCII value
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start, End: l.pos}
}

// skipWhitespaceAndComments skips whitespace, // line comments, and /* */ block comments.
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		// Comments
		if ch == '/' && l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			// Line comment //
			if next == '/' {
				l.pos += 2
				for l.pos < len(l.input) && l.input[l.pos] != '\n' {
					l.pos++
				}
				continue
			}
			// Block comment /* */
			if next == '*' {
				l.pos += 2
				for l.pos+1 < len(l.input) {
					if l.input[l.pos] == '*' && l.input[l.pos+1] == '/' {
						l.pos += 2
						break
					}
					l.pos++
				}
				// If we hit end without closing, just stop
				if l.pos >= len(l.input) {
					return
				}
				continue
			}
		}

		break
	}
}

// scanString scans a single or double quoted string with escape sequences.
func (l *Lexer) scanString(quote byte) Token {
	l.pos++ // skip opening quote
	var buf []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == quote {
			l.pos++ // skip closing quote
			return Token{Type: tokString, Str: string(buf), Loc: l.start, End: l.pos}
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.pos++ // skip backslash
			esc := l.input[l.pos]
			l.pos++
			switch esc {
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case '\\':
				buf = append(buf, '\\')
			case '\'':
				buf = append(buf, '\'')
			case '"':
				buf = append(buf, '"')
			case '0':
				buf = append(buf, 0)
			case 'b':
				buf = append(buf, '\b')
			case 'f':
				buf = append(buf, '\f')
			case 'v':
				buf = append(buf, '\v')
			case 'u':
				// \uXXXX
				r := l.parseHex4()
				if r >= 0 {
					var ubuf [4]byte
					n := utf8.EncodeRune(ubuf[:], rune(r))
					buf = append(buf, ubuf[:n]...)
				} else {
					buf = append(buf, '\\', 'u')
				}
			case 'x':
				// \xHH
				if l.pos+1 < len(l.input) {
					h1 := hexVal(l.input[l.pos])
					h2 := hexVal(l.input[l.pos+1])
					if h1 >= 0 && h2 >= 0 {
						buf = append(buf, byte(h1<<4|h2))
						l.pos += 2
					} else {
						buf = append(buf, '\\', 'x')
					}
				} else {
					buf = append(buf, '\\', 'x')
				}
			default:
				// Unknown escape — include as-is
				buf = append(buf, '\\', esc)
			}
			continue
		}
		buf = append(buf, ch)
		l.pos++
	}
	// Unterminated string — return what we have
	return Token{Type: tokString, Str: string(buf), Loc: l.start, End: l.pos}
}

// parseHex4 parses 4 hex digits for \uXXXX escape. Returns -1 on failure.
func (l *Lexer) parseHex4() int {
	if l.pos+3 >= len(l.input) {
		return -1
	}
	val := 0
	for i := 0; i < 4; i++ {
		h := hexVal(l.input[l.pos+i])
		if h < 0 {
			return -1
		}
		val = val<<4 | h
	}
	l.pos += 4
	return val
}

// hexVal returns the hex digit value or -1.
func hexVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}

// scanNumber scans an integer, float, hex, or scientific notation number.
// Handles optional leading minus sign.
func (l *Lexer) scanNumber() Token {
	start := l.pos

	// Optional leading minus
	if l.pos < len(l.input) && l.input[l.pos] == '-' {
		l.pos++
	}

	// Hex: 0x or 0X
	if l.pos+1 < len(l.input) && l.input[l.pos] == '0' && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: tokNumber, Str: l.input[start:l.pos], Loc: l.start, End: l.pos}
	}

	// Integer part (may be empty for .5)
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}

	// Fractional part
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	// Exponent part
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}

	return Token{Type: tokNumber, Str: l.input[start:l.pos], Loc: l.start, End: l.pos}
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// scanRegex scans a regex literal /pattern/flags.
func (l *Lexer) scanRegex() Token {
	l.pos++ // skip opening /
	var pattern []byte
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '/' {
			l.pos++ // skip closing /
			// Scan flags
			flagStart := l.pos
			for l.pos < len(l.input) {
				c := l.input[l.pos]
				if c == 'g' || c == 'i' || c == 'm' || c == 's' || c == 'u' || c == 'y' {
					l.pos++
				} else {
					break
				}
			}
			flags := l.input[flagStart:l.pos]
			return Token{Type: tokRegex, Str: string(pattern) + "/" + flags, Loc: l.start, End: l.pos}
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			// Escaped character in regex
			pattern = append(pattern, ch, l.input[l.pos+1])
			l.pos += 2
			continue
		}
		if ch == '\n' {
			break // unterminated regex
		}
		pattern = append(pattern, ch)
		l.pos++
	}
	// Unterminated — return as division operator
	l.pos = l.start + 1
	return Token{Type: int('/'), Str: "/", Loc: l.start, End: l.pos}
}

// canBeRegex returns true if a '/' at the current position could start a regex literal
// rather than being a division operator.
func (l *Lexer) canBeRegex() bool {
	// Simple heuristic: regex can appear at start of input or after punctuation/operators
	// that cannot end an expression. After identifiers, numbers, or closing brackets, it's division.
	if l.start == 0 {
		return true
	}
	// Look at what came before (skip whitespace backwards)
	i := l.pos - 1
	for i >= 0 && (l.input[i] == ' ' || l.input[i] == '\t' || l.input[i] == '\n' || l.input[i] == '\r') {
		i--
	}
	if i < 0 {
		return true
	}
	prev := l.input[i]
	// After these, '/' is division
	switch prev {
	case ')', ']', '}':
		return false
	case '_':
		return false
	}
	if (prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') {
		return false
	}
	return true
}

// scanIdentOrKeyword scans an identifier or keyword.
func (l *Lexer) scanIdentOrKeyword() Token {
	start := l.pos

	// First character: letter, _, or $, or Unicode letter
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError && size <= 1 {
		l.pos++
	} else {
		l.pos += size
	}

	// Subsequent characters: letter, digit, _, $, or Unicode
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '$' {
			l.pos++
			continue
		}
		if ch >= 0x80 {
			r, size := utf8.DecodeRuneInString(l.input[l.pos:])
			if r != utf8.RuneError && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
				l.pos += size
				continue
			}
		}
		break
	}

	word := l.input[start:l.pos]

	// Check keyword map (case-sensitive for mongosh)
	if kw, ok := keywords[word]; ok {
		return Token{Type: kw, Str: word, Loc: l.start, End: l.pos}
	}

	return Token{Type: tokIdent, Str: word, Loc: l.start, End: l.pos}
}

// isIdentStart returns true if c can start an identifier.
func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
