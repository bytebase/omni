// Package parser - select.go implements T-SQL SELECT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSelectStmt parses a full SELECT statement.
//
// BNF: mssql/parser/bnf/select-transact-sql.bnf
//
//	<SELECT statement> ::=
//	    [ WITH { [ XMLNAMESPACES , ] [ <common_table_expression> [ , ...n ] ] } ]
//	    <query_expression>
//	    [ ORDER BY <order_by_expression> ]
//	    [ <FOR Clause> ]
//	    [ OPTION ( <query_hint> [ , ...n ] ) ]
//
//	<query_expression> ::=
//	    { <query_specification> | ( <query_expression> ) }
//	    [  { UNION [ ALL ] | EXCEPT | INTERSECT }
//	        <query_specification> | ( <query_expression> ) [ ...n ] ]
//
//	<query_specification> ::=
//	SELECT [ ALL | DISTINCT ]
//	    [ TOP ( expression ) [ PERCENT ] [ WITH TIES ] ]
//	    <select_list>
//	    [ INTO new_table ]
//	    [ FROM { <table_source> } [ , ...n ] ]
//	    [ WHERE <search_condition> ]
//	    [ <GROUP BY> ]
//	    [ HAVING <search_condition> ]
//	    [ WINDOW windowDefinition [ , windowDefinition ]* ]
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error) {
	// Entering a SELECT (top-level or subquery) resets the search-condition
	// depth so that the select list / scalar subexpressions are parsed as
	// scalar expressions even if we entered from an outer WHERE or HAVING.
	// The SELECT's own WHERE / HAVING / JOIN ON clauses will increment the
	// depth back for their duration.
	savedDepth := p.searchCondDepth
	p.searchCondDepth = 0
	defer func() { p.searchCondDepth = savedDepth }()

	loc := p.pos()

	// WITH clause (CTE)
	var withClause *nodes.WithClause
	if p.cur.Type == kwWITH {
		var err error
		withClause, err = p.parseWithClause()
		if err != nil {
			return nil, err
		}
	}

	if p.cur.Type != kwSELECT {
		return nil, nil
	}
	p.advance() // consume SELECT

	// Completion: after SELECT keyword → target list candidates
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		p.addTokenCandidate(kwDISTINCT)
		p.addTokenCandidate(kwTOP)
		p.addTokenCandidate(kwALL)
		p.addTokenCandidate('*')
		return nil, errCollecting
	}

	stmt := &nodes.SelectStmt{
		WithClause: withClause,
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// ALL | DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		stmt.Distinct = true
		// Completion: after DISTINCT → columnref context
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate('*')
			return nil, errCollecting
		}
	} else if _, ok := p.match(kwALL); ok {
		stmt.All = true
	}

	// TOP clause
	if p.cur.Type == kwTOP {
		var err error
		stmt.Top, err = p.parseTopClause()
		if err != nil {
			return nil, err
		}
		// Completion: after TOP N → columnref context
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate('*')
			return nil, errCollecting
		}
	}

	// Target list — at least one result column is required per T-SQL.
	var err error
	stmt.TargetList, err = p.parseTargetList()
	if err != nil {
		return nil, err
	}
	if stmt.TargetList == nil || len(stmt.TargetList.Items) == 0 {
		return nil, p.unexpectedToken()
	}

	// INTO
	if _, ok := p.match(kwINTO); ok {
		stmt.IntoTable, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if stmt.IntoTable == nil {
			return nil, p.newParseError(p.cur.Loc, "expected table name after INTO")
		}
	}

	// FROM
	if _, ok := p.match(kwFROM); ok {
		// Completion: after FROM → table_ref context
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		stmt.FromClause, err = p.parseFromClause()
		if err != nil {
			return nil, err
		}
		if stmt.FromClause == nil || len(stmt.FromClause.Items) == 0 {
			return nil, p.unexpectedToken()
		}
		// Completion: after FROM clause (before WHERE/JOIN/ORDER etc.)
		if p.collectMode() {
			p.addTokenCandidate(kwWHERE)
			p.addTokenCandidate(kwJOIN)
			p.addTokenCandidate(kwINNER)
			p.addTokenCandidate(kwLEFT)
			p.addTokenCandidate(kwRIGHT)
			p.addTokenCandidate(kwCROSS)
			p.addTokenCandidate(kwFULL)
			p.addTokenCandidate(kwOUTER)
			p.addTokenCandidate(kwAPPLY)
			p.addTokenCandidate(kwORDER)
			p.addTokenCandidate(kwGROUP)
			p.addTokenCandidate(kwHAVING)
			p.addTokenCandidate(kwUNION)
			p.addTokenCandidate(kwFOR)
			p.addTokenCandidate(kwOPTION)
			p.addTokenCandidate(kwPIVOT)
			p.addTokenCandidate(kwUNPIVOT)
			return nil, errCollecting
		}
	}

	// WHERE
	if _, ok := p.match(kwWHERE); ok {
		// Completion: after WHERE → columnref context
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		p.enterSearchCondition()
		stmt.WhereClause, err = p.parseExpr()
		p.leaveSearchCondition()
		if err != nil {
			return nil, err
		}
		if stmt.WhereClause == nil {
			return nil, p.unexpectedToken()
		}
	}

	// GROUP BY [ALL]
	if p.cur.Type == kwGROUP {
		p.advance()
		if _, err := p.expect(kwBY); err == nil {
			// Completion: after GROUP BY → columnref context
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				p.addRuleCandidate("func_name")
				return nil, errCollecting
			}
			if _, ok := p.match(kwALL); ok {
				stmt.GroupByAll = true
			}
			stmt.GroupByClause, err = p.parseGroupByList()
			if err != nil {
				return nil, err
			}
			if stmt.GroupByClause == nil || len(stmt.GroupByClause.Items) == 0 {
				return nil, p.unexpectedToken()
			}
			// Completion: after GROUP BY list → clause keywords
			if p.collectMode() {
				p.addTokenCandidate(kwHAVING)
				p.addTokenCandidate(kwORDER)
				p.addTokenCandidate(kwFOR)
				p.addTokenCandidate(kwOPTION)
				p.addTokenCandidate(kwUNION)
				p.addTokenCandidate(kwINTERSECT)
				p.addTokenCandidate(kwEXCEPT)
				return nil, errCollecting
			}
		}
	}

	// HAVING
	if _, ok := p.match(kwHAVING); ok {
		// Completion: after HAVING → columnref context
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		p.enterSearchCondition()
		stmt.HavingClause, err = p.parseExpr()
		p.leaveSearchCondition()
		if err != nil {
			return nil, err
		}
		if stmt.HavingClause == nil {
			return nil, p.unexpectedToken()
		}
	}

	// WINDOW clause (named window definitions)
	if p.cur.Type == kwWINDOW {
		stmt.WindowClause, err = p.parseWindowClause()
		if err != nil {
			return nil, err
		}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err == nil {
			// Completion: after ORDER BY → columnref context
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				p.addRuleCandidate("func_name")
				return nil, errCollecting
			}
			stmt.OrderByClause, err = p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			if stmt.OrderByClause == nil || len(stmt.OrderByClause.Items) == 0 {
				return nil, p.unexpectedToken()
			}
			// Completion: after ORDER BY list → sort direction, OFFSET, etc.
			if p.collectMode() {
				p.addTokenCandidate(kwASC)
				p.addTokenCandidate(kwDESC)
				p.addTokenCandidate(kwOFFSET)
				p.addTokenCandidate(kwFOR)
				p.addTokenCandidate(kwOPTION)
				return nil, errCollecting
			}
		}
	}

	// OFFSET ... FETCH
	if p.cur.Type == kwOFFSET {
		p.advance()
		// Completion: after OFFSET → numeric context (no specific candidates)
		stmt.OffsetClause, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		// Consume optional ROWS/ROW
		if p.cur.Type == kwROWS || p.cur.Type == kwROW {
			p.advance()
		}
		// FETCH NEXT n ROWS ONLY
		if p.cur.Type == kwFETCH {
			fetchLoc := p.pos()
			p.advance()
			// NEXT or FIRST
			if p.cur.Type == kwNEXT || p.cur.Type == kwFIRST {
				p.advance()
			}
			count, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			// ROWS/ROW
			if p.cur.Type == kwROWS || p.cur.Type == kwROW {
				p.advance()
			}
			// ONLY
			if p.cur.Type == kwONLY {
				p.advance()
			}
			stmt.FetchClause = &nodes.FetchClause{
				Count: count,
				Loc:   nodes.Loc{Start: fetchLoc, End: p.prevEnd()},
			}
		}
	}

	// FOR XML / FOR JSON / FOR BROWSE
	if p.cur.Type == kwFOR {
		next := p.peekNext()
		if next.Type == kwXML || next.Type == kwJSON || next.Type == kwBROWSE {
			stmt.ForClause, err = p.parseForClause()
			if err != nil {
				return nil, err
			}
		}
	}

	// Completion: after OPTION clause or FOR clause position
	// This is checked here in case we fall through without FOR/OPTION.

	// OPTION clause
	if p.cur.Type == kwOPTION {
		stmt.OptionClause, err = p.parseOptionClause()
		if err != nil {
			return nil, err
		}
	}

	// UNION / INTERSECT / EXCEPT
	if p.cur.Type == kwUNION || p.cur.Type == kwINTERSECT || p.cur.Type == kwEXCEPT {
		return p.parseSetOperation(stmt)
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSetOperation parses UNION/INTERSECT/EXCEPT.
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

	// Completion: after UNION/INTERSECT/EXCEPT → ALL, SELECT
	if p.collectMode() {
		p.addTokenCandidate(kwALL)
		p.addTokenCandidate(kwSELECT)
		return nil, errCollecting
	}

	all := false
	if _, ok := p.match(kwALL); ok {
		all = true
		// Completion: after UNION ALL → SELECT
		if p.collectMode() {
			p.addTokenCandidate(kwSELECT)
			return nil, errCollecting
		}
	}

	right, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.unexpectedToken()
	}
	return &nodes.SelectStmt{
		Op:   op,
		All:  all,
		Larg: left,
		Rarg: right,
		Loc:  left.Loc,
	}, nil
}

