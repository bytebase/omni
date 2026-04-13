package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// Binding power constants for Pratt parsing (low to high).
const (
	bpNone       = 0
	bpOr         = 10 // OR
	bpAnd        = 20 // AND
	bpNot        = 30 // NOT (prefix)
	bpComparison = 40 // =, <>, <, >, <=, >=, != + IS, IN, BETWEEN, LIKE
	bpAdd        = 50 // +, -, ||
	bpMul        = 60 // *, /, %
	bpUnary      = 70 // unary -, +
	bpPostfix    = 80 // ::, COLLATE, OVER
	bpAccess     = 90 // :, [], .
)

// ---------------------------------------------------------------------------
// Core Pratt loop
// ---------------------------------------------------------------------------

// parseExpr is the entry point for expression parsing. It parses an
// expression using Pratt parsing / precedence climbing.
func (p *Parser) parseExpr() (ast.Node, error) {
	return p.parseExprPrec(bpNone + 1)
}

// parseExprPrec parses an expression with the given minimum binding power.
func (p *Parser) parseExprPrec(minBP int) (ast.Node, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}

	for {
		bp, op, special, ok := p.infixBindingPower()
		if !ok || bp < minBP {
			break
		}

		if special != "" {
			switch special {
			case "IS":
				left, err = p.parseIsExpr(left)
			case "BETWEEN":
				left, err = p.parseBetweenExpr(left)
			case "NOT_BETWEEN":
				left, err = p.parseBetweenExpr(left)
			case "IN":
				left, err = p.parseInExpr(left)
			case "NOT_IN":
				left, err = p.parseInExpr(left)
			case "LIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpLike, false)
			case "NOT_LIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpLike, true)
			case "ILIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpILike, false)
			case "NOT_ILIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpILike, true)
			case "RLIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpRLike, false)
			case "NOT_RLIKE":
				left, err = p.parseLikeExpr(left, ast.LikeOpRLike, true)
			case "REGEXP":
				left, err = p.parseLikeExpr(left, ast.LikeOpRegexp, false)
			case "NOT_REGEXP":
				left, err = p.parseLikeExpr(left, ast.LikeOpRegexp, true)
			case "DOUBLE_COLON":
				left, err = p.parseCastPostfix(left)
			case "COLLATE":
				left, err = p.parseCollatePostfix(left)
			case "OVER":
				left, err = p.parseOverPostfix(left)
			case "COLON_ACCESS":
				left, err = p.parseColonAccess(left)
			case "BRACKET_ACCESS":
				left, err = p.parseBracketAccess(left)
			case "DOT_ACCESS":
				left, err = p.parseDotAccess(left)
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		// Regular binary operator: consume the operator, parse right side.
		startLoc := p.cur.Loc
		p.advance()
		right, err := p.parseExprPrec(bp + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: startLoc.Start, End: ast.NodeLoc(right).End},
		}
	}

	return left, nil
}

