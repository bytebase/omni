package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// Binding power constants for Pratt parsing (low to high).
// Doris operator precedence from DorisParser.g4:
//
//	OR / || < XOR < AND / && < NOT (prefix) < predicates/comparison <
//	bitwise | < bitwise & < shifts << >> < add +- < mul */% DIV MOD <
//	bitwise ^ < unary (-, +, ~) < primary
const (
	bpNone    = 0
	bpOr      = 10 // OR, ||
	bpXor     = 15 // XOR
	bpAnd     = 20 // AND, &&
	bpNot     = 25 // NOT (prefix)
	bpPred    = 30 // IS, BETWEEN, IN, LIKE, REGEXP/RLIKE
	bpCmp     = 35 // =, <, >, <=, >=, <>, !=, <=>
	bpBitOr   = 40 // |
	bpBitAnd  = 45 // &
	bpShift   = 50 // <<, >>
	bpAdd     = 55 // +, -
	bpMul     = 60 // *, /, %, DIV
	bpBitXor  = 65 // ^
	bpUnary   = 70 // unary -, +, ~
	bpPrimary = 80 // atoms
)

// ---------------------------------------------------------------------------
// Core Pratt loop
// ---------------------------------------------------------------------------

// parseExpr is the entry point for expression parsing.
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
				left, err = p.parseBetweenExpr(left, false)
			case "NOT_BETWEEN":
				left, err = p.parseBetweenExpr(left, true)
			case "IN":
				left, err = p.parseInExpr(left, false)
			case "NOT_IN":
				left, err = p.parseInExpr(left, true)
			case "LIKE":
				left, err = p.parseLikeExpr(left, false)
			case "NOT_LIKE":
				left, err = p.parseLikeExpr(left, true)
			case "REGEXP":
				left, err = p.parseRegexpExpr(left, false)
			case "NOT_REGEXP":
				left, err = p.parseRegexpExpr(left, true)
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		// Regular binary operator: consume the operator, parse right side.
		p.advance()
		right, err := p.parseExprPrec(bp + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}
	}

	return left, nil
}

// infixBindingPower returns (bp, op, special, ok) for the current token.
// If special is non-empty, the caller must dispatch to the named handler
// instead of building a BinaryExpr.
func (p *Parser) infixBindingPower() (int, ast.BinaryOp, string, bool) {
	switch p.cur.Kind {
	// Logical
	case kwOR, tokDoublePipes:
		return bpOr, ast.BinOr, "", true
	case kwXOR:
		return bpXor, ast.BinXor, "", true
	case kwAND, tokLogicalAnd:
		return bpAnd, ast.BinAnd, "", true

	// NOT as infix prefix for BETWEEN, IN, LIKE, REGEXP/RLIKE
	case kwNOT:
		next := p.peekNext()
		switch next.Kind {
		case kwBETWEEN:
			return bpPred, 0, "NOT_BETWEEN", true
		case kwIN:
			return bpPred, 0, "NOT_IN", true
		case kwLIKE:
			return bpPred, 0, "NOT_LIKE", true
		case kwREGEXP, kwRLIKE:
			return bpPred, 0, "NOT_REGEXP", true
		}
		return 0, 0, "", false

	// Predicates
	case kwIS:
		return bpPred, 0, "IS", true
	case kwBETWEEN:
		return bpPred, 0, "BETWEEN", true
	case kwIN:
		return bpPred, 0, "IN", true
	case kwLIKE:
		return bpPred, 0, "LIKE", true
	case kwREGEXP, kwRLIKE:
		return bpPred, 0, "REGEXP", true

	// Comparison operators
	case int('='):
		return bpCmp, ast.BinEq, "", true
	case tokNotEq:
		return bpCmp, ast.BinNe, "", true
	case int('<'):
		return bpCmp, ast.BinLt, "", true
	case int('>'):
		return bpCmp, ast.BinGt, "", true
	case tokLessEq:
		return bpCmp, ast.BinLe, "", true
	case tokGreaterEq:
		return bpCmp, ast.BinGe, "", true
	case tokNullSafeEq:
		return bpCmp, ast.BinNullSafeEq, "", true

	// Bitwise OR
	case int('|'):
		return bpBitOr, ast.BinBitOr, "", true

	// Bitwise AND
	case int('&'):
		return bpBitAnd, ast.BinBitAnd, "", true

	// Shifts
	case tokShiftLeft:
		return bpShift, ast.BinShiftLeft, "", true
	case tokShiftRight:
		return bpShift, ast.BinShiftRight, "", true

	// Addition / subtraction
	case int('+'):
		return bpAdd, ast.BinAdd, "", true
	case int('-'):
		return bpAdd, ast.BinSub, "", true

	// Multiplication / division / modulo
	case int('*'):
		return bpMul, ast.BinMul, "", true
	case int('/'):
		return bpMul, ast.BinDiv, "", true
	case int('%'):
		return bpMul, ast.BinMod, "", true
	case kwDIV:
		return bpMul, ast.BinIntDiv, "", true

	// Bitwise XOR
	case int('^'):
		return bpBitXor, ast.BinBitXor, "", true
	}

	return 0, 0, "", false
}

