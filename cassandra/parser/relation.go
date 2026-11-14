package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseRelationElement parses one condition of a WHERE clause
// (Parser.g:452 relationOrExpression): a relation, or expr(index, term) on
// a custom index.
func (p *Parser) parseRelationElement() (ast.ExprNode, error) {
	// 'expr(' is one token of Cassandra's lexer, spelled in lower case and
	// with no space before the parenthesis (Parser.g:459).
	if p.cur.Type == tokIDENT && p.cur.Str == "expr" {
		if next := p.peekNext(); next.Type == tokLPAREN && next.Loc == p.cur.End {
			return p.parseCustomIndexExpr()
		}
	}
	return p.parseRelation()
}

// parseCustomIndexExpr parses 'expr(' idxName ',' term ')' (Parser.g:457).
func (p *Parser) parseCustomIndexExpr() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // expr
	p.advance() // (
	if !isNameToken(p.cur) {
		return nil, p.errorf("expected index name, got %s", p.tokenDesc())
	}
	idx, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOMMA); err != nil {
		return nil, err
	}
	val, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.CustomIndexExpr{Index: idx, Value: val, Loc: p.makeLoc(start)}, nil
}

// isRelationType reports whether typ is a relationType (Parser.g:1765).
func isRelationType(typ int) bool {
	switch typ {
	case tokEQ, tokLT, tokLTE, tokGT, tokGTE, tokNE:
		return true
	}
	return false
}

// parseRelation parses a relation (Parser.g:1774).
//
//	cident relationType term
//	cident LIKE term
//	cident IS NOT NULL
//	TOKEN tupleOfIdentifiers relationType term
//	cident IN (inMarker | singleColumnInValues)
//	cident CONTAINS KEY? term
//	cident '[' term ']' relationType term
//	tupleOfIdentifiers (IN ... | relationType (tupleLiteral | markerForTuple))
//	'(' relation ')'
func (p *Parser) parseRelation() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.curLoc()

	if p.cur.Type == tokLPAREN {
		if p.startsTupleOfIdentifiers() {
			return p.parseMultiColumnRelation(start)
		}
		p.advance() // (
		rel, err := p.parseRelation()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return rel, nil
	}

	if isWord(p.cur, "token") && p.peekNext().Type == tokLPAREN {
		nameTok := p.cur
		p.advance() // TOKEN
		cols, err := p.parseTupleOfIdentifiers()
		if err != nil {
			return nil, err
		}
		fn := &ast.FunctionCall{
			Name: &ast.Identifier{Name: nameTok.Str, Loc: ast.Loc{Start: nameTok.Loc, End: nameTok.End}},
			Args: cols,
			Loc:  p.makeLoc(start),
		}
		op := p.cur.Str
		if !isRelationType(p.cur.Type) {
			return nil, p.errorf("expected comparison operator after TOKEN(...), got %s", p.tokenDesc())
		}
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ast.BinaryExpr{Left: fn, Op: op, Right: right, Loc: p.makeLoc(start)}, nil
	}

	if !isNameToken(p.cur) {
		return nil, p.errorf("expected column name in relation, got %s", p.tokenDesc())
	}
	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	switch {
	case isRelationType(p.cur.Type):
		op := p.cur.Str
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ast.BinaryExpr{Left: col, Op: op, Right: right, Loc: p.makeLoc(start)}, nil

	case isWord(p.cur, "like"):
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ast.BinaryExpr{Left: col, Op: "LIKE", Right: right, Loc: p.makeLoc(start)}, nil

	case p.cur.Type == tokIS:
		p.advance()
		if err := p.expectKeyword(tokNOT); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokNULL); err != nil {
			return nil, err
		}
		return &ast.IsNotNullExpr{Column: col, Loc: p.makeLoc(start)}, nil

	case p.cur.Type == tokIN:
		p.advance()
		values, marker, err := p.parseInValues()
		if err != nil {
			return nil, err
		}
		return &ast.InExpr{Column: col, Values: values, Marker: marker, Loc: p.makeLoc(start)}, nil

	case p.cur.Type == tokCONTAINS:
		p.advance()
		isKey := p.match(tokKEY)
		val, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ast.ContainsExpr{Column: col, Value: val, IsKey: isKey, Loc: p.makeLoc(start)}, nil

	case p.cur.Type == tokLBRACK:
		p.advance()
		key, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		left := &ast.IndexAccess{Collection: col, Index: key, Loc: p.makeLoc(start)}
		op := p.cur.Str
		if !isRelationType(p.cur.Type) {
			return nil, p.errorf("expected comparison operator, got %s", p.tokenDesc())
		}
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(start)}, nil
	}
	return nil, p.errorf("expected operator after column in WHERE clause, got %s", p.tokenDesc())
}

