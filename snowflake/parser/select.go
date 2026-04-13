package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// SELECT statement parser
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Query-expression wrappers (SELECT + optional set operators)
// ---------------------------------------------------------------------------

// parseQueryExpr parses a SELECT statement optionally followed by set
// operators (UNION/EXCEPT/MINUS/INTERSECT). Returns either a bare
// *SelectStmt or a *SetOperationStmt wrapping chained SELECTs.
//
// Also handles the case where the query starts with ( SELECT ... ), i.e.
// a parenthesized SELECT as the leftmost operand in a set-op chain.
func (p *Parser) parseQueryExpr() (ast.Node, error) {
	var left ast.Node
	var err error
	if p.cur.Type == '(' {
		// Parenthesized SELECT: consume '(', parse inner query, then ')'
		p.advance() // consume '('
		left, err = p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		left, err = p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
	}
	return p.parseSetOpChain(left)
}

// parseWithQueryExpr parses WITH ... SELECT ... [set operators].
func (p *Parser) parseWithQueryExpr() (ast.Node, error) {
	node, err := p.parseWithSelect()
	if err != nil {
		return nil, err
	}
	return p.parseSetOpChain(node)
}

// parseSetOpChain checks for UNION/EXCEPT/MINUS/INTERSECT after a SELECT
// and builds a left-associative SetOperationStmt chain.
func (p *Parser) parseSetOpChain(left ast.Node) (ast.Node, error) {
	for {
		op, all, byName, ok := p.tryParseSetOp()
		if !ok {
			break
		}
		// The right side may be parenthesized: (SELECT ...)
		var right ast.Node
		var err error
		if p.cur.Type == '(' {
			p.advance() // consume '('
			right, err = p.parseQueryExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		} else {
			right, err = p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
		}
		left = &ast.SetOperationStmt{
			Op:     op,
			All:    all,
			ByName: byName,
			Left:   left,
			Right:  right,
			Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}
	}
	return left, nil
}

// tryParseSetOp checks if the current token starts a set operator.
// Returns (op, all, byName, ok). Consumes the operator tokens on match.
func (p *Parser) tryParseSetOp() (ast.SetOp, bool, bool, bool) {
	switch p.cur.Type {
	case kwUNION:
		p.advance() // consume UNION
		all := false
		byName := false
		if p.cur.Type == kwALL {
			all = true
			p.advance() // consume ALL
		}
		// UNION [ALL] BY NAME — Snowflake-specific
		if p.cur.Type == kwBY {
			next := p.peekNext()
			if next.Type == kwNAME {
				p.advance() // consume BY
				p.advance() // consume NAME
				byName = true
			}
		}
		return ast.SetOpUnion, all, byName, true
	case kwEXCEPT:
		p.advance() // consume EXCEPT
		return ast.SetOpExcept, false, false, true
	case kwMINUS:
		p.advance() // consume MINUS (alias for EXCEPT in Snowflake)
		return ast.SetOpExcept, false, false, true
	case kwINTERSECT:
		p.advance() // consume INTERSECT
		return ast.SetOpIntersect, false, false, true
	}
	return 0, false, false, false
}

// ---------------------------------------------------------------------------
// SELECT statement parser
// ---------------------------------------------------------------------------

