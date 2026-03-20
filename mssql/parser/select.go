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
func (p *Parser) parseSelectStmt() *nodes.SelectStmt {
	loc := p.pos()

	// WITH clause (CTE)
	var withClause *nodes.WithClause
	if p.cur.Type == kwWITH {
		withClause = p.parseWithClause()
	}

	if p.cur.Type != kwSELECT {
		return nil
	}
	p.advance() // consume SELECT

	stmt := &nodes.SelectStmt{
		WithClause: withClause,
		Loc:        nodes.Loc{Start: loc},
	}

	// ALL | DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		stmt.Distinct = true
	} else if _, ok := p.match(kwALL); ok {
		stmt.All = true
	}

	// TOP clause
	if p.cur.Type == kwTOP {
		stmt.Top = p.parseTopClause()
	}

	// Target list
	stmt.TargetList = p.parseTargetList()

	// INTO
	if _, ok := p.match(kwINTO); ok {
		stmt.IntoTable , _ = p.parseTableRef()
	}

	// FROM
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause = p.parseFromClause()
	}

	// WHERE
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause, _ = p.parseExpr()
	}

	// GROUP BY [ALL]
	if p.cur.Type == kwGROUP {
		p.advance()
		if _, err := p.expect(kwBY); err == nil {
			if _, ok := p.match(kwALL); ok {
				stmt.GroupByAll = true
			}
			stmt.GroupByClause = p.parseGroupByList()
		}
	}

	// HAVING
	if _, ok := p.match(kwHAVING); ok {
		stmt.HavingClause, _ = p.parseExpr()
	}

	// WINDOW clause (named window definitions)
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "WINDOW") {
		stmt.WindowClause = p.parseWindowClause()
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err == nil {
			stmt.OrderByClause = p.parseOrderByList()
		}
	}

	// OFFSET ... FETCH
	if p.cur.Type == kwOFFSET {
		p.advance()
		stmt.OffsetClause, _ = p.parseExpr()
		// Consume optional ROWS/ROW
		if p.cur.Type == kwROWS || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "row")) {
			p.advance()
		}
		// FETCH NEXT n ROWS ONLY
		if p.cur.Type == kwFETCH {
			fetchLoc := p.pos()
			p.advance()
			// NEXT or FIRST
			if p.cur.Type == tokIDENT && (strings.EqualFold(p.cur.Str, "next") || strings.EqualFold(p.cur.Str, "first")) {
				p.advance()
			}
			count, _ := p.parseExpr()
			// ROWS/ROW
			if p.cur.Type == kwROWS || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "row")) {
				p.advance()
			}
			// ONLY
			if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "only") {
				p.advance()
			}
			stmt.FetchClause = &nodes.FetchClause{
				Count: count,
				Loc:   nodes.Loc{Start: fetchLoc},
			}
		}
	}

	// FOR XML / FOR JSON / FOR BROWSE
	if p.cur.Type == kwFOR {
		next := p.peekNext()
		if next.Type == kwXML || next.Type == kwJSON || next.Type == kwBROWSE {
			stmt.ForClause = p.parseForClause()
		}
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		stmt.OptionClause = p.parseOptionClause()
	}

	// UNION / INTERSECT / EXCEPT
	if p.cur.Type == kwUNION || p.cur.Type == kwINTERSECT || p.cur.Type == kwEXCEPT {
		return p.parseSetOperation(stmt)
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSetOperation parses UNION/INTERSECT/EXCEPT.
func (p *Parser) parseSetOperation(left *nodes.SelectStmt) *nodes.SelectStmt {
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
	}

	right := p.parseSelectStmt()
	return &nodes.SelectStmt{
		Op:   op,
		All:  all,
		Larg: left,
		Rarg: right,
		Loc:  left.Loc,
	}
}

