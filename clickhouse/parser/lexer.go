package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Lexer tokenizes ClickHouse SQL source text one token at a time.
type Lexer struct {
	src string
	pos int
	err error
}

// NewLexer creates a streaming ClickHouse SQL lexer over src.
func NewLexer(src string) *Lexer {
	return &Lexer{src: src}
}

// Err returns the first lexical error encountered by the lexer, if any.
func (l *Lexer) Err() error {
	return l.err
}

// NextToken scans and returns the next ClickHouse SQL token. On a lexical
// error, NextToken records the first error (retrievable via Err), advances
// past the malformed span so callers always terminate, and returns a
// TokIllegal token.
func (l *Lexer) NextToken() Token {
	before := l.pos
	tok, err := l.nextToken()
	if err == nil {
		return tok
	}
	if l.err == nil {
		l.err = err
	}
	pos := before
	if pe, ok := err.(*ParseError); ok {
		pos = pe.Position
	}
	if l.pos <= before && l.pos < len(l.src) {
		l.pos = before + 1
	}
	if l.pos <= pos && l.pos < len(l.src) {
		l.pos = pos + 1
	}
	if pos > l.pos {
		pos = l.pos
	}
	return Token{
		Type: TokIllegal,
		Raw:  l.src[pos:l.pos],
		Pos:  pos,
		End:  l.pos,
	}
}

// Tokenize scans all tokens in src up to and including TokEOF, returning the
// first lexical error if any.
func Tokenize(src string) ([]Token, error) {
	return tokenize(src)
}

func tokenize(src string) ([]Token, error) {
	l := &Lexer{src: src}
	var out []Token
	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		out = append(out, tok)
		if tok.Type == TokEOF {
			return out, nil
		}
	}
}

func (l *Lexer) nextToken() (Token, error) {
	if err := l.skipWhitespaceAndComments(); err != nil {
		return Token{}, err
	}
	if l.pos >= len(l.src) {
		return Token{Type: TokEOF, Pos: l.pos, End: l.pos}, nil
	}
	start := l.pos
	c := l.src[l.pos]

	// Quoted strings
	if c == '\'' {
		return l.scanSingleQuoted()
	}
	// Dollar-quoted raw strings: $tag$...$tag$
	if c == '$' {
		if tok, ok, err := l.tryScanDollarQuoted(); ok || err != nil {
			return tok, err
		}
		l.pos++
		return Token{}, &ParseError{Message: "unexpected '$'", Position: start}
	}
	// Backtick or double-quoted identifiers
	if c == '`' || c == '"' {
		return l.scanQuotedIdent(c)
	}
	// Hex / binary string literals: X'...' or B'...'
	if (c == 'X' || c == 'x' || c == 'B' || c == 'b') && l.pos+1 < len(l.src) && l.src[l.pos+1] == '\'' {
		return l.scanPrefixedString(c)
	}
	// Numbers (ClickHouse numbers must start with a digit; leading '.' is always TokDot)
	if isDigit(c) {
		return l.scanNumber()
	}
	// Identifiers and keywords
	if isIdentStart(c) {
		l.pos++
		for l.pos < len(l.src) && isIdentPart(l.src[l.pos]) {
			l.pos++
		}
		raw := l.src[start:l.pos]
		return Token{
			Type: lookupKeyword(raw),
			Raw:  raw,
			Val:  raw,
			Pos:  start,
			End:  l.pos,
		}, nil
	}

	// Two-character operators
	if l.pos+1 < len(l.src) {
		two := l.src[l.pos : l.pos+2]
		var tt TokenType
		switch two {
		case "->":
			tt = TokArrow
		case "::":
			tt = TokDoubleColon
		case "||":
			tt = TokConcat
		case "==":
			tt = TokEqDouble
		case "!=", "<>":
			tt = TokNotEq
		case "<=":
			tt = TokLE
		case ">=":
			tt = TokGE
		}
		if tt != 0 {
			l.pos += 2
			return Token{Type: tt, Raw: two, Val: two, Pos: start, End: l.pos}, nil
		}
	}

	// Single-character punctuation/operators
	l.pos++
	var tt TokenType
	switch c {
	case '(':
		tt = TokLParen
	case ')':
		tt = TokRParen
	case '[':
		tt = TokLBracket
	case ']':
		tt = TokRBracket
	case '{':
		tt = TokLBrace
	case '}':
		tt = TokRBrace
	case ',':
		tt = TokComma
	case '.':
		tt = TokDot
	case ';':
		tt = TokSemicolon
	case ':':
		tt = TokColon
	case '?':
		tt = TokQuestion
	case '+':
		tt = TokPlus
	case '-':
		tt = TokMinus
	case '*':
		tt = TokStar
	case '/':
		tt = TokSlash
	case '%':
		tt = TokPercent
	case '=':
		tt = TokEq
	case '<':
		tt = TokLT
	case '>':
		tt = TokGT
	default:
		return Token{}, &ParseError{
			Message:  fmt.Sprintf("unexpected character %q", c),
			Position: start,
		}
	}
	raw := l.src[start:l.pos]
	return Token{Type: tt, Raw: raw, Val: raw, Pos: start, End: l.pos}, nil
}

