package pg

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/yacc"
)

// Statement is the result of parsing a single SQL statement.
type Statement struct {
	// Text is the SQL text including trailing semicolon if present.
	Text string
	// AST is the inner statement node (e.g. *ast.SelectStmt). Nil for empty statements.
	AST ast.Node

	// ByteStart is the inclusive start byte offset in the original SQL.
	ByteStart int
	// ByteEnd is the exclusive end byte offset in the original SQL.
	ByteEnd int

	// Start is the start position (line:column) in the original SQL.
	Start Position
	// End is the exclusive end position (line:column) in the original SQL.
	End Position
}

// Position represents a location in source text.
type Position struct {
	// Line is 1-based line number.
	Line int
	// Column is 1-based column in bytes.
	Column int
}

// Empty returns true if this statement has no meaningful content.
func (s *Statement) Empty() bool {
	return s.AST == nil
}

// Parse splits and parses a SQL string into statements.
// Each statement includes the text, AST, and byte/line positions.
func Parse(sql string) ([]Statement, error) {
	list, err := yacc.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return nil, nil
	}

	// Build a line index for position calculation.
	lineIndex := buildLineIndex(sql)

	// The grammar returns statement nodes directly (not wrapped in RawStmt).
	// We split the SQL text by semicolons to determine each statement's boundaries.
	ranges := splitRanges(sql, len(list.Items))

	var stmts []Statement
	for i, item := range list.Items {
		if item == nil {
			continue
		}

		// Unwrap RawStmt if present (future-proofing).
		node := item
		if raw, ok := item.(*ast.RawStmt); ok {
			node = raw.Stmt
		}

		start := ranges[i][0]
		end := ranges[i][1]
		text := sql[start:end]

		// Start position points to the first non-whitespace character.
		contentStart := start
		for contentStart < end && isSpace(sql[contentStart]) {
			contentStart++
		}

		stmts = append(stmts, Statement{
			Text:      text,
			AST:       node,
			ByteStart: start,
			ByteEnd:   end,
			Start:     offsetToPosition(lineIndex, contentStart),
			End:       offsetToPosition(lineIndex, end),
		})
	}
	return stmts, nil
}

// splitRanges splits SQL text into n statement ranges by semicolons.
// Returns [n][2]int where each pair is [start, end) byte offsets.
// It handles quoted strings, comments, and dollar-quoted strings.
func splitRanges(sql string, n int) [][2]int {
	ranges := make([][2]int, 0, n)
	start := 0
	i := 0

	for i < len(sql) && len(ranges) < n-1 {
		switch sql[i] {
		case '\'':
			// Skip single-quoted string.
			i++
			for i < len(sql) {
				if sql[i] == '\'' {
					i++
					if i < len(sql) && sql[i] == '\'' {
						i++ // escaped quote
						continue
					}
					break
				}
				i++
			}
		case '"':
			// Skip double-quoted identifier.
			i++
			for i < len(sql) {
				if sql[i] == '"' {
					i++
					if i < len(sql) && sql[i] == '"' {
						i++ // escaped quote
						continue
					}
					break
				}
				i++
			}
		case '-':
			// Skip line comment.
			if i+1 < len(sql) && sql[i+1] == '-' {
				for i < len(sql) && sql[i] != '\n' {
					i++
				}
			} else {
				i++
			}
		case '/':
			// Skip block comment.
			if i+1 < len(sql) && sql[i+1] == '*' {
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
			} else {
				i++
			}
		case '$':
			// Skip dollar-quoted string.
			tag := parseDollarTag(sql, i)
			if tag != "" {
				i += len(tag)
				end := findDollarEnd(sql, i, tag)
				if end >= 0 {
					i = end + len(tag)
				}
			} else {
				i++
			}
		case ';':
			// Found statement boundary.
			end := i + 1 // include the semicolon
			ranges = append(ranges, [2]int{start, end})
			start = end
			i++
		default:
			i++
		}
	}

	// Last statement gets the rest.
	if len(ranges) < n {
		end := len(sql)
		ranges = append(ranges, [2]int{start, end})
	}

	return ranges
}

func parseDollarTag(sql string, pos int) string {
	if pos >= len(sql) || sql[pos] != '$' {
		return ""
	}
	end := pos + 1
	for end < len(sql) {
		if sql[end] == '$' {
			return sql[pos : end+1]
		}
		c := sql[end]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9' && end > pos+1)) {
			break
		}
		end++
	}
	// $$ (empty tag)
	if end == pos+1 && end < len(sql) && sql[end] == '$' {
		return "$$"
	}
	return ""
}

func findDollarEnd(sql string, pos int, tag string) int {
	for i := pos; i <= len(sql)-len(tag); i++ {
		if sql[i] == '$' && sql[i:i+len(tag)] == tag {
			return i
		}
	}
	return -1
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// lineIndex stores the byte offset of each line start.
// lineIndex[0] = 0 (line 1 starts at byte 0).
type lineIndex []int

func buildLineIndex(s string) lineIndex {
	idx := lineIndex{0}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			idx = append(idx, i+1)
		}
	}
	return idx
}

func offsetToPosition(idx lineIndex, offset int) Position {
	// Binary search for the line containing offset.
	lo, hi := 0, len(idx)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if idx[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return Position{
		Line:   lo + 1,               // 1-based
		Column: offset - idx[lo] + 1, // 1-based
	}
}
