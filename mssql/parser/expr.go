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
func (p *Parser) parseExpr() nodes.ExprNode {
	return p.parseOr()
}

// parseOr parses OR expressions (lowest precedence).
//
//	or_expr = and_expr { OR and_expr }
func (p *Parser) parseOr() nodes.ExprNode {
	left := p.parseAnd()
	for p.cur.Type == kwOR {
		loc := p.pos()
		p.advance()
		right := p.parseAnd()
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpOr,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

// parseAnd parses AND expressions.
//
//	and_expr = not_expr { AND not_expr }
func (p *Parser) parseAnd() nodes.ExprNode {
	left := p.parseNot()
	for p.cur.Type == kwAND {
		loc := p.pos()
		p.advance()
		right := p.parseNot()
		left = &nodes.BinaryExpr{
			Op:    nodes.BinOpAnd,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

// parseNot parses NOT expressions.
//
//	not_expr = NOT not_expr | comparison_expr
func (p *Parser) parseNot() nodes.ExprNode {
	if p.cur.Type == kwNOT {
		loc := p.pos()
		p.advance()
		operand := p.parseNot()
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryNot,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc},
		}
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
func (p *Parser) parseComparison() nodes.ExprNode {
	left := p.parseAddition()
	if left == nil {
		return nil
	}

	// IS [NOT] NULL
	if p.cur.Type == kwIS {
		loc := p.pos()
		p.advance()
		not := false
		if p.cur.Type == kwNOT {
			not = true
			p.advance()
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
				Loc:      nodes.Loc{Start: loc},
			}
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
		low := p.parseAddition()
		if _, err := p.expect(kwAND); err != nil {
			return left
		}
		high := p.parseAddition()
		return &nodes.BetweenExpr{
			Expr: left,
			Low:  low,
			High: high,
			Not:  not,
			Loc:  nodes.Loc{Start: loc},
		}
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
			return left
		}
		var items []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			expr := p.parseExpr()
			items = append(items, expr)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		_, _ = p.expect(')')
		return &nodes.InExpr{
			Expr: left,
			List: &nodes.List{Items: items},
			Not:  not,
			Loc:  nodes.Loc{Start: loc},
		}
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
		pattern := p.parseAddition()
		var escape nodes.ExprNode
		if p.cur.Type == kwESCAPE {
			p.advance()
			escape = p.parseAddition()
		}
		return &nodes.LikeExpr{
			Expr:    left,
			Pattern: pattern,
			Escape:  escape,
			Not:     not,
			Loc:     nodes.Loc{Start: loc},
		}
	}

	// Comparison operators
	op, ok := comparisonOp(p.cur.Type)
	if ok {
		loc := p.pos()
		p.advance()

		// Check for ANY/SOME/ALL subquery comparison
		if p.cur.Type == kwALL || p.cur.Type == kwSOME || p.cur.Type == kwANY {
			quantifier := strings.ToUpper(p.cur.Str)
			p.advance() // consume ALL/SOME/ANY
			if p.cur.Type == '(' {
				p.advance() // consume (
				subquery := p.parseSelectStmt()
				_, _ = p.expect(')')
				return &nodes.SubqueryComparisonExpr{
					Left:       left,
					Op:         op,
					Quantifier: quantifier,
					Subquery:   subquery,
					Loc:        nodes.Loc{Start: loc},
				}
			}
		}

		right := p.parseAddition()
		return &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}

	return left
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
func (p *Parser) parseAddition() nodes.ExprNode {
	left := p.parseMultiplication()
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
		right := p.parseMultiplication()
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	// Postfix COLLATE / AT TIME ZONE (may be chained)
	for {
		if p.cur.Type == kwCOLLATE {
			left = p.parseCollateExpr(left)
		} else if p.cur.Type == kwAT {
			left = p.parseAtTimeZoneExpr(left)
		} else {
			break
		}
	}
	return left
}

// parseCollateExpr parses a postfix COLLATE operator on an expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/collations
//
//	expr COLLATE { <collation_name> | database_default }
//
//	<collation_name> ::=
//	    { Windows_collation_name } | { SQL_collation_name }
func (p *Parser) parseCollateExpr(expr nodes.ExprNode) nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume COLLATE
	// Collation name is an identifier (may contain underscores, numbers)
	// or the special keyword database_default
	var collation string
	if p.cur.Type == tokIDENT || p.isIdentLike() {
		collation = p.cur.Str
		p.advance()
	} else {
		return expr
	}
	node := &nodes.CollateExpr{
		Expr:      expr,
		Collation: collation,
		Loc:       nodes.Loc{Start: loc},
	}
	node.Loc.End = p.pos()
	return node
}

// parseAtTimeZoneExpr parses a postfix AT TIME ZONE expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/at-time-zone-transact-sql
//
//	inputdate AT TIME ZONE timezone
//
// Can be chained: expr AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'
func (p *Parser) parseAtTimeZoneExpr(expr nodes.ExprNode) nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume AT
	if p.cur.Type != kwTIME {
		// Not AT TIME ZONE — return original expression
		// (AT used in some other context is handled elsewhere)
		return expr
	}
	p.advance() // consume TIME
	if _, err := p.expect(kwZONE); err != nil {
		return expr
	}
	tz := p.parseMultiplication()
	node := &nodes.AtTimeZoneExpr{
		Expr:     expr,
		TimeZone: tz,
		Loc:      nodes.Loc{Start: loc},
	}
	node.Loc.End = p.pos()
	return node
}

// parseMultiplication parses multiplication/division/modulo.
//
//	multiplication_expr = unary_expr { ('*' | '/' | '%') unary_expr }
func (p *Parser) parseMultiplication() nodes.ExprNode {
	left := p.parseUnary()
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
		right := p.parseUnary()
		left = &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: loc},
		}
	}
	return left
}

