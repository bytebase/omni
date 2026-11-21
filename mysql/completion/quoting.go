package completion

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/mysql/parser"
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
// identifier the user is typing, outer delimiters are omitted while embedded
// backticks are still doubled so they cannot prematurely close the identifier.
func QuoteIdentifierIfNeeded(name string, caretInBacktick bool) string {
	return quoteIdentifierIfNeeded(name, caretInBacktick, isReservedKeyword)
}

func quoteIdentifierIfNeeded(name string, caretInBacktick bool, isReserved func(string) bool) string {
	if caretInBacktick {
		return strings.ReplaceAll(name, "`", "``")
	}
	if !identifierNeedsQuoting(name, isReserved) {
		return name
	}
	return backtickIdentifier(name)
}

func identifierNeedsQuoting(name string, isReserved func(string) bool) bool {
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
	return isReserved != nil && isReserved(name)
}

// extractIdentTokenAt extracts the identifier prefix being typed at
// cursorOffset, the byte offset where that token begins (before any opening
// backtick when caret is inside a backtick-quoted identifier), and whether the
// caret sits inside an open backtick identifier.
func extractIdentTokenAt(sql string, cursorOffset int) (prefix string, tokenStart int, inBacktick bool) {
	end := max(0, min(cursorOffset, len(sql)))
	if CaretInsideBacktickIdentifier(sql, end) {
		masked := stripStringsAndComments(sql[:end])
		for i := end - 1; i >= 0; i-- {
			if masked[i] == '`' {
				if i > 0 && masked[i-1] == '`' {
					i--
					continue
				}
				raw := sql[i+1 : end]
				return strings.ReplaceAll(raw, "``", "`"), i, true
			}
		}
	}
	i := end
	for i > 0 {
		if c := sql[i-1]; c < utf8.RuneSelf {
			if !isIdentByte(c) {
				break
			}
			i--
			continue
		}
		r, size := utf8.DecodeLastRuneInString(sql[:i])
		if r == utf8.RuneError || !isIdentRune(r) {
			break
		}
		i -= size
	}
	return sql[i:end], i, false
}

// cleanSQLForBacktickCollect returns a copy of sql where the open backtick
// token spanning [tokenStart, cursorOffset) (up to its matching closing backtick
// if present, else cursorOffset) is replaced by spaces so lexer-based C3
// collection at tokenStart never stalls on an unterminated backtick token.
func cleanSQLForBacktickCollect(sql string, tokenStart, cursorOffset int) string {
	n := len(sql)
	tokenStart = max(0, min(tokenStart, n))
	end := max(tokenStart, min(cursorOffset, n))
	if tokenStart < n && sql[tokenStart] == '`' {
		j := tokenStart + 1
		foundClose := false
		for j < n {
			if sql[j] == '`' {
				if j+1 < n && sql[j+1] == '`' {
					j += 2
					continue
				}
				j++
				foundClose = true
				break
			}
			j++
		}
		if foundClose && j >= end {
			end = j
		}
	}
	return sql[:tokenStart] + strings.Repeat(" ", end-tokenStart) + sql[end:]
}

// identBeforeDot scans left from pos expecting (optional whitespace/comments) '.'
// (optional whitespace/comments) <identifier>, returning the identifier text,
// the byte index where it starts, and whether the pattern matched. The
// identifier may be bare or backtick-quoted (with doubled “ unescaped).
func identBeforeDot(sql string, pos int) (ident string, start int, ok bool) {
	masked := stripStringsAndComments(sql)
	pos = max(0, min(pos, len(masked)))
	i := pos
	for i > 0 && isSpaceByte(masked[i-1]) {
		i--
	}
	if i == 0 || masked[i-1] != '.' {
		return "", pos, false
	}
	i-- // consume '.'
	for i > 0 && isSpaceByte(masked[i-1]) {
		i--
	}
	if i == 0 {
		return "", pos, false
	}
	if masked[i-1] == '`' {
		end := i - 1
		j := end - 1
		for j >= 0 {
			if masked[j] == '`' {
				if j > 0 && masked[j-1] == '`' {
					j -= 2
					continue
				}
				return strings.ReplaceAll(sql[j+1:end], "``", "`"), j, true
			}
			j--
		}
		return "", pos, false
	}
	end := i
	for i > 0 {
		if c := masked[i-1]; c < utf8.RuneSelf {
			if !isIdentByte(c) {
				break
			}
			i--
			continue
		}
		r, size := utf8.DecodeLastRuneInString(masked[:i])
		if r == utf8.RuneError || !isIdentRune(r) {
			break
		}
		i -= size
	}
	if i == end {
		return "", pos, false
	}
	return sql[i:end], i, true
}

func caretAtStatementStart(statement string, pos int) bool {
	pos = max(0, min(pos, len(statement)))
	before := stripStringsAndComments(statement)[:pos]
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
	type withFrame struct {
		idx   int
		depth int
	}
	var withStack []withFrame
	depth := 0
	startIdx := 0
	for i, tok := range toks {
		switch tok.punct {
		case ';':
			startIdx = i + 1
			depth = 0
			withStack = withStack[:0]
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
			for len(withStack) > 0 && withStack[len(withStack)-1].depth > depth {
				withStack = withStack[:len(withStack)-1]
			}
		default:
			if tok.upper == "WITH" && (i == startIdx || (i > 0 && toks[i-1].punct == '(')) {
				withStack = append(withStack, withFrame{idx: i, depth: depth})
			}
		}
	}
	if len(withStack) == 0 {
		return false
	}
	top := withStack[len(withStack)-1]
	if depth != top.depth {
		return false
	}
	withIdx, withDepth := top.idx, top.depth
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

// CaretInsideBacktickIdentifier reports whether byte offset pos in statement
// sits inside an open backtick-quoted identifier (an odd number of backticks
// precede it).
func CaretInsideBacktickIdentifier(statement string, pos int) bool {
	pos = max(0, min(pos, len(statement)))
	masked := stripStringsAndComments(statement)
	count := 0
	for i := 0; i < pos; i++ {
		if masked[i] == '`' {
			count++
		}
	}
	return count%2 == 1
}

// stripStringsAndComments blanks string literals and ordinary comments while
// preserving executable-comment SQL (/*! ... */, /*T! ... */), backtick
// identifiers, newlines, and original byte offsets using the engine lexer.
func stripStringsAndComments(statement string) string {
	out := make([]byte, len(statement))
	for i := 0; i < len(statement); i++ {
		if statement[i] == '\n' {
			out[i] = '\n'
		} else {
			out[i] = ' '
		}
	}
	lex := parser.NewLexer(statement)
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
		if tok.Loc < 0 || tok.End > len(statement) || tok.Loc >= tok.End {
			continue
		}
		if isStringLiteralSpan(statement[tok.Loc:tok.End]) {
			continue
		}
		copy(out[tok.Loc:tok.End], statement[tok.Loc:tok.End])
	}
	return string(out)
}

func isStringLiteralSpan(raw string) bool {
	if raw == "" {
		return false
	}
	if raw[0] == '\'' || raw[0] == '"' {
		return true
	}
	return len(raw) >= 2 && (raw[0] == 'b' || raw[0] == 'B' || raw[0] == 'x' || raw[0] == 'X') && raw[1] == '\''
}

func isSpaceByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
