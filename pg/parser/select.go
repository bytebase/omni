package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// Set operation precedence levels for precedence climbing.
const (
	setOpPrecNone      = 0
	setOpPrecUnion     = 1 // UNION, EXCEPT
	setOpPrecIntersect = 2 // INTERSECT (binds tighter)
)

// ---------------------------------------------------------------------------
// parseSelectNoParens parses a top-level SELECT (not wrapped in parentheses).
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	select_no_parens:
//	    simple_select
//	    | select_clause sort_clause
//	    | select_clause opt_sort_clause for_locking_clause opt_select_limit
//	    | select_clause opt_sort_clause select_limit opt_for_locking_clause
//	    | with_clause select_clause ...
func (p *Parser) parseSelectNoParens() (*nodes.SelectStmt, error) {
	loc := p.pos()
	// 1. Optional WITH clause
	var withClause *nodes.WithClause
	if p.cur.Type == WITH || p.cur.Type == WITH_LA {
		var err error
		withClause, err = p.parseWithClause()
		if err != nil {
			return nil, err
		}
	}

	// 2. Parse select_clause (handles UNION/INTERSECT/EXCEPT)
	stmt, err := p.parseSelectClause(setOpPrecNone)
	if err != nil {
		return nil, err
	}
	if stmt == nil {
		return nil, nil
	}

	// 3. Optional ORDER BY
	if p.cur.Type == ORDER {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		stmt.SortClause, err = p.parseSortByList()
		if err != nil {
			return nil, err
		}
	}

	// 4. Parse LIMIT/OFFSET and FOR locking in either order
	if err := p.parseSelectOptions(stmt); err != nil {
		return nil, err
	}

	if withClause != nil {
		stmt.WithClause = withClause
	}

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseSelectWithParens parses a parenthesized SELECT statement.
//
//	select_with_parens:
//	    '(' select_no_parens ')'
//	    | '(' select_with_parens ')'
func (p *Parser) parseSelectWithParens() (*nodes.SelectStmt, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if p.collectMode() {
		p.addTokenCandidate('(')
		p.addTokenCandidate(SELECT)
		p.addTokenCandidate(VALUES)
		p.addTokenCandidate(TABLE)
		p.addTokenCandidate(WITH)
		return nil, nil
	}
	// Delegate to parseSelectNoParens unconditionally. PG's grammar says
	//
	//	select_with_parens: '(' select_no_parens ')' | '(' select_with_parens ')'
	//
	// but the second production is already reachable via select_no_parens →
	// select_clause → parseSelectClausePrimary → '(' → parseSelectWithParens.
	// Short-circuiting on cur == '(' here (as a hand-written recursion)
	// skips parseSelectClause's set-op precedence climbing, so
	//
	//	SELECT * FROM ((SELECT 1) UNION (SELECT 2))
	//
	// failed: the inner '(SELECT 1)' was consumed as a nested
	// select_with_parens and the outer expected ')' instead of UNION.
	// Left-factoring to always go through parseSelectNoParens lets
	// select_clause handle nested select_with_parens as the first operand
	// of a set op, matching PG's LALR reduction.
	stmt, err := p.parseSelectNoParens()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseSelectClause parses a select_clause using precedence climbing
// for set operations (UNION/INTERSECT/EXCEPT).
//
//	select_clause:
//	    simple_select
//	    | select_with_parens
func (p *Parser) parseSelectClause(minPrec int) (*nodes.SelectStmt, error) {
	left, err := p.parseSelectClausePrimary()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}

	for {
		var prec int
		var op nodes.SetOperation
		switch p.cur.Type {
		case UNION:
			prec, op = setOpPrecUnion, nodes.SETOP_UNION
		case EXCEPT:
			prec, op = setOpPrecUnion, nodes.SETOP_EXCEPT
		case INTERSECT:
			prec, op = setOpPrecIntersect, nodes.SETOP_INTERSECT
		default:
			return left, nil
		}
		if prec < minPrec {
			return left, nil
		}
		p.advance() // consume UNION/INTERSECT/EXCEPT
		all := p.parseSetQuantifier()
		right, err := p.parseSelectClause(prec + 1) // left-associative
		if err != nil {
			return nil, err
		}
		result := &nodes.SelectStmt{
			Op:   op,
			Larg: left,
			Rarg: right,
		}
		if all == nodes.SET_QUANTIFIER_ALL {
			result.All = true
		}
		result.Loc = nodes.Loc{Start: left.Loc.Start, End: p.prev.End}
		left = result
	}
}

// parseSelectClausePrimary parses a primary select clause (leaf or parenthesized).
func (p *Parser) parseSelectClausePrimary() (*nodes.SelectStmt, error) {
	if p.cur.Type == '(' {
		return p.parseSelectWithParens()
	}
	return p.parseSimpleSelectLeaf()
}

// parseSimpleSelectLeaf parses a leaf simple_select (SELECT, VALUES, TABLE).
//
//	simple_select:
//	    SELECT [ALL|DISTINCT [ON (...)]] target_list [INTO ...] [FROM ...]
//	           [WHERE ...] [GROUP BY ...] [HAVING ...] [WINDOW ...]
//	    | VALUES (...)
//	    | TABLE relation_expr
func (p *Parser) parseSimpleSelectLeaf() (*nodes.SelectStmt, error) {
	switch p.cur.Type {
	case SELECT:
		return p.parseSimpleSelectCore()
	case VALUES:
		return p.parseValuesClause()
	case TABLE:
		return p.parseTableCmd()
	default:
		return nil, nil
	}
}

// parseSimpleSelectCore parses the core of a SELECT statement after the SELECT keyword.
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
//	    [ * | expression [ [ AS ] output_name ] [, ...] ]
//	    [ INTO new_table ]
//	    [ FROM from_item [, ...] ]
//	    [ WHERE condition ]
//	    [ GROUP BY [ ALL | DISTINCT ] grouping_element [, ...] ]
//	    [ HAVING condition ]
//	    [ WINDOW window_name AS ( window_definition ) [, ...] ]
func (p *Parser) parseSimpleSelectCore() (*nodes.SelectStmt, error) {
	loc := p.pos()
	p.advance() // consume SELECT
	stmt := &nodes.SelectStmt{}

	// ALL / DISTINCT
	if p.cur.Type == ALL {
		p.advance()
	} else if p.cur.Type == DISTINCT {
		p.advance()
		if p.cur.Type == ON {
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			var err error
			stmt.DistinctClause, err = p.parseExprListFull()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		} else {
			stmt.DistinctClause = &nodes.List{Items: []nodes.Node{nil}}
		}
	}

	// target_list (opt_target_list)
	var err error
	stmt.TargetList, err = p.parseTargetList()
	if err != nil {
		return nil, err
	}

	if p.collectMode() {
		// After target list, valid continuations:
		for _, t := range []int{
			INTO, FROM, WHERE, GROUP_P, HAVING, WINDOW,
			ORDER, LIMIT, OFFSET, FETCH, FOR,
			UNION, EXCEPT, INTERSECT, ';',
		} {
			p.addTokenCandidate(t)
		}
		// Also valid: expression operators (AS for alias, comma for next target)
		p.addTokenCandidate(',')
		p.addTokenCandidate(AS)
		return stmt, nil
	}

	// INTO clause
	if p.cur.Type == INTO {
		stmt.IntoClause, err = p.parseIntoClause()
		if err != nil {
			return nil, err
		}
	}

	// FROM clause
	if p.cur.Type == FROM {
		p.advance()
		stmt.FromClause, err = p.parseFromListFull()
		if err != nil {
			return nil, err
		}
		if stmt.FromClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.collectMode() {
			for _, t := range []int{
				WHERE, GROUP_P, HAVING, WINDOW,
				ORDER, LIMIT, OFFSET, FETCH, FOR,
				UNION, EXCEPT, INTERSECT, ';',
			} {
				p.addTokenCandidate(t)
			}
			return stmt, nil
		}
	}

	// WHERE clause
	if p.cur.Type == WHERE {
		p.advance()
		stmt.WhereClause, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && stmt.WhereClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.collectMode() {
			for _, t := range []int{
				GROUP_P, HAVING, WINDOW,
				ORDER, LIMIT, OFFSET, FETCH, FOR,
				UNION, EXCEPT, INTERSECT, ';',
			} {
				p.addTokenCandidate(t)
			}
			return stmt, nil
		}
	}

	// GROUP BY clause
	if p.cur.Type == GROUP_P {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		// Check for ALL or DISTINCT
		if p.cur.Type == ALL {
			p.advance()
			// GROUP BY ALL is treated the same as GROUP BY (all is default)
		} else if p.cur.Type == DISTINCT {
			p.advance()
			stmt.GroupDistinct = true
		}
		stmt.GroupClause, err = p.parseGroupByList()
		if err != nil {
			return nil, err
		}
		if p.collectMode() {
			for _, t := range []int{
				HAVING, WINDOW, ORDER, LIMIT, OFFSET, FETCH, FOR,
				UNION, EXCEPT, INTERSECT, ';',
			} {
				p.addTokenCandidate(t)
			}
			return stmt, nil
		}
	}

	// HAVING clause
	if p.cur.Type == HAVING {
		p.advance()
		stmt.HavingClause, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
	}

	// WINDOW clause
	stmt.WindowClause, err = p.parseWindowClause()
	if err != nil {
		return nil, err
	}

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseValuesClause parses a VALUES clause.
//
//	values_clause:
//	    VALUES '(' expr_list ')' [',' '(' expr_list ')' ...]
func (p *Parser) parseValuesClause() (*nodes.SelectStmt, error) {
	loc := p.pos()
	p.advance() // consume VALUES
	stmt := &nodes.SelectStmt{}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	first, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.ValuesLists = &nodes.List{Items: []nodes.Node{first}}

	for p.cur.Type == ',' {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprListFull()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.ValuesLists.Items = append(stmt.ValuesLists.Items, exprs)
	}

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseTableCmd parses TABLE relation_expr (shorthand for SELECT * FROM ...).
//
//	TABLE relation_expr
func (p *Parser) parseTableCmd() (*nodes.SelectStmt, error) {
	loc := p.pos()
	p.advance() // consume TABLE
	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	cr := &nodes.ColumnRef{
		Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{Loc: nodes.NoLoc()}}},
		Loc:    nodes.NoLoc(),
	}
	rt := &nodes.ResTarget{
		Val: cr,
		Loc: nodes.NoLoc(),
	}
	stmt := &nodes.SelectStmt{
		TargetList: &nodes.List{Items: []nodes.Node{rt}},
		FromClause: &nodes.List{Items: []nodes.Node{rel}},
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseSetQuantifier parses ALL | DISTINCT | EMPTY for set operations.
//
//	set_quantifier:
//	    ALL | DISTINCT | /* EMPTY */
func (p *Parser) parseSetQuantifier() nodes.SetQuantifier {
	if p.cur.Type == ALL {
		p.advance()
		return nodes.SET_QUANTIFIER_ALL
	}
	if p.cur.Type == DISTINCT {
		p.advance()
		return nodes.SET_QUANTIFIER_DISTINCT
	}
	return nodes.SET_QUANTIFIER_DEFAULT
}

// parseSelectOptions parses ORDER BY, LIMIT/OFFSET, and FOR locking clauses.
// These can appear in various orders.
func (p *Parser) parseSelectOptions(stmt *nodes.SelectStmt) error {
	// Parse LIMIT/OFFSET/FETCH and FOR locking in either order
	for {
		if p.collectMode() {
			for _, t := range []int{LIMIT, OFFSET, FETCH, FOR, ';'} {
				p.addTokenCandidate(t)
			}
			return nil
		}
		if p.cur.Type == LIMIT || p.cur.Type == OFFSET || p.cur.Type == FETCH {
			if err := p.parseSelectLimit(stmt); err != nil {
				return err
			}
		} else if p.cur.Type == FOR {
			locking, err := p.parseForLockingClause()
			if err != nil {
				return err
			}
			if locking != nil {
				stmt.LockingClause = locking
			}
		} else {
			break
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// LIMIT / OFFSET / FETCH

// parseSelectLimit parses LIMIT, OFFSET, and FETCH clauses.
//
//	select_limit:
//	    limit_clause offset_clause
//	    | offset_clause limit_clause
//	    | limit_clause
//	    | offset_clause
func (p *Parser) parseSelectLimit(stmt *nodes.SelectStmt) error {
	for {
		switch p.cur.Type {
		case LIMIT:
			p.advance()
			if p.cur.Type == ALL {
				p.advance()
				stmt.LimitCount = &nodes.A_Const{Isnull: true}
				stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
			} else {
				var err error
				stmt.LimitCount, err = p.parseAExpr(0)
				if err != nil {
					return err
				}
				stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
				// Check for deprecated LIMIT #,# syntax
				if p.cur.Type == ',' {
					p.advance()
					// In MySQL-compatible syntax: LIMIT offset, count
					// Swap: what we parsed as count is actually offset
					stmt.LimitOffset = stmt.LimitCount
					stmt.LimitCount, err = p.parseAExpr(0)
					if err != nil {
						return err
					}
				}
			}
		case OFFSET:
			p.advance()
			// Could be: OFFSET select_offset_value
			// or: OFFSET select_fetch_first_value row_or_rows
			expr, err := p.parseAExpr(0)
			if err != nil {
				return err
			}
			// Check if followed by ROW or ROWS (SQL:2008 syntax)
			if p.cur.Type == ROW || p.cur.Type == ROWS {
				p.advance()
			}
			stmt.LimitOffset = expr
		case FETCH:
			if err := p.parseFetchClause(stmt); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// parseFetchClause parses FETCH FIRST/NEXT ... ROWS ONLY/WITH TIES.
//
//	FETCH first_or_next select_fetch_first_value row_or_rows ONLY
//	| FETCH first_or_next select_fetch_first_value row_or_rows WITH TIES
//	| FETCH first_or_next row_or_rows ONLY
//	| FETCH first_or_next row_or_rows WITH TIES
func (p *Parser) parseFetchClause(stmt *nodes.SelectStmt) error {
	p.advance() // consume FETCH
	// first_or_next: FIRST | NEXT (ignored semantically)
	if p.cur.Type == FIRST_P || p.cur.Type == NEXT {
		p.advance()
	}

	// Check if next is ROW/ROWS (no count specified, defaults to 1)
	if p.cur.Type == ROW || p.cur.Type == ROWS {
		p.advance()
		stmt.LimitCount = &nodes.A_Const{Val: &nodes.Integer{Ival: 1}}
	} else {
		// select_fetch_first_value: c_expr | '+' c_expr | '-' c_expr
		var err error
		if p.cur.Type == '+' {
			p.advance()
			stmt.LimitCount, err = p.parseCExpr()
			if err != nil {
				return err
			}
		} else if p.cur.Type == '-' {
			p.advance()
			limitExpr, err := p.parseCExpr()
			if err != nil {
				return err
			}
			stmt.LimitCount = doNegate(limitExpr)
		} else {
			stmt.LimitCount, err = p.parseCExpr()
			if err != nil {
				return err
			}
		}
		// row_or_rows: ROW | ROWS
		if p.cur.Type == ROW || p.cur.Type == ROWS {
			p.advance()
		}
	}

	// ONLY or WITH TIES
	if p.cur.Type == WITH {
		p.advance()
		if _, err := p.expect(TIES); err != nil {
			return err
		}
		stmt.LimitOption = nodes.LIMIT_OPTION_WITH_TIES
	} else {
		if _, err := p.expect(ONLY); err != nil {
			return err
		}
		stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
	}
	return nil
}

// ---------------------------------------------------------------------------
// FOR locking clause

// parseForLockingClause parses FOR UPDATE/SHARE/NO KEY UPDATE/KEY SHARE clauses.
//
//	for_locking_clause:
//	    for_locking_items
//	    | FOR READ ONLY
func (p *Parser) parseForLockingClause() (*nodes.List, error) {
	if p.cur.Type != FOR {
		return nil, nil
	}

	// Check for FOR READ ONLY
	next := p.peekNext()
	if next.Type == READ {
		p.advance() // consume FOR
		p.advance() // consume READ
		if _, err := p.expect(ONLY); err != nil {
			return nil, err
		}
		return nil, nil // FOR READ ONLY = no locking
	}

	var items []nodes.Node
	for p.cur.Type == FOR {
		item, err := p.parseForLockingItem()
		if err != nil {
			return nil, err
		}
		if item != nil {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseForLockingItem parses a single FOR UPDATE/SHARE item.
//
//	for_locking_item:
//	    for_locking_strength locked_rels_list opt_nowait_or_skip
func (p *Parser) parseForLockingItem() (*nodes.LockingClause, error) {
	forLoc := p.pos()
	strength := p.parseForLockingStrength()
	if strength < 0 {
		return nil, nil
	}

	lc := &nodes.LockingClause{
		Strength: int(strength),
	}

	// locked_rels_list: OF qualified_name_list | EMPTY
	if p.cur.Type == OF {
		p.advance()
		var err error
		lc.LockedRels, err = p.parseQualifiedNameList()
		if err != nil {
			return nil, err
		}
	}

	// opt_nowait_or_skip
	if p.cur.Type == NOWAIT {
		p.advance()
		lc.WaitPolicy = int(nodes.LockWaitError)
	} else if p.cur.Type == SKIP {
		p.advance()
		if _, err := p.expect(LOCKED); err != nil {
			return nil, err
		}
		lc.WaitPolicy = int(nodes.LockWaitSkip)
	} else {
		lc.WaitPolicy = int(nodes.LockWaitBlock)
	}

	lc.Loc = nodes.Loc{Start: forLoc, End: p.prev.End}
	return lc, nil
}

// parseForLockingStrength parses FOR UPDATE/NO KEY UPDATE/SHARE/KEY SHARE.
func (p *Parser) parseForLockingStrength() int64 {
	if p.cur.Type != FOR {
		return -1
	}
	p.advance() // consume FOR

	switch p.cur.Type {
	case UPDATE:
		p.advance()
		return int64(nodes.LCS_FORUPDATE)
	case NO:
		p.advance() // consume NO
		p.expect(KEY)
		p.expect(UPDATE)
		return int64(nodes.LCS_FORNOKEYUPDATE)
	case SHARE:
		p.advance()
		return int64(nodes.LCS_FORSHARE)
	case KEY:
		p.advance() // consume KEY
		p.expect(SHARE)
		return int64(nodes.LCS_FORKEYSHARE)
	default:
		return -1
	}
}

// ---------------------------------------------------------------------------
// WITH clause (Common Table Expressions)

// parseWithClause parses a WITH clause.
//
//	with_clause:
//	    WITH cte_list
//	    | WITH_LA cte_list
//	    | WITH RECURSIVE cte_list
func (p *Parser) parseWithClause() (*nodes.WithClause, error) {
	loc := p.pos()
	// Record CTE position for completion context (before consuming WITH).
	if p.completing {
		p.addCTEPosition(loc)
	}
	p.advance() // consume WITH or WITH_LA

	wc := &nodes.WithClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == RECURSIVE {
		p.advance()
		wc.Recursive = true
	}

	var err error
	wc.Ctes, err = p.parseCTEList()
	if err != nil {
		return nil, err
	}
	wc.Loc.End = p.pos()
	return wc, nil
}

// parseCTEList parses a comma-separated list of CTEs.
func (p *Parser) parseCTEList() (*nodes.List, error) {
	first, err := p.parseCommonTableExpr()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		cte, err := p.parseCommonTableExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, cte)
	}
	return &nodes.List{Items: items}, nil
}

// parseCommonTableExpr parses a single CTE.
//
//	common_table_expr:
//	    name opt_name_list AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
//	    | name '(' name_list ')' AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
func (p *Parser) parseCommonTableExpr() (*nodes.CommonTableExpr, error) {
	loc := p.pos()
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	cte := &nodes.CommonTableExpr{
		Ctename: name,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Check for name '(' name_list ')' form
	if p.cur.Type == '(' {
		p.advance()
		colnames, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		cte.Aliascolnames = colnames
	}

	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	// opt_materialized
	cte.Ctematerialized = int(p.parseOptMaterialized())

	// '(' PreparableStmt ')'
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	ctequery, err := p.parsePreparableStmt()
	if err != nil {
		return nil, err
	}
	cte.Ctequery = ctequery
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// opt_search_clause
	var searchErr error
	cte.SearchClause, searchErr = p.parseOptSearchClause()
	if searchErr != nil {
		return nil, searchErr
	}
	// opt_cycle_clause
	var cycleErr error
	cte.CycleClause, cycleErr = p.parseOptCycleClause()
	if cycleErr != nil {
		return nil, cycleErr
	}

	cte.Loc.End = p.pos()
	return cte, nil
}

// parseOptMaterialized parses MATERIALIZED / NOT MATERIALIZED / empty.
func (p *Parser) parseOptMaterialized() nodes.CTEMaterialize {
	if p.cur.Type == MATERIALIZED {
		p.advance()
		return nodes.CTEMaterializeAlways
	}
	if p.cur.Type == NOT {
		next := p.peekNext()
		if next.Type == MATERIALIZED {
			p.advance() // consume NOT
			p.advance() // consume MATERIALIZED
			return nodes.CTEMaterializeNever
		}
	}
	return nodes.CTEMaterializeDefault
}

// parseOptSearchClause parses an optional SEARCH clause in a recursive CTE.
//
// Ref: https://www.postgresql.org/docs/17/queries-with.html
//
//	opt_search_clause:
//	    SEARCH DEPTH FIRST_P BY columnList SET ColId
//	    | SEARCH BREADTH FIRST_P BY columnList SET ColId
//	    | /* EMPTY */
func (p *Parser) parseOptSearchClause() (nodes.Node, error) {
	if p.cur.Type != SEARCH {
		return nil, nil
	}
	p.advance() // consume SEARCH

	breadthFirst := false
	switch p.cur.Type {
	case DEPTH:
		p.advance()
		breadthFirst = false
	case BREADTH:
		p.advance()
		breadthFirst = true
	}
	if _, err := p.expect(FIRST_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}
	colList := p.parseColumnList()
	if _, err := p.expect(SET); err != nil {
		return nil, err
	}
	seqCol, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	return &nodes.CTESearchClause{
		SearchColList:      colList,
		SearchBreadthFirst: breadthFirst,
		SearchSeqColumn:    seqCol,
		Loc:                nodes.NoLoc(),
	}, nil
}

// parseOptCycleClause parses an optional CYCLE clause in a recursive CTE.
//
// Ref: https://www.postgresql.org/docs/17/queries-with.html
//
//	opt_cycle_clause:
//	    CYCLE columnList SET ColId TO AexprConst DEFAULT AexprConst USING ColId
//	    | CYCLE columnList SET ColId USING ColId
//	    | /* EMPTY */
func (p *Parser) parseOptCycleClause() (nodes.Node, error) {
	if p.cur.Type != CYCLE {
		return nil, nil
	}
	p.advance() // consume CYCLE

	colList := p.parseColumnList()
	if _, err := p.expect(SET); err != nil {
		return nil, err
	}
	markCol, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	var markValue, markDefault nodes.Node

	if p.cur.Type == TO {
		// CYCLE ... SET col TO val DEFAULT val USING col
		p.advance() // consume TO
		markValue, err = p.parseCExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(DEFAULT); err != nil {
			return nil, err
		}
		markDefault, err = p.parseCExpr()
		if err != nil {
			return nil, err
		}
	} else {
		// CYCLE ... SET col USING col (implicit TRUE/FALSE)
		markValue = makeBoolAConst(1)
		markDefault = makeBoolAConst(0)
	}

	if _, err := p.expect(USING); err != nil {
		return nil, err
	}
	pathCol, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	return &nodes.CTECycleClause{
		CycleColList:     colList,
		CycleMarkColumn:  markCol,
		CycleMarkValue:   markValue,
		CycleMarkDefault: markDefault,
		CyclePathColumn:  pathCol,
		Loc:              nodes.NoLoc(),
	}, nil
}

// ---------------------------------------------------------------------------
// INTO clause

// parseIntoClause parses SELECT INTO clause.
//
//	into_clause:
//	    INTO OptTempTableName
func (p *Parser) parseIntoClause() (*nodes.IntoClause, error) {
	intoLoc := p.pos()
	p.advance() // consume INTO

	rv, err := p.parseOptTempTableName()
	if err != nil {
		return nil, err
	}
	return &nodes.IntoClause{
		Rel:      rv,
		OnCommit: nodes.ONCOMMIT_NOOP,
		Loc:      nodes.Loc{Start: intoLoc, End: p.prev.End},
	}, nil
}

// parseOptTempTableName parses the target table name for SELECT INTO.
//
//	OptTempTableName:
//	    TEMPORARY opt_table qualified_name
//	    | TEMP opt_table qualified_name
//	    | LOCAL TEMPORARY opt_table qualified_name
//	    | LOCAL TEMP opt_table qualified_name
//	    | GLOBAL TEMPORARY opt_table qualified_name
//	    | GLOBAL TEMP opt_table qualified_name
//	    | UNLOGGED opt_table qualified_name
//	    | TABLE qualified_name
//	    | qualified_name
func (p *Parser) parseOptTempTableName() (*nodes.RangeVar, error) {
	loc := p.pos()
	persistence := byte('p')

	switch p.cur.Type {
	case TEMPORARY, TEMP:
		p.advance()
		persistence = 't'
		p.parseOptTable()
	case LOCAL, GLOBAL:
		p.advance()
		if p.cur.Type == TEMPORARY || p.cur.Type == TEMP {
			p.advance()
		}
		persistence = 't'
		p.parseOptTable()
	case UNLOGGED:
		p.advance()
		persistence = 'u'
		p.parseOptTable()
	case TABLE:
		p.advance()
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = persistence
	rv.Loc = nodes.Loc{Start: loc, End: p.pos()}
	return rv, nil
}

// parseOptTable consumes optional TABLE keyword.
func (p *Parser) parseOptTable() {
	if p.cur.Type == TABLE {
		p.advance()
	}
}

// ---------------------------------------------------------------------------
// FROM clause

// parseFromListFull parses a comma-separated list of table references with join support.
func (p *Parser) parseFromListFull() (*nodes.List, error) {
	first, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if first == nil {
		return nil, nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		item, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

// parseTableRef parses a table reference with optional joins.
//
//	table_ref:
//	    relation_expr opt_alias_clause [tablesample_clause]
//	    | select_with_parens opt_alias_clause
//	    | func_table func_alias_clause
//	    | LATERAL_P func_table func_alias_clause
//	    | LATERAL_P select_with_parens opt_alias_clause
//	    | joined_table
//	    | '(' joined_table ')' opt_alias_clause
//	    | xmltable opt_alias_clause
//	    | json_table opt_alias_clause
func (p *Parser) parseTableRef() (nodes.Node, error) {
	left, err := p.parseTableRefPrimary()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}

	// Handle joins (left-recursive in yacc, loop in recursive descent)
	for {
		joined, err := p.tryParseJoin(left)
		if err != nil {
			return nil, err
		}
		if joined == nil {
			break
		}
		left = joined
	}

	return left, nil
}

// parseTableRefPrimary parses a primary (non-join) table reference.
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	table_ref:
//	    relation_expr opt_alias_clause [tablesample_clause]
//	    | func_table func_alias_clause
//	    | xmltable opt_alias_clause
//	    | json_table opt_alias_clause
//	    | LATERAL_P ...
//	    | select_with_parens opt_alias_clause
//	    | '(' joined_table ')' opt_alias_clause
func (p *Parser) parseTableRefPrimary() (nodes.Node, error) {
	if p.collectMode() {
		p.addRuleCandidate("relation_expr")
		p.addRuleCandidate("qualified_name")
		p.addTokenCandidate('(')       // subquery
		p.addTokenCandidate(LATERAL_P) // LATERAL
		return nil, errCollecting
	}
	switch p.cur.Type {
	case '(':
		return p.parseParenTableRef()
	case LATERAL_P:
		return p.parseLateralTableRef()
	case XMLTABLE:
		xt, err := p.parseXmlTable()
		if err != nil {
			return nil, err
		}
		n := xt.(*nodes.RangeTableFunc)
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
			n.Loc.End = alias.Loc.End
		}
		return n, nil
	case JSON_TABLE:
		n, err := p.parseJsonTable()
		if err != nil {
			return nil, err
		}
		alias := p.parseOptAliasClause()
		if n != nil && alias != nil {
			n.Alias = alias
			n.Loc.End = alias.Loc.End
		}
		return n, nil
	default:
		// Could be relation_expr or func_table
		// relation_expr starts with qualified_name, ONLY, or ONLY '('
		// func_table starts with func_expr_windowless or ROWS FROM
		if p.cur.Type == ONLY {
			return p.parseRelationExprWithAlias()
		}
		if p.cur.Type == ROWS {
			return p.parseRowsFromTable()
		}
		if p.isColId() {
			return p.parseRelationOrFuncTable()
		}
		// func_table: func_expr_common_subexpr tokens (USER, CURRENT_USER, etc.)
		// In PostgreSQL, table_ref has a separate func_table alternative that
		// matches func_expr_windowless, which includes func_expr_common_subexpr.
		// For example, SELECT * FROM user parses USER as CURRENT_USER via this path.
		if p.isFuncExprCommonSubexprStart() {
			funcLoc := p.pos()
			funcExpr, err := p.parseFuncExprWindowless()
			if err != nil {
				return nil, err
			}
			return p.finishFuncTable(funcExpr, funcLoc), nil
		}
		return nil, nil
	}
}

// parseParenTableRef handles '(' ... ')' in FROM clause.
// Could be: '(' select_no_parens ')' or '(' joined_table ')' opt_alias
func (p *Parser) parseParenTableRef() (nodes.Node, error) {
	loc := p.pos()
	if p.parenBeginsSubquery() {
		stmt, err := p.parseSelectWithParens()
		if err != nil {
			return nil, err
		}
		alias := p.parseOptAliasClause()
		return &nodes.RangeSubselect{
			Subquery: stmt,
			Alias:    alias,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}

	// '(' joined_table ')' opt_alias_clause
	p.advance() // consume '('
	if p.collectMode() {
		// After '(' in FROM, valid starts include both subquery keywords
		// and table references for joined tables.
		for _, t := range []int{SELECT, VALUES, WITH, TABLE, '('} {
			p.addTokenCandidate(t)
		}
		p.addRuleCandidate("relation_expr")
		p.addRuleCandidate("qualified_name")
		return nil, nil
	}
	inner, err := p.parseJoinedTable()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	alias := p.parseOptAliasClause()
	if j, ok := inner.(*nodes.JoinExpr); ok && alias != nil {
		j.Alias = alias
		j.Loc.End = alias.Loc.End
	}
	return inner, nil
}

// parseJoinedTable parses a PG joined_table nonterminal:
//
//	joined_table:
//	    '(' joined_table ')'
//	  | table_ref CROSS JOIN table_ref
//	  | table_ref [join_type] JOIN table_ref join_qual
//	  | table_ref NATURAL [join_type] JOIN table_ref
//
// A joined_table must contain at least one JOIN operator — this is what
// distinguishes the '(' joined_table ')' alternative of table_ref from
// the select_with_parens / LATERAL alternatives. omni's parseTableRef
// is the generic table_ref entry point and happily returns a single
// RangeVar or RangeSubselect when no JOIN follows; reusing it inside
// '(' ... ')' would wrongly accept things like `FROM (a)` or
// `FROM ((SELECT 1) x)` which PG rejects because the content is not a
// joined_table.
//
// Implementation note: parseTableRef already parses the primary +
// left-recursive JOIN loop (see parseTableRef comment). A valid
// joined_table is therefore "parseTableRef result whose root is a
// JoinExpr"; we enforce that assertion here.
func (p *Parser) parseJoinedTable() (nodes.Node, error) {
	node, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if p.collectMode() {
		return node, nil
	}
	if node == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, ok := node.(*nodes.JoinExpr); !ok {
		return nil, p.syntaxErrorAtCur()
	}
	return node, nil
}

// parenBeginsSubquery reports whether the '(' at p.cur begins a
// select_with_parens (a parenthesized SELECT/VALUES/WITH/TABLE
// subquery) rather than a '(' joined_table ')' parenthesized joined
// table.
//
// A single-token peek can't disambiguate these:
//
//	'(' SELECT ... ')'         — subquery
//	'(' a JOIN b ON ... ')'    — joined_table (peekNext is an identifier)
//	'(' '(' SELECT ... ')' ')' — subquery (select_with_parens wraps select_with_parens)
//	'(' '(' a JOIN b ')' ')'   — joined_table (joined_table wraps joined_table)
//	'(' '(' SELECT 1 ')' JOIN ... ')' — joined_table (inner is subquery, outer joins it)
//
// The last two both start '((' but go to different branches. Walking
// past all opening parens and looking at the first non-'(' token is
// therefore insufficient — a nested '(SELECT ...)' could be either a
// wrapped subquery or one operand of an outer join.
//
// Instead, we snapshot the token stream and scan the matched paren
// block. The outer is select_with_parens iff the depth-1 content is
// either a SELECT/VALUES/WITH/TABLE lead or a single nested
// paren-block that is itself select_with_parens (recursive). Any JOIN-
// family keyword at depth 1, or any trailing tokens after a nested
// paren-block, routes to the joined_table path.
func (p *Parser) parenBeginsSubquery() bool {
	if p.cur.Type != '(' {
		return false
	}
	snap := p.snapshotTokenStream()
	defer p.restoreTokenStream(snap)
	return p.consumeMatchedParenIsSubquery()
}

// consumeMatchedParenIsSubquery: p.cur is at '('. Consumes through
// the matched ')' (inclusive) and returns true iff the paren block is
// a select_with_parens. Must only be called inside a snapshot/restore
// scope because it mutates the token stream.
func (p *Parser) consumeMatchedParenIsSubquery() bool {
	if p.cur.Type != '(' || p.cur.Type == 0 {
		return false
	}
	p.advance() // consume '('

	var isSubquery bool
	switch p.cur.Type {
	case SELECT, VALUES, WITH, WITH_LA, TABLE:
		isSubquery = true
	case '(':
		nested := p.consumeMatchedParenIsSubquery()
		// After recursion, p.cur is past the nested ')'. Classify the
		// outer by what follows the nested paren-block:
		//   ')'                            — outer just wraps nested
		//   UNION / INTERSECT / EXCEPT     — set-op continuation; nested
		//                                    is the first operand of a
		//                                    select_no_parens, so outer
		//                                    is select_with_parens iff
		//                                    nested was.
		//   JOIN / CROSS / NATURAL / ...   — nested is the first operand
		//                                    of a joined_table; outer is
		//                                    a joined_table wrap.
		//   anything else (alias, etc.)    — default to joined_table path;
		//                                    parseJoinedTable will surface
		//                                    a proper syntax error if the
		//                                    content isn't actually a
		//                                    joined_table.
		switch p.cur.Type {
		case ')', UNION, INTERSECT, EXCEPT:
			isSubquery = nested
		default:
			isSubquery = false
		}
	default:
		isSubquery = false
	}

	// Consume through the matching outer ')'.
	depth := 1
	for depth > 0 && p.cur.Type != 0 {
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				p.advance()
				return isSubquery
			}
		}
		p.advance()
	}
	return isSubquery
}

// parseLateralTableRef handles LATERAL prefix.
//
// Ref: gram.y:13611-13620 (PG 17 table_ref LATERAL alternatives)
//
//	table_ref:
//	    ...
//	  | LATERAL_P func_table func_alias_clause
//	  | LATERAL_P select_with_parens opt_alias_clause
//	  | LATERAL_P xmltable opt_alias_clause
//	  | LATERAL_P json_table opt_alias_clause
//
// Dispatch technique: T6 (dedicated nonterminal dispatch per PG
// production). The four LATERAL-prefixed variants have disjoint leading
// tokens after LATERAL: `(` (select_with_parens), XMLTABLE, JSON_TABLE,
// and everything else (func_table via func_expr_windowless, which
// admits ROWS FROM and func_expr_common_subexpr starters).
//
// Notably absent: bare `LATERAL <relation>` is NOT a grammar production
// (relation_expr's LATERAL-less variant is the only way to spell a
// plain relation reference). We enforce that by post-filtering the
// func_table branch — if `parseFuncExprWindowless` returned a bare
// ColumnRef (no parenthesized argument list), the content wasn't a
// valid func_application and we reject it with a syntax error matching
// PG 17.
func (p *Parser) parseLateralTableRef() (nodes.Node, error) {
	lateralLoc := p.pos()
	p.advance() // consume LATERAL

	// LATERAL select_with_parens opt_alias_clause
	if p.cur.Type == '(' {
		stmt, err := p.parseSelectWithParens()
		if err != nil {
			return nil, err
		}
		alias := p.parseOptAliasClause()
		return &nodes.RangeSubselect{
			Lateral:  true,
			Subquery: stmt,
			Alias:    alias,
			Loc:      nodes.Loc{Start: lateralLoc, End: p.prev.End},
		}, nil
	}

	// LATERAL ROWS FROM (...) — a sub-production of func_table
	if p.cur.Type == ROWS {
		rf, err := p.parseRowsFromTable()
		if err != nil {
			return nil, err
		}
		if rf, ok := rf.(*nodes.RangeFunction); ok {
			rf.Lateral = true
			rf.Loc.Start = lateralLoc // extend Loc to include LATERAL keyword
		}
		return rf, nil
	}

	// LATERAL xmltable opt_alias_clause
	if p.cur.Type == XMLTABLE {
		xt, err := p.parseXmlTable()
		if err != nil {
			return nil, err
		}
		n := xt.(*nodes.RangeTableFunc)
		n.Lateral = true
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
			n.Loc.End = alias.Loc.End
		}
		return n, nil
	}

	// LATERAL json_table opt_alias_clause
	if p.cur.Type == JSON_TABLE {
		n, err := p.parseJsonTable()
		if err != nil {
			return nil, err
		}
		if n != nil {
			n.Lateral = true
			alias := p.parseOptAliasClause()
			if alias != nil {
				n.Alias = alias
				n.Loc.End = alias.Loc.End
			}
		}
		return n, nil
	}

	// LATERAL func_table func_alias_clause
	//
	// func_table = func_expr_windowless opt_ordinality (gram.y:13730).
	// func_expr_windowless = func_application | func_expr_common_subexpr
	// (gram.y:14470). Neither alternative admits a bare identifier — a
	// plain `LATERAL u` is therefore not a grammar production and must
	// be rejected.
	funcStart := p.cur
	funcExpr, err := p.parseFuncExprWindowless()
	if err != nil {
		return nil, err
	}
	// Reject bare relation references (ColumnRef without a following
	// argument list) — these are not valid func_application shapes.
	// Matches PG 17's rejection of `FROM t, LATERAL u`.
	if _, isColRef := funcExpr.(*nodes.ColumnRef); isColRef {
		return nil, p.syntaxErrorAtTok(funcStart)
	}
	rf := &nodes.RangeFunction{
		Lateral:   true,
		Functions: &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{funcExpr}}}},
	}
	p.parseOptOrdinality(rf)
	p.parseFuncAliasClause(rf)
	rf.Loc = nodes.Loc{Start: lateralLoc, End: p.prev.End}
	return rf, nil
}

// parseRowsFromTable parses ROWS FROM (...) [WITH ORDINALITY].
func (p *Parser) parseRowsFromTable() (nodes.Node, error) {
	rowsLoc := p.pos()
	p.advance() // consume ROWS
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var items []nodes.Node
	for {
		funcExpr, err := p.parseFuncExprWindowless()
		if err != nil {
			return nil, err
		}
		var colDef *nodes.List
		if p.cur.Type == AS && p.peekNext().Type == '(' {
			p.advance() // AS
			p.advance() // (
			colDef, _ = p.parseTableFuncElementList()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		item := &nodes.List{Items: []nodes.Node{funcExpr}}
		if colDef != nil {
			item.Items = append(item.Items, colDef)
		}
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	ordinality := false
	if p.cur.Type == WITH_LA && p.peekNext().Type == ORDINALITY {
		p.advance() // WITH_LA
		p.advance() // ORDINALITY
		ordinality = true
	}

	rf := &nodes.RangeFunction{
		IsRowsfrom: true,
		Ordinality: ordinality,
		Functions:  &nodes.List{Items: items},
	}
	p.parseFuncAliasClause(rf)
	rf.Loc = nodes.Loc{Start: rowsLoc, End: p.prev.End}
	return rf, nil
}

// parseTableFuncElementList is in define.go

// parseRelationOrFuncTable determines if the current position starts a relation_expr
// or a func_table, and parses accordingly.
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	table_ref:
//	    relation_expr opt_alias_clause [tablesample_clause]
//	    | func_table func_alias_clause
//
//	func_table:
//	    func_expr_windowless opt_ordinality
//	    | ROWS FROM '(' rowsfrom_list ')' opt_ordinality
//
// Disambiguation: parse the qualified name, then if '(' follows it's a function call,
// otherwise it's a relation reference.
func (p *Parser) parseRelationOrFuncTable() (nodes.Node, error) {
	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	// name( → function call
	if p.cur.Type == '(' {
		funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
		funcExpr, err := p.parseFuncApplication(funcName, loc)
		if err != nil {
			return nil, err
		}
		return p.finishFuncTable(funcExpr, loc), nil
	}

	// name.something
	if p.cur.Type == '.' {
		p.advance() // consume '.'
		// After "name.", we may have reached the cursor — emit rule candidates.
		if p.collectMode() {
			p.addRuleCandidate("relation_expr")
			p.addRuleCandidate("qualified_name")
			p.addRuleCandidate("func_name")
		}
		attr, err := p.parseAttrName()
		if err != nil {
			return nil, err
		}

		// schema.func(
		if p.cur.Type == '(' {
			funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}}
			funcExpr, err := p.parseFuncApplication(funcName, loc)
			if err != nil {
				return nil, err
			}
			return p.finishFuncTable(funcExpr, loc), nil
		}

		// catalog.schema.name
		if p.cur.Type == '.' {
			p.advance()
			// After "name.attr.", emit rule candidates for 3-part names.
			if p.collectMode() {
				p.addRuleCandidate("relation_expr")
				p.addRuleCandidate("qualified_name")
				p.addRuleCandidate("func_name")
			}
			attr2, err := p.parseAttrName()
			if err != nil {
				return nil, err
			}

			// catalog.schema.func(
			if p.cur.Type == '(' {
				funcName := &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2},
				}}
				funcExpr, err := p.parseFuncApplication(funcName, loc)
				if err != nil {
					return nil, err
				}
				return p.finishFuncTable(funcExpr, loc), nil
			}

			// 3-part qualified relation name
			names := &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2},
			}}
			rv3 := makeRangeVarFromNames(names)
			rv3.Loc = nodes.Loc{Start: loc, End: p.pos()}
			return p.finishRelationTable(rv3), nil
		}

		// 2-part qualified relation name
		names := &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: name}, &nodes.String{Str: attr},
		}}
		rv2 := makeRangeVarFromNames(names)
		rv2.Loc = nodes.Loc{Start: loc, End: p.pos()}
		return p.finishRelationTable(rv2), nil
	}

	// Simple relation name
	names := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
	rv := makeRangeVarFromNames(names)

	// Check for '*' (include child tables) — relation_expr: qualified_name '*'
	if p.cur.Type == '*' {
		p.advance()
		rv.Inh = true
	}

	rv.Loc = nodes.Loc{Start: loc, End: p.pos()}
	return p.finishRelationTable(rv), nil
}

