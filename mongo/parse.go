// Package mongo provides a parser for MongoDB shell (mongosh) commands.
package mongo

import (
	"github.com/bytebase/omni/mongo/ast"
	"github.com/bytebase/omni/mongo/parser"
	"github.com/bytebase/omni/review"
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
	// Column is 1-based column in code points.
	Column int
}

// Empty returns true if this statement has no meaningful content.
func (s *Statement) Empty() bool {
	return s.AST == nil
}

// ParseResult holds the outcome of a best-effort parse.
type ParseResult struct {
	Statements []Statement
	Errors     []*parser.ParseError
}

// Parse splits and parses a mongosh input string into statements.
// Each statement includes the text, AST, and byte/line positions.
//
// Parse is strict: it returns a non-nil error (the first one encountered)
// when any statement in the input fails to parse, so callers never silently
// lose statements. Use ParseBestEffort to recover partial results together
// with all errors.
func Parse(input string) ([]Statement, error) {
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		return nil, result.Errors[0]
	}
	return result.Statements, nil
}

// ParseBestEffort splits and parses a mongosh input string, recovering from
// errors. It always returns as many successfully parsed statements as possible,
// along with any errors encountered.
func ParseBestEffort(input string) *ParseResult {
	pr := parser.ParseBestEffort(input)
	result := &ParseResult{Errors: pr.Errors}
	if len(pr.Nodes) == 0 {
		return result
	}

	lineIndex := review.Index(input)
	for _, node := range pr.Nodes {
		loc := node.GetLoc()
		start := loc.Start
		end := loc.End
		j := end
		for j < len(input) && isSpace(input[j]) {
			j++
		}
		if j < len(input) && input[j] == ';' {
			end = j + 1
		}
		result.Statements = append(result.Statements, Statement{
			Text:      input[start:end],
			AST:       node,
			ByteStart: start,
			ByteEnd:   end,
			Start:     positionAt(lineIndex, start),
			End:       positionAt(lineIndex, end),
		})
	}
	return result
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// positionAt converts a byte offset into a Position through the shared index:
// lines are 1-based, columns are 1-based and counted in code points, the
// units Bytebase's Position uses.
func positionAt(idx *review.Text, offset int) Position {
	line, column := idx.Position(offset)
	return Position{Line: line, Column: column}
}
