package parser

import (
	"fmt"
	"sort"

	"github.com/bytebase/omni/cassandra/ast"
)

// ParseError represents a syntax error during CQL parsing.
type ParseError struct {
	Message string
	Loc     ast.Loc
	Line    int
	Column  int
	Near    string
}

func (e *ParseError) Error() string {
	if e.Near != "" {
		return fmt.Sprintf("line %d column %d: %s at or near %q", e.Line, e.Column, e.Message, e.Near)
	}
	return fmt.Sprintf("line %d column %d: %s", e.Line, e.Column, e.Message)
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

func offsetToLineCol(idx lineIndex, offset int) (int, int) {
	line := sort.SearchInts(idx, offset+1)
	col := offset - idx[line-1] + 1
	return line, col
}

func locFromOffsets(start, end int) ast.Loc {
	return ast.Loc{Start: start, End: end}
}