// ---------------------------------------------------------------------------
// Prefix dispatch
// ---------------------------------------------------------------------------

// parsePrefixExpr handles prefix unary operators (-, +, ~, NOT) and delegates
// everything else to parsePrimaryExpr.
func (p *Parser) parsePrefixExpr() (ast.Node, error) {
	switch p.cur.Kind {
	case int('-'):
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

	case int('+'):
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

	case int('~'):
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryBitNot,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwNOT:
		// NOT as a prefix unary operator. Only if the next token does NOT form
		// an infix predicate (NOT BETWEEN, NOT IN, NOT LIKE, NOT REGEXP/RLIKE).
		// In the prefix case we are at the start of an expression, so NOT always
		// acts as prefix.
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
		return p.parseExistsExpr()

	case kwINTERVAL:
		return p.parseIntervalExpr()

	default:
		return p.parsePrimaryExpr()
	}
}

// ---------------------------------------------------------------------------
// Primary expression dispatch
// ---------------------------------------------------------------------------

// parsePrimaryExpr parses atomic expressions: literals, identifiers,
// function calls, parenthesized expressions, CASE, CAST, etc.
func (p *Parser) parsePrimaryExpr() (ast.Node, error) {
	switch p.cur.Kind {
	case tokInt:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitInt,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokFloat:
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
			Loc:   tok.Loc,
		}, nil

	case kwFALSE:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBool,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case kwNULL:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitNull,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case int('('):
		return p.parseParenExprOrSubquery()

	case kwCASE:
		return p.parseCaseExpr()

	case kwCAST:
		return p.parseCastExpr(false)

	case kwTRY_CAST:
		return p.parseCastExpr(true)

	// Nilary functions: CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP,
	// CURRENT_USER, SESSION_USER, LOCALTIME, LOCALTIMESTAMP.
	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP,
		kwCURRENT_USER, kwSESSION_USER, kwLOCALTIME, kwLOCALTIMESTAMP:
		return p.parseNilaryFunction()

	case tokPlaceholder:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitNull, // placeholder treated as a null-ish literal
			Value: "?",
			Loc:   tok.Loc,
		}, nil

	default:
		// Identifier-based: column ref or function call.
		if p.isExprIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// isExprIdentToken reports whether the current token is usable as an identifier
// in expression context: tokIdent, tokQuotedIdent, or a non-reserved keyword.
func (p *Parser) isExprIdentToken() bool {
	switch p.cur.Kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		return p.cur.Kind >= 700 && !IsReserved(p.cur.Kind)
	}
}

// ---------------------------------------------------------------------------
// Primary expression helpers
// ---------------------------------------------------------------------------

