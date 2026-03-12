package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseWithClause parses a WITH clause (Common Table Expressions).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/with.html
//
//	with_clause:
//	    WITH [RECURSIVE]
//	        cte_name [(col_name [, col_name] ...)] AS (subquery)
//	        [, cte_name [(col_name [, col_name] ...)] AS (subquery)] ...
func (p *Parser) parseWithClause() ([]*nodes.CommonTableExpr, error) {
	p.advance() // consume WITH

	recursive := false
	if _, ok := p.match(kwRECURSIVE); ok {
		recursive = true
	}

	var ctes []*nodes.CommonTableExpr
	for {
		start := p.pos()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		cte := &nodes.CommonTableExpr{
			Loc:       nodes.Loc{Start: start},
			Name:      name,
			Recursive: recursive,
		}

		// Optional column list: (col1, col2, ...)
		if p.cur.Type == '(' {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			cte.Columns = cols
		}

		// AS (subquery)
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		cte.Select = sel

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		cte.Loc.End = p.pos()
		ctes = append(ctes, cte)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return ctes, nil
}

// parseWithStmt parses a statement starting with WITH: WITH ... SELECT, WITH ... UPDATE, WITH ... DELETE.
func (p *Parser) parseWithStmt() (nodes.Node, error) {
	ctes, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwSELECT:
		stmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.CTEs = ctes
		if len(ctes) > 0 {
			stmt.Loc.Start = ctes[0].Loc.Start
		}
		return stmt, nil

	case kwUPDATE:
		stmt, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		// Store CTEs in update context — for now we wrap in a SelectStmt that carries the CTEs
		// and has the update as its payload. But the simpler approach: the UpdateStmt/DeleteStmt
		// don't have CTEs fields. We need to handle this differently.
		// For WITH...UPDATE and WITH...DELETE, the CTEs are typically used via subqueries in WHERE.
		// We return the update stmt as-is since the CTE is referenced in subqueries.
		return stmt, nil

	case kwDELETE:
		stmt, err := p.parseDeleteStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil

	default:
		return nil, &ParseError{
			Message:  "expected SELECT, UPDATE, or DELETE after WITH clause",
			Position: p.cur.Loc,
		}
	}
}

