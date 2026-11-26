package analysis

import (
	"strconv"

	"github.com/bytebase/omni/spanner/ast"
)

// RewriteResultLimit sets the LIMIT of the outermost query in statement to at
// most limitCount, lowering a larger literal LIMIT and adding one where the
// query has none. It never wraps the statement in a subquery, where Spanner
// rejects a WITH clause, so a query it cannot rewrite — a non-literal LIMIT, a
// shape the parser does not read, a FOR UPDATE that a read-only transaction
// refuses anyway — is returned unchanged.

// rewritePipeResultLimit caps a pipe query at limitCount rows by rewriting its
// terminal `|> LIMIT` operator, or appending `|> LIMIT n` when the pipeline has
// no terminal limit. A terminal non-literal limit (`|> LIMIT @n`) is left alone
// unless force is set (BigQuery's wrapper used to enforce a literal outer cap),
// as a non-literal LIMIT is for a plain query; `|> LIMIT ALL` is capped by
// appending an operator, which the grammar and the engine allow.
func rewritePipeResultLimit(statement string, pipe *ast.Subpipeline, limitCount int, force bool) string {
	if n := len(pipe.Ops); n > 0 {
		if lim, ok := pipe.Ops[n-1].(*ast.PipeLimitOffset); ok && lim.Limit != nil {
			if count, ok := lim.Limit.(*ast.Literal); ok && count.Kind == ast.LitInt {
				if count.Ival <= int64(limitCount) {
					return statement
				}
				return statement[:count.Loc.Start] + strconv.Itoa(limitCount) + statement[count.Loc.End:]
			}
			if !force {
				// Spanner leaves a non-literal LIMIT alone, as it does for a
				// plain query.
				return statement
			}
			// BigQuery must enforce a literal cap; a second pipe LIMIT does.
		}
	}
	end := pipe.Loc.End
	return statement[:end] + " |> LIMIT " + strconv.Itoa(limitCount) + statement[end:]
}

func RewriteResultLimit(statement string, limitCount int) string {
	file, err := parseFile(statement)
	if err != nil || len(file.Stmts) != 1 {
		return statement
	}
	switch n := file.Stmts[0].(type) {
	case *ast.SubpipelineStmt:
		return rewritePipeResultLimit(statement, n.Pipeline, limitCount, false)
	case *ast.PipeStatement:
		return rewritePipeResultLimit(statement, n.Pipeline, limitCount, false)
	case *ast.QueryStmt:
		if n.ForUpdate {
			return statement
		}
		// A pipe query's LIMIT is a `|> LIMIT` operator, not a query-level
		// clause, so it is capped by lowering that operator or appending one.
		if n.Pipe != nil {
			return rewritePipeResultLimit(statement, n.Pipe, limitCount, false)
		}
		if n.LimitAll {
			// LIMIT ALL [OFFSET m] (googlesql.tm:8504) has no count to lower,
			// and a second LIMIT after it is a syntax error, so it is left
			// unchanged, as a non-literal LIMIT is.
			return statement
		}
		if n.Limit == nil {
			end := n.Loc.End
			return statement[:end] + " LIMIT " + strconv.Itoa(limitCount) + statement[end:]
		}
		count, ok := n.Limit.(*ast.Literal)
		if !ok || count.Kind != ast.LitInt || count.Ival <= int64(limitCount) {
			return statement
		}
		return statement[:count.Loc.Start] + strconv.Itoa(limitCount) + statement[count.Loc.End:]
	}
	return statement
}