// parseWithClause parses WITH [XMLNAMESPACES(...),] cte_name [(col_list)] AS (select), ...
//
// BNF: mssql/parser/bnf/select-transact-sql.bnf
//
//	[ WITH { [ XMLNAMESPACES , ] [ <common_table_expression> [ , ...n ] ] } ]
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/with-common-table-expression-transact-sql
func (p *Parser) parseWithClause() *nodes.WithClause {
	loc := p.pos()
	p.advance() // consume WITH

	wc := &nodes.WithClause{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional XMLNAMESPACES (...)
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "XMLNAMESPACES") {
		wc.XmlNamespaces = p.parseXmlNamespaces()
		p.match(',') // consume comma between XMLNAMESPACES and CTEs
	}

	var ctes []nodes.Node
	for p.cur.Type != kwSELECT && p.cur.Type != tokEOF && p.cur.Type != ';' {
		cte := p.parseCTE()
		if cte == nil {
			break
		}
		ctes = append(ctes, cte)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	wc.CTEs = &nodes.List{Items: ctes}
	wc.Loc.End = p.pos()
	return wc
}

// parseXmlNamespaces parses XMLNAMESPACES ( namespace_decl [, ...n] ).
//
//	XMLNAMESPACES ( uri AS prefix [, ...n] | DEFAULT uri [, ...n] )
func (p *Parser) parseXmlNamespaces() *nodes.List {
	p.advance() // consume XMLNAMESPACES
	if _, err := p.expect('('); err != nil {
		return nil
	}

	var decls []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		loc := p.pos()
		decl := &nodes.XmlNamespaceDecl{Loc: nodes.Loc{Start: loc}}

		if _, ok := p.match(kwDEFAULT); ok {
			// DEFAULT 'uri'
			decl.IsDefault = true
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				decl.URI = p.cur.Str
				p.advance()
			}
		} else {
			// 'uri' AS prefix
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
		decl.Loc.End = p.pos()
		decls = append(decls, decl)

		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')

	if len(decls) == 0 {
		return nil
	}
	return &nodes.List{Items: decls}
}

// parseCTE parses a single CTE: name [(columns)] AS (query).
func (p *Parser) parseCTE() *nodes.CommonTableExpr {
	loc := p.pos()
	name, ok := p.parseIdentifier()
	if !ok {
		return nil
	}

	cte := &nodes.CommonTableExpr{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			colName, ok := p.parseIdentifier()
			if !ok {
				break
			}
			cols = append(cols, &nodes.String{Str: colName})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		_, _ = p.expect(')')
		cte.Columns = &nodes.List{Items: cols}
	}

	if _, err := p.expect(kwAS); err != nil {
		return cte
	}
	if _, err := p.expect('('); err != nil {
		return cte
	}

	cte.Query = p.parseSelectStmt()
	_, _ = p.expect(')')

	cte.Loc.End = p.pos()
	return cte
}

// parseTopClause parses TOP (expr) [PERCENT] [WITH TIES].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/top-transact-sql
func (p *Parser) parseTopClause() *nodes.TopClause {
	loc := p.pos()
	p.advance() // consume TOP

	tc := &nodes.TopClause{
		Loc: nodes.Loc{Start: loc},
	}

	// TOP (expr) or TOP literal
	if p.cur.Type == '(' {
		p.advance()
		tc.Count, _ = p.parseExpr()
		_, _ = p.expect(')')
	} else {
		tc.Count, _ = p.parsePrimary()
	}

	// PERCENT
	if _, ok := p.match(kwPERCENT); ok {
		tc.Percent = true
	}

	// WITH TIES
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == tokIDENT && strings.EqualFold(next.Str, "ties") {
			p.advance() // consume WITH
			p.advance() // consume TIES
			tc.WithTies = true
		}
	}

	tc.Loc.End = p.pos()
	return tc
}

