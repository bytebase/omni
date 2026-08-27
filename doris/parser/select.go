package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// SELECT statement parser (T1.4)
// ---------------------------------------------------------------------------

// parseSelectStmt parses a full SELECT statement:
//
//	SELECT [DISTINCT|ALL] select_list
//	  [FROM table_references]
//	  [WHERE condition]
//	  [GROUP BY expr, ...]
//	  [HAVING condition]
//	  [QUALIFY condition]
//	  [ORDER BY expr [ASC|DESC] [NULLS FIRST|LAST], ...]
//	  [LIMIT count [OFFSET offset]]
func (p *Parser) parseSelectStmt() (*ast.SelectStmt, error) {
	selectTok, err := p.expect(kwSELECT)
	if err != nil {
		return nil, err
	}

	stmt := &ast.SelectStmt{
		Loc: ast.Loc{Start: selectTok.Loc.Start},
	}

	// DISTINCT / ALL
	if p.cur.Kind == kwDISTINCT {
		p.advance()
		stmt.Distinct = true
	} else if p.cur.Kind == kwALL {
		p.advance()
		stmt.All = true
	}

	// SELECT list
	items, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	stmt.Items = items

	// FROM clause
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY clause
	if p.cur.Kind == kwGROUP {
		groupBy, withRollup, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
		stmt.GroupByWithRollup = withRollup
	}

	// HAVING clause
	if p.cur.Kind == kwHAVING {
		p.advance() // consume HAVING
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// QUALIFY clause (Doris extension)
	if p.cur.Kind == kwQUALIFY {
		p.advance() // consume QUALIFY
		qualify, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Qualify = qualify
	}

	// ORDER BY clause
	if p.cur.Kind == kwORDER {
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

	// LIMIT / OFFSET
	if p.cur.Kind == kwLIMIT {
		limit, offset, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
		stmt.Offset = offset
	}

	// Set End location to the last consumed token.
	stmt.Loc.End = p.prev.Loc.End

	return stmt, nil
}

// ---------------------------------------------------------------------------
// Set operations (UNION / INTERSECT / EXCEPT / MINUS) — T1.7
// ---------------------------------------------------------------------------

// parseSetOpTail checks whether the token stream continues with a set
// operator (UNION, INTERSECT, EXCEPT, MINUS). If so, it loops consuming
// set-op tokens and building a left-associative SetOpStmt tree.
//
// Precedence note: The SQL standard gives INTERSECT higher precedence than
// UNION/EXCEPT. This implementation uses a two-level loop:
//
//  1. An inner loop collects a chain of INTERSECT clauses (higher priority).
//  2. An outer loop then folds UNION/EXCEPT operators at lower priority.
//
// This matches Doris/MySQL behaviour where INTERSECT binds more tightly.
func (p *Parser) parseSetOpTail(left ast.Node) (ast.Node, error) {
	// Collect the left side of any pending UNION/EXCEPT operation, starting
	// by giving INTERSECT first crack.
	left, err := p.parseIntersectChain(left)
	if err != nil {
		return nil, err
	}

	// Outer loop: UNION / EXCEPT (including MINUS alias) – left-associative.
	for {
		var op ast.SetOperator
		switch p.cur.Kind {
		case kwUNION:
			op = ast.SetUnion
		case kwEXCEPT, kwMINUS:
			op = ast.SetExcept
		default:
			return left, nil
		}

		opTok := p.advance() // consume UNION / EXCEPT / MINUS
		_ = opTok

		// Optional ALL / DISTINCT quantifier
		all := false
		if p.cur.Kind == kwALL {
			p.advance()
			all = true
		} else if p.cur.Kind == kwDISTINCT {
			p.advance()
			// DISTINCT is the default — all stays false
		}

		// Parse the right-hand side: a SELECT or a parenthesized operand.
		var rightSelect ast.Node
		var err error
		if p.cur.Kind == int('(') {
			rightSelect, err = p.parseParenQueryOperand()
		} else if p.cur.Kind == kwSELECT {
			rightSelect, err = p.parseSelectStmt()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		if err != nil {
			return nil, err
		}
		// Give INTERSECT a chance to grab more operands on the right.
		right, err := p.parseIntersectChain(rightSelect)
		if err != nil {
			return nil, err
		}

		left = &ast.SetOpStmt{
			Op:    op,
			All:   all,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}
	}
}

// parseParenQueryStmt parses a top-level parenthesized query — (SELECT 1),
// ((SELECT 1)), (SELECT 1) UNION (SELECT 2) — all engine-verified accepts.
// The parens are grouping only: the inner query node is returned directly,
// so analysis sees an ordinary SelectStmt/SetOpStmt.
func (p *Parser) parseParenQueryStmt() (ast.Node, error) {
	inner, err := p.parseParenQueryOperand()
	if err != nil {
		return nil, err
	}
	return p.parseQueryTail(inner)
}

// parseQueryTail applies the set-operator tail to a parsed query operand,
// then any trailing ORDER BY / LIMIT the operands did not consume. The
// engine is lenient here — SELECT 1 LIMIT 5 ORDER BY 1, SELECT 1 UNION
// (SELECT 2) LIMIT 5, and even repeated groups like SELECT 1 ORDER BY 1
// ORDER BY 2 or (SELECT 1 ORDER BY 1) ORDER BY 1 are all engine-verified
// accepts — so the attach loops, and where a clause repeats the outer
// (last, semantically governing) one wins.
func (p *Parser) parseQueryTail(node ast.Node) (ast.Node, error) {
	node, err := p.parseSetOpTail(node)
	if err != nil {
		return nil, err
	}
	return p.parseTrailingQueryClauses(node)
}

// parseTrailingQueryClauses parses trailing ORDER BY / LIMIT groups
// following a query and attaches them to the node, extending its range over
// the consumed clauses.
func (p *Parser) parseTrailingQueryClauses(node ast.Node) (ast.Node, error) {
	switch node.(type) {
	case *ast.SelectStmt, *ast.SetOpStmt:
	default:
		return node, nil
	}
	for p.cur.Kind == kwORDER || p.cur.Kind == kwLIMIT {
		var orderBy []*ast.OrderByItem
		var limit, offset ast.Node
		if p.cur.Kind == kwORDER {
			p.advance() // consume ORDER
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			var err error
			orderBy, err = p.parseOrderByList()
			if err != nil {
				return nil, err
			}
		}
		if p.cur.Kind == kwLIMIT {
			var err error
			limit, offset, err = p.parseLimitClause()
			if err != nil {
				return nil, err
			}
		}
		switch n := node.(type) {
		case *ast.SelectStmt:
			if len(orderBy) > 0 {
				n.OrderBy = orderBy
			}
			if limit != nil {
				n.Limit = limit
				n.Offset = offset
			}
			n.Loc.End = p.prev.Loc.End
		case *ast.SetOpStmt:
			if len(orderBy) > 0 {
				n.OrderBy = orderBy
			}
			if limit != nil {
				n.Limit = limit
				n.Offset = offset
			}
			n.Loc.End = p.prev.Loc.End
		}
	}
	return node, nil
}

// parseParenQueryOperand parses one parenthesized query operand, nesting
// freely: '(' followed by a SELECT query, a WITH query, or another
// parenthesized operand, then ')'.
func (p *Parser) parseParenQueryOperand() (ast.Node, error) {
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}
	var inner ast.Node
	switch p.cur.Kind {
	case int('('):
		inner, err = p.parseParenQueryOperand()
	case kwSELECT:
		inner, err = p.parseSelectStmt()
	case kwWITH:
		inner, err = p.parseWithSelect()
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	// A set-op tail (and, after one, trailing clauses) may follow any operand
	// inside the parens — ((SELECT 1) UNION SELECT 2) and (SELECT 1 UNION
	// (SELECT 2) LIMIT 5) are engine-verified accepts.
	inner, err = p.parseQueryTail(inner)
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	// The parens are grouping only, but the statement's source range must
	// still cover them, or consumers slicing by NodeLoc drop the delimiters.
	switch n := inner.(type) {
	case *ast.SelectStmt:
		n.Loc.Start = openTok.Loc.Start
		n.Loc.End = closeTok.Loc.End
	case *ast.SetOpStmt:
		n.Loc.Start = openTok.Loc.Start
		n.Loc.End = closeTok.Loc.End
	}
	return inner, nil
}

// parseIntersectChain collects a left-associative chain of INTERSECT clauses
// starting from an already-parsed left node.
func (p *Parser) parseIntersectChain(left ast.Node) (ast.Node, error) {
	for p.cur.Kind == kwINTERSECT {
		p.advance() // consume INTERSECT

		// Optional ALL / DISTINCT quantifier
		all := false
		if p.cur.Kind == kwALL {
			p.advance()
			all = true
		} else if p.cur.Kind == kwDISTINCT {
			p.advance()
		}

		var right ast.Node
		var err error
		if p.cur.Kind == int('(') {
			right, err = p.parseParenQueryOperand()
		} else if p.cur.Kind == kwSELECT {
			right, err = p.parseSelectStmt()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		if err != nil {
			return nil, err
		}

		left = &ast.SetOpStmt{
			Op:    ast.SetIntersect,
			All:   all,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}
	}
	return left, nil
}

// ---------------------------------------------------------------------------
// SELECT list
// ---------------------------------------------------------------------------

// parseSelectList parses comma-separated SELECT items.
func (p *Parser) parseSelectList() ([]*ast.SelectItem, error) {
	var items []*ast.SelectItem

	item, err := p.parseSelectItem()
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		item, err = p.parseSelectItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// parseSelectItem parses one item in the SELECT list:
//   - *                  — all columns
//   - table.*            — all columns from a specific table
//   - expr [AS alias]    — a computed expression with optional alias
func (p *Parser) parseSelectItem() (*ast.SelectItem, error) {
	startLoc := p.cur.Loc

	// Bare star: SELECT *  (optionally followed by EXCEPT(col, ...) for column exclusion)
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		item := &ast.SelectItem{
			Star: true,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}
		// SELECT * EXCEPT(col, col, ...)
		if p.cur.Kind == kwEXCEPT && p.peekNext().Kind == int('(') {
			p.advance() // consume EXCEPT
			p.advance() // consume '('
			for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
				name, _, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				item.ExceptColumns = append(item.ExceptColumns, name)
				if p.cur.Kind == int(',') {
					p.advance()
				}
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			item.Loc.End = p.prev.Loc.End
		}
		return item, nil
	}

	// Qualified star: ident.* or ident.ident.*. This lookahead exists ONLY
	// for the star form — any other qualified name must go through the
	// general expression path below, because an expression can continue
	// after it (t.a + 1, t.arr[1], db.f(1) OVER (...)); returning a bare
	// ColumnRef here would hand that continuation to the trailing-token
	// swallow. The multipart identifier is re-parsed after the rollback —
	// a few tokens, same cost class as the paren-lambda speculation.
	if p.isSelectIdentToken() && p.peekNext().Kind == int('.') {
		saved := p.save()
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		if p.cur.Kind == int('.') && p.peekNext().Kind == int('*') {
			p.advance() // consume '.'
			p.advance() // consume '*'
			return &ast.SelectItem{
				Star:      true,
				TableName: name,
				Loc:       ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}, nil
		}
		p.restore(saved)
	}

	// General expression
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	item := &ast.SelectItem{
		Expr: expr,
		Loc:  ast.Loc{Start: startLoc.Start},
	}

	alias, aliased := p.parseOptionalAlias(true)
	if aliased {
		item.Alias = alias
		item.Aliased = true
	}

	item.Loc.End = p.prev.Loc.End
	return item, nil
}

// isSelectIdentToken reports whether the current token can start an identifier
// in the SELECT list context (same as isExprIdentToken).
func (p *Parser) isSelectIdentToken() bool {
	return p.isExprIdentToken()
}

// parseOptionalAlias parses an optional alias after an expression in SELECT
// or FROM context. Returns the alias string, or empty if no alias.
//
// Alias forms:
//   - AS identifier
//   - AS 'string' (SELECT items only: identifierOrText allows AS "20%",
//     while a table alias is a strictIdentifier and rejects strings)
//   - identifier (implicit, if not a clause keyword)
//
// stringOK selects between those two grammar rules.
//
// The second return reports whether an alias was present at all: the empty
// string alias AS ” is engine-valid and distinct from "no alias", so the
// value alone cannot carry presence.
func (p *Parser) parseOptionalAlias(stringOK bool) (string, bool) {
	// Explicit: AS alias
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		if stringOK && p.cur.Kind == tokString {
			return p.advance().Str, true
		}
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return "", false
		}
		return name, true
	}

	// Implicit alias: current token is an identifier or non-reserved keyword
	// that does NOT start a clause.
	if p.isAliasIdentToken() {
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return "", false
		}
		return name, true
	}

	return "", false
}

// isAliasIdentToken reports whether the current token can be used as an
// implicit alias. Must be an identifier-like token that is NOT a clause keyword.
func (p *Parser) isAliasIdentToken() bool {
	switch p.cur.Kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		if p.cur.Kind >= 700 && !IsReserved(p.cur.Kind) && !isSelectClauseKeyword(p.cur.Kind) {
			return true
		}
		return false
	}
}

