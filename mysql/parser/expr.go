package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// Precedence levels for Pratt parsing (low to high, matching MySQL).
const (
	precNone       = 0
	precAssign     = 1  // :=
	precOr         = 2  // OR, ||
	precXor        = 3  // XOR
	precAnd        = 4  // AND, &&
	precNot        = 5  // NOT (prefix)
	precComparison = 6  // =, <=>, >=, >, <=, <, <>, !=, IS, LIKE, REGEXP, IN, BETWEEN
	precBitOr      = 7  // |
	precBitAnd     = 8  // &
	precShift      = 9  // <<, >>
	precAdd        = 10 // +, -
	precMul        = 11 // *, /, DIV, %, MOD
	precBitXor     = 12 // ^
	precUnary      = 13 // -, ~, !
	precCollate    = 14 // COLLATE
	precJsonAccess = 15 // ->, ->>
)

// parseExpr parses an expression using Pratt parsing / precedence climbing.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/expressions.html
func (p *Parser) parseExpr() (nodes.ExprNode, error) {
	return p.parseExprPrec(precAssign)
}

// parseExprPrec parses an expression with the given minimum precedence.
func (p *Parser) parseExprPrec(minPrec int) (nodes.ExprNode, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}

	for {
		// MEMBER OF special handling (not in infixPrecedence since it's keyword-based)
		if p.cur.Type == kwMEMBER {
			left, err = p.parseMemberOfExpr(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		prec, binOp, ok := p.infixPrecedence()
		if !ok || prec < minPrec {
			break
		}

		// Handle special infix operators
		switch {
		case p.cur.Type == kwIS:
			left, err = p.parseIsExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwBETWEEN || (p.cur.Type == kwNOT && p.peekNext().Type == kwBETWEEN):
			left, err = p.parseBetweenExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwIN || (p.cur.Type == kwNOT && p.peekNext().Type == kwIN):
			left, err = p.parseInExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwLIKE || (p.cur.Type == kwNOT && p.peekNext().Type == kwLIKE):
			left, err = p.parseLikeExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwSOUNDS:
			left, err = p.parseSoundsLikeExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwREGEXP || p.cur.Type == kwRLIKE ||
			(p.cur.Type == kwNOT && (p.peekNext().Type == kwREGEXP || p.peekNext().Type == kwRLIKE)):
			left, err = p.parseRegexpExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwCOLLATE:
			colStart := p.pos()
			p.advance() // consume COLLATE
			collation, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			left = &nodes.CollateExpr{
				Loc:       nodes.Loc{Start: colStart, End: p.pos()},
				Expr:      left,
				Collation: collation,
			}
			continue
		}

		// Regular binary operator
		opStart := p.pos()
		p.advance() // consume operator

		// Right-associative for assignment
		nextPrec := prec + 1
		if binOp == nodes.BinOpAssign {
			nextPrec = prec
		}

		right, err := p.parseExprPrec(nextPrec)
		if err != nil {
			return nil, err
		}

		left = &nodes.BinaryExpr{
			Loc:   nodes.Loc{Start: opStart},
			Op:    binOp,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

// infixPrecedence returns the precedence and operator for the current token
// if it is an infix operator.
func (p *Parser) infixPrecedence() (int, nodes.BinaryOp, bool) {
	switch p.cur.Type {
	case tokAssign:
		return precAssign, nodes.BinOpAssign, true
	case kwOR:
		return precOr, nodes.BinOpOr, true
	case kwXOR:
		return precXor, nodes.BinOpBitXor, true
	case kwAND:
		return precAnd, nodes.BinOpAnd, true

	// NOT as infix: NOT IN, NOT LIKE, NOT BETWEEN, NOT REGEXP
	case kwNOT:
		switch p.peekNext().Type {
		case kwIN, kwLIKE, kwBETWEEN, kwREGEXP, kwRLIKE:
			return precComparison, 0, true // handled specially
		}
		return 0, 0, false

	// Comparison operators
	case '=':
		return precComparison, nodes.BinOpEq, true
	case tokNullSafeEq:
		return precComparison, nodes.BinOpNullSafeEq, true
	case tokNotEq:
		return precComparison, nodes.BinOpNe, true
	case '<':
		return precComparison, nodes.BinOpLt, true
	case '>':
		return precComparison, nodes.BinOpGt, true
	case tokLessEq:
		return precComparison, nodes.BinOpLe, true
	case tokGreaterEq:
		return precComparison, nodes.BinOpGe, true
	case kwIS:
		return precComparison, 0, true // handled specially
	case kwLIKE:
		return precComparison, 0, true
	case kwREGEXP, kwRLIKE:
		return precComparison, 0, true
	case kwIN:
		return precComparison, 0, true
	case kwBETWEEN:
		return precComparison, 0, true
	case kwSOUNDS:
		return precComparison, 0, true

	// Bit operators
	case '|':
		return precBitOr, nodes.BinOpBitOr, true
	case '&':
		return precBitAnd, nodes.BinOpBitAnd, true
	case tokShiftLeft:
		return precShift, nodes.BinOpShiftLeft, true
	case tokShiftRight:
		return precShift, nodes.BinOpShiftRight, true

	// Arithmetic
	case '+':
		return precAdd, nodes.BinOpAdd, true
	case '-':
		return precAdd, nodes.BinOpSub, true
	case '*':
		return precMul, nodes.BinOpMul, true
	case '/':
		return precMul, nodes.BinOpDiv, true
	case kwDIV:
		return precMul, nodes.BinOpDivInt, true
	case '%':
		return precMul, nodes.BinOpMod, true
	case kwMOD:
		return precMul, nodes.BinOpMod, true

	case '^':
		return precBitXor, nodes.BinOpBitXor, true

	case kwCOLLATE:
		return precCollate, 0, true

	// JSON column-path operators
	case tokJsonExtract:
		return precJsonAccess, nodes.BinOpJsonExtract, true
	case tokJsonUnquote:
		return precJsonAccess, nodes.BinOpJsonUnquote, true
	}

	return 0, 0, false
}

// parsePrefixExpr parses prefix/primary expressions.
func (p *Parser) parsePrefixExpr() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case '-':
		return p.parseUnaryExpr(nodes.UnaryMinus)
	case '+':
		// Unary plus — just parse the operand
		p.advance()
		return p.parsePrefixExpr()
	case '~':
		return p.parseUnaryExpr(nodes.UnaryBitNot)
	case '!':
		return p.parseUnaryExpr(nodes.UnaryNot)
	case kwNOT:
		return p.parseUnaryExpr(nodes.UnaryNot)
	case kwEXISTS_KW:
		return p.parseExistsExpr()
	case kwCASE:
		return p.parseCaseExpr()
	case kwINTERVAL:
		return p.parseIntervalExpr()
	case kwMATCH:
		return p.parseMatchExpr()
	case kwBINARY:
		return p.parseUnaryExpr(nodes.UnaryBinary)
	default:
		return p.parsePrimaryExpr()
	}
}

// parseUnaryExpr parses a unary expression.
func (p *Parser) parseUnaryExpr(op nodes.UnaryOp) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume operator

	operand, err := p.parseExprPrec(precUnary)
	if err != nil {
		return nil, err
	}

	return &nodes.UnaryExpr{
		Loc:     nodes.Loc{Start: start},
		Op:      op,
		Operand: operand,
	}, nil
}

// parsePrimaryExpr parses atoms: literals, column refs, subqueries, func calls, parenthesized exprs.
func (p *Parser) parsePrimaryExpr() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Ival}, nil

	case tokFCONST:
		tok := p.advance()
		return &nodes.FloatLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Str}, nil

	case tokSCONST:
		tok := p.advance()
		return &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Str}, nil

	case tokXCONST:
		tok := p.advance()
		return &nodes.HexLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Str}, nil

	case tokBCONST:
		tok := p.advance()
		return &nodes.BitLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Str}, nil

	case kwTRUE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: true}, nil

	case kwFALSE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: false}, nil

	case kwNULL:
		tok := p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: tok.Loc}}, nil

	case kwDEFAULT:
		return p.parseDefaultExpr()

	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwCURRENT_USER, kwLOCALTIME, kwLOCALTIMESTAMP:
		tok := p.advance()
		name := strings.ToUpper(tok.Str)
		fc := &nodes.FuncCallExpr{Loc: nodes.Loc{Start: tok.Loc}, Name: name}
		// Optional parentheses with optional fsp argument
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type != ')' {
				// Parse arguments (e.g., fractional seconds precision)
				for {
					arg, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					fc.Args = append(fc.Args, arg)
					if p.cur.Type != ',' {
						break
					}
					p.advance() // consume ','
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		fc.Loc.End = p.pos()
		return fc, nil

	case '*':
		tok := p.advance()
		return &nodes.StarExpr{Loc: nodes.Loc{Start: tok.Loc}}, nil

	case '(':
		return p.parseParenExpr()

	case kwCAST:
		return p.parseCastExpr()

	case kwCONVERT:
		return p.parseConvertExpr()

	case kwROW:
		return p.parseRowConstructor()

	case kwVALUES:
		return p.parseValuesFunc()

	default:
		// Variable reference
		if p.isVariableRef() {
			return p.parseVariableRef()
		}

		// Identifier — could be column ref or function call
		if p.isIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, &ParseError{
			Message:  "expected expression",
			Position: p.cur.Loc,
		}
	}
}