// finishFuncTable wraps a parsed function expression as a RangeFunction (func_table)
// and parses opt_ordinality and func_alias_clause.
func (p *Parser) finishFuncTable(funcExpr nodes.Node, startLoc int) nodes.Node {
	rf := &nodes.RangeFunction{
		Functions: &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{funcExpr}}}},
	}
	p.parseOptOrdinality(rf)
	p.parseFuncAliasClause(rf)
	rf.Loc = nodes.Loc{Start: startLoc, End: p.prev.End}
	return rf
}

// finishRelationTable adds opt_alias_clause and optional tablesample_clause to a RangeVar.
func (p *Parser) finishRelationTable(rel *nodes.RangeVar) nodes.Node {
	alias := p.parseOptAliasClause()
	if alias != nil {
		rel.Alias = alias
		rel.Loc.End = alias.Loc.End
	}

	// TABLESAMPLE clause
	if p.cur.Type == TABLESAMPLE {
		ts, _ := p.parseTableSampleClause()
		if ts != nil {
			ts.Relation = rel
			return ts
		}
	}

	return rel
}

// parseOptOrdinality parses optional WITH ORDINALITY suffix for func_table.
//
//	opt_ordinality: WITH_LA ORDINALITY | /* EMPTY */
func (p *Parser) parseOptOrdinality(rf *nodes.RangeFunction) {
	if p.cur.Type == WITH_LA && p.peekNext().Type == ORDINALITY {
		p.advance() // WITH_LA
		p.advance() // ORDINALITY
		rf.Ordinality = true
	}
}

