package mysqlcomplete

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func isIdentByte(b byte) bool {
	return b == '_' || b == '$' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func isIdentRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$'
}

func nextIdentRune(s string, i int) (int, bool) {
	if s[i] < utf8.RuneSelf {
		return 1, isIdentByte(s[i])
	}
	r, size := utf8.DecodeRuneInString(s[i:])
	return size, r != utf8.RuneError && isIdentRune(r)
}

func backtickIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func isObjectIdentifierCandidate(t CandidateType) bool {
	switch t {
	case CandidateDatabase,
		CandidateTable,
		CandidateView,
		CandidateColumn,
		CandidateIndex,
		CandidateTrigger,
		CandidateEvent,
		CandidateProcedure:
		return true
	default:
		return false
	}
}

// QuoteIdentifierIfNeeded backtick-quotes an object identifier when it would
// otherwise be invalid as a bare identifier, so that accepting the completion
// inserts valid SQL. When the caret already sits inside a backtick-quoted
// identifier the user is typing, the name is returned unquoted.
func QuoteIdentifierIfNeeded(name string, caretInBacktick bool, isReservedKeyword func(string) bool) string {
	if caretInBacktick || !identifierNeedsQuoting(name, isReservedKeyword) {
		return name
	}
	return backtickIdentifier(name)
}

func identifierNeedsQuoting(name string, isReservedKeyword func(string) bool) bool {
	if name == "" {
		return true
	}
	for _, r := range name {
		if !isIdentRune(r) {
			return true
		}
	}
	firstRune, _ := utf8.DecodeRuneInString(name)
	if unicode.IsDigit(firstRune) {
		return true
	}
	return isReservedKeyword != nil && isReservedKeyword(name)
}

// CaretAtStatementStart reports whether pos in statement sits at the start of a
// statement (or after a semicolon or completed WITH CTE list).
func CaretAtStatementStart(statement string, pos int) bool {
	pos = max(0, min(pos, len(statement)))
	before := StripStringsAndComments(statement)[:pos]
	before = strings.TrimRightFunc(before, isIdentRune)
	before = strings.TrimRight(before, " \t\r\n")
	if before == "" || before[len(before)-1] == ';' {
		return true
	}
	return caretAfterWithClause(before)
}

func caretAfterWithClause(before string) bool {
	toks := scanQualTokens(before)
	if len(toks) == 0 || toks[len(toks)-1].punct != ')' {
		return false
	}
	depth := 0
	startIdx := 0
	withIdx := -1
	withDepth := -1
	for i, tok := range toks {
		switch tok.punct {
		case ';':
			startIdx = i + 1
			depth = 0
			withIdx = -1
			withDepth = -1
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
			if depth < withDepth {
				withIdx = -1
				withDepth = -1
			}
		default:
			if tok.upper == "WITH" && (i == startIdx || (i > 0 && toks[i-1].punct == '(')) {
				withIdx = i
				withDepth = depth
			}
		}
	}
	if withIdx < 0 || depth != withDepth {
		return false
	}
	d := withDepth
	sawCTEBody := false
	for i := withIdx + 1; i < len(toks); i++ {
		tok := toks[i]
		switch tok.punct {
		case '(':
			if d == withDepth && i > 0 && toks[i-1].upper == "AS" {
				sawCTEBody = true
			}
			d++
		case ')':
			if d > withDepth {
				d--
			}
		case ',':
			if d == withDepth {
				sawCTEBody = false
			}
		default:
			if d == withDepth && sawCTEBody && tok.isWord {
				return false
			}
		}
	}
	return sawCTEBody
}

// CaretInsideBacktickIdentifier reports whether the caret sits inside an open
// backtick-quoted identifier (an odd number of backticks precede it).
func CaretInsideBacktickIdentifier(statement string, pos int) bool {
	pos = max(0, min(pos, len(statement)))
	masked := StripStringsAndComments(statement)
	count := 0
	for i := 0; i < pos; i++ {
		if masked[i] == '`' {
			count++
		}
	}
	return count%2 == 1
}

// StripStringsAndComments returns statement with the contents of string literals
// ('...', "..."), line comments (-- , #), and block comments (/* */) replaced by
// spaces, preserving byte length (newlines kept) so caret offsets stay valid.
// Backtick-quoted identifiers are left intact.
func StripStringsAndComments(statement string) string {
	out := []byte(statement)
	n := len(out)
	blank := func(i int) {
		if out[i] != '\n' {
			out[i] = ' '
		}
	}
	i := 0
	for i < n {
		switch c := out[i]; {
		case c == '\'' || c == '"':
			q := c
			blank(i)
			i++
			for i < n {
				if out[i] == '\\' && i+1 < n {
					blank(i)
					blank(i + 1)
					i += 2
					continue
				}
				if out[i] == q {
					if i+1 < n && out[i+1] == q {
						blank(i)
						blank(i + 1)
						i += 2
						continue
					}
					blank(i)
					i++
					break
				}
				blank(i)
				i++
			}
		case c == '#':
			for i < n && out[i] != '\n' {
				out[i] = ' '
				i++
			}
		case c == '-' && i+1 < n && out[i+1] == '-' && (i+2 >= n || isSpaceByte(out[i+2])):
			for i < n && out[i] != '\n' {
				out[i] = ' '
				i++
			}
		case c == '/' && i+1 < n && out[i+1] == '*':
			blank(i)
			blank(i + 1)
			i += 2
			for i < n {
				if out[i] == '*' && i+1 < n && out[i+1] == '/' {
					blank(i)
					blank(i + 1)
					i += 2
					break
				}
				blank(i)
				i++
			}
		case c == '`':
			i++
			for i < n {
				if out[i] == '`' {
					if i+1 < n && out[i+1] == '`' {
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
		default:
			i++
		}
	}
	return string(out)
}

func isSpaceByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