// parseWithClause parses WITH [XMLNAMESPACES(...),] cte_name [(col_list)] AS (select), ...
//
// BNF: mssql/parser/bnf/select-transact-sql.bnf
//
//	[ WITH { [ XMLNAMESPACES , ] [ <common_table_expression> [ , ...n ] ] } ]
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/with-common-table-expression-transact-sql
func (p *Parser) parseWithClause() (*nodes.WithClause, error) {
	loc := p.pos()
	// Record CTE position for completion module
	if p.completing {
		p.addCTEPosition(loc)
	}
	p.advance() // consume WITH

	// Completion: after WITH → identifier context for CTE name
	if p.collectMode() {
		p.addRuleCandidate("cte_name")
		return nil, errCollecting
	}

	wc := &nodes.WithClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional XMLNAMESPACES (...)
	if p.cur.Type == kwXMLNAMESPACES {
		var err error
		wc.XmlNamespaces, err = p.parseXmlNamespaces()
		if err != nil {
			return nil, err
		}
		p.match(',') // consume comma between XMLNAMESPACES and CTEs
	}

	// CTE list terminator is multi-token (SELECT / INSERT / UPDATE / DELETE /
	// MERGE); parseCTE returns nil when the current token does not start a CTE
	// name, which lets us break the loop cleanly at non-SELECT DML statements.
	// parseCommaList's single-terminator contract does not fit here.
	var ctes []nodes.Node
	for p.cur.Type != kwSELECT && p.cur.Type != tokEOF && p.cur.Type != ';' {
		cte, err := p.parseCTE()
		if err != nil {
			return nil, err
		}
		if cte == nil {
			break
		}
		ctes = append(ctes, cte)
		if _, ok := p.match(','); !ok {
			break
		}
		// Reject trailing comma: after consuming ',', we must have another CTE
		// before a terminator.
		if p.cur.Type == kwSELECT || p.cur.Type == tokEOF || p.cur.Type == ';' {
			return nil, p.unexpectedToken()
		}
	}
	wc.CTEs = &nodes.List{Items: ctes}
	wc.Loc.End = p.prevEnd()
	return wc, nil
}

// parseXmlNamespaces parses XMLNAMESPACES ( namespace_decl [, ...n] ).
//
//	XMLNAMESPACES ( uri AS prefix [, ...n] | DEFAULT uri [, ...n] )
func (p *Parser) parseXmlNamespaces() (*nodes.List, error) {
	p.advance() // consume XMLNAMESPACES
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	decls, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		loc := p.pos()
		decl := &nodes.XmlNamespaceDecl{Loc: nodes.Loc{Start: loc, End: -1}}
		if _, ok := p.match(kwDEFAULT); ok {
			decl.IsDefault = true
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				decl.URI = p.cur.Str
				p.advance()
			}
		} else {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				decl.URI = p.cur.Str
				p.advance()
			}
			if _, ok := p.match(kwAS); ok {
				if name, ok := p.parseIdentifier(); ok {
					decl.Prefix = name
				}
			}
		}
		decl.Loc.End = p.prevEnd()
		return decl, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if len(decls) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: decls}, nil
}

// parseCTE parses a single CTE: name [(columns)] AS (query).
func (p *Parser) parseCTE() (*nodes.CommonTableExpr, error) {
	loc := p.pos()
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, nil
	}

	cte := &nodes.CommonTableExpr{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		// Completion: WITH cte (|) → identifier context for column names
		if p.collectMode() {
			p.addRuleCandidate("cte_column_name")
			return nil, errCollecting
		}
		cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			colName, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			return &nodes.String{Str: colName}, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		cte.Columns = &nodes.List{Items: cols}
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: WITH cte AS (|) → SELECT keyword
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return nil, errCollecting
	}

	var err error
	cte.Query, err = p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if cte.Query == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	cte.Loc.End = p.prevEnd()
	return cte, nil
}

// parseTopClause parses TOP (expr) [PERCENT] [WITH TIES].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/top-transact-sql
func (p *Parser) parseTopClause() (*nodes.TopClause, error) {
	loc := p.pos()
	p.advance() // consume TOP

	// Completion: after TOP → numeric context (no specific rule candidates)
	// Just return without candidates since this is a numeric context.

	tc := &nodes.TopClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// TOP (expr) or TOP literal
	var err error
	if p.cur.Type == '(' {
		p.advance()
		tc.Count, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if tc.Count == nil {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		tc.Count, err = p.parsePrimary()
		if err != nil {
			return nil, err
		}
		if tc.Count == nil {
			return nil, p.unexpectedToken()
		}
	}

	// PERCENT
	if _, ok := p.match(kwPERCENT); ok {
		tc.Percent = true
	}

	// WITH TIES
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwTIES {
			p.advance() // consume WITH
			p.advance() // consume TIES
			tc.WithTies = true
		}
	}

	tc.Loc.End = p.prevEnd()
	return tc, nil
}

