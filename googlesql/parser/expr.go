package parser

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// This file is the `expressions` DAG node. It implements GoogleSQL's full
// expression grammar — the precedence chain (expression /
// expression_higher_prec_than_and / and_expression), every primary form, the
// operator / predicate / constructor families, function calls with all their
// clauses, CASE / CAST / EXTRACT, array / struct constructors, field & array
// access, INTERVAL, AT TIME ZONE, named arguments, lambdas, and the subquery
// keywords (ARRAY(…), EXISTS(…), parenthesized query) — as a hand-written
// recursive-descent + precedence-climbing parser over the token stream,
// producing ast.Node expression trees.
//
// Legacy ZetaSQL ANTLR rules (GoogleSQLParser.g4 §2.17 / §2.18), a hand-port of
// Google's open-source GoogleSQL reference. The implementation is adjudicated
// against the LIVE Cloud Spanner emulator oracle (oracle.md), not the literal
// ANTLR alternative ordering — ANTLR's left-recursive alternative order for the
// one big `expression_higher_prec_than_and` rule does NOT reflect operator
// precedence (it is a bison-port artifact; the real precedence lives in the
// referenced bison %left/%nonassoc declarations and the official GoogleSQL
// operator-precedence table). This parser uses standard precedence climbing,
// verified against the emulator. Two oracle-confirmed precedence facts:
//
//	P1 — the comparison family is NON-ASSOCIATIVE. After one comparison-family
//	     operator (comparative / [NOT] LIKE / [NOT] IN / [NOT] BETWEEN /
//	     IS [NOT] {NULL|TRUE|FALSE|UNKNOWN} / IS [NOT] DISTINCT FROM) you may not
//	     chain another at the same level: `a = b = c`, `1 < 2 < 3`,
//	     `x IN (1) IN (2)`, `1 IS NULL IS NULL`, `'a' LIKE 'b' LIKE 'c'`,
//	     `1 BETWEEN 0 AND 2 BETWEEN 0 AND 1` all REJECT. (Oracle: "Syntax error"
//	     / "Expression to the left of IS must be parenthesized".)
//	P2 — prefix NOT binds LOOSER than comparison. `NOT a = b` is `NOT (a = b)`,
//	     `NOT a IS NULL` is `NOT (a IS NULL)` (oracle: both accept). NOT sits
//	     between the comparison level and AND in precedence.

// Binding powers for precedence climbing (lowest → highest). bpCompare is the
// non-associative comparison-family level; it is handled specially in the loop
// (it never recurses back into another comparison). bpNot is the prefix-NOT
// level, between comparison and AND (P2).
const (
	bpNone    = iota
	bpOr      // OR
	bpAnd     // AND
	bpNot     // NOT (prefix) — P2
	bpCompare // = != <> < <= > >=, LIKE, IN, BETWEEN, IS, DISTINCT (non-assoc, P1)
	bpBitOr   // |
	bpBitXor  // ^
	bpBitAnd  // &
	bpConcat  // ||  (BOOL_OR_SYMBOL)
	bpShift   // << >>
	bpAdd     // + -
	bpMul     // * /
	bpUnary   // unary + - ~
	bpAccess  // postfix . [] .(path)
)

// ParseExpression parses a complete GoogleSQL expression from a standalone
// string, returning the ast.Node and any ParseErrors. Trailing tokens after the
// expression are reported as an error (the expression is still returned). It is
// the string-input counterpart of parseExpr, for tests and callers (and the
// differential oracle harness) that hold an expression string rather than a
// token stream. Mirrors ParseDataType / snowflake's parser entry points.
func ParseExpression(input string) (ast.Node, []ParseError) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()

	node, err := p.parseExpr()
	if err != nil {
		return nil, errToSlice(err)
	}
	if p.cur.Type != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Type)
		}
		return node, []ParseError{{Loc: p.cur.Loc, Msg: "unexpected token after expression: " + text}}
	}
	// Surface any best-effort errors collected during the parse (e.g. the
	// grammar's notify-style "SELECT as a query argument" alternatives, which
	// parse a tree but flag a diagnostic).
	return node, p.errors
}

// errToSlice converts a parse error into the []ParseError result shape.
func errToSlice(err error) []ParseError {
	if pe, ok := err.(*ParseError); ok {
		return []ParseError{*pe}
	}
	return []ParseError{{Msg: err.Error()}}
}

// ---------------------------------------------------------------------------
// Precedence-climbing core
// ---------------------------------------------------------------------------

// parseExpr is the expression entry point. It parses a full `expression`,
// lowest precedence first: OR. The precedence is layered as explicit recursive
// descents for the bands where associativity matters distinctly —
//
//	parseExpr → OR (left-assoc loop)
//	          → AND (left-assoc loop)
//	          → NOT (prefix, P2)
//	          → comparison family (NON-associative single step, P1)
//	          → bitwise/shift/arithmetic (left-assoc precedence-climbing loop)
//	          → unary prefix → primary → postfix access
//
// — so the comparison family is parsed by exactly ONE step (parseComparison),
// never folded twice at the same level. That is what makes `a = b = c`,
// `1 < 2 < 3`, `x IN (1) IN (2)`, `1 IS NULL IS NULL`, etc. reject (P1): after
// one comparison operator, no second comparison operator is accepted; the
// trailing operator is left as an unexpected token.
func (p *Parser) parseExpr() (ast.Node, error) {
	return p.parseOr()
}

// parseOr parses the OR level: `and_expr (OR and_expr)*` (left-associative).
func (p *Parser) parseOr() (ast.Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwOR {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: ast.BinOr, Left: left, Right: right, Loc: spanNodes(left, right)}
	}
	return left, nil
}

// parseAnd parses the AND level: `not_expr (AND not_expr)*` (left-associative,
// matching the grammar's and_expression).
func (p *Parser) parseAnd() (ast.Node, error) {
	left, err := p.parseNotExpr()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == kwAND {
		p.advance()
		right, err := p.parseNotExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: ast.BinAnd, Left: left, Right: right, Loc: spanNodes(left, right)}
	}
	return left, nil
}

// parseNotExpr parses prefix NOT (P2: NOT binds looser than comparison, so its
// operand is a full comparison-level expression). `NOT NOT x` chains. A
// non-NOT token falls through to the comparison level.
func (p *Parser) parseNotExpr() (ast.Node, error) {
	if p.cur.Type == kwNOT {
		tok := p.advance()
		operand, err := p.parseNotExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: ast.UnaryNot, Expr: operand, Loc: locFromTo2(tok.Loc, operand)}, nil
	}
	return p.parseComparison()
}

// parseComparison parses the NON-ASSOCIATIVE comparison family (P1): a single
// higher-precedence operand, then AT MOST ONE comparison-family operator
// (comparative / [NOT] LIKE / [NOT] IN / [NOT] BETWEEN / IS […] /
// IS […] DISTINCT FROM) whose RHS is again a higher-precedence operand. After
// one operator it returns — a second comparison operator at the same level is
// NOT consumed, so the caller (and ultimately parseSingle/ParseExpression)
// reports it as an unexpected token. The higher band (bitwise/shift/arith) is
// parsed by parseBinaryExpr.
func (p *Parser) parseComparison() (ast.Node, error) {
	left, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case int('='):
		return p.parseCompareRHS(left, ast.CmpEq)
	case tokNotEqual, tokNotEqual2:
		return p.parseCompareRHS(left, ast.CmpNe)
	case int('<'):
		return p.parseCompareRHS(left, ast.CmpLt)
	case tokLessEqual:
		return p.parseCompareRHS(left, ast.CmpLe)
	case int('>'):
		return p.parseCompareRHS(left, ast.CmpGt)
	case tokGreaterEqual:
		return p.parseCompareRHS(left, ast.CmpGe)
	case kwIS:
		return p.parseIs(left)
	case kwBETWEEN:
		return p.parseBetween(left, false)
	case kwIN:
		return p.parseIn(left, false)
	case kwLIKE:
		return p.parseLike(left, false)
	case kwNOT:
		// NOT { LIKE | IN | BETWEEN } — a comparison-family lead-in.
		switch p.peekNext().Type {
		case kwLIKE, kwIN, kwBETWEEN:
			p.advance() // consume NOT
			switch p.cur.Type {
			case kwLIKE:
				return p.parseLike(left, true)
			case kwIN:
				return p.parseIn(left, true)
			case kwBETWEEN:
				return p.parseBetween(left, true)
			}
		}
		// A bare NOT not followed by a comparison-family keyword is not an infix
		// operator here; leave it (the caller surfaces it).
		return left, nil
	}
	return left, nil
}

// parseCompareRHS consumes a comparative operator and its RHS. The plain RHS is
// parseBinaryExpr(bpBitOr) — NOT another comparison — so the family stays
// non-associative (P1). The quantified form (`op {ANY|SOME|ALL} <rhs>`, the
// any_some_all production) is also accepted here: GoogleSQL allows
// `expr comparative_operator any_some_all (subquery | list | UNNEST)` exactly as
// it allows `expr LIKE any_some_all ...`. (Verified against the live Spanner
// emulator: all six operators × ANY/SOME/ALL × subquery/list/UNNEST parse —
// see TestExprDifferential and divergence #201.)
func (p *Parser) parseCompareRHS(left ast.Node, op ast.CompareOp) (ast.Node, error) {
	p.advance() // consume the operator

	// Quantified form: comparative_operator { ANY | SOME | ALL } <rhs>.
	switch p.cur.Type {
	case kwANY, kwSOME, kwALL:
		quant := strings.ToUpper(TokenName(p.cur.Type))
		p.advance()
		rhs, err := p.parseQuantifiedRHS(nodeLoc(left).Start)
		if err != nil {
			return nil, err
		}
		return &ast.CompareExpr{
			Op:            op,
			Left:          left,
			Quantifier:    quant,
			QuantValues:   rhs.values,
			QuantUnnest:   rhs.unnest,
			QuantSubquery: rhs.subquery,
			Loc:           ast.Loc{Start: nodeLoc(left).Start, End: rhs.end},
		}, nil
	}

	right, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	return &ast.CompareExpr{Op: op, Left: left, Right: right, Loc: spanNodes(left, right)}, nil
}

// quantifiedRHS is the parsed right-hand side of a quantified operator
// (`{ANY|SOME|ALL} <rhs>`), shared by the comparative and LIKE forms. Exactly
// one of values / unnest / subquery is set; end is the offset past the RHS.
type quantifiedRHS struct {
	values   []ast.Node // parenthesized expression list (>= 1)
	unnest   ast.Node   // UNNEST(...) call
	subquery ast.Node   // parenthesized query
	end      int        // end offset past the RHS
}

