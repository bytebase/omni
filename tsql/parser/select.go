// Package parser - select.go implements T-SQL SELECT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
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
		p.advance()
		if _, err := p.expect('('); err == nil {
			var hints []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				hint := p.parseExpr()
				if hint != nil {
					hints = append(hints, hint)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			stmt.OptionClause = &nodes.List{Items: hints}
		}
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
		return p.parsePivotUnpivot(result)
	}

	// Alias
	alias := p.parseOptionalAlias()
	if alias != "" {
		ref.Alias = alias
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
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/xml/for-xml-sql-server
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
			// PATH('element')
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type != ')' {
					p.parseExpr() // consume the argument
				}
				_, _ = p.expect(')')
			}
		}
	} else if p.cur.Type == kwJSON {
		fc.Mode = nodes.ForJSON
		p.advance()
		// AUTO or PATH
		if p.isIdentLike() || p.cur.Type == kwPATH {
			fc.SubMode = strings.ToUpper(p.cur.Str)
			p.advance()
		}
	}

	fc.Loc.End = p.pos()
	return fc
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
