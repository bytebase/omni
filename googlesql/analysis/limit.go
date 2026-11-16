package analysis

import (
	"strconv"

	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
)

// RewriteResultLimit sets the LIMIT of the outermost query in statement to at
// most limitCount, lowering a larger literal LIMIT and adding one where the
// query has none. It never wraps the statement in a subquery, where Spanner
// rejects a WITH clause, so a query it cannot rewrite — a non-literal LIMIT, a
// shape the parser does not read, a FOR UPDATE that a read-only transaction
// refuses anyway — is returned unchanged.
func RewriteResultLimit(statement string, limitCount int) string {
	file, err := parser.Parse(statement)
	if err != nil || len(file.Stmts) != 1 {
		return statement
	}
	query, ok := file.Stmts[0].(*ast.QueryStmt)
	if !ok || query.ForUpdate {
		return statement
	}
	if query.Limit == nil {
		end := query.Loc.End
		return statement[:end] + " LIMIT " + strconv.Itoa(limitCount) + statement[end:]
	}
	count, ok := query.Limit.(*ast.Literal)
	if !ok || count.Kind != ast.LitInt || count.Ival <= int64(limitCount) {
		return statement
	}
	return statement[:count.Loc.Start] + strconv.Itoa(limitCount) + statement[count.Loc.End:]
}
