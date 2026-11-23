package parser

import (
	"unicode"
	"unicode/utf8"
)

// Lexer is a streaming lexical scanner for Databricks/Spark SQL following
// SqlBaseLexer.g4, including nested block comments, /*+ ... */ optimizer hints,
// backtick identifiers, raw strings (R'...' / R"..."), numeric suffixes
// (L/S/Y/F/D/BD), and URI identifiers.
type Lexer struct {
	src    string
	pos    int
	inHint bool
	err    error
}

// NewLexer returns a new Lexer scanning src from byte offset 0.
func NewLexer(src string) *Lexer {
	return &Lexer{src: src}
}

// Err returns the first fatal lexical error encountered (such as an unclosed
// block comment), or nil if none occurred.
func (l *Lexer) Err() error {
	return l.err
}

// Tokenize lexes a Databricks/Spark SQL string into a slice of Tokens ending
// with TokEOF.
func Tokenize(src string) ([]Token, error) {
	return lex(src)
}

// lex tokenizes a Databricks/Spark SQL string into a slice of Tokens ending
// with tokEOF.
func lex(src string) ([]Token, error) {
	l := NewLexer(src)
	var toks []Token
	for {
		tok := l.NextToken()
		if l.err != nil {
			return nil, l.err
		}
		toks = append(toks, tok)
		if tok.Kind == tokEOF {
			return toks, nil
		}
	}
}

