package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// Operator precedence levels for Pratt parsing, matching PostgreSQL's yacc grammar.
// Higher numbers bind tighter.
const (
	precNone       = 0
	precOr         = 1
	precAnd        = 2
	precNot        = 3  // prefix NOT
	precIs         = 4  // IS NULL, IS TRUE, ISNULL, NOTNULL
	precComparison = 5  // =, <, >, <=, >=, <>
	precIn         = 6  // BETWEEN, IN, LIKE, ILIKE, SIMILAR, NOT_LA
	precOp         = 7  // Op, OPERATOR()
	precAdd        = 8  // +, -
	precMul        = 9  // *, /, %
	precExp        = 10 // ^
	precAt         = 11 // AT TIME ZONE
	precCollate    = 12 // COLLATE
	precUnary      = 13 // unary +, -
	precTypecast   = 14 // ::
)

// ---------------------------------------------------------------------------
// parseAExpr parses a full expression using Pratt (precedence climbing) parsing.
//
// Ref: https://www.postgresql.org/docs/17/sql-expressions.html
//
//	a_expr: c_expr | a_expr op a_expr | ...
func (p *Parser) parseAExpr(minPrec int) nodes.Node {
	loc := p.pos()
	left := p.parseAExprAtom()
	if left == nil {
		return nil
	}
	for {
		if p.collectMode() {
			// Emit all infix operators with precedence >= minPrec.
			p.addInfixCandidates(minPrec)
			return left
		}
		prec := p.aExprInfixPrec()
		if prec < minPrec || prec == precNone {
			break
		}
		left = p.parseAExprInfix(left, prec)
		if left == nil {
			return nil
		}
		setNodeLoc(left, loc, p.pos())
	}
	return left
}

// aExprInfixPrec returns the infix/postfix precedence of the current token
// in an a_expr context. Returns precNone if the token is not an infix operator.
func (p *Parser) aExprInfixPrec() int {
	switch p.cur.Type {
	case OR:
		return precOr
	case AND:
		return precAnd
	case IS, ISNULL, NOTNULL:
		return precIs
	case '<', '>', '=':
		return precComparison
	case LESS_EQUALS, GREATER_EQUALS, NOT_EQUALS:
		return precComparison
	case BETWEEN, IN_P, LIKE, ILIKE, SIMILAR:
		return precIn
	case NOT_LA:
		return precIn
	case Op:
		return precOp
	case '+', '-':
		return precAdd
	case '*', '/', '%':
		return precMul
	case '^':
		return precExp
	case AT:
		return precAt
	case COLLATE:
		return precCollate
	case TYPECAST:
		return precTypecast
	case '[':
		return precTypecast // subscript binds same as typecast
	default:
		return precNone
	}
}

// parseAExprAtom parses the prefix/atom part of an a_expr.
// Handles prefix operators (NOT, unary +/-) and delegates to parseCExpr.
func (p *Parser) parseAExprAtom() nodes.Node {
	loc := p.pos()
	switch p.cur.Type {
	case NOT:
		p.advance()
		arg := p.parseAExpr(precNot)
		if arg == nil {
			return nil
		}
		n := &nodes.BoolExpr{
			Boolop: nodes.NOT_EXPR,
			Args:   &nodes.List{Items: []nodes.Node{arg}},
			Loc:    nodes.Loc{Start: loc, End: p.pos()},
		}
		return n
	case '+':
		// Unary plus: just return the operand (yacc does the same)
		p.advance()
		return p.parseAExpr(precUnary)
	case '-':
		// Unary minus
		p.advance()
		arg := p.parseAExpr(precUnary)
		if arg == nil {
			return nil
		}
		n := doNegate(arg)
		setNodeLoc(n, loc, p.pos())
		return n
	default:
		return p.parseCExpr()
	}
}

// parseAExprInfix handles an infix or postfix operator in a_expr.
func (p *Parser) parseAExprInfix(left nodes.Node, prec int) nodes.Node {
	switch p.cur.Type {
	case OR:
		p.advance()
		right := p.parseAExpr(prec + 1)
		return makeBoolExpr(nodes.OR_EXPR, left, right)
	case AND:
		p.advance()
		right := p.parseAExpr(prec + 1)
		return makeBoolExpr(nodes.AND_EXPR, left, right)

	// --- IS / ISNULL / NOTNULL ---
	case IS:
		return p.parseIsPostfix(left)
	case ISNULL:
		p.advance()
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NULL}
	case NOTNULL:
		p.advance()
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NOT_NULL}

	// --- Comparison operators ---
	case '<', '>', '=':
		tok := p.advance()
		opStr := string(rune(tok.Type))
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, opStr, left, right)
	case LESS_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: "<="}}}
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, "<=", left, right)
	case GREATER_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: ">="}}}
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, ">=", left, right)
	case NOT_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: "<>"}}}
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, "<>", left, right)

	// --- BETWEEN, IN, LIKE, ILIKE, SIMILAR ---
	case BETWEEN:
		return p.parseBetweenExpr(left, false)
	case IN_P:
		return p.parseInExpr(left, false)
	case LIKE:
		return p.parseLikeExpr(left, false)
	case ILIKE:
		return p.parseIlikeExpr(left, false)
	case SIMILAR:
		return p.parseSimilarExpr(left, false)
	case NOT_LA:
		// NOT_LA means NOT followed by BETWEEN/IN/LIKE/ILIKE/SIMILAR
		p.advance() // consume NOT
		switch p.cur.Type {
		case BETWEEN:
			return p.parseBetweenExpr(left, true)
		case IN_P:
			return p.parseInExpr(left, true)
		case LIKE:
			return p.parseLikeExpr(left, true)
		case ILIKE:
			return p.parseIlikeExpr(left, true)
		case SIMILAR:
			return p.parseSimilarExpr(left, true)
		default:
			return nil // shouldn't happen due to NOT_LA conversion logic
		}

	// --- Custom operators (Op, OPERATOR()) ---
	case Op:
		tok := p.advance()
		opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: tok.Str}}}
		// Check for subquery_Op sub_type pattern: expr Op ANY/ALL/SOME (subquery_or_expr)
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return &nodes.A_Expr{Kind: nodes.AEXPR_OP, Name: opName, Lexpr: left, Rexpr: right}

	// --- Arithmetic ---
	case '+', '-':
		tok := p.advance()
		opStr := string(rune(tok.Type))
		// Check for subquery_Op sub_type
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
			return p.parseSubqueryOp(left, opName)
		}
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, opStr, left, right)
	case '*', '/', '%':
		tok := p.advance()
		right := p.parseAExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right)
	case '^':
		p.advance()
		right := p.parseAExpr(prec + 1) // left-assoc in PostgreSQL
		return makeSimpleAExpr(nodes.AEXPR_OP, "^", left, right)

	// --- AT TIME ZONE ---
	case AT:
		return p.parseAtTimeZone(left)

	// --- COLLATE ---
	case COLLATE:
		p.advance()
		collname, err := p.parseAnyName()
		if err != nil {
			return nil
		}
		return &nodes.CollateClause{Arg: left, Collname: collname}

	// --- TYPECAST (::) ---
	case TYPECAST:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil
		}
		return &nodes.TypeCast{Arg: left, TypeName: tn, Loc: nodes.NoLoc()}

	// --- Array subscript ---
	case '[':
		return p.parseSubscript(left)
	}

	return left
}