// infixBindingPower returns (bp, op, special, ok) for the current token.
// If special is non-empty, the caller must dispatch to the named handler
// instead of building a BinaryExpr.
func (p *Parser) infixBindingPower() (int, ast.BinaryOp, string, bool) {
	switch p.cur.Type {
	// Logical
	case kwOR:
		return bpOr, ast.BinOr, "", true
	case kwAND:
		return bpAnd, ast.BinAnd, "", true

	// NOT as infix prefix for BETWEEN, IN, LIKE, ILIKE, RLIKE, REGEXP
	case kwNOT:
		next := p.peekNext()
		switch next.Type {
		case kwBETWEEN:
			return bpComparison, 0, "NOT_BETWEEN", true
		case kwIN:
			return bpComparison, 0, "NOT_IN", true
		case kwLIKE:
			return bpComparison, 0, "NOT_LIKE", true
		case kwILIKE:
			return bpComparison, 0, "NOT_ILIKE", true
		case kwRLIKE:
			return bpComparison, 0, "NOT_RLIKE", true
		default:
			// Check for REGEXP via tokIdent
			if next.Type == tokIdent && strings.ToUpper(next.Str) == "REGEXP" {
				return bpComparison, 0, "NOT_REGEXP", true
			}
		}
		return 0, 0, "", false

	// Comparison operators
	case '=':
		return bpComparison, ast.BinEq, "", true
	case tokNotEq:
		return bpComparison, ast.BinNe, "", true
	case '<':
		return bpComparison, ast.BinLt, "", true
	case '>':
		return bpComparison, ast.BinGt, "", true
	case tokLessEq:
		return bpComparison, ast.BinLe, "", true
	case tokGreaterEq:
		return bpComparison, ast.BinGe, "", true
	case kwIS:
		return bpComparison, 0, "IS", true
	case kwBETWEEN:
		return bpComparison, 0, "BETWEEN", true
	case kwIN:
		return bpComparison, 0, "IN", true
	case kwLIKE:
		return bpComparison, 0, "LIKE", true
	case kwILIKE:
		return bpComparison, 0, "ILIKE", true
	case kwRLIKE:
		return bpComparison, 0, "RLIKE", true

	// Arithmetic
	case '+':
		return bpAdd, ast.BinAdd, "", true
	case '-':
		return bpAdd, ast.BinSub, "", true
	case tokConcat:
		return bpAdd, ast.BinConcat, "", true
	case '*':
		return bpMul, ast.BinMul, "", true
	case '/':
		return bpMul, ast.BinDiv, "", true
	case '%':
		return bpMul, ast.BinMod, "", true

	// Postfix: :: cast
	case tokDoubleColon:
		return bpPostfix, 0, "DOUBLE_COLON", true

	// Postfix: COLLATE
	case kwCOLLATE:
		return bpPostfix, 0, "COLLATE", true

	// Postfix: OVER
	case kwOVER:
		return bpPostfix, 0, "OVER", true

	// Access: : (JSON path)
	case ':':
		return bpAccess, 0, "COLON_ACCESS", true

	// Access: [ (array subscript)
	case '[':
		return bpAccess, 0, "BRACKET_ACCESS", true

	// Access: . (dot path) — only for semi-structured chaining after :/[]
	case '.':
		// Dot access is tricky: we only want to treat . as an access operator
		// when the left side is an AccessExpr (semi-structured chaining).
		// Otherwise . is not an infix expression operator (it's handled in
		// identifier parsing for column refs).
		return bpAccess, 0, "DOT_ACCESS", true
	}

	// Check for REGEXP via tokIdent (not a keyword in F2)
	if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "REGEXP" {
		return bpComparison, 0, "REGEXP", true
	}

	return 0, 0, "", false
}

// ---------------------------------------------------------------------------
// Prefix dispatch
// ---------------------------------------------------------------------------

// parsePrefixExpr handles prefix unary operators (-, +, NOT) and delegates
// everything else to parsePrimaryExpr.
func (p *Parser) parsePrefixExpr() (ast.Node, error) {
	switch p.cur.Type {
	case '-':
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryMinus,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case '+':
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryPlus,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwNOT:
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpNot)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryNot,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwEXISTS:
		existsTok := p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var query ast.Node
		var err error
		if p.cur.Type == kwWITH {
			query, err = p.parseWithQueryExpr()
		} else {
			query, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.ExistsExpr{Query: query, Loc: ast.Loc{Start: existsTok.Loc.Start, End: closeTok.Loc.End}}, nil

	default:
		return p.parsePrimaryExpr()
	}
}

// ---------------------------------------------------------------------------
// Primary expression dispatch
// ---------------------------------------------------------------------------

// parsePrimaryExpr parses atomic expressions: literals, identifiers,
// function calls, parenthesized expressions, CASE, CAST, IFF, array/JSON
// literals, and lambda expressions.
func (p *Parser) parsePrimaryExpr() (ast.Node, error) {
	switch p.cur.Type {
	case tokInt:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitInt,
			Value: tok.Str,
			Ival:  tok.Ival,
			Loc:   tok.Loc,
		}, nil

	case tokFloat, tokReal:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitFloat,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokString:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitString,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case kwTRUE:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBool,
			Value: tok.Str,
			Bval:  true,
			Loc:   tok.Loc,
		}, nil

	case kwFALSE:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBool,
			Value: tok.Str,
			Bval:  false,
			Loc:   tok.Loc,
		}, nil

	case kwNULL:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitNull,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case '(':
		return p.parseParenExpr()

	case kwCASE:
		return p.parseCaseExpr()

	case kwCAST:
		return p.parseCastExpr()

	case kwTRY_CAST:
		return p.parseTryCastExpr()

	case kwIFF:
		return p.parseIffExpr()

	case '[':
		return p.parseArrayLiteral()

	case '{':
		return p.parseJsonLiteral()

	case '*':
		tok := p.advance()
		return &ast.StarExpr{Loc: tok.Loc}, nil

	default:
		// Lambda detection: single ident followed by ->
		if p.isIdentToken() && p.peekNext().Type == tokArrow {
			return p.parseLambdaExpr()
		}

		// Identifier-based: column ref or function call
		if p.isIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// isIdentToken reports whether the current token is usable as an identifier
