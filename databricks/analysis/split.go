// Package analysis provides parser-backed SQL analysis for Databricks SQL:
// statement splitting (including compound BEGIN ... END blocks), statement
// classification, read-only query validation, EXPLAIN rewriting and
// inspection, result-limit rewriting, and view dependency extraction / DDL
// qualification.
package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/analysis"
)

// Statement is one segment of a split SQL script.
type Statement = analysis.Statement

// Split splits script into Databricks SQL segments at top-level semicolons
// using databricks/parser's Spark 4.0.5 lexer, keeping compound BEGIN ... END
// scripting blocks intact. Each segment's Text includes leading whitespace and
// comments since the previous segment, the terminating semicolon when present,
// and any same-line trailing comment after the semicolon, so byte offsets stay
// contiguous with script[Start:End]. Empty segments (lone semicolons,
// whitespace, or comment-only chunks) are retained with Empty set to true.
func Split(script string) ([]Statement, error) {
	return Engine.Split(script), nil
}

// TrimStatement drops surrounding whitespace and trailing semicolons from
// statement so a split segment's terminator does not reach the server.
func TrimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
