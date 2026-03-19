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
func (p *Parser) parseInsertStmt(withClause *nodes.WithClause) *nodes.InsertStmt {
	loc := p.pos()
	p.advance() // consume INSERT
	p.expect(INTO)

	// insert_target: qualified_name [AS ColId]
	rvLoc := p.pos()
	names, _ := p.parseQualifiedName()
	if names == nil {
		return nil // collect-mode: parseQualifiedName already emitted rule candidates
	}
	rv := makeRangeVarFromNames(names)
	if p.cur.Type == AS {
		p.advance()
		alias, _ := p.parseColId()
		rv.Alias = &nodes.Alias{Aliasname: alias}
	}
	rv.Loc = nodes.Loc{Start: rvLoc, End: p.pos()}

	if p.collectMode() {
		// After table name, valid continuations for insert_rest:
		for _, t := range []int{'(', DEFAULT, OVERRIDING, SELECT, VALUES, TABLE, WITH} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// insert_rest
	stmt := p.parseInsertRest()
	if stmt == nil {
		return nil
	}
	stmt.Relation = rv

	if p.collectMode() {
		// After insert_rest, valid continuations:
		for _, t := range []int{ON, RETURNING, ';'} {
			p.addTokenCandidate(t)
		}
		return nil
	}

	// opt_on_conflict
	if p.cur.Type == ON {
		next := p.peekNext()
		if next.Type == CONFLICT {
			stmt.OnConflictClause = p.parseOnConflict()
		}
	}

	// returning_clause
	stmt.ReturningList = p.parseReturningClause()

	if withClause != nil {
		stmt.WithClause = withClause
		loc = withClause.Loc.Start
	}

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt
}

// parseInsertRest parses the body of an INSERT statement after the target table.
func (p *Parser) parseInsertRest() *nodes.InsertStmt {
	if p.collectMode() {
		// Valid first tokens for insert_rest:
		for _, t := range []int{DEFAULT, '(', OVERRIDING, SELECT, VALUES, TABLE, WITH} {
			p.addTokenCandidate(t)
		}
		return nil
	}
	switch p.cur.Type {
	case DEFAULT:
		// DEFAULT VALUES
		p.advance()
		p.expect(VALUES)
		return &nodes.InsertStmt{}

	case '(':
		// '(' insert_column_list ')' [OVERRIDING override_kind VALUE] SelectStmt
		p.advance() // consume '('
		cols := p.parseInsertColumnList()
		p.expect(')')

		var override nodes.OverridingKind
		if p.cur.Type == OVERRIDING {
			override = p.parseOverriding()
		}

		selectStmt := p.parseSelectStmt()
		return &nodes.InsertStmt{
			Cols:       cols,
			Override:   override,
			SelectStmt: selectStmt,
		}

	case OVERRIDING:
		// OVERRIDING override_kind VALUE SelectStmt
		override := p.parseOverriding()
		selectStmt := p.parseSelectStmt()
		return &nodes.InsertStmt{
			Override:   override,
			SelectStmt: selectStmt,
		}

	default:
		// SelectStmt
		selectStmt := p.parseSelectStmt()
		return &nodes.InsertStmt{
			SelectStmt: selectStmt,
		}
	}
}

// parseOverriding parses OVERRIDING {USER|SYSTEM} VALUE.
func (p *Parser) parseOverriding() nodes.OverridingKind {
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
	p.expect(VALUE_P)
	return kind
}

// parseSelectStmt parses SelectStmt: select_no_parens | select_with_parens.
func (p *Parser) parseSelectStmt() nodes.Node {
	if p.cur.Type == '(' {
		return p.parseSelectWithParens()
	}
	return p.parseSelectNoParens()
}

// parseInsertColumnList parses a comma-separated list of insert columns.
//
//	insert_column_list:
//	    insert_column_item [',' insert_column_item ...]
func (p *Parser) parseInsertColumnList() *nodes.List {
	first := p.parseInsertColumnItem()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseInsertColumnItem())
	}
	return &nodes.List{Items: items}
}

