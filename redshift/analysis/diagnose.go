package analysis

import (
	"errors"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/parser"
)

// SyntaxError is the syntax diagnostic Diagnose reports.
type SyntaxError struct {
	// Offset is the error's byte offset in the script.
	Offset int
	// Message is the parser's message.
	Message string
}

// Diagnose returns the script's syntax error, or nil when it parses.
//
// It parses the script whole, as the editor shows it: an error in any statement
// is the error, and its offset is in the script's coordinates.
func Diagnose(script string) *SyntaxError {
	_, err := redshift.Parse(script)
	if err == nil {
		return nil
	}
	var parseErr *parser.ParseError
	if errors.As(err, &parseErr) {
		return &SyntaxError{Offset: parseErr.Position, Message: parseErr.Message}
	}
	return &SyntaxError{Message: err.Error()}
}
