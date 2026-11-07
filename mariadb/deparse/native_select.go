package deparse

import ast "github.com/bytebase/omni/mariadb/ast"

// ApplySelectAliases applies view or CTE column names to the stored select list.
func ApplySelectAliases(stmt *ast.SelectStmt, names []string) {
	if stmt == nil || len(names) == 0 {
		return
	}
	leaf := ast.LeftmostQueryLeaf(stmt)
	if leaf == nil {
		return
	}
	for i, name := range names {
		if i >= len(leaf.TargetList) {
			break
		}
		if rt, ok := leaf.TargetList[i].(*ast.ResTarget); ok {
			rt.Name = name
		} else {
			leaf.TargetList[i] = &ast.ResTarget{Val: leaf.TargetList[i], Name: name}
		}
	}
}
func selectWithAliases(stmt *ast.SelectStmt, names []string) *ast.SelectStmt {
	ApplySelectAliases(stmt, names)
	return stmt
}

func expandNamedWindows(stmt *ast.SelectStmt) {
	if len(stmt.WindowClause) == 0 {
		return
	}
	named := map[string]*ast.WindowDef{}
	for _, wd := range stmt.WindowClause {
		named[wd.Name] = wd
	}
	for i, target := range stmt.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok && rt.Name == "" {
			rt.Name = autoAlias(rt.Val, deparseExpr(rt.Val), i+1)
		}
	}
	ast.Inspect(stmt, func(n ast.Node) bool {
		if q, ok := n.(*ast.SelectStmt); ok && q != stmt {
			return false
		}
		if call, ok := n.(*ast.FuncCallExpr); ok && call.Over != nil {
			if wd := named[call.Over.RefName]; wd != nil {
				copy := *wd
				call.Over = &copy
			}
		}
		return true
	})
}
