package mssql

import (
	"strings"

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

	// scanGap emits an Empty segment for each top-level semicolon in
	// sql[prevByte:limit] (skipping comments) and reports whether an
	// unterminated block comment was encountered.
	scanGap := func(limit int) bool {
		unterminatedComment := false
		i := prevByte
		for i < limit {
			switch sql[i] {
			case '-':
				if i+1 < limit && sql[i+1] == '-' {
					i += 2
					for i < limit && sql[i] != '\n' {
						i++
					}
					if i < limit && sql[i] == '\n' {
						i++
					}
				} else {
					i++
				}
			case '/':
				if i+1 < limit && sql[i+1] == '*' {
					end, closed := skipTSQLBlockComment(sql[:limit], i)
					i = end
					if !closed {
						unterminatedComment = true
					}
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
		return unterminatedComment
	}

	for _, os := range stmts {
		endByte := os.ByteEnd
		if endByte < prevByte {
			endByte = prevByte
		}
		if endByte > len(sql) {
			endByte = len(sql)
		}

		// Preserve any lone semicolons before the AST node's start token as
		// Empty segments so parsed and fallback paths agree on separator counts.
		contentStart := nodeLoc(os.AST).Start
		if contentStart > prevByte && contentStart <= endByte {
			scanGap(contentStart)
		}

		text := sql[prevByte:endByte]
		_, isGo := os.AST.(*ast.GoStmt)
		empty := isGo || os.Empty() || strings.TrimSpace(text) == ""
		emit(endByte, empty)
	}

	if prevByte < len(sql) {
		unterminatedComment := scanGap(len(sql))
		if prevByte < len(sql) && !isOnlyTSQLSpace(sql[prevByte:]) {
			emit(len(sql), !unterminatedComment)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// SplitBySemicolon splits a T-SQL script lexically at top-level semicolons and
// standalone GO batch separator lines without parsing. It skips single-quoted
// strings (with ” escapes), double-quoted identifiers ("" escapes), bracketed
// identifiers (]] escapes), line comments (--), and nested block comments
// (/* ... */). It never fails; unterminated strings, identifiers, or block
// comments stay non-empty so callers that parse each segment surface the exact
// lexical error.
func SplitBySemicolon(sql string) []Segment {
	if len(sql) == 0 {
		return nil
	}

	idx := review.Index(sql)
	var result []Segment
	start := 0
	prevEnd := Position{Line: 1, Column: 1}
	hasContent := false
	lastContentEnd := 0
	atLineStart := true

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

	i := 0
	for i < len(sql) {
		if atLineStart {
			atLineStart = false
			if goEnd, ok := matchStandaloneGoLine(sql, i); ok {
				if hasContent {
					emit(lastContentEnd, false)
				}
				emit(goEnd, true)
				i = goEnd
				continue
			}
		}

		ch := sql[i]
		switch ch {
		case '\n':
			i++
			atLineStart = true
		case ' ', '\t', '\r', '\f', '\v':
			i++
		case '-':
			if i+1 < len(sql) && sql[i+1] == '-' {
				i += 2
				for i < len(sql) && sql[i] != '\n' {
					i++
				}
				if i < len(sql) && sql[i] == '\n' {
					i++
					atLineStart = true
				}
			} else {
				i++
				hasContent = true
				lastContentEnd = i
			}
		case '/':
			if i+1 < len(sql) && sql[i+1] == '*' {
				end, closed := skipTSQLBlockComment(sql, i)
				i = end
				if !closed {
					// Unterminated block comment is malformed SQL, not empty trivia.
					hasContent = true
					lastContentEnd = i
				}
			} else {
				i++
				hasContent = true
				lastContentEnd = i
			}
		case '\'':
			i = skipTSQLDelimited(sql, i, '\'')
			hasContent = true
			lastContentEnd = i
		case '"':
			i = skipTSQLDelimited(sql, i, '"')
			hasContent = true
			lastContentEnd = i
		case '[':
			i = skipTSQLDelimited(sql, i, ']')
			hasContent = true
			lastContentEnd = i
		case ';':
			i++
			emit(i, !hasContent)
		default:
			i++
			hasContent = true
			lastContentEnd = i
		}
	}

	if start < len(sql) && !isOnlyTSQLSpace(sql[start:]) {
		emit(len(sql), !hasContent)
	}
	return result
}

// matchStandaloneGoLine checks whether sql[lineStart:] is a standalone T-SQL GO
// batch separator line: optional horizontal whitespace, case-insensitive "GO",
// optional positive integer repeat count, and only horizontal whitespace or a
// line comment until the end of the line or input. It returns the byte offset
// immediately after GO (or its repeat count).
func matchStandaloneGoLine(sql string, lineStart int) (int, bool) {
	j := lineStart
	for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
		j++
	}
	if j+2 > len(sql) || (sql[j] != 'g' && sql[j] != 'G') || (sql[j+1] != 'o' && sql[j+1] != 'O') {
		return 0, false
	}
	j += 2
	if !isGoBoundary(sql, j) {
		return 0, false
	}
	goEnd := j

	// Optional repeat count after whitespace.
	wsStart := j
	for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
		j++
	}
	if j > wsStart && j < len(sql) && sql[j] >= '0' && sql[j] <= '9' {
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
		for j < len(sql) && (sql[j] == ' ' || sql[j] == '\t') {
			j++
		}
	}

	if j == len(sql) || sql[j] == '\r' || sql[j] == '\n' || (sql[j] == '-' && j+1 < len(sql) && sql[j+1] == '-') {
		return goEnd, true
	}
	return 0, false
}

func isGoBoundary(sql string, pos int) bool {
	if pos >= len(sql) {
		return true
	}
	switch sql[pos] {
	case ' ', '\t', '\r', '\n':
		return true
	case '-':
		return pos+1 < len(sql) && sql[pos+1] == '-'
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

func isOnlyTSQLSpace(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
		default:
			return false
		}
	}
	return true
}
