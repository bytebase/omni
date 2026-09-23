package oracle

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/review"
)

// Statement is the result of parsing a single Oracle SQL statement or PL/SQL
// unit from a script.
type Statement struct {
	Text string
	AST  ast.Node

	ByteStart int
	ByteEnd   int

	Start Position
	End   Position
}

// Position represents a location in source text.
type Position struct {
	Line   int // 1-based
	Column int // 1-based, code points
}

// Empty returns true if this statement has no meaningful content.
func (s *Statement) Empty() bool {
	return s.AST == nil
}

// Parse splits and parses an Oracle SQL script into statements.
func Parse(sql string) ([]Statement, error) {
	segments := parser.Split(sql)
	if len(segments) == 0 {
		return nil, nil
	}

	lineIndex := review.Index(sql)
	stmts := make([]Statement, 0, len(segments))
	for _, seg := range segments {
		list, err := parser.ParseRange(sql, seg.ByteStart, seg.ByteEnd)
		if err != nil {
			return nil, err
		}

		var node ast.Node
		if list != nil {
			for _, item := range list.Items {
				raw, ok := item.(*ast.RawStmt)
				if ok {
					node = raw.Stmt
					break
				}
			}
		}

		contentStart := seg.ByteStart
		for contentStart < seg.ByteEnd && isSpace(sql[contentStart]) {
			contentStart++
		}

		stmts = append(stmts, Statement{
			Text:      seg.Text,
			AST:       node,
			ByteStart: seg.ByteStart,
			ByteEnd:   seg.ByteEnd,
			Start:     positionAt(lineIndex, contentStart),
			End:       positionAt(lineIndex, seg.ByteEnd),
		})
	}
	return stmts, nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

// positionAt converts a byte offset into a Position through the shared index:
// lines are 1-based, columns are 1-based and counted in code points, the
// units Bytebase's Position uses.
func positionAt(idx *review.Text, offset int) Position {
	line, column := idx.Position(offset)
	return Position{Line: line, Column: column}
}