func (l *Lexer) skipWhitespaceAndComments() error {
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		// Whitespace
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\v' || c == '\f' {
			l.pos++
			continue
		}
		// Single-line comment: --
		if c == '-' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.src) && l.src[l.pos] != '\n' && l.src[l.pos] != '\r' {
				l.pos++
			}
			continue
		}
		// Single-line comment: #! or # (ClickHouse supports #! and # line comments)
		if c == '#' {
			l.pos++
			for l.pos < len(l.src) && l.src[l.pos] != '\n' && l.src[l.pos] != '\r' {
				l.pos++
			}
			continue
		}
		// Multi-line comment: /* ... */ (ClickHouse supports nested block comments)
		if c == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '*' {
			start := l.pos
			l.pos += 2
			depth := 1
			for l.pos < len(l.src) && depth > 0 {
				if l.src[l.pos] == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '*' {
					depth++
					l.pos += 2
				} else if l.src[l.pos] == '*' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '/' {
					depth--
					l.pos += 2
				} else {
					l.pos++
				}
			}
			if depth > 0 {
				return &ParseError{Message: "unterminated block comment", Position: start}
			}
			continue
		}
		break
	}
	return nil
}

func (l *Lexer) scanSingleQuoted() (Token, error) {
	start := l.pos
	l.pos++ // skip opening '
	var b strings.Builder
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == '\\' {
			if err := l.scanEscape(&b, start, "unterminated string literal"); err != nil {
				return Token{}, err
			}
			continue
		}
		if c == '\'' {
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '\'' {
				b.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{
				Type: TokString,
				Raw:  l.src[start:l.pos],
				Val:  b.String(),
				Pos:  start,
				End:  l.pos,
			}, nil
		}
		b.WriteByte(c)
		l.pos++
	}
	return Token{}, &ParseError{Message: "unterminated string literal", Position: start}
}

func (l *Lexer) tryScanDollarQuoted() (Token, bool, error) {
	start := l.pos
	i := l.pos + 1
	for i < len(l.src) && isIdentPart(l.src[i]) {
		i++
	}
	if i >= len(l.src) || l.src[i] != '$' {
		return Token{}, false, nil
	}
	tag := l.src[start : i+1]
	l.pos = i + 1
	endIdx := strings.Index(l.src[l.pos:], tag)
	if endIdx < 0 {
		l.pos = len(l.src)
		return Token{}, true, &ParseError{Message: fmt.Sprintf("unterminated dollar-quoted string %s", tag), Position: start}
	}
	val := l.src[l.pos : l.pos+endIdx]
	l.pos += endIdx + len(tag)
	return Token{
		Type: TokString,
		Raw:  l.src[start:l.pos],
		Val:  val,
		Pos:  start,
		End:  l.pos,
	}, true, nil
}