// parseUnary parses unary operators.
//
//	unary_expr = ('+' | '-' | '~') unary_expr | primary_expr
func (p *Parser) parseUnary() nodes.ExprNode {
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
		operand := p.parseUnary()
		return &nodes.UnaryExpr{
			Op:      op,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc},
		}
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
func (p *Parser) parsePrimary() nodes.ExprNode {
	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitInteger,
			Ival: tok.Ival,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokFCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitFloat,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitString,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokNSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type:    nodes.LitString,
			Str:     tok.Str,
			IsNChar: true,
			Loc:     nodes.Loc{Start: tok.Loc},
		}
	case kwNULL:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitNull,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case kwDEFAULT:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitDefault,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case tokSYSVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc},
		}
	case '*':
		tok := p.advance()
		return &nodes.StarExpr{
			Loc: nodes.Loc{Start: tok.Loc},
		}
	case '(':
		loc := p.pos()
		p.advance()
		expr := p.parseExpr()
		_, _ = p.expect(')')
		return &nodes.ParenExpr{
			Expr: expr,
			Loc:  nodes.Loc{Start: loc},
		}
	case kwCAST:
		return p.parseCast()
	case kwCONVERT:
		return p.parseConvert()
	case kwCASE:
		return p.parseCaseExpr()
	case kwCOALESCE:
		return p.parseCoalesce()
	case kwNULLIF:
		return p.parseNullif()
	case kwIIF:
		return p.parseIif()
	case kwEXISTS:
		return p.parseExists()
	case kwTRY_CAST:
		return p.parseTryCast()
	case kwTRY_CONVERT:
		return p.parseTryConvert()
	case tokIDENT:
		return p.parseIdentExpr()
	default:
		// Many keywords can also be used as identifiers or function names in T-SQL.
		if p.isIdentLike() {
			next := p.peekNext()
			if next.Type == '(' || next.Type == '.' || next.Type == '=' {
				return p.parseIdentExpr()
			}
		}
		return nil
	}
}

// parseCast parses CAST(expr AS data_type).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/cast-and-convert-transact-sql
//
//	CAST '(' expr AS data_type ')'
func (p *Parser) parseCast() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil
	}
	expr := p.parseExpr()
	if _, err := p.expect(kwAS); err != nil {
		return nil
	}
	dt := p.parseDataType()
	_, _ = p.expect(')')
	return &nodes.CastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseTryCast parses TRY_CAST(expr AS data_type).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/try-cast-transact-sql
//
//	TRY_CAST '(' expr AS data_type ')'
func (p *Parser) parseTryCast() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume TRY_CAST
	if _, err := p.expect('('); err != nil {
		return nil
	}
	expr := p.parseExpr()
	if _, err := p.expect(kwAS); err != nil {
		return nil
	}
	dt := p.parseDataType()
	_, _ = p.expect(')')
	return &nodes.TryCastExpr{
		Expr:     expr,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseConvert parses CONVERT(data_type, expr [, style]).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/cast-and-convert-transact-sql
//
//	CONVERT '(' data_type ',' expr [',' style] ')'
func (p *Parser) parseConvert() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CONVERT
	if _, err := p.expect('('); err != nil {
		return nil
	}
	dt := p.parseDataType()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	expr := p.parseExpr()
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		style = p.parseExpr()
	}
	_, _ = p.expect(')')
	return &nodes.ConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseTryConvert parses TRY_CONVERT(data_type, expr [, style]).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/try-convert-transact-sql