// NextToken returns the next token from the input, or TokEOF at end of input
// (or after a fatal lexical error, inspectable via Err).
func (l *Lexer) NextToken() Token {
	if l.err != nil {
		return Token{Kind: tokEOF, Start: len(l.src), End: len(l.src)}
	}
	src := l.src
	n := len(src)
	i := l.pos

	for i < n {
		c := src[i]

		// ASCII whitespace fast-path.
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			i++
			continue
		}
		if c >= utf8.RuneSelf {
			r, sz := utf8.DecodeRuneInString(src[i:])
			if isSparkSpace(r) {
				i += sz
				continue
			}
		}

		// Line comment: -- ... (with \\\n line continuation as in Spark).
		if c == '-' && i+1 < n && src[i+1] == '-' {
			i += 2
			for i < n {
				if src[i] == '\\' && i+1 < n && src[i+1] == '\n' {
					i += 2
					continue
				}
				if src[i] == '\n' || src[i] == '\r' {
					break
				}
				i++
			}
			continue
		}

		// Block comment or optimizer hint /*+ ... */.
		if c == '/' && i+1 < n && src[i+1] == '*' {
			if i+2 < n && src[i+2] == '+' {
				l.pos = i + 3
				l.inHint = true
				return Token{Kind: tokHintStart, Text: "/*+", Start: i, End: i + 3}
			}
			start := i
			depth := 1
			i += 2
			for i < n && depth > 0 {
				if src[i] == '/' && i+1 < n && src[i+1] == '*' {
					depth++
					i += 2
				} else if src[i] == '*' && i+1 < n && src[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth != 0 {
				l.pos = n
				l.err = &ParseError{Message: "unclosed block comment", Position: start}
				return Token{Kind: tokEOF, Start: n, End: n}
			}
			continue
		}

		// End of hint */.
		if c == '*' && i+1 < n && src[i+1] == '/' && l.inHint {
			l.pos = i + 2
			l.inHint = false
			return Token{Kind: tokHintEnd, Text: "*/", Start: i, End: i + 2}
		}

		// Raw string literal: R'...' or R"..." (case-insensitive R).
		// If no closing quote exists, fall through so `R` lexes as an
		// identifier and the quote lexes as tokInvalid (matching Spark's
		// ANTLR UNRECOGNIZED fallback for raw SET/ADD/REFRESH tails).
		if (c == 'R' || c == 'r') && i+1 < n && (src[i+1] == '\'' || src[i+1] == '"') {
			q := src[i+1]
			j := i + 2
			for j < n && src[j] != q {
				j++
			}
			if j < n {
				j++
				l.pos = j
				return Token{Kind: tokString, Text: src[i:j], Start: i, End: j}
			}
		}

		// Single-quoted string literal: '...' with backslash escapes.
		// Unmatched quotes fall back to a single-byte tokInvalid token
		// (Spark's ANTLR UNRECOGNIZED) so raw SET/ADD/REFRESH tails accept
		// values containing apostrophes while SQL expressions reject them.
		if c == '\'' {
			start := i
			j := i + 1
			for j < n {
				if src[j] == '\\' && j+1 < n {
					j += 2
					continue
				}
				if src[j] == '\'' {
					j++
					l.pos = j
					return Token{Kind: tokString, Text: src[start:j], Start: start, End: j}
				}
				j++
			}
			l.pos = start + 1
			return Token{Kind: tokInvalid, Text: src[start : start+1], Start: start, End: start + 1}
		}

		// Double-quoted string literal: "..." with backslash escapes.
		if c == '"' {
			start := i
			j := i + 1
			for j < n {
				if src[j] == '\\' && j+1 < n {
					j += 2
					continue
				}
				if src[j] == '"' {
					j++
					l.pos = j
					return Token{Kind: tokDoubleQuoted, Text: src[start:j], Start: start, End: j}
				}
				j++
			}
			l.pos = start + 1
			return Token{Kind: tokInvalid, Text: src[start : start+1], Start: start, End: start + 1}
		}

		// Backtick-quoted identifier: `...` with `` escape.
		if c == '`' {
			start := i
			j := i + 1
			for j < n {
				if src[j] == '`' {
					if j+1 < n && src[j+1] == '`' {
						j += 2
						continue
					}
					j++
					l.pos = j
					return Token{Kind: tokBackquoted, Text: src[start:j], Start: start, End: j}
				}
				j++
			}
			l.pos = start + 1
			return Token{Kind: tokInvalid, Text: src[start : start+1], Start: start, End: start + 1}
		}

		// Numeric literal or identifier starting with digit (Spark IDENTIFIER allows leading digits).
		if isASCIIDigit(c) || (c == '.' && i+1 < n && isASCIIDigit(src[i+1])) {
			tok := scanNumberOrIdent(src, i)
			l.pos = tok.End
			return tok
		}

		// Identifier or keyword (including Unicode letters and scheme:// URI paths).
		if isIdentStartByte(src, i) {
			start := i
			for i < n {
				r, sz := utf8.DecodeRuneInString(src[i:])
				if isIdentPartRune(r) {
					i += sz
				} else {
					break
				}
			}
			// Spark URI path identifier: UNICODE_LETTER+ '://' ...
			if i+3 <= n && src[i:i+3] == "://" {
				allLetters := true
				for _, r := range src[start:i] {
					if !unicode.IsLetter(r) {
						allLetters = false
						break
					}
				}
				if allLetters {
					j := i + 3
					for j < n {
						r, sz := utf8.DecodeRuneInString(src[j:])
						if isURIPartRune(r) {
							j += sz
						} else {
							break
						}
					}
					if j > i+3 {
						l.pos = j
						return Token{Kind: tokRawWord, Text: src[start:j], Start: start, End: j}
					}
				}
			}
			word := src[start:i]
			kind := tokRawWord
			if isASCIIIdent(word) {
				kind = lookupKeyword(word)
			}
			l.pos = i
			return Token{Kind: kind, Text: word, Start: start, End: i}
		}

		// Multi-character and single-character operators/punctuation.
		start := i
		switch c {
		case ';':
			l.pos = i + 1
			return Token{Kind: tokSemicolon, Text: ";", Start: start, End: i + 1}
		case '(':
			l.pos = i + 1
			return Token{Kind: tokLParen, Text: "(", Start: start, End: i + 1}
		case ')':
			l.pos = i + 1
			return Token{Kind: tokRParen, Text: ")", Start: start, End: i + 1}
		case ',':
			l.pos = i + 1
			return Token{Kind: tokComma, Text: ",", Start: start, End: i + 1}
		case '.':
			l.pos = i + 1
			return Token{Kind: tokDot, Text: ".", Start: start, End: i + 1}
		case '[':
			l.pos = i + 1
			return Token{Kind: tokLBracket, Text: "[", Start: start, End: i + 1}
		case ']':
			l.pos = i + 1
			return Token{Kind: tokRBracket, Text: "]", Start: start, End: i + 1}
		case '+':
			l.pos = i + 1
			return Token{Kind: tokPlus, Text: "+", Start: start, End: i + 1}
		case '-':
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokArrow, Text: "->", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokMinus, Text: "-", Start: start, End: i + 1}
		case '*':
			l.pos = i + 1
			return Token{Kind: tokAsterisk, Text: "*", Start: start, End: i + 1}
		case '/':
			l.pos = i + 1
			return Token{Kind: tokSlash, Text: "/", Start: start, End: i + 1}
		case '%':
			l.pos = i + 1
			return Token{Kind: tokPercent, Text: "%", Start: start, End: i + 1}
		case '~':
			l.pos = i + 1
			return Token{Kind: tokTilde, Text: "~", Start: start, End: i + 1}
		case '&':
			l.pos = i + 1
			return Token{Kind: tokAmpersand, Text: "&", Start: start, End: i + 1}
		case '^':
			l.pos = i + 1
			return Token{Kind: tokHat, Text: "^", Start: start, End: i + 1}
		case '?':
			l.pos = i + 1
			return Token{Kind: tokQuestion, Text: "?", Start: start, End: i + 1}
		case ':':
			if i+1 < n && src[i+1] == ':' {
				l.pos = i + 2
				return Token{Kind: tokDoubleColon, Text: "::", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokColon, Text: ":", Start: start, End: i + 1}
		case '|':
			if i+1 < n && src[i+1] == '|' {
				l.pos = i + 2
				return Token{Kind: tokConcatPipe, Text: "||", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokOperatorPipe, Text: "|>", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokPipe, Text: "|", Start: start, End: i + 1}
		case '=':
			if i+1 < n && src[i+1] == '=' {
				l.pos = i + 2
				return Token{Kind: tokEq, Text: "==", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokFatArrow, Text: "=>", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokEq, Text: "=", Start: start, End: i + 1}
		case '!':
			if i+1 < n && src[i+1] == '=' {
				l.pos = i + 2
				return Token{Kind: tokNeqBang, Text: "!=", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokLTE, Text: "!>", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '<' {
				l.pos = i + 2
				return Token{Kind: tokGTE, Text: "!<", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokBang, Text: "!", Start: start, End: i + 1}
		case '<':
			if i+2 < n && src[i+1] == '=' && src[i+2] == '>' {
				l.pos = i + 3
				return Token{Kind: tokNullSafeEq, Text: "<=>", Start: start, End: i + 3}
			}
			if i+1 < n && src[i+1] == '=' {
				l.pos = i + 2
				return Token{Kind: tokLTE, Text: "<=", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokNeq, Text: "<>", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '<' {
				l.pos = i + 2
				return Token{Kind: tokShiftLeft, Text: "<<", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokLT, Text: "<", Start: start, End: i + 1}
		case '>':
			if i+2 < n && src[i+1] == '>' && src[i+2] == '>' {
				l.pos = i + 3
				return Token{Kind: tokShiftRightU, Text: ">>>", Start: start, End: i + 3}
			}
			if i+1 < n && src[i+1] == '>' {
				l.pos = i + 2
				return Token{Kind: tokShiftRight, Text: ">>", Start: start, End: i + 2}
			}
			if i+1 < n && src[i+1] == '=' {
				l.pos = i + 2
				return Token{Kind: tokGTE, Text: ">=", Start: start, End: i + 2}
			}
			l.pos = i + 1
			return Token{Kind: tokGT, Text: ">", Start: start, End: i + 1}
		default:
			// Match Spark's UNRECOGNIZED lexer token so raw wildcard tails
			// (such as SET spark.test = abc@def or ${var}) can consume arbitrary
			// characters while ordinary SQL grammar rules reject tokInvalid.
			_, sz := utf8.DecodeRuneInString(src[i:])
			l.pos = i + sz
			return Token{Kind: tokInvalid, Text: src[start : i+sz], Start: start, End: i + sz}
		}
	}

	l.pos = n
	return Token{Kind: tokEOF, Start: n, End: n}
}

func scanNumberOrIdent(src string, start int) Token {
	n := len(src)
	i := start
	hasDot := false
	hasExp := false
	hasExpSign := false

	dotPos := -1
	if src[i] == '.' {
		hasDot = true
		dotPos = i
		i++
		for i < n && isASCIIDigit(src[i]) {
			i++
		}
	} else {
		for i < n && isASCIIDigit(src[i]) {
			i++
		}
		if i < n && src[i] == '.' {
			dotPos = i
			hasDot = true
			i++
			for i < n && isASCIIDigit(src[i]) {
				i++
			}
		}
	}

	// Optional exponent: E[+-]?DIGIT+
	if i < n && (src[i] == 'E' || src[i] == 'e') {
		j := i + 1
		signed := false
		if j < n && (src[j] == '+' || src[j] == '-') {
			signed = true
			j++
		}
		if j < n && isASCIIDigit(src[j]) {
			hasExp = true
			hasExpSign = signed
			i = j
			for i < n && isASCIIDigit(src[i]) {
				i++
			}
		}
	}

	// Numeric type suffixes (L, S, Y, F, D, BD), case-insensitive, only if not followed by ident part.
	if i < n {
		c := upperASCII(src[i])
		if c == 'B' && i+1 < n && upperASCII(src[i+1]) == 'D' && !isIdentPartAt(src, i+2) {
			return Token{Kind: tokBigDecimal, Text: src[start : i+2], Start: start, End: i + 2}
		}
		if !isIdentPartAt(src, i+1) {
			switch c {
			case 'F':
				return Token{Kind: tokFloat, Text: src[start : i+1], Start: start, End: i + 1}
			case 'D':
				return Token{Kind: tokDouble, Text: src[start : i+1], Start: start, End: i + 1}
			case 'L':
				if !hasDot && !hasExp {
					return Token{Kind: tokBigInt, Text: src[start : i+1], Start: start, End: i + 1}
				}
			case 'S':
				if !hasDot && !hasExp {
					return Token{Kind: tokSmallInt, Text: src[start : i+1], Start: start, End: i + 1}
				}
			case 'Y':
				if !hasDot && !hasExp {
					return Token{Kind: tokTinyInt, Text: src[start : i+1], Start: start, End: i + 1}
				}
			}
		}
	}

	// If no dot and no sign in exponent and followed by identifier characters,
	// Spark's IDENTIFIER rule wins by longest match (e.g. 123abc, 1e2data).
	if !hasDot && !hasExpSign && isIdentPartAt(src, i) {
		for i < n {
			r, sz := utf8.DecodeRuneInString(src[i:])
			if isIdentPartRune(r) {
				i += sz
			} else {
				break
			}
		}
		word := src[start:i]
		kind := tokRawWord
		if isASCIIIdent(word) {
			kind = tokIdent
		}
		return Token{Kind: kind, Text: word, Start: start, End: i}
	}

	// Spark's isValidDecimal(): a decimal/exponent literal with a dot cannot be
	// immediately followed by an identifier character (e.g. `1.0L`, `1.0Y`,
	// `t.123abc`, `t.1e2data`). Fall back to DOT (leading-dot) or INT (digit-leading).
	if hasDot && isIdentPartAt(src, i) {
		if src[start] == '.' {
			return Token{Kind: tokDot, Text: ".", Start: start, End: start + 1}
		}
		return Token{Kind: tokInt, Text: src[start:dotPos], Start: start, End: dotPos}
	}

	kind := tokInt
	if hasExp {
		kind = tokExponent
	} else if hasDot {
		kind = tokDecimal
	}
	return Token{Kind: kind, Text: src[start:i], Start: start, End: i}
}

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// isASCIIIdent reports whether s contains only Spark's accepted unquoted
// identifier alphabet ([A-Za-z0-9_]), matching PostProcessor.exitUnquotedIdentifier.
func isASCIIIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_') {
			return false
		}
	}
	return true
}

func upperASCII(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}
	return b
}

func isIdentStartByte(src string, i int) bool {
	if i >= len(src) {
		return false
	}
	b := src[i]
	if b < utf8.RuneSelf {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
	}
	r, _ := utf8.DecodeRuneInString(src[i:])
	return unicode.IsLetter(r)
}

func isIdentPartAt(src string, i int) bool {
	if i >= len(src) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(src[i:])
	return isIdentPartRune(r)
}

// isSparkSpace matches Spark's SqlBaseLexer.g4 WS rule exactly (excluding
// U+0085 and U+2029, which unicode.IsSpace accepts but Spark rejects).
func isSparkSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\f', '\r', '\v',
		'\u00A0', '\u1680',
		'\u2000', '\u2001', '\u2002', '\u2003', '\u2004', '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200A',
		'\u2028', '\u202F', '\u205F', '\u3000':
		return true
	}
	return false
}

func isIdentPartRune(r rune) bool {
	return r == '_' || (r >= '0' && r <= '9') || unicode.IsLetter(r)
}

func isURIPartRune(r rune) bool {
	if isIdentPartRune(r) {
		return true
	}
	switch r {
	case '/', '-', '.', '?', '=', '&', '#', '%':
		return true
	}
	return false
}
