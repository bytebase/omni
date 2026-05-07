// Package parser - expr.go implements T-SQL expression parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseExpr parses an expression using precedence climbing.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/expressions-transact-sql
//
//	expr = or_expr
func (p *Parser) parseExpr() (nodes.ExprNode, error) {
	return p.parseOr()
}

func exprLocStart(expr nodes.ExprNode, fallback int) int {
	loc := nodes.NodeLoc(expr)
	if loc.Start >= 0 {
		return loc.Start
	}
	return fallback
}

// parseOr parses OR expressions (lowest precedence).
//
//	or_expr = and_expr { OR and_expr }
func (p *Parser) parseOr() (nodes.ExprNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwOR {
		loc := p.pos()
		p.advance()
		// Completion: after OR → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpOr,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}
	}
	return left, nil
}

// parseAnd parses AND expressions.
//
//	and_expr = not_expr { AND not_expr }
func (p *Parser) parseAnd() (nodes.ExprNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwAND {
		loc := p.pos()
		p.advance()
		// Completion: after AND → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}
	}
	return left, nil
}

// parseNot parses NOT expressions.
//
//	not_expr = NOT not_expr
//	         | fulltext_predicate   -- CONTAINS(...) / FREETEXT(...)
//	         | comparison_expr
func (p *Parser) parseNot() (nodes.ExprNode, error) {
	if p.cur.Type == kwNOT {
		loc := p.pos()
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		if operand == nil {
			return nil, p.unexpectedToken()
		}
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryNot,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	if p.inSearchCondition() && (p.cur.Type == kwCONTAINS || p.cur.Type == kwFREETEXT) {
		return p.parseFullTextPredicate()
	}
	return p.parseComparison()
}

// parseComparison parses comparison expressions, including IS, BETWEEN, IN, LIKE.
//
//	comparison_expr = addition_expr
//	    [ IS [NOT] NULL
//	    | [NOT] BETWEEN addition_expr AND addition_expr
//	    | [NOT] IN '(' expr_list ')'
//	    | [NOT] LIKE addition_expr [ESCAPE addition_expr]
//	    | comparison_op addition_expr ]
func (p *Parser) parseComparison() (nodes.ExprNode, error) {
	left, err := p.parseAddition()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}

	// IS [NOT] NULL
	if p.cur.Type == kwIS {
		loc := p.pos()
		p.advance()
		// Completion: after IS → NULL, NOT
		if p.collectMode() {
			p.addTokenCandidate(kwNULL)
			p.addTokenCandidate(kwNOT)
			return nil, errCollecting
		}
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
			// Completion: after IS NOT → NULL
			if p.collectMode() {
				p.addTokenCandidate(kwNULL)
				return nil, errCollecting
			}
		}
		if p.cur.Type == kwNULL {
			p.advance()
			testType := nodes.IsNull
			if not {
				testType = nodes.IsNotNull
			}
			return &nodes.IsExpr{
				Expr:     left,
				TestType: testType,
				Loc:      nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
			}, nil
		}
	}

	// BETWEEN
	if p.cur.Type == kwBETWEEN || (p.cur.Type == kwNOT && p.peekNext().Type == kwBETWEEN) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume BETWEEN
		// Completion: after BETWEEN → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		low, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		if low == nil {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		// Completion: after BETWEEN x AND → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		high, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		if high == nil {
			return nil, p.unexpectedToken()
		}
		return &nodes.BetweenExpr{
			Expr: left,
			Low:  low,
			High: high,
			Not:  not,
			Loc:  nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}, nil
	}

	// [NOT] IN
	if p.cur.Type == kwIN || (p.cur.Type == kwNOT && p.peekNext().Type == kwIN) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume IN
		if _, err := p.expect('('); err != nil {
			return left, nil
		}
		// Completion: after IN ( → columnref, func_name, SELECT
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate(kwSELECT)
			return nil, errCollecting
		}
		// Check for subquery: IN (SELECT ...)
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			sub, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			_, _ = p.expect(')')
			return &nodes.InExpr{
				Expr:     left,
				Subquery: &nodes.SubqueryExpr{Query: sub, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}},
				Not:      not,
				Loc:      nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
			}, nil
		}
		items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				p.addRuleCandidate("func_name")
				p.addTokenCandidate(kwSELECT)
				return nil, errCollecting
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.unexpectedToken()
			}
			return expr, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.InExpr{
			Expr: left,
			List: &nodes.List{Items: items},
			Not:  not,
			Loc:  nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}, nil
	}

	// [NOT] LIKE
	if p.cur.Type == kwLIKE || (p.cur.Type == kwNOT && p.peekNext().Type == kwLIKE) {
		loc := p.pos()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
		}
		p.advance() // consume LIKE
		// Completion: after LIKE → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		pattern, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		if pattern == nil {
			return nil, p.unexpectedToken()
		}
		var escape nodes.ExprNode
		if p.cur.Type == kwESCAPE {
			p.advance()
			escape, err = p.parseAddition()
			if err != nil {
				return nil, err
			}
		}
		return &nodes.LikeExpr{
			Expr:    left,
			Pattern: pattern,
			Escape:  escape,
			Not:     not,
			Loc:     nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}, nil
	}

	// Comparison operators
	op, ok := comparisonOp(p.cur.Type)
	if ok {
		loc := p.pos()
		p.advance()

		// Completion: after comparison operator → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}

		// Check for ANY/SOME/ALL subquery comparison
		if p.cur.Type == kwALL || p.cur.Type == kwSOME || p.cur.Type == kwANY {
			quantifier := strings.ToUpper(p.cur.Str)
			p.advance() // consume ALL/SOME/ANY
			if p.cur.Type == '(' {
				p.advance() // consume (
				subquery, err := p.parseSelectStmt()
				if err != nil {
					return nil, err
				}
				_, _ = p.expect(')')
				return &nodes.SubqueryComparisonExpr{
					Left:       left,
					Op:         op,
					Quantifier: quantifier,
					Subquery:   subquery,
					Loc:        nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
				}, nil
			}
		}

		right, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		return &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}, nil
	}

	return left, nil
}