// parseIsPostfix handles IS ... postfix expressions on left.
func (p *Parser) parseIsPostfix(left nodes.Node) nodes.Node {
	p.advance() // consume IS

	negated := false
	if p.cur.Type == NOT {
		negated = true
		p.advance()
	}

	switch p.cur.Type {
	case NULL_P:
		p.advance()
		if negated {
			return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NOT_NULL}
		}
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NULL}
	case TRUE_P:
		p.advance()
		bt := nodes.IS_TRUE
		if negated {
			bt = nodes.IS_NOT_TRUE
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}
	case FALSE_P:
		p.advance()
		bt := nodes.IS_FALSE
		if negated {
			bt = nodes.IS_NOT_FALSE
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}
	case UNKNOWN:
		p.advance()
		bt := nodes.IS_UNKNOWN
		if negated {
			bt = nodes.IS_NOT_UNKNOWN
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}
	case DISTINCT:
		// IS [NOT] DISTINCT FROM a_expr
		p.advance() // consume DISTINCT
		p.expect(FROM)
		right := p.parseAExpr(precIs + 1)
		kind := nodes.AEXPR_DISTINCT
		if negated {
			kind = nodes.AEXPR_NOT_DISTINCT
		}
		return makeSimpleAExpr(kind, "=", left, right)
	case DOCUMENT_P:
		p.advance()
		xmlExpr := &nodes.XmlExpr{
			Op:       nodes.IS_DOCUMENT,
			Args:     &nodes.List{Items: []nodes.Node{left}},
			Loc: nodes.NoLoc(),
		}
		if negated {
			return &nodes.BoolExpr{
				Boolop:   nodes.NOT_EXPR,
				Args:     &nodes.List{Items: []nodes.Node{xmlExpr}},
				Loc: nodes.NoLoc(),
			}
		}
		return xmlExpr
	case JSON:
		// IS [NOT] JSON [VALUE|ARRAY|OBJECT|SCALAR] [WITH UNIQUE [KEYS]]
		return p.parseJsonIsPredicate(left, negated)
	case NORMALIZED:
		// IS [NOT] NORMALIZED (default NFC form)
		p.advance()
		fc := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "is_normalized"),
			Args:       &nodes.List{Items: []nodes.Node{left, makeStringConst("NFC")}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
		if negated {
			return &nodes.BoolExpr{
				Boolop:   nodes.NOT_EXPR,
				Args:     &nodes.List{Items: []nodes.Node{fc}},
				Loc: nodes.NoLoc(),
			}
		}
		return fc
	case NFC, NFD, NFKC, NFKD:
		// IS [NOT] NFC/NFD/NFKC/NFKD NORMALIZED
		form := p.parseUnicodeNormalForm()
		p.expect(NORMALIZED)
		fc := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "is_normalized"),
			Args:       &nodes.List{Items: []nodes.Node{left, makeStringConst(form)}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
		if negated {
			return &nodes.BoolExpr{
				Boolop:   nodes.NOT_EXPR,
				Args:     &nodes.List{Items: []nodes.Node{fc}},
				Loc: nodes.NoLoc(),
			}
		}
		return fc
	default:
		return left
	}
}

// parseBetweenExpr parses BETWEEN [SYMMETRIC] b_expr AND a_expr.
func (p *Parser) parseBetweenExpr(left nodes.Node, negated bool) nodes.Node {
	p.advance() // consume BETWEEN

	symmetric := false
	if p.cur.Type == SYMMETRIC {
		p.advance()
		symmetric = true
	} else {
		// opt_asymmetric
		if p.cur.Type == ASYMMETRIC {
			p.advance()
		}
	}

	// Lower bound is b_expr (to avoid ambiguity with AND)
	lower := p.parseBExpr(0)
	p.expect(AND)
	upper := p.parseAExpr(precIn + 1)

	var kind nodes.A_Expr_Kind
	var op string
	if symmetric {
		if negated {
			kind = nodes.AEXPR_NOT_BETWEEN_SYM
			op = "NOT BETWEEN SYMMETRIC"
		} else {
			kind = nodes.AEXPR_BETWEEN_SYM
			op = "BETWEEN SYMMETRIC"
		}
	} else {
		if negated {
			kind = nodes.AEXPR_NOT_BETWEEN
			op = "NOT BETWEEN"
		} else {
			kind = nodes.AEXPR_BETWEEN
			op = "BETWEEN"
		}
	}

	return &nodes.A_Expr{
		Kind:     kind,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    &nodes.List{Items: []nodes.Node{lower, upper}},
	}
}

// parseInExpr parses IN (expr_list) or IN (subquery).
func (p *Parser) parseInExpr(left nodes.Node, negated bool) nodes.Node {
	p.advance() // consume IN

	op := "="
	if negated {
		op = "<>"
	}

	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // consume '('

	// Try to determine if this is a subquery or expression list.
	// Subqueries start with SELECT, VALUES, WITH, or '(' followed by another select.
	if p.isSelectStart() {
		subquery := p.parseSelectStmtForExpr()
		p.expect(')')
		// IN subquery: ANY_SUBLINK; NOT IN wraps in BoolExpr NOT (matches yacc)
		sub := &nodes.SubLink{
			SubLinkType: int(nodes.ANY_SUBLINK),
			Testexpr:    left,
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}
		if negated {
			return &nodes.BoolExpr{
				Boolop: nodes.NOT_EXPR,
				Args:   &nodes.List{Items: []nodes.Node{sub}},
			}
		}
		return sub
	}

	// Expression list
	exprs := p.parseExprListFull()
	p.expect(')')
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_IN,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    exprs,
	}
}

// parseLikeExpr parses [NOT] LIKE a_expr [ESCAPE a_expr].
func (p *Parser) parseLikeExpr(left nodes.Node, negated bool) nodes.Node {
	p.advance() // consume LIKE
	right := p.parseAExpr(precIn + 1)

	op := "~~"
	if negated {
		op = "!~~"
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr := p.parseAExpr(precIn + 1)
		escapedPattern := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "like_escape"),
			Args:       &nodes.List{Items: []nodes.Node{right, escapeExpr}},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}
		return &nodes.A_Expr{
			Kind:     nodes.AEXPR_LIKE,
			Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
			Lexpr:    left,
			Rexpr:    escapedPattern,
		}
	}

	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_LIKE,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    right,
	}
}

// parseIlikeExpr parses [NOT] ILIKE a_expr [ESCAPE a_expr].
func (p *Parser) parseIlikeExpr(left nodes.Node, negated bool) nodes.Node {
	p.advance() // consume ILIKE
	right := p.parseAExpr(precIn + 1)

	op := "~~*"
	if negated {
		op = "!~~*"
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr := p.parseAExpr(precIn + 1)
		escapedPattern := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "iclike_escape"),
			Args:       &nodes.List{Items: []nodes.Node{right, escapeExpr}},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}
		return &nodes.A_Expr{
			Kind:     nodes.AEXPR_ILIKE,
			Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
			Lexpr:    left,
			Rexpr:    escapedPattern,
		}
	}

	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_ILIKE,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    right,
	}
}

// parseSimilarExpr parses [NOT] SIMILAR TO a_expr [ESCAPE a_expr].
func (p *Parser) parseSimilarExpr(left nodes.Node, negated bool) nodes.Node {
	p.advance() // consume SIMILAR
	p.expect(TO)
	right := p.parseAExpr(precIn + 1)

	if negated {
		if p.cur.Type == ESCAPE {
			p.advance()
			escapeExpr := p.parseAExpr(precIn + 1)
			// NOT SIMILAR with ESCAPE
			escapedPattern := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
				Args:       &nodes.List{Items: []nodes.Node{right, escapeExpr}},
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Loc: nodes.NoLoc(),
			}
			return &nodes.A_Expr{
				Kind:     nodes.AEXPR_SIMILAR,
				Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "!~"}}},
				Lexpr:    left,
				Rexpr:    escapedPattern,
			}
		}
		// NOT SIMILAR without ESCAPE: wrap in similar_to_escape with no escape char
		wrapped := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
			Args:       &nodes.List{Items: []nodes.Node{right}},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}
		return &nodes.A_Expr{
			Kind:     nodes.AEXPR_SIMILAR,
			Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "!~"}}},
			Lexpr:    left,
			Rexpr:    wrapped,
		}
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr := p.parseAExpr(precIn + 1)
		escapedPattern := &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
			Args:       &nodes.List{Items: []nodes.Node{right, escapeExpr}},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}
		return &nodes.A_Expr{
			Kind:     nodes.AEXPR_SIMILAR,
			Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "~"}}},
			Lexpr:    left,
			Rexpr:    escapedPattern,
		}
	}

	wrapped := &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
		Args:       &nodes.List{Items: []nodes.Node{right}},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_SIMILAR,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "~"}}},
		Lexpr:    left,
		Rexpr:    wrapped,
	}
}

// parseSubqueryOp handles expr op ANY/ALL/SOME (subquery_or_expr).
func (p *Parser) parseSubqueryOp(left nodes.Node, opName *nodes.List) nodes.Node {
	subType := nodes.ANY_SUBLINK
	if p.cur.Type == ALL {
		subType = nodes.ALL_SUBLINK
	}
	p.advance() // consume ANY/ALL/SOME

	if p.cur.Type != '(' {
		return nil
	}
	p.advance()

	if p.isSelectStart() {
		subquery := p.parseSelectStmtForExpr()
		p.expect(')')
		return &nodes.SubLink{
			SubLinkType: int(subType),
			Testexpr:    left,
			OperName:    opName,
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}
	}

	// ANY/ALL (expr) -- array form
	expr := p.parseAExpr(0)
	p.expect(')')

	kind := nodes.AEXPR_OP_ANY
	if subType == nodes.ALL_SUBLINK {
		kind = nodes.AEXPR_OP_ALL
	}
	return &nodes.A_Expr{
		Kind:     kind,
		Name:     opName,
		Lexpr:    left,
		Rexpr:    expr,
	}
}

