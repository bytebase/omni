package pg

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
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
// Text boundaries are derived from RawStmt.Loc: each statement's text
// spans from the end of the previous statement to just past the semicolon
// (or end of input) following the current statement.
func Parse(sql string) ([]Statement, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return nil, nil
	}

	lineIndex := buildLineIndex(sql)

	var stmts []Statement
	prevEnd := 0 // byte offset where previous statement's text ended

	for _, item := range list.Items {
		if item == nil {
			continue
		}

		raw, ok := item.(*ast.RawStmt)
		if !ok {
			continue
		}

		// Text starts where the previous statement ended (includes leading whitespace/comments).
		start := prevEnd

		// Text ends after the semicolon following the statement, or at the Loc.End if no semicolon.
		end := raw.Loc.End
		// Scan past trailing whitespace to find the semicolon.
		j := end
		for j < len(sql) && isSpace(sql[j]) {
			j++
		}
		if j < len(sql) && sql[j] == ';' {
			end = j + 1
		}

		// Start position points to the first non-whitespace character.
		contentStart := start
		for contentStart < end && isSpace(sql[contentStart]) {
			contentStart++
		}

		stmts = append(stmts, Statement{
			Text:      sql[start:end],
			AST:       raw.Stmt,
			ByteStart: start,
			ByteEnd:   end,
			Start:     offsetToPosition(lineIndex, contentStart),
			End:       offsetToPosition(lineIndex, end),
		})

		prevEnd = end
	}
	return stmts, nil
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
