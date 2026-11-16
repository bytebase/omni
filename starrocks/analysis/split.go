package analysis

import (
	"github.com/bytebase/omni/starrocks/parser"
)

// Statement is one segment of a script. Text keeps the trailing ';' the engine
// splits on, because that is the text a caller logs and executes.
type Statement struct {
	// Text is the segment's source text, including its trailing ';' when one
	// follows it.
	Text string
	// Empty reports whether the segment holds no statement (whitespace,
	// comments, or a lone delimiter).
	Empty bool
	// Start is the segment's first byte in the script.
	Start int
	// End is the byte just past the segment, including a trailing ';'.
	End int
}

// SplitStatements splits a script into the statements a caller executes. The
// split is lexical and never fails; a segment with no statement is reported
// with Empty set.
func SplitStatements(script string) []Statement {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]Statement, 0, len(segments))
	for _, segment := range segments {
		text := segment.Text
		end := segment.ByteEnd
		if end < len(script) && script[end] == ';' {
			text = script[segment.ByteStart : end+1]
			end++
		}
		statements = append(statements, Statement{
			Text:  text,
			Empty: segment.Empty(),
			Start: segment.ByteStart,
			End:   end,
		})
	}
	return statements
}
