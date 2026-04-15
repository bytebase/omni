package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// DELETE statement parser
// ---------------------------------------------------------------------------

// parseDeleteStmt parses a DELETE statement:
//
//	DELETE FROM table [alias] [USING source [, source ...]] [WHERE cond]
func (p *Parser) parseDeleteStmt() (*ast.DeleteStmt, error) {
	deleteTok, err := p.expect(kwDELETE)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.DeleteStmt{
		Target: target,
		Loc:    ast.Loc{Start: deleteTok.Loc.Start},
	}

	// Optional table alias
	alias, hasAlias := p.parseOptionalAlias()
	if hasAlias {
		stmt.Alias = alias
	}

	// Optional USING clause
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		using, err := p.parseTableOrQueryList()
		if err != nil {
			return nil, err
		}
		stmt.Using = using
	}

	// Optional WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTableOrQueryList parses a comma-separated list of table references or
// subqueries. Used by DELETE USING.
//
//	table_or_query := object_name [alias] | '(' subquery ')' [alias]
func (p *Parser) parseTableOrQueryList() ([]ast.Node, error) {
	var items []ast.Node

	item, err := p.parseFromItem()
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		item, err = p.parseFromItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
