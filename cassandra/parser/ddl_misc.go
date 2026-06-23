package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateFunction() (*ast.CreateFunctionStmt, error) {
	start := p.curLoc()
	p.advance() // FUNCTION

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

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
	} else if p.cur.Type == tokRETURNS && p.peekNext().Type == tokNULL {
		stmt.ReturnMode = ast.ReturnNullOnNull
		p.advance() // RETURNS
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

	// RETURNS dataType
	if err := p.expectKeyword(tokRETURNS); err != nil {
		return nil, err
	}

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

	var argTypes []*ast.DataType
	if p.cur.Type == tokLPAREN {
		p.advance()
		if p.cur.Type != tokRPAREN {
			argTypes, err = p.parseTypeList()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	}

	return &ast.DropFunctionStmt{
		IfExists: ifExists,
		Name:     name,
		ArgTypes: argTypes,
		Loc:      p.makeLoc(start),
	}, nil
}

func (p *Parser) parseCreateAggregate() (*ast.CreateAggregateStmt, error) {
	start := p.curLoc()
	p.advance() // AGGREGATE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

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

	var finalfunc *ast.Identifier
	if p.cur.Type == tokFINALFUNC {
		p.advance()
		ff, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		finalfunc = ff
	}

	var initcond ast.ExprNode
	if p.cur.Type == tokINITCOND {
		p.advance()
		ic, err := p.parseInitCondDefinition()
		if err != nil {
			return nil, err
		}
		initcond = ic
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

	var argTypes []*ast.DataType
	if p.cur.Type == tokLPAREN {
		p.advance()
		if p.cur.Type != tokRPAREN {
			var err error
			argTypes, err = p.parseTypeList()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	}

	return &ast.DropAggregateStmt{
		IfExists: ifExists,
		Name:     name,
		ArgTypes: argTypes,
		Loc:      p.makeLoc(start),
	}, nil
}

// parseTypeList parses a comma-separated list of data types.
func (p *Parser) parseTypeList() ([]*ast.DataType, error) {
	var types []*ast.DataType
	for {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		types = append(types, dt)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return types, nil
}

func (p *Parser) parseCreateTrigger() (*ast.CreateTriggerStmt, error) {
	start := p.curLoc()
	p.advance() // TRIGGER

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

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
		Table:       table,
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
