package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// Selectors are what SELECT returns (Parser.g:304-445): a column, a literal,
// a function call, arithmetic on them, element and field selection, a type
// hint, a tuple, list, set or map of selectors.

// parseSelectors parses '*' or selector (',' selector)* (Parser.g:304).
func (p *Parser) parseSelectors() ([]*ast.SelectElement, error) {
	if p.cur.Type == tokSTAR {
		start := p.curLoc()
		p.advance()
		star := &ast.StarExpr{Loc: p.makeLoc(start)}
		return []*ast.SelectElement{{Expr: star, Loc: p.makeLoc(start)}}, nil
	}
	var elements []*ast.SelectElement
	for {
		el, err := p.parseSelector()
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
		if !p.match(tokCOMMA) {
			return elements, nil
		}
	}
}

// parseSelector parses unaliasedSelector (AS noncol_ident)? (Parser.g:309).
func (p *Parser) parseSelector() (*ast.SelectElement, error) {
	start := p.curLoc()
	expr, err := p.parseUnaliasedSelector()
	if err != nil {
		return nil, err
	}
	alias, err := p.parseOptionalAlias()
	if err != nil {
		return nil, err
	}
	return &ast.SelectElement{Expr: expr, Alias: alias, Loc: p.makeLoc(start)}, nil
}

// parseUnaliasedSelector parses selectionMultiplication joined by + and -
// (Parser.g:314 unaliasedSelector, 318 selectionAddition).
func (p *Parser) parseUnaliasedSelector() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	left, err := p.parseSelectionMultiplication()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		op := p.cur.Str
		p.advance()
		right, err := p.parseSelectionMultiplication()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(left.GetLoc().Start)}
	}
	return left, nil
}

// parseSelectionMultiplication parses selectionGroup joined by *, / and %
// (Parser.g:325).
func (p *Parser) parseSelectionMultiplication() (ast.ExprNode, error) {
	left, err := p.parseSelectionGroup()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokSTAR || p.cur.Type == tokSLASH || p.cur.Type == tokPERCENT {
		op := p.cur.Str
		p.advance()
		right, err := p.parseSelectionGroup()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right, Loc: p.makeLoc(left.GetLoc().Start)}
	}
	return left, nil
}

// parseSelectionGroup parses - selectionGroup, or selectionGroupWithoutField
// followed by field and element selections (Parser.g:333-361).
func (p *Parser) parseSelectionGroup() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	if p.cur.Type == tokMINUS {
		start := p.curLoc()
		p.advance()
		operand, err := p.parseSelectionGroup()
		if err != nil {
			return nil, err
		}
		return negate(operand, start, p.prevEnd()), nil
	}
	expr, err := p.parseSelectionGroupWithoutField()
	if err != nil {
		return nil, err
	}
	return p.parseSelectorModifiers(expr)
}

// parseSelectorModifiers parses the field selections (.field) and element
// or slice selections ([term], [term..term]) after receiver
// (Parser.g:343 selectorModifier). They chain left to right, in a loop.
func (p *Parser) parseSelectorModifiers(receiver ast.ExprNode) (ast.ExprNode, error) {
	start := receiver.GetLoc().Start
	for {
		switch p.cur.Type {
		case tokDOT:
			p.advance()
			if !isUnreservedName(p.cur) {
				return nil, p.errorf("expected field name, got %s", p.tokenDesc())
			}
			field, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			receiver = &ast.DotAccess{Object: receiver, Field: field, Loc: p.makeLoc(start)}
		case tokLBRACK:
			p.advance()
			sel, err := p.parseCollectionSubSelection(receiver, start)
			if err != nil {
				return nil, err
			}
			receiver = sel
		default:
			return receiver, nil
		}
	}
}

