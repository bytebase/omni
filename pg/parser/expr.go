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
func (p *Parser) parseAExpr(minPrec int) (nodes.Node, error) {
	loc := p.pos()
	left, err := p.parseAExprAtom()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}
	for {
		if p.collectMode() {
			// Emit all infix operators with precedence >= minPrec.
			p.addInfixCandidates(minPrec)
			return left, nil
		}
		prec := p.aExprInfixPrec()
		if prec < minPrec || prec == precNone {
			break
		}
		left, err = p.parseAExprInfix(left, prec)
		if err != nil {
			return nil, err
		}
		if left == nil {
			return nil, nil
		}
		setNodeLoc(left, loc, p.pos())
	}
	return left, nil
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
func (p *Parser) parseAExprAtom() (nodes.Node, error) {
	loc := p.pos()
	switch p.cur.Type {
	case NOT:
		p.advance()
		arg, err := p.parseAExpr(precNot)
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		n := &nodes.BoolExpr{
			Boolop: nodes.NOT_EXPR,
			Args:   &nodes.List{Items: []nodes.Node{arg}},
			Loc:    nodes.Loc{Start: loc, End: p.pos()},
		}
		return n, nil
	case '+':
		// Unary plus: just return the operand (yacc does the same)
		p.advance()
		return p.parseAExpr(precUnary)
	case '-':
		// Unary minus
		p.advance()
		arg, err := p.parseAExpr(precUnary)
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		n := doNegate(arg)
		setNodeLoc(n, loc, p.pos())
		return n, nil
	default:
		return p.parseCExpr()
	}
}

// parseAExprInfix handles an infix or postfix operator in a_expr.
func (p *Parser) parseAExprInfix(left nodes.Node, prec int) (nodes.Node, error) {
	switch p.cur.Type {
	case OR:
		p.advance()
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeBoolExpr(nodes.OR_EXPR, left, right), nil
	case AND:
		p.advance()
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeBoolExpr(nodes.AND_EXPR, left, right), nil

	// --- IS / ISNULL / NOTNULL ---
	case IS:
		return p.parseIsPostfix(left)
	case ISNULL:
		p.advance()
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NULL}, nil
	case NOTNULL:
		p.advance()
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NOT_NULL}, nil

	// --- Comparison operators ---
	case '<', '>', '=':
		tok := p.advance()
		opStr := string(rune(tok.Type))
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, opStr, left, right), nil
	case LESS_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: "<="}}}
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, "<=", left, right), nil
	case GREATER_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: ">="}}}
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, ">=", left, right), nil
	case NOT_EQUALS:
		p.advance()
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: "<>"}}}
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, "<>", left, right), nil

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
			return nil, p.syntaxErrorAtCur() // shouldn't happen due to NOT_LA conversion logic
		}

	// --- Custom operators (Op, OPERATOR()) ---
	case Op:
		tok := p.advance()
		opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: tok.Str}}}
		// Check for subquery_Op sub_type pattern: expr Op ANY/ALL/SOME (subquery_or_expr)
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.A_Expr{Kind: nodes.AEXPR_OP, Name: opName, Lexpr: left, Rexpr: right}, nil

	// --- Arithmetic ---
	case '+', '-':
		tok := p.advance()
		opStr := string(rune(tok.Type))
		// Check for subquery_Op sub_type
		if p.cur.Type == ANY || p.cur.Type == SOME || p.cur.Type == ALL {
			opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
			return p.parseSubqueryOp(left, opName)
		}
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, opStr, left, right), nil
	case '*', '/', '%':
		tok := p.advance()
		right, err := p.parseAExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right), nil
	case '^':
		p.advance()
		right, err := p.parseAExpr(prec + 1) // left-assoc in PostgreSQL
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, "^", left, right), nil

	// --- AT TIME ZONE ---
	case AT:
		return p.parseAtTimeZone(left)

	// --- COLLATE ---
	case COLLATE:
		p.advance()
		collname, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		return &nodes.CollateClause{Arg: left, Collname: collname}, nil

	// --- TYPECAST (::) ---
	case TYPECAST:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		return &nodes.TypeCast{Arg: left, TypeName: tn, Loc: nodes.NoLoc()}, nil

	// --- Array subscript ---
	case '[':
		return p.parseSubscript(left)
	}

	return left, nil
}

// parseIsPostfix handles IS ... postfix expressions on left.
func (p *Parser) parseIsPostfix(left nodes.Node) (nodes.Node, error) {
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
			return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NOT_NULL}, nil
		}
		return &nodes.NullTest{Arg: left, Nulltesttype: nodes.IS_NULL}, nil
	case TRUE_P:
		p.advance()
		bt := nodes.IS_TRUE
		if negated {
			bt = nodes.IS_NOT_TRUE
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}, nil
	case FALSE_P:
		p.advance()
		bt := nodes.IS_FALSE
		if negated {
			bt = nodes.IS_NOT_FALSE
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}, nil
	case UNKNOWN:
		p.advance()
		bt := nodes.IS_UNKNOWN
		if negated {
			bt = nodes.IS_NOT_UNKNOWN
		}
		return &nodes.BooleanTest{Arg: left, Booltesttype: bt}, nil
	case DISTINCT:
		// IS [NOT] DISTINCT FROM a_expr
		p.advance() // consume DISTINCT
		if _, err := p.expect(FROM); err != nil {
			return nil, err
		}
		right, err := p.parseAExpr(precIs + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		kind := nodes.AEXPR_DISTINCT
		if negated {
			kind = nodes.AEXPR_NOT_DISTINCT
		}
		return makeSimpleAExpr(kind, "=", left, right), nil
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
			}, nil
		}
		return xmlExpr, nil
	case JSON:
		// IS [NOT] JSON [VALUE|ARRAY|OBJECT|SCALAR] [WITH UNIQUE [KEYS]]
		return p.parseJsonIsPredicate(left, negated), nil
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
			}, nil
		}
		return fc, nil
	case NFC, NFD, NFKC, NFKD:
		// IS [NOT] NFC/NFD/NFKC/NFKD NORMALIZED
		form := p.parseUnicodeNormalForm()
		if _, err := p.expect(NORMALIZED); err != nil {
			return nil, err
		}
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
			}, nil
		}
		return fc, nil
	default:
		return left, nil
	}
}

