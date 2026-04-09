// Package parser implements a hand-written Snowflake SQL lexer and (in
// future DAG nodes) a recursive-descent parser. F2 ships only the lexer.
package parser

import "github.com/bytebase/omni/snowflake/ast"

// LexError describes a single lexing failure with its source location.
//
// Lex errors are non-fatal: the lexer collects them in a slice (Lexer.Errors)
// and emits a synthetic tokInvalid token at each failure site so consumers
// can choose to halt on the first error or proceed with best-effort parsing.
type LexError struct {
	Loc ast.Loc
	Msg string
}

// Standard lex error messages. Tests assert against these constants so a
// reword in one place propagates everywhere.
const (
	errUnterminatedString  = "unterminated string literal"
	errUnterminatedDollar  = "unterminated $$ string literal"
	errUnterminatedComment = "unterminated block comment"
	errUnterminatedQuoted  = "unterminated quoted identifier"
	errInvalidByte         = "invalid byte"
)

// ParseError describes a single parse error with its source location.
//
// Loc uses the same shape as LexError so consumers can handle both
// uniformly. Line/column conversion is a caller-side concern via
// LineTable (defined in linetable.go).
type ParseError struct {
	Loc ast.Loc
	Msg string
}

// Error implements the error interface. Returns just the message — line
// and column are omitted here to keep ParseError a pure data carrier.
// Callers that want formatted "msg (line N, col M)" output should use a
// LineTable to convert Loc.Start into a (line, col) pair.
func (e *ParseError) Error() string {
	return e.Msg
}
