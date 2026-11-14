package analysis

import (
	"fmt"

	omniast "github.com/bytebase/omni/tidb/ast"
	omnideparse "github.com/bytebase/omni/tidb/deparse"
	"github.com/bytebase/omni/tidb/parser"
)

// RewriteResultLimit rewrites one SELECT, including a set operation, so it
// returns at most limit rows, and restores it from the AST. An integer LIMIT
// larger than limit becomes limit, a non-constant one is kept, and a SELECT
// with no LIMIT gets one. A statement that is not a SELECT is returned
// unchanged.
//
// An error (a parse error, other than one statement) means the caller should
// wrap the statement instead.
func RewriteResultLimit(statement string, limit int) (string, error) {
	list, err := parser.Parse(statement)
	if err != nil {
		return "", fmt.Errorf("failed to parse tidb statement: %s: %w", statement, err)
	}
	if list == nil || len(list.Items) != 1 {
		return "", fmt.Errorf("expect one single statement in the query, %s", statement)
	}
	// omni folds UNION/INTERSECT/EXCEPT into SelectStmt, so a single case covers
	// both a plain SELECT and a set operation; the trailing LIMIT lives on the
	// outermost SelectStmt either way.
	stmt, ok := list.Items[0].(*omniast.SelectStmt)
	if !ok {
		return statement, nil
	}
	if stmt.Limit != nil && stmt.Limit.Count != nil {
		if v, ok := stmt.Limit.Count.(*omniast.IntLit); ok {
			userLimit := int(v.Value)
			if limit < userLimit {
				userLimit = limit
			}
			stmt.Limit.Count = &omniast.IntLit{Value: int64(userLimit)}
		}
	} else {
		stmt.Limit = &omniast.Limit{
			Count: &omniast.IntLit{Value: int64(limit)},
		}
	}
	// String literals are restored without a charset introducer.
	return omnideparse.RestoreSelectWith(stmt, omnideparse.RestoreStringWithoutCharset), nil
}