// parseBetweenExpr parses BETWEEN [SYMMETRIC] b_expr AND a_expr.
func (p *Parser) parseBetweenExpr(left nodes.Node, negated bool) (nodes.Node, error) {
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
	lower, err := p.parseBExpr(0)
	if err != nil {
		return nil, err
	}
	if lower == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(AND); err != nil {
		return nil, err
	}
	upper, err := p.parseAExpr(precIn + 1)
	if err != nil {
		return nil, err
	}
	if upper == nil {
		return nil, p.syntaxErrorAtCur()
	}

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
	}, nil
}

// parseInExpr parses IN (expr_list) or IN (subquery).
func (p *Parser) parseInExpr(left nodes.Node, negated bool) (nodes.Node, error) {
	p.advance() // consume IN

	op := "="
	if negated {
		op = "<>"
	}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume '('

	// Try to determine if this is a subquery or expression list.
	// Subqueries start with SELECT, VALUES, WITH, or '(' followed by another select.
	if p.isSelectStart() {
		subquery, err := p.parseSelectStmtForExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
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
			}, nil
		}
		return sub, nil
	}

	// Expression list
	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_IN,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    exprs,
	}, nil
}

// parseLikeExpr parses [NOT] LIKE a_expr [ESCAPE a_expr].
func (p *Parser) parseLikeExpr(left nodes.Node, negated bool) (nodes.Node, error) {
	p.advance() // consume LIKE
	right, err := p.parseAExpr(precIn + 1)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.syntaxErrorAtCur()
	}

	op := "~~"
	if negated {
		op = "!~~"
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr, err := p.parseAExpr(precIn + 1)
		if err != nil {
			return nil, err
		}
		if escapeExpr == nil {
			return nil, p.syntaxErrorAtCur()
		}
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
		}, nil
	}

	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_LIKE,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    right,
	}, nil
}

// parseIlikeExpr parses [NOT] ILIKE a_expr [ESCAPE a_expr].
func (p *Parser) parseIlikeExpr(left nodes.Node, negated bool) (nodes.Node, error) {
	p.advance() // consume ILIKE
	right, err := p.parseAExpr(precIn + 1)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.syntaxErrorAtCur()
	}

	op := "~~*"
	if negated {
		op = "!~~*"
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr, err := p.parseAExpr(precIn + 1)
		if err != nil {
			return nil, err
		}
		if escapeExpr == nil {
			return nil, p.syntaxErrorAtCur()
		}
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
		}, nil
	}

	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_ILIKE,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr:    left,
		Rexpr:    right,
	}, nil
}

// parseSimilarExpr parses [NOT] SIMILAR TO a_expr [ESCAPE a_expr].
func (p *Parser) parseSimilarExpr(left nodes.Node, negated bool) (nodes.Node, error) {
	p.advance() // consume SIMILAR
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	right, err := p.parseAExpr(precIn + 1)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if negated {
		if p.cur.Type == ESCAPE {
			p.advance()
			escapeExpr, err := p.parseAExpr(precIn + 1)
			if err != nil {
				return nil, err
			}
			if escapeExpr == nil {
				return nil, p.syntaxErrorAtCur()
			}
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
			}, nil
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
		}, nil
	}

	if p.cur.Type == ESCAPE {
		p.advance()
		escapeExpr, err := p.parseAExpr(precIn + 1)
		if err != nil {
			return nil, err
		}
		if escapeExpr == nil {
			return nil, p.syntaxErrorAtCur()
		}
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
		}, nil
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
	}, nil
}

// parseSubqueryOp handles expr op ANY/ALL/SOME (subquery_or_expr).
func (p *Parser) parseSubqueryOp(left nodes.Node, opName *nodes.List) (nodes.Node, error) {
	subType := nodes.ANY_SUBLINK
	if p.cur.Type == ALL {
		subType = nodes.ALL_SUBLINK
	}
	p.advance() // consume ANY/ALL/SOME

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if p.isSelectStart() {
		subquery, err := p.parseSelectStmtForExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.SubLink{
			SubLinkType: int(subType),
			Testexpr:    left,
			OperName:    opName,
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}, nil
	}

	// ANY/ALL (expr) -- array form
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	kind := nodes.AEXPR_OP_ANY
	if subType == nodes.ALL_SUBLINK {
		kind = nodes.AEXPR_OP_ALL
	}
	return &nodes.A_Expr{
		Kind:     kind,
		Name:     opName,
		Lexpr:    left,
		Rexpr:    expr,
	}, nil
}

