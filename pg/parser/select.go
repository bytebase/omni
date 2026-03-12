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
func (p *Parser) parseSelectNoParens() *nodes.SelectStmt {
	// 1. Optional WITH clause
	var withClause *nodes.WithClause
	if p.cur.Type == WITH || p.cur.Type == WITH_LA {
		withClause = p.parseWithClause()
	}

	// 2. Parse select_clause (handles UNION/INTERSECT/EXCEPT)
	stmt := p.parseSelectClause(setOpPrecNone)
	if stmt == nil {
		return nil
	}

	// 3. Optional ORDER BY
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		stmt.SortClause = p.parseSortByList()
	}

	// 4. Parse LIMIT/OFFSET and FOR locking in either order
	p.parseSelectOptions(stmt)

	if withClause != nil {
		stmt.WithClause = withClause
	}

	return stmt
}

// parseSelectWithParens parses a parenthesized SELECT statement.
//
//	select_with_parens:
//	    '(' select_no_parens ')'
//	    | '(' select_with_parens ')'
func (p *Parser) parseSelectWithParens() *nodes.SelectStmt {
	p.expect('(')
	var stmt *nodes.SelectStmt
	if p.cur.Type == '(' {
		stmt = p.parseSelectWithParens()
	} else {
		stmt = p.parseSelectNoParens()
	}
	p.expect(')')
	return stmt
}

// parseSelectClause parses a select_clause using precedence climbing
// for set operations (UNION/INTERSECT/EXCEPT).
//
//	select_clause:
//	    simple_select
//	    | select_with_parens
func (p *Parser) parseSelectClause(minPrec int) *nodes.SelectStmt {
	left := p.parseSelectClausePrimary()
	if left == nil {
		return nil
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
			return left
		}
		if prec < minPrec {
			return left
		}
		p.advance() // consume UNION/INTERSECT/EXCEPT
		all := p.parseSetQuantifier()
		right := p.parseSelectClause(prec + 1) // left-associative
		result := &nodes.SelectStmt{
			Op:   op,
			Larg: left,
			Rarg: right,
		}
		if all == nodes.SET_QUANTIFIER_ALL {
			result.All = true
		}
		left = result
	}
}

