package analysis

import (
	"strings"

	"github.com/bytebase/omni/doris/parser"
)

// SyntaxError is one diagnostic Diagnose reports.
type SyntaxError struct {
	// Offset is the error's byte offset in the script.
	Offset int
	// End is the byte offset just past the offending text, or -1.
	End int
	// Message is the parser's message.
	Message string
}

// Diagnose returns the syntax diagnostics of a script. The parser's
// "statement parsing is not yet supported" dispatch stubs are filtered out:
// they report a parser gap, not an error in the SQL.
func Diagnose(script string) []SyntaxError {
	diagnostics := parser.Diagnose(script)
	out := make([]SyntaxError, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		if strings.HasSuffix(diagnostic.Msg, "statement parsing is not yet supported") {
			continue
		}
		out = append(out, SyntaxError{Offset: diagnostic.Loc.Start, End: diagnostic.Loc.End, Message: diagnostic.Msg})
	}
	return out
}