//
//	TRY_CONVERT '(' data_type ',' expr [',' style] ')'
func (p *Parser) parseTryConvert() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume TRY_CONVERT
	if _, err := p.expect('('); err != nil {
		return nil
	}
	dt := p.parseDataType()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	expr := p.parseExpr()
	var style nodes.ExprNode
	if _, ok := p.match(','); ok {
		style = p.parseExpr()
	}
	_, _ = p.expect(')')
	return &nodes.TryConvertExpr{
		DataType: dt,
		Expr:     expr,
		Style:    style,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseCaseExpr parses a CASE expression (both simple and searched).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/case-transact-sql
//
//	CASE [expr] WHEN expr THEN expr { WHEN expr THEN expr } [ELSE expr] END
func (p *Parser) parseCaseExpr() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume CASE

	var arg nodes.ExprNode
	// Simple CASE vs searched CASE
	if p.cur.Type != kwWHEN {
		arg = p.parseExpr()
	}

	var whenList []nodes.Node
	for p.cur.Type == kwWHEN {
		wloc := p.pos()
		p.advance() // consume WHEN
		cond := p.parseExpr()
		if _, err := p.expect(kwTHEN); err != nil {
			break
		}
		result := p.parseExpr()
		whenList = append(whenList, &nodes.CaseWhen{
			Condition: cond,
			Result:    result,
			Loc:       nodes.Loc{Start: wloc},
		})
	}

	var elseExpr nodes.ExprNode
	if _, ok := p.match(kwELSE); ok {
		elseExpr = p.parseExpr()
	}

	_, _ = p.expect(kwEND)

	return &nodes.CaseExpr{
		Arg:      arg,
		WhenList: &nodes.List{Items: whenList},
		ElseExpr: elseExpr,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseCoalesce parses COALESCE(expr, expr, ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/coalesce-transact-sql
//
//	COALESCE '(' expr { ',' expr } ')'
func (p *Parser) parseCoalesce() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume COALESCE
	if _, err := p.expect('('); err != nil {
		return nil
	}
	var args []nodes.Node
	for {
		arg := p.parseExpr()
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.CoalesceExpr{
		Args: &nodes.List{Items: args},
		Loc:  nodes.Loc{Start: loc},
	}
}

// parseNullif parses NULLIF(expr, expr).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/nullif-transact-sql
//
//	NULLIF '(' expr ',' expr ')'
func (p *Parser) parseNullif() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume NULLIF
	if _, err := p.expect('('); err != nil {
		return nil
	}
	left := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	right := p.parseExpr()
	_, _ = p.expect(')')
	return &nodes.NullifExpr{
		Left:  left,
		Right: right,
		Loc:   nodes.Loc{Start: loc},
	}
}

// parseIif parses IIF(condition, true_val, false_val).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/logical-functions-iif-transact-sql
//
//	IIF '(' expr ',' expr ',' expr ')'
func (p *Parser) parseIif() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume IIF
	if _, err := p.expect('('); err != nil {
		return nil
	}
	cond := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	trueVal := p.parseExpr()
	if _, err := p.expect(','); err != nil {
		return nil
	}
	falseVal := p.parseExpr()
	_, _ = p.expect(')')
	return &nodes.IifExpr{
		Condition: cond,
		TrueVal:   trueVal,
		FalseVal:  falseVal,
		Loc:       nodes.Loc{Start: loc},
	}
}

// parseExists parses EXISTS(subquery).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/exists-transact-sql
//
//	EXISTS '(' select_stmt ')'
func (p *Parser) parseExists() nodes.ExprNode {
	loc := p.pos()
	p.advance() // consume EXISTS
	if _, err := p.expect('('); err != nil {
		return nil
	}
	query := p.parseStmt()
	_, _ = p.expect(')')
	return &nodes.ExistsExpr{
		Subquery: query,
		Loc:      nodes.Loc{Start: loc},
	}
}

// parseFuncCall parses a function call after the opening paren has been seen.
//
//	func_name '(' [DISTINCT] [expr_list | '*'] ')' [OVER ...]
func (p *Parser) parseFuncCall(name string, loc int) nodes.ExprNode {
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Object: name, Loc: nodes.Loc{Start: loc}},
		Loc:  nodes.Loc{Start: loc},
	}

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		_, _ = p.expect(')')
		if p.cur.Type == kwOVER {
			fc.Over = p.parseOverClause()
		}
		return fc
	}

	if p.cur.Type == ')' {
		p.advance()
		if p.cur.Type == kwOVER {
			fc.Over = p.parseOverClause()
		}
		return fc
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	var args []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		arg := p.parseExpr()
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	fc.Args = &nodes.List{Items: args}

	_, _ = p.expect(')')

	// Check for WITHIN GROUP (ORDER BY ...) clause
	//
	// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/string-agg-transact-sql
	//
	//  STRING_AGG ( expression, separator ) WITHIN GROUP ( ORDER BY order_by_expression [ ASC | DESC ] )
	if p.cur.Type == kwWITHIN {
		fc.Within = p.parseWithinGroupClause()
	}

	// Check for OVER clause
	if p.cur.Type == kwOVER {
		fc.Over = p.parseOverClause()
	}

	return fc
}