// parseTargetList parses a comma-separated list of result columns.
func (p *Parser) parseTargetList() (*nodes.List, error) {
	var targets []nodes.Node
	for {
		// Completion: at start of each target list item → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate('*')
			return nil, errCollecting
		}
		targetLoc := p.pos()

		// T-SQL `column_alias = expression` and `@var = expression` forms.
		// Disambiguated by 2-token lookahead: { ident | @var } followed by '='
		// at the head of a select list item can only be the assignment form
		// (a bare expression that starts with `ident =` would be a top-level
		// comparison, which is not a valid select-list element shape).
		if p.peekNext().Type == '=' {
			if p.cur.Type == tokVARIABLE {
				name := p.cur.Str
				p.advance() // consume @var
				p.advance() // consume =
				rhs, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				assign := &nodes.SelectAssign{
					Variable: name,
					Value:    rhs,
					Loc:      nodes.Loc{Start: targetLoc, End: p.prevEnd()},
				}
				targets = append(targets, &nodes.ResTarget{
					Val: assign,
					Loc: nodes.Loc{Start: targetLoc, End: p.prevEnd()},
				})
				if _, ok := p.match(','); !ok {
					break
				}
				continue
			}
			if p.isIdentLike() && !p.isBareAliasExcluded() {
				aliasLoc := p.pos()
				alias := p.cur.Str
				p.advance() // consume ident
				p.advance() // consume =
				rhs, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				targets = append(targets, &nodes.ResTarget{
					Name: alias,
					Val:  rhs,
					Loc:  nodes.Loc{Start: targetLoc, End: p.prevEnd()},
				})
				if p.completing {
					p.addSelectAliasPosition(aliasLoc)
				}
				if _, ok := p.match(','); !ok {
					break
				}
				continue
			}
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			if len(targets) > 0 {
				// Just consumed ',' with no expr to follow — trailing comma.
				return nil, p.unexpectedToken()
			}
			break
		}
		target := &nodes.ResTarget{
			Val: expr,
			Loc: nodes.Loc{Start: targetLoc, End: -1},
		}
		// Check for alias: AS name or just name (but not keywords that start clauses).
		// Reserved keywords (CoreKeyword) cannot be unquoted aliases per T-SQL —
		// writing `SELECT x AS FROM` requires `AS [FROM]` which the lexer emits
		// as tokIDENT and thus already passes isIdentLike.
		if _, ok := p.match(kwAS); ok {
			if p.isIdentLike() {
				aliasLoc := p.pos()
				target.Name = p.cur.Str
				p.advance()
				if p.completing {
					p.addSelectAliasPosition(aliasLoc)
				}
			}
		} else if p.isIdentLike() && !p.isBareAliasExcluded() {
			aliasLoc := p.pos()
			target.Name = p.cur.Str
			p.advance()
			if p.completing {
				p.addSelectAliasPosition(aliasLoc)
			}
		}
		target.Loc.End = p.prevEnd()
		targets = append(targets, target)
		if _, ok := p.match(','); !ok {
			break
		}
		// Completion: after comma in target list → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate('*')
			return nil, errCollecting
		}
	}
	return &nodes.List{Items: targets}, nil
}

// parseFromClause parses a FROM clause table source list.
func (p *Parser) parseFromClause() (*nodes.List, error) {
	var sources []nodes.Node
	for {
		source, err := p.parseTableSource()
		if err != nil {
			return nil, err
		}
		if source == nil {
			break
		}
		sources = append(sources, source)
		if _, ok := p.match(','); !ok {
			break
		}
		// Completion: after comma in FROM list → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
	}
	return &nodes.List{Items: sources}, nil
}

// parseTableSource parses a single table source (table, subquery, join).
func (p *Parser) parseTableSource() (nodes.TableExpr, error) {
	left, err := p.parsePrimaryTableSource()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}

	// Completion: after primary table source → JOIN keyword candidates
	if p.collectMode() {
		p.addTokenCandidate(kwJOIN)
		p.addTokenCandidate(kwINNER)
		p.addTokenCandidate(kwLEFT)
		p.addTokenCandidate(kwRIGHT)
		p.addTokenCandidate(kwCROSS)
		p.addTokenCandidate(kwFULL)
		p.addTokenCandidate(kwOUTER)
		p.addTokenCandidate(kwWHERE)
		p.addTokenCandidate(kwORDER)
		p.addTokenCandidate(kwGROUP)
		p.addTokenCandidate(kwHAVING)
		p.addTokenCandidate(kwUNION)
		p.addTokenCandidate(kwFOR)
		p.addTokenCandidate(kwOPTION)
		p.addTokenCandidate(kwAPPLY)
		return nil, errCollecting
	}

	// Parse joins
	for {
		joinLoc := p.pos()
		jt, ok := p.matchJoinType()
		if !ok {
			break
		}
		// Completion: after JOIN keyword → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		right, err := p.parsePrimaryTableSource()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		start := nodes.NodeLoc(left).Start
		if start < 0 {
			start = joinLoc
		}
		join := &nodes.JoinClause{
			Type:  jt,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: start, End: -1},
		}
		// ON condition (not for CROSS JOIN / CROSS APPLY / OUTER APPLY)
		if jt != nodes.JoinCross && jt != nodes.JoinCrossApply && jt != nodes.JoinOuterApply {
			if _, ok := p.match(kwON); ok {
				// Completion: after ON → columnref
				if p.collectMode() {
					p.addRuleCandidate("columnref")
					p.addRuleCandidate("func_name")
					return nil, errCollecting
				}
				p.enterSearchCondition()
				join.Condition, err = p.parseExpr()
				p.leaveSearchCondition()
				if err != nil {
					return nil, err
				}
				if join.Condition == nil {
					return nil, p.unexpectedToken()
				}
			}
		}
		join.Loc.End = p.prevEnd()
		left = join
	}

	return left, nil
}

// parsePrimaryTableSource parses a base table, subquery, or function call as table source.
func (p *Parser) parsePrimaryTableSource() (nodes.TableExpr, error) {
	// Completion: at start of primary table source → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Parenthesised derived table: (SELECT ...) or (VALUES ...)
	if p.cur.Type == '(' {
		loc := p.pos()
		p.advance()

		var inner nodes.TableExpr
		if p.cur.Type == kwVALUES {
			vc, err := p.parseValuesClause()
			if err != nil {
				return nil, err
			}
			inner = vc
		} else {
			sub, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if sub == nil {
				return nil, p.unexpectedToken()
			}
			inner = &nodes.SubqueryExpr{
				Query: sub,
				Loc:   nodes.Loc{Start: loc, End: -1},
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		// Alias [ ( col1, col2, ... ) ]
		alias, cols, err := p.parseAliasAndOptionalColumnList()
		if err != nil {
			return nil, err
		}
		if alias != "" {
			result := &nodes.AliasedTableRef{
				Table:   inner,
				Alias:   alias,
				Columns: cols,
				Loc:     nodes.Loc{Start: loc, End: -1},
			}
			return p.parsePivotUnpivot(result)
		}
		return p.parsePivotUnpivot(inner)
	}

	// Rowset functions: OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
	if p.cur.Type == kwOPENROWSET || p.cur.Type == kwOPENQUERY || p.cur.Type == kwOPENJSON ||
		p.cur.Type == kwOPENDATASOURCE || p.cur.Type == kwOPENXML {
		return p.parseRowsetFunction()
	}

	// Full-text rowset functions: CONTAINSTABLE, FREETEXTTABLE
	if p.cur.Type == kwCONTAINSTABLE || p.cur.Type == kwFREETEXTTABLE {
		return p.parseFullTextTableRef()
	}

	// Semantic table functions: SEMANTICKEYPHRASETABLE,
	// SEMANTICSIMILARITYTABLE, SEMANTICSIMILARITYDETAILSTABLE
	if p.cur.Type == kwSEMANTICKEYPHRASETABLE ||
		p.cur.Type == kwSEMANTICSIMILARITYTABLE ||
		p.cur.Type == kwSEMANTICSIMILARITYDETAILSTABLE {
		return p.parseSemanticTableRef()
	}

	// T-SQL table variable: @t [alias] or @t.Method(args) [alias (cols)]
	if p.cur.Type == tokVARIABLE {
		return p.parseVariableTableSource()
	}

	// Table reference
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, nil
	}

	// Check if this is a function call (table-valued function)
	if p.cur.Type == '(' {
		tvf, err := p.parseTableValuedFunction(ref)
		if err != nil {
			return nil, err
		}
		return p.parsePivotUnpivot(tvf)
	}

	// TABLESAMPLE
	if p.cur.Type == kwTABLESAMPLE {
		ts, err := p.parseTableSampleClause()
		if err != nil {
			return nil, err
		}
		alias := p.parseOptionalAlias()
		if alias != "" {
			ref.Alias = alias
		}
		result := &nodes.AliasedTableRef{
			Table:       ref,
			Alias:       ref.Alias,
			TableSample: ts,
			Loc:         ref.Loc,
		}
		// Table hints after TABLESAMPLE
		if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
			var err error
			result.Hints, err = p.parseTableHints()
			if err != nil {
				return nil, err
			}
		}
		return p.parsePivotUnpivot(result)
	}

	// Alias
	alias := p.parseOptionalAlias()
	if alias != "" {
		ref.Alias = alias
		ref.Loc.End = p.prevEnd()
	}

	// Table hints: WITH (NOLOCK), WITH (INDEX(idx1)), etc.
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		ref.Hints, err = p.parseTableHints()
		if err != nil {
			return nil, err
		}
		ref.Loc.End = p.prevEnd()
	}

	return p.parsePivotUnpivot(ref)
}