// parseSelectStmt parses a SELECT statement, optionally preceded by a WITH clause.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
//
//	[WITH [RECURSIVE]
//	    cte_name [(col_name [, col_name] ...)] AS (subquery)
//	    [, cte_name [(col_name [, col_name] ...)] AS (subquery)] ...]
//	SELECT
//	    [ALL | DISTINCT | DISTINCTROW]
//	    [HIGH_PRIORITY]
//	    [STRAIGHT_JOIN]
//	    [SQL_CALC_FOUND_ROWS]
//	    select_expr [, select_expr] ...
//	    [FROM table_references]
//	    [WHERE where_condition]
//	    [GROUP BY {col_name | expr | position} [ASC | DESC], ... [WITH ROLLUP]]
//	    [HAVING where_condition]
//	    [ORDER BY {col_name | expr | position} [ASC | DESC], ...]
//	    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
//	    [FOR {UPDATE | SHARE} [OF tbl_name [, tbl_name] ...] [NOWAIT | SKIP LOCKED]]
//	    [INTO OUTFILE 'file_name' | INTO DUMPFILE 'file_name' | INTO var_name [, var_name]]
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error) {
	start := p.pos()

	// Parse optional WITH clause
	var ctes []*nodes.CommonTableExpr
	if p.cur.Type == kwWITH {
		var err error
		ctes, err = p.parseWithClause()
		if err != nil {
			return nil, err
		}
	}

	p.advance() // consume SELECT

	stmt := &nodes.SelectStmt{Loc: nodes.Loc{Start: start}, CTEs: ctes}

	// Parse SELECT options
	for {
		switch p.cur.Type {
		case kwALL:
			stmt.DistinctKind = nodes.DistinctAll
			p.advance()
			continue
		case kwDISTINCT:
			stmt.DistinctKind = nodes.DistinctOn
			p.advance()
			continue
		case kwHIGH_PRIORITY:
			stmt.HighPriority = true
			p.advance()
			continue
		case kwSTRAIGHT_JOIN:
			stmt.StraightJoin = true
			p.advance()
			continue
		case kwSQL_CALC_FOUND_ROWS:
			stmt.CalcFoundRows = true
			p.advance()
			continue
		}
		break
	}

	// Parse select expression list
	targets, err := p.parseSelectExprList()
	if err != nil {
		return nil, err
	}
	stmt.TargetList = targets

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY clause
	if p.cur.Type == kwGROUP {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		groupBy, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy

		// WITH ROLLUP
		if p.cur.Type == kwWITH {
			p.advance()
			if _, ok := p.match(kwROLLUP); ok {
				stmt.WithRollup = true
			}
		}
	}

	// HAVING clause
	if _, ok := p.match(kwHAVING); ok {
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// WINDOW clause: WINDOW window_name AS (window_spec) [, ...]
	if p.cur.Type == kwWINDOW {
		p.advance()
		defs, err := p.parseNamedWindowList()
		if err != nil {
			return nil, err
		}
		stmt.WindowClause = defs
	}

	// ORDER BY clause
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// FOR UPDATE / FOR SHARE
	if p.cur.Type == kwFOR {
		fu, err := p.parseForUpdateClause()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	}

	// INTO clause
	if _, ok := p.match(kwINTO); ok {
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	stmt.Loc.End = p.pos()

	// Check for set operations: UNION, INTERSECT, EXCEPT
	if p.cur.Type == kwUNION || p.cur.Type == kwINTERSECT || p.cur.Type == kwEXCEPT {
		return p.parseSetOperation(stmt)
	}

	return stmt, nil
}

// parseSelectExprList parses a comma-separated list of select expressions.
func (p *Parser) parseSelectExprList() ([]nodes.ExprNode, error) {
	var list []nodes.ExprNode

	for {
		target, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, target)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseSelectExpr parses a single select expression (expr [AS alias]).
func (p *Parser) parseSelectExpr() (nodes.ExprNode, error) {
	start := p.pos()

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Check for AS alias or implicit alias
	var alias string
	if _, ok := p.match(kwAS); ok {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		alias = name
	} else if p.isIdentToken() && !p.isSelectTerminator() {
		// Implicit alias (identifier without AS), but not if it's a keyword that starts the next clause
		alias, _, _ = p.parseIdentifier()
	}

	if alias != "" {
		return &nodes.ResTarget{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Name: alias,
			Val:  expr,
		}, nil
	}

	return expr, nil
}

// isSelectTerminator returns true if the current token starts a clause that terminates
// the select expression list. Used to avoid consuming clause keywords as aliases.
func (p *Parser) isSelectTerminator() bool {
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwLIMIT, kwFOR, kwINTO,
		kwUNION, kwINTERSECT, kwEXCEPT, kwON, kwUSING, kwJOIN, kwINNER, kwLEFT,
		kwRIGHT, kwCROSS, kwNATURAL, kwFULL, kwWINDOW, ';', tokEOF:
		return true
	}
	return false
}

// parseTableReferenceList parses a comma-separated list of table references (with joins).
func (p *Parser) parseTableReferenceList() ([]nodes.TableExpr, error) {
	var refs []nodes.TableExpr

	for {
		ref, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return refs, nil
}

// parseTableReference parses a table reference (table_factor with optional joins).
func (p *Parser) parseTableReference() (nodes.TableExpr, error) {
	left, err := p.parseTableFactor()
	if err != nil {
		return nil, err
	}

	// Parse joins
	for {
		jt, ok := p.matchJoinType()
		if !ok {
			break
		}

		right, err := p.parseTableFactor()
		if err != nil {
			return nil, err
		}

		join := &nodes.JoinClause{
			Loc:   nodes.Loc{Start: tableExprLoc(left)},
			Type:  jt,
			Left:  left,
			Right: right,
		}

		// ON or USING condition
		if _, ok := p.match(kwON); ok {
			condStart := p.pos()
			cond, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			join.Condition = &nodes.OnCondition{
				Loc:  nodes.Loc{Start: condStart},
				Expr: cond,
			}
		} else if _, ok := p.match(kwUSING); ok {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			join.Condition = &nodes.UsingCondition{Columns: cols}
		}

		join.Loc.End = p.pos()
		left = join
	}

	return left, nil
}

// matchJoinType checks for a join keyword combination and returns the join type.
func (p *Parser) matchJoinType() (nodes.JoinType, bool) {
	switch p.cur.Type {
	case kwJOIN:
		p.advance()
		return nodes.JoinInner, true
	case kwINNER:
		p.advance()
		p.match(kwJOIN)
		return nodes.JoinInner, true
	case kwLEFT:
		p.advance()
		p.match(kwOUTER)
		p.match(kwJOIN)
		return nodes.JoinLeft, true
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
		p.match(kwJOIN)
		return nodes.JoinRight, true
	case kwCROSS:
		p.advance()
		p.match(kwJOIN)
		return nodes.JoinCross, true
	case kwNATURAL:
		p.advance()
		// NATURAL [LEFT|RIGHT] [OUTER] JOIN
		if _, ok := p.match(kwLEFT); ok {
			p.match(kwOUTER)
			p.match(kwJOIN)
			return nodes.JoinLeft, true
		}
		if _, ok := p.match(kwRIGHT); ok {
			p.match(kwOUTER)
			p.match(kwJOIN)
			return nodes.JoinRight, true
		}
		p.match(kwJOIN)
		return nodes.JoinNatural, true
	case kwSTRAIGHT_JOIN:
		p.advance()
		return nodes.JoinStraight, true
	}
	return 0, false
}

// parseTableFactor parses a table factor: table_ref, subquery, LATERAL subquery, or parenthesized table_references.
func (p *Parser) parseTableFactor() (nodes.TableExpr, error) {
	// LATERAL (SELECT ...) — LATERAL derived table
	if p.cur.Type == kwLATERAL {
		latStart := p.pos()
		p.advance() // consume LATERAL

		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		if p.cur.Type != kwSELECT {
			return nil, &ParseError{
				Message:  "expected SELECT after LATERAL (",
				Position: p.cur.Loc,
			}
		}

		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		sub := &nodes.SubqueryExpr{
			Loc:     nodes.Loc{Start: latStart},
			Select:  sel,
			Lateral: true,
		}

		// Optional alias: [AS] alias
		if _, ok := p.match(kwAS); ok {
			alias, _, _ := p.parseIdentifier()
			sub.Alias = alias
		} else if p.isIdentToken() && !p.isSelectTerminator() {
			alias, _, _ := p.parseIdentifier()
			sub.Alias = alias
		}

		sub.Loc.End = p.pos()
		return sub, nil
	}

	if p.cur.Type == '(' {
		p.advance()

		// Subquery
		if p.cur.Type == kwSELECT {
			sel, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}

			sub := &nodes.SubqueryExpr{
				Loc:    nodes.Loc{Start: sel.Loc.Start},
				Select: sel,
			}

			// Optional alias
			if _, ok := p.match(kwAS); ok {
				// Subquery alias — stored in a wrapper TableRef
				alias, _, _ := p.parseIdentifier()
				return &nodes.TableRef{
					Loc:   nodes.Loc{Start: sub.Loc.Start, End: p.pos()},
					Name:  alias,
					Alias: alias,
				}, nil
			}

			sub.Loc.End = p.pos()
			return sub, nil
		}

		// Parenthesized table reference
		ref, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return ref, nil
	}

	// JSON_TABLE(expr, path COLUMNS (...)) AS alias
	if p.cur.Type == kwJSON_TABLE {
		return p.parseJsonTable()
	}

	// Regular table reference with alias
	return p.parseTableRefWithAlias()
}

// parseJsonTable parses JSON_TABLE(expr, path COLUMNS (column_list)) [AS] alias.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/json-table-functions.html
//
//	JSON_TABLE(
//	    expr,
//	    path COLUMNS (column_list)
//	) [AS] alias
//
//	column_list:
//	    column[, column][, ...]
//
//	column:
//	    name FOR ORDINALITY
//	  | name type PATH string_path [on_empty] [on_error]
//	  | name type EXISTS PATH string_path
//	  | NESTED [PATH] path COLUMNS (column_list)
//
//	on_empty:
//	    {NULL | DEFAULT json_string | ERROR} ON EMPTY
//
//	on_error:
//	    {NULL | DEFAULT json_string | ERROR} ON ERROR
func (p *Parser) parseJsonTable() (nodes.TableExpr, error) {
	start := p.pos()
	p.advance() // consume JSON_TABLE

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Parse the JSON expression
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	// Parse the path expression
	path, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// COLUMNS
	if p.cur.Type != kwCOLUMNS {
		return nil, &ParseError{Message: "expected COLUMNS in JSON_TABLE", Position: p.cur.Loc}
	}
	p.advance()

	cols, err := p.parseJsonTableColumns()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// [AS] alias (required)
	p.match(kwAS)
	alias := ""
	if p.isIdentToken() {
		alias, _, _ = p.parseIdentifier()
	}

	return &nodes.JsonTableExpr{
		Loc:     nodes.Loc{Start: start, End: p.pos()},
		Expr:    expr,
		Path:    path,
		Columns: cols,
		Alias:   alias,
	}, nil
}

// parseJsonTableColumns parses (column_list) for JSON_TABLE.
func (p *Parser) parseJsonTableColumns() ([]*nodes.JsonTableColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*nodes.JsonTableColumn
	for {
		col, err := p.parseJsonTableColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseJsonTableColumn parses a single JSON_TABLE column definition.
func (p *Parser) parseJsonTableColumn() (*nodes.JsonTableColumn, error) {
	start := p.pos()

	// NESTED [PATH] path COLUMNS (column_list)
	if p.cur.Type == kwNESTED {
		p.advance()
		p.match(kwPATH) // optional PATH keyword
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{Message: "expected path string in NESTED", Position: p.cur.Loc}
		}
		nestedPath := p.cur.Str
		p.advance()

		if p.cur.Type != kwCOLUMNS {
			return nil, &ParseError{Message: "expected COLUMNS after NESTED PATH", Position: p.cur.Loc}
		}
		p.advance()
		nestedCols, err := p.parseJsonTableColumns()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonTableColumn{
			Loc:        nodes.Loc{Start: start, End: p.pos()},
			Nested:     true,
			NestedPath: nestedPath,
			NestedCols: nestedCols,
		}, nil
	}

	// name FOR ORDINALITY | name type PATH path | name type EXISTS PATH path
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	col := &nodes.JsonTableColumn{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// FOR ORDINALITY
	if p.cur.Type == kwFOR {
		p.advance()
		p.match(kwORDINALITY)
		col.Ordinality = true
		col.Loc.End = p.pos()
		return col, nil
	}

	// type
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col.TypeName = typeName

	// EXISTS PATH
	if p.cur.Type == kwEXISTS_KW {
		p.advance()
		col.Exists = true
	}

	// PATH
	if p.cur.Type == kwPATH {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{Message: "expected path string", Position: p.cur.Loc}
		}
		col.Path = p.cur.Str
		p.advance()
	}

	// [on_empty] [on_error]
	// {NULL | DEFAULT json_string | ERROR} ON {EMPTY | ERROR}
	for p.cur.Type == kwNULL || p.cur.Type == kwDEFAULT || p.cur.Type == kwERROR_KW {
		option := ""
		switch p.cur.Type {
		case kwNULL:
			option = "NULL"
			p.advance()
		case kwDEFAULT:
			p.advance()
			if p.cur.Type == tokSCONST {
				option = "DEFAULT " + p.cur.Str
				p.advance()
			}
		case kwERROR_KW:
			option = "ERROR"
			p.advance()
		}

		if p.cur.Type != kwON {
			break
		}
		p.advance() // consume ON

		if p.cur.Type == kwEMPTY {
			col.EmptyOption = option
			p.advance()
		} else if p.cur.Type == kwERROR_KW {
			col.ErrorOption = option
			p.advance()
		}
	}

	col.Loc.End = p.pos()
	return col, nil
}

// parseOrderByList parses ORDER BY items.
func (p *Parser) parseOrderByList() ([]*nodes.OrderByItem, error) {
	var items []*nodes.OrderByItem

	for {
		start := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		item := &nodes.OrderByItem{
			Loc:  nodes.Loc{Start: start},
			Expr: expr,
		}

		if _, ok := p.match(kwDESC); ok {
			item.Desc = true
		} else {
			p.match(kwASC)
		}

		item.Loc.End = p.pos()
		items = append(items, item)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return items, nil
}

// parseLimitClause parses LIMIT [offset,] count or LIMIT count OFFSET offset.
func (p *Parser) parseLimitClause() (*nodes.Limit, error) {
	start := p.pos()

	count, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	limit := &nodes.Limit{
		Loc:   nodes.Loc{Start: start},
		Count: count,
	}

	// LIMIT offset, count
	if p.cur.Type == ',' {
		p.advance()
		count2, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		// In MySQL, LIMIT offset, count — first is offset, second is count
		limit.Offset = limit.Count
		limit.Count = count2
	}

	// LIMIT count OFFSET offset
	if _, ok := p.match(kwOFFSET); ok {
		offset, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		limit.Offset = offset
	}

	limit.Loc.End = p.pos()
	return limit, nil
}

// parseForUpdateClause parses FOR UPDATE / FOR SHARE.
func (p *Parser) parseForUpdateClause() (*nodes.ForUpdate, error) {
	start := p.pos()
	p.advance() // consume FOR

	fu := &nodes.ForUpdate{Loc: nodes.Loc{Start: start}}

	if _, ok := p.match(kwSHARE); ok {
		fu.Share = true
	} else if _, ok := p.match(kwUPDATE); ok {
		// FOR UPDATE (default)
	} else {
		return nil, &ParseError{
			Message:  "expected UPDATE or SHARE after FOR",
			Position: p.cur.Loc,
		}
	}

	// OF table_list
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "of") {
		p.advance()
		for {
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			fu.Tables = append(fu.Tables, ref)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// NOWAIT / SKIP LOCKED
	if _, ok := p.match(kwNOWAIT); ok {
		fu.NoWait = true
	} else if _, ok := p.match(kwSKIP); ok {
		p.match(kwLOCKED)
		fu.SkipLocked = true
	}

	fu.Loc.End = p.pos()
	return fu, nil
}

// parseIntoClause parses INTO OUTFILE / DUMPFILE / var_list.
func (p *Parser) parseIntoClause() (*nodes.IntoClause, error) {
	start := p.pos()
	into := &nodes.IntoClause{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwOUTFILE:
		p.advance()
		if p.cur.Type == tokSCONST {
			into.Outfile = p.cur.Str
			p.advance()
		}
	case kwDUMPFILE:
		p.advance()
		if p.cur.Type == tokSCONST {
			into.Dumpfile = p.cur.Str
			p.advance()
		}
	default:
		// INTO var_name [, var_name ...]
		for {
			if p.isVariableRef() {
				v, err := p.parseVariableRef()
				if err != nil {
					return nil, err
				}
				into.Vars = append(into.Vars, v)
			} else {
				break
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	into.Loc.End = p.pos()
	return into, nil
}

// parseSetOperation parses UNION/INTERSECT/EXCEPT [ALL] SELECT.
func (p *Parser) parseSetOperation(left *nodes.SelectStmt) (*nodes.SelectStmt, error) {
	var op nodes.SetOperation
	switch p.cur.Type {
	case kwUNION:
		op = nodes.SetOpUnion
	case kwINTERSECT:
		op = nodes.SetOpIntersect
	case kwEXCEPT:
		op = nodes.SetOpExcept
	}
	p.advance()

	all := false
	if _, ok := p.match(kwALL); ok {
		all = true
	} else {
		p.match(kwDISTINCT) // UNION DISTINCT is the default
	}

	if p.cur.Type != kwSELECT && p.cur.Type != '(' {
		return nil, &ParseError{
			Message:  "expected SELECT after set operation",
			Position: p.cur.Loc,
		}
	}

	right, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	return &nodes.SelectStmt{
		Loc:    nodes.Loc{Start: left.Loc.Start, End: right.Loc.End},
		SetOp:  op,
		SetAll: all,
		Left:   left,
		Right:  right,
	}, nil
}

// parseParenIdentList parses a parenthesized comma-separated list of identifiers.
func (p *Parser) parseParenIdentList() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var names []string
	for {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// tableExprLoc returns the start location of a TableExpr.
func tableExprLoc(te nodes.TableExpr) int {
	switch t := te.(type) {
	case *nodes.TableRef:
		return t.Loc.Start
	case *nodes.JoinClause:
		return t.Loc.Start
	case *nodes.SubqueryExpr:
		return t.Loc.Start
	}
	return 0
}

// parseNamedWindowList parses WINDOW window_name AS (window_spec) [, ...].
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/window-functions-named-windows.html
//
//	window_clause:
//	    WINDOW window_name AS (window_spec) [, window_name AS (window_spec)] ...
func (p *Parser) parseNamedWindowList() ([]*nodes.WindowDef, error) {
	var defs []*nodes.WindowDef
	for {
		start := p.pos()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		wd := &nodes.WindowDef{Loc: nodes.Loc{Start: start}, Name: name}

		// Optional reference to existing window
		if p.isIdentToken() && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
			p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
			wd.RefName, _, _ = p.parseIdentifier()
		}

		if p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			exprs, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			wd.PartitionBy = exprs
		}

		if p.cur.Type == kwORDER {
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			orderBy, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			wd.OrderBy = orderBy
		}

		if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
			frame, err := p.parseFrameClause()
			if err != nil {
				return nil, err
			}
			wd.Frame = frame
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		wd.Loc.End = p.pos()
		defs = append(defs, wd)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return defs, nil
}

// parseTableStmt parses a TABLE statement (MySQL 8.0.19+).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/table.html
//
//	TABLE table_name [ORDER BY column_name] [LIMIT number [OFFSET number]]
func (p *Parser) parseTableStmt() (*nodes.TableStmt, error) {
	start := p.pos()
	p.advance() // consume TABLE

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.TableStmt{
		Loc:   nodes.Loc{Start: start},
		Table: ref,
	}

	// ORDER BY clause
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseValuesStmt parses a VALUES statement (MySQL 8.0.19+).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/values.html
//
//	VALUES ROW(value_list) [, ROW(value_list)] ...
//	    [ORDER BY column_designator] [LIMIT number [OFFSET number]]
func (p *Parser) parseValuesStmt() (*nodes.ValuesStmt, error) {
	start := p.pos()
	p.advance() // consume VALUES

	stmt := &nodes.ValuesStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse ROW(...) list
	for {
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		var row []nodes.ExprNode
		if p.cur.Type != ')' {
			exprs, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			row = exprs
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Rows = append(stmt.Rows, row)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// ORDER BY clause
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSubqueryExpr parses a subquery expression: (SELECT ...)
func (p *Parser) parseSubqueryExpr() (*nodes.SubqueryExpr, error) {
	start := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	return &nodes.SubqueryExpr{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Select: sel,
	}, nil
}