func (l *Lexer) scanQuotedIdent(quote byte) (Token, error) {
	start := l.pos
	l.pos++ // skip opening quote
	var b strings.Builder
	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if c == '\\' {
			if err := l.scanEscape(&b, start, "unterminated quoted identifier"); err != nil {
				return Token{}, err
			}
			continue
		}
		if c == quote {
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == quote {
				b.WriteByte(quote)
				l.pos += 2
				continue
			}
			l.pos++
			return Token{
				Type: TokQuotedIdent,
				Raw:  l.src[start:l.pos],
				Val:  b.String(),
				Pos:  start,
				End:  l.pos,
			}, nil
		}
		b.WriteByte(c)
		l.pos++
	}
	return Token{}, &ParseError{Message: "unterminated quoted identifier", Position: start}
}

// scanEscape decodes a ClickHouse backslash escape sequence starting at l.pos ('\\').
func (l *Lexer) scanEscape(b *strings.Builder, tokenStart int, unterminatedMsg string) error {
	if l.pos+1 >= len(l.src) {
		l.pos = len(l.src)
		return &ParseError{Message: unterminatedMsg, Position: tokenStart}
	}
	l.pos++ // skip '\\'
	c := l.src[l.pos]
	l.pos++
	switch c {
	case 'a':
		b.WriteByte('\a')
	case 'b':
		b.WriteByte('\b')
	case 'e':
		b.WriteByte(0x1b)
	case 'f':
		b.WriteByte('\f')
	case 'n':
		b.WriteByte('\n')
	case 'r':
		b.WriteByte('\r')
	case 't':
		b.WriteByte('\t')
	case 'v':
		b.WriteByte('\v')
	case '0':
		b.WriteByte(0)
	case '\\', '\'', '"', '`', '/', '=':
		b.WriteByte(c)
	case 'x', 'X':
		if l.pos+2 <= len(l.src) && isHexDigit(l.src[l.pos]) && isHexDigit(l.src[l.pos+1]) {
			v, _ := strconv.ParseUint(l.src[l.pos:l.pos+2], 16, 8)
			b.WriteByte(byte(v))
			l.pos += 2
		} else {
			b.WriteByte('\\')
			b.WriteByte(c)
		}
	case 'u':
		if l.pos+4 <= len(l.src) &&
			isHexDigit(l.src[l.pos]) && isHexDigit(l.src[l.pos+1]) &&
			isHexDigit(l.src[l.pos+2]) && isHexDigit(l.src[l.pos+3]) {
			v, _ := strconv.ParseUint(l.src[l.pos:l.pos+4], 16, 32)
			b.WriteRune(rune(v))
			l.pos += 4
		} else {
			b.WriteByte('\\')
			b.WriteByte(c)
		}
	default:
		b.WriteByte('\\')
		b.WriteByte(c)
	}
	return nil
}

func (l *Lexer) scanPrefixedString(prefix byte) (Token, error) {
	start := l.pos
	isHex := prefix == 'X' || prefix == 'x'
	l.pos += 2 // skip prefix and '
	valStart := l.pos
	badPos := -1
	var badChar byte
	for l.pos < len(l.src) && l.src[l.pos] != '\'' {
		c := l.src[l.pos]
		if badPos < 0 && ((isHex && !isHexDigit(c)) || (!isHex && c != '0' && c != '1')) {
			badPos = l.pos
			badChar = c
		}
		l.pos++
	}
	if l.pos >= len(l.src) {
		if badPos >= 0 {
			if isHex {
				return Token{}, &ParseError{Message: fmt.Sprintf("invalid hex digit %q in hex string literal", badChar), Position: badPos}
			}
			return Token{}, &ParseError{Message: fmt.Sprintf("invalid binary digit %q in binary string literal", badChar), Position: badPos}
		}
		return Token{}, &ParseError{Message: "unterminated string literal", Position: start}
	}
	val := l.src[valStart:l.pos]
	l.pos++ // skip closing '
	if badPos >= 0 {
		if isHex {
			return Token{}, &ParseError{Message: fmt.Sprintf("invalid hex digit %q in hex string literal", badChar), Position: badPos}
		}
		return Token{}, &ParseError{Message: fmt.Sprintf("invalid binary digit %q in binary string literal", badChar), Position: badPos}
	}
	tt := TokBinString
	if isHex {
		tt = TokHexString
	}
	return Token{
		Type: tt,
		Raw:  l.src[start:l.pos],
		Val:  val,
		Pos:  start,
		End:  l.pos,
	}, nil
}

