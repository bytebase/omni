// select.go implements SELECT query parsing for the PartiQL parser.
//
// This file handles the SFW (SELECT-FROM-WHERE) query form and all its
// optional clauses: selectClause, fromClause (delegated to from.go),
// letClause (stub), whereClauseSelect, groupClause, havingClause,
// orderByClause, limitClause, offsetByClause.
//
// Grammar references cite PartiQLParser.g4 line numbers.
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseSFWQuery parses a full SELECT-FROM-WHERE query and returns it
// wrapped in an ast.SubLink (since SelectStmt is a StmtNode, not an
// ExprNode, and the caller expects ExprNode).
//
// Grammar: exprSelect#SfwQuery (lines 456-465).
//
//	selectClause fromClause letClause? whereClauseSelect?
//	groupClause? havingClause? orderByClause? limitClause? offsetByClause?
func (p *Parser) parseSFWQuery() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	stmt := &ast.SelectStmt{}

	// SELECT clause (required).
	if err := p.parseSelectClause(stmt); err != nil {
		return nil, err
	}

	// FROM clause (required).
	if p.cur.Type != tokFROM {
		return nil, &ParseError{
			Message: fmt.Sprintf("expected FROM, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	from, err := p.parseFromClause()
	if err != nil {
		return nil, err
	}
	stmt.From = from

	// LET clause (optional, deferred).
	if p.cur.Type == tokLET {
		return nil, p.deferredFeature("LET", "parser-let-pivot (DAG node 12)")
	}

	// WHERE clause (optional).
	// Grammar: whereClauseSelect (line 300-301).
	if p.cur.Type == tokWHERE {
		p.advance() // consume WHERE
		where, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// GROUP BY clause (optional).
	// Grammar: groupClause (lines 262-263).
	if p.cur.Type == tokGROUP {
		// Peek to distinguish GROUP BY from GROUP AS (group alias without
		// GROUP BY is not valid at this position).
		next := p.peekNext()
		if next.Type == tokPARTIAL || next.Type == tokBY {
			gb, err := p.parseGroupByClause()
			if err != nil {
				return nil, err
			}
			stmt.GroupBy = gb
		}
	}

	// HAVING clause (optional).
	// Grammar: havingClause (lines 294-295).
	if p.cur.Type == tokHAVING {
		p.advance() // consume HAVING
		having, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}

	// ORDER BY clause (optional).
	// Grammar: orderByClause (lines 250-251).
	if p.cur.Type == tokORDER {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT clause (optional).
	// Grammar: limitClause (lines 306-307).
	if p.cur.Type == tokLIMIT {
		p.advance() // consume LIMIT
		limit, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// OFFSET clause (optional).
	// Grammar: offsetByClause (lines 303-304).
	if p.cur.Type == tokOFFSET {
		p.advance() // consume OFFSET
		offset, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		stmt.Offset = offset
	}

	// Compute the Loc spanning the entire SFW query.
	stmt.Loc = ast.Loc{Start: start, End: p.prev.Loc.End}

	return &ast.SubLink{
		Stmt: stmt,
		Loc:  stmt.Loc,
	}, nil
}

// parseSelectClause parses the SELECT clause and populates the
// appropriate fields on stmt (Quantifier, Star, Targets, Value, Pivot).
//
// Grammar: selectClause (lines 216-221):
//
//	SELECT setQuantifierStrategy? ASTERISK        # SelectAll
//	SELECT setQuantifierStrategy? projectionItems # SelectItems
//	SELECT setQuantifierStrategy? VALUE expr      # SelectValue
//	PIVOT pivot=expr AT at=expr                   # SelectPivot
func (p *Parser) parseSelectClause(stmt *ast.SelectStmt) error {
	// PIVOT is a separate production (no SELECT keyword).
	if p.cur.Type == tokPIVOT {
		return p.deferredFeature("PIVOT", "parser-let-pivot (DAG node 12)")
	}

	if p.cur.Type != tokSELECT {
		return &ParseError{
			Message: fmt.Sprintf("expected SELECT, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	p.advance() // consume SELECT

	// Optional DISTINCT/ALL quantifier.
	// Grammar: setQuantifierStrategy (lines 229-232).
	if p.cur.Type == tokDISTINCT {
		stmt.Quantifier = ast.QuantifierDistinct
		p.advance()
	} else if p.cur.Type == tokALL {
		stmt.Quantifier = ast.QuantifierAll
		p.advance()
	}

	// SELECT * (SelectAll).
	if p.cur.Type == tokASTERISK {
		stmt.Star = true
		p.advance()
		return nil
	}

	// SELECT VALUE expr (SelectValue).
	if p.cur.Type == tokVALUE {
		p.advance() // consume VALUE
		val, err := p.parseSelectExpr()
		if err != nil {
			return err
		}
		stmt.Value = val
		return nil
	}

	// SELECT items (SelectItems).
	targets, err := p.parseProjectionItems()
	if err != nil {
		return err
	}
	stmt.Targets = targets
	return nil
}

// parseProjectionItems parses a comma-separated list of projection items.
// Each item is `expr [AS? alias]`.
//
// Grammar: projectionItems (line 223-224), projectionItem (line 226-227).
func (p *Parser) parseProjectionItems() ([]*ast.TargetEntry, error) {
	var items []*ast.TargetEntry
	for {
		item, err := p.parseProjectionItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	return items, nil
}

// parseProjectionItem parses a single `expr [AS? alias]` item.
//
// Grammar: projectionItem (line 226-227):
//
//	expr (AS? symbolPrimitive)?
func (p *Parser) parseProjectionItem() (*ast.TargetEntry, error) {
	expr, err := p.parseSelectExpr()
	if err != nil {
		return nil, err
	}
	te := &ast.TargetEntry{
		Expr: expr,
		Loc:  expr.GetLoc(),
	}

	// Optional alias: AS symbolPrimitive, or bare symbolPrimitive.
	// The grammar allows omitting AS (projectionItem rule line 227).
	if p.cur.Type == tokAS {
		p.advance() // consume AS
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		te.Alias = &name
		te.Loc = ast.Loc{Start: te.Loc.Start, End: nameLoc.End}
	} else if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		// Bare alias (no AS keyword). We need to be careful:
		// only consume if the next thing is an unquoted/quoted identifier
		// and NOT a keyword that starts the next clause.
		// The grammar rule says `AS? symbolPrimitive` so bare aliases
		// are allowed, but we must not consume identifiers that are
		// really the start of the FROM clause etc.
		// In practice, we only accept a bare alias if the identifier
		// is followed by a comma, FROM, or other clause boundary.
		// Actually, the PartiQL grammar does allow bare alias with
		// `AS?` making AS optional, so we parse it.
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		te.Alias = &name
		te.Loc = ast.Loc{Start: te.Loc.Start, End: nameLoc.End}
	}

	return te, nil
}

// parseOrderByClause parses ORDER BY with one or more sort specs.
//
// Grammar: orderByClause (lines 250-251):
//
//	ORDER BY orderSortSpec (COMMA orderSortSpec)*
func (p *Parser) parseOrderByClause() ([]*ast.OrderByItem, error) {
	if _, err := p.expect(tokORDER); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}
	var items []*ast.OrderByItem
	for {
		item, err := p.parseOrderSortSpec()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	return items, nil
}

// parseOrderSortSpec parses a single ORDER BY item.
//
// Grammar: orderSortSpec (lines 253-254):
//
//	expr dir=(ASC|DESC)? (NULLS nulls=(FIRST|LAST))?
func (p *Parser) parseOrderSortSpec() (*ast.OrderByItem, error) {
	expr, err := p.parseSelectExpr()
	if err != nil {
		return nil, err
	}
	item := &ast.OrderByItem{
		Expr: expr,
		Loc:  expr.GetLoc(),
	}

	// Optional ASC/DESC.
	if p.cur.Type == tokASC {
		item.Desc = false
		item.Loc = ast.Loc{Start: item.Loc.Start, End: p.cur.Loc.End}
		p.advance()
	} else if p.cur.Type == tokDESC {
		item.Desc = true
		item.Loc = ast.Loc{Start: item.Loc.Start, End: p.cur.Loc.End}
		p.advance()
	}

	// Optional NULLS FIRST/LAST.
	if p.cur.Type == tokNULLS {
		p.advance() // consume NULLS
		switch p.cur.Type {
		case tokFIRST:
			item.NullsFirst = true
			item.NullsExplicit = true
			item.Loc = ast.Loc{Start: item.Loc.Start, End: p.cur.Loc.End}
			p.advance()
		case tokLAST:
			item.NullsFirst = false
			item.NullsExplicit = true
			item.Loc = ast.Loc{Start: item.Loc.Start, End: p.cur.Loc.End}
			p.advance()
		default:
			return nil, &ParseError{
				Message: fmt.Sprintf("expected FIRST or LAST after NULLS, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
	}

	return item, nil
}

// parseGroupByClause parses GROUP [PARTIAL] BY items [GROUP AS alias].
//
// Grammar: groupClause (lines 262-263):
//
//	GROUP PARTIAL? BY groupKey (COMMA groupKey)* groupAlias?
func (p *Parser) parseGroupByClause() (*ast.GroupByClause, error) {
	start := p.cur.Loc.Start
	p.advance() // consume GROUP

	partial := false
	if p.cur.Type == tokPARTIAL {
		partial = true
		p.advance()
	}

	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}

	var items []*ast.GroupByItem
	for {
		item, err := p.parseGroupKey()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}

	end := p.prev.Loc.End

	// Optional GROUP AS alias.
	// Grammar: groupAlias (lines 265-266).
	var groupAs *string
	if p.cur.Type == tokGROUP {
		p.advance() // consume GROUP
		if _, err := p.expect(tokAS); err != nil {
			return nil, err
		}
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		groupAs = &name
		end = nameLoc.End
	}

	return &ast.GroupByClause{
		Partial: partial,
		Items:   items,
		GroupAs: groupAs,
		Loc:     ast.Loc{Start: start, End: end},
	}, nil
}

// parseGroupKey parses a single GROUP BY key item.
//
// Grammar: groupKey (lines 268-269):
//
//	key=exprSelect (AS symbolPrimitive)?
func (p *Parser) parseGroupKey() (*ast.GroupByItem, error) {
	expr, err := p.parseSelectExpr()
	if err != nil {
		return nil, err
	}
	item := &ast.GroupByItem{
		Expr: expr,
		Loc:  expr.GetLoc(),
	}

	if p.cur.Type == tokAS {
		p.advance() // consume AS
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		item.Alias = &name
		item.Loc = ast.Loc{Start: item.Loc.Start, End: nameLoc.End}
	}

	return item, nil
}
