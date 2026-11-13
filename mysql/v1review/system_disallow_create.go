package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemProcedureDisallowCreate reports each top-level CREATE PROCEDURE that
// names its procedure. Statements nested in a routine body are not judged.
// Every finding anchors the whole statement. Ported from Bytebase's V1 rule
// SYSTEM_PROCEDURE_DISALLOW_CREATE.
func SystemProcedureDisallowCreate(stmts []ast.Node) []v1review.Finding {
	return created(stmts, v1review.ProcedureCreated, func(node ast.Node) (ast.Loc, bool) {
		n, ok := node.(*ast.CreateFunctionStmt)
		return locOf(n), ok && n.IsProcedure && n.Name != nil
	})
}

// SystemFunctionDisallowCreate reports each top-level CREATE FUNCTION that
// names its function, a loadable function's CREATE AGGREGATE FUNCTION
// included. Statements nested in a routine body are not judged. Every
// finding anchors the whole statement. Ported from Bytebase's V1 rule
// SYSTEM_FUNCTION_DISALLOW_CREATE.
func SystemFunctionDisallowCreate(stmts []ast.Node) []v1review.Finding {
	return created(stmts, v1review.FunctionCreated, func(node ast.Node) (ast.Loc, bool) {
		n, ok := node.(*ast.CreateFunctionStmt)
		return locOf(n), ok && !n.IsProcedure && n.Name != nil
	})
}

// SystemEventDisallowCreate reports each top-level CREATE EVENT that names
// its event. Statements nested in a routine body are not judged. Every
// finding anchors the whole statement. Ported from Bytebase's V1 rule
// SYSTEM_EVENT_DISALLOW_CREATE.
func SystemEventDisallowCreate(stmts []ast.Node) []v1review.Finding {
	return created(stmts, v1review.EventCreated, func(node ast.Node) (ast.Loc, bool) {
		n, ok := node.(*ast.CreateEventStmt)
		if !ok {
			return ast.Loc{}, false
		}
		return n.Loc, n.Name != ""
	})
}

// SystemViewDisallowCreate reports each top-level CREATE VIEW that names its
// view. Statements nested in a routine body are not judged. Every finding
// anchors the whole statement. Ported from Bytebase's V1 rule
// SYSTEM_VIEW_DISALLOW_CREATE.
func SystemViewDisallowCreate(stmts []ast.Node) []v1review.Finding {
	return created(stmts, v1review.ViewCreated, func(node ast.Node) (ast.Loc, bool) {
		n, ok := node.(*ast.CreateViewStmt)
		if !ok {
			return ast.Loc{}, false
		}
		return n.Loc, n.Name != nil
	})
}

func locOf(n *ast.CreateFunctionStmt) ast.Loc {
	if n == nil {
		return ast.Loc{}
	}
	return n.Loc
}

// created reports, with code, each statement match says creates a forbidden
// object, anchored at the location match returns.
func created(stmts []ast.Node, code v1review.Code, match func(ast.Node) (ast.Loc, bool)) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if loc, ok := match(stmt); ok {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: code, Start: start, End: end})
		}
	}
	return findings
}