// parseAtTimeZone handles AT TIME ZONE a_expr and AT LOCAL.
func (p *Parser) parseAtTimeZone(left nodes.Node) nodes.Node {
	p.advance() // consume AT

	if p.cur.Type == LOCAL {
		// AT LOCAL => AT TIME ZONE 'DEFAULT' (handled as function call)
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "timezone"),
			Args:       &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.String{Str: "DEFAULT"}}, left}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	p.expect(TIME)
	p.expect(ZONE)
	tz := p.parseAExpr(precAt + 1)
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "timezone"),
		Args:       &nodes.List{Items: []nodes.Node{tz, left}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseSubscript handles array subscript and slice: expr[idx] or expr[lo:hi].
func (p *Parser) parseSubscript(left nodes.Node) nodes.Node {
	p.advance() // consume '['

	// Check for empty lower bound followed by ':'
	if p.cur.Type == ':' {
		p.advance()
		var uidx nodes.Node
		if p.cur.Type != ']' {
			uidx = p.parseAExpr(0)
		}
		p.expect(']')
		return &nodes.A_Indirection{
			Arg:         left,
			Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{IsSlice: true, Uidx: uidx}}},
		}
	}

	idx := p.parseAExpr(0)

	if p.cur.Type == ':' {
		// Slice
		p.advance()
		var uidx nodes.Node
		if p.cur.Type != ']' {
			uidx = p.parseAExpr(0)
		}
		p.expect(']')
		return &nodes.A_Indirection{
			Arg:         left,
			Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{IsSlice: true, Lidx: idx, Uidx: uidx}}},
		}
	}

	// Simple subscript
	p.expect(']')
	return &nodes.A_Indirection{
		Arg:         left,
		Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{Uidx: idx}}},
	}
}

// ---------------------------------------------------------------------------
// parseBExpr parses a restricted expression (no boolean, no IN/LIKE/BETWEEN, no IS NULL etc).
//
//	b_expr: c_expr | b_expr op b_expr | b_expr IS [NOT] DISTINCT FROM b_expr | ...
func (p *Parser) parseBExpr(minPrec int) nodes.Node {
	left := p.parseBExprAtom()
	if left == nil {
		return nil
	}
	for {
		prec := p.bExprInfixPrec()
		if prec < minPrec || prec == precNone {
			break
		}
		left = p.parseBExprInfix(left, prec)
		if left == nil {
			return nil
		}
	}
	return left
}

// parseBExprAtom handles prefix operators for b_expr.
func (p *Parser) parseBExprAtom() nodes.Node {
	switch p.cur.Type {
	case '+':
		p.advance()
		return p.parseBExpr(precUnary)
	case '-':
		p.advance()
		arg := p.parseBExpr(precUnary)
		if arg == nil {
			return nil
		}
		return doNegate(arg)
	default:
		return p.parseCExpr()
	}
}

// bExprInfixPrec returns precedence for b_expr infix operators.
func (p *Parser) bExprInfixPrec() int {
	switch p.cur.Type {
	case IS:
		return precIs
	case '<', '>', '=':
		return precComparison
	case LESS_EQUALS, GREATER_EQUALS, NOT_EQUALS:
		return precComparison
	case Op:
		return precOp
	case '+', '-':
		return precAdd
	case '*', '/', '%':
		return precMul
	case '^':
		return precExp
	case TYPECAST:
		return precTypecast
	default:
		return precNone
	}
}

// parseBExprInfix handles infix operators in b_expr.
func (p *Parser) parseBExprInfix(left nodes.Node, prec int) nodes.Node {
	switch p.cur.Type {
	case IS:
		p.advance()
		negated := false
		if p.cur.Type == NOT {
			negated = true
			p.advance()
		}
		if p.cur.Type == DISTINCT {
			p.advance()
			p.expect(FROM)
			right := p.parseBExpr(precIs + 1)
			kind := nodes.AEXPR_DISTINCT
			if negated {
				kind = nodes.AEXPR_NOT_DISTINCT
			}
			return makeSimpleAExpr(kind, "=", left, right)
		}
		if p.cur.Type == DOCUMENT_P {
			p.advance()
			xmlExpr := &nodes.XmlExpr{
				Op:       nodes.IS_DOCUMENT,
				Args:     &nodes.List{Items: []nodes.Node{left}},
				Loc: nodes.NoLoc(),
			}
			if negated {
				return &nodes.BoolExpr{
					Boolop: nodes.NOT_EXPR,
					Args:   &nodes.List{Items: []nodes.Node{xmlExpr}},
				}
			}
			return xmlExpr
		}
		return left
	case '<', '>', '=':
		tok := p.advance()
		right := p.parseBExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right)
	case LESS_EQUALS:
		p.advance()
		right := p.parseBExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, "<=", left, right)
	case GREATER_EQUALS:
		p.advance()
		right := p.parseBExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, ">=", left, right)
	case NOT_EQUALS:
		p.advance()
		right := p.parseBExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, "<>", left, right)
	case Op:
		tok := p.advance()
		opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: tok.Str}}}
		right := p.parseBExpr(prec + 1)
		return &nodes.A_Expr{Kind: nodes.AEXPR_OP, Name: opName, Lexpr: left, Rexpr: right}
	case '+', '-', '*', '/', '%', '^':
		tok := p.advance()
		right := p.parseBExpr(prec + 1)
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right)
	case TYPECAST:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil
		}
		return &nodes.TypeCast{Arg: left, TypeName: tn, Loc: nodes.NoLoc()}
	}
	return left
}

