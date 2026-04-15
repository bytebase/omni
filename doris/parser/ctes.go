package parser

import (
	"github.com/bytebase/omni/doris/ast"
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
	withTok, err := p.expect(kwWITH)
	if err != nil {
		return nil, err
	}

	with := &ast.WithClause{
		Loc: ast.Loc{Start: withTok.Loc.Start},
	}

	// Optional RECURSIVE keyword.
	if p.cur.Kind == kwRECURSIVE {
		p.advance()
		with.Recursive = true
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

	// The trailing SELECT statement (WITH must be followed by SELECT).
	stmt, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	// Attach the WITH clause and extend the overall Loc to cover WITH.
	stmt.With = with
	stmt.Loc.Start = withTok.Loc.Start

	return stmt, nil
}

// parseCTE parses one CTE definition:
//
//	cte_name [(col1, col2, ...)] AS ( query )
func (p *Parser) parseCTE() (*ast.CTE, error) {
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
	// a plain SELECT.
	var innerStmt ast.Node
	if p.cur.Kind == kwWITH {
		innerStmt, err = p.parseWithSelect()
	} else {
		innerStmt, err = p.parseSelectStmt()
	}
	if err != nil {
		return nil, err
	}
	cte.Query = innerStmt

	// Closing paren.
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	cte.Loc.End = p.prev.Loc.End
	return cte, nil
}
