package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `parser-select` DAG node (with select.go,
// relation.go, ctes.go and window.go): it implements Trino's `queryTerm` and
// `queryPrimary` grammar rules — the set-operation layer (UNION / INTERSECT /
// EXCEPT) and the four query primaries (a SELECT block, `TABLE name`,
// `VALUES …`, and a parenthesized `( queryNoWith )`).
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	queryTerm
//	    : queryPrimary                                              # queryTermDefault
//	    | left INTERSECT setQuantifier? right                       # setOperation
//	    | left (UNION|EXCEPT) setQuantifier? right                  # setOperation ;
//	queryPrimary
//	    : querySpecification          # queryPrimaryDefault
//	    | TABLE qualifiedName         # table
//	    | VALUES expression (, …)*    # inlineTable
//	    | ( queryNoWith )             # subquery ;
//
// Oracle-confirmed precedence & associativity (S1 in select.go), reflected in
// the two-level parse: INTERSECT binds TIGHTER than UNION/EXCEPT (ANTLR lists
// the INTERSECT alternative first), and within each level operators are
// LEFT-associative. So `a UNION b EXCEPT c` == `(a UNION b) EXCEPT c` and
// `a INTERSECT b UNION c` == `(a INTERSECT b) UNION c`. parseQueryTerm parses a
// UNION/EXCEPT chain whose operands are INTERSECT chains.
//
// Divergence #6 (D3 in select.go): the `CORRESPONDING` set-op modifier valid in
// Trino 481 is ABSENT from the legacy grammar (where the only modifier is
// setQuantifier = DISTINCT/ALL). It is a deferred P1 extension; not parsed here
// (omni rejects `UNION CORRESPONDING`, Trino 481 accepts it). Recorded in the
// migration divergence ledger.

// SetOperation is a binary set operation `left <op> [quantifier] right`
// (the setOperation queryTerm alternative). Op is "UNION", "INTERSECT", or
// "EXCEPT". Quantifier is "" (no modifier), "DISTINCT", or "ALL". Left and Right
// are the operand query terms.
type SetOperation struct {
	Op         string // "UNION", "INTERSECT", or "EXCEPT"
	Quantifier string // "", "DISTINCT", or "ALL"
	Left       QueryNode
	Right      QueryNode
	Loc        ast.Loc
}

func (n *SetOperation) Span() ast.Loc { return n.Loc }
func (*SetOperation) queryNode()      {}

// parseQueryTerm parses a `queryTerm`: a left-associative UNION/EXCEPT chain
// whose operands are INTERSECT chains (S1). This is the entry point for a query
// body, a parenthesized query, a CTE body, and a relation subquery.
func (p *Parser) parseQueryTerm() (QueryNode, error) {
	left, err := p.parseIntersectTerm()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwUNION || p.cur.Kind == kwEXCEPT {
		opTok := p.advance() // consume UNION / EXCEPT
		quant := p.parseSetQuantifier()
		right, err := p.parseIntersectTerm()
		if err != nil {
			return nil, err
		}
		left = &SetOperation{
			Op:         opTok.Str,
			Quantifier: quant,
			Left:       left,
			Right:      right,
			Loc:        ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseIntersectTerm parses the INTERSECT precedence level: a left-associative
// INTERSECT chain whose operands are query primaries. INTERSECT binds tighter
// than UNION/EXCEPT (S1).
func (p *Parser) parseIntersectTerm() (QueryNode, error) {
	left, err := p.parseQueryPrimary()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwINTERSECT {
		opTok := p.advance() // consume INTERSECT
		quant := p.parseSetQuantifier()
		right, err := p.parseQueryPrimary()
		if err != nil {
			return nil, err
		}
		left = &SetOperation{
			Op:         opTok.Str,
			Quantifier: quant,
			Left:       left,
			Right:      right,
			Loc:        ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseSetQuantifier consumes an optional `DISTINCT | ALL` set-operation
// modifier and returns its spelling ("" when absent). The CORRESPONDING modifier
// (Trino 481) is a deferred P1 extension and is intentionally NOT consumed here
// (D3) — leaving it surfaces as the syntax error a legacy-scope parser reports.
func (p *Parser) parseSetQuantifier() string {
	if tok, ok := p.match(kwDISTINCT, kwALL); ok {
		return tok.Str
	}
	return ""
}

// parseQueryPrimary parses a `queryPrimary`: a SELECT block, `TABLE
// qualifiedName`, `VALUES …`, or a parenthesized `( queryNoWith )`.
func (p *Parser) parseQueryPrimary() (QueryNode, error) {
	switch p.cur.Kind {
	case kwSELECT:
		return p.parseQuerySpecification()
	case kwTABLE:
		return p.parseTableQuery()
	case kwVALUES:
		return p.parseValuesQuery()
	case int('('):
		return p.parseParenQuery()
	default:
		return nil, p.queryPrimaryError()
	}
}

// parseTableQuery parses the `TABLE qualifiedName` query primary (TABLE is
// current): shorthand for `SELECT * FROM qualifiedName`.
func (p *Parser) parseTableQuery() (*TableQuery, error) {
	tableTok := p.advance() // consume TABLE
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &TableQuery{
		Name: name,
		Loc:  ast.Loc{Start: tableTok.Loc.Start, End: name.Loc.End},
	}, nil
}

// parseValuesQuery parses the `VALUES expression (, expression)*` query primary
// (VALUES is current): an inline row set. Each row is an expression — commonly a
// rowConstructor for a multi-column row (`VALUES (1, 'a'), (2, 'b')`), parsed by
// the expression layer.
func (p *Parser) parseValuesQuery() (*ValuesQuery, error) {
	valuesTok := p.advance() // consume VALUES
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	rows := []Expr{first}
	end := first.Span().End
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		rows = append(rows, next)
		end = next.Span().End
	}
	return &ValuesQuery{
		Rows: rows,
		Loc:  ast.Loc{Start: valuesTok.Loc.Start, End: end},
	}, nil
}

// parseParenQuery parses the `( queryNoWith )` query primary (the `subquery`
// rule; '(' is current). The inner content is a queryNoWith — a query without
// its own WITH but possibly carrying its own ORDER BY / LIMIT / set-operation.
// A WITH inside the parentheses is NOT a queryNoWith, so it is reported as a
// syntax error (Trino rejects `(WITH … SELECT …)` in this position; a CTE there
// must be a relation subquery `( query )`, handled in relation.go).
func (p *Parser) parseParenQuery() (*ParenQuery, error) {
	openTok := p.advance() // consume '('
	inner, err := p.parseQueryNoWith(nil, p.cur.Loc.Start)
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	inner.Loc.End = closeTok.Loc.Start
	return &ParenQuery{
		Inner: inner,
		Loc:   ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// queryPrimaryError reports a token that cannot begin a queryPrimary (SELECT /
// TABLE / VALUES / '(') where one is required.
func (p *Parser) queryPrimaryError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{Loc: p.cur.Loc, Msg: "expected a query (SELECT/TABLE/VALUES), found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "expected a query (SELECT/TABLE/VALUES), found " + text}
}
