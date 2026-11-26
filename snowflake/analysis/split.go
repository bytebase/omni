package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/parser"
)

// Statement is one statement of a script, as SplitStatements returns it.
type Statement = analysis.Statement

// SplitStatements splits a script into statements with parser.Split, keeping
// each statement's trailing ';' in its text and range. It returns nil for an
// empty script.
//
// A segment holding no statement (a lone ';', a comment) is returned with Empty
// set, so a caller that accounts for every delimiter sees one entry per ';'; a
// caller that executes or classifies statements skips it.
func SplitStatements(statement string) []Statement {
	var stmts []Statement
	for _, seg := range parser.Split(statement) {
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