// parseInValues parses what follows IN for one column: a bind marker
// (Parser.g:1832 inMarker), or '(' (term (',' term)*)? ')'
// (Parser.g:1842 singleColumnInValues).
func (p *Parser) parseInValues() ([]ast.ExprNode, *ast.BindMarker, error) {
	if p.cur.Type == tokQMARK || p.cur.Type == tokCOLON {
		m, err := p.parseBindMarker()
		return nil, m, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, nil, err
	}
	var values []ast.ExprNode
	if p.cur.Type != tokRPAREN {
		var err error
		if values, err = p.parseTermList(); err != nil {
			return nil, nil, err
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, nil, err
	}
	return values, nil, nil
}

// startsTupleOfIdentifiers reports whether the '(' at p.cur opens a
// tupleOfIdentifiers, '(' cident (',' cident)* ')', rather than a relation
// in parentheses.
func (p *Parser) startsTupleOfIdentifiers() bool {
	s := p.mark()
	defer p.reset(s)
	p.advance() // (
	if !isNameToken(p.cur) {
		return false
	}
	p.advance()
	return p.cur.Type == tokCOMMA || p.cur.Type == tokRPAREN
}

// parseTupleOfIdentifiers parses '(' cident (',' cident)* ')'
// (Parser.g:1837).
func (p *Parser) parseTupleOfIdentifiers() ([]ast.ExprNode, error) {
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	var cols []ast.ExprNode
	for {
		if !isNameToken(p.cur) {
			return nil, p.errorf("expected column name, got %s", p.tokenDesc())
		}
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if !p.match(tokCOMMA) {
			break
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseMultiColumnRelation parses a relation on a tuple of columns
// (Parser.g:1786-1806):
//
//	(a, b) IN ()
//	(a, b) IN ?                       inMarkerForTuple
//	(a, b) IN ((1, 2), (3, 4))        tupleOfTupleLiterals
//	(a, b) IN (?, :m)                 tupleOfMarkersForTuples
//	(a, b) > (1, 2)                   tupleLiteral
//	(a, b) > ?                        markerForTuple
func (p *Parser) parseMultiColumnRelation(start int) (ast.ExprNode, error) {
	cols, err := p.parseTupleOfIdentifiers()
	if err != nil {
		return nil, err
	}

	if p.match(tokIN) {
		rel := &ast.TupleInExpr{Columns: cols}
		switch {
		case p.cur.Type == tokQMARK || p.cur.Type == tokCOLON:
			if rel.Marker, err = p.parseBindMarker(); err != nil {
				return nil, err
			}
		default:
			if _, err := p.expect(tokLPAREN); err != nil {
				return nil, err
			}
			switch {
			case p.cur.Type == tokRPAREN:
			case p.cur.Type == tokQMARK || p.cur.Type == tokCOLON:
				for {
					m, err := p.parseBindMarker()
					if err != nil {
						return nil, err
					}
					rel.Markers = append(rel.Markers, m)
					if !p.match(tokCOMMA) {
						break
					}
				}
			default:
				for {
					t, err := p.parseTupleLiteral()
					if err != nil {
						return nil, err
					}
					rel.Tuples = append(rel.Tuples, t)
					if !p.match(tokCOMMA) {
						break
					}
				}
			}
			if _, err := p.expect(tokRPAREN); err != nil {
				return nil, err
			}
		}
		rel.Loc = p.makeLoc(start)
		return rel, nil
	}

	op := p.cur.Str
	if !isRelationType(p.cur.Type) {
		return nil, p.errorf("expected operator or IN after tuple columns, got %s", p.tokenDesc())
	}
	p.advance()
	rel := &ast.TupleCompareExpr{Columns: cols, Op: op}
	if p.cur.Type == tokQMARK || p.cur.Type == tokCOLON {
		if rel.Marker, err = p.parseBindMarker(); err != nil {
			return nil, err
		}
	} else {
		t, err := p.parseTupleLiteral()
		if err != nil {
			return nil, err
		}
		rel.Values = t.Elements
	}
	rel.Loc = p.makeLoc(start)
	return rel, nil
}
