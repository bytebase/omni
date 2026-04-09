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