// parseAliasIdentifier parses a single identifier that serves as an alias.
// This accepts identifiers and non-reserved keywords.
func (p *Parser) parseAliasIdentifier() (string, ast.Loc, error) {
	tok := p.cur
	switch tok.Kind {
	case tokIdent, tokQuotedIdent:
		p.advance()
		return tok.Str, tok.Loc, nil
	default:
		// Non-reserved keywords may be used as aliases.
		if tok.Kind >= 700 && !IsReserved(tok.Kind) {
			p.advance()
			return tok.Str, tok.Loc, nil
		}
		return "", ast.Loc{}, p.syntaxErrorAtCur()
	}
}

// isSelectClauseKeyword returns true for keywords that start SQL clauses
// and should NOT be consumed as implicit aliases in SELECT context.
func isSelectClauseKeyword(t int) bool {
	switch t {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwQUALIFY, kwORDER,
		kwLIMIT, kwOFFSET, kwUNION, kwEXCEPT, kwINTERSECT,
		kwINTO, kwON, kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL,
		kwCROSS, kwNATURAL, kwWITH, kwSELECT, kwSET,
		kwFOR, kwLOCK, kwUSING, kwOUTER:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// FROM clause
// ---------------------------------------------------------------------------

// parseFromClause parses comma-separated FROM items.
// The FROM keyword has already been consumed by the caller.
func (p *Parser) parseFromClause() ([]ast.Node, error) {
	var items []ast.Node

	item, err := p.parseFromItem()
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	for p.cur.Kind == int(',') {
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
//   - otherwise → parseTableRef (ObjectName + alias)
func (p *Parser) parsePrimarySource() (ast.Node, error) {
	startLoc := p.cur.Loc

	// Parenthesized: subquery
	if p.cur.Kind == int('(') {
		next := p.peekNext()
		if next.Kind == kwSELECT || next.Kind == kwWITH {
			// Subquery in FROM: (SELECT ...) [AS] alias
			openTok := p.advance() // consume '('
			subq, err := p.parseSubqueryPlaceholder(openTok.Loc.Start)
			if err != nil {
				return nil, err
			}
			ref := &ast.TableRef{
				Loc: ast.Loc{Start: startLoc.Start},
			}
			alias, _ := p.parseOptionalAlias(false)
			if alias != "" {
				ref.Alias = alias
			}
			// Legacy consumers still read the raw text out of Name.Parts[0];
			// Subquery is the discriminator — table names like `selected` or
			// a quoted `SELECT` must never be mistaken for query text.
			ref.Subquery = subq
			ref.Name = &ast.ObjectName{
				Parts: []string{subq.RawText},
				Loc:   subq.Loc,
			}
			ref.Loc.End = p.prev.Loc.End
			return ref, nil
		}
		// Parenthesized from-item: (t1 JOIN t2 ON ...) or implicit-cross-join list (t1, t2, t3).
		p.advance() // consume '('
		inner, err := p.parseFromItem()
		if err != nil {
			return nil, err
		}
		// Implicit cross-join: comma-separated table list inside parens.
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			right, err := p.parseFromItem()
			if err != nil {
				return nil, err
			}
			inner = &ast.JoinClause{
				Left:  inner,
				Right: right,
				Type:  ast.JoinCross,
				Loc:   ast.Loc{Start: ast.NodeLoc(inner).Start, End: ast.NodeLoc(right).End},
			}
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		return inner, nil
	}

	// Default: simple table reference (ObjectName + alias)
	return p.parseTableRef()
}

// parseTableRef parses one simple table reference with its suffixes, in
// grammar order (relationPrimary): name — or a table-valued function call —
// then TABLET(...), then the alias, then TABLESAMPLE(...) [REPEATABLE seed].
func (p *Parser) parseTableRef() (*ast.TableRef, error) {
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	ref := &ast.TableRef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	// Table-valued function: BACKENDS(), numbers("number" = "10"). The grammar
	// (#tableValuedFunction) takes a single identifier as the function name.
	if p.cur.Kind == int('(') && len(name.Parts) == 1 {
		p.advance() // consume '('
		fc := &ast.FuncCallExpr{
			Name: name,
			Loc:  ast.Loc{Start: name.Loc.Start},
		}
		if p.cur.Kind != int(')') {
			args, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			fc.Args = args
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
		ref.Func = fc
	}

	// TABLET(id, ...) sits between the name and the alias (grammar:
	// tabletList? tableAlias). Only TABLET followed by '(' is the clause.
	if p.cur.Kind == kwTABLET && p.peekNext().Kind == int('(') {
		p.advance() // consume TABLET
		p.advance() // consume '('
		for {
			idTok, err := p.expect(tokInt)
			if err != nil {
				return nil, err
			}
			ref.TabletIDs = append(ref.TabletIDs, idTok.Ival)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	alias, _ := p.parseOptionalAlias(false)
	if alias != "" {
		ref.Alias = alias
	}

	// TABLESAMPLE(...) [REPEATABLE seed] follows the alias (grammar:
	// tableAlias sample?).
	if p.cur.Kind == kwTABLESAMPLE && p.peekNext().Kind == int('(') {
		sample, err := p.parseTableSample()
		if err != nil {
			return nil, err
		}
		ref.Sample = sample
	}

	ref.Loc.End = p.prev.Loc.End
	return ref, nil
}

// parseTableSample parses TABLESAMPLE(n ROWS | n PERCENT | ) [REPEATABLE seed].
// On entry cur is TABLESAMPLE with '(' next.
func (p *Parser) parseTableSample() (*ast.TableSample, error) {
	p.advance() // consume TABLESAMPLE
	p.advance() // consume '('

	sample := &ast.TableSample{}
	if p.cur.Kind != int(')') {
		valTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		sample.Value = &ast.Literal{Kind: ast.LitInt, Value: valTok.Str, Loc: valTok.Loc}
		switch p.cur.Kind {
		case kwROWS, kwPERCENT:
			sample.Unit = strings.ToUpper(p.advance().Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	if p.cur.Kind == kwREPEATABLE {
		p.advance() // consume REPEATABLE
		seedTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		sample.Seed = &ast.Literal{Kind: ast.LitInt, Value: seedTok.Str, Loc: seedTok.Loc}
	}
	return sample, nil
}

// ---------------------------------------------------------------------------
// JOIN chain (basic for T1.4)
// ---------------------------------------------------------------------------

// parseJoinChain builds a left-associative JoinClause tree from any
// JOIN keywords following the left source.
func (p *Parser) parseJoinChain(left ast.Node) (ast.Node, error) {
	for {
		joinType, natural, ok := p.parseJoinKeywords()
		if !ok {
			break
		}

		// Skip optional Doris execution hints: [shuffle], [broadcast], etc.
		hints := p.parseJoinHints()

		right, err := p.parsePrimarySource()
		if err != nil {
			return nil, err
		}

		join := &ast.JoinClause{
			Type:    joinType,
			Left:    left,
			Right:   right,
			Natural: natural,
			Hints:   hints,
			Loc:     ast.Loc{Start: ast.NodeLoc(left).Start},
		}

		// Parse join condition
		switch {
		case joinType == ast.JoinCross || natural:
			// CROSS JOIN and NATURAL JOIN: no condition required

		case p.cur.Kind == kwON:
			p.advance() // consume ON
			onExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			join.On = onExpr

		case p.cur.Kind == kwUSING:
			p.advance() // consume USING
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			var cols []string
			colName, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			cols = append(cols, colName)
			for p.cur.Kind == int(',') {
				p.advance() // consume ','
				colName, _, err = p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				cols = append(cols, colName)
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			join.Using = cols
		}

		join.Loc.End = p.prev.Loc.End
		left = join
	}
	return left, nil
}

// parseJoinHints skips any Doris execution hint between join keywords and
// the right-side table reference. Hints have the form [identifier] or
// /*+ identifier */ (comment-style). Only bracket-style is handled here.
// Returns the consumed hint identifiers (may be empty).
func (p *Parser) parseJoinHints() []string {
	var hints []string
	for p.cur.Kind == int('[') {
		p.advance() // consume '['
		// Collect tokens until ']'
		for p.cur.Kind != int(']') && p.cur.Kind != tokEOF {
			if p.cur.Kind == tokIdent || p.cur.Kind == tokQuotedIdent ||
				(p.cur.Kind >= 700) {
				hints = append(hints, p.cur.Str)
			}
			p.advance()
		}
		if p.cur.Kind == int(']') {
			p.advance() // consume ']'
		}
	}
	return hints
}

// parseJoinKeywords checks whether the current token position starts a
// JOIN keyword sequence. If so, it consumes the tokens and returns
// (joinType, natural, true). If not, returns (0, false, false)
// without consuming any tokens.
func (p *Parser) parseJoinKeywords() (ast.JoinType, bool, bool) {
	// NATURAL [LEFT|RIGHT|FULL] [OUTER] JOIN
	if p.cur.Kind == kwNATURAL {
		next := p.peekNext()
		if next.Kind == kwJOIN {
			p.advance() // consume NATURAL
			p.advance() // consume JOIN
			return ast.JoinInner, true, true
		}
		if next.Kind == kwLEFT || next.Kind == kwRIGHT || next.Kind == kwFULL {
			p.advance() // consume NATURAL
			jt := p.consumeDirectionAndJoin()
			return jt, true, true
		}
		return 0, false, false
	}

	// INNER JOIN
	if p.cur.Kind == kwINNER {
		next := p.peekNext()
		if next.Kind == kwJOIN {
			p.advance() // consume INNER
			p.advance() // consume JOIN
			return ast.JoinInner, false, true
		}
		return 0, false, false
	}

	// LEFT variants: LEFT [OUTER] JOIN | LEFT SEMI JOIN | LEFT ANTI JOIN
	if p.cur.Kind == kwLEFT {
		next := p.peekNext()
		switch next.Kind {
		case kwSEMI:
			p.advance() // consume LEFT
			p.advance() // consume SEMI
			if p.cur.Kind != kwJOIN {
				return 0, false, false
			}
			p.advance() // consume JOIN
			return ast.JoinLeftSemi, false, true
		case kwANTI:
			p.advance() // consume LEFT
			p.advance() // consume ANTI
			if p.cur.Kind != kwJOIN {
				return 0, false, false
			}
			p.advance() // consume JOIN
			return ast.JoinLeftAnti, false, true
		case kwJOIN, kwOUTER:
			jt := p.consumeDirectionAndJoin()
			if jt == ast.JoinLeft {
				return jt, false, true
			}
		}
		return 0, false, false
	}

	// RIGHT variants: RIGHT [OUTER] JOIN | RIGHT SEMI JOIN | RIGHT ANTI JOIN
	if p.cur.Kind == kwRIGHT {
		next := p.peekNext()
		switch next.Kind {
		case kwSEMI:
			p.advance() // consume RIGHT
			p.advance() // consume SEMI
			if p.cur.Kind != kwJOIN {
				return 0, false, false
			}
			p.advance() // consume JOIN
			return ast.JoinRightSemi, false, true
		case kwANTI:
			p.advance() // consume RIGHT
			p.advance() // consume ANTI
			if p.cur.Kind != kwJOIN {
				return 0, false, false
			}
			p.advance() // consume JOIN
			return ast.JoinRightAnti, false, true
		case kwJOIN, kwOUTER:
			jt := p.consumeDirectionAndJoin()
			if jt == ast.JoinRight {
				return jt, false, true
			}
		}
		return 0, false, false
	}

	// FULL [OUTER] JOIN
	if p.cur.Kind == kwFULL {
		next := p.peekNext()
		if next.Kind == kwJOIN || next.Kind == kwOUTER {
			jt := p.consumeDirectionAndJoin()
			if jt == ast.JoinFull {
				return jt, false, true
			}
		}
		return 0, false, false
	}

	// CROSS JOIN
	if p.cur.Kind == kwCROSS {
		next := p.peekNext()
		if next.Kind == kwJOIN {
			p.advance() // consume CROSS
			p.advance() // consume JOIN
			return ast.JoinCross, false, true
		}
		return 0, false, false
	}

	// Bare JOIN (= INNER)
	if p.cur.Kind == kwJOIN {
		p.advance() // consume JOIN
		return ast.JoinInner, false, true
	}

	return 0, false, false
}

// consumeDirectionAndJoin consumes LEFT/RIGHT/FULL [OUTER] JOIN and returns
// the corresponding JoinType. The caller must have already checked that
// p.cur.Kind is kwLEFT, kwRIGHT, or kwFULL.
func (p *Parser) consumeDirectionAndJoin() ast.JoinType {
	dir := p.cur.Kind
	p.advance() // consume LEFT/RIGHT/FULL

	// Optional OUTER
	if p.cur.Kind == kwOUTER {
		p.advance() // consume OUTER
	}

	// Expect JOIN
	if p.cur.Kind != kwJOIN {
		return ast.JoinInner // sentinel for "not a join"
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

// ---------------------------------------------------------------------------
// GROUP BY clause
// ---------------------------------------------------------------------------

// parseGroupByClause parses the grouping specification after GROUP BY: a
// plain expression list (optionally followed by WITH ROLLUP), or exactly one
// of CUBE(...) / GROUPING SETS (...). Returns the grouping items and whether
// WITH ROLLUP was present.
func (p *Parser) parseGroupByClause() ([]ast.Node, bool, error) {
	p.advance() // consume GROUP
	if _, err := p.expect(kwBY); err != nil {
		return nil, false, err
	}

	// CUBE is a reserved keyword, so `GROUP BY CUBE(a, b)` cannot reach the
	// ordinary expression path the way ROLLUP and GROUPING SETS do.
	if p.cur.Kind == kwCUBE && p.peekNext().Kind == int('(') {
		fc, err := p.parseGroupingElementCall()
		if err != nil {
			return nil, false, err
		}
		// CUBE(...) is the entire grouping specification — the engine rejects
		// both `GROUP BY CUBE(a), b` and `GROUP BY CUBE(a), CUBE(b)`. Without
		// this check, returning here would silently discard everything after
		// the comma and hand downstream analysis an incomplete GROUP BY.
		if p.cur.Kind == int(',') {
			return nil, false, p.syntaxErrorAtCur()
		}
		return []ast.Node{fc}, false, nil
	}

	// GROUPING SETS (...) — like CUBE, the entire grouping specification.
	// GROUPING alone stays an ordinary identifier (it is non-reserved), so a
	// column named grouping still groups normally.
	if p.cur.Kind == kwGROUPING && p.peekNext().Kind == kwSETS {
		gs, err := p.parseGroupingSets()
		if err != nil {
			return nil, false, err
		}
		if p.cur.Kind == int(',') {
			return nil, false, p.syntaxErrorAtCur()
		}
		return []ast.Node{gs}, false, nil
	}

	list, err := p.parseExprList()
	if err != nil {
		return nil, false, err
	}

	// Optional WITH ROLLUP — the grammar allows it only on the plain
	// expression-list form, not after CUBE / GROUPING SETS.
	if p.cur.Kind == kwWITH && p.peekNext().Kind == kwROLLUP {
		p.advance() // consume WITH
		p.advance() // consume ROLLUP
		return list, true, nil
	}
	return list, false, nil
}

// parseGroupingSets parses GROUPING SETS ((a, b), (a), ()). On entry cur is
// GROUPING with SETS next. Every set is itself parenthesized and may be empty.
func (p *Parser) parseGroupingSets() (ast.Node, error) {
	startTok := p.advance() // consume GROUPING
	p.advance()             // consume SETS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	gs := &ast.GroupingSetsExpr{}
	for {
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		var set []ast.Node
		if p.cur.Kind != int(')') {
			exprs, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			set = exprs
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		gs.Sets = append(gs.Sets, set)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	gs.Loc = ast.Loc{Start: startTok.Loc.Start, End: closeTok.Loc.End}
	return gs, nil
}

// ---------------------------------------------------------------------------
// LIMIT / OFFSET clause
// ---------------------------------------------------------------------------

// parseLimitClause parses LIMIT count [OFFSET offset].
// Returns (limit, offset, error) where offset may be nil.
func (p *Parser) parseLimitClause() (ast.Node, ast.Node, error) {
	p.advance() // consume LIMIT

	limitExpr, err := p.parseExpr()
	if err != nil {
		return nil, nil, err
	}

	var offsetExpr ast.Node
	switch p.cur.Kind {
	case int(','):
		// MySQL form LIMIT offset, row_count: the expression parsed first is
		// the offset and the one after the comma is the count.
		p.advance() // consume ','
		offsetExpr = limitExpr
		limitExpr, err = p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
	case kwOFFSET:
		p.advance() // consume OFFSET
		offsetExpr, err = p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
	}

	return limitExpr, offsetExpr, nil
}