// isFuncExprCommonSubexprStart returns true if the current token can start
// a func_expr_common_subexpr production. This is used by parseTableRefPrimary
// to route tokens like USER, CURRENT_USER, etc. to the func_table path,
// matching PostgreSQL's table_ref grammar where func_table is a separate
// alternative from relation_expr.
func (p *Parser) isFuncExprCommonSubexprStart() bool {
	switch p.cur.Type {
	case CAST, COALESCE, GREATEST, LEAST, NULLIF,
		EXTRACT, NORMALIZE, OVERLAY, POSITION, SUBSTRING, TRIM,
		GROUPING, COLLATION, TREAT,
		// JSON functions
		JSON, JSON_OBJECT, JSON_ARRAY, JSON_SCALAR, JSON_SERIALIZE,
		JSON_QUERY, JSON_EXISTS, JSON_VALUE, JSON_OBJECTAGG, JSON_ARRAYAGG,
		// XML functions
		XMLCONCAT, XMLELEMENT, XMLEXISTS, XMLFOREST, XMLPARSE, XMLPI, XMLROOT, XMLSERIALIZE,
		// SQL value functions
		CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP, LOCALTIME, LOCALTIMESTAMP,
		CURRENT_ROLE, CURRENT_USER, SESSION_USER, USER, CURRENT_CATALOG, CURRENT_SCHEMA, SYSTEM_USER:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// parseCExpr parses primary expressions (atoms).
//
// Ref: gram.y c_expr rule
//
//	c_expr:
//	    columnref | AexprConst | PARAM opt_indirection
//	    | '(' a_expr ')' opt_indirection
//	    | func_expr | select_with_parens
//	    | EXISTS select_with_parens
//	    | case_expr | ARRAY select_with_parens | ARRAY array_expr
//	    | explicit_row | implicit_row
func (p *Parser) parseCExpr() nodes.Node {
	loc := p.pos()
	n := p.parseCExprInner()
	if n != nil {
		setNodeLoc(n, loc, p.pos())
	}
	return n
}

func (p *Parser) parseCExprInner() nodes.Node {
	if p.collectMode() {
		p.cachedCollect("parseCExprInner", func() {
			// Expression-starting tokens.
			exprTokens := []int{
				ICONST, FCONST, SCONST, BCONST, XCONST,
				TRUE_P, FALSE_P, NULL_P, PARAM,
				'(', EXISTS, ARRAY, CASE, ROW, DEFAULT,
				CAST, COALESCE, GREATEST, LEAST, NULLIF,
				EXTRACT, NORMALIZE, OVERLAY, POSITION, SUBSTRING, TRIM,
				GROUPING, COLLATION, TREAT,
				JSON, JSON_OBJECT, JSON_ARRAY, JSON_SCALAR, JSON_SERIALIZE,
				JSON_QUERY, JSON_EXISTS, JSON_VALUE, JSON_OBJECTAGG, JSON_ARRAYAGG,
				XMLCONCAT, XMLELEMENT, XMLEXISTS, XMLFOREST, XMLPARSE, XMLPI, XMLROOT, XMLSERIALIZE,
				CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP, LOCALTIME, LOCALTIMESTAMP,
				CURRENT_ROLE, CURRENT_USER, SESSION_USER, USER, CURRENT_CATALOG, CURRENT_SCHEMA, SYSTEM_USER,
				NOT, '+', '-',
			}
			for _, t := range exprTokens {
				p.addTokenCandidate(t)
			}
			// Identifiers/functions are also valid expression starts.
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			// Unreserved and col_name keywords are valid as identifiers in expressions.
			p.addTokenCandidate(IDENT)
			p.addKeywordsByCategory(UnreservedKeyword, ColNameKeyword)
		})
		return nil
	}
	switch p.cur.Type {
	// --- Constants ---
	case ICONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.Integer{Ival: tok.Ival}}
	case FCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.Float{Fval: tok.Str}}
	case SCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.String{Str: tok.Str}}
	case BCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.BitString{Bsval: tok.Str}}
	case XCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.BitString{Bsval: tok.Str}}
	case TRUE_P:
		p.advance()
		return &nodes.A_Const{Val: &nodes.Boolean{Boolval: true}}
	case FALSE_P:
		p.advance()
		return &nodes.A_Const{Val: &nodes.Boolean{Boolval: false}}
	case NULL_P:
		p.advance()
		return &nodes.A_Const{Isnull: true}

	// --- Parameter ($1, $2, ...) ---
	case PARAM:
		tok := p.advance()
		ref := &nodes.ParamRef{Number: int(tok.Ival)}
		indir := p.parseOptIndirection()
		if indir != nil && len(indir.Items) > 0 {
			return &nodes.A_Indirection{Arg: ref, Indirection: indir}
		}
		return ref

	// --- Parenthesized expr, implicit row, or subquery ---
	case '(':
		return p.parseParenExprOrRow()

	// --- EXISTS ---
	case EXISTS:
		return p.parseExistsExpr()

	// --- ARRAY ---
	case ARRAY:
		return p.parseArrayCExpr()

	// --- CASE ---
	case CASE:
		return p.parseCaseExpr()

	// --- ROW ---
	case ROW:
		return p.parseExplicitRow()

	// --- DEFAULT ---
	case DEFAULT:
		p.advance()
		return &nodes.SetToDefault{Loc: nodes.NoLoc()}

	// --- func_expr_common_subexpr keywords ---
	case CAST:
		return p.parseCastExpr()
	case COALESCE:
		return p.parseCoalesceExpr()
	case GREATEST:
		return p.parseMinMaxExpr(nodes.IS_GREATEST)
	case LEAST:
		return p.parseMinMaxExpr(nodes.IS_LEAST)
	case NULLIF:
		return p.parseNullIfExpr()
	case EXTRACT:
		return p.parseExtractExpr()
	case NORMALIZE:
		return p.parseNormalizeExpr()
	case OVERLAY:
		return p.parseOverlayExpr()
	case POSITION:
		return p.parsePositionExpr()
	case SUBSTRING:
		return p.parseSubstringExpr()
	case TRIM:
		return p.parseTrimExpr()
	case GROUPING:
		return p.parseGroupingExpr()
	case COLLATION:
		return p.parseCollationForExpr()
	case TREAT:
		return p.parseTreatExpr()

	// --- JSON functions ---
	case JSON:
		return p.parseJsonParseExpr()
	case JSON_OBJECT:
		return p.parseJsonObjectExpr()
	case JSON_ARRAY:
		return p.parseJsonArrayExpr()
	case JSON_SCALAR:
		return p.parseJsonScalarExpr()
	case JSON_SERIALIZE:
		return p.parseJsonSerializeExpr()
	case JSON_QUERY:
		return p.parseJsonQueryExpr()
	case JSON_EXISTS:
		return p.parseJsonExistsExpr()
	case JSON_VALUE:
		return p.parseJsonValueFuncExpr()
	case JSON_OBJECTAGG:
		return p.parseJsonObjectAgg()
	case JSON_ARRAYAGG:
		return p.parseJsonArrayAgg()

	// --- XML functions ---
	case XMLCONCAT:
		return p.parseXmlConcat()
	case XMLELEMENT:
		return p.parseXmlElement()
	case XMLEXISTS:
		return p.parseXmlExists()
	case XMLFOREST:
		return p.parseXmlForest()
	case XMLPARSE:
		return p.parseXmlParse()
	case XMLPI:
		return p.parseXmlPI()
	case XMLROOT:
		return p.parseXmlRoot()
	case XMLSERIALIZE:
		return p.parseXmlSerialize()

	// --- SQL value functions (no parens) ---
	case CURRENT_DATE:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_DATE, Typmod: -1, Loc: nodes.NoLoc()}
	case CURRENT_TIME:
		return p.parseSVFWithOptionalPrecision(nodes.SVFOP_CURRENT_TIME, nodes.SVFOP_CURRENT_TIME_N)
	case CURRENT_TIMESTAMP:
		return p.parseSVFWithOptionalPrecision(nodes.SVFOP_CURRENT_TIMESTAMP, nodes.SVFOP_CURRENT_TIMESTAMP_N)
	case LOCALTIME:
		return p.parseSVFWithOptionalPrecision(nodes.SVFOP_LOCALTIME, nodes.SVFOP_LOCALTIME_N)
	case LOCALTIMESTAMP:
		return p.parseSVFWithOptionalPrecision(nodes.SVFOP_LOCALTIMESTAMP, nodes.SVFOP_LOCALTIMESTAMP_N)
	case CURRENT_ROLE:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_ROLE, Typmod: -1, Loc: nodes.NoLoc()}
	case CURRENT_USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_USER, Typmod: -1, Loc: nodes.NoLoc()}
	case SESSION_USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_SESSION_USER, Typmod: -1, Loc: nodes.NoLoc()}
	case USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_USER, Typmod: -1, Loc: nodes.NoLoc()}
	case CURRENT_CATALOG:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_CATALOG, Typmod: -1, Loc: nodes.NoLoc()}
	case CURRENT_SCHEMA:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_SCHEMA, Typmod: -1, Loc: nodes.NoLoc()}
	case SYSTEM_USER:
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "system_user"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}

	default:
		// --- AexprConst: type-casted string constants (e.g., int '42', interval '1 day') ---
		if p.isAExprConstTypeCast() {
			return p.parseTypeCastedConst()
		}
		// --- func_expr: func_application or columnref ---
		if p.isColId() || p.isTypeFunctionName() {
			return p.parseColumnRefOrFuncCall()
		}
	}

	return nil
}

// parseParenExprOrRow handles '(' ... ')' which could be:
// - Parenthesized expression: (a_expr) opt_indirection
// - Implicit row: (a_expr, a_expr, ...)
// - Subquery: (SELECT ...)
func (p *Parser) parseParenExprOrRow() nodes.Node {
	p.advance() // consume '('

	// Check for subquery
	if p.isSelectStart() {
		subquery := p.parseSelectStmtForExpr()
		p.expect(')')
		// Check for indirection after subquery
		if p.cur.Type == '.' || p.cur.Type == '[' {
			indir, err := p.parseIndirection()
			if err == nil && indir != nil && len(indir.Items) > 0 {
				sub := &nodes.SubLink{
					SubLinkType: int(nodes.EXPR_SUBLINK),
					Subselect:   subquery,
					Loc: nodes.NoLoc(),
				}
				return &nodes.A_Indirection{Arg: sub, Indirection: indir}
			}
		}
		return &nodes.SubLink{
			SubLinkType: int(nodes.EXPR_SUBLINK),
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}
	}

	// Parse first expression
	first := p.parseAExpr(0)

	if p.cur.Type == ',' {
		// Implicit row: (expr, expr, ...)
		items := []nodes.Node{first}
		for p.cur.Type == ',' {
			p.advance()
			expr := p.parseAExpr(0)
			items = append(items, expr)
		}
		p.expect(')')
		if len(items) < 2 {
			// Single element + trailing comma should not happen but handle gracefully
			return first
		}
		return &nodes.RowExpr{
			Args:      &nodes.List{Items: items},
			RowFormat: nodes.COERCE_IMPLICIT_CAST,
			Loc: nodes.NoLoc(),
		}
	}

	// Parenthesized expression
	p.expect(')')
	// Check for indirection
	indir := p.parseOptIndirection()
	if indir != nil && len(indir.Items) > 0 {
		return &nodes.A_Indirection{Arg: first, Indirection: indir}
	}
	return first
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() nodes.Node {
	p.advance() // consume EXISTS
	p.expect('(')
	subquery := p.parseSelectStmtForExpr()
	p.expect(')')
	return &nodes.SubLink{
		SubLinkType: int(nodes.EXISTS_SUBLINK),
		Subselect:   subquery,
		Loc: nodes.NoLoc(),
	}
}

// parseArrayCExpr parses ARRAY [...] or ARRAY (subquery).
func (p *Parser) parseArrayCExpr() nodes.Node {
	p.advance() // consume ARRAY

	if p.cur.Type == '[' {
		arr := p.parseArrayExpr()
		if arr != nil {
			if a, ok := arr.(*nodes.A_ArrayExpr); ok {
				a.Loc = nodes.NoLoc()
			}
		}
		return arr
	}

	if p.cur.Type == '(' {
		p.advance()
		subquery := p.parseSelectStmtForExpr()
		p.expect(')')
		return &nodes.SubLink{
			SubLinkType: int(nodes.ARRAY_SUBLINK),
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}
	}

	return nil
}

// parseArrayExpr parses array literal: '[' expr_list ']' or '[' array_expr_list ']' or '[]'.
//
//	array_expr:
//	    '[' expr_list ']'
//	    | '[' array_expr_list ']'
//	    | '[' ']'
func (p *Parser) parseArrayExpr() nodes.Node {
	p.advance() // consume '['

	if p.cur.Type == ']' {
		p.advance()
		return &nodes.A_ArrayExpr{Loc: nodes.NoLoc()}
	}

	// Check if first element is a nested array
	if p.cur.Type == '[' {
		// array_expr_list
		first := p.parseArrayExpr()
		items := []nodes.Node{first}
		for p.cur.Type == ',' {
			p.advance()
			items = append(items, p.parseArrayExpr())
		}
		p.expect(']')
		return &nodes.A_ArrayExpr{Elements: &nodes.List{Items: items}}
	}

	// expr_list
	exprs := p.parseExprListFull()
	p.expect(']')
	return &nodes.A_ArrayExpr{Elements: exprs}
}

// parseCaseExpr parses CASE ... END.
//
// Ref: https://www.postgresql.org/docs/17/sql-expressions.html#SYNTAX-EXPRESSION-EVAL
//
//	CASE [expr] WHEN expr THEN expr [WHEN ...] [ELSE expr] END
func (p *Parser) parseCaseExpr() nodes.Node {
	p.advance() // consume CASE

	// Optional case argument (simple CASE)
	var caseArg nodes.Node
	if p.cur.Type != WHEN {
		caseArg = p.parseAExpr(0)
	}

	// Parse WHEN clauses
	var whens []nodes.Node
	for p.cur.Type == WHEN {
		p.advance() // consume WHEN
		expr := p.parseAExpr(0)
		p.expect(THEN)
		result := p.parseAExpr(0)
		whens = append(whens, &nodes.CaseWhen{
			Expr:     expr,
			Result:   result,
			Loc: nodes.NoLoc(),
		})
	}

	// Optional ELSE
	var defResult nodes.Node
	if p.cur.Type == ELSE {
		p.advance()
		defResult = p.parseAExpr(0)
	}

	p.expect(END_P)

	return &nodes.CaseExpr{
		Arg:       caseArg,
		Args:      &nodes.List{Items: whens},
		Defresult: defResult,
		Loc: nodes.NoLoc(),
	}
}

// parseExplicitRow parses ROW(...).
//
//	explicit_row: ROW '(' expr_list ')' | ROW '(' ')'
func (p *Parser) parseExplicitRow() nodes.Node {
	p.advance() // consume ROW
	p.expect('(')

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.RowExpr{
			RowFormat: nodes.COERCE_EXPLICIT_CALL,
			Loc: nodes.NoLoc(),
		}
	}

	exprs := p.parseExprListFull()
	p.expect(')')
	return &nodes.RowExpr{
		Args:      exprs,
		RowFormat: nodes.COERCE_EXPLICIT_CALL,
		Loc: nodes.NoLoc(),
	}
}

