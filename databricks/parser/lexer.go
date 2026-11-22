package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// lex tokenizes a Databricks/Spark SQL string into a slice of Tokens ending
// with tokEOF. It follows SqlBaseLexer.g4, including nested block comments,
// /*+ ... */ optimizer hints, backtick identifiers, raw strings (R'...' /
// R"..."), numeric suffixes (L/S/Y/F/D/BD), and URI identifiers.
func lex(src string) ([]Token, error) {
	var toks []Token
	n := len(src)
	i := 0
	inHint := false

outer:
	for i < n {
		c := src[i]

		// ASCII whitespace fast-path.
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			i++
			continue
		}
		if c >= utf8.RuneSelf {
			r, sz := utf8.DecodeRuneInString(src[i:])
			if unicode.IsSpace(r) {
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
				toks = append(toks, Token{Kind: tokHintStart, Text: "/*+", Start: i, End: i + 3})
				i += 3
				inHint = true
				continue
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
				return nil, &ParseError{Message: "unclosed block comment", Position: start}
			}
			continue
		}

		// End of hint */.
		if c == '*' && i+1 < n && src[i+1] == '/' && inHint {
			toks = append(toks, Token{Kind: tokHintEnd, Text: "*/", Start: i, End: i + 2})
			i += 2
			inHint = false
			continue
		}

		// Raw string literal: R'...' or R"..." (case-insensitive R).
		if (c == 'R' || c == 'r') && i+1 < n && (src[i+1] == '\'' || src[i+1] == '"') {
			start := i
			q := src[i+1]
			i += 2
			for i < n && src[i] != q {
				i++
			}
			if i >= n {
				return nil, &ParseError{Message: "unterminated raw string literal", Position: start}
			}
			i++
			toks = append(toks, Token{Kind: tokString, Text: src[start:i], Start: start, End: i})
			continue
		}

		// Single-quoted string literal: '...' with backslash escapes.
		if c == '\'' {
			start := i
			i++
			for i < n {
				if src[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if src[i] == '\'' {
					i++
					toks = append(toks, Token{Kind: tokString, Text: src[start:i], Start: start, End: i})
					continue outer
				}
				i++
			}
			return nil, &ParseError{Message: "unterminated string literal", Position: start}
		}

		// Double-quoted string / identifier: "..." with backslash escapes.
		if c == '"' {
			start := i
			i++
			for i < n {
				if src[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if src[i] == '"' {
					i++
					toks = append(toks, Token{Kind: tokDoubleQuoted, Text: src[start:i], Start: start, End: i})
					continue outer
				}
				i++
			}
			return nil, &ParseError{Message: "unterminated double-quoted literal", Position: start}
		}

		// Backtick-quoted identifier: `...` with `` escape.
		if c == '`' {
			start := i
			i++
			for i < n {
				if src[i] == '`' {
					if i+1 < n && src[i+1] == '`' {
						i += 2
						continue
					}
					i++
					toks = append(toks, Token{Kind: tokBackquoted, Text: src[start:i], Start: start, End: i})
					continue outer
				}
				i++
			}
			return nil, &ParseError{Message: "unterminated backquoted identifier", Position: start}
		}

		// Numeric literal or identifier starting with digit (Spark IDENTIFIER allows leading digits).
		if isASCIIDigit(c) || (c == '.' && i+1 < n && isASCIIDigit(src[i+1])) {
			tok := scanNumberOrIdent(src, i)
			toks = append(toks, tok)
			i = tok.End
			continue
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
						i = j
						toks = append(toks, Token{Kind: tokIdent, Text: src[start:i], Start: start, End: i})
						continue
					}
				}
			}
			word := src[start:i]
			kind := lookupKeyword(word)
			toks = append(toks, Token{Kind: kind, Text: word, Start: start, End: i})
			continue
		}

		// Multi-character and single-character operators/punctuation.
		start := i
		switch c {
		case ';':
			i++
			toks = append(toks, Token{Kind: tokSemicolon, Text: ";", Start: start, End: i})
		case '(':
			i++
			toks = append(toks, Token{Kind: tokLParen, Text: "(", Start: start, End: i})
		case ')':
			i++
			toks = append(toks, Token{Kind: tokRParen, Text: ")", Start: start, End: i})
		case ',':
			i++
			toks = append(toks, Token{Kind: tokComma, Text: ",", Start: start, End: i})
		case '.':
			i++
			toks = append(toks, Token{Kind: tokDot, Text: ".", Start: start, End: i})
		case '[':
			i++
			toks = append(toks, Token{Kind: tokLBracket, Text: "[", Start: start, End: i})
		case ']':
			i++
			toks = append(toks, Token{Kind: tokRBracket, Text: "]", Start: start, End: i})
		case '+':
			i++
			toks = append(toks, Token{Kind: tokPlus, Text: "+", Start: start, End: i})
		case '-':
			if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokArrow, Text: "->", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokMinus, Text: "-", Start: start, End: i})
			}
		case '*':
			i++
			toks = append(toks, Token{Kind: tokAsterisk, Text: "*", Start: start, End: i})
		case '/':
			i++
			toks = append(toks, Token{Kind: tokSlash, Text: "/", Start: start, End: i})
		case '%':
			i++
			toks = append(toks, Token{Kind: tokPercent, Text: "%", Start: start, End: i})
		case '~':
			i++
			toks = append(toks, Token{Kind: tokTilde, Text: "~", Start: start, End: i})
		case '&':
			i++
			toks = append(toks, Token{Kind: tokAmpersand, Text: "&", Start: start, End: i})
		case '^':
			i++
			toks = append(toks, Token{Kind: tokHat, Text: "^", Start: start, End: i})
		case '?':
			i++
			toks = append(toks, Token{Kind: tokQuestion, Text: "?", Start: start, End: i})
		case ':':
			if i+1 < n && src[i+1] == ':' {
				i += 2
				toks = append(toks, Token{Kind: tokDoubleColon, Text: "::", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokColon, Text: ":", Start: start, End: i})
			}
		case '|':
			if i+1 < n && src[i+1] == '|' {
				i += 2
				toks = append(toks, Token{Kind: tokConcatPipe, Text: "||", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokOperatorPipe, Text: "|>", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokPipe, Text: "|", Start: start, End: i})
			}
		case '=':
			if i+1 < n && src[i+1] == '=' {
				i += 2
				toks = append(toks, Token{Kind: tokEq, Text: "==", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokFatArrow, Text: "=>", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokEq, Text: "=", Start: start, End: i})
			}
		case '!':
			if i+1 < n && src[i+1] == '=' {
				i += 2
				toks = append(toks, Token{Kind: tokNeqBang, Text: "!=", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokLTE, Text: "!>", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '<' {
				i += 2
				toks = append(toks, Token{Kind: tokGTE, Text: "!<", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokBang, Text: "!", Start: start, End: i})
			}
		case '<':
			if i+2 < n && src[i+1] == '=' && src[i+2] == '>' {
				i += 3
				toks = append(toks, Token{Kind: tokNullSafeEq, Text: "<=>", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '=' {
				i += 2
				toks = append(toks, Token{Kind: tokLTE, Text: "<=", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokNeq, Text: "<>", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '<' {
				i += 2
				toks = append(toks, Token{Kind: tokShiftLeft, Text: "<<", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokLT, Text: "<", Start: start, End: i})
			}
		case '>':
			if i+2 < n && src[i+1] == '>' && src[i+2] == '>' {
				i += 3
				toks = append(toks, Token{Kind: tokShiftRightU, Text: ">>>", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '>' {
				i += 2
				toks = append(toks, Token{Kind: tokShiftRight, Text: ">>", Start: start, End: i})
			} else if i+1 < n && src[i+1] == '=' {
				i += 2
				toks = append(toks, Token{Kind: tokGTE, Text: ">=", Start: start, End: i})
			} else {
				i++
				toks = append(toks, Token{Kind: tokGT, Text: ">", Start: start, End: i})
			}
		default:
			_, sz := utf8.DecodeRuneInString(src[i:])
			return nil, &ParseError{Message: fmt.Sprintf("unexpected character %q", src[start:start+sz]), Position: start}
		}
	}

	toks = append(toks, Token{Kind: tokEOF, Start: n, End: n})
	return toks, nil
}

func scanNumberOrIdent(src string, start int) Token {
	n := len(src)
	i := start
	hasDot := false
	hasExp := false
	hasExpSign := false

	if src[i] == '.' {
		hasDot = true
		i++
		for i < n && isASCIIDigit(src[i]) {
			i++
		}
	} else {
		for i < n && isASCIIDigit(src[i]) {
			i++
		}
		if i < n && src[i] == '.' {
			// Check Spark isValidDecimal: decimal followed by letter/digit/_ is not valid decimal,
			// unless followed by digits or exponent or F/D/BD suffix.
			j := i + 1
			for j < n && isASCIIDigit(src[j]) {
				j++
			}
			// Accept dot as part of decimal if next is not ident-start (except exponent/suffix) or if digits follow.
			if j > i+1 || j >= n || !isIdentStartByte(src, j) || isNumericSuffixOrExp(src, j) {
				hasDot = true
				i = j
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
		return Token{Kind: tokIdent, Text: src[start:i], Start: start, End: i}
	}

	kind := tokInt
	if hasExp {
		kind = tokExponent
	} else if hasDot {
		kind = tokDecimal
	}
	return Token{Kind: kind, Text: src[start:i], Start: start, End: i}
}

func isNumericSuffixOrExp(src string, i int) bool {
	if i >= len(src) {
		return false
	}
	c := upperASCII(src[i])
	if c == 'E' {
		j := i + 1
		if j < len(src) && (src[j] == '+' || src[j] == '-') {
			j++
		}
		return j < len(src) && isASCIIDigit(src[j])
	}
	if c == 'F' || c == 'D' {
		return !isIdentPartAt(src, i+1)
	}
	if c == 'B' && i+1 < len(src) && upperASCII(src[i+1]) == 'D' {
		return !isIdentPartAt(src, i+2)
	}
	return false
}

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
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

func isIdentPartRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
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