// parseTargetList parses a comma-separated list of result columns.
func (p *Parser) parseTargetList() *nodes.List {
	var targets []nodes.Node
	for {
		targetLoc := p.pos()
		expr, _ := p.parseExpr()
		if expr == nil {
			break
		}
		target := &nodes.ResTarget{
			Val: expr,
			Loc: nodes.Loc{Start: targetLoc},
		}
		// Check for alias: AS name or just name (but not keywords that start clauses)
		if _, ok := p.match(kwAS); ok {
			if p.isIdentLike() {
				target.Name = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == tokIDENT {
			target.Name = p.cur.Str
			p.advance()
		}
		target.Loc.End = p.pos()
		targets = append(targets, target)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: targets}
}

// parseFromClause parses a FROM clause table source list.
func (p *Parser) parseFromClause() *nodes.List {
	var sources []nodes.Node
	for {
		source := p.parseTableSource()
		if source == nil {
			break
		}
		sources = append(sources, source)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: sources}
}

// parseTableSource parses a single table source (table, subquery, join).
func (p *Parser) parseTableSource() nodes.TableExpr {
	left := p.parsePrimaryTableSource()
	if left == nil {
		return nil
	}

	// Parse joins
	for {
		jt, ok := p.matchJoinType()
		if !ok {
			break
		}
		right := p.parsePrimaryTableSource()
		join := &nodes.JoinClause{
			Type:  jt,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: p.pos()},
		}
		// ON condition (not for CROSS JOIN / CROSS APPLY / OUTER APPLY)
		if jt != nodes.JoinCross && jt != nodes.JoinCrossApply && jt != nodes.JoinOuterApply {
			if _, ok := p.match(kwON); ok {
				join.Condition, _ = p.parseExpr()
			}
		}
		left = join
	}

	return left
}

// parsePrimaryTableSource parses a base table, subquery, or function call as table source.
func (p *Parser) parsePrimaryTableSource() nodes.TableExpr {
	// Subquery: (SELECT ...)
	if p.cur.Type == '(' {
		loc := p.pos()
		p.advance()
		sub := p.parseSelectStmt()
		_, _ = p.expect(')')

		subExpr := &nodes.SubqueryExpr{
			Query: sub,
			Loc:   nodes.Loc{Start: loc},
		}

		// Alias
		alias := p.parseOptionalAlias()
		if alias != "" {
			result := &nodes.AliasedTableRef{
				Table: subExpr,
				Alias: alias,
				Loc:   nodes.Loc{Start: loc},
			}
			return p.parsePivotUnpivot(result)
		}
		return p.parsePivotUnpivot(subExpr)
	}

	// Rowset functions: OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
	if p.cur.Type == kwOPENROWSET || p.cur.Type == kwOPENQUERY || p.cur.Type == kwOPENJSON ||
		p.cur.Type == kwOPENDATASOURCE || p.cur.Type == kwOPENXML {
		return p.parseRowsetFunction()
	}

	// Table reference
	ref , _ := p.parseTableRef()
	if ref == nil {
		return nil
	}

	// Check if this is a function call (table-valued function)
	if p.cur.Type == '(' {
		return p.parsePivotUnpivot(p.parseTableValuedFunction(ref))
	}

	// TABLESAMPLE
	if p.cur.Type == kwTABLESAMPLE {
		ts := p.parseTableSampleClause()
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
			result.Hints = p.parseTableHints()
		}
		return p.parsePivotUnpivot(result)
	}

	// Alias
	alias := p.parseOptionalAlias()
	if alias != "" {
		ref.Alias = alias
	}

	// Table hints: WITH (NOLOCK), WITH (INDEX(idx1)), etc.
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		ref.Hints = p.parseTableHints()
	}

	return p.parsePivotUnpivot(ref)
}

// parsePivotUnpivot checks for and parses PIVOT or UNPIVOT after a table source.
func (p *Parser) parsePivotUnpivot(source nodes.TableExpr) nodes.TableExpr {
	if p.cur.Type == kwPIVOT {
		return p.parsePivotExpr(source)
	}
	if p.cur.Type == kwUNPIVOT {
		return p.parseUnpivotExpr(source)
	}
	return source
}

