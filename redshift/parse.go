package redshift

import (
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
	"github.com/bytebase/omni/review"
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
	// Column is 1-based column in code points.
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

	lineIndex := review.Index(sql)

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
			Start:     positionAt(lineIndex, contentStart),
			End:       positionAt(lineIndex, end),
		})

		prevEnd = end
	}
	return stmts, nil
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
