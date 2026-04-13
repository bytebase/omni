// Package parser implements a hand-written Doris SQL lexer and
// recursive-descent parser. F4 ships only the entry-point framework;
// concrete statement parsing is added by Tier 1+ DAG nodes.
package parser

import "github.com/bytebase/omni/doris/ast"

// ParseError describes a single parse error with its source location.
//
// Loc uses the same shape as LexError so consumers can handle both
// uniformly.
type ParseError struct {
	Loc ast.Loc
	Msg string
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	return e.Msg
}
