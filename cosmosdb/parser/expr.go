package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

// Precedence levels for the Pratt parser.
const (
	precTernary    = 1  // ? :
	precCoalesce   = 2  // ??
	precOr         = 3  // OR
	precAnd        = 4  // AND
	precInBetween  = 5  // IN, BETWEEN, LIKE
	precComparison = 6  // = != <> < <= > >=
	precConcat     = 7  // ||
	precBitOr      = 8  // |
	precBitXor     = 9  // ^
	precBitAnd     = 10 // &
	precShift      = 11 // << >> >>>
	precAdditive   = 12 // + -
	precMult       = 13 // * / %
	precUnary      = 14 // NOT ~ + -
	precPostfix    = 15 // . []
)

// parseExpr is the main Pratt expression parser.
func (p *Parser) parseExpr(minPrec int) (nodes.ExprNode, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}

	for {
		prec, ok := p.infixPrecedence()
		if !ok || prec < minPrec {
			break
		}
		left, err = p.parseInfixExpr(left, prec)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

// parsePrefixExpr handles prefix/primary expressions.
func (p *Parser) parsePrefixExpr() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokTRUE:
		loc := p.cur.Loc
		p.advance()
		return &nodes.BoolLit{Val: true, Loc: nodes.Loc{Start: loc, End: loc + 4}}, nil
	case tokFALSE:
		loc := p.cur.Loc
		p.advance()
		return &nodes.BoolLit{Val: false, Loc: nodes.Loc{Start: loc, End: loc + 5}}, nil
	case tokNULL:
		loc := p.cur.Loc
		p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: loc, End: loc + 4}}, nil
	case tokUNDEFINED:
		loc := p.cur.Loc
		p.advance()
		return &nodes.UndefinedLit{Loc: nodes.Loc{Start: loc, End: loc + 9}}, nil
	case tokINFINITY:
		loc := p.cur.Loc
		p.advance()
		return &nodes.InfinityLit{Loc: nodes.Loc{Start: loc, End: loc + 8}}, nil
	case tokNAN:
		loc := p.cur.Loc
		p.advance()
		return &nodes.NanLit{Loc: nodes.Loc{Start: loc, End: loc + 3}}, nil

	case tokICONST, tokFCONST, tokHCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.NumberLit{Val: val, Loc: nodes.Loc{Start: loc, End: loc + len(val)}}, nil

	case tokSCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.StringLit{Val: val, Loc: nodes.Loc{Start: loc, End: p.prev.Loc + len(p.prev.Str)}}, nil

	case tokDCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.StringLit{Val: val, Loc: nodes.Loc{Start: loc, End: p.prev.Loc + len(p.prev.Str)}}, nil

	case tokPARAM:
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &nodes.ParamRef{Name: name, Loc: nodes.Loc{Start: loc, End: p.prev.Loc + len(p.prev.Str)}}, nil

	// Unary operators.
	case tokNOT:
		loc := p.cur.Loc
		p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.UnaryExpr{Op: "NOT", Operand: operand, Loc: nodes.Loc{Start: loc, End: endLoc}}, nil

	case tokBITNOT:
		loc := p.cur.Loc
		p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.UnaryExpr{Op: "~", Operand: operand, Loc: nodes.Loc{Start: loc, End: endLoc}}, nil

	case tokPLUS:
		loc := p.cur.Loc
		p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.UnaryExpr{Op: "+", Operand: operand, Loc: nodes.Loc{Start: loc, End: endLoc}}, nil

	case tokMINUS:
		loc := p.cur.Loc
		p.advance()
		operand, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.UnaryExpr{Op: "-", Operand: operand, Loc: nodes.Loc{Start: loc, End: endLoc}}, nil

	// EXISTS(SELECT ...)
	case tokEXISTS:
		return p.parseExistsExpr()

	// ARRAY(SELECT ...)
	case tokARRAY:
		// Check if it's ARRAY(...) - function form - vs identifier use.
		if p.peekNext().Type == tokLPAREN {
			return p.parseArrayExpr()
		}
		// Otherwise treat as identifier.
		return p.parseIdentExprOrFuncCall()

	// UDF.name(args)
	case tokUDF:
		return p.parseUDFCall()

	// Parenthesized expression or subquery.
	case tokLPAREN:
		return p.parseParenExpr()

	// Array literal [...]
	case tokLBRACK:
		return p.parseCreateArrayExpr()

	// Object literal {...}
	case tokLBRACE:
		return p.parseCreateObjectExpr()

	default:
		// Identifier / function call / keyword-as-identifier.
		if isIdentLike(p.cur.Type) || p.cur.Type == tokIDENT {
			return p.parseIdentExprOrFuncCall()
		}
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q in expression", p.cur.Str),
			Pos:     p.cur.Loc,
		}
	}
}

