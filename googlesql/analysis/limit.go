package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
)

// StatementWithResultLimit returns statement rewritten so the server returns at
// most limit rows.
//
// BigQuery's caller has already confirmed the statement is a SELECT, so it is
// wrapped in a CTE: this package does not rewrite BigQuery's LIMIT placement
// from the AST. Spanner rejects a WITH clause outside a subquery, so its rewrite
// moves the outermost query's own LIMIT and never wraps.
func StatementWithResultLimit(statement string, limit int, dialect Dialect) string {
	if dialect == DialectBigQuery {
		// A statement that is only a pipe suffix is not a query expression, so
		// a CTE wrapper around it is not valid SQL. Cap it with a pipe LIMIT.
		if pipe, ok := suffixOnlyPipeline(statement); ok {
			if limit <= 0 {
				return statement
			}
			return rewritePipeResultLimit(statement, pipe, limit, true)
		}
		return wrapResultLimit(statement, limit)
	}
	return RewriteResultLimit(statement, limit)
}

// suffixOnlyPipeline returns the pipeline of a statement that is only a pipe
// suffix or a table-returning statement with one, when statement parses as
// exactly that.
func suffixOnlyPipeline(statement string) (*ast.Subpipeline, bool) {
	file, err := parser.Parse(statement)
	if err != nil || len(file.Stmts) != 1 {
		return nil, false
	}
	switch n := file.Stmts[0].(type) {
	case *ast.SubpipelineStmt:
		return n.Pipeline, true
	case *ast.PipeStatement:
		return n.Pipeline, true
	}
	return nil, false
}

// wrapResultLimit is BigQuery's result cap. A limit of zero or less omits the
// clause, which leaves the wrapper capping nothing.
func wrapResultLimit(statement string, limit int) string {
	limitPart := ""
	if limit > 0 {
		limitPart = fmt.Sprintf(" LIMIT %d", limit)
	}

	// The CTE holds one statement, without its terminator. The splitter returns
	// each statement with the ';' and anything after it excluded, so a trailing
	// comment that follows the terminator does not end up inside the CTE.
	inner := trimStatement(statement)
	if segments := parser.Split(statement); len(segments) == 1 {
		inner = strings.TrimSpace(statement[segments[0].ByteStart:segments[0].ByteEnd])
	}

	// A trailing comment would otherwise comment out the closing parenthesis;
	// start the wrapper's own SQL on a new line instead.
	closing := ""
	if hasTrailingComment(inner) {
		closing = "\n"
	}
	return fmt.Sprintf("WITH result AS (%s%s) SELECT * FROM result%s;", inner, closing, limitPart)
}

// hasTrailingComment reports whether statement ends in trivia after its last
// token, which is whitespace or a comment. A line comment there would swallow
// the character that follows it on the same line.
func hasTrailingComment(statement string) bool {
	tokens, _ := parser.Tokenize(statement)
	end := 0
	for _, token := range tokens {
		if token.Type == 0 {
			break
		}
		end = token.Loc.End
	}
	return strings.TrimSpace(statement[end:]) != ""
}

// trimStatement removes the whitespace and trailing semicolons the CTE wrapper
// would otherwise leave stray inside itself.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}

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
	file, err := parser.Parse(statement)
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
