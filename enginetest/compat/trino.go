package compat

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// trinoSyntaxError is the errorName of Trino's syntax errors.
const trinoSyntaxError = "SYNTAX_ERROR"

// A TrinoError is what a Trino server's answer to a failed statement says,
// as its client protocol's QueryError carries it.
type TrinoError struct {
	// Name is errorName, such as SYNTAX_ERROR or TABLE_NOT_FOUND.
	Name string
	// Line and Column are errorLocation: from 1, Column in code points, as
	// ANTLR counts them. 0 is no location.
	Line, Column int
}

// TrinoVerdict reads a Trino server's answer to a statement: whether Trino
// takes its syntax. err is what running the statement returned, and read
// returns what the server said, or false for an error the server did not
// send (a dropped connection, a timeout, an HTTP status).
//
// SYNTAX_ERROR is a rejection, at the byte of its location where it has one.
// The parser raises it, for text the grammar does not derive and for what
// its AST builder refuses (an identifier that starts with a digit). So do a
// few checks after the parse that read only the text, which omni's parser
// makes too: a name with too many parts, an INTERVAL qualifier whose fields
// are out of order, two isolation levels. Every other error comes after the
// parse and is an acceptance, as is a statement that ran.
//
// An answer that is no verdict comes back as an error, so the caller fails
// rather than records agreement: an error the server did not send, one
// with no errorName, and a location outside the statement.
func TrinoVerdict(statement string, err error, read func(error) (TrinoError, bool)) (string, error) {
	if err == nil {
		return Accept, nil
	}
	e, ok := read(err)
	if !ok {
		return "", err
	}
	if e.Name == "" {
		return "", fmt.Errorf("Trino answered with no errorName: %w", err)
	}
	if e.Name != trinoSyntaxError {
		return Accept, nil
	}
	if e.Line == 0 && e.Column == 0 {
		return Reject, nil
	}
	at, ok := trinoByte(statement, e.Line, e.Column)
	if !ok {
		return "", fmt.Errorf("Trino put the error at line %d, column %d, outside the statement: %w", e.Line, e.Column, err)
	}
	return RejectAt(at), nil
}

// trinoByte returns the byte offset of a line and column into statement,
// both from 1, the column in code points. ANTLR starts a line after each
// \n. The column just past the last character is the end of input.
func trinoByte(statement string, line, column int) (int, bool) {
	if line < 1 || column < 1 {
		return 0, false
	}
	at := 0
	for l := 1; l < line; l++ {
		nl := strings.IndexByte(statement[at:], '\n')
		if nl < 0 {
			return 0, false
		}
		at += nl + 1
	}
	for c := 1; c < column; c++ {
		if at >= len(statement) || statement[at] == '\n' {
			return 0, false
		}
		_, size := utf8.DecodeRuneInString(statement[at:])
		at += size
	}
	return at, true
}
