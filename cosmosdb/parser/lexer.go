// Package parser implements a hand-written parser for Azure Cosmos DB NoSQL SQL API.
package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Token type constants.
const (
	tokEOF = 0
)

// Literal token types.
const (
	tokICONST = iota + 1000 // integer constant
	tokFCONST               // float constant
	tokHCONST               // hex constant
	tokSCONST               // single-quoted string
	tokDCONST               // double-quoted string
	tokIDENT                 // identifier
	tokPARAM                 // @param
)

// Operator token types.
const (
	tokDOT      = iota + 2000
	tokCOMMA    // ,
	tokCOLON    // :
	tokLPAREN   // (
	tokRPAREN   // )
	tokLBRACK   // [
	tokRBRACK   // ]
	tokLBRACE   // {
	tokRBRACE   // }
	tokSTAR     // *
	tokPLUS     // +
	tokMINUS    // -
	tokDIV      // /
	tokMOD      // %
	tokEQ       // =
	tokNE       // !=
	tokNE2      // <>
	tokLT       // <
	tokLE       // <=
	tokGT       // >
	tokGE       // >=
	tokBITAND   // &
	tokBITOR    // |
	tokBITXOR   // ^
	tokBITNOT   // ~
	tokLSHIFT   // <<
	tokRSHIFT   // >>
	tokURSHIFT  // >>>
	tokCONCAT   // ||
	tokCOALESCE // ??
	tokQUESTION // ?
)

// Keyword token types.
const (
	tokSELECT   = iota + 3000
	tokFROM      // FROM
	tokWHERE     // WHERE
	tokAND       // AND
	tokOR        // OR
	tokNOT       // NOT
	tokIN        // IN
	tokBETWEEN   // BETWEEN
	tokLIKE      // LIKE
	tokESCAPE    // ESCAPE
	tokAS        // AS
	tokJOIN      // JOIN
	tokTOP       // TOP
	tokDISTINCT  // DISTINCT
	tokVALUE     // VALUE
	tokORDER     // ORDER
	tokBY        // BY
	tokGROUP     // GROUP
	tokHAVING    // HAVING
	tokOFFSET    // OFFSET
	tokLIMIT     // LIMIT
	tokASC       // ASC
	tokDESC      // DESC
	tokEXISTS    // EXISTS
	tokTRUE      // TRUE
	tokFALSE     // FALSE
	tokNULL      // NULL
	tokUNDEFINED // UNDEFINED
	tokUDF       // UDF
	tokARRAY     // ARRAY
	tokROOT      // ROOT
	tokRANK      // RANK
	tokINFINITY  // Infinity (case-sensitive)
	tokNAN       // NaN (case-sensitive)
)

// keywords maps lowercase keyword strings to their token types.
var keywords = map[string]int{
	"select":    tokSELECT,
	"from":      tokFROM,
	"where":     tokWHERE,
	"and":       tokAND,
	"or":        tokOR,
	"not":       tokNOT,
	"in":        tokIN,
	"between":   tokBETWEEN,
	"like":      tokLIKE,
	"escape":    tokESCAPE,
	"as":        tokAS,
	"join":      tokJOIN,
	"top":       tokTOP,
	"distinct":  tokDISTINCT,
	"value":     tokVALUE,
	"order":     tokORDER,
	"by":        tokBY,
	"group":     tokGROUP,
	"having":    tokHAVING,
	"offset":    tokOFFSET,
	"limit":     tokLIMIT,
	"asc":       tokASC,
	"desc":      tokDESC,
	"exists":    tokEXISTS,
	"true":      tokTRUE,
	"false":     tokFALSE,
	"null":      tokNULL,
	"undefined": tokUNDEFINED,
	"udf":       tokUDF,
	"array":     tokARRAY,
	"root":      tokROOT,
	"rank":      tokRANK,
}

// Token represents a single lexer token.
type Token struct {
	Type int
	Str  string
	Loc  int // byte offset in original input
}

// Lexer is the hand-written tokenizer for Cosmos DB SQL.
type Lexer struct {
	input string
	pos   int   // current read position (next byte to consume)
	start int   // start position of token being scanned
	Err   error // first error encountered
}