// parseIdentExpr parses an identifier that could be a column ref or function call.
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
	start := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for function call: name(
	if p.cur.Type == '(' {
		// But first check if this could be a schema-qualified function: name.name(
		// No — MySQL doesn't use schema.func() syntax in expressions (use schema.func())
		return p.parseFuncCall(start, "", name)
	}

	// Check for qualified name: name.name or name.name.name
	if p.cur.Type == '.' {
		p.advance()

		// table.*
		if p.cur.Type == '*' {
			p.advance()
			return &nodes.ColumnRef{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Table: name,
				Star:  true,
			}, nil
		}

		name2, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		// name.name( — schema.func()
		if p.cur.Type == '(' {
			return p.parseFuncCall(start, name, name2)
		}

		// Check for third part: schema.table.col or schema.table.*
		if p.cur.Type == '.' {
			p.advance()
			if p.cur.Type == '*' {
				p.advance()
				return &nodes.ColumnRef{
					Loc:    nodes.Loc{Start: start, End: p.pos()},
					Schema: name,
					Table:  name2,
					Star:   true,
				}, nil
			}
			name3, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &nodes.ColumnRef{
				Loc:    nodes.Loc{Start: start, End: p.pos()},
				Schema: name,
				Table:  name2,
				Column: name3,
			}, nil
		}

		// table.col
		return &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.pos()},
			Table:  name,
			Column: name2,
		}, nil
	}

	// Plain column ref
	return &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Column: name,
	}, nil
}