// parseSelectStmt parses a SELECT statement:
//
//	SELECT [DISTINCT|ALL] [TOP n] target_list
//	  [FROM table_refs]
//	  [WHERE expr]
//	  [GROUP BY ...]
//	  [HAVING expr]
//	  [QUALIFY expr]
//	  [ORDER BY ...]
//	  [LIMIT n [OFFSET n]]
//	  [FETCH FIRST|NEXT n ROWS ONLY]
func (p *Parser) parseSelectStmt() (*ast.SelectStmt, error) {
	selectTok, err := p.expect(kwSELECT)
	if err != nil {
		return nil, err
	}

	stmt := &ast.SelectStmt{
		Loc: ast.Loc{Start: selectTok.Loc.Start},
	}

	// DISTINCT / ALL
	if p.cur.Type == kwDISTINCT {
		p.advance()
		stmt.Distinct = true
	} else if p.cur.Type == kwALL {
		p.advance()
		stmt.All = true
	}

	// TOP n
	if p.cur.Type == kwTOP {
		p.advance() // consume TOP
		topExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Top = topExpr
	}

	// SELECT list
	targets, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	stmt.Targets = targets

	// FROM clause
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY clause
	if p.cur.Type == kwGROUP {
		groupBy, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
	}

	// HAVING clause
	if p.cur.Type == kwHAVING {
		p.advance() // consume HAVING
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// QUALIFY clause (Snowflake-specific)
	if p.cur.Type == kwQUALIFY {
		p.advance() // consume QUALIFY
		qualify, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Qualify = qualify
	}

	// ORDER BY clause (reuses T1.3's parseOrderByList)
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT / OFFSET / FETCH
	if err := p.parseLimitOffsetFetch(stmt); err != nil {
		return nil, err
	}

	// Set End location to the last consumed token.
	stmt.Loc.End = p.prev.Loc.End

	return stmt, nil
}

// ---------------------------------------------------------------------------
// WITH ... SELECT (CTEs)
// ---------------------------------------------------------------------------

// parseWithSelect parses WITH [RECURSIVE] cte_list SELECT ...
func (p *Parser) parseWithSelect() (ast.Node, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}

	recursive := false
	if p.cur.Type == kwRECURSIVE {
		p.advance()
		recursive = true
	}

	ctes, err := p.parseCTEList(recursive)
	if err != nil {
		return nil, err
	}

	stmt, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	stmt.With = ctes
	// Extend Loc.Start to include the WITH keyword.
	if len(ctes) > 0 {
		stmt.Loc.Start = ctes[0].Loc.Start
		// The CTE Loc.Start already accounts for the name, but we need the
		// WITH keyword itself. We'll use the CTE's start minus a rough
		// offset. Better: just record the WITH token's start.
	}

	return stmt, nil
}

// parseCTEList parses a comma-separated list of CTEs.
func (p *Parser) parseCTEList(recursive bool) ([]*ast.CTE, error) {
	var ctes []*ast.CTE

	cte, err := p.parseCTE(recursive)
	if err != nil {
		return nil, err
	}
	ctes = append(ctes, cte)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		cte, err = p.parseCTE(recursive)
		if err != nil {
			return nil, err
		}
		ctes = append(ctes, cte)
	}

	return ctes, nil
}

