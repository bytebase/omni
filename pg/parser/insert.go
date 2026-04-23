package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ON CONFLICT action constants matching PostgreSQL's values.
const (
	ONCONFLICT_NONE    = 0
	ONCONFLICT_NOTHING = 1
	ONCONFLICT_UPDATE  = 2
)

// parseInsertStmt parses an INSERT statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-insert.html
//
//	[WITH ...] INSERT INTO insert_target insert_rest [ON CONFLICT ...] [RETURNING ...]
//
//	insert_target:
//	    qualified_name
//	    | qualified_name AS ColId
//
//	insert_rest:
//	    SelectStmt
//	    | OVERRIDING override_kind VALUE SelectStmt
//	    | '(' insert_column_list ')' SelectStmt
//	    | '(' insert_column_list ')' OVERRIDING override_kind VALUE SelectStmt
//	    | DEFAULT VALUES
func (p *Parser) parseInsertStmt(withClause *nodes.WithClause) (*nodes.InsertStmt, error) {
	loc := p.pos()
	p.advance() // consume INSERT
	if _, err := p.expect(INTO); err != nil {
		return nil, err
	}

	// insert_target: qualified_name [AS ColId]
	rvLoc := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if names == nil {
		return nil, nil // collect-mode: parseQualifiedName already emitted rule candidates
	}
	rv := makeRangeVarFromNames(names)
	if p.cur.Type == AS {
		aliasLoc := p.pos()
		p.advance()
		alias, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		rv.Alias = &nodes.Alias{Aliasname: alias, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
	}
	rv.Loc = nodes.Loc{Start: rvLoc, End: p.pos()}

	if p.collectMode() {
		// After table name, valid continuations for insert_rest:
		for _, t := range []int{'(', DEFAULT, OVERRIDING, SELECT, VALUES, TABLE, WITH} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// insert_rest
	stmt, err := p.parseInsertRest()
	if err != nil {
		return nil, err
	}
	if stmt == nil {
		return nil, nil
	}
	stmt.Relation = rv

	if p.collectMode() {
		// After insert_rest, valid continuations:
		for _, t := range []int{ON, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// opt_on_conflict
	if p.cur.Type == ON {
		next := p.peekNext()
		if next.Type == CONFLICT {
			var err error
			stmt.OnConflictClause, err = p.parseOnConflict()
			if err != nil {
				return nil, err
			}
		}
	}

	// returning_clause
	var retErr error
	stmt.ReturningList, retErr = p.parseReturningClause()
	if retErr != nil {
		return nil, retErr
	}

	if withClause != nil {
		stmt.WithClause = withClause
		loc = withClause.Loc.Start
	}

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseInsertRest parses the body of an INSERT statement after the target table.
func (p *Parser) parseInsertRest() (*nodes.InsertStmt, error) {
	if p.collectMode() {
		// Valid first tokens for insert_rest:
		for _, t := range []int{DEFAULT, '(', OVERRIDING, SELECT, VALUES, TABLE, WITH} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	switch p.cur.Type {
	case DEFAULT:
		// DEFAULT VALUES
		p.advance()
		if _, err := p.expect(VALUES); err != nil {
			return nil, err
		}
		return &nodes.InsertStmt{}, nil

	case '(':
		if next := p.peekNext(); next.Type == '(' || isSelectStartToken(next.Type) {
			selectStmt, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			return &nodes.InsertStmt{
				SelectStmt: selectStmt,
			}, nil
		}

		// '(' insert_column_list ')' [OVERRIDING override_kind VALUE] SelectStmt
		p.advance() // consume '('
		cols, err := p.parseInsertColumnList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		var override nodes.OverridingKind
		if p.cur.Type == OVERRIDING {
			var err error
			override, err = p.parseOverriding()
			if err != nil {
				return nil, err
			}
		}

		selectStmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.InsertStmt{
			Cols:       cols,
			Override:   override,
			SelectStmt: selectStmt,
		}, nil

	case OVERRIDING:
		// OVERRIDING override_kind VALUE SelectStmt
		override, err := p.parseOverriding()
		if err != nil {
			return nil, err
		}
		selectStmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.InsertStmt{
			Override:   override,
			SelectStmt: selectStmt,
		}, nil

	default:
		// SelectStmt
		selectStmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.InsertStmt{
			SelectStmt: selectStmt,
		}, nil
	}
}

// parseOverriding parses OVERRIDING {USER|SYSTEM} VALUE.
func (p *Parser) parseOverriding() (nodes.OverridingKind, error) {
	p.advance() // consume OVERRIDING
	var kind nodes.OverridingKind
	switch p.cur.Type {
	case USER:
		p.advance()
		kind = nodes.OVERRIDING_USER_VALUE
	case SYSTEM_P:
		p.advance()
		kind = nodes.OVERRIDING_SYSTEM_VALUE
	}
	if _, err := p.expect(VALUE_P); err != nil {
		return kind, err
	}
	return kind, nil
}

// parseSelectStmt parses SelectStmt: select_no_parens | select_with_parens.
func (p *Parser) parseSelectStmt() (nodes.Node, error) {
	if p.cur.Type == '(' {
		return p.parseSelectWithParens()
	}
	return p.parseSelectNoParens()
}

// parseInsertColumnList parses a comma-separated list of insert columns.
//
//	insert_column_list:
//	    insert_column_item [',' insert_column_item ...]
func (p *Parser) parseInsertColumnList() (*nodes.List, error) {
	first, err := p.parseInsertColumnItem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		item, err := p.parseInsertColumnItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

// parseInsertColumnItem parses a single insert column.
//
//	insert_column_item:
//	    ColId opt_indirection
func (p *Parser) parseInsertColumnItem() (*nodes.ResTarget, error) {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	ind, err := p.parseOptIndirection()
	if err != nil {
		return nil, err
	}
	return &nodes.ResTarget{
		Name:        name,
		Indirection: ind,
		Loc:         nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseOnConflict parses ON CONFLICT clause.
//
//	opt_on_conflict:
//	    ON CONFLICT DO NOTHING
//	    | ON CONFLICT DO UPDATE SET set_clause_list where_clause
//	    | ON CONFLICT '(' index_params ')' [WHERE a_expr] DO {NOTHING | UPDATE SET ...}
//	    | ON CONFLICT ON CONSTRAINT name DO {NOTHING | UPDATE SET ...}
func (p *Parser) parseOnConflict() (*nodes.OnConflictClause, error) {
	loc := p.pos()
	p.advance() // consume ON
	p.advance() // consume CONFLICT

	occ := &nodes.OnConflictClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Determine which form we have
	switch p.cur.Type {
	case '(':
		// '(' index_params ')' [WHERE a_expr] DO ...
		inferLoc := p.pos()
		p.advance() // consume '('
		indexElems := p.parseIndexParams()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		infer := &nodes.InferClause{
			IndexElems: indexElems,
			Loc:        nodes.Loc{Start: inferLoc, End: -1},
		}

		// Optional WHERE clause for partial index predicate
		if p.cur.Type == WHERE {
			p.advance()
			var err error
			infer.WhereClause, err = p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
		}

		infer.Loc.End = p.pos()
		occ.Infer = infer

	case ON:
		// ON CONSTRAINT name
		inferLoc := p.pos()
		p.advance() // consume ON (second one)
		if _, err := p.expect(CONSTRAINT); err != nil {
			return nil, err
		}
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		occ.Infer = &nodes.InferClause{
			Conname: name,
			Loc:     nodes.Loc{Start: inferLoc, End: p.pos()},
		}

	default:
		// No infer clause: DO NOTHING or DO UPDATE
	}

	// DO {NOTHING | UPDATE SET ...}
	if _, err := p.expect(DO); err != nil {
		return nil, err
	}
	if p.cur.Type == NOTHING {
		p.advance()
		occ.Action = 1 // ONCONFLICT_NOTHING
	} else if p.cur.Type == UPDATE {
		p.advance()
		if _, err := p.expect(SET); err != nil {
			return nil, err
		}
		occ.Action = 2 // ONCONFLICT_UPDATE
		var err error
		occ.TargetList, err = p.parseSetClauseList()
		if err != nil {
			return nil, err
		}
		occ.WhereClause, err = p.parseWhereClause()
		if err != nil {
			return nil, err
		}
	}

	occ.Loc.End = p.pos()
	return occ, nil
}

// parseReturningClause parses an optional RETURNING clause.
//
//	returning_clause:
//	    RETURNING target_list
//	    | /* EMPTY */
func (p *Parser) parseReturningClause() (*nodes.List, error) {
	if p.cur.Type != RETURNING {
		return nil, nil
	}
	p.advance()
	return p.parseTargetList()
}

// parseSetClauseList parses a comma-separated list of SET clauses (for UPDATE/ON CONFLICT).
//
//	set_clause_list:
//	    set_clause [',' set_clause ...]
//
//	set_clause:
//	    set_target '=' a_expr
//	    | '(' set_target_list ')' '=' a_expr
func (p *Parser) parseSetClauseList() (*nodes.List, error) {
	var items []nodes.Node
	first, err := p.parseSetClause()
	if err != nil {
		return nil, err
	}
	items = append(items, first...)
	for p.cur.Type == ',' {
		p.advance()
		clause, err := p.parseSetClause()
		if err != nil {
			return nil, err
		}
		items = append(items, clause...)
	}
	return &nodes.List{Items: items}, nil
}

// parseSetClause parses a single SET clause. Returns a slice because
// multi-column assignment (a, b) = expr expands to multiple ResTargets.
func (p *Parser) parseSetClause() ([]nodes.Node, error) {
	if p.collectMode() {
		p.addTokenCandidate('(')
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	if p.cur.Type == '(' {
		// Multi-column: '(' set_target_list ')' '=' a_expr
		p.advance() // consume '('
		targets, err := p.parseSetTargetList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}

		ncolumns := len(targets.Items)
		var result []nodes.Node
		for i, t := range targets.Items {
			rt := t.(*nodes.ResTarget)
			rt.Val = &nodes.MultiAssignRef{
				Source:   expr,
				Colno:    i + 1,
				Ncolumns: ncolumns,
				Loc:      nodes.NodeLoc(expr),
			}
			result = append(result, rt)
		}
		return result, nil
	}

	// Single column: set_target '=' a_expr
	rt, err := p.parseSetTarget()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	rt.Val, err = p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	return []nodes.Node{rt}, nil
}

// parseSetTarget parses a SET target column.
//
//	set_target:
//	    ColId opt_indirection
func (p *Parser) parseSetTarget() (*nodes.ResTarget, error) {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	ind, err := p.parseOptIndirection()
	if err != nil {
		return nil, err
	}
	return &nodes.ResTarget{
		Name:        name,
		Indirection: ind,
		Loc:         nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseSetTargetList parses a comma-separated list of SET targets.
func (p *Parser) parseSetTargetList() (*nodes.List, error) {
	first, err := p.parseSetTarget()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		target, err := p.parseSetTarget()
		if err != nil {
			return nil, err
		}
		items = append(items, target)
	}
	return &nodes.List{Items: items}, nil
}

// parseIndexParams and parseIndexElem are defined in create_index.go