// in expression context: tokIdent, tokQuotedIdent, or a non-reserved keyword.
func (p *Parser) isIdentToken() bool {
	switch p.cur.Type {
	case tokIdent, tokQuotedIdent:
		return true
	}
	return p.cur.Type >= 700 && !keywordReserved[p.cur.Type]
}

// ---------------------------------------------------------------------------
// Primary expression helpers
// ---------------------------------------------------------------------------

// parseIdentExpr parses an identifier that could be a column ref, function
// call, or a multi-part dotted name followed by function call/star.
func (p *Parser) parseIdentExpr() (ast.Node, error) {
	first, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	// Check for function call: name(
	if p.cur.Type == '(' {
		objName := ast.ObjectName{Name: first, Loc: first.Loc}
		return p.parseFuncCall(objName, first.Loc)
	}

	// Check for qualified name: name.name or name.name.name
	if p.cur.Type == '.' {
		// Could be schema.table.column or table.column or table.*
		// But be careful: dot access for semi-structured data is handled in
		// the infix loop (DOT_ACCESS). Here we only handle standard SQL
		// dotted column references.
		return p.parseDottedIdentOrFunc(first)
	}

	// Simple column ref (1 part).
	return &ast.ColumnRef{
		Parts: []ast.Ident{first},
		Loc:   first.Loc,
	}, nil
}

// parseDottedIdentOrFunc handles dotted names after the first ident.
// Returns a ColumnRef, StarExpr, or FuncCallExpr depending on what follows.
func (p *Parser) parseDottedIdentOrFunc(first ast.Ident) (ast.Node, error) {
	parts := []ast.Ident{first}

	for p.cur.Type == '.' {
		p.advance() // consume '.'

		// table.*
		if p.cur.Type == '*' {
			starTok := p.advance()
			var qualifier *ast.ObjectName
			switch len(parts) {
			case 1:
				qualifier = &ast.ObjectName{
					Name: parts[0],
					Loc:  parts[0].Loc,
				}
			case 2:
				qualifier = &ast.ObjectName{
					Schema: parts[0],
					Name:   parts[1],
					Loc:    ast.Loc{Start: parts[0].Loc.Start, End: parts[1].Loc.End},
				}
			case 3:
				qualifier = &ast.ObjectName{
					Database: parts[0],
					Schema:   parts[1],
					Name:     parts[2],
					Loc:      ast.Loc{Start: parts[0].Loc.Start, End: parts[2].Loc.End},
				}
			}
			return &ast.StarExpr{
				Qualifier: qualifier,
				Loc:       ast.Loc{Start: first.Loc.Start, End: starTok.Loc.End},
			}, nil
		}

		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		parts = append(parts, ident)

		// Check for function call: schema.func(
		if p.cur.Type == '(' {
			var objName ast.ObjectName
			switch len(parts) {
			case 2:
				objName = ast.ObjectName{
					Schema: parts[0],
					Name:   parts[1],
					Loc:    ast.Loc{Start: parts[0].Loc.Start, End: parts[1].Loc.End},
				}
			case 3:
				objName = ast.ObjectName{
					Database: parts[0],
					Schema:   parts[1],
					Name:     parts[2],
					Loc:      ast.Loc{Start: parts[0].Loc.Start, End: parts[2].Loc.End},
				}
			default:
				objName = ast.ObjectName{
					Name: parts[len(parts)-1],
					Loc:  parts[len(parts)-1].Loc,
				}
			}
			return p.parseFuncCall(objName, ast.Loc{Start: first.Loc.Start})
		}

		if len(parts) >= 4 {
			break // max 4 parts for column ref
		}
	}

	return &ast.ColumnRef{
		Parts: parts,
		Loc:   ast.Loc{Start: parts[0].Loc.Start, End: parts[len(parts)-1].Loc.End},
	}, nil
}