// NewLexer creates a new Lexer for the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token from the input.
func (l *Lexer) Next() Token {
	l.skipWhitespaceAndComments()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// String literals.
	if ch == '\'' {
		return l.scanString('\'', tokSCONST)
	}
	if ch == '"' {
		return l.scanString('"', tokDCONST)
	}

	// Parameter reference.
	if ch == '@' {
		return l.scanParam()
	}

	// Numbers: digit or hex 0x.
	if ch >= '0' && ch <= '9' {
		return l.scanNumber()
	}

	// Identifiers and keywords.
	if isIdentStart(ch) {
		return l.scanIdentOrKeyword()
	}

	// Operators.
	return l.scanOperator()
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}
		// Line comment: -- to EOL
		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			l.pos += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}
		break
	}
}

func (l *Lexer) scanString(quote byte, tokType int) Token {
	l.pos++ // skip opening quote
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\\' {
			l.pos++
			if l.pos >= len(l.input) {
				l.Err = fmt.Errorf("unexpected end of input in string escape at position %d", l.start)
				return Token{Type: tokEOF, Loc: l.start}
			}
			esc := l.input[l.pos]
			l.pos++
			switch esc {
			case 'b':
				buf.WriteByte('\b')
			case 't':
				buf.WriteByte('\t')
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 'f':
				buf.WriteByte('\f')
			case '"':
				buf.WriteByte('"')
			case '\'':
				buf.WriteByte('\'')
			case '\\':
				buf.WriteByte('\\')
			case '/':
				buf.WriteByte('/')
			case 'u':
				if l.pos+4 > len(l.input) {
					l.Err = fmt.Errorf("incomplete unicode escape at position %d", l.pos-2)
					return Token{Type: tokEOF, Loc: l.start}
				}
				hex := l.input[l.pos : l.pos+4]
				var val rune
				for _, h := range []byte(hex) {
					val <<= 4
					switch {
					case h >= '0' && h <= '9':
						val |= rune(h - '0')
					case h >= 'a' && h <= 'f':
						val |= rune(h-'a') + 10
					case h >= 'A' && h <= 'F':
						val |= rune(h-'A') + 10
					default:
						l.Err = fmt.Errorf("invalid unicode escape \\u%s at position %d", hex, l.pos-2)
						return Token{Type: tokEOF, Loc: l.start}
					}
				}
				l.pos += 4
				buf.WriteRune(val)
			default:
				l.Err = fmt.Errorf("invalid escape \\%c at position %d", esc, l.pos-1)
				return Token{Type: tokEOF, Loc: l.start}
			}
			continue
		}
		if ch == quote {
			l.pos++
			return Token{Type: tokType, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated string starting at position %d", l.start)
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) scanParam() Token {
	l.pos++ // skip @
	start := l.pos
	for l.pos < len(l.input) && isIdentContinue(l.input[l.pos]) {
		l.pos++
	}
	if l.pos == start {
		l.Err = fmt.Errorf("expected parameter name after @ at position %d", l.start)
		return Token{Type: tokEOF, Loc: l.start}
	}
	return Token{Type: tokPARAM, Str: l.input[start:l.pos], Loc: l.start}
}

func (l *Lexer) scanNumber() Token {
	// Hex: 0x or 0X
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		start := l.pos
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		if l.pos == start {
			l.Err = fmt.Errorf("expected hex digits after 0x at position %d", l.start)
			return Token{Type: tokEOF, Loc: l.start}
		}
		return Token{Type: tokHCONST, Str: l.input[l.start:l.pos], Loc: l.start}
	}

	isFloat := false
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}
	// Decimal point.
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			l.pos++
		}
	}
	// Scientific notation.
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

	tokType := tokICONST
	if isFloat {
		tokType = tokFCONST
	}
	return Token{Type: tokType, Str: l.input[l.start:l.pos], Loc: l.start}
}

func (l *Lexer) scanIdentOrKeyword() Token {
	for l.pos < len(l.input) && isIdentContinue(l.input[l.pos]) {
		l.pos++
	}
	raw := l.input[l.start:l.pos]

	// Case-sensitive keywords: Infinity and NaN.
	if raw == "Infinity" {
		return Token{Type: tokINFINITY, Str: raw, Loc: l.start}
	}
	if raw == "NaN" {
		return Token{Type: tokNAN, Str: raw, Loc: l.start}
	}

	lower := strings.ToLower(raw)
	if tokType, ok := keywords[lower]; ok {
		return Token{Type: tokType, Str: raw, Loc: l.start}
	}
	return Token{Type: tokIDENT, Str: raw, Loc: l.start}
}

