package analysis

import (
	"github.com/bytebase/omni/doris/parser"
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
// split is lexical and never fails; a segment with no statement — a lone ';',
// whitespace, a comment — is reported with Empty set, so a caller that
// accounts for every delimiter sees one entry per ';'.
func SplitStatements(script string) []Statement {
	segments := parser.SplitAll(script)
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

// firstLexError returns the first lexical error the whole of text carries, as
// a ParseError, or nil.
//
// parser.Split suppresses the errors of the chunks it drops, and a chunk that
// holds no statement is dropped, so a buffer such as "/*" (an unterminated
// block comment) would otherwise reach a caller as "no statement" instead of
// invalid SQL. ValidateQuery and ResolveQuerySpan check the whole text before
// splitting, so they fail closed on it.
func firstLexError(text string) error {
	l := parser.NewLexer(text)
	for l.NextToken().Kind != tokEOF {
	}
	errs := l.Errors()
	if len(errs) == 0 {
		return nil
	}
	return &parser.ParseError{Position: errs[0].Loc.Start, End: errs[0].Loc.End, Message: errs[0].Msg}
}
