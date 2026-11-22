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
	defer p.restorePLSQLCode(p.setPLSQLCode(false))

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

	// INTO (ORA-00925)
	if p.cur.Type != kwINTO {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// Table name
	if !p.isIdentLike() && p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	if p.isIdentLike() {
		var parseErr775 error
		stmt.Table, parseErr775 = p.parseReservedCheckedObjectName()
		if parseErr775 !=

			// Partition extension clause
			nil {
			return nil, parseErr775
		}
	}

	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		var parseErr776 error
		stmt.PartitionExt, parseErr776 = p.parsePartitionExtClause(false)
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

	// A column list, unless the parentheses hold the subquery:
	// INSERT INTO t (SELECT ...) compiles on Oracle 23ai.
	if p.cur.Type == '(' && !p.atQuery() {
		var parseErr779 error
		stmt.Columns, parseErr779 = p.parseParenColumnList()
		if parseErr779 !=

			// VALUES (...) or subquery or SET or BY NAME/POSITION
			nil {
			return nil, parseErr779
		}
	}

	// A returning clause follows the VALUES and SET forms alone: after a
	// subquery Oracle 23ai reads RETURNING as the query's last table alias
	// and fails at what follows it (ORA-03049).
	returns := false
	switch {
	case p.cur.Type == kwVALUES:
		returns = true
		p.advance()
		if p.cur.Type != '(' {
			// INSERT ... VALUES record — the PL/SQL record form (engine-
			// verified, incl. FORALL collection elements: VALUES arr(i)).
			var parseErrRec error
			stmt.ValuesRecord, parseErrRec = p.parseExpr()
			if parseErrRec != nil {
				return nil, parseErrRec
			}
			if stmt.ValuesRecord == nil {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		p.advance()
		var parseErr780 error
		stmt.Values, parseErr780 = p.parseExprList()
		if parseErr780 != nil {
			return nil, parseErr780
		}
		if stmt.Values == nil || stmt.Values.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	case p.cur.Type == kwSET:
		// INSERT ... SET col = expr, ...
		returns = true
		p.advance()
		var parseErr781 error
		stmt.SetClauses, parseErr781 = p.parseSetClausesUnique(true)
		if parseErr781 != nil {
			return nil, parseErr781

			// BY { NAME | POSITION } subquery
		}
	case p.cur.Type == kwBY:

		p.advance()
		switch {
		case p.isKeywordStr("NAME"):
			stmt.ByName = true
		case p.isKeywordStr("POSITION"):
			stmt.ByPosition = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if !p.atQuery() {
			return nil, p.syntaxErrorAtCur()
		}
		var parseErr782 error
		stmt.Select, parseErr782 = p.parseSelectStmt()
		if parseErr782 != nil {
			return nil, parseErr782
		}
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH || p.cur.Type == '(':
		var parseErr783 error
		stmt.Select, parseErr783 = p.parseSelectStmt()
		if parseErr783 != nil {
			return nil, parseErr783
		}
	default:
		// The rows to insert: VALUES, a subquery, SET or BY (ORA-00926).
		return nil, p.syntaxErrorAtCur()
	}

	if returns && (p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN) {
		var parseErr784 error
		stmt.Returning, stmt.ReturningBulk, parseErr784 = p.parseReturningClause()
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
	} else {
		return nil, p.syntaxErrorAtCur()
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
		if p.cur.Type == kwTHEN {
			return nil, p.syntaxErrorAtCur()
		}
		cond, parseErr788 := p.parseExpr()
		if parseErr788 != nil {
			return nil, parseErr788
		}
		if cond == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != kwTHEN {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		// One or more INTO clauses after THEN
		intoCount := 0
		for p.cur.Type == kwINTO {
			clause, parseErr789 := p.parseInsertIntoClause()
			if parseErr789 != nil {
				return nil, parseErr789
			}
			clause.When = cond
			stmt.MultiTable.Items = append(stmt.MultiTable.Items, clause)
			intoCount++
		}
		if intoCount == 0 {
			return nil, p.syntaxErrorAtCur()
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
	} else {
		return nil, p.syntaxErrorAtCur()
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
		clause.Table, parseErr792 = p.parseReservedCheckedObjectName()
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
		if p.cur.Type != '(' {
			// Oracle's SQL parser shares the VALUES record rule before rejecting
			// non-parenthesized VALUES on multi-table INSERT, so VALUES 1)
			// reports the ')' after the expression rather than 1.
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			_ = expr
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr794 error
		clause.Values, parseErr794 = p.parseExprList()
		if parseErr794 != nil {
			return nil, parseErr794
		}
		if clause.Values == nil || clause.Values.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
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
			return nil, p.syntaxErrorAtCur()
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
	if list.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseReturningClause parses RETURNING/RETURN expr [, ...] INTO bind [, ...].
//
//	returning_clause::=
//	    { RETURNING | RETURN } expr [, expr ]...
//	    INTO data_item [, data_item ]...
//
// UPDATE's returning clause, and on Oracle 23ai DELETE's and INSERT's too,
// takes OLD column and NEW column as items: the value before or after the
// change. OLD and NEW then name a column (ORA-03050 at INTO in RETURNING
// OLD INTO, ORA-00921 at the end), and the expression goes on after it, as
// in RETURNING OLD a + 1. The AST keeps the expression; OLD.a is a column of
// a table named OLD.
func (p *Parser) parseReturningClause() (*nodes.List, bool, error) {
	p.advance() // consume RETURNING or RETURN
	list := &nodes.List{}
	for {
		if p.isReturningOldNewStart() {
			p.advance() // consume OLD or NEW
			if !p.isIdentLike() {
				return nil, false, p.syntaxErrorAtCur()
			}
			if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
				return nil, false, err
			}
		}
		if p.cur.Type == kwINTO {
			return nil, false, p.syntaxErrorAtCur()
		}
		item, err := p.parseExpr()
		if err != nil {
			return nil, false, err
		}
		if item == nil {
			return nil, false, p.syntaxErrorAtCur()
		}
		list.Items = append(list.Items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// { INTO | BULK COLLECT INTO } targets, appended to the returned list;
	// bulk reports BULK COLLECT.
	bulk, err := p.parseOptionalBulkCollect()
	if err != nil {
		return nil, false, err
	}
	if p.cur.Type == kwINTO {
		p.advance()
		binds, parseErr797 := p.parseExprList()
		if parseErr797 !=
			// Append the INTO targets to the returning list
			nil {
			return nil, false, parseErr797
		}
		if err := p.checkTargets(binds); err != nil {
			return nil, false, err
		}

		list.Items = append(list.Items, binds.Items...)
	}
	return list, bulk, nil
}

// isReturningOldNewStart reports whether the current token is the OLD or
// NEW of a returning item, not a column or table of that name: OLD.a is a
// column of a table named OLD.
func (p *Parser) isReturningOldNewStart() bool {
	return (p.isKeywordStr("OLD") || p.isKeywordStr("NEW")) && p.peekNext().Type != '.'
}

// parseErrorLogClause parses LOG ERRORS [INTO table] [(tag)] [REJECT LIMIT {n|UNLIMITED}].
//
//	error_logging_clause::=
//	    LOG ERRORS [ INTO [ schema. ] table ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
//
// LOG takes ERRORS (ORA-02000), INTO a table (ORA-00903), the parentheses
// an expression (ORA-00936), and REJECT LIMIT (ORA-30652) an integer or
// UNLIMITED: Oracle 23ai refuses 1+1, -1, 1.5, a bind or a name there
// (ORA-30645, ORA-02017, ORA-03048).
func (p *Parser) parseErrorLogClause() (*nodes.ErrorLogClause, error) {
	start := p.pos()
	p.advance() // consume LOG

	if p.cur.Type != kwERRORS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	elc := &nodes.ErrorLogClause{
		Loc: nodes.Loc{Start: start},
	}

	var err error
	if p.cur.Type == kwINTO {
		p.advance()
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		if elc.Into, err = p.parseObjectName(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		if elc.Tag, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if elc.Tag == nil || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	if p.cur.Type == kwREJECT {
		p.advance()
		if p.cur.Type != kwLIMIT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != tokICONST && !p.isKeywordStr("UNLIMITED") {
			return nil, p.syntaxErrorAtCur()
		}
		if elc.Reject, err = p.parsePrimary(); err != nil {
			return nil, err
		}
	}

	elc.Loc.End = p.prev.End
	return elc, nil
}