// parseFuncCall parses a function call starting after the opening '('.
// The name and startLoc identify the function name and its starting position.
func (p *Parser) parseFuncCall(name ast.ObjectName, startLoc ast.Loc) (*ast.FuncCallExpr, error) {
	p.advance() // consume '('

	funcName := strings.ToUpper(name.Name.Name)
	fc := &ast.FuncCallExpr{
		Name: name,
		Loc:  ast.Loc{Start: startLoc.Start},
	}

	// COUNT(*) special handling
	if funcName == "COUNT" && p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
		return fc, nil
	}

	// DISTINCT handling for aggregates like COUNT(DISTINCT x)
	if p.cur.Type == kwDISTINCT {
		p.advance()
		fc.Distinct = true
	}

	// TRIM special handling: TRIM( expr ) or TRIM( LEADING|TRAILING|BOTH expr FROM expr )
	// We parse TRIM as a regular function call with arguments.

	if p.cur.Type != ')' {
		args, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		fc.Args = args
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End

	// Check for WITHIN GROUP (ORDER BY ...)
	if p.cur.Type == kwWITHIN {
		p.advance() // consume WITHIN
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwORDER); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderItems, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		fc.OrderBy = orderItems
		closeOrd, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeOrd.Loc.End
	}

	return fc, nil
}

