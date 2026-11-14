package analysis

import (
	"errors"

	"github.com/bytebase/omni/tidb/parser"
)

// SyntaxError is the first syntax error Diagnose finds, in the statement
// that holds it.
type SyntaxError struct {
	// ByteStart and ByteEnd are the statement's range in the script, as
	// parser.Split reports it, without the ";" that may follow.
	ByteStart int
	ByteEnd   int
	// Offset is the error's byte offset in the statement's text, which runs
	// through its ";" when one follows.
	Offset      int
	Message     string
	RelatedText string
}

// Diagnose splits a script for an editor and parses each statement,
// returning the first syntax error, or nil. An error that is not a syntax
// error is returned as is.
func Diagnose(statement string) (*SyntaxError, error) {
	for _, seg := range parser.Split(statement) {
		if seg.Empty() {
			continue
		}
		end := seg.ByteEnd
		if end < len(statement) && statement[end] == ';' {
			end++
		}
		_, err := parser.Parse(statement[seg.ByteStart:end])
		if err == nil {
			continue
		}
		var parseErr *parser.ParseError
		if !errors.As(err, &parseErr) {
			return nil, err
		}
		return &SyntaxError{
			ByteStart:   seg.ByteStart,
			ByteEnd:     seg.ByteEnd,
			Offset:      parseErr.Position,
			Message:     parseErr.Message,
			RelatedText: parseErr.RelatedText,
		}, nil
	}
	return nil, nil
}