// parsePivotUnpivot checks for and parses PIVOT or UNPIVOT after a table source.
func (p *Parser) parsePivotUnpivot(source nodes.TableExpr) (nodes.TableExpr, error) {
	if p.cur.Type == kwPIVOT {
		return p.parsePivotExpr(source)
	}
	if p.cur.Type == kwUNPIVOT {
		return p.parseUnpivotExpr(source)
	}
	return source, nil
}

// parsePivotExpr parses PIVOT (agg_func(col) FOR pivot_col IN ([v1],[v2],...)) AS alias.
func (p *Parser) parsePivotExpr(source nodes.TableExpr) (*nodes.PivotExpr, error) {
	loc := p.pos()
	p.advance() // consume PIVOT

	pivot := &nodes.PivotExpr{
		Source: source,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside PIVOT (...) → aggregate function context
	if p.collectMode() {
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	// Parse aggregate function call
	var err error
	pivot.AggFunc, err = p.parseExpr()
	if err != nil {
		return nil, err
	}

	// FOR column
	if _, ok := p.match(kwFOR); ok {
		if name, ok := p.parseIdentifier(); ok {
			pivot.ForCol = name
		}
	}

	// IN ([v1], [v2], ...)
	if _, ok := p.match(kwIN); ok {
		if _, err := p.expect('('); err == nil {
			vals, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
				if p.collectMode() {
					p.addRuleCandidate("columnref")
					return nil, errCollecting
				}
				name, ok := p.parseIdentifier()
				if !ok {
					return nil, p.unexpectedToken()
				}
				return &nodes.String{Str: name}, nil
			})
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			pivot.InValues = &nodes.List{Items: vals}
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// AS alias
	alias := p.parseOptionalAlias()
	pivot.Alias = alias

	pivot.Loc.End = p.prevEnd()
	return pivot, nil
}

// parseUnpivotExpr parses UNPIVOT (value_col FOR pivot_col IN ([c1],[c2],...)) AS alias.
func (p *Parser) parseUnpivotExpr(source nodes.TableExpr) (*nodes.UnpivotExpr, error) {
	loc := p.pos()
	p.advance() // consume UNPIVOT

	unpivot := &nodes.UnpivotExpr{
		Source: source,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside UNPIVOT (...) → column context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}

	// value column name
	if name, ok := p.parseIdentifier(); ok {
		unpivot.ValueCol = name
	}

	// FOR column
	if _, ok := p.match(kwFOR); ok {
		if name, ok := p.parseIdentifier(); ok {
			unpivot.ForCol = name
		}
	}

	// IN ([c1], [c2], ...)
	if _, ok := p.match(kwIN); ok {
		if _, err := p.expect('('); err == nil {
			cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
				if p.collectMode() {
					p.addRuleCandidate("columnref")
					return nil, errCollecting
				}
				name, ok := p.parseIdentifier()
				if !ok {
					return nil, p.unexpectedToken()
				}
				return &nodes.String{Str: name}, nil
			})
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			unpivot.InCols = &nodes.List{Items: cols}
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// AS alias
	alias := p.parseOptionalAlias()
	unpivot.Alias = alias

	unpivot.Loc.End = p.prevEnd()
	return unpivot, nil
}

// parseTableSampleClause parses TABLESAMPLE (size PERCENT|ROWS) [REPEATABLE (seed)].
func (p *Parser) parseTableSampleClause() (*nodes.TableSampleClause, error) {
	loc := p.pos()
	p.advance() // consume TABLESAMPLE

	ts := &nodes.TableSampleClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var err error
	ts.Size, err = p.parseExpr()
	if err != nil {
		return nil, err
	}

	// PERCENT or ROWS
	if _, ok := p.match(kwPERCENT); ok {
		ts.Unit = "PERCENT"
	} else if _, ok := p.match(kwROWS); ok {
		ts.Unit = "ROWS"
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// REPEATABLE (seed)
	if p.matchIdentCI("REPEATABLE") {
		if _, err := p.expect('('); err == nil {
			ts.Repeatable, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}

	ts.Loc.End = p.prevEnd()
	return ts, nil
}

// parseTableValuedFunction parses a table-valued function call after the name.
func (p *Parser) parseTableValuedFunction(ref *nodes.TableRef) (nodes.TableExpr, error) {
	loc := ref.Loc.Start
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: ref,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	args, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		return arg, nil
	})
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		fc.Args = &nodes.List{Items: args}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()

	alias, cols, err := p.parseAliasAndOptionalColumnList()
	if err != nil {
		return nil, err
	}
	aliasEnd := p.prevEnd()
	return &nodes.AliasedTableRef{
		Table:   fc,
		Alias:   alias,
		Columns: cols,
		Loc:     nodes.Loc{Start: loc, End: aliasEnd},
	}, nil
}

// parseOptionalAlias parses an optional alias (AS name or just name).
// Reserved keywords (CoreKeyword) cannot be unquoted aliases — a bracketed
// form like `AS [FROM]` is already tokIDENT at the lexer level.
func (p *Parser) parseOptionalAlias() string {
	if _, ok := p.match(kwAS); ok {
		if p.isIdentLike() {
			name := p.cur.Str
			p.advance()
			return name
		}
		return ""
	}
	// Bare alias - accept identifiers and context keywords, but NOT clause/statement keywords
	if p.isIdentLike() && !p.isSelectClauseIdent() && !p.isBareAliasExcluded() {
		name := p.cur.Str
		p.advance()
		return name
	}
	return ""
}

// parseAliasAndOptionalColumnList parses the `[ AS alias ] [ ( col1, col2, ... ) ]`
// tail that appears after a derived-table-like source (subquery, TVF, VALUES, rowset).
// The column list is only consumed when a non-empty alias was parsed.
// Returns ("", nil, nil) when no alias is present.
func (p *Parser) parseAliasAndOptionalColumnList() (string, *nodes.List, error) {
	alias := p.parseOptionalAlias()
	if alias == "" {
		return "", nil, nil
	}
	if p.cur.Type != '(' {
		return alias, nil, nil
	}
	p.advance() // consume '('
	cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		return &nodes.String{Str: name}, nil
	})
	if err != nil {
		return alias, nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return alias, nil, err
	}
	return alias, &nodes.List{Items: cols}, nil
}

// isSelectClauseIdent returns true if the current token is a contextual
// keyword that starts a SELECT clause and should not be consumed as a bare alias.
func (p *Parser) isSelectClauseIdent() bool {
	return p.cur.Type == kwWINDOW
}

// isBareAliasExcluded returns true if the current token is a context keyword
// that should NOT be consumed as a bare alias because it starts a new statement
// or has special meaning at statement boundaries.
func (p *Parser) isBareAliasExcluded() bool {
	switch p.cur.Type {
	case kwGO:
		return true
	}
	return false
}

