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
	// Record CTE position for completion context (before consuming WITH).
	if p.completing {
		p.addCTEPosition(p.pos())
	}
	p.advance() // consume WITH

	// Completion: after WITH keyword, offer RECURSIVE and identifier context for CTE name.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwRECURSIVE)
		return nil, &ParseError{Message: "collecting"}
	}

	recursive := false
	if _, ok := p.match(kwRECURSIVE); ok {
		recursive = true
	}

	var ctes []*nodes.CommonTableExpr
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
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

		// Completion: inside CTE AS (, offer SELECT keyword.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
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

	case kwINSERT, kwREPLACE:
		isReplace := p.cur.Type == kwREPLACE
		stmt, err := p.parseInsertOrReplace(isReplace)
		if err != nil {
			return nil, err
		}
		return stmt, nil

	case kwUPDATE:
		stmt, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil

	case kwDELETE:
		stmt, err := p.parseDeleteStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil

	default:
		return nil, &ParseError{
			Message:  "expected SELECT, INSERT, UPDATE, or DELETE after WITH clause",
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

	// Completion: after SELECT keyword, offer select-expression candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwDISTINCT)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.SelectStmt{Loc: nodes.Loc{Start: start}, CTEs: ctes}

	// Parse SELECT options
	for {
		switch p.cur.Type {
		case kwALL:
			stmt.DistinctKind = nodes.DistinctAll
			p.advance()
			continue
		case kwDISTINCT, kwDISTINCTROW:
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
		case kwSQL_SMALL_RESULT:
			stmt.SmallResult = true
			p.advance()
			continue
		case kwSQL_BIG_RESULT:
			stmt.BigResult = true
			p.advance()
			continue
		case kwSQL_BUFFER_RESULT:
			stmt.BufferResult = true
			p.advance()
			continue
		case kwSQL_NO_CACHE:
			stmt.NoCache = true
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

	// INTO clause (position 1: after select_expr, before FROM)
	if p.cur.Type == kwINTO {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause — MySQL allows WHERE without FROM (e.g., SELECT 1 WHERE 1=1)
	if p.cur.Type == kwWHERE {
		p.advance()
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where

		// Reject duplicate WHERE
		if p.cur.Type == kwWHERE {
			return nil, &ParseError{
				Message:  "duplicate WHERE clause",
				Position: p.cur.Loc,
			}
		}
	}

	// GROUP BY clause — MySQL allows GROUP BY without FROM
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

		// Completion: after GROUP BY list, offer follow-set keywords.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwHAVING)
			p.addTokenCandidate(kwORDER)
			p.addTokenCandidate(kwLIMIT)
			p.addTokenCandidate(kwWITH)
			return nil, &ParseError{Message: "collecting"}
		}

		// WITH ROLLUP
		if p.cur.Type == kwWITH {
			p.advance()
			if _, ok := p.match(kwROLLUP); ok {
				stmt.WithRollup = true
			}
		}

		// Reject duplicate GROUP BY
		if p.cur.Type == kwGROUP {
			return nil, &ParseError{
				Message:  "duplicate GROUP BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// HAVING clause — MySQL allows HAVING without FROM
	if p.cur.Type == kwHAVING {
		p.advance()
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

	// INTO clause (position 2: after HAVING/WINDOW, before ORDER BY)
	if p.cur.Type == kwINTO && stmt.Into == nil {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
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

		// WITH ROLLUP (MySQL 8.0.12+)
		if p.cur.Type == kwWITH {
			next := p.peekNext()
			if next.Type == kwROLLUP {
				p.advance() // consume WITH
				p.advance() // consume ROLLUP
				stmt.OrderByWithRollup = true
			}
		}

		// Reject duplicate ORDER BY
		if p.cur.Type == kwORDER {
			return nil, &ParseError{
				Message:  "duplicate ORDER BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if p.cur.Type == kwFOR {
		fu, err := p.parseForUpdateClause()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	} else if p.cur.Type == kwLOCK {
		fu, err := p.parseLockInShareMode()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	}

	// INTO clause (position 3: after locking clause)
	if p.cur.Type == kwINTO && stmt.Into == nil {
		p.advance()
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

// parseSelectStmtBase parses a single SELECT (no set operations consumed).
// Used by the set operation precedence parser to get the right-hand operand.
func (p *Parser) parseSelectStmtBase() (*nodes.SelectStmt, error) {
	// Handle parenthesized select: (SELECT ...)
	if p.cur.Type == '(' {
		p.advance()
		inner, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return inner, nil
	}

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

	// Completion: after SELECT keyword, offer select-expression candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwDISTINCT)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.SelectStmt{Loc: nodes.Loc{Start: start}, CTEs: ctes}

	// Parse SELECT options
	for {
		switch p.cur.Type {
		case kwALL:
			stmt.DistinctKind = nodes.DistinctAll
			p.advance()
			continue
		case kwDISTINCT, kwDISTINCTROW:
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
		case kwSQL_SMALL_RESULT:
			stmt.SmallResult = true
			p.advance()
			continue
		case kwSQL_BIG_RESULT:
			stmt.BigResult = true
			p.advance()
			continue
		case kwSQL_BUFFER_RESULT:
			stmt.BufferResult = true
			p.advance()
			continue
		case kwSQL_NO_CACHE:
			stmt.NoCache = true
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

	// INTO clause (position 1: after select_expr, before FROM)
	if p.cur.Type == kwINTO {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause — MySQL allows WHERE without FROM (e.g., SELECT 1 WHERE 1=1)
	if p.cur.Type == kwWHERE {
		p.advance()
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where

		// Reject duplicate WHERE
		if p.cur.Type == kwWHERE {
			return nil, &ParseError{
				Message:  "duplicate WHERE clause",
				Position: p.cur.Loc,
			}
		}
	}

	// GROUP BY clause — MySQL allows GROUP BY without FROM
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

		// Completion: after GROUP BY list, offer follow-set keywords.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwHAVING)
			p.addTokenCandidate(kwORDER)
			p.addTokenCandidate(kwLIMIT)
			p.addTokenCandidate(kwWITH)
			return nil, &ParseError{Message: "collecting"}
		}

		// WITH ROLLUP
		if p.cur.Type == kwWITH {
			p.advance()
			if _, ok := p.match(kwROLLUP); ok {
				stmt.WithRollup = true
			}
		}

		// Reject duplicate GROUP BY
		if p.cur.Type == kwGROUP {
			return nil, &ParseError{
				Message:  "duplicate GROUP BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// HAVING clause — MySQL allows HAVING without FROM
	if p.cur.Type == kwHAVING {
		p.advance()
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// WINDOW clause
	if p.cur.Type == kwWINDOW {
		p.advance()
		defs, err := p.parseNamedWindowList()
		if err != nil {
			return nil, err
		}
		stmt.WindowClause = defs
	}

	// INTO clause (position 2: after HAVING/WINDOW, before ORDER BY)
	if p.cur.Type == kwINTO && stmt.Into == nil {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
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

		// WITH ROLLUP (MySQL 8.0.12+)
		if p.cur.Type == kwWITH {
			next := p.peekNext()
			if next.Type == kwROLLUP {
				p.advance() // consume WITH
				p.advance() // consume ROLLUP
				stmt.OrderByWithRollup = true
			}
		}

		// Reject duplicate ORDER BY
		if p.cur.Type == kwORDER {
			return nil, &ParseError{
				Message:  "duplicate ORDER BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if p.cur.Type == kwFOR {
		fu, err := p.parseForUpdateClause()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	} else if p.cur.Type == kwLOCK {
		fu, err := p.parseLockInShareMode()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	}

	// INTO clause (position 3: after locking clause)
	if p.cur.Type == kwINTO && stmt.Into == nil {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSelectExprList parses a comma-separated list of select expressions.
func (p *Parser) parseSelectExprList() ([]nodes.ExprNode, error) {
	var list []nodes.ExprNode

	for {
		// Completion: at the start of each select expression, offer column/function candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, &ParseError{Message: "collecting"}
		}

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
		aliasLoc := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		alias = name
		if p.completing {
			p.addSelectAliasPosition(aliasLoc)
		}
	} else if p.isIdentToken() && !p.isSelectTerminator() {
		// Implicit alias (identifier without AS), but not if it's a keyword that starts the next clause
		aliasLoc := p.pos()
		alias, _, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		if p.completing {
			p.addSelectAliasPosition(aliasLoc)
		}
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
		kwRIGHT, kwCROSS, kwNATURAL, kwFULL, kwWINDOW, kwWITH, kwLOCK, ';', tokEOF:
		return true
	}
	return false
}

// parseTableReferenceList parses a comma-separated list of table references (with joins).
func (p *Parser) parseTableReferenceList() ([]nodes.TableExpr, error) {
	var refs []nodes.TableExpr

	for {
		// Completion: at the start of each table reference, offer table_ref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
		}

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

	// Completion: after a table factor, offer join keywords and post-FROM keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwWHERE)
		p.addTokenCandidate(kwJOIN)
		p.addTokenCandidate(kwLEFT)
		p.addTokenCandidate(kwRIGHT)
		p.addTokenCandidate(kwCROSS)
		p.addTokenCandidate(kwNATURAL)
		p.addTokenCandidate(kwINNER)
		p.addTokenCandidate(kwSTRAIGHT_JOIN)
		p.addTokenCandidate(kwORDER)
		p.addTokenCandidate(kwGROUP)
		p.addTokenCandidate(kwHAVING)
		p.addTokenCandidate(kwLIMIT)
		p.addTokenCandidate(kwUNION)
		p.addTokenCandidate(kwFOR)
		return nil, &ParseError{Message: "collecting"}
	}

	// Parse joins
	for {
		jt, ok, err := p.matchJoinType()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		// Completion: after JOIN keyword, offer table_ref candidates for the right-side table.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
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
		condStart := p.pos()
		if _, ok := p.match(kwON); ok {
			// Completion: after ON keyword, offer columnref candidates for join condition.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, &ParseError{Message: "collecting"}
			}
			condStart = p.pos()
			cond, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			join.Condition = &nodes.OnCondition{
				Loc:  nodes.Loc{Start: condStart, End: p.pos()},
				Expr: cond,
			}
		} else if _, ok := p.match(kwUSING); ok {
			// For USING, we manually consume '(' so we can add a completion checkpoint inside.
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			// Completion: after USING (, offer columnref candidates for shared columns.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, &ParseError{Message: "collecting"}
			}
			var cols []string
			for {
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				cols = append(cols, name)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			join.Condition = &nodes.UsingCondition{Loc: nodes.Loc{Start: condStart, End: p.pos()}, Columns: cols}
		}

		join.Loc.End = p.pos()
		left = join
	}

	return left, nil
}

// matchJoinType checks for a join keyword combination and returns the join type.
// Returns an error if a join modifier keyword (INNER, LEFT, etc.) is found but
// the required JOIN keyword is missing.
func (p *Parser) matchJoinType() (nodes.JoinType, bool, error) {
	switch p.cur.Type {
	case kwJOIN:
		p.advance()
		return nodes.JoinInner, true, nil
	case kwINNER:
		p.advance()
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinInner, true, nil
	case kwLEFT:
		p.advance()
		p.match(kwOUTER)
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinLeft, true, nil
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinRight, true, nil
	case kwCROSS:
		p.advance()
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinCross, true, nil
	case kwNATURAL:
		p.advance()
		// NATURAL [LEFT|RIGHT] [OUTER] JOIN
		if _, ok := p.match(kwLEFT); ok {
			p.match(kwOUTER)
			if _, err := p.expect(kwJOIN); err != nil {
				return 0, false, err
			}
			return nodes.JoinNaturalLeft, true, nil
		}
		if _, ok := p.match(kwRIGHT); ok {
			p.match(kwOUTER)
			if _, err := p.expect(kwJOIN); err != nil {
				return 0, false, err
			}
			return nodes.JoinNaturalRight, true, nil
		}
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinNatural, true, nil
	case kwSTRAIGHT_JOIN:
		p.advance()
		return nodes.JoinStraight, true, nil
	}
	return 0, false, nil
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
			alias, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		} else if p.isIdentToken() && !p.isSelectTerminator() {
			alias, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		}

		sub.Loc.End = p.pos()
		return sub, nil
	}

	if p.cur.Type == '(' {
		startPos := p.pos()
		p.advance()

		// Subquery (derived table): (SELECT ...) or (WITH ... SELECT ...)
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			sel, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}

			sub := &nodes.SubqueryExpr{
				Loc:    nodes.Loc{Start: startPos},
				Select: sel,
			}

			// Optional alias: [AS] alias
			if _, ok := p.match(kwAS); ok {
				alias, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				sub.Alias = alias
			} else if p.isIdentToken() && !p.isSelectTerminator() {
				alias, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				sub.Alias = alias
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

	// DUAL — MySQL reserved keyword used as a dummy table in FROM clause.
	// After registration as a reserved keyword, parseIdent rejects it,
	// so we handle it explicitly here.
	if p.cur.Type == kwDUAL {
		tok := p.advance()
		return &nodes.TableRef{
			Loc:  nodes.Loc{Start: tok.Loc, End: p.pos()},
			Name: "dual",
		}, nil
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

	// [AS] alias (required for JSON_TABLE)
	p.match(kwAS)
	alias, _, aErr := p.parseIdent()
	if aErr != nil {
		return nil, aErr
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
	name, _, err := p.parseIdent()
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

		// Completion: after ORDER BY expression, offer ASC/DESC/LIMIT and comma.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwASC)
			p.addTokenCandidate(kwDESC)
			p.addTokenCandidate(kwLIMIT)
			p.addTokenCandidate(',')
			return nil, &ParseError{Message: "collecting"}
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

	// Completion: after FOR keyword, offer UPDATE and SHARE.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwUPDATE)
		p.addTokenCandidate(kwSHARE)
		return nil, &ParseError{Message: "collecting"}
	}

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

	// Completion: after FOR UPDATE/SHARE, offer OF, NOWAIT, SKIP.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwNOWAIT)
		p.addTokenCandidate(kwSKIP)
		return nil, &ParseError{Message: "collecting"}
	}

	// OF table_list
	if p.cur.Type == kwOF {
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

// parseLockInShareMode parses LOCK IN SHARE MODE (legacy syntax).
//
//	LOCK IN SHARE MODE
func (p *Parser) parseLockInShareMode() (*nodes.ForUpdate, error) {
	start := p.pos()
	p.advance() // consume LOCK

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSHARE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwMODE); err != nil {
		return nil, err
	}

	return &nodes.ForUpdate{
		Loc:             nodes.Loc{Start: start, End: p.pos()},
		Share:           true,
		LockInShareMode: true,
	}, nil
}

// parseIntoClause parses INTO OUTFILE / DUMPFILE / var_list.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
//
//	into_option: {
//	    INTO OUTFILE 'file_name'
//	        [CHARACTER SET charset_name]
//	        [{FIELDS | COLUMNS}
//	            [TERMINATED BY 'string']
//	            [[OPTIONALLY] ENCLOSED BY 'char']
//	            [ESCAPED BY 'char']
//	        ]
//	        [LINES
//	            [STARTING BY 'string']
//	            [TERMINATED BY 'string']
//	        ]
//	  | INTO DUMPFILE 'file_name'
//	  | INTO var_name [, var_name] ...
//	}
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

		// [CHARACTER SET charset_name]
		if p.cur.Type == kwCHARACTER {
			p.advance()
			p.match(kwSET)
			charset, _, csErr := p.parseIdent()
			if csErr != nil {
				return nil, csErr
			}
			into.Charset = charset
		} else if p.cur.Type == kwCHARSET {
			p.advance()
			charset, _, csErr := p.parseIdent()
			if csErr != nil {
				return nil, csErr
			}
			into.Charset = charset
		}

		// [{FIELDS | COLUMNS} ...]
		if p.cur.Type == kwFIELDS || p.cur.Type == kwCOLUMNS {
			p.advance()
			into.HasFieldsClause = true

			// [TERMINATED BY 'string']
			if p.cur.Type == kwTERMINATED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					into.FieldsTerminatedBy = p.cur.Str
					p.advance()
				}
			}

			// [[OPTIONALLY] ENCLOSED BY 'char']
			if p.cur.Type == kwOPTIONALLY {
				p.advance()
				into.FieldsOptionalEncl = true
			}
			if p.cur.Type == kwENCLOSED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					into.FieldsEnclosedBy = p.cur.Str
					p.advance()
				}
			}

			// [ESCAPED BY 'char']
			if p.cur.Type == kwESCAPED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					into.FieldsEscapedBy = p.cur.Str
					p.advance()
				}
			}
		}

		// [LINES ...]
		if p.cur.Type == kwLINES {
			p.advance()
			into.HasLinesClause = true

			// [STARTING BY 'string']
			if p.cur.Type == kwSTARTING {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					into.LinesStartingBy = p.cur.Str
					p.advance()
				}
			}

			// [TERMINATED BY 'string']
			if p.cur.Type == kwTERMINATED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					into.LinesTerminatedBy = p.cur.Str
					p.advance()
				}
			}
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

// setOpPrecedence returns the precedence for a set operation.
// INTERSECT has higher precedence than UNION/EXCEPT in MySQL 8.0.31+.
func setOpPrecedence(tokType int) int {
	switch tokType {
	case kwINTERSECT:
		return 2
	case kwUNION, kwEXCEPT:
		return 1
	}
	return 0
}

// parseSetOperation parses UNION/INTERSECT/EXCEPT [ALL|DISTINCT] SELECT chains
// with INTERSECT binding tighter than UNION/EXCEPT.
func (p *Parser) parseSetOperation(left *nodes.SelectStmt) (*nodes.SelectStmt, error) {
	return p.parseSetOpWithPrecedence(left, 1)
}

// parseSetOpWithPrecedence implements precedence climbing for set operations.
func (p *Parser) parseSetOpWithPrecedence(left *nodes.SelectStmt, minPrec int) (*nodes.SelectStmt, error) {
	for {
		prec := setOpPrecedence(p.cur.Type)
		if prec < minPrec {
			break
		}

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

		// Completion: after UNION/INTERSECT/EXCEPT keyword, offer ALL, DISTINCT, SELECT.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwALL)
			p.addTokenCandidate(kwDISTINCT)
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
		}

		all := false
		if _, ok := p.match(kwALL); ok {
			all = true
		} else {
			p.match(kwDISTINCT)
		}

		// Completion: after ALL/DISTINCT, offer SELECT.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
		}

		if p.cur.Type != kwSELECT && p.cur.Type != '(' {
			if p.cur.Type == tokEOF {
				return nil, p.syntaxErrorAtCur()
			}
			return nil, &ParseError{
				Message:  "expected SELECT after set operation",
				Position: p.cur.Loc,
			}
		}

		right, err := p.parseSelectStmtBase()
		if err != nil {
			return nil, err
		}

		// If next op has higher precedence, recurse to consume it first
		for {
			nextPrec := setOpPrecedence(p.cur.Type)
			if nextPrec <= prec {
				break
			}
			right, err = p.parseSetOpWithPrecedence(right, nextPrec)
			if err != nil {
				return nil, err
			}
		}

		// Hoist trailing ORDER BY / LIMIT from the right-hand SELECT to the
		// set-operation node. In MySQL, ORDER BY and LIMIT after a set
		// operation (without parentheses) apply to the entire combined result,
		// not just the last SELECT.
		var orderBy []*nodes.OrderByItem
		var limit *nodes.Limit
		rightLeaf := right
		for rightLeaf.SetOp != nodes.SetOpNone && rightLeaf.Right != nil {
			rightLeaf = rightLeaf.Right
		}
		if len(rightLeaf.OrderBy) > 0 {
			orderBy = rightLeaf.OrderBy
			rightLeaf.OrderBy = nil
		}
		if rightLeaf.Limit != nil {
			limit = rightLeaf.Limit
			rightLeaf.Limit = nil
		}

		left = &nodes.SelectStmt{
			Loc:     nodes.Loc{Start: left.Loc.Start, End: right.Loc.End},
			SetOp:   op,
			SetAll:  all,
			Left:    left,
			Right:   right,
			OrderBy: orderBy,
			Limit:   limit,
		}
	}

	return left, nil
}

// parseParenIdentList parses a parenthesized comma-separated list of identifiers.
func (p *Parser) parseParenIdentList() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var names []string
	for {
		name, _, err := p.parseIdent()
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
		name, _, err := p.parseIdent()
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
			var refErr error
			wd.RefName, _, refErr = p.parseIdent()
			if refErr != nil {
				return nil, refErr
			}
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
// parseTableStmt parses a TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/table.html
//
//	TABLE table_name
//	    [ORDER BY column_name]
//	    [LIMIT number [OFFSET number]]
//	    [INTO OUTFILE 'file_name'
//	        [{FIELDS | COLUMNS} ...]
//	        [LINES ...]
//	    | INTO DUMPFILE 'file_name'
//	    | INTO var_name [, var_name] ...]
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

	// INTO clause
	if p.cur.Type == kwINTO {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
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

// parseIndexHints parses an optional list of index hints following a table reference.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/index-hints.html
//
//	index_hint_list:
//	    index_hint [index_hint] ...
//
//	index_hint:
//	    USE {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] ([index_list])
//	  | {IGNORE|FORCE} {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] (index_list)
//
//	index_list:
//	    index_name [, index_name] ...
func (p *Parser) parseIndexHints() ([]*nodes.IndexHint, error) {
	var hints []*nodes.IndexHint
	for p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		hint, err := p.parseIndexHint()
		if err != nil {
			return hints, err
		}
		hints = append(hints, hint)
	}
	return hints, nil
}

// parseIndexHint parses a single index hint.
//
//	index_hint:
//	    USE {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] ([index_list])
//	  | {IGNORE|FORCE} {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] (index_list)
func (p *Parser) parseIndexHint() (*nodes.IndexHint, error) {
	start := p.pos()
	hint := &nodes.IndexHint{
		Loc: nodes.Loc{Start: start},
	}

	// Parse hint type: USE | FORCE | IGNORE
	switch p.cur.Type {
	case kwUSE:
		hint.Type = nodes.HintUse
	case kwFORCE:
		hint.Type = nodes.HintForce
	case kwIGNORE:
		hint.Type = nodes.HintIgnore
	default:
		return nil, &ParseError{Message: "expected USE, FORCE, or IGNORE", Position: p.cur.Loc}
	}
	p.advance()

	// Parse INDEX | KEY
	if p.cur.Type != kwINDEX && p.cur.Type != kwKEY {
		return nil, &ParseError{Message: "expected INDEX or KEY", Position: p.cur.Loc}
	}
	p.advance()

	// Optional FOR {JOIN | ORDER BY | GROUP BY}
	hint.Scope = nodes.HintScopeAll
	if p.cur.Type == kwFOR {
		p.advance()
		switch p.cur.Type {
		case kwJOIN:
			hint.Scope = nodes.HintScopeJoin
			p.advance()
		case kwORDER:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			hint.Scope = nodes.HintScopeOrderBy
		case kwGROUP:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			hint.Scope = nodes.HintScopeGroupBy
		default:
			return nil, &ParseError{Message: "expected JOIN, ORDER BY, or GROUP BY", Position: p.cur.Loc}
		}
	}

	// Parse (index_list) — required parens, list may be empty for USE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside index hint parentheses, offer index_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("index_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type != ')' {
		// Parse comma-separated index names
		for {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			hint.Indexes = append(hint.Indexes, name)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	hint.Loc.End = p.pos()
	return hint, nil
}