// parseParenExpr parses a parenthesized expression or subquery.
// If it encounters SELECT or WITH after the open paren, it parses a
// subquery and wraps it in SubqueryExpr.
func (p *Parser) parseParenExpr() (ast.Node, error) {
	openTok := p.advance() // consume '('
	startLoc := openTok.Loc.Start

	// Subquery: (SELECT ...) or (WITH ... SELECT ...), optionally with set ops
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		var query ast.Node
		var err error
		if p.cur.Type == kwWITH {
			query, err = p.parseWithQueryExpr()
		} else {
			query, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.SubqueryExpr{Query: query, Loc: ast.Loc{Start: startLoc, End: closeTok.Loc.End}}, nil
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}

	return &ast.ParenExpr{
		Expr: expr,
		Loc:  ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseCaseExpr parses a CASE expression (simple or searched).
func (p *Parser) parseCaseExpr() (*ast.CaseExpr, error) {
	caseTok := p.advance() // consume CASE

	ce := &ast.CaseExpr{
		Loc: ast.Loc{Start: caseTok.Loc.Start},
	}

	// Determine if this is a simple CASE (CASE expr WHEN ...) or
	// searched CASE (CASE WHEN ...).
	if p.cur.Type != kwWHEN {
		// Simple CASE: parse the operand expression.
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Kind = ast.CaseSimple
		ce.Operand = operand
	} else {
		ce.Kind = ast.CaseSearched
	}

	// Parse WHEN clauses.
	for p.cur.Type == kwWHEN {
		whenTok := p.advance() // consume WHEN
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		result, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Whens = append(ce.Whens, &ast.WhenClause{
			Cond:   cond,
			Result: result,
			Loc:    ast.Loc{Start: whenTok.Loc.Start, End: ast.NodeLoc(result).End},
		})
	}

	if len(ce.Whens) == 0 {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected WHEN in CASE expression",
		}
	}

	// Optional ELSE.
	if p.cur.Type == kwELSE {
		p.advance() // consume ELSE
		elseExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Else = elseExpr
	}

	endTok, err := p.expect(kwEND)
	if err != nil {
		return nil, err
	}
	ce.Loc.End = endTok.Loc.End

	return ce, nil
}

// parseCastExpr parses CAST( expr AS type ).
func (p *Parser) parseCastExpr() (*ast.CastExpr, error) {
	castTok := p.advance() // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Expr:     expr,
		TypeName: typeName,
		Loc:      ast.Loc{Start: castTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseTryCastExpr parses TRY_CAST( expr AS type ).
func (p *Parser) parseTryCastExpr() (*ast.CastExpr, error) {
	castTok := p.advance() // consume TRY_CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Expr:     expr,
		TypeName: typeName,
		TryCast:  true,
		Loc:      ast.Loc{Start: castTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseIffExpr parses IFF( cond, then, else ).
func (p *Parser) parseIffExpr() (*ast.IffExpr, error) {
	iffTok := p.advance() // consume IFF
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	thenExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	elseExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	return &ast.IffExpr{
		Cond: cond,
		Then: thenExpr,
		Else: elseExpr,
		Loc:  ast.Loc{Start: iffTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseArrayLiteral parses [ elem, elem, ... ].
func (p *Parser) parseArrayLiteral() (*ast.ArrayLiteralExpr, error) {
	openTok := p.advance() // consume '['

	arr := &ast.ArrayLiteralExpr{
		Loc: ast.Loc{Start: openTok.Loc.Start},
	}

	if p.cur.Type != ']' {
		elems, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		arr.Elements = elems
	}

	closeTok, err := p.expect(']')
	if err != nil {
		return nil, err
	}
	arr.Loc.End = closeTok.Loc.End

	return arr, nil
}

// parseJsonLiteral parses { key: value, key: value, ... }.
// Keys can be quoted strings or identifiers.
func (p *Parser) parseJsonLiteral() (*ast.JsonLiteralExpr, error) {
	openTok := p.advance() // consume '{'

	jl := &ast.JsonLiteralExpr{
		Loc: ast.Loc{Start: openTok.Loc.Start},
	}

	if p.cur.Type != '}' {
		for {
			pairLoc := p.cur.Loc
			var key string
			if p.cur.Type == tokString {
				tok := p.advance()
				key = tok.Str
			} else if p.isIdentToken() {
				tok := p.advance()
				key = tok.Str
			} else {
				return nil, &ParseError{
					Loc: p.cur.Loc,
					Msg: "expected string or identifier key in JSON literal",
				}
			}

			if _, err := p.expect(':'); err != nil {
				return nil, err
			}

			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			jl.Pairs = append(jl.Pairs, ast.KeyValuePair{
				Key:   key,
				Value: val,
				Loc:   ast.Loc{Start: pairLoc.Start, End: ast.NodeLoc(val).End},
			})

			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}

	closeTok, err := p.expect('}')
	if err != nil {
		return nil, err
	}
	jl.Loc.End = closeTok.Loc.End

	return jl, nil
}

// parseLambdaExpr parses a lambda expression: x -> body or (x, y) -> body.
// Called when we know the current token is an ident followed by ->.
func (p *Parser) parseLambdaExpr() (*ast.LambdaExpr, error) {
	startLoc := p.cur.Loc

	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	params := []ast.Ident{ident}

	// Consume ->
	if _, err := p.expect(tokArrow); err != nil {
		return nil, err
	}

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.LambdaExpr{
		Params: params,
		Body:   body,
		Loc:    ast.Loc{Start: startLoc.Start, End: ast.NodeLoc(body).End},
	}, nil
}

// ---------------------------------------------------------------------------
// Infix/postfix handlers
// ---------------------------------------------------------------------------

// parseIsExpr parses IS [NOT] NULL or IS [NOT] DISTINCT FROM expr.
func (p *Parser) parseIsExpr(left ast.Node) (ast.Node, error) {
	p.advance() // consume IS

	notFlag := false
	if p.cur.Type == kwNOT {
		p.advance()
		notFlag = true
	}

	// IS [NOT] DISTINCT FROM expr
	if p.cur.Type == kwDISTINCT {
		p.advance() // consume DISTINCT
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		right, err := p.parseExprPrec(bpComparison + 1)
		if err != nil {
			return nil, err
		}
		return &ast.IsExpr{
			Expr:         left,
			Not:          notFlag,
			Null:         false,
			DistinctFrom: right,
			Loc:          ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}, nil
	}

	// IS [NOT] NULL
	if p.cur.Type != kwNULL {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected NULL or DISTINCT after IS [NOT]",
		}
	}
	nullTok := p.advance()

	return &ast.IsExpr{
		Expr: left,
		Not:  notFlag,
		Null: true,
		Loc:  ast.Loc{Start: ast.NodeLoc(left).Start, End: nullTok.Loc.End},
	}, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left ast.Node) (ast.Node, error) {
	notFlag := false
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		notFlag = true
	}

	p.advance() // consume BETWEEN

	// Parse low bound at bpComparison+1 to avoid re-consuming AND as an
	// infix binary operator.
	low, err := p.parseExprPrec(bpComparison + 1)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected AND in BETWEEN expression",
		}
	}

	high, err := p.parseExprPrec(bpComparison + 1)
	if err != nil {
		return nil, err
	}

	return &ast.BetweenExpr{
		Expr: left,
		Low:  low,
		High: high,
		Not:  notFlag,
		Loc:  ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(high).End},
	}, nil
}

// parseInExpr parses [NOT] IN ( values ).
func (p *Parser) parseInExpr(left ast.Node) (ast.Node, error) {
	notFlag := false
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		notFlag = true
	}

	p.advance() // consume IN

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Subquery: IN (SELECT ...) or IN (WITH ... SELECT ...), optionally with set ops
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		var query ast.Node
		var err error
		if p.cur.Type == kwWITH {
			query, err = p.parseWithQueryExpr()
		} else {
			query, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		openLoc := ast.NodeLoc(left).Start
		subq := &ast.SubqueryExpr{Query: query, Loc: ast.Loc{Start: openLoc, End: closeTok.Loc.End}}
		return &ast.InExpr{
			Expr:   left,
			Values: []ast.Node{subq},
			Not:    notFlag,
			Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
		}, nil
	}

	values, err := p.parseExprList()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}

	return &ast.InExpr{
		Expr:   left,
		Values: values,
		Not:    notFlag,
		Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
	}, nil
}

// parseLikeExpr parses [NOT] LIKE/ILIKE/RLIKE/REGEXP pattern [ESCAPE esc]
// or LIKE ANY (...). The NOT has already been consumed if present (notFlag=true).
// The likeOp and notFlag are determined by the caller (infixBindingPower).
func (p *Parser) parseLikeExpr(left ast.Node, likeOp ast.LikeOp, notFlag bool) (ast.Node, error) {
	// If we have NOT, it was recognized in infixBindingPower; need to
	// consume it here if present.
	if notFlag && p.cur.Type == kwNOT {
		p.advance() // consume NOT
	}

	// Consume the LIKE/ILIKE/RLIKE keyword (or REGEXP ident).
	if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "REGEXP" {
		p.advance()
	} else {
		p.advance() // consume kwLIKE, kwILIKE, kwRLIKE
	}

	// LIKE ANY (...)
	if likeOp == ast.LikeOpLike && p.cur.Type == kwANY {
		p.advance() // consume ANY
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		values, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.LikeExpr{
			Expr:      left,
			Op:        likeOp,
			Not:       notFlag,
			Any:       true,
			AnyValues: values,
			Loc:       ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
		}, nil
	}

	// ILIKE ANY (...)
	if likeOp == ast.LikeOpILike && p.cur.Type == kwANY {
		p.advance() // consume ANY
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		values, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.LikeExpr{
			Expr:      left,
			Op:        likeOp,
			Not:       notFlag,
			Any:       true,
			AnyValues: values,
			Loc:       ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
		}, nil
	}

	// Regular LIKE/ILIKE/RLIKE/REGEXP pattern
	pattern, err := p.parseExprPrec(bpComparison + 1)
	if err != nil {
		return nil, err
	}

	le := &ast.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Op:      likeOp,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}

	// Optional ESCAPE clause (only for LIKE and ILIKE)
	if p.cur.Type == kwESCAPE {
		p.advance() // consume ESCAPE
		escExpr, err := p.parseExprPrec(bpComparison + 1)
		if err != nil {
			return nil, err
		}
		le.Escape = escExpr
		le.Loc.End = ast.NodeLoc(escExpr).End
	}

	return le, nil
}

// parseCastPostfix parses the :: type cast postfix operator.
func (p *Parser) parseCastPostfix(left ast.Node) (ast.Node, error) {
	p.advance() // consume ::

	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	return &ast.CastExpr{
		Expr:       left,
		TypeName:   typeName,
		ColonColon: true,
		Loc:        ast.Loc{Start: ast.NodeLoc(left).Start, End: typeName.Loc.End},
	}, nil
}

// parseCollatePostfix parses COLLATE collation_name.
func (p *Parser) parseCollatePostfix(left ast.Node) (ast.Node, error) {
	p.advance() // consume COLLATE

	// Collation can be a string literal or an identifier.
	var collation string
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		collation = tok.Str
		return &ast.CollateExpr{
			Expr:      left,
			Collation: collation,
			Loc:       ast.Loc{Start: ast.NodeLoc(left).Start, End: tok.Loc.End},
		}, nil
	default:
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		collation = ident.Name
		return &ast.CollateExpr{
			Expr:      left,
			Collation: collation,
			Loc:       ast.Loc{Start: ast.NodeLoc(left).Start, End: ident.Loc.End},
		}, nil
	}
}

