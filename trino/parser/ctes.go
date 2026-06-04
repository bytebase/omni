package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `parser-select` DAG node (with select.go, setops.go,
// relation.go and window.go): it implements Trino's common-table-expression
// grammar — the `with` and `namedQuery` rules.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	with      : WITH RECURSIVE? namedQuery (, namedQuery)* ;
//	namedQuery: name=identifier columnAliases? AS ( query ) ;
//
// The WITH FUNCTION inline-routine prefix (`WITH functionSpecification …`,
// the withFunction rule) is the parser-routines DAG node, NOT a CTE; parseQuery
// (select.go) routes a `WITH FUNCTION …` away from this file. Here a leading
// WITH always introduces a CTE list.
//
// Divergence #4 (D2 in select.go): the `WITH SESSION prop = v <query>` prefix
// valid in Trino 481 is ABSENT from the legacy grammar (where WITH is CTE-only).
// It is a deferred P1 extension; not parsed here. A `WITH SESSION …` therefore
// fails as a malformed CTE (SESSION is not a valid namedQuery name followed by
// AS), which is the legacy-scope behavior. Recorded in the migration divergence
// ledger.

// WithClause is a `WITH [RECURSIVE] namedQuery (, namedQuery)*` clause (the with
// rule). Recursive marks the RECURSIVE keyword; CTEs is the non-empty list of
// named queries.
type WithClause struct {
	Recursive bool
	CTEs      []NamedQuery
	Loc       ast.Loc
}

// NamedQuery is one `name [columnAliases] AS ( query )` common table expression
// (the namedQuery rule). Name is the CTE name; ColumnAliases is the optional
// explicit column list; Query is the parsed CTE body.
type NamedQuery struct {
	Name          *ast.Identifier
	ColumnAliases []*ast.Identifier
	Query         *Query
	Loc           ast.Loc
}

// parseWithClause parses a `WITH [RECURSIVE] namedQuery (, namedQuery)*` clause
// (WITH is current; the caller has already confirmed it is not WITH FUNCTION).
func (p *Parser) parseWithClause() (*WithClause, error) {
	withTok := p.advance() // consume WITH
	wc := &WithClause{Loc: ast.Loc{Start: withTok.Loc.Start}}

	if _, ok := p.match(kwRECURSIVE); ok {
		wc.Recursive = true
	}

	first, err := p.parseNamedQuery()
	if err != nil {
		return nil, err
	}
	wc.CTEs = []NamedQuery{first}
	wc.Loc.End = first.Loc.End
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseNamedQuery()
		if err != nil {
			return nil, err
		}
		wc.CTEs = append(wc.CTEs, next)
		wc.Loc.End = next.Loc.End
	}
	return wc, nil
}

// parseNamedQuery parses one `name=identifier columnAliases? AS ( query )`
// (the namedQuery rule). The body is wrapped in mandatory parentheses and is a
// full `query` (so a CTE may itself have a WITH, set-operation, ORDER BY, etc.).
func (p *Parser) parseNamedQuery() (NamedQuery, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return NamedQuery{}, err
	}
	nq := NamedQuery{Name: name, Loc: name.Loc}

	if p.cur.Kind == int('(') {
		cols, _, err := p.parseColumnAliases()
		if err != nil {
			return NamedQuery{}, err
		}
		nq.ColumnAliases = cols
	}

	if _, err := p.expect(kwAS); err != nil {
		return NamedQuery{}, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return NamedQuery{}, err
	}
	q, err := p.parseQuery()
	if err != nil {
		return NamedQuery{}, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return NamedQuery{}, err
	}
	nq.Query = q
	nq.Loc.End = closeTok.Loc.End
	return nq, nil
}