// parseFuncCall parses a function call after the name has been consumed.
func (p *Parser) parseFuncCall(start int, schema, name string) (nodes.ExprNode, error) {
	p.advance() // consume '('

	fc := &nodes.FuncCallExpr{
		Loc:    nodes.Loc{Start: start},
		Schema: schema,
		Name:   strings.ToUpper(name),
	}

	// Handle special function forms
	upperName := strings.ToUpper(name)

	// COUNT(*)
	if upperName == "COUNT" && p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		goto overCheck
	}

	// Aggregate with DISTINCT: COUNT(DISTINCT ...), SUM(DISTINCT ...), etc.
	if p.cur.Type == kwDISTINCT {
		fc.Distinct = true
		p.advance()
	}

	// TRIM special syntax: TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str)
	if upperName == "TRIM" {
		return p.parseTrimFunc(fc)
	}

	// SUBSTRING special syntax: SUBSTRING(str, pos, len) or SUBSTRING(str FROM pos FOR len)
	if upperName == "SUBSTRING" || upperName == "SUBSTR" {
		return p.parseSubstringFunc(fc)
	}

	// GROUP_CONCAT special syntax
	if upperName == "GROUP_CONCAT" {
		return p.parseGroupConcatFunc(fc)
	}

	// Empty argument list
	if p.cur.Type == ')' {
		p.advance()
		goto overCheck
	}

	// Regular argument list
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

overCheck:

	// OVER clause for window functions
	//
	// Ref: https://dev.mysql.com/doc/refman/8.0/en/window-functions-usage.html
	//
	//   over_clause:
	//       {OVER (window_spec) | OVER window_name}
	//
	//   window_spec:
	//       [window_name] [partition_clause] [order_clause] [frame_clause]
	//
	//   partition_clause:
	//       PARTITION BY expr [, expr] ...
	//
	//   order_clause:
	//       ORDER BY expr [ASC|DESC] [, expr [ASC|DESC]] ...
	//
	//   frame_clause:
	//       frame_units frame_extent
	//
	//   frame_units:
	//       {ROWS | RANGE | GROUPS}
	//
	//   frame_extent:
	//       {frame_start | frame_between}
	//
	//   frame_between:
	//       BETWEEN frame_start AND frame_end
	//
	//   frame_start, frame_end: {
	//       CURRENT ROW
	//     | UNBOUNDED PRECEDING
	//     | UNBOUNDED FOLLOWING
	//     | expr PRECEDING
	//     | expr FOLLOWING
	//   }
	if p.cur.Type == kwOVER {
		p.advance()
		wd, err := p.parseOverClause()
		if err != nil {
			return nil, err
		}
		fc.Over = wd
	}

	fc.Loc.End = p.pos()
	return fc, nil
}