// ---------------------------------------------------------------------------
// func_expr_common_subexpr implementations

// parseCastExpr parses CAST(a_expr AS Typename).
func (p *Parser) parseCastExpr() nodes.Node {
	p.advance() // consume CAST
	p.expect('(')
	arg := p.parseAExpr(0)
	p.expect(AS)
	tn, err := p.parseTypename()
	if err != nil {
		return nil
	}
	p.expect(')')
	return &nodes.TypeCast{Arg: arg, TypeName: tn, Loc: nodes.NoLoc()}
}

// parseCoalesceExpr parses COALESCE(expr_list).
func (p *Parser) parseCoalesceExpr() nodes.Node {
	p.advance() // consume COALESCE
	p.expect('(')
	exprs := p.parseExprListFull()
	p.expect(')')
	return &nodes.CoalesceExpr{Args: exprs}
}

// parseMinMaxExpr parses GREATEST(expr_list) or LEAST(expr_list).
func (p *Parser) parseMinMaxExpr(op nodes.MinMaxOp) nodes.Node {
	p.advance() // consume GREATEST/LEAST
	p.expect('(')
	exprs := p.parseExprListFull()
	p.expect(')')
	return &nodes.MinMaxExpr{Op: op, Args: exprs}
}

// parseNullIfExpr parses NULLIF(a_expr, a_expr).
func (p *Parser) parseNullIfExpr() nodes.Node {
	p.advance() // consume NULLIF
	p.expect('(')
	left := p.parseAExpr(0)
	p.expect(',')
	right := p.parseAExpr(0)
	p.expect(')')
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_NULLIF,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
		Lexpr:    left,
		Rexpr:    right,
	}
}

// parseExtractExpr parses EXTRACT(field FROM a_expr).
func (p *Parser) parseExtractExpr() nodes.Node {
	p.advance() // consume EXTRACT
	p.expect('(')

	if p.cur.Type == ')' {
		// EXTRACT() with no args
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "extract"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	field := p.parseExtractArg()
	p.expect(FROM)
	expr := p.parseAExpr(0)
	p.expect(')')

	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "extract"),
		Args:       &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.String{Str: field}}, expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseExtractArg parses extract field names.
func (p *Parser) parseExtractArg() string {
	switch p.cur.Type {
	case YEAR_P:
		p.advance()
		return "year"
	case MONTH_P:
		p.advance()
		return "month"
	case DAY_P:
		p.advance()
		return "day"
	case HOUR_P:
		p.advance()
		return "hour"
	case MINUTE_P:
		p.advance()
		return "minute"
	case SECOND_P:
		p.advance()
		return "second"
	case SCONST:
		tok := p.advance()
		return tok.Str
	default:
		// IDENT
		tok := p.advance()
		return tok.Str
	}
}