func (l *Lexer) scanOperator() Token {
	ch := l.input[l.pos]
	l.pos++

	// Two-character lookahead operators.
	if l.pos < len(l.input) {
		next := l.input[l.pos]
		switch {
		case ch == '!' && next == '=':
			l.pos++
			return Token{Type: tokNE, Str: "!=", Loc: l.start}
		case ch == '<' && next == '>':
			l.pos++
			return Token{Type: tokNE2, Str: "<>", Loc: l.start}
		case ch == '<' && next == '=':
			l.pos++
			return Token{Type: tokLE, Str: "<=", Loc: l.start}
		case ch == '<' && next == '<':
			l.pos++
			return Token{Type: tokLSHIFT, Str: "<<", Loc: l.start}
		case ch == '>' && next == '=':
			l.pos++
			return Token{Type: tokGE, Str: ">=", Loc: l.start}
		case ch == '>' && next == '>':
			l.pos++
			// Check for >>> (unsigned right shift).
			if l.pos < len(l.input) && l.input[l.pos] == '>' {
				l.pos++
				return Token{Type: tokURSHIFT, Str: ">>>", Loc: l.start}
			}
			return Token{Type: tokRSHIFT, Str: ">>", Loc: l.start}
		case ch == '|' && next == '|':
			l.pos++
			return Token{Type: tokCONCAT, Str: "||", Loc: l.start}
		case ch == '?' && next == '?':
			l.pos++
			return Token{Type: tokCOALESCE, Str: "??", Loc: l.start}
		}
	}

	// Single character operators.
	switch ch {
	case '.':
		return Token{Type: tokDOT, Str: ".", Loc: l.start}
	case ',':
		return Token{Type: tokCOMMA, Str: ",", Loc: l.start}
	case ':':
		return Token{Type: tokCOLON, Str: ":", Loc: l.start}
	case '(':
		return Token{Type: tokLPAREN, Str: "(", Loc: l.start}
	case ')':
		return Token{Type: tokRPAREN, Str: ")", Loc: l.start}
	case '[':
		return Token{Type: tokLBRACK, Str: "[", Loc: l.start}
	case ']':
		return Token{Type: tokRBRACK, Str: "]", Loc: l.start}
	case '{':
		return Token{Type: tokLBRACE, Str: "{", Loc: l.start}
	case '}':
		return Token{Type: tokRBRACE, Str: "}", Loc: l.start}
	case '*':
		return Token{Type: tokSTAR, Str: "*", Loc: l.start}
	case '+':
		return Token{Type: tokPLUS, Str: "+", Loc: l.start}
	case '-':
		return Token{Type: tokMINUS, Str: "-", Loc: l.start}
	case '/':
		return Token{Type: tokDIV, Str: "/", Loc: l.start}
	case '%':
		return Token{Type: tokMOD, Str: "%", Loc: l.start}
	case '=':
		return Token{Type: tokEQ, Str: "=", Loc: l.start}
	case '<':
		return Token{Type: tokLT, Str: "<", Loc: l.start}
	case '>':
		return Token{Type: tokGT, Str: ">", Loc: l.start}
	case '&':
		return Token{Type: tokBITAND, Str: "&", Loc: l.start}
	case '|':
		return Token{Type: tokBITOR, Str: "|", Loc: l.start}
	case '^':
		return Token{Type: tokBITXOR, Str: "^", Loc: l.start}
	case '~':
		return Token{Type: tokBITNOT, Str: "~", Loc: l.start}
	case '?':
		return Token{Type: tokQUESTION, Str: "?", Loc: l.start}
	default:
		_, size := utf8.DecodeRuneInString(l.input[l.start:])
		_ = size
		l.Err = fmt.Errorf("unexpected character %q at position %d", string(ch), l.start)
		return Token{Type: tokEOF, Loc: l.start}
	}
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentContinue(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
