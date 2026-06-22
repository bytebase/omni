// Package cassandra provides a parser for Apache Cassandra CQL (Cassandra Query Language).
package cassandra

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/cassandra/parser"
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
	Column int // 1-based, bytes
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

	idx := buildLineIndex(sql)
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
