package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseInsertStmt parses an INSERT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/INSERT.html
//
//	INSERT [ hint ] INTO table [ alias ] [ ( col1, col2, ... ) ]
//	  { VALUES ( expr, ... ) | subquery }
//	  [ RETURNING expr [, ...] INTO bind [, ...] ]
//	  [ LOG ERRORS [ INTO table ] [ ( tag ) ] [ REJECT LIMIT { integer | UNLIMITED } ] ]
//
//	INSERT [ hint ] ALL
//	  INTO table [ ( col1, ... ) ] VALUES ( expr, ... )
//	  [ INTO table [ ( col1, ... ) ] VALUES ( expr, ... ) ]*
//	  subquery
//
//	INSERT [ hint ] FIRST
//	  WHEN condition THEN INTO table [ ( col1, ... ) ] VALUES ( expr, ... ) [INTO ...]
//	  [ WHEN condition THEN INTO table ... ]*
//	  [ ELSE INTO table [ ( col1, ... ) ] VALUES ( expr, ... ) [INTO ...] ]
//	  subquery
func (p *Parser) parseInsertStmt() *nodes.InsertStmt {
	start := p.pos()
	stmt := &nodes.InsertStmt{
		Loc: nodes.Loc{Start: start},
	}

	p.advance() // consume INSERT

	// Hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// INSERT ALL ...
	if p.cur.Type == kwALL {
		p.advance()
		return p.parseMultiTableInsert(stmt, nodes.INSERT_ALL)
	}

	// INSERT FIRST ...
	if p.cur.Type == kwFIRST {
		p.advance()
		return p.parseMultiTableInsert(stmt, nodes.INSERT_FIRST)
	}

	// Single-table INSERT INTO ...
	stmt.InsertType = nodes.INSERT_SINGLE

	// INTO
	if p.cur.Type == kwINTO {
		p.advance()
	}

	// Table name
	if p.isIdentLike() {
		stmt.Table = p.parseObjectName()
	}

	// Optional alias
	if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// Optional column list
	if p.cur.Type == '(' {
		stmt.Columns = p.parseParenColumnList()
	}

	// VALUES (...) or subquery
	if p.cur.Type == kwVALUES {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			stmt.Values = p.parseExprList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	} else if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		stmt.Select = p.parseSelectStmt()
	}

	// RETURNING clause
	if p.cur.Type == kwRETURNING {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS clause
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMultiTableInsert parses INSERT ALL or INSERT FIRST multi-table inserts.
func (p *Parser) parseMultiTableInsert(stmt *nodes.InsertStmt, insertType nodes.InsertType) *nodes.InsertStmt {
	stmt.InsertType = insertType
	stmt.MultiTable = &nodes.List{}

	if insertType == nodes.INSERT_FIRST {
		// FIRST: WHEN cond THEN INTO ... [WHEN ...] [ELSE INTO ...] subquery
		for p.cur.Type == kwWHEN {
			p.advance() // consume WHEN
			cond := p.parseExpr()
			if p.cur.Type == kwTHEN {
				p.advance()
			}
			// One or more INTO clauses after THEN
			for p.cur.Type == kwINTO {
				clause := p.parseInsertIntoClause()
				clause.When = cond
				stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
			}
		}
		// ELSE
		if p.cur.Type == kwELSE {
			p.advance()
			for p.cur.Type == kwINTO {
				clause := p.parseInsertIntoClause()
				// When is nil for ELSE clauses
				stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
			}
		}
	} else {
		// ALL: INTO ... INTO ... subquery
		for p.cur.Type == kwINTO {
			clause := p.parseInsertIntoClause()
			stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
		}
	}

	// Trailing subquery
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		stmt.Subquery = p.parseSelectStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseInsertIntoClause parses a single INTO clause for multi-table INSERT.
//
//	INTO table [ ( col1, ... ) ] VALUES ( expr, ... )
func (p *Parser) parseInsertIntoClause() *nodes.InsertIntoClause {
	start := p.pos()
	clause := &nodes.InsertIntoClause{
		Loc: nodes.Loc{Start: start},
	}

	p.advance() // consume INTO

	// Table name
	if p.isIdentLike() {
		clause.Table = p.parseObjectName()
	}

	// Optional column list
	if p.cur.Type == '(' {
		clause.Columns = p.parseParenColumnList()
	}

	// VALUES (...)
	if p.cur.Type == kwVALUES {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			clause.Values = p.parseExprList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	clause.Loc.End = p.pos()
	return clause
}

// parseParenColumnList parses ( col1, col2, ... ) as a list of ColumnRef nodes.
func (p *Parser) parseParenColumnList() *nodes.List {
	p.advance() // consume '('
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		colStart := p.pos()
		name := p.parseIdentifier()
		list.Items = append(list.Items, &nodes.ColumnRef{
			Column: name,
			Loc:    nodes.Loc{Start: colStart, End: p.pos()},
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type == ')' {
		p.advance()
	}
	return list
}

// parseReturningClause parses RETURNING expr [, ...] INTO bind [, ...].
func (p *Parser) parseReturningClause() *nodes.List {
	p.advance() // consume RETURNING
	list := p.parseExprList()
	// INTO bind variables
	if p.cur.Type == kwINTO {
		p.advance()
		binds := p.parseExprList()
		// Append the INTO targets to the returning list
		list.Items = append(list.Items, binds.Items...)
	}
	return list
}

// parseErrorLogClause parses LOG ERRORS [INTO table] [(tag)] [REJECT LIMIT {n|UNLIMITED}].
func (p *Parser) parseErrorLogClause() *nodes.ErrorLogClause {
	start := p.pos()
	p.advance() // consume LOG

	if p.cur.Type == kwERRORS {
		p.advance()
	}

	elc := &nodes.ErrorLogClause{
		Loc: nodes.Loc{Start: start},
	}

	// INTO table
	if p.cur.Type == kwINTO {
		p.advance()
		if p.isIdentLike() {
			elc.Into = p.parseObjectName()
		}
	}

	// ( tag )
	if p.cur.Type == '(' {
		p.advance()
		elc.Tag = p.parseExpr()
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// REJECT LIMIT { n | UNLIMITED }
	if p.cur.Type == kwREJECT {
		p.advance()
		if p.cur.Type == kwLIMIT {
			p.advance()
		}
		elc.Reject = p.parseExpr()
	}

	elc.Loc.End = p.pos()
	return elc
}
