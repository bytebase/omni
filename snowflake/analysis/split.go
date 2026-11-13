package analysis

import "github.com/bytebase/omni/snowflake/parser"

// Statement is one statement of a script, as SplitStatements returns it.
type Statement struct {
	// Text is the statement's text, including its trailing ';' when present.
	Text string
	// Empty reports a statement with no SQL, only whitespace and comments.
	Empty bool
	// Start and End are the byte offsets of Text in the script.
	Start int
	End   int
}

// SplitStatements splits a script into statements with parser.Split, keeping
// each statement's trailing ';' in its text and range. It returns nil when the
// script holds no statement.
func SplitStatements(statement string) []Statement {
	segs := parser.Split(statement)
	if len(segs) == 0 {
		return nil
	}

	stmts := make([]Statement, 0, len(segs))
	for _, seg := range segs {
		// parser.Split excludes the ';' from the segment and points ByteEnd at it.
		text := seg.Text
		end := seg.ByteEnd
		if end < len(statement) && statement[end] == ';' {
			text = statement[seg.ByteStart : end+1]
			end++
		}
		stmts = append(stmts, Statement{
			Text:  text,
			Empty: seg.Empty(),
			Start: seg.ByteStart,
			End:   end,
		})
	}
	return stmts
}