// parseTrimFunc parses TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str).
// parseMemberOfExpr parses value MEMBER OF(json_array).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html
//
//	value MEMBER OF(json_array)
func (p *Parser) parseMemberOfExpr(value nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume MEMBER

	// Expect OF
	if p.cur.Type != tokIDENT || !eqFold(p.cur.Str, "of") {
		return nil, &ParseError{
			Message:  "expected OF after MEMBER",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume OF

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	array, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.MemberOfExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Value: value,
		Array: array,
	}, nil
}

// parseOverClause parses OVER (window_spec) or OVER window_name.
func (p *Parser) parseOverClause() (*nodes.WindowDef, error) {
	start := p.pos()

	// OVER window_name (identifier reference)
	if p.cur.Type != '(' {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &nodes.WindowDef{
			Loc:     nodes.Loc{Start: start, End: p.pos()},
			RefName: name,
		}, nil
	}

	// OVER (window_spec)
	p.advance() // consume '('

	wd := &nodes.WindowDef{Loc: nodes.Loc{Start: start}}

	// Optional window_name reference
	if p.isIdentToken() && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
		p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
		wd.RefName, _, _ = p.parseIdentifier()
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		wd.PartitionBy = exprs
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		wd.OrderBy = orderBy
	}

	// frame_clause: {ROWS | RANGE | GROUPS} frame_extent
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		frame, err := p.parseFrameClause()
		if err != nil {
			return nil, err
		}
		wd.Frame = frame
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	wd.Loc.End = p.pos()
	return wd, nil
}

// parseFrameClause parses a window frame clause.
func (p *Parser) parseFrameClause() (*nodes.WindowFrame, error) {
	start := p.pos()
	frame := &nodes.WindowFrame{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwROWS:
		frame.Type = nodes.FrameRows
	case kwRANGE:
		frame.Type = nodes.FrameRange
	case kwGROUPS:
		frame.Type = nodes.FrameGroups
	}
	p.advance()

	// frame_extent: frame_start | BETWEEN frame_start AND frame_end
	if _, ok := p.match(kwBETWEEN); ok {
		startBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		endBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.End = endBound
	} else {
		// Single frame_start
		startBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound
	}

	frame.Loc.End = p.pos()
	return frame, nil
}

// parseFrameBound parses a window frame bound.
func (p *Parser) parseFrameBound() (*nodes.WindowFrameBound, error) {
	start := p.pos()
	bound := &nodes.WindowFrameBound{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwCURRENT {
		p.advance()
		p.match(kwROW)
		bound.Type = nodes.BoundCurrentRow
	} else if p.cur.Type == kwUNBOUNDED {
		p.advance()
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundUnboundedPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundUnboundedFollowing
		} else {
			return nil, &ParseError{
				Message:  "expected PRECEDING or FOLLOWING after UNBOUNDED",
				Position: p.cur.Loc,
			}
		}
	} else {
		// expr PRECEDING | expr FOLLOWING
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		bound.Offset = expr
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundFollowing
		} else {
			return nil, &ParseError{
				Message:  "expected PRECEDING or FOLLOWING",
				Position: p.cur.Loc,
			}
		}
	}

	bound.Loc.End = p.pos()
	return bound, nil
}

func (p *Parser) parseTrimFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Check for LEADING, TRAILING, BOTH
	switch p.cur.Type {
	case kwLEADING, kwTRAILING, kwBOTH:
		// Include as first "argument" indicator via the function name
		fc.Name = "TRIM_" + strings.ToUpper(p.cur.Str)
		p.advance()
	}

	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == kwFROM {
		p.advance()
		// remstr FROM str
		str, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = []nodes.ExprNode{arg, str}
	} else {
		fc.Args = []nodes.ExprNode{arg}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.pos()
	return fc, nil
}