// infixPrecedence returns the precedence and whether the current token is an infix operator.
func (p *Parser) infixPrecedence() (int, bool) {
	switch p.cur.Type {
	case tokQUESTION:
		return precTernary, true
	case tokCOALESCE:
		return precCoalesce, true
	case tokOR:
		return precOr, true
	case tokAND:
		return precAnd, true
	case tokIN:
		return precInBetween, true
	case tokBETWEEN:
		return precInBetween, true
	case tokLIKE:
		return precInBetween, true
	case tokNOT:
		// NOT IN, NOT BETWEEN, NOT LIKE.
		next := p.peekNext()
		if next.Type == tokIN || next.Type == tokBETWEEN || next.Type == tokLIKE {
			return precInBetween, true
		}
		return 0, false
	case tokEQ, tokNE, tokNE2, tokLT, tokLE, tokGT, tokGE:
		return precComparison, true
	case tokCONCAT:
		return precConcat, true
	case tokBITOR:
		return precBitOr, true
	case tokBITXOR:
		return precBitXor, true
	case tokBITAND:
		return precBitAnd, true
	case tokLSHIFT, tokRSHIFT, tokURSHIFT:
		return precShift, true
	case tokPLUS, tokMINUS:
		return precAdditive, true
	case tokSTAR, tokDIV, tokMOD:
		return precMult, true
	case tokDOT:
		return precPostfix, true
	case tokLBRACK:
		return precPostfix, true
	default:
		return 0, false
	}
}

// parseInfixExpr parses the infix portion of an expression.
func (p *Parser) parseInfixExpr(left nodes.ExprNode, prec int) (nodes.ExprNode, error) {
	startLoc := locStart(left)

	switch p.cur.Type {
	case tokQUESTION:
		// Ternary: cond ? then : else
		p.advance()
		thenExpr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokCOLON); err != nil {
			return nil, err
		}
		elseExpr, err := p.parseExpr(precTernary)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.TernaryExpr{
			Cond: left, Then: thenExpr, Else: elseExpr,
			Loc: nodes.Loc{Start: startLoc, End: endLoc},
		}, nil

	case tokNOT:
		// NOT IN / NOT BETWEEN / NOT LIKE
		p.advance() // consume NOT
		switch p.cur.Type {
		case tokIN:
			return p.parseInExpr(left, true, startLoc)
		case tokBETWEEN:
			return p.parseBetweenExpr(left, true, startLoc)
		case tokLIKE:
			return p.parseLikeExpr(left, true, startLoc)
		default:
			return nil, &ParseError{
				Message: fmt.Sprintf("expected IN, BETWEEN, or LIKE after NOT, got %q", p.cur.Str),
				Pos:     p.cur.Loc,
			}
		}

	case tokIN:
		return p.parseInExpr(left, false, startLoc)

	case tokBETWEEN:
		return p.parseBetweenExpr(left, false, startLoc)

	case tokLIKE:
		return p.parseLikeExpr(left, false, startLoc)

	case tokDOT:
		p.advance()
		prop, _, err := p.parsePropertyName()
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.DotAccessExpr{
			Expr: left, Property: prop,
			Loc: nodes.Loc{Start: startLoc, End: endLoc},
		}, nil

	case tokLBRACK:
		p.advance()
		index, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		rbracket, err := p.expect(tokRBRACK)
		if err != nil {
			return nil, err
		}
		return &nodes.BracketAccessExpr{
			Expr: left, Index: index,
			Loc: nodes.Loc{Start: startLoc, End: rbracket.Loc + 1},
		}, nil

	default:
		// Binary operator.
		op := p.cur.Str
		tokType := p.cur.Type
		p.advance()

		// Normalize operator string.
		switch tokType {
		case tokAND:
			op = "AND"
		case tokOR:
			op = "OR"
		case tokCONCAT:
			op = "||"
		case tokCOALESCE:
			op = "??"
		}

		// Right associativity for coalesce; left-assoc for everything else.
		nextPrec := prec + 1
		if tokType == tokCOALESCE {
			nextPrec = prec
		}

		right, err := p.parseExpr(nextPrec)
		if err != nil {
			return nil, err
		}
		endLoc := p.prev.Loc + len(p.prev.Str)
		return &nodes.BinaryExpr{
			Op: op, Left: left, Right: right,
			Loc: nodes.Loc{Start: startLoc, End: endLoc},
		}, nil
	}
}

