package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseSelect parses a SELECT statement.
//
//	SELECT [DISTINCT] [JSON] selectElements
//	  FROM fromSpecElement
//	  [WHERE relationElements]
//	  [ORDER BY orderSpecElement (',' orderSpecElement)*]
//	  [LIMIT decimal]
//	  [ALLOW FILTERING]
func (p *Parser) parseSelect() (*ast.SelectStmt, error) {
	start := p.curLoc()
	if err := p.expectKeyword(tokSELECT); err != nil {
		return nil, err
	}

	stmt := &ast.SelectStmt{}

	// Optional DISTINCT.
	if p.cur.Type == tokDISTINCT {
		stmt.Distinct = true
		p.advance()
	}

	// Optional JSON.
	if p.cur.Type == tokJSON {
		stmt.JSON = true
		p.advance()
	}

	// Parse select elements.
	elements, err := p.parseSelectElements()
	if err != nil {
		return nil, err
	}
	stmt.Elements = elements

	// FROM tableName.
	if err := p.expectKeyword(tokFROM); err != nil {
		return nil, err
	}
	from, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.From = from

	// Optional WHERE clause.
	if p.cur.Type == tokWHERE {
		where, err := p.parseWhereClause()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// Optional GROUP BY clause.
	if p.cur.Type == tokGROUP {
		groupBy, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
	}

	// Optional ORDER BY clause.
	if p.cur.Type == tokORDER {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// Optional PER PARTITION LIMIT clause.
	if p.cur.Type == tokPER {
		perPartLimit, err := p.parsePerPartitionLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.PerPartitionLimit = perPartLimit
	}

	// Optional LIMIT clause.
	if p.cur.Type == tokLIMIT {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// Optional ALLOW FILTERING.
	if p.cur.Type == tokALLOW {
		p.advance()
		if err := p.expectKeyword(tokFILTERING); err != nil {
			return nil, err
		}
		stmt.AllowFiltering = true
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseSelectElements parses the select clause elements:
//
//	'*' | selectElement (',' selectElement)*
func (p *Parser) parseSelectElements() ([]*ast.SelectElement, error) {
	// Handle bare '*'.
	if p.cur.Type == tokSTAR {
		elemStart := p.curLoc()
		p.advance()
		star := &ast.StarExpr{Loc: p.makeLoc(elemStart)}
		elem := &ast.SelectElement{Expr: star, Loc: p.makeLoc(elemStart)}
		return []*ast.SelectElement{elem}, nil
	}

	var elements []*ast.SelectElement
	first, err := p.parseSelectElement()
	if err != nil {
		return nil, err
	}
	elements = append(elements, first)

	for p.cur.Type == tokCOMMA {
		p.advance()
		el, err := p.parseSelectElement()
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
	}
	return elements, nil
}

// parseSelectElement parses a single select element:
//
//	IDENT '.' '*'          -> DotAccess with StarExpr-like field
//	IDENT [AS IDENT]       -> Identifier with optional alias
//	functionCall [AS IDENT] -> FunctionCall with optional alias
func (p *Parser) parseSelectElement() (*ast.SelectElement, error) {
	elemStart := p.curLoc()

	// CAST(expr AS type) in SELECT element
	if p.cur.Type == tokCAST {
		castExpr, err := p.parseCast()
		if err != nil {
			return nil, err
		}
		alias, err := p.parseOptionalAlias()
		if err != nil {
			return nil, err
		}
		return &ast.SelectElement{Expr: castExpr, Alias: alias, Loc: p.makeLoc(elemStart)}, nil
	}

	// We need an identifier-like token to start.
	if !isIdentLike(p.cur.Type) {
		return nil, p.errorf("expected column name or function call in SELECT, got %s", p.tokenDesc())
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for IDENT '.' '*' (qualified star, e.g. ks.*).
	if p.cur.Type == tokDOT {
		next := p.peekNext()
		if next.Type == tokSTAR {
			p.advance() // consume '.'
			starStart := p.curLoc()
			p.advance() // consume '*'
			starField := &ast.Identifier{Name: "*", Loc: p.makeLoc(starStart)}
			dotAccess := &ast.DotAccess{Object: name, Field: starField, Loc: p.makeLoc(elemStart)}
			elem := &ast.SelectElement{Expr: dotAccess, Loc: p.makeLoc(elemStart)}
			return elem, nil
		}
	}

	// Check for function call: IDENT '(' ... ')'.
	if p.cur.Type == tokLPAREN {
		fc, err := p.parseFunctionCallWithName(name)
		if err != nil {
			return nil, err
		}
		alias, err := p.parseOptionalAlias()
		if err != nil {
			return nil, err
		}
		elem := &ast.SelectElement{Expr: fc, Alias: alias, Loc: p.makeLoc(elemStart)}
		return elem, nil
	}

	// Plain identifier, optionally with alias.
	alias, err := p.parseOptionalAlias()
	if err != nil {
		return nil, err
	}
	elem := &ast.SelectElement{Expr: name, Alias: alias, Loc: p.makeLoc(elemStart)}
	return elem, nil
}

// parseOptionalAlias parses an optional AS IDENT clause.
func (p *Parser) parseOptionalAlias() (*ast.Identifier, error) {
	if p.cur.Type != tokAS {
		return nil, nil
	}
	p.advance() // consume AS
	return p.parseIdentifier()
}

// parseOrderByClause parses ORDER BY orderSpecElement (',' orderSpecElement)*.
func (p *Parser) parseOrderByClause() ([]*ast.OrderByElement, error) {
	if err := p.expectKeyword(tokORDER); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokBY); err != nil {
		return nil, err
	}

	var elements []*ast.OrderByElement
	first, err := p.parseOrderByElement()
	if err != nil {
		return nil, err
	}
	elements = append(elements, first)

	for p.cur.Type == tokCOMMA {
		p.advance()
		el, err := p.parseOrderByElement()
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
	}
	return elements, nil
}

// parseOrderByElement parses a single ORDER BY element:
//
//	IDENT [ASC | DESC]
//	IDENT ANN OF vectorLiteral [LIMIT DECIMAL]
func (p *Parser) parseOrderByElement() (*ast.OrderByElement, error) {
	elemStart := p.curLoc()

	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	elem := &ast.OrderByElement{Column: col}

	// Check for ANN OF ... ordering.
	if p.cur.Type == tokANN {
		p.advance() // consume ANN
		if err := p.expectKeyword(tokOF); err != nil {
			return nil, err
		}
		elem.IsANN = true

		vec, err := p.parseVectorLiteral()
		if err != nil {
			return nil, err
		}
		elem.AnnVector = vec

		// Optional LIMIT within ANN ordering.
		if p.cur.Type == tokLIMIT {
			p.advance()
			limit, err := p.parseConstant()
			if err != nil {
				return nil, err
			}
			elem.AnnLimit = limit
		}

		elem.Loc = p.makeLoc(elemStart)
		return elem, nil
	}

	// Optional ASC / DESC.
	switch p.cur.Type {
	case tokASC:
		elem.Direction = "ASC"
		p.advance()
	case tokDESC:
		elem.Direction = "DESC"
		p.advance()
	}

	elem.Loc = p.makeLoc(elemStart)
	return elem, nil
}

// parseVectorLiteral parses a vector literal: '[' constant (',' constant)* ']'.
func (p *Parser) parseVectorLiteral() (*ast.VectorLit, error) {
	start := p.curLoc()
	if _, err := p.expect(tokLBRACK); err != nil {
		return nil, err
	}

	var elements []ast.ExprNode
	first, err := p.parseConstant()
	if err != nil {
		return nil, err
	}
	elements = append(elements, first)

	for p.cur.Type == tokCOMMA {
		p.advance()
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		elements = append(elements, val)
	}

	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}

	return &ast.VectorLit{Elements: elements, Loc: p.makeLoc(start)}, nil
}

// parseGroupByClause parses GROUP BY column (',' column)*.
func (p *Parser) parseGroupByClause() ([]*ast.Identifier, error) {
	if err := p.expectKeyword(tokGROUP); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokBY); err != nil {
		return nil, err
	}

	var cols []*ast.Identifier
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cols = append(cols, first)

	for p.cur.Type == tokCOMMA {
		p.advance()
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

// parsePerPartitionLimitClause parses PER PARTITION LIMIT decimal.
func (p *Parser) parsePerPartitionLimitClause() (ast.ExprNode, error) {
	if err := p.expectKeyword(tokPER); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokPARTITION); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokLIMIT); err != nil {
		return nil, err
	}
	return p.parseLimitValue()
}

// parseLimitClause parses LIMIT decimal.
func (p *Parser) parseLimitClause() (ast.ExprNode, error) {
	p.advance() // consume LIMIT
	return p.parseLimitValue()
}

// parseLimitValue parses an integer literal, bind marker, or named bind marker.
func (p *Parser) parseLimitValue() (ast.ExprNode, error) {
	if p.cur.Type == tokQMARK {
		m := &ast.BindMarker{Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		return m, nil
	}
	if p.cur.Type == tokCOLON {
		start := p.curLoc()
		p.advance()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.BindMarker{Name: name.Name, Loc: p.makeLoc(start)}, nil
	}
	if p.cur.Type == tokINTEGER {
		lit := &ast.IntegerLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		return lit, nil
	}
	return nil, p.errorf("expected integer literal or bind marker for LIMIT, got %s", p.tokenDesc())
}