// matchJoinType matches and consumes a join keyword sequence, returning the join type.
func (p *Parser) matchJoinType() (nodes.JoinType, bool) {
	switch p.cur.Type {
	case kwINNER:
		p.advance()
		_, _ = p.expect(kwJOIN)
		return nodes.JoinInner, true
	case kwJOIN:
		p.advance()
		return nodes.JoinInner, true
	case kwLEFT:
		p.advance()
		p.match(kwOUTER)
		_, _ = p.expect(kwJOIN)
		return nodes.JoinLeft, true
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
		_, _ = p.expect(kwJOIN)
		return nodes.JoinRight, true
	case kwFULL:
		p.advance()
		p.match(kwOUTER)
		_, _ = p.expect(kwJOIN)
		return nodes.JoinFull, true
	case kwCROSS:
		p.advance()
		// CROSS JOIN vs CROSS APPLY
		if _, ok := p.match(kwAPPLY); ok {
			return nodes.JoinCrossApply, true
		}
		_, _ = p.expect(kwJOIN)
		return nodes.JoinCross, true
	case kwOUTER:
		// OUTER APPLY
		if p.peekNext().Type == kwAPPLY {
			p.advance() // consume OUTER
			p.advance() // consume APPLY
			return nodes.JoinOuterApply, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// parseForClause parses FOR BROWSE, FOR XML, or FOR JSON.
//
// BNF: mssql/parser/bnf/select-transact-sql.bnf
//
//	[ FOR { BROWSE | <XML> | <JSON> } ]
//
//	<XML> ::=
//	XML
//	{
//	    { RAW [ ( 'ElementName' ) ] | AUTO }
//	    [
//	        <CommonDirectivesForXML>
//	        [ , { XMLDATA | XMLSCHEMA [ ( 'TargetNameSpaceURI' ) ] } ]
//	        [ , ELEMENTS [ XSINIL | ABSENT ] ]
//	    ]
//	  | EXPLICIT
//	    [
//	        <CommonDirectivesForXML>
//	        [ , XMLDATA ]
//	    ]
//	  | PATH [ ( 'ElementName' ) ]
//	    [
//	        <CommonDirectivesForXML>
//	        [ , ELEMENTS [ XSINIL | ABSENT ] ]
//	    ]
//	}
//
//	<CommonDirectivesForXML> ::=
//	[ , BINARY BASE64 ]
//	[ , TYPE ]
//	[ , ROOT [ ( 'RootName' ) ] ]
//
//	<JSON> ::=
//	JSON
//	{
//	    { AUTO | PATH }
//	    [
//	        [ , ROOT [ ( 'RootName' ) ] ]
//	        [ , INCLUDE_NULL_VALUES ]
//	        [ , WITHOUT_ARRAY_WRAPPER ]
//	    ]
//	}
func (p *Parser) parseForClause() (*nodes.ForClause, error) {
	loc := p.pos()
	p.advance() // consume FOR

	// Completion: after FOR → XML, JSON, BROWSE
	if p.collectMode() {
		p.addTokenCandidate(kwXML)
		p.addTokenCandidate(kwJSON)
		p.addTokenCandidate(kwBROWSE)
		return nil, errCollecting
	}

	fc := &nodes.ForClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwBROWSE {
		fc.Mode = nodes.ForBrowse
		p.advance()
		fc.Loc.End = p.prevEnd()
		return fc, nil
	}

	if p.cur.Type == kwXML {
		fc.Mode = nodes.ForXML
		p.advance()
		// Completion: after FOR XML → PATH, RAW, AUTO, EXPLICIT
		if p.collectMode() {
			p.addRuleCandidate("xml_mode") // PATH, RAW, AUTO, EXPLICIT
			return nil, errCollecting
		}
		// RAW, AUTO, EXPLICIT, PATH — SqlScriptDOM ForXmlMode enum.
		if p.isValidOption(forXmlModes) {
			fc.SubMode = strings.ToUpper(p.cur.Str)
			p.advance()
			// RAW('ElementName') or PATH('ElementName')
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokSCONST {
					fc.ElementName = p.cur.Str
					p.advance()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
		}
		// Parse comma-separated XML options
		if err := p.parseForXmlOptions(fc); err != nil {
			return nil, err
		}
	} else if p.cur.Type == kwJSON {
		fc.Mode = nodes.ForJSON
		p.advance()
		// Completion: after FOR JSON → PATH, AUTO
		if p.collectMode() {
			p.addRuleCandidate("json_mode") // PATH, AUTO
			return nil, errCollecting
		}
		// AUTO or PATH — SqlScriptDOM ForJsonMode enum.
		if p.isValidOption(forJsonModes) {
			fc.SubMode = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		// Parse comma-separated JSON options
		if err := p.parseForJsonOptions(fc); err != nil {
			return nil, err
		}
	}

	fc.Loc.End = p.prevEnd()
	return fc, nil
}

// forXmlOptions defines the valid options for FOR XML clauses.
// Matches SqlScriptDOM XmlForClauseOptionsHelper: BINARY, TYPE, ROOT,
// XMLDATA, XMLSCHEMA, ELEMENTS.
var forXmlOptions = newOptionSet(
	kwBINARY, kwTYPE, kwROOT, kwXMLDATA, kwXMLSCHEMA, kwELEMENTS,
)

// parseForXmlOptions parses the comma-separated options after FOR XML {RAW|AUTO|EXPLICIT|PATH}.
//
//	[ , BINARY BASE64 ]
//	[ , TYPE ]
//	[ , ROOT [ ( 'RootName' ) ] ]
//	[ , { XMLDATA | XMLSCHEMA [ ( 'TargetNameSpaceURI' ) ] } ]
//	[ , ELEMENTS [ XSINIL | ABSENT ] ]
func (p *Parser) parseForXmlOptions(fc *nodes.ForClause) error {
	for {
		if p.cur.Type != ',' {
			return nil
		}
		// Peek: after the comma, validate that the next token is a known FOR XML option.
		next := p.peekNext()
		if !p.isValidOptionToken(forXmlOptions, next) {
			return nil
		}
		p.advance() // consume comma
		switch {
		case p.matchIdentCI("BINARY"):
			// BINARY BASE64
			p.matchIdentCI("BASE64")
			fc.BinaryBase64 = true
		case p.cur.Type == kwTYPE:
			p.advance()
			fc.Type = true
		case p.matchIdentCI("ROOT"):
			fc.Root = true
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokSCONST {
					fc.RootName = p.cur.Str
					p.advance()
				}
				_, _ = p.expect(')')
			}
		case p.matchIdentCI("XMLDATA"):
			fc.XmlData = true
		case p.matchIdentCI("XMLSCHEMA"):
			fc.XmlSchema = true
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokSCONST {
					fc.XmlSchemaURI = p.cur.Str
					p.advance()
				}
				_, _ = p.expect(')')
			}
		case p.matchIdentCI("ELEMENTS"):
			fc.Elements = true
			if p.matchIdentCI("XSINIL") {
				fc.ElementsMode = "XSINIL"
			} else if p.matchIdentCI("ABSENT") {
				fc.ElementsMode = "ABSENT"
			}
		default:
			return nil
		}
	}
}

// forJsonOptions defines the valid options for FOR JSON clauses.
// Matches SqlScriptDOM JsonForClauseOptionsHelper: ROOT, INCLUDE_NULL_VALUES,
// WITHOUT_ARRAY_WRAPPER.
var forJsonOptions = newOptionSet(
	kwROOT, kwINCLUDE_NULL_VALUES, kwWITHOUT_ARRAY_WRAPPER,
)

// parseForJsonOptions parses the comma-separated options after FOR JSON {AUTO|PATH}.
//
//	[ , ROOT [ ( 'RootName' ) ] ]
//	[ , INCLUDE_NULL_VALUES ]
//	[ , WITHOUT_ARRAY_WRAPPER ]
func (p *Parser) parseForJsonOptions(fc *nodes.ForClause) error {
	for {
		if p.cur.Type != ',' {
			return nil
		}
		// Peek: after the comma, validate that the next token is a known FOR JSON option.
		next := p.peekNext()
		if !p.isValidOptionToken(forJsonOptions, next) {
			return nil
		}
		p.advance() // consume comma
		switch {
		case p.matchIdentCI("ROOT"):
			fc.Root = true
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokSCONST {
					fc.RootName = p.cur.Str
					p.advance()
				}
				_, _ = p.expect(')')
			}
		case p.matchIdentCI("INCLUDE_NULL_VALUES"):
			fc.IncludeNullValues = true
		case p.matchIdentCI("WITHOUT_ARRAY_WRAPPER"):
			fc.WithoutArrayWrapper = true
		default:
			return nil
		}
	}
}

// parseExprList parses a comma-separated list of expressions.
// parseExprList parses a non-empty comma-separated list of expressions,
// stopping at the first non-expression token (which the caller then matches
// against its own terminator). At least one expression is required and a
// trailing comma is rejected — every current call site (VALUES row,
// ROLLUP/CUBE args, PARTITION BY) needs both invariants per SqlScriptDOM.
func (p *Parser) parseExprList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			// Empty list at head or trailing-comma tail — both invalid.
			return nil, p.unexpectedToken()
		}
		items = append(items, expr)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseGroupByList parses a GROUP BY list which may contain GROUPING SETS, ROLLUP, CUBE.
func (p *Parser) parseGroupByList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		// GROUPING SETS (...)
		if p.cur.Type == kwGROUPING {
			next := p.peekNext()
			if next.Type == kwSETS {
				loc := p.pos()
				p.advance() // consume GROUPING
				p.advance() // consume SETS
				if _, err := p.expect('('); err == nil {
					sets, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
						return p.parseGroupingSet()
					})
					if err != nil {
						return nil, err
					}
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
					items = append(items, &nodes.GroupingSetsExpr{
						Sets: &nodes.List{Items: sets},
						Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
					})
				}
				if _, ok := p.match(','); !ok {
					break
				}
				continue
			}
		}
		// ROLLUP (...)
		if p.cur.Type == kwROLLUP {
			loc := p.pos()
			p.advance() // consume ROLLUP
			if _, err := p.expect('('); err == nil {
				args, err := p.parseExprList()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				items = append(items, &nodes.RollupExpr{
					Args: args,
					Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
				})
			}
			if _, ok := p.match(','); !ok {
				break
			}
			continue
		}
		// CUBE (...)
		if p.cur.Type == kwCUBE {
			loc := p.pos()
			p.advance() // consume CUBE
			if _, err := p.expect('('); err == nil {
				args, err := p.parseExprList()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				items = append(items, &nodes.CubeExpr{
					Args: args,
					Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
				})
			}
			if _, ok := p.match(','); !ok {
				break
			}
			continue
		}
		// Regular expression
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			break
		}
		items = append(items, expr)
		if _, ok := p.match(','); !ok {
			break
		}
		// Completion: after comma in GROUP BY → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseGroupingSet parses a single grouping set: () or (expr, expr, ...) or just expr.
