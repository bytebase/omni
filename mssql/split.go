package mssql

import (
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// Segment represents a single statement segment from splitting a T-SQL script.
// Byte ranges are contiguous ([ByteStart, ByteEnd)), Start/End are 1-based line
// and code-point column positions, and Empty marks GO batch separators, lone
// semicolons, and whitespace/comment-only segments.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
	Start     Position
	End       Position
	Empty     bool
}

// Split splits a T-SQL script into statement segments.
//
// When Parse succeeds, Split derives statement boundaries from the parsed AST so
// semicolon-less T-SQL statements and compound bodies (such as CREATE PROCEDURE
// with internal semicolons) stay intact while GO batch separators, lone
// semicolons, and trailing comment-only segments become Empty segments. When
// Parse fails on incomplete or invalid SQL, Split falls back to
// SplitBySemicolon.
//
// Returned segments maintain byte contiguity (seg[0].ByteStart == 0,
// seg[i].ByteStart == seg[i-1].ByteEnd) and position adjacency
// (seg[0].Start == {1, 1}, seg[i].Start == seg[i-1].End).
func Split(sql string) []Segment {
	if len(sql) == 0 {
		return nil
	}
	stmts, err := Parse(sql)
	if err != nil {
		return SplitBySemicolon(sql)
	}

	idx := review.Index(sql)
	result := make([]Segment, 0, len(stmts))
	prevByte := 0
	prevEnd := Position{Line: 1, Column: 1}

	emit := func(endByte int, empty bool) {
		endPos := positionAt(idx, endByte)
		result = append(result, Segment{
			Text:      sql[prevByte:endByte],
			ByteStart: prevByte,
			ByteEnd:   endByte,
			Start:     prevEnd,
			End:       endPos,
			Empty:     empty,
		})
		prevByte = endByte
		prevEnd = endPos
	}

	// emitPreGapSemicolons emits an Empty segment for each top-level semicolon
	// in sql[prevByte:limit] (skipping comments) so parsed and fallback paths
	// agree on separator counts.
	emitPreGapSemicolons := func(limit int) {
		i := prevByte
		for i < limit {
			switch sql[i] {
			case '-':
				if isTSQLLineComment(sql, i, limit) {
					i, _ = skipTSQLLineComment(sql, i, limit)
				} else {
					i++
				}
			case '/':
				if i+1 < limit && sql[i+1] == '*' {
					end, _ := skipTSQLBlockComment(sql[:limit], i)
					i = end
				} else {
					i++
				}
			case ';':
				i++
				emit(i, true)
			default:
				i++
			}
		}
	}

	for _, os := range stmts {
		endByte := os.ByteEnd
		if endByte < prevByte {
			endByte = prevByte
		}
		if endByte > len(sql) {
			endByte = len(sql)
		}

		contentStart := nodeLoc(os.AST).Start
		if contentStart > prevByte && contentStart <= endByte {
			emitPreGapSemicolons(contentStart)
		}

		text := sql[prevByte:endByte]
		_, isGo := os.AST.(*ast.GoStmt)
		empty := isGo || os.Empty() || isOnlyTSQLSpace(text)
		emit(endByte, empty)
	}

	// Process any unparsed tail through the lexical splitter so lone semicolons
	// and comments become Empty segments while malformed trailing tokens (such
	// as unterminated strings or block comments) stay non-empty.
	if prevByte < len(sql) {
		atLineStart := prevByte == 0 || sql[prevByte-1] == '\n' || sql[prevByte-1] == '\r'
		result = append(result, splitBySemicolonFrom(sql, idx, prevByte, prevEnd, atLineStart)...)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// SplitBySemicolon splits a T-SQL script lexically at top-level semicolons and
// standalone GO batch separator lines without parsing. It skips single-quoted
// strings (with ” escapes), double-quoted identifiers ("" escapes), bracketed
// identifiers (]] escapes), line comments (--, honoring SQL Server's --(*
// exception), and nested block comments (/* ... */). It never fails;
// unterminated strings, identifiers, or block comments stay non-empty so
// callers that parse each segment surface the exact lexical error.
func SplitBySemicolon(sql string) []Segment {
	if len(sql) == 0 {
		return nil
	}
	return splitBySemicolonFrom(sql, review.Index(sql), 0, Position{Line: 1, Column: 1}, true)
}

func splitBySemicolonFrom(sql string, idx *review.Text, start int, prevEnd Position, atLineStart bool) []Segment {
	var result []Segment
	hasContent := false
	lastContentEnd := start

	emit := func(end int, empty bool) {
		endPos := positionAt(idx, end)
		result = append(result, Segment{
			Text:      sql[start:end],
			ByteStart: start,
			ByteEnd:   end,
			Start:     prevEnd,
			End:       endPos,
			Empty:     empty,
		})
		start = end
		prevEnd = endPos
		hasContent = false
		lastContentEnd = end
	}

	i := start
	for i < len(sql) {
		ch := sql[i]
		if ch == '\n' || ch == '\r' {
			i++
			atLineStart = true
			continue
		}
		if size, ok := matchTSQLInlineSpace(sql, i); ok {
			// Horizontal whitespace preserves GO eligibility at line start.
			i += size
			continue
		}
		switch ch {
		case '-':
			if isTSQLLineComment(sql, i, len(sql)) {
				var endedLine bool
				i, endedLine = skipTSQLLineComment(sql, i, len(sql))
				if endedLine {
					atLineStart = true
				}
			} else if strings.HasPrefix(sql[i:], "--(*") {
				// SQL Server rejects --(* as an illegal ODBC initiator rather
				// than a line comment; advance past all four bytes as content.
				i += 4
				hasContent = true
				lastContentEnd = i
				atLineStart = false
			} else {
				i++
				hasContent = true
				lastContentEnd = i
				atLineStart = false
			}
		case '/':
			if i+1 < len(sql) && sql[i+1] == '*' {
				commentStart := i
				end, closed := skipTSQLBlockComment(sql, i)
				i = end
				if !closed {
					// Unterminated block comment is malformed SQL, not empty trivia.
					hasContent = true
					lastContentEnd = i
					atLineStart = false
				} else if strings.ContainsAny(sql[commentStart:end], "\r\n") {
					// A closed multi-line block comment leaves the scanner at
					// line-leading trivia on its closing line, preserving GO eligibility.
					atLineStart = true
				}
				// A single-line closed block comment preserves atLineStart without
				// rescanning leading comments on each iteration.
			} else {
				i++
				hasContent = true
				lastContentEnd = i
				atLineStart = false
			}
		case '\'':
			i = skipTSQLDelimited(sql, i, '\'')
			hasContent = true
			lastContentEnd = i
			atLineStart = false
		case '"':
			i = skipTSQLDelimited(sql, i, '"')
			hasContent = true
			lastContentEnd = i
			atLineStart = false
		case '[':
			i = skipTSQLDelimited(sql, i, ']')
			hasContent = true
			lastContentEnd = i
			atLineStart = false
		case ';':
			i++
			emit(i, !hasContent)
			atLineStart = false
		case 'g', 'G':
			if atLineStart {
				atLineStart = false
				if goEnd, ok := matchGoAt(sql, i); ok {
					if hasContent {
						emit(lastContentEnd, false)
					}
					emit(goEnd, true)
					i = goEnd
					continue
				}
			}
			i++
			hasContent = true
			lastContentEnd = i
		default:
			i++
			hasContent = true
			lastContentEnd = i
			atLineStart = false
		}
	}

	if start < len(sql) && !isOnlyTSQLSpace(sql[start:]) {
		emit(len(sql), !hasContent)
	}
	return result
}

// isTSQLLineComment reports whether sql[i:limit] opens a valid T-SQL line
// comment. SQL Server rejects "--(*" (OdbcInitiator) instead of treating it as
// a comment.
func isTSQLLineComment(sql string, i, limit int) bool {
	return i+1 < limit && sql[i] == '-' && sql[i+1] == '-' && !strings.HasPrefix(sql[i:limit], "--(*")
}

// skipTSQLLineComment advances past a T-SQL "--" line comment through its line
// terminator (\r\n, \n, or \r) and reports whether a line break was consumed.
func skipTSQLLineComment(sql string, i, limit int) (int, bool) {
	i += 2
	for i < limit && sql[i] != '\n' && sql[i] != '\r' {
		i++
	}
	if i >= limit {
		return i, false
	}
	if sql[i] == '\r' {
		i++
		if i < limit && sql[i] == '\n' {
			i++
		}
		return i, true
	}
	i++ // '\n'
	return i, true
}

// matchGoAt checks whether sql[pos:] (already positioned at 'g'/'G' after
// line-leading whitespace or closed block comments) is a standalone T-SQL GO
// batch separator: case-insensitive "GO", optional positive integer repeat
// count, and only horizontal whitespace, closed block comments, or a line
// comment until the end of the line or input. It returns the byte offset
// immediately after GO (or its repeat count).
func matchGoAt(sql string, pos int) (int, bool) {
	if pos+2 > len(sql) || (sql[pos] != 'g' && sql[pos] != 'G') || (sql[pos+1] != 'o' && sql[pos+1] != 'O') {
		return 0, false
	}
	j := pos + 2
	if !isGoBoundary(sql, j) {
		return 0, false
	}
	goEnd := j

	j, sawNewline, ok := skipGoInlineTrivia(sql, j)
	if !ok {
		return 0, false
	}
	if sawNewline {
		return goEnd, true
	}
	if j > goEnd && j < len(sql) && sql[j] >= '0' && sql[j] <= '9' {
		numStart := j
		nonZero := false
		for j < len(sql) && sql[j] >= '0' && sql[j] <= '9' {
			if sql[j] != '0' {
				nonZero = true
			}
			j++
		}
		if !nonZero || !isGoBoundary(sql, j) {
			return 0, false
		}
		if j > numStart {
			goEnd = j
		}
		j, sawNewline, ok = skipGoInlineTrivia(sql, j)
		if !ok {
			return 0, false
		}
		if sawNewline {
			return goEnd, true
		}
	}

	if j == len(sql) || sql[j] == '\r' || sql[j] == '\n' || isTSQLLineComment(sql, j, len(sql)) {
		return goEnd, true
	}
	return 0, false
}

// skipGoInlineTrivia skips horizontal whitespace and closed block comments on
// a GO line. If a closed block comment contains a line break (\r or \n), the
// GO line has ended and sawNewline is returned true so tokens on the comment's
// closing line are not validated against the GO line.
func skipGoInlineTrivia(sql string, j int) (pos int, sawNewline bool, ok bool) {
	for j < len(sql) {
		if size, isSpace := matchTSQLInlineSpace(sql, j); isSpace {
			j += size
			continue
		}
		if sql[j] == '/' && j+1 < len(sql) && sql[j+1] == '*' {
			end, closed := skipTSQLBlockComment(sql, j)
			if !closed {
				return 0, false, false
			}
			if strings.ContainsAny(sql[j:end], "\r\n") {
				return end, true, true
			}
			j = end
			continue
		}
		break
	}
	return j, false, true
}

func isGoBoundary(sql string, pos int) bool {
	if pos >= len(sql) {
		return true
	}
	if sql[pos] == '\r' || sql[pos] == '\n' {
		return true
	}
	if _, ok := matchTSQLInlineSpace(sql, pos); ok {
		return true
	}
	switch sql[pos] {
	case '-':
		return isTSQLLineComment(sql, pos, len(sql))
	case '/':
		return pos+1 < len(sql) && sql[pos+1] == '*'
	default:
		return false
	}
}

// skipTSQLBlockComment skips a T-SQL block comment starting at "/*", honoring
// nested /* ... */ pairs. It returns the next byte offset and whether the
// comment was properly terminated.
func skipTSQLBlockComment(sql string, i int) (int, bool) {
	i += 2
	depth := 1
	for i < len(sql) && depth > 0 {
		if sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*' {
			depth++
			i += 2
		} else if sql[i] == '*' && i+1 < len(sql) && sql[i+1] == '/' {
			depth--
			i += 2
		} else {
			i++
		}
	}
	return i, depth == 0
}

// skipTSQLDelimited skips a T-SQL string literal ('...'), double-quoted
// identifier ("..."), or bracketed identifier ([...]), where doubling close
// escapes it.
func skipTSQLDelimited(sql string, i int, close byte) int {
	i++ // skip opener
	for i < len(sql) {
		if sql[i] == close {
			i++
			if i < len(sql) && sql[i] == close {
				i++
				continue
			}
			return i
		}
		i++
	}
	return i
}

// matchTSQLInlineSpace checks whether sql[i:] starts with a SQL Server
// non-newline whitespace character (WS_CHAR_WO_NEWLINE in TSql160.g) and
// returns its UTF-8 byte length.
func matchTSQLInlineSpace(sql string, i int) (int, bool) {
	ch := sql[i]
	if ch < 0x80 {
		if ch <= 0x20 && ch != '\n' && ch != '\r' {
			return 1, true
		}
		return 0, false
	}
	r, size := utf8.DecodeRuneInString(sql[i:])
	if isTSQLInlineSpaceRune(r) {
		return size, true
	}
	return 0, false
}

func isTSQLInlineSpaceRune(r rune) bool {
	if r <= 0x20 {
		return r != '\n' && r != '\r'
	}
	switch {
	case r == 0x0085, r == 0x00a0, r == 0x1680,
		r >= 0x2000 && r <= 0x200b,
		r == 0x2028, r == 0x2029, r == 0x202f, r == 0x205f, r == 0x3000:
		return true
	default:
		return false
	}
}

func isOnlyTSQLSpace(s string) bool {
	for i := 0; i < len(s); {
		if s[i] == '\n' || s[i] == '\r' {
			i++
			continue
		}
		size, ok := matchTSQLInlineSpace(s, i)
		if !ok {
			return false
		}
		i += size
	}
	return true
}
