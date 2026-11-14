package analysis

import (
	"errors"
	"strings"
	"unicode"

	"github.com/bytebase/omni/mysql/parser"
)

// SyntaxError is the first syntax error Diagnose finds.
type SyntaxError struct {
	// Text is the script as parsed: the input, with ";" appended when its last
	// statement lacked one.
	Text string
	// Offset is the error's byte offset in Text, or -1 when the parser
	// reported no position.
	Offset  int
	Message string
}

// Diagnose parses a script for an editor and returns its first syntax error,
// or nil. The last statement may omit its ";".
func Diagnose(statement string) *SyntaxError {
	trimmed := strings.TrimRightFunc(statement, unicode.IsSpace)
	if len(trimmed) > 0 && !strings.HasSuffix(trimmed, ";") {
		statement += ";"
	}

	_, err := parser.Parse(statement)
	if err == nil {
		return nil
	}

	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return &SyntaxError{Text: statement, Offset: -1, Message: err.Error()}
	}
	return &SyntaxError{Text: statement, Offset: parseErr.Position, Message: parseErr.Message}
}