// parseWithinGroupClause parses WITHIN GROUP (ORDER BY ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/functions/string-agg-transact-sql
//
//	WITHIN GROUP ( ORDER BY order_by_expression [ ASC | DESC ] [ ,...n ] )
func (p *Parser) parseWithinGroupClause() *nodes.List {
	p.advance() // consume WITHIN
	// Expect GROUP keyword (it's not a reserved keyword, so use matchIdentCI)
	if !p.matchIdentCI("GROUP") {
		return nil
	}
	if _, err := p.expect('('); err != nil {
		return nil
	}
	// Expect ORDER BY
	if _, err := p.expect(kwORDER); err != nil {
		p.expect(')')
		return nil
	}
	if _, err := p.expect(kwBY); err != nil {
		p.expect(')')
		return nil
	}

	var orders []nodes.Node
	for {
		expr := p.parseExpr()
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		orders = append(orders, &nodes.OrderByItem{
			Expr:    expr,
			SortDir: dir,
			Loc:     nodes.Loc{Start: p.pos()},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: orders}
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
func (p *Parser) parseOverClause() *nodes.OverClause {
	loc := p.pos()
	p.advance() // consume OVER
	if _, err := p.expect('('); err != nil {
		return nil
	}

	over := &nodes.OverClause{
		Loc: nodes.Loc{Start: loc},
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over
		}
		var parts []nodes.Node
		for {
			expr := p.parseExpr()
			parts = append(parts, expr)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		over.PartitionBy = &nodes.List{Items: parts}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return over
		}
		var orders []nodes.Node
		for {
			expr := p.parseExpr()
			dir := nodes.SortDefault
			if _, ok := p.match(kwASC); ok {
				dir = nodes.SortAsc
			} else if _, ok := p.match(kwDESC); ok {
				dir = nodes.SortDesc
			}
			orders = append(orders, &nodes.OrderByItem{
				Expr:    expr,
				SortDir: dir,
				Loc:     nodes.Loc{Start: p.pos()},
			})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		over.OrderBy = &nodes.List{Items: orders}
	}

	// ROWS | RANGE | GROUPS window frame clause
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		over.WindowFrame = p.parseWindowFrame()
	}

	over.Loc.End = p.pos()
	_, _ = p.expect(')')
	return over
}

// parseWindowFrame parses a window frame specification (ROWS/RANGE/GROUPS).
//
//	{ ROWS | RANGE | GROUPS } <window frame extent>
//
//	<window frame extent> ::=
//	  {   <window frame preceding>
//	    | BETWEEN <window frame bound> AND <window frame bound>
//	  }
func (p *Parser) parseWindowFrame() *nodes.WindowFrame {
	frame := &nodes.WindowFrame{
		Loc: nodes.Loc{Start: p.pos()},
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

	if p.cur.Type == kwBETWEEN {
		// BETWEEN <bound> AND <bound>
		p.advance()
		frame.Start = p.parseWindowFrameBound()
		if _, err := p.expect(kwAND); err != nil {
			frame.Loc.End = p.pos()
			return frame
		}
		frame.End = p.parseWindowFrameBound()
	} else {
		// Short syntax: <window frame preceding> (implies AND CURRENT ROW)
		frame.Start = p.parseWindowFrameBound()
	}

	frame.Loc.End = p.pos()
	return frame
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
func (p *Parser) parseWindowFrameBound() *nodes.WindowBound {
	bound := &nodes.WindowBound{
		Loc: nodes.Loc{Start: p.pos()},
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
		bound.Offset = p.parseExpr()
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundFollowing
		}
	}

	bound.Loc.End = p.pos()
	return bound
}
