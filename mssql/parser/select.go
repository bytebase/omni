// Package parser - select.go implements T-SQL SELECT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSelectStmt parses a full SELECT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/select-transact-sql
//
//	[WITH cte_list]
//	SELECT [ALL | DISTINCT] [TOP ...]
//	    target_list
//	    [INTO table]
//	    [FROM table_source_list]
//	    [WHERE expr]
//	    [GROUP BY expr_list]
//	    [HAVING expr]
//	    [ORDER BY order_item_list]
//	    [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
//	    [FOR XML|JSON ...]
//	    [OPTION (...)]
//	    [UNION|INTERSECT|EXCEPT ...]
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
		stmt.IntoTable = p.parseTableRef()
	}

	// FROM
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause = p.parseFromClause()
	}

	// WHERE
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseExpr()
	}

	// GROUP BY
	if p.cur.Type == kwGROUP {
		p.advance()
		if _, err := p.expect(kwBY); err == nil {
			stmt.GroupByClause = p.parseGroupByList()
		}
	}

	// HAVING
	if _, ok := p.match(kwHAVING); ok {
		stmt.HavingClause = p.parseExpr()
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
		stmt.OffsetClause = p.parseExpr()
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
			count := p.parseExpr()
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

	// FOR XML / FOR JSON (only if followed by XML or JSON, not FOR UPDATE/READ_ONLY)
	if p.cur.Type == kwFOR {
		next := p.peekNext()
		if next.Type == kwXML || next.Type == kwJSON {
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

// parseWithClause parses WITH cte_name [(col_list)] AS (select), ...
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/with-common-table-expression-transact-sql
func (p *Parser) parseWithClause() *nodes.WithClause {
	loc := p.pos()
	p.advance() // consume WITH

	wc := &nodes.WithClause{
		Loc: nodes.Loc{Start: loc},
	}

	var ctes []nodes.Node
	for {
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
		tc.Count = p.parseExpr()
		_, _ = p.expect(')')
	} else {
		tc.Count = p.parsePrimary()
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
		expr := p.parseExpr()
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
				join.Condition = p.parseExpr()
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
	ref := p.parseTableRef()
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
	pivot.AggFunc = p.parseExpr()

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

	ts.Size = p.parseExpr()

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
			ts.Repeatable = p.parseExpr()
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
			arg := p.parseExpr()
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
	if p.cur.Type == tokIDENT {
		name := p.cur.Str
		p.advance()
		return name
	}
	return ""
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

// parseForClause parses FOR XML or FOR JSON.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/select-for-clause-transact-sql
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
		expr := p.parseExpr()
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
		expr := p.parseExpr()
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
			expr := p.parseExpr()
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
	expr := p.parseExpr()
	if expr != nil {
		return &nodes.List{Items: []nodes.Node{expr}}
	}
	return &nodes.List{Items: nil}
}

// parseOrderByList parses ORDER BY items.
func (p *Parser) parseOrderByList() *nodes.List {
	var items []nodes.Node
	for {
		oloc := p.pos()
		expr := p.parseExpr()
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
			hint.IntValue = p.parsePrimary()
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
// Returns a String node with the structured hint text.
func (p *Parser) parseQueryHint() nodes.Node {
	var buf strings.Builder

	// Handle keyword-based hints
	switch {
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "RECOMPILE"):
		buf.WriteString("RECOMPILE")
		p.advance()

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "OPTIMIZE"):
		buf.WriteString("OPTIMIZE")
		p.advance()
		if p.cur.Type == kwFOR {
			buf.WriteString(" FOR")
			p.advance()
			// OPTIMIZE FOR UNKNOWN
			if p.isIdentLike() && strings.EqualFold(p.cur.Str, "UNKNOWN") {
				buf.WriteString(" UNKNOWN")
				p.advance()
			} else if p.cur.Type == '(' {
				// OPTIMIZE FOR ( @var = val | UNKNOWN [,...n] )
				buf.WriteString(" (")
				p.advance()
				first := true
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if !first {
						buf.WriteString(", ")
					}
					first = false
					if p.cur.Type == tokVARIABLE {
						buf.WriteString(p.cur.Str)
						p.advance()
						if p.cur.Type == '=' {
							buf.WriteString(" = ")
							p.advance()
							// literal value
							if p.cur.Str != "" {
								buf.WriteString(p.cur.Str)
							}
							p.advance()
						} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "UNKNOWN") {
							buf.WriteString(" UNKNOWN")
							p.advance()
						}
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				p.expect(')')
				buf.WriteString(")")
			}
		}

	case p.isIdentLike() && (strings.EqualFold(p.cur.Str, "LOOP") ||
		strings.EqualFold(p.cur.Str, "HASH") ||
		strings.EqualFold(p.cur.Str, "MERGE")):
		buf.WriteString(strings.ToUpper(p.cur.Str))
		p.advance()
		// JOIN | UNION | GROUP
		if p.cur.Type == kwJOIN {
			buf.WriteString(" JOIN")
			p.advance()
		} else if p.cur.Type == kwUNION {
			buf.WriteString(" UNION")
			p.advance()
		} else if p.cur.Type == kwGROUP {
			buf.WriteString(" GROUP")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "CONCAT"):
		buf.WriteString("CONCAT")
		p.advance()
		if p.cur.Type == kwUNION {
			buf.WriteString(" UNION")
			p.advance()
		}

	case p.cur.Type == kwORDER:
		buf.WriteString("ORDER")
		p.advance()
		if p.cur.Type == kwGROUP {
			buf.WriteString(" GROUP")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "FORCE"):
		buf.WriteString("FORCE")
		p.advance()
		if p.cur.Type == kwORDER {
			buf.WriteString(" ORDER")
			p.advance()
		} else if p.isIdentLike() {
			buf.WriteString(" ")
			buf.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "MAXDOP"):
		buf.WriteString("MAXDOP")
		p.advance()
		if p.cur.Type == tokICONST {
			buf.WriteString(" ")
			buf.WriteString(p.cur.Str)
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "MAXRECURSION"):
		buf.WriteString("MAXRECURSION")
		p.advance()
		if p.cur.Type == tokICONST {
			buf.WriteString(" ")
			buf.WriteString(p.cur.Str)
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "FAST"):
		buf.WriteString("FAST")
		p.advance()
		if p.cur.Type == tokICONST {
			buf.WriteString(" ")
			buf.WriteString(p.cur.Str)
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "EXPAND"):
		buf.WriteString("EXPAND")
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "VIEWS") {
			buf.WriteString(" VIEWS")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "KEEP"):
		buf.WriteString("KEEP")
		p.advance()
		if p.cur.Type == kwPLAN {
			buf.WriteString(" PLAN")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "KEEPFIXED"):
		buf.WriteString("KEEPFIXED")
		p.advance()
		if p.cur.Type == kwPLAN {
			buf.WriteString(" PLAN")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "ROBUST"):
		buf.WriteString("ROBUST")
		p.advance()
		if p.cur.Type == kwPLAN {
			buf.WriteString(" PLAN")
			p.advance()
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "PARAMETERIZATION"):
		buf.WriteString("PARAMETERIZATION")
		p.advance()
		if p.isIdentLike() {
			buf.WriteString(" ")
			buf.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
		}

	case p.cur.Type == kwUSE:
		buf.WriteString("USE")
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "HINT") {
			buf.WriteString(" HINT")
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				first := true
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if !first {
						buf.WriteString(", ")
					}
					first = false
					if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
						buf.WriteString("'")
						buf.WriteString(p.cur.Str)
						buf.WriteString("'")
						p.advance()
					}
					p.match(',')
				}
				p.expect(')')
			}
		} else if p.cur.Type == kwPLAN {
			buf.WriteString(" PLAN")
			p.advance()
			// N'xml_plan'
			if p.cur.Type == tokNSCONST || p.cur.Type == tokSCONST {
				buf.WriteString(" ")
				buf.WriteString(p.cur.Str)
				p.advance()
			}
		}

	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "DISABLE"):
		buf.WriteString("DISABLE")
		p.advance()
		if p.isIdentLike() {
			buf.WriteString(" ")
			buf.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
		}

	case p.cur.Type == kwTABLE:
		buf.WriteString("TABLE")
		p.advance()
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "HINT") {
			buf.WriteString(" HINT")
			p.advance()
			// Parse TABLE HINT(exposed_object_name, hint [,...n])
			//
			// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-query
			//
			//   TABLE HINT ( exposed_object_name [ , <table_hint> [ , ...n ] ] )
			//   <table_hint> ::= NOLOCK | READUNCOMMITTED | UPDLOCK | REPEATABLEREAD
			//                  | SERIALIZABLE | READCOMMITTED | TABLOCK | TABLOCKX
			//                  | PAGLOCK | ROWLOCK | NOWAIT | READPAST | XLOCK
			//                  | SNAPSHOT | NOEXPAND | HOLDLOCK | KEEPIDENTITY
			//                  | KEEPDEFAULTS | IGNORE_CONSTRAINTS | IGNORE_TRIGGERS
			//                  | INDEX ( index_val [, ...n] )
			//                  | FORCESEEK [ ( index_val ( col [, ...n] ) ) ]
			//                  | FORCESCAN | SPATIAL_WINDOW_MAX_CELLS = n
			if p.cur.Type == '(' {
				p.advance() // consume '('
				buf.WriteByte('(')
				// exposed_object_name (may be multi-part)
				first := true
				for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
					if !first {
						buf.WriteByte('.')
					}
					buf.WriteString(p.cur.Str)
					p.advance()
					if p.cur.Type == '.' {
						p.advance()
					}
					first = false
				}
				// Parse comma-separated hints
				for p.cur.Type == ',' {
					buf.WriteString(", ")
					p.advance() // consume ','
					if p.isIdentLike() || p.cur.Type >= kwADD {
						hintName := strings.ToUpper(p.cur.Str)
						buf.WriteString(hintName)
						p.advance()
						// INDEX(...), FORCESEEK(...) may have parenthesized args
						if p.cur.Type == '(' {
							depth := 1
							p.advance()
							buf.WriteByte('(')
							for depth > 0 && p.cur.Type != tokEOF {
								if p.cur.Type == '(' {
									depth++
									buf.WriteByte('(')
								} else if p.cur.Type == ')' {
									depth--
									if depth > 0 {
										buf.WriteByte(')')
									}
								} else {
									buf.WriteString(p.cur.Str)
								}
								if depth > 0 {
									p.advance()
									if p.cur.Type == ',' && depth > 0 {
										buf.WriteString(", ")
									}
								}
							}
							buf.WriteByte(')')
							p.advance() // consume final ')'
						}
						// SPATIAL_WINDOW_MAX_CELLS = n
						if p.cur.Type == '=' {
							p.advance()
							buf.WriteString(" = ")
							if p.cur.Str != "" {
								buf.WriteString(p.cur.Str)
							}
							p.advance()
						}
					}
				}
				buf.WriteByte(')')
				p.match(')') // consume final ')'
			}
		}

	default:
		// Unknown or keyword-based hint: try to parse as expression
		if p.isIdentLike() {
			buf.WriteString(strings.ToUpper(p.cur.Str))
			p.advance()
			// Handle optional value after hint name (e.g., MAX_GRANT_PERCENT = 10)
			if p.cur.Type == '=' {
				p.advance()
				buf.WriteString(" = ")
				if p.cur.Str != "" {
					buf.WriteString(p.cur.Str)
				}
				p.advance()
			}
		} else {
			// Fallback: parse as expression
			expr := p.parseExpr()
			if expr != nil {
				return expr
			}
			return nil
		}
	}

	if buf.Len() == 0 {
		return nil
	}
	return &nodes.String{Str: buf.String()}
}
