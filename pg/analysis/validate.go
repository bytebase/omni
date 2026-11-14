package analysis

import (
	"github.com/bytebase/omni/pg/ast"
)

// ValidateQuery reports whether parsed statements may run in the SQL editor's
// read-only mode, and whether they all return data rather than set session
// state. A nil node is skipped.
//
// Validation rules:
//  1. Allow: SELECT statements
//  2. Allow: EXPLAIN statements (but not EXPLAIN ANALYZE for non-SELECT)
//  3. Allow: SHOW/SET statements (SET is considered executable)
//  4. Allow: CTEs whose every term is a read (SELECT/VALUES/TABLE/set-op/RECURSIVE);
//     reject a SELECT with a data-modifying CTE term — a fail-safe whitelist.
//  5. Reject: SELECT ... INTO (creates a table) and all other statements (DDL, DML).
func ValidateQuery(nodes []ast.Node) (readOnly bool, returnsData bool) {
	var hasExecute bool

	for _, node := range nodes {
		if node == nil {
			continue
		}

		switch n := node.(type) {
		case *ast.SelectStmt:
			// SELECT is read-only unless it writes: SELECT ... INTO creates a table,
			// or a data-modifying CTE smuggles a write into the read path.
			if isWriteSelect(n) {
				return false, false
			}

		case *ast.ExplainStmt:
			if isExplainAnalyzeStmt(n) {
				// EXPLAIN ANALYZE executes the query, so it must be a read-only SELECT.
				sel, ok := n.Query.(*ast.SelectStmt)
				if !ok || isWriteSelect(sel) {
					return false, false
				}
			}

		case *ast.VariableSetStmt:
			hasExecute = true

		case *ast.VariableShowStmt:
			// SHOW is allowed

		default:
			// All other statements (DDL, DML) are not allowed
			return false, false
		}
	}

	return true, !hasExecute
}

// isWriteSelect reports whether a SelectStmt actually writes — and so must not take
// the read-only editor path. A SELECT writes if it has an INTO target (SELECT ...
// INTO creates a table) or any of its CTE terms is data-modifying.
//
// The read-only gate uses it; classifyQueryType shares containsWriteCTE so the
// permission check treats the same SELECTs as DML.
func isWriteSelect(n *ast.SelectStmt) bool {
	if n == nil {
		return false
	}
	if hasIntoClause(n) {
		return true
	}
	return containsWriteCTE(n)
}