// parseIdentExpr parses an identifier that could be a column ref or function
// call (name followed by '(').
func (p *Parser) parseIdentExpr() (ast.Node, error) {
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for function call: name(
	if p.cur.Kind == int('(') {
		return p.parseFuncCall(name)
	}

	// Simple or qualified column reference.
	return &ast.ColumnRef{
		Name: name,
		Loc:  name.Loc,
	}, nil
}

// parseFuncCall parses a function call starting at the opening '('.
// The function name has already been parsed as an ObjectName.
func (p *Parser) parseFuncCall(name *ast.ObjectName) (*ast.FuncCallExpr, error) {
	p.advance() // consume '('

	funcName := strings.ToUpper(name.Parts[len(name.Parts)-1])
	fc := &ast.FuncCallExpr{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	// COUNT(*) special handling
	if funcName == "COUNT" && p.cur.Kind == int('*') {
		p.advance()
		fc.Star = true
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
		return fc, nil
	}

	// DISTINCT handling for aggregates like COUNT(DISTINCT x)
	if p.cur.Kind == kwDISTINCT {
		p.advance()
		fc.Distinct = true
	}

	if p.cur.Kind != int(')') {
		args, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		fc.Args = args

		// Optional ORDER BY within aggregate functions (e.g., GROUP_CONCAT)
		if p.cur.Kind == kwORDER {
			p.advance() // consume ORDER
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			orderItems, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			fc.OrderBy = orderItems
		}

		// Optional SEPARATOR for GROUP_CONCAT
		if p.cur.Kind == kwSEPARATOR {
			p.advance() // consume SEPARATOR
			if p.cur.Kind != tokString {
				return nil, &ParseError{
					Loc: p.cur.Loc,
					Msg: "expected string after SEPARATOR",
				}
			}
			fc.Separator = p.advance().Str
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End

	return fc, nil
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() ([]ast.Node, error) {
	var list []ast.Node
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}

// parseOrderByList parses a comma-separated list of ORDER BY items.
func (p *Parser) parseOrderByList() ([]*ast.OrderByItem, error) {
	var items []*ast.OrderByItem
	for {
		item, err := p.parseOrderByItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseOrderByItem parses one ORDER BY item: expr [ASC|DESC] [NULLS FIRST|LAST].
func (p *Parser) parseOrderByItem() (*ast.OrderByItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	item := &ast.OrderByItem{
		Expr: expr,
		Loc:  ast.NodeLoc(expr),
	}

	// Optional ASC/DESC
	if p.cur.Kind == kwASC {
		p.advance()
		item.Desc = false
		item.Loc.End = p.prev.Loc.End
	} else if p.cur.Kind == kwDESC {
		p.advance()
		item.Desc = true
		item.Loc.End = p.prev.Loc.End
	}

	// Optional NULLS FIRST/LAST
	if p.cur.Kind == kwNULLS {
		p.advance() // consume NULLS
		switch p.cur.Kind {
		case kwFIRST:
			p.advance()
			b := true
			item.NullsFirst = &b
			item.Loc.End = p.prev.Loc.End
		case kwLAST:
			p.advance()
			b := false
			item.NullsFirst = &b
			item.Loc.End = p.prev.Loc.End
		default:
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "expected FIRST or LAST after NULLS",
			}
		}
	}

	return item, nil
}

// parseParenExprOrSubquery parses a parenthesized expression or subquery.
// For subqueries (SELECT/WITH after open paren), we create a placeholder
// SubqueryExpr that stores the raw text.
func (p *Parser) parseParenExprOrSubquery() (ast.Node, error) {
	openTok := p.advance() // consume '('

	// Subquery: (SELECT ...) or (WITH ... SELECT ...)
	if p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
		return p.parseSubqueryPlaceholder(openTok.Loc.Start)
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	return &ast.ParenExpr{
		Expr: expr,
		Loc:  ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseSubqueryPlaceholder consumes tokens until the matching closing ')' for
// a subquery expression and returns a SubqueryExpr with the raw text.
// startOffset is the byte offset of the opening '(' token.
func (p *Parser) parseSubqueryPlaceholder(startOffset int) (*ast.SubqueryExpr, error) {
	depth := 1
	subStart := p.cur.Loc.Start
	var subEnd int
	for depth > 0 && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
			if depth == 0 {
				subEnd = p.cur.Loc.Start
				break
			}
		}
		if depth > 0 {
			p.advance()
		}
	}

	if depth != 0 {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "unterminated subquery",
		}
	}

	rawText := p.input[subStart:subEnd]
	closeTok := p.advance() // consume ')'

	return &ast.SubqueryExpr{
		RawText: strings.TrimSpace(rawText),
		Loc:     ast.Loc{Start: startOffset, End: closeTok.Loc.End},
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
	if p.cur.Kind != kwWHEN {
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
	for p.cur.Kind == kwWHEN {
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
	if p.cur.Kind == kwELSE {
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

// parseCastExpr parses CAST(expr AS type) or TRY_CAST(expr AS type).
func (p *Parser) parseCastExpr(tryCast bool) (*ast.CastExpr, error) {
	castTok := p.advance() // consume CAST or TRY_CAST
	if _, err := p.expect(int('(')); err != nil {
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
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Expr:     expr,
		TypeName: typeName,
		TryCast:  tryCast,
		Loc:      ast.Loc{Start: castTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseNilaryFunction parses zero-argument keyword functions like
// CURRENT_DATE, CURRENT_TIMESTAMP, etc.
func (p *Parser) parseNilaryFunction() (ast.Node, error) {
	tok := p.advance()
	name := strings.ToUpper(tok.Str)
	objName := &ast.ObjectName{
		Parts: []string{name},
		Loc:   tok.Loc,
	}

	// Some of these may optionally have empty parens: CURRENT_TIMESTAMP()
	fc := &ast.FuncCallExpr{
		Name: objName,
		Loc:  tok.Loc,
	}
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
	}
	return fc, nil
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() (ast.Node, error) {
	existsTok := p.advance() // consume EXISTS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	subq, err := p.parseSubqueryPlaceholder(existsTok.Loc.Start)
	if err != nil {
		return nil, err
	}

	return &ast.ExistsExpr{
		Subquery: subq,
		Loc:      ast.Loc{Start: existsTok.Loc.Start, End: subq.Loc.End},
	}, nil
}

// parseIntervalExpr parses INTERVAL expr unit.
func (p *Parser) parseIntervalExpr() (ast.Node, error) {
	intervalTok := p.advance() // consume INTERVAL

	value, err := p.parseExprPrec(bpPrimary)
	if err != nil {
		return nil, err
	}

	// Parse interval unit keyword
	unit, err := p.parseIntervalUnit()
	if err != nil {
		return nil, err
	}

	return &ast.IntervalExpr{
		Value: value,
		Unit:  unit,
		Loc:   ast.Loc{Start: intervalTok.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseIntervalUnit parses and returns the interval unit keyword as a string.
func (p *Parser) parseIntervalUnit() (string, error) {
	switch p.cur.Kind {
	case kwYEAR, kwMONTH, kwWEEK, kwDAY, kwHOUR, kwMINUTE, kwSECOND:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	default:
		// Allow any identifier-like token as a unit for forward compatibility
		if p.cur.Kind == tokIdent {
			tok := p.advance()
			return strings.ToUpper(tok.Str), nil
		}
		return "", &ParseError{
			Loc: p.cur.Loc,
			Msg: fmt.Sprintf("expected interval unit, got %q", p.cur.Str),
		}
	}
}

// ---------------------------------------------------------------------------
// Infix/postfix handlers
// ---------------------------------------------------------------------------

// parseIsExpr parses IS [NOT] NULL/TRUE/FALSE.
func (p *Parser) parseIsExpr(left ast.Node) (ast.Node, error) {
	p.advance() // consume IS

	notFlag := false
	if p.cur.Kind == kwNOT {
		p.advance()
		notFlag = true
	}

	var isWhat string
	switch p.cur.Kind {
	case kwNULL:
		isWhat = "NULL"
	case kwTRUE:
		isWhat = "TRUE"
	case kwFALSE:
		isWhat = "FALSE"
	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected NULL, TRUE, or FALSE after IS [NOT]",
		}
	}
	endTok := p.advance()

	return &ast.IsExpr{
		Expr:   left,
		Not:    notFlag,
		IsWhat: isWhat,
		Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: endTok.Loc.End},
	}, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume BETWEEN

	// Parse low bound at bpPred+1 to avoid re-consuming AND as an
	// infix binary operator.
	low, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected AND in BETWEEN expression",
		}
	}

	high, err := p.parseExprPrec(bpPred + 1)
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

// parseInExpr parses [NOT] IN (values...) or [NOT] IN (subquery).
func (p *Parser) parseInExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume IN

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Subquery: IN (SELECT ...) or IN (WITH ... SELECT ...)
	if p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
		subq, err := p.parseSubqueryPlaceholder(ast.NodeLoc(left).Start)
		if err != nil {
			return nil, err
		}
		return &ast.InExpr{
			Expr:   left,
			Values: []ast.Node{subq},
			Not:    notFlag,
			Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: subq.Loc.End},
		}, nil
	}

	values, err := p.parseExprList()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(int(')'))
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

// parseLikeExpr parses [NOT] LIKE pattern [ESCAPE esc].
func (p *Parser) parseLikeExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume LIKE

	pattern, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	le := &ast.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}

	// Optional ESCAPE clause
	if p.cur.Kind == kwESCAPE {
		p.advance() // consume ESCAPE
		esc, err := p.parseExprPrec(bpPred + 1)
		if err != nil {
			return nil, err
		}
		le.Escape = esc
		le.Loc.End = ast.NodeLoc(esc).End
	}

	return le, nil
}

// parseRegexpExpr parses [NOT] REGEXP|RLIKE pattern.
func (p *Parser) parseRegexpExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume REGEXP or RLIKE

	pattern, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	return &ast.RegexpExpr{
		Expr:    left,
		Pattern: pattern,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}, nil
}
