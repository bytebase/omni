package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/UPDATE.html
//
//	UPDATE [ hint ] table [ alias ]
//	    SET col = expr [, col = expr ...]
//	    [ WHERE condition ]
//	    [ RETURNING expr [, ...] INTO var [, ...] ]
//	    [ LOG ERRORS ... ]
func (p *Parser) parseUpdateStmt() *nodes.UpdateStmt {
	start := p.pos()
	p.advance() // consume UPDATE

	stmt := &nodes.UpdateStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// SET
	if p.cur.Type == kwSET {
		p.advance()
		stmt.SetClauses = p.parseSetClauses()
	}

	// WHERE
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause = p.parseExpr()
	}

	// RETURNING ... INTO ...
	if p.cur.Type == kwRETURNING {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSetClauses parses comma-separated SET clauses.
func (p *Parser) parseSetClauses() *nodes.List {
	list := &nodes.List{}

	for {
		sc := p.parseSetClause()
		if sc == nil {
			break
		}
		list.Items = append(list.Items, sc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list
}

// parseSetClause parses a single SET clause: col = expr.
func (p *Parser) parseSetClause() *nodes.SetClause {
	start := p.pos()

	// Multi-column: (col1, col2) = (subquery)
	if p.cur.Type == '(' {
		p.advance()
		sc := &nodes.SetClause{
			Columns: &nodes.List{},
			Loc:     nodes.Loc{Start: start},
		}
		for {
			col := p.parseColumnRef()
			if col != nil {
				sc.Columns.Items = append(sc.Columns.Items, col)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		if p.cur.Type == '=' {
			p.advance()
		}
		sc.Value = p.parseExpr()
		sc.Loc.End = p.pos()
		return sc
	}

	// Single column: col = expr
	if !p.isIdentLike() {
		return nil
	}

	col := p.parseColumnRef()
	sc := &nodes.SetClause{
		Column: col,
		Loc:    nodes.Loc{Start: start},
	}

	if p.cur.Type == '=' {
		p.advance()
	}

	sc.Value = p.parseExpr()
	sc.Loc.End = p.pos()
	return sc
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/DELETE.html
//
//	DELETE [ hint ] [ FROM ] table [ alias ]
//	    [ WHERE condition ]
//	    [ RETURNING expr [, ...] INTO var [, ...] ]
//	    [ LOG ERRORS ... ]
func (p *Parser) parseDeleteStmt() *nodes.DeleteStmt {
	start := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// Optional FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// WHERE
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause = p.parseExpr()
	}

	// RETURNING ... INTO ...
	if p.cur.Type == kwRETURNING {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
