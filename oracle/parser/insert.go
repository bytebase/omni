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
func (p *Parser) parseInsertStmt() (*nodes.InsertStmt, error) {
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
			Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
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
		var parseErr775 error
		stmt.Table, parseErr775 = p.parseObjectName()
		if parseErr775 !=

			// Partition extension clause
			nil {
			return nil, parseErr775
		}
	}

	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		var parseErr776 error
		stmt.PartitionExt, parseErr776 = p.parsePartitionExtClause()
		if parseErr776 !=

			// @dblink
			nil {
			return nil, parseErr776
		}
	}

	if p.cur.Type == '@' {
		p.advance()
		var parseErr777 error
		stmt.Dblink, parseErr777 = p.parseIdentifier()
		if parseErr777 !=

			// Optional alias
			nil {
			return nil, parseErr777
		}
	}

	if p.isTableAliasCandidate() {
		var parseErr778 error
		stmt.Alias, parseErr778 = p.parseAlias()
		if parseErr778 !=

			// Optional column list
			nil {
			return nil, parseErr778
		}
	}

	if p.cur.Type == '(' {
		var parseErr779 error
		stmt.Columns, parseErr779 = p.parseParenColumnList()
		if parseErr779 !=

			// VALUES (...) or subquery or SET or BY NAME/POSITION
			nil {
			return nil, parseErr779
		}
	}

	switch {
	case p.cur.Type == kwVALUES:
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr780 error
			stmt.Values, parseErr780 = p.parseExprList()
			if parseErr780 != nil {
				return nil, parseErr780
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	case p.cur.Type == kwSET:
		// INSERT ... SET col = expr, ...
		p.advance()
		var parseErr781 error
		stmt.SetClauses, parseErr781 = p.parseSetClauses()
		if parseErr781 != nil {
			return nil, parseErr781

			// BY { NAME | POSITION } subquery
		}
	case p.cur.Type == kwBY:

		p.advance()
		if p.cur.Type == kwNAME {
			stmt.ByName = true
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "POSITION" {
			stmt.ByPosition = true
			p.advance()
		}
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
			var parseErr782 error
			stmt.Select, parseErr782 = p.parseSelectStmt()
			if parseErr782 != nil {
				return nil, parseErr782
			}
		}
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(':
		var parseErr783 error
		stmt.Select, parseErr783 = p.parseSelectStmt()
		if parseErr783 !=

			// RETURNING / RETURN clause
			nil {
			return nil, parseErr783
		}
	}

	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		var parseErr784 error
		stmt.Returning, parseErr784 = p.parseReturningClause()
		if parseErr784 !=

			// LOG ERRORS clause
			nil {
			return nil, parseErr784
		}
	}

	if p.cur.Type == kwLOG {
		var parseErr785 error
		stmt.ErrorLog, parseErr785 = p.parseErrorLogClause()
		if parseErr785 != nil {
			return nil, parseErr785
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMultiTableInsert parses INSERT ALL unconditional multi-table inserts.
//
//	ALL into_clause [ values_clause ] [, into_clause [ values_clause ] ]...
//	subquery
func (p *Parser) parseMultiTableInsert(stmt *nodes.InsertStmt, insertType nodes.InsertType) (*nodes.InsertStmt, error) {
	stmt.InsertType = insertType
	stmt.MultiTable = &nodes.List{}

	// ALL: INTO ... INTO ... subquery
	for p.cur.Type == kwINTO {
		clause, parseErr786 := p.parseInsertIntoClause()
		if parseErr786 != nil {
			return nil, parseErr786
		}
		stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
	}

	// Trailing subquery
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
		var parseErr787 error
		stmt.Subquery, parseErr787 = p.parseSelectStmt()
		if parseErr787 != nil {
			return nil, parseErr787
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseConditionalInsert(stmt *nodes.InsertStmt, insertType nodes.InsertType) (*nodes.InsertStmt, error) {
	stmt.InsertType = insertType
	stmt.MultiTable = &nodes.List{}

	for p.cur.Type == kwWHEN {
		p.advance() // consume WHEN
		cond, parseErr788 := p.parseExpr()
		if parseErr788 != nil {
			return nil, parseErr788
		}
		if p.cur.Type == kwTHEN {
			p.advance()
		}
		// One or more INTO clauses after THEN
		for p.cur.Type == kwINTO {
			clause, parseErr789 := p.parseInsertIntoClause()
			if parseErr789 != nil {
				return nil, parseErr789
			}
			clause.When = cond
			stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
		}
	}
	// ELSE
	if p.cur.Type == kwELSE {
		p.advance()
		for p.cur.Type == kwINTO {
			clause, parseErr790 := p.parseInsertIntoClause()
			if parseErr790 !=
				// When is nil for ELSE clauses
				nil {
				return nil, parseErr790
			}

			stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
		}
	}

	// Trailing subquery
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(' {
		var parseErr791 error
		stmt.Subquery, parseErr791 = p.parseSelectStmt()
		if parseErr791 != nil {
			return nil, parseErr791
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseInsertIntoClause parses a single INTO clause for multi-table INSERT.
//
//	INTO [ schema. ] { table | view } [ ( column [, column ]... ) ]
//	[ VALUES ( expr, ... ) ]
func (p *Parser) parseInsertIntoClause() (*nodes.InsertIntoClause, error) {
	start := p.pos()
	clause := &nodes.InsertIntoClause{
		Loc: nodes.Loc{Start: start},
	}

	p.advance() // consume INTO

	// Table name
	if p.isIdentLike() {
		var parseErr792 error
		clause.Table, parseErr792 = p.parseObjectName()
		if parseErr792 !=

			// Optional column list
			nil {
			return nil, parseErr792
		}
	}

	if p.cur.Type == '(' {
		var parseErr793 error
		clause.Columns, parseErr793 = p.parseParenColumnList()
		if parseErr793 !=

			// VALUES (...)
			nil {
			return nil, parseErr793
		}
	}

	if p.cur.Type == kwVALUES {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr794 error
			clause.Values, parseErr794 = p.parseExprList()
			if parseErr794 != nil {
				return nil, parseErr794
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseParenColumnList parses ( col1, col2, ... ) as a list of ColumnRef nodes.
func (p *Parser) parseParenColumnList() (*nodes.List, error) {
	p.advance() // consume '('
	list := &nodes.List{}
	for {
		if !p.isIdentLike() {
			break
		}
		colStart := p.pos()
		name, parseErr795 := p.parseIdentifier()
		if parseErr795 != nil {
			return nil, parseErr795
		}
		list.Items = append(list.Items, &nodes.ColumnRef{
			Column: name,
			Loc:    nodes.Loc{Start: colStart, End: p.prev.End},
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type == ')' {
		p.advance()
	}
	return list, nil
}

// parseReturningClause parses RETURNING/RETURN expr [, ...] INTO bind [, ...].
//
//	returning_clause::=
//	    { RETURNING | RETURN } expr [, expr ]...
//	    INTO data_item [, data_item ]...
func (p *Parser) parseReturningClause() (*nodes.List, error) {
	p.advance() // consume RETURNING or RETURN
	list, parseErr796 := p.parseExprList()
	if parseErr796 !=
		// INTO bind variables
		nil {
		return nil, parseErr796
	}

	if p.cur.Type == kwINTO {
		p.advance()
		binds, parseErr797 := p.parseExprList()
		if parseErr797 !=
			// Append the INTO targets to the returning list
			nil {
			return nil, parseErr797
		}

		list.Items = append(list.Items, binds.Items...)
	}
	return list, nil
}

// parseErrorLogClause parses LOG ERRORS [INTO table] [(tag)] [REJECT LIMIT {n|UNLIMITED}].
//
//	error_logging_clause::=
//	    LOG ERRORS [ INTO [ schema. ] table ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
func (p *Parser) parseErrorLogClause() (*nodes.ErrorLogClause, error) {
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
			var parseErr798 error
			elc.Into, parseErr798 = p.parseObjectName()
			if parseErr798 !=

				// ( tag )
				nil {
				return nil, parseErr798
			}
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr799 error
		elc.Tag, parseErr799 = p.parseExpr()
		if parseErr799 != nil {
			return nil, parseErr799
		}
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
		var parseErr800 error
		elc.Reject, parseErr800 = p.parseExpr()
		if parseErr800 != nil {
			return nil, parseErr800
		}
	}

	elc.Loc.End = p.prev.End
	return elc, nil
}