// parseRelationExprWithAlias parses ONLY qualified_name with alias.
func (p *Parser) parseRelationExprWithAlias() (nodes.Node, error) {
	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	alias := p.parseOptAliasClause()
	if rel != nil && alias != nil {
		rel.Alias = alias
		rel.Loc.End = alias.Loc.End
	}
	return rel, nil
}

// ---------------------------------------------------------------------------
// Relation expression

// parseRelationExpr parses a relation expression (table reference).
//
//	relation_expr:
//	    qualified_name
//	    | qualified_name '*'
//	    | ONLY qualified_name
//	    | ONLY '(' qualified_name ')'
func (p *Parser) parseRelationExpr() (*nodes.RangeVar, error) {
	if p.collectMode() {
		p.addRuleCandidate("relation_expr")
		// Also emit ONLY token since it's a valid prefix
		p.addTokenCandidate(ONLY)
		// Emit qualified_name candidates too (identifiers/keywords)
		p.addRuleCandidate("qualified_name")
		return nil, errCollecting
	}
	loc := p.pos()
	if p.cur.Type == ONLY {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			names, err := p.parseQualifiedName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			rv := makeRangeVarFromNames(names)
			rv.Inh = false
			rv.Loc = nodes.Loc{Start: loc, End: p.pos()}
			return rv, nil
		}
		names, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		rv := makeRangeVarFromNames(names)
		rv.Inh = false
		rv.Loc = nodes.Loc{Start: loc, End: p.pos()}
		return rv, nil
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)

	// Check for '*' (include child tables)
	if p.cur.Type == '*' {
		p.advance()
		rv.Inh = true
	}

	rv.Loc = nodes.Loc{Start: loc, End: p.pos()}
	return rv, nil
}