// parseInsertColumnItem parses a single insert column.
//
//	insert_column_item:
//	    ColId opt_indirection
func (p *Parser) parseInsertColumnItem() *nodes.ResTarget {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil
	}
	loc := p.pos()
	name, _ := p.parseColId()
	ind, _ := p.parseOptIndirection()
	return &nodes.ResTarget{
		Name:        name,
		Indirection: ind,
		Loc:         nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseOnConflict parses ON CONFLICT clause.
//
//	opt_on_conflict:
//	    ON CONFLICT DO NOTHING
//	    | ON CONFLICT DO UPDATE SET set_clause_list where_clause
//	    | ON CONFLICT '(' index_params ')' [WHERE a_expr] DO {NOTHING | UPDATE SET ...}
//	    | ON CONFLICT ON CONSTRAINT name DO {NOTHING | UPDATE SET ...}
func (p *Parser) parseOnConflict() *nodes.OnConflictClause {
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
		p.expect(')')

		infer := &nodes.InferClause{
			IndexElems: indexElems,
			Loc:        nodes.Loc{Start: inferLoc, End: -1},
		}

		// Optional WHERE clause for partial index predicate
		if p.cur.Type == WHERE {
			p.advance()
			infer.WhereClause = p.parseAExpr(0)
		}

		infer.Loc.End = p.pos()
		occ.Infer = infer

	case ON:
		// ON CONSTRAINT name
		inferLoc := p.pos()
		p.advance() // consume ON (second one)
		p.expect(CONSTRAINT)
		name, _ := p.parseName()
		occ.Infer = &nodes.InferClause{
			Conname: name,
			Loc:     nodes.Loc{Start: inferLoc, End: p.pos()},
		}

	default:
		// No infer clause: DO NOTHING or DO UPDATE
	}

	// DO {NOTHING | UPDATE SET ...}
	p.expect(DO)
	if p.cur.Type == NOTHING {
		p.advance()
		occ.Action = 1 // ONCONFLICT_NOTHING
	} else if p.cur.Type == UPDATE {
		p.advance()
		p.expect(SET)
		occ.Action = 2 // ONCONFLICT_UPDATE
		occ.TargetList = p.parseSetClauseList()
		occ.WhereClause = p.parseWhereClause()
	}

	occ.Loc.End = p.pos()
	return occ
}

// parseReturningClause parses an optional RETURNING clause.
//
//	returning_clause:
//	    RETURNING target_list
//	    | /* EMPTY */
func (p *Parser) parseReturningClause() *nodes.List {
	if p.cur.Type != RETURNING {
		return nil
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
func (p *Parser) parseSetClauseList() *nodes.List {
	var items []nodes.Node
	first := p.parseSetClause()
	items = append(items, first...)
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseSetClause()...)
	}
	return &nodes.List{Items: items}
}

// parseSetClause parses a single SET clause. Returns a slice because
// multi-column assignment (a, b) = expr expands to multiple ResTargets.
func (p *Parser) parseSetClause() []nodes.Node {
	if p.collectMode() {
		p.addTokenCandidate('(')
		p.addRuleCandidate("columnref")
		return nil
	}
	if p.cur.Type == '(' {
		// Multi-column: '(' set_target_list ')' '=' a_expr
		p.advance() // consume '('
		targets := p.parseSetTargetList()
		p.expect(')')
		p.expect('=')
		expr := p.parseAExpr(0)

		ncolumns := len(targets.Items)
		var result []nodes.Node
		for i, t := range targets.Items {
			rt := t.(*nodes.ResTarget)
			rt.Val = &nodes.MultiAssignRef{
				Source:   expr,
				Colno:    i + 1,
				Ncolumns: ncolumns,
			}
			result = append(result, rt)
		}
		return result
	}

	// Single column: set_target '=' a_expr
	rt := p.parseSetTarget()
	p.expect('=')
	rt.Val = p.parseAExpr(0)
	return []nodes.Node{rt}
}

// parseSetTarget parses a SET target column.
//
//	set_target:
//	    ColId opt_indirection
func (p *Parser) parseSetTarget() *nodes.ResTarget {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil
	}
	loc := p.pos()
	name, _ := p.parseColId()
	ind, _ := p.parseOptIndirection()
	return &nodes.ResTarget{
		Name:        name,
		Indirection: ind,
		Loc:         nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseSetTargetList parses a comma-separated list of SET targets.
func (p *Parser) parseSetTargetList() *nodes.List {
	first := p.parseSetTarget()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseSetTarget())
	}
	return &nodes.List{Items: items}
}

// parseIndexParams and parseIndexElem are defined in create_index.go
