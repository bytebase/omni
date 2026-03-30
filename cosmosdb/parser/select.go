package parser

import (
	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// parseTopClause parses TOP n.
func (p *Parser) parseTopClause() (*int, error) {
	p.advance() // consume TOP
	val, _, err := p.parseIntLiteral()
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// parseTargetList parses a comma-separated list of target entries.
func (p *Parser) parseTargetList() ([]*nodes.TargetEntry, error) {
	var targets []*nodes.TargetEntry
	for {
		entry, err := p.parseTargetEntry()
		if err != nil {
			return nil, err
		}
		targets = append(targets, entry)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return targets, nil
}

// parseTargetEntry parses: scalar_expression [AS? alias].
func (p *Parser) parseTargetEntry() (*nodes.TargetEntry, error) {
	startLoc := p.pos()
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	entry := &nodes.TargetEntry{Expr: expr}
	endLoc := p.prev.Loc + len(p.prev.Str)

	// Optional alias: AS identifier or just identifier (if not a clause start).
	if p.cur.Type == tokAS {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		entry.Alias = &name
		endLoc = p.prev.Loc + len(p.prev.Str)
	} else if p.cur.Type == tokIDENT {
		// Implicit alias (no AS keyword) - only if it looks like an alias, not a clause keyword.
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		entry.Alias = &name
		endLoc = loc + len(name)
	}

	entry.Loc = nodes.Loc{Start: startLoc, End: endLoc}
	return entry, nil
}

// parseGroupBy parses GROUP BY expr, ...
func (p *Parser) parseGroupBy() ([]nodes.ExprNode, error) {
	p.advance() // consume GROUP
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}

	var exprs []nodes.ExprNode
	for {
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return exprs, nil
}

// parseOrderBy parses ORDER BY sort_expr, ...
func (p *Parser) parseOrderBy() ([]*nodes.SortExpr, error) {
	p.advance() // consume ORDER
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}

	var sorts []*nodes.SortExpr
	for {
		se, err := p.parseSortExpr()
		if err != nil {
			return nil, err
		}
		sorts = append(sorts, se)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return sorts, nil
}

// parseSortExpr parses: expr [ASC|DESC] | RANK expr.
func (p *Parser) parseSortExpr() (*nodes.SortExpr, error) {
	startLoc := p.pos()

	// RANK expr
	if p.cur.Type == tokRANK {
		p.advance()
		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.SortExpr{
			Expr: expr,
			Rank: true,
			Loc:  nodes.Loc{Start: startLoc, End: endLoc},
		}, nil
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	endLoc := p.prev.Loc + len(p.prev.Str)

	se := &nodes.SortExpr{Expr: expr}
	if p.cur.Type == tokASC {
		p.advance()
		endLoc = p.prev.Loc + len(p.prev.Str)
	} else if p.cur.Type == tokDESC {
		se.Desc = true
		p.advance()
		endLoc = p.prev.Loc + len(p.prev.Str)
	}

	se.Loc = nodes.Loc{Start: startLoc, End: endLoc}
	return se, nil
}

// parseOffsetLimit parses OFFSET n LIMIT m.
func (p *Parser) parseOffsetLimit() (*nodes.OffsetLimitClause, error) {
	startLoc := p.pos()
	p.advance() // consume OFFSET

	offset, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLIMIT); err != nil {
		return nil, err
	}

	limit, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	endLoc := p.prev.Loc + len(p.prev.Str)
	return &nodes.OffsetLimitClause{
		Offset: offset,
		Limit:  limit,
		Loc:    nodes.Loc{Start: startLoc, End: endLoc},
	}, nil
}
