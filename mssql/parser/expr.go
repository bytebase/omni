// Package parser - expr.go implements T-SQL expression parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

var (
	jsonArrayOption   = newOptionSet().withIdents("ARRAY")
	jsonWrapperOption = newOptionSet().withIdents("WRAPPER")
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

// parseScalarExpr parses an expression in a value position, where the search
// condition around it does not reach: a boolean is no value, so SQL Server
// refuses IF ABS(1 AND 2) = 1 and IF CAST(1 AND 2 AS int) = 1 with Msg 4145
// and IF CASE 1 WHEN 1 AND 2 THEN 1 END = 1 with Msg 156 (measured). A
// construct that opens value slots clears the depth for them, as
// parseSelectStmtAfterWith does for a subquery's SELECT.
func (p *Parser) parseScalarExpr() (nodes.ExprNode, error) {
	return p.scalar(p.parseExpr)
}

// scalar runs parse with the search condition cleared: an operand of a
// predicate or of an operator is a value, so WHERE a = (1 = 1),
// WHERE a + (a = 1) = 2 and WHERE a IN (1 = 1) are refused (measured).
func (p *Parser) scalar(parse func() (nodes.ExprNode, error)) (nodes.ExprNode, error) {
	if p.searchCondDepth == 0 {
		return parse()
	}
	saved := p.searchCondDepth
	p.searchCondDepth = 0
	defer func() { p.searchCondDepth = saved }()
	return parse()
}

// isConditionParen reports whether e is a parenthesized search condition,
// such as (a = 1) in WHERE (a = 1) AND a > 0. It is a predicate, not a
// value, so no operator takes it: WHERE (a = 1) = 1 is refused (measured).
func isConditionParen(e nodes.ExprNode) bool {
	paren, ok := e.(*nodes.ParenExpr)
	return ok && procIsCondition(paren.Expr)
}

// parseOr parses OR expressions (lowest precedence).
//
//	or_expr = and_expr { OR and_expr }
//
// SQL Server reads AND, OR and NOT only in a search condition: a boolean is no
// value there, so SELECT 1 OR 2 and UPDATE a SET a = 1 AND 2 are Msg 156
// (measured), and the three operators stop an expression outside one.
func (p *Parser) parseOr() (nodes.ExprNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwOR && p.inSearchCondition() {
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
	for p.cur.Type == kwAND && p.inSearchCondition() {
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	if p.cur.Type == kwNOT && p.inSearchCondition() {
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
	// SQL Server reads a comparison, IS NULL, BETWEEN, IN and LIKE only in a
	// search condition (booleanExpression, TSql160.g): SELECT 1 = 1,
	// SELECT 1 IS NULL and UPDATE a SET a = 1 = 1 are refused (measured).
	// Elsewhere the operators end the expression.
	if !p.inSearchCondition() || isConditionParen(left) {
		return left, nil
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
		// IS [NOT] DISTINCT FROM, which SQL Server 2022 takes only in a
		// search condition: IF 1 IS DISTINCT FROM 2 is accepted and
		// SELECT 1 IS DISTINCT FROM 2 is an error.
		if p.cur.Type == kwDISTINCT && p.inSearchCondition() {
			p.advance() // consume DISTINCT
			if _, err := p.expect(kwFROM); err != nil {
				return nil, err
			}
			testType := nodes.IsDistinctFrom
			if not {
				testType = nodes.IsNotDistinctFrom
			}
			switch {
			case p.cur.Type == kwNULL:
				// IS [NOT] DISTINCT FROM NULL
				// (booleanExpressionNullPredicateFromDistinctPredicate)
				p.advance()
				return &nodes.IsExpr{
					Expr:     left,
					TestType: testType,
					Value:    &nodes.Literal{Type: nodes.LitNull, Loc: nodes.Loc{Start: p.prev.Loc, End: p.prevEnd()}},
					Loc:      nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
				}, nil
			case (p.cur.Type == kwALL || p.cur.Type == kwSOME || p.cur.Type == kwANY) && p.peekNext().Type == '(':
				// IS [NOT] DISTINCT FROM { ALL | SOME | ANY } ( <subquery> )
				// (subqueryComparisonDistinctPredicate)
				quantifier := strings.ToUpper(p.cur.Str)
				p.advance() // consume ALL/SOME/ANY
				p.advance() // consume (
				subquery, err := p.parseSubquerySelect()
				if err != nil {
					return nil, err
				}
				if subquery == nil {
					return nil, p.unexpectedToken()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				distinct := nodes.DistinctFrom
				if not {
					distinct = nodes.DistinctNotFrom
				}
				return &nodes.SubqueryComparisonExpr{
					Left:       left,
					Quantifier: quantifier,
					Distinct:   distinct,
					Subquery:   subquery,
					Loc:        nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
				}, nil
			}
			// IS [NOT] DISTINCT FROM <expression>
			// (booleanExpressionDistinctPredicate)
			value, err := p.scalar(p.parseAddition)
			if err != nil {
				return nil, err
			}
			if value == nil {
				return nil, p.unexpectedToken()
			}
			return &nodes.IsExpr{
				Expr:     left,
				TestType: testType,
				Value:    value,
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
		low, err := p.scalar(p.parseAddition)
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
		high, err := p.scalar(p.parseAddition)
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
			sub, err := p.parseSubquerySelect()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
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
			expr, err := p.parseScalarExpr()
			if err != nil {
				return nil, err
			}
			// IN's list holds expressions, and a star is none: IF 1 IN (*)
			// SELECT 1 is Msg 102 (measured).
			if expr == nil {
				return nil, p.unexpectedToken()
			}
			if _, star := expr.(*nodes.StarExpr); star {
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
		pattern, err := p.scalar(p.parseAddition)
		if err != nil {
			return nil, err
		}
		if pattern == nil {
			return nil, p.unexpectedToken()
		}
		var escape nodes.ExprNode
		braced := false
		if p.cur.Type == kwESCAPE {
			p.advance()
			escape, err = p.scalar(p.parseAddition)
			if err != nil {
				return nil, err
			}
			// ESCAPE takes an expression: IF 1 LIKE 1 ESCAPE SELECT 1 is
			// Msg 156 (measured).
			if escape == nil {
				return nil, p.syntaxErrorAtCur()
			}
		} else if p.cur.Type == '{' {
			// likePredicate takes the ODBC escape { ESCAPE <expr> } beside
			// the bare escapeExpression (bnf/TSql170.g).
			p.advance() // consume {
			if _, err := p.expect(kwESCAPE); err != nil {
				return nil, err
			}
			escape, err = p.scalar(p.parseAddition)
			if err != nil {
				return nil, err
			}
			if escape == nil {
				return nil, p.unexpectedToken()
			}
			if _, err := p.expect('}'); err != nil {
				return nil, err
			}
			braced = true
		}
		return &nodes.LikeExpr{
			Expr:         left,
			Pattern:      pattern,
			Escape:       escape,
			Not:          not,
			EscapeBraced: braced,
			Loc:          nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
		}, nil
	}

	// Comparison operators. In a search condition SQL Server also reads the
	// operators ScriptDOM's comparisonOperator writes with a space between
	// their characters: ! =, ! <, ! >, < =, < > and > =. The lexer joins them
	// only when they touch, and elsewhere they are errors
	// (SELECT 1 ! = 2 is refused).
	op, ok := p.splitComparisonOp()
	tokens := 2
	if !ok {
		op, ok = comparisonOp(p.cur.Type)
		tokens = 1
	}
	if ok {
		loc := p.pos()
		for i := 0; i < tokens; i++ {
			p.advance()
		}

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
				subquery, err := p.parseSubquerySelect()
				if err != nil {
					return nil, err
				}
				if subquery == nil {
					return nil, p.unexpectedToken()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				return &nodes.SubqueryComparisonExpr{
					Left:       left,
					Op:         op,
					Quantifier: quantifier,
					Subquery:   subquery,
					Loc:        nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()},
				}, nil
			}
		}

		right, err := p.scalar(p.parseAddition)
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

// splitComparisonOp reads a comparison operator whose two characters the
// lexer made two tokens, which SQL Server takes in a search condition:
// ! =, ! <, ! >, < =, < > and > =. It reads the two tokens without advancing
// and reports false otherwise.
func (p *Parser) splitComparisonOp() (nodes.BinaryOp, bool) {
	if !p.inSearchCondition() {
		return 0, false
	}
	return p.splitComparisonTokens()
}

// splitComparisonTokens reads the two-token comparison operators of
// splitComparisonOp anywhere: comparisonOperator (TSql160.g:31264) writes
// each as two tokens, Bang ( EqualsSign | LessThan | GreaterThan ),
// LessThan ( EqualsSign | GreaterThan )? and GreaterThan ( EqualsSign )?.
func (p *Parser) splitComparisonTokens() (nodes.BinaryOp, bool) {
	first, second := p.cur.Type, p.peekNext().Type
	var op nodes.BinaryOp
	switch first {
	case '!':
		switch second {
		case '=':
			op = nodes.BinOpNeq
		case '<':
			op = nodes.BinOpNotLt
		case '>':
			op = nodes.BinOpNotGt
		default:
			return 0, false
		}
	case '<':
		switch second {
		case '=':
			op = nodes.BinOpLte
		case '>':
			op = nodes.BinOpNeq
		default:
			return 0, false
		}
	case '>':
		if second != '=' {
			return 0, false
		}
		op = nodes.BinOpGte
	default:
		return 0, false
	}
	return op, true
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

// parseAddition parses addition, subtraction, the bitwise operators, and AT
// TIME ZONE, which share one precedence level and associate to the left.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/operator-precedence-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/at-time-zone-transact-sql
//
//	addition_expr = multiplication_expr
//	    { ('+' | '-' | '&' | '|' | '^') multiplication_expr
//	    | AT TIME ZONE multiplication_expr }
//
// SQL Server and ScriptDOM read d + 1 AT TIME ZONE 'UTC' as the conversion of
// d + 1, and d AT TIME ZONE 'UTC' + 1 as a sum whose left operand is the
// conversion. COLLATE is not at this level: it binds to the operand before
// it (parseCollatedPrimary).
func (p *Parser) parseAddition() (nodes.ExprNode, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}
	for {
		if isConditionParen(left) {
			return left, nil
		}
		if p.cur.Type == kwAT {
			if !p.atTimeZoneAhead() {
				return left, nil
			}
			left, err = p.parseAtTimeZoneExpr(left)
			if err != nil {
				return nil, err
			}
			continue
		}
		if p.cur.Type != '+' && p.cur.Type != '-' && p.cur.Type != '&' && p.cur.Type != '|' && p.cur.Type != '^' {
			return left, nil
		}
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
		right, err := p.scalar(p.parseMultiplication)
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
}

// atTimeZoneAhead reports whether the current AT starts AT TIME ZONE. AT is
// also a column alias, so AT followed by anything but TIME ZONE is left to
// the caller: SELECT x AT FROM t names the column AT, and SQL Server and
// ScriptDOM reject SELECT x AT TIME y FROM t at TIME, the token after the
// alias.
func (p *Parser) atTimeZoneAhead() bool {
	return p.cur.Type == kwAT && p.peekNext().Type == kwTIME && p.peekSecond().Type == kwZONE
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
	collation, err := p.parseCollationName()
	if err != nil {
		return nil, err
	}
	node := &nodes.CollateExpr{
		Expr:      expr,
		Collation: collation,
		Loc:       nodes.Loc{Start: exprLocStart(expr, loc), End: p.prevEnd()},
	}
	node.Loc.End = p.prevEnd()
	return node, nil
}

// parseCollationName parses the collation name after COLLATE: a Windows or
// SQL collation name, or database_default. SQL Server takes only a regular
// identifier there, so a bracketed or double-quoted name, which the lexer
// also returns as tokIDENT, is rejected (Msg 102), as are keywords and
// strings.
func (p *Parser) parseCollationName() (string, error) {
	if !p.isIdentLike() || p.source[p.cur.Loc] == '[' || p.source[p.cur.Loc] == '"' {
		return "", p.unexpectedToken()
	}
	name := p.cur.Str
	p.advance()
	return name, nil
}

// parseAtTimeZoneExpr parses a postfix AT TIME ZONE expression. The current
// token is the AT of AT TIME ZONE (atTimeZoneAhead).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/at-time-zone-transact-sql
//
//	inputdate AT TIME ZONE timezone
//
// Can be chained: expr AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'
func (p *Parser) parseAtTimeZoneExpr(expr nodes.ExprNode) (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume AT
	p.advance() // consume TIME
	p.advance() // consume ZONE
	tz, err := p.scalar(p.parseMultiplication)
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
//	multiplication_expr = unary_expr { ('*' | '/' | '%' | '<<' | '>>') unary_expr }
//
// The shift operators of SQL Server 2022 bind as tightly as '*'
// (expressionMultiplication's LeftShift and RightShift, TSql160.g).
func (p *Parser) parseMultiplication() (nodes.ExprNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for !isConditionParen(left) && (p.cur.Type == '*' || p.cur.Type == '/' || p.cur.Type == '%' || p.cur.Type == tokSHIFTLEFT || p.cur.Type == tokSHIFTRIGHT) {
		loc := p.pos()
		var op nodes.BinaryOp
		switch p.cur.Type {
		case '*':
			op = nodes.BinOpMul
		case '/':
			op = nodes.BinOpDiv
		case '%':
			op = nodes.BinOpMod
		case tokSHIFTLEFT:
			op = nodes.BinOpShiftLeft
		case tokSHIFTRIGHT:
			op = nodes.BinOpShiftRight
		}
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		right, err := p.scalar(p.parseUnary)
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
//	unary_expr = ('+' | '-' | '~') unary_expr | collated_primary
func (p *Parser) parseUnary() (nodes.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
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
		operand, err := p.scalar(p.parseUnary)
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
	return p.parseCollatedPrimary()
}

// parseCollatedPrimary parses a primary expression and the COLLATE clause
// that may follow it.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/collations
//
//	collated_primary = primary_expr [ COLLATE collation_name ]
//
// COLLATE binds tighter than any operator and collates the operand before
// it: 'a' COLLATE x + 'b' collates 'a', 'a' + 'b' COLLATE x collates 'b', and
// -x COLLATE x negates the collated x. SQL Server accepts COLLATE before a
// binary operator, and ScriptDOM keeps the collation on the operand's
// PrimaryExpression. An operand takes one COLLATE: x COLLATE a COLLATE b is
// Msg 156 at the second COLLATE, while (x COLLATE a) COLLATE b collates a
// new, parenthesized operand.
func (p *Parser) parseCollatedPrimary() (nodes.ExprNode, error) {
	expr, err := p.parsePrimary()
	if err != nil || expr == nil || isConditionParen(expr) {
		return expr, err
	}
	if p.searchCondDepth > 0 {
		saved := p.searchCondDepth
		p.searchCondDepth = 0
		defer func() { p.searchCondDepth = saved }()
	}
	expr, err = p.parseInstanceMethods(expr)
	if err != nil || expr == nil || p.cur.Type != kwCOLLATE {
		return expr, err
	}
	expr, err = p.parseCollateExpr(expr)
	if err != nil {
		return nil, err
	}
	if p.cur.Type == kwCOLLATE {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
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
	// A primary is a value, and the search condition around it does not
	// reach inside: WHERE CASE WHEN a = 1 THEN a = 1 END = 1 and
	// WHERE IIF(a = 1, a = 1, 0) = 1 are refused (measured). A parenthesis
	// may hold a search condition of its own (WHERE (a = 1) AND a > 0), and
	// UPDATE ( column ), TSEQUAL and MATCH are predicates that check the
	// depth themselves.
	if p.searchCondDepth > 0 && !p.searchConditionPrimaryAhead() {
		saved := p.searchCondDepth
		p.searchCondDepth = 0
		defer func() { p.searchCondDepth = saved }()
	}
	// UPDATE ( column ) tests a column in a DML trigger's search condition
	// (updateCall, TSql160.g). SQL Server refuses it elsewhere: IF UPDATE(a)
	// SELECT 1 is rejected in compat.tsv, and accepted in a trigger's body
	// (measured, mssql/catalog/testdata/programs/simulate.tsv).
	if p.cur.Type == kwUPDATE && p.dmlTrigger && p.inSearchCondition() && p.peekNext().Type == '(' {
		return p.parseUpdateCall()
	}
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
	case tokMONEY:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitMoney,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	case tokPSEUDOCOL:
		return p.parsePseudoColumn()
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
			sub, err := p.parseSubquerySelect()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.SubqueryExpr{
				Query: sub,
				Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
			}, nil
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		// Parentheses hold an expression, and a star is no expression outside
		// a select list: SELECT (), SELECT (*) and IF 1 IN (*) SELECT 1 are
		// Msg 102 (measured).
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, star := expr.(*nodes.StarExpr); star {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
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
	case kwTSEQUAL:
		// tsEqualCall (TSql160.g:30641): TSEQUAL ( expression , expression ),
		// a predicate of a search condition and nothing else. SQL Server
		// refuses it in a select list (Msg 156), with other than two
		// arguments (Msg 102, measured), and reads each argument as an
		// expression rather than a condition: IF TSEQUAL(1 AND 2, 3) SELECT 1
		// is Msg 156 too.
		if !p.inSearchCondition() || p.peekNext().Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		loc := p.pos()
		name := p.advance() // consume TSEQUAL
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		savedDepth := p.searchCondDepth
		p.searchCondDepth = 0
		first, err := p.parseExpr()
		if err == nil {
			if first == nil {
				err = p.syntaxErrorAtCur()
			} else {
				_, err = p.expect(',')
			}
		}
		var second nodes.ExprNode
		if err == nil {
			second, err = p.parseExpr()
			if err == nil && second == nil {
				err = p.syntaxErrorAtCur()
			}
		}
		p.searchCondDepth = savedDepth
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// Neither DISTINCT, WITHIN GROUP nor OVER: this is a predicate of its
		// own, not the generic function call.
		return &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: name.Str, Loc: nodes.Loc{Start: loc, End: name.End}},
			Args: &nodes.List{Items: []nodes.Node{first, second}},
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case kwTRY_CAST:
		if p.peekNext().Type == '(' {
			return p.parseTryCast()
		}
		return p.parseIdentExpr()
	case kwTRY_CONVERT:
		return p.parseTryConvert()
	case kwLEFT, kwRIGHT:
		if p.peekNext().Type == '(' {
			return p.parseIdentExpr()
		}
		return nil, nil
	// CURRENT_DATE and CURRENT_TIME are ScriptDOM's parameterlessCall for
	// another engine: SQL Server 2022 refuses both with Msg 156 (measured).
	case kwCURRENT_TIMESTAMP, kwCURRENT_USER, kwSESSION_USER,
		kwSYSTEM_USER, kwUSER:
		return p.parseNiladicKeywordFunc()
	case '.':
		// A name may start with a dot, which names an empty first part, and
		// $PARTITION's call follows one: SELECT .$PARTITION.a(1) is accepted
		// (measured).
		loc := p.pos()
		p.advance() // consume .
		if p.cur.Type != tokPSEUDOCOL {
			return nil, p.syntaxErrorAtCur()
		}
		expr, err := p.parsePseudoColumn()
		if err != nil {
			return nil, err
		}
		if fc, ok := expr.(*nodes.FuncCallExpr); ok {
			fc.Loc.Start = loc
			if fc.Name != nil {
				fc.Name.Loc.Start = loc
			}
		}
		return expr, nil
	case '{':
		// ODBC escape: { fn <function> } or { d|t|ts '<literal>' }
		return p.parseOdbcEscape()
	case kwIDENTITYCOL, kwROWGUIDCOL:
		// IDENTITYCOL / ROWGUIDCOL reference the identity / rowguid column of a
		// table without naming it. Keyword-only column references; no parens.
		tok := p.advance()
		return &nodes.ColumnRef{
			Column: tok.Str,
			Loc:    nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
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

// searchConditionPrimaryAhead reports whether the primary at the current
// token may read the search condition around it: a parenthesis, which may
// hold one, and the predicates UPDATE ( column ), TSEQUAL ( ... ) and
// MATCH ( ... ).
func (p *Parser) searchConditionPrimaryAhead() bool {
	switch p.cur.Type {
	case '(', kwUPDATE, kwTSEQUAL:
		return true
	}
	return p.peekNext().Type == '(' && p.isValidOptionToken(graphMatchWord, p.cur)
}

// parseNiladicKeywordFunc handles SQL-standard niladic functions that are
// CoreKeywords and take no parentheses: CURRENT_TIMESTAMP, CURRENT_USER, etc.
func (p *Parser) parseNiladicKeywordFunc() (nodes.ExprNode, error) {
	tok := p.advance()
	return &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Object: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.prevEnd()}},
		Loc:  nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
	}, nil
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
	expr, err := p.parseScalarExpr()
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	expr, err := p.parseScalarExpr()
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	expr, err := p.parseRequiredExpr()
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
		style, err = p.parseRequiredExpr()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	expr, err := p.parseRequiredExpr()
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
		style, err = p.parseRequiredExpr()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.TryConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseRequiredExpr parses an expression where the grammar requires one, as
// after a comma in an argument list: parseExpr returns nil when no
// expression starts here.
func (p *Parser) parseRequiredExpr() (nodes.ExprNode, error) {
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return e, nil
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
		arg, err = p.parseScalarExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	// CASE takes one WHEN at least, THEN after each, a value after ELSE,
	// and END (caseExpression, TSql160.g; Msg 102 / 156 otherwise).
	if p.cur.Type != kwWHEN && !p.collectMode() {
		return nil, p.syntaxErrorAtCur()
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
		// The searched WHEN is a condition, the simple one's a value.
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
		// A searched CASE's WHEN is a search condition, and SQL Server refuses
		// an expression of no boolean type there: CASE WHEN 1 THEN 1 END is
		// Msg 4145 (measured).
		if arg == nil && !procIsCondition(cond) {
			return nil, p.newParseError(wloc, "an expression of non-boolean type specified in a context where a condition is expected")
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		// Completion: after CASE WHEN x THEN → columnref, func_name
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		result, err := p.parseScalarExpr()
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
		if elseExpr == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}

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
//	COALESCE '(' expr ',' expr { ',' expr } ')'
//
// COALESCE takes two arguments at least (coalesceExpression, TSql160.g:
// expression (Comma expression)+); COALESCE(1) is Msg 102 (measured,
// mssql/catalog/testdata/queries/simulate.tsv).
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
		arg, err := p.parseScalarExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			if len(args) < 2 {
				return nil, p.unexpectedToken()
			}
			break
		}
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	left, err := p.parseScalarExpr()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	right, err := p.parseRequiredExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	// IIF's first argument is a search condition: IIF(1 = 1 AND 2 = 2, 1, 0)
	// is accepted and IIF(1, 1, 0) is Msg 4145 (measured).
	p.enterSearchCondition()
	cond, err := p.parseExpr()
	p.leaveSearchCondition()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.unexpectedToken()
	}
	if !procIsCondition(cond) {
		return nil, p.newParseError(nodes.NodeLoc(cond).Start, "an expression of non-boolean type specified in a context where a condition is expected")
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	trueVal, err := p.parseRequiredExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	falseVal, err := p.parseRequiredExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	query, err := p.parseSubquerySelect()
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, p.unexpectedToken()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	n := &nodes.NextValueForExpr{
		Sequence: ref,
		Loc:      nodes.Loc{Start: loc, End: p.prevEnd()},
	}
	if p.cur.Type == kwOVER {
		over, err := p.parseOverClause()
		if err != nil {
			return nil, err
		}
		n.Over = over
		n.Loc.End = p.prevEnd()
	}
	return n, nil
}

func (p *Parser) parseParseExpr(name string, loc int) (nodes.ExprNode, error) {
	p.advance() // consume (
	n := &nodes.ParseExpr{
		Try: strings.EqualFold(name, "TRY_PARSE"),
		Loc: nodes.Loc{Start: loc, End: -1},
	}
	expr, err := p.parseScalarExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.unexpectedToken()
	}
	n.Expr = expr
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	n.DataType = dt
	if p.cur.Type == kwUSING {
		p.advance()
		culture, err := p.parseScalarExpr()
		if err != nil {
			return nil, err
		}
		if culture == nil {
			return nil, p.unexpectedToken()
		}
		n.Culture = culture
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	n.Loc.End = p.prevEnd()
	return n, nil
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
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
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
		// LEFT, RIGHT and TRIM take an argument, and SQL Server refuses the
		// empty list: SELECT LEFT() is Msg 102 and SELECT TRIM() is Msg 189
		// (measured).
		switch strings.ToUpper(name) {
		case "LEFT", "RIGHT", "TRIM":
			return nil, p.syntaxErrorAtCur()
		}
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

	if strings.EqualFold(name, "TRIM") {
		return p.parseTrimFuncCall(fc)
	}
	if isJsonFunction(name) {
		return p.parseJsonFuncCall(fc, name)
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	argIndex := 0
	args, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
		if p.collectMode() {
			p.addExpressionCandidates()
			return nil, errCollecting
		}
		if argIndex == 0 && isDatePartFunction(name) && p.isIdentLike() {
			tok := p.advance()
			argIndex++
			return &nodes.DatePart{
				Name: tok.Str,
				Loc:  nodes.Loc{Start: tok.Loc, End: tok.End},
			}, nil
		}
		arg, err := p.parseScalarExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		argIndex++
		return arg, nil
	})
	if err != nil {
		return nil, err
	}
	fc.Args = &nodes.List{Items: args}

	// LEFT and RIGHT take exactly two arguments, and SQL Server counts them
	// while parsing: SELECT LEFT(1) is Msg 174 (measured).
	if n := strings.ToUpper(name); n == "LEFT" || n == "RIGHT" {
		if len(args) != 2 {
			return nil, &ParseError{Message: "the " + strings.ToLower(n) + " function requires 2 arguments", Position: loc}
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()

	p.parseNullTreatment(fc, name)

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

func (p *Parser) parseTrimFuncCall(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	if trimOption, ok := p.consumeTrimOption(); ok {
		fc.TrimOption = trimOption
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if first == nil {
		return nil, p.unexpectedToken()
	}
	items := []nodes.Node{first}
	if p.cur.Type == kwFROM {
		p.advance()
		second, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if second == nil {
			return nil, p.unexpectedToken()
		}
		items = append(items, second)
	}
	fc.Args = &nodes.List{Items: items}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()
	return fc, nil
}

func (p *Parser) consumeTrimOption() (string, bool) {
	if !p.isKeywordOrIdent() {
		return "", false
	}
	switch strings.ToUpper(p.cur.Str) {
	case "LEADING", "TRAILING", "BOTH":
		option := strings.ToUpper(p.cur.Str)
		p.advance()
		return option, true
	default:
		return "", false
	}
}

func (p *Parser) parseNullTreatment(fc *nodes.FuncCallExpr, name string) {
	if !isNullTreatmentFunction(name) || !p.isKeywordOrIdent() {
		return
	}
	treatment := strings.ToUpper(p.cur.Str)
	if treatment != "IGNORE" && treatment != "RESPECT" {
		return
	}
	next := p.peekNext()
	if next.Type != kwNULLS {
		return
	}
	p.advance()
	p.advance()
	fc.NullTreatment = treatment
	fc.Loc.End = p.prevEnd()
}

func isDatePartFunction(name string) bool {
	switch strings.ToUpper(name) {
	case "DATEADD", "DATEDIFF", "DATEDIFF_BIG", "DATEPART", "DATENAME", "DATETRUNC", "DATE_BUCKET":
		return true
	default:
		return false
	}
}

func isNullTreatmentFunction(name string) bool {
	switch strings.ToUpper(name) {
	case "FIRST_VALUE", "LAST_VALUE", "LAG", "LEAD":
		return true
	default:
		return false
	}
}

func isJsonFunction(name string) bool {
	switch strings.ToUpper(name) {
	case "JSON_OBJECT", "JSON_ARRAY", "JSON_VALUE", "JSON_QUERY":
		return true
	default:
		return false
	}
}

func (p *Parser) parseJsonFuncCall(fc *nodes.FuncCallExpr, name string) (nodes.ExprNode, error) {
	switch strings.ToUpper(name) {
	case "JSON_OBJECT":
		if err := p.parseJsonObjectArgs(fc); err != nil {
			return nil, err
		}
	case "JSON_ARRAY":
		if err := p.parseJsonArrayArgs(fc); err != nil {
			return nil, err
		}
	case "JSON_VALUE":
		if err := p.parseJsonValueArgs(fc); err != nil {
			return nil, err
		}
	case "JSON_QUERY":
		if err := p.parseJsonQueryArgs(fc); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()
	return fc, nil
}

func (p *Parser) parseJsonObjectArgs(fc *nodes.FuncCallExpr) error {
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF && !p.isJsonNullClauseStart() && !p.isReturningStart() {
		loc := p.pos()
		key, err := p.parseExpr()
		if err != nil {
			return err
		}
		if key == nil {
			return p.unexpectedToken()
		}
		if _, err := p.expect(':'); err != nil {
			return err
		}
		value, err := p.parseExpr()
		if err != nil {
			return err
		}
		if value == nil {
			return p.unexpectedToken()
		}
		items = append(items, &nodes.JsonKeyValueExpr{
			Key:   key,
			Value: value,
			Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	fc.Args = &nodes.List{Items: items}
	if err := p.parseJsonNullClause(fc); err != nil {
		return err
	}
	return p.parseJsonReturningClause(fc)
}

func (p *Parser) parseJsonArrayArgs(fc *nodes.FuncCallExpr) error {
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF && !p.isJsonNullClauseStart() && !p.isReturningStart() {
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		if expr == nil {
			return p.unexpectedToken()
		}
		items = append(items, expr)
		if _, ok := p.match(','); !ok {
			break
		}
		// A comma needs a value after it: JSON_ARRAY(1,) is refused.
		if p.cur.Type == ')' {
			return p.unexpectedToken()
		}
	}
	fc.Args = &nodes.List{Items: items}
	if err := p.parseJsonNullClause(fc); err != nil {
		return err
	}
	return p.parseJsonReturningClause(fc)
}

func (p *Parser) parseJsonValueArgs(fc *nodes.FuncCallExpr) error {
	var items []nodes.Node
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil {
		return p.unexpectedToken()
	}
	items = append(items, expr)
	if _, ok := p.match(','); ok {
		path, err := p.parseExpr()
		if err != nil {
			return err
		}
		if path == nil {
			return p.unexpectedToken()
		}
		items = append(items, path)
	}
	fc.Args = &nodes.List{Items: items}
	return p.parseJsonReturningClause(fc)
}

func (p *Parser) parseJsonQueryArgs(fc *nodes.FuncCallExpr) error {
	var items []nodes.Node
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil {
		return p.unexpectedToken()
	}
	items = append(items, expr)
	if _, ok := p.match(','); ok {
		path, err := p.parseExpr()
		if err != nil {
			return err
		}
		if path == nil {
			return p.unexpectedToken()
		}
		items = append(items, path)
	}
	fc.Args = &nodes.List{Items: items}
	return p.parseJsonWithArrayWrapperClause(fc)
}

func (p *Parser) parseJsonNullClause(fc *nodes.FuncCallExpr) error {
	if !p.isJsonNullClauseStart() {
		return nil
	}
	clause := strings.ToUpper(p.cur.Str)
	p.advance()
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	if _, err := p.expect(kwNULL); err != nil {
		return err
	}
	fc.JsonNullClause = clause
	return nil
}

func (p *Parser) parseJsonReturningClause(fc *nodes.FuncCallExpr) error {
	if !p.isReturningStart() {
		return nil
	}
	p.advance()
	dt, err := p.parseDataType()
	if err != nil {
		return err
	}
	if dt == nil {
		return p.unexpectedToken()
	}
	fc.ReturnType = dt
	return nil
}

func (p *Parser) parseJsonWithArrayWrapperClause(fc *nodes.FuncCallExpr) error {
	if p.cur.Type != kwWITH {
		return nil
	}
	p.advance()
	if _, err := p.expectOption(jsonArrayOption); err != nil {
		return err
	}
	if _, err := p.expectOption(jsonWrapperOption); err != nil {
		return err
	}
	fc.WithArrayWrapper = true
	return nil
}

func (p *Parser) isJsonNullClauseStart() bool {
	if p.cur.Type != kwNULL && p.cur.Type != kwABSENT {
		return false
	}
	return p.peekNext().Type == kwON
}

func (p *Parser) isReturningStart() bool {
	return p.cur.Type == kwRETURNING
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
		return nil, err
	}
	if _, err := p.expect(kwBY); err != nil {
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	spec, err := p.parseWindowSpecBody(true)
	if err != nil {
		return nil, err
	}
	over.WindowName, over.PartitionBy, over.OrderBy, over.WindowFrame = spec.ref, spec.partitionBy, spec.orderBy, spec.frame

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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

	// SQL Server has ROWS and RANGE, not GROUPS (Msg 102).
	switch p.cur.Type {
	case kwROWS:
		frame.Type = nodes.FrameRows
	case kwRANGE:
		frame.Type = nodes.FrameRange
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// Completion: after ROWS/RANGE → BETWEEN, UNBOUNDED, CURRENT
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
			return nil, err
		}
		frame.End, err = p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
	} else {
		// One bound, which SQL Server's parser takes but for UNBOUNDED
		// FOLLOWING (Msg 102).
		var err error
		frame.Start, err = p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
		if frame.Start.Type == nodes.BoundUnboundedFollowing {
			return nil, &ParseError{Message: "syntax error at or near \"FOLLOWING\"", Position: p.prev.Loc}
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
//	    | unsigned_integer PRECEDING
//	    | unsigned_integer FOLLOWING
//	    | CURRENT ROW
//	  }
//
// The offset is an integer literal: SQL Server refuses a variable or an
// expression there (Msg 102; windowFrameBound, TSql160.g).
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
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	case p.cur.Type == kwCURRENT:
		p.advance()
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		bound.Type = nodes.BoundCurrentRow
	case p.cur.Type == tokICONST:
		bound.Offset = &nodes.Literal{Type: nodes.LitInteger, Ival: p.cur.Ival, Str: p.cur.Str, Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundFollowing
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	bound.Loc.End = p.prevEnd()
	return bound, nil
}

// windowSpec is the body of an OVER ( ... ) or a WINDOW definition.
type windowSpec struct {
	ref         string
	partitionBy *nodes.List
	orderBy     *nodes.List
	frame       *nodes.WindowFrame
}

// parseWindowSpecBody parses
//
//	[ window_name ] [ PARTITION BY expr [ ,...n ] ] [ ORDER BY order_item [ ,...n ] ] [ frame ]
//
// up to the closing parenthesis, which it leaves. SQL Server takes a frame
// only after a window name, PARTITION BY or ORDER BY (Msg 102 on OVER (ROWS
// ...)), and BY and a non-empty list after PARTITION and ORDER. In OVER, a
// window name needs something after it: OVER (w) is refused (Msg 102),
// OVER w is the form.
func (p *Parser) parseWindowSpecBody(over bool) (windowSpec, error) {
	var spec windowSpec
	if p.isIdentLike() && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER && p.cur.Type != kwROWS && p.cur.Type != kwRANGE {
		spec.ref = p.cur.Str
		p.advance()
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return spec, err
		}
		// Completion: after PARTITION BY → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return spec, errCollecting
		}
		list, err := p.parseExprList()
		if err != nil {
			return spec, err
		}
		if list == nil || len(list.Items) == 0 {
			return spec, p.syntaxErrorAtCur()
		}
		spec.partitionBy = list
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return spec, err
		}
		// Completion: after ORDER BY in OVER → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return spec, errCollecting
		}
		list, err := p.parseOrderByList()
		if err != nil {
			return spec, err
		}
		if list == nil || len(list.Items) == 0 {
			return spec, p.syntaxErrorAtCur()
		}
		spec.orderBy = list
	}

	// ROWS | RANGE window frame clause
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		if spec.ref == "" && spec.partitionBy == nil && spec.orderBy == nil {
			return spec, p.syntaxErrorAtCur()
		}
		frame, err := p.parseWindowFrame()
		if err != nil {
			return spec, err
		}
		spec.frame = frame
	}
	if over && spec.ref != "" && spec.partitionBy == nil && spec.orderBy == nil && spec.frame == nil {
		return spec, p.syntaxErrorAtCur()
	}
	return spec, nil
}

// parseUpdateCall parses UPDATE ( column ) as a call of UPDATE with the
// column as its argument.
func (p *Parser) parseUpdateCall() (nodes.ExprNode, error) {
	loc := p.pos()
	nameTok := p.advance() // consume UPDATE
	p.advance()            // consume (
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	col := p.advance()
	arg := &nodes.ColumnRef{Column: col.Str, Loc: nodes.Loc{Start: col.Loc, End: p.prevEnd()}}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Object: "UPDATE", Loc: nodes.Loc{Start: nameTok.Loc, End: nameTok.Loc + len("UPDATE")}},
		Args: &nodes.List{Items: []nodes.Node{arg}},
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseInstanceMethods parses the methods and properties that may follow a
// primary expression: expressionWithClrElements (TSql160.g:31441), primary
// ( . identifier [ ( args ) ] )*. EVENTDATA().value(path, type) is the common
// one (measured, mssql/catalog/testdata/programs/simulate.tsv), and SQL Server
// takes it on a variable, a subquery, a literal, CASE and COALESCE too:
// @x.value('.', 'int'), (SELECT ... FOR XML PATH(”), TYPE).value(...) and
// @x.a.b parse (measured, compat.tsv). A property has no parentheses and nil
// Args. A dotted name's methods are read with the name.
func (p *Parser) parseInstanceMethods(expr nodes.ExprNode) (nodes.ExprNode, error) {
	switch e := expr.(type) {
	case *nodes.FuncCallExpr, *nodes.CastExpr, *nodes.ConvertExpr, *nodes.TryCastExpr, *nodes.TryConvertExpr,
		*nodes.ParenExpr, *nodes.SubqueryExpr, *nodes.VariableRef, *nodes.CaseExpr, *nodes.CoalesceExpr,
		*nodes.NullifExpr, *nodes.IifExpr:
	case *nodes.Literal:
		if e.Type == nodes.LitDefault {
			return expr, nil
		}
	default:
		return expr, nil
	}
	for p.cur.Type == '.' && IsIdentTokenType(p.peekNext().Type) {
		start := nodes.NodeLoc(expr).Start
		p.advance() // consume .
		method := p.cur.Str
		p.advance() // consume the method's name
		mc := &nodes.MethodCallExpr{Expr: expr, Method: method}
		if p.cur.Type == '(' {
			p.advance() // consume (
			args, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
				return p.parseExpr()
			})
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			mc.Args = &nodes.List{Items: args}
		}
		mc.Loc = nodes.Loc{Start: start, End: p.prevEnd()}
		expr = mc
	}
	return expr, nil
}