// parseCollectionSubSelection parses what follows '[' in a selector, through
// the ']': term, term .. term?, or .. term (Parser.g:353).
func (p *Parser) parseCollectionSubSelection(receiver ast.ExprNode, start int) (ast.ExprNode, error) {
	var from, to ast.ExprNode
	slice := false
	var err error
	if p.match(tokRANGE) {
		slice = true
		if to, err = p.parseTerm(); err != nil {
			return nil, err
		}
	} else {
		if from, err = p.parseTerm(); err != nil {
			return nil, err
		}
		if p.match(tokRANGE) {
			slice = true
			if p.cur.Type != tokRBRACK {
				if to, err = p.parseTerm(); err != nil {
					return nil, err
				}
			}
		}
	}
	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}
	if slice {
		return &ast.SliceAccess{Collection: receiver, From: from, To: to, Loc: p.makeLoc(start)}, nil
	}
	return &ast.IndexAccess{Collection: receiver, Index: from, Loc: p.makeLoc(start)}, nil
}

// startsSelectionGroupWithoutField reports whether tok can start a
// selectionGroupWithoutField (Parser.g:364). FROM, which Cassandra
// reserves, never does: SELECT (a) FROM t selects a in parentheses.
func startsSelectionGroupWithoutField(tok, next Token) bool {
	if isSignedConstant(tok, next) {
		return true
	}
	switch tok.Type {
	case tokNULL, tokQMARK, tokCOLON, tokLPAREN, tokLBRACK, tokLBRACE:
		return true
	case tokFROM:
		return false
	}
	return isConstantToken(tok) || isIdentLike(tok.Type)
}