// parseAtTimeZone handles AT TIME ZONE a_expr and AT LOCAL.
func (p *Parser) parseAtTimeZone(left nodes.Node) (nodes.Node, error) {
	p.advance() // consume AT

	if p.cur.Type == LOCAL {
		// AT LOCAL => AT TIME ZONE 'DEFAULT' (handled as function call)
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "timezone"),
			Args:       &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.String{Str: "DEFAULT"}}, left}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	if _, err := p.expect(TIME); err != nil {
		return nil, err
	}
	if _, err := p.expect(ZONE); err != nil {
		return nil, err
	}
	tz, err := p.parseAExpr(precAt + 1)
	if err != nil {
		return nil, err
	}
	if tz == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "timezone"),
		Args:       &nodes.List{Items: []nodes.Node{tz, left}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
}

// parseSubscript handles array subscript and slice: expr[idx] or expr[lo:hi].
func (p *Parser) parseSubscript(left nodes.Node) (nodes.Node, error) {
	p.advance() // consume '['

	// Check for empty lower bound followed by ':'
	if p.cur.Type == ':' {
		p.advance()
		var uidx nodes.Node
		if p.cur.Type != ']' {
			var err error
			uidx, err = p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(']'); err != nil {
			return nil, err
		}
		return &nodes.A_Indirection{
			Arg:         left,
			Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{IsSlice: true, Uidx: uidx}}},
		}, nil
	}

	idx, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ':' {
		// Slice
		p.advance()
		var uidx nodes.Node
		if p.cur.Type != ']' {
			uidx, err = p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(']'); err != nil {
			return nil, err
		}
		return &nodes.A_Indirection{
			Arg:         left,
			Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{IsSlice: true, Lidx: idx, Uidx: uidx}}},
		}, nil
	}

	// Simple subscript
	if _, err := p.expect(']'); err != nil {
		return nil, err
	}
	return &nodes.A_Indirection{
		Arg:         left,
		Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Indices{Uidx: idx}}},
	}, nil
}

// ---------------------------------------------------------------------------
// parseBExpr parses a restricted expression (no boolean, no IN/LIKE/BETWEEN, no IS NULL etc).
//
//	b_expr: c_expr | b_expr op b_expr | b_expr IS [NOT] DISTINCT FROM b_expr | ...
func (p *Parser) parseBExpr(minPrec int) (nodes.Node, error) {
	left, err := p.parseBExprAtom()
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, nil
	}
	for {
		prec := p.bExprInfixPrec()
		if prec < minPrec || prec == precNone {
			break
		}
		left, err = p.parseBExprInfix(left, prec)
		if err != nil {
			return nil, err
		}
		if left == nil {
			return nil, nil
		}
	}
	return left, nil
}

// parseBExprAtom handles prefix operators for b_expr.
func (p *Parser) parseBExprAtom() (nodes.Node, error) {
	switch p.cur.Type {
	case '+':
		p.advance()
		return p.parseBExpr(precUnary)
	case '-':
		p.advance()
		arg, err := p.parseBExpr(precUnary)
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return doNegate(arg), nil
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
func (p *Parser) parseBExprInfix(left nodes.Node, prec int) (nodes.Node, error) {
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
			if _, err := p.expect(FROM); err != nil {
				return nil, err
			}
			right, err := p.parseBExpr(precIs + 1)
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, p.syntaxErrorAtCur()
			}
			kind := nodes.AEXPR_DISTINCT
			if negated {
				kind = nodes.AEXPR_NOT_DISTINCT
			}
			return makeSimpleAExpr(kind, "=", left, right), nil
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
				}, nil
			}
			return xmlExpr, nil
		}
		return left, nil
	case '<', '>', '=':
		tok := p.advance()
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right), nil
	case LESS_EQUALS:
		p.advance()
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, "<=", left, right), nil
	case GREATER_EQUALS:
		p.advance()
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, ">=", left, right), nil
	case NOT_EQUALS:
		p.advance()
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, "<>", left, right), nil
	case Op:
		tok := p.advance()
		opName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: tok.Str}}}
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.A_Expr{Kind: nodes.AEXPR_OP, Name: opName, Lexpr: left, Rexpr: right}, nil
	case '+', '-', '*', '/', '%', '^':
		tok := p.advance()
		right, err := p.parseBExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return makeSimpleAExpr(nodes.AEXPR_OP, string(rune(tok.Type)), left, right), nil
	case TYPECAST:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		if tn == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.TypeCast{Arg: left, TypeName: tn, Loc: nodes.NoLoc()}, nil
	}
	return left, nil
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
func (p *Parser) parseCExpr() (nodes.Node, error) {
	loc := p.pos()
	n, err := p.parseCExprInner()
	if err != nil {
		return nil, err
	}
	if n != nil {
		setNodeLoc(n, loc, p.pos())
	}
	return n, nil
}

func (p *Parser) parseCExprInner() (nodes.Node, error) {
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
		return nil, errCollecting
	}
	switch p.cur.Type {
	// --- Constants ---
	case ICONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.Integer{Ival: tok.Ival}}, nil
	case FCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.Float{Fval: tok.Str}}, nil
	case SCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.String{Str: tok.Str}}, nil
	case BCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.BitString{Bsval: tok.Str}}, nil
	case XCONST:
		tok := p.advance()
		return &nodes.A_Const{Val: &nodes.BitString{Bsval: tok.Str}}, nil
	case TRUE_P:
		p.advance()
		return &nodes.A_Const{Val: &nodes.Boolean{Boolval: true}}, nil
	case FALSE_P:
		p.advance()
		return &nodes.A_Const{Val: &nodes.Boolean{Boolval: false}}, nil
	case NULL_P:
		p.advance()
		return &nodes.A_Const{Isnull: true}, nil

	// --- Parameter ($1, $2, ...) ---
	case PARAM:
		tok := p.advance()
		ref := &nodes.ParamRef{Number: int(tok.Ival)}
		indir, err := p.parseOptIndirection()
		if err != nil {
			return nil, err
		}
		if indir != nil && len(indir.Items) > 0 {
			return &nodes.A_Indirection{Arg: ref, Indirection: indir}, nil
		}
		return ref, nil

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
		return &nodes.SetToDefault{Loc: nodes.NoLoc()}, nil

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
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_DATE, Typmod: -1, Loc: nodes.NoLoc()}, nil
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
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_ROLE, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case CURRENT_USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_USER, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case SESSION_USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_SESSION_USER, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case USER:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_USER, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case CURRENT_CATALOG:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_CATALOG, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case CURRENT_SCHEMA:
		p.advance()
		return &nodes.SQLValueFunction{Op: nodes.SVFOP_CURRENT_SCHEMA, Typmod: -1, Loc: nodes.NoLoc()}, nil
	case SYSTEM_USER:
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "system_user"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil

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

	return nil, nil
}

