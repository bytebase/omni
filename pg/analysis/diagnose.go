package analysis

import (
	"errors"
	"strings"
	"unicode"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/parser"
)

// DiagnoseError is the first syntax error Diagnose finds.
type DiagnoseError struct {
	// Text is the script as parsed: the input, with ";" appended when its last
	// statement lacked one.
	Text string
	// ParseError is set when the parser returned its *parser.ParseError; any
	// other error carries only Message.
	ParseError bool
	// Offset is the ParseError's byte offset in Text, -1 when the parser did
	// not locate it.
	Offset  int
	Message string
}

// Diagnose parses a script for an editor and returns its first syntax error,
// or nil. The last statement may omit its ";".
func Diagnose(statement string) *DiagnoseError {
	trimmedStatement := strings.TrimRightFunc(statement, unicode.IsSpace)
	if len(trimmedStatement) > 0 && !strings.HasSuffix(trimmedStatement, ";") {
		// Add a semicolon to the end of the statement to allow users to omit the semicolon
		// for the last statement in the script.
		statement += ";"
	}

	_, err := omnipg.Parse(statement)
	if err == nil {
		return nil
	}

	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return &DiagnoseError{Text: statement, Offset: -1, Message: err.Error()}
	}
	return &DiagnoseError{Text: statement, ParseError: true, Offset: parseErr.Position, Message: parseErr.Message}
}
