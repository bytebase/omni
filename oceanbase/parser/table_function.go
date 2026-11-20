package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// This file reads OceanBase's two table functions in a FROM clause, RB_ITERATE
// and UNNEST. MySQL has neither. Both take an optional alias, written with or
// without AS, and a column list in parentheses after it.
//
// Ref: rb_iterate_expr and unnest_expr in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6, lines 25536-25583). The rule's rows in testdata/compat.tsv hold
// each alternative, and its near-misses, to the pinned server.

// parseRbIterate parses rb_iterate_expr, at its RB_ITERATE:
//
//	RB_ITERATE '(' simple_expr ')' [AS] relation_name ['(' relation_name ')']
func (p *Parser) parseRbIterate() (*nodes.RbIterateExpr, error) {
	expr := &nodes.RbIterateExpr{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // RB_ITERATE
	arg, err := p.parseTableFunctionArg()
	if err != nil {
		return nil, err
	}
	expr.Arg = arg
	if expr.Alias, err = p.parseTableFunctionAlias(); err != nil {
		return nil, err
	}
	if expr.Alias != "" && p.cur.Type == '(' {
		p.advance()
		column, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		expr.Column = column
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	expr.Loc.End = p.prev.End
	return expr, nil
}

// parseUnnest parses unnest_expr, at its UNNEST:
//
//	UNNEST '(' simple_expr_list ')' [AS] relation_name ['(' relation_name_list ')']
func (p *Parser) parseUnnest() (*nodes.UnnestExpr, error) {
	expr := &nodes.UnnestExpr{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // UNNEST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		arg, err := p.parseExprPrec(precCollate) // simple_expr
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		expr.Args = append(expr.Args, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	var err error
	if expr.Alias, err = p.parseTableFunctionAlias(); err != nil {
		return nil, err
	}
	if expr.Alias != "" && p.cur.Type == '(' {
		cols, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		expr.Columns = cols
	}
	expr.Loc.End = p.prev.End
	return expr, nil
}

// parseTableCollection parses tbl_name's TABLE '(' simple_expr ')' [[AS] relation_name]
// (sql_parser_mysql_mode.y, v4.3.5_CE_BP6, lines 14476-14487).
func (p *Parser) parseTableCollection() (*nodes.TableCollectionExpr, error) {
	expr := &nodes.TableCollectionExpr{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // TABLE
	arg, err := p.parseTableFunctionArg()
	if err != nil {
		return nil, err
	}
	expr.Expr = arg
	if expr.Alias, err = p.parseTableFunctionAlias(); err != nil {
		return nil, err
	}
	expr.Loc.End = p.prev.End
	return expr, nil
}

// parseTableFunctionArg reads the '(' and one simple_expr of a table function,
// which is where both RB_ITERATE and UNNEST open their arguments.
func (p *Parser) parseTableFunctionArg() (nodes.ExprNode, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	arg, err := p.parseExprPrec(precCollate) // simple_expr
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return arg, nil
}

// parseTableFunctionAlias reads the [AS] relation_name a table function takes,
// and reports "" where none stands.
func (p *Parser) parseTableFunctionAlias() (string, error) {
	if _, ok := p.match(kwAS); ok {
		alias, _, err := p.parseIdent()
		return alias, err
	}
	if !p.isIdentToken() || p.isSelectTerminator() {
		return "", nil
	}
	alias, _, err := p.parseIdent()
	return alias, err
}