// parseOverPostfix parses OVER ( window_spec ) and attaches it to the
// preceding function call.
func (p *Parser) parseOverPostfix(left ast.Node) (ast.Node, error) {
	fc, ok := left.(*ast.FuncCallExpr)
	if !ok {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "OVER can only follow a function call",
		}
	}

	p.advance() // consume OVER

	ws, err := p.parseOverClause()
	if err != nil {
		return nil, err
	}
	fc.Over = ws
	fc.Loc.End = ws.Loc.End

	return fc, nil
}

// parseOverClause parses ( PARTITION BY ... ORDER BY ... frame ).
func (p *Parser) parseOverClause() (*ast.WindowSpec, error) {
	openTok, err := p.expect('(')
	if err != nil {
		return nil, err
	}

	ws := &ast.WindowSpec{
		Loc: ast.Loc{Start: openTok.Loc.Start},
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		parts, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		ws.PartitionBy = parts
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderItems, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		ws.OrderBy = orderItems
	}

	// Window frame: ROWS | RANGE | GROUPS
	if p.cur.Type == kwROWS || p.isFrameKeyword() {
		frame, err := p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
		ws.Frame = frame
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	ws.Loc.End = closeTok.Loc.End

	return ws, nil
}

// isFrameKeyword checks if the current token is RANGE or GROUPS. Since
// RANGE is not a keyword in F2, we check by ident comparison. GROUPS is
// a keyword (kwGROUPS). ROWS is also a keyword (kwROWS).
func (p *Parser) isFrameKeyword() bool {
	if p.cur.Type == kwGROUPS {
		return true
	}
	if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "RANGE" {
		return true
	}
	return false
}

// parseWindowFrame parses ROWS|RANGE|GROUPS BETWEEN start AND end.
func (p *Parser) parseWindowFrame() (*ast.WindowFrame, error) {
	startLoc := p.cur.Loc

	var kind ast.WindowFrameKind
	switch {
	case p.cur.Type == kwROWS:
		kind = ast.FrameRows
		p.advance()
	case p.cur.Type == kwGROUPS:
		kind = ast.FrameGroups
		p.advance()
	case p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "RANGE":
		kind = ast.FrameRange
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	frame := &ast.WindowFrame{
		Kind: kind,
		Loc:  ast.Loc{Start: startLoc.Start},
	}

	// BETWEEN start AND end
	if p.cur.Type == kwBETWEEN {
		p.advance() // consume BETWEEN

		startBound, err := p.parseWindowBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound

		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}

		endBound, err := p.parseWindowBound()
		if err != nil {
			return nil, err
		}
		frame.End = endBound
		frame.Loc.End = p.prev.Loc.End
	} else {
		// Single bound (no BETWEEN): just the start bound
		startBound, err := p.parseWindowBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound
		frame.Loc.End = p.prev.Loc.End
	}

	return frame, nil
}

// parseWindowBound parses one window frame bound:
//   - UNBOUNDED PRECEDING
//   - UNBOUNDED FOLLOWING
//   - CURRENT ROW
//   - N PRECEDING
//   - N FOLLOWING
func (p *Parser) parseWindowBound() (ast.WindowBound, error) {
	// UNBOUNDED PRECEDING / UNBOUNDED FOLLOWING
	if p.isUnbounded() {
		p.advance() // consume UNBOUNDED

		if p.isPreceding() {
			p.advance()
			return ast.WindowBound{Kind: ast.BoundUnboundedPreceding}, nil
		}
		if p.isFollowing() {
			p.advance()
			return ast.WindowBound{Kind: ast.BoundUnboundedFollowing}, nil
		}
		return ast.WindowBound{}, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected PRECEDING or FOLLOWING after UNBOUNDED",
		}
	}

	// CURRENT ROW
	if p.cur.Type == kwCURRENT {
		p.advance() // consume CURRENT
		if p.cur.Type == kwROW {
			p.advance()
			return ast.WindowBound{Kind: ast.BoundCurrentRow}, nil
		}
		return ast.WindowBound{}, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected ROW after CURRENT",
		}
	}

	// N PRECEDING / N FOLLOWING
	offset, err := p.parseExprPrec(bpComparison + 1)
	if err != nil {
		return ast.WindowBound{}, err
	}

	if p.isPreceding() {
		p.advance()
		return ast.WindowBound{Kind: ast.BoundPreceding, Offset: offset}, nil
	}
	if p.isFollowing() {
		p.advance()
		return ast.WindowBound{Kind: ast.BoundFollowing, Offset: offset}, nil
	}

	return ast.WindowBound{}, &ParseError{
		Loc: p.cur.Loc,
		Msg: "expected PRECEDING or FOLLOWING",
	}
}