// parsePivotExpr parses PIVOT (agg_func(col) FOR pivot_col IN ([v1],[v2],...)) AS alias.
func (p *Parser) parsePivotExpr(source nodes.TableExpr) *nodes.PivotExpr {
	loc := p.pos()
	p.advance() // consume PIVOT

	pivot := &nodes.PivotExpr{
		Source: source,
		Loc:    nodes.Loc{Start: loc},
	}

	if _, err := p.expect('('); err != nil {
		return pivot
	}

	// Parse aggregate function call
	pivot.AggFunc, _ = p.parseExpr()

	// FOR column
	if _, ok := p.match(kwFOR); ok {
		if name, ok := p.parseIdentifier(); ok {
			pivot.ForCol = name
		}
	}

	// IN ([v1], [v2], ...)
	if _, ok := p.match(kwIN); ok {
		if _, err := p.expect('('); err == nil {
			var vals []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if name, ok := p.parseIdentifier(); ok {
					vals = append(vals, &nodes.String{Str: name})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			pivot.InValues = &nodes.List{Items: vals}
		}
	}

	_, _ = p.expect(')') // closing paren of PIVOT(...)

	// AS alias
	alias := p.parseOptionalAlias()
	pivot.Alias = alias

	pivot.Loc.End = p.pos()
	return pivot
}

// parseUnpivotExpr parses UNPIVOT (value_col FOR pivot_col IN ([c1],[c2],...)) AS alias.
func (p *Parser) parseUnpivotExpr(source nodes.TableExpr) *nodes.UnpivotExpr {
	loc := p.pos()
	p.advance() // consume UNPIVOT

	unpivot := &nodes.UnpivotExpr{
		Source: source,
		Loc:    nodes.Loc{Start: loc},
	}

	if _, err := p.expect('('); err != nil {
		return unpivot
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
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if name, ok := p.parseIdentifier(); ok {
					cols = append(cols, &nodes.String{Str: name})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			unpivot.InCols = &nodes.List{Items: cols}
		}
	}

	_, _ = p.expect(')') // closing paren of UNPIVOT(...)

	// AS alias
	alias := p.parseOptionalAlias()
	unpivot.Alias = alias

	unpivot.Loc.End = p.pos()
	return unpivot
}

// parseTableSampleClause parses TABLESAMPLE (size PERCENT|ROWS) [REPEATABLE (seed)].
func (p *Parser) parseTableSampleClause() *nodes.TableSampleClause {
	loc := p.pos()
	p.advance() // consume TABLESAMPLE

	ts := &nodes.TableSampleClause{
		Loc: nodes.Loc{Start: loc},
	}

	if _, err := p.expect('('); err != nil {
		return ts
	}

	ts.Size, _ = p.parseExpr()

	// PERCENT or ROWS
	if _, ok := p.match(kwPERCENT); ok {
		ts.Unit = "PERCENT"
	} else if _, ok := p.match(kwROWS); ok {
		ts.Unit = "ROWS"
	}

	_, _ = p.expect(')')

	// REPEATABLE (seed)
	if p.matchIdentCI("REPEATABLE") {
		if _, err := p.expect('('); err == nil {
			ts.Repeatable, _ = p.parseExpr()
			_, _ = p.expect(')')
		}
	}

	ts.Loc.End = p.pos()
	return ts
}

// parseTableValuedFunction parses a table-valued function call after the name.
func (p *Parser) parseTableValuedFunction(ref *nodes.TableRef) nodes.TableExpr {
	loc := ref.Loc.Start
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: ref,
		Loc:  nodes.Loc{Start: loc},
	}

	if p.cur.Type != ')' {
		var args []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg, _ := p.parseExpr()
			args = append(args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		fc.Args = &nodes.List{Items: args}
	}
	_, _ = p.expect(')')

	alias := p.parseOptionalAlias()
	if alias != "" {
		return &nodes.AliasedTableRef{
			Table: fc,
			Alias: alias,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return &nodes.AliasedTableRef{
		Table: fc,
		Loc:   nodes.Loc{Start: loc},
	}
}

// parseOptionalAlias parses an optional alias (AS name or just name).
func (p *Parser) parseOptionalAlias() string {
	if _, ok := p.match(kwAS); ok {
		if p.isIdentLike() {
			name := p.cur.Str
			p.advance()
			return name
		}
		return ""
	}
	// Bare alias - only if it's an identifier and NOT a clause keyword
	if p.cur.Type == tokIDENT && !p.isSelectClauseIdent() {
		name := p.cur.Str
		p.advance()
		return name
	}
	return ""
}

// isSelectClauseIdent returns true if the current identifier token is a contextual
// keyword that starts a SELECT clause and should not be consumed as a bare alias.
func (p *Parser) isSelectClauseIdent() bool {
	if p.cur.Type != tokIDENT {
		return false
	}
	upper := strings.ToUpper(p.cur.Str)
	return upper == "WINDOW"
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
func (p *Parser) parseForClause() *nodes.ForClause {
	loc := p.pos()
	p.advance() // consume FOR

	fc := &nodes.ForClause{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwBROWSE {
		fc.Mode = nodes.ForBrowse
		p.advance()
		fc.Loc.End = p.pos()
		return fc
	}

	if p.cur.Type == kwXML {
		fc.Mode = nodes.ForXML
		p.advance()
		// RAW, AUTO, EXPLICIT, PATH
		if p.isIdentLike() || p.cur.Type == kwRAW || p.cur.Type == kwPATH {
			fc.SubMode = strings.ToUpper(p.cur.Str)
			p.advance()
			// RAW('ElementName') or PATH('ElementName')
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokSCONST {
					fc.ElementName = p.cur.Str
					p.advance()
				}
				_, _ = p.expect(')')
			}
		}
		// Parse comma-separated XML options
		p.parseForXmlOptions(fc)
	} else if p.cur.Type == kwJSON {
		fc.Mode = nodes.ForJSON
		p.advance()
		// AUTO or PATH
		if p.isIdentLike() || p.cur.Type == kwPATH {
			fc.SubMode = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		// Parse comma-separated JSON options
		p.parseForJsonOptions(fc)
	}

	fc.Loc.End = p.pos()
	return fc
}

// parseForXmlOptions parses the comma-separated options after FOR XML {RAW|AUTO|EXPLICIT|PATH}.
//
//	[ , BINARY BASE64 ]
//	[ , TYPE ]
//	[ , ROOT [ ( 'RootName' ) ] ]
//	[ , { XMLDATA | XMLSCHEMA [ ( 'TargetNameSpaceURI' ) ] } ]
//	[ , ELEMENTS [ XSINIL | ABSENT ] ]
func (p *Parser) parseForXmlOptions(fc *nodes.ForClause) {
	for {
		if _, ok := p.match(','); !ok {
			return
		}
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
			return
		}
	}
}

// parseForJsonOptions parses the comma-separated options after FOR JSON {AUTO|PATH}.
//
//	[ , ROOT [ ( 'RootName' ) ] ]
//	[ , INCLUDE_NULL_VALUES ]
//	[ , WITHOUT_ARRAY_WRAPPER ]
func (p *Parser) parseForJsonOptions(fc *nodes.ForClause) {
	for {
		if _, ok := p.match(','); !ok {
			return
		}
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
			return
		}
	}
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() *nodes.List {
	var items []nodes.Node
	for {
		expr, _ := p.parseExpr()
		if expr == nil {
			break
		}
		items = append(items, expr)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}

// parseGroupByList parses a GROUP BY list which may contain GROUPING SETS, ROLLUP, CUBE.
func (p *Parser) parseGroupByList() *nodes.List {
	var items []nodes.Node
	for {
		// GROUPING SETS (...)
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "GROUPING") {
			next := p.peekNext()
			if next.Type == tokIDENT && strings.EqualFold(next.Str, "SETS") {
				loc := p.pos()
				p.advance() // consume GROUPING
				p.advance() // consume SETS
				if _, err := p.expect('('); err == nil {
					var sets []nodes.Node
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						set := p.parseGroupingSet()
						sets = append(sets, set)
						if _, ok := p.match(','); !ok {
							break
						}
					}
					_, _ = p.expect(')')
					items = append(items, &nodes.GroupingSetsExpr{
						Sets: &nodes.List{Items: sets},
						Loc:  nodes.Loc{Start: loc, End: p.pos()},
					})
				}
				if _, ok := p.match(','); !ok {
					break
				}
				continue
			}
		}
		// ROLLUP (...)
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "ROLLUP") {
			loc := p.pos()
			p.advance() // consume ROLLUP
			if _, err := p.expect('('); err == nil {
				args := p.parseExprList()
				_, _ = p.expect(')')
				items = append(items, &nodes.RollupExpr{
					Args: args,
					Loc:  nodes.Loc{Start: loc, End: p.pos()},
				})
			}
			if _, ok := p.match(','); !ok {
				break
			}
			continue
		}
		// CUBE (...)
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CUBE") {
			loc := p.pos()
			p.advance() // consume CUBE
			if _, err := p.expect('('); err == nil {
				args := p.parseExprList()
				_, _ = p.expect(')')
				items = append(items, &nodes.CubeExpr{
					Args: args,
					Loc:  nodes.Loc{Start: loc, End: p.pos()},
				})
			}
			if _, ok := p.match(','); !ok {
				break
			}
			continue
		}
		// Regular expression
		expr, _ := p.parseExpr()
		if expr == nil {
			break
		}
		items = append(items, expr)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
}

// parseGroupingSet parses a single grouping set: () or (expr, expr, ...) or just expr.
func (p *Parser) parseGroupingSet() *nodes.List {
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			// Empty set ()
			p.advance()
			return &nodes.List{Items: nil}
		}
		var items []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			expr, _ := p.parseExpr()
			if expr != nil {
				items = append(items, expr)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		_, _ = p.expect(')')
		return &nodes.List{Items: items}
	}
	// Single expression as a set
	expr, _ := p.parseExpr()
	if expr != nil {
		return &nodes.List{Items: []nodes.Node{expr}}
	}
	return &nodes.List{Items: nil}
}

// parseWindowClause parses WINDOW window_name AS (window_spec) [, ...].
//
//	WINDOW window_name AS ( [ existing_window_name ]
//	    [ PARTITION BY expr [,...n] ]
//	    [ ORDER BY order_item [,...n] ]
//	    [ <window_frame> ]
//	)
func (p *Parser) parseWindowClause() *nodes.List {
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
			Loc:  nodes.Loc{Start: loc},
		}

		if _, err := p.expect(kwAS); err != nil {
			defs = append(defs, def)
			break
		}
		if _, err := p.expect('('); err != nil {
			defs = append(defs, def)
			break
		}

		// Optional existing_window_name (must be an ident not followed by keyword like PARTITION, ORDER)
		if p.cur.Type == tokIDENT && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
			p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
			// Check if this looks like a reference name (ident not a keyword)
			if !strings.EqualFold(p.cur.Str, "PARTITION") &&
				!strings.EqualFold(p.cur.Str, "ORDER") &&
				!strings.EqualFold(p.cur.Str, "ROWS") &&
				!strings.EqualFold(p.cur.Str, "RANGE") &&
				!strings.EqualFold(p.cur.Str, "GROUPS") &&
				p.cur.Type != ')' {
				next := p.peekNext()
				// If next token is a clause keyword or ), this is a refname
				if next.Type == kwPARTITION || next.Type == kwORDER ||
					next.Type == kwROWS || next.Type == kwRANGE || next.Type == kwGROUPS ||
					next.Type == ')' {
					def.RefName = p.cur.Str
					p.advance()
				}
			}
		}

		// PARTITION BY
		if p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect(kwBY); err == nil {
				def.PartitionBy = p.parseExprList()
			}
		}

		// ORDER BY
		if p.cur.Type == kwORDER {
			p.advance()
			if _, err := p.expect(kwBY); err == nil {
				def.OrderBy = p.parseOrderByList()
			}
		}

		// Window frame: ROWS | RANGE | GROUPS
		if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
			def.Frame, _ = p.parseWindowFrame()
		}

		_, _ = p.expect(')')
		def.Loc.End = p.pos()
		defs = append(defs, def)

		if _, ok := p.match(','); !ok {
			break
		}
	}

	if len(defs) == 0 {
		return nil
	}
	return &nodes.List{Items: defs}
}

