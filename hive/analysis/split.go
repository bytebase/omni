// Package analysis provides parser-backed SQL analysis for Hive: statement
// splitting, EXPLAIN rewriting and plan format detection, and read-only query
// classification.
package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/hive/parser"
)

// Statement is one segment of a split SQL script.
type Statement = analysis.Statement

// Split splits script into non-empty HiveQL statements using the Hive 4.0.1
// lexer. Each segment's Text includes leading whitespace and comments since the
// previous segment and the terminating semicolon when present, so byte offsets
// stay contiguous with script[Start:End].
func Split(script string) ([]Statement, error) {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil, nil
	}
	statements := make([]Statement, 0, len(segments))
	for _, seg := range segments {
		statements = append(statements, Statement{
			Text:  seg.Text,
			Empty: false,
			Start: seg.ByteStart,
			End:   seg.ByteEnd,
		})
	}
	return statements, nil
}

// TrimStatement drops surrounding whitespace and trailing semicolons from
// statement so a split segment's terminator does not reach the server.
func TrimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
