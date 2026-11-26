// Package analysis provides parser-backed SQL analysis for ClickHouse:
// statement splitting, EXPLAIN rewriting and inspection, result-limit
// wrapping, and read-only query validation.
package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/analysis"
)

// Statement is one segment of a split SQL script.
type Statement = analysis.Statement

// Split splits script into ClickHouse SQL segments at top-level semicolons
// using clickhouse/parser's lexer. Each segment's Text includes leading
// whitespace and comments since the previous segment, the terminating
// semicolon when present, and any same-line trailing line comment after the
// semicolon, so byte offsets stay contiguous with script[Start:End]. Empty
// segments (lone semicolons, whitespace, or comment-only chunks) are retained
// with Empty set to true.
func Split(script string) ([]Statement, error) {
	return Engine.Split(script), nil
}

// TrimStatement drops surrounding whitespace and trailing semicolons from
// statement.
func TrimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}
