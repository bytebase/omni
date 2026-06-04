package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-select` DAG node. It implements GoogleSQL's
// WITH / CTE clause (GoogleSQLParser.g4 §2.13 with_clause / aliased_query /
// recursion_depth_modifier), a hand-port of Google's open-source ZetaSQL
// reference. The WITH clause is parsed by parseWithClause, called from
// parseQuery (select.go) when a query begins with WITH.
//
// Grammar:
//
//	with_clause: WITH [RECURSIVE] aliased_query (, aliased_query)*
//	aliased_query: identifier [( column, … )] AS ( query ) [recursion_depth_modifier]
//	recursion_depth_modifier:
//	  WITH DEPTH [AS alias] [BETWEEN <bound> AND <bound> | MAX <bound>]
//
// RECURSIVE parses in the union grammar (BigQuery supports recursive CTEs);
// Spanner feature-rejects it ("RECURSIVE is not supported in the WITH clause" —
// divergence ledger #11), which is a feature-reject, NOT a grammar verdict, so
// the parser accepts RECURSIVE. The Spanner-only `cte_name ( column, … )`
// explicit-column form also parses in the union grammar.
//
// Disambiguation from the inline WITH(...) expression: a query-leading WITH is
// always `WITH [RECURSIVE] identifier …` (an aliased-query list); the inline
// WITH expression is `WITH ( id AS expr, …, body )` and is owned by the
// expressions node in expression position. parseQuery only reaches parseWithClause
// at statement/query head, never in expression position, so there is no conflict.

// parseWithClause parses a with_clause. WITH is the current token. It returns
// the WithClause; the following query body is parsed by the caller (parseQuery).
func (p *Parser) parseWithClause() (*ast.WithClause, error) {
	withTok := p.advance() // WITH
	wc := &ast.WithClause{Loc: withTok.Loc}

	if p.cur.Type == kwRECURSIVE {
		p.advance()
		wc.Recursive = true
	}

	for {
		cte, err := p.parseCTE()
		if err != nil {
			return nil, err
		}
		wc.CTEs = append(wc.CTEs, cte)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	wc.Loc.End = p.prev.Loc.End
	return wc, nil
}

// parseCTE parses one aliased_query:
//
//	identifier [( column, … )] AS ( query ) [recursion_depth_modifier]
//
// The current token begins the CTE name.
func (p *Parser) parseCTE() (*ast.CTE, error) {
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	cte := &ast.CTE{Name: name, Loc: nameTok.Loc}

	// Optional explicit column list `( column, … )` (Spanner form). Detect it by
	// a '(' that is NOT the `AS (query)` parenthesis: here the '(' comes BEFORE
	// the AS keyword.
	if p.cur.Type == int('(') {
		p.advance() // '('
		for {
			colTok, err := p.expectIdentifier()
			if err != nil {
				return nil, err
			}
			col, err := p.identifierText(colTok)
			if err != nil {
				return nil, err
			}
			cte.Columns = append(cte.Columns, col)
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	cte.Query = query
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	cte.Loc.End = closeTok.Loc.End

	// Optional recursion-depth modifier: WITH DEPTH ….
	if p.cur.Type == kwWITH && p.peekNext().Type == kwDEPTH {
		depth, err := p.parseRecursionDepth()
		if err != nil {
			return nil, err
		}
		cte.Depth = depth
		cte.Loc.End = depth.Loc.End
	}

	return cte, nil
}

// parseRecursionDepth parses a recursion_depth_modifier:
//
//	WITH DEPTH [AS alias] [BETWEEN <bound> AND <bound> | MAX <bound>]
//
// WITH is the current token (DEPTH confirmed by caller). Each bound is an
// integer / parameter / system-variable / UNBOUNDED
// (possibly_unbounded_int_literal_or_parameter); we accept any expression there
// for parse parity (the bound's exact shape is not consumed by bytebase).
func (p *Parser) parseRecursionDepth() (*ast.RecursionDepth, error) {
	withTok := p.advance()  // WITH
	depthTok := p.advance() // DEPTH
	rd := &ast.RecursionDepth{Loc: ast.Loc{Start: withTok.Loc.Start, End: depthTok.Loc.End}}

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
		rd.Alias = alias
		rd.Loc.End = aliasTok.Loc.End
	}

	switch p.cur.Type {
	case kwBETWEEN:
		p.advance()
		low, err := p.parseDepthBound()
		if err != nil {
			return nil, err
		}
		rd.Lower = low
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		high, err := p.parseDepthBound()
		if err != nil {
			return nil, err
		}
		rd.Upper = high
		rd.Loc.End = nodeLoc(high).End
	case kwMAX:
		p.advance()
		max, err := p.parseDepthBound()
		if err != nil {
			return nil, err
		}
		rd.Max = max
		rd.Loc.End = nodeLoc(max).End
	}
	return rd, nil
}

// parseDepthBound parses one recursion-depth bound
// (possibly_unbounded_int_literal_or_parameter): an integer literal, a `@param`
// / `?` parameter, a `@@sysvar` system variable, or the UNBOUNDED keyword. It is
// deliberately NOT parseExpr: in `BETWEEN low AND high` the AND is the BETWEEN
// separator, and a general expression parse of `low` would greedily swallow
// `AND high` as a logical conjunction. UNBOUNDED is modeled as an Identifier so
// the bound is always a Node.
func (p *Parser) parseDepthBound() (ast.Node, error) {
	switch p.cur.Type {
	case kwUNBOUNDED:
		tok := p.advance()
		return &ast.Identifier{Name: "UNBOUNDED", Loc: tok.Loc}, nil
	case tokInteger:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}, nil
	case int('@'):
		return p.parseNamedParameter()
	case int('?'):
		tok := p.advance()
		return &ast.Parameter{Positional: true, Loc: tok.Loc}, nil
	case tokAtAt:
		return p.parseSystemVariable()
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
