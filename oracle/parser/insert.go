package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseInsertStmt parses an INSERT statement.
//
// BNF: oracle/parser/bnf/INSERT.bnf
//
//	INSERT [ hint ]
//	    { single_table_insert | multi_table_insert }
//
//	single_table_insert::=
//	    insert_into_clause
//	    { values_clause | subquery | insert_set_clause | by_name_position_clause }
//	    [ returning_clause ]
//	    [ error_logging_clause ]
//
//	insert_into_clause::=
//	    INTO [ schema. ] { table | view | materialized_view }
//	    [ partition_extension_clause ]
//	    [ @ dblink ]
//	    [ t_alias ]
//	    [ ( column [, column ]... ) ]
//
//	DML_table_expression_clause::=
//	    [ schema. ] { table | view | materialized_view }
//	    [ partition_extension_clause ]
//	    [ @ dblink ]
//	    [ t_alias ]
//	  | ( subquery [ subquery_restriction_clause ] ) [ t_alias ]
//	  | TABLE ( collection_expression ) [ (+) ] [ t_alias ]
//
//	partition_extension_clause::=
//	    PARTITION ( partition )
//	  | PARTITION FOR ( partition_key_value [, partition_key_value ]... )
//	  | SUBPARTITION ( subpartition )
//	  | SUBPARTITION FOR ( subpartition_key_value [, subpartition_key_value ]... )
//
//	subquery_restriction_clause::=
//	    WITH { READ ONLY | CHECK OPTION [ CONSTRAINT constraint ] }
//
//	values_clause::=
//	    VALUES ( expr [, expr ]... )
//
//	insert_set_clause::=
//	    SET column = { expr | DEFAULT } [, column = { expr | DEFAULT } ]...
//
//	by_name_position_clause::=
//	    BY { NAME | POSITION } subquery
//
//	returning_clause::=
//	    { RETURNING | RETURN } expr [, expr ]...
//	    INTO data_item [, data_item ]...
//
//	multi_table_insert::=
//	    { ALL into_clause [ values_clause ] [, into_clause [ values_clause ] ]...
//	    | conditional_insert_clause
//	    }
//	    subquery
//
//	conditional_insert_clause::=
//	    [ ALL | FIRST ]
//	    WHEN condition THEN
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]...
//	    [ WHEN condition THEN
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]... ]...
//	    [ ELSE
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]... ]
//
//	into_clause::=
//	    INTO [ schema. ] { table | view } [ ( column [, column ]... ) ]
//
//	error_logging_clause::=
//	    LOG ERRORS [ INTO [ schema. ] table ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
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

	// INSERT ALL ... (unconditional multi-table)
	if p.cur.Type == kwALL {
		next := p.peekNext()
		if next.Type == kwINTO {
			// INSERT ALL INTO ... INTO ... subquery
			p.advance()
			return p.parseMultiTableInsert(stmt, nodes.INSERT_ALL)
		}
		// INSERT ALL WHEN ... (conditional with ALL prefix)
		if next.Type == kwWHEN {
			p.advance()
			return p.parseConditionalInsert(stmt, nodes.INSERT_ALL)
		}
		// Fallback: treat as unconditional
		p.advance()
		return p.parseMultiTableInsert(stmt, nodes.INSERT_ALL)
	}

	// INSERT FIRST WHEN ...
	if p.cur.Type == kwFIRST {
		p.advance()
		return p.parseConditionalInsert(stmt, nodes.INSERT_FIRST)
	}

	// Conditional insert without ALL/FIRST prefix: INSERT WHEN ...
	if p.cur.Type == kwWHEN {
		return p.parseConditionalInsert(stmt, nodes.INSERT_ALL)
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

	// Partition extension clause
	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		stmt.PartitionExt = p.parsePartitionExtClause()
	}

	// @dblink
	if p.cur.Type == '@' {
		p.advance()
		stmt.Dblink = p.parseIdentifier()
	}

	// Optional alias
	if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// Optional column list
	if p.cur.Type == '(' {
		stmt.Columns = p.parseParenColumnList()
	}

	// VALUES (...) or subquery or SET or BY NAME/POSITION
	switch {
	case p.cur.Type == kwVALUES:
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			stmt.Values = p.parseExprList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	case p.cur.Type == kwSET:
		// INSERT ... SET col = expr, ...
		p.advance()
		stmt.SetClauses = p.parseSetClauses()
	case p.cur.Type == kwBY:
		// BY { NAME | POSITION } subquery
		p.advance()
		if p.cur.Type == kwNAME {
			stmt.ByName = true
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "POSITION" {
			stmt.ByPosition = true
			p.advance()
		}
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
			stmt.Select = p.parseSelectStmt()
		}
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(':
		stmt.Select = p.parseSelectStmt()
	}

	// RETURNING / RETURN clause
	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS clause
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseMultiTableInsert parses INSERT ALL unconditional multi-table inserts.
//
//	ALL into_clause [ values_clause ] [, into_clause [ values_clause ] ]...
//	subquery
func (p *Parser) parseMultiTableInsert(stmt *nodes.InsertStmt, insertType nodes.InsertType) *nodes.InsertStmt {
	stmt.InsertType = insertType
	stmt.MultiTable = &nodes.List{}

	// ALL: INTO ... INTO ... subquery
	for p.cur.Type == kwINTO {
		clause := p.parseInsertIntoClause()
		stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
	}

	// Trailing subquery
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
		stmt.Subquery = p.parseSelectStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseConditionalInsert parses conditional INSERT (WHEN ... THEN INTO ...).
//
//	conditional_insert_clause::=
//	    [ ALL | FIRST ]
//	    WHEN condition THEN
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]...
//	    [ WHEN condition THEN
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]... ]...
//	    [ ELSE
//	        into_clause [ values_clause ]
//	        [, into_clause [ values_clause ] ]... ]
//	    subquery
func (p *Parser) parseConditionalInsert(stmt *nodes.InsertStmt, insertType nodes.InsertType) *nodes.InsertStmt {
	stmt.InsertType = insertType
	stmt.MultiTable = &nodes.List{}

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

	// Trailing subquery
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
		stmt.Subquery = p.parseSelectStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseInsertIntoClause parses a single INTO clause for multi-table INSERT.
//
//	INTO [ schema. ] { table | view } [ ( column [, column ]... ) ]
//	[ VALUES ( expr, ... ) ]
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

// parseReturningClause parses RETURNING/RETURN expr [, ...] INTO bind [, ...].
//
//	returning_clause::=
//	    { RETURNING | RETURN } expr [, expr ]...
//	    INTO data_item [, data_item ]...
func (p *Parser) parseReturningClause() *nodes.List {
	p.advance() // consume RETURNING or RETURN
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
//
//	error_logging_clause::=
//	    LOG ERRORS [ INTO [ schema. ] table ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
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