// parseSelectClausePrimary parses a primary select clause (leaf or parenthesized).
func (p *Parser) parseSelectClausePrimary() *nodes.SelectStmt {
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
func (p *Parser) parseSimpleSelectLeaf() *nodes.SelectStmt {
	switch p.cur.Type {
	case SELECT:
		return p.parseSimpleSelectCore()
	case VALUES:
		return p.parseValuesClause()
	case TABLE:
		return p.parseTableCmd()
	default:
		return nil
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
func (p *Parser) parseSimpleSelectCore() *nodes.SelectStmt {
	p.advance() // consume SELECT
	stmt := &nodes.SelectStmt{}

	// ALL / DISTINCT
	if p.cur.Type == ALL {
		p.advance()
	} else if p.cur.Type == DISTINCT {
		p.advance()
		if p.cur.Type == ON {
			p.advance()
			p.expect('(')
			stmt.DistinctClause = p.parseExprListFull()
			p.expect(')')
		} else {
			stmt.DistinctClause = &nodes.List{Items: []nodes.Node{nil}}
		}
	}

	// target_list (opt_target_list)
	stmt.TargetList = p.parseTargetList()

	// INTO clause
	if p.cur.Type == INTO {
		stmt.IntoClause = p.parseIntoClause()
	}

	// FROM clause
	if p.cur.Type == FROM {
		p.advance()
		stmt.FromClause = p.parseFromListFull()
	}

	// WHERE clause
	if p.cur.Type == WHERE {
		p.advance()
		stmt.WhereClause = p.parseAExpr(0)
	}

	// GROUP BY clause
	if p.cur.Type == GROUP_P {
		p.advance()
		p.expect(BY)
		// Check for ALL or DISTINCT
		if p.cur.Type == ALL {
			p.advance()
			// GROUP BY ALL is treated the same as GROUP BY (all is default)
		} else if p.cur.Type == DISTINCT {
			p.advance()
			stmt.GroupDistinct = true
		}
		stmt.GroupClause = p.parseGroupByList()
	}

	// HAVING clause
	if p.cur.Type == HAVING {
		p.advance()
		stmt.HavingClause = p.parseAExpr(0)
	}

	// WINDOW clause
	stmt.WindowClause = p.parseWindowClause()

	return stmt
}

// parseValuesClause parses a VALUES clause.
//
//	values_clause:
//	    VALUES '(' expr_list ')' [',' '(' expr_list ')' ...]
func (p *Parser) parseValuesClause() *nodes.SelectStmt {
	p.advance() // consume VALUES
	stmt := &nodes.SelectStmt{}

	p.expect('(')
	first := p.parseExprListFull()
	p.expect(')')
	stmt.ValuesLists = &nodes.List{Items: []nodes.Node{first}}

	for p.cur.Type == ',' {
		p.advance()
		p.expect('(')
		exprs := p.parseExprListFull()
		p.expect(')')
		stmt.ValuesLists.Items = append(stmt.ValuesLists.Items, exprs)
	}

	return stmt
}

// parseTableCmd parses TABLE relation_expr (shorthand for SELECT * FROM ...).
//
//	TABLE relation_expr
func (p *Parser) parseTableCmd() *nodes.SelectStmt {
	p.advance() // consume TABLE
	rel := p.parseRelationExpr()
	cr := &nodes.ColumnRef{
		Fields:   &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
		Location: -1,
	}
	rt := &nodes.ResTarget{
		Val:      cr,
		Location: -1,
	}
	return &nodes.SelectStmt{
		TargetList: &nodes.List{Items: []nodes.Node{rt}},
		FromClause: &nodes.List{Items: []nodes.Node{rel}},
	}
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
func (p *Parser) parseSelectOptions(stmt *nodes.SelectStmt) {
	// Parse LIMIT/OFFSET/FETCH and FOR locking in either order
	for {
		if p.cur.Type == LIMIT || p.cur.Type == OFFSET || p.cur.Type == FETCH {
			p.parseSelectLimit(stmt)
		} else if p.cur.Type == FOR {
			locking := p.parseForLockingClause()
			if locking != nil {
				stmt.LockingClause = locking
			}
		} else {
			break
		}
	}
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
func (p *Parser) parseSelectLimit(stmt *nodes.SelectStmt) {
	for {
		switch p.cur.Type {
		case LIMIT:
			p.advance()
			if p.cur.Type == ALL {
				p.advance()
				stmt.LimitCount = &nodes.A_Const{Isnull: true}
				stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
			} else {
				stmt.LimitCount = p.parseAExpr(0)
				stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
				// Check for deprecated LIMIT #,# syntax
				if p.cur.Type == ',' {
					p.advance()
					// In MySQL-compatible syntax: LIMIT offset, count
					// Swap: what we parsed as count is actually offset
					stmt.LimitOffset = stmt.LimitCount
					stmt.LimitCount = p.parseAExpr(0)
				}
			}
		case OFFSET:
			p.advance()
			// Could be: OFFSET select_offset_value
			// or: OFFSET select_fetch_first_value row_or_rows
			expr := p.parseAExpr(0)
			// Check if followed by ROW or ROWS (SQL:2008 syntax)
			if p.cur.Type == ROW || p.cur.Type == ROWS {
				p.advance()
			}
			stmt.LimitOffset = expr
		case FETCH:
			p.parseFetchClause(stmt)
		default:
			return
		}
	}
}

// parseFetchClause parses FETCH FIRST/NEXT ... ROWS ONLY/WITH TIES.
//
//	FETCH first_or_next select_fetch_first_value row_or_rows ONLY
//	| FETCH first_or_next select_fetch_first_value row_or_rows WITH TIES
//	| FETCH first_or_next row_or_rows ONLY
//	| FETCH first_or_next row_or_rows WITH TIES
func (p *Parser) parseFetchClause(stmt *nodes.SelectStmt) {
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
		if p.cur.Type == '+' {
			p.advance()
			stmt.LimitCount = p.parseCExpr()
		} else if p.cur.Type == '-' {
			p.advance()
			stmt.LimitCount = doNegate(p.parseCExpr())
		} else {
			stmt.LimitCount = p.parseCExpr()
		}
		// row_or_rows: ROW | ROWS
		if p.cur.Type == ROW || p.cur.Type == ROWS {
			p.advance()
		}
	}

	// ONLY or WITH TIES
	if p.cur.Type == WITH {
		p.advance()
		p.expect(TIES)
		stmt.LimitOption = nodes.LIMIT_OPTION_WITH_TIES
	} else {
		p.expect(ONLY)
		stmt.LimitOption = nodes.LIMIT_OPTION_COUNT
	}
}

// ---------------------------------------------------------------------------
// FOR locking clause

// parseForLockingClause parses FOR UPDATE/SHARE/NO KEY UPDATE/KEY SHARE clauses.
//
//	for_locking_clause:
//	    for_locking_items
//	    | FOR READ ONLY
func (p *Parser) parseForLockingClause() *nodes.List {
	if p.cur.Type != FOR {
		return nil
	}

	// Check for FOR READ ONLY
	next := p.peekNext()
	if next.Type == READ {
		p.advance() // consume FOR
		p.advance() // consume READ
		p.expect(ONLY)
		return nil // FOR READ ONLY = no locking
	}

	var items []nodes.Node
	for p.cur.Type == FOR {
		item := p.parseForLockingItem()
		if item != nil {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseForLockingItem parses a single FOR UPDATE/SHARE item.
//
//	for_locking_item:
//	    for_locking_strength locked_rels_list opt_nowait_or_skip
func (p *Parser) parseForLockingItem() *nodes.LockingClause {
	strength := p.parseForLockingStrength()
	if strength < 0 {
		return nil
	}

	lc := &nodes.LockingClause{
		Strength: int(strength),
	}

	// locked_rels_list: OF qualified_name_list | EMPTY
	if p.cur.Type == OF {
		p.advance()
		rels, _ := p.parseQualifiedNameList()
		lc.LockedRels = rels
	}

	// opt_nowait_or_skip
	if p.cur.Type == NOWAIT {
		p.advance()
		lc.WaitPolicy = int(nodes.LockWaitError)
	} else if p.cur.Type == SKIP {
		p.advance()
		p.expect(LOCKED)
		lc.WaitPolicy = int(nodes.LockWaitSkip)
	} else {
		lc.WaitPolicy = int(nodes.LockWaitBlock)
	}

	return lc
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
func (p *Parser) parseWithClause() *nodes.WithClause {
	loc := p.pos()
	p.advance() // consume WITH or WITH_LA

	wc := &nodes.WithClause{
		Location: nodes.ParseLoc(loc),
	}

	if p.cur.Type == RECURSIVE {
		p.advance()
		wc.Recursive = true
	}

	wc.Ctes = p.parseCTEList()
	return wc
}

// parseCTEList parses a comma-separated list of CTEs.
func (p *Parser) parseCTEList() *nodes.List {
	first := p.parseCommonTableExpr()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseCommonTableExpr())
	}
	return &nodes.List{Items: items}
}

// parseCommonTableExpr parses a single CTE.
//
//	common_table_expr:
//	    name opt_name_list AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
//	    | name '(' name_list ')' AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
func (p *Parser) parseCommonTableExpr() *nodes.CommonTableExpr {
	loc := p.pos()
	name, _ := p.parseName()

	cte := &nodes.CommonTableExpr{
		Ctename:  name,
		Location: nodes.ParseLoc(loc),
	}

	// Check for name '(' name_list ')' form
	if p.cur.Type == '(' {
		p.advance()
		colnames, _ := p.parseNameList()
		p.expect(')')
		cte.Aliascolnames = colnames
	}

	p.expect(AS)

	// opt_materialized
	cte.Ctematerialized = int(p.parseOptMaterialized())

	// '(' PreparableStmt ')'
	p.expect('(')
	cte.Ctequery = p.parseSelectNoParens()
	p.expect(')')

	// opt_search_clause (rarely used, skip for now)
	// opt_cycle_clause (rarely used, skip for now)

	return cte
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

// ---------------------------------------------------------------------------
// INTO clause

// parseIntoClause parses SELECT INTO clause.
//
//	into_clause:
//	    INTO OptTempTableName
func (p *Parser) parseIntoClause() *nodes.IntoClause {
	p.advance() // consume INTO

	rv := p.parseOptTempTableName()
	return &nodes.IntoClause{
		Rel:      rv,
		OnCommit: nodes.ONCOMMIT_NOOP,
	}
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
func (p *Parser) parseOptTempTableName() *nodes.RangeVar {
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

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = persistence
	return rv
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
func (p *Parser) parseFromListFull() *nodes.List {
	first := p.parseTableRef()
	if first == nil {
		return nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		item := p.parseTableRef()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}
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
func (p *Parser) parseTableRef() nodes.Node {
	left := p.parseTableRefPrimary()
	if left == nil {
		return nil
	}

	// Handle joins (left-recursive in yacc, loop in recursive descent)
	for {
		joined := p.tryParseJoin(left)
		if joined == nil {
			break
		}
		left = joined
	}

	return left
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
func (p *Parser) parseTableRefPrimary() nodes.Node {
	switch p.cur.Type {
	case '(':
		return p.parseParenTableRef()
	case LATERAL_P:
		return p.parseLateralTableRef()
	case XMLTABLE:
		n := p.parseXmlTable().(*nodes.RangeTableFunc)
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
		}
		return n
	case JSON_TABLE:
		n := p.parseJsonTable()
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
		}
		return n
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
		return nil
	}
}

// parseParenTableRef handles '(' ... ')' in FROM clause.
// Could be: '(' select_no_parens ')' or '(' joined_table ')' opt_alias
func (p *Parser) parseParenTableRef() nodes.Node {
	// Peek to see if this is a subquery or joined table
	// If next token after '(' is SELECT/VALUES/WITH/TABLE, it's a subquery
	next := p.peekNext()
	if next.Type == SELECT || next.Type == VALUES || next.Type == WITH || next.Type == TABLE || next.Type == '(' {
		// Could be select_with_parens → subquery
		stmt := p.parseSelectWithParens()
		alias := p.parseOptAliasClause()
		return &nodes.RangeSubselect{
			Subquery: stmt,
			Alias:    alias,
		}
	}

	// '(' joined_table ')' opt_alias_clause
	p.advance() // consume '('
	inner := p.parseTableRef()
	p.expect(')')
	alias := p.parseOptAliasClause()
	if j, ok := inner.(*nodes.JoinExpr); ok && alias != nil {
		j.Alias = alias
	}
	return inner
}

// parseLateralTableRef handles LATERAL prefix.
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	LATERAL_P func_table func_alias_clause
//	| LATERAL_P select_with_parens opt_alias_clause
//	| LATERAL_P xmltable opt_alias_clause
//	| LATERAL_P json_table opt_alias_clause
func (p *Parser) parseLateralTableRef() nodes.Node {
	p.advance() // consume LATERAL

	if p.cur.Type == '(' {
		// LATERAL select_with_parens
		stmt := p.parseSelectWithParens()
		alias := p.parseOptAliasClause()
		return &nodes.RangeSubselect{
			Lateral:  true,
			Subquery: stmt,
			Alias:    alias,
		}
	}

	if p.cur.Type == ROWS {
		rf := p.parseRowsFromTable()
		if rf, ok := rf.(*nodes.RangeFunction); ok {
			rf.Lateral = true
		}
		return rf
	}

	// LATERAL xmltable opt_alias_clause
	if p.cur.Type == XMLTABLE {
		n := p.parseXmlTable().(*nodes.RangeTableFunc)
		n.Lateral = true
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
		}
		return n
	}

	// LATERAL json_table opt_alias_clause
	if p.cur.Type == JSON_TABLE {
		n := p.parseJsonTable()
		n.Lateral = true
		alias := p.parseOptAliasClause()
		if alias != nil {
			n.Alias = alias
		}
		return n
	}

	// LATERAL func_table func_alias_clause
	funcExpr := p.parseFuncExprWindowless()
	rf := &nodes.RangeFunction{
		Lateral:   true,
		Functions: &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{funcExpr}}}},
	}
	p.parseOptOrdinality(rf)
	p.parseFuncAliasClause(rf)
	return rf
}

// parseRowsFromTable parses ROWS FROM (...) [WITH ORDINALITY].
func (p *Parser) parseRowsFromTable() nodes.Node {
	p.advance() // consume ROWS
	p.expect(FROM)
	p.expect('(')

	var items []nodes.Node
	for {
		funcExpr := p.parseFuncExprWindowless()
		var colDef *nodes.List
		if p.cur.Type == AS && p.peekNext().Type == '(' {
			p.advance() // AS
			p.advance() // (
			colDef = p.parseTableFuncElementList()
			p.expect(')')
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
	p.expect(')')

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
	return rf
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
func (p *Parser) parseRelationOrFuncTable() nodes.Node {
	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil
	}

	// name( → function call
	if p.cur.Type == '(' {
		funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
		funcExpr := p.parseFuncApplication(funcName, loc)
		return p.finishFuncTable(funcExpr)
	}

	// name.something
	if p.cur.Type == '.' {
		p.advance() // consume '.'
		attr, _ := p.parseAttrName()

		// schema.func(
		if p.cur.Type == '(' {
			funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}}
			funcExpr := p.parseFuncApplication(funcName, loc)
			return p.finishFuncTable(funcExpr)
		}

		// catalog.schema.name
		if p.cur.Type == '.' {
			p.advance()
			attr2, _ := p.parseAttrName()

			// catalog.schema.func(
			if p.cur.Type == '(' {
				funcName := &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2},
				}}
				funcExpr := p.parseFuncApplication(funcName, loc)
				return p.finishFuncTable(funcExpr)
			}

			// 3-part qualified relation name
			names := &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2},
			}}
			return p.finishRelationTable(makeRangeVarFromNames(names))
		}

		// 2-part qualified relation name
		names := &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: name}, &nodes.String{Str: attr},
		}}
		return p.finishRelationTable(makeRangeVarFromNames(names))
	}

	// Simple relation name
	names := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
	rv := makeRangeVarFromNames(names)

	// Check for '*' (include child tables) — relation_expr: qualified_name '*'
	if p.cur.Type == '*' {
		p.advance()
		rv.Inh = true
	}

	return p.finishRelationTable(rv)
}

