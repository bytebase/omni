package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-update.html
//
//	[ WITH [ RECURSIVE ] with_query [, ...] ]
//	UPDATE [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//	    SET { column_name = { expression | DEFAULT } |
//	          ( column_name [, ...] ) = [ ROW ] ( { expression | DEFAULT } [, ...] ) |
//	          ( column_name [, ...] ) = ( sub-SELECT )
//	        } [, ...]
//	    [ FROM from_item [, ...] ]
//	    [ WHERE condition | WHERE CURRENT OF cursor_name ]
//	    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
func (p *Parser) parseUpdateStmt(withClause *nodes.WithClause) *nodes.UpdateStmt {
	loc := p.pos()
	p.advance() // consume UPDATE

	// relation_expr_opt_alias (SET must not be consumed as alias)
	rv := p.parseRelationExpr()
	if rv == nil {
		return nil // collect-mode: parseRelationExpr already emitted rule candidates
	}
	if p.cur.Type == AS {
		p.advance()
		name, _ := p.parseColId()
		rv.Alias = &nodes.Alias{Aliasname: name}
	} else if p.isColId() && p.cur.Type != SET {
		name, _ := p.parseColId()
		rv.Alias = &nodes.Alias{Aliasname: name}
	}

	// SET set_clause_list
	p.expect(SET)
	targetList := p.parseSetClauseList()

	if p.collectMode() {
		// After SET clause, valid continuations:
		for _, t := range []int{',', FROM, WHERE, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// from_clause: FROM from_list | EMPTY
	var fromClause *nodes.List
	if p.cur.Type == FROM {
		p.advance()
		fromClause = p.parseFromListFull()
		if p.collectMode() {
			for _, t := range []int{WHERE, RETURNING, ';'} {
				p.addTokenCandidate(t)
			}
			return nil
		}
	}

	// where_or_current_clause
	whereClause := p.parseWhereOrCurrentClause()
	if p.collectMode() {
		for _, t := range []int{RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// returning_clause
	returningList := p.parseReturningClause()

	stmt := &nodes.UpdateStmt{
		Relation:      rv,
		TargetList:    targetList,
		FromClause:    fromClause,
		WhereClause:   whereClause,
		ReturningList: returningList,
	}
	if withClause != nil {
		stmt.WithClause = withClause
		loc = withClause.Loc.Start
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-delete.html
//
//	[ WITH [ RECURSIVE ] with_query [, ...] ]
//	DELETE FROM [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//	    [ USING from_item [, ...] ]
//	    [ WHERE condition | WHERE CURRENT OF cursor_name ]
//	    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
func (p *Parser) parseDeleteStmt(withClause *nodes.WithClause) *nodes.DeleteStmt {
	loc := p.pos()
	p.advance() // consume DELETE
	p.expect(FROM)

	// relation_expr_opt_alias
	rv := p.parseRelationExprOptAlias()
	if rv == nil {
		return nil // collect-mode: parseRelationExpr already emitted rule candidates
	}

	if p.collectMode() {
		// After relation, valid continuations:
		for _, t := range []int{USING, WHERE, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// using_clause: USING from_list | EMPTY
	var usingClause *nodes.List
	if p.cur.Type == USING {
		p.advance()
		usingClause = p.parseFromListFull()
		if p.collectMode() {
			for _, t := range []int{WHERE, RETURNING, ';'} {
				p.addTokenCandidate(t)
			}
			return nil
		}
	}

	// where_or_current_clause
	whereClause := p.parseWhereOrCurrentClause()
	if p.collectMode() {
		for _, t := range []int{RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// returning_clause
	returningList := p.parseReturningClause()

	stmt := &nodes.DeleteStmt{
		Relation:      rv,
		UsingClause:   usingClause,
		WhereClause:   whereClause,
		ReturningList: returningList,
	}
	if withClause != nil {
		stmt.WithClause = withClause
		loc = withClause.Loc.Start
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt
}

// parseMergeStmt parses a MERGE statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-merge.html
//
//	[ WITH [ RECURSIVE ] with_query [, ...] ]
//	MERGE INTO [ ONLY ] target_table_name [ * ] [ [ AS ] target_alias ]
//	USING data_source ON join_condition
//	when_clause [...]
//	[ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
func (p *Parser) parseMergeStmt(withClause *nodes.WithClause) *nodes.MergeStmt {
	loc := p.pos()
	p.advance() // consume MERGE
	p.expect(INTO)

	// relation_expr_opt_alias
	rv := p.parseRelationExprOptAlias()
	if rv == nil {
		return nil // collect-mode: parseRelationExpr already emitted rule candidates
	}

	if p.collectMode() {
		p.addTokenCandidate(USING)
		return nil
	}

	// USING table_ref
	p.expect(USING)
	sourceRelation := p.parseTableRef()

	if p.collectMode() {
		p.addTokenCandidate(ON)
		return nil
	}

	// ON a_expr
	p.expect(ON)
	joinCondition, _ := p.parseAExpr(0)

	if p.collectMode() {
		p.addTokenCandidate(WHEN)
		p.addTokenCandidate(RETURNING)
		return nil
	}

	// merge_when_list
	mergeWhenClauses := p.parseMergeWhenList()

	// returning_clause
	returningList := p.parseReturningClause()

	stmt := &nodes.MergeStmt{
		Relation:         rv,
		SourceRelation:   sourceRelation,
		JoinCondition:    joinCondition,
		MergeWhenClauses: mergeWhenClauses,
		ReturningList:    returningList,
	}
	if withClause != nil {
		stmt.WithClause = withClause
		loc = withClause.Loc.Start
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt
}

// parseMergeWhenList parses one or more WHEN clauses in a MERGE statement.
//
//	merge_when_list:
//	    merge_when_clause [merge_when_clause ...]
func (p *Parser) parseMergeWhenList() *nodes.List {
	var items []nodes.Node
	for p.cur.Type == WHEN {
		items = append(items, p.parseMergeWhenClause())
	}
	return &nodes.List{Items: items}
}

// parseMergeWhenClause parses a single WHEN clause in a MERGE statement.
//
//	merge_when_clause:
//	    merge_when_tgt_matched opt_merge_when_condition THEN merge_update
//	    | merge_when_tgt_matched opt_merge_when_condition THEN merge_delete
//	    | merge_when_tgt_not_matched opt_merge_when_condition THEN merge_insert
//	    | merge_when_tgt_matched opt_merge_when_condition THEN DO NOTHING
//	    | merge_when_tgt_not_matched opt_merge_when_condition THEN DO NOTHING
func (p *Parser) parseMergeWhenClause() *nodes.MergeWhenClause {
	p.advance() // consume WHEN

	// Determine match kind
	var kind nodes.MergeMatchKind
	if p.cur.Type == MATCHED {
		// WHEN MATCHED
		p.advance()
		kind = nodes.MERGE_WHEN_MATCHED
	} else if p.cur.Type == NOT {
		p.advance() // consume NOT
		p.expect(MATCHED)
		// Check for BY SOURCE or BY TARGET
		if p.cur.Type == BY {
			p.advance() // consume BY
			if p.cur.Type == SOURCE {
				p.advance()
				kind = nodes.MERGE_WHEN_NOT_MATCHED_BY_SOURCE
			} else {
				// TARGET
				p.advance()
				kind = nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET
			}
		} else {
			// WHEN NOT MATCHED (no BY clause) = BY TARGET
			kind = nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET
		}
	}

	// opt_merge_when_condition: AND a_expr | EMPTY
	var condition nodes.Node
	if p.cur.Type == AND {
		p.advance()
		condition, _ = p.parseAExpr(0)
	}

	// THEN
	p.expect(THEN)

	// Dispatch based on action keyword
	if p.collectMode() {
		for _, t := range []int{UPDATE, DELETE_P, INSERT, DO} {
			p.addTokenCandidate(t)
		}
		return nil
	}
	switch p.cur.Type {
	case UPDATE:
		// merge_update: UPDATE SET set_clause_list
		p.advance()
		p.expect(SET)
		targetList := p.parseSetClauseList()
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_UPDATE,
			Override:    nodes.OVERRIDING_NOT_SET,
			TargetList:  targetList,
		}

	case DELETE_P:
		// merge_delete: DELETE
		p.advance()
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_DELETE,
			Override:    nodes.OVERRIDING_NOT_SET,
		}

	case INSERT:
		// merge_insert
		return p.parseMergeInsert(kind, condition)

	case DO:
		// DO NOTHING
		p.advance()
		p.expect(NOTHING)
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_NOTHING,
		}

	default:
		return &nodes.MergeWhenClause{
			Kind:      kind,
			Condition: condition,
		}
	}
}

// parseMergeInsert parses the INSERT action within a MERGE WHEN clause.
//
//	merge_insert:
//	    INSERT merge_values_clause
//	    | INSERT OVERRIDING override_kind VALUE merge_values_clause
//	    | INSERT '(' insert_column_list ')' merge_values_clause
//	    | INSERT '(' insert_column_list ')' OVERRIDING override_kind VALUE merge_values_clause
//	    | INSERT DEFAULT VALUES
func (p *Parser) parseMergeInsert(kind nodes.MergeMatchKind, condition nodes.Node) *nodes.MergeWhenClause {
	p.advance() // consume INSERT

	clause := &nodes.MergeWhenClause{
		Kind:        kind,
		Condition:   condition,
		CommandType: nodes.CMD_INSERT,
		Override:    nodes.OVERRIDING_NOT_SET,
	}

	switch p.cur.Type {
	case DEFAULT:
		// INSERT DEFAULT VALUES
		p.advance()
		p.expect(VALUES)
		return clause

	case '(':
		// INSERT '(' insert_column_list ')' [OVERRIDING ...] merge_values_clause
		p.advance() // consume '('
		clause.TargetList = p.parseInsertColumnList()
		p.expect(')')
		if p.cur.Type == OVERRIDING {
			clause.Override = p.parseOverriding()
		}
		clause.Values = p.parseMergeValuesClause()
		return clause

	case OVERRIDING:
		// INSERT OVERRIDING override_kind VALUE merge_values_clause
		clause.Override = p.parseOverriding()
		clause.Values = p.parseMergeValuesClause()
		return clause

	default:
		// INSERT merge_values_clause (VALUES (...))
		clause.Values = p.parseMergeValuesClause()
		return clause
	}
}

// parseMergeValuesClause parses VALUES '(' expr_list ')' in a MERGE INSERT.
//
//	merge_values_clause:
//	    VALUES '(' expr_list ')'
func (p *Parser) parseMergeValuesClause() *nodes.List {
	p.expect(VALUES)
	p.expect('(')
	list := p.parseExprList()
	p.expect(')')
	return list
}
