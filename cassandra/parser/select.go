package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseSelect parses a SELECT statement.
//
//	SELECT [JSON] [DISTINCT] selectors
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

	elements, err := p.parseSelectClause(stmt)
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

// parseSelectClause parses (JSON)? (DISTINCT)? selectors, up to the FROM
// that ends them. JSON and DISTINCT may also name a column: each is read as
// a keyword only where the selectors after it parse and FROM follows them,
// as the predicates (K_JSON selectClause)=> and (K_DISTINCT selectors)=>
// decide (Parser.g:277, 300). SELECT DISTINCT FROM t selects the column
// distinct.
func (p *Parser) parseSelectClause(stmt *ast.SelectStmt) ([]*ast.SelectElement, error) {
	if p.cur.Type == tokJSON {
		s := p.mark()
		p.advance()
		elements, err := p.parseDistinctSelectors(stmt)
		if err == nil {
			stmt.JSON = true
			return elements, nil
		}
		p.reset(s)
	}
	return p.parseDistinctSelectors(stmt)
}

// parseDistinctSelectors parses (DISTINCT)? selectors, up to FROM.
func (p *Parser) parseDistinctSelectors(stmt *ast.SelectStmt) ([]*ast.SelectElement, error) {
	if p.cur.Type == tokDISTINCT {
		s := p.mark()
		p.advance()
		elements, err := p.parseSelectors()
		if err == nil && p.cur.Type == tokFROM {
			stmt.Distinct = true
			return elements, nil
		}
		p.reset(s)
	}
	elements, err := p.parseSelectors()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokFROM {
		return nil, p.errorf("expected FROM, got %s", p.tokenDesc())
	}
	return elements, nil
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

// parseOrderByElement parses a single ORDER BY element (Parser.g:462):
//
//	cident (ANN OF term)? (ASC | DESC)?
//
// A LIMIT after ANN OF is the statement's.
func (p *Parser) parseOrderByElement() (*ast.OrderByElement, error) {
	elemStart := p.curLoc()

	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	elem := &ast.OrderByElement{Column: col}

	if p.cur.Type == tokANN {
		p.advance() // consume ANN
		if err := p.expectKeyword(tokOF); err != nil {
			return nil, err
		}
		elem.IsANN = true
		vec, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		if list, ok := vec.(*ast.ListLit); ok {
			vec = &ast.VectorLit{Elements: list.Elements, Loc: list.Loc}
		}
		elem.AnnVector = vec
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
	return p.parseIntValue()
}

// parseLimitClause parses LIMIT decimal.
func (p *Parser) parseLimitClause() (ast.ExprNode, error) {
	p.advance() // consume LIMIT
	return p.parseIntValue()
}

// parseIntValue parses intValue (Parser.g:1601), the value of LIMIT, PER
// PARTITION LIMIT, TTL and TIMESTAMP: an INTEGER, a named marker or a
// positional one. INTEGER takes a '-' written against its digits
// (Lexer.g:333), so LIMIT -1 parses and LIMIT - 1 does not.
func (p *Parser) parseIntValue() (ast.ExprNode, error) {
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
	if p.cur.Type == tokMINUS {
		if next := p.peekNext(); next.Type == tokINTEGER && next.Loc == p.cur.End {
			start := p.cur.Loc
			p.advance()
			lit := &ast.IntegerLit{Val: "-" + p.cur.Str, Loc: ast.Loc{Start: start, End: p.cur.End}}
			p.advance()
			return lit, nil
		}
	}
	if p.cur.Type == tokINTEGER {
		lit := &ast.IntegerLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		return lit, nil
	}
	return nil, p.errorf("expected integer literal or bind marker, got %s", p.tokenDesc())
}