// isUnbounded checks if the current token is the UNBOUNDED keyword.
// UNBOUNDED is not in F2's keyword table, so we check via tokIdent.
func (p *Parser) isUnbounded() bool {
	return p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "UNBOUNDED"
}

// isPreceding checks if the current token is the PRECEDING keyword.
// PRECEDING is not in F2's keyword table, so we check via tokIdent.
func (p *Parser) isPreceding() bool {
	return p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "PRECEDING"
}

// isFollowing checks if the current token is the FOLLOWING keyword.
// FOLLOWING is not in F2's keyword table, so we check via tokIdent.
func (p *Parser) isFollowing() bool {
	return p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "FOLLOWING"
}

// ---------------------------------------------------------------------------
// Access operators
// ---------------------------------------------------------------------------

// parseColonAccess parses : field (JSON path access).
func (p *Parser) parseColonAccess(left ast.Node) (ast.Node, error) {
	p.advance() // consume ':'

	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &ast.AccessExpr{
		Expr:  left,
		Kind:  ast.AccessColon,
		Field: ident,
		Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ident.Loc.End},
	}, nil
}

// parseBracketAccess parses [ index ] (array subscript).
func (p *Parser) parseBracketAccess(left ast.Node) (ast.Node, error) {
	p.advance() // consume '['

	idx, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(']')
	if err != nil {
		return nil, err
	}

	return &ast.AccessExpr{
		Expr:  left,
		Kind:  ast.AccessBracket,
		Index: idx,
		Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
	}, nil
}