func comparisonOp(tok int) (nodes.BinaryOp, bool) {
	switch tok {
	case '=':
		return nodes.BinOpEq, true
	case tokNOTEQUAL:
		return nodes.BinOpNeq, true
	case '<':
		return nodes.BinOpLt, true
	case '>':
		return nodes.BinOpGt, true
	case tokLESSEQUAL:
		return nodes.BinOpLte, true
	case tokGREATEQUAL:
		return nodes.BinOpGte, true
	case tokNOTLESS:
		return nodes.BinOpNotLt, true
	case tokNOTGREATER:
		return nodes.BinOpNotGt, true
	default:
		return 0, false
	}
}

// parseAddition parses addition/subtraction and bitwise operators,
// followed by an optional postfix COLLATE clause.
//
//	addition_expr = multiplication_expr { ('+' | '-' | '&' | '|' | '^') multiplication_expr } [ COLLATE collation_name ]
func (p *Parser) parseAddition() (nodes.ExprNode, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == '+' || p.cur.Type == '-' || p.cur.Type == '&' || p.cur.Type == '|' || p.cur.Type == '^' {
		loc := p.pos()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '+':
			op = nodes.BinOpAdd
		case '-':
			op = nodes.BinOpSub
		case '&':
			op = nodes.BinOpBitAnd
		case '|':
			op = nodes.BinOpBitOr
		case '^':
			op = nodes.BinOpBitXor
		}
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		right, err := p.parseMultiplication()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}
	}
	// Postfix COLLATE / AT TIME ZONE (may be chained)
	for {
		if p.cur.Type == kwCOLLATE {
			var err error
			left, err = p.parseCollateExpr(left)
			if err != nil {
				return nil, err
			}
		} else if p.cur.Type == kwAT {
			var err error
			left, err = p.parseAtTimeZoneExpr(left)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return left, nil
}

// parseCollateExpr parses a postfix COLLATE operator on an expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/collations
//
//	expr COLLATE { <collation_name> | database_default }
//
//	<collation_name> ::=
//	    { Windows_collation_name } | { SQL_collation_name }
func (p *Parser) parseCollateExpr(expr nodes.ExprNode) (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume COLLATE
	// Collation name is an identifier (may contain underscores, numbers)
	// or the special keyword database_default
	var collation string
	if p.isIdentLike() {
		collation = p.cur.Str
		p.advance()
	} else {
		return nil, p.unexpectedToken()
	}
	node := &nodes.CollateExpr{
		Expr:      expr,
		Collation: collation,
		Loc:       nodes.Loc{Start: exprLocStart(expr, loc), End: p.prevEnd()},
	}
	node.Loc.End = p.prevEnd()
	return node, nil
}

// parseAtTimeZoneExpr parses a postfix AT TIME ZONE expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/at-time-zone-transact-sql
//
//	inputdate AT TIME ZONE timezone
//
// Can be chained: expr AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'
func (p *Parser) parseAtTimeZoneExpr(expr nodes.ExprNode) (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume AT
	if p.cur.Type != kwTIME {
		// Not AT TIME ZONE — return original expression
		// (AT used in some other context is handled elsewhere)
		return expr, nil
	}
	p.advance() // consume TIME
	if _, err := p.expect(kwZONE); err != nil {
		return expr, nil
	}
	tz, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}
	if tz == nil {
		return nil, p.unexpectedToken()
	}
	node := &nodes.AtTimeZoneExpr{
		Expr:     expr,
		TimeZone: tz,
		Loc:      nodes.Loc{Start: exprLocStart(expr, loc), End: p.prevEnd()},
	}
	node.Loc.End = p.prevEnd()
	return node, nil
}