// parseOrderByList parses ORDER BY items.
func (p *Parser) parseOrderByList() *nodes.List {
	var items []nodes.Node
	for {
		oloc := p.pos()
		expr, _ := p.parseExpr()
		if expr == nil {
			break
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		items = append(items, &nodes.OrderByItem{
			Expr:    expr,
			SortDir: dir,
			Loc:     nodes.Loc{Start: oloc},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
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
func (p *Parser) parseTableHints() *nodes.List {
	p.advance() // consume WITH
	if _, err := p.expect('('); err != nil {
		return nil
	}

	var hints []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		hint := p.parseTableHint()
		if hint == nil {
			break
		}
		hints = append(hints, hint)
		// Optional comma between hints
		p.match(',')
	}
	_, _ = p.expect(')')

	if len(hints) == 0 {
		return nil
	}
	return &nodes.List{Items: hints}
}

// parseTableHint parses a single table hint.
func (p *Parser) parseTableHint() *nodes.TableHint {
	loc := p.pos()

	// INDEX hint: INDEX ( values ) or INDEX = ( value )
	if p.cur.Type == kwINDEX {
		p.advance()
		hint := &nodes.TableHint{
			Name: "INDEX",
			Loc:  nodes.Loc{Start: loc},
		}
		if _, ok := p.match('='); ok {
			// INDEX = ( value )
			if _, err := p.expect('('); err == nil {
				var vals []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					vals = append(vals, p.parseIndexValue())
					if _, ok := p.match(','); !ok {
						break
					}
				}
				_, _ = p.expect(')')
				hint.IndexValues = &nodes.List{Items: vals}
			}
		} else if p.cur.Type == '(' {
			// INDEX ( values )
			p.advance()
			var vals []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				vals = append(vals, p.parseIndexValue())
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			hint.IndexValues = &nodes.List{Items: vals}
		}
		hint.Loc.End = p.pos()
		return hint
	}

	// Check for keyword-based hints that are lexer keywords
	switch p.cur.Type {
	case kwHOLDLOCK:
		p.advance()
		return &nodes.TableHint{Name: "HOLDLOCK", Loc: nodes.Loc{Start: loc, End: p.pos()}}
	case kwNOLOCK:
		p.advance()
		return &nodes.TableHint{Name: "NOLOCK", Loc: nodes.Loc{Start: loc, End: p.pos()}}
	case kwNOWAIT:
		p.advance()
		return &nodes.TableHint{Name: "NOWAIT", Loc: nodes.Loc{Start: loc, End: p.pos()}}
	}

	// All remaining hints are identifiers (not lexer keywords)
	if !p.isIdentLike() {
		return nil
	}

	name := strings.ToUpper(p.cur.Str)
	switch name {
	case "FORCESEEK":
		p.advance()
		hint := &nodes.TableHint{
			Name: "FORCESEEK",
			Loc:  nodes.Loc{Start: loc},
		}
		// Optional: FORCESEEK ( index_value ( col1, col2, ... ) )
		if p.cur.Type == '(' {
			p.advance()
			// index value
			idxVal := p.parseIndexValue()
			hint.IndexValues = &nodes.List{Items: []nodes.Node{idxVal}}
			// ( col1, col2, ... )
			if p.cur.Type == '(' {
				p.advance()
				var cols []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if colName, ok := p.parseIdentifier(); ok {
						cols = append(cols, &nodes.String{Str: colName})
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				_, _ = p.expect(')')
				hint.ForceSeekColumns = &nodes.List{Items: cols}
			}
			_, _ = p.expect(')') // outer paren
		}
		hint.Loc.End = p.pos()
		return hint

	case "SPATIAL_WINDOW_MAX_CELLS":
		p.advance()
		hint := &nodes.TableHint{
			Name: "SPATIAL_WINDOW_MAX_CELLS",
			Loc:  nodes.Loc{Start: loc},
		}
		if _, ok := p.match('='); ok {
			hint.IntValue, _ = p.parsePrimary()
		}
		hint.Loc.End = p.pos()
		return hint

	case "FORCESCAN", "NOEXPAND",
		"PAGLOCK", "READCOMMITTED", "READCOMMITTEDLOCK",
		"READPAST", "READUNCOMMITTED", "REPEATABLEREAD",
		"ROWLOCK", "SERIALIZABLE", "SNAPSHOT",
		"TABLOCK", "TABLOCKX", "UPDLOCK", "XLOCK",
		"KEEPIDENTITY", "KEEPDEFAULTS", "IGNORE_CONSTRAINTS", "IGNORE_TRIGGERS":
		p.advance()
		return &nodes.TableHint{Name: name, Loc: nodes.Loc{Start: loc, End: p.pos()}}

	default:
		return nil
	}
}

// parseIndexValue parses an index value (identifier or integer).
func (p *Parser) parseIndexValue() nodes.Node {
	if p.cur.Type == tokICONST {
		val := &nodes.Integer{Ival: p.cur.Ival}
		p.advance()
		return val
	}
	if name, ok := p.parseIdentifier(); ok {
		return &nodes.String{Str: name}
	}
	return &nodes.String{Str: ""}
}

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
func (p *Parser) parseOptionClause() *nodes.List {
	p.advance() // consume OPTION
	if _, err := p.expect('('); err != nil {
		return nil
	}

	var hints []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		hint := p.parseQueryHint()
		if hint != nil {
			hints = append(hints, hint)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')

	if len(hints) == 0 {
		return nil
	}
	return &nodes.List{Items: hints}
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
func (p *Parser) parseQueryHint() nodes.Node {
	loc := p.pos()

	switch {
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "RECOMPILE"):
		p.advance()
		return &nodes.QueryHint{Kind: "RECOMPILE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "OPTIMIZE"):
		p.advance()
		if p.cur.Type == kwFOR {
			p.advance()
			if p.isIdentLike() && strings.EqualFold(p.cur.Str, "UNKNOWN") {
				p.advance()
				return &nodes.QueryHint{Kind: "OPTIMIZE FOR UNKNOWN", Loc: nodes.Loc{Start: loc, End: p.pos()}}
			}
			if p.cur.Type == '(' {
				p.advance()
				var params []nodes.Node
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					param := p.parseOptimizeForParam()
					if param != nil {
						params = append(params, param)
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				_, _ = p.expect(')')
				hint := &nodes.QueryHint{Kind: "OPTIMIZE FOR", Loc: nodes.Loc{Start: loc, End: p.pos()}}
				if len(params) > 0 {
					hint.Params = &nodes.List{Items: params}
				}
				return hint
			}
		}
		return &nodes.QueryHint{Kind: "OPTIMIZE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && (strings.EqualFold(p.cur.Str, "LOOP") ||
		strings.EqualFold(p.cur.Str, "HASH") ||
		strings.EqualFold(p.cur.Str, "MERGE")):
		prefix := strings.ToUpper(p.cur.Str)
		p.advance()
		if p.cur.Type == kwJOIN {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " JOIN", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		} else if p.cur.Type == kwUNION {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " UNION", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		} else if p.cur.Type == kwGROUP {
			p.advance()
			return &nodes.QueryHint{Kind: prefix + " GROUP", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: prefix, Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "CONCAT"):
		p.advance()
		if p.cur.Type == kwUNION {
			p.advance()
			return &nodes.QueryHint{Kind: "CONCAT UNION", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "CONCAT", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.cur.Type == kwORDER:
		p.advance()
		if p.cur.Type == kwGROUP {
			p.advance()
			return &nodes.QueryHint{Kind: "ORDER GROUP", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "ORDER", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "FORCE"):
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
			return &nodes.QueryHint{Kind: "FORCE ORDER", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		} else if p.isIdentLike() {
			suffix := strings.ToUpper(p.cur.Str)
			p.advance()
			return &nodes.QueryHint{Kind: "FORCE " + suffix, Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "FORCE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "MAXDOP"):
		p.advance()
		hint := &nodes.QueryHint{Kind: "MAXDOP", Loc: nodes.Loc{Start: loc}}
		if p.cur.Type == tokICONST {
			hint.Value, _ = p.parsePrimary()
		}
		hint.Loc.End = p.pos()
		return hint

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "MAXRECURSION"):
		p.advance()
		hint := &nodes.QueryHint{Kind: "MAXRECURSION", Loc: nodes.Loc{Start: loc}}
		if p.cur.Type == tokICONST {
			hint.Value, _ = p.parsePrimary()
		}
		hint.Loc.End = p.pos()
		return hint

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "FAST"):
		p.advance()
		hint := &nodes.QueryHint{Kind: "FAST", Loc: nodes.Loc{Start: loc}}
		if p.cur.Type == tokICONST {
			hint.Value, _ = p.parsePrimary()
		}
		hint.Loc.End = p.pos()
		return hint

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "QUERYTRACEON"):
		p.advance()
		hint := &nodes.QueryHint{Kind: "QUERYTRACEON", Loc: nodes.Loc{Start: loc}}
		if p.cur.Type == tokICONST {
			hint.Value, _ = p.parsePrimary()
		}
		hint.Loc.End = p.pos()
		return hint

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "EXPAND"):
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "VIEWS") {
			p.advance()
			return &nodes.QueryHint{Kind: "EXPAND VIEWS", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "EXPAND", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "KEEP"):
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "KEEP PLAN", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "KEEP", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "KEEPFIXED"):
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "KEEPFIXED PLAN", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "KEEPFIXED", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "ROBUST"):
		p.advance()
		if p.cur.Type == kwPLAN {
			p.advance()
			return &nodes.QueryHint{Kind: "ROBUST PLAN", Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "ROBUST", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "PARAMETERIZATION"):
		p.advance()
		hint := &nodes.QueryHint{Kind: "PARAMETERIZATION", Loc: nodes.Loc{Start: loc}}
		if p.isIdentLike() {
			hint.StrValue = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		hint.Loc.End = p.pos()
		return hint

	case p.cur.Type == kwUSE:
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "HINT") {
			p.advance()
			hint := &nodes.QueryHint{Kind: "USE HINT", Loc: nodes.Loc{Start: loc}}
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
				_, _ = p.expect(')')
				if len(hintNames) > 0 {
					hint.Params = &nodes.List{Items: hintNames}
				}
			}
			hint.Loc.End = p.pos()
			return hint
		} else if p.cur.Type == kwPLAN {
			p.advance()
			hint := &nodes.QueryHint{Kind: "USE PLAN", Loc: nodes.Loc{Start: loc}}
			if p.cur.Type == tokNSCONST || p.cur.Type == tokSCONST {
				hint.StrValue = p.cur.Str
				p.advance()
			}
			hint.Loc.End = p.pos()
			return hint
		}
		return &nodes.QueryHint{Kind: "USE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "DISABLE"):
		p.advance()
		if p.isIdentLike() {
			suffix := strings.ToUpper(p.cur.Str)
			p.advance()
			return &nodes.QueryHint{Kind: "DISABLE " + suffix, Loc: nodes.Loc{Start: loc, End: p.pos()}}
		}
		return &nodes.QueryHint{Kind: "DISABLE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	case p.cur.Type == kwTABLE:
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "HINT") {
			p.advance()
			hint := &nodes.QueryHint{Kind: "TABLE HINT", Loc: nodes.Loc{Start: loc}}
			if p.cur.Type == '(' {
				p.advance()
				// Parse exposed_object_name as a TableRef
				hint.TableName , _ = p.parseTableRef()
				// Parse comma-separated table hints reusing parseTableHint()
				var hints []nodes.Node
				for p.cur.Type == ',' {
					p.advance()
					th := p.parseTableHint()
					if th != nil {
						hints = append(hints, th)
					}
				}
				if len(hints) > 0 {
					hint.TableHints = &nodes.List{Items: hints}
				}
				_, _ = p.expect(')')
			}
			hint.Loc.End = p.pos()
			return hint
		}
		return &nodes.QueryHint{Kind: "TABLE", Loc: nodes.Loc{Start: loc, End: p.pos()}}

	default:
		// Unknown hint with name = value pattern (e.g., MAX_GRANT_PERCENT = 10)
		if p.isIdentLike() {
			name := strings.ToUpper(p.cur.Str)
			p.advance()
			hint := &nodes.QueryHint{Kind: name, Loc: nodes.Loc{Start: loc}}
			if p.cur.Type == '=' {
				p.advance()
				hint.Value, _ = p.parsePrimary()
			}
			hint.Loc.End = p.pos()
			return hint
		}
		return nil
	}
}

// parseOptimizeForParam parses a single OPTIMIZE FOR parameter.
//
//	@variable_name { UNKNOWN | = literal_constant }
func (p *Parser) parseOptimizeForParam() *nodes.OptimizeForParam {
	if p.cur.Type != tokVARIABLE {
		return nil
	}
	loc := p.pos()
	param := &nodes.OptimizeForParam{
		Variable: p.cur.Str,
		Loc:      nodes.Loc{Start: loc},
	}
	p.advance()
	if p.cur.Type == '=' {
		p.advance()
		param.Value, _ = p.parsePrimary()
	} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "UNKNOWN") {
		param.Unknown = true
		p.advance()
	}
	param.Loc.End = p.pos()
	return param
}
