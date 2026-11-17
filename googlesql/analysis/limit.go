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
		return wrapResultLimit(statement, limit)
	}
	return RewriteResultLimit(statement, limit)
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