// parseRelationExprOptAlias parses relation_expr with optional alias.
//
//	relation_expr_opt_alias:
//	    relation_expr
//	    | relation_expr ColId
//	    | relation_expr AS ColId
func (p *Parser) parseRelationExprOptAlias() (*nodes.RangeVar, error) {
	rv, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	if rv == nil {
		return nil, nil
	}

	if p.cur.Type == AS {
		aliasLoc := p.pos()
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		rv.Alias = &nodes.Alias{Aliasname: name, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
		rv.Loc.End = p.prev.End
	} else if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		aliasLoc := p.pos()
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		rv.Alias = &nodes.Alias{Aliasname: name, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
		rv.Loc.End = p.prev.End
	}

	return rv, nil
}

// ---------------------------------------------------------------------------
// Alias clauses

// parseOptAliasClause parses an optional alias clause.
//
//	opt_alias_clause:
//	    alias_clause | EMPTY
func (p *Parser) parseOptAliasClause() *nodes.Alias {
	if p.cur.Type == AS {
		aliasLoc := p.pos()
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return nil
		}
		alias := &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			colnames, _ := p.parseNameList()
			p.expect(')')
			alias.Colnames = colnames
		}
		alias.Loc = nodes.Loc{Start: aliasLoc, End: p.prev.End}
		return alias
	}

	// ColId (without AS) - but only if it's not a reserved keyword that starts a clause
	if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		aliasLoc := p.pos()
		name, _ := p.parseColId()
		alias := &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			colnames, _ := p.parseNameList()
			p.expect(')')
			alias.Colnames = colnames
		}
		alias.Loc = nodes.Loc{Start: aliasLoc, End: p.prev.End}
		return alias
	}

	return nil
}

