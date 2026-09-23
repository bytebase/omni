// Package cassandra provides a parser for Apache Cassandra CQL (Cassandra Query Language).
package cassandra

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/cassandra/parser"
	"github.com/bytebase/omni/review"
)

// Statement represents a single parsed CQL statement with position information.
type Statement struct {
	Text      string
	AST       ast.Node
	ByteStart int
	ByteEnd   int
	Start     Position
	End       Position
}

// Position represents a line/column location in source text.
type Position struct {
	Line   int // 1-based
	Column int // 1-based, code points
}

// Empty reports whether the statement is empty (no AST).
func (s Statement) Empty() bool {
	return s.AST == nil
}

// Parse parses a CQL input containing one or more statements.
func Parse(sql string) ([]Statement, error) {
	if strings.TrimSpace(sql) == "" {
		return nil, nil
	}

	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list.Len() == 0 {
		return nil, nil
	}

	idx := review.Index(sql)
	var stmts []Statement
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok {
			continue
		}
		byteStart := raw.StmtLocation
		byteEnd := byteStart + raw.StmtLen
		// Trim trailing whitespace from statement text.
		for byteEnd > byteStart && isSpace(sql[byteEnd-1]) {
			byteEnd--
		}
		stmts = append(stmts, Statement{
			Text:      sql[byteStart:byteEnd],
			AST:       raw.Stmt,
			ByteStart: byteStart,
			ByteEnd:   byteEnd,
			Start:     positionAt(idx, byteStart),
			End:       positionAt(idx, byteEnd),
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
