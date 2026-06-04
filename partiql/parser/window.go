// Package parser — window function-call parsing (DAG node parser-window).
//
// This file implements the `windowFunction` grammar rule
// (PartiQLParser.g4:589-591) for LAG and LEAD, together with the mandatory
// `over` clause (g4:276-286). It is reached from parsePrimaryBase's single
// dispatch case for tokLAG/tokLEAD (exprprimary.go).
//
// LAG and LEAD are the ONLY window functions in the PartiQL grammar — there
// is no ROW_NUMBER / RANK / DENSE_RANK / NTILE (confirmed in the analysis
// doc, "Window Functions", and the grammar: the windowFunction rule lists
// only func=(LAG|LEAD)).
package parser

import (
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// parseWindowFuncExpr parses a LAG / LEAD window function call. The current
// token must be tokLAG or tokLEAD.
//
// Grammar (PartiQLParser.g4:589-591):
//
//	windowFunction
//	  : func=(LAG|LEAD) PAREN_LEFT expr ( COMMA expr (COMMA expr)? )? PAREN_RIGHT over # LagLeadFunction
//
//	over (g4:276-278)
//	  : OVER PAREN_LEFT windowPartitionList? windowSortSpecList? PAREN_RIGHT
//
//	windowPartitionList (g4:280-282)
//	  : PARTITION BY expr (COMMA expr)*
//
//	windowSortSpecList (g4:284-286)
//	  : ORDER BY orderSortSpec (COMMA orderSortSpec)*
//
// Two structural facts the generic FuncCall path cannot express, so this rule
// is hand-written rather than routed through parseFuncCallArgs:
//
//  1. ARITY 1-3. The argument list is `expr ( COMMA expr (COMMA expr)? )?` —
//     exactly one, two, or three expressions. The bare `LAG()` (zero args) and
//     `LAG(a, b, c, d)` (four+ args) forms are rejected. parseFuncCallArgs
//     would wrongly accept both (it permits any count, including zero).
//
//  2. OVER IS MANDATORY. The `over` clause is NOT optional in the
//     windowFunction rule (no `?` suffix), so `LAG(x)` without OVER is a
//     parse error. parseFuncCallArgs returns after the close paren and would
//     leave the missing OVER undetected at this level.
//
// All accept/reject verdicts below were confirmed against the legacy ANTLR
// parser (truth2, the runnable antlr_fallback oracle) by wrapping each form as
// a `SELECT <expr> FROM t` projection and feeding it to the generated
// PartiQLParser; see the PR body for the full oracle table.
//
//	LAG(x) OVER (ORDER BY y)                      -> {Name:LAG Args:[x] Over:{… OrderBy:[y]}}
//	LAG(x, 1) OVER (ORDER BY y)                   -> 2 args
//	LAG(x, 1, 0) OVER (PARTITION BY a ORDER BY b) -> 3 args, partition + order
//	LAG(x) OVER (PARTITION BY a)                  -> partition only (no ORDER BY)  [DIVERGENCE — see below]
//	LAG(x) OVER ()                                -> empty over                    [DIVERGENCE — see below]
//	LEAD(...) — identical in every respect to LAG.
//
//	LAG(x)                  REJECT (OVER is mandatory)
//	LAG()                   REJECT (>=1 arg required)
//	LAG(x, 1, 0, 9)         REJECT (<=3 args)
//	LAG(DISTINCT x) OVER …  REJECT (no setQuantifierStrategy in windowFunction)
//	LAG(*) OVER …           REJECT (no star; the arg is a plain expr)
//	OVER (ORDER BY)         REJECT (sort list needs >=1 spec)
//	OVER (PARTITION a …)    REJECT (PARTITION requires BY)
//	OVER (PARTITION BY)     REJECT (partition list needs >=1 expr)
//	OVER (ORDER BY y PARTITION BY a) REJECT (PARTITION must precede ORDER BY)
//	OVER ORDER BY y         REJECT (OVER requires its own parens)
//
// DIVERGENCE (parser vs semantic analyzer): the official partiql-lang-kotlin
// docs state ORDER BY is *required* for LAG/LEAD ("ORDER BY sub-clause has to
// be specified in order to use LAG function"). The executable ANTLR grammar —
// our authoritative parse-stage oracle — nonetheless ACCEPTS `OVER ()` and
// `OVER (PARTITION BY a)` with no ORDER BY, because `over` makes
// windowSortSpecList optional at the syntax level. The ORDER-BY requirement is
// therefore a *semantic* rule enforced after parsing, not a grammar rule. As a
// parser node we mirror the grammar (accept now, leave the semantic check to a
// later stage), exactly as parser-aggregates / parseExtractExpr treat their
// docs-level constraints as non-parse concerns. Flagged to the divergence
// ledger.
//
// The AST node is the generic *ast.FuncCall, distinguished by a non-nil Over
// field holding the *ast.WindowSpec (ast/exprs.go).
func (p *Parser) parseWindowFuncExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	name := strings.ToUpper(p.cur.Str)
	p.advance() // consume LAG / LEAD

	// Argument list: PAREN_LEFT expr ( COMMA expr (COMMA expr)? )? PAREN_RIGHT.
	// One to three general expressions. parseExprTop rejects `*` and a bare
	// `)`, so `LAG(*)` and `LAG()` reject here; a fourth argument trips the
	// PAREN_RIGHT expectation below, so `LAG(a,b,c,d)` rejects too.
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	args := make([]ast.ExprNode, 0, 3)
	for {
		arg, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.cur.Type != tokCOMMA || len(args) == 3 {
			break
		}
		p.advance() // consume COMMA
	}
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, err
	}

	// Mandatory OVER clause.
	over, end, err := p.parseOverClause()
	if err != nil {
		return nil, err
	}

	return &ast.FuncCall{
		Name: name,
		Args: args,
		Over: over,
		Loc:  ast.Loc{Start: start, End: end},
	}, nil
}

