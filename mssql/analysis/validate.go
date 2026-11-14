package analysis

import (
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

// ValidateQuery reports whether parsed statements may run in the SQL editor's
// read-only mode, and whether they all return data. Only SELECT without INTO
// and SET of the read-only session options are allowed; a session option
// returns no data, so it clears allQuery. DECLARE and SET @variable are
// refused. No statements is not valid.
func ValidateQuery(stmts []mssql.Statement) (valid bool, allQuery bool) {
	if len(stmts) == 0 {
		return false, false
	}
	allQuery = true
	for _, s := range stmts {
		if s.Empty() {
			continue
		}
		if isReadOnlySelect(s.AST) {
			continue
		}
		if !isReadOnlySessionOption(s.AST) {
			return false, false
		}
		allQuery = false
	}
	return true, allQuery
}

// readOnlySessionOptions are the SET options a read-only request may carry: they
// write nothing and add no result set, unlike SET STATISTICS XML.
var readOnlySessionOptions = map[string]bool{
	"TRANSACTION ISOLATION LEVEL": true,
	"NOCOUNT":                     true,
}

// isReadOnlySessionOption reports whether node is a SET whose options are all
// read-only session options. One SET can name several, as in SET NOCOUNT, FMTONLY ON.
func isReadOnlySessionOption(node ast.Node) bool {
	set, ok := node.(*ast.SetOptionStmt)
	if !ok {
		return false
	}
	for _, option := range set.Options {
		if !readOnlySessionOptions[option] {
			return false
		}
	}
	return true
}

// isReadOnlySelect returns true for SELECT statements without an INTO target.
// Non-SELECT nodes and SELECT ... INTO are rejected — INTO materialises a new
// table which is a DDL-like side effect.
func isReadOnlySelect(node ast.Node) bool {
	sel, ok := node.(*ast.SelectStmt)
	if !ok {
		return false
	}
	return !HasSelectInto(sel)
}

// HasSelectInto reports whether sel (or any branch of its set operations)
// carries an INTO clause.
func HasSelectInto(sel *ast.SelectStmt) bool {
	return selectIntoTarget(sel) != nil
}
