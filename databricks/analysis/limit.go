package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytebase/omni/databricks/parser"
)

// StatementWithResultLimit bounds a Databricks SQL query to return at most
// limit rows using databricks/parser's AST offsets.
//
// For a single parsed query (SELECT, WITH ... SELECT, VALUES, TABLE, UNION,
// FROM-first pipeline):
//   - An existing outermost integer LIMIT <= limit is preserved unchanged.
//   - An existing outermost integer LIMIT > limit or LIMIT ALL is lowered to limit.
//   - A non-integer outermost LIMIT expression (e.g. LIMIT 10 + 20) is wrapped
//     in a subquery so the cap still holds.
//   - An outermost OFFSET without LIMIT gets `LIMIT <limit> ` spliced before
//     OFFSET (Spark/Databricks requires LIMIT before OFFSET).
//   - Otherwise ` LIMIT <limit>` is spliced at the end of the query span
//     (before any trailing semicolon or comment).
//
// Non-query statements (SHOW, DESCRIBE, EXPLAIN, DDL, DML, compound blocks),
// external script transforms, non-positive limits, multi-statement scripts,
// and unparseable statements are returned unchanged.
func StatementWithResultLimit(statement string, limit int) string {
	if limit <= 0 {
		return statement
	}
	res, err := parser.ParseDatabricks(statement)
	if err != nil || res == nil || len(res.Statements) != 1 {
		return statement
	}
	ps, ok := res.Statements[0].(*parser.ParsedStmt)
	if !ok || ps == nil || ps.Kind != parser.StmtQuery || ps.HasTransform || ps.HasEmptyPipe {
		return statement
	}

	// Existing outermost LIMIT clause.
	if ps.LimitLoc.End > ps.LimitLoc.Start && ps.LimitLoc.End <= len(statement) {
		if ps.LimitAll {
			return statement[:ps.LimitLoc.Start] + strconv.Itoa(limit) + statement[ps.LimitLoc.End:]
		}
		raw := strings.TrimSpace(statement[ps.LimitLoc.Start:ps.LimitLoc.End])
		if existing, err := strconv.Atoi(raw); err == nil && existing >= 0 {
			if existing <= limit {
				return statement
			}
			return statement[:ps.LimitLoc.Start] + strconv.Itoa(limit) + statement[ps.LimitLoc.End:]
		}
		return WrapResultLimit(statement, limit)
	}

	// Outermost OFFSET without LIMIT: LIMIT must precede OFFSET in Spark SQL.
	if ps.OffsetLoc.End > ps.OffsetLoc.Start && ps.OffsetLoc.Start <= len(statement) {
		return statement[:ps.OffsetLoc.Start] + fmt.Sprintf("LIMIT %d ", limit) + statement[ps.OffsetLoc.Start:]
	}

	end := ps.Loc.End
	if end <= 0 || end > len(statement) {
		return statement
	}
	if ps.HasPipe {
		candidate := statement[:end] + fmt.Sprintf(" |> LIMIT %d", limit) + statement[end:]
		if res2, err := parser.ParseDatabricks(candidate); err == nil && res2 != nil && len(res2.Statements) == 1 {
			return candidate
		}
		return WrapResultLimit(statement, limit)
	}
	candidate := statement[:end] + fmt.Sprintf(" LIMIT %d", limit) + statement[end:]
	// Confirm the appended LIMIT binds to the outermost query organization rather
	// than being absorbed by a trailing FROM-first branch (e.g. `SELECT * FROM t
	// UNION ALL FROM u SELECT *`), falling back to a wrapper subquery otherwise.
	if res2, err := parser.ParseDatabricks(candidate); err == nil && res2 != nil && len(res2.Statements) == 1 {
		if ps2, ok := res2.Statements[0].(*parser.ParsedStmt); ok && ps2 != nil && ps2.LimitLoc.End > ps2.LimitLoc.Start {
			return candidate
		}
	}
	return WrapResultLimit(statement, limit)
}

// WrapResultLimit trims leading whitespace and trailing whitespace, semicolons,
// and comments using databricks/parser's lexer and wraps statement in a
// Databricks subquery capped at limit rows.
func WrapResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	return fmt.Sprintf("SELECT * FROM (%s) AS result LIMIT %d;", wrapped, limit)
}

func trimStatementForWrapper(statement string) string {
	trimmed := strings.TrimLeftFunc(statement, unicode.IsSpace)
	l := parser.NewLexer(trimmed)
	lastEnd := 0
	for {
		tok := l.NextToken()
		if l.Err() != nil || tok.Kind == parser.TokEOF {
			break
		}
		if tok.Kind == parser.TokSemicolon {
			continue
		}
		lastEnd = tok.End
	}
	return trimmed[:lastEnd]
}