// parseSelectionGroupWithoutField parses a column, literal, function call,
// type hint, tuple or parenthesised selector, list, set or map
// (Parser.g:364).
func (p *Parser) parseSelectionGroupWithoutField() (ast.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	switch {
	case p.cur.Type == tokLPAREN:
		// (selectionTypeHint)=> (Parser.g:368): a type in parentheses
		// before a selector is a type hint; anything else in them is a
		// tuple or a parenthesised selector.
		if hint, ok, err := p.tryTypeHint(startsSelectionGroupWithoutField, p.parseSelectionGroupWithoutField); ok || err != nil {
			return hint, err
		}
		return p.parseSelectionTuple()
	case p.cur.Type == tokLBRACK:
		return p.parseSelectionList()
	case p.cur.Type == tokLBRACE:
		return p.parseSelectionMapOrSet()
	case p.cur.Type == tokQMARK || p.cur.Type == tokCOLON:
		return p.parseBindMarker()
	case p.cur.Type == tokNULL:
		tok := p.cur
		p.advance()
		return &ast.NullLit{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case isConstantToken(p.cur):
		return p.parseConstant()
	case isSignedConstant(p.cur, p.peekNext()):
		// A sign written against a number is part of its constant
		// (Lexer.g:337), so (int) -1 is a type hint on -1.
		start := p.curLoc()
		p.advance() // -
		operand, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		return negate(operand, start, p.prevEnd()), nil
	case isIdentLike(p.cur.Type):
		return p.parseSelectionNameOrFunction()
	}
	return nil, p.errorf("expected selector, got %s", p.tokenDesc())
}

// parseSelectionNameOrFunction parses a column name (sident) or a
// selectionFunction (Parser.g:418): COUNT(*), WRITETIME, MAXWRITETIME or TTL
// of a column, CAST, or a function of selectors.
func (p *Parser) parseSelectionNameOrFunction() (ast.ExprNode, error) {
	next := p.peekNext()
	if next.Type == tokLPAREN {
		switch {
		case isWord(p.cur, "count"):
			s := p.mark()
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			p.advance() // (
			if p.match(tokSTAR) {
				if _, err := p.expect(tokRPAREN); err != nil {
					return nil, err
				}
				return &ast.FunctionCall{Name: name, Star: true, Loc: p.makeLoc(name.Loc.Start)}, nil
			}
			p.reset(s)
		case isWord(p.cur, "writetime"), isWord(p.cur, "maxwritetime"), isWord(p.cur, "ttl"):
			return p.parseWritetimeOrTTL()
		case p.cur.Type == tokCAST:
			return p.parseCastWith(p.parseUnaliasedSelector)
		}
		if isFunctionNameToken(p.cur) {
			return p.parseSelectionFunction()
		}
	}
	if next.Type == tokDOT && isFunctionNameToken(p.cur) {
		// keyspace.function(...) or column.field.
		s := p.mark()
		p.advance() // keyspace
		p.advance() // .
		isFunc := isFunctionNameToken(p.cur) && p.peekNext().Type == tokLPAREN
		p.reset(s)
		if isFunc {
			return p.parseSelectionFunction()
		}
	}
	if !isNameToken(p.cur) {
		return nil, p.errorf("expected selector, got %s", p.tokenDesc())
	}
	return p.parseIdentifier()
}

// parseWritetimeOrTTL parses WRITETIME, MAXWRITETIME or TTL '(' sident
// selectorModifier ')' (Parser.g:420-422): the argument is a column, with
// field and element selections.
func (p *Parser) parseWritetimeOrTTL() (ast.ExprNode, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	if !isNameToken(p.cur) {
		return nil, p.errorf("expected column name, got %s", p.tokenDesc())
	}
	col, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	arg, err := p.parseSelectorModifiers(col)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.FunctionCall{Name: name, Args: []ast.ExprNode{arg}, Loc: p.makeLoc(name.Loc.Start)}, nil
}

// parseSelectionFunction parses functionName '(' (unaliasedSelector
// (',' unaliasedSelector)*)? ')' (Parser.g:423, 434).
func (p *Parser) parseSelectionFunction() (ast.ExprNode, error) {
	name, err := p.parseFunctionName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	var args []ast.ExprNode
	if p.cur.Type != tokRPAREN {
		if args, err = p.parseSelectorList(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.FunctionCall{Name: name, Args: args, Loc: p.makeLoc(name.Loc.Start)}, nil
}

// parseSelectorList parses unaliasedSelector (',' unaliasedSelector)*.
func (p *Parser) parseSelectorList() ([]ast.ExprNode, error) {
	var list []ast.ExprNode
	for {
		s, err := p.parseUnaliasedSelector()
		if err != nil {
			return nil, err
		}
		list = append(list, s)
		if !p.match(tokCOMMA) {
			return list, nil
		}
	}
}

// parseSelectionTuple parses '(' unaliasedSelector (',' unaliasedSelector)*
// ')' (Parser.g:402). One selector in parentheses is a TupleLit of one:
// Cassandra decides later whether it is a tuple.
func (p *Parser) parseSelectionTuple() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // (
	elems, err := p.parseSelectorList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseSelectionList parses '[' (unaliasedSelector (','
// unaliasedSelector)*)? ']' (Parser.g:379).
func (p *Parser) parseSelectionList() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // [
	var elems []ast.ExprNode
	if p.cur.Type != tokRBRACK {
		var err error
		if elems, err = p.parseSelectorList(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}
	return &ast.ListLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseSelectionMapOrSet parses '{' '}', a set of selectors, or a map of
// selectors to selectors, which is also how a UDT is written
// (Parser.g:385-400). {} is an empty MapLit.
func (p *Parser) parseSelectionMapOrSet() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // {
	if p.match(tokRBRACE) {
		return &ast.MapLit{Loc: p.makeLoc(start)}, nil
	}
	first, err := p.parseUnaliasedSelector()
	if err != nil {
		return nil, err
	}
	if p.match(tokCOLON) {
		v, err := p.parseUnaliasedSelector()
		if err != nil {
			return nil, err
		}
		keys, values := []ast.ExprNode{first}, []ast.ExprNode{v}
		for p.match(tokCOMMA) {
			k, err := p.parseUnaliasedSelector()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokCOLON); err != nil {
				return nil, err
			}
			v, err := p.parseUnaliasedSelector()
			if err != nil {
				return nil, err
			}
			keys, values = append(keys, k), append(values, v)
		}
		if _, err := p.expect(tokRBRACE); err != nil {
			return nil, err
		}
		return &ast.MapLit{Keys: keys, Values: values, Loc: p.makeLoc(start)}, nil
	}
	elems := []ast.ExprNode{first}
	for p.match(tokCOMMA) {
		e, err := p.parseUnaliasedSelector()
		if err != nil {
			return nil, err
		}
		elems = append(elems, e)
	}
	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	return &ast.SetLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}