// finishFuncTable wraps a parsed function expression as a RangeFunction (func_table)
// and parses opt_ordinality and func_alias_clause.
func (p *Parser) finishFuncTable(funcExpr nodes.Node) nodes.Node {
	rf := &nodes.RangeFunction{
		Functions: &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{funcExpr}}}},
	}
	p.parseOptOrdinality(rf)
	p.parseFuncAliasClause(rf)
	return rf
}

// finishRelationTable adds opt_alias_clause and optional tablesample_clause to a RangeVar.
func (p *Parser) finishRelationTable(rel *nodes.RangeVar) nodes.Node {
	alias := p.parseOptAliasClause()
	if alias != nil {
		rel.Alias = alias
	}

	// TABLESAMPLE clause
	if p.cur.Type == TABLESAMPLE {
		ts := p.parseTableSampleClause()
		ts.Relation = rel
		return ts
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
func (p *Parser) parseRelationExprWithAlias() nodes.Node {
	rel := p.parseRelationExpr()
	alias := p.parseOptAliasClause()
	if rel != nil && alias != nil {
		rel.Alias = alias
	}
	return rel
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
func (p *Parser) parseRelationExpr() *nodes.RangeVar {
	if p.cur.Type == ONLY {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			names, _ := p.parseQualifiedName()
			p.expect(')')
			rv := makeRangeVarFromNames(names)
			rv.Inh = false
			return rv
		}
		names, _ := p.parseQualifiedName()
		rv := makeRangeVarFromNames(names)
		rv.Inh = false
		return rv
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil
	}
	rv := makeRangeVarFromNames(names)

	// Check for '*' (include child tables)
	if p.cur.Type == '*' {
		p.advance()
		rv.Inh = true
	}

	return rv
}

// parseRelationExprOptAlias parses relation_expr with optional alias.
//
//	relation_expr_opt_alias:
//	    relation_expr
//	    | relation_expr ColId
//	    | relation_expr AS ColId
func (p *Parser) parseRelationExprOptAlias() *nodes.RangeVar {
	rv := p.parseRelationExpr()
	if rv == nil {
		return nil
	}

	if p.cur.Type == AS {
		p.advance()
		name, _ := p.parseColId()
		rv.Alias = &nodes.Alias{Aliasname: name}
	} else if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		name, _ := p.parseColId()
		rv.Alias = &nodes.Alias{Aliasname: name}
	}

	return rv
}

// ---------------------------------------------------------------------------
// Alias clauses

// parseOptAliasClause parses an optional alias clause.
//
//	opt_alias_clause:
//	    alias_clause | EMPTY
func (p *Parser) parseOptAliasClause() *nodes.Alias {
	if p.cur.Type == AS {
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
		return alias
	}

	// ColId (without AS) - but only if it's not a reserved keyword that starts a clause
	if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		name, _ := p.parseColId()
		alias := &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			colnames, _ := p.parseNameList()
			p.expect(')')
			alias.Colnames = colnames
		}
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
			colDef := p.parseTableFuncElementList()
			p.expect(')')
			rf.Coldeflist = colDef
			return
		}
		p.advance() // AS
		name, _ := p.parseColId()
		rf.Alias = &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			p.parseFuncAliasParenContents(rf)
			p.expect(')')
		}
		return
	}

	if p.isColId() && !p.isReservedForClause() && !p.isJoinKeyword() {
		name, _ := p.parseColId()
		rf.Alias = &nodes.Alias{Aliasname: name}
		if p.cur.Type == '(' {
			p.advance()
			p.parseFuncAliasParenContents(rf)
			p.expect(')')
		}
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
		colDef := p.parseTableFuncElementList()
		rf.Coldeflist = colDef
	}
}