func (l *Lexer) scanNumber() (Token, error) {
	start := l.pos
	isFloat := false

	// Hex (0x...) or Binary (0b...)
	if l.src[l.pos] == '0' && l.pos+1 < len(l.src) {
		next := l.src[l.pos+1]
		if next == 'x' || next == 'X' {
			l.pos += 2
			digits := 0
			for l.pos < len(l.src) && (isHexDigit(l.src[l.pos]) || l.src[l.pos] == '_') {
				if l.src[l.pos] != '_' {
					digits++
				}
				l.pos++
			}
			if digits == 0 {
				return Token{}, &ParseError{Message: "invalid hexadecimal literal", Position: start}
			}
			if l.pos < len(l.src) && l.src[l.pos] == '.' {
				isFloat = true
				l.pos++
				for l.pos < len(l.src) && (isHexDigit(l.src[l.pos]) || l.src[l.pos] == '_') {
					l.pos++
				}
			}
			if l.pos < len(l.src) && (l.src[l.pos] == 'p' || l.src[l.pos] == 'P') {
				isFloat = true
				l.pos++
				if l.pos < len(l.src) && (l.src[l.pos] == '+' || l.src[l.pos] == '-') {
					l.pos++
				}
				expDigits := 0
				for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
					expDigits++
					l.pos++
				}
				if expDigits == 0 {
					return Token{}, &ParseError{Message: "invalid hexadecimal float exponent", Position: start}
				}
			}
			tt := TokInt
			if isFloat {
				tt = TokFloat
			}
			raw := l.src[start:l.pos]
			return Token{Type: tt, Raw: raw, Val: raw, Pos: start, End: l.pos}, nil
		}
		if next == 'b' || next == 'B' {
			l.pos += 2
			digits := 0
			for l.pos < len(l.src) && (l.src[l.pos] == '0' || l.src[l.pos] == '1' || l.src[l.pos] == '_') {
				if l.src[l.pos] != '_' {
					digits++
				}
				l.pos++
			}
			if digits == 0 {
				return Token{}, &ParseError{Message: "invalid binary literal", Position: start}
			}
			raw := l.src[start:l.pos]
			return Token{Type: TokInt, Raw: raw, Val: raw, Pos: start, End: l.pos}, nil
		}
	}

	// Leading dot float (.5)
	if l.src[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.src) && (isDigit(l.src[l.pos]) || l.src[l.pos] == '_') {
			l.pos++
		}
	} else {
		for l.pos < len(l.src) && (isDigit(l.src[l.pos]) || l.src[l.pos] == '_') {
			l.pos++
		}
		if l.pos < len(l.src) && l.src[l.pos] == '.' {
			isFloat = true
			l.pos++
			for l.pos < len(l.src) && (isDigit(l.src[l.pos]) || l.src[l.pos] == '_') {
				l.pos++
			}
		}
	}

	// Optional exponent: e/E [+-]? digits
	if l.pos < len(l.src) && (l.src[l.pos] == 'e' || l.src[l.pos] == 'E') {
		save := l.pos
		l.pos++
		if l.pos < len(l.src) && (l.src[l.pos] == '+' || l.src[l.pos] == '-') {
			l.pos++
		}
		expDigits := 0
		for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
			expDigits++
			l.pos++
		}
		if expDigits > 0 {
			isFloat = true
		} else {
			l.pos = save
		}
	}

	// ClickHouse rejects identifier characters immediately glued to a number (e.g. `1.zzz`).
	if l.pos < len(l.src) && isIdentStart(l.src[l.pos]) {
		return Token{}, &ParseError{Message: fmt.Sprintf("invalid numeric literal %q", l.src[start:l.pos+1]), Position: start}
	}

	tt := TokInt
	if isFloat {
		tt = TokFloat
	}
	raw := l.src[start:l.pos]
	return Token{Type: tt, Raw: raw, Val: raw, Pos: start, End: l.pos}, nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}