// parseSubstringFunc parses SUBSTRING(str, pos, len) or SUBSTRING(str FROM pos [FOR len]).
func (p *Parser) parseSubstringFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	fc.Args = append(fc.Args, arg)

	if p.cur.Type == kwFROM {
		p.advance()
		pos, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, pos)

		if p.cur.Type == kwFOR {
			p.advance()
			length, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, length)
		}
	} else if p.cur.Type == ',' {
		p.advance()
		pos, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, pos)

		if p.cur.Type == ',' {
			p.advance()
			length, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, length)
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.pos()
	return fc, nil
}

// parseGroupConcatFunc parses GROUP_CONCAT([DISTINCT] expr [, expr ...] [ORDER BY ...] [SEPARATOR str]).
func (p *Parser) parseGroupConcatFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	if p.cur.Type == ')' {
		p.advance()
		fc.Loc.End = p.pos()
		return fc, nil
	}

	// Parse arguments
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, arg)

		if p.cur.Type == ',' {
			// Peek to see if next is ORDER or SEPARATOR — if so, don't consume comma
			next := p.peekNext()
			if next.Type == kwORDER || next.Type == kwSEPARATOR {
				break
			}
			p.advance()
		} else {
			break
		}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		fc.OrderBy = orderBy
	}

	// SEPARATOR
	if p.cur.Type == kwSEPARATOR {
		p.advance()
		sep, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Separator = sep
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.pos()
	return fc, nil
}

// parseParenExpr parses a parenthesized expression or subquery.
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume '('

	// Check for subquery
	if p.cur.Type == kwSELECT {
		sub, err := p.parseSubqueryExpr()
		if err != nil {
			return nil, err
		}
		if _, errP := p.expect(')'); errP != nil {
			return nil, errP
		}
		sub.Loc.Start = start
		sub.Loc.End = p.pos()
		return sub, nil
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.ParenExpr{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Expr: expr,
	}, nil
}

// parseCastExpr parses CAST(expr AS type).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html
//
//	CAST(expr AS data_type)
func (p *Parser) parseCastExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CAST

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

	dt, err := p.parseCastDataType()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.CastExpr{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Expr:     expr,
		TypeName: dt,
	}, nil
}

// parseCastDataType parses data types valid in CAST expressions.
// CAST supports a subset of types: BINARY, CHAR, DATE, DATETIME, DECIMAL, SIGNED, UNSIGNED, TIME, JSON, etc.
func (p *Parser) parseCastDataType() (*nodes.DataType, error) {
	start := p.pos()

	// SIGNED [INTEGER]
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "signed") {
		p.advance()
		p.match(kwINTEGER)
		p.match(kwINT)
		return &nodes.DataType{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "SIGNED"}, nil
	}

	// UNSIGNED [INTEGER]
	if p.cur.Type == kwUNSIGNED {
		p.advance()
		p.match(kwINTEGER)
		p.match(kwINT)
		return &nodes.DataType{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "UNSIGNED"}, nil
	}

	return p.parseDataType()
}

// parseConvertExpr parses CONVERT(expr, type) or CONVERT(expr USING charset).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html
//
//	CONVERT(expr, type)
//	CONVERT(expr USING transcoding_name)
func (p *Parser) parseConvertExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CONVERT

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	conv := &nodes.ConvertExpr{
		Loc:  nodes.Loc{Start: start},
		Expr: expr,
	}

	if p.cur.Type == kwUSING {
		p.advance()
		charset, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		conv.Charset = charset
	} else if p.cur.Type == ',' {
		p.advance()
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		conv.TypeName = dt
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	conv.Loc.End = p.pos()
	return conv, nil
}

// parseIsExpr parses IS [NOT] NULL / TRUE / FALSE / UNKNOWN.
func (p *Parser) parseIsExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume IS

	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}

	is := &nodes.IsExpr{
		Loc:  nodes.Loc{Start: start},
		Not:  not,
		Expr: left,
	}

	switch p.cur.Type {
	case kwNULL:
		is.Test = nodes.IsNull
		p.advance()
	case kwTRUE:
		is.Test = nodes.IsTrue
		p.advance()
	case kwFALSE:
		is.Test = nodes.IsFalse
		p.advance()
	default:
		// IS [NOT] UNKNOWN — check for identifier "unknown"
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "unknown") {
			is.Test = nodes.IsUnknown
			p.advance()
		} else {
			return nil, &ParseError{
				Message:  "expected NULL, TRUE, FALSE, or UNKNOWN after IS",
				Position: p.cur.Loc,
			}
		}
	}

	is.Loc.End = p.pos()
	return is, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume BETWEEN

	low, err := p.parseExprPrec(precAdd) // parse at higher precedence to avoid AND confusion
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, err
	}

	high, err := p.parseExprPrec(precAdd)
	if err != nil {
		return nil, err
	}

	return &nodes.BetweenExpr{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Not:  not,
		Expr: left,
		Low:  low,
		High: high,
	}, nil
}