// parseFuncAliasClause parses a function alias clause and applies it to a RangeFunction.
//
// Ref: gram.y func_alias_clause
//
//	func_alias_clause:
//	    alias_clause                              (AS ColId ['(' name_list ')'] | ColId ['(' name_list ')'])
//	    | AS '(' TableFuncElementList ')'         (column defs without alias name)
//	    | AS ColId '(' TableFuncElementList ')'   (alias name + column type defs)
//	    | ColId '(' TableFuncElementList ')'      (alias name + column type defs, no AS)
//	    | /* EMPTY */
func (p *Parser) parseFuncAliasClause(rf *nodes.RangeFunction) {
	if p.cur.Type == AS {
		next := p.peekNext()
		if next.Type == '(' {
			// AS '(' TableFuncElementList ')'
			p.advance() // AS
			p.advance() // (
			colDef, _ := p.parseTableFuncElementList()
			p.expect(')')
			rf.Coldeflist = colDef
			return
		}
		aliasLoc := p.pos() // AS token start
		p.advance() // AS
		name, _ := p.parseColId()
		rf.Alias = &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			p.parseFuncAliasParenContents(rf)
			p.expect(')')
		}
		rf.Alias.Loc = nodes.Loc{Start: aliasLoc, End: p.prev.End}
		return
	}

	if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		aliasLoc := p.pos()
		name, _ := p.parseColId()
		rf.Alias = &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			p.parseFuncAliasParenContents(rf)
			p.expect(')')
		}
		rf.Alias.Loc = nodes.Loc{Start: aliasLoc, End: p.prev.End}
	}
}

