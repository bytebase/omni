package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateFunction() (*ast.CreateFunctionStmt, error) {
	start := p.curLoc()
	p.advance() // FUNCTION

	ifNotExists := p.parseIfNotExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateFunctionStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	// Parameter list
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	if p.cur.Type != tokRPAREN {
		for {
			pStart := p.curLoc()
			pName, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			pType, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.Params = append(stmt.Params, &ast.FunctionParam{
				Name: pName,
				Type: pType,
				Loc:  p.makeLoc(pStart),
			})
			if !p.match(tokCOMMA) {
				break
			}
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	// Return mode: CALLED ON NULL INPUT | RETURNS NULL ON NULL INPUT
	if p.cur.Type == tokCALLED {
		stmt.ReturnMode = ast.ReturnCalledOnNull
		p.advance() // CALLED
		if err := p.expectKeyword(tokON); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokNULL); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokINPUT); err != nil {
			return nil, err
		}
	} else if p.cur.Type == tokRETURNS {
		p.advance() // RETURNS
		if p.cur.Type == tokNULL {
			stmt.ReturnMode = ast.ReturnNullOnNull
			p.advance() // NULL
			if err := p.expectKeyword(tokON); err != nil {
				return nil, err
			}
			if err := p.expectKeyword(tokNULL); err != nil {
				return nil, err
			}
			if err := p.expectKeyword(tokINPUT); err != nil {
				return nil, err
			}
		}
		// RETURNS <type>
		if err := p.expectKeyword(tokRETURNS); err != nil {
			// If we already consumed RETURNS for RETURNS NULL, expect another RETURNS for return type.
			// But actually in the grammar: returnMode RETURNS dataType
			// So if we consumed RETURNS for "RETURNS NULL ON NULL INPUT", we need another RETURNS.
			return nil, err
		}
	}

	// Actually, let me re-read the grammar:
	// returnMode: (CALLED | RETURNS NULL) ON NULL INPUT
	// Then: RETURNS dataType LANGUAGE language AS codeBlock
	// So RETURNS is used twice: once in returnMode and once for return type.
	// Let me fix the approach.

	// If we haven't seen return mode yet (came here via else path),
	// the return type RETURNS should follow.
	// But the above code already handled returnMode, so now we expect:
	// RETURNS dataType LANGUAGE language AS codeBlock
	// However, if returnMode was "RETURNS NULL ON NULL INPUT", we already consumed RETURNS once.
	// The grammar expects another RETURNS keyword for the return type.

	retType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	stmt.ReturnType = retType

	if err := p.expectKeyword(tokLANGUAGE); err != nil {
		return nil, err
	}
	lang, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Language = lang

	if err := p.expectKeyword(tokAS); err != nil {
		return nil, err
	}

	// Body: code block or string literal
	switch p.cur.Type {
	case tokCODEBLOCK:
		stmt.Body = &ast.CodeBlock{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
	case tokSTRING:
		stmt.Body = &ast.StringLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
	default:
		return nil, p.errorf("expected code block or string literal for function body, got %s", p.tokenDesc())
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropFunction() (*ast.DropFunctionStmt, error) {
	start := p.curLoc()
	p.advance() // FUNCTION

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropFunctionStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

func (p *Parser) parseCreateAggregate() (*ast.CreateAggregateStmt, error) {
	start := p.curLoc()
	p.advance() // AGGREGATE

	ifNotExists := p.parseIfNotExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	paramType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokSFUNC); err != nil {
		return nil, err
	}
	sfunc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokSTYPE); err != nil {
		return nil, err
	}
	stype, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokFINALFUNC); err != nil {
		return nil, err
	}
	finalfunc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokINITCOND); err != nil {
		return nil, err
	}
	initcond, err := p.parseInitCondDefinition()
	if err != nil {
		return nil, err
	}

	return &ast.CreateAggregateStmt{
		IfNotExists: ifNotExists,
		Name:        name,
		ParamType:   paramType,
		SFunc:       sfunc,
		SType:       stype,
		FinalFunc:   finalfunc,
		InitCond:    initcond,
		Loc:         p.makeLoc(start),
	}, nil
}

func (p *Parser) parseInitCondDefinition() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokLBRACE:
		return p.parseMapOrSetLiteral()
	case tokLPAREN:
		return p.parseInitCondListOrNested()
	default:
		return p.parseConstant()
	}
}

func (p *Parser) parseInitCondListOrNested() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // (
	var elems []ast.ExprNode
	for p.cur.Type != tokRPAREN && p.cur.Type != tokEOF {
		elem, err := p.parseInitCondDefinition()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
		if !p.match(tokCOMMA) {
			break
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

func (p *Parser) parseDropAggregate() (*ast.DropAggregateStmt, error) {
	start := p.curLoc()
	p.advance() // AGGREGATE

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropAggregateStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

func (p *Parser) parseCreateTrigger() (*ast.CreateTriggerStmt, error) {
	start := p.curLoc()
	p.advance() // TRIGGER

	ifNotExists := p.parseIfNotExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokUSING); err != nil {
		return nil, err
	}

	usingClass, err := p.parseConstant()
	if err != nil {
		return nil, err
	}

	return &ast.CreateTriggerStmt{
		IfNotExists: ifNotExists,
		Name:        name,
		UsingClass:  usingClass,
		Loc:         p.makeLoc(start),
	}, nil
}

func (p *Parser) parseDropTrigger() (*ast.DropTriggerStmt, error) {
	start := p.curLoc()
	p.advance() // TRIGGER

	ifExists := p.parseIfExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropTriggerStmt{
		IfExists: ifExists,
		Name:     name,
		Table:    table,
		Loc:      p.makeLoc(start),
	}, nil
}
