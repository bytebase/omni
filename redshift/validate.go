package redshift

import (
	"strings"

	"github.com/bytebase/omni/redshift/ast"
)

// ValidateSQLForEditor parses Redshift SQL and reports whether it is safe for
// read-only SQL Editor execution, and whether it is a plain EXPLAIN plan query.
func ValidateSQLForEditor(sql string) (bool, bool, error) {
	stmts, err := Parse(sql)
	if err != nil {
		return false, false, err
	}

	seen := false
	hasExplain := false
	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		seen = true
		readonly, explain := isReadOnlyEditorStatement(stmt.AST)
		if !readonly {
			return false, false, nil
		}
		hasExplain = hasExplain || explain
	}
	return seen, hasExplain, nil
}

func isReadOnlyEditorStatement(node ast.Node) (bool, bool) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		if n.IntoClause != nil {
			return false, false
		}
		return true, false
	case *ast.RedshiftShowStmt, *ast.VariableShowStmt:
		return true, false
	case *ast.VariableSetStmt:
		return true, false
	case *ast.ExplainStmt:
		if explainHasAnalyze(n) {
			return false, false
		}
		if !isReadOnlyExplainQuery(n.Query) {
			return false, false
		}
		return true, true
	default:
		return false, false
	}
}

func isReadOnlyExplainQuery(node ast.Node) bool {
	switch node.(type) {
	case *ast.SelectStmt:
		return node.(*ast.SelectStmt).IntoClause == nil
	default:
		return false
	}
}

func explainHasAnalyze(stmt *ast.ExplainStmt) bool {
	if stmt == nil || stmt.Options == nil {
		return false
	}
	for _, item := range stmt.Options.Items {
		def, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		if strings.EqualFold(def.Defname, "analyze") {
			return true
		}
	}
	return false
}