// parseParenExprOrRow handles '(' ... ')' which could be:
// - Parenthesized expression: (a_expr) opt_indirection
// - Implicit row: (a_expr, a_expr, ...)
// - Subquery: (SELECT ...)
func (p *Parser) parseParenExprOrRow() (nodes.Node, error) {
	p.advance() // consume '('

	// Check for subquery
	if p.isSelectStart() {
		subquery, err := p.parseSelectStmtForExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// Check for indirection after subquery
		if p.cur.Type == '.' || p.cur.Type == '[' {
			indir, err := p.parseIndirection()
			if err == nil && indir != nil && len(indir.Items) > 0 {
				sub := &nodes.SubLink{
					SubLinkType: int(nodes.EXPR_SUBLINK),
					Subselect:   subquery,
					Loc: nodes.NoLoc(),
				}
				return &nodes.A_Indirection{Arg: sub, Indirection: indir}, nil
			}
		}
		return &nodes.SubLink{
			SubLinkType: int(nodes.EXPR_SUBLINK),
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Parse first expression
	first, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ',' {
		// Implicit row: (expr, expr, ...)
		items := []nodes.Node{first}
		for p.cur.Type == ',' {
			p.advance()
			expr, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			items = append(items, expr)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if len(items) < 2 {
			// Single element + trailing comma should not happen but handle gracefully
			return first, nil
		}
		return &nodes.RowExpr{
			Args:      &nodes.List{Items: items},
			RowFormat: nodes.COERCE_IMPLICIT_CAST,
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Parenthesized expression
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	// Check for indirection
	indir, err := p.parseOptIndirection()
	if err != nil {
		return nil, err
	}
	if indir != nil && len(indir.Items) > 0 {
		return &nodes.A_Indirection{Arg: first, Indirection: indir}, nil
	}
	return first, nil
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() (nodes.Node, error) {
	p.advance() // consume EXISTS
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	subquery, err := p.parseSelectStmtForExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.SubLink{
		SubLinkType: int(nodes.EXISTS_SUBLINK),
		Subselect:   subquery,
		Loc: nodes.NoLoc(),
	}, nil
}

// parseArrayCExpr parses ARRAY [...] or ARRAY (subquery).
func (p *Parser) parseArrayCExpr() (nodes.Node, error) {
	p.advance() // consume ARRAY

	if p.cur.Type == '[' {
		arr, err := p.parseArrayExpr()
		if err != nil {
			return nil, err
		}
		if arr != nil {
			if a, ok := arr.(*nodes.A_ArrayExpr); ok {
				a.Loc = nodes.NoLoc()
			}
		}
		return arr, nil
	}

	if p.cur.Type == '(' {
		p.advance()
		subquery, err := p.parseSelectStmtForExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.SubLink{
			SubLinkType: int(nodes.ARRAY_SUBLINK),
			Subselect:   subquery,
			Loc: nodes.NoLoc(),
		}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// parseArrayExpr parses array literal: '[' expr_list ']' or '[' array_expr_list ']' or '[]'.
//
//	array_expr:
//	    '[' expr_list ']'
//	    | '[' array_expr_list ']'
//	    | '[' ']'
func (p *Parser) parseArrayExpr() (nodes.Node, error) {
	p.advance() // consume '['

	if p.cur.Type == ']' {
		p.advance()
		return &nodes.A_ArrayExpr{Loc: nodes.NoLoc()}, nil
	}

	// Check if first element is a nested array
	if p.cur.Type == '[' {
		// array_expr_list
		first, err := p.parseArrayExpr()
		if err != nil {
			return nil, err
		}
		items := []nodes.Node{first}
		for p.cur.Type == ',' {
			p.advance()
			elem, err := p.parseArrayExpr()
			if err != nil {
				return nil, err
			}
			items = append(items, elem)
		}
		if _, err := p.expect(']'); err != nil {
			return nil, err
		}
		return &nodes.A_ArrayExpr{Elements: &nodes.List{Items: items}}, nil
	}

	// expr_list
	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(']'); err != nil {
		return nil, err
	}
	return &nodes.A_ArrayExpr{Elements: exprs}, nil
}

// parseCaseExpr parses CASE ... END.
//
// Ref: https://www.postgresql.org/docs/17/sql-expressions.html#SYNTAX-EXPRESSION-EVAL
//
//	CASE [expr] WHEN expr THEN expr [WHEN ...] [ELSE expr] END
func (p *Parser) parseCaseExpr() (nodes.Node, error) {
	p.advance() // consume CASE

	// Optional case argument (simple CASE)
	var caseArg nodes.Node
	if p.cur.Type != WHEN {
		var err error
		caseArg, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
	}

	// Parse WHEN clauses
	var whens []nodes.Node
	for p.cur.Type == WHEN {
		p.advance() // consume WHEN
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(THEN); err != nil {
			return nil, err
		}
		result, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
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
		var err error
		defResult, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(END_P); err != nil {
		return nil, err
	}

	return &nodes.CaseExpr{
		Arg:       caseArg,
		Args:      &nodes.List{Items: whens},
		Defresult: defResult,
		Loc: nodes.NoLoc(),
	}, nil
}

// parseExplicitRow parses ROW(...).
//
//	explicit_row: ROW '(' expr_list ')' | ROW '(' ')'
func (p *Parser) parseExplicitRow() (nodes.Node, error) {
	p.advance() // consume ROW
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.RowExpr{
			RowFormat: nodes.COERCE_EXPLICIT_CALL,
			Loc: nodes.NoLoc(),
		}, nil
	}

	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.RowExpr{
		Args:      exprs,
		RowFormat: nodes.COERCE_EXPLICIT_CALL,
		Loc: nodes.NoLoc(),
	}, nil
}

// ---------------------------------------------------------------------------
// func_expr_common_subexpr implementations

// parseCastExpr parses CAST(a_expr AS Typename).
func (p *Parser) parseCastExpr() (nodes.Node, error) {
	p.advance() // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	arg, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.TypeCast{Arg: arg, TypeName: tn, Loc: nodes.NoLoc()}, nil
}

// parseCoalesceExpr parses COALESCE(expr_list).
func (p *Parser) parseCoalesceExpr() (nodes.Node, error) {
	p.advance() // consume COALESCE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CoalesceExpr{Args: exprs}, nil
}

// parseMinMaxExpr parses GREATEST(expr_list) or LEAST(expr_list).
func (p *Parser) parseMinMaxExpr(op nodes.MinMaxOp) (nodes.Node, error) {
	p.advance() // consume GREATEST/LEAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.MinMaxExpr{Op: op, Args: exprs}, nil
}

// parseNullIfExpr parses NULLIF(a_expr, a_expr).
func (p *Parser) parseNullIfExpr() (nodes.Node, error) {
	p.advance() // consume NULLIF
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	left, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	right, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.A_Expr{
		Kind:     nodes.AEXPR_NULLIF,
		Name:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}},
		Lexpr:    left,
		Rexpr:    right,
	}, nil
}

// parseExtractExpr parses EXTRACT(field FROM a_expr).
func (p *Parser) parseExtractExpr() (nodes.Node, error) {
	p.advance() // consume EXTRACT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		// EXTRACT() with no args
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "extract"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	field := p.parseExtractArg()
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "extract"),
		Args:       &nodes.List{Items: []nodes.Node{&nodes.A_Const{Val: &nodes.String{Str: field}}, expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
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
func (p *Parser) parseNormalizeExpr() (nodes.Node, error) {
	p.advance() // consume NORMALIZE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ',' {
		p.advance()
		form := p.parseUnicodeNormalForm()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "normalize"),
			Args:       &nodes.List{Items: []nodes.Node{expr, &nodes.A_Const{Val: &nodes.String{Str: form}}}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "normalize"),
		Args:       &nodes.List{Items: []nodes.Node{expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
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
func (p *Parser) parseOverlayExpr() (nodes.Node, error) {
	p.advance() // consume OVERLAY
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Try SQL-standard syntax: a_expr PLACING a_expr FROM a_expr [FOR a_expr]
	first, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if p.cur.Type == PLACING {
		p.advance()
		second, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(FROM); err != nil {
			return nil, err
		}
		third, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		args := []nodes.Node{first, second, third}
		if p.cur.Type == FOR {
			p.advance()
			fourth, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			args = append(args, fourth)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "overlay"),
			Args:       &nodes.List{Items: args},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Function-call syntax: overlay(arg, arg, ...)
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("overlay"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}, nil
}

// parsePositionExpr parses POSITION(b_expr IN b_expr).
func (p *Parser) parsePositionExpr() (nodes.Node, error) {
	p.advance() // consume POSITION
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "position"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	first, err := p.parseBExpr(0)
	if err != nil {
		return nil, err
	}
	if p.cur.Type == IN_P {
		p.advance()
		second, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// Note: PG reverses the arguments
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "position"),
			Args:       &nodes.List{Items: []nodes.Node{second, first}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Function-call syntax
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("position"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc: nodes.NoLoc(),
	}, nil
}

// parseSubstringExpr parses SUBSTRING(... FROM ... [FOR ...]) or SUBSTRING(expr_list).
func (p *Parser) parseSubstringExpr() (nodes.Node, error) {
	p.advance() // consume SUBSTRING
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		p.advance()
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	first, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	// Check for SQL-standard syntax with FROM/FOR/SIMILAR
	if p.cur.Type == FROM {
		p.advance()
		second, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if p.cur.Type == FOR {
			p.advance()
			third, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "substring"),
				Args:       &nodes.List{Items: []nodes.Node{first, second, third}},
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Loc: nodes.NoLoc(),
			}, nil
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, second}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	if p.cur.Type == FOR {
		p.advance()
		second, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if p.cur.Type == FROM {
			p.advance()
			third, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "substring"),
				Args:       &nodes.List{Items: []nodes.Node{first, third, second}},
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Loc: nodes.NoLoc(),
			}, nil
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, makeIntConst(1), second}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	if p.cur.Type == SIMILAR {
		p.advance()
		second, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ESCAPE); err != nil {
			return nil, err
		}
		third, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname:   makeFuncName("pg_catalog", "substring"),
			Args:       &nodes.List{Items: []nodes.Node{first, second, third}},
			FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Comma-separated function call form
	args := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(args) > 1 {
		// Regular function call form: substring(str, start, len)
		return &nodes.FuncCall{
			Funcname:   makeFuncName("substring"),
			Args:       &nodes.List{Items: args},
			FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
			Loc: nodes.NoLoc(),
		}, nil
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "substring"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
}

// parseTrimExpr parses TRIM([BOTH|LEADING|TRAILING] [expr FROM] expr_list).
func (p *Parser) parseTrimExpr() (nodes.Node, error) {
	p.advance() // consume TRIM
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

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
	args, err := p.parseTrimList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", funcName),
		Args:       args,
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
}

// parseTrimList parses trim_list.
//
//	trim_list:
//	    a_expr FROM expr_list
//	    | FROM expr_list
//	    | expr_list
func (p *Parser) parseTrimList() (*nodes.List, error) {
	if p.cur.Type == FROM {
		p.advance()
		return p.parseExprListFull()
	}

	first, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if p.cur.Type == FROM {
		p.advance()
		rest, err := p.parseExprListFull()
		if err != nil {
			return nil, err
		}
		// Prepend the trim character
		items := make([]nodes.Node, 0, len(rest.Items)+1)
		items = append(items, first)
		items = append(items, rest.Items...)
		return &nodes.List{Items: items}, nil
	}

	// Plain expression list
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}, nil
}

// parseGroupingExpr parses GROUPING(expr_list).
func (p *Parser) parseGroupingExpr() (nodes.Node, error) {
	p.advance() // consume GROUPING
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	exprs, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.GroupingFunc{Args: exprs}, nil
}

// parseCollationForExpr parses COLLATION FOR (a_expr).
func (p *Parser) parseCollationForExpr() (nodes.Node, error) {
	p.advance() // consume COLLATION
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "pg_collation_for"),
		Args:       &nodes.List{Items: []nodes.Node{expr}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc: nodes.NoLoc(),
	}, nil
}

// parseTreatExpr parses TREAT(a_expr AS Typename).
func (p *Parser) parseTreatExpr() (nodes.Node, error) {
	p.advance() // consume TREAT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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
	}, nil
}

// parseSVFWithOptionalPrecision parses SQL value functions that take optional (int) precision.
func (p *Parser) parseSVFWithOptionalPrecision(noArgOp, withArgOp nodes.SVFOp) (nodes.Node, error) {
	p.advance() // consume keyword
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ICONST {
			prec := p.cur.Ival
			p.advance()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.SQLValueFunction{Op: withArgOp, Typmod: int32(prec), Loc: nodes.NoLoc()}, nil
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	return &nodes.SQLValueFunction{Op: noArgOp, Typmod: -1, Loc: nodes.NoLoc()}, nil
}

// ---------------------------------------------------------------------------
// Function call parsing (minimal for batch 3, expanded in batch 4)

// parseColumnRefOrFuncCall disambiguates between column reference and function call.
// If followed by '(', it's a function call. Otherwise, column reference.
func (p *Parser) parseColumnRefOrFuncCall() (nodes.Node, error) {
	// Check for AexprConst: func_name Sconst (type-casted literal)
	// This is handled by isAExprConstTypeCast / parseTypeCastedConst above.

	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil, err
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
		}, nil
	}

	// Check for function call
	if p.cur.Type == '(' {
		funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
		return p.parseFuncApplication(funcName, loc)
	}

	// Check for qualified name: schema.func(...) or schema.table.column
	if p.cur.Type == '.' {
		p.advance()
		// After "name.", we're in a column reference context.
		if p.collectMode() {
			p.addRuleCandidate("columnref")
		}
		if p.cur.Type == '*' {
			// schema.*
			p.advance()
			return &nodes.ColumnRef{
				Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.A_Star{}}},
			}, nil
		}

		attr, err := p.parseAttrName()
		if err != nil {
			return nil, err
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
			}, nil
		}

		// schema.func(...)
		if p.cur.Type == '(' {
			funcName := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}}
			return p.parseFuncApplication(funcName, loc)
		}

		// schema.table.column or schema.column
		if p.cur.Type == '.' {
			p.advance()
			// After "name.attr.", emit columnref for 3+ part references.
			if p.collectMode() {
				p.addRuleCandidate("columnref")
			}
			if p.cur.Type == '*' {
				p.advance()
				return &nodes.ColumnRef{
					Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.A_Star{}}},
				}, nil
			}
			attr2, err := p.parseAttrName()
			if err != nil {
				return nil, err
			}
			// Support 4+ part column references (e.g. db.schema.table.column)
			// by continuing to consume dot-separated names, matching PostgreSQL's
			// columnref: ColId indirection rule which allows arbitrary depth.
			fields := []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}, &nodes.String{Str: attr2}}
			for p.cur.Type == '.' {
				p.advance()
				if p.cur.Type == '*' {
					p.advance()
					fields = append(fields, &nodes.A_Star{})
					return &nodes.ColumnRef{
						Fields: &nodes.List{Items: fields},
					}, nil
				}
				attrN, err := p.parseAttrName()
				if err != nil {
					return nil, err
				}
				fields = append(fields, &nodes.String{Str: attrN})
			}
			return &nodes.ColumnRef{
				Fields: &nodes.List{Items: fields},
			}, nil
		}

		return &nodes.ColumnRef{
			Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}, &nodes.String{Str: attr}}},
		}, nil
	}

	// Simple column reference
	return &nodes.ColumnRef{
		Fields:   &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
	}, nil
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
func (p *Parser) parseFuncApplication(funcName *nodes.List, loc int) (nodes.Node, error) {
	p.advance() // consume '('

	fc := &nodes.FuncCall{
		Funcname: funcName,
	}

	// If cursor is inside function parens, emit expression candidates.
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		p.addTokenCandidate('*')
		p.addTokenCandidate(DISTINCT)
	}

	if p.cur.Type == ')' {
		p.advance()
	} else if p.cur.Type == '*' {
		p.advance()
		fc.AggStar = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == DISTINCT {
		p.advance()
		fc.AggDistinct = true
		var err error
		fc.Args, err = p.parseFuncArgListFull()
		if err != nil {
			return nil, err
		}
		fc.AggOrder, err = p.parseOptSortClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == ALL {
		p.advance()
		var err error
		fc.Args, err = p.parseFuncArgListFull()
		if err != nil {
			return nil, err
		}
		fc.AggOrder, err = p.parseOptSortClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == VARIADIC {
		p.advance()
		fc.FuncVariadic = true
		varArg, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = &nodes.List{Items: []nodes.Node{varArg}}
		fc.AggOrder, err = p.parseOptSortClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		// Regular argument list
		args, err := p.parseFuncArgListFull()
		if err != nil {
			return nil, err
		}

		if p.cur.Type == ',' && args != nil {
			// Check for VARIADIC after regular args
			p.advance()
			if p.cur.Type == VARIADIC {
				p.advance()
				fc.FuncVariadic = true
				variadicArg, err := p.parseFuncArgExpr()
				if err != nil {
					return nil, err
				}
				args.Items = append(args.Items, variadicArg)
			}
		}

		fc.Args = args
		fc.AggOrder, err = p.parseOptSortClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// within_group_clause: WITHIN GROUP '(' sort_clause ')'
	if p.cur.Type == WITHIN {
		p.advance()
		if _, err := p.expect(GROUP_P); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if p.cur.Type == ORDER {
			p.advance()
			if _, err := p.expect(BY); err != nil {
				return nil, err
			}
			var err error
			fc.AggOrder, err = p.parseSortByList()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.AggWithinGroup = true
	}

	// filter_clause: FILTER '(' WHERE a_expr ')'
	if p.cur.Type == FILTER {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect(WHERE); err != nil {
			return nil, err
		}
		var err error
		fc.AggFilter, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// over_clause: OVER window_specification | OVER ColId
	if p.cur.Type == OVER {
		var err error
		fc.Over, err = p.parseOverClause()
		if err != nil {
			return nil, err
		}
	}

	return fc, nil
}

// parseOptSortClause parses optional ORDER BY sort clause (used in aggregates).
// Returns nil if no ORDER BY present.
func (p *Parser) parseOptSortClause() (*nodes.List, error) {
	if p.cur.Type != ORDER {
		return nil, nil
	}
	p.advance() // consume ORDER
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}
	return p.parseSortByList()
}

// parseSortByList parses a comma-separated list of sortby items.
func (p *Parser) parseSortByList() (*nodes.List, error) {
	first, err := p.parseSortBy()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		sb, err := p.parseSortBy()
		if err != nil {
			return nil, err
		}
		items = append(items, sb)
	}
	return &nodes.List{Items: items}, nil
}

// parseSortBy parses a single sortby item: a_expr [ASC|DESC] [NULLS FIRST|LAST].
func (p *Parser) parseSortBy() (nodes.Node, error) {
	loc := p.pos()
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

// parseOverClause parses OVER window_specification or OVER existing_window_name.
//
// Ref: gram.y over_clause
//
//	OVER window_specification
//	| OVER ColId
func (p *Parser) parseOverClause() (nodes.Node, error) {
	p.advance() // consume OVER
	if p.cur.Type == '(' {
		return p.parseWindowSpecification()
	}
	// OVER ColId - references an existing window by name
	loc := p.pos()
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.WindowDef{
		Name:         name,
		FrameOptions: nodes.FRAMEOPTION_DEFAULTS,
		Loc:          nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseWindowSpecification parses a window specification: ( [name] [PARTITION BY ...] [ORDER BY ...] [frame] ).
func (p *Parser) parseWindowSpecification() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume '('
	wd := &nodes.WindowDef{Loc: nodes.Loc{Start: loc, End: -1}}

	// Optional existing window name
	if p.isColId() && p.cur.Type != PARTITION && p.cur.Type != ORDER && p.cur.Type != RANGE && p.cur.Type != ROWS && p.cur.Type != GROUPS {
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		wd.Refname = name
	}

	// PARTITION BY
	if p.cur.Type == PARTITION {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		var err error
		wd.PartitionClause, err = p.parseExprListFull()
		if err != nil {
			return nil, err
		}
	}

	// ORDER BY
	if p.cur.Type == ORDER {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		var err error
		wd.OrderClause, err = p.parseSortByList()
		if err != nil {
			return nil, err
		}
	}

	// Frame clause
	if p.cur.Type == RANGE || p.cur.Type == ROWS || p.cur.Type == GROUPS {
		if err := p.parseFrameClause(wd); err != nil {
			return nil, err
		}
	} else {
		wd.FrameOptions = nodes.FRAMEOPTION_DEFAULTS
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	wd.Loc.End = p.pos()
	return wd, nil
}

// parseFrameClause parses window frame specification.
func (p *Parser) parseFrameClause(wd *nodes.WindowDef) error {
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
		var err error
		frameOptions, wd.StartOffset, err = p.parseFrameBound(frameOptions, true)
		if err != nil {
			return err
		}
		if _, err := p.expect(AND); err != nil {
			return err
		}
		frameOptions, wd.EndOffset, err = p.parseFrameBound(frameOptions, false)
		if err != nil {
			return err
		}
	} else {
		var err error
		frameOptions, wd.StartOffset, err = p.parseFrameBound(frameOptions, true)
		if err != nil {
			return err
		}
		// Default end bound: CURRENT ROW
		frameOptions |= nodes.FRAMEOPTION_END_CURRENT_ROW
	}

	// EXCLUDE clause
	if p.cur.Type == EXCLUDE {
		p.advance()
		switch p.cur.Type {
		case CURRENT_P:
			p.advance()
			if _, err := p.expect(ROW); err != nil {
				return err
			}
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_CURRENT_ROW
		case GROUP_P:
			p.advance()
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_GROUP
		case TIES:
			p.advance()
			frameOptions |= nodes.FRAMEOPTION_EXCLUDE_TIES
		case NO:
			p.advance()
			if _, err := p.expect(OTHERS); err != nil {
				return err
			}
			// No exclusion - default
		}
	}

	wd.FrameOptions = frameOptions
	return nil
}

// parseFrameBound parses a single frame bound (start or end).
func (p *Parser) parseFrameBound(opts int, isStart bool) (int, nodes.Node, error) {
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
		return opts, nil, nil
	}

	if p.cur.Type == CURRENT_P {
		p.advance()
		if _, err := p.expect(ROW); err != nil {
			return 0, nil, err
		}
		if isStart {
			opts |= nodes.FRAMEOPTION_START_CURRENT_ROW
		} else {
			opts |= nodes.FRAMEOPTION_END_CURRENT_ROW
		}
		return opts, nil, nil
	}

	// expression PRECEDING/FOLLOWING
	expr, err := p.parseAExpr(0)
	if err != nil {
		return 0, nil, err
	}
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
	return opts, expr, nil
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
func (p *Parser) parseTypeCastedConst() (nodes.Node, error) {
	if p.cur.Type == INTERVAL {
		return p.parseIntervalConst()
	}

	tn, err := p.parseSimpleTypename()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == SCONST {
		tok := p.advance()
		return &nodes.TypeCast{
			Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
			TypeName: tn,
			Loc: nodes.NoLoc(),
		}, nil
	}

	// Not a type-casted constant; this shouldn't happen if isAExprConstTypeCast was correct
	return nil, p.syntaxErrorAtCur()
}

// parseIntervalConst parses interval 'literal' [qualifier].
func (p *Parser) parseIntervalConst() (nodes.Node, error) {
	tn, err := p.parseIntervalType()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == SCONST {
		tok := p.advance()
		// Parse optional interval qualifier after the string
		optInterval, err := p.parseOptInterval()
		if err != nil {
			return nil, err
		}
		if optInterval != nil {
			tn.Typmods = optInterval
		}
		return &nodes.TypeCast{
			Arg:      &nodes.A_Const{Val: &nodes.String{Str: tok.Str}},
			TypeName: tn,
			Loc: nodes.NoLoc(),
		}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// Expression list helpers

// parseExprListFull parses a comma-separated list of a_expr.
//
//	expr_list: a_expr | expr_list ',' a_expr
func (p *Parser) parseExprListFull() (*nodes.List, error) {
	first, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if first == nil {
		return &nodes.List{}, nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			break
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}, nil
}

// parseFuncArgExpr parses a single function argument (supports named args).
//
//	func_arg_expr:
//	    a_expr
//	    | param_name COLON_EQUALS a_expr
//	    | param_name EQUALS_GREATER a_expr
func (p *Parser) parseFuncArgExpr() (nodes.Node, error) {
	// Check for named argument: name := expr or name => expr
	if p.isColId() {
		next := p.peekNext()
		if next.Type == COLON_EQUALS || next.Type == EQUALS_GREATER {
			nameTok := p.advance() // consume name
			p.advance()            // consume := or =>
			arg, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			return &nodes.NamedArgExpr{
				Name:      nameTok.Str,
				Arg:       arg,
				Argnumber: -1,
				Loc: nodes.Loc{Start: nameTok.Loc, End: -1},
			}, nil
		}
	}
	return p.parseAExpr(0)
}

// parseFuncArgListFull parses a comma-separated list of func_arg_expr.
func (p *Parser) parseFuncArgListFull() (*nodes.List, error) {
	first, err := p.parseFuncArgExpr()
	if err != nil {
		return nil, err
	}
	if first == nil {
		return &nodes.List{}, nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		// Check if next is VARIADIC (stop here, caller handles it)
		next := p.peekNext()
		if next.Type == VARIADIC {
			break
		}
		p.advance()
		expr, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			break
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}, nil
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
func (p *Parser) parseSelectStmtForExpr() (nodes.Node, error) {
	return p.parseSelectNoParens()
}

// parseTargetList parses a comma-separated target list (SELECT expressions).
func (p *Parser) parseTargetList() (*nodes.List, error) {
	first, err := p.parseTargetEl()
	if err != nil {
		return nil, err
	}
	if first == nil {
		return nil, nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		el, err := p.parseTargetEl()
		if err != nil {
			return nil, err
		}
		if el == nil {
			break
		}
		items = append(items, el)
	}
	return &nodes.List{Items: items}, nil
}

// parseTargetEl parses a single target list element.
func (p *Parser) parseTargetEl() (nodes.Node, error) {
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
		}, nil
	}

	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, nil
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
	return rt, nil
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
func (p *Parser) parseFromList() (*nodes.List, error) {
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
