package parser

import (
	"fmt"
	"strings"
)

// Lexer tokenizes HiveQL input according to Apache Hive 4.0.1's HiveLexer.g and
// HiveLexerParent.g (with case-insensitive keyword/token matching like
// ANTLRNoCaseStringStream).
type Lexer struct {
	input string
	pos   int
	err   *ParseError
}

// NewLexer returns a Lexer positioned at the start of input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Err returns the first lexical error encountered by the lexer, if any.
func (l *Lexer) Err() error {
	if l.err == nil {
		return nil
	}
	return l.err
}

func (l *Lexer) fail(pos int, msg string, end int) Token {
	if l.err == nil {
		l.err = &ParseError{Message: msg, Position: pos}
	}
	l.pos = end
	return Token{Type: TokInvalid, Str: l.input[pos:end], Loc: Loc{Start: pos, End: end}}
}

// Tokenize lexes input into a slice of tokens ending with TokEOF.
func Tokenize(input string) ([]Token, error) {
	return lex(input)
}

func lex(input string) ([]Token, error) {
	l := NewLexer(input)
	var toks []Token
	for {
		tok := l.NextToken()
		if tok.Type == TokInvalid {
			return nil, l.err
		}
		toks = append(toks, tok)
		if tok.Type == TokEOF {
			return toks, nil
		}
	}
}