// The empty form `()` is a valid T-SQL grand-total marker and so is the one site
// in the parser that uses commaListAllowEmpty for a paren-delimited list.
func (p *Parser) parseGroupingSet() (*nodes.List, error) {
	if p.cur.Type == '(' {
		p.advance()
		items, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.unexpectedToken()
			}
			return expr, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: items}, nil
	}
	// Single expression as a set
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr != nil {
		return &nodes.List{Items: []nodes.Node{expr}}, nil
	}
	return &nodes.List{Items: nil}, nil
}

// parseWindowClause parses WINDOW window_name AS (window_spec) [, ...].
//
//	WINDOW window_name AS ( [ existing_window_name ]
//	    [ PARTITION BY expr [,...n] ]
//	    [ ORDER BY order_item [,...n] ]
//	    [ <window_frame> ]
//	)
func (p *Parser) parseWindowClause() (*nodes.List, error) {
	p.advance() // consume WINDOW

	var defs []nodes.Node
	for {
		loc := p.pos()
		name, ok := p.parseIdentifier()
		if !ok {
			break
		}
		def := &nodes.WindowDef{
			Name: name,
			Loc:  nodes.Loc{Start: loc, End: -1},
		}

		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		// Optional existing_window_name (must be an ident not followed by keyword like PARTITION, ORDER)
		if p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
			p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS &&
			p.cur.Type != ')' && p.isIdentLike() {
			next := p.peekNext()
			// If next token is a clause keyword or ), this is a refname
			if next.Type == kwPARTITION || next.Type == kwORDER ||
				next.Type == kwROWS || next.Type == kwRANGE || next.Type == kwGROUPS ||
				next.Type == ')' {
				def.RefName = p.cur.Str
				p.advance()
			}
		}

		// PARTITION BY
		if p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect(kwBY); err == nil {
				var err error
				def.PartitionBy, err = p.parseExprList()
				if err != nil {
					return nil, err
				}
			}
		}

		// ORDER BY
		if p.cur.Type == kwORDER {
			p.advance()
			if _, err := p.expect(kwBY); err == nil {
				var err error
				def.OrderBy, err = p.parseOrderByList()
				if err != nil {
					return nil, err
				}
			}
		}

		// Window frame: ROWS | RANGE | GROUPS
		if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
			var err error
			def.Frame, err = p.parseWindowFrame()
			if err != nil {
				return nil, err
			}
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		def.Loc.End = p.prevEnd()
		defs = append(defs, def)

		if _, ok := p.match(','); !ok {
			break
		}
	}

	if len(defs) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: defs}, nil
}

// parseOrderByList parses ORDER BY items.
func (p *Parser) parseOrderByList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		// Completion: at start of each ORDER BY item → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		oloc := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			if len(items) > 0 {
				// Just consumed ',' with no expr to follow — trailing comma.
				return nil, p.unexpectedToken()
			}
			break
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		// Completion: after ORDER BY item (direction) → next valid keywords
		if p.collectMode() {
			p.addTokenCandidate(kwASC)
			p.addTokenCandidate(kwDESC)
			p.addTokenCandidate(kwOFFSET)
			p.addTokenCandidate(kwFOR)
			p.addTokenCandidate(kwOPTION)
			return nil, errCollecting
		}
		items = append(items, &nodes.OrderByItem{
			Expr:    expr,
			SortDir: dir,
			Loc:     nodes.Loc{Start: oloc, End: p.prevEnd()},
		})
		if _, ok := p.match(','); !ok {
			break
		}
		// Completion: after comma in ORDER BY → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseTableHints parses WITH ( <table_hint> [ [ , ] ...n ] ).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-table
//
//	WITH ( <table_hint> [ [ , ] ...n ] )
//
//	<table_hint> ::=
//	{ NOEXPAND
//	  | INDEX ( <index_value> [ , ...n ] ) | INDEX = ( <index_value> )
//	  | FORCESEEK [ ( <index_value> ( <index_column_name> [ , ... ] ) ) ]
//	  | FORCESCAN
//	  | HOLDLOCK
//	  | NOLOCK
//	  | NOWAIT
//	  | PAGLOCK
//	  | READCOMMITTED
//	  | READCOMMITTEDLOCK
//	  | READPAST
//	  | READUNCOMMITTED
//	  | REPEATABLEREAD
//	  | ROWLOCK
//	  | SERIALIZABLE
//	  | SNAPSHOT
//	  | SPATIAL_WINDOW_MAX_CELLS = <integer_value>
//	  | TABLOCK
//	  | TABLOCKX
//	  | UPDLOCK
//	  | XLOCK
//	}
func (p *Parser) parseTableHints() (*nodes.List, error) {
	p.advance() // consume WITH
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside WITH (...) → table hint keywords
	if p.collectMode() {
		p.addRuleCandidate("table_hint")
		return nil, errCollecting
	}

	// Table hints allow an OPTIONAL comma between hints, so this list is not
	// a standard comma-separated list (parseCommaList would mandate commas).
	// We do however reject a trailing comma before ')', matching SqlScriptDOM.
	var hints []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Completion: at start of each hint slot → table hint keywords
		if p.collectMode() {
			p.addRuleCandidate("table_hint")
			return nil, errCollecting
		}
		hint, err := p.parseTableHint()
		if err != nil {
			return nil, err
		}
		if hint == nil {
			break
		}
		hints = append(hints, hint)
		// Optional comma between hints. If a comma is consumed, another hint
		// MUST follow — reject `WITH (NOLOCK,)`.
		if _, ok := p.match(','); ok {
			if p.collectMode() {
				p.addRuleCandidate("table_hint")
				return nil, errCollecting
			}
			if p.cur.Type == ')' || p.cur.Type == tokEOF {
				return nil, p.unexpectedToken()
			}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(hints) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: hints}, nil
}

