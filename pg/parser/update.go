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
func (p *Parser) parseUpdateStmt(withClause *nodes.WithClause) (*nodes.UpdateStmt, error) {
	loc := p.pos()
	p.advance() // consume UPDATE

	// relation_expr_opt_alias (SET must not be consumed as alias)
	rv, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	if rv == nil {
		return nil, nil // collect-mode: parseRelationExpr already emitted rule candidates
	}
	if p.cur.Type == AS {
		aliasLoc := p.pos()
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		rv.Alias = &nodes.Alias{Aliasname: name, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
	} else if p.isColId() && p.cur.Type != SET {
		aliasLoc := p.pos()
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		rv.Alias = &nodes.Alias{Aliasname: name, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
	}

	// SET set_clause_list
	if _, err := p.expect(SET); err != nil {
		return nil, err
	}
	targetList, err := p.parseSetClauseList()
	if err != nil {
		return nil, err
	}

	if p.collectMode() {
		// After SET clause, valid continuations:
		for _, t := range []int{',', FROM, WHERE, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// from_clause: FROM from_list | EMPTY
	var fromClause *nodes.List
	if p.cur.Type == FROM {
		p.advance()
		fromClause, err = p.parseFromListFull()
		if err != nil {
			return nil, err
		}
		if fromClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.collectMode() {
			for _, t := range []int{WHERE, RETURNING, ';'} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
	}

	// where_or_current_clause
	whereClause, err := p.parseWhereOrCurrentClause()
	if err != nil {
		return nil, err
	}
	if p.collectMode() {
		for _, t := range []int{RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// returning_clause
	returningList, err := p.parseReturningClause()
	if err != nil {
		return nil, err
	}

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
	return stmt, nil
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
func (p *Parser) parseDeleteStmt(withClause *nodes.WithClause) (*nodes.DeleteStmt, error) {
	loc := p.pos()
	p.advance() // consume DELETE
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}

	// relation_expr_opt_alias
	rv, err := p.parseRelationExprOptAlias()
	if err != nil {
		return nil, err
	}
	if rv == nil {
		return nil, nil // collect-mode: parseRelationExpr already emitted rule candidates
	}

	if p.collectMode() {
		// After relation, valid continuations:
		for _, t := range []int{USING, WHERE, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// using_clause: USING from_list | EMPTY
	var usingClause *nodes.List
	if p.cur.Type == USING {
		p.advance()
		usingClause, err = p.parseFromListFull()
		if err != nil {
			return nil, err
		}
		if usingClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.collectMode() {
			for _, t := range []int{WHERE, RETURNING, ';'} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
	}

	// where_or_current_clause
	whereClause, err := p.parseWhereOrCurrentClause()
	if err != nil {
		return nil, err
	}
	if p.collectMode() {
		for _, t := range []int{RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// returning_clause
	returningList, err := p.parseReturningClause()
	if err != nil {
		return nil, err
	}

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
	return stmt, nil
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
func (p *Parser) parseMergeStmt(withClause *nodes.WithClause) (*nodes.MergeStmt, error) {
	loc := p.pos()
	p.advance() // consume MERGE
	if _, err := p.expect(INTO); err != nil {
		return nil, err
	}

	// relation_expr_opt_alias
	rv, err := p.parseRelationExprOptAlias()
	if err != nil {
		return nil, err
	}
	if rv == nil {
		return nil, nil // collect-mode: parseRelationExpr already emitted rule candidates
	}

	if p.collectMode() {
		p.addTokenCandidate(USING)
		return nil, errCollecting
	}

	// USING table_ref
	if _, err := p.expect(USING); err != nil {
		return nil, err
	}
	sourceRelation, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	if p.collectMode() {
		p.addTokenCandidate(ON)
		return nil, errCollecting
	}

	// ON a_expr
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	joinCondition, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.collectMode() {
		p.addTokenCandidate(WHEN)
		p.addTokenCandidate(RETURNING)
		return nil, errCollecting
	}

	// merge_when_list
	mergeWhenClauses, err := p.parseMergeWhenList()
	if err != nil {
		return nil, err
	}

	// returning_clause
	returningList, err := p.parseReturningClause()
	if err != nil {
		return nil, err
	}

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
	return stmt, nil
}

// parseMergeWhenList parses one or more WHEN clauses in a MERGE statement.
//
//	merge_when_list:
//	    merge_when_clause [merge_when_clause ...]
func (p *Parser) parseMergeWhenList() (*nodes.List, error) {
	var items []nodes.Node
	for p.cur.Type == WHEN {
		clause, err := p.parseMergeWhenClause()
		if err != nil {
			return nil, err
		}
		items = append(items, clause)
	}
	return &nodes.List{Items: items}, nil
}

// parseMergeWhenClause parses a single WHEN clause in a MERGE statement.
//
//	merge_when_clause:
//	    merge_when_tgt_matched opt_merge_when_condition THEN merge_update
//	    | merge_when_tgt_matched opt_merge_when_condition THEN merge_delete
//	    | merge_when_tgt_not_matched opt_merge_when_condition THEN merge_insert
//	    | merge_when_tgt_matched opt_merge_when_condition THEN DO NOTHING
//	    | merge_when_tgt_not_matched opt_merge_when_condition THEN DO NOTHING
func (p *Parser) parseMergeWhenClause() (*nodes.MergeWhenClause, error) {
	p.advance() // consume WHEN

	// Determine match kind
	var kind nodes.MergeMatchKind
	if p.cur.Type == MATCHED {
		// WHEN MATCHED
		p.advance()
		kind = nodes.MERGE_WHEN_MATCHED
	} else if p.cur.Type == NOT {
		p.advance() // consume NOT
		if _, err := p.expect(MATCHED); err != nil {
			return nil, err
		}
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
		var err error
		condition, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
	}

	// THEN
	if _, err := p.expect(THEN); err != nil {
		return nil, err
	}

	// Dispatch based on action keyword
	if p.collectMode() {
		for _, t := range []int{UPDATE, DELETE_P, INSERT, DO} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	switch p.cur.Type {
	case UPDATE:
		// merge_update: UPDATE SET set_clause_list
		p.advance()
		if _, err := p.expect(SET); err != nil {
			return nil, err
		}
		targetList, err := p.parseSetClauseList()
		if err != nil {
			return nil, err
		}
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_UPDATE,
			Override:    nodes.OVERRIDING_NOT_SET,
			TargetList:  targetList,
		}, nil

	case DELETE_P:
		// merge_delete: DELETE
		p.advance()
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_DELETE,
			Override:    nodes.OVERRIDING_NOT_SET,
		}, nil

	case INSERT:
		// merge_insert
		return p.parseMergeInsert(kind, condition)

	case DO:
		// DO NOTHING
		p.advance()
		if _, err := p.expect(NOTHING); err != nil {
			return nil, err
		}
		return &nodes.MergeWhenClause{
			Kind:        kind,
			Condition:   condition,
			CommandType: nodes.CMD_NOTHING,
		}, nil

	default:
		return &nodes.MergeWhenClause{
			Kind:      kind,
			Condition: condition,
		}, nil
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
func (p *Parser) parseMergeInsert(kind nodes.MergeMatchKind, condition nodes.Node) (*nodes.MergeWhenClause, error) {
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
		if _, err := p.expect(VALUES); err != nil {
			return nil, err
		}
		return clause, nil

	case '(':
		// INSERT '(' insert_column_list ')' [OVERRIDING ...] merge_values_clause
		p.advance() // consume '('
		var err error
		clause.TargetList, err = p.parseInsertColumnList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if p.cur.Type == OVERRIDING {
			clause.Override, err = p.parseOverriding()
			if err != nil {
				return nil, err
			}
		}
		clause.Values, err = p.parseMergeValuesClause()
		if err != nil {
			return nil, err
		}
		return clause, nil

	case OVERRIDING:
		// INSERT OVERRIDING override_kind VALUE merge_values_clause
		var err error
		clause.Override, err = p.parseOverriding()
		if err != nil {
			return nil, err
		}
		clause.Values, err = p.parseMergeValuesClause()
		if err != nil {
			return nil, err
		}
		return clause, nil

	default:
		// INSERT merge_values_clause (VALUES (...))
		var err error
		clause.Values, err = p.parseMergeValuesClause()
		if err != nil {
			return nil, err
		}
		return clause, nil
	}
}

// parseMergeValuesClause parses VALUES '(' expr_list ')' in a MERGE INSERT.
//
//	merge_values_clause:
//	    VALUES '(' expr_list ')'
func (p *Parser) parseMergeValuesClause() (*nodes.List, error) {
	if _, err := p.expect(VALUES); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list := p.parseExprList()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}