// parseMultiplication parses multiplication/division/modulo.
//
//	multiplication_expr = unary_expr { ('*' | '/' | '%') unary_expr }
func (p *Parser) parseMultiplication() (nodes.ExprNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == '*' || p.cur.Type == '/' || p.cur.Type == '%' {
		loc := p.pos()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '*':
			op = nodes.BinOpMul
		case '/':
			op = nodes.BinOpDiv
		case '%':
			op = nodes.BinOpMod
		}
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.unexpectedToken()
		}
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}
	}
	return left, nil
}

// parseUnary parses unary operators.
//
//	unary_expr = ('+' | '-' | '~') unary_expr | primary_expr
func (p *Parser) parseUnary() (nodes.ExprNode, error) {
	if p.cur.Type == '+' || p.cur.Type == '-' || p.cur.Type == '~' {
		loc := p.pos()
		var op nodes.UnaryOp
		switch p.cur.Type {
		case '+':
			op = nodes.UnaryPlus
		case '-':
			op = nodes.UnaryMinus
		case '~':
			op = nodes.UnaryBitNot
		}
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if operand == nil {
			return nil, p.unexpectedToken()
		}
		return &nodes.UnaryExpr{
			Op:      op,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	return p.parsePrimary()
}

// parsePrimary parses primary expressions: literals, variables, identifiers,
// function calls, CAST, CONVERT, CASE, COALESCE, NULLIF, IIF, EXISTS, parenthesized.
//
//	primary_expr = integer | float | string | nstring | NULL | DEFAULT
//	             | variable | system_variable | '*'
//	             | '(' expr ')'
//	             | CAST '(' expr AS data_type ')'
//	             | CONVERT '(' data_type ',' expr [',' style] ')'
//	             | CASE ...
//	             | COALESCE '(' expr_list ')'
//	             | NULLIF '(' expr ',' expr ')'
//	             | IIF '(' expr ',' expr ',' expr ')'
//	             | EXISTS '(' select_stmt ')'
//	             | TRY_CAST '(' expr AS data_type ')'
//	             | TRY_CONVERT '(' data_type ',' expr [',' style] ')'
//	             | ident_expr
func (p *Parser) parsePrimary() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitInteger,
			Ival: tok.Ival,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokFCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitFloat,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitString,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokNSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type:    nodes.LitString,
			Str:     tok.Str,
			IsNChar: true,
			Loc:     nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case kwNULL:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitNull,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case kwDEFAULT:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitDefault,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokSYSVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case '*':
		tok := p.advance()
		return &nodes.StarExpr{
			Loc: nodes.Loc{Start: tok.Loc, End: tok.End},
		}, nil
	case '(':
		loc := p.pos()
		p.advance()
		// Completion: after ( in expression → columnref, func_name, subquery keywords
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			p.addTokenCandidate(kwSELECT)
			return nil, errCollecting
		}
		// Scalar subquery: (SELECT ...)
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			sub, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			_, _ = p.expect(')')
			return &nodes.SubqueryExpr{
				Query: sub,
				Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
			}, nil
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		_, _ = p.expect(')')
		return &nodes.ParenExpr{
			Expr: expr,
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case kwCAST:
		if p.peekNext().Type == '(' {
			return p.parseCast()
		}
		return p.parseIdentExpr()
	case kwCONVERT:
		return p.parseConvert()
	case kwCASE:
		return p.parseCaseExpr()
	case kwCOALESCE:
		return p.parseCoalesce()
	case kwNULLIF:
		return p.parseNullif()
	case kwIIF:
		if p.peekNext().Type == '(' {
			return p.parseIif()
		}
		return p.parseIdentExpr()
	case kwNEXT:
		if p.peekNext().Type == kwVALUE {
			return p.parseNextValueFor()
		}
		return p.parseIdentExpr()
	case kwEXISTS:
		return p.parseExists()
	case kwTRY_CAST:
		if p.peekNext().Type == '(' {
			return p.parseTryCast()
		}
		return p.parseIdentExpr()
	case kwTRY_CONVERT:
		return p.parseTryConvert()
	case tokIDENT:
		return p.parseIdentExpr()
	default:
		// T-SQL permits ContextKeyword tokens (e.g. VALUE, ENCRYPTION,
		// SNAPSHOT, ...) as unquoted identifiers in expression position.
		// CoreKeyword tokens (FROM, SELECT, WHERE, ...) are NOT permitted —
		// they must be bracketed (in which case the lexer emits tokIDENT) or
		// quoted. The keyword classification table in lexer.go is the single
		// source of truth; no parallel whitelist is kept here.
		if isContextKeyword(p.cur.Type) {
			return p.parseIdentExpr()
		}
		return nil, nil
	}
}