// ---------------------------------------------------------------------------
// JOIN parsing

// tryParseJoin attempts to parse a JOIN following a table reference.
// Returns nil if no join is found.
func (p *Parser) tryParseJoin(left nodes.Node) nodes.Node {
	switch p.cur.Type {
	case CROSS:
		// CROSS JOIN
		p.advance() // consume CROSS
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		return &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
		}

	case JOIN:
		// [INNER] JOIN
		p.advance() // consume JOIN
		right := p.parseTableRefPrimary()
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
		}
		p.parseJoinQual(j)
		return j

	case INNER_P:
		p.advance() // consume INNER
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_INNER,
			Larg:     left,
			Rarg:     right,
		}
		p.parseJoinQual(j)
		return j

	case LEFT:
		p.advance() // consume LEFT
		p.parseOptOuter()
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_LEFT,
			Larg:     left,
			Rarg:     right,
		}
		p.parseJoinQual(j)
		return j

	case RIGHT:
		p.advance() // consume RIGHT
		p.parseOptOuter()
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_RIGHT,
			Larg:     left,
			Rarg:     right,
		}
		p.parseJoinQual(j)
		return j

	case FULL:
		p.advance() // consume FULL
		p.parseOptOuter()
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		j := &nodes.JoinExpr{
			Jointype: nodes.JOIN_FULL,
			Larg:     left,
			Rarg:     right,
		}
		p.parseJoinQual(j)
		return j

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
		p.expect(JOIN)
		right := p.parseTableRefPrimary()
		return &nodes.JoinExpr{
			Jointype:  jt,
			IsNatural: true,
			Larg:      left,
			Rarg:      right,
		}

	default:
		return nil
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
func (p *Parser) parseJoinQual(j *nodes.JoinExpr) {
	if p.cur.Type == USING {
		p.advance()
		p.expect('(')
		names, _ := p.parseNameList()
		p.expect(')')
		j.UsingClause = names

		// join_using_alias: AS ColId | EMPTY
		if p.cur.Type == AS {
			p.advance()
			aliasName, _ := p.parseColId()
			j.JoinUsing = &nodes.Alias{Aliasname: aliasName}
		}
	} else if p.cur.Type == ON {
		p.advance()
		j.Quals = p.parseAExpr(0)
	}
}

