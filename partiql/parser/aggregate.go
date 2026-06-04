// Package parser — aggregate function-call parsing (DAG node 14,
// parser-aggregates).
//
// This file implements the `aggregate` grammar rule (PartiQLParser.g4:577-580)
// for COUNT(*), COUNT([DISTINCT|ALL] expr), and
// SUM/AVG/MIN/MAX([DISTINCT|ALL] expr). It is reached from
// parsePrimaryBase's single dispatch case for
// tokCOUNT/tokSUM/tokAVG/tokMIN/tokMAX (exprprimary.go).
package parser

import (
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// parseAggregateExpr parses an aggregate function call. The current token
// must be one of COUNT / SUM / AVG / MIN / MAX (the five reserved aggregate
// keyword tokens).
//
// Grammar (PartiQLParser.g4:577-580):
//
//	aggregate
//	  : func=COUNT PAREN_LEFT ASTERISK PAREN_RIGHT                                      # CountAll
//	  | func=(COUNT|MAX|MIN|SUM|AVG) PAREN_LEFT setQuantifierStrategy? expr PAREN_RIGHT # AggregateBase
//
// where setQuantifierStrategy is `DISTINCT | ALL` (g4:229-232).
//
// COUNT is special: in addition to the two `aggregate` alternatives above it
// ALSO appears in the FunctionCallReserved list (g4:613), and in exprPrimary
// the `aggregate` alternative is tried before `functionCall` (g4:524 vs 526).
// ANTLR therefore parses COUNT through `aggregate` whenever it can (the star
// form, or `[DISTINCT|ALL] expr`) and falls through to `functionCall` only for
// the residual forms the aggregate rule cannot match — the empty-argument
// `COUNT()` and the multi-argument `COUNT(a, b, …)`. SUM/AVG/MIN/MAX are NOT in
// FunctionCallReserved, so they have no such fallback: their only legal shape
// is `( [DISTINCT|ALL] expr )`.
//
// The resulting forms and verdicts (all confirmed against the legacy ANTLR
// parser as the differential oracle, truth2):
//
//	COUNT(*)              -> {Name:COUNT Star:true}           (CountAll)
//	COUNT([DISTINCT|ALL] expr) -> {Name:COUNT Quantifier:… Args:[expr]} (AggregateBase)
//	COUNT(expr)           -> {Name:COUNT Args:[expr]}         (AggregateBase, 1 arg)
//	COUNT()               -> {Name:COUNT}                     (functionCall fallback)
//	COUNT(a, b, …)        -> {Name:COUNT Args:[a b …]}        (functionCall fallback)
//	SUM|AVG|MIN|MAX([DISTINCT|ALL] expr) -> {Name:… Quantifier:… Args:[expr]}
//
//	COUNT(DISTINCT *) / COUNT(ALL *)   REJECT (star carries no quantifier)
//	COUNT(*, …)                        REJECT (star is a lone argument)
//	COUNT(DISTINCT) / COUNT(ALL)       REJECT (quantifier needs an expr)
//	COUNT(DISTINCT a, b)               REJECT (quantifier form is single-expr)
//	SUM(*) / SUM() / SUM(a, b) / SUM(DISTINCT) … REJECT (no fallback for SUM/AVG/MIN/MAX)
//
// The AST node is the generic *ast.FuncCall, distinguished by its Star and
// Quantifier fields (ast/exprs.go).
func (p *Parser) parseAggregateExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	name := strings.ToUpper(p.cur.Str)
	isCount := p.cur.Type == tokCOUNT
	p.advance() // consume the aggregate keyword

	// COUNT fallback: when the token after `(` is neither the star nor a
	// DISTINCT/ALL quantifier, COUNT cannot match either `aggregate`
	// alternative and ANTLR falls through to FunctionCallReserved (g4:613
	// lists COUNT). That path is the ordinary `name(arg, …)` form, so we
	// reuse parseFuncCallArgs verbatim — `(` still current — which permits
	// zero args (`COUNT()`) and multiple args (`COUNT(a, b)`). A single arg
	// produces the same AST as AggregateBase's one-expr form, so the two
	// readings coincide. SUM/AVG/MIN/MAX have no such fallback (absent from
	// FunctionCallReserved); they always take the aggregate path below and so
	// reject empty/multi-arg/star forms.
	if isCount && p.cur.Type == tokPAREN_LEFT {
		next := p.peekNext().Type
		if next != tokASTERISK && next != tokDISTINCT && next != tokALL {
			args, endOff, err := p.parseFuncCallArgs()
			if err != nil {
				return nil, err
			}
			return &ast.FuncCall{
				Name: name,
				Args: args,
				Loc:  ast.Loc{Start: start, End: endOff},
			}, nil
		}
	}

	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}

	// CountAll: COUNT ( * ). Only COUNT admits the star, and the star is a
	// complete argument — no quantifier may precede it and nothing may follow
	// it but the close paren. SUM/AVG/MIN/MAX(*) is rejected because the star
	// is not a valid expression in the AggregateBase `expr` position (the
	// non-COUNT path below calls parseExprTop, which rejects `*`).
	if isCount && p.cur.Type == tokASTERISK {
		p.advance() // consume *
		endOff := p.cur.Loc.End
		if _, err := p.expect(tokPAREN_RIGHT); err != nil {
			return nil, err
		}
		return &ast.FuncCall{
			Name: name,
			Star: true,
			Loc:  ast.Loc{Start: start, End: endOff},
		}, nil
	}

	// Optional setQuantifierStrategy (DISTINCT|ALL). The AggregateBase rule
	// then requires EXACTLY ONE expr (no star, no comma-list) — for COUNT as
	// well as SUM/AVG/MIN/MAX. When absent (only reachable for SUM/AVG/MIN/MAX
	// here, since the no-quantifier COUNT cases were routed to the fallback
	// above), the same single-expr requirement applies.
	quantifier := ast.QuantifierNone
	switch p.cur.Type {
	case tokDISTINCT:
		quantifier = ast.QuantifierDistinct
		p.advance()
	case tokALL:
		quantifier = ast.QuantifierAll
		p.advance()
	}

	// `[DISTINCT|ALL]? expr` — exactly one expression, then the close paren.
	// A trailing comma (multi-arg) or a bare `)` / `*` (no valid expr) both
	// reject here: parseExprTop rejects `)`/`*`, and the PAREN_RIGHT
	// expectation rejects a trailing `, …`. This single tail serves the
	// quantified COUNT forms and every SUM/AVG/MIN/MAX form.
	arg, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	endOff := p.cur.Loc.End
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, err
	}
	return &ast.FuncCall{
		Name:       name,
		Quantifier: quantifier,
		Args:       []ast.ExprNode{arg},
		Loc:        ast.Loc{Start: start, End: endOff},
	}, nil
}
