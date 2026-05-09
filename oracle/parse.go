package oracle

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
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
	Line   int
	Column int
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

	lineIndex := buildLineIndex(sql)
	stmts := make([]Statement, 0, len(segments))
	for _, seg := range segments {
		list, err := parser.Parse(strings.Repeat(" ", seg.ByteStart) + seg.Text)
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
			Start:     offsetToPosition(lineIndex, contentStart),
			End:       offsetToPosition(lineIndex, seg.ByteEnd),
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
	if offset < 0 {
		offset = 0
	}
	line := 0
	for line+1 < len(idx) && idx[line+1] <= offset {
		line++
	}
	return Position{
		Line:   line + 1,
		Column: offset - idx[line] + 1,
	}
}
