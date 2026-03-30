// Package cosmosdb provides a parser for Azure Cosmos DB NoSQL SQL API queries.
package cosmosdb

import (
	"sort"

	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/parser"
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
	Column int // 1-based, bytes
}

func (s Statement) Empty() bool {
	return s.AST == nil
}

func Parse(sql string) ([]Statement, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list.Len() == 0 {
		return nil, nil
	}
	idx := buildLineIndex(sql)
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
			Start:     offsetToPosition(idx, byteStart),
			End:       offsetToPosition(idx, byteEnd),
		})
	}
	return stmts, nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

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
	line := sort.SearchInts(idx, offset+1)
	col := offset - idx[line-1] + 1
	return Position{Line: line, Column: col}
}
