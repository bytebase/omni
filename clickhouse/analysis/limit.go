package analysis

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/bytebase/omni/clickhouse/parser"
)

// StatementWithResultLimit bounds a subqueryable ClickHouse statement with a
// CTE wrapper (`WITH result AS (...) SELECT * FROM result LIMIT N;`).
// Statements that cannot stand as a subquery (SHOW, DESCRIBE, INSERT, EXPLAIN,
// WITH ... INSERT, queries with top-level FORMAT or INTO OUTFILE, etc.) are
// returned unchanged.
func StatementWithResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	if !isSubqueryable(wrapped) {
		return statement
	}
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", wrapped, limit)
}

// WrapResultLimit trims whitespace, trailing semicolons, and trailing comments
// from statement and unconditionally wraps it in ClickHouse's result-limit CTE.
func WrapResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", wrapped, limit)
}

// trimStatementForWrapper strips leading whitespace and trailing whitespace,
// semicolons, and comments (line, `#`, and nested block comments) using
// clickhouse/parser's lexer so the statement can sit cleanly inside a CTE.
func trimStatementForWrapper(statement string) string {
	trimmed := strings.TrimLeftFunc(statement, unicode.IsSpace)
	l := parser.NewLexer(trimmed)
	lastEnd := 0
	for {
		tok := l.NextToken()
		if tok.Type == parser.TokEOF {
			break
		}
		if tok.Type == parser.TokSemicolon {
			continue
		}
		lastEnd = tok.End
	}
	return trimmed[:lastEnd]
}

// isSubqueryable reports whether statement can stand inside a ClickHouse CTE.
// It requires a single parsed SELECT/UNION query without top-level FORMAT or
// INTO OUTFILE clauses, which ClickHouse does not permit inside subqueries.
func isSubqueryable(statement string) bool {
	if statement == "" {
		return false
	}
	stmts, err := parser.Parse(statement)
	if err != nil || len(stmts) != 1 {
		return false
	}
	sel, ok := stmts[0].(*parser.SelectUnionStmt)
	if !ok {
		return false
	}
	if sel.Output != nil && (sel.Output.Format != "" || sel.Output.Outfile != "") {
		return false
	}
	return true
}