// ---------------------------------------------------------------------------
// TABLESAMPLE clause

// parseTableSampleClause parses TABLESAMPLE method (args) [REPEATABLE (seed)].
//
//	tablesample_clause:
//	    TABLESAMPLE func_name '(' expr_list ')' opt_repeatable_clause
func (p *Parser) parseTableSampleClause() *nodes.RangeTableSample {
	loc := p.pos()
	p.advance() // consume TABLESAMPLE

	method, _ := p.parseFuncName()
	p.expect('(')
	args := p.parseExprListFull()
	p.expect(')')

	ts := &nodes.RangeTableSample{
		Method:   method,
		Args:     args,
		Location: nodes.ParseLoc(loc),
	}

	// opt_repeatable_clause
	if p.cur.Type == REPEATABLE {
		p.advance()
		p.expect('(')
		ts.Repeatable = p.parseAExpr(0)
		p.expect(')')
	}

	return ts
}

// ---------------------------------------------------------------------------
// GROUP BY clause

// parseGroupByList parses a GROUP BY list.
//
//	group_by_list:
//	    group_by_item [',' group_by_item ...]
func (p *Parser) parseGroupByList() *nodes.List {
	first := p.parseGroupByItem()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseGroupByItem())
	}
	return &nodes.List{Items: items}
}

