// Package cosmosdb provides a parser for Azure Cosmos DB NoSQL SQL API queries.
package cosmosdb

import (
	"strings"

	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/parser"
	"github.com/bytebase/omni/review"
)

type Statement struct {
	Text      string
	AST       ast.Node
	ByteStart int
	ByteEnd   int
	Start     Position
	End       Position
}

type Position struct {
	Line   int // 1-based
	Column int // 1-based, code points
}

func (s Statement) Empty() bool {
	return s.AST == nil
}

func Parse(sql string) ([]Statement, error) {
	// Return empty result for blank input rather than a parse error.
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
