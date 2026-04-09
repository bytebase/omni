package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// parseMathOp00 parses the string-concatenation layer.
// PartiQL is unusual here: || binds LOOSER than +/- per
// PartiQLParser.g4 lines 494-497. That contradicts most SQL dialects.
//
// Grammar: mathOp00 (lines 494-497):
//
//	lhs=mathOp00 op=CONCAT rhs=mathOp01
//	| parent_=mathOp01
func (p *Parser) parseMathOp00() (ast.ExprNode, error) {
	left, err := p.parseMathOp01()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokCONCAT {
		p.advance()
		right, err := p.parseMathOp01()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpConcat,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseMathOp01 parses the additive layer (+, -).
//
// Grammar: mathOp01 (lines 499-502):
//
//	lhs=mathOp01 op=(PLUS|MINUS) rhs=mathOp02
//	| parent_=mathOp02
func (p *Parser) parseMathOp01() (ast.ExprNode, error) {
	left, err := p.parseMathOp02()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		op := ast.BinOpAdd
		if p.cur.Type == tokMINUS {
			op = ast.BinOpSub
		}
		p.advance()
		right, err := p.parseMathOp02()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseMathOp02 parses the multiplicative layer (*, /, %).
//
// Grammar: mathOp02 (lines 504-507):
//
//	lhs=mathOp02 op=(PERCENT|ASTERISK|SLASH_FORWARD) rhs=valueExpr
//	| parent_=valueExpr
func (p *Parser) parseMathOp02() (ast.ExprNode, error) {
	left, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokASTERISK || p.cur.Type == tokSLASH_FORWARD || p.cur.Type == tokPERCENT {
		var op ast.BinOp
		switch p.cur.Type {
		case tokASTERISK:
			op = ast.BinOpMul
		case tokSLASH_FORWARD:
			op = ast.BinOpDiv
		case tokPERCENT:
			op = ast.BinOpMod
		}
		p.advance()
		right, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseValueExpr parses the unary-sign layer (+expr, -expr). Unary
// signs are right-associative so they stack naturally via recursion.
//
// Grammar: valueExpr (lines 509-512):
//
//	sign=(PLUS|MINUS) rhs=valueExpr
//	| parent_=exprPrimary
func (p *Parser) parseValueExpr() (ast.ExprNode, error) {
	if p.cur.Type == tokPLUS || p.cur.Type == tokMINUS {
		start := p.cur.Loc.Start
		op := ast.UnOpPos
		if p.cur.Type == tokMINUS {
			op = ast.UnOpNeg
		}
		p.advance()
		operand, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:      op,
			Operand: operand,
			Loc:     ast.Loc{Start: start, End: operand.GetLoc().End},
		}, nil
	}
	return p.parsePrimary()
}

// isComparisonOp reports whether the given token type is one of the
// six comparison operators from exprPredicate#PredicateComparison
// (line 485).
func isComparisonOp(t int) bool {
	switch t {
	case tokEQ, tokNEQ, tokLT, tokLT_EQ, tokGT, tokGT_EQ:
		return true
	}
	return false
}

// tokToComparisonOp maps a comparison token type to its ast.BinOp.
func tokToComparisonOp(t int) ast.BinOp {
	switch t {
	case tokEQ:
		return ast.BinOpEq
	case tokNEQ:
		return ast.BinOpNotEq
	case tokLT:
		return ast.BinOpLt
	case tokLT_EQ:
		return ast.BinOpLtEq
	case tokGT:
		return ast.BinOpGt
	case tokGT_EQ:
		return ast.BinOpGtEq
	}
	return ast.BinOpInvalid
}

// parsePredicate parses the predicate layer. Handles 5 grammar
// alternatives plus the NOT-prefix form for IN/LIKE/BETWEEN:
//
//	comparison (=, <>, <, <=, >, >=)
//	IS [NOT] {NULL|MISSING|TRUE|FALSE}
//	[NOT] IN (expr, expr, ...)
//	[NOT] LIKE mathOp00 [ESCAPE expr]
//	[NOT] BETWEEN mathOp00 AND mathOp00
//
// Grammar: exprPredicate (lines 484-492).
//
// IS type support is RESTRICTED to the 4 values the AST supports
// (NULL/MISSING/TRUE/FALSE). The grammar allows `IS <any type>` but
// ast.IsExpr.Type is a 4-value enum. Any other IS form returns a
// syntax error.
func (p *Parser) parsePredicate() (ast.ExprNode, error) {
	left, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	for {
		startLoc := left.GetLoc().Start
		switch {
		case isComparisonOp(p.cur.Type):
			op := tokToComparisonOp(p.cur.Type)
			p.advance()
			right, err := p.parseMathOp00()
			if err != nil {
				return nil, err
			}
			left = &ast.BinaryExpr{
				Op:    op,
				Left:  left,
				Right: right,
				Loc:   ast.Loc{Start: startLoc, End: right.GetLoc().End},
			}

		case p.cur.Type == tokIS:
			p.advance()
			not := p.match(tokNOT)
			isExpr, err := p.parseIsBody(left, not, startLoc)
			if err != nil {
				return nil, err
			}
			left = isExpr

		case p.cur.Type == tokNOT:
			// NOT IN, NOT LIKE, NOT BETWEEN (lookahead).
			nextType := p.peekNext().Type
			if nextType != tokIN && nextType != tokLIKE && nextType != tokBETWEEN {
				return left, nil
			}
			p.advance() // consume NOT
			switch p.cur.Type {
			case tokIN:
				left, err = p.parseInBody(left, true, startLoc)
			case tokLIKE:
				left, err = p.parseLikeBody(left, true, startLoc)
			case tokBETWEEN:
				left, err = p.parseBetweenBody(left, true, startLoc)
			}
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokIN:
			left, err = p.parseInBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokLIKE:
			left, err = p.parseLikeBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		case p.cur.Type == tokBETWEEN:
			left, err = p.parseBetweenBody(left, false, startLoc)
			if err != nil {
				return nil, err
			}

		default:
			return left, nil
		}
	}
}

// parseIsBody parses the RHS of `expr IS [NOT] {NULL|MISSING|TRUE|FALSE}`.
// The caller has already consumed the IS token and optionally NOT.
func (p *Parser) parseIsBody(left ast.ExprNode, not bool, startLoc int) (*ast.IsExpr, error) {
	var isType ast.IsType
	switch p.cur.Type {
	case tokNULL:
		isType = ast.IsTypeNull
	case tokMISSING:
		isType = ast.IsTypeMissing
	case tokTRUE:
		isType = ast.IsTypeTrue
	case tokFALSE:
		isType = ast.IsTypeFalse
	default:
		return nil, &ParseError{
			Message: "IS predicate requires NULL, MISSING, TRUE, or FALSE",
			Loc:     p.cur.Loc,
		}
	}
	end := p.cur.Loc.End
	p.advance()
	return &ast.IsExpr{
		Expr: left,
		Type: isType,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: end},
	}, nil
}

// parseInBody parses the RHS of `expr [NOT] IN ...`. The caller has
// NOT yet consumed IN.
//
// Grammar: exprPredicate#PredicateIn (lines 487-488).
//
// Foundation implements the parenthesized form only (the common case).
// The un-parenthesized `NOT? IN mathOp00` form is grammar-legal but rare
// and currently returns a parse error pointing at the missing `(`.
func (p *Parser) parseInBody(left ast.ExprNode, not bool, startLoc int) (*ast.InExpr, error) {
	p.advance() // consume IN
	if p.cur.Type != tokPAREN_LEFT {
		return nil, &ParseError{
			Message: "expected ( after IN",
			Loc:     p.cur.Loc,
		}
	}
	p.advance() // consume (
	var list []ast.ExprNode
	if p.cur.Type != tokPAREN_RIGHT {
		for {
			item, err := p.parseMathOp00()
			if err != nil {
				return nil, err
			}
			list = append(list, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.InExpr{
		Expr: left,
		List: list,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: rp.Loc.End},
	}, nil
}

// parseLikeBody parses the RHS of `expr [NOT] LIKE pattern [ESCAPE escape]`.
// The caller has NOT yet consumed LIKE.
//
// Grammar: exprPredicate#PredicateLike (line 489).
func (p *Parser) parseLikeBody(left ast.ExprNode, not bool, startLoc int) (*ast.LikeExpr, error) {
	p.advance() // consume LIKE
	pattern, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	var escape ast.ExprNode
	end := pattern.GetLoc().End
	if p.cur.Type == tokESCAPE {
		p.advance()
		escape, err = p.parseMathOp00()
		if err != nil {
			return nil, err
		}
		end = escape.GetLoc().End
	}
	return &ast.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Escape:  escape,
		Not:     not,
		Loc:     ast.Loc{Start: startLoc, End: end},
	}, nil
}

// parseBetweenBody parses the RHS of `expr [NOT] BETWEEN lower AND upper`.
// The caller has NOT yet consumed BETWEEN.
//
// Grammar: exprPredicate#PredicateBetween (line 490).
func (p *Parser) parseBetweenBody(left ast.ExprNode, not bool, startLoc int) (*ast.BetweenExpr, error) {
	p.advance() // consume BETWEEN
	lower, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokAND); err != nil {
		return nil, err
	}
	upper, err := p.parseMathOp00()
	if err != nil {
		return nil, err
	}
	return &ast.BetweenExpr{
		Expr: left,
		Low:  lower,
		High: upper,
		Not:  not,
		Loc:  ast.Loc{Start: startLoc, End: upper.GetLoc().End},
	}, nil
}

// parseOr parses the OR layer (left-associative).
//
// Grammar: exprOr (lines 469-472):
//
//	lhs=exprOr OR rhs=exprAnd
//	| parent_=exprAnd
func (p *Parser) parseOr() (ast.ExprNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokOR {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpOr,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseAnd parses the AND layer (left-associative).
//
// Grammar: exprAnd (lines 474-477):
//
//	lhs=exprAnd AND rhs=exprNot
//	| parent_=exprNot
func (p *Parser) parseAnd() (ast.ExprNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == tokAND {
		p.advance()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    ast.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End},
		}
	}
	return left, nil
}

// parseNot parses the NOT layer (right-associative prefix).
//
// Grammar: exprNot (lines 479-482):
//
//	<assoc=right> NOT rhs=exprNot
//	| parent_=exprPredicate
func (p *Parser) parseNot() (ast.ExprNode, error) {
	if p.cur.Type == tokNOT {
		start := p.cur.Loc.Start
		p.advance()
		operand, err := p.parseNot() // right-associative recursion
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:      ast.UnOpNot,
			Operand: operand,
			Loc:     ast.Loc{Start: start, End: operand.GetLoc().End},
		}, nil
	}
	return p.parsePredicate()
}
