package parser

import "github.com/bytebase/omni/trino/ast"

// This file implements rootQueryWithSession, the `WITH SESSION name = value …
// <query>` statement prefix (statementDefault). It is only a statement prefix:
// a subquery's WITH is the CTE form.

// SessionProperty is one `qualifiedName = expression` of a WITH SESSION prefix.
type SessionProperty struct {
	Name  *ast.QualifiedName
	Value Expr
	Loc   ast.Loc
}

// QueryWithSessionStmt is `WITH SESSION prop (, prop)* rootQuery`.
type QueryWithSessionStmt struct {
	Properties []SessionProperty
	Query      *Query
	Loc        ast.Loc
}

func (n *QueryWithSessionStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *QueryWithSessionStmt) Span() ast.Loc    { return n.Loc }

var _ ast.Node = (*QueryWithSessionStmt)(nil)

// parseQueryWithSessionStmt parses `WITH SESSION prop (, prop)* query`. WITH is
// current and SESSION is the next token.
func (p *Parser) parseQueryWithSessionStmt() (ast.Node, error) {
	withTok := p.advance() // consume WITH
	p.advance()            // consume SESSION
	s := &QueryWithSessionStmt{Loc: ast.Loc{Start: withTok.Loc.Start}}
	first, err := p.parseSessionProperty()
	if err != nil {
		return nil, err
	}
	s.Properties = append(s.Properties, first)
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSessionProperty()
		if err != nil {
			return nil, err
		}
		s.Properties = append(s.Properties, next)
	}
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	s.Query = q
	s.Loc.End = q.Loc.End
	return s, nil
}

// parseSessionProperty parses `qualifiedName = expression`.
func (p *Parser) parseSessionProperty() (SessionProperty, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return SessionProperty{}, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return SessionProperty{}, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return SessionProperty{}, err
	}
	return SessionProperty{
		Name:  name,
		Value: value,
		Loc:   ast.Loc{Start: name.Loc.Start, End: value.Span().End},
	}, nil
}
