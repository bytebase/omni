package analysis

import (
	"github.com/bytebase/omni/redshift/ast"
)

func hasExplainAnalyze(stmt *ast.ExplainStmt) bool {
	if stmt == nil || stmt.Options == nil {
		return false
	}
	for _, item := range stmt.Options.Items {
		def, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		if def.Defname == "analyze" {
			return true
		}
	}
	return false
}

func isMixedQuery(m ColumnSet) (bool, bool) {
	hasSystem, hasUser := false, false
	for table := range m {
		if isSystemResource(table) {
			hasSystem = true
		} else {
			hasUser = true
		}
	}
	if hasSystem && hasUser {
		return false, true
	}
	return !hasUser && hasSystem, false
}

func isSystemResource(resource ColumnResource) bool {
	if isSystemSchema(resource.Schema) {
		return true
	}
	if resource.Schema == "" && isSystemView(resource.Table) {
		return true
	}
	if resource.Schema == "" && isSystemTable(resource.Table) {
		return true
	}
	return false
}
