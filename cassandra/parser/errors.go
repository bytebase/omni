package parser

import (
	"fmt"

	"github.com/bytebase/omni/cassandra/ast"
)

// ParseError represents a syntax error during CQL parsing.
type ParseError struct {
	Message string
	Loc     ast.Loc
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Loc.Start, e.Message)
}

func locFromOffsets(start, end int) ast.Loc {
	return ast.Loc{Start: start, End: end}
}
