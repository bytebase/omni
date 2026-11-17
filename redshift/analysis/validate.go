package analysis

import (
	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
)

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows: (readOnly, returnsData,
// error).
//
// Decision is AST-based: a SELECT that carries INTO creates a table, and a CTE
// that holds DML writes, so neither runs read-only. EXPLAIN ANALYZE runs the
// statement it plans, which is read-only only when that statement is a plain
// SELECT.
//
// A syntax error is returned; a caller rejects malformed read-only SQL before
// execution.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	parsed, err := redshift.Parse(statement)
	if err != nil {
		return false, false, err
	}

	seen := false
	returnsData = true
	for _, stmt := range parsed {
		if stmt.Empty() {
			continue
		}
		seen = true
		readOnly, stmtReturnsData := validateEditorStatement(stmt.AST)
		if !readOnly {
			return false, false, nil
		}
		if !stmtReturnsData {
			returnsData = false
		}
	}
	return seen, returnsData, nil
}

func validateEditorStatement(node ast.Node) (bool, bool) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		if hasSelectIntoClause(n) || hasDMLInWithClause(n.WithClause) {
			return false, false
		}
		return true, true
	case *ast.ExplainStmt:
		if hasExplainAnalyze(n) {
			if selectStmt, ok := n.Query.(*ast.SelectStmt); ok && !hasSelectIntoClause(selectStmt) && !hasDMLInWithClause(selectStmt.WithClause) {
				return true, false
			}
			return false, false
		}
		return true, true
	case *ast.RedshiftShowStmt, *ast.RedshiftDescStmt, *ast.VariableShowStmt:
		return true, true
	case *ast.VariableSetStmt:
		return true, false
	default:
		return false, false
	}
}

func hasDMLInWithClause(withClause *ast.WithClause) bool {
	if withClause == nil || withClause.Ctes == nil {
		return false
	}
	for _, item := range withClause.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		if hasDMLNode(cte.Ctequery) {
			return true
		}
	}
	return false
}

func hasDMLNode(node ast.Node) bool {
	hasDML := false
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			hasDML = true
			return false
		default:
			return !hasDML
		}
	})
	return hasDML
}
