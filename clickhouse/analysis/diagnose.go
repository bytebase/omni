package analysis

import (
	"errors"

	"github.com/bytebase/omni/clickhouse/parser"
)

// DiagnoseError is the first syntax error Diagnose finds.
type DiagnoseError struct {
	// Offset is the error's byte offset in the script, -1 when the parser
	// did not locate it.
	Offset  int
	Message string
}

// Diagnose parses a script for an editor with ClickHouse's grammar, one
// statement at a time, and returns its first syntax error, or nil.
func Diagnose(script string) *DiagnoseError {
	for _, seg := range parser.Split(script) {
		if seg.Empty() {
			continue
		}
		_, err := parser.Parse(seg.Text)
		if err == nil {
			continue
		}
		var pe *parser.ParseError
		if errors.As(err, &pe) {
			return &DiagnoseError{Offset: seg.ByteStart + pe.Position, Message: pe.Message}
		}
		return &DiagnoseError{Offset: -1, Message: err.Error()}
	}
	return nil
}
