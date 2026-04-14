package parser

import (
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
		groupBy, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
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

	// Bare star: SELECT *
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		return &ast.SelectItem{
			Star: true,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil
	}

	// Try to detect table.* pattern:
	// If we see an identifier followed by '.', it might be table.* or table.col.
	// We need to check for qualified star: ident.* or ident.ident.*
	if p.isSelectIdentToken() && p.peekNext().Kind == int('.') {
		// Save state to potentially backtrack.
		// Parse the multipart identifier first, then check for .*
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}

		// Check for .* after the multipart identifier
		if p.cur.Kind == int('.') && p.peekNext().Kind == int('*') {
			p.advance() // consume '.'
			p.advance() // consume '*'
			return &ast.SelectItem{
				Star:      true,
				TableName: name,
				Loc:       ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}, nil
		}

		// Not a qualified star — the multipart identifier is a column ref or
		// function call. Check if it's a function call.
		if p.cur.Kind == int('(') {
			fc, err := p.parseFuncCall(name)
			if err != nil {
				return nil, err
			}
			item := &ast.SelectItem{
				Expr: fc,
				Loc:  ast.Loc{Start: startLoc.Start},
			}
			alias := p.parseOptionalAlias()
			if alias != "" {
				item.Alias = alias
			}
			item.Loc.End = p.prev.Loc.End
			return item, nil
		}

		// Plain column reference — check for alias.
		colRef := &ast.ColumnRef{
			Name: name,
			Loc:  name.Loc,
		}
		item := &ast.SelectItem{
			Expr: colRef,
			Loc:  ast.Loc{Start: startLoc.Start},
		}
		alias := p.parseOptionalAlias()
		if alias != "" {
			item.Alias = alias
		}
		item.Loc.End = p.prev.Loc.End
		return item, nil
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

	alias := p.parseOptionalAlias()
	if alias != "" {
		item.Alias = alias
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
//   - identifier (implicit, if not a clause keyword)
func (p *Parser) parseOptionalAlias() string {
	// Explicit: AS alias
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return ""
		}
		return name
	}

	// Implicit alias: current token is an identifier or non-reserved keyword
	// that does NOT start a clause.
	if p.isAliasIdentToken() {
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return ""
		}
		return name
	}

	return ""
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
			alias := p.parseOptionalAlias()
			if alias != "" {
				ref.Alias = alias
			}
			// Store the subquery raw text in the table name for now
			ref.Name = &ast.ObjectName{
				Parts: []string{subq.RawText},
				Loc:   subq.Loc,
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
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		return inner, nil
	}

	// Default: simple table reference (ObjectName + alias)
	return p.parseTableRef()
}

// parseTableRef parses one simple table reference: object_name [AS alias].
func (p *Parser) parseTableRef() (*ast.TableRef, error) {
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	ref := &ast.TableRef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	alias := p.parseOptionalAlias()
	if alias != "" {
		ref.Alias = alias
	}

	ref.Loc.End = p.prev.Loc.End
	return ref, nil
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

// parseGroupByClause parses GROUP BY expr, expr, ...
// Returns the list of GROUP BY expressions.
func (p *Parser) parseGroupByClause() ([]ast.Node, error) {
	p.advance() // consume GROUP
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	return p.parseExprList()
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
	if p.cur.Kind == kwOFFSET {
		p.advance() // consume OFFSET
		offsetExpr, err = p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
	}

	return limitExpr, offsetExpr, nil
}