// parseOverClause parses the `over` rule (PartiQLParser.g4:276-278):
//
//	over : OVER PAREN_LEFT windowPartitionList? windowSortSpecList? PAREN_RIGHT
//
// Both sub-lists are optional and, when present, must appear in the fixed
// order PARTITION-BY then ORDER-BY (an ORDER BY followed by a PARTITION BY is a
// parse error — the ANTLR oracle rejects `OVER (ORDER BY y PARTITION BY a)`).
// `OVER ()` with neither list is accepted by the grammar.
//
// Returns the populated *ast.WindowSpec and the 0-based byte End offset of the
// closing paren.
func (p *Parser) parseOverClause() (*ast.WindowSpec, int, error) {
	if _, err := p.expect(tokOVER); err != nil {
		return nil, 0, err
	}
	start := p.prev.Loc.Start
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, 0, err
	}

	spec := &ast.WindowSpec{}

	// Optional windowPartitionList: PARTITION BY expr (COMMA expr)*.
	if p.cur.Type == tokPARTITION {
		partition, err := p.parseWindowPartitionList()
		if err != nil {
			return nil, 0, err
		}
		spec.PartitionBy = partition
	}

	// Optional windowSortSpecList: ORDER BY orderSortSpec (COMMA orderSortSpec)*.
	// Reuses the existing parseOrderByClause helper (select.go), which parses
	// the identical `ORDER BY orderSortSpec (COMMA orderSortSpec)*` production.
	if p.cur.Type == tokORDER {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, 0, err
		}
		spec.OrderBy = orderBy
	}

	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, 0, err
	}
	spec.Loc = ast.Loc{Start: start, End: rp.Loc.End}
	return spec, rp.Loc.End, nil
}

// parseWindowPartitionList parses the `windowPartitionList` rule
// (PartiQLParser.g4:280-282):
//
//	windowPartitionList : PARTITION BY expr (COMMA expr)*
//
// At least one expression is required (the leading `expr` is not optional), so
// `PARTITION BY` with no expression and a trailing `PARTITION BY a,` both
// reject. PARTITION must be immediately followed by BY.
func (p *Parser) parseWindowPartitionList() ([]ast.ExprNode, error) {
	if _, err := p.expect(tokPARTITION); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokBY); err != nil {
		return nil, err
	}
	var exprs []ast.ExprNode
	for {
		expr, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume COMMA
	}
	return exprs, nil
}