// parseGroupByItem parses a single GROUP BY item.
//
//	group_by_item:
//	    a_expr
//	    | empty_grouping_set
//	    | cube_clause
//	    | rollup_clause
//	    | grouping_sets_clause
func (p *Parser) parseGroupByItem() nodes.Node {
	switch p.cur.Type {
	case '(':
		// empty_grouping_set: '(' ')'
		next := p.peekNext()
		if next.Type == ')' {
			p.advance() // (
			p.advance() // )
			return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_EMPTY, Location: -1}
		}
		return p.parseAExpr(0)

	case CUBE:
		p.advance()
		p.expect('(')
		content := p.parseExprListFull()
		p.expect(')')
		return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_CUBE, Content: content, Location: -1}

	case ROLLUP:
		p.advance()
		p.expect('(')
		content := p.parseExprListFull()
		p.expect(')')
		return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_ROLLUP, Content: content, Location: -1}

	case GROUPING:
		// Check for GROUPING SETS
		next := p.peekNext()
		if next.Type == SETS {
			p.advance() // GROUPING
			p.advance() // SETS
			p.expect('(')
			content := p.parseGroupByList()
			p.expect(')')
			return &nodes.GroupingSet{Kind: nodes.GROUPING_SET_SETS, Content: content, Location: -1}
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
func (p *Parser) parseWhereClause() nodes.Node {
	if p.cur.Type != WHERE {
		return nil
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
func (p *Parser) parseWhereOrCurrentClause() nodes.Node {
	if p.cur.Type != WHERE {
		return nil
	}
	p.advance()

	if p.cur.Type == CURRENT_P {
		next := p.peekNext()
		if next.Type == OF {
			p.advance() // CURRENT
			p.advance() // OF
			name, _ := p.parseCursorName()
			return &nodes.CurrentOfExpr{
				CursorName: name,
			}
		}
	}

	return p.parseAExpr(0)
}
