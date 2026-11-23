package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/parser"
)

// Statement is one statement of a script, as SplitStatements returns it.
type Statement = analysis.Statement

// SplitStatements splits a script into statements with parser.Split, keeping
// each statement's trailing ';' in its text and range. It returns nil when the
// script holds no statement.
//
// A segment holding no statement (a lone ';', a comment) is skipped, so every
// Statement it returns has Empty false: its callers (the driver, explain,
// ValidateQuery) execute or classify statements and never counted delimiters.
func SplitStatements(statement string) []Statement {
	var stmts []Statement
	for _, seg := range parser.Split(statement) {
		if seg.Empty() {
			continue
		}
		// parser.Split excludes the ';' from the segment and points ByteEnd at it.
		text := seg.Text
		end := seg.ByteEnd
		if end < len(statement) && statement[end] == ';' {
			text = statement[seg.ByteStart : end+1]
			end++
		}
		stmts = append(stmts, Statement{
			Text:  text,
			Start: seg.ByteStart,
			End:   end,
		})
	}
	return stmts
}