// parseInExpr parses: [NOT] IN (expr, ...)
func (p *Parser) parseInExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume IN
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	var list []nodes.ExprNode
	if p.cur.Type != tokRPAREN {
		for {
			expr, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			list = append(list, expr)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}

	rparen, err := p.expect(tokRPAREN)
	if err != nil {
		return nil, err
	}

	return &nodes.InExpr{
		Expr: left, List: list, Not: not,
		Loc: nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
	}, nil
}

// parseBetweenExpr parses: [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume BETWEEN

	// Parse low bound with precComparison to avoid consuming AND.
	low, err := p.parseExpr(precComparison)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokAND); err != nil {
		return nil, err
	}

	high, err := p.parseExpr(precComparison)
	if err != nil {
		return nil, err
	}

	endLoc := p.prev.Loc + len(p.prev.Str)
	return &nodes.BetweenExpr{
		Expr: left, Low: low, High: high, Not: not,
		Loc: nodes.Loc{Start: startLoc, End: endLoc},
	}, nil
}

// parseLikeExpr parses: [NOT] LIKE pattern [ESCAPE escape].
func (p *Parser) parseLikeExpr(left nodes.ExprNode, not bool, startLoc int) (nodes.ExprNode, error) {
	p.advance() // consume LIKE

	pattern, err := p.parseExpr(precComparison)
	if err != nil {
		return nil, err
	}

	var escape nodes.ExprNode
	if p.cur.Type == tokESCAPE {
		p.advance()
		escape, err = p.parseExpr(precComparison)
		if err != nil {
			return nil, err
		}
	}

	endLoc := p.prev.Loc + len(p.prev.Str)
	return &nodes.LikeExpr{
		Expr: left, Pattern: pattern, Escape: escape, Not: not,
		Loc: nodes.Loc{Start: startLoc, End: endLoc},
	}, nil
}

// parseExistsExpr parses EXISTS(SELECT ...).
func (p *Parser) parseExistsExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume EXISTS
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	sub, err := p.parseSelect()
	if err != nil {
		return nil, err
	}
	rparen, err := p.expect(tokRPAREN)
	if err != nil {
		return nil, err
	}
	return &nodes.ExistsExpr{
		Select: sub,
		Loc:    nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
	}, nil
}

// parseArrayExpr parses ARRAY(SELECT ...).
func (p *Parser) parseArrayExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume ARRAY
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	sub, err := p.parseSelect()
	if err != nil {
		return nil, err
	}
	rparen, err := p.expect(tokRPAREN)
	if err != nil {
		return nil, err
	}
	return &nodes.ArrayExpr{
		Select: sub,
		Loc:    nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
	}, nil
}

