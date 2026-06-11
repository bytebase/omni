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
	switch {
	case p.cur.Type == '(':
		// Parenthesized query body: consume '(', parse inner query, then ')'.
		// The inner query may itself begin with WITH (a CTE / RECURSIVE), e.g.
		// CREATE VIEW v AS ( WITH cte AS (...) SELECT ... ), so dispatch on WITH
		// rather than always recursing into the bare-SELECT path.
		p.advance() // consume '('
		left, err = p.parseQueryBody()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(')'); err != nil {
			return nil, err
		}
	case p.cur.Type == kwWITH:
		left, err = p.parseWithQueryExpr()
		if err != nil {
			return nil, err
		}
	default:
		left, err = p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
	}
	return p.parseSetOpChain(left)
}

// parseQueryBody parses a query expression that may begin with WITH (a CTE /
// RECURSIVE block) or be a plain SELECT / set-operation / parenthesized query.
// It is the shared dispatch used wherever a query can appear after a '(' — the
// parenthesized operand of a set-op chain, a CTE body, and a parenthesized view
// body — so a leading WITH is never mistaken for a bare SELECT.
func (p *Parser) parseQueryBody() (ast.Node, error) {
	if p.cur.Type == kwWITH {
		return p.parseWithQueryExpr()
	}
	return p.parseQueryExpr()
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
//	  [FETCH [FIRST|NEXT] n [ROW|ROWS] [ONLY]]
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
		topExpr, err := p.parseTopCount()
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

	// START WITH ... CONNECT BY ... (hierarchical query). Snowflake places
	// these after WHERE; START WITH is optional and, when present, precedes
	// CONNECT BY.
	if err := p.parseConnectBy(stmt); err != nil {
		return nil, err
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

// parseTopCount parses the count following TOP. Snowflake documents TOP n
// as taking a constant, so this deliberately parses a single primary
// expression (number literal, $variable, or parenthesized expression)
// rather than a full expression: in `SELECT TOP 125 * FROM t` the `*` is
// the star select target, not a multiplication operator applied to 125.
func (p *Parser) parseTopCount() (ast.Node, error) {
	switch p.cur.Type {
	case tokInt, tokFloat, tokReal, tokVariable, '(':
		return p.parsePrimaryExpr()
	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected a number after TOP",
		}
	}
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

// parseValuesClause parses a VALUES query expression used as a row source:
//
//	VALUES (expr, …) [, (expr, …) …]
//
// It is reached from parsePrimarySource when a parenthesized derived table
// opens with VALUES (i.e. `( VALUES … )`). Row parsing is delegated to the
// shared parseValuesRows (the same code path INSERT uses); its loop makes
// progress only on a `, (` and so cannot spin. Rows need not be uniform in
// arity — Snowflake accepts ragged VALUES lists.
func (p *Parser) parseValuesClause() (ast.Node, error) {
	start := p.cur.Loc.Start
	rows, err := p.parseValuesRows()
	if err != nil {
		return nil, err
	}
	return &ast.ValuesClause{
		Rows: rows,
		Loc:  ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// SELECT list
// ---------------------------------------------------------------------------

// parseSelectList parses comma-separated SELECT targets.
//
// Snowflake permits a trailing comma right before the clause that follows the
// SELECT list (e.g. `SELECT a, b, FROM t`). After consuming a comma, if the
// next token terminates the list (FROM or another query clause, a `)`, `;`, or
// EOF) we stop without parsing — and erroring on — an empty target.
func (p *Parser) parseSelectList() ([]*ast.SelectTarget, error) {
	var targets []*ast.SelectTarget

	target, err := p.parseSelectTarget()
	if err != nil {
		return nil, err
	}
	targets = append(targets, target)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		if p.selectListTerminator() {
			// Trailing comma before a clause terminator — allowed.
			break
		}
		target, err = p.parseSelectTarget()
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

// selectListTerminator reports whether the current token ends the SELECT list.
// Used to accept a trailing comma: a comma immediately followed by one of these
// tokens is a permitted trailing comma rather than an empty list item. It is
// deliberately narrow — only genuine list-ending tokens — so a stray comma in
// the middle of the list (e.g. `SELECT a, , b`) still errors.
func (p *Parser) selectListTerminator() bool {
	switch p.cur.Type {
	case tokEOF, ')', ';':
		return true
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwQUALIFY, kwORDER,
		kwLIMIT, kwOFFSET, kwFETCH, kwUNION, kwEXCEPT, kwMINUS, kwINTERSECT,
		kwINTO, kwSTART, kwCONNECT:
		return true
	}
	return false
}

// parseSelectTarget parses one item in the SELECT list:
//   - * [ILIKE 'pattern'] [EXCLUDE (col, ...)] [REPLACE (expr AS col, ...)] [RENAME (col AS alias, ...)]
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

	// Star ILIKE transform: ILIKE is an infix operator, so for
	// `* ILIKE '<pattern>'` (and `tbl.* ILIKE ...`) the expression parser has
	// already bound the star into LikeExpr(StarExpr, pattern) before this
	// function can see it. A star is not a scalar operand, so in a SELECT
	// list that shape is unambiguously Snowflake's star ILIKE
	// column-transform, not a boolean expression — unwrap it here, at the
	// select-target boundary, so other expression contexts are untouched.
	// The unwrap keys on the exact documented shape (plain ILIKE with a
	// string-literal pattern); NOT/ANY/ESCAPE variants and non-literal
	// patterns are not transforms and stay expressions.
	if like, ok := expr.(*ast.LikeExpr); ok &&
		like.Op == ast.LikeOpILike && !like.Not && !like.Any && like.Escape == nil {
		if _, isStar := like.Expr.(*ast.StarExpr); isStar {
			if pat, isLit := like.Pattern.(*ast.Literal); isLit && pat.Kind == ast.LitString {
				expr = like.Expr
				target.Ilike = pat
			}
		}
	}

	// Check if the expression is a star (* or qualifier.*)
	if _, ok := expr.(*ast.StarExpr); ok {
		target.Star = true
		target.Expr = expr

		// Star column-transforms, in Snowflake's documented order: ILIKE
		// (unwrapped above), then EXCLUDE, then REPLACE, then RENAME (any
		// subset may be present; the docs forbid combining ILIKE with
		// EXCLUDE, which the parser over-accepts — the combination is
		// positionally in order and parses soundly).
		//   ILIKE '<pattern>'
		//   EXCLUDE <col> | EXCLUDE (<col>, ...)
		//   REPLACE (<expr> AS <col>, ...)
		//   RENAME <col> AS <alias> | RENAME (<col> AS <alias>, ...)
		if p.cur.Type == kwEXCLUDE {
			p.advance() // consume EXCLUDE
			if err := p.parseStarExclude(target); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == kwREPLACE {
			p.advance() // consume REPLACE
			if err := p.parseStarReplace(target); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == kwRENAME {
			p.advance() // consume RENAME
			if err := p.parseStarRename(target); err != nil {
				return nil, err
			}
		}
		// A transform keyword still pending here is out of documented order
		// (e.g. `* RENAME (...) REPLACE (...)`, or ILIKE after any other
		// transform — ILIKE must come first, and in first position it was
		// already consumed by the expression parse and unwrapped above).
		// parseSingle ignores tokens after a completed statement, so without
		// this check the rest of the statement — including FROM — would be
		// dropped silently. Fail loudly instead.
		switch p.cur.Type {
		case kwILIKE, kwEXCLUDE, kwREPLACE, kwRENAME:
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "star column-transforms must appear in ILIKE/EXCLUDE, REPLACE, RENAME order",
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

// parseStarExclude parses the EXCLUDE column-transform that may follow a star
// in a SELECT list. The EXCLUDE keyword has already been consumed. Accepts
// either a single bare column or a parenthesized comma-separated list:
//
//	EXCLUDE <col>
//	EXCLUDE (<col>, <col>, ...)
func (p *Parser) parseStarExclude(target *ast.SelectTarget) error {
	if p.cur.Type == '(' {
		p.advance() // consume '('
		for {
			// Progress guard: each iteration must consume at least one token
			// (parseIdent advances or errors; the comma is consumed below), so
			// a stalled cursor means a malformed list rather than an infinite
			// loop.
			before := p.cur.Loc.Start
			col, err := p.parseIdent()
			if err != nil {
				return err
			}
			target.Exclude = append(target.Exclude, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
			if p.cur.Loc.Start <= before {
				return &ParseError{Loc: p.cur.Loc, Msg: "malformed EXCLUDE list"}
			}
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
		return nil
	}
	// Single bare column.
	col, err := p.parseIdent()
	if err != nil {
		return err
	}
	target.Exclude = append(target.Exclude, col)
	return nil
}

// parseStarRename parses the RENAME column-transform that may follow a star in
// a SELECT list. The RENAME keyword has already been consumed. Accepts either a
// single bare `col AS alias` or a parenthesized comma-separated list of them:
//
//	RENAME <col> AS <alias>
//	RENAME (<col> AS <alias>, <col> AS <alias>, ...)
func (p *Parser) parseStarRename(target *ast.SelectTarget) error {
	if p.cur.Type == '(' {
		p.advance() // consume '('
		for {
			// Progress guard: see parseStarExclude — the cursor must strictly
			// advance each iteration.
			before := p.cur.Loc.Start
			pair, err := p.parseStarRenamePair()
			if err != nil {
				return err
			}
			target.Rename = append(target.Rename, pair)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
			if p.cur.Loc.Start <= before {
				return &ParseError{Loc: p.cur.Loc, Msg: "malformed RENAME list"}
			}
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
		return nil
	}
	// Single bare `col AS alias`.
	pair, err := p.parseStarRenamePair()
	if err != nil {
		return err
	}
	target.Rename = append(target.Rename, pair)
	return nil
}

// parseStarRenamePair parses one `<col> AS <alias>` pair of a RENAME transform.
// The AS keyword is required (Snowflake's documented syntax).
func (p *Parser) parseStarRenamePair() (ast.StarRename, error) {
	col, err := p.parseIdent()
	if err != nil {
		return ast.StarRename{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return ast.StarRename{}, err
	}
	alias, err := p.parseIdent()
	if err != nil {
		return ast.StarRename{}, err
	}
	return ast.StarRename{Col: col, Alias: alias}, nil
}

// parseStarReplace parses the REPLACE column-transform that may follow a star
// in a SELECT list. The REPLACE keyword has already been consumed. Unlike
// EXCLUDE and RENAME, Snowflake documents only the parenthesized form (the
// replacement is an arbitrary expression, so the parentheses are required):
//
//	REPLACE (<expr> AS <col> [, <expr> AS <col>, ...])
func (p *Parser) parseStarReplace(target *ast.SelectTarget) error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		// Progress guard: see parseStarExclude — the cursor must strictly
		// advance each iteration.
		before := p.cur.Loc.Start
		pair, err := p.parseStarReplacePair()
		if err != nil {
			return err
		}
		target.Replace = append(target.Replace, pair)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		if p.cur.Loc.Start <= before {
			return &ParseError{Loc: p.cur.Loc, Msg: "malformed REPLACE list"}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
}

// parseStarReplacePair parses one `<expr> AS <col>` pair of a REPLACE
// transform. The AS keyword is required (Snowflake's documented syntax); Col
// must name an existing column of the star expansion, which the parser cannot
// check — it records whatever identifier follows AS.
func (p *Parser) parseStarReplacePair() (ast.StarReplace, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return ast.StarReplace{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return ast.StarReplace{}, err
	}
	col, err := p.parseIdent()
	if err != nil {
		return ast.StarReplace{}, err
	}
	return ast.StarReplace{Expr: expr, Col: col}, nil
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

	// Parenthesized: subquery, VALUES row-source, or parenthesized from-item.
	if p.cur.Type == '(' {
		next := p.peekNext()
		if next.Type == kwSELECT || next.Type == kwWITH || next.Type == kwVALUES {
			// Derived table in FROM: ( SELECT … ) | ( VALUES … ) [AS] alias
			p.advance() // consume '('
			var query ast.Node
			var err error
			switch p.cur.Type {
			case kwWITH:
				query, err = p.parseWithSelect()
			case kwVALUES:
				query, err = p.parseValuesClause()
			default:
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
			if err := p.parseTableSuffix(ref); err != nil {
				return nil, err
			}
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
		if err := p.parseTableSuffix(ref); err != nil {
			return nil, err
		}
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

	// Bare (unparenthesized) VALUES row source: FROM VALUES (1, 'a'), (2, 'b').
	// Snowflake accepts the inline VALUES list without the surrounding
	// parentheses that the derived-table form above requires. VALUES is a
	// reserved keyword, so `VALUES (` in FROM position is unambiguous.
	if p.cur.Type == kwVALUES && p.peekNext().Type == '(' {
		query, err := p.parseValuesClause()
		if err != nil {
			return nil, err
		}
		ref := &ast.TableRef{
			Subquery: query,
			Loc:      ast.Loc{Start: startLoc.Start},
		}
		if err := p.parseTableSuffix(ref); err != nil {
			return nil, err
		}
		return ref, nil
	}

	// $N result-set table reference: FROM $1. A bare positional $-reference in
	// FROM position names the result set of the preceding statement (the source
	// of a `->>` result-pipe). A qualified $N (t.$1) is a column expression, not
	// a source, and never reaches FROM position, so only the bare token is
	// handled here.
	if p.cur.Type == tokVariable {
		dollar := p.parseDollarRef(nil)
		ref := &ast.TableRef{
			DollarN: dollar,
			Loc:     ast.Loc{Start: startLoc.Start, End: dollar.Loc.End},
		}
		if err := p.parseTableSuffix(ref); err != nil {
			return nil, err
		}
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
	// NATURAL [INNER | {LEFT|RIGHT|FULL} [OUTER]] JOIN
	if p.cur.Type == kwNATURAL {
		next := p.peekNext()
		// NATURAL JOIN
		if next.Type == kwJOIN {
			p.advance() // consume NATURAL
			p.advance() // consume JOIN
			return ast.JoinInner, true, false, true
		}
		// NATURAL INNER JOIN
		if next.Type == kwINNER {
			p.advance() // consume NATURAL
			p.advance() // consume INNER
			if p.cur.Type == kwJOIN {
				p.advance() // consume JOIN
				return ast.JoinInner, true, false, true
			}
			return 0, false, false, false
		}
		// NATURAL LEFT/RIGHT/FULL [OUTER] JOIN
		if next.Type == kwLEFT || next.Type == kwRIGHT || next.Type == kwFULL {
			p.advance() // consume NATURAL
			jt, _ := p.consumeDirectionAndJoin()
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
			jt, _ := p.consumeDirectionAndJoin()
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

	// INNER [DIRECTED] JOIN — Snowflake's documented placement of DIRECTED is
	// between the join-type keyword and JOIN (t1 INNER DIRECTED JOIN t2).
	if p.cur.Type == kwINNER {
		next := p.peekNext()
		if next.Type == kwJOIN {
			p.advance() // consume INNER
			p.advance() // consume JOIN
			return ast.JoinInner, false, false, true
		}
		if next.Type == kwDIRECTED {
			p.advance() // consume INNER
			p.advance() // consume DIRECTED
			if p.cur.Type == kwJOIN {
				p.advance() // consume JOIN
				return ast.JoinInner, false, true, true
			}
			return 0, false, false, false
		}
		return 0, false, false, false
	}

	// LEFT [OUTER] [DIRECTED] JOIN
	if p.cur.Type == kwLEFT {
		jt, directed := p.consumeDirectionAndJoin()
		if jt == ast.JoinLeft {
			return jt, false, directed, true
		}
		return 0, false, false, false
	}

	// RIGHT [OUTER] [DIRECTED] JOIN
	if p.cur.Type == kwRIGHT {
		jt, directed := p.consumeDirectionAndJoin()
		if jt == ast.JoinRight {
			return jt, false, directed, true
		}
		return 0, false, false, false
	}

	// FULL [OUTER] [DIRECTED] JOIN
	if p.cur.Type == kwFULL {
		jt, directed := p.consumeDirectionAndJoin()
		if jt == ast.JoinFull {
			return jt, false, directed, true
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

// consumeDirectionAndJoin consumes LEFT/RIGHT/FULL [OUTER] [DIRECTED] JOIN
// and returns the corresponding JoinType plus whether the DIRECTED keyword
// was present. The caller must have already checked that p.cur.Type is
// kwLEFT, kwRIGHT, or kwFULL.
// If the sequence doesn't end with JOIN, this returns JoinInner as a sentinel
// (the caller checks the return value).
func (p *Parser) consumeDirectionAndJoin() (ast.JoinType, bool) {
	dir := p.cur.Type
	p.advance() // consume LEFT/RIGHT/FULL

	// Optional OUTER
	if p.cur.Type == kwOUTER {
		p.advance() // consume OUTER
	}

	// Optional DIRECTED (Snowflake places it directly before JOIN).
	directed := false
	if p.cur.Type == kwDIRECTED {
		p.advance() // consume DIRECTED
		directed = true
	}

	// Expect JOIN
	if p.cur.Type != kwJOIN {
		// Not a join sequence — but we already consumed tokens.
		// This is an issue; for safety, return a sentinel and let caller handle.
		return ast.JoinInner, false
	}
	p.advance() // consume JOIN

	switch dir {
	case kwLEFT:
		return ast.JoinLeft, directed
	case kwRIGHT:
		return ast.JoinRight, directed
	case kwFULL:
		return ast.JoinFull, directed
	default:
		return ast.JoinInner, directed
	}
}

// parseTableRef parses one simple table reference: object_name followed by
// the Snowflake table-attached clause chain (AT/BEFORE, CHANGES,
// MATCH_RECOGNIZE, PIVOT/UNPIVOT, alias, SAMPLE).
func (p *Parser) parseTableRef() (*ast.TableRef, error) {
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	ref := &ast.TableRef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	if err := p.parseTableSuffix(ref); err != nil {
		return nil, err
	}
	return ref, nil
}

// parseTableSuffix parses the chain of Snowflake clauses that may follow a
// table primary, in the documented source order:
//
//	AT|BEFORE → CHANGES → MATCH_RECOGNIZE → PIVOT|UNPIVOT → [AS] alias → SAMPLE
//
// The aggregate alias position sits between PIVOT/UNPIVOT and SAMPLE (matching
// `table1 AS t1 SAMPLE (25)` from the docs). ref already has its source set
// (Name / Subquery / FuncCall); this fills the clause fields, Alias, and the
// final End location. Each branch is guarded by a keyword, so the function
// makes progress only when a clause is actually present and is otherwise a
// no-op (leaving ref.Loc.End at the source's end).
func (p *Parser) parseTableSuffix(ref *ast.TableRef) error {
	ref.Loc.End = p.prev.Loc.End

	// AT | BEFORE ( … )
	if p.cur.Type == kwAT || p.cur.Type == kwBEFORE {
		tt, err := p.parseTimeTravelClause()
		if err != nil {
			return err
		}
		ref.TimeTravel = tt
		ref.Loc.End = tt.Loc.End
	}

	// CHANGES ( … ) AT|BEFORE ( … ) [END ( … )]
	if p.cur.Type == kwCHANGES {
		ch, err := p.parseChangesClause()
		if err != nil {
			return err
		}
		ref.Changes = ch
		ref.Loc.End = ch.Loc.End
	}

	// MATCH_RECOGNIZE ( … )
	if p.cur.Type == kwMATCH_RECOGNIZE {
		mr, err := p.parseMatchRecognizeClause()
		if err != nil {
			return err
		}
		ref.MatchRecognize = mr
		ref.Loc.End = mr.Loc.End
		// MATCH_RECOGNIZE consumes its own trailing alias; only SAMPLE may
		// follow (per the legacy grammar's clause ordering).
		return p.parseTrailingSample(ref)
	}

	// PIVOT ( … ) | UNPIVOT ( … ) — possibly chained (`src PIVOT(...)
	// PIVOT(...)`, official docs): each subsequent clause applies to the
	// result of the previous one, so on the second and later clauses the ref
	// is re-rooted in place with the prior contents moved to Nested.
	pivoted := false
	for p.cur.Type == kwPIVOT || p.cur.Type == kwUNPIVOT {
		if ref.Pivot != nil || ref.Unpivot != nil {
			inner := &ast.TableRef{}
			*inner = *ref
			*ref = ast.TableRef{
				Nested: inner,
				Loc:    ast.Loc{Start: inner.Loc.Start, End: inner.Loc.End},
			}
		}
		switch p.cur.Type {
		case kwPIVOT:
			pv, err := p.parsePivotClause()
			if err != nil {
				return err
			}
			ref.Pivot = pv
			ref.Loc.End = pv.Loc.End
		case kwUNPIVOT:
			uv, err := p.parseUnpivotClause()
			if err != nil {
				return err
			}
			ref.Unpivot = uv
			ref.Loc.End = uv.Loc.End
		}
		pivoted = true
	}
	if pivoted {
		// PIVOT/UNPIVOT consume their own trailing alias; only SAMPLE may
		// follow.
		return p.parseTrailingSample(ref)
	}

	// [AS] alias
	if alias, has := p.parseOptionalAlias(); has {
		ref.Alias = alias
		ref.Loc.End = p.prev.Loc.End

		// Derived column list: alias (c1, c2, …). Only valid on a derived
		// table (subquery / VALUES); a parenthesis after a plain table name is
		// not a column list and is left for the caller. Requires an explicit
		// alias to precede it, which the IsEmpty guard above ensures.
		if ref.Subquery != nil && p.cur.Type == '(' {
			cols, err := p.parseDerivedColumnList()
			if err != nil {
				return err
			}
			ref.Columns = cols
			ref.Loc.End = p.prev.Loc.End
		}
	}

	// SAMPLE | TABLESAMPLE ( … )
	return p.parseTrailingSample(ref)
}

// parseDerivedColumnList parses a parenthesized derived-table column list:
//
//	( col1 [, col2 …] )
//
// The caller must have already verified p.cur is '('. The list must be
// non-empty (an empty "()" is rejected). Each iteration consumes one
// identifier and either a ',' (continue) or ')' (stop), so the loop always
// makes progress.
func (p *Parser) parseDerivedColumnList() ([]ast.Ident, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var cols []ast.Ident
	for {
		col, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseTrailingSample parses an optional SAMPLE / TABLESAMPLE clause and
// updates ref. It is a no-op when no sampling clause is present.
func (p *Parser) parseTrailingSample(ref *ast.TableRef) error {
	if p.cur.Type == kwSAMPLE || p.cur.Type == kwTABLESAMPLE {
		s, err := p.parseSampleClause()
		if err != nil {
			return err
		}
		ref.Sample = s
		ref.Loc.End = s.Loc.End
	}
	return nil
}

// ---------------------------------------------------------------------------
// START WITH / CONNECT BY (hierarchical queries)
// ---------------------------------------------------------------------------

// parseConnectBy parses the optional hierarchical-query clauses:
//
//	[ START WITH <condition> ] CONNECT BY [PRIOR] <condition> [ , … ]
//
// Snowflake documents START WITH before CONNECT BY. PRIOR appears as a
// prefix operator inside the condition expressions (handled by the shared
// expression parser). A bare START WITH with no following CONNECT BY is a
// syntax error, surfaced by the CONNECT BY expect below.
func (p *Parser) parseConnectBy(stmt *ast.SelectStmt) error {
	if p.cur.Type != kwSTART && p.cur.Type != kwCONNECT {
		return nil
	}

	// Optional START WITH <condition>.
	if p.cur.Type == kwSTART {
		p.advance() // consume START
		if _, err := p.expect(kwWITH); err != nil {
			return err
		}
		cond, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.StartWith = cond
	}

	// CONNECT BY <condition> [ , … ] (required when hierarchical clauses appear).
	if _, err := p.expect(kwCONNECT); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	// Enable PRIOR-as-prefix only for the duration of the CONNECT BY conditions.
	saved := p.inConnectBy
	p.inConnectBy = true
	conds, err := p.parseExprList()
	p.inConnectBy = saved
	if err != nil {
		return err
	}
	stmt.ConnectBy = conds

	return nil
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

	// FETCH [FIRST|NEXT] n [ROW|ROWS] [ONLY]
	//
	// FIRST/NEXT, ROW/ROWS, and ONLY are all optional noise words in
	// Snowflake: `FETCH 123` is equivalent to `FETCH FIRST 123 ROWS ONLY`.
	if p.cur.Type == kwFETCH {
		fetchTok := p.advance() // consume FETCH

		// Optional FIRST or NEXT
		if p.cur.Type == kwFIRST || p.cur.Type == kwNEXT {
			p.advance() // consume FIRST or NEXT
		}

		// Parse count expression
		countExpr, err := p.parseExpr()
		if err != nil {
			return err
		}

		// Optional ROWS or ROW
		if p.cur.Type == kwROWS || p.cur.Type == kwROW {
			p.advance() // consume ROWS or ROW
		}

		// Optional ONLY
		if p.cur.Type == kwONLY {
			p.advance() // consume ONLY
		}

		stmt.Fetch = &ast.FetchClause{
			Count: countExpr,
			// p.prev is the last token consumed: ONLY, ROW/ROWS, or the
			// end of the count expression.
			Loc: ast.Loc{Start: fetchTok.Loc.Start, End: p.prev.Loc.End},
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
