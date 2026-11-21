package analysis

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/oracle/parser"
)

// Statement is one statement of a script, as SplitStatements returns it.
type Statement = analysis.Statement

// SplitStatements splits a script into statements with parser.Split. SQL*Plus
// commands are dropped, a MATCH_RECOGNIZE whose DEFINE line parser.Split took
// for the SQL*Plus DEFINE command is rejoined, and trailing whitespace and
// comments are trimmed from each statement except embedded Java or JavaScript
// source.
func SplitStatements(statement string) []Statement {
	return splitStatements(statement, true)
}

// SplitStatementsForCompletion is SplitStatements without the trailing trim,
// so a statement keeps the text up to the cursor.
func SplitStatementsForCompletion(statement string) []Statement {
	return splitStatements(statement, false)
}

func splitStatements(statement string, trim bool) []Statement {
	segments := parser.Split(statement)

	result := make([]Statement, 0, len(segments))
	for i := 0; i < len(segments); i++ {
		seg := segments[i]
		if seg.Kind == parser.SegmentSQLPlusCommand {
			continue
		}
		byteStart := seg.ByteStart
		byteEnd := seg.ByteEnd
		if end, ok := matchRecognizeDefineMergeEnd(statement, segments, i); ok {
			byteEnd = end
			i += 2
		}
		// Embedded Java or JavaScript source does not lex as SQL: its "--"
		// decrement would read as a comment and cut the source short.
		if trim && seg.Kind != parser.SegmentEmbeddedSource {
			byteEnd = byteStart + trimSegmentTrailingHidden(statement[byteStart:byteEnd])
		}
		result = append(result, Statement{
			Text:  statement[byteStart:byteEnd],
			Empty: seg.Empty(),
			Start: byteStart,
			End:   byteEnd,
		})
	}
	return result
}

func matchRecognizeDefineMergeEnd(statement string, segments []parser.Segment, index int) (int, bool) {
	if index+2 >= len(segments) {
		return 0, false
	}
	current := segments[index]
	define := segments[index+1]
	suffix := segments[index+2]
	if current.Kind != parser.SegmentSQL || define.Kind != parser.SegmentSQLPlusCommand || suffix.Kind != parser.SegmentSQL {
		return 0, false
	}
	if !strings.Contains(strings.ToUpper(current.Text), "MATCH_RECOGNIZE") {
		return 0, false
	}
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(define.Text)), "DEFINE ") {
		return 0, false
	}
	merged := statement[current.ByteStart:suffix.ByteEnd]
	if _, err := parser.Parse(merged); err != nil {
		return 0, false
	}
	return suffix.ByteEnd, true
}

func trimSegmentTrailingHidden(text string) int {
	lexer := parser.NewLexer(text)
	lastTokenEnd := 0
	for {
		token := lexer.NextToken()
		if token.Type == 0 {
			break
		}
		lastTokenEnd = token.End
	}
	if lexer.Err != nil || lastTokenEnd == 0 {
		return len(text)
	}
	return lastTokenEnd
}