// parseDotAccess parses .field (dot path chaining for semi-structured data).
func (p *Parser) parseDotAccess(left ast.Node) (ast.Node, error) {
	p.advance() // consume '.'

	// In dot access context, allow '$' prefix identifiers and general identifiers.
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &ast.AccessExpr{
		Expr:  left,
		Kind:  ast.AccessDot,
		Field: ident,
		Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ident.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() ([]ast.Node, error) {
	var exprs []ast.Node

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprs = append(exprs, expr)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

// parseOrderByList parses ORDER BY items: expr [ASC|DESC] [NULLS FIRST|LAST].
func (p *Parser) parseOrderByList() ([]*ast.OrderItem, error) {
	var items []*ast.OrderItem

	for {
		item, err := p.parseOrderItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return items, nil
}

// parseOrderItem parses a single ORDER BY item.
func (p *Parser) parseOrderItem() (*ast.OrderItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	item := &ast.OrderItem{
		Expr: expr,
		Loc:  ast.NodeLoc(expr),
	}

	// ASC / DESC
	if p.cur.Type == kwASC {
		p.advance()
		item.Desc = false
		item.Loc.End = p.prev.Loc.End
	} else if p.cur.Type == kwDESC {
		p.advance()
		item.Desc = true
		item.Loc.End = p.prev.Loc.End
	}

	// NULLS FIRST / NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance() // consume NULLS
		if p.cur.Type == kwFIRST {
			p.advance()
			nf := true
			item.NullsFirst = &nf
			item.Loc.End = p.prev.Loc.End
		} else if p.cur.Type == kwLAST {
			p.advance()
			nf := false
			item.NullsFirst = &nf
			item.Loc.End = p.prev.Loc.End
		} else {
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "expected FIRST or LAST after NULLS",
			}
		}
	}

	return item, nil
}

// ---------------------------------------------------------------------------
// Public helper
// ---------------------------------------------------------------------------

// ParseExpr parses an expression from a standalone string. Returns the
// parsed AST node and any errors. Useful for tests and callers that have
// an expression string but not a token stream.
func ParseExpr(input string) (ast.Node, []ParseError) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	node, err := p.parseExpr()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}

	if p.cur.Type != tokEOF {
		return node, []ParseError{{
			Loc: p.cur.Loc,
			Msg: "unexpected token after expression: " + tokenDesc(p.cur),
		}}
	}

	return node, nil
}

// tokenDesc returns a human-readable description of a token for error messages.
func tokenDesc(tok Token) string {
	if tok.Str != "" {
		return strconv.Quote(tok.Str)
	}
	return TokenName(tok.Type)
}
