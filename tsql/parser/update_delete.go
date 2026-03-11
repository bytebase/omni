// Package parser - update_delete.go implements T-SQL UPDATE and DELETE statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/update-transact-sql
//
//	UPDATE [TOP (...)] table SET col=expr, ... [OUTPUT ...] [FROM ...] [WHERE ...]
func (p *Parser) parseUpdateStmt() *nodes.UpdateStmt {
	loc := p.pos()
	p.advance() // consume UPDATE

	stmt := &nodes.UpdateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top = p.parseTopClause()
	}

	// Table name
	stmt.Relation = p.parseTableRef()

	// SET clause
	if _, err := p.expect(kwSET); err == nil {
		stmt.SetClause = p.parseSetClauseList()
	}

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause = p.parseFromClause()
	}

	// WHERE clause
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/delete-transact-sql
//
//	DELETE [TOP (...)] [FROM] table [OUTPUT ...] [FROM ...] [WHERE ...]
func (p *Parser) parseDeleteStmt() *nodes.DeleteStmt {
	loc := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top = p.parseTopClause()
	}

	// Optional FROM before table name
	p.match(kwFROM)

	// Table name
	stmt.Relation = p.parseTableRef()

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	// FROM clause (second FROM for join)
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause = p.parseFromClause()
	}

	// WHERE clause
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSetClauseList parses a comma-separated list of SET assignments.
//
//	set_clause_list = set_clause { ',' set_clause }
//	set_clause = column_ref '=' expr | @variable '=' expr
func (p *Parser) parseSetClauseList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseSetClause()
		if item == nil {
			break
		}
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}

// parseSetClause parses a single SET assignment: column = expr or @var = expr.
func (p *Parser) parseSetClause() *nodes.SetExpr {
	loc := p.pos()

	se := &nodes.SetExpr{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == tokVARIABLE {
		se.Variable = p.cur.Str
		p.advance()
	} else {
		// Parse column reference (not a full expression - just the column name)
		se.Column = p.parseSetTarget()
		if se.Column == nil {
			return nil
		}
	}

	// Expect =
	if _, err := p.expect('='); err != nil {
		return nil
	}

	se.Value = p.parseExpr()
	se.Loc.End = p.pos()
	return se
}

// parseSetTarget parses the left side of a SET assignment (qualified column name).
func (p *Parser) parseSetTarget() *nodes.ColumnRef {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil
	}

	ref := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: loc},
	}

	// Check for qualified name: table.column
	if p.cur.Type == '.' {
		p.advance()
		col, ok := p.parseIdentifier()
		if ok {
			ref.Table = name
			ref.Column = col
		}
	}

	return ref
}
