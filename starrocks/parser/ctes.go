package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// WITH / CTE parsing (T1.6)
// ---------------------------------------------------------------------------

// parseWithSelect parses a WITH clause followed by a SELECT statement:
//
//	WITH [RECURSIVE] cte_name [(col, ...)] AS (query)
//	  [, cte_name2 [(col, ...)] AS (query2)]
//	SELECT ...
//
// The WITH keyword has NOT yet been consumed when this is called.
func (p *Parser) parseWithSelect() (*ast.SelectStmt, error) {
	with, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	// The trailing SELECT statement (WITH must be followed by SELECT here).
	stmt, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	// Attach the WITH clause and extend the overall Loc to cover WITH.
	stmt.With = with
	stmt.Loc.Start = with.Loc.Start

	return stmt, nil
}

// parseWithClause parses WITH [RECURSIVE] cte (',' cte)* and returns the clause,
// leaving the cursor at the statement keyword that follows it. The WITH keyword
// has NOT yet been consumed when this is called.
func (p *Parser) parseWithClause() (*ast.WithClause, error) {
	withTok, err := p.expect(kwWITH)
	if err != nil {
		return nil, err
	}

	with := &ast.WithClause{
		Loc: ast.Loc{Start: withTok.Loc.Start},
	}

	// Parse comma-separated CTE definitions.
	cte, err := p.parseCTE()
	if err != nil {
		return nil, err
	}
	with.CTEs = append(with.CTEs, cte)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		cte, err = p.parseCTE()
		if err != nil {
			return nil, err
		}
		with.CTEs = append(with.CTEs, cte)
	}

	with.Loc.End = p.prev.Loc.End
	return with, nil
}

// parseWithStatement parses a top-level WITH clause that may prefix a SELECT
// (the common case), a DELETE, or an UPDATE — StarRocks allows a CTE before all
// three (the grammar's deleteStatement/updateStatement both take withClause?).
// WITH … INSERT is not allowed and falls through to the SELECT path, which
// fails with a syntax error.
func (p *Parser) parseWithStatement() (ast.Node, error) {
	with, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	switch p.cur.Kind {
	case kwDELETE:
		deleteTok := p.advance() // consume DELETE
		del, err := p.parseDeleteStmt(deleteTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		del.With = with
		del.Loc.Start = with.Loc.Start
		return del, nil
	case kwUPDATE:
		updateTok := p.advance() // consume UPDATE
		upd, err := p.parseUpdateStmt(updateTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		upd.With = with
		upd.Loc.Start = with.Loc.Start
		return upd, nil
	case int('('):
		// withClause? queryNoWith, where queryNoWith may be a parenthesized
		// subquery (StarRocks.g4:2090-2118).
		paren, err := p.parseParenSelect()
		if err != nil {
			return nil, err
		}
		query, err := p.parseSetOpTail(paren)
		if err != nil {
			return nil, err
		}
		if err := p.settleQueryInto(query); err != nil {
			return nil, err
		}
		attachWith(query, with)
		return query, nil
	}

	stmt, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	// A WITH … SELECT may continue with set operators (UNION/INTERSECT/EXCEPT).
	query, err := p.parseSetOpTail(stmt)
	if err != nil {
		return nil, err
	}
	if err := p.settleQueryInto(query); err != nil {
		return nil, err
	}
	attachWith(query, with)
	return query, nil
}

// parseCTE parses one CTE definition:
//
//	cte_name [(col1, col2, ...)] AS ( query )
func (p *Parser) parseCTE() (*ast.CTE, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	startLoc := p.cur.Loc

	// CTE name — must be a plain or quoted identifier.
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	cte := &ast.CTE{
		Name: name,
		Loc:  ast.Loc{Start: startLoc.Start},
	}

	// Optional column list: (col1, col2, ...)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('

		colName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cte.Columns = append(cte.Columns, colName)

		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			colName, _, err = p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			cte.Columns = append(cte.Columns, colName)
		}

		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	// AS keyword.
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Opening paren of the inner query.
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Inner SELECT statement. If it starts with WITH, recurse; otherwise parse
	// a plain SELECT followed by optional set-operator tail (UNION/INTERSECT/EXCEPT).
	var innerStmt ast.Node
	switch p.cur.Kind {
	case kwWITH:
		innerStmt, err = p.parseWithStatement()
		if err == nil && !isQueryBody(innerStmt) {
			err = p.syntaxErrorAtCur()
		}
	case int('('):
		var paren *ast.ParenSelect
		paren, err = p.parseParenSelect()
		if err == nil {
			innerStmt, err = p.parseSetOpTail(paren)
		}
	default:
		var sel *ast.SelectStmt
		sel, err = p.parseSelectStmt()
		if err == nil {
			innerStmt, err = p.parseSetOpTail(sel)
		}
	}
	if err != nil {
		return nil, err
	}
	// A CTE body is a queryRelation, which has no outfile clause.
	if queryHasInto(innerStmt) {
		return nil, p.syntaxErrorAtCur()
	}
	cte.Query = innerStmt

	// Closing paren.
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	cte.Loc.End = p.prev.Loc.End
	return cte, nil
}