// NextToken returns the next token from the input. On a lexical error, it
// records the first error (accessible via Err), advances past the offending
// span, and returns a TokInvalid token. At end of input, it returns TokEOF.
func (l *Lexer) NextToken() Token {
	input := l.input
	n := len(input)
	i := l.pos
	for i < n {
		c := input[i]
		// Whitespace: (' '|'\r'|'\t'|'\n') -> HIDDEN
		if c == ' ' || c == '\r' || c == '\t' || c == '\n' {
			i++
			continue
		}
		// Line comment: '--' (~('\n'|'\r'))* -> HIDDEN
		if c == '-' && i+1 < n && input[i+1] == '-' {
			i += 2
			for i < n && input[i] != '\n' && input[i] != '\r' {
				i++
			}
			continue
		}
		// Block comment / QUERY_HINT: '/*' (QUERY_HINT | .)* '*/'
		if c == '/' && i+1 < n && input[i+1] == '*' {
			start := i
			end, err := scanBlockComment(input, i)
			if err != nil {
				return l.fail(start, "unterminated block comment", n)
			}
			// If third byte is '+' (and comment is at least "/*+*/"), emit QUERY_HINT.
			if end-start >= 5 && input[start+2] == '+' {
				l.pos = end
				return Token{
					Type: TokQueryHint,
					Str:  input[start+3 : end-2],
					Loc:  Loc{Start: start, End: end},
				}
			}
			i = end
			continue
		}

		// StringLiteral: ( '"' ( ~('"'|'\\') | '\\' . )* '"' | '\'' ( ~('\''|'\\') | '\\' . )* '\'' )+
		if c == '\'' || c == '"' {
			start := i
			for i < n && (input[i] == '\'' || input[i] == '"') {
				q := input[i]
				i++
				closed := false
				for i < n {
					if input[i] == '\\' {
						if i+1 >= n {
							return l.fail(start, "unterminated string literal", n)
						}
						i += 2
						continue
					}
					if input[i] == q {
						i++
						closed = true
						break
					}
					i++
				}
				if !closed {
					return l.fail(start, "unterminated string literal", n)
				}
			}
			l.pos = i
			return Token{
				Type: TokStringLiteral,
				Str:  input[start:i],
				Loc:  Loc{Start: start, End: i},
			}
		}

		// QuotedIdentifier: '`' ( '``' | ~'`' )+ '`'
		if c == '`' {
			start := i
			i++
			var sb strings.Builder
			closed := false
			for i < n {
				if input[i] == '`' {
					if i+1 < n && input[i+1] == '`' {
						sb.WriteByte('`')
						i += 2
						continue
					}
					i++
					closed = true
					break
				}
				sb.WriteByte(input[i])
				i++
			}
			if !closed {
				return l.fail(start, "invalid backtick-quoted identifier", n)
			}
			if sb.Len() == 0 {
				return l.fail(start, "invalid backtick-quoted identifier", i)
			}
			l.pos = i
			return Token{
				Type:   TokIdentifier,
				Str:    sb.String(),
				Quoted: true,
				Loc:    Loc{Start: start, End: i},
			}
		}

		// CharSetName: '_' (Letter | Digit | '_' | '-' | '.' | ':')+
		if c == '_' {
			start := i
			i++
			for i < n && isCharSetNameRest(input[i]) {
				i++
			}
			if i == start+1 {
				return l.fail(start, fmt.Sprintf("unexpected character %q", c), start+1)
			}
			l.pos = i
			return Token{
				Type: TokCharSetName,
				Str:  input[start:i],
				Loc:  Loc{Start: start, End: i},
			}
		}

		// Special dollar-delimited keywords: $KEY$, $VALUE$, $ELEM$
		if c == '$' {
			if rest := matchDollarKeyword(input, i); rest > 0 {
				tt, _ := lookupKeyword(input[i : i+rest])
				l.pos = i + rest
				return Token{Type: tt, Str: input[i : i+rest], Loc: Loc{Start: i, End: i + rest}}
			}
			l.pos = i + 1
			return Token{Type: TokDollar, Str: "$", Loc: Loc{Start: i, End: i + 1}}
		}

		// Digits: CharSetLiteral (0X...), IntegralLiteral, NumberLiteral, ByteLengthLiteral, Number, or Identifier starting with Digit.
		if isDigit(c) {
			tok, next := scanDigitStart(input, i)
			l.pos = next
			return tok
		}

		// Letters: Keyword (including hyphenated FAST-FORWARD, CHERRY-PICK, ORPHAN-FILES) or Identifier.
		if isLetter(c) {
			start := i
			i++
			for i < n && (isLetter(input[i]) || isDigit(input[i]) || input[i] == '_') {
				i++
			}
			// Check hyphenated Hive 4 keywords.
			if i < n && input[i] == '-' {
				word := strings.ToUpper(input[start:i])
				if word == "FAST" || word == "CHERRY" || word == "ORPHAN" {
					j := i + 1
					for j < n && (isLetter(input[j]) || isDigit(input[j]) || input[j] == '_') {
						j++
					}
					if tt, ok := lookupKeyword(input[start:j]); ok {
						l.pos = j
						return Token{Type: tt, Str: input[start:j], Loc: Loc{Start: start, End: j}}
					}
				}
			}
			word := input[start:i]
			l.pos = i
			if tt, ok := lookupKeyword(word); ok {
				return Token{Type: tt, Str: word, Loc: Loc{Start: start, End: i}}
			}
			return Token{Type: TokIdentifier, Str: word, Loc: Loc{Start: start, End: i}}
		}

		// Multi-character and single-character operators/punctuation.
		start := i
		switch c {
		case '<':
			if i+2 < n && input[i+1] == '=' && input[i+2] == '>' {
				l.pos = i + 3
				return Token{Type: TokEqualNS, Str: "<=>", Loc: Loc{Start: start, End: i + 3}}
			} else if i+1 < n && input[i+1] == '=' {
				l.pos = i + 2
				return Token{Type: TokLessThanOrEqualTo, Str: "<=", Loc: Loc{Start: start, End: i + 2}}
			} else if i+1 < n && input[i+1] == '>' {
				l.pos = i + 2
				return Token{Type: TokNotEqual, Str: "<>", Loc: Loc{Start: start, End: i + 2}}
			}
			l.pos = i + 1
			return Token{Type: TokLessThan, Str: "<", Loc: Loc{Start: start, End: i + 1}}
		case '>':
			if i+1 < n && input[i+1] == '=' {
				l.pos = i + 2
				return Token{Type: TokGreaterThanOrEqualTo, Str: ">=", Loc: Loc{Start: start, End: i + 2}}
			}
			l.pos = i + 1
			return Token{Type: TokGreaterThan, Str: ">", Loc: Loc{Start: start, End: i + 1}}
		case '=':
			if i+1 < n && input[i+1] == '=' {
				l.pos = i + 2
				return Token{Type: TokEqual, Str: "==", Loc: Loc{Start: start, End: i + 2}}
			}
			l.pos = i + 1
			return Token{Type: TokEqual, Str: "=", Loc: Loc{Start: start, End: i + 1}}
		case '!':
			if i+1 < n && input[i+1] == '=' {
				l.pos = i + 2
				return Token{Type: TokNotEqual, Str: "!=", Loc: Loc{Start: start, End: i + 2}}
			}
			// HiveLexerParent.g: KW_NOT : 'NOT' | '!';
			l.pos = i + 1
			return Token{Type: KW_NOT, Str: "!", Loc: Loc{Start: start, End: i + 1}}
		case '|':
			if i+1 < n && input[i+1] == '|' {
				l.pos = i + 2
				return Token{Type: TokConcatenate, Str: "||", Loc: Loc{Start: start, End: i + 2}}
			}
			l.pos = i + 1
			return Token{Type: TokBitwiseOr, Str: "|", Loc: Loc{Start: start, End: i + 1}}
		case '.':
			l.pos = i + 1
			return Token{Type: TokDot, Str: ".", Loc: Loc{Start: start, End: i + 1}}
		case ':':
			l.pos = i + 1
			return Token{Type: TokColon, Str: ":", Loc: Loc{Start: start, End: i + 1}}
		case ',':
			l.pos = i + 1
			return Token{Type: TokComma, Str: ",", Loc: Loc{Start: start, End: i + 1}}
		case ';':
			l.pos = i + 1
			return Token{Type: TokSemicolon, Str: ";", Loc: Loc{Start: start, End: i + 1}}
		case '(':
			l.pos = i + 1
			return Token{Type: TokLParen, Str: "(", Loc: Loc{Start: start, End: i + 1}}
		case ')':
			l.pos = i + 1
			return Token{Type: TokRParen, Str: ")", Loc: Loc{Start: start, End: i + 1}}
		case '[':
			l.pos = i + 1
			return Token{Type: TokLSquare, Str: "[", Loc: Loc{Start: start, End: i + 1}}
		case ']':
			l.pos = i + 1
			return Token{Type: TokRSquare, Str: "]", Loc: Loc{Start: start, End: i + 1}}
		case '{':
			l.pos = i + 1
			return Token{Type: TokLCurly, Str: "{", Loc: Loc{Start: start, End: i + 1}}
		case '}':
			l.pos = i + 1
			return Token{Type: TokRCurly, Str: "}", Loc: Loc{Start: start, End: i + 1}}
		case '/':
			l.pos = i + 1
			return Token{Type: TokDivide, Str: "/", Loc: Loc{Start: start, End: i + 1}}
		case '+':
			l.pos = i + 1
			return Token{Type: TokPlus, Str: "+", Loc: Loc{Start: start, End: i + 1}}
		case '-':
			l.pos = i + 1
			return Token{Type: TokMinus, Str: "-", Loc: Loc{Start: start, End: i + 1}}
		case '*':
			l.pos = i + 1
			return Token{Type: TokStar, Str: "*", Loc: Loc{Start: start, End: i + 1}}
		case '%':
			l.pos = i + 1
			return Token{Type: TokMod, Str: "%", Loc: Loc{Start: start, End: i + 1}}
		case '&':
			l.pos = i + 1
			return Token{Type: TokAmpersand, Str: "&", Loc: Loc{Start: start, End: i + 1}}
		case '~':
			l.pos = i + 1
			return Token{Type: TokTilde, Str: "~", Loc: Loc{Start: start, End: i + 1}}
		case '^':
			l.pos = i + 1
			return Token{Type: TokBitwiseXor, Str: "^", Loc: Loc{Start: start, End: i + 1}}
		case '?':
			l.pos = i + 1
			return Token{Type: TokQuestion, Str: "?", Loc: Loc{Start: start, End: i + 1}}
		default:
			return l.fail(start, fmt.Sprintf("unexpected character %q", c), start+1)
		}
	}
	l.pos = n
	return Token{Type: TokEOF, Loc: Loc{Start: n, End: n}}
}