// parseInExpr parses [NOT] IN (value_list) or [NOT] IN (subquery).
func (p *Parser) parseInExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume IN

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	inExpr := &nodes.InExpr{
		Loc:  nodes.Loc{Start: start},
		Not:  not,
		Expr: left,
	}

	// Check for subquery
	if p.cur.Type == kwSELECT {
		sub, err := p.parseSubqueryExpr()
		if err != nil {
			return nil, err
		}
		if _, errP := p.expect(')'); errP != nil {
			return nil, errP
		}
		inExpr.Select = sub.Select
		inExpr.Loc.End = p.pos()
		return inExpr, nil
	}

	// Value list
	for {
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		inExpr.List = append(inExpr.List, val)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	inExpr.Loc.End = p.pos()
	return inExpr, nil
}

// parseLikeExpr parses [NOT] LIKE pattern [ESCAPE escape_char].
func (p *Parser) parseLikeExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume LIKE

	pattern, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}

	like := &nodes.LikeExpr{
		Loc:     nodes.Loc{Start: start},
		Not:     not,
		Expr:    left,
		Pattern: pattern,
	}

	if _, ok := p.match(kwESCAPE); ok {
		esc, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		like.Escape = esc
	}

	like.Loc.End = p.pos()
	return like, nil
}

// parseRegexpExpr parses [NOT] REGEXP pattern.
func (p *Parser) parseRegexpExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume REGEXP or RLIKE

	pattern, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}

	// Represent as BinaryExpr with BinOpRegexp
	expr := &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Op:    nodes.BinOpRegexp,
		Left:  left,
		Right: pattern,
	}

	if not {
		return &nodes.UnaryExpr{
			Loc:     nodes.Loc{Start: start, End: p.pos()},
			Op:      nodes.UnaryNot,
			Operand: expr,
		}, nil
	}

	return expr, nil
}

// parseCaseExpr parses a CASE expression.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/flow-control-functions.html#operator_case
//
//	CASE [operand]
//	    WHEN condition THEN result
//	    [WHEN condition THEN result ...]
//	    [ELSE result]
//	END
func (p *Parser) parseCaseExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CASE

	ce := &nodes.CaseExpr{Loc: nodes.Loc{Start: start}}

	// Simple CASE: CASE operand WHEN ...
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Operand = operand
	}

	// WHEN clauses
	for p.cur.Type == kwWHEN {
		whenStart := p.pos()
		p.advance() // consume WHEN
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
		ce.Whens = append(ce.Whens, &nodes.CaseWhen{
			Loc:    nodes.Loc{Start: whenStart, End: p.pos()},
			Cond:   cond,
			Result: result,
		})
	}

	// ELSE
	if _, ok := p.match(kwELSE); ok {
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Default = def
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}

	ce.Loc.End = p.pos()
	return ce, nil
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume EXISTS

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	sub, err := p.parseSubqueryExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.ExistsExpr{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Select: sub.Select,
	}, nil
}

// parseIntervalExpr parses INTERVAL value unit.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/expressions.html#temporal-intervals
//
//	INTERVAL expr unit
func (p *Parser) parseIntervalExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume INTERVAL

	val, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	// Parse unit: DAY, HOUR, MINUTE, SECOND, MONTH, YEAR, etc.
	unit, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.IntervalExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Value: val,
		Unit:  strings.ToUpper(unit),
	}, nil
}

