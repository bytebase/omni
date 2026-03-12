package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseFuncExprWindowless parses func_expr_windowless.
//
// func_expr_windowless is func_application or func_expr_common_subexpr,
// but without OVER/FILTER/WITHIN GROUP (no window function support).
//
// It is used in places where window functions are not allowed, such as
// index expressions, partition key expressions, func_table, etc.
//
// Ref: gram.y func_expr_windowless
//
//	func_expr_windowless:
//	    func_application
//	    | func_expr_common_subexpr
func (p *Parser) parseFuncExprWindowless() nodes.Node {
	// func_expr_common_subexpr keywords are handled by parseCExpr which
	// already dispatches COALESCE, GREATEST, LEAST, CAST, etc.
	// For func_application, we need a name followed by '('.
	// We delegate to parseCExpr which handles both cases.
	return p.parseCExpr()
}

// parseWindowClause parses an optional WINDOW clause in a SELECT statement.
//
// Ref: gram.y opt_window_clause / window_clause
//
//	opt_window_clause:
//	    WINDOW window_definition_list
//	    | /* EMPTY */
func (p *Parser) parseWindowClause() *nodes.List {
	if p.cur.Type != WINDOW {
		return nil
	}
	p.advance() // consume WINDOW
	return p.parseWindowDefinitionList()
}

// parseWindowDefinitionList parses a comma-separated list of window definitions.
//
//	window_definition_list:
//	    window_definition
//	    | window_definition_list ',' window_definition
func (p *Parser) parseWindowDefinitionList() *nodes.List {
	first := p.parseWindowDefinition()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseWindowDefinition())
	}
	return &nodes.List{Items: items}
}

// parseWindowDefinition parses a single named window definition.
//
//	window_definition:
//	    ColId AS window_specification
func (p *Parser) parseWindowDefinition() *nodes.WindowDef {
	name, _ := p.parseColId()
	p.expect(AS)
	wd := p.parseWindowSpecification().(*nodes.WindowDef)
	wd.Name = name
	wd.Loc = nodes.NoLoc()
	return wd
}