// parseFuncAliasParenContents disambiguates between name_list and TableFuncElementList
// inside the parentheses of a func_alias_clause. If the first ColId is followed by
// another identifier/type (not ',' or ')'), it's a TableFuncElementList; otherwise name_list.
func (p *Parser) parseFuncAliasParenContents(rf *nodes.RangeFunction) {
	// Peek: after the first ColId, is the next token ',' or ')' (name_list)
	// or something else (TableFuncElementList with types)?
	next := p.peekNext()
	if next.Type == ',' || next.Type == ')' {
		// name_list: ColId [, ColId ...]
		colnames, _ := p.parseNameList()
		rf.Alias.Colnames = colnames
	} else {
		// TableFuncElementList: ColId Typename [, ColId Typename ...]
		colDef, _ := p.parseTableFuncElementList()
		rf.Coldeflist = colDef
	}
}

// ---------------------------------------------------------------------------
// JOIN parsing

// tryParseJoin attempts to parse a JOIN following a table reference.
// Returns nil if no join is found.
func (p *Parser) tryParseJoin(left nodes.Node) (nodes.Node, error) {
	switch p.cur.Type {
	case CROSS:
		// CROSS JOIN
		p.advance() // consume CROSS
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
			Loc:      nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End},
		}, nil

	case JOIN:
		// [INNER] JOIN
		p.advance() // consume JOIN
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
		}
		if err := p.parseJoinQual(j); err != nil {
			return nil, err
		}
		j.Loc = nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End}
		return j, nil

	case INNER_P:
		p.advance() // consume INNER
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
		}
		if err := p.parseJoinQual(j); err != nil {
			return nil, err
		}
		j.Loc = nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End}
		return j, nil

	case LEFT:
		p.advance() // consume LEFT
		p.parseOptOuter()
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_LEFT,
			Larg:     left,
			Rarg:     right,
		}
		if err := p.parseJoinQual(j); err != nil {
			return nil, err
		}
		j.Loc = nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End}
		return j, nil

	case RIGHT:
		p.advance() // consume RIGHT
		p.parseOptOuter()
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_RIGHT,
			Larg:     left,
			Rarg:     right,
		}
		if err := p.parseJoinQual(j); err != nil {
			return nil, err
		}
		j.Loc = nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End}
		return j, nil

	case FULL:
		p.advance() // consume FULL
		p.parseOptOuter()
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_FULL,
			Larg:     left,
			Rarg:     right,
		}
		if err := p.parseJoinQual(j); err != nil {
			return nil, err
		}
		j.Loc = nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End}
		return j, nil

	case NATURAL:
		p.advance() // consume NATURAL
		jt := nodes.JOIN_INNER
		switch p.cur.Type {
		case LEFT:
			p.advance()
			p.parseOptOuter()
			jt = nodes.JOIN_LEFT
		case RIGHT:
			p.advance()
			p.parseOptOuter()
			jt = nodes.JOIN_RIGHT
		case FULL:
			p.advance()
			p.parseOptOuter()
			jt = nodes.JOIN_FULL
		case INNER_P:
			p.advance()
		}
		if _, err := p.expect(JOIN); err != nil {
			return nil, err
		}
		right, err := p.parseTableRefPrimary()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.JoinExpr{
			Jointype:  jt,
			IsNatural: true,
			Larg:      left,
			Rarg:      right,
			Loc:       nodes.Loc{Start: nodes.NodeLoc(left).Start, End: p.prev.End},
		}, nil

	default:
		return nil, nil
	}
}

