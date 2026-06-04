package parser

import "github.com/bytebase/omni/trino/ast"

// ParseError describes a single parse error with its source location.
//
// Loc uses the same {Start,End} byte-offset shape as LexError (tokens.go) so
// consumers can handle lex and parse failures uniformly. Line/column
// conversion is a caller-side concern; bytebase's Diagnose maps Loc.Start back
// to a (line, col) pair against the original source.
//
// ParseError is a pure data carrier: Error returns just the message. This
// mirrors snowflake/parser.ParseError and doris/parser.ParseError so the omni
// parser family stays uniform across engines.
type ParseError struct {
	Loc ast.Loc
	Msg string
}

// Error implements the error interface, returning the message only. Callers
// that want a formatted "msg (line N, col M)" string convert Loc.Start with a
// line table.
func (e *ParseError) Error() string {
	return e.Msg
}