// scanBlockComment scans nested /* ... */ comments matching HiveLexerParent.g's
// recursive QUERY_HINT rule.
func scanBlockComment(input string, start int) (int, error) {
	depth := 0
	i := start
	n := len(input)
	for i < n {
		if i+1 < n && input[i] == '/' && input[i+1] == '*' {
			depth++
			i += 2
			continue
		}
		if i+1 < n && input[i] == '*' && input[i+1] == '/' {
			depth--
			i += 2
			if depth == 0 {
				return i, nil
			}
			continue
		}
		i++
	}
	return 0, &ParseError{Message: "unterminated block comment", Position: start}
}

func matchDollarKeyword(input string, start int) int {
	for _, kw := range []string{"$KEY$", "$VALUE$", "$ELEM$"} {
		if start+len(kw) <= len(input) && strings.EqualFold(input[start:start+len(kw)], kw) {
			return len(kw)
		}
	}
	return 0
}

// scanDigitStart resolves tokens that start with a digit using ANTLR's
// longest-match / first-rule priority:
// CharSetLiteral (0X...), IntegralLiteral, NumberLiteral, ByteLengthLiteral,
// Number, and Identifier (which may begin with a digit in Hive).
func scanDigitStart(input string, start int) (Token, int) {
	n := len(input)
	// Measure identifier span (Digit | Letter | '_')+
	idEnd := start
	for idEnd < n && (isLetter(input[idEnd]) || isDigit(input[idEnd]) || input[idEnd] == '_') {
		idEnd++
	}

	// Check CharSetLiteral: '0' ('x'|'X') HexDigit+
	hexEnd := start
	if input[start] == '0' && start+2 < n && (input[start+1] == 'x' || input[start+1] == 'X') && isHexDigit(input[start+2]) {
		hexEnd = start + 3
		for hexEnd < n && isHexDigit(input[hexEnd]) {
			hexEnd++
		}
	}

	// Measure base digits
	digEnd := start
	for digEnd < n && isDigit(input[digEnd]) {
		digEnd++
	}

	// Measure Number: Digit+ ( '.' Digit* Exponent? | Exponent )?
	numEnd := digEnd
	hasDotOrExp := false
	if numEnd < n && input[numEnd] == '.' {
		hasDotOrExp = true
		numEnd++
		for numEnd < n && isDigit(input[numEnd]) {
			numEnd++
		}
		if e, ok := scanExponent(input, numEnd); ok {
			numEnd = e
		}
	} else if e, ok := scanExponent(input, numEnd); ok {
		hasDotOrExp = true
		numEnd = e
	}

	// Check NumberLiteral: Number ('D' | 'BD') (case-insensitive)
	numLitEnd := 0
	if numEnd < n && (input[numEnd] == 'd' || input[numEnd] == 'D') {
		numLitEnd = numEnd + 1
	} else if numEnd+1 < n && (input[numEnd] == 'b' || input[numEnd] == 'B') && (input[numEnd+1] == 'd' || input[numEnd+1] == 'D') {
		numLitEnd = numEnd + 2
	}

	// Check IntegralLiteral: Digit+ ('L' | 'S' | 'Y') (case-insensitive)
	intLitEnd := 0
	if !hasDotOrExp && digEnd < n {
		switch input[digEnd] {
		case 'l', 'L', 's', 'S', 'y', 'Y':
			intLitEnd = digEnd + 1
		}
	}

	// Check ByteLengthLiteral: Digit+ ('B' | 'K' | 'M' | 'G') (case-insensitive)
	byteLitEnd := 0
	if !hasDotOrExp && digEnd < n {
		switch input[digEnd] {
		case 'b', 'B', 'k', 'K', 'm', 'M', 'g', 'G':
			byteLitEnd = digEnd + 1
		}
	}

	// Pick longest match; on tie, HiveLexerParent.g rule order wins:
	// CharSetLiteral > IntegralLiteral > NumberLiteral > ByteLengthLiteral > Number > Identifier.
	bestType := TokNumber
	bestEnd := numEnd
	tryCandidate := func(tt TokenType, end int, higherPriority bool) {
		if end > bestEnd || (end == bestEnd && higherPriority) {
			bestType = tt
			bestEnd = end
		}
	}
	if hexEnd > start {
		tryCandidate(TokCharSetLiteral, hexEnd, true)
	}
	if intLitEnd > start {
		tryCandidate(TokIntegralLiteral, intLitEnd, bestType != TokCharSetLiteral)
	}
	if numLitEnd > start {
		tryCandidate(TokNumberLiteral, numLitEnd, bestType != TokCharSetLiteral && bestType != TokIntegralLiteral)
	}
	if byteLitEnd > start {
		tryCandidate(TokByteLengthLiteral, byteLitEnd, bestType == TokNumber || bestType == TokIdentifier)
	}
	if idEnd > bestEnd {
		bestType = TokIdentifier
		bestEnd = idEnd
	}
	return Token{Type: bestType, Str: input[start:bestEnd], Loc: Loc{Start: start, End: bestEnd}}, bestEnd
}

func scanExponent(input string, pos int) (int, bool) {
	n := len(input)
	if pos >= n || (input[pos] != 'e' && input[pos] != 'E') {
		return pos, false
	}
	j := pos + 1
	if j < n && (input[j] == '+' || input[j] == '-') {
		j++
	}
	digStart := j
	for j < n && isDigit(input[j]) {
		j++
	}
	if j == digStart {
		return pos, false
	}
	return j, true
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isCharSetNameRest(c byte) bool {
	return isLetter(c) || isDigit(c) || c == '_' || c == '-' || c == '.' || c == ':'
}