// parseNormalizeExpr parses NORMALIZE(a_expr [, form]).
func (p *Parser) parseNormalizeExpr() nodes.Node {
	p.advance() // consume NORMALIZE
	p.expect('(')
	expr := p.parseAExpr(0)

	if p.cur.Type == ',' {
		p.advance()
		form := p.parseUnicodeNormalForm()
		p.expect(')')
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "normalize"),
			Args:       &nodes.List{Items: []nodes.Node{expr, &nodes.A_Const{Val: &nodes.String{Str: form}}}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "normalize"),
		Args:       &nodes.List{Items: []nodes.Node{expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseUnicodeNormalForm parses NFC | NFD | NFKC | NFKD.
func (p *Parser) parseUnicodeNormalForm() string {
	switch p.cur.Type {
	case NFC:
		p.advance()
		return "NFC"
	case NFD:
		p.advance()
		return "NFD"
	case NFKC:
		p.advance()
		return "NFKC"
	case NFKD:
		p.advance()
		return "NFKD"
	default:
		tok := p.advance()
		return tok.Str
	}
}

// parseOverlayExpr parses OVERLAY(a_expr PLACING a_expr FROM a_expr [FOR a_expr]).
func (p *Parser) parseOverlayExpr() nodes.Node {
	p.advance() // consume OVERLAY
	p.expect('(')

	// Try SQL-standard syntax: a_expr PLACING a_expr FROM a_expr [FOR a_expr]
	first := p.parseAExpr(0)
	if p.cur.Type == PLACING {
		p.advance()
		second := p.parseAExpr(0)
		p.expect(FROM)
		third := p.parseAExpr(0)
		args := []nodes.Node{first, second, third}
		if p.cur.Type == FOR {
			p.advance()
			fourth := p.parseAExpr(0)
			args = append(args, fourth)
		}
		p.expect(')')
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "overlay"),
			Args:       &nodes.List{Items: args},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	// Function-call syntax: overlay(arg, arg, ...)
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		args = append(args, p.parseFuncArgExpr())
	}
	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("overlay"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}
}

// parsePositionExpr parses POSITION(b_expr IN b_expr).
func (p *Parser) parsePositionExpr() nodes.Node {
	p.advance() // consume POSITION
	p.expect('(')

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "position"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	first := p.parseBExpr(0)
	if p.cur.Type == IN_P {
		p.advance()
		second := p.parseBExpr(0)
		p.expect(')')
		// Note: PG reverses the arguments
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "position"),
			Args:       &nodes.List{Items: []nodes.Node{second, first}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	// Function-call syntax
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		args = append(args, p.parseFuncArgExpr())
	}
	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("position"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}
}

// parseSubstringExpr parses SUBSTRING(... FROM ... [FOR ...]) or SUBSTRING(expr_list).
func (p *Parser) parseSubstringExpr() nodes.Node {
	p.advance() // consume SUBSTRING
	p.expect('(')

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	first := p.parseAExpr(0)

	// Check for SQL-standard syntax with FROM/FOR/SIMILAR
	if p.cur.Type == FROM {
		p.advance()
		second := p.parseAExpr(0)
		if p.cur.Type == FOR {
			p.advance()
			third := p.parseAExpr(0)
			p.expect(')')
			return &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "substring"),
				Args:       &nodes.List{Items: []nodes.Node{first, second, third}},
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Loc: nodes.NoLoc(),
			}
		}
		p.expect(')')
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, second}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	if p.cur.Type == FOR {
		p.advance()
		second := p.parseAExpr(0)
		if p.cur.Type == FROM {
			p.advance()
			third := p.parseAExpr(0)
			p.expect(')')
			return &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "substring"),
				Args:       &nodes.List{Items: []nodes.Node{first, third, second}},
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Loc: nodes.NoLoc(),
			}
		}
		p.expect(')')
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, makeIntConst(1), second}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	if p.cur.Type == SIMILAR {
		p.advance()
		second := p.parseAExpr(0)
		p.expect(ESCAPE)
		third := p.parseAExpr(0)
		p.expect(')')
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, second, third}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}
	}

	// Comma-separated function call form
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		args = append(args, p.parseFuncArgExpr())
	}
	p.expect(')')

	if len(args) > 1 {
		// Regular function call form: substring(str, start, len)
		return &nodes.FuncCall{
			Funcname:   makeFuncName("substring"),
			Args:       &nodes.List{Items: args},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "substring"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseTrimExpr parses TRIM([BOTH|LEADING|TRAILING] [expr FROM] expr_list).
func (p *Parser) parseTrimExpr() nodes.Node {
	p.advance() // consume TRIM
	p.expect('(')

	funcName := "btrim" // default is BOTH
	switch p.cur.Type {
	case BOTH:
		p.advance()
		funcName = "btrim"
	case LEADING:
		p.advance()
		funcName = "ltrim"
	case TRAILING:
		p.advance()
		funcName = "rtrim"
	}

	// Parse trim_list
	args := p.parseTrimList()

	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", funcName),
		Args:       args,
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseTrimList parses trim_list.
//
//	trim_list:
//	    a_expr FROM expr_list
//	    | FROM expr_list
//	    | expr_list
func (p *Parser) parseTrimList() *nodes.List {
	if p.cur.Type == FROM {
		p.advance()
		return p.parseExprListFull()
	}

	first := p.parseAExpr(0)
	if p.cur.Type == FROM {
		p.advance()
		rest := p.parseExprListFull()
		// Prepend the trim character
		items := make([]nodes.Node, 0, len(rest.Items)+1)
		items = append(items, first)
		items = append(items, rest.Items...)
		return &nodes.List{Items: items}
	}

	// Plain expression list
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseAExpr(0))
	}
	return &nodes.List{Items: items}
}

// parseGroupingExpr parses GROUPING(expr_list).
func (p *Parser) parseGroupingExpr() nodes.Node {
	p.advance() // consume GROUPING
	p.expect('(')
	exprs := p.parseExprListFull()
	p.expect(')')
	return &nodes.GroupingFunc{Args: exprs}
}

// parseCollationForExpr parses COLLATION FOR (a_expr).
func (p *Parser) parseCollationForExpr() nodes.Node {
	p.advance() // consume COLLATION
	p.expect(FOR)
	p.expect('(')
	expr := p.parseAExpr(0)
	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "pg_collation_for"),
		Args:       &nodes.List{Items: []nodes.Node{expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}
}

// parseTreatExpr parses TREAT(a_expr AS Typename).
func (p *Parser) parseTreatExpr() nodes.Node {
	p.advance() // consume TREAT
	p.expect('(')
	expr := p.parseAExpr(0)
	p.expect(AS)
	tn, err := p.parseTypename()
	if err != nil {
		return nil
	}
	p.expect(')')
	funcNameStr := "treat"
	if tn.Names != nil && len(tn.Names.Items) > 0 {
		if nameNode, ok := tn.Names.Items[len(tn.Names.Items)-1].(*nodes.String); ok {
			funcNameStr = nameNode.Str
		}
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", funcNameStr),
		Args:       &nodes.List{Items: []nodes.Node{expr}},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}
}

// parseSVFWithOptionalPrecision parses SQL value functions that take optional (int) precision.
func (p *Parser) parseSVFWithOptionalPrecision(noArgOp, withArgOp nodes.SVFOp) nodes.Node {
	p.advance() // consume keyword
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ICONST {
			prec := p.cur.Ival
			p.advance()
			p.expect(')')
			return &nodes.SQLValueFunction{Op: withArgOp, Typmod: int32(prec), Loc: nodes.NoLoc()}
		}
		p.expect(')')
	}
	return &nodes.SQLValueFunction{Op: noArgOp, Typmod: -1, Loc: nodes.NoLoc()}
}

// ---------------------------------------------------------------------------
// Function call parsing (minimal for batch 3, expanded in batch 4)

// parseColumnRefOrFuncCall disambiguates between column reference and function call.
// If followed by '(', it's a function call. Otherwise, column reference.
func (p *Parser) parseColumnRefOrFuncCall() nodes.Node {
	// Check for AexprConst: func_name Sconst (type-casted literal)
	// This is handled by isAExprConstTypeCast / parseTypeCastedConst above.

	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil
	}

	// Check for func_name Sconst (type-casted constant like: jsonb '[]', text 'hello')
	if p.cur.Type == SCONST {
		tok := p.advance()
		tn := &nodes.TypeName{
			Names:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
			TypeOid: 0,
			Loc:     nodes.NoLoc(),
		}
		return &nodes.TypeCast{
			Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
			TypeName: tn,
			Loc:      nodes.NoLoc(),
		}
	}

	// Check for function call
	if p.cur.Type == '(' {
		funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
		return p.parseFuncApplication(funcName, loc)
	}

	// Check for qualified name: schema.func(...) or schema.table.column
	if p.cur.Type == '.' {
		p.advance()
		if p.cur.Type == '*' {
			// schema.*
			p.advance()
			return &nodes.ColumnRef{
				Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.A_Star{}}},
			}
		}

		attr, err := p.parseAttrName()
		if err != nil {
			return nil
		}

		// schema.typename Sconst (e.g., pg_catalog.int4 '42')
		if p.cur.Type == SCONST {
			tok := p.advance()
			tn := &nodes.TypeName{
				Names:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}},
				TypeOid: 0,
				Loc:     nodes.NoLoc(),
			}
			return &nodes.TypeCast{
				Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
				TypeName: tn,
				Loc:      nodes.NoLoc(),
			}
		}

		// schema.func(...)
		if p.cur.Type == '(' {
			funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}}
			return p.parseFuncApplication(funcName, loc)
		}

		// schema.table.column or schema.column
		if p.cur.Type == '.' {
			p.advance()
			if p.cur.Type == '*' {
				p.advance()
				return &nodes.ColumnRef{
					Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.A_Star{}}},
				}
			}
			attr2, err := p.parseAttrName()
			if err != nil {
				return nil
			}
			return &nodes.ColumnRef{
				Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2}}},
			}
		}

		return &nodes.ColumnRef{
			Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}},
		}
	}

	// Simple column reference
	return &nodes.ColumnRef{
		Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
	}
}

// parseFuncApplication parses a function call after the name has been read.
//
// Ref: gram.y func_application rule
//
//	func_name '(' ')' | func_name '(' func_arg_list opt_sort_clause ')'
//	| func_name '(' VARIADIC func_arg_expr opt_sort_clause ')'
//	| func_name '(' func_arg_list ',' VARIADIC func_arg_expr opt_sort_clause ')'
//	| func_name '(' '*' ')' | func_name '(' DISTINCT func_arg_list opt_sort_clause ')'
//	| func_name '(' ALL func_arg_list opt_sort_clause ')'
func (p *Parser) parseFuncApplication(funcName *nodes.List, loc int) nodes.Node {
	p.advance() // consume '('

	fc := &nodes.FuncCall{
		Funcname: funcName,
	}

	if p.cur.Type == ')' {
		p.advance()
	} else if p.cur.Type == '*' {
		p.advance()
		fc.AggStar = true
		p.expect(')')
	} else if p.cur.Type == DISTINCT {
		p.advance()
		fc.AggDistinct = true
		fc.Args = p.parseFuncArgListFull()
		fc.AggOrder = p.parseOptSortClause()
		p.expect(')')
	} else if p.cur.Type == ALL {
		p.advance()
		fc.Args = p.parseFuncArgListFull()
		fc.AggOrder = p.parseOptSortClause()
		p.expect(')')
	} else if p.cur.Type == VARIADIC {
		p.advance()
		fc.FuncVariadic = true
		fc.Args = &nodes.List{Items: []nodes.Node{p.parseFuncArgExpr()}}
		fc.AggOrder = p.parseOptSortClause()
		p.expect(')')
	} else {
		// Regular argument list
		args := p.parseFuncArgListFull()

		if p.cur.Type == ',' && args != nil {
			// Check for VARIADIC after regular args
			p.advance()
			if p.cur.Type == VARIADIC {
				p.advance()
				fc.FuncVariadic = true
				variadicArg := p.parseFuncArgExpr()
				args.Items = append(args.Items, variadicArg)
			}
		}

		fc.Args = args
		fc.AggOrder = p.parseOptSortClause()
		p.expect(')')
	}

	// within_group_clause: WITHIN GROUP '(' sort_clause ')'
	if p.cur.Type == WITHIN {
		p.advance()
		p.expect(GROUP_P)
		p.expect('(')
		if p.cur.Type == ORDER {
			p.advance()
			p.expect(BY)
			fc.AggOrder = p.parseSortByList()
		}
		p.expect(')')
		fc.AggWithinGroup = true
	}

	// filter_clause: FILTER '(' WHERE a_expr ')'
	if p.cur.Type == FILTER {
		p.advance()
		p.expect('(')
		p.expect(WHERE)
		fc.AggFilter = p.parseAExpr(0)
		p.expect(')')
	}

	// over_clause: OVER window_specification | OVER ColId
	if p.cur.Type == OVER {
		fc.Over = p.parseOverClause()
	}

	return fc
}