// parseCast parses CAST(expr AS data_type).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/cast-and-convert-transact-sql
//
//	CAST '(' expr AS data_type ')'
func (p *Parser) parseCast() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after CAST( → columnref, func_name
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	// Completion: after CAST(x AS → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	_, _ = p.expect(')')
	return &nodes.CastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseTryCast parses TRY_CAST(expr AS data_type).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/try-cast-transact-sql
//
//	TRY_CAST '(' expr AS data_type ')'
func (p *Parser) parseTryCast() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume TRY_CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after TRY_CAST( → columnref, func_name
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	// Completion: after TRY_CAST(x AS → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	_, _ = p.expect(')')
	return &nodes.TryCastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseConvert parses CONVERT(data_type, expr [, style]).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/cast-and-convert-transact-sql
//
//	CONVERT '(' data_type ',' expr [',' style] ')'
func (p *Parser) parseConvert() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume CONVERT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after CONVERT( → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		style, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	_, _ = p.expect(')')
	return &nodes.ConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseTryConvert parses TRY_CONVERT(data_type, expr [, style]).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/try-convert-transact-sql
//
//	TRY_CONVERT '(' data_type ',' expr [',' style] ')'
func (p *Parser) parseTryConvert() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume TRY_CONVERT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after TRY_CONVERT( → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		style, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	_, _ = p.expect(')')
	return &nodes.TryConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseCaseExpr parses a CASE expression (both simple and searched).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/case-transact-sql
//
//	CASE [expr] WHEN expr THEN expr { WHEN expr THEN expr } [ELSE expr] END
func (p *Parser) parseCaseExpr() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume CASE

	var arg nodes.ExprNode
	// Simple CASE vs searched CASE
	if p.cur.Type != kwWHEN {
		var err error
		arg, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	var whenList []nodes.Node
	for p.cur.Type == kwWHEN {
		wloc := p.pos()
		p.advance() // consume WHEN
		// Completion: after CASE [expr] WHEN → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		// In a searched CASE (no arg) the WHEN predicate is a
		// <search_condition>. In a simple CASE (with arg) it is a
		// scalar expression compared against arg.
		if arg == nil {
			p.enterSearchCondition()
		}
		cond, err := p.parseExpr()
		if arg == nil {
			p.leaveSearchCondition()
		}
		if err != nil {
			return nil, err
		}
		if cond == nil {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect(kwTHEN); err != nil {
			break
		}
		// Completion: after CASE WHEN x THEN → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		result, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if result == nil {
			return nil, p.unexpectedToken()
		}
		whenList = append(whenList, &nodes.CaseWhen{
			Condition: cond,
			Result:    result,
			Loc:       nodes.Loc{Start: wloc, End: p.prevEnd()},
		})
	}

	var elseExpr nodes.ExprNode
	if _, ok := p.match(kwELSE); ok {
		// Completion: after ELSE → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		var err error
		elseExpr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	_, _ = p.expect(kwEND)

	return &nodes.CaseExpr{
		Arg:      arg,
		WhenList: &nodes.List{Items: whenList},
		ElseExpr: elseExpr,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseCoalesce parses COALESCE(expr, expr, ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/coalesce-transact-sql
//
//	COALESCE '(' expr { ',' expr } ')'
func (p *Parser) parseCoalesce() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume COALESCE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after COALESCE( → columnref, func_name
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	var args []nodes.Node
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
	}
	_, _ = p.expect(')')
	return &nodes.CoalesceExpr{
		Args: &nodes.List{Items: args},
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseNullif parses NULLIF(expr, expr).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/nullif-transact-sql
//
//	NULLIF '(' expr ',' expr ')'
func (p *Parser) parseNullif() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume NULLIF
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after NULLIF( → columnref, func_name
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	left, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	right, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, _ = p.expect(')')
	return &nodes.NullifExpr{
		Left:  left,
		Right: right,
		Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseIif parses IIF(condition, true_val, false_val).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/logical-functions-iif-transact-sql
//
//	IIF '(' expr ',' expr ',' expr ')'
func (p *Parser) parseIif() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume IIF
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after IIF( → columnref, func_name
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	trueVal, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	falseVal, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, _ = p.expect(')')
	return &nodes.IifExpr{
		Condition: cond,
		TrueVal:   trueVal,
		FalseVal:  falseVal,
		Loc:       nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseExists parses EXISTS(subquery).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/exists-transact-sql
//
//	EXISTS '(' select_stmt ')'
func (p *Parser) parseExists() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume EXISTS
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: after EXISTS( → SELECT
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return nil, errCollecting
	}
	query, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, p.unexpectedToken()
	}
	_, _ = p.expect(')')
	return &nodes.ExistsExpr{
		Subquery: query,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

func (p *Parser) parseNextValueFor() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume NEXT
	if _, err := p.expect(kwVALUE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if p.collectMode() {
		p.addRuleCandidate("sequence_ref")
		return nil, errCollecting
	}
	ref, err := p.parseSequenceRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.unexpectedToken()
	}
	return &nodes.ColumnRef{
		Server:   ref.Server,
		Database: ref.Database,
		Schema:   ref.Schema,
		Table:    ref.Schema,
		Column:   ref.Object,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseFuncCall parses a function call after the opening paren has been seen.
//
//	func_name '(' [DISTINCT] [expr_list | '*'] ')' [OVER ...]
func (p *Parser) parseFuncCall(name string, loc int) (nodes.ExprNode, error) {
	nameEnd := p.prevEnd()
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Object: name, Loc: nodes.Loc{Start: loc, End: nameEnd}},
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		_, _ = p.expect(')')
		fc.Loc.End = p.prevEnd()
		if p.cur.Type == kwOVER {
			var err error
			fc.Over, err = p.parseOverClause()
			if err != nil {
				return nil, err
			}
			fc.Loc.End = p.prevEnd()
		}
		return fc, nil
	}

	if p.collectMode() {
		p.addExpressionCandidates()
		return nil, errCollecting
	}

	if p.cur.Type == ')' {
		p.advance()
		fc.Loc.End = p.prevEnd()
		if p.cur.Type == kwOVER {
			var err error
			fc.Over, err = p.parseOverClause()
			if err != nil {
				return nil, err
			}
			fc.Loc.End = p.prevEnd()
		}
		return fc, nil
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	args, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
		if p.collectMode() {
			p.addExpressionCandidates()
			return nil, errCollecting
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		return arg, nil
	})
	if err != nil {
		return nil, err
	}
	fc.Args = &nodes.List{Items: args}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()

	// Check for WITHIN GROUP (ORDER BY ...) clause
	//
	// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/string-agg-transact-sql
	//
	//  STRING_AGG ( expression, separator ) WITHIN GROUP ( ORDER BY order_by_expression [ ASC | DESC ] )
	if p.cur.Type == kwWITHIN {
		var err error
		fc.Within, err = p.parseWithinGroupClause()
		if err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()
	}

	// Check for OVER clause
	if p.cur.Type == kwOVER {
		var err error
		fc.Over, err = p.parseOverClause()
		if err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()
	}

	return fc, nil
}

// parseWithinGroupClause parses WITHIN GROUP (ORDER BY ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/string-agg-transact-sql
//
//	WITHIN GROUP ( ORDER BY order_by_expression [ ASC | DESC ] [ ,...n ] )
func (p *Parser) parseWithinGroupClause() (*nodes.List, error) {
	p.advance() // consume WITHIN
	// Expect GROUP keyword (it's not a reserved keyword, so use matchIdentCI)
	if !p.matchIdentCI("GROUP") {
		return nil, nil
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Expect ORDER BY
	if _, err := p.expect(kwORDER); err != nil {
		_, _ = p.expect(')')
		return nil, err
	}
	if _, err := p.expect(kwBY); err != nil {
		_, _ = p.expect(')')
		return nil, err
	}

	var orders []nodes.Node
	for {
		oloc := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		orders = append(orders, &nodes.OrderByItem{
			Expr:    expr,
			SortDir: dir,
			Loc:     nodes.Loc{Start: oloc, End: p.prevEnd()},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: orders}, nil
}

// parseOverClause parses an OVER clause for window functions.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/select-over-clause-transact-sql
//
//	OVER '(' [PARTITION BY expr_list] [ORDER BY order_item_list] [<ROW or RANGE clause>] ')'
//
//	<ROW or RANGE clause> ::=
//	  { ROWS | RANGE | GROUPS } <window frame extent>
//
//	<window frame extent> ::=
//	  {   <window frame preceding>
//	    | <window frame between>
//	  }
//
//	<window frame between> ::=
//	  BETWEEN <window frame bound> AND <window frame bound>
//
//	<window frame bound> ::=
//	  {   <window frame preceding>
//	    | <window frame following>
//	  }
//
//	<window frame preceding> ::=
//	  {
//	      UNBOUNDED PRECEDING
//	    | <unsigned_value_specification> PRECEDING
//	    | CURRENT ROW
//	  }
//
//	<window frame following> ::=
//	  {
//	      UNBOUNDED FOLLOWING
//	    | <unsigned_value_specification> FOLLOWING
//	    | CURRENT ROW
//	  }
//
//	<unsigned value specification> ::=
//	  { <unsigned integer literal> }
func (p *Parser) parseOverClause() (*nodes.OverClause, error) {
	loc := p.pos()
	p.advance() // consume OVER

	// OVER window_name (reference to named window, no parentheses)
	if p.cur.Type != '(' && p.isIdentLike() {
		over := &nodes.OverClause{
			WindowName: p.cur.Str,
			Loc:        nodes.Loc{Start: loc, End: p.prevEnd()},
		}
		p.advance()
		over.Loc.End = p.prevEnd()
		return over, nil
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside OVER (...) → PARTITION, ORDER
	if p.collectMode() {
		p.addTokenCandidate(kwPARTITION)
		p.addTokenCandidate(kwORDER)
		return nil, errCollecting
	}

	over := &nodes.OverClause{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over, nil
		}
		// Completion: after PARTITION BY → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		list, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		over.PartitionBy = list
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over, nil
		}
		// Completion: after ORDER BY in OVER → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		var orders []nodes.Node
		for {
			oloc := p.pos()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			dir := nodes.SortDefault
			if _, ok := p.match(kwASC); ok {
				dir = nodes.SortAsc
			} else if _, ok := p.match(kwDESC); ok {
				dir = nodes.SortDesc
			}
			orders = append(orders, &nodes.OrderByItem{
				Expr:    expr,
				SortDir: dir,
				Loc:     nodes.Loc{Start: oloc, End: p.prevEnd()},
			})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		over.OrderBy = &nodes.List{Items: orders}
	}

	// ROWS | RANGE | GROUPS window frame clause
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		var err error
		over.WindowFrame, err = p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
	}

	_, _ = p.expect(')')
	over.Loc.End = p.prevEnd()
	return over, nil
}

// parseWindowFrame parses a window frame specification (ROWS/RANGE/GROUPS).
//
//	{ ROWS | RANGE | GROUPS } <window frame extent>
//
//	<window frame extent> ::=
//	  {   <window frame preceding>
//	    | BETWEEN <window frame bound> AND <window frame bound>
//	  }
func (p *Parser) parseWindowFrame() (*nodes.WindowFrame, error) {
	frame := &nodes.WindowFrame{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	switch p.cur.Type {
	case kwROWS:
		frame.Type = nodes.FrameRows
	case kwRANGE:
		frame.Type = nodes.FrameRange
	case kwGROUPS:
		frame.Type = nodes.FrameGroups
	}
	p.advance()

	// Completion: after ROWS/RANGE/GROUPS → BETWEEN, UNBOUNDED, CURRENT
	if p.collectMode() {
		p.addTokenCandidate(kwBETWEEN)
		p.addTokenCandidate(kwUNBOUNDED)
		p.addTokenCandidate(kwCURRENT)
		return nil, errCollecting
	}

	if p.cur.Type == kwBETWEEN {
		// BETWEEN <bound> AND <bound>
		p.advance()
		var err error
		frame.Start, err = p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAND); err != nil {
			frame.Loc.End = p.prevEnd()
			return frame, nil
		}
		frame.End, err = p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
	} else {
		// Short syntax: <window frame preceding> (implies AND CURRENT ROW)
		var err error
		frame.Start, err = p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
	}

	frame.Loc.End = p.prevEnd()
	return frame, nil
}

// parseWindowFrameBound parses a single window frame bound.
//
//	<window frame bound> ::=
//	  {   UNBOUNDED PRECEDING
//	    | UNBOUNDED FOLLOWING
//	    | <unsigned_value_specification> PRECEDING
//	    | <unsigned_value_specification> FOLLOWING
//	    | CURRENT ROW
//	  }
func (p *Parser) parseWindowFrameBound() (*nodes.WindowBound, error) {
	bound := &nodes.WindowBound{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	switch {
	case p.cur.Type == kwUNBOUNDED:
		p.advance()
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundUnboundedPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundUnboundedFollowing
		}
	case p.cur.Type == kwCURRENT:
		p.advance()
		p.matchIdentCI("ROW")
		bound.Type = nodes.BoundCurrentRow
	default:
		// <unsigned_value_specification> PRECEDING | FOLLOWING
		var err error
		bound.Offset, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundFollowing
		}
	}

	bound.Loc.End = p.prevEnd()
	return bound, nil
}
