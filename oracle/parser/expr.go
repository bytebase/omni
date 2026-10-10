package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

func exprLocStart(expr nodes.ExprNode, fallback int) int {
	loc := nodes.NodeLoc(expr)
	if loc.Start >= 0 {
		return loc.Start
	}
	return fallback
}

// Precedence levels for Pratt parsing.
const (
	precNone    = 0
	precOr      = 1
	precAnd     = 2
	precNot     = 3
	precIs      = 4
	precComp    = 5
	precLike    = 6
	precConcat  = 7
	precAdd     = 8
	precMul     = 9
	precUnary   = 10
	precExpon   = 11
	precPrimary = 12
)

// parseExpr parses an expression using Pratt parsing / precedence climbing.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/About-SQL-Expressions.html
//
//	expression ::= or_expr
//	or_expr    ::= and_expr { OR and_expr }
//	and_expr   ::= not_expr { AND not_expr }
//	not_expr   ::= NOT not_expr | comparison_expr
//	...
func (p *Parser) parseExpr() (nodes.ExprNode, error) {
	return p.parseExprPrec(precOr)
}

// parseExprPrec parses expressions at or above the given precedence level.
func (p *Parser) parseExprPrec(minPrec int) (nodes.ExprNode, error) {
	left, parseErr695 := p.parsePrefix()
	if parseErr695 != nil {
		return nil, parseErr695

		// Handle MODEL cell reference subscript: expr[dim1, dim2, ...]
	}
	if left == nil {
		return nil, nil
	}
	var parseErr696 error

	left, parseErr696 = p.parseSubscriptIfPresent(left)
	if parseErr696 !=

		// Oracle legacy outer join: column_ref(+)
		nil {
		return nil, parseErr696
	}
	var parseErr697 error

	left, parseErr697 = p.parseOuterJoinIfPresent(left)
	if parseErr697 != nil {
		return nil, parseErr697
	}

	for {
		if p.startsDotPostfix() {
			var err error
			left, err = p.parseDotPostfix(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		if p.startsCursorAttributePostfix() {
			var err error
			left, err = p.parseCursorAttributePostfix(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		if p.startsAtTimeZoneExpr() {
			var err error
			left, err = p.parseAtTimeZoneExpr(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		if prec, ok := p.postfixInfo(); ok && prec >= minPrec {
			var err error
			left, err = p.parsePostfix(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		prec, op, isBool := p.infixInfo()
		if prec < minPrec {
			break
		}

		if isBool {
			var parseErr698 error
			left, parseErr698 = p.parseBoolInfix(left, prec, op)
			if parseErr698 != nil {
				return nil, parseErr698
			}
		} else {
			var parseErr699 error
			left, parseErr699 = p.parseBinaryInfix(left, prec, op)
			if parseErr699 !=

				// Check for postfix operators: IS, BETWEEN, IN, LIKE, NOT BETWEEN/IN/LIKE
				nil {
				return nil, parseErr699
			}
		}
	}
	if p.cur.Type == kwMULTISET {
		var parseErr701 error
		left, parseErr701 = p.parseMultisetOp(left)
		if parseErr701 != nil {
			return nil,

				// infixInfo returns the precedence, operator string, and whether it's a boolean op
				// for the current token if it's an infix operator.
				parseErr701
		}
	}

	return left, nil
}

func (p *Parser) startsDotPostfix() bool {
	if p.cur.Type != '.' {
		return false
	}
	next := p.peekNext()
	return next.Type == tokIDENT || next.Type == tokQIDENT || next.Type >= 2000
}

func (p *Parser) parseDotPostfix(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := exprLocStart(left, p.pos())
	p.advance() // consume '.'
	memberStart := p.pos()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != '(' {
		return &nodes.FieldAccessExpr{
			Expr:  left,
			Field: name,
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: memberStart, End: p.prev.End}},
		Args:     &nodes.List{Items: []nodes.Node{left}},
		Loc:      nodes.Loc{Start: start},
	}
	p.advance() // consume '('
	if p.cur.Type != ')' {
		args, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if args != nil {
			fc.Args.Items = append(fc.Args.Items, args.Items...)
		}
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	fc.Loc.End = p.prev.End
	return p.parseFuncCallPostfix(fc)
}

func (p *Parser) startsCursorAttributePostfix() bool {
	if p.cur.Type != '%' {
		return false
	}
	next := p.peekNext()
	return p.isCursorAttributeToken(next)
}

func (p *Parser) isCursorAttributeToken(tok Token) bool {
	switch tok.Str {
	case "FOUND", "ISOPEN", "NOTFOUND", "ROWCOUNT":
		return tok.Type == tokIDENT || tok.Type == tokQIDENT || tok.Type >= 2000
	default:
		return false
	}
}

func (p *Parser) parseCursorAttributePostfix(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := exprLocStart(left, p.pos())
	p.advance() // consume '%'
	attr, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if attr == "" {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.CursorAttributeExpr{
		Cursor:    left,
		Attribute: attr,
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

func (p *Parser) startsAtTimeZoneExpr() bool {
	if !p.isIdentLikeStr("AT") {
		return false
	}
	next := p.peekNext()
	return p.isIdentLikeStrAt(next, "LOCAL") || p.isIdentLikeStrAt(next, "TIME")
}

func (p *Parser) parseAtTimeZoneExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := exprLocStart(left, p.pos())
	p.advance() // consume AT

	if p.cur.Type == kwLOCAL {
		p.advance()
		return &nodes.FuncCallExpr{
			FuncName: &nodes.ObjectName{Name: "AT_LOCAL", Loc: nodes.Loc{Start: start, End: p.prev.End}},
			Args:     &nodes.List{Items: []nodes.Node{left}},
			Loc:      nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if !p.isIdentLikeStr("TIME") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwZONE {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	zone, err := p.parseExprPrec(precPrimary)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: "AT_TIME_ZONE", Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{Items: []nodes.Node{left, zone}},
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

func (p *Parser) infixInfo() (int, string, bool) {
	switch p.cur.Type {
	case kwOR:
		return precOr, "OR", true
	case kwAND:
		return precAnd, "AND", true
	case '=':
		return precComp, "=", false
	case tokNOTEQ:
		return precComp, p.cur.Str, false
	case '<':
		return precComp, "<", false
	case '>':
		return precComp, ">", false
	case tokLESSEQ:
		return precComp, "<=", false
	case tokGREATEQ:
		return precComp, ">=", false
	case tokCONCAT:
		return precConcat, "||", false
	case '+':
		return precAdd, "+", false
	case '-':
		return precAdd, "-", false
	case '*':
		return precMul, "*", false
	case '/':
		return precMul, "/", false
	case tokEXPON:
		return precExpon, "**", false
	}
	return precNone, "", false
}

func (p *Parser) postfixInfo() (int, bool) {
	switch p.cur.Type {
	case kwIS:
		return precIs, true
	case kwBETWEEN, kwIN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4:
		return precLike, true
	case kwNOT:
		switch p.peekNext().Type {
		case kwBETWEEN, kwIN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4:
			return precLike, true
		}
	}
	return precNone, false
}

// parseBinaryInfix parses a binary infix expression.
func (p *Parser) parseBinaryInfix(left nodes.ExprNode, prec int, op string) (nodes.ExprNode, error) {
	locStart := p.pos()
	p.advance() // consume operator

	var right nodes.ExprNode
	if p.cur.Type == kwANY || p.cur.Type == kwALL {
		var parseErr702 error
		right, parseErr702 = p.parseQuantifiedComparisonOperand()
		if parseErr702 != nil {
			return nil, parseErr702
		}
		return &nodes.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: locStart, End: p.prev.End},
		}, nil
	}

	// Right-associative for exponentiation
	nextPrec := prec + 1
	if op == "**" {
		nextPrec = prec
	}

	right, parseErr702 := p.parseExprPrec(nextPrec)
	if parseErr702 != nil {
		return nil, parseErr702
	}
	if right == nil {

		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.BinaryExpr{
		Op:    op,
		Left:  left,
		Right: right,
		Loc:   nodes.Loc{Start: locStart, End: p.prev.End},
	}, nil
}

func (p *Parser) parseQuantifiedComparisonOperand() (nodes.ExprNode, error) {
	start := p.pos()
	quantifier := p.advance()
	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: quantifier.Str, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		subquery, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		fc.Args.Items = append(fc.Args.Items, &nodes.SubqueryExpr{
			Subquery: subquery,
			Loc:      nodes.Loc{Start: start, End: p.prev.End},
		})
	} else {
		args, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if args != nil {
			fc.Args.Items = append(fc.Args.Items, args.Items...)
		}
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseBoolInfix parses AND/OR boolean expressions.
func (p *Parser) parseBoolInfix(left nodes.ExprNode, prec int, op string) (nodes.ExprNode, error) {
	start := p.pos()
	var boolop nodes.BoolExprType
	if op == "AND" {
		boolop = nodes.BOOL_AND
	} else {
		boolop = nodes.BOOL_OR
	}

	p.advance() // consume AND/OR

	right, parseErr703 := p.parseExprPrec(prec + 1)
	if parseErr703 != nil {
		return nil, parseErr703
	}
	if right == nil {

		return nil, p.syntaxErrorAtCur(

		// Flatten: if left is the same bool type, merge args
		)
	}

	if be, ok := left.(*nodes.BoolExpr); ok && be.Boolop == boolop {
		be.Args.Items = append(be.Args.Items, right)
		return be, nil
	}

	return &nodes.BoolExpr{
		Boolop: boolop,
		Args:   &nodes.List{Items: []nodes.Node{left, right}},
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parsePrefix parses a prefix expression (unary operators and primary expressions).
func (p *Parser) parsePrefix() (nodes.ExprNode, error) {
	start := p.pos()

	switch p.cur.Type {
	case kwNOT:
		p.advance()
		operand, parseErr704 := p.parseExprPrec(precNot)
		if parseErr704 != nil {
			return nil, parseErr704
		}
		if operand == nil {

			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.BoolExpr{
			Boolop: nodes.BOOL_NOT,
			Args:   &nodes.List{Items: []nodes.Node{operand}},
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case '-':
		p.advance()
		operand, parseErr705 := p.parseExprPrec(precUnary)
		if parseErr705 != nil {
			return nil, parseErr705
		}
		if operand == nil {

			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.UnaryExpr{
			Op:      "-",
			Operand: operand,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case '+':
		p.advance()
		operand, parseErr706 := p.parseExprPrec(precUnary)
		if parseErr706 != nil {
			return nil, parseErr706
		}
		if operand == nil {

			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.UnaryExpr{
			Op:      "+",
			Operand: operand,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwPRIOR:
		p.advance()
		operand, parseErr707 := p.parseExprPrec(precUnary)
		if parseErr707 != nil {
			return nil, parseErr707
		}
		if operand == nil {

			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.UnaryExpr{
			Op:      "PRIOR",
			Operand: operand,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwCONNECT_BY_ROOT:
		p.advance()
		operand, parseErr708 := p.parseExprPrec(precUnary)
		if parseErr708 != nil {
			return nil, parseErr708
		}
		if operand == nil {

			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.UnaryExpr{
			Op:      "CONNECT_BY_ROOT",
			Operand: operand,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	default:
		return p.parsePrimary()
	}
}

// parsePrimary parses a primary expression (literals, identifiers, function calls, etc).
func (p *Parser) parsePrimary() (nodes.ExprNode, error) {
	start := p.pos()

	// A keyword Oracle does not reserve, followed by the legacy outer-join
	// marker (+), is a column, whatever construct the word would otherwise
	// open: WHERE cast(+) = u.a. A few words Oracle still reads as their
	// construct, and it rejects the + inside it. Pseudo-columns stay
	// pseudo-columns: Oracle reads SYSTIMESTAMP(+) as SYSTIMESTAMP(precision)
	// and rejects it (ORA-30088) even when the table has such a column.
	if p.cur.Type >= 2000 && !isOracleSQLReservedKeyword(p.cur) &&
		!isOraclePseudoColumnKeyword(p.cur.Type) && p.nextIsOuterJoinMarker() {
		if keepsConstructBeforeOuterJoin(p.cur.Type) {
			p.advance() // consume the keyword
			p.advance() // consume (
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseIdentExpr()
	}

	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		return &nodes.NumberLiteral{
			Val:  tok.Str,
			Ival: tok.Ival,
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case tokFCONST:
		tok := p.advance()
		return &nodes.NumberLiteral{
			Val:     tok.Str,
			IsFloat: true,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case tokSCONST:
		tok := p.advance()
		return &nodes.StringLiteral{
			Val: tok.Str,
			Loc: nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case tokNCHARLIT:
		tok := p.advance()
		return &nodes.StringLiteral{
			Val:     tok.Str,
			IsNChar: true,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwNULL:
		p.advance()
		return &nodes.NullLiteral{
			Loc: nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case tokBIND:
		return p.parseBindVariable()

	case '*':
		p.advance()
		return &nodes.Star{
			Loc: nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case '(':
		return p.parseParenExpr()

	case kwCASE:
		return p.parseCaseOrColumn()

	// CAST, DECODE, and INTERVAL are not reserved; without the '(' or the
	// string literal that opens the function or literal, they name a column.
	case kwCAST:
		if p.peekNext().Type != '(' {
			return p.parseIdentExpr()
		}
		return p.parseCastExpr()

	case kwDATE, kwTIMESTAMP:
		if p.peekNext().Type == tokSCONST {
			return p.parseDateTimeLiteral()
		}
		return p.parseIdentExpr()

	case kwDECODE:
		if p.peekNext().Type != '(' {
			return p.parseIdentExpr()
		}
		return p.parseDecodeExpr()

	case kwEXISTS:
		return p.parseExistsExpr()

	case kwCURSOR:
		return p.parseCursorExpr()

	case kwTREAT:
		return p.parseTreatExpr()

	case kwXMLELEMENT:
		return p.parseXmlElement()
	case kwXMLFOREST:
		return p.parseXmlGenericFunc("XMLFOREST")
	case kwXMLAGG:
		return p.parseXmlAgg()
	case kwXMLROOT:
		return p.parseXmlRoot()
	case kwXMLPARSE:
		return p.parseXmlContentFunc("XMLPARSE")
	case kwXMLSERIALIZE:
		return p.parseXmlSerialize()

	case kwJSON_OBJECT:
		return p.parseJsonObjectOrArray("JSON_OBJECT")
	case kwJSON_ARRAY:
		return p.parseJsonObjectOrArray("JSON_ARRAY")
	case kwJSON_VALUE:
		return p.parseJsonPathFunc("JSON_VALUE")
	case kwJSON_QUERY:
		return p.parseJsonPathFunc("JSON_QUERY")
	case kwJSON_EXISTS:
		return p.parseJsonPathFunc("JSON_EXISTS")
	case kwJSON_MERGEPATCH:
		return p.parseJsonPathFunc("JSON_MERGEPATCH")

	case kwINTERVAL:
		if p.peekNext().Type != tokSCONST {
			return p.parseIdentExpr()
		}
		return p.parseIntervalExpr()

	default:
		// Pseudo columns
		if p.isPseudoColumn() {
			return p.parsePseudoColumn()
		}

		// Identifier — could be column ref, function call, or keyword-as-identifier.
		// A clause keyword Oracle reserves (FROM, WHERE, ...) ends the operand.
		// The ones it does not reserve (FETCH, JOIN, MODEL, OFFSET, USING) name
		// columns in practice and start their clauses only after a complete
		// expression or table reference, never where an operand begins.
		// TestOracleNonReservedKeywordsAsColumns checks them on the engine.
		if isOracleClauseStarterKeyword(p.cur.Type) && isOracleSQLReservedKeyword(p.cur) {
			return nil, nil
		}
		if p.isIdentLike() {
			return p.parseIdentExpr()
		}

		return nil, nil
	}
}

// parseCaseOrColumn parses CASE as Oracle reads it, since Oracle does not
// reserve the word: a CASE expression when WHEN follows it (searched CASE) or
// follows the selector expression after it (simple CASE, whose selector may
// be any expression, Boolean ones included in 23ai), and a column named CASE
// otherwise (ORDER BY case NULLS FIRST, WHERE case LIKEC 'x%', case + 1).
// decideCaseTokens settles every CASE of the range in one linear pass, so no
// expression is parsed twice to decide.
func (p *Parser) parseCaseOrColumn() (nodes.ExprNode, error) {
	if p.caseKinds == nil {
		p.caseKinds = decideCaseTokens(p.source, p.cur.Loc, p.lexer.end)
	}
	opensExpression, decided := p.caseKinds[p.cur.Loc]
	if decided && !opensExpression {
		return p.parseIdentExpr()
	}
	return p.parseCaseExpr()
}

// decideCaseTokens lexes source[start:end] and reports, for each CASE token,
// whether it opens a CASE expression. A CASE opens one when the first token
// after it that can end an operand at its own depth is WHEN: the scan skips
// parenthesized groups and the CASE expressions that start after it, and
// stops at WHEN, THEN, ELSE, END, ')', ',', ';', and the reserved words that
// begin a clause. Deciding right to left with a jump table (stop) visits each
// token a constant number of times.
func decideCaseTokens(source string, start, end int) map[int]bool {
	lexer := NewLexerRange(source, start, end)
	var toks []Token
	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF || lexer.Err != nil {
			break
		}
		toks = append(toks, tok)
	}
	n := len(toks)

	// closing[i] is the index of the ')' or ']' that closes the '(' or '['
	// at i (a MODEL cell reference such as s['Mouse Pad', 1998]), or n.
	closing := make([]int, n)
	var open []int
	for i, tok := range toks {
		switch tok.Type {
		case '(', '[':
			closing[i] = n
			open = append(open, i)
		case ')', ']':
			want := '('
			if tok.Type == ']' {
				want = '['
			}
			if len(open) > 0 && toks[open[len(open)-1]].Type == int(want) {
				closing[open[len(open)-1]] = i
				open = open[:len(open)-1]
			}
		}
	}

	// stop[i] is the index of the first token from i on, at i's depth, that
	// ends an operand; n stands for the end of the range.
	stop := make([]int, n+1)
	stop[n] = n
	after := func(i int) int {
		if i+1 >= n {
			return n
		}
		return stop[i+1]
	}
	kinds := make(map[int]bool)
	for i := n - 1; i >= 0; i-- {
		tok := toks[i]
		// A word after '.' is a name component (t.case, t.end), and END where
		// an operand must start is a column named END (1 + end, THEN end):
		// neither opens nor closes anything for the scan.
		plain := i > 0 && (toks[i-1].Type == '.' ||
			tok.Type == kwEND && expectsOperandAfter(toks[i-1]))
		switch {
		case plain:
			stop[i] = after(i)
		case isCaseScanStop(tok):
			stop[i] = i
		case tok.Type == '(' || tok.Type == '[':
			stop[i] = after(closing[i])
		case tok.Type == kwCASE:
			// s is where the scan for WHEN stops; -1 when the token after
			// CASE cannot start a selector, so CASE is a column (case IS NULL).
			s := -1
			switch {
			case i+1 >= n:
			case toks[i+1].Type == kwWHEN:
				s = i + 1
			case !canStartCaseSelector(toks[i+1]):
			default:
				s = after(i)
			}
			if s >= 0 && s < n && toks[s].Type == kwWHEN {
				kinds[tok.Loc] = true
				if e := caseExprEnd(toks, stop, s); e >= 0 {
					stop[i] = after(e)
				} else {
					stop[i] = s
				}
			} else {
				kinds[tok.Loc] = false
				stop[i] = after(i)
			}
		default:
			stop[i] = after(i)
		}
	}
	return kinds
}

// caseExprEnd follows WHEN ... THEN ... [WHEN ... THEN ...]... [ELSE ...] END
// from the WHEN at w through the stop table and returns the index of END, or
// -1 when the CASE expression is malformed.
func caseExprEnd(toks []Token, stop []int, w int) int {
	n := len(toks)
	next := func(i int) int {
		if i+1 >= n {
			return n
		}
		return stop[i+1]
	}
	for {
		then := next(w)
		if then >= n || toks[then].Type != kwTHEN {
			return -1
		}
		r := next(then)
		if r >= n {
			return -1
		}
		switch toks[r].Type {
		case kwWHEN:
			w = r
		case kwELSE:
			if e := next(r); e < n && toks[e].Type == kwEND {
				return e
			}
			return -1
		case kwEND:
			return r
		default:
			return -1
		}
	}
}

// canStartCaseSelector reports whether tok can begin the selector of a simple
// CASE. Operators and keywords that only follow an operand (=, IS, IN, AND,
// ...) and the clause words that end one cannot; END can, as a column named
// END, which Oracle accepts there.
func canStartCaseSelector(tok Token) bool {
	switch tok.Type {
	case '=', '<', '>', '*', '/', '.', tokLESSEQ, tokGREATEQ, tokNOTEQ,
		tokCONCAT, tokEXPON, kwIS, kwIN, kwLIKE, kwBETWEEN, kwAND, kwOR,
		kwASC, kwDESC:
		return false
	case kwEND:
		return true
	}
	return !isCaseScanStop(tok)
}

// expectsOperandAfter reports whether an operand must follow tok, so an END
// right after it is a column named END rather than the end of a CASE. Only
// tokens that work at expression level matter: an END inside parentheses is
// skipped with its group. The unary PRIOR and CONNECT_BY_ROOT, ESCAPE, OF
// (MEMBER OF, SUBMULTISET OF), and ZONE (AT TIME ZONE) all take an operand.
func expectsOperandAfter(tok Token) bool {
	switch tok.Type {
	case '+', '-', '*', '/', '=', '<', '>', '(', ',', '.',
		tokCONCAT, tokEXPON, tokLESSEQ, tokGREATEQ, tokNOTEQ,
		kwCASE, kwWHEN, kwTHEN, kwELSE, kwAND, kwOR, kwNOT,
		kwPRIOR, kwCONNECT_BY_ROOT, kwESCAPE, kwOF, kwZONE,
		kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4, kwBETWEEN,
		kwSELECT, kwWHERE, kwBY, kwHAVING, kwON, kwSET, kwRETURN:
		return true
	}
	return false
}

// isCaseScanStop reports the tokens that end an operand for decideCaseTokens.
func isCaseScanStop(tok Token) bool {
	switch tok.Type {
	case kwWHEN, kwTHEN, kwELSE, kwEND, ')', ']', ',', ';',
		kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT,
		kwMINUS, kwINTO, kwAS, kwCONNECT, kwSTART, kwSET, kwVALUES, kwFOR:
		return true
	}
	return false
}

// keepsConstructBeforeOuterJoin reports the non-reserved keywords Oracle
// parses as their own construct even when (+) follows, rejecting word(+)
// (ORA-00931, ORA-00936), so they cannot name an outer-joined column
// unquoted. Engine-checked by TestOracleNonReservedKeywordsAsColumns.
func keepsConstructBeforeOuterJoin(tokenType int) bool {
	switch tokenType {
	case kwJSON, kwJSON_ARRAY, kwJSON_EXISTS, kwJSON_MERGEPATCH, kwJSON_OBJECT,
		kwJSON_QUERY, kwJSON_TABLE, kwJSON_VALUE, kwTREAT,
		kwXMLELEMENT, kwXMLFOREST, kwXMLROOT:
		return true
	}
	return false
}

// nextIsOuterJoinMarker reports whether the tokens after the current one are
// the legacy outer-join marker (+).
func (p *Parser) nextIsOuterJoinMarker() bool {
	toks := p.peekAhead(3)
	return toks[0].Type == '(' && toks[1].Type == '+' && toks[2].Type == ')'
}

// parseDateTimeLiteral parses ANSI datetime literals such as DATE '2020-01-01'.
func (p *Parser) parseDateTimeLiteral() (nodes.ExprNode, error) {
	start := p.pos()
	typeTok := p.advance()
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	valueTok := p.advance()
	return &nodes.DateTimeLiteral{
		TypeName: typeTok.Str,
		Val:      valueTok.Str,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// isIntervalField returns true if the current token is an interval date field keyword.
func (p *Parser) isIntervalField() bool {
	if !p.isIdentLike() {
		return false
	}
	switch p.cur.Str {
	case "YEAR", "MONTH", "DAY", "HOUR", "MINUTE", "SECOND":
		return true
	}
	return false
}

// parseIntervalExpr parses an INTERVAL literal expression.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Literals.html
//
//	INTERVAL 'value' { YEAR | MONTH | DAY | HOUR | MINUTE | SECOND }
//	  [ ( precision ) ] [ TO { YEAR | MONTH | DAY | HOUR | MINUTE | SECOND } [ ( precision ) ] ]
func (p *Parser) parseIntervalExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume INTERVAL

	// Parse the value expression (typically a string literal like '1')
	value, parseErr709 := p.parseExprPrec(precConcat)
	if parseErr709 != nil {
		return nil,

			// Parse the FROM field
			parseErr709
	}

	var from, to string

	if p.isIntervalField() {
		from = p.cur.Str
		p.advance()
		// Optional precision: ( n )
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// Optional TO field
	if p.cur.Type == kwTO {
		p.advance()
		if p.isIntervalField() {
			to = p.cur.Str
			p.advance()
			// Optional precision: ( n )
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == tokICONST {
					p.advance()
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	return &nodes.IntervalExpr{
		Value: value,
		From:  from,
		To:    to,
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseIdentExpr parses an identifier-starting expression.
// It could be a function call (name(...)), a column ref (name or table.column), etc.
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
	start := p.pos()
	name1, parseErr710 := p.parseIdentifier()
	if parseErr710 != nil {
		return nil, parseErr710

		// Check for function call: name(
	}
	if name1 == "" {
		return nil, nil
	}

	if p.cur.Type == '(' {
		if p.peekNext().Type == '+' {
			return p.parseOuterJoinIfPresent(&nodes.ColumnRef{
				Column: name1,
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
			})
		}
		if name1 == "EXTRACT" {
			return p.parseExtractExpr(start)
		}
		if name1 == "XMLQUERY" {
			return p.parseXMLQueryExpr(start)
		}
		if name1 == "XMLCAST" {
			return p.parseXMLCastExpr(start)
		}
		// But first check if this looks like a schema-qualified function: name.name(
		// No — just name( is the function call case
		return p.parseFuncCall(name1, "", start)
	}

	// Check for schema.name or table.column
	if p.cur.Type == '.' {
		p.advance() // consume '.'

		// table.*
		if p.cur.Type == '*' {
			p.advance()
			return &nodes.ColumnRef{
				Table:  name1,
				Column: "*",
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		name2, parseErr711 := p.parseIdentifier()
		if parseErr711 != nil {
			return nil, parseErr711
		}
		if name2 == "" {
			return &nodes.ColumnRef{
				Column: name1,
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// schema.func( ?
		if p.cur.Type == '(' {
			if p.peekNext().Type == '+' {
				return p.parseOuterJoinIfPresent(&nodes.ColumnRef{
					Table:  name1,
					Column: name2,
					Loc:    nodes.Loc{Start: start, End: p.prev.End},
				})
			}
			return p.parseFuncCall(name2, name1, start)
		}

		// schema.table.column or schema.table.*
		if p.cur.Type == '.' {
			p.advance()
			if p.cur.Type == '*' {
				p.advance()
				return &nodes.ColumnRef{
					Schema: name1,
					Table:  name2,
					Column: "*",
					Loc:    nodes.Loc{Start: start, End: p.prev.End},
				}, nil
			}
			name3, parseErr712 := p.parseIdentifier()
			if parseErr712 != nil {
				return nil, parseErr712
			}
			if p.cur.Type == '(' && p.peekNext().Type == '+' {
				return p.parseOuterJoinIfPresent(&nodes.ColumnRef{
					Schema: name1,
					Table:  name2,
					Column: name3,
					Loc:    nodes.Loc{Start: start, End: p.prev.End},
				})
			}
			if p.cur.Type == '(' {
				return p.parseFuncCall(name3, name1+"."+name2, start)
			}
			return &nodes.ColumnRef{
				Schema: name1,
				Table:  name2,
				Column: name3,
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// table.column
		return &nodes.ColumnRef{
			Table:  name1,
			Column: name2,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// Simple column reference
	return &nodes.ColumnRef{
		Column: name1,
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseExtractExpr parses Oracle EXTRACT(datetime_field FROM expr).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/EXTRACT-datetime.html
//
//	EXTRACT({ YEAR | MONTH | DAY | HOUR | MINUTE | SECOND } FROM expr)
func (p *Parser) parseExtractExpr(start int) (nodes.ExprNode, error) {
	p.advance() // consume '('

	if !(p.isIntervalField() && p.peekNext().Type == kwFROM) {
		fc := &nodes.FuncCallExpr{
			FuncName: &nodes.ObjectName{Name: "EXTRACT", Loc: nodes.Loc{Start: start, End: p.prev.End}},
			Args:     &nodes.List{},
			Loc:      nodes.Loc{Start: start},
		}
		if p.cur.Type != ')' {
			for {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if arg == nil {
					return nil, p.syntaxErrorAtCur()
				}
				fc.Args.Items = append(fc.Args.Items, arg)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		fc.Loc.End = p.prev.End
		return p.parseFuncCallPostfix(fc)
	}

	field, parseErr713 := p.parseIdentifier()
	if parseErr713 != nil {
		return nil, parseErr713
	}
	if field == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != kwFROM {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	expr, parseErr714 := p.parseExpr()
	if parseErr714 != nil {
		return nil, parseErr714
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	return &nodes.ExtractExpr{
		Field: field,
		Expr:  expr,
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

func (p *Parser) parseXMLQueryExpr(start int) (nodes.ExprNode, error) {
	return p.parseXMLSpecialFunc("XMLQUERY", start)
}

func (p *Parser) parseXMLCastExpr(start int) (nodes.ExprNode, error) {
	return p.parseXMLSpecialFunc("XMLCAST", start)
}

func (p *Parser) parseXMLSpecialFunc(name string, start int) (nodes.ExprNode, error) {
	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}
	p.advance() // consume '('
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			if depth == 0 {
				break
			}
			depth--
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	fc.Loc.End = p.prev.End
	return p.parseFuncCallPostfix(fc)
}

// parseSubscriptIfPresent checks if the current token is '[' and parses a
// MODEL cell reference subscript: expr[dim1, dim2, ...].
// Returns the original expression if no '[' is found.
func (p *Parser) parseSubscriptIfPresent(expr nodes.ExprNode) (nodes.ExprNode, error) {
	if p.cur.Type != '[' {
		return expr, nil
	}

	// Determine start position and function name from the expression
	var start int
	funcNameEnd := -1
	funcName := ""
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		funcName = e.Column
		start = e.Loc.Start
		funcNameEnd = e.Loc.End
	default:
		start = p.pos()
	}
	p.advance() // consume '['

	args := &nodes.List{}
	for {
		if p.cur.Type == ']' {
			break
		}
		arg, parseErr713 := p.parseExpr()
		if parseErr713 != nil {
			return nil, parseErr713
		}
		if arg == nil {
			break
		}
		args.Items = append(args.Items, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ']' {
		p.advance()
	}

	// Represent subscript access as a FuncCallExpr (reusing the existing node type).
	// This is a MODEL cell reference: measure[dim1, dim2]
	return &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: funcName, Loc: nodes.Loc{Start: start, End: funcNameEnd}},
		Args:     args,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseOuterJoinIfPresent checks for Oracle legacy outer join syntax: column_ref(+).
// If the current tokens are '(' '+' ')', marks the ColumnRef with OuterJoin=true.
func (p *Parser) parseOuterJoinIfPresent(expr nodes.ExprNode) (nodes.ExprNode, error) {
	cr, ok := expr.(*nodes.ColumnRef)
	if !ok {
		return expr, nil
	}
	if p.cur.Type != '(' {
		return expr, nil
	}
	next := p.peekNext()
	if next.Type != '+' {
		return expr, nil
	}
	p.advance() // consume '('
	p.advance() // consume '+'
	if p.cur.Type == ')' {
		p.advance() // consume ')'
	}
	cr.OuterJoin = true
	cr.Loc.End = p.prev.End
	return cr, nil
}

// parseFuncCall parses a function call after the name has been consumed.
// The opening '(' is the current token.
func (p *Parser) parseFuncCall(name, schema string, start int) (nodes.ExprNode, error) {
	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{
			Schema: schema,
			Name:   name,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		},
		Args: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	p.advance() // consume '('

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if p.cur.Type == ')' {
			p.advance()
		}
		// Still check for KEEP/OVER after COUNT(*)
		return p.parseFuncCallPostfix(fc)
	}

	// DISTINCT
	if p.cur.Type == kwDISTINCT || p.cur.Type == kwUNIQUE {
		fc.Distinct = true
		p.advance()
	}

	// ALL
	if p.cur.Type == kwALL {
		p.advance()
	}

	if schema == "" && name == "TRIM" {
		if err := p.parseTrimFuncArgs(fc); err != nil {
			return nil, err
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		return p.parseFuncCallPostfix(fc)
	}

	// Arguments
	if p.cur.Type != ')' {
		for {
			arg, parseErr714 := p.parseExpr()
			if parseErr714 != nil {
				return nil, parseErr714
			}
			if p.cur.Type == tokASSOC {
				named, err := p.parseNamedArgExpr(arg)
				if err != nil {
					return nil, err
				}
				arg = named
			}
			if schema == "" && name == "TRIM" && p.cur.Type == kwFROM {
				if arg != nil {
					fc.Args.Items = append(fc.Args.Items, arg)
				}
				p.advance()
				source, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if source == nil {
					return nil, p.syntaxErrorAtCur()
				}
				fc.Args.Items = append(fc.Args.Items, source)
				break
			}
			if arg != nil {
				fc.Args.Items = append(fc.Args.Items, arg)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

func (p *Parser) parseTrimFuncArgs(fc *nodes.FuncCallExpr) error {
	if p.cur.Type == ')' {
		return nil
	}

	arg, err := p.parseExpr()
	if err != nil {
		return err
	}
	if arg == nil {
		return p.syntaxErrorAtCur()
	}
	fc.Args.Items = append(fc.Args.Items, arg)

	if p.isTrimModeArg(arg) && p.cur.Type != kwFROM && p.cur.Type != ',' && p.cur.Type != ')' {
		charArg, err := p.parseExpr()
		if err != nil {
			return err
		}
		if charArg == nil {
			return p.syntaxErrorAtCur()
		}
		fc.Args.Items = append(fc.Args.Items, charArg)
	}

	if p.cur.Type == kwFROM {
		p.advance()
		source, err := p.parseExpr()
		if err != nil {
			return err
		}
		if source == nil {
			return p.syntaxErrorAtCur()
		}
		fc.Args.Items = append(fc.Args.Items, source)
		return nil
	}

	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseExpr()
		if err != nil {
			return err
		}
		if arg == nil {
			return p.syntaxErrorAtCur()
		}
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	return nil
}

func (p *Parser) isTrimModeArg(arg nodes.ExprNode) bool {
	ref, ok := arg.(*nodes.ColumnRef)
	if !ok || ref.Schema != "" || ref.Table != "" {
		return false
	}
	switch ref.Column {
	case "LEADING", "TRAILING", "BOTH":
		return true
	default:
		return false
	}
}

func (p *Parser) parseNamedArgExpr(arg nodes.ExprNode) (nodes.ExprNode, error) {
	ref, ok := arg.(*nodes.ColumnRef)
	if !ok || ref.Schema != "" || ref.Table != "" || ref.Column == "" {
		return nil, p.syntaxErrorAtCur()
	}
	start := ref.Loc.Start
	name := ref.Column
	p.advance() // consume =>
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.NamedArgExpr{
		Name: name,
		Expr: value,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseFuncCallPostfix checks for WITHIN GROUP, KEEP, and OVER clauses
// after a function call's closing parenthesis.
func (p *Parser) parseFuncCallPostfix(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// WITHIN GROUP (ORDER BY ...)
	if p.cur.Type == kwWITHIN {
		var parseErr715 error
		fc.OrderBy, parseErr715 = p.parseWithinGroup()
		if parseErr715 !=

			// KEEP (DENSE_RANK FIRST/LAST ORDER BY ...)
			nil {
			return nil, parseErr715
		}
	}

	if p.cur.Type == kwKEEP {
		var parseErr716 error
		fc.KeepClause, parseErr716 = p.parseKeepClause()
		if parseErr716 !=

			// OVER (analytic window specification)
			nil {
			return nil, parseErr716
		}
	}

	if p.cur.Type == kwOVER {
		var parseErr717 error
		fc.Over, parseErr717 = p.parseOverClause()
		if parseErr717 != nil {
			return nil, parseErr717
		}
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseOverClause parses an analytic function's OVER clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Analytic-Functions.html
//
//	OVER ( [ partition_by_clause ] [ order_by_clause [ windowing_clause ] ] )
//	OVER window_name
func (p *Parser) parseOverClause() (*nodes.WindowSpec, error) {
	start := p.pos()
	p.advance() // consume OVER

	ws := &nodes.WindowSpec{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		// OVER window_name
		if p.isIdentLike() {
			var parseErr718 error
			ws.WindowName, parseErr718 = p.parseIdentifier()
			if parseErr718 != nil {
				return nil, parseErr718
			}
		}
		ws.Loc.End = p.prev.End
		return ws, nil
	}

	p.advance() // consume '('

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		ws.PartitionBy = &nodes.List{}
		for {
			expr, parseErr719 := p.parseExpr()
			if parseErr719 != nil {
				return nil, parseErr719
			}
			if expr != nil {
				ws.PartitionBy.Items = append(ws.PartitionBy.Items, expr)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		var parseErr720 error
		ws.OrderBy, parseErr720 = p.parseOrderByList()
		if parseErr720 !=

			// Windowing clause: ROWS | RANGE | GROUPS
			nil {
			return nil, parseErr720
		}
	}

	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		var parseErr721 error
		ws.Frame, parseErr721 = p.parseWindowFrame()
		if parseErr721 != nil {
			return nil, parseErr721
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	ws.Loc.End = p.prev.End
	return ws, nil
}

// parseWindowFrame parses a window frame specification.
//
//	{ ROWS | RANGE | GROUPS }
//	  { BETWEEN bound AND bound | bound }
func (p *Parser) parseWindowFrame() (*nodes.WindowFrame, error) {
	start := p.pos()
	wf := &nodes.WindowFrame{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwROWS:
		wf.Type = nodes.WINDOW_ROWS
	case kwRANGE:
		wf.Type = nodes.WINDOW_RANGE
	case kwGROUPS:
		wf.Type = nodes.WINDOW_GROUPS
	}
	p.advance() // consume ROWS/RANGE/GROUPS

	if p.cur.Type == kwBETWEEN {
		p.advance()
		var // consume BETWEEN
		parseErr722 error
		wf.Start, parseErr722 = p.parseWindowBound()
		if parseErr722 != nil {
			return nil, parseErr722
		}
		if p.cur.Type == kwAND {
			p.advance() // consume AND
		}
		var parseErr723 error
		wf.End, parseErr723 = p.parseWindowBound()
		if parseErr723 !=

			// Single bound (start only)
			nil {
			return nil, parseErr723
		}
	} else {
		var parseErr724 error

		wf.Start, parseErr724 = p.parseWindowBound()
		if parseErr724 != nil {
			return nil, parseErr724
		}
	}

	wf.Loc.End = p.prev.End
	return wf, nil
}

// parseWindowBound parses a window frame bound.
//
//	UNBOUNDED PRECEDING | UNBOUNDED FOLLOWING
//	CURRENT ROW
//	expr PRECEDING | expr FOLLOWING
func (p *Parser) parseWindowBound() (*nodes.WindowBound, error) {
	start := p.pos()
	wb := &nodes.WindowBound{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwUNBOUNDED {
		p.advance() // consume UNBOUNDED
		if p.cur.Type == kwPRECEDING {
			wb.Type = nodes.WINDOW_UNBOUNDED_PRECEDING
			p.advance()
		} else if p.cur.Type == kwFOLLOWING {
			wb.Type = nodes.WINDOW_UNBOUNDED_FOLLOWING
			p.advance()
		}
	} else if p.cur.Type == kwCURRENT {
		p.advance() // consume CURRENT
		wb.Type = nodes.WINDOW_CURRENT_ROW
		if p.cur.Type == kwROW {
			p.advance() // consume ROW
		}
	} else {
		var parseErr725 error
		// expr PRECEDING | expr FOLLOWING
		wb.Value, parseErr725 = p.parseExprPrec(precAdd)
		if parseErr725 != nil {
			return nil, parseErr725
		}
		if p.cur.Type == kwPRECEDING {
			wb.Type = nodes.WINDOW_VALUE_PRECEDING
			p.advance()
		} else if p.cur.Type == kwFOLLOWING {
			wb.Type = nodes.WINDOW_VALUE_FOLLOWING
			p.advance()
		}
	}

	wb.Loc.End = p.prev.End
	return wb, nil
}

// parseKeepClause parses a KEEP (DENSE_RANK FIRST/LAST ORDER BY ...) clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/FIRST.html
//
//	KEEP ( DENSE_RANK { FIRST | LAST } ORDER BY sort_list )
func (p *Parser) parseKeepClause() (*nodes.KeepClause, error) {
	start := p.pos()
	p.advance() // consume KEEP

	kc := &nodes.KeepClause{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == '(' {
		p.advance() // consume '('
	}

	// DENSE_RANK
	if p.cur.Type == kwDENSE_RANK {
		p.advance()
	}

	// FIRST or LAST
	if p.cur.Type == kwFIRST {
		kc.IsFirst = true
		p.advance()
	} else if p.cur.Type == kwLAST {
		kc.IsFirst = false
		p.advance()
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		var parseErr726 error
		kc.OrderBy, parseErr726 = p.parseOrderByList()
		if parseErr726 != nil {
			return nil, parseErr726
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	kc.Loc.End = p.prev.End
	return kc, nil
}

// parseWithinGroup parses a WITHIN GROUP (ORDER BY ...) clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/LISTAGG.html
//
//	WITHIN GROUP ( ORDER BY sort_list )
func (p *Parser) parseWithinGroup() (*nodes.List, error) {
	p.advance() // consume WITHIN

	// GROUP
	if p.cur.Type == kwGROUP {
		p.advance()
	}

	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance() // consume '('

	var orderBy *nodes.List
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		var parseErr727 error
		orderBy, parseErr727 = p.parseOrderByList()
		if parseErr727 != nil {
			return nil, parseErr727
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return orderBy, nil
}

// parseParenExpr parses a parenthesized expression or subquery.
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume '('

	// Detect scalar subquery: (SELECT ...)
	if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
		sub, parseErr728 := p.parseSelectStmt()
		if parseErr728 != nil {
			return nil, parseErr728
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		return &nodes.SubqueryExpr{
			Subquery: sub,
			Loc:      nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	inner, parseErr729 := p.parseExpr()
	if parseErr729 != nil {
		return nil, parseErr729
	}
	if inner == nil {
		if p.cur.Type == ')' {
			p.advance()
		}
		return &nodes.ParenExpr{
			Loc: nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == ',' {
		items := &nodes.List{Items: []nodes.Node{inner}}
		for p.cur.Type == ',' {
			p.advance()
			item, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if item == nil {
				return nil, p.syntaxErrorAtCur()
			}
			items.Items = append(items.Items, item)
		}
		if p.cur.Type == ')' {
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.FuncCallExpr{
			FuncName: &nodes.ObjectName{Name: "", Loc: nodes.Loc{Start: start, End: p.prev.End}},
			Args:     items,
			Loc:      nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return &nodes.ParenExpr{
		Expr: inner,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseCaseExpr parses a CASE expression (simple or searched).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CASE-Expressions.html
//
//	CASE [ expr ]
//	    WHEN condition THEN result
//	    [ WHEN condition THEN result ... ]
//	    [ ELSE default ]
//	END
func (p *Parser) parseCaseExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CASE

	// Simple CASE: CASE expr WHEN ...
	// Searched CASE: CASE WHEN ...
	var arg nodes.ExprNode
	if p.cur.Type != kwWHEN {
		var parseErr730 error
		arg, parseErr730 = p.parseExpr()
		if parseErr730 != nil {
			return nil, parseErr730
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	return p.parseCaseWhens(start, arg)
}

// parseCaseWhens parses the WHEN ... [ELSE ...] END of a CASE expression that
// starts at start, after CASE and the selector arg (nil for a searched CASE)
// have been consumed.
func (p *Parser) parseCaseWhens(start int, arg nodes.ExprNode) (nodes.ExprNode, error) {
	ce := &nodes.CaseExpr{
		Arg:   arg,
		Whens: &nodes.List{},
		Loc:   nodes.Loc{Start: start},
	}

	for p.cur.Type == kwWHEN {
		whenStart := p.pos()
		p.advance() // consume WHEN
		cond, parseErr731 := p.parseExpr()
		if parseErr731 != nil {
			return nil, parseErr731
		}
		if cond == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != kwTHEN {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		result, parseErr732 := p.parseExpr()
		if parseErr732 != nil {
			return nil, parseErr732
		}
		if result == nil {
			return nil, p.syntaxErrorAtCur()
		}
		ce.Whens.Items = append(ce.Whens.Items, &nodes.CaseWhen{
			Condition: cond,
			Result:    result,
			Loc:       nodes.Loc{Start: whenStart, End: p.prev.End},
		})
	}
	if ce.Whens.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwELSE {
		p.advance()
		var parseErr733 error
		ce.Default, parseErr733 = p.parseExpr()
		if parseErr733 != nil {
			return nil, parseErr733
		}
		if ce.Default == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type != kwEND {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	ce.Loc.End = p.prev.End
	return ce, nil
}

// parseCastExpr parses a CAST expression.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CAST.html
//
//	CAST ( expr AS datatype )
func (p *Parser) parseCastExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CAST

	if p.cur.Type == '(' {
		p.advance()
	}

	arg, parseErr734 := p.parseExpr()
	if parseErr734 != nil {
		return nil, parseErr734
	}

	if p.cur.Type == kwAS {
		p.advance()
	}

	typeName, parseErr735 := p.parseTypeName()
	if parseErr735 != nil {
		return nil, parseErr735
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return &nodes.CastExpr{
		Arg:      arg,
		TypeName: typeName,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseDecodeExpr parses Oracle's DECODE function.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/DECODE.html
//
//	DECODE ( expr, search, result [, search, result ...] [, default] )
func (p *Parser) parseDecodeExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume DECODE

	if p.cur.Type == '(' {
		p.advance()
	}

	de := &nodes.DecodeExpr{
		Pairs: &nodes.List{},
		Loc:   nodes.Loc{Start: start},
	}
	var parseErr736 error

	// First arg is the expression to decode
	de.Arg, parseErr736 = p.parseExpr()
	if parseErr736 !=

		// Parse search, result pairs
		nil {
		return nil, parseErr736
	}

	for p.cur.Type == ',' {
		p.advance() // consume ','

		search, parseErr737 := p.parseExpr()
		if parseErr737 != nil {
			return nil, parseErr737

			// This is the default value (odd argument at end)
		}

		if p.cur.Type != ',' {

			de.Default = search
			break
		}
		p.advance() // consume ','

		result, parseErr738 := p.parseExpr()
		if parseErr738 != nil {
			return nil, parseErr738
		}

		de.Pairs.Items = append(de.Pairs.Items, &nodes.DecodePair{
			Search: search,
			Result: result,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		})
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	de.Loc.End = p.prev.End
	return de, nil
}

// parseExistsExpr parses an EXISTS subquery expression.
//
//	EXISTS ( subquery )
func (p *Parser) parseExistsExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume EXISTS

	expr := &nodes.ExistsExpr{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	// Parse the inner SELECT statement.
	if p.cur.Type != kwSELECT && p.cur.Type != kwWITH {
		return nil, p.syntaxErrorAtCur()
	}
	var parseErr739 error
	expr.Subquery, parseErr739 = p.parseSelectStmt()
	if parseErr739 != nil {
		return nil, parseErr739
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	expr.Loc.End = p.prev.End
	return expr, nil
}

// parsePostfix parses postfix expression operators: IS, BETWEEN, IN, LIKE, NOT BETWEEN/IN/LIKE.
func (p *Parser) parsePostfix(left nodes.ExprNode) (nodes.ExprNode, error) {
	if left == nil {
		return nil, nil
	}

	switch p.cur.Type {
	case kwIS:
		return p.parseIsExpr(left)

	case kwBETWEEN:
		return p.parseBetweenExpr(left, false)

	case kwIN:
		return p.parseInExpr(left, false)

	case kwLIKE:
		return p.parseLikeExpr(left, false, nodes.LIKE_LIKE)
	case kwLIKEC:
		return p.parseLikeExpr(left, false, nodes.LIKE_LIKEC)
	case kwLIKE2:
		return p.parseLikeExpr(left, false, nodes.LIKE_LIKE2)
	case kwLIKE4:
		return p.parseLikeExpr(left, false, nodes.LIKE_LIKE4)

	case kwNOT:
		// NOT BETWEEN, NOT IN, NOT LIKE
		next := p.peekNext()
		switch next.Type {
		case kwBETWEEN:
			p.advance() // consume NOT
			return p.parseBetweenExpr(left, true)
		case kwIN:
			p.advance() // consume NOT
			return p.parseInExpr(left, true)
		case kwLIKE:
			p.advance() // consume NOT
			return p.parseLikeExpr(left, true, nodes.LIKE_LIKE)
		case kwLIKEC:
			p.advance() // consume NOT
			return p.parseLikeExpr(left, true, nodes.LIKE_LIKEC)
		case kwLIKE2:
			p.advance() // consume NOT
			return p.parseLikeExpr(left, true, nodes.LIKE_LIKE2)
		case kwLIKE4:
			p.advance() // consume NOT
			return p.parseLikeExpr(left, true, nodes.LIKE_LIKE4)
		}
	}

	return left, nil
}

// parseIsExpr parses IS [NOT] NULL / IS [NOT] NAN / etc.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Conditions.html
//
//	expr IS [NOT] { NULL | NAN | INFINITE | EMPTY | JSON | OF ... | A SET }
func (p *Parser) parseIsExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume IS

	not := false
	if p.cur.Type == kwNOT {
		not = true
		p.advance()
	}

	test := ""
	switch p.cur.Type {
	case kwNULL:
		test = "NULL"
		p.advance()
	case kwJSON:
		test = "JSON"
		p.advance()
	default:
		if p.isIdentLike() {
			test = p.cur.Str
			p.advance()
		}
	}
	if test == "" {

		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.IsExpr{
		Expr: left,
		Test: test,
		Not:  not,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/BETWEEN-Condition.html
//
//	expr [NOT] BETWEEN low AND high
func (p *Parser) parseBetweenExpr(left nodes.ExprNode, not bool) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume BETWEEN

	// Parse low bound at higher precedence to avoid consuming AND as boolean
	low, parseErr740 := p.parseExprPrec(precConcat)
	if parseErr740 != nil {
		return nil, parseErr740
	}
	if low == nil {

		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwAND {
		p.advance()
	} else {

		return nil, p.syntaxErrorAtCur()
	}

	high, parseErr741 := p.parseExprPrec(precConcat)
	if parseErr741 != nil {
		return nil, parseErr741
	}
	if high == nil {

		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.BetweenExpr{
		Expr: left,
		Low:  low,
		High: high,
		Not:  not,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseInExpr parses [NOT] IN ( list | subquery ).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/IN-Condition.html
//
//	expr [NOT] IN ( expr_list | subquery )
func (p *Parser) parseInExpr(left nodes.ExprNode, not bool) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume IN

	list := &nodes.List{}

	if p.cur.Type == '(' {
		p.advance()
		// IN (SELECT ...) or IN (WITH ...) — subquery
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			sub, parseErr742 := p.parseSelectStmt()
			if parseErr742 != nil {
				return nil, parseErr742
			}
			subExpr := &nodes.SubqueryExpr{
				Subquery: sub,
				Loc:      nodes.Loc{Start: start, End: p.prev.End},
			}
			list.Items = append(list.Items, subExpr)
		} else {
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				item, parseErr743 := p.parseExpr()
				if parseErr743 != nil {
					return nil, parseErr743
				}
				if item == nil {

					return nil, p.syntaxErrorAtCur()
				}
				list.Items = append(list.Items, item)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		} else {

			return nil, p.syntaxErrorAtCur()
		}
	} else {

		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.InExpr{
		Expr: left,
		List: list,
		Not:  not,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseLikeExpr parses [NOT] LIKE pattern [ESCAPE escape_char].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Pattern-matching-Conditions.html
//
//	expr [NOT] LIKE pattern [ ESCAPE escape_char ]
func (p *Parser) parseLikeExpr(left nodes.ExprNode, not bool, likeType nodes.LikeType) (nodes.ExprNode, error) {
	start := exprLocStart(left, p.pos())
	p.advance() // consume LIKE/LIKEC/LIKE2/LIKE4

	pattern, parseErr744 := p.parseExprPrec(precConcat)
	if parseErr744 != nil {
		return nil, parseErr744
	}
	if pattern == nil {

		return nil, p.syntaxErrorAtCur()
	}

	var escape nodes.ExprNode
	if p.cur.Type == kwESCAPE {
		p.advance()
		var parseErr745 error
		escape, parseErr745 = p.parseExprPrec(precConcat)
		if parseErr745 != nil {
			return nil, parseErr745
		}
		if escape == nil {

			return nil, p.syntaxErrorAtCur()
		}
	}

	return &nodes.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Escape:  escape,
		Not:     not,
		Type:    likeType,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseCursorExpr parses a CURSOR(subquery) expression.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CURSOR-Expressions.html
//
//	CURSOR ( subquery )
func (p *Parser) parseCursorExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CURSOR

	if p.cur.Type != '(' {
		return &nodes.CursorExpr{Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
	}
	p.advance() // consume '('

	subSel, parseErr746 := p.parseSelectStmt()
	if parseErr746 != nil {
		return nil, parseErr746
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return &nodes.CursorExpr{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseTreatExpr parses a TREAT(expr AS type) expression.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/TREAT.html
//
//	TREAT ( expr AS [ REF ] type )
func (p *Parser) parseTreatExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume TREAT

	if p.cur.Type != '(' {
		return &nodes.TreatExpr{Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
	}
	p.advance() // consume '('

	expr, parseErr747 := p.parseExpr()
	if parseErr747 != nil {
		return nil, parseErr747
	}

	if p.cur.Type == kwAS {
		p.advance()
	}

	// Skip optional REF
	if p.cur.Type == kwREF {
		p.advance()
	}

	typeName, parseErr748 := p.parseTypeName()
	if parseErr748 != nil {
		return nil, parseErr748
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return &nodes.TreatExpr{
		Expr:     expr,
		TypeName: typeName,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseMultisetOp parses MULTISET UNION/INTERSECT/EXCEPT operations.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/MULTISET-UNION.html
//
//	expr MULTISET { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] expr
func (p *Parser) parseMultisetOp(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume MULTISET

	op := ""
	switch p.cur.Type {
	case kwUNION:
		op = "UNION"
		p.advance()
	case kwINTERSECT:
		op = "INTERSECT"
		p.advance()
	case kwEXCEPT:
		op = "EXCEPT"
		p.advance()
	default:
		return left, nil
	}

	all := false
	if p.cur.Type == kwALL {
		all = true
		p.advance()
	} else if p.cur.Type == kwDISTINCT {
		p.advance()
	}

	right, parseErr749 := p.parseExprPrec(precComp)
	if parseErr749 != nil {
		return nil, parseErr749
	}

	return &nodes.MultisetExpr{
		Op:    op,
		Left:  left,
		Right: right,
		All:   all,
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseXmlElement parses XMLELEMENT(NAME tag, expr, ...).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/XMLELEMENT.html
//
//	XMLELEMENT ( [ NAME ] identifier_or_string, expr [, ...] )
func (p *Parser) parseXmlElement() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume XMLELEMENT

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: "XMLELEMENT", Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// Optional NAME keyword
	if p.cur.Type == kwNAME {
		p.advance()
	}

	// Element name — identifier or quoted identifier
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		arg, parseErr750 := p.parseExpr()
		if parseErr750 != nil {
			return nil, parseErr750
		}
		if arg != nil {
			fc.Args.Items = append(fc.Args.Items, arg)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseXmlGenericFunc parses XML functions with standard argument syntax.
// Used for XMLFOREST, XMLROOT.
func (p *Parser) parseXmlGenericFunc(name string) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume keyword

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Handle keyword-value pairs like VERSION '1.0'
		arg, parseErr751 := p.parseExpr()
		if parseErr751 != nil {
			return nil, parseErr751
		}
		if arg != nil {
			fc.Args.Items = append(fc.Args.Items, arg)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseXmlRoot parses XMLROOT(xml_expr, VERSION version_string | NO VALUE).
func (p *Parser) parseXmlRoot() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume XMLROOT

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: "XMLROOT", Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// XML expression
	arg, parseErr752 := p.parseExpr()
	if parseErr752 != nil {
		return nil, parseErr752
	}
	if arg != nil {
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	if p.cur.Type == ',' {
		p.advance()
	}

	// VERSION keyword-value pair: VERSION string_literal | VERSION NO VALUE
	if p.isIdentLikeStr("VERSION") {
		p.advance() // consume VERSION
		// VERSION NO VALUE or VERSION 'string'
		versionArg, parseErr753 := p.parseExpr()
		if parseErr753 != nil {
			return nil, parseErr753
		}
		if versionArg != nil {
			fc.Args.Items = append(fc.Args.Items, versionArg)
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseXmlAgg parses XMLAGG(expr ORDER BY ...).
//
//	XMLAGG ( expr [ ORDER BY sort_list ] )
func (p *Parser) parseXmlAgg() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume XMLAGG

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: "XMLAGG", Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// Expression argument
	arg, parseErr754 := p.parseExpr()
	if parseErr754 != nil {
		return nil, parseErr754
	}
	if arg != nil {
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		var parseErr755 error
		fc.OrderBy, parseErr755 = p.parseOrderByList()
		if parseErr755 != nil {
			return nil, parseErr755
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseXmlContentFunc parses XMLPARSE(CONTENT expr).
func (p *Parser) parseXmlContentFunc(name string) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume keyword

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// CONTENT or DOCUMENT keyword — skip
	if p.cur.Type == kwCONTENT || p.isIdentLikeStr("DOCUMENT") {
		p.advance()
	}

	arg, parseErr756 := p.parseExpr()
	if parseErr756 != nil {
		return nil, parseErr756
	}
	if arg != nil {
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseJsonObjectOrArray parses JSON_OBJECT or JSON_ARRAY expressions.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_OBJECT.html
//
//	JSON_OBJECT ( [ key_expr VALUE value_expr | key:value ] [, ...] )
//	JSON_ARRAY ( [ expr [, ...] ] )
func (p *Parser) parseJsonObjectOrArray(name string) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume keyword

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		var arg nodes.ExprNode
		var parseErr757 error
		if name == "JSON_OBJECT" {
			arg, parseErr757 = p.parseExprPrec(precComp)
		} else {
			arg, parseErr757 = p.parseExpr()
		}
		if parseErr757 != nil {
			return nil, parseErr757
		}
		if arg != nil {
			fc.Args.Items = append(fc.Args.Items, arg)
		}
		// Skip VALUE/IS keyword (JSON_OBJECT key VALUE val or key IS val syntax)
		if p.isIdentLikeStr("VALUE") || (name == "JSON_OBJECT" && p.cur.Type == kwIS) {
			p.advance()
			val, parseErr758 := p.parseExpr()
			if parseErr758 != nil {
				return nil, parseErr758
			}
			if val != nil {
				fc.Args.Items = append(fc.Args.Items, val)
			}
		}
		// Skip FORMAT JSON
		if p.cur.Type == kwFORMAT {
			p.advance()
			if p.cur.Type == kwJSON {
				p.advance()
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Skip trailing keywords: NULL ON NULL, ABSENT ON NULL, RETURNING clause, STRICT/LAX
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		p.advance()
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseJsonPathFunc parses JSON path functions: JSON_VALUE, JSON_QUERY, JSON_EXISTS, JSON_MERGEPATCH.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_VALUE.html
//
//	JSON_VALUE ( expr, path_string [ RETURNING type ] [ error_clause ] )
//	JSON_QUERY ( expr, path_string [ RETURNING type ] [ wrapper_clause ] [ error_clause ] )
//	JSON_EXISTS ( expr, path_string [ error_clause ] )
//	JSON_MERGEPATCH ( expr, expr [ RETURNING type ] )
func (p *Parser) parseJsonPathFunc(name string) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume keyword

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// First argument (JSON expr)
	arg, parseErr759 := p.parseExpr()
	if parseErr759 != nil {
		return nil, parseErr759
	}
	if arg != nil {
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	if p.cur.Type == ',' {
		p.advance()
		// Second argument (path or patch expr)
		arg2, parseErr760 := p.parseExpr()
		if parseErr760 != nil {
			return nil, parseErr760
		}
		if arg2 != nil {
			fc.Args.Items = append(fc.Args.Items, arg2)
		}
	}

	// Skip trailing keywords: RETURNING type, error clauses, wrapper clauses
	// These may include keywords like RETURNING, ERROR, NULL, DEFAULT, EMPTY, etc.
	// We consume everything until the closing paren.
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			if depth == 0 {
				break
			}
			depth--
		}
		p.advance()
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}

// parseXmlSerialize parses XMLSERIALIZE(CONTENT expr AS type).
func (p *Parser) parseXmlSerialize() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume XMLSERIALIZE

	fc := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: "XMLSERIALIZE", Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type != '(' {
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	p.advance() // consume '('

	// CONTENT or DOCUMENT keyword — skip
	if p.cur.Type == kwCONTENT || p.isIdentLikeStr("DOCUMENT") {
		p.advance()
	}

	arg, parseErr761 := p.parseExpr()
	if parseErr761 != nil {
		return nil, parseErr761
	}
	if arg != nil {
		fc.Args.Items = append(fc.Args.Items, arg)
	}

	// AS type
	if p.cur.Type == kwAS {
		p.advance()
		typeName, parseErr762 := p.parseTypeName()
		if parseErr762 != nil {
			return nil, parseErr762
		}
		if typeName != nil {
			fc.Args.Items = append(fc.Args.Items, typeName)
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return p.parseFuncCallPostfix(fc)
}
