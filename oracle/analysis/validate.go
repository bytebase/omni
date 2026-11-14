package analysis

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

// ValidateQuery reports whether a script may run in the SQL editor's
// read-only mode, and whether it returns rows. Only SELECT and EXPLAIN PLAN
// are allowed; a SQL*Plus command or an empty script is not, and neither is
// a query whose WITH clause declares PL/SQL functions or procedures, which
// can write through PRAGMA AUTONOMOUS_TRANSACTION or EXECUTE IMMEDIATE. The
// error is Parse's.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	for _, segment := range parser.Split(statement) {
		if segment.Kind == parser.SegmentSQLPlusCommand && !segment.Empty() {
			return false, false, nil
		}
	}

	list, err := Parse(statement)
	if err != nil {
		return false, false, err
	}
	if list == nil || len(list.Items) == 0 {
		return false, false, nil
	}

	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			return false, false, nil
		}
		switch raw.Stmt.(type) {
		case *ast.SelectStmt, *ast.ExplainPlanStmt:
		default:
			return false, false, nil
		}
		if declaresPLSQL(raw.Stmt) {
			return false, false, nil
		}
	}

	return true, true, nil
}

// declaresPLSQL reports whether any WITH clause in node declares PL/SQL.
func declaresPLSQL(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if w, ok := n.(*ast.WithClause); ok && w.PLSQLDeclarations != nil && len(w.PLSQLDeclarations.Items) > 0 {
			found = true
		}
		return !found
	})
	return found
}