// parseMatchExpr parses MATCH (col_list) AGAINST (expr [modifier]).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/fulltext-search.html
//
//	MATCH (col1 [, col2, ...]) AGAINST (expr [search_modifier])
//
//	search_modifier:
//	  {
//	       IN NATURAL LANGUAGE MODE
//	     | IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION
//	     | IN BOOLEAN MODE
//	     | WITH QUERY EXPANSION
//	  }
func (p *Parser) parseMatchExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume MATCH

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	me := &nodes.MatchExpr{Loc: nodes.Loc{Start: start}}

	// Column list
	for {
		ref, err := p.parseColumnRef()
		if err != nil {
			return nil, err
		}
		me.Columns = append(me.Columns, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAGAINST); err != nil {
		return nil, err
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Parse against expression with high enough precedence to avoid consuming
	// IN (which is precComparison) as part of the expression — IN here starts
	// the search modifier, not an IN-list comparison.
	against, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}
	me.Against = against

	// Search modifier
	if p.cur.Type == kwIN {
		p.advance()
		if p.cur.Type == kwNATURAL {
			// IN NATURAL LANGUAGE MODE [WITH QUERY EXPANSION]
			p.advance() // consume NATURAL
			if _, err := p.expect(kwLANGUAGE); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwMODE); err != nil {
				return nil, err
			}
			// Check for optional WITH QUERY EXPANSION
			if p.cur.Type == kwWITH && p.peekNext().Type == kwQUERY {
				p.advance() // consume WITH
				p.advance() // consume QUERY
				if _, err := p.expect(kwEXPANSION); err != nil {
					return nil, err
				}
				me.Modifier = "IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION"
			} else {
				me.Modifier = "IN NATURAL LANGUAGE MODE"
			}
		} else {
			// IN BOOLEAN MODE
			if _, err := p.expect(kwBOOLEAN); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwMODE); err != nil {
				return nil, err
			}
			me.Modifier = "IN BOOLEAN MODE"
		}
	} else if p.cur.Type == kwWITH && p.peekNext().Type == kwQUERY {
		// WITH QUERY EXPANSION
		p.advance() // consume WITH
		p.advance() // consume QUERY
		if _, err := p.expect(kwEXPANSION); err != nil {
			return nil, err
		}
		me.Modifier = "WITH QUERY EXPANSION"
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	me.Loc.End = p.pos()
	return me, nil
}

// parseSoundsLikeExpr parses expr SOUNDS LIKE expr.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html
//
//	expr1 SOUNDS LIKE expr2
//
// This is equivalent to SOUNDEX(expr1) = SOUNDEX(expr2).
func (p *Parser) parseSoundsLikeExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume SOUNDS

	if _, err := p.expect(kwLIKE); err != nil {
		return nil, err
	}

	right, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}

	return &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Op:    nodes.BinOpSoundsLike,
		Left:  left,
		Right: right,
	}, nil
}

// parseRowConstructor parses a ROW(expr, expr, ...) constructor.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/row-subqueries.html
//
//	ROW(val1, val2, ..., valN)
func (p *Parser) parseRowConstructor() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume ROW

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	row := &nodes.RowExpr{Loc: nodes.Loc{Start: start}}

	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		row.Items = append(row.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	row.Loc.End = p.pos()
	return row, nil
}

// parseDefaultExpr parses DEFAULT or DEFAULT(col_name).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_default
//
//	DEFAULT
//	DEFAULT(col_name)
func (p *Parser) parseDefaultExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume DEFAULT

	de := &nodes.DefaultExpr{Loc: nodes.Loc{Start: start}}

	// Check for DEFAULT(col_name) form
	if p.cur.Type == '(' {
		p.advance() // consume '('
		col, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		de.Column = col
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	de.Loc.End = p.pos()
	return de, nil
}

// parseValuesFunc parses VALUES(col_name) used in INSERT ... ON DUPLICATE KEY UPDATE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_values
//
//	VALUES(col_name)
func (p *Parser) parseValuesFunc() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume VALUES

	// VALUES without '(' is handled elsewhere (e.g., VALUES row constructor in INSERT)
	// Here we handle the function form VALUES(col_name)
	if p.cur.Type != '(' {
		// Not a function call — this shouldn't be reached from parsePrimaryExpr
		// since VALUES without '(' is not a valid expression
		return nil, &ParseError{
			Message:  "expected '(' after VALUES in expression context",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume '('

	col, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	fc := &nodes.FuncCallExpr{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: "VALUES",
		Args: []nodes.ExprNode{col},
	}
	return fc, nil
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() ([]nodes.ExprNode, error) {
	var list []nodes.ExprNode
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}