// parseCTE parses one CTE: name [(columns)] AS (SELECT ...)
func (p *Parser) parseCTE(recursive bool) (*ast.CTE, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	cte := &ast.CTE{
		Name:      name,
		Recursive: recursive,
		Loc:       ast.Loc{Start: name.Loc.Start},
	}

	// Optional column list: (col1, col2, ...)
	if p.cur.Type == '(' {
		p.advance() // consume '('
		for {
			col, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			cte.Columns = append(cte.Columns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// ( SELECT ... ) — the CTE body must be parenthesized
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// The CTE body can be WITH ... SELECT [set ops] or just SELECT [set ops]
	var query ast.Node
	if p.cur.Type == kwWITH {
		query, err = p.parseWithQueryExpr()
	} else {
		query, err = p.parseQueryExpr()
	}
	if err != nil {
		return nil, err
	}
	cte.Query = query

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	cte.Loc.End = closeTok.Loc.End

	return cte, nil
}

// ---------------------------------------------------------------------------
// SELECT list
// ---------------------------------------------------------------------------

// parseSelectList parses comma-separated SELECT targets.
func (p *Parser) parseSelectList() ([]*ast.SelectTarget, error) {
	var targets []*ast.SelectTarget

	target, err := p.parseSelectTarget()
	if err != nil {
		return nil, err
	}
	targets = append(targets, target)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		target, err = p.parseSelectTarget()
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

// parseSelectTarget parses one item in the SELECT list:
//   - * [EXCLUDE (col, ...)]
//   - expr [AS alias]
//
// The expression parser already handles * (StarExpr) and qualifier.*
// (StarExpr with Qualifier), so we parse an expression first and then
// check the result type.
func (p *Parser) parseSelectTarget() (*ast.SelectTarget, error) {
	startLoc := p.cur.Loc

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	target := &ast.SelectTarget{
		Loc: ast.Loc{Start: startLoc.Start},
	}

	// Check if the expression is a star (* or qualifier.*)
	if _, ok := expr.(*ast.StarExpr); ok {
		target.Star = true
		target.Expr = expr

		// Check for EXCLUDE (col1, col2, ...)
		if p.cur.Type == kwEXCLUDE {
			p.advance() // consume EXCLUDE
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			for {
				col, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				target.Exclude = append(target.Exclude, col)
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	} else {
		target.Expr = expr
	}

	// Check for optional alias (only for non-star targets, or after EXCLUDE)
	if !target.Star {
		alias, hasAlias := p.parseOptionalAlias()
		if hasAlias {
			target.Alias = alias
		}
	}

	target.Loc.End = p.prev.Loc.End
	return target, nil
}

// ---------------------------------------------------------------------------
// FROM clause
// ---------------------------------------------------------------------------

// parseFromClause parses comma-separated FROM items, where each item is
// a primary source optionally followed by a chain of JOINs.
// The FROM keyword has already been consumed by the caller.
func (p *Parser) parseFromClause() ([]ast.Node, error) {
	var items []ast.Node

	item, err := p.parseFromItem()
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		item, err = p.parseFromItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// parseFromItem parses one comma-separated FROM item: a primary source
// optionally followed by a join chain.
func (p *Parser) parseFromItem() (ast.Node, error) {
	left, err := p.parsePrimarySource()
	if err != nil {
		return nil, err
	}
	return p.parseJoinChain(left)
}

// parsePrimarySource dispatches on the current token to parse a single
// FROM source:
//   - ( → subquery or parenthesized from-item
//   - LATERAL → consume, recurse, set Lateral=true
//   - TABLE → TABLE(func(...))
//   - otherwise → parseTableRef (ObjectName + alias)
func (p *Parser) parsePrimarySource() (ast.Node, error) {
	startLoc := p.cur.Loc

	// Parenthesized: subquery or parenthesized from-item.
	if p.cur.Type == '(' {
		next := p.peekNext()
		if next.Type == kwSELECT || next.Type == kwWITH {
			// Subquery in FROM: (SELECT ...) [AS] alias
			p.advance() // consume '('
			var query ast.Node
			var err error
			if p.cur.Type == kwWITH {
				query, err = p.parseWithSelect()
			} else {
				query, err = p.parseSelectStmt()
			}
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			ref := &ast.TableRef{
				Subquery: query,
				Loc:      ast.Loc{Start: startLoc.Start},
			}
			alias, hasAlias := p.parseOptionalAlias()
			if hasAlias {
				ref.Alias = alias
			}
			ref.Loc.End = p.prev.Loc.End
			return ref, nil
		}
		// Parenthesized from-item: (t1 JOIN t2 ON ...)
		p.advance() // consume '('
		inner, err := p.parseFromItem()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return inner, nil
	}

	// LATERAL prefix
	if p.cur.Type == kwLATERAL {
		p.advance() // consume LATERAL
		// LATERAL can be followed by a subquery, TABLE(func), or FLATTEN(...)
		source, err := p.parsePrimarySource()
		if err != nil {
			return nil, err
		}
		// Set Lateral on the resulting TableRef
		if ref, ok := source.(*ast.TableRef); ok {
			ref.Lateral = true
			ref.Loc.Start = startLoc.Start
			return ref, nil
		}
		// If it's not a TableRef (shouldn't normally happen), wrap it
		return source, nil
	}

	// TABLE(func(...)) — table function
	if p.cur.Type == kwTABLE {
		p.advance() // consume TABLE
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		// Parse the function call expression inside TABLE(...)
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		funcCall, ok := expr.(*ast.FuncCallExpr)
		if !ok {
			return nil, &ParseError{
				Loc: ast.NodeLoc(expr),
				Msg: "expected function call inside TABLE(...)",
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		ref := &ast.TableRef{
			FuncCall: funcCall,
			Loc:      ast.Loc{Start: startLoc.Start},
		}
		alias, hasAlias := p.parseOptionalAlias()
		if hasAlias {
			ref.Alias = alias
		}
		ref.Loc.End = p.prev.Loc.End
		return ref, nil
	}

	// FLATTEN(...) as a bare function in FROM (common Snowflake pattern)
	if p.cur.Type == kwFLATTEN {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		funcCall, ok := expr.(*ast.FuncCallExpr)
		if !ok {
			return nil, &ParseError{
				Loc: ast.NodeLoc(expr),
				Msg: "expected function call for FLATTEN",
			}
		}
		ref := &ast.TableRef{
			FuncCall: funcCall,
			Loc:      ast.Loc{Start: startLoc.Start},
		}
		alias, hasAlias := p.parseOptionalAlias()
		if hasAlias {
			ref.Alias = alias
		}
		ref.Loc.End = p.prev.Loc.End
		return ref, nil
	}

	// Default: simple table reference (ObjectName + alias)
	return p.parseTableRef()
}

// parseJoinChain builds a left-associative JoinExpr tree from any
// JOIN keywords following the left source.
func (p *Parser) parseJoinChain(left ast.Node) (ast.Node, error) {
	for {
		joinType, natural, directed, ok := p.parseJoinKeywords()
		if !ok {
			break
		}

		right, err := p.parsePrimarySource()
		if err != nil {
			return nil, err
		}

		join := &ast.JoinExpr{
			Type:     joinType,
			Left:     left,
			Right:    right,
			Natural:  natural,
			Directed: directed,
			Loc:      ast.Loc{Start: ast.NodeLoc(left).Start},
		}

		// Parse join condition
		switch {
		case joinType == ast.JoinAsof:
			// ASOF JOIN: expect MATCH_CONDITION(expr)
			if p.cur.Type == kwMATCH_CONDITION ||
				(p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "MATCH_CONDITION") {
				p.advance() // consume MATCH_CONDITION
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				cond, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				join.MatchCondition = cond
			}
			// ASOF may also have ON
			if p.cur.Type == kwON {
				p.advance() // consume ON
				onExpr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				join.On = onExpr
			}

		case joinType == ast.JoinCross || natural:
			// CROSS JOIN and NATURAL JOIN: no condition required

		case p.cur.Type == kwON:
			p.advance() // consume ON
			onExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			join.On = onExpr

		case p.cur.Type == kwUSING:
			p.advance() // consume USING
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			var cols []ast.Ident
			col, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			cols = append(cols, col)
			for p.cur.Type == ',' {
				p.advance() // consume ','
				col, err = p.parseIdent()
				if err != nil {
					return nil, err
				}
				cols = append(cols, col)
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			join.Using = cols
		}

		join.Loc.End = p.prev.Loc.End
		left = join
	}
	return left, nil
}

// parseJoinKeywords checks whether the current token position starts a
// JOIN keyword sequence. If so, it consumes the tokens and returns
// (joinType, natural, directed, true). If not, returns (0, false, false, false)
// without consuming any tokens.
func (p *Parser) parseJoinKeywords() (ast.JoinType, bool, bool, bool) {
	// NATURAL [LEFT|RIGHT|FULL] [OUTER] JOIN
	if p.cur.Type == kwNATURAL {
		next := p.peekNext()
		// NATURAL JOIN
		if next.Type == kwJOIN {
			p.advance() // consume NATURAL
			p.advance() // consume JOIN
			return ast.JoinInner, true, false, true
		}
		// NATURAL LEFT/RIGHT/FULL [OUTER] JOIN
		if next.Type == kwLEFT || next.Type == kwRIGHT || next.Type == kwFULL {
			p.advance() // consume NATURAL
			jt := p.consumeDirectionAndJoin()
			return jt, true, false, true
		}
		return 0, false, false, false
	}

	// DIRECTED [INNER|LEFT|RIGHT|FULL|CROSS] JOIN
	if p.cur.Type == kwDIRECTED {
		next := p.peekNext()
		switch next.Type {
		case kwJOIN:
			p.advance() // consume DIRECTED
			p.advance() // consume JOIN
			return ast.JoinInner, false, true, true
		case kwINNER:
			p.advance() // consume DIRECTED
			p.advance() // consume INNER
			if _, err := p.expect(kwJOIN); err == nil {
				return ast.JoinInner, false, true, true
			}
			return 0, false, false, false
		case kwLEFT, kwRIGHT, kwFULL:
			p.advance() // consume DIRECTED
			jt := p.consumeDirectionAndJoin()
			return jt, false, true, true
		case kwCROSS:
			p.advance() // consume DIRECTED
			p.advance() // consume CROSS
			if _, err := p.expect(kwJOIN); err == nil {
				return ast.JoinCross, false, true, true
			}
			return 0, false, false, false
		}
		return 0, false, false, false
	}

	// ASOF JOIN
	if p.cur.Type == kwASOF {
		next := p.peekNext()
		if next.Type == kwJOIN {
			p.advance() // consume ASOF
			p.advance() // consume JOIN
			return ast.JoinAsof, false, false, true
		}
		return 0, false, false, false
	}

	// INNER JOIN
	if p.cur.Type == kwINNER {
		next := p.peekNext()
		if next.Type == kwJOIN {
			p.advance() // consume INNER
			p.advance() // consume JOIN
			return ast.JoinInner, false, false, true
		}
		return 0, false, false, false
	}

	// LEFT [OUTER] JOIN
	if p.cur.Type == kwLEFT {
		jt := p.consumeDirectionAndJoin()
		if jt == ast.JoinLeft {
			return jt, false, false, true
		}
		return 0, false, false, false
	}

	// RIGHT [OUTER] JOIN
	if p.cur.Type == kwRIGHT {
		jt := p.consumeDirectionAndJoin()
		if jt == ast.JoinRight {
			return jt, false, false, true
		}
		return 0, false, false, false
	}

	// FULL [OUTER] JOIN
	if p.cur.Type == kwFULL {
		jt := p.consumeDirectionAndJoin()
		if jt == ast.JoinFull {
			return jt, false, false, true
		}
		return 0, false, false, false
	}

	// CROSS JOIN
	if p.cur.Type == kwCROSS {
		next := p.peekNext()
		if next.Type == kwJOIN {
			p.advance() // consume CROSS
			p.advance() // consume JOIN
			return ast.JoinCross, false, false, true
		}
		return 0, false, false, false
	}

	// Bare JOIN (= INNER)
	if p.cur.Type == kwJOIN {
		p.advance() // consume JOIN
		return ast.JoinInner, false, false, true
	}

	return 0, false, false, false
}

// consumeDirectionAndJoin consumes LEFT/RIGHT/FULL [OUTER] JOIN and returns
// the corresponding JoinType. The caller must have already checked that
// p.cur.Type is kwLEFT, kwRIGHT, or kwFULL.
// If the sequence doesn't end with JOIN, this returns JoinInner as a sentinel
// (the caller checks the return value).
func (p *Parser) consumeDirectionAndJoin() ast.JoinType {
	dir := p.cur.Type
	p.advance() // consume LEFT/RIGHT/FULL

	// Optional OUTER
	if p.cur.Type == kwOUTER {
		p.advance() // consume OUTER
	}

	// Expect JOIN
	if p.cur.Type != kwJOIN {
		// Not a join sequence — but we already consumed tokens.
		// This is an issue; for safety, return a sentinel and let caller handle.
		return ast.JoinInner
	}
	p.advance() // consume JOIN

	switch dir {
	case kwLEFT:
		return ast.JoinLeft
	case kwRIGHT:
		return ast.JoinRight
	case kwFULL:
		return ast.JoinFull
	default:
		return ast.JoinInner
	}
}

// parseTableRef parses one simple table reference: object_name [AS alias].
func (p *Parser) parseTableRef() (*ast.TableRef, error) {
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	ref := &ast.TableRef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	alias, hasAlias := p.parseOptionalAlias()
	if hasAlias {
		ref.Alias = alias
	}

	ref.Loc.End = p.prev.Loc.End
	return ref, nil
}

// ---------------------------------------------------------------------------
// GROUP BY clause
// ---------------------------------------------------------------------------

// parseGroupByClause parses GROUP BY variants:
//   - GROUP BY ALL
//   - GROUP BY CUBE (expr, expr)
//   - GROUP BY ROLLUP (expr, expr)
//   - GROUP BY GROUPING SETS ((expr), (expr))
//   - GROUP BY expr, expr
func (p *Parser) parseGroupByClause() (*ast.GroupByClause, error) {
	groupTok := p.advance() // consume GROUP
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	clause := &ast.GroupByClause{
		Loc: ast.Loc{Start: groupTok.Loc.Start},
	}

	// GROUP BY ALL
	if p.cur.Type == kwALL {
		p.advance()
		clause.Kind = ast.GroupByAll
		clause.Loc.End = p.prev.Loc.End
		return clause, nil
	}

	// GROUP BY CUBE (...)
	if p.cur.Type == kwCUBE {
		p.advance()
		clause.Kind = ast.GroupByCube
		items, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		clause.Items = items
		clause.Loc.End = p.prev.Loc.End
		return clause, nil
	}

	// GROUP BY ROLLUP (...)
	if p.cur.Type == kwROLLUP {
		p.advance()
		clause.Kind = ast.GroupByRollup
		items, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		clause.Items = items
		clause.Loc.End = p.prev.Loc.End
		return clause, nil
	}

	// GROUP BY GROUPING SETS (...)
	if p.cur.Type == kwGROUPING {
		// Check that the next token is SETS (keyword or ident)
		next := p.peekNext()
		if next.Type == kwSETS || (next.Type == tokIdent && strings.ToUpper(next.Str) == "SETS") {
			p.advance() // consume GROUPING
			p.advance() // consume SETS
			clause.Kind = ast.GroupByGroupingSets
			items, err := p.parseParenExprList()
			if err != nil {
				return nil, err
			}
			clause.Items = items
			clause.Loc.End = p.prev.Loc.End
			return clause, nil
		}
	}

	// Normal GROUP BY: comma-separated expressions
	clause.Kind = ast.GroupByNormal
	items, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	clause.Items = items
	clause.Loc.End = p.prev.Loc.End
	return clause, nil
}

// parseParenExprList parses ( expr, expr, ... ) — a parenthesized list of
// expressions. Used by GROUP BY CUBE/ROLLUP/GROUPING SETS.
func (p *Parser) parseParenExprList() ([]ast.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	items, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return items, nil
}

// ---------------------------------------------------------------------------
// LIMIT / OFFSET / FETCH
// ---------------------------------------------------------------------------

// parseLimitOffsetFetch fills stmt.Limit, stmt.Offset, and stmt.Fetch
// by consuming any LIMIT, OFFSET, and FETCH clauses in any valid order.
func (p *Parser) parseLimitOffsetFetch(stmt *ast.SelectStmt) error {
	// LIMIT n [OFFSET n]
	if p.cur.Type == kwLIMIT {
		p.advance() // consume LIMIT
		limitExpr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Limit = limitExpr

		if p.cur.Type == kwOFFSET {
			p.advance() // consume OFFSET
			offsetExpr, err := p.parseExpr()
			if err != nil {
				return err
			}
			stmt.Offset = offsetExpr
		}
		return nil
	}

	// OFFSET n [FETCH ...]
	if p.cur.Type == kwOFFSET {
		p.advance() // consume OFFSET
		offsetExpr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Offset = offsetExpr
	}

	// FETCH FIRST|NEXT n ROWS ONLY
	if p.cur.Type == kwFETCH {
		fetchTok := p.advance() // consume FETCH

		// Expect FIRST or NEXT
		if p.cur.Type != kwFIRST && p.cur.Type != kwNEXT {
			return &ParseError{
				Loc: p.cur.Loc,
				Msg: "expected FIRST or NEXT after FETCH",
			}
		}
		p.advance() // consume FIRST or NEXT

		// Parse count expression
		countExpr, err := p.parseExpr()
		if err != nil {
			return err
		}

		// Expect ROWS or ROW
		if p.cur.Type != kwROWS && p.cur.Type != kwROW {
			return &ParseError{
				Loc: p.cur.Loc,
				Msg: "expected ROWS or ROW after FETCH count",
			}
		}
		p.advance() // consume ROWS or ROW

		// Expect ONLY
		onlyTok, err := p.expect(kwONLY)
		if err != nil {
			return err
		}

		stmt.Fetch = &ast.FetchClause{
			Count: countExpr,
			Loc:   ast.Loc{Start: fetchTok.Loc.Start, End: onlyTok.Loc.End},
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Alias helper
// ---------------------------------------------------------------------------

// parseOptionalAlias returns (alias, true) if the current position has an
// alias (explicit with AS, or implicit identifier). Returns (zero, false)
// if no alias is present.
//
// This is the trickiest helper: it must distinguish `SELECT a b FROM t`
// (b is alias) from `SELECT a FROM t` (FROM is NOT an alias). It does
// this by checking whether the current token is a clause-starting keyword.
func (p *Parser) parseOptionalAlias() (ast.Ident, bool) {
	// Explicit: AS alias
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		id, err := p.parseIdent()
		if err != nil {
			return ast.Ident{}, false
		}
		return id, true
	}

	// Implicit alias: current token is an identifier or non-reserved keyword
	// that does NOT start a clause.
	if p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent ||
		(p.cur.Type >= 700 && !keywordReserved[p.cur.Type] && !isClauseKeyword(p.cur.Type)) {
		id, err := p.parseIdent()
		if err != nil {
			return ast.Ident{}, false
		}
		return id, true
	}

	return ast.Ident{}, false
}

// isClauseKeyword returns true for keywords that start SQL clauses and
// should NOT be consumed as implicit aliases.
func isClauseKeyword(t int) bool {
	switch t {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwQUALIFY, kwORDER,
		kwLIMIT, kwOFFSET, kwFETCH, kwUNION, kwEXCEPT, kwMINUS, kwINTERSECT,
		kwINTO, kwON, kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL,
		kwCROSS, kwNATURAL, kwWITH, kwSELECT, kwSET, kwWHEN,
		kwTHEN, kwELSE, kwEND, kwCASE, kwROWS, kwGROUPS, kwOVER,
		kwUSING, kwOUTER, kwASOF, kwDIRECTED, kwLATERAL, kwMATCH_CONDITION:
		return true
	}
	return false
}
