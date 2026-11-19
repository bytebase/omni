package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads a derived analytic view: the inline analytic view in a
// FROM clause and the subav_factoring_clause of a WITH clause. Both take a
// sub_av_clause, which names a persistent analytic view and the hierarchies,
// filters and measures the query reads of it.
//
// BNF: oracle/parser/bnf/SELECT.bnf (sub_av_clause, hierarchies_clause,
// filter_clauses, add_meas_clause, inline_analytic_view,
// subav_factoring_clause),
// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SELECT.html
//
// The pinned Oracle 23.26.3 reads the clause as the page's diagrams draw it
// and the page's examples write it: sub_av_clause starts at USING, and
// ANALYTIC VIEW opens the inline view and the WITH element, not the clause.
// SELECT * FROM ANALYTIC VIEW(ANALYTIC VIEW USING a) is ORA-02000, "missing
// USING keyword", at the token inside the parens (byte 28), while
// ANALYTIC VIEW(USING a) is read and its name resolved (ORA-18307), and
// WITH a ANALYTIC VIEW AS(USING a) is read too (ORA-18385 there, a check of
// the parsed statement).

// isInlineAnalyticViewAhead reports whether the tokens at the cursor open a
// derived analytic view: ANALYTIC VIEW, which no table reference starts
// with.
func (p *Parser) isInlineAnalyticViewAhead() bool {
	return p.isKeywordStr("ANALYTIC") && p.peekNext().Type == kwVIEW
}

// parseInlineAnalyticViewRef reads ANALYTIC VIEW ( sub_av_clause ) with its
// optional alias.
//
//	inline_analytic_view:
//	    ANALYTIC VIEW ( sub_av_clause ) [ inline_av_alias ]
func (p *Parser) parseInlineAnalyticViewRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume ANALYTIC
	p.advance() // consume VIEW
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume (

	clause, err := p.parseSubAVClause()
	if err != nil {
		return nil, err
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}

	ref := &nodes.InlineAnalyticView{
		Clause: clause,
		Loc:    nodes.Loc{Start: start},
	}
	if p.cur.Type == kwAS {
		p.advance()
		if ref.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	} else if p.isTableAliasCandidate() {
		if ref.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	}
	ref.Loc.End = p.prev.End
	return p.parsePivotSuffix(ref)
}

// parseSubAVClause reads a sub_av_clause, the clause in an inline analytic
// view's parentheses and in a WITH subav_factoring_clause.
//
//	sub_av_clause:
//	    USING [ schema. ] base_av_name
//	    [ hierarchies_clause ]
//	    [ filter_clauses ]
//	    [ add_meas_clause ]
func (p *Parser) parseSubAVClause() (*nodes.SubAVClause, error) {
	start := p.pos()
	clause := &nodes.SubAVClause{Loc: nodes.Loc{Start: start}}

	if err := p.anExpect(kwUSING); err != nil {
		return nil, err
	}
	var err error
	if clause.Using, err = p.anParseObjectName(); err != nil {
		return nil, err
	}

	if p.isKeywordStr("HIERARCHIES") {
		if err := p.parseSubAVHierarchies(clause); err != nil {
			return nil, err
		}
	}
	if p.isKeywordStr("FILTER") {
		if err := p.parseSubAVFilters(clause); err != nil {
			return nil, err
		}
	}
	if p.isKeywordStr("ADD") && p.anNextIsWord("MEASURES") {
		if err := p.parseSubAVAddMeasures(clause); err != nil {
			return nil, err
		}
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseSubAVHierarchies reads hierarchies_clause, whose list the page draws
// as entirely optional:
//
//	HIERARCHIES ( [ [ attr_dim_alias. ] hier_alias [, hier_alias ]... ] )
func (p *Parser) parseSubAVHierarchies(clause *nodes.SubAVClause) error {
	optStart := p.pos()
	p.advance() // consume HIERARCHIES
	if err := p.anExpect('('); err != nil {
		return err
	}
	if p.cur.Type != ')' {
		for {
			if _, err := p.anParseDottedName(2); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if err := p.anExpect(')'); err != nil {
		return err
	}
	clause.Options = anAppendItems(clause.Options, &nodes.List{Items: []nodes.Node{p.anOption("HIERARCHIES", optStart)}})
	return nil
}

// parseSubAVFilters reads filter_clauses: one FILTER FACT clause around a
// comma-separated list of filter_clause.
//
//	FILTER FACT ( filter_clause [, filter_clause ]... )
//	filter_clause: hier_ids TO condition
//	hier_ids: hier_id [, hier_id ]...
func (p *Parser) parseSubAVFilters(clause *nodes.SubAVClause) error {
	optStart := p.pos()
	p.advance() // consume FILTER
	if err := p.anExpectWord("FACT"); err != nil {
		return err
	}
	if err := p.anExpect('('); err != nil {
		return err
	}
	for {
		// hier_ids: a comma between names belongs to the group, since a
		// filter_clause's own comma follows its condition.
		for {
			if _, err := p.anParseDottedName(2); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.anExpect(kwTO); err != nil {
			return err
		}
		exprStart := p.pos()
		expr, parseErr := p.parseExpr()
		_ = expr
		if parseErr != nil {
			return parseErr
		}
		if p.pos() == exprStart {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return err
	}
	clause.Options = anAppendItems(clause.Options, &nodes.List{Items: []nodes.Node{p.anOption("FILTER FACT", optStart)}})
	return nil
}

// parseSubAVAddMeasures reads add_meas_clause, whose items are cube_meas:
// a measure name and a base or calculated measure clause.
//
//	ADD MEASURES ( cube_meas [, cube_meas ]... )
//	cube_meas: meas_name { base_meas_clause | calc_meas_clause }
//	base_meas_clause: [ FACT column_name ] [ aggr_function ]
//	calc_meas_clause: AS ( expression )
func (p *Parser) parseSubAVAddMeasures(clause *nodes.SubAVClause) error {
	optStart := p.pos()
	p.advance() // consume ADD
	p.advance() // consume MEASURES
	if err := p.anExpect('('); err != nil {
		return err
	}
	for {
		if _, err := p.anParseName(); err != nil {
			return err
		}
		switch {
		case p.cur.Type == kwAS:
			p.advance()
			// The calculated measure's expression has a grammar of its own:
			// its parentheses are read whole, and must not be empty.
			if err := p.skipAnalyticViewExpression(); err != nil {
				return err
			}
		default:
			if p.isKeywordStr("FACT") {
				p.advance()
				if _, err := p.anParseName(); err != nil {
					return err
				}
			}
			if p.cur.Type != tokQIDENT && p.isIdentLike() && analyticAggrFunctions[p.cur.Str] {
				aggr, parseErr := p.parseAnalyticAggrFunction()
				_ = aggr
				if parseErr != nil {
					return parseErr
				}
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return err
	}
	clause.Options = anAppendItems(clause.Options, &nodes.List{Items: []nodes.Node{p.anOption("ADD MEASURES", optStart)}})
	return nil
}
