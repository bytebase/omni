// Package analysis provides parser-backed SQL analysis for ClickHouse:
// statement splitting, EXPLAIN rewriting and inspection, result-limit
// wrapping, and read-only query validation.
package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
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
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil, nil
	}
	statements := make([]Statement, 0, len(segments))
	for _, seg := range segments {
		statements = append(statements, Statement{
			Text:  seg.Text,
			Empty: seg.Empty(),
			Start: seg.ByteStart,
			End:   seg.ByteEnd,
		})
	}
	return statements, nil
}