// parseTableHint parses a single table hint.
func (p *Parser) parseTableHint() (*nodes.TableHint, error) {
	loc := p.pos()

	// INDEX hint: INDEX ( values ) or INDEX = ( value )
	if p.cur.Type == kwINDEX {
		p.advance()
		hint := &nodes.TableHint{
			Name: "INDEX",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		parseIndexVals := func() error {
			vals, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
				return p.parseIndexValue()
			})
			if err != nil {
				return err
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			hint.IndexValues = &nodes.List{Items: vals}
			return nil
		}
		if _, ok := p.match('='); ok {
			// INDEX = ( value )
			if _, err := p.expect('('); err == nil {
				if err := parseIndexVals(); err != nil {
					return nil, err
				}
			}
		} else if p.cur.Type == '(' {
			// INDEX ( values )
			p.advance()
			if err := parseIndexVals(); err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil
	}

	// Check for keyword-based hints that are lexer keywords
	switch p.cur.Type {
	case kwHOLDLOCK:
		p.advance()
		return &nodes.TableHint{Name: "HOLDLOCK", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case kwNOLOCK:
		p.advance()
		return &nodes.TableHint{Name: "NOLOCK", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case kwNOWAIT:
		p.advance()
		return &nodes.TableHint{Name: "NOWAIT", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}

	// All remaining hints are identifiers (not lexer keywords)
	if !p.isIdentLike() {
		return nil, nil
	}

	name := strings.ToUpper(p.cur.Str)
	switch name {
	case "FORCESEEK":
		p.advance()
		hint := &nodes.TableHint{
			Name: "FORCESEEK",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		// Optional: FORCESEEK ( index_value ( col1, col2, ... ) )
		if p.cur.Type == '(' {
			p.advance()
			// index value
			idxVal, err := p.parseIndexValue()
			if err != nil {
				return nil, err
			}
			hint.IndexValues = &nodes.List{Items: []nodes.Node{idxVal}}
			// ( col1, col2, ... )
			if p.cur.Type == '(' {
				p.advance()
				cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
					colName, ok := p.parseIdentifier()
					if !ok {
						return nil, p.unexpectedToken()
					}
					return &nodes.String{Str: colName}, nil
				})
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				hint.ForceSeekColumns = &nodes.List{Items: cols}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case "SPATIAL_WINDOW_MAX_CELLS":
		p.advance()
		hint := &nodes.TableHint{
			Name: "SPATIAL_WINDOW_MAX_CELLS",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		if _, ok := p.match('='); ok {
			var err error
			hint.IntValue, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case "FORCESCAN", "NOEXPAND",
		"PAGLOCK", "READCOMMITTED", "READCOMMITTEDLOCK",
		"READPAST", "READUNCOMMITTED", "REPEATABLEREAD",
		"ROWLOCK", "SERIALIZABLE", "SNAPSHOT",
		"TABLOCK", "TABLOCKX", "UPDLOCK", "XLOCK",
		"KEEPIDENTITY", "KEEPDEFAULTS", "IGNORE_CONSTRAINTS", "IGNORE_TRIGGERS":
		p.advance()
		return &nodes.TableHint{Name: name, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	default:
		return nil, nil
	}
}

// parseIndexValue parses an index value (identifier or integer).
func (p *Parser) parseIndexValue() (nodes.Node, error) {
	if p.cur.Type == tokICONST {
		val := &nodes.Integer{Ival: p.cur.Ival}
		p.advance()
		return val, nil
	}
	if name, ok := p.parseIdentifier(); ok {
		return &nodes.String{Str: name}, nil
	}
	return &nodes.String{Str: ""}, nil
}

// queryHintForceDisableSuffix defines valid suffixes after FORCE/DISABLE in query hints.
// FORCE ORDER is handled separately; the remaining are EXTERNALPUSHDOWN and SCALEOUTEXECUTION.
var queryHintForceDisableSuffix = newOptionSet().withIdents(
	"EXTERNALPUSHDOWN",
	"SCALEOUTEXECUTION",
)

// queryHintParameterizationMode defines valid modes after PARAMETERIZATION.
var queryHintParameterizationMode = newOptionSet().withIdents(
	"SIMPLE",
	"FORCED",
)

// queryHintDefaultNames defines valid hint names for the default (name = value) path
// in OPTION clauses. These are hints from SqlScriptDOM that are not handled by
// dedicated switch cases (RECOMPILE, OPTIMIZE, LOOP/HASH/MERGE, CONCAT, ORDER,
// FORCE, MAXDOP, MAXRECURSION, FAST, QUERYTRACEON, EXPAND, KEEP, KEEPFIXED,
// ROBUST, PARAMETERIZATION, USE, DISABLE, TABLE are all handled above).
//
// From MonoOptimizerHintHelper: IGNORE_NONCLUSTERED_COLUMNSTORE_INDEX, NO_PERFORMANCE_SPOOL
// From IntegerOptimizerHintHelper: CARDINALITY_TUNER_LIMIT
// From DoubleOptimizerHintHelper: MAX_GRANT_PERCENT, MIN_GRANT_PERCENT
// From StringOptimizerHintHelper: LABEL
// From PlanOptimizerHintHelper: SHRINKDB, ALTERCOLUMN, CHECKCONSTRAINTS
var queryHintDefaultNames = newOptionSet().withIdents(
	"IGNORE_NONCLUSTERED_COLUMNSTORE_INDEX",
	"NO_PERFORMANCE_SPOOL",
	"CARDINALITY_TUNER_LIMIT",
	"MAX_GRANT_PERCENT",
	"MIN_GRANT_PERCENT",
	"LABEL",
	"SHRINKDB",
	"ALTERCOLUMN",
	"CHECKCONSTRAINTS",
)

// parseOptionClause parses OPTION ( query_hint [,...n] ).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/option-clause-transact-sql
//
//	OPTION ( query_hint [ ,...n ] )
//
//	query_hint ::=
//	    { HASH | ORDER } GROUP
//	  | { CONCAT | HASH | MERGE } UNION
//	  | { LOOP | MERGE | HASH } JOIN
//	  | EXPAND VIEWS
//	  | FAST number_rows
//	  | FORCE ORDER
//	  | { FORCE | DISABLE } EXTERNALPUSHDOWN
//	  | { FORCE | DISABLE } SCALEOUTEXECUTION
//	  | IGNORE_NONCLUSTERED_COLUMNSTORE_INDEX
//	  | KEEP PLAN
//	  | KEEPFIXED PLAN
//	  | MAX_GRANT_PERCENT = percent
//	  | MIN_GRANT_PERCENT = percent
//	  | MAXDOP number_of_processors
//	  | MAXRECURSION number
//	  | NO_PERFORMANCE_SPOOL
//	  | OPTIMIZE FOR ( @variable_name { UNKNOWN | = literal } [ , ...n ] )
//	  | OPTIMIZE FOR UNKNOWN
//	  | PARAMETERIZATION { SIMPLE | FORCED }
//	  | QUERYTRACEON trace_flag
//	  | RECOMPILE
//	  | ROBUST PLAN
//	  | USE HINT ( 'hint_name' [ , ...n ] )
//	  | USE PLAN N'xml_plan'
//	  | TABLE HINT ( exposed_object_name [ , hint [ , ...n ] ] )
func (p *Parser) parseOptionClause() (*nodes.List, error) {
	p.advance() // consume OPTION
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside OPTION (...) → query hint keywords
	if p.collectMode() {
		p.addRuleCandidate("query_hint")
		return nil, errCollecting
	}

	var hints []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		hint, err := p.parseQueryHint()
		if err != nil {
			return nil, err
		}
		if hint != nil {
			hints = append(hints, hint)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(hints) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: hints}, nil
}

// parseQueryHint parses a single query hint within an OPTION clause.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-query
//
//	query_hint ::=
//	    { HASH | ORDER } GROUP
//	  | { CONCAT | HASH | MERGE } UNION
//	  | { LOOP | MERGE | HASH } JOIN
//	  | EXPAND VIEWS
//	  | FAST number_rows
//	  | FORCE ORDER
//	  | { FORCE | DISABLE } EXTERNALPUSHDOWN
//	  | { FORCE | DISABLE } SCALEOUTEXECUTION
//	  | IGNORE_NONCLUSTERED_COLUMNSTORE_INDEX
//	  | KEEP PLAN
//	  | KEEPFIXED PLAN
//	  | MAX_GRANT_PERCENT = percent
//	  | MIN_GRANT_PERCENT = percent
//	  | MAXDOP number_of_processors
//	  | MAXRECURSION number
//	  | NO_PERFORMANCE_SPOOL
//	  | OPTIMIZE FOR ( @variable_name { UNKNOWN | = literal } [ , ...n ] )
//	  | OPTIMIZE FOR UNKNOWN
//	  | PARAMETERIZATION { SIMPLE | FORCED }
//	  | QUERYTRACEON trace_flag
//	  | RECOMPILE
//	  | ROBUST PLAN
//	  | USE HINT ( 'hint_name' [ , ...n ] )
//	  | USE PLAN N'xml_plan'
//	  | TABLE HINT ( exposed_object_name [ , <table_hint> [ , ...n ] ] )
func (p *Parser) parseQueryHint() (nodes.Node, error) {
	loc := p.pos()

	switch {
	case p.cur.Type == kwRECOMPILE:
		p.advance()
		return &nodes.QueryHint{Kind: "RECOMPILE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwOPTIMIZE:
		p.advance()
		if p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type == kwUNKNOWN {
				p.advance()
				return &nodes.QueryHint{Kind: "OPTIMIZE FOR UNKNOWN", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
			}
			if p.cur.Type == '(' {
				p.advance()
				var params []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					param, err := p.parseOptimizeForParam()
					if err != nil {
						return nil, err
					}
					if param != nil {
						params = append(params, param)
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				hint := &nodes.QueryHint{Kind: "OPTIMIZE FOR", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
				if len(params) > 0 {
					hint.Params = &nodes.List{Items: params}
				}
				return hint, nil
			}
		}
		return &nodes.QueryHint{Kind: "OPTIMIZE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwLOOP || p.cur.Type == kwHASH || p.cur.Type == kwMERGE:
		prefix := strings.ToUpper(p.cur.Str)
		p.advance()
		if p.cur.Type == kwJOIN {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " JOIN", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		} else if p.cur.Type == kwUNION {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " UNION", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		} else if p.cur.Type == kwGROUP {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " GROUP", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: prefix, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwCONCAT:
		p.advance()
		if p.cur.Type == kwUNION {
			p.advance()
			return &nodes.QueryHint{Kind: "CONCAT UNION", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "CONCAT", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwORDER:
		p.advance()
		if p.cur.Type == kwGROUP {
			p.advance()
			return &nodes.QueryHint{Kind: "ORDER GROUP", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "ORDER", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwFORCE:
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
			return &nodes.QueryHint{Kind: "FORCE ORDER", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		} else if p.isValidOption(queryHintForceDisableSuffix) {
			suffix := strings.ToUpper(p.cur.Str)
			p.advance()
			return &nodes.QueryHint{Kind: "FORCE " + suffix, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "FORCE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwMAXDOP:
		p.advance()
		hint := &nodes.QueryHint{Kind: "MAXDOP", Loc: nodes.Loc{Start: loc, End: -1}}
		if p.cur.Type == tokICONST {
			var err error
			hint.Value, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case p.cur.Type == kwMAXRECURSION:
		p.advance()
		hint := &nodes.QueryHint{Kind: "MAXRECURSION", Loc: nodes.Loc{Start: loc, End: -1}}
		if p.cur.Type == tokICONST {
			var err error
			hint.Value, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case p.cur.Type == kwFAST:
		p.advance()
		hint := &nodes.QueryHint{Kind: "FAST", Loc: nodes.Loc{Start: loc, End: -1}}
		if p.cur.Type == tokICONST {
			var err error
			hint.Value, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case p.cur.Type == kwQUERYTRACEON:
		p.advance()
		hint := &nodes.QueryHint{Kind: "QUERYTRACEON", Loc: nodes.Loc{Start: loc, End: -1}}
		if p.cur.Type == tokICONST {
			var err error
			hint.Value, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case p.cur.Type == kwEXPAND:
		p.advance()
		if p.cur.Type == kwVIEWS {
			p.advance()
			return &nodes.QueryHint{Kind: "EXPAND VIEWS", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "EXPAND", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwKEEP:
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "KEEP PLAN", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "KEEP", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwKEEPFIXED:
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "KEEPFIXED PLAN", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "KEEPFIXED", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwROBUST:
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "ROBUST PLAN", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "ROBUST", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwPARAMETERIZATION:
		p.advance()
		hint := &nodes.QueryHint{Kind: "PARAMETERIZATION", Loc: nodes.Loc{Start: loc, End: -1}}
		if p.isValidOption(queryHintParameterizationMode) {
			hint.StrValue = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		hint.Loc.End = p.prevEnd()
		return hint, nil

	case p.cur.Type == kwUSE:
		p.advance()
		if p.cur.Type == kwHINT {
			p.advance()
			hint := &nodes.QueryHint{Kind: "USE HINT", Loc: nodes.Loc{Start: loc, End: -1}}
			if p.cur.Type == '(' {
				p.advance()
				var hintNames []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
						hintNames = append(hintNames, &nodes.String{Str: p.cur.Str})
						p.advance()
					}
					p.match(',')
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				if len(hintNames) > 0 {
					hint.Params = &nodes.List{Items: hintNames}
				}
			}
			hint.Loc.End = p.prevEnd()
			return hint, nil
		} else if p.cur.Type == kwPLAN {
			p.advance()
			hint := &nodes.QueryHint{Kind: "USE PLAN", Loc: nodes.Loc{Start: loc, End: -1}}
			if p.cur.Type == tokNSCONST || p.cur.Type == tokSCONST {
				hint.StrValue = p.cur.Str
				p.advance()
			}
			hint.Loc.End = p.prevEnd()
			return hint, nil
		}
		return &nodes.QueryHint{Kind: "USE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwDISABLE:
		p.advance()
		if p.isValidOption(queryHintForceDisableSuffix) {
			suffix := strings.ToUpper(p.cur.Str)
			p.advance()
			return &nodes.QueryHint{Kind: "DISABLE " + suffix, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
		}
		return &nodes.QueryHint{Kind: "DISABLE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	case p.cur.Type == kwTABLE:
		p.advance()
		if p.cur.Type == kwHINT {
			p.advance()
			hint := &nodes.QueryHint{Kind: "TABLE HINT", Loc: nodes.Loc{Start: loc, End: -1}}
			if p.cur.Type == '(' {
				p.advance()
				// Parse exposed_object_name as a TableRef
				var err error
				hint.TableName, err = p.parseTableRef()
				if err != nil {
					return nil, err
				}
				// After the table reference, zero-or-more hints may follow, each
				// introduced by a leading comma. When the first comma is present
				// the tail is a normal strict comma-separated list terminated by
				// ')' — empty slots and trailing commas are rejected.
				if _, ok := p.match(','); ok {
					hints, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
						th, err := p.parseTableHint()
						if err != nil {
							return nil, err
						}
						if th == nil {
							return nil, p.unexpectedToken()
						}
						return th, nil
					})
					if err != nil {
						return nil, err
					}
					hint.TableHints = &nodes.List{Items: hints}
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
			hint.Loc.End = p.prevEnd()
			return hint, nil
		}
		return &nodes.QueryHint{Kind: "TABLE", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil

	default:
		// Known hint with name = value pattern (e.g., MAX_GRANT_PERCENT = 10)
		if p.isValidOption(queryHintDefaultNames) {
			name := strings.ToUpper(p.cur.Str)
			p.advance()
			hint := &nodes.QueryHint{Kind: name, Loc: nodes.Loc{Start: loc, End: -1}}
			if p.cur.Type == '=' {
				p.advance()
				var err error
				hint.Value, err = p.parsePrimary()
				if err != nil {
					return nil, err
				}
			}
			hint.Loc.End = p.prevEnd()
			return hint, nil
		}
		return nil, nil
	}
}

// parseOptimizeForParam parses a single OPTIMIZE FOR parameter.
//
//	@variable_name { UNKNOWN | = literal_constant }
func (p *Parser) parseOptimizeForParam() (*nodes.OptimizeForParam, error) {
	if p.cur.Type != tokVARIABLE {
		return nil, nil
	}
	loc := p.pos()
	param := &nodes.OptimizeForParam{
		Variable: p.cur.Str,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}
	p.advance()
	if p.cur.Type == '=' {
		p.advance()
		var err error
		param.Value, err = p.parsePrimary()
		if err != nil {
			return nil, err
		}
	} else if p.cur.Type == kwUNKNOWN {
		param.Unknown = true
		p.advance()
	}
	param.Loc.End = p.prevEnd()
	return param, nil
}
