package analysis

import (
	"github.com/bytebase/omni/mariadb/ast"
)

// ValidateQuery reports whether parsed statements may run in the SQL editor's
// read-only mode, and whether they return rows rather than only setting
// session state. Only SELECT, EXPLAIN, SHOW, SET and DESCRIBE are allowed, and
// an INSERT, REPLACE or DELETE when it carries RETURNING: the write answers with
// a result set, which the editor reads like a query's.
// EXPLAIN ANALYZE is treated as non-read-only since it actually executes the query.
func ValidateQuery(nodes []ast.Node) (readOnly bool, returnsData bool) {
	hasExecute := false
	returnsRows := false
	readOnly = true
	for _, node := range nodes {
		switch stmt := node.(type) {
		case *ast.SelectStmt:
			// INTO OUTFILE and INTO DUMPFILE write a file on the server, so
			// they are writes (pg refuses its own SELECT ... INTO likewise).
			// INTO <variable> returns no rows: the SET case, not a write.
			writesFile, assignsVariables := selectIntoTargets(stmt)
			if writesFile {
				return false, false
			}
			if assignsVariables {
				hasExecute = true
			}
		case *ast.ExplainStmt:
			if stmt.Analyze {
				readOnly = false
			}
		case *ast.ShowStmt:
			// SHOW is always allowed.
		case *ast.SetStmt, *ast.SetPasswordStmt, *ast.SetDefaultRoleStmt, *ast.SetRoleStmt, *ast.SetResourceGroupStmt, *ast.SetTransactionStmt:
			hasExecute = true
		case *ast.InsertStmt:
			// INSERT and REPLACE share a node; Returning is set for both.
			if len(stmt.Returning) == 0 {
				return false, false
			}
			returnsRows, readOnly = true, false
		case *ast.DeleteStmt:
			// MariaDB's DELETE ... RETURNING is single-table only, which the
			// parser enforces; a DELETE without it still answers no rows.
			if len(stmt.Returning) == 0 {
				return false, false
			}
			returnsRows, readOnly = true, false
		default:
			return false, false
		}
	}
	return readOnly, returnsRows || !hasExecute
}

// selectIntoTargets reports what a SELECT's INTO clause targets. It searches
// the whole statement rather than its root because the parser attaches INTO to
// an arm of a set operation or to a parenthesized query, not necessarily to the
// node the caller holds — the same reason pg's intoClause walks its arms.
//
// A subquery cannot carry INTO in MySQL, so reaching into the tree cannot
// over-match a legal statement; if the grammar ever allowed it, the answer here
// would be conservative rather than permissive.
func selectIntoTargets(n *ast.SelectStmt) (writesFile, assignsVariables bool) {
	ast.Inspect(n, func(node ast.Node) bool {
		sel, ok := node.(*ast.SelectStmt)
		if !ok || sel.Into == nil {
			return true
		}
		// The clause form decides, not the filename it carries: those fields
		// hold the decoded path, so "INTO OUTFILE ''" leaves both empty while
		// still being a file write.
		if len(sel.Into.Vars) > 0 {
			assignsVariables = true
		} else {
			writesFile = true
		}
		return true
	})
	return writesFile, assignsVariables
}
