package compat

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// TrinoParsingException is the class of the exception Trino's parser
// throws, as an error's failureInfo names it: SqlParser's for text the
// grammar does not derive, and AstBuilder's for what the grammar derives
// and the AST refuses, such as an identifier that starts with a digit.
const TrinoParsingException = "io.trino.sql.parser.ParsingException"

// trinoSyntaxError is the errorName of every ParsingException.
const trinoSyntaxError = "SYNTAX_ERROR"

// A TrinoError is what a Trino server's answer to a failed statement says,
// as its client protocol's QueryError carries it.
type TrinoError struct {
	// Name is errorName, such as SYNTAX_ERROR or TABLE_NOT_FOUND.
	Name string
	// Type is failureInfo.type, the class of the exception.
	Type string
	// Line and Column are errorLocation: from 1, Column in code points, as
	// ANTLR counts them. 0 is no location.
	Line, Column int
}

// TrinoVerdict reads a Trino server's answer to a statement: whether its
// parser accepts it. err is what running the statement returned, and read
// returns what the server said, or false for an error the server did not
// send (a dropped connection, a timeout, an HTTP status).
//
// The parser's refusal is a ParsingException, at the byte of its location.
// Its errorName, SYNTAX_ERROR, is not enough: the analyzer and the metadata
// layer raise SYNTAX_ERROR too, after the parse, as for a table name with
// too many dots, and those statements parsed. So does every other error,
// and a statement that ran.
//
// An answer that is no verdict comes back as an error, so the caller fails
// rather than records agreement: an error the server did not send, one with
// no exception class, a ParsingException under another name, and a
// location outside the statement.
func TrinoVerdict(statement string, err error, read func(error) (TrinoError, bool)) (string, error) {
	if err == nil {
		return Accept, nil
	}
	e, ok := read(err)
	if !ok {
		return "", err
	}
	if e.Type == "" {
		return "", fmt.Errorf("Trino answered %s with no exception class: %w", e.Name, err)
	}
	if e.Type != TrinoParsingException {
		return Accept, nil
	}
	if e.Name != trinoSyntaxError {
		return "", fmt.Errorf("Trino's parser answered %s, not %s: %w", e.Name, trinoSyntaxError, err)
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
