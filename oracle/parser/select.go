package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseSelectStmt parses a SELECT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	[ with_clause ]
//	SELECT [ hint ] [ ALL | DISTINCT | UNIQUE ]
//	    select_list
//	    [ FROM table_reference [, ...] ]
//	    [ WHERE condition ]
//	    [ hierarchical_query_clause ]
//	    [ GROUP BY expr [, ...] ]
//	    [ HAVING condition ]
//	    [ ORDER BY sort_key [, ...] ]
//	    [ FOR UPDATE ... ]
//	    [ { UNION [ALL] | INTERSECT | MINUS } select ]
//	    [ OFFSET n { ROW | ROWS } ]
//	    [ FETCH { FIRST | NEXT } n { ROW | ROWS } { ONLY | WITH TIES } ]
func (p *Parser) parseSelectStmt() *nodes.SelectStmt {
	start := p.pos()
	sel := &nodes.SelectStmt{
		TargetList: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// WITH clause
	if p.cur.Type == kwWITH {
		sel.WithClause = p.parseWithClause()
	}

	if p.cur.Type != kwSELECT {
		sel.Loc.End = p.pos()
		return sel
	}
	p.advance() // consume SELECT

	// Hints
	if p.cur.Type == tokHINT {
		sel.Hints = &nodes.List{}
		sel.Hints.Items = append(sel.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// ALL | DISTINCT | UNIQUE
	switch p.cur.Type {
	case kwALL:
		sel.All = true
		p.advance()
	case kwDISTINCT:
		sel.Distinct = true
		p.advance()
	case kwUNIQUE:
		sel.UniqueKw = true
		sel.Distinct = true
		p.advance()
	}

	// Select list
	sel.TargetList = p.parseSelectList()

	// FROM
	if p.cur.Type == kwFROM {
		p.advance()
		sel.FromClause = p.parseFromClause()
	}

	// PIVOT / UNPIVOT (parsed after FROM, before WHERE)
	if p.cur.Type == kwPIVOT {
		sel.Pivot = p.parsePivotClause()
	} else if p.cur.Type == kwUNPIVOT {
		sel.Unpivot = p.parseUnpivotClause()
	}

	// WHERE
	if p.cur.Type == kwWHERE {
		p.advance()
		sel.WhereClause = p.parseExpr()
	}

	// START WITH / CONNECT BY (either order)
	if p.cur.Type == kwSTART || p.cur.Type == kwCONNECT {
		sel.Hierarchical = p.parseHierarchicalClause()
	}

	// GROUP BY
	if p.cur.Type == kwGROUP {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		sel.GroupClause = p.parseExprList()
	}

	// HAVING
	if p.cur.Type == kwHAVING {
		p.advance()
		sel.HavingClause = p.parseExpr()
	}

	// MODEL clause
	if p.cur.Type == kwMODEL {
		sel.ModelClause = p.parseModelClause()
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		sel.OrderBy = p.parseOrderByList()
	}

	// FOR UPDATE
	if p.cur.Type == kwFOR {
		sel.ForUpdate = p.parseForUpdateClause()
	}

	// OFFSET / FETCH FIRST
	if p.cur.Type == kwOFFSET || p.cur.Type == kwFETCH {
		sel.FetchFirst = p.parseFetchFirstClause()
	}

	// Set operations: UNION, INTERSECT, MINUS
	switch p.cur.Type {
	case kwUNION:
		p.advance()
		sel.Op = nodes.SETOP_UNION
		if p.cur.Type == kwALL {
			sel.SetAll = true
			p.advance()
		}
		sel.Rarg = p.parseSelectStmt()
	case kwINTERSECT:
		p.advance()
		sel.Op = nodes.SETOP_INTERSECT
		if p.cur.Type == kwALL {
			sel.SetAll = true
			p.advance()
		}
		sel.Rarg = p.parseSelectStmt()
	case kwMINUS:
		p.advance()
		sel.Op = nodes.SETOP_MINUS
		sel.Rarg = p.parseSelectStmt()
	}

	sel.Loc.End = p.pos()
	return sel
}

// parseSelectList parses the select list (target expressions).
func (p *Parser) parseSelectList() *nodes.List {
	list := &nodes.List{}

	for {
		rt := p.parseResTarget()
		if rt != nil {
			list.Items = append(list.Items, rt)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return list
}

// parseResTarget parses a single target expression (with optional alias).
func (p *Parser) parseResTarget() *nodes.ResTarget {
	start := p.pos()
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	rt := &nodes.ResTarget{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// Optional alias: AS name or just name
	if p.cur.Type == kwAS {
		p.advance()
		rt.Name = p.parseIdentifier()
	} else if p.isAliasCandidate() {
		rt.Name = p.parseIdentifier()
	}

	rt.Loc.End = p.pos()
	return rt
}

// isAliasCandidate returns true if the current token can be an implicit alias.
// We exclude keywords that start clauses to avoid consuming FROM, WHERE, etc. as aliases.
func (p *Parser) isAliasCandidate() bool {
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	// Disallow clause-starting keywords as implicit aliases
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT,
		kwMINUS, kwFOR, kwCONNECT, kwSTART, kwFETCH, kwOFFSET, kwON,
		kwLEFT, kwRIGHT, kwINNER, kwOUTER, kwCROSS, kwFULL, kwNATURAL,
		kwJOIN, kwWHEN, kwTHEN, kwELSE, kwEND, kwAND, kwOR, kwNOT,
		kwIS, kwIN, kwBETWEEN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4,
		kwINTO, kwVALUES, kwSET, kwRETURNING, kwPIVOT, kwUNPIVOT,
		kwMODEL, kwWITH, kwKEEP, kwOVER:
		return false
	}
	return false
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() *nodes.List {
	list := &nodes.List{}
	for {
		expr := p.parseExpr()
		if expr == nil {
			break
		}
		list.Items = append(list.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseFromClause parses a FROM clause (comma-separated table references).
func (p *Parser) parseFromClause() *nodes.List {
	list := &nodes.List{}

	for {
		tref := p.parseTableRef()
		if tref == nil {
			break
		}

		// Check for JOINs
		tref = p.parseJoinContinuation(tref)

		list.Items = append(list.Items, tref)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list
}

// parseTableRef parses a single table reference.
func (p *Parser) parseTableRef() nodes.TableExpr {
	start := p.pos()

	// Subquery: ( SELECT ... )
	if p.cur.Type == '(' {
		return p.parseSubqueryRef(start)
	}

	// Table name
	if !p.isIdentLike() {
		return nil
	}

	name := p.parseObjectName()
	tr := &nodes.TableRef{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// Flashback query: VERSIONS BETWEEN ... AND ... or AS OF SCN/TIMESTAMP
	if p.cur.Type == kwVERSIONS || (p.cur.Type == kwAS && p.peekNext().Type == kwOF) {
		tr.Flashback = p.parseFlashbackClause()
	}

	// SAMPLE [BLOCK] (percent) [SEED (value)]
	if p.cur.Type == kwSAMPLE {
		tr.Sample = p.parseSampleClause()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		// Only consume AS as alias intro if next is not OF (flashback already handled)
		p.advance()
		tr.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		tr.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	tr.Loc.End = p.pos()
	return tr
}

// isTableAliasCandidate checks if current token can be a table alias.
func (p *Parser) isTableAliasCandidate() bool {
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	return false
}

// parseSubqueryRef parses a subquery in FROM: ( SELECT ... ) alias.
func (p *Parser) parseSubqueryRef(start int) nodes.TableExpr {
	p.advance() // consume '('

	subSel := p.parseSelectStmt()

	if p.cur.Type == ')' {
		p.advance()
	}

	ref := &nodes.SubqueryRef{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start},
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		ref.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		ref.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	ref.Loc.End = p.pos()
	return ref
}

// parseJoinContinuation parses any JOIN clauses that follow a table reference.
func (p *Parser) parseJoinContinuation(left nodes.TableExpr) nodes.TableExpr {
	for {
		jt, ok := p.matchJoinType()
		if !ok {
			break
		}

		right := p.parseTableRef()

		jc := &nodes.JoinClause{
			Type:  jt,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: p.pos()},
		}

		// ON condition
		if p.cur.Type == kwON {
			p.advance()
			jc.On = p.parseExpr()
		}

		// USING ( col1, col2, ... )
		if p.cur.Type == kwUSING {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				jc.Using = &nodes.List{}
				for {
					name := p.parseIdentifier()
					if name != "" {
						jc.Using.Items = append(jc.Using.Items, &nodes.String{Str: name})
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}

		jc.Loc.End = p.pos()
		left = jc
	}
	return left
}

// matchJoinType tries to match a JOIN keyword sequence.
// Returns the JoinType and true if matched, or false if no JOIN found.
func (p *Parser) matchJoinType() (nodes.JoinType, bool) {
	natural := false
	if p.cur.Type == kwNATURAL {
		natural = true
		p.advance()
	}

	switch p.cur.Type {
	case kwJOIN:
		p.advance()
		if natural {
			return nodes.JOIN_NATURAL_INNER, true
		}
		return nodes.JOIN_INNER, true

	case kwINNER:
		p.advance()
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_INNER, true
		}
		return nodes.JOIN_INNER, true

	case kwLEFT:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_LEFT, true
		}
		return nodes.JOIN_LEFT, true

	case kwRIGHT:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_RIGHT, true
		}
		return nodes.JOIN_RIGHT, true

	case kwFULL:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_FULL, true
		}
		return nodes.JOIN_FULL, true

	case kwCROSS:
		p.advance()
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		return nodes.JOIN_CROSS, true
	}

	if natural {
		// NATURAL without a recognized join keyword — treat as NATURAL INNER JOIN
		return nodes.JOIN_NATURAL_INNER, true
	}

	return 0, false
}

// parseHierarchicalClause parses START WITH / CONNECT BY clauses (either order).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Hierarchical-Queries.html
//
//	[ START WITH condition ] CONNECT BY [ NOCYCLE ] condition
//	CONNECT BY [ NOCYCLE ] condition [ START WITH condition ]
func (p *Parser) parseHierarchicalClause() *nodes.HierarchicalClause {
	start := p.pos()
	hc := &nodes.HierarchicalClause{
		Loc: nodes.Loc{Start: start},
	}

	// START WITH first, then CONNECT BY
	if p.cur.Type == kwSTART {
		p.advance() // START
		if p.cur.Type == kwWITH {
			p.advance() // WITH
		}
		hc.StartWith = p.parseExpr()
	}

	if p.cur.Type == kwCONNECT {
		p.advance() // CONNECT
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.isIdentLikeStr("NOCYCLE") {
			hc.IsNocycle = true
			p.advance()
		}
		hc.ConnectBy = p.parseExpr()
	}

	// START WITH may come after CONNECT BY
	if hc.StartWith == nil && p.cur.Type == kwSTART {
		p.advance()
		if p.cur.Type == kwWITH {
			p.advance()
		}
		hc.StartWith = p.parseExpr()
	}

	hc.Loc.End = p.pos()
	return hc
}

// parseOrderByList parses a comma-separated list of ORDER BY items.
func (p *Parser) parseOrderByList() *nodes.List {
	list := &nodes.List{}

	for {
		sb := p.parseSortBy()
		if sb == nil {
			break
		}
		list.Items = append(list.Items, sb)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list
}

// parseSortBy parses a single ORDER BY item.
//
//	sort_key [ ASC | DESC ] [ NULLS { FIRST | LAST } ]
func (p *Parser) parseSortBy() *nodes.SortBy {
	start := p.pos()
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	sb := &nodes.SortBy{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		sb.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		sb.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			sb.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			sb.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	sb.Loc.End = p.pos()
	return sb
}

// parseForUpdateClause parses FOR UPDATE [OF ...] [NOWAIT | WAIT n | SKIP LOCKED].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html#GUID-CFA006CA-6FF1-4972-821E-6996142A51C6
func (p *Parser) parseForUpdateClause() *nodes.ForUpdateClause {
	start := p.pos()
	p.advance() // consume FOR

	if p.cur.Type == kwUPDATE {
		p.advance()
	}

	fu := &nodes.ForUpdateClause{
		Loc: nodes.Loc{Start: start},
	}

	// OF table_list
	if p.cur.Type == kwOF {
		p.advance()
		fu.Tables = p.parseExprList()
	}

	// NOWAIT | WAIT n | SKIP LOCKED
	switch p.cur.Type {
	case kwNOWAIT:
		fu.NoWait = true
		p.advance()
	case kwWAIT:
		p.advance()
		fu.Wait = p.parseExpr()
	case kwSKIP:
		p.advance()
		if p.cur.Type == kwLOCKED {
			fu.SkipLocked = true
			p.advance()
		}
	}

	fu.Loc.End = p.pos()
	return fu
}

// parseFetchFirstClause parses OFFSET/FETCH FIRST clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	[ OFFSET n { ROW | ROWS } ]
//	FETCH { FIRST | NEXT } [ n [ PERCENT ] ] { ROW | ROWS } { ONLY | WITH TIES }
func (p *Parser) parseFetchFirstClause() *nodes.FetchFirstClause {
	start := p.pos()
	fc := &nodes.FetchFirstClause{
		Loc: nodes.Loc{Start: start},
	}

	// OFFSET n ROWS
	if p.cur.Type == kwOFFSET {
		p.advance()
		fc.Offset = p.parseExpr()
		// ROW | ROWS
		if p.cur.Type == kwROW || p.cur.Type == kwROWS {
			p.advance()
		}
	}

	// FETCH FIRST|NEXT
	if p.cur.Type == kwFETCH {
		p.advance()
		// FIRST | NEXT
		if p.cur.Type == kwFIRST || p.cur.Type == kwNEXT {
			p.advance()
		}

		// count expression
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			fc.Count = p.parseExpr()
		}

		// PERCENT
		if p.cur.Type == kwPERCENT {
			fc.Percent = true
			p.advance()
		}

		// ROW | ROWS
		if p.cur.Type == kwROW || p.cur.Type == kwROWS {
			p.advance()
		}

		// ONLY | WITH TIES
		if p.cur.Type == kwONLY {
			p.advance()
		} else if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == kwTIES {
				fc.WithTies = true
				p.advance()
			}
		}
	}

	fc.Loc.End = p.pos()
	return fc
}

// parseWithClause parses a WITH clause (common table expressions).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	WITH [ RECURSIVE ] cte_name [ ( col1, col2, ... ) ] AS ( subquery ) [, ...]
func (p *Parser) parseWithClause() *nodes.WithClause {
	start := p.pos()
	p.advance() // consume WITH

	wc := &nodes.WithClause{
		CTEs: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	if p.cur.Type == kwRECURSIVE {
		wc.Recursive = true
		p.advance()
	}

	for {
		cte := p.parseCTE()
		if cte == nil {
			break
		}
		wc.CTEs.Items = append(wc.CTEs.Items, cte)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	wc.Loc.End = p.pos()
	return wc
}

// parseCTE parses a single common table expression.
func (p *Parser) parseCTE() *nodes.CTE {
	if !p.isIdentLike() {
		return nil
	}

	start := p.pos()
	name := p.parseIdentifier()

	cte := &nodes.CTE{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		cte.Columns = &nodes.List{}
		for {
			col := p.parseIdentifier()
			if col != "" {
				cte.Columns.Items = append(cte.Columns.Items, &nodes.String{Str: col})
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// AS
	if p.cur.Type == kwAS {
		p.advance()
	}

	// ( subquery )
	if p.cur.Type == '(' {
		p.advance()
		cte.Query = p.parseSelectStmt()
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	cte.Loc.End = p.pos()
	return cte
}

// parsePivotClause parses a PIVOT clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	pivot_clause ::=
//	    PIVOT (
//	        aggregate_function ( expr ) [ [ AS ] c_alias ]
//	        [, aggregate_function ( expr ) [ [ AS ] c_alias ] ] ...
//	        FOR { column | ( column [, column ] ... ) }
//	        IN ( { { expr | ( expr [, expr] ... ) } [ [ AS ] c_alias ] } [, ...] )
//	    )
func (p *Parser) parsePivotClause() *nodes.PivotClause {
	start := p.pos()
	p.advance() // consume PIVOT

	pc := &nodes.PivotClause{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		pc.Loc.End = p.pos()
		return pc
	}
	p.advance() // consume '('

	// Parse aggregate function list (before FOR keyword)
	pc.AggFuncs = &nodes.List{}
	for {
		agg := p.parseResTarget()
		if agg != nil {
			pc.AggFuncs.Items = append(pc.AggFuncs.Items, agg)
		}
		if p.cur.Type != ',' {
			break
		}
		// Peek ahead to see if this comma separates aggregates (before FOR)
		// or if we've reached the FOR keyword
		if p.cur.Type == kwFOR {
			break
		}
		p.advance() // consume ','
		// If the next token is FOR, we went too far — this shouldn't happen
		// because FOR is not an expression start, but be safe
		if p.cur.Type == kwFOR {
			break
		}
	}

	// FOR column | ( column, column, ... )
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if p.cur.Type == '(' {
			// Multi-column: ( col1, col2, ... )
			p.advance() // consume '('
			colList := &nodes.List{}
			for {
				col := p.parseColumnRef()
				if col != nil {
					colList.Items = append(colList.Items, col)
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
			if len(colList.Items) == 1 {
				if e, ok := colList.Items[0].(nodes.ExprNode); ok {
					pc.ForCol = e
				}
			} else {
				pc.ForCols = colList
			}
		} else {
			// Single column reference (not a full expression, to avoid consuming IN)
			pc.ForCol = p.parseColumnRef()
		}
	}

	// IN ( ... )
	if p.cur.Type == kwIN {
		p.advance() // consume IN
		if p.cur.Type == '(' {
			p.advance() // consume '('
			pc.InList = p.parsePivotInList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance() // consume outer ')'
	}

	pc.Loc.End = p.pos()
	return pc
}

// parsePivotInList parses the IN list of a PIVOT clause.
// Each item is: expr [ [ AS ] c_alias ] | ( expr, expr, ... ) [ [ AS ] c_alias ]
func (p *Parser) parsePivotInList() *nodes.List {
	list := &nodes.List{}
	for {
		start := p.pos()
		var expr nodes.ExprNode

		if p.cur.Type == '(' {
			// Tuple value: ( expr, expr, ... )
			p.advance() // consume '('
			tupleList := p.parseExprList()
			if p.cur.Type == ')' {
				p.advance()
			}
			// Wrap tuple as a ParenExpr containing the first item
			// For multi-value, store as a List via a special representation
			if tupleList.Len() == 1 {
				if e, ok := tupleList.Items[0].(nodes.ExprNode); ok {
					expr = &nodes.ParenExpr{Expr: e, Loc: nodes.Loc{Start: start, End: p.pos()}}
				}
			} else {
				// Use a FuncCallExpr with empty name to represent a row/tuple
				args := &nodes.List{Items: tupleList.Items}
				expr = &nodes.FuncCallExpr{
					FuncName: &nodes.ObjectName{Name: "", Loc: nodes.Loc{Start: start}},
					Args:     args,
					Loc:      nodes.Loc{Start: start, End: p.pos()},
				}
			}
		} else {
			expr = p.parseExpr()
		}

		if expr == nil {
			break
		}

		rt := &nodes.ResTarget{
			Expr: expr,
			Loc:  nodes.Loc{Start: start},
		}

		// Optional alias: [ AS ] c_alias
		if p.cur.Type == kwAS {
			p.advance()
			rt.Name = p.parseIdentifier()
		} else if p.isAliasCandidate() {
			rt.Name = p.parseIdentifier()
		}

		rt.Loc.End = p.pos()
		list.Items = append(list.Items, rt)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseUnpivotClause parses an UNPIVOT clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	unpivot_clause ::=
//	    UNPIVOT [ { INCLUDE | EXCLUDE } NULLS ]
//	    (
//	        column
//	        FOR column
//	        IN ( column [ [ AS ] literal ] [, column [ [ AS ] literal ] ] ... )
//	    )
func (p *Parser) parseUnpivotClause() *nodes.UnpivotClause {
	start := p.pos()
	p.advance() // consume UNPIVOT

	uc := &nodes.UnpivotClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ INCLUDE | EXCLUDE ] NULLS
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type == kwNULLS {
			p.advance()
		}
		uc.IncludeNulls = true
	} else if p.isIdentLikeStr("EXCLUDE") {
		p.advance()
		if p.cur.Type == kwNULLS {
			p.advance()
		}
		// IncludeNulls stays false (EXCLUDE is the default)
	}

	if p.cur.Type != '(' {
		uc.Loc.End = p.pos()
		return uc
	}
	p.advance() // consume '('

	// Value column(s)
	uc.ValueCol = p.parseColumnRef()

	// FOR pivot_column(s)
	if p.cur.Type == kwFOR {
		p.advance()
		uc.PivotCol = p.parseColumnRef()
	}

	// IN ( column [ AS literal ], ... )
	if p.cur.Type == kwIN {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			uc.InList = &nodes.List{}
			for {
				start := p.pos()
				col := p.parseColumnRef()
				if col == nil || col.Column == "" {
					break
				}
				rt := &nodes.ResTarget{
					Expr: col,
					Loc:  nodes.Loc{Start: start},
				}
				// Optional AS alias (can be identifier or string literal)
				if p.cur.Type == kwAS {
					p.advance()
					if p.cur.Type == tokSCONST {
						rt.Name = p.cur.Str
						p.advance()
					} else {
						rt.Name = p.parseIdentifier()
					}
				}
				rt.Loc.End = p.pos()
				uc.InList.Items = append(uc.InList.Items, rt)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance() // consume outer ')'
	}

	uc.Loc.End = p.pos()
	return uc
}

// parseModelClause parses an Oracle MODEL clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	model_clause ::=
//	    MODEL
//	      [ cell_reference_options ]
//	      [ return_rows_clause ]
//	      [ reference_model ]...
//	      main_model
//
//	cell_reference_options ::=
//	    { IGNORE NAV | KEEP NAV }
//	    { UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE }
//
//	return_rows_clause ::=
//	    RETURN { UPDATED | ALL } ROWS
//
//	reference_model ::=
//	    REFERENCE reference_model_name ON ( subquery )
//	        model_column_clauses
//	        [ cell_reference_options ]
//
//	main_model ::=
//	    [ MAIN main_model_name ]
//	        model_column_clauses
//	        [ cell_reference_options ]
//	        model_rules_clause
//
//	model_column_clauses ::=
//	    [ PARTITION BY ( expr [ [ AS ] c_alias ] [, ...] ) ]
//	    DIMENSION BY ( expr [ [ AS ] c_alias ] [, ...] )
//	    MEASURES ( expr [ [ AS ] c_alias ] [, ...] )
//
//	model_rules_clause ::=
//	    [ RULES ]
//	    [ { UPDATE | UPSERT [ ALL ] } ]
//	    [ { AUTOMATIC | SEQUENTIAL } ORDER ]
//	    [ model_iterate_clause ]
//	    ( cell_assignment [, ...] )
//
//	model_iterate_clause ::=
//	    ITERATE ( number ) [ UNTIL ( condition ) ]
//
//	cell_assignment ::=
//	    measure_column [ dimension_subscripts ] = expr
//
//	single_column_for_loop ::=
//	    FOR dimension_column
//	      { IN ( { literal [, ...] | subquery } )
//	      | [ LIKE pattern ] FROM literal TO literal { INCREMENT | DECREMENT } literal
//	      }
//
//	multi_column_for_loop ::=
//	    FOR ( dimension_column [, ...] ) IN
//	      ( ( literal [, ...] ) [, ( literal [, ...] ) ]...
//	      | subquery
//	      )
func (p *Parser) parseModelClause() *nodes.ModelClause {
	start := p.pos()
	p.advance() // consume MODEL

	mc := &nodes.ModelClause{
		Loc: nodes.Loc{Start: start},
	}

	// cell_reference_options (optional, before RETURN/REFERENCE/MAIN/PARTITION/DIMENSION)
	mc.CellRefOptions = p.parseModelCellRefOptions()

	// RETURN { UPDATED | ALL } ROWS
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if p.cur.Type == kwUPDATED {
			mc.ReturnRows = "UPDATED"
			p.advance()
		} else if p.cur.Type == kwALL {
			mc.ReturnRows = "ALL"
			p.advance()
		}
		if p.cur.Type == kwROWS {
			p.advance()
		}
	}

	// REFERENCE models (zero or more)
	for p.cur.Type == kwREFERENCE {
		ref := p.parseModelRefModel()
		mc.RefModels = append(mc.RefModels, ref)
	}

	// main_model
	mc.MainModel = p.parseModelMainModel()

	mc.Loc.End = p.pos()
	return mc
}

// parseModelCellRefOptions parses optional cell_reference_options.
func (p *Parser) parseModelCellRefOptions() *nodes.ModelCellRefOptions {
	opts := &nodes.ModelCellRefOptions{Loc: nodes.Loc{Start: p.pos()}}
	found := false

	// IGNORE NAV | KEEP NAV
	if p.isIdentLikeStr("IGNORE") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.IgnoreNav = true
		found = true
	} else if p.isIdentLikeStr("KEEP") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.KeepNav = true
		found = true
	}

	// UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE
	if p.cur.Type == kwUNIQUE {
		p.advance()
		if p.cur.Type == kwDIMENSION {
			opts.UniqueDimension = true
			p.advance()
			found = true
		} else if p.isIdentLikeStr("SINGLE") {
			p.advance()
			if p.cur.Type == kwREFERENCE {
				p.advance()
			}
			opts.UniqueSingleRef = true
			found = true
		}
	}

	if !found {
		return nil
	}

	opts.Loc.End = p.pos()
	return opts
}

// parseModelRefModel parses a REFERENCE model.
func (p *Parser) parseModelRefModel() *nodes.ModelRefModel {
	start := p.pos()
	p.advance() // consume REFERENCE

	ref := &nodes.ModelRefModel{
		Loc: nodes.Loc{Start: start},
	}

	// reference_model_name
	ref.Name = p.parseIdentifier()

	// ON ( subquery )
	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			ref.Subquery = p.parseSelectStmt()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// model_column_clauses
	ref.ColumnClauses = p.parseModelColumnClauses()

	// optional cell_reference_options
	ref.CellRefOptions = p.parseModelCellRefOptions()

	ref.Loc.End = p.pos()
	return ref
}

// parseModelMainModel parses the main_model.
func (p *Parser) parseModelMainModel() *nodes.ModelMainModel {
	start := p.pos()
	mm := &nodes.ModelMainModel{
		Loc: nodes.Loc{Start: start},
	}

	// [ MAIN main_model_name ]
	if p.cur.Type == kwMAIN {
		p.advance()
		mm.Name = p.parseIdentifier()
	}

	// model_column_clauses
	mm.ColumnClauses = p.parseModelColumnClauses()

	// optional cell_reference_options
	mm.CellRefOptions = p.parseModelCellRefOptions()

	// model_rules_clause
	mm.RulesClause = p.parseModelRulesClause()

	mm.Loc.End = p.pos()
	return mm
}

// parseModelColumnClauses parses PARTITION BY / DIMENSION BY / MEASURES.
func (p *Parser) parseModelColumnClauses() *nodes.ModelColumnClauses {
	start := p.pos()
	cc := &nodes.ModelColumnClauses{
		Loc: nodes.Loc{Start: start},
	}

	// [ PARTITION BY ( expr [AS alias], ... ) ]
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			cc.PartitionBy = p.parseModelColumnList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// DIMENSION BY ( expr [AS alias], ... )
	if p.cur.Type == kwDIMENSION {
		p.advance() // DIMENSION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			cc.DimensionBy = p.parseModelColumnList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// MEASURES ( expr [AS alias], ... )
	if p.cur.Type == kwMEASURES {
		p.advance() // MEASURES
		if p.cur.Type == '(' {
			p.advance()
			cc.Measures = p.parseModelColumnList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	cc.Loc.End = p.pos()
	return cc
}

// parseModelColumnList parses a comma-separated list of expr [ [AS] alias ] for MODEL clauses.
func (p *Parser) parseModelColumnList() *nodes.List {
	list := &nodes.List{}
	for {
		rt := p.parseResTarget()
		if rt == nil {
			break
		}
		list.Items = append(list.Items, rt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseModelRulesClause parses the model_rules_clause.
func (p *Parser) parseModelRulesClause() *nodes.ModelRulesClause {
	start := p.pos()
	rc := &nodes.ModelRulesClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ RULES ]
	if p.cur.Type == kwRULES {
		p.advance()
	}

	// [ UPDATE | UPSERT [ALL] ]
	if p.cur.Type == kwUPDATE {
		rc.UpdateMode = "UPDATE"
		p.advance()
	} else if p.cur.Type == kwUPSERT {
		p.advance()
		if p.cur.Type == kwALL {
			rc.UpdateMode = "UPSERT ALL"
			p.advance()
		} else {
			rc.UpdateMode = "UPSERT"
		}
	}

	// [ AUTOMATIC | SEQUENTIAL ] ORDER
	if p.cur.Type == kwAUTOMATIC {
		rc.OrderMode = "AUTOMATIC"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	} else if p.cur.Type == kwSEQUENTIAL {
		rc.OrderMode = "SEQUENTIAL"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	}

	// [ ITERATE ( number ) [ UNTIL ( condition ) ] ]
	if p.cur.Type == kwITERATE {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			rc.Iterate = p.parseExpr()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// [ UNTIL ( condition ) ]
		if p.cur.Type == kwUNTIL {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				rc.Until = p.parseExpr()
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	// ( cell_assignment [, ...] )
	if p.cur.Type == '(' {
		p.advance()
		rc.Rules = &nodes.List{}
		for {
			rule := p.parseModelRule()
			if rule == nil {
				break
			}
			rc.Rules.Items = append(rc.Rules.Items, rule)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	rc.Loc.End = p.pos()
	return rc
}

// parseModelRule parses a single cell_assignment: measure_column[dim1, dim2] = expr.
func (p *Parser) parseModelRule() *nodes.ModelRule {
	start := p.pos()

	// Parse the left side: measure_column [ dim_subscripts ]
	// The left side is an identifier (possibly qualified) followed by [ subscripts ]
	cellRef := p.parseModelCellRef()
	if cellRef == nil {
		return nil
	}

	rule := &nodes.ModelRule{
		CellRef: cellRef,
		Loc:     nodes.Loc{Start: start},
	}

	// = expr
	if p.cur.Type == '=' {
		p.advance()
		rule.Expr = p.parseExpr()
	}

	rule.Loc.End = p.pos()
	return rule
}

// parseModelCellRef parses a model cell reference: measure[dim1, dim2, ...].
// Dimension subscripts can be expressions, ANY, or FOR loops (single/multi column).
func (p *Parser) parseModelCellRef() nodes.ExprNode {
	start := p.pos()

	// Parse the measure column name as a column reference
	if !p.isIdentLike() {
		return nil
	}
	name := p.parseIdentifier()
	col := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: start, End: p.pos()},
	}

	// [ dim_subscripts ]
	if p.cur.Type != '[' {
		return col
	}
	p.advance() // consume '['

	// Parse subscript expressions (comma-separated)
	args := &nodes.List{}
	for {
		if p.cur.Type == ']' {
			break
		}
		expr := p.parseExpr()
		if expr == nil {
			break
		}
		args.Items = append(args.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ']' {
		p.advance()
	}

	// Represent as a FuncCallExpr with the column name and subscript args
	return &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start}},
		Args:     args,
		Loc:      nodes.Loc{Start: start, End: p.pos()},
	}
}

// parseSampleClause parses a SAMPLE clause on a table reference.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	sample_clause ::=
//	    SAMPLE [ BLOCK ] ( sample_percent )
//	    [ SEED ( seed_value ) ]
func (p *Parser) parseSampleClause() *nodes.SampleClause {
	start := p.pos()
	p.advance() // consume SAMPLE

	sc := &nodes.SampleClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ BLOCK ]
	if p.cur.Type == kwBLOCK {
		sc.Block = true
		p.advance()
	}

	// ( percent )
	if p.cur.Type == '(' {
		p.advance()
		sc.Percent = p.parseExpr()
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// [ SEED ( value ) ]
	if p.cur.Type == kwSEED {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			sc.Seed = p.parseExpr()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	sc.Loc.End = p.pos()
	return sc
}

// parseFlashbackClause parses a flashback query clause on a table reference.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	flashback_query_clause ::=
//	    { VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
//	    | AS OF { SCN | TIMESTAMP } expr
//	    }
func (p *Parser) parseFlashbackClause() *nodes.FlashbackClause {
	start := p.pos()
	fc := &nodes.FlashbackClause{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwVERSIONS {
		// VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
		fc.IsVersions = true
		p.advance() // consume VERSIONS

		if p.cur.Type == kwBETWEEN {
			p.advance() // consume BETWEEN
		}

		// SCN | TIMESTAMP
		if p.cur.Type == kwSCN {
			fc.Type = "SCN"
			p.advance()
		} else if p.cur.Type == kwTIMESTAMP {
			fc.Type = "TIMESTAMP"
			p.advance()
		}

		// low expr — parse above AND precedence to not consume the BETWEEN...AND
		fc.VersionsLow = p.parseExprPrec(precNot)

		// AND
		if p.cur.Type == kwAND {
			p.advance()
		}

		// high expr
		fc.VersionsHigh = p.parseExpr()

	} else if p.cur.Type == kwAS {
		// AS OF { SCN | TIMESTAMP } expr
		p.advance() // consume AS
		if p.cur.Type == kwOF {
			p.advance() // consume OF
		}

		// SCN | TIMESTAMP
		if p.cur.Type == kwSCN {
			fc.Type = "SCN"
			p.advance()
		} else if p.cur.Type == kwTIMESTAMP {
			fc.Type = "TIMESTAMP"
			p.advance()
		}

		// expr
		fc.Expr = p.parseExpr()
	}

	fc.Loc.End = p.pos()
	return fc
}