// parseParenExpr parses a parenthesized expression or subquery.
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume (

	// Check if it's a subquery.
	if p.cur.Type == tokSELECT {
		sub, err := p.parseSelect()
		if err != nil {
			return nil, err
		}
		rparen, err := p.expect(tokRPAREN)
		if err != nil {
			return nil, err
		}
		return &nodes.SubLink{
			Select: sub,
			Loc:    nodes.Loc{Start: startLoc, End: rparen.Loc + 1},
		}, nil
	}

	// Regular parenthesized expression.
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseCreateArrayExpr parses [...].
func (p *Parser) parseCreateArrayExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume [

	var elements []nodes.ExprNode
	if p.cur.Type != tokRBRACK {
		for {
			elem, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}

	rbracket, err := p.expect(tokRBRACK)
	if err != nil {
		return nil, err
	}

	return &nodes.CreateArrayExpr{
		Elements: elements,
		Loc:      nodes.Loc{Start: startLoc, End: rbracket.Loc + 1},
	}, nil
}

// parseCreateObjectExpr parses {...}.
func (p *Parser) parseCreateObjectExpr() (nodes.ExprNode, error) {
	startLoc := p.cur.Loc
	p.advance() // consume {

	var fields []*nodes.ObjectFieldPair
	if p.cur.Type != tokRBRACE {
		for {
			field, err := p.parseObjectFieldPair()
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}

	rbrace, err := p.expect(tokRBRACE)
	if err != nil {
		return nil, err
	}

	return &nodes.CreateObjectExpr{
		Fields: fields,
		Loc:    nodes.Loc{Start: startLoc, End: rbrace.Loc + 1},
	}, nil
}

// parseObjectFieldPair parses key: value inside an object literal.
func (p *Parser) parseObjectFieldPair() (*nodes.ObjectFieldPair, error) {
	startLoc := p.pos()

	// Key can be a string literal or identifier/property name.
	var key string
	switch p.cur.Type {
	case tokSCONST, tokDCONST:
		key = p.cur.Str
		p.advance()
	default:
		var err error
		key, _, err = p.parsePropertyName()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}

	value, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	endLoc := p.prev.Loc + len(p.prev.Str)
	return &nodes.ObjectFieldPair{
		Key: key, Value: value,
		Loc: nodes.Loc{Start: startLoc, End: endLoc},
	}, nil
}

// parseIdentExprOrFuncCall parses an identifier that may be followed by ( for a function call.
func (p *Parser) parseIdentExprOrFuncCall() (nodes.ExprNode, error) {
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for function call.
	if p.cur.Type == tokLPAREN {
		return p.parseFuncCall(name, loc)
	}

	return &nodes.ColumnRef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: loc + len(name)},
	}, nil
}

// locStart extracts the start position from any ExprNode.
func locStart(e nodes.ExprNode) int {
	switch n := e.(type) {
	case *nodes.ColumnRef:
		return n.Loc.Start
	case *nodes.DotAccessExpr:
		return n.Loc.Start
	case *nodes.BracketAccessExpr:
		return n.Loc.Start
	case *nodes.BinaryExpr:
		return n.Loc.Start
	case *nodes.UnaryExpr:
		return n.Loc.Start
	case *nodes.TernaryExpr:
		return n.Loc.Start
	case *nodes.InExpr:
		return n.Loc.Start
	case *nodes.BetweenExpr:
		return n.Loc.Start
	case *nodes.LikeExpr:
		return n.Loc.Start
	case *nodes.FuncCall:
		return n.Loc.Start
	case *nodes.UDFCall:
		return n.Loc.Start
	case *nodes.ExistsExpr:
		return n.Loc.Start
	case *nodes.ArrayExpr:
		return n.Loc.Start
	case *nodes.CreateArrayExpr:
		return n.Loc.Start
	case *nodes.CreateObjectExpr:
		return n.Loc.Start
	case *nodes.ParamRef:
		return n.Loc.Start
	case *nodes.SubLink:
		return n.Loc.Start
	case *nodes.StringLit:
		return n.Loc.Start
	case *nodes.NumberLit:
		return n.Loc.Start
	case *nodes.BoolLit:
		return n.Loc.Start
	case *nodes.NullLit:
		return n.Loc.Start
	case *nodes.UndefinedLit:
		return n.Loc.Start
	case *nodes.InfinityLit:
		return n.Loc.Start
	case *nodes.NanLit:
		return n.Loc.Start
	default:
		return -1
	}
}