// parseOptOuter consumes optional OUTER keyword.
func (p *Parser) parseOptOuter() {
	if p.cur.Type == OUTER_P {
		p.advance()
	}
}

// parseJoinQual parses ON condition or USING clause for a join.
//
//	join_qual:
//	    USING '(' name_list ')' join_using_alias
//	    | ON a_expr
func (p *Parser) parseJoinQual(j *nodes.JoinExpr) error {
	if p.cur.Type == USING {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return err
		}
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil
		}
		names, err := p.parseNameList()
		if err != nil {
			return err
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
		j.UsingClause = names

		// join_using_alias: AS ColId | EMPTY
		if p.cur.Type == AS {
			aliasLoc := p.pos()
			p.advance()
			aliasName, err := p.parseColId()
			if err != nil {
				return err
			}
			j.JoinUsing = &nodes.Alias{Aliasname: aliasName, Loc: nodes.Loc{Start: aliasLoc, End: p.prev.End}}
		}
	} else if p.cur.Type == ON {
		p.advance()
		var err error
		j.Quals, err = p.parseAExpr(0)
		if err != nil {
			return err
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// TABLESAMPLE clause

// parseTableSampleClause parses TABLESAMPLE method (args) [REPEATABLE (seed)].
//
//	tablesample_clause:
//	    TABLESAMPLE func_name '(' expr_list ')' opt_repeatable_clause
func (p *Parser) parseTableSampleClause() (*nodes.RangeTableSample, error) {
	loc := p.pos()
	p.advance() // consume TABLESAMPLE

	method, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	args, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	ts := &nodes.RangeTableSample{
		Method: method,
		Args:   args,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// opt_repeatable_clause
	if p.cur.Type == REPEATABLE {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		ts.Repeatable, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	ts.Loc.End = p.pos()
	return ts, nil
}

// ---------------------------------------------------------------------------
// GROUP BY clause

// parseGroupByList parses a GROUP BY list.
//
//	group_by_list:
//	    group_by_item [',' group_by_item ...]
func (p *Parser) parseGroupByList() (*nodes.List, error) {
	first, err := p.parseGroupByItem()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && first == nil {
		return nil, p.syntaxErrorAtCur()
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		item, err := p.parseGroupByItem()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && item == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

// parseGroupByItem parses a single GROUP BY item.
//
//	group_by_item:
//	    a_expr
//	    | empty_grouping_set
//	    | cube_clause
//	    | rollup_clause
//	    | grouping_sets_clause
func (p *Parser) parseGroupByItem() (nodes.Node, error) {
	switch p.cur.Type {
	case '(':
		// empty_grouping_set: '(' ')'
		next := p.peekNext()
		if next.Type == ')' {
			p.advance() // (
			p.advance() // )
			return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_EMPTY, Loc: nodes.NoLoc()}, nil
		}
		return p.parseAExpr(0)

	case CUBE:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		content, err := p.parseExprListFull()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_CUBE, Content: content, Loc: nodes.NoLoc()}, nil

	case ROLLUP:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		content, err := p.parseExprListFull()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_ROLLUP, Content: content, Loc: nodes.NoLoc()}, nil

	case GROUPING:
		// Check for GROUPING SETS
		next := p.peekNext()
		if next.Type == SETS {
			p.advance() // GROUPING
			p.advance() // SETS
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			content, err := p.parseGroupByList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_SETS, Content: content, Loc: nodes.NoLoc()}, nil
		}
		// Not GROUPING SETS; parse as expression (could be GROUPING(...) function)
		return p.parseAExpr(0)

	default:
		return p.parseAExpr(0)
	}
}

// ---------------------------------------------------------------------------
// WHERE clause helpers

// parseWhereClause parses an optional WHERE clause.
func (p *Parser) parseWhereClause() (nodes.Node, error) {
	if p.cur.Type != WHERE {
		return nil, nil
	}
	p.advance()
	return p.parseAExpr(0)
}

// parseWhereOrCurrentClause parses WHERE clause including WHERE CURRENT OF.
//
//	where_or_current_clause:
//	    WHERE a_expr
//	    | WHERE CURRENT_P OF cursor_name
//	    | /* EMPTY */
func (p *Parser) parseWhereOrCurrentClause() (nodes.Node, error) {
	if p.cur.Type != WHERE {
		return nil, nil
	}
	p.advance()

	if p.cur.Type == CURRENT_P {
		next := p.peekNext()
		if next.Type == OF {
			currentLoc := p.pos()
			p.advance() // CURRENT
			p.advance() // OF
			name, err := p.parseCursorName()
			if err != nil {
				return nil, err
			}
			return &nodes.CurrentOfExpr{
				CursorName: name,
				Loc:        nodes.Loc{Start: currentLoc, End: p.prev.End},
			}, nil
		}
	}

	return p.parseAExpr(0)
}