// parseOptSortClause parses optional ORDER BY sort clause (used in aggregates).
// Returns nil if no ORDER BY present.
func (p *Parser) parseOptSortClause() *nodes.List {
	if p.cur.Type != ORDER {
		return nil
	}
	p.advance() // consume ORDER
	p.expect(BY)
	return p.parseSortByList()
}

// parseSortByList parses a comma-separated list of sortby items.
func (p *Parser) parseSortByList() *nodes.List {
	first := p.parseSortBy()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseSortBy())
	}
	return &nodes.List{Items: items}
}

// parseSortBy parses a single sortby item: a_expr [ASC|DESC] [NULLS FIRST|LAST].
func (p *Parser) parseSortBy() nodes.Node {
	loc := p.pos()
	expr := p.parseAExpr(0)
	dir := nodes.SORTBY_DEFAULT
	if p.cur.Type == ASC {
		p.advance()
		dir = nodes.SORTBY_ASC
	} else if p.cur.Type == DESC {
		p.advance()
		dir = nodes.SORTBY_DESC
	} else if p.cur.Type == USING {
		p.advance()
		dir = nodes.SORTBY_USING
		// consume operator
		p.parseAllOp()
	}

	nullsOrder := nodes.SORTBY_NULLS_DEFAULT
	if p.cur.Type == NULLS_LA {
		p.advance()
		if p.cur.Type == FIRST_P {
			p.advance()
			nullsOrder = nodes.SORTBY_NULLS_FIRST
		} else if p.cur.Type == LAST_P {
			p.advance()
			nullsOrder = nodes.SORTBY_NULLS_LAST
		}
	}

	return &nodes.SortBy{
		Node:        expr,
		SortbyDir:   dir,
		SortbyNulls: nullsOrder,
		Loc:         nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseOverClause parses OVER window_specification or OVER existing_window_name.
//
// Ref: gram.y over_clause
//
//	OVER window_specification
//	| OVER ColId
func (p *Parser) parseOverClause() nodes.Node {
	p.advance() // consume OVER
	if p.cur.Type == '(' {
		return p.parseWindowSpecification()
	}
	// OVER ColId - references an existing window by name
	loc := p.pos()
	name, _ := p.parseColId()
	return &nodes.WindowDef{
		Name:         name,
		FrameOptions: nodes.FRAMEOPTION_DEFAULTS,
		Loc:          nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseWindowSpecification parses a window specification: ( [name] [PARTITION BY ...] [ORDER BY ...] [frame] ).
func (p *Parser) parseWindowSpecification() nodes.Node {
	loc := p.pos()
	p.advance() // consume '('
	wd := &nodes.WindowDef{Loc: nodes.Loc{Start: loc, End: -1}}

	// Optional existing window name
	if p.isColId() && p.cur.Type != PARTITION && p.cur.Type != ORDER && p.cur.Type != RANGE && p.cur.Type != ROWS && p.cur.Type != GROUPS {
		name, _ := p.parseColId()
		wd.Refname = name
	}

	// PARTITION BY
	if p.cur.Type == PARTITION {
		p.advance()
		p.expect(BY)
		wd.PartitionClause = p.parseExprListFull()
	}

	// ORDER BY
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		wd.OrderClause = p.parseSortByList()
	}

	// Frame clause
	if p.cur.Type == RANGE || p.cur.Type == ROWS || p.cur.Type == GROUPS {
		p.parseFrameClause(wd)
	} else {
		wd.FrameOptions = nodes.FRAMEOPTION_DEFAULTS
	}

	p.expect(')')
	wd.Loc.End = p.pos()
	return wd
}

// parseFrameClause parses window frame specification.
func (p *Parser) parseFrameClause(wd *nodes.WindowDef) {
	frameOptions := nodes.FRAMEOPTION_NONDEFAULT

	switch p.cur.Type {
	case RANGE:
		p.advance()
		frameOptions |= nodes.FRAMEOPTION_RANGE
	case ROWS:
		p.advance()
		frameOptions |= nodes.FRAMEOPTION_ROWS
	case GROUPS:
		p.advance()
		frameOptions |= nodes.FRAMEOPTION_GROUPS
	}

	if p.cur.Type == BETWEEN {
		p.advance()
		frameOptions |= nodes.FRAMEOPTION_BETWEEN
		frameOptions, wd.StartOffset = p.parseFrameBound(frameOptions, true)
		p.expect(AND)
		frameOptions, wd.EndOffset = p.parseFrameBound(frameOptions, false)
	} else {
		frameOptions, wd.StartOffset = p.parseFrameBound(frameOptions, true)
		// Default end bound: CURRENT ROW
		frameOptions |= nodes.FRAMEOPTION_END_CURRENT_ROW
	}

	// EXCLUDE clause
	if p.cur.Type == EXCLUDE {
		p.advance()
		switch p.cur.Type {
		case CURRENT_P:
			p.advance()
			p.expect(ROW)
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_CURRENT_ROW
		case GROUP_P:
			p.advance()
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_GROUP
		case TIES:
			p.advance()
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_TIES
		case NO:
			p.advance()
			p.expect(OTHERS)
			// No exclusion - default
		}
	}

	wd.FrameOptions = frameOptions
}

// parseFrameBound parses a single frame bound (start or end).
func (p *Parser) parseFrameBound(opts int, isStart bool) (int, nodes.Node) {
	if p.cur.Type == UNBOUNDED {
		p.advance()
		if p.cur.Type == PRECEDING {
			p.advance()
			if isStart {
				opts |= nodes.FRAMEOPTION_START_UNBOUNDED_PRECEDING
			} else {
				opts |= nodes.FRAMEOPTION_END_UNBOUNDED_PRECEDING
			}
		} else if p.cur.Type == FOLLOWING {
			p.advance()
			if isStart {
				opts |= nodes.FRAMEOPTION_START_UNBOUNDED_FOLLOWING
			} else {
				opts |= nodes.FRAMEOPTION_END_UNBOUNDED_FOLLOWING
			}
		}
		return opts, nil
	}

	if p.cur.Type == CURRENT_P {
		p.advance()
		p.expect(ROW)
		if isStart {
			opts |= nodes.FRAMEOPTION_START_CURRENT_ROW
		} else {
			opts |= nodes.FRAMEOPTION_END_CURRENT_ROW
		}
		return opts, nil
	}

	// expression PRECEDING/FOLLOWING
	expr := p.parseAExpr(0)
	if p.cur.Type == PRECEDING {
		p.advance()
		if isStart {
			opts |= nodes.FRAMEOPTION_START_OFFSET_PRECEDING
		} else {
			opts |= nodes.FRAMEOPTION_END_OFFSET_PRECEDING
		}
	} else if p.cur.Type == FOLLOWING {
		p.advance()
		if isStart {
			opts |= nodes.FRAMEOPTION_START_OFFSET_FOLLOWING
		} else {
			opts |= nodes.FRAMEOPTION_END_OFFSET_FOLLOWING
		}
	}
	return opts, expr
}

// ---------------------------------------------------------------------------
// AexprConst type-casted string constants

// isAExprConstTypeCast checks if we're looking at a type-casted constant
// like: int '42' or interval '1 day' or typename 'literal'.
func (p *Parser) isAExprConstTypeCast() bool {
	// ConstTypename Sconst: Numeric/Bit/Character/DateTime/JSON types followed by SCONST
	switch p.cur.Type {
	case INT_P, INTEGER, SMALLINT, BIGINT, REAL, FLOAT_P, DOUBLE_P,
		DECIMAL_P, DEC, NUMERIC,
		BIT, CHARACTER, CHAR_P, VARCHAR, NATIONAL, NCHAR,
		BOOLEAN_P, JSON,
		TIMESTAMP, TIME, INTERVAL:
		return true
	}
	return false
}

// parseTypeCastedConst parses ConstTypename Sconst or ConstInterval Sconst opt_interval.
func (p *Parser) parseTypeCastedConst() nodes.Node {
	if p.cur.Type == INTERVAL {
		return p.parseIntervalConst()
	}

	tn, err := p.parseSimpleTypename()
	if err != nil {
		return nil
	}

	if p.cur.Type == SCONST {
		tok := p.advance()
		return &nodes.TypeCast{
			Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
			TypeName: tn,
			Loc: nodes.NoLoc(),
		}
	}

	// Not a type-casted constant; this shouldn't happen if isAExprConstTypeCast was correct
	return nil
}

// parseIntervalConst parses interval 'literal' [qualifier].
func (p *Parser) parseIntervalConst() nodes.Node {
	tn, err := p.parseIntervalType()
	if err != nil {
		return nil
	}

	if p.cur.Type == SCONST {
		tok := p.advance()
		// Parse optional interval qualifier after the string
		optInterval := p.parseOptInterval()
		if optInterval != nil {
			tn.Typmods = optInterval
		}
		return &nodes.TypeCast{
			Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
			TypeName: tn,
			Loc: nodes.NoLoc(),
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Expression list helpers

// parseExprListFull parses a comma-separated list of a_expr.
//
//	expr_list: a_expr | expr_list ',' a_expr
func (p *Parser) parseExprListFull() *nodes.List {
	first := p.parseAExpr(0)
	if first == nil {
		return &nodes.List{}
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		expr := p.parseAExpr(0)
		if expr == nil {
			break
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}
}

// parseFuncArgExpr parses a single function argument (supports named args).
//
//	func_arg_expr:
//	    a_expr
//	    | param_name COLON_EQUALS a_expr
//	    | param_name EQUALS_GREATER a_expr
func (p *Parser) parseFuncArgExpr() nodes.Node {
	// Check for named argument: name := expr or name => expr
	if p.isColId() {
		next := p.peekNext()
		if next.Type == COLON_EQUALS || next.Type == EQUALS_GREATER {
			nameTok := p.advance() // consume name
			p.advance()            // consume := or =>
			arg := p.parseAExpr(0)
			return &nodes.NamedArgExpr{
				Name:      nameTok.Str,
				Arg:       arg,
				Argnumber: -1,
				Loc: nodes.Loc{Start: nameTok.Loc, End: -1},
			}
		}
	}
	return p.parseAExpr(0)
}

// parseFuncArgListFull parses a comma-separated list of func_arg_expr.
func (p *Parser) parseFuncArgListFull() *nodes.List {
	first := p.parseFuncArgExpr()
	if first == nil {
		return &nodes.List{}
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		// Check if next is VARIADIC (stop here, caller handles it)
		next := p.peekNext()
		if next.Type == VARIADIC {
			break
		}
		p.advance()
		expr := p.parseFuncArgExpr()
		if expr == nil {
			break
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}
}

// ---------------------------------------------------------------------------
// Subquery helpers - delegates to select.go

// isSelectStart returns true if the current token could start a SELECT statement.
func (p *Parser) isSelectStart() bool {
	switch p.cur.Type {
	case SELECT, VALUES, WITH, TABLE:
		return true
	}
	return false
}

// parseSelectStmtForExpr parses a select statement in expression context.
// Delegates to the full SELECT parser in select.go.
func (p *Parser) parseSelectStmtForExpr() nodes.Node {
	return p.parseSelectNoParens()
}

// parseTargetList parses a comma-separated target list (SELECT expressions).
func (p *Parser) parseTargetList() *nodes.List {
	first := p.parseTargetEl()
	if first == nil {
		return nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		el := p.parseTargetEl()
		if el == nil {
			break
		}
		items = append(items, el)
	}
	return &nodes.List{Items: items}
}

// parseTargetEl parses a single target list element.
func (p *Parser) parseTargetEl() nodes.Node {
	loc := p.pos()
	if p.cur.Type == '*' {
		p.advance()
		end := p.pos()
		return &nodes.ResTarget{
			Val: &nodes.ColumnRef{
				Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
				Loc:    nodes.Loc{Start: loc, End: end},
			},
			Loc: nodes.Loc{Start: loc, End: end},
		}
	}

	expr := p.parseAExpr(0)
	if expr == nil {
		return nil
	}

	rt := &nodes.ResTarget{
		Val: expr,
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional alias: AS ColLabel | ColId (without AS)
	if p.cur.Type == AS {
		aliasLoc := p.pos()
		p.advance()
		label, err := p.parseColLabel()
		if err == nil {
			rt.Name = label
			if p.completing {
				p.addSelectAliasPosition(aliasLoc)
			}
		}
	} else if p.isColId() && !p.isReservedForClause() {
		aliasLoc := p.pos()
		label, err := p.parseColId()
		if err == nil {
			rt.Name = label
			if p.completing {
				p.addSelectAliasPosition(aliasLoc)
			}
		}
	}

	rt.Loc.End = p.pos()
	return rt
}

// isReservedForClause checks if the current token is a keyword that starts a clause
// and should not be consumed as an alias.
func (p *Parser) isReservedForClause() bool {
	switch p.cur.Type {
	case FROM, WHERE, GROUP_P, HAVING, ORDER, LIMIT, OFFSET, FETCH,
		FOR, UNION, INTERSECT, EXCEPT, INTO, WINDOW:
		return true
	}
	return false
}

// parseFromList is a compatibility wrapper that calls parseFromListFull.
func (p *Parser) parseFromList() *nodes.List {
	return p.parseFromListFull()
}

// isJoinKeyword checks if the current token is a join-related keyword.
func (p *Parser) isJoinKeyword() bool {
	switch p.cur.Type {
	case JOIN, CROSS, LEFT, RIGHT, FULL, INNER_P, NATURAL, ON:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Helper functions

// makeSimpleAExpr creates an A_Expr with a simple string operator name.
func makeSimpleAExpr(kind nodes.A_Expr_Kind, op string, lexpr, rexpr nodes.Node) nodes.Node {
	return &nodes.A_Expr{
		Kind:  kind,
		Name:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr: lexpr,
		Rexpr: rexpr,
	}
}

// makeBoolExpr creates a BoolExpr node.
func makeBoolExpr(boolop nodes.BoolExprType, arg1, arg2 nodes.Node) nodes.Node {
	be := &nodes.BoolExpr{
		Boolop: boolop,
		Args:   &nodes.List{},
	}
	if arg1 != nil {
		be.Args.Items = append(be.Args.Items, arg1)
	}
	if arg2 != nil {
		be.Args.Items = append(be.Args.Items, arg2)
	}
	return be
}

// setNodeLoc sets the Loc field on any node that has one, using type assertions.
func setNodeLoc(n nodes.Node, start, end int) {
	loc := nodes.Loc{Start: start, End: end}
	switch v := n.(type) {
	case *nodes.A_Expr:
		v.Loc = loc
	case *nodes.BoolExpr:
		v.Loc = loc
	case *nodes.NullTest:
		v.Loc = loc
	case *nodes.BooleanTest:
		v.Loc = loc
	case *nodes.ColumnRef:
		v.Loc = loc
	case *nodes.FuncCall:
		v.Loc = loc
	case *nodes.A_Const:
		v.Loc = loc
	case *nodes.A_ArrayExpr:
		v.Loc = loc
	case *nodes.CoalesceExpr:
		v.Loc = loc
	case *nodes.GroupingFunc:
		v.Loc = loc
	case *nodes.TypeCast:
		v.Loc = loc
	case *nodes.SubLink:
		v.Loc = loc
	case *nodes.CollateClause:
		v.Loc = loc
	case *nodes.ParamRef:
		v.Loc = loc
	case *nodes.CaseExpr:
		v.Loc = loc
	case *nodes.RowExpr:
		v.Loc = loc
	case *nodes.MinMaxExpr:
		v.Loc = loc
	case *nodes.XmlExpr:
		v.Loc = loc
	case *nodes.SQLValueFunction:
		v.Loc = loc
	case *nodes.RangeVar:
		v.Loc = loc
	case *nodes.SetToDefault:
		v.Loc = loc
	}
}

// doNegate negates a numeric constant in place or creates a unary minus A_Expr.
func doNegate(n nodes.Node) nodes.Node {
	if ac, ok := n.(*nodes.A_Const); ok {
		if i, ok := ac.Val.(*nodes.Integer); ok {
			i.Ival = -i.Ival
			return n
		}
		if f, ok := ac.Val.(*nodes.Float); ok {
			if len(f.Fval) > 0 && f.Fval[0] == '-' {
				f.Fval = f.Fval[1:]
			} else {
				f.Fval = "-" + f.Fval
			}
			return n
		}
	}
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_OP,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "-"}}},
		Rexpr:    n,
	}
}

// makeFuncName creates a qualified function name list.
func makeFuncName(parts ...string) *nodes.List {
	items := make([]nodes.Node, len(parts))
	for i, part := range parts {
		items[i] = &nodes.String{Str: part}
	}
	return &nodes.List{Items: items}
}
