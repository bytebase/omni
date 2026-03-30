// Package mongo provides a parser for MongoDB shell (mongosh) commands.
package mongo

import (
	"github.com/bytebase/omni/mongo/ast"
	"github.com/bytebase/omni/mongo/parser"
)

// Statement is the result of parsing a single mongosh command.
type Statement struct {
	// Text is the command text including trailing semicolon if present.
	Text string
	// AST is the parsed statement node. Nil for empty statements.
	AST ast.Node

	// ByteStart is the inclusive start byte offset in the original input.
	ByteStart int
	// ByteEnd is the exclusive end byte offset in the original input.
	ByteEnd int

	// Start is the start position (line:column) in the original input.
	Start Position
	// End is the exclusive end position (line:column) in the original input.
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

// Parse splits and parses a mongosh input string into statements.
// Each statement includes the text, AST, and byte/line positions.
func Parse(input string) ([]Statement, error) {
	nodes, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}

	lineIndex := buildLineIndex(input)

	var stmts []Statement
	for _, node := range nodes {
		loc := node.GetLoc()

		// Determine text boundaries — expand to include trailing semicolons/whitespace
		start := loc.Start
		end := loc.End
		// Scan past trailing whitespace to find semicolon
		j := end
		for j < len(input) && isSpace(input[j]) {
			j++
		}
		if j < len(input) && input[j] == ';' {
			end = j + 1
		}

		stmts = append(stmts, Statement{
			Text:      input[start:end],
			AST:       node,
			ByteStart: start,
			ByteEnd:   end,
			Start:     offsetToPosition(lineIndex, start),
			End:       offsetToPosition(lineIndex, end),
		})
	}
	return stmts, nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// lineIndex stores the byte offset of each line start.
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
		Line:   lo + 1,              // 1-based
		Column: offset - idx[lo] + 1, // 1-based
	}
}