// parseQuantifiedRHS parses the RHS after an `ANY|SOME|ALL` quantifier has been
// consumed: an optional leading `@{...}`/`@N` hint, then either `UNNEST(...)`, a
// parenthesized query, or a parenthesized expression list (>= 1). lhsStart is the
// start offset of the overall expression's left operand, used for the span. This
// is the any_some_all RHS shared by `expr op {ANY|SOME|ALL} ...` and
// `expr LIKE {ANY|SOME|ALL} ...` (the grammar's anysomeall list/subquery RHS).
func (p *Parser) parseQuantifiedRHS(lhsStart int) (quantifiedRHS, error) {
	var out quantifiedRHS
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return out, herr
		}
	}
	if p.cur.Type == kwUNNEST {
		unnest, err := p.parseUnnestExpression()
		if err != nil {
			return out, err
		}
		out.unnest = unnest
		out.end = nodeLoc(unnest).End
		return out, nil
	}
	if _, err := p.expect(int('(')); err != nil {
		return out, err
	}
	if p.atQueryStart() {
		sub, end, err := p.parseSubqueryBody(lhsStart)
		if err != nil {
			return out, err
		}
		out.subquery = sub
		out.end = end
		return out, nil
	}
	values, endLoc, err := p.parseExprListThroughParen()
	if err != nil {
		return out, err
	}
	out.values = values
	out.end = endLoc.End
	return out, nil
}

// parseBinaryExpr parses the left-associative arithmetic / bitwise / shift /
// concat band by precedence climbing with the given minimum binding power
// (lowest in this band is bpBitOr). Below this band lives the comparison family
// (handled by the caller); above it live unary and postfix (handled by
// parsePrefixExpr).
func (p *Parser) parseBinaryExpr(minBP int) (ast.Node, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	for {
		op, bp, ok := p.binaryOpAt()
		if !ok || bp < minBP {
			break
		}
		p.advance() // consume the operator
		right, err := p.parseBinaryExpr(bp + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right, Loc: spanNodes(left, right)}
	}
	return left, nil
}

// binaryOpAt returns the BinaryOp, binding power, and ok for the current token
// as a left-associative binary operator in the bitwise/shift/arithmetic band.
// (OR / AND / comparison are handled by their dedicated layers, not here.)
func (p *Parser) binaryOpAt() (ast.BinaryOp, int, bool) {
	switch p.cur.Type {
	case int('|'):
		return ast.BinBitOr, bpBitOr, true
	case int('^'):
		return ast.BinBitXor, bpBitXor, true
	case int('&'):
		return ast.BinBitAnd, bpBitAnd, true
	case tokBoolOr:
		return ast.BinConcat, bpConcat, true
	case tokShiftLeft:
		return ast.BinShiftLeft, bpShift, true
	case tokShiftRight:
		return ast.BinShiftRight, bpShift, true
	case int('+'):
		return ast.BinAdd, bpAdd, true
	case int('-'):
		return ast.BinSub, bpAdd, true
	case int('*'):
		return ast.BinMul, bpMul, true
	case int('/'):
		return ast.BinDiv, bpMul, true
	}
	return 0, 0, false
}

// parseIs parses `IS [NOT] { NULL | TRUE | FALSE | UNKNOWN }` or
// `IS [NOT] DISTINCT FROM expr` (is_operator / distinct_operator). The IS
// keyword is the current token. Non-associative (P1): a DISTINCT FROM RHS is
// parsed at bpCompare+1.
func (p *Parser) parseIs(left ast.Node) (ast.Node, error) {
	p.advance() // consume IS
	not := false
	if p.cur.Type == kwNOT {
		p.advance()
		not = true
	}
	// DISTINCT FROM form.
	if p.cur.Type == kwDISTINCT {
		p.advance() // DISTINCT
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		right, err := p.parseBinaryExpr(bpBitOr)
		if err != nil {
			return nil, err
		}
		return &ast.IsExpr{Expr: left, Not: not, DistinctFrom: right, Loc: spanNodes(left, right)}, nil
	}
	// Predicate form: NULL / TRUE / FALSE / UNKNOWN.
	var pred string
	switch p.cur.Type {
	case kwNULL:
		pred = "NULL"
	case kwTRUE:
		pred = "TRUE"
	case kwFALSE:
		pred = "FALSE"
	case kwUNKNOWN:
		pred = "UNKNOWN"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return &ast.IsExpr{Expr: left, Not: not, Pred: pred, Loc: locFromTo(left, tok)}, nil
}

// parseBetween parses `[NOT] BETWEEN low AND high` (between_operator). The
// BETWEEN keyword is the current token; not records a preceding NOT. low and
// high are parsed at the EHP (higher-than-AND) level so a bare OR/AND inside is
// not consumed (the grammar's "Expression in BETWEEN must be parenthesized"
// alternative — verified: `a BETWEEN 1 OR 2` rejects). Concretely, low and high
// parse at bpCompare+1: above AND/OR/NOT and above the comparison family, so the
// AND separator is the BETWEEN's own AND, not a logical conjunction.
func (p *Parser) parseBetween(left ast.Node, not bool) (ast.Node, error) {
	p.advance() // consume BETWEEN
	low, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAND); err != nil {
		return nil, err
	}
	high, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	return &ast.BetweenExpr{Expr: left, Low: low, High: high, Not: not, Loc: spanNodes(left, high)}, nil
}

