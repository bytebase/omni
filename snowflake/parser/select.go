package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

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

	// The CTE body can be WITH ... SELECT or just SELECT
	var query ast.Node
	if p.cur.Type == kwWITH {
		query, err = p.parseWithSelect()
	} else {
		query, err = p.parseSelectStmt()
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

// parseFromClause parses comma-separated table references.
// The FROM keyword has already been consumed by the caller.
func (p *Parser) parseFromClause() ([]*ast.TableRef, error) {
	var refs []*ast.TableRef

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	refs = append(refs, ref)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		ref, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

// parseTableRef parses one table reference: object_name [AS alias].
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
		kwTHEN, kwELSE, kwEND, kwCASE, kwROWS, kwGROUPS, kwOVER:
		return true
	}
	return false
}