// parseIn parses `[NOT] IN <rhs>` (in_operator). The IN keyword is the current
// token. RHS forms: `UNNEST(...)`, a parenthesized query, or a parenthesized
// expression list (>= 1). A leading `@{...}` hint is allowed before a
// non-UNNEST RHS (the grammar errors on a hint before UNNEST — we record the
// diagnostic but still parse, matching the notify-style oracle behavior).
func (p *Parser) parseIn(left ast.Node, not bool) (ast.Node, error) {
	inTok := p.advance() // consume IN

	// Optional hint @{...} / @N before the RHS.
	hadHint := false
	if p.cur.Type == int('@') {
		hadHint = true
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	out := &ast.InExpr{Expr: left, Not: not}

	// UNNEST(...) RHS.
	if p.cur.Type == kwUNNEST {
		if hadHint {
			p.errors = append(p.errors, ParseError{
				Loc: inTok.Loc,
				Msg: "Syntax error: HINTs cannot be specified on IN clause with UNNEST",
			})
		}
		unnest, err := p.parseUnnestExpression()
		if err != nil {
			return nil, err
		}
		out.Unnest = unnest
		out.Loc = spanNodes(left, unnest)
		return out, nil
	}

	// Parenthesized RHS: a query, or a (one-or-more) expression list.
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if p.atQueryStart() {
		sub, end, err := p.parseSubqueryBody(inTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		out.Subquery = sub
		out.Loc = ast.Loc{Start: nodeLoc(left).Start, End: end}
		return out, nil
	}
	// Expression list.
	values, endLoc, err := p.parseExprListThroughParen()
	if err != nil {
		return nil, err
	}
	out.Values = values
	out.Loc = ast.Loc{Start: nodeLoc(left).Start, End: endLoc.End}
	return out, nil
}

// parseLike parses `[NOT] LIKE <rhs>` (like_operator). The LIKE keyword is the
// current token. Plain form: `LIKE pattern`. Quantified form: `LIKE
// {ANY|SOME|ALL} [hint] <list|UNNEST|query>`.
func (p *Parser) parseLike(left ast.Node, not bool) (ast.Node, error) {
	p.advance() // consume LIKE
	out := &ast.LikeExpr{Expr: left, Not: not}

	// Quantified form: ANY / SOME / ALL. Shares the RHS grammar (hint + UNNEST /
	// subquery / list) with the comparative-operator quantified form via
	// parseQuantifiedRHS.
	switch p.cur.Type {
	case kwANY, kwSOME, kwALL:
		out.Quantifier = strings.ToUpper(TokenName(p.cur.Type))
		p.advance()
		rhs, err := p.parseQuantifiedRHS(nodeLoc(left).Start)
		if err != nil {
			return nil, err
		}
		out.QuantValues = rhs.values
		out.QuantUnnest = rhs.unnest
		out.QuantSubquery = rhs.subquery
		out.Loc = ast.Loc{Start: nodeLoc(left).Start, End: rhs.end}
		return out, nil
	}

	// Plain form: a single pattern at the comparison-RHS level (non-assoc).
	pattern, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	out.Pattern = pattern
	out.Loc = spanNodes(left, pattern)
	return out, nil
}

// ---------------------------------------------------------------------------
// Prefix / primary dispatch
// ---------------------------------------------------------------------------

// parsePrefixExpr handles the prefix unary operators (+ - ~) and otherwise
// delegates to parsePrimaryExpr followed by postfix-access folding. (Prefix NOT
// is handled higher up by parseNotExpr — P2 — because it binds looser than the
// comparison family, while + - ~ bind tighter than multiplication.) Unary + - ~
// take the official GoogleSQL precedence: above multiplicative (level 2), so
// `-a*b` groups as `(-a)*b` — semantically identical to `-(a*b)`, so this choice
// is unobservable through the oracle and follows the documented precedence
// table (truth1).
func (p *Parser) parsePrefixExpr() (ast.Node, error) {
	switch p.cur.Type {
	case int('+'):
		return p.parseUnary(ast.UnaryPlus)
	case int('-'):
		return p.parseUnary(ast.UnaryMinus)
	case int('~'):
		return p.parseUnary(ast.UnaryBitNot)
	}

	prim, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	return p.parsePostfix(prim)
}

// parseUnary consumes a prefix unary operator and its operand at bpUnary, so a
// following higher-or-equal-precedence operator (another unary, or postfix
// access) binds to the operand and `- - 1` / `2 - - 1` / `-a[0]` parse.
func (p *Parser) parseUnary(op ast.UnaryOp) (ast.Node, error) {
	tok := p.advance()
	operand, err := p.parseBinaryExpr(bpUnary)
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpr{Op: op, Expr: operand, Loc: locFromTo2(tok.Loc, operand)}, nil
}

// parsePostfix folds postfix-access operators (`. field`, `. (path)`,
// `[ index ]`) onto base. Access is left-associative and the highest precedence,
// so it always binds before returning to the caller.
func (p *Parser) parsePostfix(base ast.Node) (ast.Node, error) {
	for {
		switch p.cur.Type {
		case int('.'):
			next, err := p.parseDotAccess(base)
			if err != nil {
				return nil, err
			}
			base = next
		case int('['):
			next, err := p.parseIndexAccess(base)
			if err != nil {
				return nil, err
			}
			base = next
		default:
			return base, nil
		}
	}
}

// parseDotAccess parses a `. field` or `. (path)` postfix on base
// (EHP DOT identifier | EHP DOT '(' path ')'). The '.' is the current token.
func (p *Parser) parseDotAccess(base ast.Node) (ast.Node, error) {
	p.advance() // consume '.'
	// Extension access: . ( path ).
	if p.cur.Type == int('(') {
		p.advance() // consume '('
		path, err := p.parsePathExpr()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		return &ast.ExtensionAccess{Expr: base, Path: path, Loc: ast.Loc{Start: nodeLoc(base).Start, End: closeTok.Loc.End}}, nil
	}
	// Field access: . identifier (any keyword permitted as a field name part).
	if !isAnyKeywordIdentifier(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	name, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	return &ast.FieldAccess{Expr: base, Field: name, Loc: locFromTo(base, tok)}, nil
}

// parseIndexAccess parses a `[ index ]` subscript on base
// (EHP LS_BRACKET expression RS_BRACKET). The '[' is the current token. The
// index expression may be an OFFSET/ORDINAL/SAFE_OFFSET/SAFE_ORDINAL function
// call — those are ordinary function calls, not special syntax.
func (p *Parser) parseIndexAccess(base ast.Node) (ast.Node, error) {
	p.advance() // consume '['
	idx, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(']'))
	if err != nil {
		return nil, err
	}
	return &ast.IndexAccess{Expr: base, Index: idx, Loc: ast.Loc{Start: nodeLoc(base).Start, End: closeTok.Loc.End}}, nil
}

// parsePrimaryExpr parses a primary (atomic) expression — everything in the
// `expression_higher_prec_than_and` rule that is not an operator application:
// literals, parameters, constructors, CASE/CAST/EXTRACT, function calls,
// identifiers/paths, parenthesized expressions, and subquery keywords.
func (p *Parser) parsePrimaryExpr() (ast.Node, error) {
	switch p.cur.Type {
	// --- literals ---
	case kwNULL:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitNull, Value: "NULL", Loc: tok.Loc}, nil
	case kwTRUE, kwFALSE:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitBool, Value: TokenName(tok.Type), Loc: tok.Loc}, nil
	case tokInteger:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}, nil
	case tokFloat:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}, nil
	case tokString:
		return p.parseStringLiteral()
	case tokBytes:
		return p.parseBytesLiteral()

	// --- typed-prefix literals (NUMERIC '…', DATE '…', JSON '…', RANGE<…> '…') ---
	case kwNUMERIC, kwDECIMAL, kwBIGNUMERIC, kwBIGDECIMAL, kwJSON,
		kwDATE, kwTIME, kwDATETIME, kwTIMESTAMP:
		// These are also valid bare type names / identifiers; here they are typed
		// literals ONLY when immediately followed by a string literal. Otherwise
		// fall through to identifier/function-call handling.
		if p.peekNext().Type == tokString {
			return p.parseTypedLiteral()
		}
		return p.parseNameOrCall()

	// --- parameters / system variables ---
	case int('@'):
		return p.parseNamedParameter()
	case tokAtAt:
		return p.parseSystemVariable()
	case int('?'):
		tok := p.advance()
		return &ast.Parameter{Positional: true, Loc: tok.Loc}, nil

	// --- CASE / CAST / EXTRACT / INTERVAL ---
	case kwCASE:
		return p.parseCaseExpr()
	case kwCAST:
		return p.parseCastExpr(false)
	case kwSAFE_CAST:
		return p.parseCastExpr(true)
	case kwEXTRACT:
		return p.parseExtractExpr()
	case kwINTERVAL:
		return p.parseIntervalExpr()

	// --- constructors ---
	case int('['):
		// Bare array constructor [ … ].
		return p.parseArrayConstructor(false, nil)
	case kwARRAY:
		return p.parseArrayOrArraySubquery()
	case kwSTRUCT:
		return p.parseStructOrTyped()
	case kwNEW:
		return p.parseNewConstructor()
	case kwREPLACE_FIELDS:
		return p.parseReplaceFields()
	case kwWITH:
		return p.parseWithExpr()
	case int('{'):
		return p.parseBracedConstructor(nil)

	// --- EXISTS subquery ---
	case kwEXISTS:
		return p.parseExistsExpr()

	// --- parenthesized expression or subquery ---
	case int('('):
		return p.parseParenOrSubquery()

	// --- keyword function names (IF/GROUPING/LEFT/RIGHT/COLLATE/RANGE) ---
	case kwIF, kwGROUPING, kwLEFT, kwRIGHT, kwCOLLATE:
		return p.parseKeywordFuncCall()
	case kwRANGE:
		// RANGE is a keyword function name AND a type prefix (RANGE<T> '…'). A
		// following '<' is the typed RANGE literal; a following '(' is the
		// keyword function call.
		if p.peekNext().Type == int('<') {
			return p.parseRangeLiteral()
		}
		return p.parseKeywordFuncCall()

	// --- identifier / path / function call ---
	default:
		if isIdentifierStart(p.cur.Type) {
			return p.parseNameOrCall()
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// ---------------------------------------------------------------------------
// Identifiers, paths, and function calls
// ---------------------------------------------------------------------------

// parseNameOrCall parses an identifier-led primary: either a path expression
// (`a`, `a.b.c`) used as a column/name reference, or a function call when a
// path is immediately followed by `(`. The grammar resolves this exactly as
// ZetaSQL's LALR action does: `path_expression '(' …` is a call; otherwise the
// path is a name reference (further `.field` / `[idx]` access is folded by
// parsePostfix). The current token begins an identifier.
func (p *Parser) parseNameOrCall() (ast.Node, error) {
	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == int('(') {
		return p.parseFuncCallSuffix(path)
	}
	if len(path.Parts) == 1 {
		return &ast.Identifier{Name: path.Parts[0], Loc: path.Loc}, nil
	}
	return path, nil
}

// parsePathExpr parses `path_expression: identifier (DOT identifier)*`. The
// first component must be an identifier-start (token identifier or non-reserved
// keyword); dotted continuations accept any keyword as a name part (matching the
// foundation's permissive path-component rule, also used by parseType). The
// current token begins the path.
func (p *Parser) parsePathExpr() (*ast.PathExpr, error) {
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	first := p.advance()
	part0, err := p.identifierText(first)
	if err != nil {
		return nil, err
	}
	path := &ast.PathExpr{Parts: []string{part0}, Loc: first.Loc}
	for p.cur.Type == int('.') {
		// Only treat '.' as a path continuation when followed by a name part. A
		// '.(' is extension access (handled by parsePostfix), and a '.' followed
		// by a non-name is a syntax error surfaced by the caller — but here we
		// stop so the caller/postfix layer handles those uniformly.
		if !isAnyKeywordIdentifier(p.peekNext().Type) {
			break
		}
		p.advance() // consume '.'
		partTok := p.advance()
		part, err := p.identifierText(partTok)
		if err != nil {
			return nil, err
		}
		path.Parts = append(path.Parts, part)
		path.Loc.End = partTok.Loc.End
	}
	return path, nil
}

// parseFuncCallSuffix parses a function call given its already-parsed name path.
// The current token is the opening '('. Implements
// function_call_expression_with_clauses_suffix: optional DISTINCT, the argument
// list (or empty / bare '*'), the null-handling / having / clamped / with-report
// modifiers, trailing ORDER BY / LIMIT, then the closing ')', and the
// post-paren hint / WITH GROUP ROWS / OVER clauses.
func (p *Parser) parseFuncCallSuffix(name *ast.PathExpr) (ast.Node, error) {
	p.advance() // consume '('
	fc := &ast.FuncCall{Name: name, Loc: name.Loc}

	if p.cur.Type == kwDISTINCT {
		p.advance()
		fc.Distinct = true
	}

	// Argument list (possibly empty, or a bare '*').
	if p.cur.Type != int(')') {
		if err := p.parseFuncArgs(fc); err != nil {
			return nil, err
		}
		// Aggregate / DP modifiers, in grammar order.
		if err := p.parseFuncModifiers(fc); err != nil {
			return nil, err
		}
	} else {
		// Empty-arg list still allows opt_having / ORDER BY / LIMIT before ')'.
		if err := p.parseEmptyArgModifiers(fc); err != nil {
			return nil, err
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End

	// Post-paren: hint? with_group_rows? over_clause?
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
		fc.Loc.End = p.prev.Loc.End
	}
	if p.cur.Type == kwWITH && p.peekNext().Type == kwGROUP {
		// WITH GROUP ROWS.
		p.advance() // WITH
		p.advance() // GROUP
		rowsTok, err := p.expect(kwROWS)
		if err != nil {
			return nil, err
		}
		fc.WithGroupRows = true
		fc.Loc.End = rowsTok.Loc.End
	}
	if p.cur.Type == kwOVER {
		win, err := p.parseOverClause()
		if err != nil {
			return nil, err
		}
		fc.Over = win
		fc.Loc.End = win.Loc.End
	}
	return fc, nil
}

// parseKeywordFuncCall parses a function call whose name is one of the keyword
// function names (function_name_from_keyword: IF / GROUPING / LEFT / RIGHT /
// COLLATE / RANGE). The current token is the keyword; it must be followed by
// '('. The keyword spelling becomes the function name.
func (p *Parser) parseKeywordFuncCall() (ast.Node, error) {
	kw := p.advance()
	name := &ast.PathExpr{Parts: []string{TokenName(kw.Type)}, Loc: kw.Loc}
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseFuncCallSuffix(name)
}

// parseFuncArgs parses a non-empty function argument list into fc.Args. The
// first token is the first argument (not ')'). Each argument is a positional
// expression (optionally `AS alias` for a TVF arg), a named argument
// (`id => expr`), a lambda, a sequence arg, or a bare '*' (only as the sole
// leading token). A bare SELECT argument is the grammar's notify-style error
// alternative: it parses (as a flagged subquery) so the call is still accepted.
func (p *Parser) parseFuncArgs(fc *ast.FuncCall) error {
	// Bare '*' argument (COUNT(*)).
	if p.cur.Type == int('*') {
		star := p.advance()
		fc.Args = append(fc.Args, &ast.StarExpr{Loc: star.Loc})
		for p.cur.Type == int(',') {
			p.advance()
			arg, err := p.parseFuncArg()
			if err != nil {
				return err
			}
			fc.Args = append(fc.Args, arg)
		}
		return nil
	}
	arg, err := p.parseFuncArg()
	if err != nil {
		return err
	}
	fc.Args = append(fc.Args, arg)
	for p.cur.Type == int(',') {
		p.advance()
		a, err := p.parseFuncArg()
		if err != nil {
			return err
		}
		fc.Args = append(fc.Args, a)
	}
	return nil
}

// parseFuncArg parses one function_call_argument: a named arg (`id => …`), a
// lambda (`params -> body` / `() -> body`), a sequence arg (`SEQUENCE path`), a
// bare SELECT (notify-style accept), or a positional expression with an optional
// `AS alias`.
func (p *Parser) parseFuncArg() (ast.Node, error) {
	// SEQUENCE path arg.
	if p.cur.Type == kwSEQUENCE {
		seqTok := p.advance()
		path, err := p.parsePathExpr()
		if err != nil {
			return nil, err
		}
		return &ast.SequenceArg{Path: path, Loc: ast.Loc{Start: seqTok.Loc.Start, End: path.Loc.End}}, nil
	}
	// Bare SELECT as an argument: the grammar's error alternative. It flags a
	// diagnostic but the call parses (oracle: `func(SELECT 1)` is ACCEPT with a
	// semantic-style message). We consume the subquery body and record the
	// notify-style message so accept/reject parity holds.
	if p.cur.Type == kwSELECT {
		startLoc := p.cur.Loc
		p.errors = append(p.errors, ParseError{
			Loc: startLoc,
			Msg: "Each function argument is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression",
		})
		// Consume the bare SELECT body up to the call's closing ')'/',' at depth 0.
		sub := p.captureBareSelectArg(startLoc)
		return sub, nil
	}

	// Named argument: identifier '=>' …  (lookahead: ident then '=>').
	if isIdentifierStart(p.cur.Type) && p.peekNext().Type == tokFatArrow {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		p.advance() // consume '=>'
		// The value may itself be a lambda.
		val, err := p.parseArgValueMaybeLambda()
		if err != nil {
			return nil, err
		}
		return &ast.NamedArg{Name: name, Value: val, Loc: ast.Loc{Start: nameTok.Loc.Start, End: nodeLoc(val).End}}, nil
	}

	// Lambda with a parenthesized parameter list: '(' id (',' id)* ')' '->' …
	// or the empty '() -> …'. Disambiguated from a parenthesized expression by a
	// trailing '->'.
	if p.cur.Type == int('(') {
		if lam, ok, err := p.tryParseParenLambda(); ok || err != nil {
			return lam, err
		}
	}

	// Positional expression, possibly a single-identifier lambda `id -> body`.
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == tokArrow {
		// `id -> body` lambda: the parsed expr must be a single identifier.
		if id, ok := expr.(*ast.Identifier); ok {
			p.advance() // consume '->'
			body, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &ast.LambdaExpr{Params: []string{id.Name}, Body: body, Loc: ast.Loc{Start: id.Loc.Start, End: nodeLoc(body).End}}, nil
		}
	}
	// Optional `AS alias` (TVF / struct-style aliased argument).
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		return &ast.StructFieldExpr{Value: expr, Alias: alias, Loc: ast.Loc{Start: nodeLoc(expr).Start, End: aliasTok.Loc.End}}, nil
	}
	return expr, nil
}

// parseArgValueMaybeLambda parses a named-argument value, which may be a lambda
// (`id => lambda` is allowed by named_argument). It tries a parenthesized lambda
// first, then a single-identifier lambda, then a plain expression.
func (p *Parser) parseArgValueMaybeLambda() (ast.Node, error) {
	if p.cur.Type == int('(') {
		if lam, ok, err := p.tryParseParenLambda(); ok || err != nil {
			return lam, err
		}
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == tokArrow {
		if id, ok := expr.(*ast.Identifier); ok {
			p.advance()
			body, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &ast.LambdaExpr{Params: []string{id.Name}, Body: body, Loc: ast.Loc{Start: id.Loc.Start, End: nodeLoc(body).End}}, nil
		}
	}
	return expr, nil
}

// parseFuncModifiers parses the post-argument aggregate / DP modifiers in
// grammar order: opt_null_handling_modifier? opt_having_or_group_by_modifier?
// clamped_between_modifier? with_report_modifier? order_by_clause?
// limit_offset_clause?. The current token is just past the argument list.
func (p *Parser) parseFuncModifiers(fc *ast.FuncCall) error {
	// IGNORE NULLS | RESPECT NULLS.
	if p.cur.Type == kwIGNORE || p.cur.Type == kwRESPECT {
		kw := p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return err
		}
		fc.NullHandling = strings.ToUpper(TokenName(kw.Type)) + " NULLS"
	}
	// HAVING MAX expr | HAVING MIN expr group-by.
	if p.cur.Type == kwHAVING {
		havTok := p.advance()
		isMax := false
		switch p.cur.Type {
		case kwMAX:
			isMax = true
		case kwMIN:
			isMax = false
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance() // MAX/MIN
		hexpr, err := p.parseExpr()
		if err != nil {
			return err
		}
		fc.Having = &ast.HavingModifier{IsMax: isMax, Expr: hexpr, Loc: ast.Loc{Start: havTok.Loc.Start, End: nodeLoc(hexpr).End}}
		// The MIN form carries a trailing group-by; consume it best-effort
		// (GROUP [hint] [AND ORDER] BY grouping_item, …). Only relevant for
		// parse parity; the items are not retained.
		if !isMax && p.cur.Type == kwGROUP {
			if err := p.skipGroupByTail(); err != nil {
				return err
			}
		}
	}
	// CLAMPED BETWEEN low AND high.
	if p.cur.Type == kwCLAMPED {
		clTok := p.advance()
		low, err := p.parseBinaryExpr(bpBitOr)
		if err != nil {
			return err
		}
		if _, err := p.expect(kwAND); err != nil {
			return err
		}
		high, err := p.parseExpr()
		if err != nil {
			return err
		}
		fc.Clamped = &ast.ClampedModifier{Low: low, High: high, Loc: ast.Loc{Start: clTok.Loc.Start, End: nodeLoc(high).End}}
	}
	// WITH REPORT options_list.
	if p.cur.Type == kwWITH && p.peekNext().Type == kwREPORT {
		p.advance() // WITH
		p.advance() // REPORT
		if err := p.skipOptionsList(); err != nil {
			return err
		}
		fc.WithReport = true
	}
	// ORDER BY.
	if p.cur.Type == kwORDER {
		items, err := p.parseOrderByClause()
		if err != nil {
			return err
		}
		fc.OrderBy = items
	}
	// LIMIT expr [OFFSET expr].
	if p.cur.Type == kwLIMIT {
		if err := p.parseLimitOffset(fc); err != nil {
			return err
		}
	}
	return nil
}

// parseEmptyArgModifiers parses the empty-argument-list tail: opt_having? ORDER
// BY? LIMIT?. Only ORDER BY / LIMIT and the having modifier are valid here
// (the empty-list alternative in the grammar). The current token is just after
// '(' with no arguments.
func (p *Parser) parseEmptyArgModifiers(fc *ast.FuncCall) error {
	if p.cur.Type == kwHAVING {
		// Reuse the full modifier parser, which handles HAVING and the trailing
		// ORDER BY / LIMIT.
		return p.parseFuncModifiers(fc)
	}
	if p.cur.Type == kwORDER {
		items, err := p.parseOrderByClause()
		if err != nil {
			return err
		}
		fc.OrderBy = items
	}
	if p.cur.Type == kwLIMIT {
		if err := p.parseLimitOffset(fc); err != nil {
			return err
		}
	}
	return nil
}

// parseLimitOffset parses `LIMIT expr [OFFSET expr]` into fc.
func (p *Parser) parseLimitOffset(fc *ast.FuncCall) error {
	p.advance() // LIMIT
	lim, err := p.parseExpr()
	if err != nil {
		return err
	}
	fc.Limit = lim
	if p.cur.Type == kwOFFSET {
		p.advance()
		off, err := p.parseExpr()
		if err != nil {
			return err
		}
		fc.LimitOffset = off
	}
	return nil
}

// parseOverClause parses `OVER window_specification` (over_clause). OVER is the
// current token. The window spec is either a bare name or a parenthesized inline
// spec.
func (p *Parser) parseOverClause() (*ast.WindowSpec, error) {
	overTok := p.advance() // OVER
	if p.cur.Type == int('(') {
		return p.parseInlineWindowSpec(overTok.Loc.Start)
	}
	// Named window reference.
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	return &ast.WindowSpec{Name: name, Loc: ast.Loc{Start: overTok.Loc.Start, End: nameTok.Loc.End}}, nil
}

// parseInlineWindowSpec parses `( [name] [PARTITION BY …] [ORDER BY …]
// [frame] )` (window_specification). The current token is '('. startOff anchors
// the Loc start (the OVER keyword).
func (p *Parser) parseInlineWindowSpec(startOff int) (*ast.WindowSpec, error) {
	p.advance() // '('
	ws := &ast.WindowSpec{Inline: true, Loc: ast.Loc{Start: startOff, End: startOff}}
	// Optional leading base-window name.
	if isIdentifierStart(p.cur.Type) && p.cur.Type != kwPARTITION {
		// Only an identifier that is NOT the start of PARTITION/ORDER/frame is a
		// base name. PARTITION is non-reserved (isIdentifierStart true) but here
		// always introduces the partition clause, so exclude it.
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		ws.Name = name
	}
	if p.cur.Type == kwPARTITION {
		parts, err := p.parsePartitionBy()
		if err != nil {
			return nil, err
		}
		ws.PartitionBy = parts
	}
	if p.cur.Type == kwORDER {
		items, err := p.parseOrderByClause()
		if err != nil {
			return nil, err
		}
		ws.OrderBy = items
	}
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE {
		frame, err := p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
		ws.Frame = frame
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ws.Loc.End = closeTok.Loc.End
	return ws, nil
}

// parsePartitionBy parses `PARTITION [hint] BY expr (, expr)*`
// (partition_by_clause). PARTITION is the current token.
func (p *Parser) parsePartitionBy() ([]ast.Node, error) {
	p.advance() // PARTITION
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	return p.parseExprCommaList()
}

// parseOrderByClause parses `ORDER [hint] BY ordering_expression (, …)*`
// (order_by_clause). ORDER is the current token. Each item is
// `expr [COLLATE c] [ASC|DESC] [NULLS FIRST|LAST]`.
func (p *Parser) parseOrderByClause() ([]*ast.OrderItem, error) {
	p.advance() // ORDER
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	var items []*ast.OrderItem
	for {
		item, err := p.parseOrderItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	return items, nil
}

// parseOrderItem parses one `expr [COLLATE c] [ASC|DESC] [NULLS FIRST|LAST]`.
func (p *Parser) parseOrderItem() (*ast.OrderItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	item := &ast.OrderItem{Expr: expr, Loc: nodeLoc(expr)}
	if p.cur.Type == kwCOLLATE {
		p.advance()
		coll, err := p.parseBinaryExpr(bpBitOr)
		if err != nil {
			return nil, err
		}
		item.Collate = coll
		item.Loc.End = nodeLoc(coll).End
	}
	if p.cur.Type == kwASC || p.cur.Type == kwDESC {
		item.HasDir = true
		item.Desc = p.cur.Type == kwDESC
		item.Loc.End = p.cur.Loc.End
		p.advance()
	}
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			v := true
			item.NullsFirst = &v
		case kwLAST:
			v := false
			item.NullsFirst = &v
		default:
			return nil, p.syntaxErrorAtCur()
		}
		item.Loc.End = p.cur.Loc.End
		p.advance()
	}
	return item, nil
}

// parseWindowFrame parses `{ROWS|RANGE} { bound | BETWEEN bound AND bound }`
// (opt_window_frame_clause). The current token is ROWS or RANGE.
func (p *Parser) parseWindowFrame() (*ast.WindowFrame, error) {
	unitTok := p.advance() // ROWS / RANGE
	frame := &ast.WindowFrame{Loc: unitTok.Loc}
	if unitTok.Type == kwRANGE {
		frame.Kind = ast.FrameRange
	} else {
		frame.Kind = ast.FrameRows
	}
	if p.cur.Type == kwBETWEEN {
		p.advance()
		start, err := p.parseWindowBound()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		end, err := p.parseWindowBound()
		if err != nil {
			return nil, err
		}
		frame.Between = true
		frame.Start = start
		frame.End = end
		frame.Loc.End = p.prev.Loc.End
		return frame, nil
	}
	start, err := p.parseWindowBound()
	if err != nil {
		return nil, err
	}
	frame.Start = start
	frame.Loc.End = p.prev.Loc.End
	return frame, nil
}

// parseWindowBound parses one window_frame_bound: `UNBOUNDED {PRECEDING|
// FOLLOWING}`, `CURRENT ROW`, or `expr {PRECEDING|FOLLOWING}`.
func (p *Parser) parseWindowBound() (ast.WindowBound, error) {
	switch p.cur.Type {
	case kwUNBOUNDED:
		p.advance()
		switch p.cur.Type {
		case kwPRECEDING:
			p.advance()
			return ast.WindowBound{Kind: ast.BoundUnboundedPreceding}, nil
		case kwFOLLOWING:
			p.advance()
			return ast.WindowBound{Kind: ast.BoundUnboundedFollowing}, nil
		default:
			return ast.WindowBound{}, p.syntaxErrorAtCur()
		}
	case kwCURRENT:
		p.advance()
		if _, err := p.expect(kwROW); err != nil {
			return ast.WindowBound{}, err
		}
		return ast.WindowBound{Kind: ast.BoundCurrentRow}, nil
	default:
		expr, err := p.parseExpr()
		if err != nil {
			return ast.WindowBound{}, err
		}
		switch p.cur.Type {
		case kwPRECEDING:
			p.advance()
			return ast.WindowBound{Kind: ast.BoundPreceding, Offset: expr}, nil
		case kwFOLLOWING:
			p.advance()
			return ast.WindowBound{Kind: ast.BoundFollowing, Offset: expr}, nil
		default:
			return ast.WindowBound{}, p.syntaxErrorAtCur()
		}
	}
}

// ---------------------------------------------------------------------------
// Literals
// ---------------------------------------------------------------------------

// parseStringLiteral parses one or more adjacent STRING_LITERAL components
// (string_literal: adjacent concatenation). GoogleSQL concatenates adjacent
// string literals into one value; the lexer emits one token per component.
// Adjacent bytes/string mixing is not handled here (the lexer typed each token);
// a bytes token after a string simply stops the string and is left for the
// caller, which is a syntax error in expression position if unexpected.
func (p *Parser) parseStringLiteral() (ast.Node, error) {
	first := p.advance()
	value := first.Str
	loc := first.Loc
	for p.cur.Type == tokString {
		comp := p.advance()
		value += comp.Str
		loc.End = comp.Loc.End
	}
	return &ast.Literal{Kind: ast.LitString, Value: value, Loc: loc}, nil
}

// parseBytesLiteral parses one or more adjacent BYTES_LITERAL components
// (bytes_literal: adjacent concatenation).
func (p *Parser) parseBytesLiteral() (ast.Node, error) {
	first := p.advance()
	value := first.Str
	loc := first.Loc
	for p.cur.Type == tokBytes {
		comp := p.advance()
		value += comp.Str
		loc.End = comp.Loc.End
	}
	return &ast.Literal{Kind: ast.LitBytes, Value: value, Loc: loc}, nil
}

// parseTypedLiteral parses a type-prefixed literal `<kw> '…'` for
// NUMERIC/DECIMAL/BIGNUMERIC/BIGDECIMAL/JSON/DATE/TIME/DATETIME/TIMESTAMP. The
// keyword is the current token and is known to be followed by a string literal.
func (p *Parser) parseTypedLiteral() (ast.Node, error) {
	kw := p.advance()
	strTok := p.advance() // the string literal (peekNext-confirmed by caller)
	return &ast.TypedLiteral{
		TypeKeyword: TokenName(kw.Type),
		Value:       strTok.Str,
		Loc:         ast.Loc{Start: kw.Loc.Start, End: strTok.Loc.End},
	}, nil
}

// parseRangeLiteral parses `RANGE<type> '…'` (range_literal: range_type
// string_literal). RANGE is the current token, followed by '<'. The element
// type is parsed by the `types` node; the trailing string literal is required.
func (p *Parser) parseRangeLiteral() (ast.Node, error) {
	startLoc := p.cur.Loc
	// Reuse the type parser for `RANGE<type>` (parseRangeType lives in the types
	// node). It consumes RANGE and the template.
	rt, err := p.parseRangeType()
	if err != nil {
		return nil, err
	}
	strTok, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	return &ast.TypedLiteral{
		TypeKeyword: rt.String(),
		Value:       strTok.Str,
		Loc:         ast.Loc{Start: startLoc.Start, End: strTok.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Parameters / system variables
// ---------------------------------------------------------------------------

// parseNamedParameter parses `@identifier` (named_parameter_expression). The
// current token is '@'. (Statement-level hints `@{`/`@N`/`@[` are handled by the
// foundation before dispatch and never reach expression parsing; an '@' here is
// always a named parameter.)
func (p *Parser) parseNamedParameter() (ast.Node, error) {
	atTok := p.advance() // '@'
	if !isAnyKeywordIdentifier(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	return &ast.Parameter{Name: name, Loc: ast.Loc{Start: atTok.Loc.Start, End: nameTok.Loc.End}}, nil
}

// parseSystemVariable parses `@@path` (system_variable_expression: ATAT
// path_expression). The current token is '@@'.
func (p *Parser) parseSystemVariable() (ast.Node, error) {
	atatTok := p.advance() // '@@'
	if !isAnyKeywordIdentifier(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	sv := &ast.SystemVariable{Loc: atatTok.Loc}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	sv.Parts = append(sv.Parts, name)
	sv.Loc.End = nameTok.Loc.End
	for p.cur.Type == int('.') {
		if !isAnyKeywordIdentifier(p.peekNext().Type) {
			break
		}
		p.advance() // '.'
		partTok := p.advance()
		part, err := p.identifierText(partTok)
		if err != nil {
			return nil, err
		}
		sv.Parts = append(sv.Parts, part)
		sv.Loc.End = partTok.Loc.End
	}
	return sv, nil
}

// ---------------------------------------------------------------------------
// CASE / CAST / EXTRACT / INTERVAL
// ---------------------------------------------------------------------------

// parseCaseExpr parses a CASE expression (case_expression). The current token
// is CASE. Simple form has an operand before the first WHEN; searched form goes
// straight to WHEN. At least one WHEN…THEN is required; ELSE is optional; END
// terminates.
func (p *Parser) parseCaseExpr() (ast.Node, error) {
	caseTok := p.advance() // CASE
	ce := &ast.CaseExpr{Loc: caseTok.Loc}

	// Simple CASE: an operand expression before WHEN.
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Operand = operand
	}

	// One or more WHEN cond THEN result.
	for p.cur.Type == kwWHEN {
		whenTok := p.advance()
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
			Cond: cond, Result: result,
			Loc: ast.Loc{Start: whenTok.Loc.Start, End: nodeLoc(result).End},
		})
	}
	if len(ce.Whens) == 0 {
		// `CASE END` / `CASE x END` with no WHEN is a syntax error (oracle:
		// `CASE END` rejects "Unexpected keyword END").
		return nil, p.syntaxErrorAtCur()
	}

	// Optional ELSE.
	if p.cur.Type == kwELSE {
		p.advance()
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

// parseCastExpr parses `CAST(expr AS type [FORMAT fmt [AT TIME ZONE tz]])` or
// the SAFE_CAST variant (cast_expression). The CAST/SAFE_CAST keyword is the
// current token. The `CAST(CAST …` / `CAST(SAFE_CAST …` error alternatives are
// surfaced as the grammar's notify-style reject: the inner reserved-keyword
// argument is not a valid expression, so it falls out as a syntax error
// (oracle: `CAST(CAST AS INT64)` rejects).
func (p *Parser) parseCastExpr(safe bool) (ast.Node, error) {
	castTok := p.advance() // CAST / SAFE_CAST
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
	dt, err := p.parseType()
	if err != nil {
		return nil, err
	}
	ce := &ast.CastExpr{
		Expr: expr,
		Type: &ast.TypeRef{Text: dt.String(), Loc: dt.Loc},
		Safe: safe,
		Loc:  castTok.Loc,
	}
	// Optional FORMAT fmt [AT TIME ZONE tz].
	if p.cur.Type == kwFORMAT {
		p.advance()
		fmtExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Format = fmtExpr
		if p.cur.Type == int('@') && p.peekNext().Type == kwTIME {
			// AT TIME ZONE — '@' here is the AT keyword? No: AT is the '@' symbol?
			// GoogleSQL AT TIME ZONE uses the keyword AT, lexed as a word. Handled
			// below via the AT-token path.
		}
		if tz, ok, err := p.tryParseAtTimeZone(); err != nil {
			return nil, err
		} else if ok {
			ce.TimeZone = tz
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ce.Loc.End = closeTok.Loc.End
	return ce, nil
}

// parseExtractExpr parses `EXTRACT(part FROM expr [AT TIME ZONE tz])`
// (extract_expression). EXTRACT is the current token.
func (p *Parser) parseExtractExpr() (ast.Node, error) {
	extractTok := p.advance() // EXTRACT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	part, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ee := &ast.ExtractExpr{Part: part, From: from, Loc: extractTok.Loc}
	if tz, ok, err := p.tryParseAtTimeZone(); err != nil {
		return nil, err
	} else if ok {
		ee.TimeZone = tz
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ee.Loc.End = closeTok.Loc.End
	return ee, nil
}

// parseIntervalExpr parses `INTERVAL expr datepart [TO datepart]`
// (interval_expression). INTERVAL is the current token. The date-part(s) are
// identifier spellings (e.g. DAY, YEAR, SECOND).
func (p *Parser) parseIntervalExpr() (ast.Node, error) {
	intervalTok := p.advance() // INTERVAL
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !isAnyKeywordIdentifier(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	fromTok := p.advance()
	fromPart, err := p.identifierText(fromTok)
	if err != nil {
		return nil, err
	}
	ie := &ast.IntervalExpr{Value: val, From: fromPart, Loc: ast.Loc{Start: intervalTok.Loc.Start, End: fromTok.Loc.End}}
	if p.cur.Type == kwTO {
		p.advance()
		if !isAnyKeywordIdentifier(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		toTok := p.advance()
		toPart, err := p.identifierText(toTok)
		if err != nil {
			return nil, err
		}
		ie.To = toPart
		ie.Loc.End = toTok.Loc.End
	}
	return ie, nil
}

// tryParseAtTimeZone parses an optional `AT TIME ZONE expr` suffix
// (opt_at_time_zone). It returns (tz, true, nil) when present, (nil, false, nil)
// when the current token is not the AT keyword. GoogleSQL's AT is a word-keyword
// — but the lexer has no dedicated AT token; the foundation's keyword set omits
// a standalone "at" keyword (the '@' AT_SYMBOL is punctuation). AT TIME ZONE is
// recognized by matching a bare identifier "AT" followed by TIME ZONE.
func (p *Parser) tryParseAtTimeZone() (ast.Node, bool, error) {
	if !p.curIsWord("AT") {
		return nil, false, nil
	}
	// Look ahead: AT TIME ZONE. TIME is a keyword; ZONE is a keyword.
	if p.peekNext().Type != kwTIME {
		return nil, false, nil
	}
	p.advance() // AT
	p.advance() // TIME
	if _, err := p.expect(kwZONE); err != nil {
		return nil, false, err
	}
	tz, err := p.parseExpr()
	if err != nil {
		return nil, false, err
	}
	return tz, true, nil
}

// curIsWord reports whether the current token is an identifier (or keyword)
// whose spelling equals word (case-insensitive). Used for context keywords that
// the lexer does not tokenize as dedicated keywords (e.g. AT in AT TIME ZONE).
func (p *Parser) curIsWord(word string) bool {
	if p.cur.Type == tokIdentifier {
		return strings.EqualFold(p.cur.Str, word)
	}
	if p.cur.Type >= keywordBase {
		return strings.EqualFold(TokenName(p.cur.Type), word)
	}
	return false
}

// ---------------------------------------------------------------------------
// Constructors: array / struct / new / braced / replace-fields / with
// ---------------------------------------------------------------------------

// parseArrayOrArraySubquery parses an ARRAY-led primary: `ARRAY(query)`
// (expression_subquery_with_keyword), `ARRAY[…]` (array_constructor with the
// ARRAY keyword), or `ARRAY<T>[…]` (typed array constructor). ARRAY is the
// current token.
func (p *Parser) parseArrayOrArraySubquery() (ast.Node, error) {
	arrayTok := p.advance() // ARRAY
	switch p.cur.Type {
	case int('('):
		// ARRAY ( query ).
		p.advance() // '('
		sub, end, err := p.parseSubqueryBodyRaw(arrayTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		return &ast.ArraySubqueryExpr{RawText: sub, Loc: ast.Loc{Start: arrayTok.Loc.Start, End: end}}, nil
	case int('['):
		// ARRAY [ … ].
		return p.parseArrayConstructor(true, nil)
	case int('<'):
		// ARRAY < T > [ … ] — typed array constructor. Re-parse the full type
		// starting from ARRAY using the types node.
		dt, err := p.parseArrayTypeFromKeyword(arrayTok)
		if err != nil {
			return nil, err
		}
		tr := &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
		return p.parseArrayConstructor(true, tr)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseArrayConstructor parses `[ … ]` element list (array_constructor). The
// current token is '['. hasArrayKeyword/elemType carry context from a leading
// ARRAY / ARRAY<T>.
func (p *Parser) parseArrayConstructor(hasArrayKeyword bool, elemType *ast.TypeRef) (ast.Node, error) {
	openTok := p.advance() // '['
	ae := &ast.ArrayExpr{HasArrayKeyword: hasArrayKeyword, ElemType: elemType}
	start := openTok.Loc.Start
	if elemType != nil {
		start = elemType.Loc.Start
	} else if hasArrayKeyword {
		start = p.prev.Loc.Start // ARRAY keyword captured by caller; approximate
	}
	if p.cur.Type != int(']') {
		elems, err := p.parseExprCommaList()
		if err != nil {
			return nil, err
		}
		ae.Elements = elems
	}
	closeTok, err := p.expect(int(']'))
	if err != nil {
		return nil, err
	}
	ae.Loc = ast.Loc{Start: start, End: closeTok.Loc.End}
	return ae, nil
}

// parseStructOrTyped parses a STRUCT-led primary: `STRUCT(args)`,
// `STRUCT<…>(args)`, or `STRUCT { … }` / `STRUCT<…> { … }` braced. STRUCT is the
// current token.
func (p *Parser) parseStructOrTyped() (ast.Node, error) {
	switch p.peekNext().Type {
	case int('('):
		// STRUCT ( args ).
		structTok := p.advance() // STRUCT
		p.advance()              // '('
		return p.parseStructBody(structTok.Loc.Start, true, nil)
	case int('{'):
		// STRUCT { … } braced (no type).
		p.advance() // STRUCT
		return p.parseBracedConstructor(nil)
	case int('<'):
		// STRUCT < … > followed by '(' (typed constructor) or '{' (braced). Use
		// parseRawType (NOT parseType): the `(args)` that follows is the
		// constructor's argument list, NOT an opt_type_parameters list — parseType
		// would greedily consume `(1)` as type parameters. parseRawType stops after
		// `STRUCT<…>`. (A field's own type params, e.g. STRUCT<x STRING(10)>, are
		// parsed inside the template by the field's parseType and are unaffected.)
		dt, err := p.parseRawType() // consumes STRUCT<…> only
		if err != nil {
			return nil, err
		}
		tr := &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
		switch p.cur.Type {
		case int('('):
			p.advance() // '('
			return p.parseStructBody(dt.Loc.Start, true, tr)
		case int('{'):
			return p.parseBracedConstructor(tr)
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseStructBody parses the parenthesized struct-constructor argument list
// (struct_constructor_prefix*). The current token is just after '('. start
// anchors Loc; hasStruct/typeRef carry context.
func (p *Parser) parseStructBody(start int, hasStruct bool, typeRef *ast.TypeRef) (ast.Node, error) {
	se := &ast.StructExpr{HasStruct: hasStruct, Type: typeRef, Loc: ast.Loc{Start: start}}
	if p.cur.Type != int(')') {
		for {
			field, err := p.parseStructArg()
			if err != nil {
				return nil, err
			}
			se.Fields = append(se.Fields, field)
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	se.Loc.End = closeTok.Loc.End
	return se, nil
}

// parseStructArg parses one `expr [AS alias]` struct constructor argument
// (struct_constructor_arg).
func (p *Parser) parseStructArg() (*ast.StructFieldExpr, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	field := &ast.StructFieldExpr{Value: expr, Loc: nodeLoc(expr)}
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		field.Alias = alias
		field.Loc.End = aliasTok.Loc.End
	}
	return field, nil
}

// parseNewConstructor parses `NEW type ( args )` (new_constructor) or
// `NEW type { … }` (braced_new_constructor). NEW is the current token. The type
// is a type_name (parsed via the types node, which yields a path/INTERVAL form).
func (p *Parser) parseNewConstructor() (ast.Node, error) {
	newTok := p.advance() // NEW
	dt, err := p.parseTypeName()
	if err != nil {
		return nil, err
	}
	tr := &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
	nc := &ast.NewConstructor{Type: tr, Loc: newTok.Loc}
	switch p.cur.Type {
	case int('('):
		p.advance() // '('
		if p.cur.Type != int(')') {
			for {
				arg, err := p.parseNewArg()
				if err != nil {
					return nil, err
				}
				nc.Args = append(nc.Args, arg)
				if p.cur.Type != int(',') {
					break
				}
				p.advance()
			}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		nc.Loc.End = closeTok.Loc.End
		return nc, nil
	case int('{'):
		braced, err := p.parseBracedConstructor(nil)
		if err != nil {
			return nil, err
		}
		nc.Braced = braced.(*ast.BracedConstructor)
		nc.Loc.End = nc.Braced.Loc.End
		return nc, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseNewArg parses one new_constructor_arg: `expr`, `expr AS id`, or
// `expr AS ( path )`.
func (p *Parser) parseNewArg() (*ast.StructFieldExpr, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	field := &ast.StructFieldExpr{Value: expr, Loc: nodeLoc(expr)}
	if p.cur.Type == kwAS {
		p.advance()
		if p.cur.Type == int('(') {
			p.advance()
			path, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			field.Alias = "(" + path.String() + ")"
			field.Loc.End = closeTok.Loc.End
			return field, nil
		}
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		field.Alias = alias
		field.Loc.End = aliasTok.Loc.End
	}
	return field, nil
}

// parseBracedConstructor parses a proto braced constructor `{ field… }`
// (braced_constructor). The current token is '{'. typeRef carries an optional
// leading struct type (struct_braced_constructor). Field internals are captured
// best-effort: each field is `path { value | : expr }` or an extension
// `( path )`; the value expressions are retained for walking. A trailing comma
// is permitted.
func (p *Parser) parseBracedConstructor(typeRef *ast.TypeRef) (ast.Node, error) {
	openTok := p.advance() // '{'
	bc := &ast.BracedConstructor{Type: typeRef, Loc: openTok.Loc}
	if typeRef != nil {
		bc.Loc.Start = typeRef.Loc.Start
	}
	for p.cur.Type != int('}') && p.cur.Type != tokEOF {
		// Extension field: ( path ).
		if p.cur.Type == int('(') {
			p.advance()
			path, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			bc.Fields = append(bc.Fields, path)
		} else {
			// path { … } | path : expr.
			lhs, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			bc.Fields = append(bc.Fields, lhs)
			switch p.cur.Type {
			case int(':'):
				p.advance()
				val, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				bc.Fields = append(bc.Fields, val)
			case int('{'):
				nested, err := p.parseBracedConstructor(nil)
				if err != nil {
					return nil, err
				}
				bc.Fields = append(bc.Fields, nested)
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}
		// Optional comma between fields (also a trailing comma before '}').
		if p.cur.Type == int(',') {
			p.advance()
		}
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	bc.Loc.End = closeTok.Loc.End
	return bc, nil
}

// parseReplaceFields parses `REPLACE_FIELDS(expr, value AS path, …)`
// (replace_fields_expression). REPLACE_FIELDS is the current token.
func (p *Parser) parseReplaceFields() (ast.Node, error) {
	rfTok := p.advance() // REPLACE_FIELDS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	base, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	rf := &ast.ReplaceFieldsExpr{Expr: base, Loc: rfTok.Loc}
	for p.cur.Type == int(',') {
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		path, err := p.parseGeneralizedPath()
		if err != nil {
			return nil, err
		}
		rf.Items = append(rf.Items, &ast.StructFieldExpr{Value: val, Alias: path, Loc: nodeLoc(val)})
	}
	if len(rf.Items) == 0 {
		// REPLACE_FIELDS requires at least one replacement arg.
		return nil, p.syntaxErrorAtCur()
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	rf.Loc.End = closeTok.Loc.End
	return rf, nil
}

// parseWithExpr parses the inline `WITH(name AS expr, …, body)` expression
// (with_expression). WITH is the current token, known (by the caller) to be
// followed by '('. Each leading entry is `name AS expr`; the final entry is the
// body expression.
func (p *Parser) parseWithExpr() (ast.Node, error) {
	// Distinguish the inline WITH(...) expression from a WITH-CTE query head: an
	// inline WITH expression requires `WITH (` and a first entry `id AS expr`.
	// (A query `WITH cte AS (...)` is not an expression and never reaches here in
	// expression position; the expression node only owns WITH(...) here.)
	if p.peekNext().Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	withTok := p.advance() // WITH
	p.advance()            // '('
	we := &ast.WithExpr{Loc: withTok.Loc}
	for {
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		bound, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		we.Vars = append(we.Vars, &ast.StructFieldExpr{Alias: name, Value: bound, Loc: ast.Loc{Start: nameTok.Loc.Start, End: nodeLoc(bound).End}})
		if _, err := p.expect(int(',')); err != nil {
			return nil, err
		}
		// After a variable, the next token decides: another `id AS …` binding, or
		// the body expression. A binding is `identifier AS`; anything else is the
		// body.
		if isIdentifierStart(p.cur.Type) && p.peekNext().Type == kwAS {
			continue
		}
		break
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	we.Body = body
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	we.Loc.End = closeTok.Loc.End
	return we, nil
}

// ---------------------------------------------------------------------------
// Parenthesized expression / subquery / EXISTS / UNNEST
// ---------------------------------------------------------------------------

// parseParenOrSubquery parses a parenthesized primary: a parenthesized query
// `( SELECT … )` → SubqueryExpr (parenthesized_query), a bare struct-tuple
// `( e1, e2, … )` with two or more elements → StructExpr (the
// struct_constructor_prefix_without_keyword form), or a plain parenthesized
// expression `( expr )` → ParenExpr. The current token is '('.
func (p *Parser) parseParenOrSubquery() (ast.Node, error) {
	openTok := p.advance() // '('

	// Parenthesized query.
	if p.atQueryStart() {
		sub, _, err := p.parseSubqueryBody(openTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		return sub, nil
	}

	// First element expression.
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	// Bare struct tuple: two or more comma-separated elements.
	if p.cur.Type == int(',') {
		se := &ast.StructExpr{HasStruct: false, Loc: openTok.Loc}
		se.Fields = append(se.Fields, &ast.StructFieldExpr{Value: first, Loc: nodeLoc(first)})
		for p.cur.Type == int(',') {
			p.advance()
			field, err := p.parseStructArg()
			if err != nil {
				return nil, err
			}
			se.Fields = append(se.Fields, field)
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		se.Loc.End = closeTok.Loc.End
		return se, nil
	}
	// Plain parenthesized expression.
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.ParenExpr{Expr: first, Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}, nil
}

// parseExistsExpr parses `EXISTS [hint] ( query )`
// (expression_subquery_with_keyword). EXISTS is the current token.
func (p *Parser) parseExistsExpr() (ast.Node, error) {
	existsTok := p.advance() // EXISTS
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	raw, end, err := p.parseSubqueryBodyRaw(existsTok.Loc.Start)
	if err != nil {
		return nil, err
	}
	return &ast.ExistsExpr{RawText: raw, Loc: ast.Loc{Start: existsTok.Loc.Start, End: end}}, nil
}

// parseUnnestExpression parses `UNNEST ( expr [AS alias] [, zip-mode] )`
// (unnest_expression) as an expression. UNNEST is the current token. It is
// represented as a FuncCall named UNNEST so downstream consumers treat it
// uniformly; the inner args carry the array expression(s).
func (p *Parser) parseUnnestExpression() (ast.Node, error) {
	unnestTok := p.advance() // UNNEST
	name := &ast.PathExpr{Parts: []string{"UNNEST"}, Loc: unnestTok.Loc}
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseFuncCallSuffix(name)
}

// atQueryStart reports whether the current token begins a query (so a following
// '(' opens a parenthesized_query rather than a parenthesized expression). A
// query head is SELECT, WITH, FROM (the FROM-query form), GRAPH (GQL), or a
// nested '(' that itself opens a query. This is the disambiguator the grammar's
// LALR(1) lookahead applies.
func (p *Parser) atQueryStart() bool {
	switch p.cur.Type {
	case kwSELECT, kwWITH, kwFROM, kwGRAPH:
		return true
	case int('('):
		// A nested '(' could open a parenthesized query — but it could also open a
		// parenthesized expression or a struct tuple. Peek one token: a SELECT/WITH
		// after the nested '(' indicates a query.
		switch p.peekNext().Type {
		case kwSELECT, kwWITH, kwFROM:
			return true
		}
		return false
	}
	return false
}

// parseSubqueryBody parses a parenthesized query body after the opening '(' has
// been consumed (the current token is the query's first token, e.g. SELECT). It
// captures the inner query as a SubqueryExpr with RawText set (the query grammar
// belongs to parser-select, which fills Query later). Returns the node and the
// end offset (past the matching ')').
func (p *Parser) parseSubqueryBody(start int) (ast.Node, int, error) {
	raw, end, err := p.parseSubqueryBodyRaw(start)
	if err != nil {
		return nil, 0, err
	}
	return &ast.SubqueryExpr{RawText: raw, Loc: ast.Loc{Start: start, End: end}}, end, nil
}

// parseSubqueryBodyRaw consumes a balanced parenthesized query body — the
// current token is the first token inside the already-consumed '(' — up to and
// including the matching ')'. It returns the inner source text (between the
// parens, trimmed) and the end offset past the ')'. The inner query is NOT
// parsed: the query/SELECT grammar is owned by the downstream parser-select
// node (which depends on this one), so the expression node only records the
// span so it accepts every subquery-bearing form the oracle accepts and so
// parser-select can re-parse RawText when it lands.
//
// Balance tracks ( ) [ ] { } so a nested subquery, array, or struct inside the
// query does not prematurely close it. An unterminated body (EOF before the
// matching ')') is a syntax error.
func (p *Parser) parseSubqueryBodyRaw(start int) (string, int, error) {
	innerStart := p.cur.Loc.Start
	depth := 1
	innerEnd := innerStart
	for {
		if p.cur.Type == tokEOF {
			return "", 0, &ParseError{Loc: p.cur.Loc, Msg: "syntax error: unterminated subquery (expected ')')"}
		}
		switch p.cur.Type {
		case int('('), int('['), int('{'):
			depth++
		case int(')'), int(']'), int('}'):
			depth--
			if depth == 0 {
				closeTok := p.advance() // consume the matching ')'
				inner := strings.TrimSpace(p.input[absIndex(p, innerStart):absIndex(p, innerEnd)])
				return inner, closeTok.Loc.End, nil
			}
		}
		innerEnd = p.cur.Loc.End
		p.advance()
	}
}

// absIndex converts an absolute source offset to an index into p.input, which
// is the current segment text. The lexer was created with NewLexerWithOffset, so
// token offsets are absolute (baseOffset + segment index); subtract baseOffset
// to index p.input. For the standalone ParseExpression path baseOffset is 0.
func absIndex(p *Parser, absOffset int) int {
	idx := absOffset - p.baseOffset
	if idx < 0 {
		idx = 0
	}
	if idx > len(p.input) {
		idx = len(p.input)
	}
	return idx
}

// ---------------------------------------------------------------------------
// Lambdas, generalized paths, list helpers, hint/options skips
// ---------------------------------------------------------------------------

// tryParseParenLambda attempts to parse a parenthesized-parameter lambda
// `( id (, id)* ) -> body` or the empty `() -> body` (lambda_argument with a
// parenthesized list). The current token is '('. It returns (node, true, nil) on
// success. If the parenthesized run is NOT a lambda (no trailing '->'), it does
// NOT consume anything and returns (nil, false, nil) so the caller falls back to
// parsing a parenthesized expression. Detection scans the balanced '(' … ')'
// span and checks the token immediately after it for '->'.
func (p *Parser) tryParseParenLambda() (ast.Node, bool, error) {
	if !p.parenRunFollowedByArrow() {
		return nil, false, nil
	}
	openTok := p.advance() // '('
	var params []string
	if p.cur.Type != int(')') {
		for {
			idTok, err := p.expectIdentifier()
			if err != nil {
				return nil, false, err
			}
			name, err := p.identifierText(idTok)
			if err != nil {
				return nil, false, err
			}
			params = append(params, name)
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}
	if _, err := p.expect(tokArrow); err != nil {
		return nil, false, err
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, false, err
	}
	return &ast.LambdaExpr{Params: params, Body: body, Loc: ast.Loc{Start: openTok.Loc.Start, End: nodeLoc(body).End}}, true, nil
}

// parenRunFollowedByArrow reports whether the balanced parenthesized run
// starting at the current '(' is immediately followed by a '->' token (the
// signature of a parenthesized-list lambda). It uses a temporary lexer scan from
// the saved cursor WITHOUT mutating the parser state. The parser's two-token
// lookahead is insufficient (the param list is arbitrary-length), so this peeks
// via a throwaway lexer cloned at the current position.
func (p *Parser) parenRunFollowedByArrow() bool {
	// Scan the remaining input from the current token's start with a fresh lexer.
	startIdx := absIndex(p, p.cur.Loc.Start)
	lx := NewLexerWithOffset(p.input[startIdx:], 0)
	depth := 0
	for {
		t := lx.NextToken()
		if t.Type == tokEOF {
			return false
		}
		switch t.Type {
		case int('('), int('['), int('{'):
			depth++
		case int(')'), int(']'), int('}'):
			depth--
			if depth == 0 {
				// The run closed; the next token decides.
				return lx.NextToken().Type == tokArrow
			}
		}
	}
}

// captureBareSelectArg consumes a bare SELECT used as a function argument (the
// grammar's notify-style error alternative) up to — but not including — the
// argument list's closing ')' or a top-level ',' at bracket depth 0. It returns
// a SubqueryExpr capturing the consumed source so the call still parses (the
// diagnostic is recorded by the caller). The current token is SELECT.
func (p *Parser) captureBareSelectArg(start ast.Loc) ast.Node {
	innerStart := p.cur.Loc.Start
	innerEnd := innerStart
	depth := 0
	for p.cur.Type != tokEOF {
		if depth == 0 && (p.cur.Type == int(')') || p.cur.Type == int(',')) {
			break
		}
		switch p.cur.Type {
		case int('('), int('['), int('{'):
			depth++
		case int(')'), int(']'), int('}'):
			depth--
		}
		innerEnd = p.cur.Loc.End
		p.advance()
	}
	raw := strings.TrimSpace(p.input[absIndex(p, innerStart):absIndex(p, innerEnd)])
	return &ast.SubqueryExpr{RawText: raw, Loc: ast.Loc{Start: start.Start, End: innerEnd}}
}

// parseGeneralizedPath parses a generalized_path_expression and renders it to a
// string (used for REPLACE_FIELDS replacement targets, which carry a path, not a
// general expression). It accepts `identifier`, `. identifier`, `. ( path )`,
// and `[ expr ]` continuations. The rendering is for the StructFieldExpr.Alias
// slot; the exact path structure is not consumed by bytebase.
func (p *Parser) parseGeneralizedPath() (string, error) {
	if !isIdentifierStart(p.cur.Type) {
		return "", p.syntaxErrorAtCur()
	}
	tok := p.advance()
	name, err := p.identifierText(tok)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString(name)
	for {
		switch p.cur.Type {
		case int('.'):
			p.advance()
			if p.cur.Type == int('(') {
				p.advance()
				path, err := p.parsePathExpr()
				if err != nil {
					return "", err
				}
				if _, err := p.expect(int(')')); err != nil {
					return "", err
				}
				b.WriteString(".(")
				b.WriteString(path.String())
				b.WriteString(")")
				continue
			}
			if !isAnyKeywordIdentifier(p.cur.Type) {
				return "", p.syntaxErrorAtCur()
			}
			part := p.advance()
			pn, err := p.identifierText(part)
			if err != nil {
				return "", err
			}
			b.WriteString(".")
			b.WriteString(pn)
		case int('['):
			p.advance()
			if _, err := p.parseExpr(); err != nil {
				return "", err
			}
			if _, err := p.expect(int(']')); err != nil {
				return "", err
			}
			b.WriteString("[...]")
		default:
			return b.String(), nil
		}
	}
}

// parseExprCommaList parses one or more comma-separated expressions (no
// surrounding delimiters consumed). The current token begins the first
// expression; parsing stops at the first non-comma boundary.
func (p *Parser) parseExprCommaList() ([]ast.Node, error) {
	var out []ast.Node
	for {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		out = append(out, e)
		if p.cur.Type != int(',') {
			return out, nil
		}
		p.advance()
	}
}

// parseExprListThroughParen parses a parenthesized expression list when the
// opening '(' has already been consumed: `expr (, expr)* )`. At least one
// expression is required (the IN / quantified-list RHS, which the grammar
// rejects when empty — oracle: `a IN ()` rejects). Returns the elements and the
// Loc through the closing ')'.
func (p *Parser) parseExprListThroughParen() ([]ast.Node, ast.Loc, error) {
	values, err := p.parseExprCommaList()
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	return values, closeTok.Loc, nil
}

// parseArrayTypeFromKeyword parses `ARRAY < type >` given the already-consumed
// ARRAY keyword token (arrTok). It is a thin wrapper that reconstructs the
// element-type parse so parseArrayType (which itself consumes ARRAY) is not
// re-entered with the keyword already gone. The current token is '<'.
func (p *Parser) parseArrayTypeFromKeyword(arrTok Token) (*DataType, error) {
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:        TypeArray,
		ElementType: elem,
		Loc:         ast.Loc{Start: arrTok.Loc.Start, End: endLoc.End},
	}, nil
}

// expectIdentifier consumes the current token if it can be an identifier
// (token identifier or non-reserved keyword), else returns a syntax error.
func (p *Parser) expectIdentifier() (Token, error) {
	if isIdentifierStart(p.cur.Type) {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// skipHint consumes a hint `@<int>` or `@{ … }` / `@[ <int> @] { … }` that may
// appear before an IN/LIKE RHS or after a function call (hint). The current
// token is '@'. It reuses the foundation's brace-balanced hint skipper. Returns
// a *ParseError for a malformed (unterminated / empty) hint body.
func (p *Parser) skipHint() *ParseError {
	if p.cur.Type != int('@') {
		return nil
	}
	hintStart := p.cur.Loc
	next := p.peekNext()
	switch {
	case next.Type == int('{'):
		p.advance() // '@'
		return p.skipBalancedBraces(hintStart)
	case next.Type == int('['):
		p.advance() // '@'
		p.advance() // '['
		for p.cur.Type != tokEOF && p.cur.Type != int(']') {
			p.advance()
		}
		if p.cur.Type != int(']') {
			return &ParseError{Loc: hintStart, Msg: "unterminated hint"}
		}
		p.advance() // ']'
		if p.cur.Type == int('{') {
			return p.skipBalancedBraces(hintStart)
		}
		return &ParseError{Loc: hintStart, Msg: "unterminated hint"}
	case next.Type == tokInteger:
		p.advance() // '@'
		p.advance() // int
	}
	return nil
}

// skipOptionsList consumes an options_list `( [entry (, entry)*] )` (used by
// WITH REPORT). The current token is '('. Entry internals are not retained;
// balance is tracked so nested parens do not close early.
func (p *Parser) skipOptionsList() error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	depth := 1
	for depth > 0 {
		if p.cur.Type == tokEOF {
			return &ParseError{Loc: p.cur.Loc, Msg: "syntax error: unterminated OPTIONS list"}
		}
		switch p.cur.Type {
		case int('('):
			depth++
		case int(')'):
			depth--
		}
		p.advance()
	}
	return nil
}

// skipGroupByTail consumes the trailing group-by of the HAVING MIN aggregate
// modifier (group_by_clause_prefix): `GROUP [hint] [AND ORDER] BY grouping_item
// (, grouping_item)*`. The current token is GROUP. Items are not retained;
// each grouping_item is parsed as an expression (the common case) up to the
// boundary. Used only for parse parity of the DP HAVING MIN form.
func (p *Parser) skipGroupByTail() error {
	p.advance() // GROUP
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return herr
		}
	}
	if p.cur.Type == kwAND {
		p.advance()
		if _, err := p.expect(kwORDER); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	for {
		if _, err := p.parseExpr(); err != nil {
			return err
		}
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	return nil
}

// ---------------------------------------------------------------------------
// Loc helpers (local to the expressions node)
// ---------------------------------------------------------------------------
//
// These compute spans directly from already-built nodes / tokens, without
// relying on ast.NodeLoc (which is owned by ast-core and does not yet enumerate
// the expression node types). nodeLoc reads the Loc field off any expression
// node this package builds.

// spanNodes returns the Loc spanning from a's start to b's end (both expression
// nodes built by this package).
func spanNodes(a, b ast.Node) ast.Loc {
	return ast.Loc{Start: nodeLoc(a).Start, End: nodeLoc(b).End}
}

// locFromTo returns the Loc spanning from node a's start to token b's end.
func locFromTo(a ast.Node, b Token) ast.Loc {
	return ast.Loc{Start: nodeLoc(a).Start, End: b.Loc.End}
}

// locFromTo2 returns the Loc spanning from loc a's start to node b's end.
func locFromTo2(a ast.Loc, b ast.Node) ast.Loc {
	return ast.Loc{Start: a.Start, End: nodeLoc(b).End}
}

// nodeLoc reads the Loc off any expression node this package constructs. It is a
// local type switch (not ast.NodeLoc) so the expressions node owns its own span
// extraction independent of ast-core's NodeLoc enumeration. An unknown node
// yields NoLoc (defensive; every node this package builds is enumerated).
func nodeLoc(n ast.Node) ast.Loc {
	switch v := n.(type) {
	case *ast.Identifier:
		return v.Loc
	case *ast.PathExpr:
		return v.Loc
	case *ast.TypeRef:
		return v.Loc
	case *ast.StarExpr:
		return v.Loc
	case *ast.Literal:
		return v.Loc
	case *ast.TypedLiteral:
		return v.Loc
	case *ast.IntervalExpr:
		return v.Loc
	case *ast.Parameter:
		return v.Loc
	case *ast.SystemVariable:
		return v.Loc
	case *ast.ParenExpr:
		return v.Loc
	case *ast.UnaryExpr:
		return v.Loc
	case *ast.BinaryExpr:
		return v.Loc
	case *ast.CompareExpr:
		return v.Loc
	case *ast.IsExpr:
		return v.Loc
	case *ast.InExpr:
		return v.Loc
	case *ast.BetweenExpr:
		return v.Loc
	case *ast.LikeExpr:
		return v.Loc
	case *ast.CaseExpr:
		return v.Loc
	case *ast.CastExpr:
		return v.Loc
	case *ast.ExtractExpr:
		return v.Loc
	case *ast.FuncCall:
		return v.Loc
	case *ast.NamedArg:
		return v.Loc
	case *ast.LambdaExpr:
		return v.Loc
	case *ast.SequenceArg:
		return v.Loc
	case *ast.ArrayExpr:
		return v.Loc
	case *ast.StructExpr:
		return v.Loc
	case *ast.StructFieldExpr:
		return v.Loc
	case *ast.NewConstructor:
		return v.Loc
	case *ast.BracedConstructor:
		return v.Loc
	case *ast.ReplaceFieldsExpr:
		return v.Loc
	case *ast.WithExpr:
		return v.Loc
	case *ast.FieldAccess:
		return v.Loc
	case *ast.IndexAccess:
		return v.Loc
	case *ast.ExtensionAccess:
		return v.Loc
	case *ast.SubqueryExpr:
		return v.Loc
	case *ast.ExistsExpr:
		return v.Loc
	case *ast.ArraySubqueryExpr:
		return v.Loc
	default:
		return ast.NoLoc()
	}
}
